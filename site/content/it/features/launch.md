---
title: Avviare Claude Code
description: Avvia Claude Code con qualsiasi account usando ccam <name|number>, ccam run o la dashboard; gli argomenti passano invariati e --dry-run mostra il comando.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: 11df08608a1e
---

`ccam work` avvia Claude Code con l'account `work`, e tutto ciò che segue il nome arriva a `claude` invariato. Usalo ovunque scriveresti `claude`: in un terminale, in uno script o con un input passato tramite pipe.

## Utilizzo

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

Indica il nome di un account o il suo numero in `ccam list`. ccam non legge gli argomenti che seguono il nome, quindi `ccam work --help` mostra l'aiuto di Claude Code. Un `--` dopo il nome è facoltativo. Standard input, codice di uscita e segnali si comportano come con `claude`.

Quando il primo argomento inizia con `-` e non è una delle opzioni di ccam, ccam sceglie l'account al posto tuo e passa tutti gli argomenti a Claude Code. Usa l'account impostato per la cartella corrente con [`ccam folder`](/features/folder), altrimenti l'account predefinito: quello impostato con `ccam use`, oppure, se non ne hai impostato uno, l'account con il numero più basso.

### La dashboard

`ccam` senza argomenti apre la [dashboard](/features/dashboard) quando standard input e standard output sono un terminale, altrimenti stampa `ccam list`. Spostati su un account (con le frecce, o digitando il suo numero) e premi `Enter` per avviarlo; `n` avvia l'account con più margine e `o` propone `--continue`, `--resume` o altri argomenti. La dashboard si chiude prima che Claude Code si avvii, così il terminale viene ripristinato prima.

