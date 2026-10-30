---
title: La dashboard
description: ccam in un terminale apre una dashboard a schermo intero con account, gruppi, regole per cartella, estensioni e controlli; ogni modifica mostra prima il piano.
updated: 2026-10-30
order: 15
translation: machine
sourceHash: d65f1b5bb0c8
---

`ccam` senza argomenti apre una dashboard a schermo intero quando standard input e standard output sono un terminale; altrimenti stampa `ccam list`. La dashboard mostra ogni account con il suo utilizzo e i suoi gruppi, e fa ciò che fanno i comandi: avviare, aggiungere, effettuare il login, cambiare i gruppi, impostare le regole delle cartelle, elencare le sessioni in background, controllare e riparare. Avviare un account, o collegarsi a una sessione in background, la chiude; aggiungere un account o effettuare il login le fa cedere il terminale per un momento, poi si torna alla dashboard.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · r to refresh · a to add · ? for shortcuts
```

La riga sotto i nomi delle schede dice in una frase a cosa serve la scheda corrente. In un terminale più largo la scheda Accounts mostra anche l'email di ogni account. In ogni scheda, la fine della riga delle schede mostra se [`auto_switch`](/features/next#auto-switch) è attivo (Linux e macOS; su Windows viene ignorato). Quando è attivo, la scheda Health elenca cosa gli serve ancora: Claude Code 2.1.214 o successivo, e il plugin di ccam attivato in ogni account che accede con un abbonamento Claude.

Ogni volta che si apre, la dashboard chiede a GitHub in background se esiste una release più recente di ccam; il primo fotogramma non la aspetta mai e, quando sei offline, non viene mostrato nulla. Se ce n'è una, una riga sotto il titolo la consiglia:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` chiede prima conferma, poi aggiorna come fa [`ccam update`](/features/update): la release viene verificata con i suoi checksum SHA-256 e sostituisce ccam in un solo passaggio. La dashboard in cui ti trovi è ancora la versione vecchia, quindi esci e avvia di nuovo `ccam`. Quando l'aggiornamento non si può fare da qui, per esempio per una copia di ccam installata dal pacchetto `.deb`, la dashboard mostra lo stesso consiglio di `ccam update` e mantiene la proposta.

Il primo fotogramma viene disegnato a partire dal solo `config.toml`, in molto meno di 100 ms; l'utilizzo compare man mano che ogni account risponde. Ogni limite ha la sua barra: la finestra di 5 ore, la finestra settimanale per tutti i modelli e la finestra settimanale propria di Fable; `-` significa che l'account non ha quella finestra. Una percentuale diventa gialla all'80% e rossa al 100%. Se una consultazione fallisce, i valori in cache vengono mostrati attenuati. Una finestra che Claude Code ha segnalato piena mostra 100% fino al suo reset. Un token di accesso scaduto mostra un `⚠ token expired` attenuato (Claude Code lo rinnova al prossimo avvio), mentre `⚠ login expired` significa che l'account ha bisogno di `ccam login`.

## Schede

