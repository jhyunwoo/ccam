package tui

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/group"
	"github.com/jhyunwoo/ccam/internal/usage"
)

const (
	cellBarW = 6
	timeW    = 7         // "12h 30m"
	cellW    = 5 + timeW // "100% 12h 30m"; a wide terminal adds the bar: "██████ 100% 12h 30m"
	lastW    = 4
)

func (d *Dashboard) current() *config.Profile {
	if len(d.cfg.Profiles) == 0 {
		return nil
	}
	return &d.cfg.Profiles[d.cursor[tabAccounts]]
}

func emailOf(p *config.Profile) string {
	switch {
	case p.API != nil:
		return p.API.Label()
	case p.Email == "":
		return "-"
	}
	return p.Email
}

// shortAgo is the Last column: now, 5m, 2h, 3d.
func shortAgo(t, now time.Time) string {
	if t.IsZero() {
		return "-"
	}
	switch d := now.Sub(t); {
	case d < time.Minute:
		return "now"
	case d < time.Hour:
		return strconv.Itoa(int(d.Minutes())) + "m"
	case d < 24*time.Hour:
		return strconv.Itoa(int(d.Hours())) + "h"
	default:
		return strconv.Itoa(int(d.Hours()/24)) + "d"
	}
}

func (d *Dashboard) accountsView(w int) ([]string, int) {
	t := d.theme
	if len(d.cfg.Profiles) == 0 {
		return []string{" No accounts yet. Press " + t.accent().Render("a") + " to add one."}, -1
	}
	nameW, emailW, slotW := len("Account"), len("Email"), 2
	for _, p := range d.cfg.Profiles {
		nameW = max(nameW, ansi.StringWidth(p.Name))
		emailW = max(emailW, ansi.StringWidth(emailOf(&p)))
		slotW = max(slotW, len(strconv.Itoa(p.Slot))+1)
	}
	nameW, emailW = min(nameW, 20), min(emailW, 24)
	const groupsMin = 12
	// 5-hour, weekly and Fable. Each cell always says when it resets; its
	// bar needs a wide terminal.
	cw := cellBarW + 1 + cellW
	bars := 1+2+slotW+1+nameW+5+3*cw+4+2+groupsMin <= w
	if !bars {
		cw = cellW
	}
	usageW := 3*cw + 4
	fixed := 1 + 2 + slotW + 1 + nameW + 5 + usageW + 2
	showEmail := fixed+emailW+2+lastW+2+groupsMin <= w
	showLast := fixed+lastW+2+groupsMin <= w || showEmail
	head := " " + strings.Repeat(" ", 2+slotW-1) + "# " + padRight("Account", nameW) + "      "
	if showEmail {
		head += padRight("Email", emailW) + "  "
	}
	head += padRight("5-hour", cw) + "  " + padRight("Weekly", cw) + "  " + padRight("Fable", cw) + "  "
	if showLast {
		head += padRight("Last", lastW) + "  "
	}
	head += "Groups"
	out := []string{t.dim().Render(head)}
	def, _ := d.cfg.DefaultProfile()
	here, _ := d.here()
	now := d.now()
	for i := range d.cfg.Profiles {
		p := &d.cfg.Profiles[i]
		sel := i == d.cursor[tabAccounts]
		cur, name := "  ", padRight(ansi.Truncate(p.Name, nameW, symMore), nameW)
		if sel {
			cur, name = t.accent().Render(symCursor)+" ", t.accent().Render(name)
		}
		marks := ""
		for _, m := range []struct {
			on  bool
			sym string
		}{{def != nil && def.Name == p.Name, t.success().Render(symOK)}, {here != nil && here.Account == p.Name, symHere}, {d.info[p.Name].running, t.accent().Render(symRunning)},
			{d.usage[p.Name] != nil && d.usage[p.Name].HoldsReset, t.accent().Render(symReset)}} {
			if m.on {
				marks += m.sym
			} else {
				marks += " "
			}
		}
		line := " " + cur + t.dim().Render(padLeft(strconv.Itoa(p.Slot)+".", slotW)) + " " + name + " " + marks + " "
		if showEmail {
			line += t.dim().Render(padRight(ansi.Truncate(emailOf(p), emailW, symMore), emailW)) + "  "
		}
		line += d.usageCell(p, bars, now) + "  "
		if showLast {
			line += t.dim().Render(padRight(shortAgo(d.info[p.Name].last, now), lastW)) + "  "
		}
		var groups []string
		for _, g := range d.cfg.GroupsOf(p.Name) {
			groups = append(groups, t.Group(g.Name))
		}
		if len(groups) == 0 {
			line += t.dim().Render("no group")
		} else {
			line += strings.Join(groups, " ")
		}
		out = append(out, line)
	}
	return out, d.cursor[tabAccounts] + 1
}

