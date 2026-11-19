---
title: Frequently asked questions
metaTitle: "ccam FAQ: multiple Claude Code accounts"
description: Answers about ccam, Claude Code accounts, usage limits, sharing, privacy and supported systems.
updated: 2026-11-19
published: 2026-10-19
order: 6
translation: source
---

## Is ccam an official Anthropic tool?

No. ccam is an unofficial open-source project and is not affiliated with or endorsed by Anthropic. It uses the `CLAUDE_CONFIG_DIR` setting that [Claude Code documents](https://code.claude.com/docs/en/env-vars) for running several accounts side by side.

## Am I allowed to use several Claude Code accounts?

That depends on the terms of each account and, for work accounts, your organization's policy: Anthropic's [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) for Pro and Max, its [Commercial Terms](https://www.anthropic.com/legal/commercial-terms) for Team and Enterprise, and its [Usage Policy](https://www.anthropic.com/legal/aup) for all. Check them before you use several accounts. ccam does not raise, merge or hide an account's limits, and every account still stops at its own. It moves your work to another of your accounts when you ask, or by itself at a limit when you turn on [`auto_switch`](/features/next#auto-switch) (Linux and macOS). Check whether that is allowed for your accounts before you use it; the decision is yours.

**Team and Enterprise seats.** Automatic account selection (`ccam next`, `--fallback`, `ccam pick`, [`ccam agent`](/features/agent), the dashboard's `n` and [`auto_switch`](/features/next#auto-switch) when you turn it on) also switches between seats of the same Claude Team or Enterprise organization. Anthropic's help center says "Each seat on a Team plan is assigned to one named individual, and only that person may use it" and "Users can't route their work through another user's account to use that user's remaining usage" ([Team plans](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); likewise for [Enterprise plans](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Using automatic selection between seats may therefore violate Anthropic's terms. Whether you use it there is your call: check with your organization and Anthropic before you rely on it. Personal plans (Pro, Max) are covered by the paragraph above.

## Does ccam get around usage limits?

Not past any account's own limit: each account keeps its limits and stops at them, and ccam does not raise, merge or hide them. ccam shows them and lets you continue on another of your accounts when one runs out; with `auto_switch` on, it does that by itself (see the next question). Between seats of one Team or Enterprise organization, continuing on another seat uses that seat's remaining usage, which Anthropic's help center rules out for another user's seat (see Team and Enterprise seats above).

## Can ccam switch accounts by itself when a session hits its limit?

Only if you turn it on. With `auto_switch = true` at the top of `~/.config/ccam/config.toml` (Linux and macOS; off by default and experimental), when Claude Code stops an interactive session (one that ccam started) at a 5-hour, weekly, Opus, Sonnet or Fable limit, ccam ends the session and continues the conversation as another of your accounts with room: another seat of the same Team or Enterprise organization, or another personal plan (Pro, Max). It never switches before the limit, never watches your usage in the background, and never moves work out of a Team or Enterprise organization or into one.

On seats it moves your work from one seat to another, which may violate Anthropic's terms (see Team and Enterprise seats above). For personal plans, Anthropic's terms do not clearly address it: the Consumer Terms forbid accessing the Services "through automated or non-human means" unless Anthropic explicitly permits it, and a move sends one message on your behalf. Turning it on is your call.

[`auto_switch`](/features/next#auto-switch) lists every condition and what a move costs.

## How do I use two Claude Code accounts on one computer?

Give each account its own configuration folder with `CLAUDE_CONFIG_DIR`. ccam does this for you: run `ccam add` for each account, then start them with `ccam personal` and `ccam work`, or with shortcuts you name yourself, such as `ccwork`. The [guide](/guides/multiple-claude-code-accounts) also shows how to do it by hand.

## How do I continue a conversation on another account?

Run `ccam next --continue` in the folder: it carries the folder's last conversation over to the next account (`--auto` for the one with the most room). Inside Claude Code, `/ccam:switch` from the [ccam plugin](/docs/install#claude-code-plugin) does the same for the session you are in. When the two accounts share conversation history (`ccam group set default --add history`, or a group of their own: `ccam group add pair --items history --members work,side`), the new account opens that very conversation; otherwise Claude Code copies it into the new account's history. Add `--brief` to start the new account on a short summary instead of the whole conversation.

## Does ccam see or send my login?

ccam never shares or copies your Claude login between accounts, and never copies the sign-ins of MCP servers. It reads an account's access token only to ask Anthropic's usage endpoint for that account's usage: for `ccam usage`, the dashboard, `ccam doctor` and automatic choice (`ccam next --auto`, `ccam pick`, `--fallback`, and `auto_switch` when a session hits its limit). It renews an expired one only with `ccam usage --refresh`; with `usage_source = "claude"` it sends no token anywhere and reads the stored login (the file, or the macOS Keychain entry) only to tell whether it still works. `ccam ext` and `ccam doctor` read the names of an account's MCP sign-ins, never their tokens. ccam has no telemetry. Apart from those lookups, the only request it makes on its own is the dashboard's check on GitHub for a newer ccam release, which carries nothing about you or your accounts.

## How do I sign in on a remote server?

Run `ccam add` or `ccam login <name>` in the SSH session as usual. Claude Code waits for the sign-in on the server, so the browser on your computer ends up on an `http://localhost:<port>/callback?code=…` page that does not load. Copy that address from the browser and paste it into the terminal: ccam passes it to Claude Code on the server, which finishes the sign-in. The address contains a one-time code, so do not share it.

## Where are my accounts stored?

In `~/.local/share/ccam/profiles/<name>/`, one folder per account. On macOS, Claude Code keeps the login itself in the Keychain, under a name derived from that folder, so `ccam rename` keeps the folder's old name there.

## What happens to my existing ~/.claude?

Nothing, unless you run `ccam import` or `ccam migrate`. Import registers the folder where it is, so plain `claude` and `ccam <name>` keep using the same account and settings. `ccam migrate` copies its settings, `CLAUDE.md`, extensions and MCP servers into the shared `default` group and only reads the folder.

## Which systems does ccam support?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine and others), macOS and Windows 10 and 11, natively or inside WSL. On Windows, sharing files such as `settings.json` needs Developer Mode; shared folders work without it, and without it the first `ccam add` creates the `default` group with MCP servers and extensions only ([sharing on Windows](/features/share#on-windows)).

## Why does an MCP server show up in only one account?

Servers added inside a session with `claude mcp add` are stored for that account only, and without `-s user` for the current folder only. ccam tells you the next time you start the account, and `ccam mcp list` and `ccam ext` show them; move the server to the shared group with `ccam mcp import <account> <server>`.

## Does ccam slow Claude Code down?

Starting Claude Code through ccam adds about 4 milliseconds, measured with `scripts/perf-check.sh` on a laptop. ccam reads one small file, checks the shared links and then replaces itself with `claude`, so by default it is not running while you work. Windows cannot replace a process, so there ccam starts Claude Code and waits for it, doing nothing else. With `relay` or `auto_switch` turned on, ccam keeps running the same way alongside each interactive session, so it can continue the conversation on another account.

## How do I update ccam?

Run `ccam update`. It downloads the latest release, checks its checksum and replaces the `ccam` file; `ccam update --check` only tells you whether a newer release exists. The dashboard also tells you when a newer release exists, and `U` installs it after asking; ccam never updates on its own. If you installed the `.deb` package, install the new `.deb` with apt instead. See [update](/docs/install#update).

## How do I uninstall ccam?

Run `ccam uninstall`. It lists what goes and asks first, then removes every account (logged out first), the shortcut commands, `~/.config/ccam`, `~/.local/share/ccam`, `~/.local/state/ccam` and the `ccam` command. An imported `~/.claude` stays, with a copy of what it shared in its groups. See [install](/docs/install#uninstall).
