package e2e

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/handoff"
)

// agentRow is a row Claude Code (2.1.314) writes into a conversation with
// the session's metadata for the conversation's agent (transcript.Agent).
func agentRow(name string) string {
	return `{"type":"agent-setting","agentSetting":"` + name + `","sessionId":"s1"}` + "\n"
}

// appendRows adds rows to the end of the conversation at path.
func appendRows(t *testing.T, path string, rows ...string) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, r := range rows {
		if _, err := f.WriteString(r); err != nil {
			t.Fatal(err)
		}
	}
}

// agentSetting sets the agent setting in account's settings.json (the
// group's, which the account shares), or takes it out for "".
func (h *harness) agentSetting(account, agent string) {
	h.t.Helper()
	path := filepath.Join(h.profileDir(account), "settings.json")
	v := map[string]any{}
	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, &v); err != nil {
			h.t.Fatal(err)
		}
	}
	v["agent"] = agent
	if agent == "" {
		delete(v, "agent")
	}
	data, _ := json.Marshal(v)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		h.t.Fatal(err)
	}
}

// The agent a conversation runs as goes on with it on the relay's hop:
// Claude Code 2.1.314 does not bring it back for an interactive --resume
// <file> --fork-session, and would go on as its default agent, with every
// tool. That is the session's --agent, else the agent setting, else the one
// the conversation's newest agent-setting row names. The rows do not win
// over the first two: a session that resumed a conversation by id runs as
// its --agent or the agent setting, and goes on writing the rows of the
// agent the conversation ran as before.
func TestRelayCarriesTheAgent(t *testing.T) {
	h, project, _ := relayHarness(t)
	for _, c := range []struct {
		name, args string
		setting    string // the agent setting personal and work share
		rows       []string
		want       []string
	}{
		{"started with --agent", "--agent reviewer", "", nil, []string{"--agent", "reviewer"}},
		{"named by the conversation", "", "", []string{agentRow("reviewer")}, []string{"--agent", "reviewer"}},
		{"the session's own over the rows", "--agent=reviewer --add-dir /tmp", "", []string{agentRow("writer")}, []string{"--agent=reviewer", "--add-dir", "/tmp"}},
		{"the agent setting over the rows", "", "reviewer", []string{agentRow("writer")}, []string{"--agent", "reviewer"}},
	} {
		h.agentSetting("personal", c.setting)
		conv := h.conversation("personal", project, "s1")
		appendRows(t, conv, c.rows...)
		shown := h.onTerminal(project, strings.TrimSpace(ccamBin+" personal "+c.args), "", h.switchEnv("personal", "work", conv, "FAKE_SWITCH_STOPPED=1")...)
		dir, args := h.called()
		want := append([]string{"--resume", conv, "--fork-session"}, c.want...)
		if dir != h.profileDir("work") || !slices.Equal(args, want) || strings.Contains(shown, "not carried over") || strings.Contains(shown, "define it differently") {
			t.Fatalf("%s: the relay started %s %q, want %q\n%s", c.name, dir, args, want, shown)
		}
	}
}

// /ccam:switch to an account that does not share the agents, plugins and
// settings of the conversation's account moves the conversation as asked,
// with its agent, and says, after the line about the move, that the
// account may not have that agent or may define it otherwise: resuming,
// Claude Code does not stop on an agent it does not have.
func TestRelaySwitchToAccountWithoutTheAgent(t *testing.T) {
	h, project, conv := relayHarness(t)
	h.ok("ccam", "group", "leave", "default", "work", "--yes")
	shown := h.onTerminal(project, ccamBin+" personal --agent reviewer", "", h.switchEnv("personal", "work", conv, "FAKE_SWITCH_STOPPED=1")...)
	dir, args := h.called()
	if want := []string{"--resume", conv, "--fork-session", "--agent", "reviewer"}; dir != h.profileDir("work") || !slices.Equal(args, want) {
		t.Fatalf("the relay started %s %q, want %q\n%s", dir, args, want, shown)
	}
	says := `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently`
	if moved := strings.Index(shown, "continuing the last conversation; work does not share personal's history"); moved < 0 || strings.Index(shown, says) < moved {
		t.Fatalf("the terminal did not say %q after the line about the move:\n%s", says, shown)
	}
}

// `ccam next --continue` typed by hand carries the agent the conversation
// runs as on too, and says so where the account may not have it, as the
// relay's hop does; an --agent given after -- wins.
func TestContinueKeepsTheAgent(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	project := filepath.Join(h.home, "app")
	os.MkdirAll(project, 0o700)
	project, _ = filepath.EvalSymlinks(project)
	conv := h.conversation("personal", project, "s1")
	appendRows(t, conv, agentRow("reviewer"))
	warns := `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer"`
	for _, c := range []struct {
		name  string
		leave bool // work leaves group default first
		args  []string
		want  []string
	}{
		{"agents shared", false, nil, []string{"--agent", "reviewer"}},
		{"an --agent of the person's", false, []string{"--", "--agent", "writer"}, []string{"--agent", "writer"}},
		{"agents of its own", true, nil, []string{"--agent", "reviewer"}},
	} {
		if c.leave {
			h.ok("ccam", "group", "leave", "default", "work", "--yes")
		}
		r := h.runIn(project, "", nil, append([]string{"ccam", "next", "--continue"}, c.args...)...)
		dir, args := h.called()
		want := append([]string{"--resume", conv, "--fork-session"}, c.want...)
		if r.code != 0 || dir != h.profileDir("work") || !slices.Equal(args, want) || strings.Contains(r.stderr, warns) != c.leave {
			t.Fatalf("%s: %s %q, want %q; %+v", c.name, dir, args, want, r)
		}
	}
}

