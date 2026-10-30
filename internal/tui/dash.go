// Package tui is the `ccam` dashboard, the `ccam add` wizard and the
// `ccam migrate` picker.
package tui

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/doctor"
	"github.com/jhyunwoo/ccam/internal/group"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// Backend does the work the dashboard asks for. A method that changes
// something returns the config as it is afterwards.
type Backend interface {
	Usage(ctx context.Context, name string, fresh bool) usage.Report
	Running(name string) bool
	// Sessions is how many ccam sessions of the account run now.
	Sessions(name string) int
	// Takeover lists the accounts that may take over from from by the rules
	// of next --auto (the same organization, not the same Claude account).
	Takeover(from string) []string
	LastUsed(name string) time.Time
	SetDefault(name string) (*config.Config, error)
	SetFolder(dir, account string) (*config.Config, error)
	UnsetFolder(dir string) (*config.Config, error)
	Rename(oldName, newName string) (*config.Config, error)
	// SetShortcut gives account the shortcut command; "" removes it.
	SetShortcut(account, shortcut string) (*config.Config, error)
	SuggestShortcut(account string) string
	CheckShortcut(account, shortcut string) error
	Remove(name string) (*config.Config, string, error)
	// SameLogin names another account (or plain claude) whose folder holds
	// the account's login, which Remove therefore does not log out; or "".
	SameLogin(name string) string
	// RemoveGroup applies ch (the group's accounts leave it), drops the
	// group and moves its folder to the trash.
	RemoveGroup(name string, ch *group.Change) (*config.Config, string, error)
	// PlanGroups builds what edit would move on disk, refusing a change
	// that cannot be applied now; ApplyGroups saves edit and moves.
	PlanGroups(edit func(*config.Config) error) (*group.Change, error)
	ApplyGroups(edit func(*config.Config) error, ch *group.Change) (*config.Config, error)
	Extensions() (ExtView, error)
	// Agents lists every account's Claude Code sessions as `ccam agents`
	// does: each account asked at once, each within its time limit.
	Agents() AgentList
	// StopSession stops account's background session id (claude stop).
	StopSession(account, id string) error
	// Doctor runs the checks; with fix it repairs what is safe, only the
	// findings of account only (doctor.General: of no account) when set.
	Doctor(fix bool, only string) ([]doctor.Finding, error)
	// NewerRelease is the latest ccam release tag when it is newer than this
	// ccam, else "" (also offline). The dashboard asks it each time it opens.
	NewerRelease(ctx context.Context) string
	// Update installs the latest release, as `ccam update` does.
	Update(ctx context.Context) error
}

// ExtView is `ccam ext` prepared for the Extensions tab.
type ExtView struct {
	Kinds    []ExtKind
	Project  []string // lines about the current folder's own extensions
	Problems []string
}

type ExtKind struct {
	Title string
	Rows  []ExtRow
}

// ExtRow is one extension and where it is loaded: Groups are group names,
// Own the accounts that have it in their own folder, Off the accounts that
// have the plugin installed but not turned on.
type ExtRow struct {
	Name    string
	Groups  []string
	Own     []string
	Off     []string
	Problem bool
}

// AgentList is `ccam agents` prepared for the Agents tab: the sessions of
// the accounts that answered, newest first, and a line for each account
// that could not be asked.
type AgentList struct {
	Sessions []Agent
	Failed   []string
}

// Agent is one session. Only a background session has an ID, the one
// `ccam attach` opens.
type Agent struct {
	Account   string
	ID        string
	Name      string
	State     string    // as `ccam agents` shows it: "blocked, idle", "open (busy)"
	Folder    string    // ~ for the home folder
	StartedAt time.Time // zero when unknown
}

// Action is what the user chose that needs the terminal; the caller does
// it after the program has exited.
type Action int

const (
	None Action = iota
	Launch
	Add
	Login
)

type Result struct {
	Action  Action
	Profile string
	Args    []string // for Launch: arguments for claude
	DryRun  bool     // for Launch: print the command only
	Add     AddAnswers
}

const (
	tabAccounts = iota
	tabAgents
	tabGroups
	tabFolders
	tabExtensions
	tabHealth
)

