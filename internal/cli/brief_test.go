package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/transcript"
)

// fakeGit makes gitStatus say out: the cli tests run inside the ccam
// repository, whose status is no test's business.
func fakeGit(t *testing.T, out string) {
	t.Helper()
	old := gitStatus
	gitStatus = func(string) string { return out }
	t.Cleanup(func() { gitStatus = old })
}

func TestWriteBrief(t *testing.T) {
	d := transcript.Digest{
		Title:      "Fix the parser",
		Prompts:    []string{"(1)", "(2)", "(3)", "(4)", "(5)", "(6)", "(7)"},
		Answer:     "(reply)",
		Unanswered: true,
		Summary:    "(summary)",
		Files:      []string{"/p/a.go", "/elsewhere/b.go"},
		Todos:      []transcript.Todo{{Text: "(one)", Status: "completed"}, {Text: "(two)", Status: "in_progress"}, {Text: "(three)", Status: "pending"}},
	}
	want := "# Summary of an earlier Claude Code conversation: \"Fix the parser\"\n\n" +
		briefIntro + "/data/s1.jsonl\n\n" +
		"## First request\n\n(1)\n\n" +
		"## Summary Claude Code wrote when it compacted the conversation\n\n(summary)\n\n" +
		"## Last requests, oldest first\n\n(2 earlier requests are left out.)\n\n1. (4)\n\n2. (5)\n\n3. (6)\n\n4. (7)\n\n" +
		"## Claude's last reply\n\n(reply)\n\n(This came before the last request above; Claude wrote no text after that request.)\n\n" +
		"## Todo list\n\n- [x] (one)\n- [ ] (two) (in progress)\n- [ ] (three)\n\n" +
		"## Files Claude changed\n\n- a.go\n- /elsewhere/b.go\n\n" +
		"## git status\n\n```\n## main\n M a.go\n```\n"
	if got := writeBrief(d, "/data/s1.jsonl", "/p", "## main\n M a.go"); got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
	// Empty parts are left out.
	want = "# Summary of an earlier Claude Code conversation\n\n" + briefIntro + "/data/s1.jsonl\n\n## Last requests, oldest first\n\n1. (only)\n"
	if got := writeBrief(transcript.Digest{Prompts: []string{"(only)"}}, "/data/s1.jsonl", "", ""); got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
	// One left out; the first request among the last ones is not repeated.
	got := writeBrief(transcript.Digest{Prompts: []string{"(1)", "(2)", "(3)", "(4)", "(5)", "(6)"}}, "/data/s1.jsonl", "", "")
	if !strings.Contains(got, "## First request\n\n(1)\n\n## Last requests, oldest first\n\n(1 earlier request is left out.)\n\n1. (3)") {
		t.Fatalf("%s", got)
	}
	if got := writeBrief(transcript.Digest{Prompts: []string{"(1)", "(2)", "(3)", "(4)"}}, "/data/s1.jsonl", "", ""); strings.Contains(got, "First request") || strings.Contains(got, "left out") {
		t.Fatalf("%s", got)
	}
}

// Every part is cut to its share, so a brief stays at a few thousand
// tokens however long the conversation was.
func TestWriteBriefStaysUnderMax(t *testing.T) {
	d := transcript.Digest{Title: strings.Repeat("t", 500), Answer: strings.Repeat("a", 50000), Unanswered: true, Summary: strings.Repeat("s", 50000)}
	for range 50 {
		d.Prompts = append(d.Prompts, strings.Repeat("p", 5000))
	}
	for i := range 100 {
		d.Files = append(d.Files, "/"+strings.Repeat(fmt.Sprintf("f%03d", i), 75))
		d.Todos = append(d.Todos, transcript.Todo{Text: strings.Repeat("x", 500), Status: "in_progress"})
	}
	git := strings.TrimSuffix(strings.Repeat(strings.Repeat("g", 200)+"\n", 100), "\n")
	text := writeBrief(d, "/"+strings.Repeat("p", 199), "/d", git)
	if n := utf8.RuneCountInString(text); n > briefMax {
		t.Fatalf("%d characters", n)
	}
	for _, s := range []string{"(45 earlier requests are left out.)", "(90 unfinished items are left out.)", "(90 other files are left out.)", "(97 more lines are left out.)", "characters left out ..."} {
		if !strings.Contains(text, s) {
			t.Errorf("no %q", s)
		}
	}
	// A git status line longer than all it may show is cut.
	if text := writeBrief(transcript.Digest{Prompts: []string{"(1)"}}, "/s.jsonl", "", strings.Repeat("g", 900)+"\n?? x"); !strings.Contains(text, "```\n"+strings.Repeat("g", gitMax-3)+"...\n```\n\n(1 more line is left out.)") {
		t.Fatalf("%s", text)
	}
}

