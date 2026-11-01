---
title: "Lancer Claude Code"
description: "Lancez Claude Code avec tout compte via ccam <name|number>, ccam run ou le tableau de bord ; les arguments passent tels quels et --dry-run affiche la commande."
updated: 2026-11-01
order: 1
translation: machine
sourceHash: 11df08608a1e
---

`ccam work` lance Claude Code avec le compte `work`, et tout ce qui suit le nom est transmis sans modification à `claude`. Utilisez-le partout où vous taperiez `claude` : dans un terminal, dans un script ou avec une entrée redirigée.

## Utilisation

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

Indiquez le nom d'un compte ou son numéro dans `ccam list`. ccam ne lit pas les arguments qui suivent le nom, donc `ccam work --help` affiche l'aide de Claude Code. Un `--` après le nom est facultatif. L'entrée standard, le code de sortie et les signaux se comportent comme avec `claude`.

Quand le premier argument commence par `-` et n'est pas l'une des options propres à ccam, ccam choisit le compte pour vous et transmet tous les arguments à Claude Code. Il utilise le compte défini pour le dossier courant avec [`ccam folder`](/features/folder), sinon le compte par défaut : celui défini avec `ccam use`, ou le compte au plus petit numéro si vous n'en avez pas défini.

### Le tableau de bord

`ccam` sans argument ouvre le [tableau de bord](/features/dashboard) quand l'entrée et la sortie standard sont un terminal, et affiche `ccam list` dans le cas contraire. Placez-vous sur un compte (avec les flèches, ou en tapant son numéro) et appuyez sur `Enter` pour le lancer ; `n` lance le compte qui a le plus de marge et `o` propose `--continue`, `--resume` ou d'autres arguments. Le tableau de bord se ferme avant le démarrage de Claude Code, si bien que le terminal est d'abord restauré.

**Sièges Team et Enterprise.** La sélection automatique de compte (`ccam next`, `--fallback`, `ccam pick`, la touche `n` du tableau de bord et, une fois que vous l'activez, [`auto_switch`](/features/next#auto-switch)) passe aussi d'un siège à l'autre au sein d'une même organisation Claude Team ou Enterprise, ce qui peut enfreindre les conditions d'Anthropic. Lisez la [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) et renseignez-vous auprès de votre organisation et d'Anthropic avant de vous y fier.

Quand le shell définit `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` ou `ANTHROPIC_AUTH_TOKEN`, Claude Code peut utiliser cette clé ou ce jeton au lieu de la connexion propre au compte : ccam vous avertit donc sur la sortie d'erreur standard avant le lancement (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Il se contente de lire ces variables ; un [compte à clé d'API](/features/api) les efface au contraire.

