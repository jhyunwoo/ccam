package cli

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/next"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/session"
	"github.com/jhyunwoo/ccam/internal/usage"
)

func usageCmd(e *env) *cobra.Command {
	var opt usage.Options
	var asJSON, verbose bool
	var watch int
	ended := make(chan os.Signal, 1) // the signal that cancelled the lookups
	cmd := &cobra.Command{
		Use:     "usage",
		Aliases: []string{"u"},
		Short:   "5-hour, weekly and per-model usage of every account, with reset times",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := spacedValue(cmd, args, "watch", "60"); err != nil {
				return err
			}
			return cobra.NoArgs(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if opt.Source != "" && opt.Source != "live" && opt.Source != "cache" && opt.Source != "claude" {
				return fmt.Errorf("--source must be live, claude or cache")
			}
			if e.dryRun && opt.Refresh {
				return fmt.Errorf("--dry-run cannot be combined with --refresh, which rewrites login tokens")
			}
			cfg, err := e.load()
			if err != nil {
				return err
			}
			if len(cfg.Profiles) == 0 {
				if asJSON {
					e.printf("[]\n") // a script still gets JSON
					return nil
				}
				e.printf("No accounts yet. Add one with `ccam add <name>`.\n")
				return nil
			}
			f := e.fetcher()
			watching := cmd.Flags().Changed("watch")
			ctx := context.Background()
			if watching || opt.Refresh {
				// A token renewal in flight must not be cut off (the server may
				// have rotated the refresh token already): with --refresh the
				// signals that end ccam end the lookups instead, and ccam exits
				// once the renewal is saved.
				sigs := []os.Signal{os.Interrupt}
				if opt.Refresh {
					sigs = endSignals
				}
				sig := make(chan os.Signal, 1)
				signal.Notify(sig, sigs...)
				defer signal.Stop(sig)
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				defer cancel()
				go func() {
					select {
					case s := <-sig:
						ended <- s
						cancel()
					case <-ctx.Done():
					}
				}()
			}
			if !watching {
				err := showUsage(ctx, e, f, cfg, opt, asJSON, verbose, false)
				select {
				case s := <-ended:
					if err == nil {
						err = exitError(128 + int(s.(syscall.Signal))) // as if the signal had ended ccam
					}
				default:
				}
				return err
			}
			if watch <= 0 {
				watch = 60
			}
			for {
				if err := showUsage(ctx, e, f, cfg, opt, asJSON, verbose, true); err != nil {
					return err
				}
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(time.Duration(watch) * time.Second):
				}
			}
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "print as JSON")
	cmd.Flags().IntVar(&watch, "watch", 60, "refresh every N seconds")
	cmd.Flags().Lookup("watch").NoOptDefVal = "60"
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show where the numbers came from")
	cmd.Flags().StringVar(&opt.Source, "source", "", "live (ccam asks the usage endpoint), claude (Claude Code does, slower) or cache (saved values) (config: usage_source)")
	cmd.Flags().BoolVar(&opt.NoCache, "no-cache", false, "ignore the 60-second cache")
	cmd.Flags().BoolVar(&opt.Refresh, "refresh", false, "refresh expired tokens (only for accounts not running; a refresh may take up to 30 s)")
	return cmd
}

// newFetcher is replaced in tests to point at a fake server.
var newFetcher = usage.NewFetcher

func (e *env) fetcher() *usage.Fetcher {
	f := newFetcher(usageStateDir(e))
	f.Home = e.paths.Home
	if cfg, err := e.load(); err == nil {
		f.Source = cfg.UsageSource
		f.Probe = claudeProbe(e, cfg)
		if e.dryRun {
			// Claude Code writes to the account's folder and may refresh its
			// login, so a dry run goes by the saved usage.
			f.Probe = func(context.Context, usage.Account) ([]byte, error) {
				return nil, errors.New("a dry run does not start Claude Code; the usage shown is the saved one")
			}
		}
	}
	return f
}

// probeExitWait is how long a probe asked to end may take before it is
// killed, and how long one that has exited may keep its output open:
// Claude Code's shutdown waits 2 s for a renewal in flight and has a 5 s
// failsafe (2.1.308).
var probeExitWait = 10 * time.Second

