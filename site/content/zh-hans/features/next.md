---
title: 切换到下一个账号
description: ccam next 以另一个账号打开 Claude Code；--auto 选择余量最多的账号，--continue 把对话带到新账号，--fallback 重新运行 -p 任务。
updated: 2026-11-01
order: 8
translation: machine
sourceHash: "79300c521748"
---

`ccam next` 无需你指定账号名，就能以另一个账号打开 Claude Code。当某个账号达到用量限制时使用它。有四种切换方式：

- 单纯的 `ccam next` 选择你上次使用的账号之后的那个账号，`--auto` 选择余量最多的账号。`--continue` 会在新账号上继续此文件夹中的上一次对话，无论两个账号是否共享历史记录。
- 在 Claude Code 中，[ccam 插件](#inside-claude-code)会在用量限制让会话停下时提醒你，`/ccam:switch` 则把对话交接出去。
- 对于运行 `claude -p` 的脚本，`ccam --fallback` 会自动完成这种切换（参见[达到限制时重新运行](#run-again-at-the-limit)）。
- 开启 [`auto_switch`](#auto-switch)（Linux 和 macOS）后，ccam 还会移动被 Claude Code 因用量限制停止的交互式会话。

**Team 和 Enterprise 席位。** 自动选择账号（`ccam next`、`--fallback`、`ccam pick`、仪表盘中的 `n`，以及你开启后的 [`auto_switch`](/features/next#auto-switch)）也会在同一个 Claude Team 或 Enterprise 组织的席位之间切换，这可能违反 Anthropic 的条款。在依赖这项功能之前，请先阅读[常见问题](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，并向你所在的组织和 Anthropic 确认。

## 用法

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

ccam 在标准错误上以 `→ <from> → <to>` 的形式打印这次切换，对 `--auto` 的选择还会附上理由，然后像 `ccam <name>` 一样以新账号启动 Claude Code。`<from>` 是你上次使用的账号；使用 `--continue` 时，则是那次对话所属的账号。`--` 之后的参数原样传给 `claude`。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--auto` | 按用量选择已登录的账号（见下文），跳过用量已达 100% 的账号 | 关闭 |
| `--max-usage <1-100>` | 与 `--auto` 一起使用：优先选择用量低于此值的账号 | `auto_max_usage`，否则为 100 |
| `--strategy headroom\|expiring\|balance` | 与 `--auto` 一起使用：5 小时用量最低的账号、5 小时时间窗口最早重置的账号，或平均每天剩余每周配额最多的账号 | `auto_strategy`，否则为 `headroom` |
| `--cross-org` | 与 `--auto` 一起使用：也选择其他组织的账号（见下文） | `cross_org`，否则为关闭 |
| `--wait[=<duration>]` | 与 `--auto` 一起使用：没有账号还有余量时，等待第一个账号重置，最多等待这么久。值写在 `=` 之后：`--wait 2h` 会被拒绝，并给出提示 `--wait takes its value after '=': --wait=2h` | 关闭；单独给出时为 `5h` |
| `-c`、`--continue` | 在新账号上继续此文件夹中的上一次对话（见下文） | 关闭 |
| `--brief` | 与 `--continue` 一起使用：以新账号从 ccam 为上一次对话写的摘要开始一个新对话，而不是恢复整个对话（[见下文](#brief)） | 关闭 |
| `-- <claude args>` | 放在 ccam 自身参数之后传给 `claude` 的参数 | 无 |
| `--dry-run` | 支持：打印 `→ <from> → <to>` 和将要执行的命令行，不启动 Claude Code，也不记录启动 | 关闭 |

`--auto` 和 `--continue` 可以一起使用。使用 `--continue` 时，`--resume` 参数会放在你自己的参数之前。

Claude Code 自己的选项，例如 `--dangerously-skip-permissions` 或 `--model`，要放在 `--` 之后。`ccam <name>` 会把名称之后的所有内容原样传下去，而 `ccam next` 会自己读取参数，所以没有 `--` 时它会以 `ccam: unknown flag: --dangerously-skip-permissions` 停止。

## 工作原理

每次通过 ccam 启动都会把时间写入 `~/.local/state/ccam/last-used/<name>`。`ccam next` 读取这些文件，检查哪些账号已登录，然后做出选择。

1. ccam 根据 `config.toml` 为每个账号构建一个 `next.Candidate`（`cli.candidates`）：账号本身、它最近一次启动的时间（`launcher.LastUsed`），以及它是否有可用的登录（`creds.Has`：`.credentials.json` 中的 Claude 登录，在 macOS 上则是钥匙串中的登录，且未完全过期、也未被 Claude Code 清除）。改用 API 密钥的文件夹（Console 登录，或用 `ccam add --api` 添加的账号）没有套餐限制，永远不会被选中。
2. 不带 `--auto` 时，ccam 按槽位号对账号排序，并找出最近一次启动的账号（`next.Rotate`）。这就是你要切换离开的账号。如果还没有任何账号启动过，则取槽位号最小的账号。
3. 从那里开始，它按槽位顺序向后查找，到最后一个槽位后回到第一个，返回第一个已登录的账号。
4. ccam 像 `ccam <name>` 一样准备并启动所选账号（`launcher.LaunchWith`）：检查账号的链接，应用 MCP 共享，注册会话，并通过 `syscall.Exec` 用 `claude` 替换 ccam。

查找从上次使用的账号开始，在它之后的第一个已登录账号处停止：

<!-- motion:next-rotate -->

`work` 是最近一次启动的，因此从它之后开始查找。`side` 没有登录，被跳过，于是选中 `lab`。如果 `lab` 也未登录，查找会绕回到 `personal`（图中的 `% len`）。轮换永远不会返回你出发时的账号：如果没有其他已登录的账号，它会以错误停止。

### --auto

`--auto` 按用量而不是槽位顺序做决定。ccam 保留已登录的账号，以与 `ccam usage` 相同的方式一次性查询它们的用量（`Fetcher.FetchAll`），然后按下面的规则选择（`cli.pickAuto`、`usage.Policy.Choose()`）：

<!-- motion:next-auto -->

账号按其 5 小时和每周时间窗口来判断；当传给 Claude Code 的参数中（`--` 之后）有 `--model` 时，还会计入该模型的每周时间窗口（名称中含有 `opus`、`sonnet`、`haiku` 或 `fable`）。ccam 把账号分为三类（`Policy.Choose()`）：

1. **有余量**：每个时间窗口都低于 `auto_max_usage`（`config.toml`，默认 100）。由策略在其中选择：`headroom`（默认）选 5 小时用量最低的账号；`expiring` 选 5 小时时间窗口最早重置的账号，让它剩余的配额在失效前被用掉；`balance` 选在每周时间窗口重置之前平均每天剩余配额最多的账号，并把每个正在运行的 ccam 会话计入其中，这样每周配额不会被浪费，并行的会话也会分散开。
2. **超出**：低于 100%，但达到或超过 `auto_max_usage`。只在没有账号有余量时使用；5 小时用量最低的排在最前。
3. **已满**：某个时间窗口已达 100%，或 5 小时用量未知。从不选择。

当 Claude Code 在某个时间窗口拒绝了请求时，该时间窗口在重置前一律按 100% 计算：[ccam 插件](#inside-claude-code)的钩子会记录它（`~/.local/state/ccam/usage/<name>.limit.json`），无论保存的数字是多少。因此刚刚达到限制的账号在重置前不会再次被选中，即使它的数字已经过时。

对所有账号生效的设置写在 `config.toml` 的顶部，位于第一个 `[[group]]`、`[[profile]]` 或 `[[folder]]` 之上。写在其中某个表下面的设置属于该表：ccam 会在每次启动时对它发出警告，而保存配置的命令在它被移走之前会拒绝执行。

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**组织。** 自动选择只会落在与你离开的账号同类的账号上：同一个 Team 或 Enterprise 组织中的席位，或者个人套餐（Pro、Max）之间，这样公司账号的工作不会在私人账号上继续，反之亦然。ccam 从每个账号的登录中得知其所属组织（`ccam login` 会记录它；较早的登录则从 `.claude.json` 读取）；ccam 无法判断组织的账号可以与任何账号匹配。登录到同一个 Claude 账号的两个文件夹共享同一份限制，所以它们永远不会互相接替。`--cross-org`，或 `config.toml` 中的 `cross_org = true`，会取消这条组织规则，但在 [`/ccam:switch`](#inside-claude-code) 之后，移到另一个组织仍要等你回答 yes。

当没有可选的账号时，ccam 以退出码 75（`EX_TEMPFAIL`）退出，这样脚本就能把"稍后再试"和失败区分开。使用 `--wait` 时，它会改为等待，直到第一个已满账号中造成阻塞的时间窗口重置（重置时间之后 30 秒，此时接口会显示新的时间窗口），打印 `ccam: no account has room; waiting until 15:10`，然后重新选择；如果没有账号说明何时重置，它会在 10 分钟后再查看。如果在等待结束前没有账号空出来，它以 75 退出。

查询遵循 `ccam usage` 的规则：60 秒内获取的数字会被复用，每个账号的查询在 3 秒后超时，没有响应或访问令牌已过期的账号按其保存的数字判断，并且用量限制记录会覆盖这些数字。重置时间已过的已保存时间窗口按 0% 计算，因此之后已经重置的账号不会被跳过。`--auto` 从不为令牌续期。与单纯的轮换不同，如果你上次使用的账号余量仍然最多，`--auto` 可能会选中它。

### --continue

1. ccam 在每个账号中查找此文件夹的上一次对话：在 `projects/<folder>/` 中的 `<session>.jsonl` 文件里（文件夹名就是解析链接后的工作文件夹，其中字母和数字以外的每个字符都换成 `-`，与 Claude Code 的命名方式一致），选出包含最新提示或回复的那一个。没有提示或回复的会话，例如只运行过 `/mcp`、`/usage` 或 `/status` 的会话，或者空文件，会被跳过。如果某个会话留下了一条记录（用量限制或 `/ccam:switch`，见下文），并且没有任何对话在它之后超过 2 分钟还有提示或回复，则改用那个会话；Claude Code 在 `/exit` 时写入的行不计算在内。
2. 对话所属的账号就是你要切换离开的账号，新账号按上述规则在同一组织内选择。当有多个账号共享那份历史记录时，它属于 `.claude.json` 表明最后运行该对话的账号（`lastSessionId`），否则属于最近使用的账号。
3. 如果新账号共享那份历史记录（同一个组文件夹），ccam 启动 `claude --resume <session id>`，也就是那次对话本身（`--continue` 会打开文件夹中最新的对话，不管是哪一次）。
4. 否则它启动 `claude --resume <the conversation's file> --fork-session`：Claude Code 以新的会话 ID 把整个对话复制到新账号自己的历史记录中，原对话保持不变。此后新账号的 `--continue` 就能找到它。
5. 除非你在 `--` 之后给出 `--agent`，ccam 还会传递对话所用的代理：文件夹的 `.claude/settings.local.json` 或 `.claude/settings.json`、或者账号的 `settings.json` 中的 `agent` 设置，没有时则是对话最后一条 `agent-setting` 行所指定的代理。Claude Code 在复制对话时（第 4 步）不会带回代理，而会以它的默认代理、带着所有工具继续对话。当新账号不共享旧账号的代理、插件和设置（即没有通过组共享相同的 `agents` 和 `plugins` 文件夹以及 `settings.json`）时，ccam 会在关于移动的那一行之后再加一行：`work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`。

ccam 会告诉你对话有多大：新账号在第一次回复时会重新读取全部内容。提示缓存从不在组织之间共享，而每个个人套餐都是一个独立的组织，所以当两个账号属于不同组织时，ccam 会加上 `prompt caches are not shared between organizations, so the first reply reads it all again`。当旧账号已登录某些 MCP 服务器、而新账号没有时，ccam 会列出它们：`work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`。在组中共享 `history` 可以避免复制，只保留一份对话，同时也共享 `/rewind` 用来恢复代码的检查点（`file-history/`）。

### --brief

`--continue` 会恢复整个对话，新账号在第一次回复时会重新读取全部内容。`ccam next --continue --brief` 则改为开始一个新对话，其起点是 ccam 不借助模型为旧对话写的一份摘要：

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

摘要按以下顺序包含这些内容，且只在有内容可写时才包含：

1. 对话的标题；
2. 对话记录文件的路径，供 Claude 在需要某个细节时搜索；
3. 你的第一个请求（当请求超过四个时）；
4. Claude Code 上一次压缩对话时写的摘要；
5. 你最近的四个请求，其中 `!` 命令和斜杠命令保持你输入时的样子；
6. Claude 的最后一次回复；
7. 待办列表；
8. Claude 最近修改的 10 个文件；
9. 该文件夹的 `git status --short --branch`，最多 800 个字符。

待办列表最多保留 10 项：先去掉已完成的项，再去掉最靠后的步骤，因此当前步骤和接下来的几步会保留下来。过长的部分会从中间截断，整个摘要保持在 16,000 个字符以内。

ccam 把它保存为 `~/.local/state/ccam/brief/<session>.md`，只有你能读取，下一次 `--brief` 会删除超过 7 天的摘要。Claude Code 以一条提示启动，并附加其中提到的文件：`Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` 你在 `--` 之后给出的提示会取代最后一句（`Then: <your prompt>`），`--` 之后的其他选项跟在它后面。

Claude Code 在 `"` 处结束一个提及，并在 `#` 处截断它，所以当摘要的路径含有其中之一时，ccam 会用相对于工作文件夹的路径来提及它；如果这个路径也含有其中之一，提示会写出文件名，由 Claude 用它的 Read 工具读取该文件，由于文件位于工作文件夹之外，Read 工具会先询问（在 Manual 模式下的 `-p` 运行中，Claude Code 会拒绝这次读取，Claude 会在没有摘要的情况下继续）。

`ccam --dry-run next --continue --brief` 会打印摘要，但不保存任何内容。旧对话保持原样，新对话只知道摘要里写的内容，所以当任务本身比你如何走到这一步的每个细节更重要时，再使用 `--brief`。没有 `--continue` 时，ccam 会以 `--brief works with --continue` 停止；当它无法读取对话或在其中找不到可以总结的内容时，它会停止，并说明不带 `--brief` 的 `ccam next --continue` 会完整恢复对话。[接力](#relay)和 [`auto_switch`](#auto-switch) 总是恢复整个对话。Claude Code 自己的 `--brief` 选项和它的其他选项一样，放在 `--` 之后。

## 在 Claude Code 中

[ccam 插件](/docs/install#claude-code-plugin)会添加两个钩子（Claude Code 2.1.78 或更高版本）：

- 当某一轮因用量限制而停止时，它会记下这次对话，并在支持桌面通知的终端（iTerm2、WezTerm、Windows Terminal、kitty、Ghostty）中显示一条通知：``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  它认识 Claude Code 视为用量限制的每一条消息：5 小时限制、每周限制、普通限制和模型限制（`You've hit your session limit`、`You've hit your limit`，以及表示模型额度的 `You've reached your …`），用量积分和支出限制，以及组织停用了该账号。

  当 Claude Code 已把被拒绝的请求写入对话时，钩子还会把已满的时间窗口一直记录到其重置为止（[用量限制记录](/features/usage#the-limit-record)），但组织设置的停用除外，因为组织管理员随时可以解除它。服务器短暂的限流或 API 密钥的速率限制不属于用量限制，不会处理；子代理的停止以及无人值守的会话（`-p` 运行和后台会话）也不会处理。

  Claude Code 2.1.234 及更高版本还会在会话中等待，并在 5 小时时间窗口重置后自动继续；这条通知让你可以选择改为切换账号；开启 [`auto_switch`](#auto-switch)（Linux 和 macOS）后，当另一个账号还有余量时，ccam 会替你切换。当没有其他账号可以接替时，也就是没有其他账号属于同一个 Team 或 Enterprise 组织，或者没有其他个人套餐（开启 `cross_org` 时，任何其他 Claude 账号都可以接替），通知会改为以 `You have no other personal plans, so ccam cannot move this conversation.` 或 `You have no other accounts in Acme, so ccam cannot move this conversation.` 结尾。
- `/ccam:switch [account]` 会记下这次对话，并在发出任何请求之前结束这一轮：输入 `/exit`，然后运行 `ccam next --continue`，在你指定的账号（或下一个账号）上继续这次对话。在还没有任何对话的会话中，它会说明没有可以继续的内容。只有你自己输入才会运行它：从插件 0.3.2 起，Claude 无法运行 `/ccam:switch`，而 Claude 安排的任务（CronCreate、ScheduleWakeup 或 `.claude/scheduled_tasks.json`）会把它作为纯文本交给 Claude。

  当你指定的账号不在对话所属的组织内（另一个 Team 或 Enterprise 组织的席位、从席位换到个人套餐或反过来，或任一方是 API 密钥账号）时，无论 `cross_org` 如何设置，`ccam next --continue` 都会先询问：`ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`。你没有指定账号时，ccam 在移到它自己挑选的账号（`cross_org = true` 时可以来自另一个组织）之前也会同样询问：`ccam: the account picked for /ccam:switch is home, a personal plan, …`。只有 `y` 或 `yes` 才会移动；在 Linux 和 macOS 上，会先丢弃提问之前输入的按键。任何其他回答，或者没有可以提问的终端，都会让对话和记录保持原样：``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work``（如果账号是 ccam 挑选的，`ccam next --continue` 会重新挑选）。与自动选择一样，ccam 无法判断组织的账号可以与任何账号匹配。

这条记录保存在 `~/.local/state/ccam/handoff/<account>/<session>.json` 中，每个会话一个文件，因此同一账号的两个会话永远不会覆盖彼此的记录。它的有效期为 12 小时，只在下一个会话即将启动时使用一次；启动失败时它会保留下来，供下次尝试。它属于会话启动时所在的文件夹，即使之后在 Claude Code 的 shell 中执行了 `cd` 也是如此。`ccam rename` 会更新其中的名称；如果它指定的账号后来被删除，ccam 会说明这一点，并照常选择。

### 接力

在 `config.toml` 的顶部设置 `relay = true` 时，ccam 会在交互式会话旁边持续运行，而不是用 Claude Code 替换自身；在 `/ccam:switch` 之后，它会自己继续那个会话的对话（不会处理其他会话留下的记录）：钩子结束 Claude Code，ccam 用那个会话的记录运行 `ccam next --continue`。没有指定账号时，它会选择剩余额度最多的账号，并计入对话所用的模型；如果因为 `cross_org = true` 而选中了另一个组织的账号，它会在把对话移过去之前先问上面的问题。在 Windows 上，钩子无法结束 Claude Code：输入 `/exit`，ccam 就会继续这个对话。对于对话所属组织之外的账号，钩子也不会结束 Claude Code：它会显示 `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`，等你输入 `/exit` 之后，ccam 会问上面的问题。代价是每个会话多一个等待中的 ccam 进程，所以默认关闭。

接力的会话会保留恢复对话时不会还原的内容：`--add-dir`、`--mcp-config`、工具列表、`--settings`、`--plugin-dir` 和系统提示相关的选项，以及对话的权限模式、会话的推理强度和对话所用的代理（见 [`--continue`](#continue)）。`--model` 和 `--name` 由 Claude Code 自己从对话中还原。其他内容（`--worktree`、`--debug`、`--remote-control`、ccam 不认识的选项）会在 `ccam: not carried over: …` 中列出。切换时从不询问共享历史的问题，切换过程中你输入的内容会被丢弃；它唯一会问的，是移到另一个组织的账号之前的那个问题。

只有当 Claude Code 按 `/ccam:switch` 的要求结束时，接力才会继续这个对话：如果你先输入了 `/exit`、在它结束时按了 Ctrl+C，或者 ccam 收到挂断或其他信号，或者终端已经不在了，都不会继续。此时记录会保留，之后可以用 `ccam next --continue` 继续这次对话。在旧版本 ccam 的接力下启动的会话，或者无人值守的会话（`-p` 运行和后台会话），不会被钩子结束，而是收到 `/exit` 的提示。`claude remote-control` 不使用接力。

### `auto_switch`

在 `config.toml` 的顶部设置 `auto_switch = true`（Linux 和 macOS）后，当 Claude Code 因用量限制停止一个会话时，ccam 会自己把这个会话移到另一个账号：插件的钩子选出另一个还有余量的账号并结束 Claude Code，而在会话旁边持续运行的 ccam 会在那个账号上继续这个对话，就像[接力](#relay)在 `/ccam:switch` 之后所做的那样。它默认关闭，属于实验性功能。

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

只有在以下条件全部满足时，ccam 才会移动会话：

1. 会话由 ccam 在终端中启动，并且 ccam 在它旁边持续运行。`auto_switch` 会开启[接力](#relay)，所以对设置之后启动的交互式会话来说这一条成立。在此之前启动的会话，或在旧版 ccam 下启动的会话，会收到一条以 ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` 结尾的通知。无人值守的会话（`-p` 运行和后台会话）、`claude remote-control` 以及不通过 ccam 启动的 Claude Code 永远不会被移动。
2. 该账号已开启 [ccam 插件](/docs/install#claude-code-plugin)，并运行 Claude Code 2.1.214 或更高版本，这些版本会把自己的进程 ID 传给钩子（`CLAUDE_PID`）。
3. Claude Code 已把被拒绝的请求连同重置时间写入对话，且针对的是该账号的 5 小时、每周、Opus 或 Sonnet 限制。支出限制、用量积分耗尽、组织设置的停用、Fable 额度（`You've reached your Fable limit.`）、服务器限流、API 密钥的速率限制以及子代理的停止，都永远不会让会话移动。
4. ccam 知道该账号的套餐，来自它的 `.claude.json` 或 `ccam login` 记录的信息。否则通知会以 `ccam does not know personal's plan, so auto_switch leaves this session.` 结尾。
5. 会话没有连接到 Remote Control。换成另一个账号后，它不会重新连接到你用来操控它的应用，所以会话保持不动，通知会说明这一点。
6. 限制在 15 分钟以后才重置。等待时间更短时，ccam 会把会话留给 Claude Code：`ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. 自达到限制以来，你还没有输入过提示。
8. 另一个账号还有余量。ccam 的选择方式与 [`--auto`](#auto) 相同（`auto_max_usage`、`auto_strategy`、对话所用的模型），并且只在与你离开的账号同类的账号中选择：同一个 Team 或 Enterprise 组织的另一个席位，或另一个个人套餐（Pro、Max）。`cross_org` 和 `--cross-org` 在这里不适用。该账号必须：
   - 已经登录，且其登录可由 Claude Code 续期；
   - 不与当前账号登录同一个 Claude 账号；
   - 不使用 API 密钥；
   - 有 ccam 知道的套餐；
   - 就 ccam 自己的用量查询所见，没有开启用量积分；
   - 已在 Claude Code 中信任此文件夹；
   - 如果对话以某个代理运行（`agent` 设置或对话指定了代理），要共享你离开的账号的代理：它的 `agents` 和 `plugins` 文件夹以及 `settings.json` 必须各自（通过组）完全相同，或者两边都没有。

   当最后这三条排除了所有还有余量的账号时，通知会指出它们：`side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`、`side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` 或 `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` 有[用量限制记录](/features/usage#the-limit-record)的时间窗口按已满计算，最近 2 分钟内有其他会话移入的账号排在其他账号之后。选择必须在大约 6 秒内完成；用量查询耗时更长时，会话保持不动。设置 `usage_source = "claude"` 时，`auto_switch` 只根据保存的数字做决定，而且 ccam 无法从 Claude Code 的查询或状态栏中得知用量积分是否开启，所以假定它们是关闭的。

然后钩子结束 Claude Code，在它旁边持续运行的 ccam 会在支持桌面通知的终端中显示一条通知，例如 `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` 它打印上面示例输出的第一行，并以选中的账号为该会话运行 `ccam next --continue`，后者打印输出中的另外两行，并像[接力](#relay)那样，以会话的选项、权限模式、推理强度和代理启动 Claude Code。新会话的第一条消息是下面这段固定文本，这是 ccam 唯一代表你发送的内容（在 `/ccam:switch` 之后的移动不发送任何内容）：

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

如果那个账号也达到了限制，同样的事会在那里再发生一次。用量限制记录会阻止这条链在某个账号重置之前回到该账号，所以当每个账号都已用满时，会话会停在最后一个账号上。

钩子通过 `agent` 设置，或者通过 Claude Code 不时以及在会话结束时写入对话、指明代理的那一行，得知对话以某个代理运行。如果到目前为止只有会话的 `--agent` 指明了代理，而选中的账号不共享这些代理，`ccam next --continue` 就不启动任何东西，并保留记录，第二条桌面通知会说明这一点：``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``。

当没有其他账号还有余量时，什么都不会移动：Claude Code 照常处理这次限制，在会话中等待或给出它的选项，通知会说明哪个账号最先恢复，例如 `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` 对于席位，它会说 `None of your other accounts in Acme has room now`。要停止一次已经开始的移动，请在 Claude Code 关闭时按 Ctrl+C：``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

一次移动是有代价的。Claude Code 会重新启动，这需要几秒钟，Claude 启动的后台 shell、开发服务器和代理会随旧会话一起停止。新账号在第一次回复时会重新读取整个对话，与 [`--continue`](#continue) 相同，ccam 会告诉你对话有多大；移动从不使用 [`--brief`](#brief)。

`auto_switch` 从不在 Claude Code 报告限制之前切换，从不在后台监视你的用量，也从不复制、续期或转交登录信息。在 Windows 上它会被忽略，`ccam doctor` 会说明这一点；在那里，设置 `relay = true` 后，`/ccam:switch` 和 `/exit` 仍能把对话移到另一个账号。

**Anthropic 的条款。** 在 Team 和 Enterprise 席位上，`auto_switch` 会把你的工作移到同一组织的另一个席位。Anthropic 帮助中心关于 [Team 套餐](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)的说明写道："Each seat on a Team plan is assigned to one named individual, and only that person may use it"（Team 套餐的每个席位都分配给一位指定的个人，只有此人可以使用），以及 "Users can't route their work through another user's account to use that user's remaining usage"（用户不能通过另一位用户的账号完成自己的工作，以使用该用户剩余的用量）；关于 [Enterprise 套餐](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)的说明也是如此。

对于个人套餐（Pro、Max），Anthropic 的条款没有说明你自己的多个账号能否在达到限制时自动相互接替。[消费者条款](https://www.anthropic.com/legal/consumer-terms)（Consumer Terms）禁止 "through automated or non-human means, whether through a bot, script, or otherwise"（通过自动化或非人工的方式，无论是通过机器人、脚本还是其他方式）访问服务，除非 Anthropic 明确允许，而一次移动会代表你发送上面那段固定消息。2026 年 11 月 12 日生效的[使用政策](https://www.anthropic.com/legal/aup)（Usage Policy）把 "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products"（故意绕过或干扰我们产品中设定的功能、限制或防护措施）列为禁止事项之一。

ccam 不声称 `auto_switch` 符合这些条款。是否开启由你自己决定，风险也由你承担：当 Anthropic 认为条款被违反时，它可能暂停一个账号，或暂停一个 Team 或 Enterprise 组织的访问权限。

开启 `auto_switch` 时，`ccam doctor` 会对没有 ccam 插件的账号以及早于 2.1.214 的 Claude Code 版本发出警告，并对席位重申 Anthropic 关于席位的规定（[检查项](/features/doctor)）。

## 达到限制时重新运行

`ccam next` 用于已经停下来的交互式会话。非交互式运行，即脚本或定时任务中的 `claude -p`，无法等你操作，所以在账号名之前加上 `--fallback`，就能让 ccam 替它切换：

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

1. `--fallback` 只能与 `-p` 或 `--print` 一起使用；否则 ccam 会停止，并提示使用 `ccam next --auto --continue`。
2. 通过管道传入的标准输入只读取一次，最多 10 MiB，并提供给每一次尝试。与 Claude Code 一样，ccam 最多等待 3 秒让输入开始（使用 `--input-format stream-json` 时除外）；如果管道一直打开却不发送任何内容，这次运行就不会得到输入，ccam 会在标准错误中说明这一点。
3. ccam 以子进程方式启动 `claude`，而不是使用 `exec`；它会暂存标准输出，标准错误则随产生随时传出。
4. 当 `claude` 以 0 退出，或因其他原因失败时，暂存的输出和退出码会原样传递。被信号结束的运行（退出码大于 128）也是如此：它是被要求停止的，不会再次运行。
5. 当 Claude Code 的消息表明达到了用量限制（[钩子](#inside-claude-code)认识的任何一条用量限制消息，也包括 `--output-format json` 或 `stream-json` 输出中的消息），或者该账号的用量显示其 5 小时、每周或 `--model` 时间窗口已达 100%（包括用量限制记录）时，这次失败就被视为达到用量限制。
6. 此时会丢弃这次尝试的输出，并以另一个已登录账号运行同一个命令：先选 5 小时用量最低的账号，然后是无法查询用量的账号，按编号顺序。用量已达 100% 的账号会被跳过，使用 API 密钥登录的账号永远不会被尝试，每个账号最多尝试一次。
7. 如果参数中有 `-c`、`--continue`、`-r` 或 `--resume`，只会尝试同一历史组中的账号，以便找到那次对话。
8. 与 `--auto` 一样，只会尝试同一组织的账号（`cross_org = true` 会取消这一限制）。

下一个账号按与 `--auto` 相同的规则选择（`auto_max_usage`、`auto_strategy`、`--model` 的时间窗口）。当没有剩余账号时，ccam 会打印最后一次尝试的输出，说明没有已登录账号还有余量，并以退出码 75 退出。在账号名之前加上 `--wait`（`ccam --fallback --wait=2h work -p …`）时，它会改为等待第一个账号重置，然后可能再次尝试每个账号；它总是先至少等待一次，所以刚刚达到限制的账号不会立刻再次运行。这个时长限制的是整次运行的总等待时间，而不是每一次等待。用空格分隔（`--wait 30m`）时会被拒绝，并给出与 `ccam next` 相同的提示。Claude Code 参数之前的 `--` 与任何启动一样有效（`ccam --fallback work -- -p …`）。`ccam --dry-run --fallback work -p …` 会打印第一次尝试的命令。

除非你开启 [`auto_switch`](#auto-switch)，交互式会话不会自动切换。使用 `ccam statusline install` 后，当 5 小时或每周用量达到 90% 时，状态栏会加上 `→ ccam next --auto --continue`（[用量限制](/features/usage)）。

## 安全性与边界情况

- 没有登录的账号永远不会被选中。如果没有其他已登录的账号，单纯的 `ccam next` 会以 `no logged-in account other than <name>` 停止。
- 使用 `--continue` 但文件夹中没有对话时：`no account has a conversation in <folder> to continue`。
- 当组织规则排除了所有账号时：`` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``。
- 使用 `--auto` 但没有已登录账号时：`no logged-in account to pick`。当每个账号都已用满或用量未知时：`` no account with known usage and headroom (see `ccam usage`) ``。两者都以退出码 75 结束。
- `--max-usage`、`--strategy`、`--wait` 和 `--cross-org` 只能与 `--auto` 一起使用，超出 1-100 的值或未知的策略会被拒绝。
- 复刻（fork）只复制对话，不复制账号的其他状态：除非共享历史记录，原账号 `file-history/` 中的检查点不会一起带过去。
- 即使带有 `--dry-run`，`--auto` 也会发出网络请求。它只读取，从不写入令牌文件。
- 所选账号通过与 `ccam <name>` 相同的路径启动，因此适用相同的检查，包括当另一个共享同一历史记录的 Claude Code 会话已在此文件夹中运行时发出的警告。
- `last-used` 只在 Claude Code 实际启动时写入。`--dry-run` 不会改动它，因此试运行不会推进轮换。
- `--fallback` 会暂存标准输出直到一次尝试结束，所以 `--output-format stream-json` 这样的流式输出会在最后一次性到达。以流式方式发送 `stream-json` 输入、并在发送更多内容前等待回复的客户端无法使用 `--fallback`。第一次尝试在达到限制之前修改过的文件会保持修改后的状态；下一次尝试会重新开始执行提示。

## 相关内容

- [用量限制](/features/usage)
- [用量限制与切换账号](/docs/usage-limits)
- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [Claude Code 达到用量限制时该怎么办](/guides/usage-limit-workflow)
