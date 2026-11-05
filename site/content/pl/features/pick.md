---
title: Wybór konta dla skryptów
description: ccam pick wypisuje konto z największym zapasem dla skryptu, daje workerom uruchomionym razem różne konta i kończy się kodem 75, gdy żadne nie ma zapasu.
updated: 2026-11-05
order: 12
translation: machine
sourceHash: "0ec25752ace9"
---

`ccam pick` wypisuje nazwę zalogowanego konta z największym zapasem, żeby skrypt mógł na nim uruchomić Claude Code. Workery uruchomione w tym samym czasie dostają każdy inne konto, dopóki jest wystarczająco wiele kont z zapasem, więc partia zadań `claude -p` rozkłada się na twoje konta, zamiast wyczerpywać jedno.

**Miejsca w planach Team i Enterprise.** Automatyczny wybór konta (`ccam next`, `--fallback`, `ccam pick`, klawisz `n` w panelu oraz, gdy go włączysz, [`auto_switch`](/features/next#auto-switch)) przełącza też między miejscami w tej samej organizacji Claude Team lub Enterprise, co może naruszać warunki Anthropic. Przeczytaj [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) i skonsultuj to ze swoją organizacją i z Anthropic, zanim zaczniesz na nim polegać.

## Użycie

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # most weekly quota left per day; fewer sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

Kody wyjścia: `0` z nazwą na standardowym wyjściu, `75` (`EX_TEMPFAIL`), gdy żadne konto nie ma zapasu lub żadne nie jest zalogowane, `1` przy innych błędach, takich jak nieprawidłowa opcja.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--json` | Wypisuje konto, jego użycie 5-godzinne i tygodniowe, źródło liczb i powód wyboru | wyłączone |
| `--model <model>` | Uwzględnia też tygodniowe okno tego modelu, np. `opus` | brak |
| `--max-usage <1-100>` | Preferuje konta poniżej tego użycia | `auto_max_usage`, w przeciwnym razie 100 |
| `--strategy headroom\|expiring\|balance` | Najniższe użycie tygodniowe w granicach 3 punktów procentowych od najniższego użycia 5-godzinnego, najbliższy reset okna 5-godzinnego albo największy tygodniowy limit pozostały na dzień | `auto_strategy`, w przeciwnym razie `headroom` |
| `--cross-org` | Wybiera też konta innych organizacji | `cross_org`, w przeciwnym razie wyłączone |
| `--wait[=<duration>]` | Gdy żadne konto nie ma zapasu, czeka na reset pierwszego, najwyżej tyle. Wartość podaje się po `=`: `--wait 30m` jest odrzucane z podpowiedzią `--wait takes its value after '=': --wait=30m` | wyłączone; `5h`, gdy podano samą opcję |
| `--dry-run` | Wybiera bez zapisywania wyboru | wyłączone |

## Jak to działa

1. ccam wyświetla listę zalogowanych kont (konta z kluczem API i loginy Console nie mają limitów i są pomijane) w organizacji konta, które uruchomiłby ten folder (zobacz [organizacje](/features/next#auto)); folder zalogowany na to samo konto Claude co to konto również jest pomijany.
2. Zakłada blokadę (`~/.local/state/ccam/picks.lock`), więc wybory uruchomione w tym samym czasie odbywają się jeden po drugim.
3. Sprawdza użycie wszystkich kont naraz, tak samo jak `ccam usage`, z 60-sekundową pamięcią podręczną. Okno, które Claude Code zgłosił jako pełne, liczy się jako 100% aż do resetu ([zapis limitu](/features/usage#the-limit-record)), więc konto, które właśnie osiągnęło limit, jest wykluczane, nawet jeśli zapisane liczby nadal pokazują zapas.
4. Konta wybrane w ciągu ostatnich 2 minut, przez `ccam pick` albo przez [`auto_switch`](/features/next#auto-switch), który przeniósł tam sesję (`~/.local/state/ccam/picks/<name>`; wybór z przyszłą datą się nie liczy), są najpierw pomijane, a [zasady wyboru](/features/next#auto) z `ccam next --auto` wybierają spośród pozostałych: najpierw poniżej `auto_max_usage`, według strategii, nigdy konto na poziomie 100%. Jeśli przez to pominięte zostaje konto, które zasady wybrałyby w przeciwnym razie, mówi o tym `reason` w `--json`: `30% of its 5 hours used; work was picked a moment ago`.
5. Jeśli każde konto z zapasem zostało niedawno wybrane, zasady wybierają ponownie spośród wszystkich, więc przy trzech kontach czwarty worker dzieli jedno z nich.
6. Wybór zostaje zapisany, blokada zwolniona, a potem wypisywana jest nazwa.

Okno 2 minut ma znaczenie, bo użycie nowego zadania pojawia się w liczbach z pewnym opóźnieniem: bez niego pięć workerów uruchomionych razem widziałoby to samo konto jako to z największym zapasem.

## Bezpieczeństwo i przypadki brzegowe

- `ccam pick` tylko odczytuje użycie i zapisuje czas wyboru; nigdy nie odnawia ani nie odczytuje tokenów w większym zakresie niż `ccam usage`.
- Wybór nie rezerwuje konta. Dwa skrypty mogą nadal wystartować na tym samym koncie po 2 minutach; użyj `--fallback`, żeby przejść dalej, gdy to konto osiągnie limit.
- Czekanie z `--wait` nie trzyma blokady, więc w tym czasie inne wybory działają dalej.
- `ccam pick` ignoruje [`ccam folder`](/features/folder): wybiera wyłącznie według użycia.

## Powiązane

- [Przełączanie na następne konto](/features/next)
- [Ponowienie po osiągnięciu limitu](/features/next#run-again-at-the-limit)
- [Limity użycia](/features/usage)
- [Konta z kluczem API](/features/api)
