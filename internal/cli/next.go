package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/next"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/transcript"
	"github.com/jhyunwoo/ccam/internal/usage"
)

func nextCmd(e *env) *cobra.Command {
	var cont, auto, brief bool
	var noteFile string
	var choice choiceFlags
	cmd := &cobra.Command{
		Use:   "next [--auto] [--continue [--brief]] [-- claude args]",
		Short: "Run the next account (--continue: pick up the last conversation)",
		Long: "Opens Claude Code as the account numbered after the last one used.\n" +
			"--continue resumes this folder's last conversation: with claude --resume <id> when the account\n" +
			"shares its history, otherwise Claude Code copies it over (--resume <file> --fork-session).\n" +
			"It names the agent the conversation runs as (--agent), which a copy does not bring back.\n" +
			"--continue follows /ccam:switch to the account it named. Unless that account, or the one ccam\n" +
			"picks when it named none, and the conversation's are seats in one Team or Enterprise organization\n" +
			"or both personal plans (a login whose plan ccam does not know counts as either), ccam asks first,\n" +
			"and with no terminal to ask on the conversation stays.\n" +
			"--brief, with --continue, starts a new conversation instead, which first reads a summary ccam\n" +
			"writes from the saved one without a model: your first and last requests, the summary Claude Code\n" +
			"wrote when it last compacted it, Claude's last reply, the todo list, the files Claude changed and\n" +
			"git status. That costs far less than resuming a long conversation, but details are left out.\n" +
			"Claude says where the work stands and waits, or follows a prompt you give after --. --dry-run\n" +
			"prints the summary instead of saving it. Claude Code's own --brief flag goes after -- as well.\n" +
			"--auto picks by usage: the lowest 5-hour usage, then weekly, among accounts below auto_max_usage\n" +
			"(config.toml, default 100), skipping full ones; --max-usage and --strategy override\n" +
			"the config, and --wait waits for an account to reset when none has room. With no\n" +
			"account to pick, ccam exits with code 75.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			var own *handoff.Ticket
			if noteFile != "" {
				// The relay's hop: that session's note, followed as it says
				// with nobody there to ask, except that a /ccam:switch out of
				// the organization pool waits for the person's yes
				// (confirmSwitch).
				t, ok := handoff.LoadFile(e.paths.State, noteFile)
				if ok {
					_, ok = cfg.Find(t.Account)
				}
				if !ok || !t.Good(time.Now()) {
					return errors.New("this session's note is gone; run `ccam next --continue` to choose the conversation")
				}
				own = t
				launcher.Unattended = true
			}
			// Logins only count where ccam chooses: a hop to a named account
			// does not ask the Keychain about every account.
			hopTo := own != nil && (own.Target != "" || own.Chosen != "")
			cands := candidates(e, cfg, !hopTo)
			var from, to *config.Profile
			var conv *conversation
			if cont {
				// Continue from the account that has the conversation.
				cwd, _ := os.Getwd()
				if conv, err = lastConversation(e, cfg, cands, cwd, own); err != nil {
					return err
				}
				from = conv.owner
			}
			why := "" // why --auto chose to
			var t *handoff.Ticket
			if conv != nil {
				t = conv.ticket
			}
			if t != nil && t.Reason == handoff.Limit && !e.dryRun {
				recordLimit(e, t)
			}
			if name := noteTarget(t, own != nil); name != "" {
				// /ccam:switch named the account, or the hook chose it.
				p, ok := cfg.Find(name)
				switch {
				case ok:
					if err := confirmSwitch(e, cfg, conv, p, false); err != nil {
						return err
					}
					to = p
				case own != nil:
					return fmt.Errorf("the note names %s, which is no longer an account; run `ccam next --continue` to choose one", name)
				default:
					fmt.Fprintf(e.stderr, "ccam: the note names %s, which is no longer an account; choosing as usual\n", name)
				}
			}
			picked := to == nil // ccam picks the account below
			switch {
			case to != nil:
			case auto:
				p, perr := choice.policy(cfg, args)
				if perr != nil {
					return perr
				}
				if p.Model == "" && own != nil {
					// A hop's arguments never name the model (resuming brings
					// it back): its weekly window counts by the conversation's.
					p.Model = transcript.Model(own.Transcript)
				}
				if t != nil {
					// The session asked to move on, or hit a limit its usage
					// may not show yet.
					cands = slices.DeleteFunc(slices.Clone(cands), func(c next.Candidate) bool { return c.Profile.Name == t.Account })
				}
				from, to, why, err = pickAuto(e, cfg, cands, from, p, choice.wait, choice.crossOrg)
			default:
				var fits func(*config.Profile) bool
				name := ""
				if from != nil {
					org := newOrgs(e, cfg, false)
					fits = func(p *config.Profile) bool { return org.fits(from, p) }
					name = from.Name
				}
				from, to, err = next.Rotate(cands, name, fits)
			}
			if err != nil {
				return err
			}
			if picked && t != nil {
				// /ccam:switch named no account (or one that is gone), and
				// cross_org lets ccam's pick leave the pool: such a move asks
				// as well.
				if err := confirmSwitch(e, cfg, conv, to, true); err != nil {
					return err
				}
			}
			agentNote := ""
			if cont && !brief {
				// The agent the conversation runs as goes on with it, also
				// when a person continues it: a resume that copies it into
				// to does not bring the agent back (Claude Code 2.1.314).
				// The relay's hop passes the session's own --agent on.
				args = launcher.WithAgent(args, conversationAgent(profile.Dir(e.paths, conv.owner), conv.cwd, conv.path))
				if agentNote, err = agentMove(e, own, conv, to, args); err != nil {
					return err
				}
			}
			switch {
			case from != nil && why != "":
				fmt.Fprintf(e.stderr, "→ %s → %s (%s)\n", from.Name, to.Name, why)
			case from != nil:
				fmt.Fprintf(e.stderr, "→ %s → %s\n", from.Name, to.Name)
			}
			if cont {
				var note string
				if brief {
					if args, note, err = briefArgs(e, conv, args); err != nil {
						return err
					}
				} else {
					prompt := ""
					if own != nil && own.Reason == handoff.Limit && own.Chosen != "" {
						prompt = handoff.ContinuePrompt // the relay's hop after a limit (D3)
						if own.Subagent != "" {
							prompt = handoff.SubagentPrompt
						}
					}
					args, note = continueArgs(e, cfg, conv, to, prompt, args)
				}
				fmt.Fprintf(e.stderr, "  %s\n", note)
				if agentNote != "" {
					fmt.Fprintf(e.stderr, "  %s\n", agentNote)
				}
				signInNote(e, from, to)
			}
			var used func()
			if t != nil {
				// Used up once claude is about to start, not before: a launch
				// that fails leaves it for the next try.
				used = func() { handoff.Remove(t.File()) }
			}
			return launcher.LaunchThen(e.paths, cfg, to.Name, args, e.dryRun, e.stdout, used)
		},
	}
	cmd.Flags().BoolVarP(&cont, "continue", "c", false, "continue this folder's last conversation on the next account")
	cmd.Flags().StringVar(&noteFile, "note", "", "with --continue: continue the session that left this note (used by the relay)")
	_ = cmd.Flags().MarkHidden("note")
	cmd.Flags().BoolVar(&brief, "brief", false, "with --continue: start a new conversation that reads ccam's summary of the last one instead of resuming it")
	cmd.Flags().BoolVar(&auto, "auto", false, "pick the account by usage (most room)")
	choice.add(cmd)
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		for _, f := range []string{"max-usage", "strategy", "wait", "cross-org"} {
			if cmd.Flags().Changed(f) && !auto {
				return fmt.Errorf("--%s works with --auto", f)
			}
		}
		for _, f := range []string{"note", "brief"} {
			if cmd.Flags().Changed(f) && !cont {
				return fmt.Errorf("--%s works with --continue", f)
			}
		}
		return spacedValue(cmd, args, "wait", DefaultWait)
	}
	return cmd
}

