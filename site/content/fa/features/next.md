---
title: جابه‌جایی به حساب بعدی
description: ccam next، Claude Code را با حساب دیگری باز می‌کند؛ --auto پرظرفیت‌ترین را می‌گیرد، --continue گفتگو را منتقل می‌کند و --fallback کارهای -p را تکرار می‌کند.
updated: 2026-11-05
order: 8
translation: machine
sourceHash: "de1fd573f1ac"
---

`ccam next` بدون آنکه نام حسابی را بدهید، Claude Code را با حساب دیگری باز می‌کند. وقتی حسابی به سقف مصرف خود رسید از آن استفاده کنید. چهار راه برای جابه‌جایی وجود دارد:

- `ccam next` ساده حساب پس از حسابی را که آخرین بار استفاده کردید برمی‌دارد و `--auto` حسابی را که بیشترین ظرفیت را دارد. `--continue` آخرین گفتگوی این پوشه را روی حساب جدید ادامه می‌دهد، چه آن دو تاریخچه را به اشتراک بگذارند چه نه.
- درون Claude Code، [پلاگین ccam](#inside-claude-code) به شما می‌گوید کی سقف مصرف یک نشست را متوقف کرده است، و `/ccam:switch` گفتگو را به حساب دیگری می‌سپارد.
- برای اسکریپت‌هایی که `claude -p` را اجرا می‌کنند، `ccam --fallback` جابه‌جایی را خودش انجام می‌دهد ([اجرای دوباره هنگام رسیدن به سقف](#run-again-at-the-limit) را ببینید).
- با روشن بودن [`auto_switch`](#auto-switch) (Linux و macOS)، ccam نشست تعاملی‌ای را هم که Claude Code در سقف مصرفش متوقف می‌کند جابه‌جا می‌کند.

**صندلی‌های Team و Enterprise.** انتخاب خودکار حساب (`ccam next`، `--fallback`، `ccam pick`، کلید `n` در داشبورد و، وقتی آن را روشن کنید، [`auto_switch`](/features/next#auto-switch)) میان صندلی‌های همان سازمان Claude Team یا Enterprise هم جابه‌جا می‌شود، که ممکن است شرایط Anthropic را نقض کند. پیش از آنکه به آن تکیه کنید، [پرسش‌های متداول](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) را بخوانید و از سازمان خود و Anthropic بپرسید.

## نحوه استفاده

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # continue this folder's last conversation on the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next -c --brief            # a new conversation that starts from a summary of the last one
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the most weekly quota left per day
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work CCAM_PROFILE=work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam جابه‌جایی را به شکل `→ <from> → <to>` در خطای استاندارد چاپ می‌کند، همراه با دلیل انتخابی که `--auto` کرده است، و سپس Claude Code را با حساب جدید به همان شکلی آغاز می‌کند که `ccam <name>` انجام می‌دهد. `<from>` حسابی است که آخرین بار استفاده کردید، یا با `--continue` حسابی که گفتگو به آن تعلق دارد. آرگومان‌های پس از `--` بدون تغییر به `claude` می‌روند.

## گزینه‌ها

| گزینه | کاری که انجام می‌دهد | پیش‌فرض |
| --- | --- | --- |
| `--auto` | حساب واردشده را بر اساس مصرف انتخاب می‌کند (پایین‌تر را ببینید) و از حساب‌های ۱۰۰٪ می‌گذرد | خاموش |
| `--max-usage <1-100>` | با `--auto`: حساب‌هایی را که مصرفشان زیر این مقدار است ترجیح می‌دهد | `auto_max_usage`، در غیر این صورت ۱۰۰ |
| `--strategy headroom\|expiring\|balance` | با `--auto`: کمترین مصرف ۵ ساعته، زودترین بازنشانی ۵ ساعته یا بیشترین سهمیه هفتگی باقی‌مانده به ازای هر روز | `auto_strategy`، در غیر این صورت `headroom` |
| `--cross-org` | با `--auto`: حساب‌های سازمان دیگری را هم انتخاب می‌کند (پایین‌تر را ببینید) | `cross_org`، در غیر این صورت خاموش |
| `--wait[=<duration>]` | با `--auto`: وقتی هیچ حسابی ظرفیت ندارد، منتظر می‌ماند تا نخستین حساب بازنشانی شود، حداکثر به این مدت. مقدار پس از `=` می‌آید: `--wait 2h` با راهنمایی `--wait takes its value after '=': --wait=2h` رد می‌شود | خاموش؛ اگر تنها داده شود `5h` |
| `-c`، `--continue` | آخرین گفتگوی این پوشه را با حساب جدید ادامه می‌دهد (پایین‌تر را ببینید) | خاموش |
| `--brief` | با `--continue`: به‌جای از سر گرفتن کل گفتگوی قبلی، گفتگوی تازه‌ای را با حساب جدید از خلاصه‌ای که ccam از آن نوشته آغاز می‌کند ([پایین‌تر را ببینید](#brief)) | خاموش |
| `-- <claude args>` | آرگومان‌هایی که پس از آرگومان‌های خود ccam به `claude` داده می‌شوند | هیچ |
| `--dry-run` | پشتیبانی می‌شود: `→ <from> → <to>` و خط فرمانی را که اجرا می‌کرد چاپ می‌کند، بدون آغاز Claude Code یا ثبت یک اجرا | خاموش |

`--auto` و `--continue` با هم کار می‌کنند. با `--continue`، آرگومان‌های `--resume` پیش از آرگومان‌های خود شما قرار می‌گیرند.

گزینه‌های خود Claude Code، مانند `--dangerously-skip-permissions` یا `--model`، پس از `--` می‌آیند. برخلاف `ccam <name>` که هر چه پس از نام بیاید را منتقل می‌کند، `ccam next` آرگومان‌هایش را خودش می‌خواند، پس بدون `--` با `ccam: unknown flag: --dangerously-skip-permissions` متوقف می‌شود.

## سازوکار

هر اجرا از طریق ccam زمان را در `~/.local/state/ccam/last-used/<name>` می‌نویسد. `ccam next` این فایل‌ها را می‌خواند، بررسی می‌کند کدام حساب‌ها وارد شده‌اند و انتخاب می‌کند.

1. ccam برای هر حساب در `config.toml` یک `next.Candidate` می‌سازد (`cli.candidates`): حساب، زمان آخرین اجرای آن (`launcher.LastUsed`) و اینکه ورودی کارآمد دارد یا نه (`creds.Has`: یک ورود Claude در `.credentials.json`، یا در macOS در Keychain، که به‌طور کامل منقضی نشده و Claude Code آن را پاک نکرده است). پوشه‌ای که به جای آن از کلید API استفاده می‌کند (یک ورود Console، یا حسابی که با `ccam add --api` اضافه شده) سقف طرح ندارد و هرگز انتخاب نمی‌شود.
2. بدون `--auto`، ccam حساب‌ها را بر اساس شماره جایگاه مرتب می‌کند و حسابی را که اخیراً از همه اجرا شده پیدا می‌کند (`next.Rotate`). این همان حسابی است که از آن جابه‌جا می‌شوید. اگر هنوز هیچ حسابی اجرا نشده باشد، پایین‌ترین جایگاه است.
3. از آنجا به ترتیب جایگاه جلو می‌رود، از آخرین جایگاه دوباره به نخستین برمی‌گردد و نخستین حساب واردشده را برمی‌گرداند.
4. ccam حساب انتخاب‌شده را دقیقاً مانند `ccam <name>` آماده و اجرا می‌کند (`launcher.LaunchWith`): پیوندهای حساب را بررسی می‌کند، اشتراک‌گذاری MCP را اعمال می‌کند، نشست را ثبت می‌کند و ccam را از طریق `syscall.Exec` با `claude` جایگزین می‌کند.

پیمایش از حسابی که آخرین بار استفاده شده آغاز می‌شود و در نخستین حساب واردشده پس از آن متوقف می‌شود:

<!-- motion:next-rotate -->

`work` اخیراً از همه اجرا شده، پس پیمایش پس از آن آغاز می‌شود. `side` ورود ندارد و رد می‌شود و `lab` انتخاب می‌شود. اگر `lab` هم از حساب خارج شده بود، پیمایش دور می‌زد و به `personal` می‌رسید (`% len` در نمودار). چرخش هرگز حسابی را که از آن شروع کرده‌اید برنمی‌گرداند: اگر حساب واردشده دیگری نباشد، با خطا متوقف می‌شود.

### --auto

`--auto` به‌جای ترتیب جایگاه، بر اساس مصرف تصمیم می‌گیرد. ccam حساب‌های واردشده را نگه می‌دارد، مصرف همه آن‌ها را هم‌زمان به همان روشی که `ccam usage` می‌پرسد جویا می‌شود (`Fetcher.FetchAll`) و طبق قواعد زیر انتخاب می‌کند (`cli.pickAuto`، `usage.Policy.Choose()`):

<!-- motion:next-auto -->

هر حساب با پنجره‌های ۵ ساعته و هفتگی‌اش سنجیده می‌شود، و وقتی `--model` میان آرگومان‌های Claude Code باشد (پس از `--`)، با پنجره هفتگی همان مدل هم (`opus`، `sonnet`، `haiku` یا `fable` در نام). ccam حساب‌ها را در سه دسته قرار می‌دهد (`Policy.Choose()`):

1. **ظرفیت دارد**: همه پنجره‌ها زیر `auto_max_usage` (`config.toml`، پیش‌فرض ۱۰۰). راهبرد از میان آن‌ها انتخاب می‌کند: `headroom` (پیش‌فرض) کمترین مصرف ۵ ساعته را برمی‌دارد، و از میان حساب‌هایی که مصرف ۵ ساعته برابر دارند، حسابی را که کمترین بخش از هفته‌اش مصرف شده است؛ `expiring` پنجره ۵ ساعته‌ای را که زودتر از همه بازنشانی می‌شود، تا سهمیه باقی‌مانده‌اش پیش از هدر رفتن مصرف شود؛ `balance` حسابی را برمی‌دارد که تا بازنشانی هفته‌اش، به ازای هر روز بیشترین سهمیه هفتگی را باقی دارد، و هر نشست ccam در حال اجرا را به حساب آن می‌گذارد، تا سهمیه هفتگی بی‌استفاده نماند و نشست‌های موازی پخش شوند.
2. **بالای حد**: زیر ۱۰۰٪ اما برابر یا بالاتر از `auto_max_usage`. فقط وقتی استفاده می‌شود که هیچ حسابی ظرفیت نداشته باشد؛ حسابی که کمترین مصرف ۵ ساعته را دارد اول می‌آید، و در صورت برابری حسابی که کمترین بخش از هفته‌اش مصرف شده است.
3. **پر**: پنجره‌ای در ۱۰۰٪، یا مصرف ۵ ساعته نامعلوم. هرگز انتخاب نمی‌شود.

پنجره‌ای که Claude Code درخواستی را در آن رد کرده تا بازنشانی‌اش ۱۰۰٪ حساب می‌شود: هوک [پلاگین ccam](#inside-claude-code) آن را ثبت می‌کند (`~/.local/state/ccam/usage/<name>.limit.json`)، هر چه عددهای ذخیره‌شده بگویند. پس حسابی که همین حالا به سقفش رسیده پیش از بازنشانی‌اش دوباره انتخاب نمی‌شود، حتی وقتی عددهایش قدیمی‌اند.

تنظیمات همه حساب‌ها در بالای `config.toml` می‌آیند، بالای نخستین `[[group]]`، `[[profile]]` یا `[[folder]]`. تنظیمی که زیر یکی از آن‌ها باشد به همان جدول تعلق دارد: ccam در هر اجرا درباره آن هشدار می‌دهد و فرمان‌هایی که پیکربندی را ذخیره می‌کنند تا جابه‌جا شدنش سر باز می‌زنند.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**سازمان‌ها.** انتخاب‌های خودکار در میان حساب‌هایی از همان نوع حسابی که از آن می‌روید می‌مانند: صندلی‌های همان سازمان Team یا Enterprise، یا طرح‌های شخصی (Pro، Max)، تا کاری که از یک حساب شرکتی شروع شده در یک حساب شخصی ادامه پیدا نکند، و برعکس. ccam سازمان هر حساب را از ورود آن می‌داند (`ccam login` آن را ثبت می‌کند؛ ورودهای قدیمی‌تر از `.claude.json` خوانده می‌شوند)؛ حسابی که ccam نتواند سازمانش را تشخیص دهد با هر حسابی جور می‌شود. دو پوشه که به همان حساب Claude وارد شده‌اند سقف‌های آن را مشترک دارند، پس هیچ‌کدام هرگز جای دیگری را نمی‌گیرد. `--cross-org`، یا `cross_org = true` در `config.toml`، قاعده سازمان را برمی‌دارد، اما پس از [`/ccam:switch`](#inside-claude-code) جابه‌جایی به سازمان دیگر همچنان منتظر «بله» شما می‌ماند.

وقتی هیچ حسابی را نتوان انتخاب کرد، ccam با کد ۷۵ (`EX_TEMPFAIL`) خارج می‌شود تا اسکریپت بتواند «بعداً دوباره امتحان کن» را از شکست تشخیص دهد. با `--wait` به جای آن صبر می‌کند تا پنجره بازدارنده نخستین حساب پر بازنشانی شود (۳۰ ثانیه پس از زمان بازنشانی، وقتی نقطه پایانی پنجره جدید را نشان می‌دهد)، `ccam: no account has room; waiting until 15:10` را چاپ می‌کند و دوباره انتخاب می‌کند؛ وقتی هیچ حسابی نگوید کی بازنشانی می‌شود، پس از ۱۰ دقیقه دوباره نگاه می‌کند. اگر پیش از پایان انتظار هیچ حسابی آزاد نشود، با کد ۷۵ خارج می‌شود.

این پرس‌وجو از قواعد `ccam usage` پیروی می‌کند: عددهایی که کمتر از ۶۰ ثانیه از عمرشان گذشته دوباره استفاده می‌شوند، مهلت پرس‌وجوی هر حساب ۳ ثانیه است، حسابی که پاسخ ندهد یا توکن دسترسی‌اش منقضی شده باشد بر اساس عددهای ذخیره‌شده‌اش سنجیده می‌شود، و ثبت سقف مصرف بر آن‌ها مقدم است. پنجره ذخیره‌شده‌ای که زمان بازنشانی‌اش گذشته باشد ۰٪ حساب می‌شود، پس حسابی که پنجره‌اش از آن زمان بازنشانی شده نادیده گرفته نمی‌شود. `--auto` هرگز توکن‌ها را تمدید نمی‌کند. برخلاف چرخش ساده، `--auto` می‌تواند همان حسابی را که آخرین بار استفاده کردید انتخاب کند، اگر هنوز بیشترین جا را داشته باشد.

### --continue

1. ccam آخرین گفتگوی این پوشه را در همه حساب‌ها جست‌وجو می‌کند: از میان فایل‌های `<session>.jsonl` در `projects/<folder>/` (نام پوشه همان پوشه کاری است، با پیوندهای حل‌شده، که هر نویسه‌ای جز حرف یا رقم در آن به `-` تبدیل شده، همان‌طور که Claude Code نام‌گذاری می‌کند)، آن که تازه‌ترین پرامپت یا پاسخ را دارد. نشست‌هایی که هیچ‌کدام را ندارند، مانند نشستی که فقط `/mcp`، `/usage` یا `/status` را اجرا کرده، یا فایلی خالی، کنار گذاشته می‌شوند. اگر نشستی یادداشتی به جا گذاشته باشد (سقف مصرف یا `/ccam:switch`، پایین‌تر را ببینید) و هیچ گفتگویی بیش از ۲ دقیقه پس از آن پرامپت یا پاسخی نداشته باشد، به‌جای آن دقیقاً همان نشست استفاده می‌شود؛ ردیف‌هایی که Claude Code هنگام `/exit` می‌نویسد حساب نمی‌شوند.
2. حسابی که گفتگو به آن تعلق دارد همان حسابی است که از آن جابه‌جا می‌شوید، و حساب جدید مانند بالا، در همان سازمان، انتخاب می‌شود. وقتی چند حساب آن تاریخچه را به اشتراک بگذارند، گفتگو به حسابی تعلق دارد که `.claude.json` آن می‌گوید آخرین بار گفتگو را اجرا کرده است (`lastSessionId`)، وگرنه به حسابی که آخرین بار استفاده شده است.
3. اگر حساب جدید آن تاریخچه را به اشتراک بگذارد (همان پوشه گروه)، ccam فرمان `claude --resume <session id>` را آغاز می‌کند، یعنی دقیقاً همان گفتگو را (`--continue` جدیدترین گفتگوی پوشه را، هر کدام که باشد، باز می‌کرد).
4. در غیر این صورت `claude --resume <the conversation's file> --fork-session` را آغاز می‌کند: Claude Code کل گفتگو را با یک شناسه (ID) نشست جدید در تاریخچه خود حساب جدید کپی می‌کند و نسخه اصلی همان‌طور که بود می‌ماند. از آن پس `--continue` حساب جدید آن را پیدا می‌کند.
5. مگر اینکه `--agent` را پس از `--` بدهید، ccam عاملی را هم که گفتگو با آن اجرا می‌شود منتقل می‌کند: تنظیم `agent` در `.claude/settings.local.json` یا `.claude/settings.json` پوشه یا در `settings.json` حساب، و در غیر این صورت عاملی که آخرین ردیف `agent-setting` گفتگو نام می‌برد. Claude Code هنگام کپی کردن یک گفتگو (مرحله ۴) عامل را برنمی‌گرداند و با عامل پیش‌فرض خود، با همه ابزارها، ادامه می‌داد. وقتی حساب جدید عامل‌ها، پلاگین‌ها و تنظیمات حساب قبلی را مشترک ندارد (همان پوشه‌های `agents` و `plugins` و همان `settings.json`، از طریق یک گروه)، ccam پس از خط مربوط به جابه‌جایی یک خط اضافه می‌کند: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

ccam می‌گوید گفتگو چقدر بزرگ است: حساب جدید با نخستین پاسخ خود همه آن را دوباره می‌خواند. حافظه‌های نهان پرامپت هرگز بین سازمان‌ها مشترک نمی‌شوند و هر طرح شخصی سازمانی جداگانه است، پس وقتی دو حساب به سازمان‌های متفاوتی تعلق دارند، ccam این را اضافه می‌کند: `prompt caches are not shared between organizations, so the first reply reads it all again`. وقتی حساب قدیمی به سرورهای MCP وارد شده باشد که حساب جدید وارد نشده است، ccam نامشان را می‌برد: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. اشتراک `history` در یک گروه از کپی جلوگیری می‌کند، یک گفتگو را نگه می‌دارد و نقطه‌های بازگشتی را که `/rewind` کد را از آن‌ها بازمی‌گرداند (`file-history/`) هم مشترک می‌کند.

### --brief

`--continue` کل گفتگو را از سر می‌گیرد و حساب جدید با نخستین پاسخ خود همه آن را دوباره می‌خواند. `ccam next --continue --brief` به‌جای آن گفتگوی تازه‌ای را آغاز می‌کند، از خلاصه‌ای که ccam بدون پرسیدن از یک مدل، از گفتگوی قدیمی می‌نویسد:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

خلاصه، به همین ترتیب و فقط جایی که چیزی برای گفتن باشد، این‌ها را در بر دارد:

1. عنوان گفتگو؛
2. مسیر رونوشت (transcript) آن، تا Claude هر وقت به جزئیاتی نیاز داشت در آن جست‌وجو کند؛
3. نخستین درخواست شما، وقتی بیش از چهار درخواست بوده باشد؛
4. خلاصه‌ای که Claude Code آخرین بار هنگام فشرده کردن گفتگو نوشته است؛
5. چهار درخواست آخر شما، با فرمان‌های `!` و فرمان‌های اسلش همان‌طور که تایپشان کردید؛
6. آخرین پاسخ Claude؛
7. فهرست کارها (todo)؛
8. ۱۰ فایلی که Claude آخر از همه تغییر داده است؛
9. `git status --short --branch` پوشه، تا ۸۰۰ نویسه.

فهرست کارها حداکثر ۱۰ مورد نگه می‌دارد: نخست موارد انجام‌شده کنار می‌روند، سپس دورترین گام‌ها، تا گام فعلی و گام‌های بعدی بمانند. بخش‌های طولانی از وسط بریده می‌شوند و کل خلاصه زیر ۱۶٬۰۰۰ نویسه می‌ماند.

ccam آن را در `~/.local/state/ccam/brief/<session>.md` ذخیره می‌کند که فقط شما می‌توانید آن را بخوانید، و `--brief` بعدی خلاصه‌های قدیمی‌تر از ۷ روز را پاک می‌کند. Claude Code با یک پرامپت آغاز می‌شود و فایلی را که در آن نام برده شده پیوست می‌کند: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` پرامپتی که پس از `--` بدهید جای جمله آخر را می‌گیرد (`Then: <your prompt>`) و گزینه‌های دیگرتان پس از `--` به دنبال آن می‌آیند.

Claude Code یک اشاره (mention) را در `"` پایان می‌دهد و در `#` می‌بُرد، پس وقتی مسیر خلاصه یکی از این دو را داشته باشد، ccam با مسیر آن نسبت به پوشه کاری به آن اشاره می‌کند؛ وقتی آن مسیر هم یکی از آن‌ها را داشته باشد، پرامپت نام فایل را می‌برد و Claude آن را با ابزار Read خود می‌خواند، که اول اجازه می‌پرسد چون فایل بیرون از پوشه‌های کاری است (در یک اجرای `-p` در حالت Manual، Claude Code خواندن را رد می‌کند و Claude بدون خلاصه ادامه می‌دهد).

`ccam --dry-run next --continue --brief` خلاصه را چاپ می‌کند و چیزی ذخیره نمی‌کند. گفتگوی قدیمی همان‌طور که بود می‌ماند و گفتگوی جدید فقط آنچه را خلاصه می‌گوید می‌داند، پس وقتی از `--brief` استفاده کنید که خود کار از همه جزئیات راهی که تا آنجا آمده‌اید مهم‌تر است. بدون `--continue`، ccam با `--brief works with --continue` متوقف می‌شود؛ وقتی نتواند گفتگو را بخواند یا چیزی برای خلاصه کردن در آن پیدا نکند، متوقف می‌شود و می‌گوید که `ccam next --continue` بدون `--brief` آن را به‌طور کامل از سر می‌گیرد. [رله](#relay) و [`auto_switch`](#auto-switch) همیشه کل گفتگو را از سر می‌گیرند. گزینه `--brief` خود Claude Code، مانند دیگر گزینه‌هایش، پس از `--` می‌آید.

## درون Claude Code

[پلاگین ccam](/docs/install#claude-code-plugin) دو هوک اضافه می‌کند (Claude Code 2.1.78 یا جدیدتر):

- وقتی یک نوبت در سقف مصرف متوقف شود، گفتگو را یادداشت می‌کند و در ترمینال‌هایی که از آن پشتیبانی می‌کنند (iTerm2، WezTerm، Windows Terminal، kitty، Ghostty) یک اعلان دسکتاپ نشان می‌دهد: ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  همه پیام‌هایی را که Claude Code سقف مصرف به حساب می‌آورد می‌شناسد: سقف‌های ۵ ساعته، هفتگی، ساده و سقف مدل (`You've hit your session limit`، `You've hit your limit`، `You've reached your …` برای سهمیه یک مدل)، اعتبارهای مصرف و سقف‌های هزینه، و سازمانی که حساب را متوقف کرده است.

  وقتی Claude Code درخواست ردشده را در گفتگو نوشته باشد، هوک پنجره پر را هم تا بازنشانی‌اش ثبت می‌کند ([ثبت سقف مصرف](/features/usage#the-limit-record))، جز برای توقفی که سازمان تعیین کرده، که مدیرانش هر زمان می‌توانند آن را بردارند. محدودیت کوتاه‌مدت سرور یا محدودیت نرخ یک کلید API سقف مصرف نیست و به حال خود رها می‌شود، و همچنین نشست‌های بدون حضور کاربر (اجراهای `-p` و نشست‌های پس‌زمینه). اگر یک subagent یا عامل گردش‌کار در سقف مصرف متوقف شود، این توقف نشست خودش به حساب می‌آید، حتی وقتی گفتگوی اصلی فقط منتظر آن است: Claude Code آن درخواست ردشده را در فایل خود آن عامل در پوشه `subagents` نشست می‌نویسد و هوک آن را از همان‌جا می‌خواند. وقتی چند عامل هم‌زمان متوقف شوند، نشست یک یادداشت و یک اعلان می‌گیرد.

  Claude Code 2.1.234 و جدیدتر هم در نشست منتظر می‌ماند و پس از بازنشانی ۵ ساعته خودش ادامه می‌دهد؛ اعلان به شما این انتخاب را می‌دهد که به‌جای آن جابه‌جا شوید، و با روشن بودن [`auto_switch`](#auto-switch) (Linux و macOS)، وقتی حساب دیگری ظرفیت داشته باشد ccam به‌جای شما جابه‌جا می‌شود. وقتی هیچ حساب دیگری نتواند کار را ادامه دهد، چون هیچ حساب دیگری به همان سازمان Team یا Enterprise تعلق ندارد، یا هیچ‌کدام طرح شخصی دیگری نیست (با روشن بودن `cross_org`، هر حساب Claude دیگری می‌تواند)، اعلان به‌جای آن با `You have no other personal plans, so ccam cannot move this conversation.` یا `You have no other accounts in Acme, so ccam cannot move this conversation.` پایان می‌یابد.
- `/ccam:switch [account]` گفتگو را یادداشت می‌کند و نوبت را پیش از فرستادن هر درخواستی پایان می‌دهد: `/exit` را تایپ کنید، سپس `ccam next --continue` را اجرا کنید تا آن را روی حسابی که نام بردید (یا روی حساب بعدی) ادامه دهید. در نشستی که هنوز گفتگویی ندارد، می‌گوید چیزی برای ادامه دادن نیست. فقط آنچه خودتان تایپ می‌کنید آن را اجرا می‌کند: از پلاگین 0.3.2 به بعد، Claude نمی‌تواند `/ccam:switch` را اجرا کند، و کاری که Claude زمان‌بندی می‌کند (CronCreate، ScheduleWakeup یا `.claude/scheduled_tasks.json`) آن را به‌صورت متن ساده به Claude می‌دهد.

  وقتی حسابی که نام می‌برید بیرون از سازمان گفتگو باشد (صندلی‌ای از سازمان Team یا Enterprise دیگری، یک طرح شخصی برای گفتگویی که روی یک صندلی است یا برعکس، یا یک حساب کلید API در هر یک از دو طرف)، `ccam next --continue` اول می‌پرسد، صرف‌نظر از مقدار `cross_org`: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. همین را پیش از حسابی که خودش انتخاب می‌کند هم می‌پرسد، وقتی شما حسابی نام نبرید؛ `cross_org = true` به آن اجازه می‌دهد این حساب را از سازمان دیگری بردارد: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. فقط `y` یا `yes` گفتگو را جابه‌جا می‌کند؛ در Linux و macOS، کلیدهایی که پیش از پرسش زده شده‌اند اول کنار گذاشته می‌شوند. هر پاسخ دیگری، یا نبودن ترمینالی برای پرسیدن، گفتگو و یادداشت را همان‌جا که هستند نگه می‌دارد: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (اگر حساب را خود ccam انتخاب کرده بود، `ccam next --continue` دوباره انتخاب می‌کند). مانند انتخاب‌های خودکار، حسابی که ccam نتواند سازمانش را تشخیص دهد با هر حسابی جور می‌شود.

یادداشت در `~/.local/state/ccam/handoff/<account>/<session>.json` می‌ماند، یک فایل برای هر نشست، پس دو نشست یک حساب هرگز یادداشت یکدیگر را بازنویسی نمی‌کنند. ۱۲ ساعت نگه داشته می‌شود و یک بار استفاده می‌شود، درست پیش از آغاز نشست بعدی؛ اگر اجرا شکست بخورد، یادداشت برای تلاش بعدی می‌ماند. یادداشت به پوشه‌ای تعلق دارد که نشست در آن آغاز شده است، حتی پس از یک `cd` در شل Claude Code. `ccam rename` نام‌های درون آن را به‌روز می‌کند؛ اگر حسابی که نام می‌برد از آن زمان حذف شده باشد، ccam این را اعلام می‌کند و طبق معمول انتخاب می‌کند.

### رله (relay)

با `relay = true` در بالای `config.toml`، ccam به‌جای جایگزین کردن خودش با Claude Code در کنار نشست‌های تعاملی در حال اجرا می‌ماند و پس از `/ccam:switch` گفتگوی همان نشست را خودش ادامه می‌دهد (به یادداشتی که نشست دیگری گذاشته باشد عمل نمی‌کند): هوک Claude Code را پایان می‌دهد و ccam فرمان `ccam next --continue` را با یادداشت همان نشست اجرا می‌کند. اگر حسابی نام برده نشده باشد، حسابی را برمی‌دارد که بیشترین ظرفیت را دارد، با در نظر گرفتن مدل گفتگو، و وقتی `cross_org = true` اجازه دهد که این حساب از سازمان دیگری باشد، پیش از جابه‌جا کردن گفتگو به آنجا همان پرسش بالا را می‌پرسد. در Windows هوک نمی‌تواند Claude Code را پایان دهد: `/exit` را تایپ کنید تا ccam گفتگو را ادامه دهد. هوک برای حسابی بیرون از سازمان گفتگو هم Claude Code را پایان نمی‌دهد: می‌گوید `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).` و ccam پس از آنکه `/exit` را تایپ کردید همان پرسش بالا را می‌پرسد. این کار برای هر نشست یک فرایند ccam منتظر هزینه دارد و به همین دلیل به‌طور پیش‌فرض خاموش است.

نشستی که این‌گونه ادامه می‌یابد چیزهایی را نگه می‌دارد که از سر گرفتن یک گفتگو بازنمی‌گرداند: گزینه‌هایی مانند `--add-dir`، `--mcp-config`، فهرست‌های ابزار، `--settings`، `--plugin-dir` و گزینه‌های پرامپت سیستمی، حالت مجوز گفتگو، سطح تلاش نشست و عاملی که گفتگو با آن اجرا می‌شود ([`--continue`](#continue) را ببینید). خود Claude Code گزینه‌های `--model` و `--name` را از گفتگو بازمی‌گرداند. هر چیز دیگری (`--worktree`، `--debug`، `--remote-control`، گزینه‌هایی که ccam نمی‌شناسد) در `ccam: not carried over: …` نام برده می‌شود. این جابه‌جایی هرگز پرسش تاریخچه مشترک را نمی‌پرسد، و آنچه در حین اجرای آن تایپ کرده‌اید کنار گذاشته می‌شود؛ تنها پرسش آن همان پرسشی است که پیش از حسابی از سازمان دیگر پرسیده می‌شود.

رله فقط وقتی گفتگو را ادامه می‌دهد که Claude Code به همان شکلی پایان یافته باشد که `/ccam:switch` خواسته بود: نه وقتی خودتان ابتدا `/exit` را تایپ کرده‌اید، هنگام پایان آن Ctrl+C زده‌اید، یا ccam یک hangup یا سیگنال دیگری گرفته است، و نه وقتی ترمینال از بین رفته است. در این صورت یادداشت می‌ماند و `ccam next --continue` آن را بعداً ادامه می‌دهد. نشستی که زیر رله نسخه قدیمی‌تری از ccam شروع شده، یا نشستی بدون حضور کاربر (اجراهای `-p` و نشست‌های پس‌زمینه)، به دست هوک پایان نمی‌یابد؛ به جای آن پیام `/exit` را می‌گیرد. `claude remote-control` بدون رله اجرا می‌شود.

### `auto_switch`

با `auto_switch = true` در بالای `config.toml` (Linux و macOS)، وقتی Claude Code نشستی را در سقف مصرف متوقف کند، ccam خودش آن را به حساب دیگری منتقل می‌کند: هوک پلاگین حساب دیگری را که ظرفیت دارد برمی‌دارد و Claude Code را پایان می‌دهد، و ccam، که در کنار نشست در حال اجراست، گفتگو را در آنجا ادامه می‌دهد، همان‌طور که [رله](#relay) پس از `/ccam:switch` انجام می‌دهد. این قابلیت به‌طور پیش‌فرض خاموش و آزمایشی است.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam فقط وقتی نشستی را جابه‌جا می‌کند که همه این شرط‌ها برقرار باشند:

1. ccam آن را در یک ترمینال آغاز کرده و در کنارش در حال اجراست. `auto_switch` [رله](#relay) را روشن می‌کند، پس این شرط برای نشست‌های تعاملی‌ای که پس از تنظیم آن آغاز می‌کنید برقرار است. نشستی که پیش از آن، یا زیر یک ccam قدیمی‌تر، آغاز شده اعلانی می‌گیرد که با ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` پایان می‌یابد. نشست‌های بدون حضور کاربر (اجراهای `-p` و نشست‌های پس‌زمینه)، `claude remote-control` و Claude Code‌ای که بدون ccam آغاز شده هرگز جابه‌جا نمی‌شوند.
2. [پلاگین ccam](/docs/install#claude-code-plugin) در آن حساب روشن است و حساب Claude Code 2.1.214 یا جدیدتر را اجرا می‌کند، که شناسه (ID) فرایند خود را به هوک می‌دهد (`CLAUDE_PID`).
3. Claude Code درخواست ردشده را، همراه با زمان بازنشانی، برای سقف ۵ ساعته، هفتگی، Opus یا Sonnet آن حساب در گفتگو نوشته است. سقف هزینه، تمام شدن اعتبارهای مصرف، توقفی که سازمان تعیین کرده، سهمیه Fable (`You've reached your Fable limit.`)، محدودیت کوتاه‌مدت سرور و محدودیت نرخ یک کلید API هرگز نشستی را جابه‌جا نمی‌کنند. سقف یک subagent یا عامل گردش‌کار به حساب می‌آید، حتی وقتی گفتگوی اصلی فقط منتظر آن است؛ اگر این سقف، سقف هفتگی Opus یا Sonnet باشد، نشست فقط وقتی جابه‌جا می‌شود که خود گفتگو با همان مدل اجرا شود.
4. ccam طرح حساب را می‌داند، از `.claude.json` آن یا از آنچه `ccam login` ثبت کرده است. در غیر این صورت اعلان با `ccam does not know personal's plan, so auto_switch leaves this session.` پایان می‌یابد.
5. نشست به Remote Control وصل نیست. با حساب دیگر دوباره به برنامه‌ای که نشست را از آن می‌گردانید وصل نمی‌شد، پس می‌ماند و اعلان این را می‌گوید.
6. سقف بیش از ۱۵ دقیقه دیگر بازنشانی می‌شود. برای انتظار کوتاه‌تر، ccam نشست را به Claude Code می‌سپارد: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. از زمان رسیدن به سقف، پرامپتی تایپ نکرده‌اید.
8. حساب دیگری ظرفیت دارد. ccam همان‌طور انتخاب می‌کند که [`--auto`](#auto) انتخاب می‌کند (`auto_max_usage`، `auto_strategy`، مدل گفتگو)، از میان حساب‌هایی از همان نوع حسابی که از آن می‌روید: صندلی دیگری از همان سازمان Team یا Enterprise، یا طرح شخصی دیگری (Pro، Max). `cross_org` و `--cross-org` در اینجا اعمال نمی‌شوند. حساب باید:
   - با ورودی وارد شده باشد که Claude Code بتواند آن را تمدید کند؛
   - به همان حساب Claude وارد نشده باشد؛
   - از کلید API استفاده نکند؛
   - طرحی داشته باشد که ccam می‌شناسد؛
   - تا جایی که پرس‌وجوی مصرف خود ccam نشان می‌دهد، اعتبارهای مصرف را روشن نداشته باشد؛
   - در Claude Code به این پوشه اعتماد کرده باشد؛
   - اگر گفتگو با یک عامل اجرا می‌شود (تنظیم `agent` یا خود گفتگو عاملی را نام می‌برد)، عامل‌های حسابی را که از آن می‌روید مشترک داشته باشد: پوشه‌های `agents` و `plugins` و `settings.json` آن حساب باید هر کدام دقیقاً همان باشند (از طریق یک گروه)، یا در هر دو وجود نداشته باشند.

   وقتی این سه شرط آخر همه حساب‌های دارای ظرفیت را کنار بگذارند، اعلان آن‌ها را نام می‌برد: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`، `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` یا `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` پنجره‌ای که [ثبت سقف مصرف](/features/usage#the-limit-record) دارد پر حساب می‌شود، و حسابی که نشست دیگری در ۲ دقیقه گذشته به آن جابه‌جا شده پس از بقیه می‌آید. انتخاب باید در حدود ۶ ثانیه انجام شود؛ وقتی پرس‌وجوهای مصرف بیشتر طول بکشند، نشست می‌ماند. با `usage_source = "claude"`، `auto_switch` فقط از عددهای ذخیره‌شده تصمیم می‌گیرد، و ccam نمی‌تواند از پرس‌وجوی Claude Code یا از خط وضعیت تشخیص دهد که اعتبارهای مصرف روشن است یا نه، پس فرض می‌کند خاموش است.

سپس هوک Claude Code را پایان می‌دهد و ccam، که در کنار آن در حال اجراست، در ترمینال‌هایی که از آن پشتیبانی می‌کنند یک اعلان دسکتاپ نشان می‌دهد، مانند `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` خط نخست خروجی نمونه بالا را چاپ می‌کند و `ccam next --continue` را برای آن نشست با حساب انتخاب‌شده اجرا می‌کند، که دو خط دیگر آن را چاپ می‌کند و Claude Code را با گزینه‌ها، حالت مجوز، سطح تلاش و عامل همان نشست آغاز می‌کند، همان‌طور که [رله](#relay) انجام می‌دهد. نخستین پیام نشست جدید یک متن ثابت است، تنها چیزی که ccam از طرف شما می‌فرستد (جابه‌جایی پس از `/ccam:switch` چیزی نمی‌فرستد):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

پس از سقف یک subagent، پیام این یکی است، چون subagentها، گردش‌کارها و فرمان‌های پس‌زمینه همراه با نشست قبلی متوقف شده‌اند:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

اگر آن حساب هم به سقفش برسد، همین اتفاق در آنجا می‌افتد. ثبت‌های سقف مصرف نمی‌گذارند این زنجیره پیش از بازنشانی یک حساب به آن برگردد، پس وقتی همه حساب‌ها پر باشند، نشست روی آخرین حساب می‌ماند.

هوک از تنظیم `agent`، یا از ردیفی که عامل را نام می‌برد و Claude Code هر از گاهی و هنگام پایان نشست در گفتگو می‌نویسد، می‌فهمد که گفتگو با یک عامل اجرا می‌شود. اگر تا این لحظه فقط `--agent` نشست آن را نام برده باشد، و حساب انتخاب‌شده عامل‌ها را مشترک نداشته باشد، `ccam next --continue` چیزی را آغاز نمی‌کند و یادداشت را نگه می‌دارد، و یک اعلان دسکتاپ دوم این را می‌گوید: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

وقتی هیچ حساب دیگری ظرفیت نداشته باشد، چیزی جابه‌جا نمی‌شود: Claude Code طبق معمول با سقف کنار می‌آید، یعنی در نشست منتظر می‌ماند یا گزینه‌هایش را پیشنهاد می‌کند، و اعلان می‌گوید کدام حساب زودتر از همه آزاد می‌شود، برای نمونه `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` برای صندلی‌ها می‌گوید `None of your other accounts in Acme has room now`. برای متوقف کردن جابه‌جایی‌ای که آغاز شده، هنگام بسته شدن Claude Code کلید Ctrl+C را بزنید: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

جابه‌جایی هزینه‌هایی دارد. Claude Code دوباره آغاز می‌شود که چند ثانیه طول می‌کشد، و شل‌های پس‌زمینه، سرورهای توسعه و عامل‌هایی که Claude آغاز کرده همراه با نشست قدیمی متوقف می‌شوند. حساب جدید با نخستین پاسخ خود کل گفتگو را دوباره می‌خواند، مانند [`--continue`](#continue)، و ccam می‌گوید گفتگو چقدر بزرگ است؛ جابه‌جایی هرگز از [`--brief`](#brief) استفاده نمی‌کند.

`auto_switch` هرگز پیش از آنکه Claude Code سقفی را گزارش کند جابه‌جا نمی‌کند، هرگز مصرف شما را در پس‌زمینه زیر نظر نمی‌گیرد و هرگز اطلاعات ورودی را کپی، تمدید یا منتقل نمی‌کند. در Windows نادیده گرفته می‌شود و `ccam doctor` این را می‌گوید؛ در آنجا، با `relay = true`، `/ccam:switch` و `/exit` همچنان گفتگو را به حساب دیگری منتقل می‌کنند.

**شرایط Anthropic.** در صندلی‌های Team و Enterprise، `auto_switch` کار شما را به صندلی دیگری از همان سازمان منتقل می‌کند. مرکز راهنمای Anthropic می‌گوید: «Each seat on a Team plan is assigned to one named individual, and only that person may use it» (هر صندلی در طرح Team به نام یک فرد مشخص اختصاص داده می‌شود و فقط همان شخص اجازه دارد از آن استفاده کند) و «Users can't route their work through another user's account to use that user's remaining usage» (کاربران نمی‌توانند کار خود را از طریق حساب کاربر دیگری انجام دهند تا از سهمیه مصرف باقی‌مانده آن کاربر استفاده کنند) ([طرح‌های Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)؛ و همین‌طور برای [طرح‌های Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)).

برای طرح‌های شخصی (Pro، Max)، شرایط Anthropic نمی‌گوید که آیا حساب‌های خودتان اجازه دارند در سقف مصرف به‌طور خودکار جای یکدیگر را بگیرند یا نه. [شرایط مصرف‌کننده (Consumer Terms)](https://www.anthropic.com/legal/consumer-terms) دسترسی به خدمات را «through automated or non-human means, whether through a bot, script, or otherwise» (از راه‌های خودکار یا غیرانسانی، چه با بات، چه با اسکریپت و چه به هر شکل دیگر) ممنوع می‌کند، مگر آنکه Anthropic صراحتاً اجازه دهد، و یک جابه‌جایی یکی از پیام‌های ثابت بالا را از طرف شما می‌فرستد. [سیاست استفاده (Usage Policy)](https://www.anthropic.com/legal/aup) که از ۱۲ نوامبر ۲۰۲۶ اجرا می‌شود، «Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products» (دور زدن عمدی قابلیت‌ها، محدودیت‌ها یا حفاظ‌های تعبیه‌شده در محصولات ما، یا مداخله در آن‌ها) را در میان کارهای ممنوع آورده است.

ccam ادعا نمی‌کند که `auto_switch` با این شرایط سازگار است. روشن کردن آن تصمیم خودتان است، و خطرش هم با خودتان است: Anthropic ممکن است وقتی باور داشته باشد شرایط نقض شده‌اند، یک حساب یا دسترسی یک سازمان Team یا Enterprise را تعلیق کند.

وقتی `auto_switch` روشن است، `ccam doctor` درباره حسابی که پلاگین ccam ندارد و درباره نسخه‌های Claude Code قدیمی‌تر از 2.1.214 هشدار می‌دهد، و برای صندلی‌ها قاعده صندلی Anthropic را تکرار می‌کند ([بررسی‌ها](/features/doctor)).

## اجرای دوباره هنگام رسیدن به سقف

`ccam next` برای نشست تعاملی‌ای است که متوقف شده است. یک اجرای غیرتعاملی، یعنی `claude -p` در یک اسکریپت یا یک کار زمان‌بندی‌شده، نمی‌تواند منتظر شما بماند، پس `--fallback` پیش از نام حساب به ccam اجازه می‌دهد به جای آن جابه‌جا شود:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` فقط همراه با `-p` یا `--print` کار می‌کند؛ در غیر این صورت ccam متوقف می‌شود و به `ccam next --auto --continue` اشاره می‌کند.
2. ورودی استاندارد لوله‌شده یک بار، تا ۱۰ مبی‌بایت، خوانده می‌شود و به هر تلاش داده می‌شود. ccam مانند Claude Code حداکثر ۳ ثانیه منتظر می‌ماند تا آغاز شود (جز با `--input-format stream-json`)؛ اگر لوله‌ای باز بماند اما چیزی نفرستد، اجرا ورودی‌ای نمی‌گیرد و ccam این را در خطای استاندارد اعلام می‌کند.
3. ccam به جای `exec`، `claude` را به‌عنوان فرایند فرزند اجرا می‌کند، خروجی استاندارد آن را نگه می‌دارد و خطای استاندارد را همان‌طور که می‌رسد عبور می‌دهد.
4. وقتی `claude` با کد ۰ خارج شود، یا به دلیل دیگری ناموفق شود، خروجی نگه‌داشته‌شده و کد خروج بدون تغییر منتقل می‌شوند. اجرایی که با یک سیگنال پایان یافته باشد (کد خروج بالاتر از ۱۲۸) هم همین‌طور: از آن خواسته شده بود متوقف شود و دوباره اجرا نمی‌شود.
5. یک شکست وقتی رسیدن به سقف مصرف به حساب می‌آید که پیام Claude Code این را بگوید، هر یک از پیام‌های سقف مصرف که [هوک](#inside-claude-code) می‌شناسد، حتی درون خروجی `--output-format json` یا `stream-json`، یا وقتی مصرف حساب نشان دهد پنجره ۵ ساعته، هفتگی یا `--model` آن به ۱۰۰٪ رسیده است (از جمله با ثبت سقف مصرف).
6. آنگاه خروجی آن تلاش کنار گذاشته می‌شود و همان فرمان با حساب واردشده دیگری اجرا می‌شود: حسابی که کمترین مصرف ۵ ساعته را دارد، سپس حساب‌هایی که مصرفشان را نمی‌توان پرسید، به ترتیب شماره. حساب‌هایی که به ۱۰۰٪ رسیده‌اند کنار گذاشته می‌شوند، حساب‌هایی که با کلید API وارد شده‌اند هرگز امتحان نمی‌شوند و هر حساب حداکثر یک بار امتحان می‌شود.
7. اگر `-c`، `--continue`، `-r` یا `--resume` میان آرگومان‌ها باشد، فقط حساب‌های همان گروه تاریخچه امتحان می‌شوند تا گفتگو پیدا شود.
8. مانند `--auto`، فقط حساب‌های همان سازمان امتحان می‌شوند (`cross_org = true` این قاعده را برمی‌دارد).

حساب بعدی با همان قواعد `--auto` انتخاب می‌شود (`auto_max_usage`، `auto_strategy` و پنجره `--model`). وقتی هیچ حسابی باقی نماند، ccam خروجی آخرین تلاش را چاپ می‌کند، می‌گوید هیچ حساب واردشده‌ای ظرفیتی باقی ندارد و با کد ۷۵ خارج می‌شود. با `--wait` پیش از نام حساب (`ccam --fallback --wait=2h work -p …`)، به جای آن منتظر می‌ماند تا نخستین حساب بازنشانی شود و سپس ممکن است همه حساب‌ها را دوباره امتحان کند؛ همیشه دست‌کم یک بار ابتدا صبر می‌کند، تا حسابی که تازه به سقفش رسیده بلافاصله دوباره اجرا نشود. این مدت کل زمان انتظار اجرا را محدود می‌کند، نه هر انتظار را. اگر با فاصله نوشته شود (`--wait 30m`)، با همان راهنمایی `ccam next` رد می‌شود. یک `--` آغازین پیش از آرگومان‌های Claude Code مانند هر اجرای دیگری کار می‌کند (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` فرمان نخستین تلاش را چاپ می‌کند.

نشست‌های تعاملی خودبه‌خود جابه‌جا نمی‌شوند، مگر اینکه [`auto_switch`](#auto-switch) را روشن کنید. با `ccam statusline install`، وقتی مصرف ۵ ساعته یا هفتگی به ۹۰٪ برسد، خط وضعیت `→ ccam next --auto --continue` را اضافه می‌کند ([سقف مصرف](/features/usage)).

## ایمنی و موارد خاص

- حساب‌های بدون ورود هرگز انتخاب نمی‌شوند. اگر حساب واردشده دیگری نباشد، `ccam next` ساده با `no logged-in account other than <name>` متوقف می‌شود.
- با `--continue` و بدون گفتگویی در پوشه: `no account has a conversation in <folder> to continue`.
- وقتی قاعده سازمان هیچ حسابی باقی نگذارد: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- با `--auto` و بدون حساب واردشده: `no logged-in account to pick`. وقتی همه حساب‌ها پر هستند یا مصرفشان نامعلوم است: `` no account with known usage and headroom (see `ccam usage`) ``. هر دو با کد خروج ۷۵ پایان می‌یابند.
- `--max-usage`، `--strategy`، `--wait` و `--cross-org` فقط با `--auto` کار می‌کنند، و مقداری بیرون از بازه 1-100 یا راهبرد ناشناخته رد می‌شود.
- انشعاب (fork) گفتگو را کپی می‌کند، نه وضعیت دیگر حساب را: نقطه‌های بازگشت `file-history/` حساب اصلی همراه آن نمی‌آیند، مگر اینکه تاریخچه مشترک باشد.
- `--auto` حتی با `--dry-run` درخواست شبکه می‌فرستد. فقط می‌خواند و هرگز در فایل‌های توکن نمی‌نویسد.
- حساب انتخاب‌شده از همان مسیر `ccam <name>` آغاز می‌شود، پس همان بررسی‌ها اعمال می‌شوند، از جمله هشداری که وقتی نشست Claude Code دیگری با همان تاریخچه مشترک از قبل در این پوشه در حال اجراست داده می‌شود.
- `last-used` فقط وقتی نوشته می‌شود که Claude Code واقعاً آغاز شود. `--dry-run` به آن دست نمی‌زند، پس اجرای آزمایشی چرخش را جلو نمی‌برد.
- `--fallback` خروجی استاندارد را تا پایان هر تلاش نگه می‌دارد، پس خروجی جریانی مانند `--output-format stream-json` یک‌جا در پایان می‌رسد. کلاینتی که ورودی `stream-json` را جریانی می‌فرستد و پیش از فرستادن بیشتر منتظر پاسخ می‌ماند نمی‌تواند از `--fallback` استفاده کند. فایل‌هایی که تلاش نخست پیش از رسیدن به سقف تغییر داده، تغییریافته می‌مانند؛ تلاش بعدی درخواست را از نو آغاز می‌کند.

## مطالب مرتبط

- [سقف مصرف](/features/usage)
- [سقف مصرف و جابه‌جایی بین حساب‌ها](/docs/usage-limits)
- [اشتراک‌گذاری سرورهای MCP، تنظیمات و تاریخچه](/docs/sharing)
- [وقتی Claude Code به سقف مصرف می‌رسد چه کنیم](/guides/usage-limit-workflow)
