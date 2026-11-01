---
title: تشغيل Claude Code
description: شغّل Claude Code بأي حساب باستخدام ccam <name|number> أو ccam run أو لوحة التحكم؛ تُمرَّر الوسائط إلى claude دون تغيير، ويعرض --dry-run الأمر.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: "11df08608a1e"
---

يشغّل `ccam work` برنامج Claude Code بالحساب `work`، وكل ما يأتي بعد الاسم يُمرَّر إلى `claude` دون تغيير. استخدمه في أي مكان كنت ستكتب فيه `claude`: في الطرفية، أو في سكربت، أو مع إدخال يصل عبر أنبوب.

## الاستخدام

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

حدّد اسم الحساب أو رقمه من `ccam list`. لا يقرأ ccam الوسائط التي تأتي بعد الاسم، لذلك يعرض `ccam work --help` مساعدة Claude Code. ووضع `--` بعد الاسم اختياري. ويتصرف الإدخال القياسي ورمز الخروج والإشارات كما مع `claude`.

عندما يبدأ الوسيط الأول بـ `-` ولا يكون من خيارات ccam نفسه، يختار ccam الحساب نيابةً عنك ويمرّر كل الوسائط إلى Claude Code. فيستخدم الحساب المعيّن للمجلد الحالي باستخدام [`ccam folder`](/features/folder)، وإلا فالحساب الافتراضي: الحساب الذي عيّنته باستخدام `ccam use`، أو صاحب أدنى رقم إن لم تعيّن واحدًا.

### لوحة التحكم

يفتح `ccam` دون وسائط [لوحة التحكم](/features/dashboard) عندما يكون الإدخال والإخراج القياسيان طرفية، ويطبع `ccam list` في غير ذلك. انتقل إلى حساب (بالأسهم، أو بكتابة رقمه) واضغط `Enter` لتشغيله؛ ويشغّل `n` الحساب الذي لديه أكبر هامش، ويعرض `o` الخيارات `--continue` أو `--resume` أو وسائط أخرى. وتُغلق لوحة التحكم قبل أن يبدأ Claude Code، فتُستعاد الطرفية أولًا.

