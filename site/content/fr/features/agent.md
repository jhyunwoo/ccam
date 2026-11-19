---
title: "Des sous-agents sur d'autres comptes"
description: "ccam agent confie une partie d'une tâche Claude Code à un autre de vos comptes et renvoie la réponse de Claude : une longue tâche puise dans plusieurs comptes."
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Les sous-agents propres à Claude Code s'exécutent avec le compte de la session qui les lance, et consomment donc l'utilisation de ce compte. `ccam agent "<task>"` exécute une tâche avec `claude -p` dans le même dossier avec un autre de vos comptes et affiche la dernière réponse de Claude. Appelé depuis une session Claude Code, il fonctionne comme un sous-agent qui consomme l'utilisation d'un autre compte, si bien qu'une longue tâche peut continuer au-delà du moment où le compte de la session aurait été épuisé.

**Sièges Team et Enterprise.** `ccam agent` choisit le compte de lui-même. Au sein d'une même organisation Claude Team ou Enterprise, il exécute la tâche sur un autre siège, ce qui peut enfreindre les conditions d'Anthropic. Lisez la [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) et renseignez-vous auprès de votre organisation et d'Anthropic avant de vous y fier.

## Utilisation

Claude Code le lance avec son outil Bash. Avec le [plugin ccam](/docs/install#claude-code-plugin) installé, demandez-le simplement avec vos mots, par exemple « utilise mes autres comptes pour certaines parties de cette tâche », ou écrivez-le dans votre CLAUDE.md. Le skill du plugin indique à Claude Code de lancer chaque agent en arrière-plan, d'écrire une tâche qui se suffit à elle-même et de donner aux agents qui modifient des fichiers leur propre partie ou leur propre worktree git.

```sh
ccam agent "find why TestParse fails and fix it; reply with what you changed"
ccam agent --model opus -- "-v breaks the parser: find out why"
ccam agent --json < task.md            # the task from standard input
ccam agent --resume 7f3c9a2e-… "now add a test for it"
ccam agent --wait=1h "…"               # nobody has room: wait for a reset
ccam --dry-run agent "…"               # the account and command it would run
```

```console
$ ccam agent "summarize what internal/usage does"
ccam: agent runs as work (12% of its 5 hours used); session 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f
internal/usage looks up the 5-hour and weekly usage of every account …
ccam: agent finished as work; follow up with: ccam agent --resume 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f "…"
```

La dernière réponse de Claude va sur la sortie standard ; les lignes qui commencent par `ccam:` vont sur la sortie d'erreur standard. `--json` affiche à la place `{account, session_id, result, is_error, accounts}`, où `accounts` liste tous les comptes avec lesquels la tâche s'est exécutée.

Codes de sortie : celui de claude (`0` quand l'agent a terminé), `75` (`EX_TEMPFAIL`) quand aucun compte autre que l'appelant n'a de marge, `1` pour les autres erreurs, comme un agent qui tente d'en lancer un autre.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--model <model>` | Le modèle de l'agent ; la fenêtre hebdomadaire de ce modèle compte aussi dans le choix | le modèle par défaut du compte |
| `--resume <id>` | Poursuivre la conversation d'un agent précédent avec une tâche de suivi (l'ID de session figure sur sa première et sa dernière ligne) | désactivé |
| `--json` | Afficher le compte, la session, la réponse et les comptes avec lesquels la tâche s'est exécutée | désactivé |
| `--max-usage <1-100>` | Préférer les comptes sous ce niveau d'utilisation | `auto_max_usage`, sinon 100 |
| `--strategy headroom\|expiring\|balance` | La [règle de choix](/features/next#auto) | `auto_strategy`, sinon `headroom` |
| `--cross-org` | S'exécuter aussi avec des comptes d'autres organisations | `cross_org`, sinon désactivé |
| `--wait[=<duration>]` | Quand aucun compte n'a de marge, attendre la réinitialisation du premier, au plus cette durée | désactivé ; `5h` s'il est donné seul |
| `--dry-run` | Afficher la commande pour le compte qu'il choisirait, sans l'exécuter | désactivé |

## Fonctionnement

1. **L'appelant est laissé de côté.** Le compte depuis lequel `ccam agent` est appelé (celui de la session que ccam a lancée, sinon celui que lance ce dossier) n'exécute jamais l'agent, pas plus qu'un autre dossier connecté au même compte Claude, qui en partage les limites. L'agent reste dans l'[organisation](/features/next#auto) de l'appelant, sauf si `--cross-org` en décide autrement.
2. **Il choisit comme `ccam pick`.** L'utilisation est consultée comme le fait [`ccam pick`](/features/pick), et les agents lancés à moins de 2 minutes d'intervalle reçoivent des comptes différents tant que suffisamment de comptes ont de la marge.
3. **Il s'exécute avec les permissions de la session, jamais plus.** Personne n'est là pour autoriser un outil dans `claude -p` : l'agent reçoit donc le mode de permission et le niveau d'effort de la session qui l'a appelé. Lancé depuis un terminal, il suit le mode par défaut de Claude Code et les règles de permission de vos paramètres. `ccam agent` n'a aucune option qui accorde davantage.
4. **À une limite d'utilisation, la conversation continue.** Quand le compte de l'agent atteint sa limite, ccam enregistre la limite comme le ferait le [hook du plugin](/features/usage#the-limit-record) et poursuit la même conversation avec le compte qui a le plus de marge : Claude Code l'y copie (`--resume <file> --fork-session`), avec comme premier message « Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete. » Chaque compte n'est essayé qu'une fois.
5. **Chaque conversation d'agent est enregistrée** dans `~/.local/state/ccam/agent/<session id>.json` (30 jours). `--resume <id>` la poursuit dans son dossier, avec le compte qui l'a exécutée s'il a encore de la marge, sinon avec un autre. [`ccam next --continue`](/features/next) laisse de côté les conversations d'agents, si bien qu'il continue de poursuivre les vôtres.

## Sécurité et cas particuliers

- Un agent ne peut pas lancer ses propres agents : ccam le marque avec `CCAM_AGENT=1`, et `ccam agent` refuse de s'y exécuter.
- Arrêter `ccam agent` (Ctrl+C, un signal, ou Claude Code qui interrompt la commande) arrête aussi l'agent.
- L'agent ne voit pas la conversation qui l'a appelé, seulement sa tâche.
- Des agents qui modifient les mêmes fichiers en même temps écrasent le travail les uns des autres. Donnez à chacun sa propre partie de la tâche, ou son propre worktree git.
- Quand la session appelante passe à un autre compte ([`auto_switch`](/features/next#auto-switch)), les agents qu'elle a lancés s'arrêtent avec elle. Leurs conversations restent : `--resume <id>` les poursuit.
- `agent` est un nom de commande : aucun nouveau compte ne peut donc s'appeler `agent` ; un compte existant de ce nom se lance toujours avec `ccam run agent`.

## Voir aussi

- [Choisir un compte pour les scripts](/features/pick)
- [Passer au compte suivant](/features/next)
- [Limites d'utilisation](/features/usage)
- [Les sessions en arrière-plan de tous les comptes](/features/agents)
