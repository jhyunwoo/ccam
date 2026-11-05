---
title: O que fazer no limite de uso do Claude Code
description: Espere o limite reiniciar ou continue na hora a mesma conversa em outra conta do Claude Code. Como funcionam os limites de 5 horas e semanal, passo a passo.
updated: 2026-11-05
order: 3
translation: machine
sourceHash: "53b1e3634ea9"
---

Quando uma conta atinge o limite de uso (muitas vezes chamado de rate limit), você tem duas opções: esperar até o limite reiniciar ou continuar agora a mesma conversa em outra conta sua. Com o ccam, a segunda opção é um único comando, `ccam next --auto --continue`.

## Quais limites o Claude Code tem?

Nos planos Pro e Max, o Claude Code tem um [limite de sessão de 5 horas e um limite semanal](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Os dois são compartilhados com o Claude na web e nos apps, e o uso conta para os dois ao mesmo tempo. Quando um deles se esgota, o Claude Code [mostra uma mensagem](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) como "You've hit your session limit · resets 3:45pm". Trocar de modelo não ajuda, porque os limites valem para todos os modelos.

## Opção 1: esperar o limite reiniciar

Se você tem uma só conta, espere. As versões recentes do Claude Code [esperam na sessão aberta](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) e continuam a tarefa sozinhas quando o limite reinicia. `/usage` dentro do Claude Code mostra quanto você já usou dos limites do seu plano.

## Opção 2: continuar agora em outra conta

Se você tem uma segunda conta, pode continuar a mesma conversa imediatamente. Configure uma vez:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Agora a linha de status mostra a conta e o uso de 5 horas dela enquanto você trabalha, então o limite não pega você de surpresa.

Quando o limite for atingido, saia do Claude Code (`/exit`) e depois:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Abra a conversa na conta com mais espaço:

```sh
ccam next --auto --continue
```

O Claude Code inicia como `work` com a conversa mais recente da pasta atual, incluindo tudo o que foi dito até agora. Compartilhar o histórico (`ccam group set default --add history`, acima) faz com que as duas contas trabalhem na mesma conversa; sem isso, o Claude Code copia a conversa para `work`, e o ccam informa quantos tokens a primeira resposta precisa ler de novo.

Com o [plugin do ccam](/docs/install#claude-code-plugin) no Claude Code, uma sessão que para no limite mostra uma notificação e registra a janela cheia, para que o `ccam next --auto` não escolha essa conta de novo antes do reinício. `/ccam:switch` (ou `/ccam:switch work`) registra a conversa para o `ccam next --continue`, para que ele continue exatamente aquela conversa.

**Assentos dos planos Team e Enterprise.** A escolha automática de conta (`ccam next`, `--fallback`, `ccam pick`, a tecla `n` do painel e, depois que você o ativa, o [`auto_switch`](/features/next#auto-switch)) também alterna entre assentos da mesma organização Claude Team ou Enterprise, o que pode violar os termos da Anthropic. Leia as [perguntas frequentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e consulte a sua organização e a Anthropic antes de contar com ela.

## O ccam troca de conta automaticamente?

Só se você ativar isso.

- Por padrão, o ccam não fica em execução enquanto você trabalha: ele inicia o Claude Code e se substitui por ele, então você decide quando trocar. `ccam next --auto` escolhe a conta pelo uso, pegando a de menor uso semanal entre as contas a até 3 pontos percentuais do menor uso de 5 horas; pula as contas que estão em 100% ou sem login e fica dentro de uma mesma organização Team ou Enterprise, ou entre planos pessoais.
- Com `relay = true` no `config.toml`, o ccam, em vez disso, continua rodando ao lado da sessão, então `/ccam:switch` continua a conversa sozinho, mantendo as opções, o modo de permissão, o nível de esforço e o agente da sessão; antes de mover a conversa para uma conta de outra organização, o ccam pergunta a você.
- Com `auto_switch = true` (Linux e macOS; experimental), o ccam também faz a troca sozinho: quando o Claude Code interrompe a sessão no limite dela e outra conta do mesmo tipo tem espaço, o ccam encerra a sessão e continua a conversa nessa conta ([condições e custos](/features/next#auto-switch)).

Algumas outras ferramentas trocam o login de uma sessão em andamento; [o ccam e outras ferramentas para trocar de conta](/compare/account-switchers) explica a diferença.

## O que ter em mente

- Execute `ccam next --continue` na mesma pasta em que você estava trabalhando; as conversas são guardadas por pasta de projeto.
- Evite executar duas contas na mesma conversa ao mesmo tempo. O ccam avisa se você tentar.
- Sem histórico compartilhado, a conversa é copiada para a outra conta e as duas cópias seguem caminhos separados; os checkpoints a partir dos quais o `/rewind` restaura o código ficam com a primeira conta.
- A primeira mensagem na nova conta pode consumir mais do limite dela do que o normal, porque essa conta ainda não tem a conversa em cache. `ccam next --continue --brief` a começa a partir de um resumo curto.
- Se você pode trocar de conta desse jeito depende dos termos que se aplicam a cada conta.

Veja [Limites de uso](/features/usage) e [Trocar para a próxima conta](/features/next) para todas as opções.
