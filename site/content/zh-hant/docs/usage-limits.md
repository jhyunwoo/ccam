---
title: 用量限制與切換帳號
description: 用 ccam usage 查看每個 Claude Code 帳號的 5 小時和每週用量，並用 ccam next --auto 開啟剩餘額度最多的帳號。
updated: 2026-11-05
order: 5
translation: machine
sourceHash: "06bd349073e3"
---

`ccam usage` 會顯示每個帳號在各項限制中已使用了多少：5 小時限制、每週限制，以及各模型自己的每週限制（例如 Fable）。`ccam next --auto` 會開啟剩餘額度最多的帳號，而 `--continue` 會把目前的對話帶過去。

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

每個帳號先有一行，顯示數字的新舊程度，接著每項限制各有一個長條，並附上距離重設的時間：5 小時區間、所有模型共用的每週區間，以及各模型自己的每週區間，例如 Fable。長條在 80% 時變成黃色，在 100% 時變成紅色。最後一行會列出 `ccam next --auto` 會開啟的帳號；設定 `auto_strategy = "expiring"` 或 `"balance"` 時，則改為顯示 `→ next --auto picks <name>: <reason>`。

**Team 和 Enterprise 席位。** 自動選擇帳號（`ccam next`、`--fallback`、`ccam pick`、儀表板的 `n`，以及你開啟後的 [`auto_switch`](/features/next#auto-switch)）也會在同一個 Claude Team 或 Enterprise 組織的席位之間切換，這可能會違反 Anthropic 的條款。在依賴這項功能之前，請先閱讀[常見問題](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，並向你所屬的組織和 Anthropic 確認。

## 數字從哪裡來

ccam 會使用每個帳號自己的登入資訊，同時向 Anthropic 的用量端點查詢所有帳號。每個帳號最多等待 3 秒，並在帳號一有回應時就顯示出來。

這個端點沒有公開的文件，可能會變更。請求失敗時，ccam 會顯示最後已知的數字，並標示為快取資料。在這些數字儲存之後已經重設的時間區間會顯示 0%。如果回應格式有所變更，`ccam doctor` 會告訴你。

當 Claude Code 因用量上限而停止工作階段時，[ccam 外掛](/docs/install#claude-code-plugin)會記錄是哪個區間已滿，以及它何時重設。在那之前，不論其他數字如何，該區間都會顯示 100%，因此 `ccam next --auto` 和 `--fallback` 不會太早再次選擇該帳號。

60 秒內取得的數字會重複使用。選項：

| 選項 | 效果 |
| --- | --- |
| `--watch[=seconds]` | 持續重新整理，預設每 60 秒一次 |
| `--json` | 機器可讀的輸出 |
| `-v` | 另外顯示每個帳號的數字來源 |
| `--no-cache` | 一律重新查詢 |
| `--source cache` | 絕不查詢；只顯示已儲存的數字 |

## 過期的登入

當帳號的存取權杖過期時，ccam 不會更新它，因為在某個 Claude Code 工作階段使用同一個登入時更新權杖，可能會讓該工作階段被登出。該帳號會顯示 ``access token expired; `ccam <name>` renews it when it starts``：用 `ccam <name>` 啟動該帳號一次，Claude Code 就會更新權杖。完全過期的登入，或已被 Claude Code 清除的登入，會顯示 `` login expired → `ccam login <name>` ``，需要重新登入。

`ccam usage --refresh` 會自行更新過期的權杖，但只針對沒有 Claude Code 工作階段在執行的帳號，不論該工作階段是如何啟動的。它絕不會更新存放在 macOS 鑰匙圈或匯入的 `~/.claude` 中的登入資訊，而且一次更新可能需要最多 30 秒。

## 狀態列中的用量

```sh
ccam statusline install
```

之後 Claude Code 會在狀態列中顯示帳號名稱及其用量，例如 `[work] · Opus · 5h 24% · 7d 41%`。達到 90% 以上時，這一行也會建議執行 `ccam next --auto --continue`。如果你已經有狀態列指令，ccam 會保留它，並在前面加上帳號名稱。狀態列也會儲存它的數字，當端點沒有回應時，`ccam usage` 會使用這些數字。

## 切換

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the most room
ccam next --continue  # continue this folder's conversation on the next account
```

`ccam next` 會略過未登入的帳號，`--auto` 則會略過已達 100% 的帳號（包括 Claude Code 回報已滿、尚未重設的區間）和屬於其他組織的帳號。兩個帳號共用對話記錄時，`--continue` 會開啟同一段對話；否則 Claude Code 會把對話複製到新的帳號。請參閱[共用](/docs/sharing)和[切換到下一個帳號](/features/next#continue)。

在 `config.toml` 的最上方設定 `auto_switch = true`（Linux 和 macOS；預設關閉）時，當 Claude Code 因上限而停止互動式工作階段，而另一個同類帳號（同一個組織中的席位，或另一個個人方案）還有餘裕，ccam 會自行完成這次切換。它何時動作、絕不會做什麼，請參閱 [`auto_switch`](/features/next#auto-switch)。

對於執行 `claude -p` 的指令碼和排程工作，請在帳號名稱前加上 `--fallback`。該帳號達到上限時，ccam 會以剩餘額度最多的帳號再次執行同一個指令，並且只輸出完成的那次執行的結果：

```sh
ccam --fallback work -p "review the diff" --output-format json
```

ccam 如何決定，請參閱[達到上限時重新執行](/features/next#run-again-at-the-limit)。

[用量限制工作流程指南](/guides/usage-limit-workflow)會帶你走過一個完整的範例。
