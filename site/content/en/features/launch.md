---
title: Launching Claude Code
description: Start Claude Code as any account with ccam <name|number>, ccam run or the dashboard; arguments go to claude unchanged, and --dry-run shows the command.
updated: 2026-11-01
order: 1
translation: source
---

`ccam work` starts Claude Code as the `work` account, and everything after the name goes to `claude` unchanged. Use it wherever you would type `claude`: in a terminal, in a script, or with input piped in.

## Usage

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Give an account's name or its number from `ccam list`. ccam does not read the arguments after the name, so `ccam work --help` shows Claude Code's help. A `--` after the name is optional. Standard input, the exit code and signals behave as with `claude`.

When the first argument starts with `-` and is not one of ccam's own flags, ccam picks the account for you and passes all arguments to Claude Code. It uses the account set for the current folder with [`ccam folder`](/features/folder), otherwise the default account: the one set with `ccam use`, or the account with the lowest number if you have not set one.

### The dashboard

`ccam` without arguments opens the [dashboard](/features/dashboard) when standard input and output are a terminal, and prints `ccam list` otherwise. Move to an account (with the arrows, or by typing its number) and press `Enter` to start it; `n` starts the account with the most room and `o` offers `--continue`, `--resume` or other arguments. The dashboard closes before Claude Code starts, so the terminal is restored first.

