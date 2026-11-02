package main

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// hookGrace is how long a session goes on after its StopFailure hook when
// the hook did not end it: the hook sends SIGTERM before it exits, so none
// comes later. Meanwhile the SIGTERM handler (main) may still end it.
const hookGrace = time.Second

// stopFailure plays a turn that Claude Code stops at a usage limit: it
// writes what Claude Code writes into the conversation, runs the ccam
// plugin's StopFailure hook (the real `ccam hook stop-failure` at FAKE_CCAM)
// with the input and environment Claude Code gives a hook, and shows the
// hook's terminalSequence as Claude Code does. Unlike Claude Code it waits
// for the hook. The rows are synthetic: Claude Code 2.1.306's keys,
// placeholder values.
//
// The conversation is FAKE_SWITCH_TRANSCRIPT, of session FAKE_SWITCH_SESSION
// (by default the file's name). Claude Code says FAKE_LIMIT_TEXT, and the
// 5-hour window resets FAKE_LIMIT_RESET (a Go duration, 1h) from now.
// FAKE_MODE writes the permission mode first, FAKE_BRIDGE a Remote Control
// connection ("off": turned off again after it), FAKE_EFFORT is the
// session's effort level, and with FAKE_LIMIT_AFTER the limit comes once
// that file is there. With FAKE_AGENT_ID the limit stops that subagent's
// loop instead: its row goes into the subagent's own conversation, and the
// hook gets its agent_id.
func stopFailure() {
	if v := os.Getenv("FAKE_HOOK"); v != "stop-failure" {
		fmt.Fprintf(os.Stderr, "fakeclaude: FAKE_HOOK=%s: only stop-failure is played\n", v)
		os.Exit(2)
	}
	ccam := os.Getenv("FAKE_CCAM")
	if ccam == "" {
		fmt.Fprintln(os.Stderr, "fakeclaude: FAKE_HOOK needs FAKE_CCAM, the ccam the hook runs")
		os.Exit(2)
	}
	if gate := os.Getenv("FAKE_LIMIT_AFTER"); gate != "" {
		// The session works on until the test lets the limit come.
		for end := time.Now().Add(20 * time.Second); time.Now().Before(end); time.Sleep(20 * time.Millisecond) {
			if _, err := os.Stat(gate); err == nil {
				break
			}
		}
	}
	conv := os.Getenv("FAKE_SWITCH_TRANSCRIPT")
	session := cmp.Or(os.Getenv("FAKE_SWITCH_SESSION"), strings.TrimSuffix(filepath.Base(conv), ".jsonl"))
	text := cmp.Or(os.Getenv("FAKE_LIMIT_TEXT"), "You've hit your limit · resets 3pm")
	after, err := time.ParseDuration(cmp.Or(os.Getenv("FAKE_LIMIT_RESET"), "1h"))
	must(err)
	var rows []string
	if mode := os.Getenv("FAKE_MODE"); mode != "" {
		rows = append(rows, fmt.Sprintf(`{"type":"permission-mode","permissionMode":%s,"sessionId":%s}`, js(mode), js(session)))
	}
	if bridge := os.Getenv("FAKE_BRIDGE"); bridge != "" {
		rows = append(rows, fmt.Sprintf(`{"type":"bridge-session","sessionId":%s,"bridgeSessionId":"bridge-synthetic","lastSequenceNum":0,`+
			`"noHistoryBackfill":false,"ownerAccountUuid":"account-synthetic","ownerOrganizationUuid":"org-synthetic"}`, js(session)))
		if bridge == "off" {
			rows = append(rows, fmt.Sprintf(`{"type":"bridge-session","sessionId":%s,"bridgeSessionId":"","lastSequenceNum":0}`, js(session)))
		}
	}
	agent := os.Getenv("FAKE_AGENT_ID")
	now := time.Now()
	limitRow := fmt.Sprintf(`{"parentUuid":null,"isSidechain":%t,"type":"assistant","uuid":"limit-synthetic","timestamp":%s,`+
		`"message":{"model":"<synthetic>","role":"assistant","type":"message","content":[{"type":"text","text":%s}]},"apiError":"usage_limit_reached",`+
		`"apiErrorParams":{"rate_limit_info":{"status":"rejected","resetsAt":%d,"rateLimitType":"five_hour","overageStatus":"rejected",`+
		`"overageDisabledReason":"out_of_credits","isUsingOverage":false}},"error":"rate_limit","isApiErrorMessage":true,"apiErrorStatus":429,"sessionId":%s}`,
		agent != "", js(now.UTC().Format(time.RFC3339Nano)), js(text), now.Add(after).Unix(), js(session))
	if agent == "" {
		rows = append(rows, limitRow)
	}
	appendRows(conv, rows)
	if agent != "" {
		// Claude Code 2.1.316 keeps a subagent's conversation below its session's.
		appendRows(filepath.Join(strings.TrimSuffix(conv, ".jsonl"), "subagents", "agent-"+agent+".jsonl"), []string{limitRow})
	}

	cwd, _ := os.Getwd()
	in := map[string]any{"session_id": session, "transcript_path": conv, "cwd": cwd, "hook_event_name": "StopFailure",
		"error": "rate_limit", "last_assistant_message": text}
	if agent != "" {
		in["agent_id"], in["agent_type"] = agent, "general-purpose"
	}
	effort := os.Getenv("FAKE_EFFORT")
	if effort != "" {
		in["effort"] = map[string]string{"level": effort}
	}
	input, err := json.Marshal(in)
	must(err)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // the plugin's timeout
	defer cancel()
	hook := exec.CommandContext(ctx, ccam, "hook", "stop-failure")
	hook.Stdin, hook.Stderr = bytes.NewReader(input), os.Stderr
	// Later entries win, so what the developer's own Claude Code exported
	// does not reach the hook. Remote Control shows in the rows alone.
	hook.Env = append(os.Environ(), "CLAUDE_PID="+strconv.Itoa(os.Getpid()), "CLAUDE_CODE_SESSION_ATTENDED=1", "CLAUDE_CODE_ENTRYPOINT=cli",
		"CLAUDE_PROJECT_DIR="+cwd, "CLAUDE_EFFORT="+effort, "CLAUDE_CODE_BRIDGE_SESSION_ID=")
	out, err := hook.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "fakeclaude: ccam hook stop-failure:", err)
	}
	var answer struct {
		TerminalSequence string `json:"terminalSequence"`
	}
	if json.Unmarshal(out, &answer) == nil && answer.TerminalSequence != "" {
		fmt.Println(answer.TerminalSequence)
	}
	time.Sleep(hookGrace)
}

// appendRows adds rows, if any, to the conversation at path.
func appendRows(path string, rows []string) {
	if len(rows) == 0 {
		return
	}
	must(os.MkdirAll(filepath.Dir(path), 0o700))
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	must(err)
	fmt.Fprintln(f, strings.Join(rows, "\n"))
	f.Close()
}

// js is s as a JSON string.
func js(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
