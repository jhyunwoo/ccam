package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/handoff"
)

// acct is an account of autoHarness: its Claude login and organization.
type acct struct{ name, email, orgType, org string }

// Personal plans, each its own organization; Pro and Max are one pool.
var (
	personalMax = acct{"personal", "me@example.com", "claude_max", "org-p"}
	workMax     = acct{"work", "dev@company.com", "claude_max", "org-w"}
	sidePro     = acct{"side", "side@example.com", "claude_pro", "org-s"}
)

// usual ends the notification of a session that is left to be carried on
// by hand.
const usual = "`ccam next --continue` continues the conversation on another account."

// autoHarness adds accts (default: personalMax and workMax) logged in with
// their plans, turns auto_switch on (and only it: auto_switch turns the
// relay on by itself), and makes every account trust the folder app, which
// it returns with links resolved. Usage lookups cannot reach a server.
func autoHarness(t *testing.T, accts ...acct) (*harness, string) {
	t.Helper()
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("auto_switch works on Linux and macOS only")
	}
	if len(accts) == 0 {
		accts = []acct{personalMax, workMax}
	}
	h := newHarness(t)
	// A saved lookup gone stale fails at once on this closed port, and the
	// same numbers stand: no token leaves. util-linux script runs $SHELL -c.
	h.env = append(h.env, "HTTPS_PROXY=http://127.0.0.1:9", "https_proxy=http://127.0.0.1:9", "NO_PROXY=", "no_proxy=", "SHELL=/bin/sh")
	for _, a := range accts {
		r := h.run("", []string{"FAKE_EMAIL=" + a.email, "FAKE_ORG_TYPE=" + a.orgType, "FAKE_ORG_UUID=" + a.org}, "ccam", "add", a.name)
		if r.code != 0 || !strings.Contains(r.stdout, a.email) {
			t.Fatalf("add %s: %+v", a.name, r)
		}
	}
	h.configure(func(c *config.Config) { c.AutoSwitch = true })
	project := filepath.Join(h.home, "app")
	os.MkdirAll(project, 0o700)
	project, _ = filepath.EvalSymlinks(project)
	for _, a := range accts {
		h.trust(a.name, project, true)
	}
	return h, project
}

// configure changes config.toml with edit.
func (h *harness) configure(edit func(*config.Config)) {
	h.t.Helper()
	path := filepath.Join(h.home, "config", "config.toml")
	cfg, err := config.Load(path)
	if err == nil {
		edit(cfg)
		err = cfg.Save(path)
	}
	if err != nil {
		h.t.Fatal(err)
	}
}

// trust writes into account's .claude.json what Claude Code keeps there once
// its trust dialog was answered in folder: projects[folder].hasTrustDialogAccepted.
func (h *harness) trust(account, folder string, trusted bool) {
	h.t.Helper()
	path := filepath.Join(h.profileDir(account), ".claude.json")
	data, _ := os.ReadFile(path)
	var v map[string]any
	if err := json.Unmarshal(data, &v); err != nil {
		h.t.Fatal(err)
	}
	projects, _ := v["projects"].(map[string]any)
	if projects == nil {
		projects = map[string]any{}
	}
	projects[folder] = map[string]any{"hasTrustDialogAccepted": trusted}
	v["projects"] = projects
	data, _ = json.Marshal(v)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		h.t.Fatal(err)
	}
}

// usageNow saves each account's 5-hour usage as ccam saves a lookup, now,
// with its week at 10%: younger than usage.CacheTTL, it stands for a
// lookup, so nothing asks a server. It returns when the 5-hour windows
// reset, in two hours.
func (h *harness) usageNow(fiveHour map[string]int) time.Time {
	h.t.Helper()
	dir := filepath.Join(h.home, "state", "usage")
	os.MkdirAll(dir, 0o700)
	now := time.Now().UTC()
	reset := now.Add(2 * time.Hour).Truncate(time.Second)
	for name, pct := range fiveHour {
		data := fmt.Appendf(nil, `{"at":%q,"limits":{"five_hour":{"percent":%d,"resets_at":%q,"known":true},"weekly":{"percent":10,"resets_at":%q,"known":true}}}`,
			now.Format(time.RFC3339Nano), pct, reset.Format(time.RFC3339), now.Add(72*time.Hour).Format(time.RFC3339))
		if err := os.WriteFile(filepath.Join(dir, name+".json"), data, 0o600); err != nil {
			h.t.Fatal(err)
		}
	}
	return reset
}

