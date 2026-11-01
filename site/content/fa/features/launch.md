---
title: اجرای Claude Code
description: Claude Code را با ccam <name|number>، ccam run یا داشبورد، با هر حسابی اجرا کنید؛ آرگومان‌ها بی‌تغییر به claude می‌رسند و --dry-run فرمان را نشان می‌دهد.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: "11df08608a1e"
---

`ccam work` اجرای Claude Code را با حساب `work` آغاز می‌کند و هر چه پس از نام بیاید بدون تغییر به `claude` می‌رسد. هر جا که `claude` را تایپ می‌کردید از آن استفاده کنید: در ترمینال، در یک اسکریپت یا با ورودی‌ای که از طریق pipe داده می‌شود.

## نحوه استفاده

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

نام یک حساب یا شماره آن در `ccam list` را بدهید. ccam آرگومان‌های پس از نام را نمی‌خواند، بنابراین `ccam work --help` راهنمای Claude Code را نشان می‌دهد. گذاشتن `--` پس از نام اختیاری است. ورودی استاندارد، کد خروج و سیگنال‌ها مانند `claude` رفتار می‌کنند.

وقتی آرگومان اول با `-` شروع شود و یکی از گزینه‌های خود ccam نباشد، ccam حساب را برای شما انتخاب می‌کند و همه آرگومان‌ها را به Claude Code می‌دهد. از حسابی استفاده می‌کند که با [`ccam folder`](/features/folder) برای پوشه فعلی تعیین شده، وگرنه از حساب پیش‌فرض: حسابی که با `ccam use` تعیین شده، یا اگر حسابی تعیین نکرده باشید، حسابی که کمترین شماره را دارد.

### داشبورد

`ccam` بدون آرگومان، وقتی ورودی و خروجی استاندارد ترمینال باشند [داشبورد](/features/dashboard) را باز می‌کند و در غیر این صورت `ccam list` را چاپ می‌کند. به یک حساب بروید (با کلیدهای جهت‌نما یا با تایپ شماره آن) و برای اجرای آن `Enter` را بزنید؛ `n` حسابی را که بیشترین ظرفیت را دارد اجرا می‌کند و `o`، `--continue`، `--resume` یا آرگومان‌های دیگر را پیشنهاد می‌دهد. داشبورد پیش از شروع Claude Code بسته می‌شود، بنابراین ترمینال ابتدا به حالت عادی برمی‌گردد.

