---
title: Gestire gli account
description: Aggiungi, collega, elenca, rinomina, rimuovi e importa account Claude Code con ccam add, ccam login, ccam list, ccam rename, ccam remove e ccam import.
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

Ogni account di ccam è una cartella di configurazione di Claude Code separata, con il proprio login. Usa questi comandi per aggiungere gli account, effettuarne il login, elencarli, rinominarli e rimuoverli. Servono anche a registrare le cartelle che hai già.

## Utilizzo

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

Dove un comando accetta `<name|number>`, puoi indicare il nome dell'account oppure il suo numero di slot (la colonna `#` in `ccam list`).

`ccam use` cambia solo quale account è quello predefinito (`default` in `config.toml`).

`ccam add` in un terminale senza opzioni apre una procedura guidata, che segue questi passi:

1. il nome dell'account, se non l'hai indicato;
2. come accede l'account: un login di Claude subito, una chiave API o un login di Claude più tardi;
3. a quali gruppi unirsi, con `default` selezionato (se non esiste ancora nessun gruppo, mostra cosa condividerà il gruppo `default`);
4. `Create a shortcut command for work? [y/N]`;
5. un riepilogo da confermare.

`Esc` torna indietro di un passo. Una qualsiasi tra `--group`, `--no-group`, `--shortcut` o `--no-login`, oppure l'esecuzione senza un terminale (per esempio da una pipe), salta la procedura guidata. In quel caso il nome è obbligatorio, e non viene creata nessuna scorciatoia a meno che `--shortcut` non ne indichi una.

Senza `--group`, un nuovo account si unisce al gruppo `default`. Il primo account lo crea, condividendo server MCP, impostazioni, `CLAUDE.md` ed estensioni; vedi [condivisione con i gruppi](/features/share). Su Windows i file si possono collegare solo in Modalità sviluppatore. Senza di essa, il primo account crea `default` con i soli `mcp` ed `extensions` e lo segnala (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), e ccam rifiuta `--group` con un gruppo che condivide impostazioni, memoria o cronologia.

`ccam list` mostra lo slot, il nome, la scorciatoia, l'email, i gruppi e l'ultimo uso, con `*` davanti all'account predefinito. L'email è quella con cui la cartella ha effettuato il login adesso (dal suo `.claude.json`); quando è diversa da quella registrata da ccam, quella registrata segue come `(was …)`:

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

Con `--json`, ogni account è un oggetto con `slot`, `shortcut`, `name`, `email` (quella del login attuale, come nella tabella), `default`, `dir`, `last_used` (omesso se l'account non è mai stato avviato), `groups` e `share`. `share` associa ogni elemento (`mcp`, `settings`, `memory`, `extensions`, `history`) al gruppo da cui proviene, oppure a `"local"`.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `ccam add --no-login` | Crea l'account senza effettuare il login; accedi in seguito con `ccam login <name>` | disattivata |
| `ccam add --group <group>` | Si unisce a questo gruppo invece che a `default`. Ripetibile; `default` viene creato se non esiste | `default` |
| `ccam add --no-group` | Non si unisce a nessun gruppo: l'account non condivide nulla | disattivata |
| `ccam add --shortcut <command>` | Crea anche questo comando scorciatoia in `~/.local/bin` | nessuna |
| `ccam add --api --key-from <source>` | Un account che usa una chiave API invece di un login; vedi [Account con chiave API](/features/api) | disattivata |
| `ccam list --json` | Stampa gli account in JSON | disattivata |
| `ccam remove -y`, `--yes` | Rimuove senza chiedere | disattivata |
| `ccam logout -y`, `--yes` | Effettua il logout senza chiedere | disattivata |
| `ccam import --name <name>` | Nome dell'account per la cartella indicata; rifiutata senza una cartella (`--name needs a folder: ccam import <folder> --name ops`) | ricavato dal nome della cartella |
| `ccam import -y`, `--yes` | Importa ogni cartella trovata senza chiedere | disattivata |
| `--dry-run` | Supportata da `ccam list` e `ccam completion`, che non modificano nulla. Rifiutata da `add`, `login`, `logout`, `use`, `rename`, `remove` e `import`, con `` --dry-run is not supported by `ccam add`; nothing was changed `` (con il nome del comando) | disattivata |

`ccam login` (a parte `--expired`), `ccam use` e `ccam rename` non hanno opzioni proprie. `ccam completion` accetta `bash`, `zsh` o `powershell`; qualsiasi altro valore fallisce con `only bash, zsh and powershell are supported`.

I nomi degli account sono lunghi da 1 a 32 caratteri tra `a-z`, `0-9`, `-` e `_`, e iniziano con una lettera o una cifra. Un nome non può essere composto solo da cifre, perché `ccam 2` indica lo slot 2. Un nome non può nemmeno essere un comando di ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), così `ccam <name>` non è mai ambiguo. L'elenco è `config.Reserved` in `internal/config`.