// usageCell is the 5-hour, weekly and Fable columns: for each its bar (with
// bars), the percent used and the time until it resets.
func (d *Dashboard) usageCell(p *config.Profile, bars bool, now time.Time) string {
	t := d.theme
	cw := cellW
	if bars {
		cw += cellBarW + 1
	}
	usageW := 3*cw + 4
	r, ok := d.usage[p.Name]
	switch {
	case !ok:
		return padRight(t.accent().Render(spinner[d.frame%len(spinner)])+t.dim().Render(" "+symMore), usageW)
	case r.Status == usage.StatusAPI:
		return padRight(t.dim().Render("no limits (API)"), usageW)
	case r.Status == usage.StatusLoggedOut:
		return padRight(t.warning().Render(symFail+" not logged in"), usageW)
	case r.Status == usage.StatusExpired && r.Limits == nil:
		// Only these need `ccam login`; an expired access token is renewed
		// when the account starts.
		for _, why := range []string{"login expired", "login may be lost"} {
			if strings.HasPrefix(r.Err, why) {
				return padRight(t.warning().Render(symWarn+" "+why), usageW)
			}
		}
		return padRight(t.dim().Render(symWarn+" token expired"), usageW)
	case r.Limits == nil:
		return padRight(t.dim().Render(symWarn+" no data"), usageW)
	}
	// Values from a cache after a failed lookup are drawn faint.
	stale := r.Status != usage.StatusOK
	cell := func(w usage.Window) string {
		if !w.Known {
			return padRight(t.dim().Render("-"), cw)
		}
		pct, left := fmt.Sprintf("%3.0f%%", w.Percent), ""
		if !w.ResetsAt.IsZero() {
			left = usage.Duration(w.ResetsAt.Sub(now))
		}
		left = padRight(" "+left, 1+timeW)
		switch {
		case stale && bars:
			return t.dim().Render(plainBar(w.Percent, cellBarW) + " " + pct + left)
		case stale:
			return t.dim().Render(pct + left)
		case bars:
			return t.bar(w.Percent, cellBarW) + " " + t.percentStyle(w.Percent).Render(pct) + t.dim().Render(left)
		}
		return t.percentStyle(w.Percent).Render(pct) + t.dim().Render(left)
	}
	l := r.Limits
	return cell(l.FiveHour) + "  " + cell(l.Weekly) + "  " + cell(modelWindow(l, "fable"))
}

// modelWindow is the model-specific weekly window whose name contains
// family, such as Fable's.
func modelWindow(l *usage.Limits, family string) usage.Window {
	for _, s := range l.Scoped {
		if strings.Contains(strings.ToLower(s.Name), family) {
			return s.Window
		}
	}
	return usage.Window{}
}

func plainBar(percent float64, width int) string {
	n := int(percent/100*float64(width) + 0.5)
	n = max(0, min(n, width))
	return strings.Repeat("█", n) + strings.Repeat("░", width-n)
}

func padLeft(s string, w int) string {
	if gap := w - ansi.StringWidth(s); gap > 0 {
		return strings.Repeat(" ", gap) + s
	}
	return s
}