**صندلی‌های Team و Enterprise.** انتخاب خودکار حساب (`ccam next`، `--fallback`، `ccam pick`، کلید `n` در داشبورد و، وقتی آن را روشن کنید، [`auto_switch`](/features/next#auto-switch)) میان صندلی‌های همان سازمان Claude Team یا Enterprise هم جابه‌جا می‌شود، که ممکن است شرایط Anthropic را نقض کند. پیش از آنکه به آن تکیه کنید، [پرسش‌های متداول](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) را بخوانید و از سازمان خود و Anthropic بپرسید.

وقتی در shell یکی از `CLAUDE_CODE_OAUTH_TOKEN`، `ANTHROPIC_API_KEY` یا `ANTHROPIC_AUTH_TOKEN` تنظیم شده باشد، Claude Code ممکن است به‌جای ورود خود حساب از آن کلید یا توکن استفاده کند، پس ccam پیش از شروع در خطای استاندارد به شما هشدار می‌دهد (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). فقط متغیرها را می‌خواند؛ یک [حساب کلید API](/features/api) به‌جای آن، آن‌ها را پاک می‌کند.

به‌طور پیش‌فرض ccam خودش را با Claude Code جایگزین می‌کند، بنابراین هیچ چیز دیگری در حال اجرا نمی‌ماند. با `relay = true` (یا در Linux و macOS با [`auto_switch = true`](/features/next#auto-switch)) در بالای `config.toml`، بالای نخستین `[[group]]`، `[[profile]]` یا `[[folder]]`، ccam در عوض نشست‌های تعاملی را به‌عنوان فرایند فرزند اجرا می‌کند و در کنارشان در حال اجرا می‌ماند. این کار به `/ccam:switch`، یا `auto_switch` در سقف مصرف، امکان می‌دهد گفتگو را خودش روی حساب دیگری ادامه دهد، با گزینه‌ها، حالت مجوز، سطح تلاش و عامل همان نشست؛ پیش از آنکه `/ccam:switch` آن را به حسابی از سازمان دیگری جابه‌جا کند، ccam از شما می‌پرسد ([جزئیات](/features/next#relay)).

## گزینه‌ها

| گزینه | کار آن | پیش‌فرض |
| --- | --- | --- |
| `--dry-run` | محیط و فرمانی را که ccam اجرا می‌کرد چاپ می‌کند و خارج می‌شود | خاموش |
| `--fallback` | برای اجراهای `-p`/`--print`: وقتی حساب به سقف مصرفش برسد، همان فرمان را دوباره با حساب واردشده‌ای که بیشترین ظرفیت را دارد اجرا می‌کند ([جزئیات](/features/next#run-again-at-the-limit)) | خاموش |
| `-v`، `--version` | `ccam <version>` را چاپ می‌کند (فقط به‌عنوان آرگومان اول) | |
| `-h`، `--help` | راهنمای ccam، به شکل `ccam --help` یا `ccam run --help` | |

قواعد آرگومان‌ها:

- ccam گزینه‌های خودش را فقط پیش از نام حساب می‌خواند. پس از نام، همه آرگومان‌ها متعلق به Claude Code هستند، از جمله `--dry-run` و `--help`.
- `--dry-run` به شکل‌های `ccam --dry-run work`، `ccam run --dry-run work` و `ccam --dry-run` (برای حساب پوشه یا حساب پیش‌فرض) کار می‌کند. `--dry-run` و `--fallback` به هر ترتیبی می‌توانند بیایند.
- `ccam run <name|number> [claude args]` شکل بلند `ccam <name>` است. `ccam run` بدون نام از حساب پوشه یا حساب پیش‌فرض استفاده می‌کند.
- نام حساب‌ها نمی‌تواند نام فرمان‌های ccam مانند `add`، `list`، `run` یا `usage` باشد، بنابراین `ccam <word>` هرگز مبهم نیست.

اجرای آزمایشی یک خط چاپ می‌کند که می‌توانید آن را در shell، یا در Windows در PowerShell، بچسبانید. این خط با چیزهایی شروع می‌شود که اجرا در محیط تغییر می‌دهد: همیشه `CLAUDE_CONFIG_DIR` و `CCAM_PROFILE`، `env -u <name>` برای متغیری که حذف می‌کند (مانند `CCAM_RELAY` به ارث رسیده)، و برای یک [حساب کلید API](/features/api)، نقطه پایانی آن، مدل آن و جایی که کلید از آن می‌آید (`ANTHROPIC_API_KEY="$WORK_KEY"`)، هرگز خود کلید:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

خط PowerShell پس از اجرا متغیرها را به حالت قبل برمی‌گرداند، همان‌طور که `VAR=x cmd` در shell انجام می‌دهد، و خط فرمان Claude Code را از طریق `Start-Process` به آن می‌دهد، پس Windows PowerShell 5.1 و PowerShell 7 آن را به یک شکل اجرا می‌کنند. وقتی یک اجرا متغیری را با مقدار خالی تنظیم کند، یک خط `#` زیر آن می‌گوید که PowerShell پیش از 7.5 چنین متغیری را به جای آن حذف می‌کند. برای `claude.cmd` که npm نصب کرده است، این خط فرمان `cmd.exe` را که اجرا به کار می‌برد اجرا می‌کند (پایین‌تر را ببینید)، و آرگومانی که اجرا آنجا رد می‌کرد اجرای آزمایشی را هم ناموفق می‌کند.

## چگونه کار می‌کند

هر اجرا، از خط فرمان تا `exec`، همین گام‌ها را طی می‌کند:

<!-- motion:launch -->

هر جعبه هنگام اجرای گام خود روشن می‌شود و خط‌های زیر آن، محیط و آرگومان‌هایی هستند که Claude Code با آن‌ها شروع می‌شود.

1. ccam در `app.Main` شروع می‌شود که `cmd/ccam/main.go` آن را فراخوانی می‌کند. اگر ccam با نام دیگری اجرا شده باشد، یک [میان‌بر](/features/shortcuts) است و مستقیم به اجراکننده می‌رود.
2. ccam پیش از دخالت cobra به `os.Args` نگاه می‌کند (`app.parse`). فقط `ccam` به‌تنهایی (داشبورد)، نام یک فرمان ccam (`config.Reserved`)، `-h`/`--help` و `-v`/`--version` به `cli.Execute` می‌روند. هر چیز دیگری یک اجراست.
3. ccam فایل `~/.config/ccam/config.toml` را می‌خواند، تنها فایلی که تجزیه می‌کند (`launcher.Launch`).
4. ccam حساب را با نام یا شماره پیدا می‌کند، یک `--` آغازین را حذف می‌کند و فایل اجرایی `claude` را پیدا می‌کند (`launcher.Prepare`، `ResolveClaude`). بدون نام، حساب پوشه را برمی‌دارد و پوشه کاری را به‌صورت رشته، و در macOS و Windows بدون توجه به بزرگی و کوچکی حروف، با ورودی‌های `[[folder]]` در `config.toml` مقایسه می‌کند (`HereProfile`). نتیجه جست‌وجو در PATH همراه با زمان تغییر فایل و یک hash از `PATH` در `~/.local/state/ccam/claude-path` ذخیره می‌شود، بنابراین ccam در هر اجرا `PATH` را جست‌وجو نمی‌کند. اگر فایل تغییر کرده یا وجود نداشته باشد، یا `PATH` تغییر کرده باشد، جست‌وجوی تازه‌ای انجام می‌شود. مسیر `claude_bin` در `config.toml` (`~/` گسترش داده می‌شود) باید یک فایل اجرایی باشد؛ در غیر این صورت اجرا پیش از نوشتن هر چیزی با `claude_bin X is not an executable file. Check claude_bin in config.toml` متوقف می‌شود. ccam فایل `claude` را که خود ccam باشد (با مقایسه inode) رد می‌کند، چون به یک حلقه بی‌پایان می‌انجامد.
5. ccam پوشه حساب `~/.local/share/ccam/profiles/<name>/` را با تنظیمات اشتراک‌گذاری آن مقایسه می‌کند و فقط `lstat` و `readlink` را فراخوانی می‌کند، یعنی یک یا دو فراخوانی سیستمی برای هر مورد (`checkLinks`، `share.Check`). پیوندی که وجود نداشته باشد در همان لحظه ساخته می‌شود. هر چیزی که ممکن است به از دست رفتن داده بینجامد اجرا را متوقف می‌کند.
6. ccam سرورهای MCP مشترک را اضافه می‌کند (`applyMCP`؛ پایین‌تر را ببینید).
7. ccam نشست را ثبت می‌کند (`session.Register`). یک نشانگر به نام PID خود، که زمان شروع فرایند را در خود دارد، در `~/.local/state/ccam/locks/profile-<name>/<pid>` می‌نویسد، هر چه `XDG_RUNTIME_DIR` بگوید، تا همه shellها نشانگرهای یکسانی ببینند. برای حسابی که تاریخچه گفتگو را به اشتراک می‌گذارد، ccam ابتدا در `history-<group>-<hash>`، که hash آن از پوشه فعلی با پیوندهای نمادین حل‌شده گرفته می‌شود، به دنبال نشستی در حال اجرا می‌گردد و آنجا هم ثبت می‌کند (`historyCheck`).

   `exec` همان PID را نگه می‌دارد، بنابراین نشانگر سپس به خود فرایند Claude Code اشاره می‌کند. ccam به‌جای قفل نگه‌داشته‌شده از یک فایل نشانگر استفاده می‌کند، چون فرایندهایی که Claude Code اجرا می‌کند قفل را به ارث می‌برند و پس از پایان نشست آن را نگه می‌دارند. نشانگری که نوشته نشود اجرا را متوقف نمی‌کند: ccam می‌گوید `ccam: could not mark the session as running (…); starting Claude Code anyway`.

   بررسی‌های دیگر ccam درباره اینکه آیا حسابی در حال اجراست (تغییر نام، حذف، تغییر گروه‌ها، `doctor --fix`، همگام‌سازی MCP) سوابق خود Claude Code را هم در پوشه حساب می‌خوانند (`sessions/<pid>.json`، `daemon.lock`)، پس Claude Code را که بدون ccam شروع شده هم می‌بینند: `claude` معمولی، یک IDE، نشست‌های پس‌زمینه و daemon آن‌ها.
8. ccam عنوان ترمینال را `claude · <name>` می‌گذارد (وقتی خروجی خطای استاندارد ترمینال باشد)، زمان اجرا را بدون `fsync` در `~/.local/state/ccam/last-used/<name>` می‌نویسد و `syscall.Exec` را فراخوانی می‌کند (`Plan.Exec`).

Claude Code با `CLAUDE_CONFIG_DIR` برابر با پوشه حساب و `CCAM_PROFILE` برابر با نام حساب اجرا می‌شود. ccam هر مقدار به ارث رسیده از این متغیرها، و `CCAM_RELAY` به ارث رسیده را، پاک می‌کند. برای پوشه `~/.claude` که با `ccam import` روی macOS ثبت شده، `CLAUDE_SECURESTORAGE_CONFIG_DIR` را هم با مقدار خالی تنظیم می‌کند تا آن حساب و `claude` معمولی از یک ورود Keychain استفاده کنند؛ در Windows وقتی Claude Code ورود را در Credential Manager نگه می‌دارد همین کار را می‌کند.

### چرا exec

ccam فرایند خودش را با `claude` جایگزین می‌کند (`syscall.Exec`). در نتیجه:

- ccam در حافظه باقی نمی‌ماند؛
- کلیدها، سیگنال‌ها و کد خروج متعلق به خود Claude Code هستند؛
- آرگومان‌ها به شکل آرایه می‌مانند، هرگز رشته shell نمی‌شوند.

Windows `exec` ندارد. آنجا ccam فایل `claude.exe` را به‌عنوان فرایند فرزند با همان محیط، آرگومان‌ها و کنسول اجرا می‌کند، منتظرش می‌ماند و با کد خروج آن خارج می‌شود (`Plan.Exec`). `claude.cmd` که npm نصب کرده باشد در `cmd.exe` با مسیر کاملش اجرا می‌شود و هر آرگومان طوری نقل‌قول می‌شود که `cmd.exe` هیچ‌کدام را نحو خودش نخواند (`&`، `%PATH%`)؛ آرگومانی که شکست خط داشته باشد را نمی‌توان این‌طور داد و رد می‌شود، با توصیه به نصب Claude Code بومی.

Ctrl+C مستقیم از کنسول به Claude Code می‌رسد، پس ccam آن را نادیده می‌گیرد. نشانگر نشست PID خود ccam را دارد که دقیقاً به اندازه Claude Code زنده می‌ماند و وقتی Claude Code خارج شود حذف می‌شود. ccam هنگام انتظار کار دیگری انجام نمی‌دهد.

کارهای پس از نشست، مانند تشخیص سرورهایی که با `claude mcp add` اضافه شده‌اند، در اجرای بعدی انجام می‌شود: ccam زمان تغییر `.claude.json` را با `~/.local/state/ccam/mcp/<name>/seen` مقایسه می‌کند و فقط وقتی فایل تغییر کرده باشد آن را تجزیه می‌کند.

### مسیر یک خط فرمان

سه خط فرمان، سه مسیر در `internal/app`:

<!-- motion:launch-route -->

میان‌بر بر اساس `argv[0]` و `ccam <name>` بر اساس آرگومان اول تشخیص داده می‌شود و هر دو به `syscall.Exec` می‌رسند. فقط یک فرمان ccam مانند `ccam usage` به cobra می‌رسد.

### سرورهای MCP مشترک

با راهبرد پیش‌فرض `flag`، ccam سرورهای MCP گروه را به‌صورت یک آرگومان، پیش از آرگومان‌های شما، به Claude Code می‌دهد:

<!-- motion:launch-mcp -->

ccam فقط وقتی فایل `mcp.json` گروه وجود داشته باشد این آرگومان را اضافه می‌کند و آن را درست پس از `claude` قرار می‌دهد (`mcp.FlagArgs`، `Plan.InsertArgs`). گزینه‌های `--mcp-config` یا `--strict-mcp-config` خود شما همچنان روی آن اعمال می‌شوند. این یک آرگومان واحد است، `--mcp-config=<path>`، چون `--mcp-config` در Claude Code چند مقدار می‌پذیرد: اگر به شکل دو آرگومان نوشته شود، پرامپت بعدی را هم به‌عنوان فایل پیکربندی می‌خواند.

ccam تعیین می‌کند که این گزینه کجا نمی‌تواند برود (`claudeVerb`). فرمان فرعی Claude Code را همان‌طور که Claude Code از آرگومان‌ها می‌خواند می‌خواند و از گزینه‌ها همراه با مقدارهایشان می‌گذرد. `remote-control` یا `rc` این گزینه را نمی‌گیرد (پایین‌تر)، و نه یک فعل نشست پس‌زمینه که اول آمده باشد (`attach`، `logs`، `stop`، `kill`، `respawn`، `rm`). `ccam work mcp list`، `get` و `remove` هم آن را نمی‌گیرند، چون Claude Code آنجا ردش می‌کند؛ به جای آن ccam در خطای استاندارد می‌گوید که سرورهای گروه در نشست‌های حساب هم بارگذاری می‌شوند:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

هر فرمان فرعی دیگر و هر نشست این گزینه را نگه می‌دارد.

با `mcp_strategy = "sync"` در `config.toml`، ccam به‌جای آن سرورها را در `.claude.json` حساب می‌نویسد و اگر حساب از قبل در حال اجرا باشد این کار را انجام نمی‌دهد. وقتی نه `mcp.json` گروه و نه `.claude.json` از آخرین بررسی تغییر نکرده باشند، اجرا هم از این کار می‌گذرد. اگر حساب سرور خودش را با همان نام یک سرور گروه داشته باشد، همگام‌سازی سرور حساب را نگه می‌دارد و یک بار این را می‌گوید. همگام‌سازی همچنین نسخه‌ای را که از `.claude.json` حذف شده برمی‌گرداند (`put back group default's MCP servers …`)، و فقط نسخه‌هایی را تغییر می‌دهد یا حذف می‌کند که خودش نوشته و از آن پس کسی تغییرشان نداده است.

### Remote Control

`ccam work rc` (یا `ccam work remote-control`) قابلیت Remote Control در Claude Code را با حساب `work` اجرا می‌کند. `claude remote-control` گزینه `--mcp-config` را پیش از فعل خود رد می‌کند و پس از آن هم هیچ‌کدام را نمی‌پذیرد، و نشست‌های آن سرورهای MCP را فقط از `.claude.json` حساب می‌خوانند. بنابراین راهبرد flag اینجا طور دیگری کار می‌کند:

<!-- motion:launch-rc -->

1. ccam فعل `remote-control` یا `rc` را به‌عنوان فعل Claude Code پیدا می‌کند، حتی پس از گزینه‌ها و مقدارهایشان (`ccam work --model opus rc`)، و هیچ گزینه‌ای اضافه نمی‌کند (`claudeVerb`، `mcp.FlagArgs`).
2. ccam زمان تغییر `mcp.json` گروه را با `sync-stamp` حساب مقایسه می‌کند و بررسی می‌کند که هر سرور گروه هنوز در `.claude.json` باشد (`mcp.Synced`). اگر `.claude.json` از قبل سرورهای فعلی گروه را داشته باشد، یا گروه هیچ سروری نداشته باشد، Remote Control بدون پرسش شروع می‌شود.
3. در غیر این صورت ccam سرورهایی را که نشست‌ها از دست می‌دهند نام می‌برد و می‌پرسد:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. با `y`، ccam آن‌ها را همان‌طور کپی می‌کند که راهبرد `sync` کپی می‌کند (`mcp.Sync`): زیر قفل Claude Code، پس از گرفتن نسخه پشتیبان و با نگه داشتن سرورهای خود حساب. با `N`، یا بدون ترمینال، Remote Control بدون آن‌ها شروع می‌شود و ccam به `mcp_strategy = "sync"` اشاره می‌کند که آن‌ها را پیش از هر اجرا کپی می‌کند.

ccam در هر اجرای Remote Control می‌پرسد تا سرورها کپی شوند. پس از کپی، فقط وقتی دوباره می‌پرسد که `mcp.json` گروه تغییر کند یا یک سرور کپی‌شده از `.claude.json` ناپدید شود. در اجراهای عادی بعدی، کپی موجود در `.claude.json` بر `--mcp-config` برتری دارد، بنابراین `ccam mcp add`، `remove` و `import` کپی‌هایی را که حساب‌های دیگر نگه می‌دارند به‌روز می‌کنند ([سرورهای MCP مشترک](/features/mcp)).

### بودجه ۱۰ میلی‌ثانیه

بودجه مسیر اجرا ۱۰ میلی‌ثانیه در صدک ۹۵ است؛ در عمل ccam حدود ۴ میلی‌ثانیه پیش از شروع Claude Code اضافه می‌کند. مسیر اجرا هرگز درخت فرمان‌های cobra را نمی‌سازد، Bubble Tea را راه نمی‌اندازد و اتصال شبکه باز نمی‌کند، و اگر مجموع زمان توابع `init()` بسته‌ها بیش از ۵ میلی‌ثانیه شود، `TestPackageInitIsCheap` شکست می‌خورد. `scripts/perf-check.sh` زمان `ccam --dry-run work`، `ccam work` و یک میان‌بر را با `hyperfine` اندازه می‌گیرد و بالای ۱۰ میلی‌ثانیه شکست می‌خورد، و `BenchmarkLaunchPath` در `internal/launcher` گام‌ها را در یک حلقه اجرا می‌کند. برای دیدن زمان هر گام، `CCAM_TRACE=1` را تنظیم کنید:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## ایمنی و حالت‌های خاص

- نام یا شماره ناشناخته با `no account "<name>"` متوقف می‌شود و شما را به `ccam list` ارجاع می‌دهد. اگر هیچ حسابی وجود نداشته باشد، ccam از شما می‌خواهد با `ccam add <name>` یکی اضافه کنید.
- اگر یک مسیر مشترک به فایل عادی تبدیل شده یا به جای اشتباهی اشاره کند، ccam با `stopped because of share link problems:` متوقف می‌شود، مشکلات را فهرست می‌کند و از شما می‌خواهد `ccam doctor --fix` را اجرا کنید. ccam در مسیر اجرا داده‌ای را جابه‌جا نمی‌کند.
- اگر نشست Claude Code دیگری که تاریخچه گفتگو را با این حساب به اشتراک می‌گذارد (همان گروه) در همان پوشه در حال اجرا باشد، ccam هشدار می‌دهد که تاریخچه ممکن است در هم آمیخته شود و می‌پرسد `Run anyway? [y/N]`. پوشه با پیوندهای نمادین حل‌شده مقایسه می‌شود. بدون ترمینال (ورودی استاندارد از `/dev/null` هم ترمینال به حساب نمی‌آید)، یا از یک کار در پس‌زمینه، هشدار را نشان می‌دهد و ادامه می‌دهد. یک اجرای تازه `-p` گفتگوی خودش را شروع می‌کند و هشداری نمی‌گیرد؛ `--version`، `--help` و فرمان‌های فرعی که گفتگویی باز نمی‌کنند، مانند `mcp`، `auth` یا `plugin`، اصلاً بررسی نمی‌شوند.
- `--dry-run` از همگام‌سازی MCP، اعلان سرور جدید، یادداشت `mcp list`، بررسی تاریخچه، نشانگر نشست، عنوان و زمان آخرین استفاده می‌گذرد و exec نمی‌کند. با این حال پیوندهای اشتراک را بررسی می‌کند و پیوندی را که وجود نداشته باشد، مانند یک اجرای عادی، می‌سازد.
- اجرای `ccam work rc` وقتی `work` از قبل در حال اجراست به `.claude.json` دست نمی‌زند، این را اعلام می‌کند و Remote Control را بدون سرورهای گروه شروع می‌کند؛ پس از پایان نشست دوباره آن را اجرا کنید تا سرورها کپی شوند. با `--dry-run` پیام را چاپ می‌کند اما نه می‌پرسد و نه کپی می‌کند.
- با `mcp_strategy = "sync"`، اگر همگام‌سازی ناموفق باشد `ccam: MCP sync failed (launching anyway): …` چاپ می‌شود و Claude Code با سرورهایی که از قبل داشت شروع می‌شود.

## مطالب مرتبط

- [فرمان‌های میان‌بر](/features/shortcuts)
- [استفاده از Remote Control با ccam](/guides/remote-control)
- [مدیریت حساب‌ها](/features/accounts)
- [مرجع فرمان‌ها](/docs/commands)
- [اشتراک‌گذاری سرورهای MCP، تنظیمات و تاریخچه](/docs/sharing)
- [توضیح CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
