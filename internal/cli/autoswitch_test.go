package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/session"
	"github.com/jhyunwoo/ccam/internal/transcript"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// auto_switch moves a session only at a limit of the account's own plan
// (a Pro or Max plan, or a seat: D7 b) that resets more than nearReset
// later, under a relay of this ccam that runs that very Claude Code, and
// not while Remote Control is connected (D11).
func TestEligible(t *testing.T) {
	base := limitFacts{on: true, kind: handoff.Account, resetsIn: 3 * time.Hour, relayed: true, child: true,
		plan: &config.Identity{OrgUUID: "o1", OrgType: "claude_max"}}
	for _, c := range []struct {
		name string
		edit func(*limitFacts)
		want stay
	}{
		{"a Max plan's 5-hour limit", nil, moves},
		{"auto_switch off, or Windows", func(f *limitFacts) { f.on = false }, stayOff},
		{"a spend limit", func(f *limitFacts) { f.kind = handoff.Credits }, stayOff},
		{"an organization's stop", func(f *limitFacts) { f.kind = handoff.Org }, stayOff},
		{"no usage limit", func(f *limitFacts) { f.kind = "" }, stayOff},
		{"an Opus limit", func(f *limitFacts) { f.kind = handoff.Model }, moves},
		{"no reset known", func(f *limitFacts) { f.resetsIn = 0 }, stayOff},
		{"a Pro plan", func(f *limitFacts) { f.plan = &config.Identity{OrgUUID: "o2", OrgType: "claude_pro"} }, moves},
		{"a Team seat", func(f *limitFacts) { f.plan = &config.Identity{OrgUUID: "acme", OrgType: "claude_team"} }, moves},
		{"an Enterprise seat", func(f *limitFacts) { f.plan = &config.Identity{OrgUUID: "big", OrgType: "claude_enterprise"} }, moves},
		{"nobody known", func(f *limitFacts) { f.plan = nil }, stayUnknownPlan},
		{"no plan known", func(f *limitFacts) { f.plan = &config.Identity{AccountUUID: "u1"} }, stayUnknownPlan},
		{"resets in 15 minutes", func(f *limitFacts) { f.resetsIn = 15 * time.Minute }, stayNear},
		{"resets in 16 minutes", func(f *limitFacts) { f.resetsIn = 16 * time.Minute }, moves},
		{"Remote Control", func(f *limitFacts) { f.remote = true }, stayRemote},
		{"Remote Control, resets in 9 minutes", func(f *limitFacts) { f.remote, f.resetsIn = true, 9*time.Minute }, stayNear},
		{"no relay", func(f *limitFacts) { f.relayed = false }, stayNotRelayed},
		{"Remote Control, no relay", func(f *limitFacts) { f.remote, f.relayed = true, false }, stayRemote},
		{"not the relay's claude", func(f *limitFacts) { f.child = false }, stayOff},
	} {
		f := base
		if c.edit != nil {
			c.edit(&f)
		}
		if got := eligible(f); got != c.want {
			t.Errorf("%s: %d, want %d", c.name, got, c.want)
		}
	}
}

func TestStayNotes(t *testing.T) {
	hit := "ccam: p1 hit its 5-hour limit (resets 3:45pm)"
	for _, c := range []struct {
		s      stay
		remote bool
		want   string
	}{
		{stayNear, false, "ccam: p1 hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it."},
		{stayNear, true, "ccam: p1 hit its 5-hour limit, which resets in 9m. Remote Control is on, so Claude Code does not wait on its own: send a prompt after the reset."},
		{stayRemote, true, hit + ". This session is connected to Remote Control and would not reconnect as another account, so ccam leaves it here."},
		{stayNotRelayed, false, hit + ". This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`."},
		{stayUnknownPlan, false, hit + ". ccam does not know p1's plan, so auto_switch leaves this session."},
		{stayOff, false, ""},
		{moves, false, ""},
	} {
		f := limitFacts{resetsIn: 9*time.Minute + 30*time.Second, remote: c.remote}
		if got := stayNote(c.s, f, "p1", "5-hour", hit); got != c.want {
			t.Errorf("%d: %q, want %q", c.s, got, c.want)
		}
	}
}

