---
title: Limiti di utilizzo
description: ccam usage mostra insieme i limiti su 5 ore e settimanali di tutti gli account, e ccam statusline mette il nome dell'account e l'utilizzo in Claude Code.
updated: 2026-11-01
order: 7
translation: machine
sourceHash: 54b296f7da35
---

`ccam usage` mostra, per ogni account, quanta parte di ciascun limite è stata usata: il limite su 5 ore, il limite settimanale e il limite settimanale proprio di ogni modello (come Fable). Ogni limite ha una barra e il tempo mancante al suo azzeramento. Usalo prima di iniziare un'attività lunga, o quando Claude Code ti avvisa che un limite è vicino, per vedere quale account ha ancora margine. `ccam statusline` è il comando per la riga di stato di Claude Code: mostra il nome dell'account e l'utilizzo dentro Claude Code e salva i numeri per `ccam usage`. Per i token e i costi stimati che stanno dietro all'utilizzo, `ccam ccusage` esegue [ccusage](https://ccusage.com) sulle conversazioni di ogni account ([sotto](#tokens-and-costs-ccam-ccusage)).

## Utilizzo

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts that are not running
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

Ogni account ha una riga con il suo numero, il suo nome e quanto sono vecchi i suoi numeri; `(cached)` significa che provengono da un file salvato e non dall'endpoint. Sotto, una riga per limite mostra una barra di 20 celle, la percentuale e il tempo mancante all'azzeramento della finestra: `5-hour`, `Weekly` (tutti i modelli), poi la finestra settimanale propria di ogni modello, come `Fable`. Una finestra che l'account non ha mostra `-`. Una percentuale diventa gialla all'80% e rossa al 100%.

Se Claude Code ha rifiutato una richiesta perché una finestra era piena, quella finestra mostra 100% fino al suo reset, qualunque cosa dicano i numeri salvati (vedi [la registrazione del limite](#the-limit-record)). In un terminale ogni account compare subito come `fetching…`, e le sue righe vengono disegnate quando risponde.

L'ultima riga, mostrata solo quando hai più di un account, indica l'account che [`ccam next --auto`](/features/next) sceglierebbe. Segue le stesse regole:

- `auto_max_usage` e `auto_strategy`;
- le sessioni in esecuzione;
- solo gli account dell'organizzazione dell'ultimo account usato, a meno che non sia impostato `cross_org`;
- ogni finestra al 100% viene saltata, così un account con la settimana esaurita non viene mai indicato.

Con la strategia predefinita `headroom` riporta `→ Most headroom: <name>`; con `expiring` o `balance` indica il motivo, come `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Posti nei piani Team ed Enterprise.** La scelta automatica dell'account (`ccam next`, `--fallback`, `ccam pick`, il tasto `n` della dashboard e, una volta attivato, [`auto_switch`](/features/next#auto-switch)) passa anche da un posto all'altro della stessa organizzazione Claude Team o Enterprise, e questo può violare i termini di Anthropic. Leggi le [domande frequenti](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e verifica con la tua organizzazione e con Anthropic prima di farci affidamento.

Un account senza login mostra `not logged in → ccam login <name>`, e un account con chiave API mostra `no limits (API)`. Un token di accesso scaduto mostra ``access token expired; `ccam <name>` renews it when it starts``, sopra gli eventuali numeri salvati: Claude Code lo rinnova al prossimo avvio. Un login scaduto del tutto, perché il suo refresh token è scaduto o perché Claude Code lo ha cancellato dopo che il server lo ha rifiutato, mostra `` login expired → `ccam login <name>` ``.

Dopo `ccam statusline install`, Claude Code mostra una riga come `[work] · Opus · 5h 24% · 7d 41%`.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--json` | Stampa i risultati in JSON: `name`, `slot`, `status`, `limits` (con `credits` quando l'account ha i crediti di utilizzo attivi), `source`, `at`, `error`, e `limited` e `limited_windows` quando una registrazione del limite porta delle finestre al 100%. Senza account stampa `[]` | disattivata |
| `--watch[=N]` | Ridisegna ogni `N` secondi finché non premi Ctrl+C. `--watch` da solo, o un valore pari o inferiore a 0, indica 60. Il valore va dopo `=`: `--watch 10` viene rifiutato con un suggerimento | disattivata |
| `-v`, `--verbose` | Aggiunge sotto ogni account una riga `source:`: da dove vengono i suoi numeri (l'endpoint di utilizzo, Claude Code, la riga di stato, un file salvato, o il limite di utilizzo segnalato da Claude Code). Quando una registrazione del limite prevale sui numeri, segue una riga come `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` | disattivata |
| `--no-cache` | Ignora la cache di 60 secondi e interroga di nuovo l'endpoint | disattivata |
| `--source live\|claude\|cache` | `live` interroga l'endpoint con il token dell'account; `claude` fa consultare l'utilizzo a Claude Code (vedi sotto); `cache` mostra solo i numeri salvati | `usage_source`, altrimenti `live` |
| `--refresh` | Rinnova i token di accesso scaduti, solo per gli account che non sono in esecuzione. Un rinnovo può richiedere fino a 30 secondi | disattivata |
| `--dry-run` | Accettata: `ccam usage` non modifica nulla tranne la propria cache. Rifiutata insieme a `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | disattivata |

Qualsiasi altro valore di `--source` fallisce con `--source must be live, claude or cache`. `ccam usage` non accetta argomenti.

### L'utilizzo letto da Claude Code stesso

Con `usage_source = "claude"` in cima a `config.toml` (sopra il primo `[[group]]`, `[[profile]]` o `[[folder]]`), oppure `--source claude`, ccam non legge affatto il token di login. Per ogni account avvia `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` nella cartella dell'account e invia sullo standard input la richiesta `get_usage` di Claude Code (Claude Code 2.1.270 o successivo). Claude Code consulta l'utilizzo con il proprio login e risponde con le finestre su 5 ore, settimanale e per modello.

`--safe-mode` impedisce l'avvio di hook, server MCP e plugin dell'account, e nessuna sessione viene salvata. Ogni account richiede un avvio di Claude Code (circa un secondo), quindi questa fonte è più lenta di `live`. Gli account vengono interrogati tutti insieme, ciascuno con un timeout di 20 secondi, e la cache di 60 secondi vale come al solito.

Un login il cui token di accesso è scaduto, o scade entro 5 minuti, ha 50 secondi, perché Claude Code lo rinnova all'avvio e un rinnovo interrotto può far perdere il login. Se ccam abbandona una consultazione, chiede a Claude Code di terminare con SIGTERM, così Claude Code completa un rinnovo che ha iniziato (Windows non ha SIGTERM, quindi lì il processo viene terminato).

Un login scaduto del tutto mostra `` login expired → `ccam login <name>` `` senza avviare Claude Code. Claude Code indica la richiesta come sperimentale, quindi se smette di rispondere, ccam mostra i numeri salvati e `ccam doctor` lo segnala. Quando Claude Code risponde che la sua consultazione è fallita, ccam mostra i numeri salvati con `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`, non un piano senza limiti.

Con `--dry-run`, ccam non avvia Claude Code, perché Claude Code scrive nella cartella dell'account; mostra invece i numeri salvati. `usage_source` accetta `live` o `claude`; `cache` è solo una scelta di `--source`.

`ccam statusline` legge dallo standard input il JSON della riga di stato di Claude Code. La sua unica opzione è `--wrap <command>`: un comando esistente per la riga di stato, davanti alla cui prima riga di output viene messo il nome dell'account. `ccam statusline install` non accetta opzioni. Entrambi rifiutano `--dry-run` con `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Token e costi: `ccam ccusage`

`ccam usage` mostra quanta parte di ogni limite ha usato un account. [ccusage](https://ccusage.com) invece conta i token e i costi stimati nei log delle conversazioni di Claude Code. Da solo cerca i log di Claude Code soltanto in `~/.claude` e `~/.config/claude` (o nelle cartelle indicate da `CLAUDE_CONFIG_DIR`), quindi l'unico account di ccam che può vedere è uno importato sul posto in `~/.claude`. `ccam ccusage` lo esegue con `CLAUDE_CONFIG_DIR` impostato sulle cartelle di ogni account che ha conversazioni, così un solo report li copre tutti:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam legge prima le proprie opzioni: `-a` (`--account`) accetta nomi o numeri di account, separati da virgole o ripetendo l'opzione, e `--dry-run` stampa il comando invece di eseguirlo. Dal primo argomento che non è un'opzione di ccam in poi, tutto passa a ccusage così come l'hai digitato, tranne `--config` per `bun x` e `npx` (vedi sotto); `ccam ccusage -- --help` mostra l'aiuto di ccusage stesso. ccam esegue `ccusage` se è installato, altrimenti `bun x ccusage`, e in mancanza di questo `npx ccusage@latest`; questi ultimi due lo scaricano da npm. Si avviano nella cartella di ccam stesso (`~/.local/state/ccam/ccusage`), non nella cartella in cui sei, quindi un progetto che sta lì non può scegliere cosa eseguono con il suo `node_modules/.bin`, un `bin` nel suo `package.json` o un registry nel suo `.npmrc`; `ccam --dry-run ccusage` mostra il `cd` in quella cartella. Quando nessuno di questi è disponibile, ti dice di installare Node.js, Bun o ccusage. Su Linux e macOS ccusage prende il posto di ccam; su Windows ccam lo attende e termina con il suo codice di uscita. Sotto WSL, ccam salta un `ccusage`, `bun` o `npx` che si trova su un'unità Windows, perché verrebbe eseguito come programma Windows senza le cartelle degli account.

- Il report copre gli account le cui cartelle hanno conversazioni (una cartella `projects`). Una conversazione che più account condividono tramite un gruppo viene contata una sola volta.
- Con `-a`, un account che condivide la propria cronologia tramite un gruppo porta con sé le conversazioni di tutto il gruppo, e ccam dice di chi sono: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- `blocks` mostra le finestre di 5 ore, e ogni account ha le sue, quindi indica un solo account quando lo usi, come sopra.
- ccusage riporta anche gli altri strumenti di programmazione che trova, come Codex, a meno che tu non ne indichi uno: `ccam ccusage claude daily` è solo Claude Code.
- Quando è ccam a eseguirlo, ccusage legge le sue impostazioni (`ccusage.json`) solo da `./.ccusage` e dalle cartelle degli account. `--config` indica un altro file: `ccam ccusage daily --config ~/.config/claude/ccusage.json`. Poiché `bun x` e `npx` si avviano nella cartella di ccam stesso, ccam passa loro `./.ccusage/ccusage.json` con `--config` quando quel file esiste e non ne indichi nessuno, e passa un `--config` relativo come percorso completo. ccusage legge poi gli altri percorsi relativi, come `--pi-path`, dalla cartella di ccam, quindi indicali per intero.
- ccusage leggerebbe come due cartelle una cartella il cui percorso contiene una virgola, quindi ccam la rifiuta e indica l'account; lascia fuori quell'account con `-a`.
- `ccusage` è un comando di ccam, quindi un account chiamato ccusage si avvia con `ccam run ccusage`; `ccam doctor` lo segnala.

## Come funziona

I numeri di `ccam usage` provengono da due fonti: l'endpoint di utilizzo non documentato di Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, chiamato con il login di ciascun account, e i valori che Claude Code passa al suo comando per la riga di stato. Tutto si trova in `internal/usage`.

1. ccam carica `config.toml` e costruisce un `usage.Account` per ogni account (`cli.usageCmd`): la sua cartella `~/.local/share/ccam/profiles/<name>/`, i suoi indicatori di sessione (insieme ai record propri di Claude Code nella cartella dell'account, dicono se è in esecuzione) e se è stato importato sul posto.
2. ccam consulta tutti gli account contemporaneamente, con una goroutine per account (`Fetcher.FetchAll`). Tutte condividono un solo `http.Client`, così le connessioni all'endpoint vengono riutilizzate. La connessione può richiedere 2 s, e ogni account ha in totale 3 s. Il client non segue i reindirizzamenti, perché un reindirizzamento invierebbe di nuovo il token altrove, e legge al massimo 1 MB della risposta.
3. Quando ogni account termina, la sua riga viene ridisegnata. Un account che fallisce o non risponde in tempo ripiega sui numeri salvati, così un account lento non blocca mai gli altri.

Nel diagramma, ogni account viene consultato nella propria goroutine, e quello che supera il tempo limite mostra invece i numeri salvati:

<!-- motion:usage-fetch -->

`personal` e `work` rispondono in tempo e mostrano numeri aggiornati. `side` impiega più di 3 s, quindi ccam mostra invece i suoi numeri salvati (`fallback()`). ccam disegna poi un blocco per account (`usage.Table()`) e sceglie l'account indicato nell'ultima riga (`usage.Best()`).

Per un singolo account, ccam esegue questi controlli in ordine (`Fetcher.Fetch`):

<!-- motion:usage-cache -->

1. Con `--source cache`, passa direttamente a `fallback()`.
2. A meno che non sia impostato `--no-cache`, un risultato aggiornato in `~/.local/state/ccam/usage/<name>.json` misurato meno di 60 secondi prima (`CacheTTL`) viene usato, valutato rispetto a questo momento: una finestra che si è azzerata nel frattempo mostra 0%, e un file datato avanti rispetto all'orologio non conta come recente. Eseguire `ccam usage` due volte di seguito costa una richiesta per account.
3. ccam legge `.credentials.json` (o il Portachiavi di macOS) senza modificarlo e decodifica solo la chiave `claudeAiOauth` (`readCredentials()`). Nessun login con abbonamento significa `not logged in`, oppure `no limits (API)` quando `.claude.json` contiene una chiave API. Un login il cui refresh token è scaduto, o che Claude Code ha cancellato, è `` login expired → `ccam login <name>` `` e non viene consultato.
4. Un token di accesso scaduto non viene usato. Senza `--refresh`, l'account mostra i numeri salvati con ``access token expired; `ccam <name>` renews it when it starts``.
5. Altrimenti ccam chiama l'endpoint. Un 401 o un 403 conta come login scaduto (`login expired: authentication failed (401)`), qualsiasi altro stato o un corpo che non riesce a leggere conta come errore, e in entrambi i casi si ripiega sui numeri salvati. Una consultazione terminata con Ctrl+C dice `cancelled`. Un 429 significa che l'endpoint ha ricevuto troppe richieste per questo account: ccam scrive `<name>.pause` con l'ora in cui chiedere di nuovo (`Retry-After`, oppure 5 minuti, al massimo un'ora), mostra i numeri salvati con `rate limited by the usage endpoint; asking again after 15:10` e non chiama l'endpoint per quell'account fino ad allora, nemmeno con `--no-cache`. Riprovare subito non farebbe che prolungare il limite.
6. Una risposta valida viene scritta in `<name>.json` con l'ora in cui è stata misurata.

Per i numeri salvati, ccam prende la consultazione salvata più recente (`<name>.json`, oppure `<name>.claude.json` dalle consultazioni di Claude Code) e, quando `<name>.statusline.json` è stato scritto dopo, vi sovrappone le finestre su 5 ore e settimanale della riga di stato (`fallback()`); le finestre per modello, che la riga di stato non ha, restano come le ha salvate la consultazione. Prima di mostrare i numeri salvati, controlla l'orario di azzeramento di ogni finestra: una finestra che si è azzerata dopo il salvataggio dei numeri mostra 0%, perché la vecchia percentuale non vale più. Per questo `personal` qui sopra mostra `0%` senza orario di azzeramento.

L'endpoint ha due formati di risposta. Quello più recente ha un array `limits` con le voci `session` (5 ore), `weekly_all` e `weekly_scoped` (per modello). Quello più vecchio ha gli oggetti `five_hour`, `seven_day`, `seven_day_opus` e `seven_day_sonnet`. ccam li legge entrambi (`ParseLive`): prevale l'array `limits`, e i vecchi oggetti completano ciò che non copre. Un corpo senza alcun limite viene segnalato come cambio di formato, di cui ti avvisa `ccam doctor`.

### La registrazione del limite

Quando Claude Code ferma un turno a un limite di utilizzo, l'hook del [plugin di ccam](/docs/install#claude-code-plugin) legge dalla conversazione la richiesta rifiutata: quale finestra era piena (5 ore, settimanale, Opus o Sonnet) e quando si azzera. Lo scrive in `~/.local/state/ccam/usage/<name>.limit.json`. ccam applica la registrazione sopra a tutto ciò che sa, da ogni fonte (`Fetch`): ogni finestra che contiene conta come 100% fino al suo reset, per quanto vecchi o bassi siano i numeri salvati, così `ccam next --auto`, `ccam pick`, `--fallback` e [`auto_switch`](/features/next#auto-switch) non scelgono di nuovo l'account prima di allora.

Se non si sa nient'altro, l'account mostra solo la registrazione (`source: the usage limit Claude Code reported` con `-v`). Le sessioni non presidiate (esecuzioni `-p` e sessioni in background) non scrivono alcuna registrazione, e nemmeno un blocco impostato da un'organizzazione (per esempio quando ha esaurito l'utilizzo o il limite di un gruppo è $0), che i suoi amministratori possono togliere in qualsiasi momento. `ccam login` in un altro account Claude elimina la registrazione; anche eliminare il file la cancella.

### Rinnovare i token

I refresh token ruotano: rinnovarne uno ti dà un nuovo refresh token e invalida quello vecchio. Se ccam rinnovasse il login mentre Claude Code è in esecuzione con lo stesso account, quella sessione verrebbe disconnessa. Per questo `ccam usage` legge i token e non li scrive mai, a meno che tu non passi `--refresh`. In quel caso ccam (`Fetcher.refresh`):

1. rifiuta gli account registrati sul posto con `ccam import`, dato che `claude` da solo potrebbe usare quel login a insaputa di ccam;
2. rifiuta i login salvati nel Portachiavi di macOS;
3. prende un file di lock, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, così due processi di ccam non rinnovano mai entrambi lo stesso token; finché è tenuto, l'account mostra `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)`. Anche `ccam rename`, `remove` e `logout` lo prendono, e si rifiutano mentre è in corso un rinnovo;
4. con il lock acquisito, si rifiuta se è in esecuzione una sessione dell'account, comunque sia stata avviata (tramite ccam, `claude` da solo, un IDE, oppure una sessione in background e il daemon dietro di essa, che rinnova da sé il login), e salta la richiesta se un altro processo ha già rinnovato il token;
5. non invia nulla mentre Claude Code tiene i propri lock di rinnovo nella cartella dell'account (`.oauth_refresh.lock`, `<folder>.lock`): in quel momento Claude Code sta rinnovando quel login;
6. invia il refresh token all'endpoint dei token e attende la risposta fino a 30 secondi, il tempo che Claude Code dà ai propri rinnovi. Una volta inviata, una richiesta non viene mai interrotta, nemmeno con Ctrl+C: il server potrebbe aver già ruotato il refresh token, e solo la sua risposta contiene il nuovo login. Ctrl+C o un hangup terminano le altre consultazioni, e ccam esce (130 dopo Ctrl+C) una volta salvata la risposta;
7. salva sotto il lock di scrittura di Claude Code (`.storage-write.lock`) e modifica solo i campi del token di `.credentials.json` così com'è in quel momento: ciò che Claude Code ha scritto nel frattempo (un accesso MCP) resta, un login che Claude Code ha rinnovato nel frattempo viene mantenuto, e un logout non viene mai annullato. Quando la risposta indica la durata del nuovo refresh token, viene salvata come `refreshTokenExpiresAt`, come fa Claude Code.

Un rinnovo che non va a buon fine dice perché, sulla riga dell'account o sopra i suoi numeri salvati: per esempio `token not renewed: a session is running`, oppure `login may be lost: …` quando una richiesta è stata inviata ma la sua risposta non si è potuta usare o salvare, con il consiglio di eseguire `ccam login <name>` se poi Claude Code ti chiede di effettuare il login.

### Il raccoglitore della riga di stato

Claude Code esegue il comando della riga di stato e gli passa tramite pipe un documento JSON. ccam ne legge il nome del modello e `rate_limits`:

<!-- motion:usage-statusline -->

1. ccam trova l'account da `CCAM_PROFILE`, che ccam imposta quando avvia Claude Code, oppure confrontando `CLAUDE_CONFIG_DIR` con le cartelle degli account (`runStatusline`).
2. ccam legge `rate_limits.five_hour` e `rate_limits.seven_day` (`ParseStatusLine`). Il JSON non contiene limiti per gli account con chiave API, né prima della prima risposta della sessione; in quel caso non viene salvato nulla.
3. ccam scrive `~/.local/state/ccam/usage/<name>.statusline.json` (`SaveStatusLine`). `ccam usage` lo usa quando l'endpoint non risponde.
4. ccam stampa `[name] · model · 5h N% · 7d N%`. Con `--wrap`, esegue invece il tuo comando tramite `sh -c` con lo stesso input e stampa il suo output come Claude Code mostra una riga di stato: ogni riga non vuota, con `[name]` davanti alla prima. Un comando che fallisce mostra solo il nome dell'account. Su Windows esegue il comando come fa Claude Code: con Git Bash se è installato, altrimenti con PowerShell.
5. Quando l'utilizzo su 5 ore o settimanale è al 90% o oltre e c'è un altro account, la riga termina con `→ ccam next --auto --continue`, per esempio `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Il comando avvolto non ha un proprio limite di tempo: Claude Code dà a un comando della riga di stato 10 minuti e ferma un'esecuzione quando inizia un aggiornamento più recente, quindi un comando lento si comporta come farebbe senza ccam. Quando Claude Code termina ccam, termina anche il comando (su Linux e macOS con ogni processo che ha avviato).

`ccam statusline install` scrive `statusLine` nel `settings.json` di ogni account e mantiene le altre chiavi dell'oggetto, come `padding`, `refreshInterval` e `hideVimModeIndicator`. Quando un account ha già un comando per la riga di stato, ccam lo mantiene come `--wrap '<old command>'`; su Windows come `--wrap-base64 <base64>`, che non richiede virgolette in nessuna delle due shell. Se la riga di stato esegue già ccam (da un altro percorso, o tra virgolette come le scriveva una versione più vecchia), ccam la sostituisce invece di avvolgerla, e mantiene il comando che avvolge.

Su Windows il percorso di ccam viene scritto con barre normali e `~` per la cartella home (`~/.local/bin/ccam.exe statusline`), che sia Git Bash sia PowerShell capiscono. Un percorso che richiede virgolette viene scritto per intero per la shell che usa Claude Code: `'…'` per Git Bash, `& '…'` per PowerShell. Un `settings.json` condiviso tramite un gruppo viene scritto una sola volta, attraverso il suo link, così lo riceve ogni account del gruppo. Eseguirlo di nuovo stampa `(already set)` per i file che usano già ccam.

## Sicurezza e casi limite

- `ccam usage` non modifica mai un file di token a meno che tu non passi `--refresh`. Quando un token di accesso è scaduto, avvia l'account una volta con `ccam <name>` e sarà Claude Code a rinnovarlo. Un login scaduto del tutto richiede `ccam login <name>`.
- `--refresh` non tocca mai gli account in esecuzione (anche quando Claude Code è stato avviato senza ccam), gli account importati o i login nel Portachiavi. Un rinnovo rifiutato mostra i numeri salvati e il motivo (`token not renewed: …`).
- L'endpoint di utilizzo non è documentato pubblicamente e potrebbe cambiare. Se succede, vedi comunque i numeri salvati, e `ccam doctor` segnala il cambio di formato.
- Senza numeri salvati e senza risposta, la riga dell'account mostra `lookup failed:` e il motivo, per esempio `no saved usage` con `--source cache`.
- I file di cache contengono solo percentuali e orari di azzeramento, mai token.

## Pagine correlate

- [Limiti di utilizzo e cambio di account](/docs/usage-limits)
- [Passare all'account successivo](/features/next)
- [Cosa fare al limite di utilizzo di Claude Code](/guides/usage-limit-workflow)
- [Riferimento dei comandi](/docs/commands)
