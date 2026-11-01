---
title: Claude Code ను ప్రారంభించడం
description: ccam <name|number>, ccam run లేదా డాష్‌బోర్డ్‌తో ఏ ఖాతాగానైనా Claude Code ప్రారంభించండి; ఆర్గ్యుమెంట్‌లు అలాగే claude కు వెళ్తాయి, --dry-run కమాండ్ చూపుతుంది.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: "11df08608a1e"
---

`ccam work` Claude Code ను `work` ఖాతాగా ప్రారంభిస్తుంది, పేరు తర్వాత ఉన్నదంతా మార్పు లేకుండా `claude` కు వెళ్తుంది. మీరు `claude` అని టైప్ చేసే ప్రతి చోటా దీన్ని ఉపయోగించండి: టెర్మినల్‌లో, స్క్రిప్ట్‌లో, లేదా పైప్ ద్వారా ఇన్‌పుట్ ఇచ్చి.

## వినియోగం

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

ఖాతా పేరును లేదా `ccam list` లోని దాని నంబర్‌ను ఇవ్వండి. పేరు తర్వాతి ఆర్గ్యుమెంట్‌లను ccam చదవదు, కాబట్టి `ccam work --help` Claude Code సహాయాన్ని చూపిస్తుంది. పేరు తర్వాత `--` ఐచ్ఛికం. స్టాండర్డ్ ఇన్‌పుట్, ఎగ్జిట్ కోడ్, సిగ్నల్స్ `claude` తో ఉన్నట్లే పనిచేస్తాయి.

మొదటి ఆర్గ్యుమెంట్ `-` తో మొదలై, ccam సొంత ఫ్లాగ్‌లలో ఒకటి కాకపోతే, ccam మీ కోసం ఖాతాను ఎంచుకుని, అన్ని ఆర్గ్యుమెంట్‌లను Claude Code కు పంపుతుంది. అది [`ccam folder`](/features/folder) తో ప్రస్తుత ఫోల్డర్‌కు సెట్ చేసిన ఖాతాను ఉపయోగిస్తుంది, లేకపోతే డిఫాల్ట్ ఖాతాను: `ccam use` తో సెట్ చేసినది, మీరు ఏదీ సెట్ చేయకపోతే అతి తక్కువ నంబర్ ఉన్న ఖాతా.

### డాష్‌బోర్డ్

స్టాండర్డ్ ఇన్‌పుట్, అవుట్‌పుట్ టెర్మినల్ అయినప్పుడు ఆర్గ్యుమెంట్‌లు లేని `ccam` [డాష్‌బోర్డ్](/features/dashboard) ను తెరుస్తుంది, లేకపోతే `ccam list` ను ప్రింట్ చేస్తుంది. ఒక ఖాతాకు వెళ్లి (బాణం కీలతో, లేదా దాని నంబర్ టైప్ చేసి) దాన్ని ప్రారంభించడానికి `Enter` నొక్కండి; `n` ఎక్కువ అవకాశం ఉన్న ఖాతాను ప్రారంభిస్తుంది, `o` `--continue`, `--resume` లేదా ఇతర ఆర్గ్యుమెంట్‌లను అందిస్తుంది. Claude Code ప్రారంభమయ్యే ముందు డాష్‌బోర్డ్ మూసుకుంటుంది, కాబట్టి టెర్మినల్ ముందుగా పునరుద్ధరించబడుతుంది.

