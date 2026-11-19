---
title: 管理帳號
description: 用 ccam add、ccam login、ccam list、ccam rename、ccam remove 和 ccam import 新增、登入、列出、重新命名、移除和匯入 Claude Code 帳號。
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

每個 ccam 帳號都是一個獨立的 Claude Code 設定資料夾，有自己的登入資訊。用這些指令來新增帳號、讓帳號登入，以及列出、重新命名和移除帳號。它們也能登錄你已經有的資料夾。

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

凡是接受 `<name|number>` 的指令，都可以指定帳號名稱或它的槽位編號（`ccam list` 中的 `#` 欄）。

`ccam use` 只會變更哪個帳號是預設帳號（`config.toml` 中的 `default`）。

在終端機中不加任何選項執行 `ccam add`，會開啟設定精靈。它會依序經過以下步驟：

1. 帳號名稱，除非你已經指定；
2. 帳號的登入方式：現在以 Claude 登入、API 金鑰，或稍後再以 Claude 登入；
3. 要加入哪些群組，其中 `default` 已勾選（還沒有任何群組時，則顯示 `default` 群組將會共用的項目）；
4. `Create a shortcut command for work? [y/N]`；
5. 讓你確認的摘要。

`Esc` 會回到上一步。只要加上 `--group`、`--no-group`、`--shortcut` 或 `--no-login` 任一個選項，或是不在終端機中執行（例如從管線執行），就會略過設定精靈。這時必須指定名稱，而且除非 `--shortcut` 指定了名稱，否則不會建立快捷指令。

不加 `--group` 時，新帳號會加入 `default` 群組。第一個帳號會建立這個群組，共用 MCP 伺服器、設定、`CLAUDE.md` 和擴充功能；請參閱[用群組共用](/features/share)。在 Windows 上，只有在開發人員模式下才能連結檔案。沒有開發人員模式時，第一個帳號建立的 `default` 只包含 `mcp` 和 `extensions`，並會告訴你（`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`）；ccam 也會拒絕以 `--group` 加入共用設定、記憶或對話記錄的群組。

`ccam list` 會顯示槽位、名稱、快捷指令、電子郵件、群組和上次使用時間，預設帳號前面會標示 `*`。電子郵件是該資料夾目前登入的帳號（取自它的 `.claude.json`）；若與 ccam 記錄的不同，記錄的電子郵件會以 `(was …)` 接在後面：

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

加上 `--json` 時，每個帳號都是一個物件，包含 `slot`、`shortcut`、`name`、`email`（目前登入的電子郵件，與表格相同）、`default`、`dir`、`last_used`（帳號從未啟動過時省略）、`groups` 和 `share`。`share` 會把每個項目（`mcp`、`settings`、`memory`、`extensions`、`history`）對應到它所來自的群組，或 `"local"`。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `ccam add --no-login` | 建立帳號但不登入；之後再用 `ccam login <name>` 登入 | 關閉 |
| `ccam add --group <group>` | 加入這個群組，而不是 `default`。可重複指定；`default` 不存在時會建立它 | `default` |
| `ccam add --no-group` | 不加入任何群組：這個帳號不共用任何項目 | 關閉 |
| `ccam add --shortcut <command>` | 同時在 `~/.local/bin` 中建立這個快捷指令 | 無 |
| `ccam add --api --key-from <source>` | 使用 API 金鑰而不是登入的帳號；請參閱 [API 金鑰帳號](/features/api) | 關閉 |
| `ccam list --json` | 以 JSON 輸出帳號 | 關閉 |
| `ccam remove -y`、`--yes` | 不詢問直接移除 | 關閉 |
| `ccam logout -y`、`--yes` | 不詢問直接登出 | 關閉 |
| `ccam import --name <name>` | 你指定的資料夾要使用的帳號名稱；未指定資料夾時會被拒絕（`--name needs a folder: ccam import <folder> --name ops`） | 從資料夾名稱推導 |
| `ccam import -y`、`--yes` | 不詢問直接匯入找到的每個資料夾 | 關閉 |
| `--dry-run` | `ccam list` 和 `ccam completion` 不會變更任何東西，所以支援它。`add`、`login`、`logout`、`use`、`rename`、`remove` 和 `import` 會拒絕它，並顯示 `` --dry-run is not supported by `ccam add`; nothing was changed ``（其中為該指令的名稱） | 關閉 |

