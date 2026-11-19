package launcher

import (
	"cmp"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/session"
)

func TestHopOK(t *testing.T) {
	start := time.Now()
	note := func(edit func(*handoff.Ticket)) *handoff.Ticket {
		n := &handoff.Ticket{Account: "p1", Reason: handoff.Switch, Child: 100, Stopped: true, At: start.Add(10 * time.Second)}
		if edit != nil {
			edit(n)
		}
		return n
	}
	base := hop{child: 100, code: 143, start: start, terminal: true}
	cancelled, closed := "switch cancelled", "closed before ccam moved"
	for _, c := range []struct {
		name string
		note *handoff.Ticket
		edit func(*hop)
		ok   bool
		why  string
	}{
		{"the hook ended claude", note(nil), nil, true, ""},
		{"/exit came first", note(nil), func(h *hop) { h.code = 0 }, false, closed},
		{"something else ended claude", note(nil), func(h *hop) { h.code = 1 }, false, closed},
		{"the relay got a signal", note(nil), func(h *hop) { h.signaled = true }, false, ""},
		{"Ctrl+C after the note", note(nil), func(h *hop) { h.interrupt = start.Add(11 * time.Second) }, false, cancelled},
		{"Ctrl+C before the note", note(nil), func(h *hop) { h.interrupt = start.Add(5 * time.Second) }, true, ""},
		{"the terminal is gone", note(nil), func(h *hop) { h.terminal = false }, false, ""},
		{"another claude's note", note(nil), func(h *hop) { h.child = 101 }, false, ""},
		{"a note from before claude started", note(func(n *handoff.Ticket) { n.At = start.Add(-time.Second) }), nil, false, ""},
		// Windows: the note names claude, but nothing was ended; the person typed /exit.
		{"nothing ended, /exit", note(func(n *handoff.Ticket) { n.Stopped = false }), func(h *hop) { h.code = 0 }, true, ""},
		{"nothing ended, another claude's note", note(func(n *handoff.Ticket) { n.Stopped = false }), func(h *hop) { h.code, h.child = 0, 101 }, false, ""},
		// Claude Code without CLAUDE_PID: nothing was ended, the person typed /exit.
		{"no process named, /exit", note(func(n *handoff.Ticket) { n.Child, n.Stopped = 0, false }), func(h *hop) { h.code = 0 }, true, ""},
		{"no process named, SIGTERM", note(func(n *handoff.Ticket) { n.Child, n.Stopped = 0, false }), nil, false, closed},
		{"stopped without a process", note(func(n *handoff.Ticket) { n.Child = 0 }), nil, false, ""},
		// A limit hops only to an account the hook chose; Ctrl+C counts from the limit.
		{"a limit, nobody chosen", note(func(n *handoff.Ticket) { n.Reason = handoff.Limit }), nil, false, ""},
		{"a limit, chosen", note(func(n *handoff.Ticket) { n.Reason, n.Chosen = handoff.Limit, "p2" }), nil, true, ""},
		{"Ctrl+C after the limit", note(func(n *handoff.Ticket) { n.Reason, n.Chosen, n.LimitAt = handoff.Limit, "p2", start.Add(2*time.Second) }),
			func(h *hop) { h.interrupt = start.Add(5 * time.Second) }, false, cancelled},
	} {
		h := base
		h.note = c.note
		if c.edit != nil {
			c.edit(&h)
		}
		ok, why := hopOK(h)
		if ok != c.ok || c.why == "" && why != "" || !strings.Contains(why, c.why) {
			t.Errorf("%s: %v %q", c.name, ok, why)
		}
	}
}

