---
title: ಕಮಾಂಡ್ ಉಲ್ಲೇಖ
metaTitle: ccam ಕಮಾಂಡ್ ಉಲ್ಲೇಖ
description: ccam add ಮತ್ತು ccam <name> ನಿಂದ ccam usage, ccam next ಮತ್ತು ccam doctor ವರೆಗೆ, ccam ನ ಎಲ್ಲಾ ಕಮಾಂಡ್‌ಗಳು ಮತ್ತು ಆಯ್ಕೆಗಳು.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: "28e1ad1dcbad"
---

ಕಮಾಂಡ್ ಏನು ಮಾಡುತ್ತದೆ ಎಂಬುದನ್ನು ನಿಜವಾಗಿ ಮಾಡದೆಯೇ ನೋಡಲು, ಅದರ ಮುಂದೆ `--dry-run` ಹಾಕಿ. ಉದಾಹರಣೆಗೆ `ccam --dry-run work`. ಇದು Claude Code ಆರಂಭಿಸುವಾಗ (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) ಮತ್ತು `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` ಹಾಗೂ `doctor --fix` ಜೊತೆ ಕೆಲಸ ಮಾಡುತ್ತದೆ. `list`, `usage`, `share` ಮತ್ತು `ext` ನಂತಹ ಓದುವ ಕಮಾಂಡ್‌ಗಳು ಇದನ್ನು ಸ್ವೀಕರಿಸುತ್ತವೆ, ಹೇಗಿದ್ದರೂ ಅವು ಏನನ್ನೂ ಬದಲಾಯಿಸುವುದಿಲ್ಲ. ಉಳಿದವು, `group rename` ಸೇರಿದಂತೆ, ಚಾಲನೆಯಾಗುವ ಬದಲು ಇದನ್ನು ನಿರಾಕರಿಸುತ್ತವೆ.

<!-- commands -->

## Claude Code ಆರಂಭಿಸುವುದು

`ccam <name|number> [arguments]` ಆ ಖಾತೆಯಾಗಿ Claude Code ಅನ್ನು ಆರಂಭಿಸುತ್ತದೆ. ಹೆಸರಿನ ನಂತರದ ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳನ್ನು ccam ಓದುವುದಿಲ್ಲ; ಅದು ಅವುಗಳನ್ನು ಅದೇ ಕ್ರಮದಲ್ಲಿ `claude` ಗೆ ರವಾನಿಸುತ್ತದೆ. ಸ್ಟ್ಯಾಂಡರ್ಡ್ ಇನ್‌ಪುಟ್, ಎಕ್ಸಿಟ್ ಕೋಡ್ ಮತ್ತು ಸಿಗ್ನಲ್‌ಗಳು `claude` ನಲ್ಲಿ ಇರುವಂತೆಯೇ ವರ್ತಿಸುತ್ತವೆ. ಹೆಸರಿನ ನಂತರ `--` ಐಚ್ಛಿಕ.

ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳಿದ್ದು ಖಾತೆಯ ಹೆಸರು ಇಲ್ಲದ `ccam`, ಉದಾಹರಣೆಗೆ `ccam -p "summarise"`, ಈಗಿನ ಫೋಲ್ಡರ್‌ಗೆ ಹೊಂದಿಸಿದ ಖಾತೆಯನ್ನು, ಇಲ್ಲದಿದ್ದರೆ ಡೀಫಾಲ್ಟ್ ಖಾತೆಯನ್ನು ಬಳಸುತ್ತದೆ. ಅವುಗಳನ್ನು `ccam folder <name> [folder]` ಮತ್ತು `ccam use <name>` ಮೂಲಕ ಹೊಂದಿಸಿ.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

