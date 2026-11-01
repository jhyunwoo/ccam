package launcher

import (
	"slices"
	"strings"
)

// How many values follow a claude flag (Claude Code 2.1.308, commander).
const (
	noValue  = iota
	oneValue // always the next argument, even one that starts with -
	optValue // the next argument unless it is a flag
	values   // the next argument, as for oneValue, then more up to the next flag
)

// carried says what a hop does with a claude flag: whether it passes it on
// and, when not, whether it says so; and how many values follow it.
func carried(flag string) (keep, quiet bool, n int) {
	switch flag {
	// Resuming does not restore these, so they go on to the next session;
	// --agent among them, which an interactive --resume <file>
	// --fork-session (the hop's resume without shared history) does not
	// bring back (Claude Code 2.1.314).
	case "--add-dir", "--mcp-config", "--allowedTools", "--allowed-tools", "--disallowedTools", "--disallowed-tools", "--tools", "--betas":
		return true, false, values
	case "--permission-mode", "--effort", "--agent", "--agents", "--append-system-prompt", "--append-system-prompt-file", "--system-prompt",
		"--system-prompt-file", "--system-prompt-snapshot", "--settings", "--setting-sources", "--plugin-dir", "--plugin-url",
		"--fallback-model", "--autocompact":
		return true, false, oneValue
	case "--dangerously-skip-permissions", "--allow-dangerously-skip-permissions", "--strict-mcp-config", "--ide", "--chrome",
		"--no-chrome", "--verbose", "--bare", "--safe-mode", "--restricted", "--disable-slash-commands", "--brief",
		"--ax-screen-reader", "--exclude-dynamic-system-prompt-sections":
		return true, false, noValue
	// They pick the conversation, which the hop does itself, or resuming
	// brings them back (model, name).
	case "-c", "--continue", "--fork-session":
		return false, true, noValue
	case "-r", "--resume", "--from-pr", "--teleport":
		return false, true, optValue
	case "--session-id", "--model", "-n", "--name":
		return false, true, oneValue
	// What a session cannot do again on its way to another account.
	case "-w", "--worktree", "-d", "--debug", "--remote-control", "--rc":
		return false, false, optValue
	case "--debug-file":
		return false, false, oneValue
	}
	return false, false, noValue // --tmux, --bg and flags ccam does not know
}

func isFlag(a string) bool { return len(a) > 1 && a[0] == '-' }

// CarryArgs picks the claude args a session started with that a hop passes
// on to the conversation's next session: those resuming does not restore
// (directories, MCP servers, tools, the agent, permission and system prompt
// flags, settings, plugins …). It drops the prompt, "--" and everything
// after it, what picks or restores the conversation (--resume, --model,
// --name …) and the rest; dropped names the flags of the rest, which a
// person may miss (--worktree, --debug, --remote-control, unknown ones).
func CarryArgs(args []string) (out, dropped []string) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			break
		}
		if !isFlag(a) {
			continue // the prompt
		}
		name, _, inline := strings.Cut(a, "=")
		keep, quiet, _ := carried(name)
		take := 0
		if !inline { // --flag=value: no more values follow, also for variadic ones
			take = takes(args, i, arity(name))
		}
		switch {
		case keep:
			out = append(out, args[i:i+1+take]...)
		case !quiet && !slices.Contains(dropped, name):
			dropped = append(dropped, name)
		}
		i += take
	}
	return out, dropped
}

// TakePrompt splits claude args into their prompt, the first argument that
// is neither an option nor an option's value (after "--", the argument
// after it), and the rest without it. Claude Code takes one prompt and
// ignores more (2.1.310), so a caller that sends a prompt of its own folds
// the person's into it. "--" and what follows the prompt after it are
// dropped, as CarryArgs drops them.
func TakePrompt(args []string) (prompt string, rest []string) {
	found := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--":
			if !found && i+1 < len(args) {
				prompt = args[i+1]
			}
			return prompt, rest
		case isFlag(a):
			take := 0
			if name, _, inline := strings.Cut(a, "="); !inline {
				take = takes(args, i, arity(name))
			}
			rest = append(rest, args[i:i+1+take]...)
			i += take
		case found:
			rest = append(rest, a)
		default:
			prompt, found = a, true
		}
	}
	return prompt, rest
}

