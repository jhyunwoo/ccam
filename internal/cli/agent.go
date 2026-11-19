package cli

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/transcript"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// agentEnv marks the environment of a claude that ccam agent runs. It goes
// on to everything that claude starts, so an agent cannot start agents of
// its own: it loads the same plugin and skill, and could otherwise spread
// over every account.
const agentEnv = "CCAM_AGENT"

// agentKeep is how long a record of an agent's conversation is kept.
const agentKeep = 30 * 24 * time.Hour

// agentRecord is what ccam keeps of each conversation an agent ran, in
// state/agent/<session id>.json: ccam agent --resume continues it from
// there, and ccam next --continue leaves it out.
type agentRecord struct {
	Account    string    `json:"account"`
	Cwd        string    `json:"cwd"`
	Transcript string    `json:"transcript"`
	Caller     string    `json:"caller,omitempty"` // the account that ran ccam agent
	At         time.Time `json:"at"`
}

// agentOutput is what --json prints.
type agentOutput struct {
	Account   string   `json:"account"`
	SessionID string   `json:"session_id"`
	Result    string   `json:"result"`
	IsError   bool     `json:"is_error"`
	Accounts  []string `json:"accounts"` // every account the task ran as, in order
}

func agentCmd(e *env) *cobra.Command {
	var asJSON bool
	var model, resume string
	var choice choiceFlags
	cmd := &cobra.Command{
		Use:   "agent [--model <model>] [--resume <id>] [--json] [--] <task>",
		Short: "Run a task as another account, so it uses that account's usage",
		Long: "Runs the task with claude -p in this folder as another logged-in account and prints\n" +
			"Claude's last reply. Claude Code calls it like a subagent that uses another account's\n" +
			"usage, so one task can go on longer than the current account's limits allow.\n\n" +
			"It never runs as the account it is called from (CCAM_PROFILE in a session ccam started,\n" +
			"else the account of this folder), nor as another folder of that Claude account, and it\n" +
			"stays among that account's Team or Enterprise organization, or personal plans\n" +
			"(--cross-org also takes the others). It picks by usage as `ccam pick` does, and agents\n" +
			"started together get different accounts. The agent works in the permission mode and\n" +
			"effort of the session that calls it, never more: run from a terminal, Claude Code's\n" +
			"default mode applies. When the agent's account reaches its usage limit, its\n" +
			"conversation goes on as another account. With no account to run as, ccam exits with\n" +
			"code 75; --wait waits for one to reset.\n\n" +
			"The task comes from the arguments, or from standard input when there are none.\n" +
			"stderr names the agent's session; --resume <id> continues that conversation with a\n" +
			"follow-up task. An agent cannot start agents of its own.\n\n" +
			"Not to be confused with `ccam agents`, which lists sessions, and `ccam attach`.",
		Example: "  ccam agent \"find why TestParse fails and fix it; reply with what changed\"\n" +
			"  ccam agent --model opus -- \"-v breaks the parser: find out why\"\n" +
			"  ccam agent --resume 7f3c9a2e-… \"now add a test for it\"\n" +
			"  ccam agent --json < task.md",
		Args: func(cmd *cobra.Command, args []string) error {
			return spacedValue(cmd, args, "wait", DefaultWait)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if os.Getenv(agentEnv) != "" {
				return errors.New("ccam agent runs inside an agent here, which cannot start agents of its own; do the task yourself")
			}
			task := strings.TrimSpace(strings.Join(args, " "))
			if task == "" && !isTerminal(e.stdin) {
				in, err := readInput(e.stdin, e.stderr, stdinWait)
				if err != nil {
					return err
				}
				task = strings.TrimSpace(string(in))
			}
			if task == "" {
				return errors.New("give the task as arguments or on standard input")
			}
			cfg, err := e.load()
			if err != nil {
				return err
			}
			p, err := choice.policy(cfg, nil)
			if err != nil {
				return err
			}
			p.Model = model
			wait, err := waitFor(choice.wait)
			if err != nil {
				return err
			}
			a := &agentRun{e: e, cfg: cfg, policy: p, wait: wait, org: newOrgs(e, cfg, choice.crossOrg), json: asJSON}
			a.caller = callerProfile(cfg)
			base := []string{"-p", "--output-format", "json"}
			if model != "" {
				base = append(base, "--model", model)
			}
			// The calling session's permission mode and effort, which a -p
			// run needs to use tools at all: nobody is there to allow them.
			base, _ = launcher.WithState(base, callerMode(e, a.caller), os.Getenv("CLAUDE_EFFORT"))
			a.base = base
			if resume != "" {
				return a.resume(resume, task)
			}
			return a.start(task)
		},
	}
	cmd.Flags().StringVar(&model, "model", "", "the model the agent uses (its weekly limit counts too), e.g. opus")
	cmd.Flags().StringVar(&resume, "resume", "", "continue the conversation of an earlier agent (the session id it printed)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "print the account, session and reply as JSON")
	choice.add(cmd)
	return cmd
}

// agentRun is one ccam agent: who called it, how accounts are chosen, and
// the claude args every run of it shares.
type agentRun struct {
	e      *env
	cfg    *config.Config
	policy usage.Policy
	wait   time.Duration
	org    *orgs
	caller *config.Profile // nil when ccam agent knows of none
	base   []string        // claude args before the session's
	json   bool
}

// callerProfile is the account ccam agent is called from: the one of the
// session ccam started (CCAM_PROFILE), else the one this folder starts.
func callerProfile(cfg *config.Config) *config.Profile {
	if p, ok := cfg.Find(os.Getenv("CCAM_PROFILE")); ok && os.Getenv("CCAM_PROFILE") != "" {
		return p
	}
	p, _ := launcher.HereProfile(cfg)
	return p
}

// callerMode is the permission mode of the Claude Code session that runs
// ccam agent (CLAUDE_CODE_SESSION_ID), read from its conversation in the
// caller's folder: under projects/ by its id, since the shell may have
// left the folder the session started in. "" when there is none.
func callerMode(e *env, caller *config.Profile) string {
	if path := callerTranscript(e, caller, os.Getenv("CLAUDE_CODE_SESSION_ID")); path != "" {
		return transcript.PermissionMode(path)
	}
	return ""
}

func callerTranscript(e *env, caller *config.Profile, id string) string {
	if caller == nil || !sessionID(id) {
		return ""
	}
	matches, _ := filepath.Glob(filepath.Join(profile.Dir(e.paths, caller), "projects", "*", id+".jsonl"))
	if len(matches) == 0 {
		return ""
	}
	return matches[0]
}

// sessionID reports whether id may name a conversation file: 1-64 of
// [A-Za-z0-9_-], so it cannot leave the folder it is looked up in.
func sessionID(id string) bool {
	if id == "" || len(id) > 64 {
		return false
	}
	for _, c := range id {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '-') {
			return false
		}
	}
	return true
}

