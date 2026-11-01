---
title: Was tun, wenn Claude Code sein Nutzungslimit erreicht
description: Auf den Reset warten oder dasselbe Gespräch sofort in einem anderen Claude Code-Konto fortsetzen. So funktionieren das 5-Stunden- und das Wochenlimit.
updated: 2026-11-01
order: 3
translation: machine
sourceHash: "760a36724feb"
---

Wenn ein Konto sein Nutzungslimit erreicht (oft Rate Limit genannt), hast du zwei Möglichkeiten: warten, bis das Limit zurückgesetzt wird, oder dasselbe Gespräch sofort in einem anderen Konto fortsetzen, das dir gehört. Mit ccam braucht das Zweite nur einen Befehl, `ccam next --auto --continue`.

## Welche Limits hat Claude Code?

In den Pro- und Max-Tarifen hat Claude Code ein [5-Stunden-Sitzungslimit und ein Wochenlimit](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Beide gelten auch für Claude im Web und in den Apps, und die Nutzung zählt gleichzeitig auf beide. Wenn eines aufgebraucht ist, [zeigt Claude Code eine Meldung](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) wie „You've hit your session limit · resets 3:45pm“. Ein anderes Modell zu wählen hilft nicht, weil die Limits für alle Modelle gelten.

## Option 1: auf den Reset warten

Wenn du nur ein Konto hast, warte. Aktuelle Versionen von Claude Code [warten in der offenen Sitzung](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) und setzen die Aufgabe von selbst fort, sobald das Limit zurückgesetzt ist. `/usage` in Claude Code zeigt, wie viel du von den Limits deines Tarifs verbraucht hast.

## Option 2: sofort in einem anderen Konto weitermachen

Wenn du ein zweites Konto hast, kannst du dasselbe Gespräch sofort fortsetzen. Richte das einmal ein:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Die Statuszeile zeigt jetzt beim Arbeiten das Konto und seine 5-Stunden-Nutzung, sodass dich ein Limit nicht überrascht.

Wenn das Limit erreicht ist, beende Claude Code (`/exit`) und führe dann aus:

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

Öffne das Gespräch im Konto mit dem meisten Spielraum:

```sh
ccam next --auto --continue
```

Claude Code startet als `work` mit dem neuesten Gespräch des aktuellen Ordners, einschließlich allem, was bisher gesagt wurde. Mit geteiltem Verlauf (`ccam group set default --add history` oben) arbeiten beide Konten am selben Gespräch; ohne ihn kopiert Claude Code das Gespräch nach `work`, und ccam sagt dir, wie viele Tokens die erste Antwort erneut lesen muss.

Mit dem [ccam-Plugin](/docs/install#claude-code-plugin) in Claude Code zeigt eine Sitzung, die an ihrem Limit anhält, eine Benachrichtigung und vermerkt das volle Fenster, sodass `ccam next --auto` dieses Konto vor seinem Zurücksetzen nicht erneut wählt. `/ccam:switch` (oder `/ccam:switch work`) merkt das Gespräch für `ccam next --continue` vor, sodass genau dieses Gespräch fortgesetzt wird.

**Team- und Enterprise-Plätze.** Die automatische Kontoauswahl (`ccam next`, `--fallback`, `ccam pick`, `n` im Dashboard und, sobald du es einschaltest, [`auto_switch`](/features/next#auto-switch)) wechselt auch zwischen Plätzen derselben Claude Team- oder Enterprise-Organisation, was gegen die Bedingungen von Anthropic verstoßen kann. Lies die [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) und kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt.

## Wechselt ccam das Konto automatisch?

Nur wenn du es einschaltest.

- Standardmäßig läuft ccam nicht, während du arbeitest: Es startet Claude Code und ersetzt sich selbst durch diesen Prozess, also entscheidest du, wann gewechselt wird. `ccam next --auto` wählt das Konto mit der geringsten 5-Stunden-Nutzung, überspringt Konten, die bei 100 % stehen oder nicht angemeldet sind, und bleibt innerhalb einer Team- oder Enterprise-Organisation oder unter privaten Tarifen.
- Mit `relay = true` in `config.toml` läuft ccam stattdessen neben der Sitzung weiter, sodass `/ccam:switch` das Gespräch von selbst fortsetzt und dabei die Optionen, den Berechtigungsmodus, die Effort-Stufe und den Agenten der Sitzung beibehält; bevor ccam das Gespräch in ein Konto einer anderen Organisation verschiebt, fragt es dich.
- Mit `auto_switch = true` (Linux und macOS; experimentell) nimmt ccam den Wechsel auch selbst vor: Wenn Claude Code die Sitzung an ihrem Limit stoppt und ein anderes Konto derselben Art freies Kontingent hat, beendet ccam die Sitzung und setzt das Gespräch dort fort ([Bedingungen und Kosten](/features/next#auto-switch)).

Manche anderen Tools tauschen stattdessen die Anmeldung einer laufenden Sitzung aus; [ccam und andere Tools zum Kontowechsel](/compare/account-switchers) erklärt den Unterschied.

## Was du beachten solltest

- Führe `ccam next --continue` im selben Ordner aus, in dem du gearbeitet hast; Gespräche werden pro Projektordner gespeichert.
- Arbeite nicht gleichzeitig mit zwei Konten am selben Gespräch. ccam warnt dich, wenn du es versuchst.
- Ohne geteilten Verlauf wird das Gespräch in das andere Konto kopiert, und die beiden Kopien gehen getrennte Wege; die Checkpoints, aus denen `/rewind` Code wiederherstellt, bleiben beim ersten Konto.
- Die erste Nachricht im neuen Konto kann mehr von dessen Limit verbrauchen als üblich, weil dieses Konto das Gespräch noch nicht im Cache hat. `ccam next --continue --brief` beginnt es stattdessen mit einer kurzen Zusammenfassung.
- Ob du so zwischen Konten wechseln darfst, hängt von den Bedingungen ab, die für jedes Konto gelten.

Alle Optionen findest du unter [Nutzungslimits](/features/usage) und [Zum nächsten Konto wechseln](/features/next).
