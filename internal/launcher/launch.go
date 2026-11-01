package launcher

import (
	"fmt"
	"hash/fnv"
	"io"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/charmbracelet/x/term"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/session"
	"github.com/jhyunwoo/ccam/internal/trace"
)

// Launch is the whole `ccam <name>` sequence: load config, prepare, check
// links, apply MCP sharing, take locks, then exec (or print when dry).
func Launch(p paths.Paths, sel string, args []string, dry bool, stdout io.Writer) error {
	trace.Step("start")
	cfg, err := config.Load(p.ConfigFile())
	if err != nil {
		return err
	}
	trace.Step("load config")
	return LaunchWith(p, cfg, sel, args, dry, stdout)
}

// LaunchWith is Launch with an already loaded config.
func LaunchWith(p paths.Paths, cfg *config.Config, sel string, args []string, dry bool, stdout io.Writer) error {
	return LaunchThen(p, cfg, sel, args, dry, stdout, nil)
}

// Unattended marks a launch nobody started by hand: the relay's hop, which
// carries ccam's own conversation on. It asks nothing but ConfirmMove's
// question, and what was typed while the conversation moved is dropped
// rather than sent to Claude Code.
var Unattended bool

// LaunchThen is LaunchWith that calls ready (unless nil or dry) once
// everything is set for claude to start.
func LaunchThen(p paths.Paths, cfg *config.Config, sel string, args []string, dry bool, stdout io.Writer, ready func()) error {
	// A setting written below a table is not read: say so, or it would
	// seem to do nothing.
	for _, u := range cfg.Unknown() {
		fmt.Fprintf(os.Stderr, "ccam: config.toml %s\n", u)
	}
	pl, err := Ready(p, cfg, sel, args, dry)
	if err != nil {
		return err
	}
	if dry {
		return pl.DryRun(stdout)
	}
	if ready != nil {
		ready()
	}
	if Unattended && discardInput(os.Stdin) {
		fmt.Fprintln(os.Stderr, "ccam: what you typed during the switch was discarded.")
	}
	if useRelay(cfg, runtime.GOOS) && relayable(pl.Argv[1:]) {
		return pl.relay(p)
	}
	return pl.Exec(p)
}

// useRelay says whether claude runs under the relay on the system goos:
// with relay = true, or auto_switch = true, which needs it, except on
// Windows, which ignores auto_switch (D6).
func useRelay(cfg *config.Config, goos string) bool {
	return cfg.Relay || cfg.AutoSwitch && goos != "windows"
}

// Ready does everything before claude starts: build the plan, check links,
// apply MCP sharing and, unless dry, register the session.
func Ready(p paths.Paths, cfg *config.Config, sel string, args []string, dry bool) (*Plan, error) {
	pl, err := Prepare(p, cfg, sel, args)
	if err != nil {
		return nil, err
	}
	trace.Step("prepare")
	if err := checkLinks(p, pl); err != nil {
		return nil, err
	}
	trace.Step("check links")
	if api := pl.Profile.API; api != nil {
		// A dry run prints where the key comes from, so it does not read one.
		if err := pl.useAPI(api, dry); err != nil {
			return nil, err
		}
	} else {
		warnAuthEnv(pl)
	}
	if err := prelaunch(p, cfg, pl, dry); err != nil {
		return nil, err
	}
	trace.Step("prelaunch")
	return pl, nil
}

// warnAuthEnv says when the environment holds a key or token that Claude
// Code would use instead of the account's own login.
func warnAuthEnv(pl *Plan) {
	for _, kv := range pl.Env {
		name, v, _ := strings.Cut(kv, "=")
		if v == "" {
			continue
		}
		switch name {
		case "CLAUDE_CODE_OAUTH_TOKEN":
			fmt.Fprintf(os.Stderr, "ccam: $%s is set, so Claude Code uses that token instead of %s's login. Unset it to use %s.\n", name, pl.Profile.Name, pl.Profile.Name)
		case "ANTHROPIC_API_KEY", "ANTHROPIC_AUTH_TOKEN":
			fmt.Fprintf(os.Stderr, "ccam: $%s is set, so Claude Code may use that key instead of %s's login. Unset it, or make the key an account of its own (ccam add <name> --api --key-from env:%s).\n", name, pl.Profile.Name, name)
		}
	}
}

// SessionDir holds markers of running sessions of a profile.
func SessionDir(p paths.Paths, name string) string {
	return filepath.Join(p.Runtime, "profile-"+name)
}

