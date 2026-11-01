---
title: Khởi động Claude Code
description: "Chạy Claude Code bằng tài khoản bất kỳ với ccam <name|number>, ccam run hoặc bảng điều khiển; đối số được chuyển nguyên vẹn cho claude, --dry-run in ra lệnh."
updated: 2026-11-01
order: 1
translation: machine
sourceHash: 11df08608a1e
---

`ccam work` khởi động Claude Code bằng tài khoản `work`, và mọi thứ sau tên được chuyển cho `claude` mà không thay đổi. Hãy dùng nó ở bất cứ đâu bạn vẫn gõ `claude`: trong terminal, trong script, hoặc khi có đầu vào được pipe vào.

## Cách dùng

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Hãy đưa tên của một tài khoản hoặc số của nó trong `ccam list`. ccam không đọc các đối số sau tên, nên `ccam work --help` hiển thị trợ giúp của Claude Code. Dấu `--` sau tên là không bắt buộc. Đầu vào chuẩn, mã thoát và tín hiệu hoạt động giống như với `claude`.

Khi đối số đầu tiên bắt đầu bằng `-` và không phải một cờ của chính ccam, ccam chọn tài khoản cho bạn và chuyển tất cả các đối số cho Claude Code. Nó dùng tài khoản được đặt cho thư mục hiện tại bằng [`ccam folder`](/features/folder), nếu không thì dùng tài khoản mặc định: tài khoản được đặt bằng `ccam use`, hoặc tài khoản có số thấp nhất nếu bạn chưa đặt.

### Bảng điều khiển

`ccam` không kèm đối số sẽ mở [bảng điều khiển](/features/dashboard) khi đầu vào và đầu ra chuẩn là terminal, và in ra `ccam list` trong các trường hợp khác. Di chuyển tới một tài khoản (bằng các phím mũi tên, hoặc gõ số của nó) và nhấn `Enter` để khởi động nó; `n` khởi động tài khoản còn nhiều hạn mức nhất và `o` đề xuất `--continue`, `--resume` hoặc các đối số khác. Bảng điều khiển đóng lại trước khi Claude Code khởi động, nên terminal được khôi phục trước.

