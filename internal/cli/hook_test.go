package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/lock"
	"github.com/jhyunwoo/ccam/internal/transcript"
	"github.com/jhyunwoo/ccam/internal/usage"
)

func hookJSON(t *testing.T, v map[string]string) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// A usage limit leaves a note and a notification; next --continue then
// carries that very conversation on as another account, never the one
// that hit the limit.
func TestLimitHookAndNextContinue(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 40, "side": 30})
	cfg, _ := config.Load(c.p.ConfigFile())
	cfg.ClaudeBin = stubClaude
	cfg.Save(c.p.ConfigFile())
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	os.MkdirAll(filepath.Dir(conv), 0o700)
	os.WriteFile(conv, []byte("{}\n"), 0o600)

	hookEnv(t)
	t.Setenv("CCAM_PROFILE", "personal")
	// A short throttle is not a usage limit.
	throttle := hookJSON(t, map[string]string{"error": "rate_limit", "last_assistant_message": "API Error: Server is temporarily limiting requests (not your usage limit)", "transcript_path": conv, "cwd": cwd})
	if out, _, code := c.run(throttle, "hook", "stop-failure"); code != 0 || out != "" {
		t.Fatalf("%d %q", code, out)
	}
	limit := hookJSON(t, map[string]string{"error": "rate_limit", "last_assistant_message": "You've hit your session limit · resets 3:45pm",
		"session_id": "s1", "transcript_path": conv, "cwd": cwd})
	out, _, code := c.run(limit, "hook", "stop-failure")
	var got map[string]string
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil || !strings.HasPrefix(got["terminalSequence"], "\x1b]9;ccam: personal hit its 5-hour limit (resets 3:45pm).") {
		t.Fatalf("%d %q", code, out)
	}
	// personal has the most room on paper, but just hit its limit.
	_, errOut, code := c.run("", "--dry-run", "next", "--auto", "--continue")
	if code != 0 || !strings.Contains(errOut, "→ personal → side") || !strings.Contains(errOut, "copies it into side") {
		t.Fatalf("%d %s", code, errOut)
	}
	// A dry run leaves the note for the real run.
	if _, ok := handoff.ForFolder(c.p.State, cwd, time.Now()); !ok {
		t.Fatal("the note is gone after a dry run")
	}
	// Without the limit in the conversation, the message alone says it.
	if _, err := os.Stat(filepath.Join(c.p.State, "usage", "personal.limit.json")); err == nil {
		t.Fatal("a limit record without the limit row")
	}
}

// stopReason is the reason of a hook's answer that ends the turn.
func stopReason(t *testing.T, out string) string {
	t.Helper()
	var got struct {
		Continue   *bool  `json:"continue"`
		StopReason string `json:"stopReason"`
	}
	if json.Unmarshal([]byte(out), &got) != nil || got.Continue == nil || *got.Continue {
		t.Fatalf("not stopped: %q", out)
	}
	return got.StopReason
}