// A long todo list keeps where the work stands. Claude keeps its plan in
// order (finished, in progress, pending), so the old last items alone
// were all pending. Finished items are left out first, the oldest first,
// then pending ones from the end of the list, items in progress last; the
// rest keep their order.
func TestWriteBriefTodos(t *testing.T) {
	for _, c := range []struct {
		list string // x finished, > in progress, . pending; item i says (i)
		note string
		kept []int
	}{
		{"xx>............", "(5 items are left out: 2 finished, 3 unfinished.)", []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
		{"x>..........", "(2 items are left out: 1 finished, 1 unfinished.)", []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}},
		{"xxxxxx>........", "(5 finished items are left out.)", []int{6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
		{"x>.........", "(1 finished item is left out.)", []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}},
		{".xxxxxxxxxxx>..", "(5 finished items are left out.)", []int{1, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
		{"...........>", "(2 unfinished items are left out.)", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12}},
		{"..........>", "(1 unfinished item is left out.)", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11}},
		{">>>>>>>>>>>>", "(2 unfinished items are left out.)", []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
		{"x>........", "", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	} {
		var d transcript.Digest
		status := map[rune]string{'x': "completed", '>': "in_progress", '.': "pending"}
		for i, r := range c.list {
			d.Todos = append(d.Todos, transcript.Todo{Text: fmt.Sprintf("(%d)", i+1), Status: status[r]})
		}
		want := "## Todo list\n\n"
		if c.note != "" {
			want += c.note + "\n"
		}
		var lines []string
		for _, i := range c.kept {
			switch c.list[i-1] {
			case 'x':
				lines = append(lines, fmt.Sprintf("- [x] (%d)", i))
			case '>':
				lines = append(lines, fmt.Sprintf("- [ ] (%d) (in progress)", i))
			default:
				lines = append(lines, fmt.Sprintf("- [ ] (%d)", i))
			}
		}
		want += strings.Join(lines, "\n") + "\n"
		if got := writeBrief(d, "/s.jsonl", "", ""); !strings.HasSuffix(got, "\n\n"+want) {
			t.Errorf("%s: got\n%s\nwant it to end in\n%s", c.list, got, want)
		}
	}
}

func TestClip(t *testing.T) {
	for _, c := range []struct{ got, want string }{
		{clip("short", 3, 2), "short"},
		{clip(strings.Repeat("x", 10), 2, 3), "xx\n[... 5 characters left out ...]\nxxx"},
		{clip("αβγδε", 2, 1), "αβ\n[... 2 characters left out ...]\nε"}, // characters, not bytes
		{cut("abcdefgh", 6), "abc..."},
		{cut("abcdef", 6), "abcdef"},
		{cut("αβγδεζη", 6), "αβγ..."},
		{cutLeft("/a/b/c/d.go", 8), ".../d.go"},
		{cutLeft("/a.go", 8), "/a.go"},
		{leftOut(1, "more line is", "more lines are"), "(1 more line is left out.)"},
		{leftOut(2, "more line is", "more lines are"), "(2 more lines are left out.)"},
	} {
		if c.got != c.want {
			t.Errorf("%q, want %q", c.got, c.want)
		}
	}
}

// The first prompt is one line (a batch claude refuses line breaks) that
// starts with no / or !, which Claude Code would read as a command. It
// @-mentions the summary, which Claude Code attaches. A path with a " (which
// ends the mention) or a # (which Claude Code reads as a line range) is
// mentioned from the working folder, which Claude Code resolves it from,
// when that path has neither: a # in the account's name is in a folder both
// share. Else the file is named for Claude to read.
func TestBriefPrompt(t *testing.T) {
	const rest = ", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here."
	cases := []struct{ file, dir, ask string }{
		{"/s/brief/s1.md", "/home/u/app", `Read @"/s/brief/s1.md"`},
		{"/home/my home/state/brief/s1.md", "/home/my home/app", `Read @"/home/my home/state/brief/s1.md"`},
		{"/home/c#dev/.local/state/ccam/brief/s1.md", "/home/c#dev/code/app", `Read @"../../.local/state/ccam/brief/s1.md"`},
		{"/home/c#dev/.local/state/ccam/brief/s1.md", "/home/c#dev", `Read @".local/state/ccam/brief/s1.md"`},
		{`/home/q"uote/state/brief/s1.md`, `/home/q"uote/app`, `Read @"../state/brief/s1.md"`},
		// The # is in a folder the working folder is not in, or there is no
		// working folder.
		{"/home/u/c#state/ccam/brief/s1.md", "/home/u/app", "Read the file /home/u/c#state/ccam/brief/s1.md"},
		{"/home/c#dev/.local/state/ccam/brief/s1.md", "/srv/app", "Read the file /home/c#dev/.local/state/ccam/brief/s1.md"},
		{`/home/q"uote/state/brief/s1.md`, "", `Read the file /home/q"uote/state/brief/s1.md`},
	}
	if runtime.GOOS == "windows" {
		cases = []struct{ file, dir, ask string }{
			{`C:\Users\me\.local\state\ccam\brief\s1.md`, `C:\Users\me\app`, `Read @"C:\Users\me\.local\state\ccam\brief\s1.md"`},
			{`C:\Users\C#dev\.local\state\ccam\brief\s1.md`, `C:\Users\C#dev\code\app`, `Read @"..\..\.local\state\ccam\brief\s1.md"`},
			{`C:\Users\C#dev\.local\state\ccam\brief\s1.md`, `C:\srv\app`, `Read the file C:\Users\C#dev\.local\state\ccam\brief\s1.md`},
			// No path leads from another drive.
			{`C:\Users\C#dev\.local\state\ccam\brief\s1.md`, `D:\app`, `Read the file C:\Users\C#dev\.local\state\ccam\brief\s1.md`},
		}
	}
	for _, c := range cases {
		wait := briefPrompt(c.file, c.dir, "")
		if want := c.ask + rest + " Then reply with one or two sentences on where the work stands, and wait for my next message."; wait != want {
			t.Errorf("%q, want %q", wait, want)
		}
		then := briefPrompt(c.file, c.dir, "run the tests")
		if want := c.ask + rest + " Then: run the tests"; then != want {
			t.Errorf("%q, want %q", then, want)
		}
		for _, p := range []string{wait, then} {
			if strings.Contains(p, "\n") || strings.HasPrefix(p, "/") || strings.HasPrefix(p, "!") {
				t.Errorf("%q", p)
			}
		}
	}
}

func TestGitStatus(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("no git")
	}
	dir := t.TempDir()
	t.Setenv("GIT_CEILING_DIRECTORIES", filepath.Dir(dir))
	if s := gitStatus(dir); s != "" {
		t.Fatalf("outside a work tree: %q", s)
	}
	if out, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("%v %s", err, out)
	}
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0o600)
	if s := gitStatus(dir); !strings.HasPrefix(s, "## ") || !strings.Contains(s, "?? a.txt") || strings.HasSuffix(s, "\n") {
		t.Fatalf("%q", s)
	}
	if s := gitStatus(filepath.Join(dir, "gone")); s != "" {
		t.Fatalf("a folder that is gone: %q", s)
	}
}

