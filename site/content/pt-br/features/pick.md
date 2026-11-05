---
title: Escolher uma conta para scripts
description: ccam pick imprime a conta com mais espaço para um script, dá contas diferentes a workers iniciados juntos e sai com 75 quando nenhuma tem espaço.
updated: 2026-11-05
order: 12
translation: machine
sourceHash: "0ec25752ace9"
---

`ccam pick` imprime o nome da conta com login feito que tem mais espaço, para que um script possa iniciar o Claude Code com ela. Workers iniciados ao mesmo tempo recebem cada um uma conta diferente, desde que contas suficientes tenham espaço, então um lote de jobs `claude -p` se espalha pelas suas contas em vez de esgotar uma só.

**Assentos dos planos Team e Enterprise.** A escolha automática de conta (`ccam next`, `--fallback`, `ccam pick`, a tecla `n` do painel e, depois que você o ativa, o [`auto_switch`](/features/next#auto-switch)) também alterna entre assentos da mesma organização Claude Team ou Enterprise, o que pode violar os termos da Anthropic. Leia as [perguntas frequentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e consulte a sua organização e a Anthropic antes de contar com ela.

## Uso

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

Códigos de saída: `0` com um nome na saída padrão, `75` (`EX_TEMPFAIL`) quando nenhuma conta tem espaço ou nenhuma tem login feito, `1` para outros erros, como uma opção inválida.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--json` | Imprime a conta, o uso de 5 horas e o semanal, de onde vêm os números e por que ela foi escolhida | desativado |
| `--model <model>` | Também conta a janela semanal desse modelo, como `opus` | nenhum |
| `--max-usage <1-100>` | Prefere as contas abaixo deste uso | `auto_max_usage`, senão 100 |
| `--strategy headroom\|expiring\|balance` | O menor uso semanal a até 3 pontos percentuais do menor uso de 5 horas, o reinício de 5 horas mais próximo ou a maior cota semanal restante por dia | `auto_strategy`, senão `headroom` |
| `--cross-org` | Também escolhe contas de outras organizações | `cross_org`, senão desativado |
| `--wait[=<duration>]` | Quando nenhuma conta tem espaço, espera o reinício da primeira, no máximo por este tempo. O valor vai depois de `=`: `--wait 30m` é recusado com a dica `--wait takes its value after '=': --wait=30m` | desativado; `5h` quando usado sozinho |
| `--dry-run` | Escolhe sem registrar a escolha | desativado |

## Como funciona

1. O ccam lista as contas com login feito (contas com chave de API e logins do Console não têm limites e ficam de fora) na organização da conta que esta pasta iniciaria (consulte [organizações](/features/next#auto)); uma pasta logada na mesma conta Claude que essa também fica de fora.
2. Ele pega um lock (`~/.local/state/ccam/picks.lock`), para que escolhas feitas ao mesmo tempo aconteçam uma depois da outra.
3. Consulta o uso de todas as contas ao mesmo tempo, do mesmo jeito que o `ccam usage`, com o cache de 60 segundos. Uma janela que o Claude Code informou como cheia conta como 100% até reiniciar ([o registro de limite](/features/usage#the-limit-record)), então uma conta que acabou de atingir o limite é excluída mesmo que os números salvos dela ainda mostrem espaço.
4. Contas escolhidas nos últimos 2 minutos, pelo `ccam pick` ou pelo [`auto_switch`](/features/next#auto-switch) ao mover uma sessão para elas (`~/.local/state/ccam/picks/<name>`; uma escolha com horário no futuro não conta), ficam de fora primeiro, e as [regras de escolha](/features/next#auto) do `ccam next --auto` escolhem entre as restantes: primeiro as abaixo de `auto_max_usage`, pela estratégia, nunca uma conta em 100%. Se isso pular a conta que as regras escolheriam de outra forma, o `reason` do `--json` informa isso: `30% of its 5 hours used; work was picked a moment ago`.
5. Se todas as contas com espaço foram escolhidas há pouco, as regras escolhem de novo entre todas elas, então, com três contas, um quarto worker divide uma delas.
6. A escolha é registrada e o lock liberado, e então o nome é impresso.

A janela de 2 minutos importa porque o uso de um job novo demora um pouco para aparecer nos números: sem ela, cinco workers iniciados juntos veriam todos a mesma conta como a de mais espaço.

## Segurança e casos especiais

- `ccam pick` só lê o uso e grava o horário da escolha; ele nunca renova nem lê tokens além do que o `ccam usage` faz.
- Uma escolha não reserva a conta. Dois scripts ainda podem começar na mesma conta depois de 2 minutos; use `--fallback` para se recuperar quando essa conta atingir o limite.
- Esperar com `--wait` não segura nenhum lock, então outras escolhas continuam nesse meio-tempo.
- `ccam pick` ignora o [`ccam folder`](/features/folder): ele escolhe só pelo uso.

## Relacionados

- [Trocar para a próxima conta](/features/next)
- [Executar de novo no limite](/features/next#run-again-at-the-limit)
- [Limites de uso](/features/usage)
- [Contas com chave de API](/features/api)
