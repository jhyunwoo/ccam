---
title: కమాండ్ రిఫరెన్స్
metaTitle: ccam కమాండ్ రిఫరెన్స్
description: ccam add, ccam <name> నుంచి ccam usage, ccam next, ccam doctor వరకు ccam యొక్క ప్రతి కమాండ్, ఆప్షన్.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: "28e1ad1dcbad"
---

ఒక కమాండ్ అసలు చేయకుండా అది ఏమి చేస్తుందో చూడటానికి, దాని ముందు `--dry-run` పెట్టండి, ఉదాహరణకు `ccam --dry-run work`. ఇది Claude Code ను ప్రారంభించేటప్పుడు (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`), అలాగే `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall`, `doctor --fix` తో పనిచేస్తుంది. `list`, `usage`, `share`, `ext` వంటి చదవడం మాత్రమే చేసే కమాండ్‌లు దాన్ని అంగీకరిస్తాయి, ఎలాగూ ఏమీ మార్చవు. మిగతావి, `group rename` తో సహా, రన్ అవ్వకుండా దాన్ని తిరస్కరిస్తాయి.

<!-- commands -->

## Claude Code ను ప్రారంభించడం

`ccam <name|number> [arguments]` ఆ ఖాతాగా Claude Code ను ప్రారంభిస్తుంది. పేరు తర్వాత ఉన్న ఆర్గ్యుమెంట్‌లను ccam చదవదు; వాటిని అదే క్రమంలో `claude` కు పంపుతుంది. స్టాండర్డ్ ఇన్‌పుట్, ఎగ్జిట్ కోడ్, సిగ్నల్స్ అన్నీ `claude` తో నేరుగా ఉన్నట్లే పనిచేస్తాయి. పేరు తర్వాత `--` ఐచ్ఛికం.

ఖాతా పేరు లేకుండా ఆర్గ్యుమెంట్‌లతో `ccam`, ఉదాహరణకు `ccam -p "summarise"`, ప్రస్తుత ఫోల్డర్‌కు సెట్ చేసిన ఖాతాను, లేకపోతే డిఫాల్ట్ ఖాతాను ఉపయోగిస్తుంది. వాటిని `ccam folder <name> [folder]`, `ccam use <name>` తో సెట్ చేయండి.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

ఖాతా పేరుకు ముందు `--fallback` అనేది `-p`/`--print` రన్‌ల కోసం. ఖాతా తన వినియోగ పరిమితిని చేరినప్పుడు, ccam అదే కమాండ్‌ను ఎక్కువ అవకాశం ఉన్న లాగిన్ అయిన ఖాతాగా మళ్లీ రన్ చేసి, పూర్తయిన రన్ output ను మాత్రమే ప్రింట్ చేస్తుంది. చూడండి [ఫోల్డర్ వారీగా ఖాతాలు](/features/folder), [పరిమితి వద్ద మళ్లీ రన్ చేయడం](/features/next#run-again-at-the-limit).

**Team, Enterprise సీట్లు.** ఆటోమేటిక్ ఖాతా ఎంపిక (`ccam next`, `--fallback`, `ccam pick`, డాష్‌బోర్డ్‌లోని `n`, మీరు ఆన్ చేసిన తర్వాత [`auto_switch`](/features/next#auto-switch)) ఒకే Claude Team లేదా Enterprise సంస్థలోని సీట్ల మధ్య కూడా మారుస్తుంది, ఇది Anthropic నిబంధనల ఉల్లంఘన కావచ్చు. దానిపై ఆధారపడే ముందు [తరచుగా అడిగే ప్రశ్నలు](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) చదివి, మీ సంస్థను, Anthropic ను సంప్రదించి నిర్ధారించుకోండి.

`ccwork` వంటి షార్ట్‌కట్ కూడా తన ఖాతాను ఇదే విధంగా ప్రారంభిస్తుంది: `ccam work --resume` చేసేదే `ccwork --resume` చేస్తుంది. కింద ఉన్న [షార్ట్‌కట్‌లు](/docs/commands#shortcuts) చూడండి.

## ఖాతాలను జోడించడం, తీసివేయడం

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

ఫోల్డర్ ఇవ్వకుండా `ccam import` రన్ చేస్తే, అది `~/.claude`, `~/.claude-*` కోసం వెతికి ప్రతిదాని గురించి అడుగుతుంది; విఫలమైన ఫోల్డర్ `Skipped <path>: <reason>` గా చూపబడుతుంది, మిగతావి నమోదు అవుతాయి. ఒక ఫోల్డర్ పేరు ఇప్పటికే వాడుకలో ఉంటే, ఆ ఫోల్డర్‌ను `ccam import <folder> --name <name>` తో విడిగా ఇంపోర్ట్ చేయండి. ఇంపోర్ట్ చేసిన ఫోల్డర్‌లు ఉన్న చోటే నమోదు అవుతాయి; ccam ఎప్పుడూ లాగిన్‌ను తరలించదు, కాపీ చేయదు.

Claude Code క్లియర్ చేసిన లాగిన్ ఉన్న ఖాతాలను కూడా `ccam login --expired` లాగిన్ చేస్తుంది. ccam నిర్వహించే ఖాతా (`ccam import` తో నమోదు చేసిన ఖాతా కాకుండా, ccam ఫోల్డర్ సృష్టించిన ఖాతా) ఫోల్డర్ ట్రాష్‌కు వెళ్లే ముందు `ccam remove` ఆ ఖాతాను లాగ్ అవుట్ చేస్తుంది, మరో ఖాతా అదే లాగిన్‌ను కలిగి ఉంటే తప్ప. ఇంపోర్ట్ చేసిన ఖాతాను తొలగిస్తే, అది తన గ్రూప్‌ల నుంచి బయటకు వెళ్తుంది, షేర్ చేసిన ప్రతి అంశం కాపీ దాని ఫోల్డర్‌లో ఉంటుంది. macOS లో `ccam rename` ఫోల్డర్ పేరును మార్చదు, ఎందుకంటే Keychain లాగిన్‌ను ఫోల్డర్ పాత్ కింద ఉంచుతుంది.

`ccam add --api` Claude లాగిన్ బదులు API కీ ఉపయోగించే ఖాతాను జోడిస్తుంది; ccam ప్రతి ప్రారంభంలో కీని చదువుతుంది, దాన్ని ఎప్పుడూ సేవ్ చేయదు. [API కీ ఖాతాలు](/features/api) చూడండి.

Claude Code ఒక ఖాతాగా రన్ అవుతున్నప్పుడు, అది ఎలా ప్రారంభమైనా సరే (బ్యాక్‌గ్రౌండ్ సెషన్‌లతో సహా), `ccam rename`, `ccam remove`, `ccam group` కమాండ్‌లు ఆ ఖాతాను మార్చడానికి నిరాకరిస్తాయి. ముందుగా ఆ సెషన్ నుంచి నిష్క్రమించండి. గ్రూప్ మార్పు ఒకటి అసంపూర్తిగా ఉన్నప్పుడు కూడా అవి నిరాకరిస్తాయి: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` ప్రతి ఖాతాను దాని స్లాట్, షార్ట్‌కట్, ఇమెయిల్, ఫోల్డర్‌తో ప్రింట్ చేస్తుంది. దానిలోని `groups` అర్రే ఖాతా గ్రూప్‌లను జాబితా చేస్తుంది, `share` ఆబ్జెక్ట్ ప్రతి అంశాన్ని (`mcp`, `settings`, `memory`, `extensions`, `history`) అది వచ్చే గ్రూప్‌కు, లేదా `local` కు మ్యాప్ చేస్తుంది.

## షార్ట్‌కట్‌లు

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

షార్ట్‌కట్ అంటే ఒక ఖాతాగా Claude Code ను ప్రారంభించే, మీరే పేరు పెట్టిన కమాండ్. `ccam work` లాగే, ఇది ప్రతి ఆర్గ్యుమెంట్‌ను మార్పు లేకుండా Claude Code కు పంపుతుంది, కాబట్టి `ccwork --resume`, `ccwork -p "hi"` మీరు ఆశించినట్లే పనిచేస్తాయి. ఇది `~/.local/bin` లో ccam కు ఉన్న లింక్.

షార్ట్‌కట్‌లు ఐచ్ఛికం. `ccam add` విజార్డ్ ఒకటి సృష్టించాలా అని అడుగుతుంది, `ccam add <name> --shortcut <command>` అడగకుండానే ఒకటి సృష్టిస్తుంది, `ccam list` ప్రతి ఖాతా షార్ట్‌కట్‌ను చూపిస్తుంది, లేకపోతే `-` చూపిస్తుంది.

పేర్లలో a-z, 0-9, `.`, `-`, `_` అక్షరాలు ఉపయోగించవచ్చు. పేరు ccam కమాండ్ కాకూడదు, మరో ఖాతా షార్ట్‌కట్ కాకూడదు, మీ `PATH` లోని మరో ప్రోగ్రామ్ పేరు కాకూడదు, ccam సృష్టించని `~/.local/bin` లోని ఫైల్ (ఉదాహరణకు ccam యొక్క మీ సొంత కాపీ లేదా లింక్) పేరు కూడా కాకూడదు. `config.toml` లో చేతితో ఎడిట్ చేసిన షార్ట్‌కట్ ఈ నియమాలను ఉల్లంఘిస్తే, `claude` లేదా `ccam` అయితే, లేదా రెండు ఖాతాలకు ఇచ్చి ఉంటే, `config.toml` లోడ్ అవ్వడం విఫలమవుతుంది, ఎర్రర్ ఆ ఖాతా పేరును చెబుతుంది.

`ccam rename` షార్ట్‌కట్ పేరును అలాగే ఉంచుతుంది, `ccam remove` ఖాతాతో పాటు షార్ట్‌కట్‌ను తొలగిస్తుంది. `ccam shims install` అన్ని ఖాతాల లింక్‌లను మళ్లీ సృష్టిస్తుంది, ఉదాహరణకు మీరు ccam ను వేరే చోటికి తరలించిన తర్వాత; ఒక షార్ట్‌కట్‌కు అడ్డుగా ఉన్న ఫైల్ మిగతా వాటిని ఆపదు.

## షేరింగ్

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

ఒక గ్రూప్ తన ఖాతాల మధ్య అంశాలను షేర్ చేస్తుంది; ఒక ఖాతా అనేక గ్రూప్‌లలో ఉండవచ్చు, కానీ ప్రతి అంశం ఒకే గ్రూప్ నుంచి వస్తుంది. [గ్రూప్‌లతో షేరింగ్](/features/share) చూడండి.

## ఇప్పటికే ఉన్న సెటప్‌ను కాపీ చేయడం

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` ఒక Claude Code ఫోల్డర్ నుంచి `settings.json`, `CLAUDE.md`, ఎక్స్‌టెన్షన్లు, యూజర్ MCP సర్వర్లను ఒక గ్రూప్‌లోకి కాపీ చేస్తుంది, కాబట్టి ఆ అంశాలను షేర్ చేసే ప్రతి ఖాతాకు అవి అందుతాయి. ఫోల్డర్ చదవబడుతుంది మాత్రమే.

- `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` మీరు పేర్కొన్న అంశాలను మాత్రమే కాపీ చేస్తాయి. అవి లేకపోతే, టెర్మినల్‌లో ఒక చెక్‌లిస్ట్ నుంచి మీరు ఎంచుకోవచ్చు, అందులో సంభాషణ చరిత్ర తప్ప అన్నీ టిక్ చేసి ఉంటాయి; ఇతర చోట్ల సంభాషణ చరిత్ర తప్ప అన్నీ కాపీ అవుతాయి.
- `~/.claude` ఖాళీగా ఉన్నా లేదా అసలు లేకపోయినా, `~/.claude.json` లో సర్వర్లు ఉంటే `--mcp` పనిచేస్తుంది.
- `-g` లేదా `--group` గ్రూప్‌ను ఎంచుకుంటుంది.
- `-y` లేదా `--yes` అడగకుండా కాపీ చేసి, ఫైల్ వేరుగా ఉన్నప్పుడు గ్రూప్ కాపీని ఉంచుతుంది.

MCP సైన్-ఇన్‌లు కాపీ అవ్వవు, ఎందుకంటే కాపీ చేసిన సైన్-ఇన్ దాన్ని కలిగి ఉన్న మిగతావాటిని లాగ్ అవుట్ చేస్తుంది; `/mcp` తో సైన్ ఇన్ చేయాల్సిన ఖాతాలను ప్లాన్ పేర్కొంటుంది. [ఇప్పటికే ఉన్న సెటప్‌ను కాపీ చేయడం](/docs/sharing#copy-an-existing-setup) చూడండి.

## MCP సర్వర్లు

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` ఆ ఖాతాకు మాత్రమే ఉన్న సర్వర్‌లను దాని గ్రూప్‌లోకి తరలిస్తుంది, ఒకే ఫోల్డర్ కోసం మాత్రమే ఉన్న సర్వర్‌లతో సహా (సాధారణ `claude mcp add`); అవి ఆ తర్వాత ప్రతి ఫోల్డర్‌లో లోడ్ అవుతాయి. గ్రూప్ ఇప్పటికే వేరే నిర్వచనం కోసం వాడుతున్న పేరు ఉన్న సర్వర్‌ను వదిలేస్తుంది; గ్రూప్ కాపీని భర్తీ చేయడానికి `--force` జోడించండి, ccam ముందుగా దాన్ని బ్యాకప్ చేస్తుంది.

## వినియోగం, ఖాతా మార్చడం

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

- `next --auto`, `--fallback`, `ccam pick` ఒకే విధంగా ఎంచుకుంటాయి: ముందుగా `auto_max_usage` (`config.toml`, డిఫాల్ట్ 100) కంటే తక్కువ ఉన్నవి, `auto_strategy` (`headroom`, `expiring` లేదా `balance`) ప్రకారం, `--model` వారపు విండోను కూడా లెక్కించి. అవి 100% వద్ద ఉన్న ఖాతాను (Claude Code నిండిందని నివేదించిన విండో దాని రీసెట్ వరకు 100% గా లెక్కించబడుతుంది) గానీ API కీ లాగిన్‌ను గానీ ఎప్పుడూ ఎంచుకోవు, అవి ప్రారంభించే ఖాతా సంస్థలోనే ఉంటాయి (ఇంకా దూరం చూడటానికి `--cross-org` లేదా `cross_org = true`).
- `--wait` తన విలువను `=` తర్వాత తీసుకుంటుంది (`--wait=2h`); ఖాళీతో రాస్తే ఒక సూచనతో తిరస్కరించబడుతుంది. ఏదీ ఎంచుకోలేనప్పుడు అవి కోడ్ 75 తో నిష్క్రమిస్తాయి, లేదా `--wait` తో మొదటి ఖాతా రీసెట్ అయ్యే వరకు వేచి ఉంటాయి.
- చరిత్రను షేర్ చేయని ఖాతాల మధ్య కూడా `ccam next --continue` సంభాషణను తీసుకువెళ్తుంది: అప్పుడు Claude Code దాన్ని కొత్త ఖాతాలోకి కాపీ చేస్తుంది. `--brief` తో, కొత్త ఖాతా బదులుగా కొత్త సంభాషణను ప్రారంభిస్తుంది, ccam మోడల్ లేకుండా రాసే చివరి సంభాషణ సారాంశం నుంచి ([`--brief`](/features/next#brief)).
- `ccam <name>` లా కాకుండా, `ccam next` తన సొంత ఆప్షన్‌లను చదువుతుంది, కాబట్టి Claude Code ఆప్షన్‌లు `--` తర్వాత వస్తాయి; అది లేకుండా ccam `unknown flag` తో ఆగిపోతుంది.
- `ccam agent` ఒక పనిని ఈ ఫోల్డర్‌లో `claude -p` తో, దాన్ని పిలిచిన ఖాతా కాకుండా వేరే ఖాతాగా, పిలిచిన సెషన్ permission mode, effort తో రన్ చేసి, Claude చివరి సమాధానాన్ని ప్రింట్ చేస్తుంది. ఆ ఖాతా పరిమితి వద్ద సంభాషణ మరో ఖాతాగా కొనసాగుతుంది ([ఇతర ఖాతాలపై subagents](/features/agent)).
- Claude Code లోపల, ప్లగిన్ యొక్క `/ccam:switch [account]` సంభాషణను ముందుకు అప్పగిస్తుంది, దాని hooks రన్ చేసేది `ccam hook`.
- `config.toml` లో `auto_switch = true` తో (Linux, macOS; డిఫాల్ట్‌గా ఆఫ్), పరిమితి hook ccam ప్రారంభించిన ఇంటరాక్టివ్ సెషన్‌ను కూడా ముగించి, అవకాశం ఉన్న మరో ఖాతాలో దాన్ని కొనసాగిస్తుంది, కానీ అదే రకమైన ఖాతాలో మాత్రమే: అదే సంస్థలోని మరో సీటు, లేదా మరో వ్యక్తిగత ప్లాన్ ([`auto_switch`](/features/next#auto-switch)).

[తదుపరి ఖాతాకు మారడం](/features/next), [స్క్రిప్ట్‌ల కోసం ఖాతాను ఎంచుకోవడం](/features/pick), [ఇతర ఖాతాలపై subagents](/features/agent), [ప్రతి ఖాతా బ్యాక్‌గ్రౌండ్ సెషన్‌లు](/features/agents) చూడండి.

## నిర్వహణ

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

`ccam update` ccam ను మార్చే ముందు రిలీజ్ SHA-256 చెక్‌సమ్‌ను తనిఖీ చేస్తుంది, `--version v0.2.0` ఒక నిర్దిష్ట రిలీజ్‌ను ఇన్‌స్టాల్ చేస్తుంది. `.deb` ప్యాకేజీ నుంచి ఇన్‌స్టాల్ చేసిన ccam కాపీని బదులుగా apt తో అప్‌డేట్ చేయండి. వివరాల కోసం [అప్‌డేట్ చేయడం](/docs/install#update) చూడండి.

`ccam uninstall` తాను తొలగించే వాటిని జాబితాగా చూపించి, ముందుగా అడుగుతుంది: `ccam remove` చేసే విధంగా ప్రతి ఖాతా, షార్ట్‌కట్ కమాండ్‌లు, ccam ఫోల్డర్‌లు, `ccam` కమాండ్. ఇంపోర్ట్ చేసిన ఫోల్డర్ ఉన్న చోటే ఉంటుంది. [అన్‌ఇన్‌స్టాల్ చేయడం](/docs/install#uninstall) చూడండి.

## ఫైళ్లు

| పాత్ | విషయాలు |
| --- | --- |
| `~/.config/ccam/config.toml` | ఖాతాలు, గ్రూప్‌లు, సెట్టింగ్‌లు (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). సెట్టింగ్‌లు పైభాగంలో, మొదటి `[[group]]`, `[[profile]]` లేదా `[[folder]]` కంటే పైన ఉండాలి: ఒక టేబుల్ కింద ఉన్న సెట్టింగ్ ఆ టేబుల్‌కే చెందుతుంది, కాబట్టి ccam ప్రతి లాంచ్‌లోనూ `ccam doctor` లోనూ దాని గురించి హెచ్చరిస్తుంది, దాన్ని పైకి తరలించే వరకు config ను సేవ్ చేసే కమాండ్‌లు నిరాకరిస్తాయి. ccam ఫైల్‌ను సేవ్ చేసినప్పుడు కామెంట్‌లు ఉంచబడవు |
| `~/.local/share/ccam/profiles/<name>/` | ప్రతి ఖాతా యొక్క `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | షేర్ చేసిన సెట్టింగ్‌లు, MCP సర్వర్లు, ఎక్స్‌టెన్షన్లు, చరిత్ర |
| `~/.local/state/ccam/` | వినియోగ కాష్, పరిమితి రికార్డులు, లాక్‌లు, నోట్స్, ఇతర స్థితి (టేబుల్ కింద ఉన్న జాబితా చూడండి) |
| `~/.local/bin/<shortcut>` | `ccwork` వంటి షార్ట్‌కట్ కమాండ్‌లు, ccam కు లింక్‌లు |

స్థితి ఫోల్డర్‌లో ఇవి ఉంటాయి:

- వినియోగ కాష్, పరిమితి రికార్డులు (`usage/<name>.json`, `usage/<name>.limit.json`);
- చివరిగా ఉపయోగించిన సమయాలు;
- సెషన్ మార్కర్లు, లాక్‌లు (`locks/`);
- ccam సృష్టించిన షార్ట్‌కట్‌లు (`shortcuts.json`);
- పూర్తి కాని గ్రూప్ మార్పు (`group-change.json`);
- `/ccam:switch`, పరిమితి hook వదిలే నోట్స్ (`handoff/<account>/<session>.json`);
- `ccam pick` లేదా `auto_switch` ప్రతి ఖాతాను చివరిగా ఎప్పుడు ఎంచుకుందో (`picks/<name>`; ఒక ఎంపిక 2 నిమిషాల పాటు లెక్కలో ఉంటుంది);
- 7 రోజులు ఉంచే `--brief` సారాంశాలు (`brief/<session>.md`).

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` లేదా `XDG_STATE_HOME` ను సెట్ చేస్తే ఈ ఫోల్డర్‌ల స్థానం మారుతుంది. Windows లో `~` అంటే `%USERPROFILE%`, షార్ట్‌కట్‌లు `<shortcut>.exe` గా ఉంటాయి.
