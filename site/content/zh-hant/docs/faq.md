---
title: 常見問題
metaTitle: "ccam 常見問題：Claude Code 多帳號"
description: 關於 ccam、Claude Code 多帳號、用量限制、共用、隱私和支援的作業系統的常見問題解答。
updated: 2026-11-17
order: 6
translation: machine
sourceHash: "9b13b37391f2"
---

## ccam 是 Anthropic 的官方工具嗎？

不是。ccam 是非官方的開放原始碼專案，與 Anthropic 沒有關聯，也未獲得 Anthropic 認可。它使用的是 [Claude Code 文件](https://code.claude.com/docs/en/env-vars)中說明的 `CLAUDE_CONFIG_DIR` 設定，用來同時執行多個帳號。

## 我可以使用多個 Claude Code 帳號嗎？

這取決於各帳號的條款；如果是工作帳號，也取決於你所屬組織的政策：Pro 和 Max 適用 Anthropic 的[消費者條款](https://www.anthropic.com/legal/consumer-terms)（Consumer Terms），Team 和 Enterprise 適用它的[商業條款](https://www.anthropic.com/legal/commercial-terms)（Commercial Terms），所有方案都適用它的[使用政策](https://www.anthropic.com/legal/aup)（Usage Policy）。使用多個帳號前請先確認這些規定。ccam 不會提高、合併或隱藏帳號的限制，每個帳號仍會在自己的限制處停下。它會在你要求時，把你的工作移到你的另一個帳號；或在你開啟 [`auto_switch`](/features/next#auto-switch)（Linux 和 macOS）後，在達到上限時自行這麼做。使用前請先確認你的帳號是否允許這麼做；決定權在你。

**Team 和 Enterprise 席位。** 自動選擇帳號（`ccam next`、`--fallback`、`ccam pick`、儀表板的 `n`，以及你開啟後的 [`auto_switch`](/features/next#auto-switch)）也會在同一個 Claude Team 或 Enterprise 組織的席位之間切換。Anthropic 說明中心關於 [Team 方案](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)的文章寫道：「Each seat on a Team plan is assigned to one named individual, and only that person may use it」（Team 方案的每個席位都指派給一位具名的個人，而且只有這個人可以使用），以及「Users can't route their work through another user's account to use that user's remaining usage」（使用者不能經由其他使用者的帳號處理自己的工作，藉此使用對方剩餘的用量）；關於 [Enterprise 方案](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)的文章也是如此。因此，在席位之間使用自動選擇可能會違反 Anthropic 的條款。要不要在這種情況下使用它由你決定：在依賴這項功能之前，請先向你所屬的組織和 Anthropic 確認。個人方案（Pro、Max）的情況請參閱上一段。

## ccam 能規避用量限制嗎？

無法超出任何帳號自己的限制：每個帳號都保有各自的限制並在達到時停下，ccam 也不會提高、合併或隱藏這些限制。ccam 會顯示這些限制，並在某個帳號用完額度時，讓你在自己的另一個帳號繼續；開啟 `auto_switch` 時，它會自行這麼做（請參閱下一個問題）。在同一個 Team 或 Enterprise 組織的席位之間，在另一個席位繼續會用掉那個席位剩餘的用量，而 Anthropic 說明中心不允許這樣使用其他使用者的席位（請參閱上方的 Team 和 Enterprise 席位）。

## ccam 能在工作階段達到上限時自行切換帳號嗎？

只有在你開啟這項功能時才會。在 `~/.config/ccam/config.toml` 的最上方設定 `auto_switch = true`（Linux 和 macOS；預設關閉，屬於實驗性功能）後，當 Claude Code 因 5 小時、每週、Opus、Sonnet 或 Fable 上限而停止一個互動式工作階段（由 ccam 啟動的工作階段）時，ccam 會結束該工作階段，並以你另一個還有餘裕的帳號接續對話：同一個 Team 或 Enterprise 組織的另一個席位，或另一個個人方案（Pro、Max）。它絕不會在達到上限之前切換，絕不會在背景監看你的用量，也絕不會把工作移出 Team 或 Enterprise 組織或移入其中。

在席位之間，它會把你的工作從一個席位移到另一個席位，這可能會違反 Anthropic 的條款（請參閱上方的 Team 和 Enterprise 席位）。至於個人方案，Anthropic 的條款並沒有明確規範：消費者條款禁止「through automated or non-human means」（透過自動化或非人為的方式）存取服務，除非 Anthropic 明確允許，而一次移動會代替你送出一則訊息。要不要開啟由你決定。

[`auto_switch`](/features/next#auto-switch) 列出了所有條件，以及一次移動的代價。

## 如何在一台電腦上使用兩個 Claude Code 帳號？

用 `CLAUDE_CONFIG_DIR` 為每個帳號指定獨立的設定資料夾。ccam 會替你完成這件事：為每個帳號執行 `ccam add`，再用 `ccam personal` 和 `ccam work`，或你自己命名的快捷指令（例如 `ccwork`）啟動。[指南](/guides/multiple-claude-code-accounts)也說明了如何手動設定。

## 如何在另一個帳號繼續對話？

在該資料夾中執行 `ccam next --continue`：它會把這個資料夾最近的一段對話帶到下一個帳號（加上 `--auto` 則選用量餘裕最多的帳號）。在 Claude Code 中，[ccam 外掛](/docs/install#claude-code-plugin)的 `/ccam:switch` 會對你目前所在的工作階段做同樣的事。當兩個帳號共用對話記錄時（`ccam group set default --add history`，或為它們另建一個群組：`ccam group add pair --items history --members work,side`），新帳號會開啟同一段對話；否則 Claude Code 會把對話複製到新帳號的記錄中。加上 `--brief`，新帳號就會從一份簡短的摘要開始，而不是整段對話。

## ccam 會看到或傳送我的登入資訊嗎？

ccam 絕不會在帳號之間共用或複製你的 Claude 登入資訊，也絕不會複製 MCP 伺服器的登入。它讀取帳號的存取權杖只有一個用途：向 Anthropic 的用量端點查詢該帳號的用量，用於 `ccam usage`、儀表板、`ccam doctor` 和自動選擇（`ccam next --auto`、`ccam pick`、`--fallback`，以及工作階段達到上限時的 `auto_switch`）。只有在執行 `ccam usage --refresh` 時才會更新過期的權杖；設定 `usage_source = "claude"` 時，它不會把權杖傳送到任何地方，只會讀取已儲存的登入資訊（檔案，或 macOS 鑰匙圈中的項目）來判斷它是否仍然有效。`ccam ext` 和 `ccam doctor` 會讀取帳號中 MCP 登入的名稱，但絕不讀取其權杖。ccam 沒有遙測功能。除了這些查詢，它唯一會自行發出的請求是儀表板向 GitHub 檢查是否有較新的 ccam 發行版本，其中不含任何關於你或你帳號的資訊。

## 如何在遠端伺服器上登入？

像平常一樣在 SSH 工作階段中執行 `ccam add` 或 `ccam login <name>`。Claude Code 會在伺服器上等待登入，所以你電腦上的瀏覽器最後會停在無法載入的 `http://localhost:<port>/callback?code=…` 頁面。從瀏覽器複製這個網址並貼到終端機中：ccam 會把它交給伺服器上的 Claude Code，由它完成登入。這個網址包含一次性代碼，請不要分享給別人。

## 我的帳號存放在哪裡？

存放在 `~/.local/share/ccam/profiles/<name>/`，每個帳號一個資料夾。在 macOS 上，Claude Code 會把登入資訊本身存放在鑰匙圈中，名稱由該資料夾衍生而來，因此 `ccam rename` 會在那裡保留資料夾原本的名稱。

## 現有的 ~/.claude 會怎麼樣？

不會有任何變化，除非你執行 `ccam import` 或 `ccam migrate`。匯入會在原位置登錄該資料夾，因此直接執行 `claude` 和 `ccam <name>` 都會繼續使用相同的帳號和設定。`ccam migrate` 會把其中的設定、`CLAUDE.md`、擴充功能和 MCP 伺服器複製到共用的 `default` 群組，而且只會讀取該資料夾。

## ccam 支援哪些系統？

Linux（Ubuntu、Debian、Fedora、Arch、Alpine 等）、macOS，以及 Windows 10 和 11（原生執行或在 WSL 中執行）。在 Windows 上，共用 `settings.json` 這類檔案需要開發人員模式；共用資料夾則不需要，而在沒有開發人員模式時，第一次 `ccam add` 建立的群組 `default` 只包含 MCP 伺服器和擴充功能（[Windows 上的共用](/features/share#on-windows)）。

## 為什麼某個 MCP 伺服器只出現在一個帳號中？

在工作階段中用 `claude mcp add` 新增的伺服器只會儲存在該帳號，而且沒有加上 `-s user` 時只適用於目前的資料夾。下次啟動該帳號時，ccam 會提醒你，`ccam mcp list` 和 `ccam ext` 也會顯示這些伺服器；用 `ccam mcp import <account> <server>` 將伺服器移到共用群組。

## ccam 會拖慢 Claude Code 嗎？

透過 ccam 啟動 Claude Code 大約會增加 4 毫秒，這是在一台筆記型電腦上用 `scripts/perf-check.sh` 測得的。ccam 會讀取一個小檔案、檢查共用連結，然後以 `claude` 取代自身，所以預設情況下，你工作時它並不在執行。Windows 無法取代行程，所以在那裡 ccam 會啟動 Claude Code 並等待它結束，除此之外什麼都不做。開啟 `relay` 或 `auto_switch` 時，ccam 會以同樣的方式在每個互動式工作階段旁持續執行，以便在另一個帳號接續對話。

## 如何更新 ccam？

執行 `ccam update`。它會下載最新的發行版本，檢查總和檢查碼，然後取代 `ccam` 檔案；`ccam update --check` 只會告訴你是否有更新的發行版本。儀表板也會在有更新的發行版本時告訴你，按下 `U` 會先詢問再安裝；ccam 不會自行更新。如果你安裝的是 `.deb` 套件，請改用 apt 安裝新的 `.deb`。請參閱[更新](/docs/install#update)。

## 如何解除安裝 ccam？

執行 `ccam uninstall`。它會先列出將要移除的內容並詢問你，然後移除每個帳號（先登出）、快捷指令、`~/.config/ccam`、`~/.local/share/ccam`、`~/.local/state/ccam` 和 `ccam` 指令。匯入的 `~/.claude` 會保留，並取回它在各群組中所共用內容的副本。請參閱[安裝](/docs/install#uninstall)。