**Team, Enterprise సీట్లు.** ఆటోమేటిక్ ఖాతా ఎంపిక (`ccam next`, `--fallback`, `ccam pick`, డాష్‌బోర్డ్‌లోని `n`, మీరు ఆన్ చేసిన తర్వాత [`auto_switch`](/features/next#auto-switch)) ఒకే Claude Team లేదా Enterprise సంస్థలోని సీట్ల మధ్య కూడా మారుస్తుంది, ఇది Anthropic నిబంధనల ఉల్లంఘన కావచ్చు. దానిపై ఆధారపడే ముందు [తరచుగా అడిగే ప్రశ్నలు](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) చదివి, మీ సంస్థను, Anthropic ను సంప్రదించి నిర్ధారించుకోండి.

షెల్‌లో `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` లేదా `ANTHROPIC_AUTH_TOKEN` సెట్ అయి ఉంటే, Claude Code ఖాతా సొంత లాగిన్ బదులు ఆ కీని లేదా టోకెన్‌ను ఉపయోగించవచ్చు, కాబట్టి ప్రారంభించే ముందు ccam స్టాండర్డ్ ఎర్రర్‌పై మిమ్మల్ని హెచ్చరిస్తుంది (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). అది వేరియబుల్స్‌ను చదువుతుంది మాత్రమే; [API కీ ఖాతా](/features/api) బదులుగా వాటిని తొలగిస్తుంది.

డిఫాల్ట్‌గా, ccam తన స్థానాన్ని Claude Code కు ఇస్తుంది, కాబట్టి ఇంకేదీ రన్ అవుతూ ఉండదు. `config.toml` పైభాగంలో, మొదటి `[[group]]`, `[[profile]]` లేదా `[[folder]]` కంటే పైన, `relay = true` (లేదా Linux, macOS లలో [`auto_switch = true`](/features/next#auto-switch)) తో, ccam బదులుగా ఇంటరాక్టివ్ సెషన్‌లను చైల్డ్ ప్రాసెస్‌గా ప్రారంభించి వాటితో పాటు రన్ అవుతూ ఉంటుంది. దాంతో `/ccam:switch`, లేదా వినియోగ పరిమితి వద్ద `auto_switch`, సెషన్ ఆప్షన్లు, permission mode, effort level, agent తో సహా సంభాషణను తనంతట తానే మరో ఖాతాలో కొనసాగించగలదు; `/ccam:switch` సంభాషణను మరో సంస్థ ఖాతాకు తరలించే ముందు, ccam మిమ్మల్ని అడుగుతుంది ([వివరాలు](/features/next#relay)).

## ఆప్షన్లు

| ఆప్షన్ | ఏమి చేస్తుంది | డిఫాల్ట్ |
| --- | --- | --- |
| `--dry-run` | ccam రన్ చేసే ఎన్విరాన్‌మెంట్, కమాండ్‌ను ప్రింట్ చేసి నిష్క్రమిస్తుంది | ఆఫ్ |
| `--fallback` | `-p`/`--print` రన్‌ల కోసం: ఖాతా తన వినియోగ పరిమితిని చేరినప్పుడు, అదే కమాండ్‌ను ఎక్కువ అవకాశం ఉన్న లాగిన్ అయిన ఖాతాగా మళ్లీ రన్ చేస్తుంది ([వివరాలు](/features/next#run-again-at-the-limit)) | ఆఫ్ |
| `-v`, `--version` | `ccam <version>` ను ప్రింట్ చేస్తుంది (మొదటి ఆర్గ్యుమెంట్‌గా మాత్రమే) | |
| `-h`, `--help` | `ccam --help` లేదా `ccam run --help` గా, ccam సహాయం | |

ఆర్గ్యుమెంట్ నియమాలు:

- ccam తన సొంత ఫ్లాగ్‌లను ఖాతా పేరుకు ముందు మాత్రమే చదువుతుంది. పేరు తర్వాత, `--dry-run`, `--help` తో సహా ప్రతి ఆర్గ్యుమెంట్ Claude Code కు చెందుతుంది.
- `--dry-run` `ccam --dry-run work`, `ccam run --dry-run work`, `ccam --dry-run` (ఫోల్డర్ ఖాతా లేదా డిఫాల్ట్ కోసం) గా పనిచేస్తుంది. `--dry-run`, `--fallback` ఏ క్రమంలోనైనా రావచ్చు.
- `ccam run <name|number> [claude args]` అనేది `ccam <name>` యొక్క పొడవైన రూపం. పేరు లేని `ccam run` ఫోల్డర్ ఖాతాను లేదా డిఫాల్ట్ ఖాతాను ఉపయోగిస్తుంది.
- ఖాతా పేర్లు `add`, `list`, `run` లేదా `usage` వంటి ccam కమాండ్ పేర్లు కాకూడదు, కాబట్టి `ccam <word>` ఎప్పుడూ అస్పష్టంగా ఉండదు.

డ్రై రన్ మీరు షెల్‌లో, లేదా Windows లో PowerShell లో పేస్ట్ చేయగల ఒక లైన్‌ను ప్రింట్ చేస్తుంది. అది ప్రారంభం ఎన్విరాన్‌మెంట్‌లో మార్చే వాటితో మొదలవుతుంది: ఎల్లప్పుడూ `CLAUDE_CONFIG_DIR`, `CCAM_PROFILE`, అది తొలగించే వేరియబుల్ కోసం `env -u <name>` (వారసత్వంగా వచ్చిన `CCAM_RELAY` వంటిది), [API కీ ఖాతా](/features/api) కోసం దాని ఎండ్‌పాయింట్, దాని మోడల్, కీ ఎక్కడి నుంచి వస్తుందో (`ANTHROPIC_API_KEY="$WORK_KEY"`), అసలు కీని ఎప్పుడూ కాదు:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

షెల్‌లో `VAR=x cmd` చేసినట్లే, PowerShell లైన్ తర్వాత వేరియబుల్స్‌ను వెనక్కి పెడుతుంది, `Start-Process` ద్వారా Claude Code కు దాని కమాండ్ లైన్‌ను అందిస్తుంది, కాబట్టి Windows PowerShell 5.1, PowerShell 7 దాన్ని ఒకే విధంగా రన్ చేస్తాయి. ప్రారంభం ఒక వేరియబుల్‌ను ఖాళీ విలువకు సెట్ చేసినప్పుడు, 7.5 కంటే ముందున్న PowerShell అలాంటి వేరియబుల్‌ను బదులుగా తొలగిస్తుందని దాని కింద ఒక `#` లైన్ చెబుతుంది. npm ఇన్‌స్టాల్ చేసిన `claude.cmd` కోసం, ఆ లైన్ ప్రారంభం ఉపయోగించే `cmd.exe` కమాండ్‌ను రన్ చేస్తుంది (కింద చూడండి), అక్కడ ప్రారంభం తిరస్కరించే ఆర్గ్యుమెంట్ డ్రై రన్‌ను కూడా విఫలం చేస్తుంది.

## ఇది ఎలా పనిచేస్తుంది

ప్రతి ప్రారంభం కమాండ్ లైన్ నుంచి `exec` వరకు అవే దశలను తీసుకుంటుంది:

<!-- motion:launch -->

ప్రతి బాక్స్ దాని దశ రన్ అయినప్పుడు వెలుగుతుంది, కింద ఉన్న లైన్లు Claude Code ప్రారంభమయ్యే ఎన్విరాన్‌మెంట్, ఆర్గ్యుమెంట్‌లు.

1. ccam `app.Main` లో ప్రారంభమవుతుంది, దాన్ని `cmd/ccam/main.go` పిలుస్తుంది. ccam మరో పేరుతో ప్రారంభించబడితే, అది ఒక [షార్ట్‌కట్](/features/shortcuts), నేరుగా లాంచర్‌కు వెళ్తుంది.
2. cobra జోక్యం చేసుకునే ముందే ccam `os.Args` ను చూస్తుంది (`app.parse`). ఒంటరి `ccam` (డాష్‌బోర్డ్), ccam కమాండ్ పేరు (`config.Reserved`), `-h`/`--help`, `-v`/`--version` మాత్రమే `cli.Execute` కు వెళ్తాయి. మిగతావన్నీ ప్రారంభమే.
3. ccam `~/.config/ccam/config.toml` ను చదువుతుంది, అది పార్స్ చేసే ఏకైక ఫైల్ అదే (`launcher.Launch`).
4. ccam పేరు లేదా నంబర్ ద్వారా ఖాతాను కనుగొని, ముందున్న ఒక `--` ను తొలగించి, `claude` ఎగ్జిక్యూటబుల్‌ను కనుగొంటుంది (`launcher.Prepare`, `ResolveClaude`). పేరు లేకపోతే, అది ఫోల్డర్ ఖాతాను తీసుకుంటుంది, working directory ని `config.toml` లోని `[[folder]]` ఎంట్రీలతో strings గా పోల్చి, macOS, Windows లలో పెద్ద చిన్న అక్షరాల తేడా లేకుండా (`HereProfile`). PATH లుకప్ ఫైల్ మార్పు సమయం, `PATH` హాష్‌తో పాటు `~/.local/state/ccam/claude-path` లో కాష్ అవుతుంది, కాబట్టి ccam ప్రతి ప్రారంభంలో `PATH` ను వెతకదు. ఫైల్ మారినా లేకపోయినా, లేదా `PATH` మారినా కొత్త లుకప్ జరుగుతుంది. `config.toml` లోని `claude_bin` పాత్ (`~/` విస్తరించబడుతుంది) ఎగ్జిక్యూటబుల్ ఫైల్ అయి ఉండాలి; లేకపోతే ఏదైనా రాయబడే ముందే ప్రారంభం `claude_bin X is not an executable file. Check claude_bin in config.toml` తో ఆగిపోతుంది. ccam తానే అయిన `claude` ను, inode ద్వారా పోల్చి, ccam తిరస్కరిస్తుంది, ఎందుకంటే అది లూప్ అవుతుంది.
5. ccam ఖాతా ఫోల్డర్ `~/.local/share/ccam/profiles/<name>/` ను దాని షేరింగ్ సెట్టింగ్‌లతో పోలుస్తుంది, `lstat`, `readlink` ను మాత్రమే పిలుస్తూ, ఒక్కో అంశానికి ఒకటి లేదా రెండు సిస్టమ్ కాల్స్ (`checkLinks`, `share.Check`). లేని లింక్ అక్కడికక్కడే సృష్టించబడుతుంది. డేటా కోల్పోయే అవకాశం ఉన్న ఏదైనా ప్రారంభాన్ని ఆపుతుంది.
6. ccam షేర్ చేసిన MCP సర్వర్లను జోడిస్తుంది (`applyMCP`; కింద చూడండి).
7. ccam సెషన్‌ను నమోదు చేస్తుంది (`session.Register`). అది తన PID పేరుతో, ప్రాసెస్ ప్రారంభ సమయాన్ని కలిగిన ఒక మార్కర్‌ను `~/.local/state/ccam/locks/profile-<name>/<pid>` లో రాస్తుంది, `XDG_RUNTIME_DIR` ఏమి చెప్పినా సరే, కాబట్టి ప్రతి షెల్ అవే మార్కర్లను చూస్తుంది. సంభాషణ చరిత్రను షేర్ చేసే ఖాతా కోసం, ccam ముందుగా `history-<group>-<hash>` లో రన్ అవుతున్న సెషన్ కోసం చూస్తుంది, ఇక్కడ హాష్ సింబాలిక్ లింక్‌లు పరిష్కరించిన ప్రస్తుత ఫోల్డర్‌ది, అక్కడ కూడా నమోదు చేస్తుంది (`historyCheck`).

   `exec` PID ను ఉంచుతుంది, కాబట్టి ఆ తర్వాత మార్కర్ Claude Code ప్రాసెస్‌నే సూచిస్తుంది. ccam పట్టుకున్న లాక్ కాకుండా మార్కర్ ఫైల్‌ను ఉపయోగిస్తుంది, ఎందుకంటే Claude Code ప్రారంభించే ప్రాసెస్‌లు లాక్‌ను వారసత్వంగా పొంది సెషన్ ముగిసిన తర్వాత కూడా దాన్ని ఉంచుకుంటాయి. రాయలేని మార్కర్ ప్రారంభాన్ని ఆపదు: ccam `ccam: could not mark the session as running (…); starting Claude Code anyway` అని చెబుతుంది.

   ఒక ఖాతా రన్ అవుతుందో లేదో చూసే ccam ఇతర తనిఖీలు (పేరు మార్పు, తీసివేత, గ్రూప్ మార్పులు, `doctor --fix`, MCP సింక్) ఖాతా ఫోల్డర్‌లోని Claude Code సొంత రికార్డులను కూడా చదువుతాయి (`sessions/<pid>.json`, `daemon.lock`), కాబట్టి అవి ccam లేకుండా ప్రారంభించిన Claude Code ను కూడా చూస్తాయి: సాధారణ `claude`, ఒక IDE, బ్యాక్‌గ్రౌండ్ సెషన్‌లు, వాటి daemon.
8. ccam టెర్మినల్ టైటిల్‌ను `claude · <name>` గా సెట్ చేస్తుంది (స్టాండర్డ్ ఎర్రర్ టెర్మినల్ అయినప్పుడు), `fsync` లేకుండా ప్రారంభ సమయాన్ని `~/.local/state/ccam/last-used/<name>` లో రాస్తుంది, `syscall.Exec` ను పిలుస్తుంది (`Plan.Exec`).

Claude Code `CLAUDE_CONFIG_DIR` ఖాతా ఫోల్డర్‌కు, `CCAM_PROFILE` ఖాతా పేరుకు సెట్ చేయబడి రన్ అవుతుంది. ccam ఈ వేరియబుల్స్‌కు వారసత్వంగా వచ్చిన విలువలను, వారసత్వంగా వచ్చిన `CCAM_RELAY` ను తొలగిస్తుంది. macOS లో `ccam import` తో నమోదు చేసిన `~/.claude` ఫోల్డర్ కోసం, ఇది `CLAUDE_SECURESTORAGE_CONFIG_DIR` ను కూడా ఖాళీ విలువకు సెట్ చేస్తుంది, తద్వారా ఆ ఖాతా, సాధారణ `claude` ఒకే Keychain లాగిన్‌ను ఉపయోగిస్తాయి; Claude Code లాగిన్‌ను Credential Manager లో ఉంచినప్పుడు Windows లో కూడా అది ఇలాగే చేస్తుంది.

### exec ఎందుకు

ccam తన సొంత ప్రాసెస్‌ను `claude` తో భర్తీ చేస్తుంది (`syscall.Exec`). దాని ఫలితంగా:

- ccam మెమరీలో ఉండదు;
- కీస్ట్రోక్‌లు, సిగ్నల్స్, ఎగ్జిట్ కోడ్ Claude Code సొంతమే;
- ఆర్గ్యుమెంట్‌లు ఎప్పుడూ షెల్ స్ట్రింగ్ కాకుండా అర్రేగానే ఉంటాయి.

Windows లో `exec` లేదు. అక్కడ ccam అదే ఎన్విరాన్‌మెంట్, ఆర్గ్యుమెంట్‌లు, console తో `claude.exe` ను చైల్డ్ ప్రాసెస్‌గా ప్రారంభించి, దాని కోసం వేచి ఉండి, దాని ఎగ్జిట్ కోడ్‌తో నిష్క్రమిస్తుంది (`Plan.Exec`). npm ఇన్‌స్టాల్ చేసిన `claude.cmd` `cmd.exe` లో రన్ అవుతుంది, దాని పూర్తి పాత్‌తో, ప్రతి ఆర్గ్యుమెంట్ కోట్ చేయబడి, కాబట్టి `cmd.exe` వాటిలో దేన్నీ తన సొంత సింటాక్స్‌గా (`&`, `%PATH%`) చదవదు; లైన్ బ్రేక్ ఉన్న ఆర్గ్యుమెంట్‌ను అలా పంపలేము, అది తిరస్కరించబడుతుంది, నేటివ్ Claude Code ను ఇన్‌స్టాల్ చేయమనే సలహాతో.

Ctrl+C console నుంచి నేరుగా Claude Code కు చేరుతుంది, కాబట్టి ccam దాన్ని పట్టించుకోదు. సెషన్ మార్కర్‌లో ccam PID ఉంటుంది, అది సరిగ్గా Claude Code ఉన్నంత కాలం ఉంటుంది, Claude Code నిష్క్రమించినప్పుడు మార్కర్ తొలగించబడుతుంది. వేచి ఉన్నప్పుడు ccam ఇంకేమీ చేయదు.

`claude mcp add` తో జోడించిన సర్వర్లను గమనించడం వంటి సెషన్ తర్వాతి పని, తదుపరి ప్రారంభంలో జరుగుతుంది: ccam `.claude.json` మార్పు సమయాన్ని `~/.local/state/ccam/mcp/<name>/seen` తో పోల్చి, ఫైల్ మారినప్పుడు మాత్రమే దాన్ని పార్స్ చేస్తుంది.

### కమాండ్ లైన్ ఎలా దారి మళ్లించబడుతుంది

మూడు కమాండ్ లైన్లు, `internal/app` ద్వారా మూడు మార్గాలు:

<!-- motion:launch-route -->

షార్ట్‌కట్ `argv[0]` ద్వారా, `ccam <name>` మొదటి ఆర్గ్యుమెంట్ ద్వారా గుర్తించబడతాయి, రెండూ `syscall.Exec` తో ముగుస్తాయి. `ccam usage` వంటి ccam కమాండ్ మాత్రమే cobra ను చేరుతుంది.

### షేర్ చేసిన MCP సర్వర్లు

డిఫాల్ట్ `flag` విధానంతో, ccam గ్రూప్ MCP సర్వర్లను మీ ఆర్గ్యుమెంట్‌ల ముందు, ఒక ఆర్గ్యుమెంట్‌గా Claude Code కు పంపుతుంది:

<!-- motion:launch-mcp -->

గ్రూప్ `mcp.json` ఉంటేనే ccam ఆర్గ్యుమెంట్‌ను జోడిస్తుంది, దాన్ని `claude` తర్వాత వెంటనే పెడుతుంది (`mcp.FlagArgs`, `Plan.InsertArgs`). ఆపై మీ సొంత `--mcp-config` లేదా `--strict-mcp-config` దానిపై అదనంగా వర్తిస్తుంది. ఇది ఒకే ఆర్గ్యుమెంట్, `--mcp-config=<path>`, ఎందుకంటే Claude Code యొక్క `--mcp-config` అనేక విలువలను తీసుకుంటుంది: రెండు ఆర్గ్యుమెంట్‌లుగా రాస్తే, తర్వాత వచ్చే ప్రాంప్ట్‌ను కూడా కాన్ఫిగ్ ఫైల్‌గా చదువుతుంది.

ఫ్లాగ్ ఎక్కడ వెళ్లకూడదో ccam నిర్ణయిస్తుంది (`claudeVerb`). అది Claude Code చేసే విధంగానే ఆర్గ్యుమెంట్‌ల నుంచి Claude Code సబ్‌కమాండ్‌ను చదువుతుంది, ఆప్షన్లను వాటి విలువలతో సహా దాటవేస్తూ. `remote-control` లేదా `rc` కు ఫ్లాగ్ రాదు (కింద చూడండి), ముందుగా ఇచ్చిన బ్యాక్‌గ్రౌండ్-సెషన్ వెర్బ్‌కు (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`) కూడా రాదు. `ccam work mcp list`, `get`, `remove` కు కూడా రాదు, ఎందుకంటే Claude Code అక్కడ దాన్ని తిరస్కరిస్తుంది; బదులుగా ఖాతా సెషన్‌లలో గ్రూప్ సర్వర్లు కూడా లోడ్ అవుతాయని ccam స్టాండర్డ్ ఎర్రర్‌పై చెబుతుంది:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

మిగతా ప్రతి సబ్‌కమాండ్, ప్రతి సెషన్ ఫ్లాగ్‌ను ఉంచుకుంటాయి.

`config.toml` లో `mcp_strategy = "sync"` తో, ccam బదులుగా సర్వర్లను ఖాతా `.claude.json` లో రాస్తుంది, ఖాతా ఇప్పటికే రన్ అవుతుంటే దాన్ని దాటవేస్తుంది. గత తనిఖీ నుంచి గ్రూప్ `mcp.json` గానీ `.claude.json` గానీ మారకపోతే కూడా ప్రారంభం ఆ పనిని దాటవేస్తుంది. ఒక గ్రూప్ సర్వర్ పేరుతోనే ఖాతాకు సొంత సర్వర్ ఉంటే, సింక్ ఖాతా సర్వర్‌నే ఉంచి, ఆ విషయం ఒక్కసారి చెబుతుంది. సింక్ `.claude.json` నుంచి పోయిన కాపీని కూడా తిరిగి పెడుతుంది (`put back group default's MCP servers …`), తాను రాసిన, ఆ తర్వాత ఎవరూ మార్చని కాపీలను మాత్రమే మారుస్తుంది లేదా తీసివేస్తుంది.

### Remote Control

`ccam work rc` (లేదా `ccam work remote-control`) Claude Code యొక్క Remote Control ను `work` గా ప్రారంభిస్తుంది. `claude remote-control` తన వెర్బ్‌కు ముందు `--mcp-config` ను తిరస్కరిస్తుంది, దాని తర్వాత ఏదీ తీసుకోదు, దాని సెషన్‌లు MCP సర్వర్లను ఖాతా `.claude.json` నుంచి మాత్రమే చదువుతాయి. కాబట్టి ఇక్కడ ఫ్లాగ్ విధానం వేరుగా పనిచేస్తుంది:

<!-- motion:launch-rc -->

1. ccam ఆప్షన్లు, వాటి విలువల తర్వాత ఉన్నా (`ccam work --model opus rc`) `remote-control` లేదా `rc` ను Claude Code వెర్బ్‌గా కనుగొని, ఏ ఫ్లాగ్‌నూ జోడించదు (`claudeVerb`, `mcp.FlagArgs`).
2. ccam గ్రూప్ `mcp.json` మార్పు సమయాన్ని ఖాతా `sync-stamp` తో పోల్చి, ప్రతి గ్రూప్ సర్వర్ ఇంకా `.claude.json` లో ఉందో లేదో తనిఖీ చేస్తుంది (`mcp.Synced`). `.claude.json` లో ఇప్పటికే గ్రూప్ ప్రస్తుత సర్వర్లు ఉంటే, లేదా గ్రూప్‌కు సర్వర్లే లేకపోతే, Remote Control ప్రశ్న లేకుండా ప్రారంభమవుతుంది.
3. లేకపోతే ccam సెషన్‌లకు అందని సర్వర్ల పేర్లు చెప్పి అడుగుతుంది:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. `y` అయితే, ccam వాటిని `sync` విధానం లాగే కాపీ చేస్తుంది (`mcp.Sync`): Claude Code లాక్ కింద, బ్యాకప్ తర్వాత, ఖాతా సొంత సర్వర్లను ఉంచుతూ. `N` అయితే, లేదా టెర్మినల్ లేకపోతే, Remote Control అవి లేకుండా ప్రారంభమవుతుంది, ccam ప్రతి ప్రారంభానికి ముందు వాటిని కాపీ చేసే `mcp_strategy = "sync"` ను సూచిస్తుంది.

సర్వర్లు కాపీ అయ్యే వరకు ప్రతి Remote Control ప్రారంభంలో ccam అడుగుతుంది. కాపీ తర్వాత, గ్రూప్ `mcp.json` మారినప్పుడు లేదా కాపీ చేసిన సర్వర్ `.claude.json` నుంచి మాయమైనప్పుడు మాత్రమే మళ్లీ అడుగుతుంది. తర్వాతి సాధారణ ప్రారంభాల్లో `.claude.json` లోని కాపీ `--mcp-config` పై గెలుస్తుంది, కాబట్టి `ccam mcp add`, `remove`, `import` ఇతర ఖాతాలు ఉంచుకున్న కాపీలను అప్‌డేట్ చేస్తాయి ([షేర్ చేసిన MCP సర్వర్లు](/features/mcp)).

### 10 ms బడ్జెట్

ప్రారంభ మార్గానికి బడ్జెట్ 95వ పర్సెంటైల్ వద్ద 10 ms; ఆచరణలో Claude Code ప్రారంభమయ్యే ముందు ccam సుమారు 4 ms జోడిస్తుంది. ప్రారంభ మార్గం ఎప్పుడూ cobra కమాండ్ ట్రీని నిర్మించదు, Bubble Tea ను ప్రారంభించదు, నెట్‌వర్క్ కనెక్షన్‌ను తెరవదు, ప్యాకేజీ `init()` ఫంక్షన్‌లు మొత్తం 5 ms కంటే ఎక్కువ తీసుకుంటే `TestPackageInitIsCheap` విఫలమవుతుంది. `scripts/perf-check.sh` `hyperfine` తో `ccam --dry-run work`, `ccam work`, ఒక షార్ట్‌కట్‌ను కొలిచి 10 ms దాటితే విఫలమవుతుంది, `internal/launcher` లోని `BenchmarkLaunchPath` దశలను ఒక లూప్‌లో రన్ చేస్తుంది. ఒక్కో దశ సమయాన్ని చూడటానికి, `CCAM_TRACE=1` ను సెట్ చేయండి:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## భద్రత, ప్రత్యేక సందర్భాలు

- తెలియని పేరు లేదా నంబర్ `no account "<name>"` తో ఆగి `ccam list` ను సూచిస్తుంది. అసలు ఖాతాలే లేకపోతే, ccam `ccam add <name>` తో ఒకటి జోడించమని అడుగుతుంది.
- ఒక షేర్డ్ పాత్ సాధారణ ఫైల్‌గా మారితే లేదా తప్పు చోటికి చూపిస్తే, ccam `stopped because of share link problems:` తో ఆగి, వాటిని జాబితా చేసి, `ccam doctor --fix` రన్ చేయమని అడుగుతుంది. ప్రారంభ మార్గంలో అది డేటాను తరలించదు.
- ఈ ఖాతాతో సంభాషణ చరిత్రను షేర్ చేసే (అదే గ్రూప్) మరో Claude Code సెషన్ అదే ఫోల్డర్‌లో రన్ అవుతుంటే, చరిత్ర కలిసిపోవచ్చని ccam హెచ్చరించి `Run anyway? [y/N]` అని అడుగుతుంది. ఫోల్డర్ సింబాలిక్ లింక్‌లు పరిష్కరించి సరిపోల్చబడుతుంది. టెర్మినల్ లేకపోతే (`/dev/null` నుంచి వచ్చే స్టాండర్డ్ ఇన్‌పుట్ టెర్మినల్ లేనట్టే), లేదా బ్యాక్‌గ్రౌండ్ job నుంచి అయితే, అది హెచ్చరికను చూపించి కొనసాగుతుంది. కొత్త `-p` రన్ తన సొంత సంభాషణను మొదలుపెడుతుంది, దానికి హెచ్చరిక రాదు; `--version`, `--help`, `mcp`, `auth` లేదా `plugin` వంటి సంభాషణను తెరవని సబ్‌కమాండ్‌లు అసలు తనిఖీ చేయబడవు.
- `--dry-run` MCP సింక్, కొత్త సర్వర్ నోటీసు, `mcp list` నోట్, చరిత్ర తనిఖీ, సెషన్ మార్కర్, టైటిల్, చివరి వినియోగ సమయాన్ని దాటవేస్తుంది, exec చేయదు. అయినా అది షేర్ లింక్‌లను తనిఖీ చేస్తుంది, సాధారణ ప్రారంభం లాగే లేని లింక్‌ను సృష్టిస్తుంది.
- `ccam work rc`, `work` ఇప్పటికే రన్ అవుతున్నప్పుడు, `.claude.json` ను ముట్టుకోకుండా, ఆ విషయం చెప్పి, గ్రూప్ సర్వర్లు లేకుండా Remote Control ను ప్రారంభిస్తుంది; వాటిని కాపీ చేయడానికి సెషన్ ముగిసిన తర్వాత మళ్లీ రన్ చేయండి. `--dry-run` తో అది సందేశాన్ని ప్రింట్ చేస్తుంది కానీ అడగదు, కాపీ చేయదు.
- `mcp_strategy = "sync"` తో, సింక్ విఫలమైతే `ccam: MCP sync failed (launching anyway): …` అని ప్రింట్ అవుతుంది, Claude Code తనకు ఇప్పటికే ఉన్న సర్వర్లతో ప్రారంభమవుతుంది.

## సంబంధిత పేజీలు

- [షార్ట్‌కట్ కమాండ్‌లు](/features/shortcuts)
- [ccam తో Remote Control ఉపయోగించడం](/guides/remote-control)
- [ఖాతాలను నిర్వహించడం](/features/accounts)
- [కమాండ్ రిఫరెన్స్](/docs/commands)
- [MCP సర్వర్లు, సెట్టింగ్‌లు, చరిత్రను షేర్ చేయండి](/docs/sharing)
- [CLAUDE_CONFIG_DIR వివరణ](/guides/claude-config-dir)
