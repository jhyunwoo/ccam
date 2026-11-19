---
title: Managing accounts
description: Add, log in, list, rename, remove and import Claude Code accounts with ccam add, ccam login, ccam list, ccam rename, ccam remove and ccam import.
updated: 2026-11-19
order: 2
translation: source
---

Every ccam account is a separate Claude Code config folder with its own login. Use these commands to add accounts, sign them in, list, rename and remove them. They also register folders you already have.

## Usage

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account that was never logged in or whose login expired
ccam logout side                           # log out: the login and its MCP sign-ins
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, logs out, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Wherever a command takes `<name|number>`, you can give the account name or its slot number (the `#` column in `ccam list`).

`ccam use` only changes which account is the default (`default` in `config.toml`).

`ccam add` on a terminal with no options opens a wizard. It goes through these steps:

1. the account name, unless you gave one;
2. how the account signs in: a Claude login now, an API key or a Claude login later;
3. which groups to join, with `default` checked (if no groups exist yet, it shows what the `default` group will share);
4. `Create a shortcut command for work? [y/N]`;
5. a summary to confirm.

`Esc` goes back a step. Any of `--group`, `--no-group`, `--shortcut` or `--no-login`, or running without a terminal (for example, from a pipe), skips the wizard. A name is then required, and no shortcut is created unless `--shortcut` names one.

Without `--group`, a new account joins the `default` group. The first account creates it, sharing MCP servers, settings, `CLAUDE.md` and extensions; see [sharing with groups](/features/share). On Windows, files can be linked only in Developer Mode. Without it, the first account creates `default` with only `mcp` and `extensions` and says so (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), and ccam refuses `--group` with a group that shares settings, memory or history.

`ccam list` shows the slot, name, shortcut, email, groups and last use, with `*` before the default account. The email is the one the folder is logged in as now (from its `.claude.json`); when that differs from the one ccam recorded, the recorded one follows as `(was …)`:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

```console
$ ccam list
   #  Name      Shortcut  Email                                     Groups   Last used
*  1  work      ccwork    other@company.com (was work@example.com)  default  2m ago
```