// claudeProbe runs Claude Code in an account's folder and asks it for the
// usage (usage.ProbeRequest): with --safe-mode, so the account's hooks, MCP
// servers and plugins do not start, and without saving a session.
func claudeProbe(e *env, cfg *config.Config) func(context.Context, usage.Account) ([]byte, error) {
	return func(ctx context.Context, a usage.Account) ([]byte, error) {
		bin, err := launcher.ResolveClaude(e.paths, cfg.ClaudeBin)
		if err != nil {
			return nil, err
		}
		c, err := launcher.ClaudeCommand(ctx, bin, "-p", "--input-format", "stream-json", "--output-format", "stream-json",
			"--verbose", "--safe-mode", "--no-session-persistence")
		if err != nil {
			return nil, err
		}
		// Its own login only, from a folder that is not a project.
		c.Env = launcher.ProfileEnv(e.paths, launcher.WithoutAuth(os.Environ()), a.Dir, a.Name)
		c.Dir = os.TempDir()
		c.Stdin = strings.NewReader(usage.ProbeRequest)
		// A probe given up (its time is over, Ctrl+C in --watch) asks Claude
		// Code to end rather than kills it: on SIGTERM it finishes a token
		// renewal it has started, and one cut off may lose the login.
		// Windows has no such signal.
		if runtime.GOOS != "windows" {
			c.Cancel = func() error { return c.Process.Signal(syscall.SIGTERM) }
		}
		c.WaitDelay = probeExitWait
		return c.Output()
	}
}

func usageStateDir(e *env) string { return filepath.Join(e.paths.State, "usage") }

func usageAccounts(e *env, cfg *config.Config) []usage.Account {
	out := make([]usage.Account, 0, len(cfg.Profiles))
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		out = append(out, usage.Account{
			Name: p.Name, Slot: p.Slot, Dir: profile.Dir(e.paths, p),
			SessionDir: launcher.SessionDir(e.paths, p.Name),
			Imported:   p.Dir != "",
			API:        p.API != nil,
		})
	}
	return out
}

func showUsage(ctx context.Context, e *env, f *usage.Fetcher, cfg *config.Config, opt usage.Options, asJSON, verbose, clear bool) error {
	accts := usageAccounts(e, cfg)
	if asJSON {
		reports := f.FetchAll(ctx, accts, opt, nil)
		enc := json.NewEncoder(e.stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(reports)
	}
	tty := isTerminal(e.stdout)
	ro := usage.RenderOptions{Color: tty && os.Getenv("NO_COLOR") == "", Verbose: verbose}
	view := make([]*usage.Report, len(accts))
	for i, a := range accts {
		view[i] = &usage.Report{Name: a.Name, Slot: a.Slot, Status: usage.StatusPending}
	}
	if clear && tty {
		fmt.Fprint(e.stdout, "\x1b[H\x1b[2J")
	}
	// On a terminal, draw placeholders and redraw as each account answers.
	live := &liveTable{w: e.stdout, tty: tty}
	if f, ok := e.stdout.(*os.File); ok && tty {
		if w, h, err := term.GetSize(f.Fd()); err == nil {
			if w > 1 {
				live.cols = w - 1
			}
			live.rows = h
		}
	}
	draw := func(final bool) {
		ro.Now = time.Now()
		live.draw(usage.Table(view, ro), final)
	}
	if tty {
		draw(false)
	}
	f.FetchAll(ctx, accts, opt, func(i int, r usage.Report) {
		view[i] = &r
		if tty {
			draw(false)
		}
	})
	if !tty || live.held {
		draw(true)
	}
	// The same rule next --auto and pick use, so the advice matches them.
	p, err := choiceFlags{}.policy(cfg, nil)
	if err != nil {
		p = usage.Policy{}
	}
	// next --auto stays in the organization of the last used account and
	// counts running sessions.
	from, _, _ := next.Rotate(candidates(e, cfg, true), "", nil)
	org := newOrgs(e, cfg, false)
	var pool []*usage.Report
	for _, r := range view {
		if prof, ok := cfg.Find(r.Name); ok && org.fits(from, prof) {
			r.Sessions = session.Count(launcher.SessionDir(e.paths, r.Name))
			pool = append(pool, r)
		}
	}
	if best := p.Choose(pool, time.Now()); best != nil && len(view) > 1 {
		if p.Strategy == "" || p.Strategy == usage.Headroom {
			e.printf("→ Most headroom: %s   (open it with ccam next --auto)\n", best.Name)
		} else {
			e.printf("→ next --auto picks %s: %s\n", best.Name, p.Why(best, pool, time.Now()))
		}
	}
	return nil
}

// liveTable draws the usage table and, on a terminal, draws it again in
// place as accounts answer, moving the cursor up over the last draw. Lines
// are cut to the terminal width: a wrapped line would take two rows, and
// moving up by the number of lines would then land too low. A draw as tall
// as the screen scrolls its top out of reach of the next move up, which
// would leave copies of the table in the scrollback; such a draw waits
// (held), and the table is drawn once more at the end.
type liveTable struct {
	w          io.Writer
	tty        bool
	cols, rows int // the width less one and the height; 0 when unknown
	drawn      int // lines of the last draw, which the next one moves over
	held       bool
}

func (t *liveTable) draw(lines []string, final bool) {
	if t.tty && !final && t.rows > 0 && len(lines) >= t.rows {
		t.held = true
		return
	}
	if t.drawn > 0 {
		fmt.Fprintf(t.w, "\x1b[%dA", t.drawn)
	}
	for _, l := range lines {
		if t.tty {
			fmt.Fprint(t.w, "\x1b[2K")
			if t.cols > 0 {
				l = ansi.Truncate(l, t.cols, "…")
			}
		}
		fmt.Fprintln(t.w, l)
	}
	t.drawn = len(lines)
}

// isTerminal reports whether w is a terminal. A character-device check is
// not enough: /dev/null is one, and `ccam add work >/dev/null` must not
// start the wizard.
func isTerminal(w any) bool {
	f, ok := w.(*os.File)
	return ok && term.IsTerminal(f.Fd())
}

func statuslineCmd(e *env) *cobra.Command {
	var wrap, wrap64 string
	cmd := &cobra.Command{
		Use:   "statusline",
		Short: "Claude Code status line: shows the account name and collects usage",
		Long:  "Claude Code runs this as its status line command. Install with `ccam statusline install`.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wrap64 != "" {
				b, err := base64.StdEncoding.DecodeString(wrap64)
				if err != nil {
					return fmt.Errorf("--wrap-base64: %w", err)
				}
				wrap = string(b)
			}
			return runStatusline(e, wrap)
		},
	}
	cmd.Flags().StringVar(&wrap, "wrap", "", "existing status line command (its output gets the account name in front)")
	// What install writes on Windows, where the command goes through Git
	// Bash or PowerShell, which quote differently: base64 needs no quotes.
	cmd.Flags().StringVar(&wrap64, "wrap-base64", "", "--wrap, base64 encoded")
	_ = cmd.Flags().MarkHidden("wrap-base64")
	cmd.AddCommand(statuslineInstallCmd(e))
	return cmd
}

