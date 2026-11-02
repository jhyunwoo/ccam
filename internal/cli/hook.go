package cli

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/lock"
	"github.com/jhyunwoo/ccam/internal/transcript"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// hookInput is the part of a Claude Code hook's JSON that ccam reads.
type hookInput struct {
	SessionID  string `json:"session_id"`
	Transcript string `json:"transcript_path"`
	Cwd        string `json:"cwd"`
	AgentID    string `json:"agent_id"` // set when a subagent's loop stopped
	// Effort is the session's effort level, which a resumed conversation
	// does not get back on its own (absent for models without one).
	Effort struct {
		Level string `json:"level"`
	} `json:"effort"`
	Error       string `json:"error"`
	LastMessage string `json:"last_assistant_message"`
	CommandArgs string `json:"command_args"`
}

// folder is where the session started. The hook's cwd follows a cd in
// Claude Code's shell; the conversation stays under the start folder.
func (in hookInput) folder() string { return cmp.Or(os.Getenv("CLAUDE_PROJECT_DIR"), in.Cwd) }

func hookCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "hook stop-failure|switch",
		Short: "Run by the ccam plugin's hooks inside Claude Code (reads the hook's JSON on stdin)",
		Long: "The ccam plugin for Claude Code runs these; you do not need to.\n\n" +
			"  stop-failure  a turn stopped at a usage limit: note the conversation, count the\n" +
			"                account as full until the limit resets, and show a terminal\n" +
			"                notification, so `ccam next --continue` can continue it\n" +
			"  switch        /ccam:switch [account]: note the conversation and end the turn\n\n" +
			"In a session that ccam did not start, stop-failure does nothing and switch says why it cannot move it.\n" +
			"A subagent that stops at a usage limit counts as its session's stop. stop-failure does nothing for an\n" +
			"unattended session (-p, background).\n\n" +
			"With auto_switch = true in config.toml (Linux, macOS), stop-failure also ends Claude Code when another\n" +
			"account like this one has room (a seat of the same Team or Enterprise organization, or another personal\n" +
			"plan), and the ccam process running alongside the session continues the conversation there.",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"stop-failure", "switch"},
		RunE: func(cmd *cobra.Command, args []string) error {
			var in hookInput
			_ = json.NewDecoder(io.LimitReader(e.stdin, 1<<20)).Decode(&in)
			account := os.Getenv("CCAM_PROFILE")
			switch args[0] {
			case "stop-failure":
				return hookStopFailure(e, in, account)
			case "switch":
				return hookSwitch(e, in, account)
			}
			return fmt.Errorf("unknown hook %q (stop-failure or switch)", args[0])
		},
	}
}

// hookOutput writes the hook's JSON answer.
func hookOutput(e *env, v map[string]any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	e.printf("%s\n", data)
	return nil
}

// limitNames name the limits of handoff.LimitKind for the notification.
var limitNames = map[string]string{"session": "5-hour", "weekly": "weekly", "usage": "usage", "model": "model", "spend": "spend", "org": "organization's usage"}

// limitPoll is how long stop-failure waits for the refused request to show
// in the conversation: Claude Code writes it out every 100 ms and does not
// wait for the hook, which may start first. Tests change it.
var limitPoll = 1500 * time.Millisecond

