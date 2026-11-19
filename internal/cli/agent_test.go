package cli

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
)

func TestSessionID(t *testing.T) {
	for id, want := range map[string]bool{
		"a2d835df-fff1-4910-ba33-4afd9ca51435": true,
		"caller_session":                       true,
		"":                                     false,
		"../x":                                 false,
		"a/b":                                  false,
		"a.b":                                  false,
		string(make([]byte, 65)):               false,
	} {
		if got := sessionID(id); got != want {
			t.Errorf("sessionID(%q) = %v", id, got)
		}
	}
	uuid := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	a, b := newSessionID(), newSessionID()
	if !uuid.MatchString(a) || a == b || !sessionID(a) {
		t.Fatalf("%s %s", a, b)
	}
}

func TestAgentResult(t *testing.T) {
	r, ok := agentResult([]byte("warming up\n{\"type\":\"result\",\"subtype\":\"success\",\"is_error\":false,\"result\":\"line 1\\nline 2\"}\n"))
	if !ok || r.Result != "line 1\nline 2" || r.IsError {
		t.Fatalf("%+v %v", r, ok)
	}
	r, ok = agentResult([]byte(`{"type":"result","subtype":"error_max_turns","is_error":true}`))
	if !ok || !r.IsError || r.Subtype != "error_max_turns" {
		t.Fatalf("%+v %v", r, ok)
	}
	if _, ok := agentResult([]byte("plain text\n{\"type\":\"system\"}\n")); ok {
		t.Fatal("no result object")
	}
}

// A run that ends with 0 but gives no result, or an error result, fails.
func TestAgentFinishFails(t *testing.T) {
	var out, errb bytes.Buffer
	a := &agentRun{e: &env{stdout: &out, stderr: &errb}}
	prof := &config.Profile{Name: "work"}
	err := a.finish(prof, "id", []string{"work"}, claudeResult{}, false, []byte("oops\n"), 0)
	var code launcher.ExitCode
	if !errors.As(err, &code) || code != 1 || out.String() != "oops\n" || !strings.Contains(errb.String(), "the agent stopped as work") {
		t.Fatalf("%v %q %q", err, out.String(), errb.String())
	}
	out.Reset()
	if err := a.finish(prof, "id", []string{"work"}, claudeResult{Type: "result", Result: "done"}, true, nil, 0); err != nil || out.String() != "done\n" {
		t.Fatalf("%v %q", err, out.String())
	}
	if err := a.finish(prof, "id", nil, claudeResult{Type: "result", IsError: true, Result: "no"}, true, nil, 0); !errors.As(err, &code) || code != 1 {
		t.Fatalf("error result: %v", err)
	}
}
