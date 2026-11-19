---
title: وكلاء فرعيون على حسابات أخرى
description: يشغّل ccam agent جزءًا من مهمة Claude Code بحساب آخر من حساباتك ويعيد ردّ Claude، فتستفيد مهمة طويلة واحدة من استخدام عدة حسابات.
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

يعمل الوكلاء الفرعيون الخاصون بـ Claude Code بحساب الجلسة التي تشغّلهم، فيستهلكون استخدام ذلك الحساب. يشغّل `ccam agent "<task>"` مهمة باستخدام `claude -p` في المجلد نفسه بحساب آخر من حساباتك ويطبع آخر رد من Claude. وعند استدعائه من جلسة Claude Code، يعمل مثل وكيل فرعي يستهلك استخدام حساب آخر، فيمكن لمهمة طويلة أن تستمر بعد النقطة التي كان سينفد عندها حساب الجلسة نفسه.

**مقاعد Team وEnterprise.** يختار `ccam agent` الحساب بنفسه. وداخل مؤسسة واحدة في Claude Team أو Enterprise يشغّل المهمة على مقعد آخر، وهذا قد يخالف شروط Anthropic. اقرأ [الأسئلة الشائعة](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) وتحقّق من ذلك مع مؤسستك ومع Anthropic قبل أن تعتمد عليه.

## الاستخدام

