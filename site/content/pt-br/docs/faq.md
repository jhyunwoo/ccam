---
title: Perguntas frequentes
metaTitle: "Perguntas frequentes: várias contas do Claude Code"
description: Respostas sobre o ccam, contas do Claude Code, limites de uso, compartilhamento, privacidade e sistemas suportados.
updated: 2026-11-19
order: 6
translation: machine
sourceHash: "090ecc3a34bd"
---

## O ccam é uma ferramenta oficial da Anthropic?

Não. O ccam é um projeto de código aberto não oficial e não é afiliado à Anthropic nem endossado por ela. Ele usa a configuração `CLAUDE_CONFIG_DIR`, que o [Claude Code documenta](https://code.claude.com/docs/en/env-vars) para executar várias contas lado a lado.

## Posso usar várias contas do Claude Code?

Isso depende dos termos de cada conta e, no caso de contas de trabalho, da política da sua organização: os [Termos do Consumidor](https://www.anthropic.com/legal/consumer-terms) da Anthropic para Pro e Max, os [Termos Comerciais](https://www.anthropic.com/legal/commercial-terms) dela para Team e Enterprise, e a [Política de Uso](https://www.anthropic.com/legal/aup) dela para todos. Verifique-os antes de usar várias contas. O ccam não aumenta, não junta nem esconde os limites de uma conta, e cada conta continua parando nos próprios limites. Ele move o seu trabalho para outra das suas contas quando você pede, ou sozinho em um limite quando você ativa o [`auto_switch`](/features/next#auto-switch) (Linux e macOS). Verifique se isso é permitido para as suas contas antes de usá-lo; a decisão é sua.

**Assentos dos planos Team e Enterprise.** A escolha automática de conta (`ccam next`, `--fallback`, `ccam pick`, [`ccam agent`](/features/agent), a tecla `n` do painel e o [`auto_switch`](/features/next#auto-switch), quando você o ativa) também alterna entre assentos da mesma organização Claude Team ou Enterprise. A central de ajuda da Anthropic diz "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada assento de um plano Team é atribuído a uma pessoa específica, e só essa pessoa pode usá-lo) e "Users can't route their work through another user's account to use that user's remaining usage" (os usuários não podem encaminhar o próprio trabalho pela conta de outro usuário para aproveitar o uso restante desse usuário) ([planos Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); o mesmo vale para os [planos Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Por isso, usar a escolha automática entre assentos pode violar os termos da Anthropic. Usá-la nesse caso é uma decisão sua: consulte a sua organização e a Anthropic antes de contar com ela. Os planos pessoais (Pro, Max) são tratados no parágrafo acima.

## O ccam contorna os limites de uso?

Não além do limite de nenhuma conta: cada conta mantém os próprios limites e para neles, e o ccam não os aumenta, não os junta nem os esconde. O ccam mostra esses limites e permite que você continue em outra das suas contas quando uma delas se esgota; com o `auto_switch` ativado, ele faz isso sozinho (veja a próxima pergunta). Entre assentos de uma mesma organização Team ou Enterprise, continuar em outro assento usa o uso restante desse assento, o que a central de ajuda da Anthropic não permite para o assento de outro usuário (veja Assentos dos planos Team e Enterprise acima).

## O ccam pode trocar de conta sozinho quando uma sessão atinge o limite?

Só se você ativar isso. Com `auto_switch = true` no início do `~/.config/ccam/config.toml` (Linux e macOS; desativado por padrão e experimental), quando o Claude Code interrompe uma sessão interativa (uma que o ccam iniciou) em um limite de 5 horas, semanal, do Opus, do Sonnet ou do Fable, o ccam encerra a sessão e continua a conversa como outra das suas contas com espaço: outro assento da mesma organização Team ou Enterprise, ou outro plano pessoal (Pro, Max). Ele nunca troca antes do limite, nunca acompanha o seu uso em segundo plano e nunca move trabalho para fora de uma organização Team ou Enterprise nem para dentro de uma.

Em assentos, ele move o seu trabalho de um assento para outro, o que pode violar os termos da Anthropic (veja Assentos dos planos Team e Enterprise acima). Para planos pessoais, os termos da Anthropic não tratam disso com clareza: os Termos do Consumidor proíbem acessar os Serviços "through automated or non-human means" (por meios automatizados ou não humanos), a menos que a Anthropic permita isso explicitamente, e uma troca envia uma mensagem em seu nome. Ativá-lo é uma decisão sua.

[`auto_switch`](/features/next#auto-switch) lista todas as condições e o que uma troca custa.

## Como usar duas contas do Claude Code no mesmo computador?

Dê a cada conta a própria pasta de configuração com `CLAUDE_CONFIG_DIR`. O ccam faz isso por você: execute `ccam add` para cada conta e depois inicie-as com `ccam personal` e `ccam work`, ou com atalhos que você mesmo nomeia, como `ccwork`. O [guia](/guides/multiple-claude-code-accounts) também mostra como fazer isso manualmente.

## Como continuar uma conversa em outra conta?

Execute `ccam next --continue` na pasta: ele leva a última conversa da pasta para a próxima conta (`--auto` para a que tem mais espaço). Dentro do Claude Code, o `/ccam:switch` do [plugin do ccam](/docs/install#claude-code-plugin) faz o mesmo com a sessão em que você está. Quando as duas contas compartilham o histórico de conversas (`ccam group set default --add history`, ou um grupo só delas: `ccam group add pair --items history --members work,side`), a nova conta abre essa mesma conversa; caso contrário, o Claude Code a copia para o histórico da nova conta. Adicione `--brief` para começar a nova conta com um resumo curto em vez da conversa inteira.

## O ccam vê ou envia meu login?

O ccam nunca compartilha nem copia seu login do Claude entre contas, e nunca copia os logins de servidores MCP. Ele lê o token de acesso de uma conta apenas para consultar o uso dessa conta no endpoint de uso da Anthropic: para `ccam usage`, o painel, `ccam doctor` e a escolha automática (`ccam next --auto`, `ccam pick`, `--fallback` e o `auto_switch` quando uma sessão atinge o limite). Ele só renova um token expirado com `ccam usage --refresh`; com `usage_source = "claude"`, ele não envia nenhum token para lugar nenhum e lê o login guardado (o arquivo ou o item do Acesso às Chaves do macOS) apenas para saber se ele ainda funciona. `ccam ext` e `ccam doctor` leem os nomes dos logins de MCP de uma conta, nunca os tokens deles. O ccam não tem telemetria. Além dessas consultas, a única requisição que ele faz por conta própria é a verificação do painel no GitHub por uma versão mais nova do ccam, que não leva nada sobre você ou suas contas.

## Como fazer login em um servidor remoto?

Execute `ccam add` ou `ccam login <name>` na sessão SSH, como de costume. O Claude Code espera o login no servidor, então o navegador do seu computador termina em uma página `http://localhost:<port>/callback?code=…` que não carrega. Copie esse endereço do navegador e cole no terminal: o ccam o repassa ao Claude Code no servidor, que conclui o login. O endereço contém um código de uso único, então não o compartilhe.

## Onde minhas contas ficam armazenadas?

Em `~/.local/share/ccam/profiles/<name>/`, uma pasta por conta. No macOS, o Claude Code guarda o próprio login no Acesso às Chaves, com um nome derivado dessa pasta, então `ccam rename` mantém ali o nome antigo da pasta.

## O que acontece com a minha pasta ~/.claude existente?

Nada, a menos que você execute `ccam import` ou `ccam migrate`. A importação registra a pasta no lugar onde ela está, então o `claude` simples e o `ccam <name>` continuam usando a mesma conta e as mesmas configurações. O `ccam migrate` copia as configurações, o `CLAUDE.md`, as extensões e os servidores MCP dela para o grupo `default` e só lê a pasta.

## Quais sistemas o ccam suporta?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine e outras distribuições), macOS e Windows 10 e 11, nativamente ou dentro do WSL. No Windows, compartilhar arquivos como `settings.json` exige o Modo de Desenvolvedor; pastas compartilhadas funcionam sem ele, e, sem ele, o primeiro `ccam add` cria o grupo `default` só com servidores MCP e extensões ([compartilhamento no Windows](/features/share#on-windows)).

## Por que um servidor MCP aparece em apenas uma conta?

Servidores adicionados dentro de uma sessão com `claude mcp add` ficam armazenados apenas para essa conta e, sem `-s user`, apenas para a pasta atual. O ccam avisa você na próxima vez que você iniciar a conta, e `ccam mcp list` e `ccam ext` os mostram; mova o servidor para o grupo compartilhado com `ccam mcp import <account> <server>`.

## O ccam deixa o Claude Code mais lento?

Iniciar o Claude Code pelo ccam acrescenta cerca de 4 milissegundos, medidos com `scripts/perf-check.sh` em um notebook. O ccam lê um arquivo pequeno, verifica os links compartilhados e depois se substitui pelo `claude`, então, por padrão, ele não fica em execução enquanto você trabalha. O Windows não consegue substituir um processo, então lá o ccam inicia o Claude Code e espera por ele, sem fazer mais nada. Com `relay` ou `auto_switch` ativado, o ccam continua rodando do mesmo jeito ao lado de cada sessão interativa, para poder continuar a conversa em outra conta.

## Como atualizar o ccam?

Execute `ccam update`. Ele baixa a versão mais recente, verifica o checksum e substitui o arquivo `ccam`; `ccam update --check` apenas informa se existe uma versão mais recente. O painel também avisa quando existe uma versão mais recente, e `U` a instala depois de perguntar; o ccam nunca se atualiza sozinho. Se você instalou o pacote `.deb`, instale o novo `.deb` com o apt. Consulte [atualização](/docs/install#update).

## Como desinstalar o ccam?

Execute `ccam uninstall`. Ele lista o que vai ser removido e pergunta antes; depois remove cada conta (desconectando-a primeiro), os atalhos, `~/.config/ccam`, `~/.local/share/ccam`, `~/.local/state/ccam` e o comando `ccam`. Um `~/.claude` importado permanece, com uma cópia do que compartilhava nos grupos dele. Consulte [instalação](/docs/install#uninstall).