// limitEnv makes the fake claude of the accounts in from (comma separated)
// stop at a usage limit in conversation conv and run the real stop-failure
// hook; the accounts it ran as, and those SIGTERM ended, go to home/runs and
// home/terms. extra comes last and wins.
func (h *harness) limitEnv(from, conv string, extra ...string) []string {
	return append([]string{"FAKE_HOOK=stop-failure", "FAKE_HOOK_FROM=" + from, "FAKE_CCAM=" + ccamBin, "FAKE_SWITCH_TRANSCRIPT=" + conv,
		"FAKE_RUNS=" + filepath.Join(h.home, "runs"), "FAKE_TERMS=" + filepath.Join(h.home, "terms")}, extra...)
}

// limited reports whether account's limit record is there.
func (h *harness) limited(account string) bool {
	return fileExists(filepath.Join(h.home, "state", "usage", account+".limit.json"))
}

// leftAlone runs personal's session in project to a usage limit (env on
// top of limitEnv) and fails unless ccam left the session where it was:
// Claude Code ended by itself, not by SIGTERM, no other account ran, the
// terminal showed says, and the note for `ccam next --continue` (s1.json)
// and personal's limit record are there.
func (h *harness) leftAlone(project, conv, says string, env ...string) {
	h.t.Helper()
	shown := h.onTerminal(project, ccamBin+" personal", "", h.limitEnv("personal", conv, env...)...)
	runs, terms := runsIn(filepath.Join(h.home, "runs")), runsIn(filepath.Join(h.home, "terms"))
	if !slices.Equal(runs, []string{"personal"}) || len(terms) != 0 || !strings.Contains(shown, says) {
		h.t.Fatalf("ran as %q, SIGTERM ended %q; want the terminal to say %q:\n%s", runs, terms, says, shown)
	}
	if notes := h.notes("personal"); !slices.Equal(notes, []string{"s1.json"}) || !h.limited("personal") {
		h.t.Fatalf("notes %q, limit recorded: %v", notes, h.limited("personal"))
	}
}

// hopped fails unless the session ran as accounts, in order, each but the
// last ended by SIGTERM (the hook's), and the terminal shown says the
// conversation moved to the last. It returns the args the last one got.
func (h *harness) hopped(shown string, accounts ...string) []string {
	h.t.Helper()
	runs, terms := runsIn(filepath.Join(h.home, "runs")), runsIn(filepath.Join(h.home, "terms"))
	last := accounts[len(accounts)-1]
	dir, args := h.called()
	if !slices.Equal(runs, accounts) || !slices.Equal(terms, accounts[:len(accounts)-1]) || dir != h.profileDir(last) ||
		!strings.Contains(shown, "Moving this conversation to "+last) {
		h.t.Fatalf("ran as %q, SIGTERM ended %q, the last claude ran in %s with %q:\n%s", runs, terms, dir, args, shown)
	}
	return args
}

// At a usage limit, auto_switch carries the conversation on by itself as
// another personal plan with room: the hook records the limit and ends
// Claude Code, and the relay (which auto_switch turns on) resumes it there
// with a prompt to go on, the arguments resuming does not restore, the
// permission mode the conversation was in and the session's effort level.
func TestAutoSwitchCarriesOn(t *testing.T) {
	h, project := autoHarness(t)
	conv := h.conversation("personal", project, "s1")
	envFile := filepath.Join(h.home, "env")
	h.usageNow(map[string]int{"personal": 5, "work": 12})
	// $$ is the shell that becomes ccam: the relay.
	line := "echo relay=$$; exec " + ccamBin + " personal --model opus --permission-mode plan --effort low --add-dir /tmp"
	shown := h.onTerminal(project, line, "", h.limitEnv("personal", conv, "FAKE_MODE=acceptEdits", "FAKE_EFFORT=high",
		"FAKE_ENV_FILE="+envFile, "FAKE_ENV_NAMES=CCAM_RELAY")...)
	args := h.hopped(shown, "personal", "work")
	if want := []string{"--resume", conv, "--fork-session", handoff.ContinuePrompt, "--add-dir", "/tmp", "--permission-mode", "acceptEdits", "--effort", "high"}; !slices.Equal(args, want) {
		t.Fatalf("work's claude got %q, want %q", args, want)
	}
	if !strings.Contains(shown, "ccam: personal hit its 5-hour limit; continuing this conversation as work") ||
		!strings.Contains(shown, "Moving this conversation to work (12% of its 5 hours used, the least).") || len(h.notes("personal")) != 0 {
		t.Fatalf("notes %q:\n%s", h.notes("personal"), shown)
	}
	// personal counts as full until its limit resets, in an hour.
	var record struct{ Windows map[string]time.Time }
	data, _ := os.ReadFile(filepath.Join(h.home, "state", "usage", "personal.limit.json"))
	err := json.Unmarshal(data, &record)
	if left := time.Until(record.Windows["five_hour"]); err != nil || left < 55*time.Minute || left > 61*time.Minute {
		t.Fatalf("limit record %s: %v", data, err)
	}
	// The hop took the relay's place: the same process, CCAM_RELAY as before.
	m := regexp.MustCompile(`relay=(\d+)`).FindStringSubmatch(shown)
	if b, _ := os.ReadFile(envFile); m == nil || string(b) != "CCAM_RELAY=2:"+m[1]+"\n" {
		t.Fatalf("work's claude had %q, the relay is %q", b, m)
	}
	if cfg, _ := os.ReadFile(filepath.Join(h.home, "config", "config.toml")); strings.Contains(string(cfg), "relay") {
		t.Fatalf("relay is on in config.toml:\n%s", cfg)
	}
}

