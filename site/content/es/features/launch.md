---
title: Iniciar Claude Code
description: Inicia Claude Code como cualquier cuenta con ccam <name|number>, ccam run o el panel; los argumentos pasan intactos a claude y --dry-run muestra el comando.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: "11df08608a1e"
---

`ccam work` inicia Claude Code con la cuenta `work`, y todo lo que va después del nombre llega a `claude` sin cambios. Úsalo donde escribirías `claude`: en una terminal, en un script o con la entrada por una tubería.

## Uso

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Indica el nombre de una cuenta o su número de `ccam list`. ccam no lee los argumentos que van después del nombre, así que `ccam work --help` muestra la ayuda de Claude Code. Un `--` después del nombre es opcional. La entrada estándar, el código de salida y las señales se comportan igual que con `claude`.

Cuando el primer argumento empieza con `-` y no es una de las opciones propias de ccam, ccam elige la cuenta por ti y pasa todos los argumentos a Claude Code. Usa la cuenta asignada a la carpeta actual con [`ccam folder`](/features/folder) o, si no hay ninguna, la cuenta predeterminada: la establecida con `ccam use`, o la cuenta con el número más bajo si no has establecido ninguna.

### El panel

`ccam` sin argumentos abre el [panel](/features/dashboard) cuando la entrada y la salida estándar son una terminal, y en otro caso imprime `ccam list`. Ve a una cuenta (con las flechas o escribiendo su número) y presiona `Enter` para iniciarla; `n` inicia la cuenta con más margen y `o` ofrece `--continue`, `--resume` u otros argumentos. El panel se cierra antes de que se inicie Claude Code, así que primero se restaura la terminal.