// Running reports whether Claude Code runs on the folder dir of account
// name: a session ccam started, or one started any other way that Claude
// Code lists in the folder (plain claude, an IDE, a background session, its
// daemon). Commands that move, rewrite or remove what Claude Code writes ask
// it right before they do. It reads the folder's records, so the launch
// path asks it only where it would rewrite .claude.json.
func Running(p paths.Paths, name, dir string) bool {
	return session.InUse(SessionDir(p, name), dir)
}

// prelaunch registers the session (the marker then names the exec'd claude
// process) and applies MCP sharing.
func prelaunch(p paths.Paths, cfg *config.Config, pl *Plan, dry bool) error {
	if err := applyMCP(p, cfg, pl, dry); err != nil {
		return err
	}
	if dry {
		return nil
	}
	if err := historyCheck(p, pl); err != nil {
		return err
	}
	pl.register(SessionDir(p, pl.Profile.Name))
	return nil
}

// register marks this process as a session in dir. After exec the marker
// names claude; when ccam waits for claude instead, Done removes it. A
// marker that cannot be written does not stop the launch: Claude Code's
// own record still shows the session to ccam's checks.
func (pl *Plan) register(dir string) {
	if err := session.Register(dir); err != nil {
		if !pl.unmarked {
			fmt.Fprintf(os.Stderr, "ccam: could not mark the session as running (%v); starting Claude Code anyway\n", err)
		}
		pl.unmarked = true
		return
	}
	pl.sessions = append(pl.sessions, dir)
}

// Done removes the session markers of a launch that ran claude with Run.
func (pl *Plan) Done() {
	for _, dir := range pl.sessions {
		session.Unregister(dir)
	}
	pl.sessions = nil
}

// HistoryDir guards one (history group, working dir) pair: two accounts
// sharing conversations in the same folder would write the same session files.
func HistoryDir(p paths.Paths, group, cwd string) string {
	h := fnv.New64a()
	h.Write([]byte(cwd))
	return filepath.Join(p.Runtime, fmt.Sprintf("history-%s-%016x", group, h.Sum64()))
}

func historyCheck(p paths.Paths, pl *Plan) error {
	group := pl.Profile.Share.Get("history")
	if group == config.Local || !opensSession(pl.UserArgs) {
		return nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	// Claude Code keeps a folder's conversations under its real path, so a
	// folder reached through a symbolic link is the same folder.
	if real, err := filepath.EvalSymlinks(cwd); err == nil {
		cwd = real
	}
	dir := HistoryDir(p, group, cwd)
	// A new -p conversation is one of its own: parallel workers go on.
	if session.Active(dir) && !newPrint(pl.UserArgs) {
		fmt.Fprintf(os.Stderr, "ccam: another Claude Code sharing conversation history (group %s) is running in this folder.\n      Continuing the same conversation from both at once can mix up the history.\n", group)
		switch {
		case Unattended:
			fmt.Fprintln(os.Stderr, "      ccam is continuing its own conversation, so it does not ask.")
		case !askTTY("Run anyway?"):
			return fmt.Errorf("cancelled")
		}
	}
	pl.register(dir)
	return nil
}

// opensSession reports whether claude args open a conversation here, which
// is what can mix up a shared history: not --version or --help, nor a
// subcommand that writes no conversation (Claude Code 2.1.308). The
// subcommand comes first: a "--" after it (`mcp add x -- cmd`) is its own.
func opensSession(args []string) bool {
	if len(args) > 0 {
		switch args[0] {
		case "agents", "auth", "auto-mode", "doctor", "gateway", "import", "install", "logs", "mcp", "plugin", "plugins",
			"purge", "rm", "setup-token", "stop", "kill", "ultrareview", "update", "upgrade":
			return false
		}
	}
	for _, a := range args {
		switch a {
		case "--":
			return true
		case "-v", "--version", "-h", "--help":
			return false
		}
	}
	return true
}

// newPrint reports whether claude args run -p in a conversation of their
// own, not one they continue (-c, --resume, --session-id …).
func newPrint(args []string) bool {
	printing := false
	for _, a := range args {
		name, _, _ := strings.Cut(a, "=")
		switch name {
		case "--":
			return printing
		case "-p", "--print":
			printing = true
		case "-c", "--continue", "-r", "--resume", "--session-id", "--from-pr", "--teleport":
			return false
		}
	}
	return printing
}

// remoteControlMCP offers to copy the group's servers into the account's
// .claude.json, the only place Remote Control sessions read them from. It
// asks only while that copy is missing, incomplete or older than the
// group's mcp.json.
func remoteControlMCP(pl *Plan, group, gdir, state string, dry bool, running func() bool) {
	if mcp.Synced(pl.Dir, gdir, state) {
		return
	}
	name := pl.Profile.Name
	names := "the shared MCP servers"
	if store, err := mcp.Load(mcp.Path(gdir)); err == nil && len(store.Servers) > 0 {
		names = "the shared MCP servers (" + strings.Join(store.Names(), ", ") + ")"
	}
	fmt.Fprintf(os.Stderr, "ccam: Remote Control sessions only get MCP servers stored in account %s's .claude.json, not %s.\n", name, names)
	switch {
	case dry:
	case running():
		fmt.Fprintf(os.Stderr, "      Account %s is running, so its .claude.json is left alone. Run this again after it exits to copy them.\n", name)
	case confirm("Copy them into .claude.json (a backup is kept; the account's own servers stay)?"):
		res, err := mcp.Sync(pl.Dir, gdir, state)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ccam: copying the MCP servers failed (starting without them): %v\n", err)
		}
		syncNotes(name, group, res)
		return
	}
	fmt.Fprintln(os.Stderr, "      Starting without them. To copy them before every launch instead, set mcp_strategy = \"sync\" in ~/.config/ccam/config.toml")
}