func (d *Dashboard) accountsKey(k tea.KeyPressMsg) tea.Cmd {
	key := k.String()
	if len(key) == 1 && key[0] >= '0' && key[0] <= '9' {
		return d.digit(key)
	}
	p := d.current()
	if p == nil {
		if key == "a" {
			d.dlg = newWizardDialog(d)
		}
		return nil
	}
	switch key {
	case "enter":
		return d.quit(Result{Action: Launch, Profile: p.Name})
	case "n":
		return d.mostRoom()
	case "o":
		d.dlg = d.launchOptions(p.Name)
	case "u":
		d.dlg = d.usageDetail(p)
	case "a":
		d.dlg = newWizardDialog(d)
	case "l":
		if p.API != nil {
			d.notice = notice{text: p.Name + " is an API key account; it has no login (its key comes from " + p.API.KeyFrom + ")", err: true}
			return nil
		}
		return d.quit(Result{Action: Login, Profile: p.Name})
	case "s":
		d.dlg = d.groupsFor(p.Name)
	case "d":
		cfg, err := d.be.SetDefault(p.Name)
		d.after(cfg, err, p.Name+" is the default account")
	case "f":
		return d.useHere(p.Name)
	case "e":
		name := p.Name
		d.dlg = newInputDialog("Rename "+name, "a-z, 0-9, - and _; its shortcut command stays the same", name, "", func(d *Dashboard, v string) (tea.Cmd, error) {
			if v == name || v == "" {
				return nil, nil
			}
			cfg, err := d.be.Rename(name, v)
			if err != nil {
				return nil, err
			}
			// The usage is the same account's; it keeps it under its new name.
			if r := d.usage[name]; r != nil {
				c := *r
				c.Name = v
				d.usage[v] = &c
			}
			d.after(cfg, nil, name+" is now "+v)
			return nil, nil
		})
	case "c":
		name, now := p.Name, p.Shortcut
		value := cmp.Or(now, d.be.SuggestShortcut(name))
		hint := "a command that starts Claude Code as " + name + "; leave it empty to remove the shortcut"
		d.dlg = newInputDialog("Shortcut for "+name, hint, value, "", func(d *Dashboard, v string) (tea.Cmd, error) {
			if v == now {
				return nil, nil
			}
			cfg, err := d.be.SetShortcut(name, v)
			if err != nil {
				return nil, err
			}
			msg := "`" + v + "` starts " + name
			if v == "" {
				msg = name + " has no shortcut now"
			}
			d.after(cfg, nil, msg)
			return nil, nil
		})
	case "x":
		name := p.Name
		body := "It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay."
		switch {
		case p.API != nil:
			body = "Its folder goes to ccam's trash; the groups' shared files stay."
		case p.Dir == "":
			if q := d.be.SameLogin(name); q != "" {
				body = "It is not logged out, since " + q + " holds the same login; its folder goes to ccam's trash and the groups' shared files stay."
			}
		}
		switch {
		case p.Dir != "" && !p.Share.Equal(config.Share{}):
			body = "It was imported, so its folder " + p.Dir + " stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it."
		case p.Dir != "":
			body = "It was imported, so its folder " + p.Dir + " stays where it is; ccam only forgets it."
		}
		d.dlg = newDestructive("Remove account "+name+"?", body, "Yes, remove "+name, func(d *Dashboard) tea.Cmd {
			return d.run("Removing "+name+"…", func() doneMsg {
				cfg, msg, err := d.be.Remove(name)
				return doneMsg{cfg: cfg, notice: msg, err: err}
			})
		})
	default:
		return d.moveCursor(k)
	}
	return nil
}

// after shows the outcome of a quick change made in place.
func (d *Dashboard) after(cfg *config.Config, err error, ok string) {
	if err != nil {
		d.notice = notice{text: err.Error(), err: true}
		return
	}
	if cfg != nil {
		d.setConfig(cfg)
	}
	d.notice = notice{text: ok}
}

// digit moves the cursor to the account with that number; digits typed
// within a second make one number while some account's number starts with
// them, so 1 then 2 goes to 12.
func (d *Dashboard) digit(key string) tea.Cmd {
	now := d.now()
	fits := func(prefix string) bool {
		for _, p := range d.cfg.Profiles {
			if strings.HasPrefix(strconv.Itoa(p.Slot), prefix) {
				return true
			}
		}
		return false
	}
	if now.Sub(d.digitAt) > time.Second || !fits(d.digits+key) {
		d.digits = ""
	}
	d.digits += key
	d.digitAt = now
	n, _ := strconv.Atoi(d.digits)
	for i, p := range d.cfg.Profiles {
		if p.Slot == n {
			d.cursor[tabAccounts] = i
			return nil
		}
	}
	if !fits(d.digits) {
		d.notice = notice{text: "no account with number " + d.digits, err: true}
		d.digits = ""
	}
	return nil
}