// briefArgs saves the summary in the state folder (0600), removes the ones
// older than a week, and starts claude with one prompt that mentions it;
// the person's prompt ends that prompt.
func TestBriefArgsSavesTheSummary(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 50, "work": 10})
	fakeGit(t, "## main\n M a.go")
	cwd := t.TempDir()
	real, _ := filepath.EvalSymlinks(cwd)
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("work"), cwd), "s1.jsonl")
	writeConversation(t, conv, time.Now(), fmt.Sprintf(`{"isSidechain":false,"type":"assistant","message":{"id":"m2","model":"claude-opus-5-5","role":"assistant","content":[{"type":"tool_use","id":"t1","name":"Edit","input":{"file_path":%q}}]}}`, filepath.Join(real, "a.go")))
	dir := filepath.Join(c.p.State, "brief")
	os.MkdirAll(dir, 0o700)
	old, recent := filepath.Join(dir, "old.md"), filepath.Join(dir, "recent.md")
	for f, age := range map[string]time.Duration{old: 8 * 24 * time.Hour, recent: 24 * time.Hour} {
		os.WriteFile(f, []byte("(brief)"), 0o600)
		at := time.Now().Add(-age)
		os.Chtimes(f, at, at)
	}
	e := &env{paths: c.p, stdout: io.Discard, stderr: io.Discard}
	args, note, err := briefArgs(e, &conversation{path: conv, cwd: cwd}, []string{"--model", "x", "run it"})
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(dir, "s1.md")
	if want := []string{briefPrompt(file, real, "run it"), "--model", "x"}; !reflect.DeepEqual(args, want) {
		t.Fatalf("%q", args)
	}
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, s := range []string{"## Last requests, oldest first\n\n1. Fix the failing tests\n", "## Claude's last reply\n\nDone.\n", "## Files Claude changed\n\n- a.go\n", "```\n## main\n M a.go\n```\n", conv} {
		if !strings.Contains(text, s) {
			t.Errorf("no %q in\n%s", s, text)
		}
	}
	if !strings.HasPrefix(text, "# Summary of an earlier Claude Code conversation\n") {
		t.Fatalf("%s", text)
	}
	if info, _ := os.Stat(file); runtime.GOOS != "windows" && info.Mode().Perm() != 0o600 {
		t.Fatalf("mode %v", info.Mode())
	}
	if _, err := os.Stat(old); err == nil {
		t.Fatal("a summary older than a week stays")
	}
	if _, err := os.Stat(recent); err != nil {
		t.Fatal("a summary of yesterday is gone")
	}
	want := "starting a new conversation with ccam's summary of the last conversation instead of resuming it (about 15 tokens)\n" +
		fmt.Sprintf("  summary (%s characters): %s", tokens(utf8.RuneCountInString(text)), file)
	if note != want {
		t.Fatalf("%q\nwant %q", note, want)
	}
}

