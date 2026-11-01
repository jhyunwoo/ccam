// Package launcher turns "profile + claude args" into an exec of claude.
//
// This is the hot path behind `ccam <name>` and shortcuts: no cobra, no TUI,
// no HTTP, one config read, lstat-only checks.
package launcher

import (
	"context"
	"errors"
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
	"sync/atomic"
	"syscall"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/share"
)

type Plan struct {
	Profile *config.Profile
	Dir     string   // CLAUDE_CONFIG_DIR
	Bin     string   // absolute path of claude
	Argv    []string // argv including argv[0]
	Env     []string
	// UserArgs are the claude args as the user gave them (one leading "--"
	// dropped), without the ones ccam adds such as --mcp-config.
	UserArgs []string
	// WorkDir is the folder a program other than claude starts in when it
	// is not the working one, as ccusage's bun x and npx do. Only the dry
	// run reads it: their caller changes to the folder itself.
	WorkDir string

	sessions []string // session folders this process is registered in
	keyVar   string   // dry run of an API account: the key variable, not set,
	keyFrom  string   // and where a launch reads its value
	unmarked bool     // a marker could not be written (said once)
	run      runState // what Run saw
}

// runState is what Run saw of claude and of the signals ccam itself got, by
// which the relay judges whether to carry the conversation on. The signal
// goroutine may still run after Run returns, hence the atomics.
type runState struct {
	child     int          // claude's process
	signaled  atomic.Bool  // ccam got SIGHUP, SIGTERM or SIGQUIT (passed on to claude)
	interrupt atomic.Int64 // when ccam last got Ctrl+C (SIGINT, not passed on), UnixNano; 0 for never
}

// Prepare resolves the profile (sel == "" means the one set for the
// current folder, else the default one) and builds the exec plan. User args
// are passed through untouched, except that one leading "--" separator is
// dropped.
func Prepare(p paths.Paths, cfg *config.Config, sel string, args []string) (*Plan, error) {
	var prof *config.Profile
	var err error
	if sel == "" {
		prof, err = HereProfile(cfg)
	} else {
		prof, err = cfg.MustFind(sel)
	}
	if err != nil {
		return nil, err
	}
	args = ClaudeArgs(args)
	bin, err := ResolveClaude(p, cfg.ClaudeBin)
	if err != nil {
		return nil, err
	}
	dir := profile.Dir(p, prof)
	if _, err := os.Stat(dir); err != nil {
		return nil, fmt.Errorf("the folder of account %s is missing (%s). Check with `ccam doctor`", prof.Name, dir)
	}
	argv := make([]string, 0, len(args)+1)
	argv = append(argv, "claude")
	argv = append(argv, args...)
	return &Plan{Profile: prof, Dir: dir, Bin: bin, Argv: argv, Env: ProfileEnv(p, os.Environ(), dir, prof.Name), UserArgs: args}, nil
}

// ClaudeArgs is what claude gets of the args after the account name: all
// of them, less one leading "--" separator.
func ClaudeArgs(args []string) []string {
	if len(args) > 0 && args[0] == "--" {
		return args[1:]
	}
	return args
}

// HereProfile is the account a launch without a name uses: the one set for
// the current folder with `ccam folder`, else the default account.
func HereProfile(cfg *config.Config) (*config.Profile, error) {
	if len(cfg.Folders) > 0 {
		if f, ok := FolderHere(cfg); ok {
			if prof, ok := cfg.Find(f.Account); ok {
				return prof, nil
			}
			return nil, fmt.Errorf("%s is set to account %q, which does not exist. Change it with `ccam folder`", f.Path, f.Account)
		}
	}
	return cfg.DefaultProfile()
}

// FolderHere finds the `ccam folder` entry for the working directory. Entries
// are stored with symlinks resolved; the working directory is resolved only
// when it does not match as it is, which keeps the common case to one getcwd.
func FolderHere(cfg *config.Config) (*config.Folder, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, false
	}
	if f, ok := cfg.FolderFor(cwd); ok {
		return f, true
	}
	if real, err := filepath.EvalSymlinks(cwd); err == nil && real != cwd {
		return cfg.FolderFor(real)
	}
	return nil, false
}

