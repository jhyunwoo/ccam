---
title: 仪表盘
description: 在终端中运行 ccam 会打开全屏仪表盘，显示账号、组、文件夹规则、扩展和检查结果；每项更改都会先显示计划。
updated: 2026-10-30
order: 15
translation: machine
sourceHash: "d65f1b5bb0c8"
---

`ccam` 不带参数运行时，如果标准输入和标准输出都是终端，就会打开全屏仪表盘；否则打印 `ccam list` 的输出。仪表盘显示每个账号及其用量和组，并能完成各个命令所做的事：启动、添加、登录、更改组、设置文件夹规则、列出后台会话、检查和修复。启动 Claude Code 或连接到后台会话时仪表盘会关闭；添加账号或登录时会暂时把终端交出去，完成后再回到仪表盘。

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · r to refresh · a to add · ? for shortcuts
```

选项卡名称下面的一行用一句话说明当前选项卡的用途。终端较宽时，Accounts 选项卡还会显示每个账号的邮箱。在每个选项卡上，选项卡所在行的末尾都会显示 [`auto_switch`](/features/next#auto-switch) 是否开启（Linux 和 macOS；在 Windows 上会被忽略）。开启时，Health 选项卡会列出这项设置还需要什么：Claude Code 2.1.214 或更高版本，以及在每个使用 Claude 订阅登录的账号中开启 ccam 插件。

仪表盘每次打开时都会在后台向 GitHub 询问是否有更新的 ccam 版本；第一帧从不等待它，你离线时什么也不显示。有新版本时，标题下面的一行会推荐它：

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` 会先询问，然后按 [`ccam update`](/features/update) 的方式更新：用 SHA-256 校验和检查该版本，并一步替换 ccam。你当前所在的仪表盘仍是旧版本，所以请退出并重新启动 `ccam`。如果无法在这里更新，例如通过 `.deb` 软件包安装的 ccam 副本，仪表盘会显示与 `ccam update` 相同的建议，并保留这个更新提示。

第一帧只根据 `config.toml` 绘制，用时远低于 100 ms；用量会在每个账号返回结果后陆续填入。每项限额都有自己的用量条：5 小时窗口、所有模型共用的每周窗口，以及 Fable 自己的每周窗口；`-` 表示该账号没有这个窗口。百分比在 80% 时变为黄色，在 100% 时变为红色。如果查询失败，缓存的值会以淡色显示。Claude Code 报告已满的时间窗口在重置前显示 100%。访问令牌过期时会以淡色显示 `⚠ token expired`（Claude Code 会在下次启动时续期），而 `⚠ login expired` 表示该账号需要运行 `ccam login`。

## 选项卡

| 选项卡 | 显示内容 | 按键 |
| --- | --- | --- |
| Accounts | 编号、名称、邮箱、5 小时、每周和 Fable 用量条、上次使用时间、组；`✔` 表示默认账号，`▸` 表示当前文件夹的账号，`•` 表示正在运行（也包括不通过 ccam 启动的 Claude Code：直接运行的 `claude`、IDE、后台会话及其守护进程） | 见下文 |
| Agents | [`ccam agents`](/features/agents) 列出的内容：每个账号的后台会话和打开的会话，以及它们的状态、文件夹和启动时间（见下文） | `enter` 连接（同 `ccam attach`），`x` 停止，`r` 重新加载 |
| Groups | 每个组、它共享的内容和它的账号，以及不在任何组中的账号 | `enter` 编辑，`n` 新建，`e` 重命名，`x` 删除 |
| Folders | [`ccam folder`](/features/folder) 的规则，以及在当前位置生效的那一条，即使通过符号链接进入该文件夹也是如此（新规则会以解析链接后的路径保存） | `n` 在此处使用某个账号，`enter` 更改，`x` 删除 |
| Extensions | [`ccam ext`](/features/ext) 显示的内容，每个扩展只列一次，并注明它来自哪些组；`off:` 列出未启用某个插件的账号；账号自己的服务器中只对当前文件夹有效的那些；`This folder's project`；以及不属于任何账号的问题 | `r` 重新加载 |
| Health | [`ccam doctor`](/features/doctor) 报告的内容，每个账号一个框（见下文） | `f` 修复选中的框，`F` 修复全部（`doctor --fix`），`↑`/`↓` 在框之间移动 |