func TestSwitchHook(t *testing.T) {
	hookEnv(t)
	c := newTestCLI(t, nil)
	cwd, _ := os.Getwd()
	conv := filepath.Join(t.TempDir(), "s2.jsonl")
	os.WriteFile(conv, []byte("{}\n"), 0o600)
	in := func(args string) string {
		return hookJSON(t, map[string]string{"session_id": "s2", "transcript_path": conv, "cwd": cwd, "command_args": args})
	}
	stopped := func(out string) string { return stopReason(t, out) }
	t.Setenv("CCAM_PROFILE", "")
	if out, _, _ := c.run(in(""), "hook", "switch"); !strings.Contains(stopped(out), "not started by ccam") {
		t.Fatal(out)
	}
	t.Setenv("CCAM_PROFILE", "work")
	if out, _, _ := c.run(in("nobody"), "hook", "switch"); !strings.Contains(stopped(out), `no account "nobody"`) {
		t.Fatal(out)
	}
	// Nothing written yet: no conversation to note or to carry on.
	empty := hookJSON(t, map[string]string{"session_id": "s0", "transcript_path": filepath.Join(t.TempDir(), "s0.jsonl"), "cwd": cwd, "command_args": "side"})
	if out, _, _ := c.run(empty, "hook", "switch"); !strings.Contains(stopped(out), "no conversation yet") || !strings.Contains(out, "`ccam side`") {
		t.Fatal(out)
	}
	if entries, _ := os.ReadDir(handoff.Dir(c.p.State, "work")); len(entries) != 0 {
		t.Fatalf("a note for no conversation: %v", entries)
	}
	out, _, _ := c.run(in("3"), "hook", "switch")
	if r := stopped(out); !strings.Contains(r, "/exit") || !strings.Contains(r, "on side") {
		t.Fatal(r)
	}
	tk, ok := handoff.ForFolder(c.p.State, cwd, time.Now())
	if !ok || tk.Account != "work" || tk.Target != "side" || tk.Reason != handoff.Switch || tk.Relay != "" || tk.Stopped {
		t.Fatalf("%+v", tk)
	}
	// The note belongs to the folder the session started in
	// (CLAUDE_PROJECT_DIR), whatever cd Claude Code's shell made since.
	start := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", start)
	moved := hookJSON(t, map[string]string{"session_id": "s3", "transcript_path": conv, "cwd": filepath.Join(start, "sub"), "command_args": "side"})
	stopped(func() string { out, _, _ := c.run(moved, "hook", "switch"); return out }())
	if tk, ok := handoff.ForFolder(c.p.State, start, time.Now()); !ok || tk.SessionID != "s3" || tk.Cwd != start {
		t.Fatalf("%+v", tk)
	}
}

// Under a relay of this ccam, /ccam:switch ends Claude Code once its answer
// has been read (on Windows it is not ended, and the relay carries on after
// /exit), and the note tells the relay which claude that was. An older
// relay (whose CCAM_RELAY is a bare pid) would not carry the note on, and a
// claude that is not the relay's own is no session it runs: both are told
// to /exit, and nothing is ended.
func TestSwitchHookUnderARelay(t *testing.T) {
	windows := runtime.GOOS == "windows"
	hookEnv(t)
	c := newTestCLI(t, nil)
	cwd, _ := os.Getwd()
	conv := filepath.Join(t.TempDir(), "s2.jsonl")
	os.WriteFile(conv, []byte("{}\n"), 0o600)
	var ended []int
	old := endClaude
	endClaude = func(pid int) error { ended = append(ended, pid); return nil }
	t.Cleanup(func() { endClaude = old })
	// This test is the relay; its child stands in for Claude Code.
	claude := exec.Command("sleep", "5")
	if windows {
		claude = exec.Command("ping", "-n", "6", "127.0.0.1") // 5 s
	}
	if err := claude.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { claude.Process.Kill(); claude.Wait() }()
	pid := claude.Process.Pid
	t.Setenv("CCAM_PROFILE", "work")
	t.Setenv("CLAUDE_PID", strconv.Itoa(pid))
	t.Setenv("CLAUDE_EFFORT", "high")
	run := func(relay string) (string, *handoff.Ticket) {
		t.Helper()
		ended = nil
		t.Setenv("CCAM_RELAY", relay)
		out, _, _ := c.run(hookJSON(t, map[string]string{"session_id": "s2", "transcript_path": conv, "cwd": cwd, "command_args": "side"}), "hook", "switch")
		tk, _ := handoff.ForFolder(c.p.State, cwd, time.Now())
		return stopReason(t, out), tk
	}
	me := "2:" + strconv.Itoa(os.Getpid())
	says, ends, most := "moving this conversation to side", []int{pid}, "moving this conversation to the account with the most room"
	if windows {
		// Nothing is ended there: the relay carries on after /exit.
		says, ends, most = "type /exit; ccam then continues this conversation on side.", nil, "type /exit; ccam then continues this conversation on the account with the most room."
	}
	if r, tk := run(me); !strings.Contains(r, says) || !reflect.DeepEqual(ended, ends) ||
		tk.Relay != me || tk.Child != pid || tk.Stopped == windows || tk.Effort != "high" {
		t.Fatalf("%q %v %+v", r, ended, tk)
	}
	for _, relay := range []string{strconv.Itoa(os.Getpid()), "2:" + strconv.Itoa(os.Getppid())} {
		if r, tk := run(relay); !strings.Contains(r, "/exit, then run `ccam next --continue`") || len(ended) != 0 || tk.Relay != "" || tk.Stopped {
			t.Fatalf("CCAM_RELAY=%s: %q %v %+v", relay, r, ended, tk)
		}
	}
	// A session nobody sits at (-p, background, a teammate) is not the one
	// the relay runs, whatever its process tree.
	t.Setenv("CLAUDE_CODE_SESSION_ATTENDED", "0")
	if r, tk := run(me); !strings.Contains(r, "/exit, then run `ccam next --continue`") || len(ended) != 0 || tk.Relay != "" {
		t.Fatalf("unattended: %q %v %+v", r, ended, tk)
	}
	t.Setenv("CLAUDE_CODE_SESSION_ATTENDED", "")
	// Claude Code before 2.1.214 has no CLAUDE_PID: the relay carries on after /exit.
	t.Setenv("CLAUDE_PID", "")
	if r, tk := run(me); !strings.Contains(r, "type /exit; ccam then continues") || len(ended) != 0 || tk.Relay != me || tk.Child != 0 || tk.Stopped {
		t.Fatalf("%q %v %+v", r, ended, tk)
	}
	// With no account named, the relay takes the one with the most room.
	t.Setenv("CLAUDE_PID", strconv.Itoa(pid))
	out, _, _ := c.run(hookJSON(t, map[string]string{"session_id": "s2", "transcript_path": conv, "cwd": cwd}), "hook", "switch")
	if r := stopReason(t, out); !strings.Contains(r, most) {
		t.Fatal(r)
	}
}

