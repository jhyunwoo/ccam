package usage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func window(p float64, reset time.Time) Window {
	return Window{Percent: p, ResetsAt: reset, Known: true}
}

// A window Claude Code refused a request in counts as full until its reset,
// whatever numbers measured up to then say and however they were found
// (a lookup at the limit's own time is not after it).
func TestLimitRecordCountsAsFull(t *testing.T) {
	ctx := context.Background()
	now := time.Now().Truncate(time.Second)
	reset := now.Add(3 * time.Hour)
	full := func(t *testing.T, r Report) {
		t.Helper()
		if l := r.Limits; l == nil || !l.FiveHour.Known || l.FiveHour.Percent != 100 || !l.FiveHour.ResetsAt.Equal(reset) {
			t.Fatalf("not full until the reset: %+v", r)
		}
		if (Policy{}).Choose([]*Report{&r}, now) != nil || !(Policy{}).FreeAt([]*Report{&r}, now).Equal(reset) {
			t.Fatalf("choice: %+v", r.Limits)
		}
	}
	t.Run("live cache 20 s old at 40%", func(t *testing.T) {
		e := newEnv(t)
		e.f.Now = func() time.Time { return now }
		a := e.account("a", "ok-new-format.json", future)
		e.f.writeCache("a", "live", cached{At: now.Add(-20 * time.Second), Limits: &Limits{FiveHour: window(40, now.Add(time.Hour))}})
		if err := SaveLimit(e.f.StateDir, "a", "five_hour", reset, now); err != nil {
			t.Fatal(err)
		}
		r := e.f.Fetch(ctx, a, Options{})
		if e.hits.Load() != 0 || r.Source != "live" {
			t.Fatalf("the cache within 60 s: %+v hits=%d", r, e.hits.Load())
		}
		full(t, r)
		r = e.f.Fetch(ctx, a, Options{NoCache: true}) // a live answer, too
		if e.hits.Load() != 1 {
			t.Fatal("--no-cache must ask")
		}
		full(t, r)
	})
	t.Run("expired token and a 3-hour-old cache at 10%", func(t *testing.T) {
		e := newEnv(t)
		e.f.Now = func() time.Time { return now }
		a := e.account("a", "ok-new-format.json", now.Add(-time.Hour).UnixMilli())
		e.f.writeCache("a", "live", cached{At: now.Add(-3 * time.Hour), Limits: &Limits{FiveHour: window(10, now.Add(time.Hour))}})
		SaveLimit(e.f.StateDir, "a", "five_hour", reset, now)
		r := e.f.Fetch(ctx, a, Options{})
		if r.Status != StatusExpired || r.Source != "cache" || e.hits.Load() != 0 {
			t.Fatalf("%+v", r)
		}
		full(t, r)
	})
	t.Run("saved values only, status line", func(t *testing.T) {
		e := newEnv(t)
		e.f.Now = func() time.Time { return now }
		a := e.account("a", "boom", future)
		SaveStatusLine(e.f.StateDir, "a", Limits{FiveHour: window(23, now.Add(time.Hour))}, now)
		SaveLimit(e.f.StateDir, "a", "five_hour", reset, now)
		for _, opt := range []Options{{Source: "cache"}, {}} {
			full(t, e.f.Fetch(ctx, a, opt))
		}
	})
	t.Run("Claude Code's lookup", func(t *testing.T) {
		e := newEnv(t)
		e.f.Now = func() time.Time { return now }
		e.f.Source = "claude"
		e.f.Probe = func(context.Context, Account) ([]byte, error) {
			return []byte(`{"type":"control_response","response":{"subtype":"success","request_id":"ccam","response":{"rate_limits_available":true,` +
				`"rate_limits":{"five_hour":{"utilization":30,"resets_at":"` + now.Add(time.Hour).UTC().Format(time.RFC3339) + `"}}}}}`), nil
		}
		a := e.account("a", "ok-new-format.json", future)
		SaveLimit(e.f.StateDir, "a", "five_hour", reset, now)
		full(t, e.f.Fetch(ctx, a, Options{}))
		full(t, e.f.Fetch(ctx, a, Options{})) // its 60 s cache
	})
	t.Run("nothing saved", func(t *testing.T) {
		e := newEnv(t)
		e.f.Now = func() time.Time { return now }
		a := e.account("a", "boom", future)
		SaveLimit(e.f.StateDir, "a", "five_hour", reset, now.Add(-time.Minute))
		r := e.f.Fetch(ctx, a, Options{})
		full(t, r)
		if r.Status != StatusError || r.Source != "limit" || !r.At.Equal(now.Add(-time.Minute)) || r.Limits.Weekly.Known {
			t.Fatalf("%+v", r)
		}
	})
	t.Run("no subscription", func(t *testing.T) {
		e := newEnv(t)
		e.f.Now = func() time.Time { return now }
		out := e.account("out", "", 0)
		api := e.account("api", "", 0)
		os.WriteFile(filepath.Join(api.Dir, ".claude.json"), []byte(`{"primaryApiKey":"sk-x"}`), 0o600)
		for _, a := range []Account{out, api} {
			SaveLimit(e.f.StateDir, a.Name, "five_hour", reset, now)
			if r := e.f.Fetch(ctx, a, Options{}); r.Limits != nil {
				t.Fatalf("%s: %+v", a.Name, r)
			}
		}
	})
}