// A summary under a folder with a # is mentioned by its path from the
// working folder with its links resolved, as Claude Code resolves its own
// working folder: from there that path names the summary. Through the link
// the path would have the #.
func TestBriefArgsMentionsFromTheFolder(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 50, "work": 10})
	fakeGit(t, "")
	root, _ := filepath.EvalSymlinks(t.TempDir())
	home := filepath.Join(root, "c#dev")
	app := filepath.Join(home, "code", "app")
	os.MkdirAll(app, 0o700)
	link := filepath.Join(root, "app")
	if err := os.Symlink(app, link); err != nil {
		t.Skip(err) // Windows without the right to make links
	}
	p := c.p
	p.State = filepath.Join(home, ".local", "state", "ccam")
	conv := filepath.Join(transcript.ProjectDir(p.ProfileDir("work"), link), "s1.jsonl")
	writeConversation(t, conv, time.Now())
	e := &env{paths: p, stdout: io.Discard, stderr: io.Discard}
	args, _, err := briefArgs(e, &conversation{path: conv, cwd: link}, nil)
	if err != nil {
		t.Fatal(err)
	}
	rel := filepath.Join("..", "..", ".local", "state", "ccam", "brief", "s1.md")
	if want := `Read @"` + rel + `", the summary ccam wrote`; !strings.HasPrefix(args[0], want) {
		t.Fatalf("%q, want it to start with %q", args[0], want)
	}
	if _, err := os.Stat(filepath.Join(app, rel)); err != nil {
		t.Fatal(err)
	}
}

