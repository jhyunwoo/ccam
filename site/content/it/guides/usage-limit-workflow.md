---
title: Cosa fare al limite di utilizzo di Claude Code
description: Aspetta il reset o continua subito la stessa conversazione su un altro account Claude Code. Come funzionano i limiti di 5 ore e settimanale, passo dopo passo.
updated: 2026-11-01
order: 3
translation: machine
sourceHash: 760a36724feb
---

Quando un account raggiunge il limite di utilizzo (spesso chiamato rate limit), hai due possibilità: aspettare che il limite si azzeri oppure continuare subito la stessa conversazione su un altro account che possiedi. Con ccam la seconda richiede un solo comando, `ccam next --auto --continue`.

## Quali limiti ha Claude Code?

Con i piani Pro e Max, Claude Code ha un [limite di sessione di 5 ore e un limite settimanale](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Entrambi sono condivisi con Claude sul web e nelle app, e l'utilizzo conta su entrambi contemporaneamente. Quando uno si esaurisce, Claude Code [mostra un messaggio](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) come "You've hit your session limit · resets 3:45pm". Cambiare modello non serve, perché i limiti valgono per tutti i modelli.

## Opzione 1: aspettare il reset

Se hai un solo account, aspetta. Le versioni recenti di Claude Code [aspettano nella sessione aperta](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) e riprendono il lavoro da sole quando il limite si azzera. `/usage` dentro Claude Code mostra quanta parte dei limiti del tuo piano hai usato.

## Opzione 2: continuare subito su un altro account

Se hai un secondo account, puoi continuare subito la stessa conversazione. Configuralo una volta:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Ora la riga di stato mostra l'account e il suo utilizzo su 5 ore mentre lavori, così il limite non ti coglie di sorpresa.

Quando raggiungi il limite, esci da Claude Code (`/exit`), poi:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Apri la conversazione nell'account con più margine:

```sh
ccam next --auto --continue
```

Claude Code si avvia come `work` con la conversazione più recente della cartella corrente, compreso tutto ciò che è stato detto finora. Condividendo la cronologia (`ccam group set default --add history` qui sopra) entrambi gli account lavorano sulla stessa conversazione; senza, Claude Code copia la conversazione in `work`, e ccam ti dice quanti token deve rileggere la prima risposta.

Con il [plugin di ccam](/docs/install#claude-code-plugin) in Claude Code, una sessione che si ferma al limite mostra una notifica e registra la finestra piena, così `ccam next --auto` non sceglie di nuovo quell'account prima del suo reset. `/ccam:switch` (oppure `/ccam:switch work`) annota la conversazione per `ccam next --continue`, così continua esattamente quella conversazione.

**Posti nei piani Team ed Enterprise.** La scelta automatica dell'account (`ccam next`, `--fallback`, `ccam pick`, il tasto `n` della dashboard e, una volta attivato, [`auto_switch`](/features/next#auto-switch)) passa anche da un posto all'altro della stessa organizzazione Claude Team o Enterprise, e questo può violare i termini di Anthropic. Leggi le [domande frequenti](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e verifica con la tua organizzazione e con Anthropic prima di farci affidamento.

## ccam cambia account automaticamente?

Solo se lo attivi.

- Per impostazione predefinita, ccam non è in esecuzione mentre lavori: avvia Claude Code e si sostituisce con esso, quindi decidi tu quando cambiare. `ccam next --auto` sceglie l'account con l'utilizzo su 5 ore più basso, salta gli account al 100% o senza login, e resta all'interno di una stessa organizzazione Team o Enterprise, oppure tra i piani personali.
- Con `relay = true` in `config.toml`, ccam invece resta in esecuzione accanto alla sessione, così `/ccam:switch` continua da solo la conversazione, mantenendo le opzioni, la modalità dei permessi, il livello di effort e l'agente della sessione; prima di spostare la conversazione su un account di un'altra organizzazione, ccam te lo chiede.
- Con `auto_switch = true` (Linux e macOS; sperimentale), ccam fa anche il cambio da solo: quando Claude Code ferma la sessione al suo limite e un altro account dello stesso tipo ha margine disponibile, ccam chiude la sessione e continua la conversazione lì ([condizioni e costi](/features/next#auto-switch)).

Alcuni altri strumenti invece cambiano il login di una sessione in corso; [ccam e gli altri strumenti per cambiare account](/compare/account-switchers) spiega la differenza.

## Cose da tenere a mente

- Esegui `ccam next --continue` nella stessa cartella in cui stavi lavorando; le conversazioni sono salvate per cartella di progetto.
- Evita di eseguire due account sulla stessa conversazione contemporaneamente. ccam ti avvisa se ci provi.
- Senza cronologia condivisa, la conversazione viene copiata nell'altro account e le due copie proseguono ciascuna per conto proprio; i checkpoint da cui `/rewind` ripristina il codice restano nel primo account.
- Il primo messaggio sul nuovo account può consumare più limite del solito, perché quell'account non ha ancora la conversazione in cache. `ccam next --continue --brief` la avvia invece su un breve riassunto.
- La possibilità di passare da un account all'altro in questo modo dipende dai termini che si applicano a ciascun account.

Vedi [Limiti di utilizzo](/features/usage) e [Passare all'account successivo](/features/next) per tutte le opzioni.
