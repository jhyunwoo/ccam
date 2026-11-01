---
title: Claude Code starten
description: Claude Code mit ccam <name|number>, ccam run oder dem Dashboard als beliebiges Konto starten; Argumente gehen unverändert an claude, --dry-run zeigt den Befehl.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: "11df08608a1e"
---

`ccam work` startet Claude Code als das Konto `work`, und alles nach dem Namen geht unverändert an `claude`. Verwende es überall dort, wo du `claude` tippen würdest: in einem Terminal, in einem Skript oder mit weitergeleiteter Eingabe.

## Verwendung

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Gib den Namen eines Kontos oder seine Nummer aus `ccam list` an. ccam liest die Argumente nach dem Namen nicht, daher zeigt `ccam work --help` die Hilfe von Claude Code. Ein `--` nach dem Namen ist optional. Standardeingabe, Exit-Code und Signale verhalten sich wie bei `claude`.

Beginnt das erste Argument mit `-` und ist keine von ccams eigenen Optionen, wählt ccam das Konto für dich und gibt alle Argumente an Claude Code weiter. Es verwendet das Konto, das mit [`ccam folder`](/features/folder) für den aktuellen Ordner festgelegt ist, sonst das Standardkonto: das mit `ccam use` festgelegte Konto, oder das Konto mit der niedrigsten Nummer, wenn du keines festgelegt hast.

### Das Dashboard

`ccam` ohne Argumente öffnet das [Dashboard](/features/dashboard), wenn Standardeingabe und -ausgabe ein Terminal sind, und gibt sonst `ccam list` aus. Gehe zu einem Konto (mit den Pfeiltasten oder indem du seine Nummer tippst) und drücke `Enter`, um es zu starten; `n` startet das Konto mit dem meisten Spielraum, und `o` bietet `--continue`, `--resume` oder andere Argumente an. Das Dashboard schließt sich, bevor Claude Code startet, sodass das Terminal zuerst wiederhergestellt wird.