**Asientos de Team y Enterprise.** La selección automática de cuenta (`ccam next`, `--fallback`, `ccam pick`, la tecla `n` del panel y, una vez que lo activas, [`auto_switch`](/features/next#auto-switch)) también pasa de un asiento a otro dentro de una misma organización de Claude Team o Enterprise, lo que puede infringir los términos de Anthropic. Lee las [preguntas frecuentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) y consulta con tu organización y con Anthropic antes de depender de ella.

Cuando la shell tiene definida `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` o `ANTHROPIC_AUTH_TOKEN`, Claude Code puede usar esa clave o ese token en lugar del inicio de sesión propio de la cuenta, así que ccam te avisa en la salida de error estándar antes de iniciar (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Solo lee las variables; una [cuenta con clave de API](/features/api) las borra en su lugar.

De forma predeterminada, ccam se reemplaza a sí mismo por Claude Code, así que no queda nada más en ejecución. Con `relay = true` (o [`auto_switch = true`](/features/next#auto-switch) en Linux y macOS) al principio de `config.toml`, antes del primer `[[group]]`, `[[profile]]` o `[[folder]]`, ccam inicia en cambio las sesiones interactivas como proceso hijo y sigue ejecutándose junto a ellas. Así, `/ccam:switch`, o `auto_switch` en un límite de uso, puede continuar la conversación en otra cuenta por sí solo, con las opciones, el modo de permisos, el nivel de esfuerzo y el agente de la sesión; antes de que `/ccam:switch` la mueva a una cuenta de otra organización, ccam te pregunta ([detalles](/features/next#relay)).

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--dry-run` | Imprime el entorno y el comando que ccam ejecutaría, y sale | desactivada |
| `--fallback` | Para ejecuciones `-p`/`--print`: cuando la cuenta llega a su límite de uso, vuelve a ejecutar el mismo comando con la cuenta con sesión iniciada que tenga más margen ([detalles](/features/next#run-again-at-the-limit)) | desactivada |
| `-v`, `--version` | Imprime `ccam <version>` (solo como primer argumento) | |
| `-h`, `--help` | La ayuda de ccam, como `ccam --help` o `ccam run --help` | |

Reglas de los argumentos:

- ccam lee sus propias opciones solo antes del nombre de la cuenta. Después del nombre, todos los argumentos son de Claude Code, incluidos `--dry-run` y `--help`.
- `--dry-run` funciona como `ccam --dry-run work`, `ccam run --dry-run work` y `ccam --dry-run` (para la cuenta de la carpeta o la predeterminada). `--dry-run` y `--fallback` pueden ir en cualquier orden.
- `ccam run <name|number> [claude args]` es la forma larga de `ccam <name>`. `ccam run` sin un nombre usa la cuenta de la carpeta o la predeterminada.
- Los nombres de cuenta no pueden ser nombres de comandos de ccam como `add`, `list`, `run` o `usage`, así que `ccam <word>` nunca es ambiguo.

La ejecución de prueba imprime una línea que podrías pegar en una shell, o en PowerShell en Windows. Empieza con lo que el inicio cambia en el entorno: siempre `CLAUDE_CONFIG_DIR` y `CCAM_PROFILE`, `env -u <name>` para una variable que descarta (como un `CCAM_RELAY` heredado) y, para una [cuenta con clave de API](/features/api), su endpoint, su modelo y de dónde sale la clave (`ANTHROPIC_API_KEY="$WORK_KEY"`), nunca la clave en sí:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

La línea de PowerShell restablece las variables después, como hace `VAR=x cmd` en una shell, y le pasa a Claude Code su línea de comandos mediante `Start-Process`, así que Windows PowerShell 5.1 y PowerShell 7 la ejecutan igual. Cuando un inicio establece una variable con un valor vacío, una línea `#` debajo indica que PowerShell anterior a 7.5 elimina esa variable en su lugar. Para un `claude.cmd` instalado por npm, la línea ejecuta el comando de `cmd.exe` que usaría el inicio (consulta más abajo), y un argumento que el inicio rechazaría allí también hace fallar la ejecución de prueba.

## Cómo funciona

Cada inicio sigue los mismos pasos, desde la línea de comandos hasta `exec`:

<!-- motion:launch -->

Cada recuadro se ilumina mientras se ejecuta su paso, y las líneas de abajo son el entorno y los argumentos con los que se inicia Claude Code.

1. ccam arranca en `app.Main`, a la que llama `cmd/ccam/main.go`. Si ccam se inició con otro nombre, es un [atajo](/features/shortcuts) y va directamente al iniciador.
2. ccam examina `os.Args` antes de que intervenga cobra (`app.parse`). Solo `ccam` a secas (el panel), un nombre de comando de ccam (`config.Reserved`), `-h`/`--help` y `-v`/`--version` van a `cli.Execute`. Todo lo demás es un inicio.
3. ccam lee `~/.config/ccam/config.toml`, el único archivo que analiza (`launcher.Launch`).
4. ccam busca la cuenta por nombre o número, descarta un `--` inicial y resuelve el ejecutable `claude` (`launcher.Prepare`, `ResolveClaude`). Sin nombre, toma la cuenta de la carpeta comparando el directorio de trabajo con las entradas `[[folder]]` de `config.toml` como cadenas, sin distinguir mayúsculas y minúsculas en macOS y Windows (`HereProfile`). La búsqueda en PATH se guarda en caché en `~/.local/state/ccam/claude-path` junto con la fecha de modificación del archivo y un hash de `PATH`, así que ccam no busca en `PATH` en cada inicio. Un archivo modificado o que falta, o un `PATH` distinto, provoca una búsqueda nueva. Una ruta `claude_bin` en `config.toml` (`~/` se expande) debe ser un archivo ejecutable; si no, el inicio se detiene con `claude_bin X is not an executable file. Check claude_bin in config.toml` antes de escribir nada. ccam rechaza un `claude` que sea el propio ccam, comparando el inodo, porque eso crearía un bucle.
5. ccam compara la carpeta de la cuenta `~/.local/share/ccam/profiles/<name>/` con su configuración de uso compartido, llamando solo a `lstat` y `readlink`, una o dos llamadas al sistema por elemento (`checkLinks`, `share.Check`). Un enlace que falta se crea en el momento. Cualquier cosa que podría perder datos detiene el inicio.
6. ccam agrega los servidores MCP compartidos (`applyMCP`; consulta más abajo).
7. ccam registra la sesión (`session.Register`). Escribe un marcador con el nombre de su PID, que contiene la hora de inicio del proceso, en `~/.local/state/ccam/locks/profile-<name>/<pid>`, diga lo que diga `XDG_RUNTIME_DIR`, para que todas las shells vean los mismos marcadores. En una cuenta que comparte el historial de conversaciones, ccam primero busca una sesión en ejecución en `history-<group>-<hash>`, donde el hash es el de la carpeta actual con los enlaces simbólicos resueltos, y también se registra ahí (`historyCheck`).

   `exec` conserva el PID, así que el marcador pasa a nombrar al propio proceso de Claude Code. ccam usa un archivo marcador y no un bloqueo retenido, porque los procesos que inicia Claude Code heredarían el bloqueo y lo conservarían después de que termine la sesión. Un marcador que no se puede escribir no detiene el inicio: ccam dice `ccam: could not mark the session as running (…); starting Claude Code anyway`.

   Las demás comprobaciones de ccam sobre si una cuenta se está ejecutando (cambio de nombre, eliminación, cambios de grupo, `doctor --fix`, la sincronización de MCP) también leen los registros propios de Claude Code en la carpeta de la cuenta (`sessions/<pid>.json`, `daemon.lock`), así que también ven Claude Code iniciado sin ccam: `claude` a secas, un IDE, sesiones en segundo plano y su daemon.
8. ccam establece el título de la terminal en `claude · <name>` (cuando la salida de error estándar es una terminal), escribe la hora de inicio en `~/.local/state/ccam/last-used/<name>` sin `fsync` y llama a `syscall.Exec` (`Plan.Exec`).

Claude Code se ejecuta con `CLAUDE_CONFIG_DIR` apuntando a la carpeta de la cuenta y `CCAM_PROFILE` con el nombre de la cuenta. ccam borra cualquier valor heredado de estas variables, y un `CCAM_RELAY` heredado. Para una carpeta `~/.claude` registrada con `ccam import` en macOS, también establece `CLAUDE_SECURESTORAGE_CONFIG_DIR` con un valor vacío, para que la cuenta y `claude` a secas usen el mismo inicio de sesión del Llavero; en Windows hace lo mismo cuando Claude Code guarda el inicio de sesión en el Administrador de credenciales.

### Por qué exec

ccam reemplaza su propio proceso por `claude` (`syscall.Exec`). Como resultado:

- ccam no se queda en memoria;
- las pulsaciones de teclas, las señales y el código de salida son los de Claude Code;
- los argumentos siguen siendo un arreglo, nunca una cadena de shell.

Windows no tiene `exec`. Allí, ccam inicia `claude.exe` como proceso hijo con el mismo entorno, los mismos argumentos y la misma consola, lo espera y sale con su código de salida (`Plan.Exec`). Un `claude.cmd` instalado por npm se ejecuta en `cmd.exe` con su ruta completa y cada argumento entre comillas, para que `cmd.exe` no lea ninguno como su propia sintaxis (`&`, `%PATH%`); un argumento con un salto de línea no se puede pasar así y se rechaza, con el consejo de instalar el Claude Code nativo.

Ctrl+C llega a Claude Code directamente desde la consola, así que ccam lo ignora. El marcador de sesión lleva el PID de ccam, que vive exactamente lo mismo que Claude Code, y se elimina cuando Claude Code termina. Mientras espera, ccam no hace nada más.

El trabajo posterior a una sesión, como detectar servidores agregados con `claude mcp add`, ocurre en el siguiente inicio: ccam compara la fecha de modificación de `.claude.json` con `~/.local/state/ccam/mcp/<name>/seen` y analiza el archivo solo cuando cambió.

### Cómo se enruta una línea de comandos

Tres líneas de comandos, tres rutas por `internal/app`:

<!-- motion:launch-route -->

Un atajo se reconoce por `argv[0]` y `ccam <name>` por el primer argumento, y ambos terminan en `syscall.Exec`. Solo un comando de ccam como `ccam usage` llega a cobra.

### Servidores MCP compartidos

Con la estrategia predeterminada `flag`, ccam pasa los servidores MCP del grupo a Claude Code como un argumento, colocado antes de los tuyos:

<!-- motion:launch-mcp -->

ccam agrega el argumento solo si existe el `mcp.json` del grupo, y lo coloca justo después de `claude` (`mcp.FlagArgs`, `Plan.InsertArgs`). Tus propias opciones `--mcp-config` o `--strict-mcp-config` se siguen aplicando encima. Es un solo argumento, `--mcp-config=<path>`, porque `--mcp-config` de Claude Code acepta varios valores: escrito como dos argumentos, también leería una instrucción posterior como un archivo de configuración.

ccam decide dónde no puede ir la opción (`claudeVerb`). Lee el subcomando de Claude Code a partir de los argumentos igual que Claude Code, omitiendo las opciones junto con sus valores. `remote-control` o `rc` no reciben la opción (más abajo), ni tampoco un verbo de sesión en segundo plano indicado primero (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`). `ccam work mcp list`, `get` y `remove` tampoco la reciben, ya que Claude Code la rechaza allí; en su lugar, ccam indica en la salida de error estándar que los servidores del grupo también se cargan en las sesiones de la cuenta:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Todos los demás subcomandos y todas las sesiones conservan la opción.

Con `mcp_strategy = "sync"` en `config.toml`, ccam escribe los servidores en el `.claude.json` de la cuenta, y lo omite si la cuenta ya está en ejecución. Un inicio también omite ese trabajo cuando ni el `mcp.json` del grupo ni `.claude.json` cambiaron desde la última comprobación. Si la cuenta tiene un servidor propio con el mismo nombre que un servidor del grupo, la sincronización conserva el servidor de la cuenta y lo indica una vez. La sincronización también vuelve a poner una copia que ya no está en `.claude.json` (`put back group default's MCP servers …`), y solo cambia o quita las copias que escribió y que nadie cambió desde entonces.

### Remote Control

`ccam work rc` (o `ccam work remote-control`) inicia el Remote Control de Claude Code como `work`. `claude remote-control` rechaza `--mcp-config` antes de su verbo y no acepta ninguno después, y sus sesiones leen los servidores MCP solo del `.claude.json` de la cuenta. Por eso la estrategia de opción funciona de otra manera aquí:

<!-- motion:launch-rc -->

1. ccam encuentra `remote-control` o `rc` como verbo de Claude Code, incluso después de opciones y sus valores (`ccam work --model opus rc`), y no añade la opción (`claudeVerb`, `mcp.FlagArgs`).
2. ccam compara la hora de modificación del `mcp.json` del grupo con el `sync-stamp` de la cuenta, y comprueba que cada servidor del grupo siga en `.claude.json` (`mcp.Synced`). Si `.claude.json` ya contiene los servidores actuales del grupo, o el grupo no tiene ninguno, Remote Control se inicia sin preguntar.
3. Si no, ccam nombra los servidores que les faltarían a las sesiones y pregunta:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. Con `y`, ccam los copia como lo hace la estrategia `sync` (`mcp.Sync`): con el bloqueo de Claude Code, después de una copia de seguridad y conservando los servidores propios de la cuenta. Con `N`, o sin una terminal, Remote Control se inicia sin ellos, y ccam sugiere `mcp_strategy = "sync"`, que los copia antes de cada inicio.

ccam pregunta en cada inicio de Remote Control hasta que los servidores se copian. Después de una copia, vuelve a preguntar solo cuando el `mcp.json` del grupo cambia o un servidor copiado desaparece de `.claude.json`. Una copia en `.claude.json` gana sobre `--mcp-config` en los inicios normales posteriores, así que `ccam mcp add`, `remove` e `import` actualizan las copias que tienen otras cuentas ([Servidores MCP compartidos](/features/mcp)).

### El presupuesto de 10 ms

El presupuesto de la ruta de inicio es de 10 ms en el percentil 95; en la práctica, ccam agrega unos 4 ms antes de que se inicie Claude Code. La ruta de inicio nunca construye el árbol de comandos de cobra, no inicia Bubble Tea ni abre una conexión de red, y `TestPackageInitIsCheap` falla si las funciones `init()` de los paquetes tardan más de 5 ms en total. `scripts/perf-check.sh` mide `ccam --dry-run work`, `ccam work` y un atajo con `hyperfine` y falla por encima de 10 ms, y `BenchmarkLaunchPath` en `internal/launcher` ejecuta los pasos en un bucle. Para ver el tiempo de cada paso, define `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Seguridad y casos límite

- Un nombre o número desconocido se detiene con `no account "<name>"` y remite a `ccam list`. Si no hay ninguna cuenta, ccam te pide que agregues una con `ccam add <name>`.
- Si una ruta compartida se convirtió en un archivo normal o apunta al lugar equivocado, ccam se detiene con `stopped because of share link problems:`, las lista y te pide que ejecutes `ccam doctor --fix`. No mueve datos en la ruta de inicio.
- Si otra sesión de Claude Code que comparte el historial de conversaciones con esta cuenta (el mismo grupo) se está ejecutando en la misma carpeta, ccam advierte que el historial se puede mezclar y pregunta `Run anyway? [y/N]`. La carpeta se compara con los enlaces simbólicos resueltos. Sin una terminal (una entrada estándar desde `/dev/null` cuenta como ninguna), o desde un trabajo en segundo plano, muestra la advertencia y continúa. Una nueva ejecución `-p` inicia una conversación propia y no recibe advertencia; `--version`, `--help` y los subcomandos que no abren ninguna conversación, como `mcp`, `auth` o `plugin`, no se comprueban en absoluto.
- `--dry-run` omite la sincronización de MCP, el aviso de servidores nuevos, la nota de `mcp list`, la comprobación del historial, el marcador de sesión, el título y la hora del último uso, y no hace exec. Sí comprueba los enlaces compartidos y crea un enlace que falta, como lo haría un inicio normal.
- `ccam work rc` cuando `work` ya está en ejecución no toca `.claude.json`, lo indica e inicia Remote Control sin los servidores del grupo; vuelve a ejecutarlo cuando termine la sesión para copiarlos. Con `--dry-run` imprime el mensaje, pero no pregunta ni copia.
- Con `mcp_strategy = "sync"`, una sincronización fallida imprime `ccam: MCP sync failed (launching anyway): …` y Claude Code se inicia con los servidores que ya tenía.

## Relacionado

- [Atajos](/features/shortcuts)
- [Usar Remote Control con ccam](/guides/remote-control)
- [Administrar cuentas](/features/accounts)
- [Referencia de comandos](/docs/commands)
- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [CLAUDE_CONFIG_DIR explicado](/guides/claude-config-dir)
