---
title: 指令參考
metaTitle: "ccam 指令參考"
description: 所有 ccam 指令和選項，從 ccam add、ccam <name> 到 ccam usage、ccam next 和 ccam doctor。
updated: 2026-11-19
order: 3
translation: machine
sourceHash: "28e1ad1dcbad"
---

在指令前加上 `--dry-run`，就能查看它會做什麼而不實際執行，例如 `ccam --dry-run work`。它適用於啟動 Claude Code（`ccam <name>`、`ccam run`、`ccam next`、`ccam attach`、`ccam ccusage`）以及 `group`、`migrate`、`folder`、`shortcut`、`pick`、`update`、`uninstall` 和 `doctor --fix`。只會讀取的指令，例如 `list`、`usage`、`share` 和 `ext`，會接受它，而且本來就不會變更任何東西。其餘指令則會拒絕它，而不會執行，`group rename` 也是如此。

<!-- commands -->

## 啟動 Claude Code

`ccam <name|number> [arguments]` 會以該帳號啟動 Claude Code。ccam 不會讀取名稱後面的參數，而是依相同順序傳給 `claude`；標準輸入、結束代碼和訊號的行為都和直接執行 `claude` 完全相同。名稱後面的 `--` 可加可不加。

只有參數而沒有帳號名稱的 `ccam`，例如 `ccam -p "summarise"`，會使用為目前資料夾設定的帳號，否則使用預設帳號。分別用 `ccam folder <name> [folder]` 和 `ccam use <name>` 設定。

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

