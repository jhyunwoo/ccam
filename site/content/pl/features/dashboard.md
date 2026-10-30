---
title: Panel
description: ccam w terminalu otwiera pełnoekranowy panel kont, grup, reguł folderów, rozszerzeń i sprawdzeń; każda zmiana najpierw pokazuje swój plan.
updated: 2026-10-30
order: 15
translation: machine
sourceHash: "d65f1b5bb0c8"
---

`ccam` bez argumentów otwiera pełnoekranowy panel, gdy standardowe wejście i wyjście są terminalem; w przeciwnym razie wypisuje `ccam list`. Panel pokazuje każde konto z jego użyciem i grupami i robi to samo co polecenia: uruchamia, dodaje, loguje, zmienia grupy, ustawia reguły folderów, wyświetla sesje w tle, sprawdza i naprawia. Uruchomienie konta albo podłączenie się do sesji w tle go zamyka; dodawanie konta lub logowanie na chwilę oddaje terminal, a potem wraca do panelu.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · r to refresh · a to add · ? for shortcuts
```

Wiersz pod nazwami kart mówi jednym zdaniem, do czego służy bieżąca karta. W szerszym terminalu karta Accounts pokazuje też adres e-mail każdego konta. Na każdej karcie koniec wiersza z nazwami kart pokazuje, czy [`auto_switch`](/features/next#auto-switch) jest włączone (Linux i macOS; w systemie Windows jest ignorowane). Gdy jest włączone, karta Health wymienia, czego mu jeszcze brakuje: Claude Code 2.1.214 lub nowszego oraz włączonego pluginu ccam na każdym koncie, które loguje się subskrypcją Claude.

Przy każdym otwarciu panel pyta w tle GitHub, czy istnieje nowsze wydanie ccam; pierwsza klatka nigdy na to nie czeka, a gdy jesteś offline, nic się nie pojawia. Gdy takie wydanie jest, wiersz pod tytułem je poleca:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` najpierw pyta, a potem aktualizuje tak jak [`ccam update`](/features/update): wydanie jest sprawdzane z jego sumami kontrolnymi SHA-256 i zastępuje ccam w jednym kroku. Otwarty panel to wciąż stara wersja, więc wyjdź i uruchom `ccam` ponownie. Gdy aktualizacji nie da się wykonać tutaj, na przykład dla kopii ccam zainstalowanej z pakietu `.deb`, panel pokazuje tę samą wskazówkę co `ccam update` i zachowuje propozycję.

Pierwsza klatka jest rysowana wyłącznie z `config.toml`, w znacznie mniej niż 100 ms; użycie pojawia się, gdy tylko dane konto odpowie. Każdy limit ma własny pasek: okno 5-godzinne, tygodniowe okno dla wszystkich modeli i własne tygodniowe okno Fable; `-` oznacza, że konto nie ma takiego okna. Procent zmienia kolor na żółty przy 80% i na czerwony przy 100%. Jeśli zapytanie się nie powiedzie, wartości z pamięci podręcznej są pokazywane przyciemnione. Okno, które Claude Code zgłosił jako pełne, pokazuje 100% aż do resetu. Wygasły token dostępu pokazuje blade `⚠ token expired` (Claude Code odnawia go przy następnym uruchomieniu), a `⚠ login expired` oznacza, że konto wymaga `ccam login`.

## Karty

