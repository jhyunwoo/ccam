package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/transcript"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// writeConversation writes a conversation with a prompt and a reply made at
// at, then rows, each a line of the file.
func writeConversation(t *testing.T, path string, at time.Time, rows ...string) {
	t.Helper()
	ts := at.UTC().Format(time.RFC3339Nano)
	lines := []string{
		fmt.Sprintf(`{"parentUuid":null,"isSidechain":false,"type":"user","message":{"role":"user","content":"Fix the failing tests"},"timestamp":%q}`, ts),
		fmt.Sprintf(`{"parentUuid":"u1","isSidechain":false,"type":"assistant","message":{"model":"claude-opus-5-5","role":"assistant","content":[{"type":"text","text":"Done."}],"usage":{"input_tokens":10,"output_tokens":5}},"timestamp":%q}`, ts),
	}
	os.MkdirAll(filepath.Dir(path), 0o700)
	if err := os.WriteFile(path, []byte(strings.Join(append(lines, rows...), "\n")+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

// appendRows adds rows to the conversation at path, as Claude Code does.
func appendRows(t *testing.T, path string, rows ...string) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, r := range rows {
		fmt.Fprintln(f, r)
	}
}

// What Claude Code writes when a session ends: no prompt, no reply.
var exitRows = []string{`{"type":"cost-state","sessionId":"s1","totalCostUSD":0.4}`, `{"type":"last-prompt","lastPrompt":"/exit","sessionId":"s1"}`}

func trueClaude(t *testing.T, c *testCLI) {
	t.Helper()
	cfg, _ := config.Load(c.p.ConfigFile())
	cfg.ClaudeBin = stubClaude
	cfg.Save(c.p.ConfigFile())
}

// Claude Code writes its rows at /exit when the session ends, which may be
// long after /ccam:switch or the limit: they are no sign the conversation
// went on, so the note still says where it goes. A prompt after the note
// is.
func TestNoteOutlastsTheRowsWrittenAtExit(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 50, "work": 10, "side": 40})
	trueClaude(t, c)
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("work"), cwd), "s1.jsonl")
	now := time.Now()
	writeConversation(t, conv, now.Add(-10*time.Minute))
	appendRows(t, conv, exitRows...)
	note := handoff.Ticket{Account: "work", SessionID: "s1", Transcript: conv, Cwd: cwd, At: now.Add(-5 * time.Minute), Reason: handoff.Switch, Target: "personal"}
	handoff.Save(c.p.State, note)
	if _, errOut, code := c.run("", "--dry-run", "next", "--continue"); code != 0 || !strings.Contains(errOut, "→ work → personal") {
		t.Fatalf("switch note: %d %s", code, errOut)
	}
	// A limit note: --auto leaves the account that hit it out (work has
	// the most room on paper).
	note.Reason, note.Target, note.Limit = handoff.Limit, "", "session"
	handoff.Save(c.p.State, note)
	if _, errOut, _ := c.run("", "--dry-run", "next", "--auto", "--continue"); !strings.Contains(errOut, "→ work → side") {
		t.Fatalf("limit note: %s", errOut)
	}
	// Another conversation in the folder that only got rows of its own
	// since does not win over the note either.
	other := filepath.Join(filepath.Dir(conv), "s2.jsonl")
	writeConversation(t, other, now.Add(-20*time.Minute), exitRows...)
	if _, errOut, _ := c.run("", "--dry-run", "next", "--auto", "--continue"); !strings.Contains(errOut, "→ work → side") {
		t.Fatalf("with another conversation: %s", errOut)
	}
	// A prompt after the note: the conversation went on, the note is old.
	note.Reason, note.Target, note.Limit = handoff.Switch, "personal", ""
	handoff.Save(c.p.State, note)
	writeConversation(t, conv, now.Add(-10*time.Minute), fmt.Sprintf(`{"isSidechain":false,"type":"user","message":{"role":"user","content":"go on"},"timestamp":%q}`, now.UTC().Format(time.RFC3339Nano)))
	if out, errOut, _ := c.run("", "--dry-run", "next", "--continue"); !strings.Contains(errOut, "→ work → side") || !strings.Contains(out, "s1.jsonl") {
		t.Fatalf("after a prompt: %s %s", errOut, out)
	}
}