var tabNames = []string{"Accounts", "Agents", "Groups", "Folders", "Extensions", "Health"}

// tabAbout says in one line what each tab is for; it shows under the tabs.
var tabAbout = []string{
	"Your Claude accounts and their usage limits; enter launches the selected one.",
	"Background agents and open sessions of every account; enter attaches to one.",
	"Groups share settings, memory, extensions, MCP or history between accounts.",
	"A folder rule picks the account that plain `ccam` starts in that folder.",
	"Skills, commands, agents, plugins, MCP and hooks, and which accounts see them.",
	"Checks Claude Code and each account's folder, login and links (ccam doctor).",
}

type (
	usageMsg struct {
		name   string
		report usage.Report
	}
	tickMsg struct{}
	doneMsg struct {
		cfg    *config.Config
		notice string
		err    error
	}
	extMsg struct {
		view ExtView
		err  error
	}
	healthMsg struct {
		findings []doctor.Finding
		err      error
		fixed    bool   // the result of f or F (doctor --fix)
		scope    string // what was repaired, for the notice
	}
	agentsMsg struct {
		gen  int // the load it answers (Dashboard.agentsGen)
		list AgentList
	}
	newerMsg   struct{ tag string }
	updatedMsg struct{ err error }
)

// notice is the line above the keys: a result or an error, cleared by the
// next key.
type notice struct {
	text string
	err  bool
}

// accInfo is what the Accounts tab shows beside the config.
type accInfo struct {
	last    time.Time
	running bool
}

// Dashboard is the full-screen `ccam`. Build it with NewDashboard.
type Dashboard struct {
	be        Backend
	cfg       *config.Config
	cwd       string // the current folder as given
	dir       string // cwd with symbolic links resolved, the form folder rules are stored in
	version   string
	goos      string // runtime.GOOS; tests set another
	themeName string
	theme     Theme
	width     int
	height    int
	tab       int
	cursor    [6]int

	usage     map[string]*usage.Report
	info      map[string]accInfo
	frame     int
	ticking   bool
	tickEvery time.Duration // 0 stops the spinner (tests)
	digits    string
	digitAt   time.Time

	groups groupsTab
	agents *AgentList // the Agents tab's list; nil while it loads
	ext    *ExtView
	extErr error
	health []doctor.Finding
	hErr   error
	loaded [6]bool
	newer  string // a newer ccam release, offered with U
	// agentsGen numbers the loads of agents: only the newest one's answer
	// counts, as an account may answer an older one late.
	agentsGen int
	// updated is set once U installed newer; this process still runs the
	// old binary until it is started again.
	updated bool
	// quitAfter is Ctrl+C pressed while busy: ccam quits once the change
	// is done, not halfway through it.
	quitAfter bool

	dlg    dialog
	help   bool
	busy   string // a change is being applied
	notice notice
	result Result
	now    func() time.Time
}

// NewDashboard shows cfg; cwd is the current folder (for ccam folder).
func NewDashboard(be Backend, cfg *config.Config, cwd, version string) *Dashboard {
	d := &Dashboard{be: be, cwd: cwd, dir: cwd, version: version, goos: runtime.GOOS, themeName: cfg.Theme, width: 80, height: 24,
		usage: map[string]*usage.Report{}, now: time.Now, tickEvery: 120 * time.Millisecond}
	if real, err := filepath.EvalSymlinks(cwd); err == nil {
		d.dir = real
	}
	d.theme = ThemeFor(cfg.Theme, true)
	d.setConfig(cfg)
	def, _ := cfg.DefaultProfile()
	for i, p := range cfg.Profiles {
		if def != nil && p.Name == def.Name {
			d.cursor[tabAccounts] = i
		}
	}
	if f, ok := d.here(); ok {
		for i, p := range cfg.Profiles {
			if p.Name == f.Account {
				d.cursor[tabAccounts] = i
			}
		}
	}
	return d
}

// Installed is the release U installed, or "".
func (d *Dashboard) Installed() string {
	if d.updated {
		return d.newer
	}
	return ""
}