**مقاعد Team وEnterprise.** يبدّل الاختيار التلقائي للحساب (`ccam next` و`--fallback` و`ccam pick` والمفتاح `n` في لوحة التحكم، و[`auto_switch`](/features/next#auto-switch) بعد أن تشغّله) أيضًا بين مقاعد المؤسسة نفسها في Claude Team أو Enterprise، وهذا قد يخالف شروط Anthropic. اقرأ [الأسئلة الشائعة](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) وتحقّق من ذلك مع مؤسستك ومع Anthropic قبل أن تعتمد عليه.

عندما يكون `CLAUDE_CODE_OAUTH_TOKEN` أو `ANTHROPIC_API_KEY` أو `ANTHROPIC_AUTH_TOKEN` معيّنًا في الصدفة، قد يستخدم Claude Code ذلك المفتاح أو الرمز بدلًا من تسجيل دخول الحساب نفسه، لذلك يحذّرك ccam في الخطأ القياسي قبل أن يبدأ (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). وهو يقرأ المتغيرات فقط؛ أما [حساب مفتاح API](/features/api) فيمسحها بدلًا من ذلك.

افتراضيًا، يستبدل ccam نفسه بـ Claude Code، فلا يبقى شيء آخر قيد التشغيل. ومع `relay = true` (أو [`auto_switch = true`](/features/next#auto-switch) على Linux وmacOS) في أعلى `config.toml`، فوق أول `[[group]]` أو `[[profile]]` أو `[[folder]]`، يشغّل ccam بدلًا من ذلك الجلسات التفاعلية بوصفها عمليات فرعية ويبقى قيد التشغيل بجانبها. وهذا يتيح لـ `/ccam:switch`، أو لـ `auto_switch` عند حد الاستخدام، متابعة المحادثة على حساب آخر بنفسه، مع خيارات الجلسة ووضع الأذونات ومستوى الجهد والوكيل؛ وقبل أن ينقلها `/ccam:switch` إلى حساب من مؤسسة أخرى، يسألك ccam ([التفاصيل](/features/next#relay)).

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--dry-run` | يطبع البيئة والأمر اللذين كان ccam سيشغّلهما، ثم يخرج | معطّل |
| `--fallback` | لتشغيلات `-p`/`--print`: عندما يبلغ الحساب حد استخدامه، يشغّل الأمر نفسه مجددًا بالحساب المسجّل الدخول الذي لديه أكبر هامش ([التفاصيل](/features/next#run-again-at-the-limit)) | معطّل |
| `-v`، `--version` | يطبع `ccam <version>` (فقط بوصفه الوسيط الأول) | |
| `-h`، `--help` | مساعدة ccam، بالشكل `ccam --help` أو `ccam run --help` | |

قواعد الوسائط:

- لا يقرأ ccam خياراته الخاصة إلا قبل اسم الحساب. وبعد الاسم، كل وسيط يخص Claude Code، بما في ذلك `--dry-run` و`--help`.
- يعمل `--dry-run` بالأشكال `ccam --dry-run work` و`ccam run --dry-run work` و`ccam --dry-run` (لحساب المجلد أو الحساب الافتراضي). ويمكن أن يأتي `--dry-run` و`--fallback` بأي ترتيب.
- `ccam run <name|number> [claude args]` هو الشكل الطويل لـ `ccam <name>`. ويستخدم `ccam run` دون اسم حساب المجلد أو الحساب الافتراضي.
- لا يجوز أن تكون أسماء الحسابات أسماء أوامر ccam مثل `add` أو `list` أو `run` أو `usage`، فلا يكون `ccam <word>` ملتبسًا أبدًا.

يطبع التشغيل التجريبي سطرًا واحدًا يمكنك لصقه في الصدفة، أو في PowerShell على Windows. ويبدأ بما يغيّره التشغيل في البيئة: `CLAUDE_CONFIG_DIR` و`CCAM_PROFILE` دائمًا، و`env -u <name>` للمتغير الذي يُسقطه (مثل `CCAM_RELAY` موروث)، ولـ [حساب مفتاح API](/features/api) نقطة النهاية والنموذج ومصدر المفتاح (`ANTHROPIC_API_KEY="$WORK_KEY"`)، لا المفتاح نفسه أبدًا:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

يعيد سطر PowerShell المتغيرات إلى ما كانت عليه بعد ذلك، كما يفعل `VAR=x cmd` في الصدفة، ويسلّم Claude Code سطر أوامره عبر `Start-Process`، فيشغّله Windows PowerShell 5.1 وPowerShell 7 بالطريقة نفسها. وعندما يعيّن التشغيل متغيرًا إلى قيمة فارغة، يذكر سطر `#` تحته أن PowerShell قبل 7.5 يزيل مثل هذا المتغير بدلًا من ذلك. ولـ `claude.cmd` المثبّت عبر npm، يشغّل السطر أمر `cmd.exe` الذي سيستخدمه التشغيل (انظر أدناه)، والوسيط الذي سيرفضه التشغيل هناك يُفشل التشغيل التجريبي أيضًا.

## كيف يعمل

يمرّ كل تشغيل بالخطوات نفسها، من سطر الأوامر إلى `exec`:

<!-- motion:launch -->

يضيء كل مربع عندما تعمل خطوته، والأسطر في الأسفل هي البيئة والوسائط التي يبدأ بها Claude Code.

1. يبدأ ccam في `app.Main`، الذي يستدعيه `cmd/ccam/main.go`. وإذا شُغّل ccam باسم آخر، فهو [أمر مختصر](/features/shortcuts) ويذهب مباشرة إلى المشغّل.
2. ينظر ccam في `os.Args` قبل أن يتدخل cobra (`app.parse`). لا يذهب إلى `cli.Execute` إلا `ccam` وحده (لوحة التحكم)، واسم أمر من أوامر ccam (`config.Reserved`)، و`-h`/`--help` و`-v`/`--version`. وكل ما عدا ذلك تشغيل.
3. يقرأ ccam الملف `~/.config/ccam/config.toml`، وهو الملف الوحيد الذي يحلّله (`launcher.Launch`).
4. يجد ccam الحساب بالاسم أو الرقم، ويحذف `--` واحدًا في البداية، ويحدد الملف التنفيذي `claude` (`launcher.Prepare` و`ResolveClaude`). ودون اسم، يأخذ حساب المجلد، بمقارنة مجلد العمل مع إدخالات `[[folder]]` في `config.toml` كسلاسل نصية، دون تمييز حالة الأحرف على macOS وWindows (`HereProfile`). تُخزَّن نتيجة البحث في PATH مؤقتًا في `~/.local/state/ccam/claude-path` مع وقت تعديل الملف وتجزئة لـ `PATH`، فلا يبحث ccam في `PATH` عند كل تشغيل. والملف المتغيّر أو المفقود، أو `PATH` المتغيّر، يؤدي إلى بحث جديد. ويجب أن يكون مسار `claude_bin` في `config.toml` (يُوسَّع `~/`) ملفًا تنفيذيًا؛ وإلا يتوقف التشغيل مع `claude_bin X is not an executable file. Check claude_bin in config.toml` قبل كتابة أي شيء. ويرفض ccam أي `claude` يكون هو ccam نفسه، بالمقارنة حسب inode، لأن ذلك سيؤدي إلى حلقة.
5. يقارن ccam مجلد الحساب `~/.local/share/ccam/profiles/<name>/` بإعدادات المشاركة الخاصة به، مستدعيًا `lstat` و`readlink` فقط، أي استدعاء نظام واحد أو اثنين لكل عنصر (`checkLinks` و`share.Check`). ويُنشأ الرابط المفقود في الحال. وأي شيء قد يؤدي إلى فقدان بيانات يوقف التشغيل.
6. يضيف ccam خوادم MCP المشتركة (`applyMCP`؛ انظر أدناه).
7. يسجّل ccam الجلسة (`session.Register`). فيكتب علامة باسم معرّف العملية (PID) الخاص به، تحتوي على وقت بدء العملية، في `~/.local/state/ccam/locks/profile-<name>/<pid>`، أيًّا كانت قيمة `XDG_RUNTIME_DIR`، فترى كل صدفة العلامات نفسها. وللحساب الذي يشارك سجل المحادثات، يبحث ccam أولًا عن جلسة تعمل في `history-<group>-<hash>`، حيث التجزئة هي تجزئة المجلد الحالي بعد حل الروابط الرمزية، ويسجّل هناك أيضًا (`historyCheck`).

   يحتفظ `exec` بمعرّف العملية، فتشير العلامة بعد ذلك إلى عملية Claude Code نفسها. ويستخدم ccam ملف علامة لا قفلًا محتجزًا، لأن العمليات التي يشغّلها Claude Code سترث القفل وتحتفظ به بعد انتهاء الجلسة. والعلامة التي يتعذر كتابتها لا توقف التشغيل: يقول ccam `ccam: could not mark the session as running (…); starting Claude Code anyway`.

   وفحوص ccam الأخرى لمعرفة هل يعمل حساب (إعادة التسمية، والإزالة، وتغييرات المجموعات، و`doctor --fix`، ومزامنة MCP) تقرأ أيضًا سجلات Claude Code نفسه في مجلد الحساب (`sessions/<pid>.json` و`daemon.lock`)، فترى كذلك Claude Code الذي بدأ دون ccam: `claude` العادي، أو بيئة تطوير، أو الجلسات الخلفية وخدمتها الخلفية.
8. يعيّن ccam عنوان الطرفية إلى `claude · <name>` (عندما يكون الخطأ القياسي طرفية)، ويكتب وقت التشغيل في `~/.local/state/ccam/last-used/<name>` دون `fsync`، ويستدعي `syscall.Exec` (`Plan.Exec`).

يعمل Claude Code مع تعيين `CLAUDE_CONFIG_DIR` إلى مجلد الحساب و`CCAM_PROFILE` إلى اسم الحساب. ويمسح ccam أي قيم موروثة لهذين المتغيرين، و`CCAM_RELAY` الموروث. وبالنسبة إلى مجلد `~/.claude` المسجل باستخدام `ccam import` على macOS، يعيّن أيضًا `CLAUDE_SECURESTORAGE_CONFIG_DIR` إلى قيمة فارغة، كي يستخدم الحساب و`claude` العادي تسجيل الدخول نفسه في سلسلة المفاتيح؛ وعلى Windows يفعل الشيء نفسه عندما يحفظ Claude Code تسجيل الدخول في Credential Manager.

### لماذا exec

يستبدل ccam عمليته بـ `claude` (`syscall.Exec`). ونتيجة لذلك:

- لا يبقى ccam في الذاكرة؛
- تكون ضغطات المفاتيح والإشارات ورمز الخروج خاصة بـ Claude Code نفسه؛
- تبقى الوسائط مصفوفة، لا نصًا للصدفة أبدًا.

لا يوجد `exec` في Windows. هناك يشغّل ccam الملف `claude.exe` عمليةً فرعية بالبيئة والوسائط ووحدة التحكم نفسها، وينتظرها ثم يخرج برمز خروجها (`Plan.Exec`). والملف `claude.cmd` الذي ثبّته npm يعمل في `cmd.exe` بمساره الكامل مع وضع كل وسيط بين علامات اقتباس، فلا يقرأ `cmd.exe` أيًّا منها على أنه من صياغته (`&` و`%PATH%`)؛ والوسيط الذي فيه فاصل أسطر لا يمكن تمريره بهذه الطريقة ويُرفض، مع نصيحة بتثبيت Claude Code الأصلي.

يصل Ctrl+C إلى Claude Code من وحدة التحكم مباشرة، لذا يتجاهله ccam. وتحمل علامة الجلسة PID الخاص بـ ccam، الذي يعيش تمامًا ما دام Claude Code يعيش، وتُزال عند خروج Claude Code. وأثناء الانتظار لا يفعل ccam أي شيء آخر.

والعمل الذي يأتي بعد الجلسة، مثل ملاحظة الخوادم المضافة باستخدام `claude mcp add`، يحدث عند التشغيل التالي: يقارن ccam وقت تعديل `.claude.json` مع `~/.local/state/ccam/mcp/<name>/seen` ولا يحلّل الملف إلا إذا تغيّر.

### كيف يُوجَّه سطر الأوامر

ثلاثة أسطر أوامر، وثلاثة مسارات عبر `internal/app`:

<!-- motion:launch-route -->

يُطابَق الأمر المختصر على `argv[0]`، و`ccam <name>` على الوسيط الأول، وكلاهما ينتهي في `syscall.Exec`. ولا يصل إلى cobra إلا أمر من أوامر ccam مثل `ccam usage`.

### خوادم MCP المشتركة

مع الاستراتيجية الافتراضية `flag`، يمرّر ccam خوادم MCP الخاصة بالمجموعة إلى Claude Code بوصفها وسيطًا، يوضع قبل وسائطك:

<!-- motion:launch-mcp -->

لا يضيف ccam الوسيط إلا إذا كان ملف `mcp.json` الخاص بالمجموعة موجودًا، ويضعه مباشرة بعد `claude` (`mcp.FlagArgs` و`Plan.InsertArgs`). ثم يظل `--mcp-config` أو `--strict-mcp-config` الخاص بك ساريًا فوق ذلك. وهو وسيط واحد، `--mcp-config=<path>`، لأن `--mcp-config` في Claude Code يأخذ عدة قيم: فلو كُتب وسيطين، لقرأ أيضًا الطلب الذي يليه على أنه ملف إعداد.

يقرر ccam أين لا يمكن وضع الخيار (`claudeVerb`). فهو يقرأ الأمر الفرعي لـ Claude Code من الوسائط كما يفعل Claude Code، متخطيًا الخيارات مع قيمها. لا يحصل `remote-control` أو `rc` على الخيار (أدناه)، ولا فعل الجلسات الخلفية الذي يأتي أولًا (`attach` و`logs` و`stop` و`kill` و`respawn` و`rm`). ولا يحصل عليه `ccam work mcp list` و`get` و`remove` أيضًا، لأن Claude Code يرفضه هناك؛ وبدلًا من ذلك يقول ccam في الخطأ القياسي إن خوادم المجموعة تُحمَّل في جلسات الحساب أيضًا:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

وكل أمر فرعي آخر وكل جلسة يحتفظان بالخيار.

مع `mcp_strategy = "sync"` في `config.toml`، يكتب ccam الخوادم في ملف `.claude.json` الخاص بالحساب بدلًا من ذلك، ويتخطى ذلك إذا كان الحساب يعمل بالفعل. ويتخطى التشغيل هذا العمل أيضًا عندما لا يكون ملف `mcp.json` الخاص بالمجموعة ولا `.claude.json` قد تغيّر منذ آخر فحص. وإذا كان لدى الحساب خادمه الخاص بالاسم نفسه الذي يحمله خادم في المجموعة، تحتفظ المزامنة بخادم الحساب وتقول ذلك مرة واحدة. وتعيد المزامنة أيضًا النسخة التي اختفت من `.claude.json` (`put back group default's MCP servers …`)، ولا تغيّر أو تزيل إلا النسخ التي كتبتها ولم يغيّرها أحد منذ ذلك الحين.

### Remote Control

يشغّل `ccam work rc` (أو `ccam work remote-control`) ميزة Remote Control في Claude Code بالحساب `work`. يرفض `claude remote-control` الخيار `--mcp-config` قبل فعله ولا يقبل أيًّا منه بعده، ولا تقرأ جلساته خوادم MCP إلا من ملف `.claude.json` الخاص بالحساب. لذلك تعمل استراتيجية الخيار هنا بطريقة مختلفة:

<!-- motion:launch-rc -->

1. يجد ccam الفعل `remote-control` أو `rc` بوصفه فعلًا لـ Claude Code، حتى بعد الخيارات وقيمها (`ccam work --model opus rc`)، ولا يضيف أي خيار (`claudeVerb` و`mcp.FlagArgs`).
2. يقارن ccam وقت تعديل ملف `mcp.json` الخاص بالمجموعة مع `sync-stamp` الخاص بالحساب، ويتحقق من أن كل خادم من المجموعة لا يزال في `.claude.json` (`mcp.Synced`). وإذا كان `.claude.json` يحتوي بالفعل على خوادم المجموعة الحالية، أو لم تكن لدى المجموعة أي خوادم، تبدأ Remote Control دون سؤال.
3. وإلا يسمّي ccam الخوادم التي ستفتقدها الجلسات ويسأل:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. عند `y`، ينسخها ccam كما تفعل استراتيجية `sync` (`mcp.Sync`): تحت قفل Claude Code، وبعد نسخة احتياطية، مع الإبقاء على خوادم الحساب نفسه. وعند `N`، أو بدون طرفية، تبدأ Remote Control بدونها، ويشير ccam إلى `mcp_strategy = "sync"`، التي تنسخها قبل كل تشغيل.

يسأل ccam عند كل تشغيل لـ Remote Control إلى أن تُنسخ الخوادم. وبعد النسخ، لا يسأل مجددًا إلا عندما يتغير ملف `mcp.json` الخاص بالمجموعة أو يختفي خادم منسوخ من `.claude.json`. وتتقدم النسخة الموجودة في `.claude.json` على `--mcp-config` في عمليات التشغيل العادية اللاحقة، لذلك تحدّث الأوامر `ccam mcp add` و`remove` و`import` النسخ التي تحتفظ بها الحسابات الأخرى ([خوادم MCP المشتركة](/features/mcp)).

### ميزانية 10 ms

ميزانية مسار التشغيل هي 10 ms عند المئين 95؛ وعمليًا يضيف ccam نحو 4 ms قبل أن يبدأ Claude Code. لا يبني مسار التشغيل أبدًا شجرة أوامر cobra، ولا يشغّل Bubble Tea، ولا يفتح اتصال شبكة، ويفشل `TestPackageInitIsCheap` إذا استغرقت دوال `init()` في الحزم أكثر من 5 ms إجمالًا. يقيس `scripts/perf-check.sh` الأوامر `ccam --dry-run work` و`ccam work` وأمرًا مختصرًا باستخدام `hyperfine` ويفشل فوق 10 ms، ويشغّل `BenchmarkLaunchPath` في `internal/launcher` الخطوات في حلقة. ولرؤية الوقت لكل خطوة، عيّن `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## الأمان والحالات الخاصة

- الاسم أو الرقم غير المعروف يتوقف مع `no account "<name>"` ويوجّهك إلى `ccam list`. وإن لم تكن هناك حسابات على الإطلاق، يطلب منك ccam إضافة واحد باستخدام `ccam add <name>`.
- إذا أصبح مسار مشترك ملفًا عاديًا أو صار يشير إلى المكان الخطأ، يتوقف ccam مع `stopped because of share link problems:`، ويسرد المشكلات، ويطلب منك تشغيل `ccam doctor --fix`. ولا ينقل بيانات في مسار التشغيل.
- إذا كانت جلسة Claude Code أخرى تشارك سجل المحادثات مع هذا الحساب (المجموعة نفسها) تعمل في المجلد نفسه، يحذّرك ccam من أن السجل قد يختلط ويسأل `Run anyway? [y/N]`. ويُطابَق المجلد بعد حل الروابط الرمزية. ودون طرفية (والإدخال القياسي من `/dev/null` يُعد بلا طرفية)، أو من مهمة في الخلفية، يعرض التحذير ويتابع. وتشغيل `-p` جديد يبدأ محادثة خاصة به فلا يحصل على تحذير؛ أما `--version` و`--help` والأوامر الفرعية التي لا تفتح محادثة، مثل `mcp` أو `auth` أو `plugin`، فلا تُفحص إطلاقًا.
- يتخطى `--dry-run` مزامنة MCP، وإشعار الخادم الجديد، وملاحظة `mcp list`، وفحص السجل، وعلامة الجلسة، والعنوان، ووقت آخر استخدام، ولا يستدعي exec. لكنه يظل يفحص روابط المشاركة، وينشئ الرابط المفقود كما يفعل التشغيل العادي.
- عند تشغيل `ccam work rc` بينما `work` يعمل بالفعل، يترك ccam ملف `.claude.json` دون تغيير، ويخبرك بذلك، ويبدأ Remote Control بدون خوادم المجموعة؛ شغّله مرة أخرى بعد انتهاء الجلسة لنسخها. ومع `--dry-run` يطبع الرسالة لكنه لا يسأل ولا ينسخ.
- مع `mcp_strategy = "sync"`، تطبع المزامنة الفاشلة `ccam: MCP sync failed (launching anyway): …` ويبدأ Claude Code بالخوادم التي كانت لديه بالفعل.

## صفحات ذات صلة

- [الأوامر المختصرة](/features/shortcuts)
- [استخدام Remote Control مع ccam](/guides/remote-control)
- [إدارة الحسابات](/features/accounts)
- [مرجع الأوامر](/docs/commands)
- [مشاركة خوادم MCP والإعدادات والسجل](/docs/sharing)
- [شرح CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