## Come funziona

Aggiungere un account crea la sua cartella, lo unisce ai suoi gruppi, salva la configurazione, effettua il login e crea la scorciatoia, in quest'ordine. Quando esegui `ccam add work --shortcut ccwork`, ccam segue questi passaggi (`addProfile` in `internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. ccam controlla il nome (`config.ValidateName`) e la scorciatoia (`shim.Check`) prima di creare qualsiasi cosa.
2. ccam crea `~/.local/share/ccam/profiles/work/` con modalità `0700` e assegna all'account il numero di slot libero più basso, così il numero di un account rimosso viene riutilizzato (`profile.Add`, `Config.NextSlot`). Il primo account diventa quello predefinito.
3. L'account si unisce ai suoi gruppi (`default`, a meno che `--group` o `--no-group` non dicano altro; il primo account crea `default`), e ccam vi collega la nuova cartella vuota, per esempio `settings.json` a `groups/default/settings.json` (`applyShare`, con `share.Build` e poi `Apply`). Se questo passaggio fallisce, la nuova cartella viene rimossa, perché una cartella rimasta bloccherebbe ogni nuovo tentativo con `folder already exists`.
4. L'account viene salvato in `~/.config/ccam/config.toml` tramite `fsutil.AtomicWrite`: un file temporaneo, `fsync`, poi `rename`, così un crash non lascia mai una configurazione scritta a metà.
5. ccam avvia `claude auth login` con `CLAUDE_CONFIG_DIR` impostata sulla cartella dell'account, e `--email` con l'ultima email dell'account, così la pagina di login è già compilata (`runLogin`). Se il login fallisce, l'account resta e ccam stampa l'errore e ti dice di eseguire `ccam login work` più tardi.
6. ccam legge `oauthAccount` da `profiles/work/.claude.json` (`profile.ReadIdentity`): l'email, gli ID dell'account e dell'organizzazione, il nome dell'organizzazione e il piano (`claude_max`, `claude_team` …). Quel file può pesare diversi megabyte, quindi ccam cerca quell'unica chiave invece di analizzare l'intero file (`TopLevelValue`). ccam conserva questi dati in `[profile.login]` in `config.toml`.
7. ccam collega `~/.local/bin/ccwork` al binario di ccam, che all'avvio cerca `argv[0]` in `config.toml`, e registra il link in `~/.local/state/ccam/shortcuts.json` (`shim.Install`). Se il link non si può creare, l'account viene comunque aggiunto, ccam spiega perché (`Could not create shortcuts: …`) e il riepilogo mostra `shortcut: ccwork (not created; see above)`.

Un login può richiedere minuti, e nel frattempo un altro processo ccam potrebbe salvare la configurazione. Per questo i comandi che restano in attesa ricaricano `config.toml`, applicano l'unica modifica (qui, l'email) e salvano (`env.update`). Per lo stesso motivo `ccam remove` ricarica la configurazione dopo la sua domanda.