// mostRoom launches the account `ccam next --auto` would choose.
func (d *Dashboard) mostRoom() tea.Cmd {
	// Like next --auto: stay with accounts of the selected account's
	// organization, and skip folders logged in to the same account.
	from := d.current()
	cands := d.be.Takeover(from.Name)
	var reports []*usage.Report
	loading := false
	for _, name := range cands {
		r := d.usage[name]
		if r == nil {
			loading = true
			continue
		}
		c := *r
		c.Sessions = d.be.Sessions(name)
		reports = append(reports, &c)
	}
	pol := usage.Policy{MaxUsage: d.cfg.AutoMaxUsage, Strategy: d.cfg.AutoStrategy}
	if pol.Check() != nil {
		pol = usage.Policy{}
	}
	best := pol.Choose(reports, d.now())
	if best == nil {
		msg := "no logged-in account has room"
		if len(cands) < len(d.cfg.Profiles) {
			msg += " among those that may take over from " + from.Name + " (its organization; cross_org = true looks further)"
		}
		if loading {
			msg = "no account has room yet; usage is still loading"
		}
		d.notice = notice{text: msg, err: true}
		return nil
	}
	return d.quit(Result{Action: Launch, Profile: best.Name})
}

func (d *Dashboard) launchOptions(name string) dialog {
	labels := []string{
		"Start a new session",
		"Continue the last conversation (--continue)",
		"Choose a conversation to resume (--resume)",
		"Start with other arguments…",
		"Show the command only (dry run)",
	}
	return &choiceDialog{title: "Launch " + name, labels: labels, onPick: func(d *Dashboard, i int) tea.Cmd {
		r := Result{Action: Launch, Profile: name}
		switch i {
		case 1:
			r.Args = []string{"--continue"}
		case 2:
			r.Args = []string{"--resume"}
		case 3:
			d.dlg = newInputDialog("Arguments for claude", "as in `ccam "+name+" <arguments>`; quote an argument with spaces", "", "-p \"…\" --model opus", func(d *Dashboard, v string) (tea.Cmd, error) {
				args, err := splitArgs(v)
				if err != nil {
					return nil, err
				}
				return d.quit(Result{Action: Launch, Profile: name, Args: args}), nil
			})
			return nil
		case 4:
			r.DryRun = true
		}
		return d.quit(r)
	}}
}

func (d *Dashboard) usageDetail(p *config.Profile) dialog {
	t := d.theme
	title := "Usage of " + p.Name
	r, ok := d.usage[p.Name]
	switch {
	case !ok:
		return &infoDialog{title: title, lines: []string{"Still loading…"}}
	case r.Status == usage.StatusAPI:
		return &infoDialog{title: title, lines: []string{"An API key account has no 5-hour or weekly limit."}}
	case r.Status == usage.StatusLoggedOut:
		return &infoDialog{title: title, lines: []string{"Not logged in. Press l to log in."}}
	case r.Limits == nil:
		return &infoDialog{title: title, lines: []string{t.warning().Render(symWarn) + " " + cmp.Or(r.Err, "no usage data")}}
	}
	now := d.now()
	row := func(label string, w usage.Window) string {
		if !w.Known {
			return padRight(label, 14) + t.dim().Render("-")
		}
		line := padRight(label, 14) + t.bar(w.Percent, 20) + " " + t.percentStyle(w.Percent).Render(fmt.Sprintf("%3.0f%%", w.Percent))
		if !w.ResetsAt.IsZero() {
			line += t.dim().Render("  resets in " + usage.Duration(w.ResetsAt.Sub(now)))
		}
		return line
	}
	lines := []string{row("5-hour", r.Limits.FiveHour), row("Weekly", r.Limits.Weekly)}
	for _, s := range r.Limits.Scoped {
		lines = append(lines, row(s.Name+" weekly", s.Window))
	}
	src := map[string]string{"live": "live, from the usage endpoint", "claude": "live, from Claude Code", "cache": "saved from the usage endpoint", "statusline": "saved by the status line",
		"limit": "the usage limit Claude Code reported"}[r.Source]
	info := "source: " + cmp.Or(src, r.Source)
	if !r.At.IsZero() {
		info += symDot + shortAgo(r.At, now) + " ago"
	}
	if r.HoldsReset {
		lines = append(lines, "", t.accent().Render(symReset)+" Holds an unused limit reset"+t.dim().Render(symDot+"/limit-reset in Claude Code uses it"))
	}
	lines = append(lines, "", t.dim().Render(info))
	if plan := p.Login.PlanName(); plan != "" {
		if p.Login.Pool() != config.Personal && p.Login.OrgName != "" {
			plan += symDot + p.Login.OrgName
		}
		lines = append(lines, t.dim().Render("plan: "+plan))
	}
	if r.Err != "" {
		lines = append(lines, t.warning().Render(symWarn+" "+r.Err))
	}
	return &infoDialog{title: title, lines: lines}
}