// KeepInstalled carries an update an earlier dashboard of this process
// made: the release is not offered again, and the notice says to restart.
func (d *Dashboard) KeepInstalled(tag string) {
	if tag == "" {
		return
	}
	d.newer, d.updated = tag, true
	d.notice = notice{text: "ccam " + tag + " is installed; quit and start ccam again to use it"}
}

// Result is valid after the program has finished.
func (d *Dashboard) Result() Result { return d.result }

// here is the folder rule for the current folder, found as a launch finds
// it (launcher.FolderHere): for the folder as given, else with links
// resolved.
func (d *Dashboard) here() (*config.Folder, bool) {
	if f, ok := d.cfg.FolderFor(d.cwd); ok {
		return f, true
	}
	return d.cfg.FolderFor(d.dir)
}

func (d *Dashboard) setConfig(cfg *config.Config) {
	d.cfg = cfg
	d.info = map[string]accInfo{}
	for _, p := range cfg.Profiles {
		d.info[p.Name] = accInfo{last: d.be.LastUsed(p.Name), running: d.be.Running(p.Name)}
	}
	for name := range d.usage {
		if _, ok := cfg.Find(name); !ok {
			delete(d.usage, name)
		}
	}
	for t := range d.cursor {
		d.cursor[t] = max(0, min(d.cursor[t], d.rowCount(t)-1))
	}
	d.groups.sync(cfg)
	// Tabs that read the disk load again when opened next.
	d.loaded[tabAgents], d.loaded[tabExtensions], d.loaded[tabHealth] = false, false, false
}

func (d *Dashboard) rowCount(tab int) int {
	switch tab {
	case tabAccounts:
		return len(d.cfg.Profiles)
	case tabAgents:
		if d.agents != nil {
			return len(d.agents.Sessions)
		}
	case tabGroups:
		return len(d.cfg.Groups)
	case tabFolders:
		return len(d.cfg.Folders)
	case tabHealth:
		return len(healthBoxes(d.health))
	}
	return 0
}

func (d *Dashboard) Init() tea.Cmd {
	var newer tea.Cmd
	if !d.updated {
		newer = func() tea.Msg { return newerMsg{d.be.NewerRelease(context.Background())} }
	}
	return tea.Batch(tea.RequestBackgroundColor, d.fetchAll(false), newer, d.tick())
}

func (d *Dashboard) fetchAll(fresh bool) tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(d.cfg.Profiles))
	for _, p := range d.cfg.Profiles {
		name := p.Name
		if fresh {
			delete(d.usage, name)
		}
		cmds = append(cmds, func() tea.Msg {
			return usageMsg{name: name, report: d.be.Usage(context.Background(), name, fresh)}
		})
	}
	return tea.Batch(cmds...)
}

// tick runs the spinner while something is loading.
func (d *Dashboard) tick() tea.Cmd {
	if d.ticking || d.tickEvery == 0 {
		return nil
	}
	d.ticking = true
	return tea.Tick(d.tickEvery, func(time.Time) tea.Msg { return tickMsg{} })
}

func (d *Dashboard) loading() bool {
	if d.busy != "" {
		return true
	}
	for _, p := range d.cfg.Profiles {
		if _, ok := d.usage[p.Name]; !ok {
			return true
		}
	}
	return d.tab == tabAgents && d.agents == nil ||
		d.tab == tabExtensions && d.ext == nil && d.extErr == nil ||
		d.tab == tabHealth && d.health == nil && d.hErr == nil
}

