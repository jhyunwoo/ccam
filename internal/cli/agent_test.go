package cli

import (
	"regexp"
	"testing"
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
