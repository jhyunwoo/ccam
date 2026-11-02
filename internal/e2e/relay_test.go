package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/transcript"
)

// relayHarness has relay = true in config.toml and personal's conversation
// s1 in the folder app; it returns the folder and the conversation.
func relayHarness(t *testing.T) (*harness, string, string) {
	t.Helper()
	h := newHarness(t)
	h.addAccounts()
	cfgPath := filepath.Join(h.home, "config", "config.toml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Relay = true
	if err := cfg.Save(cfgPath); err != nil {
		t.Fatal(err)
	}
	project := filepath.Join(h.home, "app")
	os.MkdirAll(project, 0o700)
	project, _ = filepath.EvalSymlinks(project)
	return h, project, h.conversation("personal", project, "s1")
}

// conversation writes account's conversation id in folder, links resolved.
func (h *harness) conversation(account, folder, id string) string {
	h.t.Helper()
	conv := filepath.Join(h.profileDir(account), "projects", transcript.Slug(folder), id+".jsonl")
	os.MkdirAll(filepath.Dir(conv), 0o700)
	at := time.Now().Add(-time.Minute).UTC().Format(time.RFC3339Nano)
	os.WriteFile(conv, fmt.Appendf(nil, `{"type":"user","message":{"role":"user","content":"Fix the tests"},"timestamp":%q}`+"\n"+
		`{"type":"assistant","message":{"model":"claude-opus-5-5","content":[{"type":"text","text":"Done."}]},"timestamp":%q}`+"\n", at, at), 0o600)
	resolved, _ := filepath.EvalSymlinks(conv)
	return resolved
}

// switchEnv makes the fake claude of from leave /ccam:switch's note for
// conv to move to to, with extra settings.
func (h *harness) switchEnv(from, to, conv string, extra ...string) []string {
	return append([]string{"FAKE_SWITCH_FROM=" + from, "FAKE_SWITCH_TO=" + to, "FAKE_SWITCH_TRANSCRIPT=" + conv,
		"FAKE_SWITCH_STATE=" + filepath.Join(h.home, "state")}, extra...)
}

// notes lists the notes account's sessions left.
func (h *harness) notes(account string) []string {
	entries, _ := os.ReadDir(filepath.Join(h.home, "state", "handoff", account))
	var names []string
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".lock") { // the stop-failure hook's lock for subagents
			names = append(names, e.Name())
		}
	}
	return names
}

// runsIn reads FAKE_RUNS: the accounts claude ran as, in order.
func runsIn(path string) []string {
	data, _ := os.ReadFile(path)
	return strings.Fields(string(data))
}

// waitFile waits up to 5 s for path to appear.
func waitFile(t *testing.T, path string) {
	t.Helper()
	for deadline := time.Now().Add(5 * time.Second); !fileExists(path); time.Sleep(20 * time.Millisecond) {
		if time.Now().After(deadline) {
			t.Fatalf("%s did not appear", path)
		}
	}
}

