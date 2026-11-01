package launcher

import (
	"reflect"
	"strings"
	"testing"
)

func TestCarryArgs(t *testing.T) {
	split := func(s string) []string { return strings.Fields(s) }
	for _, c := range []struct{ in, out, dropped string }{
		// Resuming restores the model and the prompt is in the conversation.
		{"fix it --model opus --permission-mode plan --add-dir /tmp", "--permission-mode plan --add-dir /tmp", ""},
		{"--model opus --verbose fix it", "--verbose", ""},
		// Variadic flags take values up to the next flag, as Claude Code reads
		// them (so a prompt right after one was one of its values); = takes no more.
		{"--add-dir /a /b --verbose --allowedTools Bash Edit", "--add-dir /a /b --verbose --allowedTools Bash Edit", ""},
		{"--add-dir /a fix it", "--add-dir /a fix it", ""},
		{"--add-dir=/a prompt --mcp-config=x.json", "--add-dir=/a --mcp-config=x.json", ""},
		// The first value is taken whatever it looks like, as the one of a flag with one value.
		{"--add-dir -odd --verbose", "--add-dir -odd --verbose", ""},
		// A value that starts with - is still the value of a flag that takes one.
		{"--append-system-prompt -terse- --effort high", "--append-system-prompt -terse- --effort high", ""},
		// "--" ends the flags: what follows is a prompt.
		{"--ide -- --add-dir /x", "--ide", ""},
		// What picks the conversation goes quietly; what cannot be done again does not.
		// The agent goes on: an interactive --resume <file> --fork-session does
		// not bring it back (Claude Code 2.1.314).
		{"-c --resume abc --fork-session --session-id 1234 -n name --agent rev --from-pr 12 --teleport", "--agent rev", ""},
		{"--agent=rev --model opus", "--agent=rev", ""},
		{"--plugin-url https://example.com/p.zip --plugin-dir ./p", "--plugin-url https://example.com/p.zip --plugin-dir ./p", ""},
		{"--resume --verbose", "--verbose", ""},
		{"-w feature --tmux --debug api,hooks --debug-file /tmp/d.log --bg --remote-control phone", "", "-w --tmux --debug --debug-file --bg --remote-control"},
		// An unknown flag goes with a word; a value after it goes as a prompt would.
		{"--shiny value --chrome", "--chrome", "--shiny"},
		{"--shiny=1 --shiny=2", "", "--shiny"},
		// Flags a hop drops take their values as Claude Code reads them (arity).
		{"--output-format -x --file a b --verbose", "--verbose", "--output-format --file"},
		{"--dangerously-skip-permissions --settings {\"a\":1} --plugin-dir p1 --plugin-dir p2", "--dangerously-skip-permissions --settings {\"a\":1} --plugin-dir p1 --plugin-dir p2", ""},
	} {
		out, dropped := CarryArgs(split(c.in))
		if strings.Join(out, " ") != c.out || strings.Join(dropped, " ") != c.dropped {
			t.Errorf("%s:\n  carried %q, want %q\n  dropped %q, want %q", c.in, out, c.out, dropped, c.dropped)
		}
	}
	// A second hop carries what the first one passed.
	first, _ := WithState(must(CarryArgs(split("--model opus --add-dir /tmp --permission-mode plan --effort low go on"))), "acceptEdits", "high")
	again, _ := WithState(must(CarryArgs(append([]string{"--resume", "s1"}, first...))), "acceptEdits", "high")
	if !reflect.DeepEqual(first, again) {
		t.Fatalf("first %q, then %q", first, again)
	}
}

func must(out, _ []string) []string { return out }

func TestWithState(t *testing.T) {
	split := func(s string) []string { return strings.Fields(s) }
	for _, c := range []struct{ in, mode, effort, out, dropped string }{
		// Nothing known: the args as they were.
		{"--permission-mode plan --effort high --dangerously-skip-permissions", "", "", "--permission-mode plan --effort high --dangerously-skip-permissions", ""},
		// The conversation's mode and the session's effort win.
		{"--permission-mode plan --add-dir /tmp --effort low", "acceptEdits", "high", "--add-dir /tmp --permission-mode acceptEdits --effort high", ""},
		{"--permission-mode=plan --effort=low", "default", "max", "--permission-mode default --effort max", ""},
		// Bypass stays allowed when it was, without forcing it on.
		{"--dangerously-skip-permissions --ide", "default", "", "--ide --permission-mode default --allow-dangerously-skip-permissions", ""},
		{"--allow-dangerously-skip-permissions", "plan", "", "--permission-mode plan --allow-dangerously-skip-permissions", ""},
		{"", "bypassPermissions", "", "--permission-mode bypassPermissions --allow-dangerously-skip-permissions", ""},
		// A session that started in bypass may turn it on again.
		{"--permission-mode bypassPermissions", "default", "", "--permission-mode default --allow-dangerously-skip-permissions", ""},
		{"--permission-mode=bypassPermissions --ide", "plan", "", "--ide --permission-mode plan --allow-dangerously-skip-permissions", ""},
		// Claude Code stops on a mode it does not know, so it is not passed.
		{"--permission-mode plan", "turbo", "", "--permission-mode plan", "--permission-mode turbo"},
	} {
		in := split(c.in)
		out, dropped := WithState(in, c.mode, c.effort)
		if strings.Join(out, " ") != c.out || strings.Join(dropped, " ") != c.dropped {
			t.Errorf("%q %q %q:\n  %q, want %q\n  dropped %q, want %q", c.in, c.mode, c.effort, out, c.out, dropped, c.dropped)
		}
		if strings.Join(in, " ") != c.in {
			t.Errorf("%q: the args given changed to %q", c.in, in)
		}
	}
}

