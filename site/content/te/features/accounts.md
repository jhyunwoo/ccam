---
title: ఖాతాలను నిర్వహించడం
description: ccam add, ccam login, ccam list, ccam rename, ccam remove, ccam import తో Claude Code ఖాతాలను జోడించండి, లాగిన్, జాబితా, పేరు మార్పు, తొలగింపు, ఇంపోర్ట్.
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

ప్రతి ccam ఖాతా తన సొంత లాగిన్‌తో ఉన్న వేరే Claude Code కాన్ఫిగ్ ఫోల్డర్. ఖాతాలను జోడించడానికి, వాటిని సైన్ ఇన్ చేయడానికి, జాబితాగా చూడటానికి, పేరు మార్చడానికి, తీసివేయడానికి ఈ కమాండ్‌లను ఉపయోగించండి. మీకు ఇప్పటికే ఉన్న ఫోల్డర్‌లను కూడా ఇవి నమోదు చేస్తాయి.

## వినియోగం

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account that was never logged in or whose login expired
ccam logout side                           # log out: the login and its MCP sign-ins
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, logs out, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

ఏ కమాండ్ అయినా `<name|number>` తీసుకునే చోట, మీరు ఖాతా పేరును లేదా దాని స్లాట్ నంబర్‌ను (`ccam list` లోని `#` కాలమ్) ఇవ్వవచ్చు.

`ccam use` డిఫాల్ట్ ఖాతా ఏదో మాత్రమే మారుస్తుంది (`config.toml` లోని `default`).

ఆప్షన్లు లేకుండా టెర్మినల్‌లో `ccam add` ఒక విజార్డ్‌ను తెరుస్తుంది. అది ఈ దశల గుండా వెళ్తుంది:

1. ఖాతా పేరు, మీరు ఇవ్వకపోతే;
2. ఖాతా ఎలా సైన్ ఇన్ అవుతుంది: ఇప్పుడే Claude లాగిన్, API కీ, లేదా తర్వాత Claude లాగిన్;
3. ఏ గ్రూప్‌లలో చేరాలి, `default` టిక్ చేసి ఉంటుంది (గ్రూప్‌లు ఇంకా లేకపోతే, `default` గ్రూప్ ఏమి షేర్ చేస్తుందో చూపిస్తుంది);
4. `Create a shortcut command for work? [y/N]`;
5. నిర్ధారించడానికి ఒక సారాంశం.

`Esc` ఒక దశ వెనక్కి వెళ్తుంది. `--group`, `--no-group`, `--shortcut`, `--no-login` లో ఏదైనా ఉంటే, లేదా టెర్మినల్ లేకుండా (ఉదాహరణకు, పైప్ నుంచి) రన్ చేస్తే, విజార్డ్ దాటవేయబడుతుంది. అప్పుడు పేరు తప్పనిసరి, `--shortcut` ఒక పేరు ఇస్తే తప్ప షార్ట్‌కట్ సృష్టించబడదు.

`--group` లేకుండా, కొత్త ఖాతా `default` గ్రూప్‌లో చేరుతుంది. మొదటి ఖాతా దాన్ని సృష్టిస్తుంది, అది MCP సర్వర్లు, సెట్టింగ్‌లు, `CLAUDE.md`, ఎక్స్‌టెన్షన్లను షేర్ చేస్తుంది; [గ్రూప్‌లతో షేరింగ్](/features/share) చూడండి. Windows లో ఫైల్‌లను Developer Mode లో మాత్రమే లింక్ చేయవచ్చు. అది లేకుండా, మొదటి ఖాతా `mcp`, `extensions` తో మాత్రమే `default` ను సృష్టించి, ఆ విషయం చెబుతుంది (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), సెట్టింగ్‌లు, memory లేదా చరిత్రను షేర్ చేసే గ్రూప్‌తో `--group` ను ccam తిరస్కరిస్తుంది.

`ccam list` స్లాట్, పేరు, షార్ట్‌కట్, ఇమెయిల్, గ్రూప్‌లు, చివరి వినియోగాన్ని చూపిస్తుంది, డిఫాల్ట్ ఖాతా ముందు `*` ఉంటుంది. ఇమెయిల్ అంటే ఫోల్డర్ ఇప్పుడు లాగిన్ అయి ఉన్న ఇమెయిల్ (దాని `.claude.json` నుంచి); ccam నమోదు చేసిన ఇమెయిల్‌కు అది భిన్నంగా ఉంటే, నమోదు చేసినది దాని తర్వాత `(was …)` గా వస్తుంది:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

```console
$ ccam list
   #  Name      Shortcut  Email                                     Groups   Last used
*  1  work      ccwork    other@company.com (was work@example.com)  default  2m ago
```

