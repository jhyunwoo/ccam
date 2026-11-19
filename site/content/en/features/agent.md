---
title: Subagents on other accounts
description: ccam agent runs part of a Claude Code task as another of your accounts and returns Claude's reply, so one long task uses the usage of several accounts.
updated: 2026-11-19
order: 17
translation: source
---

Claude Code's own subagents run as the account of the session that starts them, so they spend that account's usage. `ccam agent "<task>"` runs a task with `claude -p` in the same folder as another of your accounts and prints Claude's last reply. Called from a Claude Code session, it works like a subagent that spends another account's usage, so one long task can go on after the session's own account would have run out.

**Team and Enterprise seats.** `ccam agent` chooses the account by itself. Within one Claude Team or Enterprise organization it runs the task on another seat, which may violate Anthropic's terms. Read the [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) and check with your organization and Anthropic before you rely on it.

## Usage

Claude Code runs it through its Bash tool. With the [ccam plugin](/docs/install#claude-code-plugin) installed, ask for it in plain words, for example "use my other accounts for the parts of this task", or put that in your CLAUDE.md. The plugin's skill tells Claude Code to start each agent in the background, to write a task that stands alone, and to give agents that edit files their own part or git worktree.

```sh
ccam agent "find why TestParse fails and fix it; reply with what you changed"
ccam agent --model opus -- "-v breaks the parser: find out why"
ccam agent --json < task.md            # the task from standard input
ccam agent --resume 7f3c9a2e-… "now add a test for it"
ccam agent --wait=1h "…"               # nobody has room: wait for a reset
ccam --dry-run agent "…"               # the account and command it would run
```

```console
$ ccam agent "summarize what internal/usage does"
ccam: agent runs as work (12% of its 5 hours used); session 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f
internal/usage looks up the 5-hour and weekly usage of every account …
ccam: agent finished as work; follow up with: ccam agent --resume 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f "…"
```

Claude's last reply goes to standard output; the lines that start with `ccam:` go to standard error. `--json` prints `{account, session_id, result, is_error, accounts}` instead, where `accounts` lists every account the task ran as.

Exit codes: claude's own (`0` when the agent finished), `75` (`EX_TEMPFAIL`) when no account other than the caller has room, `1` for other errors such as an agent trying to start another.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--model <model>` | The agent's model; that model's weekly window counts in the choice too | the account's default model |
| `--resume <id>` | Continue the conversation of an earlier agent with a follow-up task (the session id on its first and last line) | off |
| `--json` | Print the account, session, reply and the accounts it ran as | off |
| `--max-usage <1-100>` | Prefer accounts below this usage | `auto_max_usage`, otherwise 100 |
| `--strategy headroom\|expiring\|balance` | The [choice rule](/features/next#auto) | `auto_strategy`, otherwise `headroom` |
| `--cross-org` | Also run as accounts of other organizations | `cross_org`, otherwise off |
| `--wait[=<duration>]` | When no account has room, wait for the first one to reset, at most this long | off; `5h` when given alone |
| `--dry-run` | Print the command for the account it would choose, without running it | off |

## How it works

1. **The caller is left out.** The account `ccam agent` is called from (the one of the session ccam started, else the one this folder starts) never runs the agent, nor does another folder logged in to the same Claude account, which shares its limits. The agent stays in the caller's [organization](/features/next#auto) unless `--cross-org` says otherwise.
2. **It picks like `ccam pick`.** Usage is looked up the way [`ccam pick`](/features/pick) does, and agents started within 2 minutes of each other get different accounts while enough have room.
3. **It runs with the session's permissions, never more.** Nobody is there to allow a tool in `claude -p`, so the agent gets the permission mode and effort of the session that called it. Run from a terminal, Claude Code's default mode and your settings' permission rules apply. `ccam agent` has no option that grants more.
4. **At a usage limit, the conversation moves on.** When the agent's account reaches its limit, ccam records the limit as the [plugin hook](/features/usage#the-limit-record) would and continues the same conversation as the account with the most room: Claude Code copies it there (`--resume <file> --fork-session`), with the first message "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." Each account is tried once.
5. **Every agent conversation is recorded** in `~/.local/state/ccam/agent/<session id>.json` (30 days). `--resume <id>` continues it in its folder, as the account that ran it when that one still has room, otherwise as another. [`ccam next --continue`](/features/next) leaves agent conversations out, so it keeps continuing yours.

## Safety and edge cases

- An agent cannot start agents of its own: ccam marks it with `CCAM_AGENT=1`, and `ccam agent` refuses to run there.
- Stopping `ccam agent` (Ctrl+C, a signal, or Claude Code interrupting the command) stops the agent too.
- The agent does not see the conversation that called it, only its task.
- Agents that edit the same files at once overwrite each other's work. Give each its own part of the task, or its own git worktree.
- When the calling session moves to another account ([`auto_switch`](/features/next#auto-switch)), the agents it started stop with it. Their conversations remain: `--resume <id>` continues them.
- `agent` is a command name, so no new account can be called `agent`; an existing one still starts with `ccam run agent`.

## Related

- [Picking an account for scripts](/features/pick)
- [Switching to the next account](/features/next)
- [Usage limits](/features/usage)
- [Background sessions of every account](/features/agents)
