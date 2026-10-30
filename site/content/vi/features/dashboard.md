---
title: Bảng điều khiển
description: "ccam trên terminal mở bảng điều khiển toàn màn hình gồm tài khoản, nhóm, quy tắc thư mục, tiện ích mở rộng và kiểm tra; mọi thay đổi đều hiện kế hoạch trước."
updated: 2026-10-30
order: 15
translation: machine
sourceHash: d65f1b5bb0c8
---

`ccam` không kèm đối số sẽ mở một bảng điều khiển toàn màn hình khi đầu vào và đầu ra chuẩn là terminal; nếu không, nó in ra `ccam list`. Bảng điều khiển hiển thị mọi tài khoản cùng mức sử dụng và các nhóm của chúng, và làm được những gì các lệnh làm: khởi động, thêm, đăng nhập, thay đổi nhóm, đặt quy tắc thư mục, liệt kê các phiên nền, kiểm tra và sửa chữa. Khởi động một tài khoản, hoặc gắn vào một phiên nền, sẽ đóng bảng điều khiển; thêm tài khoản hoặc đăng nhập thì tạm nhường terminal trong chốc lát, rồi quay lại bảng điều khiển.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · r to refresh · a to add · ? for shortcuts
```

Dòng dưới tên các thẻ cho biết trong một câu thẻ hiện tại dùng để làm gì. Trên terminal rộng hơn, thẻ Accounts còn hiển thị email của từng tài khoản. Trên mọi thẻ, cuối dòng chứa các thẻ cho biết [`auto_switch`](/features/next#auto-switch) có đang bật hay không (Linux và macOS; trên Windows nó bị bỏ qua). Khi nó bật, thẻ Health liệt kê những gì tính năng này còn cần: Claude Code 2.1.214 trở lên, và plugin ccam được bật trong mọi tài khoản đăng nhập bằng gói đăng ký Claude.

Mỗi lần mở, bảng điều khiển hỏi GitHub ở chế độ nền xem có bản phát hành ccam mới hơn không; khung hình đầu tiên không bao giờ chờ việc này, và khi không có mạng thì không hiển thị gì. Khi có bản mới, một dòng dưới tiêu đề sẽ đề xuất nó:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` hỏi trước, rồi cập nhật theo cách của [`ccam update`](/features/update): bản phát hành được kiểm tra với checksum SHA-256 của nó và thay thế ccam trong một bước. Bảng điều khiển bạn đang dùng vẫn là phiên bản cũ, nên hãy thoát rồi chạy lại `ccam`. Khi không thể cập nhật tại đây, ví dụ với một bản ccam được cài từ gói `.deb`, bảng điều khiển hiển thị cùng lời khuyên như `ccam update` và vẫn giữ lời đề nghị.

Khung hình đầu tiên được vẽ chỉ từ `config.toml`, trong chưa đến 100 ms; mức sử dụng được điền vào khi từng tài khoản phản hồi. Mỗi giới hạn có thanh riêng: khoảng 5 giờ, khoảng hàng tuần cho mọi mô hình, và khoảng hàng tuần riêng của Fable; `-` nghĩa là tài khoản không có khoảng đó. Một tỷ lệ phần trăm chuyển sang màu vàng ở 80% và màu đỏ ở 100%. Nếu một lần tra cứu thất bại, các giá trị từ bộ nhớ đệm được hiển thị mờ. Một khoảng mà Claude Code báo là đã đầy hiển thị 100% cho đến khi được đặt lại. Một access token đã hết hạn hiển thị `⚠ token expired` mờ (Claude Code gia hạn nó ở lần khởi động tiếp theo), còn `⚠ login expired` nghĩa là tài khoản cần `ccam login`.

## Các thẻ

