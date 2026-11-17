---
title: Domande frequenti
metaTitle: "Domande frequenti: più account Claude Code"
description: Risposte su ccam, account Claude Code, limiti di utilizzo, condivisione, privacy e sistemi supportati.
updated: 2026-11-17
order: 6
translation: machine
sourceHash: "9b13b37391f2"
---

## ccam è uno strumento ufficiale di Anthropic?

No. ccam è un progetto open source non ufficiale e non è affiliato ad Anthropic né approvato da Anthropic. Usa l'impostazione `CLAUDE_CONFIG_DIR`, che [Claude Code documenta](https://code.claude.com/docs/en/env-vars) per eseguire più account in parallelo.

## Posso usare più account Claude Code?

Dipende dai termini di ciascun account e, per gli account di lavoro, dalle policy della tua organizzazione: i [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) di Anthropic per Pro e Max, i suoi [Commercial Terms](https://www.anthropic.com/legal/commercial-terms) per Team ed Enterprise, e la sua [Usage Policy](https://www.anthropic.com/legal/aup) per tutti. Controllali prima di usare più account. ccam non aumenta, non unisce e non nasconde i limiti di un account, e ogni account si ferma comunque ai propri. Sposta il tuo lavoro su un altro dei tuoi account quando lo chiedi, oppure da solo a un limite quando attivi [`auto_switch`](/features/next#auto-switch) (Linux e macOS). Verifica se questo è consentito per i tuoi account prima di usarlo; la decisione è tua.

**Posti nei piani Team ed Enterprise.** La scelta automatica dell'account (`ccam next`, `--fallback`, `ccam pick`, il tasto `n` della dashboard e, se lo attivi, [`auto_switch`](/features/next#auto-switch)) passa anche da un posto all'altro della stessa organizzazione Claude Team o Enterprise. Il centro assistenza di Anthropic dice: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (ogni posto di un piano Team è assegnato nominativamente a una sola persona, e solo quella persona può usarlo) e "Users can't route their work through another user's account to use that user's remaining usage" (gli utenti non possono far passare il proprio lavoro attraverso l'account di un altro utente per sfruttare l'utilizzo residuo di quell'utente) ([piani Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); lo stesso vale per i [piani Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Usare la scelta automatica tra posti diversi può quindi violare i termini di Anthropic. Decidi tu se usarla lì: verifica con la tua organizzazione e con Anthropic prima di farci affidamento. I piani personali (Pro, Max) sono trattati nel paragrafo qui sopra.

## ccam aggira i limiti di utilizzo?

Non oltre il limite di un account: ogni account mantiene i propri limiti e si ferma a essi, e ccam non li aumenta, non li unisce e non li nasconde. ccam te li mostra e ti permette di continuare su un altro dei tuoi account quando uno li esaurisce; con `auto_switch` attivo, lo fa da solo (vedi la domanda successiva). Tra i posti di una stessa organizzazione Team o Enterprise, continuare su un altro posto usa l'utilizzo residuo di quel posto, cosa che il centro assistenza di Anthropic esclude per il posto di un altro utente (vedi Posti nei piani Team ed Enterprise qui sopra).

## ccam può cambiare account da solo quando una sessione raggiunge il limite?

Solo se lo attivi. Con `auto_switch = true` in cima a `~/.config/ccam/config.toml` (Linux e macOS; disattivato per impostazione predefinita e sperimentale), quando Claude Code ferma una sessione interattiva (una avviata da ccam) a un limite di 5 ore, settimanale, Opus, Sonnet o Fable, ccam chiude la sessione e continua la conversazione come un altro dei tuoi account con margine disponibile: un altro posto della stessa organizzazione Team o Enterprise, o un altro piano personale (Pro, Max). Non cambia mai account prima del limite, non controlla mai il tuo utilizzo in background e non sposta mai il lavoro fuori da un'organizzazione Team o Enterprise né dentro una di esse.

Sui posti sposta il tuo lavoro da un posto all'altro, cosa che può violare i termini di Anthropic (vedi Posti nei piani Team ed Enterprise qui sopra). Per i piani personali, i termini di Anthropic non affrontano chiaramente la questione: i Consumer Terms vietano di accedere ai Servizi "through automated or non-human means" (con mezzi automatizzati o non umani) salvo esplicita autorizzazione di Anthropic, e uno spostamento invia un messaggio per tuo conto. Attivarlo è una tua decisione.

[`auto_switch`](/features/next#auto-switch) elenca ogni condizione e quanto costa uno spostamento.

## Come si usano due account Claude Code su un solo computer?

Assegna a ogni account la propria cartella di configurazione con `CLAUDE_CONFIG_DIR`. ccam lo fa per te: esegui `ccam add` per ogni account, poi avviali con `ccam personal` e `ccam work`, oppure con scorciatoie a cui dai tu il nome, come `ccwork`. La [guida](/guides/multiple-claude-code-accounts) mostra anche come farlo a mano.

## Come continuo una conversazione su un altro account?

Esegui `ccam next --continue` nella cartella: porta l'ultima conversazione della cartella nell'account successivo (`--auto` per quello con più margine). Dentro Claude Code, `/ccam:switch` del [plugin di ccam](/docs/install#claude-code-plugin) fa lo stesso per la sessione in cui ti trovi. Quando i due account condividono la cronologia delle conversazioni (`ccam group set default --add history`, o un gruppo tutto loro: `ccam group add pair --items history --members work,side`), il nuovo account apre proprio quella conversazione; altrimenti Claude Code la copia nella cronologia del nuovo account. Aggiungi `--brief` per avviare il nuovo account su un breve riassunto invece che sull'intera conversazione.

## ccam vede o invia il mio login?

ccam non condivide né copia mai il tuo login di Claude tra gli account, e non copia mai gli accessi dei server MCP. Legge il token di accesso di un account solo per chiedere all'endpoint di utilizzo di Anthropic i dati di utilizzo di quell'account: per `ccam usage`, la dashboard, `ccam doctor` e la scelta automatica (`ccam next --auto`, `ccam pick`, `--fallback`, e `auto_switch` quando una sessione raggiunge il limite). Rinnova un token scaduto solo con `ccam usage --refresh`; con `usage_source = "claude"` non invia alcun token da nessuna parte e legge il login salvato (il file, o la voce nel Portachiavi di macOS) solo per capire se funziona ancora. `ccam ext` e `ccam doctor` leggono i nomi degli accessi MCP di un account, mai i loro token. ccam non ha telemetria. A parte queste consultazioni, l'unica richiesta che fa di sua iniziativa è il controllo della dashboard su GitHub per una release più recente di ccam, che non contiene nulla su di te o sui tuoi account.

## Come accedo su un server remoto?

Esegui `ccam add` o `ccam login <name>` nella sessione SSH come al solito. Claude Code attende l'accesso sul server, quindi il browser sul tuo computer finisce su una pagina `http://localhost:<port>/callback?code=…` che non si carica. Copia quell'indirizzo dal browser e incollalo nel terminale: ccam lo passa a Claude Code sul server, che completa l'accesso. L'indirizzo contiene un codice monouso, quindi non condividerlo.

## Dove vengono salvati i miei account?

In `~/.local/share/ccam/profiles/<name>/`, una cartella per account. Su macOS, Claude Code conserva il login vero e proprio nel Portachiavi, con un nome derivato da quella cartella, quindi `ccam rename` mantiene lì il vecchio nome della cartella.

## Cosa succede alla mia cartella ~/.claude esistente?

Niente, a meno che tu non esegua `ccam import` o `ccam migrate`. L'importazione registra la cartella dove si trova, quindi `claude` da solo e `ccam <name>` continuano a usare lo stesso account e le stesse impostazioni. `ccam migrate` copia le sue impostazioni, `CLAUDE.md`, le estensioni e i server MCP nel gruppo condiviso `default`, e si limita a leggere la cartella.

## Quali sistemi supporta ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine e altri), macOS e Windows 10 e 11, in modo nativo o all'interno di WSL. Su Windows, condividere file come `settings.json` richiede la Modalità sviluppatore; le cartelle condivise funzionano anche senza, e senza di essa il primo `ccam add` crea il gruppo `default` con i soli server MCP ed estensioni ([condivisione su Windows](/features/share#on-windows)).

## Perché un server MCP compare in un solo account?

I server aggiunti all'interno di una sessione con `claude mcp add` vengono salvati solo per quell'account e, senza `-s user`, solo per la cartella corrente. ccam te lo segnala la volta successiva che avvii l'account, e `ccam mcp list` e `ccam ext` li mostrano; sposta il server nel gruppo condiviso con `ccam mcp import <account> <server>`.

## ccam rallenta Claude Code?

Avviare Claude Code tramite ccam aggiunge circa 4 millisecondi, misurati con `scripts/perf-check.sh` su un portatile. ccam legge un piccolo file, controlla i link condivisi e poi si sostituisce con `claude`, quindi per impostazione predefinita non è in esecuzione mentre lavori. Windows non può sostituire un processo, quindi lì ccam avvia Claude Code e ne attende la fine, senza fare altro. Con `relay` o `auto_switch` attivi, ccam resta in esecuzione allo stesso modo accanto a ogni sessione interattiva, così può continuare la conversazione su un altro account.

## Come si aggiorna ccam?

Esegui `ccam update`. Scarica l'ultima release, ne verifica il checksum e sostituisce il file `ccam`; `ccam update --check` ti dice solo se esiste una release più recente. Anche la dashboard ti avvisa quando esiste una release più recente, e `U` la installa dopo aver chiesto conferma; ccam non si aggiorna mai da solo. Se hai installato il pacchetto `.deb`, installa invece il nuovo `.deb` con apt. Vedi [aggiornamento](/docs/install#update).

## Come si disinstalla ccam?

Esegui `ccam uninstall`. Elenca cosa verrà rimosso e chiede prima conferma, poi rimuove ogni account (dopo averlo disconnesso), i comandi scorciatoia, `~/.config/ccam`, `~/.local/share/ccam`, `~/.local/state/ccam` e il comando `ccam`. Un `~/.claude` importato resta, con una copia di ciò che condivideva nei suoi gruppi. Vedi [installazione](/docs/install#uninstall).