// InsertArgs puts ccam-generated args (e.g. --mcp-config) before the user's,
// so user flags of the same kind still apply on top.
func (pl *Plan) InsertArgs(extra ...string) {
	argv := make([]string, 0, len(pl.Argv)+len(extra))
	argv = append(argv, pl.Argv[0])
	argv = append(argv, extra...)
	pl.Argv = append(argv, pl.Argv[1:]...)
}

// ProfileEnv is the environment claude runs with for a profile.
func ProfileEnv(p paths.Paths, environ []string, dir, name string) []string {
	return buildEnv(environ, dir, name, secureStore && creds.DefaultStore(dir, p.Home))
}

// secureStore says the platform has a login store Claude Code names after
// the config folder: the macOS Keychain, Windows Credential Manager.
const secureStore = runtime.GOOS == "darwin" || runtime.GOOS == "windows"

// buildEnv sets CLAUDE_CONFIG_DIR. defaultStore keeps an imported ~/.claude
// on the unsuffixed secure-storage entry that plain `claude` uses (macOS
// Keychain, Windows Credential Manager); with the file store the variable
// names the same folder. A Remote Control session id inherited from a Claude
// Code that is connected is dropped: Claude Code sets it when its own
// connects, and the hook goes by it.
func buildEnv(environ []string, dir, name string, defaultStore bool) []string {
	env := make([]string, 0, len(environ)+3)
	for _, kv := range environ {
		if strings.HasPrefix(kv, "CLAUDE_CONFIG_DIR=") || strings.HasPrefix(kv, "CCAM_PROFILE=") ||
			strings.HasPrefix(kv, "CCAM_RELAY=") || strings.HasPrefix(kv, "CLAUDE_SECURESTORAGE_CONFIG_DIR=") ||
			strings.HasPrefix(kv, "CLAUDE_CODE_BRIDGE_SESSION_ID=") {
			continue
		}
		env = append(env, kv)
	}
	env = append(env, "CLAUDE_CONFIG_DIR="+dir, "CCAM_PROFILE="+name)
	if defaultStore {
		env = append(env, "CLAUDE_SECURESTORAGE_CONFIG_DIR=")
	}
	return env
}

// DryRun prints what would be executed, as a line for sh, or for
// PowerShell on Windows: the variables a launch sets or clears, then the
// command, started in WorkDir if the plan has one. The key of an API
// account is shown by its source, never read.
// It fails, as the launch would, where claude could not be started so.
func (pl *Plan) DryRun(w io.Writer) error { return pl.dryRun(w, os.Environ()) }

