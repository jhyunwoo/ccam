---
title: తదుపరి ఖాతాకు మారడం
description: ccam next Claude Code ను మరో ఖాతాగా తెరుస్తుంది; --auto ఎక్కువ అవకాశం ఉన్నది, --continue సంభాషణను తీసుకువెళ్తుంది, --fallback -p jobs మళ్లీ నడుపుతుంది.
updated: 2026-11-05
order: 8
translation: machine
sourceHash: "de1fd573f1ac"
---

`ccam next` మీరు పేరు చెప్పకుండానే Claude Code ను మరో ఖాతాగా తెరుస్తుంది. ఒక ఖాతా తన వినియోగ పరిమితిని చేరుకున్నప్పుడు దీన్ని ఉపయోగించండి. మారడానికి నాలుగు మార్గాలు ఉన్నాయి:

- సాధారణ `ccam next` మీరు చివరిగా ఉపయోగించిన ఖాతా తర్వాతి ఖాతాను తీసుకుంటుంది, `--auto` ఎక్కువ అవకాశం ఉన్న ఖాతాను తీసుకుంటుంది. `--continue` ఈ ఫోల్డర్ యొక్క చివరి సంభాషణను కొత్త ఖాతాలో కొనసాగిస్తుంది, ఆ రెండూ చరిత్రను షేర్ చేసినా చేయకపోయినా.
- Claude Code లోపల, ఒక పరిమితి సెషన్‌ను ఆపినప్పుడు [ccam ప్లగిన్](#inside-claude-code) మీకు చెబుతుంది, `/ccam:switch` సంభాషణను ముందుకు అప్పగిస్తుంది.
- `claude -p` రన్ చేసే స్క్రిప్ట్‌ల కోసం, `ccam --fallback` ఆ మార్పును తనంతట తానే చేస్తుంది (చూడండి [పరిమితి వద్ద మళ్లీ రన్ చేయడం](#run-again-at-the-limit)).
- [`auto_switch`](#auto-switch) ఆన్ చేసి ఉంటే (Linux, macOS), Claude Code తన పరిమితి వద్ద ఆపిన ఇంటరాక్టివ్ సెషన్‌ను కూడా ccam తరలిస్తుంది.

**Team, Enterprise సీట్లు.** ఆటోమేటిక్ ఖాతా ఎంపిక (`ccam next`, `--fallback`, `ccam pick`, డాష్‌బోర్డ్‌లోని `n`, మీరు ఆన్ చేసిన తర్వాత [`auto_switch`](/features/next#auto-switch)) ఒకే Claude Team లేదా Enterprise సంస్థలోని సీట్ల మధ్య కూడా మారుస్తుంది, ఇది Anthropic నిబంధనల ఉల్లంఘన కావచ్చు. దానిపై ఆధారపడే ముందు [తరచుగా అడిగే ప్రశ్నలు](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) చదివి, మీ సంస్థను, Anthropic ను సంప్రదించి నిర్ధారించుకోండి.

## వినియోగం

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

ccam మార్పును స్టాండర్డ్ ఎర్రర్‌లో `→ <from> → <to>` గా, `--auto` ఎంపికకు కారణంతో సహా ప్రింట్ చేసి, ఆపై `ccam <name>` చేసే విధంగానే కొత్త ఖాతాగా Claude Code ను ప్రారంభిస్తుంది. `<from>` అంటే మీరు చివరిగా ఉపయోగించిన ఖాతా, లేదా `--continue` తో సంభాషణ చెందిన ఖాతా. `--` తర్వాతి ఆర్గ్యుమెంట్‌లు మార్పు లేకుండా `claude` కు వెళ్తాయి.

## ఆప్షన్లు

| ఆప్షన్ | ఏమి చేస్తుంది | డిఫాల్ట్ |
| --- | --- | --- |
| `--auto` | వినియోగం ఆధారంగా (కింద చూడండి) లాగిన్ అయిన ఖాతాను ఎంచుకుంటుంది, 100% వద్ద ఉన్న ఖాతాలను వదిలేస్తుంది | ఆఫ్ |
| `--max-usage <1-100>` | `--auto` తో: ఈ వినియోగం కంటే తక్కువ ఉన్న ఖాతాలకు ప్రాధాన్యం ఇస్తుంది | `auto_max_usage`, లేకపోతే 100 |
| `--strategy headroom\|expiring\|balance` | `--auto` తో: అతి తక్కువ 5 గంటల వినియోగం, అన్నిటికంటే ముందు వచ్చే 5 గంటల రీసెట్, లేదా రోజుకు ఎక్కువగా మిగిలి ఉన్న వారపు కోటా | `auto_strategy`, లేకపోతే `headroom` |
| `--cross-org` | `--auto` తో: మరో సంస్థకు చెందిన ఖాతాలను కూడా ఎంచుకుంటుంది (కింద చూడండి) | `cross_org`, లేకపోతే ఆఫ్ |
| `--wait[=<duration>]` | `--auto` తో: ఏ ఖాతాలోనూ అవకాశం లేనప్పుడు, మొదటి ఖాతా రీసెట్ అయ్యే వరకు, గరిష్ఠంగా ఇంత సేపు వేచి ఉంటుంది. విలువ `=` తర్వాత వస్తుంది: `--wait 2h` `--wait takes its value after '=': --wait=2h` అనే సూచనతో తిరస్కరించబడుతుంది | ఆఫ్; ఒంటరిగా ఇస్తే `5h` |
| `-c`, `--continue` | ఈ ఫోల్డర్ యొక్క చివరి సంభాషణను కొత్త ఖాతాలో కొనసాగిస్తుంది (కింద చూడండి) | ఆఫ్ |
| `--brief` | `--continue` తో: మొత్తం సంభాషణను resume చేయడానికి బదులుగా, పాత సంభాషణ గురించి ccam రాసిన సారాంశం నుంచి కొత్త ఖాతాగా కొత్త సంభాషణను ప్రారంభిస్తుంది ([కింద చూడండి](#brief)) | ఆఫ్ |
| `-- <claude args>` | ccam సొంత ఆర్గ్యుమెంట్‌ల తర్వాత `claude` కు పంపే ఆర్గ్యుమెంట్‌లు | ఏమీ లేవు |
| `--dry-run` | మద్దతు ఉంది: Claude Code ను ప్రారంభించకుండా, ప్రారంభాన్ని నమోదు చేయకుండా, `→ <from> → <to>` ను, తాను రన్ చేసే కమాండ్ లైన్‌ను ప్రింట్ చేస్తుంది | ఆఫ్ |

`--auto`, `--continue` కలిసి పనిచేస్తాయి. `--continue` తో, `--resume` ఆర్గ్యుమెంట్‌లు మీ సొంత ఆర్గ్యుమెంట్‌ల ముందు పెట్టబడతాయి.

`--dangerously-skip-permissions` లేదా `--model` వంటి Claude Code సొంత ఆప్షన్‌లు `--` తర్వాత వస్తాయి. పేరు తర్వాత ఉన్న ప్రతిదాన్నీ ముందుకు పంపే `ccam <name>` లా కాకుండా, `ccam next` తన ఆర్గ్యుమెంట్‌లను తానే చదువుతుంది, కాబట్టి `--` లేకుండా అది `ccam: unknown flag: --dangerously-skip-permissions` తో ఆగిపోతుంది.

## ఇది ఎలా పనిచేస్తుంది

ccam ద్వారా జరిగే ప్రతి ప్రారంభం సమయాన్ని `~/.local/state/ccam/last-used/<name>` లో రాస్తుంది. `ccam next` ఈ ఫైళ్లను చదివి, ఏ ఖాతాలు లాగిన్ అయి ఉన్నాయో తనిఖీ చేసి, ఎంచుకుంటుంది.

1. ccam `config.toml` నుంచి ప్రతి ఖాతాకు ఒక `next.Candidate` ను తయారు చేస్తుంది (`cli.candidates`): ఖాతా, దాని చివరి ప్రారంభ సమయం (`launcher.LastUsed`), దానికి పనిచేసే లాగిన్ ఉందా లేదా (`creds.Has`: `.credentials.json` లో, లేదా macOS లో Keychain లో, పూర్తిగా గడువు ముగియని, Claude Code క్లియర్ చేయని Claude లాగిన్). బదులుగా API కీ ఉపయోగించే ఫోల్డర్ (Console లాగిన్, లేదా `ccam add --api` తో జోడించిన ఖాతా) కు ప్లాన్ పరిమితులు ఉండవు, అది ఎప్పుడూ ఎంచుకోబడదు.
2. `--auto` లేకుండా, ccam ఖాతాలను స్లాట్ నంబర్ ప్రకారం క్రమబద్ధీకరించి, అత్యంత ఇటీవల ప్రారంభించిన దాన్ని కనుగొంటుంది (`next.Rotate`). మీరు దాని నుంచే మారుతున్నారు. ఇంకా ఏ ఖాతా ప్రారంభించబడకపోతే, అది అతి తక్కువ స్లాట్.
3. అక్కడి నుంచి అది స్లాట్ క్రమంలో ముందుకు వెళ్తుంది, చివరి స్లాట్ నుంచి మొదటిదానికి తిరిగి చుట్టుకుంటుంది, లాగిన్ అయిన మొదటి ఖాతాను తిరిగి ఇస్తుంది.
4. ccam ఎంచుకున్న ఖాతాను `ccam <name>` చేసే విధంగానే సిద్ధం చేసి ప్రారంభిస్తుంది (`launcher.LaunchWith`): ఖాతా లింక్‌లను తనిఖీ చేస్తుంది, MCP షేరింగ్‌ను వర్తింపజేస్తుంది, సెషన్‌ను నమోదు చేస్తుంది, `syscall.Exec` ద్వారా ccam ను `claude` తో భర్తీ చేస్తుంది.

ఈ ప్రయాణం చివరిగా ఉపయోగించిన ఖాతా వద్ద మొదలై, దాని తర్వాత లాగిన్ అయిన మొదటి ఖాతా వద్ద ఆగుతుంది:

<!-- motion:next-rotate -->

`work` అత్యంత ఇటీవల ప్రారంభించబడింది, కాబట్టి ప్రయాణం దాని తర్వాత మొదలవుతుంది. `side` కు లాగిన్ లేదు కాబట్టి వదిలేయబడుతుంది, `lab` ఎంచుకోబడుతుంది. `lab` కూడా లాగ్ అవుట్ అయి ఉంటే, ప్రయాణం `personal` కు చుట్టుకునేది (డయాగ్రామ్‌లోని `% len`). రొటేషన్ మీరు మొదలుపెట్టిన ఖాతాను ఎప్పుడూ తిరిగి ఇవ్వదు: లాగిన్ అయిన మరో ఖాతా లేకపోతే, అది ఒక ఎర్రర్‌తో ఆగుతుంది.

### --auto

`--auto` స్లాట్ క్రమం బదులు వినియోగం ఆధారంగా నిర్ణయిస్తుంది. ccam లాగిన్ అయిన ఖాతాలను ఉంచుకుని, `ccam usage` చేసే విధంగానే వాటన్నింటి వినియోగాన్ని ఒకేసారి అడిగి (`Fetcher.FetchAll`), కింది నియమాల ప్రకారం ఎంచుకుంటుంది (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

ఒక ఖాతాను దాని 5 గంటల, వారపు విండోలతో అంచనా వేస్తారు, Claude Code కోసం ఇచ్చిన ఆర్గ్యుమెంట్‌లలో (`--` తర్వాత) `--model` ఉంటే ఆ మోడల్ వారపు విండోతో కూడా (పేరులో `opus`, `sonnet`, `haiku` లేదా `fable`). ccam ఖాతాలను మూడు వర్గాలుగా విభజిస్తుంది (`Policy.Choose()`):

1. **అవకాశం ఉన్నవి**: ప్రతి విండో `auto_max_usage` (`config.toml`, డిఫాల్ట్ 100) కంటే తక్కువ. వాటిలో strategy ఎంచుకుంటుంది: `headroom` (డిఫాల్ట్) అతి తక్కువ 5 గంటల వినియోగాన్ని తీసుకుంటుంది, 5 గంటల వినియోగం సమానంగా ఉన్న ఖాతాలలో తన వారంలో అతి తక్కువ వాడిన ఖాతాను తీసుకుంటుంది; `expiring` అన్నిటికంటే ముందు రీసెట్ అయ్యే 5 గంటల విండోను తీసుకుంటుంది, కాబట్టి దాని మిగిలిన quota కోల్పోక ముందే ఉపయోగించబడుతుంది; `balance` వారం రీసెట్ అయ్యే వరకు రోజుకు ఎక్కువ వారపు quota మిగిలి ఉన్న ఖాతాను తీసుకుంటుంది, రన్ అవుతున్న ప్రతి ccam సెషన్‌ను దానికి వ్యతిరేకంగా లెక్కిస్తూ, కాబట్టి వారపు quota వృథా కాదు, సమాంతర సెషన్‌లు పంచబడతాయి.
2. **మించినవి**: 100% కంటే తక్కువ, కానీ `auto_max_usage` వద్ద లేదా అంతకంటే ఎక్కువ. ఏ ఖాతాలోనూ అవకాశం లేనప్పుడే ఉపయోగించబడతాయి; అతి తక్కువ 5 గంటల వినియోగం ఉన్నది ముందుగా వస్తుంది, సమానంగా ఉంటే తన వారంలో అతి తక్కువ వాడినది.
3. **నిండినవి**: 100% వద్ద ఒక విండో, లేదా తెలియని 5 గంటల వినియోగం. ఎప్పుడూ ఎంచుకోబడవు.

Claude Code ఒక విండోలో అభ్యర్థనను తిరస్కరించినప్పుడు, ఆ విండో దాని రీసెట్ వరకు 100% గా లెక్కించబడుతుంది: సేవ్ చేసిన సంఖ్యలు ఏమి చెప్పినా, [ccam ప్లగిన్](#inside-claude-code) hook దాన్ని నమోదు చేస్తుంది (`~/.local/state/ccam/usage/<name>.limit.json`). కాబట్టి ఇప్పుడే పరిమితిని చేరిన ఖాతా, దాని సంఖ్యలు పాతవైనా, రీసెట్‌కు ముందు మళ్లీ ఎంచుకోబడదు.

అన్ని ఖాతాల సెట్టింగ్‌లు `config.toml` పైభాగంలో, మొదటి `[[group]]`, `[[profile]]` లేదా `[[folder]]` కంటే పైన ఉండాలి. వాటిలో ఒకదాని కింద ఉన్న సెట్టింగ్ ఆ టేబుల్‌కే చెందుతుంది: ccam ప్రతి లాంచ్‌లో దాని గురించి హెచ్చరిస్తుంది, దాన్ని తరలించే వరకు కాన్ఫిగ్‌ను సేవ్ చేసే కమాండ్‌లు నిరాకరిస్తాయి.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**సంస్థలు.** ఆటోమేటిక్ ఎంపికలు మీరు వదిలే ఖాతా లాంటి రకానికే చెందిన ఖాతాల మధ్యే ఉంటాయి: అదే Team లేదా Enterprise సంస్థలోని సీట్లు, లేదా వ్యక్తిగత ప్లాన్‌లు (Pro, Max), కాబట్టి కంపెనీ ఖాతా నుంచి వచ్చిన పని వ్యక్తిగత ఖాతాలో కొనసాగదు, అలాగే వ్యక్తిగత పని కంపెనీ ఖాతాలో కొనసాగదు. ప్రతి ఖాతా సంస్థను ccam దాని లాగిన్ నుంచి తెలుసుకుంటుంది (`ccam login` దాన్ని నమోదు చేస్తుంది; పాత లాగిన్‌లు `.claude.json` నుంచి చదవబడతాయి); ccam సంస్థను చెప్పలేని ఖాతా ఏ ఖాతాతోనైనా సరిపోలవచ్చు. ఒకే Claude ఖాతాలోకి లాగిన్ అయిన రెండు ఫోల్డర్‌లు దాని పరిమితులను పంచుకుంటాయి, కాబట్టి ఒకటి మరొకదాని స్థానాన్ని ఎప్పుడూ తీసుకోదు. `--cross-org`, లేదా `config.toml` లో `cross_org = true`, సంస్థ నియమాన్ని ఎత్తివేస్తుంది, కానీ [`/ccam:switch`](#inside-claude-code) తర్వాత మరో సంస్థకు తరలింపు ఇప్పటికీ మీ అవును కోసం వేచి ఉంటుంది.

ఏదీ ఎంచుకోలేనప్పుడు, ccam కోడ్ 75 (`EX_TEMPFAIL`) తో నిష్క్రమిస్తుంది, కాబట్టి స్క్రిప్ట్ "తర్వాత ప్రయత్నించు" ను వైఫల్యం నుంచి వేరు చేయగలదు. `--wait` తో, బదులుగా అది మొదటి నిండిన ఖాతా అడ్డుకుంటున్న విండో రీసెట్ అయ్యే వరకు నిద్రిస్తుంది (రీసెట్ సమయం తర్వాత 30 సెకన్లకు, ఎండ్‌పాయింట్ కొత్త విండోను చూపించినప్పుడు), `ccam: no account has room; waiting until 15:10` అని ప్రింట్ చేసి, మళ్లీ ఎంచుకుంటుంది; ఏ ఖాతా ఎప్పుడు రీసెట్ అవుతుందో చెప్పనప్పుడు, 10 నిమిషాల తర్వాత మళ్లీ చూస్తుంది. వేచి ఉండటం ముగిసే లోపు ఏదీ ఖాళీ కాకపోతే, 75 తో నిష్క్రమిస్తుంది.

ఈ లుకప్ `ccam usage` నియమాలను పాటిస్తుంది: 60 సెకన్ల కంటే తక్కువ పాతవైన సంఖ్యలు మళ్లీ ఉపయోగించబడతాయి, ప్రతి ఖాతా లుకప్ 3 సెకన్ల తర్వాత timeout అవుతుంది, సమాధానం ఇవ్వని లేదా యాక్సెస్ టోకెన్ గడువు ముగిసిన ఖాతాను దాని సేవ్ చేసిన సంఖ్యలతో అంచనా వేస్తుంది, వాటిపై పరిమితి రికార్డ్ ప్రాధాన్యం పొందుతుంది. రీసెట్ సమయం దాటిన సేవ్ చేసిన విండో 0% గా లెక్కించబడుతుంది, కాబట్టి ఆ తర్వాత విండో రీసెట్ అయిన ఖాతా దాటవేయబడదు. `--auto` ఎప్పుడూ టోకెన్‌లను పునరుద్ధరించదు. సాధారణ రొటేషన్‌లా కాకుండా, మీరు చివరిగా ఉపయోగించిన ఖాతాకు ఇంకా ఎక్కువ అవకాశం ఉంటే `--auto` దాన్నే ఎంచుకోవచ్చు.

### --continue

1. ccam ప్రతి ఖాతాలో ఈ ఫోల్డర్ యొక్క చివరి సంభాషణ కోసం చూస్తుంది: `projects/<folder>/` లోని `<session>.jsonl` ఫైళ్లలో (ఫోల్డర్ పేరు అంటే లింక్‌లు పరిష్కరించిన working ఫోల్డర్, అందులో అక్షరం లేదా అంకె కాని ప్రతి అక్షరం `-` గా మార్చబడుతుంది, Claude Code పేరు పెట్టే విధంగానే), అత్యంత కొత్త prompt లేదా సమాధానం ఉన్నది. అవి లేని సెషన్‌లు, ఉదాహరణకు `/mcp`, `/usage` లేదా `/status` మాత్రమే రన్ చేసిన సెషన్, లేదా ఖాళీ ఫైల్, వదిలేయబడతాయి. ఒక సెషన్ నోట్ వదిలి ఉంటే (వినియోగ పరిమితి లేదా `/ccam:switch`, కింద చూడండి), దాని తర్వాత 2 నిమిషాలకు మించి ఏ సంభాషణలోనూ prompt లేదా సమాధానం లేకపోతే, బదులుగా సరిగ్గా ఆ సెషన్ ఉపయోగించబడుతుంది; `/exit` వద్ద Claude Code రాసే వరుసలు లెక్కలోకి రావు.
2. అది చెందిన ఖాతా నుంచే మీరు మారుతున్నారు, కొత్త ఖాతా పైన చెప్పినట్లుగా, అదే సంస్థలో ఎంచుకోబడుతుంది. అనేక ఖాతాలు ఆ చరిత్రను షేర్ చేసినప్పుడు, ఆ సంభాషణను చివరిగా రన్ చేసిందని తన `.claude.json` చెప్పే ఖాతాకు (`lastSessionId`), లేకపోతే చివరిగా ఉపయోగించిన ఖాతాకు అది చెందుతుంది.
3. కొత్త ఖాతా ఆ చరిత్రను షేర్ చేస్తే (అదే గ్రూప్ ఫోల్డర్), ccam `claude --resume <session id>` ను, అంటే సరిగ్గా ఆ సంభాషణనే, ప్రారంభిస్తుంది (`--continue` అయితే ఫోల్డర్‌లో ఏది అత్యంత కొత్తదో దాన్ని తెరుస్తుంది).
4. లేకపోతే అది `claude --resume <the conversation's file> --fork-session` ను ప్రారంభిస్తుంది: Claude Code మొత్తం సంభాషణను, కొత్త సెషన్ ID తో, కొత్త ఖాతా సొంత చరిత్రలోకి కాపీ చేస్తుంది, అసలుది ఉన్నట్లే ఉంటుంది. అప్పటి నుంచి కొత్త ఖాతా `--continue` దాన్ని కనుగొంటుంది.
5. మీరు `--` తర్వాత `--agent` ఇవ్వకపోతే, సంభాషణ ఏ agent గా నడుస్తుందో ఆ agent ను కూడా ccam అందిస్తుంది: ఫోల్డర్ యొక్క `.claude/settings.local.json` లేదా `.claude/settings.json` లోని, లేదా ఖాతా `settings.json` లోని `agent` సెట్టింగ్, లేకపోతే సంభాషణ చివరి `agent-setting` వరుస పేర్కొనే agent. సంభాషణను కాపీ చేసేటప్పుడు (దశ 4) Claude Code agent ను తిరిగి తీసుకురాదు, అది తన డిఫాల్ట్ agent గా, అన్ని టూల్‌లతో, కొనసాగేది. కొత్త ఖాతా పాత ఖాతా agents, ప్లగిన్‌లు, సెట్టింగ్‌లను షేర్ చేయనప్పుడు (గ్రూప్ ద్వారా అవే `agents`, `plugins` ఫోల్డర్‌లు, `settings.json`), తరలింపు గురించిన లైన్ తర్వాత ccam ఒక లైన్ జోడిస్తుంది: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

సంభాషణ ఎంత పెద్దదో ccam చెబుతుంది: కొత్త ఖాతా తన మొదటి సమాధానంతో దాన్నంతా మళ్లీ చదువుతుంది. prompt caches సంస్థల మధ్య ఎప్పుడూ షేర్ కావు, ప్రతి వ్యక్తిగత ప్లాన్ తనకు తానే ఒక సంస్థ, కాబట్టి రెండు ఖాతాలు వేర్వేరు సంస్థలకు చెందినప్పుడు ccam `prompt caches are not shared between organizations, so the first reply reads it all again` అని జోడిస్తుంది. పాత ఖాతా సైన్ ఇన్ అయిన MCP సర్వర్లకు కొత్త ఖాతా సైన్ ఇన్ కాకపోతే, ccam వాటి పేర్లు చెబుతుంది: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. ఒక గ్రూప్‌లో `history` ను షేర్ చేయడం కాపీని నివారిస్తుంది, ఒకే సంభాషణను ఉంచుతుంది, `/rewind` కోడ్‌ను పునరుద్ధరించే checkpoints (`file-history/`) ను కూడా షేర్ చేస్తుంది.

### --brief

`--continue` మొత్తం సంభాషణను resume చేస్తుంది, కొత్త ఖాతా తన మొదటి సమాధానంతో దాన్నంతా మళ్లీ చదువుతుంది. `ccam next --continue --brief` బదులుగా కొత్త సంభాషణను ప్రారంభిస్తుంది, ccam మోడల్‌ను అడగకుండా పాత సంభాషణ గురించి రాసే సారాంశం నుంచి:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

సారాంశంలో, ఈ క్రమంలో, చెప్పడానికి ఏదైనా ఉన్న చోట మాత్రమే ఇవి ఉంటాయి:

1. సంభాషణ శీర్షిక;
2. దాని transcript పాత్, ఏదైనా వివరాలు అవసరమైనప్పుడు Claude వెతకడానికి;
3. మీ మొదటి అభ్యర్థన, నాలుగు కంటే ఎక్కువ ఉన్నప్పుడు;
4. సంభాషణను చివరిసారి compact చేసినప్పుడు Claude Code రాసిన సారాంశం;
5. మీ చివరి నాలుగు అభ్యర్థనలు, `!` కమాండ్‌లు, slash కమాండ్‌లు మీరు టైప్ చేసినట్లే;
6. Claude చివరి సమాధానం;
7. todo జాబితా;
8. Claude చివరిగా మార్చిన 10 ఫైల్‌లు;
9. ఫోల్డర్ యొక్క `git status --short --branch`, 800 అక్షరాల వరకు.

todo జాబితా గరిష్ఠంగా 10 అంశాలను ఉంచుతుంది: ముందుగా పూర్తయిన అంశాలు తీసేయబడతాయి, తర్వాత అన్నింటికంటే దూరంగా ఉన్న దశలు, కాబట్టి ప్రస్తుత దశ, తర్వాతి దశలు మిగులుతాయి. పొడవైన భాగాలు మధ్యలో కత్తిరించబడతాయి, మొత్తం సారాంశం 16,000 అక్షరాల లోపే ఉంటుంది.

ccam దాన్ని `~/.local/state/ccam/brief/<session>.md` గా సేవ్ చేస్తుంది, దాన్ని మీరు మాత్రమే చదవగలరు, తర్వాతి `--brief` 7 రోజుల కంటే పాత సారాంశాలను తొలగిస్తుంది. Claude Code ఒక prompt తో ప్రారంభమవుతుంది, అది పేర్కొనే ఫైల్‌ను జత చేస్తుంది: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` `--` తర్వాత మీరు ఇచ్చే prompt చివరి వాక్యం స్థానంలోకి వస్తుంది (`Then: <your prompt>`), `--` తర్వాత మీ ఇతర ఆప్షన్లు దాని తర్వాత వస్తాయి.

Claude Code ఒక mention ను `"` వద్ద ముగిస్తుంది, `#` వద్ద కత్తిరిస్తుంది, కాబట్టి సారాంశం పాత్‌లో వాటిలో ఏదైనా ఉంటే, ccam దాన్ని working ఫోల్డర్‌కు సాపేక్షంగా ఉన్న పాత్‌తో పేర్కొంటుంది; ఆ పాత్‌లో కూడా ఒకటి ఉంటే, prompt ఫైల్ పేరు చెబుతుంది, Claude దాన్ని తన Read టూల్‌తో చదువుతుంది, ఫైల్ working ఫోల్డర్‌ల బయట ఉన్నందున అది ముందుగా అడుగుతుంది (Manual mode లోని `-p` రన్‌లో, Claude Code ఆ చదవడాన్ని నిరాకరిస్తుంది, Claude సారాంశం లేకుండానే కొనసాగుతుంది).

`ccam --dry-run next --continue --brief` సారాంశాన్ని ప్రింట్ చేస్తుంది, ఏదీ సేవ్ చేయదు. పాత సంభాషణ అలాగే ఉంటుంది, కొత్తదానికి సారాంశం చెప్పేది మాత్రమే తెలుసు, కాబట్టి మీరు అక్కడికి ఎలా చేరారో ప్రతి వివరం కంటే పని ముఖ్యమైనప్పుడు `--brief` ను ఉపయోగించండి. `--continue` లేకుండా, ccam `--brief works with --continue` తో ఆగిపోతుంది; సంభాషణను చదవలేనప్పుడు లేదా అందులో సంగ్రహించడానికి ఏమీ కనిపించనప్పుడు, అది ఆగి, `--brief` లేని `ccam next --continue` దాన్ని మొత్తంగా resume చేస్తుందని చెబుతుంది. [రిలే](#relay), [`auto_switch`](#auto-switch) ఎల్లప్పుడూ మొత్తం సంభాషణను resume చేస్తాయి. Claude Code సొంత `--brief` ఆప్షన్, దాని ఇతర ఆప్షన్‌ల లాగే, `--` తర్వాత వస్తుంది.

## Claude Code లోపల

[ccam ప్లగిన్](/docs/install#claude-code-plugin) రెండు hooks ను జోడిస్తుంది (Claude Code 2.1.78 లేదా తర్వాతి వెర్షన్):

- ఒక turn వినియోగ పరిమితి వద్ద ఆగినప్పుడు, అది సంభాషణను గుర్తుపెట్టి, సపోర్ట్ చేసే టెర్మినల్‌లలో (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty) డెస్క్‌టాప్ నోటిఫికేషన్‌ను చూపిస్తుంది: ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  Claude Code వినియోగ పరిమితిగా లెక్కించే ప్రతి సందేశం దానికి తెలుసు: 5 గంటల, వారపు, సాధారణ, మోడల్ పరిమితులు (`You've hit your session limit`, `You've hit your limit`, ఒక మోడల్ కోటా కోసం `You've reached your …`), usage credits, spend limits, ఖాతాను ఆపిన సంస్థ.

  తిరస్కరించిన అభ్యర్థనను Claude Code సంభాషణలో రాసినప్పుడు, hook నిండిన విండోను దాని రీసెట్ వరకు కూడా నమోదు చేస్తుంది ([పరిమితి రికార్డ్](/features/usage#the-limit-record)), సంస్థ విధించిన ఆపివేత తప్ప, ఎందుకంటే దాని అడ్మిన్‌లు దాన్ని ఎప్పుడైనా ఎత్తివేయగలరు. సర్వర్ చేసే చిన్న throttle లేదా API కీ rate limit వినియోగ పరిమితి కాదు, వాటిని అలాగే వదిలేస్తుంది, ఎవరూ పర్యవేక్షించని సెషన్‌లను (`-p` రన్‌లు, బ్యాక్‌గ్రౌండ్ సెషన్‌లు) కూడా. వినియోగ పరిమితి వద్ద ఆగిపోయే subagent లేదా workflow agent, ప్రధాన సంభాషణ దాని కోసం కేవలం ఎదురుచూస్తున్నప్పుడు కూడా, దాని సెషన్ ఆగిపోవడంగానే లెక్కించబడుతుంది: Claude Code ఆ తిరస్కరించిన అభ్యర్థనను సెషన్ `subagents` ఫోల్డర్‌లోని ఆ agent సొంత ఫైల్‌లో రాస్తుంది, hook దాన్ని అక్కడి నుంచి చదువుతుంది. చాలా agents ఒకేసారి ఆగిపోతే, సెషన్‌కు ఒకే నోట్, ఒకే నోటిఫికేషన్ వస్తాయి.

  Claude Code 2.1.234, తర్వాతి వెర్షన్‌లు కూడా సెషన్‌లో వేచి ఉండి, 5 గంటల రీసెట్ తర్వాత తమంతట తామే కొనసాగుతాయి; బదులుగా మారే ఎంపికను నోటిఫికేషన్ మీకు ఇస్తుంది, [`auto_switch`](#auto-switch) ఆన్‌లో ఉంటే (Linux, macOS), మరో ఖాతాలో అవకాశం ఉన్నప్పుడు ccam మీ కోసం మారుస్తుంది. ఏ ఇతర ఖాతా కూడా పనిని తీసుకోలేనప్పుడు, అంటే అదే Team లేదా Enterprise సంస్థకు చెందిన మరో ఖాతా లేనప్పుడు, లేదా మరో వ్యక్తిగత ప్లాన్ ఏదీ లేనప్పుడు (`cross_org` ఆన్‌లో ఉంటే, మరే Claude ఖాతా అయినా తీసుకోవచ్చు), నోటిఫికేషన్ బదులుగా `You have no other personal plans, so ccam cannot move this conversation.` లేదా `You have no other accounts in Acme, so ccam cannot move this conversation.` తో ముగుస్తుంది.
- `/ccam:switch [account]` సంభాషణను గుర్తుపెట్టి, ఏ అభ్యర్థనా చేయక ముందే turn ను ముగిస్తుంది: `/exit` టైప్ చేసి, ఆపై మీరు పేర్కొన్న ఖాతాలో (లేదా తర్వాతి ఖాతాలో) దాన్ని కొనసాగించడానికి `ccam next --continue` రన్ చేయండి. ఇంకా సంభాషణ లేని సెషన్‌లో, కొనసాగించడానికి ఏమీ లేదని అది చెబుతుంది. మీరు టైప్ చేసేది మాత్రమే దీన్ని రన్ చేస్తుంది: ప్లగిన్ 0.3.2 నుంచి, Claude `/ccam:switch` ను రన్ చేయలేదు, Claude షెడ్యూల్ చేసే పని (CronCreate, ScheduleWakeup లేదా `.claude/scheduled_tasks.json`) దాన్ని Claude కు సాధారణ టెక్స్ట్‌గా పంపుతుంది.

  మీరు పేర్కొన్న ఖాతా సంభాషణ సంస్థ బయట ఉన్నప్పుడు (మరో Team లేదా Enterprise సంస్థ సీటు, సీటుకు బదులుగా వ్యక్తిగత ప్లాన్ లేదా దానికి వ్యతిరేకంగా, లేదా ఏ వైపైనా API కీ ఖాతా), `cross_org` ఏమి చెప్పినా, `ccam next --continue` ముందుగా అడుగుతుంది: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. మీరు ఏ పేరూ ఇవ్వనప్పుడు అది తానే ఎంచుకునే ఖాతాకు ముందు కూడా ఇదే అడుగుతుంది, `cross_org = true` ఉంటే దాన్ని మరో సంస్థ నుంచి తీసుకోవచ్చు: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. `y` లేదా `yes` మాత్రమే దాన్ని తరలిస్తుంది; Linux, macOS లలో, ప్రశ్నకు ముందు టైప్ చేసిన కీలు ముందుగా వదిలేయబడతాయి. ఇంకే సమాధానమైనా, లేదా అడగడానికి టెర్మినల్ లేకపోయినా, సంభాషణ, నోట్ ఉన్న చోటే ఉంటాయి: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (తానే ఎంచుకున్న తర్వాత, `ccam next --continue` మళ్లీ ఎంచుకుంటుంది). ఆటోమేటిక్ ఎంపికల లాగే, ccam సంస్థను చెప్పలేని ఖాతా దేనికైనా సరిపోతుంది.

నోట్ `~/.local/state/ccam/handoff/<account>/<session>.json` లో ఉంటుంది, ప్రతి సెషన్‌కు ఒక ఫైల్, కాబట్టి ఒకే ఖాతా యొక్క రెండు సెషన్‌లు ఒకదాని నోట్‌ను మరొకటి ఎప్పుడూ ఓవర్‌రైట్ చేయవు. అది 12 గంటల పాటు ఉంచబడుతుంది, తర్వాతి సెషన్ ప్రారంభం కాబోతున్నప్పుడు ఒక్కసారే ఉపయోగించబడుతుంది; విఫలమైన ప్రారంభం దాన్ని తర్వాతి ప్రయత్నం కోసం అలాగే వదిలేస్తుంది. Claude Code షెల్‌లో `cd` చేసిన తర్వాత కూడా, అది సెషన్ ప్రారంభమైన ఫోల్డర్‌కే చెందుతుంది. `ccam rename` అందులోని పేర్లను అప్‌డేట్ చేస్తుంది; అందులో పేర్కొన్న ఖాతా ఆ తర్వాత తీసివేయబడితే, ccam ఆ విషయం చెప్పి, ఎప్పటిలాగే ఎంచుకుంటుంది.

### రిలే

`config.toml` పైభాగంలో `relay = true` తో, ccam తన స్థానాన్ని Claude Code కు ఇవ్వకుండా ఇంటరాక్టివ్ సెషన్‌లతో పాటు రన్ అవుతూ ఉంటుంది, `/ccam:switch` తర్వాత ఆ సెషన్ సంభాషణను తనంతట తానే కొనసాగిస్తుంది (మరో సెషన్ వదిలిన నోట్‌పై అది చర్య తీసుకోదు): hook Claude Code ను ముగిస్తుంది, ccam ఆ సెషన్ నోట్‌తో `ccam next --continue` ను రన్ చేస్తుంది. ఖాతా పేరు ఇవ్వకపోతే, సంభాషణ మోడల్‌ను కూడా లెక్కించి, ఎక్కువ అవకాశం ఉన్న ఖాతాను తీసుకుంటుంది, `cross_org = true` అది మరో సంస్థ ఖాతా అయ్యేందుకు అనుమతించినప్పుడు, సంభాషణను అక్కడికి తరలించే ముందు పైన ఉన్న ప్రశ్నను అడుగుతుంది. Windows లో hook Claude Code ను ముగించలేదు: `/exit` టైప్ చేయండి, ccam దాన్ని కొనసాగిస్తుంది. సంభాషణ సంస్థ బయటి ఖాతా కోసం కూడా hook దాన్ని ముగించదు: అది `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).` అని చెబుతుంది, మీరు `/exit` టైప్ చేసిన తర్వాత ccam పైన ఉన్న ప్రశ్నను అడుగుతుంది. దానికి ప్రతి సెషన్‌కు వేచి ఉండే ఒక ccam ప్రాసెస్ ఖర్చవుతుంది, అందుకే అది డిఫాల్ట్‌గా ఆఫ్‌లో ఉంటుంది.

సంభాషణను resume చేయడం పునరుద్ధరించని వాటిని రిలే చేసిన సెషన్ ఉంచుకుంటుంది: `--add-dir`, `--mcp-config`, టూల్ జాబితాలు, `--settings`, `--plugin-dir`, system prompt ఆప్షన్‌ల వంటి ఆప్షన్లు, సంభాషణ permission mode, సెషన్ effort level, సంభాషణ నడిచే agent ([`--continue`](#continue) చూడండి). `--model`, `--name` ను Claude Code తానే సంభాషణ నుంచి పునరుద్ధరిస్తుంది. మిగతావి (`--worktree`, `--debug`, `--remote-control`, ccam కు తెలియని ఆప్షన్లు) `ccam: not carried over: …` లో పేర్కొనబడతాయి. ఈ మార్పు ఎప్పుడూ షేర్డ్-చరిత్ర ప్రశ్నను అడగదు, అది రన్ అవుతున్నప్పుడు మీరు టైప్ చేసింది వదిలేయబడుతుంది; దాని ఏకైక ప్రశ్న మరో సంస్థ ఖాతాకు ముందు అడిగేదే.

`/ccam:switch` అడిగిన విధంగా Claude Code ముగిసినప్పుడు మాత్రమే రిలే సంభాషణను కొనసాగిస్తుంది: మీరు ముందుగా `/exit` టైప్ చేసినప్పుడు, అది ముగుస్తున్నప్పుడు Ctrl+C నొక్కినప్పుడు, లేదా ccam కు hangup లేదా మరో signal వచ్చినప్పుడు కాదు, టెర్మినల్ పోయినప్పుడు కూడా కాదు. అప్పుడు నోట్ అలాగే ఉంటుంది, `ccam next --continue` దాన్ని తర్వాత కొనసాగిస్తుంది. ccam పాత వెర్షన్ రిలే కింద ప్రారంభమైన సెషన్‌ను, లేదా ఎవరూ పర్యవేక్షించని సెషన్‌ను (`-p` రన్‌లు, బ్యాక్‌గ్రౌండ్ సెషన్‌లు), hook ముగించదు; బదులుగా దానికి `/exit` సందేశం వస్తుంది. `claude remote-control` రిలే లేకుండా రన్ అవుతుంది.

### `auto_switch`

`config.toml` పైభాగంలో `auto_switch = true` తో (Linux, macOS), Claude Code ఒక సెషన్‌ను వినియోగ పరిమితి వద్ద ఆపినప్పుడు ccam దాన్ని తనంతట తానే మరో ఖాతాకు తరలిస్తుంది: ప్లగిన్ hook అవకాశం ఉన్న మరో ఖాతాను ఎంచుకుని Claude Code ను ముగిస్తుంది, సెషన్‌తో పాటు రన్ అవుతున్న ccam, `/ccam:switch` తర్వాత [రిలే](#relay) చేసినట్లే, సంభాషణను అక్కడ కొనసాగిస్తుంది. ఇది డిఫాల్ట్‌గా ఆఫ్‌లో ఉంటుంది, ప్రయోగాత్మకం.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ఇవన్నీ నిజమైనప్పుడు మాత్రమే ccam ఒక సెషన్‌ను తరలిస్తుంది:

1. ccam దాన్ని టెర్మినల్‌లో ప్రారంభించి, దానితో పాటు రన్ అవుతూ ఉంది. `auto_switch` [రిలే](#relay) ను ఆన్ చేస్తుంది, కాబట్టి దాన్ని సెట్ చేసిన తర్వాత మీరు ప్రారంభించే ఇంటరాక్టివ్ సెషన్‌లకు ఇది నిజం. ముందు ప్రారంభమైన సెషన్‌కు, లేదా పాత ccam కింద ప్రారంభమైన సెషన్‌కు ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` తో ముగిసే నోటిఫికేషన్ వస్తుంది. ఎవరూ పర్యవేక్షించని సెషన్‌లు (`-p` రన్‌లు, బ్యాక్‌గ్రౌండ్ సెషన్‌లు), `claude remote-control`, ccam లేకుండా ప్రారంభించిన Claude Code ఎప్పుడూ తరలించబడవు.
2. ఖాతాలో [ccam ప్లగిన్](/docs/install#claude-code-plugin) ఆన్‌లో ఉంది, అది Claude Code 2.1.214 లేదా తర్వాతి వెర్షన్‌ను రన్ చేస్తుంది, అది తన ప్రాసెస్ ID ని hook కు అందిస్తుంది (`CLAUDE_PID`).
3. ఆ ఖాతా 5 గంటల, వారపు, Opus లేదా Sonnet పరిమితి కోసం, Claude Code తిరస్కరించిన అభ్యర్థనను అది రీసెట్ అయ్యే సమయంతో సహా సంభాషణలో రాసింది. spend limit, usage credits అయిపోవడం, సంస్థ విధించిన ఆపివేత, Fable allowance (`You've reached your Fable limit.`), సర్వర్ throttle, API కీ rate limit ఎప్పుడూ సెషన్‌ను తరలించవు. ఒక subagent లేదా workflow agent పరిమితి లెక్కలోకి వస్తుంది, ప్రధాన సంభాషణ దాని కోసం కేవలం ఎదురుచూస్తున్నప్పుడు కూడా; అది వారపు Opus లేదా Sonnet పరిమితి అయితే, సంభాషణే ఆ మోడల్‌పై నడుస్తున్నప్పుడు మాత్రమే సెషన్ తరలుతుంది.
4. ఖాతా ప్లాన్ ccam కు తెలుసు, దాని `.claude.json` నుంచి లేదా `ccam login` నమోదు చేసిన దాని నుంచి. లేకపోతే, నోటిఫికేషన్ `ccam does not know personal's plan, so auto_switch leaves this session.` తో ముగుస్తుంది.
5. సెషన్ Remote Control కు కనెక్ట్ అయి లేదు. మరో ఖాతాగా అది మీరు దాన్ని నడిపే యాప్‌కు మళ్లీ కనెక్ట్ అవ్వదు, కాబట్టి అది అలాగే ఉంటుంది, నోటిఫికేషన్ ఆ విషయం చెబుతుంది.
6. పరిమితి ఇప్పటి నుంచి 15 నిమిషాల కంటే ఎక్కువ సమయం తర్వాత రీసెట్ అవుతుంది. అంతకంటే తక్కువ నిరీక్షణకు ccam సెషన్‌ను Claude Code కే వదిలేస్తుంది: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. పరిమితి తర్వాత మీరు ఏ prompt నూ టైప్ చేయలేదు.
8. మరో ఖాతాలో అవకాశం ఉంది. ccam [`--auto`](#auto) లాగే ఎంచుకుంటుంది (`auto_max_usage`, `auto_strategy`, సంభాషణ మోడల్), మీరు వదిలే ఖాతా లాంటి రకానికే చెందిన ఖాతాల మధ్య: అదే Team లేదా Enterprise సంస్థలోని మరో సీటు, లేదా మరో వ్యక్తిగత ప్లాన్ (Pro, Max). `cross_org`, `--cross-org` ఇక్కడ వర్తించవు. ఆ ఖాతా:
   - Claude Code రెన్యూ చేయగల లాగిన్‌తో లాగిన్ అయి ఉండాలి;
   - అదే Claude ఖాతాలోకి లాగిన్ అయి ఉండకూడదు;
   - API కీ ఉపయోగించకూడదు;
   - ccam కు తెలిసిన ప్లాన్ కలిగి ఉండాలి;
   - ccam సొంత వినియోగ లుకప్ చూపినంత వరకు, usage credits ఆన్‌లో ఉండకూడదు;
   - Claude Code లో ఈ ఫోల్డర్‌ను trust చేసి ఉండాలి;
   - సంభాషణ ఒక agent గా నడుస్తుంటే (`agent` సెట్టింగ్ లేదా సంభాషణ ఒక agent ను పేర్కొంటుంది), మీరు వదిలే ఖాతా యొక్క agents ను షేర్ చేయాలి: దాని `agents`, `plugins` ఫోల్డర్‌లు, దాని `settings.json` ప్రతిదీ సరిగ్గా అదే అయి ఉండాలి (గ్రూప్ ద్వారా), లేదా రెండింటిలోనూ లేకుండా ఉండాలి.

   ఈ చివరి మూడూ అవకాశం ఉన్న ప్రతి ఖాతాను వదిలేసినప్పుడు, నోటిఫికేషన్ వాటి పేర్లు చెబుతుంది: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` లేదా `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` [పరిమితి రికార్డ్](/features/usage#the-limit-record) ఉన్న విండో నిండినట్లుగా లెక్కించబడుతుంది, గత 2 నిమిషాల్లో మరో సెషన్ తరలించబడిన ఖాతా మిగతావాటి తర్వాత వస్తుంది. ఎంపిక సుమారు 6 సెకన్లలో జరగాలి; వినియోగ లుకప్‌లకు అంతకంటే ఎక్కువ సమయం పడితే, సెషన్ అలాగే ఉంటుంది. `usage_source = "claude"` తో, `auto_switch` సేవ్ చేసిన సంఖ్యల ఆధారంగా మాత్రమే నిర్ణయిస్తుంది, usage credits ఆన్‌లో ఉన్నాయో లేదో ccam Claude Code లుకప్ నుంచి గానీ స్టేటస్ లైన్ నుంచి గానీ తెలుసుకోలేదు, కాబట్టి అవి ఆఫ్‌లో ఉన్నాయని భావిస్తుంది.

అప్పుడు hook Claude Code ను ముగిస్తుంది, దానితో పాటు రన్ అవుతున్న ccam, మద్దతు ఇచ్చే టెర్మినల్స్‌లో డెస్క్‌టాప్ నోటిఫికేషన్ చూపిస్తుంది, ఉదాహరణకు `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` అది పైన ఉన్న ఉదాహరణ అవుట్‌పుట్‌లోని మొదటి లైన్‌ను ప్రింట్ చేసి, ఎంచుకున్న ఖాతాగా ఆ సెషన్ కోసం `ccam next --continue` ను రన్ చేస్తుంది, అది మిగతా రెండు లైన్‌లను ప్రింట్ చేసి, [రిలే](#relay) లాగే, సెషన్ ఆప్షన్లు, permission mode, effort level, agent తో Claude Code ను ప్రారంభిస్తుంది. కొత్త సెషన్ మొదటి సందేశం ఒక స్థిరమైన టెక్స్ట్, ccam మీ తరఫున పంపేది ఇది ఒక్కటే (`/ccam:switch` తర్వాత జరిగే తరలింపు ఏమీ పంపదు):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

subagent పరిమితి తర్వాత అది ఇదవుతుంది, ఎందుకంటే subagents, workflows, బ్యాక్‌గ్రౌండ్ కమాండ్‌లు పాత సెషన్‌తో పాటే ఆగిపోయాయి:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

ఆ ఖాతా కూడా తన పరిమితిని చేరితే, అక్కడ కూడా ఇదే జరుగుతుంది. పరిమితి రికార్డులు ఈ గొలుసు రీసెట్‌కు ముందే ఒక ఖాతాకు తిరిగి వెళ్లకుండా చూస్తాయి, కాబట్టి ప్రతి ఖాతా నిండినప్పుడు సెషన్ చివరి ఖాతాలోనే ఉంటుంది.

సంభాషణ ఒక agent గా నడుస్తుందని hook `agent` సెట్టింగ్ నుంచి, లేదా Claude Code అప్పుడప్పుడు, సెషన్ ముగిసినప్పుడు సంభాషణలో రాసే, agent ను పేర్కొనే వరుస నుంచి తెలుసుకుంటుంది. ఇప్పటివరకు సెషన్ `--agent` మాత్రమే దాన్ని పేర్కొంటే, ఎంచుకున్న ఖాతా agents ను షేర్ చేయకపోతే, `ccam next --continue` ఏదీ ప్రారంభించదు, నోట్‌ను ఉంచుతుంది, రెండో డెస్క్‌టాప్ నోటిఫికేషన్ ఆ విషయం చెబుతుంది: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

మరే ఖాతాలోనూ అవకాశం లేనప్పుడు, ఏదీ తరలించబడదు: Claude Code ఎప్పటిలాగే పరిమితిని నిర్వహిస్తుంది, సెషన్‌లో వేచి ఉంటుంది లేదా తన ఆప్షన్‌లను అందిస్తుంది, నోటిఫికేషన్ ఏ ఖాతా ముందుగా ఖాళీ అవుతుందో చెబుతుంది, ఉదాహరణకు `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` సీట్ల కోసం అది `None of your other accounts in Acme has room now` అని చెబుతుంది. ప్రారంభమైన తరలింపును ఆపడానికి, Claude Code మూసుకుంటున్నప్పుడు Ctrl+C నొక్కండి: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

ఒక తరలింపుకు ఖర్చులు ఉన్నాయి. Claude Code మళ్లీ ప్రారంభమవుతుంది, దానికి కొన్ని సెకన్లు పడుతుంది, Claude ప్రారంభించిన బ్యాక్‌గ్రౌండ్ షెల్‌లు, dev servers, agents పాత సెషన్‌తో పాటు ఆగిపోతాయి. [`--continue`](#continue) లాగే, కొత్త ఖాతా తన మొదటి సమాధానంతో మొత్తం సంభాషణను మళ్లీ చదువుతుంది, అది ఎంత పెద్దదో ccam చెబుతుంది; తరలింపు ఎప్పుడూ [`--brief`](#brief) ను ఉపయోగించదు.

`auto_switch` Claude Code పరిమితిని నివేదించక ముందు ఎప్పుడూ మారదు, బ్యాక్‌గ్రౌండ్‌లో మీ వినియోగాన్ని ఎప్పుడూ గమనించదు, లాగిన్‌ను ఎప్పుడూ కాపీ చేయదు, రెన్యూ చేయదు, మరొకరికి అందించదు. Windows లో అది విస్మరించబడుతుంది, `ccam doctor` ఆ విషయం చెబుతుంది; అక్కడ, `relay = true` తో, `/ccam:switch`, `/exit` ఇప్పటికీ సంభాషణను మరో ఖాతాకు తరలిస్తాయి.

**Anthropic నిబంధనలు.** Team, Enterprise సీట్లలో, `auto_switch` మీ పనిని అదే సంస్థలోని మరో సీటుకు తరలిస్తుంది. Anthropic సహాయ కేంద్రం "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team ప్లాన్‌లోని ప్రతి సీటు పేరుతో నిర్దేశించిన ఒక వ్యక్తికి కేటాయించబడుతుంది, ఆ వ్యక్తి మాత్రమే దాన్ని ఉపయోగించవచ్చు) అని, "Users can't route their work through another user's account to use that user's remaining usage" (యూజర్లు మరో యూజర్‌కు మిగిలి ఉన్న వినియోగాన్ని వాడుకోవడానికి తమ పనిని ఆ యూజర్ ఖాతా ద్వారా మళ్లించకూడదు) అని చెబుతుంది ([Team ప్లాన్‌లు](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise ప్లాన్‌లకు](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) కూడా ఇదే వర్తిస్తుంది).

వ్యక్తిగత ప్లాన్‌ల (Pro, Max) విషయంలో, పరిమితి వద్ద మీ సొంత ఖాతాలు ఒకదాని పనిని మరొకటి ఆటోమేటిక్‌గా తీసుకోవచ్చో లేదో Anthropic నిబంధనలు చెప్పవు. Anthropic స్పష్టంగా అనుమతిస్తే తప్ప, [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) సేవలను "through automated or non-human means, whether through a bot, script, or otherwise" (bot, script లేదా మరే విధంగానైనా, ఆటోమేటెడ్ లేదా మానవేతర మార్గాల ద్వారా) ఉపయోగించడాన్ని నిషేధిస్తాయి, ఒక తరలింపు పైన ఉన్న స్థిరమైన సందేశాల్లో ఒకదాన్ని మీ తరఫున పంపుతుంది. నవంబర్ 12, 2026 నుంచి అమలులోకి వచ్చే [Usage Policy](https://www.anthropic.com/legal/aup) తాను నిషేధించే వాటిలో "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (మా ఉత్పత్తుల్లో ఏర్పాటు చేసిన సామర్థ్యాలు, పరిమితులు లేదా రక్షణలను ఉద్దేశపూర్వకంగా దాటవేయడం లేదా వాటికి అడ్డుపడటం) ను జాబితా చేస్తుంది.

`auto_switch` ఈ నిబంధనలకు అనుగుణంగా ఉందని ccam చెప్పదు. దాన్ని ఆన్ చేయడం మీ నిర్ణయం, ఆ రిస్క్ కూడా మీదే: నిబంధనలు ఉల్లంఘించబడ్డాయని Anthropic భావించినప్పుడు, అది ఒక ఖాతాను, లేదా Team లేదా Enterprise సంస్థ యాక్సెస్‌ను నిలిపివేయవచ్చు.

`auto_switch` ఆన్‌లో ఉన్నప్పుడు, ccam ప్లగిన్ లేని ఖాతా గురించి, 2.1.214 కంటే పాత Claude Code వెర్షన్‌ల గురించి `ccam doctor` హెచ్చరిస్తుంది, సీట్ల కోసం Anthropic సీటు నియమాన్ని మళ్లీ చెబుతుంది ([తనిఖీలు](/features/doctor)).

## పరిమితి వద్ద మళ్లీ రన్ చేయడం

`ccam next` ఆగిపోయిన ఇంటరాక్టివ్ సెషన్ కోసం. నాన్-ఇంటరాక్టివ్ రన్, అంటే స్క్రిప్ట్‌లో లేదా షెడ్యూల్ చేసిన job లో ఉన్న `claude -p`, మీ కోసం వేచి ఉండలేదు, కాబట్టి ఖాతా పేరుకు ముందు `--fallback` పెడితే ccam దాని కోసం మారుస్తుంది:

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

1. `--fallback` `-p` లేదా `--print` తో మాత్రమే పనిచేస్తుంది; లేకపోతే ccam ఆగిపోయి `ccam next --auto --continue` ను సూచిస్తుంది.
2. పైప్ చేసిన స్టాండర్డ్ ఇన్‌పుట్ 10 MiB వరకు ఒక్కసారి చదవబడి ప్రతి ప్రయత్నానికీ ఇవ్వబడుతుంది. Claude Code లాగే, అది మొదలవడానికి ccam గరిష్ఠంగా 3 సెకన్లు వేచి ఉంటుంది (`--input-format stream-json` తో తప్ప); పైప్ తెరిచి ఉండి ఏమీ పంపకపోతే, రన్‌కు ఇన్‌పుట్ అందదు, ccam ఆ విషయాన్ని స్టాండర్డ్ ఎర్రర్‌పై నమోదు చేస్తుంది.
3. ccam `exec` కు బదులుగా `claude` ను చైల్డ్ ప్రాసెస్‌గా ప్రారంభించి, దాని స్టాండర్డ్ అవుట్‌పుట్‌ను ఆపి ఉంచుతుంది, స్టాండర్డ్ ఎర్రర్‌ను వచ్చినట్లే అందిస్తుంది.
4. `claude` 0 తో నిష్క్రమించినప్పుడు, లేదా వేరే కారణంతో విఫలమైనప్పుడు, ఆపి ఉంచిన అవుట్‌పుట్, ఎగ్జిట్ కోడ్ మార్పు లేకుండా అందించబడతాయి. signal తో ముగిసిన రన్ (128 కంటే ఎక్కువ ఎగ్జిట్ కోడ్) కూడా అంతే: దాన్ని ఆపమని అడిగారు, అది మళ్లీ రన్ అవ్వదు.
5. Claude Code సందేశం అలా చెప్పినప్పుడు, అంటే [hook](#inside-claude-code) కు తెలిసిన వినియోగ-పరిమితి సందేశాల్లో ఏదైనా, `--output-format json` లేదా `stream-json` అవుట్‌పుట్ లోపల కూడా, లేదా ఖాతా వినియోగం దాని 5 గంటల, వారపు లేదా `--model` విండోను 100% వద్ద చూపినప్పుడు (పరిమితి రికార్డ్‌తో సహా), వైఫల్యాన్ని వినియోగ పరిమితిగా పరిగణిస్తారు.
6. అప్పుడు ఆ ప్రయత్నం అవుట్‌పుట్ వదిలేయబడి, అదే కమాండ్ మరో లాగిన్ అయిన ఖాతాగా రన్ అవుతుంది: అతి తక్కువ 5 గంటల వినియోగం ఉన్నది, ఆ తర్వాత వినియోగాన్ని లుకప్ చేయలేని ఖాతాలు, నంబర్ క్రమంలో. 100% వద్ద ఉన్న ఖాతాలు వదిలేయబడతాయి, API కీ లాగిన్ ఉన్న ఖాతాలు ఎప్పుడూ ప్రయత్నించబడవు, ప్రతి ఖాతా గరిష్ఠంగా ఒక్కసారే ప్రయత్నించబడుతుంది.
7. ఆర్గ్యుమెంట్‌లలో `-c`, `--continue`, `-r` లేదా `--resume` ఉంటే, సంభాషణ దొరికేలా, అదే history గ్రూప్‌లోని ఖాతాలు మాత్రమే ప్రయత్నించబడతాయి.
8. `--auto` లాగే, అదే సంస్థకు చెందిన ఖాతాలు మాత్రమే ప్రయత్నించబడతాయి (`cross_org = true` దీన్ని ఎత్తివేస్తుంది).

తదుపరి ఖాతా `--auto` యొక్క అవే నియమాలతో ఎంచుకోబడుతుంది (`auto_max_usage`, `auto_strategy`, `--model` విండో). ఏ ఖాతా మిగలనప్పుడు, ccam చివరి ప్రయత్నం అవుట్‌పుట్‌ను ప్రింట్ చేసి, లాగిన్ అయిన ఏ ఖాతాలోనూ అవకాశం మిగల్లేదని చెప్పి, కోడ్ 75 తో నిష్క్రమిస్తుంది. ఖాతా పేరుకు ముందు `--wait` తో (`ccam --fallback --wait=2h work -p …`), బదులుగా అది మొదటి ఖాతా రీసెట్ అయ్యే వరకు వేచి ఉండి, ఆ తర్వాత అన్ని ఖాతాలను మళ్లీ ప్రయత్నించవచ్చు; ఇప్పుడే పరిమితిని చేరిన ఖాతా వెంటనే మళ్లీ రన్ కాకుండా, అది ఎప్పుడూ ముందుగా కనీసం ఒక్కసారి నిద్రిస్తుంది. ఈ వ్యవధి ఒక్కో నిరీక్షణను కాదు, రన్ మొత్తం నిరీక్షణ సమయాన్ని పరిమితం చేస్తుంది. ఖాళీతో రాస్తే (`--wait 30m`), `ccam next` ఇచ్చే అదే సూచనతో తిరస్కరించబడుతుంది. Claude Code ఆర్గ్యుమెంట్‌ల ముందు `--` ఏ ప్రారంభంలోనైనా లాగే పనిచేస్తుంది (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` మొదటి ప్రయత్నం కమాండ్‌ను ప్రింట్ చేస్తుంది.

మీరు [`auto_switch`](#auto-switch) ఆన్ చేస్తే తప్ప ఇంటరాక్టివ్ సెషన్‌లు తమంతట తామే మారవు. `ccam statusline install` తో, 5 గంటల లేదా వారపు వినియోగం 90% కి చేరగానే స్టేటస్ లైన్ `→ ccam next --auto --continue` ను జోడిస్తుంది ([వినియోగ పరిమితులు](/features/usage)).

## భద్రత, ప్రత్యేక సందర్భాలు

- లాగిన్ లేని ఖాతాలు ఎప్పుడూ ఎంచుకోబడవు. లాగిన్ అయిన మరో ఖాతా లేకపోతే, సాధారణ `ccam next` `no logged-in account other than <name>` తో ఆగుతుంది.
- `--continue` తో, ఫోల్డర్‌లో సంభాషణ లేనప్పుడు: `no account has a conversation in <folder> to continue`.
- సంస్థ నియమం ఏ ఖాతానూ మిగల్చనప్పుడు: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- `--auto` తో, లాగిన్ అయిన ఖాతా ఏదీ లేనప్పుడు: `no logged-in account to pick`. ప్రతి ఖాతా నిండిపోయినప్పుడు లేదా దాని వినియోగం తెలియనప్పుడు: `` no account with known usage and headroom (see `ccam usage`) ``. రెండూ exit కోడ్ 75 తో ముగుస్తాయి.
- `--max-usage`, `--strategy`, `--wait`, `--cross-org` `--auto` తో మాత్రమే పనిచేస్తాయి, 1-100 బయటి విలువ లేదా తెలియని strategy తిరస్కరించబడుతుంది.
- fork సంభాషణను కాపీ చేస్తుంది, ఖాతా ఇతర స్థితిని కాదు: చరిత్ర షేర్ కాకపోతే, అసలు ఖాతా `file-history/` లోని checkpoints వెంట రావు.
- `--dry-run` తో కూడా `--auto` నెట్‌వర్క్ అభ్యర్థనలు చేస్తుంది. అది చదువుతుంది మాత్రమే, టోకెన్ ఫైళ్లలో ఎప్పుడూ రాయదు.
- ఎంచుకున్న ఖాతా `ccam <name>` ఉన్న అదే మార్గం ద్వారా ప్రారంభమవుతుంది, కాబట్టి అవే తనిఖీలు వర్తిస్తాయి, అదే చరిత్రను షేర్ చేసే మరో Claude Code సెషన్ ఈ ఫోల్డర్‌లో ఇప్పటికే రన్ అవుతున్నప్పుడు వచ్చే హెచ్చరికతో సహా.
- Claude Code నిజంగా ప్రారంభమైనప్పుడు మాత్రమే `last-used` రాయబడుతుంది. `--dry-run` దాన్ని ముట్టుకోదు, కాబట్టి డ్రై రన్ రొటేషన్‌ను ముందుకు జరపదు.
- ఒక ప్రయత్నం ముగిసే వరకు `--fallback` స్టాండర్డ్ అవుట్‌పుట్‌ను ఆపి ఉంచుతుంది, కాబట్టి `--output-format stream-json` వంటి స్ట్రీమింగ్ అవుట్‌పుట్ చివర్లో ఒకేసారి వస్తుంది. `stream-json` ఇన్‌పుట్‌ను స్ట్రీమ్ చేసి, ఇంకా పంపే ముందు సమాధానాల కోసం వేచి ఉండే క్లయింట్ `--fallback` ను ఉపయోగించలేదు. పరిమితిని చేరక ముందు మొదటి ప్రయత్నం మార్చిన ఫైల్‌లు మారినట్లే ఉంటాయి; తదుపరి ప్రయత్నం prompt ను మళ్లీ మొదలుపెడుతుంది.

## సంబంధిత పేజీలు

- [వినియోగ పరిమితులు](/features/usage)
- [వినియోగ పరిమితులు, ఖాతాలు మార్చడం](/docs/usage-limits)
- [MCP సర్వర్లు, సెట్టింగ్‌లు, చరిత్రను షేర్ చేయండి](/docs/sharing)
- [Claude Code వినియోగ పరిమితిని చేరినప్పుడు ఏమి చేయాలి](/guides/usage-limit-workflow)
