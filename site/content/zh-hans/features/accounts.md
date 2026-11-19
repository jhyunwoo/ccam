---
title: 管理账号
description: 用 ccam add、ccam login、ccam list、ccam rename、ccam remove 和 ccam import 添加、登录、列出、重命名、删除和导入 Claude Code 账号。
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

每个 ccam 账号都是一个独立的 Claude Code 配置文件夹，拥有自己的登录。用这些命令可以添加账号、为账号登录、列出、重命名和删除账号。它们还能注册你已有的文件夹。

## 用法

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account that was never logged in or whose login expired
ccam logout side                           # log out: the login and its MCP sign-ins
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, logs out, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

凡是命令接受 `<name|number>` 的地方，你都可以给出账号名或它的槽位号（`ccam list` 中的 `#` 列）。

`ccam use` 只更改哪个账号是默认账号（`config.toml` 中的 `default`）。

在终端中不带选项运行 `ccam add` 会打开向导。它会依次经过以下步骤：

1. 账号名称，除非你已经给出；
2. 账号的登录方式：现在用 Claude 登录、使用 API 密钥，或稍后再用 Claude 登录；
3. 要加入哪些组，其中 `default` 已勾选（如果还没有任何组，则显示 `default` 组将要共享的内容）；
4. `Create a shortcut command for work? [y/N]`；
5. 一份供你确认的摘要。

按 `Esc` 返回上一步。只要给出 `--group`、`--no-group`、`--shortcut` 或 `--no-login` 中的任何一个，或者不在终端中运行（例如输入来自管道），就会跳过向导。此时必须指定名称，并且除非 `--shortcut` 指定了名称，否则不会创建快捷命令。

不带 `--group` 时，新账号会加入 `default` 组。第一个账号会创建这个组，共享 MCP 服务器、设置、`CLAUDE.md` 和扩展；参见[用组共享](/features/share)。在 Windows 上，只有在开发人员模式下才能链接文件。未开启该模式时，第一个账号创建的 `default` 只包含 `mcp` 和 `extensions`，并会说明这一点（`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`），而且 ccam 会拒绝用 `--group` 加入共享设置、memory 或历史记录的组。

`ccam list` 显示槽位、名称、快捷命令、邮箱、组和最近使用时间，默认账号前面标有 `*`。邮箱是该文件夹当前登录的邮箱（取自它的 `.claude.json`）；当它与 ccam 记录的邮箱不同时，记录的邮箱会以 `(was …)` 的形式跟在后面：

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

```console
$ ccam list
   #  Name      Shortcut  Email                                     Groups   Last used
*  1  work      ccwork    other@company.com (was work@example.com)  default  2m ago
```

使用 `--json` 时，每个账号是一个对象，包含 `slot`、`shortcut`、`name`、`email`（与表格中一样，是当前登录的邮箱）、`default`、`dir`、`last_used`（账号从未启动过时省略）、`groups` 和 `share`。`share` 把每一项（`mcp`、`settings`、`memory`、`extensions`、`history`）映射到提供它的组，或映射到 `"local"`。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `ccam add --no-login` | 创建账号但不登录；之后用 `ccam login <name>` 登录 | 关闭 |
| `ccam add --group <group>` | 加入这个组，而不是 `default`。可重复使用；`default` 不存在时会被创建 | `default` |
| `ccam add --no-group` | 不加入任何组：账号不共享任何内容 | 关闭 |
| `ccam add --shortcut <command>` | 同时在 `~/.local/bin` 中创建这个快捷命令 | 无 |
| `ccam add --api --key-from <source>` | 使用 API 密钥而不是登录的账号；参见 [API 密钥账号](/features/api) | 关闭 |
| `ccam list --json` | 以 JSON 输出账号 | 关闭 |
| `ccam remove -y`、`--yes` | 不询问直接删除 | 关闭 |
| `ccam logout -y`、`--yes` | 不询问直接退出登录 | 关闭 |
| `ccam import --name <name>` | 你所指定文件夹的账号名；不指定文件夹时会被拒绝（`--name needs a folder: ccam import <folder> --name ops`） | 由文件夹名推导 |
| `ccam import -y`、`--yes` | 不询问直接导入找到的所有文件夹 | 关闭 |
| `--dry-run` | `ccam list` 和 `ccam completion` 支持它，因为它们不做任何更改。`add`、`login`、`logout`、`use`、`rename`、`remove` 和 `import` 会拒绝它，并输出 `` --dry-run is not supported by `ccam add`; nothing was changed ``（其中是对应的命令名） | 关闭 |

