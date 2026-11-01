---
title: Claude Code-ஐத் தொடங்குதல்
description: ccam <name|number>, ccam run அல்லது டாஷ்போர்டு மூலம் எந்தக் கணக்காகவும் Claude Code-ஐத் தொடங்குங்கள்; arguments மாறாமல் claude-க்குச் செல்லும்.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: 11df08608a1e
---

`ccam work` Claude Code-ஐ `work` கணக்காகத் தொடங்குகிறது, பெயருக்குப் பின் உள்ள அனைத்தும் மாற்றமின்றி `claude`-க்குச் செல்லும். நீங்கள் `claude` என்று தட்டச்சு செய்யும் எந்த இடத்திலும் இதைப் பயன்படுத்துங்கள்: டெர்மினலில், script-இல், அல்லது pipe மூலம் உள்ளீடு கொடுக்கும்போது.

## பயன்பாடு

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

ஒரு கணக்கின் பெயரையோ `ccam list`-இல் உள்ள அதன் எண்ணையோ கொடுங்கள். பெயருக்குப் பின் உள்ள arguments-ஐ ccam படிப்பதில்லை, அதனால் `ccam work --help` Claude Code-இன் உதவியைக் காட்டும். பெயருக்குப் பின் `--` விருப்பத்தேர்வு மட்டுமே. standard input, exit code, signals ஆகியவை `claude`-இல் உள்ளது போலவே செயல்படும்.

முதல் argument `-` என்று தொடங்கி, ccam-இன் சொந்த flags-இல் ஒன்றாக இல்லாதபோது, ccam உங்களுக்காகக் கணக்கைத் தேர்ந்தெடுத்து, எல்லா arguments-ஐயும் Claude Code-க்கு அனுப்பும். அது [`ccam folder`](/features/folder) மூலம் தற்போதைய கோப்புறைக்கு அமைக்கப்பட்ட கணக்கைப் பயன்படுத்தும், இல்லையெனில் இயல்புநிலைக் கணக்கை: `ccam use` மூலம் அமைத்தது, அல்லது நீங்கள் ஒன்றை அமைக்கவில்லை என்றால் மிகக் குறைந்த எண் உள்ள கணக்கு.

### டாஷ்போர்டு

standard input-உம் output-உம் டெர்மினலாக இருக்கும்போது, arguments இல்லாத `ccam` [டாஷ்போர்டைத்](/features/dashboard) திறக்கும், இல்லையென்றால் `ccam list`-ஐ அச்சிடும். ஒரு கணக்குக்குச் சென்று (அம்புக்குறிகளால், அல்லது அதன் எண்ணைத் தட்டச்சு செய்து) அதைத் தொடங்க `Enter` அழுத்துங்கள்; `n` அதிக இடம் உள்ள கணக்கைத் தொடங்கும், `o` `--continue`, `--resume` அல்லது பிற arguments-ஐ வழங்கும். Claude Code தொடங்கும் முன் டாஷ்போர்டு மூடிக்கொள்கிறது, அதனால் டெர்மினல் முதலில் மீட்டமைக்கப்படும்.

