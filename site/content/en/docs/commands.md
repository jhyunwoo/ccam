---
title: Command reference
metaTitle: ccam command reference
description: Every ccam command and option, from ccam add and ccam <name> to ccam usage, ccam next and ccam doctor.
updated: 2026-11-19
published: 2026-10-19
order: 3
translation: source
---

Put `--dry-run` before a command to see what it would do without doing it, for example `ccam --dry-run work`. It works when starting Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) and with `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` and `doctor --fix`. Commands that only read, such as `list`, `usage`, `share` and `ext`, accept it and change nothing anyway. The rest refuse it instead of running, including `group rename`.

<!-- commands -->

## Starting Claude Code

`ccam <name|number> [arguments]` starts Claude Code as that account. ccam does not read the arguments after the name; it passes them to `claude` in the same order, and standard input, the exit code and signals behave exactly as with `claude` itself. A `--` after the name is optional.

`ccam` with arguments but no account name, such as `ccam -p "summarise"`, uses the account set for the current folder, otherwise the default account. Set them with `ccam folder <name> [folder]` and `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` before the account name is for `-p`/`--print` runs. When the account reaches its usage limit, ccam runs the same command again as the logged-in account with the most room and prints only the output of the run that finished. See [Accounts by folder](/features/folder) and [Run again at the limit](/features/next#run-again-at-the-limit).

**Team and Enterprise seats.** Automatic account selection (`ccam next`, `--fallback`, `ccam pick`, the dashboard's `n` and, once you turn it on, [`auto_switch`](/features/next#auto-switch)) also switches between seats of the same Claude Team or Enterprise organization, which may violate Anthropic's terms. Read the [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) and check with your organization and Anthropic before you rely on it.

A shortcut such as `ccwork` starts its account the same way: `ccwork --resume` does what `ccam work --resume` does. See [Shortcuts](/docs/commands#shortcuts) below.

## Adding and removing accounts

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` without a folder looks for `~/.claude` and `~/.claude-*` and asks about each one; a folder that fails is reported as `Skipped <path>: <reason>` and the others are registered. If a folder's name is already taken, import that folder on its own with `ccam import <folder> --name <name>`. Imported folders are registered where they are; ccam never moves or copies a login.

`ccam login --expired` also logs in accounts whose login Claude Code cleared. `ccam remove` logs out a managed account (one whose folder ccam created, as opposed to one registered with `ccam import`) before its folder goes to the trash, unless another account holds the same login. Removing an imported account takes it out of its groups and leaves a copy of each shared item in its folder. On macOS, `ccam rename` keeps the folder's name, since the Keychain keeps the login under the folder's path.

`ccam add --api` adds an account that uses an API key instead of a Claude login; ccam reads the key at each launch and never stores it. See [API key accounts](/features/api).

`ccam rename`, `ccam remove` and the `ccam group` commands refuse to change an account while Claude Code is running as that account, however it was started (including background sessions). Quit the session first. They also refuse while a group change is unfinished: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` prints every account with its slot, shortcut, email and folder. Its `groups` array lists the account's groups, and its `share` object maps each item (`mcp`, `settings`, `memory`, `extensions`, `history`) to the group it comes from, or `local`.

## Shortcuts

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

A shortcut is a command you name yourself that starts Claude Code as one account. Like `ccam work`, it passes every argument to Claude Code unchanged, so `ccwork --resume` and `ccwork -p "hi"` work as you would expect. It is a link to ccam in `~/.local/bin`.

Shortcuts are optional. The `ccam add` wizard asks whether to create one, `ccam add <name> --shortcut <command>` creates one without asking, and `ccam list` shows each account's shortcut, or `-` when it has none.

Names use the characters a-z, 0-9, `.`, `-` and `_`. A name cannot be a ccam command, another account's shortcut, or the name of another program on your `PATH`, nor a file in `~/.local/bin` that ccam did not create, such as your own copy or link of ccam. A shortcut edited by hand in `config.toml` that breaks these rules, is `claude` or `ccam`, or is given to two accounts makes `config.toml` fail to load, naming the account.

`ccam rename` keeps the shortcut name, and `ccam remove` deletes the shortcut together with the account. `ccam shims install` recreates the links of all accounts, for example after you have moved ccam; a file in the way of one shortcut does not stop the others.

## Sharing

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

A group shares items between its accounts; an account can be in several groups, but each item comes from one group only. See [sharing with groups](/features/share).

## Copying an existing setup

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` copies `settings.json`, `CLAUDE.md`, extensions and user MCP servers from a Claude Code folder into a group, so every account that shares those items gets them. The folder is only read.

- `--settings`, `--memory`, `--extensions`, `--mcp` and `--history` copy only the items you name. Without them, a checklist on a terminal lets you choose, with everything except conversation history checked; elsewhere everything except conversation history is copied.
- `--mcp` also works when `~/.claude` is empty or missing but `~/.claude.json` has servers.
- `-g` or `--group` picks the group.
- `-y` or `--yes` copies without asking and keeps the group's copy when a file differs.

MCP sign-ins are not copied, since a copied sign-in logs the other holders out; the plan names the accounts that sign in with `/mcp`. See [copying an existing setup](/docs/sharing#copy-an-existing-setup).

## MCP servers

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` moves the servers that only that account has into its group, including those it has for one folder only (plain `claude mcp add`), which then load in every folder. A server whose name the group already uses for a different definition is skipped; add `--force` to replace the group's copy, which ccam backs up first.

## Usage and switching

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam sends no token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --continue --brief    # a new conversation from a summary of the last one
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agent "…"                  # from Claude Code: a subagent that runs as another account
ccam agent --resume <id> "…"    # follow that agent up
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
ccam ccusage claude daily       # Claude Code's tokens and costs, every account (ccusage)
ccam ccusage -a work blocks --active   # one account's current 5-hour block
```

- `next --auto`, `--fallback` and `ccam pick` choose the same way: below `auto_max_usage` (`config.toml`, default 100) first, by `auto_strategy` (`headroom`, `expiring` or `balance`), counting the `--model` weekly window. They never choose an account at 100% (a window Claude Code reported full counts as 100% until its reset) or an API key login, and they stay within the organization of the account they start from (`--cross-org` or `cross_org = true` to look further).
- `--wait` takes its value after `=` (`--wait=2h`); written with a space it is refused with a hint. When none can be chosen they exit with code 75, or, with `--wait`, wait for the first account to reset.
- `ccam next --continue` carries a conversation over even between accounts that do not share history: Claude Code then copies it into the new account. With `--brief`, the new account starts a new conversation instead, from a summary ccam writes of the last one without a model ([`--brief`](/features/next#brief)).
- Unlike `ccam <name>`, `ccam next` reads its own options, so Claude Code's options go after `--`; without it ccam stops with `unknown flag`.
- `ccam agent` runs a task with `claude -p` in this folder as another account than the one it is called from, in the calling session's permission mode and effort, and prints Claude's last reply. At that account's limit the conversation goes on as another ([Subagents on other accounts](/features/agent)).
- Inside Claude Code, the plugin's `/ccam:switch [account]` hands a conversation on, and `ccam hook` is what its hooks run.
- With `auto_switch = true` in `config.toml` (Linux and macOS; off by default), the limit hook also ends an interactive session that ccam started and continues it on another account with room, but only one of the same kind: another seat of the same organization, or another personal plan ([`auto_switch`](/features/next#auto-switch)).

See [Switching to the next account](/features/next), [Picking an account for scripts](/features/pick), [Subagents on other accounts](/features/agent) and [Background sessions of every account](/features/agents).

## Maintenance

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam uninstall       # remove every account, ccam's folders and ccam itself
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` checks the release's SHA-256 checksum before it replaces ccam, and `--version v0.2.0` installs a particular release. A copy of ccam installed from the `.deb` package is updated with apt instead. See [update](/docs/install#update).

`ccam uninstall` lists what it removes and asks first: every account as `ccam remove` does, the shortcut commands, ccam's folders and the `ccam` command. An imported folder stays where it is. See [uninstall](/docs/install#uninstall).

## Files

| Path | Contents |
| --- | --- |
| `~/.config/ccam/config.toml` | Accounts, groups and settings (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). Settings go at the top, above the first `[[group]]`, `[[profile]]` or `[[folder]]`: one below a table belongs to that table, so ccam warns about it at every launch and in `ccam doctor`, and commands that save the config refuse until it is moved. Comments are not kept when ccam saves the file |
| `~/.local/share/ccam/profiles/<name>/` | Each account's `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | Shared settings, MCP servers, extensions and history |
| `~/.local/state/ccam/` | Usage cache, limit records, locks, notes and other state (see the list below the table) |
| `~/.local/bin/<shortcut>` | Shortcut commands such as `ccwork`, links to ccam |

The state folder holds:

- the usage cache and limit records (`usage/<name>.json`, `usage/<name>.limit.json`);
- last-used times;
- session markers and locks (`locks/`);
- the shortcuts ccam created (`shortcuts.json`);
- a group change that has not finished (`group-change.json`);
- the notes `/ccam:switch` and the limit hook leave (`handoff/<account>/<session>.json`);
- when `ccam pick` or `auto_switch` last chose each account (`picks/<name>`; a choice counts for 2 minutes);
- the summaries of `--brief`, kept 7 days (`brief/<session>.md`).

Setting `XDG_CONFIG_HOME`, `XDG_DATA_HOME` or `XDG_STATE_HOME` changes where these folders are. On Windows, `~` is `%USERPROFILE%` and shortcuts are `<shortcut>.exe`.
