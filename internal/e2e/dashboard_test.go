package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// A dashboard left open does next to nothing. Bubble Tea draws on a clock
// even when nothing changed, and at its default 60 frames a second an idle
// dashboard woke about 290 times a second (at 20, about 100).
func TestIdleDashboardRarelyWakes(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("counts wakeups in /proc")
	}
	h := newHarness(t)
	for _, name := range []string{"work", "personal"} {
		h.ok("ccam", "add", name, "--no-login")
	}
	pidFile := filepath.Join(h.home, "dashboard.pid")
	term := h.startTerminal(h.home, "stty rows 24 cols 80; echo $$ > "+pidFile+"; exec "+ccamBin)
	deadline := time.Now().Add(10 * time.Second)
	for !strings.Contains(term.out.String(), "personal") {
		if time.Now().After(deadline) {
			t.Fatalf("the dashboard did not show the accounts:\n%s", term.out.String())
		}
		time.Sleep(20 * time.Millisecond)
	}
	data, err := os.ReadFile(pidFile)
	if err != nil {
		t.Fatal(err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second) // the accounts' usage and the release check settle
	const sample = 2 * time.Second
	before := wakeups(t, pid)
	time.Sleep(sample)
	rate := float64(wakeups(t, pid)-before) / sample.Seconds()
	t.Logf("the idle dashboard woke %.0f times a second", rate)
	term.in.Write([]byte("q"))
	if _, err := term.wait(10 * time.Second); err != nil {
		t.Fatal(err)
	}
	if rate > 170 {
		t.Fatalf("the idle dashboard woke %.0f times a second, more than 170", rate)
	}
}

// wakeups is how many times the threads of process pid gave up the CPU to
// wait (voluntary context switches).
func wakeups(t *testing.T, pid int) int {
	t.Helper()
	tasks, err := filepath.Glob(fmt.Sprintf("/proc/%d/task/*/status", pid))
	if err != nil || len(tasks) == 0 {
		t.Fatalf("no threads of process %d: %v", pid, err)
	}
	n := 0
	for _, task := range tasks {
		data, err := os.ReadFile(task)
		if err != nil {
			continue // a thread that ended
		}
		for _, line := range strings.Split(string(data), "\n") {
			if v, ok := strings.CutPrefix(line, "voluntary_ctxt_switches:"); ok {
				c, _ := strconv.Atoi(strings.TrimSpace(v))
				n += c
			}
		}
	}
	return n
}