// A lookup of the usage endpoint made liftDelay or more after the limit
// that shows the window below 100% means the limit was lifted before its
// reset (Claude Code's /limit-reset): the record no longer counts for that
// window. Numbers measured before then, and the status line, do not lift it.
func TestLimitRecordLiftedByALaterLookup(t *testing.T) {
	ctx := context.Background()
	now := time.Now().Truncate(time.Second)
	limitAt := now.Add(-26 * time.Hour)
	weekReset := now.Add(18 * time.Hour)
	week := func(p float64) *Limits {
		return &Limits{FiveHour: window(0, time.Time{}), Weekly: window(p, weekReset),
			Scoped: []Scoped{{Name: "Fable", Window: window(34, weekReset)}}}
	}
	setup := func(t *testing.T) (*env, Account) {
		e := newEnv(t)
		e.f.Now = func() time.Time { return now }
		if err := SaveLimit(e.f.StateDir, "a", "seven_day", weekReset, limitAt); err != nil {
			t.Fatal(err)
		}
		return e, e.account("a", "boom", future)
	}
	weekly := func(t *testing.T, r Report, want float64, limited bool) {
		t.Helper()
		if r.Limits == nil || r.Limits.Weekly.Percent != want || !r.Limits.Weekly.ResetsAt.Equal(weekReset) {
			t.Fatalf("weekly: %+v", r.Limits)
		}
		if got := len(r.LimitedWindows) > 0; got != limited {
			t.Fatalf("limited %v: %+v", got, r)
		}
	}
	for _, src := range []string{"live", "claude"} {
		t.Run(src+" after the limit", func(t *testing.T) {
			e, a := setup(t)
			e.f.writeCache("a", src, cached{At: limitAt.Add(liftDelay), Limits: week(17)})
			for _, opt := range []Options{{Source: "cache"}, {}} {
				r := e.f.Fetch(ctx, a, opt)
				weekly(t, r, 17, false)
				if (Policy{}).Full(&r) {
					t.Fatal("chosen as full")
				}
			}
		})
	}
	t.Run("before, at or right after the limit", func(t *testing.T) {
		for _, at := range []time.Time{limitAt.Add(-time.Hour), limitAt, limitAt.Add(liftDelay - time.Second)} {
			e, a := setup(t)
			e.f.writeCache("a", "live", cached{At: at, Limits: week(17)})
			weekly(t, e.f.Fetch(ctx, a, Options{Source: "cache"}), 100, true)
		}
	})
	t.Run("still full after the limit", func(t *testing.T) {
		e, a := setup(t)
		e.f.writeCache("a", "live", cached{At: now.Add(-time.Hour), Limits: week(100)})
		weekly(t, e.f.Fetch(ctx, a, Options{Source: "cache"}), 100, true)
	})
	t.Run("status line only", func(t *testing.T) {
		// Its time is when the command ran, not when Claude Code measured.
		e, a := setup(t)
		e.f.writeCache("a", "live", cached{At: limitAt.Add(-time.Hour), Limits: week(90)})
		SaveStatusLine(e.f.StateDir, "a", Limits{FiveHour: window(5, now.Add(time.Hour)), Weekly: window(17, weekReset)}, now.Add(-time.Minute))
		weekly(t, e.f.Fetch(ctx, a, Options{Source: "cache"}), 100, true)
	})
	t.Run("a later lookup lifts it under a newer status line", func(t *testing.T) {
		e, a := setup(t)
		e.f.writeCache("a", "live", cached{At: limitAt.Add(time.Hour), Limits: week(17)})
		SaveStatusLine(e.f.StateDir, "a", Limits{FiveHour: window(5, now.Add(time.Hour)), Weekly: window(20, weekReset)}, now.Add(-time.Minute))
		r := e.f.Fetch(ctx, a, Options{Source: "cache"})
		weekly(t, r, 20, false)
	})
	t.Run("the newest lookup decides", func(t *testing.T) {
		e, a := setup(t)
		e.f.writeCache("a", "live", cached{At: limitAt.Add(time.Hour), Limits: week(17)})
		e.f.writeCache("a", "claude", cached{At: limitAt.Add(2 * time.Hour), Limits: week(100)})
		weekly(t, e.f.Fetch(ctx, a, Options{Source: "cache"}), 100, true)
	})
	t.Run("a new limit after the lookup", func(t *testing.T) {
		e, a := setup(t)
		e.f.writeCache("a", "live", cached{At: limitAt.Add(time.Hour), Limits: week(17)})
		SaveLimit(e.f.StateDir, "a", "seven_day", weekReset, now.Add(-time.Minute))
		weekly(t, e.f.Fetch(ctx, a, Options{Source: "cache"}), 100, true)
	})
	t.Run("only the windows the lookup shows below 100%", func(t *testing.T) {
		e, a := setup(t)
		SaveLimit(e.f.StateDir, "a", "five_hour", now.Add(time.Hour), limitAt)
		l := week(17)
		l.FiveHour = window(100, now.Add(time.Hour))
		e.f.writeCache("a", "live", cached{At: limitAt.Add(time.Hour), Limits: l})
		r := e.f.Fetch(ctx, a, Options{Source: "cache"})
		weekly(t, r, 17, true)
		if r.Limits.FiveHour.Percent != 100 || len(r.LimitedWindows) != 1 || r.LimitedWindows[0] != "5-hour" {
			t.Fatalf("%+v", r)
		}
	})
}