// What --brief cannot do ends ccam before claude starts.
func TestBriefArgsFails(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 50, "work": 10})
	fakeGit(t, "")
	e := &env{paths: c.p, stdout: io.Discard, stderr: io.Discard}
	cwd := t.TempDir()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("work"), cwd), "s1.jsonl")
	if _, _, err := briefArgs(e, &conversation{path: conv, cwd: cwd}, nil); err == nil ||
		!strings.HasPrefix(err.Error(), "--brief could not read the conversation (") || !strings.HasSuffix(err.Error(), "; `ccam next --continue` without --brief resumes it whole") {
		t.Fatalf("%v", err)
	}
	writeConversation(t, conv, time.Now())
	os.MkdirAll(c.p.State, 0o700)
	os.WriteFile(filepath.Join(c.p.State, "brief"), nil, 0o600) // no folder can be made there
	if _, _, err := briefArgs(e, &conversation{path: conv, cwd: cwd}, nil); err == nil || !strings.HasPrefix(err.Error(), "--brief could not save the summary (") {
		t.Fatalf("%v", err)
	}
}

// A dry run prints the summary and the launch, and writes nothing.
func TestNextContinueBrief(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 50, "work": 10})
	trueClaude(t, c)
	fakeGit(t, "## main\n M x.go")
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("work"), cwd), "s1.jsonl")
	writeConversation(t, conv, time.Now().Add(-10*time.Minute))
	out, errOut, code := c.run("", "--dry-run", "next", "--continue", "--brief", "--", "--model", "sonnet", "run the tests")
	if code != 0 || !strings.Contains(errOut, "→ work → personal\n  starting a new conversation with ccam's summary of the last conversation instead of resuming it (about 15 tokens)\n  summary (") {
		t.Fatalf("%d %s", code, errOut)
	}
	// The Windows dry run is a PowerShell line, which quotes otherwise.
	for _, s := range []string{"## Last requests, oldest first\n\n1. Fix the failing tests\n", "## Claude's last reply\n\nDone.\n", "```\n## main\n M x.go\n```\n\n",
		"the summary ccam wrote of my earlier Claude Code conversation", "Then: run the tests", "s1.md", "--model sonnet"} {
		if !strings.Contains(out, s) {
			t.Errorf("no %q in\n%s", s, out)
		}
	}
	if strings.Contains(out, "--resume") {
		t.Fatalf("%s", out)
	}
	if entries, _ := os.ReadDir(filepath.Join(c.p.State, "brief")); len(entries) > 0 {
		t.Fatal("a dry run wrote a summary")
	}
}

func TestNextBriefNeedsContinue(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 50, "work": 10})
	if _, errOut, code := c.run("", "--dry-run", "next", "--brief"); code == 0 || !strings.Contains(errOut, "--brief works with --continue") {
		t.Fatalf("%d %s", code, errOut)
	}
}

// A conversation of rows the digest does not know is refused, not resumed
// whole: the person chose --brief for its cost. The note stays for the
// next try.
func TestBriefFindsNothingToSummarize(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 50, "work": 10})
	// A real run, in which claude cannot be found: nothing replaces the test.
	cfg, _ := config.Load(c.p.ConfigFile())
	cfg.ClaudeBin = "ccam-test-no-such-claude"
	cfg.Save(c.p.ConfigFile())
	fakeGit(t, "")
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("work"), cwd), "s1.jsonl")
	os.MkdirAll(filepath.Dir(conv), 0o700)
	os.WriteFile(conv, []byte(`{"isSidechain":false,"type":"user","message":{"role":"user","content":"<command-name>/mcp</command-name>\n<command-message>mcp</command-message>\n<command-args></command-args>"}}`+"\n"+
		`{"isSidechain":false,"type":"assistant","message":{"id":"m1","model":"claude-opus-5-5","role":"assistant","content":[{"type":"tool_use","id":"t1","name":"Bash","input":{"command":"ls"}}]}}`+"\n"), 0o600)
	if _, _, ok := transcript.Latest(filepath.Dir(conv), nil); !ok {
		t.Fatal("Latest does not find it")
	}
	handoff.Save(c.p.State, handoff.Ticket{Account: "work", SessionID: "s1", Transcript: conv, Cwd: cwd, At: time.Now(), Reason: handoff.Switch, Target: "personal"})
	_, errOut, code := c.run("", "next", "--continue", "--brief")
	if code == 0 || !strings.Contains(errOut, "--brief found nothing to summarize in the conversation (a newer Claude Code may save it differently); `ccam next --continue` without --brief resumes it whole") {
		t.Fatalf("%d %s", code, errOut)
	}
	if _, err := os.Stat(handoff.Path(c.p.State, "work", "s1")); err != nil {
		t.Fatal("the note was used up")
	}
}
