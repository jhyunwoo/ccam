---
title: "Gérer les comptes"
description: "Ajoutez, connectez, listez, renommez, retirez et importez des comptes Claude Code : ccam add, ccam login, ccam list, ccam rename, ccam remove et ccam import."
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

Chaque compte ccam est un dossier de configuration Claude Code distinct, avec sa propre connexion. Utilisez ces commandes pour ajouter des comptes, les connecter, les lister, les renommer et les supprimer. Elles enregistrent aussi les dossiers que vous avez déjà.

## Utilisation

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account that was never logged in or whose login expired
ccam logout side                           # log out: the login and its MCP sign-ins
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, logs out, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Partout où une commande accepte `<name|number>`, vous pouvez indiquer le nom du compte ou son numéro d'emplacement (la colonne `#` de `ccam list`).

`ccam use` ne fait que changer le compte par défaut (`default` dans `config.toml`).

`ccam add` sans option dans un terminal ouvre un assistant. Il passe par ces étapes :

1. le nom du compte, sauf si vous en avez donné un ;
2. la façon dont le compte se connecte : une connexion Claude maintenant, une clé d'API ou une connexion Claude plus tard ;
3. les groupes à rejoindre, avec `default` coché (s'il n'existe encore aucun groupe, il affiche ce que le groupe `default` partagera) ;
4. `Create a shortcut command for work? [y/N]` ;
5. un récapitulatif à confirmer.

`Esc` revient d'une étape en arrière. L'une des options `--group`, `--no-group`, `--shortcut` ou `--no-login`, ou une exécution sans terminal (par exemple depuis un tube), fait sauter l'assistant. Un nom est alors obligatoire, et aucun raccourci n'est créé sauf si `--shortcut` en nomme un.

Sans `--group`, un nouveau compte rejoint le groupe `default`. Le premier compte le crée, avec le partage des serveurs MCP, des paramètres, de `CLAUDE.md` et des extensions ; voir [Partager avec des groupes](/features/share). Sous Windows, les fichiers ne peuvent être liés qu'en Mode développeur. Sans lui, le premier compte crée `default` avec seulement `mcp` et `extensions` et le signale (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), et ccam refuse `--group` avec un groupe qui partage les paramètres, la mémoire ou l'historique.

`ccam list` affiche l'emplacement, le nom, le raccourci, l'adresse e-mail, les groupes et la dernière utilisation, avec `*` devant le compte par défaut. L'adresse e-mail est celle avec laquelle le dossier est connecté maintenant (d'après son `.claude.json`) ; quand elle diffère de celle que ccam a enregistrée, celle-ci suit sous la forme `(was …)` :

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

```console
$ ccam list
   #  Name      Shortcut  Email                                     Groups   Last used
*  1  work      ccwork    other@company.com (was work@example.com)  default  2m ago
```

