package transcript

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSlug(t *testing.T) {
	long := "/home/me/" + strings.Repeat("very-long-folder-name/", 12) + "project"
	cases := map[string]string{
		"/home/jhyunwoo/workspace/ccam": "-home-jhyunwoo-workspace-ccam",
		"/a/b_c.d e":                    "-a-b-c-d-e",
		"/x/한":                          "-x--",
		"/x/😀":                          "-x---",
		`C:\Users\me\code`:              "C--Users-me-code",
		// Claude Code 2.1.304 cuts at 200 and adds the path's hash.
		long: "-home-me-very-long-folder-name-very-long-folder-name-very-long-folder-name-very-long-folder-name-very-long-folder-name-very-long-folder-name-very-long-folder-name-very-long-folder-name-very-long-folde-g15bg1",
	}
	for in, want := range cases {
		if got := Slug(in); got != want {
			t.Errorf("Slug(%q) = %q, want %q", in, got, want)
		}
	}
}

func write(t *testing.T, path, data string, mod time.Time) {
	t.Helper()
	os.MkdirAll(filepath.Dir(path), 0o700)
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	os.Chtimes(path, mod, mod)
}

func TestLatestSizeAndTitle(t *testing.T) {
	cfg := t.TempDir()
	cwd := t.TempDir()
	dir := ProjectDir(cfg, cwd)
	now := time.Now()
	write(t, filepath.Join(dir, "old.jsonl"), "{}\n", now.Add(-time.Hour))
	conv := strings.Join([]string{
		`{"type":"user","message":{"role":"user","content":"hi"}}`,
		`{"type":"assistant","message":{"usage":{"input_tokens":5,"cache_read_input_tokens":100,"cache_creation_input_tokens":20,"output_tokens":7}}}`,
		`{"type":"ai-title","aiTitle":"Fix the tests"}`,
		`{"type":"assistant","isSidechain":true,"message":{"usage":{"input_tokens":999999}}}`,
		`{"type":"system","subtype":"turn_duration"}`,
	}, "\n") + "\n"
	write(t, filepath.Join(dir, "abc-123.jsonl"), conv, now)
	os.MkdirAll(filepath.Join(dir, "abc-123", "subagents"), 0o700) // folders are not conversations
	path, _, ok := Latest(dir)
	if !ok || ID(path) != "abc-123" {
		t.Fatalf("%q %v", path, ok)
	}
	if n := Size(path); n != 132 {
		t.Fatalf("size %d (the subagent's reply is not the conversation's)", n)
	}
	if got := Title(path); got != "Fix the tests" {
		t.Fatal(got)
	}
	write(t, path, conv+`{"type":"custom-title","customTitle":"release prep"}`+"\n", now)
	if got := Title(path); got != "release prep" {
		t.Fatal(got)
	}
	if _, _, ok := Latest(filepath.Join(cfg, "projects", "nothing")); ok {
		t.Fatal("no folder")
	}
}

// The limit fixtures hold rows of the shape Claude Code wrote at real
// usage limits on a Team seat: the refused request's row (model
// "<synthetic>", apiError "usage_limit_reached", rate_limit_info) with its
// keys and rate-limit values as observed, 2.1.304 (two refusals on
// 2026-10-20) and 2.1.306 (2026-10-21T14:51:46Z). Ids, paths and texts are
// placeholders; the prompt and reply rows around them are made up.
func fixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestLastLimitReadsTheRefusedRequest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "c.jsonl")
	write(t, path, fixture(t, "limit-2.1.306-five_hour.jsonl"), time.Now())
	l, ok := LastLimit(path, time.Time{})
	at := time.Date(2026, 10, 21, 14, 51, 46, 515e6, time.UTC)
	if !ok || l.Window != "five_hour" || !l.ResetsAt.Equal(time.Unix(1792607400, 0)) || !l.At.Equal(at) || l.Model != "claude-opus-5-5" {
		t.Fatalf("%+v %v", l, ok)
	}
	// Only a limit hit at since or later counts.
	if _, ok := LastLimit(path, at.Add(time.Second)); ok {
		t.Fatal("a limit before since")
	}
	if _, ok := LastLimit(path, at.Add(-time.Minute)); !ok {
		t.Fatal("a limit after since")
	}
	// 2.1.304 wrote the same row; the newer of two refusals counts.
	write(t, path, fixture(t, "limit-2.1.304-team.jsonl"), time.Now())
	if l, ok := LastLimit(path, time.Time{}); !ok || !l.At.Equal(time.Date(2026, 10, 20, 13, 32, 33, 105e6, time.UTC)) || l.Model != "claude-opus-5-5" {
		t.Fatalf("2.1.304: %+v %v", l, ok)
	}
}