// syncNotes says what Sync did besides copying.
func syncNotes(name, group string, res mcp.Result) {
	if len(res.Own) > 0 {
		fmt.Fprintf(os.Stderr, "ccam: account %s has its own MCP servers %s, so group %s's servers with those names are not copied there\n"+
			"      To use the group's, remove the account's (claude mcp remove -s user <name>); to share the account's: ccam mcp import %s <name> --force\n",
			name, strings.Join(res.Own, ", "), group, name)
	}
	if len(res.Restored) > 0 {
		fmt.Fprintf(os.Stderr, "ccam: put back group %s's MCP servers that were gone from account %s's .claude.json: %s\n"+
			"      To drop one for every account: ccam mcp remove <name>\n", group, name, strings.Join(res.Restored, ", "))
	}
}

// confirm asks a yes/no question on the terminal; without one the answer is no.
// A variable so tests can answer it.
var confirm = func(question string) bool {
	if !attended() {
		return false
	}
	return askTTY(question)
}

// attended reports whether stdin is a terminal ccam may ask on: not
// /dev/null or NUL, which are character devices too, and not the terminal
// of a job in the background, which reading would stop (SIGTTIN). A
// variable so tests can have a terminal.
var attended = func() bool {
	fd := os.Stdin.Fd()
	return term.IsTerminal(fd) && foreground(fd)
}

// askTTY asks on the terminal; without one it proceeds (the warning is shown).
func askTTY(question string) bool {
	if !attended() {
		return true
	}
	return ask(question)
}

// ask puts question on the terminal and reads the answer: true for y or yes.
func ask(question string) bool {
	fmt.Fprintf(os.Stderr, "%s [y/N] ", question)
	// A terminal in line mode returns one line per read (at most 4096
	// bytes), so a long answer does not spill into claude's input.
	var buf [4096]byte
	n, _ := os.Stdin.Read(buf[:])
	ans := strings.ToLower(strings.TrimSpace(string(buf[:n])))
	return ans == "y" || ans == "yes"
}

// ConfirmMove asks the person whether to move a conversation to an account
// outside its organization pool, which a /ccam:switch fired by a job the
// model scheduled may have named. It asks in the relay's hop too, Unattended
// or not, and only a y or yes it reads moves the conversation: without a
// terminal ccam may ask on, the answer is no. On Linux and macOS, keys typed
// before the question (into Claude Code, or while it ended) are dropped
// first, so they cannot answer it. Windows has no such flush, but there the
// relay carries a conversation on only after the person's own /exit.
func ConfirmMove(question string) bool {
	if !attended() {
		return false
	}
	discardInput(os.Stdin)
	return ask(question)
}

