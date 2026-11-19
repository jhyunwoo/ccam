---
title: 命令参考
metaTitle: "ccam 命令参考"
description: ccam 的所有命令和选项，从 ccam add、ccam <name> 到 ccam usage、ccam next 和 ccam doctor。
updated: 2026-11-19
order: 3
translation: machine
sourceHash: "28e1ad1dcbad"
---

在命令前加上 `--dry-run`，可以查看它会做什么而不实际执行，例如 `ccam --dry-run work`。它适用于启动 Claude Code（`ccam <name>`、`ccam run`、`ccam next`、`ccam attach`、`ccam ccusage`）以及 `group`、`migrate`、`folder`、`shortcut`、`pick`、`update`、`uninstall` 和 `doctor --fix`。只读取内容的命令（例如 `list`、`usage`、`share` 和 `ext`）也接受它，反正它们不会更改任何东西。其余命令会拒绝它，而不会执行，`group rename` 也是如此。

<!-- commands -->

## 启动 Claude Code

`ccam <name|number> [arguments]` 以该账号启动 Claude Code。ccam 不会读取账号名之后的参数，而是按相同顺序把它们传给 `claude`；标准输入、退出码和信号的行为与直接运行 `claude` 完全相同。账号名之后的 `--` 是可选的。

`ccam` 带参数但不带账号名时（例如 `ccam -p "summarise"`），会使用为当前文件夹设置的账号，否则使用默认账号。分别用 `ccam folder <name> [folder]` 和 `ccam use <name>` 设置。

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

