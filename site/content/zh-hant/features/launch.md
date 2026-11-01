---
title: 啟動 Claude Code
description: 用 ccam <name|number>、ccam run 或儀表板以任一帳號啟動 Claude Code；參數原封不動地傳給 claude，--dry-run 會顯示要執行的指令。
updated: 2026-11-01
order: 1
translation: machine
sourceHash: "11df08608a1e"
---

`ccam work` 會以 `work` 帳號啟動 Claude Code，名稱後面的所有內容都原封不動地傳給 `claude`。凡是你會輸入 `claude` 的地方都可以使用它：在終端機中、在腳本中，或搭配管線輸入。

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

指定帳號的名稱，或它在 `ccam list` 中的編號。ccam 不會讀取名稱後面的參數，所以 `ccam work --help` 顯示的是 Claude Code 的說明。名稱後面的 `--` 可加可不加。標準輸入、結束代碼和訊號的行為都和 `claude` 相同。

當第一個參數以 `-` 開頭，而且不是 ccam 自己的旗標時，ccam 會替你選擇帳號，並把所有參數傳給 Claude Code。它會使用以 [`ccam folder`](/features/folder) 為目前資料夾設定的帳號，否則使用預設帳號：用 `ccam use` 設定的帳號；如果你沒有設定，則是編號最小的帳號。

### 儀表板

不帶參數的 `ccam` 在標準輸入和輸出都是終端機時會開啟[儀表板](/features/dashboard)，否則會印出 `ccam list`。移到某個帳號（用方向鍵，或輸入它的編號）再按 `Enter` 即可啟動它；`n` 會啟動剩餘額度最多的帳號，`o` 則讓你選擇 `--continue`、`--resume` 或其他參數。儀表板會在 Claude Code 啟動前關閉，所以終端機會先還原。

