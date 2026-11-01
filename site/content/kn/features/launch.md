---
title: Claude Code ಆರಂಭಿಸುವುದು
description: ccam <name|number>, ccam run ಅಥವಾ ಡ್ಯಾಶ್‌ಬೋರ್ಡ್ ಮೂಲಕ ಯಾವುದೇ ಖಾತೆಯಾಗಿ Claude Code ಆರಂಭಿಸಿ; ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳು ಬದಲಾಗದೆ claude ಗೆ ಹೋಗುತ್ತವೆ.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: "11df08608a1e"
---

`ccam work` `work` ಖಾತೆಯಾಗಿ Claude Code ಅನ್ನು ಆರಂಭಿಸುತ್ತದೆ, ಮತ್ತು ಹೆಸರಿನ ನಂತರದ ಎಲ್ಲವೂ ಬದಲಾವಣೆ ಇಲ್ಲದೆ `claude` ಗೆ ಹೋಗುತ್ತದೆ. ನೀವು `claude` ಎಂದು ಟೈಪ್ ಮಾಡುವ ಎಲ್ಲೆಡೆ ಇದನ್ನು ಬಳಸಿ: ಟರ್ಮಿನಲ್‌ನಲ್ಲಿ, ಸ್ಕ್ರಿಪ್ಟ್‌ನಲ್ಲಿ, ಅಥವಾ ಪೈಪ್ ಮೂಲಕ ಇನ್‌ಪುಟ್ ನೀಡುವಾಗ.

## ಬಳಕೆ

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

ಖಾತೆಯ ಹೆಸರನ್ನು ಅಥವಾ `ccam list` ನಲ್ಲಿರುವ ಅದರ ಸಂಖ್ಯೆಯನ್ನು ನೀಡಿ. ಹೆಸರಿನ ನಂತರದ ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳನ್ನು ccam ಓದುವುದಿಲ್ಲ, ಹೀಗಾಗಿ `ccam work --help` Claude Code ನ ಸಹಾಯವನ್ನು ತೋರಿಸುತ್ತದೆ. ಹೆಸರಿನ ನಂತರ `--` ಐಚ್ಛಿಕ. ಸ್ಟ್ಯಾಂಡರ್ಡ್ ಇನ್‌ಪುಟ್, ಎಕ್ಸಿಟ್ ಕೋಡ್ ಮತ್ತು ಸಿಗ್ನಲ್‌ಗಳು `claude` ನಲ್ಲಿ ಇರುವಂತೆಯೇ ವರ್ತಿಸುತ್ತವೆ.

ಮೊದಲ ಆರ್ಗ್ಯುಮೆಂಟ್ `-` ನಿಂದ ಆರಂಭವಾಗಿ ccam ನ ಸ್ವಂತ ಫ್ಲ್ಯಾಗ್‌ಗಳಲ್ಲಿ ಒಂದಲ್ಲದಿದ್ದರೆ, ccam ನಿಮಗಾಗಿ ಖಾತೆಯನ್ನು ಆಯ್ಕೆಮಾಡಿ ಎಲ್ಲಾ ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳನ್ನು Claude Code ಗೆ ರವಾನಿಸುತ್ತದೆ. ಅದು [`ccam folder`](/features/folder) ಮೂಲಕ ಈಗಿನ ಫೋಲ್ಡರ್‌ಗೆ ಹೊಂದಿಸಿದ ಖಾತೆಯನ್ನು ಬಳಸುತ್ತದೆ, ಇಲ್ಲದಿದ್ದರೆ ಡೀಫಾಲ್ಟ್ ಖಾತೆಯನ್ನು: `ccam use` ಮೂಲಕ ಹೊಂದಿಸಿದ ಖಾತೆ, ಅಥವಾ ನೀವು ಹೊಂದಿಸದಿದ್ದರೆ ಅತ್ಯಂತ ಕಡಿಮೆ ಸಂಖ್ಯೆಯ ಖಾತೆ.

### ಡ್ಯಾಶ್‌ಬೋರ್ಡ್

ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳಿಲ್ಲದ `ccam`, ಸ್ಟ್ಯಾಂಡರ್ಡ್ ಇನ್‌ಪುಟ್ ಮತ್ತು ಔಟ್‌ಪುಟ್ ಟರ್ಮಿನಲ್ ಆಗಿದ್ದಾಗ [ಡ್ಯಾಶ್‌ಬೋರ್ಡ್](/features/dashboard) ಅನ್ನು ತೆರೆಯುತ್ತದೆ, ಇಲ್ಲದಿದ್ದರೆ `ccam list` ಅನ್ನು ಮುದ್ರಿಸುತ್ತದೆ. ಒಂದು ಖಾತೆಗೆ ಹೋಗಿ (ಬಾಣದ ಕೀಲಿಗಳಿಂದ, ಅಥವಾ ಅದರ ಸಂಖ್ಯೆಯನ್ನು ಟೈಪ್ ಮಾಡಿ) ಅದನ್ನು ಆರಂಭಿಸಲು `Enter` ಒತ್ತಿ; `n` ಹೆಚ್ಚು ಅವಕಾಶ ಇರುವ ಖಾತೆಯನ್ನು ಆರಂಭಿಸುತ್ತದೆ ಮತ್ತು `o` `--continue`, `--resume` ಅಥವಾ ಇತರ ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳನ್ನು ನೀಡುತ್ತದೆ. Claude Code ಆರಂಭವಾಗುವ ಮೊದಲು ಡ್ಯಾಶ್‌ಬೋರ್ಡ್ ಮುಚ್ಚುತ್ತದೆ, ಹೀಗಾಗಿ ಟರ್ಮಿನಲ್ ಮೊದಲು ಮೊದಲಿನ ಸ್ಥಿತಿಗೆ ಮರಳುತ್ತದೆ.