// Only the main thread's refusal at a usage limit counts: not a subagent's,
// not another API error, not a row of another shape.
func TestLastLimitSkipsOtherRows(t *testing.T) {
	base := fixture(t, "limit-2.1.306-five_hour.jsonl")
	lines := strings.Split(strings.TrimSpace(base), "\n")
	row := lines[len(lines)-1]
	newer := strings.Replace(row, "2026-10-21T14:51:46.515Z", "2026-10-21T14:59:00.000Z", 1)
	cases := map[string]string{
		"subagent":       strings.Replace(newer, `"isSidechain":false`, `"isSidechain":true`, 1),
		"credits needed": strings.Replace(newer, `"usage_limit_reached"`, `"model_requires_usage_credits"`, 1),
		"warning only":   strings.Replace(newer, `"status":"rejected","resetsAt"`, `"status":"allowed_warning","resetsAt"`, 1),
		"params a text":  strings.Replace(newer, `"apiErrorParams":{"rate_limit_info":`, `"apiErrorParams":"x","y":{"rate_limit_info":`, 1),
		"not an error":   strings.Replace(newer, `"isApiErrorMessage":true`, `"isApiErrorMessage":false`, 1),
	}
	dir := t.TempDir()
	for name, other := range cases {
		if other == newer {
			t.Fatalf("%s: the fixture changed", name)
		}
		path := filepath.Join(dir, "c.jsonl")
		write(t, path, base+other+"\n", time.Now())
		if l, ok := LastLimit(path, time.Time{}); !ok || !l.At.Equal(time.Date(2026, 10, 21, 14, 51, 46, 515e6, time.UTC)) {
			t.Errorf("%s: %+v %v", name, l, ok)
		}
		write(t, path, strings.Join(lines[:len(lines)-1], "\n")+"\n"+other+"\n", time.Now())
		if l, ok := LastLimit(path, time.Time{}); ok {
			t.Errorf("%s alone: %+v", name, l)
		}
	}
}

// Claude Code may still be writing the conversation: a half-written last
// line is skipped, and a limit row that is not complete yet is not found.
func TestLastLimitWhileWriting(t *testing.T) {
	base := fixture(t, "limit-2.1.306-five_hour.jsonl")
	path := filepath.Join(t.TempDir(), "c.jsonl")
	write(t, path, base+`{"parentUuid":"l1","isSidechain":false,"type":"system","subtype":"turn_dur`, time.Now())
	if _, ok := LastLimit(path, time.Time{}); !ok {
		t.Fatal("a half-written line after the limit")
	}
	cut := strings.TrimSpace(base)
	write(t, path, cut[:len(cut)-40], time.Now())
	if l, ok := LastLimit(path, time.Time{}); ok {
		t.Fatalf("a half-written limit row: %+v", l)
	}
}

// Only the end of a long conversation is read.
func TestLastLimitReadsTheEnd(t *testing.T) {
	filler := strings.Repeat(`{"type":"system","subtype":"turn_duration","content":"`+strings.Repeat("x", 900)+`"}`+"\n", 1300)
	path := filepath.Join(t.TempDir(), "c.jsonl")
	write(t, path, filler+fixture(t, "limit-2.1.306-five_hour.jsonl"), time.Now())
	if _, ok := LastLimit(path, time.Time{}); !ok {
		t.Fatal("a limit at the end of a conversation over 1 MiB")
	}
	write(t, path, fixture(t, "limit-2.1.306-five_hour.jsonl")+filler, time.Now())
	if _, ok := LastLimit(path, time.Time{}); ok {
		t.Fatal("a limit over 1 MiB before the end")
	}
}

