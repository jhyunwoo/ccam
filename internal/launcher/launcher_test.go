package launcher

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/session"
)

// exeSuffix makes a file found by PATH lookups: Windows looks for .exe.
var exeSuffix = map[bool]string{true: ".exe"}[runtime.GOOS == "windows"]

func setup(t testing.TB) (paths.Paths, *config.Config) {
	root := t.TempDir()
	os.Setenv("CCAM_HOME", root)
	t.Cleanup(func() { os.Unsetenv("CCAM_HOME") })
	p := paths.FromEnv()
	fake := filepath.Join(root, "fake-claude")
	os.WriteFile(fake, []byte("#!/bin/sh\n"), 0o755)
	os.MkdirAll(p.ProfileDir("personal"), 0o700)
	os.MkdirAll(p.ProfileDir("work"), 0o700)
	cfg := &config.Config{Default: "work", ClaudeBin: fake, Profiles: []config.Profile{
		{Name: "personal", Slot: 1},
		{Name: "work", Slot: 2},
	}}
	return p, cfg
}

func TestPrepareSelectsProfile(t *testing.T) {
	p, cfg := setup(t)
	cases := []struct {
		sel     string
		args    []string
		profile string
		argv    []string
	}{
		{"work", []string{"--resume"}, "work", []string{"claude", "--resume"}},
		{"1", []string{"--model", "opus", "--continue"}, "personal", []string{"claude", "--model", "opus", "--continue"}},
		{"", []string{"-p", "résume ça"}, "work", []string{"claude", "-p", "résume ça"}},
		{"work", []string{"--", "--version"}, "work", []string{"claude", "--version"}},
		{"work", []string{"--help"}, "work", []string{"claude", "--help"}},
		{"work", []string{"-p", "corrige les tests — ü", "--output-format", "json"}, "work", []string{"claude", "-p", "corrige les tests — ü", "--output-format", "json"}},
		{"work", []string{"a", "--", "b"}, "work", []string{"claude", "a", "--", "b"}},
	}
	for _, c := range cases {
		pl, err := Prepare(p, cfg, c.sel, c.args)
		if err != nil {
			t.Fatal(err)
		}
		if pl.Profile.Name != c.profile || !reflect.DeepEqual(pl.Argv, c.argv) {
			t.Errorf("%q %q: got %s %q", c.sel, c.args, pl.Profile.Name, pl.Argv)
		}
		if pl.Dir != p.ProfileDir(c.profile) {
			t.Errorf("dir %s", pl.Dir)
		}
	}
}

func TestPrepareUsesTheFoldersAccount(t *testing.T) {
	p, cfg := setup(t)
	here := t.TempDir()
	t.Chdir(here)
	cwd, _ := os.Getwd()
	cfg.SetFolder(filepath.Dir(cwd), "personal")
	if pl, err := Prepare(p, cfg, "", nil); err != nil || pl.Profile.Name != "personal" {
		t.Fatalf("no name inside the folder: %v %v", pl, err)
	}
	if pl, err := Prepare(p, cfg, "work", nil); err != nil || pl.Profile.Name != "work" {
		t.Fatalf("a name wins over the folder: %v %v", pl, err)
	}
	cfg.SetFolder(cwd, "gone")
	if _, err := Prepare(p, cfg, "", nil); err == nil || !strings.Contains(err.Error(), "ccam folder") {
		t.Fatalf("missing account: %v", err)
	}
}

func TestPrepareUnknownProfile(t *testing.T) {
	p, cfg := setup(t)
	if _, err := Prepare(p, cfg, "nope", nil); err == nil {
		t.Fatal("want error")
	}
}

