package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/session"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// exitNoRoom is EX_TEMPFAIL from sysexits.h: no account has room now, so a
// script may try again later.
const exitNoRoom = 75

// codedError ends ccam with its message and a specific exit code.
type codedError struct {
	code int
	msg  string
}

func (e codedError) Error() string { return e.msg }

func noRoom(msg string) error { return codedError{exitNoRoom, msg} }

// choiceFlags are the flags that override how an account is chosen.
type choiceFlags struct {
	maxUsage float64
	strategy string
	wait     string
	crossOrg bool
}

// DefaultWait is how long a bare --wait waits at most: one 5-hour window.
const DefaultWait = "5h"

func (c *choiceFlags) add(cmd *cobra.Command) {
	cmd.Flags().Float64Var(&c.maxUsage, "max-usage", 0, "prefer accounts below this usage, 1-100 (config: auto_max_usage)")
	cmd.Flags().StringVar(&c.strategy, "strategy", "", "headroom (least weekly usage within 3 points of the lowest 5-hour usage), expiring (soonest reset) or balance (weekly quota that would go unused first) (config: auto_strategy)")
	cmd.Flags().StringVar(&c.wait, "wait", "", "when no account has room, wait for one, at most this long (default "+DefaultWait+")")
	cmd.Flags().Lookup("wait").NoOptDefVal = DefaultWait
	cmd.Flags().BoolVar(&c.crossOrg, "cross-org", false, "also choose accounts of other organizations (config: cross_org)")
}

// spacedValue catches a value written after a space for a flag whose value
// is optional (--wait 30m, --watch 10): pflag takes such a value only after
// "=", so the flag keeps its default and the value would become an
// argument (with next, claude's first prompt). Words after "--" are left
// alone.
func spacedValue(cmd *cobra.Command, args []string, name, def string) error {
	if n := cmd.ArgsLenAtDash(); n >= 0 {
		args = args[:n]
	}
	f := cmd.Flags().Lookup(name)
	if f == nil || !f.Changed || f.Value.String() != def || len(args) == 0 {
		return nil
	}
	if name != "watch" {
		return waitHint(args[0])
	}
	if _, err := strconv.Atoi(args[0]); err == nil {
		return fmt.Errorf("--%s takes its value after '=': --%s=%s", name, name, args[0])
	}
	return nil
}

// waitHint says how to write v, a word after a bare --wait that looks like
// its value, so that waitFor takes it; nil when v does not look like one.
func waitHint(v string) error {
	d, err := time.ParseDuration(v)
	_, isNumber := strconv.Atoi(v)
	switch {
	case err == nil && d > 0:
		return fmt.Errorf("--wait takes its value after '=': --wait=%s", v)
	case err == nil, isNumber == nil:
		// A number without a unit, or a zero that waitFor refuses: which
		// unit was meant is not known.
		return fmt.Errorf("--wait takes a duration after '=', such as --wait=30m or --wait=2h")
	}
	return nil
}

// orgs keeps automatic choices within one organization: accounts of one
// team or enterprise organization, or all personal plans, stand in for
// each other. A folder logged in to the same Claude account as another
// shares its limits, so it never stands in for it.
type orgs struct {
	e     *env
	cross bool
	ids   map[string]*config.Identity
	live  map[string]bool // the identity is what .claude.json says now
}

func newOrgs(e *env, cfg *config.Config, cross bool) *orgs {
	return &orgs{e: e, cross: cross || cfg.CrossOrg, ids: map[string]*config.Identity{}, live: map[string]bool{}}
}

// of is who p is logged in as (profile.LoggedInAs), read once.
func (o *orgs) of(p *config.Profile) *config.Identity {
	if id, ok := o.ids[p.Name]; ok {
		return id
	}
	id, live := profile.LoggedInAs(o.e.paths, p)
	o.ids[p.Name], o.live[p.Name] = id, live
	return id
}

// takeover lists the accounts that may take over from the account named
// from, by the rules of next --auto.
func takeover(e *env, cfg *config.Config, from string) []string {
	org := newOrgs(e, cfg, false)
	f, ok := cfg.Find(from)
	var out []string
	for i := range cfg.Profiles {
		if p := &cfg.Profiles[i]; !ok || org.fits(f, p) {
			out = append(out, p.Name)
		}
	}
	return out
}