// After a limit, the size is the last real reply's: the refused request's
// row reads and writes nothing.
func TestSizeAndModelSkipTheLimitRow(t *testing.T) {
	path := filepath.Join(t.TempDir(), "c.jsonl")
	write(t, path, fixture(t, "limit-2.1.306-five_hour.jsonl"), time.Now())
	if n := Size(path); n != 3+98000+1200+450 {
		t.Fatalf("size %d", n)
	}
	if m := Model(path); m != "claude-opus-5-5" {
		t.Fatalf("model %q", m)
	}
	lines := strings.Split(strings.TrimSpace(fixture(t, "limit-2.1.306-five_hour.jsonl")), "\n")
	write(t, path, lines[len(lines)-1]+"\n", time.Now())
	if n, m := Size(path), Model(path); n != 0 || m != "" {
		t.Fatalf("only the limit row: size %d, model %q", n, m)
	}
}

// The permission mode is the newer of the permission-mode rows and the
// mode a person's prompt was sent in.
func TestPermissionMode(t *testing.T) {
	path := filepath.Join(t.TempDir(), "c.jsonl")
	base := fixture(t, "limit-2.1.306-five_hour.jsonl")
	write(t, path, base, time.Now())
	if m := PermissionMode(path); m != "acceptEdits" {
		t.Fatalf("the prompt after the permission-mode row: %q", m)
	}
	write(t, path, base+`{"type":"permission-mode","permissionMode":"plan","sessionId":"s"}`+"\n", time.Now())
	if m := PermissionMode(path); m != "plan" {
		t.Fatalf("a newer permission-mode row: %q", m)
	}
	// Prompts Claude Code sends itself do not say what the person chose, nor
	// do a local command and the summary written after compacting (made-up
	// rows of those shapes, with a mode added).
	others := `{"type":"user","message":{"role":"user","content":"(done)"},"permissionMode":"bypassPermissions","origin":{"kind":"task-notification"}}` + "\n" +
		`{"type":"user","isMeta":true,"message":{"role":"user","content":"(peer)"},"permissionMode":"bypassPermissions","origin":{"kind":"peer"}}` + "\n" +
		`{"type":"user","isSidechain":true,"message":{"role":"user","content":"(task)"},"permissionMode":"bypassPermissions"}` + "\n" +
		`{"type":"user","message":{"role":"user","content":"<command-name>/usage</command-name>\n<command-args></command-args>"},"permissionMode":"bypassPermissions"}` + "\n" +
		`{"type":"user","isCompactSummary":true,"isVisibleInTranscriptOnly":true,"message":{"role":"user","content":"(summary)"},"permissionMode":"bypassPermissions"}` + "\n"
	write(t, path, base+others, time.Now())
	if m := PermissionMode(path); m != "acceptEdits" {
		t.Fatalf("prompts not from the person: %q", m)
	}
	// 2.1.304's conversation has no permission-mode row; its prompts say.
	write(t, path, fixture(t, "limit-2.1.304-team.jsonl"), time.Now())
	if m := PermissionMode(path); m != "default" {
		t.Fatalf("no permission-mode row: %q", m)
	}
	if m := PermissionMode(filepath.Join(t.TempDir(), "none.jsonl")); m != "" {
		t.Fatalf("no conversation: %q", m)
	}
}

