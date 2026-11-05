---
title: 为脚本选择账号
description: ccam pick 为脚本输出余量最多的账号，让同时启动的并行任务各用不同的账号，并在没有账号还有余量时以 75 退出。
updated: 2026-11-05
order: 12
translation: machine
sourceHash: "0ec25752ace9"
---

`ccam pick` 输出已登录账号中余量最多的那个账号的名称，这样脚本就能用它启动 Claude Code。只要有足够多的账号还有余量，同时启动的并行任务就会各自拿到不同的账号，这样一批 `claude -p` 任务会分散到你的各个账号上，而不是耗尽其中一个。

**Team 和 Enterprise 席位。** 自动选择账号（`ccam next`、`--fallback`、`ccam pick`、仪表盘中的 `n`，以及你开启后的 [`auto_switch`](/features/next#auto-switch)）也会在同一个 Claude Team 或 Enterprise 组织的席位之间切换，这可能违反 Anthropic 的条款。在依赖这项功能之前，请先阅读[常见问题](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，并向你所在的组织和 Anthropic 确认。

## 用法

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

退出码：在标准输出上给出名称时为 `0`；没有账号还有余量或没有已登录账号时为 `75`（`EX_TEMPFAIL`）；其他错误（例如错误的选项）为 `1`。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--json` | 输出账号、它的 5 小时和每周用量、这些数字的来源，以及选中它的理由 | 关闭 |
| `--model <model>` | 同时计入该模型的每周时间窗口，例如 `opus` | 无 |
| `--max-usage <1-100>` | 优先选择用量低于此值的账号 | `auto_max_usage`，否则为 100 |
| `--strategy headroom\|expiring\|balance` | 在 5 小时用量与最低值相差不超过 3 个百分点的账号中每周用量最低的账号、5 小时时间窗口最早重置的账号，或平均每天剩余每周配额最多的账号 | `auto_strategy`，否则为 `headroom` |
| `--cross-org` | 也选择其他组织的账号 | `cross_org`，否则为关闭 |
| `--wait[=<duration>]` | 没有账号还有余量时，等待第一个账号重置，最多等待这么久。值写在 `=` 之后：`--wait 30m` 会被拒绝，并给出提示 `--wait takes its value after '=': --wait=30m` | 关闭；单独给出时为 `5h` |
| `--dry-run` | 只做选择，不记录这次选择 | 关闭 |

## 工作原理

1. ccam 列出此文件夹将要启动的账号所属组织中已登录的账号（API 密钥账号和 Console 登录没有限制，因此不包括在内；参见[组织](/features/next#auto)）；与该账号登录到同一个 Claude 账号的文件夹也会被排除。
2. 它获取一个锁（`~/.local/state/ccam/picks.lock`），所以同时运行的选择会一个接一个地进行。
3. 它以与 `ccam usage` 相同的方式、使用 60 秒缓存，一次性查询每个账号的用量。Claude Code 报告已满的时间窗口在重置前按 100% 计算（[用量限制记录](/features/usage#the-limit-record)），因此刚刚达到限制的账号即使保存的数字仍显示还有余量，也会被排除在外。
4. 先排除最近 2 分钟内被选中过的账号，无论是被 `ccam pick` 选中，还是被 [`auto_switch`](/features/next#auto-switch) 把会话移到那里（`~/.local/state/ccam/picks/<name>`；时间戳位于未来的选择不计算在内），再用 `ccam next --auto` 的[选择规则](/features/next#auto)在其余账号中选择：优先选择低于 `auto_max_usage` 的账号，按策略选择，从不选择已达 100% 的账号。如果因此跳过了规则本来会选中的账号，`--json` 的 `reason` 会说明这一点：`30% of its 5 hours used; work was picked a moment ago`。
5. 如果每个有余量的账号最近都被选中过，规则会重新在所有这些账号中选择，因此在三个账号的情况下，第四个并行任务会与其中一个共用账号。
6. 记录这次选择并释放锁，然后输出名称。

这 2 分钟的时间窗口很重要，因为新任务的用量需要一段时间才会体现在数字中：如果没有它，同时启动的五个并行任务都会把同一个账号看作余量最多的账号。

## 安全性与边界情况

- `ccam pick` 只读取用量并写入选择时间；除了 `ccam usage` 本身会做的之外，它从不为令牌续期或读取令牌。
- 选择并不会预留账号。2 分钟后，两个脚本仍可能在同一个账号上启动；可以用 `--fallback` 在该账号达到限制时恢复。
- 使用 `--wait` 等待时不持有锁，所以其他选择在此期间照常进行。
- `ccam pick` 会忽略 [`ccam folder`](/features/folder)：它只按用量选择。

## 相关内容

- [切换到下一个账号](/features/next)
- [达到限制时重新运行](/features/next#run-again-at-the-limit)
- [用量限制](/features/usage)
- [API 密钥账号](/features/api)
