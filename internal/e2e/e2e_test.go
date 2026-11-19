// Package e2e builds the real ccam binary and drives it against a fake
// claude inside a temporary CCAM_HOME.
package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	ccamBin   string
	fakeBin   string // testdata/fakeclaude, the stand-in for claude
	buildOnce sync.Once
	buildErr  error
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "ccam-e2e")
	if err != nil {
		panic(err)
	}
	ccamBin = filepath.Join(dir, "ccam"+exe)
	fakeBin = filepath.Join(dir, "fakeclaude"+exe)
	if raceBuild {
		// A -race binary waits a second before it exits 0 (GORACE
		// atexit_sleep_ms), which over the hundreds of ccam and fake claude
		// runs took the package past go test's 10-minute limit. The
		// developer's own GORACE options come after and win.
		os.Setenv("GORACE", strings.TrimSpace("atexit_sleep_ms=0 "+os.Getenv("GORACE")))
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

// buildCcam runs inside a test (not TestMain) so the files it reads are
// recorded by the go test cache.
func buildCcam(t *testing.T) {
	buildOnce.Do(func() {
		// Binaries built elsewhere, for machines without Go (a Windows host
		// driven from WSL): CCAM_E2E_BIN, CCAM_E2E_FAKE and CCAM_E2E_ROOT.
		if bin, fake := os.Getenv("CCAM_E2E_BIN"), os.Getenv("CCAM_E2E_FAKE"); bin != "" && fake != "" {
			ccamBin, fakeBin = bin, fake
			return
		}
		touchSources()
		for bin, pkg := range map[string]string{ccamBin: "./cmd/ccam", fakeBin: "./testdata/fakeclaude"} {
			// -buildvcs=false: CI containers may refuse git on the checkout (ownership check).
			args := []string{"build", "-buildvcs=false", "-o", bin}
			if raceBuild {
				args = append(args, "-race")
			}
			build := exec.Command("go", append(args, pkg)...)
			build.Dir = repoRoot()
			if out, err := build.CombinedOutput(); err != nil {
				buildErr = fmt.Errorf("build of %s failed: %v\n%s", pkg, err, out)
				return
			}
		}
	})
	if buildErr != nil {
		t.Fatal(buildErr)
	}
}

// touchSources opens every Go source file so the go test cache, which tracks
// files a test reads, invalidates this package when the binary would change.
func touchSources() {
	filepath.WalkDir(repoRoot(), func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && (strings.HasSuffix(path, ".go") || d.Name() == "go.mod" || d.Name() == "go.sum") {
			if f, err := os.Open(path); err == nil {
				f.Close()
			}
		}
		return nil
	})
}

func repoRoot() string {
	if root := os.Getenv("CCAM_E2E_ROOT"); root != "" {
		return root
	}
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..")
}

// systemPath is the PATH after the test's bin folder: system folders only.
var systemPath = map[bool]string{true: `C:\Windows\System32;C:\Windows`, false: "/usr/bin:/bin"}[runtime.GOOS == "windows"]

// perms is whether the platform has Unix permission bits (Windows does not).
var perms = runtime.GOOS != "windows"

// exe is the file name suffix of executables.
var exe = map[bool]string{true: ".exe"}[runtime.GOOS == "windows"]

// sessionEnv are what Claude Code and ccam give the processes of a session,
// which the developer's own may have exported: ccam's hooks and relay read
// them, so no test inherits them.
var sessionEnv = []string{"CCAM_PROFILE", "CCAM_RELAY", "CLAUDE_CONFIG_DIR", "CLAUDE_PID", "CLAUDE_EFFORT", "CLAUDE_PROJECT_DIR",
	"CLAUDE_CODE_SESSION_ATTENDED", "CLAUDE_CODE_ENTRYPOINT", "CLAUDE_CODE_BRIDGE_SESSION_ID", "CLAUDE_CODE_SESSION_ID", "CCAM_AGENT"}