// The limit counts only when it ended the conversation's last turn: once
// a person has sent a prompt after it, the turn that prompt began is the
// one a hook asks about, and the older refusal (possibly another account's,
// in a shared history) is not its limit.
func TestLastLimitEndsTheLastTurn(t *testing.T) {
	base := fixture(t, "limit-2.1.306-five_hour.jsonl")
	path := filepath.Join(t.TempDir(), "c.jsonl")
	prompt := `{"parentUuid":"l1","isSidechain":false,"type":"user","message":{"role":"user","content":"(again)"},"uuid":"u2","timestamp":"2026-10-21T14:52:10.000Z","permissionMode":"default","origin":{"kind":"human"}}` + "\n"
	throttle := `{"parentUuid":"u2","isSidechain":false,"type":"assistant","message":{"model":"<synthetic>","role":"assistant","content":[{"type":"text","text":"API Error: Server is temporarily limiting requests (not your usage limit)"}]},"isApiErrorMessage":true,"error":"rate_limit","timestamp":"2026-10-21T14:52:40.000Z"}` + "\n"
	write(t, path, base+prompt+throttle, time.Now())
	if l, ok := LastLimit(path, time.Time{}); ok {
		t.Fatalf("a refusal before the last prompt: %+v", l)
	}
	// Rows Claude Code adds after a limit do not begin a turn of the person's,
	// and neither do a local command (/usage), which never reaches the model,
	// and the summary Claude Code writes after compacting.
	after := `{"type":"user","message":{"role":"user","content":"(done)"},"origin":{"kind":"task-notification"}}` + "\n" +
		`{"type":"user","isMeta":true,"message":{"role":"user","content":"(caveat)"}}` + "\n" +
		`{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"t1","content":"ok"}]}}` + "\n" +
		`{"type":"system","subtype":"turn_duration"}` + "\n" + `{"type":"file-history-snapshot","snapshot":{}}` + "\n" +
		usageCommand + compactSummary
	write(t, path, base+after, time.Now())
	l, ok := LastLimit(path, time.Time{})
	if !ok || l.Text != "You've hit your session limit · resets 3:30am (Asia/Seoul)" {
		t.Fatalf("rows after the limit: %+v %v", l, ok)
	}
}

// A conversation counts by its last prompt or reply. Sessions that only
// ran local commands (2.1.306 writes them as user rows, or system rows),
// empty files, and rows Claude Code writes at exit say nothing about it.
func TestLatestSkipsSessionsWithoutAConversation(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	ts := func(d time.Duration) string { return now.Add(d).UTC().Format(time.RFC3339) }
	earlier := `{"type":"user","message":{"role":"user","content":"earlier"},"timestamp":"` + ts(-time.Hour) + `"}` + "\n"
	for _, c := range []struct {
		name, body string
		counts     bool
	}{
		{"empty", "", false},
		{"mcp only", `{"type":"last-prompt","lastPrompt":"/mcp"}
{"parentUuid":null,"isSidechain":false,"type":"attachment","attachment":{"type":"hook_success"}}
{"parentUuid":"t1","isSidechain":false,"type":"user","message":{"role":"user","content":"<local-command-caveat>Caveat</local-command-caveat>"},"isMeta":true}
{"parentUuid":"c1","isSidechain":false,"type":"user","message":{"role":"user","content":"<command-name>/mcp</command-name>\n<command-args></command-args>"}}
{"parentUuid":"c2","isSidechain":false,"type":"user","message":{"role":"user","content":"<local-command-stdout>MCP dialog dismissed</local-command-stdout>"}}
{"type":"cost-state"}
`, false},
		{"status only", `{"type":"mode","mode":"normal"}
{"parentUuid":null,"type":"system","subtype":"local_command","content":"<command-name>/status</command-name>"}
`, false},
		{"a subagent's rows", `{"type":"assistant","isSidechain":true,"message":{"content":[]},"timestamp":"` + ts(0) + `"}` + "\n", false},
		{"prompt refused at a limit", `{"type":"user","message":{"role":"user","content":"Fix the tests"},"timestamp":"` + ts(-time.Minute) + `"}
{"type":"assistant","message":{"model":"<synthetic>","content":[{"type":"text","text":"You've hit your session limit"}]},"isApiErrorMessage":true,"timestamp":"` + ts(-time.Minute) + `"}
`, true},
		{"tool result", `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","content":"ok"}]},"timestamp":"` + ts(-time.Minute) + `"}` + "\n", true},
	} {
		dir := t.TempDir()
		old := filepath.Join(dir, "s1.jsonl")
		write(t, old, earlier, now.Add(-time.Hour))
		f := filepath.Join(dir, "s2.jsonl")
		write(t, f, c.body, now)
		path, _, ok := Latest(dir)
		if want := map[bool]string{true: f, false: old}[c.counts]; !ok || path != want {
			t.Errorf("%s: %s, want %s", c.name, filepath.Base(path), filepath.Base(want))
		}
	}
	// Rows at exit make s1's file the newest, but s2 went on later.
	dir := t.TempDir()
	write(t, filepath.Join(dir, "s1.jsonl"), earlier+`{"type":"cost-state"}`+"\n", now)
	write(t, filepath.Join(dir, "s2.jsonl"), `{"type":"assistant","message":{"content":[]},"timestamp":"`+ts(-10*time.Minute)+`"}`+"\n", now.Add(-10*time.Minute))
	if path, at, ok := Latest(dir); !ok || ID(path) != "s2" || !at.Equal(now.Add(-10*time.Minute)) {
		t.Fatalf("%s %v", path, at)
	}
	// File times kept to the second: s4's message is the later one although
	// its file looks older than s3's message.
	dir = t.TempDir()
	write(t, filepath.Join(dir, "s3.jsonl"), `{"type":"assistant","message":{"content":[]},"timestamp":"`+now.Add(900*time.Millisecond).UTC().Format(time.RFC3339Nano)+`"}`+"\n", now.Add(time.Second))
	write(t, filepath.Join(dir, "s4.jsonl"), `{"type":"assistant","message":{"content":[]},"timestamp":"`+now.Add(950*time.Millisecond).UTC().Format(time.RFC3339Nano)+`"}`+"\n", now)
	if path, _, _ := Latest(dir); ID(path) != "s4" {
		t.Fatalf("%s", path)
	}
}

