---
title: Giới hạn sử dụng
description: "ccam usage hiển thị giới hạn 5 giờ và hàng tuần của mọi tài khoản cùng lúc, và ccam statusline đưa tên tài khoản cùng mức sử dụng vào Claude Code."
updated: 2026-11-01
order: 7
translation: machine
sourceHash: 54b296f7da35
---

`ccam usage` cho biết, với mọi tài khoản, mỗi giới hạn đã được dùng bao nhiêu: giới hạn 5 giờ, giới hạn hàng tuần và giới hạn hàng tuần riêng của từng mô hình (như Fable). Mỗi giới hạn có một thanh và thời gian đến lần đặt lại. Hãy dùng nó trước khi bắt đầu một tác vụ dài, hoặc khi Claude Code báo rằng sắp chạm giới hạn, để xem tài khoản nào còn hạn mức. `ccam statusline` là lệnh dòng trạng thái cho Claude Code: nó hiển thị tên tài khoản và mức sử dụng bên trong Claude Code và lưu các con số cho `ccam usage`. Để xem số token và chi phí ước tính đằng sau mức sử dụng đó, `ccam ccusage` chạy [ccusage](https://ccusage.com) trên các cuộc hội thoại của mọi tài khoản ([bên dưới](#tokens-and-costs-ccam-ccusage)).

## Cách dùng

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

Mỗi tài khoản có một dòng gồm số thứ tự, tên và độ cũ của các con số; `(cached)` nghĩa là chúng lấy từ một tệp đã lưu, không phải từ endpoint. Bên dưới, mỗi giới hạn có một dòng với thanh 20 ô, tỷ lệ phần trăm và thời gian đến khi khoảng đó được đặt lại: `5-hour`, `Weekly` (mọi mô hình), rồi đến khoảng hàng tuần riêng của từng mô hình, như `Fable`. Khoảng mà tài khoản không có thì hiển thị `-`. Tỷ lệ phần trăm chuyển sang màu vàng ở 80% và màu đỏ ở 100%.

Nếu Claude Code đã từ chối một yêu cầu vì một khoảng đã đầy, khoảng đó hiển thị 100% cho đến khi được đặt lại, bất kể các con số đã lưu nói gì (xem [bản ghi giới hạn](#the-limit-record)). Trên terminal, mọi tài khoản xuất hiện ngay cùng lúc dưới dạng `fetching…`, và các dòng của từng tài khoản được vẽ khi nó phản hồi.

Dòng cuối, chỉ xuất hiện khi bạn có nhiều hơn một tài khoản, nêu tên tài khoản mà [`ccam next --auto`](/features/next) sẽ chọn. Nó tuân theo cùng các quy tắc:

- `auto_max_usage` và `auto_strategy`;
- các phiên đang chạy;
- chỉ các tài khoản thuộc tổ chức của tài khoản dùng gần nhất, trừ khi đặt `cross_org`;
- mọi khoảng ở mức 100% đều bị bỏ qua, nên một tài khoản đã đầy hạn mức tuần không bao giờ được nêu tên.

Với chiến lược mặc định `headroom`, dòng này là `→ Most headroom: <name>`; với `expiring` hoặc `balance`, nó nêu lý do, chẳng hạn `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Seat của gói Team và Enterprise.** Việc tự động chọn tài khoản (`ccam next`, `--fallback`, `ccam pick`, phím `n` trên bảng điều khiển và, khi bạn bật nó, [`auto_switch`](/features/next#auto-switch)) cũng chuyển đổi giữa các seat của cùng một tổ chức Claude Team hoặc Enterprise, điều này có thể vi phạm điều khoản của Anthropic. Hãy đọc [Câu hỏi thường gặp](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) và xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó.

Tài khoản chưa đăng nhập hiển thị `not logged in → ccam login <name>`, còn tài khoản dùng khóa API hiển thị `no limits (API)`. Một access token đã hết hạn hiển thị ``access token expired; `ccam <name>` renews it when it starts``, phía trên các con số đã lưu, nếu có: Claude Code gia hạn nó ở lần khởi động tiếp theo. Một thông tin đăng nhập đã hết hạn hoàn toàn, vì refresh token của nó đã hết hạn hoặc Claude Code đã xóa nó sau khi máy chủ từ chối, hiển thị `` login expired → `ccam login <name>` ``.

Sau `ccam statusline install`, Claude Code hiển thị một dòng như `[work] · Opus · 5h 24% · 7d 41%`.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--json` | In các báo cáo dưới dạng JSON: `name`, `slot`, `status`, `limits` (kèm `credits` khi tài khoản bật tín dụng sử dụng), `source`, `at`, `error`, và `limited` cùng `limited_windows` khi một bản ghi giới hạn đặt các khoảng thành 100%. Khi không có tài khoản nào, nó in ra `[]` | tắt |
| `--watch[=N]` | Vẽ lại mỗi `N` giây cho đến khi bạn nhấn Ctrl+C. `--watch` không kèm giá trị, hoặc giá trị từ 0 trở xuống, nghĩa là 60. Giá trị đặt sau `=`: `--watch 10` bị từ chối kèm gợi ý | tắt |
| `-v`, `--verbose` | Thêm một dòng `source:` dưới mỗi tài khoản: nơi các con số của nó được lấy (endpoint mức sử dụng, Claude Code, dòng trạng thái, một tệp đã lưu, hoặc giới hạn sử dụng mà Claude Code đã báo). Khi một bản ghi giới hạn ghi đè lên các con số, một dòng như `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` sẽ theo sau | tắt |
| `--no-cache` | Bỏ qua bộ nhớ đệm 60 giây và hỏi lại endpoint | tắt |
| `--source live\|claude\|cache` | `live` hỏi endpoint bằng token của tài khoản; `claude` để Claude Code tra cứu (bên dưới); `cache` chỉ hiển thị các con số đã lưu | `usage_source`, nếu không thì `live` |
| `--refresh` | Gia hạn các access token đã hết hạn, chỉ cho những tài khoản không đang chạy. Một lần gia hạn có thể mất tới 30 giây | tắt |
| `--dry-run` | Được chấp nhận: `ccam usage` không thay đổi gì ngoài bộ nhớ đệm của chính nó. Bị từ chối khi đi cùng `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | tắt |

Mọi giá trị `--source` khác đều thất bại với `--source must be live, claude or cache`. `ccam usage` không nhận đối số.

### Mức sử dụng từ chính Claude Code

Với `usage_source = "claude"` ở đầu `config.toml` (phía trên `[[group]]`, `[[profile]]` hoặc `[[folder]]` đầu tiên), hoặc `--source claude`, ccam hoàn toàn không đọc token đăng nhập. Với mỗi tài khoản, nó khởi động `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` trong thư mục của tài khoản và gửi yêu cầu `get_usage` của Claude Code qua đầu vào chuẩn (Claude Code 2.1.270 trở lên). Claude Code tra cứu mức sử dụng bằng thông tin đăng nhập của chính nó và trả lời với khoảng 5 giờ, khoảng hàng tuần và khoảng theo mô hình.

`--safe-mode` ngăn các hook, máy chủ MCP và plugin của tài khoản khởi động, và không có phiên nào được lưu. Mỗi tài khoản cần một lần khởi động Claude Code (khoảng một giây), nên nguồn này chậm hơn `live`. Các tài khoản được hỏi cùng lúc, mỗi tài khoản với thời gian chờ tối đa 20 giây, và bộ nhớ đệm 60 giây vẫn áp dụng như thường.

Một thông tin đăng nhập có access token đã hết hạn, hoặc sẽ hết hạn trong vòng 5 phút, được cho 50 giây, vì Claude Code gia hạn nó khi khởi động và một lần gia hạn bị cắt ngang có thể làm mất thông tin đăng nhập. Nếu ccam bỏ cuộc một lần tra cứu, nó yêu cầu Claude Code thoát bằng SIGTERM, để Claude Code hoàn tất lần gia hạn đã bắt đầu (Windows không có SIGTERM, nên ở đó tiến trình bị chấm dứt).

Một thông tin đăng nhập đã hết hạn hoàn toàn hiển thị `` login expired → `ccam login <name>` `` mà không khởi động Claude Code. Claude Code đánh dấu yêu cầu này là thử nghiệm, nên nếu nó ngừng trả lời, ccam hiển thị các con số đã lưu và `ccam doctor` báo điều đó. Khi Claude Code trả lời rằng lần tra cứu của chính nó thất bại, ccam hiển thị các con số đã lưu kèm `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`, chứ không phải một gói không có giới hạn.

Với `--dry-run`, ccam không khởi động Claude Code, vì Claude Code ghi vào thư mục của tài khoản; thay vào đó nó hiển thị các con số đã lưu. `usage_source` nhận `live` hoặc `claude`; `cache` chỉ là một lựa chọn của `--source`.

`ccam statusline` đọc JSON dòng trạng thái của Claude Code từ đầu vào chuẩn. Tùy chọn duy nhất của nó là `--wrap <command>`: một lệnh dòng trạng thái có sẵn, dòng đầu ra đầu tiên của lệnh đó sẽ được thêm tên tài khoản vào phía trước. `ccam statusline install` không nhận tùy chọn. Cả hai đều từ chối `--dry-run` với `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Token và chi phí: `ccam ccusage`

`ccam usage` cho biết một tài khoản đã dùng bao nhiêu trong mỗi giới hạn. Còn [ccusage](https://ccusage.com) thì đếm token và chi phí ước tính trong nhật ký hội thoại của Claude Code. Khi chạy riêng, nó chỉ tìm nhật ký của Claude Code trong `~/.claude` và `~/.config/claude` (hoặc trong các thư mục mà `CLAUDE_CONFIG_DIR` chỉ định), nên trong số các tài khoản của ccam, nó thấy nhiều nhất một tài khoản: tài khoản được nhập tại chỗ ở `~/.claude`. `ccam ccusage` chạy nó với `CLAUDE_CONFIG_DIR` được đặt thành thư mục của mọi tài khoản có hội thoại, để một báo cáo bao quát tất cả:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam đọc các tùy chọn của chính nó trước: `-a` (`--account`) nhận tên hoặc số của tài khoản, phân tách bằng dấu phẩy hoặc lặp lại tùy chọn, và `--dry-run` in ra lệnh thay vì chạy nó. Kể từ đối số đầu tiên không thuộc các tùy chọn này, mọi thứ được chuyển cho ccusage đúng như bạn gõ, trừ `--config` khi chạy bằng `bun x` và `npx` (xem bên dưới); `ccam ccusage -- --help` hiển thị trợ giúp riêng của ccusage. ccam chạy `ccusage` đã cài đặt, nếu không có thì `bun x ccusage`, nếu vẫn không có thì `npx ccusage@latest`; hai cách sau tải ccusage từ npm. Hai cách này khởi động trong thư mục riêng của ccam (`~/.local/state/ccam/ccusage`), không phải trong thư mục bạn đang dùng, nên một dự án ở đó không thể chọn thứ chúng chạy bằng `node_modules/.bin` của nó, một `bin` trong `package.json` của nó hay một registry trong `.npmrc` của nó; `ccam --dry-run ccusage` hiển thị lệnh `cd` vào thư mục của ccam. Khi không có cách nào, ccam bảo bạn cài Node.js, Bun hoặc ccusage. Trên Linux và macOS, ccusage thế chỗ ccam; trên Windows, ccam chờ nó và thoát với mã thoát của nó. Bên trong WSL, ccam bỏ qua `ccusage`, `bun` hoặc `npx` nằm trên ổ đĩa Windows, vì chúng sẽ chạy như một chương trình Windows mà không có các thư mục tài khoản.

- Báo cáo bao gồm các tài khoản có hội thoại trong thư mục của mình (một thư mục `projects`). Một cuộc hội thoại mà các tài khoản dùng chung qua nhóm chỉ được tính một lần.
- Với `-a`, một tài khoản dùng chung lịch sử qua nhóm sẽ kéo theo hội thoại của cả nhóm, và ccam cho biết đó là hội thoại của những ai: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- `blocks` hiển thị các khoảng 5 giờ, và mỗi tài khoản có khoảng riêng, nên hãy chỉ định một tài khoản cho báo cáo này, như ở trên.
- ccusage còn báo cáo các công cụ lập trình khác mà nó tìm thấy, như Codex, trừ khi bạn chỉ định một công cụ: `ccam ccusage claude daily` chỉ gồm Claude Code.
- Khi được ccam chạy, ccusage chỉ đọc cài đặt của nó (`ccusage.json`) từ `./.ccusage` và các thư mục tài khoản. `--config` chỉ định một tệp khác: `ccam ccusage daily --config ~/.config/claude/ccusage.json`. Vì `bun x` và `npx` khởi động trong thư mục riêng của ccam, ccam truyền `./.ccusage/ccusage.json` cho chúng bằng `--config` khi tệp đó tồn tại và bạn không chỉ định tệp nào, và chuyển một `--config` tương đối thành đường dẫn đầy đủ. Khi đó ccusage đọc các đường dẫn tương đối khác, chẳng hạn `--pi-path`, từ thư mục của ccam, nên hãy ghi chúng bằng đường dẫn đầy đủ.
- ccusage sẽ đọc một thư mục có dấu phẩy trong đường dẫn thành hai thư mục, nên ccam từ chối thư mục đó và nêu tên tài khoản; hãy loại tài khoản đó ra bằng `-a`.
- `ccusage` là một lệnh của ccam, nên tài khoản tên ccusage được khởi động bằng `ccam run ccusage`; `ccam doctor` có nhắc điều này.

## Cách hoạt động

Các con số mà `ccam usage` hiển thị đến từ hai nơi: endpoint mức sử dụng không được công bố tài liệu của Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, được gọi bằng thông tin đăng nhập riêng của từng tài khoản, và các giá trị mà Claude Code trao cho lệnh dòng trạng thái của nó. Toàn bộ nằm trong `internal/usage`.

1. ccam tải `config.toml` và tạo một `usage.Account` cho mỗi tài khoản (`cli.usageCmd`): thư mục `~/.local/share/ccam/profiles/<name>/` của nó, các tệp đánh dấu phiên của nó (cùng với các bản ghi riêng của Claude Code trong thư mục tài khoản, chúng cho biết nó có đang chạy không) và việc nó có được nhập tại chỗ hay không.
2. ccam tra cứu mọi tài khoản cùng lúc, một goroutine cho mỗi tài khoản (`Fetcher.FetchAll`). Tất cả dùng chung một `http.Client`, nên các kết nối tới endpoint được dùng lại. Việc kết nối có thể mất 2 giây, và mỗi tài khoản có tổng cộng 3 giây. Client không đi theo chuyển hướng, vì một lần chuyển hướng sẽ gửi token tới nơi khác, và chỉ đọc tối đa 1 MB của phản hồi.
3. Khi từng tài khoản hoàn tất, hàng của nó được vẽ lại. Tài khoản thất bại hoặc không phản hồi kịp sẽ quay về các con số đã lưu, nên một tài khoản chậm không bao giờ giữ chân các tài khoản khác.

Trong sơ đồ, mỗi tài khoản được tra cứu trong goroutine riêng, và tài khoản bị quá thời gian sẽ hiển thị các con số đã lưu thay vào đó:

<!-- motion:usage-fetch -->

`personal` và `work` phản hồi kịp và hiển thị số liệu trực tiếp. `side` mất hơn 3 giây, nên ccam hiển thị các con số đã lưu của nó thay vào đó (`fallback()`). Sau đó ccam vẽ một khối cho mỗi tài khoản (`usage.Table()`) và chọn tài khoản được nêu tên ở dòng cuối (`usage.Best()`).

Với một tài khoản, ccam thực hiện các kiểm tra sau theo thứ tự (`Fetcher.Fetch`):

<!-- motion:usage-cache -->

1. Với `--source cache`, nó đi thẳng tới `fallback()`.
2. Trừ khi đặt `--no-cache`, một kết quả trực tiếp trong `~/.local/state/ccam/usage/<name>.json` được đo cách đây chưa đến 60 giây (`CacheTTL`) sẽ được dùng, được đánh giá theo thời điểm hiện tại: một khoảng đã được đặt lại kể từ đó hiển thị 0%, và một tệp có ngày sau giờ của đồng hồ không được coi là mới. Chạy `ccam usage` hai lần liên tiếp chỉ tốn một yêu cầu cho mỗi tài khoản.
3. ccam đọc `.credentials.json` (hoặc Chuỗi khóa (Keychain) của macOS) mà không thay đổi nó, và chỉ giải mã khóa `claudeAiOauth` (`readCredentials()`). Không có thông tin đăng nhập gói đăng ký nghĩa là `not logged in`, hoặc `no limits (API)` khi `.claude.json` chứa một khóa API. Một thông tin đăng nhập có refresh token đã hết hạn, hoặc đã bị Claude Code xóa, là `` login expired → `ccam login <name>` `` và không được tra cứu.
4. Access token đã hết hạn không được dùng. Không có `--refresh`, tài khoản hiển thị các con số đã lưu kèm ``access token expired; `ccam <name>` renews it when it starts``.
5. Nếu không, ccam gọi endpoint. Mã 401 hoặc 403 được tính là đăng nhập hết hạn (`login expired: authentication failed (401)`), mọi trạng thái khác hoặc nội dung không đọc được được tính là lỗi, và cả hai đều quay về các con số đã lưu. Một lần tra cứu bị kết thúc bằng Ctrl+C báo `cancelled`. Mã 429 nghĩa là endpoint đã nhận quá nhiều yêu cầu cho tài khoản này: ccam ghi `<name>.pause` với thời điểm được hỏi lại (`Retry-After`, hoặc 5 phút, tối đa một giờ), hiển thị các con số đã lưu kèm `rate limited by the usage endpoint; asking again after 15:10`, và không gọi endpoint cho tài khoản đó cho đến lúc ấy, kể cả với `--no-cache`. Thử lại ngay lập tức chỉ khiến giới hạn tần suất kéo dài thêm.
6. Một phản hồi hợp lệ được ghi vào `<name>.json` cùng thời điểm đo.

Với các con số đã lưu, ccam lấy lần tra cứu đã lưu mới nhất (`<name>.json`, hoặc `<name>.claude.json` từ các lần tra cứu của Claude Code) và, khi `<name>.statusline.json` được ghi sau nó, đặt khoảng 5 giờ và khoảng hàng tuần của dòng trạng thái đè lên (`fallback()`); các khoảng theo mô hình, thứ mà dòng trạng thái không có, vẫn giữ như lần tra cứu đã lưu. Trước khi hiển thị các con số đã lưu, nó kiểm tra thời điểm đặt lại của từng khoảng: một khoảng đã được đặt lại kể từ khi các con số được lưu sẽ hiển thị 0%, vì tỷ lệ phần trăm cũ không còn đúng. Đó là lý do `personal` ở trên hiển thị `0%` mà không có thời gian đặt lại.

Endpoint có hai định dạng phản hồi. Định dạng mới hơn có một mảng `limits` với các mục `session` (5 giờ), `weekly_all` và `weekly_scoped` (theo từng mô hình). Định dạng cũ hơn có các đối tượng `five_hour`, `seven_day`, `seven_day_opus` và `seven_day_sonnet`. ccam đọc cả hai (`ParseLive`): mảng `limits` được ưu tiên, và các đối tượng cũ bổ sung những gì mảng đó không có. Một nội dung hoàn toàn không có giới hạn nào được báo là định dạng đã thay đổi, và `ccam doctor` sẽ báo cho bạn điều đó.

### Bản ghi giới hạn

Khi Claude Code dừng một lượt ở giới hạn sử dụng, hook của [plugin ccam](/docs/install#claude-code-plugin) đọc yêu cầu bị từ chối từ cuộc hội thoại: khoảng nào đã đầy (5 giờ, hàng tuần, Opus hoặc Sonnet) và khi nào nó được đặt lại. Nó ghi điều đó vào `~/.local/state/ccam/usage/<name>.limit.json`. ccam áp dụng bản ghi chồng lên mọi thứ khác nó biết, từ mọi nguồn (`Fetch`): mỗi khoảng trong đó được tính là 100% cho đến khi được đặt lại, bất kể các con số đã lưu cũ hay thấp đến đâu, nên `ccam next --auto`, `ccam pick`, `--fallback` và [`auto_switch`](/features/next#auto-switch) không chọn lại tài khoản trước lúc đó.

Khi không biết gì khác, chỉ riêng bản ghi là thứ tài khoản hiển thị (`source: the usage limit Claude Code reported` với `-v`). Các phiên không có người theo dõi (các lần chạy `-p` và các phiên nền) không ghi bản ghi nào, việc dừng do tổ chức đặt ra cũng vậy (ví dụ khi tổ chức hết mức sử dụng hoặc giới hạn của một nhóm là $0), điều mà quản trị viên của nó có thể gỡ bất cứ lúc nào. `ccam login` vào một tài khoản Claude khác sẽ bỏ bản ghi; xóa tệp cũng xóa được nó.

### Gia hạn token

Refresh token được xoay vòng: gia hạn một token sẽ cho bạn một refresh token mới và vô hiệu hóa token cũ. Nếu ccam gia hạn thông tin đăng nhập trong lúc Claude Code đang chạy bằng cùng tài khoản đó, phiên đó sẽ bị đăng xuất. Vì vậy `ccam usage` đọc token và không bao giờ ghi chúng, trừ khi bạn truyền `--refresh`. Khi đó ccam (`Fetcher.refresh`):

1. từ chối các tài khoản được đăng ký tại chỗ bằng `ccam import`, vì `claude` thông thường có thể dùng thông tin đăng nhập đó mà ccam không biết;
2. từ chối thông tin đăng nhập được lưu trong Chuỗi khóa (Keychain) của macOS;
3. lấy một tệp khóa, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, để hai tiến trình ccam không bao giờ cùng gia hạn một token; trong khi khóa đang bị giữ, tài khoản hiển thị `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)`. `ccam rename`, `remove` và `logout` cũng lấy khóa này, và từ chối khi một lần gia hạn đang diễn ra;
4. trong khi giữ khóa, từ chối nếu có một phiên của tài khoản đang chạy, bất kể nó được khởi động theo cách nào (qua ccam, `claude` thông thường, một IDE, hoặc một phiên nền và daemon đứng sau nó, thứ tự gia hạn thông tin đăng nhập), và bỏ qua yêu cầu nếu một tiến trình khác đã gia hạn token;
5. không gửi gì khi Claude Code đang giữ các khóa gia hạn của riêng nó trong thư mục tài khoản (`.oauth_refresh.lock`, `<folder>.lock`): lúc đó Claude Code đang gia hạn chính thông tin đăng nhập ấy;
6. gửi refresh token tới endpoint token và chờ phản hồi tối đa 30 giây, khoảng thời gian Claude Code dành cho các lần gia hạn của chính nó. Một yêu cầu khi đã gửi đi thì không bao giờ bị cắt ngang, kể cả bằng Ctrl+C: máy chủ có thể đã xoay vòng refresh token, và chỉ phản hồi của nó mang thông tin đăng nhập mới. Ctrl+C hoặc tín hiệu hangup kết thúc các lần tra cứu khác, và ccam thoát (130 sau Ctrl+C) khi phản hồi đã được lưu;
7. lưu dưới khóa ghi của Claude Code (`.storage-write.lock`) và chỉ thay đổi các trường token của `.credentials.json` ở trạng thái lúc đó: những gì Claude Code đã ghi trong lúc này (một lần đăng nhập MCP) vẫn được giữ, một thông tin đăng nhập mà Claude Code đã gia hạn trong lúc này được giữ, và một lần đăng xuất không bao giờ bị hoàn tác. Khi phản hồi cho biết thời hạn của refresh token mới, nó được lưu dưới dạng `refreshTokenExpiresAt`, như Claude Code làm.

Một lần gia hạn không thành công sẽ cho biết lý do, trên dòng của tài khoản hoặc phía trên các con số đã lưu của nó: chẳng hạn `token not renewed: a session is running`, hoặc `login may be lost: …` khi một yêu cầu đã được gửi nhưng phản hồi của nó không dùng hoặc không lưu được, kèm lời khuyên chạy `ccam login <name>` nếu sau đó Claude Code yêu cầu bạn đăng nhập.

### Bộ thu thập từ dòng trạng thái

Claude Code chạy lệnh dòng trạng thái và pipe cho nó một tài liệu JSON. ccam đọc tên mô hình và `rate_limits` từ đó:

<!-- motion:usage-statusline -->

1. ccam tìm tài khoản từ `CCAM_PROFILE`, biến mà ccam đặt khi khởi động Claude Code, hoặc nếu không có thì bằng cách so khớp `CLAUDE_CONFIG_DIR` với các thư mục tài khoản (`runStatusline`).
2. ccam đọc `rate_limits.five_hour` và `rate_limits.seven_day` (`ParseStatusLine`). JSON không có giới hạn tần suất đối với tài khoản dùng khóa API, hoặc trước phản hồi đầu tiên của phiên; khi đó không có gì được lưu.
3. ccam ghi `~/.local/state/ccam/usage/<name>.statusline.json` (`SaveStatusLine`). `ccam usage` dùng tệp này khi endpoint không phản hồi.
4. ccam in ra `[name] · model · 5h N% · 7d N%`. Với `--wrap`, nó chạy lệnh của bạn qua `sh -c` với cùng đầu vào và in ra đầu ra của nó theo cách Claude Code hiển thị một dòng trạng thái: mọi dòng không trống, với `[name]` đứng trước dòng đầu tiên. Một lệnh thất bại chỉ hiển thị tên tài khoản. Trên Windows, nó chạy lệnh theo cách Claude Code làm: bằng Git Bash khi đã được cài, nếu không thì bằng PowerShell.
5. Khi mức sử dụng 5 giờ hoặc hàng tuần ở mức 90% trở lên và có một tài khoản khác, dòng này kết thúc bằng `→ ccam next --auto --continue`, ví dụ `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Lệnh được bao bọc không có giới hạn thời gian riêng: Claude Code cho một lệnh dòng trạng thái 10 phút và dừng một lần chạy khi một lần cập nhật mới hơn bắt đầu, nên một lệnh chậm hiển thị giống như khi không có ccam. Khi Claude Code kết thúc ccam, lệnh cũng kết thúc (trên Linux và macOS là cùng mọi tiến trình nó đã khởi động).

`ccam statusline install` ghi `statusLine` vào `settings.json` của từng tài khoản và giữ các khóa khác của đối tượng, như `padding`, `refreshInterval` và `hideVimModeIndicator`. Khi một tài khoản đã có lệnh dòng trạng thái, ccam giữ nó dưới dạng `--wrap '<old command>'`; trên Windows là `--wrap-base64 <base64>`, cách không cần đặt trong ngoặc ở shell nào. Nếu dòng trạng thái đã chạy ccam (từ đường dẫn khác, hoặc được đặt trong ngoặc theo cách một phiên bản cũ hơn đã ghi), ccam thay thế nó thay vì bao bọc, và giữ lại lệnh mà nó bao bọc.

Trên Windows, đường dẫn tới ccam được ghi bằng dấu gạch chéo xuôi và `~` cho thư mục home (`~/.local/bin/ccam.exe statusline`), cách viết mà cả Git Bash lẫn PowerShell đều hiểu. Một đường dẫn cần dấu ngoặc được ghi đầy đủ cho shell mà Claude Code dùng: `'…'` cho Git Bash, `& '…'` cho PowerShell. Một `settings.json` được dùng chung qua nhóm chỉ được ghi một lần, qua liên kết của nó, nên mọi tài khoản trong nhóm đều nhận được. Chạy lại lệnh sẽ in ra `(already set)` cho các tệp đã dùng ccam.

## An toàn và các trường hợp đặc biệt

- `ccam usage` không bao giờ thay đổi tệp token trừ khi bạn truyền `--refresh`. Khi một access token đã hết hạn, hãy khởi động tài khoản một lần bằng `ccam <name>` và Claude Code sẽ gia hạn nó. Một thông tin đăng nhập đã hết hạn hoàn toàn cần `ccam login <name>`.
- `--refresh` không bao giờ động đến tài khoản đang chạy (kể cả khi Claude Code được khởi động không qua ccam), tài khoản đã nhập hay thông tin đăng nhập trong Chuỗi khóa (Keychain). Một lần gia hạn bị từ chối sẽ hiển thị các con số đã lưu và lý do (`token not renewed: …`).
- Endpoint mức sử dụng không được công bố tài liệu chính thức và có thể thay đổi. Khi nó thay đổi, bạn vẫn thấy các con số đã lưu, và `ccam doctor` báo việc định dạng đã thay đổi.
- Khi không có con số đã lưu và không có phản hồi, dòng của tài khoản đó hiển thị `lookup failed:` kèm lý do, ví dụ `no saved usage` với `--source cache`.
- Các tệp bộ nhớ đệm chỉ chứa tỷ lệ phần trăm và thời điểm đặt lại, không bao giờ chứa token.

## Liên quan

- [Giới hạn sử dụng và chuyển đổi tài khoản](/docs/usage-limits)
- [Chuyển sang tài khoản tiếp theo](/features/next)
- [Làm gì khi Claude Code chạm giới hạn sử dụng](/guides/usage-limit-workflow)
- [Tham chiếu lệnh](/docs/commands)
