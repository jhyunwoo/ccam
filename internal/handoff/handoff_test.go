package handoff

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/transcript"
)

func TestLimitKind(t *testing.T) {
	cases := map[string]string{
		"You've hit your session limit · resets 3:45pm":                                             "session",
		"You’ve hit your weekly limit · resets Oct 23, 9am":                                         "weekly",
		"You've hit your Opus limit · resets Oct 23":                                                "model",
		"You've hit your Fable limit":                                                               "model",
		"You've hit your monthly spend limit":                                                       "spend",
		"API Error: Server is temporarily limiting requests (not your usage limit)":                 "",
		"API Error: Request rejected (429) · Number of request tokens has exceeded your rate limit": "",
	}
	for in, want := range cases {
		if got := LimitKind(in); got != want {
			t.Errorf("%q: %q, want %q", in, got, want)
		}
	}
	if got := Reset("You've hit your session limit · resets 3:45pm"); got != "3:45pm" {
		t.Fatal(got)
	}
	if Reset("no reset here") != "" {
		t.Fatal("reset")
	}
}

func TestForFolder(t *testing.T) {
	state, cwd := t.TempDir(), t.TempDir()
	conv := filepath.Join(t.TempDir(), "s.jsonl")
	os.WriteFile(conv, nil, 0o600)
	now := time.Now()
	for _, tk := range []Ticket{
		{Account: "a", SessionID: "s1", Cwd: cwd, Transcript: conv, At: now.Add(-time.Hour), Reason: Limit},
		{Account: "b", SessionID: "s2", Cwd: cwd, Transcript: conv, At: now.Add(-time.Minute), Reason: Switch, Target: "c"},
		{Account: "c", SessionID: "s3", Cwd: t.TempDir(), Transcript: conv, At: now, Reason: Limit},              // another folder
		{Account: "d", SessionID: "s4", Cwd: cwd, Transcript: conv + "x", At: now, Reason: Limit},                // gone
		{Account: "e", SessionID: "s5", Cwd: cwd, Transcript: conv, At: now.Add(-13 * time.Hour), Reason: Limit}, // too old
	} {
		if err := Save(state, tk); err != nil {
			t.Fatal(err)
		}
	}
	got, ok := ForFolder(state, cwd, now)
	if !ok || got.Account != "b" || got.Target != "c" || got.File() != Path(state, "b", "s2") {
		t.Fatalf("%+v", got)
	}
	Remove(got.File())
	if got, _ := ForFolder(state, cwd, now); got.Account != "a" {
		t.Fatalf("%+v", got)
	}
	// The note of an older ccam, one file for the account, is not read.
	os.WriteFile(filepath.Join(state, "handoff", "f.json"), fmt.Appendf(nil, `{"account":"f","transcript_path":%q,"cwd":%q,"at":%q,"reason":"switch"}`,
		conv, cwd, now.Format(time.RFC3339Nano)), 0o600)
	if got, _ := ForFolder(state, cwd, now); got.Account != "a" {
		t.Fatalf("the old note: %+v", got)
	}
}

// One account's sessions each keep a note of their own: a usage limit
// stops all of them at once.
func TestNotesPerSession(t *testing.T) {
	state := t.TempDir()
	now := time.Now()
	for _, id := range []string{"s1", "s2", "../../escape", ""} {
		if err := Save(state, Ticket{Account: "a", SessionID: id, At: now, Reason: Switch}); err != nil {
			t.Fatal(err)
		}
	}
	entries, _ := os.ReadDir(Dir(state, "a"))
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	if !reflect.DeepEqual(names, []string{"______escape.json", "s1.json", "s2.json", "unknown.json"}) {
		t.Fatalf("%q", names)
	}
	if err := Save(state, Ticket{Account: "../x", SessionID: "s1"}); err == nil {
		t.Fatal("an account that is no folder name")
	}
	// A note of a limit with an account chosen says why it was chosen.
	if err := Save(state, Ticket{Account: "a", SessionID: "s1", At: now, Reason: Limit, Chosen: "b", Why: "12% of its 5 hours used, the least"}); err != nil {
		t.Fatal(err)
	}
	t1, ok1 := LoadFile(state, Path(state, "a", "s1"))
	if !ok1 || t1.SessionID != "s1" || t1.Chosen != "b" || t1.Why != "12% of its 5 hours used, the least" {
		t.Fatalf("%+v", t1)
	}
	Remove(t1.File())
	if _, ok := LoadFile(state, Path(state, "a", "s1")); ok {
		t.Fatal("removed")
	}
	if t2, ok := LoadFile(state, Path(state, "a", "s2")); !ok || t2.SessionID != "s2" {
		t.Fatal("Remove took another session's note")
	}
	// Only notes are loaded, nothing else the path might name.
	other := filepath.Join(t.TempDir(), "a", "s2.json")
	os.MkdirAll(filepath.Dir(other), 0o700)
	os.WriteFile(other, []byte(`{"account":"a"}`), 0o600)
	for _, p := range []string{other, filepath.Join(state, "handoff", "s2.json"), filepath.Join(state, "handoff", "a", "..", "..", "x.json"), "handoff/a/s2.json"} {
		if _, ok := LoadFile(state, p); ok {
			t.Errorf("loaded %s", p)
		}
	}
	// Notes past MaxAge go when the account saves another, and so do the
	// stop-failure hook's locks beside them.
	old := now.Add(-13 * time.Hour)
	lock := filepath.Join(Dir(state, "a"), "s2.lock")
	os.WriteFile(lock, nil, 0o600)
	os.Chtimes(Path(state, "a", "s2"), old, old)
	os.Chtimes(lock, old, old)
	Save(state, Ticket{Account: "a", SessionID: "s3", At: now})
	if _, err := os.Stat(Path(state, "a", "s2")); err == nil {
		t.Fatal("an old note stays")
	}
	if _, err := os.Stat(lock); err == nil {
		t.Fatal("an old lock stays")
	}
}