// hookEnv makes hook tests independent of the session that runs them: an
// attended interactive session without a relay or Remote Control, a
// terminal with OSC 9, and no wait for a limit row that is not there.
func hookEnv(t *testing.T) {
	t.Setenv("CLAUDE_CODE_SESSION_ATTENDED", "")
	t.Setenv("CLAUDE_CODE_ENTRYPOINT", "")
	t.Setenv("CLAUDE_CODE_BRIDGE_SESSION_ID", "")
	t.Setenv("CLAUDE_PROJECT_DIR", "")
	t.Setenv("CLAUDE_PID", "")
	t.Setenv("CLAUDE_EFFORT", "")
	t.Setenv("CCAM_RELAY", "")
	t.Setenv("KITTY_WINDOW_ID", "")
	t.Setenv("TERM_PROGRAM", "iTerm.app")
	t.Setenv("TERM", "xterm-256color")
	old := limitPoll
	limitPoll = 0
	t.Cleanup(func() { limitPoll = old })
}

// limitConversation writes the 2.1.306 fixture conversation, whose last row
// is Claude Code's refusal at the five_hour limit, refused at at and
// resetting at reset; edits are old, new pairs applied after that.
func limitConversation(t *testing.T, path string, at, reset time.Time, edits ...string) {
	t.Helper()
	data, err := os.ReadFile("../transcript/testdata/limit-2.1.306-five_hour.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	s := strings.NewReplacer("2026-10-21T14:51:46.515Z", at.UTC().Format(time.RFC3339Nano), "1792607400", strconv.FormatInt(reset.Unix(), 10)).Replace(string(data))
	if strings.Count(s, strconv.FormatInt(reset.Unix(), 10)) != 2 {
		t.Fatal("the fixture changed; update limitConversation")
	}
	for i := 0; i+1 < len(edits); i += 2 {
		s = strings.ReplaceAll(s, edits[i], edits[i+1])
	}
	os.MkdirAll(filepath.Dir(path), 0o700)
	if err := os.WriteFile(path, []byte(s), 0o600); err != nil {
		t.Fatal(err)
	}
}

// With the refused request in the conversation, the hook records the
// account as full until Claude Code's reset time and says when that is.
func TestLimitHookRecordsTheLimit(t *testing.T) {
	hookEnv(t)
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 40})
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	now := time.Now().Truncate(time.Second)
	reset := now.Add(3*time.Hour + 20*time.Minute)
	// The message the fixture's row has, as Claude Code passes it to the hook.
	msg := "You've hit your session limit · resets 3:30am (Asia/Seoul)"
	limitConversation(t, conv, now.Add(-time.Second), reset)
	t.Setenv("CCAM_PROFILE", "personal")
	in := func(msg string, extra ...string) string {
		v := map[string]string{"error": "rate_limit", "last_assistant_message": msg, "session_id": "s1", "transcript_path": conv, "cwd": cwd}
		for i := 0; i+1 < len(extra); i += 2 {
			v[extra[i]] = extra[i+1]
		}
		return hookJSON(t, v)
	}
	record := filepath.Join(c.p.State, "usage", "personal.limit.json")
	nothing := func(what string, out string) {
		t.Helper()
		if _, err := os.Stat(record); err == nil || out != "" {
			t.Fatalf("%s: out %q, record %v", what, out, err)
		}
		if _, ok := handoff.ForFolder(c.p.State, cwd, time.Now()); ok {
			t.Fatalf("%s: a note", what)
		}
	}
	// Nobody sits at a -p or background session.
	t.Setenv("CLAUDE_CODE_SESSION_ATTENDED", "0")
	out, _, _ := c.run(in(msg), "hook", "stop-failure")
	nothing("unattended", out)
	t.Setenv("CLAUDE_CODE_SESSION_ATTENDED", "")
	t.Setenv("CLAUDE_CODE_ENTRYPOINT", "sdk-cli")
	out, _, _ = c.run(in(msg), "hook", "stop-failure")
	nothing("-p entrypoint", out)
	t.Setenv("CLAUDE_CODE_ENTRYPOINT", "cli")

	out, _, code := c.run(in(msg, "effort", ""), "hook", "stop-failure")
	var got map[string]string
	want := "\x1b]9;ccam: personal hit its 5-hour limit (resets " + handoff.FormatReset(reset, now) + "). `ccam next --continue`"
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil || !strings.HasPrefix(got["terminalSequence"], want) {
		t.Fatalf("%d %q, want %q", code, out, want)
	}
	if _, ok := handoff.ForFolder(c.p.State, cwd, time.Now()); !ok {
		t.Fatal("no note")
	}
	out, _, _ = c.run("", "usage", "--json")
	var reports []usage.Report
	if err := json.Unmarshal([]byte(out), &reports); err != nil || reports[0].Limits.FiveHour.Percent != 100 || !reports[0].Limits.FiveHour.ResetsAt.Equal(reset) {
		t.Fatalf("%v %s", err, out)
	}
	// A spend limit stops the account at a full window too.
	os.Remove(record)
	spend := "You've hit your org's monthly spend limit · ask your admin to raise it · your session limit resets 3:30am"
	limitConversation(t, conv, now.Add(-time.Second), reset, msg, spend)
	if out, _, _ = c.run(in(spend), "hook", "stop-failure"); !fileThere(record) || !strings.Contains(out, "personal hit its 5-hour limit") {
		t.Fatalf("spend limit: %q", out)
	}
	// An organization's stop is lifted by its admins, not at the reset.
	os.Remove(record)
	stop := "Your org is out of usage · contact your admin"
	limitConversation(t, conv, now.Add(-time.Second), reset, msg, stop)
	if out, _, _ = c.run(in(stop), "hook", "stop-failure"); fileThere(record) || !strings.Contains(out, "personal hit its organization's usage limit") {
		t.Fatalf("org: %q", out)
	}
	// The Fable allowance and usage credits are no windows ccam chooses by.
	bare := "You've hit your limit · resets 3pm"
	limitConversation(t, conv, now.Add(-time.Second), reset, msg, bare, `"rateLimitType":"five_hour"`, `"rateLimitType":"overage"`)
	if out, _, _ = c.run(in(bare), "hook", "stop-failure"); fileThere(record) || !strings.Contains(out, "personal hit its usage limit (resets 3pm)") {
		t.Fatalf("overage: %q", out)
	}
}