`--json` తో, ప్రతి ఖాతా `slot`, `shortcut`, `name`, `email` (పట్టికలో లాగే, ప్రస్తుత లాగిన్‌ది), `default`, `dir`, `last_used` (ఖాతా ఎప్పుడూ ప్రారంభించబడకపోతే వదిలేయబడుతుంది), `groups`, `share` ఉన్న ఒక ఆబ్జెక్ట్. `share` ప్రతి అంశాన్ని (`mcp`, `settings`, `memory`, `extensions`, `history`) అది వచ్చే గ్రూప్‌కు, లేదా `"local"` కు మ్యాప్ చేస్తుంది.

## ఆప్షన్లు

| ఆప్షన్ | ఏమి చేస్తుంది | డిఫాల్ట్ |
| --- | --- | --- |
| `ccam add --no-login` | లాగిన్ చేయకుండా ఖాతాను సృష్టిస్తుంది; తర్వాత `ccam login <name>` తో లాగిన్ చేయండి | ఆఫ్ |
| `ccam add --group <group>` | `default` బదులు ఈ గ్రూప్‌లో చేరుతుంది. మళ్లీ మళ్లీ ఇవ్వవచ్చు; `default` లేకపోతే సృష్టించబడుతుంది | `default` |
| `ccam add --no-group` | ఏ గ్రూప్‌లోనూ చేరదు: ఖాతా ఏదీ షేర్ చేయదు | ఆఫ్ |
| `ccam add --shortcut <command>` | `~/.local/bin` లో ఈ షార్ట్‌కట్ కమాండ్‌ను కూడా సృష్టిస్తుంది | ఏమీ లేదు |
| `ccam add --api --key-from <source>` | లాగిన్ బదులు API కీ ఉపయోగించే ఖాతా; [API కీ ఖాతాలు](/features/api) చూడండి | ఆఫ్ |
| `ccam list --json` | ఖాతాలను JSON గా ప్రింట్ చేస్తుంది | ఆఫ్ |
| `ccam remove -y`, `--yes` | అడగకుండా తీసివేస్తుంది | ఆఫ్ |
| `ccam logout -y`, `--yes` | అడగకుండా లాగ్ అవుట్ చేస్తుంది | ఆఫ్ |
| `ccam import --name <name>` | మీరు ఇచ్చే ఫోల్డర్ కోసం ఖాతా పేరు; ఫోల్డర్ లేకుండా తిరస్కరించబడుతుంది (`--name needs a folder: ccam import <folder> --name ops`) | ఫోల్డర్ పేరు నుంచి తీసుకోబడుతుంది |
| `ccam import -y`, `--yes` | కనిపించిన ప్రతి ఫోల్డర్‌ను అడగకుండా ఇంపోర్ట్ చేస్తుంది | ఆఫ్ |
| `--dry-run` | ఏమీ మార్చని `ccam list`, `ccam completion` దీనికి మద్దతు ఇస్తాయి. `add`, `login`, `logout`, `use`, `rename`, `remove`, `import` దీన్ని `` --dry-run is not supported by `ccam add`; nothing was changed `` (కమాండ్ పేరుతో) అని తిరస్కరిస్తాయి | ఆఫ్ |

`ccam login` (`--expired` తప్ప), `ccam use`, `ccam rename` లకు సొంత ఆప్షన్లు లేవు. `ccam completion` `bash`, `zsh` లేదా `powershell` తీసుకుంటుంది; మరేదైనా `only bash, zsh and powershell are supported` తో విఫలమవుతుంది.

ఖాతా పేర్లు `a-z`, `0-9`, `-`, `_` అక్షరాలతో 1 నుంచి 32 అక్షరాల వరకు ఉంటాయి, ఒక అక్షరం లేదా అంకెతో మొదలవుతాయి. పేరు మొత్తం అంకెలు కాకూడదు, ఎందుకంటే `ccam 2` అంటే స్లాట్ 2. పేరు ccam కమాండ్ కూడా కాకూడదు (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), కాబట్టి `ccam <name>` ఎప్పుడూ అస్పష్టంగా ఉండదు. ఈ జాబితా `internal/config` లోని `config.Reserved`.

## ఇది ఎలా పనిచేస్తుంది

ఖాతాను జోడించడం దాని ఫోల్డర్‌ను సృష్టిస్తుంది, దాని గ్రూప్‌లలో చేరుస్తుంది, కాన్ఫిగ్‌ను సేవ్ చేస్తుంది, సైన్ ఇన్ చేస్తుంది, షార్ట్‌కట్‌ను సృష్టిస్తుంది, ఈ క్రమంలో. మీరు `ccam add work --shortcut ccwork` రన్ చేసినప్పుడు, ccam ఈ దశలను తీసుకుంటుంది (`internal/cli/profiles.go` లోని `addProfile`):