账号名之前的 `--fallback` 用于 `-p`/`--print` 运行。当账号达到用量限制时，ccam 会以余量最多的已登录账号重新运行同一个命令，并只打印最终完成的那次运行的输出。参见[按文件夹使用账号](/features/folder)和[达到限制时重新运行](/features/next#run-again-at-the-limit)。

**Team 和 Enterprise 席位。** 自动选择账号（`ccam next`、`--fallback`、`ccam pick`、仪表盘中的 `n`，以及你开启后的 [`auto_switch`](/features/next#auto-switch)）也会在同一个 Claude Team 或 Enterprise 组织的席位之间切换，这可能违反 Anthropic 的条款。在依赖这项功能之前，请先阅读[常见问题](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，并向你所在的组织和 Anthropic 确认。

`ccwork` 这样的快捷命令也以同样的方式启动对应账号：`ccwork --resume` 的效果与 `ccam work --resume` 相同。参见下文的[快捷命令](/docs/commands#shortcuts)。

## 添加和删除账号

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

不指定文件夹的 `ccam import` 会查找 `~/.claude` 和 `~/.claude-*`，并逐个询问；失败的文件夹会报告为 `Skipped <path>: <reason>`，其余的照常注册。如果某个文件夹的名称已被占用，请用 `ccam import <folder> --name <name>` 单独导入该文件夹。导入的文件夹在原位置注册；ccam 从不移动或复制登录信息。

`ccam login --expired` 也会为登录已被 Claude Code 清除的账号重新登录。`ccam remove` 会在受管账号（由 ccam 创建文件夹的账号，而不是用 `ccam import` 注册的账号）的文件夹移入回收文件夹之前先让它退出登录，除非另一个账号持有同一登录。删除导入的账号时，会把它移出所在的组，并在它的文件夹中为每个共享项留下一份副本。在 macOS 上，`ccam rename` 会保留文件夹的名称，因为钥匙串按文件夹的路径保存登录信息。

`ccam add --api` 会添加一个使用 API 密钥而不是 Claude 登录的账号；ccam 在每次启动时读取密钥，从不保存它。参见 [API 密钥账号](/features/api)。

当 Claude Code 正以某个账号运行时，`ccam rename`、`ccam remove` 和各个 `ccam group` 命令会拒绝更改该账号，无论它是怎样启动的（包括后台会话）。请先退出该会话。组的更改尚未完成时，它们也会拒绝执行：``the last group change stopped halfway; finish it with `ccam doctor --fix` first``。

`ccam list --json` 会输出每个账号及其槽位、快捷命令、邮箱和文件夹。其中的 `groups` 数组列出账号所在的组，`share` 对象则把每一项（`mcp`、`settings`、`memory`、`extensions`、`history`）映射到提供它的组，或映射到 `local`。

## 快捷命令

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

快捷命令是你自己命名的命令，用来以某个账号启动 Claude Code。它和 `ccam work` 一样，会把所有参数原样传给 Claude Code，所以 `ccwork --resume` 和 `ccwork -p "hi"` 都会按预期工作。它是 `~/.local/bin` 中指向 ccam 的链接。

快捷命令是可选的。`ccam add` 向导会询问是否创建；`ccam add <name> --shortcut <command>` 则不询问直接创建。`ccam list` 会显示每个账号的快捷命令，没有时显示 `-`。

名称可以使用 a-z、0-9、`.`、`-` 和 `_`。名称不能是 ccam 的命令、其他账号的快捷命令，也不能与 `PATH` 上的其他程序同名，或是 `~/.local/bin` 中并非由 ccam 创建的文件，例如你自己复制或链接的 ccam。如果在 `config.toml` 中手动编辑的快捷命令违反这些规则、是 `claude` 或 `ccam`，或者被分配给两个账号，`config.toml` 就会加载失败，并指出是哪个账号。

`ccam rename` 会保留快捷命令的名称，`ccam remove` 会连同账号一起删除快捷命令。`ccam shims install` 会重新创建所有账号的链接，例如在你移动了 ccam 之后；某个快捷命令的位置被文件占用时，不会影响其他快捷命令。

## 共享

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

组在其账号之间共享各项；一个账号可以属于多个组，但每一项只能来自一个组。参见[用组共享](/features/share)。

## 复制现有配置

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` 会把 Claude Code 文件夹中的 `settings.json`、`CLAUDE.md`、扩展和用户 MCP 服务器复制到一个组中，让共享这些项的每个账号都能用上。该文件夹只会被读取。

- `--settings`、`--memory`、`--extensions`、`--mcp` 和 `--history` 只复制你指定的项。不加这些选项时，在终端中会显示一个清单供你选择，除对话历史外的所有项都已勾选；在其他情况下，会复制除对话历史以外的所有内容。
- 即使 `~/.claude` 为空或不存在，只要 `~/.claude.json` 中有服务器，`--mcp` 也能使用。
- `-g` 或 `--group` 用于选择组。
- `-y` 或 `--yes` 不询问直接复制，文件内容不同时保留组中的副本。

MCP 登录不会被复制，因为复制过去的登录会让其他持有者退出登录；计划会列出需要用 `/mcp` 登录的账号。参见[复制现有配置](/docs/sharing#copy-an-existing-setup)。

## MCP 服务器

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` 会把只有该账号拥有的服务器移到它的组中，也包括它只为某一个文件夹添加的服务器（普通的 `claude mcp add`），这些服务器随后会在所有文件夹中加载。如果组中已有同名但定义不同的服务器，该服务器会被跳过；加上 `--force` 可替换组中的副本，ccam 会先备份它。

## 用量与切换

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam sends no token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --continue --brief    # a new conversation from a summary of the last one
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agent "…"                  # from Claude Code: a subagent that runs as another account
ccam agent --resume <id> "…"    # follow that agent up
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
ccam ccusage claude daily       # Claude Code's tokens and costs, every account (ccusage)
ccam ccusage -a work blocks --active   # one account's current 5-hour block
```

- `next --auto`、`--fallback` 和 `ccam pick` 的选择方式相同：优先选择低于 `auto_max_usage`（`config.toml`，默认 100）的账号，按 `auto_strategy`（`headroom`、`expiring` 或 `balance`）选择，并计入 `--model` 的每周时间窗口。它们从不选择已达 100% 的账号（Claude Code 报告已满的时间窗口在重置前按 100% 计算）或 API 密钥登录，并且只在出发账号所属的组织内选择（用 `--cross-org` 或 `cross_org = true` 扩大范围）。
- `--wait` 的值写在 `=` 之后（`--wait=2h`）；用空格分隔时会被拒绝，并给出提示。当没有可选的账号时，它们以退出码 75 退出，或者在使用 `--wait` 时等待第一个账号重置。
- 即使两个账号不共享历史记录，`ccam next --continue` 也能把对话带过去：此时由 Claude Code 把对话复制到新账号中。使用 `--brief` 时，新账号会改为开始一个新对话，其起点是 ccam 不借助模型为上一个对话写的摘要（[`--brief`](/features/next#brief)）。
- 与 `ccam <name>` 不同，`ccam next` 会读取自己的选项，所以 Claude Code 的选项要放在 `--` 之后；没有它时 ccam 会以 `unknown flag` 停止。
- `ccam agent` 在此文件夹中用 `claude -p` 以另一个账号（而不是调用它的账号）运行任务，沿用调用它的会话的权限模式和推理强度，并输出 Claude 的最后一条回复。该账号达到限制时，对话会以另一个账号继续（[其他账号上的子代理](/features/agent)）。
- 在 Claude Code 中，插件的 `/ccam:switch [account]` 用于交接对话，`ccam hook` 则是插件的钩子所运行的命令。
- 在 `config.toml` 中设置 `auto_switch = true`（仅限 Linux 和 macOS；默认关闭）时，用量限制钩子还会结束由 ccam 启动的交互式会话，并在另一个还有余量的账号上继续它，但只限同类账号：同一组织的另一个席位，或另一个个人套餐（[`auto_switch`](/features/next#auto-switch)）。

参见[切换到下一个账号](/features/next)、[为脚本选择账号](/features/pick)、[其他账号上的子代理](/features/agent)和[所有账号的后台会话](/features/agents)。

## 维护

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam uninstall       # remove every account, ccam's folders and ccam itself
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` 在替换 ccam 之前会校验版本的 SHA-256 校验和，`--version v0.2.0` 用于安装特定版本。通过 `.deb` 软件包安装的 ccam 副本则用 apt 更新。参见[更新](/docs/install#update)。

`ccam uninstall` 会先列出要删除的内容并征求确认：像 `ccam remove` 那样删除每个账号，以及快捷命令、ccam 的文件夹和 `ccam` 命令。导入的文件夹会保留在原位置。参见[卸载](/docs/install#uninstall)。

## 文件

| 路径 | 内容 |
| --- | --- |
| `~/.config/ccam/config.toml` | 账号、组和设置（`theme`、`usage_source`、`relay`、`auto_switch`、`cross_org` …）。设置写在文件顶部，位于第一个 `[[group]]`、`[[profile]]` 或 `[[folder]]` 之上：写在某个表下面的设置属于该表，所以 ccam 会在每次启动时以及在 `ccam doctor` 中对它发出警告，而保存配置的命令在它被移走之前会拒绝执行。ccam 保存该文件时不会保留注释 |
| `~/.local/share/ccam/profiles/<name>/` | 每个账号的 `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | 共享的设置、MCP 服务器、扩展和历史记录 |
| `~/.local/state/ccam/` | 用量缓存、用量限制记录、锁、记录和其他状态（见表格下方的列表） |
| `~/.local/bin/<shortcut>` | `ccwork` 这样的快捷命令，即指向 ccam 的链接 |

状态文件夹中保存：

- 用量缓存和用量限制记录（`usage/<name>.json`、`usage/<name>.limit.json`）；
- 最近使用时间；
- 会话标记和锁（`locks/`）；
- ccam 创建的快捷命令（`shortcuts.json`）；
- 尚未完成的组更改（`group-change.json`）；
- `/ccam:switch` 和用量限制钩子留下的记录（`handoff/<account>/<session>.json`）；
- `ccam pick` 或 `auto_switch` 最近一次选择各个账号的时间（`picks/<name>`；一次选择在 2 分钟内有效）；
- `--brief` 的摘要，保留 7 天（`brief/<session>.md`）。

设置 `XDG_CONFIG_HOME`、`XDG_DATA_HOME` 或 `XDG_STATE_HOME` 可以改变这些文件夹的位置。在 Windows 上，`~` 是 `%USERPROFILE%`，快捷命令是 `<shortcut>.exe`。
