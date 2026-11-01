---
title: حدود الاستخدام
description: يعرض ccam usage حدّي الخمس ساعات والأسبوع لكل الحسابات دفعة واحدة، ويضع ccam statusline اسم الحساب واستخدامه داخل Claude Code.
updated: 2026-11-01
order: 7
translation: machine
sourceHash: "54b296f7da35"
---

يعرض `ccam usage`، لكل حساب، مقدار ما استُخدم من كل حد: حد الخمس ساعات، والحد الأسبوعي، والحد الأسبوعي الخاص بكل نموذج (مثل Fable). ولكل حد شريط والوقت المتبقي حتى إعادة ضبطه. استخدمه قبل أن تبدأ مهمة طويلة، أو عندما يخبرك Claude Code أن حدًا ما يقترب، لترى أي حساب لا تزال لديه سعة. و`ccam statusline` هو أمر سطر الحالة في Claude Code: يعرض اسم الحساب واستخدامه داخل Claude Code ويحفظ الأرقام لـ `ccam usage`. ولمعرفة الرموز (tokens) والتكاليف التقديرية التي وراء هذا الاستخدام، يشغّل `ccam ccusage` أداة [ccusage](https://ccusage.com) على محادثات كل الحسابات ([أدناه](#tokens-and-costs-ccam-ccusage)).

## الاستخدام

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts that are not running
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

يحصل كل حساب على سطر فيه رقمه واسمه وعمر أرقامه؛ وتعني `(cached)` أنها تأتي من ملف محفوظ، لا من نقطة النهاية. وتحته يعرض سطر لكل حد شريطًا من 20 خانة والنسبة المئوية والوقت المتبقي حتى إعادة ضبط النافذة: `5-hour`، ثم `Weekly` (كل النماذج)، ثم النافذة الأسبوعية الخاصة بكل نموذج، مثل `Fable`. وتظهر النافذة التي لا يملكها الحساب بالرمز `-`. وتتحول النسبة إلى الأصفر عند 80% وإلى الأحمر عند 100%.

إذا رفض Claude Code طلبًا لأن نافذة امتلأت، تعرض تلك النافذة 100% حتى إعادة ضبطها، مهما قالت الأرقام المحفوظة (راجع [سجل الحد](#the-limit-record)). في الطرفية تظهر كل الحسابات فورًا بالحالة `fetching…`، وتُرسم أسطر كل حساب عندما يصل رده.

يذكر السطر الأخير، الذي لا يظهر إلا إذا كان لديك أكثر من حساب، الحساب الذي سيختاره [`ccam next --auto`](/features/next). وهو يتبع القواعد نفسها:

- `auto_max_usage` و`auto_strategy`؛
- الجلسات العاملة؛
- حسابات مؤسسة آخر حساب مستخدَم فقط، ما لم يُعيَّن `cross_org`؛
- تُتخطى كل نافذة بلغت 100%، فلا يُذكر أبدًا حساب امتلأ أسبوعه.

ومع الاستراتيجية الافتراضية `headroom` يكون `→ Most headroom: <name>`؛ ومع `expiring` أو `balance` يذكر السبب، مثل `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**مقاعد Team وEnterprise.** يبدّل الاختيار التلقائي للحساب (`ccam next` و`--fallback` و`ccam pick` والمفتاح `n` في لوحة التحكم، و[`auto_switch`](/features/next#auto-switch) بعد أن تشغّله) أيضًا بين مقاعد المؤسسة نفسها في Claude Team أو Enterprise، وهذا قد يخالف شروط Anthropic. اقرأ [الأسئلة الشائعة](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) وتحقّق من ذلك مع مؤسستك ومع Anthropic قبل أن تعتمد عليه.

يعرض الحساب غير المسجَّل دخوله `not logged in → ccam login <name>`، ويعرض حساب مفتاح API الرسالة `no limits (API)`. ويعرض رمز الوصول المنتهي ``access token expired; `ccam <name>` renews it when it starts``، فوق الأرقام المحفوظة إن وُجدت: يجدّده Claude Code عند التشغيل التالي. أما تسجيل الدخول الذي انتهى بالكامل، لأن رمز التحديث الخاص به انتهى أو لأن Claude Code مسحه بعد أن رفضه الخادم، فيعرض `` login expired → `ccam login <name>` ``.

بعد `ccam statusline install`، يعرض Claude Code سطرًا مثل `[work] · Opus · 5h 24% · 7d 41%`.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--json` | يطبع التقارير بصيغة JSON: `name` و`slot` و`status` و`limits` (مع `credits` عندما تكون أرصدة الاستخدام مفعّلة للحساب) و`source` و`at` و`error`، و`limited` و`limited_windows` عندما يجعل سجلُّ حد نوافذَ بقيمة 100%. ودون حسابات يطبع `[]` | معطّل |
| `--watch[=N]` | يعيد الرسم كل `N` ثانية حتى تضغط Ctrl+C. و`--watch` وحده، أو قيمة 0 أو أقل، يعني 60. وتأتي القيمة بعد `=`: فيُرفض `--watch 10` مع تلميح | معطّل |
| `-v`، `--verbose` | يضيف سطر `source:` تحت كل حساب: مصدر أرقامه (نقطة نهاية الاستخدام، أو Claude Code، أو سطر الحالة، أو ملف محفوظ، أو حد الاستخدام الذي أبلغ عنه Claude Code). وعندما يتقدّم سجلُّ حد على الأرقام، يليه سطر مثل `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` | معطّل |
| `--no-cache` | يتجاهل ذاكرة التخزين المؤقت ذات الستين ثانية ويطلب من نقطة النهاية مرة أخرى | معطّل |
| `--source live\|claude\|cache` | يطلب `live` من نقطة النهاية باستخدام رمز الحساب؛ ويجعل `claude` برنامج Claude Code يستعلم عنه (انظر أدناه)؛ ويعرض `cache` الأرقام المحفوظة فقط | `usage_source`، وإلا `live` |
| `--refresh` | يجدّد رموز الوصول المنتهية، للحسابات التي لا تعمل فقط. وقد يستغرق التجديد ما يصل إلى 30 ثانية | معطّل |
| `--dry-run` | مقبول: لا يغيّر `ccam usage` شيئًا سوى ذاكرة التخزين المؤقت الخاصة به. ويُرفض مع `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | معطّل |

أي قيمة أخرى لـ `--source` تفشل مع `--source must be live, claude or cache`. ولا يأخذ `ccam usage` أي وسائط.

### الاستخدام من Claude Code نفسه

مع `usage_source = "claude"` في أعلى `config.toml` (فوق أول `[[group]]` أو `[[profile]]` أو `[[folder]]`)، أو `--source claude`، لا يقرأ ccam رمز تسجيل الدخول إطلاقًا. فلكل حساب يشغّل `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` في مجلد الحساب، ويرسل طلب `get_usage` الخاص بـ Claude Code عبر الإدخال القياسي (Claude Code 2.1.270 أو أحدث). فيستعلم Claude Code عن الاستخدام بتسجيل دخوله هو، ويجيب بنوافذ الخمس ساعات والأسبوع والنماذج.

ويمنع `--safe-mode` خطافات الحساب وخوادم MCP والمكوّنات الإضافية من البدء، ولا تُحفظ أي جلسة. ويحتاج كل حساب إلى تشغيل واحد لـ Claude Code (نحو ثانية)، لذلك هذا المصدر أبطأ من `live`. وتُسأل الحسابات كلها في آن واحد، مع مهلة 20 ثانية لكل منها، وتنطبق الذاكرة المؤقتة لمدة 60 ثانية كالمعتاد.

وتسجيل الدخول الذي انتهى رمز الوصول الخاص به، أو سينتهي خلال 5 دقائق، يحصل على 50 ثانية، لأن Claude Code يجدّده عند بدئه، والتجديد الذي يُقطع قد يُضيّع تسجيل الدخول. وإذا تخلّى ccam عن استعلام، يطلب من Claude Code أن يخرج بالإشارة SIGTERM، فيُكمل Claude Code تجديدًا بدأه (لا توجد SIGTERM في Windows، لذا تُنهى العملية هناك).

وتسجيل الدخول الذي انتهى بالكامل يعرض `` login expired → `ccam login <name>` `` دون تشغيل Claude Code. ويصنّف Claude Code هذا الطلب على أنه تجريبي، لذلك إذا توقف عن الاستجابة، يعرض ccam الأرقام المحفوظة ويبلّغ عنه `ccam doctor`. وعندما يجيب Claude Code بأن استعلامه هو فشل، يعرض ccam الأرقام المحفوظة مع `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`، لا خطة بلا حدود.

ومع `--dry-run`، لا يشغّل ccam برنامج Claude Code، لأن Claude Code يكتب في مجلد الحساب؛ بل يعرض الأرقام المحفوظة بدلًا من ذلك. ويقبل `usage_source` القيمة `live` أو `claude`؛ أما `cache` فهو خيار لـ `--source` فقط.

يقرأ `ccam statusline` بيانات JSON الخاصة بسطر الحالة في Claude Code من الإدخال القياسي. وخياره الوحيد هو `--wrap <command>`: أمر سطر حالة موجود يُضاف اسم الحساب أمام أول سطر من مخرجاته. ولا يأخذ `ccam statusline install` أي خيارات. ويرفض كلاهما `--dry-run` مع `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## الرموز والتكاليف: `ccam ccusage`

يعرض `ccam usage` مقدار ما استخدمه الحساب من كل حد. أما [ccusage](https://ccusage.com) فيحصي الرموز (tokens) والتكاليف التقديرية في سجلات محادثات Claude Code. وعند تشغيله وحده لا يبحث عن سجلات Claude Code إلا في `~/.claude` و`~/.config/claude` (أو في المجلدات التي يحددها `CLAUDE_CONFIG_DIR`)، لذا فحساب ccam الوحيد الذي يستطيع رؤيته هو حساب مستورد في مكانه عند `~/.claude`. ويشغّله `ccam ccusage` مع ضبط `CLAUDE_CONFIG_DIR` على مجلدات كل الحسابات التي فيها محادثات، فيغطيها كلها تقرير واحد:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

يقرأ ccam خياراته هو أولًا: يأخذ `-a` (`--account`) أسماء الحسابات أو أرقامها، مفصولة بفواصل أو بتكرار الخيار، ويطبع `--dry-run` الأمر بدلًا من تشغيله. وبدءًا من أول وسيط ليس من خياراته، يذهب كل شيء إلى ccusage كما كتبته، باستثناء `--config` مع `bun x` و`npx` (أدناه)؛ ويعرض `ccam ccusage -- --help` نص المساعدة الخاص بـ ccusage. ويشغّل ccam الأمر `ccusage` إذا كان مثبّتًا، وإلا `bun x ccusage`، فإن تعذّر ذلك `npx ccusage@latest`؛ وهذان الأخيران يجلبانه من npm. ويبدأ تشغيلهما في مجلد ccam نفسه (`~/.local/state/ccam/ccusage`)، لا في المجلد الذي أنت فيه، فلا يستطيع مشروعٌ موجود فيه أن يحدد ما يشغّلانه عبر `node_modules/.bin` الخاص به أو `bin` في `package.json` الخاص به أو سجل حزم (registry) في `.npmrc` الخاص به؛ ويعرض `ccam --dry-run ccusage` الأمر `cd` إلى ذلك المجلد. وعندما لا يتوفر أي منها، يطلب منك تثبيت Node.js أو Bun أو ccusage. على Linux وmacOS يحل ccusage محل ccam؛ أما على Windows فينتظره ccam ثم ينتهي برمز خروج ccusage. وتحت WSL، يتخطى ccam أي `ccusage` أو `bun` أو `npx` موجود على محرك أقراص Windows، لأنه سيعمل كبرنامج Windows دون مجلدات الحسابات.

- يغطي التقرير الحسابات التي في مجلداتها محادثات (مجلد `projects`). والمحادثة التي تتشاركها حسابات عبر مجموعة تُحتسب مرة واحدة.
- مع `-a`، إذا كان الحساب يتشارك سجل محادثاته عبر مجموعة، فإنه يجلب معه محادثات المجموعة كلها، ويذكر ccam لمن هي: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- يعرض `blocks` نوافذ الخمس ساعات، ولكل حساب نوافذه الخاصة، لذا حدّد له حسابًا واحدًا، كما في المثال أعلاه.
- يشمل تقرير ccusage أيضًا أدوات البرمجة الأخرى التي يجدها، مثل Codex، ما لم تحدد واحدة منها: فالأمر `ccam ccusage claude daily` يعرض Claude Code وحده.
- ما دام ccam يشغّله، لا يقرأ ccusage إعداداته (`ccusage.json`) إلا من `./.ccusage` ومن مجلدات الحسابات. ويحدد `--config` ملفًا آخر: `ccam ccusage daily --config ~/.config/claude/ccusage.json`. ولأن `bun x` و`npx` يبدأ تشغيلهما في مجلد ccam نفسه، يمرّر لهما ccam الملف `./.ccusage/ccusage.json` مع `--config` عندما يكون ذلك الملف موجودًا ولا تحدد أنت ملفًا، ويمرّر `--config` النسبي بوصفه مسارًا كاملًا. وعندها يقرأ ccusage المسارات النسبية الأخرى، مثل `--pi-path`، من مجلد ccam، لذا اكتبها كاملة.
- كان ccusage سيقرأ المجلد الذي يحتوي مساره على فاصلة كأنه مجلدان، لذا يرفضه ccam ويذكر اسم الحساب؛ فاستبعد ذلك الحساب باستخدام `-a`.
- بما أن `ccusage` أمر من أوامر ccam، يُشغَّل الحساب المسمّى ccusage بالأمر `ccam run ccusage`؛ ويذكر `ccam doctor` ذلك.

## كيف يعمل

تأتي أرقام `ccam usage` من مكانين: نقطة نهاية الاستخدام غير الموثقة لدى Anthropic، أي `GET https://api.anthropic.com/api/oauth/usage`، التي تُستدعى بتسجيل الدخول الخاص بكل حساب، والقيم التي يمرّرها Claude Code إلى أمر سطر الحالة. وكل ذلك موجود في `internal/usage`.

1. يحمّل ccam الملف `config.toml` وينشئ `usage.Account` واحدًا لكل حساب (`cli.usageCmd`): مجلده `~/.local/share/ccam/profiles/<name>/`، وعلامات جلساته (التي تبيّن، مع سجلات Claude Code نفسه في مجلد الحساب، هل يعمل)، وهل استُورد في مكانه.
2. يستعلم ccam عن كل الحسابات في آن واحد، بروتين متزامن (goroutine) واحد لكل حساب (`Fetcher.FetchAll`). وكلها تتشارك `http.Client` واحدًا، فيُعاد استخدام الاتصالات بنقطة النهاية. قد يستغرق الاتصال 2 s، ولكل حساب 3 s إجمالًا. ولا يتبع العميل إعادة التوجيه، لأن إعادة التوجيه سترسل الرمز إلى مكان آخر، ويقرأ 1 MB على الأكثر من الرد.
3. عندما ينتهي كل حساب، يُعاد رسم صفه. والحساب الذي يفشل أو لا يستجيب في الوقت المحدد يلجأ إلى الأرقام المحفوظة، فلا يؤخّر حساب بطيء الحسابات الأخرى أبدًا.

في المخطط، يُستعلم عن كل حساب في روتينه الخاص، والحساب الذي تنتهي مهلته يعرض أرقامه المحفوظة بدلًا من ذلك:

<!-- motion:usage-fetch -->

يستجيب `personal` و`work` في الوقت المحدد ويعرضان أرقامًا حية. أما `side` فيستغرق أكثر من 3 s، لذلك يعرض ccam أرقامه المحفوظة بدلًا من ذلك (`fallback()`). ثم يرسم ccam كتلة لكل حساب (`usage.Table()`)، ويختار الحساب المذكور في السطر الأخير (`usage.Best()`).

لحساب واحد، يجري ccam هذه الفحوص بالترتيب (`Fetcher.Fetch`):

<!-- motion:usage-cache -->

1. مع `--source cache`، ينتقل مباشرة إلى `fallback()`.
2. ما لم يُعيَّن `--no-cache`، تُستخدم النتيجة الحية الموجودة في `~/.local/state/ccam/usage/<name>.json` إذا قيست قبل أقل من 60 ثانية (`CacheTTL`)، ويُحكم عليها بحسب الوقت الحالي: فالنافذة التي أُعيد ضبطها منذ ذلك الحين تظهر 0%، والملف المؤرخ بعد الوقت الحالي للساعة لا يُعدّ حديثًا. فتشغيل `ccam usage` مرتين متتاليتين يكلّف طلبًا واحدًا لكل حساب.
3. يقرأ ccam الملف `.credentials.json` (أو سلسلة المفاتيح على macOS) دون تغييره، ولا يفك إلا المفتاح `claudeAiOauth` (`readCredentials()`). وعدم وجود تسجيل دخول باشتراك يعني `not logged in`، أو `no limits (API)` عندما يحتوي `.claude.json` على مفتاح API. وتسجيل الدخول الذي انتهى رمز التحديث الخاص به، أو الذي مسحه Claude Code، هو `` login expired → `ccam login <name>` `` ولا يُستعلم عنه.
4. لا يُستخدم رمز الوصول المنتهي. ودون `--refresh`، يعرض الحساب أرقامه المحفوظة مع ``access token expired; `ccam <name>` renews it when it starts``.
5. وإلا يستدعي ccam نقطة النهاية. يُعدّ الرد 401 أو 403 تسجيل دخول منتهيًا (`login expired: authentication failed (401)`)، وأي حالة أخرى أو محتوى لا يستطيع قراءته يُعدّ خطأ، وفي الحالتين يلجأ إلى الأرقام المحفوظة. والاستعلام الذي ينهيه Ctrl+C يقول `cancelled`. ويعني الرد 429 أن نقطة النهاية تلقّت طلبات أكثر مما ينبغي لهذا الحساب: يكتب ccam الملف `<name>.pause` مع وقت السؤال التالي (`Retry-After`، أو 5 دقائق، وساعة على الأكثر)، ويعرض الأرقام المحفوظة مع `rate limited by the usage endpoint; asking again after 15:10`، ولا يستدعي نقطة النهاية لذلك الحساب حتى ذلك الحين، حتى مع `--no-cache`. فإعادة المحاولة فورًا لن تؤدي إلا إلى إطالة حد المعدّل.
6. يُكتب الرد السليم في `<name>.json` مع وقت قياسه.

بالنسبة إلى الأرقام المحفوظة، يأخذ ccam أحدث استعلام محفوظ (`<name>.json`، أو `<name>.claude.json` من استعلامات Claude Code)، وعندما يكون `<name>.statusline.json` قد كُتب بعده، يضع نافذتي الخمس ساعات والأسبوع من سطر الحالة فوقه (`fallback()`)؛ أما نوافذ النماذج، التي لا يملكها سطر الحالة، فتبقى كما حفظها الاستعلام. وقبل أن يعرض الأرقام المحفوظة، يتحقق من وقت إعادة ضبط كل نافذة: النافذة التي أُعيد ضبطها منذ حفظ الأرقام تظهر بقيمة 0%، لأن النسبة القديمة لم تعد صحيحة. ولهذا يعرض `personal` أعلاه `0%` دون وقت لإعادة الضبط.

لنقطة النهاية تنسيقان للاستجابة. الأحدث فيه مصفوفة `limits` تحتوي على المدخلات `session` (الخمس ساعات) و`weekly_all` و`weekly_scoped` (لكل نموذج). والأقدم فيه الكائنات `five_hour` و`seven_day` و`seven_day_opus` و`seven_day_sonnet`. يقرأ ccam كليهما (`ParseLive`): تتقدم مصفوفة `limits`، وتكمل الكائنات القديمة ما لا تغطيه. والمحتوى الذي لا يحتوي على أي حدود يُبلَّغ عنه على أنه تنسيق متغيّر، ويخبرك `ccam doctor` بذلك.

### سجل الحد

عندما يوقف Claude Code دورًا عند حد استخدام، يقرأ خطاف [المكوّن الإضافي لـ ccam](/docs/install#claude-code-plugin) الطلب المرفوض من المحادثة: أي نافذة امتلأت (الخمس ساعات، أو الأسبوع، أو Opus، أو Sonnet) ومتى تُعاد ضبطها. ويكتب ذلك في `~/.local/state/ccam/usage/<name>.limit.json`. ويطبّق ccam السجل فوق كل ما يعرفه غير ذلك، من كل مصدر (`Fetch`): كل نافذة فيه تُحتسب 100% حتى إعادة ضبطها، مهما كانت الأرقام المحفوظة قديمة أو منخفضة، فلا يختار `ccam next --auto` و`ccam pick` و`--fallback` و[`auto_switch`](/features/next#auto-switch) الحساب مجددًا قبل ذلك.

وعندما لا يُعرف شيء آخر، يكون السجل وحده ما يعرضه الحساب (`source: the usage limit Claude Code reported` مع `-v`). ولا تكتب الجلسات غير المراقَبة (تشغيلات `-p` والجلسات الخلفية) أي سجل، ولا الإيقاف الذي تفرضه مؤسسة (مثلًا عندما ينفد استخدامها أو يكون حد مجموعة مضبوطًا على $0)، الذي يستطيع مسؤولوها رفعه في أي وقت. ويُسقط `ccam login` إلى حساب Claude آخر السجل؛ وحذف الملف يمسحه أيضًا.

### تجديد الرموز

رموز التحديث تتبدّل: تجديد واحد منها يمنحك رمز تحديث جديدًا ويبطل القديم. لو جدّد ccam تسجيل الدخول بينما يعمل Claude Code بالحساب نفسه، لسُجّل خروج تلك الجلسة. لذلك يقرأ `ccam usage` الرموز ولا يكتبها أبدًا، ما لم تمرّر `--refresh`. عندها فإن ccam (`Fetcher.refresh`):

1. يرفض الحسابات المسجلة في مكانها باستخدام `ccam import`، لأن `claude` العادي قد يستخدم ذلك التسجيل دون علم ccam؛
2. يرفض تسجيلات الدخول المخزنة في سلسلة المفاتيح على macOS؛
3. يأخذ ملف قفل، `~/.local/state/ccam/locks/profile-<name>.refresh.lock`، فلا تجدّد عمليتان من ccam الرمز نفسه كلتاهما أبدًا؛ وما دام القفل ممسوكًا، يعرض الحساب `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)`. ويأخذه أيضًا `ccam rename` و`remove` و`logout`، وترفض العمل ما دام تجديد جاريًا؛
4. تحت القفل، يرفض إذا كانت جلسة للحساب تعمل، أيًّا كانت طريقة تشغيلها (عبر ccam، أو `claude` العادي، أو بيئة تطوير، أو جلسة خلفية والخدمة الخلفية التي وراءها، التي تجدّد تسجيل الدخول بنفسها)، ويتخطى الطلب إذا كانت عملية أخرى قد جدّدت الرمز بالفعل؛
5. لا يرسل شيئًا ما دام Claude Code يمسك بأقفال التجديد الخاصة به في مجلد الحساب (`.oauth_refresh.lock` و`<folder>.lock`): فـ Claude Code يجدّد تسجيل الدخول ذاك في تلك اللحظة؛
6. يرسل رمز التحديث إلى نقطة نهاية الرموز وينتظر الرد حتى 30 ثانية، وهي المدة التي يمنحها Claude Code لتجديداته. وبعد إرسال الطلب لا يُقطع أبدًا، حتى بـ Ctrl+C: فالخادم ربما بدّل رمز التحديث بالفعل، ولا يحمل تسجيلَ الدخول الجديد إلا رده. وينهي Ctrl+C أو قطع الاتصال الاستعلامات الأخرى، ويخرج ccam (بالرمز 130 بعد Ctrl+C) بمجرد حفظ الرد؛
7. يحفظ تحت قفل الكتابة الخاص بـ Claude Code (`.storage-write.lock`) ولا يغيّر إلا حقول الرموز في `.credentials.json` بحالته عندئذٍ: فما كتبه Claude Code في الأثناء (تسجيل دخول MCP) يبقى، وتسجيل الدخول الذي جدّده Claude Code في الأثناء يُحتفظ به، ولا يُلغى تسجيل الخروج أبدًا. وعندما يذكر الرد عمر رمز التحديث الجديد، يُخزَّن بوصفه `refreshTokenExpiresAt`، كما يفعل Claude Code.

التجديد الذي لا يتم يذكر السبب، في سطر الحساب أو فوق أرقامه المحفوظة: `token not renewed: a session is running` مثلًا، أو `login may be lost: …` عندما أُرسل طلب لكن تعذّر استخدام رده أو حفظه، مع نصيحة بتشغيل `ccam login <name>` إذا طلب منك Claude Code بعد ذلك تسجيل الدخول.

### جامع بيانات سطر الحالة

يشغّل Claude Code أمر سطر الحالة ويمرّر إليه مستند JSON. يقرأ ccam منه اسم النموذج و`rate_limits`:

<!-- motion:usage-statusline -->

1. يجد ccam الحساب من `CCAM_PROFILE`، الذي يعيّنه ccam عندما يشغّل Claude Code، وإلا فبمطابقة `CLAUDE_CONFIG_DIR` مع مجلدات الحسابات (`runStatusline`).
2. يقرأ ccam القيمتين `rate_limits.five_hour` و`rate_limits.seven_day` (`ParseStatusLine`). لا يحتوي JSON على حدود للطلبات في حسابات مفاتيح API، أو قبل أول رد في الجلسة؛ وعندها لا يُحفظ شيء.
3. يكتب ccam الملف `~/.local/state/ccam/usage/<name>.statusline.json` (`SaveStatusLine`). ويستخدمه `ccam usage` عندما لا تستجيب نقطة النهاية.
4. يطبع ccam السطر `[name] · model · 5h N% · 7d N%`. ومع `--wrap`، يشغّل أمرك عبر `sh -c` بالإدخال نفسه ويطبع مخرجاته كما يعرض Claude Code سطر الحالة: كل سطر غير فارغ، مع `[name]` أمام الأول. والأمر الذي يفشل يعرض اسم الحساب وحده. وعلى Windows يشغّل الأمر كما يفعل Claude Code: باستخدام Git Bash إذا كان مثبّتًا، وإلا باستخدام PowerShell.
5. عندما يكون الاستخدام في نافذة الخمس ساعات أو النافذة الأسبوعية 90% أو أكثر ويوجد حساب آخر، ينتهي السطر بـ `→ ccam next --auto --continue`، مثلًا `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

ليس للأمر المغلَّف حد زمني خاص به: فـ Claude Code يمنح أمر سطر الحالة 10 دقائق ويوقف التشغيل عندما يبدأ تحديث أحدث، فيظهر الأمر البطيء كما كان سيظهر دون ccam. وعندما ينهي Claude Code برنامج ccam، ينتهي الأمر أيضًا (على Linux وmacOS مع كل عملية بدأها).

يكتب `ccam statusline install` المدخل `statusLine` في ملف `settings.json` لكل حساب ويحتفظ بالمفاتيح الأخرى للكائن، مثل `padding` و`refreshInterval` و`hideVimModeIndicator`. وعندما يكون لدى الحساب أمر لسطر الحالة بالفعل، يحتفظ به ccam بالشكل `--wrap '<old command>'`؛ وعلى Windows بالشكل `--wrap-base64 <base64>`، الذي لا يحتاج إلى علامات اقتباس في أي من الصدفتين. وإذا كان سطر الحالة يشغّل ccam بالفعل (من مسار آخر، أو مقتبسًا بالطريقة التي كتبه بها إصدار أقدم)، يستبدله ccam بدلًا من تغليفه، ويحتفظ بالأمر الذي يغلّفه.

وعلى Windows يُكتب مسار ccam بشرطات مائلة أمامية و`~` لمجلد المنزل (`~/.local/bin/ccam.exe statusline`)، وهو ما يفهمه Git Bash وPowerShell كلاهما. والمسار الذي يحتاج إلى علامات اقتباس يُكتب كاملًا للصدفة التي يستخدمها Claude Code: `'…'` لـ Git Bash، و`& '…'` لـ PowerShell. وملف `settings.json` المشترك عبر مجموعة يُكتب مرة واحدة، عبر رابطه، فتحصل عليه كل حسابات المجموعة. وتشغيله مرة أخرى يطبع `(already set)` للملفات التي تستخدم ccam بالفعل.

## الأمان والحالات الخاصة

- لا يغيّر `ccam usage` أبدًا ملف رموز ما لم تمرّر `--refresh`. وعندما ينتهي رمز وصول، شغّل الحساب مرة واحدة باستخدام `ccam <name>` وسيجدّده Claude Code. أما تسجيل الدخول الذي انتهى بالكامل فيحتاج إلى `ccam login <name>`.
- لا يلمس `--refresh` أبدًا الحسابات التي تعمل (حتى عندما يبدأ Claude Code دون ccam) أو الحسابات المستوردة أو تسجيلات الدخول في سلسلة المفاتيح. والتجديد المرفوض يعرض الأرقام المحفوظة والسبب (`token not renewed: …`).
- نقطة نهاية الاستخدام غير موثقة علنًا وقد تتغير. وعندما يحدث ذلك، تظل ترى الأرقام المحفوظة، ويبلّغ `ccam doctor` عن تغيّر التنسيق.
- عندما لا توجد أرقام محفوظة ولا رد، يعرض سطر الحساب `lookup failed:` والسبب، مثلًا `no saved usage` مع `--source cache`.
- لا تحتوي ملفات ذاكرة التخزين المؤقت إلا على النسب وأوقات إعادة الضبط، ولا تحتوي على رموز أبدًا.

## صفحات ذات صلة

- [حدود الاستخدام والتبديل بين الحسابات](/docs/usage-limits)
- [التبديل إلى الحساب التالي](/features/next)
- [ماذا تفعل عندما يبلغ Claude Code حد الاستخدام](/guides/usage-limit-workflow)
- [مرجع الأوامر](/docs/commands)
