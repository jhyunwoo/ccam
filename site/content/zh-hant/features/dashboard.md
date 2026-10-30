---
title: 儀表板
description: 在終端機中執行 ccam 會開啟全螢幕儀表板，顯示帳號、群組、資料夾規則、擴充功能和檢查結果；每項變更都會先顯示計畫。
updated: 2026-10-30
order: 15
translation: machine
sourceHash: "d65f1b5bb0c8"
---

不帶參數的 `ccam` 在標準輸入和輸出都是終端機時會開啟全螢幕儀表板，否則會印出 `ccam list`。儀表板會顯示每個帳號及其用量和群組，並能完成各個指令的工作：啟動、新增、登入、變更群組、設定資料夾規則、列出背景工作階段、檢查和修復。啟動帳號或連接到背景工作階段時，儀表板會關閉；新增帳號或登入時則會暫時把終端機交出去，完成後再回到儀表板。

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

分頁名稱下方的一行用一句話說明目前分頁的用途。在較寬的終端機上，Accounts 分頁還會顯示每個帳號的電子郵件。在每個分頁上，分頁所在那一行的結尾都會顯示 [`auto_switch`](/features/next#auto-switch) 是否開啟（Linux 和 macOS；在 Windows 上會被忽略）。開啟時，Health 分頁會列出這項設定還需要什麼：Claude Code 2.1.214 或更新版本，以及在每個以 Claude 訂閱登入的帳號中開啟 ccam 外掛。

儀表板每次開啟時都會在背景向 GitHub 詢問是否有更新的 ccam 發行版本；第一個畫面從不等待它，你離線時則不顯示任何內容。有新版本時，標題下方的一行會推薦它：

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` 會先詢問，再以 [`ccam update`](/features/update) 的方式更新：以 SHA-256 校驗碼檢查該發行版本，並一次取代 ccam。你目前所在的儀表板仍是舊版本，所以請退出並重新啟動 `ccam`。若無法在這裡更新，例如透過 `.deb` 套件安裝的 ccam 副本，儀表板會顯示與 `ccam update` 相同的建議，並保留這個更新提示。

第一個畫面只根據 `config.toml` 繪製，遠低於 100 ms；用量會在每個帳號回應時陸續填入。每項限制都有自己的長條：5 小時區間、所有模型共用的每週區間，以及 Fable 自己的每週區間；`-` 表示該帳號沒有這個區間。百分比在 80% 時變成黃色，在 100% 時變成紅色。如果查詢失敗，快取的數值會以淡色顯示。Claude Code 回報已滿的區間在重設前都會顯示 100%。存取權杖過期時會以淡色顯示 `⚠ token expired`（Claude Code 會在下次啟動時更新它），而 `⚠ login expired` 表示該帳號需要執行 `ccam login`。

## 分頁

| 分頁 | 顯示的內容 | 按鍵 |
| --- | --- | --- |
| Accounts | 編號、名稱、電子郵件、5 小時、每週和 Fable 用量長條、上次使用時間、群組；`✔` 為預設帳號，`▸` 為這個資料夾的帳號，`•` 為執行中（也包括不透過 ccam 啟動的 Claude Code：直接執行 `claude`、IDE、背景工作階段及其常駐程式） | 見下文 |
| Agents | [`ccam agents`](/features/agents) 列出的內容：所有帳號的背景工作階段和開啟中的工作階段，以及它們的狀態、資料夾和啟動時間（見下文） | `enter` 連接（如同 `ccam attach`），`x` 停止，`r` 重新載入 |
| Groups | 每個群組、它共用的項目和它的帳號，以及不屬於任何群組的帳號 | `enter` 編輯，`n` 新增，`e` 重新命名，`x` 移除 |
| Folders | [`ccam folder`](/features/folder) 規則，以及適用於這裡的規則，即使你是透過符號連結來到這個資料夾也一樣（新規則會以解析連結後的路徑儲存） | `n` 在這裡使用某個帳號，`enter` 變更，`x` 移除 |
| Extensions | [`ccam ext`](/features/ext) 顯示的內容，每個擴充功能只列一次，並附上它來自的群組；`off:` 列出外掛未開啟的帳號，帳號自己的伺服器中只適用於這個資料夾的那些，`This folder's project`，以及不屬於任何帳號的問題 | `r` 重新載入 |
| Health | [`ccam doctor`](/features/doctor) 回報的內容，每個帳號一個方框（見下文） | `f` 修復選取的方框，`F` 修復全部（`doctor --fix`），`↑`/`↓` 在方框之間移動 |

`←`/`→` 或 `Tab` 切換分頁，`?` 列出目前分頁的按鍵，`r` 會略過短期快取重新整理用量，並重新載入 Agents、Extensions 和 Health 分頁，`q` 或 `Esc` 則結束。

## Accounts 分頁的按鍵

| 按鍵 | 作用 |
| --- | --- |
| `Enter` | 以選取的帳號啟動 Claude Code |
| `1`–`9` | 移到該編號的帳號；快速輸入 `1` `2` 可移到 12。只有 `Enter` 會啟動 |
| `↑`/`↓`、`k`/`j` | 移動 |
| `n` | 啟動 [`ccam next --auto`](/features/next) 會選擇的帳號 |
| `o` | 帶選項啟動：`--continue`、`--resume`、其他參數，或試執行 |
| `u` | 用量詳細資料：5 小時、每週和各模型的區間，以及它們的重設時間和數字的來源 |
| `a` | 新增帳號（使用下文的設定精靈） |
| `l` | 重新登入 |
| `s` | 選擇帳號的群組 |
| `d` | 設為預設帳號 |
| `f` | 在目前的資料夾中使用它 |
| `e` | 重新命名 |
| `c` | 設定或移除它的快捷指令 |
| `x` | 移除它（會先詢問，預設選取 No） |

**Team 和 Enterprise 席位。** 自動選擇帳號（`ccam next`、`--fallback`、`ccam pick`、儀表板的 `n`，以及你開啟後的 [`auto_switch`](/features/next#auto-switch)）也會在同一個 Claude Team 或 Enterprise 組織的席位之間切換，這可能會違反 Anthropic 的條款。在依賴這項功能之前，請先閱讀[常見問題](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，並向你所屬的組織和 Anthropic 確認。

## Health 分頁

開啟 Health 分頁時，它會執行 [`ccam doctor`](/features/doctor) 的檢查並繪製相同的畫面：第一行為每個帳號顯示一個標記，下方每個帳號各有一個方框。方框內每一列是一個檢查過的項目，`Folder`、`Login` 和 `Share links` 排在最前面，標記為 `✔` 正常、`ℹ` 資訊、`⚠` 警告或 `✘` 錯誤，讓你一眼看出哪個帳號在哪裡有問題。不屬於任何帳號的檢查結果放在 `General` 方框中。

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

`↑`/`↓` 在方框之間移動，選取的方框會以強調色邊框顯示。`f` 只修復選取的方框，方式與 `ccam doctor --fix` 相同：在帳號的方框上只修復該帳號，在 `General` 方框上修復不屬於任何帳號的檢查結果（外掛路徑和快捷指令）；其他帳號只會重新檢查。它會先詢問（``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``），沒有可修復項目的方框只會顯示 `nothing to repair in <name>`。`F` 會像 `ccam doctor --fix` 一樣一次修復全部；`r` 會重新檢查。

## Agents 分頁

Agents 分頁會顯示 [`ccam agents`](/features/agents) 印出的內容：所有帳號的背景工作階段和開啟中的工作階段，以及每個工作階段所屬的帳號。第一次開啟這個分頁時、按下 `r` 時、停止工作階段後，以及設定變更後，它會同時詢問每個帳號的 Claude Code（`claude agents --json`，每個的逾時為 15 秒），在所有帳號回應之前會顯示 `Asking every account's Claude Code for its sessions…`。

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

`State` 是背景工作階段的狀態以及它的行程正在做什麼；互動式工作階段則顯示 `open`，和 `ccam agents` 相同。如果無法查詢某個帳號，清單下方會有一行 `✘` 顯示 Claude Code 的回應，例如上面的 `lab`，它的 agent view 已關閉。已完成的背景工作階段不會列出（`ccam agents --all` 會列出它們）；當所有帳號都已回應而且沒有任何工作階段時，分頁會顯示 ``No sessions. Start one in the background with `ccam <account> --bg "…"`.``。終端機較窄時，分頁會依序省略 `Started`、`Name`、`Folder`，並從左側截斷過長的資料夾路徑。

在背景工作階段上按 `Enter` 會關閉儀表板，並在這個終端機中以它的帳號開啟該工作階段，方式和 `ccam attach <id>` 相同；儀表板不會再回來（在 Claude Code 中，`←` 會前往它的 agent view，`Ctrl+Z` 會回到 Shell）。開啟中的工作階段沒有 ID，所以 `Enter` 和 `x` 只會提示需要選擇背景工作階段。`x` 會先詢問，預設選取 `No, keep it`：``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.``。選擇是時，ccam 會以該帳號執行 `claude stop <id>`，顯示 ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` 或無法停止的原因，然後重新讀取清單。

## 變更群組

在 Groups 分頁中按 `Enter` 會開啟編輯器：左邊是項目，右邊是帳號，用 `Tab` 切換，用 `Space` 勾選或取消。其中某個帳號已經從另一個群組取得的項目會顯示為灰色並附上原因，會重複取得某個項目的帳號也一樣；一個項目只能來自一個群組。按 `Enter` 會顯示計畫：

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

如果某個檔案和群組的副本內容不同，ccam 會先詢問你要怎麼做（保留群組的版本、使用我的版本，或合併）。接著變更會在背景以和 [`ccam group`](/features/share) 相同的程式碼執行：先儲存設定，再依序移動每個帳號的資料，原始檔案則放到 `.ccam-backup/`；退出 MCP 共用的帳號會在它的 `.claude.json` 中保留群組的伺服器。

會移動執行中帳號資料的變更會被拒絕，並列出該帳號的名稱，即使該帳號是在計畫開啟期間才啟動的也一樣；因另一個 ccam 行程變更了群組而過時的計畫，會以 `the groups changed while this plan waited; nothing was moved, make the change again` 拒絕。當另一個 ccam 行程正在變更群組時，新的變更、重新命名或移除會顯示 `another ccam process is changing the groups right now`；若上一次變更中途停止，它們會等到你在 Health 分頁按下 `F`（或執行 `ccam doctor --fix`）完成它之後。

## 新增帳號

`a` 會開啟 `ccam add` 在終端機中也使用的設定精靈：

1. 帳號名稱。
2. 登入方式：現在以 Claude 訂閱登入（在瀏覽器中）、[API 金鑰](/features/api)，或稍後再以 Claude 訂閱登入。
3. 使用 API 金鑰時：ccam 每次啟動時從哪裡讀取金鑰（`env:NAME` 或 `command:COMMAND`）、端點、金鑰的傳送方式和模型。
4. 要加入的群組，其中 `default` 已勾選（還沒有任何群組時，則是 `default` 群組將會共用的項目）。
5. 快捷指令，如果你需要的話。
6. 讓你確認的摘要。

`Esc` 會回到上一步，在第一步時則會關閉設定精靈。

## 移除帳號

`x` 會在任何動作之前先詢問，預設選取 No，並說明移除會做什麼：`It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` API 金鑰帳號沒有可撤銷的登入；另一個帳號持有同一組登入時，也不會登出。匯入的帳號會保留它的資料夾：`It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## 色彩

色彩依循 Claude Code 終端機介面的角色配色，所以 ccam 放在它旁邊看起來很協調；群組名稱會根據名稱取得一種顏色。ccam 會向終端機查詢背景色，並選擇淺色或深色主題。你也可以在 `config.toml` 的最上方、第一個 `[[group]]`、`[[profile]]` 或 `[[folder]]` 之前自行設定：

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

`-daltonized` 主題適用於紅綠色盲，`-ansi` 主題則使用終端機本身的 16 種顏色。`NO_COLOR` 會關閉色彩。

## 安全性和特殊情況

- 儀表板呼叫的函式和指令相同，所以每條規則都成立：移動資料之前先顯示計畫、拒絕變更執行中的帳號、保留備份，而且登入權杖絕不會離開帳號。
- 儀表板沒有試執行：不帶名稱的 `ccam --dry-run` 會印出這個資料夾的帳號會如何啟動，在某個帳號上按 `o` 則可以只對這一次啟動進行試執行。
- 在 Claude Code 啟動或登入開啟之前，終端機會先還原。
- 畫面太小放不下對話框時，會先省略分頁的列；過長的行會以 `…` 截斷，而不是換行。畫面較矮時，接著會省略空白行和分頁說明，讓按鍵提示那一行保持可見。小於 40 欄或 8 列時，儀表板只會顯示 `Too small: ccam needs 40x8 (now 36x7)`。
- 在套用變更、移除、更新或修復期間按 `Ctrl+C`，會等它完成：狀態列會顯示 `ccam quits when it is done.`
- `U` 更新 ccam 之後，儀表板從新增帳號或登入返回時，更新提示不會再出現。

## 相關內容

- [用群組共用](/features/share)
- [啟動 Claude Code](/features/launch)
- [管理帳號](/features/accounts)
- [用量限制](/features/usage)
