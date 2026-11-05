---
title: "Limites d'utilisation et changement de compte"
description: "Consultez l'utilisation sur 5 heures et hebdomadaire de chaque compte Claude Code avec ccam usage, et ouvrez celui qui a le plus de marge avec ccam next --auto."
updated: 2026-11-05
order: 5
translation: machine
sourceHash: "06bd349073e3"
---

`ccam usage` indique, pour chaque compte, quelle part de chaque limite est utilisée : la limite sur 5 heures, la limite hebdomadaire et la limite hebdomadaire propre à chaque modèle (comme Fable). `ccam next --auto` ouvre le compte qui a le plus de marge, et `--continue` reprend la conversation en cours.

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

Chaque compte a une ligne qui indique l'ancienneté de ses chiffres, puis une barre par limite avec le temps restant avant sa réinitialisation : la fenêtre de 5 heures, la fenêtre hebdomadaire de tous les modèles et la fenêtre hebdomadaire propre à chaque modèle, comme Fable. Une barre passe au jaune à 80 % et au rouge à 100 %. La dernière ligne nomme le compte que `ccam next --auto` ouvrirait ; avec `auto_strategy = "expiring"` ou `"balance"`, elle indique à la place `→ next --auto picks <name>: <reason>`.

**Sièges Team et Enterprise.** La sélection automatique de compte (`ccam next`, `--fallback`, `ccam pick`, la touche `n` du tableau de bord et, une fois que vous l'activez, [`auto_switch`](/features/next#auto-switch)) passe aussi d'un siège à l'autre au sein d'une même organisation Claude Team ou Enterprise, ce qui peut enfreindre les conditions d'Anthropic. Lisez la [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) et renseignez-vous auprès de votre organisation et d'Anthropic avant de vous y fier.

## D'où viennent les chiffres

ccam interroge le point de terminaison d'utilisation d'Anthropic pour tous les comptes en même temps, avec la propre connexion de chaque compte. Il attend au plus 3 secondes par compte et affiche chaque compte dès qu'il répond.

Ce point de terminaison n'est pas documenté publiquement et peut changer. Quand une requête échoue, ccam affiche les derniers chiffres connus et les signale comme provenant du cache. Une fenêtre réinitialisée depuis l'enregistrement de ces chiffres affiche 0 %. `ccam doctor` vous prévient si le format de la réponse a changé.

Quand Claude Code arrête une session à une limite d'utilisation, le [plugin ccam](/docs/install#claude-code-plugin) enregistre quelle fenêtre était pleine et quand elle se réinitialise. Jusque-là, cette fenêtre affiche 100 % quels que soient les autres chiffres : `ccam next --auto` et `--fallback` ne choisissent donc pas ce compte de nouveau trop tôt.

Les chiffres récupérés il y a moins de 60 secondes sont réutilisés. Options :

| Option | Effet |
| --- | --- |
| `--watch[=seconds]` | Actualiser en continu, toutes les 60 secondes par défaut |
| `--json` | Sortie lisible par une machine |
| `-v` | Afficher aussi la provenance des chiffres de chaque compte |
| `--no-cache` | Toujours redemander |
| `--source cache` | Ne jamais demander ; afficher uniquement les chiffres enregistrés |

## Connexions expirées

Quand le jeton d'accès d'un compte a expiré, ccam ne le renouvelle pas, car le renouveler pendant qu'une session Claude Code utilise la même connexion pourrait déconnecter cette session. Le compte affiche ``access token expired; `ccam <name>` renews it when it starts`` : lancez-le une fois avec `ccam <name>` et Claude Code le renouvellera. Une connexion entièrement expirée, ou que Claude Code a effacée, affiche `` login expired → `ccam login <name>` `` et demande une nouvelle connexion.

`ccam usage --refresh` renouvelle lui-même les jetons expirés, mais uniquement pour les comptes sans session Claude Code en cours, quelle que soit la façon dont cette session a été lancée. Il ne renouvelle jamais les connexions stockées dans le trousseau macOS ni dans un `~/.claude` importé, et un renouvellement peut prendre jusqu'à 30 secondes.

## L'utilisation dans la ligne d'état

```sh
ccam statusline install
```

Claude Code affiche alors le nom du compte et son utilisation dans la ligne d'état, par exemple `[work] · Opus · 5h 24% · 7d 41%`. À partir de 90 %, la ligne suggère aussi `ccam next --auto --continue`. Si vous avez déjà une commande de ligne d'état, ccam la conserve et place le nom du compte devant. La ligne d'état enregistre aussi ses chiffres, que `ccam usage` utilise quand le point de terminaison ne répond pas.

## Changer de compte

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the most room
ccam next --continue  # continue this folder's conversation on the next account
```

`ccam next` ignore les comptes qui ne sont pas connectés, et `--auto` ignore les comptes à 100 % (y compris une fenêtre que Claude Code a signalée pleine, jusqu'à sa réinitialisation) et ceux d'une autre organisation. Quand les deux comptes partagent l'historique des conversations, `--continue` ouvre la même conversation ; sinon, Claude Code la copie dans le nouveau compte. Voir [Partage](/docs/sharing) et [Passer au compte suivant](/features/next#continue).

Avec `auto_switch = true` en haut de `config.toml` (Linux et macOS ; désactivé par défaut), ccam effectue lui-même ce changement quand Claude Code arrête une session interactive à sa limite et qu'un autre compte du même type (un siège de la même organisation, ou une autre offre personnelle) a de la marge. Voir [`auto_switch`](/features/next#auto-switch) pour savoir quand il agit et ce qu'il ne fait jamais.

Pour les scripts et les tâches planifiées qui exécutent `claude -p`, placez `--fallback` avant le nom du compte. Quand ce compte atteint sa limite, ccam relance la même commande avec le compte qui a le plus de marge, et seule la sortie de l'exécution qui a abouti est affichée :

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Voir [Relancer à la limite](/features/next#run-again-at-the-limit) pour savoir comment ccam décide.

Le [guide sur la limite d'utilisation](/guides/usage-limit-workflow) présente un exemple complet.
