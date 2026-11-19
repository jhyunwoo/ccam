---
title: Riferimento dei comandi
metaTitle: Riferimento dei comandi di ccam
description: Tutti i comandi e le opzioni di ccam, da ccam add e ccam <name> fino a ccam usage, ccam next e ccam doctor.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: 28e1ad1dcbad
---

Metti `--dry-run` prima di un comando per vedere cosa farebbe senza eseguirlo, per esempio `ccam --dry-run work`. Funziona quando avvii Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) e con `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` e `doctor --fix`. I comandi che si limitano a leggere, come `list`, `usage`, `share` ed `ext`, lo accettano e comunque non modificano nulla. Gli altri lo rifiutano invece di essere eseguiti, compreso `group rename`.

<!-- commands -->

## Avviare Claude Code

`ccam <name|number> [arguments]` avvia Claude Code con quell'account. ccam non legge gli argomenti che seguono il nome; li passa a `claude` nello stesso ordine, e standard input, codice di uscita e segnali si comportano esattamente come con `claude` stesso. Un `--` dopo il nome è facoltativo.

`ccam` con argomenti ma senza nome di account, come `ccam -p "summarise"`, usa l'account impostato per la cartella corrente, altrimenti l'account predefinito. Impostali con `ccam folder <name> [folder]` e `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` prima del nome dell'account serve per le esecuzioni `-p`/`--print`. Quando l'account raggiunge il suo limite di utilizzo, ccam riesegue lo stesso comando con l'account connesso con più margine e stampa solo l'output dell'esecuzione portata a termine. Vedi [Account per cartella](/features/folder) e [Rieseguire al limite](/features/next#run-again-at-the-limit).