// The Fable allowance (seven_day_overage_included, Claude Code's "Fable
// limit") is the Fable week: its record fills that window only, so a Fable
// conversation goes elsewhere and an Opus one may still come.
func TestLimitRecordFillsTheFableWeek(t *testing.T) {
	e := newEnv(t)
	now := time.Now().Truncate(time.Second)
	e.f.Now = func() time.Time { return now }
	a := e.account("a", "boom", future)
	e.f.writeCache("a", "live", cached{At: now.Add(-time.Hour), Limits: &Limits{FiveHour: window(10, now.Add(time.Hour)),
		Weekly: window(40, now.Add(48*time.Hour)), Scoped: []Scoped{{Name: "Fable", Window: window(97, now.Add(48*time.Hour))}}}})
	if err := SaveLimit(e.f.StateDir, "a", "seven_day_overage_included", now.Add(48*time.Hour), now.Add(-time.Minute)); err != nil {
		t.Fatal(err)
	}
	r := e.f.Fetch(context.Background(), a, Options{Source: "cache"})
	s := r.Limits.Scoped
	if len(s) != 1 || s[0].Window.Percent != 100 || r.Limits.Weekly.Percent != 40 || len(r.LimitedWindows) != 1 || r.LimitedWindows[0] != "Fable" {
		t.Fatalf("%+v %+v", r.Limits, r.LimitedWindows)
	}
	if !(Policy{Model: "claude-fable-5-2"}).Full(&r) || (Policy{Model: "claude-opus-5-5"}).Full(&r) || (Policy{}).Full(&r) {
		t.Fatal("full by model")
	}
	for w, fam := range map[string]string{"seven_day_overage_included": "fable", "seven_day_opus": "opus", "seven_day_sonnet": "sonnet", "seven_day": "", "five_hour": "", "overage": ""} {
		if got := WindowModel(w); got != fam {
			t.Errorf("WindowModel(%s) = %q, want %q", w, got, fam)
		}
	}
}