// harness is one isolated ccam installation.
type harness struct {
	t    *testing.T
	home string
	log  string
	env  []string
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	buildCcam(t)
	home := t.TempDir()
	h := &harness{t: t, home: home, log: filepath.Join(home, "claude.log")}
	os.MkdirAll(filepath.Join(home, "config"), 0o700)
	os.WriteFile(filepath.Join(home, "config", "config.toml"), []byte(fmt.Sprintf("claude_bin = %q\n", fakeBin)), 0o600)
	env := slices.DeleteFunc(os.Environ(), func(kv string) bool {
		name, _, _ := strings.Cut(kv, "=")
		return slices.Contains(sessionEnv, name)
	})
	h.env = append(env,
		"CCAM_HOME="+home,
		"FAKE_CLAUDE_LOG="+h.log,
		// Only system folders after bin: shortcut checks look for clashes on
		// PATH, and the developer's own ~/.local/bin must not leak in.
		"PATH="+filepath.Join(home, "bin")+string(os.PathListSeparator)+systemPath,
	)
	return h
}

type result struct {
	stdout, stderr string
	code           int
}

// run executes argv[0] (ccam, ccwork …) with extra env.
func (h *harness) run(stdin string, extraEnv []string, argv ...string) result {
	h.t.Helper()
	return h.runIn("", stdin, extraEnv, argv...)
}

// runIn is run in the working directory dir ("" = the test's).
func (h *harness) runIn(dir, stdin string, extraEnv []string, argv ...string) result {
	h.t.Helper()
	bin := ccamBin
	if argv[0] != "ccam" {
		bin = filepath.Join(h.home, "bin", argv[0]+exe)
	}
	cmd := exec.Command(bin, argv[1:]...)
	cmd.Dir = dir
	cmd.Env = append(append([]string{}, h.env...), extraEnv...)
	cmd.Stdin = strings.NewReader(stdin)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	} else if err != nil {
		h.t.Fatalf("%v: %v", argv, err)
	}
	return result{out.String(), errb.String(), code}
}

func (h *harness) ok(argv ...string) string {
	h.t.Helper()
	r := h.run("", nil, argv...)
	if r.code != 0 {
		h.t.Fatalf("%v exited %d\nstdout: %s\nstderr: %s", argv, r.code, r.stdout, r.stderr)
	}
	return r.stdout
}

// called returns the env lines and args the fake claude received.
func (h *harness) called() (configDir string, args []string) {
	h.t.Helper()
	data, err := os.ReadFile(h.log)
	if err != nil {
		h.t.Fatalf("claude was not called: %v", err)
	}
	os.Remove(h.log)
	parts := strings.Split(strings.TrimSuffix(string(data), "\x00"), "\x00")
	configDir = strings.TrimPrefix(parts[0], "CLAUDE_CONFIG_DIR=")
	return configDir, parts[2:]
}

func (h *harness) profileDir(name string) string {
	return filepath.Join(h.home, "data", "profiles", name)
}

func (h *harness) addAccounts() {
	h.t.Helper()
	for _, a := range []struct{ name, email string }{{"personal", "me@example.com"}, {"work", "dev@company.com"}} {
		r := h.run("", []string{"FAKE_EMAIL=" + a.email}, "ccam", "add", a.name)
		if r.code != 0 || !strings.Contains(r.stdout, a.email) {
			h.t.Fatalf("add %s: %+v", a.name, r)
		}
	}
}

func TestAddLogsInWithoutShortcutByDefault(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	if entries, _ := os.ReadDir(filepath.Join(h.home, "bin")); len(entries) != 0 {
		t.Fatalf("shortcuts are opt-in, found %v", entries)
	}
	out := h.ok("ccam", "list")
	if !strings.Contains(out, "me@example.com") || !strings.Contains(out, "dev@company.com") {
		t.Fatalf("list:\n%s", out)
	}
	info, _ := os.Stat(h.profileDir("work"))
	if perms && info.Mode().Perm() != 0o700 {
		t.Fatalf("profile dir perm %v", info.Mode())
	}
}

