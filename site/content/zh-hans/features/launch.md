---
title: 启动 Claude Code
description: 用 ccam <name|number>、ccam run 或仪表盘以任意账号启动 Claude Code；参数原样传给 claude，--dry-run 会显示将执行的命令。
updated: 2026-11-01
order: 1
translation: machine
sourceHash: "11df08608a1e"
---

`ccam work` 以 `work` 账号启动 Claude Code，账号名之后的所有内容都原样传给 `claude`。凡是你会输入 `claude` 的地方都可以用它：在终端中、在脚本中，或者通过管道输入。

## 用法

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

给出账号名或它在 `ccam list` 中的编号。ccam 不会读取账号名之后的参数，所以 `ccam work --help` 显示的是 Claude Code 的帮助。账号名之后的 `--` 是可选的。标准输入、退出码和信号的行为与 `claude` 相同。

当第一个参数以 `-` 开头且不是 ccam 自己的选项时，ccam 会替你选择账号，并把全部参数传给 Claude Code。它使用用 [`ccam folder`](/features/folder) 为当前文件夹设置的账号，否则使用默认账号：用 `ccam use` 设置的账号；如果没有设置，则是编号最小的账号。

### 仪表盘

不带参数的 `ccam` 在标准输入和输出都是终端时打开[仪表盘](/features/dashboard)，否则打印 `ccam list`。移动到某个账号（用方向键，或输入它的编号），按 `Enter` 启动它；`n` 启动余量最多的账号，`o` 提供 `--continue`、`--resume` 或其他参数。仪表盘会在 Claude Code 启动之前关闭，因此终端会先恢复原状。

