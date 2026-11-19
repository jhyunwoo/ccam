---
title: مدیریت حساب‌ها
description: افزودن، ورود، فهرست، تغییر نام، حذف و وارد کردن حساب‌های Claude Code با ccam add، ccam login، ccam list، ccam rename، ccam remove و ccam import.
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

هر حساب ccam یک پوشه پیکربندی جداگانه Claude Code با ورود مخصوص خودش است. با این فرمان‌ها حساب اضافه می‌کنید، با آن‌ها وارد می‌شوید و آن‌ها را فهرست، تغییر نام و حذف می‌کنید. این فرمان‌ها پوشه‌هایی را هم که از قبل دارید ثبت می‌کنند.

## نحوه استفاده

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

هر جا فرمانی `<name|number>` بپذیرد، می‌توانید نام حساب یا شماره جایگاه آن (ستون `#` در `ccam list`) را بدهید.

`ccam use` فقط حساب پیش‌فرض را تغییر می‌دهد (`default` در `config.toml`).

`ccam add` در ترمینال و بدون گزینه یک دستیار باز می‌کند. دستیار این مراحل را طی می‌کند:

1. نام حساب، اگر آن را نداده باشید؛
2. روش ورود حساب: ورود به Claude همین حالا، کلید API یا ورود به Claude بعداً؛
3. گروه‌هایی که حساب به آن‌ها می‌پیوندد، با `default` علامت‌خورده (اگر هنوز گروهی وجود نداشته باشد، نشان می‌دهد گروه `default` چه چیزهایی را به اشتراک خواهد گذاشت)؛
4. `Create a shortcut command for work? [y/N]`؛
5. خلاصه‌ای برای تأیید.

`Esc` یک مرحله به عقب برمی‌گردد. هر یک از `--group`، `--no-group`، `--shortcut` یا `--no-login`، یا اجرا بدون ترمینال (برای مثال از یک pipe)، دستیار را کنار می‌گذارد. در این حالت نام الزامی است و میان‌بری ساخته نمی‌شود مگر اینکه `--shortcut` نامی برای آن بدهد.

