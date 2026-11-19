---
title: Quản lý tài khoản
description: "Thêm, đăng nhập, liệt kê, đổi tên, xóa và nhập tài khoản Claude Code bằng ccam add, ccam login, ccam list, ccam rename, ccam remove và ccam import."
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

Mỗi tài khoản ccam là một thư mục cấu hình Claude Code riêng với thông tin đăng nhập riêng. Các lệnh này tạo, đăng nhập, liệt kê, đổi tên và xóa tài khoản, cũng như đăng ký các thư mục bạn đã có.

## Cách dùng

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account that was never logged in or whose login expired
ccam logout side                           # log out: the login and its MCP sign-ins
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, logs out, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Ở bất kỳ lệnh nào nhận `<name|number>`, bạn có thể đưa tên tài khoản hoặc số slot của nó (cột `#` trong `ccam list`).

`ccam use` chỉ thay đổi tài khoản nào là mặc định (`default` trong `config.toml`).

`ccam add` chạy trên terminal mà không có tùy chọn nào sẽ mở một trình hướng dẫn. Nó đi qua các bước sau:

1. tên tài khoản, trừ khi bạn đã đưa tên;
2. cách tài khoản đăng nhập: đăng nhập Claude ngay bây giờ, một API key hoặc đăng nhập Claude sau;
3. các nhóm để tham gia, với `default` đã được đánh dấu (nếu chưa có nhóm nào, nó hiển thị những gì nhóm `default` sẽ dùng chung);
4. `Create a shortcut command for work? [y/N]`;
5. một bản tóm tắt để xác nhận.

`Esc` quay lại một bước. Bất kỳ tùy chọn nào trong `--group`, `--no-group`, `--shortcut` hoặc `--no-login`, hoặc chạy không có terminal (ví dụ từ một pipe), đều bỏ qua trình hướng dẫn. Khi đó bắt buộc phải có tên, và không có lệnh tắt nào được tạo trừ khi `--shortcut` nêu tên một lệnh.

Không có `--group`, tài khoản mới tham gia nhóm `default`. Tài khoản đầu tiên tạo nhóm này, dùng chung máy chủ MCP, cài đặt, `CLAUDE.md` và tiện ích mở rộng; xem [dùng chung bằng nhóm](/features/share). Trên Windows, chỉ có thể liên kết tệp khi bật Developer Mode. Khi không có chế độ này, tài khoản đầu tiên tạo `default` chỉ với `mcp` và `extensions` và thông báo điều đó (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), và ccam từ chối `--group` với một nhóm dùng chung settings, memory hoặc history.

`ccam list` hiển thị slot, tên, lệnh tắt, email, các nhóm và lần dùng gần nhất, với `*` đứng trước tài khoản mặc định. Email là email mà thư mục đang đăng nhập hiện tại (lấy từ `.claude.json` của nó); khi email đó khác với email ccam đã ghi lại, email đã ghi lại được hiển thị phía sau dưới dạng `(was …)`:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

```console
$ ccam list
   #  Name      Shortcut  Email                                     Groups   Last used
*  1  work      ccwork    other@company.com (was work@example.com)  default  2m ago
```