**Team 和 Enterprise 席位。** 自动选择账号（`ccam next`、`--fallback`、`ccam pick`、仪表盘中的 `n`，以及你开启后的 [`auto_switch`](/features/next#auto-switch)）也会在同一个 Claude Team 或 Enterprise 组织的席位之间切换，这可能违反 Anthropic 的条款。在依赖这项功能之前，请先阅读[常见问题](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，并向你所在的组织和 Anthropic 确认。

当 shell 设置了 `CLAUDE_CODE_OAUTH_TOKEN`、`ANTHROPIC_API_KEY` 或 `ANTHROPIC_AUTH_TOKEN` 时，Claude Code 可能会使用这个密钥或令牌，而不是账号自己的登录，所以 ccam 会在启动之前在标准错误上发出警告（`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`）。它只读取这些变量；[API 密钥账号](/features/api)则会清除它们。

默认情况下，ccam 会用 Claude Code 替换自身，因此不会有其他东西继续运行。在 `config.toml` 的顶部（第一个 `[[group]]`、`[[profile]]` 或 `[[folder]]` 之上）设置 `relay = true`（或在 Linux 和 macOS 上设置 [`auto_switch = true`](/features/next#auto-switch)）时，ccam 会改为以子进程方式启动交互式会话，并在它们旁边持续运行。这样 `/ccam:switch`，或在用量限制时的 `auto_switch`，就能自行在另一个账号上继续对话，并沿用会话的选项、权限模式、推理强度和代理；在 `/ccam:switch` 把对话移到另一个组织的账号之前，ccam 会先询问你（[详情](/features/next#relay)）。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--dry-run` | 打印 ccam 将使用的环境变量和命令，然后退出 | 关闭 |
| `--fallback` | 用于 `-p`/`--print` 运行：当账号达到用量限制时，以余量最多的已登录账号重新运行同一个命令（[详情](/features/next#run-again-at-the-limit)） | 关闭 |
| `-v`、`--version` | 打印 `ccam <version>`（仅作为第一个参数时） | |
| `-h`、`--help` | ccam 的帮助，形式为 `ccam --help` 或 `ccam run --help` | |

参数规则：

- ccam 只读取账号名之前它自己的选项。账号名之后的每个参数都属于 Claude Code，包括 `--dry-run` 和 `--help`。
- `--dry-run` 可以写成 `ccam --dry-run work`、`ccam run --dry-run work` 和 `ccam --dry-run`（用于文件夹的账号或默认账号）。`--dry-run` 和 `--fallback` 的顺序可以任意。
- `ccam run <name|number> [claude args]` 是 `ccam <name>` 的完整形式。不带账号名的 `ccam run` 使用当前文件夹的账号或默认账号。
- 账号名不能是 `add`、`list`、`run` 或 `usage` 这样的 ccam 命令名，因此 `ccam <word>` 永远不会有歧义。

试运行会打印一行可以直接粘贴到 shell 中（在 Windows 上则是 PowerShell 中）的命令。它以这次启动对环境所做的更改开头：总是包含 `CLAUDE_CONFIG_DIR` 和 `CCAM_PROFILE`；对于它丢弃的变量（例如继承来的 `CCAM_RELAY`）显示 `env -u <name>`；对于 [API 密钥账号](/features/api)，显示其端点、模型以及密钥从哪里获取（`ANTHROPIC_API_KEY="$WORK_KEY"`），从不显示密钥本身：

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

PowerShell 这一行会在之后恢复这些变量，就像 shell 中的 `VAR=x cmd` 一样，并通过 `Start-Process` 把命令行交给 Claude Code，因此 Windows PowerShell 5.1 和 PowerShell 7 的运行方式相同。当启动把某个变量设为空值时，下方会有一行 `#` 说明 7.5 之前的 PowerShell 会改为删除这样的变量。对于 npm 安装的 `claude.cmd`，这一行运行的是启动时会使用的 `cmd.exe` 命令（见下文），启动时会拒绝的参数也会让试运行失败。

## 工作原理

每次启动都经过相同的步骤，从命令行一直到 `exec`：

<!-- motion:launch -->

每个方框会在对应步骤执行时亮起，下面的几行是 Claude Code 启动时使用的环境变量和参数。

1. ccam 从 `app.Main` 开始运行，它由 `cmd/ccam/main.go` 调用。如果 ccam 是以其他名称启动的，那就是一个[快捷命令](/features/shortcuts)，直接进入启动器。
2. ccam 在 cobra 介入之前查看 `os.Args`（`app.parse`）。只有单独的 `ccam`（仪表盘）、ccam 命令名（`config.Reserved`）、`-h`/`--help` 和 `-v`/`--version` 会进入 `cli.Execute`。其他一切都是启动。
3. ccam 读取 `~/.config/ccam/config.toml`，这是它解析的唯一文件（`launcher.Launch`）。
4. ccam 按名称或编号找到账号，去掉开头的一个 `--`，并解析 `claude` 可执行文件（`launcher.Prepare`、`ResolveClaude`）。不指定名称时，它会取当前文件夹的账号，把工作目录与 `config.toml` 中的 `[[folder]]` 条目按字符串进行比较，在 macOS 和 Windows 上不区分大小写（`HereProfile`）。PATH 查找结果连同文件的修改时间和 `PATH` 的哈希值一起缓存在 `~/.local/state/ccam/claude-path` 中，因此 ccam 不会在每次启动时搜索 `PATH`。文件变化或不存在，或者 `PATH` 发生变化时，会重新查找。`config.toml` 中的 `claude_bin` 路径（会展开 `~/`）必须是可执行文件；否则启动会在写入任何内容之前以 `claude_bin X is not an executable file. Check claude_bin in config.toml` 停止。如果 `claude` 就是 ccam 自身（按 inode 比较），ccam 会拒绝它，因为那会造成循环。
5. ccam 把账号文件夹 `~/.local/share/ccam/profiles/<name>/` 与它的共享设置进行比较，只调用 `lstat` 和 `readlink`，每项一到两次系统调用（`checkLinks`、`share.Check`）。缺少的链接会当场创建。任何可能丢失数据的情况都会中止启动。
6. ccam 添加共享的 MCP 服务器（`applyMCP`；见下文）。
7. ccam 记录会话（`session.Register`）。它把一个以其 PID 命名、内容为进程启动时间的标记写入 `~/.local/state/ccam/locks/profile-<name>/<pid>`，无论 `XDG_RUNTIME_DIR` 是什么，这样每个 shell 看到的都是同样的标记。对于共享对话历史的账号，ccam 会先在 `history-<group>-<hash>` 中查找正在运行的会话（hash 由解析符号链接后的当前文件夹计算），并同样在那里注册（`historyCheck`）。

   `exec` 会保留 PID，因此该标记随后指向的就是 Claude Code 进程本身。ccam 使用标记文件而不是持有锁，是因为 Claude Code 启动的进程会继承锁，并在会话结束后继续持有它。无法写入标记时不会中止启动：ccam 会显示 `ccam: could not mark the session as running (…); starting Claude Code anyway`。

   ccam 判断账号是否正在运行的其他检查（重命名、删除、组更改、`doctor --fix`、MCP 同步）还会读取 Claude Code 自己在账号文件夹中的记录（`sessions/<pid>.json`、`daemon.lock`），因此也能发现不通过 ccam 启动的 Claude Code：直接运行的 `claude`、IDE、后台会话及其守护进程。
8. ccam 把终端标题设为 `claude · <name>`（当标准错误是终端时），不经 `fsync` 把启动时间写入 `~/.local/state/ccam/last-used/<name>`，然后调用 `syscall.Exec`（`Plan.Exec`）。

Claude Code 运行时，`CLAUDE_CONFIG_DIR` 设为账号文件夹，`CCAM_PROFILE` 设为账号名。ccam 会清除这两个变量继承来的值，以及继承来的 `CCAM_RELAY`。对于在 macOS 上用 `ccam import` 注册的 `~/.claude` 文件夹，它还会把 `CLAUDE_SECURESTORAGE_CONFIG_DIR` 设为空值，让该账号和直接运行的 `claude` 使用钥匙串中的同一个登录；在 Windows 上，当 Claude Code 把登录保存在凭据管理器中时，它也会这样做。

### 为什么用 exec

ccam 用 `claude` 替换自己的进程（`syscall.Exec`）。因此：

- ccam 不会留在内存中；
- 按键、信号和退出码都直接属于 Claude Code；
- 参数始终是一个数组，而不是 shell 字符串。

Windows 没有 `exec`。在 Windows 上，ccam 以子进程方式启动 `claude.exe`，使用相同的环境变量、参数和控制台，等待它结束，并以它的退出码退出（`Plan.Exec`）。npm 安装的 `claude.cmd` 会在 `cmd.exe` 中以完整路径运行，每个参数都加上引号，使 `cmd.exe` 不会把其中任何内容当作自己的语法（`&`、`%PATH%`）；含有换行符的参数无法这样传递，会被拒绝，并建议安装原生的 Claude Code。

Ctrl+C 会直接从控制台到达 Claude Code，所以 ccam 忽略它。会话标记记录的是 ccam 的 PID，它的存活时间与 Claude Code 完全相同，并在 Claude Code 退出时被删除。等待期间，ccam 不做任何其他事情。

会话结束后的工作（例如发现用 `claude mcp add` 添加的服务器）在下次启动时进行：ccam 把 `.claude.json` 的修改时间与 `~/.local/state/ccam/mcp/<name>/seen` 进行比较，只有文件变化时才解析它。

### 命令行如何分流

三条命令行，在 `internal/app` 中走三条路径：

<!-- motion:launch-route -->

快捷命令按 `argv[0]` 匹配，`ccam <name>` 按第一个参数匹配，两者最终都进入 `syscall.Exec`。只有 `ccam usage` 这样的 ccam 命令才会到达 cobra。

### 共享的 MCP 服务器

在默认的 `flag` 策略下，ccam 把组中的 MCP 服务器作为一个参数传给 Claude Code，放在你的参数之前：

<!-- motion:launch-mcp -->

ccam 只有在组的 `mcp.json` 存在时才添加这个参数，并把它放在 `claude` 之后（`mcp.FlagArgs`、`Plan.InsertArgs`）。你自己的 `--mcp-config` 或 `--strict-mcp-config` 仍会在此基础上生效。它是单个参数 `--mcp-config=<path>`，因为 Claude Code 的 `--mcp-config` 接受多个值：如果写成两个参数，它会把后面的提示也当作配置文件读取。

ccam 决定哪些地方不能加这个参数（`claudeVerb`）。它像 Claude Code 一样从参数中读取 Claude Code 的子命令，跳过选项及其值。`remote-control` 或 `rc` 不加这个参数（见下文），放在最前面的后台会话动词（`attach`、`logs`、`stop`、`kill`、`respawn`、`rm`）也不加。`ccam work mcp list`、`get` 和 `remove` 同样不加，因为 Claude Code 在那里会拒绝它；ccam 会改为在标准错误中说明，组的服务器也会在该账号的会话中加载：

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

其他所有子命令和所有会话都会保留这个参数。

在 `config.toml` 中设置 `mcp_strategy = "sync"` 时，ccam 改为把服务器写入账号的 `.claude.json`；如果账号已在运行，则跳过这一步。如果自上次检查以来组的 `mcp.json` 和 `.claude.json` 都没有变化，启动时也会跳过这项工作。如果账号自己有一个与组服务器同名的服务器，同步会保留账号的服务器，并说明一次。同步还会把从 `.claude.json` 中消失的副本放回去（`put back group default's MCP servers …`），并且只更改或删除它自己写入、此后没有被任何人更改的副本。

### Remote Control

`ccam work rc`（或 `ccam work remote-control`）以账号 `work` 启动 Claude Code 的 Remote Control。`claude remote-control` 会拒绝出现在其动词之前的 `--mcp-config`，在动词之后也不接受，而且它的会话只从账号的 `.claude.json` 读取 MCP 服务器。因此 flag 策略在这里的工作方式不同：

<!-- motion:launch-rc -->

1. ccam 识别出作为 Claude Code 动词的 `remote-control` 或 `rc`，即使它位于选项及其值之后（`ccam work --model opus rc`），并且不添加这个参数（`claudeVerb`、`mcp.FlagArgs`）。
2. ccam 把组的 `mcp.json` 的修改时间与账号的 `sync-stamp` 比较，并检查组的每个服务器是否仍在 `.claude.json` 中（`mcp.Synced`）。如果 `.claude.json` 已经包含组当前的服务器，或者组中没有服务器，Remote Control 会直接启动，不会询问。
3. 否则，ccam 会列出会话将缺少的服务器并询问：

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. 回答 `y` 时，ccam 按 `sync` 策略的方式复制这些服务器（`mcp.Sync`）：持有 Claude Code 的锁，先做备份，并保留账号自己的服务器。回答 `N` 或没有终端时，Remote Control 会在没有这些服务器的情况下启动，ccam 会提示 `mcp_strategy = "sync"`，它会在每次启动前复制这些服务器。

在服务器被复制之前，ccam 每次启动 Remote Control 都会询问。复制之后，只有在组的 `mcp.json` 发生变化，或者某个已复制的服务器从 `.claude.json` 中消失时，它才会再次询问。在之后的普通启动中，`.claude.json` 中的副本优先于 `--mcp-config`，因此 `ccam mcp add`、`remove` 和 `import` 会更新其他账号持有的副本（[共享的 MCP 服务器](/features/mcp)）。

### 10 ms 预算

启动路径的预算是第 95 百分位 10 ms；实际上，ccam 在 Claude Code 启动前大约增加 4 ms。启动路径从不构建 cobra 命令树、启动 Bubble Tea 或打开网络连接；如果各包 `init()` 函数总共耗时超过 5 ms，`TestPackageInitIsCheap` 就会失败。`scripts/perf-check.sh` 用 `hyperfine` 测量 `ccam --dry-run work`、`ccam work` 和一个快捷命令，超过 10 ms 即失败；`internal/launcher` 中的 `BenchmarkLaunchPath` 会循环运行这些步骤。要查看每一步的耗时，请设置 `CCAM_TRACE=1`：

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## 安全性与边界情况

- 未知的名称或编号会以 `no account "<name>"` 停止，并提示使用 `ccam list`。如果一个账号都没有，ccam 会请你用 `ccam add <name>` 添加一个。
- 如果某个共享路径变成了普通文件或指向错误的位置，ccam 会以 `stopped because of share link problems:` 停止，列出这些问题，并请你运行 `ccam doctor --fix`。它不会在启动路径上移动数据。
- 如果另一个与此账号共享对话历史（同一组）的 Claude Code 会话正在同一文件夹中运行，ccam 会警告历史记录可能混在一起，并询问 `Run anyway? [y/N]`。文件夹在解析符号链接后进行匹配。没有终端时（标准输入来自 `/dev/null` 也算没有终端），或者从后台作业运行时，它会显示警告并继续。新的 `-p` 运行会开始自己的对话，不会收到警告；`--version`、`--help` 以及不打开对话的子命令（例如 `mcp`、`auth` 或 `plugin`）完全不做检查。
- `--dry-run` 会跳过 MCP 同步、新服务器提示、`mcp list` 说明、历史检查、会话标记、标题和最近使用时间，也不会执行 exec。它仍会检查共享链接，并像正常启动一样创建缺少的链接。
- 如果在 `work` 已在运行时执行 `ccam work rc`，ccam 不会改动 `.claude.json`，会说明这一点，并在没有组中服务器的情况下启动 Remote Control；会话退出后再运行一次即可复制它们。使用 `--dry-run` 时，它只打印这条消息，既不询问也不复制。
- 使用 `mcp_strategy = "sync"` 时，同步失败会打印 `ccam: MCP sync failed (launching anyway): …`，Claude Code 会以它原有的服务器启动。

## 相关内容

- [快捷命令](/features/shortcuts)
- [用 ccam 使用 Remote Control](/guides/remote-control)
- [管理账号](/features/accounts)
- [命令参考](/docs/commands)
- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [CLAUDE_CONFIG_DIR 详解](/guides/claude-config-dir)