// TestArgumentPassThrough covers every row of the PLAN.md argument table.
func TestArgumentPassThrough(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "use", "personal")
	h.ok("ccam", "shortcut", "personal", "ccpersonal")
	h.ok("ccam", "shortcut", "work", "ccwork")
	cases := []struct {
		argv    []string
		profile string
		want    []string
	}{
		{[]string{"ccam", "work", "--resume"}, "work", []string{"--resume"}},
		{[]string{"ccam", "work", "-p", "corrige les tests — ü", "--output-format", "json"}, "work", []string{"-p", "corrige les tests — ü", "--output-format", "json"}},
		{[]string{"ccwork", "--model", "opus", "--continue"}, "work", []string{"--model", "opus", "--continue"}},
		{[]string{"ccam", "work", "--help"}, "work", []string{"--help"}},
		{[]string{"ccam", "-p", "résume ça"}, "personal", []string{"-p", "résume ça"}},
		{[]string{"ccam", "work", "--", "--version"}, "work", []string{"--version"}},
		{[]string{"ccam", "work", "--version"}, "work", []string{"--version"}},
		{[]string{"ccam", "2"}, "work", nil},
		{[]string{"ccam", "run", "work", "--", "-c"}, "work", []string{"-c"}},
		{[]string{"ccpersonal", "it's \"quoted\"", ""}, "personal", []string{"it's \"quoted\"", ""}},
	}
	for _, c := range cases {
		h.ok(c.argv...)
		dir, args := h.called()
		if dir != h.profileDir(c.profile) {
			t.Errorf("%q: config dir %s", c.argv, dir)
		}
		if !(len(args) == 0 && len(c.want) == 0) && !reflect.DeepEqual(args, c.want) {
			t.Errorf("%q: args %q want %q", c.argv, args, c.want)
		}
	}
}

func TestStdinAndExitCodePassThrough(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	stdinFile := filepath.Join(h.home, "stdin.txt")
	r := h.run("log line\n", []string{"FAKE_CLAUDE_STDIN=" + stdinFile, "FAKE_EXIT=7"}, "ccam", "work", "-p", "trouve la cause")
	if r.code != 7 {
		t.Fatalf("exit code %d", r.code)
	}
	if b, _ := os.ReadFile(stdinFile); string(b) != "log line\n" {
		t.Fatalf("stdin %q", b)
	}
}

func TestDryRunAndErrors(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	os.Remove(h.log) // left behind by `auth login`
	out := h.ok("ccam", "--dry-run", "work", "-c")
	want, args := "CLAUDE_CONFIG_DIR="+h.profileDir("work"), " -c\n"
	if runtime.GOOS == "windows" {
		want, args = "$env:CLAUDE_CONFIG_DIR='"+h.profileDir("work")+"';", " -ArgumentList '-c' "
	}
	if !strings.Contains(out, want) || !strings.Contains(out, args) {
		t.Fatalf("dry run: %s", out)
	}
	if _, err := os.Stat(h.log); err == nil {
		t.Fatal("dry run must not exec")
	}
	r := h.run("", nil, "ccam", "nobody")
	if r.code == 0 || !strings.Contains(r.stderr, "nobody") {
		t.Fatalf("unknown profile: %+v", r)
	}
	r = h.run("", nil, "ccam", "add", "list")
	if r.code == 0 {
		t.Fatal("reserved name accepted")
	}
}