**Team மற்றும் Enterprise இருக்கைகள்.** தானியங்கிக் கணக்குத் தேர்வு (`ccam next`, `--fallback`, `ccam pick`, டாஷ்போர்டின் `n`, நீங்கள் இயக்கிய பிறகு [`auto_switch`](/features/next#auto-switch)) ஒரே Claude Team அல்லது Enterprise நிறுவனத்தின் இருக்கைகளுக்கு இடையேயும் மாறுகிறது, இது Anthropic-இன் விதிமுறைகளை மீறக்கூடும். அதை நம்பியிருக்கும் முன், [அடிக்கடி கேட்கப்படும் கேள்விகளைப்](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) படித்து, உங்கள் நிறுவனத்திடமும் Anthropic-இடமும் கேட்டு உறுதிசெய்துகொள்ளுங்கள்.

shell-இல் `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` அல்லது `ANTHROPIC_AUTH_TOKEN` அமைக்கப்பட்டிருந்தால், Claude Code கணக்கின் சொந்த உள்நுழைவுக்குப் பதிலாக அந்த key அல்லது token-ஐப் பயன்படுத்தக்கூடும், அதனால் தொடங்கும் முன் ccam standard error-இல் உங்களை எச்சரிக்கிறது (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). அது variables-ஐப் படிக்க மட்டுமே செய்கிறது; ஒரு [API key கணக்கு](/features/api) அதற்குப் பதிலாக அவற்றை நீக்குகிறது.

இயல்பாக, ccam தன் இடத்தில் Claude Code-ஐ இயக்குகிறது, அதனால் வேறு எதுவும் தொடர்ந்து இயங்குவதில்லை. `config.toml`-இன் மேற்பகுதியில், முதல் `[[group]]`, `[[profile]]` அல்லது `[[folder]]`-க்கு மேலே, `relay = true` (அல்லது Linux, macOS-இல் [`auto_switch = true`](/features/next#auto-switch)) இருந்தால், ccam அதற்குப் பதிலாக interactive அமர்வுகளை ஒரு child process ஆகத் தொடங்கி, அவற்றுடன் சேர்ந்து இயங்கிக்கொண்டிருக்கும். அதனால் `/ccam:switch`, அல்லது பயன்பாட்டு வரம்பில் `auto_switch`, அமர்வின் options, permission mode, effort level, agent ஆகியவற்றுடன் உரையாடலை வேறொரு கணக்கில் தானாகவே தொடர முடியும்; `/ccam:switch` உரையாடலை வேறொரு நிறுவனத்தின் கணக்குக்கு நகர்த்தும் முன், ccam உங்களிடம் கேட்கிறது ([விவரங்கள்](/features/next#relay)).

## Options

| Option | என்ன செய்யும் | இயல்புநிலை |
| --- | --- | --- |
| `--dry-run` | ccam இயக்கும் சூழலையும் கட்டளையையும் அச்சிட்டு வெளியேறும் | off |
| `--fallback` | `-p`/`--print` இயக்கங்களுக்கு: கணக்கு தன் பயன்பாட்டு வரம்பை எட்டும்போது, அதே கட்டளையை அதிக இடம் உள்ள உள்நுழைந்த கணக்காக மீண்டும் இயக்கும் ([விவரங்கள்](/features/next#run-again-at-the-limit)) | off |
| `-v`, `--version` | `ccam <version>`-ஐ அச்சிடும் (முதல் argument ஆக மட்டும்) | |
| `-h`, `--help` | `ccam --help` அல்லது `ccam run --help` ஆக, ccam-இன் உதவி | |

Argument விதிகள்:

- ccam தன் சொந்த flags-ஐக் கணக்குப் பெயருக்கு முன் மட்டுமே படிக்கிறது. பெயருக்குப் பின், `--dry-run` மற்றும் `--help` உட்பட, ஒவ்வொரு argument-உம் Claude Code-க்கு உரியது.
- `--dry-run` `ccam --dry-run work`, `ccam run --dry-run work`, `ccam --dry-run` (கோப்புறையின் கணக்கு அல்லது இயல்புநிலைக் கணக்குக்கு) ஆக வேலை செய்யும். `--dry-run`, `--fallback` எந்த வரிசையிலும் வரலாம்.
- `ccam run <name|number> [claude args]` என்பது `ccam <name>`-இன் நீண்ட வடிவம். பெயர் இல்லாத `ccam run` கோப்புறையின் கணக்கை அல்லது இயல்புநிலைக் கணக்கைப் பயன்படுத்தும்.
- கணக்குப் பெயர்கள் `add`, `list`, `run`, `usage` போன்ற ccam கட்டளைப் பெயர்களாக இருக்க முடியாது, அதனால் `ccam <word>` ஒருபோதும் குழப்பமாக இருக்காது.

dry run, shell-இல், அல்லது Windows-இல் PowerShell-இல் ஒட்டக்கூடிய ஒரு வரியை அச்சிடும். தொடக்கம் சூழலில் மாற்றுவதுடன் அது தொடங்குகிறது: எப்போதும் `CLAUDE_CONFIG_DIR`, `CCAM_PROFILE`, அது நீக்கும் variable-க்கு `env -u <name>` (மரபுரிமையாகப் பெற்ற `CCAM_RELAY` போன்றது), ஒரு [API key கணக்குக்கு](/features/api) அதன் endpoint, அதன் model, key எங்கிருந்து வருகிறது என்பது (`ANTHROPIC_API_KEY="$WORK_KEY"`), key-ஐயே ஒருபோதும் அல்ல:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

shell-இல் `VAR=x cmd` செய்வது போலவே, PowerShell வரி பிறகு variables-ஐ முன்பிருந்தபடி மீட்டமைக்கிறது, Claude Code-க்கு அதன் கட்டளை வரியை `Start-Process` வழியாகக் கொடுக்கிறது, அதனால் Windows PowerShell 5.1-உம் PowerShell 7-உம் அதை ஒரே மாதிரி இயக்குகின்றன. ஒரு தொடக்கம் ஒரு variable-ஐ வெற்று மதிப்புக்கு அமைக்கும்போது, 7.5-க்கு முந்தைய PowerShell அத்தகைய variable-ஐ அதற்குப் பதிலாக நீக்கிவிடும் என்று அதன் கீழே ஒரு `#` வரி சொல்கிறது. npm நிறுவிய `claude.cmd`-க்கு, தொடக்கம் பயன்படுத்தும் `cmd.exe` கட்டளையை அந்த வரி இயக்குகிறது (கீழே பாருங்கள்), அங்கே தொடக்கம் மறுக்கும் argument dry run-ஐயும் தோல்வியடையச் செய்யும்.

## இது எப்படி வேலை செய்கிறது

ஒவ்வொரு தொடக்கமும், கட்டளை வரியிலிருந்து `exec` வரை, அதே படிகளை எடுக்கிறது:

<!-- motion:launch -->

ஒவ்வொரு படியும் இயங்கும்போது அதன் பெட்டி ஒளிரும், கீழே உள்ள வரிகள் Claude Code தொடங்கும் சூழலும் arguments-உம் ஆகும்.

1. ccam `app.Main`-இல் தொடங்குகிறது, அதை `cmd/ccam/main.go` அழைக்கிறது. ccam வேறொரு பெயரில் தொடங்கப்பட்டிருந்தால், அது ஒரு [குறுக்குவழி](/features/shortcuts), நேரடியாக launcher-க்குச் செல்லும்.
2. cobra ஈடுபடும் முன் ccam `os.Args`-ஐப் பார்க்கிறது (`app.parse`). தனியாக `ccam` (டாஷ்போர்டு), ஒரு ccam கட்டளைப் பெயர் (`config.Reserved`), `-h`/`--help`, `-v`/`--version` ஆகியவை மட்டுமே `cli.Execute`-க்குச் செல்லும். மற்ற அனைத்தும் தொடக்கமே.
3. ccam `~/.config/ccam/config.toml`-ஐப் படிக்கிறது; அது parse செய்யும் ஒரே கோப்பு இதுதான் (`launcher.Launch`).
4. ccam கணக்கைப் பெயராலோ எண்ணாலோ கண்டுபிடித்து, முன்னால் உள்ள ஒரு `--`-ஐ நீக்கி, `claude` executable-ஐக் கண்டறிகிறது (`launcher.Prepare`, `ResolveClaude`). பெயர் இல்லையென்றால், அது கோப்புறையின் கணக்கை எடுக்கிறது, working directory-ஐ `config.toml`-இன் `[[folder]]` பதிவுகளுடன் strings ஆக ஒப்பிட்டு, macOS-இலும் Windows-இலும் பெரிய, சிறிய எழுத்து வேறுபாட்டைப் பொருட்படுத்தாமல் (`HereProfile`). PATH தேடல் கோப்பின் மாற்ற நேரத்துடனும் `PATH`-இன் ஒரு hash உடனும் சேர்த்து `~/.local/state/ccam/claude-path`-இல் cache செய்யப்படுகிறது, அதனால் ஒவ்வொரு தொடக்கத்திலும் ccam `PATH`-ஐத் தேடுவதில்லை. மாறிய அல்லது இல்லாத கோப்பு, அல்லது மாறிய `PATH`, புதிய தேடலைத் தூண்டும். `config.toml`-இல் உள்ள `claude_bin` பாதை (`~/` விரிவாக்கப்படும்) இயக்கக்கூடிய கோப்பாக இருக்க வேண்டும்; இல்லையெனில் எதுவும் எழுதப்படும் முன் தொடக்கம் `claude_bin X is not an executable file. Check claude_bin in config.toml` என்று நின்றுவிடும். ccam-ஆகவே இருக்கும் `claude`-ஐ, inode மூலம் ஒப்பிட்டு, ccam மறுக்கிறது, ஏனெனில் அது முடிவில்லாச் சுழற்சியை ஏற்படுத்தும்.
5. ccam கணக்குக் கோப்புறை `~/.local/share/ccam/profiles/<name>/`-ஐ அதன் பகிர்வு அமைப்புகளுடன் ஒப்பிடுகிறது; அது `lstat` மற்றும் `readlink`-ஐ மட்டுமே அழைக்கிறது, ஒவ்வோர் உருப்படிக்கும் ஒன்று அல்லது இரண்டு system calls (`checkLinks`, `share.Check`). இல்லாத இணைப்பு அந்த இடத்திலேயே உருவாக்கப்படும். தரவை இழக்கக்கூடிய எதுவும் தொடக்கத்தை நிறுத்தும்.
6. ccam பகிரப்பட்ட MCP சர்வர்களைச் சேர்க்கிறது (`applyMCP`; கீழே பாருங்கள்).
7. ccam அமர்வைப் பதிவு செய்கிறது (`session.Register`). அது தன் PID-இன் பெயரில், process தொடக்க நேரத்தைக் கொண்ட ஒரு குறிப்பானை, `XDG_RUNTIME_DIR` என்ன சொன்னாலும், `~/.local/state/ccam/locks/profile-<name>/<pid>`-இல் எழுதுகிறது, அதனால் ஒவ்வொரு shell-உம் அதே குறிப்பான்களைப் பார்க்கும். உரையாடல் வரலாற்றைப் பகிரும் கணக்குக்கு, ccam முதலில் `history-<group>-<hash>`-இல் இயங்கும் அமர்வைத் தேடுகிறது, இங்கே hash, symbolic links தீர்க்கப்பட்ட தற்போதைய கோப்புறையினுடையது, அங்கேயும் பதிவு செய்கிறது (`historyCheck`).

   `exec` PID-ஐ வைத்திருக்கிறது, அதனால் குறிப்பான் பின்னர் Claude Code process-ஐயே பெயரிடுகிறது. ccam பிடித்து வைத்த lock-க்குப் பதிலாகக் குறிப்பான் கோப்பைப் பயன்படுத்துகிறது, ஏனெனில் Claude Code தொடங்கும் processes ஒரு lock-ஐப் பெற்று, அமர்வு முடிந்த பிறகும் அதை வைத்திருக்கும். எழுத முடியாத குறிப்பான் தொடக்கத்தை நிறுத்தாது: ccam `ccam: could not mark the session as running (…); starting Claude Code anyway` என்று சொல்கிறது.

   ஒரு கணக்கு இயங்குகிறதா என்பதற்கான ccam-இன் மற்ற சரிபார்ப்புகள் (மறுபெயரிடுதல், நீக்குதல், குழு மாற்றங்கள், `doctor --fix`, MCP sync) கணக்குக் கோப்புறையில் உள்ள Claude Code-இன் சொந்தப் பதிவுகளையும் (`sessions/<pid>.json`, `daemon.lock`) படிக்கின்றன, அதனால் ccam இல்லாமல் தொடங்கிய Claude Code-ஐயும் அவை பார்க்கின்றன: வெறும் `claude`, ஒரு IDE, பின்னணி அமர்வுகள், அவற்றின் daemon.
8. ccam டெர்மினல் தலைப்பை `claude · <name>` என அமைக்கிறது (standard error டெர்மினலாக இருக்கும்போது), தொடக்க நேரத்தை `fsync` இல்லாமல் `~/.local/state/ccam/last-used/<name>`-இல் எழுதுகிறது, `syscall.Exec`-ஐ அழைக்கிறது (`Plan.Exec`).

Claude Code, கணக்குக் கோப்புறைக்கு அமைக்கப்பட்ட `CLAUDE_CONFIG_DIR` உடனும் கணக்குப் பெயருக்கு அமைக்கப்பட்ட `CCAM_PROFILE` உடனும் இயங்குகிறது. இந்த variables-க்கு மரபுரிமையாகப் பெற்ற மதிப்புகள் எதையும், மரபுரிமையாகப் பெற்ற `CCAM_RELAY`-ஐயும் ccam அழிக்கிறது. macOS-இல் `ccam import` மூலம் பதிவுசெய்யப்பட்ட `~/.claude` கோப்புறைக்கு, அது `CLAUDE_SECURESTORAGE_CONFIG_DIR`-ஐ வெற்று மதிப்புக்கும் அமைக்கிறது, அதனால் அந்தக் கணக்கும் வெறும் `claude`-உம் ஒரே Keychain உள்நுழைவைப் பயன்படுத்தும்; Claude Code உள்நுழைவை Credential Manager-இல் வைத்திருக்கும்போது, Windows-இலும் அது இதையே செய்கிறது.

### ஏன் exec

ccam தன் சொந்த process-ஐ `claude`-ஆல் மாற்றுகிறது (`syscall.Exec`). இதன் விளைவாக:

- ccam நினைவகத்தில் தங்குவதில்லை;
- விசை அழுத்தங்கள், signals, exit code ஆகியவை Claude Code-இன் சொந்தமானவை;
- arguments ஒருபோதும் shell string ஆகாமல் array ஆகவே இருக்கும்.

Windows-இல் `exec` இல்லை. அங்கே ccam அதே சூழல், arguments, console-உடன் `claude.exe`-ஐ child process ஆகத் தொடங்கி, அதற்காகக் காத்திருந்து, அதன் exit code-உடன் வெளியேறுகிறது (`Plan.Exec`). npm நிறுவிய `claude.cmd`, அதன் முழுப் பாதையுடனும், `cmd.exe` எதையும் தன் சொந்த syntax ஆகப் படிக்காதபடி (`&`, `%PATH%`) ஒவ்வொரு argument-உம் மேற்கோள் குறிகளுக்குள் வைக்கப்பட்டும், `cmd.exe`-இல் இயங்குகிறது; வரி முறிவு உள்ள argument-ஐ அப்படி அனுப்ப முடியாது, அது மறுக்கப்படும், native Claude Code-ஐ நிறுவும் ஆலோசனையுடன்.

Ctrl+C console-இலிருந்து நேரடியாக Claude Code-ஐ அடைகிறது, அதனால் ccam அதைப் புறக்கணிக்கிறது. அமர்வு marker ccam-இன் PID-ஐக் கொண்டிருக்கிறது, அது சரியாக Claude Code இருக்கும் வரை இருக்கும், Claude Code வெளியேறும்போது marker நீக்கப்படுகிறது. காத்திருக்கும்போது ccam வேறு எதுவும் செய்வதில்லை.

அமர்வுக்குப் பிந்தைய வேலை, எடுத்துக்காட்டாக `claude mcp add` மூலம் சேர்க்கப்பட்ட சர்வர்களைக் கவனிப்பது, அடுத்த தொடக்கத்தில் நடக்கும்: ccam `.claude.json`-இன் மாற்ற நேரத்தை `~/.local/state/ccam/mcp/<name>/seen`-உடன் ஒப்பிட்டு, அது மாறியிருந்தால் மட்டுமே கோப்பை parse செய்கிறது.

### கட்டளை வரி எப்படி வழிநடத்தப்படுகிறது

மூன்று கட்டளை வரிகள், `internal/app` வழியாக மூன்று வழிகள்:

<!-- motion:launch-route -->

குறுக்குவழி `argv[0]`-இலும், `ccam <name>` முதல் argument-இலும் பொருத்தப்படுகின்றன, இரண்டும் `syscall.Exec`-இல் முடிகின்றன. `ccam usage` போன்ற ccam கட்டளை மட்டுமே cobra-வை அடைகிறது.

### பகிரப்பட்ட MCP சர்வர்கள்

இயல்புநிலை `flag` உத்தியுடன், ccam குழுவின் MCP சர்வர்களை உங்களுடையவற்றுக்கு முன் வைக்கப்படும் ஒரு argument ஆக Claude Code-க்கு அனுப்புகிறது:

<!-- motion:launch-mcp -->

குழுவின் `mcp.json` இருந்தால் மட்டுமே ccam argument-ஐச் சேர்க்கிறது, அதை `claude`-க்கு உடனே பின் வைக்கிறது (`mcp.FlagArgs`, `Plan.InsertArgs`). அப்போதும் உங்கள் சொந்த `--mcp-config` அல்லது `--strict-mcp-config` அதற்கு மேல் பொருந்தும். அது `--mcp-config=<path>` என்ற ஒரே argument, ஏனெனில் Claude Code-இன் `--mcp-config` பல மதிப்புகளை ஏற்கும்: இரண்டு arguments ஆக எழுதினால், அதைத் தொடரும் prompt-ஐயும் config கோப்பாகப் படித்துவிடும்.

flag எங்கே போக முடியாது என்பதை ccam தீர்மானிக்கிறது (`claudeVerb`). Claude Code செய்வது போலவே, options-ஐ அவற்றின் மதிப்புகளுடன் சேர்த்துத் தவிர்த்து, arguments-இலிருந்து Claude Code-இன் subcommand-ஐ அது படிக்கிறது. `remote-control` அல்லது `rc`-க்கு flag கிடைக்காது (கீழே), முதலில் கொடுக்கப்பட்ட பின்னணி அமர்வு verb-க்கும் (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`) கிடைக்காது. Claude Code அங்கே அதை மறுப்பதால், `ccam work mcp list`, `get`, `remove` ஆகியவற்றுக்கும் கிடைக்காது; அதற்குப் பதிலாக, குழுவின் சர்வர்கள் கணக்கின் அமர்வுகளிலும் ஏற்றப்படுகின்றன என்று ccam standard error-இல் சொல்கிறது:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

மற்ற ஒவ்வொரு subcommand-உம் ஒவ்வொரு அமர்வும் flag-ஐ வைத்திருக்கும்.

`config.toml`-இல் `mcp_strategy = "sync"` இருந்தால், ccam அதற்குப் பதிலாகச் சர்வர்களைக் கணக்கின் `.claude.json`-இல் எழுதும், கணக்கு ஏற்கெனவே இயங்கினால் அதைத் தவிர்க்கும். கடைசிச் சரிபார்ப்புக்குப் பிறகு குழுவின் `mcp.json`-உம் `.claude.json`-உம் மாறவில்லை என்றாலும் ஒரு தொடக்கம் இந்த வேலையைத் தவிர்க்கும். ஒரு குழுச் சர்வரின் அதே பெயரில் கணக்குக்குச் சொந்தமாக ஒரு சர்வர் இருந்தால், sync கணக்கின் சர்வரையே வைத்துக்கொண்டு அதை ஒருமுறை சொல்கிறது. `.claude.json`-இலிருந்து காணாமல் போன நகலையும் sync திரும்ப வைக்கிறது (`put back group default's MCP servers …`), தான் எழுதி, அதன் பிறகு யாரும் மாற்றாத நகல்களை மட்டுமே மாற்றுகிறது அல்லது நீக்குகிறது.

### Remote Control

`ccam work rc` (அல்லது `ccam work remote-control`) Claude Code-இன் Remote Control-ஐ `work` கணக்காகத் தொடங்குகிறது. `claude remote-control` தன் verb-க்கு முன் `--mcp-config`-ஐ மறுக்கிறது, அதற்குப் பின் எதையும் ஏற்பதில்லை; அதன் அமர்வுகள் MCP சர்வர்களைக் கணக்கின் `.claude.json`-இலிருந்து மட்டுமே படிக்கின்றன. அதனால் இங்கே flag உத்தி வேறு விதமாகச் செயல்படுகிறது:

<!-- motion:launch-rc -->

1. options-க்கும் அவற்றின் மதிப்புகளுக்கும் பிறகும் கூட (`ccam work --model opus rc`), ccam `remote-control` அல்லது `rc`-ஐ Claude Code-இன் verb ஆகக் கண்டுபிடித்து, flag எதையும் சேர்க்காது (`claudeVerb`, `mcp.FlagArgs`).
2. ccam குழுவின் `mcp.json` மாற்ற நேரத்தைக் கணக்கின் `sync-stamp`-உடன் ஒப்பிட்டு, ஒவ்வொரு குழுச் சர்வரும் இன்னும் `.claude.json`-இல் உள்ளதா என்று சரிபார்க்கிறது (`mcp.Synced`). `.claude.json`-இல் ஏற்கெனவே குழுவின் தற்போதைய சர்வர்கள் இருந்தால், அல்லது குழுவிடம் சர்வர்களே இல்லையென்றால், Remote Control கேள்வி இல்லாமல் தொடங்கும்.
3. இல்லையெனில், அமர்வுகளுக்குக் கிடைக்காமல் போகும் சர்வர்களை ccam பெயரிட்டுக் கேட்கும்:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. `y` என்றால், ccam அவற்றை `sync` உத்தி செய்வது போலவே நகலெடுக்கிறது (`mcp.Sync`): Claude Code-இன் lock-இன் கீழ், காப்புப்பிரதி எடுத்த பிறகு, கணக்கின் சொந்தச் சர்வர்களை வைத்துக்கொண்டு. `N` என்றால், அல்லது டெர்மினல் இல்லாவிட்டால், Remote Control அவை இல்லாமல் தொடங்கும், ccam ஒவ்வொரு தொடக்கத்துக்கு முன்பும் அவற்றை நகலெடுக்கும் `mcp_strategy = "sync"`-ஐச் சுட்டிக்காட்டும்.

சர்வர்கள் நகலெடுக்கப்படும் வரை ஒவ்வொரு Remote Control தொடக்கத்திலும் ccam கேட்கும். நகலெடுத்த பிறகு, குழுவின் `mcp.json` மாறும்போதோ நகலெடுத்த சர்வர் `.claude.json`-இலிருந்து காணாமல் போகும்போதோ மட்டுமே மீண்டும் கேட்கும். பின்வரும் சாதாரணத் தொடக்கங்களில் `.claude.json`-இல் உள்ள நகல் `--mcp-config`-ஐ விட முன்னுரிமை பெறும், அதனால் `ccam mcp add`, `remove`, `import` ஆகியவை மற்ற கணக்குகள் வைத்திருக்கும் நகல்களைப் புதுப்பிக்கின்றன ([பகிரப்பட்ட MCP சர்வர்கள்](/features/mcp)).

### 10 ms வரம்பு

தொடக்கப் பாதைக்கான வரம்பு 95-ஆவது percentile-இல் 10 ms; நடைமுறையில், Claude Code தொடங்கும் முன் ccam சுமார் 4 ms மட்டுமே சேர்க்கிறது. தொடக்கப் பாதை ஒருபோதும் cobra கட்டளை மரத்தை உருவாக்குவதோ, Bubble Tea-வைத் தொடங்குவதோ, network இணைப்பைத் திறப்பதோ இல்லை, package `init()` functions மொத்தம் 5 ms-க்கு மேல் எடுத்தால் `TestPackageInitIsCheap` தோல்வியடையும். `scripts/perf-check.sh` `ccam --dry-run work`, `ccam work`, ஒரு குறுக்குவழி ஆகியவற்றை `hyperfine` மூலம் அளந்து, 10 ms-க்கு மேல் இருந்தால் தோல்வியடைகிறது, `internal/launcher`-இல் உள்ள `BenchmarkLaunchPath` படிகளை ஒரு loop-இல் இயக்குகிறது. ஒவ்வொரு படியின் நேரத்தையும் பார்க்க, `CCAM_TRACE=1` அமையுங்கள்:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## பாதுகாப்பும் விளிம்பு நிலைகளும்

- தெரியாத பெயர் அல்லது எண் `no account "<name>"` என்று நின்று `ccam list`-ஐச் சுட்டிக்காட்டும். கணக்குகளே இல்லையென்றால், `ccam add <name>` மூலம் ஒன்றைச் சேர்க்க ccam கேட்கும்.
- பகிரப்பட்ட பாதை சாதாரணக் கோப்பாக மாறியிருந்தால் அல்லது தவறான இடத்தைச் சுட்டினால், ccam `stopped because of share link problems:` என்று நின்று, அவற்றைப் பட்டியலிட்டு, `ccam doctor --fix`-ஐ இயக்கச் சொல்லும். தொடக்கப் பாதையில் அது தரவை நகர்த்துவதில்லை.
- இந்தக் கணக்குடன் உரையாடல் வரலாற்றைப் பகிரும் (அதே குழு) மற்றொரு Claude Code அமர்வு அதே கோப்புறையில் இயங்கிக்கொண்டிருந்தால், வரலாறு கலக்கக்கூடும் என்று ccam எச்சரித்து, `Run anyway? [y/N]` என்று கேட்கும். symbolic links தீர்க்கப்பட்ட நிலையில் கோப்புறை பொருத்தப்படுகிறது. டெர்மினல் இல்லாவிட்டால் (`/dev/null`-இலிருந்து வரும் standard input டெர்மினல் இல்லாததாகவே கணக்கிடப்படும்), அல்லது பின்னணி job-இலிருந்து, எச்சரிக்கையைக் காட்டித் தொடரும். புதிய `-p` இயக்கம் தனக்கென ஓர் உரையாடலைத் தொடங்குவதால் எச்சரிக்கை பெறாது; `--version`, `--help`, `mcp`, `auth`, `plugin` போன்ற உரையாடலைத் திறக்காத subcommands சரிபார்க்கப்படவே மாட்டா.
- `--dry-run` MCP sync, புதிய சர்வர் அறிவிப்பு, `mcp list` குறிப்பு, வரலாற்றுச் சரிபார்ப்பு, அமர்வுக் குறிப்பான், தலைப்பு, கடைசிப் பயன்பாட்டு நேரம் ஆகியவற்றைத் தவிர்க்கும், exec செய்யாது. அது இன்னும் பகிர்வு இணைப்புகளைச் சரிபார்க்கும், சாதாரணத் தொடக்கம் செய்வது போலவே இல்லாத இணைப்பை உருவாக்கும்.
- `ccam work rc`, `work` ஏற்கெனவே இயங்கிக்கொண்டிருக்கும்போது, `.claude.json`-ஐத் தொடாமல் விட்டு, அதைச் சொல்லி, குழுவின் சர்வர்கள் இல்லாமல் Remote Control-ஐத் தொடங்கும்; அவற்றை நகலெடுக்க அமர்வு முடிந்த பிறகு மீண்டும் இயக்குங்கள். `--dry-run` உடன் அது செய்தியை அச்சிடும், ஆனால் கேட்கவோ நகலெடுக்கவோ செய்யாது.
- `mcp_strategy = "sync"` உடன், தோல்வியடைந்த sync `ccam: MCP sync failed (launching anyway): …` என்று அச்சிடும், Claude Code ஏற்கெனவே இருந்த சர்வர்களுடன் தொடங்கும்.

## தொடர்புடையவை

- [குறுக்குவழிக் கட்டளைகள்](/features/shortcuts)
- [ccam உடன் Remote Control-ஐப் பயன்படுத்துதல்](/guides/remote-control)
- [கணக்குகளை நிர்வகித்தல்](/features/accounts)
- [கட்டளைக் குறிப்பு](/docs/commands)
- [MCP சர்வர்கள், அமைப்புகள், வரலாற்றைப் பகிர்தல்](/docs/sharing)
- [CLAUDE_CONFIG_DIR விளக்கம்](/guides/claude-config-dir)