**Team ಮತ್ತು Enterprise ಸೀಟುಗಳು.** ಸ್ವಯಂಚಾಲಿತ ಖಾತೆ ಆಯ್ಕೆ (`ccam next`, `--fallback`, `ccam pick`, ಡ್ಯಾಶ್‌ಬೋರ್ಡ್‌ನ `n` ಮತ್ತು, ನೀವು ಆನ್ ಮಾಡಿದ ನಂತರ, [`auto_switch`](/features/next#auto-switch)) ಒಂದೇ Claude Team ಅಥವಾ Enterprise ಸಂಸ್ಥೆಯ ಸೀಟುಗಳ ನಡುವೆಯೂ ಬದಲಾಯಿಸುತ್ತದೆ, ಇದು Anthropic ನ ನಿಯಮಗಳನ್ನು ಉಲ್ಲಂಘಿಸಬಹುದು. ಇದನ್ನು ಅವಲಂಬಿಸುವ ಮೊದಲು [ಪದೇ ಪದೇ ಕೇಳುವ ಪ್ರಶ್ನೆಗಳನ್ನು](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) ಓದಿ ಮತ್ತು ನಿಮ್ಮ ಸಂಸ್ಥೆ ಹಾಗೂ Anthropic ಬಳಿ ವಿಚಾರಿಸಿ.

ಶೆಲ್‌ನಲ್ಲಿ `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` ಅಥವಾ `ANTHROPIC_AUTH_TOKEN` ಹೊಂದಿಸಿದ್ದರೆ, Claude Code ಖಾತೆಯ ಸ್ವಂತ ಲಾಗಿನ್ ಬದಲು ಆ ಕೀ ಅಥವಾ ಟೋಕನ್ ಅನ್ನು ಬಳಸಬಹುದು, ಹೀಗಾಗಿ ccam ಆರಂಭಿಸುವ ಮೊದಲು ಸ್ಟ್ಯಾಂಡರ್ಡ್ ಎರರ್‌ನಲ್ಲಿ ನಿಮ್ಮನ್ನು ಎಚ್ಚರಿಸುತ್ತದೆ (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). ಅದು ವೇರಿಯಬಲ್‌ಗಳನ್ನು ಓದುತ್ತದೆ ಅಷ್ಟೇ; [API ಕೀ ಖಾತೆ](/features/api) ಬದಲಿಗೆ ಅವುಗಳನ್ನು ತೆರವುಗೊಳಿಸುತ್ತದೆ.

ಡೀಫಾಲ್ಟ್ ಆಗಿ, ccam ತನ್ನ ಜಾಗವನ್ನು Claude Code ಗೆ ಬಿಟ್ಟುಕೊಡುತ್ತದೆ, ಹೀಗಾಗಿ ಬೇರೇನೂ ರನ್ ಆಗುತ್ತಿರುವುದಿಲ್ಲ. `config.toml` ನ ಮೇಲ್ಭಾಗದಲ್ಲಿ, ಮೊದಲ `[[group]]`, `[[profile]]` ಅಥವಾ `[[folder]]` ಗಿಂತ ಮೇಲೆ, `relay = true` (ಅಥವಾ Linux ಮತ್ತು macOS ನಲ್ಲಿ [`auto_switch = true`](/features/next#auto-switch)) ಇದ್ದರೆ, ccam ಬದಲಿಗೆ ಇಂಟರ್‌ಆಕ್ಟಿವ್ ಸೆಷನ್‌ಗಳನ್ನು ಚೈಲ್ಡ್ ಪ್ರೊಸೆಸ್ ಆಗಿ ಆರಂಭಿಸಿ ಅವುಗಳ ಜೊತೆಯಲ್ಲೇ ರನ್ ಆಗುತ್ತಿರುತ್ತದೆ. ಇದರಿಂದ `/ccam:switch`, ಅಥವಾ ಬಳಕೆಯ ಮಿತಿಯಲ್ಲಿ `auto_switch`, ಸೆಷನ್‌ನ ಆಯ್ಕೆಗಳು, ಅನುಮತಿ ಮೋಡ್, effort ಮಟ್ಟ ಮತ್ತು ಏಜೆಂಟ್‌ನೊಂದಿಗೆ, ತಾನಾಗಿಯೇ ಸಂಭಾಷಣೆಯನ್ನು ಇನ್ನೊಂದು ಖಾತೆಯಲ್ಲಿ ಮುಂದುವರಿಸಬಲ್ಲದು; `/ccam:switch` ಸಂಭಾಷಣೆಯನ್ನು ಇನ್ನೊಂದು ಸಂಸ್ಥೆಯ ಖಾತೆಗೆ ಸ್ಥಳಾಂತರಿಸುವ ಮೊದಲು, ccam ನಿಮ್ಮನ್ನು ಕೇಳುತ್ತದೆ ([ವಿವರಗಳು](/features/next#relay)).

## ಆಯ್ಕೆಗಳು

| ಆಯ್ಕೆ | ಏನು ಮಾಡುತ್ತದೆ | ಡೀಫಾಲ್ಟ್ |
| --- | --- | --- |
| `--dry-run` | ccam ರನ್ ಮಾಡುವ ಪರಿಸರ ಮತ್ತು ಕಮಾಂಡ್ ಅನ್ನು ಮುದ್ರಿಸಿ ನಿರ್ಗಮಿಸುತ್ತದೆ | ಆಫ್ |
| `--fallback` | `-p`/`--print` ರನ್‌ಗಳಿಗಾಗಿ: ಖಾತೆ ತನ್ನ ಬಳಕೆಯ ಮಿತಿ ತಲುಪಿದಾಗ, ಅದೇ ಕಮಾಂಡ್ ಅನ್ನು ಹೆಚ್ಚು ಅವಕಾಶ ಇರುವ ಲಾಗಿನ್ ಆದ ಖಾತೆಯಾಗಿ ಮತ್ತೆ ರನ್ ಮಾಡುತ್ತದೆ ([ವಿವರಗಳು](/features/next#run-again-at-the-limit)) | ಆಫ್ |
| `-v`, `--version` | `ccam <version>` ಅನ್ನು ಮುದ್ರಿಸುತ್ತದೆ (ಮೊದಲ ಆರ್ಗ್ಯುಮೆಂಟ್ ಆಗಿ ಮಾತ್ರ) | |
| `-h`, `--help` | ccam ನ ಸಹಾಯ, `ccam --help` ಅಥವಾ `ccam run --help` ಆಗಿ | |

ಆರ್ಗ್ಯುಮೆಂಟ್ ನಿಯಮಗಳು:

- ccam ತನ್ನ ಸ್ವಂತ ಫ್ಲ್ಯಾಗ್‌ಗಳನ್ನು ಖಾತೆಯ ಹೆಸರಿನ ಮೊದಲು ಮಾತ್ರ ಓದುತ್ತದೆ. ಹೆಸರಿನ ನಂತರ, `--dry-run` ಮತ್ತು `--help` ಸೇರಿದಂತೆ ಪ್ರತಿ ಆರ್ಗ್ಯುಮೆಂಟ್ Claude Code ಗೆ ಸೇರಿದ್ದು.
- `--dry-run` `ccam --dry-run work`, `ccam run --dry-run work` ಮತ್ತು `ccam --dry-run` (ಫೋಲ್ಡರ್‌ನ ಖಾತೆ ಅಥವಾ ಡೀಫಾಲ್ಟ್‌ಗಾಗಿ) ಆಗಿ ಕೆಲಸ ಮಾಡುತ್ತದೆ. `--dry-run` ಮತ್ತು `--fallback` ಯಾವುದೇ ಕ್ರಮದಲ್ಲಿ ಬರಬಹುದು.
- `ccam run <name|number> [claude args]` ಎಂಬುದು `ccam <name>` ನ ದೀರ್ಘ ರೂಪ. ಹೆಸರಿಲ್ಲದ `ccam run` ಫೋಲ್ಡರ್‌ನ ಖಾತೆಯನ್ನು ಅಥವಾ ಡೀಫಾಲ್ಟ್ ಖಾತೆಯನ್ನು ಬಳಸುತ್ತದೆ.
- ಖಾತೆಯ ಹೆಸರುಗಳು `add`, `list`, `run` ಅಥವಾ `usage` ನಂತಹ ccam ಕಮಾಂಡ್ ಹೆಸರುಗಳಾಗಿರಬಾರದು, ಹೀಗಾಗಿ `ccam <word>` ಎಂದಿಗೂ ಅಸ್ಪಷ್ಟವಾಗುವುದಿಲ್ಲ.

ಡ್ರೈ ರನ್ ನೀವು ಶೆಲ್‌ಗೆ, ಅಥವಾ Windows ನಲ್ಲಿ PowerShell ಗೆ, ಪೇಸ್ಟ್ ಮಾಡಬಹುದಾದ ಒಂದು ಸಾಲನ್ನು ಮುದ್ರಿಸುತ್ತದೆ. ಆರಂಭವು ಪರಿಸರದಲ್ಲಿ ಏನನ್ನು ಬದಲಾಯಿಸುತ್ತದೆ ಎಂಬುದರಿಂದ ಅದು ಶುರುವಾಗುತ್ತದೆ: ಯಾವಾಗಲೂ `CLAUDE_CONFIG_DIR` ಮತ್ತು `CCAM_PROFILE`, ಕೈಬಿಡುವ ವೇರಿಯಬಲ್‌ಗೆ `env -u <name>` (ಉದಾಹರಣೆಗೆ ಆನುವಂಶಿಕವಾಗಿ ಬಂದ `CCAM_RELAY`), ಮತ್ತು, [API ಕೀ ಖಾತೆ](/features/api) ಗೆ, ಅದರ ಎಂಡ್‌ಪಾಯಿಂಟ್, ಅದರ ಮಾಡೆಲ್ ಮತ್ತು ಕೀ ಎಲ್ಲಿಂದ ಬರುತ್ತದೆ ಎಂಬುದು (`ANTHROPIC_API_KEY="$WORK_KEY"`), ಕೀಯನ್ನೇ ಎಂದಿಗೂ ಅಲ್ಲ:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

ಶೆಲ್‌ನಲ್ಲಿ `VAR=x cmd` ಮಾಡುವಂತೆಯೇ, PowerShell ಸಾಲು ನಂತರ ವೇರಿಯಬಲ್‌ಗಳನ್ನು ಮೊದಲಿನಂತೆ ಮರಳಿ ಇರಿಸುತ್ತದೆ, ಮತ್ತು Claude Code ಗೆ ಅದರ ಕಮಾಂಡ್ ಲೈನ್ ಅನ್ನು `Start-Process` ಮೂಲಕ ನೀಡುತ್ತದೆ, ಹೀಗಾಗಿ Windows PowerShell 5.1 ಮತ್ತು PowerShell 7 ಅದನ್ನು ಒಂದೇ ರೀತಿಯಲ್ಲಿ ರನ್ ಮಾಡುತ್ತವೆ. ಆರಂಭವು ಒಂದು ವೇರಿಯಬಲ್ ಅನ್ನು ಖಾಲಿ ಮೌಲ್ಯಕ್ಕೆ ಹೊಂದಿಸಿದಾಗ, 7.5 ಕ್ಕಿಂತ ಹಿಂದಿನ PowerShell ಅಂತಹ ವೇರಿಯಬಲ್ ಅನ್ನು ಬದಲಿಗೆ ತೆಗೆದುಹಾಕುತ್ತದೆ ಎಂದು ಅದರ ಕೆಳಗಿನ `#` ಸಾಲು ಹೇಳುತ್ತದೆ. npm ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿದ `claude.cmd` ಗೆ, ಆರಂಭ ಬಳಸುವ `cmd.exe` ಕಮಾಂಡ್ ಅನ್ನು ಸಾಲು ರನ್ ಮಾಡುತ್ತದೆ (ಕೆಳಗೆ ನೋಡಿ), ಮತ್ತು ಅಲ್ಲಿ ಆರಂಭ ನಿರಾಕರಿಸುವ ಆರ್ಗ್ಯುಮೆಂಟ್ ಡ್ರೈ ರನ್ ಅನ್ನೂ ವಿಫಲಗೊಳಿಸುತ್ತದೆ.

## ಇದು ಹೇಗೆ ಕೆಲಸ ಮಾಡುತ್ತದೆ

ಪ್ರತಿ ಆರಂಭವೂ ಕಮಾಂಡ್ ಲೈನ್‌ನಿಂದ `exec` ವರೆಗೆ ಅದೇ ಹಂತಗಳನ್ನು ನಡೆಸುತ್ತದೆ:

<!-- motion:launch -->

ಪ್ರತಿ ಹಂತ ನಡೆಯುವಾಗ ಅದರ ಚೌಕ ಬೆಳಗುತ್ತದೆ, ಮತ್ತು ಕೆಳಗಿನ ಸಾಲುಗಳು Claude Code ಆರಂಭವಾಗುವ ಪರಿಸರ ಮತ್ತು ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳು.

1. ccam `app.Main` ನಲ್ಲಿ ಆರಂಭವಾಗುತ್ತದೆ, ಅದನ್ನು `cmd/ccam/main.go` ಕರೆಯುತ್ತದೆ. ccam ಬೇರೆ ಹೆಸರಿನಲ್ಲಿ ಆರಂಭವಾಗಿದ್ದರೆ, ಅದು ಒಂದು [ಶಾರ್ಟ್‌ಕಟ್](/features/shortcuts) ಮತ್ತು ನೇರವಾಗಿ ಲಾಂಚರ್‌ಗೆ ಹೋಗುತ್ತದೆ.
2. cobra ತೊಡಗುವ ಮೊದಲೇ ccam `os.Args` ಅನ್ನು ನೋಡುತ್ತದೆ (`app.parse`). ಕೇವಲ `ccam` (ಡ್ಯಾಶ್‌ಬೋರ್ಡ್), ccam ಕಮಾಂಡ್ ಹೆಸರು (`config.Reserved`), `-h`/`--help` ಮತ್ತು `-v`/`--version` ಮಾತ್ರ `cli.Execute` ಗೆ ಹೋಗುತ್ತವೆ. ಉಳಿದೆಲ್ಲವೂ ಆರಂಭ.
3. ccam `~/.config/ccam/config.toml` ಅನ್ನು ಓದುತ್ತದೆ, ಅದು ಪಾರ್ಸ್ ಮಾಡುವ ಏಕೈಕ ಫೈಲ್ (`launcher.Launch`).
4. ccam ಹೆಸರು ಅಥವಾ ಸಂಖ್ಯೆಯಿಂದ ಖಾತೆಯನ್ನು ಕಂಡುಹಿಡಿಯುತ್ತದೆ, ಮುಂದಿರುವ ಒಂದು `--` ಅನ್ನು ತೆಗೆದುಹಾಕುತ್ತದೆ, ಮತ್ತು `claude` ಎಕ್ಸಿಕ್ಯೂಟಬಲ್ ಅನ್ನು ಕಂಡುಹಿಡಿಯುತ್ತದೆ (`launcher.Prepare`, `ResolveClaude`). ಹೆಸರಿಲ್ಲದಿದ್ದರೆ, ಅದು ಕೆಲಸದ ಡೈರೆಕ್ಟರಿಯನ್ನು `config.toml` ನ `[[folder]]` ನಮೂದುಗಳೊಂದಿಗೆ ಸ್ಟ್ರಿಂಗ್‌ಗಳಾಗಿ, macOS ಮತ್ತು Windows ನಲ್ಲಿ ದೊಡ್ಡ-ಸಣ್ಣ ಅಕ್ಷರಗಳನ್ನು ಲೆಕ್ಕಿಸದೆ, ಹೋಲಿಸಿ ಫೋಲ್ಡರ್‌ನ ಖಾತೆಯನ್ನು ತೆಗೆದುಕೊಳ್ಳುತ್ತದೆ (`HereProfile`). PATH ಹುಡುಕಾಟವನ್ನು ಫೈಲ್‌ನ ಮಾರ್ಪಾಡು ಸಮಯ ಮತ್ತು `PATH` ನ ಒಂದು ಹ್ಯಾಶ್‌ನೊಂದಿಗೆ `~/.local/state/ccam/claude-path` ನಲ್ಲಿ ಕ್ಯಾಶ್ ಮಾಡಲಾಗುತ್ತದೆ, ಹೀಗಾಗಿ ccam ಪ್ರತಿ ಆರಂಭದಲ್ಲೂ `PATH` ಅನ್ನು ಹುಡುಕುವುದಿಲ್ಲ. ಫೈಲ್ ಬದಲಾದರೆ ಅಥವಾ ಇಲ್ಲದಿದ್ದರೆ, ಅಥವಾ `PATH` ಬದಲಾದರೆ, ಹೊಸ ಹುಡುಕಾಟ ನಡೆಯುತ್ತದೆ. `config.toml` ನಲ್ಲಿರುವ `claude_bin` ಪಾತ್ (`~/` ಅನ್ನು ವಿಸ್ತರಿಸಲಾಗುತ್ತದೆ) ಕಾರ್ಯಗತಗೊಳಿಸಬಹುದಾದ ಫೈಲ್ ಆಗಿರಬೇಕು; ಇಲ್ಲದಿದ್ದರೆ ಏನನ್ನಾದರೂ ಬರೆಯುವ ಮೊದಲೇ ಆರಂಭ `claude_bin X is not an executable file. Check claude_bin in config.toml` ಎಂದು ನಿಲ್ಲುತ್ತದೆ. ccam ಸ್ವತಃ ccam ಆಗಿರುವ `claude` ಅನ್ನು inode ಮೂಲಕ ಹೋಲಿಸಿ ನಿರಾಕರಿಸುತ್ತದೆ, ಏಕೆಂದರೆ ಅದು ಲೂಪ್ ಆಗುತ್ತದೆ.
5. ccam ಖಾತೆಯ ಫೋಲ್ಡರ್ `~/.local/share/ccam/profiles/<name>/` ಅನ್ನು ಅದರ ಹಂಚಿಕೆ ಸೆಟ್ಟಿಂಗ್‌ಗಳೊಂದಿಗೆ ಹೋಲಿಸುತ್ತದೆ, `lstat` ಮತ್ತು `readlink` ಅನ್ನು ಮಾತ್ರ ಕರೆದು, ಪ್ರತಿ ಐಟಂಗೆ ಒಂದು ಅಥವಾ ಎರಡು ಸಿಸ್ಟಮ್ ಕಾಲ್‌ಗಳು (`checkLinks`, `share.Check`). ಇಲ್ಲದ ಲಿಂಕ್ ಅನ್ನು ಅಲ್ಲೇ ರಚಿಸಲಾಗುತ್ತದೆ. ಡೇಟಾ ಕಳೆದುಹೋಗಬಹುದಾದ ಯಾವುದಾದರೂ ಆರಂಭವನ್ನು ನಿಲ್ಲಿಸುತ್ತದೆ.
6. ccam ಹಂಚಿಕೊಂಡ MCP ಸರ್ವರ್‌ಗಳನ್ನು ಸೇರಿಸುತ್ತದೆ (`applyMCP`; ಕೆಳಗೆ ನೋಡಿ).
7. ccam ಸೆಷನ್ ಅನ್ನು ದಾಖಲಿಸುತ್ತದೆ (`session.Register`). ಅದು ತನ್ನ PID ಹೆಸರಿನ, ಪ್ರಕ್ರಿಯೆಯ ಆರಂಭ ಸಮಯವನ್ನು ಹೊಂದಿರುವ ಒಂದು ಗುರುತನ್ನು `~/.local/state/ccam/locks/profile-<name>/<pid>` ಗೆ ಬರೆಯುತ್ತದೆ, `XDG_RUNTIME_DIR` ಏನೇ ಹೇಳಿದರೂ, ಹೀಗಾಗಿ ಪ್ರತಿ ಶೆಲ್ ಅದೇ ಗುರುತುಗಳನ್ನು ನೋಡುತ್ತದೆ. ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸವನ್ನು ಹಂಚಿಕೊಳ್ಳುವ ಖಾತೆಗೆ, ccam ಮೊದಲು `history-<group>-<hash>` ನಲ್ಲಿ ಚಾಲನೆಯಲ್ಲಿರುವ ಸೆಷನ್‌ಗಾಗಿ ಹುಡುಕುತ್ತದೆ (ಹ್ಯಾಶ್ ಸಿಂಬಾಲಿಕ್ ಲಿಂಕ್‌ಗಳನ್ನು ಪರಿಹರಿಸಿದ ಈಗಿನ ಫೋಲ್ಡರ್‌ನದು), ಮತ್ತು ಅಲ್ಲಿಯೂ ನೋಂದಾಯಿಸುತ್ತದೆ (`historyCheck`).

   `exec` PID ಅನ್ನು ಉಳಿಸಿಕೊಳ್ಳುತ್ತದೆ, ಹೀಗಾಗಿ ಆ ಗುರುತು ನಂತರ Claude Code ಪ್ರಕ್ರಿಯೆಯನ್ನೇ ಸೂಚಿಸುತ್ತದೆ. ccam ಹಿಡಿದಿಟ್ಟ ಲಾಕ್ ಬದಲು ಗುರುತು ಫೈಲ್ ಅನ್ನು ಬಳಸುತ್ತದೆ, ಏಕೆಂದರೆ Claude Code ಆರಂಭಿಸುವ ಪ್ರಕ್ರಿಯೆಗಳು ಲಾಕ್ ಅನ್ನು ಪಡೆದು ಸೆಷನ್ ಮುಗಿದ ನಂತರವೂ ಅದನ್ನು ಇಟ್ಟುಕೊಳ್ಳುತ್ತವೆ. ಬರೆಯಲಾಗದ ಗುರುತು ಆರಂಭವನ್ನು ನಿಲ್ಲಿಸುವುದಿಲ್ಲ: ccam `ccam: could not mark the session as running (…); starting Claude Code anyway` ಎಂದು ತಿಳಿಸುತ್ತದೆ.

   ಖಾತೆ ಚಾಲನೆಯಲ್ಲಿದೆಯೇ ಎಂಬ ccam ನ ಇತರ ಪರಿಶೀಲನೆಗಳು (ಮರುಹೆಸರು, ತೆಗೆದುಹಾಕುವಿಕೆ, ಗುಂಪು ಬದಲಾವಣೆಗಳು, `doctor --fix`, MCP ಸಿಂಕ್) ಖಾತೆಯ ಫೋಲ್ಡರ್‌ನಲ್ಲಿರುವ Claude Code ನ ಸ್ವಂತ ದಾಖಲೆಗಳನ್ನೂ ಓದುತ್ತವೆ (`sessions/<pid>.json`, `daemon.lock`), ಹೀಗಾಗಿ ಅವು ccam ಇಲ್ಲದೆ ಆರಂಭಿಸಿದ Claude Code ಅನ್ನೂ ನೋಡುತ್ತವೆ: ಸಾಮಾನ್ಯ `claude`, IDE, ಹಿನ್ನೆಲೆ ಸೆಷನ್‌ಗಳು ಮತ್ತು ಅವುಗಳ daemon.
8. ccam ಟರ್ಮಿನಲ್ ಶೀರ್ಷಿಕೆಯನ್ನು `claude · <name>` ಗೆ ಹೊಂದಿಸುತ್ತದೆ (ಸ್ಟ್ಯಾಂಡರ್ಡ್ ಎರರ್ ಟರ್ಮಿನಲ್ ಆಗಿದ್ದಾಗ), ಆರಂಭದ ಸಮಯವನ್ನು `fsync` ಇಲ್ಲದೆ `~/.local/state/ccam/last-used/<name>` ಗೆ ಬರೆಯುತ್ತದೆ, ಮತ್ತು `syscall.Exec` ಅನ್ನು ಕರೆಯುತ್ತದೆ (`Plan.Exec`).

Claude Code `CLAUDE_CONFIG_DIR` ಅನ್ನು ಖಾತೆಯ ಫೋಲ್ಡರ್‌ಗೆ ಮತ್ತು `CCAM_PROFILE` ಅನ್ನು ಖಾತೆಯ ಹೆಸರಿಗೆ ಹೊಂದಿಸಿ ಚಾಲನೆಯಾಗುತ್ತದೆ. ಈ ವೇರಿಯಬಲ್‌ಗಳ ಯಾವುದೇ ಆನುವಂಶಿಕ ಮೌಲ್ಯಗಳನ್ನು ccam ತೆರವುಗೊಳಿಸುತ್ತದೆ, ಮತ್ತು ಆನುವಂಶಿಕವಾಗಿ ಬಂದ `CCAM_RELAY` ಅನ್ನೂ. macOS ನಲ್ಲಿ `ccam import` ಮೂಲಕ ನೋಂದಾಯಿಸಿದ `~/.claude` ಫೋಲ್ಡರ್‌ಗೆ, ಅದು `CLAUDE_SECURESTORAGE_CONFIG_DIR` ಅನ್ನೂ ಖಾಲಿ ಮೌಲ್ಯಕ್ಕೆ ಹೊಂದಿಸುತ್ತದೆ, ಹೀಗಾಗಿ ಆ ಖಾತೆ ಮತ್ತು ಸಾಮಾನ್ಯ `claude` ಒಂದೇ Keychain ಲಾಗಿನ್ ಬಳಸುತ್ತವೆ; Claude Code ಲಾಗಿನ್ ಅನ್ನು Credential Manager ನಲ್ಲಿ ಇಟ್ಟಾಗ Windows ನಲ್ಲೂ ಅದು ಇದನ್ನೇ ಮಾಡುತ್ತದೆ.

### exec ಏಕೆ

ccam ತನ್ನ ಸ್ವಂತ ಪ್ರಕ್ರಿಯೆಯನ್ನು `claude` ನಿಂದ ಬದಲಾಯಿಸುತ್ತದೆ (`syscall.Exec`). ಇದರ ಪರಿಣಾಮವಾಗಿ:

- ccam ಮೆಮೊರಿಯಲ್ಲಿ ಉಳಿಯುವುದಿಲ್ಲ;
- ಕೀಸ್ಟ್ರೋಕ್‌ಗಳು, ಸಿಗ್ನಲ್‌ಗಳು ಮತ್ತು ಎಕ್ಸಿಟ್ ಕೋಡ್ Claude Code ನದೇ ಆಗಿರುತ್ತವೆ;
- ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳು ಎಂದಿಗೂ ಶೆಲ್ ಸ್ಟ್ರಿಂಗ್ ಆಗದೆ ಅರೇ ಆಗಿಯೇ ಉಳಿಯುತ್ತವೆ.

Windows ನಲ್ಲಿ `exec` ಇಲ್ಲ. ಅಲ್ಲಿ, ccam ಅದೇ ಪರಿಸರ, ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳು ಮತ್ತು console ಜೊತೆ `claude.exe` ಅನ್ನು ಚೈಲ್ಡ್ ಪ್ರಕ್ರಿಯೆಯಾಗಿ ಆರಂಭಿಸುತ್ತದೆ, ಅದಕ್ಕಾಗಿ ಕಾಯುತ್ತದೆ ಮತ್ತು ಅದರ ಎಕ್ಸಿಟ್ ಕೋಡ್‌ನೊಂದಿಗೆ ನಿರ್ಗಮಿಸುತ್ತದೆ (`Plan.Exec`). npm ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿದ `claude.cmd`, ಅದರ ಪೂರ್ಣ ಪಾತ್ ಮತ್ತು ಉದ್ಧರಣ ಚಿಹ್ನೆಗಳಲ್ಲಿ ಇರಿಸಿದ ಪ್ರತಿ ಆರ್ಗ್ಯುಮೆಂಟ್‌ನೊಂದಿಗೆ `cmd.exe` ನಲ್ಲಿ ರನ್ ಆಗುತ್ತದೆ, ಹೀಗಾಗಿ `cmd.exe` ಅವುಗಳಲ್ಲಿ ಯಾವುದನ್ನೂ ತನ್ನದೇ ಸಿಂಟ್ಯಾಕ್ಸ್ (`&`, `%PATH%`) ಎಂದು ಓದುವುದಿಲ್ಲ; ಸಾಲು ವಿರಾಮ (line break) ಇರುವ ಆರ್ಗ್ಯುಮೆಂಟ್ ಅನ್ನು ಆ ರೀತಿ ರವಾನಿಸಲಾಗುವುದಿಲ್ಲ ಮತ್ತು ಅದನ್ನು ನಿರಾಕರಿಸಲಾಗುತ್ತದೆ, ನೇಟಿವ್ Claude Code ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡುವ ಸಲಹೆಯೊಂದಿಗೆ.

Ctrl+C console ನಿಂದ ನೇರವಾಗಿ Claude Code ಅನ್ನು ತಲುಪುತ್ತದೆ, ಆದ್ದರಿಂದ ccam ಅದನ್ನು ನಿರ್ಲಕ್ಷಿಸುತ್ತದೆ. ಸೆಷನ್ ಗುರುತು ccam ನ PID ಅನ್ನು ಹೊಂದಿರುತ್ತದೆ, ಅದು Claude Code ಇರುವಷ್ಟೇ ಕಾಲ ಜೀವಿಸುತ್ತದೆ, ಮತ್ತು Claude Code ನಿರ್ಗಮಿಸಿದಾಗ ಅದನ್ನು ತೆಗೆದುಹಾಕಲಾಗುತ್ತದೆ. ಕಾಯುತ್ತಿರುವಾಗ ccam ಬೇರೇನೂ ಮಾಡುವುದಿಲ್ಲ.

ಸೆಷನ್ ನಂತರದ ಕೆಲಸ, ಉದಾಹರಣೆಗೆ `claude mcp add` ಮೂಲಕ ಸೇರಿಸಿದ ಸರ್ವರ್‌ಗಳನ್ನು ಗಮನಿಸುವುದು, ಮುಂದಿನ ಆರಂಭದಲ್ಲಿ ನಡೆಯುತ್ತದೆ: ccam `.claude.json` ನ ಮಾರ್ಪಾಡು ಸಮಯವನ್ನು `~/.local/state/ccam/mcp/<name>/seen` ಜೊತೆ ಹೋಲಿಸುತ್ತದೆ ಮತ್ತು ಅದು ಬದಲಾದಾಗ ಮಾತ್ರ ಫೈಲ್ ಅನ್ನು ಪಾರ್ಸ್ ಮಾಡುತ್ತದೆ.

### ಕಮಾಂಡ್ ಲೈನ್ ಹೇಗೆ ಮಾರ್ಗ ಪಡೆಯುತ್ತದೆ

ಮೂರು ಕಮಾಂಡ್ ಲೈನ್‌ಗಳು, `internal/app` ಮೂಲಕ ಮೂರು ಮಾರ್ಗಗಳು:

<!-- motion:launch-route -->

ಶಾರ್ಟ್‌ಕಟ್ ಅನ್ನು `argv[0]` ಮೇಲೆ, `ccam <name>` ಅನ್ನು ಮೊದಲ ಆರ್ಗ್ಯುಮೆಂಟ್ ಮೇಲೆ ಹೊಂದಿಸಲಾಗುತ್ತದೆ, ಮತ್ತು ಎರಡೂ `syscall.Exec` ನಲ್ಲಿ ಕೊನೆಗೊಳ್ಳುತ್ತವೆ. `ccam usage` ನಂತಹ ccam ಕಮಾಂಡ್ ಮಾತ್ರ cobra ಅನ್ನು ತಲುಪುತ್ತದೆ.

### ಹಂಚಿಕೊಂಡ MCP ಸರ್ವರ್‌ಗಳು

ಡೀಫಾಲ್ಟ್ `flag` ತಂತ್ರದಲ್ಲಿ, ccam ಗುಂಪಿನ MCP ಸರ್ವರ್‌ಗಳನ್ನು ಒಂದು ಆರ್ಗ್ಯುಮೆಂಟ್ ಆಗಿ, ನಿಮ್ಮ ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳ ಮೊದಲು ಇರಿಸಿ, Claude Code ಗೆ ರವಾನಿಸುತ್ತದೆ:

<!-- motion:launch-mcp -->

ಗುಂಪಿನ `mcp.json` ಇದ್ದರೆ ಮಾತ್ರ ccam ಆರ್ಗ್ಯುಮೆಂಟ್ ಅನ್ನು ಸೇರಿಸುತ್ತದೆ, ಮತ್ತು ಅದನ್ನು `claude` ನ ತಕ್ಷಣ ನಂತರ ಇರಿಸುತ್ತದೆ (`mcp.FlagArgs`, `Plan.InsertArgs`). ಆಗ ನಿಮ್ಮದೇ `--mcp-config` ಅಥವಾ `--strict-mcp-config` ಅದರ ಮೇಲೆ ಅನ್ವಯವಾಗುತ್ತದೆ. ಇದು `--mcp-config=<path>` ಎಂಬ ಒಂದೇ ಆರ್ಗ್ಯುಮೆಂಟ್, ಏಕೆಂದರೆ Claude Code ನ `--mcp-config` ಹಲವು ಮೌಲ್ಯಗಳನ್ನು ತೆಗೆದುಕೊಳ್ಳುತ್ತದೆ: ಎರಡು ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳಾಗಿ ಬರೆದರೆ, ಅದು ನಂತರ ಬರುವ ಪ್ರಾಂಪ್ಟ್ ಅನ್ನೂ ಕಾನ್ಫಿಗ್ ಫೈಲ್ ಎಂದು ಓದುತ್ತದೆ.

ಫ್ಲ್ಯಾಗ್ ಎಲ್ಲಿ ಹೋಗಬಾರದು ಎಂಬುದನ್ನು ccam ನಿರ್ಧರಿಸುತ್ತದೆ (`claudeVerb`). ಅದು Claude Code ಓದುವ ರೀತಿಯಲ್ಲೇ, ಆಯ್ಕೆಗಳನ್ನು ಅವುಗಳ ಮೌಲ್ಯಗಳೊಂದಿಗೆ ಬಿಟ್ಟು, ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳಿಂದ Claude Code ನ ಉಪ-ಕಮಾಂಡ್ ಅನ್ನು ಓದುತ್ತದೆ. `remote-control` ಅಥವಾ `rc` ಗೆ ಫ್ಲ್ಯಾಗ್ ಸಿಗುವುದಿಲ್ಲ (ಕೆಳಗೆ), ಮೊದಲು ನೀಡಿದ ಹಿನ್ನೆಲೆ-ಸೆಷನ್ ಉಪ-ಕಮಾಂಡ್‌ಗೂ (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`) ಸಿಗುವುದಿಲ್ಲ. `ccam work mcp list`, `get` ಮತ್ತು `remove` ಗೂ ಸಿಗುವುದಿಲ್ಲ, ಏಕೆಂದರೆ Claude Code ಅಲ್ಲಿ ಅದನ್ನು ನಿರಾಕರಿಸುತ್ತದೆ; ಬದಲಿಗೆ ಗುಂಪಿನ ಸರ್ವರ್‌ಗಳು ಖಾತೆಯ ಸೆಷನ್‌ಗಳಲ್ಲೂ ಲೋಡ್ ಆಗುತ್ತವೆ ಎಂದು ccam ಸ್ಟ್ಯಾಂಡರ್ಡ್ ಎರರ್‌ನಲ್ಲಿ ತಿಳಿಸುತ್ತದೆ:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

ಉಳಿದ ಪ್ರತಿ ಉಪ-ಕಮಾಂಡ್ ಮತ್ತು ಪ್ರತಿ ಸೆಷನ್ ಫ್ಲ್ಯಾಗ್ ಅನ್ನು ಉಳಿಸಿಕೊಳ್ಳುತ್ತದೆ.

`config.toml` ನಲ್ಲಿ `mcp_strategy = "sync"` ಇದ್ದರೆ, ccam ಬದಲಿಗೆ ಸರ್ವರ್‌ಗಳನ್ನು ಖಾತೆಯ `.claude.json` ಗೆ ಬರೆಯುತ್ತದೆ, ಮತ್ತು ಖಾತೆ ಈಗಾಗಲೇ ಚಾಲನೆಯಲ್ಲಿದ್ದರೆ ಅದನ್ನು ಬಿಟ್ಟುಬಿಡುತ್ತದೆ. ಕೊನೆಯ ಪರಿಶೀಲನೆಯ ನಂತರ ಗುಂಪಿನ `mcp.json` ಆಗಲಿ `.claude.json` ಆಗಲಿ ಬದಲಾಗಿಲ್ಲದಿದ್ದಾಗಲೂ ಆರಂಭವು ಆ ಕೆಲಸವನ್ನು ಬಿಟ್ಟುಬಿಡುತ್ತದೆ. ಗುಂಪಿನ ಸರ್ವರ್‌ನ ಅದೇ ಹೆಸರಿನಲ್ಲಿ ಖಾತೆಗೆ ತನ್ನದೇ ಸರ್ವರ್ ಇದ್ದರೆ, ಸಿಂಕ್ ಖಾತೆಯ ಸರ್ವರ್ ಅನ್ನು ಉಳಿಸಿಕೊಂಡು ಅದನ್ನು ಒಮ್ಮೆ ತಿಳಿಸುತ್ತದೆ. ಸಿಂಕ್ `.claude.json` ನಿಂದ ಕಾಣೆಯಾದ ಪ್ರತಿಯನ್ನೂ ಮರಳಿ ಇರಿಸುತ್ತದೆ (`put back group default's MCP servers …`), ಮತ್ತು ತಾನು ಬರೆದು ಅಂದಿನಿಂದ ಯಾರೂ ಬದಲಾಯಿಸದ ಪ್ರತಿಗಳನ್ನು ಮಾತ್ರ ಬದಲಾಯಿಸುತ್ತದೆ ಅಥವಾ ತೆಗೆದುಹಾಕುತ್ತದೆ.

### Remote Control

`ccam work rc` (ಅಥವಾ `ccam work remote-control`) Claude Code ನ Remote Control ಅನ್ನು `work` ಆಗಿ ಆರಂಭಿಸುತ್ತದೆ. `claude remote-control` ತನ್ನ ಉಪ-ಕಮಾಂಡ್‌ಗಿಂತ ಮೊದಲು `--mcp-config` ಅನ್ನು ನಿರಾಕರಿಸುತ್ತದೆ ಮತ್ತು ನಂತರ ಯಾವುದನ್ನೂ ತೆಗೆದುಕೊಳ್ಳುವುದಿಲ್ಲ, ಮತ್ತು ಅದರ ಸೆಷನ್‌ಗಳು MCP ಸರ್ವರ್‌ಗಳನ್ನು ಖಾತೆಯ `.claude.json` ನಿಂದ ಮಾತ್ರ ಓದುತ್ತವೆ. ಹೀಗಾಗಿ ಇಲ್ಲಿ ಫ್ಲ್ಯಾಗ್ ವಿಧಾನ ಬೇರೆ ರೀತಿಯಲ್ಲಿ ಕೆಲಸ ಮಾಡುತ್ತದೆ:

<!-- motion:launch-rc -->

1. ccam `remote-control` ಅಥವಾ `rc` ಅನ್ನು Claude Code ನ ಉಪ-ಕಮಾಂಡ್ ಆಗಿ ಕಂಡುಹಿಡಿಯುತ್ತದೆ, ಆಯ್ಕೆಗಳು ಮತ್ತು ಅವುಗಳ ಮೌಲ್ಯಗಳ ನಂತರವೂ (`ccam work --model opus rc`), ಮತ್ತು ಯಾವುದೇ ಫ್ಲ್ಯಾಗ್ ಸೇರಿಸುವುದಿಲ್ಲ (`claudeVerb`, `mcp.FlagArgs`).
2. ccam ಗುಂಪಿನ `mcp.json` ನ ಮಾರ್ಪಾಡಿನ ಸಮಯವನ್ನು ಖಾತೆಯ `sync-stamp` ನೊಂದಿಗೆ ಹೋಲಿಸುತ್ತದೆ, ಮತ್ತು ಗುಂಪಿನ ಪ್ರತಿ ಸರ್ವರ್ ಇನ್ನೂ `.claude.json` ನಲ್ಲಿ ಇದೆಯೇ ಎಂದು ಪರಿಶೀಲಿಸುತ್ತದೆ (`mcp.Synced`). `.claude.json` ನಲ್ಲಿ ಈಗಾಗಲೇ ಗುಂಪಿನ ಪ್ರಸ್ತುತ ಸರ್ವರ್‌ಗಳಿದ್ದರೆ, ಅಥವಾ ಗುಂಪಿನಲ್ಲಿ ಯಾವುದೇ ಸರ್ವರ್ ಇಲ್ಲದಿದ್ದರೆ, Remote Control ಯಾವುದೇ ಪ್ರಶ್ನೆಯಿಲ್ಲದೆ ಆರಂಭವಾಗುತ್ತದೆ.
3. ಇಲ್ಲದಿದ್ದರೆ ಸೆಷನ್‌ಗಳಿಗೆ ಸಿಗದ ಸರ್ವರ್‌ಗಳ ಹೆಸರುಗಳನ್ನು ccam ತಿಳಿಸಿ ಹೀಗೆ ಕೇಳುತ್ತದೆ:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. `y` ನೀಡಿದರೆ, ccam ಅವುಗಳನ್ನು `sync` ವಿಧಾನದಂತೆಯೇ ನಕಲಿಸುತ್ತದೆ (`mcp.Sync`): Claude Code ನ ಲಾಕ್ ಅಡಿಯಲ್ಲಿ, ಬ್ಯಾಕಪ್ ನಂತರ, ಖಾತೆಯ ಸ್ವಂತ ಸರ್ವರ್‌ಗಳನ್ನು ಉಳಿಸಿಕೊಂಡು. `N` ನೀಡಿದರೆ, ಅಥವಾ ಟರ್ಮಿನಲ್ ಇಲ್ಲದಿದ್ದರೆ, Remote Control ಅವುಗಳಿಲ್ಲದೆ ಆರಂಭವಾಗುತ್ತದೆ, ಮತ್ತು ccam `mcp_strategy = "sync"` ಅನ್ನು ಸೂಚಿಸುತ್ತದೆ, ಅದು ಪ್ರತಿ ಆರಂಭದ ಮೊದಲು ಅವುಗಳನ್ನು ನಕಲಿಸುತ್ತದೆ.

ಸರ್ವರ್‌ಗಳು ನಕಲಾಗುವವರೆಗೆ ccam ಪ್ರತಿ Remote Control ಆರಂಭದಲ್ಲಿ ಕೇಳುತ್ತದೆ. ನಕಲಿನ ನಂತರ, ಗುಂಪಿನ `mcp.json` ಬದಲಾದಾಗ ಅಥವಾ ನಕಲಿಸಿದ ಸರ್ವರ್ `.claude.json` ನಿಂದ ಕಾಣೆಯಾದಾಗ ಮಾತ್ರ ಅದು ಮತ್ತೆ ಕೇಳುತ್ತದೆ. ನಂತರದ ಸಾಮಾನ್ಯ ಆರಂಭಗಳಲ್ಲಿ `.claude.json` ನಲ್ಲಿನ ನಕಲು `--mcp-config` ಮೇಲೆ ಗೆಲ್ಲುತ್ತದೆ, ಹೀಗಾಗಿ `ccam mcp add`, `remove` ಮತ್ತು `import` ಇತರ ಖಾತೆಗಳಲ್ಲಿರುವ ನಕಲುಗಳನ್ನು ನವೀಕರಿಸುತ್ತವೆ ([ಹಂಚಿಕೊಂಡ MCP ಸರ್ವರ್‌ಗಳು](/features/mcp)).

### 10 ms ಮಿತಿ

ಆರಂಭದ ಮಾರ್ಗದ ಮಿತಿ 95ನೇ percentile ನಲ್ಲಿ 10 ms; ಪ್ರಾಯೋಗಿಕವಾಗಿ Claude Code ಆರಂಭವಾಗುವ ಮೊದಲು ccam ಸುಮಾರು 4 ms ಸೇರಿಸುತ್ತದೆ. ಆರಂಭದ ಮಾರ್ಗ ಎಂದಿಗೂ cobra ಕಮಾಂಡ್ ಟ್ರೀ ಅನ್ನು ನಿರ್ಮಿಸುವುದಿಲ್ಲ, Bubble Tea ಅನ್ನು ಆರಂಭಿಸುವುದಿಲ್ಲ ಅಥವಾ ನೆಟ್‌ವರ್ಕ್ ಸಂಪರ್ಕ ತೆರೆಯುವುದಿಲ್ಲ, ಮತ್ತು ಪ್ಯಾಕೇಜ್ `init()` ಫಂಕ್ಷನ್‌ಗಳು ಒಟ್ಟು 5 ms ಗಿಂತ ಹೆಚ್ಚು ತೆಗೆದುಕೊಂಡರೆ `TestPackageInitIsCheap` ವಿಫಲವಾಗುತ್ತದೆ. `scripts/perf-check.sh` `hyperfine` ಮೂಲಕ `ccam --dry-run work`, `ccam work` ಮತ್ತು ಒಂದು ಶಾರ್ಟ್‌ಕಟ್ ಅನ್ನು ಅಳೆಯುತ್ತದೆ ಮತ್ತು 10 ms ಮೀರಿದರೆ ವಿಫಲವಾಗುತ್ತದೆ, ಮತ್ತು `internal/launcher` ನಲ್ಲಿರುವ `BenchmarkLaunchPath` ಹಂತಗಳನ್ನು ಲೂಪ್‌ನಲ್ಲಿ ನಡೆಸುತ್ತದೆ. ಪ್ರತಿ ಹಂತದ ಸಮಯವನ್ನು ನೋಡಲು `CCAM_TRACE=1` ಹೊಂದಿಸಿ:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## ಸುರಕ್ಷತೆ ಮತ್ತು ವಿಶೇಷ ಸಂದರ್ಭಗಳು

- ಅಪರಿಚಿತ ಹೆಸರು ಅಥವಾ ಸಂಖ್ಯೆ `no account "<name>"` ಎಂದು ನಿಲ್ಲುತ್ತದೆ ಮತ್ತು `ccam list` ಅನ್ನು ಸೂಚಿಸುತ್ತದೆ. ಯಾವುದೇ ಖಾತೆ ಇಲ್ಲದಿದ್ದರೆ, `ccam add <name>` ಮೂಲಕ ಒಂದನ್ನು ಸೇರಿಸಲು ccam ಕೇಳುತ್ತದೆ.
- ಹಂಚಿಕೊಂಡ ಪಾತ್ ಸಾಮಾನ್ಯ ಫೈಲ್ ಆಗಿದ್ದರೆ ಅಥವಾ ತಪ್ಪು ಸ್ಥಳಕ್ಕೆ ತೋರಿಸಿದರೆ, ccam `stopped because of share link problems:` ಎಂದು ನಿಂತು, ಅವುಗಳನ್ನು ಪಟ್ಟಿ ಮಾಡಿ, `ccam doctor --fix` ರನ್ ಮಾಡಲು ಕೇಳುತ್ತದೆ. ಆರಂಭದ ಮಾರ್ಗದಲ್ಲಿ ಅದು ಡೇಟಾವನ್ನು ಸ್ಥಳಾಂತರಿಸುವುದಿಲ್ಲ.
- ಈ ಖಾತೆಯೊಂದಿಗೆ ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸ ಹಂಚಿಕೊಳ್ಳುವ (ಅದೇ ಗುಂಪು) ಇನ್ನೊಂದು Claude Code ಸೆಷನ್ ಅದೇ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ ಚಾಲನೆಯಲ್ಲಿದ್ದರೆ, ಇತಿಹಾಸ ಬೆರೆಯಬಹುದು ಎಂದು ccam ಎಚ್ಚರಿಸಿ `Run anyway? [y/N]` ಎಂದು ಕೇಳುತ್ತದೆ. ಸಿಂಬಾಲಿಕ್ ಲಿಂಕ್‌ಗಳನ್ನು ಪರಿಹರಿಸಿ ಫೋಲ್ಡರ್ ಅನ್ನು ಹೊಂದಿಸಲಾಗುತ್ತದೆ. ಟರ್ಮಿನಲ್ ಇಲ್ಲದಿದ್ದರೆ (`/dev/null` ನಿಂದ ಬರುವ ಸ್ಟ್ಯಾಂಡರ್ಡ್ ಇನ್‌ಪುಟ್ ಅನ್ನು ಟರ್ಮಿನಲ್ ಇಲ್ಲ ಎಂದೇ ಪರಿಗಣಿಸಲಾಗುತ್ತದೆ), ಅಥವಾ ಹಿನ್ನೆಲೆಯ job ನಿಂದ, ಅದು ಎಚ್ಚರಿಕೆ ತೋರಿಸಿ ಮುಂದುವರಿಯುತ್ತದೆ. ಹೊಸ `-p` ರನ್ ತನ್ನದೇ ಸಂಭಾಷಣೆಯನ್ನು ಆರಂಭಿಸುತ್ತದೆ ಮತ್ತು ಯಾವುದೇ ಎಚ್ಚರಿಕೆ ಪಡೆಯುವುದಿಲ್ಲ; `--version`, `--help` ಮತ್ತು ಯಾವುದೇ ಸಂಭಾಷಣೆ ತೆರೆಯದ `mcp`, `auth` ಅಥವಾ `plugin` ನಂತಹ ಉಪ-ಕಮಾಂಡ್‌ಗಳನ್ನು ಪರಿಶೀಲಿಸುವುದೇ ಇಲ್ಲ.
- `--dry-run` MCP ಸಿಂಕ್, ಹೊಸ ಸರ್ವರ್ ಸೂಚನೆ, `mcp list` ಟಿಪ್ಪಣಿ, ಇತಿಹಾಸ ಪರಿಶೀಲನೆ, ಸೆಷನ್ ಗುರುತು, ಶೀರ್ಷಿಕೆ ಮತ್ತು ಕೊನೆಯ ಬಳಕೆಯ ಸಮಯವನ್ನು ಬಿಟ್ಟುಬಿಡುತ್ತದೆ, ಮತ್ತು exec ಮಾಡುವುದಿಲ್ಲ. ಆದರೂ ಅದು ಹಂಚಿಕೆ ಲಿಂಕ್‌ಗಳನ್ನು ಪರಿಶೀಲಿಸುತ್ತದೆ, ಮತ್ತು ಸಾಮಾನ್ಯ ಆರಂಭದಂತೆಯೇ ಇಲ್ಲದ ಲಿಂಕ್ ಅನ್ನು ರಚಿಸುತ್ತದೆ.
- `ccam work rc`, `work` ಈಗಾಗಲೇ ಚಾಲನೆಯಲ್ಲಿರುವಾಗ, `.claude.json` ಅನ್ನು ಮುಟ್ಟುವುದಿಲ್ಲ, ಅದನ್ನು ತಿಳಿಸುತ್ತದೆ, ಮತ್ತು ಗುಂಪಿನ ಸರ್ವರ್‌ಗಳಿಲ್ಲದೆ Remote Control ಅನ್ನು ಆರಂಭಿಸುತ್ತದೆ; ಅವುಗಳನ್ನು ನಕಲಿಸಲು ಸೆಷನ್ ಮುಗಿದ ನಂತರ ಮತ್ತೆ ರನ್ ಮಾಡಿ. `--dry-run` ಇದ್ದಾಗ ಅದು ಸಂದೇಶವನ್ನು ಮುದ್ರಿಸುತ್ತದೆ, ಆದರೆ ಕೇಳುವುದೂ ಇಲ್ಲ, ನಕಲಿಸುವುದೂ ಇಲ್ಲ.
- `mcp_strategy = "sync"` ಇದ್ದಾಗ, ವಿಫಲವಾದ ಸಿಂಕ್ `ccam: MCP sync failed (launching anyway): …` ಎಂದು ಮುದ್ರಿಸುತ್ತದೆ ಮತ್ತು Claude Code ಈಗಾಗಲೇ ಇದ್ದ ಸರ್ವರ್‌ಗಳೊಂದಿಗೆ ಆರಂಭವಾಗುತ್ತದೆ.

## ಸಂಬಂಧಿತ

- [ಶಾರ್ಟ್‌ಕಟ್ ಕಮಾಂಡ್‌ಗಳು](/features/shortcuts)
- [ccam ಜೊತೆಗೆ Remote Control ಬಳಸುವುದು](/guides/remote-control)
- [ಖಾತೆಗಳ ನಿರ್ವಹಣೆ](/features/accounts)
- [ಕಮಾಂಡ್ ಉಲ್ಲೇಖ](/docs/commands)
- [MCP ಸರ್ವರ್‌ಗಳು, ಸೆಟ್ಟಿಂಗ್‌ಗಳು ಮತ್ತು ಇತಿಹಾಸವನ್ನು ಹಂಚಿಕೊಳ್ಳಿ](/docs/sharing)
- [CLAUDE_CONFIG_DIR ವಿವರಣೆ](/guides/claude-config-dir)
