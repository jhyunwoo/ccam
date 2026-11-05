---
title: Zum nächsten Konto wechseln
description: ccam next startet Claude Code als anderes Konto; --auto wählt das mit dem meisten Spielraum, --continue nimmt das Gespräch mit, --fallback wiederholt -p-Jobs.
updated: 2026-11-05
order: 8
translation: machine
sourceHash: "a320c160caef"
---

`ccam next` öffnet Claude Code mit einem anderen Konto, ohne dass du es benennen musst. Verwende es, wenn ein Konto sein Nutzungslimit erreicht. Es gibt vier Wege zu wechseln:

- Einfaches `ccam next` nimmt das Konto nach dem zuletzt verwendeten, und `--auto` nimmt das mit dem meisten Spielraum. `--continue` setzt das letzte Gespräch dieses Ordners im neuen Konto fort, ob die beiden den Verlauf teilen oder nicht.
- In Claude Code meldet dir das [ccam-Plugin](#inside-claude-code), wenn ein Limit eine Sitzung anhält, und `/ccam:switch` gibt das Gespräch weiter.
- Für Skripte, die `claude -p` ausführen, erledigt `ccam --fallback` den Wechsel von selbst (siehe [Am Limit erneut ausführen](#run-again-at-the-limit)).
- Mit eingeschaltetem [`auto_switch`](#auto-switch) (Linux und macOS) verschiebt ccam auch eine interaktive Sitzung, die Claude Code an ihrem Limit stoppt.

**Team- und Enterprise-Plätze.** Die automatische Kontoauswahl (`ccam next`, `--fallback`, `ccam pick`, `n` im Dashboard und, sobald du es einschaltest, [`auto_switch`](/features/next#auto-switch)) wechselt auch zwischen Plätzen derselben Claude Team- oder Enterprise-Organisation, was gegen die Bedingungen von Anthropic verstoßen kann. Lies die [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) und kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt.

## Verwendung

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # continue this folder's last conversation on the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next -c --brief            # a new conversation that starts from a summary of the last one
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the most weekly quota left per day
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work CCAM_PROFILE=work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam gibt den Wechsel als `→ <from> → <to>` auf der Standardfehlerausgabe aus, bei einer Wahl mit `--auto` mit dem Grund, und startet dann Claude Code mit dem neuen Konto auf dieselbe Weise wie `ccam <name>`. `<from>` ist das zuletzt verwendete Konto oder mit `--continue` das Konto, zu dem das Gespräch gehört. Argumente nach `--` gehen unverändert an `claude`.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--auto` | Das angemeldete Konto nach Nutzung wählen (siehe unten) und Konten bei 100 % überspringen | aus |
| `--max-usage <1-100>` | Mit `--auto`: Konten unter dieser Nutzung bevorzugen | `auto_max_usage`, andernfalls 100 |
| `--strategy headroom\|expiring\|balance` | Mit `--auto`: die niedrigste Wochennutzung unter den Konten, deren 5-Stunden-Nutzung höchstens 3 Prozentpunkte über der niedrigsten liegt, die früheste 5-Stunden-Zurücksetzung oder das meiste verbleibende Wochenkontingent pro Tag | `auto_strategy`, andernfalls `headroom` |
| `--cross-org` | Mit `--auto`: auch Konten einer anderen Organisation wählen (siehe unten) | `cross_org`, andernfalls aus |
| `--wait[=<duration>]` | Mit `--auto`: Wenn kein Konto Spielraum hat, auf die Zurücksetzung des ersten warten, höchstens so lange. Der Wert steht nach `=`: `--wait 2h` wird mit dem Hinweis `--wait takes its value after '=': --wait=2h` abgelehnt | aus; `5h`, wenn ohne Wert angegeben |
| `-c`, `--continue` | Das letzte Gespräch dieses Ordners im neuen Konto fortsetzen (siehe unten) | aus |
| `--brief` | Mit `--continue`: im neuen Konto ein neues Gespräch aus ccams Zusammenfassung des letzten beginnen, statt das ganze Gespräch fortzusetzen ([siehe unten](#brief)) | aus |
| `-- <claude args>` | Argumente, die nach ccams eigenen an `claude` übergeben werden | keine |
| `--dry-run` | Wird unterstützt: gibt `→ <from> → <to>` und die Befehlszeile aus, die ausgeführt würde, ohne Claude Code zu starten oder einen Start zu erfassen | aus |

`--auto` und `--continue` lassen sich kombinieren. Mit `--continue` werden die `--resume`-Argumente vor deine eigenen gesetzt.

Eigene Optionen von Claude Code wie `--dangerously-skip-permissions` oder `--model` gehören hinter `--`. Anders als `ccam <name>`, das alles nach dem Namen weiterreicht, liest `ccam next` seine Argumente selbst und bricht ohne das `--` mit `ccam: unknown flag: --dangerously-skip-permissions` ab.

## So funktioniert es

Jeder Start über ccam schreibt die Uhrzeit nach `~/.local/state/ccam/last-used/<name>`. `ccam next` liest diese Dateien, prüft, welche Konten angemeldet sind, und wählt.

1. ccam baut aus `config.toml` pro Konto einen `next.Candidate` (`cli.candidates`): das Konto, seine letzte Startzeit (`launcher.LastUsed`) und ob es eine funktionierende Anmeldung hat (`creds.Has`: eine Claude-Anmeldung in `.credentials.json` oder unter macOS im Schlüsselbund, die weder vollständig abgelaufen ist noch von Claude Code gelöscht wurde). Ein Ordner, der stattdessen einen API-Schlüssel verwendet (eine Console-Anmeldung oder ein mit `ccam add --api` hinzugefügtes Konto), hat keine Tariflimits und wird nie gewählt.
2. Ohne `--auto` sortiert ccam die Konten nach Slot-Nummer und findet das zuletzt gestartete (`next.Rotate`). Von diesem Konto wechselst du weg. Wurde noch kein Konto gestartet, ist es der niedrigste Slot.
3. Von dort geht es in Slot-Reihenfolge weiter, springt vom letzten Slot zurück zum ersten und gibt das erste angemeldete Konto zurück.
4. ccam bereitet das gewählte Konto genau so vor und startet es, wie es `ccam <name>` tun würde (`launcher.LaunchWith`): Es prüft die Links des Kontos, wendet das MCP-Teilen an, registriert die Sitzung und ersetzt ccam über `syscall.Exec` durch `claude`.

Der Durchlauf beginnt beim zuletzt verwendeten Konto und endet beim ersten angemeldeten Konto danach:

<!-- motion:next-rotate -->

`work` wurde zuletzt gestartet, daher beginnt der Durchlauf danach. `side` hat keine Anmeldung und wird übersprungen, und `lab` wird gewählt. Wäre auch `lab` abgemeldet, würde der Durchlauf zu `personal` zurückspringen (`% len` im Diagramm). Die Rotation gibt nie das Konto zurück, von dem du gestartet bist: Gibt es kein anderes angemeldetes Konto, endet sie mit einem Fehler.

### --auto

`--auto` entscheidet nach Nutzung statt nach Slot-Reihenfolge. ccam behält die angemeldeten Konten, fragt ihre Nutzung auf dieselbe Weise wie `ccam usage` gleichzeitig ab (`Fetcher.FetchAll`) und wählt nach den Regeln unten (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

Ein Konto wird nach seinem 5-Stunden- und seinem Wochenfenster beurteilt und, wenn `--model` unter den Argumenten für Claude Code steht (hinter `--`), auch nach dem Wochenfenster dieses Modells (`opus`, `sonnet`, `haiku` oder `fable` im Namen). ccam teilt die Konten in drei Stufen ein (`Policy.Choose()`):

1. **Spielraum**: jedes Fenster unter `auto_max_usage` (`config.toml`, Standard 100). Die Strategie wählt unter ihnen: `headroom` (Standard) nimmt unter den Konten, deren 5-Stunden-Nutzung höchstens 3 Prozentpunkte über der niedrigsten liegt, das, dessen Woche am wenigsten verbraucht ist, sodass ein Konto bei 12 % eines bei 10 % schlägt, wenn von seiner Woche weniger verbraucht ist, eines bei 14 % aber nicht; `expiring` nimmt das 5-Stunden-Fenster, das am frühesten zurückgesetzt wird, damit sein restliches Kontingent genutzt wird, bevor es verfällt; `balance` nimmt das Konto mit dem meisten verbleibenden Wochenkontingent pro Tag bis zur Zurücksetzung seiner Woche und rechnet jede laufende ccam-Sitzung dagegen, damit kein Wochenkontingent ungenutzt bleibt und sich parallele Sitzungen verteilen.
2. **Darüber**: unter 100 %, aber bei oder über `auto_max_usage`. Wird nur verwendet, wenn kein Konto Spielraum hat; das Konto, das `headroom` unter ihnen nehmen würde, kommt zuerst.
3. **Voll**: ein Fenster bei 100 % oder unbekannte 5-Stunden-Nutzung. Wird nie gewählt.

Ein Zeitfenster zählt bis zu seinem Zurücksetzen als 100 %, wenn Claude Code dort eine Anfrage abgelehnt hat: Der Hook des [ccam-Plugins](#inside-claude-code) hält das fest (`~/.local/state/ccam/usage/<name>.limit.json`), egal was die gespeicherten Zahlen sagen. Ein Konto, das gerade sein Limit erreicht hat, wird also vor seinem Zurücksetzen nicht wieder gewählt, auch wenn seine Zahlen alt sind.

Einstellungen für alle Konten gehören an den Anfang von `config.toml`, vor das erste `[[group]]`, `[[profile]]` oder `[[folder]]`. Eine Einstellung unter einer dieser Tabellen gehört zu ihr: ccam warnt bei jedem Start davor, und Befehle, die die Konfiguration speichern, lehnen ab, bis sie verschoben ist.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organisationen.** Automatische Wahlen bleiben bei Konten derselben Art wie dem, das du verlässt: Plätze in derselben Team- oder Enterprise-Organisation oder private Tarife (Pro, Max), sodass Arbeit aus einem Firmenkonto nicht in einem privaten weitergeht oder umgekehrt. ccam kennt die Organisation jedes Kontos aus seiner Anmeldung (`ccam login` hält sie fest; ältere Anmeldungen werden aus `.claude.json` gelesen); ein Konto, dessen Organisation ccam nicht erkennen kann, kann mit jedem Konto zusammenpassen. Zwei Ordner, die beim selben Claude-Konto angemeldet sind, teilen sich dessen Limits, daher übernimmt nie einer vom anderen. `--cross-org` oder `cross_org = true` in `config.toml` hebt die Organisationsregel auf, aber nach [`/ccam:switch`](#inside-claude-code) wartet ein Wechsel in eine andere Organisation weiterhin auf dein Ja.

Wenn nichts gewählt werden kann, beendet sich ccam mit Code 75 (`EX_TEMPFAIL`), damit ein Skript „später erneut versuchen“ von einem Fehler unterscheiden kann. Mit `--wait` schläft es stattdessen, bis das blockierende Fenster des ersten vollen Kontos zurückgesetzt wird (30 Sekunden nach dem Zurücksetzungszeitpunkt, wenn der Endpunkt das neue Fenster zeigt), gibt `ccam: no account has room; waiting until 15:10` aus und wählt erneut; wenn kein Konto angibt, wann es zurückgesetzt wird, sieht es nach 10 Minuten wieder nach. Wird vor dem Ende der Wartezeit nichts frei, beendet es sich mit 75.

Die Abfrage folgt den Regeln von `ccam usage`: Zahlen, die weniger als 60 Sekunden alt sind, werden wiederverwendet, die Abfrage jedes Kontos hat ein Timeout von 3 Sekunden, ein Konto, das nicht antwortet oder dessen Zugriffstoken abgelaufen ist, wird nach seinen gespeicherten Zahlen beurteilt, und ein Limit-Eintrag hat Vorrang vor ihnen. Ein gespeichertes Zeitfenster, dessen Zurücksetzungszeitpunkt vorbei ist, zählt als 0 %, damit ein Konto, dessen Fenster inzwischen zurückgesetzt wurde, nicht übersprungen wird. `--auto` erneuert nie Tokens. Anders als die einfache Rotation kann `--auto` das zuletzt verwendete Konto wählen, wenn es noch den meisten Spielraum hat.

### --continue

1. ccam sucht in jedem Konto nach dem letzten Gespräch dieses Ordners: unter den `<session>.jsonl`-Dateien in `projects/<folder>/` (der Ordnername ist der Arbeitsordner mit aufgelösten Links, in dem jedes Zeichen außer Buchstaben und Ziffern zu `-` wird, so wie Claude Code ihn benennt) die mit dem neuesten Prompt oder der neuesten Antwort. Sitzungen ohne beides, etwa eine Sitzung, die nur `/mcp`, `/usage` oder `/status` ausgeführt hat, oder eine leere Datei, werden übersprungen. Hat eine Sitzung eine Notiz hinterlassen (ein Nutzungslimit oder `/ccam:switch`, siehe unten) und hat kein Gespräch einen Prompt oder eine Antwort mehr als 2 Minuten danach, wird stattdessen genau diese Sitzung verwendet; die Zeilen, die Claude Code bei `/exit` schreibt, zählen nicht.
2. Das Konto, zu dem es gehört, ist das, von dem du wechselst, und das neue Konto wird wie oben gewählt, innerhalb derselben Organisation. Teilen sich mehrere Konten diesen Verlauf, gehört es dem Konto, dessen `.claude.json` angibt, dass es das Gespräch zuletzt ausgeführt hat (`lastSessionId`), andernfalls dem zuletzt verwendeten.
3. Teilt das neue Konto diesen Verlauf (derselbe Gruppenordner), startet ccam `claude --resume <session id>`, also genau dieses Gespräch (`--continue` würde das jeweils neueste Gespräch im Ordner öffnen).
4. Andernfalls startet es `claude --resume <the conversation's file> --fork-session`: Claude Code kopiert das ganze Gespräch unter einer neuen Sitzungs-ID in den eigenen Verlauf des neuen Kontos, und das Original bleibt, wie es war. Von da an findet `--continue` des neuen Kontos es.
5. Sofern du nicht `--agent` nach `--` angibst, übergibt ccam auch den Agenten, mit dem das Gespräch läuft: die Einstellung `agent` aus `.claude/settings.local.json` oder `.claude/settings.json` des Ordners oder aus der `settings.json` des Kontos, andernfalls den Agenten, den die letzte `agent-setting`-Zeile des Gesprächs nennt. Claude Code holt den Agenten nicht zurück, wenn es ein Gespräch kopiert (Schritt 4), und würde als sein Standard-Agent mit allen Werkzeugen weitermachen. Teilt das neue Konto die Agenten, Plugins und Einstellungen des alten nicht (dieselben Ordner `agents` und `plugins` und dieselbe `settings.json`, über eine Gruppe), fügt ccam nach der Zeile über den Wechsel eine Zeile hinzu: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

ccam sagt, wie groß das Gespräch ist: Das neue Konto liest es mit seiner ersten Antwort komplett erneut. Prompt-Caches werden nie zwischen Organisationen geteilt, und jeder private Tarif ist eine eigene Organisation; gehören die beiden Konten zu verschiedenen, ergänzt ccam daher `prompt caches are not shared between organizations, so the first reply reads it all again`. Ist das alte Konto bei MCP-Servern angemeldet, bei denen das neue es nicht ist, nennt ccam sie: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. Wenn eine Gruppe `history` teilt, entfällt die Kopie, es bleibt ein Gespräch, und auch die Checkpoints, aus denen `/rewind` Code wiederherstellt (`file-history/`), werden geteilt.

### --brief

`--continue` setzt das ganze Gespräch fort, und das neue Konto liest es mit seiner ersten Antwort komplett erneut. `ccam next --continue --brief` beginnt stattdessen ein neues Gespräch, ausgehend von einer Zusammenfassung, die ccam vom alten schreibt, ohne ein Modell zu fragen:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

Die Zusammenfassung enthält, in dieser Reihenfolge und nur, wo es etwas zu sagen gibt:

1. den Titel des Gesprächs;
2. den Pfad seines Transkripts, damit Claude darin suchen kann, wenn es ein Detail braucht;
3. deine erste Anfrage, wenn es mehr als vier gab;
4. die Zusammenfassung, die Claude Code beim letzten Komprimieren des Gesprächs geschrieben hat;
5. deine letzten vier Anfragen, mit `!`-Befehlen und Slash-Befehlen so, wie du sie getippt hast;
6. Claudes letzte Antwort;
7. die Todo-Liste;
8. die 10 Dateien, die Claude zuletzt geändert hat;
9. `git status --short --branch` des Ordners, bis zu 800 Zeichen.

Die Todo-Liste behält höchstens 10 Einträge: Erledigte Einträge fallen zuerst weg, dann die Schritte, die am weitesten in der Zukunft liegen, sodass der aktuelle Schritt und die nächsten bleiben. Lange Teile werden in der Mitte gekürzt, und die ganze Zusammenfassung bleibt unter 16.000 Zeichen.

ccam speichert sie als `~/.local/state/ccam/brief/<session>.md`, nur für dich lesbar, und das nächste `--brief` entfernt Zusammenfassungen, die älter als 7 Tage sind. Claude Code startet mit einem Prompt und hängt die Datei an, die er erwähnt: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` Ein Prompt, den du nach `--` angibst, tritt an die Stelle des letzten Satzes (`Then: <your prompt>`), und deine anderen Optionen nach `--` folgen ihm.

Claude Code beendet eine Erwähnung an einem `"` und schneidet sie an einem `#` ab; enthält der Pfad der Zusammenfassung eines davon, erwähnt ccam sie daher über ihren Pfad relativ zum Arbeitsordner; enthält auch dieser Pfad eines, nennt der Prompt die Datei, und Claude liest sie mit seinem Read-Werkzeug, das vorher nachfragt, weil die Datei außerhalb der Arbeitsordner liegt (in einem `-p`-Lauf im manuellen Modus verweigert Claude Code das Lesen, und Claude macht ohne die Zusammenfassung weiter).

`ccam --dry-run next --continue --brief` gibt die Zusammenfassung aus und speichert nichts. Das alte Gespräch bleibt, wie es war, und das neue weiß nur, was in der Zusammenfassung steht; nimm `--brief` also, wenn die Aufgabe wichtiger ist als jedes Detail des Wegs dorthin. Ohne `--continue` bricht ccam mit `--brief works with --continue` ab; kann es das Gespräch nicht lesen oder findet darin nichts zum Zusammenfassen, bricht es ab und sagt, dass `ccam next --continue` ohne `--brief` es vollständig fortsetzt. Das [Relay](#relay) und [`auto_switch`](#auto-switch) setzen immer das ganze Gespräch fort. Die eigene Option `--brief` von Claude Code gehört wie seine anderen Optionen hinter `--`.

## In Claude Code

Das [ccam-Plugin](/docs/install#claude-code-plugin) fügt zwei Hooks hinzu (Claude Code 2.1.78 oder neuer):

- Wenn eine Runde an einem Nutzungslimit anhält, merkt es das Gespräch vor und zeigt in Terminals, die das unterstützen (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty), eine Desktop-Benachrichtigung: ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  Es kennt jede Meldung, die Claude Code als Nutzungslimit zählt: die 5-Stunden-, Wochen-, einfachen und Modelllimits (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` für das Kontingent eines Modells), Nutzungsguthaben und Ausgabenlimits sowie eine Organisation, die das Konto gestoppt hat.

  Hat Claude Code die abgelehnte Anfrage in das Gespräch geschrieben, hält der Hook außerdem das volle Zeitfenster bis zu seinem Zurücksetzen fest ([der Limit-Eintrag](/features/usage#the-limit-record)), außer bei einem Stopp durch eine Organisation, den ihre Admins jederzeit aufheben können. Eine kurze Drosselung durch den Server oder das Ratenlimit eines API-Schlüssels ist kein Nutzungslimit und bleibt unbeachtet, ebenso unbeaufsichtigte Sitzungen (`-p`-Läufe und Hintergrundsitzungen). Ein Subagent oder Workflow-Agent, der an einem Nutzungslimit stoppt, zählt als Stopp seiner Sitzung, auch während das Hauptgespräch nur auf ihn wartet: Claude Code schreibt diese abgelehnte Anfrage in die eigene Datei des Agenten im Ordner `subagents` der Sitzung, und der Hook liest sie dort. Stoppen mehrere Agenten gleichzeitig, bekommt die Sitzung eine Notiz und eine Benachrichtigung.

  Claude Code 2.1.234 und neuer wartet außerdem in der Sitzung und macht nach einer 5-Stunden-Zurücksetzung von selbst weiter; die Benachrichtigung gibt dir die Wahl, stattdessen zu wechseln, und mit eingeschaltetem [`auto_switch`](#auto-switch) (Linux und macOS) wechselt ccam für dich, wenn ein anderes Konto freies Kontingent hat. Wenn kein anderes Konto übernehmen darf, weil kein anderes Konto zur selben Team- oder Enterprise-Organisation gehört oder keines ein anderer privater Tarif ist (mit eingeschaltetem `cross_org` darf es jedes andere Claude-Konto), endet die Benachrichtigung stattdessen mit `You have no other personal plans, so ccam cannot move this conversation.` oder `You have no other accounts in Acme, so ccam cannot move this conversation.`
- `/ccam:switch [account]` merkt das Gespräch vor und beendet die Runde, bevor eine Anfrage gestellt wird: Tippe `/exit` und führe dann `ccam next --continue` aus, um es mit dem Konto fortzusetzen, das du genannt hast (oder mit dem nächsten). In einer Sitzung, die noch kein Gespräch hat, meldet es, dass es nichts fortzusetzen gibt. Nur was du selbst tippst, führt es aus: Ab Plugin-Version 0.3.2 kann Claude `/ccam:switch` nicht ausführen, und ein Job, den Claude plant (CronCreate, ScheduleWakeup oder `.claude/scheduled_tasks.json`), gibt es als reinen Text an Claude weiter.

  Liegt das Konto, das du nennst, außerhalb der Organisation des Gesprächs (ein Platz einer anderen Team- oder Enterprise-Organisation, ein privater Tarif statt eines Platzes oder umgekehrt oder ein API-Schlüssel-Konto auf einer der beiden Seiten), fragt `ccam next --continue` zuerst nach, unabhängig davon, was `cross_org` sagt: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. Dasselbe fragt es vor einem Konto, das es selbst wählt, wenn du keines nennst, und das es mit `cross_org = true` aus einer anderen Organisation nehmen darf: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. Nur `y` oder `yes` verschiebt es; unter Linux und macOS werden Tasten, die vor der Frage getippt wurden, zuerst verworfen. Jede andere Antwort oder kein Terminal, auf dem gefragt werden kann, lässt das Gespräch und die Notiz, wo sie sind: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (nach einer Wahl wählt `ccam next --continue` erneut). Wie bei automatischen Wahlen passt ein Konto, dessen Organisation ccam nicht erkennen kann, zu jedem.

Die Notiz liegt in `~/.local/state/ccam/handoff/<account>/<session>.json`, eine Datei pro Sitzung, sodass zwei Sitzungen eines Kontos sich nie gegenseitig überschreiben. Sie wird 12 Stunden lang aufbewahrt und einmal verwendet, wenn die nächste Sitzung gerade starten soll; ein Start, der fehlschlägt, lässt sie für den nächsten Versuch liegen. Sie gehört zu dem Ordner, in dem die Sitzung gestartet wurde, auch nach einem `cd` in der Shell von Claude Code. `ccam rename` aktualisiert die Namen darin; wurde das Konto, das sie nennt, inzwischen entfernt, sagt ccam das und wählt wie gewohnt.

### Relay

Mit `relay = true` am Anfang von `config.toml` läuft ccam neben interaktiven Sitzungen weiter, statt sich selbst durch Claude Code zu ersetzen, und setzt nach `/ccam:switch` das Gespräch dieser Sitzung von selbst fort (auf eine Notiz, die eine andere Sitzung hinterlassen hat, reagiert es nicht): Der Hook beendet Claude Code, und ccam führt `ccam next --continue` mit der Notiz dieser Sitzung aus. Ohne genanntes Konto nimmt es das mit dem meisten Spielraum und berücksichtigt dabei das Modell des Gesprächs, und wenn `cross_org = true` zulässt, dass das ein Konto einer anderen Organisation ist, stellt es die Frage von oben, bevor es das Gespräch dorthin verschiebt. Unter Windows kann der Hook Claude Code nicht beenden: Tippe `/exit`, und ccam setzt das Gespräch fort. Ebenso wenig beendet der Hook Claude Code für ein Konto außerhalb der Organisation des Gesprächs: Er meldet `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`, und ccam stellt die Frage von oben, sobald du `/exit` getippt hast. Das kostet einen wartenden ccam-Prozess pro Sitzung, deshalb ist es standardmäßig aus.

Die weitergeführte Sitzung behält, was das Fortsetzen eines Gesprächs nicht wiederherstellt: Optionen wie `--add-dir`, `--mcp-config`, die Werkzeuglisten, `--settings`, `--plugin-dir` und die Systemprompt-Optionen, den Berechtigungsmodus des Gesprächs, die Effort-Stufe der Sitzung und den Agenten, mit dem das Gespräch läuft (siehe [`--continue`](#continue)). Claude Code selbst stellt `--model` und `--name` aus dem Gespräch wieder her. Alles andere (`--worktree`, `--debug`, `--remote-control`, Optionen, die ccam nicht kennt) wird in `ccam: not carried over: …` genannt. Der Wechsel stellt nie die Frage zum geteilten Verlauf, und was du während des Laufs getippt hast, wird verworfen; seine einzige Frage ist die vor einem Konto einer anderen Organisation.

Das Relay setzt das Gespräch nur fort, wenn Claude Code so endete, wie `/ccam:switch` es verlangt hat: nicht, wenn du zuerst `/exit` getippt, beim Beenden Ctrl+C gedrückt hast oder ccam ein Hangup oder ein anderes Signal bekommen hat, und nicht, wenn das Terminal weg ist. Die Notiz bleibt dann liegen, und `ccam next --continue` setzt das Gespräch später fort. Eine Sitzung, die unter dem Relay einer älteren Version von ccam gestartet wurde, oder eine unbeaufsichtigte (`-p`-Läufe und Hintergrundsitzungen) wird vom Hook nicht beendet; sie bekommt stattdessen die Meldung zu `/exit`. `claude remote-control` läuft ohne Relay.

### `auto_switch`

Mit `auto_switch = true` am Anfang von `config.toml` (Linux und macOS) verschiebt ccam eine Sitzung von selbst in ein anderes Konto, wenn Claude Code sie an einem Nutzungslimit stoppt: Der Hook des Plugins wählt ein anderes Konto mit freiem Kontingent und beendet Claude Code, und ccam, das neben der Sitzung weiterläuft, setzt das Gespräch dort fort, so wie das [Relay](#relay) nach `/ccam:switch`. Es ist standardmäßig aus und experimentell.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam verschiebt eine Sitzung nur, wenn all das zutrifft:

1. ccam hat sie in einem Terminal gestartet und läuft neben ihr weiter. `auto_switch` schaltet das [Relay](#relay) ein, daher trifft das auf die interaktiven Sitzungen zu, die du nach dem Einschalten startest. Eine Sitzung, die vorher oder unter einem älteren ccam gestartet wurde, bekommt eine Benachrichtigung, die mit ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` endet. Unbeaufsichtigte Sitzungen (`-p`-Läufe und Hintergrundsitzungen), `claude remote-control` und Claude Code, das ohne ccam gestartet wurde, werden nie verschoben.
2. Das Konto hat das [ccam-Plugin](/docs/install#claude-code-plugin) eingeschaltet und verwendet Claude Code 2.1.214 oder neuer, das dem Hook seine Prozess-ID übergibt (`CLAUDE_PID`).
3. Claude Code hat die abgelehnte Anfrage mit dem Zeitpunkt ihres Zurücksetzens in das Gespräch geschrieben, und zwar für das 5-Stunden-, Wochen-, Opus- oder Sonnet-Limit dieses Kontos. Ein Ausgabenlimit, aufgebrauchtes Nutzungsguthaben, ein Stopp durch eine Organisation, das Fable-Kontingent (`You've reached your Fable limit.`), eine Drosselung durch den Server und das Ratenlimit eines API-Schlüssels verschieben nie eine Sitzung. Das Limit eines Subagenten oder Workflow-Agenten zählt, auch während das Hauptgespräch nur auf ihn wartet; ist es ein wöchentliches Opus- oder Sonnet-Limit, wird die Sitzung nur verschoben, wenn das Gespräch selbst mit diesem Modell läuft.
4. ccam kennt den Tarif des Kontos, aus seiner `.claude.json` oder aus dem, was `ccam login` festgehalten hat. Andernfalls endet die Benachrichtigung mit `ccam does not know personal's plan, so auto_switch leaves this session.`
5. Die Sitzung ist nicht mit Remote Control verbunden. Als anderes Konto würde sie sich nicht wieder mit der App verbinden, von der aus du sie steuerst, daher bleibt sie, und die Benachrichtigung sagt das.
6. Das Limit wird in mehr als 15 Minuten zurückgesetzt. Bei einer kürzeren Wartezeit überlässt ccam die Sitzung Claude Code: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. Du hast seit dem Limit keinen Prompt getippt.
8. Ein anderes Konto hat freies Kontingent. ccam wählt wie [`--auto`](#auto) (`auto_max_usage`, `auto_strategy`, das Modell des Gesprächs), und zwar unter Konten derselben Art wie dem, das du verlässt: einem anderen Platz derselben Team- oder Enterprise-Organisation oder einem anderen privaten Tarif (Pro, Max). `cross_org` und `--cross-org` gelten hier nicht. Das Konto muss:
   - mit einer Anmeldung angemeldet sein, die Claude Code erneuern kann;
   - nicht beim selben Claude-Konto angemeldet sein;
   - keinen API-Schlüssel verwenden;
   - einen Tarif haben, den ccam kennt;
   - kein Nutzungsguthaben eingeschaltet haben, soweit ccams eigene Nutzungsabfrage das zeigt;
   - diesem Ordner in Claude Code vertraut haben;
   - wenn das Gespräch mit einem Agenten läuft (die Einstellung `agent` oder das Gespräch nennt einen), die Agenten des Kontos teilen, das du verlässt: Seine Ordner `agents` und `plugins` und seine `settings.json` müssen jeweils genau dieselben sein (über eine Gruppe) oder in beiden fehlen.

   Wenn die letzten drei Bedingungen jedes Konto mit freiem Kontingent ausschließen, nennt die Benachrichtigung sie: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` oder `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` Ein Zeitfenster mit einem [Limit-Eintrag](/features/usage#the-limit-record) zählt als voll, und ein Konto, zu dem eine andere Sitzung in den letzten 2 Minuten verschoben wurde, kommt nach den anderen. Die Wahl muss innerhalb von etwa 6 Sekunden fallen; dauern die Nutzungsabfragen länger, bleibt die Sitzung. Mit `usage_source = "claude"` entscheidet `auto_switch` allein anhand gespeicherter Zahlen, und ccam kann weder aus der Abfrage von Claude Code noch aus der Statuszeile erkennen, ob Nutzungsguthaben eingeschaltet ist, und nimmt daher an, dass es ausgeschaltet ist.

Dann beendet der Hook Claude Code, und ccam, das daneben weiterläuft, zeigt in Terminals, die das unterstützen, eine Desktop-Benachrichtigung wie `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` Es gibt die erste Zeile der Beispielausgabe oben aus und führt `ccam next --continue` für diese Sitzung mit dem gewählten Konto aus, das deren andere beiden Zeilen ausgibt und Claude Code mit den Optionen, dem Berechtigungsmodus, der Effort-Stufe und dem Agenten der Sitzung startet, so wie das [Relay](#relay). Die erste Nachricht der neuen Sitzung ist ein fester Text, das Einzige, was ccam in deinem Namen sendet (ein Wechsel nach `/ccam:switch` sendet nichts):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

Nach dem Limit eines Subagenten ist es dieser, weil die Subagenten, Workflows und Hintergrundbefehle mit der alten Sitzung gestoppt wurden:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

Erreicht auch dieses Konto sein Limit, passiert dort dasselbe. Die Limit-Einträge verhindern, dass die Kette vor dessen Zurücksetzen zu einem früheren Konto zurückkehrt; sind alle Konten voll, bleibt die Sitzung daher beim letzten.

Dass ein Gespräch mit einem Agenten läuft, erfährt der Hook aus der Einstellung `agent` oder aus der Zeile, die den Agenten nennt und die Claude Code von Zeit zu Zeit und beim Ende der Sitzung in das Gespräch schreibt. Nennt ihn bisher nur das `--agent` der Sitzung und teilt das gewählte Konto die Agenten nicht, startet `ccam next --continue` nichts und behält die Notiz, und eine zweite Desktop-Benachrichtigung sagt das: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

Hat kein anderes Konto freies Kontingent, wird nichts verschoben: Claude Code behandelt das Limit wie gewohnt, wartet in der Sitzung oder bietet seine Optionen an, und die Benachrichtigung sagt, welches Konto zuerst wieder frei wird, zum Beispiel `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` Bei Plätzen heißt es `None of your other accounts in Acme has room now`. Um einen begonnenen Wechsel abzubrechen, drücke Ctrl+C, während Claude Code sich beendet: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

Ein Wechsel hat seinen Preis. Claude Code startet neu, was einige Sekunden dauert, und die Hintergrund-Shells, Entwicklungsserver und Agenten, die Claude gestartet hat, enden mit der alten Sitzung. Das neue Konto liest mit seiner ersten Antwort das ganze Gespräch erneut, wie bei [`--continue`](#continue), und ccam sagt, wie groß es ist; ein Wechsel verwendet nie [`--brief`](#brief).

`auto_switch` wechselt nie, bevor Claude Code ein Limit meldet, überwacht deine Nutzung nie im Hintergrund und kopiert, erneuert oder übergibt nie eine Anmeldung. Unter Windows wird es ignoriert, und `ccam doctor` sagt das; dort verschieben `/ccam:switch` und `/exit` mit `relay = true` ein Gespräch weiterhin in ein anderes Konto.

**Bedingungen von Anthropic.** Auf Team- und Enterprise-Plätzen verschiebt `auto_switch` deine Arbeit auf einen anderen Platz derselben Organisation. Im Hilfe-Center von Anthropic ([Team-Tarife](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); ebenso für [Enterprise-Tarife](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)) heißt es: „Each seat on a Team plan is assigned to one named individual, and only that person may use it“ (jeder Platz in einem Team-Tarif ist einer namentlich benannten Person zugewiesen, und nur diese Person darf ihn nutzen) und „Users can't route their work through another user's account to use that user's remaining usage“ (Benutzer dürfen ihre Arbeit nicht über das Konto eines anderen Benutzers leiten, um dessen verbleibendes Nutzungskontingent zu verwenden).

Für private Tarife (Pro, Max) sagen die Bedingungen von Anthropic nicht, ob deine eigenen Konten an einem Limit automatisch füreinander einspringen dürfen. Die [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) verbieten, auf die Dienste „through automated or non-human means, whether through a bot, script, or otherwise“ (über automatisierte oder nicht-menschliche Mittel, ob durch einen Bot, ein Skript oder anderweitig) zuzugreifen, sofern Anthropic es nicht ausdrücklich erlaubt, und ein Wechsel sendet in deinem Namen eine der festen Nachrichten oben. Die [Usage Policy](https://www.anthropic.com/legal/aup), die am 12. November 2026 in Kraft tritt, nennt „Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products“ (Fähigkeiten, Beschränkungen oder Schutzmechanismen in unseren Produkten absichtlich umgehen oder beeinträchtigen) unter den Dingen, die sie verbietet.

ccam behauptet nicht, dass `auto_switch` diese Bedingungen einhält. Ob du es einschaltest, ist deine Entscheidung, ebenso das Risiko: Anthropic kann ein Konto oder den Zugang einer Team- oder Enterprise-Organisation sperren, wenn es die Bedingungen für verletzt hält.

Solange `auto_switch` eingeschaltet ist, warnt `ccam doctor` vor einem Konto ohne das ccam-Plugin und vor Claude Code-Versionen älter als 2.1.214 und wiederholt für Plätze die Platzregel von Anthropic ([Prüfungen](/features/doctor)).

## Am Limit erneut ausführen

`ccam next` ist für eine interaktive Sitzung gedacht, die angehalten hat. Ein nicht interaktiver Durchlauf, `claude -p` in einem Skript oder einem geplanten Job, kann nicht auf dich warten, deshalb lässt `--fallback` vor dem Kontonamen ccam für ihn wechseln:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` funktioniert nur zusammen mit `-p` oder `--print`; andernfalls hält ccam an und verweist auf `ccam next --auto --continue`.
2. Weitergeleitete Standardeingabe wird einmal gelesen, bis zu 10 MiB, und jedem Versuch übergeben. Wie Claude Code wartet ccam höchstens 3 Sekunden, bis sie beginnt (außer mit `--input-format stream-json`); bleibt eine Pipe offen, sendet aber nichts, bekommt der Lauf keine Eingabe, und ccam vermerkt das auf der Standardfehlerausgabe.
3. ccam startet `claude` als Kindprozess statt mit `exec`, hält seine Standardausgabe zurück und gibt die Standardfehlerausgabe sofort weiter.
4. Wenn `claude` mit 0 endet oder aus einem anderen Grund fehlschlägt, werden die zurückgehaltene Ausgabe und der Exit-Code unverändert weitergegeben. Ebenso bei einem Lauf, der durch ein Signal beendet wurde (ein Exit-Code über 128): Er sollte aufhören und wird nicht erneut ausgeführt.
5. Ein Fehlschlag zählt als Nutzungslimit, wenn die Meldung von Claude Code das sagt, mit einer der Nutzungslimit-Meldungen, die der [Hook](#inside-claude-code) kennt, auch innerhalb einer Ausgabe mit `--output-format json` oder `stream-json`, oder wenn die Nutzung des Kontos sein 5-Stunden-, Wochen- oder `--model`-Fenster bei 100 % zeigt (einschließlich eines Limit-Eintrags).
6. Dann wird die Ausgabe dieses Versuchs verworfen, und derselbe Befehl läuft mit einem anderen angemeldeten Konto: dem, das `--auto` wählen würde, danach Konten, deren Nutzung sich nicht abfragen lässt, in der Reihenfolge ihrer Nummern. Konten bei 100 % werden übersprungen, Konten mit einer Anmeldung per API-Schlüssel werden nie versucht, und jedes Konto wird höchstens einmal versucht.
7. Mit `-c`, `--continue`, `-r` oder `--resume` unter den Argumenten werden nur Konten in derselben Verlaufsgruppe versucht, damit das Gespräch gefunden wird.
8. Es werden nur Konten derselben Organisation versucht, wie bei `--auto` (`cross_org = true` hebt das auf).

Das nächste Konto wird nach denselben Regeln wie bei `--auto` gewählt (`auto_max_usage`, `auto_strategy`, das `--model`-Fenster). Wenn kein Konto mehr übrig ist, gibt ccam die Ausgabe des letzten Versuchs aus, meldet, dass kein angemeldetes Konto mehr Spielraum hat, und beendet sich mit Code 75. Mit `--wait` vor dem Kontonamen (`ccam --fallback --wait=2h work -p …`) wartet es stattdessen, bis das erste Konto zurückgesetzt wird, und kann dann jedes Konto erneut versuchen; es schläft immer zuerst mindestens einmal, damit ein Konto, das gerade sein Limit erreicht hat, nicht sofort wieder läuft. Die Dauer begrenzt die gesamte Wartezeit des Laufs, nicht jede einzelne Wartezeit. Mit Leerzeichen geschrieben (`--wait 30m`), wird es mit demselben Hinweis wie bei `ccam next` abgelehnt. Ein führendes `--` vor den Argumenten für Claude Code funktioniert wie bei jedem Start (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` gibt den Befehl des ersten Versuchs aus.

Interaktive Sitzungen werden nicht von selbst gewechselt, es sei denn, du schaltest [`auto_switch`](#auto-switch) ein. Mit `ccam statusline install` ergänzt die Statuszeile den Hinweis `→ ccam next --auto --continue`, sobald die 5-Stunden- oder Wochennutzung 90 % erreicht ([Nutzungslimits](/features/usage)).

## Sicherheit und Sonderfälle

- Konten ohne Anmeldung werden nie gewählt. Gibt es kein anderes angemeldetes Konto, hält einfaches `ccam next` mit `no logged-in account other than <name>` an.
- Mit `--continue` und keinem Gespräch im Ordner: `no account has a conversation in <folder> to continue`.
- Wenn die Organisationsregel kein Konto übrig lässt: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Mit `--auto` und keinem angemeldeten Konto: `no logged-in account to pick`. Wenn jedes Konto voll oder seine Nutzung unbekannt ist: `` no account with known usage and headroom (see `ccam usage`) ``. Beide enden mit Exit-Code 75.
- `--max-usage`, `--strategy`, `--wait` und `--cross-org` funktionieren nur mit `--auto`, und ein Wert außerhalb von 1-100 oder eine unbekannte Strategie wird abgelehnt.
- Ein Fork kopiert das Gespräch, nicht den übrigen Zustand des Kontos: Checkpoints aus `file-history/` des ursprünglichen Kontos kommen nicht mit, es sei denn, der Verlauf wird geteilt.
- `--auto` stellt auch mit `--dry-run` Netzwerkanfragen. Es liest nur und schreibt nie Token-Dateien.
- Das gewählte Konto startet über denselben Pfad wie `ccam <name>`, daher gelten dieselben Prüfungen, einschließlich der Warnung, wenn in diesem Ordner bereits eine andere Claude Code-Sitzung läuft, die denselben Verlauf teilt.
- `last-used` wird nur geschrieben, wenn Claude Code tatsächlich startet. `--dry-run` lässt es unverändert, daher rückt ein Probelauf die Rotation nicht weiter.
- `--fallback` hält die Standardausgabe zurück, bis ein Versuch endet, daher kommt eine fortlaufende Ausgabe wie `--output-format stream-json` am Ende in einem Stück an. Ein Client, der Eingaben als `stream-json` streamt und auf Antworten wartet, bevor er mehr sendet, kann `--fallback` nicht verwenden. Dateien, die der erste Versuch vor dem Erreichen des Limits geändert hat, bleiben geändert; der nächste Versuch beginnt den Prompt von vorn.

## Verwandte Themen

- [Nutzungslimits](/features/usage)
- [Nutzungslimits und Kontowechsel](/docs/usage-limits)
- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Was tun, wenn Claude Code sein Nutzungslimit erreicht](/guides/usage-limit-workflow)