func (pl *Plan) dryRun(w io.Writer, environ []string) error {
	set, unset := pl.envChanges(environ)
	if runtime.GOOS == "windows" {
		// PowerShell: the variables are put back afterwards, as sh's
		// VAR=x cmd would, and claude gets its command line through
		// Start-Process unchanged: `& claude args` would be quoted
		// differently by Windows PowerShell 5.1 and PowerShell 7.
		file, line, err := commandLine(pl.Bin, pl.Argv[1:])
		if err != nil {
			return err
		}
		var names, parts, empty []string
		for _, k := range unset {
			names = append(names, "'"+k+"'")
			parts = append(parts, "Remove-Item Env:"+k+" -ErrorAction SilentlyContinue;")
		}
		for _, kv := range set {
			k, v, _ := strings.Cut(kv, "=")
			names = append(names, "'"+k+"'")
			// The value of an assignment must be a string, so it is always quoted.
			parts = append(parts, "$env:"+k+"="+PSString(v)+";")
			if v == "" {
				empty = append(empty, k)
			}
		}
		if pl.keyVar != "" {
			v := PSString("<the first line of: " + strings.TrimPrefix(pl.keyFrom, "command:") + ">")
			if name, ok := strings.CutPrefix(pl.keyFrom, "env:"); ok {
				v = "$env:" + name
			}
			names = append(names, "'"+pl.keyVar+"'")
			parts = append(parts, "$env:"+pl.keyVar+"="+v+";")
		}
		start := "Start-Process -FilePath " + PSString(file)
		if line != "" {
			start += " -ArgumentList " + PSString(line)
		}
		if pl.WorkDir != "" {
			start += " -WorkingDirectory " + PSString(pl.WorkDir)
		}
		fmt.Fprintln(w, "& { $saved = @{}; foreach ($k in "+strings.Join(names, ",")+") { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; "+
			strings.Join(parts, " ")+" try { $p = "+start+" -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } "+
			"finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }")
		if len(empty) > 0 {
			fmt.Fprintf(w, "# The launch sets %s to an empty value; Windows PowerShell 5.1 (and PowerShell before 7.5) removes a variable set to '' instead, so there this line differs from the launch.\n", strings.Join(empty, ", "))
		}
		return nil
	}
	var parts []string
	if len(unset) > 0 {
		parts = append(parts, "env")
		for _, k := range unset {
			parts = append(parts, "-u", k)
		}
	}
	for _, kv := range set {
		k, v, _ := strings.Cut(kv, "=")
		parts = append(parts, k+"="+ShellQuote(v))
	}
	if pl.keyVar != "" {
		v := `"$(` + strings.TrimPrefix(pl.keyFrom, "command:") + ` | head -n 1)"`
		if name, ok := strings.CutPrefix(pl.keyFrom, "env:"); ok {
			v = `"$` + name + `"`
		}
		parts = append(parts, pl.keyVar+"="+v)
	}
	parts = append(parts, ShellQuote(pl.Bin))
	for _, a := range pl.Argv[1:] {
		parts = append(parts, ShellQuote(a))
	}
	line := strings.Join(parts, " ")
	if pl.WorkDir != "" {
		// A subshell, so that the shell the line is pasted into stays where
		// it is.
		line = "(cd " + ShellQuote(pl.WorkDir) + " && " + line + ")"
	}
	fmt.Fprintln(w, line)
	return nil
}

// envChanges compares the plan's environment with the one ccam inherited:
// set holds what a launch adds or changes (CLAUDE_CONFIG_DIR and
// CCAM_PROFILE always, so the line stands on its own), unset the names it
// drops. Inherited values a launch keeps are not printed, so a key in the
// shell is never echoed.
func (pl *Plan) envChanges(environ []string) (set, unset []string) {
	before := make(map[string]string, len(environ))
	for _, kv := range environ {
		k, v, _ := strings.Cut(kv, "=")
		before[k] = v
	}
	kept := make(map[string]bool, len(pl.Env))
	for _, kv := range pl.Env {
		k, v, _ := strings.Cut(kv, "=")
		kept[k] = true
		if old, ok := before[k]; !ok || old != v || k == "CLAUDE_CONFIG_DIR" || k == "CCAM_PROFILE" {
			set = append(set, kv)
		}
	}
	for k := range before {
		if !kept[k] && k != pl.keyVar {
			unset = append(unset, k)
		}
	}
	slices.Sort(unset)
	return set, unset
}

// PSString is s as a PowerShell string literal: in single quotes, each
// quote inside doubled. PowerShell takes the typographic single quotes
// (U+2018 to U+201B) for ' as well, so they are doubled too.
func PSString(s string) string { return "'" + psQuotes.Replace(s) + "'" }

var psQuotes = strings.NewReplacer("'", "''", "\u2018", "\u2018\u2018", "\u2019", "\u2019\u2019", "\u201a", "\u201a\u201a", "\u201b", "\u201b\u201b")

// BeforeExec runs right before exec, e.g. to flush a CPU profile.
var BeforeExec func()

// ExitCode is returned when claude ran and ended with a non-zero code that
// ccam must exit with, without a message of its own.
type ExitCode int

func (e ExitCode) Error() string { return fmt.Sprintf("claude exited with code %d", int(e)) }