**Team- und Enterprise-Plätze.** Die automatische Kontoauswahl (`ccam next`, `--fallback`, `ccam pick`, `n` im Dashboard und, sobald du es einschaltest, [`auto_switch`](/features/next#auto-switch)) wechselt auch zwischen Plätzen derselben Claude Team- oder Enterprise-Organisation, was gegen die Bedingungen von Anthropic verstoßen kann. Lies die [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) und kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt.

Wenn in der Shell `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` oder `ANTHROPIC_AUTH_TOKEN` gesetzt ist, kann Claude Code diesen Schlüssel oder dieses Token statt der eigenen Anmeldung des Kontos verwenden, daher warnt ccam dich vor dem Start auf der Standardfehlerausgabe (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Es liest die Variablen nur; ein [API-Schlüssel-Konto](/features/api) entfernt sie stattdessen.

Standardmäßig ersetzt sich ccam selbst durch Claude Code, sodass nichts anderes weiterläuft. Mit `relay = true` (oder, unter Linux und macOS, [`auto_switch = true`](/features/next#auto-switch)) am Anfang von `config.toml`, vor dem ersten `[[group]]`, `[[profile]]` oder `[[folder]]`, startet ccam interaktive Sitzungen stattdessen als Kindprozess und läuft neben ihnen weiter. So können `/ccam:switch` oder, an einem Nutzungslimit, `auto_switch` das Gespräch von selbst in einem anderen Konto fortsetzen, mit den Optionen, dem Berechtigungsmodus, der Effort-Stufe und dem Agenten der Sitzung; bevor `/ccam:switch` es in ein Konto einer anderen Organisation verschiebt, fragt ccam dich ([Details](/features/next#relay)).

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--dry-run` | Umgebung und Befehl ausgeben, die ccam ausführen würde, dann beenden | aus |
| `--fallback` | Für `-p`/`--print`-Durchläufe: Wenn das Konto sein Nutzungslimit erreicht, denselben Befehl erneut mit dem angemeldeten Konto mit dem meisten Spielraum ausführen ([Details](/features/next#run-again-at-the-limit)) | aus |
| `-v`, `--version` | `ccam <version>` ausgeben (nur als erstes Argument) | |
| `-h`, `--help` | Die Hilfe von ccam, als `ccam --help` oder `ccam run --help` | |

Regeln für Argumente:

- ccam liest seine eigenen Optionen nur vor dem Kontonamen. Nach dem Namen gehört jedes Argument Claude Code, auch `--dry-run` und `--help`.
- `--dry-run` funktioniert als `ccam --dry-run work`, `ccam run --dry-run work` und `ccam --dry-run` (für das Konto des Ordners oder das Standardkonto). `--dry-run` und `--fallback` können in beliebiger Reihenfolge stehen.
- `ccam run <name|number> [claude args]` ist die Langform von `ccam <name>`. `ccam run` ohne Namen verwendet das Konto des Ordners oder das Standardkonto.
- Kontonamen dürfen keine ccam-Befehlsnamen wie `add`, `list`, `run` oder `usage` sein, daher ist `ccam <word>` nie mehrdeutig.

Der Probelauf gibt eine Zeile aus, die du in eine Shell einfügen könntest, oder unter Windows in PowerShell. Sie beginnt mit dem, was der Start an der Umgebung ändert: immer `CLAUDE_CONFIG_DIR` und `CCAM_PROFILE`, `env -u <name>` für eine Variable, die er entfernt (etwa ein geerbtes `CCAM_RELAY`), und bei einem [API-Schlüssel-Konto](/features/api) dessen Endpunkt, sein Modell und woher der Schlüssel kommt (`ANTHROPIC_API_KEY="$WORK_KEY"`), nie den Schlüssel selbst:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

Die PowerShell-Zeile setzt die Variablen danach zurück, wie es `VAR=x cmd` in einer Shell tut, und übergibt Claude Code seine Befehlszeile über `Start-Process`, sodass Windows PowerShell 5.1 und PowerShell 7 sie gleich ausführen. Setzt ein Start eine Variable auf einen leeren Wert, sagt eine `#`-Zeile darunter, dass PowerShell vor 7.5 eine solche Variable stattdessen entfernt. Für ein von npm installiertes `claude.cmd` führt die Zeile den `cmd.exe`-Befehl aus, den der Start verwenden würde (siehe unten), und ein Argument, das der Start dort ablehnen würde, lässt auch den Probelauf fehlschlagen.

## So funktioniert es

Jeder Start durchläuft dieselben Schritte, von der Befehlszeile bis `exec`:

<!-- motion:launch -->

Jedes Feld leuchtet auf, während sein Schritt läuft, und die Zeilen darunter sind die Umgebung und die Argumente, mit denen Claude Code startet.

1. ccam startet in `app.Main`, das `cmd/ccam/main.go` aufruft. Wurde ccam unter einem anderen Namen gestartet, ist es ein [Kurzbefehl](/features/shortcuts) und geht direkt an den Starter.
2. ccam wertet `os.Args` aus, bevor cobra beteiligt ist (`app.parse`). Nur `ccam` allein (das Dashboard), ein ccam-Befehlsname (`config.Reserved`), `-h`/`--help` und `-v`/`--version` gehen an `cli.Execute`. Alles andere ist ein Start.
3. ccam liest `~/.config/ccam/config.toml`, die einzige Datei, die es parst (`launcher.Launch`).
4. ccam findet das Konto über Namen oder Nummer, entfernt ein führendes `--` und ermittelt die ausführbare Datei `claude` (`launcher.Prepare`, `ResolveClaude`). Ohne Namen nimmt es das Konto des Ordners und vergleicht dazu das Arbeitsverzeichnis als Zeichenkette mit den `[[folder]]`-Einträgen von `config.toml`, unter macOS und Windows ohne Beachtung der Groß- und Kleinschreibung (`HereProfile`). Die Suche im PATH wird zusammen mit der Änderungszeit der Datei und einem Hash von `PATH` in `~/.local/state/ccam/claude-path` zwischengespeichert, sodass ccam nicht bei jedem Start `PATH` durchsucht. Eine geänderte oder fehlende Datei oder ein geänderter `PATH` löst eine neue Suche aus. Ein `claude_bin`-Pfad in `config.toml` (`~/` wird erweitert) muss eine ausführbare Datei sein; andernfalls hält der Start mit `claude_bin X is not an executable file. Check claude_bin in config.toml` an, bevor irgendetwas geschrieben wird. ccam lehnt ein `claude` ab, das ccam selbst ist (Vergleich über den Inode), weil das eine Endlosschleife wäre.
5. ccam vergleicht den Kontoordner `~/.local/share/ccam/profiles/<name>/` mit seinen Freigabeeinstellungen und ruft dabei nur `lstat` und `readlink` auf, ein oder zwei Systemaufrufe pro Element (`checkLinks`, `share.Check`). Ein fehlender Link wird sofort angelegt. Alles, was Daten kosten könnte, hält den Start an.
6. ccam fügt die geteilten MCP-Server hinzu (`applyMCP`; siehe unten).
7. ccam erfasst die Sitzung (`session.Register`). Es schreibt eine nach seiner PID benannte Markierung mit der Startzeit des Prozesses nach `~/.local/state/ccam/locks/profile-<name>/<pid>`, unabhängig von `XDG_RUNTIME_DIR`, sodass jede Shell dieselben Markierungen sieht. Bei einem Konto, das den Gesprächsverlauf teilt, sucht ccam zuerst in `history-<group>-<hash>` nach einer laufenden Sitzung, wobei der Hash aus dem aktuellen Ordner mit aufgelösten symbolischen Links gebildet wird, und registriert sich auch dort (`historyCheck`).

   `exec` behält die PID, daher bezeichnet die Markierung danach den Claude Code-Prozess selbst. ccam verwendet eine Markierungsdatei statt einer gehaltenen Sperre, weil Prozesse, die Claude Code startet, eine Sperre erben und sie nach dem Ende der Sitzung behalten würden. Eine Markierung, die sich nicht schreiben lässt, hält den Start nicht an: ccam meldet `ccam: could not mark the session as running (…); starting Claude Code anyway`.

   Die anderen Prüfungen von ccam, ob ein Konto läuft (Umbenennen, Entfernen, Gruppenänderungen, `doctor --fix`, die MCP-Synchronisierung), lesen außerdem die eigenen Einträge von Claude Code im Kontoordner (`sessions/<pid>.json`, `daemon.lock`), sodass sie auch Claude Code sehen, das ohne ccam gestartet wurde: einfaches `claude`, eine IDE, Hintergrundsitzungen und ihren Daemon.
8. ccam setzt den Terminaltitel auf `claude · <name>` (wenn die Standardfehlerausgabe ein Terminal ist), schreibt die Startzeit ohne `fsync` nach `~/.local/state/ccam/last-used/<name>` und ruft `syscall.Exec` auf (`Plan.Exec`).

Claude Code läuft mit `CLAUDE_CONFIG_DIR` auf den Kontoordner und `CCAM_PROFILE` auf den Kontonamen gesetzt. ccam leert geerbte Werte dieser Variablen und ein geerbtes `CCAM_RELAY`. Für einen mit `ccam import` registrierten Ordner `~/.claude` unter macOS setzt es außerdem `CLAUDE_SECURESTORAGE_CONFIG_DIR` auf einen leeren Wert, damit das Konto und einfaches `claude` dieselbe Anmeldung im Schlüsselbund verwenden; unter Windows tut es dasselbe, wenn Claude Code die Anmeldung in der Anmeldeinformationsverwaltung speichert.

### Warum exec

ccam ersetzt seinen eigenen Prozess durch `claude` (`syscall.Exec`). Dadurch gilt:

- ccam bleibt nicht im Speicher;
- Tastendrücke, Signale und Exit-Code gehören Claude Code selbst;
- Argumente bleiben ein Array, nie ein Shell-String.

Windows hat kein `exec`. Dort startet ccam `claude.exe` als Kindprozess mit derselben Umgebung, denselben Argumenten und derselben Konsole, wartet darauf und beendet sich mit dessen Exit-Code (`Plan.Exec`). Ein von npm installiertes `claude.cmd` läuft in `cmd.exe` mit seinem vollständigen Pfad, und jedes Argument wird so in Anführungszeichen gesetzt, dass `cmd.exe` keines davon als eigene Syntax liest (`&`, `%PATH%`); ein Argument mit Zeilenumbruch lässt sich so nicht übergeben und wird abgelehnt, mit dem Rat, das native Claude Code zu installieren.

Strg+C erreicht Claude Code direkt über die Konsole, deshalb ignoriert ccam es. Die Sitzungsmarkierung trägt die PID von ccam, die genau so lange lebt wie Claude Code, und wird entfernt, wenn Claude Code sich beendet. Während es wartet, tut ccam sonst nichts.

Arbeit nach einer Sitzung, etwa das Bemerken von Servern, die mit `claude mcp add` hinzugefügt wurden, geschieht beim nächsten Start: ccam vergleicht die Änderungszeit von `.claude.json` mit `~/.local/state/ccam/mcp/<name>/seen` und parst die Datei nur, wenn sie sich geändert hat.

### Wie eine Befehlszeile weitergeleitet wird

Drei Befehlszeilen, drei Wege durch `internal/app`:

<!-- motion:launch-route -->

Ein Kurzbefehl wird über `argv[0]` erkannt, `ccam <name>` über das erste Argument, und beide enden in `syscall.Exec`. Nur ein ccam-Befehl wie `ccam usage` erreicht cobra.

### Geteilte MCP-Server

Mit der Standardstrategie `flag` übergibt ccam die MCP-Server der Gruppe als Argument an Claude Code, vor deinen eigenen:

<!-- motion:launch-mcp -->

ccam fügt das Argument nur hinzu, wenn die `mcp.json` der Gruppe existiert, und setzt es direkt hinter `claude` (`mcp.FlagArgs`, `Plan.InsertArgs`). Dein eigenes `--mcp-config` oder `--strict-mcp-config` gilt dann zusätzlich weiterhin. Es ist ein einziges Argument, `--mcp-config=<path>`, weil `--mcp-config` von Claude Code mehrere Werte nimmt: Als zwei Argumente geschrieben, würde es auch einen folgenden Prompt als Konfigurationsdatei lesen.

ccam entscheidet, wo die Option nicht hin darf (`claudeVerb`). Es liest den Unterbefehl von Claude Code so aus den Argumenten, wie Claude Code es tut, und überspringt dabei Optionen samt ihren Werten. `remote-control` oder `rc` bekommt keine Option (siehe unten), ebenso wenig ein zuerst angegebenes Verb für Hintergrundsitzungen (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`). Auch `ccam work mcp list`, `get` und `remove` bekommen keine, da Claude Code sie dort ablehnt; stattdessen meldet ccam auf der Standardfehlerausgabe, dass die Server der Gruppe auch in den Sitzungen des Kontos geladen werden:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Jeder andere Unterbefehl und jede Sitzung behält die Option.

Mit `mcp_strategy = "sync"` in `config.toml` schreibt ccam die Server stattdessen in die `.claude.json` des Kontos und lässt das aus, wenn das Konto bereits läuft. Ein Start lässt die Arbeit auch aus, wenn sich seit der letzten Prüfung weder die `mcp.json` der Gruppe noch `.claude.json` geändert hat. Hat das Konto einen eigenen Server mit demselben Namen wie ein Gruppenserver, behält die Synchronisierung den Server des Kontos und sagt das einmal. Sie legt außerdem eine Kopie zurück, die aus `.claude.json` verschwunden ist (`put back group default's MCP servers …`), und ändert oder entfernt nur die Kopien, die sie selbst geschrieben hat und die seitdem niemand geändert hat.

### Remote Control

`ccam work rc` (oder `ccam work remote-control`) startet die Remote Control von Claude Code als `work`. `claude remote-control` lehnt `--mcp-config` vor seinem Verb ab und nimmt danach keins an, und seine Sitzungen lesen MCP-Server nur aus der `.claude.json` des Kontos. Deshalb arbeitet die Flag-Strategie hier anders:

<!-- motion:launch-rc -->

1. ccam findet `remote-control` oder `rc` als Verb von Claude Code, selbst nach Optionen und ihren Werten (`ccam work --model opus rc`), und fügt keine Option hinzu (`claudeVerb`, `mcp.FlagArgs`).
2. ccam vergleicht die Änderungszeit der `mcp.json` der Gruppe mit dem `sync-stamp` des Kontos und prüft, dass jeder Gruppenserver noch in `.claude.json` steht (`mcp.Synced`). Enthält `.claude.json` bereits die aktuellen Server der Gruppe oder hat die Gruppe keine, startet Remote Control ohne Rückfrage.
3. Andernfalls nennt ccam die Server, die den Sitzungen fehlen würden, und fragt:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. Bei `y` kopiert ccam sie wie die Strategie `sync` (`mcp.Sync`): unter der Sperre von Claude Code, nach einer Sicherung und ohne die eigenen Server des Kontos anzutasten. Bei `N` oder ohne Terminal startet Remote Control ohne sie, und ccam verweist auf `mcp_strategy = "sync"`, das sie vor jedem Start kopiert.

ccam fragt bei jedem Start von Remote Control, bis die Server kopiert sind. Nach einer Kopie fragt es erst wieder, wenn sich die `mcp.json` der Gruppe ändert oder ein kopierter Server aus `.claude.json` verschwindet. Eine Kopie in `.claude.json` gewinnt bei späteren normalen Starts gegen `--mcp-config`, deshalb aktualisieren `ccam mcp add`, `remove` und `import` die Kopien, die andere Konten haben ([Geteilte MCP-Server](/features/mcp)).

### Das Budget von 10 ms

Das Budget für den Startpfad liegt bei 10 ms im 95. Perzentil; in der Praxis fügt ccam etwa 4 ms hinzu, bevor Claude Code startet. Der Startpfad baut nie den cobra-Befehlsbaum auf, startet nie Bubble Tea und öffnet keine Netzwerkverbindung, und `TestPackageInitIsCheap` schlägt fehl, wenn die `init()`-Funktionen der Pakete zusammen mehr als 5 ms brauchen. `scripts/perf-check.sh` misst `ccam --dry-run work`, `ccam work` und einen Kurzbefehl mit `hyperfine` und schlägt oberhalb von 10 ms fehl, und `BenchmarkLaunchPath` in `internal/launcher` führt die Schritte in einer Schleife aus. Um die Zeit pro Schritt zu sehen, setze `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Sicherheit und Sonderfälle

- Ein unbekannter Name oder eine unbekannte Nummer hält mit `no account "<name>"` an und verweist auf `ccam list`. Gibt es gar keine Konten, bittet ccam dich, mit `ccam add <name>` eines hinzuzufügen.
- Ist ein geteilter Pfad zu einer gewöhnlichen Datei geworden oder zeigt er an die falsche Stelle, hält ccam mit `stopped because of share link problems:` an, listet die Probleme auf und bittet dich, `ccam doctor --fix` auszuführen. Auf dem Startpfad verschiebt es keine Daten.
- Läuft im selben Ordner eine andere Claude Code-Sitzung, die den Gesprächsverlauf mit diesem Konto teilt (dieselbe Gruppe), warnt ccam, dass sich der Verlauf vermischen kann, und fragt `Run anyway? [y/N]`. Der Ordner wird mit aufgelösten symbolischen Links verglichen. Ohne Terminal (Standardeingabe aus `/dev/null` zählt als keines) oder aus einem Hintergrundjob zeigt es die Warnung und macht weiter. Ein neuer `-p`-Lauf beginnt ein eigenes Gespräch und bekommt keine Warnung; `--version`, `--help` und Unterbefehle, die kein Gespräch öffnen, etwa `mcp`, `auth` oder `plugin`, werden gar nicht geprüft.
- `--dry-run` überspringt die MCP-Synchronisierung, den Hinweis auf neue Server, den Hinweis zu `mcp list`, die Verlaufsprüfung, die Sitzungsmarkierung, den Titel und die Zeit der letzten Verwendung und führt kein exec aus. Die Freigabelinks prüft es trotzdem und legt einen fehlenden Link an wie ein normaler Start.
- `ccam work rc` bei bereits laufendem `work` lässt `.claude.json` unverändert, weist darauf hin und startet Remote Control ohne die Server der Gruppe; führe es nach dem Ende der Sitzung erneut aus, um sie zu kopieren. Mit `--dry-run` gibt es die Meldung aus, fragt aber nicht und kopiert nichts.
- Mit `mcp_strategy = "sync"` gibt eine fehlgeschlagene Synchronisierung `ccam: MCP sync failed (launching anyway): …` aus, und Claude Code startet mit den Servern, die es bereits hatte.

## Verwandte Themen

- [Kurzbefehle](/features/shortcuts)
- [Remote Control mit ccam nutzen](/guides/remote-control)
- [Konten verwalten](/features/accounts)
- [Befehlsreferenz](/docs/commands)
- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [CLAUDE_CONFIG_DIR erklärt](/guides/claude-config-dir)