// With nobody to move to, the notification says who frees first, and that
// Claude Code waits within a day (beyond, it shows its menu instead).
func TestNoRoomNote(t *testing.T) {
	now := time.Now()
	at := func(d time.Duration) string { return handoff.FormatReset(now.Add(d), now) }
	soon, later := 3*time.Hour, 72*time.Hour
	hit := handoff.HitLimit("p1", "5-hour", at(soon))
	weekly := handoff.HitLimit("p1", "weekly", at(later))
	personal := &config.Identity{OrgUUID: "o1", OrgType: "claude_max"}
	full := target{full: true, first: "p2", freeAt: now.Add(time.Hour)}
	for _, c := range []struct {
		name  string
		tg    target
		id    *config.Identity
		hit   string
		reset time.Duration
		want  string
	}{
		{"personal plans", full, personal, hit, soon, hit + ". None of your other personal plans has room now; p2 frees up first, at " + at(time.Hour) +
			". Claude Code waits for p1 here; /ccam:switch p2 after " + at(time.Hour) + " moves it sooner."},
		{"seats", full, &config.Identity{OrgUUID: "acme", OrgName: "Acme", OrgType: "claude_team"}, hit, soon,
			hit + ". None of your other accounts in Acme has room now; p2 frees up first, at " + at(time.Hour) + ". Claude Code waits for p1 here; /ccam:switch p2 after " + at(time.Hour) + " moves it sooner."},
		{"seats of an unnamed organization", full, &config.Identity{OrgUUID: "big", OrgType: "claude_enterprise"}, hit, soon,
			hit + ". None of your other accounts in this organization has room now; p2 frees up first, at " + at(time.Hour) + ". Claude Code waits for p1 here; /ccam:switch p2 after " + at(time.Hour) + " moves it sooner."},
		{"a reset more than a day away", full, personal, weekly, later,
			weekly + ". None of your other personal plans has room now; p2 frees up first, at " + at(time.Hour) + ". /ccam:switch p2 after " + at(time.Hour) + " moves it there."},
		{"nobody frees first", target{full: true, first: "p2", freeAt: now.Add(4 * time.Hour)}, personal, hit, soon,
			hit + ". None of your other personal plans has room now. Claude Code waits for p1 here."},
		{"nobody frees first, a day away", target{full: true, first: "p2", freeAt: now.Add(80 * time.Hour)}, personal, weekly, later,
			weekly + ". None of your other personal plans has room now."},
		{"one untrusted", target{full: true, untrusted: []string{"p3"}}, personal, hit, soon,
			hit + ". p3 has not trusted this folder yet, so the conversation was not moved there. Open the folder once as p3 to allow it."},
		{"two untrusted", target{untrusted: []string{"p3", "p4"}}, personal, hit, soon,
			hit + ". p3 and p4 have not trusted this folder yet, so the conversation was not moved there. Open the folder once as each of them to allow it."},
		{"three untrusted", target{untrusted: []string{"p3", "p4", "p5"}}, personal, hit, soon,
			hit + ". p3, p4 and p5 have not trusted this folder yet, so the conversation was not moved there. Open the folder once as each of them to allow it."},
		// An account with room is left out with its usage credits on (D4):
		// the person may still move there.
		{"usage credits on", target{full: true, first: "p2", freeAt: now.Add(time.Hour), credits: []string{"p3"}}, personal, hit, soon,
			hit + ". p3 has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch p3 moves it."},
		{"usage credits on in two", target{credits: []string{"p3", "p4"}}, personal, hit, soon,
			hit + ". p3 and p4 have room but their usage credits are on, so auto_switch did not move the conversation there; /ccam:switch p3 moves it."},
		// One with room that may lack the conversation's agent, or define it
		// otherwise, is left out too.
		{"the agent's", target{full: true, first: "p2", freeAt: now.Add(time.Hour), agent: "reviewer", otherAgents: []string{"p3"}}, personal, hit, soon,
			hit + `. This conversation runs as agent "reviewer", and p3 does not share p1's agents, plugins and settings, so auto_switch did not move it there.`},
		{"the agent's in two", target{agent: "reviewer", otherAgents: []string{"p3", "p4"}}, personal, hit, soon,
			hit + `. This conversation runs as agent "reviewer", and p3 and p4 do not share p1's agents, plugins and settings, so auto_switch did not move it there.`},
		{"nothing known of the others", target{}, personal, hit, soon, ""},
	} {
		if got := c.tg.note(c.hit, "p1", c.id, now.Add(c.reset), now); got != c.want {
			t.Errorf("%s:\n got %q\nwant %q", c.name, got, c.want)
		}
	}
}

// logInTrusting writes the oauthAccount Claude Code keeps in .claude.json,
// and the folders the account trusted.
func logInTrusting(t *testing.T, dir, account, org, orgType string, folders ...string) {
	t.Helper()
	projects := map[string]any{}
	for _, f := range folders {
		projects[f] = map[string]any{"hasTrustDialogAccepted": true}
	}
	data, err := json.Marshal(map[string]any{"projects": projects, "oauthAccount": map[string]string{"accountUuid": account,
		"emailAddress": account + "@example.com", "organizationUuid": org, "organizationName": "Org " + org, "organizationType": orgType}})
	if err == nil {
		err = os.WriteFile(filepath.Join(dir, ".claude.json"), data, 0o600)
	}
	if err != nil {
		t.Fatal(err)
	}
}

// saveUsage saves account name's usage as a lookup at at left it.
func saveUsage(t *testing.T, c *testCLI, name string, fiveHour int, at time.Time, credits bool) {
	t.Helper()
	path := filepath.Join(c.p.State, "usage", name+".json")
	os.MkdirAll(filepath.Dir(path), 0o700)
	data := fmt.Sprintf(`{"at":%q,"limits":{"five_hour":{"percent":%d,"resets_at":"2030-01-01T00:00:00Z","known":true},"weekly":{"percent":10,"resets_at":"2030-01-05T00:00:00Z","known":true},"credits":%v}}`,
		at.UTC().Format(time.RFC3339Nano), fiveHour, credits)
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
}

// saveLogin writes account name's login with a refresh token (expired at
// refreshEnd when that is not zero).
func saveLogin(t *testing.T, c *testCLI, name string, fiveHour int, refresh string, refreshEnd time.Time) {
	t.Helper()
	cred := fmt.Sprintf(`{"claudeAiOauth":{"accessToken":"pct-%d","refreshToken":%q,"expiresAt":%d,"refreshTokenExpiresAt":%d}}`,
		fiveHour, refresh, time.Now().Add(time.Hour).UnixMilli(), refreshEnd.UnixMilli())
	if err := os.WriteFile(filepath.Join(c.p.ProfileDir(name), ".credentials.json"), []byte(cred), 0o600); err != nil {
		t.Fatal(err)
	}
}