`ccam login` esegue i passaggi 5 e 6 per un account esistente. Claude Code usa l'email solo come suggerimento per la pagina di login, quindi puoi ritrovarti connesso a un account Claude diverso da prima; ccam confronta il nuovo login con `[profile.login]` e lo segnala (`ccam: this login is other@company.com, but work was dev@company.com before.`). Tiene il nuovo login e ti dice come annullarlo: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

Un login in un altro account Claude elimina anche la [registrazione del limite](/features/usage#the-limit-record) dell'account, perché quei limiti erano dell'altro account. Quando un altro account ccam ha effettuato il login nello stesso account Claude, ccam avverte che i due condividono i limiti di utilizzo.

`ccam login --expired` effettua il login, uno dopo l'altro, in ogni account che non ha mai effettuato il login o il cui login è scaduto: il suo refresh token è scaduto, oppure Claude Code ha cancellato il login dopo che il server lo ha rifiutato. Prima che si apra il primo browser, controlla che `config.toml` si possa salvare.

`ccam logout` esegue `claude auth logout` nell'account e cancella l'email salvata; il record del login resta, così il login successivo viene confrontato con esso. Claude Code revoca il login ed elimina ciò che conserva per l'account, compresi gli accessi ai server MCP. Per questo ccam ti avvisa e chiede prima conferma (`--yes` salta la domanda) quando:

- il login contiene accessi MCP;
- l'account è in esecuzione (anche le sue sessioni vengono disconnesse);
- la cartella di un altro account contiene lo stesso login (una cartella copiata), che il logout revoca a sua volta.

Per esempio:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

Per accedere con un altro account Claude basta `ccam login <name>`, che mantiene gli accessi MCP.

`ccam list` aggiunge una colonna `Organization` quando gli account appartengono a più di un'organizzazione (i piani personali mostrano `personal`), e `ccam list --json` ha `plan` e `org` per ogni account.

### Accedere su una macchina remota

Tramite SSH, la parte di `claude auth login` che passa dal browser può non funzionare. Claude Code attende il login su un server di callback sulla macchina remota (`http://localhost:<port>/callback`), ma il browser che completa il login è in esecuzione sul tuo computer, dove quell'indirizzo non porta da nessuna parte e il browser mostra un errore come `ERR_CONNECTION_REFUSED`. L'indirizzo contiene comunque il codice monouso, quindi ccam può consegnarlo:

<!-- motion:accounts-relay -->

ccam non cede il tuo terminale a Claude Code (`runLogin`). Collega lo standard input di Claude Code a una pipe e vi copia ogni riga che digiti (`relayLogin`), così un codice incollato al prompt `Paste code here` di Claude Code arriva comunque. Quando una riga è un indirizzo di callback, ccam lo accetta solo se l'host è `localhost`, `127.0.0.1` o `[::1]`, è indicata una porta, il percorso è `/callback` e sono presenti sia `code` sia `state`; `http://` si può omettere (`callbackURL`). ccam richiede poi quell'indirizzo sul server con un timeout di 10 secondi e senza seguire i reindirizzamenti, come avrebbe fatto il browser (`deliverCallback`). La richiesta raggiunge il server di callback del `claude auth login` in attesa (la linea tratteggiata nel diagramma). Claude Code riceve il codice e salva il login, su Linux in `.credentials.json` nella cartella dell'account. Una riga incollata non può mai far contattare a ccam un altro host.

ccam stampa questo prima che inizi il login:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

Se riesce, stampa `→ Passed the sign-in to Claude Code on this machine.`

Mentre il login attende un input in un terminale, ccam disattiva l'eco e la modalità di riga del terminale e mostra un `*` per ogni carattere che digiti, così un codice o un indirizzo di callback incollato non compare sullo schermo. ccam gestisce da sé la modifica della riga (`readLine`): `Enter` invia la riga, `Backspace` cancella un carattere e `Ctrl+D` su una riga vuota termina l'input. Il terminale viene ripristinato quando il login finisce, e anche con `Ctrl+C`, `SIGTERM` o `SIGHUP`, che continuano a terminare ccam. Senza un terminale su standard input e standard error, le righe vengono copiate così come arrivano.

### Rinominare, rimuovere e importare

`ccam rename` e `ccam remove` agiscono sulla cartella dell'account e su ciò che ccam conserva altrove con il nome dell'account. ccam elenca questi percorsi in `nameState`:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- le cache dell'utilizzo `~/.local/state/ccam/usage/<name>.json` e `<name>.statusline.json`;
- la registrazione del limite `<name>.limit.json`;
- le note di sessione in `~/.local/state/ccam/handoff/<name>/`.

`ccam remove work` li gestisce in questo ordine:

<!-- motion:accounts-remove -->

1. `ccam remove` per prima cosa si rifiuta se l'account è in esecuzione, poi chiede conferma e ricontrolla dopo la tua risposta (così intercetta una sessione avviata mentre la domanda era aperta).
2. Controlla che `config.toml` si possa salvare prima di modificare qualsiasi cosa.
3. Se l'account condivide le estensioni, ccam riscrive i file `installed_plugins.json` e `known_marketplaces.json` del gruppo (`RewritePluginPaths`), perché Claude Code salva percorsi assoluti dei plugin che possono passare per la cartella di questo account.
4. Poi un account gestito (uno la cui cartella è stata creata da ccam, a differenza di uno registrato con `ccam import`) viene disconnesso: `claude auth logout` viene eseguito nella sua cartella, così nel cestino non resta un login funzionante (né, su macOS, in un elemento del portachiavi con il nome della cartella, che un account creato in seguito con quel nome erediterebbe). Questo passaggio viene saltato per un account con chiave API e quando la cartella di un altro account contiene lo stesso login, che il logout revocherebbe a sua volta. Un logout che fallisce viene segnalato e la rimozione prosegue.
5. ccam sposta la cartella in `~/.local/share/ccam/trash/<name>-<time>/` (`profile.Trash`), l'account esce da `config.toml`, i percorsi di `nameState` vengono eliminati, così un account creato in seguito con lo stesso nome non li eredita, e il link della scorciatoia viene rimosso: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

Se alla fine `config.toml` non si può salvare, la cartella torna fuori dal cestino. Per recuperare un account rimosso, sposta di nuovo la cartella fuori dal cestino ed effettua di nuovo il login.

`ccam rename` rinomina la cartella e riscrive i percorsi dei plugin nel suo `plugins/` dalla vecchia cartella alla nuova. Sposta i percorsi di `nameState` sul nuovo nome e mantiene l'account come predefinito se lo era. Aggiorna anche le note di sessione, così `ccam next --continue` le trova ancora. Il comando scorciatoia mantiene il suo nome. Se la riscrittura dei plugin fallisce, o `config.toml` non si può salvare, la cartella torna al nome precedente.

Su macOS, e su Windows quando Claude Code conserva i login in Gestione credenziali, la cartella stessa mantiene il vecchio nome, perché il login e gli accessi MCP sono conservati sotto il percorso della cartella e una cartella spostata risulterebbe disconnessa. La configurazione registra allora la cartella come `folder` accanto al nuovo nome, e ccam dice `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

`ccam import` registra una cartella dove si trova, come `dir` in `config.toml`. Non la sposta né la copia mai: un refresh token copiato farebbe sì che le copie si disconnettano a vicenda.

Senza una cartella come argomento, cerca le cartelle `~/.claude` e `~/.claude-*` non registrate e chiede conferma per ciascuna. I nomi derivano dalla cartella (`~/.claude` diventa `main`, `~/.claude-work` diventa `work`, `~/.claude-2` diventa `acct2`). Se ccam non riesce a ricavare un nome, o il nome è già occupato, elenca la cartella con il comando `--name` da eseguire. Una cartella trovata che non riesce viene segnalata come `Skipped <path>: <reason>` e le altre vengono importate.

ccam fa prima tutte le domande, poi ricarica `config.toml` e registra le cartelle, così le modifiche salvate nel frattempo da altri comandi ccam vengono mantenute. Le righe `✓` arrivano dopo il salvataggio.

Senza `CLAUDE_CONFIG_DIR`, Claude Code tiene `.claude.json` nella tua cartella home accanto a `~/.claude`, quindi quando importi `~/.claude` e non ha un proprio `.claude.json`, ccam crea il link `~/.claude/.claude.json` verso `~/.claude.json`. In questo modo `claude` da solo e `ccam main` vedono lo stesso account.

Se `~/.claude` ha un proprio `.claude.json` accanto a `~/.claude.json`, ccam li lascia entrambi e lo segnala. `ccam main` userebbe quello interno e `claude` da solo quello esterno, tenendo separati i dati del login e i server MCP. Per tenerne uno solo, sostituisci il file interno con un link (lo segnala anche [doctor](/features/doctor)).

Quando la cartella contiene lo stesso login di un account registrato (una cartella copiata), ccam lo segnala: chi lo rinnova per primo disconnette l'altro, quindi esegui `ccam login <name>` prima di usarla.

Gli account importati non fanno parte di nessun gruppo finché non li aggiungi a uno (`ccam group join <group> <name>`), e non viene creata nessuna scorciatoia.

`ccam completion bash|zsh` stampa uno script di completamento che completa anche i nomi degli account.

## Sicurezza e casi limite

- `ccam rename` e `ccam remove` si rifiutano mentre l'account è in esecuzione (`refuseRunning`), con ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``, perché la sessione continua a scrivere nella cartella. Il controllo copre Claude Code comunque sia stato avviato: tramite ccam, `claude` da solo, un IDE, le sessioni in background e il loro daemon.
- `ccam rename`, `ccam remove` e `ccam logout` si rifiutano mentre un altro processo ccam rinnova o modifica quel login: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- Mentre una modifica dei gruppi è in corso in un altro processo ccam, `ccam rename` e `ccam remove` si rifiutano con `another ccam process is changing the groups right now; try again when it is done`; dopo una rimasta a metà, con ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.
- `ccam remove` su un account importato mantiene la cartella e il suo login, per `claude` da solo. Se l'account fa parte di gruppi, prima li lascia come fa `ccam group leave`: ogni elemento condiviso diventa una copia di quello del gruppo nella cartella, e i server MCP del gruppo finiscono nel suo `.claude.json`. Poi ne viene annullata la registrazione: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. I backup del suo `.claude.json` creati da ccam vengono spostati in `<folder>/.ccam-backup/mcp-<time>`; se non si possono spostare, restano in `~/.local/state/ccam/mcp/<name>/backups` e il messaggio lo dice. La tua `~/.claude` non viene spostata nel cestino.
- Due account non possono avere la stessa cartella: un `config.toml` che indica due volte la stessa cartella non si carica.
- `ccam add` rifiuta un nome la cui cartella esiste già in `profiles/` e rimanda a `ccam import`. Su macOS si ferma anche quando il portachiavi contiene ancora un login o accessi MCP per il percorso della nuova cartella, lasciati da una versione precedente di ccam, e indica il comando `security delete-generic-password` che li rimuove. Se `config.toml` non si può salvare, la nuova cartella viene rimossa di nuovo.
- `ccam import` rifiuta una cartella già registrata, un percorso che non è una cartella e un nome già in uso. ccam avverte anche quando la cartella è una copia di `~/.claude` (lo stesso login), anche se `~/.claude` stessa non è registrata.
- `ccam logout` chiede prima conferma quando l'account è in esecuzione; `ccam login` non controlla.
- Se l'inoltro non riesce a raggiungere Claude Code, per esempio dopo che il login è terminato, stampa `ccam: could not pass the sign-in to Claude Code:` con il motivo. Un errore HTTP di solito significa che il codice è scaduto; avvia di nuovo il login.

## Pagine correlate

- [Avviare Claude Code](/features/launch)
- [Comandi scorciatoia](/features/shortcuts)
- [Condividere server MCP, impostazioni e cronologia](/docs/sharing)
- [Usare più account Claude Code su un solo computer](/guides/multiple-claude-code-accounts)