| Karta | Co pokazuje | Klawisze |
| --- | --- | --- |
| Accounts | numer, nazwa, adres e-mail, paski 5-godzinny, tygodniowy i Fable, ostatnie użycie, grupy; `✔` konto domyślne, `▸` konto tego folderu, `•` uruchomione (także Claude Code uruchomiony bez ccam: zwykłe `claude`, IDE, sesje w tle i ich demon) | poniżej |
| Agents | to, co wyświetla [`ccam agents`](/features/agents): sesje w tle i otwarte sesje każdego konta, z ich stanem, folderem i czasem startu (poniżej) | `enter` podłączenie (jak `ccam attach`), `x` zatrzymanie, `r` ponowne wczytanie |
| Groups | każda grupa, co współdzieli i jej konta, oraz konta, które nie należą do żadnej grupy | `enter` edycja, `n` nowa, `e` zmiana nazwy, `x` usunięcie |
| Folders | reguły [`ccam folder`](/features/folder) i ta, która obowiązuje tutaj, nawet gdy trafiłeś do folderu przez link symboliczny (nowe reguły są zapisywane z rozwiązanymi linkami) | `n` użyj tu konta, `enter` zmiana, `x` usunięcie |
| Extensions | to, co pokazuje [`ccam ext`](/features/ext), każde rozszerzenie raz, z grupami, z których pochodzi, `off:` konta, na których wtyczka nie jest włączona, własne serwery konta, które dotyczą tylko tego folderu, `This folder's project` oraz problemy, które nie należą do żadnego konta | `r` ponowne wczytanie |
| Health | to, co zgłasza [`ccam doctor`](/features/doctor), jako ramka dla każdego konta (poniżej) | `f` naprawa zaznaczonej ramki, `F` naprawa wszystkiego (`doctor --fix`), `↑`/`↓` przejście między ramkami |

`←`/`→` lub `Tab` przełączają karty, `?` wyświetla klawisze bieżącej karty, `r` odświeża użycie z pominięciem krótkiej pamięci podręcznej i ponownie wczytuje karty Agents, Extensions i Health, a `q` lub `Esc` zamyka panel.

## Klawisze na karcie Accounts

| Klawisz | Działanie |
| --- | --- |
| `Enter` | Uruchamia Claude Code jako zaznaczone konto |
| `1`–`9` | Przechodzi do konta o tym numerze; wpisz szybko `1` `2`, żeby wybrać 12. Uruchamia tylko `Enter` |
| `↑`/`↓`, `k`/`j` | Przesuwa zaznaczenie |
| `n` | Uruchamia konto, które wybrałby [`ccam next --auto`](/features/next) |
| `o` | Uruchamia z opcjami: `--continue`, `--resume`, inne argumenty albo próbne uruchomienie |
| `u` | Szczegóły użycia: okna 5-godzinne, tygodniowe i modeli z ich resetami oraz źródło liczb |
| `a` | Dodaje konto (kreator opisany niżej) |
| `l` | Loguje ponownie |
| `s` | Wybiera grupy konta |
| `d` | Ustawia je jako konto domyślne |
| `f` | Używa go w bieżącym folderze |
| `e` | Zmienia jego nazwę |
| `c` | Ustawia lub usuwa jego polecenie skrótu |
| `x` | Usuwa je (najpierw pyta, z zaznaczonym No) |

