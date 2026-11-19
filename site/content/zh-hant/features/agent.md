---
title: 其他帳號上的子代理
description: ccam agent 以你的另一個帳號執行 Claude Code 任務的一部分並傳回 Claude 的回覆，讓一個長任務能用上多個帳號的用量。
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Claude Code 本身的子代理會以啟動它們的工作階段所用的帳號執行，因此花的是該帳號的用量。`ccam agent "<task>"` 會在同一個資料夾中用 `claude -p`，以你的另一個帳號執行任務，並印出 Claude 的最後一則回覆。從 Claude Code 工作階段呼叫時，它就像一個花用另一個帳號用量的子代理，所以即使工作階段本身的帳號早該用完，一個長任務也能繼續下去。

**Team 和 Enterprise 席位。** `ccam agent` 會自行選擇帳號。在同一個 Claude Team 或 Enterprise 組織內，它會在另一個席位上執行任務，這可能會違反 Anthropic 的條款。在依賴這項功能之前，請先閱讀[常見問題](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，並向你所屬的組織和 Anthropic 確認。

## 用法

Claude Code 會透過它的 Bash 工具執行它。安裝[ccam 外掛](/docs/install#claude-code-plugin)後，用一般的話提出要求即可，例如「用我的其他帳號來處理這個任務的各個部分」，或把這句話寫進你的 CLAUDE.md。外掛的技能會告訴 Claude Code 在背景啟動每個代理、寫出能獨立理解的任務，並給修改檔案的代理各自的部分或各自的 git worktree。

```sh
ccam agent "find why TestParse fails and fix it; reply with what you changed"
ccam agent --model opus -- "-v breaks the parser: find out why"
ccam agent --json < task.md            # the task from standard input
ccam agent --resume 7f3c9a2e-… "now add a test for it"
ccam agent --wait=1h "…"               # nobody has room: wait for a reset
ccam --dry-run agent "…"               # the account and command it would run
```

```console
$ ccam agent "summarize what internal/usage does"
ccam: agent runs as work (12% of its 5 hours used); session 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f
internal/usage looks up the 5-hour and weekly usage of every account …
ccam: agent finished as work; follow up with: ccam agent --resume 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f "…"
```

Claude 的最後一則回覆會寫到標準輸出；以 `ccam:` 開頭的行會寫到標準錯誤。`--json` 則改為印出 `{account, session_id, result, is_error, accounts}`，其中 `accounts` 列出任務執行時用過的每個帳號。

結束代碼：claude 本身的代碼（代理完成時為 `0`），`75`（`EX_TEMPFAIL`）表示除了呼叫端之外沒有帳號還有額度，`1` 表示其他錯誤，例如代理嘗試啟動另一個代理。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--model <model>` | 代理使用的模型；選擇時也會計入該模型的每週區間 | 帳號的預設模型 |
| `--resume <id>` | 以後續任務繼續先前某個代理的對話（即它第一行和最後一行中的工作階段 ID） | 關閉 |
| `--json` | 印出帳號、工作階段、回覆，以及任務執行時用過的帳號 | 關閉 |
| `--max-usage <1-100>` | 優先選擇用量低於此值的帳號 | `auto_max_usage`，否則為 100 |
| `--strategy headroom\|expiring\|balance` | [選擇規則](/features/next#auto) | `auto_strategy`，否則為 `headroom` |
| `--cross-org` | 也以其他組織的帳號執行 | `cross_org`，否則為關閉 |
| `--wait[=<duration>]` | 沒有帳號還有額度時，等待第一個帳號重設，最多等這麼久 | 關閉；單獨指定時為 `5h` |
| `--dry-run` | 印出它會選擇的帳號所對應的指令，但不執行 | 關閉 |

## 運作方式

1. **呼叫端會被排除。** 呼叫 `ccam agent` 的帳號（ccam 啟動的工作階段所用的帳號，否則為這個資料夾會啟動的帳號）絕不會執行代理，登入同一個 Claude 帳號、與它共用上限的其他資料夾也不會。除非 `--cross-org` 另有指定，代理會留在呼叫端所屬的[組織](/features/next#auto)內。
2. **它像 `ccam pick` 一樣選擇。** 用量的查詢方式與 [`ccam pick`](/features/pick) 相同；只要有足夠多還有額度的帳號，彼此相隔 2 分鐘內啟動的代理就會拿到不同的帳號。
3. **它以工作階段的權限執行，絕不更多。** 在 `claude -p` 中沒有人可以允許工具，所以代理會沿用呼叫它的工作階段的權限模式和推理強度。從終端機執行時，適用 Claude Code 的預設模式和你設定中的權限規則。`ccam agent` 沒有任何給予更多權限的選項。
4. **達到用量上限時，對話會繼續下去。** 代理的帳號達到上限時，ccam 會像[外掛的 hook](/features/usage#the-limit-record) 那樣記錄這次上限，並以剩餘額度最多的帳號繼續同一段對話：Claude Code 會把對話複製過去（`--resume <file> --fork-session`），第一則訊息是「Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.」每個帳號只會嘗試一次。
5. **每段代理對話都會被記錄**在 `~/.local/state/ccam/agent/<session id>.json` 中（保留 30 天）。`--resume <id>` 會在它的資料夾中繼續該對話：執行它的帳號仍有額度時以該帳號繼續，否則以另一個帳號繼續。[`ccam next --continue`](/features/next) 不會納入代理對話，所以它繼續的仍是你的對話。

## 安全性和特殊情況

- 代理無法啟動自己的代理：ccam 會以 `CCAM_AGENT=1` 標記它，而 `ccam agent` 在那裡會拒絕執行。
- 停止 `ccam agent`（Ctrl+C、訊號，或 Claude Code 中斷該指令）也會停止代理。
- 代理看不到呼叫它的對話，只看得到自己的任務。
- 同時編輯相同檔案的代理會覆寫彼此的工作。請給每個代理任務中各自的部分，或各自的 git worktree。
- 呼叫端工作階段移到另一個帳號時（[`auto_switch`](/features/next#auto-switch)），它啟動的代理也會跟著停止。它們的對話仍會保留：`--resume <id>` 可以繼續它們。
- `agent` 是指令名稱，所以新帳號不能命名為 `agent`；已有的同名帳號仍可用 `ccam run agent` 啟動。

## 相關內容

- [為指令碼選擇帳號](/features/pick)
- [切換到下一個帳號](/features/next)
- [用量限制](/features/usage)
- [所有帳號的背景工作階段](/features/agents)
