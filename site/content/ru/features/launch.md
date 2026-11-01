---
title: Запуск Claude Code
description: Запускайте Claude Code от имени любого аккаунта через ccam <name|number>, ccam run или панель; аргументы идут в claude как есть, --dry-run выводит команду.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: "11df08608a1e"
---

`ccam work` запускает Claude Code от имени аккаунта `work`, а всё, что идёт после имени, передаётся в `claude` без изменений. Используйте его везде, где вы набрали бы `claude`: в терминале, в скрипте или с вводом через конвейер.

## Использование

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

Укажите имя аккаунта или его номер из `ccam list`. ccam не читает аргументы после имени, поэтому `ccam work --help` показывает справку Claude Code. Указывать `--` после имени необязательно. Стандартный ввод, код возврата и сигналы работают так же, как с `claude`.

Если первый аргумент начинается с `-` и не является собственным флагом ccam, ccam сам выбирает аккаунт и передаёт все аргументы в Claude Code. Он использует аккаунт, заданный для текущей папки командой [`ccam folder`](/features/folder), а если его нет — аккаунт по умолчанию: заданный командой `ccam use` или, если вы его не задали, аккаунт с наименьшим номером.

### Панель

`ccam` без аргументов открывает [панель](/features/dashboard), если стандартный ввод и вывод подключены к терминалу, а в остальных случаях выводит `ccam list`. Перейдите к аккаунту (стрелками или набрав его номер) и нажмите `Enter`, чтобы запустить его; `n` запускает аккаунт с наибольшим запасом, а `o` предлагает `--continue`, `--resume` или другие аргументы. Панель закрывается до запуска Claude Code, поэтому терминал сначала восстанавливается.

