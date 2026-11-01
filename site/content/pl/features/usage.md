---
title: Limity użycia
description: ccam usage pokazuje naraz 5-godzinne i tygodniowe limity każdego konta, a ccam statusline wyświetla w Claude Code nazwę konta i jego użycie.
updated: 2026-11-01
order: 7
translation: machine
sourceHash: "54b296f7da35"
---

`ccam usage` pokazuje dla każdego konta, jaka część każdego limitu jest wykorzystana: 5-godzinnego limitu, tygodniowego limitu i własnego tygodniowego limitu każdego modelu (takiego jak Fable). Każdy limit ma pasek i czas do jego resetu. Używaj go przed rozpoczęciem długiego zadania albo gdy Claude Code informuje, że limit jest blisko, żeby zobaczyć, które konto ma jeszcze zapas. `ccam statusline` to polecenie paska stanu dla Claude Code: pokazuje nazwę konta i użycie w Claude Code i zapisuje liczby dla `ccam usage`. Tokeny i szacowane koszty stojące za tym użyciem pokazuje `ccam ccusage`, które uruchamia [ccusage](https://ccusage.com) na rozmowach każdego konta ([niżej](#tokens-and-costs-ccam-ccusage)).

## Użycie

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts that are not running
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

Każde konto ma wiersz ze swoim numerem, nazwą i informacją, jak stare są jego liczby; `(cached)` oznacza, że pochodzą z zapisanego pliku, a nie z endpointu. Pod nim każdy limit ma własny wiersz z 20-polowym paskiem, procentem i czasem do resetu okna: `5-hour`, `Weekly` (wszystkie modele), a potem własne tygodniowe okno każdego modelu, takie jak `Fable`. Okno, którego konto nie ma, pokazuje `-`. Procent zmienia kolor na żółty przy 80% i na czerwony przy 100%.

Jeśli Claude Code odrzucił żądanie, bo okno było pełne, to okno pokazuje 100% aż do resetu, bez względu na zapisane liczby (zobacz [zapis limitu](#the-limit-record)). W terminalu każde konto pojawia się od razu jako `fetching…`, a jego wiersze są rysowane, gdy odpowie.

Ostatni wiersz, widoczny tylko wtedy, gdy masz więcej niż jedno konto, wskazuje konto, które wybrałby [`ccam next --auto`](/features/next). Stosuje te same zasady:

- `auto_max_usage` i `auto_strategy`;
- uruchomione sesje;
- tylko konta organizacji ostatnio używanego konta, chyba że ustawiono `cross_org`;
- każde okno na poziomie 100% jest pomijane, więc konto z wyczerpanym tygodniem nigdy nie jest wskazywane.

Przy domyślnej strategii `headroom` brzmi `→ Most headroom: <name>`; przy `expiring` lub `balance` podaje przyczynę, na przykład `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Miejsca w planach Team i Enterprise.** Automatyczny wybór konta (`ccam next`, `--fallback`, `ccam pick`, klawisz `n` w panelu oraz, gdy go włączysz, [`auto_switch`](/features/next#auto-switch)) przełącza też między miejscami w tej samej organizacji Claude Team lub Enterprise, co może naruszać warunki Anthropic. Przeczytaj [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) i skonsultuj to ze swoją organizacją i z Anthropic, zanim zaczniesz na nim polegać.

Konto, które nie jest zalogowane, pokazuje `not logged in → ccam login <name>`, a konto z kluczem API pokazuje `no limits (API)`. Wygasły token dostępu pokazuje ``access token expired; `ccam <name>` renews it when it starts``, nad zapisanymi liczbami, jeśli są: Claude Code odnowi go przy następnym uruchomieniu. Login, który całkowicie wygasł, bo wygasł jego token odświeżania albo Claude Code wyczyścił go po odrzuceniu przez serwer, pokazuje `` login expired → `ccam login <name>` ``.

Po `ccam statusline install` Claude Code pokazuje wiersz taki jak `[work] · Opus · 5h 24% · 7d 41%`.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--json` | Wypisuje raporty jako JSON: `name`, `slot`, `status`, `limits` (z `credits`, gdy konto ma włączone kredyty użycia), `source`, `at`, `error` oraz `limited` i `limited_windows`, gdy zapis limitu ustawia okna na 100%. Bez kont wypisuje `[]` | wyłączone |
| `--watch[=N]` | Odświeża co `N` sekund, aż naciśniesz Ctrl+C. Samo `--watch` albo wartość 0 lub mniejsza oznacza 60. Wartość podaje się po `=`: `--watch 10` jest odrzucane z podpowiedzią | wyłączone |
| `-v`, `--verbose` | Dodaje pod każdym kontem wiersz `source:`: skąd pochodzą jego liczby (endpoint użycia, Claude Code, pasek stanu, zapisany plik albo limit użycia zgłoszony przez Claude Code). Gdy zapis limitu ma pierwszeństwo przed liczbami, pojawia się po nim wiersz taki jak `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` | wyłączone |
| `--no-cache` | Ignoruje 60-sekundową pamięć podręczną i pyta endpoint ponownie | wyłączone |
| `--source live\|claude\|cache` | `live` pyta endpoint z tokenem konta; `claude` zleca sprawdzenie Claude Code (poniżej); `cache` pokazuje tylko zapisane liczby | `usage_source`, w przeciwnym razie `live` |
| `--refresh` | Odnawia wygasłe tokeny dostępu, tylko dla kont, które nie są uruchomione. Odnowienie może potrwać do 30 sekund | wyłączone |
| `--dry-run` | Akceptowana: `ccam usage` nie zmienia niczego poza własną pamięcią podręczną. Odrzucana razem z `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | wyłączone |

Każda inna wartość `--source` kończy się błędem `--source must be live, claude or cache`. `ccam usage` nie przyjmuje argumentów.

### Użycie z samego Claude Code

Z `usage_source = "claude"` na górze `config.toml` (nad pierwszym `[[group]]`, `[[profile]]` lub `[[folder]]`) albo z `--source claude` ccam w ogóle nie odczytuje tokena logowania. Dla każdego konta uruchamia `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` w folderze konta i wysyła na standardowe wejście żądanie `get_usage` Claude Code (Claude Code 2.1.270 lub nowszy). Claude Code sprawdza użycie z własnym loginem i odpowiada oknami 5-godzinnym, tygodniowym i modeli.

Dzięki `--safe-mode` hooki, serwery MCP i pluginy konta się nie uruchamiają, a żadna sesja nie jest zapisywana. Każde konto wymaga jednego uruchomienia Claude Code (około sekundy), więc to źródło jest wolniejsze niż `live`. Konta są pytane naraz, każde z limitem czasu 20 sekund, a 60-sekundowa pamięć podręczna działa jak zwykle.

Login, którego token dostępu wygasł lub wygaśnie w ciągu 5 minut, dostaje 50 sekund, bo Claude Code odnawia go przy starcie, a przerwane odnowienie może stracić login. Jeśli ccam zrezygnuje z zapytania, prosi Claude Code o zakończenie sygnałem SIGTERM, więc Claude Code dokończy rozpoczęte odnowienie (Windows nie ma SIGTERM, więc tam proces jest kończony).

Login, który całkowicie wygasł, pokazuje `` login expired → `ccam login <name>` `` bez uruchamiania Claude Code. Claude Code oznacza to żądanie jako eksperymentalne, więc jeśli przestanie odpowiadać, ccam pokazuje zapisane liczby, a `ccam doctor` to zgłasza. Gdy Claude Code odpowie, że jego własne sprawdzenie się nie powiodło, ccam pokazuje zapisane liczby z komunikatem `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`, a nie plan bez limitów. Z `--dry-run` ccam nie uruchamia Claude Code, bo Claude Code zapisuje w folderze konta; zamiast tego pokazuje zapisane liczby. `usage_source` przyjmuje `live` lub `claude`; `cache` to wyłącznie wartość `--source`.

`ccam statusline` odczytuje JSON paska stanu Claude Code ze standardowego wejścia. Jego jedyna opcja to `--wrap <command>`: istniejące polecenie paska stanu, przed którego pierwszym wierszem wyniku pojawia się nazwa konta. `ccam statusline install` nie przyjmuje opcji. Oba odrzucają `--dry-run` komunikatem `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Tokeny i koszty: `ccam ccusage`

`ccam usage` pokazuje, jaką część każdego limitu wykorzystało konto. [ccusage](https://ccusage.com) liczy natomiast tokeny i szacowane koszty w logach rozmów Claude Code. Uruchomione samodzielnie szuka logów Claude Code tylko w `~/.claude` i `~/.config/claude` (albo w folderach wskazanych przez `CLAUDE_CONFIG_DIR`), więc jedyne konto ccam, które może zobaczyć, to konto zaimportowane w miejscu w `~/.claude`. `ccam ccusage` uruchamia je z `CLAUDE_CONFIG_DIR` ustawionym na foldery wszystkich kont, które mają rozmowy, więc jeden raport obejmuje je wszystkie:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam najpierw odczytuje własne opcje: `-a` (`--account`) przyjmuje nazwy lub numery kont rozdzielone przecinkami, a samą opcję można powtarzać; `--dry-run` wypisuje polecenie zamiast je uruchamiać. Od pierwszego argumentu, który nie jest opcją ccam, wszystko trafia do ccusage tak, jak to wpisałeś, z wyjątkiem `--config` dla `bun x` i `npx` (poniżej); `ccam ccusage -- --help` pokazuje pomoc samego ccusage. ccam uruchamia `ccusage`, jeśli jest zainstalowane, w przeciwnym razie `bun x ccusage`, a w ostateczności `npx ccusage@latest`; te dwa ostatnie pobierają ccusage z npm. Uruchamiają się one we własnym folderze ccam (`~/.local/state/ccam/ccusage`), a nie w folderze, w którym jesteś, więc znajdujący się tam projekt nie może wybrać, co uruchamiają, za pomocą swojego `node_modules/.bin`, `bin` w swoim `package.json` ani rejestru w swoim `.npmrc`; `ccam --dry-run ccusage` pokazuje `cd` do tego folderu. Gdy żadne z nich nie jest dostępne, ccam podpowiada, żeby zainstalować Node.js, Bun albo ccusage. W systemach Linux i macOS ccusage zajmuje miejsce ccam; w systemie Windows ccam czeka, aż ccusage się zakończy, i kończy działanie z jego kodem wyjścia. W WSL ccam pomija `ccusage`, `bun` lub `npx` leżące na dysku Windows, bo taki program uruchomiłby się jako program Windows bez folderów kont.

- Raport obejmuje konta, których foldery zawierają rozmowy (folder `projects`). Rozmowa, którą konta współdzielą przez grupę, jest liczona tylko raz.
- Z `-a` konto, które współdzieli swoją historię przez grupę, wnosi rozmowy całej grupy, a ccam podaje, czyje to rozmowy: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- `blocks` pokazuje 5-godzinne okna, a każde konto ma własne, więc podaj przy nim jedno konto, jak wyżej.
- ccusage raportuje też inne znalezione narzędzia do kodowania, takie jak Codex, chyba że wskażesz jedno: `ccam ccusage claude daily` to sam Claude Code.
- Gdy uruchamia je ccam, ccusage czyta swoje ustawienia (`ccusage.json`) tylko z `./.ccusage` i z folderów kont. `--config` wskazuje inny plik: `ccam ccusage daily --config ~/.config/claude/ccusage.json`. Ponieważ `bun x` i `npx` uruchamiają się we własnym folderze ccam, ccam przekazuje im `./.ccusage/ccusage.json` przez `--config`, gdy ten plik istnieje i nie wskazujesz żadnego pliku, a względne `--config` zamienia na pełną ścieżkę. ccusage czyta wtedy inne ścieżki względne, takie jak `--pi-path`, z folderu ccam, więc podawaj je w całości.
- ccusage odczytałoby folder, którego ścieżka zawiera przecinek, jako dwa foldery, więc ccam go odrzuca i podaje nazwę konta; pomiń to konto za pomocą `-a`.
- `ccusage` to polecenie ccam, więc konto o nazwie ccusage uruchamia się przez `ccam run ccusage`; `ccam doctor` o tym informuje.

## Jak to działa

Liczby, które pokazuje `ccam usage`, pochodzą z dwóch miejsc: z nieudokumentowanego endpointu użycia Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, wywoływanego z własnym loginem każdego konta, oraz z wartości, które Claude Code przekazuje poleceniu paska stanu. Cały ten kod jest w `internal/usage`.

1. ccam wczytuje `config.toml` i tworzy jeden `usage.Account` na konto (`cli.usageCmd`): jego folder `~/.local/share/ccam/profiles/<name>/`, znaczniki sesji (razem z własnymi zapisami Claude Code w folderze konta mówią, czy konto działa) i informację, czy zostało zaimportowane w miejscu.
2. ccam sprawdza wszystkie konta naraz, po jednej gorutynie na konto (`Fetcher.FetchAll`). Wszystkie współdzielą jeden `http.Client`, więc połączenia z endpointem są używane ponownie. Nawiązanie połączenia może trwać 2 s, a każde konto ma łącznie 3 s. Klient nie podąża za przekierowaniami, bo przekierowanie wysłałoby token gdzie indziej, i odczytuje najwyżej 1 MB odpowiedzi.
3. Gdy każde konto kończy, jego wiersz jest odświeżany. Konto, którego zapytanie się nie powiedzie lub które nie odpowie na czas, korzysta z zapisanych liczb, więc jedno wolne konto nigdy nie wstrzymuje pozostałych.

Na diagramie każde konto jest sprawdzane we własnej gorutynie, a to, które przekroczy czas, pokazuje zamiast tego zapisane liczby:

<!-- motion:usage-fetch -->

`personal` i `work` odpowiadają na czas i pokazują aktualne liczby. `side` potrzebuje więcej niż 3 s, więc ccam pokazuje zamiast tego jego zapisane liczby (`fallback()`). Potem rysuje blok dla każdego konta (`usage.Table()`) i wybiera konto wskazane w ostatnim wierszu (`usage.Best()`).

Dla jednego konta ccam wykonuje te sprawdzenia w kolejności (`Fetcher.Fetch`):

<!-- motion:usage-cache -->

1. Z `--source cache` przechodzi od razu do `fallback()`.
2. Jeśli nie ustawiono `--no-cache`, aktualny wynik w `~/.local/state/ccam/usage/<name>.json` zmierzony mniej niż 60 sekund temu (`CacheTTL`) jest używany i oceniany według bieżącej chwili: okno, które się w międzyczasie zresetowało, pokazuje 0%, a plik z datą późniejszą niż czas zegara nie liczy się jako świeży. Dwukrotne uruchomienie `ccam usage` pod rząd kosztuje jedno zapytanie na konto.
3. ccam odczytuje `.credentials.json` (lub pęk kluczy macOS) bez zmieniania go i dekoduje tylko klucz `claudeAiOauth` (`readCredentials()`). Brak loginu subskrypcji oznacza `not logged in` albo `no limits (API)`, gdy `.claude.json` zawiera klucz API. Login, którego token odświeżania wygasł albo który wyczyścił Claude Code, to `` login expired → `ccam login <name>` `` i nie jest sprawdzany.
4. Wygasły token dostępu nie jest używany. Bez `--refresh` konto pokazuje swoje zapisane liczby z ``access token expired; `ccam <name>` renews it when it starts``.
5. W przeciwnym razie ccam wywołuje endpoint. Kod 401 lub 403 jest traktowany jako wygasły login (`login expired: authentication failed (401)`), każdy inny status lub treść, której nie da się odczytać, jako błąd, i w obu przypadkach używane są zapisane liczby. Zapytanie przerwane przez Ctrl+C wyświetla `cancelled`. Kod 429 oznacza, że endpoint otrzymał zbyt wiele zapytań dla tego konta: ccam zapisuje `<name>.pause` z czasem kolejnego zapytania (`Retry-After` albo 5 minut, najwyżej godzina), pokazuje zapisane liczby z komunikatem `rate limited by the usage endpoint; asking again after 15:10` i do tego czasu nie wywołuje endpointu dla tego konta, nawet z `--no-cache`. Natychmiastowe ponawianie tylko przedłużałoby ograniczenie.
6. Poprawna odpowiedź jest zapisywana w `<name>.json` razem z czasem pomiaru.

Dla zapisanych liczb ccam bierze najnowsze zapisane sprawdzenie (`<name>.json` albo `<name>.claude.json` ze sprawdzeń przez Claude Code), a gdy `<name>.statusline.json` zapisano później, nakłada na nie 5-godzinne i tygodniowe okno z paska stanu (`fallback()`); okna modeli, których pasek stanu nie ma, zostają takie, jak zapisało je sprawdzenie. Zanim pokaże zapisane liczby, sprawdza czas resetu każdego okna: okno, które zresetowało się od zapisania liczb, pokazuje 0%, bo stary procent już nie obowiązuje. Dlatego `personal` powyżej pokazuje `0%` bez czasu resetu.

Endpoint ma dwa formaty odpowiedzi. Nowszy ma tablicę `limits` z wpisami `session` (5-godzinny), `weekly_all` i `weekly_scoped` (dla poszczególnych modeli). Starszy ma obiekty `five_hour`, `seven_day`, `seven_day_opus` i `seven_day_sonnet`. ccam odczytuje oba (`ParseLive`): tablica `limits` ma pierwszeństwo, a stare obiekty uzupełniają to, czego ona nie obejmuje. Treść bez żadnych limitów jest zgłaszana jako zmieniony format, o czym informuje `ccam doctor`.

### Zapis limitu

Gdy Claude Code zatrzyma turę na limicie użycia, hook [pluginu ccam](/docs/install#claude-code-plugin) odczytuje z rozmowy odrzucone żądanie: które okno było pełne (5-godzinne, tygodniowe, Opus lub Sonnet) i kiedy się resetuje. Zapisuje to w `~/.local/state/ccam/usage/<name>.limit.json`. ccam nakłada ten zapis na wszystko inne, co wie, z każdego źródła (`Fetch`): każde okno w zapisie liczy się jako 100% aż do resetu, bez względu na to, jak stare lub niskie są zapisane liczby, więc `ccam next --auto`, `ccam pick`, `--fallback` i [`auto_switch`](/features/next#auto-switch) nie wybiorą tego konta wcześniej.

Gdy nic innego nie wiadomo, konto pokazuje sam zapis (`source: the usage limit Claude Code reported` z `-v`). Sesje nienadzorowane (przebiegi `-p` i sesje w tle) nie tworzą zapisu, podobnie jak zatrzymanie ustawione przez organizację (na przykład gdy skończyło się jej użycie albo limit grupy wynosi $0), które jej administratorzy mogą w każdej chwili cofnąć. `ccam login` na inne konto Claude usuwa zapis; usunięcie pliku także go czyści.

### Odnawianie tokenów

Tokeny odświeżania są rotowane: odnowienie daje nowy token odświeżania i unieważnia stary. Gdyby ccam odnowił login, gdy Claude Code działa jako to samo konto, ta sesja zostałaby wylogowana. Dlatego `ccam usage` odczytuje tokeny i nigdy ich nie zapisuje, chyba że podasz `--refresh`. Wtedy ccam (`Fetcher.refresh`):

1. odrzuca konta zarejestrowane w miejscu przez `ccam import`, bo zwykłe `claude` może używać tego loginu bez wiedzy ccam;
2. odrzuca loginy przechowywane w pęku kluczy macOS;
3. bierze plik blokady `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, żeby dwa procesy ccam nigdy nie odnawiały tego samego tokenu jednocześnie; dopóki jest trzymany, konto pokazuje `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)`. `ccam rename`, `remove` i `logout` też go biorą i odmawiają, gdy trwa odnowienie;
4. pod blokadą odmawia, jeśli działa sesja tego konta, niezależnie od sposobu uruchomienia (przez ccam, zwykłe `claude`, IDE albo sesja w tle i stojący za nią demon, który sam odnawia login), i pomija zapytanie, jeśli inny proces już odnowił token;
5. nic nie wysyła, dopóki Claude Code trzyma własne blokady odnowienia w folderze konta (`.oauth_refresh.lock`, `<folder>.lock`): Claude Code właśnie odnawia ten login;
6. wysyła token odświeżania do endpointu tokenów i czeka na odpowiedź do 30 sekund, tyle, ile Claude Code daje własnym odnowieniom. Wysłane żądanie nigdy nie jest przerywane, nawet przez Ctrl+C: serwer mógł już wymienić token odświeżania, a nowy login zawiera tylko jego odpowiedź. Ctrl+C lub rozłączenie kończy pozostałe zapytania, a ccam kończy działanie (130 po Ctrl+C), gdy odpowiedź zostanie zapisana;
7. zapisuje pod blokadą zapisu Claude Code (`.storage-write.lock`) i zmienia tylko pola tokenów w `.credentials.json` w jego bieżącej postaci: to, co Claude Code w międzyczasie zapisał (logowanie MCP), zostaje, login odnowiony w międzyczasie przez Claude Code jest zachowywany, a wylogowanie nigdy nie jest cofane. Gdy odpowiedź podaje czas życia nowego tokenu odświeżania, jest on zapisywany jako `refreshTokenExpiresAt`, tak jak robi to Claude Code.

Odnowienie, które się nie udało, podaje przyczynę w wierszu konta albo nad jego zapisanymi liczbami: na przykład `token not renewed: a session is running` albo `login may be lost: …`, gdy żądanie zostało wysłane, ale jego odpowiedzi nie dało się użyć ani zapisać, z radą, żeby uruchomić `ccam login <name>`, jeśli Claude Code potem poprosi o zalogowanie.

### Zbieranie danych z paska stanu

Claude Code uruchamia polecenie paska stanu i przekazuje mu potokiem dokument JSON. ccam odczytuje z niego nazwę modelu i `rate_limits`:

<!-- motion:usage-statusline -->

1. ccam znajduje konto na podstawie `CCAM_PROFILE`, które ustawia przy uruchamianiu Claude Code, a w przeciwnym razie przez dopasowanie `CLAUDE_CONFIG_DIR` do folderów kont (`runStatusline`).
2. ccam odczytuje `rate_limits.five_hour` i `rate_limits.seven_day` (`ParseStatusLine`). Dla kont z kluczem API oraz przed pierwszą odpowiedzią w sesji JSON nie zawiera limitów; wtedy nic nie jest zapisywane.
3. ccam zapisuje `~/.local/state/ccam/usage/<name>.statusline.json` (`SaveStatusLine`). `ccam usage` korzysta z niego, gdy endpoint nie odpowiada.
4. ccam wypisuje `[name] · model · 5h N% · 7d N%`. Z `--wrap` uruchamia twoje polecenie przez `sh -c` z tym samym wejściem i wypisuje jego wynik tak, jak Claude Code pokazuje pasek stanu: każdy niepusty wiersz, z `[name]` przed pierwszym. Polecenie, które się nie powiedzie, pokazuje samą nazwę konta. W Windows uruchamia polecenie tak jak Claude Code: przez Git Bash, jeśli jest zainstalowany, a w przeciwnym razie przez PowerShell.
5. Gdy 5-godzinne lub tygodniowe użycie wynosi 90% lub więcej i istnieje inne konto, wiersz kończy się `→ ccam next --auto --continue`, na przykład `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Opakowane polecenie nie ma własnego limitu czasu: Claude Code daje poleceniu paska stanu 10 minut i zatrzymuje przebieg, gdy zaczyna się nowsza aktualizacja, więc wolne polecenie wygląda tak, jak bez ccam. Gdy Claude Code kończy ccam, polecenie też się kończy (na Linux i macOS razem z każdym procesem, który uruchomiło).

`ccam statusline install` zapisuje `statusLine` w `settings.json` każdego konta i zachowuje pozostałe klucze obiektu, takie jak `padding`, `refreshInterval` i `hideVimModeIndicator`. Gdy konto ma już polecenie paska stanu, ccam zachowuje je jako `--wrap '<old command>'`; w Windows jako `--wrap-base64 <base64>`, które w żadnej powłoce nie wymaga cudzysłowów. Jeśli pasek stanu już uruchamia ccam (z innej ścieżki albo w cudzysłowach zapisanych przez starszą wersję), ccam go zastępuje zamiast opakowywać i zachowuje polecenie, które on opakowuje.

W Windows ścieżka do ccam jest zapisywana z ukośnikami zwykłymi i `~` w miejsce folderu domowego (`~/.local/bin/ccam.exe statusline`), co rozumieją zarówno Git Bash, jak i PowerShell. Ścieżka wymagająca cudzysłowów jest zapisywana w całości dla powłoki, której używa Claude Code: `'…'` dla Git Bash, `& '…'` dla PowerShell. `settings.json` współdzielony przez grupę jest zapisywany raz, przez swój link, więc dostaje go każde konto w grupie. Ponowne uruchomienie wypisuje `(already set)` dla plików, które już używają ccam.

## Bezpieczeństwo i przypadki brzegowe

- `ccam usage` nigdy nie zmienia pliku z tokenami, chyba że podasz `--refresh`. Gdy token dostępu wygaśnie, uruchom konto raz przez `ccam <name>`, a Claude Code go odnowi. Login, który całkowicie wygasł, wymaga `ccam login <name>`.
- `--refresh` nigdy nie rusza działających kont (nawet gdy Claude Code uruchomiono bez ccam), zaimportowanych kont ani loginów z pęku kluczy. Odrzucone odnowienie pokazuje zapisane liczby i przyczynę (`token not renewed: …`).
- Endpoint użycia nie jest publicznie udokumentowany i może się zmienić. Gdy tak się stanie, nadal widzisz zapisane liczby, a `ccam doctor` zgłasza zmianę formatu.
- Gdy nie ma zapisanych liczb ani odpowiedzi, wiersz konta pokazuje `lookup failed:` i przyczynę, na przykład `no saved usage` z `--source cache`.
- Pliki pamięci podręcznej zawierają tylko procenty i czasy resetu, nigdy tokeny.

## Powiązane

- [Limity użycia i przełączanie kont](/docs/usage-limits)
- [Przełączanie na następne konto](/features/next)
- [Co zrobić, gdy Claude Code osiągnie limit użycia](/guides/usage-limit-workflow)
- [Opis poleceń](/docs/commands)