// useHere makes the account the one a launch without a name uses in the
// current folder.
func (d *Dashboard) useHere(name string) tea.Cmd {
	set := func(d *Dashboard) tea.Cmd {
		cfg, err := d.be.SetFolder(d.dir, name)
		d.after(cfg, err, name+" is used in "+d.dir)
		return nil
	}
	for _, f := range d.cfg.Folders {
		if config.SamePath(f.Path, d.dir) && f.Account != name {
			d.dlg = newConfirm("Use "+name+" in "+d.dir+"?", "This folder uses "+f.Account+" now.", "Yes, use "+name, set)
			return nil
		}
	}
	return set(d)
}

// groupsFor is the checklist of the groups account is in.
func (d *Dashboard) groupsFor(account string) dialog {
	var in []string
	for _, g := range d.cfg.GroupsOf(account) {
		in = append(in, g.Name)
	}
	return &groupsDialog{account: account, orig: in, list: groupChecklist(GroupInfos(d.cfg), in)}
}

type groupsDialog struct {
	account string
	orig    []string
	list    checklist
	err     string
}

func (g *groupsDialog) update(d *Dashboard, msg tea.Msg) tea.Cmd {
	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return nil
	}
	if g.list.key(k.String()) {
		g.err = ""
		return nil
	}
	switch k.String() {
	case "esc", "q":
		d.dlg = nil
	case "enter":
		next := g.list.items()
		if slices.Equal(next, g.orig) {
			d.dlg = nil
			return nil
		}
		if g.err = overlap(GroupInfos(d.cfg), next); g.err != "" {
			return nil
		}
		d.dlg = nil
		return d.planChange("Groups of "+g.account, SetGroups(g.account, next), g.account+" is in "+cmp.Or(strings.Join(next, ", "), "no group"))
	}
	return nil
}

func (g *groupsDialog) view(d *Dashboard, w int) []string {
	t := d.theme
	out := dialogTop(t, w, "Groups of "+g.account)
	if len(g.list.keys) == 0 {
		return append(out, "  No groups yet. Create one on the Groups tab (n).")
	}
	members := map[string]string{}
	for _, gr := range d.cfg.Groups {
		members[gr.Name] = strings.Join(gr.Members, ", ")
	}
	var b strings.Builder
	g.list.styled(&b, t, func(name string) string {
		note := members[name]
		if slices.Contains(g.orig, name) != g.list.selected[name] {
			note += " " + t.warning().Render("(changed)")
		}
		return note
	})
	for _, l := range strings.Split(strings.TrimRight(b.String(), "\n"), "\n") {
		out = append(out, " "+l)
	}
	if g.err != "" {
		out = append(out, "", "  "+t.error().Render(symFail+" "+g.err))
	}
	return out
}

func (g *groupsDialog) keys() string {
	return "space to join or leave" + symDot + "enter to review the plan" + symDot + "esc to cancel"
}