// A subagent that stops at a usage limit stops the session's work too: the
// main thread may only wait for it. auto_switch carries the conversation on
// as for the main thread's limit, with a prompt that says what stopped.
func TestAutoSwitchSubagent(t *testing.T) {
	h, project := autoHarness(t)
	conv := h.conversation("personal", project, "s1")
	h.usageNow(map[string]int{"personal": 5, "work": 12})
	shown := h.onTerminal(project, ccamBin+" personal", "", h.limitEnv("personal", conv, "FAKE_AGENT_ID=a1")...)
	if args, want := h.hopped(shown, "personal", "work"), []string{"--resume", conv, "--fork-session", handoff.SubagentPrompt}; !slices.Equal(args, want) {
		t.Fatalf("work's claude got %q, want %q", args, want)
	}
	if !h.limited("personal") || len(h.notes("personal")) != 0 {
		t.Fatalf("limit recorded: %v, notes %q", h.limited("personal"), h.notes("personal"))
	}
}

// Seats of one Team organization stand in for each other, as in next
// --auto (D7 b): a seat at its limit moves on to another seat with room.
func TestAutoSwitchSameTeamOrg(t *testing.T) {
	h, project := autoHarness(t, acct{"personal", "me@example.com", "claude_team", "org-1"}, acct{"work", "dev@company.com", "claude_team", "org-1"})
	conv := h.conversation("personal", project, "s1")
	h.usageNow(map[string]int{"personal": 5, "work": 40})
	h.hopped(h.onTerminal(project, ccamBin+" personal", "", h.limitEnv("personal", conv)...), "personal", "work")
}

// Never across organizations (D7 b): a seat does not move to a seat of
// another organization nor to a personal plan, and a personal plan not to a
// seat, however much room they have, also with cross_org on, which only
// next --auto, pick and --fallback follow. With no other account of its
// pool, `ccam next --continue` would refuse too, and the notification says
// so; with cross_org on it would not, and the notification is the usual one.
func TestAutoSwitchOtherOrg(t *testing.T) {
	seat := acct{"personal", "me@example.com", "claude_team", "org-1"}
	for _, c := range []struct {
		name     string
		from, to acct
		crossOrg bool
		says     string
	}{
		{"a seat to a seat of another organization", seat, acct{"work", "dev@company.com", "claude_team", "org-2"}, false, "You have no other accounts in Acme, so ccam cannot move this conversation."},
		{"a seat to a personal plan", seat, workMax, false, "You have no other accounts in Acme, so ccam cannot move this conversation."},
		{"a personal plan to a seat", personalMax, acct{"work", "dev@company.com", "claude_team", "org-1"}, false, "You have no other personal plans, so ccam cannot move this conversation."},
		{"a personal plan to a seat, cross_org on", personalMax, acct{"work", "dev@company.com", "claude_team", "org-1"}, true, usual},
	} {
		t.Run(c.name, func(t *testing.T) {
			h, project := autoHarness(t, c.from, c.to)
			h.configure(func(cfg *config.Config) { cfg.CrossOrg = c.crossOrg })
			conv := h.conversation("personal", project, "s1")
			h.usageNow(map[string]int{"personal": 5, "work": 12})
			h.leftAlone(project, conv, c.says)
		})
	}
}

