---
title: Límites de uso
description: ccam usage muestra a la vez los límites de 5 horas y semanal de todas las cuentas, y ccam statusline pone el nombre de la cuenta y su uso en Claude Code.
updated: 2026-11-01
order: 7
translation: machine
sourceHash: "54b296f7da35"
---

`ccam usage` muestra, para cada cuenta, cuánto se ha usado de cada límite: el límite de 5 horas, el límite semanal y el límite semanal propio de cada modelo (como Fable). Cada límite tiene una barra y el tiempo que falta para que se restablezca. Úsalo antes de empezar una tarea larga, o cuando Claude Code te avise de que un límite está cerca, para ver qué cuenta tiene margen. `ccam statusline` es el comando de línea de estado para Claude Code: muestra el nombre de la cuenta y su uso dentro de Claude Code y guarda las cifras para `ccam usage`. Para los tokens y costes estimados que hay detrás del uso, `ccam ccusage` ejecuta [ccusage](https://ccusage.com) sobre las conversaciones de todas las cuentas ([más abajo](#tokens-and-costs-ccam-ccusage)).

## Uso

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts that are not running
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

Cada cuenta tiene una línea con su número, su nombre y la antigüedad de sus cifras; `(cached)` significa que vienen de un archivo guardado, no del endpoint. Debajo, una línea por límite muestra una barra de 20 celdas, el porcentaje y el tiempo hasta que se restablece la ventana: `5-hour`, `Weekly` (todos los modelos) y luego la ventana semanal propia de cada modelo, como `Fable`. Una ventana que la cuenta no tiene muestra `-`. Un porcentaje se vuelve amarillo al 80 % y rojo al 100 %.

Si Claude Code rechazó una solicitud porque una ventana estaba llena, esa ventana muestra 100 % hasta que se restablece, digan lo que digan las cifras guardadas (consulta [el registro de límite](#the-limit-record)). En una terminal todas las cuentas aparecen de inmediato como `fetching…`, y sus líneas se dibujan a medida que responden.

La última línea, que solo aparece si tienes más de una cuenta, nombra la cuenta que elegiría [`ccam next --auto`](/features/next). Sigue las mismas reglas:

- `auto_max_usage` y `auto_strategy`;
- las sesiones en ejecución;
- solo cuentas de la organización de la última cuenta usada, salvo que `cross_org` esté activado;
- se omite toda ventana al 100 %, así que nunca se nombra una cuenta con la semana llena.

Con la estrategia predeterminada `headroom` dice `→ Most headroom: <name>`; con `expiring` o `balance` da el motivo, como `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Asientos de Team y Enterprise.** La selección automática de cuenta (`ccam next`, `--fallback`, `ccam pick`, la tecla `n` del panel y, una vez que lo activas, [`auto_switch`](/features/next#auto-switch)) también pasa de un asiento a otro dentro de una misma organización de Claude Team o Enterprise, lo que puede infringir los términos de Anthropic. Lee las [preguntas frecuentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) y consulta con tu organización y con Anthropic antes de depender de ella.

Una cuenta sin sesión iniciada muestra `not logged in → ccam login <name>`, y una cuenta con clave de API muestra `no limits (API)`. Un token de acceso vencido muestra ``access token expired; `ccam <name>` renews it when it starts``, encima de las cifras guardadas, si las hay: Claude Code lo renueva en el siguiente inicio. Un inicio de sesión que venció del todo, porque su token de actualización venció o Claude Code lo borró después de que el servidor lo rechazara, muestra `` login expired → `ccam login <name>` ``.

Después de `ccam statusline install`, Claude Code muestra una línea como `[work] · Opus · 5h 24% · 7d 41%`.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--json` | Imprime los informes en JSON: `name`, `slot`, `status`, `limits` (con `credits` cuando la cuenta tiene activados los créditos de uso), `source`, `at`, `error`, y `limited` y `limited_windows` cuando un registro de límite pone ventanas al 100 %. Sin cuentas, imprime `[]` | desactivada |
| `--watch[=N]` | Vuelve a dibujar cada `N` segundos hasta que pulses Ctrl+C. `--watch` solo, o un valor de 0 o menos, significa 60. El valor va después de `=`: `--watch 10` se rechaza con una sugerencia | desactivada |
| `-v`, `--verbose` | Añade una línea `source:` bajo cada cuenta: de dónde salieron sus cifras (el endpoint de uso, Claude Code, la línea de estado, un archivo guardado o el límite de uso que informó Claude Code). Cuando un registro de límite tiene prioridad sobre las cifras, sigue una línea como `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` | desactivada |
| `--no-cache` | Ignora la caché de 60 segundos y vuelve a consultar el endpoint | desactivada |
| `--source live\|claude\|cache` | `live` consulta el endpoint con el token de la cuenta; `claude` hace que Claude Code lo consulte (ver más abajo); `cache` solo muestra las cifras guardadas | `usage_source`, o `live` |
| `--refresh` | Renueva los tokens de acceso vencidos, solo de las cuentas que no se están ejecutando. Una renovación puede tardar hasta 30 segundos | desactivada |
| `--dry-run` | Se acepta: `ccam usage` no cambia nada salvo su propia caché. Se rechaza junto con `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | desactivada |

Cualquier otro valor de `--source` falla con `--source must be live, claude or cache`. `ccam usage` no acepta argumentos.

### El uso desde el propio Claude Code

Con `usage_source = "claude"` al principio de `config.toml` (antes del primer `[[group]]`, `[[profile]]` o `[[folder]]`), o `--source claude`, ccam no lee en absoluto el token de inicio de sesión. Para cada cuenta inicia `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` en la carpeta de la cuenta y envía por la entrada estándar la solicitud `get_usage` de Claude Code (Claude Code 2.1.270 o posterior). Claude Code consulta el uso con su propio inicio de sesión y responde con las ventanas de 5 horas, semanal y por modelo.

`--safe-mode` impide que se inicien los hooks, los servidores MCP y los plugins de la cuenta, y no se guarda ninguna sesión. Cada cuenta necesita un inicio de Claude Code (alrededor de un segundo), así que esta fuente es más lenta que `live`. Las cuentas se consultan a la vez, cada una con un tiempo de espera de 20 segundos, y la caché de 60 segundos se aplica como siempre.

Un inicio de sesión cuyo token de acceso venció, o vence en menos de 5 minutos, recibe 50 segundos, ya que Claude Code lo renueva al iniciar y una renovación cortada puede perder el inicio de sesión. Si ccam abandona una consulta, le pide a Claude Code que termine con SIGTERM, para que Claude Code acabe una renovación que haya empezado (Windows no tiene SIGTERM, así que allí el proceso se termina).

Un inicio de sesión que venció del todo muestra `` login expired → `ccam login <name>` `` sin iniciar Claude Code. Claude Code marca la solicitud como experimental, así que si deja de responder, ccam muestra las cifras guardadas y `ccam doctor` lo informa. Cuando Claude Code responde que su propia consulta falló, ccam muestra las cifras guardadas con `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`, no un plan sin límites.

Con `--dry-run`, ccam no inicia Claude Code, porque Claude Code escribe en la carpeta de la cuenta; en su lugar muestra las cifras guardadas. `usage_source` acepta `live` o `claude`; `cache` solo es una opción de `--source`.

`ccam statusline` lee de la entrada estándar el JSON de la línea de estado de Claude Code. Su única opción es `--wrap <command>`: un comando de línea de estado existente a cuya primera línea de salida se le antepone el nombre de la cuenta. `ccam statusline install` no acepta opciones. Ambos rechazan `--dry-run` con `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Tokens y costes: `ccam ccusage`

`ccam usage` muestra cuánto ha usado una cuenta de cada límite. [ccusage](https://ccusage.com), en cambio, contabiliza los tokens y costes estimados en los registros de conversaciones de Claude Code. Por sí solo busca los registros de Claude Code únicamente en `~/.claude` y `~/.config/claude` (o en las carpetas que indique `CLAUDE_CONFIG_DIR`), así que la única cuenta de ccam que puede ver es una importada en su ubicación en `~/.claude`. `ccam ccusage` lo ejecuta con `CLAUDE_CONFIG_DIR` apuntando a las carpetas de todas las cuentas que tienen conversaciones, de modo que un solo informe las cubre todas:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam lee primero sus propias opciones: `-a` (`--account`) acepta nombres o números de cuenta, separados por comas o repitiendo la opción, y `--dry-run` imprime el comando en lugar de ejecutarlo. A partir del primer argumento que no es una de ellas, todo pasa a ccusage tal como lo escribiste, salvo `--config` para `bun x` y `npx` (ver más abajo); `ccam ccusage -- --help` muestra la ayuda del propio ccusage. ccam ejecuta `ccusage` si está instalado; si no, `bun x ccusage`, y si tampoco, `npx ccusage@latest`; estos dos últimos lo descargan de npm. Ambos se inician en la carpeta propia de ccam (`~/.local/state/ccam/ccusage`), no en la carpeta en la que estás, así que un proyecto que esté allí no puede elegir lo que ejecutan con su `node_modules/.bin`, un `bin` en su `package.json` o un registro en su `.npmrc`; `ccam --dry-run ccusage` muestra el `cd` a esa carpeta. Cuando ninguno está disponible, ccam te indica que instales Node.js, Bun o ccusage. En Linux y macOS, ccusage ocupa el lugar de ccam; en Windows, ccam lo espera y termina con el código de salida de ccusage. En WSL, ccam se salta un `ccusage`, `bun` o `npx` que esté en una unidad de Windows, porque se ejecutaría como programa de Windows sin las carpetas de las cuentas.

- El informe cubre las cuentas cuyas carpetas tienen conversaciones (una carpeta `projects`). Una conversación que varias cuentas comparten a través de un grupo se contabiliza una sola vez.
- Con `-a`, una cuenta que comparte su historial a través de un grupo trae consigo las conversaciones de todo el grupo, y ccam dice de quiénes son: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- `blocks` muestra ventanas de 5 horas, y cada cuenta tiene las suyas, así que indica una sola cuenta para él, como arriba.
- ccusage también informa de las demás herramientas de programación que encuentra, como Codex, salvo que nombres una: `ccam ccusage claude daily` muestra solo Claude Code.
- Mientras ccam lo ejecuta, ccusage lee su configuración (`ccusage.json`) solo desde `./.ccusage` y las carpetas de las cuentas. `--config` indica otro archivo: `ccam ccusage daily --config ~/.config/claude/ccusage.json`. Como `bun x` y `npx` se inician en la carpeta propia de ccam, ccam les pasa `./.ccusage/ccusage.json` con `--config` cuando ese archivo existe y no indicas ninguno, y pasa un `--config` relativo como ruta completa. Entonces ccusage lee otras rutas relativas, como `--pi-path`, desde la carpeta de ccam, así que indícalas completas.
- ccusage leería una carpeta cuya ruta tiene una coma como si fueran dos carpetas, así que ccam la rechaza y nombra la cuenta; deja esa cuenta fuera con `-a`.
- `ccusage` es un comando de ccam, así que una cuenta llamada ccusage se inicia con `ccam run ccusage`; `ccam doctor` lo indica.

## Cómo funciona

Las cifras de `ccam usage` salen de dos lugares: el endpoint de uso no documentado de Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, al que se llama con el inicio de sesión propio de cada cuenta, y los valores que Claude Code le pasa a su comando de línea de estado. Todo esto está en `internal/usage`.

1. ccam carga `config.toml` y crea un `usage.Account` por cuenta (`cli.usageCmd`): su carpeta `~/.local/share/ccam/profiles/<name>/`, sus marcadores de sesión (junto con los registros propios de Claude Code en la carpeta de la cuenta, indican si está en ejecución) y si se importó en su ubicación.
2. ccam consulta todas las cuentas a la vez, con una goroutine por cuenta (`Fetcher.FetchAll`). Todas comparten un `http.Client`, así que las conexiones al endpoint se reutilizan. La conexión puede tardar 2 s, y cada cuenta tiene 3 s en total. El cliente no sigue redirecciones, ya que una redirección reenviaría el token a otro sitio, y lee como máximo 1 MB de la respuesta.
3. A medida que termina cada cuenta, su fila se vuelve a dibujar. Una cuenta que falla o no responde a tiempo recurre a las cifras guardadas, así que una cuenta lenta nunca retrasa a las demás.

En el diagrama, cada cuenta se consulta en su propia goroutine, y la que agota el tiempo muestra en su lugar sus cifras guardadas:

<!-- motion:usage-fetch -->

`personal` y `work` responden a tiempo y muestran cifras en vivo. `side` tarda más de 3 s, así que ccam muestra en su lugar sus cifras guardadas (`fallback()`). Después, ccam dibuja un bloque por cuenta (`usage.Table()`) y elige la cuenta que se nombra en la última línea (`usage.Best()`).

Para una cuenta, ccam hace estas comprobaciones en orden (`Fetcher.Fetch`):

<!-- motion:usage-cache -->

1. Con `--source cache`, va directamente a `fallback()`.
2. Salvo que se indique `--no-cache`, un resultado en vivo en `~/.local/state/ccam/usage/<name>.json` medido hace menos de 60 segundos (`CacheTTL`) se usa, evaluado a la hora actual: una ventana que se restableció desde entonces muestra 0 %, y un archivo con fecha posterior a la del reloj no cuenta como reciente. Ejecutar `ccam usage` dos veces seguidas cuesta una solicitud por cuenta.
3. ccam lee `.credentials.json` (o el Llavero de macOS) sin cambiarlo y decodifica solo la clave `claudeAiOauth` (`readCredentials()`). Si no hay un inicio de sesión de suscripción, se muestra `not logged in`, o `no limits (API)` cuando `.claude.json` contiene una clave de API. Un inicio de sesión cuyo token de actualización venció, o que Claude Code borró, es `` login expired → `ccam login <name>` `` y no se consulta.
4. Un token de acceso vencido no se usa. Sin `--refresh`, la cuenta muestra sus cifras guardadas con ``access token expired; `ccam <name>` renews it when it starts``.
5. En caso contrario, ccam llama al endpoint. Un 401 o un 403 cuenta como inicio de sesión vencido (`login expired: authentication failed (401)`), cualquier otro estado o un cuerpo que no puede leer cuenta como error, y en ambos casos se recurre a las cifras guardadas. Una consulta terminada con Ctrl+C dice `cancelled`. Un 429 significa que el endpoint ha recibido demasiadas solicitudes para esta cuenta: ccam escribe `<name>.pause` con la hora para volver a preguntar (`Retry-After`, o 5 minutos, como máximo una hora), muestra las cifras guardadas con `rate limited by the usage endpoint; asking again after 15:10` y no llama al endpoint para esa cuenta hasta entonces, ni siquiera con `--no-cache`. Reintentar de inmediato solo prolongaría el límite de solicitudes.
6. Una respuesta correcta se escribe en `<name>.json` con la hora en que se midió.

Para las cifras guardadas, ccam toma la consulta guardada más reciente (`<name>.json`, o `<name>.claude.json` de las consultas de Claude Code) y, cuando `<name>.statusline.json` se escribió después, superpone las ventanas de 5 horas y semanal de la línea de estado (`fallback()`); las ventanas por modelo, que la línea de estado no tiene, se quedan como las guardó la consulta. Antes de mostrar cifras guardadas, comprueba la hora de restablecimiento de cada ventana: una ventana que se restableció después de guardar las cifras muestra 0 %, porque el porcentaje anterior ya no es válido. Por eso `personal`, más arriba, muestra `0%` sin hora de restablecimiento.

El endpoint tiene dos formatos de respuesta. El más nuevo tiene un arreglo `limits` con entradas `session` (5 horas), `weekly_all` y `weekly_scoped` (por modelo). El más antiguo tiene los objetos `five_hour`, `seven_day`, `seven_day_opus` y `seven_day_sonnet`. ccam lee ambos (`ParseLive`): el arreglo `limits` tiene prioridad, y los objetos antiguos completan lo que no cubre. Un cuerpo sin ningún límite se informa como un cambio de formato, del que te avisa `ccam doctor`.

### El registro de límite

Cuando Claude Code detiene un turno en un límite de uso, el hook del [plugin de ccam](/docs/install#claude-code-plugin) lee de la conversación la solicitud rechazada: qué ventana estaba llena (5 horas, semanal, Opus o Sonnet) y cuándo se restablece. Lo escribe en `~/.local/state/ccam/usage/<name>.limit.json`. ccam aplica el registro por encima de todo lo demás que sabe, de cualquier fuente (`Fetch`): cada ventana que contiene cuenta como 100 % hasta su restablecimiento, por antiguas o bajas que sean las cifras guardadas, así que `ccam next --auto`, `ccam pick`, `--fallback` y [`auto_switch`](/features/next#auto-switch) no vuelven a elegir la cuenta antes de entonces.

Si no se sabe nada más, el registro es lo único que muestra la cuenta (`source: the usage limit Claude Code reported` con `-v`). Las sesiones desatendidas (ejecuciones `-p` y sesiones en segundo plano) no escriben ningún registro, ni tampoco una detención establecida por una organización (por ejemplo, cuando se queda sin uso disponible o el límite de un grupo es $0), que sus administradores pueden levantar en cualquier momento. `ccam login` en otra cuenta de Claude descarta el registro; eliminar el archivo también lo borra.

### Renovar tokens

Los tokens de actualización rotan: al renovar uno obtienes un nuevo token de actualización y el anterior queda invalidado. Si ccam renovara el inicio de sesión mientras Claude Code se está ejecutando con la misma cuenta, esa sesión quedaría cerrada. Por eso `ccam usage` lee los tokens y nunca los escribe, salvo que pases `--refresh`. Entonces ccam (`Fetcher.refresh`):

1. rechaza las cuentas registradas en su ubicación con `ccam import`, ya que `claude` a secas puede usar ese inicio de sesión sin que ccam lo sepa;
2. rechaza los inicios de sesión guardados en el Llavero de macOS;
3. toma un archivo de bloqueo, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, para que dos procesos de ccam nunca renueven ambos el mismo token; mientras está tomado, la cuenta muestra `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)`. `ccam rename`, `remove` y `logout` también lo toman, y se niegan mientras hay una renovación en curso;
4. con el bloqueo tomado, se niega si hay una sesión de la cuenta en ejecución, sin importar cómo se haya iniciado (mediante ccam, `claude` a secas, un IDE, o una sesión en segundo plano y el daemon que la respalda, que renueva el inicio de sesión por sí mismo), y omite la solicitud si otro proceso ya renovó el token;
5. no envía nada mientras Claude Code mantiene sus propios bloqueos de renovación en la carpeta de la cuenta (`.oauth_refresh.lock`, `<folder>.lock`): Claude Code está renovando ese inicio de sesión en ese momento;
6. envía el token de actualización al endpoint de tokens y espera hasta 30 segundos la respuesta, el tiempo que Claude Code da a sus propias renovaciones. Una vez enviada, una solicitud nunca se interrumpe, ni siquiera con Ctrl+C: puede que el servidor ya haya rotado el token de actualización, y solo su respuesta trae el nuevo inicio de sesión. Ctrl+C o un hangup terminan las demás consultas, y ccam sale (130 tras Ctrl+C) una vez guardada la respuesta;
7. guarda con el bloqueo de escritura de Claude Code (`.storage-write.lock`) y cambia solo los campos de token de `.credentials.json` tal como esté en ese momento: lo que Claude Code escribió mientras tanto (un inicio de sesión MCP) se conserva, un inicio de sesión que Claude Code renovó mientras tanto se mantiene, y un cierre de sesión nunca se deshace. Cuando la respuesta indica la duración del nuevo token de actualización, se guarda como `refreshTokenExpiresAt`, igual que hace Claude Code.

Una renovación que no se completa indica el motivo, en la línea de la cuenta o encima de sus cifras guardadas: `token not renewed: a session is running`, por ejemplo, o `login may be lost: …` cuando se envió una solicitud pero su respuesta no se pudo usar ni guardar, con el consejo de ejecutar `ccam login <name>` si Claude Code te pide luego que inicies sesión.

### El recolector de la línea de estado

Claude Code ejecuta el comando de línea de estado y le pasa un documento JSON por una tubería. ccam lee de él el nombre del modelo y `rate_limits`:

<!-- motion:usage-statusline -->

1. ccam identifica la cuenta a partir de `CCAM_PROFILE`, que ccam define al iniciar Claude Code, o si no, comparando `CLAUDE_CONFIG_DIR` con las carpetas de las cuentas (`runStatusline`).
2. ccam lee `rate_limits.five_hour` y `rate_limits.seven_day` (`ParseStatusLine`). El JSON no tiene límites de uso en las cuentas con clave de API, ni antes de la primera respuesta de la sesión; en esos casos no se guarda nada.
3. ccam escribe `~/.local/state/ccam/usage/<name>.statusline.json` (`SaveStatusLine`). `ccam usage` lo usa cuando el endpoint no responde.
4. ccam imprime `[name] · model · 5h N% · 7d N%`. Con `--wrap`, ejecuta tu comando mediante `sh -c` con la misma entrada e imprime su salida como Claude Code muestra una línea de estado: cada línea que no está en blanco, con `[name]` delante de la primera. Un comando que falla muestra solo el nombre de la cuenta. En Windows ejecuta el comando como lo hace Claude Code: con Git Bash si está instalado y, si no, con PowerShell.
5. Cuando el uso de 5 horas o semanal está al 90 % o más y hay otra cuenta, la línea termina con `→ ccam next --auto --continue`, por ejemplo `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

El comando envuelto no tiene un límite de tiempo propio: Claude Code da 10 minutos a un comando de línea de estado y detiene una ejecución cuando empieza una actualización más reciente, así que un comando lento se muestra igual que sin ccam. Cuando Claude Code termina ccam, el comando también termina (en Linux y macOS, con todos los procesos que inició).

`ccam statusline install` escribe `statusLine` en el `settings.json` de cada cuenta y conserva las demás claves del objeto, como `padding`, `refreshInterval` y `hideVimModeIndicator`. Cuando una cuenta ya tiene un comando de línea de estado, ccam lo conserva como `--wrap '<old command>'`; en Windows, como `--wrap-base64 <base64>`, que no necesita comillas en ninguna de las dos shells. Si la línea de estado ya ejecuta ccam (desde otra ruta, o entre comillas como lo escribía una versión anterior), ccam la reemplaza en lugar de envolverla, y conserva el comando que envuelve.

En Windows, la ruta a ccam se escribe con barras normales y `~` para la carpeta personal (`~/.local/bin/ccam.exe statusline`), que entienden tanto Git Bash como PowerShell. Una ruta que necesita comillas se escribe completa para la shell que usa Claude Code: `'…'` para Git Bash, `& '…'` para PowerShell. Un `settings.json` compartido a través de un grupo se escribe una sola vez, a través de su enlace, así que todas las cuentas del grupo lo reciben. Al volver a ejecutarlo, imprime `(already set)` para los archivos que ya usan ccam.

## Seguridad y casos límite

- `ccam usage` nunca cambia un archivo de tokens salvo que pases `--refresh`. Cuando un token de acceso venció, inicia la cuenta una vez con `ccam <name>` y Claude Code lo renueva. Un inicio de sesión que venció del todo necesita `ccam login <name>`.
- `--refresh` nunca toca las cuentas en ejecución (incluso cuando Claude Code se inició sin ccam), las cuentas importadas ni los inicios de sesión del Llavero. Una renovación rechazada muestra las cifras guardadas y el motivo (`token not renewed: …`).
- El endpoint de uso no está documentado públicamente y puede cambiar. Si cambia, sigues viendo las cifras guardadas, y `ccam doctor` informa el cambio de formato.
- Sin cifras guardadas y sin respuesta, la línea de la cuenta muestra `lookup failed:` y el motivo, por ejemplo `no saved usage` con `--source cache`.
- Los archivos de caché solo contienen porcentajes y horas de restablecimiento, nunca tokens.

## Relacionado

- [Límites de uso y cambio de cuenta](/docs/usage-limits)
- [Cambiar a la siguiente cuenta](/features/next)
- [Qué hacer cuando Claude Code llega al límite de uso](/guides/usage-limit-workflow)
- [Referencia de comandos](/docs/commands)
