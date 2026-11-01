package e2e

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
)

// Two seats of one Team organization (the fake claude names every
// organization Acme); personalMax and sidePro are personal plans.
var (
	acmeSeat  = acct{"work", "dev@company.com", "claude_team", "org-acme"}
	acmeSeat2 = acct{"work2", "dev2@company.com", "claude_team", "org-acme"}
)

// orgHarness adds accts logged in with their plans and organizations, as
// autoHarness does, turns the relay on, and returns the folder app with
// links resolved. Unlike autoHarness it runs on every system. As there,
// usage lookups cannot reach a server.
func orgHarness(t *testing.T, accts ...acct) (*harness, string) {
	t.Helper()
	h := newHarness(t)
	h.env = append(h.env, "HTTPS_PROXY=http://127.0.0.1:9", "https_proxy=http://127.0.0.1:9", "NO_PROXY=", "no_proxy=")
	for _, a := range accts {
		r := h.run("", []string{"FAKE_EMAIL=" + a.email, "FAKE_ORG_TYPE=" + a.orgType, "FAKE_ORG_UUID=" + a.org}, "ccam", "add", a.name)
		if r.code != 0 || !strings.Contains(r.stdout, a.email) {
			t.Fatalf("add %s: %+v", a.name, r)
		}
	}
	h.configure(func(c *config.Config) { c.Relay = true })
	project := filepath.Join(h.home, "app")
	os.MkdirAll(project, 0o700)
	project, _ = filepath.EvalSymlinks(project)
	return h, project
}

// shows waits up to 10 s for the terminal to show s.
func (term *terminal) shows(t *testing.T, s string) {
	t.Helper()
	for deadline := time.Now().Add(10 * time.Second); !strings.Contains(term.out.String(), s); time.Sleep(20 * time.Millisecond) {
		if time.Now().After(deadline) {
			t.Fatalf("the terminal did not show %q:\n%s", s, term.out.String())
		}
	}
}

// Claude Code expands a /ccam:switch that a job the model scheduled fires
// like one the person typed, so the relay moves the conversation out of
// its account's organization only after the person answers yes on the
// terminal. What was typed before the question does not answer it; with no,
// Claude Code does not start as the other account, the note stays and ccam
// says how to go on.
func TestRelayAsksBeforeAnotherOrganization(t *testing.T) {
	h, project := orgHarness(t, acmeSeat, personalMax)
	conv := h.conversation("work", project, "s1")
	runs := filepath.Join(h.home, "runs")
	ask := "ccam: /ccam:switch named personal, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to personal? [y/N]"

	// A note for which the hook ended Claude Code, as hooks before this
	// version did for any account. A y typed before the question (here
	// while Claude Code ran) is dropped; the answer is n.
	term := h.startTerminal(project, ccamBin+" work", h.switchEnv("work", "personal", conv, "FAKE_SWITCH_STOPPED=1", "FAKE_EXIT_DELAY=1.5", "FAKE_RUNS="+runs)...)
	waitFile(t, filepath.Join(h.home, "state", "handoff", "work", "s1.json"))
	time.Sleep(200 * time.Millisecond)
	term.in.Write([]byte("y\n"))
	term.shows(t, ask)
	term.in.Write([]byte("n\n"))
	shown, err := term.wait(20 * time.Second)
	if err != nil {
		t.Fatal(err)
	}
	stays := "ccam: the conversation stays on work, a Claude Team seat in Acme: it moves to personal, a personal plan, only when you answer yes on a terminal. " +
		"`ccam next --continue` asks again, and `ccam work --resume s1` continues it as work"
	if !slices.Equal(runsIn(runs), []string{"work"}) || !slices.Equal(h.notes("work"), []string{"s1.json"}) || !strings.Contains(shown, stays) {
		t.Fatalf("ran as %q, notes %q:\n%s", runsIn(runs), h.notes("work"), shown)
	}

	// The note as the hook now leaves it for an account outside the pool:
	// Claude Code goes on until the person's /exit. Answered y, the relay
	// carries the conversation on there.
	os.Remove(runs)
	term = h.startTerminal(project, ccamBin+" work", h.switchEnv("work", "personal", conv, "FAKE_EXIT_DELAY=0.5", "FAKE_RUNS="+runs)...)
	term.shows(t, ask)
	term.in.Write([]byte("y\n"))
	if shown, err = term.wait(20 * time.Second); err != nil {
		t.Fatal(err)
	}
	if dir, args := h.called(); dir != h.profileDir("personal") || !slices.Equal(args, []string{"--resume", conv, "--fork-session"}) ||
		!slices.Equal(runsIn(runs), []string{"work", "personal"}) || len(h.notes("work")) != 0 {
		t.Fatalf("ran as %q, the last claude in %s with %q, notes %q:\n%s", runsIn(runs), dir, args, h.notes("work"), shown)
	}
}