// A session that only ran /mcp, /usage or /status, or an empty file, has
// no conversation to carry on; the conversation before it does.
func TestContinueSkipsSessionsWithoutAConversation(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 50, "work": 10, "side": 40})
	trueClaude(t, c)
	cwd, _ := os.Getwd()
	dir := transcript.ProjectDir(c.p.ProfileDir("work"), cwd)
	conv := filepath.Join(dir, "s1.jsonl")
	writeConversation(t, conv, time.Now().Add(-8*time.Minute))
	old := time.Now().Add(-8 * time.Minute)
	os.Chtimes(conv, old, old)
	for name, body := range map[string]string{
		// The rows of a session that only ran /mcp (2.1.306, as in a real one).
		"mcp": `{"type":"last-prompt","lastPrompt":"/mcp"}
{"type":"mode","mode":"normal"}
{"parentUuid":null,"isSidechain":false,"type":"attachment","attachment":{"type":"hook_success"}}
{"parentUuid":"t1","isSidechain":false,"type":"user","message":{"role":"user","content":"<local-command-caveat>Caveat: the messages below were generated by the user while running local commands.</local-command-caveat>"},"isMeta":true}
{"parentUuid":"c1","isSidechain":false,"type":"user","message":{"role":"user","content":"<command-name>/mcp</command-name>\n<command-message>mcp</command-message>\n<command-args></command-args>"}}
{"parentUuid":"c2","isSidechain":false,"type":"user","message":{"role":"user","content":"<local-command-stdout>MCP dialog dismissed</local-command-stdout>"}}
{"type":"cost-state"}
`,
		"usage": `{"type":"mode","mode":"normal"}
{"type":"permission-mode","permissionMode":"default"}
{"parentUuid":null,"type":"system","subtype":"local_command","content":"<command-name>/usage</command-name>"}
{"parentUuid":"s1","type":"system","subtype":"local_command","content":"<local-command-stdout>Usage dialog dismissed</local-command-stdout>"}
{"type":"last-prompt","lastPrompt":"/usage"}
{"type":"cost-state"}
`,
		"empty": "",
	} {
		path := filepath.Join(dir, "s2.jsonl")
		os.WriteFile(path, []byte(body), 0o600)
		if out, errOut, code := c.run("", "--dry-run", "next", "--continue"); code != 0 || !strings.Contains(out, "s1.jsonl") {
			t.Fatalf("%s: %d %s %s", name, code, out, errOut)
		}
		os.Remove(path)
	}
}