// fits reports whether cand may take over from from (nil: anything goes).
func (o *orgs) fits(from, cand *config.Profile) bool {
	if from == nil || from.Name == cand.Name {
		return true
	}
	a, b := o.of(from), o.of(cand)
	if config.SameAccount(a, b) {
		return false
	}
	return o.cross || config.SamePool(a, b)
}

// none is the message when the rule left no account: it names the
// organization and how to look further.
func (o *orgs) none(from *config.Profile, msg string) string {
	if from == nil || o.cross {
		return msg
	}
	if label := o.of(from).OrgLabel(); label != "" {
		return fmt.Sprintf("%s among %s's (%s); --cross-org also looks at the others", msg, from.Name, label)
	}
	return msg
}

// policy is the config's choice rule with the flags and claude args on top.
func (c choiceFlags) policy(cfg *config.Config, args []string) (usage.Policy, error) {
	p := usage.Policy{MaxUsage: cfg.AutoMaxUsage, Strategy: cfg.AutoStrategy, Model: modelArg(args)}
	if c.maxUsage != 0 {
		p.MaxUsage = c.maxUsage
	}
	if c.strategy != "" {
		p.Strategy = c.strategy
	}
	if err := p.Check(); err != nil {
		return p, err
	}
	return p, nil
}

// waitFor parses a --wait value; "" means do not wait.
func waitFor(v string) (time.Duration, error) {
	if v == "" {
		return 0, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("--wait takes a duration such as 30m or 2h, not %q", v)
	}
	return d, nil
}

// modelArg returns the --model value among claude args.
func modelArg(args []string) string {
	for i, a := range args {
		if a == "--" {
			break
		}
		if v, ok := strings.CutPrefix(a, "--model="); ok {
			return v
		}
		if a == "--model" && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

// Clock and sleep, replaced by tests.
var (
	clock = time.Now
	sleep = time.Sleep
)

// afterReset is how long ccam waits past a reset before looking again: the
// endpoint takes a moment to show the new window.
const afterReset = 30 * time.Second

// lookAgain is how long ccam waits when no full account says when it resets.
const lookAgain = 10 * time.Minute

// untilRoom calls try until it chooses an account. When try finds none and
// wait is set, it sleeps until the first full account resets (per p.FreeAt)
// and tries again, until wait runs out; then, or without wait, it fails with
// exit code 75 and the message none.
func untilRoom(e *env, p usage.Policy, wait time.Duration, none string, try func() (*config.Profile, []*usage.Report, error)) (*config.Profile, error) {
	deadline := clock().Add(wait)
	for {
		prof, reports, err := try()
		if err != nil || prof != nil {
			return prof, err
		}
		now := clock()
		if wait <= 0 {
			return nil, noRoom(none)
		}
		at := p.FreeAt(reports, now)
		if at.IsZero() {
			at = now.Add(lookAgain)
		} else {
			at = at.Add(afterReset)
		}
		if at.After(deadline) {
			return nil, noRoom(none + fmt.Sprintf(" (and none gets room before --wait ends at %s)", deadline.Local().Format("15:04")))
		}
		fmt.Fprintf(e.stderr, "ccam: no account has room; waiting until %s\n", at.Local().Format("15:04"))
		sleep(at.Sub(now))
	}
}

// reportsOf fetches the usage of profs at once, in their order.
func reportsOf(e *env, cfg *config.Config, profs []*config.Profile, opt usage.Options) []*usage.Report {
	return reportsWith(context.Background(), e, e.fetcher(), cfg, profs, opt)
}

// reportsWith is reportsOf with the fetcher f, until ctx ends. The caller
// makes f, so this may run on a goroutine that outlives it (autoTarget).
func reportsWith(ctx context.Context, e *env, f *usage.Fetcher, cfg *config.Config, profs []*config.Profile, opt usage.Options) []*usage.Report {
	all := usageAccounts(e, cfg)
	var accts []usage.Account
	for _, p := range profs {
		for _, a := range all {
			if a.Name == p.Name {
				accts = append(accts, a)
			}
		}
	}
	reports := f.FetchAll(ctx, accts, opt, nil)
	out := make([]*usage.Report, len(reports))
	for i := range reports {
		reports[i].Sessions = session.Count(launcher.SessionDir(e.paths, reports[i].Name))
		out[i] = &reports[i]
	}
	return out
}