func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.width, d.height = msg.Width, msg.Height
		return d, nil
	case tea.BackgroundColorMsg:
		d.theme = ThemeFor(d.themeName, msg.IsDark())
		return d, nil
	case tickMsg:
		d.ticking = false
		d.frame++
		if d.loading() {
			return d, d.tick()
		}
		return d, nil
	case usageMsg:
		r := msg.report
		d.usage[msg.name] = &r
		return d, nil
	case doneMsg:
		d.busy = ""
		if msg.cfg != nil {
			d.setConfig(msg.cfg)
		}
		if msg.err != nil {
			d.notice = notice{text: msg.err.Error(), err: true}
		} else if msg.notice != "" {
			d.notice = notice{text: msg.notice}
		}
		if d.quitAfter {
			return d, tea.Quit
		}
		return d, d.openTab()
	case agentsMsg:
		if msg.gen != d.agentsGen {
			return d, nil
		}
		d.agents = &msg.list
		// The list may be shorter than when the cursor last moved.
		d.cursor[tabAgents] = max(0, min(d.cursor[tabAgents], d.rowCount(tabAgents)-1))
		return d, nil
	case extMsg:
		d.ext, d.extErr = &msg.view, msg.err
		return d, nil
	case newerMsg:
		if !d.updated { // once installed, the release stays the one installed
			d.newer = strings.TrimPrefix(msg.tag, "v") // shown like d.version, which has no v
		}
		return d, nil
	case updatedMsg:
		d.busy = ""
		if msg.err != nil {
			d.notice = notice{text: msg.err.Error(), err: true}
		} else {
			d.updated = true
			d.notice = notice{text: "ccam " + d.newer + " is installed; quit and start ccam again to use it"}
		}
		if d.quitAfter {
			return d, tea.Quit
		}
		return d, nil
	case healthMsg:
		d.health, d.hErr = msg.findings, msg.err
		if d.health == nil && d.hErr == nil {
			d.health = []doctor.Finding{}
		}
		// The boxes may be fewer than when the cursor last moved (or none,
		// while the check ran).
		d.cursor[tabHealth] = max(0, min(d.cursor[tabHealth], d.rowCount(tabHealth)-1))
		if msg.fixed {
			d.busy, d.loaded[tabHealth] = "", true
			if msg.err == nil {
				d.notice = notice{text: "repaired what could be repaired" + msg.scope}
			}
			if d.quitAfter {
				return d, tea.Quit
			}
		}
		return d, nil
	case tea.PasteMsg:
		if d.dlg != nil {
			return d, d.updateDialog(msg)
		}
		return d, nil
	case tea.KeyPressMsg:
		return d, d.key(msg)
	}
	return d, nil
}

// quit ends the program with r.
func (d *Dashboard) quit(r Result) tea.Cmd {
	d.result = r
	return tea.Quit
}

func (d *Dashboard) key(k tea.KeyPressMsg) tea.Cmd {
	if d.busy != "" {
		// A change is being applied; wait for it.
		if k.String() == "ctrl+c" {
			d.quitAfter = true
		}
		return nil
	}
	if k.String() == "ctrl+c" {
		return tea.Quit
	}
	if d.dlg != nil {
		return d.updateDialog(k)
	}
	d.notice = notice{}
	if d.help {
		d.help = false
		return nil
	}
	switch k.String() {
	case "q":
		return tea.Quit
	case "esc":
		if d.tab == tabGroups && d.groups.editing {
			d.groups.editing = false
			return nil
		}
		return tea.Quit
	case "?":
		d.help = true
		return nil
	case "U":
		if d.newer != "" && !d.updated && !(d.tab == tabGroups && d.groups.editing) {
			d.dlg = d.confirmUpdate()
			return nil
		}
	case "tab", "right":
		if d.tab == tabGroups && d.groups.editing && k.String() == "tab" {
			break
		}
		d.tab = (d.tab + 1) % len(tabNames)
		return d.openTab()
	case "shift+tab", "left":
		if d.tab == tabGroups && d.groups.editing && k.String() == "shift+tab" {
			break
		}
		d.tab = (d.tab + len(tabNames) - 1) % len(tabNames)
		return d.openTab()
	case "r":
		d.loaded[tabAgents], d.loaded[tabExtensions], d.loaded[tabHealth] = false, false, false
		d.info = map[string]accInfo{}
		for _, p := range d.cfg.Profiles {
			d.info[p.Name] = accInfo{last: d.be.LastUsed(p.Name), running: d.be.Running(p.Name)}
		}
		return tea.Batch(d.fetchAll(true), d.openTab(), d.tick())
	}
	switch d.tab {
	case tabAccounts:
		return d.accountsKey(k)
	case tabAgents:
		return d.agentsKey(k)
	case tabGroups:
		return d.groupsKey(k)
	case tabFolders:
		return d.foldersKey(k)
	case tabHealth:
		return d.healthKey(k)
	}
	return d.moveCursor(k)
}