// subagentConversation writes conversation conv of session s1 without a
// limit in it, and below it that of its subagent id, whose last row is
// Claude Code's refusal at the five_hour limit (as limitConversation, with
// edits); it returns the subagent's.
func subagentConversation(t *testing.T, conv, id string, at, reset time.Time, edits ...string) string {
	t.Helper()
	path := filepath.Join(strings.TrimSuffix(conv, ".jsonl"), "subagents", "agent-"+id+".jsonl")
	limitConversation(t, path, at, reset, append([]string{`"isSidechain":false`, `"isSidechain":true`}, edits...)...)
	limitConversation(t, conv, at, reset)
	data, _ := os.ReadFile(conv)
	var keep []string
	for _, l := range strings.SplitAfter(string(data), "\n") {
		if !strings.Contains(l, "usage_limit_reached") {
			keep = append(keep, l)
		}
	}
	if err := os.WriteFile(conv, []byte(strings.Join(keep, "")), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// A subagent that stops at a usage limit stops its session's work too: its
// own conversation has the limit, which the hook records and notes as the
// main thread's. Subagents side by side meet the limit one after another,
// or at once; the session gets one notification.
func TestLimitHookSubagent(t *testing.T) {
	hookEnv(t)
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 40})
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	now := time.Now().Truncate(time.Second)
	reset := now.Add(3*time.Hour + 20*time.Minute)
	msg := "You've hit your session limit · resets 3:30am (Asia/Seoul)" // the fixture's
	subagentConversation(t, conv, "a2", now.Add(-time.Second), reset)
	subagentConversation(t, conv, "a1", now.Add(-time.Second), reset)
	t.Setenv("CCAM_PROFILE", "personal")
	in := func(agent string) string {
		return `{"agent_id":"` + agent + `","agent_type":"general-purpose",` + hookJSON(t, map[string]string{"error": "rate_limit",
			"last_assistant_message": msg, "session_id": "s1", "transcript_path": conv, "cwd": cwd})[1:]
	}
	record := filepath.Join(c.p.State, "usage", "personal.limit.json")
	notePath := handoff.Path(c.p.State, "personal", "s1")

	out, _, code := c.run(in("a1"), "hook", "stop-failure")
	var got map[string]string
	want := "\x1b]9;ccam: personal hit its 5-hour limit (resets " + handoff.FormatReset(reset, now) + "). `ccam next --continue`"
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil || !strings.HasPrefix(got["terminalSequence"], want) {
		t.Fatalf("%d %q, want %q", code, out, want)
	}
	n, ok := handoff.LoadFile(c.p.State, notePath)
	if !ok || n.Subagent != "a1" || n.Transcript != conv || n.Window != "five_hour" || !n.ResetsAt.Equal(reset) || !fileThere(record) {
		t.Fatalf("note %+v, record there: %v", n, fileThere(record))
	}
	// The next subagent at that limit: the person knows already.
	if out, _, _ = c.run(in("a2"), "hook", "stop-failure"); out != "" {
		t.Fatalf("a second notification: %q", out)
	}
	// Another hook is at a subagent's limit of this session right now.
	os.Remove(notePath)
	l, err := lock.TryFlock(strings.TrimSuffix(notePath, ".json") + ".lock")
	if err != nil {
		t.Fatal(err)
	}
	out, _, _ = c.run(in("a1"), "hook", "stop-failure")
	l.Unlock()
	if _, ok := handoff.LoadFile(c.p.State, notePath); out != "" || ok {
		t.Fatalf("while another hook runs: %q, note %v", out, ok)
	}
	// A subagent's stop does not read the main thread's limit: the message
	// alone says it.
	os.Remove(record)
	limitConversation(t, conv, now.Add(-time.Second), reset)
	if out, _, _ = c.run(in("a3"), "hook", "stop-failure"); !strings.Contains(out, "personal hit its 5-hour limit (resets 3:30am)") || fileThere(record) {
		t.Fatalf("no conversation of the subagent: %q, record there: %v", out, fileThere(record))
	}
}