**Team 和 Enterprise 席位。** 自動選擇帳號（`ccam next`、`--fallback`、`ccam pick`、儀表板的 `n`，以及你開啟後的 [`auto_switch`](/features/next#auto-switch)）也會在同一個 Claude Team 或 Enterprise 組織的席位之間切換，這可能會違反 Anthropic 的條款。在依賴這項功能之前，請先閱讀[常見問題](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，並向你所屬的組織和 Anthropic 確認。

當 Shell 設定了 `CLAUDE_CODE_OAUTH_TOKEN`、`ANTHROPIC_API_KEY` 或 `ANTHROPIC_AUTH_TOKEN` 時，Claude Code 可能會使用那個金鑰或權杖，而不是帳號自己的登入，所以 ccam 會在啟動前於標準錯誤輸出中警告你（`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`）。它只會讀取這些變數；[API 金鑰帳號](/features/api)則會清除它們。

預設情況下，ccam 會以 Claude Code 取代自己，所以不會有其他東西繼續執行。在 `config.toml` 的最上方、第一個 `[[group]]`、`[[profile]]` 或 `[[folder]]` 之前設定 `relay = true`（或在 Linux 和 macOS 上設定 [`auto_switch = true`](/features/next#auto-switch)）時，ccam 則會以子行程啟動互動式工作階段，並在它旁邊持續執行。這讓 `/ccam:switch`，或達到用量上限時的 `auto_switch`，能自行在另一個帳號接續對話，並沿用工作階段的選項、權限模式、推理強度和代理；在 `/ccam:switch` 把對話移到另一個組織的帳號之前，ccam 會先詢問你（[詳細說明](/features/next#relay)）。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--dry-run` | 印出 ccam 會使用的環境變數和指令，然後結束 | 關閉 |
| `--fallback` | 用於 `-p`/`--print` 執行：帳號達到用量上限時，以剩餘額度最多的已登入帳號再次執行同一個指令（[詳細說明](/features/next#run-again-at-the-limit)） | 關閉 |
| `-v`、`--version` | 印出 `ccam <version>`（只能作為第一個參數） | |
| `-h`、`--help` | ccam 的說明，可用 `ccam --help` 或 `ccam run --help` | |

參數規則：

- ccam 只讀取帳號名稱前面的自身旗標。名稱後面的每個參數都屬於 Claude Code，包括 `--dry-run` 和 `--help`。
- `--dry-run` 可以寫成 `ccam --dry-run work`、`ccam run --dry-run work` 和 `ccam --dry-run`（用於該資料夾的帳號或預設帳號）。`--dry-run` 和 `--fallback` 的順序不拘。
- `ccam run <name|number> [claude args]` 是 `ccam <name>` 的完整寫法。不帶名稱的 `ccam run` 會使用該資料夾的帳號或預設帳號。
- 帳號名稱不能是 `add`、`list`、`run` 或 `usage` 這類 ccam 指令名稱，所以 `ccam <word>` 永遠不會有歧義。

試執行會印出一行可以直接貼到 shell（在 Windows 上則是 PowerShell）的指令。它的開頭是啟動時對環境所做的變更：一律包含 `CLAUDE_CONFIG_DIR` 和 `CCAM_PROFILE`，要移除的變數（例如繼承而來的 `CCAM_RELAY`）會顯示為 `env -u <name>`，[API 金鑰帳號](/features/api)則會顯示其端點、模型以及金鑰的來源（`ANTHROPIC_API_KEY="$WORK_KEY"`），絕不顯示金鑰本身：

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

PowerShell 的這一行會在之後把變數還原，就像 shell 中的 `VAR=x cmd` 一樣，並透過 `Start-Process` 把命令列交給 Claude Code，所以 Windows PowerShell 5.1 和 PowerShell 7 的執行方式相同。啟動時若把某個變數設為空值，下方會有一行 `#` 說明 PowerShell 7.5 之前的版本會改為移除這種變數。對於 npm 安裝的 `claude.cmd`，這一行會執行啟動時所用的 `cmd.exe` 指令（見下文），而啟動時會在那裡被拒絕的參數，也會讓試執行失敗。

## 運作方式

每次啟動都經過相同的步驟，從命令列一直到 `exec`：

<!-- motion:launch -->

每個方塊會在其步驟執行時亮起，下方的各行是 Claude Code 啟動時使用的環境變數和參數。

1. ccam 從 `cmd/ccam/main.go` 所呼叫的 `app.Main` 開始。如果 ccam 是以其他名稱啟動的，那就是[快捷指令](/features/shortcuts)，會直接交給啟動器。
2. ccam 在 cobra 介入之前檢查 `os.Args`（`app.parse`）。只有單獨的 `ccam`（儀表板）、ccam 指令名稱（`config.Reserved`）、`-h`/`--help` 和 `-v`/`--version` 會交給 `cli.Execute`。其他一切都是啟動。
3. ccam 讀取 `~/.config/ccam/config.toml`，這是它唯一解析的檔案（`launcher.Launch`）。
4. ccam 依名稱或編號找到帳號，去掉開頭的一個 `--`，並找出 `claude` 執行檔（`launcher.Prepare`、`ResolveClaude`）。沒有名稱時，它會取用該資料夾的帳號，以字串比對工作目錄和 `config.toml` 中的 `[[folder]]` 項目，在 macOS 和 Windows 上不分大小寫（`HereProfile`）。PATH 查詢的結果會連同檔案的修改時間和 `PATH` 的雜湊值快取在 `~/.local/state/ccam/claude-path`，所以 ccam 不會在每次啟動時搜尋 `PATH`。檔案變更或不存在，或 `PATH` 變更時，會重新查詢。`config.toml` 中的 `claude_bin` 路徑（會展開 `~/`）必須是可執行檔；否則啟動會在寫入任何東西之前以 `claude_bin X is not an executable file. Check claude_bin in config.toml` 停止。如果 `claude` 就是 ccam 本身（以 inode 比對），ccam 會拒絕，因為那樣會無限循環。
5. ccam 把帳號資料夾 `~/.local/share/ccam/profiles/<name>/` 和它的共用設定比對，過程中只呼叫 `lstat` 和 `readlink`，每個項目一到兩次系統呼叫（`checkLinks`、`share.Check`）。缺少的連結會當場建立。任何可能遺失資料的情況都會中止啟動。
6. ccam 加入共用的 MCP 伺服器（`applyMCP`；見下文）。
7. ccam 記錄工作階段（`session.Register`）。它會在 `~/.local/state/ccam/locks/profile-<name>/<pid>` 寫入一個以其 PID 命名、記錄行程啟動時間的標記，不論 `XDG_RUNTIME_DIR` 的值為何，讓每個 shell 都看到相同的標記。對於共用對話記錄的帳號，ccam 會先在 `history-<group>-<hash>`（hash 取自解開符號連結後的目前資料夾）中尋找正在執行的工作階段，並同樣在那裡登錄（`historyCheck`）。

   `exec` 會保留 PID，所以標記之後指的就是 Claude Code 行程本身。ccam 使用標記檔而不是持有鎖，是因為 Claude Code 啟動的行程會繼承鎖，並在工作階段結束後繼續持有它。無法寫入標記時不會中止啟動：ccam 會顯示 `ccam: could not mark the session as running (…); starting Claude Code anyway`。

   ccam 其他判斷帳號是否正在執行的檢查（重新命名、移除、群組變更、`doctor --fix`、MCP 同步）也會讀取 Claude Code 自己在帳號資料夾中的記錄（`sessions/<pid>.json`、`daemon.lock`），所以也看得到不透過 ccam 啟動的 Claude Code：直接執行 `claude`、IDE、背景工作階段及其常駐程式。
8. ccam 把終端機標題設為 `claude · <name>`（當標準錯誤是終端機時），在不 `fsync` 的情況下把啟動時間寫入 `~/.local/state/ccam/last-used/<name>`，然後呼叫 `syscall.Exec`（`Plan.Exec`）。

Claude Code 執行時，`CLAUDE_CONFIG_DIR` 會設為帳號資料夾，`CCAM_PROFILE` 會設為帳號名稱。ccam 會清除這些變數繼承而來的任何值，以及繼承而來的 `CCAM_RELAY`。在 macOS 上，對於用 `ccam import` 登錄的 `~/.claude` 資料夾，它還會把 `CLAUDE_SECURESTORAGE_CONFIG_DIR` 設為空值，讓該帳號和直接執行的 `claude` 使用同一個鑰匙圈登入；在 Windows 上，若 Claude Code 把登入資訊存放在認證管理員中，也會這樣做。

### 為什麼用 exec

ccam 會以 `claude` 取代自己的行程（`syscall.Exec`）。因此：

- ccam 不會留在記憶體中；
- 按鍵、訊號和結束代碼都直接屬於 Claude Code；
- 參數始終是陣列，絕不會變成 shell 字串。

Windows 沒有 `exec`。在那裡，ccam 會以相同的環境變數、參數和主控台，把 `claude.exe` 作為子行程啟動，等待它結束，並以它的結束代碼結束（`Plan.Exec`）。npm 安裝的 `claude.cmd` 會以完整路徑在 `cmd.exe` 中執行，每個參數都會加上引號，讓 `cmd.exe` 不會把其中任何一個當成自己的語法（`&`、`%PATH%`）；含有換行的參數無法這樣傳遞，因此會被拒絕，並建議你安裝原生的 Claude Code。

Ctrl+C 會直接從主控台送到 Claude Code，所以 ccam 會忽略它。工作階段標記帶有 ccam 的 PID，它的存活時間和 Claude Code 完全相同，並在 Claude Code 結束時移除。等待期間，ccam 不做任何其他事。

工作階段之後的工作，例如發現用 `claude mcp add` 新增的伺服器，會在下次啟動時進行：ccam 會把 `.claude.json` 的修改時間和 `~/.local/state/ccam/mcp/<name>/seen` 比對，只在檔案變更時才解析它。

### 命令列如何分派

三種命令列，在 `internal/app` 中走三條路徑：

<!-- motion:launch-route -->

快捷指令依 `argv[0]` 比對，`ccam <name>` 依第一個參數比對，兩者最後都是 `syscall.Exec`。只有 `ccam usage` 這類 ccam 指令才會進入 cobra。

### 共用的 MCP 伺服器

使用預設的 `flag` 策略時，ccam 會以參數形式把群組的 MCP 伺服器傳給 Claude Code，放在你的參數之前：

<!-- motion:launch-mcp -->

ccam 只在群組的 `mcp.json` 存在時才加入這個參數，並把它放在 `claude` 的正後方（`mcp.FlagArgs`、`Plan.InsertArgs`）。你自己的 `--mcp-config` 或 `--strict-mcp-config` 仍會疊加套用。它是單一參數 `--mcp-config=<path>`，因為 Claude Code 的 `--mcp-config` 接受多個值：如果寫成兩個參數，後面的提示詞也會被當成設定檔讀取。

ccam 會決定哪些地方不能加上這個旗標（`claudeVerb`）。它會像 Claude Code 一樣從參數中讀取 Claude Code 的子指令，並略過選項及其值。`remote-control` 或 `rc` 不會加上旗標（見下文），放在最前面的背景工作階段動詞（`attach`、`logs`、`stop`、`kill`、`respawn`、`rm`）也不會。`ccam work mcp list`、`get` 和 `remove` 同樣不會加上，因為 Claude Code 在那裡會拒絕它；ccam 會改為在標準錯誤中說明，群組的伺服器也會在該帳號的工作階段中載入：

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

其他所有子指令和每個工作階段都會保留這個旗標。

在 `config.toml` 中設定 `mcp_strategy = "sync"` 時，ccam 會改為把伺服器寫入帳號的 `.claude.json`；如果帳號已經在執行，就略過這一步。當群組的 `mcp.json` 和 `.claude.json` 自上次檢查以來都沒有變更時，啟動也會略過這項工作。如果帳號自己有一個和群組伺服器同名的伺服器，同步會保留帳號的伺服器，並告知一次。同步也會把從 `.claude.json` 中消失的副本放回去（`put back group default's MCP servers …`），而且只變更或移除它自己寫入、且之後沒有人改動過的副本。

### Remote Control

`ccam work rc`（或 `ccam work remote-control`）會以帳號 `work` 啟動 Claude Code 的 Remote Control。`claude remote-control` 會拒絕出現在動詞之前的 `--mcp-config`，在動詞之後也不接受，而且它的工作階段只從帳號的 `.claude.json` 讀取 MCP 伺服器。因此 flag 策略在這裡的運作方式不同：

<!-- motion:launch-rc -->

1. ccam 會找出作為 Claude Code 動詞的 `remote-control` 或 `rc`，即使它出現在選項及其值之後也一樣（`ccam work --model opus rc`），並且不加入旗標（`claudeVerb`、`mcp.FlagArgs`）。
2. ccam 會比較群組 `mcp.json` 的修改時間和帳號的 `sync-stamp`，並確認每個群組伺服器都仍在 `.claude.json` 中（`mcp.Synced`）。如果 `.claude.json` 已經含有群組目前的伺服器，或群組沒有伺服器，Remote Control 會直接啟動，不會詢問。
3. 否則，ccam 會列出工作階段將缺少的伺服器並詢問：

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. 回答 `y` 時，ccam 會用 `sync` 策略的方式複製它們（`mcp.Sync`）：取得 Claude Code 的鎖、先建立備份，並保留帳號自己的伺服器。回答 `N` 或沒有終端機時，Remote Control 會在沒有這些伺服器的情況下啟動，ccam 會提示 `mcp_strategy = "sync"`，它會在每次啟動前複製這些伺服器。

在伺服器被複製之前，ccam 每次啟動 Remote Control 都會詢問。複製之後，只有在群組的 `mcp.json` 變更，或某個已複製的伺服器從 `.claude.json` 中消失時，才會再次詢問。在之後的一般啟動中，`.claude.json` 中的副本優先於 `--mcp-config`，所以 `ccam mcp add`、`remove` 和 `import` 會更新其他帳號持有的副本（[共用的 MCP 伺服器](/features/mcp)）。

### 10 ms 的預算

啟動流程的預算是第 95 百分位 10 ms；實際上，ccam 在 Claude Code 啟動前大約增加 4 ms。啟動流程絕不建立 cobra 指令樹、不啟動 Bubble Tea，也不開啟網路連線；如果套件的 `init()` 函式總共花費超過 5 ms，`TestPackageInitIsCheap` 就會失敗。`scripts/perf-check.sh` 用 `hyperfine` 測量 `ccam --dry-run work`、`ccam work` 和一個快捷指令，超過 10 ms 就失敗；`internal/launcher` 中的 `BenchmarkLaunchPath` 則會以迴圈執行這些步驟。若要查看每個步驟花費的時間，請設定 `CCAM_TRACE=1`：

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## 安全性和特殊情況

- 未知的名稱或編號會以 `no account "<name>"` 停止，並提示你查看 `ccam list`。完全沒有帳號時，ccam 會請你用 `ccam add <name>` 新增一個。
- 如果某個共用路徑變成了一般檔案或指向錯誤的位置，ccam 會以 `stopped because of share link problems:` 停止，列出這些問題，並請你執行 `ccam doctor --fix`。它不會在啟動流程中移動資料。
- 如果另一個和這個帳號共用對話記錄（同一個群組）的 Claude Code 工作階段正在同一個資料夾中執行，ccam 會警告對話記錄可能混在一起，並詢問 `Run anyway? [y/N]`。資料夾會以解開符號連結後的路徑比對。沒有終端機時（標準輸入來自 `/dev/null` 也算沒有），或從背景工作執行時，它會顯示警告並繼續。新的 `-p` 執行會開始自己的對話，不會收到警告；`--version`、`--help` 以及不會開啟對話的子指令（例如 `mcp`、`auth` 或 `plugin`）則完全不檢查。
- `--dry-run` 會略過 MCP 同步、新伺服器通知、`mcp list` 的說明、對話記錄檢查、工作階段標記、標題和上次使用時間，也不會 exec。它仍會檢查共用連結，並像一般啟動一樣建立缺少的連結。
- 如果在 `work` 已經在執行時執行 `ccam work rc`，ccam 不會碰 `.claude.json`，會告訴你這一點，並在沒有群組伺服器的情況下啟動 Remote Control；工作階段結束後再執行一次即可複製它們。使用 `--dry-run` 時，它只會印出訊息，既不詢問也不複製。
- 使用 `mcp_strategy = "sync"` 時，同步失敗會印出 `ccam: MCP sync failed (launching anyway): …`，Claude Code 會以它原本就有的伺服器啟動。

## 相關內容

- [快捷指令](/features/shortcuts)
- [透過 ccam 使用 Remote Control](/guides/remote-control)
- [管理帳號](/features/accounts)
- [指令參考](/docs/commands)
- [共用 MCP 伺服器、設定和對話記錄](/docs/sharing)
- [CLAUDE_CONFIG_DIR 說明](/guides/claude-config-dir)
