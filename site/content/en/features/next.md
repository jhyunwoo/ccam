---
title: Switching to the next account
description: ccam next opens Claude Code as another account; --auto picks the most room, --continue carries the conversation over, and --fallback reruns -p jobs.
updated: 2026-11-05
order: 8
translation: source
---

`ccam next` opens Claude Code as another account without you naming it. Use it when an account hits its usage limit. There are four ways to switch:

- Plain `ccam next` takes the account after the one you used last, and `--auto` takes the one with the most room. `--continue` continues this folder's last conversation on the new account, whether or not the two share history.
- Inside Claude Code, the [ccam plugin](#inside-claude-code) tells you when a limit stops a session, and `/ccam:switch` hands the conversation on.
- For scripts that run `claude -p`, `ccam --fallback` makes the switch by itself (see [Run again at the limit](#run-again-at-the-limit)).
- With [`auto_switch`](#auto-switch) turned on (Linux and macOS), ccam also moves an interactive session that Claude Code stops at its limit.

**Team and Enterprise seats.** Automatic account selection (`ccam next`, `--fallback`, `ccam pick`, the dashboard's `n` and, once you turn it on, [`auto_switch`](/features/next#auto-switch)) also switches between seats of the same Claude Team or Enterprise organization, which may violate Anthropic's terms. Read the [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) and check with your organization and Anthropic before you rely on it.

## Usage

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # continue this folder's last conversation on the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next -c --brief            # a new conversation that starts from a summary of the last one
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the most weekly quota left per day
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work CCAM_PROFILE=work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam prints the switch as `→ <from> → <to>` on standard error, with the reason for an `--auto` choice, then starts Claude Code as the new account the same way `ccam <name>` does. `<from>` is the account you used last, or with `--continue` the account the conversation belongs to. Arguments after `--` go to `claude` unchanged.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--auto` | Pick the logged-in account by usage (see below), skipping accounts at 100% | off |
| `--max-usage <1-100>` | With `--auto`: prefer accounts below this usage | `auto_max_usage`, otherwise 100 |
| `--strategy headroom\|expiring\|balance` | With `--auto`: the least weekly usage within 3 points of the lowest 5-hour usage, the soonest 5-hour reset or the most weekly quota left per day | `auto_strategy`, otherwise `headroom` |
| `--cross-org` | With `--auto`: also choose accounts of another organization (see below) | `cross_org`, otherwise off |
| `--wait[=<duration>]` | With `--auto`: when no account has room, wait for the first one to reset, at most this long. The value goes after `=`: `--wait 2h` is refused with the hint `--wait takes its value after '=': --wait=2h` | off; `5h` when given alone |
| `-c`, `--continue` | Continue this folder's last conversation on the new account (see below) | off |
| `--brief` | With `--continue`: start a new conversation as the new account from ccam's summary of the last one, instead of resuming all of it ([see below](#brief)) | off |
| `-- <claude args>` | Arguments passed to `claude` after ccam's own | none |
| `--dry-run` | Supported: prints `→ <from> → <to>` and the command line it would run, without starting Claude Code or recording a launch | off |

`--auto` and `--continue` work together. With `--continue`, the `--resume` arguments are placed before your own.

Claude Code's own options, such as `--dangerously-skip-permissions` or `--model`, go after `--`. Unlike `ccam <name>`, which passes everything after the name on, `ccam next` reads its arguments itself, so without the `--` it stops with `ccam: unknown flag: --dangerously-skip-permissions`.

## How it works

Every launch through ccam writes the time to `~/.local/state/ccam/last-used/<name>`. `ccam next` reads these files, checks which accounts are logged in, and chooses.

1. ccam builds one `next.Candidate` per account from `config.toml` (`cli.candidates`): the account, its last launch time (`launcher.LastUsed`) and whether it has a working login (`creds.Has`: a Claude login in `.credentials.json`, or on macOS in the Keychain, that has not fully expired or been cleared by Claude Code). A folder that uses an API key instead (a Console login, or an account added with `ccam add --api`) has no plan limits and is never chosen.
2. Without `--auto`, ccam sorts the accounts by slot number and finds the one launched most recently (`next.Rotate`). That is the account you are switching away from. If no account has been launched yet, it is the lowest slot.
3. From there it walks forward in slot order, wrapping from the last slot back to the first, and returns the first account that is logged in.
4. ccam prepares and starts the chosen account exactly as `ccam <name>` would (`launcher.LaunchWith`): it checks the account's links, applies MCP sharing, registers the session and replaces ccam with `claude` through `syscall.Exec`.

The walk starts at the account used last and stops at the first logged-in account after it:

<!-- motion:next-rotate -->

`work` was launched most recently, so the walk starts after it. `side` has no login and is skipped, and `lab` is chosen. Had `lab` been logged out too, the walk would wrap around to `personal` (`% len` in the diagram). Rotation never returns the account you started from: with no other logged-in account, it stops with an error.

### --auto

`--auto` decides by usage instead of slot order. ccam keeps the logged-in accounts, asks for all their usage at once the same way `ccam usage` does (`Fetcher.FetchAll`), and chooses by the rules below (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

An account is judged by its 5-hour and weekly windows, and, when `--model` is among the arguments for Claude Code (after `--`), by that model's weekly window too (`opus`, `sonnet`, `haiku` or `fable` in the name). ccam sorts the accounts into three tiers (`Policy.Choose()`):

1. **Room**: every window below `auto_max_usage` (`config.toml`, default 100). The strategy chooses among them: `headroom` (default) takes, of the accounts within 3 points of the lowest 5-hour usage, the one with the least of its week used, so an account at 12% beats one at 10% when less of its week is used, but one at 14% does not; `expiring` takes the 5-hour window that resets soonest, so its remaining quota is used before it is lost; `balance` takes the account with the most weekly quota left per day until its week resets, counting each running ccam session against it, so weekly quota is not left unused and parallel sessions spread out.
2. **Over**: below 100% but at or over `auto_max_usage`. Used only when no account has room; the one `headroom` would take among them comes first.
3. **Full**: a window at 100%, or unknown 5-hour usage. Never chosen.

A window counts as 100% until its reset when Claude Code refused a request there: the [ccam plugin](#inside-claude-code)'s hook records it (`~/.local/state/ccam/usage/<name>.limit.json`), whatever the saved numbers say. So an account that has just hit its limit is not chosen again before its reset, even when its numbers are old.

Settings for all accounts go at the top of `config.toml`, above the first `[[group]]`, `[[profile]]` or `[[folder]]`. A setting below one of them belongs to that table: ccam warns about it at every launch, and commands that save the config refuse until it is moved.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organizations.** Automatic choices stay among accounts of the same kind as the one you are leaving: seats in the same Team or Enterprise organization, or personal plans (Pro, Max), so work from a company account does not continue on a private one or the other way around. ccam knows each account's organization from its login (`ccam login` records it; older logins are read from `.claude.json`); an account whose organization ccam cannot tell can be matched with any account. Two folders logged in to the same Claude account share its limits, so one never takes over from the other. `--cross-org`, or `cross_org = true` in `config.toml`, lifts the organization rule, but after [`/ccam:switch`](#inside-claude-code) a move to another organization still waits for your yes.

When nothing can be chosen, ccam exits with code 75 (`EX_TEMPFAIL`), so a script can tell "try later" from a failure. With `--wait`, it instead sleeps until the first full account's blocking window resets (30 seconds after the reset time, when the endpoint shows the new window), prints `ccam: no account has room; waiting until 15:10`, and chooses again; when no account says when it resets, it looks again after 10 minutes. If nothing frees up before the wait ends, it exits with 75.

The lookup follows the rules of `ccam usage`: numbers less than 60 seconds old are reused, each account's lookup times out after 3 seconds, an account that does not answer, or whose access token has expired, is judged by its saved numbers, and a limit record overrides them. A saved window whose reset time has passed counts as 0%, so an account whose window has reset since is not skipped. `--auto` never renews tokens. Unlike plain rotation, `--auto` can choose the account you used last if it still has the most room.

### --continue

1. ccam looks for this folder's last conversation in every account: of the `<session>.jsonl` files in `projects/<folder>/` (the folder name is the working folder, links resolved, with every character other than a letter or digit turned into `-`, as Claude Code names it), the one with the newest prompt or reply. Sessions without one, such as a session that only ran `/mcp`, `/usage` or `/status`, or an empty file, are skipped. If a session left a note (a usage limit or `/ccam:switch`, see below) and no conversation has a prompt or reply more than 2 minutes after it, that exact session is used instead; the rows Claude Code writes at `/exit` do not count.
2. The account it belongs to is the one you switch from, and the new account is chosen as above, within the same organization. When several accounts share that history, it belongs to the one whose `.claude.json` says it ran the conversation last (`lastSessionId`), otherwise to the one used last.
3. If the new account shares that history (the same group folder), ccam starts `claude --resume <session id>`, that very conversation (`--continue` would open whichever is newest in the folder).
4. Otherwise it starts `claude --resume <the conversation's file> --fork-session`: Claude Code copies the whole conversation, under a new session ID, into the new account's own history, and the original stays as it was. From then on the new account's `--continue` finds it.
5. Unless you give `--agent` after `--`, ccam also passes on the agent the conversation runs as: the `agent` setting of the folder's `.claude/settings.local.json` or `.claude/settings.json` or of the account's `settings.json`, else the agent the conversation's last `agent-setting` row names. Claude Code does not bring the agent back when it copies a conversation (step 4), and would go on as its default agent, with every tool. When the new account does not share the old one's agents, plugins and settings (the same `agents` and `plugins` folders and `settings.json`, through a group), ccam adds a line after the one about the move: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

ccam says how big the conversation is: the new account reads all of it again with its first reply. Prompt caches are never shared between organizations, and every personal plan is an organization of its own, so when the two accounts belong to different ones ccam adds `prompt caches are not shared between organizations, so the first reply reads it all again`. When the old account is signed in to MCP servers the new one is not, ccam names them: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. Sharing `history` in a group avoids the copy, keeps one conversation, and also shares the checkpoints `/rewind` restores code from (`file-history/`).

### --brief

`--continue` resumes the whole conversation, and the new account reads all of it again with its first reply. `ccam next --continue --brief` starts a new conversation instead, from a summary that ccam writes of the old one without asking a model:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

The summary holds, in this order and only where there is something to say:

1. the conversation's title;
2. the path of its transcript, for Claude to search when it needs a detail;
3. your first request, when there were more than four;
4. the summary Claude Code wrote when it last compacted the conversation;
5. your last four requests, with `!` commands and slash commands as you typed them;
6. Claude's last reply;
7. the todo list;
8. the 10 files Claude changed last;
9. the folder's `git status --short --branch`, up to 800 characters.

The todo list keeps at most 10 items: finished items are dropped first, then the steps furthest ahead, so the current step and the next ones stay. Long parts are cut in the middle, and the whole summary stays under 16,000 characters.

ccam saves it as `~/.local/state/ccam/brief/<session>.md`, readable only by you, and the next `--brief` removes summaries older than 7 days. Claude Code starts with one prompt, and attaches the file it mentions: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` A prompt you give after `--` takes the place of the last sentence (`Then: <your prompt>`), and your other options after `--` follow it.

Claude Code ends a mention at a `"` and cuts it at a `#`, so when the summary's path has either, ccam mentions it by its path relative to the working folder; when that path has one too, the prompt names the file, and Claude reads it with its Read tool, which asks first because the file is outside the working folders (in a `-p` run in Manual mode, Claude Code denies the read and Claude goes on without the summary).

`ccam --dry-run next --continue --brief` prints the summary and saves nothing. The old conversation stays as it was, and the new one knows only what the summary says, so use `--brief` when the task matters more than every detail of how you got there. Without `--continue`, ccam stops with `--brief works with --continue`; when it cannot read the conversation or finds nothing in it to summarize, it stops and says that `ccam next --continue` without `--brief` resumes it whole. The [relay](#relay) and [`auto_switch`](#auto-switch) always resume the whole conversation. Claude Code's own `--brief` option goes after `--`, like its other options.

## Inside Claude Code

The [ccam plugin](/docs/install#claude-code-plugin) adds two hooks (Claude Code 2.1.78 or later):

- When a turn stops at a usage limit, it notes the conversation and shows a desktop notification in terminals that support one (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  It knows every message Claude Code counts as a usage limit: the 5-hour, weekly, plain and model limits (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` for a model's allowance), usage credits and spend limits, and an organization that stopped the account.

  When Claude Code has written the refused request into the conversation, the hook also records the full window until its reset ([the limit record](/features/usage#the-limit-record)), except for a stop set by an organization, which its admins can lift at any time. A short server throttle or an API key's rate limit is not a usage limit and is left alone, and so are unattended sessions (`-p` runs and background sessions). A subagent or workflow agent that stops at a usage limit counts as its session's stop, even while the main conversation only waits for it: Claude Code writes that refused request into the agent's own file in the session's `subagents` folder, and the hook reads it there. When several agents stop at once, the session gets one note and one notification.

  Claude Code 2.1.234 and later also waits in the session and continues by itself after a 5-hour reset; the notification gives you the choice to switch instead, and with [`auto_switch`](#auto-switch) on (Linux and macOS), ccam switches for you when another account has room. When no other account may take over, because no other account belongs to the same Team or Enterprise organization, or none is another personal plan (with `cross_org` on, any other Claude account may), the notification ends with `You have no other personal plans, so ccam cannot move this conversation.` or `You have no other accounts in Acme, so ccam cannot move this conversation.` instead.
- `/ccam:switch [account]` notes the conversation and ends the turn before any request is made: type `/exit`, then run `ccam next --continue` to continue it on the account you named (or on the next one). In a session that has no conversation yet, it says there is nothing to continue. Only what you type runs it: from plugin 0.3.2, Claude cannot run `/ccam:switch`, and a job Claude schedules (CronCreate, ScheduleWakeup or `.claude/scheduled_tasks.json`) passes it to Claude as plain text.

  When the account you name is outside the conversation's organization (a seat of another Team or Enterprise organization, a personal plan for a seat or the other way around, or an API key account on either side), `ccam next --continue` asks first, whatever `cross_org` says: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. It asks the same before an account it picks itself when you name none, which `cross_org = true` lets it take from another organization: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. Only `y` or `yes` moves it; on Linux and macOS, keys typed before the question are dropped first. Any other answer, or no terminal to ask on, leaves the conversation and the note where they are: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (after a pick, `ccam next --continue` picks again). As with automatic choices, an account whose organization ccam cannot tell fits any.

The note lives in `~/.local/state/ccam/handoff/<account>/<session>.json`, one file per session, so two sessions of one account never overwrite each other's. It is kept for 12 hours and used once, when the next session is about to start; a launch that fails leaves it for the next try. It belongs to the folder the session started in, even after a `cd` in Claude Code's shell. `ccam rename` updates the names in it; if the account it names has since been removed, ccam says so and chooses as usual.

### Relay

With `relay = true` at the top of `config.toml`, ccam keeps running alongside interactive sessions instead of replacing itself with Claude Code, and after `/ccam:switch` it continues that session's conversation by itself (a note another session left is not acted on): the hook ends Claude Code, and ccam runs `ccam next --continue` with that session's note. Without an account named, it takes the one with the most room, counting the conversation's model, and when `cross_org = true` lets that be an account of another organization, it asks the question above before it moves the conversation there. On Windows the hook cannot end Claude Code: type `/exit`, and ccam continues it. Nor does the hook end it for an account outside the conversation's organization: it says `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`, and ccam asks the question above once you have typed `/exit`. That costs one waiting ccam process per session, which is why it is off by default.

The relayed session keeps what resuming a conversation does not restore: options such as `--add-dir`, `--mcp-config`, the tool lists, `--settings`, `--plugin-dir` and the system prompt options, the conversation's permission mode, the session's effort level and the agent the conversation runs as (see [`--continue`](#continue)). Claude Code itself restores `--model` and `--name` from the conversation. Anything else (`--worktree`, `--debug`, `--remote-control`, options ccam does not know) is named in `ccam: not carried over: …`. The switch never asks the shared-history question, and what you typed while it ran is dropped; its only question is the one before an account of another organization.

The relay continues the conversation only when Claude Code ended the way `/ccam:switch` asked it to: not when you typed `/exit` first, pressed Ctrl+C while it ended, or ccam got a hangup or another signal, and not when the terminal is gone. The note then stays, and `ccam next --continue` continues it later. A session started under the relay of an older version of ccam, or an unattended one (`-p` runs and background sessions), is not ended by the hook; it gets the `/exit` message instead. `claude remote-control` runs without a relay.

### `auto_switch`

With `auto_switch = true` at the top of `config.toml` (Linux and macOS), ccam moves a session to another account by itself when Claude Code stops it at a usage limit: the plugin's hook picks another account with room and ends Claude Code, and ccam, running alongside the session, continues the conversation there, the way the [relay](#relay) does after `/ccam:switch`. It is off by default and experimental.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam moves a session only when all of these hold:

1. ccam started it in a terminal and keeps running alongside it. `auto_switch` turns the [relay](#relay) on, so this holds for the interactive sessions you start after setting it. A session started before, or under an older ccam, gets a notification that ends ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` Unattended sessions (`-p` runs and background sessions), `claude remote-control` and Claude Code started without ccam are never moved.
2. The account has the [ccam plugin](/docs/install#claude-code-plugin) turned on and runs Claude Code 2.1.214 or later, which passes its process ID to the hook (`CLAUDE_PID`).
3. Claude Code wrote the refused request into the conversation, with the time it resets, for that account's 5-hour, weekly, Opus or Sonnet limit. A spend limit, usage credits running out, a stop set by an organization, the Fable allowance (`You've reached your Fable limit.`), a server throttle and an API key's rate limit never move a session. A subagent's or workflow agent's limit counts, even while the main conversation only waits for it; when that is a weekly Opus or Sonnet limit, the session moves only if the conversation itself runs on that model.
4. ccam knows the account's plan, from its `.claude.json` or from what `ccam login` recorded. If not, the notification ends `ccam does not know personal's plan, so auto_switch leaves this session.`
5. The session is not connected to Remote Control. As another account it would not reconnect to the app you drive it from, so it stays, and the notification says so.
6. The limit resets more than 15 minutes from now. For a shorter wait ccam leaves the session to Claude Code: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. You have not typed a prompt since the limit.
8. Another account has room. ccam chooses as [`--auto`](#auto) does (`auto_max_usage`, `auto_strategy`, the conversation's model), among accounts of the same kind as the one you leave: another seat of the same Team or Enterprise organization, or another personal plan (Pro, Max). `cross_org` and `--cross-org` do not apply here. The account must:
   - be logged in with a login Claude Code can renew;
   - not be logged in to the same Claude account;
   - not use an API key;
   - have a plan ccam knows;
   - not have usage credits on, as far as ccam's own usage lookup shows;
   - have trusted this folder in Claude Code;
   - if the conversation runs as an agent (the `agent` setting or the conversation names one), share the agents of the account you leave: its `agents` and `plugins` folders and its `settings.json` must each be the very same (through a group), or missing in both.

   When those last three leave out every account with room, the notification names them: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` or `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` A window with a [limit record](/features/usage#the-limit-record) counts as full, and an account another session moved to in the last 2 minutes comes after the others. The choice has to be made within about 6 seconds; when the usage lookups take longer, the session stays. With `usage_source = "claude"`, `auto_switch` decides from saved numbers alone, and ccam cannot tell from Claude Code's lookup or the status line whether usage credits are on, so it assumes they are off.

Then the hook ends Claude Code, and ccam, running alongside it, shows a desktop notification in terminals that support one, such as `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` It prints the first line of the example output above and runs `ccam next --continue` for that session as the chosen account, which prints its other two lines and starts Claude Code with the session's options, permission mode, effort level and agent, as the [relay](#relay) does. The new session's first message is a fixed text, the only thing ccam sends on your behalf (a move after `/ccam:switch` sends nothing):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

After a subagent's limit it is this one, because the subagents, workflows and background commands stopped with the old session:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

If that account reaches its limit too, the same happens there. The limit records keep the chain from going back to an account before its reset, so when every account is full the session stays on the last one.

The hook learns that a conversation runs as an agent from the `agent` setting, or from the row naming the agent that Claude Code writes into the conversation from time to time and when the session ends. If so far only the session's `--agent` names it, and the chosen account does not share the agents, `ccam next --continue` starts nothing and keeps the note, and a second desktop notification says so: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

When no other account has room, nothing moves: Claude Code handles the limit as usual, waiting in the session or offering its options, and the notification says which account frees up first, for example `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` For seats it says `None of your other accounts in Acme has room now`. To stop a move that has started, press Ctrl+C while Claude Code closes: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

A move has costs. Claude Code restarts, which takes a few seconds, and the background shells, dev servers and agents that Claude started stop with the old session. The new account reads the whole conversation again with its first reply, as with [`--continue`](#continue), and ccam says how big it is; a move never uses [`--brief`](#brief).

`auto_switch` never switches before Claude Code reports a limit, never watches your usage in the background, and never copies, renews or passes on a login. On Windows it is ignored, and `ccam doctor` says so; there, with `relay = true`, `/ccam:switch` and `/exit` still move a conversation to another account.

**Anthropic's terms.** On Team and Enterprise seats, `auto_switch` moves your work to another seat of the same organization. Anthropic's help center says "Each seat on a Team plan is assigned to one named individual, and only that person may use it" and "Users can't route their work through another user's account to use that user's remaining usage" ([Team plans](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); likewise for [Enterprise plans](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)).

For personal plans (Pro, Max), Anthropic's terms do not say whether your own accounts may take over from each other automatically at a limit. The [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) forbid accessing the Services "through automated or non-human means, whether through a bot, script, or otherwise" unless Anthropic explicitly permits it, and a move sends one of the fixed messages above on your behalf. The [Usage Policy](https://www.anthropic.com/legal/aup) that takes effect on November 12, 2026 lists "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" among the things it prohibits.

ccam does not claim that `auto_switch` complies with these terms. Turning it on is your call, and so is the risk: Anthropic may suspend an account, or a Team or Enterprise organization's access, when it believes the terms are broken.

While `auto_switch` is on, `ccam doctor` warns about an account without the ccam plugin and about Claude Code versions older than 2.1.214, and for seats repeats Anthropic's seat rule ([checks](/features/doctor)).

## Run again at the limit

`ccam next` is for an interactive session that has stopped. A non-interactive run, `claude -p` in a script or a scheduled job, cannot wait for you, so `--fallback` before the account name lets ccam switch for it:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` works only together with `-p` or `--print`; otherwise ccam stops and points to `ccam next --auto --continue`.
2. Piped standard input is read once, up to 10 MiB, and given to every attempt. Like Claude Code, ccam waits at most 3 seconds for it to start (except with `--input-format stream-json`); if a pipe stays open but sends nothing, the run gets no input and ccam notes this on standard error.
3. ccam starts `claude` as a child process instead of `exec`, holds its standard output back, and passes standard error through as it comes.
4. When `claude` exits with 0, or fails for another reason, the held output and the exit code are passed on unchanged. So is a run ended by a signal (an exit code above 128): it was asked to stop, and is not run again.
5. A failure counts as the usage limit when Claude Code's message says so, any of the usage-limit messages the [hook](#inside-claude-code) knows, also inside `--output-format json` or `stream-json` output, or when the account's usage shows its 5-hour, weekly or `--model` window at 100% (a limit record included).
6. Then the output of that attempt is dropped and the same command runs as another logged-in account: the one `--auto` would choose, then accounts whose usage cannot be looked up, in number order. Accounts at 100% are skipped, accounts with an API key login are never tried, and each account is tried at most once.
7. With `-c`, `--continue`, `-r` or `--resume` among the arguments, only accounts in the same history group are tried, so the conversation can be found.
8. Only accounts of the same organization are tried, as with `--auto` (`cross_org = true` lifts this).

The next account is chosen by the same rules as `--auto` (`auto_max_usage`, `auto_strategy`, the `--model` window). When no account is left, ccam prints the last attempt's output, says that no logged-in account has room left, and exits with code 75. With `--wait` before the account name (`ccam --fallback --wait=2h work -p …`), it waits for the first account to reset instead, then may try every account again; it always sleeps at least once first, so an account that has just hit its limit is not run again at once. The duration limits the total waiting time of the run, not each wait. Written with a space (`--wait 30m`), it is refused with the same hint as `ccam next`. A leading `--` before Claude Code's arguments works as for any launch (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` prints the command of the first attempt.

Interactive sessions are not switched on their own unless you turn on [`auto_switch`](#auto-switch). With `ccam statusline install`, the status line adds `→ ccam next --auto --continue` once the 5-hour or weekly usage reaches 90% ([Usage limits](/features/usage)).

## Safety and edge cases

- Accounts without a login are never chosen. With no other logged-in account, plain `ccam next` stops with `no logged-in account other than <name>`.
- With `--continue` and no conversation in the folder: `no account has a conversation in <folder> to continue`.
- When the organization rule leaves no account: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- With `--auto` and no logged-in account: `no logged-in account to pick`. When every account is full or its usage is unknown: `` no account with known usage and headroom (see `ccam usage`) ``. Both end with exit code 75.
- `--max-usage`, `--strategy`, `--wait` and `--cross-org` work only with `--auto`, and a value outside 1-100 or an unknown strategy is refused.
- A fork copies the conversation, not the account's other state: checkpoints of the original account's `file-history/` do not come along unless history is shared.
- `--auto` makes network requests even with `--dry-run`. It only reads, and never writes token files.
- The chosen account starts through the same path as `ccam <name>`, so the same checks apply, including the warning when another Claude Code session sharing the same history is already running in this folder.
- `last-used` is only written when Claude Code actually starts. `--dry-run` leaves it alone, so a dry run does not move the rotation forward.
- `--fallback` holds standard output until an attempt ends, so streaming output such as `--output-format stream-json` arrives in one piece at the end. A client that streams `stream-json` input and waits for replies before it sends more cannot use `--fallback`. Files the first attempt changed before it hit the limit stay changed; the next attempt starts the prompt again.

## Related

- [Usage limits](/features/usage)
- [Usage limits and switching accounts](/docs/usage-limits)
- [Share MCP servers, settings and history](/docs/sharing)
- [What to do when Claude Code hits its usage limit](/guides/usage-limit-workflow)
