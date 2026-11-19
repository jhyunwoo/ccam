---
title: Konten verwalten
description: "Claude Code-Konten anlegen, anmelden, auflisten, umbenennen, entfernen und importieren: mit ccam add, login, list, rename, remove und import."
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

Jedes ccam-Konto ist ein eigener Claude Code-Konfigurationsordner mit eigener Anmeldung. Mit diesen Befehlen legst du Konten an, meldest sie an, listest sie auf, benennst sie um und entfernst sie. Sie registrieren auch Ordner, die du bereits hast.

## Verwendung

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account that was never logged in or whose login expired
ccam logout side                           # log out: the login and its MCP sign-ins
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, logs out, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Wo ein Befehl `<name|number>` erwartet, kannst du den Kontonamen oder seine Slot-Nummer angeben (die Spalte `#` in `ccam list`).

`ccam use` ändert nur, welches Konto das Standardkonto ist (`default` in `config.toml`).

`ccam add` ohne Optionen öffnet in einem Terminal einen Assistenten. Er geht diese Schritte durch:

1. den Kontonamen, sofern du keinen angegeben hast;
2. wie sich das Konto anmeldet: jetzt mit einer Claude-Anmeldung, mit einem API-Schlüssel oder später mit einer Claude-Anmeldung;
3. welchen Gruppen es beitritt, mit angehaktem `default` (gibt es noch keine Gruppen, zeigt er, was die Gruppe `default` teilen wird);
4. `Create a shortcut command for work? [y/N]`;
5. eine Zusammenfassung zum Bestätigen.

`Esc` geht einen Schritt zurück. Jede der Optionen `--group`, `--no-group`, `--shortcut` oder `--no-login` oder ein Aufruf ohne Terminal (zum Beispiel aus einer Pipe) überspringt den Assistenten. Dann ist ein Name Pflicht, und ein Kurzbefehl wird nur angelegt, wenn `--shortcut` einen nennt.

Ohne `--group` tritt ein neues Konto der Gruppe `default` bei. Das erste Konto legt sie an und teilt darin MCP-Server, Einstellungen, `CLAUDE.md` und Erweiterungen; siehe [Teilen mit Gruppen](/features/share). Unter Windows lassen sich Dateien nur im Entwicklermodus verlinken. Ohne ihn legt das erste Konto `default` nur mit `mcp` und `extensions` an und sagt das (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), und ccam lehnt `--group` mit einer Gruppe ab, die Einstellungen, Memory oder Verlauf teilt.

`ccam list` zeigt Slot, Name, Kurzbefehl, E-Mail-Adresse, Gruppen und letzte Verwendung, mit `*` vor dem Standardkonto. Die E-Mail-Adresse ist die, mit der der Ordner gerade angemeldet ist (aus seiner `.claude.json`); weicht sie von der ab, die ccam festgehalten hat, folgt die festgehaltene als `(was …)`:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

```console
$ ccam list
   #  Name      Shortcut  Email                                     Groups   Last used
*  1  work      ccwork    other@company.com (was work@example.com)  default  2m ago
```

