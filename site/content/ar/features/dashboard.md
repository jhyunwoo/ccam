---
title: لوحة التحكم
description: يفتح ccam في الطرفية لوحة تحكم بملء الشاشة للحسابات والمجموعات وقواعد المجلدات والإضافات والفحوص؛ ويعرض كل تغيير خطته أولًا.
updated: 2026-10-30
order: 15
translation: machine
sourceHash: "d65f1b5bb0c8"
---

يفتح `ccam` دون وسائط لوحة تحكم بملء الشاشة عندما يكون الإدخال والإخراج القياسيان طرفية، ويطبع `ccam list` في غير ذلك. تعرض لوحة التحكم كل حساب مع استخدامه ومجموعاته، وتفعل ما تفعله الأوامر: التشغيل والإضافة وتسجيل الدخول وتغيير المجموعات وضبط قواعد المجلدات وسرد الجلسات الخلفية والفحص والإصلاح. ويغلقها تشغيل حساب، أو الاتصال بجلسة خلفية؛ أما إضافة حساب أو تسجيل الدخول فيسلّم الطرفية لحظةً ثم يعود إلى لوحة التحكم.

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

يقول السطر الموجود تحت أسماء التبويبات في جملة واحدة ما الغرض من التبويب الحالي. وفي الطرفية الأعرض يعرض تبويب Accounts أيضًا البريد الإلكتروني لكل حساب. وفي كل تبويب، تبيّن نهاية سطر التبويبات ما إذا كان [`auto_switch`](/features/next#auto-switch) مشغّلًا (على Linux وmacOS؛ أما Windows فيتجاهل هذا الإعداد). وعندما يكون مشغّلًا، يسرد تبويب Health ما لا يزال يحتاج إليه: Claude Code 2.1.214 أو أحدث، والمكوّن الإضافي لـ ccam مفعّلًا في كل حساب يسجّل الدخول باشتراك Claude.

في كل مرة تُفتح فيها لوحة التحكم، تسأل GitHub في الخلفية عن وجود إصدار أحدث من ccam؛ ولا ينتظر الإطار الأول ذلك أبدًا، وعندما تكون دون اتصال لا يظهر شيء. وإذا وُجد إصدار أحدث، أوصى به سطر تحت العنوان:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

يسأل `U` أولًا، ثم يحدّث بالطريقة التي يحدّث بها [`ccam update`](/features/update): يُتحقق من الإصدار مقابل مجاميع SHA-256 الخاصة به ويحل محل ccam في خطوة واحدة. ولوحة التحكم التي أنت فيها لا تزال الإصدار القديم، لذا اخرج وشغّل `ccam` من جديد. وعندما يتعذر التحديث من هنا، كما في نسخة ccam المثبّتة من حزمة `.deb` مثلًا، تعرض لوحة التحكم النصيحة نفسها التي يعرضها `ccam update` وتُبقي العرض قائمًا.

يُرسم الإطار الأول من `config.toml` وحده، في أقل بكثير من 100 ms؛ ثم يُملأ الاستخدام كلما أجاب حساب. ولكل حد شريطه الخاص: نافذة الخمس ساعات، والنافذة الأسبوعية لكل النماذج، والنافذة الأسبوعية الخاصة بـ Fable؛ ويعني `-` أن الحساب لا يملك تلك النافذة. وتتحول النسبة إلى الأصفر عند 80% وإلى الأحمر عند 100%. وإذا فشل الاستعلام، تُعرض القيم المخزنة مؤقتًا باهتةً. والنافذة التي أبلغ Claude Code بامتلائها تعرض 100% حتى إعادة ضبطها. ويظهر رمز الوصول المنتهي بعلامة `⚠ token expired` باهتة (يجدّده Claude Code عند التشغيل التالي)، بينما يعني `⚠ login expired` أن الحساب يحتاج إلى `ccam login`.

## التبويبات

| التبويب | ما يعرضه | المفاتيح |
| --- | --- | --- |
| Accounts | الرقم والاسم والبريد الإلكتروني وأشرطة الخمس ساعات والأسبوع وFable وآخر استخدام والمجموعات؛ `✔` الحساب الافتراضي، و`▸` حساب هذا المجلد، و`•` يعمل (بما في ذلك Claude Code الذي بدأ دون ccam: `claude` العادي، أو بيئة تطوير، أو الجلسات الخلفية وخدمتها الخلفية) | أدناه |
| Agents | ما يسرده [`ccam agents`](/features/agents): الجلسات الخلفية والجلسات المفتوحة لكل الحسابات، مع حالتها ومجلدها ووقت بدئها (أدناه) | `enter` اتصال (مثل `ccam attach`)، و`x` إيقاف، و`r` إعادة تحميل |
| Groups | كل مجموعة وما تشاركه وحساباتها، والحسابات التي لا تنتمي إلى أي مجموعة | `enter` تحرير، `n` جديدة، `e` إعادة تسمية، `x` إزالة |
| Folders | قواعد [`ccam folder`](/features/folder) والقاعدة التي تنطبق هنا، حتى عندما تصل إلى المجلد عبر رابط رمزي (وتُحفظ القواعد الجديدة بعد حل الروابط) | `n` استخدام حساب هنا، `enter` تغيير، `x` إزالة |
| Extensions | ما يعرضه [`ccam ext`](/features/ext)، كل إضافة مرة واحدة مع المجموعات التي تأتي منها، و`off:` الحسابات التي لم يُشغَّل فيها مكوّن إضافي، وخوادم الحساب الخاصة التي تنطبق على هذا المجلد فقط، و`This folder's project`، والمشكلات التي لا تخص أي حساب | `r` إعادة التحميل |
| Health | ما يبلّغ عنه [`ccam doctor`](/features/doctor)، على شكل صندوق لكل حساب (أدناه) | `f` إصلاح الصندوق المحدد، و`F` إصلاح الكل (`doctor --fix`)، و`↑`/`↓` للتنقل بين الصناديق |

يبدّل `←`/`→` أو `Tab` بين التبويبات، ويسرد `?` مفاتيح التبويب الحالي، ويحدّث `r` الاستخدام متجاوزًا ذاكرة التخزين المؤقت القصيرة ويعيد تحميل تبويبات Agents وExtensions وHealth، ويُنهي `q` أو `Esc` لوحة التحكم.

## المفاتيح في تبويب Accounts

| المفتاح | ما يفعله |
| --- | --- |
| `Enter` | يشغّل Claude Code بالحساب المحدد |
| `1`–`9` | ينتقل إلى الحساب الذي يحمل ذلك الرقم؛ اكتب `1` `2` بسرعة للحساب 12. ولا يشغّل الحسابَ إلا `Enter` |
| `↑`/`↓`، `k`/`j` | التنقل |
| `n` | يشغّل الحساب الذي سيختاره [`ccam next --auto`](/features/next) |
| `o` | يشغّل مع خيارات: `--continue` أو `--resume` أو وسائط أخرى، أو تشغيل تجريبي |
| `u` | تفاصيل الاستخدام: نوافذ الخمس ساعات والأسبوع والنماذج مع أوقات إعادة ضبطها ومصدر الأرقام |
| `a` | يضيف حسابًا (المعالج أدناه) |
| `l` | يسجّل الدخول من جديد |
| `s` | يختار مجموعات الحساب |
| `d` | يجعله الحساب الافتراضي |
| `f` | يستخدمه في المجلد الحالي |
| `e` | يعيد تسميته |
| `c` | يعيّن أمره المختصر أو يزيله |
| `x` | يزيله (يسأل أولًا، ويبدأ التحديد عند No) |

**مقاعد Team وEnterprise.** يبدّل الاختيار التلقائي للحساب (`ccam next` و`--fallback` و`ccam pick` والمفتاح `n` في لوحة التحكم، و[`auto_switch`](/features/next#auto-switch) بعد أن تشغّله) أيضًا بين مقاعد المؤسسة نفسها في Claude Team أو Enterprise، وهذا قد يخالف شروط Anthropic. اقرأ [الأسئلة الشائعة](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) وتحقّق من ذلك مع مؤسستك ومع Anthropic قبل أن تعتمد عليه.

## تبويب Health

يجري تبويب Health فحوص [`ccam doctor`](/features/doctor) عند فتحه ويرسم العرض نفسه: سطر أول فيه علامة لكل حساب، ثم صندوق لكل حساب. وكل صف في الصندوق مجال واحد جرى فحصه، أولها `Folder` و`Login` و`Share links`، وعليه علامة `✔` سليم، أو `ℹ` معلومة، أو `⚠` تحذير، أو `✘` خطأ، فترى فورًا أي حساب فيه مشكلة وأين. وتوضع النتائج التي لا تخص أي حساب في صندوق `General`.

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

ينتقل `↑`/`↓` من صندوق إلى آخر، ويكون للصندوق المحدد إطار بلون مميز. ويصلح `f` الصندوق المحدد وحده، كما يفعل `ccam doctor --fix`: ففي صندوق حساب يصلح ذلك الحساب وحده، وفي صندوق `General` يصلح النتائج التي لا تخص أي حساب (مسارات الإضافات والأوامر المختصرة)؛ أما بقية الحسابات فيُعاد فحصها فقط. ويسأل أولًا (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``)، والصندوق الذي لا شيء فيه للإصلاح يقول فقط `nothing to repair in <name>`. ويصلح `F` كل شيء دفعة واحدة، كما يفعل `ccam doctor --fix`، ويعيد `r` الفحص.

## تبويب Agents

يعرض تبويب Agents ما يطبعه [`ccam agents`](/features/agents): الجلسات الخلفية والجلسات المفتوحة لكل الحسابات، مع الحساب الذي تنتمي إليه كل منها. ويسأل Claude Code في كل الحسابات في آن واحد (`claude agents --json`، مع مهلة 15 ثانية لكل منها) عند أول فتح للتبويب، وعند `r`، وبعد إيقاف جلسة، وبعد أي تغيير في الإعدادات، ويقول `Asking every account's Claude Code for its sessions…` حتى تجيب كلها.

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

`State` هي حالة الجلسة الخلفية وما تفعله عمليتها، أو `open` للجلسة التفاعلية، كما في `ccam agents`. وإذا تعذّر الاستعلام من حساب ما، يعرض سطر `✘` تحت القائمة ما قاله Claude Code، مثل `lab` أعلاه، الذي أُوقف عرض الوكلاء فيه. وتُستبعد الجلسات الخلفية المنتهية (يسردها `ccam agents --all`)، وعندما تجيب كل الحسابات ولا تكون لأي منها جلسة، يقول التبويب ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` وفي الطرفية الضيقة يُسقط التبويب `Started`، ثم `Name`، ثم `Folder`، ويقصّ مسارات المجلدات الطويلة من اليسار.

يُغلق `Enter` على جلسة خلفية لوحة التحكم ويفتح الجلسة في هذه الطرفية بحسابها، كما يفعل `ccam attach <id>`؛ ولا تعود لوحة التحكم (في Claude Code، ينتقل `←` إلى عرض الوكلاء فيه و`Ctrl+Z` إلى الصدفة). والجلسة المفتوحة ليس لها معرّف، لذا يكتفي `Enter` و`x` بالقول إن المطلوب جلسة خلفية. ويسأل `x` أولًا، مع تحديد `No, keep it`: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` وعند الموافقة، يشغّل ccam الأمر `claude stop <id>` بذلك الحساب، ويقول ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` أو سبب تعذّر ذلك، ثم يقرأ القائمة من جديد.

## تغيير المجموعات

في تبويب Groups، يفتح `Enter` المحرر: العناصر على اليسار، والحسابات على اليمين، و`Tab` للتبديل بينهما، و`Space` للتحديد أو إلغائه. والعنصر الذي يحصل عليه أحد الحسابات بالفعل من مجموعة أخرى يظهر باهتًا مع السبب، وكذلك الحساب الذي سيحصل على عنصر مرتين؛ فالعنصر لا يأتي إلا من مجموعة واحدة. ويعرض `Enter` الخطة:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

إذا اختلف ملف عن نسخة المجموعة، يسألك ccam أولًا عمّا يفعله (الإبقاء على نسخة المجموعة، أو استخدام نسختك، أو الدمج). ثم يجري التغيير في الخلفية بالشيفرة نفسها التي يستخدمها [`ccam group`](/features/share): يُحفظ الإعداد أولًا، وتُنقل بيانات كل حساب بالترتيب، وتذهب الأصول إلى `.ccam-backup/`؛ والحساب الذي يغادر مشاركة MCP يحتفظ بخوادم المجموعة في `.claude.json` الخاص به.

والتغيير الذي سينقل بيانات حساب يعمل يُرفض مع ذكر اسم الحساب، حتى عندما يبدأ الحساب بينما الخطة مفتوحة، والخطة التي أصبحت قديمة لأن عملية ccam أخرى غيّرت المجموعات تُرفض مع `the groups changed while this plan waited; nothing was moved, make the change again`. وبينما تغيّر عملية ccam أخرى المجموعات، يقول التغيير الجديد أو إعادة التسمية أو الإزالة `another ccam process is changing the groups right now`؛ وبعد تغيير توقف في منتصفه، تنتظر حتى يُكمله `F` في تبويب Health (أو `ccam doctor --fix`).

## إضافة حساب

يفتح `a` المعالج نفسه الذي يستخدمه `ccam add` في الطرفية:

1. اسم الحساب.
2. طريقة تسجيل دخوله: اشتراك Claude الآن (في المتصفح)، أو [مفتاح API](/features/api)، أو اشتراك Claude لاحقًا.
3. لمفتاح API: المكان الذي يقرأ منه ccam المفتاح عند كل تشغيل (`env:NAME` أو `command:COMMAND`)، ونقطة النهاية، وطريقة إرسال المفتاح، والنموذج.
4. المجموعات التي ينضم إليها، مع تحديد `default` (وإذا لم توجد أي مجموعات بعد، ما ستشاركه المجموعة `default`).
5. أمر مختصر، إن أردت واحدًا.
6. ملخص لتأكيده.

يعود `Esc` خطوة واحدة إلى الوراء، ويغلق المعالج إذا كنت في الخطوة الأولى.

## إزالة حساب

يسأل `x` قبل أن يحدث أي شيء، مع تحديد No، ويوضح ما تفعله الإزالة: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` ولحساب مفتاح API لا يوجد تسجيل دخول لإبطاله، وعندما يحمل حساب آخر تسجيل الدخول نفسه لا يُسجَّل خروجه. والحساب المستورد يحتفظ بمجلده: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## الألوان

تتبع الألوان أدوار الألوان في واجهة الطرفية لدى Claude Code، فيبدو ccam منسجمًا بجانبه؛ وتأخذ أسماء المجموعات لونًا مشتقًا من أسمائها. ويسأل ccam الطرفية عن لون خلفيتها ويختار سمة فاتحة أو داكنة. ويمكنك تعيين سمة بنفسك في أعلى `config.toml`، فوق أول `[[group]]` أو `[[profile]]` أو `[[folder]]`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

السمات `-daltonized` مخصصة لعمى اللونين الأحمر والأخضر، والسمات `-ansi` تستخدم ألوان الطرفية الستة عشر الخاصة بها. ويوقف `NO_COLOR` الألوان.

## الأمان والحالات الخاصة

- تستدعي لوحة التحكم الدوال نفسها التي تستدعيها الأوامر، فتسري كل القواعد: الخطط قبل نقل البيانات، ورفض العمل على الحسابات التي تعمل، والنسخ الاحتياطية، ورموز تسجيل الدخول التي لا تغادر حسابها أبدًا.
- ليس في لوحة التحكم تشغيل تجريبي: يطبع `ccam --dry-run` دون اسم حساب كيف سيبدأ حساب هذا المجلد، ويعرض `o` على حساب ما تشغيلًا تجريبيًا لهذا التشغيل وحده.
- تُستعاد الطرفية قبل أن يبدأ Claude Code أو يُفتح تسجيل الدخول.
- إذا كانت الشاشة أصغر من أن تتسع لمربع حوار، تُحذف صفوف التبويب أولًا؛ وتُقطع الأسطر مع `…` بدلًا من لفّها. وعلى الشاشة القصيرة تُحذف بعد ذلك الأسطر الفارغة ووصف التبويب، فيبقى سطر تلميحات المفاتيح ظاهرًا. وتحت 40 عمودًا أو 8 صفوف لا تعرض لوحة التحكم إلا `Too small: ccam needs 40x8 (now 36x7)`.
- الضغط على `Ctrl+C` أثناء تطبيق تغيير أو إزالة أو تحديث أو إصلاح ينتظر اكتماله: ويقول سطر الحالة `ccam quits when it is done.`
- بعد أن يحدّث `U` برنامج ccam، لا يعود العرض عندما ترجع لوحة التحكم من إضافة حساب أو تسجيل دخول.

## صفحات ذات صلة

- [المشاركة عبر المجموعات](/features/share)
- [تشغيل Claude Code](/features/launch)
- [إدارة الحسابات](/features/accounts)
- [حدود الاستخدام](/features/usage)
