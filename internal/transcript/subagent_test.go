package transcript

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// A subagent's conversation sits below its session's: subagents/ for the
// Agent tool's, subagents/workflows/<run>/ for a workflow's (Claude Code
// 2.1.314-2.1.316). An agent id that is not a plain name finds nothing.
func TestSubagentPath(t *testing.T) {
	dir := t.TempDir()
	conv := filepath.Join(dir, "s1.jsonl")
	agent := filepath.Join(dir, "s1", "subagents", "agent-a1.jsonl")
	workflow := filepath.Join(dir, "s1", "subagents", "workflows", "wf_1", "agent-a2.jsonl")
	write(t, conv, "{}\n", time.Now())
	write(t, agent, "{}\n", time.Now())
	write(t, workflow, "{}\n", time.Now())
	write(t, filepath.Join(dir, "s2", "subagents", "agent-a3.jsonl"), "{}\n", time.Now())
	for id, want := range map[string]string{"a1": agent, "a2": workflow, "a3": "", "": "", "../s2/subagents/agent-a3": "", "a*": ""} {
		if got := SubagentPath(conv, id); got != want {
			t.Errorf("%q: %q, want %q", id, got, want)
		}
	}
}

// A subagent's refusal at a usage limit is a row of its own conversation,
// isSidechain like every row there, with the model of its last reply.
func TestSubagentLimit(t *testing.T) {
	base := fixture(t, "limit-2.1.306-five_hour.jsonl")
	side := strings.ReplaceAll(base, `"isSidechain":false`, `"isSidechain":true`)
	if side == base {
		t.Fatal("the fixture changed")
	}
	path := filepath.Join(t.TempDir(), "agent-a1.jsonl")
	write(t, path, side, time.Now())
	at := time.Date(2026, 10, 21, 14, 51, 46, 515e6, time.UTC)
	l, ok := SubagentLimit(path, time.Time{})
	if !ok || l.Window != "five_hour" || !l.ResetsAt.Equal(time.Unix(1792607400, 0)) || !l.At.Equal(at) || l.Model != "claude-opus-5-5" ||
		l.Text != "You've hit your session limit · resets 3:30am (Asia/Seoul)" {
		t.Fatalf("%+v %v", l, ok)
	}
	if _, ok := SubagentLimit(path, at.Add(time.Second)); ok {
		t.Fatal("a limit before since")
	}
	// The main thread's rows are LastLimit's.
	if l, ok := LastLimit(path, time.Time{}); ok {
		t.Fatalf("LastLimit read a subagent's row: %+v", l)
	}
	write(t, path, base, time.Now())
	if l, ok := SubagentLimit(path, time.Time{}); ok {
		t.Fatalf("a main-thread row: %+v", l)
	}
}

// For a limit outside the conversation (a subagent's), only a prompt of
// the person's sent after it counts as moving on.
func TestPromptedSince(t *testing.T) {
	base := fixture(t, "limit-2.1.306-five_hour.jsonl") // its prompt was sent at 14:50:00
	path := filepath.Join(t.TempDir(), "c.jsonl")
	write(t, path, base, time.Now())
	sent := time.Date(2026, 10, 21, 14, 50, 0, 0, time.UTC)
	if PromptedSince(path, sent) || !PromptedSince(path, sent.Add(-time.Second)) {
		t.Fatal("the fixture's prompt")
	}
	later := sent.Add(time.Hour)
	stamp := `"timestamp":"2026-10-21T15:50:01.000Z"`
	for name, row := range map[string]string{
		"task notification": `{"type":"user",` + stamp + `,"message":{"role":"user","content":"(done)"},"origin":{"kind":"task-notification"}}`,
		"tool result":       `{"type":"user",` + stamp + `,"message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"t1","content":"ok"}]}}`,
		"subagent":          `{"type":"user",` + stamp + `,"isSidechain":true,"message":{"role":"user","content":"(task)"}}`,
	} {
		write(t, path, base+row+"\n", time.Now())
		if PromptedSince(path, later) {
			t.Errorf("%s: a prompt", name)
		}
	}
	write(t, path, base+`{"type":"user",`+stamp+`,"message":{"role":"user","content":"(go on)"},"origin":{"kind":"human"}}`+"\n", time.Now())
	if !PromptedSince(path, later) || PromptedSince(path, later.Add(2*time.Second)) {
		t.Fatal("a typed prompt")
	}
	if PromptedSince(filepath.Join(t.TempDir(), "none.jsonl"), later) {
		t.Fatal("no conversation")
	}
}
