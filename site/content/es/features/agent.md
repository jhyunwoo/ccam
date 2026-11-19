---
title: Subagentes en otras cuentas
description: "ccam agent ejecuta parte de una tarea de Claude Code con otra de tus cuentas y devuelve la respuesta de Claude, para que una tarea larga use varias cuentas."
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Los subagentes propios de Claude Code se ejecutan con la cuenta de la sesión que los inicia, así que gastan el uso de esa cuenta. `ccam agent "<task>"` ejecuta una tarea con `claude -p` en la misma carpeta con otra de tus cuentas e imprime la última respuesta de Claude. Llamado desde una sesión de Claude Code, funciona como un subagente que gasta el uso de otra cuenta, así que una tarea larga puede seguir después de que la cuenta propia de la sesión se hubiera agotado.

**Asientos de Team y Enterprise.** `ccam agent` elige la cuenta por sí mismo. Dentro de una misma organización de Claude Team o Enterprise, ejecuta la tarea en otro asiento, lo que puede infringir los términos de Anthropic. Lee las [preguntas frecuentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) y consulta con tu organización y con Anthropic antes de depender de ello.

## Uso

Claude Code lo ejecuta con su herramienta Bash. Con el [plugin de ccam](/docs/install#claude-code-plugin) instalado, pídelo con tus propias palabras, por ejemplo "usa mis otras cuentas para partes de esta tarea", o ponlo en tu CLAUDE.md. El skill del plugin le indica a Claude Code que inicie cada agente en segundo plano, que escriba una tarea que se entienda por sí sola y que dé a los agentes que editan archivos su propia parte o su propio worktree de git.

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

La última respuesta de Claude va a la salida estándar; las líneas que empiezan por `ccam:` van a la salida de error estándar. `--json` imprime en su lugar `{account, session_id, result, is_error, accounts}`, donde `accounts` enumera todas las cuentas con las que se ejecutó la tarea.

Códigos de salida: el del propio claude (`0` cuando el agente terminó), `75` (`EX_TEMPFAIL`) cuando ninguna cuenta distinta de la que lo llama tiene margen, `1` para otros errores, como un agente que intenta iniciar otro.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--model <model>` | El modelo del agente; la ventana semanal de ese modelo también cuenta en la elección | el modelo predeterminado de la cuenta |
| `--resume <id>` | Continúa la conversación de un agente anterior con una tarea de seguimiento (el ID de sesión aparece en su primera y su última línea) | desactivada |
| `--json` | Imprime la cuenta, la sesión, la respuesta y las cuentas con las que se ejecutó | desactivada |
| `--max-usage <1-100>` | Prefiere las cuentas por debajo de este uso | `auto_max_usage`, o 100 |
| `--strategy headroom\|expiring\|balance` | La [regla de elección](/features/next#auto) | `auto_strategy`, o `headroom` |
| `--cross-org` | Se ejecuta también con cuentas de otras organizaciones | `cross_org`, o desactivada |
| `--wait[=<duration>]` | Cuando ninguna cuenta tiene margen, espera a que se restablezca la primera, como máximo este tiempo | desactivada; `5h` si se da sola |
| `--dry-run` | Imprime el comando para la cuenta que elegiría, sin ejecutarlo | desactivada |

## Cómo funciona

1. **La cuenta que lo llama queda fuera.** La cuenta desde la que se llama a `ccam agent` (la de la sesión que inició ccam o, si no, la que inicia esta carpeta) nunca ejecuta el agente, ni tampoco otra carpeta con la sesión iniciada en la misma cuenta de Claude, que comparte sus límites. El agente se queda en la [organización](/features/next#auto) de quien lo llama, salvo que `--cross-org` diga lo contrario.
2. **Elige como `ccam pick`.** El uso se consulta como lo hace [`ccam pick`](/features/pick), y los agentes iniciados con menos de 2 minutos de diferencia reciben cuentas distintas siempre que haya suficientes con margen.
3. **Se ejecuta con los permisos de la sesión, nunca con más.** En `claude -p` no hay nadie para permitir una herramienta, así que el agente recibe el modo de permisos y el nivel de esfuerzo de la sesión que lo llamó. Ejecutado desde una terminal, se aplican el modo predeterminado de Claude Code y las reglas de permisos de tu configuración. `ccam agent` no tiene ninguna opción que conceda más.
4. **Al llegar a un límite de uso, la conversación sigue.** Cuando la cuenta del agente llega a su límite, ccam registra el límite como lo haría el [hook del plugin](/features/usage#the-limit-record) y continúa la misma conversación con la cuenta que tiene más margen: Claude Code la copia allí (`--resume <file> --fork-session`), con el primer mensaje "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." Cada cuenta se intenta una sola vez.
5. **Cada conversación de agente queda registrada** en `~/.local/state/ccam/agent/<session id>.json` (30 días). `--resume <id>` la continúa en su carpeta, con la cuenta que la ejecutó si todavía tiene margen o, si no, con otra. [`ccam next --continue`](/features/next) deja fuera las conversaciones de agentes, así que sigue continuando las tuyas.

## Seguridad y casos límite

- Un agente no puede iniciar agentes propios: ccam lo marca con `CCAM_AGENT=1`, y `ccam agent` se niega a ejecutarse ahí.
- Detener `ccam agent` (Ctrl+C, una señal o Claude Code interrumpiendo el comando) también detiene el agente.
- El agente no ve la conversación que lo llamó, solo su tarea.
- Los agentes que editan los mismos archivos a la vez sobrescriben el trabajo de los demás. Da a cada uno su propia parte de la tarea, o su propio worktree de git.
- Cuando la sesión que los llamó pasa a otra cuenta ([`auto_switch`](/features/next#auto-switch)), los agentes que inició se detienen con ella. Sus conversaciones permanecen: `--resume <id>` las continúa.
- `agent` es un nombre de comando, así que ninguna cuenta nueva puede llamarse `agent`; una existente todavía se inicia con `ccam run agent`.

## Relacionado

- [Elegir una cuenta para scripts](/features/pick)
- [Cambiar a la siguiente cuenta](/features/next)
- [Límites de uso](/features/usage)
- [Sesiones en segundo plano de todas las cuentas](/features/agents)
