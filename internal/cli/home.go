package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/doctor"
	"github.com/jhyunwoo/ccam/internal/ext"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/group"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/session"
	"github.com/jhyunwoo/ccam/internal/shim"
	"github.com/jhyunwoo/ccam/internal/tui"
	"github.com/jhyunwoo/ccam/internal/update"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// interactive reports whether both ends are a terminal.
func interactive(e *env) bool { return isTerminal(e.stdin) && isTerminal(e.stdout) }

// runHome is bare `ccam`: the dashboard on a terminal, the list otherwise.
// Launching ends it; adding an account or logging in leaves the dashboard
// for the terminal and comes back to it.
func runHome(e *env) error {
	if !interactive(e) {
		return listProfiles(e, false)
	}
	installed := "" // a release U installed: this process still runs the old one
	for {
		cfg, err := e.load()
		if err != nil {
			return err
		}
		cwd, _ := os.Getwd()
		// Nothing may print over the dashboard: the backend's notes go nowhere.
		quiet := *e
		quiet.stdout, quiet.stderr = io.Discard, io.Discard
		renew := newRenewals()
		dash := tui.NewDashboard(backend{e: &quiet, renew: renew}, cfg, cwd, version())
		dash.KeepInstalled(installed)
		final, err := tea.NewProgram(dash, tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout)).Run()
		renew.end()
		if err != nil {
			return err
		}
		installed = final.(*tui.Dashboard).Installed()
		res := final.(*tui.Dashboard).Result()
		// The dashboard may have been open for a while; act on the current config.
		if cfg, err = e.load(); err != nil {
			return err
		}
		switch res.Action {
		case tui.Launch:
			return launcher.LaunchWith(e.paths, cfg, res.Profile, res.Args, e.dryRun || res.DryRun, e.stdout)
		case tui.Add:
			err := errDryRun
			if !e.dryRun {
				err = addFromAnswers(e, res.Add)
			}
			if err != nil {
				fmt.Fprintln(e.stderr, "ccam:", err)
			}
		case tui.Login:
			if e.dryRun {
				fmt.Fprintln(e.stderr, "ccam:", errDryRun)
			} else if p, err := loginProfile(e, res.Profile); err != nil {
				fmt.Fprintln(e.stderr, "ccam:", err)
			} else {
				e.printf("✓ %s logged in: %s\n", p.Name, p.Email)
			}
		default:
			return nil
		}
		e.printf("Press Enter to return to ccam.")
		_, _ = e.lines().readLine()
	}
}

// addInteractive runs the wizard, then the same steps as `ccam add`.
func addInteractive(e *env, cfg *config.Config, name string) error {
	check := func(name, shortcut string) error { return checkShortcut(e, cfg, name, shortcut) }
	ans, err := tui.RunAddWizard(cfg, name, shim.Suggest, check)
	if err != nil {
		return err
	}
	return addFromAnswers(e, ans)
}

func addFromAnswers(e *env, a tui.AddAnswers) error {
	if a.API != nil {
		if _, err := ReadKeyForCheck(a.API); err != nil {
			fmt.Fprintf(e.stderr, "Note: %v. Set it before starting the account.\n", err)
		}
	}
	if a.Login {
		e.printf("→ Log in in your browser. ccam continues when you are done.\n")
	}
	return addProfile(e, newAccount{name: a.Name, groups: a.Groups, items: a.Items, login: a.Login, shortcut: a.Shortcut, api: a.API})
}

// backend carries out what the dashboard asks for, with the same code as
// the commands.
type backend struct {
	e     *env
	renew *renewals // nil: r renews no login
}

// renewals are the lookups the dashboard's r makes, which renew an expired
// login as `ccam usage --refresh` does. A renewal cut off once sent may lose
// the login (the server may have rotated the refresh token), so runHome ends
// them before it goes on: lookups not yet renewing stop, and end waits for a
// renewal in flight.
type renewals struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func newRenewals() *renewals {
	r := &renewals{}
	r.ctx, r.cancel = context.WithCancel(context.Background())
	return r
}

