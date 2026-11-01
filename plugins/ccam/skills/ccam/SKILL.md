---
name: ccam
description: Use when the user has several Claude Code accounts managed by ccam and asks about usage limits, which account has room, switching accounts after a limit, running parallel `claude -p` workers on different accounts, or why a skill, command or MCP server is missing in one account. Triggers include "which account has room", "I hit my limit", "run these tasks on my other accounts", "spread the work over my accounts", "ccam".
---

# Working with ccam

ccam gives each Claude Code account its own `CLAUDE_CONFIG_DIR` and starts Claude Code as the account you name. It is installed as the `ccam` command. These commands are safe to run from this session; they only read, except where noted.

## See usage

```sh
ccam usage --json     # 5-hour and weekly usage of every account
ccam list --json      # accounts, emails, what they share
```

In the JSON, `limits.five_hour.percent` and `limits.weekly.percent` are 0-100, and `resets_at` says when each window starts over. `status` is `api` for API key accounts, which have no limits.

## When the current account reaches its limit

You cannot move this session to another account yourself, and Claude Code does not run `/ccam:switch` from a prompt you schedule; the user can. Tell them to type `/ccam:switch` (or `/ccam:switch <account>`): it notes this conversation and ends the turn. With `relay = true` (or `auto_switch = true`, on Linux and macOS) in ccam's `config.toml`, ccam then ends this session and continues the conversation on the other account by itself (on Windows, once they type `/exit`), with the permission mode, effort level, agent and options such as `--add-dir` it had. For an account they name in another organization than this session's, ccam does not end the session: they type `/exit`, and ccam asks them on the terminal before it moves the conversation. It asks the same before an account of another organization that it picks itself, which only `cross_org = true` allows. Otherwise they type `/exit` and run, in this folder:

```sh
ccam next --continue          # as the account /ccam:switch named, else the next one (asks first for another organization)
ccam next --auto --continue   # or the account with the most room
```

When a session stops at a usage limit, the ccam plugin notes it the same way, and `--auto` leaves that account out. If the accounts share conversation history, this very conversation continues; if not, Claude Code copies it into the new account and the original stays as it is.

`ccam next --continue --brief` starts a new conversation as the other account instead, which reads a summary ccam writes of this one without a model rather than all of it: that costs the other account much less, but details are left out. Suggest it for long conversations.

With `auto_switch = true` in ccam's `config.toml` (Linux, macOS; off by default), ccam continues the conversation on another account by itself when Claude Code stops the session at a usage limit and another account like this one has room. The new session then starts with: "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." When you see that message, carry on from where the conversation stopped. Background shells, dev servers and agents started before the move have stopped; start again only what the task still needs.

## Run work on other accounts

For non-interactive work, start a new `claude -p` through ccam. Pick the account with `ccam pick`; calls made close together return different accounts, so parallel workers spread out:

```sh
ccam "$(ccam pick)" -p "fix the failing test in pkg/parse" --output-format json
```

Or let ccam switch when an account hits its limit mid-job:

```sh
ccam --fallback "$(ccam pick)" -p "…" --output-format json
```

- Exit code 75 means no account has room now. Do not retry in a loop; report it, or add `--wait` (`ccam pick --wait`, `ccam --fallback --wait …`), which waits for the first account to reset, up to 5 hours by default.
- `ccam pick --model opus` also counts that model's weekly limit.
- Give each worker its own git worktree when they edit the same repository.
- Output of a `--fallback` run is held until it finishes, so stream-json arrives at the end.

## A skill, command or MCP server is missing in one account

```sh
ccam ext            # what every account loads, and from where
ccam ext work --json
ccam doctor         # broken share links, logins, permissions
```

`ccam ext` shows each skill, command, agent, plugin, MCP server and hook per account, with its source (a share group name, `account`, `local` for one folder only, or `project`), and lists problems such as broken links, skill folders without `SKILL.md`, MCP servers defined differently in two places or only in one account, and servers another account is signed in to. Suggest `ccam doctor --fix` for broken share links and `ccam mcp import <account> <server>` to share a server the account added on its own.

- Plain `claude mcp add` gives one account a server for one folder only; `ccam mcp add` (or `claude mcp add -s user`, then `ccam mcp import`) shares it.
- A server that needs a sign-in (OAuth) is signed in once per account, with `/mcp` in a session of that account; ccam never copies sign-ins. claude.ai connectors belong to each claude.ai account.
- With the default `mcp_strategy` (flag), `claude mcp list` does not show the group's servers (sessions get them with `--mcp-config`): use `ccam ext` or `ccam mcp list` to check. Under `mcp_strategy = "sync"` ccam copies them into the account's `.claude.json`, so `claude mcp list` shows them.
- Servers left in plain claude's `~/.claude.json` load in no account until `ccam migrate --mcp` copies them; suggest it, do not run it.

## Do not

- Do not read or print `.credentials.json` or any token; ccam never needs them in a conversation.
- Do not run `ccam add`, `ccam login`, `ccam group` (other than `ccam group` and `ccam group show`, which only read), `ccam remove` or `ccam migrate` yourself: they log in through a browser or move files. Suggest the command to the user instead.
- Do not turn on `auto_switch` or `relay` in ccam's `config.toml` yourself; that is the user's choice. On Team and Enterprise seats `auto_switch` moves work between seats of the organization, which Anthropic's seat rules may not allow: point the user to https://ccam.dev/en/features/next/#auto-switch.
