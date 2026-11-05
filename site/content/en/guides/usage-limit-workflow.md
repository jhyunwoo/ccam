---
title: What to do when Claude Code hits its usage limit
description: Wait for the reset or continue the same conversation on another Claude Code account right away. How the 5-hour and weekly limits work, step by step.
updated: 2026-11-05
published: 2026-10-19
order: 3
translation: source
---

When an account reaches its usage limit (often called a rate limit), you have two choices: wait until the limit resets, or continue the same conversation now on another account you own. With ccam the second takes one command, `ccam next --auto --continue`.

## Which limits does Claude Code have?

On Pro and Max plans, Claude Code has a [5-hour session limit and a weekly limit](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Both are shared with Claude on the web and in the apps, and usage counts against both at the same time. When one runs out, Claude Code [shows a message](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) such as "You've hit your session limit · resets 3:45pm". Switching models does not help, because the limits cover all models.

## Option 1: wait for the reset

If you have one account, wait. Recent versions of Claude Code [wait in the open session](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) and continue the task on their own once the limit resets. `/usage` inside Claude Code shows how much of your plan's limits you have used.

## Option 2: continue on another account now

If you have a second account, you can continue the same conversation right away. Set it up once:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

The status line now shows the account and its 5-hour usage while you work, so a limit does not come as a surprise.

When the limit is reached, leave Claude Code (`/exit`), then:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Open the conversation in the account with the most room:

```sh
ccam next --auto --continue
```

Claude Code starts as `work` with the most recent conversation of the current folder, including everything said so far. Sharing history (`ccam group set default --add history` above) means both accounts work on the same conversation; without it, Claude Code copies the conversation into `work`, and ccam tells you how many tokens the first reply has to read again.

With the [ccam plugin](/docs/install#claude-code-plugin) in Claude Code, a session that stops at its limit shows a notification and records the full window, so `ccam next --auto` does not pick that account again before its reset. `/ccam:switch` (or `/ccam:switch work`) notes the conversation for `ccam next --continue`, so it continues exactly that conversation.

**Team and Enterprise seats.** Automatic account selection (`ccam next`, `--fallback`, `ccam pick`, the dashboard's `n` and, once you turn it on, [`auto_switch`](/features/next#auto-switch)) also switches between seats of the same Claude Team or Enterprise organization, which may violate Anthropic's terms. Read the [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) and check with your organization and Anthropic before you rely on it.

## Does ccam switch accounts automatically?

Only if you turn it on.

- By default, ccam is not running while you work: it starts Claude Code and replaces itself with it, so you decide when to switch. `ccam next --auto` picks by usage, taking the least weekly usage among the accounts within 3 points of the lowest 5-hour usage; it skips accounts that are at 100% or not logged in, and stays within one Team or Enterprise organization, or among personal plans.
- With `relay = true` in `config.toml`, ccam keeps running alongside the session instead, so `/ccam:switch` continues the conversation by itself, keeping the session's options, permission mode, effort level and agent; before it moves the conversation to an account of another organization, ccam asks you.
- With `auto_switch = true` (Linux and macOS; experimental), ccam also makes the switch itself: when Claude Code stops the session at its limit and another account of the same kind has room, ccam ends the session and continues the conversation there ([conditions and costs](/features/next#auto-switch)).

Some other tools swap the login of a running session instead; [ccam and other account switchers](/compare/account-switchers) explains the difference.

## Things to keep in mind

- Run `ccam next --continue` in the same folder you worked in; conversations are stored per project folder.
- Avoid running two accounts on the same conversation at the same time. ccam warns you if you try.
- Without shared history, the conversation is copied to the other account and the two copies go their own ways; the checkpoints `/rewind` restores code from stay with the first account.
- The first message on the new account may use more of its limit than usual, because that account has not cached the conversation yet. `ccam next --continue --brief` starts it on a short summary instead.
- Whether you may switch between accounts like this depends on the terms that apply to each account.

See [Usage limits](/features/usage) and [Switching to the next account](/features/next) for every option.