| Scheda | Cosa mostra | Tasti |
| --- | --- | --- |
| Accounts | numero, nome, email, barre delle 5 ore, settimanale e di Fable, ultimo uso, gruppi; `✔` l'account predefinito, `▸` l'account di questa cartella, `•` in esecuzione (anche Claude Code avviato senza ccam: `claude` da solo, un IDE, le sessioni in background e il loro daemon) | vedi sotto |
| Agents | ciò che elenca [`ccam agents`](/features/agents): le sessioni in background e le sessioni aperte di ogni account, con il loro stato, la cartella e l'ora di avvio (vedi sotto) | `enter` collegati (come `ccam attach`), `x` ferma, `r` ricarica |
| Groups | ogni gruppo, cosa condivide e i suoi account, e gli account che non fanno parte di alcun gruppo | `enter` modifica, `n` nuovo, `e` rinomina, `x` rimuovi |
| Folders | le regole di [`ccam folder`](/features/folder) e quella che vale qui, anche quando sei arrivato nella cartella tramite un link simbolico (le nuove regole vengono salvate con i link risolti) | `n` usa un account qui, `enter` cambia, `x` rimuovi |
| Extensions | ciò che mostra [`ccam ext`](/features/ext), ogni estensione una sola volta con i gruppi da cui proviene, `off:` gli account in cui un plugin non è attivato, i server propri dell'account che valgono solo per questa cartella, `This folder's project`, e i problemi che non appartengono a nessun account | `r` ricarica |
| Health | ciò che segnala [`ccam doctor`](/features/doctor), con un riquadro per account (vedi sotto) | `f` ripara il riquadro selezionato, `F` ripara tutto (`doctor --fix`), `↑`/`↓` passa da un riquadro all'altro |

`←`/`→` o `Tab` cambiano scheda, `?` elenca i tasti della scheda corrente, `r` aggiorna l'utilizzo ignorando la cache breve e ricarica le schede Agents, Extensions e Health, e `q` o `Esc` escono.

## Tasti della scheda Accounts

| Tasto | Cosa fa |
| --- | --- |
| `Enter` | Avvia Claude Code con l'account selezionato |
| `1`–`9` | Si sposta sull'account con quel numero; digita `1` `2` in rapida successione per il 12. Solo `Enter` avvia |
| `↑`/`↓`, `k`/`j` | Sposta la selezione |
| `n` | Avvia l'account che sceglierebbe [`ccam next --auto`](/features/next) |
| `o` | Avvia con opzioni: `--continue`, `--resume`, altri argomenti, oppure un dry run |
| `u` | Dettagli dell'utilizzo: le finestre su 5 ore, settimanali e per modello, con i loro azzeramenti e la provenienza dei numeri |
| `a` | Aggiunge un account (la procedura guidata qui sotto) |
| `l` | Esegue di nuovo il login |
| `s` | Sceglie i gruppi dell'account |
| `d` | Lo rende l'account predefinito |
| `f` | Lo usa nella cartella corrente |
| `e` | Lo rinomina |
| `c` | Imposta o rimuove il suo comando scorciatoia |
| `x` | Lo rimuove (chiede prima, con No selezionato) |

**Posti nei piani Team ed Enterprise.** La scelta automatica dell'account (`ccam next`, `--fallback`, `ccam pick`, il tasto `n` della dashboard e, una volta attivato, [`auto_switch`](/features/next#auto-switch)) passa anche da un posto all'altro della stessa organizzazione Claude Team o Enterprise, e questo può violare i termini di Anthropic. Leggi le [domande frequenti](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e verifica con la tua organizzazione e con Anthropic prima di farci affidamento.

## La scheda Health

La scheda Health esegue i controlli di [`ccam doctor`](/features/doctor) quando la apri e disegna la stessa vista: una prima riga con un simbolo per account, poi un riquadro per account. Ogni riga di un riquadro è un'area controllata, prima `Folder`, `Login` e `Share links`, segnata con `✔` OK, `ℹ` informazione, `⚠` un avviso o `✘` un errore, così vedi subito quale account ha un problema e dove. I risultati che non appartengono a nessun account finiscono in un riquadro `General`.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` passano da un riquadro all'altro, e il riquadro selezionato ha il bordo nel colore di accento. `f` ripara solo il riquadro selezionato, come farebbe `ccam doctor --fix`: sul riquadro di un account ripara solo quell'account, e sul riquadro `General` i risultati che non appartengono a nessun account (percorsi dei plugin e scorciatoie); tutti gli altri account vengono solo controllati di nuovo. Chiede prima (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), e un riquadro senza nulla da riparare dice solo `nothing to repair in <name>`. `F` ripara tutto in una volta, come `ccam doctor --fix`, e `r` controlla di nuovo.

## La scheda Agents

La scheda Agents mostra ciò che stampa [`ccam agents`](/features/agents): le sessioni in background e le sessioni aperte di ogni account, con l'account a cui appartiene ciascuna. Interroga contemporaneamente il Claude Code di ogni account (`claude agents --json`, ciascuno con un timeout di 15 secondi) quando apri la scheda per la prima volta, con `r`, dopo un arresto e dopo una modifica alla configurazione, e mostra `Asking every account's Claude Code for its sessions…` finché tutti non hanno risposto.

```console
  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Background agents and open sessions of every account; enter attaches to one.

    Account   ID        Name   State          Folder      Started
 ❯  work      -         -      open (busy)    ~/code/web  now
    personal  49ec0028  probe  blocked, idle  ~/code/app  12m

 ✘ lab: could not list its sessions: exit status 1 'claude agents --json' is
   not available in this environment.

 enter to attach · x to stop · r to reload · ? for shortcuts
```

`State` è lo stato di una sessione in background e ciò che sta facendo il suo processo, oppure `open` per una sessione interattiva, come in `ccam agents`. Se non è possibile interrogare un account, una riga `✘` sotto l'elenco mostra ciò che ha detto Claude Code, come `lab` qui sopra, che ha la vista degli agenti disattivata. Le sessioni in background terminate vengono escluse (`ccam agents --all` le elenca), e quando tutti gli account hanno risposto e nessuno ha una sessione, la scheda dice ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` In un terminale stretto la scheda toglie `Started`, poi `Name`, poi `Folder`, e accorcia da sinistra i percorsi delle cartelle lunghi.

`Enter` su una sessione in background chiude la dashboard e apre la sessione in questo terminale come il suo account, come fa `ccam attach <id>`; la dashboard non ritorna (in Claude Code, `←` porta alla sua vista degli agenti e `Ctrl+Z` alla shell). Una sessione aperta non ha un ID, quindi `Enter` e `x` dicono solo che serve una sessione in background. `x` chiede prima conferma, con `No, keep it` selezionato: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` Se rispondi sì, ccam esegue `claude stop <id>` come quell'account, dice ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` oppure perché non ha potuto, e rilegge l'elenco.

