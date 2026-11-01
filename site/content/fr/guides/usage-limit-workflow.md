---
title: "Que faire quand Claude Code atteint sa limite"
description: "Attendez la réinitialisation ou poursuivez la même conversation sur un autre compte Claude Code. Les limites sur 5 heures et hebdomadaires, étape par étape."
updated: 2026-11-01
order: 3
translation: machine
sourceHash: 760a36724feb
---

Quand un compte atteint sa limite d'utilisation (souvent appelée « rate limit »), vous avez deux options : attendre que la limite se réinitialise, ou poursuivre dès maintenant la même conversation sur un autre compte qui vous appartient. Avec ccam, la seconde ne demande qu'une commande, `ccam next --auto --continue`.

## Quelles sont les limites de Claude Code ?

Avec les forfaits Pro et Max, Claude Code a une [limite de session sur 5 heures et une limite hebdomadaire](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Les deux sont partagées avec Claude sur le web et dans les applications, et l'utilisation compte pour les deux à la fois. Quand l'une est épuisée, Claude Code [affiche un message](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) comme « You've hit your session limit · resets 3:45pm ». Changer de modèle n'aide pas, car les limites couvrent tous les modèles.

## Option 1 : attendre la réinitialisation

Si vous n'avez qu'un compte, attendez. Les versions récentes de Claude Code [attendent dans la session ouverte](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) et reprennent la tâche d'elles-mêmes une fois la limite réinitialisée. `/usage` dans Claude Code indique quelle part des limites de votre forfait vous avez utilisée.

## Option 2 : poursuivre tout de suite sur un autre compte

Si vous avez un second compte, vous pouvez continuer immédiatement la même conversation. Configurez-le une fois :

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

La ligne d'état affiche désormais le compte et son utilisation sur 5 heures pendant que vous travaillez : la limite ne vous prend donc pas par surprise.

Quand la limite est atteinte, quittez Claude Code (`/exit`), puis :

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Ouvrez la conversation dans le compte qui a le plus de marge :

```sh
ccam next --auto --continue
```

Claude Code démarre en tant que `work` avec la conversation la plus récente du dossier courant, y compris tout ce qui a été dit jusque-là. Partager l'historique (`ccam group set default --add history` ci-dessus) fait que les deux comptes travaillent sur la même conversation ; sans cela, Claude Code copie la conversation dans `work`, et ccam vous indique combien de jetons la première réponse doit relire.

Avec le [plugin ccam](/docs/install#claude-code-plugin) dans Claude Code, une session qui s'arrête à sa limite affiche une notification et enregistre la fenêtre pleine, si bien que `ccam next --auto` ne choisit pas ce compte de nouveau avant sa réinitialisation. `/ccam:switch` (ou `/ccam:switch work`) note la conversation pour `ccam next --continue`, qui poursuit donc précisément cette conversation.

**Sièges Team et Enterprise.** La sélection automatique de compte (`ccam next`, `--fallback`, `ccam pick`, la touche `n` du tableau de bord et, une fois que vous l'activez, [`auto_switch`](/features/next#auto-switch)) passe aussi d'un siège à l'autre au sein d'une même organisation Claude Team ou Enterprise, ce qui peut enfreindre les conditions d'Anthropic. Lisez la [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) et renseignez-vous auprès de votre organisation et d'Anthropic avant de vous y fier.

## ccam change-t-il de compte automatiquement ?

Seulement si vous l'activez.

- Par défaut, ccam ne tourne pas pendant que vous travaillez : il lance Claude Code et se remplace par lui, si bien que c'est vous qui décidez quand changer. `ccam next --auto` choisit le compte dont l'utilisation sur 5 heures est la plus faible, ignore les comptes à 100 % ou non connectés, et reste au sein d'une même organisation Team ou Enterprise, ou parmi les offres personnelles.
- Avec `relay = true` dans `config.toml`, ccam continue plutôt de tourner à côté de la session, si bien que `/ccam:switch` poursuit la conversation de lui-même, en conservant les options, le mode de permission, le niveau d'effort et l'agent de la session ; avant de déplacer la conversation vers un compte d'une autre organisation, ccam vous le demande.
- Avec `auto_switch = true` (Linux et macOS ; expérimental), ccam effectue aussi le changement lui-même : quand Claude Code arrête la session à sa limite et qu'un autre compte du même type a de la marge, ccam met fin à la session et poursuit la conversation sur ce compte ([conditions et coûts](/features/next#auto-switch)).

Certains autres outils échangent plutôt la connexion d'une session en cours ; [ccam et les autres sélecteurs de compte](/compare/account-switchers) explique la différence.

## Points à retenir

- Exécutez `ccam next --continue` dans le dossier où vous travailliez ; les conversations sont stockées par dossier de projet.
- Évitez d'utiliser deux comptes sur la même conversation en même temps. ccam vous avertit si vous essayez.
- Sans historique partagé, la conversation est copiée dans l'autre compte et les deux copies évoluent chacune de leur côté ; les points de contrôle à partir desquels `/rewind` restaure le code restent dans le premier compte.
- Le premier message sur le nouveau compte peut consommer plus de sa limite que d'habitude, car ce compte n'a pas encore la conversation en cache. `ccam next --continue --brief` le démarre plutôt sur un court résumé.
- Le droit de passer ainsi d'un compte à l'autre dépend des conditions qui s'appliquent à chaque compte.

Voir [Limites d'utilisation](/features/usage) et [Passer au compte suivant](/features/next) pour toutes les options.
