package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

// printJSON fakes `claude -p --output-format json` as ccam agent runs it.
// It writes the conversation under --session-id in the config folder's
// projects/ (a copy of the --resume file first, as --fork-session makes),
// adds a line to FAKE_AGENT_LOG (account, then the arguments), and prints
// the result object. An account in FAKE_LIMITED gets Claude Code's limit
// row (five_hour, resetting in FAKE_LIMIT_RESET, 1h by default) and an
// error result instead of FAKE_STDOUT's answer. It exits with the code
// claude would.
func printJSON(args []string, dir, profile string) {
	id := after(args, "--session-id")
	cwd, err := os.Getwd()
	must(err)
	if real, err := filepath.EvalSymlinks(cwd); err == nil {
		cwd = real
	}
	path := filepath.Join(dir, "projects", slug(cwd), id+".jsonl")
	must(os.MkdirAll(filepath.Dir(path), 0o700))
	var conv []byte
	if from := after(args, "--resume"); from != "" {
		if !slices.Contains(args, "--fork-session") {
			fmt.Fprintln(os.Stderr, "fakeclaude: --session-id with --resume also needs --fork-session")
			os.Exit(1)
		}
		conv, err = os.ReadFile(from)
		must(err)
	}
	prompt := ""
	if i := slices.Index(args, "--"); i >= 0 && i+1 < len(args) {
		prompt = args[i+1]
	}
	now := time.Now().UTC()
	conv = fmt.Appendf(conv, `{"type":"user","message":{"role":"user","content":%s},"permissionMode":%s,"sessionId":%s,"timestamp":%s}`+"\n",
		js(prompt), js(cmp.Or(after(args, "--permission-mode"), "default")), js(id), js(now.Format(time.RFC3339Nano)))
	limited := slices.Contains(strings.Split(os.Getenv("FAKE_LIMITED"), ","), profile)
	text := cmp.Or(os.Getenv("FAKE_STDOUT"), "done") + " (" + profile + ")"
	if limited {
		text = cmp.Or(os.Getenv("FAKE_LIMIT_TEXT"), "You've hit your session limit · resets 3pm")
		reset, err := time.ParseDuration(cmp.Or(os.Getenv("FAKE_LIMIT_RESET"), "1h"))
		must(err)
		conv = fmt.Appendf(conv, `{"isSidechain":false,"type":"assistant","timestamp":%s,`+
			`"message":{"model":"<synthetic>","role":"assistant","type":"message","content":[{"type":"text","text":%s}]},"apiError":"usage_limit_reached",`+
			`"apiErrorParams":{"rate_limit_info":{"status":"rejected","resetsAt":%d,"rateLimitType":"five_hour"}},"error":"rate_limit","isApiErrorMessage":true,"apiErrorStatus":429,"sessionId":%s}`+"\n",
			js(now.Format(time.RFC3339Nano)), js(text), now.Add(reset).Unix(), js(id))
	} else {
		conv = fmt.Appendf(conv, `{"isSidechain":false,"type":"assistant","timestamp":%s,"message":{"model":"claude-opus-5-5","role":"assistant","content":[{"type":"text","text":%s}]},"sessionId":%s}`+"\n",
			js(now.Format(time.RFC3339Nano)), js(text), js(id))
	}
	must(os.WriteFile(path, conv, 0o600))
	if log := os.Getenv("FAKE_AGENT_LOG"); log != "" {
		f, err := os.OpenFile(log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
		must(err)
		fmt.Fprintf(f, "%s\t%s\n", profile, strings.Join(args, " "))
		f.Close()
	}
	out, _ := json.Marshal(map[string]any{"type": "result", "subtype": "success", "is_error": limited, "result": text, "session_id": id})
	fmt.Println(string(out))
	if limited {
		os.Exit(1)
	}
	code, _ := strconv.Atoi(os.Getenv("FAKE_EXIT"))
	os.Exit(code)
}

// jsonPrint reports whether args run claude -p with JSON output.
func jsonPrint(args []string) bool {
	return slices.Contains(args, "-p") && after(args, "--output-format") == "json"
}

// after is the argument after flag in args, "" when there is none.
func after(args []string, flag string) string {
	for i, a := range args {
		if a == "--" {
			break
		}
		if a == flag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

// slug is the folder Claude Code keeps a working folder's conversations
// in under projects/ (ccam's transcript.Slug, for short paths).
func slug(dir string) string {
	var b strings.Builder
	for _, u := range utf16.Encode([]rune(dir)) {
		if u < 0x80 && (u >= 'a' && u <= 'z' || u >= 'A' && u <= 'Z' || u >= '0' && u <= '9') {
			b.WriteByte(byte(u))
		} else {
			b.WriteByte('-')
		}
	}
	return b.String()
}