Par défaut, ccam se remplace par Claude Code, si bien que rien d'autre ne continue de tourner. Avec `relay = true` (ou [`auto_switch = true`](/features/next#auto-switch) sous Linux et macOS) en haut de `config.toml`, au-dessus du premier `[[group]]`, `[[profile]]` ou `[[folder]]`, ccam lance plutôt les sessions interactives comme processus enfant et continue de tourner à côté d'elles. Cela permet à `/ccam:switch`, ou à `auto_switch` à une limite d'utilisation, de poursuivre la conversation sur un autre compte de lui-même, avec les options, le mode de permission, le niveau d'effort et l'agent de la session ; avant que `/ccam:switch` la déplace vers un compte d'une autre organisation, ccam vous le demande ([détails](/features/next#relay)).

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--dry-run` | Afficher l'environnement et la commande que ccam exécuterait, puis quitter | désactivé |
| `--fallback` | Pour les exécutions `-p`/`--print` : quand le compte atteint sa limite d'utilisation, relancer la même commande avec le compte connecté qui a le plus de marge ([détails](/features/next#run-again-at-the-limit)) | désactivé |
| `-v`, `--version` | Afficher `ccam <version>` (uniquement en premier argument) | |
| `-h`, `--help` | L'aide de ccam, sous la forme `ccam --help` ou `ccam run --help` | |

Règles des arguments :

- ccam ne lit ses propres options qu'avant le nom du compte. Après le nom, chaque argument appartient à Claude Code, y compris `--dry-run` et `--help`.
- `--dry-run` fonctionne sous la forme `ccam --dry-run work`, `ccam run --dry-run work` et `ccam --dry-run` (pour le compte du dossier ou le compte par défaut). `--dry-run` et `--fallback` peuvent venir dans n'importe quel ordre.
- `ccam run <name|number> [claude args]` est la forme longue de `ccam <name>`. `ccam run` sans nom utilise le compte du dossier ou le compte par défaut.
- Les noms de compte ne peuvent pas être des noms de commandes ccam comme `add`, `list`, `run` ou `usage`, si bien que `ccam <word>` n'est jamais ambigu.

La simulation affiche une ligne que vous pourriez coller dans un shell, ou dans PowerShell sous Windows. Elle commence par ce que le lancement change dans l'environnement : toujours `CLAUDE_CONFIG_DIR` et `CCAM_PROFILE`, `env -u <name>` pour une variable qu'il retire (comme un `CCAM_RELAY` hérité), et, pour un [compte à clé d'API](/features/api), son point de terminaison, son modèle et la provenance de la clé (`ANTHROPIC_API_KEY="$WORK_KEY"`), jamais la clé elle-même :

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

La ligne PowerShell remet ensuite les variables en place, comme le fait `VAR=x cmd` dans un shell, et transmet à Claude Code sa ligne de commande via `Start-Process`, si bien que Windows PowerShell 5.1 et PowerShell 7 l'exécutent de la même façon. Quand un lancement donne à une variable une valeur vide, une ligne `#` en dessous indique que PowerShell avant la version 7.5 supprime une telle variable à la place. Pour un `claude.cmd` installé par npm, la ligne exécute la commande `cmd.exe` qu'utiliserait le lancement (voir plus bas), et un argument que le lancement y refuserait fait aussi échouer la simulation.

## Fonctionnement

Chaque lancement suit les mêmes étapes, de la ligne de commande jusqu'à `exec` :

<!-- motion:launch -->

Chaque case s'allume pendant l'exécution de son étape, et les lignes en dessous sont l'environnement et les arguments avec lesquels Claude Code démarre.

1. ccam démarre dans `app.Main`, qu'appelle `cmd/ccam/main.go`. Si ccam a été lancé sous un autre nom, il s'agit d'un [raccourci](/features/shortcuts), qui va directement au lanceur.
2. ccam examine `os.Args` avant toute intervention de cobra (`app.parse`). Seuls `ccam` seul (le tableau de bord), un nom de commande ccam (`config.Reserved`), `-h`/`--help` et `-v`/`--version` vont à `cli.Execute`. Tout le reste est un lancement.
3. ccam lit `~/.config/ccam/config.toml`, le seul fichier qu'il analyse (`launcher.Launch`).
4. ccam trouve le compte par son nom ou son numéro, retire un éventuel `--` initial et résout l'exécutable `claude` (`launcher.Prepare`, `ResolveClaude`). Sans nom, il prend le compte du dossier, en comparant le répertoire de travail aux entrées `[[folder]]` de `config.toml` sous forme de chaînes, sans tenir compte de la casse sous macOS et Windows (`HereProfile`). La recherche dans le PATH est mise en cache dans `~/.local/state/ccam/claude-path` avec la date de modification du fichier et un hachage du `PATH`, si bien que ccam ne parcourt pas le `PATH` à chaque lancement. Un fichier modifié ou absent, ou un `PATH` modifié, déclenche une nouvelle recherche. Un chemin `claude_bin` dans `config.toml` (`~/` est développé) doit être un fichier exécutable ; sinon, le lancement s'arrête avec `claude_bin X is not an executable file. Check claude_bin in config.toml` avant que quoi que ce soit ne soit écrit. ccam refuse un `claude` qui serait ccam lui-même, comparé par inode, car cela tournerait en boucle.
5. ccam compare le dossier du compte `~/.local/share/ccam/profiles/<name>/` avec ses paramètres de partage, en n'appelant que `lstat` et `readlink`, un ou deux appels système par élément (`checkLinks`, `share.Check`). Un lien manquant est créé sur-le-champ. Tout ce qui pourrait faire perdre des données arrête le lancement.
6. ccam ajoute les serveurs MCP partagés (`applyMCP` ; voir plus bas).
7. ccam enregistre la session (`session.Register`). Il écrit un marqueur nommé d'après son PID, contenant l'heure de démarrage du processus, dans `~/.local/state/ccam/locks/profile-<name>/<pid>`, quoi qu'indique `XDG_RUNTIME_DIR`, si bien que chaque shell voit les mêmes marqueurs. Pour un compte qui partage l'historique des conversations, ccam cherche d'abord une session en cours dans `history-<group>-<hash>`, où le hachage est celui du dossier courant avec les liens symboliques résolus, et s'y enregistre aussi (`historyCheck`).

   `exec` conserve le PID, si bien que le marqueur désigne ensuite le processus Claude Code lui-même. ccam utilise un fichier marqueur plutôt qu'un verrou tenu, car les processus lancés par Claude Code hériteraient du verrou et le garderaient après la fin de la session. Un marqueur qui ne peut pas être écrit n'arrête pas le lancement : ccam affiche `ccam: could not mark the session as running (…); starting Claude Code anyway`.

   Les autres vérifications de ccam qui déterminent si un compte est en cours d'exécution (renommage, suppression, modifications de groupes, `doctor --fix`, la synchronisation MCP) lisent aussi les propres enregistrements de Claude Code dans le dossier du compte (`sessions/<pid>.json`, `daemon.lock`), si bien qu'elles voient aussi Claude Code lancé sans ccam : `claude` seul, un IDE, les sessions en arrière-plan et leur démon.
8. ccam définit le titre du terminal à `claude · <name>` (quand la sortie d'erreur standard est un terminal), écrit l'heure de lancement dans `~/.local/state/ccam/last-used/<name>` sans `fsync`, et appelle `syscall.Exec` (`Plan.Exec`).

Claude Code s'exécute avec `CLAUDE_CONFIG_DIR` défini sur le dossier du compte et `CCAM_PROFILE` défini sur le nom du compte. ccam efface toute valeur héritée de ces variables, ainsi qu'un `CCAM_RELAY` hérité. Pour un dossier `~/.claude` enregistré avec `ccam import` sous macOS, il définit aussi `CLAUDE_SECURESTORAGE_CONFIG_DIR` à une valeur vide, afin que le compte et `claude` seul utilisent la même connexion du trousseau ; sous Windows, il fait de même quand Claude Code conserve la connexion dans le Gestionnaire d'identification.

### Pourquoi exec

ccam remplace son propre processus par `claude` (`syscall.Exec`). Par conséquent :

- ccam ne reste pas en mémoire ;
- les frappes au clavier, les signaux et le code de sortie sont ceux de Claude Code ;
- les arguments restent un tableau, jamais une chaîne de shell.

Windows n'a pas d'`exec`. Là, ccam démarre `claude.exe` comme processus enfant avec le même environnement, les mêmes arguments et la même console, l'attend et quitte avec son code de sortie (`Plan.Exec`). Un `claude.cmd` installé par npm s'exécute dans `cmd.exe` avec son chemin complet et chaque argument entre guillemets, pour que `cmd.exe` n'en lise aucun comme sa propre syntaxe (`&`, `%PATH%`) ; un argument contenant un saut de ligne ne peut pas être transmis ainsi et est refusé, avec le conseil d'installer la version native de Claude Code.

Ctrl+C atteint Claude Code directement depuis la console : ccam l'ignore donc. Le marqueur de session porte le PID de ccam, qui vit exactement aussi longtemps que Claude Code, et il est supprimé quand Claude Code se termine. Pendant cette attente, ccam ne fait rien d'autre.

Le travail après une session, comme repérer les serveurs ajoutés avec `claude mcp add`, a lieu au lancement suivant : ccam compare la date de modification de `.claude.json` avec `~/.local/state/ccam/mcp/<name>/seen` et n'analyse le fichier que s'il a changé.

### Comment une ligne de commande est aiguillée

Trois lignes de commande, trois chemins à travers `internal/app` :

<!-- motion:launch-route -->

Un raccourci est reconnu sur `argv[0]`, `ccam <name>` sur le premier argument, et les deux aboutissent à `syscall.Exec`. Seule une commande ccam comme `ccam usage` atteint cobra.

### Serveurs MCP partagés

Avec la stratégie `flag` par défaut, ccam transmet les serveurs MCP du groupe à Claude Code sous forme d'argument, placé avant les vôtres :

<!-- motion:launch-mcp -->

ccam n'ajoute l'argument que si le `mcp.json` du groupe existe, et le place juste après `claude` (`mcp.FlagArgs`, `Plan.InsertArgs`). Vos propres `--mcp-config` ou `--strict-mcp-config` s'appliquent alors toujours par-dessus. C'est un seul argument, `--mcp-config=<path>`, car le `--mcp-config` de Claude Code accepte plusieurs valeurs : écrit en deux arguments, il lirait aussi une invite qui suit comme fichier de configuration.

ccam décide où l'option ne peut pas aller (`claudeVerb`). Il lit la sous-commande de Claude Code dans les arguments comme le fait Claude Code, en sautant les options avec leurs valeurs. `remote-control` ou `rc` ne reçoit pas l'option (voir plus bas), pas plus qu'un verbe de session en arrière-plan donné en premier (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`). `ccam work mcp list`, `get` et `remove` n'en reçoivent pas non plus, car Claude Code la refuse à cet endroit ; ccam indique plutôt sur la sortie d'erreur standard que les serveurs du groupe se chargent aussi dans les sessions du compte :

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Toutes les autres sous-commandes et toutes les sessions gardent l'option.

Avec `mcp_strategy = "sync"` dans `config.toml`, ccam écrit plutôt les serveurs dans le `.claude.json` du compte, et ne le fait pas si le compte est déjà en cours d'exécution. Un lancement saute aussi ce travail quand ni le `mcp.json` du groupe ni `.claude.json` n'ont changé depuis la dernière vérification. Si le compte possède son propre serveur portant le même nom qu'un serveur du groupe, la synchronisation garde celui du compte et le signale une fois. Elle remet aussi en place une copie qui a disparu de `.claude.json` (`put back group default's MCP servers …`), et ne modifie ou ne supprime que les copies qu'elle a écrites et que personne n'a modifiées depuis.

### Remote Control

`ccam work rc` (ou `ccam work remote-control`) démarre le Remote Control de Claude Code en tant que `work`. `claude remote-control` refuse `--mcp-config` avant son verbe et n'en accepte aucun après, et ses sessions ne lisent les serveurs MCP que dans le `.claude.json` du compte. La stratégie de l'option fonctionne donc différemment ici :

<!-- motion:launch-rc -->

1. ccam trouve `remote-control` ou `rc` comme verbe de Claude Code, même après des options et leurs valeurs (`ccam work --model opus rc`), et n'ajoute aucune option (`claudeVerb`, `mcp.FlagArgs`).
2. ccam compare l'heure de modification du `mcp.json` du groupe avec le `sync-stamp` du compte, et vérifie que chaque serveur du groupe se trouve toujours dans `.claude.json` (`mcp.Synced`). Si `.claude.json` contient déjà les serveurs actuels du groupe, ou si le groupe n'en a aucun, Remote Control démarre sans poser de question.
3. Sinon, ccam nomme les serveurs qui manqueraient aux sessions et demande :

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. Avec `y`, ccam les copie comme le fait la stratégie `sync` (`mcp.Sync`) : sous le verrou de Claude Code, après une sauvegarde, en conservant les serveurs propres du compte. Avec `N`, ou sans terminal, Remote Control démarre sans eux, et ccam indique `mcp_strategy = "sync"`, qui les copie avant chaque lancement.

ccam pose la question à chaque lancement de Remote Control tant que les serveurs ne sont pas copiés. Après une copie, il ne repose la question que lorsque le `mcp.json` du groupe change ou qu'un serveur copié disparaît de `.claude.json`. Une copie dans `.claude.json` l'emporte sur `--mcp-config` lors des lancements normaux suivants, donc `ccam mcp add`, `remove` et `import` mettent à jour les copies que détiennent les autres comptes ([Serveurs MCP partagés](/features/mcp)).

### Le budget de 10 ms

Le budget du chemin de lancement est de 10 ms au 95e centile ; en pratique, ccam ajoute environ 4 ms avant le démarrage de Claude Code. Le chemin de lancement ne construit jamais l'arborescence de commandes cobra, ne démarre pas Bubble Tea et n'ouvre aucune connexion réseau, et `TestPackageInitIsCheap` échoue si les fonctions `init()` des paquets prennent plus de 5 ms au total. `scripts/perf-check.sh` mesure `ccam --dry-run work`, `ccam work` et un raccourci avec `hyperfine` et échoue au-delà de 10 ms, et `BenchmarkLaunchPath` dans `internal/launcher` exécute les étapes en boucle. Pour voir le temps de chaque étape, définissez `CCAM_TRACE=1` :

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Sécurité et cas particuliers

- Un nom ou un numéro inconnu s'arrête avec `no account "<name>"` et renvoie vers `ccam list`. Sans aucun compte, ccam vous demande d'en ajouter un avec `ccam add <name>`.
- Si un chemin partagé est devenu un fichier ordinaire ou pointe au mauvais endroit, ccam s'arrête avec `stopped because of share link problems:`, liste les problèmes et vous demande de lancer `ccam doctor --fix`. Il ne déplace pas de données sur le chemin de lancement.
- Si une autre session Claude Code qui partage l'historique des conversations avec ce compte (même groupe) est en cours dans le même dossier, ccam avertit que l'historique peut se mélanger et demande `Run anyway? [y/N]`. Le dossier est comparé avec les liens symboliques résolus. Sans terminal (une entrée standard venant de `/dev/null` compte comme absence de terminal), ou depuis une tâche en arrière-plan, il affiche l'avertissement et continue. Une nouvelle exécution `-p` démarre sa propre conversation et ne reçoit pas d'avertissement ; `--version`, `--help` et les sous-commandes qui n'ouvrent aucune conversation, comme `mcp`, `auth` ou `plugin`, ne sont pas vérifiés du tout.
- `--dry-run` saute la synchronisation MCP, l'avis de nouveau serveur, la note de `mcp list`, la vérification de l'historique, le marqueur de session, le titre et l'heure de dernière utilisation, et n'exécute pas exec. Il vérifie tout de même les liens de partage, et crée un lien manquant comme le ferait un lancement normal.
- `ccam work rc` alors que `work` est déjà en cours d'exécution laisse `.claude.json` tel quel, le signale, et démarre Remote Control sans les serveurs du groupe ; relancez-le une fois la session terminée pour les copier. Avec `--dry-run`, il affiche le message mais ne pose pas de question et ne copie rien.
- Avec `mcp_strategy = "sync"`, une synchronisation échouée affiche `ccam: MCP sync failed (launching anyway): …` et Claude Code démarre avec les serveurs qu'il avait déjà.

## Voir aussi

- [Commandes raccourcis](/features/shortcuts)
- [Utiliser Remote Control avec ccam](/guides/remote-control)
- [Gérer les comptes](/features/accounts)
- [Référence des commandes](/docs/commands)
- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [CLAUDE_CONFIG_DIR expliqué](/guides/claude-config-dir)
