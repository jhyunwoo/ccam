package handoff

import (
	"os"
	"strings"
)

// ContinuePrompt is the first message of a conversation that the relay
// carried on as another account after a usage limit (never after
// /ccam:switch, nor for `ccam next --continue` by hand).
const ContinuePrompt = "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete."

// SubagentPrompt is ContinuePrompt after a subagent's limit: the session
// ended while subagents, workflows or background commands ran, and the
// results they were to bring never came.
const SubagentPrompt = "A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete."

// WindowLabel names a window of Claude Code's limit row for messages:
// 5-hour, weekly, Opus, Sonnet or Fable; "" for another.
func WindowLabel(window string) string {
	switch window {
	case "five_hour":
		return "5-hour"
	case "seven_day":
		return "weekly"
	case "seven_day_opus":
		return "Opus"
	case "seven_day_sonnet":
		return "Sonnet"
	case "seven_day_overage_included":
		return "Fable"
	}
	return ""
}

// HitLimit begins ccam's messages about a usage limit: "ccam: p1 hit its
// 5-hour limit (resets 3:45pm)"; without a reset, the part in parentheses
// is left out.
func HitLimit(account, label, reset string) string {
	msg := "ccam: " + account + " hit its " + label + " limit"
	if reset != "" {
		msg += " (resets " + reset + ")"
	}
	return msg
}

// Notification is msg as a desktop notification in the escape sequence the
// terminal understands: OSC 99 for kitty, OSC 777 for Ghostty, urxvt and
// Warp, OSC 9 for the rest (iTerm2, WezTerm, Windows Terminal …).
func Notification(msg string) string {
	msg = strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return ' '
		}
		return r
	}, msg)
	tp := os.Getenv("TERM_PROGRAM")
	switch {
	case os.Getenv("KITTY_WINDOW_ID") != "":
		return "\x1b]99;;" + msg + "\x1b\\"
	case tp == "ghostty" || tp == "WarpTerminal" || strings.HasPrefix(os.Getenv("TERM"), "rxvt"):
		return "\x1b]777;notify;ccam;" + msg + "\x07"
	}
	return "\x1b]9;" + msg + "\x07"
}
