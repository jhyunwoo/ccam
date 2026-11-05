package cli

import (
	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/next"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// pickAuto chooses the logged-in account by usage (usage.Policy: below
// auto_max_usage first, by auto_strategy; full ones never) among those that
// may take over from from (nil: the most recently used one), and says why.
// With wait it waits for a full account to reset.
func pickAuto(e *env, cfg *config.Config, cands []next.Candidate, from *config.Profile, p usage.Policy, wait string, crossOrg bool) (*config.Profile, *config.Profile, string, error) {
	if from == nil {
		from, _, _ = next.Rotate(cands, "", nil)
	}
	org := newOrgs(e, cfg, crossOrg)
	var profs []*config.Profile
	for _, c := range cands {
		if c.LoggedIn && org.fits(from, c.Profile) {
			profs = append(profs, c.Profile)
		}
	}
	if len(profs) == 0 {
		return nil, nil, "", noRoom(org.none(from, "no logged-in account to pick"))
	}
	d, err := waitFor(wait)
	if err != nil {
		return nil, nil, "", err
	}
	why := ""
	to, err := untilRoom(e, p, d, "no account with known usage and headroom (see `ccam usage`)", func() (*config.Profile, []*usage.Report, error) {
		reports := reportsOf(e, cfg, profs, usage.Options{})
		if best := p.Choose(reports, clock()); best != nil {
			why = p.Why(best, reports, clock())
			prof, _ := cfg.Find(best.Name)
			return prof, reports, nil
		}
		return nil, reports, nil
	})
	return from, to, why, err
}
