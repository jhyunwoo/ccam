//go:build !windows

package e2e

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// Ctrl+C or SIGTERM to ccam agent ends the agent as well: a tool that
// interrupts only ccam must not leave it running on another account.
func TestAgentEndsWithCcam(t *testing.T) {
	for _, sig := range []syscall.Signal{syscall.SIGINT, syscall.SIGTERM} {
		h := seatsHarness(t)
		cmd := exec.Command(ccamBin, "agent", "task")
		cmd.Env = append(append([]string{}, h.env...), "CCAM_PROFILE=c4", "FAKE_SLEEP=20", "FAKE_RECORD=interactive")
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		// The fake lists itself in its folder while it runs.
		var worker string
		for deadline := time.Now().Add(5 * time.Second); worker == "" && time.Now().Before(deadline); time.Sleep(20 * time.Millisecond) {
			for _, name := range []string{"c1", "c2", "c3", "cj"} {
				if m, _ := filepath.Glob(filepath.Join(h.profileDir(name), "sessions", "*.json")); len(m) > 0 {
					worker = m[0]
				}
			}
		}
		if worker == "" {
			cmd.Process.Kill()
			t.Fatalf("%v: the agent did not start", sig)
		}
		pid, _ := strconv.Atoi(strings.TrimSuffix(filepath.Base(worker), ".json"))
		cmd.Process.Signal(sig)
		done := make(chan error, 1)
		go func() { done <- cmd.Wait() }()
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			cmd.Process.Kill()
			t.Fatalf("%v: ccam agent did not end", sig)
		}
		if cmd.ProcessState.ExitCode() == 0 {
			t.Fatalf("%v: exit 0", sig)
		}
		if pid <= 0 || syscall.Kill(pid, 0) == nil {
			t.Fatalf("%v: the agent (pid %d) still runs", sig, pid)
		}
	}
}
