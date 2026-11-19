---
title: 常见问题
metaTitle: "ccam 常见问题：Claude Code 多账号"
description: 关于 ccam、Claude Code 多账号、用量限制、共享、隐私和支持的操作系统的常见问题解答。
updated: 2026-11-19
order: 6
translation: machine
sourceHash: "090ecc3a34bd"
---

## ccam 是 Anthropic 的官方工具吗？

不是。ccam 是一个非官方的开源项目，与 Anthropic 没有关联，也未获得 Anthropic 的认可。它使用的是 [Claude Code 文档](https://code.claude.com/docs/en/env-vars)中说明的 `CLAUDE_CONFIG_DIR` 设置，用于并行运行多个账号。

## 使用多个 Claude Code 账号是否被允许？

这取决于每个账号适用的条款；如果是工作账号，还取决于你所在组织的政策：Pro 和 Max 适用 Anthropic 的[消费者条款](https://www.anthropic.com/legal/consumer-terms)（Consumer Terms），Team 和 Enterprise 适用其[商业条款](https://www.anthropic.com/legal/commercial-terms)（Commercial Terms），所有账号都适用其[使用政策](https://www.anthropic.com/legal/aup)（Usage Policy）。在使用多个账号之前，请先确认这些规定。ccam 不会提高、合并或隐藏账号的用量限制，每个账号仍会在达到自己的限制时停止。在你要求时，它会把你的工作移到你的另一个账号上；如果你开启了 [`auto_switch`](/features/next#auto-switch)（Linux 和 macOS），它也会在达到限制时自行这样做。在使用之前，请先确认你的账号是否允许这样做；决定由你自己做出。

**Team 和 Enterprise 席位。** 自动选择账号（`ccam next`、`--fallback`、`ccam pick`、[`ccam agent`](/features/agent)、仪表盘中的 `n`，以及你开启后的 [`auto_switch`](/features/next#auto-switch)）也会在同一个 Claude Team 或 Enterprise 组织的席位之间切换。Anthropic 帮助中心关于 [Team 套餐](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)的说明写道："Each seat on a Team plan is assigned to one named individual, and only that person may use it"（Team 套餐的每个席位都分配给一位指定的个人，只有此人可以使用），以及 "Users can't route their work through another user's account to use that user's remaining usage"（用户不能通过另一位用户的账号完成自己的工作，以使用该用户剩余的用量）；关于 [Enterprise 套餐](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)的说明也是如此。因此，在席位之间使用自动选择账号可能违反 Anthropic 的条款。在这些席位之间是否使用，由你自己决定：在依赖这项功能之前，请先向你所在的组织和 Anthropic 确认。个人套餐（Pro、Max）的情况请参见上一段。

## ccam 能绕过用量限制吗？

不能越过任何账号自己的限制：每个账号都保留各自的用量限制，并在达到时停止，ccam 不会提高、合并或隐藏这些限制。ccam 会显示这些限制，并在某个账号用完时，让你在自己的另一个账号上继续；开启 `auto_switch` 后，它会自行这样做（见下一个问题）。在同一个 Team 或 Enterprise 组织的席位之间，在另一个席位上继续会使用该席位剩余的用量，而对于属于另一位用户的席位，Anthropic 帮助中心不允许这样做（见上文的 Team 和 Enterprise 席位）。

## ccam 能在会话达到用量限制时自行切换账号吗？

只有在你开启后才会。在 `~/.config/ccam/config.toml` 的顶部设置 `auto_switch = true`（仅限 Linux 和 macOS；默认关闭，属于实验性功能）后，当 Claude Code 因 5 小时、每周、Opus、Sonnet 或 Fable 限制停止一个交互式会话（由 ccam 启动的会话）时，ccam 会结束该会话，并以你的另一个还有余量的账号继续这个对话：同一个 Team 或 Enterprise 组织的另一个席位，或另一个个人套餐（Pro、Max）。它从不在达到限制之前切换，从不在后台监视你的用量，也从不把工作移出或移入一个 Team 或 Enterprise 组织。

在席位上，它会把你的工作从一个席位移到另一个席位，这可能违反 Anthropic 的条款（见上文的 Team 和 Enterprise 席位）。对于个人套餐，Anthropic 的条款没有明确涉及这一点：消费者条款禁止 "through automated or non-human means"（通过自动化或非人工的方式）访问服务，除非 Anthropic 明确允许，而一次移动会代表你发送一条消息。是否开启由你自己决定。

[`auto_switch`](/features/next#auto-switch) 列出了每一项条件以及一次移动的代价。

## 如何在一台电脑上使用两个 Claude Code 账号？

用 `CLAUDE_CONFIG_DIR` 为每个账号指定单独的配置文件夹。ccam 会替你完成这件事：为每个账号运行 `ccam add`，然后用 `ccam personal` 和 `ccam work`，或者你自己命名的快捷命令（例如 `ccwork`）启动它们。[指南](/guides/multiple-claude-code-accounts)也介绍了如何手动设置。

## 如何在另一个账号上继续对话？

在该文件夹中运行 `ccam next --continue`：它会把该文件夹最近的一次对话带到下一个账号（加上 `--auto` 则选择剩余额度最多的账号）。在 Claude Code 中，[ccam 插件](/docs/install#claude-code-plugin)的 `/ccam:switch` 会对你当前所在的会话做同样的事。如果这两个账号共享对话历史（`ccam group set default --add history`，或者为它们单独建一个组：`ccam group add pair --items history --members work,side`），新账号会打开的正是那一次对话；否则 Claude Code 会把它复制到新账号的历史记录中。加上 `--brief`，新账号会以一份简短的摘要开始，而不是整个对话。

## ccam 会查看或发送我的登录信息吗？

ccam 从不在账号之间共享或复制你的 Claude 登录信息，也从不复制 MCP 服务器的登录。它读取账号的访问令牌，只是为了向 Anthropic 的用量接口查询该账号的用量：用于 `ccam usage`、仪表盘、`ccam doctor` 和自动选择（`ccam next --auto`、`ccam pick`、`--fallback`，以及会话达到限制时的 `auto_switch`）。只有在运行 `ccam usage --refresh` 时才会续期过期的令牌；设置 `usage_source = "claude"` 时，它不会把令牌发送到任何地方，读取保存的登录信息（该文件或 macOS 钥匙串中的条目）也只是为了判断它是否仍然有效。`ccam ext` 和 `ccam doctor` 会读取账号的 MCP 登录的名称，但从不读取其令牌。ccam 没有遥测。除了这些查询，它唯一会自行发出的请求是仪表盘向 GitHub 检查是否有更新的 ccam 版本，其中不包含任何关于你或你的账号的信息。

## 如何在远程服务器上登录？

像平常一样在 SSH 会话中运行 `ccam add` 或 `ccam login <name>`。Claude Code 在服务器上等待登录，所以你电脑上的浏览器最后会停在一个无法加载的 `http://localhost:<port>/callback?code=…` 页面。从浏览器复制这个地址并粘贴到终端中：ccam 会把它交给服务器上的 Claude Code，由它完成登录。这个地址包含一次性代码，请不要分享给他人。

## 我的账号存储在哪里？

在 `~/.local/share/ccam/profiles/<name>/` 中，每个账号一个文件夹。在 macOS 上，Claude Code 把登录信息本身保存在钥匙串中，条目名称由该文件夹推导而来，所以 `ccam rename` 会在那里保留文件夹的旧名称。

## 我现有的 ~/.claude 会怎样？

不会有任何变化，除非你运行 `ccam import` 或 `ccam migrate`。导入会在原位置注册该文件夹，因此直接运行 `claude` 和 `ccam <name>` 都会继续使用同一个账号和同一套设置。`ccam migrate` 会把其中的设置、`CLAUDE.md`、扩展和 MCP 服务器复制到共享的 `default` 组中，并且只读取该文件夹。

## ccam 支持哪些系统？

Linux（Ubuntu、Debian、Fedora、Arch、Alpine 等）、macOS 以及 Windows 10 和 11（原生运行或在 WSL 中运行）。在 Windows 上，共享 `settings.json` 这样的文件需要开启开发人员模式；共享文件夹则无需开启；未开启时，第一次 `ccam add` 创建的 `default` 组只包含 MCP 服务器和扩展（[在 Windows 上共享](/features/share#on-windows)）。

## 为什么某个 MCP 服务器只出现在一个账号中？

在会话中用 `claude mcp add` 添加的服务器只为该账号保存，不加 `-s user` 时还只对当前文件夹有效。下次启动该账号时，ccam 会提醒你，`ccam mcp list` 和 `ccam ext` 也会显示它们；用 `ccam mcp import <account> <server>` 把服务器移到共享组。

## ccam 会拖慢 Claude Code 吗？

通过 ccam 启动 Claude Code 大约增加 4 毫秒，这是在一台笔记本电脑上用 `scripts/perf-check.sh` 测得的。ccam 读取一个小文件，检查共享链接，然后用 `claude` 替换自身，所以默认情况下，在你工作时它并不在运行。Windows 无法替换进程，所以在 Windows 上 ccam 启动 Claude Code 后会等待它结束，除此之外什么也不做。开启 `relay` 或 `auto_switch` 后，ccam 会以同样的方式在每个交互式会话旁边持续运行，以便在另一个账号上继续这个对话。

## 如何更新 ccam？

运行 `ccam update`。它会下载最新版本，校验其校验和，然后替换 `ccam` 文件；`ccam update --check` 只告诉你是否有更新的版本。仪表盘也会在有更新版本时告诉你，按 `U` 会先询问再安装；ccam 从不自行更新。如果你安装的是 `.deb` 软件包，请改用 apt 安装新的 `.deb`。参见[更新](/docs/install#update)。

## 如何卸载 ccam？

运行 `ccam uninstall`。它会先列出将要删除的内容并征求确认，然后删除每个账号（先退出登录）、快捷命令、`~/.config/ccam`、`~/.local/share/ccam`、`~/.local/state/ccam` 和 `ccam` 命令。导入的 `~/.claude` 会保留，并拿回一份它在各个组中共享过的内容的副本。参见[安装](/docs/install#uninstall)。
