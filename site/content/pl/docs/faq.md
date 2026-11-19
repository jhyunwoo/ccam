---
title: Najczęściej zadawane pytania
metaTitle: "FAQ ccam: kilka kont Claude Code"
description: Odpowiedzi o ccam, kontach Claude Code, limitach użycia, współdzieleniu, prywatności i obsługiwanych systemach.
updated: 2026-11-19
order: 6
translation: machine
sourceHash: "090ecc3a34bd"
---

## Czy ccam to oficjalne narzędzie Anthropic?

Nie. ccam jest nieoficjalnym projektem o otwartym kodzie i nie jest powiązany z Anthropic ani przez Anthropic popierany. Korzysta z ustawienia `CLAUDE_CONFIG_DIR`, które [Claude Code dokumentuje](https://code.claude.com/docs/en/env-vars) do uruchamiania kilku kont obok siebie.

## Czy wolno mi używać kilku kont Claude Code?

To zależy od warunków każdego konta, a w przypadku kont służbowych także od zasad twojej organizacji: dla Pro i Max od [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) Anthropic, dla Team i Enterprise od jego [Commercial Terms](https://www.anthropic.com/legal/commercial-terms), a dla wszystkich od jego [Usage Policy](https://www.anthropic.com/legal/aup). Sprawdź je, zanim zaczniesz używać kilku kont. ccam nie podnosi, nie łączy ani nie ukrywa limitów konta, a każde konto nadal zatrzymuje się na swoich. Przenosi twoją pracę na inne twoje konto, gdy o to poprosisz, albo sam na limicie, gdy włączysz [`auto_switch`](/features/next#auto-switch) (Linux i macOS). Zanim z tego skorzystasz, sprawdź, czy jest to dozwolone dla twoich kont; decyzja należy do ciebie.

**Miejsca w planach Team i Enterprise.** Automatyczny wybór konta (`ccam next`, `--fallback`, `ccam pick`, [`ccam agent`](/features/agent), klawisz `n` w panelu oraz [`auto_switch`](/features/next#auto-switch), gdy go włączysz) przełącza też między miejscami w tej samej organizacji Claude Team lub Enterprise. Centrum pomocy Anthropic podaje: „Each seat on a Team plan is assigned to one named individual, and only that person may use it” (każde miejsce w planie Team jest przypisane do jednej imiennie wskazanej osoby i tylko ta osoba może z niego korzystać) oraz „Users can't route their work through another user's account to use that user's remaining usage” (użytkownicy nie mogą kierować swojej pracy przez konto innego użytkownika, aby wykorzystać pozostały limit użycia tego użytkownika) ([plany Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); podobnie w przypadku [planów Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Dlatego korzystanie z automatycznego wyboru między miejscami może naruszać warunki Anthropic. Czy korzystać z niego w tej sytuacji, decydujesz ty: zanim zaczniesz na nim polegać, skonsultuj to ze swoją organizacją i z Anthropic. Planów prywatnych (Pro, Max) dotyczy akapit powyżej.

## Czy ccam pozwala obejść limity użycia?

Nie ponad własny limit żadnego konta: każde konto ma swoje limity i zatrzymuje się na nich, a ccam ich nie podnosi, nie łączy ani nie ukrywa. ccam je pokazuje i pozwala kontynuować pracę na innym twoim koncie, gdy na jednym limit się wyczerpie; z włączonym `auto_switch` robi to sam (zobacz następne pytanie). Między miejscami jednej organizacji Team lub Enterprise kontynuowanie pracy na innym miejscu zużywa pozostały limit tego miejsca, a centrum pomocy Anthropic wyklucza to w przypadku miejsca innego użytkownika (zobacz Miejsca w planach Team i Enterprise powyżej).

## Czy ccam może sam przełączyć konto, gdy sesja osiągnie limit?

Tylko jeśli to włączysz. Z `auto_switch = true` na górze `~/.config/ccam/config.toml` (Linux i macOS; domyślnie wyłączone i eksperymentalne), gdy Claude Code zatrzyma sesję interaktywną (uruchomioną przez ccam) na limicie 5-godzinnym, tygodniowym, Opus, Sonnet lub Fable, ccam kończy sesję i kontynuuje rozmowę jako inne twoje konto, które ma jeszcze zapas: inne miejsce w tej samej organizacji Team lub Enterprise albo inny plan prywatny (Pro, Max). Nigdy nie przełącza przed osiągnięciem limitu, nigdy nie śledzi twojego użycia w tle i nigdy nie przenosi pracy poza organizację Team lub Enterprise ani do niej.

W przypadku miejsc przenosi twoją pracę z jednego miejsca na inne, co może naruszać warunki Anthropic (zobacz Miejsca w planach Team i Enterprise powyżej). W przypadku planów prywatnych warunki Anthropic nie odnoszą się do tego jednoznacznie: Consumer Terms zabraniają korzystania z usług „through automated or non-human means” (w sposób zautomatyzowany lub inny niż przez człowieka), chyba że Anthropic wyraźnie na to pozwoli, a przeniesienie wysyła jedną wiadomość w twoim imieniu. Czy to włączyć, decydujesz ty.

[`auto_switch`](/features/next#auto-switch) wymienia wszystkie warunki i to, ile kosztuje przeniesienie.

## Jak używać dwóch kont Claude Code na jednym komputerze?

Daj każdemu kontu własny folder konfiguracji przez `CLAUDE_CONFIG_DIR`. ccam robi to za ciebie: uruchom `ccam add` dla każdego konta, a potem uruchamiaj je przez `ccam personal` i `ccam work` albo przez skróty o nazwach wybranych przez ciebie, na przykład `ccwork`. [Poradnik](/guides/multiple-claude-code-accounts) pokazuje też, jak zrobić to ręcznie.

## Jak kontynuować rozmowę na innym koncie?

Uruchom `ccam next --continue` w danym folderze: przenosi ostatnią rozmowę z tego folderu na następne konto (`--auto` wybiera to z największym zapasem). W Claude Code `/ccam:switch` z [pluginu ccam](/docs/install#claude-code-plugin) robi to samo dla sesji, w której jesteś. Gdy oba konta współdzielą historię rozmów (`ccam group set default --add history` albo osobna grupa tylko dla nich: `ccam group add pair --items history --members work,side`), nowe konto otwiera dokładnie tę rozmowę; w przeciwnym razie Claude Code kopiuje ją do historii nowego konta. Dodaj `--brief`, żeby nowe konto zaczęło od krótkiego podsumowania zamiast całej rozmowy.

## Czy ccam widzi albo wysyła mój login?

ccam nigdy nie współdzieli ani nie kopiuje twojego loginu Claude między kontami i nigdy nie kopiuje logowań serwerów MCP. Odczytuje token dostępu konta tylko po to, żeby zapytać endpoint użycia Anthropic o użycie tego konta: dla `ccam usage`, panelu, `ccam doctor` i automatycznego wyboru (`ccam next --auto`, `ccam pick`, `--fallback` oraz `auto_switch`, gdy sesja osiągnie limit). Wygasły token odnawia tylko przy `ccam usage --refresh`; z `usage_source = "claude"` nigdzie nie wysyła tokena, a zapisany login (plik albo wpis w pęku kluczy macOS) odczytuje tylko po to, żeby sprawdzić, czy nadal działa. `ccam ext` i `ccam doctor` odczytują nazwy logowań MCP konta, nigdy ich tokeny. ccam nie ma telemetrii. Poza tymi zapytaniami jedyne żądanie, które wysyła sam, to sprawdzenie przez panel w GitHub, czy jest nowsze wydanie ccam, które nie zawiera niczego o tobie ani o twoich kontach.

## Jak zalogować się na zdalnym serwerze?

Uruchom `ccam add` lub `ccam login <name>` w sesji SSH jak zwykle. Claude Code czeka na logowanie na serwerze, więc przeglądarka na twoim komputerze kończy na stronie `http://localhost:<port>/callback?code=…`, która się nie wczytuje. Skopiuj ten adres z przeglądarki i wklej go do terminala: ccam przekaże go do Claude Code na serwerze, który dokończy logowanie. Adres zawiera jednorazowy kod, więc nikomu go nie udostępniaj.

## Gdzie są przechowywane moje konta?

W `~/.local/share/ccam/profiles/<name>/`, jeden folder na konto. Na macOS Claude Code przechowuje sam login w pęku kluczy, pod nazwą utworzoną na podstawie tego folderu, więc `ccam rename` zachowuje tam starą nazwę folderu.

## Co się stanie z moim istniejącym ~/.claude?

Nic, chyba że uruchomisz `ccam import` lub `ccam migrate`. Import rejestruje folder tam, gdzie jest, więc zwykłe `claude` i `ccam <name>` nadal używają tego samego konta i ustawień. `ccam migrate` kopiuje jego ustawienia, `CLAUDE.md`, rozszerzenia i serwery MCP do grupy `default`, a folder tylko odczytuje.

## Jakie systemy obsługuje ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine i inne), macOS oraz Windows 10 i 11, natywnie lub w WSL. W Windows współdzielenie plików takich jak `settings.json` wymaga Trybu dewelopera; współdzielone foldery działają bez niego, a bez niego pierwsze `ccam add` tworzy grupę `default` tylko z serwerami MCP i rozszerzeniami ([współdzielenie w Windows](/features/share#on-windows)).

## Dlaczego serwer MCP pojawia się tylko na jednym koncie?

Serwery dodane w trakcie sesji przez `claude mcp add` są zapisywane tylko dla tego konta, a bez `-s user` tylko dla bieżącego folderu. ccam poinformuje cię o tym przy następnym uruchomieniu konta, a `ccam mcp list` i `ccam ext` je pokazują; przenieś serwer do współdzielonej grupy poleceniem `ccam mcp import <account> <server>`.

## Czy ccam spowalnia Claude Code?

Uruchamianie Claude Code przez ccam dodaje około 4 milisekund, co zmierzono skryptem `scripts/perf-check.sh` na laptopie. ccam odczytuje jeden mały plik, sprawdza współdzielone linki, a potem zastępuje się procesem `claude`, więc domyślnie nie działa w czasie twojej pracy. Windows nie potrafi zastąpić procesu, więc tam ccam uruchamia Claude Code i czeka na jego zakończenie, nie robiąc nic więcej. Z włączonym `relay` lub `auto_switch` ccam działa dalej w ten sam sposób obok każdej sesji interaktywnej, żeby móc kontynuować rozmowę na innym koncie.

## Jak zaktualizować ccam?

Uruchom `ccam update`. To polecenie pobiera najnowsze wydanie, sprawdza jego sumę kontrolną i podmienia plik `ccam`, a `ccam update --check` tylko informuje, czy jest dostępne nowsze wydanie. Panel także informuje, gdy jest nowsze wydanie, a `U` instaluje je po zapytaniu; ccam nigdy nie aktualizuje się sam. Jeśli ccam pochodzi z pakietu `.deb`, zainstaluj zamiast tego nowy `.deb` przez apt. Szczegóły znajdziesz w części [aktualizacja](/docs/install#update).

## Jak odinstalować ccam?

Uruchom `ccam uninstall`. Wypisuje, co zostanie usunięte, i najpierw pyta o zgodę, a potem usuwa każde konto (najpierw je wylogowując), skróty, `~/.config/ccam`, `~/.local/share/ccam`, `~/.local/state/ccam` i polecenie `ccam`. Zaimportowany `~/.claude` zostaje, z kopią tego, co współdzielił w swoich grupach. Zobacz [instalacja](/docs/install#uninstall).