// A /ccam:switch that names no account leaves the choice to ccam, and
// cross_org lets that choice leave the organization pool. The relay's hop
// then asks the same question before it moves the conversation there: with
// no, it stays; with yes, it moves. A choice within the pool moves at once.
func TestRelayAsksBeforeAPickInAnotherOrganization(t *testing.T) {
	h, project := orgHarness(t, acmeSeat, acmeSeat2, personalMax)
	h.configure(func(c *config.Config) { c.CrossOrg = true })
	h.usageNow(map[string]int{"work": 50, "work2": 40, "personal": 5})
	conv := h.conversation("work", project, "s1")
	runs := filepath.Join(h.home, "runs")
	ask := "ccam: the account picked for /ccam:switch is personal, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to personal? [y/N]"
	stays := "ccam: the conversation stays on work, a Claude Team seat in Acme: it moves to personal, a personal plan, only when you answer yes on a terminal. " +
		"`ccam next --continue` picks again, and `ccam work --resume s1` continues it as work"
	for _, answer := range []string{"n", "y"} {
		os.Remove(runs)
		term := h.startTerminal(project, ccamBin+" work", h.switchEnv("work", "", conv, "FAKE_SWITCH_STOPPED=1", "FAKE_RUNS="+runs)...)
		term.shows(t, ask)
		term.in.Write([]byte(answer + "\n"))
		shown, err := term.wait(20 * time.Second)
		if err != nil {
			t.Fatal(err)
		}
		switch {
		case answer == "n" && (!slices.Equal(runsIn(runs), []string{"work"}) || !slices.Equal(h.notes("work"), []string{"s1.json"}) || !strings.Contains(shown, stays)):
			t.Fatalf("no: ran as %q, notes %q:\n%s", runsIn(runs), h.notes("work"), shown)
		case answer == "y" && (!slices.Equal(runsIn(runs), []string{"work", "personal"}) || len(h.notes("work")) != 0):
			t.Fatalf("yes: ran as %q, notes %q:\n%s", runsIn(runs), h.notes("work"), shown)
		}
	}

	h.usageNow(map[string]int{"work": 50, "work2": 40, "personal": 90})
	os.Remove(runs)
	shown := h.onTerminal(project, ccamBin+" work", "", h.switchEnv("work", "", conv, "FAKE_SWITCH_STOPPED=1", "FAKE_RUNS="+runs)...)
	if strings.Contains(shown, "[y/N]") || !slices.Equal(runsIn(runs), []string{"work", "work2"}) {
		t.Fatalf("within the pool: ran as %q:\n%s", runsIn(runs), shown)
	}
}

// Within the organization pool the relay moves the conversation as before,
// asking nothing: to another seat of the organization, or from a personal
// plan to another.
func TestRelayMovesWithinThePoolWithoutAsking(t *testing.T) {
	h, project := orgHarness(t, acmeSeat, acmeSeat2, personalMax, sidePro)
	runs := filepath.Join(h.home, "runs")
	for _, c := range []struct{ from, to string }{{"work", "work2"}, {"personal", "side"}} {
		os.Remove(runs)
		conv := h.conversation(c.from, project, "s-"+c.from)
		shown := h.onTerminal(project, ccamBin+" "+c.from, "", h.switchEnv(c.from, c.to, conv, "FAKE_SWITCH_STOPPED=1", "FAKE_RUNS="+runs)...)
		dir, args := h.called()
		if strings.Contains(shown, "[y/N]") || !slices.Equal(runsIn(runs), []string{c.from, c.to}) || dir != h.profileDir(c.to) ||
			!slices.Equal(args, []string{"--resume", conv, "--fork-session"}) {
			t.Fatalf("%s to %s: ran as %q, the last claude in %s with %q:\n%s", c.from, c.to, runsIn(runs), dir, args, shown)
		}
	}
}

