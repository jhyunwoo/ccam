---
description: Continue this conversation on another ccam account (the next one, or the one you name)
argument-hint: "[account]"
disable-model-invocation: true
---

The ccam plugin's hook handles this command before it reaches you: it notes this conversation for `ccam next --continue` and ends the turn. If you are reading this, the hook did not run (Claude Code before 2.1.78, the ccam plugin's hooks turned off, or `ccam` not on PATH).

Tell the user, briefly: type `/exit`, then run `ccam next --continue` in this folder (add `--auto` to take the account with the most room). ccam then opens this folder's last conversation as another account; Claude Code copies it over when that account does not share the history.

If the user named an account (see the arguments), say that without the hook `ccam next --continue` cannot know that name: it takes the next account in order, or with `--auto` the one with the most room. To go to the named account, they can fix what kept the hook from running and type `/ccam:switch <account>` again.
