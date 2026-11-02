---
title: Passare all'account successivo
description: ccam next apre Claude Code con un altro account; --auto sceglie quello con più margine, --continue porta con sé la conversazione e --fallback riesegue i job -p.
updated: 2026-11-03
order: 8
translation: machine
sourceHash: "89e8dea9c55d"
---

`ccam next` apre Claude Code con un altro account senza che tu debba indicarne il nome. Usalo quando un account raggiunge il suo limite di utilizzo. Ci sono quattro modi per cambiare account:

- `ccam next` da solo prende l'account che segue quello usato per ultimo, e `--auto` prende quello con più margine. `--continue` continua l'ultima conversazione di questa cartella sul nuovo account, che i due condividano o no la cronologia.
- Dentro Claude Code, il [plugin di ccam](#inside-claude-code) ti avvisa quando un limite ferma una sessione, e `/ccam:switch` passa la conversazione a un altro account.
- Per gli script che eseguono `claude -p`, `ccam --fallback` fa il cambio da solo (vedi [Rieseguire al limite](#run-again-at-the-limit)).
- Con [`auto_switch`](#auto-switch) attivo (Linux e macOS), ccam sposta anche una sessione interattiva che Claude Code ferma al suo limite.

**Posti nei piani Team ed Enterprise.** La scelta automatica dell'account (`ccam next`, `--fallback`, `ccam pick`, il tasto `n` della dashboard e, una volta attivato, [`auto_switch`](/features/next#auto-switch)) passa anche da un posto all'altro della stessa organizzazione Claude Team o Enterprise, e questo può violare i termini di Anthropic. Leggi le [domande frequenti](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e verifica con la tua organizzazione e con Anthropic prima di farci affidamento.

## Utilizzo

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

ccam stampa il cambio come `→ <from> → <to>` sullo standard error, con il motivo di una scelta `--auto`, poi avvia Claude Code con il nuovo account come fa `ccam <name>`. `<from>` è l'account usato per ultimo, oppure, con `--continue`, l'account a cui appartiene la conversazione. Gli argomenti dopo `--` arrivano a `claude` invariati.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--auto` | Sceglie l'account con login attivo in base all'utilizzo (vedi sotto), saltando gli account al 100% | disattivata |
| `--max-usage <1-100>` | Con `--auto`: preferisce gli account sotto questo utilizzo | `auto_max_usage`, altrimenti 100 |
| `--strategy headroom\|expiring\|balance` | Con `--auto`: l'utilizzo su 5 ore più basso, l'azzeramento su 5 ore più vicino oppure la maggiore quota settimanale residua per giorno | `auto_strategy`, altrimenti `headroom` |
| `--cross-org` | Con `--auto`: sceglie anche account di un'altra organizzazione (vedi sotto) | `cross_org`, altrimenti disattivata |
| `--wait[=<duration>]` | Con `--auto`: quando nessun account ha margine, attende l'azzeramento del primo, al massimo per questo tempo. Il valore va dopo `=`: `--wait 2h` viene rifiutato con il suggerimento `--wait takes its value after '=': --wait=2h` | disattivata; `5h` se indicata da sola |
| `-c`, `--continue` | Continua l'ultima conversazione di questa cartella sul nuovo account (vedi sotto) | disattivata |
| `--brief` | Con `--continue`: avvia con il nuovo account una nuova conversazione dal riassunto che ccam fa dell'ultima, invece di riprenderla tutta ([vedi sotto](#brief)) | disattivata |
| `-- <claude args>` | Argomenti passati a `claude` dopo quelli di ccam | nessuno |
| `--dry-run` | Supportata: stampa `→ <from> → <to>` e la riga di comando che eseguirebbe, senza avviare Claude Code né registrare un avvio | disattivata |

`--auto` e `--continue` funzionano insieme. Con `--continue`, gli argomenti di `--resume` vengono messi prima dei tuoi.

Le opzioni proprie di Claude Code, come `--dangerously-skip-permissions` o `--model`, vanno dopo `--`. A differenza di `ccam <name>`, che passa oltre tutto ciò che segue il nome, `ccam next` legge da sé i propri argomenti, quindi senza `--` si ferma con `ccam: unknown flag: --dangerously-skip-permissions`.

## Come funziona

Ogni avvio tramite ccam scrive l'ora in `~/.local/state/ccam/last-used/<name>`. `ccam next` legge questi file, controlla quali account hanno effettuato il login e sceglie.

1. ccam costruisce un `next.Candidate` per ogni account a partire da `config.toml` (`cli.candidates`): l'account, l'ora del suo ultimo avvio (`launcher.LastUsed`) e se ha un login funzionante (`creds.Has`: un login di Claude in `.credentials.json` o, su macOS, nel Portachiavi, che non sia scaduto del tutto né sia stato cancellato da Claude Code). Una cartella che usa invece una chiave API (un login della Console, o un account aggiunto con `ccam add --api`) non ha limiti di piano e non viene mai scelta.
2. Senza `--auto`, ccam ordina gli account per numero di slot e trova quello avviato più di recente (`next.Rotate`). È l'account da cui stai passando a un altro. Se nessun account è ancora stato avviato, è lo slot più basso.
3. Da lì procede in avanti nell'ordine degli slot, ricominciando dal primo dopo l'ultimo, e restituisce il primo account con un login.
4. ccam prepara e avvia l'account scelto esattamente come farebbe `ccam <name>` (`launcher.LaunchWith`): controlla i link dell'account, applica la condivisione degli MCP, registra la sessione e sostituisce ccam con `claude` tramite `syscall.Exec`.

Il percorso parte dall'account usato per ultimo e si ferma al primo account con login che lo segue:

<!-- motion:next-rotate -->

`work` è stato avviato più di recente, quindi il percorso parte da quello successivo. `side` non ha un login e viene saltato, e viene scelto `lab`. Se anche `lab` fosse stato disconnesso, il percorso sarebbe ricominciato da `personal` (`% len` nel diagramma). La rotazione non restituisce mai l'account da cui sei partito: se non c'è un altro account con login, si ferma con un errore.

### --auto

`--auto` decide in base all'utilizzo invece che all'ordine degli slot. ccam tiene gli account con login, chiede l'utilizzo di tutti contemporaneamente nello stesso modo di `ccam usage` (`Fetcher.FetchAll`) e sceglie secondo le regole qui sotto (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

Un account viene valutato in base alle finestre di 5 ore e settimanale e, quando `--model` è tra gli argomenti per Claude Code (dopo `--`), anche in base alla finestra settimanale di quel modello (`opus`, `sonnet`, `haiku` o `fable` nel nome). ccam divide gli account in tre livelli (`Policy.Choose()`):

1. **Con margine**: ogni finestra sotto `auto_max_usage` (`config.toml`, predefinito 100). La strategia sceglie tra questi: `headroom` (predefinita) prende l'utilizzo su 5 ore più basso; `expiring` prende la finestra di 5 ore che si azzera prima, così la sua quota residua viene usata prima di andare persa; `balance` prende l'account con più quota settimanale residua per ogni giorno che manca all'azzeramento della sua settimana, contando a suo carico ogni sessione ccam in esecuzione, così la quota settimanale non resta inutilizzata e le sessioni in parallelo si distribuiscono.
2. **Oltre la soglia**: sotto il 100% ma pari o superiore a `auto_max_usage`. Usati solo quando nessun account ha margine; viene prima quello con l'utilizzo su 5 ore più basso.
3. **Pieni**: una finestra al 100%, oppure utilizzo su 5 ore sconosciuto. Mai scelti.

Una finestra conta come 100% fino al suo reset quando Claude Code vi ha rifiutato una richiesta: l'hook del [plugin di ccam](#inside-claude-code) lo registra (`~/.local/state/ccam/usage/<name>.limit.json`), qualunque cosa dicano i numeri salvati. Così un account che ha appena raggiunto il suo limite non viene scelto di nuovo prima del reset, anche quando i suoi numeri sono vecchi.

Le impostazioni per tutti gli account vanno in cima a `config.toml`, sopra il primo `[[group]]`, `[[profile]]` o `[[folder]]`. Un'impostazione sotto uno di questi appartiene a quella tabella: ccam la segnala a ogni avvio, e i comandi che salvano la configurazione si rifiutano finché non viene spostata.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organizzazioni.** Le scelte automatiche restano tra account dello stesso tipo di quello che lasci: posti della stessa organizzazione Team o Enterprise, oppure piani personali (Pro, Max), così il lavoro di un account aziendale non prosegue su uno privato né viceversa. ccam conosce l'organizzazione di ogni account dal suo login (`ccam login` la registra; per i login più vecchi viene letta da `.claude.json`); un account di cui ccam non riesce a sapere l'organizzazione può essere abbinato a qualsiasi account. Due cartelle con il login nello stesso account Claude ne condividono i limiti, quindi una non subentra mai all'altra. `--cross-org`, oppure `cross_org = true` in `config.toml`, toglie la regola dell'organizzazione, ma dopo [`/ccam:switch`](#inside-claude-code) uno spostamento verso un'altra organizzazione aspetta comunque il tuo sì.

Quando non si può scegliere nulla, ccam esce con codice 75 (`EX_TEMPFAIL`), così uno script può distinguere "riprova più tardi" da un errore. Con `--wait`, invece, attende finché la finestra bloccante del primo account pieno non si azzera (30 secondi dopo l'orario di azzeramento, quando l'endpoint mostra la nuova finestra), stampa `ccam: no account has room; waiting until 15:10` e sceglie di nuovo; quando nessun account dice quando si azzera, ricontrolla dopo 10 minuti. Se nulla si libera prima della fine dell'attesa, esce con 75.

La richiesta segue le regole di `ccam usage`: i numeri ottenuti meno di 60 secondi prima vengono riutilizzati, la consultazione di ogni account va in timeout dopo 3 secondi, un account che non risponde, o il cui token di accesso è scaduto, viene valutato in base ai numeri salvati, e una registrazione del limite prevale su di essi. Una finestra salvata il cui orario di azzeramento è passato conta come 0%, così un account la cui finestra si è nel frattempo azzerata non viene scartato. `--auto` non rinnova mai i token. A differenza della semplice rotazione, `--auto` può scegliere l'account usato per ultimo se è ancora quello con più margine.

### --continue

1. ccam cerca l'ultima conversazione di questa cartella in ogni account: tra i file `<session>.jsonl` in `projects/<folder>/` (il nome della cartella è la cartella di lavoro, con i link risolti e ogni carattere che non sia una lettera o una cifra trasformato in `-`, come la chiama Claude Code), quello con il prompt o la risposta più recente. Le sessioni che non ne hanno, come una sessione che ha solo eseguito `/mcp`, `/usage` o `/status`, o un file vuoto, vengono saltate. Se una sessione ha lasciato una nota (un limite di utilizzo o `/ccam:switch`, vedi sotto) e nessuna conversazione ha un prompt o una risposta più di 2 minuti dopo, viene usata proprio quella sessione; le righe che Claude Code scrive a `/exit` non contano.
2. L'account a cui appartiene è quello da cui passi a un altro, e il nuovo account viene scelto come sopra, dentro la stessa organizzazione. Quando più account condividono quella cronologia, appartiene a quello il cui `.claude.json` dice che ha eseguito la conversazione per ultimo (`lastSessionId`), altrimenti a quello usato per ultimo.
3. Se il nuovo account condivide quella cronologia (la stessa cartella di gruppo), ccam avvia `claude --resume <session id>`, proprio quella conversazione (`--continue` aprirebbe la più recente della cartella, qualunque sia).
4. Altrimenti avvia `claude --resume <the conversation's file> --fork-session`: Claude Code copia l'intera conversazione, con un nuovo ID di sessione, nella cronologia propria del nuovo account, e l'originale resta com'era. Da quel momento il `--continue` del nuovo account la trova.
5. A meno che tu non indichi `--agent` dopo `--`, ccam passa anche l'agente con cui gira la conversazione: l'impostazione `agent` del `.claude/settings.local.json` o del `.claude/settings.json` della cartella o del `settings.json` dell'account, altrimenti l'agente indicato dall'ultima riga `agent-setting` della conversazione. Claude Code non ripristina l'agente quando copia una conversazione (passo 4), e proseguirebbe come il suo agente predefinito, con tutti gli strumenti. Quando il nuovo account non condivide agenti, plugin e impostazioni di quello vecchio (le stesse cartelle `agents` e `plugins` e lo stesso `settings.json`, tramite un gruppo), ccam aggiunge una riga dopo quella sullo spostamento: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

ccam dice quanto è grande la conversazione: il nuovo account la rilegge tutta con la sua prima risposta. Le cache dei prompt non vengono mai condivise tra organizzazioni, e ogni piano personale è un'organizzazione a sé, quindi quando i due account appartengono a organizzazioni diverse ccam aggiunge `prompt caches are not shared between organizations, so the first reply reads it all again`. Quando il vecchio account ha effettuato l'accesso a server MCP a cui il nuovo non ha accesso, ccam li indica: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. Condividere `history` in un gruppo evita la copia, mantiene un'unica conversazione e condivide anche i checkpoint da cui `/rewind` ripristina il codice (`file-history/`).

### --brief

`--continue` riprende l'intera conversazione, e il nuovo account la rilegge tutta con la sua prima risposta. `ccam next --continue --brief` avvia invece una nuova conversazione, da un riassunto che ccam scrive della vecchia senza interpellare un modello:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

Il riassunto contiene, in quest'ordine e solo dove c'è qualcosa da dire:

1. il titolo della conversazione;
2. il percorso della sua trascrizione, in cui Claude può cercare quando gli serve un dettaglio;
3. la tua prima richiesta, quando ce n'erano più di quattro;
4. il riassunto che Claude Code ha scritto l'ultima volta che ha compattato la conversazione;
5. le tue ultime quattro richieste, con i comandi `!` e i comandi slash come li hai digitati;
6. l'ultima risposta di Claude;
7. la lista delle attività;
8. i 10 file che Claude ha modificato per ultimi;
9. il `git status --short --branch` della cartella, fino a 800 caratteri.

La lista delle attività tiene al massimo 10 voci: vengono tolte per prime le voci completate, poi i passi più lontani, così restano il passo corrente e i successivi. Le parti lunghe vengono tagliate a metà, e l'intero riassunto resta sotto i 16.000 caratteri.

ccam lo salva come `~/.local/state/ccam/brief/<session>.md`, leggibile solo da te, e il `--brief` successivo rimuove i riassunti più vecchi di 7 giorni. Claude Code parte con un solo prompt, e allega il file che menziona: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` Un prompt che indichi dopo `--` prende il posto dell'ultima frase (`Then: <your prompt>`), e le altre opzioni dopo `--` lo seguono.

Claude Code termina una menzione a un `"` e la tronca a un `#`, quindi quando il percorso del riassunto contiene uno dei due, ccam lo menziona con il percorso relativo alla cartella di lavoro; quando anche quel percorso ne contiene uno, il prompt nomina il file, e Claude lo legge con il suo strumento Read, che chiede prima conferma perché il file è fuori dalle cartelle di lavoro (in un'esecuzione `-p` in modalità Manual, Claude Code nega la lettura e Claude prosegue senza il riassunto).

`ccam --dry-run next --continue --brief` stampa il riassunto e non salva nulla. La vecchia conversazione resta com'era, e quella nuova sa solo ciò che dice il riassunto, quindi usa `--brief` quando conta più il compito che ogni dettaglio di come ci sei arrivato. Senza `--continue`, ccam si ferma con `--brief works with --continue`; quando non riesce a leggere la conversazione o non vi trova nulla da riassumere, si ferma e dice che `ccam next --continue` senza `--brief` la riprende per intero. Il [relay](#relay) e [`auto_switch`](#auto-switch) riprendono sempre l'intera conversazione. L'opzione `--brief` propria di Claude Code va dopo `--`, come le sue altre opzioni.

## Dentro Claude Code

Il [plugin di ccam](/docs/install#claude-code-plugin) aggiunge due hook (Claude Code 2.1.78 o successivo):

- Quando un turno si ferma a un limite di utilizzo, annota la conversazione e mostra una notifica desktop nei terminali che la supportano (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  Riconosce ogni messaggio che Claude Code considera un limite di utilizzo: i limiti delle 5 ore, settimanale, semplice e di modello (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` per la quota di un modello), i crediti di utilizzo e i limiti di spesa, e un'organizzazione che ha fermato l'account.

  Quando Claude Code ha scritto la richiesta rifiutata nella conversazione, l'hook registra anche la finestra piena fino al suo reset ([la registrazione del limite](/features/usage#the-limit-record)), tranne per un blocco impostato da un'organizzazione, che i suoi amministratori possono togliere in qualsiasi momento. Un breve rallentamento del server o il limite di frequenza di una chiave API non è un limite di utilizzo e viene ignorato, così come le sessioni non presidiate (esecuzioni `-p` e sessioni in background). Un subagente o un agente di workflow che si ferma a un limite di utilizzo conta come arresto della sua sessione, anche mentre la conversazione principale si limita ad aspettarlo: Claude Code scrive quella richiesta rifiutata nel file dell'agente stesso, nella cartella `subagents` della sessione, e l'hook la legge lì. Quando più agenti si fermano insieme, la sessione riceve una sola nota e una sola notifica.

  Claude Code 2.1.234 e successivi inoltre attendono nella sessione e proseguono da soli dopo l'azzeramento delle 5 ore; la notifica ti dà la possibilità di cambiare account invece di aspettare, e con [`auto_switch`](#auto-switch) attivo (Linux e macOS), ccam cambia account per te quando un altro account ha margine. Quando nessun altro account può subentrare, perché nessun altro account appartiene alla stessa organizzazione Team o Enterprise, o nessuno è un altro piano personale (con `cross_org` attivo, può farlo qualsiasi altro account Claude), la notifica termina invece con `You have no other personal plans, so ccam cannot move this conversation.` o `You have no other accounts in Acme, so ccam cannot move this conversation.`
- `/ccam:switch [account]` annota la conversazione e chiude il turno prima che venga fatta qualsiasi richiesta: digita `/exit`, poi esegui `ccam next --continue` per continuarla sull'account che hai indicato (o su quello successivo). In una sessione che non ha ancora una conversazione, dice che non c'è nulla da continuare. Lo esegue solo ciò che digiti tu: dalla versione 0.3.2 del plugin, Claude non può eseguire `/ccam:switch`, e un job che Claude pianifica (CronCreate, ScheduleWakeup o `.claude/scheduled_tasks.json`) lo passa a Claude come testo semplice.

  Quando l'account che indichi è fuori dall'organizzazione della conversazione (un posto di un'altra organizzazione Team o Enterprise, un piano personale invece di un posto o viceversa, oppure un account con chiave API da una delle due parti), `ccam next --continue` chiede prima, qualunque cosa dica `cross_org`: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. Chiede lo stesso prima di un account che sceglie da solo quando non ne indichi nessuno, e che `cross_org = true` gli permette di prendere da un'altra organizzazione: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. Solo `y` o `yes` la sposta; su Linux e macOS, i tasti premuti prima della domanda vengono prima scartati. Qualsiasi altra risposta, o l'assenza di un terminale su cui chiedere, lascia la conversazione e la nota dove sono: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (dopo una scelta, `ccam next --continue` sceglie di nuovo). Come per le scelte automatiche, un account di cui ccam non riesce a sapere l'organizzazione va bene con qualsiasi account.

La nota resta in `~/.local/state/ccam/handoff/<account>/<session>.json`, un file per sessione, così due sessioni dello stesso account non sovrascrivono mai l'una la nota dell'altra. Viene conservata per 12 ore e usata una sola volta, quando la sessione successiva sta per partire; un avvio che fallisce la lascia per il tentativo successivo. Appartiene alla cartella in cui è partita la sessione, anche dopo un `cd` nella shell di Claude Code. `ccam rename` aggiorna i nomi al suo interno; se nel frattempo l'account che indica è stato rimosso, ccam lo segnala e sceglie come al solito.

### Relay

Con `relay = true` in cima a `config.toml`, ccam resta in esecuzione accanto alle sessioni interattive invece di sostituirsi con Claude Code, e dopo `/ccam:switch` continua da solo la conversazione di quella sessione (una nota lasciata da un'altra sessione non viene presa in carico): l'hook chiude Claude Code, e ccam esegue `ccam next --continue` con la nota di quella sessione. Senza un account indicato, prende quello con più margine, tenendo conto del modello della conversazione, e quando `cross_org = true` permette che sia un account di un'altra organizzazione, pone la domanda qui sopra prima di spostarvi la conversazione. Su Windows l'hook non può chiudere Claude Code: digita `/exit`, e ccam la continua. Né l'hook chiude Claude Code per un account fuori dall'organizzazione della conversazione: dice `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`, e ccam pone la domanda qui sopra una volta che hai digitato `/exit`. Questo costa un processo ccam in attesa per ogni sessione, ed è per questo che è disattivato per impostazione predefinita.

La sessione che prosegue mantiene ciò che la ripresa di una conversazione non ripristina: opzioni come `--add-dir`, `--mcp-config`, gli elenchi degli strumenti, `--settings`, `--plugin-dir` e le opzioni del prompt di sistema, la modalità dei permessi della conversazione, il livello di effort della sessione e l'agente con cui gira la conversazione (vedi [`--continue`](#continue)). Claude Code stesso ripristina `--model` e `--name` dalla conversazione. Tutto il resto (`--worktree`, `--debug`, `--remote-control`, opzioni che ccam non conosce) viene indicato in `ccam: not carried over: …`. Il cambio non pone mai la domanda sulla cronologia condivisa, e ciò che hai digitato mentre era in corso viene scartato; la sua unica domanda è quella prima di un account di un'altra organizzazione.

Il relay continua la conversazione solo quando Claude Code è terminato come `/ccam:switch` gli ha chiesto: non quando hai digitato prima `/exit`, hai premuto Ctrl+C mentre terminava, o ccam ha ricevuto un hangup o un altro segnale, e non quando il terminale non c'è più. La nota allora resta, e `ccam next --continue` la continua più tardi. Una sessione avviata sotto il relay di una versione più vecchia di ccam, o una non presidiata (esecuzioni `-p` e sessioni in background), non viene chiusa dall'hook; riceve invece il messaggio con `/exit`. `claude remote-control` viene eseguito senza relay.

### `auto_switch`

Con `auto_switch = true` in cima a `config.toml` (Linux e macOS), ccam sposta da solo una sessione su un altro account quando Claude Code la ferma a un limite di utilizzo: l'hook del plugin sceglie un altro account con margine e chiude Claude Code, e ccam, in esecuzione accanto alla sessione, continua lì la conversazione, come fa il [relay](#relay) dopo `/ccam:switch`. È disattivato per impostazione predefinita ed è sperimentale.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam sposta una sessione solo quando valgono tutte queste condizioni:

1. ccam l'ha avviata in un terminale e resta in esecuzione accanto a essa. `auto_switch` attiva il [relay](#relay), quindi questo vale per le sessioni interattive che avvii dopo averlo impostato. Una sessione avviata prima, o sotto un ccam più vecchio, riceve una notifica che termina con ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` Le sessioni non presidiate (esecuzioni `-p` e sessioni in background), `claude remote-control` e Claude Code avviato senza ccam non vengono mai spostati.
2. L'account ha il [plugin di ccam](/docs/install#claude-code-plugin) attivato ed esegue Claude Code 2.1.214 o successivo, che passa all'hook l'ID del proprio processo (`CLAUDE_PID`).
3. Claude Code ha scritto nella conversazione la richiesta rifiutata, con l'ora del reset, per il limite delle 5 ore, settimanale, di Opus o di Sonnet di quell'account. Un limite di spesa, l'esaurimento dei crediti di utilizzo, un blocco impostato da un'organizzazione, la quota di Fable (`You've reached your Fable limit.`), un rallentamento del server e il limite di frequenza di una chiave API non spostano mai una sessione. Il limite di un subagente o di un agente di workflow conta, anche mentre la conversazione principale si limita ad aspettarlo; quando si tratta di un limite settimanale di Opus o di Sonnet, la sessione si sposta solo se la conversazione stessa usa quel modello.
4. ccam conosce il piano dell'account, dal suo `.claude.json` o da ciò che ha registrato `ccam login`. Altrimenti la notifica termina con `ccam does not know personal's plan, so auto_switch leaves this session.`
5. La sessione non è collegata a Remote Control. Con un altro account non si ricollegherebbe all'app da cui la controlli, quindi resta dov'è, e la notifica lo dice.
6. Il limite si azzera tra più di 15 minuti. Per un'attesa più breve ccam lascia la sessione a Claude Code: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. Non hai digitato un prompt dopo il limite.
8. Un altro account ha margine. ccam sceglie come fa [`--auto`](#auto) (`auto_max_usage`, `auto_strategy`, il modello della conversazione), tra gli account dello stesso tipo di quello che lasci: un altro posto della stessa organizzazione Team o Enterprise, o un altro piano personale (Pro, Max). `cross_org` e `--cross-org` qui non valgono. L'account deve:
   - avere un login che Claude Code può rinnovare;
   - non essere collegato allo stesso account Claude;
   - non usare una chiave API;
   - avere un piano noto a ccam;
   - non avere i crediti di utilizzo attivi, per quanto risulta dalla consultazione dell'utilizzo di ccam stesso;
   - aver dato fiducia a questa cartella in Claude Code;
   - se la conversazione gira con un agente (l'impostazione `agent` o la conversazione ne indica uno), condividere gli agenti dell'account che lasci: le sue cartelle `agents` e `plugins` e il suo `settings.json` devono essere ciascuno esattamente gli stessi (tramite un gruppo), oppure mancare in entrambi.

   Quando queste ultime tre condizioni escludono tutti gli account con margine, la notifica li indica: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` oppure `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` Una finestra con una [registrazione del limite](/features/usage#the-limit-record) conta come piena, e un account verso cui un'altra sessione si è spostata negli ultimi 2 minuti viene dopo gli altri. La scelta deve avvenire entro circa 6 secondi; quando le consultazioni dell'utilizzo richiedono di più, la sessione resta dov'è. Con `usage_source = "claude"`, `auto_switch` decide solo in base ai numeri salvati, e ccam non può sapere né dalla consultazione di Claude Code né dalla riga di stato se i crediti di utilizzo sono attivi, quindi li considera disattivati.

Poi l'hook chiude Claude Code, e ccam, in esecuzione accanto a esso, mostra una notifica desktop nei terminali che la supportano, come `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` Stampa la prima riga dell'output di esempio qui sopra ed esegue `ccam next --continue` per quella sessione con l'account scelto, che stampa le sue altre due righe e avvia Claude Code con le opzioni, la modalità dei permessi, il livello di effort e l'agente della sessione, come fa il [relay](#relay). Il primo messaggio della nuova sessione è un testo fisso, l'unica cosa che ccam invia per tuo conto (uno spostamento dopo `/ccam:switch` non invia nulla):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

Dopo il limite di un subagente è quest'altro, perché i subagenti, i workflow e i comandi in background si sono fermati con la vecchia sessione:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

Se anche quell'account raggiunge il suo limite, lì succede la stessa cosa. Le registrazioni dei limiti impediscono alla catena di tornare a un account prima del suo reset, quindi quando tutti gli account sono pieni la sessione resta sull'ultimo.

L'hook sa che una conversazione gira con un agente dall'impostazione `agent`, o dalla riga che indica l'agente e che Claude Code scrive nella conversazione di tanto in tanto e quando la sessione termina. Se finora solo il `--agent` della sessione lo indica, e l'account scelto non condivide gli agenti, `ccam next --continue` non avvia nulla e conserva la nota, e una seconda notifica desktop lo dice: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

Quando nessun altro account ha margine, non si sposta nulla: Claude Code gestisce il limite come al solito, attendendo nella sessione od offrendo le sue opzioni, e la notifica dice quale account si libera per primo, per esempio `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` Per i posti dice `None of your other accounts in Acme has room now`. Per fermare uno spostamento già iniziato, premi Ctrl+C mentre Claude Code si chiude: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

Uno spostamento ha dei costi. Claude Code si riavvia, il che richiede qualche secondo, e le shell in background, i server di sviluppo e gli agenti avviati da Claude si fermano con la vecchia sessione. Il nuovo account rilegge l'intera conversazione con la sua prima risposta, come con [`--continue`](#continue), e ccam dice quanto è grande; uno spostamento non usa mai [`--brief`](#brief).

`auto_switch` non cambia mai account prima che Claude Code segnali un limite, non controlla mai il tuo utilizzo in background, e non copia, rinnova o passa mai un login. Su Windows viene ignorato, e `ccam doctor` lo dice; lì, con `relay = true`, `/ccam:switch` e `/exit` spostano comunque una conversazione su un altro account.

**I termini di Anthropic.** Sui posti Team ed Enterprise, `auto_switch` sposta il tuo lavoro su un altro posto della stessa organizzazione. Il centro assistenza di Anthropic dice: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (ogni posto di un piano Team è assegnato nominativamente a una sola persona, e solo quella persona può usarlo) e "Users can't route their work through another user's account to use that user's remaining usage" (gli utenti non possono far passare il proprio lavoro attraverso l'account di un altro utente per sfruttare l'utilizzo residuo di quell'utente) ([piani Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); lo stesso vale per i [piani Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)).

Per i piani personali (Pro, Max), i termini di Anthropic non dicono se i tuoi account possano subentrare automaticamente l'uno all'altro a un limite. I [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) vietano di accedere ai Servizi "through automated or non-human means, whether through a bot, script, or otherwise" (con mezzi automatizzati o non umani, tramite un bot, uno script o altro) salvo esplicita autorizzazione di Anthropic, e uno spostamento invia per tuo conto uno dei messaggi fissi qui sopra. La [Usage Policy](https://www.anthropic.com/legal/aup) che entra in vigore il 12 novembre 2026 elenca "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (aggirare o interferire intenzionalmente con funzionalità, restrizioni o protezioni stabilite nei nostri prodotti) tra le cose che vieta.

ccam non afferma che `auto_switch` rispetti questi termini. Attivarlo è una tua decisione, e lo è anche il rischio: Anthropic può sospendere un account, o l'accesso di un'organizzazione Team o Enterprise, quando ritiene che i termini siano stati violati.

Mentre `auto_switch` è attivo, `ccam doctor` avvisa per un account senza il plugin di ccam e per le versioni di Claude Code più vecchie della 2.1.214, e per i posti ripete la regola di Anthropic sui posti ([controlli](/features/doctor)).

## Rieseguire al limite

`ccam next` serve per una sessione interattiva che si è fermata. Un'esecuzione non interattiva, `claude -p` in uno script o in un job pianificato, non può aspettarti, quindi `--fallback` prima del nome dell'account permette a ccam di cambiare account al posto suo:

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

1. `--fallback` funziona solo insieme a `-p` o `--print`; altrimenti ccam si ferma e rimanda a `ccam next --auto --continue`.
2. Lo standard input da una pipe viene letto una volta sola, fino a 10 MiB, e passato a ogni tentativo. Come Claude Code, ccam attende al massimo 3 secondi che inizi ad arrivare (tranne con `--input-format stream-json`); se una pipe resta aperta ma non invia nulla, l'esecuzione non riceve input e ccam lo segnala sullo standard error.
3. ccam avvia `claude` come processo figlio invece di usare `exec`, trattiene il suo standard output e lascia passare lo standard error man mano che arriva.
4. Quando `claude` esce con 0, o fallisce per un altro motivo, l'output trattenuto e il codice di uscita vengono passati invariati. Lo stesso vale per un'esecuzione terminata da un segnale (un codice di uscita sopra 128): le è stato chiesto di fermarsi, e non viene rieseguita.
5. Un errore conta come limite di utilizzo quando il messaggio di Claude Code lo dice, con uno qualsiasi dei messaggi di limite di utilizzo che l'[hook](#inside-claude-code) riconosce, anche dentro l'output `--output-format json` o `stream-json`, oppure quando l'utilizzo dell'account mostra la finestra di 5 ore, quella settimanale o quella di `--model` al 100% (compresa una registrazione del limite).
6. A quel punto l'output di quel tentativo viene scartato e lo stesso comando viene eseguito con un altro account connesso: quello con l'utilizzo su 5 ore più basso, poi gli account il cui utilizzo non si può consultare, in ordine di numero. Gli account al 100% vengono saltati, gli account con un login a chiave API non vengono mai provati, e ogni account viene provato al massimo una volta.
7. Con `-c`, `--continue`, `-r` o `--resume` tra gli argomenti vengono provati solo gli account dello stesso gruppo di cronologia, così la conversazione può essere trovata.
8. Vengono provati solo gli account della stessa organizzazione, come con `--auto` (`cross_org = true` toglie questa regola).

L'account successivo viene scelto con le stesse regole di `--auto` (`auto_max_usage`, `auto_strategy`, la finestra di `--model`). Quando non resta alcun account, ccam stampa l'output dell'ultimo tentativo, dice che nessun account connesso ha ancora margine ed esce con codice 75. Con `--wait` prima del nome dell'account (`ccam --fallback --wait=2h work -p …`), attende invece l'azzeramento del primo account, poi può provare di nuovo ogni account; attende sempre almeno una volta prima, così un account che ha appena raggiunto il limite non viene rieseguito subito. La durata limita il tempo di attesa totale dell'esecuzione, non ciascuna attesa. Scritta con uno spazio (`--wait 30m`), viene rifiutata con lo stesso suggerimento di `ccam next`. Un `--` iniziale prima degli argomenti di Claude Code funziona come per qualsiasi avvio (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` stampa il comando del primo tentativo.

Le sessioni interattive non cambiano account da sole, a meno che tu non attivi [`auto_switch`](#auto-switch). Con `ccam statusline install`, la riga di stato aggiunge `→ ccam next --auto --continue` quando l'utilizzo su 5 ore o settimanale arriva al 90% ([Limiti di utilizzo](/features/usage)).

## Sicurezza e casi limite

- Gli account senza login non vengono mai scelti. Se non c'è un altro account con login, `ccam next` da solo si ferma con `no logged-in account other than <name>`.
- Con `--continue` e nessuna conversazione nella cartella: `no account has a conversation in <folder> to continue`.
- Quando la regola dell'organizzazione non lascia alcun account: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Con `--auto` e nessun account con login: `no logged-in account to pick`. Quando ogni account è pieno o ha un utilizzo sconosciuto: `` no account with known usage and headroom (see `ccam usage`) ``. In entrambi i casi il codice di uscita è 75.
- `--max-usage`, `--strategy`, `--wait` e `--cross-org` funzionano solo con `--auto`, e un valore fuori da 1-100 o una strategia sconosciuta vengono rifiutati.
- Un fork copia la conversazione, non il resto dello stato dell'account: i checkpoint nel `file-history/` dell'account originale non vengono con essa, a meno che la cronologia non sia condivisa.
- `--auto` effettua richieste di rete anche con `--dry-run`. Si limita a leggere e non scrive mai i file dei token.
- L'account scelto si avvia tramite lo stesso percorso di `ccam <name>`, quindi valgono gli stessi controlli, compreso l'avviso quando in questa cartella è già in esecuzione un'altra sessione di Claude Code che condivide la stessa cronologia.
- `last-used` viene scritto solo quando Claude Code si avvia davvero. `--dry-run` non lo tocca, quindi un dry run non fa avanzare la rotazione.
- `--fallback` trattiene lo standard output finché un tentativo non termina, quindi un output in streaming come `--output-format stream-json` arriva tutto insieme alla fine. Un client che invia input `stream-json` in streaming e attende le risposte prima di inviarne altro non può usare `--fallback`. I file che il primo tentativo ha modificato prima di raggiungere il limite restano modificati; il tentativo successivo riparte dall'inizio del prompt.

## Pagine correlate

- [Limiti di utilizzo](/features/usage)
- [Limiti di utilizzo e cambio di account](/docs/usage-limits)
- [Condividere server MCP, impostazioni e cronologia](/docs/sharing)
- [Cosa fare al limite di utilizzo di Claude Code](/guides/usage-limit-workflow)