func runStatusline(e *env, wrap string) error {
	input, err := readAll(e.in)
	if err != nil {
		return err
	}
	name := os.Getenv("CCAM_PROFILE")
	if name == "" {
		name = profileByConfigDir(e, os.Getenv("CLAUDE_CONFIG_DIR"))
	}
	limits, model, ok, _ := usage.ParseStatusLine(input)
	if ok && name != "" {
		_ = usage.SaveStatusLine(usageStateDir(e), name, limits, time.Now())
	}
	label := name
	if label == "" {
		label = "claude"
	}
	if wrap != "" {
		// As Claude Code shows a status line: every line that is not
		// blank, trimmed, and nothing when the command failed. The label
		// goes before the first.
		out, err := runShell(wrap, input)
		var rows []string
		if err == nil {
			for _, line := range strings.Split(out, "\n") {
				if line = strings.TrimSpace(line); line != "" {
					rows = append(rows, line)
				}
			}
		}
		if len(rows) == 0 {
			rows = []string{""}
		}
		rows[0] = strings.TrimSpace("[" + label + "] " + rows[0])
		e.printf("%s\n", strings.Join(rows, "\n"))
		return nil
	}
	parts := []string{"[" + label + "]"}
	if model != "" {
		parts = append(parts, model)
	}
	if limits.FiveHour.Known {
		parts = append(parts, fmt.Sprintf("5h %.0f%%", limits.FiveHour.Percent))
	}
	if limits.Weekly.Known {
		parts = append(parts, fmt.Sprintf("7d %.0f%%", limits.Weekly.Percent))
	}
	if nearLimit(limits) && hasOtherAccount(e) {
		parts = append(parts, "→ ccam next --auto --continue")
	}
	e.printf("%s\n", strings.Join(parts, " · "))
	return nil
}

// hintAt is the usage at which the status line suggests the next account.
const hintAt = 90

func nearLimit(l usage.Limits) bool {
	return l.FiveHour.Known && l.FiveHour.Percent >= hintAt || l.Weekly.Known && l.Weekly.Percent >= hintAt
}

func hasOtherAccount(e *env) bool {
	cfg, err := e.load()
	return err == nil && len(cfg.Profiles) > 1
}

func profileByConfigDir(e *env, dir string) string {
	if dir == "" {
		return ""
	}
	cfg, err := e.load()
	if err != nil {
		return ""
	}
	for i := range cfg.Profiles {
		if profile.Dir(e.paths, &cfg.Profiles[i]) == dir {
			return cfg.Profiles[i].Name
		}
	}
	return ""
}