// With no other personal plan with room, the session stays and Claude Code
// waits for the reset; the notification says which account frees first
// (personal's limit resets after work's window here).
func TestAutoSwitchNoRoom(t *testing.T) {
	h, project := autoHarness(t)
	conv := h.conversation("personal", project, "s1")
	at := handoff.FormatReset(h.usageNow(map[string]int{"personal": 5, "work": 100}), time.Now())
	h.leftAlone(project, conv, "None of your other personal plans has room now; work frees up first, at "+at+
		". Claude Code waits for personal here; /ccam:switch work after "+at+" moves it sooner.", "FAKE_LIMIT_RESET=3h")
}

// A session connected to Remote Control is not moved (D11): it would not
// reconnect as another account. Once Remote Control was turned off (the
// conversation's last bridge row names no bridge session), it is.
func TestAutoSwitchRemoteControl(t *testing.T) {
	h, project := autoHarness(t)
	conv := h.conversation("personal", project, "s1")
	h.usageNow(map[string]int{"personal": 5, "work": 12})
	h.leftAlone(project, conv, "This session is connected to Remote Control and would not reconnect as another account, so ccam leaves it here.", "FAKE_BRIDGE=on")

	h, project = autoHarness(t)
	conv = h.conversation("personal", project, "s1")
	h.usageNow(map[string]int{"personal": 5, "work": 12})
	h.hopped(h.onTerminal(project, ccamBin+" personal", "", h.limitEnv("personal", conv, "FAKE_BRIDGE=off")...), "personal", "work")
}

// An account that has not trusted the folder is no target: Claude Code
// would stop at its trust question with nobody there.
func TestAutoSwitchUntrusted(t *testing.T) {
	h, project := autoHarness(t)
	h.trust("work", project, false)
	conv := h.conversation("personal", project, "s1")
	h.usageNow(map[string]int{"personal": 5, "work": 12})
	h.leftAlone(project, conv, "work has not trusted this folder yet, so the conversation was not moved there. Open the folder once as work to allow it.")
}