// hookStopFailure notes a session that stopped at a usage limit, and when
// Claude Code wrote the limit into the conversation, records the account as
// full until the limit resets (usage.SaveLimit), so that no automatic
// choice takes it again before then. With auto_switch on, it then ends
// Claude Code when another account has room (autoSwitch), and the relay
// carries the conversation on. A short server throttle or an API key's
// rate limit stops a turn the same way, but is not a usage limit, and is
// left alone; so are sessions nobody sits at (-p, background).
//
// A subagent's loop (agent_id) that stops at a limit counts as the
// session's: the main thread may sit waiting for it, and never make the
// request that would stop it too. Its limit is in its own conversation
// (transcript.SubagentPath). Subagents side by side stop one after another,
// or at once: one hook at a time handles them, and the person gets one
// notification for a limit.
func hookStopFailure(e *env, in hookInput, account string) error {
	if account == "" || in.Error != "rate_limit" || unattended() {
		return nil
	}
	start := time.Now()
	notePath := handoff.Path(e.paths.State, account, in.SessionID)
	read, path := transcript.LastLimit, in.Transcript
	var prev *handoff.Ticket // the note before a subagent's stop
	if in.AgentID != "" {
		l, err := lock.TryFlock(strings.TrimSuffix(notePath, ".json") + ".lock")
		if errors.Is(err, lock.ErrBusy) {
			return nil // another subagent's stop is being handled
		}
		if err == nil {
			defer l.Unlock()
		}
		if cur, ok := handoff.LoadFile(e.paths.State, notePath); ok {
			if moving(cur) {
				return nil // Claude Code is ending
			}
			prev = cur
		}
		read, path = transcript.SubagentLimit, transcript.SubagentPath(in.Transcript, in.AgentID)
	}
	var lim transcript.Limit
	ticket := func(kind string) handoff.Ticket {
		return handoff.Ticket{Account: account, SessionID: in.SessionID, Transcript: in.Transcript,
			Cwd: in.folder(), At: time.Now(), Reason: handoff.Limit, Limit: kind, Effort: cmp.Or(in.Effort.Level, os.Getenv("CLAUDE_EFFORT")),
			Window: lim.Window, ResetsAt: lim.ResetsAt, LimitAt: lim.At, Subagent: in.AgentID}
	}
	note := func(kind string) error {
		// The note of a move under way stays: the relay follows it.
		if cur, ok := handoff.LoadFile(e.paths.State, notePath); ok && moving(cur) {
			return nil
		}
		return handoff.Save(e.paths.State, ticket(kind))
	}
	// The note comes first when the message says enough: Claude Code may
	// close before the limit shows in the conversation, and end the hook.
	kind := handoff.LimitKind(in.LastMessage)
	if kind != "" {
		if err := note(kind); err != nil {
			return err
		}
	}
	lim, row := waitLimit(read, path, start, in.LastMessage)
	switch {
	case row:
		// Again, now with the window and when it resets.
		kind = cmp.Or(kind, "usage")
		if err := note(kind); err != nil {
			return err
		}
	case kind == "":
		return nil
	}
	// An organization's stop (out of usage, a group at $0, an allocation
	// turned off) is not recorded: its admins can lift it at any time.
	own := row && handoff.Classify(in.LastMessage) != handoff.Org
	recorded := false
	if own && usage.LimitWindow(lim.Window) {
		recorded = usage.SaveLimit(usageStateDir(e), account, lim.Window, lim.ResetsAt, lim.At) == nil
	}
	name, reset := limitNames[kind], handoff.Reset(in.LastMessage)
	if n := handoff.WindowLabel(lim.Window); n != "" && own {
		name = n
		if !lim.ResetsAt.IsZero() {
			reset = handoff.FormatReset(lim.ResetsAt, time.Now())
		}
	}
	hit := handoff.HitLimit(account, name, reset)
	msg := hit + ". `ccam next --continue` continues the conversation on another account."
	if none := noTakeover(e, account); none != "" {
		// That command would refuse, and auto_switch has nowhere to go.
		msg = hit + ". " + none
	} else if recorded && (in.AgentID == "" || subagentHolds(lim.Window, transcript.Model(in.Transcript))) {
		// auto_switch leaves only an account recorded as full: no session
		// comes back to it before the reset.
		moved, say := autoSwitch(e, in, ticket(kind), lim, start, hit)
		if moved {
			return nil // Claude Code is ending: nothing written now would reach the screen
		}
		msg = cmp.Or(say, msg)
	}
	if prev != nil && row && prev.Reason == handoff.Limit && prev.Window == lim.Window && prev.ResetsAt.Equal(lim.ResetsAt) {
		return nil // an earlier stop at this limit said it
	}
	// Claude Code shows nothing else a StopFailure hook prints.
	return hookOutput(e, map[string]any{"terminalSequence": handoff.Notification(msg)})
}

// noTakeover ends the notification when no other account may take over from
// account by the rule of `ccam next --continue` (takeover: its pool, or with
// cross_org any other Claude account; logged in or not), which then refuses.
// It is "" when one may.
func noTakeover(e *env, account string) string {
	cfg, err := e.load()
	if err != nil || slices.ContainsFunc(takeover(e, cfg, account), func(n string) bool { return n != account }) {
		return ""
	}
	var id *config.Identity
	if p, ok := cfg.Find(account); ok {
		id = newOrgs(e, cfg, false).of(p)
	}
	return "You have no other " + poolName(id) + ", so ccam cannot move this conversation."
}

// unattended reports whether nobody sits at the session: Claude Code sets
// CLAUDE_CODE_SESSION_ATTENDED=0 for -p, background and daemon sessions,
// and an "sdk-" entrypoint for -p ones (also before it had that variable).
func unattended() bool {
	return os.Getenv("CLAUDE_CODE_SESSION_ATTENDED") == "0" || strings.HasPrefix(os.Getenv("CLAUDE_CODE_ENTRYPOINT"), "sdk-")
}