**Seat của gói Team và Enterprise.** Việc tự động chọn tài khoản (`ccam next`, `--fallback`, `ccam pick`, phím `n` trên bảng điều khiển và, khi bạn bật nó, [`auto_switch`](/features/next#auto-switch)) cũng chuyển đổi giữa các seat của cùng một tổ chức Claude Team hoặc Enterprise, điều này có thể vi phạm điều khoản của Anthropic. Hãy đọc [Câu hỏi thường gặp](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) và xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó.

Khi shell có đặt `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` hoặc `ANTHROPIC_AUTH_TOKEN`, Claude Code có thể dùng khóa hoặc token đó thay cho thông tin đăng nhập riêng của tài khoản, nên ccam cảnh báo bạn trên đầu ra lỗi chuẩn trước khi khởi động (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Nó chỉ đọc các biến này; một [tài khoản API key](/features/api) thì xóa chúng đi.

Theo mặc định, ccam tự thay thế mình bằng Claude Code, nên không có gì khác tiếp tục chạy. Với `relay = true` (hoặc [`auto_switch = true`](/features/next#auto-switch) trên Linux và macOS) ở đầu `config.toml`, phía trên `[[group]]`, `[[profile]]` hoặc `[[folder]]` đầu tiên, ccam thay vào đó khởi động các phiên tương tác như một tiến trình con và tiếp tục chạy song song với chúng. Nhờ vậy `/ccam:switch`, hoặc `auto_switch` khi chạm giới hạn sử dụng, có thể tự tiếp tục cuộc hội thoại trên một tài khoản khác, với các tùy chọn, chế độ quyền, mức effort và agent của phiên; trước khi `/ccam:switch` chuyển nó sang một tài khoản của tổ chức khác, ccam sẽ hỏi bạn ([chi tiết](/features/next#relay)).

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--dry-run` | In ra biến môi trường và lệnh mà ccam sẽ chạy, rồi thoát | tắt |
| `--fallback` | Cho các lần chạy `-p`/`--print`: khi tài khoản chạm giới hạn sử dụng, chạy lại cùng lệnh đó bằng tài khoản đã đăng nhập còn nhiều hạn mức nhất ([chi tiết](/features/next#run-again-at-the-limit)) | tắt |
| `-v`, `--version` | In ra `ccam <version>` (chỉ khi là đối số đầu tiên) | |
| `-h`, `--help` | Trợ giúp của ccam, dưới dạng `ccam --help` hoặc `ccam run --help` | |

Quy tắc về đối số:

- ccam chỉ đọc cờ của chính nó ở trước tên tài khoản. Sau tên, mọi đối số đều thuộc về Claude Code, kể cả `--dry-run` và `--help`.
- `--dry-run` hoạt động dưới dạng `ccam --dry-run work`, `ccam run --dry-run work` và `ccam --dry-run` (cho tài khoản của thư mục hoặc tài khoản mặc định). `--dry-run` và `--fallback` có thể đứng theo thứ tự bất kỳ.
- `ccam run <name|number> [claude args]` là dạng đầy đủ của `ccam <name>`. `ccam run` không kèm tên sẽ dùng tài khoản của thư mục hoặc tài khoản mặc định.
- Tên tài khoản không được trùng với tên lệnh của ccam như `add`, `list`, `run` hoặc `usage`, nên `ccam <word>` không bao giờ mơ hồ.

Lần chạy thử in ra một dòng mà bạn có thể dán vào shell, hoặc vào PowerShell trên Windows. Dòng này bắt đầu bằng những gì lần khởi động thay đổi trong biến môi trường: luôn có `CLAUDE_CONFIG_DIR` và `CCAM_PROFILE`, `env -u <name>` cho một biến bị bỏ đi (chẳng hạn `CCAM_RELAY` được thừa hưởng), và, với một [tài khoản API key](/features/api), endpoint, mô hình của nó và nơi lấy khóa (`ANTHROPIC_API_KEY="$WORK_KEY"`), không bao giờ là chính khóa:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

Dòng PowerShell đặt lại các biến sau đó, như `VAR=x cmd` làm trong shell, và giao dòng lệnh cho Claude Code qua `Start-Process`, nên Windows PowerShell 5.1 và PowerShell 7 chạy nó theo cùng một cách. Khi một lần khởi động đặt một biến thành giá trị rỗng, một dòng `#` bên dưới cho biết PowerShell trước 7.5 sẽ xóa biến đó thay vào đó. Với `claude.cmd` được npm cài, dòng này chạy lệnh `cmd.exe` mà lần khởi động sẽ dùng (xem bên dưới), và một đối số mà lần khởi động sẽ từ chối ở đó cũng làm lần chạy thử thất bại.

## Cách hoạt động

Mỗi lần khởi động đi qua cùng các bước, từ dòng lệnh đến `exec`:

<!-- motion:launch -->

Mỗi ô sáng lên khi bước của nó chạy, và các dòng bên dưới là biến môi trường và đối số mà Claude Code khởi động cùng.

1. ccam bắt đầu trong `app.Main`, hàm mà `cmd/ccam/main.go` gọi. Nếu ccam được khởi động dưới một tên khác, đó là một [lệnh tắt](/features/shortcuts) và đi thẳng tới trình khởi chạy.
2. ccam xem `os.Args` trước khi cobra tham gia (`app.parse`). Chỉ `ccam` đứng một mình (bảng điều khiển), một tên lệnh của ccam (`config.Reserved`), `-h`/`--help` và `-v`/`--version` mới đi tới `cli.Execute`. Mọi thứ khác là một lần khởi động.
3. ccam đọc `~/.config/ccam/config.toml`, tệp duy nhất nó phân tích (`launcher.Launch`).
4. ccam tìm tài khoản theo tên hoặc số, bỏ một dấu `--` đứng đầu, và tìm tệp thực thi `claude` (`launcher.Prepare`, `ResolveClaude`). Khi không có tên, nó lấy tài khoản của thư mục, so sánh thư mục làm việc với các mục `[[folder]]` của `config.toml` dưới dạng chuỗi, không phân biệt hoa thường trên macOS và Windows (`HereProfile`). Kết quả tra cứu PATH được lưu đệm trong `~/.local/state/ccam/claude-path` cùng với thời điểm sửa đổi của tệp và một hash của `PATH`, nên ccam không phải tìm trong `PATH` ở mỗi lần khởi động. Tệp bị thay đổi hoặc không còn, hoặc `PATH` thay đổi, sẽ kích hoạt một lần tra cứu mới. Đường dẫn `claude_bin` trong `config.toml` (`~/` được mở rộng) phải là một tệp thực thi được; nếu không, việc khởi động dừng với `claude_bin X is not an executable file. Check claude_bin in config.toml` trước khi bất cứ thứ gì được ghi. ccam từ chối một `claude` chính là ccam, so sánh theo inode, vì điều đó sẽ tạo vòng lặp.
5. ccam so sánh thư mục tài khoản `~/.local/share/ccam/profiles/<name>/` với thiết lập dùng chung của nó, chỉ gọi `lstat` và `readlink`, một hoặc hai lời gọi hệ thống cho mỗi mục (`checkLinks`, `share.Check`). Liên kết còn thiếu được tạo ngay lập tức. Bất cứ thứ gì có thể làm mất dữ liệu đều dừng việc khởi động.
6. ccam thêm các máy chủ MCP dùng chung (`applyMCP`; xem bên dưới).
7. ccam ghi lại phiên (`session.Register`). Nó ghi một tệp đánh dấu đặt tên theo PID của nó, chứa thời điểm bắt đầu tiến trình, vào `~/.local/state/ccam/locks/profile-<name>/<pid>`, bất kể `XDG_RUNTIME_DIR` là gì, để mọi shell thấy cùng các tệp đánh dấu. Với tài khoản dùng chung lịch sử hội thoại, ccam trước tiên tìm một phiên đang chạy trong `history-<group>-<hash>`, trong đó hash là của thư mục hiện tại với các liên kết tượng trưng đã được phân giải, và cũng đăng ký ở đó (`historyCheck`).

   `exec` giữ nguyên PID, nên khi đó tệp đánh dấu chỉ đúng tiến trình Claude Code. ccam dùng tệp đánh dấu chứ không giữ khóa, vì các tiến trình do Claude Code khởi động sẽ thừa hưởng khóa và giữ nó sau khi phiên kết thúc. Một tệp đánh dấu không ghi được không làm dừng việc khởi động: ccam báo `ccam: could not mark the session as running (…); starting Claude Code anyway`.

   Các kiểm tra khác của ccam xem một tài khoản có đang chạy không (đổi tên, xóa, thay đổi nhóm, `doctor --fix`, đồng bộ MCP) cũng đọc các bản ghi riêng của Claude Code trong thư mục tài khoản (`sessions/<pid>.json`, `daemon.lock`), nên chúng cũng thấy Claude Code được khởi động không qua ccam: `claude` thông thường, một IDE, các phiên nền và daemon của chúng.
8. ccam đặt tiêu đề terminal thành `claude · <name>` (khi đầu ra lỗi chuẩn là terminal), ghi thời điểm khởi động vào `~/.local/state/ccam/last-used/<name>` mà không `fsync`, và gọi `syscall.Exec` (`Plan.Exec`).

Claude Code chạy với `CLAUDE_CONFIG_DIR` trỏ tới thư mục tài khoản và `CCAM_PROFILE` là tên tài khoản. ccam xóa mọi giá trị được thừa hưởng của các biến này, và cả `CCAM_RELAY` được thừa hưởng. Với một thư mục `~/.claude` được đăng ký bằng `ccam import` trên macOS, nó còn đặt `CLAUDE_SECURESTORAGE_CONFIG_DIR` thành giá trị rỗng, để tài khoản và `claude` thông thường dùng cùng một thông tin đăng nhập trong Chuỗi khóa (Keychain); trên Windows, nó làm điều tương tự khi Claude Code giữ thông tin đăng nhập trong Credential Manager.

### Vì sao dùng exec

ccam thay tiến trình của chính nó bằng `claude` (`syscall.Exec`). Nhờ vậy:

- ccam không ở lại trong bộ nhớ;
- phím bấm, tín hiệu và mã thoát là của chính Claude Code;
- đối số luôn là một mảng, không bao giờ là chuỗi shell.

Windows không có `exec`. Ở đó, ccam khởi động `claude.exe` như một tiến trình con với cùng biến môi trường, đối số và console, chờ nó và thoát với mã thoát của nó (`Plan.Exec`). Một `claude.cmd` do npm cài chạy trong `cmd.exe` với đường dẫn đầy đủ của nó và mọi đối số được đặt trong ngoặc kép, để `cmd.exe` không đọc đối số nào như cú pháp của riêng nó (`&`, `%PATH%`); một đối số có ngắt dòng không thể truyền theo cách đó và bị từ chối, kèm lời khuyên cài Claude Code bản gốc.

Ctrl+C đến thẳng Claude Code từ console, nên ccam bỏ qua nó. Tệp đánh dấu phiên mang PID của ccam, tiến trình sống đúng bằng thời gian của Claude Code, và bị xóa khi Claude Code thoát. Trong lúc chờ, ccam không làm gì khác.

Công việc sau một phiên, như phát hiện các máy chủ được thêm bằng `claude mcp add`, diễn ra ở lần khởi động tiếp theo: ccam so sánh thời điểm sửa đổi của `.claude.json` với `~/.local/state/ccam/mcp/<name>/seen` và chỉ phân tích tệp khi nó đã thay đổi.

### Cách một dòng lệnh được định tuyến

Ba dòng lệnh, ba đường đi qua `internal/app`:

<!-- motion:launch-route -->

Lệnh tắt được nhận ra qua `argv[0]`, `ccam <name>` qua đối số đầu tiên, và cả hai đều kết thúc ở `syscall.Exec`. Chỉ một lệnh của ccam như `ccam usage` mới đến cobra.

### Máy chủ MCP dùng chung

Với chiến lược mặc định `flag`, ccam truyền các máy chủ MCP của nhóm cho Claude Code dưới dạng một đối số, đặt trước các đối số của bạn:

<!-- motion:launch-mcp -->

ccam chỉ thêm đối số này nếu `mcp.json` của nhóm tồn tại, và đặt nó ngay sau `claude` (`mcp.FlagArgs`, `Plan.InsertArgs`). Khi đó `--mcp-config` hoặc `--strict-mcp-config` của riêng bạn vẫn được áp dụng chồng lên. Đó là một đối số duy nhất, `--mcp-config=<path>`, vì `--mcp-config` của Claude Code nhận nhiều giá trị: nếu viết thành hai đối số, nó sẽ đọc cả lời nhắc theo sau như một tệp cấu hình.

ccam quyết định nơi cờ này không được đặt (`claudeVerb`). Nó đọc lệnh con của Claude Code từ các đối số theo cách Claude Code làm, bỏ qua các tùy chọn cùng giá trị của chúng. `remote-control` hoặc `rc` không nhận cờ (bên dưới), một động từ phiên nền đứng đầu (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`) cũng vậy. `ccam work mcp list`, `get` và `remove` cũng không nhận, vì Claude Code từ chối nó ở đó; thay vào đó ccam báo trên đầu ra lỗi chuẩn rằng các máy chủ của nhóm cũng được tải trong các phiên của tài khoản:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Mọi lệnh con khác và mọi phiên đều giữ cờ.

Với `mcp_strategy = "sync"` trong `config.toml`, ccam ghi các máy chủ vào `.claude.json` của tài khoản thay vào đó, và bỏ qua việc này nếu tài khoản đang chạy. Một lần khởi động cũng bỏ qua việc này khi cả `mcp.json` của nhóm lẫn `.claude.json` đều không thay đổi kể từ lần kiểm tra trước. Nếu tài khoản có máy chủ riêng trùng tên với một máy chủ của nhóm, đồng bộ giữ máy chủ của tài khoản và báo điều đó một lần. Đồng bộ cũng đặt lại một bản sao đã mất khỏi `.claude.json` (`put back group default's MCP servers …`), và chỉ thay đổi hoặc xóa những bản sao nó đã ghi mà chưa ai thay đổi kể từ đó.

### Remote Control

`ccam work rc` (hoặc `ccam work remote-control`) khởi động Remote Control của Claude Code bằng tài khoản `work`. `claude remote-control` từ chối `--mcp-config` đặt trước động từ của nó và không nhận cờ này sau động từ, và các phiên của nó chỉ đọc máy chủ MCP từ `.claude.json` của tài khoản. Vì vậy chiến lược flag hoạt động khác ở đây:

<!-- motion:launch-rc -->

1. ccam tìm thấy `remote-control` hoặc `rc` là động từ của Claude Code, kể cả sau các tùy chọn và giá trị của chúng (`ccam work --model opus rc`), và không thêm cờ (`claudeVerb`, `mcp.FlagArgs`).
2. ccam so sánh thời điểm sửa đổi của `mcp.json` của nhóm với `sync-stamp` của tài khoản, và kiểm tra rằng mọi máy chủ của nhóm vẫn còn trong `.claude.json` (`mcp.Synced`). Nếu `.claude.json` đã chứa các máy chủ hiện tại của nhóm, hoặc nhóm không có máy chủ nào, Remote Control khởi động mà không hỏi.
3. Nếu không, ccam nêu tên các máy chủ mà các phiên sẽ thiếu và hỏi:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. Khi chọn `y`, ccam sao chép chúng theo cách của chiến lược `sync` (`mcp.Sync`): dưới khóa của Claude Code, sau khi sao lưu, và giữ lại các máy chủ riêng của tài khoản. Khi chọn `N`, hoặc khi không có terminal, Remote Control khởi động mà không có chúng, và ccam gợi ý `mcp_strategy = "sync"`, cách này sao chép chúng trước mỗi lần khởi động.

ccam hỏi ở mỗi lần khởi động Remote Control cho đến khi các máy chủ được sao chép. Sau khi sao chép, nó chỉ hỏi lại khi `mcp.json` của nhóm thay đổi hoặc một máy chủ đã sao chép biến mất khỏi `.claude.json`. Bản sao trong `.claude.json` được ưu tiên hơn `--mcp-config` ở các lần khởi động bình thường sau đó, nên `ccam mcp add`, `remove` và `import` cập nhật các bản sao mà tài khoản khác đang giữ ([Máy chủ MCP dùng chung](/features/mcp)).

### Ngân sách 10 ms

Ngân sách cho đường khởi động là 10 ms ở phân vị thứ 95; trên thực tế ccam thêm khoảng 4 ms trước khi Claude Code khởi động. Đường khởi động không bao giờ dựng cây lệnh cobra, khởi động Bubble Tea hay mở kết nối mạng, và `TestPackageInitIsCheap` thất bại nếu các hàm `init()` của gói mất tổng cộng hơn 5 ms. `scripts/perf-check.sh` đo `ccam --dry-run work`, `ccam work` và một lệnh tắt bằng `hyperfine` và thất bại khi vượt 10 ms, còn `BenchmarkLaunchPath` trong `internal/launcher` chạy các bước trong một vòng lặp. Để xem thời gian của từng bước, hãy đặt `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## An toàn và các trường hợp đặc biệt

- Một tên hoặc số không xác định sẽ dừng với `no account "<name>"` và chỉ tới `ccam list`. Khi chưa có tài khoản nào, ccam yêu cầu bạn thêm một tài khoản bằng `ccam add <name>`.
- Nếu một đường dẫn dùng chung đã trở thành tệp thông thường hoặc trỏ sai chỗ, ccam dừng với `stopped because of share link problems:`, liệt kê chúng, và yêu cầu bạn chạy `ccam doctor --fix`. Nó không di chuyển dữ liệu trên đường khởi động.
- Nếu một phiên Claude Code khác dùng chung lịch sử hội thoại với tài khoản này (cùng nhóm) đang chạy trong cùng thư mục, ccam cảnh báo rằng lịch sử có thể bị trộn lẫn và hỏi `Run anyway? [y/N]`. Thư mục được so khớp với các liên kết tượng trưng đã được phân giải. Khi không có terminal (đầu vào chuẩn từ `/dev/null` được tính là không có), hoặc từ một tác vụ chạy nền, nó hiển thị cảnh báo và tiếp tục. Một lần chạy `-p` mới bắt đầu một cuộc hội thoại riêng và không nhận cảnh báo; `--version`, `--help` và các lệnh con không mở cuộc hội thoại nào, như `mcp`, `auth` hoặc `plugin`, hoàn toàn không được kiểm tra.
- `--dry-run` bỏ qua việc đồng bộ MCP, thông báo máy chủ mới, ghi chú `mcp list`, kiểm tra lịch sử, tệp đánh dấu phiên, tiêu đề và thời điểm dùng gần nhất, và không exec. Nó vẫn kiểm tra liên kết dùng chung, và tạo liên kết còn thiếu như một lần khởi động bình thường.
- `ccam work rc` khi `work` đang chạy sẽ không động đến `.claude.json`, báo điều đó, và khởi động Remote Control mà không có các máy chủ của nhóm; hãy chạy lại sau khi phiên kết thúc để sao chép chúng. Với `--dry-run`, nó in thông báo nhưng không hỏi và không sao chép.
- Với `mcp_strategy = "sync"`, việc đồng bộ thất bại sẽ in ra `ccam: MCP sync failed (launching anyway): …` và Claude Code khởi động với các máy chủ nó đã có.

## Liên quan

- [Lệnh tắt](/features/shortcuts)
- [Dùng Remote Control với ccam](/guides/remote-control)
- [Quản lý tài khoản](/features/accounts)
- [Tham chiếu lệnh](/docs/commands)
- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Giải thích về CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
