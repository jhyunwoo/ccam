---
title: Uruchamianie Claude Code
description: Uruchom Claude Code na dowolnym koncie przez ccam <name|number>, ccam run lub panel; argumenty idą do claude bez zmian, --dry-run pokazuje polecenie.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: "11df08608a1e"
---

`ccam work` uruchamia Claude Code jako konto `work`, a wszystko po nazwie trafia do `claude` bez zmian. Używaj go wszędzie tam, gdzie wpisałbyś `claude`: w terminalu, w skrypcie albo z wejściem przekazanym potokiem.

## Użycie

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

Podaj nazwę konta albo jego numer z `ccam list`. ccam nie odczytuje argumentów po nazwie, więc `ccam work --help` pokazuje pomoc Claude Code. `--` po nazwie jest opcjonalne. Standardowe wejście, kod wyjścia i sygnały działają tak jak w `claude`.

Gdy pierwszy argument zaczyna się od `-` i nie jest żadną z własnych opcji ccam, ccam sam wybiera konto i przekazuje wszystkie argumenty do Claude Code. Używa konta ustawionego dla bieżącego folderu przez [`ccam folder`](/features/folder), a jeśli go nie ma, konta domyślnego: tego ustawionego przez `ccam use`, a jeśli go nie ustawiłeś, konto z najniższym numerem.

### Panel

`ccam` bez argumentów otwiera [panel](/features/dashboard), gdy standardowe wejście i wyjście są terminalem, a w przeciwnym razie wypisuje `ccam list`. Przejdź do konta (strzałkami albo wpisując jego numer) i naciśnij `Enter`, żeby je uruchomić; `n` uruchamia konto z największym zapasem, a `o` proponuje `--continue`, `--resume` lub inne argumenty. Panel zamyka się przed startem Claude Code, więc terminal zostaje najpierw przywrócony.

