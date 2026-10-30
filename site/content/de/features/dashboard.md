---
title: Das Dashboard
description: ccam öffnet im Terminal ein Vollbild-Dashboard mit Konten, Gruppen, Ordnerregeln, Erweiterungen und Prüfungen; jede Änderung zeigt zuerst ihren Plan.
updated: 2026-10-30
order: 15
translation: machine
sourceHash: "d65f1b5bb0c8"
---

`ccam` ohne Argumente öffnet ein Vollbild-Dashboard, wenn Standardeingabe und -ausgabe ein Terminal sind; andernfalls gibt es `ccam list` aus. Das Dashboard zeigt jedes Konto mit seiner Nutzung und seinen Gruppen und erledigt, was die Befehle tun: starten, hinzufügen, anmelden, Gruppen ändern, Ordnerregeln setzen, Hintergrundsitzungen auflisten, prüfen und reparieren. Ein Start oder das Verbinden mit einer Hintergrundsitzung schließt es; beim Hinzufügen eines Kontos oder bei der Anmeldung gibt es das Terminal kurz ab und kehrt dann zum Dashboard zurück.

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

Die Zeile unter den Tab-Namen sagt in einem Satz, wofür der aktuelle Tab da ist. In einem breiteren Terminal zeigt der Tab Accounts außerdem die E-Mail-Adresse jedes Kontos. Auf jedem Tab zeigt das Ende der Zeile mit den Tab-Namen, ob [`auto_switch`](/features/next#auto-switch) eingeschaltet ist (Linux und macOS; Windows ignoriert die Einstellung). Ist es eingeschaltet, listet der Tab Health auf, was dafür noch fehlt: Claude Code 2.1.214 oder neuer und das eingeschaltete ccam-Plugin in jedem Konto, das sich mit einem Claude-Abo anmeldet.

Bei jedem Öffnen fragt das Dashboard im Hintergrund bei GitHub nach, ob es ein neueres ccam-Release gibt; das erste Bild wartet nie darauf, und wenn du offline bist, wird nichts angezeigt. Gibt es eines, empfiehlt eine Zeile unter dem Titel das Update:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` fragt zuerst nach und aktualisiert dann so wie [`ccam update`](/features/update): Das Release wird gegen seine SHA-256-Prüfsummen geprüft und ersetzt ccam in einem Schritt. Das Dashboard, in dem du gerade bist, ist noch die alte Version, also beende es und starte `ccam` neu. Wenn das Update hier nicht möglich ist, etwa bei einer aus dem `.deb`-Paket installierten Kopie von ccam, zeigt das Dashboard denselben Hinweis wie `ccam update` und behält das Angebot bei.

Das erste Bild entsteht allein aus `config.toml`, in deutlich unter 100 ms; die Nutzung erscheint, sobald das jeweilige Konto antwortet. Jedes Limit hat seinen eigenen Balken: das 5-Stunden-Fenster, das Wochenfenster für alle Modelle und das eigene Wochenfenster von Fable; `-` bedeutet, dass das Konto kein solches Fenster hat. Ein Prozentwert wird bei 80 % gelb und bei 100 % rot. Schlägt eine Abfrage fehl, werden zwischengespeicherte Werte blass dargestellt. Ein Zeitfenster, das Claude Code als voll gemeldet hat, zeigt bis zu seinem Zurücksetzen 100 %. Ein abgelaufenes Zugriffstoken zeigt ein blasses `⚠ token expired` (Claude Code erneuert es beim nächsten Start), während `⚠ login expired` bedeutet, dass das Konto `ccam login` braucht.

## Tabs

| Tab | Was er zeigt | Tasten |
| --- | --- | --- |
| Accounts | Nummer, Name, E-Mail-Adresse, Balken für 5 Stunden, Woche und Fable, letzte Verwendung, Gruppen; `✔` das Standardkonto, `▸` das Konto dieses Ordners, `•` läuft (auch Claude Code, das ohne ccam gestartet wurde: einfaches `claude`, eine IDE, Hintergrundsitzungen und ihr Daemon) | siehe unten |
| Agents | was [`ccam agents`](/features/agents) auflistet: die Hintergrundsitzungen und offenen Sitzungen jedes Kontos, mit Zustand, Ordner und Startzeit (siehe unten) | `enter` verbinden (wie `ccam attach`), `x` stoppen, `r` neu laden |
| Groups | jede Gruppe, was sie teilt und ihre Konten, sowie die Konten ohne Gruppe | `enter` bearbeiten, `n` neu, `e` umbenennen, `x` entfernen |
| Folders | die Regeln von [`ccam folder`](/features/folder) und die, die hier gilt, selbst wenn du den Ordner über einen symbolischen Link erreicht hast (neue Regeln werden mit aufgelösten Links gespeichert) | `n` ein Konto hier verwenden, `enter` ändern, `x` entfernen |
| Extensions | was [`ccam ext`](/features/ext) zeigt, jede Erweiterung einmal mit den Gruppen, aus denen sie kommt, `off:` die Konten, in denen ein Plugin nicht eingeschaltet ist, die eigenen Server des Kontos, die nur für diesen Ordner gelten, `This folder's project` und Probleme, die zu keinem Konto gehören | `r` neu laden |
| Health | was [`ccam doctor`](/features/doctor) meldet, als Kasten pro Konto (siehe unten) | `f` den ausgewählten Kasten reparieren, `F` alles reparieren (`doctor --fix`), `↑`/`↓` zwischen Kästen wechseln |

`←`/`→` oder `Tab` wechseln den Tab, `?` listet die Tasten des aktuellen Tabs auf, `r` aktualisiert die Nutzung unter Umgehung des kurzen Caches und lädt die Tabs Agents, Extensions und Health neu, und `q` oder `Esc` beendet.

## Tasten auf dem Tab Accounts

| Taste | Wirkung |
| --- | --- |
| `Enter` | Claude Code als das ausgewählte Konto starten |
| `1`–`9` | Zum Konto mit dieser Nummer gehen; tippe `1` `2` schnell hintereinander für 12. Nur `Enter` startet |
| `↑`/`↓`, `k`/`j` | Bewegen |
| `n` | Das Konto starten, das [`ccam next --auto`](/features/next) wählen würde |
| `o` | Mit Optionen starten: `--continue`, `--resume`, andere Argumente oder ein Probelauf |
| `u` | Nutzungsdetails: 5-Stunden-, Wochen- und Modellfenster mit ihren Zurücksetzungen und der Herkunft der Zahlen |
| `a` | Ein Konto hinzufügen (der Assistent unten) |
| `l` | Erneut anmelden |
| `s` | Die Gruppen des Kontos wählen |
| `d` | Es zum Standardkonto machen |
| `f` | Es im aktuellen Ordner verwenden |
| `e` | Es umbenennen |
| `c` | Seinen Kurzbefehl setzen oder entfernen |
| `x` | Es entfernen (fragt zuerst nach, mit No vorausgewählt) |

**Team- und Enterprise-Plätze.** Die automatische Kontoauswahl (`ccam next`, `--fallback`, `ccam pick`, `n` im Dashboard und, sobald du es einschaltest, [`auto_switch`](/features/next#auto-switch)) wechselt auch zwischen Plätzen derselben Claude Team- oder Enterprise-Organisation, was gegen die Bedingungen von Anthropic verstoßen kann. Lies die [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) und kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt.

## Der Tab Health

Der Tab Health führt beim Öffnen die Prüfungen von [`ccam doctor`](/features/doctor) aus und zeigt dieselbe Ansicht: eine erste Zeile mit einem Zeichen pro Konto, dann einen Kasten pro Konto. Jede Zeile in einem Kasten ist ein geprüfter Bereich, zuerst `Folder`, `Login` und `Share links`, markiert mit `✔` in Ordnung, `ℹ` Information, `⚠` eine Warnung oder `✘` ein Fehler, sodass du auf einen Blick siehst, welches Konto wo ein Problem hat. Befunde, die zu keinem Konto gehören, kommen in einen Kasten `General`.

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

`↑`/`↓` wechseln von Kasten zu Kasten, und der ausgewählte Kasten hat einen Rahmen in der Akzentfarbe. `f` repariert nur den ausgewählten Kasten, so wie `ccam doctor --fix` es tun würde: im Kasten eines Kontos nur dieses Konto, im Kasten `General` die Befunde, die zu keinem Konto gehören (Plugin-Pfade und Kurzbefehle); alle anderen Konten werden nur erneut geprüft. Vorher kommt eine Rückfrage (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), und ein Kasten, in dem es nichts zu reparieren gibt, meldet nur `nothing to repair in <name>`. `F` repariert alles auf einmal, wie `ccam doctor --fix`, und `r` prüft erneut.

## Der Tab Agents

Der Tab Agents zeigt, was [`ccam agents`](/features/agents) ausgibt: die Hintergrundsitzungen und offenen Sitzungen jedes Kontos, jeweils mit dem Konto, zu dem sie gehören. Er fragt das Claude Code jedes Kontos gleichzeitig ab (`claude agents --json`, jeweils mit einem Timeout von 15 Sekunden), wenn du den Tab zum ersten Mal öffnest, bei `r`, nach einem Stopp und nach einer Änderung der Konfiguration, und zeigt `Asking every account's Claude Code for its sessions…`, bis alle geantwortet haben.

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

`State` ist der Zustand einer Hintergrundsitzung und das, was ihr Prozess gerade tut, oder `open` für eine interaktive Sitzung, wie in `ccam agents`. Lässt sich ein Konto nicht abfragen, zeigt eine `✘`-Zeile unter der Liste, was Claude Code gesagt hat, wie bei `lab` oben, dessen Agent-Ansicht ausgeschaltet ist. Beendete Hintergrundsitzungen werden ausgelassen (`ccam agents --all` listet sie auf), und wenn alle Konten geantwortet haben und keines eine Sitzung hat, zeigt der Tab ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` In einem schmalen Terminal lässt der Tab zuerst `Started` weg, dann `Name`, dann `Folder`, und kürzt lange Ordnerpfade von links.

`Enter` auf einer Hintergrundsitzung schließt das Dashboard und öffnet die Sitzung in diesem Terminal als ihr Konto, so wie `ccam attach <id>`; das Dashboard kommt nicht zurück (in Claude Code führt `←` zu seiner Agent-Ansicht und `Ctrl+Z` zur Shell). Eine offene Sitzung hat keine ID, daher sagen `Enter` und `x` dort nur, dass eine Hintergrundsitzung nötig ist. `x` fragt vorher nach, mit `No, keep it` vorausgewählt: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` Bei Ja führt ccam `claude stop <id>` als dieses Konto aus, meldet ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` oder warum es nicht ging, und liest die Liste neu ein.

## Gruppen ändern

Auf dem Tab Groups öffnet `Enter` den Editor: links die Elemente, rechts die Konten, `Tab` zum Wechseln, `Space` zum Umschalten. Ein Element, das eines der Konten bereits von einer anderen Gruppe bekommt, ist mit dem Grund ausgegraut, ebenso ein Konto, das ein Element doppelt bekäme; ein Element kann nur aus einer Gruppe kommen. `Enter` zeigt den Plan:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Weicht eine Datei von der Kopie der Gruppe ab, fragt ccam zuerst, was zu tun ist (die der Gruppe behalten, meine verwenden oder zusammenführen). Dann läuft die Änderung im Hintergrund mit demselben Code wie [`ccam group`](/features/share): Die Konfiguration wird zuerst gespeichert, die Daten jedes Kontos werden der Reihe nach verschoben, und die Originale kommen nach `.ccam-backup/`; ein Konto, das das Teilen von MCP verlässt, behält die Server der Gruppe in seiner `.claude.json`.

Eine Änderung, die die Daten eines laufenden Kontos verschieben würde, wird mit dem Namen des Kontos abgelehnt, selbst wenn das Konto gestartet wurde, während der Plan offen war, und ein Plan, der veraltet ist, weil ein anderer ccam-Prozess die Gruppen geändert hat, wird mit `the groups changed while this plan waited; nothing was moved, make the change again` abgelehnt. Während ein anderer ccam-Prozess die Gruppen ändert, melden eine neue Änderung, ein Umbenennen oder Entfernen `another ccam process is changing the groups right now`; nach einer Änderung, die auf halbem Weg stehen geblieben ist, warten sie, bis `F` auf dem Tab Health (oder `ccam doctor --fix`) sie abgeschlossen hat.

## Ein Konto hinzufügen

`a` öffnet den Assistenten, den auch `ccam add` in einem Terminal verwendet:

1. Der Kontoname.
2. Wie es sich anmeldet: jetzt mit einem Claude-Abo (im Browser), mit einem [API-Schlüssel](/features/api) oder später mit einem Claude-Abo.
3. Bei einem API-Schlüssel: woher ccam den Schlüssel bei jedem Start liest (`env:NAME` oder `command:COMMAND`), der Endpunkt, wie der Schlüssel gesendet wird, und das Modell.
4. Die Gruppen zum Beitreten, mit angehaktem `default` (gibt es noch keine Gruppen, was die Gruppe `default` teilen wird).
5. Ein Kurzbefehl, wenn du einen willst.
6. Eine Zusammenfassung zum Bestätigen.

`Esc` geht einen Schritt zurück und schließt den Assistenten im ersten Schritt.

## Ein Konto entfernen

`x` fragt nach, bevor irgendetwas passiert, mit No vorausgewählt, und sagt, was das Entfernen bewirkt: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` Bei einem Konto mit API-Schlüssel gibt es keine Anmeldung zum Widerrufen, und wenn ein anderes Konto dieselbe Anmeldung hat, wird es nicht abgemeldet. Ein importiertes Konto behält seinen Ordner: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## Farben

Die Farben folgen den Rollen der Terminal-Oberfläche von Claude Code, sodass ccam daneben vertraut aussieht; Gruppennamen erhalten eine Farbe, die sich aus ihrem Namen ergibt. ccam fragt das Terminal nach seinem Hintergrund und wählt ein helles oder dunkles Design. Du kannst eines selbst am Anfang von `config.toml` festlegen, vor dem ersten `[[group]]`, `[[profile]]` oder `[[folder]]`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Die Designs mit `-daltonized` sind für Rot-Grün-Farbenblindheit gedacht, und die Designs mit `-ansi` verwenden die 16 eigenen Farben des Terminals. `NO_COLOR` schaltet die Farben aus.

## Sicherheit und Sonderfälle

- Das Dashboard ruft dieselben Funktionen auf wie die Befehle, daher gilt jede Regel: Pläne, bevor Daten bewegt werden, Ablehnungen bei laufenden Konten, Sicherungen und Anmeldetokens, die ein Konto nie verlassen.
- Das Dashboard hat keinen Probelauf: `ccam --dry-run` ohne Namen gibt aus, wie das Konto dieses Ordners starten würde, und `o` bei einem Konto bietet einen Probelauf genau dieses einen Starts an.
- Das Terminal wird wiederhergestellt, bevor Claude Code startet oder sich die Anmeldung öffnet.
- Bei einem Bildschirm, der für einen Dialog zu klein ist, fallen zuerst die Zeilen des Tabs weg; Zeilen werden mit `…` abgeschnitten statt umbrochen. Auf einem niedrigen Bildschirm fallen als Nächstes Leerzeilen und die Tab-Beschreibung weg, sodass die Zeile mit den Tastenhinweisen sichtbar bleibt. Unter 40 Spalten oder 8 Zeilen zeigt das Dashboard nur `Too small: ccam needs 40x8 (now 36x7)`.
- `Ctrl+C`, während eine Änderung, ein Entfernen, ein Update oder eine Reparatur angewendet wird, wartet darauf: Die Statuszeile sagt `ccam quits when it is done.`
- Nachdem `U` ccam aktualisiert hat, kommt das Angebot nicht wieder, wenn das Dashboard vom Hinzufügen eines Kontos oder einer Anmeldung zurückkehrt.

## Verwandte Themen

- [Teilen mit Gruppen](/features/share)
- [Claude Code starten](/features/launch)
- [Konten verwalten](/features/accounts)
- [Nutzungslimits](/features/usage)