## Cambiare i gruppi

Nella scheda Groups, `Enter` apre l'editor: gli elementi a sinistra, gli account a destra, `Tab` per passare dall'uno all'altro, `Space` per attivare o disattivare. Un elemento che uno degli account riceve già da un altro gruppo viene mostrato disattivato con il motivo, così come un account che riceverebbe un elemento due volte; un elemento può provenire da un solo gruppo. `Enter` mostra il piano:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Se un file è diverso dalla copia del gruppo, ccam chiede prima cosa fare (tenere quella del gruppo, usare la mia o unire). Poi la modifica viene eseguita in background con lo stesso codice di [`ccam group`](/features/share): prima viene salvata la configurazione, poi i dati di ogni account vengono spostati in ordine, e gli originali vanno in `.ccam-backup/`; un account che lascia la condivisione MCP mantiene i server del gruppo nel suo `.claude.json`.

Una modifica che sposterebbe i dati di un account in esecuzione viene rifiutata indicando il nome dell'account, anche quando l'account è stato avviato mentre il piano era aperto, e un piano diventato superato perché un altro processo ccam ha cambiato i gruppi viene rifiutato con `the groups changed while this plan waited; nothing was moved, make the change again`. Mentre un altro processo ccam cambia i gruppi, una nuova modifica, rinomina o rimozione dice `another ccam process is changing the groups right now`; dopo un cambio rimasto a metà, attendono finché `F` nella scheda Health (o `ccam doctor --fix`) non lo ha completato.

## Aggiungere un account

`a` apre la procedura guidata che `ccam add` usa anche in un terminale:

1. Il nome dell'account.
2. Come accede: un abbonamento Claude subito (nel browser), una [chiave API](/features/api), oppure un abbonamento Claude più tardi.
3. Per una chiave API: da dove ccam legge la chiave a ogni avvio (`env:NAME` o `command:COMMAND`), l'endpoint, come viene inviata la chiave e il modello.
4. I gruppi a cui unirsi, con `default` selezionato (se non esiste ancora nessun gruppo, cosa condividerà il gruppo `default`).
5. Un comando scorciatoia, se ne vuoi uno.
6. Un riepilogo da confermare.

`Esc` torna indietro di un passo e, al primo passo, chiude la procedura guidata.

## Rimuovere un account

`x` chiede conferma prima che succeda qualsiasi cosa, con No selezionato, e dice cosa fa la rimozione: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` Per un account con chiave API non c'è un login da revocare, e quando un altro account ha lo stesso login, l'account non viene disconnesso. Un account importato mantiene la sua cartella: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## Colori

I colori seguono i ruoli dell'interfaccia da terminale di Claude Code, così ccam si integra bene accanto a esso; i nomi dei gruppi ricevono un colore ricavato dal loro nome. ccam chiede al terminale il suo colore di sfondo e sceglie un tema chiaro o scuro. Puoi impostarne uno tu in cima a `config.toml`, sopra il primo `[[group]]`, `[[profile]]` o `[[folder]]`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

I temi `-daltonized` sono pensati per il daltonismo rosso-verde e i temi `-ansi` usano i 16 colori del terminale stesso. `NO_COLOR` disattiva i colori.

## Sicurezza e casi limite

- La dashboard chiama le stesse funzioni dei comandi, quindi valgono tutte le regole: piani prima di spostare i dati, rifiuti per gli account in esecuzione, backup, e token di login che non lasciano mai un account.
- La dashboard non ha un dry run: `ccam --dry-run` senza un nome stampa come partirebbe l'account di questa cartella, e `o` su un account offre un dry run di quel solo avvio.
- Il terminale viene ripristinato prima che Claude Code si avvii o che si apra il login.
- Uno schermo troppo piccolo per una finestra di dialogo toglie prima le righe della scheda; le righe vengono tagliate con `…` invece di andare a capo. Su uno schermo basso spariscono poi le righe vuote e la descrizione della scheda, così la riga dei tasti resta visibile. Sotto le 40 colonne o le 8 righe la dashboard mostra solo `Too small: ccam needs 40x8 (now 36x7)`.
- `Ctrl+C` mentre si sta applicando una modifica, una rimozione, un aggiornamento o una riparazione ne attende la fine: la riga di stato dice `ccam quits when it is done.`
- Dopo che `U` ha aggiornato ccam, la proposta non ricompare quando la dashboard ritorna dall'aggiunta di un account o da un login.

## Pagine correlate

- [Condivisione con i gruppi](/features/share)
- [Avviare Claude Code](/features/launch)
- [Gestire gli account](/features/accounts)
- [Limiti di utilizzo](/features/usage)
