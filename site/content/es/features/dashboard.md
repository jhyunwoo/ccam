---
title: El panel
description: ccam en una terminal abre un panel a pantalla completa con cuentas, grupos, reglas de carpeta, extensiones y comprobaciones; cada cambio muestra antes su plan.
updated: 2026-10-30
order: 15
translation: machine
sourceHash: "d65f1b5bb0c8"
---

`ccam` sin argumentos abre un panel a pantalla completa cuando la entrada y la salida estándar son una terminal; si no, imprime `ccam list`. El panel muestra cada cuenta con su uso y sus grupos, y hace lo mismo que los comandos: iniciar, agregar, iniciar sesión, cambiar grupos, definir reglas de carpeta, listar sesiones en segundo plano, comprobar y reparar. Al iniciar una cuenta, o al conectarse a una sesión en segundo plano, se cierra; al agregar una cuenta o iniciar sesión, cede la terminal por un momento y luego vuelve al panel.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · r to refresh · a to add · ? for shortcuts
```

La línea bajo los nombres de las pestañas dice en una frase para qué sirve la pestaña actual. En una terminal más ancha, la pestaña Accounts también muestra el correo electrónico de cada cuenta. En todas las pestañas, el final de la línea con los nombres de las pestañas muestra si [`auto_switch`](/features/next#auto-switch) está activado (Linux y macOS; Windows ignora este ajuste). Cuando lo está, la pestaña Health lista qué le falta todavía: Claude Code 2.1.214 o posterior, y el plugin de ccam activado en cada cuenta que inicia sesión con una suscripción de Claude.

Cada vez que se abre, el panel pregunta a GitHub en segundo plano si existe una versión más reciente de ccam; el primer fotograma nunca la espera y, cuando no tienes conexión, no se muestra nada. Si la hay, una línea bajo el título la recomienda:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` pregunta primero y luego actualiza como lo hace [`ccam update`](/features/update): la versión se comprueba con sus sumas de comprobación SHA-256 y reemplaza ccam en un solo paso. El panel en el que estás sigue siendo la versión anterior, así que sal y vuelve a iniciar `ccam`. Cuando la actualización no se puede hacer aquí, por ejemplo con una copia de ccam instalada desde el paquete `.deb`, el panel muestra el mismo consejo que `ccam update` y mantiene la oferta.

El primer fotograma se dibuja solo a partir de `config.toml`, en mucho menos de 100 ms; el uso se completa a medida que responde cada cuenta. Cada límite tiene su propia barra: la ventana de 5 horas, la ventana semanal de todos los modelos y la ventana semanal propia de Fable; `-` significa que la cuenta no tiene esa ventana. Un porcentaje se vuelve amarillo al 80 % y rojo al 100 %. Si una consulta falla, los valores en caché se muestran atenuados. Una ventana que Claude Code informó llena muestra 100 % hasta su restablecimiento. Un token de acceso caducado muestra un `⚠ token expired` atenuado (Claude Code lo renueva en el siguiente inicio), mientras que `⚠ login expired` significa que la cuenta necesita `ccam login`.

## Pestañas

| Pestaña | Qué muestra | Teclas |
| --- | --- | --- |
| Accounts | número, nombre, correo electrónico, barras de 5 horas, semanal y de Fable, último uso, grupos; `✔` la predeterminada, `▸` la cuenta de esta carpeta, `•` en ejecución (también Claude Code iniciado sin ccam: `claude` a secas, un IDE, sesiones en segundo plano y su daemon) | más abajo |
| Agents | lo que lista [`ccam agents`](/features/agents): las sesiones en segundo plano y las sesiones abiertas de todas las cuentas, con su estado, carpeta y hora de inicio (más abajo) | `enter` conectarse (como `ccam attach`), `x` detener, `r` volver a cargar |
| Groups | cada grupo, lo que comparte y sus cuentas, y las cuentas que no están en ningún grupo | `enter` editar, `n` nuevo, `e` renombrar, `x` quitar |
| Folders | las reglas de [`ccam folder`](/features/folder) y la que se aplica aquí, incluso cuando llegaste a la carpeta mediante un enlace simbólico (las reglas nuevas se guardan con los enlaces resueltos) | `n` usar una cuenta aquí, `enter` cambiar, `x` quitar |
| Extensions | lo que muestra [`ccam ext`](/features/ext), cada extensión una vez con los grupos de los que proviene, `off:` las cuentas en las que un plugin no está activado, los servidores propios de la cuenta que se aplican solo a esta carpeta, `This folder's project` y los problemas que no pertenecen a ninguna cuenta | `r` volver a cargar |
| Health | lo que informa [`ccam doctor`](/features/doctor), con un recuadro por cuenta (más abajo) | `f` reparar el recuadro seleccionado, `F` reparar todo (`doctor --fix`), `↑`/`↓` pasar de un recuadro a otro |

