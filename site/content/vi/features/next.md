---
title: Chuyển sang tài khoản tiếp theo
description: "ccam next mở Claude Code bằng tài khoản khác; --auto chọn tài khoản còn nhiều hạn mức nhất, --continue chuyển hội thoại sang, --fallback chạy lại tác vụ -p."
updated: 2026-11-03
order: 8
translation: machine
sourceHash: "89e8dea9c55d"
---

`ccam next` mở Claude Code bằng một tài khoản khác mà bạn không cần nêu tên. Hãy dùng nó khi một tài khoản chạm giới hạn sử dụng. Có bốn cách để chuyển đổi:

- `ccam next` không kèm tùy chọn chọn tài khoản đứng sau tài khoản bạn dùng gần nhất, còn `--auto` chọn tài khoản còn nhiều hạn mức nhất. `--continue` tiếp tục cuộc hội thoại gần nhất của thư mục này trên tài khoản mới, dù hai tài khoản có dùng chung lịch sử hay không.
- Bên trong Claude Code, [plugin ccam](#inside-claude-code) báo cho bạn khi một giới hạn dừng một phiên, và `/ccam:switch` chuyển cuộc hội thoại đi tiếp.
- Với các script chạy `claude -p`, `ccam --fallback` tự thực hiện việc đổi (xem [Chạy lại khi chạm giới hạn](#run-again-at-the-limit)).
- Khi bật [`auto_switch`](#auto-switch) (Linux và macOS), ccam còn chuyển một phiên tương tác mà Claude Code dừng lại ở giới hạn của nó.

**Seat của gói Team và Enterprise.** Việc tự động chọn tài khoản (`ccam next`, `--fallback`, `ccam pick`, phím `n` trên bảng điều khiển và, khi bạn bật nó, [`auto_switch`](/features/next#auto-switch)) cũng chuyển đổi giữa các seat của cùng một tổ chức Claude Team hoặc Enterprise, điều này có thể vi phạm điều khoản của Anthropic. Hãy đọc [Câu hỏi thường gặp](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) và xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó.

## Cách dùng

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

ccam in việc chuyển đổi dưới dạng `→ <from> → <to>` ra đầu ra lỗi chuẩn, kèm lý do của một lựa chọn `--auto`, rồi khởi động Claude Code bằng tài khoản mới theo cùng cách `ccam <name>` làm. `<from>` là tài khoản bạn dùng gần nhất, hoặc với `--continue` là tài khoản sở hữu cuộc hội thoại. Các đối số sau `--` được chuyển nguyên vẹn cho `claude`.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--auto` | Chọn tài khoản đã đăng nhập theo mức sử dụng (xem bên dưới), bỏ qua các tài khoản ở mức 100% | tắt |
| `--max-usage <1-100>` | Với `--auto`: ưu tiên các tài khoản có mức sử dụng dưới giá trị này | `auto_max_usage`, nếu không thì 100 |
| `--strategy headroom\|expiring\|balance` | Với `--auto`: mức sử dụng 5 giờ thấp nhất, khoảng 5 giờ được đặt lại sớm nhất, hoặc nhiều hạn mức hàng tuần còn lại nhất cho mỗi ngày | `auto_strategy`, nếu không thì `headroom` |
| `--cross-org` | Với `--auto`: chọn cả các tài khoản thuộc tổ chức khác (xem bên dưới) | `cross_org`, nếu không thì tắt |
| `--wait[=<duration>]` | Với `--auto`: khi không tài khoản nào còn hạn mức, chờ tài khoản đầu tiên được đặt lại, tối đa chừng này thời gian. Giá trị đặt sau `=`: `--wait 2h` bị từ chối kèm gợi ý `--wait takes its value after '=': --wait=2h` | tắt; `5h` khi chỉ ghi cờ |
| `-c`, `--continue` | Tiếp tục cuộc hội thoại gần nhất của thư mục này trên tài khoản mới (xem bên dưới) | tắt |
| `--brief` | Cùng với `--continue`: bắt đầu một cuộc hội thoại mới dưới tài khoản mới từ bản tóm tắt mà ccam viết về cuộc hội thoại trước, thay vì tiếp tục toàn bộ nó ([xem bên dưới](#brief)) | tắt |
| `-- <claude args>` | Các đối số được chuyển cho `claude` sau các đối số của chính ccam | không có |
| `--dry-run` | Được hỗ trợ: in ra `→ <from> → <to>` và dòng lệnh sẽ chạy, mà không khởi động Claude Code hay ghi lại lần khởi động | tắt |

`--auto` và `--continue` dùng được cùng nhau. Với `--continue`, các đối số `--resume` được đặt trước các đối số của bạn.

Các tùy chọn riêng của Claude Code, như `--dangerously-skip-permissions` hoặc `--model`, đặt sau `--`. Khác với `ccam <name>`, vốn chuyển tiếp mọi thứ sau tên, `ccam next` tự đọc các đối số của nó, nên nếu thiếu `--` nó sẽ dừng với `ccam: unknown flag: --dangerously-skip-permissions`.

## Cách hoạt động

Mỗi lần khởi động qua ccam đều ghi thời điểm vào `~/.local/state/ccam/last-used/<name>`. `ccam next` đọc các tệp này, kiểm tra những tài khoản nào đã đăng nhập, rồi chọn.

1. ccam tạo một `next.Candidate` cho mỗi tài khoản từ `config.toml` (`cli.candidates`): tài khoản, thời điểm khởi động gần nhất của nó (`launcher.LastUsed`) và việc nó có thông tin đăng nhập còn hoạt động hay không (`creds.Has`: một thông tin đăng nhập Claude trong `.credentials.json`, hoặc trên macOS là trong Chuỗi khóa (Keychain), chưa hết hạn hoàn toàn và chưa bị Claude Code xóa). Một thư mục dùng API key thay vào đó (một lần đăng nhập Console, hoặc một tài khoản được thêm bằng `ccam add --api`) không có giới hạn gói và không bao giờ được chọn.
2. Không có `--auto`, ccam sắp xếp các tài khoản theo số slot và tìm tài khoản được khởi động gần nhất (`next.Rotate`). Đó là tài khoản bạn đang chuyển đi. Nếu chưa có tài khoản nào được khởi động, đó là slot thấp nhất.
3. Từ đó nó đi tiếp theo thứ tự slot, quay vòng từ slot cuối về slot đầu, và trả về tài khoản đầu tiên đã đăng nhập.
4. ccam chuẩn bị và khởi động tài khoản được chọn đúng như `ccam <name>` (`launcher.LaunchWith`): nó kiểm tra các liên kết của tài khoản, áp dụng việc dùng chung MCP, đăng ký phiên và thay ccam bằng `claude` qua `syscall.Exec`.

Vòng duyệt bắt đầu ở tài khoản dùng gần nhất và dừng ở tài khoản đã đăng nhập đầu tiên sau nó:

<!-- motion:next-rotate -->

`work` được khởi động gần nhất, nên vòng duyệt bắt đầu sau nó. `side` chưa đăng nhập nên bị bỏ qua, và `lab` được chọn. Nếu `lab` cũng chưa đăng nhập, vòng duyệt sẽ quay về `personal` (`% len` trong sơ đồ). Vòng xoay không bao giờ trả về tài khoản bạn bắt đầu: khi không có tài khoản nào khác đã đăng nhập, nó dừng lại với một lỗi.

### --auto

`--auto` quyết định theo mức sử dụng thay vì thứ tự slot. ccam giữ lại các tài khoản đã đăng nhập, hỏi mức sử dụng của tất cả cùng lúc theo cùng cách `ccam usage` làm (`Fetcher.FetchAll`), và chọn theo các quy tắc bên dưới (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

Một tài khoản được đánh giá theo khoảng 5 giờ và khoảng hàng tuần của nó, và khi `--model` nằm trong các đối số cho Claude Code (sau `--`), theo cả khoảng hàng tuần của mô hình đó (`opus`, `sonnet`, `haiku` hoặc `fable` trong tên). ccam chia các tài khoản thành ba bậc (`Policy.Choose()`):

1. **Còn hạn mức**: mọi khoảng đều dưới `auto_max_usage` (`config.toml`, mặc định 100). Chiến lược chọn trong số này: `headroom` (mặc định) lấy mức sử dụng 5 giờ thấp nhất; `expiring` lấy khoảng 5 giờ được đặt lại sớm nhất, để hạn mức còn lại được dùng trước khi mất; `balance` lấy tài khoản còn nhiều hạn mức hàng tuần nhất cho mỗi ngày cho đến khi tuần của nó được đặt lại, có tính cả từng phiên ccam đang chạy trên tài khoản đó, để hạn mức hàng tuần không bị bỏ phí và các phiên song song được phân tán.
2. **Vượt ngưỡng**: dưới 100% nhưng bằng hoặc vượt `auto_max_usage`. Chỉ được dùng khi không tài khoản nào còn hạn mức; tài khoản có mức sử dụng 5 giờ thấp nhất được ưu tiên.
3. **Đầy**: có một khoảng ở mức 100%, hoặc chưa rõ mức sử dụng 5 giờ. Không bao giờ được chọn.

Một khoảng được tính là 100% cho đến khi được đặt lại khi Claude Code đã từ chối một yêu cầu ở đó: hook của [plugin ccam](#inside-claude-code) ghi lại điều này (`~/.local/state/ccam/usage/<name>.limit.json`), bất kể các con số đã lưu nói gì. Vì vậy một tài khoản vừa chạm giới hạn sẽ không được chọn lại trước khi đặt lại, kể cả khi các con số của nó đã cũ.

Cài đặt cho mọi tài khoản đặt ở đầu `config.toml`, phía trên `[[group]]`, `[[profile]]` hoặc `[[folder]]` đầu tiên. Một cài đặt nằm dưới một trong số đó thuộc về bảng đó: ccam cảnh báo về nó mỗi lần khởi động, và các lệnh lưu cấu hình sẽ từ chối cho đến khi nó được dời lên.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Tổ chức.** Các lựa chọn tự động chỉ giữ ở những tài khoản cùng loại với tài khoản bạn đang rời đi: các seat trong cùng tổ chức Team hoặc Enterprise, hoặc các gói cá nhân (Pro, Max), để công việc từ một tài khoản công ty không tiếp tục trên một tài khoản riêng tư hay ngược lại. ccam biết tổ chức của mỗi tài khoản từ thông tin đăng nhập của nó (`ccam login` ghi lại; các lần đăng nhập cũ hơn được đọc từ `.claude.json`); một tài khoản mà ccam không xác định được tổ chức thì có thể ghép với bất kỳ tài khoản nào. Hai thư mục đăng nhập vào cùng một tài khoản Claude dùng chung giới hạn của nó, nên thư mục này không bao giờ thay cho thư mục kia. `--cross-org`, hoặc `cross_org = true` trong `config.toml`, bỏ quy tắc tổ chức này, nhưng sau [`/ccam:switch`](#inside-claude-code), việc chuyển sang một tổ chức khác vẫn phải chờ bạn trả lời yes.

Khi không chọn được tài khoản nào, ccam thoát với mã 75 (`EX_TEMPFAIL`), để script phân biệt được "thử lại sau" với một lỗi. Với `--wait`, thay vào đó nó ngủ cho đến khi khoảng đang chặn của tài khoản đầy đầu tiên được đặt lại (30 giây sau thời điểm đặt lại, khi endpoint hiển thị khoảng mới), in ra `ccam: no account has room; waiting until 15:10`, rồi chọn lại; khi không tài khoản nào cho biết lúc nào được đặt lại, nó kiểm tra lại sau 10 phút. Nếu không có tài khoản nào trống ra trước khi hết thời gian chờ, nó thoát với mã 75.

Việc tra cứu tuân theo các quy tắc của `ccam usage`: các con số chưa quá 60 giây được dùng lại, lượt tra cứu của mỗi tài khoản hết thời gian chờ sau 3 giây, tài khoản không phản hồi hoặc có access token đã hết hạn sẽ được đánh giá theo các con số đã lưu, và một bản ghi giới hạn ghi đè lên chúng. Một khoảng đã lưu mà thời điểm đặt lại đã qua được tính là 0%, nên tài khoản có khoảng giới hạn đã được đặt lại kể từ đó sẽ không bị bỏ qua. `--auto` không bao giờ gia hạn token. Khác với vòng xoay thông thường, `--auto` có thể chọn tài khoản bạn dùng gần nhất nếu nó vẫn còn nhiều hạn mức nhất.

### --continue

1. ccam tìm cuộc hội thoại gần nhất của thư mục này trong mọi tài khoản: trong các tệp `<session>.jsonl` ở `projects/<folder>/` (tên thư mục là thư mục làm việc, với các liên kết đã được phân giải, và mọi ký tự không phải chữ cái hay chữ số được đổi thành `-`, theo cách Claude Code đặt tên), tệp có lời nhắc hoặc câu trả lời mới nhất. Các phiên không có lời nhắc hay câu trả lời nào, chẳng hạn một phiên chỉ chạy `/mcp`, `/usage` hoặc `/status`, hoặc một tệp trống, sẽ bị bỏ qua. Nếu một phiên để lại một ghi chú (một giới hạn sử dụng hoặc `/ccam:switch`, xem bên dưới) và không có cuộc hội thoại nào có lời nhắc hoặc câu trả lời sau đó hơn 2 phút, chính phiên đó được dùng; các dòng mà Claude Code ghi khi `/exit` không được tính.
2. Tài khoản sở hữu cuộc hội thoại là tài khoản bạn chuyển đi, và tài khoản mới được chọn như trên, trong cùng tổ chức. Khi nhiều tài khoản dùng chung lịch sử đó, cuộc hội thoại thuộc về tài khoản có `.claude.json` cho biết nó đã chạy cuộc hội thoại gần nhất (`lastSessionId`), nếu không thì thuộc về tài khoản được dùng gần nhất.
3. Nếu tài khoản mới dùng chung lịch sử đó (cùng thư mục nhóm), ccam khởi động `claude --resume <session id>`, tức chính cuộc hội thoại đó (`--continue` sẽ mở cuộc hội thoại nào mới nhất trong thư mục).
4. Nếu không, nó khởi động `claude --resume <the conversation's file> --fork-session`: Claude Code sao chép toàn bộ cuộc hội thoại, dưới một ID phiên mới, vào lịch sử riêng của tài khoản mới, và bản gốc vẫn giữ nguyên. Từ đó trở đi, `--continue` của tài khoản mới sẽ tìm thấy nó.
5. Trừ khi bạn đưa `--agent` sau `--`, ccam cũng truyền tiếp agent mà cuộc hội thoại đang chạy dưới đó: cài đặt `agent` trong `.claude/settings.local.json` hoặc `.claude/settings.json` của thư mục, hoặc trong `settings.json` của tài khoản, nếu không có thì là agent mà dòng `agent-setting` cuối cùng của cuộc hội thoại nêu tên. Claude Code không khôi phục agent khi sao chép một cuộc hội thoại (bước 4), và sẽ tiếp tục dưới agent mặc định của nó, với mọi công cụ. Khi tài khoản mới không dùng chung agent, plugin và cài đặt của tài khoản cũ (cùng các thư mục `agents` và `plugins` và `settings.json`, qua một nhóm), ccam thêm một dòng sau dòng về việc chuyển: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

ccam cho biết cuộc hội thoại lớn đến đâu: tài khoản mới đọc lại toàn bộ nó với câu trả lời đầu tiên. Bộ nhớ đệm prompt không bao giờ được dùng chung giữa các tổ chức, và mỗi gói cá nhân là một tổ chức riêng, nên khi hai tài khoản thuộc về các tổ chức khác nhau, ccam thêm `prompt caches are not shared between organizations, so the first reply reads it all again`. Khi tài khoản cũ đã đăng nhập vào những máy chủ MCP mà tài khoản mới chưa đăng nhập, ccam nêu tên chúng: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. Dùng chung `history` trong một nhóm tránh được việc sao chép, giữ một cuộc hội thoại duy nhất, và còn dùng chung các checkpoint mà `/rewind` dùng để khôi phục mã (`file-history/`).

### --brief

`--continue` tiếp tục toàn bộ cuộc hội thoại, và tài khoản mới đọc lại tất cả với câu trả lời đầu tiên. Thay vào đó, `ccam next --continue --brief` bắt đầu một cuộc hội thoại mới, từ một bản tóm tắt mà ccam viết về cuộc hội thoại cũ mà không hỏi mô hình nào:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

Bản tóm tắt chứa, theo thứ tự này và chỉ khi có điều gì để nói:

1. tiêu đề của cuộc hội thoại;
2. đường dẫn tới bản ghi (transcript) của nó, để Claude tìm khi cần một chi tiết;
3. yêu cầu đầu tiên của bạn, khi có nhiều hơn bốn yêu cầu;
4. bản tóm tắt mà Claude Code viết khi nén (compact) cuộc hội thoại lần gần nhất;
5. bốn yêu cầu gần nhất của bạn, với các lệnh `!` và lệnh gạch chéo đúng như bạn đã gõ;
6. câu trả lời gần nhất của Claude;
7. danh sách việc cần làm (todo);
8. 10 tệp mà Claude thay đổi gần nhất;
9. `git status --short --branch` của thư mục, tối đa 800 ký tự.

Danh sách việc cần làm giữ tối đa 10 mục: các mục đã xong bị bỏ trước, rồi đến các bước xa nhất phía trước, để bước hiện tại và các bước kế tiếp được giữ lại. Các phần dài bị cắt ở giữa, và toàn bộ bản tóm tắt dưới 16.000 ký tự.

ccam lưu nó thành `~/.local/state/ccam/brief/<session>.md`, chỉ bạn đọc được, và lần `--brief` tiếp theo xóa các bản tóm tắt cũ hơn 7 ngày. Claude Code khởi động với một prompt, và đính kèm tệp mà prompt nhắc tới: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` Một prompt bạn đưa sau `--` thay cho câu cuối (`Then: <your prompt>`), và các tùy chọn khác của bạn sau `--` đi theo sau nó.

Claude Code kết thúc một lần nhắc (mention) tại dấu `"` và cắt nó tại dấu `#`, nên khi đường dẫn của bản tóm tắt có một trong hai ký tự đó, ccam nhắc nó bằng đường dẫn tính từ thư mục làm việc; khi đường dẫn đó cũng có, prompt nêu tên tệp, và Claude đọc nó bằng công cụ Read, công cụ này hỏi trước vì tệp nằm ngoài các thư mục làm việc (trong một lần chạy `-p` ở chế độ Manual, Claude Code từ chối việc đọc và Claude tiếp tục mà không có bản tóm tắt).

`ccam --dry-run next --continue --brief` in bản tóm tắt và không lưu gì. Cuộc hội thoại cũ giữ nguyên, và cuộc hội thoại mới chỉ biết những gì bản tóm tắt nói, nên hãy dùng `--brief` khi nhiệm vụ quan trọng hơn mọi chi tiết về cách bạn đã đi tới đó. Thiếu `--continue`, ccam dừng với `--brief works with --continue`; khi không đọc được cuộc hội thoại hoặc không tìm thấy gì trong đó để tóm tắt, nó dừng và cho biết `ccam next --continue` không kèm `--brief` sẽ tiếp tục toàn bộ cuộc hội thoại. [Relay](#relay) và [`auto_switch`](#auto-switch) luôn tiếp tục toàn bộ cuộc hội thoại. Tùy chọn `--brief` của chính Claude Code đặt sau `--`, như các tùy chọn khác của nó.

## Bên trong Claude Code

[Plugin ccam](/docs/install#claude-code-plugin) thêm hai hook (Claude Code 2.1.78 trở lên):

- Khi một lượt dừng lại vì giới hạn sử dụng, nó ghi chú cuộc hội thoại và hiện một thông báo trên màn hình ở các terminal hỗ trợ (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  Nó nhận ra mọi thông báo mà Claude Code coi là giới hạn sử dụng: giới hạn 5 giờ, hàng tuần, giới hạn thông thường và giới hạn của mô hình (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` cho hạn mức của một mô hình), tín dụng sử dụng và giới hạn chi tiêu, và một tổ chức đã dừng tài khoản.

  Khi Claude Code đã ghi yêu cầu bị từ chối vào cuộc hội thoại, hook còn ghi lại khoảng đã đầy cho đến khi nó được đặt lại ([bản ghi giới hạn](/features/usage#the-limit-record)), trừ trường hợp tổ chức dừng tài khoản, điều mà quản trị viên của nó có thể gỡ bất cứ lúc nào. Một lần máy chủ điều tiết ngắn hoặc giới hạn tốc độ của API key không phải là giới hạn sử dụng và không bị đụng tới, các phiên không có người theo dõi (các lần chạy `-p` và các phiên nền) cũng vậy. Một subagent hoặc agent của workflow dừng lại ở giới hạn sử dụng được tính là phiên của nó dừng lại, kể cả khi cuộc hội thoại chính chỉ đang chờ nó: Claude Code ghi yêu cầu bị từ chối đó vào tệp riêng của agent trong thư mục `subagents` của phiên, và hook đọc nó ở đó. Khi nhiều agent dừng lại cùng lúc, phiên chỉ nhận một ghi chú và một thông báo.

  Claude Code 2.1.234 trở lên còn chờ trong phiên và tự tiếp tục sau khi khoảng 5 giờ được đặt lại; thông báo cho bạn lựa chọn chuyển tài khoản thay vì chờ, và khi bật [`auto_switch`](#auto-switch) (Linux và macOS), ccam đổi giúp bạn khi một tài khoản khác còn hạn mức. Khi không tài khoản nào khác được phép tiếp quản, vì không có tài khoản nào khác thuộc cùng tổ chức Team hoặc Enterprise, hoặc không có gói cá nhân nào khác (khi bật `cross_org`, mọi tài khoản Claude khác đều được), thông báo thay vào đó kết thúc bằng `You have no other personal plans, so ccam cannot move this conversation.` hoặc `You have no other accounts in Acme, so ccam cannot move this conversation.`
- `/ccam:switch [account]` ghi chú cuộc hội thoại và kết thúc lượt trước khi có yêu cầu nào được gửi: gõ `/exit`, rồi chạy `ccam next --continue` để tiếp tục nó trên tài khoản bạn đã nêu tên (hoặc trên tài khoản kế tiếp). Trong một phiên chưa có cuộc hội thoại nào, nó báo rằng không có gì để tiếp tục. Chỉ những gì bạn gõ mới chạy được nó: từ plugin 0.3.2, Claude không thể chạy `/ccam:switch`, và một tác vụ do Claude lên lịch (CronCreate, ScheduleWakeup hoặc `.claude/scheduled_tasks.json`) chuyển nó cho Claude dưới dạng văn bản thường.

  Khi tài khoản bạn nêu tên nằm ngoài tổ chức của cuộc hội thoại (một seat của tổ chức Team hoặc Enterprise khác, một gói cá nhân thay cho một seat hoặc ngược lại, hoặc một tài khoản API key ở một trong hai phía), `ccam next --continue` hỏi trước, bất kể `cross_org` nói gì: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. Nó cũng hỏi như vậy trước khi chuyển sang một tài khoản nó tự chọn khi bạn không nêu tên tài khoản nào, tài khoản mà `cross_org = true` cho phép nó lấy từ một tổ chức khác: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. Chỉ `y` hoặc `yes` mới chuyển nó; trên Linux và macOS, các phím đã gõ trước câu hỏi sẽ bị bỏ trước. Bất kỳ câu trả lời nào khác, hoặc không có terminal để hỏi, sẽ để cuộc hội thoại và ghi chú ở nguyên chỗ cũ: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (sau một lần tự chọn, `ccam next --continue` sẽ chọn lại). Giống như với các lựa chọn tự động, một tài khoản mà ccam không xác định được tổ chức thì có thể ghép với bất kỳ tài khoản nào.

Ghi chú nằm trong `~/.local/state/ccam/handoff/<account>/<session>.json`, mỗi phiên một tệp, nên hai phiên của cùng một tài khoản không bao giờ ghi đè lên nhau. Nó được giữ trong 12 giờ và chỉ được dùng một lần, khi phiên tiếp theo sắp bắt đầu; một lần khởi động thất bại sẽ để nguyên nó cho lần thử sau. Ghi chú thuộc về thư mục nơi phiên bắt đầu, kể cả sau một lệnh `cd` trong shell của Claude Code. `ccam rename` cập nhật các tên trong đó; nếu tài khoản mà ghi chú nêu tên đã bị xóa kể từ đó, ccam sẽ báo điều này và chọn như bình thường.

### Relay

Với `relay = true` ở đầu `config.toml`, ccam tiếp tục chạy song song với các phiên tương tác thay vì tự thay thế mình bằng Claude Code, và sau `/ccam:switch` nó tự tiếp tục cuộc hội thoại của phiên đó (ghi chú do một phiên khác để lại sẽ không được xử lý): hook kết thúc Claude Code, và ccam chạy `ccam next --continue` với ghi chú của phiên đó. Khi không nêu tên tài khoản, nó chọn tài khoản còn nhiều hạn mức nhất, có tính cả mô hình của cuộc hội thoại, và khi `cross_org = true` khiến tài khoản đó thuộc một tổ chức khác, nó hỏi câu hỏi ở trên trước khi chuyển cuộc hội thoại sang đó. Trên Windows, hook không thể kết thúc Claude Code: hãy gõ `/exit`, và ccam sẽ tiếp tục nó. Hook cũng không kết thúc nó khi tài khoản nằm ngoài tổ chức của cuộc hội thoại: hook báo `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`, và ccam hỏi câu hỏi ở trên sau khi bạn đã gõ `/exit`. Việc này tốn một tiến trình ccam đang chờ cho mỗi phiên, đó là lý do nó tắt theo mặc định.

Phiên được chuyển tiếp giữ lại những gì việc tiếp tục một cuộc hội thoại không khôi phục: các tùy chọn như `--add-dir`, `--mcp-config`, danh sách công cụ, `--settings`, `--plugin-dir` và các tùy chọn system prompt, chế độ quyền của cuộc hội thoại, mức effort của phiên và agent mà cuộc hội thoại đang chạy dưới đó (xem [`--continue`](#continue)). Chính Claude Code khôi phục `--model` và `--name` từ cuộc hội thoại. Mọi thứ khác (`--worktree`, `--debug`, `--remote-control`, các tùy chọn ccam không biết) được nêu tên trong `ccam: not carried over: …`. Việc chuyển đổi không bao giờ hỏi câu hỏi về lịch sử dùng chung, và những gì bạn gõ trong lúc nó chạy sẽ bị bỏ; câu hỏi duy nhất của nó là câu hỏi trước khi chuyển sang một tài khoản của tổ chức khác.

Relay chỉ tiếp tục cuộc hội thoại khi Claude Code kết thúc theo cách `/ccam:switch` yêu cầu: không tiếp tục khi bạn đã gõ `/exit` trước, nhấn Ctrl+C trong lúc nó kết thúc, hoặc ccam nhận tín hiệu hangup hay một tín hiệu khác, và cũng không khi terminal đã mất. Khi đó ghi chú được giữ lại, và `ccam next --continue` sẽ tiếp tục nó sau. Một phiên được khởi động dưới relay của một phiên bản ccam cũ hơn, hoặc một phiên không có người theo dõi (các lần chạy `-p` và các phiên nền), không bị hook kết thúc; thay vào đó nó nhận thông báo `/exit`. `claude remote-control` chạy không có relay.

### `auto_switch`

Với `auto_switch = true` ở đầu `config.toml` (Linux và macOS), ccam tự chuyển một phiên sang tài khoản khác khi Claude Code dừng phiên đó ở giới hạn sử dụng: hook của plugin chọn một tài khoản khác còn hạn mức và kết thúc Claude Code, và ccam, đang chạy song song với phiên, tiếp tục cuộc hội thoại ở đó, như [relay](#relay) làm sau `/ccam:switch`. Tính năng này mặc định tắt và đang thử nghiệm.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam chỉ chuyển một phiên khi tất cả các điều kiện sau đều đúng:

1. ccam đã khởi động phiên đó trong một terminal và tiếp tục chạy song song với nó. `auto_switch` bật [relay](#relay), nên điều này đúng với các phiên tương tác bạn khởi động sau khi đặt nó. Một phiên khởi động trước đó, hoặc dưới một ccam cũ hơn, nhận một thông báo kết thúc bằng ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` Các phiên không có người theo dõi (các lần chạy `-p` và các phiên nền), `claude remote-control` và Claude Code khởi động không qua ccam không bao giờ bị chuyển.
2. Tài khoản đã bật [plugin ccam](/docs/install#claude-code-plugin) và chạy Claude Code 2.1.214 trở lên, phiên bản truyền ID tiến trình của nó cho hook (`CLAUDE_PID`).
3. Claude Code đã ghi yêu cầu bị từ chối vào cuộc hội thoại, kèm thời điểm đặt lại, cho giới hạn 5 giờ, hàng tuần, Opus hoặc Sonnet của tài khoản đó. Giới hạn chi tiêu, việc hết tín dụng sử dụng, việc tổ chức dừng tài khoản, hạn mức Fable (`You've reached your Fable limit.`), việc máy chủ điều tiết và giới hạn tốc độ của API key không bao giờ làm một phiên bị chuyển. Giới hạn của một subagent hoặc agent của workflow vẫn được tính, kể cả khi cuộc hội thoại chính chỉ đang chờ nó; khi đó là giới hạn hàng tuần của Opus hoặc Sonnet, phiên chỉ được chuyển nếu chính cuộc hội thoại chạy trên mô hình đó.
4. ccam biết gói của tài khoản, từ `.claude.json` của nó hoặc từ những gì `ccam login` đã ghi lại. Nếu không, thông báo kết thúc bằng `ccam does not know personal's plan, so auto_switch leaves this session.`
5. Phiên không được kết nối với Remote Control. Dưới một tài khoản khác, nó sẽ không kết nối lại với ứng dụng mà bạn dùng để điều khiển nó, nên nó ở lại, và thông báo nói rõ điều đó.
6. Giới hạn được đặt lại sau hơn 15 phút kể từ bây giờ. Với thời gian chờ ngắn hơn, ccam để phiên lại cho Claude Code: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. Bạn chưa gõ prompt nào kể từ khi chạm giới hạn.
8. Một tài khoản khác còn hạn mức. ccam chọn như [`--auto`](#auto) (`auto_max_usage`, `auto_strategy`, mô hình của cuộc hội thoại), trong số các tài khoản cùng loại với tài khoản bạn rời đi: một seat khác của cùng tổ chức Team hoặc Enterprise, hoặc một gói cá nhân khác (Pro, Max). `cross_org` và `--cross-org` không áp dụng ở đây. Tài khoản phải:
   - được đăng nhập bằng một đăng nhập mà Claude Code có thể gia hạn;
   - không được đăng nhập vào cùng tài khoản Claude;
   - không dùng API key;
   - có gói mà ccam biết;
   - không bật tín dụng sử dụng, theo những gì phép tra cứu mức sử dụng của chính ccam cho thấy;
   - đã tin cậy thư mục này trong Claude Code;
   - nếu cuộc hội thoại chạy dưới một agent (cài đặt `agent` hoặc cuộc hội thoại nêu tên một agent), dùng chung các agent của tài khoản bạn rời đi: các thư mục `agents` và `plugins` và tệp `settings.json` của tài khoản đó đều phải được dùng chung y hệt (qua một nhóm), hoặc không có ở cả hai.

   Khi ba điều kiện cuối loại bỏ mọi tài khoản còn hạn mức, thông báo nêu tên chúng: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` hoặc `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` Một khoảng có [bản ghi giới hạn](/features/usage#the-limit-record) được tính là đã đầy, và một tài khoản mà một phiên khác đã chuyển sang trong 2 phút gần nhất xếp sau các tài khoản khác. Lựa chọn phải được đưa ra trong khoảng 6 giây; khi các lượt tra cứu mức sử dụng lâu hơn, phiên ở lại. Với `usage_source = "claude"`, `auto_switch` quyết định chỉ từ các con số đã lưu, và ccam không biết được tín dụng sử dụng có bật hay không từ phép tra cứu của Claude Code hay từ dòng trạng thái, nên nó giả định chúng đang tắt.

Sau đó hook kết thúc Claude Code, và ccam, đang chạy song song với nó, hiện một thông báo trên màn hình ở các terminal hỗ trợ, chẳng hạn `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` Nó in dòng đầu tiên của đầu ra ví dụ ở trên và chạy `ccam next --continue` cho phiên đó dưới tài khoản đã chọn, lệnh này in hai dòng còn lại của nó và khởi động Claude Code với các tùy chọn, chế độ quyền, mức effort và agent của phiên, như [relay](#relay) làm. Tin nhắn đầu tiên của phiên mới là một văn bản cố định, thứ duy nhất ccam gửi thay cho bạn (một lần chuyển sau `/ccam:switch` không gửi gì):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

Sau giới hạn của một subagent thì là văn bản này, vì các subagent, workflow và lệnh chạy nền đã dừng cùng với phiên cũ:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

Nếu tài khoản đó cũng chạm giới hạn, điều tương tự xảy ra ở đó. Các bản ghi giới hạn ngăn chuỗi chuyển quay lại một tài khoản trước khi nó được đặt lại, nên khi mọi tài khoản đều đầy, phiên ở lại trên tài khoản cuối cùng.

Hook biết một cuộc hội thoại chạy dưới một agent nhờ cài đặt `agent`, hoặc nhờ dòng nêu tên agent mà Claude Code thỉnh thoảng ghi vào cuộc hội thoại và ghi khi phiên kết thúc. Nếu đến lúc đó chỉ có `--agent` của phiên nêu tên agent, và tài khoản được chọn không dùng chung các agent, `ccam next --continue` không khởi động gì và giữ lại ghi chú, và một thông báo trên màn hình thứ hai cho biết điều đó: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

Khi không tài khoản nào khác còn hạn mức, không có gì được chuyển: Claude Code xử lý giới hạn như thường lệ, chờ trong phiên hoặc đưa ra các lựa chọn của nó, và thông báo cho biết tài khoản nào được giải phóng trước, ví dụ `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` Với các seat, nó ghi `None of your other accounts in Acme has room now`. Để dừng một lần chuyển đã bắt đầu, hãy nhấn Ctrl+C trong lúc Claude Code đóng lại: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

Một lần chuyển có cái giá của nó. Claude Code khởi động lại, mất vài giây, và các shell nền, máy chủ phát triển và agent mà Claude đã khởi động sẽ dừng cùng phiên cũ. Tài khoản mới đọc lại toàn bộ cuộc hội thoại với câu trả lời đầu tiên, như với [`--continue`](#continue), và ccam cho biết nó lớn đến đâu; một lần chuyển không bao giờ dùng [`--brief`](#brief).

`auto_switch` không bao giờ đổi trước khi Claude Code báo một giới hạn, không bao giờ theo dõi mức sử dụng của bạn ở chế độ nền, và không bao giờ sao chép, gia hạn hay chuyển giao một đăng nhập. Trên Windows nó bị bỏ qua, và `ccam doctor` cho biết điều đó; ở đó, với `relay = true`, `/ccam:switch` và `/exit` vẫn chuyển một cuộc hội thoại sang tài khoản khác.

**Điều khoản của Anthropic.** Trên các seat Team và Enterprise, `auto_switch` chuyển công việc của bạn sang một seat khác của cùng tổ chức. Trung tâm trợ giúp của Anthropic ghi rằng "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (mỗi seat của gói Team được gán cho đích danh một cá nhân, và chỉ người đó được phép sử dụng) và "Users can't route their work through another user's account to use that user's remaining usage" (người dùng không được chuyển công việc của mình qua tài khoản của một người dùng khác để dùng phần hạn mức còn lại của người đó) ([gói Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); tương tự đối với [gói Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)).

Với các gói cá nhân (Pro, Max), điều khoản của Anthropic không nói liệu các tài khoản của chính bạn có được tự động tiếp quản cho nhau khi chạm giới hạn hay không. [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) cấm truy cập Dịch vụ "through automated or non-human means, whether through a bot, script, or otherwise" (bằng phương tiện tự động hoặc phi con người, dù qua bot, script hay cách khác) trừ khi Anthropic cho phép rõ ràng, và một lần chuyển gửi một trong các tin nhắn cố định ở trên thay cho bạn. [Usage Policy](https://www.anthropic.com/legal/aup) có hiệu lực từ ngày 12 tháng 11 năm 2026 liệt kê "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (cố ý vượt qua hoặc can thiệp vào các khả năng, hạn chế hoặc rào chắn được thiết lập trong sản phẩm của chúng tôi) trong số những điều nó cấm.

ccam không khẳng định rằng `auto_switch` tuân thủ các điều khoản này. Bật nó hay không là quyết định của bạn, và rủi ro cũng vậy: Anthropic có thể đình chỉ một tài khoản, hoặc quyền truy cập của một tổ chức Team hoặc Enterprise, khi họ cho rằng điều khoản bị vi phạm.

Khi `auto_switch` đang bật, `ccam doctor` cảnh báo về một tài khoản không có plugin ccam và về các phiên bản Claude Code cũ hơn 2.1.214, và với các seat thì nhắc lại quy tắc về seat của Anthropic ([các bước kiểm tra](/features/doctor)).

## Chạy lại khi chạm giới hạn

`ccam next` dành cho một phiên tương tác đã dừng. Một lần chạy không tương tác, `claude -p` trong script hoặc tác vụ định kỳ, không thể chờ bạn, nên `--fallback` đặt trước tên tài khoản cho phép ccam đổi tài khoản thay cho nó:

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

1. `--fallback` chỉ hoạt động cùng với `-p` hoặc `--print`; nếu không, ccam dừng lại và chỉ tới `ccam next --auto --continue`.
2. Đầu vào chuẩn từ pipe được đọc một lần, tối đa 10 MiB, và đưa cho mọi lần thử. Giống như Claude Code, ccam chờ tối đa 3 giây để nó bắt đầu (trừ khi dùng `--input-format stream-json`); nếu một pipe vẫn mở nhưng không gửi gì, lần chạy không nhận đầu vào và ccam ghi chú điều này trên đầu ra lỗi chuẩn.
3. ccam khởi động `claude` như một tiến trình con thay vì `exec`, giữ lại đầu ra chuẩn của nó, và chuyển tiếp đầu ra lỗi chuẩn ngay khi có.
4. Khi `claude` thoát với 0, hoặc thất bại vì lý do khác, đầu ra đã giữ và mã thoát được chuyển tiếp nguyên vẹn. Một lần chạy bị kết thúc bởi tín hiệu (mã thoát lớn hơn 128) cũng vậy: nó đã được yêu cầu dừng, và không được chạy lại.
5. Một lần thất bại được tính là do giới hạn sử dụng khi thông báo của Claude Code nói vậy, bất kỳ thông báo giới hạn sử dụng nào mà [hook](#inside-claude-code) nhận ra, kể cả bên trong đầu ra `--output-format json` hoặc `stream-json`, hoặc khi mức sử dụng của tài khoản cho thấy khoảng 5 giờ, hàng tuần hoặc của `--model` ở mức 100% (tính cả bản ghi giới hạn).
6. Khi đó đầu ra của lần thử đó bị bỏ đi và cùng lệnh chạy bằng một tài khoản đã đăng nhập khác: tài khoản có mức sử dụng 5 giờ thấp nhất, rồi đến các tài khoản không tra cứu được mức sử dụng, theo thứ tự số. Tài khoản ở mức 100% bị bỏ qua, tài khoản đăng nhập bằng API key không bao giờ được thử, và mỗi tài khoản được thử nhiều nhất một lần.
7. Khi có `-c`, `--continue`, `-r` hoặc `--resume` trong các đối số, chỉ các tài khoản trong cùng nhóm lịch sử được thử, để có thể tìm thấy cuộc hội thoại.
8. Chỉ các tài khoản thuộc cùng tổ chức được thử, như với `--auto` (`cross_org = true` bỏ giới hạn này).

Tài khoản tiếp theo được chọn theo cùng các quy tắc như `--auto` (`auto_max_usage`, `auto_strategy`, khoảng của `--model`). Khi không còn tài khoản nào, ccam in ra đầu ra của lần thử cuối, báo rằng không còn tài khoản đã đăng nhập nào còn hạn mức, và thoát với mã 75. Với `--wait` đặt trước tên tài khoản (`ccam --fallback --wait=2h work -p …`), thay vào đó nó chờ tài khoản đầu tiên được đặt lại, rồi có thể thử lại mọi tài khoản; nó luôn ngủ ít nhất một lần trước, nên tài khoản vừa chạm giới hạn sẽ không bị chạy lại ngay. Khoảng thời gian này giới hạn tổng thời gian chờ của lần chạy, không phải từng lần chờ. Viết với dấu cách (`--wait 30m`) thì bị từ chối kèm cùng gợi ý như `ccam next`. Dấu `--` đứng trước các đối số của Claude Code hoạt động như với mọi lần khởi động (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` in ra lệnh của lần thử đầu tiên.

Các phiên tương tác không tự đổi tài khoản trừ khi bạn bật [`auto_switch`](#auto-switch). Với `ccam statusline install`, dòng trạng thái thêm `→ ccam next --auto --continue` khi mức sử dụng 5 giờ hoặc hàng tuần chạm 90% ([Giới hạn sử dụng](/features/usage)).

## An toàn và các trường hợp đặc biệt

- Tài khoản chưa đăng nhập không bao giờ được chọn. Khi không có tài khoản nào khác đã đăng nhập, `ccam next` thông thường dừng lại với `no logged-in account other than <name>`.
- Với `--continue` và không có cuộc hội thoại nào trong thư mục: `no account has a conversation in <folder> to continue`.
- Khi quy tắc tổ chức không để lại tài khoản nào: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Với `--auto` và không có tài khoản nào đã đăng nhập: `no logged-in account to pick`. Khi mọi tài khoản đều đã đầy hoặc chưa rõ mức sử dụng: `` no account with known usage and headroom (see `ccam usage`) ``. Cả hai đều kết thúc với mã thoát 75.
- `--max-usage`, `--strategy`, `--wait` và `--cross-org` chỉ dùng được với `--auto`, và giá trị ngoài khoảng 1-100 hoặc một chiến lược không xác định sẽ bị từ chối.
- Một bản fork sao chép cuộc hội thoại, không phải trạng thái khác của tài khoản: các checkpoint trong `file-history/` của tài khoản gốc không đi theo, trừ khi lịch sử được dùng chung.
- `--auto` gửi yêu cầu mạng ngay cả với `--dry-run`. Nó chỉ đọc, và không bao giờ ghi vào tệp token.
- Tài khoản được chọn khởi động qua cùng đường đi với `ccam <name>`, nên các kiểm tra tương tự được áp dụng, bao gồm cảnh báo khi một phiên Claude Code khác dùng chung lịch sử đang chạy trong thư mục này.
- `last-used` chỉ được ghi khi Claude Code thực sự khởi động. `--dry-run` không động đến nó, nên chạy thử không đẩy vòng xoay tiến lên.
- `--fallback` giữ đầu ra chuẩn cho đến khi một lần thử kết thúc, nên đầu ra dạng luồng như `--output-format stream-json` đến một lần vào cuối. Một client truyền đầu vào `stream-json` theo luồng và chờ câu trả lời trước khi gửi thêm không thể dùng `--fallback`. Các tệp mà lần thử đầu tiên đã thay đổi trước khi chạm giới hạn vẫn giữ nguyên thay đổi; lần thử tiếp theo bắt đầu lại lời nhắc từ đầu.

## Liên quan

- [Giới hạn sử dụng](/features/usage)
- [Giới hạn sử dụng và chuyển đổi tài khoản](/docs/usage-limits)
- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Làm gì khi Claude Code chạm giới hạn sử dụng](/guides/usage-limit-workflow)