Với `--json`, mỗi tài khoản là một đối tượng gồm `slot`, `shortcut`, `name`, `email` (của thông tin đăng nhập hiện tại, như trong bảng), `default`, `dir`, `last_used` (bị bỏ đi khi tài khoản chưa từng được khởi động), `groups` và `share`. `share` ánh xạ từng mục (`mcp`, `settings`, `memory`, `extensions`, `history`) tới nhóm mà mục đó đến từ, hoặc tới `"local"`.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `ccam add --no-login` | Tạo tài khoản mà không đăng nhập; đăng nhập sau bằng `ccam login <name>` | tắt |
| `ccam add --group <group>` | Tham gia nhóm này thay vì `default`. Có thể lặp lại; `default` được tạo khi chưa tồn tại | `default` |
| `ccam add --no-group` | Không tham gia nhóm nào: tài khoản không dùng chung gì | tắt |
| `ccam add --shortcut <command>` | Tạo thêm lệnh tắt này trong `~/.local/bin` | không có |
| `ccam add --api --key-from <source>` | Tài khoản dùng API key thay cho đăng nhập; xem [Tài khoản API key](/features/api) | tắt |
| `ccam list --json` | In các tài khoản dưới dạng JSON | tắt |
| `ccam remove -y`, `--yes` | Xóa mà không hỏi | tắt |
| `ccam logout -y`, `--yes` | Đăng xuất mà không hỏi | tắt |
| `ccam import --name <name>` | Tên tài khoản cho thư mục bạn đưa vào; bị từ chối nếu không có thư mục (`--name needs a folder: ccam import <folder> --name ops`) | lấy từ tên thư mục |
| `ccam import -y`, `--yes` | Nhập mọi thư mục tìm thấy mà không hỏi | tắt |
| `--dry-run` | Được `ccam list` và `ccam completion` hỗ trợ, vì chúng không thay đổi gì. Bị `add`, `login`, `logout`, `use`, `rename`, `remove` và `import` từ chối, với thông báo `` --dry-run is not supported by `ccam add`; nothing was changed `` (nêu tên lệnh) | tắt |

`ccam login` (ngoài `--expired`), `ccam use` và `ccam rename` không có tùy chọn riêng. `ccam completion` nhận `bash`, `zsh` hoặc `powershell`; giá trị khác sẽ thất bại với `only bash, zsh and powershell are supported`.

Tên tài khoản dài từ 1 đến 32 ký tự gồm `a-z`, `0-9`, `-` và `_`, bắt đầu bằng một chữ cái hoặc chữ số. Tên không được chỉ gồm chữ số, vì `ccam 2` nghĩa là slot 2. Tên cũng không được trùng với một lệnh của ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), nên `ccam <name>` không bao giờ mơ hồ. Danh sách này là `config.Reserved` trong `internal/config`.

## Cách hoạt động

