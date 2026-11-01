---
title: Przełączanie na następne konto
description: ccam next otwiera Claude Code jako inne konto; --auto wybiera konto z największym zapasem, --continue przenosi rozmowę, --fallback ponawia zadania -p.
updated: 2026-11-01
order: 8
translation: machine
sourceHash: "79300c521748"
---

`ccam next` otwiera Claude Code jako inne konto bez podawania jego nazwy. Używaj go, gdy konto osiągnie limit użycia. Są cztery sposoby przełączenia:

- Samo `ccam next` wybiera konto następne po ostatnio używanym, a `--auto` wybiera konto z największym zapasem. `--continue` kontynuuje ostatnią rozmowę z tego folderu na nowym koncie, niezależnie od tego, czy oba konta współdzielą historię.
- W Claude Code [plugin ccam](#inside-claude-code) informuje, gdy limit zatrzyma sesję, a `/ccam:switch` przekazuje rozmowę dalej.
- W skryptach, które uruchamiają `claude -p`, `ccam --fallback` przełącza konto samo (zobacz [Ponowienie po osiągnięciu limitu](#run-again-at-the-limit)).
- Z włączonym [`auto_switch`](#auto-switch) (Linux i macOS) ccam przenosi też sesję interaktywną, którą Claude Code zatrzyma na jej limicie.

**Miejsca w planach Team i Enterprise.** Automatyczny wybór konta (`ccam next`, `--fallback`, `ccam pick`, klawisz `n` w panelu oraz, gdy go włączysz, [`auto_switch`](/features/next#auto-switch)) przełącza też między miejscami w tej samej organizacji Claude Team lub Enterprise, co może naruszać warunki Anthropic. Przeczytaj [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) i skonsultuj to ze swoją organizacją i z Anthropic, zanim zaczniesz na nim polegać.

## Użycie

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # continue this folder's last conversation on the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next -c --brief            # a new conversation that starts from a summary of the last one
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the most weekly quota left per day
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work CCAM_PROFILE=work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam wypisuje przełączenie jako `→ <from> → <to>` na standardowe wyjście błędów, przy wyborze `--auto` z powodem, a potem uruchamia Claude Code jako nowe konto tak samo jak `ccam <name>`. `<from>` to konto, którego używałeś ostatnio, a z `--continue` konto, do którego należy rozmowa. Argumenty po `--` trafiają do `claude` bez zmian.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--auto` | Wybiera zalogowane konto według użycia (zobacz niżej), pomijając konta na poziomie 100% | wyłączone |
| `--max-usage <1-100>` | Z `--auto`: preferuje konta poniżej tego użycia | `auto_max_usage`, w przeciwnym razie 100 |
| `--strategy headroom\|expiring\|balance` | Z `--auto`: najniższe użycie 5-godzinne, najbliższy reset okna 5-godzinnego albo największy tygodniowy limit pozostały na dzień | `auto_strategy`, w przeciwnym razie `headroom` |
| `--cross-org` | Z `--auto`: wybiera też konta innej organizacji (zobacz niżej) | `cross_org`, w przeciwnym razie wyłączone |
| `--wait[=<duration>]` | Z `--auto`: gdy żadne konto nie ma zapasu, czeka na reset pierwszego, najwyżej tyle. Wartość podaje się po `=`: `--wait 2h` jest odrzucane z podpowiedzią `--wait takes its value after '=': --wait=2h` | wyłączone; `5h`, gdy podano samą opcję |
| `-c`, `--continue` | Kontynuuje ostatnią rozmowę z tego folderu na nowym koncie (zobacz niżej) | wyłączone |
| `--brief` | Z `--continue`: zaczyna na nowym koncie nową rozmowę od podsumowania poprzedniej napisanego przez ccam, zamiast wznawiać ją całą ([zobacz niżej](#brief)) | wyłączone |
| `-- <claude args>` | Argumenty przekazywane do `claude` po argumentach ccam | brak |
| `--dry-run` | Obsługiwana: wypisuje `→ <from> → <to>` i wiersz polecenia, który zostałby uruchomiony, bez startu Claude Code i bez zapisywania uruchomienia | wyłączone |

`--auto` i `--continue` działają razem. Z `--continue` argumenty `--resume` są umieszczane przed twoimi.

Własne opcje Claude Code, takie jak `--dangerously-skip-permissions` czy `--model`, podaje się po `--`. W przeciwieństwie do `ccam <name>`, które przekazuje dalej wszystko po nazwie, `ccam next` samo odczytuje swoje argumenty, więc bez `--` zatrzymuje się z `ccam: unknown flag: --dangerously-skip-permissions`.

## Jak to działa

Każde uruchomienie przez ccam zapisuje czas w `~/.local/state/ccam/last-used/<name>`. `ccam next` odczytuje te pliki, sprawdza, które konta są zalogowane, i wybiera.

1. ccam tworzy jeden `next.Candidate` na konto z `config.toml` (`cli.candidates`): konto, czas jego ostatniego uruchomienia (`launcher.LastUsed`) i informację, czy ma działający login (`creds.Has`: login Claude w `.credentials.json` albo, na macOS, w pęku kluczy, który nie wygasł całkowicie ani nie został wyczyszczony przez Claude Code). Folder, który zamiast tego używa klucza API (login Console albo konto dodane przez `ccam add --api`), nie ma limitów planu i nigdy nie jest wybierany.
2. Bez `--auto` ccam sortuje konta według numeru slotu i znajduje to uruchomione najpóźniej (`next.Rotate`). To konto, z którego się przełączasz. Jeśli żadne konto nie było jeszcze uruchamiane, jest to najniższy slot.
3. Od niego idzie dalej w kolejności slotów, wracając od ostatniego slotu do pierwszego, i zwraca pierwsze zalogowane konto.
4. ccam przygotowuje i uruchamia wybrane konto dokładnie tak jak `ccam <name>` (`launcher.LaunchWith`): sprawdza linki konta, stosuje współdzielenie MCP, rejestruje sesję i zastępuje ccam przez `claude` za pomocą `syscall.Exec`.

Przejście zaczyna się od ostatnio używanego konta i zatrzymuje na pierwszym zalogowanym koncie po nim:

<!-- motion:next-rotate -->

`work` uruchomiono najpóźniej, więc przejście zaczyna się po nim. `side` nie ma loginu i zostaje pominięte, a wybrane zostaje `lab`. Gdyby `lab` też było wylogowane, przejście wróciłoby na początek do `personal` (`% len` na diagramie). Rotacja nigdy nie zwraca konta, od którego zaczęła: gdy nie ma innego zalogowanego konta, kończy się błędem.

### --auto

`--auto` decyduje według użycia zamiast kolejności slotów. ccam zostawia zalogowane konta, pyta o użycie wszystkich naraz tak samo jak `ccam usage` (`Fetcher.FetchAll`) i wybiera według poniższych zasad (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

Konto jest oceniane według okna 5-godzinnego i tygodniowego, a gdy wśród argumentów dla Claude Code (po `--`) jest `--model`, także według tygodniowego okna tego modelu (`opus`, `sonnet`, `haiku` lub `fable` w nazwie). ccam dzieli konta na trzy poziomy (`Policy.Choose()`):

1. **Z zapasem**: każde okno poniżej `auto_max_usage` (`config.toml`, domyślnie 100). Spośród nich wybiera strategia: `headroom` (domyślna) bierze najniższe użycie 5-godzinne; `expiring` bierze okno 5-godzinne, które zresetuje się najwcześniej, żeby pozostały limit został zużyty, zanim przepadnie; `balance` bierze konto z największym tygodniowym limitem pozostałym na każdy dzień do resetu tygodnia, licząc przeciw niemu każdą uruchomioną sesję ccam, żeby tygodniowy limit nie przepadał niewykorzystany, a równoległe sesje się rozkładały.
2. **Powyżej progu**: poniżej 100%, ale na poziomie `auto_max_usage` lub wyżej. Używane tylko wtedy, gdy żadne konto nie ma zapasu; jako pierwsze idzie to z najniższym użyciem 5-godzinnym.
3. **Pełne**: okno na poziomie 100% albo nieznane użycie 5-godzinne. Nigdy niewybierane.

Okno liczy się jako 100% aż do resetu, gdy Claude Code odrzucił w nim żądanie: hook [pluginu ccam](#inside-claude-code) to zapisuje (`~/.local/state/ccam/usage/<name>.limit.json`), bez względu na zapisane liczby. Dzięki temu konto, które właśnie osiągnęło limit, nie zostanie wybrane ponownie przed resetem, nawet gdy jego liczby są stare.

Ustawienia dla wszystkich kont umieszcza się na górze `config.toml`, nad pierwszym `[[group]]`, `[[profile]]` lub `[[folder]]`. Ustawienie pod jedną z tych tabel należy do niej: ccam ostrzega o nim przy każdym uruchomieniu, a polecenia zapisujące konfigurację odmawiają, dopóki nie zostanie przeniesione.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organizacje.** Automatyczne wybory pozostają wśród kont tego samego rodzaju co to, które opuszczasz: miejsc w tej samej organizacji Team lub Enterprise albo planów prywatnych (Pro, Max), więc praca z konta firmowego nie jest kontynuowana na prywatnym ani odwrotnie. ccam zna organizację każdego konta z jego loginu (`ccam login` ją zapisuje; starsze loginy są odczytywane z `.claude.json`); konto, którego organizacji ccam nie potrafi ustalić, może zostać dopasowane do dowolnego konta. Dwa foldery zalogowane na to samo konto Claude dzielą jego limity, więc jeden nigdy nie przejmuje pracy po drugim. `--cross-org` albo `cross_org = true` w `config.toml` znosi zasadę organizacji, ale po [`/ccam:switch`](#inside-claude-code) przeniesienie do innej organizacji nadal czeka na twoje „tak”.

Gdy nie da się nic wybrać, ccam kończy działanie z kodem 75 (`EX_TEMPFAIL`), więc skrypt może odróżnić „spróbuj później” od błędu. Z `--wait` zamiast tego czeka, aż zresetuje się okno blokujące pierwsze pełne konto (30 sekund po czasie resetu, gdy endpoint pokazuje już nowe okno), wypisuje `ccam: no account has room; waiting until 15:10` i wybiera ponownie; gdy żadne konto nie podaje czasu resetu, sprawdza ponownie po 10 minutach. Jeśli nic się nie zwolni przed końcem oczekiwania, kończy działanie z kodem 75.

Zapytanie działa według zasad `ccam usage`: liczby pobrane mniej niż 60 sekund temu są używane ponownie, zapytanie o każde konto ma limit czasu 3 sekund, konto, które nie odpowiada lub którego token dostępu wygasł, jest oceniane na podstawie zapisanych liczb, a zapis limitu ma nad nimi pierwszeństwo. Zapisane okno, którego czas resetu minął, liczy się jako 0%, więc konto, którego okno w międzyczasie się zresetowało, nie zostaje pominięte. `--auto` nigdy nie odnawia tokenów. W odróżnieniu od zwykłej rotacji `--auto` może wybrać ostatnio używane konto, jeśli nadal ma największy zapas.

### --continue

1. ccam szuka ostatniej rozmowy z tego folderu na każdym koncie: spośród plików `<session>.jsonl` w `projects/<folder>/` (nazwa folderu to folder roboczy z rozwiązanymi linkami, w którym każdy znak inny niż litera lub cyfra zamieniono na `-`, tak jak nazywa go Claude Code) tego z najnowszym promptem lub odpowiedzią. Sesje bez nich, na przykład sesja, która tylko uruchomiła `/mcp`, `/usage` lub `/status`, albo pusty plik, są pomijane. Jeśli sesja zostawiła notatkę (limit użycia lub `/ccam:switch`, zobacz niżej) i żadna rozmowa nie ma promptu ani odpowiedzi późniejszych niż 2 minuty po niej, używana jest dokładnie ta sesja; wiersze, które Claude Code zapisuje przy `/exit`, się nie liczą.
2. Konto, do którego ona należy, jest kontem, z którego się przełączasz, a nowe konto jest wybierane jak powyżej, w obrębie tej samej organizacji. Gdy tę historię współdzieli kilka kont, należy ona do tego, którego `.claude.json` mówi, że ostatnio prowadziło tę rozmowę (`lastSessionId`), a w przeciwnym razie do ostatnio używanego.
3. Jeśli nowe konto współdzieli tę historię (ten sam folder grupy), ccam uruchamia `claude --resume <session id>`, czyli właśnie tę rozmowę (`--continue` otworzyłoby najnowszą rozmowę w folderze, jakakolwiek by to była).
4. W przeciwnym razie uruchamia `claude --resume <the conversation's file> --fork-session`: Claude Code kopiuje całą rozmowę, pod nowym identyfikatorem sesji, do własnej historii nowego konta, a oryginał pozostaje bez zmian. Od tej chwili `--continue` nowego konta ją znajduje.
5. Jeśli nie podasz `--agent` po `--`, ccam przekazuje też agenta, z którym działa rozmowa: ustawienie `agent` z `.claude/settings.local.json` lub `.claude/settings.json` folderu albo z `settings.json` konta, a w przeciwnym razie agenta wskazanego w ostatnim wierszu `agent-setting` rozmowy. Claude Code nie przywraca agenta, gdy kopiuje rozmowę (krok 4), i działałby dalej jako swój domyślny agent, ze wszystkimi narzędziami. Gdy nowe konto nie współdzieli agentów, pluginów i ustawień starego (tych samych folderów `agents` i `plugins` oraz `settings.json`, przez grupę), ccam dodaje po wierszu o przeniesieniu kolejny wiersz: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

ccam podaje, jak duża jest rozmowa: nowe konto czyta ją całą ponownie przy pierwszej odpowiedzi. Pamięć podręczna promptów nigdy nie jest współdzielona między organizacjami, a każdy plan prywatny jest osobną organizacją, więc gdy oba konta należą do różnych, ccam dodaje `prompt caches are not shared between organizations, so the first reply reads it all again`. Gdy stare konto jest zalogowane do serwerów MCP, do których nowe nie jest, ccam je wymienia: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. Współdzielenie `history` w grupie pozwala uniknąć kopii, zachowuje jedną rozmowę i współdzieli też punkty kontrolne, z których `/rewind` przywraca kod (`file-history/`).

### --brief

`--continue` wznawia całą rozmowę, a nowe konto czyta ją całą ponownie przy pierwszej odpowiedzi. `ccam next --continue --brief` zaczyna zamiast tego nową rozmowę, od podsumowania starej, które ccam pisze bez pytania modelu:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

Podsumowanie zawiera, w tej kolejności i tylko tam, gdzie jest coś do powiedzenia:

1. tytuł rozmowy;
2. ścieżkę jej transkryptu, żeby Claude mógł w nim szukać, gdy potrzebuje szczegółu;
3. twoją pierwszą prośbę, gdy było ich więcej niż cztery;
4. podsumowanie, które Claude Code napisał przy ostatnim kompaktowaniu rozmowy;
5. twoje ostatnie cztery prośby, z poleceniami `!` i poleceniami z ukośnikiem tak, jak je wpisałeś;
6. ostatnią odpowiedź Claude;
7. listę zadań;
8. 10 plików, które Claude zmienił ostatnio;
9. `git status --short --branch` folderu, do 800 znaków.

Lista zadań zachowuje najwyżej 10 pozycji: najpierw odpadają pozycje ukończone, potem kroki najdalej w przyszłości, więc bieżący krok i kolejne zostają. Długie części są skracane w środku, a całe podsumowanie mieści się w 16 000 znaków.

ccam zapisuje je jako `~/.local/state/ccam/brief/<session>.md`, czytelne tylko dla ciebie, a następne `--brief` usuwa podsumowania starsze niż 7 dni. Claude Code startuje z jednym promptem i dołącza plik, o którym on wspomina: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` Prompt podany po `--` zastępuje ostatnie zdanie (`Then: <your prompt>`), a twoje pozostałe opcje po `--` idą za nim. Claude Code kończy wzmiankę na `"` i ucina ją na `#`, więc gdy ścieżka podsumowania zawiera któryś z tych znaków, ccam podaje ją jako ścieżkę względem folderu roboczego; gdy i ta ścieżka go zawiera, prompt wymienia plik z nazwy, a Claude czyta go swoim narzędziem Read, które najpierw pyta o zgodę, bo plik leży poza folderami roboczymi (w przebiegu `-p` w trybie Manual Claude Code odmawia odczytu i Claude kontynuuje bez podsumowania).

`ccam --dry-run next --continue --brief` wypisuje podsumowanie i niczego nie zapisuje. Stara rozmowa zostaje bez zmian, a nowa wie tylko to, co mówi podsumowanie, więc używaj `--brief`, gdy zadanie jest ważniejsze niż każdy szczegół tego, jak do niego doszedłeś. Bez `--continue` ccam zatrzymuje się z `--brief works with --continue`; gdy nie może odczytać rozmowy albo nie znajduje w niej niczego do podsumowania, zatrzymuje się i mówi, że `ccam next --continue` bez `--brief` wznowi ją w całości. [Relay](#relay) i [`auto_switch`](#auto-switch) zawsze wznawiają całą rozmowę. Własna opcja `--brief` Claude Code idzie po `--`, jak jego pozostałe opcje.

## W Claude Code

[Plugin ccam](/docs/install#claude-code-plugin) dodaje dwa hooki (Claude Code 2.1.78 lub nowszy):

- Gdy tura zatrzyma się na limicie użycia, zapisuje rozmowę i pokazuje powiadomienie na pulpicie w terminalach, które to obsługują (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  Zna każdy komunikat, który Claude Code uznaje za limit użycia: limity 5-godzinne, tygodniowe, zwykłe i modeli (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` dla przydziału modelu), kredyty użycia i limity wydatków oraz organizację, która zatrzymała konto.

  Gdy Claude Code zapisał odrzucone żądanie w rozmowie, hook zapisuje też pełne okno aż do jego resetu ([zapis limitu](/features/usage#the-limit-record)), z wyjątkiem zatrzymania ustawionego przez organizację, które jej administratorzy mogą w każdej chwili cofnąć. Krótkie ograniczenie po stronie serwera albo limit szybkości klucza API nie jest limitem użycia i zostaje pominięte, podobnie jak zatrzymanie subagenta i sesje nienadzorowane (przebiegi `-p` i sesje w tle).

  Claude Code 2.1.234 i nowsze także czekają w sesji i po resecie okna 5-godzinnego same kontynuują; powiadomienie daje ci wybór, żeby zamiast tego przełączyć konto, a z włączonym [`auto_switch`](#auto-switch) (Linux i macOS) ccam przełącza za ciebie, gdy inne konto ma zapas. Gdy żadne inne konto nie może przejąć pracy, bo żadne inne konto nie należy do tej samej organizacji Team lub Enterprise albo żadne nie jest innym planem prywatnym (z włączonym `cross_org` może to zrobić każde inne konto Claude), powiadomienie kończy się zamiast tego `You have no other personal plans, so ccam cannot move this conversation.` albo `You have no other accounts in Acme, so ccam cannot move this conversation.`
- `/ccam:switch [account]` zapisuje rozmowę i kończy turę, zanim zostanie wysłane jakiekolwiek żądanie: wpisz `/exit`, a potem uruchom `ccam next --continue`, żeby kontynuować ją na koncie, które podałeś (albo na następnym). W sesji, która nie ma jeszcze rozmowy, mówi, że nie ma czego kontynuować. Uruchamia je wyłącznie to, co sam wpiszesz: od wersji 0.3.2 pluginu Claude nie może uruchomić `/ccam:switch`, a zadanie zaplanowane przez Claude (CronCreate, ScheduleWakeup lub `.claude/scheduled_tasks.json`) przekazuje je Claude jako zwykły tekst.

  Gdy podane przez ciebie konto jest poza organizacją rozmowy (miejsce w innej organizacji Team lub Enterprise, plan prywatny zamiast miejsca lub odwrotnie, albo konto z kluczem API po którejkolwiek stronie), `ccam next --continue` najpierw pyta, niezależnie od tego, co mówi `cross_org`: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. Pyta tak samo przed kontem, które sam wybiera, gdy żadnego nie podasz, a które `cross_org = true` pozwala mu wziąć z innej organizacji: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. Przenosi ją tylko `y` lub `yes`; w systemach Linux i macOS klawisze wciśnięte przed pytaniem są najpierw odrzucane. Każda inna odpowiedź albo brak terminala, w którym można zapytać, zostawia rozmowę i notatkę tam, gdzie są: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (po wyborze `ccam next --continue` wybiera ponownie). Tak jak przy automatycznych wyborach, konto, którego organizacji ccam nie potrafi ustalić, pasuje do każdego.

Notatka jest przechowywana w `~/.local/state/ccam/handoff/<account>/<session>.json`, jeden plik na sesję, więc dwie sesje jednego konta nigdy nie nadpisują sobie notatek. Jest trzymana przez 12 godzin i używana jeden raz, gdy ma wystartować następna sesja; uruchomienie, które się nie powiedzie, zostawia ją na następną próbę. Należy ona do folderu, w którym sesja się zaczęła, nawet po `cd` w powłoce Claude Code. `ccam rename` aktualizuje zawarte w niej nazwy; jeśli wskazane w niej konto zostało w międzyczasie usunięte, ccam o tym informuje i wybiera jak zwykle.

### Relay

Z `relay = true` na górze `config.toml` ccam działa dalej obok sesji interaktywnych, zamiast zastępować się procesem Claude Code, i po `/ccam:switch` sam kontynuuje rozmowę tej sesji (na notatkę zostawioną przez inną sesję nie reaguje): hook kończy Claude Code, a ccam uruchamia `ccam next --continue` z notatką tej sesji. Bez podanego konta wybiera to z największym zapasem, uwzględniając model rozmowy, a gdy `cross_org = true` pozwala, by było to konto z innej organizacji, zadaje powyższe pytanie, zanim przeniesie tam rozmowę. W Windows hook nie może zakończyć Claude Code: wpisz `/exit`, a ccam ją kontynuuje. Hook nie kończy Claude Code także dla konta spoza organizacji rozmowy: wyświetla `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`, a ccam zadaje powyższe pytanie, gdy wpiszesz `/exit`. Kosztuje to jeden czekający proces ccam na sesję, dlatego domyślnie jest wyłączone.

Przekazana sesja zachowuje to, czego nie przywraca wznowienie rozmowy: opcje takie jak `--add-dir`, `--mcp-config`, listy narzędzi, `--settings`, `--plugin-dir` i opcje promptu systemowego, tryb uprawnień rozmowy, poziom wysiłku sesji i agenta, z którym działa rozmowa (zobacz [`--continue`](#continue)). Sam Claude Code przywraca z rozmowy `--model` i `--name`. Wszystko inne (`--worktree`, `--debug`, `--remote-control`, opcje, których ccam nie zna) jest wymieniane w `ccam: not carried over: …`. Przełączenie nigdy nie zadaje pytania o współdzieloną historię, a to, co wpisałeś w jego trakcie, jest odrzucane; jedyne pytanie, jakie zadaje, to pytanie przed kontem z innej organizacji.

Relay kontynuuje rozmowę tylko wtedy, gdy Claude Code zakończył się tak, jak poprosiło o to `/ccam:switch`: nie wtedy, gdy najpierw wpisałeś `/exit`, nacisnąłeś Ctrl+C w trakcie kończenia albo ccam dostał sygnał rozłączenia lub inny sygnał, ani wtedy, gdy terminala już nie ma. Notatka zostaje wtedy, a `ccam next --continue` kontynuuje ją później. Sesji uruchomionej pod relay starszej wersji ccam ani sesji nienadzorowanej (przebiegi `-p` i sesje w tle) hook nie kończy; zamiast tego dostaje ona komunikat o `/exit`. `claude remote-control` działa bez relay.

### `auto_switch`

Z `auto_switch = true` na górze `config.toml` (Linux i macOS) ccam sam przenosi sesję na inne konto, gdy Claude Code zatrzyma ją na limicie użycia: hook pluginu wybiera inne konto z zapasem i kończy Claude Code, a ccam, działający obok sesji, kontynuuje rozmowę na tamtym koncie, tak jak [relay](#relay) robi to po `/ccam:switch`. Domyślnie jest wyłączone i jest eksperymentalne.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam przenosi sesję tylko wtedy, gdy spełnione są wszystkie te warunki:

1. ccam uruchomił ją w terminalu i działa dalej obok niej. `auto_switch` włącza [relay](#relay), więc dotyczy to sesji interaktywnych, które uruchomisz po jego ustawieniu. Sesja uruchomiona wcześniej albo pod starszym ccam dostaje powiadomienie zakończone ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` Sesje nienadzorowane (przebiegi `-p` i sesje w tle), `claude remote-control` i Claude Code uruchomiony bez ccam nigdy nie są przenoszone.
2. Konto ma włączony [plugin ccam](/docs/install#claude-code-plugin) i używa Claude Code 2.1.214 lub nowszego, który przekazuje hookowi identyfikator swojego procesu (`CLAUDE_PID`).
3. Claude Code zapisał w rozmowie odrzucone żądanie, z czasem resetu, dla 5-godzinnego, tygodniowego, Opus lub Sonnet limitu tego konta. Limit wydatków, wyczerpanie kredytów użycia, zatrzymanie ustawione przez organizację, przydział Fable (`You've reached your Fable limit.`), ograniczenie po stronie serwera, limit szybkości klucza API i zatrzymanie subagenta nigdy nie przenoszą sesji.
4. ccam zna plan konta, z jego `.claude.json` albo z tego, co zapisał `ccam login`. W przeciwnym razie powiadomienie kończy się `ccam does not know personal's plan, so auto_switch leaves this session.`
5. Sesja nie jest połączona z Remote Control. Jako inne konto nie połączyłaby się ponownie z aplikacją, z której nią sterujesz, więc zostaje, a powiadomienie o tym mówi.
6. Limit resetuje się za więcej niż 15 minut. Przy krótszym czekaniu ccam zostawia sesję Claude Code: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. Od osiągnięcia limitu nie wpisałeś żadnego promptu.
8. Inne konto ma zapas. ccam wybiera tak jak [`--auto`](#auto) (`auto_max_usage`, `auto_strategy`, model rozmowy), spośród kont tego samego rodzaju co to, które opuszczasz: inne miejsce w tej samej organizacji Team lub Enterprise albo inny plan prywatny (Pro, Max). `cross_org` i `--cross-org` nie mają tu zastosowania. Konto musi:
   - być zalogowane loginem, który Claude Code potrafi odnowić;
   - nie być zalogowane na to samo konto Claude;
   - nie używać klucza API;
   - mieć plan znany ccam;
   - nie mieć włączonych kredytów użycia, na ile pokazuje to własne zapytanie ccam o użycie;
   - mieć ten folder oznaczony jako zaufany w Claude Code;
   - jeśli rozmowa działa z agentem (wskazuje go ustawienie `agent` albo rozmowa), współdzielić agenty konta, które opuszczasz: jego foldery `agents` i `plugins` oraz jego `settings.json` muszą być, każde z osobna, dokładnie te same (przez grupę) albo nie istnieć w żadnym z obu kont.

   Gdy te trzy ostatnie warunki wykluczają każde konto z zapasem, powiadomienie je wymienia: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` albo `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` Okno z [zapisem limitu](/features/usage#the-limit-record) liczy się jako pełne, a konto, na które inna sesja przeniosła się w ciągu ostatnich 2 minut, ustępuje pozostałym. Wybór musi zapaść w ciągu około 6 sekund; gdy zapytania o użycie trwają dłużej, sesja zostaje. Z `usage_source = "claude"` `auto_switch` decyduje wyłącznie na podstawie zapisanych liczb, a ccam nie potrafi ustalić z zapytania Claude Code ani z paska stanu, czy kredyty użycia są włączone, więc zakłada, że są wyłączone.

Wtedy hook kończy Claude Code, a ccam, działający obok niego, pokazuje powiadomienie na pulpicie w terminalach, które to obsługują, na przykład `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` Wypisuje pierwszy wiersz przykładowego wyjścia powyżej i uruchamia `ccam next --continue` dla tej sesji jako wybrane konto, co wypisuje jego dwa pozostałe wiersze i uruchamia Claude Code z opcjami sesji, trybem uprawnień, poziomem wysiłku i agentem, tak jak [relay](#relay). Pierwsza wiadomość nowej sesji to ten stały tekst, jedyne, co ccam wysyła w twoim imieniu (przeniesienie po `/ccam:switch` niczego nie wysyła):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

Jeśli i to konto osiągnie limit, tam dzieje się to samo. Zapisy limitów nie pozwalają łańcuchowi wrócić na konto przed jego resetem, więc gdy wszystkie konta są pełne, sesja zostaje na ostatnim.

Hook dowiaduje się, że rozmowa działa z agentem, z ustawienia `agent` albo z wiersza wskazującego agenta, który Claude Code co jakiś czas i przy końcu sesji zapisuje w rozmowie. Jeśli na razie wskazuje go tylko `--agent` sesji, a wybrane konto nie współdzieli agentów, `ccam next --continue` niczego nie uruchamia i zachowuje notatkę, a drugie powiadomienie na pulpicie o tym informuje: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

Gdy żadne inne konto nie ma zapasu, nic się nie przenosi: Claude Code obsługuje limit jak zwykle, czekając w sesji albo proponując swoje opcje, a powiadomienie mówi, które konto zwolni się pierwsze, na przykład `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` W przypadku miejsc mówi `None of your other accounts in Acme has room now`. Aby zatrzymać przeniesienie, które już się zaczęło, naciśnij Ctrl+C, gdy Claude Code się zamyka: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

Przeniesienie ma swoje koszty. Claude Code uruchamia się ponownie, co trwa kilka sekund, a powłoki w tle, serwery deweloperskie i agenty uruchomione przez Claude zatrzymują się razem ze starą sesją. Nowe konto czyta całą rozmowę ponownie przy pierwszej odpowiedzi, jak przy [`--continue`](#continue), a ccam podaje, jak jest duża; przeniesienie nigdy nie używa [`--brief`](#brief).

`auto_switch` nigdy nie przełącza, zanim Claude Code zgłosi limit, nigdy nie śledzi twojego użycia w tle i nigdy nie kopiuje, nie odnawia ani nie przekazuje loginu. W systemie Windows jest ignorowane i `ccam doctor` o tym informuje; tam, z `relay = true`, `/ccam:switch` i `/exit` nadal przenoszą rozmowę na inne konto.

**Warunki Anthropic.** Na miejscach Team i Enterprise `auto_switch` przenosi twoją pracę na inne miejsce w tej samej organizacji. Centrum pomocy Anthropic podaje: „Each seat on a Team plan is assigned to one named individual, and only that person may use it” (każde miejsce w planie Team jest przypisane do jednej imiennie wskazanej osoby i tylko ta osoba może z niego korzystać) oraz „Users can't route their work through another user's account to use that user's remaining usage” (użytkownicy nie mogą kierować swojej pracy przez konto innego użytkownika, aby wykorzystać pozostały limit użycia tego użytkownika) ([plany Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); podobnie w przypadku [planów Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)).

W przypadku planów prywatnych (Pro, Max) warunki Anthropic nie mówią, czy twoje własne konta mogą automatycznie przejmować pracę od siebie nawzajem na limicie. [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) zabraniają korzystania z usług „through automated or non-human means, whether through a bot, script, or otherwise” (w sposób zautomatyzowany lub inny niż przez człowieka, czy to za pomocą bota, skryptu, czy w inny sposób), chyba że Anthropic wyraźnie na to pozwoli, a przeniesienie wysyła w twoim imieniu powyższą stałą wiadomość. [Usage Policy](https://www.anthropic.com/legal/aup), która wchodzi w życie 12 listopada 2026 r., wymienia wśród zakazanych działań „Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products” (celowe obchodzenie lub zakłócanie funkcji, ograniczeń lub zabezpieczeń ustanowionych w naszych produktach).

ccam nie twierdzi, że `auto_switch` jest zgodny z tymi warunkami. Czy go włączyć, decydujesz ty, i ty ponosisz ryzyko: Anthropic może zawiesić konto albo dostęp organizacji Team lub Enterprise, gdy uzna, że warunki zostały naruszone.

Gdy `auto_switch` jest włączony, `ccam doctor` ostrzega o koncie bez pluginu ccam i o wersjach Claude Code starszych niż 2.1.214, a w przypadku miejsc przypomina zasadę Anthropic dotyczącą miejsc ([sprawdzenia](/features/doctor)).

## Ponowienie po osiągnięciu limitu

`ccam next` służy do sesji interaktywnej, która się zatrzymała. Przebieg nieinteraktywny, `claude -p` w skrypcie lub zaplanowanym zadaniu, nie może na ciebie czekać, więc `--fallback` przed nazwą konta pozwala ccam przełączyć konto za niego:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` działa tylko razem z `-p` lub `--print`; w przeciwnym razie ccam zatrzymuje się i odsyła do `ccam next --auto --continue`.
2. Standardowe wejście z potoku jest odczytywane raz, do 10 MiB, i przekazywane każdej próbie. Tak jak Claude Code, ccam czeka najwyżej 3 sekundy, aż zacznie ono napływać (z wyjątkiem `--input-format stream-json`); jeśli potok pozostaje otwarty, ale nic nie wysyła, przebieg nie dostaje żadnego wejścia, a ccam informuje o tym na standardowym wyjściu błędów.
3. ccam uruchamia `claude` jako proces potomny zamiast przez `exec`, wstrzymuje jego standardowe wyjście, a standardowe wyjście błędów przekazuje na bieżąco.
4. Gdy `claude` zakończy się kodem 0 albo nie powiedzie się z innego powodu, wstrzymane wyjście i kod wyjścia są przekazywane bez zmian. Tak samo przebieg zakończony sygnałem (kod wyjścia powyżej 128): poproszono go o zatrzymanie, więc nie jest uruchamiany ponownie.
5. Błąd jest uznawany za limit użycia, gdy tak mówi komunikat Claude Code, czyli dowolny z komunikatów o limicie użycia znanych [hookowi](#inside-claude-code), także w wyjściu `--output-format json` lub `stream-json`, albo gdy użycie konta pokazuje okno 5-godzinne, tygodniowe lub `--model` na poziomie 100% (łącznie z zapisem limitu).
6. Wtedy wyjście tej próby jest odrzucane, a to samo polecenie uruchamia się na innym zalogowanym koncie: najpierw na tym z najniższym użyciem 5-godzinnym, potem na kontach, których użycia nie da się sprawdzić, według numerów. Konta na poziomie 100% są pomijane, konta z loginem kluczem API nigdy nie są próbowane, a każde konto jest próbowane najwyżej raz.
7. Gdy wśród argumentów jest `-c`, `--continue`, `-r` lub `--resume`, próbowane są tylko konta z tej samej grupy historii, żeby dało się znaleźć rozmowę.
8. Próbowane są tylko konta z tej samej organizacji, tak jak w `--auto` (`cross_org = true` znosi to ograniczenie).

Następne konto jest wybierane według tych samych zasad co w `--auto` (`auto_max_usage`, `auto_strategy`, okno `--model`). Gdy nie zostanie żadne konto, ccam wypisuje wyjście ostatniej próby, informuje, że żadne zalogowane konto nie ma już zapasu, i kończy działanie z kodem 75. Z `--wait` przed nazwą konta (`ccam --fallback --wait=2h work -p …`) zamiast tego czeka na reset pierwszego konta, a potem może ponownie spróbować każdego konta; zawsze najpierw czeka co najmniej raz, żeby konto, które właśnie osiągnęło limit, nie zostało od razu uruchomione ponownie. Czas oczekiwania ogranicza łączny czas czekania przebiegu, a nie każde oczekiwanie z osobna. Zapisane ze spacją (`--wait 30m`) jest odrzucane z tą samą podpowiedzią co w `ccam next`. Początkowe `--` przed argumentami Claude Code działa jak przy każdym uruchomieniu (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` wypisuje polecenie pierwszej próby.

Sesje interaktywne nie przełączają się same, chyba że włączysz [`auto_switch`](#auto-switch). Po `ccam statusline install` pasek stanu dodaje `→ ccam next --auto --continue`, gdy użycie 5-godzinne lub tygodniowe osiągnie 90% ([Limity użycia](/features/usage)).

## Bezpieczeństwo i przypadki brzegowe

- Konta bez loginu nigdy nie są wybierane. Gdy nie ma innego zalogowanego konta, samo `ccam next` kończy się komunikatem `no logged-in account other than <name>`.
- Z `--continue`, gdy w folderze nie ma żadnej rozmowy: `no account has a conversation in <folder> to continue`.
- Gdy zasada organizacji nie zostawia żadnego konta: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Z `--auto`, gdy nie ma zalogowanego konta: `no logged-in account to pick`. Gdy każde konto jest pełne lub ma nieznane użycie: `` no account with known usage and headroom (see `ccam usage`) ``. W obu przypadkach kod wyjścia to 75.
- `--max-usage`, `--strategy`, `--wait` i `--cross-org` działają tylko z `--auto`, a wartość spoza zakresu 1-100 lub nieznana strategia jest odrzucana.
- Fork kopiuje rozmowę, a nie pozostały stan konta: punkty kontrolne z `file-history/` pierwotnego konta nie są przenoszone, chyba że historia jest współdzielona.
- `--auto` wysyła zapytania sieciowe nawet z `--dry-run`. Tylko odczytuje i nigdy nie zapisuje plików z tokenami.
- Wybrane konto startuje tą samą ścieżką co `ccam <name>`, więc obowiązują te same sprawdzenia, w tym ostrzeżenie, gdy w tym folderze działa już inna sesja Claude Code współdzieląca tę samą historię.
- `last-used` jest zapisywany tylko wtedy, gdy Claude Code rzeczywiście wystartuje. `--dry-run` go nie zmienia, więc próbne uruchomienie nie przesuwa rotacji.
- `--fallback` wstrzymuje standardowe wyjście do końca próby, więc wyjście strumieniowe, takie jak `--output-format stream-json`, przychodzi w całości na końcu. Klient, który przesyła strumieniowo wejście `stream-json` i czeka na odpowiedzi, zanim wyśle więcej, nie może używać `--fallback`. Pliki zmienione przez pierwszą próbę przed osiągnięciem limitu pozostają zmienione; następna próba zaczyna prompt od nowa.

## Powiązane

- [Limity użycia](/features/usage)
- [Limity użycia i przełączanie kont](/docs/usage-limits)
- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Co zrobić, gdy Claude Code osiągnie limit użycia](/guides/usage-limit-workflow)