// Opus and Sonnet records fill the model's own window: the one whose name
// has the family in it, else a new one.
func TestLimitRecordFillsModelWindows(t *testing.T) {
	e := newEnv(t)
	now := time.Now().Truncate(time.Second)
	e.f.Now = func() time.Time { return now }
	a := e.account("a", "boom", future)
	e.f.writeCache("a", "live", cached{At: now.Add(-time.Hour), Limits: &Limits{FiveHour: window(10, now.Add(time.Hour)),
		Scoped: []Scoped{{Name: "Fable", Window: window(30, now.Add(48*time.Hour))}, {Name: "Claude Opus", Window: window(20, now.Add(48*time.Hour))}}}})
	SaveLimit(e.f.StateDir, "a", "seven_day_opus", now.Add(50*time.Hour), now)
	SaveLimit(e.f.StateDir, "a", "seven_day_sonnet", now.Add(60*time.Hour), now)
	r := e.f.Fetch(context.Background(), a, Options{Source: "cache"})
	s := r.Limits.Scoped
	if len(s) != 3 || s[0].Name != "Fable" || s[0].Window.Percent != 30 || s[1].Name != "Claude Opus" || s[1].Window.Percent != 100 ||
		!s[1].Window.ResetsAt.Equal(now.Add(50*time.Hour)) || s[2].Name != "Sonnet" || s[2].Window.Percent != 100 {
		t.Fatalf("%+v", s)
	}
	if !(Policy{Model: "opus"}).Full(&r) || !(Policy{Model: "claude-sonnet-5-5"}).Full(&r) || (Policy{Model: "fable"}).Full(&r) || (Policy{}).Full(&r) {
		t.Fatal("full by model")
	}
	// The saved cache itself is left as it was.
	if c, _ := e.f.readCache("a", "live"); c.Limits.Scoped[1].Window.Percent != 20 || len(c.Limits.Scoped) != 2 {
		t.Fatalf("cache changed: %+v", c.Limits)
	}
}