// A rename carries the account's notes along: a note of the old name is the
// new name's, with its conversation where the folder moved to, and a note
// that named the old name as its target names the new one.
func TestRenameKeepsTheNotes(t *testing.T) {
	hookEnv(t)
	c := newTestCLI(t, map[string]int{"personal": 50, "work": 10, "side": 40})
	trueClaude(t, c)
	cwd, _ := os.Getwd()
	switchTo := func(from, conv, target string) {
		t.Helper()
		writeConversation(t, conv, time.Now())
		t.Setenv("CCAM_PROFILE", from)
		t.Setenv("CCAM_RELAY", "")
		in := hookJSON(t, map[string]string{"session_id": transcript.ID(conv), "transcript_path": conv, "cwd": cwd, "command_args": target})
		if out, _, code := c.run(in, "hook", "switch"); code != 0 || !strings.Contains(out, "/exit") {
			t.Fatalf("%d %s", code, out)
		}
	}
	switchTo("personal", filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl"), "side")
	if _, errOut, code := c.run("", "rename", "personal", "home"); code != 0 {
		t.Fatalf("%d %s", code, errOut)
	}
	if _, errOut, _ := c.run("", "--dry-run", "next", "--continue"); !strings.Contains(errOut, "→ home → side") {
		t.Fatalf("the renamed account's note: %s", errOut)
	}
	time.Sleep(10 * time.Millisecond) // the newer note
	switchTo("work", filepath.Join(transcript.ProjectDir(c.p.ProfileDir("work"), cwd), "s2.jsonl"), "home")
	if _, errOut, code := c.run("", "rename", "home", "personal"); code != 0 {
		t.Fatalf("%d %s", code, errOut)
	}
	if _, errOut, _ := c.run("", "--dry-run", "next", "--continue"); !strings.Contains(errOut, "→ work → personal") {
		t.Fatalf("a note naming the renamed account: %s", errOut)
	}
}

// The relay's hop follows its session's note as it says: to the account
// /ccam:switch named or the one chosen at a limit (Chosen), which nobody's
// own `ccam next --continue` follows. With the note or that account gone it
// ends with 1 and starts nothing.
func TestNextHopFollowsItsNote(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 40, "side": 30})
	trueClaude(t, c)
	t.Cleanup(func() { launcher.Unattended = false })
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	writeConversation(t, conv, time.Now())
	note := handoff.Ticket{Account: "personal", SessionID: "s1", Transcript: conv, Cwd: cwd, At: time.Now(), Reason: handoff.Limit, Limit: "session", Chosen: "side"}
	handoff.Save(c.p.State, note)
	path := handoff.Path(c.p.State, "personal", "s1")
	resolved, _ := filepath.EvalSymlinks(conv)
	// After a limit, the conversation goes on with ccam's prompt (D3).
	prompt := launcher.ShellQuote(handoff.ContinuePrompt)
	if runtime.GOOS == "windows" {
		prompt = `"` + handoff.ContinuePrompt + `"` // the dry run shows the Windows command line
	}
	out, errOut, code := c.run("", "--dry-run", "next", "--continue", "--note", path, "--", "--add-dir", "/tmp", "--permission-mode", "plan")
	if code != 0 || !strings.Contains(errOut, "→ personal → side") ||
		!strings.Contains(out, "--resume "+resolved+" --fork-session "+prompt+" --add-dir /tmp --permission-mode plan") {
		t.Fatalf("%d %s\n%s", code, errOut, out)
	}
	// After a subagent's limit, the prompt says what stopped with it.
	note.Subagent = "a1"
	handoff.Save(c.p.State, note)
	prompt = launcher.ShellQuote(handoff.SubagentPrompt)
	if runtime.GOOS == "windows" {
		prompt = `"` + handoff.SubagentPrompt + `"`
	}
	if out, errOut, code = c.run("", "--dry-run", "next", "--continue", "--note", path); code != 0 || !strings.Contains(out, "--fork-session "+prompt) {
		t.Fatalf("a subagent's limit: %d %s\n%s", code, errOut, out)
	}
	note.Subagent = ""
	handoff.Save(c.p.State, note)
	// By hand, Chosen is not followed: the next account, or by usage
	// without the one the limit stopped. Nor is a prompt sent.
	if out, errOut, _ := c.run("", "--dry-run", "next", "--continue"); !strings.Contains(errOut, "→ personal → work") || strings.Contains(out, "moved this conversation") {
		t.Fatalf("by hand: %s\n%s", errOut, out)
	}
	usage.SaveLimit(filepath.Join(c.p.State, "usage"), "side", "five_hour", time.Now().Add(time.Hour), time.Now())
	if _, errOut, _ := c.run("", "--dry-run", "next", "--auto", "--continue"); !strings.Contains(errOut, "→ personal → work") {
		t.Fatalf("by hand with --auto: %s", errOut)
	}
	gone := func(what, want string, args ...string) {
		t.Helper()
		out, errOut, code := c.run("", append([]string{"--dry-run", "next", "--continue", "--note"}, args...)...)
		if code != 1 || out != "" || !strings.Contains(errOut, want) {
			t.Fatalf("%s: %d %q %s", what, code, out, errOut)
		}
	}
	gone("another file", "this session's note is gone", path+"x")
	other := filepath.Join(t.TempDir(), "handoff", "personal", "s1.json")
	os.MkdirAll(filepath.Dir(other), 0o700)
	data, _ := os.ReadFile(path)
	os.WriteFile(other, data, 0o600)
	gone("a file outside the notes", "this session's note is gone", other)
	note.Chosen = "nobody"
	handoff.Save(c.p.State, note)
	gone("an account that is gone", "the note names nobody, which is no longer an account; run `ccam next --continue` to choose one", path)
	if _, err := os.Stat(path); err != nil {
		t.Fatal("the hop took the note although it went nowhere")
	}
	// The person asked to move (/ccam:switch): no prompt either.
	handoff.Save(c.p.State, handoff.Ticket{Account: "personal", SessionID: "s1", Transcript: conv, Cwd: cwd, At: time.Now(), Reason: handoff.Switch, Target: "side"})
	out, errOut, code = c.run("", "--dry-run", "next", "--continue", "--note", path, "--", "--add-dir", "/tmp")
	if code != 0 || !strings.Contains(out, "--resume "+resolved+" --fork-session --add-dir /tmp") {
		t.Fatalf("a switch: %d %s\n%s", code, errOut, out)
	}
}

