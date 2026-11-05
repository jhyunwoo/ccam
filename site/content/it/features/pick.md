---
title: Scegliere un account per gli script
description: ccam pick stampa l'account con più margine per uno script, dà account diversi ai worker avviati insieme ed esce con 75 quando nessuno ha margine.
updated: 2026-11-05
order: 12
translation: machine
sourceHash: "0ec25752ace9"
---

`ccam pick` stampa il nome dell'account con login che ha più margine, così uno script può avviarci Claude Code. I worker avviati nello stesso momento ricevono ciascuno un account diverso finché abbastanza account hanno margine, così un gruppo di job `claude -p` si distribuisce sui tuoi account invece di esaurirne uno.

**Posti nei piani Team ed Enterprise.** La scelta automatica dell'account (`ccam next`, `--fallback`, `ccam pick`, il tasto `n` della dashboard e, una volta attivato, [`auto_switch`](/features/next#auto-switch)) passa anche da un posto all'altro della stessa organizzazione Claude Team o Enterprise, e questo può violare i termini di Anthropic. Leggi le [domande frequenti](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e verifica con la tua organizzazione e con Anthropic prima di farci affidamento.

## Utilizzo

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # most weekly quota left per day; fewer sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

Codici di uscita: `0` con un nome sullo standard output, `75` (`EX_TEMPFAIL`) quando nessun account ha margine o nessuno ha il login, `1` per gli altri errori, come un'opzione non valida.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--json` | Stampa l'account, il suo utilizzo su 5 ore e settimanale, la provenienza dei numeri e il motivo della scelta | disattivata |
| `--model <model>` | Conta anche la finestra settimanale di quel modello, come `opus` | nessuno |
| `--max-usage <1-100>` | Preferisce gli account sotto questo utilizzo | `auto_max_usage`, altrimenti 100 |
| `--strategy headroom\|expiring\|balance` | L'utilizzo settimanale più basso tra gli account entro 3 punti percentuali dall'utilizzo su 5 ore più basso, l'azzeramento su 5 ore più vicino oppure la maggiore quota settimanale residua per giorno | `auto_strategy`, altrimenti `headroom` |
| `--cross-org` | Sceglie anche account di altre organizzazioni | `cross_org`, altrimenti disattivata |
| `--wait[=<duration>]` | Quando nessun account ha margine, attende l'azzeramento del primo, al massimo per questo tempo. Il valore va dopo `=`: `--wait 30m` viene rifiutato con il suggerimento `--wait takes its value after '=': --wait=30m` | disattivata; `5h` se indicata da sola |
| `--dry-run` | Sceglie senza registrare la scelta | disattivata |

## Come funziona

1. ccam elenca gli account con login (gli account con chiave API e i login della Console non hanno limiti e vengono esclusi) nell'organizzazione dell'account che questa cartella avvierebbe (vedi [organizzazioni](/features/next#auto)); viene esclusa anche una cartella che ha il login nello stesso account Claude di quell'account.
2. Acquisisce un lock (`~/.local/state/ccam/picks.lock`), così le scelte eseguite nello stesso momento avvengono una dopo l'altra.
3. Consulta l'utilizzo di tutti gli account contemporaneamente, come fa `ccam usage`, con la cache di 60 secondi. Una finestra che Claude Code ha segnalato piena conta come 100% fino al suo reset ([la registrazione del limite](/features/usage#the-limit-record)), così un account che ha appena raggiunto il suo limite viene escluso anche se i numeri salvati mostrano ancora margine.
4. Gli account scelti negli ultimi 2 minuti, da `ccam pick` o da [`auto_switch`](/features/next#auto-switch) che vi ha spostato una sessione (`~/.local/state/ccam/picks/<name>`; una scelta con un orario nel futuro non conta) vengono esclusi per primi, e le [regole di scelta](/features/next#auto) di `ccam next --auto` scelgono tra i restanti: prima quelli sotto `auto_max_usage`, secondo la strategia, mai un account al 100%. Se così viene saltato l'account che le regole sceglierebbero altrimenti, il `reason` di `--json` lo dice: `30% of its 5 hours used; work was picked a moment ago`.
5. Se ogni account con margine è stato scelto di recente, le regole scelgono di nuovo tra tutti, così con tre account un quarto worker ne condivide uno.
6. La scelta viene registrata e il lock rilasciato, poi viene stampato il nome.

La finestra di 2 minuti conta perché l'utilizzo di un nuovo job impiega un po' a comparire nei numeri: senza di essa, cinque worker avviati insieme vedrebbero tutti lo stesso account come quello con più margine.

## Sicurezza e casi limite

- `ccam pick` legge solo l'utilizzo e scrive l'ora della scelta; non rinnova né legge mai i token oltre a quanto fa `ccam usage`.
- Una scelta non riserva l'account. Due script possono comunque partire sullo stesso account dopo 2 minuti; usa `--fallback` per recuperare quando quell'account raggiunge il limite.
- L'attesa con `--wait` non tiene alcun lock, quindi nel frattempo le altre scelte proseguono.
- `ccam pick` ignora [`ccam folder`](/features/folder): sceglie solo in base all'utilizzo.

## Pagine correlate

- [Passare all'account successivo](/features/next)
- [Rieseguire al limite](/features/next#run-again-at-the-limit)
- [Limiti di utilizzo](/features/usage)
- [Account con chiave API](/features/api)
