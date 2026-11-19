---
title: Subagentes em outras contas
description: "ccam agent executa parte de uma tarefa do Claude Code em outra conta sua e devolve a resposta do Claude: uma tarefa longa aproveita o uso de várias contas."
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Os subagentes do próprio Claude Code rodam como a conta da sessão que os inicia, então gastam o uso dessa conta. `ccam agent "<task>"` executa uma tarefa com `claude -p` na mesma pasta como outra das suas contas e imprime a última resposta do Claude. Chamado de uma sessão do Claude Code, ele funciona como um subagente que gasta o uso de outra conta, então uma tarefa longa pode continuar depois do ponto em que a própria conta da sessão já teria se esgotado.

**Assentos dos planos Team e Enterprise.** O `ccam agent` escolhe a conta sozinho. Dentro de uma organização Claude Team ou Enterprise, ele executa a tarefa em outro assento, o que pode violar os termos da Anthropic. Leia as [perguntas frequentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e consulte a sua organização e a Anthropic antes de contar com ele.

## Uso

O Claude Code o executa pela ferramenta Bash dele. Com o [plugin do ccam](/docs/install#claude-code-plugin) instalado, peça isso com palavras simples, por exemplo "use minhas outras contas para partes desta tarefa", ou coloque isso no seu CLAUDE.md. A skill do plugin diz ao Claude Code para iniciar cada agente em segundo plano, escrever uma tarefa que se sustente sozinha e dar aos agentes que editam arquivos uma parte própria ou um git worktree próprio.

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

A última resposta do Claude vai para a saída padrão; as linhas que começam com `ccam:` vão para a saída de erro padrão. Com `--json`, ele imprime `{account, session_id, result, is_error, accounts}` em vez disso, em que `accounts` lista todas as contas com que a tarefa rodou.

Códigos de saída: o do próprio claude (`0` quando o agente terminou), `75` (`EX_TEMPFAIL`) quando nenhuma conta além da que chamou tem espaço, `1` para outros erros, como um agente tentando iniciar outro.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--model <model>` | O modelo do agente; a janela semanal desse modelo também conta na escolha | o modelo padrão da conta |
| `--resume <id>` | Continua a conversa de um agente anterior com uma nova tarefa (o ID de sessão na primeira e na última linha dele) | desativado |
| `--json` | Imprime a conta, a sessão, a resposta e as contas com que ele rodou | desativado |
| `--max-usage <1-100>` | Prefere as contas abaixo deste uso | `auto_max_usage`, senão 100 |
| `--strategy headroom\|expiring\|balance` | A [regra de escolha](/features/next#auto) | `auto_strategy`, senão `headroom` |
| `--cross-org` | Também roda como contas de outras organizações | `cross_org`, senão desativado |
| `--wait[=<duration>]` | Quando nenhuma conta tem espaço, espera o reinício da primeira, no máximo por este tempo | desativado; `5h` quando usado sozinho |
| `--dry-run` | Imprime o comando para a conta que escolheria, sem executá-lo | desativado |

## Como funciona

1. **Quem chama fica de fora.** A conta da qual o `ccam agent` é chamado (a da sessão que o ccam iniciou, senão a que esta pasta inicia) nunca executa o agente, nem outra pasta logada na mesma conta Claude, que compartilha os limites dela. O agente fica na [organização](/features/next#auto) de quem chama, a menos que `--cross-org` diga o contrário.
2. **Ele escolhe como o `ccam pick`.** O uso é consultado do mesmo jeito que o [`ccam pick`](/features/pick) faz, e agentes iniciados com até 2 minutos de diferença recebem contas diferentes enquanto houver contas suficientes com espaço.
3. **Ele roda com as permissões da sessão, nunca mais.** Ninguém está lá para permitir uma ferramenta no `claude -p`, então o agente recebe o modo de permissão e o nível de esforço da sessão que o chamou. Executado de um terminal, valem o modo padrão do Claude Code e as regras de permissão das suas configurações. O `ccam agent` não tem nenhuma opção que conceda mais.
4. **No limite de uso, a conversa segue adiante.** Quando a conta do agente atinge o limite, o ccam registra o limite como o [hook do plugin](/features/usage#the-limit-record) faria e continua a mesma conversa como a conta com mais espaço: o Claude Code a copia para lá (`--resume <file> --fork-session`), com a primeira mensagem "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." Cada conta é tentada uma vez.
5. **Toda conversa de agente é registrada** em `~/.local/state/ccam/agent/<session id>.json` (30 dias). `--resume <id>` a continua na pasta dela, como a conta que a executou, se essa ainda tiver espaço, senão como outra. O [`ccam next --continue`](/features/next) deixa de fora as conversas de agentes, então continua retomando as suas.

## Segurança e casos especiais

- Um agente não pode iniciar agentes próprios: o ccam o marca com `CCAM_AGENT=1`, e o `ccam agent` se recusa a rodar ali.
- Parar o `ccam agent` (Ctrl+C, um sinal ou o Claude Code interrompendo o comando) também para o agente.
- O agente não vê a conversa que o chamou, só a tarefa dele.
- Agentes que editam os mesmos arquivos ao mesmo tempo sobrescrevem o trabalho um do outro. Dê a cada um uma parte própria da tarefa, ou um git worktree próprio.
- Quando a sessão que chamou passa para outra conta ([`auto_switch`](/features/next#auto-switch)), os agentes que ela iniciou param junto com ela. As conversas deles continuam lá: `--resume <id>` as continua.
- `agent` é um nome de comando, então nenhuma conta nova pode se chamar `agent`; uma que já existe ainda inicia com `ccam run agent`.

## Relacionados

- [Escolher uma conta para scripts](/features/pick)
- [Trocar para a próxima conta](/features/next)
- [Limites de uso](/features/usage)
- [Sessões em segundo plano de todas as contas](/features/agents)