// Made-up rows of shapes Claude Code (2.1.306-2.1.310) writes, with
// placeholder text: a local command (/usage: its caveat, the command and
// its output) and the summary it writes after compacting a conversation.
const (
	usageCommand = `{"parentUuid":"l1","isSidechain":false,"type":"user","message":{"role":"user","content":"<local-command-caveat>(caveat)</local-command-caveat>"},"isMeta":true}` + "\n" +
		`{"parentUuid":"c1","isSidechain":false,"type":"user","message":{"role":"user","content":"<command-name>/usage</command-name>\n            <command-message>usage</command-message>\n            <command-args></command-args>"}}` + "\n" +
		`{"parentUuid":"c2","isSidechain":false,"type":"user","message":{"role":"user","content":"<local-command-stdout>(dialog dismissed)</local-command-stdout>"}}` + "\n"
	compactSummary = `{"parentUuid":null,"isSidechain":false,"type":"user","isCompactSummary":true,"isVisibleInTranscriptOnly":true,"message":{"role":"user","content":"(summary)"},"uuid":"cs1"}` + "\n"
)

// A prompt is what a person sent; nothing Claude Code sends itself.
func TestPrompt(t *testing.T) {
	for row, want := range map[string]bool{
		`{"type":"user","message":{"role":"user","content":"(prompt)"},"origin":{"kind":"human"}}`:                                 true,
		`{"type":"user","message":{"role":"user","content":"(prompt)"}}`:                                                           true,
		`{"type":"user","message":{"role":"user","content":[{"type":"text","text":"(prompt)"}]}}`:                                  true,
		`{"type":"user","isCompactSummary":true,"isVisibleInTranscriptOnly":true,"message":{"role":"user","content":"(summary)"}}`: false,
		`{"type":"user","isMeta":true,"message":{"role":"user","content":"(continue)"},"origin":{"kind":"auto-continuation"}}`:     false,
		`{"type":"user","message":{"role":"user","content":"(done)"},"origin":{"kind":"task-notification"}}`:                       false,
		`{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"t1","content":"ok"}]}}`:           false,
		`{"type":"user","message":{"role":"user","content":"<command-name>/usage</command-name>\n<command-args></command-args>"}}`: false,
		`{"type":"user","message":{"role":"user","content":"<local-command-stdout>(dialog dismissed)</local-command-stdout>"}}`:    false,
		`{"type":"user","isSidechain":true,"message":{"role":"user","content":"(task)"}}`:                                          false,
		`{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"(reply)"}]}}`:                         false,
		// A slash command the person ran reaches the model, and so does a !
		// command (Claude Code 2.1.310 queries the model for it by default);
		// its output, a finished task without origin and Esc do not.
		`{"type":"user","message":{"role":"user","content":"<command-message>review</command-message>\n<command-name>/review</command-name>"},"origin":{"kind":"human"}}`: true,
		`{"type":"user","message":{"role":"user","content":[{"type":"image","source":{"type":"base64","data":"AAAA"}}]}}`:                                                 true,
		`{"type":"user","message":{"role":"user","content":"<bash-input>ls</bash-input>"}}`:                                                                               true,
		`{"type":"user","message":{"role":"user","content":"<bash-stdout>(files)</bash-stdout><bash-stderr></bash-stderr>"}}`:                                             false,
		`{"type":"user","message":{"role":"user","content":"<task-notification>(done)</task-notification>"}}`:                                                             false,
		`{"type":"user","message":{"role":"user","content":[{"type":"text","text":"[Request interrupted by user]"}]}}`:                                                    false,
		`{"type":"user","message":{"role":"user","content":" "},"origin":{"kind":"human"}}`:                                                                               false,
	} {
		if got := prompt([]byte(row)); got != want {
			t.Errorf("%s: %v", row, got)
		}
	}
}

