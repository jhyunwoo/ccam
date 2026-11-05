---
title: Co zrobić, gdy Claude Code osiągnie limit użycia
description: Poczekaj na reset albo od razu kontynuuj tę samą rozmowę na innym koncie Claude Code. Jak działają limity 5-godzinne i tygodniowe, krok po kroku.
updated: 2026-11-05
order: 3
translation: machine
sourceHash: "53b1e3634ea9"
---

Gdy konto osiągnie limit użycia (często nazywany rate limitem), masz dwie możliwości: poczekać, aż limit się zresetuje, albo od razu kontynuować tę samą rozmowę na innym swoim koncie. Z ccam to drugie wymaga jednego polecenia, `ccam next --auto --continue`.

## Jakie limity ma Claude Code?

W planach Pro i Max Claude Code ma [5-godzinny limit sesji i limit tygodniowy](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Oba są wspólne z Claude w przeglądarce i w aplikacjach, a użycie liczy się do obu jednocześnie. Gdy jeden się wyczerpie, Claude Code [wyświetla komunikat](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) w rodzaju „You've hit your session limit · resets 3:45pm”. Zmiana modelu nie pomaga, bo limity obejmują wszystkie modele.

## Opcja 1: poczekaj na reset

Jeśli masz jedno konto, poczekaj. Nowsze wersje Claude Code [czekają w otwartej sesji](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) i same wznawiają zadanie po zresetowaniu limitu. `/usage` w Claude Code pokazuje, jaką część limitów twojego planu już wykorzystałeś.

## Opcja 2: kontynuuj od razu na innym koncie

Jeśli masz drugie konto, możesz od razu kontynuować tę samą rozmowę. Skonfiguruj to raz:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Pasek stanu pokazuje teraz podczas pracy konto i jego 5-godzinne użycie, więc limit nie będzie zaskoczeniem.

Gdy limit zostanie osiągnięty, wyjdź z Claude Code (`/exit`), a następnie:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Otwórz rozmowę na koncie z największym zapasem:

```sh
ccam next --auto --continue
```

Claude Code uruchomi się jako `work` z najnowszą rozmową z bieżącego folderu, łącznie ze wszystkim, co zostało w niej dotąd powiedziane. Współdzielona historia (`ccam group set default --add history` powyżej) sprawia, że oba konta pracują nad tą samą rozmową; bez niej Claude Code kopiuje rozmowę do `work`, a ccam podaje, ile tokenów pierwsza odpowiedź będzie musiała przeczytać ponownie.

Z [pluginem ccam](/docs/install#claude-code-plugin) w Claude Code sesja zatrzymana na limicie pokazuje powiadomienie i zapisuje pełne okno, więc `ccam next --auto` nie wybierze tego konta ponownie przed resetem. `/ccam:switch` (lub `/ccam:switch work`) zapisuje rozmowę dla `ccam next --continue`, więc kontynuuje on dokładnie tę rozmowę.

**Miejsca w planach Team i Enterprise.** Automatyczny wybór konta (`ccam next`, `--fallback`, `ccam pick`, klawisz `n` w panelu oraz, gdy go włączysz, [`auto_switch`](/features/next#auto-switch)) przełącza też między miejscami w tej samej organizacji Claude Team lub Enterprise, co może naruszać warunki Anthropic. Przeczytaj [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) i skonsultuj to ze swoją organizacją i z Anthropic, zanim zaczniesz na nim polegać.

## Czy ccam przełącza konta automatycznie?

Tylko jeśli to włączysz.

- Domyślnie ccam nie działa w trakcie twojej pracy: uruchamia Claude Code i zastępuje się nim, więc to ty decydujesz, kiedy przełączyć. `ccam next --auto` wybiera konto według użycia, biorąc spośród kont w granicach 3 punktów procentowych od najniższego użycia 5-godzinnego to z najniższym użyciem tygodniowym; pomija konta, które mają 100% albo nie są zalogowane, i pozostaje w obrębie jednej organizacji Team lub Enterprise albo wśród planów prywatnych.
- Z `relay = true` w `config.toml` ccam zamiast tego działa dalej obok sesji, więc `/ccam:switch` sam kontynuuje rozmowę, zachowując opcje sesji, tryb uprawnień, poziom wysiłku i agenta; zanim przeniesie rozmowę na konto z innej organizacji, ccam cię o to pyta.
- Z `auto_switch = true` (Linux i macOS; funkcja eksperymentalna) ccam także sam wykonuje przełączenie: gdy Claude Code zatrzyma sesję na jej limicie, a inne konto tego samego rodzaju ma jeszcze zapas, ccam kończy sesję i kontynuuje rozmowę na tamtym koncie ([warunki i koszty](/features/next#auto-switch)).

Niektóre inne narzędzia zamiast tego podmieniają login działającej sesji; [ccam i inne narzędzia do przełączania kont](/compare/account-switchers) wyjaśnia różnicę.

## O czym pamiętać

- Uruchamiaj `ccam next --continue` w tym samym folderze, w którym pracowałeś; rozmowy są przechowywane osobno dla każdego folderu projektu.
- Nie uruchamiaj dwóch kont na tej samej rozmowie jednocześnie. ccam ostrzeże cię, jeśli spróbujesz.
- Bez współdzielonej historii rozmowa jest kopiowana na drugie konto i obie kopie toczą się dalej osobno; punkty kontrolne, z których `/rewind` przywraca kod, zostają na pierwszym koncie.
- Pierwsza wiadomość na nowym koncie może zużyć więcej jego limitu niż zwykle, bo to konto nie ma jeszcze rozmowy w pamięci podręcznej. `ccam next --continue --brief` zaczyna ją zamiast tego od krótkiego podsumowania.
- To, czy wolno ci w ten sposób przełączać się między kontami, zależy od warunków, które obowiązują każde konto.

Wszystkie opcje znajdziesz w częściach [Limity użycia](/features/usage) i [Przełączanie na następne konto](/features/next).