`ccam login`（除了 `--expired`）、`ccam use` 和 `ccam rename` 沒有自己的選項。`ccam completion` 接受 `bash`、`zsh` 或 `powershell`；其他值會失敗並顯示 `only bash, zsh and powershell are supported`。

帳號名稱為 1 到 32 個字元，可使用 `a-z`、`0-9`、`-` 和 `_`，並以字母或數字開頭。名稱不能全是數字，因為 `ccam 2` 代表槽位 2。名稱也不能是 ccam 的指令（`add`、`list`、`ls`、`run`、`use`、`next`、`usage`、`u`、`share`、`mcp`、`login`、`logout`、`rename`、`remove`、`rm`、`import`、`migrate`、`doctor`、`shims`、`shortcut`、`completion`、`help`、`version`、`statusline`、`update`、`uninstall`、`folder`、`pick`、`ext`、`group`、`groups`、`hook`、`agent`、`agents`、`attach`、`ccusage`），所以 `ccam <name>` 永遠不會有歧義。這份清單是 `internal/config` 中的 `config.Reserved`。

## 運作方式

新增帳號時，會依序建立它的資料夾、加入它的群組、儲存設定、登入，以及建立快捷指令。執行 `ccam add work --shortcut ccwork` 時，ccam 會採取以下步驟（`internal/cli/profiles.go` 中的 `addProfile`）：

<!-- motion:accounts-add -->

1. 在建立任何東西之前，ccam 會檢查名稱（`config.ValidateName`）和快捷指令（`shim.Check`）。
2. ccam 以權限 `0700` 建立 `~/.local/share/ccam/profiles/work/`，並給帳號最小的可用槽位編號，因此已移除帳號的編號會被重複使用（`profile.Add`、`Config.NextSlot`）。第一個帳號會成為預設帳號。
3. 帳號加入它的群組（除非 `--group` 或 `--no-group` 另有指定，否則為 `default`；第一個帳號會建立 `default`），接著 ccam 把新建的空資料夾連結到這些群組，例如把 `settings.json` 連結到 `groups/default/settings.json`（`applyShare`，先用 `share.Build` 再用 `Apply`）。如果這一步失敗，新資料夾會被移除，因為殘留的資料夾會讓之後每次重試都因 `folder already exists` 而失敗。
4. 帳號會透過 `fsutil.AtomicWrite` 儲存到 `~/.config/ccam/config.toml`：先寫入暫存檔、`fsync`，再 `rename`，所以當機時絕不會留下只寫了一半的設定。
5. ccam 會在 `CLAUDE_CONFIG_DIR` 設為帳號資料夾的狀態下啟動 `claude auth login`，並以 `--email` 帶入帳號上次使用的電子郵件，讓登入頁面自動填好（`runLogin`）。如果登入失敗，帳號仍會保留，ccam 會印出錯誤，並告訴你之後執行 `ccam login work`。
6. ccam 從 `profiles/work/.claude.json` 讀取 `oauthAccount`（`profile.ReadIdentity`）：電子郵件、帳號和組織的 ID、組織名稱以及方案（`claude_max`、`claude_team` …）。這個檔案可能有好幾 MB，所以 ccam 只掃描那一個鍵，而不解析整個檔案（`TopLevelValue`）。ccam 會把這些資訊保存在 `config.toml` 的 `[profile.login]` 中。
7. ccam 把 `~/.local/bin/ccwork` 連結到 ccam 執行檔（ccam 啟動時會在 `config.toml` 中查詢 `argv[0]`），並把這個連結記錄在 `~/.local/state/ccam/shortcuts.json` 中（`shim.Install`）。如果無法建立連結，帳號仍會新增，ccam 會說明原因（`Could not create shortcuts: …`），摘要則顯示 `shortcut: ccwork (not created; see above)`。

登入可能要花上幾分鐘，這段期間另一個 ccam 行程可能會儲存設定。因此需要等待的指令會重新載入 `config.toml`，套用那一項變更（這裡是電子郵件），再儲存（`env.update`）。基於同樣的理由，`ccam remove` 會在詢問之後重新載入。