// After a limit, only a prompt of the person's counts as moving on: the
// rows Claude Code adds by itself do not.
func TestPromptedAfter(t *testing.T) {
	base := fixture(t, "limit-2.1.306-five_hour.jsonl")
	at := time.Date(2026, 10, 21, 14, 51, 46, 515e6, time.UTC)
	path := filepath.Join(t.TempDir(), "c.jsonl")
	write(t, path, base, time.Now())
	if PromptedAfter(path, at) {
		t.Fatal("nothing after the limit")
	}
	// Made-up rows of the shapes Claude Code writes after a limit.
	for name, row := range map[string]string{
		"turn duration":     `{"type":"system","subtype":"turn_duration","durationMs":1}`,
		"Remote Control":    `{"type":"bridge-session","sessionId":"s","bridgeSessionId":"cse_test","lastSequenceNum":3}`,
		"cost":              `{"type":"cost-state","sessionId":"s"}`,
		"last prompt":       `{"type":"last-prompt","lastPrompt":"(prompt)","sessionId":"s"}`,
		"file history":      `{"type":"file-history-snapshot","messageId":"m","snapshot":{},"isSnapshotUpdate":false}`,
		"queued":            `{"type":"queue-operation","operation":"enqueue","content":"(prompt)","sessionId":"s"}`,
		"attachment":        `{"type":"attachment","attachment":{"type":"queued_command","prompt":"(prompt)","origin":{"kind":"human"}}}`,
		"task notification": `{"type":"user","message":{"role":"user","content":"(done)"},"origin":{"kind":"task-notification"}}`,
		"local command":     strings.TrimSpace(usageCommand),
		"compact summary":   strings.TrimSpace(compactSummary),
		"tool result":       `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"t1","content":"ok"}]}}`,
		"subagent":          `{"type":"user","isSidechain":true,"message":{"role":"user","content":"(task)"}}`,
		"! command output":  `{"type":"user","message":{"role":"user","content":"<bash-stdout>(output)</bash-stdout><bash-stderr></bash-stderr>"}}`,
	} {
		write(t, path, base+row+"\n", time.Now())
		if PromptedAfter(path, at) {
			t.Errorf("%s: a prompt", name)
		}
	}
	for name, row := range map[string]string{
		"typed":          `{"type":"user","message":{"role":"user","content":"(go on)"},"origin":{"kind":"human"}}`,
		"without origin": `{"type":"user","message":{"role":"user","content":"(go on)"}}`,
		// Claude Code queries the model for it, and counts it as the
		// person taking over its wait for the limit, as a typed prompt.
		"! command": `{"type":"user","message":{"role":"user","content":"<bash-input>git status</bash-input>"}}`,
	} {
		write(t, path, base+row+"\n", time.Now())
		if !PromptedAfter(path, at) {
			t.Errorf("%s: no prompt", name)
		}
	}
	// A prompt before the limit is the one the limit refused.
	lines := strings.SplitAfter(base, "\n")
	write(t, path, strings.Join(lines[:len(lines)-2], "")+`{"type":"user","message":{"role":"user","content":"(go on)"}}`+"\n"+lines[len(lines)-2], time.Now())
	if PromptedAfter(path, at) {
		t.Fatal("a prompt before the limit")
	}
	// A limit that is not in the end of the conversation: ccam cannot tell.
	// No prompt is in the part read either: the limit row has left it, or
	// the limit there is another one.
	filler := strings.Repeat(`{"type":"system","subtype":"turn_duration","content":"`+strings.Repeat("x", 900)+`"}`+"\n", 1300)
	write(t, path, base+filler, time.Now())
	if !PromptedAfter(path, at) {
		t.Error("a limit before the part read")
	}
	write(t, path, lines[len(lines)-2]+`{"type":"system","subtype":"turn_duration","durationMs":1}`+"\n", time.Now())
	if !PromptedAfter(path, at.Add(time.Second)) {
		t.Error("another limit")
	}
	if !PromptedAfter(filepath.Join(t.TempDir(), "none.jsonl"), at) {
		t.Error("no conversation")
	}
}

