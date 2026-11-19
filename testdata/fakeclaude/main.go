// Command fakeclaude stands in for claude in the end-to-end tests. It
// records its environment and arguments (NUL separated) in FAKE_CLAUDE_LOG
// and fakes `auth login` by writing an oauthAccount into .claude.json and a
// login into .credentials.json (`auth logout` deletes that file).
// The other FAKE_* variables make it misbehave on purpose. Like Claude
// Code, it ends with 143 on SIGTERM, and with FAKE_RECORD it lists itself
// in its config folder while it runs (sessions/<pid>.json). With FAKE_HOOK
// it stops a turn at a usage limit and runs the real ccam hook (hook.go).
package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	unlist := func() {}
	if kind := os.Getenv("FAKE_RECORD"); kind != "" {
		unlist = listSelf(os.Getenv("CLAUDE_CONFIG_DIR"), kind)
	}
	if s := os.Getenv("FAKE_ONLY_SLEEP"); s != "" {
		sleep(s)
		unlist()
		return
	}
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGTERM)
	go func() {
		<-term
		if path := os.Getenv("FAKE_TERMS"); path != "" {
			// The accounts whose runs SIGTERM ended, in order.
			if f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600); err == nil {
				fmt.Fprintln(f, os.Getenv("CCAM_PROFILE"))
				f.Close()
			}
		}
		os.Exit(143)
	}()
	args := ucrtArgs()
	dir := os.Getenv("CLAUDE_CONFIG_DIR")
	profile := os.Getenv("CCAM_PROFILE")
	var log strings.Builder
	fmt.Fprintf(&log, "CLAUDE_CONFIG_DIR=%s\x00CCAM_PROFILE=%s\x00", dir, profile)
	for _, a := range args {
		log.WriteString(a + "\x00")
	}
	path := os.Getenv("FAKE_CLAUDE_LOG")
	if path == "" {
		fmt.Fprintln(os.Stderr, "FAKE_CLAUDE_LOG not set")
		os.Exit(2)
	}
	must(os.WriteFile(path, []byte(log.String()), 0o600))
	if s := os.Getenv("FAKE_DETACH"); s != "" && slices.Contains(args, "--bg") {
		// Like `claude --bg`: the session goes on in a process that outlives
		// this one, listed in the folder as Claude Code's daemon and the
		// session it runs, and this one returns once that runs (daemon.lock
		// names it; listSelf writes it last).
		self, err := os.Executable()
		must(err)
		bg := exec.Command(self)
		bg.Env = append(os.Environ(), "FAKE_ONLY_SLEEP="+s, "FAKE_RECORD=bg", "FAKE_DAEMON=1")
		must(bg.Start())
		named := fmt.Sprintf(`"pid":%d,`, bg.Process.Pid)
		for i := 0; i < 500; i++ {
			if data, _ := os.ReadFile(filepath.Join(dir, "daemon.lock")); strings.Contains(string(data), named) {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		fmt.Println("Started background session 49ec0028")
		os.Exit(0)
	}
	if path := os.Getenv("FAKE_ENV_FILE"); path != "" {
		// NAME=value lines for the variables named in FAKE_ENV_NAMES.
		var b strings.Builder
		for _, name := range strings.Split(os.Getenv("FAKE_ENV_NAMES"), ",") {
			if v, ok := os.LookupEnv(name); ok {
				fmt.Fprintf(&b, "%s=%s\n", name, v)
			}
		}
		must(os.WriteFile(path, []byte(b.String()), 0o600))
	}
	if runs := os.Getenv("FAKE_RUNS"); runs != "" {
		f, err := os.OpenFile(runs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
		must(err)
		fmt.Fprintln(f, profile)
		f.Close()
	}
	stopped := false
	if from := os.Getenv("FAKE_SWITCH_FROM"); from != "" && from == profile {
		// What `ccam hook switch` leaves for /ccam:switch: the note of
		// session FAKE_SWITCH_SESSION (the conversation's name by default).
		// With FAKE_SWITCH_STOPPED it says the hook ends this claude, and
		// it ends with SIGTERM as the hook makes it, unless FAKE_EXIT says
		// how it ends instead (the person's /exit was first).
		conv := os.Getenv("FAKE_SWITCH_TRANSCRIPT")
		session := cmp.Or(os.Getenv("FAKE_SWITCH_SESSION"), strings.TrimSuffix(filepath.Base(conv), ".jsonl"))
		stopped = os.Getenv("FAKE_SWITCH_STOPPED") != ""
		child := 0
		if stopped {
			child = os.Getpid()
		}
		if mode := os.Getenv("FAKE_MODE"); mode != "" {
			// The permission mode the conversation was left in.
			f, err := os.OpenFile(conv, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
			must(err)
			fmt.Fprintf(f, `{"type":"permission-mode","permissionMode":%q,"sessionId":%q}`+"\n", mode, session)
			f.Close()
		}
		cwd, _ := os.Getwd()
		ticket := fmt.Sprintf(`{"account":%q,"session_id":%q,"transcript_path":%q,"cwd":%q,"at":%q,"reason":"switch","target":%q,"relay":%q,"child":%d,"stopped":%v}`,
			from, session, conv, cwd, time.Now().Format(time.RFC3339Nano), os.Getenv("FAKE_SWITCH_TO"),
			cmp.Or(os.Getenv("FAKE_SWITCH_RELAY"), os.Getenv("CCAM_RELAY")), child, stopped)
		dir := filepath.Join(os.Getenv("FAKE_SWITCH_STATE"), "handoff", from)
		must(os.MkdirAll(dir, 0o700))
		must(os.WriteFile(filepath.Join(dir, session+".json"), []byte(ticket), 0o600))
		if name := os.Getenv("FAKE_KILL_PARENT"); name != "" {
			// A signal to the relaying ccam: HUP, INT or TERM.
			sig := map[string]os.Signal{"HUP": syscall.SIGHUP, "INT": os.Interrupt, "TERM": syscall.SIGTERM}[name]
			if p, err := os.FindProcess(os.Getppid()); err == nil && sig != nil {
				must(p.Signal(sig))
			}
		}
	}
	if path := os.Getenv("FAKE_TTY_INPUT"); path != "" && profile != os.Getenv("FAKE_SWITCH_FROM") {
		// What waits on the terminal of the session that does not switch,
		// read for 0.2 s out of line mode.
		c := exec.Command("sh", "-c", "stty -icanon min 0 time 2; cat; stty icanon")
		c.Stdin = os.Stdin
		out, _ := c.Output()
		must(os.WriteFile(path, out, 0o600))
	}
	if slices.Contains(args, "--input-format") && slices.Contains(args, "stream-json") {
		// The get_usage control request: answer from FAKE_USAGE_DIR/<account>.json.
		in, _ := io.ReadAll(os.Stdin)
		if strings.Contains(string(in), "get_usage") {
			data, err := os.ReadFile(filepath.Join(os.Getenv("FAKE_USAGE_DIR"), profile+".json"))
			if err != nil {
				data = []byte(`{"type":"control_response","response":{"subtype":"success","request_id":"ccam","response":{"rate_limits_available":false,"rate_limits":null}}}`)
			}
			fmt.Println(`{"type":"system","subtype":"init"}`)
			os.Stdout.Write(append(data, '\n'))
		}
		return
	}
	if len(args) >= 2 && args[0] == "agents" && args[1] == "--json" {
		// FAKE_AGENTS_DIR/<account>.json holds what the account lists.
		data, err := os.ReadFile(filepath.Join(os.Getenv("FAKE_AGENTS_DIR"), profile+".json"))
		if err != nil {
			data = []byte("[]")
		}
		os.Stdout.Write(data)
		return
	}
	if len(args) >= 2 && args[0] == "auth" && args[1] == "login" {
		email := os.Getenv("FAKE_EMAIL")
		if email == "" {
			email = "user@example.com"
		}
		// The account is named after the email unless FAKE_ACCOUNT says
		// otherwise; FAKE_LOGIN_ARGS records the arguments.
		account := os.Getenv("FAKE_ACCOUNT")
		if account == "" {
			account = "uuid-" + email
		}
		if path := os.Getenv("FAKE_LOGIN_ARGS"); path != "" {
			must(os.WriteFile(path, []byte(strings.Join(args, " ")), 0o600))
		}
		// A seat of Team organization org-1 unless FAKE_ORG_UUID and
		// FAKE_ORG_TYPE say otherwise.
		must(os.WriteFile(filepath.Join(dir, ".claude.json"), fmt.Appendf(nil,
			`{"numStartups":1,"oauthAccount":{"emailAddress":%q,"accountUuid":%q,"organizationUuid":%q,"organizationName":"Acme","organizationType":%q}}`,
			email, account, cmp.Or(os.Getenv("FAKE_ORG_UUID"), "org-1"), cmp.Or(os.Getenv("FAKE_ORG_TYPE"), "claude_team")), 0o600))
		// Like Claude Code, a login replaces only the login in the login
		// file; the MCP sign-ins there stay.
		cred := filepath.Join(dir, ".credentials.json")
		doc := map[string]json.RawMessage{}
		if data, err := os.ReadFile(cred); err == nil {
			json.Unmarshal(data, &doc)
		}
		doc["claudeAiOauth"] = json.RawMessage(`{"accessToken":"fake","expiresAt":4102444800000}`)
		data, _ := json.Marshal(doc)
		must(os.WriteFile(cred, data, 0o600))
	}
	if len(args) >= 2 && args[0] == "auth" && args[1] == "logout" {
		// Like Claude Code, a logout deletes the login file, MCP sign-ins
		// and all.
		must(os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(`{}`), 0o600))
		if err := os.Remove(filepath.Join(dir, ".credentials.json")); err != nil && !os.IsNotExist(err) {
			must(err)
		}
	}
	if probe := os.Getenv("FAKE_SESSION_PROBE"); probe != "" {
		// With exec the marker carries this process's PID; where ccam waits
		// for claude instead (Windows), it carries ccam's, the parent's.
		held := "free"
		for _, pid := range []int{os.Getpid(), os.Getppid()} {
			if _, err := os.Stat(filepath.Join(probe, strconv.Itoa(pid))); err == nil {
				held = "held"
			}
		}
		must(os.WriteFile(probe+".result", []byte(held+"\n"), 0o600))
	}
	if path := os.Getenv("FAKE_CLAUDE_STDIN"); path != "" {
		data, err := io.ReadAll(os.Stdin)
		must(err)
		must(os.WriteFile(path, data, 0o600))
	}
	if os.Getenv("FAKE_HOOK") != "" && profile != "" && slices.Contains(strings.Split(os.Getenv("FAKE_HOOK_FROM"), ","), profile) {
		stopFailure()
	}
	if s := os.Getenv("FAKE_SLEEP"); s != "" {
		sleep(s)
	}
	if s := os.Getenv("FAKE_BG"); s != "" {
		// A child that outlives claude, like a Bash-tool background job.
		self, err := os.Executable()
		must(err)
		bg := exec.Command(self)
		bg.Env = append(os.Environ(), "FAKE_ONLY_SLEEP="+s)
		must(bg.Start())
	}
	if jsonPrint(args) && after(args, "--session-id") != "" {
		printJSON(args, dir, profile) // ccam agent's run
	}
	if slices.Contains(strings.Split(os.Getenv("FAKE_LIMITED"), ","), profile) && profile != "" {
		fmt.Println(cmp.Or(os.Getenv("FAKE_LIMIT_TEXT"), "You've hit your limit · resets 3pm"))
		os.Exit(1)
	}
	if out := os.Getenv("FAKE_STDOUT"); out != "" {
		fmt.Println(out + " (" + profile + ")")
	}
	if s := os.Getenv("FAKE_EXIT_DELAY"); s != "" {
		sleep(s) // a SIGTERM meanwhile ends it with 143
	}
	if stopped && os.Getenv("FAKE_EXIT") == "" {
		if p, err := os.FindProcess(os.Getpid()); err == nil && p.Signal(syscall.SIGTERM) == nil {
			time.Sleep(10 * time.Second) // the handler exits
		}
	}
	code, _ := strconv.Atoi(os.Getenv("FAKE_EXIT"))
	unlist()
	os.Exit(code)
}

// listSelf lists this process where Claude Code lists its own while they
// run: sessions/<pid>.json in its config folder dir, and daemon.lock for the
// daemon (FAKE_DAEMON), with the start time as Claude Code records it. It
// returns what removes them again, as Claude Code does when it exits.
func listSelf(dir, kind string) func() {
	pid, start := os.Getpid(), procStart()
	session := filepath.Join(dir, "sessions", strconv.Itoa(pid)+".json")
	must(os.MkdirAll(filepath.Dir(session), 0o700))
	must(os.WriteFile(session, fmt.Appendf(nil, `{"pid":%d,"sessionId":"fake-%d","kind":%q,"procStart":%q}`, pid, pid, kind, start), 0o644))
	files := []string{session}
	if os.Getenv("FAKE_DAEMON") != "" {
		key := "procStart"
		if runtime.GOOS == "windows" {
			key = "procStartFt"
		}
		lock := filepath.Join(dir, "daemon.lock")
		must(os.WriteFile(lock, fmt.Appendf(nil, `{"pid":%d,"version":"2.1.308","origin":"transient",%q:%q}`, pid, key, start), 0o644))
		files = append(files, lock)
	}
	return func() {
		for _, f := range files {
			os.Remove(f)
		}
	}
}

func sleep(s string) {
	n, _ := strconv.ParseFloat(s, 64)
	time.Sleep(time.Duration(n * float64(time.Second)))
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "fakeclaude:", err)
		os.Exit(2)
	}
}