// SetGroups makes account a member of exactly groups.
func SetGroups(account string, groups []string) func(*config.Config) error {
	return func(c *config.Config) error {
		for _, g := range c.GroupsOf(account) {
			if !slices.Contains(groups, g.Name) {
				if err := c.LeaveGroup(account, g.Name); err != nil {
					return err
				}
			}
		}
		for _, g := range groups {
			if _, err := c.JoinGroup(account, g, false); err != nil {
				return err
			}
		}
		return nil
	}
}

// planChange shows what edit moves, asks about conflicts and the plan, and
// applies it in the background.
func (d *Dashboard) planChange(title string, edit func(*config.Config) error, done string) tea.Cmd {
	ch, err := d.be.PlanGroups(edit)
	if err != nil {
		d.notice = notice{text: err.Error(), err: true}
		return nil
	}
	apply := func(d *Dashboard) tea.Cmd {
		return d.run("Applying…", func() doneMsg {
			cfg, err := d.be.ApplyGroups(edit, ch)
			return doneMsg{cfg: cfg, notice: done, err: err}
		})
	}
	if ch.Empty() {
		return apply(d)
	}
	return d.askConflicts(ch.Conflicts(), 0, func(d *Dashboard) tea.Cmd {
		d.dlg = &planDialog{title: title, lines: planLines(d.theme, ch), onYes: apply}
		return nil
	})
}

// askConflicts asks about each conflict in turn, then runs then.
func (d *Dashboard) askConflicts(cs []group.Conflict, i int, then func(*Dashboard) tea.Cmd) tea.Cmd {
	if i == len(cs) {
		return then(d)
	}
	c := cs[i]
	labels := make([]string, len(c.Candidates))
	cur := 0
	for j, r := range c.Candidates {
		labels[j] = r.String()
		if r == c.Choice {
			cur = j
		}
	}
	grp := c.Group
	if c.Seed != "" {
		grp = c.Seed + " (moves into the group first)"
	}
	d.dlg = &choiceDialog{
		title:  fmt.Sprintf("%s: %s differs from the group's copy (%d of %d)", c.Account, c.Entry.Name, i+1, len(cs)),
		body:   "mine:  " + c.Mine + "\ngroup: " + grp,
		labels: labels, cursor: cur,
		onPick: func(d *Dashboard, j int) tea.Cmd {
			c.Choice = c.Candidates[j]
			return d.askConflicts(cs, i+1, then)
		},
	}
	return nil
}

// planLines is a change as ● account / ⎿ step lines.
func planLines(t Theme, ch *group.Change) []string {
	var buf bytes.Buffer
	ch.Describe(&buf)
	var out []string
	for _, l := range strings.Split(strings.TrimRight(buf.String(), "\n"), "\n") {
		switch {
		case !strings.HasPrefix(l, " "):
			name, rest, _ := strings.Cut(l, ": ")
			out = append(out, t.accent().Render(symItem)+" "+t.bold().Render(name)+t.dim().Render("  "+rest))
		case strings.HasPrefix(l, "  Note: "):
			out = append(out, "  "+t.warning().Render(symWarn)+" "+strings.TrimPrefix(l, "  Note: "))
		default:
			out = append(out, "  "+t.dim().Render(symDetail)+"  "+strings.TrimPrefix(l, "  "))
		}
	}
	return out
}

// splitArgs splits a command line at spaces, keeping "…" and '…' together
// (without the quotes), as a shell would for plain arguments.
func splitArgs(s string) ([]string, error) {
	var out []string
	var cur strings.Builder
	quote, in := rune(0), false
	for _, r := range s {
		switch {
		case quote != 0 && r == quote:
			quote = 0
		case quote != 0:
			cur.WriteRune(r)
		case r == '"' || r == '\'':
			quote, in = r, true
		case r == ' ' || r == '\t':
			if in {
				out = append(out, cur.String())
				cur.Reset()
				in = false
			}
		default:
			cur.WriteRune(r)
			in = true
		}
	}
	if quote != 0 {
		return nil, fmt.Errorf("a %c is not closed", quote)
	}
	if in {
		out = append(out, cur.String())
	}
	return out, nil
}
