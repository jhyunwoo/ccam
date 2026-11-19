---
title: 其他账号上的子代理
description: ccam agent 以你的另一个账号运行 Claude Code 任务的一部分并返回 Claude 的回复，让一个长任务用上多个账号的用量。
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Claude Code 自己的子代理以启动它们的会话所用的账号运行，因此消耗的是该账号的用量。`ccam agent "<task>"` 在同一个文件夹中用 `claude -p` 以你的另一个账号运行任务，并输出 Claude 的最后一条回复。从 Claude Code 会话中调用时，它就像一个消耗另一个账号用量的子代理，因此即使会话自己的账号本该用完，一个长任务也能继续进行。

**Team 和 Enterprise 席位。** `ccam agent` 会自行选择账号。在同一个 Claude Team 或 Enterprise 组织内，它会在另一个席位上运行任务，这可能违反 Anthropic 的条款。在依赖这项功能之前，请先阅读[常见问题](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，并向你所在的组织和 Anthropic 确认。

## 用法

Claude Code 通过它的 Bash 工具运行它。安装[ccam 插件](/docs/install#claude-code-plugin)后，用平常的话提出要求即可，例如“用我的其他账号来完成这个任务的各个部分”，也可以把这句话写进你的 CLAUDE.md。插件的技能会让 Claude Code 在后台启动每个代理，写出能独立理解的任务，并给修改文件的代理分配各自的部分或各自的 git worktree。

```sh
ccam agent "find why TestParse fails and fix it; reply with what you changed"
ccam agent --model opus -- "-v breaks the parser: find out why"
ccam agent --json < task.md            # the task from standard input
ccam agent --resume 7f3c9a2e-… "now add a test for it"
ccam agent --wait=1h "…"               # nobody has room: wait for a reset
ccam --dry-run agent "…"               # the account and command it would run
```

```console
$ ccam agent "summarize what internal/usage does"
ccam: agent runs as work (12% of its 5 hours used); session 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f
internal/usage looks up the 5-hour and weekly usage of every account …
ccam: agent finished as work; follow up with: ccam agent --resume 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f "…"
```

Claude 的最后一条回复写到标准输出；以 `ccam:` 开头的行写到标准错误。`--json` 则改为输出 `{account, session_id, result, is_error, accounts}`，其中 `accounts` 列出任务运行时用过的每个账号。

退出码：claude 自己的退出码（代理完成时为 `0`）；除调用方之外没有账号还有余量时为 `75`（`EX_TEMPFAIL`）；其他错误（例如代理试图启动另一个代理）为 `1`。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--model <model>` | 代理使用的模型；选择时也会计入该模型的每周时间窗口 | 账号的默认模型 |
| `--resume <id>` | 用一个后续任务继续之前某个代理的对话（即它第一行和最后一行中的会话 ID） | 关闭 |
| `--json` | 输出账号、会话、回复，以及任务运行时用过的账号 | 关闭 |
| `--max-usage <1-100>` | 优先选择用量低于此值的账号 | `auto_max_usage`，否则为 100 |
| `--strategy headroom\|expiring\|balance` | [选择规则](/features/next#auto) | `auto_strategy`，否则为 `headroom` |
| `--cross-org` | 也以其他组织的账号运行 | `cross_org`，否则为关闭 |
| `--wait[=<duration>]` | 没有账号还有余量时，等待第一个账号重置，最多等待这么久 | 关闭；单独给出时为 `5h` |
| `--dry-run` | 输出它会选择的账号对应的命令，但不运行 | 关闭 |

## 工作原理

1. **调用方被排除在外。** 调用 `ccam agent` 的账号（ccam 启动的会话所用的账号，否则为此文件夹将要启动的账号）从不运行代理，登录到同一个 Claude 账号、与它共享限制的其他文件夹也不会。除非 `--cross-org` 另有指定，代理会留在调用方所在的[组织](/features/next#auto)内。
2. **它像 `ccam pick` 一样选择。** 用量的查询方式与 [`ccam pick`](/features/pick) 相同；只要有足够多的账号还有余量，彼此相隔 2 分钟以内启动的代理就会拿到不同的账号。
3. **它以会话的权限运行，绝不更多。** 在 `claude -p` 中没有人能批准工具，因此代理会沿用调用它的会话的权限模式和推理强度。从终端运行时，适用 Claude Code 的默认模式和你设置中的权限规则。`ccam agent` 没有任何授予更多权限的选项。
4. **达到用量限制时，对话会继续进行。** 当代理的账号达到限制时，ccam 会像[插件钩子](/features/usage#the-limit-record)那样记录这次限制，并以余量最多的账号继续同一个对话：Claude Code 把对话复制过去（`--resume <file> --fork-session`），第一条消息为“Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.”每个账号只尝试一次。
5. **每个代理对话都会被记录**在 `~/.local/state/ccam/agent/<session id>.json` 中（保留 30 天）。`--resume <id>` 会在其所在文件夹中继续该对话：运行它的账号仍有余量时以该账号继续，否则以另一个账号继续。[`ccam next --continue`](/features/next) 不考虑代理对话，因此它继续的仍是你的对话。

## 安全性与边界情况

- 代理不能启动自己的代理：ccam 用 `CCAM_AGENT=1` 标记它，`ccam agent` 在那里会拒绝运行。
- 停止 `ccam agent`（Ctrl+C、信号，或 Claude Code 中断该命令）也会停止代理。
- 代理看不到调用它的对话，只能看到自己的任务。
- 同时编辑相同文件的代理会覆盖彼此的工作。请给每个代理分配任务中各自的部分，或各自的 git worktree。
- 当调用方会话移到另一个账号时（[`auto_switch`](/features/next#auto-switch)），它启动的代理也会随之停止。它们的对话仍然保留：`--resume <id>` 可以继续它们。
- `agent` 是一个命令名，因此新账号不能命名为 `agent`；已有的同名账号仍可用 `ccam run agent` 启动。

## 相关内容

- [为脚本选择账号](/features/pick)
- [切换到下一个账号](/features/next)
- [用量限制](/features/usage)
- [所有账号的后台会话](/features/agents)