// newSessionID makes a random (version 4) UUID, the form Claude Code's
// --session-id takes.
func newSessionID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// candidates are the accounts an agent may run as: logged in, not the
// caller, and fitting the caller (not its Claude account, its pool).
func (a *agentRun) candidates() []*config.Profile {
	var out []*config.Profile
	for i := range a.cfg.Profiles {
		prof := &a.cfg.Profiles[i]
		if a.caller != nil && prof.Name == a.caller.Name {
			continue
		}
		if creds.Has(profile.Dir(a.e.paths, prof), a.e.paths.Home) && a.org.fits(a.caller, prof) {
			out = append(out, prof)
		}
	}
	return out
}

// none is the message when no account may run the agent.
func (a *agentRun) none(msg string) string {
	if a.caller == nil {
		return msg
	}
	return a.org.none(a.caller, msg+" other than "+a.caller.Name)
}

// start runs task in a new conversation as the account with the most room.
func (a *agentRun) start(task string) error {
	profs := a.candidates()
	if len(profs) == 0 {
		return noRoom(a.none("no logged-in account to run the agent as"))
	}
	var picked *usage.Report
	prof, err := untilRoom(a.e, a.policy, a.wait, a.none("no account with known usage and headroom to run the agent as"), func() (*config.Profile, []*usage.Report, error) {
		prof, r, reports, err := pickOnce(a.e, a.cfg, profs, a.policy)
		picked = r
		return prof, reports, err
	})
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	return a.run(prof, a.policy.Use(picked, clock()), cwd, "", task)
}

// resume continues the conversation of the agent session id with task: as
// the account that ran it when that account may and has room, else as
// another, which Claude Code copies the conversation into. It runs in the
// folder that conversation belongs to.
func (a *agentRun) resume(id, task string) error {
	rec, ok := loadAgentRecord(a.e, id)
	if !ok {
		return fmt.Errorf("no agent conversation %s; ccam agent names its conversation on stderr when it starts and ends", id)
	}
	if _, err := os.Stat(rec.Transcript); err != nil {
		return fmt.Errorf("the conversation of agent %s is gone (%s)", id, rec.Transcript)
	}
	if err := os.Chdir(rec.Cwd); err != nil {
		return fmt.Errorf("the folder of agent %s: %w", id, err)
	}
	var prof *config.Profile
	why := ""
	for _, p := range a.candidates() {
		if p.Name == rec.Account {
			if r := reportsOf(a.e, a.cfg, []*config.Profile{p}, usage.Options{}); !a.policy.Full(r[0]) {
				prof, why = p, a.policy.Use(r[0], clock())
			}
		}
	}
	if prof == nil {
		profs := a.candidates()
		if len(profs) == 0 {
			return noRoom(a.none("no logged-in account to run the agent as"))
		}
		var picked *usage.Report
		var err error
		prof, err = untilRoom(a.e, a.policy, a.wait, a.none("no account with known usage and headroom to run the agent as"), func() (*config.Profile, []*usage.Report, error) {
			prof, r, reports, err := pickOnce(a.e, a.cfg, profs, a.policy)
			picked = r
			return prof, reports, err
		})
		if err != nil {
			return err
		}
		why = a.policy.Use(picked, clock())
	}
	return a.run(prof, why, rec.Cwd, rec.Transcript, task)
}

