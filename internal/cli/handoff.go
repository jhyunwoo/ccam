package cli

import (
	"fmt"
	"path/filepath"
	"slices"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/next"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/transcript"
)

// conversation is the one `ccam next --continue` carries on.
type conversation struct {
	path   string          // the .jsonl file, links resolved
	owner  *config.Profile // the account it belongs to (see lastConversation)
	cwd    string
	ticket *handoff.Ticket // the note the session left, if it is what continues
}

// lastConversation finds the conversation of the folder cwd that went on
// last (transcript.Latest) in any account. Accounts that share history see
// the same files; of them, the one Claude Code says ran it last owns it
// (ranLast), else the one used last. note, when not nil, is the note to
// carry on whatever else went on since (the relay passes its session's own).
func lastConversation(e *env, cfg *config.Config, cands []next.Candidate, cwd string, note *handoff.Ticket) (*conversation, error) {
	order := slices.Clone(cands)
	slices.SortStableFunc(order, func(a, b next.Candidate) int { return b.LastUsed.Compare(a.LastUsed) })
	var dirs []string
	sharers := map[string][]*config.Profile{} // the accounts that see each folder, used last first
	for _, c := range order {
		dir, err := filepath.EvalSymlinks(transcript.ProjectDir(profile.Dir(e.paths, c.Profile), cwd))
		if err != nil {
			continue
		}
		if sharers[dir] == nil {
			dirs = append(dirs, dir)
		}
		sharers[dir] = append(sharers[dir], c.Profile)
	}
	var best *conversation
	var at time.Time
	for _, dir := range dirs {
		// An agent's conversation (ccam agent) is not the person's to continue.
		skip := func(path string) bool { return agentConversation(e, path) }
		if path, t, ok := transcript.Latest(dir, skip); ok && (best == nil || t.After(at)) {
			owner := sharers[dir][0]
			if p := ranLast(e, sharers[dir], transcript.ID(path)); p != nil {
				owner = p
			}
			best, at = &conversation{path: path, owner: owner, cwd: cwd}, t
		}
	}
	// A note left by the session (a usage limit, /ccam:switch) names the
	// exact conversation, unless a conversation went on more than 2 minutes
	// after it. Rows Claude Code writes at /exit do not count: the person
	// may type /exit long after the note.
	t, ok := note, note != nil
	if !ok {
		t, ok = handoff.ForFolder(e.paths.State, cwd, time.Now())
	}
	if ok {
		if p, found := cfg.Find(t.Account); found && (note != nil || best == nil || !at.After(t.At.Add(2*time.Minute))) {
			path := t.Transcript
			if real, err := filepath.EvalSymlinks(path); err == nil {
				path = real
			}
			return &conversation{path: path, owner: p, cwd: cwd, ticket: t}, nil
		}
	}
	if best == nil {
		return nil, fmt.Errorf("no account has a conversation in %s to continue", cwd)
	}
	return best, nil
}

// ranLast is the account of profs (which share a history, used last first)
// that last ran the conversation with session id, by what Claude Code
// keeps in each account's own .claude.json when a session ends:
// projects[<folder>].lastSessionId. nil when none says so, or when the
// folder is no shared one.
func ranLast(e *env, profs []*config.Profile, id string) *config.Profile {
	if len(profs) < 2 {
		return nil
	}
	for _, p := range profs {
		var projects map[string]struct {
			LastSessionID string `json:"lastSessionId"`
		}
		if ok, err := profile.DecodeFileKey(filepath.Join(profile.Dir(e.paths, p), ".claude.json"), "projects", &projects); err != nil || !ok {
			continue
		}
		for _, pr := range projects {
			if pr.LastSessionID == id {
				return p
			}
		}
	}
	return nil
}

// continueArgs are the claude arguments that carry conv on as to, and a
// line about it: --resume <id> when to sees the conversation's folder
// (--continue could open a newer conversation there); else Claude Code
// copies the conversation into to (--fork-session), leaving the original as
// it is. A prompt other than "" goes right after them, before args: after
// a flag that takes several values (--add-dir a b) Claude Code would read
// it as one of them.
func continueArgs(e *env, cfg *config.Config, conv *conversation, to *config.Profile, prompt string, args []string) ([]string, string) {
	what := "the last conversation"
	if title := transcript.Title(conv.path); title != "" {
		what = fmt.Sprintf("%q", title)
	}
	if n := transcript.Size(conv.path); n > 0 {
		what += fmt.Sprintf(" (about %s tokens)", tokens(n))
		if org := newOrgs(e, cfg, false); otherOrg(org.of(conv.owner), org.of(to)) {
			what += "; prompt caches are not shared between organizations, so the first reply reads it all again"
		}
	}
	var resume []string
	note := "continuing " + what
	if mine, err := filepath.EvalSymlinks(transcript.ProjectDir(profile.Dir(e.paths, to), conv.cwd)); err == nil && mine == filepath.Dir(conv.path) {
		resume = []string{"--resume", transcript.ID(conv.path)}
	} else {
		resume = []string{"--resume", conv.path, "--fork-session"}
		note = fmt.Sprintf("continuing %s; %s does not share %s's history, so Claude Code copies it into %s", what, to.Name, conv.owner.Name, to.Name)
	}
	if prompt != "" {
		resume = append(resume, prompt)
	}
	return append(resume, args...), note
}

// otherOrg reports whether a and b are known to be in different
// organizations, which keep their prompt caches apart. Each personal plan
// is an organization of its own.
func otherOrg(a, b *config.Identity) bool {
	return a != nil && b != nil && a.OrgUUID != "" && b.OrgUUID != "" && a.OrgUUID != b.OrgUUID
}

// tokens writes n as 950, 12k or 1.2M.
func tokens(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1e6)
	case n >= 1000:
		return fmt.Sprintf("%dk", (n+500)/1000)
	}
	return fmt.Sprint(n)
}