| Thẻ | Nội dung hiển thị | Phím |
| --- | --- | --- |
| Accounts | số, tên, email, thanh 5 giờ, hàng tuần và Fable, lần dùng gần nhất, nhóm; `✔` tài khoản mặc định, `▸` tài khoản của thư mục này, `•` đang chạy (kể cả Claude Code được khởi động không qua ccam: `claude` thông thường, một IDE, các phiên nền và daemon của chúng) | bên dưới |
| Agents | những gì [`ccam agents`](/features/agents) liệt kê: các phiên nền và các phiên đang mở của mọi tài khoản, cùng trạng thái, thư mục và thời điểm bắt đầu của chúng (bên dưới) | `enter` gắn vào (như `ccam attach`), `x` dừng, `r` tải lại |
| Groups | từng nhóm, những gì nhóm dùng chung và các tài khoản của nhóm, cùng các tài khoản không thuộc nhóm nào | `enter` sửa, `n` tạo mới, `e` đổi tên, `x` xóa |
| Folders | các quy tắc [`ccam folder`](/features/folder) và quy tắc áp dụng ở đây, kể cả khi bạn đến thư mục qua một liên kết tượng trưng (quy tắc mới được lưu với các liên kết đã được phân giải) | `n` dùng một tài khoản ở đây, `enter` thay đổi, `x` xóa |
| Extensions | những gì [`ccam ext`](/features/ext) hiển thị, mỗi tiện ích mở rộng một lần cùng các nhóm mà nó đến từ, `off:` các tài khoản mà một plugin chưa được bật, các máy chủ riêng của tài khoản chỉ áp dụng cho thư mục này, `This folder's project`, và các vấn đề không thuộc tài khoản nào | `r` tải lại |
| Health | những gì [`ccam doctor`](/features/doctor) báo cáo, dưới dạng một hộp cho mỗi tài khoản (bên dưới) | `f` sửa chữa hộp đang chọn, `F` sửa chữa tất cả (`doctor --fix`), `↑`/`↓` di chuyển giữa các hộp |

`←`/`→` hoặc `Tab` chuyển thẻ, `?` liệt kê các phím của thẻ hiện tại, `r` làm mới mức sử dụng, bỏ qua bộ nhớ đệm ngắn hạn, và tải lại các thẻ Agents, Extensions và Health, và `q` hoặc `Esc` để thoát.

## Phím trên thẻ Accounts

| Phím | Tác dụng |
| --- | --- |
| `Enter` | Khởi động Claude Code bằng tài khoản đang chọn |
| `1`–`9` | Di chuyển tới tài khoản có số đó; gõ nhanh `1` `2` để tới 12. Chỉ `Enter` mới khởi động |
| `↑`/`↓`, `k`/`j` | Di chuyển |
| `n` | Khởi động tài khoản mà [`ccam next --auto`](/features/next) sẽ chọn |
| `o` | Khởi động với tùy chọn: `--continue`, `--resume`, các đối số khác, hoặc chạy thử |
| `u` | Chi tiết mức sử dụng: khoảng 5 giờ, hàng tuần và theo mô hình, cùng thời điểm đặt lại và nguồn của các con số |
| `a` | Thêm một tài khoản (trình hướng dẫn bên dưới) |
| `l` | Đăng nhập lại |
| `s` | Chọn các nhóm của tài khoản |
| `d` | Đặt nó làm tài khoản mặc định |
| `f` | Dùng nó trong thư mục hiện tại |
| `e` | Đổi tên nó |
| `c` | Đặt hoặc xóa lệnh tắt của nó |
| `x` | Xóa nó (hỏi trước, với No được chọn sẵn) |

