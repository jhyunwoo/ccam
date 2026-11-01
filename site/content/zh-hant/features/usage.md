---
title: 用量限制
description: ccam usage 一次顯示每個帳號的 5 小時和每週限制，ccam statusline 則在 Claude Code 中顯示帳號名稱和用量。
updated: 2026-11-01
order: 7
translation: machine
sourceHash: "54b296f7da35"
---

`ccam usage` 會顯示每個帳號在各項限制中已使用了多少：5 小時限制、每週限制，以及各模型自己的每週限制（例如 Fable）。每項限制都有一個長條，並附上距離重設的時間。在開始一項長時間的工作之前，或 Claude Code 告訴你快要達到限制時，可以用它查看哪個帳號還有剩餘額度。`ccam statusline` 是 Claude Code 的狀態列指令：它會在 Claude Code 中顯示帳號名稱和用量，並為 `ccam usage` 儲存這些數字。若想查看用量背後的 token 數和預估費用，`ccam ccusage` 會對每個帳號的對話執行 [ccusage](https://ccusage.com)（[見下文](#tokens-and-costs-ccam-ccusage)）。

## 用法

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts that are not running
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

每個帳號先有一行，顯示編號、名稱和數字的新舊程度；`(cached)` 代表數字來自已儲存的檔案，而不是端點。下面每項限制佔一行，顯示 20 格的長條、百分比和距離區間重設的時間：`5-hour`、`Weekly`（所有模型），接著是各模型自己的每週區間，例如 `Fable`。帳號沒有的區間會顯示 `-`。百分比在 80% 時變成黃色，在 100% 時變成紅色。

如果 Claude Code 因為某個區間已滿而拒絕了請求，該區間在重設前都會顯示 100%，不論已儲存的數字是多少（請參閱[上限記錄](#the-limit-record)）。在終端機中，所有帳號會先一起顯示為 `fetching…`，並在各帳號回應時繪製它的各行。

最後一行只在你有多個帳號時才會出現，它會列出 [`ccam next --auto`](/features/next) 會選擇的帳號。它遵循相同的規則：

- `auto_max_usage` 和 `auto_strategy`；
- 執行中的工作階段；
- 除非設定了 `cross_org`，否則只考慮上次使用的帳號所屬組織的帳號；
- 任何區間達到 100% 的帳號都會略過，所以每週額度已用完的帳號絕不會被列出。

使用預設的 `headroom` 策略時，這一行為 `→ Most headroom: <name>`；使用 `expiring` 或 `balance` 時則會說明原因，例如 `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`。

**Team 和 Enterprise 席位。** 自動選擇帳號（`ccam next`、`--fallback`、`ccam pick`、儀表板的 `n`，以及你開啟後的 [`auto_switch`](/features/next#auto-switch)）也會在同一個 Claude Team 或 Enterprise 組織的席位之間切換，這可能會違反 Anthropic 的條款。在依賴這項功能之前，請先閱讀[常見問題](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，並向你所屬的組織和 Anthropic 確認。

未登入的帳號會顯示 `not logged in → ccam login <name>`，API 金鑰帳號則顯示 `no limits (API)`。存取權杖過期時會顯示 ``access token expired; `ccam <name>` renews it when it starts``，如果有已儲存的數字，則顯示在它們上方：Claude Code 會在下次啟動時更新它。登入已完全過期（因為更新權杖過期，或伺服器拒絕後 Claude Code 已清除它）時，會顯示 `` login expired → `ccam login <name>` ``。

執行 `ccam statusline install` 之後，Claude Code 會顯示 `[work] · Opus · 5h 24% · 7d 41%` 這樣的一行。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--json` | 以 JSON 輸出報告：`name`、`slot`、`status`、`limits`（帳號開啟用量點數時含有 `credits`）、`source`、`at`、`error`，以及上限記錄把區間設為 100% 時的 `limited` 和 `limited_windows`。沒有任何帳號時會印出 `[]` | 關閉 |
| `--watch[=N]` | 每 `N` 秒重新繪製，直到你按下 Ctrl+C。只寫 `--watch`，或值為 0 以下時，代表 60。值要寫在 `=` 之後：`--watch 10` 會被拒絕並附上提示 | 關閉 |
| `-v`、`--verbose` | 在每個帳號下加上 `source:` 一行，說明數字的來源（用量端點、Claude Code、狀態列、已儲存的檔案，或 Claude Code 回報的用量上限）。當上限記錄覆蓋這些數字時，接著會有 `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` 這樣的一行 | 關閉 |
| `--no-cache` | 忽略 60 秒快取，重新查詢端點 | 關閉 |
| `--source live\|claude\|cache` | `live` 以帳號的權杖查詢端點；`claude` 讓 Claude Code 查詢（見下文）；`cache` 只顯示已儲存的數字 | `usage_source`，否則為 `live` |
| `--refresh` | 更新過期的存取權杖，只針對沒有在執行的帳號。一次更新可能需要最多 30 秒 | 關閉 |
| `--dry-run` | 接受：`ccam usage` 除了自己的快取以外不變更任何東西。和 `--refresh` 一起使用時會被拒絕：`--dry-run cannot be combined with --refresh, which rewrites login tokens` | 關閉 |

其他任何 `--source` 值都會失敗並顯示 `--source must be live, claude or cache`。`ccam usage` 不接受參數。

### 由 Claude Code 本身查詢用量

在 `config.toml` 的最上方（第一個 `[[group]]`、`[[profile]]` 或 `[[folder]]` 之前）設定 `usage_source = "claude"`，或使用 `--source claude` 時，ccam 完全不會讀取登入權杖。它會在每個帳號的資料夾中啟動 `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence`，並透過標準輸入送出 Claude Code 的 `get_usage` 請求（需要 Claude Code 2.1.270 或更新版本）。Claude Code 會用自己的登入查詢用量，並回傳 5 小時、每週和各模型的區間。

`--safe-mode` 會讓帳號的 hook、MCP 伺服器和外掛不啟動，也不會儲存任何工作階段。每個帳號都需要啟動一次 Claude Code（大約一秒），所以這個來源比 `live` 慢。所有帳號會同時查詢，每個的逾時為 20 秒，60 秒快取也照常適用。

存取權杖已過期或將在 5 分鐘內過期的登入會有 50 秒，因為 Claude Code 會在啟動時更新它，而被中斷的更新可能會讓登入遺失。如果 ccam 放棄某次查詢，它會以 SIGTERM 要求 Claude Code 結束，讓 Claude Code 完成已經開始的更新（Windows 沒有 SIGTERM，所以在那裡會直接終止該行程）。

完全過期的登入會顯示 `` login expired → `ccam login <name>` ``，而不啟動 Claude Code。Claude Code 把這個請求標示為實驗性功能，所以如果它不再回應，ccam 會顯示已儲存的數字，`ccam doctor` 也會回報。當 Claude Code 回答它自己的查詢失敗時，ccam 會顯示已儲存的數字並附上 `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`，而不是顯示一個沒有限制的方案。

使用 `--dry-run` 時，ccam 不會啟動 Claude Code，因為 Claude Code 會寫入帳號的資料夾；它會改為顯示已儲存的數字。`usage_source` 接受 `live` 或 `claude`；`cache` 只能作為 `--source` 的選項。

`ccam statusline` 從標準輸入讀取 Claude Code 的狀態列 JSON。它唯一的選項是 `--wrap <command>`：一個現有的狀態列指令，ccam 會在它輸出的第一行前面加上帳號名稱。`ccam statusline install` 沒有選項。兩者都會以 `` --dry-run is not supported by `ccam statusline`; nothing was changed `` 拒絕 `--dry-run`。

## Token 和費用：`ccam ccusage`

`ccam usage` 會顯示一個帳號已使用了每項限制的多少。[ccusage](https://ccusage.com) 則是計算 Claude Code 對話記錄檔中的 token 數和預估費用。它單獨執行時只會在 `~/.claude` 和 `~/.config/claude`（或 `CLAUDE_CONFIG_DIR` 指定的資料夾）中尋找 Claude Code 的記錄檔，所以它唯一看得到的 ccam 帳號，是在 `~/.claude` 就地匯入的那個帳號。`ccam ccusage` 執行它時，會把 `CLAUDE_CONFIG_DIR` 設為所有有對話的帳號的資料夾，讓一份報告涵蓋全部帳號：

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam 會先讀取自己的選項：`-a`（`--account`）接受帳號名稱或編號，可以用逗號分隔，也可以重複指定這個選項；`--dry-run` 只會印出指令而不執行。從第一個不屬於這些選項的參數開始，所有內容都會照你輸入的樣子傳給 ccusage，只有透過 `bun x` 和 `npx` 執行時的 `--config` 例外（見下文）；`ccam ccusage -- --help` 會顯示 ccusage 本身的說明。如果已安裝 `ccusage`，ccam 會執行它，否則執行 `bun x ccusage`，再不行就執行 `npx ccusage@latest`；後兩者會從 npm 下載 ccusage。它們會在 ccam 自己的資料夾（`~/.local/state/ccam/ccusage`）中啟動，而不是在目前的資料夾，因此目前資料夾中的專案無法透過它的 `node_modules/.bin`、`package.json` 中的 `bin` 或 `.npmrc` 中的 registry 決定它們要執行什麼；`ccam --dry-run ccusage` 會顯示進入該資料夾的 `cd`。三者都無法使用時，ccam 會請你安裝 Node.js、Bun 或 ccusage。在 Linux 和 macOS 上，ccusage 會取代 ccam；在 Windows 上，ccam 會等待它結束，並以它的結束代碼結束。在 WSL 中，ccam 會略過 Windows 磁碟機上的 `ccusage`、`bun` 或 `npx`，因為它們會以 Windows 程式執行，拿不到帳號資料夾。

- 報告涵蓋資料夾中有對話（也就是有 `projects` 資料夾）的帳號。多個帳號透過群組共用的對話只會計算一次。
- 使用 `-a` 時，透過群組共用對話記錄的帳號會帶入整個群組的對話，ccam 也會說明是哪些帳號的對話：`ccam: work shares its conversation history with personal, so the report counts their conversations too`。
- `blocks` 會顯示 5 小時區間，而每個帳號都有自己的區間，所以查看這份報告時請像上面那樣只指定一個帳號。
- ccusage 也會回報它找到的其他程式設計工具（例如 Codex），除非你指定其中一個：`ccam ccusage claude daily` 只會統計 Claude Code。
- 由 ccam 執行時，ccusage 只會從 `./.ccusage` 和各帳號資料夾讀取它的設定（`ccusage.json`）。`--config` 可以指定另一個檔案：`ccam ccusage daily --config ~/.config/claude/ccusage.json`。由於 `bun x` 和 `npx` 會在 ccam 自己的資料夾中啟動，當 `./.ccusage/ccusage.json` 存在而你沒有指定檔案時，ccam 會用 `--config` 把這個檔案傳給它們，並把相對路徑的 `--config` 改成完整路徑再傳入。這時 ccusage 會以 ccam 的資料夾為基準讀取其他相對路徑（例如 `--pi-path`），所以這些路徑請寫成完整路徑。
- ccusage 會把路徑中含有逗號的資料夾當成兩個資料夾來讀取，所以 ccam 會拒絕它並指出是哪個帳號；請用 `-a` 把這個帳號排除在外。
- `ccusage` 是 ccam 的指令，所以名為 ccusage 的帳號要用 `ccam run ccusage` 啟動；`ccam doctor` 會提醒這一點。

## 運作方式

`ccam usage` 顯示的數字來自兩個地方：Anthropic 沒有公開文件的用量端點 `GET https://api.anthropic.com/api/oauth/usage`（以每個帳號自己的登入資訊呼叫），以及 Claude Code 交給狀態列指令的值。這些全都位於 `internal/usage`。

1. ccam 載入 `config.toml`，並為每個帳號建立一個 `usage.Account`（`cli.usageCmd`）：它的資料夾 `~/.local/share/ccam/profiles/<name>/`、它的工作階段標記（搭配 Claude Code 自己在帳號資料夾中的記錄，用來判斷是否正在執行），以及它是否為就地匯入。
2. ccam 同時查詢所有帳號，每個帳號一個 goroutine（`Fetcher.FetchAll`）。它們共用同一個 `http.Client`，所以連到端點的連線會被重複使用。建立連線最多 2 秒，每個帳號總共有 3 秒。用戶端不跟隨重新導向，因為重新導向會把權杖重新送到別的地方，而且最多只讀取 1 MB 的回應。
3. 每個帳號完成時，會重新繪製它的那一列。失敗或沒有及時回應的帳號會改用已儲存的數字，所以一個緩慢的帳號絕不會拖住其他帳號。

在圖中，每個帳號都在自己的 goroutine 中查詢，逾時的那個帳號改為顯示已儲存的數字：

<!-- motion:usage-fetch -->

`personal` 和 `work` 及時回應，顯示即時的數字。`side` 花了超過 3 秒，所以 ccam 改為顯示它已儲存的數字（`fallback()`）。接著 ccam 為每個帳號繪製一個區塊（`usage.Table()`），並選出最後一行列出的帳號（`usage.Best()`）。

對於單一帳號，ccam 依序進行以下檢查（`Fetcher.Fetch`）：

<!-- motion:usage-cache -->

1. 加上 `--source cache` 時，直接交給 `fallback()`。
2. 除非設定了 `--no-cache`，否則 `~/.local/state/ccam/usage/<name>.json` 中測量時間不到 60 秒（`CacheTTL`）的即時結果會被使用，並以現在的時間判斷：之後已經重設的區間會顯示 0%，日期晚於時鐘時間的檔案不算是新的。連續執行兩次 `ccam usage`，每個帳號只會花一次請求。
3. ccam 讀取 `.credentials.json`（或 macOS 鑰匙圈）而不變更它，只解碼 `claudeAiOauth` 鍵（`readCredentials()`）。沒有訂閱登入時代表 `not logged in`；如果 `.claude.json` 中有 API 金鑰，則為 `no limits (API)`。更新權杖已過期或已被 Claude Code 清除的登入，會是 `` login expired → `ccam login <name>` ``，而且不會查詢。
4. 過期的存取權杖不會被使用。沒有 `--refresh` 時，帳號會顯示已儲存的數字，並附上 ``access token expired; `ccam <name>` renews it when it starts``。
5. 否則 ccam 會呼叫端點。401 或 403 視為登入過期（`login expired: authentication failed (401)`），其他任何狀態或無法讀取的回應內容視為錯誤，兩者都會改用已儲存的數字。被 Ctrl+C 結束的查詢會顯示 `cancelled`。429 表示端點已收到這個帳號過多的請求：ccam 會寫入 `<name>.pause`，記錄可以再次查詢的時間（`Retry-After`，或 5 分鐘，最多一小時），以 `rate limited by the usage endpoint; asking again after 15:10` 顯示已儲存的數字，並且在那之前不會為該帳號呼叫端點，即使加上 `--no-cache` 也一樣。立即重試只會延長頻率限制。
6. 正常的回應會連同測量時間寫入 `<name>.json`。

對於已儲存的數字，ccam 會取用最新的已儲存查詢結果（`<name>.json`，或來自 Claude Code 查詢的 `<name>.claude.json`），而當 `<name>.statusline.json` 寫入的時間比它晚時，會把狀態列的 5 小時和每週區間覆蓋在它上面（`fallback()`）；狀態列沒有的模型區間則維持查詢時儲存的值。在顯示已儲存的數字之前，它會檢查每個時間區間的重設時間：在數字儲存之後已經重設的時間區間會顯示 0%，因為舊的百分比已不再成立。這就是上面的 `personal` 顯示 `0%` 且沒有重設時間的原因。

端點有兩種回應格式。較新的格式有一個 `limits` 陣列，包含 `session`（5 小時）、`weekly_all` 和 `weekly_scoped`（各模型）項目。較舊的格式有 `five_hour`、`seven_day`、`seven_day_opus` 和 `seven_day_sonnet` 物件。ccam 兩種都能讀取（`ParseLive`）：以 `limits` 陣列為準，舊物件補上它沒有涵蓋的部分。完全沒有限制資料的回應內容會被回報為格式已變更，`ccam doctor` 會告訴你這件事。

### 上限記錄

當 Claude Code 因用量上限而停止某一輪時，[ccam 外掛](/docs/install#claude-code-plugin)的 hook 會從對話中讀取被拒絕的請求：哪個區間已滿（5 小時、每週、Opus 或 Sonnet），以及它何時重設。它會把這些寫入 `~/.local/state/ccam/usage/<name>.limit.json`。ccam 會把這份記錄套用在它從每個來源得知的其他所有內容之上（`Fetch`）：記錄中的每個區間在重設前都算作 100%，不論已儲存的數字多舊或多低，所以 `ccam next --auto`、`ccam pick`、`--fallback` 和 [`auto_switch`](/features/next#auto-switch) 在那之前都不會再次選擇該帳號。

沒有其他已知資訊時，帳號顯示的就只有這份記錄（加上 `-v` 時為 `source: the usage limit Claude Code reported`）。無人看管的工作階段（`-p` 執行和背景工作階段）不會寫入記錄，組織設定的停用（例如用量用完，或群組的上限設為 $0）也不會，因為管理員隨時可以解除它。`ccam login` 登入另一個 Claude 帳號時會捨棄這份記錄；刪除該檔案也能清除它。

### 更新權杖

更新權杖會輪替：更新一次會得到新的更新權杖，並讓舊的失效。如果 ccam 在 Claude Code 以同一個帳號執行時更新了登入，那個工作階段就會被登出。所以 `ccam usage` 只讀取權杖、絕不寫入，除非你加上 `--refresh`。這時 ccam 會（`Fetcher.refresh`）：

1. 拒絕以 `ccam import` 就地登錄的帳號，因為直接執行的 `claude` 可能在 ccam 不知情的情況下使用那份登入；
2. 拒絕存放在 macOS 鑰匙圈中的登入資訊；
3. 取得鎖檔 `~/.local/state/ccam/locks/profile-<name>.refresh.lock`，讓兩個 ccam 行程絕不會同時更新同一個權杖；持有鎖期間，該帳號會顯示 `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)`。`ccam rename`、`remove` 和 `logout` 也會取得這個鎖，並在更新進行中時拒絕執行；
4. 在持有鎖的情況下，如果該帳號有工作階段正在執行就拒絕，不論它是如何啟動的（透過 ccam、直接執行 `claude`、IDE，或背景工作階段及其背後會自行更新登入的常駐程式），如果另一個行程已經更新了權杖就略過請求；
5. 當 Claude Code 在帳號資料夾中持有它自己的更新鎖（`.oauth_refresh.lock`、`<folder>.lock`）時，不送出任何請求：Claude Code 正在更新那份登入；
6. 把更新權杖送到權杖端點，並最多等待 30 秒的回應，這是 Claude Code 給自己的更新的時間。請求一旦送出就絕不中斷，即使按 Ctrl+C 也一樣：伺服器可能已經輪替了更新權杖，只有它的回應帶有新的登入。Ctrl+C 或掛斷會結束其他查詢，ccam 則會在回應儲存後結束（按 Ctrl+C 後為 130）；
7. 在 Claude Code 的寫入鎖（`.storage-write.lock`）之下儲存，而且只變更當時 `.credentials.json` 的權杖欄位：Claude Code 在這段期間寫入的內容（MCP 登入）會保留，Claude Code 在這段期間更新的登入會被保留，登出也絕不會被復原。當回應提供新更新權杖的有效期限時，會像 Claude Code 一樣儲存為 `refreshTokenExpiresAt`。

沒有成功的更新會在帳號那一行或其已儲存的數字上方說明原因：例如 `token not renewed: a session is running`，或者當請求已送出但其回應無法使用或儲存時，顯示 `login may be lost: …`，並建議你在 Claude Code 要求登入時執行 `ccam login <name>`。

### 狀態列收集器

Claude Code 會執行狀態列指令，並透過管線傳給它一份 JSON 文件。ccam 從中讀取模型名稱和 `rate_limits`：

<!-- motion:usage-statusline -->

1. ccam 從 `CCAM_PROFILE`（ccam 啟動 Claude Code 時會設定）找出帳號，否則就把 `CLAUDE_CONFIG_DIR` 和各帳號資料夾比對（`runStatusline`）。
2. ccam 讀取 `rate_limits.five_hour` 和 `rate_limits.seven_day`（`ParseStatusLine`）。API 金鑰帳號的 JSON 中沒有頻率限制，工作階段的第一個回應之前也沒有；這時不會儲存任何東西。
3. ccam 寫入 `~/.local/state/ccam/usage/<name>.statusline.json`（`SaveStatusLine`）。端點沒有回應時，`ccam usage` 會使用它。
4. ccam 印出 `[name] · model · 5h N% · 7d N%`。加上 `--wrap` 時，它會以相同的輸入透過 `sh -c` 執行你的指令，並以 Claude Code 顯示狀態列的方式印出它的輸出：每個非空白的行，第一行前面加上 `[name]`。失敗的指令只會顯示帳號名稱。在 Windows 上，它會用和 Claude Code 相同的方式執行指令：有安裝 Git Bash 時使用 Git Bash，否則使用 PowerShell。
5. 當 5 小時或每週用量達到 90% 以上，而且還有其他帳號時，這一行的結尾會加上 `→ ccam next --auto --continue`，例如 `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`。

被包裝的指令沒有自己的時間限制：Claude Code 給狀態列指令 10 分鐘，並在較新的更新開始時停止正在執行的那次，所以緩慢的指令會像沒有 ccam 時一樣顯示。當 Claude Code 結束 ccam 時，該指令也會結束（在 Linux 和 macOS 上，連同它啟動的每個行程）。

`ccam statusline install` 會把 `statusLine` 寫入每個帳號的 `settings.json`，並保留該物件的其他鍵，例如 `padding`、`refreshInterval` 和 `hideVimModeIndicator`。如果帳號已經有狀態列指令，ccam 會以 `--wrap '<old command>'` 保留它；在 Windows 上則為 `--wrap-base64 <base64>`，在兩種 shell 中都不需要加引號。如果狀態列已經在執行 ccam（來自其他路徑，或以較舊版本寫入時的引號方式），ccam 會取代它而不是包裝它，並保留它所包裝的指令。

在 Windows 上，ccam 的路徑會以正斜線和代表主資料夾的 `~` 寫入（`~/.local/bin/ccam.exe statusline`），Git Bash 和 PowerShell 都能理解。需要引號的路徑會依 Claude Code 使用的 shell 寫成完整形式：Git Bash 為 `'…'`，PowerShell 為 `& '…'`。透過群組共用的 `settings.json` 只會經由它的連結寫入一次，群組中的每個帳號都會取得這項設定。再次執行時，對已經使用 ccam 的檔案會印出 `(already set)`。

## 安全性和特殊情況

- 除非你加上 `--refresh`，否則 `ccam usage` 絕不變更權杖檔案。存取權杖過期時，用 `ccam <name>` 啟動該帳號一次，Claude Code 就會更新它。完全過期的登入則需要 `ccam login <name>`。
- `--refresh` 絕不碰正在執行的帳號（即使是不透過 ccam 啟動 Claude Code 的情況）、匯入的帳號或鑰匙圈中的登入資訊。被拒絕的更新會顯示已儲存的數字和原因（`token not renewed: …`）。
- 用量端點沒有公開的文件，可能會變更。變更時你仍會看到已儲存的數字，`ccam doctor` 也會回報格式的變更。
- 既沒有已儲存的數字也沒有回應時，該帳號的那一行會顯示 `lookup failed:` 和原因，例如加上 `--source cache` 時的 `no saved usage`。
- 快取檔案只存放百分比和重設時間，絕不存放權杖。

## 相關內容

- [用量限制與切換帳號](/docs/usage-limits)
- [切換到下一個帳號](/features/next)
- [Claude Code 達到用量限制時該怎麼辦](/guides/usage-limit-workflow)
- [指令參考](/docs/commands)
