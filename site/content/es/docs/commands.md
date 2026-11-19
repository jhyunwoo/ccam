---
title: Referencia de comandos
metaTitle: Referencia de comandos de ccam
description: Todos los comandos y opciones de ccam, desde ccam add y ccam <name> hasta ccam usage, ccam next y ccam doctor.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: "28e1ad1dcbad"
---

Pon `--dry-run` antes de un comando para ver lo que haría sin hacerlo, por ejemplo `ccam --dry-run work`. Funciona al iniciar Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) y con `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` y `doctor --fix`. Los comandos que solo leen, como `list`, `usage`, `share` y `ext`, la aceptan y de todos modos no cambian nada. Los demás la rechazan en lugar de ejecutarse, incluido `group rename`.

<!-- commands -->

## Iniciar Claude Code

`ccam <name|number> [arguments]` inicia Claude Code con esa cuenta. ccam no lee los argumentos que van después del nombre; los pasa a `claude` en el mismo orden, y la entrada estándar, el código de salida y las señales se comportan exactamente igual que con `claude`. Un `--` después del nombre es opcional.

`ccam` con argumentos pero sin nombre de cuenta, como `ccam -p "summarise"`, usa la cuenta asignada a la carpeta actual o, si no hay ninguna, la cuenta predeterminada. Establécelas con `ccam folder <name> [folder]` y `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` antes del nombre de la cuenta es para ejecuciones `-p`/`--print`. Cuando la cuenta llega a su límite de uso, ccam vuelve a ejecutar el mismo comando con la cuenta con sesión iniciada que tenga más margen e imprime solo la salida de la ejecución que terminó. Consulta [Cuentas por carpeta](/features/folder) y [Volver a ejecutar al llegar al límite](/features/next#run-again-at-the-limit).

**Asientos de Team y Enterprise.** La selección automática de cuenta (`ccam next`, `--fallback`, `ccam pick`, la tecla `n` del panel y, una vez que lo activas, [`auto_switch`](/features/next#auto-switch)) también pasa de un asiento a otro dentro de una misma organización de Claude Team o Enterprise, lo que puede infringir los términos de Anthropic. Lee las [preguntas frecuentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) y consulta con tu organización y con Anthropic antes de depender de ella.

Un atajo como `ccwork` inicia su cuenta de la misma forma: `ccwork --resume` hace lo mismo que `ccam work --resume`. Consulta [Atajos](/docs/commands#shortcuts) más abajo.

## Agregar y quitar cuentas

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

`ccam import` sin una carpeta busca `~/.claude` y `~/.claude-*` y pregunta por cada una; una carpeta que falla se informa como `Skipped <path>: <reason>` y las demás se registran. Si el nombre de una carpeta ya está en uso, importa esa carpeta por separado con `ccam import <folder> --name <name>`. Las carpetas importadas se registran donde están; ccam nunca mueve ni copia un inicio de sesión.

`ccam login --expired` también inicia sesión en las cuentas cuyo inicio de sesión borró Claude Code. `ccam remove` cierra la sesión de una cuenta gestionada (una cuya carpeta creó ccam, a diferencia de una registrada con `ccam import`) antes de que su carpeta vaya a la papelera, salvo que otra cuenta tenga el mismo inicio de sesión. Quitar una cuenta importada la saca de sus grupos y deja en su carpeta una copia de cada elemento compartido. En macOS, `ccam rename` conserva el nombre de la carpeta, ya que el Llavero guarda el inicio de sesión con la ruta de la carpeta.

`ccam add --api` agrega una cuenta que usa una clave de API en lugar de un inicio de sesión de Claude; ccam lee la clave en cada inicio y nunca la guarda. Consulta [Cuentas con clave de API](/features/api).

`ccam rename`, `ccam remove` y los comandos `ccam group` se niegan a cambiar una cuenta mientras Claude Code se está ejecutando con esa cuenta, sin importar cómo se haya iniciado (incluidas las sesiones en segundo plano). Cierra primero la sesión. También se niegan mientras un cambio de grupo esté sin terminar: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` muestra cada cuenta con su número de posición, atajo, correo electrónico y carpeta. Su arreglo `groups` lista los grupos de la cuenta, y su objeto `share` asigna a cada elemento (`mcp`, `settings`, `memory`, `extensions`, `history`) el grupo del que proviene, o `local`.

## Atajos

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Un atajo es un comando al que tú le pones nombre y que inicia Claude Code con una cuenta concreta. Igual que `ccam work`, pasa todos los argumentos a Claude Code sin cambios, así que `ccwork --resume` y `ccwork -p "hi"` funcionan como esperas. Es un enlace a ccam en `~/.local/bin`.

Los atajos son opcionales. El asistente de `ccam add` pregunta si quieres crear uno, `ccam add <name> --shortcut <command>` crea uno sin preguntar y `ccam list` muestra el atajo de cada cuenta, o `-` si no tiene.

Los nombres usan los caracteres a-z, 0-9, `.`, `-` y `_`. Un nombre no puede ser un comando de ccam, el atajo de otra cuenta ni el nombre de otro programa que esté en tu `PATH`, ni un archivo de `~/.local/bin` que ccam no haya creado, como tu propia copia o enlace de ccam. Un atajo editado a mano en `config.toml` que incumple estas reglas, que es `claude` o `ccam`, o que se da a dos cuentas hace que `config.toml` no se pueda cargar, y se nombra la cuenta.

`ccam rename` conserva el nombre del atajo, y `ccam remove` elimina el atajo junto con la cuenta. `ccam shims install` vuelve a crear los enlaces de todas las cuentas, por ejemplo después de mover ccam; un archivo que estorba a un atajo no detiene los demás.

## Compartir

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Un grupo comparte elementos entre sus cuentas; una cuenta puede estar en varios grupos, pero cada elemento proviene de un solo grupo. Consulta [compartir con grupos](/features/share).

## Copiar una configuración existente

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` copia `settings.json`, `CLAUDE.md`, las extensiones y los servidores MCP de usuario de una carpeta de Claude Code a un grupo, para que todas las cuentas que comparten esos elementos los tengan. La carpeta solo se lee.

- `--settings`, `--memory`, `--extensions`, `--mcp` y `--history` copian solo los elementos que indicas. Sin ellas, en una terminal una lista de selección te permite elegir, con todo marcado excepto el historial de conversaciones; en otros casos se copia todo excepto el historial de conversaciones.
- `--mcp` también funciona cuando `~/.claude` está vacío o no existe pero `~/.claude.json` tiene servidores.
- `-g` o `--group` elige el grupo.
- `-y` o `--yes` copia sin preguntar y conserva la copia del grupo cuando un archivo es distinto.

Los inicios de sesión MCP no se copian, ya que un inicio de sesión copiado cierra la sesión de los demás que lo tienen; el plan nombra las cuentas que inician sesión con `/mcp`. Consulta [copiar una configuración existente](/docs/sharing#copy-an-existing-setup).

## Servidores MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` mueve a su grupo los servidores que solo tiene esa cuenta, incluidos los que tiene solo para una carpeta (`claude mcp add` a secas), que entonces se cargan en todas las carpetas. Si el grupo ya usa el nombre de un servidor para una definición distinta, ese servidor se omite; agrega `--force` para reemplazar la copia del grupo, de la que ccam hace antes una copia de seguridad.

## Uso y cambio de cuenta

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

- `next --auto`, `--fallback` y `ccam pick` eligen de la misma forma: primero por debajo de `auto_max_usage` (`config.toml`, 100 por defecto), según `auto_strategy` (`headroom`, `expiring` o `balance`), contando la ventana semanal de `--model`. Nunca eligen una cuenta al 100 % (una ventana que Claude Code informó llena cuenta como 100 % hasta su restablecimiento) ni un inicio de sesión con clave de API, y se quedan dentro de la organización de la cuenta desde la que parten (`--cross-org` o `cross_org = true` para buscar más allá).
- `--wait` toma su valor después de `=` (`--wait=2h`); escrito con un espacio se rechaza con una sugerencia. Cuando no se puede elegir ninguna, salen con el código 75 o, con `--wait`, esperan a que se restablezca la primera cuenta.
- `ccam next --continue` lleva una conversación a otra cuenta incluso entre cuentas que no comparten el historial: Claude Code la copia entonces en la nueva cuenta. Con `--brief`, la nueva cuenta empieza en cambio una conversación nueva, a partir de un resumen de la anterior que ccam escribe sin un modelo ([`--brief`](/features/next#brief)).
- A diferencia de `ccam <name>`, `ccam next` lee sus propias opciones, así que las opciones de Claude Code van después de `--`; sin él, ccam se detiene con `unknown flag`.
- `ccam agent` ejecuta una tarea con `claude -p` en esta carpeta con una cuenta distinta de aquella desde la que se llama, en el modo de permisos y con el nivel de esfuerzo de la sesión que lo llama, e imprime la última respuesta de Claude. Al llegar al límite de esa cuenta, la conversación sigue con otra ([Subagentes en otras cuentas](/features/agent)).
- Dentro de Claude Code, `/ccam:switch [account]` del plugin pasa una conversación a otra cuenta, y `ccam hook` es lo que ejecutan sus hooks.
- Con `auto_switch = true` en `config.toml` (Linux y macOS; desactivado por defecto), el hook de límite también termina una sesión interactiva que inició ccam y la continúa en otra cuenta con margen, pero solo en una del mismo tipo: otro asiento de la misma organización u otro plan personal ([`auto_switch`](/features/next#auto-switch)).

Consulta [Cambiar a la siguiente cuenta](/features/next), [Elegir una cuenta para scripts](/features/pick), [Subagentes en otras cuentas](/features/agent) y [Sesiones en segundo plano de todas las cuentas](/features/agents).

## Mantenimiento

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

`ccam update` comprueba la suma de verificación SHA-256 de la versión antes de reemplazar ccam, y `--version v0.2.0` instala una versión concreta. Si instalaste ccam desde el paquete `.deb`, actualízalo con apt. Consulta [actualizar](/docs/install#update).

`ccam uninstall` lista lo que quita y pregunta primero: todas las cuentas, como hace `ccam remove`, los atajos, las carpetas de ccam y el comando `ccam`. Una carpeta importada se queda donde está. Consulta [desinstalar](/docs/install#uninstall).

## Archivos

| Ruta | Contenido |
| --- | --- |
| `~/.config/ccam/config.toml` | Cuentas, grupos y ajustes (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). Los ajustes van al principio, antes del primer `[[group]]`, `[[profile]]` o `[[folder]]`: uno que está debajo de una tabla pertenece a esa tabla, así que ccam avisa de ello en cada inicio y en `ccam doctor`, y los comandos que guardan la configuración se niegan hasta que se mueva. Los comentarios no se conservan cuando ccam guarda el archivo |
| `~/.local/share/ccam/profiles/<name>/` | El `CLAUDE_CONFIG_DIR` de cada cuenta |
| `~/.local/share/ccam/groups/<group>/` | Configuración, servidores MCP, extensiones e historial compartidos |
| `~/.local/state/ccam/` | Caché de uso, registros de límite, bloqueos, notas y otro estado (consulta la lista debajo de la tabla) |
| `~/.local/bin/<shortcut>` | Atajos como `ccwork`, enlaces a ccam |

La carpeta de estado contiene:

- la caché de uso y los registros de límite (`usage/<name>.json`, `usage/<name>.limit.json`);
- las horas del último uso;
- los marcadores de sesión y los bloqueos (`locks/`);
- los atajos que creó ccam (`shortcuts.json`);
- un cambio de grupo sin terminar (`group-change.json`);
- las notas que dejan `/ccam:switch` y el hook de límite (`handoff/<account>/<session>.json`);
- cuándo `ccam pick` o `auto_switch` eligió por última vez cada cuenta (`picks/<name>`; una elección cuenta durante 2 minutos);
- los resúmenes de `--brief`, que se guardan 7 días (`brief/<session>.md`).

Definir `XDG_CONFIG_HOME`, `XDG_DATA_HOME` o `XDG_STATE_HOME` cambia la ubicación de estas carpetas. En Windows, `~` es `%USERPROFILE%` y los atajos son `<shortcut>.exe`.
