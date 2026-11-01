package cli

import (
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
	"github.com/jhyunwoo/ccam/internal/transcript"
)

// A /ccam:switch note that names an account outside the pool of the
// conversation's account waits for a yes on the terminal, cross_org or
// not, and so does one that names none when ccam picks such an account;
// without a terminal the answer is no, whatever stdin holds. An API
// key account (ccam's own, or a Console login) is in no pool, on either
// side of the move. Nothing else asks: another seat of the organization,
// the account itself, a dry run, the account ccam chose after a limit, or
// a login whose plan ccam cannot tell.
func TestConfirmSwitchOnlyOutsideThePool(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 20, "side": 30})
	logIn(t, c.p.ProfileDir("personal"), "u1", "acme", "claude_team")
	logIn(t, c.p.ProfileDir("work"), "u2", "acme", "claude_team")
	logIn(t, c.p.ProfileDir("side"), "u3", "home", "claude_max")
	// ccam asks on the process's stdin, which go test may leave a terminal.
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	w.WriteString("y\n")
	oldIn := os.Stdin
	os.Stdin = r
	t.Cleanup(func() { os.Stdin = oldIn; r.Close(); w.Close() })
	e, cfg := c.env()
	cfg.CrossOrg = true
	cfg.Profiles = append(cfg.Profiles, config.Profile{Name: "gw", Slot: 4, API: &config.API{BaseURL: "https://llm.example.com/v1", KeyFrom: "env:GW_KEY"}},
		config.Profile{Name: "console", Slot: 5})
	for _, name := range []string{"gw", "console"} {
		os.MkdirAll(c.p.ProfileDir(name), 0o700)
	}
	os.WriteFile(filepath.Join(c.p.ProfileDir("console"), ".claude.json"), []byte(`{"primaryApiKey":"sk-ant-x"}`), 0o600)
	on := func(account, target string) *conversation {
		p, _ := cfg.Find(account)
		return &conversation{path: filepath.Join(transcript.ProjectDir(c.p.ProfileDir(account), t.TempDir()), "s1.jsonl"), owner: p,
			ticket: &handoff.Ticket{Account: account, SessionID: "s1", At: time.Now(), Reason: handoff.Switch, Target: target}}
	}
	stays := func(conv *conversation, to string, says string) {
		t.Helper()
		p, _ := cfg.Find(to)
		picked, again := conv.ticket.Target == "", "asks"
		if picked {
			again = "picks"
		}
		err := confirmSwitch(e, cfg, conv, p, picked)
		from := conv.owner.Name
		if err == nil || !strings.HasPrefix(err.Error(), "the conversation stays on "+says+" only when you answer yes on a terminal.") ||
			!strings.HasSuffix(err.Error(), "`ccam next --continue` "+again+" again, and `ccam "+from+" --resume s1` continues it as "+from) {
			t.Fatalf("%s to %s: %v", from, p.Name, err)
		}
	}
	stays(on("work", "side"), "side", "work, a Claude Team seat in Org acme: it moves to side, a personal plan,")
	stays(on("work", "gw"), "gw", "work, a Claude Team seat in Org acme: it moves to gw, an API key account,")
	stays(on("work", "console"), "console", "work, a Claude Team seat in Org acme: it moves to console, an API key account,")
	stays(on("gw", "work"), "work", "gw, an API key account: it moves to work, a Claude Team seat in Org acme,")
	stays(on("gw", "console"), "console", "gw, an API key account: it moves to console, an API key account,")
	// /ccam:switch named no account, and cross_org let ccam pick one outside the pool.
	stays(on("work", ""), "side", "work, a Claude Team seat in Org acme: it moves to side, a personal plan,")

	conv := on("work", "side")
	moves := func(what, to string) {
		t.Helper()
		p, _ := cfg.Find(to)
		if err := confirmSwitch(e, cfg, conv, p, false); err != nil {
			t.Fatalf("%s: %v", what, err)
		}
	}
	moves("another seat of the organization", "personal")
	moves("the account itself", "work")
	e.dryRun = true
	moves("a dry run", "side")
	e.dryRun, conv.ticket.Reason, conv.ticket.Target, conv.ticket.Chosen = false, handoff.Limit, "", "side"
	moves("the account ccam chose after a limit", "side")
	conv.ticket.Reason, conv.ticket.Target, conv.ticket.Chosen = handoff.Switch, "side", ""
	os.Remove(filepath.Join(c.p.ProfileDir("side"), ".claude.json"))
	moves("a login whose plan ccam cannot tell", "side")
}

// Under a relay, /ccam:switch to an account outside the pool leaves
// Claude Code running: the hook cannot tell a /ccam:switch the person typed
// from one that a job the model scheduled fired, and ending the session for
// such a job would let it end the session again at every run. The note
// waits for the person's own /exit, after which the relay asks. Within the
// pool the hook ends Claude Code as before (except on Windows, where it
// never does).
func TestSwitchHookLeavesClaudeCodeRunningOutsideThePool(t *testing.T) {
	hookEnv(t)
	c := newTestCLI(t, nil)
	logIn(t, c.p.ProfileDir("personal"), "u1", "acme", "claude_team")
	logIn(t, c.p.ProfileDir("work"), "u2", "acme", "claude_team")
	logIn(t, c.p.ProfileDir("side"), "u3", "home", "claude_max")
	cwd, _ := os.Getwd()
	conv := filepath.Join(t.TempDir(), "s2.jsonl")
	os.WriteFile(conv, []byte("{}\n"), 0o600)
	var ended []int
	old := endClaude
	endClaude = func(pid int) error { ended = append(ended, pid); return nil }
	t.Cleanup(func() { endClaude = old })
	// This test is the relay; its child stands in for Claude Code.
	claude := exec.Command("sleep", "5")
	if runtime.GOOS == "windows" {
		claude = exec.Command("ping", "-n", "6", "127.0.0.1") // 5 s
	}
	if err := claude.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { claude.Process.Kill(); claude.Wait() }()
	pid, me := claude.Process.Pid, "2:"+strconv.Itoa(os.Getpid())
	t.Setenv("CCAM_PROFILE", "work")
	t.Setenv("CLAUDE_PID", strconv.Itoa(pid))
	t.Setenv("CCAM_RELAY", me)
	run := func(to string) (string, *handoff.Ticket) {
		t.Helper()
		ended = nil
		out, _, _ := c.run(hookJSON(t, map[string]string{"session_id": "s2", "transcript_path": conv, "cwd": cwd, "command_args": to}), "hook", "switch")
		tk, _ := handoff.ForFolder(c.p.State, cwd, time.Now())
		return stopReason(t, out), tk
	}
	r, tk := run("side")
	if r != "ccam: type /exit; ccam then asks whether to continue this conversation on side, a personal plan (work is a Claude Team seat in Org acme)." ||
		len(ended) != 0 || tk.Target != "side" || tk.Relay != me || tk.Child != pid || tk.Stopped {
		t.Fatalf("%q %v %+v", r, ended, tk)
	}
	if runtime.GOOS == "windows" {
		return
	}
	if r, tk := run("personal"); !strings.Contains(r, "moving this conversation to personal") || !reflect.DeepEqual(ended, []int{pid}) || !tk.Stopped {
		t.Fatalf("within the pool: %q %v %+v", r, ended, tk)
	}
}
