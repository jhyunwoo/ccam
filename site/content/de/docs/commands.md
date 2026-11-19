---
title: Befehlsreferenz
metaTitle: "ccam-Befehlsreferenz"
description: Alle ccam-Befehle und -Optionen, von ccam add und ccam <name> bis zu ccam usage, ccam next und ccam doctor.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: "28e1ad1dcbad"
---

Setze `--dry-run` vor einen Befehl, um zu sehen, was er tun würde, ohne es zu tun, zum Beispiel `ccam --dry-run work`. Das funktioniert beim Starten von Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) sowie mit `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` und `doctor --fix`. Befehle, die nur lesen, etwa `list`, `usage`, `share` und `ext`, nehmen die Option an und ändern ohnehin nichts. Die übrigen lehnen sie ab, statt zu laufen, auch `group rename`.

<!-- commands -->

## Claude Code starten

`ccam <name|number> [arguments]` startet Claude Code mit diesem Konto. ccam liest die Argumente nach dem Namen nicht; es gibt sie in derselben Reihenfolge an `claude` weiter, und Standardeingabe, Exit-Code und Signale verhalten sich genau wie bei `claude` selbst. Ein `--` nach dem Namen ist optional.

`ccam` mit Argumenten, aber ohne Kontonamen, etwa `ccam -p "summarise"`, verwendet das Konto, das für den aktuellen Ordner festgelegt ist, sonst das Standardkonto. Lege sie mit `ccam folder <name> [folder]` und `ccam use <name>` fest.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` vor dem Kontonamen ist für `-p`/`--print`-Durchläufe gedacht. Wenn das Konto sein Nutzungslimit erreicht, führt ccam denselben Befehl erneut mit dem angemeldeten Konto mit dem meisten Spielraum aus und gibt nur die Ausgabe des Durchlaufs aus, der fertig wurde. Siehe [Konten nach Ordner](/features/folder) und [Am Limit erneut ausführen](/features/next#run-again-at-the-limit).

**Team- und Enterprise-Plätze.** Die automatische Kontoauswahl (`ccam next`, `--fallback`, `ccam pick`, `n` im Dashboard und, sobald du es einschaltest, [`auto_switch`](/features/next#auto-switch)) wechselt auch zwischen Plätzen derselben Claude Team- oder Enterprise-Organisation, was gegen die Bedingungen von Anthropic verstoßen kann. Lies die [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) und kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt.

Ein Kurzbefehl wie `ccwork` startet sein Konto auf dieselbe Weise: `ccwork --resume` macht dasselbe wie `ccam work --resume`. Siehe [Kurzbefehle](/docs/commands#shortcuts) weiter unten.

## Konten hinzufügen und entfernen

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` ohne Ordner sucht nach `~/.claude` und `~/.claude-*` und fragt bei jedem nach; ein Ordner, bei dem es fehlschlägt, wird als `Skipped <path>: <reason>` gemeldet, und die übrigen werden registriert. Ist der Name eines Ordners schon vergeben, importiere diesen Ordner einzeln mit `ccam import <folder> --name <name>`. Importierte Ordner werden dort registriert, wo sie liegen; ccam verschiebt oder kopiert nie eine Anmeldung.

`ccam login --expired` meldet auch Konten an, deren Anmeldung Claude Code gelöscht hat. `ccam remove` meldet ein verwaltetes Konto (eines, dessen Ordner ccam angelegt hat, im Gegensatz zu einem mit `ccam import` registrierten) ab, bevor sein Ordner in den Papierkorb kommt, außer ein anderes Konto hat dieselbe Anmeldung. Das Entfernen eines importierten Kontos nimmt es aus seinen Gruppen und hinterlässt in seinem Ordner eine Kopie jedes geteilten Elements. Unter macOS behält `ccam rename` den Namen des Ordners bei, da der Schlüsselbund die Anmeldung unter dem Pfad des Ordners speichert.

`ccam add --api` fügt ein Konto hinzu, das einen API-Schlüssel statt einer Claude-Anmeldung verwendet; ccam liest den Schlüssel bei jedem Start und speichert ihn nie. Siehe [API-Schlüssel-Konten](/features/api).

