---
title: Ein Konto für Skripte wählen
description: "ccam pick gibt das Konto mit dem meisten Spielraum für Skripte aus, verteilt gleichzeitig gestartete Worker auf Konten und endet mit 75 ohne Spielraum."
updated: 2026-11-05
order: 12
translation: machine
sourceHash: "0ec25752ace9"
---

`ccam pick` gibt den Namen des angemeldeten Kontos mit dem meisten Spielraum aus, damit ein Skript Claude Code damit starten kann. Gleichzeitig gestartete Worker erhalten jeweils ein anderes Konto, solange genug Konten Spielraum haben, sodass sich ein Stapel von `claude -p`-Jobs auf deine Konten verteilt, statt ein einziges aufzubrauchen.

**Team- und Enterprise-Plätze.** Die automatische Kontoauswahl (`ccam next`, `--fallback`, `ccam pick`, `n` im Dashboard und, sobald du es einschaltest, [`auto_switch`](/features/next#auto-switch)) wechselt auch zwischen Plätzen derselben Claude Team- oder Enterprise-Organisation, was gegen die Bedingungen von Anthropic verstoßen kann. Lies die [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) und kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt.

## Verwendung

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

Exit-Codes: `0` mit einem Namen auf der Standardausgabe, `75` (`EX_TEMPFAIL`), wenn kein Konto Spielraum hat oder keines angemeldet ist, `1` bei anderen Fehlern wie einer ungültigen Option.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--json` | Das Konto, seine 5-Stunden- und Wochennutzung, die Herkunft der Zahlen und den Grund der Wahl ausgeben | aus |
| `--model <model>` | Auch das Wochenfenster dieses Modells berücksichtigen, etwa `opus` | keins |
| `--max-usage <1-100>` | Konten unter dieser Nutzung bevorzugen | `auto_max_usage`, andernfalls 100 |
| `--strategy headroom\|expiring\|balance` | Die niedrigste Wochennutzung unter den Konten, deren 5-Stunden-Nutzung höchstens 3 Prozentpunkte über der niedrigsten liegt, die früheste 5-Stunden-Zurücksetzung oder das meiste verbleibende Wochenkontingent pro Tag | `auto_strategy`, andernfalls `headroom` |
| `--cross-org` | Auch Konten anderer Organisationen wählen | `cross_org`, andernfalls aus |
| `--wait[=<duration>]` | Wenn kein Konto Spielraum hat, auf die Zurücksetzung des ersten warten, höchstens so lange. Der Wert steht nach `=`: `--wait 30m` wird mit dem Hinweis `--wait takes its value after '=': --wait=30m` abgelehnt | aus; `5h`, wenn ohne Wert angegeben |
| `--dry-run` | Wählen, ohne die Wahl zu erfassen | aus |

## So funktioniert es

1. ccam listet die angemeldeten Konten (API-Schlüssel-Konten und Console-Anmeldungen haben keine Limits und werden ausgelassen) in der Organisation des Kontos auf, das dieser Ordner starten würde (siehe [Organisationen](/features/next#auto)); ein Ordner, der beim selben Claude-Konto angemeldet ist wie dieses, wird ebenfalls ausgelassen.
2. Es setzt eine Sperre (`~/.local/state/ccam/picks.lock`), sodass gleichzeitig laufende Auswahlen nacheinander erfolgen.
3. Es fragt die Nutzung aller Konten gleichzeitig ab, genauso wie `ccam usage`, mit dem 60-Sekunden-Cache. Ein Zeitfenster, das Claude Code als voll gemeldet hat, zählt bis zu seinem Zurücksetzen als 100 % ([der Limit-Eintrag](/features/usage#the-limit-record)), sodass ein Konto, das gerade sein Limit erreicht hat, ausgeschlossen wird, auch wenn seine gespeicherten Zahlen noch Spielraum zeigen.
4. Konten, die in den letzten 2 Minuten gewählt wurden, von `ccam pick` oder von [`auto_switch`](/features/next#auto-switch), das eine Sitzung dorthin verschoben hat (`~/.local/state/ccam/picks/<name>`; eine Wahl mit einem Zeitstempel in der Zukunft zählt nicht), werden zuerst ausgelassen, und die [Auswahlregeln](/features/next#auto) von `ccam next --auto` wählen unter den übrigen: zuerst unter `auto_max_usage`, nach der Strategie, nie ein Konto bei 100 %. Überspringt das das Konto, das die Regeln sonst wählen würden, sagt das der `reason` von `--json`: `30% of its 5 hours used; work was picked a moment ago`.
5. Wurden alle Konten mit Spielraum kürzlich gewählt, wählen die Regeln wieder unter allen, sodass sich bei drei Konten ein vierter Worker eines davon teilt.
6. Die Wahl wird erfasst und die Sperre freigegeben, dann wird der Name ausgegeben.

Das 2-Minuten-Fenster ist wichtig, weil die Nutzung eines neuen Jobs eine Weile braucht, bis sie in den Zahlen erscheint: Ohne es würden fünf gleichzeitig gestartete Worker alle dasselbe Konto als das mit dem meisten Spielraum sehen.

## Sicherheit und Sonderfälle

- `ccam pick` liest nur die Nutzung und schreibt den Zeitpunkt der Wahl; es erneuert oder liest Tokens nie über das hinaus, was `ccam usage` tut.
- Eine Wahl reserviert das Konto nicht. Zwei Skripte können nach 2 Minuten trotzdem mit demselben Konto starten; verwende `--fallback`, um weiterzumachen, wenn dieses Konto sein Limit erreicht.
- Warten mit `--wait` hält keine Sperre, sodass andere Auswahlen währenddessen weiterlaufen.
- `ccam pick` ignoriert [`ccam folder`](/features/folder): Es wählt nur nach Nutzung.

## Verwandte Themen

- [Zum nächsten Konto wechseln](/features/next)
- [Am Limit erneut ausführen](/features/next#run-again-at-the-limit)
- [Nutzungslimits](/features/usage)
- [API-Schlüssel-Konten](/features/api)
