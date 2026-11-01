package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
)

func ccusageCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "ccusage [-a name[,name…]] [ccusage args]",
		Short: "Run ccusage, the token and cost report, over every account's conversations",
		Long: "Runs ccusage (https://ccusage.com), which reports the tokens and estimated costs in\n" +
			"Claude Code's conversation logs, with CLAUDE_CONFIG_DIR set to the folders of every\n" +
			"account, so one report covers them all. A conversation that accounts share through\n" +
			"their group counts once.\n\n" +
			"  -a, --account name[,name…]  only these accounts (names or numbers). An account\n" +
			"                              that shares its history through a group brings the\n" +
			"                              group's conversations, and ccam says so.\n\n" +
			"ccam's own options come first; from the first other argument on, everything goes to\n" +
			"ccusage (`ccam ccusage -- --help` shows its own help). `blocks` shows 5-hour\n" +
			"windows, which each account has on its own, so name one account for it.\n" +
			"ccusage also reports the other coding tools it finds (Codex and others) unless you\n" +
			"name one, as in `ccam ccusage claude daily`, and reads its settings (ccusage.json)\n" +
			"only from ./.ccusage and the account folders: `--config <file>` names another.\n" +
			"ccam runs an installed ccusage, else `bun x ccusage`, else `npx ccusage@latest`.\n" +
			"Those two start in ccam's own folder, so the folder you are in cannot choose what\n" +
			"they run. ./.ccusage/ccusage.json still applies: ccam passes it to them, and a\n" +
			"relative --config, by its full path. ccusage then reads other relative paths,\n" +
			"such as a --pi-path, from ccam's folder, so give those in full.",
		Example: "  ccam ccusage\n  ccam ccusage daily --breakdown\n  ccam ccusage claude monthly --json\n" +
			"  ccam ccusage -a work blocks --active\n  ccam ccusage -a work,side weekly",
		// The options are ccusage's, so cobra parses none: ccam reads its
		// own up to the first argument that is not one of them.
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var sels []string
			add := func(flag, names string) error {
				n := len(sels)
				for _, s := range strings.Split(names, ",") {
					if s = strings.TrimSpace(s); s != "" {
						sels = append(sels, s)
					}
				}
				if len(sels) == n {
					return fmt.Errorf("%s needs account names, as in `ccam ccusage -a work daily`", flag)
				}
				return nil
			}
		options:
			for len(args) > 0 {
				switch a := args[0]; {
				case a == "--":
					args = args[1:]
					break options
				case a == "-h" || a == "--help":
					return cmd.Help()
				case a == "--dry-run": // `ccam --dry-run ccusage` puts it here too
					e.dryRun, args = true, args[1:]
				case a == "-a" || a == "--account":
					names := ""
					if len(args) > 1 {
						names = args[1]
					}
					if err := add(a, names); err != nil {
						return err
					}
					args = args[2:]
				case strings.HasPrefix(a, "--account=") || strings.HasPrefix(a, "-a="):
					flag, names, _ := strings.Cut(a, "=")
					if err := add(flag, names); err != nil {
						return err
					}
					args = args[1:]
				default:
					break options
				}
			}
			cfg, err := e.load()
			if err != nil {
				return err
			}
			dirs, err := ccusageDirs(e, cfg, sels)
			if err != nil {
				return err
			}
			bin, argv, err := ccusageCommand()
			if err != nil {
				return err
			}
			// bun x and npx take the program they run (node_modules/.bin, a
			// package.json bin) and the registry and answers they fetch it
			// with (.npmrc) from the folder they start in and the folders
			// above it, so a repository someone else wrote would choose them.
			// They start in ccam's own folder instead, and get the settings
			// ccusage reads in the working folder by their full path.
			runnerDir := ""
			if argv[0] != "ccusage" {
				runnerDir = filepath.Join(e.paths.State, "ccusage")
				if args, err = ccusageConfigArgs(args); err != nil {
					return err
				}
			}
			argv = append(argv, args...)
			env := slices.DeleteFunc(os.Environ(), func(kv string) bool { return strings.HasPrefix(kv, "CLAUDE_CONFIG_DIR=") })
			env = append(env, "CLAUDE_CONFIG_DIR="+strings.Join(dirs, ","))
			if e.dryRun {
				// Printed as the dry run of a launch is: the variable, then the
				// command, in the folder a runner starts in. It fails only
				// where ExecProgram would.
				if err := (&launcher.Plan{Bin: bin, Argv: argv, Env: env, WorkDir: runnerDir}).DryRun(e.stdout); err != nil {
					return fmt.Errorf("%s is a batch file, and cmd.exe cannot be given an argument with a line break", bin)
				}
				return nil
			}
			if runnerDir != "" {
				if env, err = enterRunnerFolder(runnerDir, env); err != nil {
					return err
				}
			}
			return launcher.ExecProgram(bin, argv, env)
		},
	}
}

