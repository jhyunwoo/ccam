---
title: Subagenti su altri account
description: ccam agent esegue parte di un'attività di Claude Code con un altro tuo account e dà la risposta di Claude, così un'attività lunga usa l'utilizzo di più account.
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

I subagenti di Claude Code girano con l'account della sessione che li avvia, quindi consumano l'utilizzo di quell'account. `ccam agent "<task>"` esegue un'attività con `claude -p` nella stessa cartella con un altro dei tuoi account e stampa l'ultima risposta di Claude. Chiamato da una sessione di Claude Code, funziona come un subagente che consuma l'utilizzo di un altro account, così un'attività lunga può proseguire anche dopo il punto in cui l'account della sessione avrebbe esaurito il suo.

**Posti nei piani Team ed Enterprise.** `ccam agent` sceglie l'account da solo. All'interno di una stessa organizzazione Claude Team o Enterprise esegue l'attività su un altro posto, e questo può violare i termini di Anthropic. Leggi le [domande frequenti](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e verifica con la tua organizzazione e con Anthropic prima di farci affidamento.

## Utilizzo

Claude Code lo esegue con il suo strumento Bash. Con il [plugin di ccam](/docs/install#claude-code-plugin) installato, chiedilo con parole semplici, per esempio "usa i miei altri account per le parti di questa attività", oppure scrivilo nel tuo CLAUDE.md. La skill del plugin dice a Claude Code di avviare ogni agente in background, di scrivere un'attività che stia in piedi da sola e di dare agli agenti che modificano file una parte propria o un proprio git worktree.

```sh
ccam agent "find why TestParse fails and fix it; reply with what you changed"
ccam agent --model opus -- "-v breaks the parser: find out why"
ccam agent --json < task.md            # the task from standard input
ccam agent --resume 7f3c9a2e-… "now add a test for it"
ccam agent --wait=1h "…"               # nobody has room: wait for a reset
ccam --dry-run agent "…"               # the account and command it would run
```

```console
$ ccam agent "summarize what internal/usage does"
ccam: agent runs as work (12% of its 5 hours used); session 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f
internal/usage looks up the 5-hour and weekly usage of every account …
ccam: agent finished as work; follow up with: ccam agent --resume 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f "…"
```

L'ultima risposta di Claude va sullo standard output; le righe che iniziano con `ccam:` vanno sullo standard error. `--json` stampa invece `{account, session_id, result, is_error, accounts}`, dove `accounts` elenca ogni account con cui è stata eseguita l'attività.

Codici di uscita: quelli di claude stesso (`0` quando l'agente ha finito), `75` (`EX_TEMPFAIL`) quando nessun account oltre a quello chiamante ha margine, `1` per gli altri errori, come un agente che prova ad avviarne un altro.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--model <model>` | Il modello dell'agente; nella scelta conta anche la finestra settimanale di quel modello | il modello predefinito dell'account |
| `--resume <id>` | Continua la conversazione di un agente precedente con un'attività successiva (l'ID di sessione sulla sua prima e ultima riga) | disattivata |
| `--json` | Stampa l'account, la sessione, la risposta e gli account con cui è stata eseguita | disattivata |
| `--max-usage <1-100>` | Preferisce gli account sotto questo utilizzo | `auto_max_usage`, altrimenti 100 |
| `--strategy headroom\|expiring\|balance` | La [regola di scelta](/features/next#auto) | `auto_strategy`, altrimenti `headroom` |
| `--cross-org` | Esegue anche con account di altre organizzazioni | `cross_org`, altrimenti disattivata |
| `--wait[=<duration>]` | Quando nessun account ha margine, attende l'azzeramento del primo, al massimo per questo tempo | disattivata; `5h` se indicata da sola |
| `--dry-run` | Stampa il comando per l'account che sceglierebbe, senza eseguirlo | disattivata |

## Come funziona

1. **Il chiamante resta escluso.** L'account da cui viene chiamato `ccam agent` (quello della sessione avviata da ccam, altrimenti quello che questa cartella avvia) non esegue mai l'agente, e nemmeno un'altra cartella con il login nello stesso account Claude, che ne condivide i limiti. L'agente resta nell'[organizzazione](/features/next#auto) del chiamante, a meno che `--cross-org` non dica altrimenti.
2. **Sceglie come `ccam pick`.** L'utilizzo viene consultato come fa [`ccam pick`](/features/pick), e gli agenti avviati a meno di 2 minuti l'uno dall'altro ricevono account diversi finché abbastanza account hanno margine.
3. **Gira con i permessi della sessione, mai di più.** In `claude -p` non c'è nessuno che possa consentire uno strumento, quindi l'agente riceve la modalità dei permessi e il livello di effort della sessione che lo ha chiamato. Se lo avvii da un terminale, valgono la modalità predefinita di Claude Code e le regole dei permessi delle tue impostazioni. `ccam agent` non ha alcuna opzione che conceda di più.
4. **A un limite di utilizzo, la conversazione prosegue.** Quando l'account dell'agente raggiunge il suo limite, ccam registra il limite come farebbe l'[hook del plugin](/features/usage#the-limit-record) e continua la stessa conversazione con l'account che ha più margine: Claude Code ve la copia (`--resume <file> --fork-session`), con il primo messaggio "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." Ogni account viene provato una sola volta.
5. **Ogni conversazione di un agente viene registrata** in `~/.local/state/ccam/agent/<session id>.json` (30 giorni). `--resume <id>` la continua nella sua cartella, con l'account che l'ha eseguita se ha ancora margine, altrimenti con un altro. [`ccam next --continue`](/features/next) esclude le conversazioni degli agenti, così continua a riprendere le tue.

## Sicurezza e casi limite

- Un agente non può avviare agenti propri: ccam lo contrassegna con `CCAM_AGENT=1`, e lì `ccam agent` si rifiuta di partire.
- Fermare `ccam agent` (Ctrl+C, un segnale o Claude Code che interrompe il comando) ferma anche l'agente.
- L'agente non vede la conversazione che lo ha chiamato, solo la sua attività.
- Gli agenti che modificano gli stessi file contemporaneamente si sovrascrivono il lavoro a vicenda. Dai a ciascuno una parte propria dell'attività, o un proprio git worktree.
- Quando la sessione chiamante passa a un altro account ([`auto_switch`](/features/next#auto-switch)), gli agenti che ha avviato si fermano con essa. Le loro conversazioni restano: `--resume <id>` le continua.
- `agent` è il nome di un comando, quindi nessun nuovo account può chiamarsi `agent`; uno già esistente si avvia ancora con `ccam run agent`.

## Pagine correlate

- [Scegliere un account per gli script](/features/pick)
- [Passare all'account successivo](/features/next)
- [Limiti di utilizzo](/features/usage)
- [Sessioni in background di ogni account](/features/agents)
