---
title: مرجع فرمان‌ها
metaTitle: مرجع فرمان‌های ccam
description: همه فرمان‌ها و گزینه‌های ccam، از ccam add و ccam <name> تا ccam usage، ccam next و ccam doctor.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: 28e1ad1dcbad
---

برای اینکه ببینید یک فرمان چه کاری انجام می‌دهد بدون آنکه واقعاً انجامش دهد، `--dry-run` را پیش از آن بگذارید، برای مثال `ccam --dry-run work`. این گزینه هنگام اجرای Claude Code (`ccam <name>`، `ccam run`، `ccam next`، `ccam attach`، `ccam ccusage`) و با `group`، `migrate`، `folder`، `shortcut`، `pick`، `update`، `uninstall` و `doctor --fix` کار می‌کند. فرمان‌هایی که فقط می‌خوانند، مانند `list`، `usage`، `share` و `ext`، آن را می‌پذیرند و در هر حال چیزی را تغییر نمی‌دهند. بقیه فرمان‌ها، از جمله `group rename`، به‌جای اجرا شدن آن را رد می‌کنند.

<!-- commands -->

## اجرای Claude Code

`ccam <name|number> [arguments]` اجرای Claude Code با آن حساب را آغاز می‌کند. ccam آرگومان‌های پس از نام را نمی‌خواند؛ آن‌ها را به همان ترتیب به `claude` می‌دهد، و ورودی استاندارد، کد خروج و سیگنال‌ها دقیقاً مانند خود `claude` رفتار می‌کنند. گذاشتن `--` پس از نام اختیاری است.

`ccam` با آرگومان ولی بدون نام حساب، مانند `ccam -p "summarise"`، از حسابی که برای پوشه فعلی تنظیم شده و در غیر این صورت از حساب پیش‌فرض استفاده می‌کند. آن‌ها را با `ccam folder <name> [folder]` و `ccam use <name>` تعیین کنید.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` پیش از نام حساب برای اجراهای `-p`/`--print` است. وقتی حساب به سقف مصرفش برسد، ccam همان فرمان را دوباره با حساب واردشده‌ای که بیشترین ظرفیت را دارد اجرا می‌کند و فقط خروجی اجرایی را که به پایان رسیده چاپ می‌کند. [حساب‌ها بر اساس پوشه](/features/folder) و [اجرای دوباره هنگام رسیدن به سقف](/features/next#run-again-at-the-limit) را ببینید.

**صندلی‌های Team و Enterprise.** انتخاب خودکار حساب (`ccam next`، `--fallback`، `ccam pick`، کلید `n` در داشبورد و، وقتی آن را روشن کنید، [`auto_switch`](/features/next#auto-switch)) میان صندلی‌های همان سازمان Claude Team یا Enterprise هم جابه‌جا می‌شود، که ممکن است شرایط Anthropic را نقض کند. پیش از آنکه به آن تکیه کنید، [پرسش‌های متداول](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) را بخوانید و از سازمان خود و Anthropic بپرسید.

میان‌بری مانند `ccwork` حساب خود را به همین شکل اجرا می‌کند: `ccwork --resume` همان کاری را می‌کند که `ccam work --resume` انجام می‌دهد. بخش [میان‌برها](/docs/commands#shortcuts) را در ادامه ببینید.

## افزودن و حذف حساب‌ها

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` بدون نام پوشه، به دنبال `~/.claude` و `~/.claude-*` می‌گردد و درباره هر کدام از شما می‌پرسد؛ پوشه‌ای که ناموفق باشد به‌صورت `Skipped <path>: <reason>` گزارش می‌شود و بقیه ثبت می‌شوند. اگر نام یک پوشه از قبل گرفته شده باشد، آن پوشه را جداگانه با `ccam import <folder> --name <name>` وارد کنید. پوشه‌های واردشده در همان جایی که هستند ثبت می‌شوند؛ ccam هرگز یک ورود را جابه‌جا یا کپی نمی‌کند.

