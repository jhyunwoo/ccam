package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/transcript"
)

// sameAgents reports whether the account folder to loads the agents of the
// account folder from the same way. Claude Code finds agents through the
// account's agents/ (its own), plugins/ (with enabledPlugins in
// settings.json) and settings.json (the agent setting); each must be the
// very same file in both, as ccam's group links make it (os.Stat follows
// them, and Windows junctions), or missing in both. A copy counts as
// another file: it may differ, or come to. The project's agents, Claude
// Code's own, managed ones and a --agents definition are the same for every
// account.
func sameAgents(from, to string) bool {
	for _, name := range []string{"agents", "plugins", "settings.json"} {
		a, errA := os.Stat(filepath.Join(from, name))
		b, errB := os.Stat(filepath.Join(to, name))
		switch {
		case errA == nil && errB == nil:
			if !os.SameFile(a, b) {
				return false
			}
		case !errors.Is(errA, fs.ErrNotExist) || !errors.Is(errB, fs.ErrNotExist):
			return false
		}
	}
	return true
}

// conversationAgent is the agent a session of the account folder dir in the
// folder project ran the conversation at path as when its claude args named
// none, as Claude Code (2.1.314) picks it: the agent setting of the
// project's .claude/settings.local.json, its .claude/settings.json or dir's
// settings.json, the first that has one; else the agent of the
// conversation's newest agent-setting row, which resuming it by id brings
// back. The setting comes first because it wins: a session that resumed a
// conversation by id runs as the setting's agent, and goes on writing rows
// for the agent the conversation ran as before. Managed settings, a
// --settings flag and plugins, which may set an agent too, are not read.
// "" for none.
func conversationAgent(dir, project, path string) string {
	files := []string{filepath.Join(dir, "settings.json")}
	if project != "" {
		files = append([]string{filepath.Join(project, ".claude", "settings.local.json"), filepath.Join(project, ".claude", "settings.json")}, files...)
	}
	for _, f := range files {
		var agent string
		if ok, _ := profile.DecodeFileKey(f, "agent", &agent); ok && agent != "" {
			return agent
		}
	}
	return transcript.Agent(path)
}

// agentMove checks the move of conv that `ccam next --continue` makes when
// the claude args args, which resume it as to, run it as an agent and to
// does not load the agents of conv's account the same way (sameAgents):
// to may lack that agent, and Claude Code, resuming, then goes on as its
// default agent, with every tool; or to may define it otherwise. The
// relay's hop after a usage limit (own), a move nobody chose, then starts
// nothing and keeps its note, and a desktop notification says so after the
// relay's that said the conversation moves. Any other move goes on, with
// the line agentMove returns.
func agentMove(e *env, own *handoff.Ticket, conv *conversation, to *config.Profile, args []string) (string, error) {
	from, agent := conv.owner, launcher.AgentArg(args)
	if agent == "" || from.Name == to.Name || sameAgents(profile.Dir(e.paths, from), profile.Dir(e.paths, to)) {
		return "", nil
	}
	differ := fmt.Sprintf("%s does not share %s's agents, plugins and settings", to.Name, from.Name)
	if own == nil || own.Reason != handoff.Limit {
		return fmt.Sprintf("%s, so it may not have agent %q or may define it differently: the conversation may go on without that agent's tool limits", differ, agent), nil
	}
	stays := fmt.Sprintf("this conversation runs as agent %q, and %s, so ccam did not move it there.", agent, differ)
	if isTerminal(e.stderr) {
		fmt.Fprint(e.stderr, handoff.Notification("ccam: "+stays))
	}
	return "", fmt.Errorf("%s Its note stays: `ccam next --continue` continues it later, and `ccam %s --resume %s` continues it as %s after the reset",
		stays, from.Name, transcript.ID(conv.path), from.Name)
}