// moveCursor handles ↑/↓ for the current tab's rows.
func (d *Dashboard) moveCursor(k tea.KeyPressMsg) tea.Cmd {
	n := d.rowCount(d.tab)
	switch k.String() {
	case "up", "k":
		d.cursor[d.tab] = max(0, d.cursor[d.tab]-1)
	case "down", "j":
		d.cursor[d.tab] = max(0, min(n-1, d.cursor[d.tab]+1))
	case "home", "g":
		d.cursor[d.tab] = 0
	case "end", "G":
		d.cursor[d.tab] = max(0, n-1)
	}
	return nil
}

// openTab loads what the current tab needs from disk the first time.
func (d *Dashboard) openTab() tea.Cmd {
	if d.loaded[d.tab] {
		return nil
	}
	d.loaded[d.tab] = true
	switch d.tab {
	case tabAgents:
		d.agents = nil
		d.agentsGen++
		gen := d.agentsGen
		return tea.Batch(func() tea.Msg { return agentsMsg{gen, d.be.Agents()} }, d.tick())
	case tabExtensions:
		d.ext, d.extErr = nil, nil
		return tea.Batch(func() tea.Msg { v, err := d.be.Extensions(); return extMsg{v, err} }, d.tick())
	case tabHealth:
		d.health, d.hErr = nil, nil
		return tea.Batch(func() tea.Msg { f, err := d.be.Doctor(false, ""); return healthMsg{findings: f, err: err} }, d.tick())
	}
	return nil
}

// run applies a change in the background; the screen shows label meanwhile.
func (d *Dashboard) run(label string, f func() doneMsg) tea.Cmd {
	d.busy = label
	return tea.Batch(func() tea.Msg { return f() }, d.tick())
}

func (d *Dashboard) View() tea.View {
	v := tea.NewView(d.render())
	v.AltScreen = true
	v.WindowTitle = "ccam"
	return v
}

// The dashboard needs minWidth x minHeight: the header, the tabs and the
// key line. In a smaller terminal it says so instead of being cut off.
const minWidth, minHeight = 40, 8

