---
title: التبديل إلى الحساب التالي
description: يفتح ccam next برنامج Claude Code بحساب آخر؛ ويختار --auto الأوسع هامشًا، وينقل --continue المحادثة إليه، ويعيد --fallback تشغيل مهام -p.
updated: 2026-11-05
order: 8
translation: machine
sourceHash: "bc9837b7a6ae"
---

يفتح `ccam next` برنامج Claude Code بحساب آخر دون أن تحدد اسمه. استخدمه عندما يبلغ حساب ما حد الاستخدام. هناك أربع طرق للتبديل:

- يأخذ `ccam next` وحده الحساب الذي يلي آخر حساب استخدمته، ويأخذ `--auto` الحساب الأكثر سعة. ويتابع `--continue` آخر محادثة في هذا المجلد على الحساب الجديد، سواء تشارك الحسابان السجل أم لا.
- داخل Claude Code، يخبرك [المكوّن الإضافي لـ ccam](#inside-claude-code) عندما يوقف الحد جلسةً، ويسلّم `/ccam:switch` المحادثة إلى حساب آخر.
- في السكربتات التي تشغّل `claude -p`، يجري `ccam --fallback` التبديل بنفسه (راجع [إعادة التشغيل عند بلوغ الحد](#run-again-at-the-limit)).
- مع تشغيل [`auto_switch`](#auto-switch) (على Linux وmacOS)، ينقل ccam أيضًا الجلسة التفاعلية التي يوقفها Claude Code عند حدها.

**مقاعد Team وEnterprise.** يبدّل الاختيار التلقائي للحساب (`ccam next` و`--fallback` و`ccam pick` والمفتاح `n` في لوحة التحكم، و[`auto_switch`](/features/next#auto-switch) بعد أن تشغّله) أيضًا بين مقاعد المؤسسة نفسها في Claude Team أو Enterprise، وهذا قد يخالف شروط Anthropic. اقرأ [الأسئلة الشائعة](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) وتحقّق من ذلك مع مؤسستك ومع Anthropic قبل أن تعتمد عليه.

## الاستخدام

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

يطبع ccam التبديل بالشكل `→ <from> → <to>` على الخطأ القياسي، مع سبب الاختيار عند استخدام `--auto`، ثم يشغّل Claude Code بالحساب الجديد بالطريقة نفسها التي يعمل بها `ccam <name>`. و`<from>` هو آخر حساب استخدمته، أو مع `--continue` الحساب الذي تنتمي إليه المحادثة. تُمرَّر الوسائط التي تأتي بعد `--` إلى `claude` دون تغيير.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--auto` | يختار الحساب المسجَّل دخوله حسب الاستخدام (راجع أدناه)، متخطيًا الحسابات التي بلغت 100% | معطّل |
| `--max-usage <1-100>` | مع `--auto`: يفضّل الحسابات التي يقل استخدامها عن هذه القيمة | `auto_max_usage`، وإلا 100 |
| `--strategy headroom\|expiring\|balance` | مع `--auto`: أدنى استخدام في نافذة الخمس ساعات، أو أقرب إعادة ضبط لنافذة الخمس ساعات، أو أكبر حصة أسبوعية متبقية لكل يوم | `auto_strategy`، وإلا `headroom` |
| `--cross-org` | مع `--auto`: يختار أيضًا حسابات مؤسسة أخرى (راجع أدناه) | `cross_org`، وإلا معطّل |
| `--wait[=<duration>]` | مع `--auto`: عندما لا يكون لأي حساب هامش، ينتظر حتى يُعاد ضبط أولها، بحد أقصى هذه المدة. وتأتي القيمة بعد `=`: فيُرفض `--wait 2h` مع التلميح `--wait takes its value after '=': --wait=2h` | معطّل؛ و`5h` عند تحديده وحده |
| `-c`، `--continue` | يتابع آخر محادثة في هذا المجلد على الحساب الجديد (راجع أدناه) | معطّل |
| `--brief` | مع `--continue`: يبدأ محادثة جديدة بالحساب الجديد من ملخص ccam للمحادثة الأخيرة، بدلًا من متابعتها كلها ([راجع أدناه](#brief)) | معطّل |
| `-- <claude args>` | وسائط تُمرَّر إلى `claude` بعد وسائط ccam نفسه | لا شيء |
| `--dry-run` | مدعوم: يطبع `→ <from> → <to>` وسطر الأمر الذي كان سيشغّله، دون تشغيل Claude Code أو تسجيل عملية تشغيل | معطّل |

يعمل `--auto` و`--continue` معًا. ومع `--continue` توضع وسائط `--resume` قبل وسائطك.

توضع خيارات Claude Code الخاصة، مثل `--dangerously-skip-permissions` أو `--model`، بعد `--`. وعلى خلاف `ccam <name>` الذي يمرّر كل ما يلي الاسم، يقرأ `ccam next` وسائطه بنفسه، لذا يتوقف من دون `--` برسالة `ccam: unknown flag: --dangerously-skip-permissions`.

## كيف يعمل

كل تشغيل عبر ccam يكتب الوقت في `~/.local/state/ccam/last-used/<name>`. يقرأ `ccam next` هذه الملفات، ويتحقق من الحسابات المسجَّل دخولها، ثم يختار.

1. ينشئ ccam كائن `next.Candidate` واحدًا لكل حساب من `config.toml` (`cli.candidates`): الحساب، ووقت آخر تشغيل له (`launcher.LastUsed`)، وهل لديه تسجيل دخول صالح (`creds.Has`: تسجيل دخول إلى Claude في `.credentials.json`، أو على macOS في سلسلة المفاتيح، لم ينتهِ بالكامل ولم يمسحه Claude Code). والمجلد الذي يستخدم مفتاح API بدلًا من ذلك (تسجيل دخول Console، أو حساب أُضيف باستخدام `ccam add --api`) ليس له حدود خطة ولا يُختار أبدًا.
2. دون `--auto`، يرتّب ccam الحسابات حسب رقم الخانة ويجد الحساب الذي شُغّل مؤخرًا (`next.Rotate`). هذا هو الحساب الذي تنتقل منه. وإذا لم يُشغَّل أي حساب بعد، فهو صاحب أدنى خانة.
3. من هناك يتقدم بترتيب الخانات، ويعود من الخانة الأخيرة إلى الأولى، ويُرجع أول حساب مسجَّل دخوله.
4. يجهّز ccam الحساب المختار ويشغّله تمامًا كما يفعل `ccam <name>` (`launcher.LaunchWith`): يتحقق من روابط الحساب، ويطبّق مشاركة MCP، ويسجّل الجلسة، ويستبدل ccam بـ `claude` عبر `syscall.Exec`.

يبدأ المرور من آخر حساب استُخدم ويتوقف عند أول حساب مسجَّل دخوله بعده:

<!-- motion:next-rotate -->

شُغّل `work` مؤخرًا، لذلك يبدأ المرور بعده. ليس لدى `side` تسجيل دخول فيُتخطّى، ويُختار `lab`. ولو كان `lab` أيضًا غير مسجَّل الدخول، لعاد المرور إلى `personal` (`% len` في المخطط). لا يُرجع التدوير أبدًا الحساب الذي بدأت منه: فإن لم يوجد حساب آخر مسجَّل دخوله، يتوقف مع خطأ.

### --auto

يقرر `--auto` حسب الاستخدام بدلًا من ترتيب الخانات. يحتفظ ccam بالحسابات المسجَّل دخولها، ويطلب استخدامها كلها في الوقت نفسه بالطريقة التي يفعلها `ccam usage` (`Fetcher.FetchAll`)، ويختار وفق القواعد أدناه (`cli.pickAuto` و`usage.Policy.Choose()`):

<!-- motion:next-auto -->

يُحكم على الحساب بنافذتيه للخمس ساعات والأسبوع، وكذلك بالنافذة الأسبوعية للنموذج عندما يكون `--model` بين وسائط Claude Code بعد `--` (`opus` أو `sonnet` أو `haiku` أو `fable` في الاسم). ويوزّع ccam الحسابات على ثلاث فئات (`Policy.Choose()`):

1. **هامش**: كل نافذة أقل من `auto_max_usage` (`config.toml`، والقيمة الافتراضية 100). وتختار الاستراتيجية من بينها: يأخذ `headroom` (الافتراضي)، من بين الحسابات التي لا يزيد استخدامها في نافذة الخمس ساعات على أدنى استخدام بأكثر من 3 نقاط مئوية، الحساب الذي استهلك أقل قدر من أسبوعه، فيتقدّم حساب عند 12% على حساب عند 10% إذا كان قد استهلك قدرًا أقل من أسبوعه، أما الحساب عند 14% فلا يتقدّم عليه؛ ويأخذ `expiring` نافذة الخمس ساعات التي تُعاد ضبطها أولًا، فتُستهلك حصتها المتبقية قبل أن تضيع؛ ويأخذ `balance` الحساب الذي تتبقى له أكبر حصة أسبوعية لكل يوم حتى إعادة ضبط أسبوعه، محتسبًا عليه كل جلسة ccam قيد التشغيل، فلا تبقى الحصة الأسبوعية دون استخدام وتتوزع الجلسات المتوازية.
2. **تجاوز**: أقل من 100% لكن عند `auto_max_usage` أو فوقه. لا يُستخدم إلا عندما لا يكون لأي حساب هامش؛ ويأتي أولًا من بينها الحساب الذي كان `headroom` سيأخذه.
3. **ممتلئ**: نافذة بلغت 100%، أو استخدام غير معروف في نافذة الخمس ساعات. لا يُختار أبدًا.

تُحتسب النافذة 100% حتى إعادة ضبطها عندما يرفض Claude Code طلبًا فيها: يسجّل ذلك خطاف [المكوّن الإضافي لـ ccam](#inside-claude-code) (`~/.local/state/ccam/usage/<name>.limit.json`)، مهما قالت الأرقام المحفوظة. فلا يُختار مجددًا قبل إعادة ضبطه حسابٌ بلغ حده للتو، حتى عندما تكون أرقامه قديمة.

توضع الإعدادات الخاصة بكل الحسابات في أعلى `config.toml`، فوق أول `[[group]]` أو `[[profile]]` أو `[[folder]]`. والإعداد الواقع تحت أحدها ينتمي إلى ذلك الجدول: يحذّر ccam منه عند كل تشغيل، وترفض الأوامر التي تحفظ الإعدادات العمل حتى يُنقل.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**المؤسسات.** تبقى الاختيارات التلقائية بين حسابات من نوع الحساب الذي تغادره نفسه: مقاعد في مؤسسة Team أو Enterprise نفسها، أو الخطط الشخصية (Pro وMax)، فلا يتابع عملٌ بدأ في حساب الشركة على حساب شخصي ولا العكس. ويعرف ccam مؤسسة كل حساب من تسجيل دخوله (يسجّلها `ccam login`؛ وتُقرأ تسجيلات الدخول الأقدم من `.claude.json`)؛ والحساب الذي لا يستطيع ccam تحديد مؤسسته يمكن أن يُطابَق مع أي حساب. والمجلدان المسجَّل دخولهما بحساب Claude نفسه يتشاركان حدوده، لذلك لا يحل أحدهما محل الآخر أبدًا. ويرفع `--cross-org`، أو `cross_org = true` في `config.toml`، قاعدة المؤسسة، لكن بعد [`/ccam:switch`](#inside-claude-code) يبقى النقل إلى مؤسسة أخرى في انتظار موافقتك.

عندما لا يمكن اختيار أي حساب، يخرج ccam بالرمز 75 (`EX_TEMPFAIL`)، فيستطيع السكربت أن يميّز "حاول لاحقًا" من الفشل. ومع `--wait` ينام بدلًا من ذلك حتى تُعاد ضبط النافذة المانعة لأول حساب ممتلئ (بعد 30 ثانية من وقت إعادة الضبط، حين تُظهر نقطة النهاية النافذة الجديدة)، ويطبع `ccam: no account has room; waiting until 15:10`، ثم يختار من جديد؛ وعندما لا يذكر أي حساب متى تُعاد ضبطه، يعيد النظر بعد 10 دقائق. وإذا لم يتحرر أي حساب قبل انتهاء الانتظار، يخرج بالرمز 75.

يتبع الاستعلام قواعد `ccam usage`: يُعاد استخدام الأرقام التي عمرها أقل من 60 ثانية، وتنتهي مهلة استعلام كل حساب بعد 3 ثوانٍ، والحساب الذي لا يستجيب أو انتهى رمز الوصول الخاص به يُحكم عليه بأرقامه المحفوظة، ويتقدّم سجلُّ الحد عليها. والنافذة المحفوظة التي مضى وقت إعادة ضبطها تُحسب 0%، فلا يُتخطى حساب أُعيد ضبط نافذته منذ ذلك الحين. لا يجدّد `--auto` الرموز أبدًا. وعلى خلاف التدوير العادي، يمكن أن يختار `--auto` آخر حساب استخدمته إذا كان لا يزال الأكثر سعة.

### --continue

1. يبحث ccam في كل حساب عن آخر محادثة في هذا المجلد: من بين ملفات `<session>.jsonl` في `projects/<folder>/` (اسم المجلد هو مجلد العمل، بعد حل الروابط، وبعد تحويل كل حرف ليس حرفًا أو رقمًا إلى `-`، كما يسمّيه Claude Code)، الملف الذي فيه أحدث طلب أو رد. وتُتخطى الجلسات التي لا طلب فيها ولا رد، مثل جلسة لم تشغّل إلا `/mcp` أو `/usage` أو `/status`، أو ملف فارغ. وإذا تركت جلسة ملاحظة (حد استخدام أو `/ccam:switch`، راجع أدناه) ولم يكن في أي محادثة طلب أو رد بعدها بأكثر من دقيقتين، تُستخدم تلك الجلسة بعينها بدلًا من ذلك؛ ولا تُحتسب الصفوف التي يكتبها Claude Code عند `/exit`.
2. الحساب الذي تنتمي إليه المحادثة هو الحساب الذي تنتقل منه، ويُختار الحساب الجديد كما سبق، داخل المؤسسة نفسها. وعندما تتشارك عدة حسابات ذلك السجل، تنتمي المحادثة إلى الحساب الذي يقول `.claude.json` الخاص به إنه شغّل المحادثة آخر مرة (`lastSessionId`)، وإلا إلى الحساب الذي استُخدم آخر مرة.
3. إذا كان الحساب الجديد يشارك ذلك السجل (مجلد المجموعة نفسه)، يشغّل ccam `claude --resume <session id>`، أي تلك المحادثة بعينها (إذ كان `--continue` سيفتح أحدث محادثة في المجلد أيًّا كانت).
4. وإلا يشغّل `claude --resume <the conversation's file> --fork-session`: ينسخ Claude Code المحادثة كلها، بمعرّف جلسة جديد، إلى سجل الحساب الجديد نفسه، ويبقى الأصل كما هو. ومن ذلك الحين يجدها `--continue` في الحساب الجديد.
5. ما لم تعطِ `--agent` بعد `--`، يمرّر ccam أيضًا الوكيل الذي تعمل به المحادثة: إعداد `agent` في `.claude/settings.local.json` أو `.claude/settings.json` الخاصين بالمجلد أو في `settings.json` الخاص بالحساب، وإلا فالوكيل الذي يسمّيه آخر صف `agent-setting` في المحادثة. إذ لا يستعيد Claude Code الوكيل عندما ينسخ محادثة (الخطوة 4)، وكان سيتابع بوكيله الافتراضي، مع كل الأدوات. وعندما لا يتشارك الحساب الجديد وكلاء الحساب القديم ومكوّناته الإضافية وإعداداته (مجلدا `agents` و`plugins` ذاتهما و`settings.json` ذاته، عبر مجموعة)، يضيف ccam سطرًا بعد السطر الخاص بالنقل: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

يخبرك ccam بحجم المحادثة: فالحساب الجديد يقرؤها كلها من جديد مع رده الأول. ولا تُشارَك ذاكرات التخزين المؤقت للمطالبات (prompt caches) بين المؤسسات أبدًا، وكل خطة شخصية مؤسسة قائمة بذاتها، لذلك عندما ينتمي الحسابان إلى مؤسستين مختلفتين يضيف ccam `prompt caches are not shared between organizations, so the first reply reads it all again`. وعندما يكون الحساب القديم مسجّل الدخول إلى خوادم MCP لم يسجّل الحساب الجديد الدخول إليها، يذكرها ccam: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. ومشاركة `history` في مجموعة تتجنب النسخ، وتُبقي محادثة واحدة، وتشارك أيضًا نقاط الحفظ التي يستعيد منها `/rewind` الشيفرة (`file-history/`).

### --brief

يتابع `--continue` المحادثة كلها، ويقرؤها الحساب الجديد كلها من جديد مع رده الأول. أما `ccam next --continue --brief` فيبدأ بدلًا من ذلك محادثة جديدة، من ملخص يكتبه ccam للمحادثة القديمة دون أن يسأل نموذجًا:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

يحوي الملخص، بهذا الترتيب وفقط حيث يوجد ما يُقال:

1. عنوان المحادثة؛
2. مسار سجلها (transcript)، ليبحث فيه Claude عندما يحتاج إلى تفصيل؛
3. طلبك الأول، إذا كانت الطلبات أكثر من أربعة؛
4. الملخص الذي كتبه Claude Code عندما ضغط المحادثة آخر مرة؛
5. طلباتك الأربعة الأخيرة، مع أوامر `!` وأوامر الشرطة المائلة كما كتبتها؛
6. آخر رد لـ Claude؛
7. قائمة المهام؛
8. آخر 10 ملفات غيّرها Claude؛
9. `git status --short --branch` للمجلد، حتى 800 حرف.

وتحتفظ قائمة المهام بعشرة عناصر على الأكثر: تُحذف العناصر المنتهية أولًا، ثم الخطوات الأبعد، فتبقى الخطوة الحالية والخطوات التالية لها. وتُقصّ الأجزاء الطويلة من وسطها، ويبقى الملخص كله دون 16,000 حرف.

يحفظه ccam في `~/.local/state/ccam/brief/<session>.md`، ولا يقرؤه أحد غيرك، ويزيل `--brief` التالي الملخصات الأقدم من 7 أيام. ويبدأ Claude Code بمطالبة واحدة، ويرفق الملف الذي تذكره: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` والمطالبة التي تعطيها بعد `--` تحل محل الجملة الأخيرة (`Then: <your prompt>`)، وتتبعها خياراتك الأخرى بعد `--`.

ويُنهي Claude Code الإشارة (mention) عند `"` ويقصّها عند `#`، لذا عندما يحوي مسار الملخص أيًّا منهما، يشير إليه ccam بمساره نسبةً إلى مجلد العمل؛ وعندما يحوي ذلك المسار أحدهما أيضًا، تسمّي المطالبة الملف، فيقرؤه Claude بأداة Read، التي تسأل أولًا لأن الملف خارج مجلدات العمل (وفي تشغيل `-p` بالوضع اليدوي (Manual)، يرفض Claude Code القراءة ويتابع Claude دون الملخص).

يطبع `ccam --dry-run next --continue --brief` الملخص ولا يحفظ شيئًا. وتبقى المحادثة القديمة كما هي، ولا تعرف الجديدة إلا ما يقوله الملخص، لذا استخدم `--brief` عندما تهمّك المهمة أكثر من كل تفصيل في الطريق إليها. ودون `--continue` يتوقف ccam برسالة `--brief works with --continue`؛ وعندما يتعذر عليه قراءة المحادثة أو لا يجد فيها ما يلخصه، يتوقف ويقول إن `ccam next --continue` دون `--brief` يتابعها كاملة. ويتابع [الترحيل](#relay) و[`auto_switch`](#auto-switch) المحادثة كاملة دائمًا. أما خيار `--brief` الخاص بـ Claude Code فيوضع بعد `--`، مثل خياراته الأخرى.

## داخل Claude Code

يضيف [المكوّن الإضافي لـ ccam](/docs/install#claude-code-plugin) خطافين (Claude Code 2.1.78 أو أحدث):

- عندما يتوقف دورٌ عند حد استخدام، يدوّن المحادثة ويعرض إشعارًا على سطح المكتب في الطرفيات التي تدعمه (iTerm2 وWezTerm وWindows Terminal وkitty وGhostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  وهو يعرف كل رسالة يعدّها Claude Code حد استخدام: حدود الخمس ساعات والأسبوع والحد العادي وحدود النماذج (`You've hit your session limit` و`You've hit your limit` و`You've reached your …` لحصة نموذج)، وأرصدة الاستخدام وحدود الإنفاق، ومؤسسة أوقفت الحساب.

  وعندما يكون Claude Code قد كتب الطلب المرفوض في المحادثة، يسجّل الخطاف أيضًا النافذة الممتلئة حتى إعادة ضبطها ([سجل الحد](/features/usage#the-limit-record))، إلا في حالة إيقاف تفرضه مؤسسة، إذ يستطيع مسؤولوها رفعه في أي وقت. أما التقييد القصير من الخادم أو حد المعدّل لمفتاح API فليس حد استخدام، ويُترك كما هو، وكذلك الجلسات غير المراقَبة (تشغيلات `-p` والجلسات الخلفية). وإذا توقف وكيل فرعي أو وكيل سير عمل عند حد استخدام، عُدّ ذلك توقفًا لجلسته، حتى لو كانت المحادثة الرئيسية تنتظره فحسب: إذ يكتب Claude Code ذلك الطلب المرفوض في ملف الوكيل الخاص به داخل مجلد `subagents` التابع للجلسة، ويقرؤه الخطاف من هناك. وعندما يتوقف عدة وكلاء في الوقت نفسه، تحصل الجلسة على ملاحظة واحدة وإشعار واحد.

  وينتظر Claude Code 2.1.234 وما بعده أيضًا في الجلسة ويتابع بنفسه بعد إعادة ضبط نافذة الخمس ساعات؛ ويمنحك الإشعار خيار التبديل بدلًا من ذلك، ومع تشغيل [`auto_switch`](#auto-switch) (على Linux وmacOS) يبدّل ccam نيابةً عنك عندما يكون لدى حساب آخر متسع. وعندما لا يجوز لأي حساب آخر أن يتولى المحادثة، لأنه لا حساب آخر ينتمي إلى مؤسسة Team أو Enterprise نفسها، أو لا يوجد حساب آخر بخطة شخصية (ومع تشغيل `cross_org` يجوز ذلك لأي حساب Claude آخر)، ينتهي الإشعار بدلًا من ذلك بـ `You have no other personal plans, so ccam cannot move this conversation.` أو `You have no other accounts in Acme, so ccam cannot move this conversation.`
- يدوّن `/ccam:switch [account]` المحادثة وينهي الدور قبل إرسال أي طلب: اكتب `/exit`، ثم شغّل `ccam next --continue` لمتابعتها على الحساب الذي حددته (أو على الحساب التالي). وفي جلسة لا محادثة فيها بعد، يقول إنه لا يوجد ما يُتابَع. ولا يشغّله إلا ما تكتبه أنت: بدءًا من الإصدار 0.3.2 من المكوّن الإضافي، لا يستطيع Claude تشغيل `/ccam:switch`، والمهمة التي يجدولها Claude (CronCreate أو ScheduleWakeup أو `.claude/scheduled_tasks.json`) تمرّره إلى Claude نصًا عاديًا.

  عندما يكون الحساب الذي تحدده خارج مؤسسة المحادثة (مقعد في مؤسسة Team أو Enterprise أخرى، أو خطة شخصية لمحادثة على مقعد أو العكس، أو حساب مفتاح API في أي من الجانبين)، يسأل `ccam next --continue` أولًا، أيًّا كانت قيمة `cross_org`: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. ويسأل السؤال نفسه قبل حساب يختاره بنفسه عندما لا تحدد حسابًا، وهو حساب يتيح له `cross_org = true` أن يأخذه من مؤسسة أخرى: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. ولا ينقلها إلا `y` أو `yes`؛ وعلى Linux وmacOS، تُهمَل أولًا المفاتيح التي كُتبت قبل السؤال. وأي إجابة أخرى، أو عدم وجود طرفية يُطرح فيها السؤال، يُبقي المحادثة والملاحظة حيث هما: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (وإذا كان ccam هو من اختار الحساب، فإن `ccam next --continue` يختار من جديد). وكما في الاختيارات التلقائية، فإن الحساب الذي لا يستطيع ccam تحديد مؤسسته يناسب أي حساب.

تُحفظ الملاحظة في `~/.local/state/ccam/handoff/<account>/<session>.json`، ملف واحد لكل جلسة، فلا تكتب جلستان لحساب واحد فوق ملاحظات بعضهما أبدًا. ويُحتفظ بها لمدة 12 ساعة وتُستخدم مرة واحدة، حين توشك الجلسة التالية على البدء؛ وإذا فشل التشغيل بقيت للمحاولة التالية. وهي تخص المجلد الذي بدأت فيه الجلسة، حتى بعد تنفيذ `cd` في صدفة Claude Code. ويحدّث `ccam rename` الأسماء فيها؛ وإذا كان الحساب الذي تذكره قد أُزيل منذ ذلك الحين، يذكر ccam ذلك ويختار كالمعتاد.

### الترحيل (Relay)

مع `relay = true` في أعلى `config.toml`، يبقى ccam قيد التشغيل بجانب الجلسات التفاعلية بدلًا من أن يستبدل نفسه بـ Claude Code، وبعد `/ccam:switch` يتابع محادثة تلك الجلسة بنفسه (ولا يعمل بملاحظة تركتها جلسة أخرى): يُنهي الخطاف Claude Code، ويشغّل ccam الأمر `ccam next --continue` مع ملاحظة تلك الجلسة. ودون تسمية حساب، يأخذ الحساب الذي لديه أكبر هامش، محتسبًا نموذج المحادثة، وعندما يسمح `cross_org = true` بأن يكون ذلك حسابًا من مؤسسة أخرى، يطرح السؤال المذكور أعلاه قبل أن ينقل المحادثة إليه. وعلى Windows لا يستطيع الخطاف إنهاء Claude Code: اكتب `/exit`، فيتابعها ccam. ولا يُنهيه الخطاف كذلك لحساب خارج مؤسسة المحادثة: بل يقول `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`، ويطرح ccam السؤال المذكور أعلاه بعد أن تكتب `/exit`. وكلفة ذلك عملية ccam منتظرة واحدة لكل جلسة، ولهذا يكون معطّلًا افتراضيًا.

تحتفظ الجلسة المنقولة بما لا تستعيده متابعة المحادثة: خيارات مثل `--add-dir` و`--mcp-config` وقوائم الأدوات و`--settings` و`--plugin-dir` وخيارات موجّه النظام، ووضع الأذونات في المحادثة، ومستوى الجهد في الجلسة، والوكيل الذي تعمل به المحادثة (راجع [`--continue`](#continue)). أما `--model` و`--name` فيستعيدهما Claude Code نفسه من المحادثة. وأي شيء آخر (`--worktree` و`--debug` و`--remote-control`، والخيارات التي لا يعرفها ccam) يُذكر في `ccam: not carried over: …`. ولا يطرح التبديل أبدًا سؤال السجل المشترك، ويُهمَل ما كتبته أثناء تنفيذه؛ وسؤاله الوحيد هو الذي يُطرح قبل حساب من مؤسسة أخرى.

لا يتابع الترحيل المحادثة إلا عندما ينتهي Claude Code بالطريقة التي طلبها `/ccam:switch`: لا عندما تكتب `/exit` أولًا، أو تضغط Ctrl+C أثناء انتهائه، أو يتلقى ccam إشارة قطع الاتصال (hangup) أو إشارة أخرى، ولا عندما تختفي الطرفية. فتبقى الملاحظة عندها، ويتابعها `ccam next --continue` لاحقًا. والجلسة التي بدأت تحت ترحيل إصدار أقدم من ccam، أو الجلسة غير المراقَبة (تشغيلات `-p` والجلسات الخلفية)، لا يُنهيها الخطاف؛ بل تتلقى رسالة `/exit` بدلًا من ذلك. ويعمل `claude remote-control` دون ترحيل.

### `auto_switch`

مع `auto_switch = true` في أعلى `config.toml` (على Linux وmacOS)، ينقل ccam الجلسة إلى حساب آخر بنفسه عندما يوقفها Claude Code عند حد استخدام: يختار خطاف المكوّن الإضافي حسابًا آخر لديه متسع ويُنهي Claude Code، ثم يتابع ccam، الذي يعمل بجانب الجلسة، المحادثة هناك، كما يفعل [الترحيل](#relay) بعد `/ccam:switch`. وهو متوقف افتراضيًا، وميزة تجريبية.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

لا ينقل ccam الجلسة إلا إذا تحققت كل هذه الشروط:

1. شغّلها ccam في طرفية ويبقى قيد التشغيل بجانبها. يشغّل `auto_switch` [الترحيل](#relay)، لذا يتحقق هذا في الجلسات التفاعلية التي تبدؤها بعد ضبطه. والجلسة التي بدأت قبل ذلك، أو تحت ccam أقدم، تتلقى إشعارًا ينتهي بـ ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` ولا تُنقل أبدًا الجلسات غير المراقَبة (تشغيلات `-p` والجلسات الخلفية)، ولا `claude remote-control`، ولا Claude Code الذي بدأ دون ccam.
2. الحساب مفعَّل فيه [المكوّن الإضافي لـ ccam](/docs/install#claude-code-plugin) ويشغّل Claude Code 2.1.214 أو أحدث، الذي يمرّر معرّف عمليته إلى الخطاف (`CLAUDE_PID`).
3. كتب Claude Code الطلب المرفوض في المحادثة، مع وقت إعادة الضبط، لحد الخمس ساعات أو الحد الأسبوعي أو حد Opus أو Sonnet لذلك الحساب. أما حد الإنفاق، ونفاد أرصدة الاستخدام، والإيقاف الذي تفرضه مؤسسة، وحصة Fable (`You've reached your Fable limit.`)، وتقييد الخادم، وحد المعدّل لمفتاح API، فلا تنقل الجلسة أبدًا. أما حد وكيل فرعي أو وكيل سير عمل فيُحتسب، حتى لو كانت المحادثة الرئيسية تنتظره فحسب؛ وإذا كان حدًا أسبوعيًا لـ Opus أو Sonnet، فلا تنتقل الجلسة إلا إذا كانت المحادثة نفسها تعمل بذلك النموذج.
4. يعرف ccam خطة الحساب، من `.claude.json` الخاص به أو مما سجّله `ccam login`. وإلا ينتهي الإشعار بـ `ccam does not know personal's plan, so auto_switch leaves this session.`
5. الجلسة غير متصلة بـ Remote Control. فبحساب آخر لن تعيد الاتصال بالتطبيق الذي تقودها منه، لذا تبقى، ويقول الإشعار ذلك.
6. يُعاد ضبط الحد بعد أكثر من 15 دقيقة من الآن. وفي الانتظار الأقصر يترك ccam الجلسة لـ Claude Code: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. لم تكتب أي مطالبة منذ بلوغ الحد.
8. لدى حساب آخر متسع. يختار ccam كما يفعل [`--auto`](#auto) (`auto_max_usage` و`auto_strategy` ونموذج المحادثة)، من بين حسابات من نوع الحساب الذي تغادره نفسه: مقعد آخر في مؤسسة Team أو Enterprise نفسها، أو خطة شخصية أخرى (Pro، Max). ولا يُطبَّق هنا `cross_org` ولا `--cross-org`. ويجب على الحساب:
   - أن يكون مسجّل الدخول بتسجيل دخول يستطيع Claude Code تجديده؛
   - ألا يكون مسجّل الدخول إلى حساب Claude نفسه؛
   - ألا يستخدم مفتاح API؛
   - أن تكون له خطة يعرفها ccam؛
   - ألا تكون أرصدة الاستخدام فيه مفعّلة، بحسب ما يُظهره استعلام الاستخدام الخاص بـ ccam؛
   - أن يكون قد وثق بهذا المجلد في Claude Code؛
   - إذا كانت المحادثة تعمل بوكيل (يسمّيه إعداد `agent` أو تسمّيه المحادثة)، أن يتشارك وكلاء الحساب الذي تغادره: يجب أن يكون كلٌّ من مجلدَي `agents` و`plugins` الخاصين به وملف `settings.json` الخاص به هو نفسه تمامًا (عبر مجموعة)، أو غائبًا في كليهما.

   وعندما تستبعد هذه الشروط الثلاثة الأخيرة كل حساب لديه متسع، يسمّيها الإشعار: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`، أو `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.`، أو `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` والنافذة التي لها [سجل حد](/features/usage#the-limit-record) تُحتسب ممتلئة، والحساب الذي نُقلت إليه جلسة أخرى خلال آخر دقيقتين يأتي بعد الحسابات الأخرى. ويجب أن يتم الاختيار خلال 6 ثوانٍ تقريبًا؛ وعندما تستغرق استعلامات الاستخدام وقتًا أطول، تبقى الجلسة. ومع `usage_source = "claude"`، يقرر `auto_switch` من الأرقام المحفوظة وحدها، ولا يستطيع ccam أن يعرف من استعلام Claude Code أو من سطر الحالة هل أرصدة الاستخدام مفعّلة، لذا يفترض أنها متوقفة.

عندها يُنهي الخطاف Claude Code، ويعرض ccam، الذي يعمل بجانبه، إشعارًا على سطح المكتب في الطرفيات التي تدعمه، مثل `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` ثم يطبع السطر الأول من المخرجات في المثال أعلاه ويشغّل `ccam next --continue` لتلك الجلسة بالحساب المختار، الذي يطبع سطريه الآخرين ويشغّل Claude Code مع خيارات الجلسة ووضع الأذونات ومستوى الجهد والوكيل، كما يفعل [الترحيل](#relay). والرسالة الأولى في الجلسة الجديدة هي نص ثابت، وهو الشيء الوحيد الذي يرسله ccam نيابةً عنك (والنقل بعد `/ccam:switch` لا يرسل شيئًا):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

وبعد أن يبلغ وكيل فرعي حد الاستخدام تكون الرسالة هذه، لأن الوكلاء الفرعيين وسير العمل والأوامر الخلفية توقفت مع الجلسة القديمة:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

وإذا بلغ ذلك الحساب حده أيضًا، يحدث الشيء نفسه هناك. وتمنع سجلات الحد السلسلة من العودة إلى حساب قبل إعادة ضبطه، لذا عندما تمتلئ كل الحسابات تبقى الجلسة على آخرها.

يعرف الخطاف أن المحادثة تعمل بوكيل من إعداد `agent`، أو من الصف الذي يسمّي الوكيل ويكتبه Claude Code في المحادثة من حين لآخر وعند انتهاء الجلسة. فإذا لم يكن يسمّيه حتى الآن إلا `--agent` الخاص بالجلسة، ولم يكن الحساب المختار يتشارك الوكلاء، فلا يشغّل `ccam next --continue` شيئًا ويحتفظ بالملاحظة، ويقول ذلك إشعارٌ ثانٍ على سطح المكتب: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

وعندما لا يكون لدى أي حساب آخر متسع، لا ينتقل شيء: يتعامل Claude Code مع الحد كالمعتاد، فينتظر في الجلسة أو يعرض خياراته، ويذكر الإشعار أي حساب يتحرر أولًا، مثل `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` وللمقاعد يقول `None of your other accounts in Acme has room now`. ولإيقاف نقل بدأ، اضغط Ctrl+C أثناء إغلاق Claude Code: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

للنقل تكاليف. يُعاد تشغيل Claude Code، وهو ما يستغرق بضع ثوانٍ، وتتوقف مع الجلسة القديمة الصدفات الخلفية وخوادم التطوير والوكلاء التي شغّلها Claude. ويقرأ الحساب الجديد المحادثة كلها من جديد مع رده الأول، كما مع [`--continue`](#continue)، ويخبرك ccam بحجمها؛ ولا يستخدم النقل [`--brief`](#brief) أبدًا.

لا يبدّل `auto_switch` أبدًا قبل أن يُبلغ Claude Code عن حد، ولا يراقب استخدامك في الخلفية أبدًا، ولا ينسخ تسجيل دخول أو يجدّده أو يمرّره أبدًا. وعلى Windows يُتجاهل، ويقول `ccam doctor` ذلك؛ وهناك، مع `relay = true`، يظل `/ccam:switch` و`/exit` ينقلان المحادثة إلى حساب آخر.

**شروط Anthropic.** في مقاعد Team وEnterprise، ينقل `auto_switch` عملك إلى مقعد آخر في المؤسسة نفسها. يقول مركز المساعدة لدى Anthropic: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (يُخصَّص كل مقعد في خطة Team لشخص واحد محدد بالاسم، ولا يجوز أن يستخدمه إلا ذلك الشخص) و"Users can't route their work through another user's account to use that user's remaining usage" (لا يجوز للمستخدمين تمرير عملهم عبر حساب مستخدم آخر للاستفادة من الاستخدام المتبقي لذلك المستخدم) ([خطط Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)؛ وبالمثل في [خطط Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)).

أما في الخطط الشخصية (Pro وMax)، فلا تقول شروط Anthropic هل يجوز لحساباتك أن يتولى بعضها عن بعض تلقائيًا عند بلوغ الحد. وتمنع [شروط المستهلك](https://www.anthropic.com/legal/consumer-terms) الوصول إلى الخدمات "through automated or non-human means, whether through a bot, script, or otherwise" (عبر وسائل آلية أو غير بشرية، سواء عبر روبوت أو سكربت أو غير ذلك) ما لم تسمح Anthropic بذلك صراحةً، والنقل يرسل إحدى الرسائل الثابتة أعلاه نيابةً عنك. وتُدرج [سياسة الاستخدام](https://www.anthropic.com/legal/aup) التي تدخل حيز التنفيذ في 12 نوفمبر 2026 عبارة "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (تجاوز القدرات أو القيود أو الضوابط المضمّنة في منتجاتنا أو التدخل فيها عمدًا) ضمن الأمور التي تحظرها.

ولا يدّعي ccam أن `auto_switch` يتوافق مع هذه الشروط. وتشغيله قرارك أنت، وكذلك المخاطرة: فقد تعلّق Anthropic حسابًا، أو وصول مؤسسة Team أو Enterprise، عندما ترى أن الشروط انتُهكت.

وما دام `auto_switch` مفعّلًا، يحذّر `ccam doctor` من الحساب الذي لا يملك المكوّن الإضافي لـ ccam ومن إصدارات Claude Code الأقدم من 2.1.214، ويكرر للمقاعد قاعدة المقاعد لدى Anthropic ([الفحوص](/features/doctor)).

## إعادة التشغيل عند بلوغ الحد

`ccam next` مخصص لجلسة تفاعلية توقفت. أما التشغيل غير التفاعلي، أي `claude -p` في سكربت أو مهمة مجدولة، فلا يستطيع انتظارك، لذا يتيح وضع `--fallback` قبل اسم الحساب أن يبدّل ccam نيابةً عنه:

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

1. لا يعمل `--fallback` إلا مع `-p` أو `--print`؛ وإلا يتوقف ccam ويشير إلى `ccam next --auto --continue`.
2. يُقرأ الإدخال القياسي المنقول عبر أنبوب مرة واحدة، بحد أقصى 10 MiB، ويُعطى لكل محاولة. ومثل Claude Code، ينتظر ccam بدء وصوله 3 ثوانٍ على الأكثر (إلا مع `--input-format stream-json`)؛ وإذا بقي الأنبوب مفتوحًا لكنه لم يرسل شيئًا، لا يحصل التشغيل على أي إدخال ويذكر ccam ذلك في الخطأ القياسي.
3. يشغّل ccam `claude` عمليةً فرعية بدلًا من `exec`، ويحجز مخرجه القياسي، ويمرّر الخطأ القياسي فور وصوله.
4. عندما يخرج `claude` بالرمز 0، أو يفشل لسبب آخر، يُمرَّر المخرج المحجوز ورمز الخروج دون تغيير. وكذلك التشغيل الذي أنهته إشارة (رمز خروج أعلى من 128): فقد طُلب منه التوقف، ولا يُشغَّل مجددًا.
5. يُعدّ الفشل بلوغًا لحد الاستخدام عندما تقول رسالة Claude Code ذلك، بأي من رسائل حد الاستخدام التي يعرفها [الخطاف](#inside-claude-code)، حتى داخل مخرجات `--output-format json` أو `stream-json`، أو عندما يُظهر استخدام الحساب أن نافذة الخمس ساعات أو الأسبوع أو `--model` بلغت 100% (بما في ذلك سجل الحد).
6. عندها يُهمَل مخرج تلك المحاولة ويُشغَّل الأمر نفسه بحساب آخر مسجّل الدخول: صاحب أدنى استخدام في نافذة الخمس ساعات، ثم الحسابات التي يتعذر الاستعلام عن استخدامها، بترتيب أرقامها. وتُتخطى الحسابات التي بلغت 100%، ولا تُجرَّب أبدًا الحسابات التي تسجّل الدخول بمفتاح API، ولا يُجرَّب أي حساب أكثر من مرة واحدة.
7. إذا كان بين الوسائط `-c` أو `--continue` أو `-r` أو `--resume`، فلا تُجرَّب إلا الحسابات الموجودة في مجموعة السجل نفسها، حتى يمكن العثور على المحادثة.
8. لا تُجرَّب إلا حسابات المؤسسة نفسها، كما مع `--auto` (ويرفع `cross_org = true` هذا القيد).

يُختار الحساب التالي بقواعد `--auto` نفسها (`auto_max_usage` و`auto_strategy` ونافذة `--model`). وعندما لا يتبقى أي حساب، يطبع ccam مخرج المحاولة الأخيرة، ويذكر أنه لم يعد لدى أي حساب مسجّل الدخول هامش متبقٍ، ويخرج بالرمز 75. ومع `--wait` قبل اسم الحساب (`ccam --fallback --wait=2h work -p …`)، ينتظر بدلًا من ذلك حتى يُعاد ضبط أول حساب، ثم قد يجرّب كل الحسابات من جديد؛ وهو ينام دائمًا مرة واحدة على الأقل أولًا، فلا يُشغَّل فورًا من جديد حساب بلغ حده للتو. وتحدّ المدة إجمالي وقت الانتظار في التشغيل، لا كل انتظار على حدة. وإذا كُتبت بمسافة (`--wait 30m`)، تُرفض مع التلميح نفسه الذي يعطيه `ccam next`. ويعمل `--` في البداية قبل وسائط Claude Code كما في أي تشغيل (`ccam --fallback work -- -p …`). ويطبع `ccam --dry-run --fallback work -p …` أمر المحاولة الأولى.

لا تُبدَّل الجلسات التفاعلية من تلقاء نفسها إلا إذا شغّلت [`auto_switch`](#auto-switch). ومع `ccam statusline install`، يضيف سطر الحالة `→ ccam next --auto --continue` بمجرد أن يبلغ الاستخدام في نافذة الخمس ساعات أو النافذة الأسبوعية 90% ([حدود الاستخدام](/features/usage)).

## الأمان والحالات الخاصة

- لا تُختار الحسابات التي ليس لها تسجيل دخول أبدًا. وإن لم يوجد حساب آخر مسجَّل دخوله، يتوقف `ccam next` وحده مع `no logged-in account other than <name>`.
- مع `--continue` وعدم وجود محادثة في المجلد: `no account has a conversation in <folder> to continue`.
- عندما لا تترك قاعدة المؤسسة أي حساب: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- مع `--auto` وعدم وجود حساب مسجَّل دخوله: `no logged-in account to pick`. وعندما يكون كل حساب ممتلئًا أو استخدامه غير معروف: `` no account with known usage and headroom (see `ccam usage`) ``. وتنتهي الحالتان برمز الخروج 75.
- لا تعمل `--max-usage` و`--strategy` و`--wait` و`--cross-org` إلا مع `--auto`، وتُرفض القيمة خارج النطاق 1-100 أو الاستراتيجية غير المعروفة.
- ينسخ التفرّع (fork) المحادثة، لا حالة الحساب الأخرى: فنقاط الحفظ في `file-history/` الخاصة بالحساب الأصلي لا تنتقل معها ما لم يكن السجل مشتركًا.
- يرسل `--auto` طلبات شبكة حتى مع `--dry-run`. وهو يقرأ فقط، ولا يكتب ملفات الرموز أبدًا.
- يبدأ الحساب المختار عبر المسار نفسه الذي يسلكه `ccam <name>`، فتنطبق الفحوص نفسها، ومنها التحذير عندما تعمل بالفعل في هذا المجلد جلسة Claude Code أخرى تشارك السجل نفسه.
- لا يُكتب `last-used` إلا عندما يبدأ Claude Code فعلًا. ويتركه `--dry-run` كما هو، فلا يدفع التشغيل التجريبي التدوير إلى الأمام.
- يحجز `--fallback` المخرج القياسي حتى تنتهي المحاولة، لذا يصل المخرج المتدفق مثل `--output-format stream-json` دفعة واحدة في النهاية. والعميل الذي يرسل مدخلات `stream-json` متدفقة وينتظر الردود قبل أن يرسل المزيد لا يستطيع استخدام `--fallback`. والملفات التي غيّرتها المحاولة الأولى قبل أن تبلغ الحد تبقى متغيّرة؛ وتبدأ المحاولة التالية الطلب من جديد.

## صفحات ذات صلة

- [حدود الاستخدام](/features/usage)
- [حدود الاستخدام والتبديل بين الحسابات](/docs/usage-limits)
- [مشاركة خوادم MCP والإعدادات والسجل](/docs/sharing)
- [ماذا تفعل عندما يبلغ Claude Code حد الاستخدام](/guides/usage-limit-workflow)
