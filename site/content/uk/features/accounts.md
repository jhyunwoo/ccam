---
title: Керування обліковими записами
description: "Додавайте, входьте, переглядайте, перейменовуйте, видаляйте й імпортуйте облікові записи Claude Code: ccam add, login, list, rename, remove і import."
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

Кожен обліковий запис ccam є окремою текою конфігурації Claude Code із власним входом. За допомогою цих команд можна додавати облікові записи, виконувати в них вхід, виводити їх список, перейменовувати й видаляти їх. Вони також реєструють теки, які у вас уже є.

## Використання

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

Там, де команда приймає `<name|number>`, можна вказати назву облікового запису або номер його слота (стовпець `#` у `ccam list`).

`ccam use` лише змінює, який обліковий запис є обліковим записом за замовчуванням (`default` у `config.toml`).

`ccam add` у терміналі без параметрів відкриває майстер. Він проходить такі кроки:

1. назва облікового запису, якщо ви її не вказали;
2. спосіб входу облікового запису: вхід Claude зараз, ключ API або вхід Claude пізніше;
3. групи, до яких приєднатися, з позначеною `default` (якщо груп ще немає, показано, що спільно використовуватиме група `default`);
4. `Create a shortcut command for work? [y/N]`;
5. підсумок для підтвердження.

`Esc` повертає на крок назад. Будь-який із параметрів `--group`, `--no-group`, `--shortcut` чи `--no-login` або запуск без термінала (наприклад, із конвеєра) пропускає майстер. Тоді назва обов'язкова, а коротка команда не створюється, якщо `--shortcut` її не вказує.

Без `--group` новий обліковий запис приєднується до групи `default`. Перший обліковий запис створює її, і вона спільно використовує сервери MCP, налаштування, `CLAUDE.md` і розширення; див. [спільне використання через групи](/features/share). У Windows файли можна пов'язувати лише в режимі розробника. Без нього перший обліковий запис створює `default` лише з `mcp` і `extensions` і повідомляє про це (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), а ccam відхиляє `--group` із групою, яка спільно використовує налаштування, пам'ять або історію.

`ccam list` показує слот, назву, коротку команду, електронну пошту, групи й останнє використання, а перед обліковим записом за замовчуванням стоїть `*`. Електронна пошта — це та, від імені якої тека увійшла зараз (з її `.claude.json`); коли вона відрізняється від записаної ccam, записана йде слідом як `(was …)`:

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

З `--json` кожен обліковий запис є об'єктом із полями `slot`, `shortcut`, `name`, `email` (поточного входу, як у таблиці), `default`, `dir`, `last_used` (відсутнє, якщо обліковий запис ще ніколи не запускали), `groups` і `share`. `share` зіставляє кожен елемент (`mcp`, `settings`, `memory`, `extensions`, `history`) з групою, з якої він береться, або з `"local"`.

## Параметри

| Параметр | Дія | За замовчуванням |
| --- | --- | --- |
| `ccam add --no-login` | Створити обліковий запис без входу; увійти пізніше через `ccam login <name>` | вимкнено |
| `ccam add --group <group>` | Приєднатися до цієї групи замість `default`. Можна повторювати; `default` створюється, якщо її немає | `default` |
| `ccam add --no-group` | Не приєднуватися до жодної групи: обліковий запис нічого не використовує спільно | вимкнено |
| `ccam add --shortcut <command>` | Також створити цю коротку команду в `~/.local/bin` | немає |
| `ccam add --api --key-from <source>` | Обліковий запис, який використовує ключ API замість входу; див. [Облікові записи з ключем API](/features/api) | вимкнено |
| `ccam list --json` | Вивести облікові записи у форматі JSON | вимкнено |
| `ccam remove -y`, `--yes` | Видалити без запитання | вимкнено |
| `ccam logout -y`, `--yes` | Вийти без запитання | вимкнено |
| `ccam import --name <name>` | Назва облікового запису для вказаної теки; без теки відхиляється (`--name needs a folder: ccam import <folder> --name ops`) | виводиться з назви теки |
| `ccam import -y`, `--yes` | Імпортувати кожну знайдену теку без запитання | вимкнено |
| `--dry-run` | Підтримується командами `ccam list` і `ccam completion`, які нічого не змінюють. Відхиляється командами `add`, `login`, `logout`, `use`, `rename`, `remove` і `import` з повідомленням `` --dry-run is not supported by `ccam add`; nothing was changed `` (з назвою команди) | вимкнено |