<!-- motion:accounts-add -->

1. ఏదైనా సృష్టించే ముందు ccam పేరును (`config.ValidateName`), షార్ట్‌కట్‌ను (`shim.Check`) తనిఖీ చేస్తుంది.
2. ccam `0700` మోడ్‌తో `~/.local/share/ccam/profiles/work/` ను సృష్టించి, ఖాతాకు అతి తక్కువ ఖాళీ స్లాట్ నంబర్‌ను ఇస్తుంది, కాబట్టి తీసివేసిన ఖాతా నంబర్ మళ్లీ ఉపయోగించబడుతుంది (`profile.Add`, `Config.NextSlot`). మొదటి ఖాతా డిఫాల్ట్ అవుతుంది.
3. ఖాతా తన గ్రూప్‌లలో చేరుతుంది (`--group` లేదా `--no-group` వేరే చెప్పకపోతే `default`; మొదటి ఖాతా `default` ను సృష్టిస్తుంది), ccam కొత్త, ఖాళీ ఫోల్డర్‌ను వాటికి లింక్ చేస్తుంది, ఉదాహరణకు `settings.json` ను `groups/default/settings.json` కు (`applyShare`, `share.Build` తో, ఆపై `Apply`). ఇది విఫలమైతే, కొత్త ఫోల్డర్ తీసివేయబడుతుంది, ఎందుకంటే మిగిలిపోయిన ఫోల్డర్ ప్రతి మళ్లీ ప్రయత్నాన్ని `folder already exists` తో అడ్డుకుంటుంది.
4. ఖాతా `fsutil.AtomicWrite` ద్వారా `~/.config/ccam/config.toml` లో సేవ్ అవుతుంది: ఒక తాత్కాలిక ఫైల్, `fsync`, ఆపై `rename`, కాబట్టి క్రాష్ అయినా సగం రాసిన కాన్ఫిగ్ ఎప్పుడూ మిగలదు.
5. ccam `CLAUDE_CONFIG_DIR` ను ఖాతా ఫోల్డర్‌కు సెట్ చేసి, లాగిన్ పేజీ నిండి ఉండేలా ఖాతా చివరి ఇమెయిల్‌తో `--email` ఇచ్చి, `claude auth login` ను ప్రారంభిస్తుంది (`runLogin`). లాగిన్ విఫలమైతే, ఖాతా అలాగే ఉంటుంది, ccam ఎర్రర్‌ను ప్రింట్ చేసి తర్వాత `ccam login work` రన్ చేయమని చెబుతుంది.
6. ccam `profiles/work/.claude.json` నుంచి `oauthAccount` ను చదువుతుంది (`profile.ReadIdentity`): ఇమెయిల్, ఖాతా, సంస్థ IDs, సంస్థ పేరు, ప్లాన్ (`claude_max`, `claude_team` …). ఆ ఫైల్ అనేక మెగాబైట్లు ఉండవచ్చు, కాబట్టి ccam మొత్తం ఫైల్‌ను పార్స్ చేయకుండా ఆ ఒక్క కీ కోసం స్కాన్ చేస్తుంది (`TopLevelValue`). ccam ఈ వివరాలను `config.toml` లోని `[profile.login]` లో నిల్వ చేస్తుంది.
7. ccam `~/.local/bin/ccwork` ను ccam బైనరీకి లింక్ చేసి (ఆ బైనరీ ప్రారంభమైనప్పుడు `config.toml` లో `argv[0]` కోసం చూస్తుంది), ఆ లింక్‌ను `~/.local/state/ccam/shortcuts.json` లో నమోదు చేస్తుంది (`shim.Install`). లింక్ సృష్టించలేకపోతే, ఖాతా అయినా జోడించబడుతుంది, ccam కారణం చెబుతుంది (`Could not create shortcuts: …`), సారాంశం `shortcut: ccwork (not created; see above)` అని చూపిస్తుంది.

లాగిన్‌కు నిమిషాలు పట్టవచ్చు, ఈలోగా మరో ccam ప్రాసెస్ కాన్ఫిగ్‌ను సేవ్ చేయవచ్చు. కాబట్టి వేచి ఉండే కమాండ్‌లు `config.toml` ను మళ్లీ లోడ్ చేసి, ఆ ఒక్క మార్పును (ఇక్కడ, ఇమెయిల్) వర్తింపజేసి సేవ్ చేస్తాయి (`env.update`). అదే కారణంతో `ccam remove` తన ప్రశ్న తర్వాత మళ్లీ లోడ్ చేస్తుంది.

