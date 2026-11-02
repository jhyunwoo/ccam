package transcript

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// SubagentPath is the conversation of the subagent agentID (a hook's
// agent_id) of the session whose conversation is conv; "" when there is
// none. Claude Code (2.1.314-2.1.316) keeps it below the session's, as
// <session id>/subagents/agent-<id>.jsonl, or deeper for a workflow's agents
// (subagents/workflows/<run>/).
func SubagentPath(conv, agentID string) string {
	if agentID == "" || strings.ContainsFunc(agentID, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_')
	}) {
		return ""
	}
	name := "agent-" + agentID + ".jsonl"
	dir := filepath.Join(strings.TrimSuffix(conv, ".jsonl"), "subagents")
	if path := filepath.Join(dir, name); regular(path) {
		return path
	}
	found := ""
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && d.Name() == name {
			found = path
			return fs.SkipAll
		}
		return nil
	})
	return found
}

// regular reports whether path is a regular file.
func regular(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}