帳號名稱前的 `--fallback` 用於 `-p`/`--print` 執行。帳號達到用量上限時，ccam 會以剩餘額度最多的已登入帳號再次執行同一個指令，並且只輸出完成的那次執行的結果。請參閱[依資料夾使用帳號](/features/folder)和[達到上限時重新執行](/features/next#run-again-at-the-limit)。

**Team 和 Enterprise 席位。** 自動選擇帳號（`ccam next`、`--fallback`、`ccam pick`、儀表板的 `n`，以及你開啟後的 [`auto_switch`](/features/next#auto-switch)）也會在同一個 Claude Team 或 Enterprise 組織的席位之間切換，這可能會違反 Anthropic 的條款。在依賴這項功能之前，請先閱讀[常見問題](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，並向你所屬的組織和 Anthropic 確認。

`ccwork` 這類快捷指令也會以同樣的方式啟動對應的帳號：`ccwork --resume` 的效果和 `ccam work --resume` 相同。請參閱下方的[快捷指令](/docs/commands#shortcuts)。

## 新增和移除帳號

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

不指定資料夾的 `ccam import` 會尋找 `~/.claude` 和 `~/.claude-*`，並逐一詢問；失敗的資料夾會以 `Skipped <path>: <reason>` 回報，其他資料夾照樣登錄。如果某個資料夾的名稱已被使用，請用 `ccam import <folder> --name <name>` 單獨匯入該資料夾。匯入的資料夾會在原位置登錄；ccam 絕不會移動或複製登入資訊。

`ccam login --expired` 也會登入那些登入資訊已被 Claude Code 清除的帳號。`ccam remove` 會先將受管理的帳號（由 ccam 建立資料夾的帳號，而非以 `ccam import` 登錄的帳號）登出，再把它的資料夾移到垃圾桶，除非另一個帳號持有同一組登入資訊。移除匯入的帳號時，會將它移出所屬的群組，並在它的資料夾中留下每個共用項目的副本。在 macOS 上，`ccam rename` 會保留資料夾的名稱，因為鑰匙圈是以資料夾的路徑保存登入資訊。

`ccam add --api` 會新增一個使用 API 金鑰而不是 Claude 登入的帳號；ccam 會在每次啟動時讀取金鑰，絕不儲存它。請參閱 [API 金鑰帳號](/features/api)。

當 Claude Code 正以某個帳號執行時（不論它是如何啟動的，背景工作階段也算在內），`ccam rename`、`ccam remove` 和 `ccam group` 系列指令會拒絕變更該帳號。請先結束該工作階段。群組變更尚未完成時，它們也會拒絕：``the last group change stopped halfway; finish it with `ccam doctor --fix` first``。

`ccam list --json` 會輸出每個帳號及其槽位、快捷指令、電子郵件和資料夾。其中的 `groups` 陣列會列出該帳號所屬的群組，`share` 物件則會把每個項目（`mcp`、`settings`、`memory`、`extensions`、`history`）對應到它所來自的群組，或 `local`。

## 快捷指令

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

快捷指令是你自己命名的指令，用來以某個帳號啟動 Claude Code。它和 `ccam work` 一樣，會把所有參數原封不動地傳給 Claude Code，所以 `ccwork --resume` 和 `ccwork -p "hi"` 都會如預期運作。它是 `~/.local/bin` 中指向 ccam 的連結。

快捷指令是選用的。`ccam add` 設定精靈會詢問是否建立；`ccam add <name> --shortcut <command>` 則不詢問直接建立。`ccam list` 會顯示每個帳號的快捷指令，沒有時顯示 `-`。

名稱可以使用 a-z、0-9、`.`、`-` 和 `_`。名稱不能是 ccam 的指令、其他帳號的快捷指令，也不能和 `PATH` 上的其他程式同名，亦不能是 `~/.local/bin` 中並非由 ccam 建立的檔案，例如你自己的 ccam 副本或連結。若在 `config.toml` 中手動編輯的快捷指令違反這些規則、是 `claude` 或 `ccam`，或同時指定給兩個帳號，`config.toml` 就會載入失敗，並指出是哪個帳號。

`ccam rename` 會保留快捷指令的名稱，`ccam remove` 會連同帳號一起刪除快捷指令。`ccam shims install` 會重新建立所有帳號的連結，例如在你移動 ccam 之後；某個快捷指令的位置被檔案佔用時，不會妨礙其他快捷指令。

## 共用

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

群組會在其中的帳號之間共用項目；一個帳號可以屬於多個群組，但每個項目只會來自一個群組。請參閱[用群組共用](/features/share)。

## 複製現有的設定

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` 會把 Claude Code 資料夾中的 `settings.json`、`CLAUDE.md`、擴充功能和使用者 MCP 伺服器複製到一個群組，讓共用這些項目的每個帳號都能使用。該資料夾只會被讀取。

- `--settings`、`--memory`、`--extensions`、`--mcp` 和 `--history` 只會複製你指定的項目。不加這些選項時，在終端機中會顯示一份核取清單讓你選擇，除了對話記錄以外的項目都已勾選；在其他情況下，會複製對話記錄以外的所有內容。
- 即使 `~/.claude` 是空的或不存在，只要 `~/.claude.json` 中有伺服器，`--mcp` 也能運作。
- `-g` 或 `--group` 用來選擇群組。
- `-y` 或 `--yes` 不詢問直接複製，檔案內容不同時保留群組中的副本。

MCP 的登入不會被複製，因為複製的登入會讓其他持有者被登出；計畫會列出需要用 `/mcp` 登入的帳號。請參閱[複製現有的設定](/docs/sharing#copy-an-existing-setup)。

## MCP 伺服器

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` 會把只有該帳號擁有的伺服器移到它的群組，也包括它只為某個資料夾設定的伺服器（一般的 `claude mcp add`），這些伺服器之後會在每個資料夾中載入。如果群組已有同名但定義不同的伺服器，該伺服器會被略過；加上 `--force` 可取代群組中的副本，ccam 會先備份它。

## 用量和切換

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

- `next --auto`、`--fallback` 和 `ccam pick` 的選擇方式相同：優先選低於 `auto_max_usage`（`config.toml`，預設 100）的帳號，依 `auto_strategy`（`headroom`、`expiring` 或 `balance`）選擇，會計入 `--model` 的每週區間。它們絕不選擇用量 100% 的帳號（Claude Code 回報已滿的區間在重設前都算作 100%）或 API 金鑰登入的帳號，而且只在起始帳號所屬的組織內選擇（加上 `--cross-org` 或設定 `cross_org = true` 可擴大範圍）。
- `--wait` 的值要寫在 `=` 之後（`--wait=2h`）；用空格分開時會被拒絕並附上提示。無法選出帳號時，它們會以代碼 75 結束，加上 `--wait` 時則等待第一個帳號重設。
- 即使帳號之間沒有共用對話記錄，`ccam next --continue` 也能把對話帶過去：此時 Claude Code 會把對話複製到新的帳號。加上 `--brief` 時，新帳號則會改為開始一段新的對話，內容來自 ccam 不靠模型為上一段對話寫的摘要（[`--brief`](/features/next#brief)）。
- 和 `ccam <name>` 不同，`ccam next` 會讀取自己的選項，所以 Claude Code 的選項要放在 `--` 後面；少了它時 ccam 會以 `unknown flag` 停止。
- `ccam agent` 會在這個資料夾中用 `claude -p`，以呼叫它的帳號以外的另一個帳號執行任務，沿用呼叫它的工作階段的權限模式和推理強度，並印出 Claude 的最後一則回覆。該帳號達到上限時，對話會以另一個帳號繼續（[其他帳號上的子代理](/features/agent)）。
- 在 Claude Code 中，外掛的 `/ccam:switch [account]` 會把對話交接出去，而 `ccam hook` 就是它的 hook 所執行的指令。
- 在 `config.toml` 中設定 `auto_switch = true`（Linux 和 macOS；預設關閉）時，上限 hook 還會結束由 ccam 啟動的互動式工作階段，並在另一個還有餘裕的帳號接續它，但只會是同類的帳號：同一個組織的另一個席位，或另一個個人方案（[`auto_switch`](/features/next#auto-switch)）。

請參閱[切換到下一個帳號](/features/next)、[為指令碼選擇帳號](/features/pick)、[其他帳號上的子代理](/features/agent)和[所有帳號的背景工作階段](/features/agents)。

## 維護

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

`ccam update` 會在取代 ccam 之前檢查發行版本的 SHA-256 總和檢查碼，`--version v0.2.0` 則可安裝特定的發行版本。從 `.deb` 套件安裝的 ccam 副本則改用 apt 更新。請參閱[更新](/docs/install#update)。

`ccam uninstall` 會先列出要移除的內容並詢問你：像 `ccam remove` 那樣移除每個帳號，以及快捷指令、ccam 的資料夾和 `ccam` 指令。匯入的資料夾會留在原處。請參閱[解除安裝](/docs/install#uninstall)。

## 檔案

| 路徑 | 內容 |
| --- | --- |
| `~/.config/ccam/config.toml` | 帳號、群組和設定（`theme`、`usage_source`、`relay`、`auto_switch`、`cross_org` …）。設定要放在最上方，位於第一個 `[[group]]`、`[[profile]]` 或 `[[folder]]` 之前：放在表格下方的設定會屬於該表格，因此 ccam 會在每次啟動時和 `ccam doctor` 中對它發出警告，而會儲存設定的指令在它被移走之前都會拒絕執行。ccam 儲存檔案時不會保留註解 |
| `~/.local/share/ccam/profiles/<name>/` | 每個帳號的 `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | 共用的設定、MCP 伺服器、擴充功能和對話記錄 |
| `~/.local/state/ccam/` | 用量快取、上限記錄、鎖定、記錄和其他狀態（請見表格下方的清單） |
| `~/.local/bin/<shortcut>` | `ccwork` 這類快捷指令，也就是指向 ccam 的連結 |

狀態資料夾包含：

- 用量快取與上限記錄（`usage/<name>.json`、`usage/<name>.limit.json`）；
- 上次使用時間；
- 工作階段標記與鎖定（`locks/`）；
- ccam 建立的快捷指令（`shortcuts.json`）；
- 尚未完成的群組變更（`group-change.json`）；
- `/ccam:switch` 和上限 hook 留下的記錄（`handoff/<account>/<session>.json`）；
- `ccam pick` 或 `auto_switch` 上次選中每個帳號的時間（`picks/<name>`；一次選擇在 2 分鐘內有效）；
- `--brief` 的摘要，保留 7 天（`brief/<session>.md`）。

設定 `XDG_CONFIG_HOME`、`XDG_DATA_HOME` 或 `XDG_STATE_HOME` 可以變更這些資料夾的位置。在 Windows 上，`~` 是 `%USERPROFILE%`，快捷指令則是 `<shortcut>.exe`。