// An account whose login Claude Code cannot renew (it emptied the refresh
// token after the server refused it) is no target, however much room its
// saved usage shows: nobody may be there to log in again. The account with
// the most room after it is.
func TestAutoSwitchDeadLogin(t *testing.T) {
	h, project := autoHarness(t, personalMax, workMax, sidePro)
	if err := os.WriteFile(filepath.Join(h.profileDir("side"), ".credentials.json"),
		[]byte(`{"claudeAiOauth":{"accessToken":"fake","refreshToken":"","expiresAt":4102444800000}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	conv := h.conversation("personal", project, "s1")
	h.usageNow(map[string]int{"personal": 5, "work": 30, "side": 12})
	h.hopped(h.onTerminal(project, ccamBin+" personal", "", h.limitEnv("personal", conv)...), "personal", "work")
}

// A session without a relay of this ccam (not on a terminal, or started
// before auto_switch was on) cannot be carried on: the hook does not end it
// and says how to go on.
func TestAutoSwitchNotRelayed(t *testing.T) {
	h, project := autoHarness(t)
	conv := h.conversation("personal", project, "s1")
	h.usageNow(map[string]int{"personal": 5, "work": 12})
	r := h.runIn(project, "", h.limitEnv("personal", conv, "PWD="+project), "ccam", "personal")
	runs, terms := runsIn(filepath.Join(h.home, "runs")), runsIn(filepath.Join(h.home, "terms"))
	says := "This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`."
	if r.code != 0 || !slices.Equal(runs, []string{"personal"}) || len(terms) != 0 || !strings.Contains(r.stdout, says) ||
		!slices.Equal(h.notes("personal"), []string{"s1.json"}) || !h.limited("personal") {
		t.Fatalf("ran as %q, SIGTERM ended %q, notes %q, limit recorded %v: %+v", runs, terms, h.notes("personal"), h.limited("personal"), r)
	}
}

// With auto_switch off, a relayed session at its limit is left as before:
// the note and the usual notification.
func TestAutoSwitchOff(t *testing.T) {
	h, project := autoHarness(t)
	h.configure(func(c *config.Config) { c.AutoSwitch, c.Relay = false, true })
	conv := h.conversation("personal", project, "s1")
	h.usageNow(map[string]int{"personal": 5, "work": 12})
	h.leftAlone(project, conv, usual)
}

// A chain of limits moves on to the next personal plan with room each
// time, never back to one before its reset (its limit record stands over
// its saved 5%), and stays on the last when none has room. With shared
// history every hop resumes the conversation by its id.
func TestAutoSwitchChainStops(t *testing.T) {
	h, project := autoHarness(t, personalMax, workMax, sidePro)
	h.ok("ccam", "group", "set", "default", "--add", "history", "--yes")
	conv := h.conversation("personal", project, "s1")
	h.usageNow(map[string]int{"personal": 5, "work": 12, "side": 30})
	shown := h.onTerminal(project, ccamBin+" personal", "", h.limitEnv("personal,work,side", conv)...)
	// The second hop does not carry the first one's prompt on again.
	if args := h.hopped(shown, "personal", "work", "side"); !slices.Equal(args, []string{"--resume", "s1", handoff.ContinuePrompt}) {
		t.Fatalf("side's claude got %q", args)
	}
	for _, says := range []string{"continuing this conversation as work", "continuing this conversation as side", "None of your other personal plans has room now"} {
		if !strings.Contains(shown, says) {
			t.Fatalf("the terminal did not say %q:\n%s", says, shown)
		}
	}
	for _, a := range []string{"personal", "work", "side"} {
		var want []string
		if a == "side" {
			want = []string{"s1.json"}
		}
		if notes := h.notes(a); !slices.Equal(notes, want) || !h.limited(a) {
			t.Fatalf("%s: notes %q, limit recorded %v", a, notes, h.limited(a))
		}
	}
}

// Every session of an account stops at its limit at once. Two of them, in
// two folders, each carry their own conversation on, and not to the same
// account: the first one's choice puts that account behind the others with
// room for a while.
func TestAutoSwitchTwoSessions(t *testing.T) {
	h, _ := autoHarness(t, personalMax, workMax, sidePro)
	type session struct {
		dir, conv, log, runs, terms string
		term                        *terminal
	}
	var sessions []*session
	for _, name := range []string{"app1", "app2"} {
		dir := filepath.Join(h.home, name)
		os.MkdirAll(dir, 0o700)
		dir, _ = filepath.EvalSymlinks(dir)
		for _, a := range []string{"work", "side"} {
			h.trust(a, dir, true)
		}
		sessions = append(sessions, &session{dir: dir, conv: h.conversation("personal", dir, "conv-"+name), log: filepath.Join(h.home, name+".log"),
			runs: filepath.Join(h.home, name+".runs"), terms: filepath.Join(h.home, name+".terms")})
	}
	h.usageNow(map[string]int{"personal": 5, "work": 12, "side": 30})
	// Claims are not locked: the second session reaches its limit once the
	// first was ended, when its choice is made.
	gate := filepath.Join(h.home, "gate")
	for i, s := range sessions {
		env := h.limitEnv("personal", s.conv, "FAKE_CLAUDE_LOG="+s.log, "FAKE_RUNS="+s.runs, "FAKE_TERMS="+s.terms)
		if i == 1 {
			env = append(env, "FAKE_LIMIT_AFTER="+gate)
		}
		s.term = h.startTerminal(s.dir, ccamBin+" personal", env...)
	}
	waitFile(t, sessions[0].terms)
	if err := os.WriteFile(gate, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	for i, s := range sessions {
		shown, err := s.term.wait(20 * time.Second)
		if err != nil {
			t.Fatal(err)
		}
		to := []string{"work", "side"}[i]
		// work is not the least used for the second session, but taken.
		why := []string{"12% of its 5 hours used, the least", "30% of its 5 hours used; work was taken by another session a moment ago"}[i]
		data, _ := os.ReadFile(s.log)
		args := strings.Split(strings.TrimSuffix(string(data), "\x00"), "\x00")
		if !slices.Equal(runsIn(s.runs), []string{"personal", to}) || !slices.Equal(runsIn(s.terms), []string{"personal"}) || len(args) < 2 ||
			args[0] != "CLAUDE_CONFIG_DIR="+h.profileDir(to) || !slices.Equal(args[2:], []string{"--resume", s.conv, "--fork-session", handoff.ContinuePrompt}) ||
			!strings.Contains(shown, "Moving this conversation to "+to+" ("+why+").") {
			t.Fatalf("%s: ran as %q, SIGTERM ended %q, the last claude %q:\n%s", s.dir, runsIn(s.runs), runsIn(s.terms), args, shown)
		}
	}
	if notes := h.notes("personal"); len(notes) != 0 {
		t.Fatalf("notes left %q", notes)
	}
}
