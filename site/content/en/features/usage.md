---
title: Usage limits
description: ccam usage shows the 5-hour and weekly limits of every account at once, and ccam statusline puts the account name and usage in Claude Code.
updated: 2026-11-01
order: 7
translation: source
---

`ccam usage` shows, for every account, how much of each limit is used: the 5-hour limit, the weekly limit and each model's own weekly limit (such as Fable). Each limit gets a bar and the time until it resets. Use it before you start a long task, or when Claude Code tells you a limit is near, to see which account has room left. `ccam statusline` is the status line command for Claude Code: it shows the account name and usage inside Claude Code and saves the numbers for `ccam usage`. For the tokens and estimated costs behind the usage, `ccam ccusage` runs [ccusage](https://ccusage.com) over every account's conversations ([below](#tokens-and-costs-ccam-ccusage)).

## Usage

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts that are not running
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

Each account gets a line with its number, its name and how old its numbers are; `(cached)` means they come from a saved file, not from the endpoint. Under it, one line per limit shows a 20-cell bar, the percentage and the time until the window resets: `5-hour`, `Weekly` (all models), then each model's own weekly window, such as `Fable`. A window the account does not have shows `-`. A percentage turns yellow at 80% and red at 100%.

If Claude Code refused a request because a window was full, that window shows 100% until it resets, whatever the saved numbers say (see [the limit record](#the-limit-record)). On a terminal every account appears at once as `fetching…`, and its lines are drawn as it answers.

The last line, shown only when you have more than one account, names the account [`ccam next --auto`](/features/next) would choose. It follows the same rules:

- `auto_max_usage` and `auto_strategy`;
- the running sessions;
- only accounts of the last used account's organization, unless `cross_org` is set;
- every window at 100% is skipped, so an account whose week is full is never named.

With the default `headroom` strategy it reads `→ Most headroom: <name>`; with `expiring` or `balance` it gives the reason, such as `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Team and Enterprise seats.** Automatic account selection (`ccam next`, `--fallback`, `ccam pick`, the dashboard's `n` and, once you turn it on, [`auto_switch`](/features/next#auto-switch)) also switches between seats of the same Claude Team or Enterprise organization, which may violate Anthropic's terms. Read the [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) and check with your organization and Anthropic before you rely on it.

An account that is not logged in shows `not logged in → ccam login <name>`, and an API key account shows `no limits (API)`. An expired access token shows ``access token expired; `ccam <name>` renews it when it starts``, above the saved numbers, if any: Claude Code renews it at the next start. A login that has fully expired, because its refresh token expired or Claude Code cleared it after the server refused it, shows `` login expired → `ccam login <name>` ``.

After `ccam statusline install`, Claude Code shows a line such as `[work] · Opus · 5h 24% · 7d 41%`.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--json` | Print the reports as JSON: `name`, `slot`, `status`, `limits` (with `credits` when the account has usage credits on), `source`, `at`, `error`, and `limited` and `limited_windows` when a limit record sets windows to 100%. With no accounts it prints `[]` | off |
| `--watch[=N]` | Redraw every `N` seconds until you press Ctrl+C. `--watch` alone, or a value of 0 or less, means 60. The value goes after `=`: `--watch 10` is refused with a hint | off |
| `-v`, `--verbose` | Add a `source:` line under each account: where its numbers came from (the usage endpoint, Claude Code, the status line, a saved file, or the usage limit Claude Code reported). When a limit record overrides the numbers, a line such as `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` follows | off |
| `--no-cache` | Ignore the 60-second cache and ask the endpoint again | off |
| `--source live\|claude\|cache` | `live` asks the endpoint with the account's token; `claude` has Claude Code look it up (below); `cache` shows saved numbers only | `usage_source`, otherwise `live` |
| `--refresh` | Renew expired access tokens, only for accounts that are not running. A renewal may take up to 30 seconds | off |
| `--dry-run` | Accepted: `ccam usage` changes nothing but its own cache. Refused together with `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | off |

Any other `--source` value fails with `--source must be live, claude or cache`. `ccam usage` takes no arguments.

### Usage from Claude Code itself

With `usage_source = "claude"` at the top of `config.toml` (above the first `[[group]]`, `[[profile]]` or `[[folder]]`), or `--source claude`, ccam does not read the login token at all. For each account it starts `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` in the account's folder and sends Claude Code's `get_usage` request on standard input (Claude Code 2.1.270 or later). Claude Code looks the usage up with its own login and answers with the 5-hour, weekly and model windows.

`--safe-mode` keeps the account's hooks, MCP servers and plugins from starting, and no session is saved. Each account needs one Claude Code start (about a second), so this source is slower than `live`. The accounts are asked at once, each with a 20-second timeout, and the 60-second cache applies as usual.

A login whose access token has expired, or expires within 5 minutes, gets 50 seconds, since Claude Code renews it as it starts and a renewal cut off may lose the login. If ccam gives up on a lookup, it asks Claude Code to exit with SIGTERM, so Claude Code finishes a renewal it has started (Windows has no SIGTERM, so there the process is terminated).

A login that has fully expired shows `` login expired → `ccam login <name>` `` without starting Claude Code. Claude Code marks the request as experimental, so if it stops answering, ccam shows the saved numbers and `ccam doctor` reports it. When Claude Code answers that its own lookup failed, ccam shows the saved numbers with `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`, not a plan without limits.

With `--dry-run`, ccam does not start Claude Code, because Claude Code writes to the account's folder; it shows the saved numbers instead. `usage_source` takes `live` or `claude`; `cache` is a `--source` choice only.

`ccam statusline` reads Claude Code's status line JSON from standard input. Its one option is `--wrap <command>`: an existing status line command whose first output line gets the account name in front. `ccam statusline install` takes no options. Both refuse `--dry-run` with `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Tokens and costs: `ccam ccusage`

`ccam usage` shows how much of each limit an account has used. [ccusage](https://ccusage.com) counts the tokens and estimated costs in Claude Code's conversation logs instead. On its own it looks for Claude Code's logs only in `~/.claude` and `~/.config/claude` (or the folders `CLAUDE_CONFIG_DIR` names), so the only ccam account it can see is one imported in place at `~/.claude`. `ccam ccusage` runs it with `CLAUDE_CONFIG_DIR` set to the folders of every account that has conversations, so one report covers them all:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam reads its own options first: `-a` (`--account`) takes account names or numbers, comma-separated or repeated, and `--dry-run` prints the command instead of running it. From the first other argument on, everything goes to ccusage as you typed it, except `--config` for `bun x` and `npx` (below); `ccam ccusage -- --help` shows ccusage's own help. ccam runs `ccusage` if it is installed, otherwise `bun x ccusage`, and failing that `npx ccusage@latest`; the last two fetch it from npm. They start in ccam's own folder (`~/.local/state/ccam/ccusage`), not in the folder you are in, so a project there cannot choose what they run with its `node_modules/.bin`, a `bin` in its `package.json` or a registry in its `.npmrc`; `ccam --dry-run ccusage` shows the `cd` into that folder. When none of them is available, it tells you to install Node.js, Bun or ccusage. On Linux and macOS ccusage takes ccam's place; on Windows ccam waits for it and ends with its exit code. Under WSL, ccam skips a `ccusage`, `bun` or `npx` on a Windows drive, which would run as a Windows program without the account folders.

- The report covers the accounts whose folders have conversations (a `projects` folder). A conversation that accounts share through a group counts once.
- With `-a`, an account that shares its history through a group brings the whole group's conversations, and ccam says whose: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- `blocks` shows 5-hour windows, and each account has its own, so name one account for it, as above.
- ccusage also reports the other coding tools it finds, such as Codex, unless you name one: `ccam ccusage claude daily` is Claude Code alone.
- While ccam runs it, ccusage reads its settings (`ccusage.json`) only from `./.ccusage` and the account folders. `--config` names another file: `ccam ccusage daily --config ~/.config/claude/ccusage.json`. As `bun x` and `npx` start in ccam's own folder, ccam passes them `./.ccusage/ccusage.json` with `--config` when that file exists and you name none, and gives a relative `--config` as a full path. ccusage then reads other relative paths, such as `--pi-path`, from ccam's folder, so give those in full.
- ccusage would read a folder whose path has a comma as two folders, so ccam refuses it and names the account; leave that account out with `-a`.
- `ccusage` is a ccam command, so an account named ccusage starts with `ccam run ccusage`; `ccam doctor` says so.

## How it works

The numbers of `ccam usage` come from two places: Anthropic's undocumented usage endpoint, `GET https://api.anthropic.com/api/oauth/usage`, called with each account's own login, and the values Claude Code hands to its status line command. All of it lives in `internal/usage`.

1. ccam loads `config.toml` and builds one `usage.Account` per account (`cli.usageCmd`): its folder `~/.local/share/ccam/profiles/<name>/`, its session markers (with Claude Code's own records in the account folder, they tell whether it is running) and whether it was imported in place.
2. ccam looks up every account at once, one goroutine per account (`Fetcher.FetchAll`). All of them share one `http.Client`, so connections to the endpoint are reused. Connecting may take 2 s, and each account gets 3 s in total. The client does not follow redirects, since a redirect would resend the token somewhere else, and reads at most 1 MB of the answer.
3. As each account finishes, its row is redrawn. An account that fails or does not answer in time falls back to saved numbers, so one slow account never holds up the others.

In the diagram, each account is looked up in its own goroutine, and the one that times out shows its saved numbers instead:

<!-- motion:usage-fetch -->

`personal` and `work` answer in time and show live numbers. `side` takes longer than 3 s, so ccam shows its saved numbers instead (`fallback()`). ccam then draws a block per account (`usage.Table()`) and picks the account named on the last line (`usage.Best()`).

For one account, ccam makes these checks in order (`Fetcher.Fetch`):

<!-- motion:usage-cache -->

1. With `--source cache`, it goes straight to `fallback()`.
2. Unless `--no-cache` is set, a live result in `~/.local/state/ccam/usage/<name>.json` measured less than 60 seconds ago (`CacheTTL`) is used, judged as of now: a window that has reset since shows 0%, and a file dated ahead of the clock does not count as fresh. Running `ccam usage` twice in a row costs one request per account.
3. ccam reads `.credentials.json` (or the macOS Keychain) without changing it and decodes only the `claudeAiOauth` key (`readCredentials()`). No subscription login means `not logged in`, or `no limits (API)` when `.claude.json` holds an API key. A login whose refresh token has expired, or that Claude Code cleared, is `` login expired → `ccam login <name>` `` and is not looked up.
4. An expired access token is not used. Without `--refresh`, the account shows its saved numbers with ``access token expired; `ccam <name>` renews it when it starts``.
5. Otherwise ccam calls the endpoint. A 401 or 403 counts as an expired login (`login expired: authentication failed (401)`), any other status or a body it cannot read counts as an error, and both fall back to saved numbers. A lookup ended by Ctrl+C says `cancelled`. A 429 means the endpoint has received too many requests for this account: ccam writes `<name>.pause` with the time to ask again (`Retry-After`, or 5 minutes, at most an hour), shows saved numbers with `rate limited by the usage endpoint; asking again after 15:10`, and does not call the endpoint for that account until then, even with `--no-cache`. Retrying immediately would only prolong the rate limit.
6. A good answer is written to `<name>.json` with the time it was measured.

For saved numbers, ccam takes the newest saved lookup (`<name>.json`, or `<name>.claude.json` from Claude Code's lookups) and, when `<name>.statusline.json` was written after it, overlays the status line's 5-hour and weekly windows on it (`fallback()`); the model windows, which the status line does not have, stay as the lookup saved them. Before it shows saved numbers, it checks each window's reset time: a window that has reset since the numbers were saved shows 0%, because the old percentage no longer holds. That is why `personal` above shows `0%` with no reset time.

The endpoint has two response formats. The newer one has a `limits` array with `session` (5-hour), `weekly_all` and `weekly_scoped` (per model) entries. The older one has `five_hour`, `seven_day`, `seven_day_opus` and `seven_day_sonnet` objects. ccam reads both (`ParseLive`): the `limits` array wins, and the old objects fill in what it does not cover. A body with no limits at all is reported as a changed format, which `ccam doctor` tells you about.

### The limit record

When Claude Code stops a turn at a usage limit, the [ccam plugin](/docs/install#claude-code-plugin)'s hook reads the refused request from the conversation: which window was full (5-hour, weekly, Opus or Sonnet) and when it resets. It writes that to `~/.local/state/ccam/usage/<name>.limit.json`. ccam applies the record on top of whatever else it knows, from every source (`Fetch`): each window in it counts as 100% until its reset, however old or low the saved numbers are, so `ccam next --auto`, `ccam pick`, `--fallback` and [`auto_switch`](/features/next#auto-switch) do not choose the account again before then.

With nothing else known, the record alone is what the account shows (`source: the usage limit Claude Code reported` with `-v`). Unattended sessions (`-p` runs and background sessions) write no record, nor does a stop set by an organization (for example, when it is out of usage or a group's limit is $0), which its admins can lift at any time. `ccam login` into another Claude account drops the record; deleting the file clears it too.

### Renewing tokens

Refresh tokens rotate: renewing one gives you a new refresh token and invalidates the old one. If ccam renewed the login while Claude Code was running as the same account, that session would be signed out. So `ccam usage` reads tokens and never writes them, unless you pass `--refresh`. Then ccam (`Fetcher.refresh`):

1. refuses accounts registered in place with `ccam import`, since plain `claude` may use that login without ccam knowing;
2. refuses logins stored in the macOS Keychain;
3. takes a lock file, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, so two ccam processes never both renew the same token; while it is held, the account shows `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)`. `ccam rename`, `remove` and `logout` take it too, and refuse while a renewal is in progress;
4. under the lock, refuses if a session of the account is running, however it was started (through ccam, plain `claude`, an IDE, or a background session and the daemon behind it, which renews the login itself), and skips the request if another process already renewed the token;
5. sends nothing while Claude Code holds its own renewal locks in the account folder (`.oauth_refresh.lock`, `<folder>.lock`): Claude Code is renewing that login right then;
6. posts the refresh token to the token endpoint and waits up to 30 seconds for the answer, the time Claude Code gives its own renewals. Once sent, a request is never interrupted, even by Ctrl+C: the server may already have rotated the refresh token, and only its response carries the new login. Ctrl+C or a hangup ends the other lookups, and ccam exits (130 after Ctrl+C) once the answer is saved;
7. saves under Claude Code's write lock (`.storage-write.lock`) and changes only the token fields of `.credentials.json` as it is then: what Claude Code wrote meanwhile (an MCP sign-in) stays, a login Claude Code renewed meanwhile is kept, and a logout is never undone. When the answer gives the new refresh token's lifetime, it is stored as `refreshTokenExpiresAt`, as Claude Code does.

A refresh that does not go through says why, on the account line or above its saved numbers: `token not renewed: a session is running`, for example, or `login may be lost: …` when a request was sent but its answer could not be used or saved, with the advice to run `ccam login <name>` if Claude Code then asks you to log in.

### The status line collector

Claude Code runs the status line command and pipes it a JSON document. ccam reads the model name and `rate_limits` from it:

<!-- motion:usage-statusline -->

1. ccam finds the account from `CCAM_PROFILE`, which ccam sets when it starts Claude Code, or otherwise by matching `CLAUDE_CONFIG_DIR` against the account folders (`runStatusline`).
2. ccam reads `rate_limits.five_hour` and `rate_limits.seven_day` (`ParseStatusLine`). The JSON has no rate limits for API key accounts, or before the session's first response; then nothing is saved.
3. ccam writes `~/.local/state/ccam/usage/<name>.statusline.json` (`SaveStatusLine`). `ccam usage` uses it when the endpoint does not answer.
4. ccam prints `[name] · model · 5h N% · 7d N%`. With `--wrap`, it runs your command through `sh -c` with the same input and prints its output the way Claude Code shows a status line: every line that is not blank, with `[name]` in front of the first. A command that fails shows the account name alone. On Windows it runs the command the way Claude Code does: with Git Bash when it is installed, otherwise with PowerShell.
5. When the 5-hour or weekly usage is at 90% or more and there is another account, the line ends with `→ ccam next --auto --continue`, for example `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

The wrapped command has no time limit of its own: Claude Code gives a status line command 10 minutes and stops a run when a newer update starts, so a slow command shows as it would without ccam. When Claude Code ends ccam, the command ends too (on Linux and macOS with every process it started).

`ccam statusline install` writes `statusLine` into each account's `settings.json` and keeps the object's other keys, such as `padding`, `refreshInterval` and `hideVimModeIndicator`. When an account already has a status line command, ccam keeps it as `--wrap '<old command>'`; on Windows as `--wrap-base64 <base64>`, which needs no quoting in either shell. If the status line already runs ccam (from another path, or quoted the way an older version wrote it), ccam replaces it instead of wrapping it, and keeps the command it wraps.

On Windows the path to ccam is written with forward slashes and `~` for the home folder (`~/.local/bin/ccam.exe statusline`), which both Git Bash and PowerShell understand. A path that needs quotes is written in full for the shell Claude Code uses: `'…'` for Git Bash, `& '…'` for PowerShell. A `settings.json` shared through a group is written once, through its link, so every account in the group gets it. Running it again prints `(already set)` for files that already use ccam.

## Safety and edge cases

- `ccam usage` never changes a token file unless you pass `--refresh`. When an access token has expired, run the account once with `ccam <name>` and Claude Code renews it. A login that has fully expired needs `ccam login <name>`.
- `--refresh` never touches running accounts (even when Claude Code was started without ccam), imported accounts or Keychain logins. A refused refresh shows the saved numbers and the reason (`token not renewed: …`).
- The usage endpoint is not publicly documented and may change. When it does, you still see saved numbers, and `ccam doctor` reports the format change.
- With no saved numbers and no answer, the account line shows `lookup failed:` and the reason, for example `no saved usage` with `--source cache`.
- The cache files hold only percentages and reset times, never tokens.

## Related

- [Usage limits and switching accounts](/docs/usage-limits)
- [Switching to the next account](/features/next)
- [What to do when Claude Code hits its usage limit](/guides/usage-limit-workflow)
- [Command reference](/docs/commands)
