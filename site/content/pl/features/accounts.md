---
title: Zarządzanie kontami
description: Dodawaj, loguj, wyświetlaj, zmieniaj nazwy, usuwaj i importuj konta Claude Code poleceniami ccam add, login, list, rename, remove i import.
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

Każde konto ccam to osobny folder konfiguracji Claude Code z własnym loginem. Tymi poleceniami dodajesz konta, logujesz się do nich, wyświetlasz je, zmieniasz ich nazwy i je usuwasz. Rejestrują też foldery, które już masz.

## Użycie

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

Gdy polecenie przyjmuje `<name|number>`, możesz podać nazwę konta albo jego numer slotu (kolumna `#` w wyniku `ccam list`).

`ccam use` zmienia tylko to, które konto jest domyślne (`default` w `config.toml`).

`ccam add` uruchomione w terminalu bez opcji otwiera kreator. Przechodzi on przez te kroki:

1. nazwa konta, chyba że już ją podałeś;
2. sposób logowania konta: login Claude teraz, klucz API albo login Claude później;
3. grupy do dołączenia, z zaznaczonym `default` (jeśli nie ma jeszcze żadnych grup, pokazuje, co będzie współdzielić grupa `default`);
4. `Create a shortcut command for work? [y/N]`;
5. podsumowanie do potwierdzenia.

`Esc` cofa o krok. Każda z opcji `--group`, `--no-group`, `--shortcut` lub `--no-login`, albo uruchomienie bez terminala (na przykład z potoku), pomija kreator. Nazwa jest wtedy wymagana, a skrót powstaje tylko wtedy, gdy podasz go w `--shortcut`.

Bez `--group` nowe konto dołącza do grupy `default`. Pierwsze konto ją tworzy, a grupa współdzieli serwery MCP, ustawienia, `CLAUDE.md` i rozszerzenia; zobacz [współdzielenie przez grupy](/features/share). W Windows pliki można linkować tylko w Trybie dewelopera. Bez niego pierwsze konto tworzy `default` tylko z `mcp` i `extensions` i o tym informuje (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), a ccam odrzuca `--group` z grupą, która współdzieli ustawienia, pamięć lub historię.

`ccam list` pokazuje slot, nazwę, skrót, adres e-mail, grupy i ostatnie użycie, z `*` przed domyślnym kontem. Adres e-mail to ten, na który folder jest teraz zalogowany (z jego `.claude.json`); gdy różni się od zapisanego przez ccam, zapisany pojawia się za nim jako `(was …)`:

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

Z `--json` każde konto jest obiektem z polami `slot`, `shortcut`, `name`, `email` (bieżącego loginu, jak w tabeli), `default`, `dir`, `last_used` (pomijanym, gdy konto nigdy nie było uruchomione), `groups` i `share`. `share` przypisuje każdy element (`mcp`, `settings`, `memory`, `extensions`, `history`) do grupy, z której pochodzi, albo do `"local"`.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `ccam add --no-login` | Tworzy konto bez logowania; zalogujesz się później przez `ccam login <name>` | wyłączone |
| `ccam add --group <group>` | Dołącza do tej grupy zamiast do `default`. Można powtarzać; `default` powstaje, gdy jeszcze nie istnieje | `default` |
| `ccam add --no-group` | Nie dołącza do żadnej grupy: konto niczego nie współdzieli | wyłączone |
| `ccam add --shortcut <command>` | Tworzy też to polecenie skrótu w `~/.local/bin` | brak |
| `ccam add --api --key-from <source>` | Konto, które zamiast loginu używa klucza API; zobacz [Konta z kluczem API](/features/api) | wyłączone |
| `ccam list --json` | Wypisuje konta jako JSON | wyłączone |
| `ccam remove -y`, `--yes` | Usuwa bez pytania | wyłączone |
| `ccam logout -y`, `--yes` | Wylogowuje bez pytania | wyłączone |
| `ccam import --name <name>` | Nazwa konta dla podanego folderu; odrzucana bez folderu (`--name needs a folder: ccam import <folder> --name ops`) | wyprowadzona z nazwy folderu |
| `ccam import -y`, `--yes` | Importuje każdy znaleziony folder bez pytania | wyłączone |
| `--dry-run` | Obsługiwana przez `ccam list` i `ccam completion`, które niczego nie zmieniają. Odrzucana przez `add`, `login`, `logout`, `use`, `rename`, `remove` i `import` komunikatem `` --dry-run is not supported by `ccam add`; nothing was changed `` (z nazwą polecenia) | wyłączone |

