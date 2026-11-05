---
title: 為指令碼選擇帳號
description: ccam pick 為指令碼印出剩餘額度最多的帳號，讓同時啟動的工作程序使用不同帳號，沒有帳號有額度時以 75 結束。
updated: 2026-11-05
order: 12
translation: machine
sourceHash: "0ec25752ace9"
---

`ccam pick` 會印出剩餘額度最多的已登入帳號名稱，讓指令碼用它啟動 Claude Code。只要有足夠多還有額度的帳號，同時啟動的工作程序就會各自拿到不同的帳號，所以一批 `claude -p` 工作會分散到你的各個帳號，而不是把一個帳號用光。

**Team 和 Enterprise 席位。** 自動選擇帳號（`ccam next`、`--fallback`、`ccam pick`、儀表板的 `n`，以及你開啟後的 [`auto_switch`](/features/next#auto-switch)）也會在同一個 Claude Team 或 Enterprise 組織的席位之間切換，這可能會違反 Anthropic 的條款。在依賴這項功能之前，請先閱讀[常見問題](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，並向你所屬的組織和 Anthropic 確認。

## 用法

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # most weekly quota left per day; fewer sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

結束代碼：`0` 表示在標準輸出上印出了名稱，`75`（`EX_TEMPFAIL`）表示沒有帳號還有額度或沒有帳號登入，`1` 表示其他錯誤，例如錯誤的旗標。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--json` | 印出帳號、它的 5 小時和每週用量、這些數字的來源，以及選擇它的原因 | 關閉 |
| `--model <model>` | 也計入該模型的每週區間，例如 `opus` | 無 |
| `--max-usage <1-100>` | 優先選擇用量低於此值的帳號 | `auto_max_usage`，否則為 100 |
| `--strategy headroom\|expiring\|balance` | 在 5 小時用量與最低值相差不超過 3 個百分點的帳號中每週用量最低的帳號、5 小時區間最早重設的帳號，或平均每天剩餘每週額度最多的帳號 | `auto_strategy`，否則為 `headroom` |
| `--cross-org` | 也選擇其他組織的帳號 | `cross_org`，否則為關閉 |
| `--wait[=<duration>]` | 沒有帳號還有額度時，等待第一個帳號重設，最多等這麼久。值要寫在 `=` 之後：`--wait 30m` 會被拒絕，並提示 `--wait takes its value after '=': --wait=30m` | 關閉；單獨指定時為 `5h` |
| `--dry-run` | 只選擇，不記錄這次選擇 | 關閉 |

## 運作方式

1. ccam 列出這個資料夾會啟動的帳號所屬組織（請參閱[組織](/features/next#auto)）中已登入的帳號（API 金鑰帳號和 Console 登入沒有限制，因此不列入）；和該帳號登入同一個 Claude 帳號的資料夾也不列入。
2. 它取得一個鎖定（`~/.local/state/ccam/picks.lock`），所以同時執行的選擇會依序進行。
3. 它以和 `ccam usage` 相同的方式一次查詢所有帳號的用量，並使用 60 秒的快取。Claude Code 回報已滿的區間在重設前都算作 100%（[上限記錄](/features/usage#the-limit-record)），所以剛達到上限的帳號即使已儲存的數字顯示還有餘裕，也不會被列入。
4. 先排除過去 2 分鐘內被選過的帳號，不論是由 `ccam pick` 選中，還是 [`auto_switch`](/features/next#auto-switch) 把工作階段移到那裡（`~/.local/state/ccam/picks/<name>`；時間戳記在未來的選取紀錄不算數），再由 `ccam next --auto` 的[選擇規則](/features/next#auto)在其餘帳號中選擇：優先選低於 `auto_max_usage` 的帳號，依策略選擇，絕不選擇用量 100% 的帳號。如果因此略過了規則原本會最先選的帳號，`--json` 的 `reason` 會說明這一點：`30% of its 5 hours used; work was picked a moment ago`。
5. 如果每個還有額度的帳號最近都被選過，規則會再次在所有帳號中選擇，所以有三個帳號時，第四個工作程序會和其他工作程序共用其中一個帳號。
6. 記錄這次選擇並釋放鎖定，然後印出名稱。

2 分鐘的時間窗很重要，因為新工作的用量要過一段時間才會反映在數字上：沒有它的話，同時啟動的五個工作程序都會把同一個帳號看成剩餘額度最多的帳號。

## 安全性和特殊情況

- `ccam pick` 只讀取用量並寫入選擇的時間；除了 `ccam usage` 會做的事之外，它絕不更新或讀取權杖。
- 選擇並不會保留帳號。2 分鐘後，兩個指令碼仍可能在同一個帳號上啟動；請使用 `--fallback`，在該帳號達到上限時恢復。
- 以 `--wait` 等待時不持有鎖定，所以其他選擇可以同時進行。
- `ccam pick` 會忽略 [`ccam folder`](/features/folder)：它只依用量選擇。

## 相關內容

- [切換到下一個帳號](/features/next)
- [達到上限時重新執行](/features/next#run-again-at-the-limit)
- [用量限制](/features/usage)
- [API 金鑰帳號](/features/api)