`←`/`→` o `Tab` cambian de pestaña, `?` lista las teclas de la pestaña actual, `r` actualiza el uso, sin pasar por la caché corta, y vuelve a cargar las pestañas Agents, Extensions y Health, y `q` o `Esc` salen.

## Teclas en la pestaña Accounts

| Tecla | Qué hace |
| --- | --- |
| `Enter` | Inicia Claude Code con la cuenta seleccionada |
| `1`–`9` | Va a la cuenta con ese número; escribe `1` `2` rápido para la 12. Solo `Enter` inicia |
| `↑`/`↓`, `k`/`j` | Mueve la selección |
| `n` | Inicia la cuenta que elegiría [`ccam next --auto`](/features/next) |
| `o` | Inicia con opciones: `--continue`, `--resume`, otros argumentos o una ejecución de prueba |
| `u` | Detalles de uso: las ventanas de 5 horas, semanal y por modelo, con sus restablecimientos y de dónde vienen las cifras |
| `a` | Agrega una cuenta (el asistente de más abajo) |
| `l` | Vuelve a iniciar sesión |
| `s` | Elige los grupos de la cuenta |
| `d` | La convierte en la cuenta predeterminada |
| `f` | La usa en la carpeta actual |
| `e` | Le cambia el nombre |
| `c` | Define o quita su atajo |
| `x` | La quita (pregunta antes, con No seleccionado) |

