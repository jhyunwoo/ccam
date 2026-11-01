---
title: Limites de uso
description: ccam usage mostra os limites de 5 horas e semanais de todas as contas de uma vez, e ccam statusline coloca o nome da conta e o uso no Claude Code.
updated: 2026-11-01
order: 7
translation: machine
sourceHash: "54b296f7da35"
---

`ccam usage` mostra, para cada conta, quanto de cada limite já foi usado: o limite de 5 horas, o limite semanal e o limite semanal próprio de cada modelo (como o Fable). Cada limite ganha uma barra e o tempo até reiniciar. Use-o antes de começar uma tarefa longa, ou quando o Claude Code avisar que um limite está próximo, para ver qual conta ainda tem espaço. `ccam statusline` é o comando de linha de status do Claude Code: ele mostra o nome da conta e o uso dentro do Claude Code e salva os números para o `ccam usage`. Para os tokens e os custos estimados por trás do uso, o `ccam ccusage` executa o [ccusage](https://ccusage.com) nas conversas de todas as contas ([veja abaixo](#tokens-and-costs-ccam-ccusage)).

## Uso

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

Cada conta ganha uma linha com o número, o nome e a idade dos seus números; `(cached)` significa que eles vêm de um arquivo salvo, não do endpoint. Abaixo dela, uma linha por limite mostra uma barra de 20 células, a porcentagem e o tempo até a janela reiniciar: `5-hour`, `Weekly` (todos os modelos) e depois a janela semanal própria de cada modelo, como `Fable`. Uma janela que a conta não tem mostra `-`. A porcentagem fica amarela em 80% e vermelha em 100%.

Se o Claude Code recusou uma solicitação porque uma janela estava cheia, essa janela mostra 100% até reiniciar, não importa o que digam os números salvos (veja [o registro de limite](#the-limit-record)). Em um terminal, todas as contas aparecem imediatamente como `fetching…`, e as linhas de cada uma são desenhadas quando ela responde.

A última linha, mostrada só quando você tem mais de uma conta, indica a conta que o [`ccam next --auto`](/features/next) escolheria. Ela segue as mesmas regras:

- `auto_max_usage` e `auto_strategy`;
- as sessões em execução;
- só contas da organização da última conta usada, a menos que `cross_org` esteja definido;
- toda janela em 100% é pulada, então uma conta com a semana esgotada nunca é indicada.

Com a estratégia padrão `headroom`, ela diz `→ Most headroom: <name>`; com `expiring` ou `balance`, ela dá o motivo, como `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Assentos dos planos Team e Enterprise.** A escolha automática de conta (`ccam next`, `--fallback`, `ccam pick`, a tecla `n` do painel e, depois que você o ativa, o [`auto_switch`](/features/next#auto-switch)) também alterna entre assentos da mesma organização Claude Team ou Enterprise, o que pode violar os termos da Anthropic. Leia as [perguntas frequentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e consulte a sua organização e a Anthropic antes de contar com ela.

Uma conta sem login mostra `not logged in → ccam login <name>`, e uma conta com chave de API mostra `no limits (API)`. Um token de acesso expirado mostra ``access token expired; `ccam <name>` renews it when it starts``, acima dos números salvos, se houver: o Claude Code o renova na próxima inicialização. Um login que expirou por completo, porque o token de renovação expirou ou porque o Claude Code o apagou depois que o servidor o recusou, mostra `` login expired → `ccam login <name>` ``.

Depois de `ccam statusline install`, o Claude Code mostra uma linha como `[work] · Opus · 5h 24% · 7d 41%`.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--json` | Imprime os relatórios em JSON: `name`, `slot`, `status`, `limits` (com `credits` quando a conta tem créditos de uso ativados), `source`, `at`, `error`, e `limited` e `limited_windows` quando um registro de limite põe janelas em 100%. Sem contas, imprime `[]` | desativado |
| `--watch[=N]` | Redesenha a cada `N` segundos até você pressionar Ctrl+C. `--watch` sozinho, ou um valor menor ou igual a 0, significa 60. O valor vai depois de `=`: `--watch 10` é recusado com uma dica | desativado |
| `-v`, `--verbose` | Adiciona uma linha `source:` abaixo de cada conta: a origem dos números dela (o endpoint de uso, o Claude Code, a linha de status, um arquivo salvo ou o limite de uso que o Claude Code informou). Quando um registro de limite prevalece sobre os números, segue uma linha como `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` | desativado |
| `--no-cache` | Ignora o cache de 60 segundos e consulta o endpoint de novo | desativado |
| `--source live\|claude\|cache` | `live` consulta o endpoint com o token da conta; `claude` faz o Claude Code consultar (veja abaixo); `cache` mostra apenas os números salvos | `usage_source`, senão `live` |
| `--refresh` | Renova os tokens de acesso expirados, só das contas que não estão em execução. Uma renovação pode levar até 30 segundos | desativado |
| `--dry-run` | Aceito: o `ccam usage` não altera nada além do próprio cache. Recusado junto com `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | desativado |

Qualquer outro valor de `--source` falha com `--source must be live, claude or cache`. `ccam usage` não aceita argumentos.

### Uso consultado pelo próprio Claude Code

Com `usage_source = "claude"` no início do `config.toml` (acima do primeiro `[[group]]`, `[[profile]]` ou `[[folder]]`), ou `--source claude`, o ccam não lê o token de login em momento algum. Para cada conta, ele inicia `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` na pasta da conta e envia a requisição `get_usage` do Claude Code pela entrada padrão (Claude Code 2.1.270 ou posterior). O Claude Code consulta o uso com o próprio login e responde com as janelas de 5 horas, semanal e por modelo.

`--safe-mode` impede que os hooks, os servidores MCP e os plugins da conta iniciem, e nenhuma sessão é salva. Cada conta exige uma inicialização do Claude Code (cerca de um segundo), então essa fonte é mais lenta que `live`. As contas são consultadas ao mesmo tempo, cada uma com um tempo limite de 20 segundos, e o cache de 60 segundos vale como sempre.

Um login cujo token de acesso expirou, ou expira em 5 minutos, recebe 50 segundos, já que o Claude Code o renova ao iniciar e uma renovação interrompida pode perder o login. Se o ccam desiste de uma consulta, ele pede ao Claude Code que saia com SIGTERM, para que o Claude Code conclua uma renovação que começou (o Windows não tem SIGTERM, então lá o processo é encerrado).

Um login que expirou por completo mostra `` login expired → `ccam login <name>` `` sem iniciar o Claude Code. O Claude Code marca a requisição como experimental, então, se ela deixar de responder, o ccam mostra os números salvos e o `ccam doctor` informa o problema. Quando o Claude Code responde que a própria consulta falhou, o ccam mostra os números salvos com `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`, e não um plano sem limites.

Com `--dry-run`, o ccam não inicia o Claude Code, porque o Claude Code grava na pasta da conta; em vez disso, mostra os números salvos. `usage_source` aceita `live` ou `claude`; `cache` é só uma opção de `--source`.

`ccam statusline` lê da entrada padrão o JSON da linha de status do Claude Code. A única opção dele é `--wrap <command>`: um comando de linha de status existente, cuja primeira linha de saída recebe o nome da conta na frente. `ccam statusline install` não aceita opções. Os dois recusam `--dry-run` com `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Tokens e custos: `ccam ccusage`

O `ccam usage` mostra quanto de cada limite uma conta já usou. Já o [ccusage](https://ccusage.com) conta os tokens e os custos estimados nos registros de conversa do Claude Code. Sozinho, ele procura os registros do Claude Code só em `~/.claude` e `~/.config/claude` (ou nas pastas que `CLAUDE_CONFIG_DIR` indica), então a única conta do ccam que ele consegue ver é uma importada no lugar em `~/.claude`. O `ccam ccusage` o executa com `CLAUDE_CONFIG_DIR` definido para as pastas de todas as contas que têm conversas, então um único relatório cobre todas elas:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

O ccam lê primeiro as próprias opções: `-a` (`--account`) recebe nomes ou números de contas, separados por vírgula ou repetindo a opção, e `--dry-run` imprime o comando em vez de executá-lo. A partir do primeiro outro argumento, tudo vai para o ccusage do jeito que você digitou, exceto `--config` para `bun x` e `npx` (veja abaixo); `ccam ccusage -- --help` mostra a ajuda do próprio ccusage. O ccam executa o `ccusage` se ele estiver instalado; caso contrário, `bun x ccusage` e, se isso falhar, `npx ccusage@latest`; os dois últimos baixam o ccusage do npm. Eles começam na pasta do próprio ccam (`~/.local/state/ccam/ccusage`), não na pasta em que você está, então um projeto que esteja ali não consegue escolher o que eles executam com o `node_modules/.bin` dele, um `bin` no `package.json` dele ou um registry no `.npmrc` dele; `ccam --dry-run ccusage` mostra o `cd` para essa pasta. Quando nenhum deles está disponível, ele diz para você instalar o Node.js, o Bun ou o ccusage. No Linux e no macOS, o ccusage toma o lugar do ccam; no Windows, o ccam espera por ele e termina com o código de saída dele. No WSL, o ccam pula um `ccusage`, `bun` ou `npx` que esteja em uma unidade do Windows, já que ele rodaria como um programa do Windows sem as pastas das contas.

- O relatório cobre as contas cujas pastas têm conversas (uma pasta `projects`). Uma conversa que várias contas compartilham por meio de um grupo conta uma vez só.
- Com `-a`, uma conta que compartilha o histórico por meio de um grupo traz as conversas do grupo inteiro, e o ccam diz de quem são: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- `blocks` mostra janelas de 5 horas, e cada conta tem as próprias, então indique uma única conta para ele, como acima.
- O ccusage também relata as outras ferramentas de programação que encontra, como o Codex, a menos que você indique uma: `ccam ccusage claude daily` mostra só o Claude Code.
- Enquanto o ccam o executa, o ccusage lê as configurações dele (`ccusage.json`) só de `./.ccusage` e das pastas das contas. `--config` indica outro arquivo: `ccam ccusage daily --config ~/.config/claude/ccusage.json`. Como `bun x` e `npx` começam na pasta do próprio ccam, o ccam passa a eles `./.ccusage/ccusage.json` com `--config` quando esse arquivo existe e você não indica nenhum, e passa um `--config` relativo como caminho completo. O ccusage então lê outros caminhos relativos, como `--pi-path`, a partir da pasta do ccam; por isso, informe esses caminhos por completo.
- O ccusage leria uma pasta cujo caminho tem uma vírgula como se fossem duas pastas, então o ccam a recusa e indica a conta; deixe essa conta de fora com `-a`.
- `ccusage` é um comando do ccam, então uma conta chamada ccusage é iniciada com `ccam run ccusage`; o `ccam doctor` diz isso.

## Como funciona

Os números do `ccam usage` vêm de dois lugares: o endpoint de uso não documentado da Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, chamado com o login de cada conta, e os valores que o Claude Code passa ao comando de linha de status. Tudo isso fica em `internal/usage`.

1. O ccam carrega o `config.toml` e monta um `usage.Account` por conta (`cli.usageCmd`): a pasta dela `~/.local/share/ccam/profiles/<name>/`, os marcadores de sessão dela (com os registros do próprio Claude Code na pasta da conta, eles dizem se ela está em execução) e se ela foi importada no lugar.
2. O ccam consulta todas as contas ao mesmo tempo, com uma goroutine por conta (`Fetcher.FetchAll`). Todas compartilham um único `http.Client`, então as conexões com o endpoint são reaproveitadas. A conexão pode levar 2 s, e cada conta tem 3 s no total. O cliente não segue redirecionamentos, já que um redirecionamento reenviaria o token para outro lugar, e lê no máximo 1 MB da resposta.
3. À medida que cada conta termina, a linha dela é redesenhada. Uma conta que falha ou não responde a tempo passa a usar os números salvos, então uma conta lenta nunca atrasa as outras.

No diagrama, cada conta é consultada na própria goroutine, e a que estoura o tempo mostra os números salvos em vez disso:

<!-- motion:usage-fetch -->

`personal` e `work` respondem a tempo e mostram números atuais. `side` leva mais de 3 s, então o ccam mostra os números salvos dela em vez disso (`fallback()`). Depois, o ccam desenha um bloco por conta (`usage.Table()`) e escolhe a conta indicada na última linha (`usage.Best()`).

Para uma conta, o ccam faz estas verificações em ordem (`Fetcher.Fetch`):

<!-- motion:usage-cache -->

1. Com `--source cache`, vai direto para `fallback()`.
2. A menos que `--no-cache` esteja definido, um resultado atual em `~/.local/state/ccam/usage/<name>.json` medido há menos de 60 segundos (`CacheTTL`) é usado, avaliado em relação ao momento atual: uma janela que reiniciou desde então mostra 0%, e um arquivo com data à frente do relógio não conta como recente. Executar `ccam usage` duas vezes seguidas custa uma solicitação por conta.
3. O ccam lê o `.credentials.json` (ou o Acesso às Chaves do macOS) sem alterá-lo e decodifica só a chave `claudeAiOauth` (`readCredentials()`). Sem um login de assinatura, o resultado é `not logged in`, ou `no limits (API)` quando o `.claude.json` contém uma chave de API. Um login cujo token de renovação expirou, ou que o Claude Code apagou, é `` login expired → `ccam login <name>` `` e não é consultado.
4. Um token de acesso expirado não é usado. Sem `--refresh`, a conta mostra os números salvos com ``access token expired; `ccam <name>` renews it when it starts``.
5. Caso contrário, o ccam chama o endpoint. Um 401 ou 403 conta como login expirado (`login expired: authentication failed (401)`), qualquer outro status ou um corpo que ele não consegue ler conta como erro, e nos dois casos ele usa os números salvos. Uma consulta encerrada por Ctrl+C diz `cancelled`. Um 429 significa que o endpoint recebeu pedidos demais para esta conta: o ccam grava `<name>.pause` com o horário para perguntar de novo (`Retry-After`, ou 5 minutos, no máximo uma hora), mostra os números salvos com `rate limited by the usage endpoint; asking again after 15:10` e não chama o endpoint para essa conta até lá, nem com `--no-cache`. Tentar de novo imediatamente só prolongaria o limite de taxa.
6. Uma resposta válida é gravada em `<name>.json` com o horário em que foi medida.

Para os números salvos, o ccam pega a consulta salva mais recente (`<name>.json`, ou `<name>.claude.json` das consultas do Claude Code) e, quando o `<name>.statusline.json` foi gravado depois dela, sobrepõe a ela as janelas de 5 horas e semanal da linha de status (`fallback()`); as janelas por modelo, que a linha de status não tem, ficam como a consulta as salvou. Antes de mostrar números salvos, ele verifica o horário de reinício de cada janela: uma janela que reiniciou depois que os números foram salvos mostra 0%, porque a porcentagem antiga não vale mais. É por isso que `personal`, acima, mostra `0%` sem horário de reinício.

O endpoint tem dois formatos de resposta. O mais novo tem um array `limits` com as entradas `session` (5 horas), `weekly_all` e `weekly_scoped` (por modelo). O mais antigo tem os objetos `five_hour`, `seven_day`, `seven_day_opus` e `seven_day_sonnet`. O ccam lê os dois (`ParseLive`): o array `limits` prevalece, e os objetos antigos completam o que ele não cobre. Um corpo sem nenhum limite é informado como mudança de formato, e o `ccam doctor` avisa sobre isso.

### O registro de limite

Quando o Claude Code interrompe um turno em um limite de uso, o hook do [plugin do ccam](/docs/install#claude-code-plugin) lê na conversa a solicitação recusada: qual janela estava cheia (5 horas, semanal, Opus ou Sonnet) e quando ela reinicia. Ele grava isso em `~/.local/state/ccam/usage/<name>.limit.json`. O ccam aplica o registro por cima de tudo o que mais souber, de todas as origens (`Fetch`): cada janela nele conta como 100% até reiniciar, por mais antigos ou baixos que sejam os números salvos, então `ccam next --auto`, `ccam pick`, `--fallback` e o [`auto_switch`](/features/next#auto-switch) não escolhem a conta de novo antes disso.

Sem mais nada conhecido, o registro sozinho é o que a conta mostra (`source: the usage limit Claude Code reported` com `-v`). Sessões sem supervisão (execuções `-p` e sessões em segundo plano) não gravam registro, nem um bloqueio definido por uma organização (por exemplo, quando ela fica sem uso ou o limite de um grupo é $0), que os administradores dela podem suspender a qualquer momento. `ccam login` em outra conta Claude descarta o registro; excluir o arquivo também o apaga.

### Renovar tokens

Os tokens de renovação são rotativos: renovar um gera um novo token de renovação e invalida o antigo. Se o ccam renovasse o login enquanto o Claude Code estivesse em execução com a mesma conta, essa sessão seria desconectada. Por isso, o `ccam usage` lê os tokens e nunca os grava, a menos que você passe `--refresh`. Nesse caso, o ccam (`Fetcher.refresh`):

1. recusa as contas registradas no lugar com `ccam import`, já que o `claude` simples pode usar esse login sem o ccam saber;
2. recusa os logins armazenados no Acesso às Chaves do macOS;
3. obtém um arquivo de bloqueio, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, para que dois processos do ccam nunca renovem ambos o mesmo token; enquanto ele está ocupado, a conta mostra `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)`. `ccam rename`, `remove` e `logout` também o obtêm e se recusam enquanto uma renovação está em andamento;
4. com o bloqueio, recusa se uma sessão da conta estiver em execução, seja qual for a forma como foi iniciada (pelo ccam, pelo `claude` simples, por uma IDE ou por uma sessão em segundo plano e o daemon por trás dela, que renova o login por conta própria), e pula a solicitação se outro processo já tiver renovado o token;
5. não envia nada enquanto o Claude Code mantém os próprios bloqueios de renovação na pasta da conta (`.oauth_refresh.lock`, `<folder>.lock`): o Claude Code está renovando aquele login naquele momento;
6. envia o token de renovação ao endpoint de tokens e espera até 30 segundos pela resposta, o tempo que o Claude Code dá às próprias renovações. Depois de enviada, uma solicitação nunca é interrompida, nem por Ctrl+C: o servidor pode já ter trocado o token de renovação, e só a resposta dele traz o novo login. Ctrl+C ou um hangup encerra as outras consultas, e o ccam sai (com 130 depois de Ctrl+C) assim que a resposta é salva;
7. salva com o bloqueio de gravação do Claude Code (`.storage-write.lock`) e altera só os campos de token do `.credentials.json` como ele estiver nesse momento: o que o Claude Code gravou nesse meio-tempo (um login de MCP) permanece, um login que o Claude Code renovou nesse meio-tempo é mantido, e um logout nunca é desfeito. Quando a resposta informa a validade do novo token de renovação, ela é guardada como `refreshTokenExpiresAt`, como o Claude Code faz.

Uma renovação que não acontece diz o motivo, na linha da conta ou acima dos números salvos: `token not renewed: a session is running`, por exemplo, ou `login may be lost: …` quando uma solicitação foi enviada mas a resposta dela não pôde ser usada ou salva, com a recomendação de executar `ccam login <name>` se o Claude Code então pedir que você faça login.

### O coletor da linha de status

O Claude Code executa o comando de linha de status e envia a ele um documento JSON por pipe. O ccam lê dele o nome do modelo e `rate_limits`:

<!-- motion:usage-statusline -->

1. O ccam encontra a conta por `CCAM_PROFILE`, que o ccam define ao iniciar o Claude Code, ou então comparando `CLAUDE_CONFIG_DIR` com as pastas das contas (`runStatusline`).
2. O ccam lê `rate_limits.five_hour` e `rate_limits.seven_day` (`ParseStatusLine`). O JSON não tem limites de taxa para contas com chave de API, nem antes da primeira resposta da sessão; nesses casos, nada é salvo.
3. O ccam grava `~/.local/state/ccam/usage/<name>.statusline.json` (`SaveStatusLine`). O `ccam usage` o usa quando o endpoint não responde.
4. O ccam imprime `[name] · model · 5h N% · 7d N%`. Com `--wrap`, ele executa o seu comando com `sh -c` e a mesma entrada e imprime a saída dele do jeito que o Claude Code mostra uma linha de status: cada linha que não está em branco, com `[name]` na frente da primeira. Um comando que falha mostra só o nome da conta. No Windows, ele executa o comando como o Claude Code faz: com o Git Bash, se estiver instalado, senão com o PowerShell.
5. Quando o uso de 5 horas ou o semanal está em 90% ou mais e existe outra conta, a linha termina com `→ ccam next --auto --continue`, por exemplo `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

O comando encapsulado não tem limite de tempo próprio: o Claude Code dá 10 minutos a um comando de linha de status e interrompe uma execução quando uma atualização mais nova começa, então um comando lento aparece como apareceria sem o ccam. Quando o Claude Code encerra o ccam, o comando também termina (no Linux e no macOS, com todos os processos que ele iniciou).

`ccam statusline install` grava `statusLine` no `settings.json` de cada conta e mantém as outras chaves do objeto, como `padding`, `refreshInterval` e `hideVimModeIndicator`. Quando uma conta já tem um comando de linha de status, o ccam o mantém como `--wrap '<old command>'`; no Windows, como `--wrap-base64 <base64>`, que não precisa de aspas em nenhum dos dois shells. Se a linha de status já executa o ccam (de outro caminho, ou entre aspas do jeito que uma versão mais antiga a gravou), o ccam a substitui em vez de encapsulá-la e mantém o comando que ela encapsula.

No Windows, o caminho do ccam é gravado com barras normais e `~` para a pasta pessoal (`~/.local/bin/ccam.exe statusline`), que tanto o Git Bash quanto o PowerShell entendem. Um caminho que precisa de aspas é gravado por inteiro para o shell que o Claude Code usa: `'…'` para o Git Bash, `& '…'` para o PowerShell. Um `settings.json` compartilhado por um grupo é gravado uma única vez, pelo link, então todas as contas do grupo o recebem. Executá-lo de novo imprime `(already set)` para os arquivos que já usam o ccam.

## Segurança e casos especiais

- O `ccam usage` nunca altera um arquivo de token, a menos que você passe `--refresh`. Quando um token de acesso expirar, inicie a conta uma vez com `ccam <name>` e o Claude Code o renova. Um login que expirou por completo precisa de `ccam login <name>`.
- `--refresh` nunca mexe em contas em execução (mesmo quando o Claude Code foi iniciado sem o ccam), contas importadas nem logins do Acesso às Chaves. Uma renovação recusada mostra os números salvos e o motivo (`token not renewed: …`).
- O endpoint de uso não é documentado publicamente e pode mudar. Quando isso acontecer, você continua vendo os números salvos, e o `ccam doctor` informa a mudança de formato.
- Sem números salvos e sem resposta, a linha da conta mostra `lookup failed:` e o motivo, por exemplo `no saved usage` com `--source cache`.
- Os arquivos de cache guardam só porcentagens e horários de reinício, nunca tokens.

## Relacionados

- [Limites de uso e troca de contas](/docs/usage-limits)
- [Trocar para a próxima conta](/features/next)
- [O que fazer no limite de uso do Claude Code](/guides/usage-limit-workflow)
- [Referência de comandos](/docs/commands)