يشغّله Claude Code عبر أداة Bash الخاصة به. وعندما يكون [المكوّن الإضافي لـ ccam](/docs/install#claude-code-plugin) مثبّتًا، اطلب ذلك بكلمات بسيطة، مثل "استخدم حساباتي الأخرى لأجزاء من هذه المهمة"، أو اكتب ذلك في ملف CLAUDE.md الخاص بك. وتطلب مهارة المكوّن الإضافي من Claude Code أن يشغّل كل وكيل في الخلفية، وأن يكتب مهمة قائمة بذاتها، وأن يعطي الوكلاء الذين يعدّلون الملفات جزءًا خاصًا بهم أو git worktree خاصًا بهم.

```sh
ccam agent "find why TestParse fails and fix it; reply with what you changed"
ccam agent --model opus -- "-v breaks the parser: find out why"
ccam agent --json < task.md            # the task from standard input
ccam agent --resume 7f3c9a2e-… "now add a test for it"
ccam agent --wait=1h "…"               # nobody has room: wait for a reset
ccam --dry-run agent "…"               # the account and command it would run
```

```console
$ ccam agent "summarize what internal/usage does"
ccam: agent runs as work (12% of its 5 hours used); session 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f
internal/usage looks up the 5-hour and weekly usage of every account …
ccam: agent finished as work; follow up with: ccam agent --resume 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f "…"
```

يذهب آخر رد من Claude إلى المخرج القياسي، أما الأسطر التي تبدأ بـ `ccam:` فتذهب إلى الخطأ القياسي. ويطبع `--json` بدلًا من ذلك `{account, session_id, result, is_error, accounts}`، حيث يسرد `accounts` كل حساب عملت به المهمة.

رموز الخروج: رمز claude نفسه (`0` عندما ينتهي الوكيل)، و`75` (`EX_TEMPFAIL`) عندما لا يكون لأي حساب غير الحساب المستدعي هامش، و`1` للأخطاء الأخرى مثل وكيل يحاول تشغيل وكيل آخر.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--model <model>` | نموذج الوكيل؛ وتُحتسب النافذة الأسبوعية لذلك النموذج في الاختيار أيضًا | النموذج الافتراضي للحساب |
| `--resume <id>` | يتابع محادثة وكيل سابق بمهمة متابعة (معرّف الجلسة موجود في سطره الأول والأخير) | معطّل |
| `--json` | يطبع الحساب والجلسة والرد والحسابات التي عملت بها المهمة | معطّل |
| `--max-usage <1-100>` | يفضّل الحسابات التي يقل استخدامها عن هذه القيمة | `auto_max_usage`، وإلا 100 |
| `--strategy headroom\|expiring\|balance` | [قاعدة الاختيار](/features/next#auto) | `auto_strategy`، وإلا `headroom` |
| `--cross-org` | يعمل أيضًا بحسابات المؤسسات الأخرى | `cross_org`، وإلا معطّل |
| `--wait[=<duration>]` | عندما لا يكون لأي حساب هامش، ينتظر حتى يُعاد ضبط أولها، بحد أقصى هذه المدة | معطّل؛ و`5h` عند تحديده وحده |
| `--dry-run` | يطبع الأمر الخاص بالحساب الذي كان سيختاره، دون تشغيله | معطّل |

## كيف يعمل

1. **يُستبعد الحساب المستدعي.** الحساب الذي يُستدعى منه `ccam agent` (حساب الجلسة التي شغّلها ccam، وإلا الحساب الذي يشغّله هذا المجلد) لا يشغّل الوكيل أبدًا، ولا يشغّله كذلك مجلد آخر مسجَّل دخوله بحساب Claude نفسه، لأنه يتشارك حدوده. ويبقى الوكيل في [مؤسسة](/features/next#auto) الحساب المستدعي ما لم يقرر `--cross-org` غير ذلك.
2. **يختار كما يختار `ccam pick`.** يُستعلَم عن الاستخدام كما يفعل [`ccam pick`](/features/pick)، والوكلاء الذين يبدؤون في غضون دقيقتين من بعضهم يحصلون على حسابات مختلفة ما دامت هناك حسابات كافية لديها هامش.
3. **يعمل بأذونات الجلسة، ولا يتجاوزها أبدًا.** لا يوجد أحد ليسمح بأداة في `claude -p`، لذلك يحصل الوكيل على وضع الأذونات ومستوى الجهد الخاصين بالجلسة التي استدعته. وعند تشغيله من طرفية، يُطبَّق الوضع الافتراضي لـ Claude Code وقواعد الأذونات في إعداداتك. وليس لـ `ccam agent` أي خيار يمنح أكثر من ذلك.
4. **عند بلوغ حد الاستخدام، تستمر المحادثة.** عندما يبلغ حساب الوكيل حده، يسجّل ccam الحد كما كان سيفعل [خطاف المكوّن الإضافي](/features/usage#the-limit-record) ويتابع المحادثة نفسها بالحساب الذي لديه أكبر هامش: ينسخها Claude Code إلى هناك (`--resume <file> --fork-session`)، مع الرسالة الأولى "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." ويُجرَّب كل حساب مرة واحدة.
5. **تُسجَّل كل محادثة وكيل** في `~/.local/state/ccam/agent/<session id>.json` (30 يومًا). ويتابعها `--resume <id>` في مجلدها، بالحساب الذي شغّلها إن كان لا يزال لديه هامش، وإلا بحساب آخر. ويستبعد [`ccam next --continue`](/features/next) محادثات الوكلاء، فيواصل متابعة محادثاتك أنت.

## الأمان والحالات الخاصة

- لا يستطيع الوكيل تشغيل وكلاء خاصين به: يميّزه ccam بـ `CCAM_AGENT=1`، ويرفض `ccam agent` العمل هناك.
- إيقاف `ccam agent` (بـ Ctrl+C أو بإشارة أو بمقاطعة Claude Code للأمر) يوقف الوكيل أيضًا.
- لا يرى الوكيل المحادثة التي استدعته، بل مهمته فقط.
- الوكلاء الذين يعدّلون الملفات نفسها في الوقت نفسه يكتب كل منهم فوق عمل الآخرين. أعطِ كل وكيل جزءًا خاصًا به من المهمة، أو git worktree خاصًا به.
- عندما تنتقل الجلسة المستدعية إلى حساب آخر ([`auto_switch`](/features/next#auto-switch))، يتوقف الوكلاء الذين شغّلتهم معها. وتبقى محادثاتهم: يتابعها `--resume <id>`.
- `agent` اسم أمر، لذلك لا يمكن تسمية أي حساب جديد `agent`؛ أما الحساب الموجود بهذا الاسم فلا يزال يبدأ بـ `ccam run agent`.

## صفحات ذات صلة

- [اختيار حساب للسكربتات](/features/pick)
- [التبديل إلى الحساب التالي](/features/next)
- [حدود الاستخدام](/features/usage)
- [الجلسات الخلفية لكل الحسابات](/features/agents)