**Posti nei piani Team ed Enterprise.** La scelta automatica dell'account (`ccam next`, `--fallback`, `ccam pick`, il tasto `n` della dashboard e, una volta attivato, [`auto_switch`](/features/next#auto-switch)) passa anche da un posto all'altro della stessa organizzazione Claude Team o Enterprise, e questo può violare i termini di Anthropic. Leggi le [domande frequenti](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e verifica con la tua organizzazione e con Anthropic prima di farci affidamento.

Una scorciatoia come `ccwork` avvia il suo account allo stesso modo: `ccwork --resume` fa ciò che fa `ccam work --resume`. Vedi [Scorciatoie](/docs/commands#shortcuts) più avanti.

## Aggiungere e rimuovere account

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

`ccam import` senza una cartella cerca `~/.claude` e `~/.claude-*` e chiede conferma per ciascuna; una cartella che non riesce viene segnalata come `Skipped <path>: <reason>` e le altre vengono registrate. Se il nome di una cartella è già occupato, importa quella cartella da sola con `ccam import <folder> --name <name>`. Le cartelle importate vengono registrate dove si trovano; ccam non sposta né copia mai un login.

`ccam login --expired` connette anche gli account il cui login è stato cancellato da Claude Code. `ccam remove` disconnette un account gestito (uno la cui cartella è stata creata da ccam, a differenza di uno registrato con `ccam import`) prima che la sua cartella finisca nel cestino, a meno che un altro account non abbia lo stesso login. Rimuovere un account importato lo toglie dai suoi gruppi e lascia nella sua cartella una copia di ogni elemento condiviso. Su macOS, `ccam rename` mantiene il nome della cartella, perché il portachiavi conserva il login sotto il percorso della cartella.

`ccam add --api` aggiunge un account che usa una chiave API invece di un login di Claude; ccam legge la chiave a ogni avvio e non la salva mai. Vedi [Account con chiave API](/features/api).

`ccam rename`, `ccam remove` e i comandi `ccam group` si rifiutano di modificare un account mentre Claude Code è in esecuzione con quell'account, comunque sia stato avviato (comprese le sessioni in background). Chiudi prima la sessione. Si rifiutano anche finché una modifica dei gruppi è incompiuta: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` stampa ogni account con il suo slot, la scorciatoia, l'email e la cartella. Il suo array `groups` elenca i gruppi dell'account, e il suo oggetto `share` associa ogni elemento (`mcp`, `settings`, `memory`, `extensions`, `history`) al gruppo da cui proviene, oppure a `local`.

## Scorciatoie

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Una scorciatoia è un comando a cui dai tu il nome e che avvia Claude Code con un determinato account. Come `ccam work`, passa ogni argomento a Claude Code invariato, quindi `ccwork --resume` e `ccwork -p "hi"` funzionano come ti aspetti. È un link a ccam in `~/.local/bin`.

Le scorciatoie sono facoltative. La procedura guidata di `ccam add` chiede se crearne una, `ccam add <name> --shortcut <command>` ne crea una senza chiedere e `ccam list` mostra la scorciatoia di ogni account, oppure `-` se non ne ha.

I nomi usano i caratteri a-z, 0-9, `.`, `-` e `_`. Un nome non può essere un comando di ccam, la scorciatoia di un altro account o il nome di un altro programma nel tuo `PATH`, né un file in `~/.local/bin` che non ha creato ccam, come una tua copia o un tuo link di ccam. Una scorciatoia modificata a mano in `config.toml` che viola queste regole, che è `claude` o `ccam`, o che è assegnata a due account impedisce di caricare `config.toml`, e l'errore indica l'account.

`ccam rename` mantiene il nome della scorciatoia, e `ccam remove` elimina la scorciatoia insieme all'account. `ccam shims install` ricrea i link di tutti gli account, per esempio dopo che hai spostato ccam; un file che ostacola una scorciatoia non blocca le altre.

## Condivisione

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Un gruppo condivide degli elementi tra i suoi account; un account può stare in più gruppi, ma ogni elemento proviene da un solo gruppo. Vedi [condivisione con i gruppi](/features/share).

## Copiare una configurazione esistente

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` copia `settings.json`, `CLAUDE.md`, le estensioni e i server MCP utente da una cartella di Claude Code in un gruppo, così ogni account che condivide quegli elementi li riceve. La cartella viene solo letta.

- `--settings`, `--memory`, `--extensions`, `--mcp` e `--history` copiano solo gli elementi indicati. Senza di esse, in un terminale una checklist ti permette di scegliere, con tutto selezionato tranne la cronologia delle conversazioni; altrove viene copiato tutto tranne la cronologia delle conversazioni.
- `--mcp` funziona anche quando `~/.claude` è vuota o manca ma `~/.claude.json` ha dei server.
- `-g` o `--group` sceglie il gruppo.
- `-y` o `--yes` copia senza chiedere e mantiene la copia del gruppo quando un file è diverso.

Gli accessi MCP non vengono copiati, perché un accesso copiato disconnette gli altri che lo possiedono; il piano indica gli account che devono accedere con `/mcp`. Vedi [copiare una configurazione esistente](/docs/sharing#copy-an-existing-setup).

## Server MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` sposta nel suo gruppo i server che ha solo quell'account, compresi quelli che ha per una sola cartella (il semplice `claude mcp add`), che da quel momento si caricano in ogni cartella. Un server il cui nome è già usato dal gruppo per una definizione diversa viene saltato; aggiungi `--force` per sostituire la copia del gruppo, di cui ccam fa prima un backup.

## Utilizzo e cambio di account

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

- `next --auto`, `--fallback` e `ccam pick` scelgono allo stesso modo: prima gli account sotto `auto_max_usage` (`config.toml`, predefinito 100), secondo `auto_strategy` (`headroom`, `expiring` o `balance`), contando la finestra settimanale di `--model`. Non scelgono mai un account al 100% (una finestra che Claude Code ha segnalato piena conta come 100% fino al suo reset) né un login con chiave API, e restano dentro l'organizzazione dell'account da cui partono (`--cross-org` o `cross_org = true` per cercare oltre).
- `--wait` prende il suo valore dopo `=` (`--wait=2h`); scritto con uno spazio viene rifiutato con un suggerimento. Quando non possono scegliere nessun account escono con codice 75 oppure, con `--wait`, attendono il reset del primo account.
- `ccam next --continue` porta una conversazione nell'altro account anche tra account che non condividono la cronologia: in quel caso Claude Code la copia nel nuovo account. Con `--brief`, il nuovo account avvia invece una nuova conversazione, da un riassunto dell'ultima che ccam scrive senza un modello ([`--brief`](/features/next#brief)).
- A differenza di `ccam <name>`, `ccam next` legge le proprie opzioni, quindi le opzioni di Claude Code vanno dopo `--`; senza, ccam si ferma con `unknown flag`.
- `ccam agent` esegue un'attività con `claude -p` in questa cartella con un account diverso da quello da cui viene chiamato, nella modalità dei permessi e con il livello di effort della sessione chiamante, e stampa l'ultima risposta di Claude. Al limite di quell'account la conversazione prosegue con un altro ([Subagenti su altri account](/features/agent)).
- Dentro Claude Code, il comando `/ccam:switch [account]` del plugin passa una conversazione a un altro account, e `ccam hook` è ciò che eseguono i suoi hook.
- Con `auto_switch = true` in `config.toml` (Linux e macOS; disattivato per impostazione predefinita), l'hook del limite chiude anche una sessione interattiva avviata da ccam e la continua su un altro account con margine, ma solo su uno dello stesso tipo: un altro posto della stessa organizzazione, o un altro piano personale ([`auto_switch`](/features/next#auto-switch)).

Vedi [Passare all'account successivo](/features/next), [Scegliere un account per gli script](/features/pick), [Subagenti su altri account](/features/agent) e [Sessioni in background di ogni account](/features/agents).

## Manutenzione

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

`ccam update` verifica il checksum SHA-256 della release prima di sostituire ccam, e `--version v0.2.0` installa una release specifica. Una copia di ccam installata dal pacchetto `.deb` si aggiorna invece con apt. Vedi [aggiornamento](/docs/install#update).

`ccam uninstall` elenca cosa rimuove e chiede prima conferma: ogni account come fa `ccam remove`, i comandi scorciatoia, le cartelle di ccam e il comando `ccam`. Una cartella importata resta dove si trova. Vedi [disinstallazione](/docs/install#uninstall).

## File

| Percorso | Contenuto |
| --- | --- |
| `~/.config/ccam/config.toml` | Account, gruppi e impostazioni (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). Le impostazioni vanno in cima, sopra il primo `[[group]]`, `[[profile]]` o `[[folder]]`: una che sta sotto una tabella appartiene a quella tabella, quindi ccam la segnala a ogni avvio e in `ccam doctor`, e i comandi che salvano la configurazione si rifiutano finché non viene spostata. I commenti non vengono conservati quando ccam salva il file |
| `~/.local/share/ccam/profiles/<name>/` | La `CLAUDE_CONFIG_DIR` di ogni account |
| `~/.local/share/ccam/groups/<group>/` | Impostazioni, server MCP, estensioni e cronologia condivisi |
| `~/.local/state/ccam/` | Cache dell'utilizzo, registrazioni dei limiti, lock, note e altro stato (vedi l'elenco sotto la tabella) |
| `~/.local/bin/<shortcut>` | Scorciatoie come `ccwork`, link a ccam |

La cartella di stato contiene:

- la cache dell'utilizzo e le registrazioni dei limiti (`usage/<name>.json`, `usage/<name>.limit.json`);
- gli orari dell'ultimo uso;
- gli indicatori di sessione e i lock (`locks/`);
- le scorciatoie create da ccam (`shortcuts.json`);
- una modifica dei gruppi non terminata (`group-change.json`);
- le note lasciate da `/ccam:switch` e dall'hook del limite (`handoff/<account>/<session>.json`);
- quando `ccam pick` o `auto_switch` hanno scelto per l'ultima volta ciascun account (`picks/<name>`; una scelta conta per 2 minuti);
- i riassunti di `--brief`, conservati per 7 giorni (`brief/<session>.md`).

Impostare `XDG_CONFIG_HOME`, `XDG_DATA_HOME` o `XDG_STATE_HOME` cambia la posizione di queste cartelle. Su Windows, `~` è `%USERPROFILE%` e le scorciatoie sono `<shortcut>.exe`.