Mit `--json` ist jedes Konto ein Objekt mit `slot`, `shortcut`, `name`, `email` (die der aktuellen Anmeldung, wie in der Tabelle), `default`, `dir`, `last_used` (fehlt, wenn das Konto nie gestartet wurde), `groups` und `share`. `share` ordnet jedem Element (`mcp`, `settings`, `memory`, `extensions`, `history`) die Gruppe zu, aus der es kommt, oder `"local"`.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `ccam add --no-login` | Das Konto ohne Anmeldung anlegen; später mit `ccam login <name>` anmelden | aus |
| `ccam add --group <group>` | Dieser Gruppe statt `default` beitreten. Mehrfach verwendbar; `default` wird angelegt, wenn es sie nicht gibt | `default` |
| `ccam add --no-group` | Keiner Gruppe beitreten: Das Konto teilt nichts | aus |
| `ccam add --shortcut <command>` | Zusätzlich diesen Kurzbefehl in `~/.local/bin` anlegen | keiner |
| `ccam add --api --key-from <source>` | Ein Konto, das einen API-Schlüssel statt einer Anmeldung verwendet; siehe [API-Schlüssel-Konten](/features/api) | aus |
| `ccam list --json` | Die Konten als JSON ausgeben | aus |
| `ccam remove -y`, `--yes` | Ohne Nachfrage entfernen | aus |
| `ccam logout -y`, `--yes` | Ohne Nachfrage abmelden | aus |
| `ccam import --name <name>` | Kontoname für den angegebenen Ordner; ohne Ordner abgelehnt (`--name needs a folder: ccam import <folder> --name ops`) | aus dem Ordnernamen abgeleitet |
| `ccam import -y`, `--yes` | Jeden gefundenen Ordner ohne Nachfrage importieren | aus |
| `--dry-run` | Wird von `ccam list` und `ccam completion` unterstützt, die nichts ändern. Abgelehnt von `add`, `login`, `logout`, `use`, `rename`, `remove` und `import`, mit `` --dry-run is not supported by `ccam add`; nothing was changed `` (mit dem jeweiligen Befehl) | aus |

`ccam login` (abgesehen von `--expired`), `ccam use` und `ccam rename` haben keine eigenen Optionen. `ccam completion` nimmt `bash`, `zsh` oder `powershell`; alles andere schlägt mit `only bash, zsh and powershell are supported` fehl.

Kontonamen bestehen aus 1 bis 32 Zeichen aus `a-z`, `0-9`, `-` und `_` und beginnen mit einem Buchstaben oder einer Ziffer. Ein Name darf nicht nur aus Ziffern bestehen, weil `ccam 2` Slot 2 bedeutet. Ein Name darf auch kein ccam-Befehl sein (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), damit `ccam <name>` nie mehrdeutig ist. Die Liste ist `config.Reserved` in `internal/config`.

## So funktioniert es

