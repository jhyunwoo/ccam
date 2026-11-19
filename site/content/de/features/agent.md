---
title: Subagenten auf anderen Konten
description: "ccam agent führt einen Teil einer Claude Code-Aufgabe mit einem anderen Konto aus und gibt Claudes Antwort zurück, so nutzt eine lange Aufgabe mehrere Konten."
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Die eigenen Subagenten von Claude Code laufen mit dem Konto der Sitzung, die sie startet, und verbrauchen daher dessen Nutzung. `ccam agent "<task>"` führt eine Aufgabe mit `claude -p` im selben Ordner mit einem anderen deiner Konten aus und gibt die letzte Antwort von Claude aus. Aus einer Claude Code-Sitzung aufgerufen, arbeitet es wie ein Subagent, der die Nutzung eines anderen Kontos verbraucht, sodass eine lange Aufgabe weiterlaufen kann, wenn das eigene Konto der Sitzung schon erschöpft wäre.

**Team- und Enterprise-Plätze.** `ccam agent` wählt das Konto selbst. Innerhalb einer Claude Team- oder Enterprise-Organisation führt es die Aufgabe auf einem anderen Platz aus, was gegen die Bedingungen von Anthropic verstoßen kann. Lies die [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) und kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt.

## Verwendung

Claude Code führt es über sein Bash-Werkzeug aus. Ist das [ccam-Plugin](/docs/install#claude-code-plugin) installiert, bitte einfach mit eigenen Worten darum, zum Beispiel „nutze meine anderen Konten für Teile dieser Aufgabe“, oder schreib das in deine CLAUDE.md. Der Skill des Plugins weist Claude Code an, jeden Agenten im Hintergrund zu starten, eine Aufgabe zu formulieren, die für sich allein steht, und Agenten, die Dateien bearbeiten, einen eigenen Teil oder einen eigenen git-Worktree zu geben.

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

Die letzte Antwort von Claude geht auf die Standardausgabe; die Zeilen, die mit `ccam:` beginnen, gehen auf die Standardfehlerausgabe. `--json` gibt stattdessen `{account, session_id, result, is_error, accounts}` aus, wobei `accounts` jedes Konto auflistet, mit dem die Aufgabe lief.

Exit-Codes: die von claude selbst (`0`, wenn der Agent fertig ist), `75` (`EX_TEMPFAIL`), wenn kein anderes Konto als das aufrufende Spielraum hat, `1` bei anderen Fehlern, etwa wenn ein Agent versucht, einen weiteren zu starten.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--model <model>` | Das Modell des Agenten; auch das Wochenfenster dieses Modells zählt bei der Wahl | das Standardmodell des Kontos |
| `--resume <id>` | Das Gespräch eines früheren Agenten mit einer Folgeaufgabe fortsetzen (die Sitzungs-ID steht in seiner ersten und letzten Zeile) | aus |
| `--json` | Konto, Sitzung, Antwort und die Konten ausgeben, mit denen die Aufgabe lief | aus |
| `--max-usage <1-100>` | Konten unter dieser Nutzung bevorzugen | `auto_max_usage`, andernfalls 100 |
| `--strategy headroom\|expiring\|balance` | Die [Auswahlregel](/features/next#auto) | `auto_strategy`, andernfalls `headroom` |
| `--cross-org` | Auch mit Konten anderer Organisationen ausführen | `cross_org`, andernfalls aus |
| `--wait[=<duration>]` | Wenn kein Konto Spielraum hat, auf die Zurücksetzung des ersten warten, höchstens so lange | aus; `5h`, wenn ohne Wert angegeben |
| `--dry-run` | Den Befehl für das Konto ausgeben, das es wählen würde, ohne ihn auszuführen | aus |

## So funktioniert es

1. **Das aufrufende Konto bleibt außen vor.** Das Konto, aus dem `ccam agent` aufgerufen wird (das der Sitzung, die ccam gestartet hat, sonst das, das dieser Ordner startet), führt den Agenten nie aus, ebenso wenig ein anderer Ordner, der beim selben Claude-Konto angemeldet ist und sich dessen Limits teilt. Der Agent bleibt in der [Organisation](/features/next#auto) des Aufrufers, sofern `--cross-org` nichts anderes sagt.
2. **Es wählt wie `ccam pick`.** Die Nutzung wird so abgefragt, wie es [`ccam pick`](/features/pick) tut, und Agenten, die innerhalb von 2 Minuten nacheinander starten, erhalten verschiedene Konten, solange genug Konten Spielraum haben.
3. **Es läuft mit den Berechtigungen der Sitzung, nie mit mehr.** In `claude -p` ist niemand da, der ein Werkzeug erlauben könnte, daher erhält der Agent den Berechtigungsmodus und die Effort-Stufe der Sitzung, die ihn aufgerufen hat. Aus einem Terminal gestartet, gelten der Standardmodus von Claude Code und die Berechtigungsregeln deiner Einstellungen. `ccam agent` hat keine Option, die mehr erlaubt.
4. **An einem Nutzungslimit geht das Gespräch weiter.** Erreicht das Konto des Agenten sein Limit, hält ccam das Limit so fest, wie es der [Plugin-Hook](/features/usage#the-limit-record) tun würde, und setzt dasselbe Gespräch mit dem Konto mit dem meisten Spielraum fort: Claude Code kopiert es dorthin (`--resume <file> --fork-session`), mit der ersten Nachricht „Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.“ Jedes Konto wird nur einmal versucht.
5. **Jedes Agenten-Gespräch wird festgehalten**, in `~/.local/state/ccam/agent/<session id>.json` (30 Tage). `--resume <id>` setzt es in seinem Ordner fort, mit dem Konto, das es ausgeführt hat, wenn dieses noch Spielraum hat, sonst mit einem anderen. [`ccam next --continue`](/features/next) lässt Agenten-Gespräche aus, sodass es weiterhin deine eigenen fortsetzt.

## Sicherheit und Sonderfälle

- Ein Agent kann keine eigenen Agenten starten: ccam markiert ihn mit `CCAM_AGENT=1`, und `ccam agent` verweigert dort die Ausführung.
- Wird `ccam agent` gestoppt (Ctrl+C, ein Signal oder Claude Code, das den Befehl unterbricht), stoppt auch der Agent.
- Der Agent sieht nicht das Gespräch, das ihn aufgerufen hat, sondern nur seine Aufgabe.
- Agenten, die gleichzeitig dieselben Dateien bearbeiten, überschreiben gegenseitig ihre Arbeit. Gib jedem seinen eigenen Teil der Aufgabe oder seinen eigenen git-Worktree.
- Wechselt die aufrufende Sitzung zu einem anderen Konto ([`auto_switch`](/features/next#auto-switch)), stoppen die Agenten, die sie gestartet hat, mit ihr. Ihre Gespräche bleiben erhalten: `--resume <id>` setzt sie fort.
- `agent` ist ein Befehlsname, daher kann kein neues Konto `agent` heißen; ein bestehendes startet weiterhin mit `ccam run agent`.

## Verwandte Themen

- [Ein Konto für Skripte wählen](/features/pick)
- [Zum nächsten Konto wechseln](/features/next)
- [Nutzungslimits](/features/usage)
- [Hintergrundsitzungen aller Konten](/features/agents)
