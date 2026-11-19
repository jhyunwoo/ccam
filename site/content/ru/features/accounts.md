---
title: Управление аккаунтами
description: "Добавляйте, входите, просматривайте, переименовывайте, удаляйте и импортируйте аккаунты Claude Code: ccam add, login, list, rename, remove и import."
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

Каждый аккаунт ccam представляет собой отдельную папку конфигурации Claude Code с собственным входом. С помощью этих команд можно добавлять аккаунты, выполнять в них вход, выводить их список, переименовывать и удалять их. Они также регистрируют уже существующие папки.

## Использование

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

Везде, где команда принимает `<name|number>`, можно указать имя аккаунта или номер его слота (столбец `#` в выводе `ccam list`).

`ccam use` только меняет аккаунт по умолчанию (`default` в `config.toml`).

`ccam add` в терминале без параметров открывает мастер. Он проходит такие шаги:

1. имя аккаунта, если вы его не указали;
2. способ входа аккаунта: вход Claude сейчас, ключ API или вход Claude позже;
3. группы, в которые нужно вступить, с отмеченной `default` (если групп ещё нет, мастер показывает, что будет делить группа `default`);
4. `Create a shortcut command for work? [y/N]`;
5. сводка для подтверждения.

`Esc` возвращает на шаг назад. Любой из параметров `--group`, `--no-group`, `--shortcut` или `--no-login`, а также запуск без терминала (например, из конвейера) отключают мастер. Тогда имя обязательно, а короткая команда создаётся, только если её имя задано через `--shortcut`.

Без `--group` новый аккаунт вступает в группу `default`. Первый аккаунт создаёт её, и она делит MCP-серверы, настройки, `CLAUDE.md` и расширения; см. раздел [общий доступ через группы](/features/share). В Windows создавать ссылки на файлы можно только в режиме разработчика. Без него первый аккаунт создаёт `default` только с `mcp` и `extensions` и сообщает об этом (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), а ccam отклоняет `--group` с группой, которая делит settings, memory или history.

`ccam list` показывает слот, имя, короткую команду, email, группы и время последнего использования, а перед аккаунтом по умолчанию ставит `*`. Email — тот, под которым папка вошла сейчас (из её `.claude.json`); если он отличается от записанного ccam, записанный следует за ним в виде `(was …)`:

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

С `--json` каждый аккаунт выводится как объект с полями `slot`, `shortcut`, `name`, `email` (текущего входа, как в таблице), `default`, `dir`, `last_used` (отсутствует, если аккаунт ни разу не запускался), `groups` и `share`. `share` сопоставляет каждому элементу (`mcp`, `settings`, `memory`, `extensions`, `history`) группу, из которой он берётся, или `"local"`.

## Параметры

| Параметр | Действие | По умолчанию |
| --- | --- | --- |
| `ccam add --no-login` | Создать аккаунт без входа; войти позже можно командой `ccam login <name>` | выкл. |
| `ccam add --group <group>` | Вступить в эту группу вместо `default`. Можно повторять; `default` создаётся, если её нет | `default` |
| `ccam add --no-group` | Не вступать ни в одну группу: аккаунт ничего не делит | выкл. |
| `ccam add --shortcut <command>` | Дополнительно создать эту короткую команду в `~/.local/bin` | нет |
| `ccam add --api --key-from <source>` | Аккаунт, который использует ключ API вместо входа; см. [Аккаунты с ключом API](/features/api) | выкл. |
| `ccam list --json` | Вывести аккаунты в формате JSON | выкл. |
| `ccam remove -y`, `--yes` | Удалить без вопросов | выкл. |
| `ccam logout -y`, `--yes` | Выйти без вопросов | выкл. |
| `ccam import --name <name>` | Имя аккаунта для указанной папки; без папки отклоняется (`--name needs a folder: ccam import <folder> --name ops`) | выводится из имени папки |
| `ccam import -y`, `--yes` | Импортировать все найденные папки без вопросов | выкл. |
| `--dry-run` | Поддерживается командами `ccam list` и `ccam completion`, которые ничего не меняют. Команды `add`, `login`, `logout`, `use`, `rename`, `remove` и `import` отклоняют его с сообщением `` --dry-run is not supported by `ccam add`; nothing was changed `` (с именем команды) | выкл. |