**Miejsca w planach Team i Enterprise.** Automatyczny wybór konta (`ccam next`, `--fallback`, `ccam pick`, klawisz `n` w panelu oraz, gdy go włączysz, [`auto_switch`](/features/next#auto-switch)) przełącza też między miejscami w tej samej organizacji Claude Team lub Enterprise, co może naruszać warunki Anthropic. Przeczytaj [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) i skonsultuj to ze swoją organizacją i z Anthropic, zanim zaczniesz na nim polegać.

## Karta Health

Karta Health po otwarciu wykonuje sprawdzenia [`ccam doctor`](/features/doctor) i rysuje ten sam widok: pierwszy wiersz ze znacznikiem dla każdego konta, a pod nim ramkę dla każdego konta. Każdy wiersz w ramce to jeden sprawdzony obszar, najpierw `Folder`, `Login` i `Share links`, oznaczony `✔` w porządku, `ℹ` informacja, `⚠` ostrzeżenie lub `✘` błąd, więc od razu widać, które konto ma problem i gdzie. Ustalenia, które nie należą do żadnego konta, trafiają do ramki `General`.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` przechodzą od ramki do ramki, a zaznaczona ramka ma obramowanie w kolorze akcentu. `f` naprawia tylko zaznaczoną ramkę, tak jak zrobiłby to `ccam doctor --fix`: w ramce konta naprawia tylko to konto, a w ramce `General` ustalenia, które nie należą do żadnego konta (ścieżki wtyczek i skróty); pozostałe konta są tylko sprawdzane ponownie. Najpierw pyta (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), a przy ramce, w której nie ma nic do naprawy, wyświetla tylko `nothing to repair in <name>`. `F` naprawia wszystko naraz, jak `ccam doctor --fix`, a `r` sprawdza ponownie.

## Karta Agents

Karta Agents pokazuje to, co wypisuje [`ccam agents`](/features/agents): sesje w tle i otwarte sesje każdego konta, razem z kontem, do którego każda należy. Pyta Claude Code wszystkich kont naraz (`claude agents --json`, każde z limitem czasu 15 sekund), gdy pierwszy raz otworzysz kartę, po `r`, po zatrzymaniu sesji i po zmianie konfiguracji, a do czasu, aż wszystkie odpowiedzą, wyświetla `Asking every account's Claude Code for its sessions…`.

```console
  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Background agents and open sessions of every account; enter attaches to one.

    Account   ID        Name   State          Folder      Started
 ❯  work      -         -      open (busy)    ~/code/web  now
    personal  49ec0028  probe  blocked, idle  ~/code/app  12m

 ✘ lab: could not list its sessions: exit status 1 'claude agents --json' is
   not available in this environment.

 enter to attach · x to stop · r to reload · ? for shortcuts
```

`State` to stan sesji w tle i to, co robi jej proces, albo `open` dla sesji interaktywnej, tak jak w `ccam agents`. Jeśli konta nie da się zapytać, wiersz `✘` pod listą pokazuje, co odpowiedział Claude Code, jak w przypadku `lab` powyżej, które ma wyłączony widok agentów. Zakończone sesje w tle są pomijane (wyświetla je `ccam agents --all`), a gdy wszystkie konta odpowiedziały i żadne nie ma sesji, karta wyświetla ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` W wąskim terminalu karta najpierw pomija `Started`, potem `Name`, potem `Folder`, a długie ścieżki folderów skraca od lewej.

`Enter` na sesji w tle zamyka panel i otwiera tę sesję w tym terminalu jako jej konto, tak jak `ccam attach <id>`; panel już nie wraca (w Claude Code `←` przechodzi do jego widoku agentów, a `Ctrl+Z` do powłoki). Otwarta sesja nie ma identyfikatora, więc `Enter` i `x` tylko informują, że potrzebna jest sesja w tle. `x` najpierw pyta, z zaznaczonym `No, keep it`: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` Po potwierdzeniu ccam uruchamia `claude stop <id>` jako to konto, wyświetla ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` albo powód, dla którego się nie udało, i ponownie wczytuje listę.

## Zmiana grup

Na karcie Groups `Enter` otwiera edytor: elementy po lewej, konta po prawej, `Tab` przełącza stronę, `Space` zaznacza lub odznacza. Element, który któreś z kont dostaje już z innej grupy, jest wyszarzony z podanym powodem, podobnie jak konto, które dostałoby jakiś element dwa razy; element może pochodzić tylko z jednej grupy. `Enter` pokazuje plan:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Jeśli plik różni się od kopii grupy, ccam najpierw pyta, co zrobić (zachować wersję grupy, użyć mojej albo scalić). Potem zmiana wykonuje się w tle tym samym kodem co [`ccam group`](/features/share): najpierw zapisywana jest konfiguracja, dane każdego konta są przenoszone po kolei, a oryginały trafiają do `.ccam-backup/`; konto, które opuszcza współdzielenie MCP, zachowuje serwery grupy w swoim `.claude.json`.

Zmiana, która przeniosłaby dane uruchomionego konta, jest odrzucana z podaniem nazwy tego konta, nawet wtedy, gdy konto uruchomiono, kiedy plan był otwarty, a plan, który stał się nieaktualny, bo inny proces ccam zmienił grupy, jest odrzucany komunikatem `the groups changed while this plan waited; nothing was moved, make the change again`. Gdy inny proces ccam zmienia grupy, nowa zmiana, zmiana nazwy lub usunięcie wyświetla `another ccam process is changing the groups right now`; po zmianie, która zatrzymała się w połowie, czekają, aż `F` na karcie Health (albo `ccam doctor --fix`) ją dokończy.

## Dodawanie konta

`a` otwiera kreator, którego używa też `ccam add` w terminalu:

1. Nazwa konta.
2. Sposób logowania: subskrypcja Claude teraz (w przeglądarce), [klucz API](/features/api) albo subskrypcja Claude później.
3. Dla klucza API: skąd ccam odczytuje klucz przy każdym uruchomieniu (`env:NAME` lub `command:COMMAND`), punkt końcowy, sposób wysyłania klucza i model.
4. Grupy do dołączenia, z zaznaczonym `default` (jeśli nie ma jeszcze żadnych grup: co będzie współdzielić grupa `default`).
5. Polecenie skrótu, jeśli go chcesz.
6. Podsumowanie do potwierdzenia.

`Esc` cofa o jeden krok, a w pierwszym kroku zamyka kreator.

## Usuwanie konta

`x` pyta, zanim cokolwiek się stanie, z zaznaczonym No, i mówi, co zrobi usunięcie: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` Konto z kluczem API nie ma loginu do unieważnienia, a gdy inne konto ma ten sam login, konto nie jest wylogowywane. Zaimportowane konto zachowuje swój folder: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## Kolory

Kolory odpowiadają rolom z terminalowego interfejsu Claude Code, więc ccam dobrze wygląda obok niego; nazwy grup dostają kolor wyliczony z nazwy. ccam pyta terminal o kolor tła i wybiera jasny lub ciemny motyw. Możesz ustawić go sam na górze `config.toml`, nad pierwszym `[[group]]`, `[[profile]]` lub `[[folder]]`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Motywy `-daltonized` są przeznaczone dla osób z zaburzeniem rozróżniania barwy czerwonej i zielonej, a motywy `-ansi` używają 16 kolorów samego terminala. `NO_COLOR` wyłącza kolory.

## Bezpieczeństwo i przypadki brzegowe

- Panel wywołuje te same funkcje co polecenia, więc obowiązuje każda zasada: plany przed przeniesieniem danych, odmowa dla uruchomionych kont, kopie zapasowe i tokeny logowania, które nigdy nie opuszczają konta.
- Panel nie ma trybu próbnego: `ccam --dry-run` bez nazwy wypisuje, jak uruchomiłoby się konto tego folderu, a `o` na koncie proponuje próbne uruchomienie tylko tego jednego startu.
- Terminal jest przywracany, zanim wystartuje Claude Code albo otworzy się logowanie.
- Gdy ekran jest za mały na okno dialogowe, najpierw znikają wiersze karty; wiersze są obcinane znakiem `…` zamiast zawijane. Na niskim ekranie jako następne znikają puste wiersze i opis karty, więc wiersz z podpowiedziami klawiszy pozostaje widoczny. Poniżej 40 kolumn lub 8 wierszy panel pokazuje tylko `Too small: ccam needs 40x8 (now 36x7)`.
- `Ctrl+C` w trakcie stosowania zmiany, usuwania, aktualizacji lub naprawy czeka na jej zakończenie: wiersz stanu mówi `ccam quits when it is done.`
- Po aktualizacji ccam przez `U` propozycja nie wraca, gdy panel wraca z dodawania konta lub logowania.

## Powiązane

- [Współdzielenie przez grupy](/features/share)
- [Uruchamianie Claude Code](/features/launch)
- [Zarządzanie kontami](/features/accounts)
- [Limity użycia](/features/usage)
