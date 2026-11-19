---
title: "Référence des commandes"
metaTitle: "Référence des commandes ccam"
description: "Toutes les commandes et options de ccam, de ccam add et ccam <name> à ccam usage, ccam next et ccam doctor."
updated: 2026-11-19
order: 3
translation: machine
sourceHash: 28e1ad1dcbad
---

Placez `--dry-run` avant une commande pour voir ce qu'elle ferait sans l'exécuter, par exemple `ccam --dry-run work`. Cela fonctionne au démarrage de Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) et avec `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` et `doctor --fix`. Les commandes qui ne font que lire, comme `list`, `usage`, `share` et `ext`, l'acceptent et ne modifient rien de toute façon. Les autres la refusent au lieu de s'exécuter, y compris `group rename`.

<!-- commands -->

## Lancer Claude Code

`ccam <name|number> [arguments]` lance Claude Code avec ce compte. ccam ne lit pas les arguments qui suivent le nom ; il les transmet à `claude` dans le même ordre, et l'entrée standard, le code de sortie et les signaux se comportent exactement comme avec `claude` lui-même. Un `--` après le nom est facultatif.

`ccam` avec des arguments mais sans nom de compte, comme `ccam -p "summarise"`, utilise le compte défini pour le dossier courant, sinon le compte par défaut. Définissez-les avec `ccam folder <name> [folder]` et `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` avant le nom du compte s'applique aux exécutions `-p`/`--print`. Quand le compte atteint sa limite d'utilisation, ccam relance la même commande avec le compte connecté qui a le plus de marge et n'affiche que la sortie de l'exécution qui a abouti. Voir [Comptes par dossier](/features/folder) et [Relancer à la limite](/features/next#run-again-at-the-limit).

**Sièges Team et Enterprise.** La sélection automatique de compte (`ccam next`, `--fallback`, `ccam pick`, la touche `n` du tableau de bord et, une fois que vous l'activez, [`auto_switch`](/features/next#auto-switch)) passe aussi d'un siège à l'autre au sein d'une même organisation Claude Team ou Enterprise, ce qui peut enfreindre les conditions d'Anthropic. Lisez la [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) et renseignez-vous auprès de votre organisation et d'Anthropic avant de vous y fier.

Un raccourci comme `ccwork` lance son compte de la même façon : `ccwork --resume` fait la même chose que `ccam work --resume`. Voir [Raccourcis](/docs/commands#shortcuts) plus bas.

## Ajouter et supprimer des comptes

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

Sans dossier, `ccam import` cherche `~/.claude` et `~/.claude-*` et pose la question pour chacun ; un dossier qui échoue est signalé par `Skipped <path>: <reason>` et les autres sont enregistrés. Si le nom d'un dossier est déjà pris, importez ce dossier seul avec `ccam import <folder> --name <name>`. Les dossiers importés sont enregistrés là où ils se trouvent ; ccam ne déplace ni ne copie jamais une connexion.

`ccam login --expired` connecte aussi les comptes dont Claude Code a effacé la connexion. `ccam remove` déconnecte un compte géré (dont ccam a créé le dossier, par opposition à un compte enregistré avec `ccam import`) avant que son dossier parte à la corbeille, sauf si un autre compte détient la même connexion. Supprimer un compte importé le retire de ses groupes et laisse dans son dossier une copie de chaque élément partagé. Sous macOS, `ccam rename` conserve le nom du dossier, car le trousseau garde la connexion sous le chemin du dossier.

`ccam add --api` ajoute un compte qui utilise une clé d'API au lieu d'une connexion Claude ; ccam lit la clé à chaque lancement et ne la stocke jamais. Voir [Comptes à clé d'API](/features/api).

`ccam rename`, `ccam remove` et les commandes `ccam group` refusent de modifier un compte tant que Claude Code s'exécute avec ce compte, quelle que soit la façon dont il a été lancé (sessions en arrière-plan comprises). Quittez d'abord la session. Elles refusent aussi tant qu'un changement de groupe est inachevé : ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` affiche chaque compte avec son numéro d'emplacement, son raccourci, son adresse e-mail et son dossier. Son tableau `groups` liste les groupes du compte, et son objet `share` associe chaque élément (`mcp`, `settings`, `memory`, `extensions`, `history`) au groupe dont il provient, ou à `local`.

## Raccourcis

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Un raccourci est une commande que vous nommez vous-même et qui lance Claude Code avec un compte donné. Comme `ccam work`, il transmet tous les arguments tels quels à Claude Code, si bien que `ccwork --resume` et `ccwork -p "hi"` fonctionnent comme prévu. C'est un lien vers ccam dans `~/.local/bin`.

Les raccourcis sont facultatifs. L'assistant de `ccam add` demande s'il faut en créer un, `ccam add <name> --shortcut <command>` en crée un sans rien demander, et `ccam list` affiche le raccourci de chaque compte, ou `-` s'il n'en a pas.

Les noms utilisent les caractères a-z, 0-9, `.`, `-` et `_`. Un nom ne peut pas être une commande ccam, le raccourci d'un autre compte, ni le nom d'un autre programme présent dans votre `PATH`, ni un fichier de `~/.local/bin` que ccam n'a pas créé, comme votre propre copie ou lien de ccam. Un raccourci modifié à la main dans `config.toml` qui enfreint ces règles, qui vaut `claude` ou `ccam`, ou qui est attribué à deux comptes empêche `config.toml` de se charger, en nommant le compte.

`ccam rename` conserve le nom du raccourci, et `ccam remove` supprime le raccourci avec le compte. `ccam shims install` recrée les liens de tous les comptes, par exemple après avoir déplacé ccam ; un fichier qui gêne un raccourci n'empêche pas les autres.

## Partage

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Un groupe partage des éléments entre ses comptes ; un compte peut appartenir à plusieurs groupes, mais chaque élément ne provient que d'un seul groupe. Voir [Partager avec des groupes](/features/share).

## Copier une configuration existante

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` copie `settings.json`, `CLAUDE.md`, les extensions et les serveurs MCP utilisateur d'un dossier Claude Code dans un groupe, pour que chaque compte qui partage ces éléments en dispose. Le dossier est seulement lu.

- Les options `--settings`, `--memory`, `--extensions`, `--mcp` et `--history` ne copient que les éléments indiqués. Sans elles, une liste à cocher vous permet de choisir dans un terminal, avec tout coché sauf l'historique des conversations ; ailleurs, tout est copié sauf l'historique des conversations.
- `--mcp` fonctionne aussi quand `~/.claude` est vide ou absent mais que `~/.claude.json` contient des serveurs.
- `-g` ou `--group` choisit le groupe.
- `-y` ou `--yes` copie sans rien demander et garde la copie du groupe quand un fichier diffère.

Les connexions MCP ne sont pas copiées, car une connexion copiée déconnecte ses autres détenteurs ; le plan nomme les comptes qui se connectent avec `/mcp`. Voir [Copier une configuration existante](/docs/sharing#copy-an-existing-setup).

## Serveurs MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` déplace dans son groupe les serveurs que seul ce compte possède, y compris ceux qu'il n'a que pour un seul dossier (`claude mcp add` seul), qui se chargent alors dans tous les dossiers. Un serveur dont le groupe utilise déjà le nom pour une autre définition est ignoré ; ajoutez `--force` pour remplacer la copie du groupe, que ccam sauvegarde d'abord.

## Utilisation et changement de compte

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

- `next --auto`, `--fallback` et `ccam pick` choisissent de la même façon : d'abord sous `auto_max_usage` (`config.toml`, 100 par défaut), selon `auto_strategy` (`headroom`, `expiring` ou `balance`), en comptant la fenêtre hebdomadaire de `--model`. Ils ne choisissent jamais un compte à 100 % (une fenêtre que Claude Code a signalée pleine compte pour 100 % jusqu'à sa réinitialisation) ni une connexion par clé d'API, et restent au sein de l'organisation du compte de départ (`--cross-org` ou `cross_org = true` pour chercher plus loin).
- `--wait` prend sa valeur après `=` (`--wait=2h`) ; écrit avec une espace, il est refusé avec une indication. Quand aucun ne peut être choisi, ils quittent avec le code 75 ou, avec `--wait`, attendent la réinitialisation du premier compte.
- `ccam next --continue` transporte une conversation même entre des comptes qui ne partagent pas l'historique : Claude Code la copie alors dans le nouveau compte. Avec `--brief`, le nouveau compte démarre plutôt une nouvelle conversation, à partir d'un résumé de la précédente que ccam rédige sans modèle ([`--brief`](/features/next#brief)).
- Contrairement à `ccam <name>`, `ccam next` lit ses propres options, donc les options de Claude Code se placent après `--` ; sans lui, ccam s'arrête avec `unknown flag`.
- `ccam agent` exécute une tâche avec `claude -p` dans ce dossier avec un autre compte que celui depuis lequel il est appelé, dans le mode de permission et le niveau d'effort de la session appelante, et affiche la dernière réponse de Claude. À la limite de ce compte, la conversation continue avec un autre ([Des sous-agents sur d'autres comptes](/features/agent)).
- Dans Claude Code, la commande `/ccam:switch [account]` du plugin passe une conversation à un autre compte, et `ccam hook` est ce que lancent ses hooks.
- Avec `auto_switch = true` dans `config.toml` (Linux et macOS ; désactivé par défaut), le hook de limite met aussi fin à une session interactive lancée par ccam et la poursuit sur un autre compte qui a de la marge, mais uniquement un compte du même type : un autre siège de la même organisation, ou une autre offre personnelle ([`auto_switch`](/features/next#auto-switch)).

Voir [Passer au compte suivant](/features/next), [Choisir un compte pour les scripts](/features/pick), [Des sous-agents sur d'autres comptes](/features/agent) et [Les sessions en arrière-plan de tous les comptes](/features/agents).

## Maintenance

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

`ccam update` vérifie la somme de contrôle SHA-256 de la version avant de remplacer ccam, et `--version v0.2.0` installe une version précise. Une copie de ccam installée depuis le paquet `.deb` se met à jour avec apt. Voir [Mettre à jour](/docs/install#update).

`ccam uninstall` liste ce qu'il supprime et demande d'abord confirmation : tous les comptes, comme le fait `ccam remove`, les raccourcis, les dossiers de ccam et la commande `ccam`. Un dossier importé reste là où il se trouve. Voir [Désinstaller](/docs/install#uninstall).

## Fichiers

| Chemin | Contenu |
| --- | --- |
| `~/.config/ccam/config.toml` | Les comptes, les groupes et les réglages (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). Les réglages se placent en haut, au-dessus du premier `[[group]]`, `[[profile]]` ou `[[folder]]` : un réglage placé sous une table appartient à cette table, ccam le signale donc à chaque lancement et dans `ccam doctor`, et les commandes qui enregistrent la configuration refusent tant qu'il n'a pas été déplacé. Les commentaires ne sont pas conservés quand ccam enregistre le fichier |
| `~/.local/share/ccam/profiles/<name>/` | Le `CLAUDE_CONFIG_DIR` de chaque compte |
| `~/.local/share/ccam/groups/<group>/` | Paramètres, serveurs MCP, extensions et historique partagés |
| `~/.local/state/ccam/` | Cache d'utilisation, relevés de limite, verrous, notes et autres états (voir la liste sous le tableau) |
| `~/.local/bin/<shortcut>` | Raccourcis comme `ccwork`, des liens vers ccam |

Le dossier d'état contient :

- le cache d'utilisation et les relevés de limite (`usage/<name>.json`, `usage/<name>.limit.json`) ;
- les heures de dernière utilisation ;
- les marqueurs de session et les verrous (`locks/`) ;
- les raccourcis créés par ccam (`shortcuts.json`) ;
- un changement de groupe inachevé (`group-change.json`) ;
- les notes que laissent `/ccam:switch` et le hook de limite (`handoff/<account>/<session>.json`) ;
- le moment où `ccam pick` ou `auto_switch` a choisi chaque compte pour la dernière fois (`picks/<name>` ; un choix compte pendant 2 minutes) ;
- les résumés de `--brief`, conservés 7 jours (`brief/<session>.md`).

Définir `XDG_CONFIG_HOME`, `XDG_DATA_HOME` ou `XDG_STATE_HOME` change l'emplacement de ces dossiers. Sous Windows, `~` correspond à `%USERPROFILE%` et les raccourcis sont des `<shortcut>.exe`.