Das Hinzufügen eines Kontos legt seinen Ordner an, tritt seinen Gruppen bei, speichert die Konfiguration, meldet an und legt den Kurzbefehl an, in dieser Reihenfolge. Wenn du `ccam add work --shortcut ccwork` ausführst, geht ccam diese Schritte durch (`addProfile` in `internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. ccam prüft den Namen (`config.ValidateName`) und den Kurzbefehl (`shim.Check`), bevor es irgendetwas anlegt.
2. ccam legt `~/.local/share/ccam/profiles/work/` mit Modus `0700` an und gibt dem Konto die niedrigste freie Slot-Nummer, sodass die Nummer eines entfernten Kontos wiederverwendet wird (`profile.Add`, `Config.NextSlot`). Das erste Konto wird zum Standardkonto.
3. Das Konto tritt seinen Gruppen bei (`default`, sofern `--group` oder `--no-group` nichts anderes sagt; das erste Konto legt `default` an), und ccam verlinkt den neuen, leeren Ordner in sie, zum Beispiel `settings.json` auf `groups/default/settings.json` (`applyShare`, mit `share.Build` und anschließend `Apply`). Schlägt das fehl, wird der neue Ordner entfernt, weil ein übrig gebliebener Ordner jeden weiteren Versuch mit `folder already exists` blockieren würde.
4. Das Konto wird über `fsutil.AtomicWrite` in `~/.config/ccam/config.toml` gespeichert: eine temporäre Datei, `fsync`, dann `rename`, sodass ein Absturz nie eine halb geschriebene Konfiguration hinterlässt.
5. ccam startet `claude auth login` mit `CLAUDE_CONFIG_DIR` auf den Kontoordner gesetzt und mit `--email` und der letzten E-Mail-Adresse des Kontos, damit die Anmeldeseite ausgefüllt ist (`runLogin`). Schlägt die Anmeldung fehl, bleibt das Konto bestehen, und ccam gibt den Fehler aus und rät dir, später `ccam login work` auszuführen.
6. ccam liest `oauthAccount` aus `profiles/work/.claude.json` (`profile.ReadIdentity`): die E-Mail-Adresse, die Konto- und Organisations-IDs, den Namen der Organisation und den Tarif (`claude_max`, `claude_team` …). Diese Datei kann mehrere Megabyte groß sein, daher sucht ccam nur nach diesem einen Schlüssel, statt die ganze Datei zu parsen (`TopLevelValue`). ccam speichert diese Angaben in `[profile.login]` in `config.toml`.
7. ccam verlinkt `~/.local/bin/ccwork` mit der ccam-Binärdatei, die beim Start `argv[0]` in `config.toml` nachschlägt, und hält den Link in `~/.local/state/ccam/shortcuts.json` fest (`shim.Install`). Lässt sich der Link nicht anlegen, wird das Konto trotzdem hinzugefügt, ccam sagt, warum (`Could not create shortcuts: …`), und die Zusammenfassung zeigt `shortcut: ccwork (not created; see above)`.

Eine Anmeldung kann Minuten dauern, und ein anderer ccam-Prozess kann die Konfiguration in der Zwischenzeit speichern. Deshalb laden Befehle, die warten, `config.toml` neu, wenden die eine Änderung an (hier die E-Mail-Adresse) und speichern (`env.update`). `ccam remove` lädt aus demselben Grund nach seiner Nachfrage neu.

`ccam login` führt die Schritte 5 und 6 für ein vorhandenes Konto aus. Claude Code nimmt die E-Mail-Adresse nur als Hinweis für die Anmeldeseite, daher kann es sein, dass du bei einem anderen Claude-Konto angemeldet bist als zuvor; ccam vergleicht die neue Anmeldung mit `[profile.login]` und sagt es (`ccam: this login is other@company.com, but work was dev@company.com before.`). Es behält die neue Anmeldung und sagt dir, wie du sie rückgängig machst: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

Eine Anmeldung in einem anderen Claude-Konto verwirft außerdem den [Limit-Eintrag](/features/usage#the-limit-record) des Kontos, da diese Limits die des anderen Kontos waren. Ist ein anderes ccam-Konto beim selben Claude-Konto angemeldet, warnt ccam, dass sich die beiden ihre Nutzungslimits teilen.

`ccam login --expired` meldet nacheinander jedes Konto an, das nie angemeldet war oder dessen Anmeldung abgelaufen ist: Sein Refresh-Token ist abgelaufen, oder Claude Code hat die Anmeldung gelöscht, nachdem der Server sie abgelehnt hat. Bevor sich der erste Browser öffnet, prüft es, ob sich `config.toml` speichern lässt.

`ccam logout` führt `claude auth logout` im Konto aus und löscht die gespeicherte E-Mail-Adresse; der Anmeldeeintrag bleibt, sodass die nächste Anmeldung mit ihm verglichen wird. Claude Code widerruft die Anmeldung und löscht, was es für das Konto speichert, einschließlich der Anmeldungen bei MCP-Servern. Deshalb warnt ccam dich und fragt zuerst nach (`--yes` überspringt die Frage), wenn:

- die Anmeldung MCP-Anmeldungen enthält;
- das Konto läuft (seine Sitzungen werden ebenfalls abgemeldet);
- der Ordner eines anderen Kontos dieselbe Anmeldung enthält (ein kopierter Ordner), die die Abmeldung ebenfalls widerruft.

Zum Beispiel:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

Um dich als ein anderes Claude-Konto anzumelden, genügt `ccam login <name>` allein, und die MCP-Anmeldungen bleiben erhalten.

`ccam list` fügt eine Spalte `Organization` hinzu, wenn die Konten zu mehr als einer Organisation gehören (private Tarife zeigen `personal`), und `ccam list --json` hat `plan` und `org` für jedes Konto.

### Anmeldung auf einem entfernten Rechner

Über SSH kann der Browserteil von `claude auth login` scheitern. Claude Code wartet auf die Anmeldung über einen Callback-Server auf dem entfernten Rechner (`http://localhost:<port>/callback`), aber der Browser, der die Anmeldung abschließt, läuft auf deinem Computer. Dort führt diese Adresse ins Leere, und der Browser zeigt einen Fehler wie `ERR_CONNECTION_REFUSED`. Die Adresse enthält trotzdem den Einmalcode, also kann ccam ihn zustellen:

<!-- motion:accounts-relay -->

ccam übergibt dein Terminal nicht an Claude Code (`runLogin`). Es verbindet die Standardeingabe von Claude Code mit einer Pipe und kopiert jede Zeile, die du tippst, hinein (`relayLogin`), sodass ein Code, den du bei der Eingabeaufforderung `Paste code here` von Claude Code einfügst, trotzdem ankommt. Ist eine Zeile eine Callback-Adresse, akzeptiert ccam sie nur, wenn der Host `localhost`, `127.0.0.1` oder `[::1]` ist, ein Port angegeben ist, der Pfad `/callback` lautet und sowohl `code` als auch `state` vorhanden sind; `http://` darf fehlen (`callbackURL`). ccam ruft diese Adresse dann auf dem Server mit einem Timeout von 10 Sekunden und ohne Weiterleitungen zu folgen ab, so wie es der Browser getan hätte (`deliverCallback`). Die Anfrage erreicht den Callback-Server des wartenden `claude auth login` (die gestrichelte Linie im Diagramm). Claude Code erhält den Code und speichert die Anmeldung, unter Linux in `.credentials.json` im Kontoordner. Eine eingefügte Zeile kann ccam nie dazu bringen, einen anderen Host zu kontaktieren.

ccam gibt dies aus, bevor die Anmeldung beginnt:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

Bei Erfolg gibt es `→ Passed the sign-in to Claude Code on this machine.` aus.

Während die Anmeldung in einem Terminal auf eine Eingabe wartet, schaltet ccam das Echo und den Zeilenmodus des Terminals aus und zeigt für jedes getippte Zeichen ein `*`, sodass ein eingefügter Code oder eine Callback-Adresse nicht auf dem Bildschirm erscheint. ccam übernimmt die Zeilenbearbeitung selbst (`readLine`): `Enter` sendet die Zeile, `Backspace` löscht ein Zeichen, und `Ctrl+D` in einer leeren Zeile beendet die Eingabe. Das Terminal wird wiederhergestellt, wenn die Anmeldung endet, und auch bei `Ctrl+C`, `SIGTERM` oder `SIGHUP`, die ccam weiterhin beenden. Ohne Terminal an Standardeingabe und Standardfehlerausgabe werden Zeilen so kopiert, wie sie ankommen.

### Umbenennen, entfernen und importieren

`ccam rename` und `ccam remove` wirken auf den Ordner des Kontos und auf das, was ccam an anderer Stelle unter dem Namen des Kontos aufbewahrt. ccam listet diese Pfade in `nameState` auf:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- die Nutzungscaches `~/.local/state/ccam/usage/<name>.json` und `<name>.statusline.json`;
- den Limit-Eintrag `<name>.limit.json`;
- die Sitzungsnotizen in `~/.local/state/ccam/handoff/<name>/`.

`ccam remove work` behandelt sie in dieser Reihenfolge:

<!-- motion:accounts-remove -->

1. `ccam remove` lehnt zuerst ab, wenn das Konto läuft, fragt dann nach und prüft nach deiner Antwort erneut (so wird eine Sitzung erkannt, die gestartet wurde, während die Frage offen war).
2. Bevor es irgendetwas ändert, prüft es, ob sich `config.toml` speichern lässt.
3. Teilt das Konto Erweiterungen, schreibt ccam die Dateien `installed_plugins.json` und `known_marketplaces.json` der Gruppe neu (`RewritePluginPaths`), weil Claude Code absolute Plugin-Pfade speichert, die durch den Ordner dieses Kontos führen können.
4. Ein verwaltetes Konto (eines, dessen Ordner ccam angelegt hat, im Gegensatz zu einem mit `ccam import` registrierten) wird als Nächstes abgemeldet: `claude auth logout` läuft in seinem Ordner, sodass keine funktionierende Anmeldung im Papierkorb wartet (oder, unter macOS, in einem nach dem Ordner benannten Schlüsselbundeintrag, den ein späteres Konto mit diesem Namen erben würde). Das entfällt bei einem API-Schlüssel-Konto und wenn der Ordner eines anderen Kontos dieselbe Anmeldung enthält, die die Abmeldung ebenfalls widerrufen würde. Eine fehlgeschlagene Abmeldung wird gemeldet, und das Entfernen geht weiter.
5. ccam verschiebt den Ordner nach `~/.local/share/ccam/trash/<name>-<time>/` (`profile.Trash`), das Konto wird aus `config.toml` entfernt, die `nameState`-Pfade werden gelöscht, damit ein späteres Konto mit demselben Namen sie nicht erbt, und der Link des Kurzbefehls wird entfernt: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

Lässt sich `config.toml` am Ende doch nicht speichern, kommt der Ordner aus dem Papierkorb zurück. Um ein entferntes Konto wiederherzustellen, verschiebst du den Ordner aus dem Papierkorb zurück und meldest dich erneut an.

`ccam rename` benennt den Ordner um und schreibt die Plugin-Pfade in seinem `plugins/` vom alten auf den neuen Ordner um. Es verschiebt die `nameState`-Pfade auf den neuen Namen und behält das Konto als Standardkonto bei, wenn es eines war. Außerdem passt es die Sitzungsnotizen an, sodass `ccam next --continue` sie weiterhin findet. Der Kurzbefehl behält seinen Namen. Schlägt das Umschreiben der Plugins fehl oder lässt sich `config.toml` nicht speichern, wird der Ordner zurückbenannt.

Unter macOS, und unter Windows, wenn Claude Code Anmeldungen in der Anmeldeinformationsverwaltung speichert, behält der Ordner selbst seinen alten Namen, weil die Anmeldung und die MCP-Anmeldungen unter dem Pfad des Ordners gespeichert sind und ein verschobener Ordner abgemeldet wäre. Die Konfiguration hält den Ordner dann als `folder` neben dem neuen Namen fest, und ccam sagt `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

`ccam import` registriert einen Ordner dort, wo er liegt, als `dir` in `config.toml`. Es verschiebt oder kopiert ihn nie: Ein kopiertes Refresh-Token würde dazu führen, dass sich die Kopien gegenseitig abmelden.

Ohne Ordnerargument sucht es nach nicht registrierten Ordnern `~/.claude` und `~/.claude-*` und fragt bei jedem nach. Die Namen ergeben sich aus dem Ordner (`~/.claude` wird `main`, `~/.claude-work` wird `work`, `~/.claude-2` wird `acct2`). Kann ccam keinen Namen ableiten oder ist der Name vergeben, listet es den Ordner mit dem `--name`-Befehl auf, den du ausführen musst. Ein gefundener Ordner, der fehlschlägt, wird als `Skipped <path>: <reason>` gemeldet, und die anderen werden importiert.

ccam stellt zuerst alle Fragen, lädt dann `config.toml` neu und registriert die Ordner, sodass Änderungen erhalten bleiben, die andere ccam-Befehle in der Zwischenzeit gespeichert haben. Die `✓`-Zeilen kommen nach dem Speichern.

Ohne `CLAUDE_CONFIG_DIR` hält Claude Code `.claude.json` in deinem Home-Ordner neben `~/.claude`. Wenn du also `~/.claude` importierst und es keine eigene `.claude.json` hat, legt ccam den Link `~/.claude/.claude.json` auf `~/.claude.json` an. Einfaches `claude` und `ccam main` sehen dann dasselbe Konto.

Hat `~/.claude` neben `~/.claude.json` eine eigene `.claude.json`, lässt ccam beide stehen und sagt das. `ccam main` würde die innere verwenden und ein einfaches `claude` die äußere, womit ihre Anmeldedaten und MCP-Server getrennt blieben. Um nur eine zu behalten, ersetze die innere Datei durch einen Link ([doctor](/features/doctor) meldet das ebenfalls).

Enthält der Ordner dieselbe Anmeldung wie ein registriertes Konto (ein kopierter Ordner), sagt ccam das: Wer sie zuerst erneuert, meldet den anderen ab, also führe `ccam login <name>` aus, bevor du ihn verwendest.

Importierte Konten sind in keiner Gruppe, bis du sie einer hinzufügst (`ccam group join <group> <name>`), und es wird kein Kurzbefehl angelegt.

`ccam completion bash|zsh` gibt ein Vervollständigungsskript aus, das auch Kontonamen vervollständigt.

## Sicherheit und Sonderfälle

- `ccam rename` und `ccam remove` lehnen ab, solange das Konto läuft (`refuseRunning`), mit ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``, weil die Sitzung weiter in den Ordner schreibt. Die Prüfung erfasst Claude Code, egal wie es gestartet wurde: über ccam, als einfaches `claude`, aus einer IDE, Hintergrundsitzungen und ihren Daemon.
- `ccam rename`, `ccam remove` und `ccam logout` lehnen ab, solange ein anderer ccam-Prozess diese Anmeldung erneuert oder ändert: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- Während in einem anderen ccam-Prozess eine Gruppenänderung läuft, lehnen `ccam rename` und `ccam remove` mit `another ccam process is changing the groups right now; try again when it is done` ab; nachdem eine mittendrin abgebrochen ist, mit ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.
- `ccam remove` behält bei einem importierten Konto den Ordner und seine Anmeldung für ein einfaches `claude`. Ist das Konto in Gruppen, verlässt es diese zuerst so wie bei `ccam group leave`: Jedes geteilte Element wird zu einer Kopie dessen der Gruppe im Ordner, und die MCP-Server der Gruppe kommen in seine `.claude.json`. Dann wird die Registrierung aufgehoben: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. Die Sicherungen seiner `.claude.json`, die ccam angelegt hat, wandern nach `<folder>/.ccam-backup/mcp-<time>`; lassen sie sich nicht verschieben, bleiben sie in `~/.local/state/ccam/mcp/<name>/backups`, und die Meldung sagt das. Dein `~/.claude` wird nicht in den Papierkorb verschoben.
- Zwei Konten können nicht denselben Ordner haben: Eine `config.toml`, die einen Ordner zweimal nennt, lässt sich nicht laden.
- `ccam add` lehnt einen Namen ab, dessen Ordner in `profiles/` bereits existiert, und verweist auf `ccam import`. Unter macOS hält es außerdem an, wenn der Schlüsselbund noch eine Anmeldung oder MCP-Anmeldungen für den Pfad des neuen Ordners enthält, die eine ältere Version von ccam hinterlassen hat, und nennt den Befehl `security delete-generic-password`, der sie entfernt. Lässt sich `config.toml` nicht speichern, wird der neue Ordner wieder entfernt.
- `ccam import` lehnt einen bereits registrierten Ordner, einen Pfad, der kein Ordner ist, und einen vergebenen Namen ab. ccam warnt außerdem, wenn der Ordner eine Kopie von `~/.claude` ist (dieselbe Anmeldung), auch wenn `~/.claude` selbst nicht registriert ist.
- `ccam logout` fragt zuerst nach, wenn das Konto läuft; `ccam login` prüft das nicht.
- Wenn die Weiterleitung Claude Code nicht erreicht, zum Beispiel nachdem die Anmeldung beendet wurde, gibt sie `ccam: could not pass the sign-in to Claude Code:` mit dem Grund aus. Ein HTTP-Fehler bedeutet meist, dass der Code abgelaufen ist; starte die Anmeldung neu.

## Verwandte Themen

- [Claude Code starten](/features/launch)
- [Kurzbefehle](/features/shortcuts)
- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Mehrere Claude Code-Konten auf einem Rechner nutzen](/guides/multiple-claude-code-accounts)