func (d *Dashboard) render() string {
	if d.width < minWidth || d.height < minHeight {
		msg := fmt.Sprintf(" Too small: ccam needs %dx%d (now %dx%d)", minWidth, minHeight, d.width, d.height)
		return ansi.Truncate(msg, max(1, d.width), "")
	}
	t := d.theme
	w := d.width
	var top []string
	var spare []int // lines of top a short screen leaves out
	head := " " + t.bold().Render("ccam") + t.dim().Render(" "+d.version+symDot+plural(len(d.cfg.Profiles), "account")+symDot+plural(len(d.cfg.Groups), "group"))
	if d.tab == tabAccounts {
		legend := t.dim().Render(t.success().Render(symOK)+" default"+symDot+symHere+" this folder"+symDot+symRunning+" running") + " "
		if gap := w - ansi.StringWidth(head) - ansi.StringWidth(legend); gap > 1 {
			head += strings.Repeat(" ", gap) + legend
		}
	}
	top = append(top, head)
	if d.newer != "" && !d.updated {
		top = append(top, " "+t.accent().Render(symItem+" ccam "+d.newer+" is available (you have "+d.version+")")+t.dim().Render(symDot+"press U to update"))
	}
	spare = append(spare, len(top))
	top = append(top, "")
	var tabs []string
	for i, n := range tabNames {
		if i == d.tab {
			tabs = append(tabs, t.badge().Render(n))
		} else {
			tabs = append(tabs, t.dim().Render(" "+n+" "))
		}
	}
	tabLine := " " + strings.Join(tabs, " ")
	if label := d.autoSwitchLabel(); label != "" {
		if gap := w - ansi.StringWidth(tabLine) - ansi.StringWidth(label) - 1; gap > 1 {
			tabLine += strings.Repeat(" ", gap) + label
		}
	}
	top = append(top, tabLine, " "+t.dim().Render(tabAbout[d.tab]), "")
	spare = append(spare, len(top)-2, len(top)-1)

	var dialog []string
	if d.dlg != nil {
		dialog = d.dlg.view(d, w)
	}
	tail := []string{""}
	switch {
	case d.busy != "":
		busy := d.busy
		if d.quitAfter {
			busy += " ccam quits when it is done."
		}
		tail = append(tail, " "+t.accent().Render(spinner[d.frame%len(spinner)])+" "+busy)
	case d.notice.err:
		// An error says what to do next; keep it whole (up to 4 lines).
		lines := wrap(symFail+" "+d.notice.text, w-2)
		for _, l := range lines[:min(len(lines), 4)] {
			tail = append(tail, " "+t.error().Render(l))
		}
	case d.notice.text != "":
		tail = append(tail, " "+t.success().Render(symOK)+" "+d.notice.text)
	}
	tail = append(tail, " "+t.dim().Render(d.keysLine(w-1)))
	// On a short screen the blank lines and the tab's description go
	// first, then the end of a long error, so that the key line stays.
	for i := len(spare) - 1; i >= 0 && len(top)+len(tail) > d.height; i-- {
		top = slices.Delete(top, spare[i], spare[i]+1)
	}
	for len(top)+len(tail) > d.height && len(tail) > 1 {
		tail = slices.Delete(tail, len(tail)-2, len(tail)-1)
	}

	// A dialog taller than the screen leaves no room for the tab, and
	// loses its first lines rather than its choices.
	room := d.height - len(top) - len(dialog) - len(tail)
	if room < 0 {
		dialog = dialog[min(len(dialog), -room):]
		room = 0
	}
	if room < 3 && d.dlg != nil {
		room = 0 // a line or two of the tab says nothing
	}
	lines, cur := d.body(w)
	if d.help {
		lines, cur = d.helpView(w), -1
	}
	lines = window(t, lines, cur, room)
	all := append(append(append(top, lines...), dialog...), tail...)
	for i, l := range all {
		if ansi.StringWidth(l) > w {
			all[i] = ansi.Truncate(l, w-1, symMore)
		}
	}
	return strings.Join(all, "\n")
}

// autoSwitchLabel is what the end of the tab line says of auto_switch:
// on or off, and nothing on Windows, which ignores it (the Health tab
// says so).
func (d *Dashboard) autoSwitchLabel() string {
	switch {
	case d.goos == "windows":
		return ""
	case d.cfg.AutoSwitch:
		return d.theme.accent().Render("auto_switch on")
	}
	return d.theme.dim().Render("auto_switch off")
}

// body is the current tab's lines and the line the cursor is on (-1: none).
func (d *Dashboard) body(w int) ([]string, int) {
	switch d.tab {
	case tabAccounts:
		return d.accountsView(w)
	case tabAgents:
		return d.agentsView(w)
	case tabGroups:
		return d.groupsView(w)
	case tabFolders:
		return d.foldersView(w)
	case tabExtensions:
		return d.extView(w), -1
	}
	return d.healthView(w)
}

// window shows at most room lines of lines, keeping cur in view, with
// "… N more" where lines are cut.
func window(t Theme, lines []string, cur, room int) []string {
	if len(lines) <= room {
		return lines
	}
	if room <= 0 {
		return nil
	}
	start := 0
	if cur >= 0 {
		start = max(0, min(cur-room/2, len(lines)-room))
	}
	end := min(len(lines), start+room)
	out := slices.Clone(lines[start:end])
	if start > 0 {
		out[0] = t.dim().Render(fmt.Sprintf("   ↑ %d more", start+1))
	}
	if end < len(lines) {
		out[len(out)-1] = t.dim().Render(fmt.Sprintf("   ↓ %d more", len(lines)-end+1))
	}
	return out
}

func plural(n int, word string) string {
	if n == 1 {
		return "1 " + word
	}
	return fmt.Sprintf("%d %ss", n, word)
}