// `ccam hook --help` says what stop-failure does with auto_switch on.
func TestHookHelpSaysWhatAutoSwitchDoes(t *testing.T) {
	c := newTestCLI(t, nil)
	out, _, _ := c.run("", "hook", "--help")
	if !strings.Contains(strings.Join(strings.Fields(out), " "), "With auto_switch = true in config.toml (Linux, macOS), stop-failure also ends Claude Code when another account like this one has room") {
		t.Fatal(out)
	}
}

// When no other account may take over from the one at its limit by the
// rule of `ccam next --continue` (its pool, or any account with cross_org;
// logged in or not), that command would refuse: the notification says so
// instead of sending the person there.
func TestLimitHookSaysWhenNoOtherAccountFits(t *testing.T) {
	hookEnv(t)
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 40}) // side is not logged in
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	reset := time.Now().Add(3 * time.Hour)
	limitConversation(t, conv, time.Now().Add(-time.Second), reset)
	t.Setenv("CCAM_PROFILE", "personal")
	hit := "ccam: personal hit its 5-hour limit (resets " + handoff.FormatReset(reset, time.Now()) + ")"
	for _, k := range []struct {
		name               string
		personal, work, sd [2]string // organization, plan
		crossOrg           bool
		want               string
	}{
		{"a personal plan among seats", [2]string{"p1", "claude_max"}, [2]string{"acme", "claude_team"}, [2]string{"acme", "claude_team"}, false,
			hit + ". You have no other personal plans, so ccam cannot move this conversation."},
		{"a seat alone in its organization", [2]string{"acme", "claude_team"}, [2]string{"other", "claude_team"}, [2]string{"s1", "claude_pro"}, false,
			hit + ". You have no other accounts in Org acme, so ccam cannot move this conversation."},
		{"cross_org", [2]string{"p1", "claude_max"}, [2]string{"acme", "claude_team"}, [2]string{"acme", "claude_team"}, true,
			hit + ". `ccam next --continue` continues the conversation on another account."},
		{"another seat, not logged in", [2]string{"acme", "claude_team"}, [2]string{"other", "claude_team"}, [2]string{"acme", "claude_team"}, false,
			hit + ". `ccam next --continue` continues the conversation on another account."},
	} {
		cfg := c.cfg()
		cfg.CrossOrg = k.crossOrg
		cfg.Save(c.p.ConfigFile())
		logIn(t, c.p.ProfileDir("personal"), "u1", k.personal[0], k.personal[1])
		logIn(t, c.p.ProfileDir("work"), "u2", k.work[0], k.work[1])
		logIn(t, c.p.ProfileDir("side"), "u3", k.sd[0], k.sd[1])
		out, _, _ := c.run(hookJSON(t, map[string]string{"error": "rate_limit", "last_assistant_message": "You've hit your session limit · resets 3:30am (Asia/Seoul)",
			"session_id": "s1", "transcript_path": conv, "cwd": cwd}), "hook", "stop-failure")
		var got map[string]string
		if err := json.Unmarshal([]byte(out), &got); err != nil || got["terminalSequence"] != "\x1b]9;"+k.want+"\x07" {
			t.Errorf("%s: %q, want %q", k.name, out, k.want)
		}
	}
}