// signInNote says which MCP servers the conversation had a sign-in for
// that the account it moves to loads without one: a sign-in (OAuth) holds
// for its account only, and ccam never copies one. It looks at the group's
// servers and, when both share their plugins, the plugins'.
func signInNote(e *env, from, to *config.Profile) {
	// The relay's hop asks no Keychain (see candidates): on macOS the
	// logins are there.
	if from == nil || from.Name == to.Name || launcher.Unattended && runtime.GOOS == "darwin" {
		return
	}
	had := mcp.SignedIn(profile.Dir(e.paths, from), e.paths.Home)
	if len(had) == 0 {
		return
	}
	has := mcp.SignedIn(profile.Dir(e.paths, to), e.paths.Home)
	var group map[string]json.RawMessage
	if f := to.Share.Get("mcp"); f != config.Local {
		if store, err := mcp.Load(mcp.Path(e.paths.GroupDir(f))); err == nil {
			group = store.Servers
		}
	}
	plugins := to.Share.Get("extensions") != config.Local && to.Share.Get("extensions") == from.Share.Get("extensions")
	var need []string
	for _, n := range slices.Sorted(maps.Keys(had)) {
		_, shared := group[n]
		if !has[n] && (shared || plugins && strings.HasPrefix(n, "plugin:")) {
			need = append(need, n)
		}
	}
	if len(need) > 0 {
		fmt.Fprintf(e.stderr, "  %s is not signed in to MCP servers that %s is signed in to: %s (sign in with /mcp there; sign-ins are not shared)\n", to.Name, from.Name, strings.Join(need, ", "))
	}
}