`ccam rename`, `ccam remove` und die `ccam group`-Befehle ändern ein Konto nicht, solange Claude Code mit diesem Konto läuft, egal wie es gestartet wurde (einschließlich Hintergrundsitzungen). Beende zuerst die Sitzung. Sie lehnen außerdem ab, solange eine Gruppenänderung nicht abgeschlossen ist: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` gibt jedes Konto mit Slot-Nummer, Kurzbefehl, E-Mail-Adresse und Ordner aus. Das Array `groups` listet die Gruppen des Kontos auf, und das Objekt `share` ordnet jedem Element (`mcp`, `settings`, `memory`, `extensions`, `history`) die Gruppe zu, aus der es kommt, oder `local`.

## Kurzbefehle

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Ein Kurzbefehl ist ein Befehl, den du selbst benennst und der Claude Code mit einem bestimmten Konto startet. Wie `ccam work` gibt er jedes Argument unverändert an Claude Code weiter, sodass `ccwork --resume` und `ccwork -p "hi"` wie erwartet funktionieren. Er ist ein Link auf ccam in `~/.local/bin`.

Kurzbefehle sind optional. Der Assistent von `ccam add` fragt, ob einer angelegt werden soll, `ccam add <name> --shortcut <command>` legt einen ohne Nachfrage an, und `ccam list` zeigt den Kurzbefehl jedes Kontos oder `-`, wenn es keinen hat.

Namen bestehen aus den Zeichen a-z, 0-9, `.`, `-` und `_`. Ein Name darf kein ccam-Befehl sein, nicht der Kurzbefehl eines anderen Kontos, nicht der Name eines anderen Programms in deinem `PATH` und keine Datei in `~/.local/bin`, die ccam nicht angelegt hat, etwa deine eigene Kopie von ccam oder ein Link darauf. Ein Kurzbefehl, der von Hand in `config.toml` bearbeitet wurde und gegen diese Regeln verstößt, `claude` oder `ccam` heißt oder zwei Konten zugewiesen ist, lässt das Laden von `config.toml` fehlschlagen, wobei das Konto genannt wird.

`ccam rename` behält den Namen des Kurzbefehls bei, und `ccam remove` löscht den Kurzbefehl zusammen mit dem Konto. `ccam shims install` erstellt die Links aller Konten neu, zum Beispiel nachdem du ccam verschoben hast; eine Datei, die einem Kurzbefehl im Weg steht, hält die anderen nicht auf.

## Teilen

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Eine Gruppe teilt Elemente zwischen ihren Konten; ein Konto kann in mehreren Gruppen sein, aber jedes Element kommt nur aus einer Gruppe. Siehe [Teilen mit Gruppen](/features/share).

## Eine vorhandene Einrichtung kopieren

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` kopiert `settings.json`, `CLAUDE.md`, Erweiterungen und Benutzer-MCP-Server aus einem Claude Code-Ordner in eine Gruppe, sodass jedes Konto, das diese Elemente teilt, sie bekommt. Der Ordner wird nur gelesen.

- `--settings`, `--memory`, `--extensions`, `--mcp` und `--history` kopieren nur die genannten Elemente. Ohne sie kannst du in einem Terminal in einer Checkliste auswählen, in der alles außer dem Gesprächsverlauf angehakt ist; anderswo wird alles außer dem Gesprächsverlauf kopiert.
- `--mcp` funktioniert auch, wenn `~/.claude` leer ist oder fehlt, `~/.claude.json` aber Server enthält.
- `-g` oder `--group` wählt die Gruppe.
- `-y` oder `--yes` kopiert ohne Nachfrage und behält die Kopie der Gruppe, wenn eine Datei abweicht.

