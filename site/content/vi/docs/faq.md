---
title: Câu hỏi thường gặp
metaTitle: "ccam FAQ: dùng nhiều tài khoản Claude Code"
description: "Giải đáp về ccam, tài khoản Claude Code, giới hạn sử dụng, việc dùng chung, quyền riêng tư và các hệ điều hành được hỗ trợ."
updated: 2026-11-17
order: 6
translation: machine
sourceHash: "9b13b37391f2"
---

## ccam có phải là công cụ chính thức của Anthropic không?

Không. ccam là dự án mã nguồn mở không chính thức, không liên kết với Anthropic và không được Anthropic xác nhận. Nó dùng thiết lập `CLAUDE_CONFIG_DIR` mà Claude Code [mô tả trong tài liệu](https://code.claude.com/docs/en/env-vars) để chạy song song nhiều tài khoản.

## Tôi có được phép dùng nhiều tài khoản Claude Code không?

Điều đó phụ thuộc vào điều khoản của từng tài khoản và, với tài khoản công việc, vào chính sách của tổ chức bạn: [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) của Anthropic cho Pro và Max, [Commercial Terms](https://www.anthropic.com/legal/commercial-terms) cho Team và Enterprise, và [Usage Policy](https://www.anthropic.com/legal/aup) cho tất cả. Hãy kiểm tra trước khi dùng nhiều tài khoản. ccam không nâng, gộp hay che giấu giới hạn của một tài khoản, và mỗi tài khoản vẫn dừng ở giới hạn của riêng nó. Nó chuyển công việc của bạn sang một tài khoản khác của bạn khi bạn yêu cầu, hoặc tự làm việc đó khi chạm giới hạn nếu bạn bật [`auto_switch`](/features/next#auto-switch) (Linux và macOS). Hãy kiểm tra xem việc đó có được phép với các tài khoản của bạn hay không trước khi dùng; quyết định là của bạn.

**Seat của gói Team và Enterprise.** Việc tự động chọn tài khoản (`ccam next`, `--fallback`, `ccam pick`, phím `n` trên bảng điều khiển và [`auto_switch`](/features/next#auto-switch) khi bạn bật nó) cũng chuyển đổi giữa các seat của cùng một tổ chức Claude Team hoặc Enterprise. Trung tâm trợ giúp của Anthropic ghi rằng "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (mỗi seat của gói Team được gán cho đích danh một cá nhân, và chỉ người đó được phép sử dụng) và "Users can't route their work through another user's account to use that user's remaining usage" (người dùng không được chuyển công việc của mình qua tài khoản của một người dùng khác để dùng phần hạn mức còn lại của người đó) ([gói Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); tương tự đối với [gói Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Vì vậy, dùng tính năng tự động chọn tài khoản giữa các seat có thể vi phạm điều khoản của Anthropic. Có dùng nó ở đó hay không là quyết định của bạn: hãy xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó. Các gói cá nhân (Pro, Max) được đề cập trong đoạn ở trên.

## ccam có giúp vượt qua giới hạn sử dụng không?

Không vượt qua giới hạn riêng của bất kỳ tài khoản nào: mỗi tài khoản giữ các giới hạn của nó và dừng lại ở đó, và ccam không nâng, gộp hay che giấu chúng. ccam hiển thị các giới hạn đó và cho phép bạn tiếp tục trên một tài khoản khác của bạn khi một tài khoản hết hạn mức; khi bật `auto_switch`, nó tự làm việc đó (xem câu hỏi tiếp theo). Giữa các seat của cùng một tổ chức Team hoặc Enterprise, tiếp tục trên một seat khác sẽ dùng phần hạn mức còn lại của seat đó, điều mà trung tâm trợ giúp của Anthropic không cho phép đối với seat của một người dùng khác (xem phần Seat của gói Team và Enterprise ở trên).

## ccam có thể tự đổi tài khoản khi một phiên chạm giới hạn không?

Chỉ khi bạn bật tính năng đó. Với `auto_switch = true` ở đầu `~/.config/ccam/config.toml` (Linux và macOS; mặc định tắt và đang thử nghiệm), khi Claude Code dừng một phiên tương tác (phiên do ccam khởi động) ở giới hạn 5 giờ, hàng tuần, Opus, Sonnet hoặc Fable, ccam kết thúc phiên và tiếp tục cuộc hội thoại dưới một tài khoản khác của bạn còn hạn mức: một seat khác của cùng tổ chức Team hoặc Enterprise, hoặc một gói cá nhân khác (Pro, Max). Nó không bao giờ đổi trước khi chạm giới hạn, không bao giờ theo dõi mức sử dụng của bạn ở chế độ nền, và không bao giờ chuyển công việc ra khỏi hay vào trong một tổ chức Team hoặc Enterprise.

Với các seat, nó chuyển công việc của bạn từ seat này sang seat khác, điều này có thể vi phạm điều khoản của Anthropic (xem phần Seat của gói Team và Enterprise ở trên). Với các gói cá nhân, điều khoản của Anthropic không đề cập rõ ràng đến việc này: Consumer Terms cấm truy cập Dịch vụ "through automated or non-human means" (bằng phương tiện tự động hoặc phi con người) trừ khi Anthropic cho phép rõ ràng, và một lần chuyển sẽ gửi một tin nhắn thay cho bạn. Bật nó hay không là quyết định của bạn.

[`auto_switch`](/features/next#auto-switch) liệt kê mọi điều kiện và chi phí của một lần chuyển.

## Làm sao để dùng hai tài khoản Claude Code trên một máy tính?

Hãy cho mỗi tài khoản một thư mục cấu hình riêng bằng `CLAUDE_CONFIG_DIR`. ccam làm việc này thay bạn: chạy `ccam add` cho từng tài khoản, rồi khởi động chúng bằng `ccam personal` và `ccam work`, hoặc bằng các lệnh tắt do bạn tự đặt tên, chẳng hạn `ccwork`. [Hướng dẫn](/guides/multiple-claude-code-accounts) cũng chỉ cách làm thủ công.

## Làm sao để tiếp tục một cuộc hội thoại trên tài khoản khác?

Chạy `ccam next --continue` trong thư mục: lệnh này chuyển cuộc hội thoại gần nhất của thư mục sang tài khoản kế tiếp (`--auto` để chọn tài khoản còn nhiều hạn mức nhất). Bên trong Claude Code, `/ccam:switch` từ [plugin ccam](/docs/install#claude-code-plugin) làm điều tương tự cho phiên bạn đang ở. Khi hai tài khoản dùng chung lịch sử hội thoại (`ccam group set default --add history`, hoặc một nhóm riêng của chúng: `ccam group add pair --items history --members work,side`), tài khoản mới mở chính cuộc hội thoại đó; nếu không, Claude Code sao chép nó vào lịch sử của tài khoản mới. Thêm `--brief` để tài khoản mới bắt đầu từ một bản tóm tắt ngắn thay vì toàn bộ cuộc hội thoại.

## ccam có xem hoặc gửi thông tin đăng nhập của tôi đi không?

ccam không bao giờ dùng chung hay sao chép thông tin đăng nhập Claude của bạn giữa các tài khoản, và không bao giờ sao chép các lần đăng nhập của máy chủ MCP. Nó chỉ đọc access token của một tài khoản để hỏi endpoint mức sử dụng của Anthropic về mức sử dụng của tài khoản đó: cho `ccam usage`, bảng điều khiển, `ccam doctor` và lựa chọn tự động (`ccam next --auto`, `ccam pick`, `--fallback`, và `auto_switch` khi một phiên chạm giới hạn). Nó chỉ gia hạn token đã hết hạn bằng `ccam usage --refresh`; với `usage_source = "claude"` nó không gửi token đi đâu cả và chỉ đọc thông tin đăng nhập đã lưu (tệp, hoặc mục trong Chuỗi khóa của macOS) để biết nó còn dùng được hay không. `ccam ext` và `ccam doctor` đọc tên các lần đăng nhập MCP của một tài khoản, không bao giờ đọc token của chúng. ccam không có telemetry. Ngoài các lượt tra cứu đó, yêu cầu duy nhất nó tự gửi là lượt kiểm tra bản phát hành ccam mới hơn trên GitHub của bảng điều khiển, và yêu cầu này không mang thông tin gì về bạn hay các tài khoản của bạn.

## Làm sao để đăng nhập trên máy chủ từ xa?

Chạy `ccam add` hoặc `ccam login <name>` trong phiên SSH như bình thường. Claude Code chờ việc đăng nhập trên máy chủ, nên trình duyệt trên máy tính của bạn dừng ở một trang `http://localhost:<port>/callback?code=…` không tải được. Hãy sao chép địa chỉ đó từ trình duyệt và dán vào terminal: ccam chuyển nó cho Claude Code trên máy chủ, và Claude Code hoàn tất việc đăng nhập. Địa chỉ này chứa một mã dùng một lần, vì vậy đừng chia sẻ nó.

## Các tài khoản của tôi được lưu ở đâu?

Trong `~/.local/share/ccam/profiles/<name>/`, mỗi tài khoản một thư mục. Trên macOS, Claude Code giữ chính thông tin đăng nhập trong Chuỗi khóa (Keychain), dưới một tên được tạo từ thư mục đó, nên `ccam rename` giữ lại tên cũ của thư mục ở đó.

## Thư mục ~/.claude hiện có của tôi sẽ ra sao?

Không có gì thay đổi, trừ khi bạn chạy `ccam import` hoặc `ccam migrate`. Lệnh import đăng ký thư mục ngay tại chỗ, nên lệnh `claude` thông thường và `ccam <name>` tiếp tục dùng cùng một tài khoản và cùng cài đặt. `ccam migrate` sao chép cài đặt, `CLAUDE.md`, tiện ích mở rộng và máy chủ MCP của thư mục đó vào nhóm `default`, và chỉ đọc thư mục.

## ccam hỗ trợ những hệ điều hành nào?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine và các bản phân phối khác), macOS và Windows 10 và 11, chạy trực tiếp hoặc bên trong WSL. Trên Windows, việc dùng chung các tệp như `settings.json` cần bật Developer Mode; thư mục dùng chung thì hoạt động mà không cần chế độ này, và khi không có nó, lần `ccam add` đầu tiên tạo nhóm `default` chỉ với máy chủ MCP và tiện ích mở rộng ([dùng chung trên Windows](/features/share#on-windows)).

## Vì sao một máy chủ MCP chỉ xuất hiện ở một tài khoản?

Máy chủ được thêm trong một phiên bằng `claude mcp add` chỉ được lưu cho tài khoản đó, và nếu không có `-s user` thì chỉ cho thư mục hiện tại. ccam sẽ báo cho bạn vào lần tới bạn khởi động tài khoản, và `ccam mcp list` cùng `ccam ext` hiển thị chúng; hãy chuyển máy chủ sang nhóm dùng chung bằng `ccam mcp import <account> <server>`.

## ccam có làm Claude Code chậm đi không?

Khởi động Claude Code thông qua ccam tốn thêm khoảng 4 mili giây, đo bằng `scripts/perf-check.sh` trên một máy tính xách tay. ccam đọc một tệp nhỏ, kiểm tra các liên kết dùng chung rồi tự thay thế mình bằng `claude`, nên theo mặc định nó không chạy trong lúc bạn làm việc. Windows không thể thay thế một tiến trình, nên ở đó ccam khởi động Claude Code và chờ nó, ngoài ra không làm gì khác. Khi bật `relay` hoặc `auto_switch`, ccam cũng tiếp tục chạy theo cách đó song song với mỗi phiên tương tác, để có thể tiếp tục cuộc hội thoại trên một tài khoản khác.

## Làm sao để cập nhật ccam?

Chạy `ccam update`. Lệnh này tải về bản phát hành mới nhất, kiểm tra mã băm rồi thay thế tệp `ccam`; `ccam update --check` chỉ cho biết có bản phát hành mới hơn hay không. Bảng điều khiển cũng báo cho bạn khi có bản phát hành mới hơn, và `U` sẽ cài nó sau khi hỏi; ccam không bao giờ tự cập nhật. Nếu bạn đã cài gói `.deb`, hãy dùng apt để cài gói `.deb` mới. Xem [cập nhật](/docs/install#update).

## Làm sao để gỡ cài đặt ccam?

Chạy `ccam uninstall`. Lệnh này liệt kê những gì sẽ bị xóa và hỏi trước, sau đó xóa mọi tài khoản (đăng xuất trước), các lệnh tắt, `~/.config/ccam`, `~/.local/share/ccam`, `~/.local/state/ccam` và lệnh `ccam`. `~/.claude` đã nhập vẫn được giữ lại, kèm một bản sao của những gì nó đã dùng chung trong các nhóm. Xem [cài đặt](/docs/install#uninstall).