// Remote Control is on when the last bridge-session row names a bridge
// session: Claude Code writes one with an empty id when it is turned off.
func TestRemoteControl(t *testing.T) {
	// Made-up rows of the shapes Claude Code 2.1.308-2.1.310 writes.
	on := `{"type":"bridge-session","sessionId":"s","bridgeSessionId":"cse_test","lastSequenceNum":3}` + "\n"
	off := `{"type":"bridge-session","sessionId":"s","bridgeSessionId":"","lastSequenceNum":0}` + "\n"
	quoted := `{"type":"user","message":{"role":"user","content":"(about {\"type\":\"bridge-session\",\"bridgeSessionId\":\"cse_x\"})"}}` + "\n"
	// A tool call's input is an object: the word is no escaped text there.
	search := `{"isSidechain":false,"type":"assistant","message":{"id":"m9","model":"claude-opus-5-5","role":"assistant","content":[{"type":"tool_use","id":"t9","name":"Grep","input":{"pattern":"bridge-session"}}]}}` + "\n"
	filler := strings.Repeat(`{"type":"system","subtype":"turn_duration","content":"`+strings.Repeat("x", 900)+`"}`+"\n", 1300)
	base := fixture(t, "limit-2.1.306-five_hour.jsonl")
	path := filepath.Join(t.TempDir(), "c.jsonl")
	for _, c := range []struct {
		name, body string
		want       bool
	}{
		{"no bridge row", base, false},
		{"connected", base + on, true},
		{"connected, then turned off", base + on + off, false},
		{"turned off, then connected", base + off + on, true},
		{"a prompt that quotes a row", base + off + quoted, false},
		{"connected, then a tool call that names bridge-session", base + on + search, true},
		{"connected long before the end", on + filler, false},
	} {
		write(t, path, c.body, time.Now())
		if got := RemoteControl(path); got != c.want {
			t.Errorf("%s: %v", c.name, got)
		}
	}
	if RemoteControl(filepath.Join(t.TempDir(), "none.jsonl")) {
		t.Error("no conversation")
	}
}

// Agent is the agent the newest agent-setting row near the end of the
// conversation names.
func TestAgent(t *testing.T) {
	// Made-up rows of the shape Claude Code 2.1.314 writes.
	row := func(name string) string {
		return `{"type":"agent-setting","agentSetting":"` + name + `","sessionId":"s"}` + "\n"
	}
	search := `{"isSidechain":false,"type":"assistant","message":{"id":"m9","model":"claude-opus-5-5","role":"assistant","content":[{"type":"tool_use","id":"t9","name":"Grep","input":{"pattern":"agent-setting"}}]}}` + "\n"
	filler := strings.Repeat(`{"type":"system","subtype":"turn_duration","content":"`+strings.Repeat("x", 900)+`"}`+"\n", 1300)
	base := fixture(t, "limit-2.1.306-five_hour.jsonl")
	path := filepath.Join(t.TempDir(), "c.jsonl")
	for _, c := range []struct{ name, body, want string }{
		{"no row", base, ""},
		{"one row", row("reviewer") + base, "reviewer"},
		{"several rows", row("reviewer") + base + row("reviewer") + row("writer") + base, "writer"},
		{"a half-written last row", base + row("reviewer") + `{"type":"agent-setting","agentSetting":"wri`, "reviewer"},
		{"a tool call that names agent-setting", base + row("reviewer") + search, "reviewer"},
		{"a row long before the end", row("reviewer") + filler, ""},
	} {
		write(t, path, c.body, time.Now())
		if got := Agent(path); got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
	if got := Agent(filepath.Join(t.TempDir(), "none.jsonl")); got != "" {
		t.Errorf("no conversation: %q", got)
	}
}
