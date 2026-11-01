---
title: Iniciar o Claude Code
description: Inicie o Claude Code com qualquer conta via ccam <name|number>, ccam run ou o painel; os argumentos chegam intactos ao claude e --dry-run mostra o comando.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: "11df08608a1e"
---

`ccam work` inicia o Claude Code com a conta `work`, e tudo o que vem depois do nome vai para o `claude` sem alterações. Use-o onde você digitaria `claude`: em um terminal, em um script ou com uma entrada enviada por pipe.

## Uso

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Informe o nome de uma conta ou o número dela em `ccam list`. O ccam não lê os argumentos depois do nome, então `ccam work --help` mostra a ajuda do Claude Code. Um `--` depois do nome é opcional. A entrada padrão, o código de saída e os sinais se comportam como no `claude`.

Quando o primeiro argumento começa com `-` e não é uma das opções do próprio ccam, o ccam escolhe a conta para você e passa todos os argumentos ao Claude Code. Ele usa a conta definida para a pasta atual com [`ccam folder`](/features/folder) ou, se não houver, a conta padrão: a definida com `ccam use` ou, se você não definiu nenhuma, a conta com o menor número.

### O painel

`ccam` sem argumentos abre o [painel](/features/dashboard) quando a entrada e a saída padrão são um terminal e, nos outros casos, imprime `ccam list`. Vá até uma conta (com as setas ou digitando o número dela) e pressione `Enter` para iniciá-la; `n` inicia a conta com mais espaço e `o` oferece `--continue`, `--resume` ou outros argumentos. O painel fecha antes de o Claude Code iniciar, então o terminal é restaurado primeiro.