func TestEnvReplacesConfigDir(t *testing.T) {
	// An inherited CCAM_RELAY would make a nested session's switch name
	// the outer relay and stop its own claude. An inherited Remote Control
	// session id (a shell of a Claude Code that is connected) would make the
	// hook take this session for a connected one: Claude Code sets it itself.
	env := buildEnv([]string{"A=1", "CLAUDE_CONFIG_DIR=/old", "CCAM_PROFILE=x", "CCAM_RELAY=2:1", "CLAUDE_SECURESTORAGE_CONFIG_DIR=/y",
		"CLAUDE_CODE_BRIDGE_SESSION_ID=cse_parent"}, "/new", "work", false)
	want := []string{"A=1", "CLAUDE_CONFIG_DIR=/new", "CCAM_PROFILE=work"}
	if !reflect.DeepEqual(env, want) {
		t.Fatalf("%q", env)
	}
	env = buildEnv(nil, "/h/.claude", "main", true)
	want = []string{"CLAUDE_CONFIG_DIR=/h/.claude", "CCAM_PROFILE=main", "CLAUDE_SECURESTORAGE_CONFIG_DIR="}
	if !reflect.DeepEqual(env, want) {
		t.Fatalf("default store: %q", env)
	}
}

// An imported ~/.claude uses the secure-storage entry plain `claude` uses,
// on Windows too: Claude Code's Credential Manager backend names its entry
// after the folder like the macOS Keychain does.
func TestDefaultStoreWhereLoginsHaveOne(t *testing.T) {
	p := paths.Paths{Home: "/h"}
	env := ProfileEnv(p, nil, filepath.Join("/h", ".claude"), "main")
	if got := slices.Contains(env, "CLAUDE_SECURESTORAGE_CONFIG_DIR="); got != (runtime.GOOS == "darwin" || runtime.GOOS == "windows") {
		t.Fatalf("%s: %q", runtime.GOOS, env)
	}
	if slices.Contains(ProfileEnv(p, nil, filepath.Join("/h", ".claude-work"), "work"), "CLAUDE_SECURESTORAGE_CONFIG_DIR=") {
		t.Fatal("only for ~/.claude")
	}
}

func TestInsertArgsGoesBeforeUserArgs(t *testing.T) {
	pl := &Plan{Argv: []string{"claude", "--mcp-config", "mine.json"}}
	pl.InsertArgs("--mcp-config", "group.json")
	want := []string{"claude", "--mcp-config", "group.json", "--mcp-config", "mine.json"}
	if !reflect.DeepEqual(pl.Argv, want) {
		t.Fatalf("%q", pl.Argv)
	}
}

func TestDryRunQuotes(t *testing.T) {
	environ := []string{"PATH=/bin", "CCAM_RELAY=1"}
	pl := &Plan{Dir: "/d", Bin: "/bin/claude", Argv: []string{"claude", "-p", "it's ünïcode", "--x"},
		Env: buildEnv(environ, "/d", "w", false)}
	var buf bytes.Buffer
	pl.dryRun(&buf, environ)
	want := `env -u CCAM_RELAY CLAUDE_CONFIG_DIR=/d CCAM_PROFILE=w /bin/claude -p 'it'\''s ünïcode' --x` + "\n"
	if runtime.GOOS == "windows" {
		want = `& { $saved = @{}; foreach ($k in 'CCAM_RELAY','CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; ` +
			`Remove-Item Env:CCAM_RELAY -ErrorAction SilentlyContinue; $env:CLAUDE_CONFIG_DIR='/d'; $env:CCAM_PROFILE='w'; ` +
			`try { $p = Start-Process -FilePath '/bin/claude' -ArgumentList '-p "it''s ünïcode" --x' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } ` +
			`finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }` + "\n"
	}
	if buf.String() != want {
		t.Fatalf("got %s", buf.String())
	}
}

