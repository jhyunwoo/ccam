package usage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/fsutil"
)

// A limit record keeps what Claude Code reported when it refused a request
// at a usage limit: that window is full until its reset. Fetch lays it over
// whatever else is known, so an account that just hit its limit is not
// chosen again before then, however old or behind the saved numbers are
// (state/usage/<name>.limit.json, written by the ccam plugin's hook). Only
// a lookup of the usage endpoint made liftDelay or more after the limit
// lifts a window, when it shows that window below 100%: the limit was
// lifted before its reset, as Claude Code's /limit-reset does.
type limitRecord struct {
	At      time.Time            `json:"at"`      // when the last refused request was
	Windows map[string]time.Time `json:"windows"` // window → when it resets
}

// liftDelay is how long after a limit a lookup must be made to lift it:
// right after, the numbers may not show the limit yet, and an automatic
// choice then (auto_switch, next --auto, the dashboard's n) must not take
// the account Claude Code just refused. A limit reset is used later.
const liftDelay = 5 * time.Minute

// limitWindows are the windows a record keeps, by the names Claude Code's
// rate_limit_info gives them, in the order they are laid over a report.
// seven_day_overage_included is the Fable allowance, the Fable week of the
// usage endpoint. Usage credits (overage) are no window ccam chooses by.
var limitWindows = []string{"five_hour", "seven_day", "seven_day_opus", "seven_day_sonnet", "seven_day_overage_included"}

// LimitWindow reports whether SaveLimit records window.
func LimitWindow(window string) bool { return slices.Contains(limitWindows, window) }

// WindowModel is the model family whose week window is (opus, sonnet or
// fable, as Policy.Model and the usage endpoint's model windows name it),
// or "" for a window of the whole account.
func WindowModel(window string) string {
	switch window {
	case "seven_day_opus":
		return "opus"
	case "seven_day_sonnet":
		return "sonnet"
	case "seven_day_overage_included": // Claude Code calls it the "Fable limit"
		return "fable"
	}
	return ""
}

func (f *Fetcher) limitPath(name string) string { return filepath.Join(f.StateDir, name+".limit.json") }

// SaveLimit records that account name's window was full at at, when Claude
// Code refused a request, until resetsAt. It keeps the account's other
// windows that have not reset by then, and records nothing for a limit that
// is already over.
func SaveLimit(stateDir, name, window string, resetsAt, at time.Time) error {
	if !LimitWindow(window) {
		return fmt.Errorf("ccam does not record the %s window", window)
	}
	if !resetsAt.After(at) {
		return nil
	}
	f := Fetcher{StateDir: stateDir}
	rec, _ := f.readLimit(name, at)
	if rec.Windows == nil {
		rec.Windows = map[string]time.Time{}
	}
	rec.Windows[window] = resetsAt.UTC()
	if at.After(rec.At) {
		rec.At = at.UTC()
	}
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return fsutil.AtomicWrite(f.limitPath(name), data, 0o600)
}

// DropLimit forgets account name's record: its limits were another
// login's.
func DropLimit(stateDir, name string) { os.Remove((&Fetcher{StateDir: stateDir}).limitPath(name)) }

// readLimit reads account name's record with only the windows that have
// not reset by now.
func (f *Fetcher) readLimit(name string, now time.Time) (limitRecord, bool) {
	data, err := os.ReadFile(f.limitPath(name))
	if err != nil {
		return limitRecord{}, false
	}
	var rec limitRecord
	if json.Unmarshal(data, &rec) != nil {
		return limitRecord{}, false
	}
	for w, at := range rec.Windows {
		if !at.After(now) {
			delete(rec.Windows, w)
		}
	}
	return rec, len(rec.Windows) > 0
}

// withLimit lays the account's limit record over r: each window in it is
// at 100% until its reset, unless a later lookup lifted it (lifted). With
// nothing else known, the record alone is what r says.
func (f *Fetcher) withLimit(a Account, r Report) Report {
	if r.Status == StatusAPI || r.Status == StatusLoggedOut {
		return r
	}
	rec, ok := f.readLimit(a.Name, f.Now())
	if !ok {
		return r
	}
	lifted := f.lifted(a.Name, rec.At)
	var l Limits
	if r.Limits != nil {
		l = *r.Limits
		l.Scoped = slices.Clone(l.Scoped)
	}
	for _, w := range limitWindows {
		at, ok := rec.Windows[w]
		if !ok || lifted(w) {
			continue
		}
		full := Window{Percent: 100, ResetsAt: at, Known: true}
		label := ""
		switch w {
		case "five_hour":
			l.FiveHour, label = full, "5-hour"
		case "seven_day":
			l.Weekly, label = full, "Weekly"
		default: // a model's week
			i := modelWindow(l.Scoped, WindowModel(w))
			if i < 0 {
				fam := WindowModel(w)
				l.Scoped = append(l.Scoped, Scoped{Name: strings.ToUpper(fam[:1]) + fam[1:]})
				i = len(l.Scoped) - 1
			}
			l.Scoped[i].Window, label = full, l.Scoped[i].Name
		}
		r.LimitedWindows = append(r.LimitedWindows, label)
	}
	if len(r.LimitedWindows) == 0 {
		return r // every window was lifted
	}
	if r.Limits == nil {
		r.At, r.Source = rec.At, "limit"
	}
	r.Limits, r.Limited = &l, rec.At
	return r
}

// modelWindow is the index in scoped of family's week, or -1.
func modelWindow(scoped []Scoped, family string) int {
	return slices.IndexFunc(scoped, func(s Scoped) bool { return strings.Contains(strings.ToLower(s.Name), family) })
}

// lifted tells which windows of account name the newest lookup of the
// usage endpoint (ccam's or Claude Code's, as saved) shows below 100%, when
// it was made liftDelay or more after the limit at limitAt. The status line
// does not count: its time is when the command ran, so a percentage another
// session showed before the limit may come with a later time.
func (f *Fetcher) lifted(name string, limitAt time.Time) func(window string) bool {
	var last *cached
	for _, src := range []string{"live", "claude"} {
		if c, ok := f.readCache(name, src); ok && (last == nil || c.At.After(last.At)) {
			last = &c
		}
	}
	return func(window string) bool {
		if last == nil || last.At.Before(limitAt.Add(liftDelay)) {
			return false
		}
		var w Window
		switch window {
		case "five_hour":
			w = last.Limits.FiveHour
		case "seven_day":
			w = last.Limits.Weekly
		default:
			if i := modelWindow(last.Limits.Scoped, WindowModel(window)); i >= 0 {
				w = last.Limits.Scoped[i].Window
			}
		}
		return w.Known && w.Percent < 100
	}
}