Thêm một tài khoản sẽ tạo thư mục của nó, tham gia các nhóm của nó, lưu cấu hình, đăng nhập và tạo lệnh tắt, theo đúng thứ tự đó. Khi bạn chạy `ccam add work --shortcut ccwork`, ccam thực hiện các bước sau (`addProfile` trong `internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. ccam kiểm tra tên (`config.ValidateName`) và lệnh tắt (`shim.Check`) trước khi tạo bất cứ thứ gì.
2. ccam tạo `~/.local/share/ccam/profiles/work/` với quyền `0700` và gán cho tài khoản số slot trống thấp nhất, nên số của một tài khoản đã xóa sẽ được dùng lại (`profile.Add`, `Config.NextSlot`). Tài khoản đầu tiên trở thành tài khoản mặc định.
3. Tài khoản tham gia các nhóm của nó (`default`, trừ khi `--group` hoặc `--no-group` nói khác; tài khoản đầu tiên tạo `default`), và ccam liên kết thư mục mới, còn trống vào các nhóm đó, ví dụ `settings.json` tới `groups/default/settings.json` (`applyShare`, với `share.Build` rồi `Apply`). Nếu bước này thất bại, thư mục mới bị xóa, vì một thư mục còn sót lại sẽ chặn mọi lần thử lại với `folder already exists`.
4. Tài khoản được lưu vào `~/.config/ccam/config.toml` qua `fsutil.AtomicWrite`: một tệp tạm, `fsync`, rồi `rename`, nên một sự cố không bao giờ để lại cấu hình mới ghi được một nửa.
5. ccam khởi động `claude auth login` với `CLAUDE_CONFIG_DIR` trỏ tới thư mục tài khoản, và `--email` với email gần nhất của tài khoản để trang đăng nhập được điền sẵn (`runLogin`). Nếu đăng nhập thất bại, tài khoản vẫn được giữ, và ccam in ra lỗi cùng lời nhắc chạy `ccam login work` sau.
6. ccam đọc `oauthAccount` từ `profiles/work/.claude.json` (`profile.ReadIdentity`): email, ID của tài khoản và của tổ chức, tên tổ chức và gói (`claude_max`, `claude_team` …). Tệp đó có thể nặng vài megabyte, nên ccam chỉ quét tìm một khóa đó thay vì phân tích cả tệp (`TopLevelValue`). ccam lưu các thông tin này trong `[profile.login]` ở `config.toml`.
7. ccam liên kết `~/.local/bin/ccwork` tới tệp thực thi ccam, và khi khởi động tệp này tra `argv[0]` trong `config.toml`, đồng thời ghi lại liên kết trong `~/.local/state/ccam/shortcuts.json` (`shim.Install`). Nếu không tạo được liên kết, tài khoản vẫn được thêm, ccam cho biết lý do (`Could not create shortcuts: …`) và bản tóm tắt hiển thị `shortcut: ccwork (not created; see above)`.

Việc đăng nhập có thể mất vài phút, và trong lúc đó một tiến trình ccam khác có thể lưu cấu hình. Vì vậy các lệnh phải chờ sẽ tải lại `config.toml`, áp dụng đúng một thay đổi (ở đây là email) rồi lưu (`env.update`). `ccam remove` cũng tải lại sau câu hỏi xác nhận vì lý do tương tự.

`ccam login` chạy bước 5 và 6 cho một tài khoản đã có. Claude Code chỉ dùng email như một gợi ý cho trang đăng nhập, nên bạn có thể đăng nhập vào một tài khoản Claude khác với trước đây; ccam so sánh lần đăng nhập mới với `[profile.login]` và báo cho bạn (`ccam: this login is other@company.com, but work was dev@company.com before.`). Nó giữ lần đăng nhập mới và cho bạn biết cách hoàn tác: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

Một lần đăng nhập vào tài khoản Claude khác cũng bỏ [bản ghi giới hạn](/features/usage#the-limit-record) của tài khoản, vì các giới hạn đó là của tài khoản kia. Khi một tài khoản ccam khác đã đăng nhập vào cùng tài khoản Claude đó, ccam cảnh báo rằng hai tài khoản dùng chung giới hạn sử dụng.

`ccam login --expired` lần lượt đăng nhập mọi tài khoản chưa từng đăng nhập hoặc có thông tin đăng nhập đã hết hạn: refresh token của nó đã hết hạn, hoặc Claude Code đã xóa thông tin đăng nhập sau khi máy chủ từ chối nó. Trước khi trình duyệt đầu tiên mở ra, lệnh này kiểm tra rằng `config.toml` có thể được lưu.

`ccam logout` chạy `claude auth logout` trong tài khoản và xóa email đã lưu; bản ghi đăng nhập vẫn được giữ, nên lần đăng nhập tiếp theo được kiểm tra với nó. Claude Code thu hồi thông tin đăng nhập và xóa những gì nó lưu cho tài khoản, kể cả các lần đăng nhập máy chủ MCP. Vì vậy ccam cảnh báo bạn và hỏi trước (`--yes` bỏ qua câu hỏi) khi:

- thông tin đăng nhập giữ các lần đăng nhập MCP;
- tài khoản đang chạy (các phiên của nó cũng bị đăng xuất);
- thư mục của một tài khoản khác giữ cùng thông tin đăng nhập (một thư mục được sao chép), thứ mà việc đăng xuất cũng sẽ thu hồi.

Ví dụ:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

Để đăng nhập bằng một tài khoản Claude khác, chỉ cần `ccam login <name>` là đủ, và nó giữ lại các lần đăng nhập MCP.

`ccam list` thêm cột `Organization` khi các tài khoản thuộc nhiều hơn một tổ chức (gói cá nhân hiển thị `personal`), và `ccam list --json` có `plan` và `org` cho mọi tài khoản.

### Đăng nhập trên máy từ xa

Qua SSH, phần trình duyệt của `claude auth login` có thể không hoạt động. Claude Code chờ việc đăng nhập trên một máy chủ callback trên máy từ xa (`http://localhost:<port>/callback`), nhưng trình duyệt hoàn tất việc đăng nhập lại chạy trên máy tính của bạn, nơi địa chỉ đó không dẫn đến đâu và trình duyệt hiển thị lỗi như `ERR_CONNECTION_REFUSED`. Địa chỉ đó vẫn mang mã dùng một lần, nên ccam có thể chuyển nó đi:

<!-- motion:accounts-relay -->

ccam không giao terminal của bạn cho Claude Code (`runLogin`). Nó nối đầu vào chuẩn của Claude Code với một pipe và sao chép từng dòng bạn gõ vào đó (`relayLogin`), nên mã được dán ở lời nhắc `Paste code here` của Claude Code vẫn đến nơi. Khi một dòng là địa chỉ callback, ccam chỉ chấp nhận nó nếu máy chủ là `localhost`, `127.0.0.1` hoặc `[::1]`, có cổng, đường dẫn là `/callback` và có cả `code` lẫn `state`; có thể bỏ `http://` (`callbackURL`). Sau đó ccam gửi yêu cầu tới địa chỉ đó trên máy chủ với thời gian chờ tối đa 10 giây và không đi theo chuyển hướng, giống như trình duyệt lẽ ra đã làm (`deliverCallback`). Yêu cầu đến máy chủ callback của `claude auth login` đang chờ (đường nét đứt trong sơ đồ). Claude Code nhận mã và lưu thông tin đăng nhập, trên Linux là trong `.credentials.json` ở thư mục tài khoản. Một dòng được dán vào không bao giờ khiến ccam liên lạc với máy chủ khác.

ccam in ra dòng này trước khi bắt đầu đăng nhập:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

Khi thành công, nó in ra `→ Passed the sign-in to Claude Code on this machine.`

Trong lúc việc đăng nhập chờ nhập liệu trên terminal, ccam tắt chế độ echo và chế độ dòng của terminal, và hiển thị một `*` cho mỗi ký tự bạn gõ, nên mã hoặc địa chỉ callback được dán vào sẽ không hiện trên màn hình. ccam tự xử lý việc chỉnh sửa dòng (`readLine`): `Enter` gửi dòng, `Backspace` xóa một ký tự, và `Ctrl+D` trên một dòng trống kết thúc việc nhập. Terminal được khôi phục khi đăng nhập kết thúc, và cả khi nhận `Ctrl+C`, `SIGTERM` hoặc `SIGHUP`, những tín hiệu vẫn kết thúc ccam như trước. Khi đầu vào chuẩn và đầu ra lỗi chuẩn không phải là terminal, các dòng được sao chép nguyên như khi chúng đến.

### Đổi tên, xóa và nhập

`ccam rename` và `ccam remove` tác động lên thư mục của tài khoản và lên những gì ccam lưu ở nơi khác dưới tên tài khoản. ccam liệt kê các đường dẫn đó trong `nameState`:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- các bộ nhớ đệm mức sử dụng `~/.local/state/ccam/usage/<name>.json` và `<name>.statusline.json`;
- bản ghi giới hạn `<name>.limit.json`;
- các ghi chú phiên trong `~/.local/state/ccam/handoff/<name>/`.

`ccam remove work` xử lý chúng theo thứ tự sau:

<!-- motion:accounts-remove -->

1. `ccam remove` trước tiên từ chối nếu tài khoản đang chạy, rồi mới hỏi, và kiểm tra lại sau khi bạn trả lời (việc này phát hiện một phiên được khởi động trong lúc lời nhắc đang mở).
2. Nó kiểm tra rằng `config.toml` có thể được lưu trước khi thay đổi bất cứ thứ gì.
3. Nếu tài khoản dùng chung tiện ích mở rộng, ccam ghi lại `installed_plugins.json` và `known_marketplaces.json` của nhóm (`RewritePluginPaths`), vì Claude Code lưu đường dẫn plugin tuyệt đối có thể đi qua thư mục của tài khoản này.
4. Tiếp theo, một tài khoản được quản lý (tài khoản có thư mục do ccam tạo, khác với tài khoản được đăng ký bằng `ccam import`) sẽ bị đăng xuất: `claude auth logout` chạy trong thư mục của nó, nên không có thông tin đăng nhập còn hiệu lực nào nằm chờ trong thùng rác (hoặc, trên macOS, trong một mục Keychain được đặt tên theo thư mục, mà một tài khoản cùng tên sau này sẽ thừa hưởng). Bước này được bỏ qua với tài khoản API key và khi thư mục của một tài khoản khác giữ cùng thông tin đăng nhập, thứ mà việc đăng xuất cũng sẽ thu hồi. Một lần đăng xuất thất bại sẽ được báo, và việc xóa vẫn tiếp tục.
5. ccam chuyển thư mục vào `~/.local/share/ccam/trash/<name>-<time>/` (`profile.Trash`), tài khoản bị gỡ khỏi `config.toml`, các đường dẫn `nameState` bị xóa để một tài khoản cùng tên sau này không thừa hưởng chúng, và liên kết lệnh tắt bị gỡ bỏ: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

Nếu rốt cuộc `config.toml` không lưu được, thư mục được đưa về từ thùng rác. Để khôi phục một tài khoản đã xóa, hãy chuyển thư mục từ thùng rác về chỗ cũ và đăng nhập lại.

`ccam rename` đổi tên thư mục và ghi lại các đường dẫn plugin trong `plugins/` của nó từ thư mục cũ sang thư mục mới. Nó chuyển các đường dẫn `nameState` sang tên mới và giữ tài khoản làm tài khoản mặc định nếu trước đó nó là mặc định. Nó cũng cập nhật các ghi chú phiên, để `ccam next --continue` vẫn tìm thấy chúng. Lệnh tắt giữ nguyên tên. Nếu việc ghi lại plugin thất bại, hoặc `config.toml` không lưu được, thư mục được đổi lại tên cũ.

Trên macOS, và trên Windows khi Claude Code lưu thông tin đăng nhập trong Credential Manager, bản thân thư mục giữ tên cũ, vì thông tin đăng nhập và các lần đăng nhập MCP được lưu theo đường dẫn của thư mục, và một thư mục bị di chuyển sẽ bị đăng xuất. Khi đó cấu hình ghi thư mục dưới dạng `folder` bên cạnh tên mới, và ccam thông báo `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

`ccam import` đăng ký một thư mục ngay tại chỗ, dưới dạng `dir` trong `config.toml`. Nó không bao giờ di chuyển hay sao chép thư mục: một refresh token bị sao chép sẽ khiến các bản sao đăng xuất lẫn nhau.

Không có đối số thư mục, nó tìm các thư mục `~/.claude` và `~/.claude-*` chưa được đăng ký và hỏi về từng thư mục. Tên được lấy từ thư mục (`~/.claude` là `main`, `~/.claude-work` là `work`, `~/.claude-2` là `acct2`). Nếu ccam không suy ra được tên, hoặc tên đã được dùng, nó liệt kê thư mục kèm lệnh `--name` cần chạy. Một thư mục tìm thấy mà thất bại sẽ được báo là `Skipped <path>: <reason>` và các thư mục còn lại vẫn được nhập.

ccam hỏi tất cả câu hỏi trước, rồi tải lại `config.toml` và đăng ký các thư mục, nên những thay đổi mà các lệnh ccam khác đã lưu trong lúc đó vẫn được giữ. Các dòng `✓` xuất hiện sau khi lưu.

Khi không có `CLAUDE_CONFIG_DIR`, Claude Code giữ `.claude.json` trong thư mục home của bạn, cạnh `~/.claude`, nên khi bạn nhập `~/.claude` mà nó không có `.claude.json` riêng, ccam tạo liên kết `~/.claude/.claude.json` tới `~/.claude.json`. Khi đó `claude` thông thường và `ccam main` thấy cùng một tài khoản.

Nếu `~/.claude` có `.claude.json` riêng bên cạnh `~/.claude.json`, ccam để nguyên cả hai và thông báo điều đó. `ccam main` sẽ dùng tệp bên trong còn `claude` thông thường dùng tệp bên ngoài, khiến thông tin đăng nhập và máy chủ MCP của chúng tách rời nhau. Để chỉ giữ một tệp, hãy thay tệp bên trong bằng một liên kết ([doctor](/features/doctor) cũng báo điều này).

Khi thư mục giữ cùng thông tin đăng nhập với một tài khoản đã đăng ký (một thư mục được sao chép), ccam thông báo điều đó: bên nào gia hạn trước sẽ đăng xuất bên kia, vì vậy hãy chạy `ccam login <name>` trước khi dùng nó.

Tài khoản được nhập không thuộc nhóm nào cho đến khi bạn thêm nó vào một nhóm (`ccam group join <group> <name>`), và không có lệnh tắt nào được tạo.

`ccam completion bash|zsh` in ra một script tự hoàn thành (completion), script này cũng hoàn thành cả tên tài khoản.

## An toàn và các trường hợp đặc biệt

- `ccam rename` và `ccam remove` từ chối khi tài khoản đang chạy (`refuseRunning`), với ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``, vì phiên vẫn tiếp tục ghi vào thư mục. Việc kiểm tra bao gồm Claude Code bất kể nó được khởi động theo cách nào: qua ccam, `claude` thông thường, một IDE, các phiên nền và daemon của chúng.
- `ccam rename`, `ccam remove` và `ccam logout` từ chối khi một tiến trình ccam khác đang gia hạn hoặc thay đổi thông tin đăng nhập đó: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- Khi một lần thay đổi nhóm đang diễn ra trong một tiến trình ccam khác, `ccam rename` và `ccam remove` từ chối với `another ccam process is changing the groups right now; try again when it is done`; sau khi một lần thay đổi bị dừng giữa chừng, với ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.
- `ccam remove` trên một tài khoản đã nhập giữ lại thư mục và thông tin đăng nhập của nó, cho `claude` thông thường. Nếu tài khoản nằm trong các nhóm, trước tiên nó rời các nhóm đó như `ccam group leave`: mọi mục dùng chung trở thành một bản sao của bản trong nhóm ở thư mục, và các máy chủ MCP của nhóm được đưa vào `.claude.json` của nó. Sau đó nó được hủy đăng ký: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. Các bản sao lưu `.claude.json` của nó mà ccam đã tạo được chuyển vào `<folder>/.ccam-backup/mcp-<time>`; nếu không chuyển được, chúng ở lại trong `~/.local/state/ccam/mcp/<name>/backups` và thông báo sẽ nói rõ. `~/.claude` của bạn không bị chuyển vào thùng rác.
- Hai tài khoản không thể có cùng một thư mục: `config.toml` nêu một thư mục hai lần sẽ không tải được.
- `ccam add` từ chối một tên mà thư mục của nó đã tồn tại trong `profiles/`, và gợi ý dùng `ccam import`. Trên macOS, nó cũng dừng lại khi Keychain vẫn còn giữ thông tin đăng nhập hoặc các lần đăng nhập MCP cho đường dẫn của thư mục mới, do một phiên bản ccam cũ hơn để lại, và nêu lệnh `security delete-generic-password` để xóa chúng. Nếu `config.toml` không lưu được, thư mục mới lại bị xóa.
- `ccam import` từ chối một thư mục đã được đăng ký, một đường dẫn không phải thư mục, và một tên đã có người dùng. ccam cũng cảnh báo khi thư mục là một bản sao của `~/.claude` (cùng thông tin đăng nhập), ngay cả khi bản thân `~/.claude` chưa được đăng ký.
- `ccam logout` hỏi trước khi tài khoản đang chạy; `ccam login` không kiểm tra.
- Nếu phần chuyển tiếp không liên lạc được với Claude Code, chẳng hạn sau khi việc đăng nhập đã kết thúc, nó in ra `ccam: could not pass the sign-in to Claude Code:` kèm lý do. Lỗi HTTP thường có nghĩa là mã đã hết hạn; hãy bắt đầu đăng nhập lại.

## Liên quan

- [Khởi động Claude Code](/features/launch)
- [Lệnh tắt](/features/shortcuts)
- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Dùng nhiều tài khoản Claude Code trên một máy tính](/guides/multiple-claude-code-accounts)