**Seat của gói Team và Enterprise.** Việc tự động chọn tài khoản (`ccam next`, `--fallback`, `ccam pick`, phím `n` trên bảng điều khiển và, khi bạn bật nó, [`auto_switch`](/features/next#auto-switch)) cũng chuyển đổi giữa các seat của cùng một tổ chức Claude Team hoặc Enterprise, điều này có thể vi phạm điều khoản của Anthropic. Hãy đọc [Câu hỏi thường gặp](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) và xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó.

## Thẻ Health

Thẻ Health chạy các kiểm tra của [`ccam doctor`](/features/doctor) khi bạn mở nó và vẽ cùng một giao diện: dòng đầu có một ký hiệu cho mỗi tài khoản, rồi đến một hộp cho mỗi tài khoản. Mỗi hàng trong hộp là một mục đã được kiểm tra, `Folder`, `Login` và `Share links` đứng đầu, được đánh dấu `✔` ổn, `ℹ` thông tin, `⚠` cảnh báo hoặc `✘` lỗi, nên bạn thấy ngay tài khoản nào có vấn đề và ở đâu. Các kết quả không thuộc tài khoản nào nằm trong hộp `General`.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` di chuyển từ hộp này sang hộp khác, và hộp đang chọn có viền màu nhấn. `f` chỉ sửa chữa hộp đang chọn, theo cách `ccam doctor --fix` sẽ làm: trên hộp của một tài khoản, nó chỉ sửa chữa tài khoản đó, còn trên hộp `General` là các kết quả không thuộc tài khoản nào (đường dẫn plugin và lệnh tắt); mọi tài khoản khác chỉ được kiểm tra lại. Nó hỏi trước (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), và một hộp không có gì để sửa chữa chỉ hiển thị `nothing to repair in <name>`. `F` sửa chữa mọi thứ cùng lúc, như `ccam doctor --fix`, và `r` kiểm tra lại.

## Thẻ Agents

Thẻ Agents hiển thị những gì [`ccam agents`](/features/agents) in ra: các phiên nền và các phiên đang mở của mọi tài khoản, cùng tài khoản mà mỗi phiên thuộc về. Nó hỏi Claude Code của mọi tài khoản cùng lúc (`claude agents --json`, mỗi tài khoản với thời gian chờ tối đa 15 giây) khi bạn mở thẻ lần đầu, khi nhấn `r`, sau khi dừng một phiên và sau khi cấu hình thay đổi, và hiển thị `Asking every account's Claude Code for its sessions…` cho đến khi tất cả đã trả lời.

```console
  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Background agents and open sessions of every account; enter attaches to one.

    Account   ID        Name   State          Folder      Started
 ❯  work      -         -      open (busy)    ~/code/web  now
    personal  49ec0028  probe  blocked, idle  ~/code/app  12m

 ✘ lab: could not list its sessions: exit status 1 'claude agents --json' is
   not available in this environment.

 enter to attach · x to stop · r to reload · ? for shortcuts
```

`State` là trạng thái của một phiên nền và việc tiến trình của nó đang làm, hoặc `open` đối với một phiên tương tác, như trong `ccam agents`. Nếu không truy vấn được một tài khoản, một dòng `✘` dưới danh sách sẽ hiển thị những gì Claude Code đã nói, như `lab` ở trên, tài khoản đã tắt chế độ xem agent. Các phiên nền đã kết thúc bị bỏ ra (`ccam agents --all` liệt kê chúng), và khi mọi tài khoản đã trả lời mà không có phiên nào, thẻ hiển thị ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` Trên terminal hẹp, thẻ bỏ `Started`, rồi `Name`, rồi `Folder`, và cắt bớt các đường dẫn thư mục dài từ bên trái.

`Enter` trên một phiên nền sẽ đóng bảng điều khiển và mở phiên đó trong terminal này dưới tài khoản của nó, như `ccam attach <id>`; bảng điều khiển không quay lại (trong Claude Code, `←` đi tới chế độ xem agent của nó và `Ctrl+Z` về shell). Một phiên đang mở không có ID, nên `Enter` và `x` chỉ báo rằng cần một phiên nền. `x` hỏi trước, với `No, keep it` được chọn sẵn: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` Nếu chọn có, ccam chạy `claude stop <id>` dưới tài khoản đó, hiển thị ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` hoặc lý do không dừng được, rồi đọc lại danh sách.

## Thay đổi nhóm

Trên thẻ Groups, `Enter` mở trình chỉnh sửa: các mục ở bên trái, các tài khoản ở bên phải, `Tab` để chuyển qua lại, `Space` để bật/tắt. Một mục mà một trong các tài khoản đã nhận từ một nhóm khác sẽ bị làm mờ kèm lý do, và một tài khoản sẽ nhận một mục hai lần cũng vậy; một mục chỉ có thể đến từ một nhóm. `Enter` hiển thị kế hoạch:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Nếu một tệp khác với bản của nhóm, ccam hỏi trước cách xử lý (giữ bản của nhóm, dùng bản của tôi hoặc gộp). Sau đó thay đổi chạy ở chế độ nền với cùng mã như [`ccam group`](/features/share): cấu hình được lưu trước, dữ liệu của từng tài khoản được di chuyển theo thứ tự, và bản gốc được đưa vào `.ccam-backup/`; một tài khoản rời việc dùng chung MCP vẫn giữ các máy chủ của nhóm trong `.claude.json` của nó.

Một thay đổi sẽ di chuyển dữ liệu của một tài khoản đang chạy bị từ chối kèm tên tài khoản đó, kể cả khi tài khoản được khởi động trong lúc kế hoạch đang mở, và một kế hoạch đã lỗi thời vì một tiến trình ccam khác đã thay đổi các nhóm sẽ bị từ chối với `the groups changed while this plan waited; nothing was moved, make the change again`. Khi một tiến trình ccam khác đang thay đổi các nhóm, một lần thay đổi, đổi tên hoặc xóa mới sẽ báo `another ccam process is changing the groups right now`; sau một lần thay đổi bị dừng giữa chừng, chúng chờ cho đến khi `F` trên thẻ Health (hoặc `ccam doctor --fix`) hoàn tất nó.

## Thêm tài khoản

`a` mở trình hướng dẫn mà `ccam add` cũng dùng trên terminal:

1. Tên tài khoản.
2. Cách nó đăng nhập: một gói đăng ký Claude ngay bây giờ (trong trình duyệt), một [API key](/features/api), hoặc một gói đăng ký Claude sau.
3. Với API key: nơi ccam đọc khóa mỗi lần khởi động (`env:NAME` hoặc `command:COMMAND`), endpoint, cách gửi khóa và mô hình.
4. Các nhóm để tham gia, với `default` đã được đánh dấu (nếu chưa có nhóm nào, đó là những gì nhóm `default` sẽ dùng chung).
5. Một lệnh tắt, nếu bạn muốn.
6. Một bản tóm tắt để xác nhận.

`Esc` quay lại một bước, và ở bước đầu tiên thì đóng trình hướng dẫn.

## Xóa tài khoản

`x` hỏi trước khi có bất cứ điều gì xảy ra, với No được chọn sẵn, và cho biết việc xóa sẽ làm gì: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` Với tài khoản API key thì không có thông tin đăng nhập nào để thu hồi, và khi một tài khoản khác giữ cùng thông tin đăng nhập thì nó không bị đăng xuất. Một tài khoản được nhập giữ lại thư mục của nó: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## Màu sắc

Màu sắc theo các vai trò trong giao diện terminal của Claude Code, nên ccam trông hài hòa khi đặt cạnh nó; tên nhóm nhận một màu dựa trên chính tên đó. ccam hỏi terminal về màu nền và chọn chủ đề sáng hoặc tối. Bạn có thể tự đặt một chủ đề ở đầu `config.toml`, phía trên `[[group]]`, `[[profile]]` hoặc `[[folder]]` đầu tiên:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Các chủ đề `-daltonized` dành cho người mù màu đỏ-lục, còn các chủ đề `-ansi` dùng 16 màu của chính terminal. `NO_COLOR` tắt màu.

## An toàn và các trường hợp đặc biệt

- Bảng điều khiển gọi cùng các hàm như các lệnh, nên mọi quy tắc đều được giữ: kế hoạch trước khi dữ liệu được di chuyển, từ chối với tài khoản đang chạy, bản sao lưu, và token đăng nhập không bao giờ rời khỏi tài khoản.
- Bảng điều khiển không có chế độ chạy thử: `ccam --dry-run` không kèm tên sẽ in ra cách tài khoản của thư mục này sẽ khởi động, và `o` trên một tài khoản cho phép chạy thử riêng lần khởi động đó.
- Terminal được khôi phục trước khi Claude Code khởi động hoặc việc đăng nhập mở ra.
- Một màn hình quá nhỏ cho hộp thoại sẽ bỏ bớt các hàng của thẻ trước; các dòng bị cắt bằng `…` thay vì xuống dòng. Trên màn hình thấp, các dòng trống và phần mô tả thẻ bị bỏ tiếp theo, để dòng gợi ý phím vẫn hiển thị. Dưới 40 cột hoặc 8 hàng, bảng điều khiển chỉ hiển thị `Too small: ccam needs 40x8 (now 36x7)`.
- `Ctrl+C` trong lúc một thay đổi, một lần xóa, một lần cập nhật hoặc một lần sửa chữa đang được áp dụng sẽ chờ nó xong: dòng trạng thái báo `ccam quits when it is done.`
- Sau khi `U` đã cập nhật ccam, lời đề nghị không xuất hiện lại khi bảng điều khiển quay về sau khi thêm tài khoản hoặc đăng nhập.

## Liên quan

- [Dùng chung bằng nhóm](/features/share)
- [Khởi động Claude Code](/features/launch)
- [Quản lý tài khoản](/features/accounts)
- [Giới hạn sử dụng](/features/usage)
