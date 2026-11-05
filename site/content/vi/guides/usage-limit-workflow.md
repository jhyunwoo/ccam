---
title: Làm gì khi Claude Code chạm giới hạn sử dụng
description: "Chờ giới hạn được đặt lại hoặc tiếp tục ngay cùng cuộc hội thoại trên một tài khoản Claude Code khác. Giới hạn 5 giờ và hàng tuần, từng bước một."
updated: 2026-11-05
order: 3
translation: machine
sourceHash: "53b1e3634ea9"
---

Khi một tài khoản chạm giới hạn sử dụng (thường gọi là rate limit), bạn có hai lựa chọn: chờ đến khi giới hạn được đặt lại, hoặc tiếp tục ngay cùng cuộc hội thoại trên một tài khoản khác của bạn. Với ccam, cách thứ hai chỉ cần một lệnh, `ccam next --auto --continue`.

## Claude Code có những giới hạn nào?

Với gói Pro và Max, Claude Code có [giới hạn phiên 5 giờ và giới hạn hàng tuần](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Cả hai đều dùng chung với Claude trên web và trong ứng dụng, và mức sử dụng được tính vào cả hai cùng lúc. Khi một trong hai hết, Claude Code [hiển thị thông báo](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) như "You've hit your session limit · resets 3:45pm". Đổi model không giúp được gì, vì giới hạn áp dụng cho mọi model.

## Cách 1: chờ giới hạn được đặt lại

Nếu bạn chỉ có một tài khoản, hãy chờ. Các phiên bản Claude Code gần đây [chờ ngay trong phiên đang mở](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) và tự tiếp tục công việc khi giới hạn được đặt lại. `/usage` trong Claude Code cho biết bạn đã dùng bao nhiêu trong giới hạn của gói.

## Cách 2: tiếp tục ngay trên tài khoản khác

Nếu có tài khoản thứ hai, bạn có thể tiếp tục cùng cuộc hội thoại ngay lập tức. Thiết lập một lần:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Giờ đây dòng trạng thái hiển thị tài khoản và mức sử dụng 5 giờ của nó trong lúc bạn làm việc, nên việc chạm giới hạn sẽ không còn bất ngờ.

Khi chạm giới hạn, hãy thoát Claude Code (`/exit`), rồi chạy:

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

Mở cuộc hội thoại trên tài khoản còn nhiều hạn mức nhất:

```sh
ccam next --auto --continue
```

Claude Code khởi động với tài khoản `work` và cuộc hội thoại gần nhất của thư mục hiện tại, bao gồm mọi nội dung đã trao đổi. Dùng chung lịch sử (`ccam group set default --add history` ở trên) nghĩa là cả hai tài khoản làm việc trên cùng một cuộc hội thoại; nếu không, Claude Code sao chép cuộc hội thoại vào `work`, và ccam cho bạn biết câu trả lời đầu tiên phải đọc lại bao nhiêu token.

Với [plugin ccam](/docs/install#claude-code-plugin) trong Claude Code, một phiên dừng lại vì chạm giới hạn sẽ hiện một thông báo và ghi lại khoảng đã đầy, nên `ccam next --auto` không chọn lại tài khoản đó trước khi đặt lại. `/ccam:switch` (hoặc `/ccam:switch work`) ghi chú cuộc hội thoại cho `ccam next --continue`, để lệnh này tiếp tục đúng cuộc hội thoại đó.

**Seat của gói Team và Enterprise.** Việc tự động chọn tài khoản (`ccam next`, `--fallback`, `ccam pick`, phím `n` trên bảng điều khiển và, khi bạn bật nó, [`auto_switch`](/features/next#auto-switch)) cũng chuyển đổi giữa các seat của cùng một tổ chức Claude Team hoặc Enterprise, điều này có thể vi phạm điều khoản của Anthropic. Hãy đọc [Câu hỏi thường gặp](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) và xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó.

## ccam có tự động đổi tài khoản không?

Chỉ khi bạn bật tính năng đó.

- Theo mặc định, ccam không chạy trong lúc bạn làm việc: nó khởi động Claude Code rồi tự thay thế mình bằng Claude Code, nên bạn quyết định khi nào đổi. `ccam next --auto` chọn theo mức sử dụng, lấy tài khoản có mức sử dụng hàng tuần thấp nhất trong số các tài khoản có mức sử dụng 5 giờ cao hơn mức thấp nhất không quá 3 điểm phần trăm; nó bỏ qua các tài khoản đã đạt 100% hoặc chưa đăng nhập, và chỉ ở trong một tổ chức Team hoặc Enterprise, hoặc giữa các gói cá nhân.
- Với `relay = true` trong `config.toml`, thay vào đó ccam tiếp tục chạy song song với phiên, nên `/ccam:switch` tự tiếp tục cuộc hội thoại, giữ lại các tùy chọn, chế độ quyền, mức effort và agent của phiên; trước khi chuyển cuộc hội thoại sang một tài khoản của tổ chức khác, ccam sẽ hỏi bạn.
- Với `auto_switch = true` (Linux và macOS; thử nghiệm), ccam còn tự thực hiện việc chuyển: khi Claude Code dừng phiên ở giới hạn của nó và một tài khoản khác cùng loại còn hạn mức, ccam kết thúc phiên và tiếp tục cuộc hội thoại ở đó ([điều kiện và chi phí](/features/next#auto-switch)).

Một số công cụ khác thì hoán đổi đăng nhập của phiên đang chạy; [ccam và các công cụ đổi tài khoản khác](/compare/account-switchers) giải thích sự khác biệt.

## Những điều cần lưu ý

- Chạy `ccam next --continue` trong đúng thư mục bạn đã làm việc; hội thoại được lưu theo từng thư mục dự án.
- Tránh chạy hai tài khoản trên cùng một cuộc hội thoại cùng lúc. ccam sẽ cảnh báo nếu bạn thử.
- Không dùng chung lịch sử, cuộc hội thoại được sao chép sang tài khoản kia và hai bản sao đi theo hướng riêng; các checkpoint mà `/rewind` dùng để khôi phục mã vẫn ở lại với tài khoản đầu tiên.
- Tin nhắn đầu tiên trên tài khoản mới có thể tốn nhiều hạn mức hơn bình thường, vì tài khoản đó chưa lưu cache cuộc hội thoại. Thay vào đó, `ccam next --continue --brief` bắt đầu nó từ một bản tóm tắt ngắn.
- Việc bạn có được đổi tài khoản như vậy hay không phụ thuộc vào điều khoản áp dụng cho từng tài khoản.

Xem [Giới hạn sử dụng](/features/usage) và [Chuyển sang tài khoản tiếp theo](/features/next) để biết mọi tùy chọn.