ಖಾತೆಯ ಹೆಸರಿನ ಮುಂದೆ ಇರುವ `--fallback`, `-p`/`--print` ರನ್‌ಗಳಿಗಾಗಿ. ಖಾತೆ ತನ್ನ ಬಳಕೆಯ ಮಿತಿ ತಲುಪಿದಾಗ, ccam ಅದೇ ಕಮಾಂಡ್ ಅನ್ನು ಹೆಚ್ಚು ಅವಕಾಶ ಇರುವ ಲಾಗಿನ್ ಆದ ಖಾತೆಯಾಗಿ ಮತ್ತೆ ರನ್ ಮಾಡುತ್ತದೆ ಮತ್ತು ಪೂರ್ಣಗೊಂಡ ರನ್‌ನ ಔಟ್‌ಪುಟ್ ಅನ್ನು ಮಾತ್ರ ಪ್ರಿಂಟ್ ಮಾಡುತ್ತದೆ. [ಫೋಲ್ಡರ್‌ವಾರು ಖಾತೆಗಳು](/features/folder) ಮತ್ತು [ಮಿತಿಯಲ್ಲಿ ಮತ್ತೆ ರನ್ ಮಾಡುವುದು](/features/next#run-again-at-the-limit) ನೋಡಿ.

**Team ಮತ್ತು Enterprise ಸೀಟುಗಳು.** ಸ್ವಯಂಚಾಲಿತ ಖಾತೆ ಆಯ್ಕೆ (`ccam next`, `--fallback`, `ccam pick`, ಡ್ಯಾಶ್‌ಬೋರ್ಡ್‌ನ `n` ಮತ್ತು, ನೀವು ಆನ್ ಮಾಡಿದ ನಂತರ, [`auto_switch`](/features/next#auto-switch)) ಒಂದೇ Claude Team ಅಥವಾ Enterprise ಸಂಸ್ಥೆಯ ಸೀಟುಗಳ ನಡುವೆಯೂ ಬದಲಾಯಿಸುತ್ತದೆ, ಇದು Anthropic ನ ನಿಯಮಗಳನ್ನು ಉಲ್ಲಂಘಿಸಬಹುದು. ಇದನ್ನು ಅವಲಂಬಿಸುವ ಮೊದಲು [ಪದೇ ಪದೇ ಕೇಳುವ ಪ್ರಶ್ನೆಗಳನ್ನು](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) ಓದಿ ಮತ್ತು ನಿಮ್ಮ ಸಂಸ್ಥೆ ಹಾಗೂ Anthropic ಬಳಿ ವಿಚಾರಿಸಿ.

`ccwork` ನಂತಹ ಶಾರ್ಟ್‌ಕಟ್ ಕೂಡ ತನ್ನ ಖಾತೆಯನ್ನು ಇದೇ ರೀತಿ ಆರಂಭಿಸುತ್ತದೆ: `ccam work --resume` ಮಾಡುವುದನ್ನೇ `ccwork --resume` ಮಾಡುತ್ತದೆ. ಕೆಳಗಿನ [ಶಾರ್ಟ್‌ಕಟ್‌ಗಳು](/docs/commands#shortcuts) ನೋಡಿ.

## ಖಾತೆಗಳನ್ನು ಸೇರಿಸುವುದು ಮತ್ತು ತೆಗೆದುಹಾಕುವುದು

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

ಫೋಲ್ಡರ್ ನೀಡದೆ `ccam import` ರನ್ ಮಾಡಿದರೆ, ಅದು `~/.claude` ಮತ್ತು `~/.claude-*` ಗಾಗಿ ಹುಡುಕುತ್ತದೆ ಮತ್ತು ಪ್ರತಿಯೊಂದರ ಬಗ್ಗೆ ಕೇಳುತ್ತದೆ; ವಿಫಲವಾದ ಫೋಲ್ಡರ್ ಅನ್ನು `Skipped <path>: <reason>` ಎಂದು ತಿಳಿಸಲಾಗುತ್ತದೆ ಮತ್ತು ಉಳಿದವು ನೋಂದಣಿಯಾಗುತ್ತವೆ. ಫೋಲ್ಡರ್‌ನ ಹೆಸರು ಈಗಾಗಲೇ ಬಳಕೆಯಲ್ಲಿದ್ದರೆ, ಆ ಫೋಲ್ಡರ್ ಅನ್ನು `ccam import <folder> --name <name>` ಮೂಲಕ ಪ್ರತ್ಯೇಕವಾಗಿ ಇಂಪೋರ್ಟ್ ಮಾಡಿ. ಇಂಪೋರ್ಟ್ ಮಾಡಿದ ಫೋಲ್ಡರ್‌ಗಳು ಇದ್ದಲ್ಲೇ ನೋಂದಣಿಯಾಗುತ್ತವೆ; ccam ಎಂದಿಗೂ ಲಾಗಿನ್ ಅನ್ನು ಸ್ಥಳಾಂತರಿಸುವುದಿಲ್ಲ ಅಥವಾ ನಕಲಿಸುವುದಿಲ್ಲ.

Claude Code ಲಾಗಿನ್ ಅನ್ನು ತೆರವುಗೊಳಿಸಿದ ಖಾತೆಗಳಿಗೂ `ccam login --expired` ಲಾಗಿನ್ ಮಾಡುತ್ತದೆ. ನಿರ್ವಹಿಸಲಾದ ಖಾತೆಯ (`ccam import` ಮೂಲಕ ನೋಂದಾಯಿಸಿದ ಖಾತೆಯಲ್ಲ, ccam ತಾನೇ ಫೋಲ್ಡರ್ ರಚಿಸಿದ ಖಾತೆ) ಫೋಲ್ಡರ್ ಕಸದ ಫೋಲ್ಡರ್‌ಗೆ ಹೋಗುವ ಮೊದಲು `ccam remove` ಆ ಖಾತೆಯನ್ನು ಲಾಗ್ ಔಟ್ ಮಾಡುತ್ತದೆ, ಇನ್ನೊಂದು ಖಾತೆ ಅದೇ ಲಾಗಿನ್ ಹೊಂದಿದ್ದರೆ ಹೊರತು. ಇಂಪೋರ್ಟ್ ಮಾಡಿದ ಖಾತೆಯನ್ನು ತೆಗೆದುಹಾಕಿದರೆ ಅದು ತನ್ನ ಗುಂಪುಗಳಿಂದ ಹೊರಬರುತ್ತದೆ, ಮತ್ತು ಹಂಚಿಕೊಂಡ ಪ್ರತಿ ಐಟಂನ ಒಂದು ಪ್ರತಿ ಅದರ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ ಉಳಿಯುತ್ತದೆ. macOS ನಲ್ಲಿ, `ccam rename` ಫೋಲ್ಡರ್‌ನ ಹೆಸರನ್ನು ಹಾಗೆಯೇ ಇಡುತ್ತದೆ, ಏಕೆಂದರೆ Keychain ಲಾಗಿನ್ ಅನ್ನು ಫೋಲ್ಡರ್‌ನ ಪಾತ್ ಅಡಿಯಲ್ಲಿ ಇಡುತ್ತದೆ.

`ccam add --api` Claude ಲಾಗಿನ್ ಬದಲು API ಕೀ ಬಳಸುವ ಖಾತೆಯನ್ನು ಸೇರಿಸುತ್ತದೆ; ccam ಪ್ರತಿ ಆರಂಭದಲ್ಲೂ ಕೀಯನ್ನು ಓದುತ್ತದೆ ಮತ್ತು ಅದನ್ನು ಎಂದಿಗೂ ಸಂಗ್ರಹಿಸುವುದಿಲ್ಲ. [API ಕೀ ಖಾತೆಗಳು](/features/api) ನೋಡಿ.

ಒಂದು ಖಾತೆಯಾಗಿ Claude Code ಚಾಲನೆಯಲ್ಲಿರುವಾಗ, ಅದನ್ನು ಹೇಗೆ ಆರಂಭಿಸಿದ್ದರೂ (ಹಿನ್ನೆಲೆ ಸೆಷನ್‌ಗಳು ಸೇರಿದಂತೆ), `ccam rename`, `ccam remove` ಮತ್ತು `ccam group` ಕಮಾಂಡ್‌ಗಳು ಆ ಖಾತೆಯನ್ನು ಬದಲಾಯಿಸಲು ನಿರಾಕರಿಸುತ್ತವೆ. ಮೊದಲು ಸೆಷನ್‌ನಿಂದ ನಿರ್ಗಮಿಸಿ. ಗುಂಪು ಬದಲಾವಣೆ ಅಪೂರ್ಣವಾಗಿರುವಾಗಲೂ ಅವು ನಿರಾಕರಿಸುತ್ತವೆ: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` ಪ್ರತಿ ಖಾತೆಯನ್ನು ಅದರ ಸ್ಲಾಟ್, ಶಾರ್ಟ್‌ಕಟ್, ಇಮೇಲ್ ಮತ್ತು ಫೋಲ್ಡರ್ ಜೊತೆ ಮುದ್ರಿಸುತ್ತದೆ. ಅದರ `groups` ಅರೇ ಖಾತೆಯ ಗುಂಪುಗಳನ್ನು ಪಟ್ಟಿ ಮಾಡುತ್ತದೆ, ಮತ್ತು ಅದರ `share` ಆಬ್ಜೆಕ್ಟ್ ಪ್ರತಿ ಐಟಂ (`mcp`, `settings`, `memory`, `extensions`, `history`) ಅನ್ನು ಅದು ಬರುವ ಗುಂಪಿಗೆ, ಅಥವಾ `local` ಗೆ ಹೊಂದಿಸುತ್ತದೆ.

## ಶಾರ್ಟ್‌ಕಟ್‌ಗಳು

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

ಶಾರ್ಟ್‌ಕಟ್ ಎಂದರೆ ಒಂದು ಖಾತೆಯಾಗಿ Claude Code ಅನ್ನು ಆರಂಭಿಸುವ, ನೀವೇ ಹೆಸರಿಡುವ ಕಮಾಂಡ್. `ccam work` ನಂತೆಯೇ, ಇದು ಪ್ರತಿ ಆರ್ಗ್ಯುಮೆಂಟ್ ಅನ್ನು ಬದಲಾವಣೆ ಇಲ್ಲದೆ Claude Code ಗೆ ರವಾನಿಸುತ್ತದೆ. ಹೀಗಾಗಿ `ccwork --resume` ಮತ್ತು `ccwork -p "hi"` ನೀವು ನಿರೀಕ್ಷಿಸಿದಂತೆಯೇ ಕೆಲಸ ಮಾಡುತ್ತವೆ. ಇದು `~/.local/bin` ನಲ್ಲಿರುವ ccam ಗೆ ಒಂದು ಲಿಂಕ್.

ಶಾರ್ಟ್‌ಕಟ್‌ಗಳು ಐಚ್ಛಿಕ. `ccam add` ವಿಜಾರ್ಡ್ ಒಂದನ್ನು ರಚಿಸಬೇಕೇ ಎಂದು ಕೇಳುತ್ತದೆ, `ccam add <name> --shortcut <command>` ಕೇಳದೆಯೇ ಒಂದನ್ನು ರಚಿಸುತ್ತದೆ, ಮತ್ತು `ccam list` ಪ್ರತಿ ಖಾತೆಯ ಶಾರ್ಟ್‌ಕಟ್ ಅನ್ನು, ಇಲ್ಲದಿದ್ದರೆ `-` ಅನ್ನು ತೋರಿಸುತ್ತದೆ.

ಹೆಸರುಗಳಲ್ಲಿ a-z, 0-9, `.`, `-` ಮತ್ತು `_` ಅಕ್ಷರಗಳನ್ನು ಬಳಸಬಹುದು. ಹೆಸರು ccam ಕಮಾಂಡ್ ಆಗಿರಬಾರದು, ಇನ್ನೊಂದು ಖಾತೆಯ ಶಾರ್ಟ್‌ಕಟ್ ಆಗಿರಬಾರದು, ನಿಮ್ಮ `PATH` ನಲ್ಲಿರುವ ಇನ್ನೊಂದು ಪ್ರೋಗ್ರಾಂನ ಹೆಸರಾಗಿರಬಾರದು, ಅಥವಾ ccam ರಚಿಸದ `~/.local/bin` ನಲ್ಲಿರುವ ಫೈಲ್ ಆಗಿರಬಾರದು, ಉದಾಹರಣೆಗೆ ccam ನ ನಿಮ್ಮದೇ ಪ್ರತಿ ಅಥವಾ ಲಿಂಕ್. `config.toml` ನಲ್ಲಿ ಕೈಯಾರೆ ಎಡಿಟ್ ಮಾಡಿದ ಶಾರ್ಟ್‌ಕಟ್ ಈ ನಿಯಮಗಳನ್ನು ಮುರಿದರೆ, `claude` ಅಥವಾ `ccam` ಆಗಿದ್ದರೆ, ಅಥವಾ ಎರಡು ಖಾತೆಗಳಿಗೆ ನೀಡಲಾಗಿದ್ದರೆ, `config.toml` ಲೋಡ್ ಆಗಲು ವಿಫಲವಾಗುತ್ತದೆ, ಮತ್ತು ದೋಷ ಆ ಖಾತೆಯನ್ನು ಹೆಸರಿಸುತ್ತದೆ.

`ccam rename` ಶಾರ್ಟ್‌ಕಟ್ ಹೆಸರನ್ನು ಹಾಗೆಯೇ ಇಡುತ್ತದೆ, ಮತ್ತು `ccam remove` ಖಾತೆಯ ಜೊತೆಗೆ ಶಾರ್ಟ್‌ಕಟ್ ಅನ್ನೂ ಅಳಿಸುತ್ತದೆ. `ccam shims install` ಎಲ್ಲಾ ಖಾತೆಗಳ ಲಿಂಕ್‌ಗಳನ್ನು ಮತ್ತೆ ರಚಿಸುತ್ತದೆ, ಉದಾಹರಣೆಗೆ ನೀವು ccam ಅನ್ನು ಬೇರೆಡೆಗೆ ಸ್ಥಳಾಂತರಿಸಿದ ನಂತರ; ಒಂದು ಶಾರ್ಟ್‌ಕಟ್‌ಗೆ ಅಡ್ಡಿಯಾಗಿರುವ ಫೈಲ್ ಉಳಿದವನ್ನು ನಿಲ್ಲಿಸುವುದಿಲ್ಲ.

## ಹಂಚಿಕೆ

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

ಒಂದು ಗುಂಪು ತನ್ನ ಖಾತೆಗಳ ನಡುವೆ ಐಟಂಗಳನ್ನು ಹಂಚಿಕೊಳ್ಳುತ್ತದೆ; ಒಂದು ಖಾತೆ ಹಲವು ಗುಂಪುಗಳಲ್ಲಿರಬಹುದು, ಆದರೆ ಪ್ರತಿ ಐಟಂ ಒಂದೇ ಗುಂಪಿನಿಂದ ಬರುತ್ತದೆ. [ಗುಂಪುಗಳೊಂದಿಗೆ ಹಂಚಿಕೆ](/features/share) ನೋಡಿ.

## ಈಗಿರುವ ಸೆಟಪ್ ಅನ್ನು ನಕಲಿಸುವುದು

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` ಒಂದು Claude Code ಫೋಲ್ಡರ್‌ನಿಂದ `settings.json`, `CLAUDE.md`, ವಿಸ್ತರಣೆಗಳು ಮತ್ತು ಬಳಕೆದಾರ MCP ಸರ್ವರ್‌ಗಳನ್ನು ಒಂದು ಗುಂಪಿಗೆ ನಕಲಿಸುತ್ತದೆ. ಹೀಗಾಗಿ ಆ ಐಟಂಗಳನ್ನು ಹಂಚಿಕೊಳ್ಳುವ ಪ್ರತಿ ಖಾತೆಗೂ ಅವು ಸಿಗುತ್ತವೆ. ಫೋಲ್ಡರ್ ಅನ್ನು ಓದಲಾಗುತ್ತದೆ ಅಷ್ಟೇ.

- `--settings`, `--memory`, `--extensions`, `--mcp` ಮತ್ತು `--history` ನೀವು ಹೆಸರಿಸಿದ ಐಟಂಗಳನ್ನು ಮಾತ್ರ ನಕಲಿಸುತ್ತವೆ. ಅವು ಇಲ್ಲದಿದ್ದರೆ, ಟರ್ಮಿನಲ್‌ನಲ್ಲಿ ಒಂದು ಚೆಕ್‌ಲಿಸ್ಟ್ ನಿಮಗೆ ಆಯ್ಕೆಮಾಡಲು ಬಿಡುತ್ತದೆ, ಅದರಲ್ಲಿ ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸವನ್ನು ಹೊರತುಪಡಿಸಿ ಎಲ್ಲವನ್ನೂ ಗುರುತಿಸಲಾಗಿರುತ್ತದೆ; ಬೇರೆಡೆ ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸವನ್ನು ಹೊರತುಪಡಿಸಿ ಎಲ್ಲವೂ ನಕಲಾಗುತ್ತದೆ.
- `~/.claude` ಖಾಲಿಯಾಗಿದ್ದರೂ ಅಥವಾ ಇಲ್ಲದಿದ್ದರೂ, `~/.claude.json` ನಲ್ಲಿ ಸರ್ವರ್‌ಗಳಿದ್ದರೆ `--mcp` ಕೆಲಸ ಮಾಡುತ್ತದೆ.
- `-g` ಅಥವಾ `--group` ಗುಂಪನ್ನು ಆಯ್ಕೆಮಾಡುತ್ತದೆ.
- `-y` ಅಥವಾ `--yes` ಕೇಳದೆ ನಕಲಿಸುತ್ತದೆ ಹಾಗೂ ಫೈಲ್ ಬೇರೆಯಾಗಿದ್ದಾಗ ಗುಂಪಿನ ನಕಲನ್ನು ಉಳಿಸಿಕೊಳ್ಳುತ್ತದೆ.

MCP ಸೈನ್-ಇನ್‌ಗಳನ್ನು ನಕಲಿಸಲಾಗುವುದಿಲ್ಲ, ಏಕೆಂದರೆ ನಕಲಿಸಿದ ಸೈನ್-ಇನ್ ಅದನ್ನು ಹೊಂದಿರುವ ಇತರರನ್ನು ಲಾಗ್ ಔಟ್ ಮಾಡುತ್ತದೆ; `/mcp` ಮೂಲಕ ಸೈನ್ ಇನ್ ಆಗುವ ಖಾತೆಗಳನ್ನು ಯೋಜನೆ ಹೆಸರಿಸುತ್ತದೆ. [ಈಗಿರುವ ಸೆಟಪ್ ಅನ್ನು ನಕಲಿಸುವುದು](/docs/sharing#copy-an-existing-setup) ನೋಡಿ.

## MCP ಸರ್ವರ್‌ಗಳು

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` ಆ ಖಾತೆಯಲ್ಲಿ ಮಾತ್ರ ಇರುವ ಸರ್ವರ್‌ಗಳನ್ನು ಅದರ ಗುಂಪಿಗೆ ಸ್ಥಳಾಂತರಿಸುತ್ತದೆ, ಒಂದೇ ಫೋಲ್ಡರ್‌ಗಾಗಿ ಮಾತ್ರ ಇರುವವು ಸೇರಿದಂತೆ (ಸಾಮಾನ್ಯ `claude mcp add`); ಅವು ನಂತರ ಪ್ರತಿ ಫೋಲ್ಡರ್‌ನಲ್ಲೂ ಲೋಡ್ ಆಗುತ್ತವೆ. ಗುಂಪು ಈಗಾಗಲೇ ಬೇರೆ ವ್ಯಾಖ್ಯಾನಕ್ಕಾಗಿ ಬಳಸುತ್ತಿರುವ ಹೆಸರಿನ ಸರ್ವರ್ ಅನ್ನು ಬಿಟ್ಟುಬಿಡಲಾಗುತ್ತದೆ; ಗುಂಪಿನ ನಕಲನ್ನು ಬದಲಾಯಿಸಲು `--force` ಸೇರಿಸಿ, ccam ಮೊದಲು ಅದರ ಬ್ಯಾಕಪ್ ತೆಗೆದುಕೊಳ್ಳುತ್ತದೆ.

## ಬಳಕೆ ಮತ್ತು ಬದಲಾವಣೆ

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

- `next --auto`, `--fallback` ಮತ್ತು `ccam pick` ಒಂದೇ ರೀತಿಯಲ್ಲಿ ಆಯ್ಕೆಮಾಡುತ್ತವೆ: ಮೊದಲು `auto_max_usage` (`config.toml`, ಡೀಫಾಲ್ಟ್ 100) ಗಿಂತ ಕೆಳಗಿರುವ ಖಾತೆಗಳು, `auto_strategy` (`headroom`, `expiring` ಅಥವಾ `balance`) ಪ್ರಕಾರ, `--model` ನ ವಾರದ ಅವಧಿಯನ್ನೂ ಎಣಿಸಿ. ಅವು 100% ತಲುಪಿದ ಖಾತೆಯನ್ನು (Claude Code ತುಂಬಿದೆ ಎಂದು ವರದಿ ಮಾಡಿದ ಅವಧಿಯನ್ನು ಅದರ ರೀಸೆಟ್‌ವರೆಗೆ 100% ಎಂದು ಎಣಿಸಲಾಗುತ್ತದೆ) ಅಥವಾ API ಕೀ ಲಾಗಿನ್ ಅನ್ನು ಎಂದಿಗೂ ಆಯ್ಕೆಮಾಡುವುದಿಲ್ಲ, ಮತ್ತು ಅವು ಆರಂಭವಾಗುವ ಖಾತೆಯ ಸಂಸ್ಥೆಯೊಳಗೇ ಉಳಿಯುತ್ತವೆ (ಅದರಾಚೆ ನೋಡಲು `--cross-org` ಅಥವಾ `cross_org = true`).
- `--wait` ತನ್ನ ಮೌಲ್ಯವನ್ನು `=` ನಂತರ ಪಡೆಯುತ್ತದೆ (`--wait=2h`); ಸ್ಪೇಸ್‌ನೊಂದಿಗೆ ಬರೆದರೆ, ಒಂದು ಸುಳಿವಿನೊಂದಿಗೆ ಅದನ್ನು ನಿರಾಕರಿಸಲಾಗುತ್ತದೆ. ಯಾವುದನ್ನೂ ಆಯ್ಕೆಮಾಡಲಾಗದಿದ್ದಾಗ ಅವು ಕೋಡ್ 75 ನೊಂದಿಗೆ ನಿರ್ಗಮಿಸುತ್ತವೆ, ಅಥವಾ, `--wait` ಇದ್ದರೆ, ಮೊದಲ ಖಾತೆ ರೀಸೆಟ್ ಆಗುವವರೆಗೆ ಕಾಯುತ್ತವೆ.
- ಇತಿಹಾಸ ಹಂಚಿಕೊಳ್ಳದ ಖಾತೆಗಳ ನಡುವೆಯೂ `ccam next --continue` ಸಂಭಾಷಣೆಯನ್ನು ಕೊಂಡೊಯ್ಯುತ್ತದೆ: ಆಗ Claude Code ಅದನ್ನು ಹೊಸ ಖಾತೆಗೆ ನಕಲಿಸುತ್ತದೆ. `--brief` ಜೊತೆ, ಹೊಸ ಖಾತೆ ಬದಲಿಗೆ ಹೊಸ ಸಂಭಾಷಣೆಯನ್ನು ಆರಂಭಿಸುತ್ತದೆ, ಹಿಂದಿನದರ ಬಗ್ಗೆ ccam ಮಾಡೆಲ್ ಇಲ್ಲದೆ ಬರೆಯುವ ಸಾರಾಂಶದಿಂದ ([`--brief`](/features/next#brief)).
- `ccam <name>` ನಂತಲ್ಲದೆ, `ccam next` ತನ್ನ ಆಯ್ಕೆಗಳನ್ನು ತಾನೇ ಓದುತ್ತದೆ, ಆದ್ದರಿಂದ Claude Code ನ ಆಯ್ಕೆಗಳು `--` ನಂತರ ಬರುತ್ತವೆ; ಅದಿಲ್ಲದೆ ccam `unknown flag` ನೊಂದಿಗೆ ನಿಲ್ಲುತ್ತದೆ.
- `ccam agent` ಈ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ `claude -p` ಮೂಲಕ ಒಂದು ಕಾರ್ಯವನ್ನು, ಅದನ್ನು ಕರೆದ ಖಾತೆಯನ್ನು ಹೊರತುಪಡಿಸಿ ಇನ್ನೊಂದು ಖಾತೆಯಾಗಿ, ಕರೆದ ಸೆಷನ್‌ನ ಅನುಮತಿ ಮೋಡ್ ಮತ್ತು effort ನಲ್ಲಿ ರನ್ ಮಾಡುತ್ತದೆ, ಮತ್ತು Claude ನ ಕೊನೆಯ ಉತ್ತರವನ್ನು ಮುದ್ರಿಸುತ್ತದೆ. ಆ ಖಾತೆಯ ಮಿತಿಯಲ್ಲಿ ಸಂಭಾಷಣೆ ಇನ್ನೊಂದು ಖಾತೆಯಾಗಿ ಮುಂದುವರಿಯುತ್ತದೆ ([ಇತರ ಖಾತೆಗಳಲ್ಲಿ subagent ಗಳು](/features/agent)).
- Claude Code ಒಳಗೆ, ಪ್ಲಗಿನ್‌ನ `/ccam:switch [account]` ಸಂಭಾಷಣೆಯನ್ನು ಮುಂದಕ್ಕೆ ಒಪ್ಪಿಸುತ್ತದೆ, ಮತ್ತು ಅದರ ಹುಕ್‌ಗಳು ರನ್ ಮಾಡುವುದು `ccam hook`.
- `config.toml` ನಲ್ಲಿ `auto_switch = true` ಇದ್ದರೆ (Linux ಮತ್ತು macOS; ಡೀಫಾಲ್ಟ್ ಆಗಿ ಆಫ್), ಮಿತಿಯ ಹುಕ್ ccam ಆರಂಭಿಸಿದ ಇಂಟರ್‌ಆಕ್ಟಿವ್ ಸೆಷನ್ ಅನ್ನು ಕೊನೆಗೊಳಿಸಿ, ಅವಕಾಶ ಇರುವ ಇನ್ನೊಂದು ಖಾತೆಯಲ್ಲಿ ಅದನ್ನು ಮುಂದುವರಿಸುತ್ತದೆ, ಆದರೆ ಅದೇ ಬಗೆಯ ಖಾತೆಯಲ್ಲಿ ಮಾತ್ರ: ಅದೇ ಸಂಸ್ಥೆಯ ಇನ್ನೊಂದು ಸೀಟು, ಅಥವಾ ಇನ್ನೊಂದು ವೈಯಕ್ತಿಕ ಪ್ಲಾನ್ ([`auto_switch`](/features/next#auto-switch)).

[ಮುಂದಿನ ಖಾತೆಗೆ ಬದಲಾಯಿಸುವುದು](/features/next), [ಸ್ಕ್ರಿಪ್ಟ್‌ಗಳಿಗಾಗಿ ಖಾತೆ ಆಯ್ಕೆ](/features/pick), [ಇತರ ಖಾತೆಗಳಲ್ಲಿ subagent ಗಳು](/features/agent) ಮತ್ತು [ಪ್ರತಿ ಖಾತೆಯ ಹಿನ್ನೆಲೆ ಸೆಷನ್‌ಗಳು](/features/agents) ನೋಡಿ.

## ನಿರ್ವಹಣೆ

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

ccam ಅನ್ನು ಬದಲಾಯಿಸುವ ಮೊದಲು `ccam update` ಬಿಡುಗಡೆಯ SHA-256 ಚೆಕ್‌ಸಮ್ ಅನ್ನು ಪರಿಶೀಲಿಸುತ್ತದೆ, ಮತ್ತು `--version v0.2.0` ಒಂದು ನಿರ್ದಿಷ್ಟ ಬಿಡುಗಡೆಯನ್ನು ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡುತ್ತದೆ. `.deb` ಪ್ಯಾಕೇಜ್‌ನಿಂದ ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿದ ccam ಪ್ರತಿಯನ್ನು ಬದಲಿಗೆ apt ಮೂಲಕ ಅಪ್‌ಡೇಟ್ ಮಾಡಲಾಗುತ್ತದೆ. [ಅಪ್‌ಡೇಟ್](/docs/install#update) ನೋಡಿ.

`ccam uninstall` ತಾನು ಏನನ್ನು ತೆಗೆದುಹಾಕುತ್ತದೆ ಎಂದು ಪಟ್ಟಿ ಮಾಡಿ ಮೊದಲು ಕೇಳುತ್ತದೆ: `ccam remove` ಮಾಡುವಂತೆ ಪ್ರತಿ ಖಾತೆ, ಶಾರ್ಟ್‌ಕಟ್ ಕಮಾಂಡ್‌ಗಳು, ccam ನ ಫೋಲ್ಡರ್‌ಗಳು ಮತ್ತು `ccam` ಕಮಾಂಡ್. ಇಂಪೋರ್ಟ್ ಮಾಡಿದ ಫೋಲ್ಡರ್ ಇದ್ದಲ್ಲೇ ಉಳಿಯುತ್ತದೆ. [ಅನ್‌ಇನ್‌ಸ್ಟಾಲ್](/docs/install#uninstall) ನೋಡಿ.

## ಫೈಲ್‌ಗಳು

| ಪಾತ್ | ಒಳಗೊಂಡಿರುವುದು |
| --- | --- |
| `~/.config/ccam/config.toml` | ಖಾತೆಗಳು, ಗುಂಪುಗಳು ಮತ್ತು ಸೆಟ್ಟಿಂಗ್‌ಗಳು (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). ಸೆಟ್ಟಿಂಗ್‌ಗಳು ಮೇಲ್ಭಾಗದಲ್ಲಿ, ಮೊದಲ `[[group]]`, `[[profile]]` ಅಥವಾ `[[folder]]` ಗಿಂತ ಮೇಲೆ ಇರಬೇಕು: ಟೇಬಲ್‌ನ ಕೆಳಗಿರುವ ಸೆಟ್ಟಿಂಗ್ ಆ ಟೇಬಲ್‌ಗೆ ಸೇರುತ್ತದೆ, ಆದ್ದರಿಂದ ccam ಪ್ರತಿ ಆರಂಭದಲ್ಲೂ ಮತ್ತು `ccam doctor` ನಲ್ಲಿ ಅದರ ಬಗ್ಗೆ ಎಚ್ಚರಿಸುತ್ತದೆ, ಮತ್ತು ಕಾನ್ಫಿಗ್ ಅನ್ನು ಉಳಿಸುವ ಕಮಾಂಡ್‌ಗಳು ಅದನ್ನು ಸ್ಥಳಾಂತರಿಸುವವರೆಗೆ ನಿರಾಕರಿಸುತ್ತವೆ. ccam ಫೈಲ್ ಅನ್ನು ಉಳಿಸುವಾಗ ಕಾಮೆಂಟ್‌ಗಳು ಉಳಿಯುವುದಿಲ್ಲ |
| `~/.local/share/ccam/profiles/<name>/` | ಪ್ರತಿ ಖಾತೆಯ `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | ಹಂಚಿಕೊಂಡ ಸೆಟ್ಟಿಂಗ್‌ಗಳು, MCP ಸರ್ವರ್‌ಗಳು, ವಿಸ್ತರಣೆಗಳು ಮತ್ತು ಇತಿಹಾಸ |
| `~/.local/state/ccam/` | ಬಳಕೆಯ ಕ್ಯಾಶ್, ಮಿತಿ ದಾಖಲೆಗಳು, ಲಾಕ್‌ಗಳು, ಟಿಪ್ಪಣಿಗಳು ಮತ್ತು ಇತರ ಸ್ಥಿತಿ (ಟೇಬಲ್‌ನ ಕೆಳಗಿನ ಪಟ್ಟಿ ನೋಡಿ) |
| `~/.local/bin/<shortcut>` | `ccwork` ನಂತಹ ಶಾರ್ಟ್‌ಕಟ್ ಕಮಾಂಡ್‌ಗಳು, ccam ಗೆ ಲಿಂಕ್‌ಗಳು |

ಸ್ಥಿತಿ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ ಇರುವುದು:

- ಬಳಕೆಯ ಕ್ಯಾಶ್ ಮತ್ತು ಮಿತಿ ದಾಖಲೆಗಳು (`usage/<name>.json`, `usage/<name>.limit.json`);
- ಕೊನೆಯ ಬಾರಿ ಬಳಸಿದ ಸಮಯಗಳು;
- ಸೆಷನ್ ಗುರುತುಗಳು ಮತ್ತು ಲಾಕ್‌ಗಳು (`locks/`);
- ccam ರಚಿಸಿದ ಶಾರ್ಟ್‌ಕಟ್‌ಗಳು (`shortcuts.json`);
- ಪೂರ್ಣಗೊಳ್ಳದ ಗುಂಪು ಬದಲಾವಣೆ (`group-change.json`);
- `/ccam:switch` ಹಾಗೂ ಮಿತಿಯ ಹುಕ್ ಬಿಡುವ ಟಿಪ್ಪಣಿಗಳು (`handoff/<account>/<session>.json`);
- `ccam pick` ಅಥವಾ `auto_switch` ಪ್ರತಿ ಖಾತೆಯನ್ನು ಕೊನೆಯದಾಗಿ ಯಾವಾಗ ಆಯ್ಕೆಮಾಡಿತು ಎಂಬುದು (`picks/<name>`; ಒಂದು ಆಯ್ಕೆ 2 ನಿಮಿಷಗಳವರೆಗೆ ಎಣಿಕೆಯಾಗುತ್ತದೆ);
- 7 ದಿನಗಳವರೆಗೆ ಇರಿಸಲಾಗುವ `--brief` ನ ಸಾರಾಂಶಗಳು (`brief/<session>.md`).

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` ಅಥವಾ `XDG_STATE_HOME` ಅನ್ನು ಹೊಂದಿಸಿದರೆ ಈ ಫೋಲ್ಡರ್‌ಗಳ ಸ್ಥಳ ಬದಲಾಗುತ್ತದೆ. Windows ನಲ್ಲಿ, `~` ಎಂದರೆ `%USERPROFILE%`, ಮತ್ತು ಶಾರ್ಟ್‌ಕಟ್‌ಗಳು `<shortcut>.exe`.