**Posti nei piani Team ed Enterprise.** La scelta automatica dell'account (`ccam next`, `--fallback`, `ccam pick`, il tasto `n` della dashboard e, una volta attivato, [`auto_switch`](/features/next#auto-switch)) passa anche da un posto all'altro della stessa organizzazione Claude Team o Enterprise, e questo può violare i termini di Anthropic. Leggi le [domande frequenti](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e verifica con la tua organizzazione e con Anthropic prima di farci affidamento.

Quando nella shell è impostata `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` o `ANTHROPIC_AUTH_TOKEN`, Claude Code potrebbe usare quella chiave o quel token invece del login dell'account, quindi ccam ti avvisa sullo standard error prima dell'avvio (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Si limita a leggere le variabili; un [account con chiave API](/features/api) invece le rimuove.

Per impostazione predefinita, ccam si sostituisce con Claude Code, quindi non resta in esecuzione nient'altro. Con `relay = true` (oppure [`auto_switch = true`](/features/next#auto-switch) su Linux e macOS) in cima a `config.toml`, sopra il primo `[[group]]`, `[[profile]]` o `[[folder]]`, ccam avvia invece le sessioni interattive come processo figlio e resta in esecuzione accanto a esse. Questo permette a `/ccam:switch`, o ad `auto_switch` al raggiungimento di un limite di utilizzo, di continuare da solo la conversazione su un altro account, con le opzioni, la modalità dei permessi, il livello di effort e l'agente della sessione; prima che `/ccam:switch` la sposti su un account di un'altra organizzazione, ccam te lo chiede ([dettagli](/features/next#relay)).

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--dry-run` | Stampa l'ambiente e il comando che ccam eseguirebbe, poi esce | disattivata |
| `--fallback` | Per le esecuzioni `-p`/`--print`: quando l'account raggiunge il suo limite di utilizzo, riesegue lo stesso comando con l'account connesso con più margine ([dettagli](/features/next#run-again-at-the-limit)) | disattivata |
| `-v`, `--version` | Stampa `ccam <version>` (solo come primo argomento) | |
| `-h`, `--help` | L'aiuto di ccam, come `ccam --help` o `ccam run --help` | |

Regole per gli argomenti:

- ccam legge le proprie opzioni solo prima del nome dell'account. Dopo il nome, ogni argomento appartiene a Claude Code, compresi `--dry-run` e `--help`.
- `--dry-run` funziona come `ccam --dry-run work`, `ccam run --dry-run work` e `ccam --dry-run` (per l'account della cartella o quello predefinito). `--dry-run` e `--fallback` possono comparire in qualsiasi ordine.
- `ccam run <name|number> [claude args]` è la forma estesa di `ccam <name>`. `ccam run` senza un nome usa l'account della cartella o quello predefinito.
- I nomi degli account non possono essere nomi di comandi di ccam come `add`, `list`, `run` o `usage`, così `ccam <word>` non è mai ambiguo.

Il dry run stampa una riga che potresti incollare in una shell, o in PowerShell su Windows. Comincia con ciò che l'avvio cambia nell'ambiente: sempre `CLAUDE_CONFIG_DIR` e `CCAM_PROFILE`, `env -u <name>` per una variabile che rimuove (come un `CCAM_RELAY` ereditato) e, per un [account con chiave API](/features/api), il suo endpoint, il suo modello e da dove arriva la chiave (`ANTHROPIC_API_KEY="$WORK_KEY"`), mai la chiave stessa:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

La riga di PowerShell rimette a posto le variabili alla fine, come fa `VAR=x cmd` in una shell, e passa a Claude Code la sua riga di comando tramite `Start-Process`, così Windows PowerShell 5.1 e PowerShell 7 la eseguono allo stesso modo. Quando un avvio imposta una variabile a un valore vuoto, una riga `#` sotto di essa avverte che PowerShell prima della 7.5 rimuove invece una variabile del genere. Per un `claude.cmd` installato da npm, la riga esegue il comando `cmd.exe` che userebbe l'avvio (vedi sotto), e un argomento che l'avvio rifiuterebbe lì fa fallire anche il dry run.

## Come funziona

Ogni avvio segue gli stessi passaggi, dalla riga di comando fino a `exec`:

<!-- motion:launch -->

Ogni riquadro si illumina mentre viene eseguito il suo passaggio, e le righe sottostanti sono l'ambiente e gli argomenti con cui parte Claude Code.

1. ccam parte in `app.Main`, chiamata da `cmd/ccam/main.go`. Se ccam è stato avviato con un altro nome, si tratta di una [scorciatoia](/features/shortcuts) e si passa direttamente all'avvio.
2. ccam esamina `os.Args` prima che entri in gioco cobra (`app.parse`). Solo `ccam` da solo (la dashboard), un nome di comando di ccam (`config.Reserved`), `-h`/`--help` e `-v`/`--version` passano a `cli.Execute`. Tutto il resto è un avvio.
3. ccam legge `~/.config/ccam/config.toml`, l'unico file che analizza (`launcher.Launch`).
4. ccam trova l'account per nome o per numero, scarta un eventuale `--` iniziale e individua l'eseguibile `claude` (`launcher.Prepare`, `ResolveClaude`). Senza nome, prende l'account della cartella, confrontando la directory di lavoro con le voci `[[folder]]` di `config.toml` come stringhe, senza distinguere maiuscole e minuscole su macOS e Windows (`HereProfile`). La ricerca nel PATH viene salvata in cache in `~/.local/state/ccam/claude-path` insieme all'ora di modifica del file e a un hash del `PATH`, così ccam non cerca nel `PATH` a ogni avvio. Un file modificato o mancante, o un `PATH` cambiato, provoca una nuova ricerca. Un percorso `claude_bin` in `config.toml` (`~/` viene espanso) deve essere un file eseguibile; altrimenti l'avvio si ferma con `claude_bin X is not an executable file. Check claude_bin in config.toml` prima che venga scritto qualsiasi cosa. ccam rifiuta un `claude` che è ccam stesso, confrontando l'inode, perché si creerebbe un ciclo.
5. ccam confronta la cartella dell'account `~/.local/share/ccam/profiles/<name>/` con le sue impostazioni di condivisione, chiamando solo `lstat` e `readlink`, una o due chiamate di sistema per elemento (`checkLinks`, `share.Check`). Un link mancante viene creato sul momento. Qualsiasi cosa che potrebbe far perdere dati interrompe l'avvio.
6. ccam aggiunge i server MCP condivisi (`applyMCP`; vedi sotto).
7. ccam registra la sessione (`session.Register`). Scrive un marcatore con il nome del suo PID, che contiene l'ora di avvio del processo, in `~/.local/state/ccam/locks/profile-<name>/<pid>`, qualunque cosa dica `XDG_RUNTIME_DIR`, così ogni shell vede gli stessi marcatori. Per un account che condivide la cronologia delle conversazioni, ccam cerca prima una sessione in esecuzione in `history-<group>-<hash>`, dove l'hash è quello della cartella corrente con i link simbolici risolti, e si registra anche lì (`historyCheck`).

   `exec` mantiene il PID, quindi il marcatore indica poi il processo di Claude Code stesso. ccam usa un file marcatore invece di un lock mantenuto, perché i processi avviati da Claude Code erediterebbero il lock e lo conserverebbero dopo la fine della sessione. Un marcatore che non si riesce a scrivere non ferma l'avvio: ccam dice `ccam: could not mark the session as running (…); starting Claude Code anyway`.

   Gli altri controlli di ccam che verificano se un account è in esecuzione (rinomina, rimozione, modifiche dei gruppi, `doctor --fix`, la sincronizzazione MCP) leggono anche i record propri di Claude Code nella cartella dell'account (`sessions/<pid>.json`, `daemon.lock`), quindi vedono anche Claude Code avviato senza ccam: `claude` da solo, un IDE, le sessioni in background e il loro daemon.
8. ccam imposta il titolo del terminale su `claude · <name>` (quando lo standard error è un terminale), scrive l'ora di avvio in `~/.local/state/ccam/last-used/<name>` senza `fsync` e chiama `syscall.Exec` (`Plan.Exec`).

Claude Code viene eseguito con `CLAUDE_CONFIG_DIR` impostata sulla cartella dell'account e `CCAM_PROFILE` impostata sul nome dell'account. ccam rimuove gli eventuali valori ereditati di queste variabili, e un `CCAM_RELAY` ereditato. Per una cartella `~/.claude` registrata con `ccam import` su macOS, imposta anche `CLAUDE_SECURESTORAGE_CONFIG_DIR` a un valore vuoto, così l'account e `claude` da solo usano lo stesso login nel Portachiavi; su Windows fa lo stesso quando Claude Code conserva il login in Gestione credenziali.

### Perché exec

ccam sostituisce il proprio processo con `claude` (`syscall.Exec`). Di conseguenza:

- ccam non resta in memoria;
- tasti premuti, segnali e codice di uscita sono quelli di Claude Code;
- gli argomenti restano un array, mai una stringa di shell.

Windows non ha `exec`. Lì ccam avvia `claude.exe` come processo figlio con lo stesso ambiente, gli stessi argomenti e la stessa console, lo attende ed esce con il suo codice di uscita (`Plan.Exec`). Un `claude.cmd` installato da npm viene eseguito in `cmd.exe` con il percorso completo e ogni argomento tra virgolette, così che `cmd.exe` non ne legga nessuno come propria sintassi (`&`, `%PATH%`); un argomento con un a capo non si può passare in quel modo e viene rifiutato, con il consiglio di installare il Claude Code nativo.

Ctrl+C arriva direttamente a Claude Code dalla console, quindi ccam lo ignora. Il marcatore di sessione porta il PID di ccam, che vive esattamente quanto Claude Code, e viene rimosso quando Claude Code termina. Mentre attende, ccam non fa nient'altro.

Il lavoro da fare dopo una sessione, come accorgersi dei server aggiunti con `claude mcp add`, avviene all'avvio successivo: ccam confronta l'ora di modifica di `.claude.json` con `~/.local/state/ccam/mcp/<name>/seen` e analizza il file solo se è cambiato.

### Come viene instradata una riga di comando

Tre righe di comando, tre percorsi attraverso `internal/app`:

<!-- motion:launch-route -->

Una scorciatoia viene riconosciuta da `argv[0]`, `ccam <name>` dal primo argomento, ed entrambe terminano in `syscall.Exec`. Solo un comando di ccam come `ccam usage` arriva a cobra.

### Server MCP condivisi

Con la strategia predefinita `flag`, ccam passa i server MCP del gruppo a Claude Code come argomento, posto prima dei tuoi:

<!-- motion:launch-mcp -->

ccam aggiunge l'argomento solo se il `mcp.json` del gruppo esiste, e lo mette subito dopo `claude` (`mcp.FlagArgs`, `Plan.InsertArgs`). Le tue opzioni `--mcp-config` o `--strict-mcp-config` continuano quindi a valere in aggiunta. È un solo argomento, `--mcp-config=<path>`, perché l'opzione `--mcp-config` di Claude Code accetta più valori: scritta come due argomenti, leggerebbe come file di configurazione anche un prompt che segue.

ccam decide dove l'opzione non può andare (`claudeVerb`). Legge il sottocomando di Claude Code dagli argomenti come fa Claude Code, saltando le opzioni insieme ai loro valori. `remote-control` o `rc` non ricevono l'opzione (vedi sotto), e nemmeno un verbo delle sessioni in background indicato per primo (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`). Non la ricevono neppure `ccam work mcp list`, `get` e `remove`, perché lì Claude Code la rifiuta; ccam dice invece sullo standard error che i server del gruppo si caricano anche nelle sessioni dell'account:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Ogni altro sottocomando e ogni sessione mantengono l'opzione.

Con `mcp_strategy = "sync"` in `config.toml`, ccam scrive invece i server nel `.claude.json` dell'account, e salta questo passaggio se l'account è già in esecuzione. Un avvio salta il lavoro anche quando né il `mcp.json` del gruppo né `.claude.json` sono cambiati dall'ultimo controllo. Se l'account ha un proprio server con lo stesso nome di un server del gruppo, la sincronizzazione mantiene il server dell'account e lo segnala una volta. La sincronizzazione rimette anche una copia sparita da `.claude.json` (`put back group default's MCP servers …`), e modifica o rimuove solo le copie che ha scritto e che nessuno ha cambiato da allora.

### Remote Control

`ccam work rc` (o `ccam work remote-control`) avvia il Remote Control di Claude Code come `work`. `claude remote-control` rifiuta `--mcp-config` prima del suo verbo e non ne accetta dopo, e le sue sessioni leggono i server MCP solo dal `.claude.json` dell'account. Per questo qui la strategia flag funziona in modo diverso:

<!-- motion:launch-rc -->

1. ccam trova `remote-control` o `rc` come verbo di Claude Code, anche dopo opzioni e i loro valori (`ccam work --model opus rc`), e non aggiunge l'opzione (`claudeVerb`, `mcp.FlagArgs`).
2. ccam confronta l'ora di modifica del `mcp.json` del gruppo con il `sync-stamp` dell'account, e controlla che ogni server del gruppo sia ancora in `.claude.json` (`mcp.Synced`). Se `.claude.json` contiene già i server attuali del gruppo, o il gruppo non ne ha, Remote Control si avvia senza domande.
3. Altrimenti ccam elenca i server che mancherebbero alle sessioni e chiede:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. Con `y`, ccam li copia come fa la strategia `sync` (`mcp.Sync`): con il lock di Claude Code, dopo un backup, mantenendo i server propri dell'account. Con `N`, o senza un terminale, Remote Control si avvia senza di loro, e ccam suggerisce `mcp_strategy = "sync"`, che li copia prima di ogni avvio.

ccam chiede a ogni avvio di Remote Control finché i server non vengono copiati. Dopo una copia, chiede di nuovo solo quando il `mcp.json` del gruppo cambia o un server copiato sparisce da `.claude.json`. Una copia in `.claude.json` prevale su `--mcp-config` nei successivi avvii normali, quindi `ccam mcp add`, `remove` e `import` aggiornano le copie che hanno gli altri account ([Server MCP condivisi](/features/mcp)).

### Il limite di 10 ms

Il budget per il percorso di avvio è di 10 ms al 95° percentile; in pratica ccam aggiunge circa 4 ms prima che Claude Code si avvii. Il percorso di avvio non costruisce mai l'albero dei comandi di cobra, non avvia Bubble Tea e non apre connessioni di rete, e `TestPackageInitIsCheap` fallisce se le funzioni `init()` dei pacchetti richiedono in totale più di 5 ms. `scripts/perf-check.sh` misura `ccam --dry-run work`, `ccam work` e una scorciatoia con `hyperfine` e fallisce sopra i 10 ms, e `BenchmarkLaunchPath` in `internal/launcher` esegue i passaggi in un ciclo. Per vedere il tempo di ogni passaggio, imposta `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Sicurezza e casi limite

- Un nome o un numero sconosciuto si ferma con `no account "<name>"` e rimanda a `ccam list`. Se non ci sono account, ccam ti chiede di aggiungerne uno con `ccam add <name>`.
- Se un percorso condiviso è diventato un file normale o punta al posto sbagliato, ccam si ferma con `stopped because of share link problems:`, elenca i problemi e ti chiede di eseguire `ccam doctor --fix`. Sul percorso di avvio non sposta dati.
- Se un'altra sessione di Claude Code che condivide la cronologia delle conversazioni con questo account (stesso gruppo) è in esecuzione nella stessa cartella, ccam avvisa che la cronologia potrebbe mescolarsi e chiede `Run anyway? [y/N]`. La cartella viene confrontata con i link simbolici risolti. Senza un terminale (uno standard input da `/dev/null` conta come assente), o da un job in background, mostra l'avviso e continua. Una nuova esecuzione `-p` inizia una conversazione propria e non riceve l'avviso; `--version`, `--help` e i sottocomandi che non aprono una conversazione, come `mcp`, `auth` o `plugin`, non vengono controllati affatto.
- `--dry-run` salta la sincronizzazione MCP, l'avviso sui nuovi server, la nota su `mcp list`, il controllo della cronologia, il marcatore di sessione, il titolo e l'ora dell'ultimo uso, e non esegue exec. Controlla comunque i link di condivisione e crea un link mancante come farebbe un avvio normale.
- `ccam work rc` mentre `work` è già in esecuzione lascia stare `.claude.json`, lo segnala e avvia Remote Control senza i server del gruppo; eseguilo di nuovo dopo la fine della sessione per copiarli. Con `--dry-run` stampa il messaggio ma non chiede e non copia.
- Con `mcp_strategy = "sync"`, una sincronizzazione fallita stampa `ccam: MCP sync failed (launching anyway): …` e Claude Code si avvia con i server che aveva già.

## Pagine correlate

- [Comandi scorciatoia](/features/shortcuts)
- [Usare Remote Control con ccam](/guides/remote-control)
- [Gestire gli account](/features/accounts)
- [Riferimento dei comandi](/docs/commands)
- [Condividere server MCP, impostazioni e cronologia](/docs/sharing)
- [Come funziona CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