**Места в тарифах Team и Enterprise.** Автоматический выбор аккаунта (`ccam next`, `--fallback`, `ccam pick`, клавиша `n` на панели и, если вы его включите, [`auto_switch`](/features/next#auto-switch)) переключается и между местами одной и той же организации Claude Team или Enterprise, что может нарушать условия Anthropic. Прочитайте [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) и уточните это у своей организации и у Anthropic, прежде чем на него полагаться.

Если в оболочке заданы `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` или `ANTHROPIC_AUTH_TOKEN`, Claude Code может использовать этот ключ или токен вместо собственного входа аккаунта, поэтому ccam предупреждает вас об этом в стандартном потоке ошибок перед запуском (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Он только читает эти переменные; [аккаунт с ключом API](/features/api) вместо этого их очищает.

По умолчанию ccam заменяет себя на Claude Code, поэтому больше ничего не продолжает работать. С `relay = true` (или, в Linux и macOS, с [`auto_switch = true`](/features/next#auto-switch)) в начале `config.toml`, над первым `[[group]]`, `[[profile]]` или `[[folder]]`, ccam вместо этого запускает интерактивные сессии как дочерний процесс и продолжает работать рядом с ними. Это позволяет `/ccam:switch` или `auto_switch` на лимите использования самому продолжить разговор в другом аккаунте вместе с параметрами сессии, режимом разрешений, уровнем усилий и агентом; прежде чем `/ccam:switch` перенесёт его в аккаунт другой организации, ccam спрашивает вас ([подробнее](/features/next#relay)).

## Параметры

| Параметр | Действие | По умолчанию |
| --- | --- | --- |
| `--dry-run` | Вывести окружение и команду, которые ccam выполнил бы, и завершиться | выкл. |
| `--fallback` | Для запусков с `-p`/`--print`: когда аккаунт достигает лимита использования, снова выполнить ту же команду от имени аккаунта с входом и наибольшим запасом ([подробнее](/features/next#run-again-at-the-limit)) | выкл. |
| `-v`, `--version` | Вывести `ccam <version>` (только как первый аргумент) | |
| `-h`, `--help` | Справка ccam, как `ccam --help` или `ccam run --help` | |

Правила разбора аргументов:

- ccam читает собственные флаги только до имени аккаунта. После имени каждый аргумент принадлежит Claude Code, включая `--dry-run` и `--help`.
- `--dry-run` работает в виде `ccam --dry-run work`, `ccam run --dry-run work` и `ccam --dry-run` (для аккаунта папки или аккаунта по умолчанию). `--dry-run` и `--fallback` можно указывать в любом порядке.
- `ccam run <name|number> [claude args]` является полной формой `ccam <name>`. `ccam run` без имени использует аккаунт папки или аккаунт по умолчанию.
- Имена аккаунтов не могут совпадать с именами команд ccam, например `add`, `list`, `run` или `usage`, поэтому `ccam <word>` никогда не бывает неоднозначным.

Пробный запуск выводит одну строку, которую можно вставить в оболочку, а в Windows — в PowerShell. Она начинается с того, что запуск меняет в окружении: всегда `CLAUDE_CONFIG_DIR` и `CCAM_PROFILE`, `env -u <name>` для переменной, которую он удаляет (например, унаследованной `CCAM_RELAY`), а для [аккаунта с ключом API](/features/api) — его эндпоинт, его модель и то, откуда берётся ключ (`ANTHROPIC_API_KEY="$WORK_KEY"`), но никогда не сам ключ:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

Строка PowerShell после этого возвращает переменные на место, как делает `VAR=x cmd` в оболочке, и передаёт Claude Code его командную строку через `Start-Process`, поэтому Windows PowerShell 5.1 и PowerShell 7 выполняют её одинаково. Когда запуск задаёт переменной пустое значение, строка `#` под ней сообщает, что PowerShell до версии 7.5 вместо этого удаляет такую переменную. Для `claude.cmd`, установленного через npm, строка выполняет команду `cmd.exe`, которую использовал бы запуск (см. ниже), а аргумент, который запуск там отклонил бы, приводит к ошибке и пробного запуска.

## Как это работает

Каждый запуск проходит одни и те же шаги, от командной строки до `exec`:

<!-- motion:launch -->

Каждый блок подсвечивается, когда выполняется его шаг, а строки под ним показывают окружение и аргументы, с которыми запускается Claude Code.

1. ccam начинает работу в `app.Main`, который вызывает `cmd/ccam/main.go`. Если ccam запущен под другим именем, это [короткая команда](/features/shortcuts), и она сразу передаётся запуску.
2. ccam просматривает `os.Args` до того, как подключается cobra (`app.parse`). В `cli.Execute` попадают только `ccam` без аргументов (панель), имя команды ccam (`config.Reserved`), `-h`/`--help` и `-v`/`--version`. Всё остальное считается запуском.
3. ccam читает `~/.config/ccam/config.toml`, единственный файл, который он разбирает (`launcher.Launch`).
4. ccam находит аккаунт по имени или номеру, отбрасывает один начальный `--` и находит исполняемый файл `claude` (`launcher.Prepare`, `ResolveClaude`). Без имени он берёт аккаунт папки, сравнивая рабочий каталог с записями `[[folder]]` из `config.toml` как строки, без учёта регистра в macOS и Windows (`HereProfile`). Результат поиска в PATH кэшируется в `~/.local/state/ccam/claude-path` вместе со временем изменения файла и хэшем `PATH`, поэтому ccam не ищет в `PATH` при каждом запуске. Если файл изменился или отсутствует или изменился `PATH`, поиск выполняется заново. Путь `claude_bin` в `config.toml` (`~/` раскрывается) должен быть исполняемым файлом; иначе запуск останавливается с сообщением `claude_bin X is not an executable file. Check claude_bin in config.toml`, прежде чем что-либо будет записано. ccam отклоняет `claude`, который является самим ccam (сравнение по inode), потому что это привело бы к зацикливанию.
5. ccam сравнивает папку аккаунта `~/.local/share/ccam/profiles/<name>/` с её настройками общего доступа, вызывая только `lstat` и `readlink`, один или два системных вызова на элемент (`checkLinks`, `share.Check`). Отсутствующая ссылка создаётся на месте. Всё, что может привести к потере данных, останавливает запуск.
6. ccam добавляет общие MCP-серверы (`applyMCP`; см. ниже).
7. ccam регистрирует сессию (`session.Register`). Он записывает маркер с именем по своему PID, содержащий время запуска процесса, в `~/.local/state/ccam/locks/profile-<name>/<pid>`, что бы ни указывала `XDG_RUNTIME_DIR`, поэтому все оболочки видят одни и те же маркеры. Для аккаунта, который делит историю разговоров, ccam сначала ищет запущенную сессию в `history-<group>-<hash>`, где хэш вычисляется от текущей папки с раскрытыми символическими ссылками, и тоже регистрируется там (`historyCheck`).

   `exec` сохраняет PID, поэтому маркер затем указывает на сам процесс Claude Code. ccam использует файл-маркер, а не удерживаемую блокировку, потому что процессы, запущенные Claude Code, унаследовали бы блокировку и сохранили бы её после завершения сессии. Если маркер не удаётся записать, запуск не останавливается: ccam сообщает `ccam: could not mark the session as running (…); starting Claude Code anyway`.

   Другие проверки ccam, запущен ли аккаунт (переименование, удаление, изменения групп, `doctor --fix`, синхронизация MCP), также читают собственные записи Claude Code в папке аккаунта (`sessions/<pid>.json`, `daemon.lock`), поэтому видят и Claude Code, запущенный без ccam: обычный `claude`, IDE, фоновые сессии и их демон.
8. ccam устанавливает заголовок терминала `claude · <name>` (если стандартный поток ошибок подключён к терминалу), записывает время запуска в `~/.local/state/ccam/last-used/<name>` без `fsync` и вызывает `syscall.Exec` (`Plan.Exec`).

Claude Code работает с `CLAUDE_CONFIG_DIR`, указывающим на папку аккаунта, и `CCAM_PROFILE`, равным имени аккаунта. ccam очищает унаследованные значения этих переменных, а также унаследованную `CCAM_RELAY`. Для папки `~/.claude`, зарегистрированной через `ccam import` в macOS, он также задаёт `CLAUDE_SECURESTORAGE_CONFIG_DIR` пустое значение, чтобы аккаунт и обычный `claude` использовали один и тот же вход из Связки ключей; в Windows он делает то же самое, когда Claude Code хранит вход в диспетчере учётных данных.

### Почему exec

ccam заменяет собственный процесс на `claude` (`syscall.Exec`). В результате:

- ccam не остаётся в памяти;
- нажатия клавиш, сигналы и код возврата принадлежат самому Claude Code;
- аргументы остаются массивом и никогда не становятся строкой оболочки.

В Windows нет `exec`. Там ccam запускает `claude.exe` как дочерний процесс с тем же окружением, аргументами и консолью, ждёт его и завершается с его кодом возврата (`Plan.Exec`). `claude.cmd`, установленный через npm, выполняется в `cmd.exe` с полным путём и каждым аргументом в кавычках, чтобы `cmd.exe` не прочитал ни один из них как собственный синтаксис (`&`, `%PATH%`); аргумент с переводом строки так передать нельзя, и он отклоняется с советом установить нативный Claude Code.

Ctrl+C доходит до Claude Code напрямую через консоль, поэтому ccam его игнорирует. Маркер сессии содержит PID ccam, который живёт ровно столько же, сколько Claude Code, и удаляется, когда Claude Code завершается. Пока ccam ждёт, он больше ничего не делает.

Работа после сессии, например обнаружение серверов, добавленных через `claude mcp add`, выполняется при следующем запуске: ccam сравнивает время изменения `.claude.json` с `~/.local/state/ccam/mcp/<name>/seen` и разбирает файл, только если он изменился.

### Как обрабатывается командная строка

Три командные строки и три пути через `internal/app`:

<!-- motion:launch-route -->

Короткая команда распознаётся по `argv[0]`, `ccam <name>` по первому аргументу, и оба пути заканчиваются `syscall.Exec`. До cobra доходит только команда ccam, например `ccam usage`.

### Общие MCP-серверы

Со стратегией по умолчанию `flag` ccam передаёт MCP-серверы группы в Claude Code в виде аргумента, который ставится перед вашими:

<!-- motion:launch-mcp -->

ccam добавляет аргумент, только если `mcp.json` группы существует, и ставит его сразу после `claude` (`mcp.FlagArgs`, `Plan.InsertArgs`). Ваши собственные `--mcp-config` или `--strict-mcp-config` при этом по-прежнему применяются поверх. Это один аргумент, `--mcp-config=<path>`, потому что `--mcp-config` в Claude Code принимает несколько значений: если записать его двумя аргументами, следующий за ним запрос тоже был бы прочитан как файл конфигурации.

ccam решает, куда флаг поставить нельзя (`claudeVerb`). Он читает подкоманду Claude Code из аргументов так же, как это делает Claude Code, пропуская параметры вместе с их значениями. `remote-control` или `rc` флаг не получает (ниже), как и команда фоновых сессий, указанная первой (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`). `ccam work mcp list`, `get` и `remove` тоже его не получают, поскольку там Claude Code его отклоняет; вместо этого ccam сообщает в стандартный поток ошибок, что серверы группы тоже загружаются в сессиях аккаунта:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Все остальные подкоманды и все сессии сохраняют флаг.

С `mcp_strategy = "sync"` в `config.toml` ccam вместо этого записывает серверы в `.claude.json` аккаунта и пропускает этот шаг, если аккаунт уже запущен. Запуск пропускает эту работу и тогда, когда с последней проверки не изменились ни `mcp.json` группы, ни `.claude.json`. Если у аккаунта есть собственный сервер с тем же именем, что и сервер группы, синхронизация сохраняет сервер аккаунта и один раз сообщает об этом. Синхронизация также возвращает копию, пропавшую из `.claude.json` (`put back group default's MCP servers …`), и меняет или удаляет только записанные ею копии, которые с тех пор никто не менял.

### Remote Control

`ccam work rc` (или `ccam work remote-control`) запускает Remote Control из Claude Code от имени аккаунта `work`. `claude remote-control` отклоняет `--mcp-config` перед своей командой и не принимает его после неё, а его сессии читают MCP-серверы только из `.claude.json` аккаунта. Поэтому стратегия с флагом здесь работает иначе:

<!-- motion:launch-rc -->

1. ccam находит `remote-control` или `rc` как команду Claude Code, даже после параметров и их значений (`ccam work --model opus rc`), и не добавляет флаг (`claudeVerb`, `mcp.FlagArgs`).
2. ccam сравнивает время изменения `mcp.json` группы с `sync-stamp` аккаунта и проверяет, что каждый сервер группы всё ещё есть в `.claude.json` (`mcp.Synced`). Если `.claude.json` уже содержит текущие серверы группы или у группы их нет, Remote Control запускается без вопроса.
3. В противном случае ccam называет серверы, которых не будет в сессиях, и спрашивает:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. При ответе `y` ccam копирует их так же, как стратегия `sync` (`mcp.Sync`): под блокировкой Claude Code, после создания резервной копии и с сохранением собственных серверов аккаунта. При ответе `N` или без терминала Remote Control запускается без них, а ccam указывает на `mcp_strategy = "sync"`, которая копирует их перед каждым запуском.

ccam спрашивает при каждом запуске Remote Control, пока серверы не скопированы. После копирования он спрашивает снова, только когда `mcp.json` группы изменится или скопированный сервер пропадёт из `.claude.json`. Копия в `.claude.json` побеждает `--mcp-config` при последующих обычных запусках, поэтому `ccam mcp add`, `remove` и `import` обновляют копии, которые хранятся у других аккаунтов ([Общие MCP-серверы](/features/mcp)).

### Бюджет в 10 мс

Бюджет пути запуска составляет 10 мс на 95-м процентиле; на практике ccam добавляет около 4 мс до запуска Claude Code. На пути запуска никогда не строится дерево команд cobra, не запускается Bubble Tea и не открывается сетевое соединение, а `TestPackageInitIsCheap` завершается ошибкой, если функции `init()` пакетов в сумме занимают больше 5 мс. `scripts/perf-check.sh` измеряет `ccam --dry-run work`, `ccam work` и короткую команду с помощью `hyperfine` и завершается ошибкой при превышении 10 мс, а `BenchmarkLaunchPath` в `internal/launcher` выполняет эти шаги в цикле. Чтобы увидеть время каждого шага, задайте `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Безопасность и особые случаи

- Неизвестное имя или номер останавливает запуск с сообщением `no account "<name>"` и ссылкой на `ccam list`. Если аккаунтов нет совсем, ccam предлагает добавить аккаунт командой `ccam add <name>`.
- Если общий путь стал обычным файлом или указывает не туда, ccam останавливается с сообщением `stopped because of share link problems:`, перечисляет проблемы и предлагает выполнить `ccam doctor --fix`. На пути запуска он не перемещает данные.
- Если в той же папке работает другая сессия Claude Code, которая делит историю разговоров с этим аккаунтом (та же группа), ccam предупреждает, что история может перемешаться, и спрашивает `Run anyway? [y/N]`. Папка сравнивается с раскрытыми символическими ссылками. Без терминала (стандартный ввод из `/dev/null` считается отсутствием терминала) или из фонового задания он показывает предупреждение и продолжает. Новый запуск `-p` начинает собственный разговор и предупреждения не получает; `--version`, `--help` и подкоманды, которые не открывают разговор, например `mcp`, `auth` или `plugin`, вообще не проверяются.
- `--dry-run` пропускает синхронизацию MCP, уведомление о новых серверах, заметку `mcp list`, проверку истории, маркер сессии, заголовок и время последнего использования и не выполняет exec. При этом он всё равно проверяет ссылки общего доступа и создаёт отсутствующую ссылку, как при обычном запуске.
- `ccam work rc`, когда `work` уже запущен, не трогает `.claude.json`, сообщает об этом и запускает Remote Control без серверов группы; чтобы скопировать их, выполните команду снова после завершения сессии. С `--dry-run` выводится сообщение, но ccam ничего не спрашивает и не копирует.
- С `mcp_strategy = "sync"` при неудачной синхронизации выводится `ccam: MCP sync failed (launching anyway): …`, и Claude Code запускается с теми серверами, которые у него уже были.

## См. также

- [Короткие команды](/features/shortcuts)
- [Remote Control с ccam](/guides/remote-control)
- [Управление аккаунтами](/features/accounts)
- [Справочник команд](/docs/commands)
- [Общие MCP-серверы, настройки и история](/docs/sharing)
- [Что такое CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