// A refusal the conversation shows from before is not this stop's: in a
// shared history it can be another account's. Here personal was refused,
// work carried the conversation on, and work's next turn stopped otherwise.
func TestLimitHookSkipsAnEarlierLimit(t *testing.T) {
	hookEnv(t)
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 40})
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("work"), cwd), "s1.jsonl")
	now := time.Now()
	limitConversation(t, conv, now.Add(-20*time.Second), now.Add(3*time.Hour))
	f, _ := os.OpenFile(conv, os.O_APPEND|os.O_WRONLY, 0o600)
	fmt.Fprintln(f, `{"isSidechain":false,"type":"user","message":{"role":"user","content":"(go on)"},"origin":{"kind":"human"}}`)
	fmt.Fprintln(f, `{"isSidechain":false,"type":"assistant","message":{"model":"<synthetic>","role":"assistant","content":[{"type":"text","text":"API Error: Server is temporarily limiting requests (not your usage limit)"}]},"isApiErrorMessage":true,"error":"rate_limit"}`)
	f.Close()
	t.Setenv("CCAM_PROFILE", "work")
	run := func(msg string) string {
		out, _, _ := c.run(hookJSON(t, map[string]string{"error": "rate_limit", "last_assistant_message": msg, "session_id": "s1", "transcript_path": conv, "cwd": cwd}), "hook", "stop-failure")
		return out
	}
	record := filepath.Join(c.p.State, "usage", "work.limit.json")
	if out := run("API Error: Server is temporarily limiting requests (not your usage limit)"); out != "" || fileThere(record) {
		t.Fatalf("throttle: %q", out)
	}
	if _, ok := handoff.ForFolder(c.p.State, cwd, time.Now()); ok {
		t.Fatal("throttle: a note")
	}
	// work's own limit, before Claude Code has written it: the message
	// still makes the note, but personal's row is not work's limit.
	if out := run("You've hit your weekly limit · resets Oct 28, 9am"); !strings.Contains(out, "work hit its weekly limit (resets Oct 28, 9am)") || fileThere(record) {
		t.Fatalf("own limit not written yet: %q", out)
	}
}