// With relay on, ccam stays beside Claude Code, and after /ccam:switch (the
// hook ends Claude Code with SIGTERM) it carries the conversation on as the
// other account by itself: with the arguments resuming does not restore,
// and the permission mode the conversation was left in.
func TestRelayCarriesASwitchOn(t *testing.T) {
	h, project, conv := relayHarness(t)
	envFile, runs := filepath.Join(h.home, "env"), filepath.Join(h.home, "runs")
	line := ccamBin + " personal --model opus --permission-mode plan --add-dir /tmp --debug api"
	shown := h.onTerminal(project, line, "", h.switchEnv("personal", "work", conv, "FAKE_SWITCH_STOPPED=1",
		"FAKE_ENV_FILE="+envFile, "FAKE_ENV_NAMES=CCAM_RELAY", "FAKE_RUNS="+runs)...)
	dir, args := h.called()
	if want := []string{"--resume", conv, "--fork-session", "--permission-mode", "plan", "--add-dir", "/tmp"}; dir != h.profileDir("work") || !reflect.DeepEqual(args, want) {
		t.Fatalf("the relay started %s %q, want %q\n%s", dir, args, want, shown)
	}
	if !strings.Contains(shown, "ccam: not carried over: --debug") || !reflect.DeepEqual(runsIn(runs), []string{"personal", "work"}) || len(h.notes("personal")) != 0 {
		t.Fatalf("runs %q, notes %q:\n%s", runsIn(runs), h.notes("personal"), shown)
	}
	if b, _ := os.ReadFile(envFile); !regexp.MustCompile(`CCAM_RELAY=2:\d+`).Match(b) {
		t.Fatalf("the session does not know it is relayed: %q", b)
	}
	// The mode the conversation was left in wins over the one it started with.
	h.onTerminal(project, line, "", h.switchEnv("personal", "work", conv, "FAKE_SWITCH_STOPPED=1", "FAKE_MODE=acceptEdits")...)
	if _, args := h.called(); !reflect.DeepEqual(args, []string{"--resume", conv, "--fork-session", "--add-dir", "/tmp", "--permission-mode", "acceptEdits"}) {
		t.Fatalf("%q", args)
	}

	// A note another session wrote (another relay, or none) is not this
	// session's /ccam:switch: the relay ends with claude.
	h.onTerminal(project, ccamBin+" personal", "", h.switchEnv("personal", "work", conv, "FAKE_SWITCH_RELAY=2:1")...)
	if dir, _ := h.called(); dir != h.profileDir("personal") {
		t.Fatalf("the relay carried on another session's note as %s", dir)
	}
}

// A usage limit stops every session of an account at once: two sessions
// of one account switch together, each carries on its own conversation in
// its own folder, and neither takes the other's note.
func TestRelayTwoSessionsOneAccount(t *testing.T) {
	h, _, _ := relayHarness(t)
	type session struct {
		dir, conv, log, runs string
		term                 *terminal
	}
	var sessions []*session
	for _, name := range []string{"app1", "app2"} {
		dir := filepath.Join(h.home, name)
		os.MkdirAll(dir, 0o700)
		dir, _ = filepath.EvalSymlinks(dir)
		sessions = append(sessions, &session{dir: dir, conv: h.conversation("personal", dir, "conv-"+name),
			log: filepath.Join(h.home, name+".log"), runs: filepath.Join(h.home, name+".runs")})
	}
	for round := 0; round < 10; round++ {
		for _, s := range sessions {
			os.Remove(s.runs)
			s.term = h.startTerminal(s.dir, ccamBin+" personal", h.switchEnv("personal", "work", s.conv,
				"FAKE_SWITCH_STOPPED=1", "FAKE_EXIT_DELAY=0.3", "FAKE_CLAUDE_LOG="+s.log, "FAKE_RUNS="+s.runs)...)
		}
		for _, s := range sessions {
			shown, err := s.term.wait(20 * time.Second)
			if err != nil {
				t.Fatal(err)
			}
			data, _ := os.ReadFile(s.log)
			args := strings.Split(strings.TrimSuffix(string(data), "\x00"), "\x00")
			if !reflect.DeepEqual(runsIn(s.runs), []string{"personal", "work"}) || len(args) < 4 || args[0] != "CLAUDE_CONFIG_DIR="+h.profileDir("work") ||
				!reflect.DeepEqual(args[2:], []string{"--resume", s.conv, "--fork-session"}) {
				t.Fatalf("round %d, %s: runs %q, last claude %q\n%s", round, s.dir, runsIn(s.runs), args, shown)
			}
		}
		if notes := h.notes("personal"); len(notes) != 0 {
			t.Fatalf("round %d: notes left %q", round, notes)
		}
	}
}