`ccam login` ఇప్పటికే ఉన్న ఖాతా కోసం దశ 5, 6 ను రన్ చేస్తుంది. Claude Code ఇమెయిల్‌ను లాగిన్ పేజీ కోసం సూచనగా మాత్రమే తీసుకుంటుంది, కాబట్టి మీరు ఇంతకు ముందు కంటే వేరే Claude ఖాతాలో సైన్ ఇన్ అవ్వవచ్చు; ccam కొత్త లాగిన్‌ను `[profile.login]` తో పోల్చి ఆ విషయం చెబుతుంది (`ccam: this login is other@company.com, but work was dev@company.com before.`). అది కొత్త లాగిన్‌ను ఉంచి, దాన్ని ఎలా రద్దు చేయాలో చెబుతుంది: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

వేరే Claude ఖాతాలోకి లాగిన్ అయితే ఆ ఖాతా యొక్క [పరిమితి రికార్డ్](/features/usage#the-limit-record) కూడా తొలగిపోతుంది, ఎందుకంటే ఆ పరిమితులు మరో ఖాతావి. మరో ccam ఖాతా అదే Claude ఖాతాలోకి లాగిన్ అయి ఉంటే, ఆ రెండూ తమ వినియోగ పరిమితులను పంచుకుంటాయని ccam హెచ్చరిస్తుంది.

`ccam login --expired` ఎప్పుడూ లాగిన్ కాని, లేదా లాగిన్ గడువు ముగిసిన ప్రతి ఖాతాలో ఒకదాని తర్వాత ఒకటి లాగిన్ చేస్తుంది: దాని రిఫ్రెష్ టోకెన్ గడువు ముగిసినా, లేదా సర్వర్ తిరస్కరించిన తర్వాత Claude Code లాగిన్‌ను క్లియర్ చేసినా. మొదటి బ్రౌజర్ తెరుచుకునే ముందు, `config.toml` ను సేవ్ చేయగలదో లేదో అది తనిఖీ చేస్తుంది.

`ccam logout` ఖాతాలో `claude auth logout` ను రన్ చేసి సేవ్ చేసిన ఇమెయిల్‌ను తొలగిస్తుంది; లాగిన్ రికార్డ్ అలాగే ఉంటుంది, కాబట్టి తర్వాతి లాగిన్ దానితో పోల్చి తనిఖీ చేయబడుతుంది. Claude Code లాగిన్‌ను రద్దు చేసి, ఆ ఖాతా కోసం తాను నిల్వ చేసిన వాటిని, MCP సర్వర్ సైన్-ఇన్‌లతో సహా, తొలగిస్తుంది. కాబట్టి ఈ సందర్భాల్లో ccam ముందుగా హెచ్చరించి అడుగుతుంది (`--yes` ఈ ప్రశ్నను దాటవేస్తుంది):

- లాగిన్‌లో MCP సైన్-ఇన్‌లు ఉన్నప్పుడు;
- ఖాతా రన్ అవుతున్నప్పుడు (దాని సెషన్‌లు కూడా సైన్ అవుట్ అవుతాయి);
- మరో ఖాతా ఫోల్డర్‌లో అదే లాగిన్ ఉన్నప్పుడు (కాపీ చేసిన ఫోల్డర్), దాన్ని కూడా లాగ్ అవుట్ రద్దు చేస్తుంది.

ఉదాహరణకు:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

మరో Claude ఖాతాగా లాగిన్ అవ్వడానికి `ccam login <name>` ఒక్కటే చాలు, అది MCP సైన్-ఇన్‌లను ఉంచుతుంది.

ఖాతాలు ఒకటి కంటే ఎక్కువ సంస్థలకు చెందినప్పుడు `ccam list` ఒక `Organization` కాలమ్‌ను జోడిస్తుంది (వ్యక్తిగత ప్లాన్‌లు `personal` అని చూపిస్తాయి), `ccam list --json` లో ప్రతి ఖాతాకు `plan`, `org` ఉంటాయి.

### రిమోట్ మెషీన్‌లో సైన్ ఇన్ చేయడం

SSH ద్వారా, `claude auth login` యొక్క బ్రౌజర్ భాగం పాడవ్వవచ్చు. Claude Code రిమోట్ మెషీన్‌లోని కాల్‌బ్యాక్ సర్వర్‌లో (`http://localhost:<port>/callback`) సైన్-ఇన్ కోసం వేచి ఉంటుంది, కానీ సైన్-ఇన్‌ను పూర్తి చేసే బ్రౌజర్ మీ కంప్యూటర్‌లో రన్ అవుతుంది, అక్కడ ఆ అడ్రస్ ఎక్కడికీ దారితీయదు, బ్రౌజర్ `ERR_CONNECTION_REFUSED` వంటి ఎర్రర్‌ను చూపిస్తుంది. అయినా ఆ అడ్రస్‌లో వన్-టైమ్ కోడ్ ఉంటుంది, కాబట్టి ccam దాన్ని అందించగలదు:

<!-- motion:accounts-relay -->

ccam మీ టెర్మినల్‌ను Claude Code కు అప్పగించదు (`runLogin`). అది Claude Code స్టాండర్డ్ ఇన్‌పుట్‌ను ఒక పైప్‌కు కలిపి, మీరు టైప్ చేసే ప్రతి లైన్‌ను అందులోకి కాపీ చేస్తుంది (`relayLogin`), కాబట్టి Claude Code యొక్క `Paste code here` ప్రాంప్ట్ వద్ద పేస్ట్ చేసిన కోడ్ కూడా చేరుతుంది. ఒక లైన్ కాల్‌బ్యాక్ అడ్రస్ అయినప్పుడు, హోస్ట్ `localhost`, `127.0.0.1` లేదా `[::1]` అయి, పోర్ట్ ఇవ్వబడి, పాత్ `/callback` అయి, `code`, `state` రెండూ ఉంటేనే ccam దాన్ని అంగీకరిస్తుంది; `http://` వదిలేయవచ్చు (`callbackURL`). ఆపై ccam బ్రౌజర్ చేసి ఉండే విధంగానే, 10 సెకన్ల timeout తో, రీడైరెక్ట్‌లను అనుసరించకుండా, సర్వర్‌లో ఆ అడ్రస్‌ను అభ్యర్థిస్తుంది (`deliverCallback`). ఆ అభ్యర్థన వేచి ఉన్న `claude auth login` యొక్క కాల్‌బ్యాక్ సర్వర్‌ను చేరుతుంది (డయాగ్రామ్‌లోని గీతల లైన్). Claude Code కోడ్‌ను అందుకుని లాగిన్‌ను సేవ్ చేస్తుంది, Linux లో ఖాతా ఫోల్డర్‌లోని `.credentials.json` లో. పేస్ట్ చేసిన లైన్ ccam ను మరో హోస్ట్‌ను సంప్రదించేలా ఎప్పుడూ చేయలేదు.

లాగిన్ మొదలయ్యే ముందు ccam ఇది ప్రింట్ చేస్తుంది:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

విజయవంతమైతే అది `→ Passed the sign-in to Claude Code on this machine.` అని ప్రింట్ చేస్తుంది.

లాగిన్ టెర్మినల్‌లో ఇన్‌పుట్ కోసం వేచి ఉన్నప్పుడు, ccam టెర్మినల్ echo ను, లైన్ మోడ్‌ను ఆపివేసి, మీరు టైప్ చేసే ప్రతి అక్షరానికి ఒక `*` చూపిస్తుంది, కాబట్టి పేస్ట్ చేసిన కోడ్ లేదా కాల్‌బ్యాక్ అడ్రస్ స్క్రీన్‌పై కనిపించదు. ccam లైన్ ఎడిటింగ్‌ను తానే చేస్తుంది (`readLine`): `Enter` లైన్‌ను పంపుతుంది, `Backspace` ఒక అక్షరాన్ని చెరిపేస్తుంది, ఖాళీ లైన్‌పై `Ctrl+D` ఇన్‌పుట్‌ను ముగిస్తుంది. లాగిన్ ముగిసినప్పుడు, అలాగే `Ctrl+C`, `SIGTERM` లేదా `SIGHUP` వచ్చినప్పుడు కూడా టెర్మినల్ పునరుద్ధరించబడుతుంది; ఇవి ఇప్పటికీ ccam ను ముగిస్తాయి. స్టాండర్డ్ ఇన్‌పుట్, స్టాండర్డ్ ఎర్రర్‌పై టెర్మినల్ లేకపోతే, లైన్‌లు వచ్చినట్లే కాపీ అవుతాయి.

### పేరు మార్చడం, తీసివేయడం, ఇంపోర్ట్ చేయడం

`ccam rename`, `ccam remove` ఖాతా ఫోల్డర్‌పై, ccam ఇతర చోట్ల ఖాతా పేరుతో ఉంచే వాటిపై పనిచేస్తాయి. ccam ఆ పాత్‌లను `nameState` లో జాబితా చేస్తుంది:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- వినియోగ కాష్‌లు `~/.local/state/ccam/usage/<name>.json`, `<name>.statusline.json`;
- పరిమితి రికార్డ్ `<name>.limit.json`;
- `~/.local/state/ccam/handoff/<name>/` లోని సెషన్ నోట్స్.

`ccam remove work` వాటిని ఈ క్రమంలో నిర్వహిస్తుంది:

<!-- motion:accounts-remove -->

1. ఖాతా రన్ అవుతుంటే `ccam remove` ముందుగా తిరస్కరిస్తుంది, ఆపై అడుగుతుంది, మీరు సమాధానం ఇచ్చిన తర్వాత మళ్లీ తనిఖీ చేస్తుంది (ప్రశ్న తెరిచి ఉన్నప్పుడు ప్రారంభమైన సెషన్‌ను ఇది గుర్తిస్తుంది).
2. ఏదైనా మార్చే ముందు `config.toml` ను సేవ్ చేయగలదో లేదో అది తనిఖీ చేస్తుంది.
3. ఖాతా ఎక్స్‌టెన్షన్లను షేర్ చేస్తే, ccam గ్రూప్ యొక్క `installed_plugins.json`, `known_marketplaces.json` ను మళ్లీ రాస్తుంది (`RewritePluginPaths`), ఎందుకంటే Claude Code ఈ ఖాతా ఫోల్డర్ ద్వారా వెళ్లే అవకాశం ఉన్న పూర్తి ప్లగిన్ పాత్‌లను నిల్వ చేస్తుంది.
4. ఆ తర్వాత ccam నిర్వహించే ఖాతా (`ccam import` తో నమోదు చేసిన ఖాతా కాకుండా, ccam ఫోల్డర్ సృష్టించిన ఖాతా) లాగ్ అవుట్ అవుతుంది: దాని ఫోల్డర్‌లో `claude auth logout` రన్ అవుతుంది, కాబట్టి పనిచేసే లాగిన్ ఏదీ ట్రాష్‌లో వేచి ఉండదు (లేదా, macOS లో, ఫోల్డర్ పేరుతో ఉన్న Keychain ఐటెమ్‌లో; అక్కడ ఉంటే ఆ పేరుతో తర్వాత వచ్చే ఖాతాకు అది వారసత్వంగా వచ్చేది). API కీ ఖాతాకు, మరో ఖాతా ఫోల్డర్‌లో అదే లాగిన్ ఉన్నప్పుడు (లాగ్ అవుట్ దాన్ని కూడా రద్దు చేస్తుంది) ఈ దశ దాటవేయబడుతుంది. విఫలమైన లాగ్ అవుట్ నివేదించబడుతుంది, తొలగింపు కొనసాగుతుంది.
5. ccam ఫోల్డర్‌ను `~/.local/share/ccam/trash/<name>-<time>/` కు తరలిస్తుంది (`profile.Trash`), ఖాతా `config.toml` నుంచి తొలగించబడుతుంది, అదే పేరుతో తర్వాత వచ్చే ఖాతాకు అవి వారసత్వంగా రాకుండా `nameState` పాత్‌లు తొలగించబడతాయి, షార్ట్‌కట్ లింక్ తీసివేయబడుతుంది: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

చివరికి `config.toml` ను సేవ్ చేయలేకపోతే, ఫోల్డర్ ట్రాష్ నుంచి తిరిగి వస్తుంది. తీసివేసిన ఖాతాను తిరిగి పొందడానికి, ఫోల్డర్‌ను ట్రాష్ నుంచి వెనక్కి తరలించి మళ్లీ లాగిన్ చేయండి.

`ccam rename` ఫోల్డర్ పేరును మారుస్తుంది, దాని `plugins/` లోని ప్లగిన్ పాత్‌లను పాత ఫోల్డర్ నుంచి కొత్తదానికి మళ్లీ రాస్తుంది. అది `nameState` పాత్‌లను కొత్త పేరుకు తరలిస్తుంది, ఖాతా డిఫాల్ట్ అయి ఉంటే అలాగే ఉంచుతుంది. `ccam next --continue` సెషన్ నోట్స్‌ను ఇంకా కనుగొనేలా వాటిని కూడా అప్‌డేట్ చేస్తుంది. షార్ట్‌కట్ కమాండ్ తన పేరును ఉంచుకుంటుంది. ప్లగిన్ పాత్‌లను మళ్లీ రాయడం విఫలమైతే, లేదా `config.toml` ను సేవ్ చేయలేకపోతే, ఫోల్డర్ పేరు వెనక్కి మార్చబడుతుంది.

macOS లో, అలాగే Claude Code లాగిన్‌లను Credential Manager లో నిల్వ చేసినప్పుడు Windows లో, ఫోల్డర్ తన పాత పేరునే ఉంచుకుంటుంది, ఎందుకంటే లాగిన్, MCP సైన్-ఇన్‌లు ఫోల్డర్ పాత్ కింద నిల్వ అవుతాయి, తరలించిన ఫోల్డర్ లాగ్ అవుట్ అవుతుంది. అప్పుడు కాన్ఫిగ్ ఆ ఫోల్డర్‌ను కొత్త పేరు పక్కన `folder` గా నమోదు చేస్తుంది, ccam `its folder stays …: the macOS Keychain keeps the login under the folder's path` అని చెబుతుంది.

`ccam import` ఒక ఫోల్డర్‌ను ఉన్న చోటే, `config.toml` లో `dir` గా నమోదు చేస్తుంది. అది దాన్ని ఎప్పుడూ తరలించదు, కాపీ చేయదు: కాపీ చేసిన రిఫ్రెష్ టోకెన్ కాపీలు ఒకదాన్ని మరొకటి లాగ్ అవుట్ చేసేలా చేస్తుంది.

ఫోల్డర్ ఆర్గ్యుమెంట్ లేకుండా, అది నమోదు కాని `~/.claude`, `~/.claude-*` ఫోల్డర్‌ల కోసం వెతికి ప్రతిదాని గురించి అడుగుతుంది. పేర్లు ఫోల్డర్ నుంచి వస్తాయి (`~/.claude` అంటే `main`, `~/.claude-work` అంటే `work`, `~/.claude-2` అంటే `acct2`). ccam పేరును తీసుకోలేకపోతే, లేదా ఆ పేరు ఇప్పటికే వాడుకలో ఉంటే, అది ఆ ఫోల్డర్‌ను రన్ చేయాల్సిన `--name` కమాండ్‌తో జాబితాలో చూపిస్తుంది. కనుగొన్న ఫోల్డర్ విఫలమైతే అది `Skipped <path>: <reason>` గా చూపబడుతుంది, మిగతావి ఇంపోర్ట్ అవుతాయి.

ccam తన ప్రశ్నలన్నీ ముందుగా అడిగి, ఆపై `config.toml` ను మళ్లీ లోడ్ చేసి ఫోల్డర్‌లను నమోదు చేస్తుంది, కాబట్టి ఈలోగా ఇతర ccam కమాండ్‌లు సేవ్ చేసిన మార్పులు అలాగే ఉంటాయి. `✓` లైన్‌లు సేవ్ తర్వాత వస్తాయి.

`CLAUDE_CONFIG_DIR` లేకుండా, Claude Code `.claude.json` ను మీ హోమ్ ఫోల్డర్‌లో `~/.claude` పక్కన ఉంచుతుంది, కాబట్టి మీరు `~/.claude` ను ఇంపోర్ట్ చేసినప్పుడు దానికి సొంత `.claude.json` లేకపోతే, ccam `~/.claude/.claude.json` నుంచి `~/.claude.json` కు లింక్‌ను సృష్టిస్తుంది. అప్పుడు సాధారణ `claude`, `ccam main` ఒకే ఖాతాను చూస్తాయి.

`~/.claude` కు `~/.claude.json` పక్కన దాని సొంత `.claude.json` ఉంటే, ccam రెండింటినీ అలాగే వదిలేసి ఆ విషయం చెబుతుంది. `ccam main` లోపలి ఫైల్‌ను, సాధారణ `claude` బయటి ఫైల్‌ను ఉపయోగిస్తాయి, దీంతో వాటి లాగిన్ వివరాలు, MCP సర్వర్లు వేరుగా ఉంటాయి. ఒక్కటే ఉంచడానికి, లోపలి ఫైల్‌ను లింక్‌తో భర్తీ చేయండి ([doctor](/features/doctor) కూడా దీన్ని నివేదిస్తుంది).

నమోదైన ఖాతాకు ఉన్న లాగిన్‌నే ఫోల్డర్ కలిగి ఉంటే (కాపీ చేసిన ఫోల్డర్), ccam ఆ విషయం చెబుతుంది: ముందుగా ఏది రెన్యూ చేస్తే అది మరొకదాన్ని లాగ్ అవుట్ చేస్తుంది, కాబట్టి దాన్ని ఉపయోగించే ముందు `ccam login <name>` రన్ చేయండి.

మీరు ఒక గ్రూప్‌కు జోడించే వరకు (`ccam group join <group> <name>`) ఇంపోర్ట్ చేసిన ఖాతాలు ఏ గ్రూప్‌లోనూ ఉండవు, షార్ట్‌కట్ సృష్టించబడదు.

`ccam completion bash|zsh` ఖాతా పేర్లను కూడా పూర్తి చేసే కంప్లీషన్ స్క్రిప్ట్‌ను ప్రింట్ చేస్తుంది.

## భద్రత, ప్రత్యేక సందర్భాలు

- సెషన్ ఫోల్డర్‌లో రాస్తూనే ఉంటుంది కాబట్టి, ఖాతా రన్ అవుతున్నప్పుడు `ccam rename`, `ccam remove` ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)`` తో తిరస్కరిస్తాయి (`refuseRunning`). Claude Code ఎలా ప్రారంభమైనా ఈ తనిఖీ వర్తిస్తుంది: ccam ద్వారా, సాధారణ `claude`, ఒక IDE, బ్యాక్‌గ్రౌండ్ సెషన్‌లు, వాటి daemon.
- మరో ccam ప్రాసెస్ ఆ లాగిన్‌ను రెన్యూ చేస్తున్నప్పుడు లేదా మారుస్తున్నప్పుడు `ccam rename`, `ccam remove`, `ccam logout` తిరస్కరిస్తాయి: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- మరో ccam ప్రాసెస్‌లో గ్రూప్ మార్పు జరుగుతున్నప్పుడు, `ccam rename`, `ccam remove` `another ccam process is changing the groups right now; try again when it is done` తో తిరస్కరిస్తాయి; ఒక మార్పు మధ్యలో ఆగిపోయిన తర్వాత, ``the last group change stopped halfway; finish it with `ccam doctor --fix` first`` తో.
- ఇంపోర్ట్ చేసిన ఖాతాపై `ccam remove` సాధారణ `claude` కోసం ఫోల్డర్‌ను, దాని లాగిన్‌ను ఉంచుతుంది. ఖాతా గ్రూప్‌లలో ఉంటే, అది ముందుగా `ccam group leave` లాగే వాటి నుంచి బయటకు వెళ్తుంది: ప్రతి షేర్డ్ అంశం ఫోల్డర్‌లో గ్రూప్ దాని కాపీ అవుతుంది, గ్రూప్ MCP సర్వర్లు దాని `.claude.json` లోకి వెళ్తాయి. ఆపై దాని నమోదు తొలగించబడుతుంది: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. ccam సృష్టించిన దాని `.claude.json` బ్యాకప్‌లు `<folder>/.ccam-backup/mcp-<time>` లోకి తరలించబడతాయి; వాటిని తరలించలేకపోతే, అవి `~/.local/state/ccam/mcp/<name>/backups` లోనే ఉంటాయి, సందేశం ఆ విషయం చెబుతుంది. మీ `~/.claude` ట్రాష్‌కు తరలించబడదు.
- రెండు ఖాతాలకు ఒకే ఫోల్డర్ ఉండకూడదు: ఒక ఫోల్డర్‌ను రెండుసార్లు పేర్కొన్న `config.toml` లోడ్ అవ్వడం విఫలమవుతుంది.
- `profiles/` లో ఫోల్డర్ ఇప్పటికే ఉన్న పేరును `ccam add` తిరస్కరించి, `ccam import` ను సూచిస్తుంది. macOS లో, ccam పాత వెర్షన్ వదిలేసిన, కొత్త ఫోల్డర్ పాత్ కోసం లాగిన్ లేదా MCP సైన్-ఇన్‌లు Keychain లో ఇంకా ఉన్నప్పుడు కూడా అది ఆగిపోతుంది, వాటిని తొలగించే `security delete-generic-password` కమాండ్‌ను చెబుతుంది. `config.toml` ను సేవ్ చేయలేకపోతే, కొత్త ఫోల్డర్ మళ్లీ తీసివేయబడుతుంది.
- ఇప్పటికే నమోదైన ఫోల్డర్‌ను, ఫోల్డర్ కాని పాత్‌ను, తీసుకున్న పేరును `ccam import` తిరస్కరిస్తుంది. ఫోల్డర్ `~/.claude` యొక్క కాపీ (అదే లాగిన్) అయితే, `~/.claude` నమోదు కాకపోయినా, ccam హెచ్చరిస్తుంది.
- ఖాతా రన్ అవుతున్నప్పుడు `ccam logout` ముందుగా అడుగుతుంది; `ccam login` తనిఖీ చేయదు.
- రిలే Claude Code ను చేరలేకపోతే, ఉదాహరణకు లాగిన్ ముగిసిన తర్వాత, అది కారణంతో `ccam: could not pass the sign-in to Claude Code:` అని ప్రింట్ చేస్తుంది. HTTP ఎర్రర్ సాధారణంగా కోడ్ గడువు ముగిసిందని అర్థం; లాగిన్‌ను మళ్లీ ప్రారంభించండి.

## సంబంధిత పేజీలు

- [Claude Code ను ప్రారంభించడం](/features/launch)
- [షార్ట్‌కట్ కమాండ్‌లు](/features/shortcuts)
- [MCP సర్వర్లు, సెట్టింగ్‌లు, చరిత్రను షేర్ చేయండి](/docs/sharing)
- [ఒకే కంప్యూటర్‌లో అనేక Claude Code ఖాతాలను ఉపయోగించండి](/guides/multiple-claude-code-accounts)