`ccam login` 會對現有的帳號執行第 5 和第 6 步。Claude Code 只把電子郵件當作登入頁面的提示，所以你登入的可能是和之前不同的 Claude 帳號；ccam 會把新的登入和 `[profile.login]` 比對，並告訴你（`ccam: this login is other@company.com, but work was dev@company.com before.`）。它會保留新的登入，並告訴你如何復原：``To undo it, run `ccam login work` again and sign in as dev@company.com.``

登入另一個 Claude 帳號時，也會捨棄這個帳號的[上限記錄](/features/usage#the-limit-record)，因為那些上限屬於另一個帳號。當另一個 ccam 帳號登入同一個 Claude 帳號時，ccam 會警告這兩者共用用量上限。

`ccam login --expired` 會逐一登入每個從未登入過、或登入已過期的帳號：更新權杖已過期，或伺服器拒絕後 Claude Code 已清除該登入。在開啟第一個瀏覽器之前，它會先確認 `config.toml` 可以儲存。

`ccam logout` 會在帳號中執行 `claude auth logout`，並清除儲存的電子郵件；登入記錄仍會保留，所以下次登入時會拿它來比對。Claude Code 會撤銷該登入，並刪除它為該帳號儲存的內容，包括 MCP 伺服器的登入。因此在以下情況，ccam 會先警告你並詢問（`--yes` 會略過詢問）：

- 登入中含有 MCP 登入；
- 帳號正在執行（它的工作階段也會被登出）；
- 另一個帳號的資料夾持有同一組登入（複製的資料夾），登出也會一併撤銷它。

例如：

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

若要以另一個 Claude 帳號登入，只要執行 `ccam login <name>` 即可，而且會保留 MCP 登入。

當帳號屬於不只一個組織時，`ccam list` 會加上 `Organization` 欄（個人方案顯示 `personal`），`ccam list --json` 也會為每個帳號提供 `plan` 和 `org`。

### 在遠端機器上登入

透過 SSH 連線時，`claude auth login` 在瀏覽器端的步驟可能會失敗。Claude Code 在遠端機器上開啟回呼伺服器（`http://localhost:<port>/callback`）等待登入，但完成登入的瀏覽器在你的電腦上執行，那個位址在你的電腦上連不到任何地方，瀏覽器會顯示 `ERR_CONNECTION_REFUSED` 之類的錯誤。這個位址仍然帶著一次性代碼，所以 ccam 可以替你轉交：

<!-- motion:accounts-relay -->

ccam 不會把你的終端機直接交給 Claude Code（`runLogin`）。它把 Claude Code 的標準輸入接到一個管線，並把你輸入的每一行複製進去（`relayLogin`），所以在 Claude Code 的 `Paste code here` 提示貼上的代碼仍然會送達。當某一行是回呼位址時，ccam 只在以下條件都成立時才接受它：主機是 `localhost`、`127.0.0.1` 或 `[::1]`，有指定連接埠，路徑是 `/callback`，而且 `code` 和 `state` 都存在；`http://` 可以省略（`callbackURL`）。接著 ccam 會像瀏覽器原本會做的那樣，在伺服器上請求該位址，逾時為 10 秒，而且不跟隨重新導向（`deliverCallback`）。這個請求會抵達正在等待的 `claude auth login` 的回呼伺服器（圖中的虛線）。Claude Code 收到代碼並儲存登入資訊，在 Linux 上是存到帳號資料夾中的 `.credentials.json`。貼上的內容絕不會讓 ccam 連線到其他主機。

ccam 在開始登入前會印出：

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

成功時會印出 `→ Passed the sign-in to Claude Code on this machine.`

登入在終端機上等待輸入時，ccam 會關閉終端機的回顯和行模式，你每輸入一個字元就顯示一個 `*`，所以貼上的代碼或回呼位址不會出現在畫面上。行編輯由 ccam 自己處理（`readLine`）：`Enter` 送出這一行，`Backspace` 刪除一個字元，在空行上按 `Ctrl+D` 結束輸入。登入結束時會還原終端機設定；收到 `Ctrl+C`、`SIGTERM` 或 `SIGHUP` 時也會還原，它們仍會結束 ccam。如果標準輸入和標準錯誤不是終端機，各行會照原樣複製。

### 重新命名、移除和匯入

`ccam rename` 和 `ccam remove` 會處理帳號的資料夾，以及 ccam 在其他地方以該帳號名稱保存的內容。ccam 在 `nameState` 中列出這些路徑：

- `~/.local/state/ccam/last-used/<name>`；
- `~/.local/state/ccam/mcp/<name>`；
- 用量快取 `~/.local/state/ccam/usage/<name>.json` 和 `<name>.statusline.json`；
- 上限記錄 `<name>.limit.json`；
- `~/.local/state/ccam/handoff/<name>/` 中的工作階段記錄。

`ccam remove work` 依以下順序處理它們：

<!-- motion:accounts-remove -->

1. 如果帳號正在執行，`ccam remove` 會先拒絕，否則會詢問你，並在你回答後再檢查一次（這樣也能發現在提示顯示期間啟動的工作階段）。
2. 在變更任何東西之前，它會先確認 `config.toml` 可以儲存。
3. 如果帳號共用擴充功能，ccam 會改寫群組的 `installed_plugins.json` 和 `known_marketplaces.json`（`RewritePluginPaths`），因為 Claude Code 儲存的是外掛的絕對路徑，而這些路徑可能經過這個帳號的資料夾。
4. 接著會將受管理的帳號（由 ccam 建立資料夾的帳號，而非以 `ccam import` 登錄的帳號）登出：在它的資料夾中執行 `claude auth logout`，讓垃圾桶中不會留下仍可使用的登入（在 macOS 上則是以資料夾命名的鑰匙圈項目，之後同名的帳號會繼承它）。API 金鑰帳號，以及另一個帳號的資料夾持有同一組登入時（登出也會撤銷它），會略過這一步。登出失敗時會回報，移除仍會繼續。
5. ccam 會把資料夾移到 `~/.local/share/ccam/trash/<name>-<time>/`（`profile.Trash`），帳號會從 `config.toml` 中移除，`nameState` 的路徑會被刪除，讓之後同名的帳號不會繼承它們，快捷指令的連結也會被移除：`✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`。

如果最後仍無法儲存 `config.toml`，資料夾會從垃圾桶移回來。若要復原已移除的帳號，把資料夾從垃圾桶移回原處，再重新登入即可。

`ccam rename` 會重新命名資料夾，並把其 `plugins/` 中的外掛路徑從舊資料夾改寫為新資料夾。它會把 `nameState` 的路徑移到新名稱下，如果該帳號原本是預設帳號，也會維持為預設帳號。它也會更新工作階段記錄，讓 `ccam next --continue` 仍然找得到它們。快捷指令的名稱不變。如果改寫外掛路徑失敗，或無法儲存 `config.toml`，資料夾會被改回原名。

在 macOS 上，以及在 Claude Code 把登入資訊存放在認證管理員中的 Windows 上，資料夾本身會保留原本的名稱，因為登入資訊和 MCP 登入是以資料夾的路徑保存的，移動資料夾會讓它被登出。此時設定會在新名稱旁把資料夾記錄為 `folder`，ccam 會顯示 `its folder stays …: the macOS Keychain keeps the login under the folder's path`。

`ccam import` 會在原位置登錄資料夾，記錄為 `config.toml` 中的 `dir`。它絕不移動或複製資料夾：複製後的更新權杖會讓各個副本互相登出。

不指定資料夾時，它會尋找尚未登錄的 `~/.claude` 和 `~/.claude-*` 資料夾，並逐一詢問。名稱取自資料夾（`~/.claude` 為 `main`，`~/.claude-work` 為 `work`，`~/.claude-2` 為 `acct2`）。如果 ccam 無法推導出名稱，或名稱已被使用，就會列出該資料夾以及應執行的 `--name` 指令。找到但失敗的資料夾會以 `Skipped <path>: <reason>` 回報，其他資料夾照樣匯入。

ccam 會先問完所有問題，再重新載入 `config.toml` 並登錄這些資料夾，因此其他 ccam 指令在這段期間儲存的變更會保留。`✓` 行會在儲存之後才出現。

沒有 `CLAUDE_CONFIG_DIR` 時，Claude Code 會把 `.claude.json` 放在家目錄中 `~/.claude` 的旁邊，所以當你匯入 `~/.claude` 而它沒有自己的 `.claude.json` 時，ccam 會建立指向 `~/.claude.json` 的連結 `~/.claude/.claude.json`。這樣直接執行 `claude` 和 `ccam main` 會看到同一個帳號。

如果 `~/.claude` 在 `~/.claude.json` 之外還有自己的 `.claude.json`，ccam 會保留兩者並告訴你。`ccam main` 會使用內層的檔案，直接執行 `claude` 則使用外層的檔案，兩者的登入詳細資料和 MCP 伺服器會彼此分開。若只想保留一份，請把內層檔案換成連結（[doctor](/features/doctor) 也會回報這一點）。

當該資料夾持有與某個已登錄帳號相同的登入時（複製的資料夾），ccam 會告訴你：先更新它的一方會讓另一方登出，所以使用前請先執行 `ccam login <name>`。

匯入的帳號在你把它加入某個群組（`ccam group join <group> <name>`）之前不屬於任何群組，也不會建立快捷指令。

`ccam completion bash|zsh` 會印出一份自動完成腳本，也能自動完成帳號名稱。

## 安全性和特殊情況

- 帳號執行時，`ccam rename` 和 `ccam remove` 會拒絕（`refuseRunning`），並顯示 ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``，因為工作階段會持續寫入該資料夾。這項檢查涵蓋以任何方式啟動的 Claude Code：透過 ccam、直接執行 `claude`、IDE、背景工作階段及其常駐程式。
- 當另一個 ccam 行程正在更新或變更該登入時，`ccam rename`、`ccam remove` 和 `ccam logout` 會拒絕：``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``。
- 當另一個 ccam 行程正在變更群組時，`ccam rename` 和 `ccam remove` 會以 `another ccam process is changing the groups right now; try again when it is done` 拒絕；若上一次群組變更中途停止，則顯示 ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``。
- 對匯入的帳號執行 `ccam remove` 會保留資料夾及其登入，供直接執行 `claude` 使用。如果帳號在群組中，它會先像 `ccam group leave` 一樣離開群組：每個共用的項目都會變成資料夾中群組內容的副本，群組的 MCP 伺服器則寫入它的 `.claude.json`。接著才取消登錄：`✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`。ccam 為它的 `.claude.json` 所做的備份會移到 `<folder>/.ccam-backup/mcp-<time>`；如果無法移動，備份會留在 `~/.local/state/ccam/mcp/<name>/backups`，訊息中也會說明。你的 `~/.claude` 不會被移到垃圾桶。
- 兩個帳號不能使用同一個資料夾：`config.toml` 中同一個資料夾出現兩次時會載入失敗。
- 如果 `profiles/` 中已有同名資料夾，`ccam add` 會拒絕該名稱，並提示你使用 `ccam import`。在 macOS 上，若鑰匙圈中仍有舊版本的 ccam 為新資料夾路徑留下的登入或 MCP 登入，它也會停止，並列出可移除它們的 `security delete-generic-password` 指令。如果無法儲存 `config.toml`，新資料夾會再被移除。
- `ccam import` 會拒絕已登錄的資料夾、不是資料夾的路徑，以及已被使用的名稱。當資料夾是 `~/.claude` 的副本（同一組登入）時，即使 `~/.claude` 本身尚未登錄，ccam 也會發出警告。
- 帳號正在執行時，`ccam logout` 會先詢問；`ccam login` 則不會檢查。
- 如果轉交無法送達 Claude Code（例如登入已經結束），會印出 `ccam: could not pass the sign-in to Claude Code:` 和原因。HTTP 錯誤通常表示代碼已過期；請重新開始登入。

## 相關內容

- [啟動 Claude Code](/features/launch)
- [快捷指令](/features/shortcuts)
- [共用 MCP 伺服器、設定和對話記錄](/docs/sharing)
- [在一台電腦上使用多個 Claude Code 帳號](/guides/multiple-claude-code-accounts)