// Run starts claude as a child process with the given standard streams,
// waits for it and returns its exit code (128+n when signal n ended it).
// Ctrl+C is not passed on: the terminal sends it to claude as well. Run
// records claude's process and the signals ccam got (pl.run).
func (pl *Plan) Run(p paths.Paths, stdin io.Reader, stdout, stderr io.Writer) (int, error) {
	setTitle(pl.Profile.Name)
	cmd := &exec.Cmd{Path: pl.Bin, Args: pl.Argv}
	if isBatch(pl.Bin) {
		c, err := batchCommand(context.Background(), pl.Bin, pl.Argv[1:])
		if err != nil {
			return 0, err
		}
		cmd = c
	}
	cmd.Env, cmd.Stdin, cmd.Stdout, cmd.Stderr = pl.Env, stdin, stdout, stderr
	sig := make(chan os.Signal, 4)
	signal.Notify(sig, passOn...)
	defer signal.Stop(sig)
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	MarkUsed(p, pl.Profile.Name)
	pl.run.child = cmd.Process.Pid
	done := make(chan struct{})
	go func() {
		for {
			select {
			case s := <-sig:
				if s == os.Interrupt {
					pl.run.interrupt.Store(time.Now().UnixNano())
					continue
				}
				pl.run.signaled.Store(true)
				_ = cmd.Process.Signal(s)
			case <-done:
				return
			}
		}
	}()
	err := cmd.Wait()
	close(done)
	var ee *exec.ExitError
	switch {
	case err == nil:
		return 0, nil
	case errors.As(err, &ee):
		if ws, ok := ee.Sys().(syscall.WaitStatus); ok && ws.Signaled() {
			return 128 + int(ws.Signal()), nil
		}
		return ee.ExitCode(), nil
	}
	return 0, err
}

// Signaled reports whether ccam got SIGHUP, SIGTERM or SIGQUIT while Run
// ran claude (and passed it on).
func (pl *Plan) Signaled() bool { return pl.run.signaled.Load() }

func setTitle(name string) {
	if runtime.GOOS == "windows" {
		return // the console shows escape sequences as text unless VT output is on
	}
	if fi, err := os.Stderr.Stat(); err == nil && fi.Mode()&os.ModeCharDevice != 0 {
		fmt.Fprintf(os.Stderr, "\x1b]0;claude · %s\a", name)
	}
}

func lastUsedPath(p paths.Paths, name string) string {
	return filepath.Join(p.State, "last-used", name)
}

// MCPStateDir keeps an account's MCP sync stamp and .claude.json backups.
func MCPStateDir(p paths.Paths, name string) string {
	return filepath.Join(p.State, "mcp", name)
}

// StatePaths lists what ccam keeps under an account's name in the state
// folder, so rename can carry it along and remove can clear it.
func StatePaths(p paths.Paths, name string) []string {
	return []string{lastUsedPath(p, name), MCPStateDir(p, name), filepath.Join(p.State, "picks", name)}
}

// MarkUsed records the launch time; losing it is harmless, so no fsync.
func MarkUsed(p paths.Paths, name string) {
	_ = fsutil.WriteFast(lastUsedPath(p, name), []byte(strconv.FormatInt(time.Now().Unix(), 10)))
}

// LastUsed returns the last launch time, or zero.
func LastUsed(p paths.Paths, name string) time.Time {
	data, err := os.ReadFile(lastUsedPath(p, name))
	if err != nil {
		return time.Time{}
	}
	sec, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(sec, 0)
}