У `ccam login` (кроме `--expired`), `ccam use` и `ccam rename` нет собственных параметров. `ccam completion` принимает `bash`, `zsh` или `powershell`; любое другое значение завершается ошибкой `only bash, zsh and powershell are supported`.

Имя аккаунта содержит от 1 до 32 символов из `a-z`, `0-9`, `-` и `_` и начинается с буквы или цифры. Имя не может состоять только из цифр, потому что `ccam 2` означает слот 2. Имя также не может совпадать с командой ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), поэтому `ccam <name>` никогда не бывает неоднозначным. Этот список задан в `config.Reserved` в `internal/config`.

## Как это работает

Добавление аккаунта создаёт его папку, вступает в его группы, сохраняет конфигурацию, выполняет вход и создаёт короткую команду — именно в таком порядке. Когда вы выполняете `ccam add work --shortcut ccwork`, ccam делает следующие шаги (`addProfile` в `internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. ccam проверяет имя (`config.ValidateName`) и короткую команду (`shim.Check`), прежде чем что-либо создать.
2. ccam создаёт `~/.local/share/ccam/profiles/work/` с правами `0700` и назначает аккаунту наименьший свободный номер слота, поэтому номер удалённого аккаунта используется повторно (`profile.Add`, `Config.NextSlot`). Первый аккаунт становится аккаунтом по умолчанию.
3. Аккаунт вступает в свои группы (`default`, если `--group` или `--no-group` не указывают иное; первый аккаунт создаёт `default`), и ccam связывает с ними новую пустую папку, например `settings.json` с `groups/default/settings.json` (`applyShare`, через `share.Build` и затем `Apply`). Если это не удалось, новая папка удаляется, потому что оставшаяся папка блокировала бы каждую повторную попытку ошибкой `folder already exists`.
4. Аккаунт сохраняется в `~/.config/ccam/config.toml` через `fsutil.AtomicWrite`: временный файл, `fsync`, затем `rename`, поэтому сбой никогда не оставляет конфигурацию записанной наполовину.
5. ccam запускает `claude auth login` с `CLAUDE_CONFIG_DIR`, указывающим на папку аккаунта, и с `--email`, содержащим последний email аккаунта, чтобы страница входа была заполнена (`runLogin`). Если вход не удался, аккаунт остаётся, а ccam выводит ошибку и предлагает позже выполнить `ccam login work`.
6. ccam читает `oauthAccount` из `profiles/work/.claude.json` (`profile.ReadIdentity`): email, идентификаторы аккаунта и организации, название организации и тариф (`claude_max`, `claude_team` …). Этот файл может занимать несколько мегабайт, поэтому ccam ищет один ключ, а не разбирает весь файл (`TopLevelValue`). ccam хранит эти данные в `[profile.login]` в `config.toml`.
7. ccam связывает `~/.local/bin/ccwork` с исполняемым файлом ccam, который при запуске ищет `argv[0]` в `config.toml`, и записывает ссылку в `~/.local/state/ccam/shortcuts.json` (`shim.Install`). Если ссылку создать не удалось, аккаунт всё равно добавляется, ccam объясняет причину (`Could not create shortcuts: …`), а сводка показывает `shortcut: ccwork (not created; see above)`.

Вход может занять несколько минут, и за это время другой процесс ccam может сохранить конфигурацию. Поэтому команды, которые ждут, заново читают `config.toml`, применяют одно изменение (здесь email) и сохраняют (`env.update`). По той же причине `ccam remove` заново читает конфигурацию после своего вопроса.

`ccam login` выполняет шаги 5 и 6 для существующего аккаунта. Claude Code использует email только как подсказку для страницы входа, поэтому вы можете оказаться в другом аккаунте Claude, чем раньше; ccam сравнивает новый вход с `[profile.login]` и сообщает об этом (`ccam: this login is other@company.com, but work was dev@company.com before.`). Он сохраняет новый вход и объясняет, как его отменить: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

Вход в другой аккаунт Claude также удаляет [запись о лимите](/features/usage#the-limit-record) аккаунта, поскольку те лимиты относились к другому аккаунту. Когда другой аккаунт ccam вошёл в тот же аккаунт Claude, ccam предупреждает, что у них общие лимиты использования.

`ccam login --expired` по очереди выполняет вход во все аккаунты, в которые вход ни разу не выполнялся или чей вход истёк: истёк его токен обновления или Claude Code очистил вход после того, как сервер его отклонил. Перед тем как откроется первый браузер, команда проверяет, что `config.toml` можно сохранить.

`ccam logout` запускает `claude auth logout` в аккаунте и удаляет сохранённый email; запись о входе остаётся, поэтому следующий вход сверяется с ней. Claude Code отзывает вход и удаляет то, что хранит для аккаунта, в том числе данные входа MCP-серверов. Поэтому ccam сначала предупреждает и спрашивает (`--yes` пропускает вопрос), когда:

- вход содержит данные входа MCP;
- аккаунт запущен (его сессии тоже выходят из системы);
- папка другого аккаунта содержит тот же вход (скопированная папка), который выход тоже отзывает.

Например:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

Чтобы войти как другой аккаунт Claude, достаточно одного `ccam login <name>`, и он сохраняет данные входа MCP.

`ccam list` добавляет столбец `Organization`, когда аккаунты принадлежат более чем одной организации (личные тарифы показываются как `personal`), а в `ccam list --json` у каждого аккаунта есть `plan` и `org`.

### Вход на удалённой машине

Через SSH браузерная часть `claude auth login` может не сработать. Claude Code ждёт завершения входа на сервере обратного вызова на удалённой машине (`http://localhost:<port>/callback`), но браузер, в котором завершается вход, работает на вашем компьютере, где этот адрес никуда не ведёт, и браузер показывает ошибку, например `ERR_CONNECTION_REFUSED`. Адрес при этом всё равно содержит одноразовый код, поэтому ccam может его доставить:

<!-- motion:accounts-relay -->

ccam не отдаёт ваш терминал Claude Code (`runLogin`). Он подключает стандартный ввод Claude Code к каналу и копирует туда каждую введённую вами строку (`relayLogin`), поэтому код, вставленный в ответ на приглашение Claude Code `Paste code here`, всё равно доходит. Если строка является адресом обратного вызова, ccam принимает её, только если хост равен `localhost`, `127.0.0.1` или `[::1]`, порт указан, путь равен `/callback` и присутствуют и `code`, и `state`; `http://` можно опустить (`callbackURL`). Затем ccam запрашивает этот адрес на сервере с тайм-аутом 10 секунд и без перехода по перенаправлениям, как сделал бы браузер (`deliverCallback`). Запрос попадает на сервер обратного вызова ожидающего `claude auth login` (пунктирная линия на схеме). Claude Code получает код и сохраняет вход, в Linux в `.credentials.json` в папке аккаунта. Вставленная строка никогда не может заставить ccam обратиться к другому хосту.

Перед началом входа ccam выводит:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

При успехе он выводит `→ Passed the sign-in to Claude Code on this machine.`

Пока вход ждёт ввода в терминале, ccam отключает эхо и построчный режим терминала и показывает `*` для каждого введённого вами символа, поэтому вставленный код или адрес обратного вызова не появляется на экране. ccam сам редактирует строку (`readLine`): `Enter` отправляет строку, `Backspace` стирает один символ, а `Ctrl+D` в пустой строке завершает ввод. Терминал восстанавливается, когда вход завершается, а также при `Ctrl+C`, `SIGTERM` или `SIGHUP`, которые по-прежнему завершают ccam. Если стандартный ввод и стандартный поток ошибок не подключены к терминалу, строки копируются по мере поступления.

### Переименование, удаление и импорт

`ccam rename` и `ccam remove` работают с папкой аккаунта и с тем, что ccam хранит под именем аккаунта в других местах. ccam перечисляет эти пути в `nameState`:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- кэши использования `~/.local/state/ccam/usage/<name>.json` и `<name>.statusline.json`;
- запись о лимите `<name>.limit.json`;
- заметки сессий в `~/.local/state/ccam/handoff/<name>/`.

`ccam remove work` обрабатывает их в таком порядке:

<!-- motion:accounts-remove -->

1. Сначала `ccam remove` отказывается работать, если аккаунт запущен, затем задаёт вопрос и проверяет ещё раз после вашего ответа (так замечается сессия, запущенная, пока вопрос был открыт).
2. Прежде чем что-либо менять, он проверяет, что `config.toml` можно сохранить.
3. Если аккаунт делит расширения, ccam переписывает файлы группы `installed_plugins.json` и `known_marketplaces.json` (`RewritePluginPaths`), потому что Claude Code хранит абсолютные пути к плагинам, которые могут проходить через папку этого аккаунта.
4. Затем из управляемого аккаунта (папку которого создал ccam, в отличие от аккаунта, зарегистрированного через `ccam import`) выполняется выход: `claude auth logout` запускается в его папке, чтобы в корзине не лежал рабочий вход (а в macOS — в элементе Связки ключей, названном по папке, который унаследовал бы будущий аккаунт с тем же именем). Это пропускается для аккаунта с ключом API и когда папка другого аккаунта содержит тот же вход, который выход тоже отозвал бы. О неудачном выходе сообщается, и удаление продолжается.
5. ccam перемещает папку в `~/.local/share/ccam/trash/<name>-<time>/` (`profile.Trash`), аккаунт удаляется из `config.toml`, пути `nameState` удаляются, чтобы будущий аккаунт с тем же именем их не унаследовал, а ссылка короткой команды удаляется: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

Если `config.toml` всё же не удаётся сохранить, папка возвращается из корзины. Чтобы восстановить удалённый аккаунт, верните папку из корзины и войдите снова.

`ccam rename` переименовывает папку и переписывает пути плагинов в её `plugins/` со старой папки на новую. Он переносит пути `nameState` на новое имя и сохраняет статус аккаунта по умолчанию, если аккаунт им был. Он также обновляет заметки сессий, чтобы `ccam next --continue` по-прежнему их находил. Короткая команда сохраняет своё имя. Если переписать пути плагинов не удалось или `config.toml` нельзя сохранить, папке возвращается прежнее имя.

В macOS, а также в Windows, когда Claude Code хранит входы в диспетчере учётных данных, сама папка сохраняет старое имя, потому что вход и данные входа MCP хранятся под путём папки, и перемещённая папка вышла бы из системы. Тогда конфигурация записывает папку как `folder` рядом с новым именем, и ccam сообщает `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

`ccam import` регистрирует папку там, где она находится, как `dir` в `config.toml`. Он никогда не перемещает и не копирует её: скопированный токен обновления заставил бы копии выводить друг друга из аккаунта.

Без аргумента с папкой он ищет незарегистрированные папки `~/.claude` и `~/.claude-*` и спрашивает про каждую. Имена берутся из папки (`~/.claude` становится `main`, `~/.claude-work` становится `work`, `~/.claude-2` становится `acct2`). Если ccam не может вывести имя или имя занято, он выводит папку вместе с командой `--name`, которую нужно выполнить. О найденной папке, которую не удалось импортировать, сообщается как `Skipped <path>: <reason>`, а остальные импортируются.

ccam сначала задаёт все вопросы, затем заново читает `config.toml` и регистрирует папки, поэтому изменения, которые другие команды ccam сохранили за это время, остаются. Строки `✓` выводятся после сохранения.

Без `CLAUDE_CONFIG_DIR` Claude Code хранит `.claude.json` в домашней папке рядом с `~/.claude`, поэтому, когда вы импортируете `~/.claude` и в ней нет собственного `.claude.json`, ccam создаёт ссылку `~/.claude/.claude.json` на `~/.claude.json`. После этого обычный `claude` и `ccam main` видят один и тот же аккаунт.

Если в `~/.claude` есть собственный `.claude.json` рядом с `~/.claude.json`, ccam оставляет оба и сообщает об этом. `ccam main` использовал бы внутренний файл, а обычный `claude` — внешний, и их данные входа и MCP-серверы были бы раздельными. Чтобы остался один, замените внутренний файл ссылкой ([doctor](/features/doctor) тоже об этом сообщает).

Когда папка содержит тот же вход, что и зарегистрированный аккаунт (скопированная папка), ccam сообщает об этом: тот, кто обновит его первым, выведет другого из системы, поэтому выполните `ccam login <name>`, прежде чем ею пользоваться.

Импортированные аккаунты не входят ни в одну группу, пока вы не добавите их в группу (`ccam group join <group> <name>`), и короткая команда для них не создаётся.

`ccam completion bash|zsh` выводит скрипт автодополнения, который также дополняет имена аккаунтов.

## Безопасность и особые случаи

- `ccam rename` и `ccam remove` отказываются работать, пока аккаунт запущен (`refuseRunning`), с сообщением ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``, потому что сессия продолжает записывать в папку. Проверка охватывает Claude Code, как бы он ни был запущен: через ccam, обычным `claude`, из IDE, фоновые сессии и их демон.
- `ccam rename`, `ccam remove` и `ccam logout` отказываются работать, пока другой процесс ccam обновляет или меняет этот вход: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- Пока в другом процессе ccam идёт изменение групп, `ccam rename` и `ccam remove` отказываются работать с сообщением `another ccam process is changing the groups right now; try again when it is done`; после изменения, остановившегося на полпути, — с ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.
- `ccam remove` для импортированного аккаунта сохраняет папку и её вход для обычного `claude`. Если аккаунт входит в группы, он сначала выходит из них так же, как при `ccam group leave`: каждый общий элемент становится копией элемента группы в папке, а MCP-серверы группы попадают в её `.claude.json`. Затем регистрация отменяется: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. Резервные копии её `.claude.json`, созданные ccam, перемещаются в `<folder>/.ccam-backup/mcp-<time>`; если переместить их не удаётся, они остаются в `~/.local/state/ccam/mcp/<name>/backups`, и сообщение об этом говорит. Ваша `~/.claude` не перемещается в корзину.
- У двух аккаунтов не может быть одной папки: `config.toml`, в котором одна папка указана дважды, не загружается.
- `ccam add` отклоняет имя, папка для которого уже существует в `profiles/`, и предлагает `ccam import`. В macOS он также останавливается, когда в Связке ключей ещё есть вход или данные входа MCP для пути новой папки, оставленные более старой версией ccam, и называет команду `security delete-generic-password`, которая их удаляет. Если `config.toml` не удаётся сохранить, новая папка снова удаляется.
- `ccam import` отклоняет уже зарегистрированную папку, путь, который не является папкой, и занятое имя. ccam также предупреждает, когда папка является копией `~/.claude` (тот же вход), даже если сама `~/.claude` не зарегистрирована.
- `ccam logout` сначала спрашивает, если аккаунт запущен; `ccam login` этого не проверяет.
- Если передать вход в Claude Code не удаётся, например после завершения входа, выводится `ccam: could not pass the sign-in to Claude Code:` с причиной. Ошибка HTTP обычно означает, что срок действия кода истёк; начните вход заново.

## См. также

- [Запуск Claude Code](/features/launch)
- [Короткие команды](/features/shortcuts)
- [Общие MCP-серверы, настройки и история](/docs/sharing)
- [Несколько аккаунтов Claude Code на одном компьютере](/guides/multiple-claude-code-accounts)