// A program that starts in a folder of its own (ccusage's bun x and npx)
// has a dry run that starts it there too: in a subshell for sh, which
// leaves the shell where it is, and with Start-Process's -WorkingDirectory
// for PowerShell.
func TestDryRunStartsInTheWorkDir(t *testing.T) {
	environ := []string{"PATH=/bin"}
	pl := &Plan{Bin: "/bin/npx", Argv: []string{"npx", "ccusage@latest", "daily"}, Env: append(environ, "CLAUDE_CONFIG_DIR=/d"),
		WorkDir: "/state/it's"}
	var buf bytes.Buffer
	if err := pl.dryRun(&buf, environ); err != nil {
		t.Fatal(err)
	}
	if want := `(cd '/state/it'\''s' && CLAUDE_CONFIG_DIR=/d /bin/npx ccusage@latest daily)` + "\n"; runtime.GOOS != "windows" && buf.String() != want {
		t.Fatalf("got  %swant %s", buf.String(), want)
	}
	if want := ` -ArgumentList 'ccusage@latest daily' -WorkingDirectory '/state/it''s' -NoNewWindow `; runtime.GOOS == "windows" && !strings.Contains(buf.String(), want) {
		t.Fatalf("got %s\nwant it to hold %s", buf.String(), want)
	}
}

// PowerShell ends a '...' string at any of its single quotes, the
// typographic ones too: each is doubled.
func TestPSString(t *testing.T) {
	for in, want := range map[string]string{
		`C:\Program Files\x`:       `'C:\Program Files\x'`,
		"it's":                     `'it''s'`,
		"":                         `''`,
		"$env:X":                   `'$env:X'`,
		"don\u2019t \u2018x\u2019": "'don\u2019\u2019t \u2018\u2018x\u2019\u2019'",
		"\u201a\u201b":             "'\u201a\u201a\u201b\u201b'",
	} {
		if got := PSString(in); got != want {
			t.Errorf("%q: %s", in, got)
		}
	}
}

func TestResolveClaudeCachesLookup(t *testing.T) {
	p, _ := setup(t)
	early, binDir := t.TempDir(), t.TempDir()
	fake := filepath.Join(binDir, "claude-x"+exeSuffix)
	os.WriteFile(fake, []byte("#!/bin/sh\n"), 0o755)
	t.Setenv("PATH", early+string(os.PathListSeparator)+binDir)
	got, err := ResolveClaude(p, "claude-x")
	if err != nil || got != fake {
		t.Fatalf("%q %v", got, err)
	}
	if data, _ := os.ReadFile(filepath.Join(p.State, "claude-path")); !strings.HasPrefix(string(data), fake+"\t") {
		t.Fatalf("cache %q", data)
	}
	// With PATH unchanged the cache answers without scanning PATH, so a
	// file added to an earlier folder is not seen yet.
	newer := filepath.Join(early, "claude-x"+exeSuffix)
	os.WriteFile(newer, []byte("#!/bin/sh\n"), 0o755)
	if got, err := ResolveClaude(p, "claude-x"); err != nil || got != fake {
		t.Fatalf("cache miss: %q %v", got, err)
	}
	// Another claude_bin is not answered from the cache.
	if _, err := ResolveClaude(p, "claude-beta"); err == nil {
		t.Fatal("cache answered for a different claude_bin")
	}
	// A changed file invalidates the cache.
	os.Chtimes(fake, time.Now().Add(time.Hour), time.Now().Add(time.Hour))
	if got, err := ResolveClaude(p, "claude-x"); err != nil || got != newer {
		t.Fatalf("stale cache should trigger a lookup: %q %v", got, err)
	}
}

// A changed PATH (nvm use, a new install placed earlier) is not answered
// from the cache, even though the old file is still there unchanged.
func TestResolveClaudeFollowsPATH(t *testing.T) {
	p, _ := setup(t)
	a, b := t.TempDir(), t.TempDir()
	for _, d := range []string{a, b} {
		os.WriteFile(filepath.Join(d, "claude-x"+exeSuffix), []byte("#!/bin/sh\n"), 0o755)
	}
	t.Setenv("PATH", a+string(os.PathListSeparator)+b)
	if got, _ := ResolveClaude(p, "claude-x"); got != filepath.Join(a, "claude-x"+exeSuffix) {
		t.Fatalf("first lookup %q", got)
	}
	t.Setenv("PATH", b+string(os.PathListSeparator)+a)
	if got, _ := ResolveClaude(p, "claude-x"); got != filepath.Join(b, "claude-x"+exeSuffix) {
		t.Fatalf("PATH now starts with %s, got the cached %q", b, got)
	}
	t.Setenv("PATH", b)
	if got, _ := ResolveClaude(p, "claude-x"); got != filepath.Join(b, "claude-x"+exeSuffix) {
		t.Fatalf("%q", got)
	}
}