// ResolveClaude finds the claude executable, which must not be ccam itself
// (that would add ccam's flags on every pass). A claude_bin with a path
// ("~/" is the home folder) must be an executable file: checked here, before
// a launch marks the account as used or takes a handoff note, since a failed
// exec says only "no such file or directory". Lookups through PATH are
// cached in the state dir with the file's mtime, the claude_bin they
// resolved and a hash of PATH; a changed or missing file, another
// claude_bin or another PATH (nvm use, a new install placed earlier)
// triggers a fresh lookup.
func ResolveClaude(p paths.Paths, bin string) (string, error) {
	if bin == "" {
		bin = "claude"
	}
	if strings.ContainsRune(bin, '/') || strings.ContainsRune(bin, filepath.Separator) {
		if rest, ok := strings.CutPrefix(bin, "~/"); ok {
			bin = filepath.Join(p.Home, rest)
		}
		fi, err := os.Stat(bin)
		if err != nil || !fi.Mode().IsRegular() || runtime.GOOS != "windows" && fi.Mode().Perm()&0o111 == 0 {
			return "", fmt.Errorf("claude_bin %s is not an executable file. Check claude_bin in config.toml", bin)
		}
		if isSelfInfo(fi) {
			return "", fmt.Errorf("%s points to ccam itself. Check claude_bin in config.toml", bin)
		}
		return bin, nil
	}
	cache := filepath.Join(p.State, "claude-path")
	key := bin + "\t" + pathHash(os.Getenv("PATH"))
	if data, err := os.ReadFile(cache); err == nil {
		path, rest, _ := strings.Cut(string(data), "\t")
		mtime, cachedKey, _ := strings.Cut(rest, "\t")
		if cachedKey == key {
			if fi, err := os.Stat(path); err == nil && strconv.FormatInt(fi.ModTime().UnixNano(), 10) == mtime {
				return path, nil
			}
		}
	}
	path, err := exec.LookPath(bin)
	if err != nil {
		return "", fmt.Errorf("cannot find the claude executable (%s). Is Claude Code installed?", bin)
	}
	path, _ = filepath.Abs(path)
	fi, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if isSelfInfo(fi) {
		return "", fmt.Errorf("%s points to ccam itself. Check claude_bin in config.toml", path)
	}
	_ = fsutil.WriteFast(cache, []byte(path+"\t"+strconv.FormatInt(fi.ModTime().UnixNano(), 10)+"\t"+key))
	return path, nil
}

// pathHash is FNV-1a of s, without allocating: PATH can be kilobytes long
// on WSL, which appends the Windows PATH.
func pathHash(s string) string {
	h := uint64(14695981039346656037)
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= 1099511628211
	}
	return strconv.FormatUint(h, 16)
}

// ClaudeCommand is exec.CommandContext for the claude at bin. On Windows a
// batch file (an npm install's claude.cmd) gets a command line built for
// cmd.exe, which would otherwise run what it finds and read the arguments
// as commands (batchCommand).
func ClaudeCommand(ctx context.Context, bin string, args ...string) (*exec.Cmd, error) {
	if isBatch(bin) {
		return batchCommand(ctx, bin, args)
	}
	return exec.CommandContext(ctx, bin, args...), nil
}

// isSelfInfo reports whether fi is the running ccam binary, through any
// symlink or hard link. Comparing inodes costs two stats, which matters on
// the launch path; resolving both paths component by component does not.
func isSelfInfo(fi os.FileInfo) bool {
	self, err := os.Executable()
	if err != nil {
		return false
	}
	a, err := os.Stat(self)
	return err == nil && os.SameFile(a, fi)
}

// ShellQuote quotes s for sh when needed.
func ShellQuote(s string) string {
	if s == "" {
		return "''"
	}
	if strings.IndexFunc(s, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("-_./=:@,+%", r))
	}) < 0 {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// checkLinks verifies share links with lstat only. Missing links are created
// on the spot; anything that could lose data is left for `ccam doctor --fix`.
func checkLinks(p paths.Paths, pl *Plan) error {
	problems := share.Check(pl.Dir, pl.Profile.Share, p.GroupsRoot())
	if len(problems) == 0 {
		return nil
	}
	rest, err := share.QuickFix(problems)
	if err != nil {
		return err
	}
	// An account's own copy of a path added to an item later waits for
	// doctor --fix; it never stops a launch.
	rest = slices.DeleteFunc(rest, func(pr share.Problem) bool { return pr.Waits(p.GroupsRoot()) })
	if len(rest) == 0 {
		return nil
	}
	msg := make([]string, 0, len(rest))
	for _, pr := range rest {
		msg = append(msg, "  "+pr.String())
	}
	return fmt.Errorf("stopped because of share link problems:\n%s\nRepair them with `ccam doctor --fix`, then run again", strings.Join(msg, "\n"))
}