// The prompt goes right after the conversation, before the carried
// arguments: after a flag of several values Claude Code would read it as
// one of them.
func TestContinueArgsCarriesThePrompt(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 40})
	e, cfg := c.env()
	cwd := t.TempDir()
	from, _ := cfg.Find("personal")
	to, _ := cfg.Find("work")
	shared := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("work"), cwd), "s1.jsonl")
	writeConversation(t, shared, time.Now())
	shared, _ = filepath.EvalSymlinks(shared)
	other := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	writeConversation(t, other, time.Now())
	other, _ = filepath.EvalSymlinks(other)
	for _, k := range []struct {
		path, prompt, want string
	}{
		{shared, "go on", "--resume s1 go on --add-dir /tmp"},
		{other, "go on", "--resume " + other + " --fork-session go on --add-dir /tmp"},
		{other, "", "--resume " + other + " --fork-session --add-dir /tmp"},
	} {
		args, _ := continueArgs(e, cfg, &conversation{path: k.path, owner: from, cwd: cwd}, to, k.prompt, []string{"--add-dir", "/tmp"})
		if got := strings.Join(args, " "); got != k.want {
			t.Errorf("got %q, want %q", got, k.want)
		}
	}
}

// A /ccam:switch hop without a named account picks by usage, and counts
// the conversation's own model: its arguments never name it.
func TestNextHopChoosesByTheConversationsModel(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 40, "side": 30})
	trueClaude(t, c)
	t.Cleanup(func() { launcher.Unattended = false })
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	writeConversation(t, conv, time.Now()) // its reply is Opus's
	handoff.Save(c.p.State, handoff.Ticket{Account: "personal", SessionID: "s1", Transcript: conv, Cwd: cwd, At: time.Now(), Reason: handoff.Switch})
	// side has the least 5-hour usage, but its Opus week is full.
	usage.SaveLimit(filepath.Join(c.p.State, "usage"), "side", "seven_day_opus", time.Now().Add(48*time.Hour), time.Now())
	path := handoff.Path(c.p.State, "personal", "s1")
	if _, errOut, code := c.run("", "--dry-run", "next", "--continue", "--note", path, "--auto", "--"); code != 0 || !strings.Contains(errOut, "→ personal → work") {
		t.Fatalf("%d %s", code, errOut)
	}
}