// autoTarget takes another account of the leaving one's pool with room:
// never the same Claude account, an API account, one without a known plan
// or a login Claude Code can renew, one with usage credits on, or one that
// never trusted the folder; recently chosen accounts last.
func TestAutoTarget(t *testing.T) {
	folder, _ := filepath.EvalSymlinks(t.TempDir())
	now := time.Now()
	for _, k := range []struct {
		name      string
		usage     map[string]int
		edit      func(t *testing.T, c *testCLI)
		want      string
		untrusted []string
		first     string // when nobody has room, the account that frees first
		credits   []string
		why       string // when not "", why want was chosen
	}{
		{"personal plans: the least used", nil, nil, "side", nil, "", nil, ""},
		{"only Team seats have room", nil, func(t *testing.T, c *testCLI) {
			logInTrusting(t, c.p.ProfileDir("work"), "u2", "acme", "claude_team", folder)
			logInTrusting(t, c.p.ProfileDir("side"), "u3", "acme", "claude_team", folder)
		}, "", nil, "", nil, ""},
		{"cross_org does not apply", nil, func(t *testing.T, c *testCLI) {
			cfg := c.cfg()
			cfg.CrossOrg = true
			cfg.Save(c.p.ConfigFile())
			logInTrusting(t, c.p.ProfileDir("work"), "u2", "acme", "claude_team", folder)
			logInTrusting(t, c.p.ProfileDir("side"), "u3", "acme", "claude_team", folder)
		}, "", nil, "", nil, ""},
		{"seats of one organization", nil, func(t *testing.T, c *testCLI) {
			logInTrusting(t, c.p.ProfileDir("personal"), "u1", "acme", "claude_team", folder)
			logInTrusting(t, c.p.ProfileDir("work"), "u2", "acme", "claude_team", folder)
			logInTrusting(t, c.p.ProfileDir("side"), "u3", "other", "claude_team", folder)
		}, "work", nil, "", nil, ""},
		{"the same Claude account", nil, func(t *testing.T, c *testCLI) {
			logInTrusting(t, c.p.ProfileDir("side"), "u1", "p1", "claude_max", folder)
		}, "work", nil, "", nil, ""},
		{"an API account", nil, func(t *testing.T, c *testCLI) {
			cfg := c.cfg()
			cfg.Profiles[2].API = &config.API{KeyFrom: "env:SIDE_KEY"}
			cfg.Save(c.p.ConfigFile())
		}, "work", nil, "", nil, ""},
		{"nobody known", nil, func(t *testing.T, c *testCLI) { os.Remove(filepath.Join(c.p.ProfileDir("side"), ".claude.json")) }, "work", nil, "", nil, ""},
		{"a refresh token that expired", nil, func(t *testing.T, c *testCLI) {
			saveLogin(t, c, "side", 10, "r", now.Add(-time.Hour))
			saveUsage(t, c, "side", 10, now, false) // fresh: Fetch reads no login
		}, "work", nil, "", nil, ""},
		{"a refresh token Claude Code emptied", nil, func(t *testing.T, c *testCLI) {
			saveLogin(t, c, "side", 10, "", time.Time{})
			saveUsage(t, c, "side", 10, now.Add(-2*time.Minute), false)
		}, "work", nil, "", nil, ""},
		{"usage credits on", nil, func(t *testing.T, c *testCLI) { saveUsage(t, c, "side", 10, now, true) }, "work", nil, "", []string{"side"}, ""},
		{"only an account with usage credits on has room", map[string]int{"personal": 90, "work": 100, "side": 10}, func(t *testing.T, c *testCLI) {
			saveUsage(t, c, "side", 10, now, true)
		}, "", nil, "work", []string{"side"}, ""},
		{"usage credits on and full", nil, func(t *testing.T, c *testCLI) { saveUsage(t, c, "side", 100, now, true) }, "work", nil, "", nil, ""},
		{"an account that never trusted the folder", nil, func(t *testing.T, c *testCLI) {
			logInTrusting(t, c.p.ProfileDir("side"), "u3", "s1", "claude_pro")
		}, "work", []string{"side"}, "", nil, ""},
		{"only an account that never trusted the folder has room", map[string]int{"personal": 90, "work": 100, "side": 10}, func(t *testing.T, c *testCLI) {
			logInTrusting(t, c.p.ProfileDir("side"), "u3", "s1", "claude_pro")
		}, "", []string{"side"}, "work", nil, ""},
		{"trust of a folder above, outside git", nil, func(t *testing.T, c *testCLI) {
			logInTrusting(t, c.p.ProfileDir("side"), "u3", "s1", "claude_pro", filepath.Dir(folder))
		}, "side", nil, "", nil, ""},
		// work is not the least used then.
		{"chosen a moment ago", nil, func(t *testing.T, c *testCLI) {
			markPicked(t, pickPath(&env{paths: c.p}, "side"), now)
		}, "work", nil, "", nil, "40% of its 5 hours used; side was taken by another session a moment ago"},
		{"chosen 3 minutes ago", nil, func(t *testing.T, c *testCLI) {
			markPicked(t, pickPath(&env{paths: c.p}, "side"), now.Add(-3*time.Minute))
		}, "side", nil, "", nil, ""},
		{"nobody has room", map[string]int{"personal": 90, "work": 100, "side": 10}, func(t *testing.T, c *testCLI) {
			usage.SaveLimit(filepath.Join(c.p.State, "usage"), "side", "five_hour", now.Add(time.Hour), now)
		}, "", nil, "side", nil, ""},
	} {
		t.Run(k.name, func(t *testing.T) {
			if k.usage == nil {
				k.usage = map[string]int{"personal": 90, "work": 40, "side": 10}
			}
			c := newTestCLI(t, k.usage)
			logInTrusting(t, c.p.ProfileDir("personal"), "u1", "p1", "claude_max", folder)
			logInTrusting(t, c.p.ProfileDir("work"), "u2", "w1", "claude_max", folder)
			logInTrusting(t, c.p.ProfileDir("side"), "u3", "s1", "claude_pro", folder)
			if k.edit != nil {
				k.edit(t, c)
			}
			e, cfg := c.env()
			from, _ := cfg.Find("personal")
			tg, err := autoTarget(context.Background(), e, cfg, from, folder, "", usage.Policy{})
			got := ""
			if tg.to != nil {
				got = tg.to.Name
			}
			if err != nil || got != k.want || !slices.Equal(tg.untrusted, k.untrusted) || tg.first != k.first || tg.full != (k.first != "") ||
				!slices.Equal(tg.credits, k.credits) {
				t.Fatalf("%v: %q %+v", err, got, tg)
			}
			if got == "side" && tg.why != "10% of its 5 hours used, the least" || k.why != "" && tg.why != k.why {
				t.Fatalf("why: %q", tg.why)
			}
		})
	}
}