// A claude_bin that points nowhere is named, and the account is not marked
// as just used: claude never ran.
func TestMissingClaudeBinIsNamed(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	cfgFile := filepath.Join(h.home, "config", "config.toml")
	data, _ := os.ReadFile(cfgFile)
	missing := filepath.Join(h.home, "gone", "claude"+exe)
	os.WriteFile(cfgFile, bytes.Replace(data, []byte(fakeBin), []byte(missing), 1), 0o600)
	os.RemoveAll(filepath.Join(h.home, "state", "last-used"))
	r := h.run("", nil, "ccam", "work", "-p", "x")
	if r.code == 0 || !strings.Contains(r.stderr, missing) || !strings.Contains(r.stderr, "claude_bin") {
		t.Fatalf("%+v", r)
	}
	if _, err := os.Stat(filepath.Join(h.home, "state", "last-used", "work")); err == nil {
		t.Fatal("marked as used although claude did not start")
	}
}

// TestDryRunNeverChangesState: a command that cannot honour --dry-run must
// refuse it rather than do the real thing.
func TestDryRunNeverChangesState(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	cfgFile := filepath.Join(h.home, "config", "config.toml")
	before, _ := os.ReadFile(cfgFile)
	for _, argv := range [][]string{
		{"remove", "work", "--yes"},
		{"add", "ghost", "--no-login"},
		{"rename", "work", "job"},
		{"use", "work"},
		{"logout", "work"},
		{"login", "work"},
		{"mcp", "add", "fs", "--", "npx", "server"},
		{"import", "--yes"},
		{"shims", "install"},
		{"statusline", "install"},
		{"usage", "--refresh"},
	} {
		r := h.run("", nil, append([]string{"ccam", "--dry-run"}, argv...)...)
		if r.code == 0 || !strings.Contains(r.stderr, "dry-run") {
			t.Errorf("--dry-run %v: %+v", argv, r)
		}
	}
	if after, _ := os.ReadFile(cfgFile); !bytes.Equal(before, after) {
		t.Fatalf("config changed:\n%s\n→\n%s", before, after)
	}
	if _, err := os.Stat(filepath.Join(h.home, "data", "groups", "default", "mcp.json")); err == nil {
		t.Fatal("mcp.json written")
	}
	if _, err := os.Stat(h.profileDir("work")); err != nil {
		t.Fatal("profile folder gone")
	}
	// Read-only commands still accept the flag; doctor --fix only reports,
	// with what --fix would set aside.
	h.ok("ccam", "--dry-run", "list")
	h.ok("ccam", "--dry-run", "doctor", "--fix")
}

func TestRemoveAndRename(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "shortcut", "work", "ccwork")
	h.ok("ccam", "rename", "work", "job")
	// On macOS the folder keeps its name: the Keychain names the login
	// after its path.
	dir := h.profileDir("job")
	if runtime.GOOS == "darwin" {
		dir = h.profileDir("work")
	}
	h.ok("ccwork")
	if got, _ := h.called(); got != dir {
		t.Fatalf("shortcut after rename: %s", got)
	}
	h.ok("ccam", "job")
	if got, _ := h.called(); got != dir {
		t.Fatalf("rename dir %s", got)
	}
	r := h.run("n\n", nil, "ccam", "remove", "job")
	if r.code != 0 || !strings.Contains(r.stdout, "Cancelled") {
		t.Fatalf("declined remove: %+v", r)
	}
	h.run("y\n", nil, "ccam", "remove", "job")
	if _, err := os.Stat(dir); err == nil {
		t.Fatal("profile dir still there")
	}
	if _, err := os.Lstat(filepath.Join(h.home, "bin", "ccwork"+exe)); err == nil {
		t.Fatal("ccwork shortcut still there")
	}
	trash, _ := os.ReadDir(filepath.Join(h.home, "data", "trash"))
	if len(trash) != 1 {
		t.Fatalf("trash %v", trash)
	}
}

