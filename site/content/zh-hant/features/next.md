---
title: 切換到下一個帳號
description: ccam next 會以另一個帳號開啟 Claude Code；--auto 選擇剩餘額度最多的帳號，--continue 把對話帶過去，--fallback 會重新執行 -p 工作。
updated: 2026-11-03
order: 8
translation: machine
sourceHash: "89e8dea9c55d"
---

`ccam next` 不需要你指定名稱，就能以另一個帳號開啟 Claude Code。當某個帳號達到用量限制時使用它。切換的方式有四種：

- 單純的 `ccam next` 會選擇你上次使用的帳號的下一個帳號，`--auto` 會選擇剩餘額度最多的帳號。`--continue` 會在新帳號接續這個資料夾的上一段對話，無論兩個帳號是否共用對話記錄。
- 在 Claude Code 中，[ccam 外掛](#inside-claude-code)會在用量上限讓工作階段停止時通知你，`/ccam:switch` 則會把對話交接出去。
- 對於執行 `claude -p` 的指令碼，`ccam --fallback` 會自行完成切換（請參閱[達到上限時重新執行](#run-again-at-the-limit)）。
- 開啟 [`auto_switch`](#auto-switch)（Linux 和 macOS）時，ccam 也會移動被 Claude Code 因上限而停止的互動式工作階段。

**Team 和 Enterprise 席位。** 自動選擇帳號（`ccam next`、`--fallback`、`ccam pick`、儀表板的 `n`，以及你開啟後的 [`auto_switch`](/features/next#auto-switch)）也會在同一個 Claude Team 或 Enterprise 組織的席位之間切換，這可能會違反 Anthropic 的條款。在依賴這項功能之前，請先閱讀[常見問題](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，並向你所屬的組織和 Anthropic 確認。

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

ccam 會在標準錯誤上以 `→ <from> → <to>` 印出切換，`--auto` 選擇時還會附上理由，然後以和 `ccam <name>` 相同的方式用新帳號啟動 Claude Code。`<from>` 是你上次使用的帳號；加上 `--continue` 時，則是這段對話所屬的帳號。`--` 後面的參數會原封不動地傳給 `claude`。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--auto` | 依用量（見下文）選擇已登入的帳號，略過用量為 100% 的帳號 | 關閉 |
| `--max-usage <1-100>` | 搭配 `--auto`：優先選擇用量低於此值的帳號 | `auto_max_usage`，否則為 100 |
| `--strategy headroom\|expiring\|balance` | 搭配 `--auto`：5 小時用量最低、5 小時區間最早重設，或平均每天剩餘每週額度最多的帳號 | `auto_strategy`，否則為 `headroom` |
| `--cross-org` | 搭配 `--auto`：也選擇其他組織的帳號（見下文） | `cross_org`，否則為關閉 |
| `--wait[=<duration>]` | 搭配 `--auto`：沒有帳號還有額度時，等待第一個帳號重設，最多等這麼久。值要寫在 `=` 之後：`--wait 2h` 會被拒絕，並提示 `--wait takes its value after '=': --wait=2h` | 關閉；單獨指定時為 `5h` |
| `-c`、`--continue` | 在新帳號接續這個資料夾的上一段對話（見下文） | 關閉 |
| `--brief` | 搭配 `--continue`：以新帳號開始一段新的對話，內容來自 ccam 為上一段對話寫的摘要，而不是恢復整段對話（[見下文](#brief)） | 關閉 |
| `-- <claude args>` | 接在 ccam 自己的參數之後傳給 `claude` 的參數 | 無 |
| `--dry-run` | 支援：印出 `→ <from> → <to>` 和要執行的命令列，不啟動 Claude Code，也不記錄啟動 | 關閉 |

`--auto` 和 `--continue` 可以一起使用。使用 `--continue` 時，`--resume` 參數會放在你自己的參數之前。

Claude Code 自己的選項，例如 `--dangerously-skip-permissions` 或 `--model`，要放在 `--` 後面。`ccam <name>` 會把名稱後面的所有內容原封不動地傳下去，`ccam next` 則會自己讀取參數，所以少了 `--` 時會以 `ccam: unknown flag: --dangerously-skip-permissions` 停止。

## 運作方式

每次透過 ccam 啟動，都會把時間寫入 `~/.local/state/ccam/last-used/<name>`。`ccam next` 會讀取這些檔案，檢查哪些帳號已登入，然後做出選擇。

1. ccam 根據 `config.toml` 為每個帳號建立一個 `next.Candidate`（`cli.candidates`）：帳號、它上次啟動的時間（`launcher.LastUsed`），以及它是否有可用的登入資訊（`creds.Has`：`.credentials.json` 中，或在 macOS 上鑰匙圈中的 Claude 登入，而且沒有完全過期，也沒有被 Claude Code 清除）。改用 API 金鑰的資料夾（Console 登入，或以 `ccam add --api` 新增的帳號）沒有方案上限，絕不會被選中。
2. 沒有 `--auto` 時，ccam 會依槽位編號排序帳號，並找出最近一次啟動的帳號（`next.Rotate`）。那就是你要切換離開的帳號。如果還沒有任何帳號啟動過，則是槽位最小的帳號。
3. 從那裡開始依槽位順序往前走，從最後一個槽位繞回第一個，並傳回第一個已登入的帳號。
4. ccam 會完全像 `ccam <name>` 一樣準備並啟動所選的帳號（`launcher.LaunchWith`）：檢查帳號的連結、套用 MCP 共用、登錄工作階段，並透過 `syscall.Exec` 以 `claude` 取代 ccam。

走訪從上次使用的帳號開始，並在它之後的第一個已登入帳號停下：

<!-- motion:next-rotate -->

`work` 是最近一次啟動的帳號，所以走訪從它之後開始。`side` 沒有登入，因此被略過，最後選擇了 `lab`。如果 `lab` 也已登出，走訪就會繞回 `personal`（圖中的 `% len`）。輪替絕不會傳回你出發時的帳號：沒有其他已登入的帳號時，它會以錯誤停止。

### --auto

`--auto` 依用量而不是槽位順序決定。ccam 保留已登入的帳號，以和 `ccam usage` 相同的方式一次查詢它們所有的用量（`Fetcher.FetchAll`），再依下列規則選擇（`cli.pickAuto`、`usage.Policy.Choose()`）：

<!-- motion:next-auto -->

帳號依它的 5 小時和每週區間判斷；當 Claude Code 的參數中（`--` 後面）有 `--model` 時，也會依該模型的每週區間判斷（名稱中含有 `opus`、`sonnet`、`haiku` 或 `fable`）。ccam 把帳號分成三個等級（`Policy.Choose()`）：

1. **有餘裕**：每個區間都低於 `auto_max_usage`（`config.toml`，預設 100）。策略會在其中選擇：`headroom`（預設）選 5 小時用量最低的帳號，`expiring` 選 5 小時區間最早重設的帳號，讓它剩下的額度在失效前用掉；`balance` 選在每週區間重設前，平均每天剩餘每週額度最多的帳號，並把每個執行中的 ccam 工作階段計入該帳號，讓每週額度不會閒置，平行的工作階段也會分散開來。
2. **超過**：低於 100%，但達到或超過 `auto_max_usage`。只有在沒有帳號有餘裕時才會使用；5 小時用量最低的帳號優先。
3. **已滿**：有區間達到 100%，或 5 小時用量未知。絕不會被選中。

當 Claude Code 在某個區間拒絕了請求時，該區間在重設前都算作 100%：[ccam 外掛](#inside-claude-code)的 hook 會記錄它（`~/.local/state/ccam/usage/<name>.limit.json`），不論已儲存的數字是多少。因此剛達到上限的帳號在重設前不會再次被選中，即使它的數字已經過時也一樣。

適用於所有帳號的設定要放在 `config.toml` 的最上方，位於第一個 `[[group]]`、`[[profile]]` 或 `[[folder]]` 之前。放在其中某個表格下方的設定會屬於該表格：ccam 會在每次啟動時對它發出警告，而會儲存設定的指令在它被移走之前都會拒絕執行。

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**組織。** 自動選擇只會在和你要離開的帳號同類的帳號之間進行：同一個 Team 或 Enterprise 組織中的席位，或者個人方案（Pro、Max），所以公司帳號的工作不會在私人帳號上繼續，反之亦然。ccam 從每個帳號的登入得知它的組織（`ccam login` 會記錄；較舊的登入則從 `.claude.json` 讀取）；ccam 無法判斷其組織的帳號可以和任何帳號搭配。登入同一個 Claude 帳號的兩個資料夾共用它的用量上限，所以其中一個絕不會接替另一個。`--cross-org`，或在 `config.toml` 中設定 `cross_org = true`，會解除組織規則，但在 [`/ccam:switch`](#inside-claude-code) 之後，移到另一個組織仍要等你回答 yes。

無法選出帳號時，ccam 會以代碼 75（`EX_TEMPFAIL`）結束，讓指令碼能分辨「稍後再試」和失敗。加上 `--wait` 時，它會改為等到第一個已滿帳號的阻擋區間重設（重設時間後 30 秒，屆時端點會顯示新的區間），印出 `ccam: no account has room; waiting until 15:10`，然後再選擇一次；沒有任何帳號說明何時重設時，它會在 10 分鐘後再查看。如果等待結束前都沒有帳號空出來，它會以 75 結束。

查詢遵循 `ccam usage` 的規則：60 秒內取得的數字會重複使用，每個帳號的查詢在 3 秒後逾時，沒有回應或存取權杖已過期的帳號會以已儲存的數字判斷，而上限記錄會覆蓋這些數字。已儲存的時間區間如果已經過了重設時間，會視為 0%，所以在那之後已經重設的帳號不會被跳過。`--auto` 絕不更新權杖。和單純的輪替不同，如果你上次使用的帳號仍然剩餘最多額度，`--auto` 也可能選擇它。

### --continue

1. ccam 會在每個帳號中尋找這個資料夾的上一段對話：在 `projects/<folder>/` 中的 `<session>.jsonl` 檔案裡（資料夾名稱是把解開連結後的工作資料夾中字母和數字以外的每個字元都換成 `-`，和 Claude Code 的命名方式相同），選出含有最新提示或回覆的那一個。沒有提示或回覆的工作階段，例如只執行了 `/mcp`、`/usage` 或 `/status` 的工作階段，或空白檔案，會被略過。如果某個工作階段留下了記錄（用量上限或 `/ccam:switch`，見下文），而且沒有任何對話在它之後超過 2 分鐘還有提示或回覆，就改用那個工作階段；Claude Code 在 `/exit` 時寫入的資料列不算在內。
2. 它所屬的帳號就是你要切換離開的帳號，新帳號則依上述方式在同一個組織內選出。當多個帳號共用那份對話記錄時，它屬於 `.claude.json` 表示最後執行過該對話的帳號（`lastSessionId`），否則屬於最近使用的帳號。
3. 如果新帳號共用那份對話記錄（同一個群組資料夾），ccam 會啟動 `claude --resume <session id>`，也就是那段對話本身（`--continue` 則會開啟資料夾中最新的那段對話，不論是哪一段）。
4. 否則它會啟動 `claude --resume <the conversation's file> --fork-session`：Claude Code 會以新的工作階段 ID，把整段對話複製到新帳號自己的對話記錄中，原本的對話維持不變。之後新帳號的 `--continue` 就找得到它。
5. 除非你在 `--` 之後指定 `--agent`，否則 ccam 也會傳遞對話所用的代理：資料夾的 `.claude/settings.local.json` 或 `.claude/settings.json`、或帳號的 `settings.json` 中的 `agent` 設定，沒有時則是對話最後一筆 `agent-setting` 資料列指定的代理。Claude Code 在複製對話時（第 4 步）不會帶回代理，而會以它的預設代理、帶著所有工具繼續對話。當新帳號沒有共用舊帳號的代理、外掛和設定（也就是沒有透過群組使用相同的 `agents` 和 `plugins` 資料夾以及 `settings.json`）時，ccam 會在關於移動的那一行之後再加一行：`work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`。

ccam 會告訴你這段對話有多大：新帳號在第一次回覆時會重新讀取全部內容。提示快取絕不會在組織之間共用，而每個個人方案都是一個獨立的組織，所以當兩個帳號屬於不同的組織時，ccam 會加上 `prompt caches are not shared between organizations, so the first reply reads it all again`。當舊帳號已登入、新帳號尚未登入某些 MCP 伺服器時，ccam 會列出它們：`work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`。在群組中共用 `history` 可以避免複製、維持同一個對話，也會共用 `/rewind` 用來還原程式碼的檢查點（`file-history/`）。

### --brief

`--continue` 會恢復整段對話，新帳號在第一次回覆時會重新讀取全部內容。`ccam next --continue --brief` 則改為開始一段新的對話，內容來自 ccam 不詢問模型、自行為舊對話寫的摘要：

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

摘要依序包含以下內容，只列出有內容可寫的項目：

1. 對話的標題；
2. 對話記錄檔的路徑，讓 Claude 需要細節時可以搜尋；
3. 你的第一個請求（請求超過四個時）；
4. Claude Code 上次壓縮對話時寫的摘要；
5. 你最後四個請求，其中 `!` 指令和斜線指令保持你輸入時的樣子；
6. Claude 的最後一則回覆；
7. 待辦清單；
8. Claude 最後變更的 10 個檔案；
9. 資料夾的 `git status --short --branch`，最多 800 個字元。

待辦清單最多保留 10 個項目：先捨棄已完成的項目，再捨棄最後面的步驟，讓目前的步驟和接下來的步驟留下。過長的部分會從中間截斷，整份摘要不超過 16,000 個字元。

ccam 會把它存為 `~/.local/state/ccam/brief/<session>.md`，只有你能讀取；下一次 `--brief` 會刪除超過 7 天的摘要。Claude Code 會以一個提示詞啟動，並附加其中提到的檔案：`Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` 你在 `--` 後面提供的提示詞會取代最後一句（`Then: <your prompt>`），`--` 後面的其他選項則接在它後面。

Claude Code 的提及會在 `"` 處結束、在 `#` 處截斷，所以摘要的路徑含有其中任何一個字元時，ccam 會改用相對於工作資料夾的路徑提及它；如果那個路徑也含有這種字元，提示詞會直接寫出檔名，由 Claude 用它的 Read 工具讀取，而因為檔案位於工作資料夾之外，讀取前會先詢問（在 Manual 模式的 `-p` 執行中，Claude Code 會拒絕讀取，Claude 會在沒有摘要的情況下繼續）。

`ccam --dry-run next --continue --brief` 會印出摘要，但不儲存任何東西。舊對話維持原樣，新對話只知道摘要中寫的內容，所以當任務本身比你如何走到這一步的每個細節更重要時，再使用 `--brief`。沒有 `--continue` 時，ccam 會以 `--brief works with --continue` 停止；無法讀取對話或在其中找不到可摘要的內容時，它會停止，並說明不加 `--brief` 的 `ccam next --continue` 會完整恢復對話。[接力](#relay)和 [`auto_switch`](#auto-switch) 一律恢復整段對話。Claude Code 自己的 `--brief` 選項和它的其他選項一樣，要放在 `--` 後面。

## 在 Claude Code 中

[ccam 外掛](/docs/install#claude-code-plugin)會新增兩個 hook（Claude Code 2.1.78 或更新版本）：

- 當某一輪因用量上限而停止時，它會記下這段對話，並在支援的終端機（iTerm2、WezTerm、Windows Terminal、kitty、Ghostty）中顯示桌面通知：``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  它認得 Claude Code 視為用量上限的每一種訊息：5 小時、每週、一般和模型上限（`You've hit your session limit`、`You've hit your limit`，以及表示模型額度的 `You've reached your …`）、用量點數和支出上限，以及組織停用了該帳號。

  當 Claude Code 已把被拒絕的請求寫入對話時，hook 也會把已滿的區間記錄到它重設為止（[上限記錄](/features/usage#the-limit-record)），但組織設定的停用除外，因為管理員隨時可以解除它。短暫的伺服器節流或 API 金鑰的速率限制不是用量上限，不會處理；無人看管的工作階段（`-p` 執行和背景工作階段）也一樣。子代理或工作流程代理因用量上限而停止時，就算是它所屬工作階段的停止，即使主對話只是在等待它：Claude Code 會把那個被拒絕的請求寫入工作階段 `subagents` 資料夾中該代理自己的檔案，hook 就從那裡讀取。多個代理同時停止時，該工作階段只會得到一筆記錄和一則通知。

  Claude Code 2.1.234 及更新版本也會在工作階段中等待，並在 5 小時區間重設後自行繼續；通知讓你可以選擇改為切換帳號；而開啟 [`auto_switch`](#auto-switch)（Linux 和 macOS）時，若另一個帳號還有餘裕，ccam 會替你切換。當沒有其他帳號可以接手時，也就是沒有其他帳號屬於同一個 Team 或 Enterprise 組織，或沒有其他個人方案時（開啟 `cross_org` 時，任何其他 Claude 帳號都可以），通知的結尾會改為 `You have no other personal plans, so ccam cannot move this conversation.` 或 `You have no other accounts in Acme, so ccam cannot move this conversation.`
- `/ccam:switch [account]` 會記下這段對話，並在發出任何請求之前結束這一輪：輸入 `/exit`，接著執行 `ccam next --continue`，就會在你指定的帳號（或下一個帳號）接續它。在還沒有任何對話的工作階段中，它會說明沒有可以接續的內容。只有你親自輸入才會執行它：從外掛 0.3.2 起，Claude 無法執行 `/ccam:switch`，而 Claude 排定的工作（CronCreate、ScheduleWakeup 或 `.claude/scheduled_tasks.json`）會把它當成純文字交給 Claude。

  當你指定的帳號不在對話所屬的組織內（另一個 Team 或 Enterprise 組織的席位、從席位換到個人方案或反過來，或任一方是 API 金鑰帳號）時，不論 `cross_org` 怎麼設定，`ccam next --continue` 都會先詢問：`ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`。你沒有指定帳號時，ccam 在移到它自己挑選的帳號（`cross_org = true` 時可以來自另一個組織）之前也會同樣詢問：`ccam: the account picked for /ccam:switch is home, a personal plan, …`。只有 `y` 或 `yes` 才會移動；在 Linux 和 macOS 上，會先捨棄提問之前輸入的按鍵。任何其他回答，或沒有可以詢問的終端機，都會讓對話和記錄留在原處：``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work``（如果帳號是 ccam 挑選的，`ccam next --continue` 會重新挑選）。和自動選擇一樣，ccam 無法判斷其組織的帳號可以和任何帳號搭配。

記錄存放在 `~/.local/state/ccam/handoff/<account>/<session>.json`，每個工作階段一個檔案，所以同一個帳號的兩個工作階段絕不會覆寫彼此的記錄。它保留 12 小時，而且只會在下一個工作階段即將啟動時使用一次；啟動失敗時會保留下來，供下次嘗試。它屬於工作階段啟動時所在的資料夾，即使之後在 Claude Code 的 shell 中執行了 `cd` 也一樣。`ccam rename` 會更新其中的名稱；如果它指名的帳號後來被移除，ccam 會告知這件事，並照常選擇。

### 接力

在 `config.toml` 的最上方設定 `relay = true` 時，ccam 會在互動式工作階段旁持續執行，而不是以 Claude Code 取代自己，並在 `/ccam:switch` 之後自行接續那個工作階段的對話（其他工作階段留下的記錄不會被處理）：hook 會結束 Claude Code，ccam 再以該工作階段的記錄執行 `ccam next --continue`。沒有指定帳號時，它會選擇剩餘額度最多的帳號，並把對話所用的模型計算在內；如果因為 `cross_org = true` 而選中了另一個組織的帳號，它會在把對話移過去之前先問上面的問題。在 Windows 上，hook 無法結束 Claude Code：請輸入 `/exit`，ccam 就會接續下去。對於對話所屬組織之外的帳號，hook 也不會結束 Claude Code：它會顯示 `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`，等你輸入 `/exit` 之後，ccam 會問上面的問題。這樣每個工作階段都要多一個等待中的 ccam 行程，所以預設為關閉。

接力的工作階段會保留繼續對話時不會還原的內容：`--add-dir`、`--mcp-config`、工具清單、`--settings`、`--plugin-dir` 和系統提示詞選項等選項，對話的權限模式、工作階段的推理強度，以及對話所用的代理（見 [`--continue`](#continue)）。`--model` 和 `--name` 則由 Claude Code 自己從對話中還原。其他任何內容（`--worktree`、`--debug`、`--remote-control`、ccam 不認得的選項）都會列在 `ccam: not carried over: …` 中。切換時絕不會詢問共用對話記錄的問題，而在切換期間輸入的內容會被捨棄；它唯一會問的，是移到另一個組織的帳號之前的那個問題。

只有當 Claude Code 依 `/ccam:switch` 的要求結束時，接力才會接續對話：如果你先輸入了 `/exit`、在它結束時按了 Ctrl+C，或 ccam 收到掛斷或其他訊號，或終端機已經不在了，都不會接續。此時記錄會保留下來，之後可以用 `ccam next --continue` 接續對話。在較舊版本 ccam 的接力下啟動的工作階段，或無人看管的工作階段（`-p` 執行和背景工作階段），不會被 hook 結束，而是會收到 `/exit` 的訊息。`claude remote-control` 不使用接力執行。

### `auto_switch`

在 `config.toml` 的最上方設定 `auto_switch = true`（Linux 和 macOS）時，當 Claude Code 因用量上限而停止工作階段，ccam 會自行把它移到另一個帳號：外掛的 hook 會選出另一個還有餘裕的帳號並結束 Claude Code，而在工作階段旁持續執行的 ccam 會在那個帳號接續對話，就像[接力](#relay)在 `/ccam:switch` 之後所做的一樣。它預設關閉，屬於實驗性功能。

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

只有在以下條件全部成立時，ccam 才會移動工作階段：

1. 工作階段是 ccam 在終端機中啟動的，而且 ccam 在它旁邊持續執行。`auto_switch` 會開啟[接力](#relay)，所以設定之後你啟動的互動式工作階段都符合這一點。在那之前或在較舊 ccam 下啟動的工作階段，會收到結尾為 ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` 的通知。無人看管的工作階段（`-p` 執行和背景工作階段）、`claude remote-control`，以及不透過 ccam 啟動的 Claude Code 絕不會被移動。
2. 帳號已開啟 [ccam 外掛](/docs/install#claude-code-plugin)，並執行 Claude Code 2.1.214 或更新版本，這些版本會把自己的行程 ID 傳給 hook（`CLAUDE_PID`）。
3. Claude Code 已把被拒絕的請求連同重設時間寫入對話，而且針對的是該帳號的 5 小時、每週、Opus 或 Sonnet 上限。支出上限、用量點數用完、組織設定的停用、Fable 額度（`You've reached your Fable limit.`）、伺服器節流以及 API 金鑰的速率限制，絕不會讓工作階段移動。子代理或工作流程代理的上限也算在內，即使主對話只是在等待它；如果那是每週的 Opus 或 Sonnet 上限，只有當對話本身以該模型執行時，工作階段才會移動。
4. ccam 從帳號的 `.claude.json` 或 `ccam login` 的記錄知道該帳號的方案。若不知道，通知的結尾會是 `ccam does not know personal's plan, so auto_switch leaves this session.`
5. 工作階段沒有連接到 Remote Control。換成另一個帳號後，它不會重新連接到你用來操作它的應用程式，所以它會留在原處，通知也會說明這一點。
6. 上限在 15 分鐘以後才會重設。等待時間較短時，ccam 會把工作階段交給 Claude Code：`ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. 達到上限之後，你還沒有輸入任何提示詞。
8. 另一個帳號還有餘裕。ccam 的選擇方式和 [`--auto`](#auto) 相同（`auto_max_usage`、`auto_strategy`、對話所用的模型），並且只在和你離開的帳號同類的帳號之間選擇：同一個 Team 或 Enterprise 組織的另一個席位，或另一個個人方案（Pro、Max）。`cross_org` 和 `--cross-org` 在這裡不適用。該帳號必須：
   - 以 Claude Code 能更新的登入資訊登入；
   - 沒有登入同一個 Claude 帳號；
   - 沒有使用 API 金鑰；
   - 有 ccam 知道的方案；
   - 就 ccam 自己的用量查詢所見，沒有開啟用量點數；
   - 已在 Claude Code 中信任這個資料夾；
   - 如果對話以某個代理執行（`agent` 設定或對話指定了代理），要共用你離開的帳號的代理：它的 `agents` 和 `plugins` 資料夾以及 `settings.json` 都必須各自（透過群組）完全相同，或兩邊都沒有。

   當最後這三項排除了所有還有餘裕的帳號時，通知會指出它們：`side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`、`side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` 或 `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` 有[上限記錄](/features/usage#the-limit-record)的區間算作已滿，而過去 2 分鐘內已有其他工作階段移入的帳號會排在其他帳號之後。選擇必須在約 6 秒內完成；用量查詢花的時間更長時，工作階段會留在原處。設定 `usage_source = "claude"` 時，`auto_switch` 只依已儲存的數字決定，而且 ccam 無法從 Claude Code 的查詢或狀態列得知用量點數是否開啟，所以會假設它們是關閉的。

接著 hook 會結束 Claude Code，而在旁邊持續執行的 ccam 會在支援的終端機中顯示桌面通知，例如 `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` 它會印出上方範例輸出中的第一行，並以選中的帳號為該工作階段執行 `ccam next --continue`，由後者印出另外兩行，再沿用工作階段的選項、權限模式、推理強度和代理啟動 Claude Code，就像[接力](#relay)一樣。新工作階段的第一則訊息是一段固定文字，這也是 ccam 唯一代替你送出的內容（在 `/ccam:switch` 之後的移動不會送出任何內容）：

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

在子代理達到上限之後則是下面這段，因為子代理、工作流程和背景指令已隨舊的工作階段一起停止：

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

如果那個帳號也達到上限，同樣的事會在那裡再發生一次。上限記錄會防止這條鏈在重設之前回到先前的帳號，所以當所有帳號都已用滿時，工作階段會留在最後一個帳號。

hook 會從 `agent` 設定，或從 Claude Code 不時以及在工作階段結束時寫入對話、指明代理的那一筆資料列，得知對話以某個代理執行。如果到目前為止只有工作階段的 `--agent` 指明了代理，而選中的帳號沒有共用這些代理，`ccam next --continue` 就不會啟動任何東西，並保留記錄，第二則桌面通知會說明這件事：``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``。

沒有其他帳號還有餘裕時，什麼都不會移動：Claude Code 會照常處理上限，在工作階段中等待或提供它的選項，通知則會說明哪個帳號最先恢復，例如 `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` 如果是席位，則會顯示 `None of your other accounts in Acme has room now`。若要停止已經開始的移動，請在 Claude Code 關閉時按 Ctrl+C：``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

移動是有代價的。Claude Code 會重新啟動，這需要幾秒鐘，而 Claude 啟動的背景 shell、開發伺服器和代理會隨著舊工作階段一起停止。新帳號在第一次回覆時會重新讀取整段對話，和 [`--continue`](#continue) 一樣，ccam 也會告訴你對話有多大；移動絕不會使用 [`--brief`](#brief)。

`auto_switch` 絕不會在 Claude Code 回報上限之前切換，絕不會在背景監看你的用量，也絕不會複製、更新或轉交登入資訊。在 Windows 上它會被忽略，`ccam doctor` 也會說明這一點；在那裡，設定 `relay = true` 時，`/ccam:switch` 和 `/exit` 仍然可以把對話移到另一個帳號。

**Anthropic 的條款。** 在 Team 和 Enterprise 席位上，`auto_switch` 會把你的工作移到同一個組織的另一個席位。Anthropic 說明中心關於 [Team 方案](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)的文章寫道：「Each seat on a Team plan is assigned to one named individual, and only that person may use it」（Team 方案的每個席位都指派給一位具名的個人，而且只有這個人可以使用），以及「Users can't route their work through another user's account to use that user's remaining usage」（使用者不能經由其他使用者的帳號處理自己的工作，藉此使用對方剩餘的用量）；關於 [Enterprise 方案](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)的文章也是如此。

至於個人方案（Pro、Max），Anthropic 的條款並沒有說明你自己的帳號能否在達到上限時自動互相接手。[消費者條款](https://www.anthropic.com/legal/consumer-terms)（Consumer Terms）禁止「through automated or non-human means, whether through a bot, script, or otherwise」（透過自動化或非人為的方式，無論是藉由機器人、指令碼或其他方式）存取服務，除非 Anthropic 明確允許，而一次移動會代替你送出上面其中一則固定訊息。將於 2026 年 11 月 12 日生效的[使用政策](https://www.anthropic.com/legal/aup)（Usage Policy）把「Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products」（故意規避或干擾我們產品中設定的功能、限制或防護措施）列為禁止事項之一。

ccam 並不主張 `auto_switch` 符合這些條款。要不要開啟由你決定，風險也由你承擔：Anthropic 認為條款遭到違反時，可能會停用帳號，或停用 Team 或 Enterprise 組織的存取權。

開啟 `auto_switch` 時，`ccam doctor` 會對沒有 ccam 外掛的帳號和早於 2.1.214 的 Claude Code 版本發出警告，並在有席位時重述 Anthropic 的席位規則（[檢查項目](/features/doctor)）。

## 達到上限時重新執行

`ccam next` 用於已經停下來的互動式工作階段。非互動式的執行，也就是指令碼或排程工作中的 `claude -p`，無法等你處理，所以在帳號名稱前加上 `--fallback`，就能讓 ccam 替它切換：

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

1. `--fallback` 只能和 `-p` 或 `--print` 一起使用；否則 ccam 會停止，並指向 `ccam next --auto --continue`。
2. 管線輸入的標準輸入只讀取一次，最多 10 MiB，並提供給每次嘗試。ccam 和 Claude Code 一樣，最多等待 3 秒讓輸入開始（使用 `--input-format stream-json` 時除外）；如果管線保持開啟卻沒有傳送任何內容，這次執行就不會有輸入，ccam 會在標準錯誤中註明這一點。
3. ccam 會把 `claude` 作為子行程啟動，而不是使用 `exec`；它會暫存標準輸出，標準錯誤則即時傳遞。
4. 當 `claude` 以 0 結束，或因其他原因失敗時，暫存的輸出和結束代碼會原樣傳回。被訊號結束的執行（結束代碼大於 128）也一樣：它是被要求停止的，不會再次執行。
5. 當 Claude Code 的訊息表明是用量上限（[hook](#inside-claude-code) 認得的任何一種用量上限訊息，也包括 `--output-format json` 或 `stream-json` 輸出中的訊息），或該帳號的用量顯示 5 小時、每週或 `--model` 區間已達 100%（包括上限記錄）時，失敗會被視為達到用量上限。
6. 這時會丟棄該次嘗試的輸出，並以另一個已登入的帳號執行同一個指令：先是 5 小時用量最低的帳號，接著是查不到用量的帳號，依編號順序。已達 100% 的帳號會被略過，使用 API 金鑰登入的帳號絕不會被嘗試，每個帳號最多嘗試一次。
7. 參數中有 `-c`、`--continue`、`-r` 或 `--resume` 時，只會嘗試同一個對話記錄群組中的帳號，這樣才找得到那段對話。
8. 和 `--auto` 一樣，只會嘗試同一個組織的帳號（`cross_org = true` 會解除這個限制）。

下一個帳號依和 `--auto` 相同的規則選擇（`auto_max_usage`、`auto_strategy`、`--model` 的區間）。沒有帳號可用時，ccam 會印出最後一次嘗試的輸出，說明沒有任何已登入的帳號還有剩餘額度，並以代碼 75 結束。在帳號名稱前加上 `--wait`（`ccam --fallback --wait=2h work -p …`）時，它會改為等待第一個帳號重設，然後可能再次嘗試每個帳號；它一定會先等待至少一次，所以剛達到上限的帳號不會立刻再次執行。這段時間限制的是整次執行的總等待時間，而不是每次等待的時間。用空格分開寫（`--wait 30m`）時，會以和 `ccam next` 相同的提示被拒絕。Claude Code 參數前面的 `--` 和任何啟動一樣有效（`ccam --fallback work -- -p …`）。`ccam --dry-run --fallback work -p …` 會印出第一次嘗試的指令。

除非你開啟 [`auto_switch`](#auto-switch)，否則互動式工作階段不會自動切換。執行 `ccam statusline install` 後，當 5 小時或每週用量達到 90% 時，狀態列會加上 `→ ccam next --auto --continue`（[用量限制](/features/usage)）。

## 安全性和特殊情況

- 沒有登入資訊的帳號絕不會被選中。沒有其他已登入的帳號時，單純的 `ccam next` 會以 `no logged-in account other than <name>` 停止。
- 加上 `--continue` 但資料夾中沒有對話時：`no account has a conversation in <folder> to continue`。
- 組織規則讓所有帳號都被排除時：`` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``。
- 加上 `--auto` 但沒有已登入的帳號時：`no logged-in account to pick`。當每個帳號都已用滿或用量未知時：`` no account with known usage and headroom (see `ccam usage`) ``。兩者都以結束代碼 75 結束。
- `--max-usage`、`--strategy`、`--wait` 和 `--cross-org` 只能搭配 `--auto` 使用，超出 1-100 的值或未知的策略會被拒絕。
- 分支只會複製對話，不會複製帳號的其他狀態：除非共用對話記錄，否則原帳號 `file-history/` 中的檢查點不會一起帶過去。
- 即使加上 `--dry-run`，`--auto` 仍會發出網路請求。它只會讀取，絕不寫入權杖檔案。
- 所選的帳號會透過和 `ccam <name>` 相同的流程啟動，所以適用相同的檢查，包括另一個共用相同對話記錄的 Claude Code 工作階段已在這個資料夾中執行時的警告。
- `last-used` 只在 Claude Code 實際啟動時才寫入。`--dry-run` 不會碰它，所以試執行不會讓輪替往前推進。
- `--fallback` 會暫存標準輸出直到一次嘗試結束，所以 `--output-format stream-json` 這類串流輸出會在最後一次全部送出。以串流傳送 `stream-json` 輸入、並等待回覆後才傳送更多內容的用戶端，無法使用 `--fallback`。第一次嘗試在達到上限前變更的檔案會維持變更後的狀態；下一次嘗試會從頭重新執行提示詞。

## 相關內容

- [用量限制](/features/usage)
- [用量限制與切換帳號](/docs/usage-limits)
- [共用 MCP 伺服器、設定和對話記錄](/docs/sharing)
- [Claude Code 達到用量限制時該怎麼辦](/guides/usage-limit-workflow)