// A record whose window has reset no longer counts, and the next record
// drops it.
func TestLimitRecordExpires(t *testing.T) {
	e := newEnv(t)
	now := time.Now().Truncate(time.Second)
	a := e.account("a", "boom", future)
	e.f.writeCache("a", "live", cached{At: now, Limits: &Limits{FiveHour: window(40, time.Time{})}})
	SaveLimit(e.f.StateDir, "a", "five_hour", now.Add(time.Hour), now)
	later := now.Add(2 * time.Hour)
	e.f.Now = func() time.Time { return later }
	if r := e.f.Fetch(context.Background(), a, Options{Source: "cache"}); r.Limits.FiveHour.Percent != 40 {
		t.Fatalf("an expired record: %+v", r.Limits)
	}
	if err := SaveLimit(e.f.StateDir, "a", "seven_day", later.Add(48*time.Hour), later); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(e.f.StateDir, "a.limit.json")
	data, _ := os.ReadFile(path)
	var rec struct{ Windows map[string]time.Time }
	if json.Unmarshal(data, &rec) != nil || len(rec.Windows) != 1 || !rec.Windows["seven_day"].Equal(later.Add(48*time.Hour)) {
		t.Fatalf("%s", data)
	}
	if info, _ := os.Stat(path); runtime.GOOS != "windows" && info.Mode().Perm() != 0o600 {
		t.Fatal(info.Mode())
	}
	// Windows ccam does not choose by, and limits that are over, are not kept.
	if err := SaveLimit(e.f.StateDir, "a", "overage", later.Add(time.Hour), later); err == nil {
		t.Fatal("overage recorded")
	}
	if err := SaveLimit(e.f.StateDir, "b", "five_hour", later.Add(-time.Minute), later); err != nil || fileExists(filepath.Join(e.f.StateDir, "b.limit.json")) {
		t.Fatal("a limit that is over was recorded")
	}
	found := false
	for _, p := range CachePaths(e.f.StateDir, "a") {
		found = found || p == path
	}
	if !found {
		t.Fatal("CachePaths must list the record, for rename and remove")
	}
}

func fileExists(path string) bool { _, err := os.Stat(path); return err == nil }

// BenchmarkFetchCache is `ccam usage --source cache` for five accounts:
// saved values only, without and with a limit record.
func BenchmarkFetchCache(b *testing.B) {
	for _, limit := range []bool{false, true} {
		b.Run(map[bool]string{false: "no-limit", true: "limit"}[limit], func(b *testing.B) {
			root := b.TempDir()
			f := NewFetcher(filepath.Join(root, "state"))
			now := time.Now()
			var accts []Account
			for i := range 5 {
				name := fmt.Sprintf("p%d", i)
				accts = append(accts, Account{Name: name, Slot: i + 1, Dir: filepath.Join(root, name)})
				f.writeCache(name, "live", cached{At: now.Add(-time.Minute), Limits: &Limits{FiveHour: window(20, now.Add(time.Hour)),
					Weekly: window(30, now.Add(48*time.Hour)), Scoped: []Scoped{{Name: "Fable", Window: window(10, now.Add(48*time.Hour))}}}})
				if limit {
					if err := SaveLimit(f.StateDir, name, "five_hour", now.Add(time.Hour), now); err != nil {
						b.Fatal(err)
					}
				}
			}
			for b.Loop() {
				f.FetchAll(context.Background(), accts, Options{Source: "cache"}, nil)
			}
		})
	}
}

// A limit record laid over a lookup says so: in --json (limited) and on
// the -v line, which otherwise names the lookup's source for a 100% it
// never reported.
func TestLimitOverALookupSaysSo(t *testing.T) {
	e := newEnv(t)
	now := time.Now().Truncate(time.Second)
	e.f.Now = func() time.Time { return now }
	a := e.account("a", "boom", future)
	// Looked up a minute before the limit.
	e.f.writeCache("a", "live", cached{At: now.Add(-6 * time.Minute), Limits: &Limits{FiveHour: window(72, now.Add(time.Hour))}})
	if err := SaveLimit(e.f.StateDir, "a", "five_hour", now.Add(2*time.Hour), now.Add(-5*time.Minute)); err != nil {
		t.Fatal(err)
	}
	r := e.f.Fetch(context.Background(), a, Options{Source: "cache"})
	if r.Source != "cache" || !r.Limited.Equal(now.Add(-5*time.Minute)) {
		t.Fatalf("%+v", r)
	}
	out := fmt.Sprint(Table([]*Report{&r}, RenderOptions{Verbose: true, Now: now}))
	if !strings.Contains(out, "5-hour full: Claude Code reported the usage limit 5m ago") {
		t.Fatalf("%s", out)
	}
	if data, _ := json.Marshal(r); !strings.Contains(string(data), `"limited":`) {
		t.Fatalf("json: %s", data)
	}
}