**Asientos de Team y Enterprise.** La selección automática de cuenta (`ccam next`, `--fallback`, `ccam pick`, la tecla `n` del panel y, una vez que lo activas, [`auto_switch`](/features/next#auto-switch)) también pasa de un asiento a otro dentro de una misma organización de Claude Team o Enterprise, lo que puede infringir los términos de Anthropic. Lee las [preguntas frecuentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) y consulta con tu organización y con Anthropic antes de depender de ella.

## La pestaña Health

La pestaña Health ejecuta las comprobaciones de [`ccam doctor`](/features/doctor) al abrirla y dibuja la misma vista: una primera línea con una marca por cuenta y luego un recuadro por cuenta. Cada fila de un recuadro es un área comprobada, primero `Folder`, `Login` y `Share links`, marcada con `✔` correcto, `ℹ` información, `⚠` una advertencia o `✘` un error, así que ves de un vistazo qué cuenta tiene un problema y dónde. Los resultados que no pertenecen a ninguna cuenta van en un recuadro `General`.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` pasan de un recuadro a otro, y el recuadro seleccionado tiene el borde del color de acento. `f` repara solo el recuadro seleccionado, como lo haría `ccam doctor --fix`: en el recuadro de una cuenta repara solo esa cuenta, y en el recuadro `General`, los resultados que no pertenecen a ninguna cuenta (rutas de plugins y atajos); las demás cuentas solo se vuelven a comprobar. Pregunta primero (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), y un recuadro sin nada que reparar solo dice `nothing to repair in <name>`. `F` repara todo a la vez, como `ccam doctor --fix`, y `r` vuelve a comprobar.

## La pestaña Agents

La pestaña Agents muestra lo que imprime [`ccam agents`](/features/agents): las sesiones en segundo plano y las sesiones abiertas de todas las cuentas, con la cuenta a la que pertenece cada una. Consulta al Claude Code de cada cuenta a la vez (`claude agents --json`, cada uno con un tiempo de espera de 15 segundos) cuando abres la pestaña por primera vez, al pulsar `r`, después de detener una sesión y después de un cambio en la configuración, y muestra `Asking every account's Claude Code for its sessions…` hasta que todas hayan respondido.

```console
  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Background agents and open sessions of every account; enter attaches to one.

    Account   ID        Name   State          Folder      Started
 ❯  work      -         -      open (busy)    ~/code/web  now
    personal  49ec0028  probe  blocked, idle  ~/code/app  12m

 ✘ lab: could not list its sessions: exit status 1 'claude agents --json' is
   not available in this environment.

 enter to attach · x to stop · r to reload · ? for shortcuts
```

`State` es el estado de una sesión en segundo plano y lo que está haciendo su proceso, o `open` para una sesión interactiva, como en `ccam agents`. Si no se puede consultar una cuenta, una línea `✘` debajo de la lista muestra lo que dijo Claude Code, como `lab` arriba, que tiene la vista de agentes desactivada. Las sesiones en segundo plano terminadas se omiten (`ccam agents --all` las lista), y cuando todas las cuentas respondieron y ninguna tiene sesiones, la pestaña dice ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` En una terminal estrecha, la pestaña quita `Started`, luego `Name` y luego `Folder`, y recorta por la izquierda las rutas de carpeta largas.

`Enter` sobre una sesión en segundo plano cierra el panel y abre la sesión en esta terminal con su cuenta, como hace `ccam attach <id>`; el panel no vuelve (en Claude Code, `←` va a su vista de agentes y `Ctrl+Z` al shell). Una sesión abierta no tiene ID, así que `Enter` y `x` solo indican que hace falta una sesión en segundo plano. `x` pregunta primero, con `No, keep it` seleccionado: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` Si respondes que sí, ccam ejecuta `claude stop <id>` con esa cuenta, dice ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` o por qué no pudo, y vuelve a leer la lista.

## Cambiar grupos

En la pestaña Groups, `Enter` abre el editor: los elementos a la izquierda, las cuentas a la derecha, `Tab` para cambiar de lado y `Space` para marcar o desmarcar. Un elemento que una de las cuentas ya recibe de otro grupo aparece atenuado junto con el motivo, igual que una cuenta que recibiría un elemento dos veces; un elemento solo puede provenir de un grupo. `Enter` muestra el plan:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Si un archivo difiere de la copia del grupo, ccam pregunta primero qué hacer (conservar el del grupo, usar el mío o combinar). Después el cambio se ejecuta en segundo plano con el mismo código que [`ccam group`](/features/share): primero se guarda la configuración, los datos de cada cuenta se mueven en orden y los originales van a `.ccam-backup/`; una cuenta que deja de compartir servidores MCP conserva los servidores del grupo en su `.claude.json`.

Un cambio que movería los datos de una cuenta en ejecución se rechaza con el nombre de la cuenta, incluso cuando la cuenta se inició mientras el plan estaba abierto, y un plan que quedó desactualizado porque otro proceso de ccam cambió los grupos se rechaza con `the groups changed while this plan waited; nothing was moved, make the change again`. Mientras otro proceso de ccam cambia los grupos, un nuevo cambio, cambio de nombre o eliminación dice `another ccam process is changing the groups right now`; después de un cambio que se quedó a medias, esperan hasta que `F` en la pestaña Health (o `ccam doctor --fix`) lo haya terminado.

## Agregar una cuenta

`a` abre el asistente que `ccam add` también usa en una terminal:

1. El nombre de la cuenta.
2. Cómo inicia sesión: con una suscripción de Claude ahora (en el navegador), con una [clave de API](/features/api) o con una suscripción de Claude más tarde.
3. Para una clave de API: de dónde lee ccam la clave en cada inicio (`env:NAME` o `command:COMMAND`), el endpoint, cómo se envía la clave y el modelo.
4. Los grupos a los que unirse, con `default` marcado (si todavía no existe ningún grupo, lo que compartirá el grupo `default`).
5. Un atajo, si quieres uno.
6. Un resumen para confirmar.

`Esc` vuelve un paso atrás y, en el primer paso, cierra el asistente.

## Quitar una cuenta

`x` pregunta antes de que ocurra nada, con No seleccionado, y dice lo que hace la eliminación: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` Para una cuenta con clave de API no hay inicio de sesión que revocar, y cuando otra cuenta tiene el mismo inicio de sesión no se cierra su sesión. Una cuenta importada conserva su carpeta: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## Colores

Los colores siguen los roles de la interfaz de terminal de Claude Code, así que ccam encaja a su lado; los nombres de grupo reciben un color según su nombre. ccam le pregunta a la terminal por su fondo y elige un tema claro u oscuro. Puedes fijar uno tú mismo al principio de `config.toml`, antes del primer `[[group]]`, `[[profile]]` o `[[folder]]`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Los temas `-daltonized` son para el daltonismo rojo-verde y los temas `-ansi` usan los 16 colores propios de la terminal. `NO_COLOR` desactiva los colores.

## Seguridad y casos límite

- El panel llama a las mismas funciones que los comandos, así que se cumplen todas las reglas: planes antes de mover datos, rechazos para las cuentas en ejecución, copias de seguridad y tokens de inicio de sesión que nunca salen de una cuenta.
- El panel no tiene ejecución de prueba: `ccam --dry-run` sin un nombre imprime cómo se iniciaría la cuenta de esta carpeta, y `o` sobre una cuenta ofrece una ejecución de prueba de ese único inicio.
- La terminal se restaura antes de que se inicie Claude Code o se abra el inicio de sesión.
- En una pantalla demasiado pequeña para un diálogo, primero se ocultan las filas de la pestaña; las líneas se recortan con `…` en lugar de partirse. En una pantalla baja se quitan después las líneas en blanco y la descripción de la pestaña, para que la línea de teclas siga visible. Por debajo de 40 columnas u 8 filas el panel solo muestra `Too small: ccam needs 40x8 (now 36x7)`.
- `Ctrl+C` mientras se aplica un cambio, una eliminación, una actualización o una reparación espera a que termine: la línea de estado dice `ccam quits when it is done.`
- Después de que `U` haya actualizado ccam, la oferta no vuelve cuando el panel regresa de agregar una cuenta o de un inicio de sesión.

## Relacionado

- [Compartir con grupos](/features/share)
- [Iniciar Claude Code](/features/launch)
- [Administrar cuentas](/features/accounts)
- [Límites de uso](/features/usage)