// While the conversation runs as an agent, autoTarget leaves out an account
// with room that does not load the agents of the one at its limit the same
// way, and names it for the notification.
func TestAutoTargetKeepsTheAgent(t *testing.T) {
	folder, _ := filepath.EvalSymlinks(t.TempDir())
	for _, k := range []struct {
		name, agent string
		usage       map[string]int
		want        string
		other       []string
	}{
		{"no agent", "", map[string]int{"personal": 90, "work": 40, "side": 10}, "side", nil},
		{"side has agents of its own", "reviewer", map[string]int{"personal": 90, "work": 40, "side": 10}, "work", []string{"side"}},
		{"only side has room", "reviewer", map[string]int{"personal": 90, "work": 100, "side": 10}, "", []string{"side"}},
	} {
		t.Run(k.name, func(t *testing.T) {
			c := newTestCLI(t, k.usage)
			logInTrusting(t, c.p.ProfileDir("personal"), "u1", "p1", "claude_max", folder)
			logInTrusting(t, c.p.ProfileDir("work"), "u2", "w1", "claude_max", folder)
			logInTrusting(t, c.p.ProfileDir("side"), "u3", "s1", "claude_pro", folder)
			if err := os.MkdirAll(filepath.Join(c.p.ProfileDir("side"), "agents"), 0o700); err != nil {
				t.Fatal(err)
			}
			e, cfg := c.env()
			from, _ := cfg.Find("personal")
			tg, err := autoTarget(context.Background(), e, cfg, from, folder, k.agent, usage.Policy{})
			got := ""
			if tg.to != nil {
				got = tg.to.Name
			}
			if err != nil || got != k.want || !slices.Equal(tg.otherAgents, k.other) || tg.agent != k.agent {
				t.Fatalf("%v: %q %+v", err, got, tg)
			}
		})
	}
}

// With usage_source = "claude" the hook goes by the saved usage, and says
// how old it is: asking Claude Code takes up to 20 s.
func TestAutoTargetStartsNoClaude(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the stand-in claude is a shell script")
	}
	c := newTestCLI(t, map[string]int{"personal": 90, "work": 40, "side": 10})
	folder, _ := filepath.EvalSymlinks(t.TempDir())
	logInTrusting(t, c.p.ProfileDir("personal"), "u1", "p1", "claude_max", folder)
	logInTrusting(t, c.p.ProfileDir("work"), "u2", "w1", "claude_max", folder)
	logInTrusting(t, c.p.ProfileDir("side"), "u3", "s1", "claude_pro", folder)
	ran := filepath.Join(t.TempDir(), "ran")
	bin := filepath.Join(t.TempDir(), "claude")
	os.WriteFile(bin, []byte("#!/bin/sh\ntouch '"+ran+"'\n"), 0o755)
	cfg := c.cfg()
	cfg.UsageSource, cfg.ClaudeBin = "claude", bin
	cfg.Save(c.p.ConfigFile())
	saved := time.Now().Add(-2 * time.Minute)
	saveUsage(t, c, "work", 40, saved, false)
	saveUsage(t, c, "side", 10, saved, false)
	e, cfg := c.env()
	from, _ := cfg.Find("personal")
	tg, err := autoTarget(context.Background(), e, cfg, from, folder, "", usage.Policy{})
	if err != nil || tg.to == nil || tg.to.Name != "side" || tg.why != "10% of its 5 hours used, the least; saved 2m ago" || fileThere(ran) {
		t.Fatalf("%v %+v, claude started: %v", err, tg, fileThere(ran))
	}
}

// markPicked writes the mark of an account picked at at.
func markPicked(t *testing.T, path string, at time.Time) {
	t.Helper()
	os.MkdirAll(filepath.Dir(path), 0o700)
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	os.Chtimes(path, at, at)
}

// slowUsage answers usage lookups after a second.
func slowUsage(t *testing.T) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)
		var pct int
		fmt.Sscanf(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer pct-"), "%d", &pct)
		fmt.Fprintf(w, `{"limits":[{"kind":"session","percent":%d,"resets_at":"2030-01-01T00:00:00Z"}]}`, pct)
	}))
	old := newFetcher
	newFetcher = func(dir string) *usage.Fetcher {
		f := usage.NewFetcher(dir)
		f.URL = srv.URL
		return f
	}
	t.Cleanup(func() { newFetcher = old })
	return srv
}

// A lookup that hangs does not keep the hook past its deadline: the choice
// is given up, and with time enough it is made.
func TestAutoTargetDeadline(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 90, "work": 40, "side": 10})
	folder, _ := filepath.EvalSymlinks(t.TempDir())
	logInTrusting(t, c.p.ProfileDir("personal"), "u1", "p1", "claude_max", folder)
	logInTrusting(t, c.p.ProfileDir("work"), "u2", "w1", "claude_max", folder)
	logInTrusting(t, c.p.ProfileDir("side"), "u3", "s1", "claude_pro", folder)
	saveUsage(t, c, "side", 10, time.Now().Add(-2*time.Minute), false)
	srv := slowUsage(t)
	defer srv.Close()
	e, cfg := c.env()
	from, _ := cfg.Find("personal")
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	begin := time.Now()
	if tg, err := autoTarget(ctx, e, cfg, from, folder, "", usage.Policy{}); err == nil || time.Since(begin) > time.Second {
		t.Fatalf("%v %+v after %v", err, tg, time.Since(begin))
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if tg, err := autoTarget(ctx, e, cfg, from, folder, "", usage.Policy{}); err != nil || tg.to == nil || tg.to.Name != "side" {
		t.Fatalf("%v %+v", err, tg)
	}
}

// claudeChild starts a process that stands in for the relay's Claude Code:
// this test process is the relay.
func claudeChild(t *testing.T) (int, string) {
	t.Helper()
	claude := exec.Command("sleep", "30")
	if runtime.GOOS == "windows" {
		claude = exec.Command("ping", "-n", "31", "127.0.0.1")
	}
	if err := claude.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { claude.Process.Kill(); claude.Wait() })
	start, _ := session.StartTime(claude.Process.Pid)
	return claude.Process.Pid, start
}

// terminated records what the hook would end instead of ending it.
type terminated struct {
	claude, relay, child int
	start                string
}