func applyMCP(p paths.Paths, cfg *config.Config, pl *Plan, dry bool) error {
	folder := pl.Profile.Share.Get("mcp")
	if folder == config.Local {
		return nil
	}
	group := folder
	if g := cfg.GroupFor(pl.Profile.Name, "mcp"); g != nil {
		group = g.Name
	}
	gdir := p.GroupDir(folder)
	state := MCPStateDir(p, pl.Profile.Name)
	// Asked only where ccam would rewrite .claude.json, which a running
	// Claude Code writes too.
	running := func() bool { return Running(p, pl.Profile.Name, pl.Dir) }
	if cfg.Strategy() == config.MCPSync {
		if dry {
			return nil
		}
		switch {
		case mcp.SyncFresh(pl.Dir, gdir, state): // nothing changed since the last sync
		case running():
			fmt.Fprintf(os.Stderr, "ccam: account %s is already running, skipping MCP sync\n", pl.Profile.Name)
		default:
			res, err := mcp.Sync(pl.Dir, gdir, state)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ccam: MCP sync failed (launching anyway): %v\n", err)
			} else {
				syncNotes(pl.Profile.Name, group, res)
			}
		}
		if l, err := mcp.NewLocal(pl.Dir, mcp.Path(gdir), filepath.Join(state, "seen")); err == nil {
			l.Shadow = nil // Sync says which servers of the account's own it leaves
			if len(l.User) > 0 {
				// A copy Sync could not take out while the account ran.
				drop, _ := mcp.Copied(pl.Dir, state)
				l.User = slices.DeleteFunc(l.User, func(n string) bool { return slices.Contains(drop, n) })
			}
			localHint(pl.Profile.Name, group, l)
		}
		return nil
	}
	flag := mcp.FlagArgs(gdir)
	if flag != nil {
		switch verb := claudeVerb(pl.Argv[1:]); verb {
		case verbRemoteControl:
			// claude remote-control refuses a --mcp-config before it and takes
			// none after it: its sessions only see the servers in .claude.json.
			flag = nil
			remoteControlMCP(pl, group, gdir, state, dry, running)
		case "attach", "logs", "stop", "kill", "respawn", "rm":
			// A background session has its servers; behind an option claude
			// reads these verbs only with an id, not a session's name.
			flag = nil
		case "mcp list", "mcp get", "mcp remove":
			// They read .claude.json, .mcp.json and plugins only.
			flag = nil
			if !dry {
				mcpVerbNote(pl.Profile.Name, group, gdir, strings.TrimPrefix(verb, "mcp "))
			}
		}
	}
	pl.InsertArgs(flag...)
	if dry {
		return nil
	}
	if l, err := mcp.NewLocal(pl.Dir, mcp.Path(gdir), filepath.Join(state, "seen")); err == nil {
		localHint(pl.Profile.Name, group, l)
	}
	return nil
}

// mcpVerbNote says that `claude mcp list|get|remove` see the account's own
// servers only, not the group's, which its sessions get with --mcp-config.
func mcpVerbNote(name, group, gdir, verb string) {
	store, err := mcp.Load(mcp.Path(gdir))
	if err != nil || len(store.Servers) == 0 {
		return
	}
	fmt.Fprintf(os.Stderr, "ccam: `claude mcp %s` sees only account %s's own MCP servers; its sessions also load group %s's: %s\n"+
		"      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.\n",
		verb, name, group, strings.Join(store.Names(), ", "))
}

// localHint names the servers NewLocal found that only this account has.
func localHint(name, group string, l mcp.Local) {
	if len(l.User) > 0 {
		fmt.Fprintf(os.Stderr, "ccam: MCP servers that only account %s has: %s\n      To share them with other accounts: ccam mcp import %s %s\n",
			name, strings.Join(l.User, ", "), name, strings.Join(l.User, " "))
	}
	if len(l.Folder) > 0 {
		var where, names []string
		for _, folder := range slices.Sorted(maps.Keys(l.Folder)) {
			for _, n := range l.Folder[folder] {
				where = append(where, n+" ("+folder+")")
				if !slices.Contains(names, n) {
					names = append(names, n)
				}
			}
		}
		fmt.Fprintf(os.Stderr, "ccam: MCP servers that account %s has for one folder only (claude mcp add without -s user): %s\n"+
			"      To share them with other accounts, in every folder: ccam mcp import %s %s\n", name, strings.Join(where, ", "), name, strings.Join(names, " "))
	}
	if len(l.Shadow) > 0 {
		fmt.Fprintf(os.Stderr, "ccam: account %s has its own MCP servers %s, which its sessions use instead of group %s's servers with those names\n"+
			"      To use the group's, remove the account's (claude mcp remove -s user <name>); to share the account's: ccam mcp import %s <name> --force\n",
			name, strings.Join(l.Shadow, ", "), group, name)
	}
}