// auto_switch carries a conversation that runs as an agent on only to an
// account that loads the same agents: one whose agents and plugins folders
// and settings.json are the very ones of the account at its limit (ccam's
// group links). The hop names the agent there. An account with room that
// does not is left out, and the notification says why.
func TestAutoSwitchKeepsTheAgent(t *testing.T) {
	h, project := autoHarness(t)
	conv := h.conversation("personal", project, "s1")
	appendRows(t, conv, agentRow("reviewer"))
	h.usageNow(map[string]int{"personal": 5, "work": 12})
	shown := h.onTerminal(project, ccamBin+" personal --agent reviewer", "", h.limitEnv("personal", conv)...)
	if args, want := h.hopped(shown, "personal", "work"), []string{"--resume", conv, "--fork-session", handoff.ContinuePrompt, "--agent", "reviewer"}; !slices.Equal(args, want) {
		t.Fatalf("work's claude got %q, want %q", args, want)
	}

	h, project = autoHarness(t)
	h.ok("ccam", "group", "leave", "default", "work", "--yes")
	conv = h.conversation("personal", project, "s1")
	appendRows(t, conv, agentRow("reviewer"))
	h.usageNow(map[string]int{"personal": 5, "work": 12})
	h.leftAlone(project, conv, `This conversation runs as agent "reviewer", and work does not share personal's agents, plugins and settings, so auto_switch did not move it there.`)

	// The agent setting alone names the agent too, before any row does.
	h, project = autoHarness(t)
	h.ok("ccam", "group", "leave", "default", "work", "--yes")
	h.agentSetting("personal", "reviewer")
	conv = h.conversation("personal", project, "s1")
	h.usageNow(map[string]int{"personal": 5, "work": 12})
	h.leftAlone(project, conv, `This conversation runs as agent "reviewer", and work does not share personal's agents, plugins and settings, so auto_switch did not move it there.`)
}

// The relay's hop after a usage limit, which nobody chose, does not carry
// a conversation that runs as an agent on to an account that does not load
// the same agents. The hook leaves such an account out when the agent
// setting or the conversation names the agent; here only the session's
// arguments do, as before Claude Code writes the first agent-setting row
// (after every 32 KiB of the conversation, and as the session ends). Nothing says the
// conversation moves once ccam knows it does not, and the desktop
// notification that said so is followed by one that says it did not. The
// note stays, and `ccam next --continue` typed later names the agent.
func TestAutoSwitchHopRefusesAnAccountWithoutTheAgent(t *testing.T) {
	h, project := autoHarness(t)
	h.ok("ccam", "group", "leave", "default", "work", "--yes")
	conv := h.conversation("personal", project, "s1")
	h.usageNow(map[string]int{"personal": 5, "work": 12})
	shown := h.onTerminal(project, ccamBin+" personal --agent reviewer", "", h.limitEnv("personal", conv)...)
	runs, terms := runsIn(filepath.Join(h.home, "runs")), runsIn(filepath.Join(h.home, "terms"))
	says := `ccam: this conversation runs as agent "reviewer", and work does not share personal's agents, plugins and settings, so ccam did not move it there.`
	next := " Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume s1` continues it as personal after the reset"
	if !slices.Equal(runs, []string{"personal"}) || !slices.Equal(terms, []string{"personal"}) || !strings.Contains(shown, says+next) ||
		!slices.Equal(h.notes("personal"), []string{"s1.json"}) {
		t.Fatalf("ran as %q, SIGTERM ended %q, notes %q; want the terminal to say %q:\n%s", runs, terms, h.notes("personal"), says+next, shown)
	}
	// The notification and the error line; no line said it moves since.
	if moved := strings.Index(shown, "Moving this conversation to work"); moved < 0 || strings.Count(shown[moved:], says) != 2 || strings.Contains(shown, "→ personal → work") {
		t.Fatalf("want a notification that says %q after the one that said it moves, and no line about the move:\n%s", says, shown)
	}
	appendRows(t, conv, agentRow("reviewer")) // as Claude Code ended
	r := h.runIn(project, "", nil, "ccam", "next", "--continue")
	dir, args := h.called()
	warns := `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer"`
	if want := []string{"--resume", conv, "--fork-session", "--agent", "reviewer"}; r.code != 0 || dir != h.profileDir("work") || !slices.Equal(args, want) || !strings.Contains(r.stderr, warns) {
		t.Fatalf("ccam next --continue: %s %q, want %q; %+v", dir, args, want, r)
	}
}