// noteTarget is the account note t says to carry its conversation on as:
// the one /ccam:switch named, or in the relay's hop the one the hook chose
// after a limit, which ccam does not choose again.
func noteTarget(t *handoff.Ticket, hop bool) string {
	switch {
	case t == nil:
		return ""
	case t.Target != "":
		return t.Target
	case hop:
		return t.Chosen
	}
	return ""
}

// confirmSwitch asks the person before conv, whose /ccam:switch note led
// to p, leaves the organization pool of the account it belongs to (apart).
// p is the account the note named or, picked, the one ccam chose when it
// named none, which cross_org may take out of the pool; cross_org does not
// lift the question. Claude Code expands /ccam:switch in the prompt of a
// job the model scheduled as if the person had typed it, so only a yes on
// the terminal, out of the model's reach, moves a conversation to another
// organization, between a seat and a personal plan, or to or from an API
// key account. A dry run asks nothing.
func confirmSwitch(e *env, cfg *config.Config, conv *conversation, p *config.Profile, picked bool) error {
	from := conv.owner
	if conv.ticket.Reason != handoff.Switch || e.dryRun {
		return nil
	}
	org := newOrgs(e, cfg, false)
	if !apart(e, org, from, p) {
		return nil
	}
	was, is := accountKind(e, org, from), accountKind(e, org, p)
	led, again := "/ccam:switch named "+p.Name, "asks"
	if picked {
		led, again = "the account picked for /ccam:switch is "+p.Name, "picks"
	}
	if launcher.ConfirmMove(fmt.Sprintf("ccam: %s, %s, but this conversation is on %s, %s. Move it to %s?", led, is, from.Name, was, p.Name)) {
		return nil
	}
	return fmt.Errorf("the conversation stays on %s, %s: it moves to %s, %s, only when you answer yes on a terminal. "+
		"`ccam next --continue` %s again, and `ccam %s --resume %s` continues it as %s", from.Name, was, p.Name, is, again, from.Name, transcript.ID(conv.path), from.Name)
}

// apart reports whether a conversation of from that moves to p leaves
// from's organization pool (config.SamePool, which a login whose plan ccam
// cannot tell fits). An API key account is in no pool: its key may belong
// to any organization, or to a gateway's, and ccam never moves a
// conversation there by itself (autoTarget, and creds.Has for next --auto).
func apart(e *env, org *orgs, from, p *config.Profile) bool {
	switch {
	case from.Name == p.Name:
		return false
	case apiKey(e, from) || apiKey(e, p):
		return true
	}
	return !config.SamePool(org.of(from), org.of(p))
}

// apiKey reports whether p uses an API key: ccam's own API account, or a
// Console login, whose key Claude Code keeps in .claude.json.
func apiKey(e *env, p *config.Profile) bool {
	return p.API != nil || creds.APIKey(profile.Dir(e.paths, p))
}

// accountKind names what p is, for the question before a move out of the
// pool.
func accountKind(e *env, org *orgs, p *config.Profile) string {
	id := org.of(p)
	switch {
	case apiKey(e, p):
		return "an API key account"
	case id.Pool() == config.Personal:
		return "a personal plan"
	case id.Pool() == "":
		return "an account whose plan ccam does not know"
	case id.OrgName != "":
		return "a " + id.PlanName() + " seat in " + id.OrgName
	}
	return "a " + id.PlanName() + " seat"
}

// recordLimit records the usage limit in the conversation of t, a limit
// note, as the hook does (usage.SaveLimit): the hook may have looked before
// Claude Code wrote the refused request out.
func recordLimit(e *env, t *handoff.Ticket) {
	if t.Limit == "org" {
		return // an organization's stop is no limit of the account (see hookStopFailure)
	}
	if l, ok := transcript.LastLimit(t.Transcript, t.At.Add(-time.Minute)); ok && usage.LimitWindow(l.Window) {
		_ = usage.SaveLimit(usageStateDir(e), t.Account, l.Window, l.ResetsAt, l.At)
	}
}

// candidates are the accounts to choose from; with logins false, LoggedIn
// is not looked up (macOS asks the Keychain for each account).
func candidates(e *env, cfg *config.Config, logins bool) []next.Candidate {
	out := make([]next.Candidate, 0, len(cfg.Profiles))
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		out = append(out, next.Candidate{
			Profile:  p,
			LastUsed: launcher.LastUsed(e.paths, p.Name),
			LoggedIn: logins && creds.Has(profile.Dir(e.paths, p), e.paths.Home),
		})
	}
	return out
}