// knownMode reports whether claude --permission-mode takes mode (2.1.308;
// it also takes default, and stops on anything else).
func knownMode(mode string) bool {
	switch mode {
	case "acceptEdits", "auto", "bypassPermissions", "default", "dontAsk", "plan", "manual":
		return true
	}
	return false
}

// WithState adds to args (from CarryArgs) the session's state that resuming
// does not bring back: the permission mode the conversation was last in,
// in place of --permission-mode and --dangerously-skip-permissions (with
// --allow-dangerously-skip-permissions when bypass was allowed or the
// session started in it, so it can be turned on again), and the effort
// level, in place of --effort. An
// empty mode or effort leaves args as they are; a mode Claude Code would
// refuse is not passed, and dropped says so.
func WithState(args []string, mode, effort string) (out, dropped []string) {
	out = slices.Clone(args)
	switch {
	case mode == "":
	case !knownMode(mode):
		dropped = append(dropped, "--permission-mode "+mode)
	default:
		bypass := mode == "bypassPermissions" || startedIn(out, "bypassPermissions") ||
			slices.Contains(out, "--dangerously-skip-permissions") || slices.Contains(out, "--allow-dangerously-skip-permissions")
		out = without(out, "--permission-mode", true)
		out = without(out, "--dangerously-skip-permissions", false)
		out = without(out, "--allow-dangerously-skip-permissions", false)
		out = append(out, "--permission-mode", mode)
		if bypass {
			out = append(out, "--allow-dangerously-skip-permissions")
		}
	}
	if effort != "" {
		out = append(without(out, "--effort", true), "--effort", effort)
	}
	return out, dropped
}

// WithAgent adds --agent agent to claude args that name no agent: agent is
// the one a session ran the conversation as without an --agent of its own
// (the agent setting, or the conversation's agent, which resuming it by id
// brings back), and an interactive --resume <file> --fork-session does not
// bring it back (Claude Code 2.1.314). An --agent the args give stays:
// Claude Code ran the session as that one, while the conversation's
// agent-setting rows may still name the agent it ran as before. It goes
// before a "--", after which everything is the prompt's. An empty agent
// leaves args as they are.
func WithAgent(args []string, agent string) []string {
	if agent == "" || AgentArg(args) != "" {
		return args
	}
	end := slices.Index(args, "--")
	if end < 0 {
		end = len(args)
	}
	return slices.Insert(slices.Clone(args), end, "--agent", agent)
}

// AgentArg is the agent claude args run as: the value of the last --agent
// before "--", as Claude Code takes it; "" when they name none.
func AgentArg(args []string) string {
	agent := ""
	for i, a := range args {
		if a == "--" {
			break
		}
		if v, ok := strings.CutPrefix(a, "--agent="); ok {
			agent = v
		} else if a == "--agent" && i+1 < len(args) {
			agent = args[i+1]
		}
	}
	return agent
}

// startedIn reports whether args give --permission-mode mode.
func startedIn(args []string, mode string) bool {
	for i, a := range args {
		if a == "--permission-mode="+mode || a == "--permission-mode" && i+1 < len(args) && args[i+1] == mode {
			return true
		}
	}
	return false
}

// without drops flag from args, and the value after it when it takes one.
func without(args []string, flag string, value bool) []string {
	out := args[:0:0]
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == flag:
			if value {
				i++
			}
		case strings.HasPrefix(args[i], flag+"="):
		default:
			out = append(out, args[i])
		}
	}
	return out
}
