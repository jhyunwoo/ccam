package e2e

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// agentRuns reads FAKE_AGENT_LOG: the account and arguments of each
// claude -p that ccam agent ran, in order.
func agentRuns(t *testing.T, path string) (accounts []string, args [][]string) {
	t.Helper()
	data, _ := os.ReadFile(path)
	os.Remove(path)
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line == "" {
			continue
		}
		name, a, _ := strings.Cut(line, "\t")
		accounts = append(accounts, name)
		args = append(args, strings.Split(a, " "))
	}
	return accounts, args
}

// argAfter is the argument after flag in args ("" for none).
func argAfter(args []string, flag string) string {
	for i, a := range args {
		if a == flag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

var sessionLine = regexp.MustCompile(`session ([0-9a-f-]{36})`)

// The agent runs as another seat of the caller's organization, never the
// caller's own account (nor a second folder of it), and agents started
// together get different seats.
func TestAgentRunsAsAnotherAccount(t *testing.T) {
	// other has the most room but another organization; dup is c1's
	// Claude account in a second folder.
	h := seatsHarness(t, seat{"other", "u6", "home", 1}, seat{"dup", "u1", "team", 1})
	log := filepath.Join(h.home, "agent.log")
	envFile := filepath.Join(h.home, "env")
	env := []string{"FAKE_AGENT_LOG=" + log, "FAKE_STDOUT=answer", "CCAM_PROFILE=c1", "CCAM_RELAY=2:1",
		"FAKE_ENV_FILE=" + envFile, "FAKE_ENV_NAMES=CCAM_AGENT,CCAM_RELAY,CCAM_PROFILE"}

	r := h.run("", env, "ccam", "agent", "fix", "the", "tests")
	if r.code != 0 || r.stdout != "answer (c3)\n" {
		t.Fatalf("agent: %+v", r)
	}
	m := sessionLine.FindStringSubmatch(r.stderr)
	if m == nil || !strings.Contains(r.stderr, "ccam: agent runs as c3 (20% of its 5 hours used); session "+m[1]) ||
		!strings.Contains(r.stderr, "ccam: agent finished as c3; follow up with: ccam agent --resume "+m[1]) {
		t.Fatalf("stderr: %s", r.stderr)
	}
	accounts, args := agentRuns(t, log)
	if len(accounts) != 1 || accounts[0] != "c3" {
		t.Fatalf("ran as %v", accounts)
	}
	got := strings.Join(args[0], " ")
	// No calling session to take a mode from: the default mode, named so
	// that the agent account's settings cannot give more.
	if !strings.Contains(got, "-p --output-format json --permission-mode default --session-id "+m[1]+" -- fix the tests") {
		t.Fatalf("args: %s", got)
	}
	// The agent knows it is one, and the caller's relay does not reach it.
	if b, _ := os.ReadFile(envFile); string(b) != "CCAM_AGENT=1\nCCAM_PROFILE=c3\n" {
		t.Fatalf("agent env: %q", b)
	}
	if _, err := os.Stat(filepath.Join(h.home, "state", "agent", m[1]+".json")); err != nil {
		t.Fatalf("no record: %v", err)
	}

	// Picked a moment ago, c3 goes last: the next agent gets cj; --json says so.
	r = h.run("", env, "ccam", "agent", "--json", "second task")
	if r.code != 0 || !strings.Contains(r.stdout, `"account": "cj"`) || !strings.Contains(r.stdout, `"result": "answer (cj)"`) || !strings.Contains(r.stdout, `"accounts": [`) {
		t.Fatalf("second agent: %+v", r)
	}
	// --cross-org also takes the other organization's account.
	r = h.run("", env, "ccam", "agent", "--cross-org", "third")
	if r.code != 0 || r.stdout != "answer (other)\n" {
		t.Fatalf("--cross-org: %+v", r)
	}
	agentRuns(t, log)

	// From a terminal (no CCAM_PROFILE) the account of the folder, c1 (the
	// default), is the one left out. The task comes from stdin.
	r = h.run("piped task\n", []string{"FAKE_AGENT_LOG=" + log}, "ccam", "agent")
	if r.code != 0 {
		t.Fatalf("stdin task: %+v", r)
	}
	if accounts, args := agentRuns(t, log); len(accounts) != 1 || accounts[0] == "c1" || accounts[0] == "dup" || !strings.HasSuffix(strings.Join(args[0], " "), "-- piped task") {
		t.Fatalf("stdin task ran as %v: %v", accounts, args)
	}
	if r := h.run("", nil, "ccam", "agent"); r.code == 0 || !strings.Contains(r.stderr, "give the task") {
		t.Fatalf("no task: %+v", r)
	}

	// A dry run says what would run and runs nothing.
	r = h.run("", env, "ccam", "--dry-run", "agent", "x")
	if r.code != 0 || !strings.Contains(r.stderr, "the agent would run as") || !strings.Contains(r.stdout, "--session-id") {
		t.Fatalf("dry run: %+v", r)
	}
	if accounts, _ := agentRuns(t, log); len(accounts) != 0 {
		t.Fatalf("dry run ran %v", accounts)
	}
}

// The agent works in the calling session's permission mode and effort,
// found by the session's id whatever folder the shell is in; an agent
// cannot start one of its own.
func TestAgentCarriesTheSessionModeAndEffort(t *testing.T) {
	h := seatsHarness(t)
	log := filepath.Join(h.home, "agent.log")
	project := t.TempDir()
	sub := filepath.Join(project, "sub")
	os.MkdirAll(sub, 0o700)
	real, _ := filepath.EvalSymlinks(project)
	conv := filepath.Join(h.profileDir("c3"), "projects", slugOf(real), "caller-session.jsonl")
	os.MkdirAll(filepath.Dir(conv), 0o700)
	os.WriteFile(conv, []byte(`{"type":"user","message":{"role":"user","content":"go"},"permissionMode":"acceptEdits","sessionId":"caller-session"}`+"\n"), 0o600)
	env := []string{"FAKE_AGENT_LOG=" + log, "CCAM_PROFILE=c3", "CLAUDE_CODE_SESSION_ID=caller-session", "CLAUDE_EFFORT=high"}

	if r := h.runIn(sub, "", env, "ccam", "agent", "task"); r.code != 0 {
		t.Fatalf("agent: %+v", r)
	}
	_, args := agentRuns(t, log)
	if got := strings.Join(args[0], " "); argAfter(args[0], "--permission-mode") != "acceptEdits" || argAfter(args[0], "--effort") != "high" {
		t.Fatalf("args: %s", got)
	}
	// An id that is no file name finds nothing: the default mode.
	env[2] = "CLAUDE_CODE_SESSION_ID=../caller-session"
	if r := h.runIn(sub, "", env, "ccam", "agent", "task"); r.code != 0 {
		t.Fatalf("agent: %+v", r)
	}
	if _, args := agentRuns(t, log); argAfter(args[0], "--permission-mode") != "default" {
		t.Fatalf("args: %v", args[0])
	}
	if r := h.runIn(sub, "", env, "ccam", "agent", "task"); !strings.Contains(r.stderr, "could not tell this session's permission mode") {
		t.Fatalf("no word of the unknown mode: %+v", r)
	}
	agentRuns(t, log)
	// A session ccam did not start is known by its config folder.
	env[0], env[1] = "FAKE_AGENT_LOG="+log, "CLAUDE_CONFIG_DIR="+h.profileDir("c3")
	env[2] = "CLAUDE_CODE_SESSION_ID=caller-session"
	if r := h.runIn(sub, "", env, "ccam", "agent", "task"); r.code != 0 {
		t.Fatalf("agent: %+v", r)
	}
	if accounts, args := agentRuns(t, log); accounts[0] == "c3" || argAfter(args[0], "--permission-mode") != "acceptEdits" {
		t.Fatalf("by config folder: %v %v", accounts, args[0])
	}
	r := h.run("", append(env, "CCAM_AGENT=1"), "ccam", "agent", "task")
	if r.code == 0 || !strings.Contains(r.stderr, "cannot start agents of its own") {
		t.Fatalf("nested: %+v", r)
	}
	if accounts, _ := agentRuns(t, log); len(accounts) != 0 {
		t.Fatalf("nested agent ran %v", accounts)
	}
}

// At a usage limit the agent's conversation goes on as another account, in
// a copy Claude Code makes (--resume <file> --fork-session), and the limit
// is recorded so the next choice leaves that account out. --resume
// continues an agent's conversation.
func TestAgentGoesOnAtTheLimit(t *testing.T) {
	h := seatsHarness(t)
	log := filepath.Join(h.home, "agent.log")
	env := []string{"FAKE_AGENT_LOG=" + log, "FAKE_STDOUT=answer", "CCAM_PROFILE=c4"}

	r := h.run("", append(env, "FAKE_LIMITED=c3"), "ccam", "agent", "long task")
	if r.code != 0 || r.stdout != "answer (cj)\n" || !strings.Contains(r.stderr, "ccam: c3 reached its usage limit; the agent goes on as cj") {
		t.Fatalf("limit: %+v", r)
	}
	accounts, args := agentRuns(t, log)
	if strings.Join(accounts, ",") != "c3,cj" {
		t.Fatalf("ran as %v", accounts)
	}
	first, second := argAfter(args[0], "--session-id"), argAfter(args[1], "--session-id")
	from := argAfter(args[1], "--resume")
	if first == "" || second == "" || first == second || filepath.Base(from) != first+".jsonl" || !strings.Contains(strings.Join(args[1], " "), "--fork-session") {
		t.Fatalf("second run: %v", args[1])
	}
	if !strings.Contains(strings.Join(args[1], " "), "Your usage limit was reached and ccam moved this conversation to another account.") {
		t.Fatalf("second run's prompt: %v", args[1])
	}
	if b, err := os.ReadFile(filepath.Join(h.home, "state", "usage", "c3.limit.json")); err != nil || !strings.Contains(string(b), "five_hour") {
		t.Fatalf("limit record: %s %v", b, err)
	}
	// The copy holds the conversation so far.
	real, _ := filepath.EvalSymlinks(mustGetwd(t))
	copied := filepath.Join(h.profileDir("cj"), "projects", slugOf(real), second+".jsonl")
	if b, err := os.ReadFile(copied); err != nil || !strings.Contains(string(b), "long task") {
		t.Fatalf("copy: %s %v", b, err)
	}

	// --resume: the account that ran it has room, so it goes on there.
	r = h.run("", env, "ccam", "agent", "--resume", second, "and a test")
	if r.code != 0 || r.stdout != "answer (cj)\n" {
		t.Fatalf("resume: %+v", r)
	}
	if accounts, args := agentRuns(t, log); len(accounts) != 1 || accounts[0] != "cj" || argAfter(args[0], "--resume") != copied {
		t.Fatalf("resume ran %v %v", accounts, args)
	}
	for _, id := range []string{"no-such-session", "../x"} {
		if r := h.run("", env, "ccam", "agent", "--resume", id, "x"); r.code == 0 || !strings.Contains(r.stderr, "no agent conversation") {
			t.Fatalf("resume %s: %+v", id, r)
		}
	}

	// Every seat but the caller at its limit: the last reply, and 75.
	r = h.run("", append(env, "FAKE_LIMITED=c1,c2,c3,cj"), "ccam", "agent", "x")
	if r.code != 75 || !strings.Contains(r.stdout, "hit your session limit") || !strings.Contains(r.stderr, "no other account the agent may use has room left") {
		t.Fatalf("all limited: %+v", r)
	}
	if accounts, _ := agentRuns(t, log); strings.Contains(strings.Join(accounts, ","), "c4") {
		t.Fatalf("ran as the caller: %v", accounts)
	}
}

// ccam next --continue continues the person's conversation, not an agent's
// that went on later in the same folder.
func TestNextContinueLeavesAgentConversationsOut(t *testing.T) {
	h := seatsHarness(t)
	dir := t.TempDir()
	real, _ := filepath.EvalSymlinks(dir)
	mine := filepath.Join(h.profileDir("c3"), "projects", slugOf(real), "my-session.jsonl")
	os.MkdirAll(filepath.Dir(mine), 0o700)
	at := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339Nano)
	os.WriteFile(mine, []byte(`{"type":"user","message":{"role":"user","content":"mine"},"timestamp":"`+at+`"}`+"\n"), 0o600)
	os.Chtimes(mine, time.Now().Add(-time.Hour), time.Now().Add(-time.Hour))

	if r := h.runIn(dir, "", []string{"CCAM_PROFILE=c3"}, "ccam", "agent", "task"); r.code != 0 {
		t.Fatalf("agent: %+v", r)
	}
	r := h.runIn(dir, "", nil, "ccam", "--dry-run", "next", "--continue")
	if r.code != 0 || !strings.Contains(r.stdout, "my-session") {
		t.Fatalf("next --continue: %+v", r)
	}
}

func mustGetwd(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return wd
}

// slugOf is the folder under projects/ for dir (transcript.Slug for short
// paths).
func slugOf(dir string) string {
	var b strings.Builder
	for _, c := range dir {
		if c < 0x80 && (c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9') {
			b.WriteRune(c)
		} else {
			b.WriteByte('-')
		}
	}
	return b.String()
}
