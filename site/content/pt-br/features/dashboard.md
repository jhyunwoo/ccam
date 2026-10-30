---
title: O painel
description: O ccam em um terminal abre um painel em tela cheia com contas, grupos, regras de pasta, extensões e verificações; cada mudança mostra o plano antes.
updated: 2026-10-30
order: 15
translation: machine
sourceHash: "d65f1b5bb0c8"
---

`ccam` sem argumentos abre um painel em tela cheia quando a entrada e a saída padrão são um terminal; caso contrário, imprime `ccam list`. O painel mostra todas as contas com o uso e os grupos de cada uma e faz o que os comandos fazem: iniciar, adicionar, fazer login, alterar grupos, definir regras de pasta, listar sessões em segundo plano, verificar e reparar. Iniciar uma conta, ou se conectar a uma sessão em segundo plano, fecha o painel; adicionar uma conta ou fazer login entrega o terminal por um momento e depois volta ao painel.

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

A linha abaixo dos nomes das abas diz em uma frase para que serve a aba atual. Em um terminal mais largo, a aba Accounts também mostra o e-mail de cada conta. Em todas as abas, o fim da linha das abas mostra se o [`auto_switch`](/features/next#auto-switch) está ativado (Linux e macOS; no Windows, ele é ignorado). Quando ele está ativado, a aba Health lista do que ele ainda precisa: o Claude Code 2.1.214 ou posterior e o plugin do ccam ativado em todas as contas que fazem login com uma assinatura do Claude.

Toda vez que abre, o painel pergunta ao GitHub em segundo plano se existe uma versão mais recente do ccam; o primeiro quadro nunca espera por isso e, quando você está sem conexão, nada é mostrado. Quando há uma, uma linha abaixo do título a recomenda:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` pergunta primeiro e depois atualiza do mesmo jeito que o [`ccam update`](/features/update): a versão é conferida com seus checksums SHA-256 e substitui o ccam em um único passo. O painel em que você está continua sendo a versão antiga, então saia e inicie o `ccam` de novo. Quando a atualização não pode ser feita aqui, por exemplo com uma cópia do ccam instalada pelo pacote `.deb`, o painel mostra o mesmo conselho que o `ccam update` e mantém a oferta.

O primeiro quadro é desenhado só a partir do `config.toml`, em bem menos de 100 ms; o uso aparece à medida que cada conta responde. Cada limite tem a própria barra: a janela de 5 horas, a janela semanal de todos os modelos e a janela semanal própria do Fable; `-` significa que a conta não tem essa janela. Uma porcentagem fica amarela em 80% e vermelha em 100%. Se uma consulta falhar, os valores do cache aparecem esmaecidos. Uma janela que o Claude Code informou como cheia mostra 100% até reiniciar. Um token de acesso expirado mostra um `⚠ token expired` esmaecido (o Claude Code o renova na próxima inicialização), enquanto `⚠ login expired` significa que a conta precisa de `ccam login`.

## Abas

| Aba | O que mostra | Teclas |
| --- | --- | --- |
| Accounts | número, nome, e-mail, barras de 5 horas, semanal e do Fable, último uso, grupos; `✔` a conta padrão, `▸` a conta desta pasta, `•` em execução (inclusive o Claude Code iniciado sem o ccam: `claude` simples, uma IDE, sessões em segundo plano e o daemon delas) | abaixo |
| Agents | o que o [`ccam agents`](/features/agents) lista: as sessões em segundo plano e as sessões abertas de todas as contas, com estado, pasta e horário de início (abaixo) | `enter` conectar (como `ccam attach`), `x` parar, `r` recarregar |
| Groups | cada grupo, o que ele compartilha e suas contas, e as contas que não estão em nenhum grupo | `enter` editar, `n` novo, `e` renomear, `x` remover |
| Folders | as regras do [`ccam folder`](/features/folder) e a que vale aqui, mesmo quando você chegou à pasta por um link simbólico (regras novas são guardadas com os links resolvidos) | `n` usar uma conta aqui, `enter` alterar, `x` remover |
| Extensions | o que o [`ccam ext`](/features/ext) mostra, cada extensão uma vez, com os grupos de onde ela vem, `off:` as contas em que um plugin não está ativado, os servidores próprios da conta que valem só para esta pasta, `This folder's project` e os problemas que não pertencem a nenhuma conta | `r` recarregar |
| Health | o que o [`ccam doctor`](/features/doctor) relata, com uma caixa por conta (abaixo) | `f` reparar a caixa selecionada, `F` reparar tudo (`doctor --fix`), `↑`/`↓` passar de uma caixa para outra |

`←`/`→` ou `Tab` trocam de aba, `?` lista as teclas da aba atual, `r` atualiza o uso, ignorando o cache curto, e carrega de novo as abas Agents, Extensions e Health, e `q` ou `Esc` sai.

## Teclas na aba Accounts

| Tecla | O que faz |
| --- | --- |
| `Enter` | Inicia o Claude Code como a conta selecionada |
| `1`–`9` | Vai para a conta com esse número; digite `1` `2` rapidamente para ir à 12. Só `Enter` inicia |
| `↑`/`↓`, `k`/`j` | Move o cursor |
| `n` | Inicia a conta que o [`ccam next --auto`](/features/next) escolheria |
| `o` | Inicia com opções: `--continue`, `--resume`, outros argumentos ou uma simulação |
| `u` | Detalhes do uso: as janelas de 5 horas, semanal e por modelo, com os reinícios e a origem dos números |
| `a` | Adiciona uma conta (o assistente abaixo) |
| `l` | Faz login novamente |
| `s` | Escolhe os grupos da conta |
| `d` | Torna a conta padrão |
| `f` | Usa a conta na pasta atual |
| `e` | Renomeia a conta |
| `c` | Define ou remove o comando de atalho dela |
| `x` | Remove a conta (pergunta antes, com No selecionado) |

**Assentos dos planos Team e Enterprise.** A escolha automática de conta (`ccam next`, `--fallback`, `ccam pick`, a tecla `n` do painel e, depois que você o ativa, o [`auto_switch`](/features/next#auto-switch)) também alterna entre assentos da mesma organização Claude Team ou Enterprise, o que pode violar os termos da Anthropic. Leia as [perguntas frequentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e consulte a sua organização e a Anthropic antes de contar com ela.

## A aba Health

A aba Health executa as verificações do [`ccam doctor`](/features/doctor) quando você a abre e desenha a mesma visão: uma primeira linha com uma marca por conta e depois uma caixa por conta. Cada linha de uma caixa é uma área verificada, primeiro `Folder`, `Login` e `Share links`, marcada com `✔` OK, `ℹ` informação, `⚠` um aviso ou `✘` um erro, para você ver de imediato qual conta tem um problema e onde. Os resultados que não pertencem a nenhuma conta vão para uma caixa `General`.

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

`↑`/`↓` passam de uma caixa para outra, e a caixa selecionada tem a borda na cor de destaque. `f` repara só a caixa selecionada, como o `ccam doctor --fix` faria: na caixa de uma conta, repara só essa conta, e na caixa `General`, os resultados que não pertencem a nenhuma conta (caminhos de plugins e atalhos); as demais contas são apenas verificadas de novo. Ele pergunta antes (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), e uma caixa sem nada para reparar só mostra `nothing to repair in <name>`. `F` repara tudo de uma vez, como o `ccam doctor --fix`, e `r` verifica de novo.

## A aba Agents

A aba Agents mostra o que o [`ccam agents`](/features/agents) imprime: as sessões em segundo plano e as sessões abertas de todas as contas, com a conta a que cada uma pertence. Ela consulta o Claude Code de todas as contas ao mesmo tempo (`claude agents --json`, cada uma com um tempo limite de 15 segundos) quando você abre a aba pela primeira vez, ao pressionar `r`, depois de parar uma sessão e depois de uma mudança na configuração, e mostra `Asking every account's Claude Code for its sessions…` até que todas tenham respondido.

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

`State` é o estado de uma sessão em segundo plano e o que o processo dela está fazendo, ou `open` para uma sessão interativa, como no `ccam agents`. Se uma conta não puder ser consultada, uma linha `✘` abaixo da lista mostra o que o Claude Code disse, como `lab` acima, que está com a visualização de agentes desativada. As sessões em segundo plano já encerradas ficam de fora (`ccam agents --all` as lista), e, quando todas as contas responderam e nenhuma tem sessões, a aba diz ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` Em um terminal estreito, a aba remove `Started`, depois `Name`, depois `Folder`, e corta caminhos de pasta longos pela esquerda.

`Enter` em uma sessão em segundo plano fecha o painel e abre a sessão neste terminal como a conta dela, como faz o `ccam attach <id>`; o painel não volta (no Claude Code, `←` vai para a visualização de agentes e `Ctrl+Z` para o shell). Uma sessão aberta não tem ID, então `Enter` e `x` só dizem que é preciso uma sessão em segundo plano. `x` pergunta antes, com `No, keep it` selecionado: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` Se a resposta for sim, o ccam executa `claude stop <id>` como essa conta, diz ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` ou por que não conseguiu, e lê a lista de novo.

## Alterar grupos

Na aba Groups, `Enter` abre o editor: os itens à esquerda, as contas à direita, `Tab` para alternar entre eles e `Space` para marcar ou desmarcar. Um item que uma das contas já recebe de outro grupo aparece esmaecido, com o motivo, assim como uma conta que receberia um item duas vezes; um item só pode vir de um grupo. `Enter` mostra o plano:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Se um arquivo difere da cópia do grupo, o ccam primeiro pergunta o que fazer (manter o do grupo, usar o meu ou mesclar). Depois, a alteração roda em segundo plano com o mesmo código do [`ccam group`](/features/share): a configuração é salva primeiro, os dados de cada conta são movidos em ordem, e os originais vão para `.ccam-backup/`; uma conta que sai do compartilhamento de MCP mantém os servidores do grupo no `.claude.json` dela.

Uma alteração que moveria os dados de uma conta em execução é recusada, com o nome da conta, mesmo quando a conta foi iniciada enquanto o plano estava aberto, e um plano que ficou desatualizado porque outro processo do ccam alterou os grupos é recusado com `the groups changed while this plan waited; nothing was moved, make the change again`. Enquanto outro processo do ccam altera os grupos, uma nova alteração, renomeação ou remoção diz `another ccam process is changing the groups right now`; depois de uma alteração que parou no meio, elas esperam até que `F` na aba Health (ou `ccam doctor --fix`) a conclua.

## Adicionar uma conta

`a` abre o assistente que o `ccam add` também usa em um terminal:

1. O nome da conta.
2. Como ela faz login: com uma assinatura do Claude agora (no navegador), com uma [chave de API](/features/api) ou com uma assinatura do Claude depois.
3. Para uma chave de API: de onde o ccam lê a chave a cada inicialização (`env:NAME` ou `command:COMMAND`), o endpoint, como a chave é enviada e o modelo.
4. Os grupos em que a conta entra, com `default` marcado (se ainda não houver grupos, o que o grupo `default` vai compartilhar).
5. Um comando de atalho, se você quiser.
6. Um resumo para confirmar.

`Esc` volta uma etapa e, na primeira etapa, fecha o assistente.

## Remover uma conta

`x` pergunta antes de qualquer coisa acontecer, com No selecionado, e diz o que a remoção faz: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` Para uma conta com chave de API, não há login a revogar, e, quando outra conta tem o mesmo login, ela não é desconectada. Uma conta importada mantém a pasta: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## Cores

As cores seguem os papéis da interface de terminal do Claude Code, para que o ccam pareça em casa ao lado dele; os nomes de grupos recebem uma cor derivada do próprio nome. O ccam pergunta ao terminal qual é a cor de fundo e escolhe um tema claro ou escuro. Para escolher um você mesmo, defina-o no início do `config.toml`, acima do primeiro `[[group]]`, `[[profile]]` ou `[[folder]]`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Os temas `-daltonized` são para daltonismo vermelho-verde, e os temas `-ansi` usam as 16 cores do próprio terminal. `NO_COLOR` desativa as cores.

## Segurança e casos especiais

- O painel chama as mesmas funções que os comandos, então todas as regras valem: planos antes de mover dados, recusas para contas em execução, backups e tokens de login que nunca saem de uma conta.
- O painel não tem simulação: `ccam --dry-run` sem um nome imprime como a conta desta pasta seria iniciada, e `o` em uma conta oferece uma simulação só daquela inicialização.
- O terminal é restaurado antes de o Claude Code iniciar ou de o login abrir.
- Em uma tela pequena demais para um diálogo, as linhas da aba são as primeiras a sair; as linhas são cortadas com `…` em vez de quebradas. Em uma tela baixa, as linhas em branco e a descrição da aba são removidas em seguida, para que a linha de dicas de teclas continue visível. Abaixo de 40 colunas ou 8 linhas, o painel mostra só `Too small: ccam needs 40x8 (now 36x7)`.
- `Ctrl+C` enquanto uma alteração, uma remoção, uma atualização ou um reparo está sendo aplicado espera que termine: a linha de status diz `ccam quits when it is done.`
- Depois que `U` atualizou o ccam, a oferta não volta quando o painel retorna de adicionar uma conta ou de um login.

## Relacionados

- [Compartilhamento com grupos](/features/share)
- [Iniciar o Claude Code](/features/launch)
- [Gerenciar contas](/features/accounts)
- [Limites de uso](/features/usage)
