package usage

import (
	"fmt"
	"strings"
	"time"
)

// Strategies for Policy.Strategy.
const (
	Headroom = "headroom" // the least weekly usage within nearFive of the lowest 5-hour usage (default)
	Expiring = "expiring" // the soonest 5-hour reset, so its quota is not left unused
	// Balance takes the account whose weekly quota would go unused first:
	// the most left per day until its reset, less for each running session.
	Balance = "balance"
)

// Policy is how `ccam next --auto`, `--fallback` and `ccam pick` choose an
// account from its usage. An account is judged by its 5-hour and weekly
// windows, and by the weekly window of the model named with --model.
type Policy struct {
	// MaxUsage is the usage (1-100) up to which an account is preferred; 0
	// means 100. Accounts at or over it are only chosen when no account is
	// below it, and accounts at 100% never.
	MaxUsage float64
	Strategy string // Headroom (also "") or Expiring
	Model    string // the --model argument, if any
}

// Check reports a policy that cannot be used, such as a typo in config.toml.
func (p Policy) Check() error {
	if p.MaxUsage < 0 || p.MaxUsage > 100 {
		return fmt.Errorf("auto_max_usage must be between 1 and 100, not %g", p.MaxUsage)
	}
	if p.Strategy != "" && p.Strategy != Headroom && p.Strategy != Expiring && p.Strategy != Balance {
		return fmt.Errorf("auto_strategy must be %q, %q or %q, not %q", Headroom, Expiring, Balance, p.Strategy)
	}
	return nil
}

func (p Policy) max() float64 {
	if p.MaxUsage <= 0 {
		return 100
	}
	return p.MaxUsage
}

// windows lists the known windows r is judged by.
func (p Policy) windows(r *Report) []Window {
	var out []Window
	for _, w := range []Window{r.Limits.FiveHour, r.Limits.Weekly} {
		if w.Known {
			out = append(out, w)
		}
	}
	if fam := family(p.Model); fam != "" {
		for _, s := range r.Limits.Scoped {
			if s.Window.Known && strings.Contains(strings.ToLower(s.Name), fam) {
				out = append(out, s.Window)
			}
		}
	}
	return out
}

// family is the model family a --model value names: opus, sonnet, haiku
// or fable, as aliases and full IDs both contain it.
func family(model string) string {
	m := strings.ToLower(model)
	for _, f := range []string{"opus", "sonnet", "haiku", "fable"} {
		if strings.Contains(m, f) {
			return f
		}
	}
	return ""
}

// usable reports whether r has known 5-hour usage and is not logged out.
func usable(r *Report) bool {
	return r != nil && r.Limits != nil && r.Limits.FiveHour.Known && r.Status != StatusLoggedOut
}

// peak is r's highest usage among the windows it is judged by.
func (p Policy) peak(r *Report) float64 {
	top := 0.0
	for _, w := range p.windows(r) {
		top = max(top, w.Percent)
	}
	return top
}

// Full reports whether r is at 100% in a window it is judged by.
func (p Policy) Full(r *Report) bool {
	return r != nil && r.Limits != nil && p.peak(r) >= 100
}

// Choose returns the account to use, or nil when every account is full or
// its usage is unknown. Below MaxUsage the strategy decides; when every
// account is at or over it, the one headroom would choose is taken.
func (p Policy) Choose(reports []*Report, now time.Time) *Report {
	rs, over := p.pool(reports)
	if over {
		return Policy{}.best(rs, now)
	}
	return p.best(rs, now)
}

// pool is what Choose chooses among: the accounts below MaxUsage or, when
// none is, those below 100% (over).
func (p Policy) pool(reports []*Report) (rs []*Report, over bool) {
	var free, overs []*Report
	for _, r := range reports {
		if !usable(r) {
			continue
		}
		switch top := p.peak(r); {
		case top >= 100:
		case top >= p.max():
			overs = append(overs, r)
		default:
			free = append(free, r)
		}
	}
	if len(free) > 0 {
		return free, false
	}
	return overs, true
}

// nearFive is how many points of 5-hour usage headroom counts as about the
// same: of the accounts within it of the lowest, the one with the least of
// its week used is taken.
const nearFive = 3

func (p Policy) headroom() bool { return p.Strategy == "" || p.Strategy == Headroom }

func (p Policy) best(rs []*Report, now time.Time) *Report {
	if p.headroom() {
		// Measured from the lowest, so the choice does not depend on the
		// order of rs: 12% may beat 10%, and 14% 12%, but 14% never 10%.
		if low := lowest(rs); low != nil {
			var near []*Report
			for _, r := range rs {
				if r.Limits.FiveHour.Percent-low.Limits.FiveHour.Percent <= nearFive {
					near = append(near, r)
				}
			}
			rs = near
		}
	}
	var best *Report
	for _, r := range rs {
		if best == nil || p.before(r, best, now) {
			best = r
		}
	}
	return best
}

// before reports whether a is a better choice than b.
func (p Policy) before(a, b *Report, now time.Time) bool {
	switch p.Strategy {
	case Expiring:
		ra, rb := resetIn(a, now), resetIn(b, now)
		if ra != rb {
			return ra < rb
		}
	case Balance:
		ua, ub := urgency(a, now), urgency(b, now)
		if ua != ub {
			return ua > ub
		}
	case "", Headroom:
		// best kept those near the lowest 5-hour usage. An unknown week
		// counts as unused, as in urgency.
		if wa, wb := a.Limits.Weekly.Percent, b.Limits.Weekly.Percent; wa != wb {
			return wa < wb
		}
	}
	return fewer(a, b)
}

