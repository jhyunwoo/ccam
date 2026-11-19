---
title: ఇతర ఖాతాలపై subagents
description: ccam agent ఒక Claude Code పనిలోని భాగాన్ని మీ మరో ఖాతాగా రన్ చేసి Claude సమాధానాన్ని తిరిగి ఇస్తుంది, కాబట్టి ఒక పొడవైన పని పలు ఖాతాల వినియోగాన్ని వాడుతుంది.
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Claude Code సొంత subagents వాటిని ప్రారంభించిన సెషన్ ఖాతాగానే రన్ అవుతాయి, కాబట్టి అవి ఆ ఖాతా వినియోగాన్నే ఖర్చు చేస్తాయి. `ccam agent "<task>"` ఒక పనిని అదే ఫోల్డర్‌లో `claude -p` తో మీ మరో ఖాతాగా రన్ చేసి, Claude చివరి సమాధానాన్ని ప్రింట్ చేస్తుంది. Claude Code సెషన్ నుంచి పిలిచినప్పుడు, ఇది మరో ఖాతా వినియోగాన్ని ఖర్చు చేసే subagent లాగా పనిచేస్తుంది, కాబట్టి సెషన్ సొంత ఖాతా అయిపోయే సమయం తర్వాత కూడా ఒక పొడవైన పని కొనసాగగలదు.

**Team, Enterprise సీట్లు.** `ccam agent` ఖాతాను తనంతట తానే ఎంచుకుంటుంది. ఒకే Claude Team లేదా Enterprise సంస్థలో అది పనిని మరో సీటుపై రన్ చేస్తుంది, ఇది Anthropic నిబంధనల ఉల్లంఘన కావచ్చు. దానిపై ఆధారపడే ముందు [తరచుగా అడిగే ప్రశ్నలు](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) చదివి, మీ సంస్థను, Anthropic ను సంప్రదించి నిర్ధారించుకోండి.

## వినియోగం