بدون `--group`، حساب جدید به گروه `default` می‌پیوندد. نخستین حساب این گروه را می‌سازد که سرورهای MCP، تنظیمات، `CLAUDE.md` و افزونه‌ها را به اشتراک می‌گذارد؛ [اشتراک‌گذاری با گروه‌ها](/features/share) را ببینید. در Windows فایل‌ها را فقط در Developer Mode می‌توان پیوند داد. بدون آن، نخستین حساب `default` را فقط با `mcp` و `extensions` می‌سازد و این را می‌گوید (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`)، و ccam گزینه `--group` را با گروهی که تنظیمات، حافظه یا تاریخچه را به اشتراک می‌گذارد رد می‌کند.

`ccam list` جایگاه، نام، میان‌بر، ایمیل، گروه‌ها و آخرین استفاده را نشان می‌دهد و پیش از حساب پیش‌فرض `*` می‌گذارد. ایمیل همان است که پوشه اکنون با آن وارد شده است (از `.claude.json` آن)؛ وقتی با ایمیلی که ccam ثبت کرده فرق داشته باشد، ایمیل ثبت‌شده به‌صورت `(was …)` پس از آن می‌آید:

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

با `--json`، هر حساب یک شیء با `slot`، `shortcut`، `name`، `email` (ایمیل ورود فعلی، مانند جدول)، `default`، `dir`، `last_used` (که اگر حساب هرگز اجرا نشده باشد حذف می‌شود)، `groups` و `share` است. `share` هر مورد (`mcp`، `settings`، `memory`، `extensions`، `history`) را به گروهی که از آن می‌آید، یا به `"local"`، نگاشت می‌کند.

## گزینه‌ها

| گزینه | کاری که انجام می‌دهد | پیش‌فرض |
| --- | --- | --- |
| `ccam add --no-login` | حساب را بدون ورود می‌سازد؛ بعداً با `ccam login <name>` وارد شوید | خاموش |
| `ccam add --group <group>` | به‌جای `default` به این گروه می‌پیوندد. قابل تکرار؛ اگر `default` وجود نداشته باشد ساخته می‌شود | `default` |
| `ccam add --no-group` | به هیچ گروهی نمی‌پیوندد: حساب چیزی را به اشتراک نمی‌گذارد | خاموش |
| `ccam add --shortcut <command>` | این فرمان میان‌بر را هم در `~/.local/bin` می‌سازد | هیچ |
| `ccam add --api --key-from <source>` | حسابی که به جای ورود از کلید API استفاده می‌کند؛ [حساب‌های کلید API](/features/api) را ببینید | خاموش |
| `ccam list --json` | حساب‌ها را به‌صورت JSON چاپ می‌کند | خاموش |
| `ccam remove -y`، `--yes` | بدون پرسش حذف می‌کند | خاموش |
| `ccam logout -y`، `--yes` | بدون پرسش از سیستم خارج می‌کند | خاموش |
| `ccam import --name <name>` | نام حساب برای پوشه‌ای که می‌دهید؛ بدون پوشه رد می‌شود (`--name needs a folder: ccam import <folder> --name ops`) | برگرفته از نام پوشه |
| `ccam import -y`، `--yes` | همه پوشه‌های پیداشده را بدون پرسش وارد می‌کند | خاموش |
| `--dry-run` | `ccam list` و `ccam completion` که چیزی را تغییر نمی‌دهند آن را می‌پذیرند. `add`، `login`، `logout`، `use`، `rename`، `remove` و `import` آن را با پیام `` --dry-run is not supported by `ccam add`; nothing was changed `` (با نام همان فرمان) رد می‌کنند | خاموش |

`ccam login` (جز `--expired`)، `ccam use` و `ccam rename` گزینه‌ای از خود ندارند. `ccam completion` مقدار `bash`، `zsh` یا `powershell` را می‌پذیرد؛ هر چیز دیگری با `only bash, zsh and powershell are supported` ناموفق می‌شود.

نام حساب‌ها ۱ تا ۳۲ نویسه از `a-z`، `0-9`، `-` و `_` است و با یک حرف یا رقم شروع می‌شود. نام نمی‌تواند فقط از رقم تشکیل شود، چون `ccam 2` یعنی جایگاه ۲. نام همچنین نمی‌تواند یکی از فرمان‌های ccam باشد (`add`، `list`، `ls`، `run`، `use`، `next`، `usage`، `u`، `share`، `mcp`، `login`، `logout`، `rename`، `remove`، `rm`، `import`، `migrate`، `doctor`، `shims`، `shortcut`، `completion`، `help`، `version`، `statusline`، `update`، `uninstall`، `folder`، `pick`، `ext`، `group`، `groups`، `hook`، `agent`، `agents`، `attach`، `ccusage`)، تا `ccam <name>` هرگز مبهم نباشد. این فهرست `config.Reserved` در `internal/config` است.

## نحوه کار

افزودن یک حساب به این ترتیب پوشه‌اش را می‌سازد، به گروه‌هایش می‌پیوندد، پیکربندی را ذخیره می‌کند، وارد می‌شود و میان‌بر را می‌سازد. وقتی `ccam add work --shortcut ccwork` را اجرا می‌کنید، ccam این مراحل را طی می‌کند (`addProfile` در `internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. ccam پیش از ساختن هر چیزی نام (`config.ValidateName`) و میان‌بر (`shim.Check`) را بررسی می‌کند.
2. ccam پوشه `~/.local/share/ccam/profiles/work/` را با مجوز `0700` می‌سازد و کوچک‌ترین شماره جایگاه آزاد را به حساب می‌دهد، بنابراین شماره یک حساب حذف‌شده دوباره استفاده می‌شود (`profile.Add`، `Config.NextSlot`). نخستین حساب، حساب پیش‌فرض می‌شود.
3. حساب به گروه‌هایش می‌پیوندد (`default`، مگر اینکه `--group` یا `--no-group` چیز دیگری بگوید؛ نخستین حساب `default` را می‌سازد) و ccam پوشه جدید و خالی را به آن‌ها پیوند می‌دهد، برای مثال `settings.json` را به `groups/default/settings.json` (`applyShare`، با `share.Build` و سپس `Apply`). اگر این کار ناموفق باشد، پوشه جدید حذف می‌شود، چون پوشه باقی‌مانده هر تلاش دوباره را با `folder already exists` متوقف می‌کرد.
4. حساب از طریق `fsutil.AtomicWrite` در `~/.config/ccam/config.toml` ذخیره می‌شود: یک فایل موقت، `fsync` و سپس `rename`، بنابراین یک خرابی هرگز پیکربندی نیمه‌نوشته بر جا نمی‌گذارد.
5. ccam فرمان `claude auth login` را با `CLAUDE_CONFIG_DIR` تنظیم‌شده روی پوشه حساب، و با `--email` همراه با آخرین ایمیل حساب، اجرا می‌کند تا صفحه ورود از پیش پر شده باشد (`runLogin`). اگر ورود ناموفق باشد، حساب باقی می‌ماند و ccam خطا را چاپ می‌کند و می‌گوید بعداً `ccam login work` را اجرا کنید.
6. ccam مقدار `oauthAccount` را از `profiles/work/.claude.json` می‌خواند (`profile.ReadIdentity`): ایمیل، شناسه‌های حساب و سازمان، نام سازمان و طرح اشتراک (`claude_max`، `claude_team` …). این فایل ممکن است چند مگابایت باشد، بنابراین ccam به‌جای تجزیه کل فایل، فقط به دنبال همان یک کلید می‌گردد (`TopLevelValue`). ccam این جزئیات را در `[profile.login]` در `config.toml` نگه می‌دارد.
7. ccam پیوند `~/.local/bin/ccwork` را به فایل اجرایی ccam می‌سازد که هنگام شروع، `argv[0]` را در `config.toml` جست‌وجو می‌کند، و پیوند را در `~/.local/state/ccam/shortcuts.json` ثبت می‌کند (`shim.Install`). اگر نتوان پیوند را ساخت، حساب همچنان اضافه می‌شود، ccam دلیلش را می‌گوید (`Could not create shortcuts: …`) و خلاصه `shortcut: ccwork (not created; see above)` را نشان می‌دهد.

ورود ممکن است چند دقیقه طول بکشد و در این مدت فرایند ccam دیگری ممکن است پیکربندی را ذخیره کند. به همین دلیل فرمان‌هایی که منتظر می‌مانند `config.toml` را دوباره می‌خوانند، همان یک تغییر (در اینجا ایمیل) را اعمال می‌کنند و ذخیره می‌کنند (`env.update`). `ccam remove` هم به همین دلیل پس از پرسش خود دوباره می‌خواند.

`ccam login` مراحل ۵ و ۶ را برای یک حساب موجود اجرا می‌کند. Claude Code ایمیل را فقط به‌عنوان راهنمایی برای صفحه ورود می‌گیرد، پس ممکن است به حساب Claude دیگری جز حساب قبلی وارد شوید؛ ccam ورود جدید را با `[profile.login]` مقایسه می‌کند و این را می‌گوید (`ccam: this login is other@company.com, but work was dev@company.com before.`). ورود جدید را نگه می‌دارد و می‌گوید چگونه آن را برگردانید: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

ورود به حساب Claude دیگر [ثبت سقف مصرف](/features/usage#the-limit-record) حساب را هم کنار می‌گذارد، چون آن سقف‌ها مال حساب دیگر بودند. وقتی حساب ccam دیگری به همان حساب Claude وارد شده باشد، ccam هشدار می‌دهد که آن دو سقف‌های مصرف مشترکی دارند.

`ccam login --expired` یکی پس از دیگری به همه حساب‌هایی وارد می‌شود که هرگز وارد نشده‌اند یا ورودشان منقضی شده است: توکن تمدیدشان منقضی شده، یا Claude Code پس از آنکه سرور ورود را رد کرد آن را پاک کرده است. پیش از باز شدن نخستین مرورگر، بررسی می‌کند که بتوان `config.toml` را ذخیره کرد.

`ccam logout` فرمان `claude auth logout` را در حساب اجرا می‌کند و ایمیل ذخیره‌شده را پاک می‌کند؛ سابقه ورود باقی می‌ماند، پس ورود بعدی با آن مقایسه می‌شود. Claude Code ورود را باطل می‌کند و آنچه را برای حساب ذخیره کرده است پاک می‌کند، از جمله ورودهای سرورهای MCP. پس ccam در این موارد ابتدا هشدار می‌دهد و می‌پرسد (`--yes` از پرسش می‌گذرد):

- ورود، ورودهای MCP را در خود دارد؛
- حساب در حال اجراست (نشست‌هایش هم از سیستم خارج می‌شوند)؛
- پوشه حساب دیگری همان ورود را دارد (یک پوشه کپی‌شده) که خروج آن را هم باطل می‌کند.

برای مثال:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

برای ورود با حساب Claude دیگر، `ccam login <name>` به‌تنهایی کافی است و ورودهای MCP را نگه می‌دارد.

وقتی حساب‌ها به بیش از یک سازمان تعلق داشته باشند، `ccam list` یک ستون `Organization` اضافه می‌کند (طرح‌های شخصی `personal` نشان می‌دهند) و `ccam list --json` برای هر حساب `plan` و `org` دارد.

### ورود روی یک ماشین راه دور

از طریق SSH، بخش مرورگر `claude auth login` ممکن است کار نکند. Claude Code منتظر ورود در یک سرور callback روی ماشین راه دور (`http://localhost:<port>/callback`) می‌ماند، اما مرورگری که ورود را تمام می‌کند روی رایانه شما اجرا می‌شود، جایی که آن نشانی به هیچ جا نمی‌رسد و مرورگر خطایی مانند `ERR_CONNECTION_REFUSED` نشان می‌دهد. نشانی همچنان کد یک‌بارمصرف را در خود دارد، پس ccam می‌تواند آن را برساند:

<!-- motion:accounts-relay -->

ccam ترمینال شما را به Claude Code نمی‌سپارد (`runLogin`). ورودی استاندارد Claude Code را به یک pipe وصل می‌کند و هر خطی را که تایپ می‌کنید در آن کپی می‌کند (`relayLogin`)، بنابراین کدی که در اعلان `Paste code here` در Claude Code بچسبانید همچنان می‌رسد. وقتی یک خط نشانی callback باشد، ccam فقط در صورتی آن را می‌پذیرد که میزبان `localhost`، `127.0.0.1` یا `[::1]` باشد، درگاه داده شده باشد، مسیر `/callback` باشد و هر دو `code` و `state` وجود داشته باشند؛ `http://` را می‌توان ننوشت (`callbackURL`). سپس ccam آن نشانی را روی سرور با مهلت ۱۰ ثانیه و بدون دنبال کردن تغییر مسیرها درخواست می‌کند، همان‌طور که مرورگر می‌کرد (`deliverCallback`). درخواست به سرور callback فرمان `claude auth login` منتظر می‌رسد (خط‌چین در نمودار). Claude Code کد را دریافت می‌کند و ورود را ذخیره می‌کند؛ روی Linux در `.credentials.json` داخل پوشه حساب. یک خط چسبانده‌شده هرگز نمی‌تواند ccam را به تماس با میزبان دیگری وادار کند.

ccam پیش از شروع ورود این را چاپ می‌کند:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

در صورت موفقیت `→ Passed the sign-in to Claude Code on this machine.` را چاپ می‌کند.

وقتی ورود در یک ترمینال منتظر ورودی است، ccam پژواک (echo) و حالت خطی ترمینال را خاموش می‌کند و برای هر نویسه‌ای که تایپ می‌کنید یک `*` نشان می‌دهد، بنابراین کد یا نشانی callback چسبانده‌شده روی صفحه دیده نمی‌شود. ccam ویرایش خط را خودش انجام می‌دهد (`readLine`): `Enter` خط را می‌فرستد، `Backspace` یک نویسه را پاک می‌کند و `Ctrl+D` روی خط خالی ورودی را تمام می‌کند. ترمینال وقتی ورود تمام شود به حالت قبل برمی‌گردد، و همچنین با `Ctrl+C`، `SIGTERM` یا `SIGHUP` که همچنان ccam را پایان می‌دهند. اگر ورودی استاندارد و خطای استاندارد به ترمینال وصل نباشند، خط‌ها همان‌طور که می‌رسند کپی می‌شوند.

### تغییر نام، حذف و وارد کردن

`ccam rename` و `ccam remove` روی پوشه حساب و روی آنچه ccam در جاهای دیگر به نام آن حساب نگه می‌دارد عمل می‌کنند. ccam این مسیرها را در `nameState` فهرست می‌کند:

- `~/.local/state/ccam/last-used/<name>`؛
- `~/.local/state/ccam/mcp/<name>`؛
- حافظه‌های نهان مصرف `~/.local/state/ccam/usage/<name>.json` و `<name>.statusline.json`؛
- ثبت سقف مصرف `<name>.limit.json`؛
- یادداشت‌های نشست در `~/.local/state/ccam/handoff/<name>/`.

`ccam remove work` آن‌ها را به این ترتیب انجام می‌دهد:

<!-- motion:accounts-remove -->

1. `ccam remove` ابتدا اگر حساب در حال اجرا باشد آن را رد می‌کند، سپس می‌پرسد و پس از پاسخ شما دوباره بررسی می‌کند (این کار نشستی را که در مدت باز بودن پرسش شروع شده باشد شناسایی می‌کند).
2. پیش از تغییر هر چیزی بررسی می‌کند که بتوان `config.toml` را ذخیره کرد.
3. اگر حساب افزونه‌ها را به اشتراک بگذارد، ccam فایل‌های `installed_plugins.json` و `known_marketplaces.json` گروه را بازنویسی می‌کند (`RewritePluginPaths`)، چون Claude Code مسیرهای مطلق افزونه‌ها را ذخیره می‌کند که ممکن است از پوشه این حساب بگذرند.
4. سپس یک حساب مدیریت‌شده (حسابی که پوشه‌اش را ccam ساخته است، نه حسابی که با `ccam import` ثبت شده) از سیستم خارج می‌شود: `claude auth logout` در پوشه آن اجرا می‌شود، پس هیچ ورود کارآمدی در سطل زباله منتظر نمی‌ماند (یا، در macOS، در یک مورد Keychain که به نام پوشه است و حساب بعدی با همان نام آن را به ارث می‌برد). این مرحله برای حساب کلید API و وقتی پوشه حساب دیگری همان ورود را دارد، که خروج آن را هم باطل می‌کرد، انجام نمی‌شود. خروجی که ناموفق باشد گزارش می‌شود و حذف ادامه می‌یابد.
5. ccam پوشه را به `~/.local/share/ccam/trash/<name>-<time>/` منتقل می‌کند (`profile.Trash`)، حساب از `config.toml` حذف می‌شود، مسیرهای `nameState` پاک می‌شوند تا حساب بعدی با همین نام آن‌ها را به ارث نبرد، و پیوند میان‌بر حذف می‌شود: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

اگر در نهایت نتوان `config.toml` را ذخیره کرد، پوشه از سطل زباله برمی‌گردد. برای بازیابی یک حساب حذف‌شده، پوشه را از سطل زباله برگردانید و دوباره وارد شوید.

`ccam rename` نام پوشه را تغییر می‌دهد و مسیرهای افزونه در `plugins/` آن را از پوشه قدیمی به پوشه جدید بازنویسی می‌کند. مسیرهای `nameState` را به نام جدید منتقل می‌کند و اگر حساب پیش‌فرض بوده باشد، آن را پیش‌فرض نگه می‌دارد. یادداشت‌های نشست را هم به‌روز می‌کند تا `ccam next --continue` همچنان آن‌ها را پیدا کند. فرمان میان‌بر نام خود را نگه می‌دارد. اگر بازنویسی افزونه‌ها ناموفق باشد، یا نتوان `config.toml` را ذخیره کرد، نام پوشه به حالت قبل برمی‌گردد.

در macOS، و در Windows وقتی Claude Code ورودها را در Credential Manager نگه می‌دارد، خود پوشه نام قدیمی‌اش را نگه می‌دارد، چون ورود و ورودهای MCP زیر مسیر پوشه ذخیره می‌شوند و پوشه جابه‌جاشده از سیستم خارج می‌شد. در این صورت پیکربندی پوشه را به‌صورت `folder` کنار نام جدید ثبت می‌کند و ccam می‌گوید `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

`ccam import` یک پوشه را در همان جایی که هست، به‌صورت `dir` در `config.toml` ثبت می‌کند. هرگز آن را جابه‌جا یا کپی نمی‌کند: یک توکن تمدید کپی‌شده باعث می‌شد نسخه‌ها یکدیگر را از حساب خارج کنند.

بدون آرگومان پوشه، به دنبال پوشه‌های ثبت‌نشده `~/.claude` و `~/.claude-*` می‌گردد و درباره هر کدام می‌پرسد. نام‌ها از پوشه گرفته می‌شوند (`~/.claude` می‌شود `main`، `~/.claude-work` می‌شود `work`، `~/.claude-2` می‌شود `acct2`). اگر ccam نتواند نامی به دست آورد، یا نام گرفته شده باشد، پوشه را همراه با فرمان `--name` که باید اجرا کنید فهرست می‌کند. پوشه پیداشده‌ای که ناموفق باشد به‌صورت `Skipped <path>: <reason>` گزارش می‌شود و بقیه وارد می‌شوند.

ccam ابتدا همه پرسش‌هایش را می‌پرسد، سپس `config.toml` را دوباره می‌خواند و پوشه‌ها را ثبت می‌کند، پس تغییراتی که فرمان‌های دیگر ccam در این مدت ذخیره کرده‌اند حفظ می‌شوند. خط‌های `✓` پس از ذخیره می‌آیند.

بدون `CLAUDE_CONFIG_DIR`، Claude Code فایل `.claude.json` را در پوشه خانه شما کنار `~/.claude` نگه می‌دارد، بنابراین وقتی `~/.claude` را وارد می‌کنید و `.claude.json` خودش را ندارد، ccam پیوند `~/.claude/.claude.json` را به `~/.claude.json` می‌سازد. در این صورت `claude` معمولی و `ccam main` یک حساب را می‌بینند.

اگر `~/.claude` کنار `~/.claude.json` یک `.claude.json` از خودش داشته باشد، ccam هر دو را دست‌نخورده می‌گذارد و این را می‌گوید. `ccam main` فایل درونی و `claude` معمولی فایل بیرونی را به کار می‌برد و جزئیات ورود و سرورهای MCP آن‌ها از هم جدا می‌ماند. برای داشتن یک فایل، فایل درونی را با یک پیوند جایگزین کنید ([doctor](/features/doctor) هم آن را گزارش می‌کند).

وقتی پوشه همان ورودِ یک حساب ثبت‌شده را داشته باشد (یک پوشه کپی‌شده)، ccam این را می‌گوید: هر کدام زودتر آن را تمدید کند دیگری را از سیستم خارج می‌کند، پس پیش از استفاده از آن `ccam login <name>` را اجرا کنید.

حساب‌های واردشده تا وقتی آن‌ها را به گروهی اضافه نکنید (`ccam group join <group> <name>`) در هیچ گروهی نیستند و میان‌بری هم برایشان ساخته نمی‌شود.

`ccam completion bash|zsh` یک اسکریپت تکمیل خودکار چاپ می‌کند که نام حساب‌ها را هم تکمیل می‌کند.

## ایمنی و موارد خاص

- `ccam rename` و `ccam remove` تا وقتی حساب در حال اجراست (`refuseRunning`) با ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)`` رد می‌کنند، چون نشست همچنان در پوشه می‌نویسد. این بررسی Claude Code را به هر شکلی که شروع شده باشد پوشش می‌دهد: از طریق ccam، `claude` معمولی، یک IDE، نشست‌های پس‌زمینه و daemon آن‌ها.
- `ccam rename`، `ccam remove` و `ccam logout` تا وقتی فرایند ccam دیگری آن ورود را تمدید یا تغییر می‌دهد رد می‌کنند: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- تا وقتی تغییر گروه در فرایند ccam دیگری در جریان است، `ccam rename` و `ccam remove` با `another ccam process is changing the groups right now; try again when it is done` رد می‌کنند؛ و پس از تغییری که نیمه‌کاره متوقف شده باشد، با ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.
- `ccam remove` روی یک حساب واردشده پوشه و ورود آن را برای `claude` معمولی نگه می‌دارد. اگر حساب در گروه‌هایی باشد، ابتدا مانند `ccam group leave` از آن‌ها بیرون می‌رود: هر مورد مشترک به یک کپی از نسخه گروه در پوشه تبدیل می‌شود و سرورهای MCP گروه در `.claude.json` آن قرار می‌گیرند. سپس ثبت آن لغو می‌شود: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. نسخه‌های پشتیبانی که ccam از `.claude.json` آن ساخته است به `<folder>/.ccam-backup/mcp-<time>` منتقل می‌شوند؛ اگر نتوان آن‌ها را منتقل کرد، در `~/.local/state/ccam/mcp/<name>/backups` می‌مانند و پیام این را می‌گوید. `~/.claude` شما به سطل زباله منتقل نمی‌شود.
- دو حساب نمی‌توانند یک پوشه داشته باشند: `config.toml`ی که یک پوشه را دو بار نام ببرد بارگذاری نمی‌شود.
- `ccam add` نامی را که پوشه‌اش از قبل در `profiles/` وجود دارد رد می‌کند و شما را به `ccam import` ارجاع می‌دهد. در macOS همچنین وقتی Keychain هنوز یک ورود یا ورودهای MCP را برای مسیر پوشه جدید نگه می‌دارد، که از نسخه قدیمی‌تری از ccam مانده است، متوقف می‌شود و فرمان `security delete-generic-password` را که آن‌ها را پاک می‌کند نام می‌برد. اگر نتوان `config.toml` را ذخیره کرد، پوشه جدید دوباره حذف می‌شود.
- `ccam import` پوشه‌ای را که از قبل ثبت شده، مسیری را که پوشه نیست و نامی را که گرفته شده رد می‌کند. ccam همچنین وقتی پوشه کپی‌ای از `~/.claude` باشد (همان ورود) هشدار می‌دهد، حتی اگر خود `~/.claude` ثبت نشده باشد.
- `ccam logout` وقتی حساب در حال اجراست ابتدا می‌پرسد؛ `ccam login` بررسی نمی‌کند.
- اگر واسطه نتواند به Claude Code برسد، برای مثال پس از پایان ورود، `ccam: could not pass the sign-in to Claude Code:` را همراه با دلیل چاپ می‌کند. خطای HTTP معمولاً یعنی کد منقضی شده است؛ ورود را دوباره شروع کنید.

## مطالب مرتبط

- [اجرای Claude Code](/features/launch)
- [فرمان‌های میان‌بر](/features/shortcuts)
- [اشتراک‌گذاری سرورهای MCP، تنظیمات و تاریخچه](/docs/sharing)
- [استفاده از چند حساب Claude Code روی یک رایانه](/guides/multiple-claude-code-accounts)
