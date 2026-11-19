---
title: Referência de comandos
metaTitle: Referência de comandos do ccam
description: Todos os comandos e opções do ccam, de ccam add e ccam <name> a ccam usage, ccam next e ccam doctor.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: "28e1ad1dcbad"
---

Coloque `--dry-run` antes de um comando para ver o que ele faria sem fazê-lo, por exemplo `ccam --dry-run work`. Funciona ao iniciar o Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) e com `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` e `doctor --fix`. Os comandos que só leem, como `list`, `usage`, `share` e `ext`, aceitam a opção e, de qualquer forma, não alteram nada. Os demais recusam a opção em vez de executar, inclusive `group rename`.

<!-- commands -->

## Iniciar o Claude Code

`ccam <name|number> [arguments]` inicia o Claude Code como essa conta. O ccam não lê os argumentos depois do nome; ele os repassa ao `claude` na mesma ordem, e a entrada padrão, o código de saída e os sinais se comportam exatamente como no próprio `claude`. Um `--` depois do nome é opcional.

`ccam` com argumentos, mas sem nome de conta, como `ccam -p "summarise"`, usa a conta definida para a pasta atual ou, se não houver, a conta padrão. Defina-as com `ccam folder <name> [folder]` e `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` antes do nome da conta serve para execuções `-p`/`--print`. Quando a conta atinge o limite de uso, o ccam executa o mesmo comando de novo com a conta logada que tem mais espaço e imprime só a saída da execução que terminou. Consulte [Contas por pasta](/features/folder) e [Executar de novo no limite](/features/next#run-again-at-the-limit).

**Assentos dos planos Team e Enterprise.** A escolha automática de conta (`ccam next`, `--fallback`, `ccam pick`, a tecla `n` do painel e, depois que você o ativa, o [`auto_switch`](/features/next#auto-switch)) também alterna entre assentos da mesma organização Claude Team ou Enterprise, o que pode violar os termos da Anthropic. Leia as [perguntas frequentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e consulte a sua organização e a Anthropic antes de contar com ela.

Um atalho como `ccwork` inicia a conta dele do mesmo jeito: `ccwork --resume` faz o mesmo que `ccam work --resume`. Consulte [Atalhos](/docs/commands#shortcuts) abaixo.

## Adicionar e remover contas

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` sem uma pasta procura `~/.claude` e `~/.claude-*` e pergunta sobre cada uma; uma pasta que falha é informada como `Skipped <path>: <reason>` e as outras são registradas. Se o nome de uma pasta já estiver em uso, importe essa pasta sozinha com `ccam import <folder> --name <name>`. As pastas importadas são registradas no lugar onde estão; o ccam nunca move nem copia um login.

`ccam login --expired` também faz login nas contas cujo login o Claude Code apagou. `ccam remove` desconecta uma conta gerenciada (uma cuja pasta o ccam criou, e não uma registrada com `ccam import`) antes de a pasta dela ir para a lixeira, a menos que outra conta tenha o mesmo login. Remover uma conta importada a tira dos grupos e deixa na pasta dela uma cópia de cada item compartilhado. No macOS, `ccam rename` mantém o nome da pasta, já que o Acesso às Chaves guarda o login pelo caminho da pasta.

`ccam add --api` adiciona uma conta que usa uma chave de API em vez de um login do Claude; o ccam lê a chave a cada início e nunca a guarda. Veja [Contas com chave de API](/features/api).

`ccam rename`, `ccam remove` e os comandos `ccam group` se recusam a alterar uma conta enquanto o Claude Code estiver em execução com essa conta, seja qual for a forma como ele foi iniciado (inclusive sessões em segundo plano). Encerre a sessão primeiro. Eles também se recusam enquanto uma alteração de grupo estiver inacabada: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` mostra todas as contas com o número (slot), atalho, e-mail e pasta de cada uma. O array `groups` lista os grupos da conta, e o objeto `share` associa cada item (`mcp`, `settings`, `memory`, `extensions`, `history`) ao grupo do qual ele vem, ou a `local`.

## Atalhos

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Um atalho é um comando que você mesmo nomeia e que inicia o Claude Code como uma determinada conta. Assim como `ccam work`, ele repassa todos os argumentos ao Claude Code sem alterações, então `ccwork --resume` e `ccwork -p "hi"` funcionam como você espera. Ele é um link para o ccam em `~/.local/bin`.

Os atalhos são opcionais. O assistente do `ccam add` pergunta se você quer criar um, `ccam add <name> --shortcut <command>` cria um sem perguntar e `ccam list` mostra o atalho de cada conta, ou `-` quando ela não tem.

Os nomes usam os caracteres a-z, 0-9, `.`, `-` e `_`. Um nome não pode ser um comando do ccam, o atalho de outra conta nem o nome de outro programa no seu `PATH`, nem um arquivo em `~/.local/bin` que o ccam não criou, como uma cópia ou um link seu do ccam. Um atalho editado à mão no `config.toml` que quebra essas regras, que é `claude` ou `ccam`, ou que foi dado a duas contas faz o `config.toml` falhar ao carregar, indicando a conta.

`ccam rename` mantém o nome do atalho, e `ccam remove` exclui o atalho junto com a conta. `ccam shims install` recria os links de todas as contas, por exemplo depois que você move o ccam; um arquivo no caminho de um atalho não impede os outros.

## Compartilhamento

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Um grupo compartilha itens entre as contas dele; uma conta pode estar em vários grupos, mas cada item vem de um único grupo. Consulte [compartilhamento com grupos](/features/share).

## Copiar uma configuração existente

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` copia `settings.json`, `CLAUDE.md`, extensões e servidores MCP de usuário de uma pasta do Claude Code para um grupo, para que todas as contas que compartilham esses itens os recebam. A pasta é apenas lida.

- `--settings`, `--memory`, `--extensions`, `--mcp` e `--history` copiam só os itens que você indicar. Sem elas, em um terminal, uma lista de seleção permite escolher, com tudo marcado exceto o histórico de conversas; fora de um terminal, tudo é copiado, exceto o histórico de conversas.
- `--mcp` também funciona quando o `~/.claude` está vazio ou não existe, mas o `~/.claude.json` tem servidores.
- `-g` ou `--group` escolhe o grupo.
- `-y` ou `--yes` copia sem perguntar e mantém a cópia do grupo quando um arquivo é diferente.

Os logins de MCP não são copiados, já que um login copiado desconecta os outros que o têm; o plano indica as contas que fazem login com `/mcp`. Consulte [copiar uma configuração existente](/docs/sharing#copy-an-existing-setup).

## Servidores MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` move para o grupo da conta os servidores que só essa conta tem, inclusive os que ela tem só para uma pasta (`claude mcp add` simples), que passam a carregar em todas as pastas. Um servidor cujo nome o grupo já usa para uma definição diferente é ignorado; adicione `--force` para substituir a cópia do grupo, da qual o ccam faz backup antes.

## Uso e troca de conta

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam sends no token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --continue --brief    # a new conversation from a summary of the last one
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agent "…"                  # from Claude Code: a subagent that runs as another account
ccam agent --resume <id> "…"    # follow that agent up
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
ccam ccusage claude daily       # Claude Code's tokens and costs, every account (ccusage)
ccam ccusage -a work blocks --active   # one account's current 5-hour block
```

- `next --auto`, `--fallback` e `ccam pick` escolhem do mesmo jeito: primeiro as contas abaixo de `auto_max_usage` (`config.toml`, padrão 100), por `auto_strategy` (`headroom`, `expiring` ou `balance`), contando a janela semanal do `--model`. Eles nunca escolhem uma conta em 100% (uma janela que o Claude Code informou como cheia conta como 100% até reiniciar) nem um login com chave de API, e ficam dentro da organização da conta de onde partem (`--cross-org` ou `cross_org = true` para ir além).
- `--wait` recebe o valor depois de `=` (`--wait=2h`); escrito com espaço, é recusado com uma dica. Quando nenhuma pode ser escolhida, saem com código 75 ou, com `--wait`, esperam o reinício da primeira conta.
- `ccam next --continue` leva uma conversa adiante mesmo entre contas que não compartilham o histórico: nesse caso, o Claude Code a copia para a nova conta. Com `--brief`, a nova conta começa uma conversa nova, a partir de um resumo da última que o ccam escreve sem usar um modelo ([`--brief`](/features/next#brief)).
- Diferente de `ccam <name>`, o `ccam next` lê as próprias opções, então as opções do Claude Code vão depois de `--`; sem isso, o ccam para com `unknown flag`.
- `ccam agent` executa uma tarefa com `claude -p` nesta pasta como uma conta diferente daquela de onde é chamado, no modo de permissão e no nível de esforço da sessão que o chama, e imprime a última resposta do Claude. No limite dessa conta, a conversa segue como outra ([Subagentes em outras contas](/features/agent)).
- Dentro do Claude Code, o `/ccam:switch [account]` do plugin passa uma conversa adiante, e `ccam hook` é o que os hooks dele executam.
- Com `auto_switch = true` no `config.toml` (Linux e macOS; desativado por padrão), o hook de limite também encerra uma sessão interativa que o ccam iniciou e a continua em outra conta com espaço, mas só em uma do mesmo tipo: outro assento da mesma organização ou outro plano pessoal ([`auto_switch`](/features/next#auto-switch)).

Veja [Trocar para a próxima conta](/features/next), [Escolher uma conta para scripts](/features/pick), [Subagentes em outras contas](/features/agent) e [Sessões em segundo plano de todas as contas](/features/agents).

## Manutenção

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam uninstall       # remove every account, ccam's folders and ccam itself
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` verifica o checksum SHA-256 da versão antes de substituir o ccam, e `--version v0.2.0` instala uma versão específica. Uma cópia do ccam instalada pelo pacote `.deb` é atualizada com o apt. Consulte [atualização](/docs/install#update).

`ccam uninstall` lista o que remove e pergunta antes: cada conta, como `ccam remove` faz, os atalhos, as pastas do ccam e o comando `ccam`. Uma pasta importada permanece onde está. Consulte [desinstalar](/docs/install#uninstall).

## Arquivos

| Caminho | Conteúdo |
| --- | --- |
| `~/.config/ccam/config.toml` | Contas, grupos e configurações (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). As configurações ficam no início, acima do primeiro `[[group]]`, `[[profile]]` ou `[[folder]]`: uma configuração abaixo de uma tabela pertence a essa tabela, então o ccam avisa sobre ela a cada inicialização e no `ccam doctor`, e os comandos que salvam a configuração se recusam até que ela seja movida. Os comentários não são mantidos quando o ccam salva o arquivo |
| `~/.local/share/ccam/profiles/<name>/` | O `CLAUDE_CONFIG_DIR` de cada conta |
| `~/.local/share/ccam/groups/<group>/` | Configurações, servidores MCP, extensões e histórico compartilhados |
| `~/.local/state/ccam/` | Cache de uso, registros de limite, bloqueios, notas e outros estados (veja a lista abaixo da tabela) |
| `~/.local/bin/<shortcut>` | Atalhos como `ccwork`, links para o ccam |

A pasta de estado guarda:

- o cache de uso e os registros de limite (`usage/<name>.json`, `usage/<name>.limit.json`);
- os horários do último uso;
- os marcadores de sessão e os bloqueios (`locks/`);
- os atalhos que o ccam criou (`shortcuts.json`);
- uma alteração de grupo que não terminou (`group-change.json`);
- as notas que o `/ccam:switch` e o hook de limite deixam (`handoff/<account>/<session>.json`);
- quando o `ccam pick` ou o `auto_switch` escolheu cada conta pela última vez (`picks/<name>`; uma escolha conta por 2 minutos);
- os resumos do `--brief`, guardados por 7 dias (`brief/<session>.md`).

Definir `XDG_CONFIG_HOME`, `XDG_DATA_HOME` ou `XDG_STATE_HOME` muda o local dessas pastas. No Windows, `~` é `%USERPROFILE%` e os atalhos são `<shortcut>.exe`.