With `--json`, each account is an object with `slot`, `shortcut`, `name`, `email` (the current login's, as in the table), `default`, `dir`, `last_used` (left out when the account was never started), `groups` and `share`. `share` maps every item (`mcp`, `settings`, `memory`, `extensions`, `history`) to the group it comes from, or to `"local"`.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `ccam add --no-login` | Create the account without logging in; log in later with `ccam login <name>` | off |
| `ccam add --group <group>` | Join this group instead of `default`. Repeatable; `default` is created when it does not exist | `default` |
| `ccam add --no-group` | Join no group: the account shares nothing | off |
| `ccam add --shortcut <command>` | Also create this shortcut command in `~/.local/bin` | none |
| `ccam add --api --key-from <source>` | An account that uses an API key instead of a login; see [API key accounts](/features/api) | off |
| `ccam list --json` | Print the accounts as JSON | off |
| `ccam remove -y`, `--yes` | Remove without asking | off |
| `ccam logout -y`, `--yes` | Log out without asking | off |
| `ccam import --name <name>` | Account name for the folder you give; refused without a folder (`--name needs a folder: ccam import <folder> --name ops`) | derived from the folder name |
| `ccam import -y`, `--yes` | Import every folder found without asking | off |
| `--dry-run` | Supported by `ccam list` and `ccam completion`, which change nothing. Refused by `add`, `login`, `logout`, `use`, `rename`, `remove` and `import`, with `` --dry-run is not supported by `ccam add`; nothing was changed `` (naming the command) | off |

`ccam login` (besides `--expired`), `ccam use` and `ccam rename` have no options of their own. `ccam completion` takes `bash`, `zsh` or `powershell`; anything else fails with `only bash, zsh and powershell are supported`.

Account names are 1 to 32 characters of `a-z`, `0-9`, `-` and `_`, starting with a letter or digit. A name cannot be all digits, because `ccam 2` means slot 2. A name cannot be a ccam command either (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), so `ccam <name>` is never ambiguous. The list is `config.Reserved` in `internal/config`.

## How it works

Adding an account creates its folder, joins its groups, saves the config, signs in and creates the shortcut, in that order. When you run `ccam add work --shortcut ccwork`, ccam takes these steps (`addProfile` in `internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. ccam checks the name (`config.ValidateName`) and the shortcut (`shim.Check`) before it creates anything.
2. ccam creates `~/.local/share/ccam/profiles/work/` with mode `0700` and gives the account the lowest free slot number, so a removed account's number is reused (`profile.Add`, `Config.NextSlot`). The first account becomes the default.
3. The account joins its groups (`default` unless `--group` or `--no-group` says otherwise; the first account creates `default`), and ccam links the new, empty folder into them, for example `settings.json` to `groups/default/settings.json` (`applyShare`, with `share.Build` and then `Apply`). If this fails, the new folder is removed, because a leftover folder would block every retry with `folder already exists`.
4. The account is saved to `~/.config/ccam/config.toml` through `fsutil.AtomicWrite`: a temporary file, `fsync`, then `rename`, so a crash never leaves a half-written config.
5. ccam starts `claude auth login` with `CLAUDE_CONFIG_DIR` set to the account folder, and `--email` with the account's last email so the login page is filled in (`runLogin`). If the login fails, the account stays, and ccam prints the error and tells you to run `ccam login work` later.
6. ccam reads `oauthAccount` from `profiles/work/.claude.json` (`profile.ReadIdentity`): the email, the account and organization IDs, the organization's name and the plan (`claude_max`, `claude_team` …). That file can be several megabytes, so ccam scans for the one key instead of parsing the whole file (`TopLevelValue`). ccam stores these details in `[profile.login]` in `config.toml`.
7. ccam links `~/.local/bin/ccwork` to the ccam binary, which looks up `argv[0]` in `config.toml` when it starts, and records the link in `~/.local/state/ccam/shortcuts.json` (`shim.Install`). If the link cannot be created, the account is still added, ccam says why (`Could not create shortcuts: …`) and the summary shows `shortcut: ccwork (not created; see above)`.

A login can take minutes, and another ccam process may save the config meanwhile. So commands that wait reload `config.toml`, apply the one change (here, the email) and save (`env.update`). `ccam remove` reloads after its prompt for the same reason.

`ccam login` runs steps 5 and 6 for an existing account. Claude Code takes the email only as a hint for the login page, so you can end up signed in to a different Claude account than before; ccam compares the new login with `[profile.login]` and says so (`ccam: this login is other@company.com, but work was dev@company.com before.`). It keeps the new login and tells you how to undo it: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

A login into another Claude account also drops the account's [limit record](/features/usage#the-limit-record), since those limits were the other account's. When another ccam account is logged in to the same Claude account, ccam warns that the two share their usage limits.

`ccam login --expired` logs in, one after another, every account that was never logged in or whose login has expired: its refresh token expired, or Claude Code cleared the login after the server refused it. Before the first browser opens, it checks that `config.toml` can be saved.

`ccam logout` runs `claude auth logout` in the account and clears the saved email; the login record stays, so the next login is checked against it. Claude Code revokes the login and deletes what it stores for the account, including the MCP server sign-ins. So ccam warns you and asks first (`--yes` skips the question) when:

- the login holds MCP sign-ins;
- the account is running (its sessions are signed out too);
- another account's folder holds the same login (a copied folder), which the logout revokes as well.

For example:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

To log in as another Claude account, `ccam login <name>` alone is enough, and it keeps the MCP sign-ins.

`ccam list` adds an `Organization` column when the accounts belong to more than one organization (personal plans show `personal`), and `ccam list --json` has `plan` and `org` for every account.

### Signing in on a remote machine

Over SSH, the browser half of `claude auth login` can break. Claude Code waits for the sign-in on a callback server on the remote machine (`http://localhost:<port>/callback`), but the browser that finishes the sign-in runs on your computer, where that address leads nowhere and the browser shows an error such as `ERR_CONNECTION_REFUSED`. The address still carries the one-time code, so ccam can deliver it:

<!-- motion:accounts-relay -->

ccam does not hand your terminal to Claude Code (`runLogin`). It connects Claude Code's standard input to a pipe and copies each line you type into it (`relayLogin`), so a code pasted at Claude Code's `Paste code here` prompt still arrives. When a line is a callback address, ccam accepts it only if the host is `localhost`, `127.0.0.1` or `[::1]`, a port is given, the path is `/callback` and both `code` and `state` are present; `http://` may be left out (`callbackURL`). ccam then requests that address on the server with a 10-second timeout and without following redirects, the way the browser would have (`deliverCallback`). The request reaches the callback server of the waiting `claude auth login` (the dashed line in the diagram). Claude Code receives the code and saves the login, on Linux in `.credentials.json` in the account folder. A pasted line can never make ccam contact another host.

ccam prints this before the login starts:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

On success it prints `→ Passed the sign-in to Claude Code on this machine.`

While the login waits for input on a terminal, ccam turns off the terminal's echo and line mode and shows a `*` for each character you type, so a pasted code or callback address does not appear on screen. ccam does the line editing itself (`readLine`): `Enter` sends the line, `Backspace` erases one character, and `Ctrl+D` on an empty line ends the input. The terminal is restored when the login ends, and also on `Ctrl+C`, `SIGTERM` or `SIGHUP`, which still end ccam. Without a terminal on standard input and standard error, lines are copied as they come.

### Rename, remove and import

`ccam rename` and `ccam remove` act on the account's folder and on what ccam keeps under the account's name elsewhere. ccam lists those paths in `nameState`:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- the usage caches `~/.local/state/ccam/usage/<name>.json` and `<name>.statusline.json`;
- the limit record `<name>.limit.json`;
- the session notes in `~/.local/state/ccam/handoff/<name>/`.

`ccam remove work` handles them in this order:

<!-- motion:accounts-remove -->

1. `ccam remove` first refuses if the account is running, then asks, and checks again after you answer (this catches a session started while the prompt was open).
2. It checks that `config.toml` can be saved before it changes anything.
3. If the account shares extensions, ccam rewrites the group's `installed_plugins.json` and `known_marketplaces.json` (`RewritePluginPaths`), because Claude Code stores absolute plugin paths that may run through this account's folder.
4. A managed account (one whose folder ccam created, as opposed to one registered with `ccam import`) is logged out next: `claude auth logout` runs in its folder, so no working login waits in the trash (or, on macOS, in a Keychain item named after the folder, which a later account of that name would inherit). That is skipped for an API key account and when another account's folder holds the same login, which the logout would revoke too. A logout that fails is reported, and the removal goes on.
5. ccam moves the folder to `~/.local/share/ccam/trash/<name>-<time>/` (`profile.Trash`), the account leaves `config.toml`, the `nameState` paths are deleted so a later account with the same name does not inherit them, and the shortcut link is removed: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

If `config.toml` cannot be saved after all, the folder comes back from the trash. To recover a removed account, move the folder back from the trash and log in again.

`ccam rename` renames the folder and rewrites the plugin paths in its `plugins/` from the old folder to the new one. It moves the `nameState` paths to the new name and keeps the account as the default if it was one. It also updates the session notes, so `ccam next --continue` still finds them. The shortcut command keeps its name. If the plugin rewrite fails, or `config.toml` cannot be saved, the folder is renamed back.

On macOS, and on Windows when Claude Code stores logins in Credential Manager, the folder itself keeps its old name, because the login and the MCP sign-ins are stored under the folder's path and a moved folder would be logged out. The config then records the folder as `folder` beside the new name, and ccam says `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

`ccam import` registers a folder where it is, as `dir` in `config.toml`. It never moves or copies it: a copied refresh token would make the copies log each other out.

Without a folder argument, it looks for unregistered `~/.claude` and `~/.claude-*` folders and asks about each. Names come from the folder (`~/.claude` is `main`, `~/.claude-work` is `work`, `~/.claude-2` is `acct2`). If ccam cannot derive a name, or the name is taken, it lists the folder with the `--name` command to run. A found folder that fails is reported as `Skipped <path>: <reason>` and the others are imported.

ccam asks all its questions first, then reloads `config.toml` and registers the folders, so changes that other ccam commands saved in the meantime are kept. The `✓` lines come after the save.

Without `CLAUDE_CONFIG_DIR`, Claude Code keeps `.claude.json` in your home folder next to `~/.claude`, so when you import `~/.claude` and it has no `.claude.json` of its own, ccam creates the link `~/.claude/.claude.json` to `~/.claude.json`. Plain `claude` and `ccam main` then see the same account.

If `~/.claude` has a `.claude.json` of its own beside `~/.claude.json`, ccam leaves both and says so. `ccam main` would use the inner one and plain `claude` the outer one, keeping their login details and MCP servers apart. To keep one, replace the inner file with a link ([doctor](/features/doctor) reports it too).

When the folder holds the same login as a registered account (a copied folder), ccam says so: whichever renews it first logs the other out, so run `ccam login <name>` before using it.

Imported accounts are in no group until you add them to one (`ccam group join <group> <name>`), and no shortcut is created.

`ccam completion bash|zsh` prints a completion script that also completes account names.

## Safety and edge cases

- `ccam rename` and `ccam remove` refuse while the account is running (`refuseRunning`), with ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``, because the session keeps writing into the folder. The check covers Claude Code however it was started: through ccam, plain `claude`, an IDE, background sessions and their daemon.
- `ccam rename`, `ccam remove` and `ccam logout` refuse while another ccam process renews or changes that login: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- While a group change is in progress in another ccam process, `ccam rename` and `ccam remove` refuse with `another ccam process is changing the groups right now; try again when it is done`; after one stopped halfway, with ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.
- `ccam remove` on an imported account keeps the folder and its login, for plain `claude`. If the account is in groups, it first leaves them as `ccam group leave` does: every shared item becomes a copy of the group's in the folder, and the group's MCP servers go into its `.claude.json`. Then it is unregistered: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. The backups of its `.claude.json` that ccam created move into `<folder>/.ccam-backup/mcp-<time>`; if they cannot be moved, they stay in `~/.local/state/ccam/mcp/<name>/backups` and the message says so. Your `~/.claude` is not moved to the trash.
- Two accounts cannot have the same folder: `config.toml` naming one folder twice fails to load.
- `ccam add` refuses a name whose folder already exists in `profiles/`, and points to `ccam import`. On macOS it also stops when the Keychain still holds a login or MCP sign-ins for the new folder's path, left by an older version of ccam, and names the `security delete-generic-password` command that removes them. If `config.toml` cannot be saved, the new folder is removed again.
- `ccam import` refuses a folder that is already registered, a path that is not a folder, and a name that is taken. ccam also warns when the folder is a copy of `~/.claude` (the same login), even if `~/.claude` itself is not registered.
- `ccam logout` asks first when the account is running; `ccam login` does not check.
- If the relay cannot reach Claude Code, for example after the login ended, it prints `ccam: could not pass the sign-in to Claude Code:` with the reason. An HTTP error usually means the code expired; start the login again.

## Related

- [Launching Claude Code](/features/launch)
- [Shortcut commands](/features/shortcuts)
- [Share MCP servers, settings and history](/docs/sharing)
- [Use multiple Claude Code accounts on one computer](/guides/multiple-claude-code-accounts)