// start reports whether a renewing lookup may begin; wg.Done must follow.
func (r *renewals) start() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.ctx.Err() != nil {
		return false
	}
	r.wg.Add(1)
	return true
}

func (r *renewals) end() {
	r.mu.Lock()
	r.cancel()
	r.mu.Unlock()
	r.wg.Wait()
}

var errDryRun = errors.New("this is a dry run (ccam --dry-run): nothing was changed")

func (b backend) changing() error {
	if b.e.dryRun {
		return errDryRun
	}
	return nil
}

func (b backend) Usage(ctx context.Context, name string, fresh bool) usage.Report {
	cfg, err := b.e.load()
	if err == nil {
		for _, a := range usageAccounts(b.e, cfg) {
			if a.Name == name {
				opt := usage.Options{NoCache: fresh}
				// r renews an expired login too, except in a dry run.
				if fresh && !b.e.dryRun && b.renew != nil && b.renew.start() {
					defer b.renew.wg.Done()
					ctx, opt.Refresh = b.renew.ctx, true
				}
				return b.e.fetcher().Fetch(ctx, a, opt)
			}
		}
		err = fmt.Errorf("no account %q", name)
	}
	return usage.Report{Name: name, Status: usage.StatusError, Err: err.Error()}
}

func (b backend) Running(name string) bool {
	cfg, err := b.e.load()
	if err != nil {
		return false
	}
	p, ok := cfg.Find(name)
	return ok && launcher.Running(b.e.paths, p.Name, profile.Dir(b.e.paths, p))
}
func (b backend) Sessions(name string) int {
	return session.Count(launcher.SessionDir(b.e.paths, name))
}

func (b backend) Takeover(from string) []string {
	cfg, err := b.e.load()
	if err != nil {
		return nil
	}
	return takeover(b.e, cfg, from)
}
func (b backend) LastUsed(name string) time.Time { return launcher.LastUsed(b.e.paths, name) }

func (b backend) SuggestShortcut(name string) string { return shim.Suggest(name) }

func (b backend) CheckShortcut(name, shortcut string) error {
	cfg, err := b.e.load()
	if err != nil {
		return err
	}
	return checkShortcut(b.e, cfg, name, shortcut)
}

func (b backend) update(change func(*config.Config) error) (*config.Config, error) {
	if err := b.changing(); err != nil {
		return nil, err
	}
	return b.e.update(change)
}

func (b backend) SetDefault(name string) (*config.Config, error) {
	return b.update(func(c *config.Config) error {
		p, err := c.MustFind(name)
		if err == nil {
			c.Default = p.Name
		}
		return err
	})
}

func (b backend) SetFolder(dir, account string) (*config.Config, error) {
	return b.update(func(c *config.Config) error {
		p, err := c.MustFind(account)
		if err != nil {
			return err
		}
		abs, err := filepath.Abs(dir)
		if err != nil {
			return err
		}
		c.SetFolder(abs, p.Name)
		return nil
	})
}

func (b backend) UnsetFolder(dir string) (*config.Config, error) {
	return b.update(func(c *config.Config) error {
		if !c.UnsetFolder(dir) {
			return fmt.Errorf("no rule for %s", dir)
		}
		return nil
	})
}

func (b backend) Rename(oldName, newName string) (*config.Config, error) {
	if err := b.changing(); err != nil {
		return nil, err
	}
	if _, _, err := renameAccount(b.e, oldName, newName); err != nil {
		return nil, err
	}
	return b.e.load()
}

func (b backend) SetShortcut(account, shortcut string) (*config.Config, error) {
	if err := b.changing(); err != nil {
		return nil, err
	}
	cfg, prof, err := b.e.loadProfile(account)
	if err != nil {
		return nil, err
	}
	if shortcut != "" {
		if err := checkShortcut(b.e, cfg, prof.Name, shortcut); err != nil {
			return nil, err
		}
	}
	old := prof.Shortcut
	prof.Shortcut = shortcut
	if err := saveShortcut(b.e, cfg, old, shortcut); err != nil {
		return nil, err
	}
	return b.e.load()
}