`ccam login`（除 `--expired` 外）、`ccam use` 和 `ccam rename` 没有自己的选项。`ccam completion` 接受 `bash`、`zsh` 或 `powershell`；其他值会失败并显示 `only bash, zsh and powershell are supported`。

账号名由 1 到 32 个 `a-z`、`0-9`、`-` 和 `_` 字符组成，必须以字母或数字开头。名称不能全是数字，因为 `ccam 2` 表示槽位 2。名称也不能是 ccam 的命令（`add`、`list`、`ls`、`run`、`use`、`next`、`usage`、`u`、`share`、`mcp`、`login`、`logout`、`rename`、`remove`、`rm`、`import`、`migrate`、`doctor`、`shims`、`shortcut`、`completion`、`help`、`version`、`statusline`、`update`、`uninstall`、`folder`、`pick`、`ext`、`group`、`groups`、`hook`、`agent`、`agents`、`attach`、`ccusage`），这样 `ccam <name>` 永远不会有歧义。这份列表是 `internal/config` 中的 `config.Reserved`。

## 工作原理

添加账号会依次创建它的文件夹、加入它的组、保存配置、登录，然后创建快捷命令。运行 `ccam add work --shortcut ccwork` 时，ccam 会执行以下步骤（`internal/cli/profiles.go` 中的 `addProfile`）：

<!-- motion:accounts-add -->

1. 在创建任何内容之前，ccam 会检查名称（`config.ValidateName`）和快捷命令（`shim.Check`）。
2. ccam 以 `0700` 权限创建 `~/.local/share/ccam/profiles/work/`，并为账号分配最小的空闲槽位号，因此已删除账号的编号会被重复使用（`profile.Add`、`Config.NextSlot`）。第一个账号会成为默认账号。
3. 账号加入它的组（除非 `--group` 或 `--no-group` 另有指定，否则为 `default`；第一个账号会创建 `default`），ccam 把新建的空文件夹链接到这些组，例如把 `settings.json` 链接到 `groups/default/settings.json`（`applyShare`，先用 `share.Build` 再用 `Apply`）。如果这一步失败，新文件夹会被删除，因为残留的文件夹会让每次重试都因 `folder already exists` 而失败。
4. 账号通过 `fsutil.AtomicWrite` 保存到 `~/.config/ccam/config.toml`：先写临时文件，执行 `fsync`，再 `rename`，因此即使崩溃也不会留下写了一半的配置。
5. ccam 在 `CLAUDE_CONFIG_DIR` 设为账号文件夹的情况下启动 `claude auth login`，并用 `--email` 传入账号上一次的邮箱，这样登录页面会预先填好（`runLogin`）。如果登录失败，账号仍会保留，ccam 会打印错误，并提示你稍后运行 `ccam login work`。
6. ccam 从 `profiles/work/.claude.json` 中读取 `oauthAccount`（`profile.ReadIdentity`）：邮箱、账号 ID 和组织 ID、组织名称以及套餐（`claude_max`、`claude_team` …）。这个文件可能有好几 MB，所以 ccam 只扫描这一个键，而不解析整个文件（`TopLevelValue`）。ccam 把这些信息保存在 `config.toml` 的 `[profile.login]` 中。
7. ccam 把 `~/.local/bin/ccwork` 链接到 ccam 可执行文件，后者启动时会在 `config.toml` 中查找 `argv[0]`；同时把这个链接记录在 `~/.local/state/ccam/shortcuts.json` 中（`shim.Install`）。如果无法创建链接，账号仍会被添加，ccam 会说明原因（`Could not create shortcuts: …`），摘要中显示 `shortcut: ccwork (not created; see above)`。