// The stop-failure hook under a relay of this ccam: with auto_switch on it
// notes the account it chose and ends Claude Code (the relay then carries
// the conversation on), or says why it leaves the session.
func TestAutoSwitchHook(t *testing.T) {
	hookEnv(t)
	t.Setenv("CLAUDE_EFFORT", "high")
	c := newTestCLI(t, map[string]int{"personal": 90, "work": 40, "side": 10})
	var calls []terminated
	failEnd := false
	var whileEnding func() // what happens meanwhile, as the hook ends Claude Code
	oldEnd, oldDeadline, oldClaude := launcher.Terminate, autoDeadline, endClaude
	launcher.Terminate = func(claude, relay, child int, start string) error {
		calls = append(calls, terminated{claude, relay, child, start})
		if whileEnding != nil {
			whileEnding()
		}
		if failEnd {
			return errors.New("no such process")
		}
		return nil
	}
	t.Cleanup(func() { launcher.Terminate, autoDeadline, endClaude = oldEnd, oldDeadline, oldClaude })
	pid, start := claudeChild(t)
	me := "2:" + strconv.Itoa(os.Getpid())
	t.Setenv("CCAM_PROFILE", "personal")
	cwd, _ := os.Getwd()
	cwd, _ = filepath.EvalSymlinks(cwd) // as Claude Code records the folders trusted
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	msg := "You've hit your session limit · resets 3:30am (Asia/Seoul)" // the fixture's
	reset := time.Now().Add(3 * time.Hour).Truncate(time.Second)
	// fresh sets the scene: three personal plans that trust this folder,
	// auto_switch on, a conversation that stopped at personal's 5-hour
	// limit, and nothing ccam wrote before.
	fresh := func(reset time.Time, edits ...string) {
		t.Helper()
		calls, failEnd, whileEnding = nil, false, nil
		for _, d := range []string{"handoff", "picks", "usage"} {
			os.RemoveAll(filepath.Join(c.p.State, d))
		}
		cfg := c.cfg()
		cfg.AutoSwitch, cfg.CrossOrg, cfg.Profiles[0].Login = true, false, nil
		if err := cfg.Save(c.p.ConfigFile()); err != nil {
			t.Fatal(err)
		}
		logInTrusting(t, c.p.ProfileDir("personal"), "u1", "p1", "claude_max", cwd)
		logInTrusting(t, c.p.ProfileDir("work"), "u2", "w1", "claude_max", cwd)
		logInTrusting(t, c.p.ProfileDir("side"), "u3", "s1", "claude_pro", cwd)
		limitConversation(t, conv, time.Now().Add(-time.Second), reset, edits...)
		t.Setenv("CCAM_RELAY", me)
		t.Setenv("CLAUDE_PID", strconv.Itoa(pid))
		t.Setenv("CLAUDE_CODE_BRIDGE_SESSION_ID", "")
	}
	// hook runs stop-failure and returns its notification's text.
	hook := func(message string) string {
		t.Helper()
		in := hookJSON(t, map[string]string{"error": "rate_limit", "last_assistant_message": message, "session_id": "s1", "transcript_path": conv, "cwd": cwd})
		out, errOut, code := c.run(in, "hook", "stop-failure")
		if code != 0 {
			t.Fatalf("%d %s", code, errOut)
		}
		if out == "" {
			return ""
		}
		var got map[string]string
		if err := json.Unmarshal([]byte(out), &got); err != nil {
			t.Fatalf("%q", out)
		}
		return strings.TrimSuffix(strings.TrimPrefix(got["terminalSequence"], "\x1b]9;"), "\x07")
	}
	note := func() *handoff.Ticket {
		t.Helper()
		n, ok := handoff.LoadFile(c.p.State, handoff.Path(c.p.State, "personal", "s1"))
		if !ok {
			t.Fatal("no note")
		}
		return n
	}
	hit := "ccam: personal hit its 5-hour limit (resets " + handoff.FormatReset(reset, time.Now()) + ")"
	usual := hit + ". `ccam next --continue` continues the conversation on another account."
	// stays checks that the hook left the session with the notification want.
	stays := func(what, got, want string) {
		t.Helper()
		if n := note(); got != want || len(calls) != 0 || n.Chosen != "" || n.Stopped || n.Relay != "" {
			t.Fatalf("%s: %q, want %q; ended %v; note %+v", what, got, want, calls, n)
		}
	}
	// switchTo is the person's /ccam:switch to; ends is what ending Claude
	// Code from its hook gives.
	switchTo := func(to string, ends error) {
		t.Helper()
		endClaude = func(int) error { return ends }
		e, _ := c.env()
		e.stdout = io.Discard
		if err := hookSwitch(e, hookInput{SessionID: "s1", Transcript: conv, Cwd: cwd, CommandArgs: to}, "personal"); err != nil {
			t.Fatal(err)
		}
	}

	if runtime.GOOS == "windows" {
		// Windows ignores auto_switch.
		fresh(reset)
		stays("Windows", hook(msg), usual)
		return
	}

	fresh(reset)
	if got := hook(msg); got != "" || !slices.Equal(calls, []terminated{{pid, os.Getpid(), pid, start}}) {
		t.Fatalf("%q %v", got, calls)
	}
	n := note()
	if n.Reason != handoff.Limit || n.Chosen != "side" || n.Relay != me || n.Child != pid || !n.Stopped || n.Why != "10% of its 5 hours used, the least" ||
		n.Window != "five_hour" || !n.ResetsAt.Equal(reset) || n.LimitAt.IsZero() || n.Effort != "high" {
		t.Fatalf("%+v", n)
	}
	if !fileThere(filepath.Join(c.p.State, "picks", "side")) || !fileThere(filepath.Join(c.p.State, "usage", "personal.limit.json")) {
		t.Fatal("no pick mark or no limit record")
	}

	// Seats of one organization stand in for each other (D7 b), with no
	// login recorded in config.toml.
	fresh(reset)
	logInTrusting(t, c.p.ProfileDir("personal"), "u1", "acme", "claude_team", cwd)
	logInTrusting(t, c.p.ProfileDir("work"), "u2", "acme", "claude_team", cwd)
	logInTrusting(t, c.p.ProfileDir("side"), "u3", "acme", "claude_team", cwd)
	if got := hook(msg); got != "" || len(calls) != 1 || note().Chosen != "side" {
		t.Fatalf("seats: %q %v %+v", got, calls, note())
	}
	// Never into another organization, nor between a seat and a personal
	// plan, whatever cross_org says. With no other account of its pool,
	// `ccam next --continue` would refuse too, and the notification says so.
	fresh(reset)
	logInTrusting(t, c.p.ProfileDir("personal"), "u1", "acme", "claude_team", cwd)
	logInTrusting(t, c.p.ProfileDir("work"), "u2", "other", "claude_team", cwd)
	stays("another organization", hook(msg), hit+". You have no other accounts in Org acme, so ccam cannot move this conversation.")
	fresh(reset)
	cfg := c.cfg()
	cfg.CrossOrg = true
	cfg.Save(c.p.ConfigFile())
	logInTrusting(t, c.p.ProfileDir("work"), "u2", "acme", "claude_team", cwd)
	logInTrusting(t, c.p.ProfileDir("side"), "u3", "acme", "claude_team", cwd)
	stays("cross_org", hook(msg), usual)

	// The plan comes from the login ccam recorded when .claude.json has none.
	fresh(reset)
	os.Remove(filepath.Join(c.p.ProfileDir("personal"), ".claude.json"))
	cfg = c.cfg()
	cfg.Profiles[0].Login = &config.Identity{AccountUUID: "u1", OrgUUID: "p1", OrgType: "claude_max"}
	cfg.Save(c.p.ConfigFile())
	if got := hook(msg); got != "" || note().Chosen != "side" {
		t.Fatalf("recorded login: %q %+v", got, note())
	}
	fresh(reset)
	os.Remove(filepath.Join(c.p.ProfileDir("personal"), ".claude.json"))
	stays("no plan known", hook(msg), hit+". ccam does not know personal's plan, so auto_switch leaves this session.")

	// Off, or a limit that is no plan's (a spend limit): as without auto_switch.
	fresh(reset)
	cfg = c.cfg()
	cfg.AutoSwitch = false
	cfg.Save(c.p.ConfigFile())
	stays("off", hook(msg), usual)
	fresh(reset)
	spend := "You've hit your org's monthly spend limit · ask your admin to raise it · your session limit resets 3:30am"
	limitConversation(t, conv, time.Now().Add(-time.Second), reset, msg, spend)
	stays("a spend limit", hook(spend), usual)
	if !fileThere(filepath.Join(c.p.State, "usage", "personal.limit.json")) {
		t.Fatal("a spend limit not recorded")
	}

	// Without the limit recorded, a hop could come back to personal before
	// the reset: nothing moves.
	fresh(reset)
	os.MkdirAll(c.p.State, 0o700)
	os.WriteFile(filepath.Join(c.p.State, "usage"), nil, 0o600)
	stays("not recorded", hook(msg), usual)

	// No relay of this ccam, or a claude that is not the relay's.
	notRelayed := hit + ". This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`."
	for _, relay := range []string{"", strconv.Itoa(os.Getpid())} {
		fresh(reset)
		t.Setenv("CCAM_RELAY", relay)
		stays("CCAM_RELAY="+relay, hook(msg), notRelayed)
	}
	fresh(reset)
	t.Setenv("CCAM_RELAY", "2:"+strconv.Itoa(os.Getppid()))
	stays("another relay", hook(msg), usual)

	// Remote Control (D11): the conversation says so, or the environment.
	remote := hit + ". This session is connected to Remote Control and would not reconnect as another account, so ccam leaves it here."
	fresh(reset)
	appendRows(t, conv, `{"type":"bridge-session","sessionId":"s1","bridgeSessionId":"cse_test","lastSequenceNum":1}`)
	stays("a bridge-session row", hook(msg), remote)
	fresh(reset)
	t.Setenv("CLAUDE_CODE_BRIDGE_SESSION_ID", "cse_test")
	stays("CLAUDE_CODE_BRIDGE_SESSION_ID", hook(msg), remote)

	// A reset in a few minutes: Claude Code's own wait is quicker.
	near := time.Now().Add(10*time.Minute + 30*time.Second)
	fresh(near)
	stays("near the reset", hook(msg), "ccam: personal hit its 5-hour limit, which resets in 10m; ccam leaves this session to wait for it.")

	// Claude Code could not be ended: the note says nothing was chosen.
	fresh(reset)
	failEnd = true
	if got, n := hook(msg), note(); got != usual || len(calls) != 1 || n.Chosen != "" || n.Stopped || n.Relay != "" {
		t.Fatalf("Terminate failed: %q %v %+v", got, calls, n)
	}
	// ... unless the person typed /ccam:switch meanwhile: its note stays. A
	// note that went meanwhile comes back.
	fresh(reset)
	failEnd, whileEnding = true, func() { switchTo("work", nil) }
	if got, n := hook(msg), note(); got != usual || len(calls) != 1 || n.Reason != handoff.Switch || n.Target != "work" || !n.Stopped {
		t.Fatalf("a switch as Terminate failed: %q %v %+v", got, calls, n)
	}
	fresh(reset)
	failEnd, whileEnding = true, func() { os.Remove(handoff.Path(c.p.State, "personal", "s1")) }
	if got, n := hook(msg), note(); got != usual || len(calls) != 1 || n.Reason != handoff.Limit || n.Chosen != "" || n.Stopped || n.Relay != "" {
		t.Fatalf("the note gone as Terminate failed: %q %v %+v", got, calls, n)
	}

	// The person typed /ccam:switch work while the hook chose (Claude Code
	// writes rows of its own for it, no prompt): the switch stands, whether
	// it ended Claude Code or not, and the hook ends nothing.
	for _, ends := range []error{nil, errors.New("no sh")} {
		fresh(reset)
		looks := newFetcher
		newFetcher = func(dir string) *usage.Fetcher {
			newFetcher = looks
			switchTo("work", ends)
			return looks(dir)
		}
		if got, n := hook(msg), note(); got != usual || len(calls) != 0 || n.Reason != handoff.Switch || n.Target != "work" || n.Stopped != (ends == nil) {
			t.Fatalf("a switch while choosing (%v): %q %v %+v", ends, got, calls, n)
		}
	}

	// Another stop-failure while Claude Code ends for a move (the person
	// sent a ! command after the limit, which was refused too): the move's
	// note stays, also once Claude Code is gone.
	gone := exec.Command("sleep", "0")
	if err := gone.Run(); err != nil {
		t.Fatal(err)
	}
	for _, claude := range []int{pid, gone.Process.Pid} {
		fresh(reset)
		t.Setenv("CLAUDE_PID", strconv.Itoa(claude))
		err := handoff.Save(c.p.State, handoff.Ticket{Account: "personal", SessionID: "s1", Transcript: conv, Cwd: cwd, At: time.Now(),
			Reason: handoff.Limit, Relay: me, Child: claude, Stopped: true, Chosen: "work"})
		if err != nil {
			t.Fatal(err)
		}
		if got, n := hook(msg), note(); got != usual || len(calls) != 0 || n.Chosen != "work" || !n.Stopped {
			t.Fatalf("a move under way (claude %d): %q %v %+v", claude, got, calls, n)
		}
	}
	// The move of another relay (a hop that was called off), or of an
	// earlier Claude Code of this one, is over: the session goes on.
	for _, old := range []handoff.Ticket{{Relay: "2:1", Child: pid}, {Relay: me, Child: gone.Process.Pid}} {
		fresh(reset)
		old.Account, old.SessionID, old.Transcript, old.Cwd, old.At = "personal", "s1", conv, cwd, time.Now().Add(-time.Hour)
		old.Reason, old.Stopped, old.Chosen = handoff.Limit, true, "work"
		if err := handoff.Save(c.p.State, old); err != nil {
			t.Fatal(err)
		}
		if got, n := hook(msg), note(); got != "" || len(calls) != 1 || n.Chosen != "side" {
			t.Fatalf("an earlier move (%s, %d): %q %v %+v", old.Relay, old.Child, got, calls, n)
		}
	}

	// Nobody has room.
	fresh(reset)
	saveUsage(t, c, "work", 100, time.Now(), false)
	saveUsage(t, c, "side", 100, time.Now(), false)
	stays("no room", hook(msg), hit+". None of your other personal plans has room now. Claude Code waits for personal here.")

	// The person sent a prompt while the hook chose: the session stays.
	fresh(reset)
	lim, _ := transcript.LastLimit(conv, time.Time{})
	appendRows(t, conv, `{"type":"user","message":{"role":"user","content":"(go on)"},"origin":{"kind":"human"}}`)
	e, _ := c.env()
	in := hookInput{SessionID: "s1", Transcript: conv, Cwd: cwd, Error: "rate_limit", LastMessage: msg}
	tk := handoff.Ticket{Account: "personal", SessionID: "s1", Transcript: conv, Cwd: cwd, At: time.Now(), Reason: handoff.Limit, Limit: "session",
		Window: lim.Window, ResetsAt: lim.ResetsAt, LimitAt: lim.At}
	if moved, say := autoSwitch(e, in, tk, lim, time.Now(), hit); moved || say != "" || len(calls) != 0 || fileThere(handoff.Path(c.p.State, "personal", "s1")) {
		t.Fatalf("prompted after the limit: %v %q %v", moved, say, calls)
	}

	// The note naming the chosen account cannot be saved: nothing is ended,
	// or the relay would find nothing to carry the conversation on with.
	fresh(reset)
	os.MkdirAll(filepath.Join(c.p.State, "handoff"), 0o700)
	if err := os.WriteFile(filepath.Join(c.p.State, "handoff", "personal"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	lim, _ = transcript.LastLimit(conv, time.Time{})
	if moved, say := autoSwitch(e, in, tk, lim, time.Now(), hit); moved || say != "" || len(calls) != 0 {
		t.Fatalf("no note saved: %v %q %v", moved, say, calls)
	}
	// A /ccam:switch from before the limit that did not end Claude Code is
	// no reason to stay.
	fresh(reset)
	lim, _ = transcript.LastLimit(conv, time.Time{})
	err := handoff.Save(c.p.State, handoff.Ticket{Account: "personal", SessionID: "s1", Transcript: conv, Cwd: cwd, At: lim.At.Add(-time.Minute),
		Reason: handoff.Switch, Target: "work", Relay: me, Child: pid})
	if err != nil {
		t.Fatal(err)
	}
	if moved, say := autoSwitch(e, in, tk, lim, time.Now(), hit); !moved || say != "" || len(calls) != 1 || note().Chosen != "side" {
		t.Fatalf("an earlier switch: %v %q %v %+v", moved, say, calls, note())
	}

	// A lookup that hangs past the deadline: no choice.
	fresh(reset)
	autoDeadline = 300 * time.Millisecond
	srv := slowUsage(t)
	defer srv.Close()
	begin := time.Now()
	stays("too slow", hook(msg), usual)
	if took := time.Since(begin); took > 2*time.Second {
		t.Fatalf("took %v", took)
	}
}

// A Claude Code that has ended (its pid not the relay's child any more) is
// not ended, and the session is not moved.
func TestHookSkipsWhenClaudeGone(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("auto_switch is ignored on Windows")
	}
	hookEnv(t)
	c := newTestCLI(t, map[string]int{"personal": 90, "work": 40, "side": 10})
	cfg := c.cfg()
	cfg.AutoSwitch = true
	cfg.Save(c.p.ConfigFile())
	var calls []terminated
	old := launcher.Terminate
	launcher.Terminate = func(claude, relay, child int, start string) error {
		calls = append(calls, terminated{claude, relay, child, start})
		return nil
	}
	t.Cleanup(func() { launcher.Terminate = old })
	cwd, _ := os.Getwd()
	cwd, _ = filepath.EvalSymlinks(cwd) // as Claude Code records the folders trusted
	for _, name := range []string{"personal", "work", "side"} {
		logInTrusting(t, c.p.ProfileDir(name), "u-"+name, "o-"+name, "claude_max", cwd)
	}
	gone := exec.Command("sleep", "0")
	if err := gone.Run(); err != nil {
		t.Fatal(err)
	}
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	reset := time.Now().Add(3 * time.Hour)
	limitConversation(t, conv, time.Now().Add(-time.Second), reset)
	t.Setenv("CCAM_PROFILE", "personal")
	t.Setenv("CCAM_RELAY", "2:"+strconv.Itoa(os.Getpid()))
	t.Setenv("CLAUDE_PID", strconv.Itoa(gone.Process.Pid))
	msg := "You've hit your session limit · resets 3:30am (Asia/Seoul)"
	out, _, _ := c.run(hookJSON(t, map[string]string{"error": "rate_limit", "last_assistant_message": msg, "session_id": "s1", "transcript_path": conv, "cwd": cwd}), "hook", "stop-failure")
	n, ok := handoff.LoadFile(c.p.State, handoff.Path(c.p.State, "personal", "s1"))
	if !strings.Contains(out, "`ccam next --continue` continues the conversation on another account.") || len(calls) != 0 || !ok || n.Chosen != "" || n.Stopped {
		t.Fatalf("%q %v %+v", out, calls, n)
	}
}

// A subagent's limit holds up the main conversation when it is a window of
// the whole account, or the week of the model that conversation runs on.
func TestSubagentHolds(t *testing.T) {
	for _, c := range []struct {
		window, model string
		holds         bool
	}{
		{"five_hour", "claude-sonnet-5-5", true},
		{"seven_day", "", true},
		{"seven_day_opus", "claude-opus-5-5", true},
		{"seven_day_opus", "claude-sonnet-5-5", false},
		{"seven_day_sonnet", "claude-sonnet-5-5", true},
		{"seven_day_sonnet", "claude-opus-5-5", false},
		{"seven_day_opus", "", false}, // no reply to tell by
	} {
		if got := subagentHolds(c.window, c.model); got != c.holds {
			t.Errorf("%s, %q: %v", c.window, c.model, got)
		}
	}
}

// A subagent that stops at a usage limit under a relay with auto_switch on:
// the hook moves the session as for the main thread's limit, and the note
// says a subagent stopped, for the relay's first prompt.
func TestAutoSwitchHookSubagent(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("auto_switch is ignored on Windows")
	}
	hookEnv(t)
	c := newTestCLI(t, map[string]int{"personal": 90, "work": 40, "side": 10})
	var calls []terminated
	old := launcher.Terminate
	launcher.Terminate = func(claude, relay, child int, start string) error {
		calls = append(calls, terminated{claude, relay, child, start})
		return nil
	}
	t.Cleanup(func() { launcher.Terminate = old })
	pid, start := claudeChild(t)
	t.Setenv("CCAM_PROFILE", "personal")
	t.Setenv("CCAM_RELAY", "2:"+strconv.Itoa(os.Getpid()))
	t.Setenv("CLAUDE_PID", strconv.Itoa(pid))
	cwd, _ := os.Getwd()
	cwd, _ = filepath.EvalSymlinks(cwd) // as Claude Code records the folders trusted
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	reset := time.Now().Add(3 * time.Hour).Truncate(time.Second)
	fixture := "You've hit your session limit · resets 3:30am (Asia/Seoul)"
	// fresh sets the scene as TestAutoSwitchHook's does, with the limit in
	// subagent a1's conversation, and returns the message Claude Code gives.
	fresh := func(msg string, edits ...string) string {
		t.Helper()
		calls = nil
		for _, d := range []string{"handoff", "picks", "usage"} {
			os.RemoveAll(filepath.Join(c.p.State, d))
		}
		cfg := c.cfg()
		cfg.AutoSwitch, cfg.Profiles[0].Login = true, nil
		if err := cfg.Save(c.p.ConfigFile()); err != nil {
			t.Fatal(err)
		}
		logInTrusting(t, c.p.ProfileDir("personal"), "u1", "p1", "claude_max", cwd)
		logInTrusting(t, c.p.ProfileDir("work"), "u2", "w1", "claude_max", cwd)
		logInTrusting(t, c.p.ProfileDir("side"), "u3", "s1", "claude_pro", cwd)
		subagentConversation(t, conv, "a1", time.Now().Add(-time.Second), reset, append([]string{fixture, msg}, edits...)...)
		return msg
	}
	hook := func(msg string) string {
		t.Helper()
		out, errOut, code := c.run(`{"agent_id":"a1","agent_type":"general-purpose",`+hookJSON(t, map[string]string{"error": "rate_limit",
			"last_assistant_message": msg, "session_id": "s1", "transcript_path": conv, "cwd": cwd})[1:], "hook", "stop-failure")
		if code != 0 {
			t.Fatalf("%d %s", code, errOut)
		}
		return out
	}
	note := func() *handoff.Ticket {
		t.Helper()
		n, ok := handoff.LoadFile(c.p.State, handoff.Path(c.p.State, "personal", "s1"))
		if !ok {
			t.Fatal("no note")
		}
		return n
	}

	if out := hook(fresh(fixture)); out != "" || !slices.Equal(calls, []terminated{{pid, os.Getpid(), pid, start}}) {
		t.Fatalf("%q %v", out, calls)
	}
	if n := note(); n.Chosen != "side" || !n.Stopped || n.Subagent != "a1" || n.Transcript != conv || n.Window != "five_hour" {
		t.Fatalf("%+v", n)
	}

	// The Opus week holds up a subagent on Opus, not the conversation on
	// Sonnet: it is left to go on.
	opus := fresh("You've hit your Opus limit · resets 3:30am (Asia/Seoul)", `"rateLimitType":"five_hour"`, `"rateLimitType":"seven_day_opus"`)
	data, _ := os.ReadFile(conv)
	os.WriteFile(conv, []byte(strings.ReplaceAll(string(data), "claude-opus-5-5", "claude-sonnet-5-5")), 0o600)
	if out := hook(opus); !strings.Contains(out, "personal hit its Opus limit") || len(calls) != 0 || note().Stopped {
		t.Fatalf("Opus week, Sonnet conversation: %q %v", out, calls)
	}

	// The person sent a prompt after the subagent stopped: the session stays.
	msg := fresh(fixture)
	appendRows(t, conv, `{"type":"user","timestamp":`+strconv.Quote(time.Now().UTC().Format(time.RFC3339Nano))+`,"message":{"role":"user","content":"(go on)"},"origin":{"kind":"human"}}`)
	if out := hook(msg); !strings.Contains(out, "`ccam next --continue`") || len(calls) != 0 || note().Stopped {
		t.Fatalf("prompted after the limit: %q %v", out, calls)
	}
}