**Team and Enterprise seats.** Automatic account selection (`ccam next`, `--fallback`, `ccam pick`, the dashboard's `n` and, once you turn it on, [`auto_switch`](/features/next#auto-switch)) also switches between seats of the same Claude Team or Enterprise organization, which may violate Anthropic's terms. Read the [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) and check with your organization and Anthropic before you rely on it.

When the shell has `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` or `ANTHROPIC_AUTH_TOKEN` set, Claude Code may use that key or token instead of the account's own login, so ccam warns you on standard error before it starts (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). It only reads the variables; an [API key account](/features/api) clears them instead.

By default, ccam replaces itself with Claude Code, so nothing else keeps running. With `relay = true` (or [`auto_switch = true`](/features/next#auto-switch) on Linux and macOS) at the top of `config.toml`, above the first `[[group]]`, `[[profile]]` or `[[folder]]`, ccam instead starts interactive sessions as a child process and keeps running alongside them. That lets `/ccam:switch`, or `auto_switch` at a usage limit, continue the conversation on another account by itself, with the session's options, permission mode, effort level and agent; before `/ccam:switch` moves it to an account of another organization, ccam asks you ([details](/features/next#relay)).

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--dry-run` | Print the environment and command ccam would run, then exit | off |
| `--fallback` | For `-p`/`--print` runs: when the account reaches its usage limit, run the same command again as the logged-in account with the most room ([details](/features/next#run-again-at-the-limit)) | off |
| `-v`, `--version` | Print `ccam <version>` (only as the first argument) | |
| `-h`, `--help` | ccam's help, as `ccam --help` or `ccam run --help` | |

Argument rules:

- ccam reads its own flags only before the account name. After the name, every argument belongs to Claude Code, including `--dry-run` and `--help`.
- `--dry-run` works as `ccam --dry-run work`, `ccam run --dry-run work` and `ccam --dry-run` (for the folder's account or the default). `--dry-run` and `--fallback` can come in either order.
- `ccam run <name|number> [claude args]` is the long form of `ccam <name>`. `ccam run` without a name uses the folder's account or the default one.
- Account names cannot be ccam command names such as `add`, `list`, `run` or `usage`, so `ccam <word>` is never ambiguous.

The dry run prints one line you could paste into a shell, or into PowerShell on Windows. It starts with what the launch changes in the environment: `CLAUDE_CONFIG_DIR` and `CCAM_PROFILE` always, `env -u <name>` for a variable it drops (such as an inherited `CCAM_RELAY`), and, for an [API key account](/features/api), its endpoint, its model and where the key comes from (`ANTHROPIC_API_KEY="$WORK_KEY"`), never the key itself:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

The PowerShell line puts the variables back afterwards, as `VAR=x cmd` does in a shell, and hands Claude Code its command line through `Start-Process`, so Windows PowerShell 5.1 and PowerShell 7 run it the same way. When a launch sets a variable to an empty value, a `#` line below it says that PowerShell before 7.5 removes such a variable instead. For a `claude.cmd` installed by npm, the line runs the `cmd.exe` command the launch would use (see below), and an argument the launch would refuse there fails the dry run too.

## How it works

Each launch takes the same steps, from the command line to `exec`:

<!-- motion:launch -->

Each box lights up as its step runs, and the lines underneath are the environment and arguments Claude Code starts with.

1. ccam starts in `app.Main`, which `cmd/ccam/main.go` calls. If ccam was started under another name, it is a [shortcut](/features/shortcuts) and goes to the launcher directly.
2. ccam looks at `os.Args` before cobra is involved (`app.parse`). Only `ccam` alone (the dashboard), a ccam command name (`config.Reserved`), `-h`/`--help` and `-v`/`--version` go to `cli.Execute`. Everything else is a launch.
3. ccam reads `~/.config/ccam/config.toml`, the only file it parses (`launcher.Launch`).
4. ccam finds the account by name or number, drops one leading `--`, and resolves the `claude` executable (`launcher.Prepare`, `ResolveClaude`). Without a name, it takes the folder's account, comparing the working directory with the `[[folder]]` entries of `config.toml` as strings, without regard to case on macOS and Windows (`HereProfile`). The PATH lookup is cached in `~/.local/state/ccam/claude-path` together with the file's modification time and a hash of `PATH`, so ccam does not search `PATH` on every launch. A changed or missing file, or a changed `PATH`, triggers a new lookup. A `claude_bin` path in `config.toml` (`~/` is expanded) must be an executable file; otherwise the launch stops with `claude_bin X is not an executable file. Check claude_bin in config.toml` before anything is written. ccam refuses a `claude` that is ccam itself, compared by inode, because that would loop.
5. ccam compares the account folder `~/.local/share/ccam/profiles/<name>/` with its sharing settings, calling only `lstat` and `readlink`, one or two system calls per item (`checkLinks`, `share.Check`). A missing link is created on the spot. Anything that could lose data stops the launch.
6. ccam adds the shared MCP servers (`applyMCP`; see below).
7. ccam records the session (`session.Register`). It writes a marker named after its PID, holding the process start time, to `~/.local/state/ccam/locks/profile-<name>/<pid>`, whatever `XDG_RUNTIME_DIR` says, so every shell sees the same markers. For an account that shares conversation history, ccam first looks for a running session in `history-<group>-<hash>`, where the hash is of the current folder with symbolic links resolved, and registers there too (`historyCheck`).

   `exec` keeps the PID, so the marker then names the Claude Code process itself. ccam uses a marker file rather than a held lock, because processes that Claude Code starts would inherit a lock and keep it after the session ends. A marker that cannot be written does not stop the launch: ccam says `ccam: could not mark the session as running (…); starting Claude Code anyway`.

   ccam's other checks of whether an account is running (rename, remove, group changes, `doctor --fix`, the MCP sync) read Claude Code's own records in the account folder as well (`sessions/<pid>.json`, `daemon.lock`), so they also see Claude Code started without ccam: plain `claude`, an IDE, background sessions and their daemon.
8. ccam sets the terminal title to `claude · <name>` (when standard error is a terminal), writes the launch time to `~/.local/state/ccam/last-used/<name>` without `fsync`, and calls `syscall.Exec` (`Plan.Exec`).

Claude Code runs with `CLAUDE_CONFIG_DIR` set to the account folder and `CCAM_PROFILE` set to the account name. ccam clears any inherited values of these variables, and an inherited `CCAM_RELAY`. For a `~/.claude` folder registered with `ccam import` on macOS, it also sets `CLAUDE_SECURESTORAGE_CONFIG_DIR` to an empty value, so that the account and plain `claude` use the same Keychain login; on Windows it does the same when Claude Code keeps the login in Credential Manager.

### Why exec

ccam replaces its own process with `claude` (`syscall.Exec`). As a result:

- ccam does not stay in memory;
- keystrokes, signals and the exit code are Claude Code's own;
- arguments stay an array, never a shell string.

Windows has no `exec`. There, ccam starts `claude.exe` as a child process with the same environment, arguments and console, waits for it and exits with its exit code (`Plan.Exec`). A `claude.cmd` that npm installed runs in `cmd.exe` with its full path and every argument quoted so that `cmd.exe` reads none of them as its own syntax (`&`, `%PATH%`); an argument with a line break cannot be passed that way and is refused, with the advice to install the native Claude Code.

Ctrl+C reaches Claude Code from the console directly, so ccam ignores it. The session marker carries ccam's PID, which lives exactly as long as Claude Code, and is removed when Claude Code exits. While it waits, ccam does nothing else.

Work after a session, such as noticing servers added with `claude mcp add`, happens at the next launch: ccam compares the modification time of `.claude.json` with `~/.local/state/ccam/mcp/<name>/seen` and parses the file only when it changed.

### How a command line is routed

Three command lines, three routes through `internal/app`:

<!-- motion:launch-route -->

A shortcut is matched on `argv[0]`, `ccam <name>` on the first argument, and both end in `syscall.Exec`. Only a ccam command such as `ccam usage` reaches cobra.

### Shared MCP servers

With the default `flag` strategy, ccam passes the group's MCP servers to Claude Code as an argument, placed before yours:

<!-- motion:launch-mcp -->

ccam adds the argument only if the group's `mcp.json` exists, and puts it right after `claude` (`mcp.FlagArgs`, `Plan.InsertArgs`). Your own `--mcp-config` or `--strict-mcp-config` then still applies on top. It is one argument, `--mcp-config=<path>`, because Claude Code's `--mcp-config` takes several values: written as two arguments, it would also read a following prompt as a config file.

ccam decides where the flag cannot go (`claudeVerb`). It reads Claude Code's subcommand from the arguments the way Claude Code does, skipping options together with their values. `remote-control` or `rc` gets no flag (below), nor does a background-session verb given first (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`). `ccam work mcp list`, `get` and `remove` get none either, since Claude Code refuses it there; instead ccam says on standard error that the group's servers load in the account's sessions too:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Every other subcommand and every session keeps the flag.

With `mcp_strategy = "sync"` in `config.toml`, ccam writes the servers into the account's `.claude.json` instead, and skips that if the account is already running. A launch skips the work too when neither the group's `mcp.json` nor `.claude.json` changed since the last check. If the account has its own server with the same name as a group server, sync keeps the account's server and says so once. Sync also puts back a copy that is gone from `.claude.json` (`put back group default's MCP servers …`), and changes or removes only the copies it wrote that nobody changed since.

### Remote Control

`ccam work rc` (or `ccam work remote-control`) starts Claude Code's Remote Control as `work`. `claude remote-control` refuses `--mcp-config` before its verb and takes none after it, and its sessions read MCP servers only from the account's `.claude.json`. So the flag strategy works differently here:

<!-- motion:launch-rc -->

1. ccam finds `remote-control` or `rc` as Claude Code's verb, even after options and their values (`ccam work --model opus rc`), and adds no flag (`claudeVerb`, `mcp.FlagArgs`).
2. ccam compares the modification time of the group's `mcp.json` with the account's `sync-stamp`, and checks that every group server is still in `.claude.json` (`mcp.Synced`). If `.claude.json` already holds the group's current servers, or the group has none, Remote Control starts without a question.
3. Otherwise ccam names the servers the sessions would miss and asks:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. On `y`, ccam copies them the way the `sync` strategy does (`mcp.Sync`): under Claude Code's lock, after a backup, keeping the account's own servers. On `N`, or without a terminal, Remote Control starts without them, and ccam points to `mcp_strategy = "sync"`, which copies them before every launch.

ccam asks on each Remote Control launch until the servers are copied. After a copy, it asks again only when the group's `mcp.json` changes or a copied server disappears from `.claude.json`. A copy in `.claude.json` wins over `--mcp-config` on later normal launches, so `ccam mcp add`, `remove` and `import` update the copies other accounts hold ([Shared MCP servers](/features/mcp)).

### The 10 ms budget

The budget for the launch path is 10 ms at the 95th percentile; in practice ccam adds about 4 ms before Claude Code starts. The launch path never builds the cobra command tree, starts Bubble Tea or opens a network connection, and `TestPackageInitIsCheap` fails if package `init()` functions take more than 5 ms in total. `scripts/perf-check.sh` measures `ccam --dry-run work`, `ccam work` and a shortcut with `hyperfine` and fails above 10 ms, and `BenchmarkLaunchPath` in `internal/launcher` runs the steps in a loop. To see the time per step, set `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Safety and edge cases

- An unknown name or number stops with `no account "<name>"` and points to `ccam list`. With no accounts at all, ccam asks you to add one with `ccam add <name>`.
- If a shared path became an ordinary file or points to the wrong place, ccam stops with `stopped because of share link problems:`, lists them, and asks you to run `ccam doctor --fix`. It does not move data on the launch path.
- If another Claude Code session that shares conversation history with this account (same group) is running in the same folder, ccam warns that the history can mix and asks `Run anyway? [y/N]`. The folder is matched with symbolic links resolved. Without a terminal (standard input from `/dev/null` counts as none), or from a job in the background, it shows the warning and continues. A new `-p` run starts a conversation of its own and gets no warning; `--version`, `--help` and subcommands that open no conversation, such as `mcp`, `auth` or `plugin`, are not checked at all.
- `--dry-run` skips the MCP sync, the new-server notice, the `mcp list` note, the history check, the session marker, the title and the last-use time, and does not exec. It still checks share links, and creates a missing link as a normal launch would.
- `ccam work rc` while `work` is already running leaves `.claude.json` alone, says so, and starts Remote Control without the group's servers; run it again after the session exits to copy them. With `--dry-run` it prints the message but neither asks nor copies.
- With `mcp_strategy = "sync"`, a failed sync prints `ccam: MCP sync failed (launching anyway): …` and Claude Code starts with the servers it already had.

## Related

- [Shortcut commands](/features/shortcuts)
- [Using Remote Control with ccam](/guides/remote-control)
- [Managing accounts](/features/accounts)
- [Command reference](/docs/commands)
- [Share MCP servers, settings and history](/docs/sharing)
- [CLAUDE_CONFIG_DIR explained](/guides/claude-config-dir)