func (b backend) SameLogin(name string) string {
	cfg, p, err := b.e.loadProfile(name)
	if err != nil {
		return ""
	}
	return sameLogin(b.e, cfg, p)
}

func (b backend) Remove(name string) (*config.Config, string, error) {
	if err := b.changing(); err != nil {
		return nil, "", err
	}
	msg, err := removeAccount(b.e, name) // refuses a running account
	if err != nil {
		return nil, "", err
	}
	cfg, err := b.e.load()
	return cfg, msg, err
}

func (b backend) PlanGroups(edit func(*config.Config) error) (*group.Change, error) {
	old, err := b.e.load()
	if err != nil {
		return nil, err
	}
	if err := group.Unfinished(b.e.paths); err != nil {
		return nil, err
	}
	next := old.Clone()
	if err := edit(next); err != nil {
		return nil, err
	}
	ch, err := group.Build(b.e.paths, old, next, group.Options{})
	if err != nil {
		return nil, err
	}
	if b.e.dryRun {
		return ch, nil
	}
	return ch, ch.Check(group.Running(b.e.paths))
}

func (b backend) ApplyGroups(edit func(*config.Config) error, ch *group.Change) (*config.Config, error) {
	// The plan dialog may have stayed open while an account started or
	// another ccam changed the groups.
	cfg, err := b.update(func(c *config.Config) error { return recheck(ch, c, edit, group.Running(b.e.paths)) })
	if err != nil {
		ch.Abort()
		return nil, err
	}
	if err := ch.Apply(); err != nil {
		return cfg, fmt.Errorf("failed while applying: %w", err)
	}
	return cfg, nil
}

func (b backend) RemoveGroup(name string, ch *group.Change) (*config.Config, string, error) {
	old, err := b.e.load()
	if err != nil {
		return nil, "", err
	}
	g, err := old.MustFindGroup(name)
	if err != nil {
		return nil, "", err
	}
	folder := g.Folder()
	trash := len(old.FolderUsers(folder)) == 1 && fsutil.Exists(b.e.paths.GroupDir(folder))
	cfg, err := b.ApplyGroups(func(c *config.Config) error { return c.RemoveGroup(name) }, ch)
	if err != nil {
		return cfg, "", err
	}
	if !trash {
		return cfg, "group " + name + " removed", nil
	}
	dst := fsutil.UniqueDir(filepath.Join(b.e.paths.TrashDir(), "group-"+folder+"-"+fsutil.Stamp(time.Now())))
	if err := fsutil.Move(b.e.paths.GroupDir(folder), dst); err != nil {
		return cfg, "", fmt.Errorf("group %s is removed, but its folder could not go to the trash: %w", name, err)
	}
	return cfg, "group " + name + " removed; its folder is in " + dst, nil
}

func (b backend) NewerRelease(ctx context.Context) string { return newerRelease(ctx) }

func (b backend) Update(ctx context.Context) error {
	if err := b.changing(); err != nil {
		return err
	}
	// The dashboard says what happened; the steps' lines are not shown.
	return b.e.runUpdate(ctx, update.Options{}, io.Discard)
}

func (b backend) Doctor(fix bool, only string) ([]doctor.Finding, error) {
	cfg, err := b.e.load()
	if err != nil {
		return nil, err
	}
	if fix {
		if err := b.changing(); err != nil {
			return nil, err
		}
	}
	self, _ := ccamPath()
	// No live usage check: the dashboard looks usage up itself.
	return doctor.Run(b.e.paths, cfg, doctor.Options{Fix: fix, FixOnly: only, CcamPath: self, Reload: b.e.load}), nil
}