// The relay says what it does after a limit: a notification of the move
// with why that account, and a line before the hop.
func TestLimitHopSays(t *testing.T) {
	now := time.Now()
	reset := now.Add(3 * time.Hour)
	note := &handoff.Ticket{Account: "p1", Reason: handoff.Limit, Window: "five_hour", ResetsAt: reset, Chosen: "p2", Why: "12% of its 5 hours used, the least"}
	notice, line := limitHop(note, now)
	if want := "ccam: p1 hit its 5-hour limit (resets " + handoff.FormatReset(reset, now) + "). Moving this conversation to p2 (12% of its 5 hours used, the least)."; notice != want {
		t.Errorf("%q, want %q", notice, want)
	}
	if line != "ccam: p1 hit its 5-hour limit; continuing this conversation as p2" {
		t.Error(line)
	}
	note.Why, note.Window, note.ResetsAt = "", "seven_day", time.Time{}
	if notice, line = limitHop(note, now); notice != "ccam: p1 hit its weekly limit. Moving this conversation to p2." || line != "ccam: p1 hit its weekly limit; continuing this conversation as p2" {
		t.Errorf("%q %q", notice, line)
	}
}

// The relay runs with relay = true, and with auto_switch = true, which
// needs it, except on Windows, which ignores auto_switch (D6).
func TestUseRelay(t *testing.T) {
	for _, c := range []struct {
		relay, auto bool
		goos        string
		want        bool
	}{
		{false, false, "linux", false},
		{true, false, "linux", true},
		{true, false, "windows", true},
		{false, true, "linux", true},
		{false, true, "darwin", true},
		{false, true, "windows", false},
		{true, true, "windows", true},
	} {
		if got := useRelay(&config.Config{Relay: c.relay, AutoSwitch: c.auto}, c.goos); got != c.want {
			t.Errorf("relay %v, auto_switch %v on %s: %v", c.relay, c.auto, c.goos, got)
		}
	}
}

func TestRelayPID(t *testing.T) {
	for v, want := range map[string]int{"2:1234": 1234, "1234": 0, "": 0, "2:": 0, "2:x": 0, "2:-5": 0, "3:12": 0} {
		if got, ok := RelayPID(v); got != want || ok != (want != 0) {
			t.Errorf("%q: %d %v", v, got, ok)
		}
	}
}

// TestHelperReexec stands in for Claude Code that starts itself again
// (as gateway settings make it): this test binary runs another copy of
// itself (or CCAM_TEST_INNER), which says its pid and sleeps.
func TestHelperReexec(t *testing.T) {
	switch os.Getenv("CCAM_TEST_REEXEC") {
	case "outer":
		c := exec.Command(cmp.Or(os.Getenv("CCAM_TEST_INNER"), os.Args[0]), "-test.run=^TestHelperReexec$")
		c.Env = append(os.Environ(), "CCAM_TEST_REEXEC=inner")
		c.Stdout = os.Stdout
		c.Run()
	case "inner":
		os.Stdout.WriteString(strconv.Itoa(os.Getpid()) + "\n")
		time.Sleep(5 * time.Second)
	}
}

