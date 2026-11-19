---
title: Subagent trên các tài khoản khác
description: ccam agent chạy một phần tác vụ Claude Code dưới tài khoản khác của bạn và trả về câu trả lời của Claude, nên một tác vụ dài dùng hạn mức của nhiều tài khoản.
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Các subagent của chính Claude Code chạy dưới tài khoản của phiên khởi động chúng, nên chúng tiêu hạn mức của tài khoản đó. `ccam agent "<task>"` chạy một tác vụ bằng `claude -p` trong cùng thư mục dưới một tài khoản khác của bạn và in câu trả lời cuối cùng của Claude. Khi được gọi từ một phiên Claude Code, nó hoạt động như một subagent tiêu hạn mức của tài khoản khác, nên một tác vụ dài có thể tiếp tục sau thời điểm mà tài khoản của chính phiên đó lẽ ra đã hết hạn mức.

**Seat của gói Team và Enterprise.** `ccam agent` tự chọn tài khoản. Trong cùng một tổ chức Claude Team hoặc Enterprise, nó chạy tác vụ trên một seat khác, điều này có thể vi phạm điều khoản của Anthropic. Hãy đọc [Câu hỏi thường gặp](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) và xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó.

## Cách dùng

Claude Code chạy nó qua công cụ Bash của mình. Khi đã cài [plugin ccam](/docs/install#claude-code-plugin), hãy yêu cầu bằng lời thường, chẳng hạn "dùng các tài khoản khác của tôi cho các phần của tác vụ này", hoặc ghi điều đó vào CLAUDE.md của bạn. Skill của plugin bảo Claude Code khởi động mỗi agent ở chế độ nền, viết một tác vụ tự đứng được một mình, và giao cho các agent sửa tệp phần việc riêng hoặc git worktree riêng của chúng.

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

Câu trả lời cuối cùng của Claude đi ra đầu ra chuẩn; các dòng bắt đầu bằng `ccam:` đi ra đầu ra lỗi chuẩn. Thay vào đó, `--json` in `{account, session_id, result, is_error, accounts}`, trong đó `accounts` liệt kê mọi tài khoản mà tác vụ đã chạy dưới.

Mã thoát: mã của chính claude (`0` khi agent hoàn tất), `75` (`EX_TEMPFAIL`) khi không tài khoản nào ngoài tài khoản gọi còn hạn mức, `1` cho các lỗi khác như một agent cố khởi động agent khác.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--model <model>` | Mô hình của agent; khoảng hàng tuần của mô hình đó cũng được tính khi chọn | mô hình mặc định của tài khoản |
| `--resume <id>` | Tiếp tục cuộc hội thoại của một agent trước đó với một tác vụ tiếp theo (ID phiên trên dòng đầu và dòng cuối của nó) | tắt |
| `--json` | In tài khoản, phiên, câu trả lời và các tài khoản mà tác vụ đã chạy dưới | tắt |
| `--max-usage <1-100>` | Ưu tiên các tài khoản có mức sử dụng dưới giá trị này | `auto_max_usage`, nếu không thì 100 |
| `--strategy headroom\|expiring\|balance` | [Quy tắc chọn](/features/next#auto) | `auto_strategy`, nếu không thì `headroom` |
| `--cross-org` | Chạy cả dưới các tài khoản thuộc tổ chức khác | `cross_org`, nếu không thì tắt |
| `--wait[=<duration>]` | Khi không tài khoản nào còn hạn mức, chờ tài khoản đầu tiên được đặt lại, tối đa chừng này thời gian | tắt; `5h` khi chỉ ghi cờ |
| `--dry-run` | In lệnh cho tài khoản mà nó sẽ chọn, không chạy lệnh đó | tắt |

## Cách hoạt động

1. **Tài khoản gọi bị loại ra.** Tài khoản mà từ đó `ccam agent` được gọi (tài khoản của phiên do ccam khởi động, nếu không thì tài khoản mà thư mục này khởi động) không bao giờ chạy agent; một thư mục khác đăng nhập vào cùng tài khoản Claude, vốn dùng chung giới hạn với nó, cũng vậy. Agent ở lại trong [tổ chức](/features/next#auto) của tài khoản gọi, trừ khi `--cross-org` nói khác.
2. **Nó chọn như `ccam pick`.** Mức sử dụng được tra cứu theo cách [`ccam pick`](/features/pick) làm, và các agent khởi động cách nhau không quá 2 phút sẽ nhận các tài khoản khác nhau miễn là có đủ tài khoản còn hạn mức.
3. **Nó chạy với quyền của phiên, không bao giờ hơn.** Trong `claude -p` không có ai để cho phép một công cụ, nên agent nhận chế độ quyền và mức effort của phiên đã gọi nó. Khi chạy từ terminal, chế độ mặc định của Claude Code và các quy tắc quyền trong thiết lập của bạn được áp dụng. `ccam agent` không có tùy chọn nào cấp nhiều quyền hơn.
4. **Khi chạm giới hạn sử dụng, cuộc hội thoại đi tiếp.** Khi tài khoản của agent chạm giới hạn, ccam ghi lại giới hạn như [hook của plugin](/features/usage#the-limit-record) sẽ làm và tiếp tục cùng cuộc hội thoại dưới tài khoản còn nhiều hạn mức nhất: Claude Code sao chép nó sang đó (`--resume <file> --fork-session`), với tin nhắn đầu tiên "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." Mỗi tài khoản chỉ được thử một lần.
5. **Mọi cuộc hội thoại của agent đều được ghi lại** trong `~/.local/state/ccam/agent/<session id>.json` (30 ngày). `--resume <id>` tiếp tục nó trong thư mục của nó, dưới tài khoản đã chạy nó nếu tài khoản đó vẫn còn hạn mức, nếu không thì dưới tài khoản khác. [`ccam next --continue`](/features/next) bỏ qua các cuộc hội thoại của agent, nên nó vẫn tiếp tục cuộc hội thoại của bạn.

## An toàn và các trường hợp đặc biệt

- Một agent không thể khởi động agent của riêng nó: ccam đánh dấu nó bằng `CCAM_AGENT=1`, và `ccam agent` từ chối chạy ở đó.
- Dừng `ccam agent` (Ctrl+C, một tín hiệu, hoặc Claude Code ngắt lệnh) cũng dừng agent.
- Agent không thấy cuộc hội thoại đã gọi nó, chỉ thấy tác vụ của nó.
- Các agent cùng lúc sửa cùng những tệp sẽ ghi đè lên công việc của nhau. Hãy giao cho mỗi agent phần việc riêng của tác vụ, hoặc git worktree riêng.
- Khi phiên gọi chuyển sang tài khoản khác ([`auto_switch`](/features/next#auto-switch)), các agent mà nó đã khởi động cũng dừng theo. Các cuộc hội thoại của chúng vẫn còn: `--resume <id>` tiếp tục chúng.
- `agent` là tên một lệnh, nên không tài khoản mới nào có thể tên là `agent`; một tài khoản có sẵn mang tên đó vẫn khởi động bằng `ccam run agent`.

## Liên quan

- [Chọn tài khoản cho script](/features/pick)
- [Chuyển sang tài khoản tiếp theo](/features/next)
- [Giới hạn sử dụng](/features/usage)
- [Phiên chạy nền của mọi tài khoản](/features/agents)