// Claude Code writes the conversation out every 100 ms and does not wait
// for the hook, so the hook reads again until the refused request shows.
func TestLimitHookWaitsForTheRow(t *testing.T) {
	hookEnv(t)
	limitPoll = 5 * time.Second
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 40})
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	now := time.Now()
	limitConversation(t, conv, now, now.Add(time.Hour))
	data, _ := os.ReadFile(conv)
	lines := strings.SplitAfter(strings.TrimSpace(string(data)), "\n")
	os.WriteFile(conv, []byte(strings.Join(lines[:len(lines)-1], "")), 0o600)
	go func() {
		time.Sleep(300 * time.Millisecond)
		os.WriteFile(conv, data, 0o600)
	}()
	t.Setenv("CCAM_PROFILE", "personal")
	start := time.Now()
	in := hookJSON(t, map[string]string{"error": "rate_limit", "last_assistant_message": "You've hit your session limit · resets 3:30am (Asia/Seoul)",
		"session_id": "s1", "transcript_path": conv, "cwd": cwd})
	c.run(in, "hook", "stop-failure")
	if took := time.Since(start); !fileThere(filepath.Join(c.p.State, "usage", "personal.limit.json")) || took > 3*time.Second {
		t.Fatalf("took %v", took)
	}
}

// A limit row that ends the last turn but has another text than the
// message of this stop is not this stop's (a hop that carried no prompt
// leaves the other account's refused row last in a shared history): the
// message still makes the note, but nothing is recorded from the row.
func TestLimitHookNeedsTheSameMessage(t *testing.T) {
	hookEnv(t)
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 40})
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("work"), cwd), "s1.jsonl")
	now := time.Now()
	limitConversation(t, conv, now.Add(-20*time.Second), now.Add(3*time.Hour))
	t.Setenv("CCAM_PROFILE", "work")
	in := hookJSON(t, map[string]string{"error": "rate_limit", "last_assistant_message": "You've hit your weekly limit · resets Oct 28, 9am",
		"session_id": "s1", "transcript_path": conv, "cwd": cwd})
	out, _, _ := c.run(in, "hook", "stop-failure")
	if !strings.Contains(out, "work hit its weekly limit (resets Oct 28, 9am)") || fileThere(filepath.Join(c.p.State, "usage", "work.limit.json")) {
		t.Fatalf("%q", out)
	}
}