// Only the relay's own claude counts: its child, or that child's child when
// both are the same program. A claude started under it otherwise (in its
// shell) does not.
func TestRelayChild(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("no process tree here (Windows: TestRelayChildOnWindows)")
	}
	sleep := exec.Command("sleep", "5")
	if err := sleep.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { sleep.Process.Kill(); sleep.Wait() }()
	if c, ok := RelayChild(sleep.Process.Pid, os.Getpid()); !ok || c != sleep.Process.Pid {
		t.Fatalf("own child: %d %v", c, ok)
	}
	if _, ok := RelayChild(sleep.Process.Pid, os.Getppid()); ok {
		t.Fatal("another process's grandchild")
	}
	// sh runs sleep: a grandchild, but another program.
	sh := exec.Command("sh", "-c", "sleep 5 & echo $!; wait")
	out, _ := sh.StdoutPipe()
	sh.Start()
	defer func() { sh.Process.Kill(); sh.Wait() }()
	var buf [32]byte
	n, _ := out.Read(buf[:])
	pid, _ := strconv.Atoi(strings.TrimSpace(string(buf[:n])))
	// Until it runs sleep, the forked child is sh still.
	sleepBin, _ := filepath.EvalSymlinks(sleep.Path)
	for deadline := time.Now().Add(2 * time.Second); time.Now().Before(deadline); time.Sleep(10 * time.Millisecond) {
		if exe, _ := session.Exe(pid); exe == sleepBin {
			break
		}
	}
	shBin, _ := exec.LookPath("sh")
	a, err1 := os.Stat(shBin)
	b, err2 := os.Stat(sleepBin)
	if _, ok := RelayChild(pid, os.Getpid()); ok && err1 == nil && err2 == nil && !os.SameFile(a, b) { // busybox is both
		t.Fatal("sh's child")
	}
	// The same program started again.
	outer := exec.Command(os.Args[0], "-test.run=^TestHelperReexec$")
	outer.Env = append(os.Environ(), "CCAM_TEST_REEXEC=outer")
	pipe, _ := outer.StdoutPipe()
	outer.Start()
	defer func() { outer.Process.Kill(); outer.Wait() }()
	n, _ = pipe.Read(buf[:])
	inner, _ := strconv.Atoi(strings.TrimSpace(string(buf[:n])))
	if c, ok := RelayChild(inner, os.Getpid()); !ok || c != outer.Process.Pid {
		t.Fatalf("started again: %d %v, want %d", c, ok, outer.Process.Pid)
	}
}

// relayable leaves out the Remote Control server (the subcommand), not a
// session with Remote Control on (the flag).
func TestRemoteControlServer(t *testing.T) {
	rc := func(args []string) bool { return claudeVerb(args) == verbRemoteControl }
	if !rc([]string{"--verbose", "remote-control"}) || !rc([]string{"rc"}) || rc([]string{"--remote-control"}) {
		t.Fatal("remote-control is the server subcommand, --remote-control a session")
	}
}

// Prepare keeps the args as the user gave them for the hop; the ones ccam
// adds are not among them.
func TestUserArgsAreTheUsers(t *testing.T) {
	p, cfg := setup(t)
	pl, err := Prepare(p, cfg, "work", []string{"--", "--add-dir", "/tmp"})
	if err != nil {
		t.Fatal(err)
	}
	pl.InsertArgs("--mcp-config=/g/mcp.json")
	if strings.Join(pl.UserArgs, " ") != "--add-dir /tmp" || pl.Argv[1] != "--mcp-config=/g/mcp.json" {
		t.Fatalf("%q %q", pl.UserArgs, pl.Argv)
	}
}

func TestHistoryQuestionArgs(t *testing.T) {
	for args, want := range map[string][2]bool{ // opens a session, new -p run
		"":                     {true, false},
		"-p task":              {true, true},
		"--print=task":         {true, true},
		"--model opus -p task": {true, true},
		"-p -c task":           {true, false},
		"-p --resume=abc":      {true, false},
		"-p --session-id x":    {true, true},
		"-p --resume=abc --fork-session --session-id x": {true, true},
		"-p -c --fork-session":                          {true, true},
		"--resume abc --fork-session":                   {true, false},
		"-p -- --resume x":                              {true, true},
		"-p -- -c":                                      {true, true},
		"--resume abc":                                  {true, false},
		"--version":                                     {false, false},
		"-h":                                            {false, false},
		"mcp list":                                      {false, false},
		"mcp add srv -- npx -y p":                       {false, false},
		"auth status":                                   {false, false},
		"plugin install x":                              {false, false},
		"attach a1":                                     {true, false},
		"remote-control":                                {true, false},
		"-- --version":                                  {true, false},
		"--add-dir mcp --verbose":                       {true, false},
	} {
		a := strings.Fields(args)
		if got := [2]bool{opensSession(a), newPrint(a)}; got != want {
			t.Errorf("%q: %v, want %v", args, got, want)
		}
	}
}
