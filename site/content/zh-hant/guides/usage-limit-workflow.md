---
title: Claude Code 達到用量限制時該怎麼辦
description: 等待額度重置，或立即在另一個 Claude Code 帳號接續同一個對話。逐步說明 5 小時限制和每週限制如何運作。
updated: 2026-11-05
order: 3
translation: machine
sourceHash: "53b1e3634ea9"
---

當一個帳號達到用量限制（常被稱為 rate limit）時，你有兩個選擇：等到限制重置，或現在就在你擁有的另一個帳號接續同一個對話。使用 ccam 時，第二種做法只需要一個指令：`ccam next --auto --continue`。

## Claude Code 有哪些限制？

在 Pro 和 Max 方案中，Claude Code 有[5 小時工作階段限制和每週限制](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan)。這兩個限制與網頁版和應用程式中的 Claude 共用，用量會同時計入兩者。其中一個用完時，Claude Code 會[顯示一則訊息](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit)，例如「You've hit your session limit · resets 3:45pm」。切換模型沒有幫助，因為這些限制涵蓋所有模型。

## 做法 1：等待重置

如果你只有一個帳號，就只能等待。較新版本的 Claude Code 會[在開啟的工作階段中等待](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset)，限制重置後自動繼續任務。在 Claude Code 中執行 `/usage` 可以查看你已用掉多少方案額度。

## 做法 2：現在就在另一個帳號接續

如果你有第二個帳號，就能立刻接續同一個對話。只需設定一次：

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

現在，工作時狀態列會顯示目前的帳號和它的 5 小時用量，達到限制時就不會措手不及。

達到限制後，先離開 Claude Code（`/exit`），然後執行：

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

在剩餘額度最多的帳號中開啟這個對話：

```sh
ccam next --auto --continue
```

Claude Code 會以 `work` 帳號啟動，並開啟目前資料夾中最近的對話，包含之前的所有內容。共用對話記錄（上面的 `ccam group set default --add history`）代表兩個帳號處理的是同一個對話；沒有共用時，Claude Code 會把對話複製到 `work`，ccam 也會告訴你第一次回覆必須重新讀取多少 token。

在 Claude Code 中安裝 [ccam 外掛](/docs/install#claude-code-plugin)後，因上限而停止的工作階段會顯示通知並記錄已滿的區間，讓 `ccam next --auto` 在重設前不會再次選擇該帳號。`/ccam:switch`（或 `/ccam:switch work`）會為 `ccam next --continue` 記下這個對話，讓它接續的正是這個對話。

**Team 和 Enterprise 席位。** 自動選擇帳號（`ccam next`、`--fallback`、`ccam pick`、儀表板的 `n`，以及你開啟後的 [`auto_switch`](/features/next#auto-switch)）也會在同一個 Claude Team 或 Enterprise 組織的席位之間切換，這可能會違反 Anthropic 的條款。在依賴這項功能之前，請先閱讀[常見問題](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，並向你所屬的組織和 Anthropic 確認。

## ccam 會自動切換帳號嗎？

只有在你開啟這項功能時才會。

- 預設情況下，你工作時 ccam 並沒有在執行：它啟動 Claude Code 後就用 Claude Code 取代自己，所以何時切換由你決定。`ccam next --auto` 會依用量選擇：在 5 小時用量與最低值相差不超過 3 個百分點的帳號中，選每週用量最低的帳號；它會略過用量已達 100% 或未登入的帳號，並且只在同一個 Team 或 Enterprise 組織內、或在個人方案之間選擇。
- 在 `config.toml` 中設定 `relay = true` 時，ccam 會改為在工作階段旁持續執行，讓 `/ccam:switch` 能自行接續對話，並沿用工作階段的選項、權限模式、推理強度和代理；在把對話移到另一個組織的帳號之前，ccam 會先詢問你。
- 設定 `auto_switch = true` 時（Linux 和 macOS；實驗性功能），ccam 也會自行完成切換：當 Claude Code 因上限而停止工作階段，而另一個同類帳號還有餘裕時，ccam 會結束這個工作階段，並在那個帳號接續對話（[條件和代價](/features/next#auto-switch)）。

另一些工具則是替換執行中工作階段的登入資訊；[ccam 與其他帳號切換工具](/compare/account-switchers)說明了兩者的差異。

## 注意事項

- 在你工作時所在的同一個資料夾中執行 `ccam next --continue`；對話是依專案資料夾儲存的。
- 避免讓兩個帳號同時處理同一個對話。你這樣做時，ccam 會發出警告。
- 沒有共用對話記錄時，對話會被複製到另一個帳號，之後兩份副本各自發展；`/rewind` 用來還原程式碼的檢查點會留在第一個帳號。
- 新帳號上的第一則訊息可能會比平常消耗更多額度，因為該帳號還沒有快取這個對話。`ccam next --continue --brief` 則會改用一份簡短的摘要來開始。
- 是否可以這樣在帳號之間切換，取決於適用於各帳號的條款。

所有選項請參閱[用量限制](/features/usage)和[切換到下一個帳號](/features/next)。
