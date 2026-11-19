---
title: مرجع الأوامر
metaTitle: مرجع أوامر ccam
description: كل أوامر ccam وخياراتها، من ccam add وccam <name> إلى ccam usage وccam next وccam doctor.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: "28e1ad1dcbad"
---

ضع `--dry-run` قبل أمر لترى ما سيفعله دون أن يفعله، مثلًا `ccam --dry-run work`. يعمل ذلك عند تشغيل Claude Code (`ccam <name>` و`ccam run` و`ccam next` و`ccam attach` و`ccam ccusage`) ومع `group` و`migrate` و`folder` و`shortcut` و`pick` و`update` و`uninstall` و`doctor --fix`. أما الأوامر التي تقرأ فقط، مثل `list` و`usage` و`share` و`ext`، فتقبله ولا تغيّر شيئًا في كل الأحوال. وترفضه بقية الأوامر بدلًا من أن تعمل، ومنها `group rename`.

<!-- commands -->

## تشغيل Claude Code

يشغّل `ccam <name|number> [arguments]` برنامج Claude Code بذلك الحساب. لا يقرأ ccam الوسائط التي تأتي بعد الاسم؛ بل يمرّرها إلى `claude` بالترتيب نفسه، ويتصرف الإدخال القياسي ورمز الخروج والإشارات تمامًا كما مع `claude` نفسه. وضع `--` بعد الاسم اختياري.