func (b backend) Agents() tui.AgentList {
	cfg, err := b.e.load()
	if err != nil {
		return tui.AgentList{Failed: []string{err.Error()}}
	}
	list, failed := listAgents(b.e, cfg, false)
	home, _ := os.UserHomeDir()
	out := tui.AgentList{Failed: failed}
	for _, a := range list {
		var started time.Time
		if a.StartedAt > 0 {
			started = time.UnixMilli(a.StartedAt)
		}
		out.Sessions = append(out.Sessions, tui.Agent{Account: a.Account, ID: a.ID, Name: a.Name,
			State: a.stateText(), Folder: tildeDir(a.Cwd, home), StartedAt: started})
	}
	return out
}

func (b backend) StopSession(account, id string) error {
	if err := b.changing(); err != nil {
		return err
	}
	return stopSession(b.e, account, id)
}

func (b backend) Extensions() (tui.ExtView, error) {
	cfg, err := b.e.load()
	if err != nil {
		return tui.ExtView{}, err
	}
	profs := make([]*config.Profile, 0, len(cfg.Profiles))
	for i := range cfg.Profiles {
		profs = append(profs, &cfg.Profiles[i])
	}
	cwd, _ := os.Getwd()
	return extView(scanExt(b.e, cfg, profs, cwd)), nil
}

// extView sorts a scan for the Extensions tab: each extension once, with
// the groups it comes from and the accounts that have their own.
func extView(rep extReport) tui.ExtView {
	var v tui.ExtView
	where := func(row *tui.ExtRow, account, source string) {
		switch {
		case source == ext.Account || source == ext.Local:
			row.Own = append(row.Own, account)
		case !slices.Contains(row.Groups, source):
			row.Groups = append(row.Groups, source)
		}
	}
	kinds := []struct {
		title string
		pick  func(ext.Set) []ext.Item
	}{
		{"Skills", func(s ext.Set) []ext.Item { return s.Skills }},
		{"Commands", func(s ext.Set) []ext.Item { return s.Commands }},
		{"Agents", func(s ext.Set) []ext.Item { return s.Agents }},
		{"Plugins", func(s ext.Set) []ext.Item { return s.Plugins }},
		{"MCP servers", func(s ext.Set) []ext.Item { return s.MCP }},
	}
	for _, k := range kinds {
		kind := tui.ExtKind{Title: k.title}
		for _, n := range ext.Names(rep.Accounts, k.pick) {
			row := tui.ExtRow{Name: n}
			for _, s := range rep.Accounts {
				for _, it := range k.pick(s) {
					switch {
					case it.Name != n:
					case it.Off:
						row.Problem = row.Problem || it.Problem != ""
						row.Off = append(row.Off, s.Name)
					default:
						row.Problem = row.Problem || it.Problem != ""
						where(&row, s.Name, it.Source)
					}
				}
			}
			kind.Rows = append(kind.Rows, row)
		}
		v.Kinds = append(v.Kinds, kind)
	}
	hooks := tui.ExtKind{Title: "Hooks"}
	for _, r := range hookTable(rep.Accounts) {
		row := tui.ExtRow{Name: strings.TrimSpace(r[0])}
		for _, s := range rep.Accounts {
			for _, h := range s.Hooks {
				if h.Event == row.Name {
					where(&row, s.Name, h.Source)
				}
			}
		}
		hooks.Rows = append(hooks.Rows, row)
	}
	v.Kinds = append(v.Kinds, hooks)
	p := rep.Project
	for _, part := range []struct {
		label string
		items []ext.Item
	}{{"skills", p.Skills}, {"commands", p.Commands}, {"agents", p.Agents}, {"MCP servers", p.MCP}} {
		if len(part.items) > 0 {
			v.Project = append(v.Project, part.label+": "+joinNames(part.items))
		}
	}
	for _, h := range p.Hooks {
		v.Project = append(v.Project, fmt.Sprintf("hook %s: %d", h.Event, h.Count))
	}
	for _, pr := range rep.Problems {
		switch {
		case pr.Group != "":
			v.Problems = append(v.Problems, "group "+pr.Group+": "+pr.Problem)
		case pr.Account != "":
			v.Problems = append(v.Problems, pr.Account+": "+pr.Problem)
		default:
			v.Problems = append(v.Problems, pr.Problem)
		}
	}
	return v
}
