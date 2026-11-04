---
title: Cambiar a la siguiente cuenta
description: ccam next abre Claude Code con otra cuenta; --auto elige la de más margen, --continue lleva consigo la conversación y --fallback reejecuta trabajos -p.
updated: 2026-11-05
order: 8
translation: machine
sourceHash: "de1fd573f1ac"
---

`ccam next` abre Claude Code con otra cuenta sin que tengas que nombrarla. Úsalo cuando una cuenta llegue a su límite de uso. Hay cuatro formas de cambiar:

- `ccam next` a secas toma la cuenta que sigue a la que usaste por última vez, y `--auto` toma la que tiene más margen. `--continue` continúa la última conversación de esta carpeta en la cuenta nueva, compartan o no el historial las dos.
- Dentro de Claude Code, el [plugin de ccam](#inside-claude-code) te avisa cuando un límite detiene una sesión, y `/ccam:switch` pasa la conversación a otra cuenta.
- Para scripts que ejecutan `claude -p`, `ccam --fallback` hace el cambio por sí solo (consulta [Volver a ejecutar al llegar al límite](#run-again-at-the-limit)).
- Con [`auto_switch`](#auto-switch) activado (Linux y macOS), ccam también mueve una sesión interactiva que Claude Code detiene en su límite.

**Asientos de Team y Enterprise.** La selección automática de cuenta (`ccam next`, `--fallback`, `ccam pick`, la tecla `n` del panel y, una vez que lo activas, [`auto_switch`](/features/next#auto-switch)) también pasa de un asiento a otro dentro de una misma organización de Claude Team o Enterprise, lo que puede infringir los términos de Anthropic. Lee las [preguntas frecuentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) y consulta con tu organización y con Anthropic antes de depender de ella.

## Uso

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # continue this folder's last conversation on the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next -c --brief            # a new conversation that starts from a summary of the last one
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the most weekly quota left per day
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work CCAM_PROFILE=work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam imprime el cambio como `→ <from> → <to>` en la salida de error estándar, con el motivo de una elección de `--auto`, y luego inicia Claude Code con la cuenta nueva igual que `ccam <name>`. `<from>` es la cuenta que usaste por última vez o, con `--continue`, la cuenta a la que pertenece la conversación. Los argumentos que van después de `--` llegan a `claude` sin cambios.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--auto` | Elige la cuenta con sesión iniciada según el uso (ver más abajo), y omite las cuentas al 100 % | desactivada |
| `--max-usage <1-100>` | Con `--auto`: prefiere las cuentas por debajo de este uso | `auto_max_usage`, o 100 |
| `--strategy headroom\|expiring\|balance` | Con `--auto`: el menor uso de 5 horas, el restablecimiento de 5 horas más próximo o la mayor cuota semanal restante por día | `auto_strategy`, o `headroom` |
| `--cross-org` | Con `--auto`: elige también cuentas de otra organización (ver más abajo) | `cross_org`, o desactivada |
| `--wait[=<duration>]` | Con `--auto`: cuando ninguna cuenta tiene margen, espera a que se restablezca la primera, como máximo este tiempo. El valor va después de `=`: `--wait 2h` se rechaza con la sugerencia `--wait takes its value after '=': --wait=2h` | desactivada; `5h` si se da sola |
| `-c`, `--continue` | Continúa la última conversación de esta carpeta en la cuenta nueva (ver más abajo) | desactivada |
| `--brief` | Con `--continue`: inicia con la cuenta nueva una conversación nueva a partir del resumen que ccam hace de la anterior, en lugar de reanudarla entera ([ver más abajo](#brief)) | desactivada |
| `-- <claude args>` | Argumentos que se pasan a `claude` después de los de ccam | ninguno |
| `--dry-run` | Se admite: imprime `→ <from> → <to>` y la línea de comandos que ejecutaría, sin iniciar Claude Code ni registrar un inicio | desactivada |

`--auto` y `--continue` se pueden combinar. Con `--continue`, los argumentos de `--resume` se colocan antes de los tuyos.

Las opciones propias de Claude Code, como `--dangerously-skip-permissions` o `--model`, van después de `--`. A diferencia de `ccam <name>`, que pasa todo lo que va después del nombre, `ccam next` lee sus argumentos por sí mismo, así que sin el `--` se detiene con `ccam: unknown flag: --dangerously-skip-permissions`.

## Cómo funciona

Cada inicio a través de ccam escribe la hora en `~/.local/state/ccam/last-used/<name>`. `ccam next` lee estos archivos, comprueba qué cuentas tienen sesión iniciada y elige.

1. ccam crea un `next.Candidate` por cada cuenta de `config.toml` (`cli.candidates`): la cuenta, la hora de su último inicio (`launcher.LastUsed`) y si tiene un inicio de sesión válido (`creds.Has`: un inicio de sesión de Claude en `.credentials.json` o, en macOS, en el Llavero, que no haya vencido del todo ni lo haya borrado Claude Code). Una carpeta que usa una clave de API en su lugar (un inicio de sesión de Console, o una cuenta agregada con `ccam add --api`) no tiene límites de plan y nunca se elige.
2. Sin `--auto`, ccam ordena las cuentas por número de posición y busca la que se inició más recientemente (`next.Rotate`). Esa es la cuenta de la que te cambias. Si todavía no se ha iniciado ninguna cuenta, es la posición más baja.
3. A partir de ahí avanza en el orden de las posiciones, volviendo de la última posición a la primera, y devuelve la primera cuenta que tiene sesión iniciada.
4. ccam prepara e inicia la cuenta elegida exactamente como lo haría `ccam <name>` (`launcher.LaunchWith`): comprueba los enlaces de la cuenta, aplica el uso compartido de MCP, registra la sesión y reemplaza ccam por `claude` mediante `syscall.Exec`.

El recorrido empieza en la última cuenta usada y se detiene en la primera cuenta con sesión iniciada que viene después:

<!-- motion:next-rotate -->

`work` fue la que se inició más recientemente, así que el recorrido empieza después de ella. `side` no tiene sesión iniciada y se omite, y se elige `lab`. Si `lab` tampoco tuviera sesión iniciada, el recorrido volvería al principio hasta `personal` (`% len` en el diagrama). La rotación nunca devuelve la cuenta desde la que empezaste: si no hay otra cuenta con sesión iniciada, se detiene con un error.

### --auto

`--auto` decide por el uso en lugar del orden de las posiciones. ccam se queda con las cuentas que tienen sesión iniciada, consulta el uso de todas a la vez igual que `ccam usage` (`Fetcher.FetchAll`) y elige según las reglas de abajo (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

Una cuenta se juzga por sus ventanas de 5 horas y semanal y, cuando `--model` está entre los argumentos para Claude Code (después de `--`), también por la ventana semanal de ese modelo (`opus`, `sonnet`, `haiku` o `fable` en el nombre). ccam reparte las cuentas en tres niveles (`Policy.Choose()`):

1. **Con margen**: todas las ventanas por debajo de `auto_max_usage` (`config.toml`, 100 por defecto). La estrategia elige entre ellas: `headroom` (predeterminada) toma el menor uso de 5 horas y, entre cuentas con el mismo uso de 5 horas, la que menos ha gastado de su semana; `expiring` toma la ventana de 5 horas que se restablece antes, para que su cuota restante se use antes de perderse; `balance` toma la cuenta con más cuota semanal restante por día hasta que se restablezca su semana, y cuenta en su contra cada sesión de ccam en ejecución, para que la cuota semanal no quede sin usar y las sesiones en paralelo se repartan.
2. **Por encima**: por debajo del 100 %, pero en `auto_max_usage` o por encima. Solo se usan cuando ninguna cuenta tiene margen; la de menor uso de 5 horas va primero y, en caso de empate, la que menos ha gastado de su semana.
3. **Llenas**: una ventana al 100 %, o un uso de 5 horas desconocido. Nunca se eligen.

Una ventana cuenta como 100 % hasta su restablecimiento cuando Claude Code rechazó una solicitud en ella: el hook del [plugin de ccam](#inside-claude-code) lo registra (`~/.local/state/ccam/usage/<name>.limit.json`), digan lo que digan las cifras guardadas. Así, una cuenta que acaba de llegar a su límite no se vuelve a elegir antes de su restablecimiento, aunque sus cifras sean antiguas.

Los ajustes para todas las cuentas van al principio de `config.toml`, antes del primer `[[group]]`, `[[profile]]` o `[[folder]]`. Un ajuste debajo de uno de ellos pertenece a esa tabla: ccam avisa de ello en cada inicio, y los comandos que guardan la configuración se niegan hasta que se mueva.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organizaciones.** Las elecciones automáticas se quedan entre cuentas del mismo tipo que la que dejas: asientos de la misma organización de Team o Enterprise, o planes personales (Pro, Max), para que el trabajo de una cuenta de la empresa no continúe en una privada ni al revés. ccam conoce la organización de cada cuenta por su inicio de sesión (`ccam login` la registra; en los inicios de sesión más antiguos se lee de `.claude.json`); una cuenta cuya organización ccam no puede determinar se puede emparejar con cualquier cuenta. Dos carpetas con la sesión iniciada en la misma cuenta de Claude comparten sus límites, así que una nunca releva a la otra. `--cross-org`, o `cross_org = true` en `config.toml`, quita la regla de la organización, pero después de [`/ccam:switch`](#inside-claude-code) un cambio a otra organización sigue esperando tu sí.

Cuando no se puede elegir nada, ccam sale con el código 75 (`EX_TEMPFAIL`), para que un script distinga "inténtalo más tarde" de un fallo. Con `--wait`, en cambio, espera hasta que se restablezca la ventana que bloquea a la primera cuenta llena (30 segundos después de la hora de restablecimiento, cuando el endpoint muestra la nueva ventana), imprime `ccam: no account has room; waiting until 15:10` y vuelve a elegir; cuando ninguna cuenta dice cuándo se restablece, vuelve a mirar a los 10 minutos. Si nada se libera antes de que acabe la espera, sale con 75.

La consulta sigue las reglas de `ccam usage`: las cifras con menos de 60 segundos se reutilizan, la consulta de cada cuenta tiene un tiempo de espera de 3 segundos, una cuenta que no responde, o cuyo token de acceso venció, se juzga por sus cifras guardadas, y un registro de límite tiene prioridad sobre ellas. Una ventana guardada cuya hora de restablecimiento ya pasó cuenta como 0 %, así que no se omite una cuenta cuya ventana se restableció desde entonces. `--auto` nunca renueva tokens. A diferencia de la rotación simple, `--auto` puede elegir la cuenta que usaste por última vez si sigue teniendo el mayor margen.

### --continue

1. ccam busca la última conversación de esta carpeta en todas las cuentas: de los archivos `<session>.jsonl` de `projects/<folder>/` (el nombre de la carpeta es la carpeta de trabajo, con los enlaces resueltos, con cada carácter que no sea una letra o un dígito convertido en `-`, como la nombra Claude Code), el que tiene el prompt o la respuesta más reciente. Las sesiones sin ninguno, como una sesión que solo ejecutó `/mcp`, `/usage` o `/status`, o un archivo vacío, se omiten. Si una sesión dejó una nota (un límite de uso o `/ccam:switch`, ver más abajo) y ninguna conversación tiene un prompt o una respuesta más de 2 minutos posterior, se usa exactamente esa sesión; las filas que Claude Code escribe al hacer `/exit` no cuentan.
2. La cuenta a la que pertenece es la cuenta de la que te cambias, y la cuenta nueva se elige como se describe arriba, dentro de la misma organización. Cuando varias cuentas comparten ese historial, pertenece a aquella cuyo `.claude.json` dice que ejecutó la conversación por última vez (`lastSessionId`); si no, a la que se usó por última vez.
3. Si la cuenta nueva comparte ese historial (la misma carpeta de grupo), ccam inicia `claude --resume <session id>`, justo esa conversación (`--continue` abriría la más reciente de la carpeta, sea cual sea).
4. Si no, inicia `claude --resume <the conversation's file> --fork-session`: Claude Code copia la conversación completa, con un nuevo ID de sesión, en el historial propio de la cuenta nueva, y el original queda como estaba. A partir de entonces, `--continue` de la cuenta nueva la encuentra.
5. A menos que indiques `--agent` después de `--`, ccam también pasa el agente con el que se ejecuta la conversación: el ajuste `agent` del `.claude/settings.local.json` o del `.claude/settings.json` de la carpeta o del `settings.json` de la cuenta y, si no, el agente que nombra la última fila `agent-setting` de la conversación. Claude Code no recupera el agente cuando copia una conversación (paso 4), y seguiría como su agente predeterminado, con todas las herramientas. Cuando la cuenta nueva no comparte los agentes, plugins y configuración de la anterior (las mismas carpetas `agents` y `plugins` y el mismo `settings.json`, a través de un grupo), ccam agrega una línea después de la del cambio: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

ccam indica el tamaño de la conversación: la cuenta nueva la vuelve a leer entera con su primera respuesta. Las cachés de prompts nunca se comparten entre organizaciones, y cada plan personal es una organización propia, así que cuando las dos cuentas pertenecen a organizaciones distintas ccam agrega `prompt caches are not shared between organizations, so the first reply reads it all again`. Cuando la cuenta anterior tiene la sesión iniciada en servidores MCP en los que la nueva no, ccam los nombra: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. Compartir `history` en un grupo evita la copia, mantiene una sola conversación y comparte también los puntos de control desde los que `/rewind` restaura el código (`file-history/`).

### --brief

`--continue` reanuda la conversación completa, y la cuenta nueva la vuelve a leer entera con su primera respuesta. `ccam next --continue --brief` inicia en cambio una conversación nueva, a partir de un resumen que ccam escribe de la anterior sin consultar a un modelo:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

El resumen contiene, en este orden y solo donde hay algo que decir:

1. el título de la conversación;
2. la ruta de su transcripción, para que Claude la consulte cuando necesite un detalle;
3. tu primera solicitud, cuando hubo más de cuatro;
4. el resumen que escribió Claude Code la última vez que compactó la conversación;
5. tus últimas cuatro solicitudes, con los comandos `!` y los comandos de barra tal como los escribiste;
6. la última respuesta de Claude;
7. la lista de tareas;
8. los 10 archivos que Claude cambió por última vez;
9. el `git status --short --branch` de la carpeta, hasta 800 caracteres.

La lista de tareas conserva como máximo 10 elementos: primero se descartan los terminados, luego los pasos más lejanos, para que se queden el paso actual y los siguientes. Las partes largas se recortan por el medio, y el resumen completo se queda por debajo de 16.000 caracteres.

ccam lo guarda como `~/.local/state/ccam/brief/<session>.md`, legible solo por ti, y el siguiente `--brief` elimina los resúmenes de más de 7 días. Claude Code se inicia con un prompt y adjunta el archivo que menciona: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` Un prompt que pongas después de `--` ocupa el lugar de la última frase (`Then: <your prompt>`), y tus demás opciones después de `--` lo siguen.

Claude Code termina una mención en un `"` y la corta en un `#`, así que cuando la ruta del resumen contiene alguno de los dos, ccam lo menciona por su ruta relativa a la carpeta de trabajo; cuando esa ruta también contiene uno, el prompt nombra el archivo, y Claude lo lee con su herramienta Read, que pregunta antes porque el archivo está fuera de las carpetas de trabajo (en una ejecución `-p` en modo Manual, Claude Code deniega la lectura y Claude sigue sin el resumen).

`ccam --dry-run next --continue --brief` imprime el resumen y no guarda nada. La conversación anterior se queda como estaba, y la nueva solo sabe lo que dice el resumen, así que usa `--brief` cuando la tarea importe más que cada detalle de cómo llegaste ahí. Sin `--continue`, ccam se detiene con `--brief works with --continue`; cuando no puede leer la conversación o no encuentra nada que resumir, se detiene e indica que `ccam next --continue` sin `--brief` la reanuda entera. El [relevo](#relay) y [`auto_switch`](#auto-switch) siempre reanudan la conversación completa. La opción `--brief` propia de Claude Code va después de `--`, como sus demás opciones.

## Dentro de Claude Code

El [plugin de ccam](/docs/install#claude-code-plugin) agrega dos hooks (Claude Code 2.1.78 o posterior):

- Cuando un turno se detiene en un límite de uso, anota la conversación y muestra una notificación de escritorio en las terminales que la admiten (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  Conoce todos los mensajes que Claude Code cuenta como límite de uso: los límites de 5 horas, semanal, simple y de modelo (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` para la asignación de un modelo), los créditos de uso y los límites de gasto, y una organización que detuvo la cuenta.

  Cuando Claude Code ha escrito la solicitud rechazada en la conversación, el hook también registra la ventana llena hasta su restablecimiento ([el registro de límite](/features/usage#the-limit-record)), salvo en una detención establecida por una organización, que sus administradores pueden levantar en cualquier momento. Una limitación breve del servidor o el límite de solicitudes de una clave de API no es un límite de uso y no se toca, y tampoco las sesiones desatendidas (ejecuciones `-p` y sesiones en segundo plano). Un subagente o un agente de un flujo de trabajo que se detiene en un límite de uso cuenta como la detención de su sesión, aunque la conversación principal solo lo esté esperando: Claude Code escribe esa solicitud rechazada en el propio archivo del agente, en la carpeta `subagents` de la sesión, y el hook la lee ahí. Cuando varios agentes se detienen a la vez, la sesión recibe una sola nota y una sola notificación.

  Claude Code 2.1.234 y posteriores también esperan en la sesión y continúan por sí solos tras un restablecimiento de 5 horas; la notificación te da la opción de cambiar de cuenta en su lugar, y con [`auto_switch`](#auto-switch) activado (Linux y macOS), ccam cambia por ti cuando otra cuenta tiene margen. Cuando ninguna otra cuenta puede tomar el relevo, porque ninguna otra cuenta pertenece a la misma organización de Team o Enterprise, o ninguna es otro plan personal (con `cross_org` activado, puede hacerlo cualquier otra cuenta de Claude), la notificación termina en cambio con `You have no other personal plans, so ccam cannot move this conversation.` o `You have no other accounts in Acme, so ccam cannot move this conversation.`
- `/ccam:switch [account]` anota la conversación y termina el turno antes de que se haga ninguna solicitud: escribe `/exit` y luego ejecuta `ccam next --continue` para continuarla en la cuenta que nombraste (o en la siguiente). En una sesión que todavía no tiene conversación, indica que no hay nada que continuar. Solo lo que escribes tú lo ejecuta: desde la versión 0.3.2 del plugin, Claude no puede ejecutar `/ccam:switch`, y un trabajo que Claude programa (CronCreate, ScheduleWakeup o `.claude/scheduled_tasks.json`) se lo pasa a Claude como texto plano.

  Cuando la cuenta que nombras está fuera de la organización de la conversación (un asiento de otra organización de Team o Enterprise, un plan personal en lugar de un asiento o al revés, o una cuenta con clave de API en cualquiera de los dos lados), `ccam next --continue` pregunta primero, diga lo que diga `cross_org`: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. Pregunta lo mismo antes de una cuenta que elige por sí mismo cuando no nombras ninguna, y que `cross_org = true` le permite tomar de otra organización: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. Solo `y` o `yes` la mueve; en Linux y macOS, primero se descartan las teclas pulsadas antes de la pregunta. Cualquier otra respuesta, o no tener una terminal en la que preguntar, deja la conversación y la nota donde están: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (después de una elección, `ccam next --continue` vuelve a elegir). Como en las elecciones automáticas, una cuenta cuya organización ccam no puede determinar encaja con cualquiera.

La nota vive en `~/.local/state/ccam/handoff/<account>/<session>.json`, un archivo por sesión, así que dos sesiones de una misma cuenta nunca sobrescriben la nota de la otra. Se conserva durante 12 horas y se usa una sola vez, cuando la siguiente sesión está a punto de iniciarse; un inicio que falla la deja para el siguiente intento. Pertenece a la carpeta en la que se inició la sesión, incluso después de un `cd` en el shell de Claude Code. `ccam rename` actualiza los nombres que contiene; si la cuenta que nombra se ha eliminado desde entonces, ccam lo indica y elige como de costumbre.

### Relevo

Con `relay = true` al principio de `config.toml`, ccam sigue ejecutándose junto a las sesiones interactivas en lugar de reemplazarse a sí mismo por Claude Code, y después de `/ccam:switch` continúa por sí solo la conversación de esa sesión (no actúa sobre una nota que haya dejado otra sesión): el hook termina Claude Code, y ccam ejecuta `ccam next --continue` con la nota de esa sesión. Si no se nombra ninguna cuenta, toma la que tiene más margen, contando el modelo de la conversación, y cuando `cross_org = true` permite que sea una cuenta de otra organización, hace la pregunta de arriba antes de mover allí la conversación. En Windows el hook no puede terminar Claude Code: escribe `/exit`, y ccam la continúa. Tampoco lo termina el hook para una cuenta fuera de la organización de la conversación: dice `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`, y ccam hace la pregunta de arriba en cuanto hayas escrito `/exit`. Eso cuesta un proceso de ccam en espera por sesión, y por eso está desactivado por defecto.

La sesión que toma el relevo conserva lo que reanudar una conversación no restaura: opciones como `--add-dir`, `--mcp-config`, las listas de herramientas, `--settings`, `--plugin-dir` y las opciones del prompt del sistema, el modo de permisos de la conversación, el nivel de esfuerzo de la sesión y el agente con el que se ejecuta la conversación (ver [`--continue`](#continue)). El propio Claude Code restaura `--model` y `--name` a partir de la conversación. Cualquier otra cosa (`--worktree`, `--debug`, `--remote-control`, opciones que ccam no conoce) se nombra en `ccam: not carried over: …`. El cambio nunca hace la pregunta del historial compartido, y lo que escribiste mientras se ejecutaba se descarta; su única pregunta es la que hace antes de una cuenta de otra organización.

El relevo solo continúa la conversación cuando Claude Code terminó como se lo pidió `/ccam:switch`: no cuando escribiste `/exit` antes, pulsaste Ctrl+C mientras terminaba o ccam recibió un hangup u otra señal, ni cuando la terminal ya no está. La nota se queda entonces, y `ccam next --continue` la continúa más tarde. Una sesión iniciada con el relevo de una versión anterior de ccam, o una desatendida (ejecuciones `-p` y sesiones en segundo plano), no la termina el hook; recibe en su lugar el mensaje de `/exit`. `claude remote-control` se ejecuta sin relevo.

### `auto_switch`

Con `auto_switch = true` al principio de `config.toml` (Linux y macOS), ccam mueve una sesión a otra cuenta por sí mismo cuando Claude Code la detiene en un límite de uso: el hook del plugin elige otra cuenta con margen y termina Claude Code, y ccam, que sigue ejecutándose junto a la sesión, continúa allí la conversación, como hace el [relevo](#relay) después de `/ccam:switch`. Está desactivado por defecto y es experimental.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam mueve una sesión solo cuando se cumplen todas estas condiciones:

1. ccam la inició en una terminal y sigue ejecutándose junto a ella. `auto_switch` activa el [relevo](#relay), así que esto se cumple para las sesiones interactivas que inicias después de activarlo. Una sesión iniciada antes, o con un ccam anterior, recibe una notificación que termina con ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` Las sesiones desatendidas (ejecuciones `-p` y sesiones en segundo plano), `claude remote-control` y Claude Code iniciado sin ccam nunca se mueven.
2. La cuenta tiene activado el [plugin de ccam](/docs/install#claude-code-plugin) y ejecuta Claude Code 2.1.214 o posterior, que pasa el ID de su proceso al hook (`CLAUDE_PID`).
3. Claude Code escribió la solicitud rechazada en la conversación, con la hora a la que se restablece, para el límite de 5 horas, semanal, de Opus o de Sonnet de esa cuenta. Un límite de gasto, el agotamiento de los créditos de uso, una detención establecida por una organización, la asignación de Fable (`You've reached your Fable limit.`), una limitación del servidor y el límite de solicitudes de una clave de API nunca mueven una sesión. El límite de un subagente o de un agente de un flujo de trabajo sí cuenta, aunque la conversación principal solo lo esté esperando; cuando es un límite semanal de Opus o de Sonnet, la sesión se mueve solo si la propia conversación se ejecuta con ese modelo.
4. ccam conoce el plan de la cuenta, por su `.claude.json` o por lo que registró `ccam login`. Si no, la notificación termina con `ccam does not know personal's plan, so auto_switch leaves this session.`
5. La sesión no está conectada a Remote Control. Con otra cuenta no volvería a conectarse a la app desde la que la manejas, así que se queda, y la notificación lo indica.
6. El límite se restablece dentro de más de 15 minutos. Para una espera más corta, ccam deja la sesión a Claude Code: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. No has escrito ningún prompt desde el límite.
8. Otra cuenta tiene margen. ccam elige como lo hace [`--auto`](#auto) (`auto_max_usage`, `auto_strategy`, el modelo de la conversación), entre cuentas del mismo tipo que la que dejas: otro asiento de la misma organización de Team o Enterprise, u otro plan personal (Pro, Max). `cross_org` y `--cross-org` no se aplican aquí. La cuenta debe:
   - tener la sesión iniciada con un inicio de sesión que Claude Code pueda renovar;
   - no tener la sesión iniciada en la misma cuenta de Claude;
   - no usar una clave de API;
   - tener un plan que ccam conozca;
   - no tener los créditos de uso activados, según lo que muestra la consulta de uso propia de ccam;
   - haber confiado en esta carpeta en Claude Code;
   - si la conversación se ejecuta con un agente (el ajuste `agent` o la conversación nombra uno), compartir los agentes de la cuenta que dejas: sus carpetas `agents` y `plugins` y su `settings.json` tienen que ser, cada uno, exactamente los mismos (a través de un grupo) o faltar en las dos.

   Cuando esas tres últimas condiciones dejan fuera a todas las cuentas con margen, la notificación las nombra: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` o `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` Una ventana con un [registro de límite](/features/usage#the-limit-record) cuenta como llena, y una cuenta a la que otra sesión se movió en los últimos 2 minutos va después de las demás. La elección tiene que hacerse en unos 6 segundos; cuando las consultas de uso tardan más, la sesión se queda. Con `usage_source = "claude"`, `auto_switch` decide solo con las cifras guardadas, y ccam no puede saber por la consulta de Claude Code ni por la línea de estado si los créditos de uso están activados, así que supone que están desactivados.

Entonces el hook termina Claude Code, y ccam, que sigue ejecutándose junto a él, muestra una notificación de escritorio en las terminales que la admiten, como `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` Imprime la primera línea de la salida de ejemplo de arriba y ejecuta `ccam next --continue` para esa sesión con la cuenta elegida, que imprime sus otras dos líneas e inicia Claude Code con las opciones, el modo de permisos, el nivel de esfuerzo y el agente de la sesión, como hace el [relevo](#relay). El primer mensaje de la nueva sesión es un texto fijo, lo único que ccam envía en tu nombre (un cambio después de `/ccam:switch` no envía nada):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

Después del límite de un subagente es este otro, porque los subagentes, los flujos de trabajo y los comandos en segundo plano se detuvieron con la sesión anterior:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

Si esa cuenta también llega a su límite, ocurre lo mismo allí. Los registros de límite impiden que la cadena vuelva a una cuenta antes de su restablecimiento, así que cuando todas las cuentas están llenas la sesión se queda en la última.

El hook sabe que una conversación se ejecuta con un agente por el ajuste `agent`, o por la fila que nombra el agente y que Claude Code escribe en la conversación de vez en cuando y al terminar la sesión. Si hasta ahora solo el `--agent` de la sesión lo nombra, y la cuenta elegida no comparte los agentes, `ccam next --continue` no inicia nada y conserva la nota, y una segunda notificación de escritorio lo indica: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

Cuando ninguna otra cuenta tiene margen, nada se mueve: Claude Code gestiona el límite como de costumbre, esperando en la sesión u ofreciendo sus opciones, y la notificación indica qué cuenta se libera primero, por ejemplo `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` Para los asientos dice `None of your other accounts in Acme has room now`. Para detener un cambio que ya empezó, pulsa Ctrl+C mientras Claude Code se cierra: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

Un cambio tiene costes. Claude Code se reinicia, lo que tarda unos segundos, y los shells en segundo plano, los servidores de desarrollo y los agentes que inició Claude se detienen con la sesión anterior. La cuenta nueva vuelve a leer la conversación completa con su primera respuesta, como con [`--continue`](#continue), y ccam indica su tamaño; un cambio nunca usa [`--brief`](#brief).

`auto_switch` nunca cambia antes de que Claude Code informe de un límite, nunca vigila tu uso en segundo plano y nunca copia, renueva ni transmite un inicio de sesión. En Windows se ignora, y `ccam doctor` lo indica; allí, con `relay = true`, `/ccam:switch` y `/exit` siguen moviendo una conversación a otra cuenta.

**Términos de Anthropic.** En los asientos de Team y Enterprise, `auto_switch` mueve tu trabajo a otro asiento de la misma organización. El centro de ayuda de Anthropic indica: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada asiento de un plan Team se asigna a una persona específica, y solo esa persona puede usarlo) y "Users can't route their work through another user's account to use that user's remaining usage" (los usuarios no pueden hacer pasar su trabajo por la cuenta de otro usuario para consumir el uso restante de ese usuario) ([planes Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); lo mismo para los [planes Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)).

Para los planes personales (Pro, Max), los términos de Anthropic no dicen si tus propias cuentas pueden tomar el relevo unas de otras automáticamente al llegar a un límite. Los [Términos para consumidores](https://www.anthropic.com/legal/consumer-terms) prohíben acceder a los Servicios "through automated or non-human means, whether through a bot, script, or otherwise" (por medios automatizados o no humanos, ya sea mediante un bot, un script o de otro modo) salvo que Anthropic lo permita expresamente, y un cambio envía en tu nombre uno de los mensajes fijos de arriba. La [Política de uso](https://www.anthropic.com/legal/aup) que entra en vigor el 12 de noviembre de 2026 incluye "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (eludir o interferir intencionadamente con las capacidades, restricciones o salvaguardas establecidas en nuestros productos) entre lo que prohíbe.

ccam no afirma que `auto_switch` cumpla estos términos. Activarlo es decisión tuya, y el riesgo también: Anthropic puede suspender una cuenta, o el acceso de una organización de Team o Enterprise, cuando considera que se han infringido los términos.

Mientras `auto_switch` está activado, `ccam doctor` advierte de una cuenta sin el plugin de ccam y de las versiones de Claude Code anteriores a 2.1.214, y para los asientos repite la regla de Anthropic sobre los asientos ([comprobaciones](/features/doctor)).

## Volver a ejecutar al llegar al límite

`ccam next` es para una sesión interactiva que se ha detenido. Una ejecución no interactiva, `claude -p` en un script o una tarea programada, no puede esperarte, así que `--fallback` antes del nombre de la cuenta permite que ccam cambie por ella:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` solo funciona junto con `-p` o `--print`; si no, ccam se detiene y remite a `ccam next --auto --continue`.
2. La entrada estándar por tubería se lee una vez, hasta 10 MiB, y se entrega a cada intento. Como Claude Code, ccam espera como máximo 3 segundos a que empiece (salvo con `--input-format stream-json`); si una tubería sigue abierta pero no envía nada, la ejecución no recibe entrada y ccam lo indica en la salida de error estándar.
3. ccam inicia `claude` como proceso hijo en lugar de con `exec`, retiene su salida estándar y deja pasar la salida de error estándar a medida que llega.
4. Cuando `claude` sale con 0, o falla por otro motivo, la salida retenida y el código de salida se transmiten sin cambios. Lo mismo pasa con una ejecución terminada por una señal (un código de salida superior a 128): se le pidió que se detuviera, y no se vuelve a ejecutar.
5. Un fallo cuenta como límite de uso cuando el mensaje de Claude Code lo indica, cualquiera de los mensajes de límite de uso que conoce el [hook](#inside-claude-code), también dentro de la salida `--output-format json` o `stream-json`, o cuando el uso de la cuenta muestra su ventana de 5 horas, semanal o de `--model` al 100 % (incluido un registro de límite).
6. Entonces se descarta la salida de ese intento y el mismo comando se ejecuta con otra cuenta con sesión iniciada: la de menor uso de 5 horas, y después las cuentas cuyo uso no se puede consultar, por orden de número. Las cuentas al 100 % se omiten, las cuentas con un inicio de sesión con clave de API nunca se prueban, y cada cuenta se prueba como máximo una vez.
7. Con `-c`, `--continue`, `-r` o `--resume` entre los argumentos, solo se prueban las cuentas del mismo grupo de historial, para que se pueda encontrar la conversación.
8. Solo se prueban cuentas de la misma organización, igual que con `--auto` (`cross_org = true` quita esta restricción).

La siguiente cuenta se elige con las mismas reglas que `--auto` (`auto_max_usage`, `auto_strategy`, la ventana de `--model`). Cuando no queda ninguna cuenta, ccam imprime la salida del último intento, indica que ninguna cuenta con sesión iniciada tiene margen y sale con el código 75. Con `--wait` antes del nombre de la cuenta (`ccam --fallback --wait=2h work -p …`), en cambio, espera a que se restablezca la primera cuenta y después puede volver a probar todas las cuentas; siempre espera al menos una vez antes, para que una cuenta que acaba de llegar a su límite no se vuelva a ejecutar de inmediato. La duración limita el tiempo total de espera de la ejecución, no cada espera. Escrito con un espacio (`--wait 30m`), se rechaza con la misma sugerencia que `ccam next`. Un `--` inicial antes de los argumentos de Claude Code funciona como en cualquier inicio (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` imprime el comando del primer intento.

Las sesiones interactivas no cambian por sí solas a menos que actives [`auto_switch`](#auto-switch). Con `ccam statusline install`, la línea de estado añade `→ ccam next --auto --continue` cuando el uso de 5 horas o semanal llega al 90 % ([Límites de uso](/features/usage)).

## Seguridad y casos límite

- Las cuentas sin inicio de sesión nunca se eligen. Si no hay otra cuenta con sesión iniciada, `ccam next` a secas se detiene con `no logged-in account other than <name>`.
- Con `--continue` y ninguna conversación en la carpeta: `no account has a conversation in <folder> to continue`.
- Cuando la regla de la organización no deja ninguna cuenta: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Con `--auto` y ninguna cuenta con sesión iniciada: `no logged-in account to pick`. Cuando todas las cuentas están llenas o su uso es desconocido: `` no account with known usage and headroom (see `ccam usage`) ``. Ambos casos terminan con el código de salida 75.
- `--max-usage`, `--strategy`, `--wait` y `--cross-org` solo funcionan con `--auto`, y se rechaza un valor fuera de 1-100 o una estrategia desconocida.
- Una bifurcación copia la conversación, no el resto del estado de la cuenta: los puntos de control de `file-history/` de la cuenta original no la acompañan, salvo que el historial se comparta.
- `--auto` hace solicitudes de red incluso con `--dry-run`. Solo lee, y nunca escribe archivos de tokens.
- La cuenta elegida se inicia por la misma ruta que `ccam <name>`, así que se aplican las mismas comprobaciones, incluida la advertencia cuando otra sesión de Claude Code que comparte el mismo historial ya se está ejecutando en esta carpeta.
- `last-used` solo se escribe cuando Claude Code se inicia de verdad. `--dry-run` no lo toca, así que una ejecución de prueba no hace avanzar la rotación.
- `--fallback` retiene la salida estándar hasta que termina un intento, así que una salida continua como `--output-format stream-json` llega de una sola vez al final. Un cliente que transmite entrada `stream-json` y espera respuestas antes de enviar más no puede usar `--fallback`. Los archivos que el primer intento cambió antes de llegar al límite siguen cambiados; el siguiente intento empieza el prompt de nuevo.

## Relacionado

- [Límites de uso](/features/usage)
- [Límites de uso y cambio de cuenta](/docs/usage-limits)
- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Qué hacer cuando Claude Code llega al límite de uso](/guides/usage-limit-workflow)