أما `ccam` مع وسائط ولكن دون اسم حساب، مثل `ccam -p "summarise"`، فيستخدم الحساب المعيّن للمجلد الحالي، وإلا فالحساب الافتراضي. عيّنهما باستخدام `ccam folder <name> [folder]` و`ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

الخيار `--fallback` قبل اسم الحساب مخصص لتشغيلات `-p`/`--print`. عندما يبلغ الحساب حد استخدامه، يشغّل ccam الأمر نفسه مجددًا بالحساب المسجّل الدخول الذي لديه أكبر هامش، ولا يطبع إلا مخرج التشغيل الذي اكتمل. راجع [الحسابات حسب المجلد](/features/folder) و[إعادة التشغيل عند بلوغ الحد](/features/next#run-again-at-the-limit).

**مقاعد Team وEnterprise.** يبدّل الاختيار التلقائي للحساب (`ccam next` و`--fallback` و`ccam pick` والمفتاح `n` في لوحة التحكم، و[`auto_switch`](/features/next#auto-switch) بعد أن تشغّله) أيضًا بين مقاعد المؤسسة نفسها في Claude Team أو Enterprise، وهذا قد يخالف شروط Anthropic. اقرأ [الأسئلة الشائعة](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) وتحقّق من ذلك مع مؤسستك ومع Anthropic قبل أن تعتمد عليه.

يشغّل الأمر المختصر مثل `ccwork` حسابه بالطريقة نفسها: يفعل `ccwork --resume` ما يفعله `ccam work --resume`. راجع [الأوامر المختصرة](/docs/commands#shortcuts) أدناه.

## إضافة الحسابات وإزالتها

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

عند تشغيل `ccam import` دون تحديد مجلد، يبحث عن `~/.claude` و`~/.claude-*` ويسألك عن كل واحد منها؛ ويُبلَّغ عن المجلد الذي يفشل بالرسالة `Skipped <path>: <reason>` وتُسجَّل المجلدات الأخرى. وإذا كان اسم مجلد مستخدمًا بالفعل، فاستورد ذلك المجلد وحده باستخدام `ccam import <folder> --name <name>`. تُسجَّل المجلدات المستوردة في أماكنها؛ ولا ينقل ccam أي تسجيل دخول أو ينسخه أبدًا.

يسجّل `ccam login --expired` أيضًا دخول الحسابات التي مسح Claude Code تسجيل دخولها. ويسجّل `ccam remove` خروج الحساب المُدار (الذي أنشأ ccam مجلده، بخلاف الحساب المسجَّل باستخدام `ccam import`) قبل أن يذهب مجلده إلى سلة المهملات، إلا إذا كان حساب آخر يحمل تسجيل الدخول نفسه. وإزالة حساب مستورد تُخرجه من مجموعاته وتترك في مجلده نسخة من كل عنصر مشترك. وعلى macOS، يحتفظ `ccam rename` باسم المجلد، لأن سلسلة المفاتيح تحفظ تسجيل الدخول تحت مسار المجلد.

يضيف `ccam add --api` حسابًا يستخدم مفتاح API بدلًا من تسجيل الدخول إلى Claude؛ ويقرأ ccam المفتاح عند كل تشغيل ولا يخزّنه أبدًا. راجع [حسابات مفاتيح API](/features/api).

ترفض الأوامر `ccam rename` و`ccam remove` وأوامر `ccam group` تغيير حساب ما دام Claude Code يعمل بذلك الحساب، أيًّا كانت طريقة تشغيله (بما في ذلك الجلسات الخلفية). أنهِ الجلسة أولًا. وترفض أيضًا ما دام تغيير في المجموعات لم يكتمل: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

يطبع `ccam list --json` كل حساب مع رقم خانته واختصاره وبريده الإلكتروني ومجلده. وتسرد المصفوفة `groups` مجموعات الحساب، ويربط الكائن `share` كل عنصر (`mcp` و`settings` و`memory` و`extensions` و`history`) بالمجموعة التي يأتي منها، أو بالقيمة `local`.

## الأوامر المختصرة

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

الأمر المختصر أمر تختار اسمه بنفسك، ويشغّل Claude Code بحساب واحد. ومثل `ccam work`، يمرّر كل وسيط إلى Claude Code دون تغيير، لذلك يعمل `ccwork --resume` و`ccwork -p "hi"` كما تتوقع. وهو رابط إلى ccam في `~/.local/bin`.

الأوامر المختصرة اختيارية. يسألك معالج `ccam add` هل تريد إنشاء واحد، وينشئ `ccam add <name> --shortcut <command>` واحدًا دون سؤال، ويعرض `ccam list` الأمر المختصر لكل حساب، أو `-` إذا لم يكن له أمر.

تستخدم الأسماء الأحرف a-z و0-9 و`.` و`-` و`_`. ولا يجوز أن يكون الاسم أمرًا من أوامر ccam، أو أمرًا مختصرًا لحساب آخر، أو اسم برنامج آخر موجود في `PATH`، ولا ملفًا في `~/.local/bin` لم ينشئه ccam، مثل نسختك الخاصة من ccam أو رابط إليه. والأمر المختصر المعدَّل يدويًا في `config.toml` الذي يخالف هذه القواعد، أو يكون `claude` أو `ccam`، أو يُعطى لحسابين، يجعل تحميل `config.toml` يفشل مع ذكر الحساب.

يحتفظ `ccam rename` باسم الأمر المختصر، ويحذف `ccam remove` الأمر المختصر مع الحساب. ويعيد `ccam shims install` إنشاء روابط كل الحسابات، مثلًا بعد نقل ccam؛ والملف الذي يعترض أمرًا مختصرًا واحدًا لا يوقف الأوامر الأخرى.

## المشاركة

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

تشارك المجموعة عناصر بين حساباتها؛ ويمكن أن يكون الحساب في عدة مجموعات، لكن كل عنصر يأتي من مجموعة واحدة فقط. راجع [المشاركة عبر المجموعات](/features/share).

## نسخ إعداد موجود

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

ينسخ `ccam migrate` الملفات `settings.json` و`CLAUDE.md` والإضافات وخوادم MCP الخاصة بالمستخدم من مجلد Claude Code إلى مجموعة، فتحصل عليها كل الحسابات التي تشارك هذه العناصر. يُقرأ المجلد فقط.

- تنسخ الخيارات `--settings` و`--memory` و`--extensions` و`--mcp` و`--history` العناصر التي تحددها فقط. ودونها تتيح لك قائمة اختيار في الطرفية أن تختار، مع تحديد كل شيء ما عدا سجل المحادثات؛ وفي غير ذلك يُنسخ كل شيء ما عدا سجل المحادثات.
- يعمل `--mcp` أيضًا عندما يكون `~/.claude` فارغًا أو غير موجود لكن `~/.claude.json` يحتوي على خوادم.
- يختار `-g` أو `--group` المجموعة.
- ينسخ `-y` أو `--yes` دون سؤال ويحتفظ بنسخة المجموعة عندما يختلف ملف.

لا تُنسخ تسجيلات دخول MCP، لأن تسجيل الدخول المنسوخ يُخرج حامليه الآخرين؛ وتذكر الخطة الحسابات التي تسجّل الدخول باستخدام `/mcp`. راجع [نسخ إعداد موجود](/docs/sharing#copy-an-existing-setup).

## خوادم MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

ينقل `ccam mcp import <account>` الخوادم التي لا يملكها إلا ذلك الحساب إلى مجموعته، بما في ذلك الخوادم التي يملكها لمجلد واحد فقط (`claude mcp add` العادي)، فتُحمَّل بعدها في كل مجلد. ويُتخطّى الخادم الذي تستخدم المجموعة اسمه بالفعل لتعريف مختلف؛ أضف `--force` لاستبدال نسخة المجموعة، ويحتفظ ccam بنسخة احتياطية منها أولًا.

## الاستخدام والتبديل

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

- تختار `next --auto` و`--fallback` و`ccam pick` بالطريقة نفسها: الحسابات التي يقل استخدامها عن `auto_max_usage` (`config.toml`، والقيمة الافتراضية 100) أولًا، وفق `auto_strategy` (`headroom` أو `expiring` أو `balance`)، مع احتساب النافذة الأسبوعية لـ `--model`. ولا تختار أبدًا حسابًا بلغ 100% (والنافذة التي أبلغ Claude Code بامتلائها تُحتسب 100% حتى إعادة ضبطها) ولا تسجيل دخول بمفتاح API، وتبقى داخل مؤسسة الحساب الذي تبدأ منه (`--cross-org` أو `cross_org = true` للبحث أبعد من ذلك).
- يأخذ `--wait` قيمته بعد `=` (`--wait=2h`)؛ وإذا كُتبت بمسافة يُرفض مع تلميح. وعندما يتعذر اختيار أي حساب تخرج بالرمز 75، أو تنتظر، مع `--wait`، حتى يُعاد ضبط أول حساب.
- ينقل `ccam next --continue` المحادثة حتى بين حسابات لا تتشارك السجل: فينسخها Claude Code حينئذٍ إلى الحساب الجديد. ومع `--brief`، يبدأ الحساب الجديد بدلًا من ذلك محادثة جديدة، من ملخص يكتبه ccam للمحادثة الأخيرة دون نموذج ([`--brief`](/features/next#brief)).
- على خلاف `ccam <name>`، يقرأ `ccam next` خياراته بنفسه، لذا توضع خيارات Claude Code بعد `--`؛ ومن دونه يتوقف ccam برسالة `unknown flag`.
- يشغّل `ccam agent` مهمة باستخدام `claude -p` في هذا المجلد بحساب غير الحساب الذي يُستدعى منه، بوضع الأذونات ومستوى الجهد الخاصين بالجلسة المستدعية، ويطبع آخر رد من Claude. وعند بلوغ ذلك الحساب حده، تستمر المحادثة بحساب آخر ([وكلاء فرعيون على حسابات أخرى](/features/agent)).
- داخل Claude Code، يسلّم الأمر `/ccam:switch [account]` من المكوّن الإضافي المحادثة إلى حساب آخر، و`ccam hook` هو ما تشغّله خطافاته.
- مع `auto_switch = true` في `config.toml` (على Linux وmacOS؛ وهو متوقف افتراضيًا)، يُنهي خطاف الحد أيضًا الجلسة التفاعلية التي شغّلها ccam ويتابعها على حساب آخر لديه متسع، لكن فقط على حساب من النوع نفسه: مقعد آخر في المؤسسة نفسها، أو خطة شخصية أخرى ([`auto_switch`](/features/next#auto-switch)).

راجع [التبديل إلى الحساب التالي](/features/next) و[اختيار حساب للسكربتات](/features/pick) و[وكلاء فرعيون على حسابات أخرى](/features/agent) و[الجلسات الخلفية لكل الحسابات](/features/agents).

## الصيانة

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

يتحقق `ccam update` من مجموع SHA-256 الخاص بالإصدار قبل أن يستبدل ccam، ويثبّت `--version v0.2.0` إصدارًا محددًا. أما نسخة ccam المثبّتة من حزمة `.deb` فتُحدَّث باستخدام apt. راجع [التحديث](/docs/install#update).

يعرض `ccam uninstall` ما سيزيله ويسأل أولًا: كل الحسابات كما يفعل `ccam remove`، والأوامر المختصرة، ومجلدات ccam، والأمر `ccam`. أما المجلد المستورد فيبقى في مكانه. راجع [إلغاء التثبيت](/docs/install#uninstall).

## الملفات

| المسار | المحتوى |
| --- | --- |
| `~/.config/ccam/config.toml` | الحسابات والمجموعات والإعدادات (`theme` و`usage_source` و`relay` و`auto_switch` و`cross_org` …). توضع الإعدادات في الأعلى، فوق أول `[[group]]` أو `[[profile]]` أو `[[folder]]`: فالإعداد الواقع تحت جدول ينتمي إلى ذلك الجدول، لذا يحذّر ccam منه عند كل تشغيل وفي `ccam doctor`، وترفض الأوامر التي تحفظ الإعدادات العمل حتى يُنقل. ولا تُحفظ التعليقات عندما يحفظ ccam الملف |
| `~/.local/share/ccam/profiles/<name>/` | مجلد `CLAUDE_CONFIG_DIR` لكل حساب |
| `~/.local/share/ccam/groups/<group>/` | الإعدادات وخوادم MCP والإضافات والسجل المشتركة |
| `~/.local/state/ccam/` | ذاكرة التخزين المؤقت للاستخدام، وسجلات الحد، والأقفال، والملاحظات، وحالات أخرى (راجع القائمة أسفل الجدول) |
| `~/.local/bin/<shortcut>` | أوامر مختصرة مثل `ccwork`، وهي روابط إلى ccam |

يحتوي مجلد الحالة على:

- ذاكرة التخزين المؤقت للاستخدام وسجلات الحد (`usage/<name>.json` و`usage/<name>.limit.json`)؛
- أوقات آخر استخدام؛
- علامات الجلسات والأقفال (`locks/`)؛
- الأوامر المختصرة التي أنشأها ccam (`shortcuts.json`)؛
- تغيير المجموعات الذي لم يكتمل (`group-change.json`)؛
- الملاحظات التي يتركها `/ccam:switch` وخطاف الحد (`handoff/<account>/<session>.json`)؛
- متى اختار `ccam pick` أو `auto_switch` كل حساب آخر مرة (`picks/<name>`؛ ويُحتسب الاختيار لمدة دقيقتين)؛
- ملخصات `--brief`، المحفوظة 7 أيام (`brief/<session>.md`).

ضبط `XDG_CONFIG_HOME` أو `XDG_DATA_HOME` أو `XDG_STATE_HOME` يغيّر مواقع هذه المجلدات. وعلى Windows، يكون `~` هو `%USERPROFILE%` وتكون الاختصارات `<shortcut>.exe`.