登录可能需要几分钟，期间另一个 ccam 进程可能会保存配置。因此需要等待的命令会重新加载 `config.toml`，只应用这一项更改（这里是邮箱），然后保存（`env.update`）。出于同样的原因，`ccam remove` 在提示之后也会重新加载。

`ccam login` 对已有账号执行第 5 步和第 6 步。Claude Code 只把邮箱当作登录页面的提示，所以你最终登录的可能是与之前不同的 Claude 账号；ccam 会把新的登录与 `[profile.login]` 比较，并指出这一点（`ccam: this login is other@company.com, but work was dev@company.com before.`）。它会保留新的登录，并告诉你如何撤销：``To undo it, run `ccam login work` again and sign in as dev@company.com.``

登录到另一个 Claude 账号还会丢弃该账号的[用量限制记录](/features/usage#the-limit-record)，因为那些限制属于另一个账号。当另一个 ccam 账号登录到同一个 Claude 账号时，ccam 会警告这两个账号共享用量限制。

`ccam login --expired` 会依次登录每个从未登录过、或登录已过期的账号：其刷新令牌已过期，或者服务器拒绝后 Claude Code 清除了该登录。在打开第一个浏览器之前，它会检查 `config.toml` 能否保存。

`ccam logout` 在账号中运行 `claude auth logout` 并清除保存的邮箱；登录记录会保留，因此下一次登录会与它进行比较。Claude Code 会吊销该登录，并删除它为该账号保存的内容，包括 MCP 服务器的登录。所以在以下情况下，ccam 会先发出警告并询问确认（`--yes` 跳过询问）：

- 登录中保存有 MCP 登录；
- 账号正在运行（它的会话也会退出登录）；
- 另一个账号的文件夹持有同一登录（复制的文件夹），退出登录也会吊销它。

例如：

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

要以另一个 Claude 账号登录，只需运行 `ccam login <name>`，并且它会保留 MCP 登录。

当账号属于不止一个组织时，`ccam list` 会增加一列 `Organization`（个人套餐显示为 `personal`），`ccam list --json` 则为每个账号提供 `plan` 和 `org`。

### 在远程机器上登录

通过 SSH 连接时，`claude auth login` 中浏览器那一半可能会出问题。Claude Code 在远程机器上的回调服务器（`http://localhost:<port>/callback`）上等待登录完成，但完成登录的浏览器运行在你的电脑上，在那里这个地址指向不了任何地方，浏览器会显示 `ERR_CONNECTION_REFUSED` 之类的错误。这个地址中仍带有一次性代码，所以 ccam 可以把它转交过去：

<!-- motion:accounts-relay -->

ccam 不会把你的终端直接交给 Claude Code（`runLogin`）。它把 Claude Code 的标准输入连接到一个管道，并把你输入的每一行复制进去（`relayLogin`），因此在 Claude Code 的 `Paste code here` 提示处粘贴的代码依然能送达。当某一行是回调地址时，ccam 只在以下条件都满足时才接受它：主机是 `localhost`、`127.0.0.1` 或 `[::1]`，指定了端口，路径为 `/callback`，并且同时带有 `code` 和 `state`；`http://` 可以省略（`callbackURL`）。然后 ccam 像浏览器那样，在服务器上请求这个地址，超时为 10 秒，且不跟随重定向（`deliverCallback`）。请求到达正在等待的 `claude auth login` 的回调服务器（图中的虚线）。Claude Code 收到代码并保存登录，在 Linux 上保存到账号文件夹中的 `.credentials.json`。粘贴的内容永远不会让 ccam 连接其他主机。

登录开始之前，ccam 会打印：

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

成功时打印 `→ Passed the sign-in to Claude Code on this machine.`

登录在终端上等待输入时，ccam 会关闭终端的回显和行模式，你每输入一个字符就显示一个 `*`，因此粘贴的代码或回调地址不会出现在屏幕上。行编辑由 ccam 自己处理（`readLine`）：`Enter` 发送这一行，`Backspace` 删除一个字符，在空行上按 `Ctrl+D` 结束输入。登录结束时会恢复终端设置；收到 `Ctrl+C`、`SIGTERM` 或 `SIGHUP` 时也会恢复，它们仍会结束 ccam。如果标准输入和标准错误不是终端，各行会按原样复制。

### 重命名、删除和导入

`ccam rename` 和 `ccam remove` 作用于账号的文件夹，以及 ccam 在其他位置以该账号名保存的内容。ccam 在 `nameState` 中列出这些路径：

- `~/.local/state/ccam/last-used/<name>`；
- `~/.local/state/ccam/mcp/<name>`；
- 用量缓存 `~/.local/state/ccam/usage/<name>.json` 和 `<name>.statusline.json`；
- 用量限制记录 `<name>.limit.json`；
- `~/.local/state/ccam/handoff/<name>/` 中的会话记录。

`ccam remove work` 按以下顺序处理它们：

<!-- motion:accounts-remove -->

1. 如果账号正在运行，`ccam remove` 会先拒绝，否则询问确认，并在你回答后再检查一次（这样能发现在提示期间启动的会话）。
2. 在做任何更改之前，它会检查 `config.toml` 能否保存。
3. 如果账号共享扩展，ccam 会改写组中的 `installed_plugins.json` 和 `known_marketplaces.json`（`RewritePluginPaths`），因为 Claude Code 保存的是插件的绝对路径，这些路径可能经过该账号的文件夹。
4. 接下来，受管账号（由 ccam 创建文件夹的账号，而不是用 `ccam import` 注册的账号）会退出登录：在它的文件夹中运行 `claude auth logout`，这样回收站中就不会留下仍然有效的登录（在 macOS 上，也不会留下以该文件夹命名、之后同名账号会继承的钥匙串条目）。对于 API 密钥账号，以及另一个账号的文件夹持有同一登录（退出登录也会吊销它）的情况，会跳过这一步。退出登录失败时会报告出来，删除继续进行。
5. ccam 把文件夹移到 `~/.local/share/ccam/trash/<name>-<time>/`（`profile.Trash`），账号从 `config.toml` 中移除，`nameState` 中的路径被删除，以免之后同名的账号继承它们，快捷命令的链接也会被删除：`✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`。

如果最终仍无法保存 `config.toml`，文件夹会从回收站移回。要恢复已删除的账号，把文件夹从回收站移回并重新登录即可。

`ccam rename` 重命名文件夹，并把其 `plugins/` 中的插件路径从旧文件夹改写为新文件夹。它会把 `nameState` 中的路径移到新名称下；如果该账号原本是默认账号，则仍保持为默认账号。它还会更新会话记录，因此 `ccam next --continue` 仍能找到它们。快捷命令保持原名。如果插件路径改写失败，或者无法保存 `config.toml`，文件夹会被改回原名。

在 macOS 上，以及在 Claude Code 把登录保存在凭据管理器中的 Windows 上，文件夹本身会保留旧名称，因为登录和 MCP 登录是按文件夹的路径保存的，移动过的文件夹会退出登录。此时配置会在新名称旁把该文件夹记录为 `folder`，ccam 会显示 `its folder stays …: the macOS Keychain keeps the login under the folder's path`。

`ccam import` 把文件夹在原位置注册为 `config.toml` 中的 `dir`。它从不移动或复制文件夹：复制的刷新令牌会让各个副本互相把对方登出。

不指定文件夹时，它会查找尚未注册的 `~/.claude` 和 `~/.claude-*` 文件夹并逐个询问。名称取自文件夹（`~/.claude` 为 `main`，`~/.claude-work` 为 `work`，`~/.claude-2` 为 `acct2`）。如果 ccam 无法推导出名称，或者名称已被占用，它会列出该文件夹，并附上应运行的 `--name` 命令。找到的文件夹导入失败时会报告为 `Skipped <path>: <reason>`，其余的照常导入。

ccam 先问完所有问题，再重新加载 `config.toml` 并注册这些文件夹，因此其他 ccam 命令在此期间保存的更改会保留。`✓` 行在保存之后才显示。

没有 `CLAUDE_CONFIG_DIR` 时，Claude Code 把 `.claude.json` 放在你的主目录中、`~/.claude` 旁边，因此当你导入 `~/.claude` 而它没有自己的 `.claude.json` 时，ccam 会创建链接 `~/.claude/.claude.json`，指向 `~/.claude.json`。这样直接运行 `claude` 和运行 `ccam main` 看到的是同一个账号。

如果 `~/.claude` 在 `~/.claude.json` 之外还有自己的 `.claude.json`，ccam 会保留两者并说明这一点。`ccam main` 会使用里面的文件，而直接运行 `claude` 会使用外面的文件，两者的登录信息和 MCP 服务器互不相通。要只保留一个，请把里面的文件换成链接（[doctor](/features/doctor) 也会报告这一点）。

当该文件夹与某个已注册账号持有同一登录时（复制的文件夹），ccam 会说明：先续期的一方会让另一方退出登录，所以使用前请运行 `ccam login <name>`。

导入的账号不属于任何组，直到你把它加入某个组（`ccam group join <group> <name>`），也不会创建快捷命令。

`ccam completion bash|zsh` 会打印一个补全脚本，它也能补全账号名。

## 安全性与边界情况

- 账号正在运行时，`ccam rename` 和 `ccam remove` 会拒绝执行（`refuseRunning`），并显示 ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``，因为会话会继续写入该文件夹。无论 Claude Code 是怎样启动的，这项检查都适用：通过 ccam、直接运行 `claude`、IDE、后台会话及其守护进程。
- 当另一个 ccam 进程正在续期或更改该登录时，`ccam rename`、`ccam remove` 和 `ccam logout` 会拒绝执行：``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``。
- 当另一个 ccam 进程正在更改组时，`ccam rename` 和 `ccam remove` 会拒绝执行，并显示 `another ccam process is changing the groups right now; try again when it is done`；组更改中途停止后，则显示 ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``。
- 对导入的账号执行 `ccam remove` 会保留文件夹及其登录，供直接运行的 `claude` 使用。如果该账号在某些组中，它会先像 `ccam group leave` 那样离开这些组：每个共享项都在文件夹中变成组内容的一份副本，组的 MCP 服务器写入它的 `.claude.json`。然后取消注册：`✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`。ccam 为它的 `.claude.json` 做的备份会移到 `<folder>/.ccam-backup/mcp-<time>`；如果无法移动，它们会留在 `~/.local/state/ccam/mcp/<name>/backups` 中，消息会说明这一点。你的 `~/.claude` 不会被移到回收站。
- 两个账号不能使用同一个文件夹：`config.toml` 中两次指定同一个文件夹会加载失败。
- 如果 `profiles/` 中已存在同名文件夹，`ccam add` 会拒绝该名称，并提示使用 `ccam import`。在 macOS 上，如果钥匙串中仍保存着旧版本 ccam 留下的、属于新文件夹路径的登录或 MCP 登录，它也会停止，并给出删除它们的 `security delete-generic-password` 命令。如果无法保存 `config.toml`，新文件夹会被再次删除。
- `ccam import` 会拒绝已注册的文件夹、不是文件夹的路径，以及已被占用的名称。当该文件夹是 `~/.claude` 的副本（同一登录）时，即使 `~/.claude` 本身未注册，ccam 也会发出警告。
- 账号正在运行时，`ccam logout` 会先询问；`ccam login` 不做检查。
- 如果转交无法到达 Claude Code（例如登录已经结束），它会打印 `ccam: could not pass the sign-in to Claude Code:` 以及原因。HTTP 错误通常表示代码已过期；请重新开始登录。

## 相关内容

- [启动 Claude Code](/features/launch)
- [快捷命令](/features/shortcuts)
- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [在一台电脑上使用多个 Claude Code 账号](/guides/multiple-claude-code-accounts)