`ccam login --expired` حساب‌هایی را هم که Claude Code ورودشان را پاک کرده است وارد می‌کند. `ccam remove` یک حساب مدیریت‌شده (حسابی که پوشه‌اش را ccam ساخته است، نه حسابی که با `ccam import` ثبت شده) را پیش از آنکه پوشه‌اش به زباله‌دان برود از سیستم خارج می‌کند، مگر اینکه حساب دیگری همان ورود را داشته باشد. حذف یک حساب واردشده آن را از گروه‌هایش بیرون می‌برد و از هر مورد مشترک یک کپی در پوشه‌اش به جا می‌گذارد. در macOS، `ccam rename` نام پوشه را نگه می‌دارد، چون Keychain ورود را زیر مسیر پوشه نگه می‌دارد.

`ccam add --api` حسابی اضافه می‌کند که به جای ورود به Claude از کلید API استفاده می‌کند؛ ccam کلید را در هر اجرا می‌خواند و هرگز ذخیره‌اش نمی‌کند. [حساب‌های کلید API](/features/api) را ببینید.

تا وقتی Claude Code با یک حساب در حال اجراست، `ccam rename`، `ccam remove` و فرمان‌های `ccam group` آن حساب را تغییر نمی‌دهند، Claude Code به هر شکلی که شروع شده باشد (از جمله نشست‌های پس‌زمینه). ابتدا نشست را ببندید. همچنین تا وقتی یک تغییر گروه نیمه‌تمام مانده است، سر باز می‌زنند: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` هر حساب را همراه با شماره جایگاه، میان‌بر، ایمیل و پوشه آن چاپ می‌کند. آرایه `groups` آن گروه‌های حساب را فهرست می‌کند و شیء `share` آن هر مورد (`mcp`، `settings`، `memory`، `extensions`، `history`) را به گروهی که از آن می‌آید، یا به `local`، نگاشت می‌کند.

## میان‌برها

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

میان‌بر فرمانی است که نامش را خودتان انتخاب می‌کنید و Claude Code را با یک حساب اجرا می‌کند. مانند `ccam work`، همه آرگومان‌ها را بدون تغییر به Claude Code می‌دهد، بنابراین `ccwork --resume` و `ccwork -p "hi"` همان‌طور که انتظار دارید کار می‌کنند. میان‌بر یک پیوند به ccam در `~/.local/bin` است.

میان‌برها اختیاری‌اند. دستیار `ccam add` می‌پرسد که میان‌بری ساخته شود یا نه، `ccam add <name> --shortcut <command>` بدون پرسش یکی می‌سازد و `ccam list` میان‌بر هر حساب را نشان می‌دهد، یا `-` اگر میان‌بری نداشته باشد.

نام‌ها از نویسه‌های a-z، 0-9، `.`، `-` و `_` تشکیل می‌شوند. نام نمی‌تواند یکی از فرمان‌های ccam، میان‌بر حساب دیگری یا نام برنامه دیگری در `PATH` شما باشد، و نه فایلی در `~/.local/bin` که ccam آن را نساخته است، مانند کپی یا پیوند خود شما به ccam. میان‌بری که دستی در `config.toml` ویرایش شده و این قواعد را می‌شکند، `claude` یا `ccam` است، یا به دو حساب داده شده است، باعث می‌شود `config.toml` بارگذاری نشود و نام حساب گفته شود.

`ccam rename` نام میان‌بر را نگه می‌دارد و `ccam remove` میان‌بر را همراه با حساب حذف می‌کند. `ccam shims install` پیوندهای همه حساب‌ها را دوباره می‌سازد، برای نمونه پس از اینکه ccam را جابه‌جا کرده‌اید؛ فایلی که سر راه یک میان‌بر باشد مانع بقیه نمی‌شود.

## اشتراک‌گذاری

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

یک گروه مواردی را بین حساب‌هایش به اشتراک می‌گذارد؛ یک حساب می‌تواند در چند گروه باشد، اما هر مورد فقط از یک گروه می‌آید. [اشتراک‌گذاری با گروه‌ها](/features/share) را ببینید.

## کپی کردن یک پیکربندی موجود

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` فایل‌های `settings.json` و `CLAUDE.md`، افزونه‌ها و سرورهای MCP کاربر را از یک پوشه Claude Code در یک گروه کپی می‌کند تا همه حساب‌هایی که این موارد را به اشتراک می‌گذارند آن‌ها را داشته باشند. پوشه فقط خوانده می‌شود.

