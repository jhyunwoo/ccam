// Package handoff keeps the notes Claude Code sessions leave when they stop
// at a usage limit or ask to move to another account (/ccam:switch), so
// that `ccam next --continue` carries on that very conversation. The ccam
// plugin's hooks write them, one per session, in handoff/<account>/ of the
// state folder; nothing runs in the background.
package handoff

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/transcript"
)

// Reasons a ticket is written for.
const (
	Limit  = "limit"  // the session stopped at a usage limit
	Switch = "switch" // /ccam:switch asked to move on
)

// Ticket is the note one session left.
type Ticket struct {
	Account    string    `json:"account"`
	SessionID  string    `json:"session_id"`
	Transcript string    `json:"transcript_path"`
	Cwd        string    `json:"cwd"`
	At         time.Time `json:"at"`
	Reason     string    `json:"reason"`
	Limit      string    `json:"limit,omitempty"`  // session, weekly, usage, model, spend or org (LimitKind)
	Target     string    `json:"target,omitempty"` // the account /ccam:switch named
	Relay      string    `json:"relay,omitempty"`  // CCAM_RELAY of the relaying ccam that is to carry the session on
	// Child is the relay's claude process (CLAUDE_PID, or its parent when
	// Claude Code started itself again); 0 when it could not be told.
	Child   int    `json:"child,omitempty"`
	Stopped bool   `json:"stopped,omitempty"` // the hook ended Child with SIGTERM
	Effort  string `json:"effort,omitempty"`  // the session's effort level, which resuming does not bring back
	// Window is the usage window Claude Code found full (as in
	// transcript.Limit), ResetsAt when it starts over, LimitAt when the
	// request was refused.
	Window   string    `json:"window,omitempty"`
	ResetsAt time.Time `json:"resets_at,omitzero"`
	LimitAt  time.Time `json:"limit_at,omitzero"`
	// Chosen is the account chosen to carry the conversation on after a
	// limit; only the relay's own hop follows it.
	Chosen string `json:"chosen,omitempty"`
	// Why says why Chosen was chosen (its usage), for the relay's message.
	Why string `json:"why,omitempty"`
	// Subagent is the agent_id of the subagent whose loop stopped at the
	// limit; "" for the main thread's.
	Subagent string `json:"subagent,omitempty"`

	file string // where the note was read from
}

// MaxAge is how long a ticket stays good.
const MaxAge = 12 * time.Hour

func root(state string) string { return filepath.Join(state, "handoff") }

// Dir holds the notes of account's sessions.
func Dir(state, account string) string { return filepath.Join(root(state), account) }

// Path is where the note of account's session lives. Claude Code's session
// ids are UUIDs; anything but letters, digits, - and _ becomes _, so the
// name stays a file in Dir.
func Path(state, account, session string) string {
	if session == "" {
		session = "unknown"
	}
	b := []byte(session)
	for i, c := range b {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' || c == '_') {
			b[i] = '_'
		}
	}
	return filepath.Join(Dir(state, account), string(b)+".json")
}

// Save writes t as its session's note, and drops the account's notes that
// are older than MaxAge.
func Save(state string, t Ticket) error {
	if t.Account == "" || t.Account != filepath.Base(t.Account) || t.Account == "." || t.Account == ".." {
		return fmt.Errorf("no account to note the session for (%q)", t.Account)
	}
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	if err := fsutil.AtomicWrite(Path(state, t.Account, t.SessionID), data, 0o600); err != nil {
		return err
	}
	now := time.Now()
	entries, _ := os.ReadDir(Dir(state, t.Account))
	for _, e := range entries {
		// The hook's locks (.lock) beside the notes go with them.
		if info, err := e.Info(); err == nil && (strings.HasSuffix(e.Name(), ".json") || strings.HasSuffix(e.Name(), ".lock")) && now.Sub(info.ModTime()) > MaxAge {
			os.Remove(filepath.Join(Dir(state, t.Account), e.Name()))
		}
	}
	return nil
}

// LoadFile reads the note at path, which must be one in state's notes
// (handoff/<account>/<session>.json).
func LoadFile(state, path string) (*Ticket, bool) {
	rel, err := filepath.Rel(root(state), path)
	if err != nil || !filepath.IsLocal(rel) {
		return nil, false
	}
	account, name, ok := strings.Cut(filepath.ToSlash(rel), "/")
	if !ok || account == "" || strings.Contains(name, "/") || !strings.HasSuffix(name, ".json") {
		return nil, false
	}
	return read(path)
}

