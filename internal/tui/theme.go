package tui

import (
	"hash/fnv"
	"image/color"
	"slices"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/jhyunwoo/ccam/internal/config"
)

// Theme holds the colors of the screens by role. The roles and values follow
// Claude Code's terminal UI, so ccam looks at home next to it; Claude's own
// orange and its marks are left out on purpose (ccam is not an official
// product).
type Theme struct {
	Name     string
	Accent   color.Color // focus, the cursor, the current tab, dialog titles
	Success  color.Color
	Warning  color.Color
	Error    color.Color
	Inactive color.Color // rows that cannot be chosen
	Subtle   color.Color // rules and borders that stay in the back
	BarFill  color.Color
	BarEmpty color.Color
	Inverse  color.Color // text on an Accent background
	groups   []color.Color
}

func hex(s string) color.Color { return lipgloss.Color(s) }

var groupHex = []string{"#dc2626", "#6a9bcc", "#16a34a", "#ca8a04", "#827dbd", "#c46686", "#0891b2"}

// ThemeFor returns the named theme; "auto" and unknown names follow dark,
// which says whether the terminal background is dark.
func ThemeFor(name string, dark bool) Theme {
	if name == "" || name == "auto" || !validTheme(name) {
		name = "light"
		if dark {
			name = "dark"
		}
	}
	t := Theme{Name: name}
	for _, h := range groupHex {
		t.groups = append(t.groups, hex(h))
	}
	switch name {
	case "dark":
		t.Accent, t.Success, t.Warning, t.Error = hex("#b1b9f9"), hex("#4eba65"), hex("#ffc107"), hex("#ff6b80")
		t.Inactive, t.Subtle, t.BarFill, t.BarEmpty, t.Inverse = hex("#999999"), hex("#505050"), hex("#b1b9f9"), hex("#505370"), hex("#000000")
	case "light":
		t.Accent, t.Success, t.Warning, t.Error = hex("#5769f7"), hex("#2c7a39"), hex("#966c1e"), hex("#ab2b3f")
		t.Inactive, t.Subtle, t.BarFill, t.BarEmpty, t.Inverse = hex("#666666"), hex("#afafaf"), hex("#5769f7"), hex("#272f6f"), hex("#ffffff")
	case "dark-daltonized":
		t.Accent, t.Success, t.Warning, t.Error = hex("#99ccff"), hex("#3399ff"), hex("#ffcc00"), hex("#ff6666")
		t.Inactive, t.Subtle, t.BarFill, t.BarEmpty, t.Inverse = hex("#999999"), hex("#505050"), hex("#99ccff"), hex("#455c73"), hex("#000000")
	case "light-daltonized":
		t.Accent, t.Success, t.Warning, t.Error = hex("#3366ff"), hex("#006699"), hex("#ff9900"), hex("#cc0000")
		t.Inactive, t.Subtle, t.BarFill, t.BarEmpty, t.Inverse = hex("#666666"), hex("#afafaf"), hex("#3366ff"), hex("#172e72"), hex("#ffffff")
	case "dark-ansi", "light-ansi":
		// The 16 colors of the terminal's own palette.
		bright := name == "dark-ansi"
		pick := func(b, n int) color.Color {
			if bright {
				return lipgloss.Color(strconv.Itoa(b))
			}
			return lipgloss.Color(strconv.Itoa(n))
		}
		t.Accent, t.Success, t.Warning, t.Error = pick(12, 4), pick(10, 2), pick(11, 3), pick(9, 1)
		t.Inactive, t.Subtle, t.BarFill, t.BarEmpty, t.Inverse = pick(7, 8), pick(7, 8), pick(11, 3), pick(7, 0), pick(0, 15)
		t.groups = []color.Color{pick(9, 1), pick(12, 4), pick(10, 2), pick(11, 3), pick(13, 5), pick(14, 6)}
	}
	return t
}

func validTheme(name string) bool { return slices.Contains(config.Themes, name) }

func (t Theme) fg(c color.Color) lipgloss.Style { return lipgloss.NewStyle().Foreground(c) }

func (t Theme) accent() lipgloss.Style  { return t.fg(t.Accent) }
func (t Theme) success() lipgloss.Style { return t.fg(t.Success) }
func (t Theme) warning() lipgloss.Style { return t.fg(t.Warning) }
func (t Theme) error() lipgloss.Style   { return t.fg(t.Error) }
func (t Theme) inactive() lipgloss.Style {
	return t.fg(t.Inactive)
}
func (t Theme) subtle() lipgloss.Style { return t.fg(t.Subtle) }

// dim is secondary text: faint, so it works on any background.
func (t Theme) dim() lipgloss.Style  { return lipgloss.NewStyle().Faint(true) }
func (t Theme) bold() lipgloss.Style { return lipgloss.NewStyle().Bold(true) }

// badge is the current tab: inverse text on the accent color.
func (t Theme) badge() lipgloss.Style {
	return lipgloss.NewStyle().Background(t.Accent).Foreground(t.Inverse).Padding(0, 1)
}

// Group colors a group name; the color comes from the name, so a group
// keeps it everywhere and across runs.
func (t Theme) Group(name string) string {
	h := fnv.New32a()
	h.Write([]byte(name))
	return t.fg(t.groups[int(h.Sum32()%uint32(len(t.groups)))]).Render(name)
}

// Symbols used on every screen.
const (
	symCursor  = "❯"
	symOK      = "✔"
	symFail    = "✘"
	symWarn    = "⚠"
	symInfo    = "ℹ"
	symHere    = "▸"
	symRunning = "•"
	symReset   = "↻"
	symDot     = " · "
	symDetail  = "⎿"
	symItem    = "●"
	symMore    = "…"
)

var spinner = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// bar draws percent (0-100) in width cells with eighth blocks.
func (t Theme) bar(percent float64, width int) string {
	eighths := int(percent/100*float64(width*8) + 0.5)
	eighths = max(0, min(eighths, width*8))
	full, part := eighths/8, eighths%8
	var b strings.Builder
	b.WriteString(strings.Repeat("█", full))
	if part > 0 {
		b.WriteString([]string{"", "▏", "▎", "▍", "▌", "▋", "▊", "▉"}[part])
	}
	used := full
	if part > 0 {
		used++
	}
	return t.fg(t.BarFill).Render(b.String()) + t.fg(t.BarEmpty).Render(strings.Repeat("░", width-used))
}

// percentStyle colors a usage figure: warning from 80%, error at 100%.
func (t Theme) percentStyle(p float64) lipgloss.Style {
	switch {
	case p >= 100:
		return t.error()
	case p >= 80:
		return t.warning()
	}
	return lipgloss.NewStyle()
}