**Assentos dos planos Team e Enterprise.** A escolha automática de conta (`ccam next`, `--fallback`, `ccam pick`, a tecla `n` do painel e, depois que você o ativa, o [`auto_switch`](/features/next#auto-switch)) também alterna entre assentos da mesma organização Claude Team ou Enterprise, o que pode violar os termos da Anthropic. Leia as [perguntas frequentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e consulte a sua organização e a Anthropic antes de contar com ela.

Quando o shell tem `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` ou `ANTHROPIC_AUTH_TOKEN` definida, o Claude Code pode usar essa chave ou esse token em vez do login da própria conta, então o ccam avisa você na saída de erro padrão antes de iniciar (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Ele só lê as variáveis; uma [conta de chave de API](/features/api) as limpa.

Por padrão, o ccam se substitui pelo Claude Code, então nada mais continua rodando. Com `relay = true` (ou [`auto_switch = true`](/features/next#auto-switch) no Linux e no macOS) no início do `config.toml`, acima do primeiro `[[group]]`, `[[profile]]` ou `[[folder]]`, o ccam, em vez disso, inicia as sessões interativas como processo filho e continua rodando ao lado delas. Isso permite que o `/ccam:switch`, ou o `auto_switch` em um limite de uso, continue a conversa em outra conta sozinho, com as opções, o modo de permissão, o nível de esforço e o agente da sessão; antes que o `/ccam:switch` a mova para uma conta de outra organização, o ccam pergunta a você ([detalhes](/features/next#relay)).

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--dry-run` | Imprime o ambiente e o comando que o ccam executaria e sai | desativado |
| `--fallback` | Para execuções `-p`/`--print`: quando a conta atinge o limite de uso, executa o mesmo comando de novo com a conta logada que tem mais espaço ([detalhes](/features/next#run-again-at-the-limit)) | desativado |
| `-v`, `--version` | Imprime `ccam <version>` (só como primeiro argumento) | |
| `-h`, `--help` | A ajuda do ccam, como `ccam --help` ou `ccam run --help` | |

Regras dos argumentos:

- O ccam só lê as próprias opções antes do nome da conta. Depois do nome, todos os argumentos pertencem ao Claude Code, incluindo `--dry-run` e `--help`.
- `--dry-run` funciona como `ccam --dry-run work`, `ccam run --dry-run work` e `ccam --dry-run` (para a conta da pasta ou a padrão). `--dry-run` e `--fallback` podem vir em qualquer ordem.
- `ccam run <name|number> [claude args]` é a forma longa de `ccam <name>`. `ccam run` sem nome usa a conta da pasta ou a padrão.
- Os nomes de conta não podem ser nomes de comandos do ccam, como `add`, `list`, `run` ou `usage`, então `ccam <word>` nunca é ambíguo.

A simulação imprime uma linha que você poderia colar em um shell, ou no PowerShell no Windows. Ela começa com o que a inicialização muda no ambiente: `CLAUDE_CONFIG_DIR` e `CCAM_PROFILE` sempre, `env -u <name>` para uma variável que ela descarta (como um `CCAM_RELAY` herdado) e, para uma [conta com chave de API](/features/api), o endpoint, o modelo e de onde vem a chave (`ANTHROPIC_API_KEY="$WORK_KEY"`), nunca a própria chave:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

A linha do PowerShell restaura as variáveis depois, como `VAR=x cmd` faz em um shell, e entrega ao Claude Code a linha de comando por meio de `Start-Process`, para que o Windows PowerShell 5.1 e o PowerShell 7 a executem do mesmo jeito. Quando uma inicialização define uma variável com um valor vazio, uma linha `#` abaixo avisa que o PowerShell anterior ao 7.5 remove uma variável assim. Para um `claude.cmd` instalado pelo npm, a linha executa o comando `cmd.exe` que a inicialização usaria (veja abaixo), e um argumento que a inicialização recusaria ali também faz a simulação falhar.

## Como funciona

Cada inicialização segue as mesmas etapas, da linha de comando até o `exec`:

<!-- motion:launch -->

Cada caixa acende quando a etapa dela é executada, e as linhas abaixo são o ambiente e os argumentos com que o Claude Code inicia.

1. O ccam começa em `app.Main`, que o `cmd/ccam/main.go` chama. Se o ccam foi iniciado com outro nome, é um [atalho](/features/shortcuts) e vai direto para o inicializador.
2. O ccam examina `os.Args` antes de o cobra entrar em ação (`app.parse`). Só o `ccam` sozinho (o painel), um nome de comando do ccam (`config.Reserved`), `-h`/`--help` e `-v`/`--version` vão para `cli.Execute`. Todo o resto é uma inicialização.
3. O ccam lê `~/.config/ccam/config.toml`, o único arquivo que ele analisa (`launcher.Launch`).
4. O ccam encontra a conta pelo nome ou número, descarta um `--` inicial e localiza o executável `claude` (`launcher.Prepare`, `ResolveClaude`). Sem nome, ele usa a conta da pasta, comparando o diretório de trabalho com as entradas `[[folder]]` do `config.toml` como strings, sem diferenciar maiúsculas e minúsculas no macOS e no Windows (`HereProfile`). A busca no PATH fica em cache em `~/.local/state/ccam/claude-path`, junto com o horário de modificação do arquivo e um hash do `PATH`, para que o ccam não percorra o `PATH` a cada inicialização. Um arquivo alterado ou ausente, ou um `PATH` alterado, provoca uma nova busca. Um caminho `claude_bin` no `config.toml` (`~/` é expandido) precisa ser um arquivo executável; caso contrário, a inicialização para com `claude_bin X is not an executable file. Check claude_bin in config.toml` antes de qualquer coisa ser gravada. O ccam recusa um `claude` que seja o próprio ccam, comparado pelo inode, porque isso entraria em loop.
5. O ccam compara a pasta da conta `~/.local/share/ccam/profiles/<name>/` com as configurações de compartilhamento dela, chamando só `lstat` e `readlink`, uma ou duas chamadas de sistema por item (`checkLinks`, `share.Check`). Um link ausente é criado na hora. Qualquer coisa que possa perder dados interrompe a inicialização.
6. O ccam adiciona os servidores MCP compartilhados (`applyMCP`; veja abaixo).
7. O ccam registra a sessão (`session.Register`). Ele grava um marcador com o nome do PID dele, contendo o horário de início do processo, em `~/.local/state/ccam/locks/profile-<name>/<pid>`, não importa o que diga `XDG_RUNTIME_DIR`, para que todos os shells vejam os mesmos marcadores. Para uma conta que compartilha o histórico de conversas, o ccam primeiro procura uma sessão em execução em `history-<group>-<hash>`, em que o hash é o da pasta atual com os links simbólicos resolvidos, e também se registra ali (`historyCheck`).

   O `exec` mantém o PID, então o marcador passa a indicar o próprio processo do Claude Code. O ccam usa um arquivo de marcador, e não um bloqueio mantido, porque os processos iniciados pelo Claude Code herdariam o bloqueio e o manteriam depois do fim da sessão. Um marcador que não pode ser gravado não impede a inicialização: o ccam diz `ccam: could not mark the session as running (…); starting Claude Code anyway`.

   As outras verificações do ccam sobre se uma conta está em execução (renomear, remover, alterar grupos, `doctor --fix`, a sincronização de MCP) também leem os registros do próprio Claude Code na pasta da conta (`sessions/<pid>.json`, `daemon.lock`), então também enxergam o Claude Code iniciado sem o ccam: `claude` simples, uma IDE, sessões em segundo plano e o daemon delas.
8. O ccam define o título do terminal como `claude · <name>` (quando a saída de erro padrão é um terminal), grava o horário da inicialização em `~/.local/state/ccam/last-used/<name>` sem `fsync` e chama `syscall.Exec` (`Plan.Exec`).

O Claude Code é executado com `CLAUDE_CONFIG_DIR` apontando para a pasta da conta e `CCAM_PROFILE` com o nome da conta. O ccam limpa qualquer valor herdado dessas variáveis, e um `CCAM_RELAY` herdado. Para uma pasta `~/.claude` registrada com `ccam import` no macOS, ele também define `CLAUDE_SECURESTORAGE_CONFIG_DIR` com um valor vazio, para que a conta e o `claude` simples usem o mesmo login do Acesso às Chaves; no Windows, ele faz o mesmo quando o Claude Code guarda o login no Gerenciador de Credenciais.

### Por que exec

O ccam substitui o próprio processo pelo `claude` (`syscall.Exec`). Com isso:

- o ccam não fica na memória;
- as teclas, os sinais e o código de saída são os do próprio Claude Code;
- os argumentos continuam sendo um array, nunca uma string de shell.

O Windows não tem `exec`. Lá, o ccam inicia o `claude.exe` como processo filho com o mesmo ambiente, os mesmos argumentos e o mesmo console, espera por ele e sai com o código de saída dele (`Plan.Exec`). Um `claude.cmd` que o npm instalou roda no `cmd.exe` com o caminho completo e cada argumento entre aspas, para que o `cmd.exe` não leia nenhum deles como sintaxe própria (`&`, `%PATH%`); um argumento com quebra de linha não pode ser passado assim e é recusado, com a recomendação de instalar o Claude Code nativo.

O Ctrl+C chega ao Claude Code direto pelo console, então o ccam o ignora. O marcador de sessão leva o PID do ccam, que vive exatamente o mesmo tempo que o Claude Code, e é removido quando o Claude Code termina. Enquanto espera, o ccam não faz mais nada.

O trabalho posterior a uma sessão, como notar servidores adicionados com `claude mcp add`, acontece na próxima inicialização: o ccam compara o horário de modificação do `.claude.json` com `~/.local/state/ccam/mcp/<name>/seen` e só analisa o arquivo quando ele mudou.

### Como uma linha de comando é encaminhada

Três linhas de comando, três caminhos por `internal/app`:

<!-- motion:launch-route -->

Um atalho é reconhecido pelo `argv[0]`, `ccam <name>` pelo primeiro argumento, e os dois terminam em `syscall.Exec`. Só um comando do ccam, como `ccam usage`, chega ao cobra.

### Servidores MCP compartilhados

Com a estratégia padrão `flag`, o ccam passa os servidores MCP do grupo ao Claude Code como um argumento, colocado antes dos seus:

<!-- motion:launch-mcp -->

O ccam só acrescenta o argumento se o `mcp.json` do grupo existir, e o coloca logo depois de `claude` (`mcp.FlagArgs`, `Plan.InsertArgs`). Suas próprias opções `--mcp-config` ou `--strict-mcp-config` continuam valendo por cima. É um único argumento, `--mcp-config=<path>`, porque o `--mcp-config` do Claude Code aceita vários valores: escrito como dois argumentos, ele também leria um prompt seguinte como arquivo de configuração.

O ccam decide onde a opção não pode ir (`claudeVerb`). Ele lê o subcomando do Claude Code nos argumentos do jeito que o Claude Code faz, pulando opções junto com os valores delas. `remote-control` ou `rc` não recebe a opção (abaixo), nem um verbo de sessão em segundo plano passado primeiro (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`). `ccam work mcp list`, `get` e `remove` também não a recebem, já que o Claude Code a recusa ali; em vez disso, o ccam diz na saída de erro padrão que os servidores do grupo também carregam nas sessões da conta:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Todos os outros subcomandos e todas as sessões mantêm a opção.

Com `mcp_strategy = "sync"` no `config.toml`, o ccam grava os servidores no `.claude.json` da conta e pula essa etapa se a conta já estiver em execução. Uma inicialização também pula esse trabalho quando nem o `mcp.json` do grupo nem o `.claude.json` mudaram desde a última verificação. Se a conta tiver um servidor próprio com o mesmo nome de um servidor do grupo, a sincronização mantém o servidor da conta e avisa uma vez. A sincronização também coloca de volta uma cópia que sumiu do `.claude.json` (`put back group default's MCP servers …`) e só altera ou remove as cópias que ela gravou e que ninguém alterou desde então.

### Remote Control

`ccam work rc` (ou `ccam work remote-control`) inicia o Remote Control do Claude Code como `work`. `claude remote-control` recusa `--mcp-config` antes do seu verbo e não aceita nenhum depois dele, e as sessões dele leem os servidores MCP só do `.claude.json` da conta. Por isso, a estratégia flag funciona de outro jeito aqui:

<!-- motion:launch-rc -->

1. O ccam encontra `remote-control` ou `rc` como verbo do Claude Code, mesmo depois de opções e dos valores delas (`ccam work --model opus rc`), e não acrescenta a opção (`claudeVerb`, `mcp.FlagArgs`).
2. O ccam compara o horário de modificação do `mcp.json` do grupo com o `sync-stamp` da conta e verifica se cada servidor do grupo ainda está no `.claude.json` (`mcp.Synced`). Se o `.claude.json` já tiver os servidores atuais do grupo, ou se o grupo não tiver nenhum, o Remote Control inicia sem perguntar.
3. Caso contrário, o ccam lista os servidores que faltariam às sessões e pergunta:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. Com `y`, o ccam os copia como a estratégia `sync` faz (`mcp.Sync`): com o bloqueio do Claude Code, depois de um backup, mantendo os servidores próprios da conta. Com `N`, ou sem um terminal, o Remote Control inicia sem eles, e o ccam indica `mcp_strategy = "sync"`, que os copia antes de cada inicialização.

O ccam pergunta a cada inicialização do Remote Control até que os servidores sejam copiados. Depois de uma cópia, ele só pergunta de novo quando o `mcp.json` do grupo muda ou um servidor copiado some do `.claude.json`. Uma cópia no `.claude.json` prevalece sobre `--mcp-config` nas inicializações normais seguintes, então `ccam mcp add`, `remove` e `import` atualizam as cópias que outras contas têm ([Servidores MCP compartilhados](/features/mcp)).

### O limite de 10 ms

O orçamento do caminho de inicialização é de 10 ms no percentil 95; na prática, o ccam acrescenta cerca de 4 ms antes de o Claude Code iniciar. O caminho de inicialização nunca monta a árvore de comandos do cobra, não inicia o Bubble Tea nem abre conexões de rede, e `TestPackageInitIsCheap` falha se as funções `init()` dos pacotes levarem mais de 5 ms no total. `scripts/perf-check.sh` mede `ccam --dry-run work`, `ccam work` e um atalho com `hyperfine` e falha acima de 10 ms, e `BenchmarkLaunchPath` em `internal/launcher` executa as etapas em loop. Para ver o tempo de cada etapa, defina `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Segurança e casos especiais

- Um nome ou número desconhecido para com `no account "<name>"` e indica `ccam list`. Sem nenhuma conta, o ccam pede que você adicione uma com `ccam add <name>`.
- Se um caminho compartilhado virou um arquivo comum ou aponta para o lugar errado, o ccam para com `stopped because of share link problems:`, lista os problemas e pede que você execute `ccam doctor --fix`. Ele não move dados no caminho de inicialização.
- Se outra sessão do Claude Code que compartilha o histórico de conversas com esta conta (mesmo grupo) estiver em execução na mesma pasta, o ccam avisa que o histórico pode se misturar e pergunta `Run anyway? [y/N]`. A pasta é comparada com os links simbólicos resolvidos. Sem um terminal (entrada padrão vinda de `/dev/null` conta como nenhum), ou a partir de um job em segundo plano, ele mostra o aviso e continua. Uma nova execução `-p` começa uma conversa própria e não recebe aviso; `--version`, `--help` e subcomandos que não abrem conversa, como `mcp`, `auth` ou `plugin`, nem são verificados.
- `--dry-run` pula a sincronização de MCP, o aviso de servidor novo, a nota do `mcp list`, a verificação do histórico, o marcador de sessão, o título e o horário do último uso, e não faz o exec. Ele ainda verifica os links de compartilhamento e cria um link ausente, como uma inicialização normal faria.
- `ccam work rc` com `work` já em execução não mexe no `.claude.json`, avisa e inicia o Remote Control sem os servidores do grupo; execute-o de novo depois que a sessão terminar para copiá-los. Com `--dry-run`, ele imprime a mensagem, mas não pergunta nem copia.
- Com `mcp_strategy = "sync"`, uma sincronização que falha imprime `ccam: MCP sync failed (launching anyway): …`, e o Claude Code inicia com os servidores que já tinha.

## Relacionados

- [Atalhos](/features/shortcuts)
- [Usar o Remote Control com o ccam](/guides/remote-control)
- [Gerenciar contas](/features/accounts)
- [Referência de comandos](/docs/commands)
- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [CLAUDE_CONFIG_DIR explicado](/guides/claude-config-dir)