// keysLine is the short key list for the current screen. Keys that do not
// fit in w columns are left out, the last ones first; "? for shortcuts"
// stays.
func (d *Dashboard) keysLine(w int) string {
	if d.dlg != nil {
		return d.dlg.keys()
	}
	if d.help {
		return "any key to close"
	}
	var keys []string
	switch d.tab {
	case tabAccounts:
		keys = []string{"enter to launch", "n for most room", "r to refresh", "a to add", "s for groups"}
	case tabAgents:
		keys = []string{"enter to attach", "x to stop", "r to reload"}
	case tabGroups:
		if d.groups.editing {
			keys = []string{"space to toggle", "tab to switch lists", "enter for the plan", "esc to close"}
		} else {
			keys = []string{"enter to edit", "n for a new group", "e to rename", "x to remove"}
		}
	case tabFolders:
		keys = []string{"n to use an account here", "enter to change", "x to remove"}
	case tabExtensions:
		keys = []string{"r to reload"}
	case tabHealth:
		keys = []string{"f to fix the selected box", "F to fix all", "r to check again"}
	}
	if !(d.tab == tabGroups && d.groups.editing) {
		const more = "? for shortcuts"
		for len(keys) > 1 && ansi.StringWidth(strings.Join(append(keys, more), symDot)) > w {
			keys = keys[:len(keys)-1]
		}
		keys = append(keys, more)
	}
	return strings.Join(keys, symDot)
}

func (d *Dashboard) helpView(w int) []string {
	t := d.theme
	rows := [][2]string{
		{"←/→ tab", "switch tabs"}, {"↑/↓ j/k", "move"}, {"r", "refresh"}, {"q esc", "quit"},
	}
	switch d.tab {
	case tabAccounts:
		rows = append(rows, [][2]string{
			{"enter", "launch the account"}, {"1-9", "go to that number (12: 1 2)"},
			{"n", "launch the most room"}, {"o", "launch with options"},
			{"u", "usage details"}, {"a", "add an account"}, {"l", "log in again"}, {"s", "the account's groups"},
			{"d", "make it the default"}, {"f", "use it in this folder"},
			{"e", "rename"}, {"c", "shortcut command"}, {"x", "remove"},
		}...)
	case tabAgents:
		rows = append(rows, [][2]string{{"enter", "attach the background session"}, {"x", "stop the background session"}}...)
	case tabGroups:
		rows = append(rows, [][2]string{
			{"enter", "edit items and accounts"}, {"space", "toggle while editing"},
			{"tab", "items or accounts"}, {"n", "new group"}, {"e", "rename"}, {"x", "remove"},
		}...)
	case tabFolders:
		rows = append(rows, [][2]string{{"n", "use an account here"}, {"enter", "change the account"}, {"x", "remove the rule"}}...)
	case tabHealth:
		rows = append(rows, [][2]string{{"↑/↓", "select a box"}, {"f", "repair the selected box"}, {"F", "repair all (doctor --fix)"}}...)
	}
	if d.newer != "" && !d.updated {
		rows = append(rows, [2]string{"U", "update ccam to " + d.newer})
	}
	cell := func(r [2]string) string { return t.accent().Render(padRight(r[0], 8)) + " " + r[1] }
	out := []string{" " + t.bold().Render("Shortcuts"), ""}
	if w < 72 {
		for _, r := range rows {
			out = append(out, "  "+cell(r))
		}
		return out
	}
	half := (len(rows) + 1) / 2
	colW := (w - 6) / 2
	for i := 0; i < half; i++ {
		line := "  " + padRight(ansi.Truncate(cell(rows[i]), colW, symMore), colW)
		if i+half < len(rows) {
			line += "  " + cell(rows[i+half])
		}
		out = append(out, line)
	}
	return out
}

// confirmUpdate asks before U installs the newer release.
func (d *Dashboard) confirmUpdate() dialog {
	return newConfirm("Update ccam "+d.version+" → "+d.newer+"?",
		"Like `ccam update`: the release is downloaded from GitHub, checked against its SHA-256 checksums and put in place of this ccam in one step. Start ccam again afterwards to use it.",
		"Yes, update", func(d *Dashboard) tea.Cmd {
			d.busy = "Updating ccam to " + d.newer + "…"
			return tea.Batch(func() tea.Msg { return updatedMsg{d.be.Update(context.Background())} }, d.tick())
		})
}