func TestLoginExpiredLogsInOnlyThoseThatNeedIt(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "add", "later", "--no-login")
	// work's refresh token has expired; personal's login works.
	os.WriteFile(filepath.Join(h.profileDir("work"), ".credentials.json"),
		[]byte(`{"claudeAiOauth":{"accessToken":"x","expiresAt":1,"refreshTokenExpiresAt":1}}`), 0o600)
	r := h.run("", []string{"FAKE_EMAIL=new@example.com"}, "ccam", "login", "--expired")
	if r.code != 0 || !strings.Contains(r.stdout, "[1/2] work (login expired)") || !strings.Contains(r.stdout, "[2/2] later (not logged in)") ||
		strings.Contains(r.stdout, "personal") {
		t.Fatalf("%+v", r)
	}
	if out := h.ok("ccam", "login", "--expired"); !strings.Contains(out, "Every account's login works.") {
		t.Fatal(out)
	}
	if r := h.run("", nil, "ccam", "login", "--expired", "work"); r.code == 0 {
		t.Fatal("--expired takes no name")
	}
}

// Claude Code takes the email only as a hint, so ccam checks whom the login
// belongs to.
func TestLoginIntoAnotherAccountIsCaught(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	// A usage limit Claude Code reported for work's account.
	limit := filepath.Join(h.home, "state", "usage", "work.limit.json")
	os.MkdirAll(filepath.Dir(limit), 0o700)
	os.WriteFile(limit, []byte(`{"at":"2026-10-22T10:00:00Z","windows":{"seven_day":"2100-01-01T00:00:00Z"}}`), 0o600)
	args := filepath.Join(h.home, "login-args")
	r := h.run("", []string{"FAKE_EMAIL=dev@company.com", "FAKE_LOGIN_ARGS=" + args}, "ccam", "login", "work")
	if b, _ := os.ReadFile(args); r.code != 0 || !strings.Contains(string(b), "--email dev@company.com") {
		t.Fatalf("%+v %s", r, b)
	}
	if _, err := os.Stat(limit); err != nil {
		t.Fatal("logging in to the same account again dropped its limit")
	}
	r = h.run("", []string{"FAKE_EMAIL=other@company.com"}, "ccam", "login", "work")
	if !strings.Contains(r.stderr, "this login is other@company.com, but work was dev@company.com before") ||
		!strings.Contains(r.stderr, "run `ccam login work` again and sign in as dev@company.com") || strings.Contains(r.stderr, "logout") {
		t.Fatalf("%+v", r)
	}
	if _, err := os.Stat(limit); err == nil {
		t.Fatal("the other account's limit stayed with work")
	}
	// Two folders logged in to one account share their limits.
	r = h.run("", []string{"FAKE_EMAIL=me@example.com"}, "ccam", "login", "work")
	if !strings.Contains(r.stderr, "personal is logged in to the same account") {
		t.Fatalf("%+v", r)
	}
	cfg, _ := os.ReadFile(filepath.Join(h.home, "config", "config.toml"))
	if !strings.Contains(string(cfg), "[profile.login]") || !strings.Contains(string(cfg), "org_type = 'claude_team'") {
		t.Fatalf("%s", cfg)
	}
}

// An account an older ccam added has no login record, only the email it
// recorded: a login into another Claude account is caught all the same.
func TestLoginOfAnAccountWithoutARecordIsChecked(t *testing.T) {
	h := newHarness(t)
	os.MkdirAll(h.profileDir("work"), 0o700)
	cfg := fmt.Sprintf("claude_bin = %q\n\n[[profile]]\nname = 'work'\nslot = 1\nemail = 'dev@company.com'\n", fakeBin)
	os.WriteFile(filepath.Join(h.home, "config", "config.toml"), []byte(cfg), 0o600)
	r := h.run("", []string{"FAKE_EMAIL=other@company.com"}, "ccam", "login", "work")
	if r.code != 0 || !strings.Contains(r.stderr, "this login is other@company.com, but work was dev@company.com before") {
		t.Fatalf("%+v", r)
	}
	r = h.run("", []string{"FAKE_EMAIL=other@company.com"}, "ccam", "login", "work")
	if strings.Contains(r.stderr, "but work was") {
		t.Fatalf("the same login again: %+v", r)
	}
}

