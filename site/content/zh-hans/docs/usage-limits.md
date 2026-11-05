---
title: 用量限制与切换账号
description: 用 ccam usage 查看每个 Claude Code 账号的 5 小时用量和每周用量，并用 ccam next --auto 打开余量最多的账号。
updated: 2026-11-05
order: 5
translation: machine
sourceHash: "06bd349073e3"
---

`ccam usage` 显示每个账号的各项限额已用掉多少：5 小时限额、每周限额，以及各模型自己的每周限额（例如 Fable）。`ccam next --auto` 打开余量最多的账号，`--continue` 会把当前对话带过去。

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

每个账号先占一行，显示这些数字是多久之前获取的，然后每项限额各有一个用量条，并附上距离重置的时间：5 小时窗口、所有模型共用的每周窗口，以及各模型自己的每周窗口，例如 Fable。用量条在 80% 时变为黄色，在 100% 时变为红色。最后一行列出 `ccam next --auto` 会打开的账号；设置 `auto_strategy = "expiring"` 或 `"balance"` 时，这一行改为 `→ next --auto picks <name>: <reason>`。

**Team 和 Enterprise 席位。** 自动选择账号（`ccam next`、`--fallback`、`ccam pick`、仪表盘中的 `n`，以及你开启后的 [`auto_switch`](/features/next#auto-switch)）也会在同一个 Claude Team 或 Enterprise 组织的席位之间切换，这可能违反 Anthropic 的条款。在依赖这项功能之前，请先阅读[常见问题](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，并向你所在的组织和 Anthropic 确认。

## 数字从哪里来

ccam 使用每个账号自己的登录，同时向 Anthropic 的用量接口查询所有账号。每个账号最多等待 3 秒，某个账号一返回结果就立即显示它。

这个接口没有公开文档，可能会发生变化。请求失败时，ccam 会显示最后一次获得的数字，并将其标记为缓存。在这些数字保存之后已经重置的时间窗口会显示 0%。如果响应格式发生了变化，`ccam doctor` 会告诉你。

当 Claude Code 因用量限制停止某个会话时，[ccam 插件](/docs/install#claude-code-plugin)会记录哪个时间窗口已满以及它何时重置。在那之前，无论其他数字是多少，该时间窗口都显示 100%，所以 `ccam next --auto` 和 `--fallback` 不会过早地再次选中这个账号。

60 秒内获取的数字会被复用。选项：

| 选项 | 作用 |
| --- | --- |
| `--watch[=seconds]` | 持续刷新，默认每 60 秒一次 |
| `--json` | 机器可读的输出 |
| `-v` | 同时显示每个账号的数字来自哪里 |
| `--no-cache` | 总是重新查询 |
| `--source cache` | 从不查询；只显示已保存的数字 |

## 过期的登录

当某个账号的访问令牌已过期时，ccam 不会为它续期，因为如果某个 Claude Code 会话正在使用同一登录，续期可能会让该会话退出登录。该账号会显示 ``access token expired; `ccam <name>` renews it when it starts``：用 `ccam <name>` 启动一次该账号，Claude Code 就会为它续期。登录已完全过期，或者已被 Claude Code 清除时，会显示 `` login expired → `ccam login <name>` ``，需要重新登录。

`ccam usage --refresh` 会自行为过期的令牌续期，但仅限于没有正在运行的 Claude Code 会话的账号，无论该会话是怎样启动的。它从不为存储在 macOS 钥匙串中的登录或已导入的 `~/.claude` 中的登录续期，并且一次续期最多可能需要 30 秒。

## 在状态栏中显示用量

```sh
ccam statusline install
```

之后 Claude Code 会在状态栏中显示账号名及其用量，例如 `[work] · Opus · 5h 24% · 7d 41%`。用量达到 90% 起，状态栏还会提示 `ccam next --auto --continue`。如果你已经有状态栏命令，ccam 会保留它，并在前面加上账号名。状态栏也会保存它的数字，当接口没有响应时，`ccam usage` 会使用这些数字。

## 切换

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the most room
ccam next --continue  # continue this folder's conversation on the next account
```

`ccam next` 会跳过未登录的账号，`--auto` 会跳过用量已达 100% 的账号（也包括 Claude Code 报告已满的时间窗口，直到其重置）以及属于其他组织的账号。如果两个账号共享对话历史，`--continue` 会打开同一个对话；否则 Claude Code 会把它复制到新账号中。参见[共享](/docs/sharing)和[切换到下一个账号](/features/next#continue)。

在 `config.toml` 的顶部设置 `auto_switch = true`（仅限 Linux 和 macOS；默认关闭）时，如果 Claude Code 因用量限制停止了一个交互式会话，并且另一个同类账号（同一组织中的席位，或另一个个人套餐）还有余量，ccam 会自己完成这次切换。它何时动作、哪些事绝不会做，请参见 [`auto_switch`](/features/next#auto-switch)。

对于运行 `claude -p` 的脚本和定时任务，请把 `--fallback` 放在账号名之前。当该账号达到限制时，ccam 会以余量最多的账号重新运行同一个命令，并且只打印最终完成的那次运行的输出：

```sh
ccam --fallback work -p "review the diff" --output-format json
```

ccam 如何做出判断，参见[达到限制时重新运行](/features/next#run-again-at-the-limit)。

[用量限制工作流程指南](/guides/usage-limit-workflow)通过一个完整示例进行了讲解。