// fewer reports whether a has used less of its 5 hours than b, or as much
// and less of its week.
func fewer(a, b *Report) bool {
	if fa, fb := a.Limits.FiveHour.Percent, b.Limits.FiveHour.Percent; fa != fb {
		return fa < fb
	}
	return a.Limits.Weekly.Percent < b.Limits.Weekly.Percent
}

// lowest is the account of rs with the least of its 5 hours used.
func lowest(rs []*Report) *Report {
	var low *Report
	for _, r := range rs {
		if low == nil || fewer(r, low) {
			low = r
		}
	}
	return low
}

// sessionCost is how much of the weekly quota per day a running session
// is taken to use, so balance spreads sessions over the accounts.
const sessionCost = 10

// urgency is the share of r's weekly quota left per day until the week
// resets, less sessionCost for each running session. Without a known week,
// a full week is assumed.
func urgency(r *Report, now time.Time) float64 {
	w := r.Limits.Weekly
	left, days := 100.0, 7.0
	if w.Known {
		left = 100 - w.Percent
		if !w.ResetsAt.IsZero() {
			days = max(w.ResetsAt.Sub(now).Hours()/24, 1.0/24)
		}
	}
	return left/days - sessionCost*float64(r.Sessions)
}

// Why says in a few words why Choose took r from reports.
func (p Policy) Why(r *Report, reports []*Report, now time.Time) string {
	if r == nil || !usable(r) {
		return ""
	}
	five := fmt.Sprintf("%.0f%% of its 5 hours used", r.Limits.FiveHour.Percent)
	if top := p.peak(r); top >= p.max() {
		if near := p.nearWhy(r, reports, five); near != "" {
			return fmt.Sprintf("every account is at %.0f%% or more; %s", p.max(), near)
		}
		return fmt.Sprintf("every account is at %.0f%% or more; it has the lowest use, %s", p.max(), five)
	}
	if p.headroom() {
		if near := p.nearWhy(r, reports, five); near != "" {
			return near
		}
	}
	switch p.Strategy {
	case Expiring:
		if at := r.Limits.FiveHour.ResetsAt; !at.IsZero() {
			return fmt.Sprintf("its 5-hour window resets first, in %s (%s)", Duration(at.Sub(now)), five)
		}
	case Balance:
		w := r.Limits.Weekly
		why := "its weekly use is unknown"
		if w.Known {
			why = fmt.Sprintf("%.0f%% of its week is left", 100-w.Percent)
			if !w.ResetsAt.IsZero() {
				why += ", resetting in " + Duration(w.ResetsAt.Sub(now))
			}
		}
		if r.Sessions > 0 {
			why += fmt.Sprintf("; %d running", r.Sessions)
		}
		return why
	}
	return five + ", the least"
}

// nearWhy says why headroom took r over the account with the least of its
// 5 hours used, or "" when r is that account (or ties with it).
func (p Policy) nearWhy(r *Report, reports []*Report, five string) string {
	rs, _ := p.pool(reports)
	low := lowest(rs)
	if low == nil || low.Limits.FiveHour.Percent >= r.Limits.FiveHour.Percent {
		return ""
	}
	week := "unknown, counted as unused"
	if r.Limits.Weekly.Known {
		week = fmt.Sprintf("%.0f%%", r.Limits.Weekly.Percent)
	}
	return fmt.Sprintf("%s, within %d points of %s's %.0f%%, with less of its week used (%s against %.0f%%)",
		five, nearFive, low.Name, low.Limits.FiveHour.Percent, week, low.Limits.Weekly.Percent)
}

// Use says in a few words how far r is by what p goes by, as Why does but
// without comparing it with other accounts: for an account taken while the
// one Choose takes was passed over.
func (p Policy) Use(r *Report, now time.Time) string {
	if r == nil || !usable(r) {
		return ""
	}
	five := fmt.Sprintf("%.0f%% of its 5 hours used", r.Limits.FiveHour.Percent)
	if p.peak(r) < p.max() {
		switch p.Strategy {
		case Expiring:
			if at := r.Limits.FiveHour.ResetsAt; !at.IsZero() {
				return fmt.Sprintf("its 5-hour window resets in %s (%s)", Duration(at.Sub(now)), five)
			}
		case Balance:
			return p.Why(r, nil, now) // it names no other account
		}
	}
	return five
}

// resetIn is the time until r's 5-hour window resets; a window that is not
// running (no reset time) sorts last.
func resetIn(r *Report, now time.Time) time.Duration {
	at := r.Limits.FiveHour.ResetsAt
	if at.IsZero() {
		return 1<<63 - 1
	}
	return at.Sub(now)
}

// FreeAt returns when the first full account gets room again: the reset of
// its latest window at 100%. It is zero when no full account has a known
// reset time.
func (p Policy) FreeAt(reports []*Report, now time.Time) time.Time {
	var first time.Time
	for _, r := range reports {
		if !usable(r) {
			continue
		}
		var at time.Time
		for _, w := range p.windows(r) {
			if w.Percent >= 100 && w.ResetsAt.After(at) {
				at = w.ResetsAt
			}
		}
		if at.After(now) && (first.IsZero() || at.Before(first)) {
			first = at
		}
	}
	return first
}
