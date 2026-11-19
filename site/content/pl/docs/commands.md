---
title: Opis poleceń
metaTitle: "Opis poleceń ccam"
description: Wszystkie polecenia i opcje ccam, od ccam add i ccam <name> po ccam usage, ccam next i ccam doctor.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: "28e1ad1dcbad"
---

Wpisz `--dry-run` przed poleceniem, żeby zobaczyć, co by zrobiło, bez wykonywania go, na przykład `ccam --dry-run work`. Działa to przy uruchamianiu Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) oraz z `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` i `doctor --fix`. Polecenia, które tylko odczytują, takie jak `list`, `usage`, `share` i `ext`, przyjmują tę opcję i i tak niczego nie zmieniają. Pozostałe odrzucają ją zamiast się wykonać, w tym `group rename`.

<!-- commands -->

## Uruchamianie Claude Code

`ccam <name|number> [arguments]` uruchamia Claude Code jako to konto. ccam nie odczytuje argumentów po nazwie; przekazuje je do `claude` w tej samej kolejności, a standardowe wejście, kod wyjścia i sygnały działają dokładnie tak jak w samym `claude`. `--` po nazwie jest opcjonalne.

`ccam` z argumentami, ale bez nazwy konta, na przykład `ccam -p "summarise"`, używa konta ustawionego dla bieżącego folderu, a jeśli go nie ma, konta domyślnego. Ustawisz je poleceniami `ccam folder <name> [folder]` i `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` przed nazwą konta jest przeznaczone dla przebiegów `-p`/`--print`. Gdy konto osiągnie limit użycia, ccam uruchamia to samo polecenie ponownie na zalogowanym koncie z największym zapasem i wypisuje tylko wyjście przebiegu, który się zakończył. Zobacz [Konta według folderu](/features/folder) i [Ponowienie po osiągnięciu limitu](/features/next#run-again-at-the-limit).

**Miejsca w planach Team i Enterprise.** Automatyczny wybór konta (`ccam next`, `--fallback`, `ccam pick`, klawisz `n` w panelu oraz, gdy go włączysz, [`auto_switch`](/features/next#auto-switch)) przełącza też między miejscami w tej samej organizacji Claude Team lub Enterprise, co może naruszać warunki Anthropic. Przeczytaj [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) i skonsultuj to ze swoją organizacją i z Anthropic, zanim zaczniesz na nim polegać.

Skrót taki jak `ccwork` uruchamia swoje konto w ten sam sposób: `ccwork --resume` robi to samo co `ccam work --resume`. Zobacz [Skróty](/docs/commands#shortcuts) poniżej.

## Dodawanie i usuwanie kont

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

`ccam import` bez podanego folderu szuka `~/.claude` i `~/.claude-*` i pyta o każdy z nich; folder, którego nie udało się zarejestrować, jest zgłaszany jako `Skipped <path>: <reason>`, a pozostałe są rejestrowane. Jeśli nazwa folderu jest już zajęta, zaimportuj ten folder osobno poleceniem `ccam import <folder> --name <name>`. Zaimportowane foldery są rejestrowane tam, gdzie są; ccam nigdy nie przenosi ani nie kopiuje loginu.

`ccam login --expired` loguje też konta, których login wyczyścił Claude Code. `ccam remove` wylogowuje zarządzane konto (takie, którego folder utworzył ccam, w odróżnieniu od konta zarejestrowanego przez `ccam import`), zanim jego folder trafi do kosza, chyba że inne konto ma ten sam login. Usunięcie zaimportowanego konta wyłącza je z jego grup i zostawia w jego folderze kopię każdego współdzielonego elementu. Na macOS `ccam rename` zachowuje nazwę folderu, bo pęk kluczy przechowuje login pod ścieżką folderu.

`ccam add --api` dodaje konto, które zamiast loginu Claude używa klucza API; ccam odczytuje klucz przy każdym uruchomieniu i nigdy go nie zapisuje. Zobacz [Konta z kluczem API](/features/api).

`ccam rename`, `ccam remove` i polecenia `ccam group` odmawiają zmiany konta, gdy Claude Code działa jako to konto, niezależnie od tego, jak go uruchomiono (w tym sesje w tle). Najpierw zakończ sesję. Odmawiają też, gdy zmiana grup nie została dokończona: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` wypisuje każde konto z jego numerem (slotem), skrótem, adresem e-mail i folderem. Jego tablica `groups` wymienia grupy konta, a obiekt `share` przypisuje każdy element (`mcp`, `settings`, `memory`, `extensions`, `history`) do grupy, z której pochodzi, albo `local`.

## Skróty

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Skrót to polecenie o nazwie wybranej przez ciebie, które uruchamia Claude Code jako jedno konto. Podobnie jak `ccam work` przekazuje każdy argument do Claude Code bez zmian, więc `ccwork --resume` i `ccwork -p "hi"` działają tak, jak się spodziewasz. Jest to link do ccam w `~/.local/bin`.

Skróty są opcjonalne. Kreator `ccam add` pyta, czy utworzyć skrót, `ccam add <name> --shortcut <command>` tworzy go bez pytania, a `ccam list` pokazuje skrót każdego konta albo `-`, gdy konto go nie ma.

Nazwy mogą zawierać znaki a-z, 0-9, `.`, `-` i `_`. Nazwa nie może być poleceniem ccam, skrótem innego konta ani nazwą innego programu w twoim `PATH`, ani plikiem w `~/.local/bin`, którego nie utworzył ccam, na przykład twoją własną kopią lub linkiem do ccam. Skrót zmieniony ręcznie w `config.toml`, który łamie te zasady, jest `claude` lub `ccam` albo jest przypisany dwóm kontom, sprawia, że `config.toml` nie daje się wczytać, a komunikat wskazuje konto.

`ccam rename` zachowuje nazwę skrótu, a `ccam remove` usuwa skrót razem z kontem. `ccam shims install` odtwarza linki wszystkich kont, na przykład po przeniesieniu ccam; plik blokujący jeden skrót nie zatrzymuje pozostałych.

## Współdzielenie

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Grupa współdzieli elementy między swoimi kontami; konto może należeć do kilku grup, ale każdy element pochodzi tylko z jednej grupy. Zobacz [współdzielenie przez grupy](/features/share).

## Kopiowanie istniejącej konfiguracji

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` kopiuje `settings.json`, `CLAUDE.md`, rozszerzenia i serwery MCP użytkownika z folderu Claude Code do grupy, żeby dostało je każde konto współdzielące te elementy. Folder jest tylko odczytywany.

- `--settings`, `--memory`, `--extensions`, `--mcp` i `--history` kopiują tylko wskazane elementy. Bez nich w terminalu pojawia się lista wyboru, na której zaznaczone jest wszystko oprócz historii rozmów; poza terminalem kopiowane jest wszystko oprócz historii rozmów.
- `--mcp` działa też wtedy, gdy `~/.claude` jest pusty lub go nie ma, ale `~/.claude.json` zawiera serwery.
- `-g` lub `--group` wybiera grupę.
- `-y` lub `--yes` kopiuje bez pytania i zachowuje kopię grupy, gdy plik się różni.

Logowania MCP nie są kopiowane, bo skopiowane logowanie wylogowuje pozostałych posiadaczy; plan wymienia konta, które logują się przez `/mcp`. Zobacz [kopiowanie istniejącej konfiguracji](/docs/sharing#copy-an-existing-setup).

## Serwery MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` przenosi do grupy konta serwery, które ma tylko to konto, także te, które ma tylko dla jednego folderu (zwykłe `claude mcp add`); ładują się one potem w każdym folderze. Serwer, którego nazwy grupa używa już dla innej definicji, jest pomijany; dodaj `--force`, aby zastąpić kopię w grupie, a ccam najpierw zrobi jej kopię zapasową.

## Użycie i przełączanie

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

- `next --auto`, `--fallback` i `ccam pick` wybierają w ten sam sposób: najpierw konta poniżej `auto_max_usage` (`config.toml`, domyślnie 100), według `auto_strategy` (`headroom`, `expiring` lub `balance`), z uwzględnieniem tygodniowego okna `--model`. Nigdy nie wybierają konta na poziomie 100% (okno, które Claude Code zgłosił jako pełne, liczy się jako 100% aż do resetu) ani loginu kluczem API i pozostają w obrębie organizacji konta, od którego zaczynają (`--cross-org` lub `cross_org = true`, żeby szukać dalej).
- `--wait` przyjmuje wartość po `=` (`--wait=2h`); zapisane ze spacją jest odrzucane z podpowiedzią. Gdy nie da się wybrać żadnego, kończą działanie z kodem 75, a z `--wait` czekają na reset pierwszego konta.
- `ccam next --continue` przenosi rozmowę nawet między kontami, które nie współdzielą historii: Claude Code kopiuje ją wtedy na nowe konto. Z `--brief` nowe konto zaczyna zamiast tego nową rozmowę, od podsumowania poprzedniej, które ccam pisze bez użycia modelu ([`--brief`](/features/next#brief)).
- W przeciwieństwie do `ccam <name>`, `ccam next` odczytuje własne opcje, więc opcje Claude Code podaje się po `--`; bez tego ccam zatrzymuje się z `unknown flag`.
- `ccam agent` uruchamia zadanie przez `claude -p` w tym folderze jako konto inne niż to, z którego go wywołano, w trybie uprawnień i z poziomem wysiłku sesji wywołującej, i wypisuje ostatnią odpowiedź Claude. Gdy to konto osiągnie limit, rozmowa toczy się dalej jako inne ([Subagenci na innych kontach](/features/agent)).
- W Claude Code polecenie pluginu `/ccam:switch [account]` przekazuje rozmowę dalej, a jego hooki uruchamiają `ccam hook`.
- Z `auto_switch = true` w `config.toml` (Linux i macOS; domyślnie wyłączone) hook limitu dodatkowo kończy sesję interaktywną uruchomioną przez ccam i kontynuuje ją na innym koncie z zapasem, ale tylko na koncie tego samego rodzaju: inne miejsce w tej samej organizacji albo inny plan prywatny ([`auto_switch`](/features/next#auto-switch)).

Zobacz [Przełączanie na następne konto](/features/next), [Wybór konta dla skryptów](/features/pick), [Subagenci na innych kontach](/features/agent) i [Sesje w tle wszystkich kont](/features/agents).

## Konserwacja

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

`ccam update` sprawdza sumę kontrolną SHA-256 wydania, zanim podmieni ccam, a `--version v0.2.0` instaluje konkretne wydanie. Kopię ccam zainstalowaną z pakietu `.deb` aktualizujesz przez apt. Szczegóły znajdziesz w części [aktualizacja](/docs/install#update).

`ccam uninstall` wypisuje, co usunie, i najpierw pyta o zgodę: każde konto tak jak `ccam remove`, skróty, foldery ccam i polecenie `ccam`. Zaimportowany folder zostaje tam, gdzie jest. Szczegóły znajdziesz w części [odinstalowanie](/docs/install#uninstall).

## Pliki

| Ścieżka | Zawartość |
| --- | --- |
| `~/.config/ccam/config.toml` | Konta, grupy i ustawienia (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). Ustawienia umieszcza się na górze, nad pierwszym `[[group]]`, `[[profile]]` lub `[[folder]]`: ustawienie pod tabelą należy do tej tabeli, więc ccam ostrzega o nim przy każdym uruchomieniu i w `ccam doctor`, a polecenia zapisujące konfigurację odmawiają, dopóki nie zostanie przeniesione. Komentarze nie są zachowywane, gdy ccam zapisuje plik |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` każdego konta |
| `~/.local/share/ccam/groups/<group>/` | Współdzielone ustawienia, serwery MCP, rozszerzenia i historia |
| `~/.local/state/ccam/` | Pamięć podręczna użycia, zapisy limitów, blokady, notatki i inny stan (zobacz listę pod tabelą) |
| `~/.local/bin/<shortcut>` | Skróty takie jak `ccwork`, czyli linki do ccam |

Folder stanu zawiera:

- pamięć podręczną użycia i zapisy limitów (`usage/<name>.json`, `usage/<name>.limit.json`);
- czasy ostatniego użycia;
- znaczniki sesji i blokady (`locks/`);
- skróty utworzone przez ccam (`shortcuts.json`);
- niedokończoną zmianę grup (`group-change.json`);
- notatki zostawiane przez `/ccam:switch` i hook limitu (`handoff/<account>/<session>.json`);
- informację, kiedy `ccam pick` lub `auto_switch` ostatnio wybrał każde konto (`picks/<name>`; wybór liczy się przez 2 minuty);
- podsumowania `--brief`, przechowywane przez 7 dni (`brief/<session>.md`).

Ustawienie `XDG_CONFIG_HOME`, `XDG_DATA_HOME` lub `XDG_STATE_HOME` zmienia położenie tych folderów. W Windows `~` to `%USERPROFILE%`, a skróty to pliki `<shortcut>.exe`.
