---
title: Limiti di utilizzo e cambio di account
description: Visualizza l'utilizzo su 5 ore e settimanale di ogni account Claude Code con ccam usage e apri quello con più margine con ccam next --auto.
updated: 2026-11-05
order: 5
translation: machine
sourceHash: "06bd349073e3"
---

`ccam usage` mostra, per ogni account, quanta parte di ciascun limite è stata usata: il limite su 5 ore, il limite settimanale e il limite settimanale proprio di ogni modello (come Fable). `ccam next --auto` apre l'account con più margine, e `--continue` porta con sé la conversazione corrente.

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

Ogni account ha una riga che indica quanto sono recenti i suoi dati, poi una barra per limite con il tempo mancante all'azzeramento: la finestra di 5 ore, la finestra settimanale per tutti i modelli e la finestra settimanale propria di ogni modello, come Fable. Una barra diventa gialla all'80% e rossa al 100%. L'ultima riga indica l'account che `ccam next --auto` aprirebbe; con `auto_strategy = "expiring"` o `"balance"` riporta invece `→ next --auto picks <name>: <reason>`.

**Posti nei piani Team ed Enterprise.** La scelta automatica dell'account (`ccam next`, `--fallback`, `ccam pick`, il tasto `n` della dashboard e, una volta attivato, [`auto_switch`](/features/next#auto-switch)) passa anche da un posto all'altro della stessa organizzazione Claude Team o Enterprise, e questo può violare i termini di Anthropic. Leggi le [domande frequenti](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e verifica con la tua organizzazione e con Anthropic prima di farci affidamento.

## Da dove vengono i numeri

ccam interroga l'endpoint di utilizzo di Anthropic per tutti gli account contemporaneamente, usando il login di ciascun account. Attende al massimo 3 secondi per account e mostra ogni account appena risponde.

Questo endpoint non è documentato pubblicamente e potrebbe cambiare. Quando una richiesta non va a buon fine, ccam mostra gli ultimi numeri noti e li segna come presi dalla cache. Una finestra che si è azzerata dopo il salvataggio di quei numeri mostra 0%. `ccam doctor` ti avvisa se il formato della risposta è cambiato.

Quando Claude Code ferma una sessione a un limite di utilizzo, il [plugin di ccam](/docs/install#claude-code-plugin) registra quale finestra era piena e quando si azzera. Fino ad allora quella finestra mostra 100% qualunque cosa dicano gli altri numeri, così `ccam next --auto` e `--fallback` non scelgono di nuovo l'account troppo presto.

I numeri ottenuti meno di 60 secondi prima vengono riutilizzati. Opzioni:

| Opzione | Effetto |
| --- | --- |
| `--watch[=seconds]` | Aggiorna di continuo, ogni 60 secondi per impostazione predefinita |
| `--json` | Output leggibile da una macchina |
| `-v` | Mostra anche da dove vengono i numeri di ogni account |
| `--no-cache` | Chiede sempre di nuovo |
| `--source cache` | Non chiede mai; mostra solo i numeri salvati |

## Login scaduti

Quando il token di accesso di un account è scaduto, ccam non lo rinnova, perché rinnovarlo mentre una sessione di Claude Code usa lo stesso login potrebbe disconnettere quella sessione. L'account mostra ``access token expired; `ccam <name>` renews it when it starts``: avvialo una volta con `ccam <name>` e sarà Claude Code a rinnovarlo. Un login scaduto del tutto, o cancellato da Claude Code, mostra `` login expired → `ccam login <name>` `` e richiede un nuovo login.

`ccam usage --refresh` rinnova da solo i token scaduti, ma solo per gli account senza una sessione di Claude Code in esecuzione, comunque sia stata avviata quella sessione. Non rinnova mai i login salvati nel Portachiavi di macOS o in una cartella `~/.claude` importata, e un rinnovo può richiedere fino a 30 secondi.

## Utilizzo nella riga di stato

```sh
ccam statusline install
```

Claude Code mostra quindi il nome dell'account e il suo utilizzo nella riga di stato, per esempio `[work] · Opus · 5h 24% · 7d 41%`. Dal 90% in su, la riga suggerisce anche `ccam next --auto --continue`. Se hai già un comando per la riga di stato, ccam lo mantiene e mette davanti il nome dell'account. La riga di stato salva anche i propri numeri, che `ccam usage` usa quando l'endpoint non risponde.

## Cambiare account

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the most room
ccam next --continue  # continue this folder's conversation on the next account
```

`ccam next` salta gli account che non hanno effettuato il login, e `--auto` salta gli account al 100% (compresa una finestra che Claude Code ha segnalato piena, fino al suo reset) e quelli di un'altra organizzazione. Quando i due account condividono la cronologia delle conversazioni, `--continue` apre la stessa conversazione; altrimenti Claude Code la copia nel nuovo account. Vedi [condivisione](/docs/sharing) e [passare all'account successivo](/features/next#continue).

Con `auto_switch = true` in cima a `config.toml` (Linux e macOS; disattivato per impostazione predefinita), ccam fa questo cambio da solo quando Claude Code ferma una sessione interattiva al suo limite e un altro account dello stesso tipo (un posto della stessa organizzazione, o un altro piano personale) ha margine disponibile. Vedi [`auto_switch`](/features/next#auto-switch) per sapere quando interviene e cosa non fa mai.

Per script e job pianificati che eseguono `claude -p`, metti `--fallback` prima del nome dell'account. Quando quell'account raggiunge il limite, ccam riesegue lo stesso comando con l'account con più margine, e viene stampato solo l'output dell'esecuzione portata a termine:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Vedi [Rieseguire al limite](/features/next#run-again-at-the-limit) per come decide ccam.

La [guida al flusso di lavoro per i limiti di utilizzo](/guides/usage-limit-workflow) illustra un esempio completo.