func read(path string) (*Ticket, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	var t Ticket
	if json.Unmarshal(data, &t) != nil {
		return nil, false
	}
	t.file = path
	return &t, true
}

// File is where the note was read from.
func (t *Ticket) File() string { return t.file }

// Remove drops the note at path (a Ticket's File), and no other.
func Remove(path string) {
	if path != "" {
		os.Remove(path)
	}
}

// Paths lists what Save writes for account, for rename and remove.
func Paths(state, account string) []string { return []string{Dir(state, account)} }

// readDir reads the notes in dir.
func readDir(dir string) []*Ticket {
	entries, _ := os.ReadDir(dir)
	var out []*Ticket
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		if t, ok := read(filepath.Join(dir, e.Name())); ok {
			out = append(out, t)
		}
	}
	return out
}

// all reads the notes of every account.
func all(state string) []*Ticket {
	accounts, _ := os.ReadDir(root(state))
	var out []*Ticket
	for _, a := range accounts {
		if a.IsDir() {
			out = append(out, readDir(filepath.Join(root(state), a.Name()))...)
		}
	}
	return out
}

// ForFolder returns the newest good note any session left in the folder cwd.
func ForFolder(state, cwd string, now time.Time) (*Ticket, bool) {
	var best *Ticket
	for _, t := range all(state) {
		if t.Good(now) && t.In(cwd) && (best == nil || t.At.After(best.At)) {
			best = t
		}
	}
	return best, best != nil
}

// ForRelay returns the newest note account's sessions left at since or
// later for the relaying ccam relay (its CCAM_RELAY value) about its claude
// process child. A note that names no process counts too: the hook cannot
// tell it on Windows, nor without CLAUDE_PID (Claude Code before 2.1.214).
func ForRelay(state, account, relay string, child int, since time.Time) (*Ticket, bool) {
	var best *Ticket
	for _, t := range readDir(Dir(state, account)) {
		if t.Relay == relay && !t.At.Before(since) && (t.Child == 0 || t.Child == child) && (best == nil || t.At.After(best.At)) {
			best = t
		}
	}
	return best, best != nil
}

// Rename follows `ccam rename` of account from to to, whose folder moved
// from fromDir to toDir, once its notes have moved to Dir(state, to): they
// now say to and their conversation's path runs through toDir, and a note
// of any account that names from as its target or chosen account names to.
func Rename(state, from, to, fromDir, toDir string) {
	for _, t := range all(state) {
		changed := false
		if t.Account == from {
			t.Account, changed = to, true
			if rest, ok := strings.CutPrefix(t.Transcript, fromDir+string(filepath.Separator)); ok {
				t.Transcript = filepath.Join(toDir, rest)
			}
		}
		if t.Target == from {
			t.Target, changed = to, true
		}
		if t.Chosen == from {
			t.Chosen, changed = to, true
		}
		if data, err := json.Marshal(t); err == nil && changed {
			_ = fsutil.AtomicWrite(t.file, data, 0o600)
		}
	}
}

// Good reports whether t is recent and its conversation is still there.
func (t *Ticket) Good(now time.Time) bool {
	return now.Sub(t.At) <= MaxAge && fsutil.Exists(t.Transcript)
}

// In reports whether t's conversation belongs to the folder cwd: the
// session started there, or Claude Code keeps the conversation under it.
func (t *Ticket) In(cwd string) bool {
	return samePath(t.Cwd, cwd) || filepath.Base(filepath.Dir(t.Transcript)) == filepath.Base(transcript.ProjectDir("", cwd))
}

func samePath(a, b string) bool {
	if a == b {
		return true
	}
	ra, err1 := filepath.EvalSymlinks(a)
	rb, err2 := filepath.EvalSymlinks(b)
	return err1 == nil && err2 == nil && ra == rb
}

// Reset is the "resets …" part of Claude Code's message, or "", without
// the time zone Claude Code adds ("3:45pm (Asia/Seoul)"), which is the
// local one.
func Reset(msg string) string {
	_, after, ok := strings.Cut(msg, "resets ")
	if !ok {
		return ""
	}
	after, _, _ = strings.Cut(after, "\n")
	after, _, _ = strings.Cut(after, " ·")
	after = strings.TrimRight(strings.TrimSpace(after), ".")
	if i := strings.LastIndex(after, " ("); i > 0 && strings.HasSuffix(after, ")") {
		after = after[:i]
	}
	return after
}
