---
title: Довідник команд
metaTitle: Довідник команд ccam
description: Усі команди та параметри ccam, від ccam add і ccam <name> до ccam usage, ccam next і ccam doctor.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: "28e1ad1dcbad"
---

Додайте `--dry-run` перед командою, щоб побачити, що вона зробила б, не виконуючи її, наприклад `ccam --dry-run work`. Це працює під час запуску Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) і з `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` та `doctor --fix`. Команди, які лише читають, як-от `list`, `usage`, `share` і `ext`, приймають його й однаково нічого не змінюють. Решта відхиляють його й не виконуються, зокрема `group rename`.

<!-- commands -->

## Запуск Claude Code

`ccam <name|number> [arguments]` запускає Claude Code від імені цього облікового запису. ccam не читає аргументи після назви; він передає їх у `claude` у тому самому порядку, а стандартне введення, код виходу та сигнали поводяться так само, як із самим `claude`. `--` після назви необов'язковий.

`ccam` з аргументами, але без назви облікового запису, наприклад `ccam -p "summarise"`, використовує обліковий запис, заданий для поточної теки, а якщо його немає — обліковий запис за замовчуванням. Задайте їх командами `ccam folder <name> [folder]` і `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` перед назвою облікового запису призначений для запусків із `-p`/`--print`. Коли обліковий запис досягає ліміту використання, ccam знову виконує ту саму команду від імені облікового запису з входом, де найбільше запасу, і виводить лише результат запуску, який завершився. Див. [Облікові записи за теками](/features/folder) і [Повторний запуск при ліміті](/features/next#run-again-at-the-limit).

**Місця в планах Team і Enterprise.** Автоматичний вибір облікового запису (`ccam next`, `--fallback`, `ccam pick`, `n` на панелі керування, а також, щойно ви його ввімкнете, [`auto_switch`](/features/next#auto-switch)) також перемикає між місцями однієї організації Claude Team або Enterprise, що може порушувати умови Anthropic. Прочитайте [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) і уточніть це у своїй організації та в Anthropic, перш ніж на нього покладатися.

Коротка команда, наприклад `ccwork`, запускає свій обліковий запис так само: `ccwork --resume` робить те саме, що й `ccam work --resume`. Див. [Короткі команди](/docs/commands#shortcuts) нижче.

## Додавання та видалення облікових записів

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

`ccam import` без теки шукає `~/.claude` і `~/.claude-*` та запитує про кожну з них; про теку, яку не вдалося зареєструвати, повідомляється як `Skipped <path>: <reason>`, а решта реєструються. Якщо назву теки вже зайнято, імпортуйте цю теку окремо командою `ccam import <folder> --name <name>`. Імпортовані теки реєструються там, де вони є; ccam ніколи не переміщує й не копіює дані входу.

`ccam login --expired` також виконує вхід в облікові записи, вхід яких очистив Claude Code. `ccam remove` виконує вихід із керованого облікового запису (того, теку якого створив ccam, на відміну від зареєстрованого через `ccam import`), перш ніж його тека потрапить до кошика, якщо тільки інший обліковий запис не має той самий вхід. Видалення імпортованого облікового запису виводить його з груп і лишає в його теці копію кожного спільного елемента. На macOS `ccam rename` зберігає назву теки, бо в'язка ключів зберігає вхід за шляхом теки.

`ccam add --api` додає обліковий запис, який використовує ключ API замість входу Claude; ccam читає ключ під час кожного запуску й ніколи його не зберігає. Див. [Облікові записи з ключем API](/features/api).

`ccam rename`, `ccam remove` і команди `ccam group` відмовляються змінювати обліковий запис, поки Claude Code працює від його імені, хоч би як його було запущено (зокрема у фонових сеансах). Спершу завершіть сеанс. Вони також відмовляються працювати, поки зміну груп не завершено: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` виводить кожен обліковий запис із його слотом, короткою командою, електронною поштою й текою. Його масив `groups` перелічує групи облікового запису, а об'єкт `share` зіставляє кожен елемент (`mcp`, `settings`, `memory`, `extensions`, `history`) з групою, з якої він береться, або з `local`.

## Короткі команди

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Коротка команда є командою з назвою, яку ви вибираєте самі, і запускає Claude Code від імені одного облікового запису. Як і `ccam work`, вона передає всі аргументи в Claude Code без змін, тож `ccwork --resume` і `ccwork -p "hi"` працюють так, як ви очікуєте. Це посилання на ccam у `~/.local/bin`.

Короткі команди необов'язкові. Майстер `ccam add` запитує, чи створити її, `ccam add <name> --shortcut <command>` створює її без запитань, а `ccam list` показує коротку команду кожного облікового запису або `-`, якщо її немає.

У назвах можна використовувати символи a-z, 0-9, `.`, `-` і `_`. Назва не може збігатися з командою ccam, з короткою командою іншого облікового запису чи з назвою іншої програми у вашому `PATH`, а також із файлом у `~/.local/bin`, який створив не ccam, наприклад вашою власною копією ccam чи посиланням на нього. Коротка команда, вручну змінена в `config.toml`, яка порушує ці правила, є `claude` або `ccam` чи призначена двом обліковим записам, не дає завантажити `config.toml`, і ccam називає обліковий запис.

`ccam rename` зберігає назву короткої команди, а `ccam remove` видаляє її разом з обліковим записом. `ccam shims install` повторно створює посилання всіх облікових записів, наприклад після того, як ви перемістили ccam; файл, що заважає одній короткій команді, не зупиняє решту.

## Спільне використання

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Група спільно використовує елементи між своїми обліковими записами; обліковий запис може бути в кількох групах, але кожен елемент береться лише з однієї групи. Див. [спільне використання через групи](/features/share).

## Копіювання наявного налаштування

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` копіює `settings.json`, `CLAUDE.md`, розширення та користувацькі сервери MCP з теки Claude Code до групи, щоб їх отримав кожен обліковий запис, який спільно використовує ці елементи. Тека лише читається.

- `--settings`, `--memory`, `--extensions`, `--mcp` і `--history` копіюють лише вказані елементи. Без них у терміналі ви вибираєте елементи в контрольному списку, де позначено все, крім історії розмов; в інших випадках копіюється все, крім історії розмов.
- `--mcp` працює й тоді, коли `~/.claude` порожня або її немає, але в `~/.claude.json` є сервери.
- `-g` або `--group` вибирає групу.
- `-y` або `--yes` копіює без запитань і залишає копію групи, якщо файл відрізняється.

Входи MCP не копіюються, бо скопійований вхід виводить із системи інших його власників; план називає облікові записи, які мають увійти через `/mcp`. Див. [копіювання наявного налаштування](/docs/sharing#copy-an-existing-setup).

## Сервери MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` переміщує до групи облікового запису сервери, які є лише в ньому, зокрема ті, що в нього є лише для однієї теки (звичайний `claude mcp add`), і після цього вони завантажуються в усіх теках. Сервер, назву якого група вже використовує для іншого визначення, пропускається; додайте `--force`, щоб замінити копію групи, а ccam спершу створить її резервну копію.

## Використання та перемикання

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

- `next --auto`, `--fallback` і `ccam pick` вибирають однаково: спершу облікові записи нижче `auto_max_usage` (`config.toml`, за замовчуванням 100), за стратегією `auto_strategy` (`headroom`, `expiring` або `balance`), з урахуванням тижневого вікна `--model`. Вони ніколи не вибирають обліковий запис на 100% (вікно, про заповнення якого повідомив Claude Code, вважається заповненим на 100% до свого скидання) чи вхід із ключем API і лишаються в межах організації облікового запису, з якого починають (`--cross-org` або `cross_org = true`, щоб шукати ширше).
- `--wait` приймає значення після `=` (`--wait=2h`); записаний через пробіл, він відхиляється з підказкою. Коли вибрати нікого, вони завершуються з кодом 75 або, з `--wait`, чекають на скидання першого облікового запису.
- `ccam next --continue` переносить розмову навіть між обліковими записами, які не мають спільної історії: тоді Claude Code копіює її в новий обліковий запис. З `--brief` новий обліковий запис натомість починає нову розмову з підсумку попередньої, який ccam пише без моделі ([`--brief`](/features/next#brief)).
- На відміну від `ccam <name>`, `ccam next` читає власні параметри, тож параметри Claude Code ставляться після `--`; без нього ccam зупиняється з `unknown flag`.
- `ccam agent` виконує завдання через `claude -p` у цій теці від імені іншого облікового запису, ніж той, з якого його викликано, у режимі дозволів і з рівнем зусиль сеансу, що його викликав, і виводить останню відповідь Claude. Коли цей обліковий запис досягає ліміту, розмова триває від імені іншого ([Субагенти в інших облікових записах](/features/agent)).
- Усередині Claude Code `/ccam:switch [account]` із плагіна передає розмову далі, а `ccam hook` — це те, що запускають його хуки.
- З `auto_switch = true` у `config.toml` (Linux і macOS; типово вимкнено) хук ліміту ще й завершує інтерактивний сеанс, який запустив ccam, і продовжує його в іншому обліковому записі, де є запас, але лише того самого типу: іншому місці тієї самої організації або іншому особистому плані ([`auto_switch`](/features/next#auto-switch)).

Див. [Перемикання на наступний обліковий запис](/features/next), [Вибір облікового запису для скриптів](/features/pick), [Субагенти в інших облікових записах](/features/agent) і [Фонові сеанси всіх облікових записів](/features/agents).

## Обслуговування

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

`ccam update` перевіряє контрольну суму SHA-256 випуску, перш ніж замінити ccam, а `--version v0.2.0` встановлює певний випуск. Копія ccam, встановлена з пакета `.deb`, натомість оновлюється через apt. Див. [оновлення](/docs/install#update).

`ccam uninstall` показує, що саме видалить, і спершу питає: кожен обліковий запис (так, як це робить `ccam remove`), короткі команди, теки ccam і команду `ccam`. Імпортована тека залишається там, де вона є. Див. [видалення](/docs/install#uninstall).

## Файли

| Шлях | Вміст |
| --- | --- |
| `~/.config/ccam/config.toml` | Облікові записи, групи та налаштування (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). Налаштування розміщуються вгорі, над першим `[[group]]`, `[[profile]]` або `[[folder]]`: налаштування під таблицею належить цій таблиці, тому ccam попереджає про нього під час кожного запуску та в `ccam doctor`, а команди, які зберігають конфігурацію, відмовляються працювати, доки його не перенесено. Коментарі не зберігаються, коли ccam записує файл |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` кожного облікового запису |
| `~/.local/share/ccam/groups/<group>/` | Спільні налаштування, сервери MCP, розширення та історія |
| `~/.local/state/ccam/` | Кеш використання, записи про ліміти, блокування, нотатки та інший стан (див. список під таблицею) |
| `~/.local/bin/<shortcut>` | Короткі команди, наприклад `ccwork`; посилання на ccam |

Тека стану містить:

- кеш використання й записи про ліміти (`usage/<name>.json`, `usage/<name>.limit.json`);
- час останнього використання;
- позначки сеансів і блокування (`locks/`);
- короткі команди, створені ccam (`shortcuts.json`);
- незавершену зміну груп (`group-change.json`);
- нотатки, які залишають `/ccam:switch` і хук ліміту (`handoff/<account>/<session>.json`);
- час, коли `ccam pick` або `auto_switch` востаннє вибрав кожен обліковий запис (`picks/<name>`; вибір діє 2 хвилини);
- підсумки `--brief`, що зберігаються 7 днів (`brief/<session>.md`).

Якщо задати `XDG_CONFIG_HOME`, `XDG_DATA_HOME` або `XDG_STATE_HOME`, розташування цих тек зміниться. У Windows `~` — це `%USERPROFILE%`, а короткі команди — це `<shortcut>.exe`.