// run runs claude -p as prof in cwd, in a new conversation, or in a copy of
// the conversation at from (--resume <file> --fork-session) when from is
// set. When the account reaches its usage limit, the conversation goes on
// the same way as another account, at most once per account.
func (a *agentRun) run(prof *config.Profile, why, cwd, from, task string) error {
	e := a.e
	tried := map[string]bool{}
	if a.caller != nil {
		tried[a.caller.Name] = true // never the caller, and only what fits it
	}
	var accounts []string
	left, waited := a.wait, false
	for {
		tried[prof.Name] = true
		accounts = append(accounts, prof.Name)
		id := newSessionID()
		args := append([]string{}, a.base...)
		args = append(args, "--session-id", id)
		if from != "" {
			args = append(args, "--resume", from, "--fork-session")
		}
		args = append(args, "--", task)
		pl, err := launcher.Ready(e.paths, a.cfg, prof.Name, args, e.dryRun)
		if err != nil {
			return err
		}
		if e.dryRun {
			if err := pl.DryRun(e.stdout); err != nil {
				return err
			}
			fmt.Fprintf(e.stderr, "ccam: the agent would run as %s; at its usage limit, as another account\n", prof.Name)
			return nil
		}
		pl.Env = append(pl.Env, agentEnv+"=1")
		pl.PassInterrupt = true
		path := filepath.Join(transcript.ProjectDir(profile.Dir(e.paths, prof), cwd), id+".jsonl")
		rec := agentRecord{Account: prof.Name, Cwd: cwd, Transcript: path, At: clock().UTC()}
		if a.caller != nil {
			rec.Caller = a.caller.Name
		}
		if err := saveAgentRecord(e, id, rec); err != nil {
			pl.Done()
			return err
		}
		if why != "" {
			fmt.Fprintf(e.stderr, "ccam: agent runs as %s (%s); session %s\n", prof.Name, why, id)
		} else {
			fmt.Fprintf(e.stderr, "ccam: agent runs as %s; session %s\n", prof.Name, id)
		}
		start := clock()
		var out bytes.Buffer
		errTail := &tail{max: 64 << 10}
		code, err := pl.Run(e.paths, bytes.NewReader(nil), &out, io.MultiWriter(e.stderr, errTail))
		pl.Done()
		if err != nil {
			return err
		}
		res, ok := agentResult(out.Bytes())
		failed := code != 0 || res.IsError
		// A run stopped by a signal was asked to end, not refused.
		if !failed || code > 128 || pl.Signaled() || !limited(e, a.cfg, prof, a.policy, out.Bytes(), errTail.buf) {
			return a.finish(prof, id, accounts, res, ok, out.Bytes(), code)
		}
		recordAgentLimit(e, prof, path, start)
		next, _ := targetIn(e, a.cfg, a.org, prof, tried, nil, a.policy)
		if next == nil && left > 0 {
			fmt.Fprintf(e.stderr, "ccam: %s reached its usage limit, and no other account the agent may use has room left\n", prof.Name)
			// The first look only finds out how long to wait (as for
			// --fallback): saved usage may still show room in an account
			// that just reached its limit.
			first, began := true, clock()
			again := map[string]bool{}
			if a.caller != nil {
				again[a.caller.Name] = true
			}
			next, err = untilRoom(e, a.policy, left, a.none("no account has room left to run the agent as"), func() (*config.Profile, []*usage.Report, error) {
				p, reports := targetIn(e, a.cfg, a.org, prof, again, nil, a.policy)
				if first {
					first = false
					return nil, reports, nil
				}
				return p, reports, nil
			})
			left -= clock().Sub(began)
			if err != nil {
				_ = a.finish(prof, id, accounts, res, ok, out.Bytes(), code)
				return err
			}
			tried, waited = again, true
		}
		if next == nil {
			_ = a.finish(prof, id, accounts, res, ok, out.Bytes(), code)
			fmt.Fprintf(e.stderr, "ccam: %s reached its usage limit, and no other account the agent may use has room left\n", prof.Name)
			return launcher.ExitCode(exitNoRoom)
		}
		if waited {
			fmt.Fprintf(e.stderr, "ccam: the agent goes on as %s\n", next.Name)
			waited = false
		} else {
			fmt.Fprintf(e.stderr, "ccam: %s reached its usage limit; the agent goes on as %s\n", prof.Name, next.Name)
		}
		prof, why, from, task = next, "", path, handoff.ContinuePrompt
	}
}

