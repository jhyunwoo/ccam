---
title: "Passer au compte suivant"
description: "ccam next ouvre Claude Code avec un autre compte ; --auto prend le plus de marge, --continue y emporte la conversation et --fallback relance les tâches -p."
updated: 2026-11-05
order: 8
translation: machine
sourceHash: "a320c160caef"
---

`ccam next` ouvre Claude Code avec un autre compte sans que vous ayez à le nommer. Utilisez-le quand un compte atteint sa limite d'utilisation. Il y a quatre façons de changer de compte :

- `ccam next` seul prend le compte qui suit celui utilisé en dernier, et `--auto` prend celui qui a le plus de marge. `--continue` poursuit la dernière conversation de ce dossier sur le nouveau compte, que les deux partagent l'historique ou non.
- Dans Claude Code, le [plugin ccam](#inside-claude-code) vous prévient quand une limite arrête une session, et `/ccam:switch` passe la conversation à un autre compte.
- Pour les scripts qui exécutent `claude -p`, `ccam --fallback` fait le changement de lui-même (voir [Relancer à la limite](#run-again-at-the-limit)).
- Avec [`auto_switch`](#auto-switch) activé (Linux et macOS), ccam déplace aussi une session interactive que Claude Code arrête à sa limite.

**Sièges Team et Enterprise.** La sélection automatique de compte (`ccam next`, `--fallback`, `ccam pick`, la touche `n` du tableau de bord et, une fois que vous l'activez, [`auto_switch`](/features/next#auto-switch)) passe aussi d'un siège à l'autre au sein d'une même organisation Claude Team ou Enterprise, ce qui peut enfreindre les conditions d'Anthropic. Lisez la [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) et renseignez-vous auprès de votre organisation et d'Anthropic avant de vous y fier.

## Utilisation

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

ccam affiche le changement sous la forme `→ <from> → <to>` sur la sortie d'erreur standard, avec la raison d'un choix `--auto`, puis lance Claude Code avec le nouveau compte de la même façon que `ccam <name>`. `<from>` est le compte utilisé en dernier ou, avec `--continue`, le compte auquel appartient la conversation. Les arguments après `--` sont transmis tels quels à `claude`.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--auto` | Choisir le compte connecté selon l'utilisation (voir plus bas), en ignorant les comptes à 100 % | désactivé |
| `--max-usage <1-100>` | Avec `--auto` : préférer les comptes sous ce niveau d'utilisation | `auto_max_usage`, sinon 100 |
| `--strategy headroom\|expiring\|balance` | Avec `--auto` : l'utilisation hebdomadaire la plus faible parmi les comptes à 3 points au plus de l'utilisation sur 5 heures la plus basse, la réinitialisation sur 5 heures la plus proche ou le plus de quota hebdomadaire restant par jour | `auto_strategy`, sinon `headroom` |
| `--cross-org` | Avec `--auto` : choisir aussi des comptes d'une autre organisation (voir plus bas) | `cross_org`, sinon désactivé |
| `--wait[=<duration>]` | Avec `--auto` : quand aucun compte n'a de marge, attendre la réinitialisation du premier, au plus cette durée. La valeur se place après `=` : `--wait 2h` est refusé avec l'indication `--wait takes its value after '=': --wait=2h` | désactivé ; `5h` s'il est donné seul |
| `-c`, `--continue` | Poursuivre la dernière conversation de ce dossier sur le nouveau compte (voir plus bas) | désactivé |
| `--brief` | Avec `--continue` : démarrer une nouvelle conversation avec le nouveau compte à partir du résumé que ccam rédige de la précédente, au lieu de la reprendre entièrement ([voir plus bas](#brief)) | désactivé |
| `-- <claude args>` | Arguments transmis à `claude` après ceux de ccam | aucun |
| `--dry-run` | Pris en charge : affiche `→ <from> → <to>` et la ligne de commande qui serait exécutée, sans lancer Claude Code ni enregistrer de lancement | désactivé |

`--auto` et `--continue` fonctionnent ensemble. Avec `--continue`, les arguments `--resume` sont placés avant les vôtres.

Les options propres à Claude Code, comme `--dangerously-skip-permissions` ou `--model`, se placent après `--`. Contrairement à `ccam <name>`, qui transmet tout ce qui suit le nom, `ccam next` lit lui-même ses arguments : sans le `--`, il s'arrête donc avec `ccam: unknown flag: --dangerously-skip-permissions`.

## Fonctionnement

Chaque lancement via ccam écrit l'heure dans `~/.local/state/ccam/last-used/<name>`. `ccam next` lit ces fichiers, vérifie quels comptes sont connectés, puis choisit.

1. ccam construit un `next.Candidate` par compte à partir de `config.toml` (`cli.candidates`) : le compte, l'heure de son dernier lancement (`launcher.LastUsed`) et s'il dispose d'une connexion valide (`creds.Has` : une connexion Claude dans `.credentials.json` ou, sous macOS, dans le trousseau, qui n'a pas entièrement expiré ni été effacée par Claude Code). Un dossier qui utilise plutôt une clé d'API (une connexion Console, ou un compte ajouté avec `ccam add --api`) n'a pas de limites d'offre et n'est jamais choisi.
2. Sans `--auto`, ccam trie les comptes par numéro d'emplacement et trouve celui lancé le plus récemment (`next.Rotate`). C'est le compte que vous quittez. Si aucun compte n'a encore été lancé, c'est l'emplacement le plus bas.
3. À partir de là, il avance dans l'ordre des emplacements, en revenant du dernier au premier, et renvoie le premier compte connecté.
4. ccam prépare et lance le compte choisi exactement comme le ferait `ccam <name>` (`launcher.LaunchWith`) : il vérifie les liens du compte, applique le partage MCP, enregistre la session et remplace ccam par `claude` via `syscall.Exec`.

Le parcours commence au compte utilisé en dernier et s'arrête au premier compte connecté qui le suit :

<!-- motion:next-rotate -->

`work` a été lancé le plus récemment, donc le parcours commence après lui. `side` n'a pas de connexion et est ignoré, et `lab` est choisi. Si `lab` avait lui aussi été déconnecté, le parcours serait revenu à `personal` (`% len` dans le diagramme). La rotation ne renvoie jamais le compte de départ : sans autre compte connecté, elle s'arrête avec une erreur.

### --auto

`--auto` décide selon l'utilisation plutôt que selon l'ordre des emplacements. ccam garde les comptes connectés, demande leur utilisation à tous en même temps de la même façon que `ccam usage` (`Fetcher.FetchAll`), et choisit selon les règles ci-dessous (`cli.pickAuto`, `usage.Policy.Choose()`) :

<!-- motion:next-auto -->

Un compte est jugé sur ses fenêtres de 5 heures et hebdomadaire et, quand `--model` figure parmi les arguments de Claude Code (après `--`), aussi sur la fenêtre hebdomadaire de ce modèle (`opus`, `sonnet`, `haiku` ou `fable` dans le nom). ccam répartit les comptes en trois niveaux (`Policy.Choose()`) :

1. **Marge** : chaque fenêtre sous `auto_max_usage` (`config.toml`, 100 par défaut). La stratégie choisit parmi eux : `headroom` (par défaut) prend, parmi les comptes à 3 points au plus de l'utilisation sur 5 heures la plus basse, celui qui a le moins consommé de sa semaine, si bien qu'un compte à 12 % l'emporte sur un compte à 10 % quand il a moins consommé de sa semaine, mais pas un compte à 14 % ; `expiring` prend la fenêtre de 5 heures qui se réinitialise le plus tôt, afin que son quota restant serve avant d'être perdu ; `balance` prend le compte auquel il reste le plus de quota hebdomadaire par jour jusqu'à la réinitialisation de sa semaine, en décomptant chaque session ccam en cours, afin que le quota hebdomadaire ne reste pas inutilisé et que les sessions parallèles se répartissent.
2. **Au-dessus** : sous 100 %, mais à `auto_max_usage` ou au-delà. Utilisés seulement quand aucun compte n'a de marge ; celui que `headroom` prendrait parmi eux passe en premier.
3. **Plein** : une fenêtre à 100 %, ou une utilisation sur 5 heures inconnue. Jamais choisis.

Une fenêtre compte pour 100 % jusqu'à sa réinitialisation quand Claude Code y a refusé une requête : le hook du [plugin ccam](#inside-claude-code) l'enregistre (`~/.local/state/ccam/usage/<name>.limit.json`), quels que soient les chiffres enregistrés. Un compte qui vient d'atteindre sa limite n'est donc pas choisi de nouveau avant sa réinitialisation, même quand ses chiffres sont anciens.

Les réglages valables pour tous les comptes se placent en haut de `config.toml`, au-dessus du premier `[[group]]`, `[[profile]]` ou `[[folder]]`. Un réglage placé sous l'un d'eux appartient à cette table : ccam vous en avertit à chaque lancement, et les commandes qui enregistrent la configuration refusent tant qu'il n'a pas été déplacé.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organisations.** Les choix automatiques restent parmi les comptes du même type que celui que vous quittez : des sièges de la même organisation Team ou Enterprise, ou des offres personnelles (Pro, Max), afin qu'un travail commencé sur un compte d'entreprise ne se poursuive pas sur un compte privé, ni l'inverse. ccam connaît l'organisation de chaque compte grâce à sa connexion (`ccam login` l'enregistre ; pour les connexions plus anciennes, elle est lue dans `.claude.json`) ; un compte dont ccam ne peut pas déterminer l'organisation peut être associé à n'importe quel compte. Deux dossiers connectés au même compte Claude partagent ses limites : l'un ne prend donc jamais le relais de l'autre. `--cross-org`, ou `cross_org = true` dans `config.toml`, lève la règle d'organisation, mais après [`/ccam:switch`](#inside-claude-code), un changement vers une autre organisation attend toujours votre accord.

Quand rien ne peut être choisi, ccam quitte avec le code 75 (`EX_TEMPFAIL`), pour qu'un script distingue « réessayer plus tard » d'un échec. Avec `--wait`, il attend plutôt que la fenêtre bloquante du premier compte plein se réinitialise (30 secondes après l'heure de réinitialisation, quand le point de terminaison montre la nouvelle fenêtre), affiche `ccam: no account has room; waiting until 15:10` et choisit à nouveau ; quand aucun compte n'indique quand il se réinitialise, il regarde de nouveau au bout de 10 minutes. Si rien ne se libère avant la fin de l'attente, il quitte avec 75.

La consultation suit les règles de `ccam usage` : les chiffres de moins de 60 secondes sont réutilisés, la consultation de chaque compte expire au bout de 3 secondes, un compte qui ne répond pas, ou dont le jeton d'accès a expiré, est jugé sur ses chiffres enregistrés, et un relevé de limite prend le pas sur eux. Une fenêtre enregistrée dont l'heure de réinitialisation est passée compte pour 0 %, si bien qu'un compte dont la fenêtre a été réinitialisée entre-temps n'est pas écarté. `--auto` ne renouvelle jamais les jetons. Contrairement à la rotation simple, `--auto` peut choisir le compte utilisé en dernier s'il a encore le plus de marge.

### --continue

1. ccam cherche la dernière conversation de ce dossier dans chaque compte : parmi les fichiers `<session>.jsonl` de `projects/<folder>/` (le nom du dossier est le dossier de travail, liens résolus, dont chaque caractère autre qu'une lettre ou un chiffre devient `-`, comme le nomme Claude Code), celui dont l'invite ou la réponse est la plus récente. Les sessions qui n'en ont pas, comme une session qui n'a exécuté que `/mcp`, `/usage` ou `/status`, ou un fichier vide, sont ignorées. Si une session a laissé une note (une limite d'utilisation ou `/ccam:switch`, voir plus bas) et qu'aucune conversation n'a d'invite ou de réponse plus de 2 minutes après elle, c'est cette session précise qui est utilisée ; les lignes que Claude Code écrit à `/exit` ne comptent pas.
2. Le compte auquel elle appartient est celui que vous quittez, et le nouveau compte est choisi comme ci-dessus, au sein de la même organisation. Quand plusieurs comptes partagent cet historique, elle appartient à celui dont le `.claude.json` indique qu'il a exécuté la conversation en dernier (`lastSessionId`), sinon à celui utilisé en dernier.
3. Si le nouveau compte partage cet historique (le même dossier de groupe), ccam lance `claude --resume <session id>`, précisément cette conversation (`--continue` ouvrirait la plus récente du dossier, quelle qu'elle soit).
4. Sinon, il lance `claude --resume <the conversation's file> --fork-session` : Claude Code copie toute la conversation, sous un nouvel ID de session, dans l'historique propre au nouveau compte, et l'original reste tel quel. Dès lors, le `--continue` du nouveau compte la trouve.
5. Sauf si vous donnez `--agent` après `--`, ccam transmet aussi l'agent avec lequel la conversation s'exécute : le paramètre `agent` du `.claude/settings.local.json` ou du `.claude/settings.json` du dossier, ou du `settings.json` du compte, sinon l'agent que nomme la dernière ligne `agent-setting` de la conversation. Claude Code ne rétablit pas l'agent quand il copie une conversation (étape 4), et continuerait avec son agent par défaut, avec tous les outils. Quand le nouveau compte ne partage pas les agents, plugins et paramètres de l'ancien (les mêmes dossiers `agents` et `plugins` et le même `settings.json`, via un groupe), ccam ajoute une ligne après celle qui concerne le changement : `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

ccam indique la taille de la conversation : le nouveau compte la relit entièrement avec sa première réponse. Les caches de prompt ne sont jamais partagés entre organisations, et chaque offre personnelle est une organisation à part : quand les deux comptes appartiennent à des organisations différentes, ccam ajoute donc `prompt caches are not shared between organizations, so the first reply reads it all again`. Quand l'ancien compte est connecté à des serveurs MCP auxquels le nouveau ne l'est pas, ccam les nomme : `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. Partager `history` dans un groupe évite la copie, conserve une seule conversation et partage aussi les points de contrôle à partir desquels `/rewind` restaure le code (`file-history/`).

### --brief

`--continue` reprend toute la conversation, et le nouveau compte la relit entièrement avec sa première réponse. `ccam next --continue --brief` démarre plutôt une nouvelle conversation, à partir d'un résumé que ccam rédige de l'ancienne sans interroger de modèle :

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

Le résumé contient, dans cet ordre et seulement quand il y a quelque chose à dire :

1. le titre de la conversation ;
2. le chemin de sa transcription, dans laquelle Claude peut chercher quand il a besoin d'un détail ;
3. votre première demande, quand il y en a eu plus de quatre ;
4. le résumé que Claude Code a rédigé lors du dernier compactage de la conversation ;
5. vos quatre dernières demandes, avec les commandes `!` et les commandes slash telles que vous les avez tapées ;
6. la dernière réponse de Claude ;
7. la liste de tâches ;
8. les 10 derniers fichiers modifiés par Claude ;
9. le `git status --short --branch` du dossier, jusqu'à 800 caractères.

La liste de tâches garde au plus 10 éléments : les éléments terminés sont retirés en premier, puis les étapes les plus lointaines, si bien que l'étape en cours et les suivantes restent. Les parties longues sont coupées au milieu, et le résumé entier reste sous 16 000 caractères.

ccam l'enregistre sous `~/.local/state/ccam/brief/<session>.md`, lisible par vous seul, et le `--brief` suivant supprime les résumés de plus de 7 jours. Claude Code démarre avec une seule invite, et joint le fichier qu'elle mentionne : `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` Une invite que vous donnez après `--` remplace la dernière phrase (`Then: <your prompt>`), et vos autres options après `--` la suivent.

Claude Code termine une mention à un `"` et la coupe à un `#` : quand le chemin du résumé contient l'un des deux, ccam le mentionne par son chemin relatif au dossier de travail ; quand ce chemin en contient un aussi, l'invite nomme le fichier, et Claude le lit avec son outil Read, qui demande d'abord l'autorisation car le fichier se trouve hors des dossiers de travail (dans une exécution `-p` en mode Manual, Claude Code refuse la lecture et Claude continue sans le résumé).

`ccam --dry-run next --continue --brief` affiche le résumé et n'enregistre rien. L'ancienne conversation reste telle quelle, et la nouvelle ne sait que ce que dit le résumé : utilisez donc `--brief` quand la tâche compte plus que chaque détail du chemin parcouru. Sans `--continue`, ccam s'arrête avec `--brief works with --continue` ; quand il ne peut pas lire la conversation ou n'y trouve rien à résumer, il s'arrête et indique que `ccam next --continue` sans `--brief` la reprend entièrement. Le [relais](#relay) et [`auto_switch`](#auto-switch) reprennent toujours la conversation entière. L'option `--brief` propre à Claude Code se place après `--`, comme ses autres options.

## Dans Claude Code

Le [plugin ccam](/docs/install#claude-code-plugin) ajoute deux hooks (Claude Code 2.1.78 ou plus récent) :

- Quand un tour s'arrête à une limite d'utilisation, il note la conversation et affiche une notification de bureau dans les terminaux qui le permettent (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty) : ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  Il connaît chaque message que Claude Code considère comme une limite d'utilisation : les limites sur 5 heures, hebdomadaire, simple et par modèle (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` pour l'allocation d'un modèle), les crédits d'utilisation et les plafonds de dépenses, et une organisation qui a suspendu le compte.

  Quand Claude Code a écrit la requête refusée dans la conversation, le hook enregistre aussi la fenêtre pleine jusqu'à sa réinitialisation ([le relevé de limite](/features/usage#the-limit-record)), sauf pour une suspension décidée par une organisation, que ses administrateurs peuvent lever à tout moment. Un bref ralentissement imposé par le serveur ou la limite de débit d'une clé d'API n'est pas une limite d'utilisation et n'est pas pris en compte, pas plus que les sessions sans surveillance (exécutions `-p` et sessions en arrière-plan). Un sous-agent ou un agent de workflow qui s'arrête à une limite d'utilisation compte comme l'arrêt de sa session, même quand la conversation principale ne fait que l'attendre : Claude Code écrit cette requête refusée dans le fichier propre à l'agent, dans le dossier `subagents` de la session, et le hook la lit là. Quand plusieurs agents s'arrêtent en même temps, la session reçoit une seule note et une seule notification.

  Claude Code 2.1.234 et les versions ultérieures attendent aussi dans la session et reprennent d'eux-mêmes après une réinitialisation de la fenêtre de 5 heures ; la notification vous laisse le choix de changer de compte à la place, et avec [`auto_switch`](#auto-switch) activé (Linux et macOS), ccam change de compte pour vous quand un autre compte a de la marge. Quand aucun autre compte ne peut prendre le relais, parce qu'aucun autre compte n'appartient à la même organisation Team ou Enterprise, ou qu'aucun n'est une autre offre personnelle (avec `cross_org` activé, n'importe quel autre compte Claude le peut), la notification se termine à la place par `You have no other personal plans, so ccam cannot move this conversation.` ou `You have no other accounts in Acme, so ccam cannot move this conversation.`
- `/ccam:switch [account]` note la conversation et termine le tour avant toute requête : tapez `/exit`, puis exécutez `ccam next --continue` pour la poursuivre sur le compte que vous avez nommé (ou sur le suivant). Dans une session qui n'a pas encore de conversation, il indique qu'il n'y a rien à poursuivre. Seul ce que vous tapez l'exécute : depuis la version 0.3.2 du plugin, Claude ne peut pas exécuter `/ccam:switch`, et une tâche que Claude planifie (CronCreate, ScheduleWakeup ou `.claude/scheduled_tasks.json`) le transmet à Claude comme simple texte.

  Quand le compte que vous nommez est en dehors de l'organisation de la conversation (un siège d'une autre organisation Team ou Enterprise, une offre personnelle à la place d'un siège ou l'inverse, ou un compte à clé d'API d'un côté ou de l'autre), `ccam next --continue` demande d'abord, quoi qu'indique `cross_org` : `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. Il pose la même question avant un compte qu'il choisit lui-même quand vous n'en nommez aucun, et que `cross_org = true` lui permet de prendre dans une autre organisation : `ccam: the account picked for /ccam:switch is home, a personal plan, …`. Seuls `y` ou `yes` la déplacent ; sous Linux et macOS, les touches tapées avant la question sont d'abord ignorées. Toute autre réponse, ou l'absence de terminal où poser la question, laisse la conversation et la note là où elles sont : ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (après un choix, `ccam next --continue` choisit de nouveau). Comme pour les choix automatiques, un compte dont ccam ne peut pas déterminer l'organisation convient à n'importe quel compte.

La note se trouve dans `~/.local/state/ccam/handoff/<account>/<session>.json`, un fichier par session, si bien que deux sessions d'un même compte n'écrasent jamais la note l'une de l'autre. Elle est conservée pendant 12 heures et ne sert qu'une fois, au moment où la session suivante va démarrer ; un lancement qui échoue la laisse pour la tentative suivante. Elle appartient au dossier dans lequel la session a démarré, même après un `cd` dans le shell de Claude Code. `ccam rename` met à jour les noms qu'elle contient ; si le compte qu'elle nomme a été supprimé entre-temps, ccam le signale et choisit comme d'habitude.

### Relais

Avec `relay = true` en haut de `config.toml`, ccam continue de tourner à côté des sessions interactives au lieu de se remplacer par Claude Code, et après `/ccam:switch` il poursuit de lui-même la conversation de cette session (une note laissée par une autre session n'est pas prise en compte) : le hook termine Claude Code, et ccam exécute `ccam next --continue` avec la note de cette session. Sans compte nommé, il prend celui qui a le plus de marge, en tenant compte du modèle de la conversation, et quand `cross_org = true` permet que ce soit un compte d'une autre organisation, il pose la question ci-dessus avant d'y déplacer la conversation. Sous Windows, le hook ne peut pas terminer Claude Code : tapez `/exit`, et ccam poursuit la conversation. Le hook ne le termine pas non plus pour un compte en dehors de l'organisation de la conversation : il affiche `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`, et ccam pose la question ci-dessus une fois que vous avez tapé `/exit`. Cela coûte un processus ccam en attente par session, c'est pourquoi c'est désactivé par défaut.

La session relayée conserve ce que la reprise d'une conversation ne restaure pas : des options comme `--add-dir`, `--mcp-config`, les listes d'outils, `--settings`, `--plugin-dir` et les options d'invite système, le mode de permission de la conversation, le niveau d'effort de la session et l'agent avec lequel la conversation s'exécute (voir [`--continue`](#continue)). Claude Code restaure lui-même `--model` et `--name` à partir de la conversation. Tout le reste (`--worktree`, `--debug`, `--remote-control`, les options que ccam ne connaît pas) est nommé dans `ccam: not carried over: …`. Le changement ne pose jamais la question de l'historique partagé, et ce que vous avez tapé pendant son exécution est abandonné ; sa seule question est celle qui précède un compte d'une autre organisation.

Le relais ne poursuit la conversation que lorsque Claude Code s'est terminé comme `/ccam:switch` le lui a demandé : pas quand vous avez tapé `/exit` d'abord, appuyé sur Ctrl+C pendant qu'il se terminait, ou que ccam a reçu un raccrochage ou un autre signal, ni quand le terminal a disparu. La note reste alors, et `ccam next --continue` la poursuit plus tard. Une session lancée sous le relais d'une ancienne version de ccam, ou une session sans surveillance (exécutions `-p` et sessions en arrière-plan), n'est pas terminée par le hook ; elle reçoit à la place le message `/exit`. `claude remote-control` s'exécute sans relais.

### `auto_switch`

Avec `auto_switch = true` en haut de `config.toml` (Linux et macOS), ccam déplace de lui-même une session vers un autre compte quand Claude Code l'arrête à une limite d'utilisation : le hook du plugin choisit un autre compte qui a de la marge et termine Claude Code, et ccam, qui tourne à côté de la session, y poursuit la conversation, comme le fait le [relais](#relay) après `/ccam:switch`. Ce réglage est désactivé par défaut et expérimental.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam ne déplace une session que si toutes ces conditions sont réunies :

1. ccam l'a lancée dans un terminal et continue de tourner à côté d'elle. `auto_switch` active le [relais](#relay), si bien que c'est le cas pour les sessions interactives que vous lancez après l'avoir réglé. Une session lancée avant, ou sous un ancien ccam, reçoit une notification qui se termine par ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` Les sessions sans surveillance (exécutions `-p` et sessions en arrière-plan), `claude remote-control` et Claude Code lancé sans ccam ne sont jamais déplacés.
2. Le compte a le [plugin ccam](/docs/install#claude-code-plugin) activé et exécute Claude Code 2.1.214 ou plus récent, qui transmet l'ID de son processus au hook (`CLAUDE_PID`).
3. Claude Code a écrit la requête refusée dans la conversation, avec l'heure de réinitialisation, pour la limite de 5 heures, hebdomadaire, Opus ou Sonnet de ce compte. Un plafond de dépenses, l'épuisement des crédits d'utilisation, une suspension décidée par une organisation, l'allocation Fable (`You've reached your Fable limit.`), un ralentissement imposé par le serveur et la limite de débit d'une clé d'API ne déplacent jamais une session. La limite d'un sous-agent ou d'un agent de workflow compte, même quand la conversation principale ne fait que l'attendre ; s'il s'agit d'une limite hebdomadaire Opus ou Sonnet, la session n'est déplacée que si la conversation elle-même tourne sur ce modèle.
4. ccam connaît l'offre du compte, d'après son `.claude.json` ou ce qu'a enregistré `ccam login`. Sinon, la notification se termine par `ccam does not know personal's plan, so auto_switch leaves this session.`
5. La session n'est pas connectée à Remote Control. Sous un autre compte, elle ne se reconnecterait pas à l'application depuis laquelle vous la pilotez : elle reste donc, et la notification le signale.
6. La limite se réinitialise dans plus de 15 minutes. Pour une attente plus courte, ccam laisse la session à Claude Code : `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. Vous n'avez pas tapé d'invite depuis la limite.
8. Un autre compte a de la marge. ccam choisit comme le fait [`--auto`](#auto) (`auto_max_usage`, `auto_strategy`, le modèle de la conversation), parmi les comptes du même type que celui que vous quittez : un autre siège de la même organisation Team ou Enterprise, ou une autre offre personnelle (Pro, Max). `cross_org` et `--cross-org` ne s'appliquent pas ici. Le compte doit :
   - être connecté avec une connexion que Claude Code peut renouveler ;
   - ne pas être connecté au même compte Claude ;
   - ne pas utiliser de clé d'API ;
   - avoir une offre que ccam connaît ;
   - ne pas avoir les crédits d'utilisation activés, d'après la propre consultation d'utilisation de ccam ;
   - avoir accordé sa confiance à ce dossier dans Claude Code ;
   - si la conversation s'exécute avec un agent (le paramètre `agent` ou la conversation en nomme un), partager les agents du compte que vous quittez : ses dossiers `agents` et `plugins` et son `settings.json` doivent chacun être exactement les mêmes (via un groupe), ou absents des deux côtés.

   Quand ces trois dernières conditions écartent tous les comptes qui ont de la marge, la notification les nomme : `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` ou `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` Une fenêtre qui a un [relevé de limite](/features/usage#the-limit-record) compte comme pleine, et un compte vers lequel une autre session a été déplacée au cours des 2 dernières minutes passe après les autres. Le choix doit se faire en 6 secondes environ ; quand les consultations d'utilisation prennent plus longtemps, la session reste. Avec `usage_source = "claude"`, `auto_switch` décide à partir des seuls chiffres enregistrés, et ccam ne peut pas savoir, ni par la consultation de Claude Code ni par la ligne d'état, si les crédits d'utilisation sont activés : il suppose donc qu'ils sont désactivés.

Le hook termine alors Claude Code, et ccam, qui tourne à côté, affiche une notification de bureau dans les terminaux qui le permettent, par exemple `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` Il affiche la première ligne de l'exemple de sortie ci-dessus et exécute `ccam next --continue` pour cette session avec le compte choisi, qui affiche ses deux autres lignes et lance Claude Code avec les options, le mode de permission, le niveau d'effort et l'agent de la session, comme le fait le [relais](#relay). Le premier message de la nouvelle session est un texte fixe, la seule chose que ccam envoie en votre nom (un changement après `/ccam:switch` n'envoie rien) :

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

Après la limite d'un sous-agent, c'est celui-ci, car les sous-agents, workflows et commandes en arrière-plan se sont arrêtés avec l'ancienne session :

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

Si ce compte atteint à son tour sa limite, la même chose se produit sur celui-ci. Les relevés de limite empêchent la chaîne de revenir à un compte avant sa réinitialisation : quand tous les comptes sont pleins, la session reste donc sur le dernier.

Le hook apprend qu'une conversation s'exécute avec un agent par le paramètre `agent`, ou par la ligne nommant l'agent que Claude Code écrit dans la conversation de temps en temps et à la fin de la session. Si jusque-là seul le `--agent` de la session le nomme, et que le compte choisi ne partage pas les agents, `ccam next --continue` ne lance rien et garde la note, et une deuxième notification de bureau le signale : ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

Quand aucun autre compte n'a de marge, rien ne bouge : Claude Code gère la limite comme d'habitude, en attendant dans la session ou en proposant ses options, et la notification indique quel compte se libère en premier, par exemple `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` Pour des sièges, elle indique `None of your other accounts in Acme has room now`. Pour arrêter un changement déjà commencé, appuyez sur Ctrl+C pendant que Claude Code se ferme : ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

Un changement a des coûts. Claude Code redémarre, ce qui prend quelques secondes, et les shells en arrière-plan, serveurs de développement et agents lancés par Claude s'arrêtent avec l'ancienne session. Le nouveau compte relit toute la conversation avec sa première réponse, comme avec [`--continue`](#continue), et ccam en indique la taille ; un changement n'utilise jamais [`--brief`](#brief).

`auto_switch` ne change jamais de compte avant que Claude Code ne signale une limite, ne surveille jamais votre utilisation en arrière-plan, et ne copie, ne renouvelle ni ne transmet jamais une connexion. Sous Windows, il est ignoré, et `ccam doctor` le signale ; avec `relay = true`, `/ccam:switch` et `/exit` y déplacent toujours une conversation vers un autre compte.

**Conditions d'Anthropic.** Sur des sièges Team et Enterprise, `auto_switch` fait passer votre travail sur un autre siège de la même organisation. Le centre d'aide d'Anthropic indique : « Each seat on a Team plan is assigned to one named individual, and only that person may use it » (chaque siège d'une offre Team est attribué à une personne nommément désignée, et seule cette personne peut l'utiliser) et « Users can't route their work through another user's account to use that user's remaining usage » (les utilisateurs ne peuvent pas faire passer leur travail par le compte d'un autre utilisateur pour profiter de l'utilisation restante de celui-ci) ([offres Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129) ; de même pour les [offres Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)).

Pour les offres personnelles (Pro, Max), les conditions d'Anthropic ne disent pas si vos propres comptes peuvent prendre automatiquement le relais les uns des autres à une limite. Les [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) interdisent d'accéder aux Services « through automated or non-human means, whether through a bot, script, or otherwise » (par des moyens automatisés ou non humains, que ce soit au moyen d'un bot, d'un script ou autrement) sauf autorisation explicite d'Anthropic, et un changement envoie en votre nom l'un des messages fixes ci-dessus. La [Usage Policy](https://www.anthropic.com/legal/aup) qui entre en vigueur le 12 novembre 2026 compte « Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products » (contourner ou perturber intentionnellement les capacités, restrictions ou garde-fous établis dans nos produits) parmi ce qu'elle interdit.

ccam n'affirme pas que `auto_switch` respecte ces conditions. L'activer relève de votre décision, tout comme le risque : Anthropic peut suspendre un compte, ou l'accès d'une organisation Team ou Enterprise, lorsqu'elle estime que les conditions ne sont pas respectées.

Tant que `auto_switch` est activé, `ccam doctor` signale un compte sans le plugin ccam et les versions de Claude Code antérieures à 2.1.214, et, pour les sièges, rappelle la règle d'Anthropic sur les sièges ([vérifications](/features/doctor)).

## Relancer à la limite

`ccam next` sert à une session interactive qui s'est arrêtée. Une exécution non interactive, `claude -p` dans un script ou une tâche planifiée, ne peut pas vous attendre : `--fallback` avant le nom du compte laisse donc ccam changer de compte pour elle :

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

1. `--fallback` ne fonctionne qu'avec `-p` ou `--print` ; sinon ccam s'arrête et renvoie vers `ccam next --auto --continue`.
2. L'entrée standard reçue par un tube est lue une seule fois, jusqu'à 10 Mio, et donnée à chaque tentative. Comme Claude Code, ccam attend au plus 3 secondes qu'elle commence (sauf avec `--input-format stream-json`) ; si un tube reste ouvert sans rien envoyer, l'exécution ne reçoit aucune entrée et ccam le signale sur la sortie d'erreur standard.
3. ccam démarre `claude` comme processus enfant au lieu d'`exec`, retient sa sortie standard et laisse passer la sortie d'erreur standard au fur et à mesure.
4. Quand `claude` se termine avec 0, ou échoue pour une autre raison, la sortie retenue et le code de sortie sont transmis sans modification. Il en va de même pour une exécution terminée par un signal (un code de sortie supérieur à 128) : on lui a demandé de s'arrêter, et elle n'est pas relancée.
5. Un échec compte comme la limite d'utilisation quand le message de Claude Code le dit, avec l'un des messages de limite d'utilisation que connaît le [hook](#inside-claude-code), y compris dans une sortie `--output-format json` ou `stream-json`, ou quand l'utilisation du compte montre sa fenêtre de 5 heures, hebdomadaire ou `--model` à 100 % (relevé de limite compris).
6. La sortie de cette tentative est alors abandonnée et la même commande s'exécute avec un autre compte connecté : celui que `--auto` choisirait, puis les comptes dont l'utilisation ne peut pas être consultée, dans l'ordre des numéros. Les comptes à 100 % sont ignorés, les comptes connectés par clé d'API ne sont jamais essayés, et chaque compte est essayé au plus une fois.
7. Avec `-c`, `--continue`, `-r` ou `--resume` parmi les arguments, seuls les comptes du même groupe d'historique sont essayés, afin que la conversation puisse être trouvée.
8. Seuls les comptes de la même organisation sont essayés, comme avec `--auto` (`cross_org = true` lève cette règle).

Le compte suivant est choisi selon les mêmes règles que `--auto` (`auto_max_usage`, `auto_strategy`, la fenêtre `--model`). Quand il ne reste aucun compte, ccam affiche la sortie de la dernière tentative, indique qu'aucun compte connecté n'a plus de marge et quitte avec le code 75. Avec `--wait` avant le nom du compte (`ccam --fallback --wait=2h work -p …`), il attend plutôt la réinitialisation du premier compte, puis peut réessayer chaque compte ; il attend toujours au moins une fois d'abord, afin qu'un compte qui vient d'atteindre sa limite ne soit pas relancé aussitôt. La durée limite le temps d'attente total de l'exécution, pas chaque attente. Écrite avec une espace (`--wait 30m`), elle est refusée avec la même indication que pour `ccam next`. Un `--` initial avant les arguments de Claude Code fonctionne comme pour tout lancement (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` affiche la commande de la première tentative.

Les sessions interactives ne changent pas de compte d'elles-mêmes, sauf si vous activez [`auto_switch`](#auto-switch). Avec `ccam statusline install`, la ligne d'état ajoute `→ ccam next --auto --continue` dès que l'utilisation sur 5 heures ou hebdomadaire atteint 90 % ([Limites d'utilisation](/features/usage)).

## Sécurité et cas particuliers

- Les comptes sans connexion ne sont jamais choisis. Sans autre compte connecté, `ccam next` seul s'arrête avec `no logged-in account other than <name>`.
- Avec `--continue` et aucune conversation dans le dossier : `no account has a conversation in <folder> to continue`.
- Quand la règle d'organisation ne laisse aucun compte : `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Avec `--auto` et aucun compte connecté : `no logged-in account to pick`. Quand chaque compte est plein ou que son utilisation est inconnue : `` no account with known usage and headroom (see `ccam usage`) ``. Les deux se terminent avec le code de sortie 75.
- `--max-usage`, `--strategy`, `--wait` et `--cross-org` ne fonctionnent qu'avec `--auto`, et une valeur hors de 1-100 ou une stratégie inconnue est refusée.
- Une bifurcation copie la conversation, pas le reste de l'état du compte : les points de contrôle du `file-history/` du compte d'origine ne suivent pas, sauf si l'historique est partagé.
- `--auto` fait des requêtes réseau même avec `--dry-run`. Il ne fait que lire et n'écrit jamais dans les fichiers de jetons.
- Le compte choisi démarre par le même chemin que `ccam <name>`, si bien que les mêmes vérifications s'appliquent, y compris l'avertissement quand une autre session Claude Code qui partage le même historique s'exécute déjà dans ce dossier.
- `last-used` n'est écrit que lorsque Claude Code démarre réellement. `--dry-run` n'y touche pas, si bien qu'une simulation ne fait pas avancer la rotation.
- `--fallback` retient la sortie standard jusqu'à la fin d'une tentative : une sortie en flux comme `--output-format stream-json` arrive donc d'un seul bloc à la fin. Un client qui envoie une entrée `stream-json` en flux et attend des réponses avant d'en envoyer davantage ne peut pas utiliser `--fallback`. Les fichiers que la première tentative a modifiés avant d'atteindre la limite restent modifiés ; la tentative suivante reprend l'invite depuis le début.

## Voir aussi

- [Limites d'utilisation](/features/usage)
- [Limites d'utilisation et changement de compte](/docs/usage-limits)
- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [Que faire quand Claude Code atteint sa limite](/guides/usage-limit-workflow)