**Miejsca w planach Team i Enterprise.** Automatyczny wybór konta (`ccam next`, `--fallback`, `ccam pick`, klawisz `n` w panelu oraz, gdy go włączysz, [`auto_switch`](/features/next#auto-switch)) przełącza też między miejscami w tej samej organizacji Claude Team lub Enterprise, co może naruszać warunki Anthropic. Przeczytaj [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) i skonsultuj to ze swoją organizacją i z Anthropic, zanim zaczniesz na nim polegać.

Gdy w powłoce ustawiono `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` lub `ANTHROPIC_AUTH_TOKEN`, Claude Code może użyć tego klucza lub tokena zamiast własnego loginu konta, więc ccam ostrzega o tym na standardowym wyjściu błędów przed startem (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Tylko odczytuje te zmienne; [konto z kluczem API](/features/api) zamiast tego je czyści.

Domyślnie ccam zastępuje się procesem Claude Code, więc nic innego nie działa dalej. Z `relay = true` (albo, na Linux i macOS, z [`auto_switch = true`](/features/next#auto-switch)) na górze `config.toml`, nad pierwszym `[[group]]`, `[[profile]]` lub `[[folder]]`, ccam zamiast tego uruchamia sesje interaktywne jako proces potomny i działa dalej obok nich. Dzięki temu `/ccam:switch` albo `auto_switch` na limicie użycia może samo kontynuować rozmowę na innym koncie, z opcjami sesji, trybem uprawnień, poziomem wysiłku i agentem; zanim `/ccam:switch` przeniesie ją na konto z innej organizacji, ccam cię o to pyta ([szczegóły](/features/next#relay)).

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--dry-run` | Wypisuje środowisko i polecenie, które ccam by uruchomił, i kończy działanie | wyłączone |
| `--fallback` | Dla przebiegów `-p`/`--print`: gdy konto osiągnie limit użycia, uruchamia to samo polecenie ponownie na zalogowanym koncie z największym zapasem ([szczegóły](/features/next#run-again-at-the-limit)) | wyłączone |
| `-v`, `--version` | Wypisuje `ccam <version>` (tylko jako pierwszy argument) | |
| `-h`, `--help` | Pomoc ccam, jako `ccam --help` lub `ccam run --help` | |

Zasady dotyczące argumentów:

- ccam odczytuje własne opcje tylko przed nazwą konta. Po nazwie każdy argument należy do Claude Code, w tym `--dry-run` i `--help`.
- `--dry-run` działa jako `ccam --dry-run work`, `ccam run --dry-run work` i `ccam --dry-run` (dla konta folderu albo domyślnego). `--dry-run` i `--fallback` mogą występować w dowolnej kolejności.
- `ccam run <name|number> [claude args]` to długa forma `ccam <name>`. `ccam run` bez nazwy używa konta folderu albo domyślnego.
- Nazwy kont nie mogą być nazwami poleceń ccam, takimi jak `add`, `list`, `run` czy `usage`, więc `ccam <word>` nigdy nie jest niejednoznaczne.

Próbne uruchomienie wypisuje jeden wiersz, który możesz wkleić do powłoki, a w Windows do PowerShell. Zaczyna się od tego, co uruchomienie zmienia w środowisku: zawsze `CLAUDE_CONFIG_DIR` i `CCAM_PROFILE`, `env -u <name>` dla usuwanej zmiennej (na przykład odziedziczonego `CCAM_RELAY`), a dla [konta z kluczem API](/features/api) jego endpoint, model i to, skąd pochodzi klucz (`ANTHROPIC_API_KEY="$WORK_KEY"`), nigdy sam klucz:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

Wiersz PowerShell przywraca potem zmienne, tak jak `VAR=x cmd` w powłoce, i przekazuje Claude Code jego wiersz poleceń przez `Start-Process`, więc Windows PowerShell 5.1 i PowerShell 7 wykonują go tak samo. Gdy uruchomienie ustawia zmienną na pustą wartość, wiersz `#` pod nim mówi, że PowerShell przed 7.5 zamiast tego usuwa taką zmienną. Dla `claude.cmd` zainstalowanego przez npm wiersz uruchamia polecenie `cmd.exe`, którego użyłoby uruchomienie (zobacz niżej), a argument, który uruchomienie by tam odrzuciło, powoduje też błąd próbnego uruchomienia.

## Jak to działa

Każde uruchomienie przechodzi te same kroki, od wiersza poleceń do `exec`:

<!-- motion:launch -->

Każde pole podświetla się, gdy wykonuje się jego krok, a wiersze pod spodem to środowisko i argumenty, z którymi startuje Claude Code.

1. ccam startuje w `app.Main`, wywoływanym przez `cmd/ccam/main.go`. Jeśli ccam uruchomiono pod inną nazwą, jest to [skrót](/features/shortcuts) i trafia od razu do programu uruchamiającego.
2. ccam analizuje `os.Args`, zanim w ogóle zaangażuje się cobra (`app.parse`). Tylko samo `ccam` (panel), nazwa polecenia ccam (`config.Reserved`), `-h`/`--help` i `-v`/`--version` trafiają do `cli.Execute`. Wszystko inne to uruchomienie.
3. ccam odczytuje `~/.config/ccam/config.toml`, jedyny plik, który parsuje (`launcher.Launch`).
4. ccam znajduje konto po nazwie lub numerze, odrzuca jedno początkowe `--` i odnajduje plik wykonywalny `claude` (`launcher.Prepare`, `ResolveClaude`). Bez nazwy bierze konto folderu, porównując katalog roboczy z wpisami `[[folder]]` w `config.toml` jako ciągi znaków, na macOS i Windows bez względu na wielkość liter (`HereProfile`). Wynik wyszukiwania w PATH jest zapisywany w pamięci podręcznej `~/.local/state/ccam/claude-path` razem z czasem modyfikacji pliku i skrótem `PATH`, więc ccam nie przeszukuje `PATH` przy każdym uruchomieniu. Zmieniony lub brakujący plik albo zmieniony `PATH` wywołuje nowe wyszukiwanie. Ścieżka `claude_bin` w `config.toml` (`~/` jest rozwijane) musi być plikiem wykonywalnym; w przeciwnym razie uruchomienie zatrzymuje się z `claude_bin X is not an executable file. Check claude_bin in config.toml`, zanim cokolwiek zostanie zapisane. ccam odrzuca `claude`, który jest samym ccam (porównanie po i-węźle), bo powstałaby pętla.
5. ccam porównuje folder konta `~/.local/share/ccam/profiles/<name>/` z jego ustawieniami współdzielenia, wywołując tylko `lstat` i `readlink`, jedno lub dwa wywołania systemowe na element (`checkLinks`, `share.Check`). Brakujący link jest tworzony od razu. Wszystko, co mogłoby spowodować utratę danych, zatrzymuje uruchomienie.
6. ccam dodaje współdzielone serwery MCP (`applyMCP`; zobacz niżej).
7. ccam rejestruje sesję (`session.Register`). Zapisuje znacznik nazwany od swojego PID, zawierający czas startu procesu, w `~/.local/state/ccam/locks/profile-<name>/<pid>`, bez względu na `XDG_RUNTIME_DIR`, więc każda powłoka widzi te same znaczniki. Dla konta, które współdzieli historię rozmów, ccam najpierw szuka działającej sesji w `history-<group>-<hash>`, gdzie skrót jest liczony z bieżącego folderu z rozwiązanymi linkami symbolicznymi, i tam też się rejestruje (`historyCheck`).

   `exec` zachowuje PID, więc znacznik wskazuje potem sam proces Claude Code. ccam używa pliku znacznika, a nie trzymanej blokady, bo procesy uruchamiane przez Claude Code odziedziczyłyby blokadę i trzymały ją po zakończeniu sesji. Znacznik, którego nie da się zapisać, nie zatrzymuje uruchomienia: ccam wypisuje `ccam: could not mark the session as running (…); starting Claude Code anyway`.

   Pozostałe sprawdzenia ccam, czy konto działa (zmiana nazwy, usuwanie, zmiany grup, `doctor --fix`, synchronizacja MCP), odczytują też własne zapisy Claude Code w folderze konta (`sessions/<pid>.json`, `daemon.lock`), więc widzą również Claude Code uruchomiony bez ccam: zwykłe `claude`, IDE, sesje w tle i ich demona.
8. ccam ustawia tytuł terminala na `claude · <name>` (gdy standardowe wyjście błędów jest terminalem), zapisuje czas uruchomienia w `~/.local/state/ccam/last-used/<name>` bez `fsync` i wywołuje `syscall.Exec` (`Plan.Exec`).

Claude Code działa z `CLAUDE_CONFIG_DIR` ustawionym na folder konta i `CCAM_PROFILE` ustawionym na nazwę konta. ccam czyści wszelkie odziedziczone wartości tych zmiennych oraz odziedziczone `CCAM_RELAY`. Dla folderu `~/.claude` zarejestrowanego przez `ccam import` na macOS ustawia też `CLAUDE_SECURESTORAGE_CONFIG_DIR` na pustą wartość, żeby konto i zwykłe `claude` używały tego samego loginu z pęku kluczy; w Windows robi to samo, gdy Claude Code przechowuje login w Menedżerze poświadczeń.

### Dlaczego exec

ccam zastępuje własny proces procesem `claude` (`syscall.Exec`). W efekcie:

- ccam nie zostaje w pamięci;
- naciśnięcia klawiszy, sygnały i kod wyjścia należą do Claude Code;
- argumenty pozostają tablicą, nigdy ciągiem dla powłoki.

Windows nie ma `exec`. Tam ccam uruchamia `claude.exe` jako proces potomny z tym samym środowiskiem, argumentami i konsolą, czeka na niego i kończy działanie z jego kodem wyjścia (`Plan.Exec`). `claude.cmd` zainstalowany przez npm działa w `cmd.exe` z pełną ścieżką i każdym argumentem w cudzysłowie, tak żeby `cmd.exe` żadnego z nich nie odczytał jako własnej składni (`&`, `%PATH%`); argumentu z podziałem wiersza nie da się w ten sposób przekazać i jest on odrzucany z radą, żeby zainstalować natywny Claude Code.

Ctrl+C trafia do Claude Code bezpośrednio z konsoli, więc ccam go ignoruje. Znacznik sesji zawiera PID ccam, który żyje dokładnie tak długo jak Claude Code, i jest usuwany, gdy Claude Code się zakończy. Podczas czekania ccam nie robi nic innego.

Prace po sesji, takie jak wykrycie serwerów dodanych przez `claude mcp add`, odbywają się przy następnym uruchomieniu: ccam porównuje czas modyfikacji `.claude.json` z `~/.local/state/ccam/mcp/<name>/seen` i parsuje plik tylko wtedy, gdy się zmienił.

### Jak kierowany jest wiersz poleceń

Trzy wiersze poleceń, trzy drogi przez `internal/app`:

<!-- motion:launch-route -->

Skrót jest rozpoznawany po `argv[0]`, `ccam <name>` po pierwszym argumencie, a obie drogi kończą się na `syscall.Exec`. Tylko polecenie ccam, takie jak `ccam usage`, trafia do cobra.

### Współdzielone serwery MCP

Przy domyślnej strategii `flag` ccam przekazuje serwery MCP grupy do Claude Code jako argument, umieszczony przed twoimi:

<!-- motion:launch-mcp -->

ccam dodaje argument tylko wtedy, gdy `mcp.json` grupy istnieje, i wstawia go zaraz po `claude` (`mcp.FlagArgs`, `Plan.InsertArgs`). Twoje własne `--mcp-config` lub `--strict-mcp-config` nadal obowiązują dodatkowo. To jeden argument, `--mcp-config=<path>`, bo `--mcp-config` w Claude Code przyjmuje kilka wartości: zapisany jako dwa argumenty odczytałby też następujący po nim prompt jako plik konfiguracji.

ccam decyduje, gdzie flaga nie może trafić (`claudeVerb`). Odczytuje podpolecenie Claude Code z argumentów tak, jak robi to Claude Code, pomijając opcje razem z ich wartościami. `remote-control` ani `rc` nie dostają flagi (poniżej), podobnie jak polecenie sesji w tle podane jako pierwsze (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`). `ccam work mcp list`, `get` i `remove` też jej nie dostają, bo Claude Code ją tam odrzuca; zamiast tego ccam informuje na standardowym wyjściu błędów, że serwery grupy ładują się też w sesjach konta:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Każde inne podpolecenie i każda sesja zachowują flagę.

Z `mcp_strategy = "sync"` w `config.toml` ccam zamiast tego zapisuje serwery w `.claude.json` konta i pomija ten krok, jeśli konto już działa. Uruchomienie pomija tę pracę także wtedy, gdy od ostatniego sprawdzenia nie zmienił się ani `mcp.json` grupy, ani `.claude.json`. Jeśli konto ma własny serwer o tej samej nazwie co serwer grupy, synchronizacja zachowuje serwer konta i mówi o tym raz. Synchronizacja przywraca też kopię, która zniknęła z `.claude.json` (`put back group default's MCP servers …`), i zmienia lub usuwa tylko te zapisane przez siebie kopie, których nikt od tego czasu nie zmienił.

### Remote Control

`ccam work rc` (lub `ccam work remote-control`) uruchamia Remote Control z Claude Code jako konto `work`. `claude remote-control` odrzuca `--mcp-config` przed swoim poleceniem i nie przyjmuje go po nim, a jego sesje odczytują serwery MCP tylko z `.claude.json` konta. Dlatego strategia flagi działa tu inaczej:

<!-- motion:launch-rc -->

1. ccam znajduje `remote-control` lub `rc` jako polecenie Claude Code, nawet po opcjach i ich wartościach (`ccam work --model opus rc`), i nie dodaje flagi (`claudeVerb`, `mcp.FlagArgs`).
2. ccam porównuje czas modyfikacji `mcp.json` grupy z `sync-stamp` konta i sprawdza, czy każdy serwer grupy nadal jest w `.claude.json` (`mcp.Synced`). Jeśli `.claude.json` zawiera już bieżące serwery grupy albo grupa nie ma żadnych, Remote Control startuje bez pytania.
3. W przeciwnym razie ccam wymienia serwery, których zabrakłoby w sesjach, i pyta:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. Po odpowiedzi `y` ccam kopiuje je tak jak strategia `sync` (`mcp.Sync`): pod blokadą Claude Code, po wykonaniu kopii zapasowej i z zachowaniem własnych serwerów konta. Po odpowiedzi `N` albo bez terminala Remote Control startuje bez nich, a ccam wskazuje `mcp_strategy = "sync"`, które kopiuje je przed każdym uruchomieniem.

ccam pyta przy każdym uruchomieniu Remote Control, dopóki serwery nie zostaną skopiowane. Po skopiowaniu pyta ponownie tylko wtedy, gdy `mcp.json` grupy się zmieni albo skopiowany serwer zniknie z `.claude.json`. Kopia w `.claude.json` wygrywa z `--mcp-config` przy późniejszych zwykłych uruchomieniach, dlatego `ccam mcp add`, `remove` i `import` aktualizują kopie przechowywane przez inne konta ([Współdzielone serwery MCP](/features/mcp)).

### Budżet 10 ms

Budżet ścieżki uruchamiania to 10 ms w 95. percentylu; w praktyce ccam dodaje około 4 ms przed startem Claude Code. Ścieżka uruchamiania nigdy nie buduje drzewa poleceń cobra, nie uruchamia Bubble Tea ani nie otwiera połączenia sieciowego, a `TestPackageInitIsCheap` kończy się niepowodzeniem, jeśli funkcje `init()` pakietów zajmują łącznie więcej niż 5 ms. `scripts/perf-check.sh` mierzy `ccam --dry-run work`, `ccam work` i skrót za pomocą `hyperfine` i zgłasza błąd powyżej 10 ms, a `BenchmarkLaunchPath` w `internal/launcher` wykonuje kroki w pętli. Aby zobaczyć czas każdego kroku, ustaw `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Bezpieczeństwo i przypadki brzegowe

- Nieznana nazwa lub numer kończy się komunikatem `no account "<name>"` i wskazaniem `ccam list`. Gdy nie ma żadnych kont, ccam prosi o dodanie konta przez `ccam add <name>`.
- Jeśli współdzielona ścieżka stała się zwykłym plikiem lub wskazuje złe miejsce, ccam zatrzymuje się z komunikatem `stopped because of share link problems:`, wymienia problemy i prosi o uruchomienie `ccam doctor --fix`. Na ścieżce uruchamiania nie przenosi danych.
- Jeśli w tym samym folderze działa inna sesja Claude Code, która współdzieli historię rozmów z tym kontem (ta sama grupa), ccam ostrzega, że historia może się wymieszać, i pyta `Run anyway? [y/N]`. Folder jest dopasowywany z rozwiązanymi linkami symbolicznymi. Bez terminala (standardowe wejście z `/dev/null` się nie liczy) albo z zadania w tle pokazuje ostrzeżenie i kontynuuje. Nowy przebieg `-p` zaczyna własną rozmowę i nie dostaje ostrzeżenia; `--version`, `--help` i podpolecenia, które nie otwierają rozmowy, takie jak `mcp`, `auth` czy `plugin`, w ogóle nie są sprawdzane.
- `--dry-run` pomija synchronizację MCP, informację o nowych serwerach, uwagę o `mcp list`, sprawdzenie historii, znacznik sesji, tytuł i czas ostatniego użycia, i nie wykonuje exec. Nadal sprawdza linki współdzielenia i tworzy brakujący link tak jak zwykłe uruchomienie.
- `ccam work rc`, gdy `work` już działa, nie rusza `.claude.json`, informuje o tym i uruchamia Remote Control bez serwerów grupy; uruchom polecenie ponownie po zakończeniu sesji, aby je skopiować. Z `--dry-run` wypisuje komunikat, ale ani nie pyta, ani nie kopiuje.
- Z `mcp_strategy = "sync"` nieudana synchronizacja wypisuje `ccam: MCP sync failed (launching anyway): …`, a Claude Code startuje z serwerami, które już miał.

## Powiązane

- [Polecenia skrótów](/features/shortcuts)
- [Remote Control z ccam](/guides/remote-control)
- [Zarządzanie kontami](/features/accounts)
- [Opis poleceń](/docs/commands)
- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [CLAUDE_CONFIG_DIR wyjaśnione](/guides/claude-config-dir)
