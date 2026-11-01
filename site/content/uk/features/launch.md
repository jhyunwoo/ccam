---
title: Запуск Claude Code
description: Запускайте Claude Code з будь-якого облікового запису через ccam <name|number>, ccam run або панель ccam; аргументи йдуть у claude як є, --dry-run дає команду.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: "11df08608a1e"
---

`ccam work` запускає Claude Code від імені облікового запису `work`, і все, що йде після назви, передається в `claude` без змін. Використовуйте його всюди, де ви ввели б `claude`: у терміналі, у скрипті або з введенням через конвеєр.

## Використання

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

Вкажіть назву облікового запису або його номер із `ccam list`. ccam не читає аргументи після назви, тож `ccam work --help` показує довідку Claude Code. `--` після назви необов'язковий. Стандартне введення, код виходу та сигнали поводяться так само, як із `claude`.

Коли перший аргумент починається з `-` і не є одним із власних прапорців ccam, ccam сам вибирає обліковий запис і передає всі аргументи в Claude Code. Він використовує обліковий запис, заданий для поточної теки через [`ccam folder`](/features/folder), а якщо його немає — обліковий запис за замовчуванням: той, який задано через `ccam use`, або обліковий запис із найменшим номером, якщо ви його не задали.

### Панель керування

`ccam` без аргументів відкриває [панель керування](/features/dashboard), коли стандартне введення та виведення є терміналом, а в інших випадках виводить `ccam list`. Перейдіть до облікового запису (стрілками або ввівши його номер) і натисніть `Enter`, щоб запустити його; `n` запускає обліковий запис, де найбільше запасу, а `o` пропонує `--continue`, `--resume` або інші аргументи. Панель керування закривається до запуску Claude Code, тож спершу термінал відновлюється.

