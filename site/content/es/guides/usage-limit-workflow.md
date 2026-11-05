---
title: Qué hacer cuando Claude Code llega al límite de uso
description: Espera al reinicio o continúa al instante la misma conversación en otra cuenta de Claude Code. Cómo funcionan los límites de 5 horas y semanal, paso a paso.
updated: 2026-11-05
order: 3
translation: machine
sourceHash: "53b1e3634ea9"
---

Cuando una cuenta llega a su límite de uso (a menudo llamado rate limit), tienes dos opciones: esperar a que el límite se reinicie o continuar ahora la misma conversación en otra cuenta tuya. Con ccam, lo segundo es un solo comando, `ccam next --auto --continue`.

## ¿Qué límites tiene Claude Code?

En los planes Pro y Max, Claude Code tiene un [límite de sesión de 5 horas y un límite semanal](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Ambos se comparten con Claude en la web y en las apps, y el uso cuenta para los dos a la vez. Cuando uno se agota, Claude Code [muestra un mensaje](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) como "You've hit your session limit · resets 3:45pm". Cambiar de modelo no ayuda, porque los límites cubren todos los modelos.

## Opción 1: esperar al reinicio

Si tienes una sola cuenta, espera. Las versiones recientes de Claude Code [esperan en la sesión abierta](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) y continúan la tarea por su cuenta cuando el límite se reinicia. `/usage` dentro de Claude Code muestra cuánto has usado de los límites de tu plan.

## Opción 2: continuar ahora en otra cuenta

Si tienes una segunda cuenta, puedes continuar de inmediato la misma conversación. Configúralo una vez:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Ahora la línea de estado muestra la cuenta y su uso de 5 horas mientras trabajas, así que el límite no te toma por sorpresa.

Cuando llegues al límite, sal de Claude Code (`/exit`) y luego:

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

Abre la conversación en la cuenta con más margen:

```sh
ccam next --auto --continue
```

Claude Code se inicia como `work` con la conversación más reciente de la carpeta actual, incluido todo lo dicho hasta ahora. Compartir el historial (`ccam group set default --add history` arriba) hace que ambas cuentas trabajen en la misma conversación; sin él, Claude Code copia la conversación en `work`, y ccam te indica cuántos tokens tiene que volver a leer la primera respuesta.

Con el [plugin de ccam](/docs/install#claude-code-plugin) en Claude Code, una sesión que se detiene en su límite muestra una notificación y registra la ventana llena, así que `ccam next --auto` no vuelve a elegir esa cuenta antes de su restablecimiento. `/ccam:switch` (o `/ccam:switch work`) anota la conversación para `ccam next --continue`, así que continúa exactamente esa conversación.

**Asientos de Team y Enterprise.** La selección automática de cuenta (`ccam next`, `--fallback`, `ccam pick`, la tecla `n` del panel y, una vez que lo activas, [`auto_switch`](/features/next#auto-switch)) también pasa de un asiento a otro dentro de una misma organización de Claude Team o Enterprise, lo que puede infringir los términos de Anthropic. Lee las [preguntas frecuentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) y consulta con tu organización y con Anthropic antes de depender de ella.

## ¿ccam cambia de cuenta automáticamente?

Solo si lo activas.

- Por defecto, ccam no se está ejecutando mientras trabajas: inicia Claude Code y se reemplaza a sí mismo por él, así que tú decides cuándo cambiar. `ccam next --auto` elige según el uso y toma, de las cuentas que están a 3 puntos porcentuales o menos del menor uso de 5 horas, la de menor uso semanal; omite las cuentas que están al 100% o sin sesión iniciada, y se queda dentro de una misma organización de Team o Enterprise, o entre planes personales.
- Con `relay = true` en `config.toml`, ccam sigue ejecutándose en cambio junto a la sesión, así que `/ccam:switch` continúa la conversación por sí solo, conservando las opciones, el modo de permisos, el nivel de esfuerzo y el agente de la sesión; antes de mover la conversación a una cuenta de otra organización, ccam te pregunta.
- Con `auto_switch = true` (Linux y macOS; experimental), ccam también hace el cambio por sí mismo: cuando Claude Code detiene la sesión en su límite y otra cuenta del mismo tipo tiene margen, ccam termina la sesión y continúa la conversación allí ([condiciones y costes](/features/next#auto-switch)).

Otras herramientas cambian el inicio de sesión de una sesión en curso; [ccam y otros selectores de cuentas](/compare/account-switchers) explica la diferencia.

## Qué tener en cuenta

- Ejecuta `ccam next --continue` en la misma carpeta en la que trabajabas; las conversaciones se guardan por carpeta de proyecto.
- Evita ejecutar dos cuentas sobre la misma conversación al mismo tiempo. ccam te avisa si lo intentas.
- Sin historial compartido, la conversación se copia a la otra cuenta y las dos copias siguen su propio camino; los puntos de control desde los que `/rewind` restaura el código se quedan con la primera cuenta.
- El primer mensaje en la cuenta nueva puede consumir más de su límite de lo habitual, porque esa cuenta todavía no tiene la conversación en caché. `ccam next --continue --brief` la empieza en cambio con un resumen breve.
- Si puedes cambiar de cuenta de esta forma depende de los términos que se aplican a cada cuenta.

Consulta [Límites de uso](/features/usage) y [Cambiar a la siguiente cuenta](/features/next) para ver todas las opciones.
