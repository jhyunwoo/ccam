---
title: Picking an account for scripts
description: ccam pick prints the account with the most room for a script, gives workers started together different accounts, and exits 75 when none has room.
updated: 2026-11-05
order: 12
translation: source
---

`ccam pick` prints the name of the logged-in account with the most room, so a script can start Claude Code with it. Workers started at the same time each get a different account as long as enough accounts have room, so a batch of `claude -p` jobs spreads over your accounts instead of draining one.

**Team and Enterprise seats.** Automatic account selection (`ccam next`, `--fallback`, `ccam pick`, the dashboard's `n` and, once you turn it on, [`auto_switch`](/features/next#auto-switch)) also switches between seats of the same Claude Team or Enterprise organization, which may violate Anthropic's terms. Read the [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) and check with your organization and Anthropic before you rely on it.

## Usage

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # most weekly quota left per day; fewer sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

Exit codes: `0` with a name on standard output, `75` (`EX_TEMPFAIL`) when no account has room or none is logged in, `1` for other errors such as a bad flag.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--json` | Print the account, its 5-hour and weekly usage, where the numbers come from and why it was chosen | off |
| `--model <model>` | Count that model's weekly window too, such as `opus` | none |
| `--max-usage <1-100>` | Prefer accounts below this usage | `auto_max_usage`, otherwise 100 |
| `--strategy headroom\|expiring\|balance` | The least weekly usage within 3 points of the lowest 5-hour usage, the soonest 5-hour reset or the most weekly quota left per day | `auto_strategy`, otherwise `headroom` |
| `--cross-org` | Also choose accounts of other organizations | `cross_org`, otherwise off |
| `--wait[=<duration>]` | When no account has room, wait for the first one to reset, at most this long. The value goes after `=`: `--wait 30m` is refused with the hint `--wait takes its value after '=': --wait=30m` | off; `5h` when given alone |
| `--dry-run` | Choose without recording the pick | off |

## How it works

1. ccam lists the logged-in accounts (API key accounts and Console logins have no limits and are left out) in the organization of the account this folder would start (see [organizations](/features/next#auto)); a folder logged in to the same Claude account as that one is left out too.
2. It takes a lock (`~/.local/state/ccam/picks.lock`), so picks running at the same time happen one after another.
3. It looks up every account's usage at once, the same way `ccam usage` does, with the 60-second cache. A window Claude Code reported full counts as 100% until its reset ([the limit record](/features/usage#the-limit-record)), so an account that has just hit its limit is excluded even if its saved numbers still show room.
4. Accounts picked within the last 2 minutes, by `ccam pick` or by [`auto_switch`](/features/next#auto-switch) moving a session there (`~/.local/state/ccam/picks/<name>`; a pick with a future timestamp does not count), are left out first, and the [choice rules](/features/next#auto) of `ccam next --auto` pick among the rest: below `auto_max_usage` first, by the strategy, never an account at 100%. If this skips the account the rules would otherwise choose, the `reason` of `--json` says so: `30% of its 5 hours used; work was picked a moment ago`.
5. If every account with room was picked recently, the rules choose among all of them again, so with three accounts, a fourth worker shares one of them.
6. The pick is recorded and the lock released, then the name is printed.

The 2-minute window matters because a new job's usage takes a while to show in the numbers: without it, five workers started together would all see the same account as the one with the most room.

## Safety and edge cases

- `ccam pick` only reads usage and writes the pick time; it never renews or reads tokens beyond what `ccam usage` does.
- A pick does not reserve the account. Two scripts may still start on the same account after 2 minutes; use `--fallback` to recover when that account hits its limit.
- Waiting with `--wait` holds no lock, so other picks go on meanwhile.
- `ccam pick` ignores [`ccam folder`](/features/folder): it chooses by usage only.

## Related

- [Switching to the next account](/features/next)
- [Run again at the limit](/features/next#run-again-at-the-limit)
- [Usage limits](/features/usage)
- [API key accounts](/features/api)