// The hop is ccam's own: it does not ask whether to run beside another
// session of the shared history, which nobody would answer.
func TestRelayHopDoesNotAsk(t *testing.T) {
	h, project, _ := relayHarness(t)
	h.ok("ccam", "group", "set", "default", "--add", "history", "--yes")
	conv := h.conversation("personal", project, "s1") // now in the shared history
	runs := filepath.Join(h.home, "runs")
	start := time.Now()
	term := h.startTerminal(project, ccamBin+" personal", h.switchEnv("personal", "work", conv, "FAKE_SWITCH_STOPPED=1", "FAKE_EXIT_DELAY=3", "FAKE_RUNS="+runs)...)
	waitFile(t, filepath.Join(h.home, "state", "handoff", "personal", "s1.json"))
	h.holdFolder(project, "work") // another session of the shared history in this folder
	shown, err := term.wait(20 * time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(shown, "Run anyway?") || !strings.Contains(shown, "ccam is continuing its own conversation, so it does not ask.") ||
		!reflect.DeepEqual(runsIn(runs), []string{"personal", "work"}) || time.Since(start) > 10*time.Second {
		t.Fatalf("runs %q after %v:\n%s", runsIn(runs), time.Since(start), shown)
	}
}

// When the person ended Claude Code first (/exit before the hook's
// SIGTERM), or a signal reached the relay itself (the terminal hung up, the
// system shuts down), or Ctrl+C came while the session ended, the relay
// does not carry on: it ends with claude, and the note stays for `ccam next
// --continue`.
func TestRelayStopsWhenThePersonOrASignalDoes(t *testing.T) {
	h, project, conv := relayHarness(t)
	runs := filepath.Join(h.home, "runs")
	for _, c := range []struct {
		name string
		env  []string
		says string
	}{
		{"/exit first", []string{"FAKE_EXIT=0"}, "ccam: Claude Code closed before ccam moved the conversation, so it stays on personal. `ccam next --continue` continues it on another account."},
		{"SIGHUP", []string{"FAKE_KILL_PARENT=HUP", "FAKE_EXIT_DELAY=2"}, ""},
		{"SIGTERM", []string{"FAKE_KILL_PARENT=TERM", "FAKE_EXIT_DELAY=2"}, ""},
		{"Ctrl+C", []string{"FAKE_KILL_PARENT=INT", "FAKE_EXIT_DELAY=0.5"}, "ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation."},
	} {
		os.Remove(runs)
		shown := h.onTerminal(project, ccamBin+" personal", "", h.switchEnv("personal", "work", conv, append(c.env, "FAKE_SWITCH_STOPPED=1", "FAKE_RUNS="+runs)...)...)
		if !reflect.DeepEqual(runsIn(runs), []string{"personal"}) || !strings.Contains(shown, c.says) || !reflect.DeepEqual(h.notes("personal"), []string{"s1.json"}) {
			t.Fatalf("%s: runs %q, notes %q:\n%s", c.name, runsIn(runs), h.notes("personal"), shown)
		}
	}
	// The note carries the conversation on later.
	if r := h.runIn(project, "", nil, "ccam", "next", "--continue"); r.code != 0 {
		t.Fatalf("%+v", r)
	}
	if dir, args := h.called(); dir != h.profileDir("work") || !reflect.DeepEqual(args, []string{"--resume", conv, "--fork-session"}) || len(h.notes("personal")) != 0 {
		t.Fatalf("%s %q, notes %q", dir, args, h.notes("personal"))
	}
}

// A terminal that is gone (its window closed while Claude Code ended) gets
// no new session.
func TestRelayDeadTerminal(t *testing.T) {
	h, project, conv := relayHarness(t)
	runs := filepath.Join(h.home, "runs")
	term := h.startTerminal(project, ccamBin+" personal", h.switchEnv("personal", "work", conv, "FAKE_SWITCH_STOPPED=1", "FAKE_EXIT_DELAY=3", "FAKE_RUNS="+runs)...)
	waitFile(t, filepath.Join(h.home, "state", "handoff", "personal", "s1.json"))
	term.cmd.Process.Kill() // script, which holds the terminal
	term.wait(5 * time.Second)
	time.Sleep(4 * time.Second) // past the end the fake claude had in mind
	if !reflect.DeepEqual(runsIn(runs), []string{"personal"}) || !reflect.DeepEqual(h.notes("personal"), []string{"s1.json"}) {
		t.Fatalf("runs %q, notes %q", runsIn(runs), h.notes("personal"))
	}
}

// What was typed while the conversation moved is dropped, not sent to the
// next Claude Code as its first prompt.
func TestRelayFlushesTypeAhead(t *testing.T) {
	h, project, conv := relayHarness(t)
	typed := filepath.Join(h.home, "typed")
	term := h.startTerminal(project, ccamBin+" personal", h.switchEnv("personal", "work", conv, "FAKE_SWITCH_STOPPED=1", "FAKE_EXIT_DELAY=1.5", "FAKE_TTY_INPUT="+typed)...)
	waitFile(t, filepath.Join(h.home, "state", "handoff", "personal", "s1.json"))
	time.Sleep(200 * time.Millisecond)
	term.in.Write([]byte("rm -rf build\n"))
	shown, err := term.wait(20 * time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := os.ReadFile(typed); err != nil || len(got) != 0 || !strings.Contains(shown, "ccam: what you typed during the switch was discarded.") {
		t.Fatalf("the next claude read %q (%v):\n%s", got, err, shown)
	}
}

// A hook of this ccam ends Claude Code only under a relay of this ccam, and
// only the claude that relay runs: a relay of an older ccam (CCAM_RELAY a
// bare pid), started before an update, would not carry the note on. On
// Windows the hook ends nothing, and only the relay's own claude is told
// that ccam carries the conversation on after /exit.
func TestRelayOldProtocol(t *testing.T) {
	windows := runtime.GOOS == "windows"
	if _, err := exec.LookPath("sh"); err != nil && !windows {
		t.Skip("needs sh")
	}
	h := newHarness(t)
	h.addAccounts()
	conv := filepath.Join(h.home, "s1.jsonl")
	os.WriteFile(conv, []byte("{}\n"), 0o600)
	in := fmt.Sprintf(`{"session_id":"s1","transcript_path":%q,"cwd":%q,"command_args":"work"}`, conv, h.home)
	for _, c := range []struct {
		relay string
		ours  bool // the relay's own claude: this test started it
	}{{"%d", false}, {"2:1", false}, {"2:%d", true}} {
		ends := c.ours && !windows
		claude := exec.Command(fakeBin) // a claude this test runs
		claude.Env = append(append([]string{}, h.env...), "FAKE_ONLY_SLEEP=5")
		if err := claude.Start(); err != nil {
			t.Fatal(err)
		}
		ended := make(chan struct{})
		go func() { claude.Wait(); close(ended) }()
		relay := c.relay
		if strings.Contains(relay, "%d") {
			relay = fmt.Sprintf(relay, os.Getpid())
		}
		r := h.run(in, []string{"CCAM_PROFILE=personal", "CCAM_RELAY=" + relay, "CLAUDE_PID=" + strconv.Itoa(claude.Process.Pid)}, "ccam", "hook", "switch")
		select {
		case <-ended:
			if !ends {
				t.Fatalf("CCAM_RELAY=%s ended claude: %+v", relay, r)
			}
		case <-time.After(2 * time.Second):
			claude.Process.Kill()
			<-ended
			if ends {
				t.Fatalf("CCAM_RELAY=%s left claude running: %+v", relay, r)
			}
		}
		want := "/exit, then run `ccam next --continue`"
		switch {
		case ends:
			want = "moving this conversation to work"
		case c.ours:
			want = "type /exit; ccam then continues this conversation on work."
		}
		if !strings.Contains(r.stdout, want) {
			t.Fatalf("CCAM_RELAY=%s: %+v", relay, r)
		}
	}
}

// `claude remote-control` serves sessions of its own: ccam replaces itself
// with it and runs no relay beside it.
func TestRelayNotForRemoteControlServer(t *testing.T) {
	h, project, _ := relayHarness(t)
	envFile := filepath.Join(h.home, "env")
	// Nor for -p runs typed on a terminal: those stay an exec.
	for _, args := range []string{"remote-control", "-p hi", "--print=hi"} {
		os.Remove(envFile)
		h.onTerminal(project, ccamBin+" personal "+args, "", "FAKE_ENV_FILE="+envFile, "FAKE_ENV_NAMES=CCAM_RELAY")
		if b, err := os.ReadFile(envFile); err != nil || strings.Contains(string(b), "CCAM_RELAY") {
			t.Fatalf("%s: %q %v", args, b, err)
		}
	}
	h.onTerminal(project, ccamBin+" personal", "", "FAKE_ENV_FILE="+envFile, "FAKE_ENV_NAMES=CCAM_RELAY")
	if b, _ := os.ReadFile(envFile); !strings.Contains(string(b), "CCAM_RELAY=2:") {
		t.Fatalf("an interactive session is relayed: %q", b)
	}
}