// The agent a session ran as by another way than its arguments (the agent
// setting, or the conversation's own) goes on when they name none. An
// --agent they give stays: Claude Code ran the session as it, also when
// the conversation's rows still name the agent it ran as before.
func TestWithAgent(t *testing.T) {
	split := func(s string) []string { return strings.Fields(s) }
	for _, c := range []struct{ in, agent, out string }{
		{"--add-dir /tmp", "", "--add-dir /tmp"},
		{"--agent rev --add-dir /tmp", "", "--agent rev --add-dir /tmp"},
		{"--add-dir /tmp", "rev", "--add-dir /tmp --agent rev"},
		{"--agent rev --add-dir /tmp", "writer", "--agent rev --add-dir /tmp"},
		{"--agent=a --verbose --agent b", "c", "--agent=a --verbose --agent b"},
		// --agents defines agents; it names none to run as.
		{`--agents {"x":{}} --verbose`, "x", `--agents {"x":{}} --verbose --agent x`},
		// After "--" every argument is the prompt's, an --agent too.
		{"--verbose -- -x --agent b", "rev", "--verbose --agent rev -- -x --agent b"},
	} {
		in := split(c.in)
		if out := WithAgent(in, c.agent); strings.Join(out, " ") != c.out {
			t.Errorf("%q %q: %q, want %q", c.in, c.agent, out, c.out)
		}
		if strings.Join(in, " ") != c.in {
			t.Errorf("%q: the args given changed to %q", c.in, in)
		}
	}
}

// The agent claude args run as: the value of the last --agent before "--".
func TestAgentArg(t *testing.T) {
	for in, want := range map[string]string{
		"":                             "",
		"--resume s1 --verbose":        "",
		"--resume s1 --agent reviewer": "reviewer",
		"--agent=a --add-dir /tmp":     "a",
		"--agent a --agent=b":          "b",
		`--agents {"x":{}} go`:         "",
		"--agent":                      "",
		"go -- --agent x":              "",
	} {
		if got := AgentArg(strings.Fields(in)); got != want {
			t.Errorf("%q: %q, want %q", in, got, want)
		}
	}
}

// The prompt among claude args, as Claude Code 2.1.310 reads them: the
// first argument that is no option or option's value, or the one after
// "--"; the rest keep their order.
func TestTakePrompt(t *testing.T) {
	split := func(s string) []string { return strings.Fields(s) }
	for _, c := range []struct{ in, prompt, rest string }{
		{"", "", ""},
		{"--model sonnet fix", "fix", "--model sonnet"},
		{"fix --model sonnet", "fix", "--model sonnet"},
		{"--add-dir /a /b", "", "--add-dir /a /b"}, // a variadic flag takes both, as claude does
		{"--add-dir=/a fix", "fix", "--add-dir=/a"},
		{"-p fix", "fix", "-p"},
		{"--debug api fix", "fix", "--debug api"},
		{"--resume", "", "--resume"},
		{"--model sonnet -- -x y", "-x", "--model sonnet"}, // after "--", a prompt; what follows it goes
		{"--", "", ""},
		{"fix more", "fix", "more"}, // Claude Code ignores a second prompt; it stays where it was
		{"fix -- -x", "fix", ""},
	} {
		in := split(c.in)
		prompt, rest := TakePrompt(in)
		if prompt != c.prompt || strings.Join(rest, " ") != c.rest {
			t.Errorf("%q: prompt %q, rest %q; want %q, %q", c.in, prompt, rest, c.prompt, c.rest)
		}
		if strings.Join(in, " ") != c.in {
			t.Errorf("%q: the args given changed to %q", c.in, in)
		}
	}
	// A prompt with spaces and an empty one are prompts all the same.
	if prompt, rest := TakePrompt([]string{"--verbose", "run the tests"}); prompt != "run the tests" || !reflect.DeepEqual(rest, []string{"--verbose"}) {
		t.Fatalf("%q %q", prompt, rest)
	}
	if prompt, rest := TakePrompt([]string{"", "--verbose"}); prompt != "" || !reflect.DeepEqual(rest, []string{"--verbose"}) {
		t.Fatalf("%q %q", prompt, rest)
	}
}