// `ccam next --continue` by hand follows a note /ccam:switch left (with
// the relay off, or after it did not move), so it asks the same question.
// Without a terminal to ask on the answer is no, whatever stdin holds, and
// cross_org, which lets ccam's own choices leave the organization, does not
// change that: Claude Code does not start and the note stays. An API key
// account is in no organization's pool, whichever side of the move it is
// on. A dry run shows the move and asks nothing.
func TestNextWithoutATerminalKeepsANoteToAnotherOrganization(t *testing.T) {
	h, project := orgHarness(t, acmeSeat, personalMax)
	h.configure(func(c *config.Config) { c.CrossOrg = true })
	h.ok("ccam", "add", "gw", "--api", "--base-url", "https://llm.example.com/v1", "--auth", "bearer", "--key-from", "env:GW_KEY")
	key := []string{"GW_KEY=sk-gw-123"}
	for _, c := range []struct{ from, to, says string }{
		{"work", "personal", "ccam: the conversation stays on work, a Claude Team seat in Acme: it moves to personal, a personal plan, only when you answer yes on a terminal."},
		{"work", "gw", "ccam: the conversation stays on work, a Claude Team seat in Acme: it moves to gw, an API key account, only when you answer yes on a terminal."},
		{"gw", "work", "ccam: the conversation stays on gw, an API key account: it moves to work, a Claude Team seat in Acme, only when you answer yes on a terminal."},
	} {
		conv := h.conversation(c.from, project, "s1")
		writeNote(t, h, c.from, conv, project, c.to)
		os.Remove(h.log) // left behind by auth login, or the case before
		r := h.runIn(project, "y\n", key, "ccam", "next", "--continue")
		if r.code != 1 || fileExists(h.log) || !slices.Equal(h.notes(c.from), []string{"s1.json"}) || !strings.Contains(r.stderr, c.says) ||
			!strings.Contains(r.stderr, "`ccam "+c.from+" --resume s1` continues it as "+c.from) {
			t.Fatalf("%s to %s: claude ran: %v, notes %q: %+v", c.from, c.to, fileExists(h.log), h.notes(c.from), r)
		}
		os.Remove(conv) // so the next case's note names the folder's only conversation
		os.Remove(filepath.Join(h.home, "state", "handoff", c.from, "s1.json"))
	}
	conv := h.conversation("work", project, "s1")
	writeNote(t, h, "work", conv, project, "personal")
	if r := h.runIn(project, "", nil, "ccam", "--dry-run", "next", "--continue"); r.code != 0 || !strings.Contains(r.stderr, "→ work → personal") || strings.Contains(r.stderr, "[y/N]") {
		t.Fatalf("dry run: %+v", r)
	}
}

// When /ccam:switch named no account, `ccam next --continue` by hand picks
// one: the next account, or with --auto the one with the most room. cross_org
// lets the pick leave the organization pool, and the move asks as for an
// account /ccam:switch named: without a terminal the note stays.
func TestNextAsksBeforeAPickInAnotherOrganization(t *testing.T) {
	h, project := orgHarness(t, acmeSeat, personalMax)
	h.configure(func(c *config.Config) { c.CrossOrg = true })
	h.usageNow(map[string]int{"work": 50, "personal": 5})
	conv := h.conversation("work", project, "s1")
	says := "ccam: the conversation stays on work, a Claude Team seat in Acme: it moves to personal, a personal plan, only when you answer yes on a terminal. " +
		"`ccam next --continue` picks again, and `ccam work --resume s1` continues it as work"
	for _, args := range [][]string{{"ccam", "next", "--continue"}, {"ccam", "next", "--auto", "--continue"}} {
		writeNote(t, h, "work", conv, project, "")
		os.Remove(h.log)
		r := h.runIn(project, "y\n", nil, args...)
		if r.code != 1 || fileExists(h.log) || !slices.Equal(h.notes("work"), []string{"s1.json"}) || !strings.Contains(r.stderr, says) {
			t.Fatalf("%q: claude ran: %v, notes %q: %+v", args, fileExists(h.log), h.notes("work"), r)
		}
	}
}
