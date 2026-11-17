---
title: Preguntas frecuentes
metaTitle: "Preguntas frecuentes: varias cuentas de Claude Code"
description: Respuestas sobre ccam, cuentas de Claude Code, límites de uso, uso compartido, privacidad y sistemas compatibles.
updated: 2026-11-17
order: 6
translation: machine
sourceHash: "9b13b37391f2"
---

## ¿ccam es una herramienta oficial de Anthropic?

No. ccam es un proyecto no oficial de código abierto y no está afiliado con Anthropic ni cuenta con su respaldo. Usa la opción `CLAUDE_CONFIG_DIR`, que [Claude Code documenta](https://code.claude.com/docs/en/env-vars) para ejecutar varias cuentas en paralelo.

## ¿Puedo usar varias cuentas de Claude Code?

Depende de los términos de cada cuenta y, en el caso de las cuentas de trabajo, de la política de tu organización: los [Términos para consumidores](https://www.anthropic.com/legal/consumer-terms) de Anthropic para Pro y Max, sus [Términos comerciales](https://www.anthropic.com/legal/commercial-terms) para Team y Enterprise, y su [Política de uso](https://www.anthropic.com/legal/aup) para todos. Revísalos antes de usar varias cuentas. ccam no aumenta, no combina ni oculta los límites de una cuenta, y cada cuenta sigue deteniéndose en los suyos. Mueve tu trabajo a otra de tus cuentas cuando lo pides, o por sí mismo en un límite cuando activas [`auto_switch`](/features/next#auto-switch) (Linux y macOS). Comprueba si eso está permitido para tus cuentas antes de usarlo; la decisión es tuya.

**Asientos de Team y Enterprise.** La selección automática de cuenta (`ccam next`, `--fallback`, `ccam pick`, la tecla `n` del panel y [`auto_switch`](/features/next#auto-switch) cuando lo activas) también pasa de un asiento a otro dentro de una misma organización de Claude Team o Enterprise. El centro de ayuda de Anthropic indica: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada asiento de un plan Team se asigna a una persona específica, y solo esa persona puede usarlo) y "Users can't route their work through another user's account to use that user's remaining usage" (los usuarios no pueden hacer pasar su trabajo por la cuenta de otro usuario para consumir el uso restante de ese usuario) ([planes Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); lo mismo para los [planes Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Por lo tanto, usar la selección automática entre asientos puede infringir los términos de Anthropic. Si la usas ahí es decisión tuya: consulta con tu organización y con Anthropic antes de depender de ella. Los planes personales (Pro, Max) se tratan en el párrafo anterior.

## ¿ccam evita los límites de uso?

No más allá del límite propio de ninguna cuenta: cada cuenta mantiene sus límites y se detiene en ellos, y ccam no los aumenta, no los combina ni los oculta. ccam los muestra y te permite continuar en otra de tus cuentas cuando una se agota; con `auto_switch` activado, lo hace por sí mismo (consulta la siguiente pregunta). Entre asientos de una misma organización de Team o Enterprise, continuar en otro asiento consume el uso restante de ese asiento, algo que el centro de ayuda de Anthropic descarta para el asiento de otro usuario (consulta Asientos de Team y Enterprise más arriba).

## ¿Puede ccam cambiar de cuenta por sí mismo cuando una sesión llega a su límite?

Solo si lo activas. Con `auto_switch = true` al principio de `~/.config/ccam/config.toml` (Linux y macOS; desactivado por defecto y experimental), cuando Claude Code detiene una sesión interactiva (una que inició ccam) en un límite de 5 horas, semanal, de Opus, de Sonnet o de Fable, ccam termina la sesión y continúa la conversación como otra de tus cuentas con margen: otro asiento de la misma organización de Team o Enterprise, u otro plan personal (Pro, Max). Nunca cambia antes del límite, nunca vigila tu uso en segundo plano y nunca saca el trabajo de una organización de Team o Enterprise ni lo mete en una.

En los asientos mueve tu trabajo de un asiento a otro, lo que puede infringir los términos de Anthropic (consulta Asientos de Team y Enterprise más arriba). Para los planes personales, los términos de Anthropic no lo abordan con claridad: los Términos para consumidores prohíben acceder a los Servicios "through automated or non-human means" (por medios automatizados o no humanos) salvo que Anthropic lo permita expresamente, y un cambio envía un mensaje en tu nombre. Activarlo es decisión tuya.

[`auto_switch`](/features/next#auto-switch) enumera todas las condiciones y lo que cuesta un cambio.

## ¿Cómo uso dos cuentas de Claude Code en una computadora?

Dale a cada cuenta su propia carpeta de configuración con `CLAUDE_CONFIG_DIR`. ccam lo hace por ti: ejecuta `ccam add` para cada cuenta y luego inícialas con `ccam personal` y `ccam work`, o con atajos a los que tú les pones nombre, como `ccwork`. La [guía](/guides/multiple-claude-code-accounts) también muestra cómo hacerlo a mano.

## ¿Cómo continúo una conversación en otra cuenta?

Ejecuta `ccam next --continue` en la carpeta: lleva la última conversación de la carpeta a la siguiente cuenta (`--auto` para la que tiene más margen). Dentro de Claude Code, `/ccam:switch` del [plugin de ccam](/docs/install#claude-code-plugin) hace lo mismo con la sesión en la que estás. Cuando las dos cuentas comparten el historial de conversaciones (`ccam group set default --add history`, o un grupo propio: `ccam group add pair --items history --members work,side`), la nueva cuenta abre esa misma conversación; si no, Claude Code la copia en el historial de la nueva cuenta. Añade `--brief` para que la nueva cuenta empiece con un resumen breve en lugar de la conversación completa.

## ¿ccam ve o envía mi inicio de sesión?

ccam nunca comparte ni copia tu inicio de sesión de Claude entre cuentas, y nunca copia los inicios de sesión de los servidores MCP. Lee el token de acceso de una cuenta solo para consultar al endpoint de uso de Anthropic el uso de esa cuenta: para `ccam usage`, el panel, `ccam doctor` y la elección automática (`ccam next --auto`, `ccam pick`, `--fallback`, y `auto_switch` cuando una sesión llega a su límite). Solo renueva uno caducado con `ccam usage --refresh`; con `usage_source = "claude"` no envía ningún token a ningún sitio y lee el inicio de sesión guardado (el archivo, o la entrada del Llavero de macOS) solo para saber si sigue funcionando. `ccam ext` y `ccam doctor` leen los nombres de los inicios de sesión MCP de una cuenta, nunca sus tokens. ccam no tiene telemetría. Aparte de esas consultas, la única solicitud que hace por su cuenta es la comprobación que hace el panel en GitHub de si hay una versión más reciente de ccam, que no lleva nada sobre ti ni sobre tus cuentas.

## ¿Cómo inicio sesión en un servidor remoto?

Ejecuta `ccam add` o `ccam login <name>` en la sesión SSH como de costumbre. Claude Code espera el inicio de sesión en el servidor, así que el navegador de tu computadora termina en una página `http://localhost:<port>/callback?code=…` que no carga. Copia esa dirección del navegador y pégala en la terminal: ccam se la pasa a Claude Code en el servidor, que completa el inicio de sesión. La dirección contiene un código de un solo uso, así que no la compartas.

## ¿Dónde se guardan mis cuentas?

En `~/.local/share/ccam/profiles/<name>/`, una carpeta por cuenta. En macOS, Claude Code guarda el inicio de sesión en sí en el Llavero, con un nombre derivado de esa carpeta, así que `ccam rename` conserva allí el nombre anterior de la carpeta.

## ¿Qué pasa con mi ~/.claude actual?

Nada, a menos que ejecutes `ccam import` o `ccam migrate`. La importación registra la carpeta donde está, así que `claude` a secas y `ccam <name>` siguen usando la misma cuenta y la misma configuración. `ccam migrate` copia su configuración, `CLAUDE.md`, sus extensiones y sus servidores MCP al grupo compartido `default`, y solo lee la carpeta.

## ¿Qué sistemas admite ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine y otros), macOS y Windows 10 y 11, de forma nativa o dentro de WSL. En Windows, compartir archivos como `settings.json` requiere el modo de desarrollador; las carpetas compartidas funcionan sin él, y sin él el primer `ccam add` crea el grupo `default` solo con servidores MCP y extensiones ([compartir en Windows](/features/share#on-windows)).

## ¿Por qué un servidor MCP aparece solo en una cuenta?

Los servidores agregados dentro de una sesión con `claude mcp add` se guardan solo para esa cuenta y, sin `-s user`, solo para la carpeta actual. ccam te avisa la próxima vez que inicies la cuenta, y `ccam mcp list` y `ccam ext` los muestran; mueve el servidor al grupo compartido con `ccam mcp import <account> <server>`.

## ¿ccam hace más lento a Claude Code?

Iniciar Claude Code a través de ccam agrega unos 4 milisegundos, medidos con `scripts/perf-check.sh` en una laptop. ccam lee un archivo pequeño, comprueba los enlaces compartidos y luego se reemplaza a sí mismo por `claude`, así que por defecto no se está ejecutando mientras trabajas. Windows no puede reemplazar un proceso, así que allí ccam inicia Claude Code y lo espera, sin hacer nada más. Con `relay` o `auto_switch` activados, ccam sigue ejecutándose de la misma forma junto a cada sesión interactiva, para poder continuar la conversación en otra cuenta.

## ¿Cómo actualizo ccam?

Ejecuta `ccam update`. Descarga la versión más reciente, comprueba su suma de verificación y reemplaza el archivo `ccam`; `ccam update --check` solo te indica si hay una versión más reciente. El panel también te avisa cuando hay una versión más reciente, y `U` la instala después de preguntar; ccam nunca se actualiza por su cuenta. Si instalaste el paquete `.deb`, instala el nuevo `.deb` con apt. Consulta [actualizar](/docs/install#update).

## ¿Cómo desinstalo ccam?

Ejecuta `ccam uninstall`. Lista lo que se va y pregunta primero; luego quita todas las cuentas (cerrando antes su sesión), los atajos, `~/.config/ccam`, `~/.local/share/ccam`, `~/.local/state/ccam` y el comando `ccam`. Un `~/.claude` importado se queda, con una copia de lo que compartía en sus grupos. Consulta [instalar](/docs/install#uninstall).
