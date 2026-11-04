---
title: Trocar para a próxima conta
description: ccam next abre o Claude Code com outra conta; --auto escolhe a de mais espaço, --continue leva a conversa junto e --fallback reexecuta jobs -p.
updated: 2026-11-05
order: 8
translation: machine
sourceHash: "de1fd573f1ac"
---

`ccam next` abre o Claude Code com outra conta sem que você precise dizer qual. Use-o quando uma conta atingir o limite de uso. Há quatro formas de trocar:

- O `ccam next` simples pega a conta seguinte à última que você usou, e `--auto` pega a que tem mais espaço. `--continue` continua a última conversa desta pasta na nova conta, compartilhem as duas o histórico ou não.
- Dentro do Claude Code, o [plugin do ccam](#inside-claude-code) avisa quando um limite interrompe uma sessão, e `/ccam:switch` passa a conversa adiante.
- Para scripts que executam `claude -p`, `ccam --fallback` faz a troca sozinho (consulte [Executar de novo no limite](#run-again-at-the-limit)).
- Com o [`auto_switch`](#auto-switch) ativado (Linux e macOS), o ccam também move uma sessão interativa que o Claude Code interrompe no limite dela.

**Assentos dos planos Team e Enterprise.** A escolha automática de conta (`ccam next`, `--fallback`, `ccam pick`, a tecla `n` do painel e, depois que você o ativa, o [`auto_switch`](/features/next#auto-switch)) também alterna entre assentos da mesma organização Claude Team ou Enterprise, o que pode violar os termos da Anthropic. Leia as [perguntas frequentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e consulte a sua organização e a Anthropic antes de contar com ela.

## Uso

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

O ccam imprime a troca como `→ <from> → <to>` na saída de erro padrão, com o motivo de uma escolha do `--auto`, e depois inicia o Claude Code com a nova conta, do mesmo jeito que `ccam <name>`. `<from>` é a última conta que você usou ou, com `--continue`, a conta à qual a conversa pertence. Os argumentos depois de `--` vão para o `claude` sem alterações.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--auto` | Escolhe a conta com login feito pelo uso (veja abaixo), pulando as contas em 100% | desativado |
| `--max-usage <1-100>` | Com `--auto`: prefere as contas abaixo deste uso | `auto_max_usage`, senão 100 |
| `--strategy headroom\|expiring\|balance` | Com `--auto`: o menor uso de 5 horas, o reinício de 5 horas mais próximo ou a maior cota semanal restante por dia | `auto_strategy`, senão `headroom` |
| `--cross-org` | Com `--auto`: também escolhe contas de outra organização (veja abaixo) | `cross_org`, senão desativado |
| `--wait[=<duration>]` | Com `--auto`: quando nenhuma conta tem espaço, espera o reinício da primeira, no máximo por este tempo. O valor vai depois de `=`: `--wait 2h` é recusado com a dica `--wait takes its value after '=': --wait=2h` | desativado; `5h` quando usado sozinho |
| `-c`, `--continue` | Continua a última conversa desta pasta na nova conta (veja abaixo) | desativado |
| `--brief` | Com `--continue`: inicia uma conversa nova na nova conta a partir do resumo que o ccam faz da última, em vez de retomá-la inteira ([veja abaixo](#brief)) | desativado |
| `-- <claude args>` | Argumentos passados ao `claude` depois dos do próprio ccam | nenhum |
| `--dry-run` | Aceito: imprime `→ <from> → <to>` e a linha de comando que executaria, sem iniciar o Claude Code nem registrar uma inicialização | desativado |

`--auto` e `--continue` funcionam juntos. Com `--continue`, os argumentos de `--resume` são colocados antes dos seus.

As opções do próprio Claude Code, como `--dangerously-skip-permissions` ou `--model`, vão depois de `--`. Diferente de `ccam <name>`, que repassa tudo o que vem depois do nome, o `ccam next` lê os próprios argumentos, então sem o `--` ele para com `ccam: unknown flag: --dangerously-skip-permissions`.

## Como funciona

Cada inicialização pelo ccam grava o horário em `~/.local/state/ccam/last-used/<name>`. `ccam next` lê esses arquivos, verifica quais contas estão com login feito e escolhe.

1. O ccam monta um `next.Candidate` por conta a partir do `config.toml` (`cli.candidates`): a conta, o horário da última inicialização (`launcher.LastUsed`) e se ela tem um login que funciona (`creds.Has`: um login do Claude no `.credentials.json`, ou no macOS no Acesso às Chaves, que não expirou por completo nem foi apagado pelo Claude Code). Uma pasta que usa uma chave de API em vez disso (um login do Console, ou uma conta adicionada com `ccam add --api`) não tem limites de plano e nunca é escolhida.
2. Sem `--auto`, o ccam ordena as contas pelo número do slot e encontra a iniciada mais recentemente (`next.Rotate`). Essa é a conta da qual você está saindo. Se nenhuma conta foi iniciada ainda, é o menor slot.
3. A partir dela, avança na ordem dos slots, voltando do último slot para o primeiro, e retorna a primeira conta com login feito.
4. O ccam prepara e inicia a conta escolhida exatamente como `ccam <name>` faria (`launcher.LaunchWith`): verifica os links da conta, aplica o compartilhamento de MCP, registra a sessão e substitui o ccam pelo `claude` com `syscall.Exec`.

A busca começa na última conta usada e para na primeira conta com login feito depois dela:

<!-- motion:next-rotate -->

`work` foi a iniciada mais recentemente, então a busca começa depois dela. `side` não tem login e é pulada, e `lab` é escolhida. Se `lab` também estivesse sem login, a busca voltaria ao início, em `personal` (`% len` no diagrama). A rotação nunca retorna a conta de onde você partiu: sem outra conta com login feito, ela para com um erro.

### --auto

`--auto` decide pelo uso em vez da ordem dos slots. O ccam mantém as contas com login feito, consulta o uso de todas ao mesmo tempo do mesmo jeito que o `ccam usage` (`Fetcher.FetchAll`) e escolhe pelas regras abaixo (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

Uma conta é avaliada pelas janelas de 5 horas e semanal e, quando `--model` está entre os argumentos para o Claude Code (depois de `--`), também pela janela semanal desse modelo (`opus`, `sonnet`, `haiku` ou `fable` no nome). O ccam divide as contas em três níveis (`Policy.Choose()`):

1. **Com espaço**: todas as janelas abaixo de `auto_max_usage` (`config.toml`, padrão 100). A estratégia escolhe entre elas: `headroom` (padrão) pega o menor uso de 5 horas e, entre contas com o mesmo uso de 5 horas, a que usou a menor parte da semana; `expiring` pega a janela de 5 horas que reinicia antes, para que a cota restante seja usada antes de se perder; `balance` pega a conta com mais cota semanal restante por dia até a semana dela reiniciar, contando cada sessão do ccam em execução contra ela, para que a cota semanal não fique sem uso e as sessões em paralelo se distribuam.
2. **Acima do limite**: abaixo de 100%, mas em `auto_max_usage` ou acima. Usadas só quando nenhuma conta tem espaço; a que tem o menor uso de 5 horas vem primeiro e, em caso de empate, a que usou a menor parte da semana.
3. **Cheias**: uma janela em 100%, ou uso de 5 horas desconhecido. Nunca escolhidas.

Uma janela conta como 100% até reiniciar quando o Claude Code recusou uma solicitação ali: o hook do [plugin do ccam](#inside-claude-code) registra isso (`~/.local/state/ccam/usage/<name>.limit.json`), não importa o que digam os números salvos. Assim, uma conta que acabou de atingir o limite não é escolhida de novo antes do reinício, mesmo quando os números dela são antigos.

As configurações para todas as contas ficam no início do `config.toml`, acima do primeiro `[[group]]`, `[[profile]]` ou `[[folder]]`. Uma configuração abaixo de um deles pertence àquela tabela: o ccam avisa sobre ela a cada inicialização, e os comandos que salvam a configuração se recusam até que ela seja movida.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organizações.** As escolhas automáticas ficam entre contas do mesmo tipo da que você está deixando: assentos da mesma organização Team ou Enterprise, ou planos pessoais (Pro, Max), para que o trabalho de uma conta da empresa não continue em uma conta particular, nem o contrário. O ccam sabe a organização de cada conta pelo login dela (`ccam login` a registra; logins mais antigos são lidos do `.claude.json`); uma conta cuja organização o ccam não consegue identificar pode ser combinada com qualquer conta. Duas pastas com login na mesma conta Claude compartilham os limites dela, então uma nunca substitui a outra. `--cross-org`, ou `cross_org = true` no `config.toml`, suspende a regra da organização, mas depois de [`/ccam:switch`](#inside-claude-code) uma troca para outra organização ainda espera o seu sim.

Quando nada pode ser escolhido, o ccam sai com código 75 (`EX_TEMPFAIL`), para que um script consiga distinguir "tente mais tarde" de uma falha. Com `--wait`, ele espera até a janela que bloqueia a primeira conta cheia reiniciar (30 segundos depois do horário de reinício, quando o endpoint mostra a nova janela), imprime `ccam: no account has room; waiting until 15:10` e escolhe de novo; quando nenhuma conta informa quando reinicia, ele verifica de novo depois de 10 minutos. Se nada se liberar antes do fim da espera, sai com 75.

A consulta segue as regras do `ccam usage`: números com menos de 60 segundos são reaproveitados, a consulta de cada conta tem um tempo limite de 3 segundos, uma conta que não responde, ou cujo token de acesso expirou, é avaliada pelos números salvos, e um registro de limite prevalece sobre eles. Uma janela salva cujo horário de reinício já passou conta como 0%, então uma conta cuja janela reiniciou nesse meio-tempo não é pulada. `--auto` nunca renova tokens. Diferente da rotação simples, `--auto` pode escolher a última conta usada se ela ainda tiver mais espaço.

### --continue

1. O ccam procura a última conversa desta pasta em todas as contas: dos arquivos `<session>.jsonl` em `projects/<folder>/` (o nome da pasta é a pasta de trabalho, com os links resolvidos e todo caractere que não seja letra ou dígito trocado por `-`, como o Claude Code a nomeia), aquele com o prompt ou a resposta mais recente. Sessões sem nenhum dos dois, como uma sessão que só executou `/mcp`, `/usage` ou `/status`, ou um arquivo vazio, são puladas. Se uma sessão deixou uma nota (um limite de uso ou `/ccam:switch`, veja abaixo) e nenhuma conversa tem um prompt ou uma resposta mais de 2 minutos depois dela, é essa sessão exata que é usada; as linhas que o Claude Code grava no `/exit` não contam.
2. A conta à qual ela pertence é a conta da qual você sai, e a nova conta é escolhida como acima, dentro da mesma organização. Quando várias contas compartilham esse histórico, ela pertence àquela cujo `.claude.json` diz que executou a conversa por último (`lastSessionId`); caso contrário, à usada por último.
3. Se a nova conta compartilha esse histórico (a mesma pasta de grupo), o ccam inicia `claude --resume <session id>`, exatamente aquela conversa (`--continue` abriria a mais recente da pasta, qualquer que fosse).
4. Caso contrário, ele inicia `claude --resume <the conversation's file> --fork-session`: o Claude Code copia a conversa inteira, com um novo ID de sessão, para o histórico da própria nova conta, e a original fica como estava. A partir daí, o `--continue` da nova conta a encontra.
5. A menos que você passe `--agent` depois de `--`, o ccam também repassa o agente com o qual a conversa roda: a configuração `agent` do `.claude/settings.local.json` ou do `.claude/settings.json` da pasta ou do `settings.json` da conta; caso contrário, o agente que a última linha `agent-setting` da conversa indica. O Claude Code não traz o agente de volta quando copia uma conversa (passo 4) e seguiria como o agente padrão dele, com todas as ferramentas. Quando a nova conta não compartilha os agentes, plugins e configurações da antiga (as mesmas pastas `agents` e `plugins` e o mesmo `settings.json`, por meio de um grupo), o ccam adiciona uma linha depois da que fala da troca: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

O ccam informa o tamanho da conversa: a nova conta lê tudo de novo na primeira resposta. Os caches de prompt nunca são compartilhados entre organizações, e cada plano pessoal é uma organização própria, então, quando as duas contas pertencem a organizações diferentes, o ccam acrescenta `prompt caches are not shared between organizations, so the first reply reads it all again`. Quando a conta antiga tem login em servidores MCP em que a nova não tem, o ccam os indica: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. Compartilhar `history` em um grupo evita a cópia, mantém uma só conversa e também compartilha os checkpoints a partir dos quais o `/rewind` restaura o código (`file-history/`).

### --brief

O `--continue` retoma a conversa inteira, e a nova conta lê tudo de novo na primeira resposta. O `ccam next --continue --brief` inicia uma conversa nova em vez disso, a partir de um resumo que o ccam escreve da antiga sem consultar um modelo:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

O resumo traz, nesta ordem e só quando há algo a dizer:

1. o título da conversa;
2. o caminho da transcrição dela, para o Claude pesquisar quando precisar de um detalhe;
3. o seu primeiro pedido, quando houve mais de quatro;
4. o resumo que o Claude Code escreveu na última vez que compactou a conversa;
5. os seus quatro últimos pedidos, com comandos `!` e comandos de barra como você os digitou;
6. a última resposta do Claude;
7. a lista de tarefas;
8. os 10 arquivos que o Claude alterou por último;
9. o `git status --short --branch` da pasta, com até 800 caracteres.

A lista de tarefas guarda no máximo 10 itens: os itens concluídos são descartados primeiro, depois as etapas mais distantes, para que a etapa atual e as seguintes fiquem. Partes longas são cortadas no meio, e o resumo inteiro fica abaixo de 16.000 caracteres.

O ccam o salva como `~/.local/state/ccam/brief/<session>.md`, legível só por você, e o próximo `--brief` remove os resumos com mais de 7 dias. O Claude Code começa com um prompt e anexa o arquivo que ele menciona: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` Um prompt que você passa depois de `--` toma o lugar da última frase (`Then: <your prompt>`), e as suas outras opções depois de `--` vêm em seguida.

O Claude Code termina uma menção em um `"` e a corta em um `#`, então, quando o caminho do resumo tem um deles, o ccam o menciona pelo caminho relativo à pasta de trabalho; quando esse caminho também tem um, o prompt cita o arquivo pelo nome, e o Claude o lê com a ferramenta Read, que pergunta antes porque o arquivo está fora das pastas de trabalho (em uma execução `-p` no modo Manual, o Claude Code nega a leitura e o Claude continua sem o resumo).

`ccam --dry-run next --continue --brief` imprime o resumo e não salva nada. A conversa antiga fica como estava, e a nova só sabe o que o resumo diz, então use `--brief` quando a tarefa importa mais do que cada detalhe de como você chegou até ali. Sem `--continue`, o ccam para com `--brief works with --continue`; quando não consegue ler a conversa ou não encontra nada nela para resumir, ele para e diz que `ccam next --continue` sem `--brief` a retoma inteira. O [relay](#relay) e o [`auto_switch`](#auto-switch) sempre retomam a conversa inteira. A opção `--brief` do próprio Claude Code vai depois de `--`, como as outras opções dele.

## Dentro do Claude Code

O [plugin do ccam](/docs/install#claude-code-plugin) adiciona dois hooks (Claude Code 2.1.78 ou posterior):

- Quando um turno para em um limite de uso, ele anota a conversa e mostra uma notificação na área de trabalho nos terminais compatíveis (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  Ele conhece todas as mensagens que o Claude Code trata como limite de uso: os limites de 5 horas, semanal, simples e de modelo (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` para a cota de um modelo), créditos de uso e limites de gasto, e uma organização que bloqueou a conta.

  Quando o Claude Code gravou a solicitação recusada na conversa, o hook também registra a janela cheia até o reinício ([o registro de limite](/features/usage#the-limit-record)), exceto em um bloqueio definido por uma organização, que os administradores dela podem suspender a qualquer momento. Uma limitação breve do servidor ou o limite de taxa de uma chave de API não é um limite de uso e é ignorado, assim como as sessões sem supervisão (execuções `-p` e sessões em segundo plano). Um subagente ou agente de workflow que para em um limite de uso conta como a parada da sua sessão, mesmo enquanto a conversa principal apenas espera por ele: o Claude Code grava essa solicitação recusada no próprio arquivo do agente, na pasta `subagents` da sessão, e o hook a lê ali. Quando vários agentes param ao mesmo tempo, a sessão recebe uma só nota e uma só notificação.

  O Claude Code 2.1.234 e posteriores também esperam na sessão e continuam sozinhos depois de um reinício de 5 horas; a notificação dá a você a opção de trocar de conta em vez disso, e, com o [`auto_switch`](#auto-switch) ativado (Linux e macOS), o ccam troca por você quando outra conta tem espaço. Quando nenhuma outra conta pode assumir, porque nenhuma outra conta pertence à mesma organização Team ou Enterprise, ou nenhuma é outro plano pessoal (com `cross_org` ativado, qualquer outra conta Claude pode), a notificação termina com `You have no other personal plans, so ccam cannot move this conversation.` ou `You have no other accounts in Acme, so ccam cannot move this conversation.` em vez disso.
- `/ccam:switch [account]` anota a conversa e encerra o turno antes que qualquer solicitação seja feita: digite `/exit` e depois execute `ccam next --continue` para continuá-la na conta que você indicou (ou na próxima). Em uma sessão que ainda não tem conversa, ele avisa que não há nada para continuar. Só o que você digita o executa: a partir da versão 0.3.2 do plugin, o Claude não consegue executar `/ccam:switch`, e um job que o Claude agenda (CronCreate, ScheduleWakeup ou `.claude/scheduled_tasks.json`) o passa ao Claude como texto simples.

  Quando a conta que você indica está fora da organização da conversa (um assento de outra organização Team ou Enterprise, um plano pessoal no lugar de um assento ou o contrário, ou uma conta com chave de API em qualquer um dos lados), `ccam next --continue` pergunta antes, não importa o que `cross_org` diga: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. Ele pergunta o mesmo antes de uma conta que ele mesmo escolhe quando você não indica nenhuma, e que `cross_org = true` permite que ele pegue de outra organização: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. Só `y` ou `yes` a move; no Linux e no macOS, as teclas digitadas antes da pergunta são descartadas primeiro. Qualquer outra resposta, ou nenhum terminal onde perguntar, deixa a conversa e a nota onde estão: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (depois de uma escolha, `ccam next --continue` escolhe de novo). Como nas escolhas automáticas, uma conta cuja organização o ccam não consegue identificar combina com qualquer uma.

A nota fica em `~/.local/state/ccam/handoff/<account>/<session>.json`, um arquivo por sessão, para que duas sessões de uma conta nunca sobrescrevam a nota uma da outra. Ela é mantida por 12 horas e usada uma vez, quando a próxima sessão está para começar; uma inicialização que falha a deixa para a próxima tentativa. Ela pertence à pasta em que a sessão começou, mesmo depois de um `cd` no shell do Claude Code. `ccam rename` atualiza os nomes nela; se a conta que ela indica foi removida desde então, o ccam avisa e escolhe como de costume.

### Relay

Com `relay = true` no início do `config.toml`, o ccam continua rodando ao lado das sessões interativas em vez de se substituir pelo Claude Code e, depois de `/ccam:switch`, continua sozinho a conversa daquela sessão (uma nota deixada por outra sessão não é usada): o hook encerra o Claude Code, e o ccam executa `ccam next --continue` com a nota daquela sessão. Sem uma conta indicada, ele pega a que tem mais espaço, contando o modelo da conversa, e, quando `cross_org = true` permite que ela seja uma conta de outra organização, faz a pergunta acima antes de mover a conversa para lá. No Windows, o hook não consegue encerrar o Claude Code: digite `/exit`, e o ccam a continua. O hook também não encerra o Claude Code para uma conta fora da organização da conversa: ele diz `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`, e o ccam faz a pergunta acima assim que você digitar `/exit`. Isso custa um processo do ccam esperando por sessão, e é por isso que vem desativado por padrão.

A sessão que dá continuidade mantém o que retomar uma conversa não restaura: opções como `--add-dir`, `--mcp-config`, as listas de ferramentas, `--settings`, `--plugin-dir` e as opções de prompt de sistema, o modo de permissão da conversa, o nível de esforço da sessão e o agente com o qual a conversa roda (veja [`--continue`](#continue)). O próprio Claude Code restaura `--model` e `--name` a partir da conversa. Qualquer outra coisa (`--worktree`, `--debug`, `--remote-control`, opções que o ccam não conhece) é indicada em `ccam: not carried over: …`. A troca nunca faz a pergunta do histórico compartilhado, e o que você digitou enquanto ela rodava é descartado; a única pergunta dela é a que vem antes de uma conta de outra organização.

O relay só continua a conversa quando o Claude Code terminou do jeito que o `/ccam:switch` pediu: não quando você digitou `/exit` antes, pressionou Ctrl+C enquanto ele terminava, ou o ccam recebeu um hangup ou outro sinal, nem quando o terminal sumiu. A nota então fica, e `ccam next --continue` a continua depois. Uma sessão iniciada pelo relay de uma versão mais antiga do ccam, ou uma sessão sem supervisão (execuções `-p` e sessões em segundo plano), não é encerrada pelo hook; ela recebe a mensagem do `/exit` em vez disso. `claude remote-control` roda sem relay.

### `auto_switch`

Com `auto_switch = true` no início do `config.toml` (Linux e macOS), o ccam move sozinho uma sessão para outra conta quando o Claude Code a interrompe em um limite de uso: o hook do plugin escolhe outra conta com espaço e encerra o Claude Code, e o ccam, rodando ao lado da sessão, continua a conversa nessa conta, do mesmo jeito que o [relay](#relay) faz depois de `/ccam:switch`. Ele vem desativado por padrão e é experimental.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

O ccam só move uma sessão quando todas estas condições valem:

1. O ccam a iniciou em um terminal e continua rodando ao lado dela. O `auto_switch` ativa o [relay](#relay), então isso vale para as sessões interativas que você inicia depois de ativá-lo. Uma sessão iniciada antes, ou por um ccam mais antigo, recebe uma notificação que termina com ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` Sessões sem supervisão (execuções `-p` e sessões em segundo plano), `claude remote-control` e o Claude Code iniciado sem o ccam nunca são movidos.
2. A conta tem o [plugin do ccam](/docs/install#claude-code-plugin) ativado e roda o Claude Code 2.1.214 ou posterior, que passa o ID do próprio processo ao hook (`CLAUDE_PID`).
3. O Claude Code gravou a solicitação recusada na conversa, com o horário em que ela reinicia, para o limite de 5 horas, semanal, do Opus ou do Sonnet dessa conta. Um limite de gasto, o fim dos créditos de uso, um bloqueio definido por uma organização, a cota do Fable (`You've reached your Fable limit.`), uma limitação do servidor e o limite de taxa de uma chave de API nunca movem uma sessão. O limite de um subagente ou agente de workflow conta, mesmo enquanto a conversa principal apenas espera por ele; quando é um limite semanal do Opus ou do Sonnet, a sessão só é movida se a própria conversa rodar nesse modelo.
4. O ccam conhece o plano da conta, pelo `.claude.json` dela ou pelo que o `ccam login` registrou. Se não conhecer, a notificação termina com `ccam does not know personal's plan, so auto_switch leaves this session.`
5. A sessão não está conectada ao Remote Control. Como outra conta, ela não se reconectaria ao app de onde você a controla, então ela fica, e a notificação diz isso.
6. O limite reinicia daqui a mais de 15 minutos. Para uma espera menor, o ccam deixa a sessão com o Claude Code: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. Você não digitou nenhum prompt desde o limite.
8. Outra conta tem espaço. O ccam escolhe como o [`--auto`](#auto) faz (`auto_max_usage`, `auto_strategy`, o modelo da conversa), entre contas do mesmo tipo da que você deixa: outro assento da mesma organização Team ou Enterprise, ou outro plano pessoal (Pro, Max). `cross_org` e `--cross-org` não se aplicam aqui. A conta precisa:
   - estar logada com um login que o Claude Code consegue renovar;
   - não estar logada na mesma conta Claude;
   - não usar uma chave de API;
   - ter um plano que o ccam conhece;
   - não ter créditos de uso ativados, até onde a própria consulta de uso do ccam mostra;
   - ter confiado nesta pasta no Claude Code;
   - se a conversa roda com um agente (a configuração `agent` ou a conversa indica um), compartilhar os agentes da conta que você está deixando: as pastas `agents` e `plugins` e o `settings.json` dela precisam ser, cada um, exatamente os mesmos (por meio de um grupo) ou não existir em nenhuma das duas.

   Quando essas três últimas condições deixam de fora todas as contas com espaço, a notificação as indica: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` ou `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` Uma janela com um [registro de limite](/features/usage#the-limit-record) conta como cheia, e uma conta para a qual outra sessão foi movida nos últimos 2 minutos fica depois das outras. A escolha precisa ser feita em cerca de 6 segundos; quando as consultas de uso demoram mais, a sessão fica. Com `usage_source = "claude"`, o `auto_switch` decide só a partir de números salvos, e o ccam não consegue saber, pela consulta do Claude Code nem pela linha de status, se os créditos de uso estão ativados, então supõe que estão desativados.

Então o hook encerra o Claude Code, e o ccam, rodando ao lado dele, mostra uma notificação na área de trabalho nos terminais compatíveis, como `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` Ele imprime a primeira linha da saída de exemplo acima e executa `ccam next --continue` para essa sessão como a conta escolhida, que imprime as outras duas linhas e inicia o Claude Code com as opções, o modo de permissão, o nível de esforço e o agente da sessão, como o [relay](#relay) faz. A primeira mensagem da nova sessão é um texto fixo, a única coisa que o ccam envia em seu nome (uma troca depois de `/ccam:switch` não envia nada):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

Depois do limite de um subagente, é esta outra, porque os subagentes, os workflows e os comandos em segundo plano pararam junto com a sessão antiga:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

Se essa conta também atingir o limite, o mesmo acontece nela. Os registros de limite impedem que a cadeia volte para uma conta antes do reinício dela, então, quando todas as contas estão cheias, a sessão fica na última.

O hook descobre que uma conversa roda com um agente pela configuração `agent` ou pela linha que indica o agente, que o Claude Code grava na conversa de tempos em tempos e quando a sessão termina. Se até então só o `--agent` da sessão o indica, e a conta escolhida não compartilha os agentes, `ccam next --continue` não inicia nada e mantém a nota, e uma segunda notificação na área de trabalho avisa isso: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

Quando nenhuma outra conta tem espaço, nada é movido: o Claude Code trata o limite como de costume, esperando na sessão ou oferecendo as opções dele, e a notificação diz qual conta libera primeiro, por exemplo `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` Para assentos, ela diz `None of your other accounts in Acme has room now`. Para cancelar uma troca que já começou, pressione Ctrl+C enquanto o Claude Code fecha: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

Uma troca tem custos. O Claude Code reinicia, o que leva alguns segundos, e os shells em segundo plano, servidores de desenvolvimento e agentes que o Claude iniciou param junto com a sessão antiga. A nova conta lê a conversa inteira de novo na primeira resposta, como com [`--continue`](#continue), e o ccam informa o tamanho dela; uma troca nunca usa [`--brief`](#brief).

O `auto_switch` nunca troca antes de o Claude Code informar um limite, nunca acompanha o seu uso em segundo plano e nunca copia, renova ou repassa um login. No Windows, ele é ignorado, e o `ccam doctor` diz isso; lá, com `relay = true`, `/ccam:switch` e `/exit` ainda movem uma conversa para outra conta.

**Termos da Anthropic.** Em assentos Team e Enterprise, o `auto_switch` move o seu trabalho para outro assento da mesma organização. A central de ajuda da Anthropic diz "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada assento de um plano Team é atribuído a uma pessoa específica, e só essa pessoa pode usá-lo) e "Users can't route their work through another user's account to use that user's remaining usage" (os usuários não podem encaminhar o próprio trabalho pela conta de outro usuário para aproveitar o uso restante desse usuário) ([planos Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); o mesmo vale para os [planos Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)).

Para planos pessoais (Pro, Max), os termos da Anthropic não dizem se as suas próprias contas podem assumir automaticamente umas pelas outras em um limite. Os [Termos do Consumidor](https://www.anthropic.com/legal/consumer-terms) proíbem acessar os Serviços "through automated or non-human means, whether through a bot, script, or otherwise" (por meios automatizados ou não humanos, seja por um bot, um script ou de outra forma), a menos que a Anthropic permita isso explicitamente, e uma troca envia uma das mensagens fixas acima em seu nome. A [Política de Uso](https://www.anthropic.com/legal/aup) que entra em vigor em 12 de novembro de 2026 lista "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (contornar ou interferir intencionalmente em recursos, restrições ou proteções estabelecidos nos nossos produtos) entre as coisas que ela proíbe.

O ccam não afirma que o `auto_switch` cumpre esses termos. Ativá-lo é uma decisão sua, e o risco também é seu: a Anthropic pode suspender uma conta, ou o acesso de uma organização Team ou Enterprise, quando entender que os termos foram violados.

Enquanto o `auto_switch` está ativado, o `ccam doctor` avisa sobre uma conta sem o plugin do ccam e sobre versões do Claude Code anteriores à 2.1.214 e, para assentos, repete a regra de assentos da Anthropic ([verificações](/features/doctor)).

## Executar de novo no limite

`ccam next` serve para uma sessão interativa que parou. Uma execução não interativa, `claude -p` em um script ou em um job agendado, não pode esperar por você, então `--fallback` antes do nome da conta deixa o ccam fazer a troca por ela:

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

1. `--fallback` só funciona junto com `-p` ou `--print`; caso contrário, o ccam para e indica `ccam next --auto --continue`.
2. A entrada padrão vinda de um pipe é lida uma vez, até 10 MiB, e entregue a cada tentativa. Como o Claude Code, o ccam espera no máximo 3 segundos para ela começar (exceto com `--input-format stream-json`); se um pipe fica aberto, mas não envia nada, a execução fica sem entrada, e o ccam avisa isso na saída de erro padrão.
3. O ccam inicia o `claude` como processo filho em vez de usar `exec`, retém a saída padrão dele e repassa a saída de erro à medida que ela chega.
4. Quando o `claude` sai com 0, ou falha por outro motivo, a saída retida e o código de saída são repassados sem alteração. O mesmo vale para uma execução encerrada por um sinal (um código de saída acima de 128): pediram que ela parasse, e ela não roda de novo.
5. Uma falha conta como limite de uso quando a mensagem do Claude Code diz isso, qualquer uma das mensagens de limite de uso que o [hook](#inside-claude-code) conhece, também dentro da saída de `--output-format json` ou `stream-json`, ou quando o uso da conta mostra a janela de 5 horas, a semanal ou a do `--model` em 100% (incluindo um registro de limite).
6. Então a saída dessa tentativa é descartada e o mesmo comando roda com outra conta logada: a que tem o menor uso de 5 horas, depois as contas cujo uso não pode ser consultado, pela ordem dos números. Contas em 100% são puladas, contas com login por chave de API nunca são tentadas, e cada conta é tentada no máximo uma vez.
7. Com `-c`, `--continue`, `-r` ou `--resume` entre os argumentos, só são tentadas contas do mesmo grupo de histórico, para que a conversa possa ser encontrada.
8. Só são tentadas contas da mesma organização, como com `--auto` (`cross_org = true` suspende isso).

A próxima conta é escolhida pelas mesmas regras do `--auto` (`auto_max_usage`, `auto_strategy`, a janela do `--model`). Quando não sobra nenhuma conta, o ccam imprime a saída da última tentativa, avisa que nenhuma conta logada tem espaço e sai com código 75. Com `--wait` antes do nome da conta (`ccam --fallback --wait=2h work -p …`), ele espera o reinício da primeira conta, depois pode tentar todas as contas de novo; ele sempre espera pelo menos uma vez antes, para que uma conta que acabou de atingir o limite não rode de novo na hora. A duração limita o tempo total de espera da execução, não cada espera. Escrita com espaço (`--wait 30m`), ela é recusada com a mesma dica do `ccam next`. Um `--` antes dos argumentos do Claude Code funciona como em qualquer inicialização (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` imprime o comando da primeira tentativa.

Sessões interativas não trocam de conta sozinhas, a menos que você ative o [`auto_switch`](#auto-switch). Com `ccam statusline install`, a linha de status acrescenta `→ ccam next --auto --continue` quando o uso de 5 horas ou o semanal chega a 90% ([Limites de uso](/features/usage)).

## Segurança e casos especiais

- Contas sem login nunca são escolhidas. Sem outra conta com login feito, o `ccam next` simples para com `no logged-in account other than <name>`.
- Com `--continue` e nenhuma conversa na pasta: `no account has a conversation in <folder> to continue`.
- Quando a regra da organização não deixa nenhuma conta: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Com `--auto` e nenhuma conta com login feito: `no logged-in account to pick`. Quando todas as contas estão cheias ou com uso desconhecido: `` no account with known usage and headroom (see `ccam usage`) ``. Os dois casos terminam com código de saída 75.
- `--max-usage`, `--strategy`, `--wait` e `--cross-org` só funcionam com `--auto`, e um valor fora de 1-100 ou uma estratégia desconhecida é recusado.
- Um fork copia a conversa, não o restante do estado da conta: os checkpoints do `file-history/` da conta original não vêm junto, a menos que o histórico seja compartilhado.
- `--auto` faz solicitações de rede mesmo com `--dry-run`. Ele só lê e nunca grava arquivos de token.
- A conta escolhida é iniciada pelo mesmo caminho que `ccam <name>`, então as mesmas verificações se aplicam, incluindo o aviso quando outra sessão do Claude Code que compartilha o mesmo histórico já está em execução nesta pasta.
- `last-used` só é gravado quando o Claude Code realmente inicia. `--dry-run` não mexe nele, então uma simulação não avança a rotação.
- `--fallback` retém a saída padrão até uma tentativa terminar, então uma saída em streaming como `--output-format stream-json` chega de uma vez no final. Um cliente que envia entrada `stream-json` em streaming e espera respostas antes de enviar mais não pode usar `--fallback`. Arquivos que a primeira tentativa alterou antes de atingir o limite continuam alterados; a tentativa seguinte começa o prompt de novo.

## Relacionados

- [Limites de uso](/features/usage)
- [Limites de uso e troca de contas](/docs/usage-limits)
- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [O que fazer no limite de uso do Claude Code](/guides/usage-limit-workflow)