**Місця в планах Team і Enterprise.** Автоматичний вибір облікового запису (`ccam next`, `--fallback`, `ccam pick`, `n` на панелі керування, а також, щойно ви його ввімкнете, [`auto_switch`](/features/next#auto-switch)) також перемикає між місцями однієї організації Claude Team або Enterprise, що може порушувати умови Anthropic. Прочитайте [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) і уточніть це у своїй організації та в Anthropic, перш ніж на нього покладатися.

Коли в оболонці задано `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` або `ANTHROPIC_AUTH_TOKEN`, Claude Code може використати цей ключ або токен замість власного входу облікового запису, тож ccam попереджає вас про це в стандартному потоці помилок перед запуском (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Він лише читає ці змінні; [обліковий запис із ключем API](/features/api) натомість їх очищає.

Типово ccam замінює себе на Claude Code, тож більше нічого не працює. З `relay = true` (або [`auto_switch = true`](/features/next#auto-switch) на Linux і macOS) угорі `config.toml`, над першим `[[group]]`, `[[profile]]` або `[[folder]]`, ccam натомість запускає інтерактивні сеанси як дочірні процеси й працює далі поруч із ними. Завдяки цьому `/ccam:switch` або `auto_switch` на ліміті використання може сам продовжити розмову в іншому обліковому записі з параметрами сеансу, режимом дозволів, рівнем зусиль і агентом; перш ніж `/ccam:switch` перенесе її в обліковий запис іншої організації, ccam запитує вас ([докладніше](/features/next#relay)).

## Параметри

| Параметр | Дія | За замовчуванням |
| --- | --- | --- |
| `--dry-run` | Вивести середовище й команду, які виконав би ccam, і завершитися | вимкнено |
| `--fallback` | Для запусків із `-p`/`--print`: коли обліковий запис досягає ліміту використання, знову виконати ту саму команду від імені облікового запису з входом, де найбільше запасу ([докладніше](/features/next#run-again-at-the-limit)) | вимкнено |
| `-v`, `--version` | Вивести `ccam <version>` (лише як перший аргумент) | |
| `-h`, `--help` | Довідка ccam, як `ccam --help` або `ccam run --help` | |

Правила для аргументів:

- ccam читає власні прапорці лише перед назвою облікового запису. Після назви кожен аргумент належить Claude Code, зокрема `--dry-run` і `--help`.
- `--dry-run` працює як `ccam --dry-run work`, `ccam run --dry-run work` і `ccam --dry-run` (для облікового запису теки або облікового запису за замовчуванням). `--dry-run` і `--fallback` можна вказувати в будь-якому порядку.
- `ccam run <name|number> [claude args]` — повна форма `ccam <name>`. `ccam run` без назви використовує обліковий запис теки або обліковий запис за замовчуванням.
- Назви облікових записів не можуть збігатися з назвами команд ccam, як-от `add`, `list`, `run` чи `usage`, тож `ccam <word>` ніколи не буває неоднозначним.

Пробний запуск виводить один рядок, який можна вставити в оболонку, а у Windows — у PowerShell. Він починається з того, що запуск змінює в середовищі: завжди `CLAUDE_CONFIG_DIR` і `CCAM_PROFILE`, `env -u <name>` для змінної, яку він відкидає (як-от успадкований `CCAM_RELAY`), а для [облікового запису з ключем API](/features/api) — його кінцеву точку, його модель і звідки береться ключ (`ANTHROPIC_API_KEY="$WORK_KEY"`), але ніколи не сам ключ:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

Рядок PowerShell потім повертає змінні назад, як це робить `VAR=x cmd` в оболонці, і передає Claude Code його командний рядок через `Start-Process`, тож Windows PowerShell 5.1 і PowerShell 7 виконують його однаково. Коли запуск задає змінній порожнє значення, рядок `#` під ним повідомляє, що PowerShell до версії 7.5 натомість видаляє таку змінну. Для `claude.cmd`, установленого через npm, рядок виконує команду `cmd.exe`, яку використав би запуск (див. нижче), а аргумент, який запуск там відхилив би, призводить до невдачі й пробного запуску.

## Як це працює

Кожен запуск проходить ті самі кроки, від командного рядка до `exec`:

<!-- motion:launch -->

Кожен блок підсвічується, коли виконується його крок, а рядки під ним — це середовище й аргументи, з якими запускається Claude Code.

1. ccam запускається в `app.Main`, яку викликає `cmd/ccam/main.go`. Якщо ccam запущено під іншою назвою, це [коротка команда](/features/shortcuts), і виклик іде безпосередньо до засобу запуску.
2. ccam переглядає `os.Args` ще до cobra (`app.parse`). Лише `ccam` без нічого (панель керування), назва команди ccam (`config.Reserved`), `-h`/`--help` і `-v`/`--version` потрапляють до `cli.Execute`. Усе інше є запуском.
3. ccam читає `~/.config/ccam/config.toml`, єдиний файл, який він розбирає (`launcher.Launch`).
4. ccam знаходить обліковий запис за назвою або номером, відкидає один початковий `--` і знаходить виконуваний файл `claude` (`launcher.Prepare`, `ResolveClaude`). Без назви він бере обліковий запис теки, порівнюючи робочий каталог із записами `[[folder]]` у `config.toml` як рядки, без урахування регістру на macOS і Windows (`HereProfile`). Результат пошуку в PATH кешується в `~/.local/state/ccam/claude-path` разом із часом зміни файлу та хешем `PATH`, тож ccam не шукає в `PATH` під час кожного запуску. Змінений або відсутній файл чи змінений `PATH` запускає новий пошук. Шлях `claude_bin` у `config.toml` (`~/` розгортається) має бути виконуваним файлом; інакше запуск зупиняється з `claude_bin X is not an executable file. Check claude_bin in config.toml`, перш ніж щось буде записано. ccam відхиляє `claude`, який є самим ccam (порівняння за inode), бо це спричинило б зациклення.
5. ccam порівнює теку облікового запису `~/.local/share/ccam/profiles/<name>/` з її налаштуваннями спільного використання, викликаючи лише `lstat` і `readlink`, один або два системні виклики на елемент (`checkLinks`, `share.Check`). Відсутнє посилання створюється одразу. Усе, що може призвести до втрати даних, зупиняє запуск.
6. ccam додає спільні сервери MCP (`applyMCP`; див. нижче).
7. ccam реєструє сеанс (`session.Register`). Він записує позначку з назвою за своїм PID, що містить час запуску процесу, у `~/.local/state/ccam/locks/profile-<name>/<pid>`, хоч би що казав `XDG_RUNTIME_DIR`, тож кожна оболонка бачить ті самі позначки. Для облікового запису, який спільно використовує історію розмов, ccam спершу шукає запущений сеанс у `history-<group>-<hash>`, де хеш обчислено від поточної теки з розв'язаними символьними посиланнями, і реєструється й там (`historyCheck`).

   `exec` зберігає PID, тож позначка потім вказує на сам процес Claude Code. ccam використовує файл-позначку, а не утримуване блокування, бо процеси, які запускає Claude Code, успадкували б блокування й утримували б його після завершення сеансу. Позначка, яку не вдається записати, не зупиняє запуск: ccam повідомляє `ccam: could not mark the session as running (…); starting Claude Code anyway`.

   Інші перевірки ccam того, чи працює обліковий запис (перейменування, видалення, зміни груп, `doctor --fix`, синхронізація MCP), також читають власні записи Claude Code у теці облікового запису (`sessions/<pid>.json`, `daemon.lock`), тож вони бачать і Claude Code, запущений без ccam: звичайний `claude`, IDE, фонові сеанси та їхній демон.
8. ccam встановлює заголовок термінала `claude · <name>` (коли стандартний потік помилок є терміналом), записує час запуску в `~/.local/state/ccam/last-used/<name>` без `fsync` і викликає `syscall.Exec` (`Plan.Exec`).

Claude Code працює зі змінною `CLAUDE_CONFIG_DIR`, що вказує на теку облікового запису, і `CCAM_PROFILE`, що містить назву облікового запису. ccam очищає будь-які успадковані значення цих змінних, а також успадкований `CCAM_RELAY`. Для теки `~/.claude`, зареєстрованої через `ccam import` на macOS, він також встановлює `CLAUDE_SECURESTORAGE_CONFIG_DIR` у порожнє значення, щоб обліковий запис і звичайний `claude` використовували той самий вхід у в'язці ключів; у Windows він робить те саме, коли Claude Code зберігає вхід у диспетчері облікових даних.

### Чому exec

ccam замінює власний процес на `claude` (`syscall.Exec`). У результаті:

- ccam не залишається в пам'яті;
- натискання клавіш, сигнали й код виходу належать самому Claude Code;
- аргументи залишаються масивом і ніколи не стають рядком оболонки.

У Windows немає `exec`. Там ccam запускає `claude.exe` як дочірній процес із тим самим середовищем, аргументами й консоллю, чекає на нього й завершується з його кодом виходу (`Plan.Exec`). `claude.cmd`, установлений через npm, виконується в `cmd.exe` з повним шляхом, а кожен аргумент береться в лапки так, щоб `cmd.exe` не прочитав жоден із них як свій синтаксис (`&`, `%PATH%`); аргумент із розривом рядка так передати не можна, і його буде відхилено з порадою встановити нативний Claude Code.

Ctrl+C доходить до Claude Code безпосередньо з консолі, тому ccam його ігнорує. Маркер сеансу містить PID ccam, який живе рівно стільки ж, скільки Claude Code, і видаляється, коли Claude Code завершується. Поки ccam чекає, він нічого іншого не робить.

Робота після сеансу, як-от виявлення серверів, доданих через `claude mcp add`, відбувається під час наступного запуску: ccam порівнює час зміни `.claude.json` з `~/.local/state/ccam/mcp/<name>/seen` і розбирає файл лише тоді, коли він змінився.

### Як маршрутизується командний рядок

Три командні рядки, три маршрути через `internal/app`:

<!-- motion:launch-route -->

Коротка команда розпізнається за `argv[0]`, `ccam <name>` — за першим аргументом, і обидва маршрути закінчуються `syscall.Exec`. Лише команда ccam, як-от `ccam usage`, доходить до cobra.

### Спільні сервери MCP

За стратегії `flag`, що діє за замовчуванням, ccam передає сервери MCP групи в Claude Code як аргумент, розміщений перед вашими:

<!-- motion:launch-mcp -->

ccam додає аргумент, лише якщо існує `mcp.json` групи, і ставить його одразу після `claude` (`mcp.FlagArgs`, `Plan.InsertArgs`). Ваші власні `--mcp-config` чи `--strict-mcp-config` після цього все одно застосовуються поверх нього. Це один аргумент, `--mcp-config=<path>`, бо `--mcp-config` у Claude Code приймає кілька значень: якби його записали як два аргументи, наступний запит також було б прочитано як файл конфігурації.

ccam вирішує, куди прапорець потрапити не може (`claudeVerb`). Він читає підкоманду Claude Code з аргументів так само, як Claude Code, пропускаючи параметри разом із їхніми значеннями. `remote-control` або `rc` прапорця не отримує (нижче), як і команда фонового сеансу, вказана першою (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`). `ccam work mcp list`, `get` і `remove` теж його не отримують, бо Claude Code там його відхиляє; натомість ccam повідомляє в стандартний потік помилок, що сервери групи також завантажуються в сеансах облікового запису:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Кожна інша підкоманда й кожен сеанс зберігають прапорець.

З `mcp_strategy = "sync"` у `config.toml` ccam натомість записує сервери в `.claude.json` облікового запису й пропускає цей крок, якщо обліковий запис уже запущено. Запуск пропускає цю роботу й тоді, коли ні `mcp.json` групи, ні `.claude.json` не змінилися після останньої перевірки. Якщо обліковий запис має власний сервер із тією самою назвою, що й сервер групи, синхронізація зберігає сервер облікового запису й один раз повідомляє про це. Синхронізація також повертає копію, яка зникла з `.claude.json` (`put back group default's MCP servers …`), і змінює чи видаляє лише ті записані нею копії, яких відтоді ніхто не змінював.

### Remote Control

`ccam work rc` (або `ccam work remote-control`) запускає Remote Control з Claude Code від імені облікового запису `work`. `claude remote-control` відхиляє `--mcp-config` перед своєю командою й не приймає його після неї, а його сеанси читають сервери MCP лише з `.claude.json` облікового запису. Тож стратегія прапорця тут працює інакше:

<!-- motion:launch-rc -->

1. ccam знаходить `remote-control` або `rc` як команду Claude Code, навіть після параметрів та їхніх значень (`ccam work --model opus rc`), і не додає прапорця (`claudeVerb`, `mcp.FlagArgs`).
2. ccam порівнює час зміни `mcp.json` групи з `sync-stamp` облікового запису й перевіряє, що кожен сервер групи досі є в `.claude.json` (`mcp.Synced`). Якщо `.claude.json` уже містить поточні сервери групи або група їх не має, Remote Control запускається без запитання.
3. Інакше ccam називає сервери, яких бракуватиме сеансам, і запитує:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. Після відповіді `y` ccam копіює їх так само, як стратегія `sync` (`mcp.Sync`): під блокуванням Claude Code, після створення резервної копії й зі збереженням власних серверів облікового запису. Після відповіді `N` або без термінала Remote Control запускається без них, а ccam радить `mcp_strategy = "sync"`, що копіює їх перед кожним запуском.

ccam запитує під час кожного запуску Remote Control, доки сервери не скопійовано. Після копіювання він запитує знову, лише коли `mcp.json` групи зміниться або скопійований сервер зникне з `.claude.json`. Копія в `.claude.json` перемагає `--mcp-config` під час наступних звичайних запусків, тож `ccam mcp add`, `remove` і `import` оновлюють копії, що є в інших облікових записах ([Спільні сервери MCP](/features/mcp)).

### Бюджет у 10 мс

Бюджет для шляху запуску — 10 мс на 95-му процентилі; на практиці ccam додає близько 4 мс до запуску Claude Code. Шлях запуску ніколи не будує дерево команд cobra, не запускає Bubble Tea і не відкриває мережевих з'єднань, а `TestPackageInitIsCheap` не проходить, якщо функції пакетів `init()` разом займають понад 5 мс. `scripts/perf-check.sh` вимірює `ccam --dry-run work`, `ccam work` і коротку команду за допомогою `hyperfine` і завершується помилкою, якщо час перевищує 10 мс, а `BenchmarkLaunchPath` в `internal/launcher` виконує ці кроки в циклі. Щоб побачити час кожного кроку, задайте `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Безпека та граничні випадки

- Невідома назва чи номер зупиняє запуск із повідомленням `no account "<name>"` і відсилає до `ccam list`. Якщо облікових записів немає взагалі, ccam пропонує додати обліковий запис через `ccam add <name>`.
- Якщо спільний шлях став звичайним файлом або вказує не туди, ccam зупиняється з повідомленням `stopped because of share link problems:`, перелічує проблеми й просить виконати `ccam doctor --fix`. На шляху запуску він не переміщує дані.
- Якщо в тій самій теці працює інший сеанс Claude Code, який спільно використовує історію розмов із цим обліковим записом (та сама група), ccam попереджає, що історія може змішатися, і запитує `Run anyway? [y/N]`. Теку зіставляють із розв'язаними символьними посиланнями. Без термінала (стандартне введення з `/dev/null` вважається його відсутністю) або із завдання у фоні він показує попередження й продовжує. Новий запуск `-p` починає власну розмову й попередження не отримує; `--version`, `--help` і підкоманди, які не відкривають розмови, як-от `mcp`, `auth` чи `plugin`, узагалі не перевіряються.
- `--dry-run` пропускає синхронізацію MCP, повідомлення про нові сервери, примітку `mcp list`, перевірку історії, позначку сеансу, заголовок і час останнього використання, а також не виконує exec. Він усе одно перевіряє спільні посилання й створює відсутнє посилання, як звичайний запуск.
- `ccam work rc`, коли `work` уже запущено, не чіпає `.claude.json`, повідомляє про це й запускає Remote Control без серверів групи; щоб скопіювати їх, виконайте команду знову після завершення сеансу. З `--dry-run` вона виводить повідомлення, але нічого не запитує й не копіює.
- З `mcp_strategy = "sync"` невдала синхронізація виводить `ccam: MCP sync failed (launching anyway): …`, і Claude Code запускається з тими серверами, які вже мав.

## Пов'язані сторінки

- [Короткі команди](/features/shortcuts)
- [Remote Control з ccam](/guides/remote-control)
- [Керування обліковими записами](/features/accounts)
- [Довідник команд](/docs/commands)
- [Спільні сервери MCP, налаштування та історія](/docs/sharing)
- [Що таке CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