// claudeResult is the part of claude -p --output-format json's result that
// ccam agent reads.
type claudeResult struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	IsError bool   `json:"is_error"`
	Result  string `json:"result"`
}

// agentResult reads the result object of claude -p --output-format json:
// the last line of out that is one. ok is false when there is none.
func agentResult(out []byte) (res claudeResult, ok bool) {
	lines := bytes.Split(bytes.TrimSpace(out), []byte("\n"))
	for i := len(lines) - 1; i >= 0; i-- {
		var r claudeResult
		if json.Unmarshal(bytes.TrimSpace(lines[i]), &r) == nil && r.Type == "result" {
			return r, true
		}
	}
	return claudeResult{}, false
}

// finish prints what the agent answered: Claude's last reply on stdout
// (claude's own output when it gave no result), or the JSON of --json, and
// on stderr how to follow up.
func (a *agentRun) finish(prof *config.Profile, id string, accounts []string, res claudeResult, ok bool, out []byte, code int) error {
	e := a.e
	failed := code != 0 || res.IsError || !ok
	switch {
	case a.json:
		o := agentOutput{Account: prof.Name, SessionID: id, Result: res.Result, IsError: failed, Accounts: accounts}
		if !ok {
			o.Result = strings.TrimSpace(string(out))
		}
		data, _ := json.MarshalIndent(o, "", "  ")
		e.printf("%s\n", data)
	case ok && res.Result != "":
		e.printf("%s\n", strings.TrimRight(res.Result, "\n"))
	case ok:
		fmt.Fprintf(e.stderr, "ccam: the agent gave no reply (%s)\n", res.Subtype)
	default:
		_, _ = e.stdout.Write(out)
	}
	follow := fmt.Sprintf("ccam agent --resume %s \"…\"", id)
	if failed {
		fmt.Fprintf(e.stderr, "ccam: the agent stopped as %s; its conversation goes on with: %s\n", prof.Name, follow)
	} else {
		fmt.Fprintf(e.stderr, "ccam: agent finished as %s; follow up with: %s\n", prof.Name, follow)
	}
	return exitWith(code)
}

// recordAgentLimit records the usage limit the agent's conversation at path
// ran into since start, from the row Claude Code wrote for it, as the hook
// does for a session: a -p run's hook does nothing, and without the record
// the next choice could take this account again while its saved usage
// still shows room. An organization's stop is not recorded.
func recordAgentLimit(e *env, prof *config.Profile, path string, start time.Time) {
	lim, ok := transcript.LastLimit(path, start.Add(-time.Second))
	if !ok || handoff.Classify(lim.Text) == handoff.Org || !usage.LimitWindow(lim.Window) {
		return
	}
	_ = usage.SaveLimit(usageStateDir(e), prof.Name, lim.Window, lim.ResetsAt, lim.At)
}

func agentDir(e *env) string { return filepath.Join(e.paths.State, "agent") }

// saveAgentRecord writes the record of the agent conversation id, and
// drops records older than agentKeep.
func saveAgentRecord(e *env, id string, rec agentRecord) error {
	dir := agentDir(e)
	if entries, err := os.ReadDir(dir); err == nil {
		now := clock()
		for _, ent := range entries {
			if info, err := ent.Info(); err == nil && now.Sub(info.ModTime()) > agentKeep {
				_ = os.Remove(filepath.Join(dir, ent.Name()))
			}
		}
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, _ := json.MarshalIndent(rec, "", "  ")
	return fsutil.AtomicWrite(filepath.Join(dir, id+".json"), append(data, '\n'), 0o600)
}

// loadAgentRecord reads the record of the agent conversation id.
func loadAgentRecord(e *env, id string) (agentRecord, bool) {
	var rec agentRecord
	if !sessionID(id) {
		return rec, false
	}
	data, err := os.ReadFile(filepath.Join(agentDir(e), id+".json"))
	if err != nil || json.Unmarshal(data, &rec) != nil || rec.Transcript == "" || rec.Cwd == "" {
		return rec, false
	}
	return rec, true
}

// agentConversation reports whether the conversation at path is one an
// agent ran, which ccam next --continue leaves out: the person's
// conversation is the one to continue, not an agent's that went on later.
func agentConversation(e *env, path string) bool {
	_, err := os.Stat(filepath.Join(agentDir(e), transcript.ID(path)+".json"))
	return err == nil
}