// The relay acts only on what its own session left: its CCAM_RELAY, its
// claude process (or none named), since it started that claude.
func TestForRelay(t *testing.T) {
	state := t.TempDir()
	start := time.Now()
	for _, tk := range []Ticket{
		{SessionID: "other-relay", Relay: "2:9", Child: 100, At: start.Add(time.Second)},
		{SessionID: "other-child", Relay: "2:1", Child: 101, At: start.Add(time.Second)},
		{SessionID: "before", Relay: "2:1", Child: 100, At: start.Add(-time.Second)},
		{SessionID: "mine", Relay: "2:1", Child: 100, Stopped: true, At: start.Add(time.Second)},
	} {
		tk.Account, tk.Reason = "a", Switch
		Save(state, tk)
	}
	if got, ok := ForRelay(state, "a", "2:1", 100, start); !ok || got.SessionID != "mine" {
		t.Fatalf("%+v", got)
	}
	Remove(Path(state, "a", "mine"))
	if got, ok := ForRelay(state, "a", "2:1", 100, start); ok {
		t.Fatalf("%+v", got)
	}
	// Windows, or Claude Code without CLAUDE_PID: no process named.
	Save(state, Ticket{Account: "a", SessionID: "unnamed", Relay: "2:1", At: start.Add(time.Second), Reason: Switch})
	if got, ok := ForRelay(state, "a", "2:1", 100, start); !ok || got.SessionID != "unnamed" {
		t.Fatalf("%+v", got)
	}
}

// A rename carries the account's notes along.
func TestRename(t *testing.T) {
	state := t.TempDir()
	from, to := filepath.Join(t.TempDir(), "personal"), filepath.Join(t.TempDir(), "home")
	conv := filepath.Join(from, "projects", "-app", "s1.jsonl")
	Save(state, Ticket{Account: "personal", SessionID: "s1", Transcript: conv, Reason: Switch, Target: "work"})
	Save(state, Ticket{Account: "work", SessionID: "s2", Transcript: "/w/s2.jsonl", Reason: Switch, Target: "personal"})
	Save(state, Ticket{Account: "side", SessionID: "s3", Transcript: "/s/s3.jsonl", Reason: Limit, Chosen: "personal"})
	os.Rename(Dir(state, "personal"), Dir(state, "home")) // as rename does with Paths
	Rename(state, "personal", "home", from, to)
	if got, _ := LoadFile(state, Path(state, "home", "s1")); got.Account != "home" || got.Transcript != filepath.Join(to, "projects", "-app", "s1.jsonl") || got.Target != "work" {
		t.Fatalf("%+v", got)
	}
	if got, _ := LoadFile(state, Path(state, "work", "s2")); got.Target != "home" || got.Transcript != "/w/s2.jsonl" {
		t.Fatalf("%+v", got)
	}
	if got, _ := LoadFile(state, Path(state, "side", "s3")); got.Chosen != "home" {
		t.Fatalf("%+v", got)
	}
}

// The hook's cwd follows a cd in Claude Code's shell; the conversation
// stays under the folder the session started in.
func TestNoteFromASubfolderBelongsToItsStartFolder(t *testing.T) {
	cwd := t.TempDir()
	real, _ := filepath.EvalSymlinks(cwd)
	conv := filepath.Join(t.TempDir(), "projects", transcript.Slug(real), "s.jsonl")
	tk := Ticket{Cwd: filepath.Join(cwd, "sub"), Transcript: conv}
	if !tk.In(cwd) || tk.In(filepath.Join(cwd, "sub", "other")) {
		t.Fatal("matched by the conversation's folder")
	}
}