MCP-Anmeldungen werden nicht kopiert, da eine kopierte Anmeldung die anderen Inhaber abmeldet; der Plan nennt die Konten, die sich mit `/mcp` anmelden. Siehe [Eine vorhandene Einrichtung kopieren](/docs/sharing#copy-an-existing-setup).

## MCP-Server

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` verschiebt die Server, die nur dieses Konto hat, in seine Gruppe, einschließlich derer, die es nur für einen Ordner hat (einfaches `claude mcp add`); diese werden dann in jedem Ordner geladen. Ein Server, dessen Namen die Gruppe bereits für eine andere Definition verwendet, wird übersprungen; mit `--force` ersetzt du die Kopie der Gruppe, die ccam vorher sichert.

## Nutzung und Wechseln

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam sends no token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --continue --brief    # a new conversation from a summary of the last one
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agent "…"                  # from Claude Code: a subagent that runs as another account
ccam agent --resume <id> "…"    # follow that agent up
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
ccam ccusage claude daily       # Claude Code's tokens and costs, every account (ccusage)
ccam ccusage -a work blocks --active   # one account's current 5-hour block
```

- `next --auto`, `--fallback` und `ccam pick` wählen auf dieselbe Weise: zuerst unter `auto_max_usage` (`config.toml`, Standard 100), nach `auto_strategy` (`headroom`, `expiring` oder `balance`), unter Berücksichtigung des Wochenfensters von `--model`. Sie wählen nie ein Konto bei 100 % (ein Fenster, das Claude Code als voll gemeldet hat, zählt bis zu seinem Zurücksetzen als 100 %) oder eine Anmeldung per API-Schlüssel und bleiben innerhalb der Organisation des Kontos, von dem sie ausgehen (`--cross-org` oder `cross_org = true`, um weiter zu suchen).
- `--wait` nimmt seinen Wert nach `=` (`--wait=2h`); mit Leerzeichen geschrieben wird es mit einem Hinweis abgelehnt. Wenn keines gewählt werden kann, beenden sie sich mit Code 75 oder warten mit `--wait`, bis das erste Konto zurückgesetzt wird.
- `ccam next --continue` nimmt ein Gespräch auch zwischen Konten mit, die keinen Verlauf teilen: Claude Code kopiert es dann in das neue Konto. Mit `--brief` beginnt das neue Konto stattdessen ein neues Gespräch, ausgehend von einer Zusammenfassung des letzten, die ccam ohne Modell schreibt ([`--brief`](/features/next#brief)).
- Anders als `ccam <name>` liest `ccam next` seine eigenen Optionen, daher gehören die Optionen von Claude Code hinter `--`; ohne es bricht ccam mit `unknown flag` ab.
- `ccam agent` führt eine Aufgabe mit `claude -p` in diesem Ordner mit einem anderen Konto als dem aus, aus dem es aufgerufen wird, im Berechtigungsmodus und mit der Effort-Stufe der aufrufenden Sitzung, und gibt die letzte Antwort von Claude aus. Am Limit dieses Kontos geht das Gespräch mit einem anderen weiter ([Subagenten auf anderen Konten](/features/agent)).
- In Claude Code übergibt `/ccam:switch [account]` aus dem Plugin ein Gespräch, und `ccam hook` ist das, was seine Hooks ausführen.
- Mit `auto_switch = true` in `config.toml` (Linux und macOS; standardmäßig aus) beendet der Limit-Hook außerdem eine interaktive Sitzung, die ccam gestartet hat, und setzt sie in einem anderen Konto mit freiem Kontingent fort, aber nur in einem derselben Art: einem anderen Platz derselben Organisation oder einem anderen privaten Tarif ([`auto_switch`](/features/next#auto-switch)).

Siehe [Zum nächsten Konto wechseln](/features/next), [Ein Konto für Skripte wählen](/features/pick), [Subagenten auf anderen Konten](/features/agent) und [Hintergrundsitzungen aller Konten](/features/agents).

## Wartung

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam uninstall       # remove every account, ccam's folders and ccam itself
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` prüft die SHA-256-Prüfsumme des Releases, bevor es ccam ersetzt, und `--version v0.2.0` installiert ein bestimmtes Release. Eine aus dem `.deb`-Paket installierte Kopie von ccam aktualisierst du stattdessen mit apt. Siehe [Aktualisieren](/docs/install#update).

`ccam uninstall` listet auf, was es entfernt, und fragt zuerst: jedes Konto so wie `ccam remove`, die Kurzbefehle, die Ordner von ccam und den Befehl `ccam`. Ein importierter Ordner bleibt, wo er ist. Siehe [Deinstallieren](/docs/install#uninstall).

## Dateien

| Pfad | Inhalt |
| --- | --- |
| `~/.config/ccam/config.toml` | Konten, Gruppen und Einstellungen (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). Einstellungen stehen am Anfang, über dem ersten `[[group]]`, `[[profile]]` oder `[[folder]]`: Eine Einstellung unter einer Tabelle gehört zu dieser Tabelle, daher warnt ccam bei jedem Start und in `ccam doctor` davor, und Befehle, die die Konfiguration speichern, lehnen ab, bis sie verschoben ist. Kommentare bleiben nicht erhalten, wenn ccam die Datei speichert |
| `~/.local/share/ccam/profiles/<name>/` | Das `CLAUDE_CONFIG_DIR` jedes Kontos |
| `~/.local/share/ccam/groups/<group>/` | Geteilte Einstellungen, MCP-Server, Erweiterungen und Verlauf |
| `~/.local/state/ccam/` | Nutzungscache, Limit-Einträge, Sperren, Notizen und weiterer Zustand (siehe die Liste unter der Tabelle) |
| `~/.local/bin/<shortcut>` | Kurzbefehle wie `ccwork`, Links auf ccam |

Der Zustandsordner enthält:

- den Nutzungscache und die Limit-Einträge (`usage/<name>.json`, `usage/<name>.limit.json`);
- die Zeitpunkte der letzten Verwendung;
- Sitzungsmarker und Sperren (`locks/`);
- die Kurzbefehle, die ccam angelegt hat (`shortcuts.json`);
- eine nicht abgeschlossene Gruppenänderung (`group-change.json`);
- die Notizen, die `/ccam:switch` und der Limit-Hook hinterlassen (`handoff/<account>/<session>.json`);
- wann `ccam pick` oder `auto_switch` jedes Konto zuletzt gewählt hat (`picks/<name>`; eine Wahl zählt 2 Minuten lang);
- die Zusammenfassungen von `--brief`, 7 Tage lang aufbewahrt (`brief/<session>.md`).

Wenn du `XDG_CONFIG_HOME`, `XDG_DATA_HOME` oder `XDG_STATE_HOME` setzt, ändert sich, wo diese Ordner liegen. Unter Windows ist `~` gleich `%USERPROFILE%`, und Kurzbefehle sind `<shortcut>.exe`.