// A setting added at the end of config.toml belongs to the last table in
// TOML, so ccam does not read it: a launch says so, and a command that
// saves the config refuses instead of deleting the line.
func TestSettingBelowATableIsNamedAndKept(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	path := filepath.Join(h.home, "config", "config.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, "\nrelay = true\n"...)
	os.WriteFile(path, data, 0o600)
	r := h.run("", nil, "ccam", "work", "--version")
	if r.code != 0 || !strings.Contains(r.stderr, "ccam: config.toml line ") || !strings.Contains(r.stderr, "relay is read as part of the [profile.login] table above it") {
		t.Fatalf("launch: %+v", r)
	}
	if dir, _ := h.called(); dir != h.profileDir("work") {
		t.Fatal("claude did not run")
	}
	r = h.run("", nil, "ccam", "use", "work")
	if r.code == 0 || !strings.Contains(r.stderr, "so it was not saved") {
		t.Fatalf("use: %+v", r)
	}
	if after, _ := os.ReadFile(path); string(after) != string(data) {
		t.Fatalf("config.toml changed:\n%s", after)
	}
	if out := h.run("", nil, "ccam", "doctor").stdout; !strings.Contains(out, "relay is read as part of") {
		t.Fatalf("doctor:\n%s", out)
	}
}

func TestAgentsAcrossAccounts(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	agents := filepath.Join(h.home, "agents")
	os.MkdirAll(agents, 0o700)
	now := time.Now().UnixMilli()
	os.WriteFile(filepath.Join(agents, "personal.json"), []byte(fmt.Sprintf(
		`[{"pid":11,"id":"49ec0028","cwd":"/tmp/app","kind":"background","startedAt":%d,"sessionId":"49ec0028-x","name":"probe","status":"idle","state":"blocked"}]`, now-60000)), 0o600)
	os.WriteFile(filepath.Join(agents, "work.json"), []byte(fmt.Sprintf(
		`[{"pid":12,"cwd":"/tmp/web","kind":"interactive","startedAt":%d,"status":"busy"}]`, now)), 0o600)
	env := []string{"FAKE_AGENTS_DIR=" + agents}
	r := h.run("", env, "ccam", "agents")
	lines := strings.Split(r.stdout, "\n")
	if r.code != 0 || len(lines) < 3 || !strings.HasPrefix(lines[1], "work") || !strings.Contains(lines[1], "open (busy)") ||
		!strings.HasPrefix(lines[2], "personal") || !strings.Contains(lines[2], "49ec0028") || !strings.Contains(lines[2], "blocked, idle") {
		t.Fatalf("%+v", r)
	}
	if r := h.run("", env, "ccam", "agents", "--json"); !strings.Contains(r.stdout, `"account": "personal"`) {
		t.Fatalf("%+v", r)
	}
	// An account that cannot be asked is not an empty list.
	broken := t.TempDir()
	for _, n := range []string{"personal", "work"} {
		os.WriteFile(filepath.Join(broken, n+".json"), []byte("not json"), 0o600)
	}
	if r := h.run("", []string{"FAKE_AGENTS_DIR=" + broken}, "ccam", "agents"); r.code == 0 || strings.Contains(r.stdout, "No sessions") ||
		!strings.Contains(r.stderr, "2 of 2 accounts could not be asked") {
		t.Fatalf("%+v", r)
	}
	os.MkdirAll(filepath.Join(h.profileDir("personal"), "jobs", "49ec0028"), 0o700)
	r = h.run("", nil, "ccam", "--dry-run", "attach", "49ec")
	if r.code != 0 || !strings.Contains(r.stdout, filepath.Join("profiles", "personal")) || !strings.Contains(r.stdout, "attach 49ec0028") {
		t.Fatalf("%+v", r)
	}
	if r := h.run("", nil, "ccam", "attach", "nope"); r.code == 0 || !strings.Contains(r.stderr, "no account has a background session") {
		t.Fatalf("%+v", r)
	}
	// Flags after the id are claude's, and part of the session name works.
	r = h.run("", nil, "ccam", "--dry-run", "attach", "49ec", "--debug")
	if r.code != 0 || !strings.Contains(r.stdout, "attach 49ec0028 --debug") {
		t.Fatalf("%+v", r)
	}
	os.WriteFile(filepath.Join(h.profileDir("personal"), "jobs", "49ec0028", "state.json"), []byte(`{"name":"Fix Login"}`), 0o600)
	os.MkdirAll(filepath.Join(h.profileDir("work"), "jobs", "77ab0011"), 0o700)
	os.WriteFile(filepath.Join(h.profileDir("work"), "jobs", "77ab0011", "state.json"), []byte(`{"name":"fix docs"}`), 0o600)
	r = h.run("", nil, "ccam", "--dry-run", "attach", "LOGIN")
	if r.code != 0 || !strings.Contains(r.stdout, filepath.Join("profiles", "personal")) || !strings.Contains(r.stdout, "attach 49ec0028") {
		t.Fatalf("%+v", r)
	}
	if r := h.run("", nil, "ccam", "attach", "fix"); r.code == 0 || !strings.Contains(r.stderr, "fits more than one session") {
		t.Fatalf("%+v", r)
	}
	// What may start an id is not read as a name; an empty one fits nothing.
	for _, sel := range []string{"dead", ""} {
		if r := h.run("", nil, "ccam", "attach", sel); r.code == 0 || !strings.Contains(r.stderr, "no account has a background session") {
			t.Fatalf("%q: %+v", sel, r)
		}
	}
}

// With usage_source = "claude", Claude Code looks the usage up itself.
func TestUsageThroughClaudeCode(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	dir := filepath.Join(h.home, "usage")
	os.MkdirAll(dir, 0o700)
	os.WriteFile(filepath.Join(dir, "work.json"), []byte(`{"type":"control_response","response":{"subtype":"success","request_id":"ccam","response":{"rate_limits_available":true,`+
		`"rate_limits":{"five_hour":{"utilization":37,"resets_at":"2099-01-01T00:00:00Z"},"seven_day":{"utilization":12,"resets_at":"2099-01-05T00:00:00Z"}}}}}`), 0o600)
	env := []string{"FAKE_USAGE_DIR=" + dir}
	r := h.run("", env, "ccam", "usage", "--source", "claude", "-v")
	if r.code != 0 || !strings.Contains(r.stdout, "37%") || !strings.Contains(r.stdout, "from Claude Code") {
		t.Fatalf("%+v", r)
	}
	// personal's login has no plan limits.
	if !strings.Contains(r.stdout, "no plan limits") {
		t.Fatalf("%+v", r)
	}
	if r := h.run("", nil, "ccam", "usage", "--source", "bogus"); r.code == 0 || !strings.Contains(r.stderr, "live, claude or cache") {
		t.Fatalf("%+v", r)
	}

	// A dry run goes by the saved usage: Claude Code would write to the folder.
	cfgFile := filepath.Join(h.home, "config", "config.toml")
	cfg, _ := os.ReadFile(cfgFile)
	os.WriteFile(cfgFile, append([]byte("usage_source = \"claude\"\n"), cfg...), 0o600)
	os.RemoveAll(filepath.Join(h.home, "state", "usage"))
	os.Remove(h.log)
	r = h.run("", env, "ccam", "--dry-run", "usage", "-v")
	if _, err := os.Stat(h.log); err == nil || strings.Contains(r.stdout, "37%") || !strings.Contains(r.stdout, "a dry run does not start Claude Code") {
		t.Fatalf("Claude Code ran: %+v", r)
	}
}
