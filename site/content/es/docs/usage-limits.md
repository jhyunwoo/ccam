---
title: Límites de uso y cambio de cuenta
description: Consulta el uso de 5 horas y semanal de cada cuenta de Claude Code con ccam usage, y abre la cuenta con más margen con ccam next --auto.
updated: 2026-11-05
order: 5
translation: machine
sourceHash: "06bd349073e3"
---

`ccam usage` muestra, para cada cuenta, cuánto se ha usado de cada límite: el límite de 5 horas, el límite semanal y el límite semanal propio de cada modelo (como Fable). `ccam next --auto` abre la cuenta con más margen, y `--continue` lleva consigo la conversación actual.

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

Cada cuenta tiene una línea con la antigüedad de sus cifras y, debajo, una barra por límite con el tiempo hasta que se restablece: la ventana de 5 horas, la ventana semanal de todos los modelos y la ventana semanal propia de cada modelo, como Fable. Una barra se vuelve amarilla al 80 % y roja al 100 %. La última línea nombra la cuenta que abriría `ccam next --auto`; con `auto_strategy = "expiring"` o `"balance"` dice en su lugar `→ next --auto picks <name>: <reason>`.

**Asientos de Team y Enterprise.** La selección automática de cuenta (`ccam next`, `--fallback`, `ccam pick`, la tecla `n` del panel y, una vez que lo activas, [`auto_switch`](/features/next#auto-switch)) también pasa de un asiento a otro dentro de una misma organización de Claude Team o Enterprise, lo que puede infringir los términos de Anthropic. Lee las [preguntas frecuentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) y consulta con tu organización y con Anthropic antes de depender de ella.

## De dónde salen las cifras

ccam consulta el endpoint de uso de Anthropic para todas las cuentas al mismo tiempo, con el inicio de sesión propio de cada cuenta. Espera como máximo 3 segundos por cuenta y muestra cada cuenta en cuanto responde.

Este endpoint no está documentado públicamente y puede cambiar. Cuando una solicitud falla, ccam muestra las últimas cifras conocidas y las marca como almacenadas en caché. Una ventana que se haya reiniciado desde que se guardaron esas cifras muestra 0%. `ccam doctor` te avisa si el formato de la respuesta cambió.

Cuando Claude Code detiene una sesión en un límite de uso, el [plugin de ccam](/docs/install#claude-code-plugin) registra qué ventana se llenó y cuándo se restablece. Hasta entonces esa ventana muestra 100 % digan lo que digan las demás cifras, así que `ccam next --auto` y `--fallback` no vuelven a elegir la cuenta demasiado pronto.

Las cifras obtenidas hace menos de 60 segundos se reutilizan. Opciones:

| Opción | Efecto |
| --- | --- |
| `--watch[=seconds]` | Actualiza de forma continua, cada 60 segundos de forma predeterminada |
| `--json` | Salida legible por máquinas |
| `-v` | Muestra también de dónde salieron las cifras de cada cuenta |
| `--no-cache` | Siempre vuelve a consultar |
| `--source cache` | Nunca consulta; solo muestra las cifras guardadas |

## Inicios de sesión vencidos

Cuando el token de acceso de una cuenta venció, ccam no lo renueva, porque renovarlo mientras una sesión de Claude Code usa el mismo inicio de sesión podría cerrar esa sesión. La cuenta muestra ``access token expired; `ccam <name>` renews it when it starts``: iníciala una vez con `ccam <name>` y Claude Code lo renueva. Un inicio de sesión que venció del todo, o que Claude Code borró, muestra `` login expired → `ccam login <name>` `` y necesita un nuevo inicio de sesión.

`ccam usage --refresh` renueva por sí mismo los tokens vencidos, pero solo de las cuentas sin ninguna sesión de Claude Code en ejecución, sin importar cómo se haya iniciado esa sesión. Nunca renueva inicios de sesión guardados en el Llavero de macOS ni en un `~/.claude` importado, y una renovación puede tardar hasta 30 segundos.

## Uso en la línea de estado

```sh
ccam statusline install
```

Entonces Claude Code muestra el nombre de la cuenta y su uso en la línea de estado, por ejemplo `[work] · Opus · 5h 24% · 7d 41%`. A partir del 90 %, la línea también sugiere `ccam next --auto --continue`. Si ya tienes un comando de línea de estado, ccam lo conserva y pone el nombre de la cuenta delante. La línea de estado también guarda sus cifras, que `ccam usage` usa cuando el endpoint no responde.

## Cambiar de cuenta

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the most room
ccam next --continue  # continue this folder's conversation on the next account
```

`ccam next` omite las cuentas que no tienen sesión iniciada, y `--auto` omite las cuentas al 100 % (incluida una ventana que Claude Code informó llena, hasta su restablecimiento) y las cuentas de otra organización. Cuando las dos cuentas comparten el historial de conversaciones, `--continue` abre la misma conversación; si no, Claude Code la copia en la nueva cuenta. Consulta [compartir](/docs/sharing) y [cambiar a la siguiente cuenta](/features/next#continue).

Con `auto_switch = true` al principio de `config.toml` (Linux y macOS; desactivado por defecto), ccam hace este cambio por sí mismo cuando Claude Code detiene una sesión interactiva en su límite y otra cuenta del mismo tipo (un asiento de la misma organización u otro plan personal) tiene margen. Consulta [`auto_switch`](/features/next#auto-switch) para saber cuándo actúa y qué no hace nunca.

Para scripts y tareas programadas que ejecutan `claude -p`, pon `--fallback` antes del nombre de la cuenta. Cuando esa cuenta llega a su límite, ccam vuelve a ejecutar el mismo comando con la cuenta con más margen, y solo se imprime la salida de la ejecución que terminó:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Consulta [Volver a ejecutar al llegar al límite](/features/next#run-again-at-the-limit) para ver cómo decide ccam.

La [guía del flujo de trabajo con límites de uso](/guides/usage-limit-workflow) recorre un ejemplo completo.
