---
title: Häufig gestellte Fragen
metaTitle: "ccam-FAQ: mehrere Claude Code-Konten"
description: Antworten zu ccam, Claude Code-Konten, Nutzungslimits, Teilen, Datenschutz und unterstützten Systemen.
updated: 2026-11-17
order: 6
translation: machine
sourceHash: "9b13b37391f2"
---

## Ist ccam ein offizielles Tool von Anthropic?

Nein. ccam ist ein inoffizielles Open-Source-Projekt und steht in keiner Verbindung zu Anthropic und wird nicht von Anthropic unterstützt. Es verwendet die Einstellung `CLAUDE_CONFIG_DIR`, die [Claude Code dafür dokumentiert](https://code.claude.com/docs/en/env-vars), mehrere Konten nebeneinander zu betreiben.

## Darf ich mehrere Claude Code-Konten nutzen?

Das hängt von den Bedingungen jedes Kontos ab und bei Arbeitskonten von den Richtlinien deiner Organisation: den [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) von Anthropic für Pro und Max, seinen [Commercial Terms](https://www.anthropic.com/legal/commercial-terms) für Team und Enterprise und seiner [Usage Policy](https://www.anthropic.com/legal/aup) für alle. Prüfe sie, bevor du mehrere Konten nutzt. ccam erhöht die Limits eines Kontos nicht, legt sie nicht zusammen und verbirgt sie nicht, und jedes Konto hält weiterhin an seinen eigenen an. Es verschiebt deine Arbeit auf ein anderes deiner Konten, wenn du darum bittest, oder von selbst an einem Limit, wenn du [`auto_switch`](/features/next#auto-switch) einschaltest (Linux und macOS). Prüfe, ob das für deine Konten erlaubt ist, bevor du es nutzt; die Entscheidung liegt bei dir.

**Team- und Enterprise-Plätze.** Die automatische Kontoauswahl (`ccam next`, `--fallback`, `ccam pick`, `n` im Dashboard und [`auto_switch`](/features/next#auto-switch), wenn du es einschaltest) wechselt auch zwischen Plätzen derselben Claude Team- oder Enterprise-Organisation. Im Hilfe-Center von Anthropic ([Team-Tarife](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); ebenso für [Enterprise-Tarife](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)) heißt es: „Each seat on a Team plan is assigned to one named individual, and only that person may use it“ (jeder Platz in einem Team-Tarif ist einer namentlich benannten Person zugewiesen, und nur diese Person darf ihn nutzen) und „Users can't route their work through another user's account to use that user's remaining usage“ (Benutzer dürfen ihre Arbeit nicht über das Konto eines anderen Benutzers leiten, um dessen verbleibendes Nutzungskontingent zu verwenden). Wenn du die automatische Auswahl zwischen Plätzen nutzt, kann das daher gegen die Bedingungen von Anthropic verstoßen. Ob du sie dort nutzt, ist deine Entscheidung: Kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt. Für private Tarife (Pro, Max) gilt der Absatz darüber.

## Umgeht ccam Nutzungslimits?

Nicht über das eigene Limit eines Kontos hinaus: Jedes Konto behält seine Limits und hält an ihnen an, und ccam erhöht sie nicht, legt sie nicht zusammen und verbirgt sie nicht. ccam zeigt sie an und lässt dich in einem anderen deiner Konten weitermachen, wenn eines erschöpft ist; mit eingeschaltetem `auto_switch` tut es das von selbst (siehe die nächste Frage). Zwischen Plätzen einer Team- oder Enterprise-Organisation verbraucht das Weitermachen auf einem anderen Platz dessen verbleibendes Nutzungskontingent, was das Hilfe-Center von Anthropic für den Platz eines anderen Benutzers ausschließt (siehe Team- und Enterprise-Plätze oben).

## Kann ccam von selbst das Konto wechseln, wenn eine Sitzung ihr Limit erreicht?

Nur wenn du es einschaltest. Mit `auto_switch = true` am Anfang von `~/.config/ccam/config.toml` (Linux und macOS; standardmäßig aus und experimentell) gilt: Wenn Claude Code eine interaktive Sitzung (eine, die ccam gestartet hat) an einem 5-Stunden-, Wochen-, Opus-, Sonnet- oder Fable-Limit stoppt, beendet ccam die Sitzung und setzt das Gespräch als ein anderes deiner Konten mit freiem Kontingent fort: einen anderen Platz derselben Team- oder Enterprise-Organisation oder einen anderen privaten Tarif (Pro, Max). Es wechselt nie vor dem Limit, überwacht deine Nutzung nie im Hintergrund und verschiebt Arbeit nie aus einer Team- oder Enterprise-Organisation heraus oder in eine hinein.

Bei Plätzen verschiebt es deine Arbeit von einem Platz auf einen anderen, was gegen die Bedingungen von Anthropic verstoßen kann (siehe Team- und Enterprise-Plätze oben). Für private Tarife regeln die Bedingungen von Anthropic die Frage nicht eindeutig: Die Consumer Terms verbieten, auf die Dienste „through automated or non-human means“ (über automatisierte oder nicht-menschliche Mittel) zuzugreifen, sofern Anthropic es nicht ausdrücklich erlaubt, und ein Wechsel sendet eine Nachricht in deinem Namen. Ob du es einschaltest, ist deine Entscheidung.

[`auto_switch`](/features/next#auto-switch) nennt jede Bedingung und was ein Wechsel kostet.

## Wie nutze ich zwei Claude Code-Konten auf einem Computer?

Gib jedem Konto mit `CLAUDE_CONFIG_DIR` einen eigenen Konfigurationsordner. ccam erledigt das für dich: Führe `ccam add` für jedes Konto aus und starte sie dann mit `ccam personal` und `ccam work` oder mit Kurzbefehlen, die du selbst benennst, etwa `ccwork`. Die [Anleitung](/guides/multiple-claude-code-accounts) zeigt auch, wie es von Hand geht.

## Wie setze ich ein Gespräch in einem anderen Konto fort?

Führe `ccam next --continue` im Ordner aus: Es nimmt das letzte Gespräch des Ordners ins nächste Konto mit (`--auto` für das Konto mit dem meisten Spielraum). Innerhalb von Claude Code macht `/ccam:switch` aus dem [ccam-Plugin](/docs/install#claude-code-plugin) dasselbe für die Sitzung, in der du bist. Teilen die beiden Konten den Gesprächsverlauf (`ccam group set default --add history` oder eine eigene Gruppe: `ccam group add pair --items history --members work,side`), öffnet das neue Konto genau dieses Gespräch; andernfalls kopiert Claude Code es in den Verlauf des neuen Kontos. Füge `--brief` hinzu, damit das neue Konto mit einer kurzen Zusammenfassung statt mit dem ganzen Gespräch beginnt.

## Sieht oder sendet ccam meine Anmeldung?

ccam teilt oder kopiert deine Claude-Anmeldung nie zwischen Konten und kopiert nie die Anmeldungen von MCP-Servern. Es liest das Zugriffstoken eines Kontos nur, um beim Nutzungs-Endpunkt von Anthropic die Nutzung dieses Kontos abzufragen: für `ccam usage`, das Dashboard, `ccam doctor` und die automatische Wahl (`ccam next --auto`, `ccam pick`, `--fallback` und `auto_switch`, wenn eine Sitzung ihr Limit erreicht). Ein abgelaufenes erneuert es nur mit `ccam usage --refresh`; mit `usage_source = "claude"` sendet es nirgendwohin ein Token und liest die gespeicherte Anmeldung (die Datei oder den Eintrag im macOS-Schlüsselbund) nur, um festzustellen, ob sie noch gültig ist. `ccam ext` und `ccam doctor` lesen die Namen der MCP-Anmeldungen eines Kontos, nie ihre Tokens. ccam hat keine Telemetrie. Abgesehen von diesen Abfragen stellt es von sich aus nur eine einzige Anfrage: die Prüfung des Dashboards bei GitHub auf ein neueres ccam-Release, die nichts über dich oder deine Konten enthält.

## Wie melde ich mich auf einem entfernten Server an?

Führe `ccam add` oder `ccam login <name>` wie gewohnt in der SSH-Sitzung aus. Claude Code wartet auf dem Server auf die Anmeldung, daher landet der Browser auf deinem Computer auf einer Seite `http://localhost:<port>/callback?code=…`, die nicht lädt. Kopiere diese Adresse aus dem Browser und füge sie im Terminal ein: ccam gibt sie an Claude Code auf dem Server weiter, das die Anmeldung abschließt. Die Adresse enthält einen Einmalcode, gib sie also nicht weiter.

## Wo werden meine Konten gespeichert?

In `~/.local/share/ccam/profiles/<name>/`, ein Ordner pro Konto. Unter macOS speichert Claude Code die Anmeldung selbst im Schlüsselbund, unter einem Namen, der von diesem Ordner abgeleitet ist, deshalb behält `ccam rename` dort den alten Namen des Ordners.

## Was passiert mit meinem vorhandenen ~/.claude?

Nichts, außer du führst `ccam import` oder `ccam migrate` aus. Der Import registriert den Ordner dort, wo er liegt, sodass einfaches `claude` und `ccam <name>` weiterhin dasselbe Konto und dieselben Einstellungen verwenden. `ccam migrate` kopiert seine Einstellungen, `CLAUDE.md`, Erweiterungen und MCP-Server in die geteilte Gruppe `default` und liest den Ordner dabei nur.

## Welche Systeme unterstützt ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine und andere), macOS sowie Windows 10 und 11, nativ oder innerhalb von WSL. Unter Windows braucht das Teilen von Dateien wie `settings.json` den Entwicklermodus; geteilte Ordner funktionieren auch ohne ihn, und ohne ihn legt das erste `ccam add` die Gruppe `default` nur mit MCP-Servern und Erweiterungen an ([Teilen unter Windows](/features/share#on-windows)).

## Warum erscheint ein MCP-Server nur in einem Konto?

Server, die innerhalb einer Sitzung mit `claude mcp add` hinzugefügt werden, werden nur für dieses Konto gespeichert, und ohne `-s user` nur für den aktuellen Ordner. ccam weist dich beim nächsten Start des Kontos darauf hin, und `ccam mcp list` und `ccam ext` zeigen sie an; verschiebe den Server mit `ccam mcp import <account> <server>` in die geteilte Gruppe.

## Macht ccam Claude Code langsamer?

Claude Code über ccam zu starten kostet etwa 4 Millisekunden, gemessen mit `scripts/perf-check.sh` auf einem Laptop. ccam liest eine kleine Datei, prüft die geteilten Links und ersetzt sich dann durch `claude`, läuft also standardmäßig nicht, während du arbeitest. Windows kann einen Prozess nicht ersetzen, deshalb startet ccam dort Claude Code und wartet darauf, ohne sonst etwas zu tun. Mit eingeschaltetem `relay` oder `auto_switch` läuft ccam auf dieselbe Weise neben jeder interaktiven Sitzung weiter, damit es das Gespräch in einem anderen Konto fortsetzen kann.

## Wie aktualisiere ich ccam?

Führe `ccam update` aus. Der Befehl lädt das neueste Release herunter, prüft seine Prüfsumme und ersetzt die Datei `ccam`; `ccam update --check` zeigt nur an, ob es ein neueres Release gibt. Auch das Dashboard zeigt dir, wenn es ein neueres Release gibt, und `U` installiert es nach einer Rückfrage; ccam aktualisiert sich nie von selbst. Wenn du das `.deb`-Paket installiert hast, installiere stattdessen das neue `.deb` mit apt. Siehe [Aktualisieren](/docs/install#update).

## Wie deinstalliere ich ccam?

Führe `ccam uninstall` aus. Es listet auf, was wegkommt, und fragt zuerst; dann entfernt es jedes Konto (vorher abgemeldet), die Kurzbefehle, `~/.config/ccam`, `~/.local/share/ccam`, `~/.local/state/ccam` und den Befehl `ccam`. Ein importiertes `~/.claude` bleibt, mit einer Kopie dessen, was es in seinen Gruppen geteilt hat. Siehe [Installation](/docs/install#uninstall).