Avec `--json`, chaque compte est un objet avec `slot`, `shortcut`, `name`, `email` (celle de la connexion actuelle, comme dans le tableau), `default`, `dir`, `last_used` (omis quand le compte n'a jamais été lancé), `groups` et `share`. `share` associe chaque élément (`mcp`, `settings`, `memory`, `extensions`, `history`) au groupe dont il provient, ou à `"local"`.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `ccam add --no-login` | Créer le compte sans se connecter ; connectez-vous plus tard avec `ccam login <name>` | désactivé |
| `ccam add --group <group>` | Rejoindre ce groupe au lieu de `default`. Répétable ; `default` est créé s'il n'existe pas | `default` |
| `ccam add --no-group` | Ne rejoindre aucun groupe : le compte ne partage rien | désactivé |
| `ccam add --shortcut <command>` | Créer aussi cette commande raccourci dans `~/.local/bin` | aucun |
| `ccam add --api --key-from <source>` | Un compte qui utilise une clé d'API au lieu d'une connexion ; voir [Comptes à clé d'API](/features/api) | désactivé |
| `ccam list --json` | Afficher les comptes en JSON | désactivé |
| `ccam remove -y`, `--yes` | Supprimer sans demander | désactivé |
| `ccam logout -y`, `--yes` | Se déconnecter sans demander | désactivé |
| `ccam import --name <name>` | Nom du compte pour le dossier indiqué ; refusé sans dossier (`--name needs a folder: ccam import <folder> --name ops`) | déduit du nom du dossier |
| `ccam import -y`, `--yes` | Importer tous les dossiers trouvés sans demander | désactivé |
| `--dry-run` | Pris en charge par `ccam list` et `ccam completion`, qui ne modifient rien. Refusé par `add`, `login`, `logout`, `use`, `rename`, `remove` et `import`, avec `` --dry-run is not supported by `ccam add`; nothing was changed `` (en nommant la commande) | désactivé |

`ccam login` (en dehors de `--expired`), `ccam use` et `ccam rename` n'ont pas d'options propres. `ccam completion` accepte `bash`, `zsh` ou `powershell` ; toute autre valeur échoue avec `only bash, zsh and powershell are supported`.

Les noms de compte comptent de 1 à 32 caractères parmi `a-z`, `0-9`, `-` et `_`, et commencent par une lettre ou un chiffre. Un nom ne peut pas être composé uniquement de chiffres, car `ccam 2` désigne l'emplacement 2. Un nom ne peut pas non plus être une commande ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), si bien que `ccam <name>` n'est jamais ambigu. La liste est `config.Reserved` dans `internal/config`.

## Fonctionnement

Ajouter un compte crée son dossier, rejoint ses groupes, enregistre la configuration, connecte le compte et crée le raccourci, dans cet ordre. Quand vous exécutez `ccam add work --shortcut ccwork`, ccam suit ces étapes (`addProfile` dans `internal/cli/profiles.go`) :

<!-- motion:accounts-add -->

1. ccam vérifie le nom (`config.ValidateName`) et le raccourci (`shim.Check`) avant de créer quoi que ce soit.
2. ccam crée `~/.local/share/ccam/profiles/work/` avec le mode `0700` et attribue au compte le plus petit numéro d'emplacement libre, de sorte que le numéro d'un compte supprimé est réutilisé (`profile.Add`, `Config.NextSlot`). Le premier compte devient le compte par défaut.
3. Le compte rejoint ses groupes (`default`, sauf si `--group` ou `--no-group` en décide autrement ; le premier compte crée `default`), et ccam y relie le nouveau dossier vide, par exemple `settings.json` vers `groups/default/settings.json` (`applyShare`, avec `share.Build` puis `Apply`). En cas d'échec, le nouveau dossier est supprimé, car un dossier restant bloquerait chaque nouvelle tentative avec `folder already exists`.
4. Le compte est enregistré dans `~/.config/ccam/config.toml` via `fsutil.AtomicWrite` : un fichier temporaire, `fsync`, puis `rename`, si bien qu'un plantage ne laisse jamais une configuration à moitié écrite.
5. ccam lance `claude auth login` avec `CLAUDE_CONFIG_DIR` défini sur le dossier du compte, et `--email` avec la dernière adresse e-mail du compte pour que la page de connexion soit préremplie (`runLogin`). Si la connexion échoue, le compte est conservé, et ccam affiche l'erreur et vous indique de lancer `ccam login work` plus tard.
6. ccam lit `oauthAccount` dans `profiles/work/.claude.json` (`profile.ReadIdentity`) : l'adresse e-mail, les identifiants du compte et de l'organisation, le nom de l'organisation et l'offre (`claude_max`, `claude_team` …). Ce fichier peut peser plusieurs mégaoctets, donc ccam recherche cette seule clé au lieu d'analyser tout le fichier (`TopLevelValue`). ccam conserve ces informations dans `[profile.login]` dans `config.toml`.
7. ccam relie `~/.local/bin/ccwork` au binaire ccam, qui recherche `argv[0]` dans `config.toml` au démarrage, et enregistre le lien dans `~/.local/state/ccam/shortcuts.json` (`shim.Install`). Si le lien ne peut pas être créé, le compte est quand même ajouté, ccam indique pourquoi (`Could not create shortcuts: …`) et le récapitulatif affiche `shortcut: ccwork (not created; see above)`.

Une connexion peut prendre plusieurs minutes, et un autre processus ccam peut enregistrer la configuration entre-temps. Les commandes qui attendent rechargent donc `config.toml`, appliquent la seule modification (ici, l'adresse e-mail) et enregistrent (`env.update`). `ccam remove` recharge après sa question pour la même raison.

`ccam login` exécute les étapes 5 et 6 pour un compte existant. Claude Code ne prend l'adresse e-mail que comme indication pour la page de connexion : vous pouvez donc vous retrouver connecté à un autre compte Claude qu'auparavant ; ccam compare la nouvelle connexion avec `[profile.login]` et le signale (`ccam: this login is other@company.com, but work was dev@company.com before.`). Il conserve la nouvelle connexion et vous indique comment l'annuler : ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

Une connexion à un autre compte Claude supprime aussi le [relevé de limite](/features/usage#the-limit-record) du compte, puisque ces limites étaient celles de l'autre compte. Quand un autre compte ccam est connecté au même compte Claude, ccam avertit que les deux partagent leurs limites d'utilisation.

`ccam login --expired` connecte, l'un après l'autre, chaque compte qui n'a jamais été connecté ou dont la connexion a expiré : son jeton de renouvellement a expiré, ou Claude Code a effacé la connexion après un refus du serveur. Avant l'ouverture du premier navigateur, il vérifie que `config.toml` peut être enregistré.

`ccam logout` exécute `claude auth logout` dans le compte et efface l'adresse e-mail enregistrée ; l'enregistrement de connexion reste, si bien que la connexion suivante est comparée avec lui. Claude Code révoque la connexion et supprime ce qu'il conserve pour le compte, y compris les connexions aux serveurs MCP. ccam vous avertit donc et pose d'abord la question (`--yes` la saute) quand :

- la connexion contient des connexions MCP ;
- le compte est en cours d'exécution (ses sessions sont déconnectées aussi) ;
- le dossier d'un autre compte contient la même connexion (un dossier copié), que la déconnexion révoque également.

Par exemple :

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

Pour se connecter avec un autre compte Claude, `ccam login <name>` seul suffit, et il conserve les connexions MCP.

`ccam list` ajoute une colonne `Organization` quand les comptes appartiennent à plusieurs organisations (les offres personnelles affichent `personal`), et `ccam list --json` contient `plan` et `org` pour chaque compte.

### Se connecter sur une machine distante

Via SSH, la partie navigateur de `claude auth login` peut échouer. Claude Code attend la connexion sur un serveur de rappel sur la machine distante (`http://localhost:<port>/callback`), mais le navigateur qui termine la connexion s'exécute sur votre ordinateur, où cette adresse ne mène nulle part, et le navigateur affiche une erreur comme `ERR_CONNECTION_REFUSED`. L'adresse contient malgré tout le code à usage unique, que ccam peut donc transmettre :

<!-- motion:accounts-relay -->

ccam ne confie pas votre terminal à Claude Code (`runLogin`). Il relie l'entrée standard de Claude Code à un tube et y copie chaque ligne que vous tapez (`relayLogin`), si bien qu'un code collé à l'invite `Paste code here` de Claude Code arrive toujours. Quand une ligne est une adresse de rappel, ccam ne l'accepte que si l'hôte est `localhost`, `127.0.0.1` ou `[::1]`, qu'un port est indiqué, que le chemin est `/callback` et que `code` et `state` sont tous deux présents ; `http://` peut être omis (`callbackURL`). ccam demande alors cette adresse sur le serveur avec un délai maximal de 10 secondes et sans suivre les redirections, comme l'aurait fait le navigateur (`deliverCallback`). La requête atteint le serveur de rappel du `claude auth login` en attente (la ligne en pointillés du diagramme). Claude Code reçoit le code et enregistre la connexion, sous Linux dans `.credentials.json` dans le dossier du compte. Une ligne collée ne peut jamais amener ccam à contacter un autre hôte.

ccam affiche ceci avant le début de la connexion :

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

En cas de réussite, il affiche `→ Passed the sign-in to Claude Code on this machine.`

Pendant que la connexion attend une saisie dans un terminal, ccam désactive l'écho et le mode ligne du terminal et affiche un `*` pour chaque caractère que vous tapez, si bien qu'un code ou une adresse de rappel collés n'apparaissent pas à l'écran. ccam gère lui-même l'édition de la ligne (`readLine`) : `Enter` envoie la ligne, `Backspace` efface un caractère, et `Ctrl+D` sur une ligne vide termine la saisie. Le terminal est restauré à la fin de la connexion, ainsi qu'à `Ctrl+C`, `SIGTERM` ou `SIGHUP`, qui terminent toujours ccam. Sans terminal sur l'entrée standard et la sortie d'erreur standard, les lignes sont copiées telles qu'elles arrivent.

### Renommer, supprimer et importer

`ccam rename` et `ccam remove` agissent sur le dossier du compte et sur ce que ccam conserve ailleurs sous le nom du compte. ccam liste ces chemins dans `nameState` :

- `~/.local/state/ccam/last-used/<name>` ;
- `~/.local/state/ccam/mcp/<name>` ;
- les caches d'utilisation `~/.local/state/ccam/usage/<name>.json` et `<name>.statusline.json` ;
- le relevé de limite `<name>.limit.json` ;
- les notes de session dans `~/.local/state/ccam/handoff/<name>/`.

`ccam remove work` les traite dans cet ordre :

<!-- motion:accounts-remove -->

1. `ccam remove` refuse d'abord si le compte est en cours d'exécution, puis pose la question, et vérifie de nouveau après votre réponse (cela détecte une session lancée pendant que la question était affichée).
2. Il vérifie que `config.toml` peut être enregistré avant de modifier quoi que ce soit.
3. Si le compte partage les extensions, ccam réécrit les fichiers `installed_plugins.json` et `known_marketplaces.json` du groupe (`RewritePluginPaths`), car Claude Code enregistre des chemins de plugins absolus qui peuvent passer par le dossier de ce compte.
4. Un compte géré (dont ccam a créé le dossier, par opposition à un compte enregistré avec `ccam import`) est ensuite déconnecté : `claude auth logout` s'exécute dans son dossier, si bien qu'aucune connexion valide n'attend dans la corbeille (ni, sous macOS, dans un élément du trousseau nommé d'après le dossier, dont hériterait un futur compte du même nom). Cette étape est sautée pour un compte à clé d'API et quand le dossier d'un autre compte contient la même connexion, que la déconnexion révoquerait aussi. Une déconnexion qui échoue est signalée, et la suppression continue.
5. ccam déplace le dossier vers `~/.local/share/ccam/trash/<name>-<time>/` (`profile.Trash`), le compte est retiré de `config.toml`, les chemins de `nameState` sont supprimés pour qu'un futur compte du même nom n'en hérite pas, et le lien du raccourci est supprimé : `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

Si `config.toml` ne peut finalement pas être enregistré, le dossier revient de la corbeille. Pour récupérer un compte supprimé, remettez le dossier en place depuis la corbeille et reconnectez-vous.

`ccam rename` renomme le dossier et réécrit les chemins de plugins dans son `plugins/` de l'ancien dossier vers le nouveau. Il déplace les chemins de `nameState` vers le nouveau nom et conserve le statut de compte par défaut s'il l'avait. Il met aussi à jour les notes de session, pour que `ccam next --continue` les trouve toujours. La commande raccourci garde son nom. Si la réécriture des plugins échoue, ou si `config.toml` ne peut pas être enregistré, le dossier reprend son ancien nom.

Sous macOS, et sous Windows quand Claude Code conserve les connexions dans le Gestionnaire d'identification, le dossier lui-même garde son ancien nom, car la connexion et les connexions MCP sont stockées sous le chemin du dossier, et un dossier déplacé serait déconnecté. La configuration enregistre alors le dossier comme `folder` à côté du nouveau nom, et ccam affiche `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

`ccam import` enregistre un dossier là où il se trouve, comme `dir` dans `config.toml`. Il ne le déplace ni ne le copie jamais : un jeton de renouvellement copié ferait que les copies se déconnectent mutuellement.

Sans dossier en argument, il cherche les dossiers `~/.claude` et `~/.claude-*` non enregistrés et pose la question pour chacun. Les noms viennent du dossier (`~/.claude` devient `main`, `~/.claude-work` devient `work`, `~/.claude-2` devient `acct2`). Si ccam ne peut pas déduire un nom, ou si le nom est déjà pris, il liste le dossier avec la commande `--name` à exécuter. Un dossier trouvé qui échoue est signalé par `Skipped <path>: <reason>` et les autres sont importés.

ccam pose d'abord toutes ses questions, puis recharge `config.toml` et enregistre les dossiers, si bien que les modifications enregistrées entre-temps par d'autres commandes ccam sont conservées. Les lignes `✓` viennent après l'enregistrement.

Sans `CLAUDE_CONFIG_DIR`, Claude Code conserve `.claude.json` dans votre dossier personnel à côté de `~/.claude` ; quand vous importez `~/.claude` et qu'il n'a pas de `.claude.json` à lui, ccam crée donc le lien `~/.claude/.claude.json` vers `~/.claude.json`. `claude` seul et `ccam main` voient alors le même compte.

Si `~/.claude` a son propre `.claude.json` à côté de `~/.claude.json`, ccam laisse les deux et le signale. `ccam main` utiliserait le fichier intérieur et `claude` seul le fichier extérieur, ce qui sépare leurs informations de connexion et leurs serveurs MCP. Pour n'en garder qu'un, remplacez le fichier intérieur par un lien ([doctor](/features/doctor) le signale aussi).

Quand le dossier contient la même connexion qu'un compte enregistré (un dossier copié), ccam le signale : celui qui la renouvelle en premier déconnecte l'autre, exécutez donc `ccam login <name>` avant de l'utiliser.

Les comptes importés n'appartiennent à aucun groupe tant que vous ne les ajoutez pas à l'un d'eux (`ccam group join <group> <name>`), et aucun raccourci n'est créé.

`ccam completion bash|zsh` affiche un script de complétion qui complète aussi les noms de compte.

## Sécurité et cas particuliers

- `ccam rename` et `ccam remove` refusent d'agir pendant que le compte s'exécute (`refuseRunning`), avec ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``, car la session continue d'écrire dans le dossier. La vérification couvre Claude Code quelle que soit la façon dont il a été lancé : via ccam, `claude` seul, un IDE, les sessions en arrière-plan et leur démon.
- `ccam rename`, `ccam remove` et `ccam logout` refusent pendant qu'un autre processus ccam renouvelle ou modifie cette connexion : ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- Pendant qu'un changement de groupe est en cours dans un autre processus ccam, `ccam rename` et `ccam remove` refusent avec `another ccam process is changing the groups right now; try again when it is done` ; après un changement arrêté à mi-chemin, avec ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.
- `ccam remove` sur un compte importé conserve le dossier et sa connexion, pour `claude` seul. Si le compte fait partie de groupes, il les quitte d'abord comme le fait `ccam group leave` : chaque élément partagé devient dans le dossier une copie de celui du groupe, et les serveurs MCP du groupe vont dans son `.claude.json`. Il est ensuite désenregistré : `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. Les sauvegardes de son `.claude.json` faites par ccam vont dans `<folder>/.ccam-backup/mcp-<time>` ; si elles ne peuvent pas être déplacées, elles restent dans `~/.local/state/ccam/mcp/<name>/backups` et le message le précise. Votre `~/.claude` n'est pas déplacé dans la corbeille.
- Deux comptes ne peuvent pas avoir le même dossier : un `config.toml` qui nomme deux fois un dossier ne se charge pas.
- `ccam add` refuse un nom dont le dossier existe déjà dans `profiles/`, et renvoie vers `ccam import`. Sous macOS, il s'arrête aussi quand le trousseau contient encore une connexion ou des connexions MCP pour le chemin du nouveau dossier, laissées par une ancienne version de ccam, et nomme la commande `security delete-generic-password` qui les supprime. Si `config.toml` ne peut pas être enregistré, le nouveau dossier est de nouveau supprimé.
- `ccam import` refuse un dossier déjà enregistré, un chemin qui n'est pas un dossier et un nom déjà pris. ccam avertit aussi quand le dossier est une copie de `~/.claude` (la même connexion), même si `~/.claude` lui-même n'est pas enregistré.
- `ccam logout` demande d'abord quand le compte est en cours d'exécution ; `ccam login` ne vérifie pas.
- Si le relais ne peut pas joindre Claude Code, par exemple une fois la connexion terminée, il affiche `ccam: could not pass the sign-in to Claude Code:` avec la raison. Une erreur HTTP signifie généralement que le code a expiré ; relancez la connexion.

## Voir aussi

- [Lancer Claude Code](/features/launch)
- [Commandes raccourcis](/features/shortcuts)
- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [Utiliser plusieurs comptes Claude Code sur un même ordinateur](/guides/multiple-claude-code-accounts)