`ccam login` (poza `--expired`), `ccam use` i `ccam rename` nie mają własnych opcji. `ccam completion` przyjmuje `bash`, `zsh` lub `powershell`; wszystko inne kończy się błędem `only bash, zsh and powershell are supported`.

Nazwy kont mają od 1 do 32 znaków spośród `a-z`, `0-9`, `-` i `_` i zaczynają się literą lub cyfrą. Nazwa nie może składać się z samych cyfr, bo `ccam 2` oznacza slot 2. Nie może też być poleceniem ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), więc `ccam <name>` nigdy nie jest niejednoznaczne. Ta lista to `config.Reserved` w `internal/config`.

## Jak to działa

Dodanie konta tworzy jego folder, dołącza je do grup, zapisuje konfigurację, loguje i tworzy skrót, w tej kolejności. Gdy uruchomisz `ccam add work --shortcut ccwork`, ccam wykonuje te kroki (`addProfile` w `internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. ccam sprawdza nazwę (`config.ValidateName`) i skrót (`shim.Check`), zanim cokolwiek utworzy.
2. ccam tworzy `~/.local/share/ccam/profiles/work/` z uprawnieniami `0700` i nadaje kontu najniższy wolny numer slotu, więc numer usuniętego konta jest używany ponownie (`profile.Add`, `Config.NextSlot`). Pierwsze konto staje się domyślnym.
3. Konto dołącza do swoich grup (`default`, chyba że `--group` lub `--no-group` mówi inaczej; pierwsze konto tworzy `default`), a ccam linkuje do nich nowy, pusty folder, na przykład `settings.json` do `groups/default/settings.json` (`applyShare`, przez `share.Build`, a potem `Apply`). Jeśli to się nie uda, nowy folder jest usuwany, bo pozostawiony folder blokowałby każdą kolejną próbę błędem `folder already exists`.
4. Konto jest zapisywane w `~/.config/ccam/config.toml` przez `fsutil.AtomicWrite`: plik tymczasowy, `fsync`, potem `rename`, więc awaria nigdy nie zostawia częściowo zapisanej konfiguracji.
5. ccam uruchamia `claude auth login` z `CLAUDE_CONFIG_DIR` ustawionym na folder konta oraz z `--email` z ostatnim adresem e-mail konta, żeby strona logowania była już wypełniona (`runLogin`). Jeśli logowanie się nie powiedzie, konto zostaje, a ccam wypisuje błąd i każe później uruchomić `ccam login work`.
6. ccam odczytuje `oauthAccount` z `profiles/work/.claude.json` (`profile.ReadIdentity`): adres e-mail, identyfikatory konta i organizacji, nazwę organizacji i plan (`claude_max`, `claude_team` …). Ten plik może mieć kilka megabajtów, więc ccam wyszukuje jeden klucz zamiast parsować cały plik (`TopLevelValue`). ccam przechowuje te dane w `[profile.login]` w `config.toml`.
7. ccam linkuje `~/.local/bin/ccwork` do pliku binarnego ccam, który przy starcie szuka `argv[0]` w `config.toml`, i zapisuje link w `~/.local/state/ccam/shortcuts.json` (`shim.Install`). Jeśli nie da się utworzyć linku, konto i tak zostaje dodane, ccam podaje przyczynę (`Could not create shortcuts: …`), a podsumowanie pokazuje `shortcut: ccwork (not created; see above)`.

Logowanie może trwać minuty, a w tym czasie inny proces ccam może zapisać konfigurację. Dlatego polecenia, które czekają, ponownie wczytują `config.toml`, stosują jedną zmianę (tutaj adres e-mail) i zapisują (`env.update`). Z tego samego powodu `ccam remove` wczytuje konfigurację ponownie po pytaniu.

`ccam login` wykonuje kroki 5 i 6 dla istniejącego konta. Claude Code traktuje adres e-mail tylko jako podpowiedź dla strony logowania, więc możesz się zalogować na inne konto Claude niż wcześniej; ccam porównuje nowy login z `[profile.login]` i o tym mówi (`ccam: this login is other@company.com, but work was dev@company.com before.`). Zachowuje nowy login i podpowiada, jak to cofnąć: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

Login na inne konto Claude usuwa też [zapis limitu](/features/usage#the-limit-record) konta, bo te limity należały do tamtego konta. Gdy inne konto ccam jest zalogowane na to samo konto Claude, ccam ostrzega, że oba dzielą limity użycia.

`ccam login --expired` loguje, jedno po drugim, każde konto, które nigdy nie było zalogowane albo którego login wygasł: wygasł jego token odświeżania albo Claude Code wyczyścił login po tym, jak serwer go odrzucił. Zanim otworzy się pierwsza przeglądarka, sprawdza, czy `config.toml` da się zapisać.

`ccam logout` uruchamia `claude auth logout` na koncie i czyści zapisany adres e-mail; zapis logowania zostaje, więc następny login jest z nim porównywany. Claude Code unieważnia login i usuwa to, co przechowuje dla konta, w tym logowania do serwerów MCP. Dlatego ccam najpierw ostrzega i pyta (`--yes` pomija pytanie), gdy:

- login zawiera logowania MCP;
- konto jest uruchomione (jego sesje też zostają wylogowane);
- folder innego konta ma ten sam login (skopiowany folder), który wylogowanie również unieważnia.

Na przykład:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

Aby zalogować się na inne konto Claude, wystarczy samo `ccam login <name>`, a logowania MCP zostają.

`ccam list` dodaje kolumnę `Organization`, gdy konta należą do więcej niż jednej organizacji (plany prywatne pokazują `personal`), a `ccam list --json` ma `plan` i `org` dla każdego konta.

### Logowanie na zdalnej maszynie

Przez SSH część `claude auth login` działająca w przeglądarce może zawieść. Claude Code czeka na zalogowanie na serwerze zwrotnym na zdalnej maszynie (`http://localhost:<port>/callback`), ale przeglądarka, która kończy logowanie, działa na twoim komputerze, gdzie ten adres prowadzi donikąd, więc przeglądarka pokazuje błąd taki jak `ERR_CONNECTION_REFUSED`. Adres nadal zawiera jednorazowy kod, więc ccam może go dostarczyć:

<!-- motion:accounts-relay -->

ccam nie oddaje twojego terminala Claude Code (`runLogin`). Podłącza standardowe wejście Claude Code do potoku i kopiuje do niego każdy wpisany przez ciebie wiersz (`relayLogin`), więc kod wklejony przy monicie `Paste code here` w Claude Code nadal dociera. Gdy wiersz jest adresem zwrotnym, ccam akceptuje go tylko wtedy, gdy host to `localhost`, `127.0.0.1` lub `[::1]`, podano port, ścieżka to `/callback` i obecne są zarówno `code`, jak i `state`; `http://` można pominąć (`callbackURL`). ccam wysyła wtedy żądanie pod ten adres na serwerze z limitem czasu 10 sekund i bez podążania za przekierowaniami, tak jak zrobiłaby to przeglądarka (`deliverCallback`). Żądanie trafia do serwera zwrotnego oczekującego `claude auth login` (linia przerywana na diagramie). Claude Code otrzymuje kod i zapisuje login, na Linux w `.credentials.json` w folderze konta. Wklejony wiersz nigdy nie może sprawić, że ccam połączy się z innym hostem.

Przed rozpoczęciem logowania ccam wypisuje:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

Po powodzeniu wypisuje `→ Passed the sign-in to Claude Code on this machine.`

Gdy logowanie czeka na dane w terminalu, ccam wyłącza echo i tryb liniowy terminala i dla każdego wpisanego znaku pokazuje `*`, więc wklejony kod lub adres zwrotny nie pojawia się na ekranie. ccam sam obsługuje edycję wiersza (`readLine`): `Enter` wysyła wiersz, `Backspace` usuwa jeden znak, a `Ctrl+D` w pustym wierszu kończy wprowadzanie. Terminal jest przywracany po zakończeniu logowania, a także po `Ctrl+C`, `SIGTERM` lub `SIGHUP`, które nadal kończą ccam. Gdy standardowe wejście i standardowe wyjście błędów nie są terminalem, wiersze są kopiowane na bieżąco.

### Zmiana nazwy, usuwanie i import

`ccam rename` i `ccam remove` działają na folderze konta i na tym, co ccam przechowuje pod nazwą konta w innych miejscach. ccam wymienia te ścieżki w `nameState`:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- pamięci podręczne użycia `~/.local/state/ccam/usage/<name>.json` i `<name>.statusline.json`;
- zapis limitu `<name>.limit.json`;
- notatki sesji w `~/.local/state/ccam/handoff/<name>/`.

`ccam remove work` obsługuje je w tej kolejności:

<!-- motion:accounts-remove -->

1. `ccam remove` najpierw odmawia, jeśli konto jest uruchomione, potem pyta i po twojej odpowiedzi sprawdza ponownie (wychwytuje to sesję uruchomioną w czasie, gdy pytanie było otwarte).
2. Zanim cokolwiek zmieni, sprawdza, czy `config.toml` da się zapisać.
3. Jeśli konto współdzieli rozszerzenia, ccam przepisuje `installed_plugins.json` i `known_marketplaces.json` grupy (`RewritePluginPaths`), bo Claude Code przechowuje bezwzględne ścieżki wtyczek, które mogą prowadzić przez folder tego konta.
4. Następnie zarządzane konto (takie, którego folder utworzył ccam, w odróżnieniu od konta zarejestrowanego przez `ccam import`) jest wylogowywane: `claude auth logout` działa w jego folderze, więc w koszu nie czeka działający login (ani, na macOS, element pęku kluczy nazwany od folderu, który odziedziczyłoby późniejsze konto o tej nazwie). Ten krok jest pomijany dla konta z kluczem API i wtedy, gdy folder innego konta ma ten sam login, który wylogowanie również by unieważniło. Nieudane wylogowanie jest zgłaszane, a usuwanie trwa dalej.
5. ccam przenosi folder do `~/.local/share/ccam/trash/<name>-<time>/` (`profile.Trash`), konto znika z `config.toml`, ścieżki `nameState` są usuwane, żeby późniejsze konto o tej samej nazwie ich nie odziedziczyło, a link skrótu jest usuwany: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

Jeśli mimo to `config.toml` nie da się zapisać, folder wraca z kosza. Aby odzyskać usunięte konto, przenieś folder z kosza z powrotem i zaloguj się ponownie.

`ccam rename` zmienia nazwę folderu i przepisuje ścieżki wtyczek w jego `plugins/` ze starego folderu na nowy. Przenosi ścieżki `nameState` pod nową nazwę i zachowuje status domyślnego konta, jeśli konto go miało. Aktualizuje też notatki sesji, więc `ccam next --continue` nadal je znajduje. Polecenie skrótu zachowuje swoją nazwę. Jeśli przepisanie ścieżek wtyczek się nie powiedzie albo nie da się zapisać `config.toml`, folder wraca do starej nazwy.

Na macOS, a także w Windows, gdy Claude Code przechowuje loginy w Menedżerze poświadczeń, sam folder zachowuje starą nazwę, bo login i logowania MCP są przechowywane pod ścieżką folderu, a przeniesiony folder zostałby wylogowany. Konfiguracja zapisuje wtedy folder jako `folder` obok nowej nazwy, a ccam informuje `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

`ccam import` rejestruje folder tam, gdzie jest, jako `dir` w `config.toml`. Nigdy go nie przenosi ani nie kopiuje: skopiowany token odświeżania sprawiłby, że kopie nawzajem by się wylogowywały.

Bez argumentu z folderem szuka niezarejestrowanych folderów `~/.claude` i `~/.claude-*` i pyta o każdy. Nazwy pochodzą od folderu (`~/.claude` to `main`, `~/.claude-work` to `work`, `~/.claude-2` to `acct2`). Jeśli ccam nie może wyprowadzić nazwy albo nazwa jest zajęta, wypisuje folder razem z poleceniem `--name` do uruchomienia. Znaleziony folder, którego nie udało się zaimportować, jest zgłaszany jako `Skipped <path>: <reason>`, a pozostałe są importowane.

ccam najpierw zadaje wszystkie pytania, potem ponownie wczytuje `config.toml` i rejestruje foldery, więc zmiany zapisane w międzyczasie przez inne polecenia ccam zostają zachowane. Wiersze `✓` pojawiają się po zapisie.

Bez `CLAUDE_CONFIG_DIR` Claude Code trzyma `.claude.json` w folderze domowym obok `~/.claude`, więc gdy importujesz `~/.claude`, a nie ma on własnego `.claude.json`, ccam tworzy link `~/.claude/.claude.json` do `~/.claude.json`. Zwykłe `claude` i `ccam main` widzą wtedy to samo konto.

Jeśli `~/.claude` ma własny `.claude.json` obok `~/.claude.json`, ccam zostawia oba i o tym informuje. `ccam main` używałby wewnętrznego, a zwykłe `claude` zewnętrznego, co rozdziela ich dane logowania i serwery MCP. Żeby został jeden, zastąp wewnętrzny plik linkiem ([doctor](/features/doctor) też to zgłasza).

Gdy folder ma ten sam login co zarejestrowane konto (skopiowany folder), ccam o tym mówi: to, które pierwsze go odnowi, wyloguje drugie, więc przed użyciem uruchom `ccam login <name>`.

Zaimportowane konta nie należą do żadnej grupy, dopóki ich do niej nie dodasz (`ccam group join <group> <name>`), i nie powstaje dla nich skrót.

`ccam completion bash|zsh` wypisuje skrypt uzupełniania, który uzupełnia też nazwy kont.

## Bezpieczeństwo i przypadki brzegowe

- `ccam rename` i `ccam remove` odmawiają działania, gdy konto działa (`refuseRunning`), z komunikatem ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``, bo sesja dalej zapisuje do folderu. Sprawdzenie obejmuje Claude Code niezależnie od sposobu uruchomienia: przez ccam, zwykłe `claude`, IDE, sesje w tle i ich demona.
- `ccam rename`, `ccam remove` i `ccam logout` odmawiają, gdy inny proces ccam odnawia lub zmienia ten login: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- Gdy w innym procesie ccam trwa zmiana grup, `ccam rename` i `ccam remove` odmawiają z komunikatem `another ccam process is changing the groups right now; try again when it is done`; po zmianie, która zatrzymała się w połowie, z komunikatem ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.
- `ccam remove` na zaimportowanym koncie zachowuje folder i jego login dla zwykłego `claude`. Jeśli konto jest w grupach, najpierw je opuszcza tak jak `ccam group leave`: każdy współdzielony element staje się w folderze kopią wersji grupy, a serwery MCP grupy trafiają do jego `.claude.json`. Potem jest wyrejestrowywane: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. Kopie zapasowe jego `.claude.json` utworzone przez ccam są przenoszone do `<folder>/.ccam-backup/mcp-<time>`; jeśli nie da się ich przenieść, zostają w `~/.local/state/ccam/mcp/<name>/backups`, a komunikat o tym mówi. Twój `~/.claude` nie trafia do kosza.
- Dwa konta nie mogą mieć tego samego folderu: `config.toml`, który podaje jeden folder dwa razy, nie daje się wczytać.
- `ccam add` odrzuca nazwę, której folder już istnieje w `profiles/`, i wskazuje `ccam import`. Na macOS zatrzymuje się też, gdy pęk kluczy nadal przechowuje login lub logowania MCP dla ścieżki nowego folderu, pozostawione przez starszą wersję ccam, i podaje polecenie `security delete-generic-password`, które je usuwa. Jeśli nie da się zapisać `config.toml`, nowy folder jest ponownie usuwany.
- `ccam import` odrzuca folder, który jest już zarejestrowany, ścieżkę, która nie jest folderem, i zajętą nazwę. ccam ostrzega też, gdy folder jest kopią `~/.claude` (ten sam login), nawet jeśli sam `~/.claude` nie jest zarejestrowany.
- `ccam logout` najpierw pyta, gdy konto jest uruchomione; `ccam login` tego nie sprawdza.
- Jeśli przekazanie nie może dotrzeć do Claude Code, na przykład po zakończeniu logowania, ccam wypisuje `ccam: could not pass the sign-in to Claude Code:` z przyczyną. Błąd HTTP zwykle oznacza, że kod wygasł; zacznij logowanie od nowa.

## Powiązane

- [Uruchamianie Claude Code](/features/launch)
- [Polecenia skrótów](/features/shortcuts)
- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Kilka kont Claude Code na jednym komputerze](/guides/multiple-claude-code-accounts)
