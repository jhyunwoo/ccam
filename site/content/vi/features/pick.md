---
title: Chọn tài khoản cho script
description: ccam pick in tài khoản còn nhiều hạn mức nhất cho script, giao tài khoản khác nhau cho các worker khởi động cùng lúc, và thoát 75 khi hết hạn mức.
updated: 2026-11-05
order: 12
translation: machine
sourceHash: "0ec25752ace9"
---

`ccam pick` in tên của tài khoản đã đăng nhập còn nhiều hạn mức nhất, để một script có thể dùng nó khởi động Claude Code. Các worker khởi động cùng lúc sẽ mỗi worker nhận một tài khoản khác nhau, miễn là có đủ tài khoản còn hạn mức, nên một loạt tác vụ `claude -p` được trải ra trên các tài khoản của bạn thay vì dùng cạn một tài khoản.

**Seat của gói Team và Enterprise.** Việc tự động chọn tài khoản (`ccam next`, `--fallback`, `ccam pick`, phím `n` trên bảng điều khiển và, khi bạn bật nó, [`auto_switch`](/features/next#auto-switch)) cũng chuyển đổi giữa các seat của cùng một tổ chức Claude Team hoặc Enterprise, điều này có thể vi phạm điều khoản của Anthropic. Hãy đọc [Câu hỏi thường gặp](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) và xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó.

## Cách dùng

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

Mã thoát: `0` kèm một tên trên đầu ra chuẩn, `75` (`EX_TEMPFAIL`) khi không tài khoản nào còn hạn mức hoặc không có tài khoản nào đăng nhập, `1` cho các lỗi khác như một cờ sai.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--json` | In tài khoản, mức sử dụng 5 giờ và hàng tuần của nó, nguồn của các con số và lý do nó được chọn | tắt |
| `--model <model>` | Tính cả khoảng hàng tuần của mô hình đó, chẳng hạn `opus` | không có |
| `--max-usage <1-100>` | Ưu tiên các tài khoản có mức sử dụng dưới giá trị này | `auto_max_usage`, nếu không thì 100 |
| `--strategy headroom\|expiring\|balance` | Mức sử dụng hàng tuần thấp nhất trong số các tài khoản có mức sử dụng 5 giờ cao hơn mức thấp nhất không quá 3 điểm phần trăm, khoảng 5 giờ được đặt lại sớm nhất, hoặc nhiều hạn mức hàng tuần còn lại nhất cho mỗi ngày | `auto_strategy`, nếu không thì `headroom` |
| `--cross-org` | Chọn cả các tài khoản thuộc tổ chức khác | `cross_org`, nếu không thì tắt |
| `--wait[=<duration>]` | Khi không tài khoản nào còn hạn mức, chờ tài khoản đầu tiên được đặt lại, tối đa chừng này thời gian. Giá trị đặt sau `=`: `--wait 30m` bị từ chối kèm gợi ý `--wait takes its value after '=': --wait=30m` | tắt; `5h` khi chỉ ghi cờ |
| `--dry-run` | Chọn mà không ghi lại lần chọn | tắt |

## Cách hoạt động

1. ccam liệt kê các tài khoản đã đăng nhập (tài khoản API key và đăng nhập Console không có giới hạn nên bị loại ra) trong tổ chức của tài khoản mà thư mục này sẽ khởi động (xem [tổ chức](/features/next#auto)); một thư mục đăng nhập vào cùng tài khoản Claude với tài khoản đó cũng bị loại ra.
2. Nó lấy một khóa (`~/.local/state/ccam/picks.lock`), nên các lần chọn chạy cùng lúc sẽ diễn ra lần lượt.
3. Nó tra cứu mức sử dụng của mọi tài khoản cùng lúc, theo cùng cách với `ccam usage`, với bộ nhớ đệm 60 giây. Một khoảng mà Claude Code báo là đã đầy được tính là 100% cho đến khi được đặt lại ([bản ghi giới hạn](/features/usage#the-limit-record)), nên một tài khoản vừa chạm giới hạn bị loại ra kể cả khi các con số đã lưu của nó vẫn cho thấy còn hạn mức.
4. Các tài khoản được chọn trong 2 phút gần nhất, bởi `ccam pick` hoặc bởi [`auto_switch`](/features/next#auto-switch) khi nó chuyển một phiên sang đó (`~/.local/state/ccam/picks/<name>`; một lượt chọn có dấu thời gian ở tương lai không được tính) bị loại ra trước, và [quy tắc chọn](/features/next#auto) của `ccam next --auto` chọn trong số còn lại: ưu tiên dưới `auto_max_usage`, theo chiến lược, không bao giờ chọn tài khoản ở mức 100%. Nếu việc này bỏ qua tài khoản mà quy tắc lẽ ra sẽ chọn, `reason` của `--json` sẽ nói rõ: `30% of its 5 hours used; work was picked a moment ago`.
5. Nếu mọi tài khoản còn hạn mức đều vừa được chọn gần đây, các quy tắc sẽ chọn lại trong tất cả, nên với ba tài khoản, worker thứ tư sẽ dùng chung một trong số đó.
6. Lần chọn được ghi lại và khóa được nhả ra, rồi tên được in ra.

Khoảng 2 phút này quan trọng vì mức sử dụng của một tác vụ mới cần một lúc mới hiện ra trong các con số: nếu không có nó, năm worker khởi động cùng lúc đều sẽ thấy cùng một tài khoản là tài khoản còn nhiều hạn mức nhất.

## An toàn và các trường hợp đặc biệt

- `ccam pick` chỉ đọc mức sử dụng và ghi thời điểm chọn; nó không bao giờ gia hạn hay đọc token ngoài những gì `ccam usage` làm.
- Một lần chọn không giữ chỗ tài khoản. Sau 2 phút, hai script vẫn có thể khởi động trên cùng một tài khoản; hãy dùng `--fallback` để phục hồi khi tài khoản đó chạm giới hạn.
- Việc chờ với `--wait` không giữ khóa, nên các lần chọn khác vẫn diễn ra trong lúc đó.
- `ccam pick` bỏ qua [`ccam folder`](/features/folder): nó chỉ chọn theo mức sử dụng.

## Liên quan

- [Chuyển sang tài khoản tiếp theo](/features/next)
- [Chạy lại khi chạm giới hạn](/features/next#run-again-at-the-limit)
- [Giới hạn sử dụng](/features/usage)
- [Tài khoản API key](/features/api)