// ccusageDirs is the folders ccusage reads for the accounts sels names, in
// that order (every account, in the config's order, when there is none):
// those with conversations, as ccusage stops when none of its folders has a
// projects folder.
func ccusageDirs(e *env, cfg *config.Config, sels []string) ([]string, error) {
	if len(cfg.Profiles) == 0 {
		return nil, errors.New("no accounts yet. Add one first with `ccam add <name>`")
	}
	var chosen []*config.Profile
	for _, s := range sels {
		p, err := cfg.MustFind(s)
		if err != nil {
			return nil, err
		}
		if !slices.Contains(chosen, p) {
			chosen = append(chosen, p)
		}
	}
	if len(sels) == 0 {
		for i := range cfg.Profiles {
			chosen = append(chosen, &cfg.Profiles[i])
		}
	}
	var dirs, names []string
	var noted []os.FileInfo // shared folders already named
	for _, p := range chosen {
		names = append(names, p.Name)
		dir := profile.Dir(e.paths, p)
		fi, err := os.Stat(filepath.Join(dir, "projects"))
		if err != nil || !fi.IsDir() {
			continue
		}
		if strings.Contains(dir, ",") {
			return nil, fmt.Errorf("the folder of %s (%s) has a comma in its path, which ccusage reads as two folders; leave the account out with -a", p.Name, dir)
		}
		if len(sels) > 0 && !slices.ContainsFunc(noted, func(n os.FileInfo) bool { return os.SameFile(n, fi) }) {
			if others := sharers(e, cfg, chosen, fi); len(others) > 0 {
				fmt.Fprintf(e.stderr, "ccam: %s shares its conversation history with %s, so the report counts their conversations too\n", p.Name, andList(others))
				noted = append(noted, fi)
			}
		}
		dirs = append(dirs, dir)
	}
	switch {
	case len(dirs) > 0:
		return dirs, nil
	case len(sels) == 0:
		return nil, errors.New("no account has Claude Code conversations yet")
	case len(names) == 1:
		return nil, fmt.Errorf("%s has no Claude Code conversations yet", names[0])
	}
	return nil, fmt.Errorf("%s have no Claude Code conversations yet", andList(names))
}

// sharers is the accounts left out of chosen whose projects folder is fi:
// a group that shares history links its folder into every member.
func sharers(e *env, cfg *config.Config, chosen []*config.Profile, fi os.FileInfo) []string {
	var names []string
	for i := range cfg.Profiles {
		q := &cfg.Profiles[i]
		if slices.Contains(chosen, q) {
			continue
		}
		if qi, err := os.Stat(filepath.Join(profile.Dir(e.paths, q), "projects")); err == nil && os.SameFile(fi, qi) {
			names = append(names, q.Name)
		}
	}
	return names
}

// ccusageCommand is how ccusage starts: an installed ccusage, else through
// bun or npm, as ccusage's guide runs it. bin is the program's path, argv
// its arguments from its name on.
func ccusageCommand() (bin string, argv []string, err error) {
	var drives []string
	if runtime.GOOS == "linux" {
		drives = windowsDrives()
	}
	windows := ""
	for _, try := range [][]string{{"ccusage"}, {"bun", "x", "ccusage"}, {"npx", "ccusage@latest"}} {
		for _, bin := range lookPaths(try[0], len(drives) > 0) {
			if onDrive(bin, drives) {
				if windows == "" {
					windows = bin
				}
				continue
			}
			return bin, try, nil
		}
	}
	if windows != "" {
		return "", nil, fmt.Errorf("%s is a Windows program, which would not get the account folders: install Node.js or Bun inside WSL, or ccusage itself (https://ccusage.com)", windows)
	}
	return "", nil, errors.New("ccusage is not installed, and neither bun nor npx is there to run it: install Node.js or Bun, or ccusage itself (https://ccusage.com)")
}