Claude Code దాన్ని తన Bash టూల్ ద్వారా రన్ చేస్తుంది. [ccam ప్లగిన్](/docs/install#claude-code-plugin) ఇన్‌స్టాల్ చేసి ఉంటే, దాన్ని సాధారణ మాటల్లో అడగండి, ఉదాహరణకు "ఈ పనిలోని భాగాల కోసం నా ఇతర ఖాతాలను ఉపయోగించు", లేదా దాన్ని మీ CLAUDE.md లో రాయండి. ప్లగిన్ skill, Claude Code కు ప్రతి agent ను బ్యాక్‌గ్రౌండ్‌లో ప్రారంభించమని, తనంతట తానే అర్థమయ్యే పనిని రాయమని, ఫైళ్లను ఎడిట్ చేసే agents కు వాటి సొంత భాగాన్ని లేదా సొంత git worktree ను ఇవ్వమని చెబుతుంది.

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

Claude చివరి సమాధానం స్టాండర్డ్ అవుట్‌పుట్‌కు వెళ్తుంది; `ccam:` తో మొదలయ్యే లైన్‌లు స్టాండర్డ్ ఎర్రర్‌కు వెళ్తాయి. `--json` బదులుగా `{account, session_id, result, is_error, accounts}` ను ప్రింట్ చేస్తుంది, ఇందులో `accounts` పని ఏ ఏ ఖాతాలుగా రన్ అయిందో ఆ ప్రతి ఖాతాను జాబితా చేస్తుంది.

నిష్క్రమణ కోడ్‌లు: claude సొంత కోడ్ (agent పూర్తయినప్పుడు `0`), పిలిచిన ఖాతా తప్ప మరే ఖాతాలోనూ అవకాశం లేనప్పుడు `75` (`EX_TEMPFAIL`), ఒక agent మరొకదాన్ని ప్రారంభించడానికి ప్రయత్నించడం వంటి ఇతర ఎర్రర్‌లకు `1`.

## ఆప్షన్లు

| ఆప్షన్ | ఏమి చేస్తుంది | డిఫాల్ట్ |
| --- | --- | --- |
| `--model <model>` | agent మోడల్; ఎంపికలో ఆ మోడల్ వారపు విండో కూడా లెక్కించబడుతుంది | ఖాతా డిఫాల్ట్ మోడల్ |
| `--resume <id>` | మునుపటి agent సంభాషణను తదుపరి పనితో కొనసాగిస్తుంది (దాని మొదటి, చివరి లైన్‌లోని సెషన్ id) | ఆఫ్ |
| `--json` | ఖాతా, సెషన్, సమాధానం, పని రన్ అయిన ఖాతాలను ప్రింట్ చేస్తుంది | ఆఫ్ |
| `--max-usage <1-100>` | ఈ వినియోగం కంటే తక్కువ ఉన్న ఖాతాలకు ప్రాధాన్యం ఇస్తుంది | `auto_max_usage`, లేకపోతే 100 |
| `--strategy headroom\|expiring\|balance` | [ఎంపిక నియమం](/features/next#auto) | `auto_strategy`, లేకపోతే `headroom` |
| `--cross-org` | ఇతర సంస్థల ఖాతాలుగా కూడా రన్ చేస్తుంది | `cross_org`, లేకపోతే ఆఫ్ |
| `--wait[=<duration>]` | ఏ ఖాతాలోనూ అవకాశం లేనప్పుడు, మొదటి ఖాతా రీసెట్ అయ్యే వరకు, గరిష్ఠంగా ఇంత సేపు వేచి ఉంటుంది | ఆఫ్; ఒంటరిగా ఇస్తే `5h` |
| `--dry-run` | అది ఎంచుకునే ఖాతా కోసం కమాండ్‌ను, రన్ చేయకుండా ప్రింట్ చేస్తుంది | ఆఫ్ |

## ఇది ఎలా పనిచేస్తుంది

1. **పిలిచిన ఖాతా వదిలేయబడుతుంది.** `ccam agent` ను పిలిచిన ఖాతా (ccam ప్రారంభించిన సెషన్ ఖాతా, లేకపోతే ఈ ఫోల్డర్ ప్రారంభించే ఖాతా) ఎప్పుడూ agent ను రన్ చేయదు, దాని పరిమితులను పంచుకునే, అదే Claude ఖాతాలోకి లాగిన్ అయిన మరో ఫోల్డర్ కూడా రన్ చేయదు. `--cross-org` వేరేగా చెప్పకపోతే, agent పిలిచిన ఖాతా [సంస్థ](/features/next#auto) లోనే ఉంటుంది.
2. **ఇది `ccam pick` లాగా ఎంచుకుంటుంది.** వినియోగాన్ని [`ccam pick`](/features/pick) చూసే విధంగానే చూస్తుంది, తగినన్ని ఖాతాలలో అవకాశం ఉన్నంత వరకు, ఒకదానికొకటి 2 నిమిషాల్లోపు ప్రారంభమైన agents కు వేర్వేరు ఖాతాలు వస్తాయి.
3. **ఇది సెషన్ అనుమతులతోనే రన్ అవుతుంది, ఎప్పుడూ అంతకు మించి కాదు.** `claude -p` లో ఒక టూల్‌ను అనుమతించడానికి ఎవరూ ఉండరు, కాబట్టి agent కు దాన్ని పిలిచిన సెషన్ permission mode, effort వస్తాయి. టెర్మినల్ నుంచి రన్ చేస్తే, Claude Code డిఫాల్ట్ mode, మీ సెట్టింగ్‌లలోని permission నియమాలు వర్తిస్తాయి. `ccam agent` కు అంతకు మించి అనుమతి ఇచ్చే ఆప్షన్ లేదు.
4. **వినియోగ పరిమితి వద్ద, సంభాషణ ముందుకు సాగుతుంది.** agent ఖాతా తన పరిమితిని చేరినప్పుడు, [ప్లగిన్ hook](/features/usage#the-limit-record) నమోదు చేసే విధంగానే ccam పరిమితిని నమోదు చేసి, ఎక్కువ అవకాశం ఉన్న ఖాతాగా అదే సంభాషణను కొనసాగిస్తుంది: Claude Code దాన్ని అక్కడికి కాపీ చేస్తుంది (`--resume <file> --fork-session`), మొదటి సందేశం "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." ప్రతి ఖాతాను ఒక్కసారే ప్రయత్నిస్తుంది.
5. **ప్రతి agent సంభాషణ నమోదు చేయబడుతుంది**, `~/.local/state/ccam/agent/<session id>.json` లో (30 రోజులు). `--resume <id>` దాన్ని దాని ఫోల్డర్‌లో కొనసాగిస్తుంది, దాన్ని రన్ చేసిన ఖాతాలో ఇంకా అవకాశం ఉంటే ఆ ఖాతాగా, లేకపోతే మరో ఖాతాగా. [`ccam next --continue`](/features/next) agent సంభాషణలను వదిలేస్తుంది, కాబట్టి అది మీ సంభాషణనే కొనసాగిస్తూ ఉంటుంది.

## భద్రత, ప్రత్యేక సందర్భాలు

- ఒక agent తన సొంత agents ను ప్రారంభించలేదు: ccam దాన్ని `CCAM_AGENT=1` తో గుర్తు పెడుతుంది, అక్కడ `ccam agent` రన్ అవ్వడానికి నిరాకరిస్తుంది.
- `ccam agent` ను ఆపితే (Ctrl+C, ఒక signal, లేదా Claude Code కమాండ్‌ను మధ్యలో ఆపడం) agent కూడా ఆగిపోతుంది.
- agent కు దాన్ని పిలిచిన సంభాషణ కనిపించదు, దాని పని మాత్రమే కనిపిస్తుంది.
- ఒకే ఫైళ్లను ఒకేసారి ఎడిట్ చేసే agents ఒకదాని పనిని మరొకటి ఓవర్‌రైట్ చేస్తాయి. ప్రతిదానికీ పనిలో దాని సొంత భాగాన్ని, లేదా దాని సొంత git worktree ను ఇవ్వండి.
- పిలిచిన సెషన్ మరో ఖాతాకు మారినప్పుడు ([`auto_switch`](/features/next#auto-switch)), అది ప్రారంభించిన agents దానితో పాటే ఆగిపోతాయి. వాటి సంభాషణలు మిగిలి ఉంటాయి: `--resume <id>` వాటిని కొనసాగిస్తుంది.
- `agent` ఒక కమాండ్ పేరు, కాబట్టి ఏ కొత్త ఖాతాకూ `agent` అని పేరు పెట్టలేరు; ఇప్పటికే ఉన్న అటువంటి ఖాతా ఇప్పటికీ `ccam run agent` తో ప్రారంభమవుతుంది.

## సంబంధిత పేజీలు

- [స్క్రిప్ట్‌ల కోసం ఖాతాను ఎంచుకోవడం](/features/pick)
- [తదుపరి ఖాతాకు మారడం](/features/next)
- [వినియోగ పరిమితులు](/features/usage)
- [ప్రతి ఖాతా బ్యాక్‌గ్రౌండ్ సెషన్‌లు](/features/agents)
