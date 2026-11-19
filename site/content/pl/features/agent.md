---
title: Subagenci na innych kontach
description: ccam agent uruchamia część zadania Claude Code jako inne twoje konto i zwraca odpowiedź Claude, więc długie zadanie korzysta z limitów użycia kilku kont.
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Własni subagenci Claude Code działają jako konto sesji, która ich uruchamia, więc zużywają limit użycia tego konta. `ccam agent "<task>"` uruchamia zadanie przez `claude -p` w tym samym folderze jako inne z twoich kont i wypisuje ostatnią odpowiedź Claude. Wywołany z sesji Claude Code działa jak subagent, który zużywa limit użycia innego konta, więc jedno długie zadanie może trwać dalej, gdy własne konto sesji już by się wyczerpało.

**Miejsca w planach Team i Enterprise.** `ccam agent` sam wybiera konto. W obrębie jednej organizacji Claude Team lub Enterprise uruchamia zadanie na innym miejscu, co może naruszać warunki Anthropic. Przeczytaj [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) i skonsultuj to ze swoją organizacją i z Anthropic, zanim zaczniesz na nim polegać.

## Użycie

Claude Code uruchamia go swoim narzędziem Bash. Gdy masz zainstalowany [plugin ccam](/docs/install#claude-code-plugin), poproś o to zwykłymi słowami, na przykład „użyj moich innych kont do części tego zadania”, albo zapisz to w swoim CLAUDE.md. Skill pluginu każe Claude Code uruchamiać każdego agenta w tle, pisać zadanie zrozumiałe bez dodatkowego kontekstu i dawać agentom, którzy edytują pliki, własną część albo własny git worktree.

```sh
ccam agent "find why TestParse fails and fix it; reply with what you changed"
ccam agent --model opus -- "-v breaks the parser: find out why"
ccam agent --json < task.md            # the task from standard input
ccam agent --resume 7f3c9a2e-… "now add a test for it"
ccam agent --wait=1h "…"               # nobody has room: wait for a reset
ccam --dry-run agent "…"               # the account and command it would run
```

```console
$ ccam agent "summarize what internal/usage does"
ccam: agent runs as work (12% of its 5 hours used); session 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f
internal/usage looks up the 5-hour and weekly usage of every account …
ccam: agent finished as work; follow up with: ccam agent --resume 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f "…"
```

Ostatnia odpowiedź Claude trafia na standardowe wyjście; wiersze zaczynające się od `ccam:` trafiają na standardowe wyjście błędów. `--json` wypisuje zamiast tego `{account, session_id, result, is_error, accounts}`, gdzie `accounts` wymienia każde konto, jako które działało zadanie.

Kody wyjścia: własny kod claude (`0`, gdy agent skończył), `75` (`EX_TEMPFAIL`), gdy żadne konto poza wywołującym nie ma zapasu, `1` przy innych błędach, takich jak próba uruchomienia innego agenta przez agenta.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--model <model>` | Model agenta; przy wyborze liczy się też tygodniowe okno tego modelu | domyślny model konta |
| `--resume <id>` | Kontynuuje rozmowę wcześniejszego agenta z kolejnym zadaniem (identyfikator sesji z jego pierwszego i ostatniego wiersza) | wyłączone |
| `--json` | Wypisuje konto, sesję, odpowiedź i konta, jako które działał | wyłączone |
| `--max-usage <1-100>` | Preferuje konta poniżej tego użycia | `auto_max_usage`, w przeciwnym razie 100 |
| `--strategy headroom\|expiring\|balance` | [Zasada wyboru](/features/next#auto) | `auto_strategy`, w przeciwnym razie `headroom` |
| `--cross-org` | Działa też jako konta innych organizacji | `cross_org`, w przeciwnym razie wyłączone |
| `--wait[=<duration>]` | Gdy żadne konto nie ma zapasu, czeka na reset pierwszego, najwyżej tyle | wyłączone; `5h`, gdy podano samą opcję |
| `--dry-run` | Wypisuje polecenie dla konta, które by wybrał, bez uruchamiania go | wyłączone |

## Jak to działa

1. **Wywołujący jest pomijany.** Konto, z którego wywołano `ccam agent` (konto sesji uruchomionej przez ccam, w przeciwnym razie to, które uruchamia ten folder), nigdy nie uruchamia agenta; nie robi tego też inny folder zalogowany na to samo konto Claude, bo dzieli z nim limity. Agent zostaje w [organizacji](/features/next#auto) wywołującego, chyba że `--cross-org` mówi inaczej.
2. **Wybiera jak `ccam pick`.** Użycie jest sprawdzane tak, jak robi to [`ccam pick`](/features/pick), a agenci uruchomieni w odstępie do 2 minut dostają różne konta, dopóki wystarczająco wiele ma zapas.
3. **Działa z uprawnieniami sesji, nigdy większymi.** W `claude -p` nikt nie może zezwolić na użycie narzędzia, więc agent dostaje tryb uprawnień i poziom wysiłku sesji, która go wywołała. Gdy uruchomisz go z terminala, obowiązuje domyślny tryb Claude Code i reguły uprawnień z twoich ustawień. `ccam agent` nie ma opcji, która dawałaby więcej.
4. **Przy limicie użycia rozmowa toczy się dalej.** Gdy konto agenta osiągnie limit, ccam zapisuje limit tak, jak zrobiłby to [hook pluginu](/features/usage#the-limit-record), i kontynuuje tę samą rozmowę jako konto z największym zapasem: Claude Code kopiuje ją tam (`--resume <file> --fork-session`) z pierwszą wiadomością „Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.” Każde konto jest próbowane raz.
5. **Każda rozmowa agenta jest zapisywana** w `~/.local/state/ccam/agent/<session id>.json` (30 dni). `--resume <id>` kontynuuje ją w jej folderze, jako konto, które ją prowadziło, jeśli nadal ma zapas, a w przeciwnym razie jako inne. [`ccam next --continue`](/features/next) pomija rozmowy agentów, więc nadal kontynuuje twoje.

## Bezpieczeństwo i przypadki brzegowe

- Agent nie może uruchamiać własnych agentów: ccam oznacza go przez `CCAM_AGENT=1`, a `ccam agent` odmawia tam działania.
- Zatrzymanie `ccam agent` (Ctrl+C, sygnał albo przerwanie polecenia przez Claude Code) zatrzymuje też agenta.
- Agent nie widzi rozmowy, która go wywołała, tylko swoje zadanie.
- Agenci, którzy jednocześnie edytują te same pliki, nadpisują sobie nawzajem pracę. Daj każdemu własną część zadania albo własny git worktree.
- Gdy sesja wywołująca przenosi się na inne konto ([`auto_switch`](/features/next#auto-switch)), uruchomieni przez nią agenci zatrzymują się razem z nią. Ich rozmowy zostają: `--resume <id>` je kontynuuje.
- `agent` jest nazwą polecenia, więc żadne nowe konto nie może nazywać się `agent`; istniejące nadal uruchomisz przez `ccam run agent`.

## Powiązane

- [Wybór konta dla skryptów](/features/pick)
- [Przełączanie na następne konto](/features/next)
- [Limity użycia](/features/usage)
- [Sesje w tle wszystkich kont](/features/agents)