- `--settings`، `--memory`، `--extensions`، `--mcp` و `--history` فقط مواردی را که نام می‌برید کپی می‌کنند. بدون آن‌ها، در ترمینال یک فهرست انتخاب به شما امکان انتخاب می‌دهد که در آن همه چیز به جز تاریخچه گفتگو علامت خورده است؛ در جاهای دیگر همه چیز به جز تاریخچه گفتگو کپی می‌شود.
- `--mcp` وقتی `~/.claude` خالی است یا وجود ندارد اما `~/.claude.json` سرور دارد هم کار می‌کند.
- `-g` یا `--group` گروه را انتخاب می‌کند.
- `-y` یا `--yes` بدون پرسش کپی می‌کند و وقتی فایلی متفاوت باشد، نسخه گروه را نگه می‌دارد.

ورودهای MCP کپی نمی‌شوند، چون یک ورود کپی‌شده بقیه دارندگان را از سیستم خارج می‌کند؛ برنامه حساب‌هایی را که باید با `/mcp` وارد شوند نام می‌برد. [کپی کردن یک پیکربندی موجود](/docs/sharing#copy-an-existing-setup) را ببینید.

## سرورهای MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` سرورهایی را که فقط آن حساب دارد به گروه آن منتقل می‌کند، از جمله آن‌هایی که فقط برای یک پوشه دارد (`claude mcp add` ساده)، که از آن پس در همه پوشه‌ها بارگذاری می‌شوند. سروری که گروه از نامش برای تعریف دیگری استفاده می‌کند رد می‌شود؛ برای جایگزین کردن نسخه گروه `--force` را اضافه کنید، ccam پیش از آن از نسخه گروه پشتیبان می‌گیرد.

## مصرف و جابه‌جایی

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam sends no token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --continue --brief    # a new conversation from a summary of the last one
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agent "…"                  # from Claude Code: a subagent that runs as another account
ccam agent --resume <id> "…"    # follow that agent up
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
ccam ccusage claude daily       # Claude Code's tokens and costs, every account (ccusage)
ccam ccusage -a work blocks --active   # one account's current 5-hour block
```

- `next --auto`، `--fallback` و `ccam pick` به یک شکل انتخاب می‌کنند: ابتدا حساب‌های زیر `auto_max_usage` (`config.toml`، پیش‌فرض ۱۰۰)، بر اساس `auto_strategy` (`headroom`، `expiring` یا `balance`)، با در نظر گرفتن پنجره هفتگی `--model`. هرگز حسابی را که به ۱۰۰٪ رسیده است (پنجره‌ای که Claude Code پر گزارش کرده تا بازنشانی‌اش ۱۰۰٪ حساب می‌شود) یا ورودی با کلید API را انتخاب نمی‌کنند و در سازمان حسابی که از آن شروع می‌کنند می‌مانند (`--cross-org` یا `cross_org = true` برای جست‌وجوی فراتر از آن).
- `--wait` مقدارش را پس از `=` می‌گیرد (`--wait=2h`)؛ اگر با فاصله نوشته شود، با یک راهنمایی رد می‌شود. وقتی هیچ حسابی را نتوان انتخاب کرد، با کد ۷۵ خارج می‌شوند، یا با `--wait` منتظر می‌مانند تا نخستین حساب بازنشانی شود.
- `ccam next --continue` یک گفتگو را حتی بین حساب‌هایی که تاریخچه مشترک ندارند هم منتقل می‌کند: در این حالت Claude Code آن را در حساب جدید کپی می‌کند. با `--brief`، حساب جدید به‌جای آن گفتگوی تازه‌ای را آغاز می‌کند، از خلاصه‌ای که ccam بدون یک مدل از گفتگوی قبلی می‌نویسد ([`--brief`](/features/next#brief)).
- برخلاف `ccam <name>`، `ccam next` گزینه‌های خودش را می‌خواند، پس گزینه‌های Claude Code پس از `--` می‌آیند؛ بدون آن ccam با `unknown flag` متوقف می‌شود.
- `ccam agent` یک کار را با `claude -p` در این پوشه با حسابی غیر از حسابی که از آن فراخوانده می‌شود اجرا می‌کند، با حالت مجوز و سطح تلاش نشست فراخواننده، و آخرین پاسخ Claude را چاپ می‌کند. در سقف آن حساب، گفتگو با حساب دیگری ادامه پیدا می‌کند ([subagentها در حساب‌های دیگر](/features/agent)).
- درون Claude Code، فرمان `/ccam:switch [account]` پلاگین یک گفتگو را به حساب دیگری می‌سپارد و `ccam hook` همان چیزی است که هوک‌های آن اجرا می‌کنند.
- با `auto_switch = true` در `config.toml` (Linux و macOS؛ به‌طور پیش‌فرض خاموش)، هوک سقف مصرف یک نشست تعاملی را هم که ccam آغاز کرده پایان می‌دهد و آن را در حساب دیگری که ظرفیت دارد ادامه می‌دهد، اما فقط در حسابی از همان نوع: صندلی دیگری از همان سازمان، یا طرح شخصی دیگری ([`auto_switch`](/features/next#auto-switch)).

[جابه‌جایی به حساب بعدی](/features/next)، [انتخاب حساب برای اسکریپت‌ها](/features/pick)، [subagentها در حساب‌های دیگر](/features/agent) و [نشست‌های پس‌زمینه همه حساب‌ها](/features/agents) را ببینید.

## نگهداری

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam uninstall       # remove every account, ccam's folders and ccam itself
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` پیش از جایگزین کردن ccam، checksum SHA-256 نسخه را بررسی می‌کند و `--version v0.2.0` یک نسخه مشخص را نصب می‌کند. نسخه‌ای از ccam که از بسته `.deb` نصب شده باشد را باید با apt به‌روز کنید. [به‌روزرسانی](/docs/install#update) را ببینید.

`ccam uninstall` فهرست چیزهایی را که حذف می‌کند نشان می‌دهد و ابتدا می‌پرسد: همه حساب‌ها همان‌طور که `ccam remove` حذف می‌کند، فرمان‌های میان‌بر، پوشه‌های ccam و فرمان `ccam`. پوشه واردشده در جای خود باقی می‌ماند. [حذف نصب](/docs/install#uninstall) را ببینید.

## فایل‌ها

| مسیر | محتوا |
| --- | --- |
| `~/.config/ccam/config.toml` | حساب‌ها، گروه‌ها و تنظیمات (`theme`، `usage_source`، `relay`، `auto_switch`، `cross_org` …). تنظیمات در بالای فایل می‌آیند، بالای نخستین `[[group]]`، `[[profile]]` یا `[[folder]]`: تنظیمی که زیر یک جدول باشد به همان جدول تعلق دارد، پس ccam در هر اجرا و در `ccam doctor` درباره‌اش هشدار می‌دهد و فرمان‌هایی که پیکربندی را ذخیره می‌کنند تا جابه‌جا شدنش سر باز می‌زنند. وقتی ccam فایل را ذخیره می‌کند، توضیحات نگه داشته نمی‌شوند |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` هر حساب |
| `~/.local/share/ccam/groups/<group>/` | تنظیمات، سرورهای MCP، افزونه‌ها و تاریخچه مشترک |
| `~/.local/state/ccam/` | حافظه نهان مصرف، ثبت‌های سقف مصرف، قفل‌ها، یادداشت‌ها و وضعیت‌های دیگر (فهرست زیر جدول را ببینید) |
| `~/.local/bin/<shortcut>` | فرمان‌های میان‌بر مانند `ccwork`، که پیوندهایی به ccam هستند |

پوشه وضعیت این‌ها را نگه می‌دارد:

- حافظه نهان مصرف و ثبت‌های سقف مصرف (`usage/<name>.json`، `usage/<name>.limit.json`)؛
- زمان‌های آخرین استفاده؛
- نشانگرهای نشست و قفل‌ها (`locks/`)؛
- میان‌برهایی که ccam ساخته است (`shortcuts.json`)؛
- تغییر گروهی که تمام نشده است (`group-change.json`)؛
- یادداشت‌هایی که `/ccam:switch` و هوک سقف مصرف به جا می‌گذارند (`handoff/<account>/<session>.json`)؛
- زمانی که `ccam pick` یا `auto_switch` آخرین بار هر حساب را انتخاب کرده است (`picks/<name>`؛ هر انتخاب ۲ دقیقه حساب می‌شود)؛
- خلاصه‌های `--brief` که ۷ روز نگه داشته می‌شوند (`brief/<session>.md`).

تنظیم `XDG_CONFIG_HOME`، `XDG_DATA_HOME` یا `XDG_STATE_HOME` محل این پوشه‌ها را تغییر می‌دهد. در Windows، `~` همان `%USERPROFILE%` است و میان‌برها `<shortcut>.exe` هستند.
