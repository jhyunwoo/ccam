package cli

import (
	"cmp"
	"context"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/session"
	"github.com/jhyunwoo/ccam/internal/transcript"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// autoDeadline bounds the choice of an account after a usage limit,
// counted from the hook's start: the plugin gives the hook 10 s, and the
// limit row may take limitPoll to show. Tests change it.
var autoDeadline = 6 * time.Second

// nearReset is how soon a reset leaves a session to Claude Code's own wait
// rather than to a restart as another account (D9). Tests change it.
var nearReset = 15 * time.Minute

// limitFacts is what the stop-failure hook knows of a session that stopped
// at a usage limit of a plan's window, which it recorded, for eligible.
type limitFacts struct {
	on       bool             // auto_switch = true, and not on Windows
	kind     string           // handoff.Classify of Claude Code's message
	resetsIn time.Duration    // until the window resets; 0 when unknown
	remote   bool             // the session is connected to Remote Control (D11)
	relayed  bool             // CCAM_RELAY is this protocol's ("2:<pid>") and CLAUDE_PID is set
	child    bool             // CLAUDE_PID is the relay's own claude, still running
	plan     *config.Identity // who the account is logged in as (orgs.of)
}

// stay says why auto_switch leaves a session where it is.
type stay int

const (
	stayOff         stay = iota // off, or no limit it moves: the usual notification
	stayNear                    // the window resets within nearReset
	stayRemote                  // connected to Remote Control (D11)
	stayNotRelayed              // no relay of this protocol, or no CLAUDE_PID
	stayUnknownPlan             // the account's plan is unknown (D7 b)
	moves                       // eligible
)

// eligible says whether auto_switch may move a session that stopped at a
// usage limit, or why it leaves it there (plan 4.3, D7 b, D11). A reset
// near at hand comes first, as Remote Control then needs a word of its own;
// Remote Control before the relay, which `claude remote-control` sessions
// never have; a claude that is not the relay's own, or has ended, gets the
// usual notification.
func eligible(f limitFacts) stay {
	switch {
	case !f.on, f.kind != handoff.Account && f.kind != handoff.Model, f.resetsIn <= 0:
		return stayOff
	case f.resetsIn <= nearReset:
		return stayNear
	case f.remote:
		return stayRemote
	case !f.relayed:
		return stayNotRelayed
	case !f.child:
		return stayOff
	case f.plan == nil || f.plan.OrgType == "":
		return stayUnknownPlan
	}
	return moves
}

// stayNote is the notification for a session that eligible leaves where it
// is, or "" for the usual one. hit begins it (handoff.HitLimit); label
// names the window.
func stayNote(s stay, f limitFacts, account, label, hit string) string {
	switch s {
	case stayNear:
		near := handoff.HitLimit(account, label, "") + ", which resets in " + usage.Duration(f.resetsIn)
		if f.remote {
			return near + ". Remote Control is on, so Claude Code does not wait on its own: send a prompt after the reset."
		}
		return near + "; ccam leaves this session to wait for it."
	case stayRemote:
		return hit + ". This session is connected to Remote Control and would not reconnect as another account, so ccam leaves it here."
	case stayNotRelayed:
		return hit + ". This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`."
	case stayUnknownPlan:
		return hit + ". ccam does not know " + account + "'s plan, so auto_switch leaves this session."
	}
	return ""
}

// target is the account autoTarget chose, or what the notification says
// when there is none.
type target struct {
	to        *config.Profile // the account to carry the conversation on as; nil: none
	why       string          // why to (usage.Policy.Why; how old saved values are)
	full      bool            // without to: an account is known to be full,
	first     string          // the one whose full window resets first,
	freeAt    time.Time       // and when
	untrusted []string        // accounts with room that never trusted the folder
	credits   []string        // accounts with room and usage credits on (D4)
	agent     string          // the agent the conversation runs as; "" for none
	// Accounts with room that do not load the agents of the account at its
	// limit the same way (sameAgents): they may lack agent, or define it
	// otherwise.
	otherAgents []string
}

// autoTarget chooses the account to carry on a conversation of folder that
// stopped at a usage limit of from (plan 4.3, D7 b): another account of
// from's pool, whatever cross_org says, logged in to another Claude account
// of a known plan, no API account, that trusts folder (mcp.SessionTrusts),
// has no usage credits on (D4), a login Claude Code can renew and, when
// the conversation runs as agent (not ""), loads from's agents the same
// way (sameAgents); chosen by p, the accounts picked within recentPick last
// (spread). Usage is looked up live, or taken as saved with usage_source =
// "claude": a hook starts no Claude Code. The choice is given up when ctx
// ends; macOS Keychain reads do not follow ctx, hence the goroutine.
func autoTarget(ctx context.Context, e *env, cfg *config.Config, from *config.Profile, folder, agent string, p usage.Policy) (target, error) {
	org := newOrgs(e, cfg, false)
	org.cross = false // cross_org never applies to auto_switch (D7)
	var profs []*config.Profile
	untrusted, otherAgents := map[string]bool{}, map[string]bool{}
	for i := range cfg.Profiles {
		c := &cfg.Profiles[i]
		if c.Name == from.Name || c.API != nil || !org.fits(from, c) {
			continue
		}
		if id := org.of(c); id == nil || id.OrgType == "" {
			continue // an unknown plan fits any pool
		}
		if !mcp.SessionTrusts(profile.Dir(e.paths, c), folder) {
			untrusted[c.Name] = true
		}
		// Resuming there, Claude Code would go on without the agent's tool
		// limits if c lacks the agent, or with c's own definition of it,
		// and nobody is there to see it.
		if agent != "" && !sameAgents(profile.Dir(e.paths, from), profile.Dir(e.paths, c)) {
			otherAgents[c.Name] = true
		}
		profs = append(profs, c)
	}
	opt := usage.Options{}
	if cfg.UsageSource == "claude" {
		opt.Source = "cache"
	}
	f, now := e.fetcher(), clock()
	done := make(chan target, 1)
	go func() {
		tg := target{agent: agent}
		var ok []*usage.Report
		for _, r := range reportsWith(ctx, e, f, cfg, profs, opt) {
			room := p.Choose([]*usage.Report{r}, now) != nil
			switch {
			case r.Limits != nil && r.Limits.Credits: // past its limits it would go on at API rates
				if room {
					tg.credits = append(tg.credits, r.Name)
				}
			case untrusted[r.Name]:
				if room {
					tg.untrusted = append(tg.untrusted, r.Name)
				}
			case otherAgents[r.Name]:
				if room {
					tg.otherAgents = append(tg.otherAgents, r.Name)
				}
			default:
				ok = append(ok, r)
			}
		}
		for {
			best := spread(e, p, ok, now)
			if best == nil {
				break
			}
			// Nobody may be there to log in again: an account whose login
			// Claude Code cannot renew is left out.
			if to, _ := cfg.Find(best.Name); creds.Has(profile.Dir(e.paths, to), e.paths.Home) {
				tg.to, tg.why = to, p.Why(best, ok, now)
				if top := p.Choose(ok, now); top != best {
					// spread passed over top, which another session took a
					// moment ago: Why would compare best with every account.
					tg.why = p.Use(best, now) + "; " + top.Name + " was taken by another session a moment ago"
				}
				if best.Source != "live" && !best.At.IsZero() && best.At.Before(now) {
					tg.why += "; saved " + usage.Duration(now.Sub(best.At)) + " ago"
				}
				done <- tg
				return
			}
			ok = slices.DeleteFunc(ok, func(r *usage.Report) bool { return r == best })
		}
		for _, r := range ok {
			tg.full = tg.full || p.Full(r)
			if at := p.FreeAt([]*usage.Report{r}, now); !at.IsZero() && (tg.first == "" || at.Before(tg.freeAt)) {
				tg.first, tg.freeAt = r.Name, at
			}
		}
		done <- tg
	}()
	var tg target
	select {
	case tg = <-done:
	case <-ctx.Done():
	}
	if err := ctx.Err(); err != nil {
		return target{}, err // a choice made past the deadline is not used
	}
	return tg, nil
}

// poolName names the accounts that stand in for one logged in as id (its
// pool) in notifications: "personal plans", "accounts in Acme" ("in this
// organization" without a name), or "accounts" when its plan is unknown.
func poolName(id *config.Identity) string {
	switch id.Pool() {
	case config.Personal:
		return "personal plans"
	case "":
		return "accounts"
	}
	return "accounts in " + cmp.Or(id.OrgName, "this organization")
}

// note is the notification when tg has no account: the accounts with room
// that never trusted the folder, have usage credits on or do not load the
// agents of from the same way, or that nobody has room now and who frees up
// first; "" when nothing is known of the others (the usual notification).
// from, logged in as id, hit its limit, which resets at reset.
func (tg target) note(hit, from string, id *config.Identity, reset, now time.Time) string {
	if n := len(tg.untrusted); n > 0 {
		verb, as := "has", tg.untrusted[0]
		if n > 1 {
			verb, as = "have", "each of them"
		}
		return hit + ". " + andList(tg.untrusted) + " " + verb + " not trusted this folder yet, so the conversation was not moved there. Open the folder once as " + as + " to allow it."
	}
	if n := len(tg.credits); n > 0 {
		verb := "has room but its"
		if n > 1 {
			verb = "have room but their"
		}
		return hit + ". " + andList(tg.credits) + " " + verb + " usage credits are on, so auto_switch did not move the conversation there; /ccam:switch " + tg.credits[0] + " moves it."
	}
	if n := len(tg.otherAgents); n > 0 {
		verb := "does"
		if n > 1 {
			verb = "do"
		}
		return hit + ". This conversation runs as agent " + strconv.Quote(tg.agent) + ", and " + andList(tg.otherAgents) + " " + verb + " not share " + from +
			"'s agents, plugins and settings, so auto_switch did not move it there."
	}
	if !tg.full {
		return ""
	}
	msg := hit + ". None of your other " + poolName(id) + " has room now"
	first, at := tg.first != "" && tg.freeAt.Before(reset), handoff.FormatReset(tg.freeAt, now)
	if first {
		msg += "; " + tg.first + " frees up first, at " + at
	}
	// Claude Code waits for a reset within a day; beyond, it shows a menu.
	waits := reset.Sub(now) <= 24*time.Hour
	switch {
	case waits && first:
		return msg + ". Claude Code waits for " + from + " here; /ccam:switch " + tg.first + " after " + at + " moves it sooner."
	case waits:
		return msg + ". Claude Code waits for " + from + " here."
	case first:
		return msg + ". /ccam:switch " + tg.first + " after " + at + " moves it there."
	}
	return msg + "."
}

// andList joins account names as "a, b and c".
func andList(list []string) string {
	n := len(list)
	if n == 1 {
		return list[0]
	}
	return strings.Join(list[:n-1], ", ") + " and " + list[n-1]
}

// autoSwitch ends Claude Code after a usage limit that the hook recorded,
// when auto_switch is on and another account has room (plan 4.3, 4.4), so
// that the relay carries the conversation on as that account. moved:
// Claude Code is ending, and the hook says nothing more. Otherwise say is
// the notification to show instead of the usual one ("" keeps it). t is
// the session's note as saved, hit the start of every message.
func autoSwitch(e *env, in hookInput, t handoff.Ticket, lim transcript.Limit, start time.Time, hit string) (moved bool, say string) {
	cfg, err := e.load()
	if err != nil {
		return false, ""
	}
	now := clock()
	f := limitFacts{on: cfg.AutoSwitch && runtime.GOOS != "windows", kind: handoff.Classify(cmp.Or(in.LastMessage, lim.Text))}
	if !lim.ResetsAt.IsZero() {
		f.resetsIn = lim.ResetsAt.Sub(now)
	}
	from, _ := cfg.Find(t.Account)
	relay, protocol := launcher.RelayPID(os.Getenv("CCAM_RELAY"))
	claude, _ := strconv.Atoi(os.Getenv("CLAUDE_PID"))
	var child int
	var started string
	if f.on { // nothing more is read while auto_switch is off
		// Claude Code keeps CLAUDE_CODE_BRIDGE_SESSION_ID in its environment,
		// which hooks inherit, while Remote Control is connected.
		f.remote = os.Getenv("CLAUDE_CODE_BRIDGE_SESSION_ID") != "" || transcript.RemoteControl(in.Transcript)
		f.relayed = protocol && claude > 0
		if f.relayed {
			var ok bool
			if child, ok = launcher.RelayChild(claude, relay); ok {
				started, f.child = session.StartTime(claude)
			}
		}
		if from != nil {
			f.plan = newOrgs(e, cfg, false).of(from)
		}
	}
	if s := eligible(f); s != moves {
		return false, stayNote(s, f, t.Account, handoff.WindowLabel(lim.Window), hit)
	}
	p := usage.Policy{MaxUsage: cfg.AutoMaxUsage, Strategy: cfg.AutoStrategy, Model: lim.Model}
	if fam, ok := strings.CutPrefix(lim.Window, "seven_day_"); ok {
		p.Model = fam // the Opus or Sonnet week: the next account's must have room
	}
	if p.Check() != nil {
		return false, ""
	}
	ctx, cancel := context.WithDeadline(context.Background(), start.Add(autoDeadline))
	defer cancel()
	tg, err := autoTarget(ctx, e, cfg, from, t.Cwd, conversationAgent(profile.Dir(e.paths, from), t.Cwd, in.Transcript), p)
	switch {
	case err != nil:
		return false, "" // too slow (a lookup or the Keychain hangs)
	case tg.to == nil:
		return false, tg.note(hit, t.Account, f.plan, lim.ResetsAt, now)
	}
	// Not once the person has moved on, nor once the process is not the one
	// looked at (it ended, or its pid is another's now). A subagent's limit
	// is not in the session's conversation, only the prompts are.
	prompted := transcript.PromptedAfter
	if in.AgentID != "" {
		prompted = transcript.PromptedSince
	}
	if prompted(in.Transcript, lim.At) {
		return false, ""
	}
	if c, ok := launcher.RelayChild(claude, relay); !ok || c != child {
		return false, ""
	}
	if s, ok := session.StartTime(claude); !ok || s != started {
		return false, ""
	}
	// Nor over a /ccam:switch since the limit (Claude Code records it as no
	// prompt), or another hook's move.
	path := handoff.Path(e.paths.State, t.Account, t.SessionID)
	if cur, ok := handoff.LoadFile(e.paths.State, path); ok && (cur.Reason == handoff.Switch && cur.At.After(lim.At) || moving(cur)) {
		return false, ""
	}
	plain := t
	t.Relay, t.Child, t.Stopped, t.Chosen, t.Why = os.Getenv("CCAM_RELAY"), child, true, tg.to.Name, tg.why
	if handoff.Save(e.paths.State, t) != nil {
		return false, ""
	}
	// Other sessions that stop now take another account first (spread).
	_ = fsutil.WriteFast(pickPath(e, tg.to.Name), []byte(now.UTC().Format(time.RFC3339)))
	if launcher.Terminate(claude, relay, child, started) != nil {
		// No hop: the person carries on by hand, unless /ccam:switch saved
		// its note over this one meanwhile.
		if cur, ok := handoff.LoadFile(e.paths.State, path); !ok || cur.Reason == handoff.Limit && cur.Chosen == t.Chosen {
			_ = handoff.Save(e.paths.State, plain)
		}
		return false, ""
	}
	return true, ""
}

// subagentHolds reports whether a subagent's limit at window holds up the
// main conversation, whose last reply came from model: a window of the
// whole account does, the week of a model (seven_day_opus,
// seven_day_sonnet) only when the conversation runs on that model, and
// otherwise goes on as it is.
func subagentHolds(window, model string) bool {
	family, ok := strings.CutPrefix(window, "seven_day_")
	return !ok || strings.Contains(model, family)
}

// moving reports whether the note t says that this session's Claude Code
// is being ended for the relay to carry the conversation on (another
// stop-failure hook's move, or /ccam:switch): stopped by a hook under this
// relay, as this CLAUDE_PID or the relay's child for it. The relay follows
// that note, so nothing is saved over it. A note of an earlier Claude Code
// of the relay names another process.
func moving(t *handoff.Ticket) bool {
	claude, _ := strconv.Atoi(os.Getenv("CLAUDE_PID"))
	if !t.Stopped || t.Relay != os.Getenv("CCAM_RELAY") || claude <= 0 {
		return false
	}
	// claude itself, also once it has ended (RelayChild no longer knows it
	// then), or its parent when Claude Code started itself again.
	relay, _ := launcher.RelayPID(t.Relay)
	child, ok := launcher.RelayChild(claude, relay)
	return t.Child == claude || ok && t.Child == child
}