`←`/`→` 或 `Tab` 切换选项卡，`?` 列出当前选项卡的按键，`r` 绕过短期缓存刷新用量，并重新加载 Agents、Extensions 和 Health 选项卡，`q` 或 `Esc` 退出。

## Accounts 选项卡中的按键

| 按键 | 作用 |
| --- | --- |
| `Enter` | 以所选账号启动 Claude Code |
| `1`–`9` | 移到该编号的账号；快速输入 `1` `2` 可移到 12 号。只有 `Enter` 才会启动 |
| `↑`/`↓`、`k`/`j` | 移动 |
| `n` | 启动 [`ccam next --auto`](/features/next) 会选择的账号 |
| `o` | 带选项启动：`--continue`、`--resume`、其他参数，或试运行 |
| `u` | 用量详情：5 小时、每周和各模型的时间窗口，以及它们的重置时间和数字的来源 |
| `a` | 添加账号（见下文的向导） |
| `l` | 重新登录 |
| `s` | 选择账号所属的组 |
| `d` | 设为默认账号 |
| `f` | 在当前文件夹中使用它 |
| `e` | 重命名 |
| `c` | 设置或删除它的快捷命令 |
| `x` | 删除它（先询问，默认选中 No） |

**Team 和 Enterprise 席位。** 自动选择账号（`ccam next`、`--fallback`、`ccam pick`、仪表盘中的 `n`，以及你开启后的 [`auto_switch`](/features/next#auto-switch)）也会在同一个 Claude Team 或 Enterprise 组织的席位之间切换，这可能违反 Anthropic 的条款。在依赖这项功能之前，请先阅读[常见问题](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，并向你所在的组织和 Anthropic 确认。

## Health 选项卡

打开 Health 选项卡时，它会运行 [`ccam doctor`](/features/doctor) 的检查并绘制相同的视图：第一行为每个账号显示一个标记，下面每个账号一个框。框内每一行是一个检查过的方面，`Folder`、`Login` 和 `Share links` 在前，标记为 `✔` 正常、`ℹ` 信息、`⚠` 警告或 `✘` 错误，让你一眼看出哪个账号在哪里有问题。不属于任何账号的结果放在 `General` 框中。

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` 在框之间移动，选中的框带有强调色边框。`f` 只修复选中的框，方式与 `ccam doctor --fix` 相同：在账号的框上只修复该账号，在 `General` 框上修复不属于任何账号的结果（插件路径和快捷命令）；其他账号只会重新检查。它会先询问（``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``），没有可修复内容的框只显示 `nothing to repair in <name>`。`F` 像 `ccam doctor --fix` 一样一次修复全部；`r` 重新检查。

## Agents 选项卡

Agents 选项卡显示 [`ccam agents`](/features/agents) 打印的内容：每个账号的后台会话和打开的会话，以及各自所属的账号。在你第一次打开该选项卡时、按 `r` 时、停止一个会话后以及配置更改后，它会同时询问每个账号的 Claude Code（`claude agents --json`，每个超时为 15 秒），在所有账号都回答之前显示 `Asking every account's Claude Code for its sessions…`。

```console
  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Background agents and open sessions of every account; enter attaches to one.

    Account   ID        Name   State          Folder      Started
 ❯  work      -         -      open (busy)    ~/code/web  now
    personal  49ec0028  probe  blocked, idle  ~/code/app  12m

 ✘ lab: could not list its sessions: exit status 1 'claude agents --json' is
   not available in this environment.

 enter to attach · x to stop · r to reload · ? for shortcuts
```

`State` 是后台会话的状态及其进程正在做的事，交互式会话则显示 `open`，与 `ccam agents` 相同。如果无法查询某个账号，列表下方会有一行 `✘`，显示 Claude Code 的回答，例如上面的 `lab`，它的代理视图已关闭。已结束的后台会话不显示（`ccam agents --all` 会列出它们）；当所有账号都已回答且都没有会话时，选项卡会显示 ``No sessions. Start one in the background with `ccam <account> --bg "…"`.``。在较窄的终端中，选项卡会依次去掉 `Started`、`Name`、`Folder` 列，并从左侧截短过长的文件夹路径。

在后台会话上按 `Enter` 会关闭仪表盘，并以该会话所属的账号在当前终端中打开它，就像 `ccam attach <id>` 那样；仪表盘不会再回来（在 Claude Code 中，`←` 进入它的代理视图，`Ctrl+Z` 回到 shell）。打开的会话没有 ID，所以对它按 `Enter` 和 `x` 只会提示需要一个后台会话。`x` 会先询问，默认选中 `No, keep it`：``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.``。选择 yes 后，ccam 会以该账号运行 `claude stop <id>`，显示 ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` 或无法停止的原因，并重新读取列表。

## 更改组

在 Groups 选项卡中，`Enter` 会打开编辑器：左边是项，右边是账号，用 `Tab` 切换，用 `Space` 勾选或取消。如果某个账号已经从另一个组获得某一项，该项会变灰并显示原因；会重复获得某一项的账号也同样变灰，因为一项只能来自一个组。`Enter` 显示计划：

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

如果某个文件与组中的副本不同，ccam 会先询问如何处理（保留组中的、使用我的，或者合并）。然后更改在后台运行，使用的代码与 [`ccam group`](/features/share) 相同：先保存配置，再按顺序移动每个账号的数据，原始文件放入 `.ccam-backup/`；退出 MCP 共享的账号会在它的 `.claude.json` 中保留组的服务器。

会移动正在运行的账号数据的更改会被拒绝，并显示该账号的名称，即使该账号是在计划打开期间才启动的也是如此；如果另一个 ccam 进程更改了组而让计划过时，计划会被拒绝，并显示 `the groups changed while this plan waited; nothing was moved, make the change again`。当另一个 ccam 进程正在更改组时，新的更改、重命名或删除会显示 `another ccam process is changing the groups right now`；组更改中途停止后，它们会等到 Health 选项卡中的 `F`（或 `ccam doctor --fix`）完成它为止。

## 添加账号

`a` 会打开 `ccam add` 在终端中同样使用的向导：

1. 账号名称。
2. 登录方式：现在登录 Claude 订阅（在浏览器中）、[API 密钥](/features/api)，或稍后再登录 Claude 订阅。
3. 使用 API 密钥时：ccam 每次启动时从哪里读取密钥（`env:NAME` 或 `command:COMMAND`）、端点、密钥的发送方式以及模型。
4. 要加入的组，默认勾选 `default`（如果还没有任何组，则显示 `default` 组将要共享的内容）。
5. 快捷命令（如果需要）。
6. 供确认的摘要。

`Esc` 返回上一步；在第一步时则关闭向导。

## 删除账号

`x` 会在执行任何操作之前先询问，默认选中 No，并说明删除会做什么：`It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` 对于 API 密钥账号，没有需要吊销的登录；当另一个账号持有同一登录时，该账号不会退出登录。导入的账号会保留其文件夹：`It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## 颜色

颜色沿用 Claude Code 终端界面中各角色的配色，因此 ccam 放在它旁边也显得协调；组名的颜色由名称决定。ccam 会向终端查询背景色，并选择浅色或深色主题。你也可以在 `config.toml` 的顶部、第一个 `[[group]]`、`[[profile]]` 或 `[[folder]]` 之上自行设置：

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

`-daltonized` 主题适用于红绿色盲，`-ansi` 主题使用终端自己的 16 种颜色。`NO_COLOR` 会关闭颜色。

## 安全性与边界情况

- 仪表盘调用的函数与各个命令相同，因此所有规则都同样适用：移动数据前先显示计划、拒绝更改正在运行的账号、保留备份，以及登录令牌永远不会离开账号。
- 仪表盘没有试运行：不带名称的 `ccam --dry-run` 会打印此文件夹的账号将如何启动，在某个账号上按 `o` 则可以只对这一次启动进行试运行。
- 在 Claude Code 启动或登录页面打开之前，终端会先恢复原状。
- 屏幕太小、放不下对话框时，会先省略选项卡中的行；过长的行用 `…` 截断，而不是换行。在较矮的屏幕上，接下来省略空行和选项卡说明，让按键提示行保持可见。宽度小于 40 列或高度小于 8 行时，仪表盘只显示 `Too small: ccam needs 40x8 (now 36x7)`。
- 正在应用更改、删除、更新或修复时按 `Ctrl+C`，会等它完成：状态行显示 `ccam quits when it is done.`
- 用 `U` 更新 ccam 之后，仪表盘从添加账号或登录返回时不会再次提供更新。

## 相关内容

- [通过组共享](/features/share)
- [启动 Claude Code](/features/launch)
- [管理账号](/features/accounts)
- [用量限制](/features/usage)