// ccusageConfigArgs is args for a runner that starts outside the working
// folder: a --config relative to that folder made absolute and, when no
// --config is given, the folder's .ccusage/ccusage.json named with one, as
// ccusage reads only the first settings file it finds and looks there
// first. Nothing after -- is an option, and a --config without a value is
// left for ccusage to report.
func ccusageConfigArgs(args []string) ([]string, error) {
	out := slices.Clone(args)
	end := slices.Index(out, "--")
	if end < 0 {
		end = len(out)
	}
	named := false
	for i := 0; i < end; i++ {
		prefix := "--config="
		if out[i] == "--config" {
			if i+1 == end || strings.HasPrefix(out[i+1], "-") {
				named = true
				continue // no value, which ccusage reports
			}
			i, prefix = i+1, ""
		} else if !strings.HasPrefix(out[i], prefix) {
			continue
		}
		named = true
		path := strings.TrimPrefix(out[i], prefix)
		if path == "" || filepath.IsAbs(path) {
			continue
		}
		abs, err := fromHere(path)
		if err != nil {
			return nil, err
		}
		out[i] = prefix + abs
	}
	if named {
		return out, nil
	}
	settings, err := filepath.Abs(filepath.Join(".ccusage", "ccusage.json"))
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(settings); err != nil {
		return out, nil
	}
	return slices.Insert(out, end, "--config", settings), nil
}

// fromHere is path, relative to the working folder, as a full path that
// names the same file from any other folder. On Windows, which takes .. by
// name, that is filepath.Abs, which also knows paths relative to a drive or
// its root. Elsewhere the path is put after the folder's as it is, not
// cleaned: a .. after a linked folder leads out of the folder the link
// names, which only the system follows.
func fromHere(path string) (string, error) {
	if runtime.GOOS == "windows" {
		return filepath.Abs(path)
	}
	here, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(here, "/") + "/" + path, nil
}

// enterRunnerFolder makes dir, ccam's own folder for bun x and npx, the
// working folder, which ExecProgram's program starts in, and returns
// environ with PWD naming it, so that nothing reads the folder ccam started
// in from there.
func enterRunnerFolder(dir string, environ []string) ([]string, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	if err := os.Chdir(dir); err != nil {
		return nil, err
	}
	environ = slices.DeleteFunc(environ, func(kv string) bool { return strings.HasPrefix(kv, "PWD=") })
	return append(environ, "PWD="+dir), nil
}

// lookPaths is where name is found on PATH: exec.LookPath's answer, or with
// all every match in PATH's order, since under WSL a Windows program may
// come before a Linux one.
func lookPaths(name string, all bool) []string {
	if !all {
		if bin, err := exec.LookPath(name); err == nil {
			return []string{bin}
		}
		return nil
	}
	var found []string
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if !filepath.IsAbs(dir) {
			continue // like exec.LookPath, never through a relative folder
		}
		if bin, err := exec.LookPath(filepath.Join(dir, name)); err == nil {
			found = append(found, bin)
		}
	}
	return found
}

// mountsFile lists the mounts; a variable for tests.
var mountsFile = "/proc/self/mounts"

// windowsDrives is where WSL mounts the Windows drives: drvfs (WSL 1), 9p
// with aname=drvfs, or virtiofs from a drvfs… tag (WSL 2). WSL appends their
// folders to PATH, and a program there runs as a Windows one, which gets
// none of the variables ccam sets unless WSLENV names them.
func windowsDrives() []string {
	data, err := os.ReadFile(mountsFile)
	if err != nil {
		return nil
	}
	var drives []string
	for _, line := range strings.Split(string(data), "\n") {
		f := strings.Fields(line)
		if len(f) < 4 || f[2] != "drvfs" && !(f[2] == "9p" && strings.Contains(f[3], "aname=drvfs")) &&
			!(f[2] == "virtiofs" && strings.HasPrefix(f[0], "drvfs")) {
			continue
		}
		drives = append(drives, strings.TrimSuffix(mountPath(f[1]), "/"))
	}
	return drives
}

// onDrive reports whether path, with its links resolved, is on one of drives.
func onDrive(path string, drives []string) bool {
	if len(drives) == 0 {
		return false
	}
	if real, err := filepath.EvalSymlinks(path); err == nil {
		path = real
	}
	for _, d := range drives {
		if path == d || strings.HasPrefix(path, d+"/") {
			return true
		}
	}
	return false
}

// mountPath undoes the octal escapes of a mount point in /proc/self/mounts
// (\040 for a space).
func mountPath(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+3 < len(s) && isOctal(s[i+1]) && isOctal(s[i+2]) && isOctal(s[i+3]) {
			b.WriteByte((s[i+1]-'0')<<6 | (s[i+2]-'0')<<3 | (s[i+3] - '0'))
			i += 3
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func isOctal(c byte) bool { return c >= '0' && c <= '7' }