`ccam login` (окрім `--expired`), `ccam use` і `ccam rename` не мають власних параметрів. `ccam completion` приймає `bash`, `zsh` або `powershell`; будь-що інше завершується помилкою `only bash, zsh and powershell are supported`.

Назви облікових записів мають від 1 до 32 символів із `a-z`, `0-9`, `-` і `_` та починаються з літери або цифри. Назва не може складатися лише з цифр, бо `ccam 2` означає слот 2. Назва також не може збігатися з командою ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), тож `ccam <name>` ніколи не буває неоднозначним. Цей список міститься в `config.Reserved` у `internal/config`.

## Як це працює

Додавання облікового запису створює його теку, приєднує його до груп, зберігає конфігурацію, виконує вхід і створює коротку команду саме в такому порядку. Коли ви виконуєте `ccam add work --shortcut ccwork`, ccam робить такі кроки (`addProfile` у `internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. ccam перевіряє назву (`config.ValidateName`) і коротку команду (`shim.Check`), перш ніж щось створити.
2. ccam створює `~/.local/share/ccam/profiles/work/` з режимом `0700` і призначає обліковому запису найменший вільний номер слота, тож номер видаленого облікового запису використовується повторно (`profile.Add`, `Config.NextSlot`). Перший обліковий запис стає обліковим записом за замовчуванням.
3. Обліковий запис приєднується до своїх груп (`default`, якщо `--group` чи `--no-group` не вказує інше; перший обліковий запис створює `default`), а ccam пов'язує з ними нову порожню теку, наприклад `settings.json` з `groups/default/settings.json` (`applyShare`, через `share.Build`, а потім `Apply`). Якщо це не вдається, нова тека видаляється, бо залишена тека блокувала б кожну повторну спробу помилкою `folder already exists`.
4. Обліковий запис зберігається в `~/.config/ccam/config.toml` через `fsutil.AtomicWrite`: тимчасовий файл, `fsync`, потім `rename`, тож збій ніколи не залишає наполовину записану конфігурацію.
5. ccam запускає `claude auth login` зі змінною `CLAUDE_CONFIG_DIR`, що вказує на теку облікового запису, і з `--email` з останньою електронною поштою облікового запису, щоб сторінку входу було заповнено (`runLogin`). Якщо вхід не вдається, обліковий запис залишається, а ccam виводить помилку й радить пізніше виконати `ccam login work`.
6. ccam читає `oauthAccount` з `profiles/work/.claude.json` (`profile.ReadIdentity`): електронну пошту, ідентифікатори облікового запису й організації, назву організації та план (`claude_max`, `claude_team` …). Цей файл може мати кілька мегабайтів, тому ccam шукає один ключ замість розбору всього файлу (`TopLevelValue`). ccam зберігає ці дані в `[profile.login]` у `config.toml`.
7. ccam пов'язує `~/.local/bin/ccwork` з виконуваним файлом ccam, який під час запуску шукає `argv[0]` у `config.toml`, і записує посилання в `~/.local/state/ccam/shortcuts.json` (`shim.Install`). Якщо посилання створити не вдається, обліковий запис усе одно додається, ccam пояснює чому (`Could not create shortcuts: …`), а підсумок показує `shortcut: ccwork (not created; see above)`.

Вхід може тривати кілька хвилин, і за цей час інший процес ccam може зберегти конфігурацію. Тому команди, які чекають, повторно завантажують `config.toml`, застосовують одну зміну (тут електронну пошту) і зберігають (`env.update`). З тієї самої причини `ccam remove` повторно завантажує конфігурацію після свого запитання.

`ccam login` виконує кроки 5 і 6 для наявного облікового запису. Claude Code сприймає електронну пошту лише як підказку для сторінки входу, тож ви можете опинитися в іншому обліковому записі Claude, ніж раніше; ccam порівнює новий вхід із `[profile.login]` і повідомляє про це (`ccam: this login is other@company.com, but work was dev@company.com before.`). Він зберігає новий вхід і пояснює, як це скасувати: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

Вхід в інший обліковий запис Claude також скидає [запис про ліміт](/features/usage#the-limit-record) облікового запису, бо ті ліміти належали іншому обліковому запису. Коли інший обліковий запис ccam увійшов у той самий обліковий запис Claude, ccam попереджає, що вони мають спільні ліміти використання.

`ccam login --expired` по черзі виконує вхід у кожен обліковий запис, у який ще ніколи не входили або вхід якого прострочено: сплив термін дії його токена оновлення або Claude Code очистив вхід після того, як сервер його відхилив. Перш ніж відкриється перший браузер, він перевіряє, що `config.toml` можна зберегти.

`ccam logout` запускає `claude auth logout` в обліковому записі й очищає збережену електронну пошту; запис про вхід лишається, тож наступний вхід звіряється з ним. Claude Code відкликає вхід і видаляє те, що зберігає для облікового запису, зокрема входи серверів MCP. Тому ccam попереджає вас і спершу запитує (`--yes` пропускає запитання), коли:

- вхід містить входи MCP;
- обліковий запис запущено (його сеанси теж виходять із системи);
- тека іншого облікового запису містить той самий вхід (скопійована тека), який вихід теж відкликає.

Наприклад:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

Щоб увійти як інший обліковий запис Claude, достатньо самого `ccam login <name>`, і він зберігає входи MCP.

`ccam list` додає стовпець `Organization`, коли облікові записи належать до кількох організацій (особисті плани показують `personal`), а `ccam list --json` містить `plan` і `org` для кожного облікового запису.

### Вхід на віддаленому комп'ютері

Через SSH браузерна частина `claude auth login` може не спрацювати. Claude Code чекає на вхід на сервері зворотного виклику на віддаленій машині (`http://localhost:<port>/callback`), але браузер, який завершує вхід, працює на вашому комп'ютері, де ця адреса нікуди не веде, і браузер показує помилку на кшталт `ERR_CONNECTION_REFUSED`. Адреса все одно містить одноразовий код, тож ccam може його доставити:

<!-- motion:accounts-relay -->

ccam не передає ваш термінал у Claude Code (`runLogin`). Він під'єднує стандартне введення Claude Code до каналу й копіює туди кожен введений вами рядок (`relayLogin`), тож код, вставлений у запит Claude Code `Paste code here`, усе одно надходить. Коли рядок є адресою зворотного виклику, ccam приймає його, лише якщо хост — `localhost`, `127.0.0.1` або `[::1]`, вказано порт, шлях — `/callback` і присутні обидва параметри `code` і `state`; `http://` можна пропустити (`callbackURL`). Потім ccam запитує цю адресу на сервері з тайм-аутом у 10 секунд і без переходу за переспрямуваннями, як це зробив би браузер (`deliverCallback`). Запит надходить до сервера зворотного виклику `claude auth login`, що очікує (пунктирна лінія на діаграмі). Claude Code отримує код і зберігає вхід, на Linux — у `.credentials.json` у теці облікового запису. Вставлений рядок ніколи не може змусити ccam звернутися до іншого хоста.

Перед початком входу ccam виводить:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

У разі успіху він виводить `→ Passed the sign-in to Claude Code on this machine.`

Поки вхід чекає на введення в терміналі, ccam вимикає відлуння та построковий режим термінала й показує `*` для кожного введеного вами символу, тож вставлений код чи адреса зворотного виклику не з'являється на екрані. ccam сам редагує рядок (`readLine`): `Enter` надсилає рядок, `Backspace` стирає один символ, а `Ctrl+D` на порожньому рядку завершує введення. Термінал відновлюється, коли вхід завершується, а також у разі `Ctrl+C`, `SIGTERM` чи `SIGHUP`, які, як і раніше, завершують ccam. Якщо стандартне введення та стандартний потік помилок не під'єднано до термінала, рядки копіюються в міру надходження.

### Перейменування, видалення та імпорт

`ccam rename` і `ccam remove` діють на теку облікового запису та на те, що ccam зберігає під назвою облікового запису в інших місцях. ccam перелічує ці шляхи в `nameState`:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- кеші використання `~/.local/state/ccam/usage/<name>.json` та `<name>.statusline.json`;
- запис про ліміт `<name>.limit.json`;
- нотатки сеансів у `~/.local/state/ccam/handoff/<name>/`.

`ccam remove work` обробляє їх у такому порядку:

<!-- motion:accounts-remove -->

1. `ccam remove` спершу відмовляється, якщо обліковий запис запущено, потім запитує підтвердження й після вашої відповіді перевіряє ще раз (так виявляється сеанс, запущений, поки запитання було відкрите).
2. Перш ніж щось змінити, він перевіряє, що `config.toml` можна зберегти.
3. Якщо обліковий запис спільно використовує розширення, ccam переписує `installed_plugins.json` і `known_marketplaces.json` групи (`RewritePluginPaths`), бо Claude Code зберігає абсолютні шляхи до плагінів, які можуть проходити через теку цього облікового запису.
4. Далі з керованого облікового запису (того, теку якого створив ccam, на відміну від зареєстрованого через `ccam import`) виконується вихід: `claude auth logout` запускається в його теці, тож жоден робочий вхід не чекає в кошику (або, на macOS, в елементі в'язки ключів, названому за текою, який успадкував би пізніший обліковий запис із цією назвою). Цей крок пропускається для облікового запису з ключем API і тоді, коли тека іншого облікового запису містить той самий вхід, який вихід теж відкликав би. Про невдалий вихід повідомляється, а видалення триває.
5. ccam переміщує теку до `~/.local/share/ccam/trash/<name>-<time>/` (`profile.Trash`), обліковий запис зникає з `config.toml`, шляхи `nameState` видаляються, щоб пізніший обліковий запис із тією самою назвою їх не успадкував, а посилання короткої команди видаляється: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

Якщо `config.toml` усе ж не вдається зберегти, тека повертається з кошика. Щоб відновити видалений обліковий запис, перемістіть теку назад із кошика й знову увійдіть.

`ccam rename` перейменовує теку й переписує шляхи до плагінів у її `plugins/` зі старої теки на нову. Він переміщує шляхи `nameState` під нову назву й лишає обліковий запис обліковим записом за замовчуванням, якщо той ним був. Він також оновлює нотатки сеансів, щоб `ccam next --continue` і далі їх знаходив. Коротка команда зберігає свою назву. Якщо переписати шляхи до плагінів не вдається або `config.toml` не можна зберегти, теку перейменовують назад.

На macOS, а також у Windows, коли Claude Code зберігає входи в диспетчері облікових даних, сама тека зберігає стару назву, бо вхід і входи MCP зберігаються за шляхом теки, а переміщена тека втратила б вхід. Тоді конфігурація записує теку як `folder` поруч із новою назвою, а ccam каже `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

`ccam import` реєструє теку там, де вона є, як `dir` у `config.toml`. Він ніколи не переміщує й не копіює її: скопійований токен оновлення змушував би копії виводити одна одну із системи.

Без аргументу-теки він шукає незареєстровані теки `~/.claude` і `~/.claude-*` та запитує про кожну. Назви походять від теки (`~/.claude` — це `main`, `~/.claude-work` — `work`, `~/.claude-2` — `acct2`). Якщо ccam не може вивести назву або назву зайнято, він показує теку разом із командою з `--name`, яку слід виконати. Про знайдену теку, яку не вдалося імпортувати, повідомляється як `Skipped <path>: <reason>`, а решта імпортуються.

ccam спершу ставить усі запитання, потім повторно завантажує `config.toml` і реєструє теки, тож зміни, які тим часом зберегли інші команди ccam, лишаються. Рядки `✓` з'являються після збереження.

Без `CLAUDE_CONFIG_DIR` Claude Code зберігає `.claude.json` у вашій домашній теці поруч із `~/.claude`, тож коли ви імпортуєте `~/.claude`, а власного `.claude.json` у ній немає, ccam створює посилання `~/.claude/.claude.json` на `~/.claude.json`. Тоді звичайний `claude` і `ccam main` бачать той самий обліковий запис.

Якщо `~/.claude` має власний `.claude.json` поруч із `~/.claude.json`, ccam лишає обидва й повідомляє про це. `ccam main` використовував би внутрішній, а звичайний `claude` — зовнішній, тримаючи їхні дані входу й сервери MCP окремо. Щоб лишився один, замініть внутрішній файл посиланням ([doctor](/features/doctor) теж про це повідомляє).

Коли тека містить той самий вхід, що й зареєстрований обліковий запис (скопійована тека), ccam повідомляє про це: той, хто оновить його першим, виведе іншого із системи, тож виконайте `ccam login <name>`, перш ніж нею користуватися.

Імпортовані облікові записи не входять до жодної групи, доки ви не додасте їх до неї (`ccam group join <group> <name>`), і коротка команда не створюється.

`ccam completion bash|zsh` виводить скрипт автодоповнення, який також доповнює назви облікових записів.

## Безпека та граничні випадки

- `ccam rename` і `ccam remove` відмовляються працювати, поки обліковий запис запущено (`refuseRunning`), з повідомленням ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``, бо сеанс і далі записує в теку. Перевірка охоплює Claude Code, хоч би як його було запущено: через ccam, звичайним `claude`, з IDE, фонові сеанси та їхній демон.
- `ccam rename`, `ccam remove` і `ccam logout` відмовляються працювати, поки інший процес ccam оновлює або змінює цей вхід: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- Поки в іншому процесі ccam триває зміна груп, `ccam rename` і `ccam remove` відмовляються з `another ccam process is changing the groups right now; try again when it is done`; після зміни, що зупинилася на півдорозі, — з ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.
- `ccam remove` для імпортованого облікового запису зберігає теку та її вхід для звичайного `claude`. Якщо обліковий запис входить до груп, він спершу виходить із них так, як це робить `ccam group leave`: кожен спільний елемент стає в теці копією елемента групи, а сервери MCP групи потрапляють до її `.claude.json`. Потім його реєстрацію скасовують: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. Резервні копії її `.claude.json`, які створив ccam, переміщуються до `<folder>/.ccam-backup/mcp-<time>`; якщо їх не вдається перемістити, вони лишаються в `~/.local/state/ccam/mcp/<name>/backups`, і повідомлення про це каже. Ваша `~/.claude` не переміщується в кошик.
- Два облікові записи не можуть мати ту саму теку: `config.toml`, де одну теку названо двічі, не завантажується.
- `ccam add` відхиляє назву, тека якої вже існує в `profiles/`, і радить `ccam import`. На macOS він також зупиняється, коли в'язка ключів досі містить вхід або входи MCP для шляху нової теки, залишені старішою версією ccam, і називає команду `security delete-generic-password`, яка їх видаляє. Якщо `config.toml` не вдається зберегти, нова тека знову видаляється.
- `ccam import` відхиляє вже зареєстровану теку, шлях, який не є текою, і зайняту назву. ccam також попереджає, коли тека є копією `~/.claude` (той самий вхід), навіть якщо саму `~/.claude` не зареєстровано.
- `ccam logout` спершу запитує, коли обліковий запис запущено; `ccam login` цього не перевіряє.
- Якщо ретранслятор не може зв'язатися з Claude Code, наприклад після завершення входу, він виводить `ccam: could not pass the sign-in to Claude Code:` із причиною. Помилка HTTP зазвичай означає, що код прострочено; почніть вхід заново.

## Пов'язані сторінки

- [Запуск Claude Code](/features/launch)
- [Короткі команди](/features/shortcuts)
- [Спільні сервери MCP, налаштування та історія](/docs/sharing)
- [Кілька облікових записів Claude Code на одному ПК](/guides/multiple-claude-code-accounts)