// A claude_bin that cannot run is named before anything of the launch
// happens; exec itself would only say "no such file or directory".
func TestPrepareRejectsAClaudeBinThatCannotRun(t *testing.T) {
	p, cfg := setup(t)
	noexec := filepath.Join(t.TempDir(), "claude")
	os.WriteFile(noexec, []byte("#!/bin/sh\n"), 0o644)
	bad := []string{filepath.Join(t.TempDir(), "missing", "claude"), t.TempDir(), "~/no/claude"}
	if runtime.GOOS != "windows" {
		bad = append(bad, noexec)
	}
	for _, bin := range bad {
		cfg.ClaudeBin = bin
		if _, err := Prepare(p, cfg, "work", nil); err == nil || !strings.Contains(err.Error(), "claude_bin") {
			t.Errorf("%s: got %v", bin, err)
		}
	}
	if !LastUsed(p, "work").IsZero() {
		t.Error("a launch that could not start marked the account as used")
	}
}

func TestClaudeBinExpandsTilde(t *testing.T) {
	p, cfg := setup(t)
	bin := filepath.Join(p.Home, "bin", "claude")
	os.MkdirAll(filepath.Dir(bin), 0o700)
	os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755)
	cfg.ClaudeBin = "~/bin/claude"
	pl, err := Prepare(p, cfg, "work", nil)
	if err != nil {
		t.Fatal(err)
	}
	if pl.Bin != bin {
		t.Fatalf("got %s", pl.Bin)
	}
}

func TestLastUsed(t *testing.T) {
	p, _ := setup(t)
	if !LastUsed(p, "work").IsZero() {
		t.Fatal("expected zero")
	}
	MarkUsed(p, "work")
	if time.Since(LastUsed(p, "work")) > time.Minute {
		t.Fatal("not recorded")
	}
}

func BenchmarkLaunchPath(b *testing.B) {
	p, cfg := setup(b)
	var names []string
	for _, prof := range cfg.Profiles {
		names = append(names, prof.Name)
	}
	if err := cfg.AddGroup("default", config.DefaultItems, names); err != nil {
		b.Fatal(err)
	}
	if err := cfg.Save(p.ConfigFile()); err != nil {
		b.Fatal(err)
	}
	// First run creates the share links; the loop measures the steady state.
	if err := Launch(p, "work", []string{"--continue"}, true, io.Discard); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Launch(p, "work", []string{"--continue"}, true, io.Discard); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkResolveClaudeCached is the lookup of a claude_bin without a
// path (the default) when the cache answers, with a WSL-sized PATH.
func BenchmarkResolveClaudeCached(b *testing.B) {
	p, _ := setup(b)
	dir := b.TempDir()
	os.WriteFile(filepath.Join(dir, "claude-x"+exeSuffix), []byte("#!/bin/sh\n"), 0o755)
	b.Setenv("PATH", strings.Repeat("/mnt/c/Program Files/Some Tool/bin"+string(os.PathListSeparator), 30)+dir)
	if _, err := ResolveClaude(p, "claude-x"); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ResolveClaude(p, "claude-x"); err != nil {
			b.Fatal(err)
		}
	}
}