// waitLimit reads with read (transcript.LastLimit, or SubagentLimit) the
// usage limit Claude Code wrote into the conversation at path for the
// request it refused with msg, giving it limitPoll to show up. A limit with
// another message is an earlier one, perhaps another account's in a shared
// history, not this stop's.
func waitLimit(read func(string, time.Time) (transcript.Limit, bool), path string, start time.Time, msg string) (transcript.Limit, bool) {
	if path == "" {
		return transcript.Limit{}, false
	}
	for {
		// The row is made before the hook starts and written soon after.
		l, ok := read(path, start.Add(-time.Minute))
		if ok && (msg == "" || strings.TrimSpace(l.Text) == strings.TrimSpace(msg)) {
			return l, true
		}
		if time.Since(start) >= limitPoll {
			return transcript.Limit{}, false
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// endClaude sends SIGTERM to the Claude Code process pid after a moment,
// once Claude Code has read the hook's answer. A variable so tests can run
// the hook under a relay.
var endClaude = func(pid int) error {
	return exec.Command("sh", "-c", "sleep 0.3; kill -TERM "+strconv.Itoa(pid)).Start()
}

// hookSwitch notes the session for /ccam:switch and ends the turn before
// any request is made. Under a relay of this ccam it also ends Claude Code,
// and the relay carries the conversation on as the other account; for an
// account outside the session's organization pool, it waits for the
// person's /exit and asks first.
func hookSwitch(e *env, in hookInput, account string) error {
	stop := func(msg string) error {
		return hookOutput(e, map[string]any{"continue": false, "stopReason": msg})
	}
	if account == "" {
		return stop("ccam: this session was not started by ccam, so ccam cannot move it. Start Claude Code with `ccam <account>`.")
	}
	target := ""
	if f := strings.Fields(in.CommandArgs); len(f) > 0 {
		target = f[0]
	}
	cfg, err := e.load()
	if err != nil {
		return stop("ccam: " + err.Error())
	}
	var to *config.Profile
	if target != "" {
		p, ok := cfg.Find(target)
		if !ok {
			return stop(fmt.Sprintf("ccam: no account %q (see `ccam list`).", target))
		}
		if target = p.Name; target == account {
			return stop("ccam: this session already runs as " + account + ".")
		}
		to = p
	}
	if !fsutil.Exists(in.Transcript) {
		// Nothing written yet: there is no conversation to carry on.
		return stop("ccam: this session has no conversation yet, so there is nothing to continue. To work as another account, type /exit and run `ccam " +
			cmp.Or(target, "next") + "`.")
	}
	t := handoff.Ticket{Account: account, SessionID: in.SessionID, Transcript: in.Transcript, Cwd: in.folder(),
		At: time.Now(), Reason: handoff.Switch, Target: target, Effort: cmp.Or(in.Effort.Level, os.Getenv("CLAUDE_EFFORT"))}
	exit := "ccam: type /exit, then run `ccam next --continue` to continue this conversation on " + cmp.Or(target, "the next account") + "."
	// Without a target the relay goes on with `ccam next --continue --auto`.
	who := cmp.Or(target, "the account with the most room")
	carry := "ccam: type /exit; ccam then continues this conversation on " + who + "."
	// A move out of the pool waits for the person's yes on the terminal
	// (confirmSwitch), and this hook cannot tell a /ccam:switch the person
	// typed from one that a job the model scheduled fired. So Claude Code is
	// not ended for it (a recurring job would end the session at every run):
	// the relay asks once the person types /exit.
	away := false
	if from, ok := cfg.Find(account); ok && to != nil {
		org := newOrgs(e, cfg, false)
		if away = apart(e, org, from, to); away {
			carry = fmt.Sprintf("ccam: type /exit; ccam then asks whether to continue this conversation on %s, %s (%s is %s).",
				to.Name, accountKind(e, org, to), account, accountKind(e, org, from))
		}
	}
	save := func(msg string) error {
		if err := handoff.Save(e.paths.State, t); err != nil {
			return stop("ccam: " + err.Error())
		}
		return stop(msg)
	}
	// Only a relay of this ccam carries the note on; an older one, started
	// before an update, would not, so its session is left to /exit. Nor
	// is a session nobody sits at (-p, background, a teammate) the one it
	// runs.
	relay, ok := launcher.RelayPID(os.Getenv("CCAM_RELAY"))
	if !ok || unattended() {
		return save(exit)
	}
	claude, _ := strconv.Atoi(os.Getenv("CLAUDE_PID"))
	if claude <= 0 {
		// Claude Code before 2.1.214 has no CLAUDE_PID, so it cannot be
		// ended from here: the relay carries on after /exit.
		t.Relay = os.Getenv("CCAM_RELAY")
		return save(carry)
	}
	child, ok := launcher.RelayChild(claude, relay)
	if !ok {
		// A claude below the relayed one (started in its shell, say) has
		// CCAM_RELAY too, but is not what the relay runs.
		return save(exit)
	}
	t.Relay, t.Child = os.Getenv("CCAM_RELAY"), child
	if runtime.GOOS == "windows" || away {
		// Claude Code cannot be ended from here on Windows, nor is it for a
		// move out of the pool: the relay carries on after /exit.
		return save(carry)
	}
	t.Stopped = true
	if err := handoff.Save(e.paths.State, t); err != nil {
		return stop("ccam: " + err.Error())
	}
	if err := endClaude(claude); err != nil {
		t.Stopped = false
		return save(carry)
	}
	return stop("ccam: moving this conversation to " + who + "…")
}