func TestRemoteControlCopiesMCPOnYes(t *testing.T) {
	root := t.TempDir()
	gdir, state, dir := filepath.Join(root, "group"), filepath.Join(root, "state"), filepath.Join(root, "profile")
	os.MkdirAll(gdir, 0o700)
	os.MkdirAll(dir, 0o700)
	os.WriteFile(filepath.Join(gdir, "mcp.json"), []byte(`{"mcpServers":{"fs":{"command":"fs-mcp"}}}`), 0o600)
	asked := 0
	old := confirm
	confirm = func(string) bool { asked++; return true }
	t.Cleanup(func() { confirm = old })
	pl := &Plan{Profile: &config.Profile{Name: "work"}, Dir: dir}
	idle := func() bool { return false }

	remoteControlMCP(pl, "default", gdir, state, false, idle)
	if local, _ := os.ReadFile(filepath.Join(dir, ".claude.json")); asked != 1 || !strings.Contains(string(local), "fs-mcp") {
		t.Fatalf("asked %d, .claude.json %s", asked, local)
	}
	remoteControlMCP(pl, "default", gdir, state, false, idle)
	if asked != 1 {
		t.Fatal("asked again although the copy is current")
	}
	later := time.Now().Add(time.Second)
	os.WriteFile(filepath.Join(gdir, "mcp.json"), []byte(`{"mcpServers":{"fs":{"command":"fs-mcp-2"}}}`), 0o600)
	os.Chtimes(filepath.Join(gdir, "mcp.json"), later, later)
	remoteControlMCP(pl, "default", gdir, state, false, func() bool { return true }) // running: no question, no write
	if asked != 1 {
		t.Fatal("asked while the account runs")
	}
	remoteControlMCP(pl, "default", gdir, state, false, idle)
	if local, _ := os.ReadFile(filepath.Join(dir, ".claude.json")); asked != 2 || !strings.Contains(string(local), "fs-mcp-2") {
		t.Fatalf("after the group changed: asked %d, .claude.json %s", asked, local)
	}
}

func TestFolderHereFollowsSymlinks(t *testing.T) {
	root := t.TempDir()
	real := filepath.Join(root, "real")
	os.MkdirAll(filepath.Join(real, "sub"), 0o700)
	if err := os.Symlink(real, filepath.Join(root, "link")); err != nil {
		t.Skip("no symlinks:", err)
	}
	resolved, _ := filepath.EvalSymlinks(real)
	cfg := &config.Config{Profiles: []config.Profile{{Name: "work", Slot: 1}}}
	cfg.SetFolder(resolved, "work")
	t.Chdir(filepath.Join(root, "link", "sub"))
	if f, ok := FolderHere(cfg); !ok || f.Account != "work" {
		t.Fatalf("through the link: %v %v", f, ok)
	}
}

// The relay's hop is ccam's own: beside another session sharing the
// history it warns and goes on, asking nothing and leaving the terminal's
// input alone; a launch by hand asks there.
func TestHistoryCheckAsksNothingInAHop(t *testing.T) {
	p, _ := setup(t)
	t.Chdir(t.TempDir())
	cwd, _ := os.Getwd()
	cwd, _ = filepath.EvalSymlinks(cwd)
	if err := session.Register(HistoryDir(p, "default", cwd)); err != nil { // the other session (this process)
		t.Fatal(err)
	}
	oldIn, oldAttended := os.Stdin, attended
	attended = func() bool { return true } // as on a terminal
	t.Cleanup(func() { os.Stdin, attended, Unattended = oldIn, oldAttended, false })
	pl := &Plan{Profile: &config.Profile{Name: "work", Share: config.Share{History: "default"}}, UserArgs: []string{"--resume", "s1"}}
	for _, hop := range []bool{true, false} {
		r, w, _ := os.Pipe()
		w.Write([]byte("y\n"))
		os.Stdin, Unattended = r, hop
		if err := historyCheck(p, pl); err != nil {
			t.Fatal(err)
		}
		w.Close()
		if rest, _ := io.ReadAll(r); hop != (string(rest) == "y\n") {
			t.Fatalf("hop %v: %q left of the answer", hop, rest)
		}
	}
}

// listedRunning writes what Claude Code lists for a session of its own in
// the config folder dir: this test process stands in for a claude that ccam
// did not start.
func listedRunning(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "sessions", strconv.Itoa(os.Getpid())+".json")
	os.MkdirAll(filepath.Dir(path), 0o700)
	if err := os.WriteFile(path, fmt.Appendf(nil, `{"pid":%d,"kind":"interactive","startedAt":%d}`, os.Getpid(), time.Now().UnixMilli()), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// With mcp_strategy = "sync" a launch rewrites .claude.json, which a running
// Claude Code writes too: an account that Claude Code runs on, started
// without ccam, is left alone like one ccam started.
func TestSyncLeavesAnAccountClaudeCodeRunsOn(t *testing.T) {
	p, cfg := setup(t)
	if _, ok := session.StartTime(os.Getpid()); !ok {
		t.Skip("no process start times on this platform")
	}
	cfg.MCPStrategy = config.MCPSync
	cfg.Profiles[1].Share = config.Share{MCP: "default"}
	os.MkdirAll(p.GroupDir("default"), 0o700)
	os.WriteFile(filepath.Join(p.GroupDir("default"), "mcp.json"), []byte(`{"mcpServers":{"fs":{"command":"fs-mcp"}}}`), 0o600)
	dir := p.ProfileDir("work")
	listed := listedRunning(t, dir)
	if _, err := Ready(p, cfg, "work", nil, false); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(filepath.Join(dir, ".claude.json")); strings.Contains(string(data), "fs-mcp") {
		t.Fatal(".claude.json rewritten under a running Claude Code")
	}
	os.Remove(listed)
	session.Unregister(SessionDir(p, "work")) // the marker of the first launch (this process)
	if _, err := Ready(p, cfg, "work", nil, false); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(filepath.Join(dir, ".claude.json")); !strings.Contains(string(data), "fs-mcp") {
		t.Fatalf("not synced once nothing runs: %s", data)
	}
}

// A marker that cannot be written (the state folder is broken) does not
// stop the launch.
func TestLaunchGoesOnWithoutAMarker(t *testing.T) {
	p, cfg := setup(t)
	if _, ok := session.StartTime(os.Getpid()); !ok {
		t.Skip("no session markers on this platform")
	}
	if err := os.WriteFile(p.Runtime, []byte("not a folder"), 0o600); err != nil {
		t.Fatal(err)
	}
	pl, err := Ready(p, cfg, "work", nil, false)
	if err != nil {
		t.Fatalf("launch stopped: %v", err)
	}
	if len(pl.sessions) != 0 {
		t.Fatalf("registered in %v", pl.sessions)
	}
}

// BenchmarkSyncLaunchUnchanged is a launch's MCP step under the sync
// strategy when nothing changed since the last one: a 1 MB .claude.json
// and a group with one server.
func BenchmarkSyncLaunchUnchanged(b *testing.B) {
	p, cfg := setup(b)
	cfg.MCPStrategy = config.MCPSync
	if err := cfg.AddGroup("g", []string{"mcp"}, []string{"work"}); err != nil {
		b.Fatal(err)
	}
	os.MkdirAll(p.GroupDir("g"), 0o700)
	os.WriteFile(filepath.Join(p.GroupDir("g"), "mcp.json"), []byte(`{"mcpServers":{"g":{"command":"g"}}}`), 0o600)
	big := `{"projects":{"/p":{"history":["` + strings.Repeat("x", 1<<20) + `"]}},"mcpServers":{}}`
	os.WriteFile(filepath.Join(p.ProfileDir("work"), ".claude.json"), []byte(big), 0o600)
	pl, err := Prepare(p, cfg, "work", nil)
	if err != nil {
		b.Fatal(err)
	}
	if err := applyMCP(p, cfg, pl, false); err != nil { // the first sync
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := applyMCP(p, cfg, pl, false); err != nil {
			b.Fatal(err)
		}
	}
}
