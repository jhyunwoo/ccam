---
title: कमांड संदर्भ
metaTitle: "ccam कमांड संदर्भ"
description: ccam add आणि ccam <name> पासून ccam usage, ccam next आणि ccam doctor पर्यंत ccam ची प्रत्येक कमांड आणि पर्याय.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: 28e1ad1dcbad
---

एखादी कमांड प्रत्यक्ष न चालवता ती काय करेल हे पाहण्यासाठी तिच्या आधी `--dry-run` लावा, उदाहरणार्थ `ccam --dry-run work`. हे Claude Code सुरू करताना (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) आणि `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` व `doctor --fix` सोबत चालते. फक्त वाचणाऱ्या कमांड्स, जसे `list`, `usage`, `share` आणि `ext`, हा पर्याय स्वीकारतात आणि तसेही काहीच बदलत नाहीत. बाकीच्या कमांड्स, `group rename` सह, चालण्याऐवजी हा पर्याय नाकारतात.

<!-- commands -->

## Claude Code सुरू करणे

`ccam <name|number> [arguments]` त्या खात्याने Claude Code सुरू करते. नावानंतरचे आर्ग्युमेंट्स ccam वाचत नाही; ते त्याच क्रमाने `claude` ला पाठवते, आणि standard input, exit code आणि signals अगदी `claude` प्रमाणेच वागतात. नावानंतर `--` लावणे ऐच्छिक आहे.

खात्याच्या नावाशिवाय पण आर्ग्युमेंट्ससह `ccam`, जसे `ccam -p "summarise"`, सध्याच्या फोल्डरसाठी सेट केलेले खाते वापरते, नाहीतर डीफॉल्ट खाते. ती `ccam folder <name> [folder]` आणि `ccam use <name>` ने सेट करा.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

खात्याच्या नावाआधी `--fallback` हे `-p`/`--print` runs साठी आहे. खाते आपल्या वापर मर्यादेला पोहोचले की ccam तीच कमांड सर्वात जास्त जागा असलेल्या लॉगिन केलेल्या खात्याने पुन्हा चालवते आणि फक्त पूर्ण झालेल्या run चा output छापते. [फोल्डरनुसार खाती](/features/folder) आणि [मर्यादेवर पुन्हा चालवणे](/features/next#run-again-at-the-limit) पाहा.

**Team आणि Enterprise सीट्स.** आपोआप होणारी खाते निवड (`ccam next`, `--fallback`, `ccam pick`, डॅशबोर्डमधील `n` आणि, तुम्ही चालू केल्यावर, [`auto_switch`](/features/next#auto-switch)) एकाच Claude Team किंवा Enterprise organization मधील सीट्सदरम्यानही खाते बदलते, आणि यामुळे Anthropic च्या अटींचे उल्लंघन होऊ शकते. त्यावर अवलंबून राहण्यापूर्वी [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) वाचा आणि तुमच्या संस्थेकडे आणि Anthropic कडे चौकशी करा.

`ccwork` सारखा शॉर्टकट त्याचे खाते याच प्रकारे सुरू करतो: `ccwork --resume` तेच करते जे `ccam work --resume` करते. खालील [शॉर्टकट](/docs/commands#shortcuts) पाहा.

## खाती जोडणे आणि काढणे

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

फोल्डर न देता `ccam import` चालवल्यास ते `~/.claude` आणि `~/.claude-*` शोधते आणि प्रत्येकाबद्दल विचारते; अयशस्वी होणाऱ्या फोल्डरसाठी `Skipped <path>: <reason>` असे सांगितले जाते आणि बाकीचे फोल्डर नोंदवले जातात. एखाद्या फोल्डरचे नाव आधीच वापरात असेल, तर तो फोल्डर `ccam import <folder> --name <name>` ने स्वतंत्रपणे इम्पोर्ट करा. इम्पोर्ट केलेले फोल्डर आहेत त्याच जागी नोंदवले जातात; ccam कधीही लॉगिन हलवत किंवा कॉपी करत नाही.

`ccam login --expired` ज्या खात्यांचे लॉगिन Claude Code ने काढून टाकले आहे त्यांनाही लॉगिन करते. `ccam remove` व्यवस्थापित खात्याचा (ज्याचा फोल्डर ccam ने तयार केला, `ccam import` ने नोंदवलेले नव्हे) फोल्डर trash मध्ये जाण्यापूर्वी ते खाते लॉग आउट करते, दुसऱ्या खात्याकडे तेच लॉगिन असेल तर सोडून. इम्पोर्ट केलेले खाते काढल्यास ते त्याच्या गटांतून बाहेर पडते आणि प्रत्येक शेअर केलेल्या आयटमची एक प्रत त्याच्या फोल्डरमध्ये ठेवली जाते. macOS वर `ccam rename` फोल्डरचे नाव तेच ठेवते, कारण Keychain लॉगिन फोल्डरच्या पाथखाली ठेवते.

`ccam add --api` असे खाते जोडते जे Claude लॉगिनऐवजी API key वापरते; ccam प्रत्येक launch ला key वाचते आणि ती कधीही साठवत नाही. [API key खाती](/features/api) पहा.

एखाद्या खात्याने Claude Code चालू असताना, ते कसेही सुरू केलेले असो (background सत्रांसह), `ccam rename`, `ccam remove` आणि `ccam group` कमांड्स त्या खात्यात बदल करण्यास नकार देतात. आधी session बंद करा. गटाचा एखादा बदल अपूर्ण असतानाही त्या नकार देतात: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` प्रत्येक खाते त्याचा slot क्रमांक, शॉर्टकट, ईमेल आणि फोल्डरसह छापते. त्यातील `groups` array खात्याच्या गटांची यादी देतो, आणि `share` object प्रत्येक घटक (`mcp`, `settings`, `memory`, `extensions`, `history`) तो ज्या गटातून येतो त्या गटाशी, किंवा `local` शी जोडतो.

## शॉर्टकट

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

शॉर्टकट ही तुम्ही स्वतः नाव दिलेली कमांड आहे, जी Claude Code एका खात्याने सुरू करते. `ccam work` प्रमाणेच ती प्रत्येक आर्ग्युमेंट न बदलता Claude Code ला देते, त्यामुळे `ccwork --resume` आणि `ccwork -p "hi"` अपेक्षेप्रमाणे चालतात. ती `~/.local/bin` मधील ccam ची लिंक आहे.

शॉर्टकट ऐच्छिक आहेत. `ccam add` चा विझार्ड शॉर्टकट तयार करायचा का ते विचारतो, `ccam add <name> --shortcut <command>` न विचारता एक तयार करते, आणि `ccam list` प्रत्येक खात्याचा शॉर्टकट दाखवते, किंवा शॉर्टकट नसेल तर `-`.

नावांमध्ये a-z, 0-9, `.`, `-` आणि `_` ही अक्षरे चालतात. नाव ccam ची कमांड, दुसऱ्या खात्याचा शॉर्टकट, किंवा तुमच्या `PATH` वरील दुसऱ्या प्रोग्रामचे नाव असू शकत नाही, आणि ccam ने न बनवलेल्या `~/.local/bin` मधील फाइलचे, जसे ccam ची तुमची स्वतःची प्रत किंवा लिंक, नावही असू शकत नाही. `config.toml` मध्ये हाताने संपादित केलेला शॉर्टकट हे नियम मोडत असेल, `claude` किंवा `ccam` असेल, किंवा दोन खात्यांना दिलेला असेल, तर `config.toml` लोड होत नाही, आणि त्रुटी त्या खात्याचे नाव सांगते.

`ccam rename` शॉर्टकटचे नाव तेच ठेवते, आणि `ccam remove` खात्यासोबत त्याचा शॉर्टकटही काढून टाकते. `ccam shims install` सर्व खात्यांच्या लिंक्स पुन्हा तयार करते, उदाहरणार्थ ccam दुसरीकडे हलवल्यानंतर; एखाद्या शॉर्टकटच्या जागी आधीच फाइल असली तरी इतर शॉर्टकट थांबत नाहीत.

## शेअरिंग

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

गट त्याच्या खात्यांमध्ये आयटम शेअर करतो; एक खाते अनेक गटांमध्ये असू शकते, पण प्रत्येक आयटम फक्त एकाच गटातून येतो. [गटांसह शेअरिंग](/features/share) पाहा.

## आधीचा सेटअप कॉपी करणे

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` एका Claude Code फोल्डरमधून `settings.json`, `CLAUDE.md`, एक्स्टेंशन्स आणि user MCP सर्व्हर एका गटात कॉपी करते, म्हणजे ते आयटम शेअर करणाऱ्या प्रत्येक खात्याला ते मिळतात. फोल्डर फक्त वाचला जातो.

- `--settings`, `--memory`, `--extensions`, `--mcp` आणि `--history` हे पर्याय फक्त तुम्ही नाव दिलेले आयटम कॉपी करतात. त्यांच्याशिवाय, टर्मिनलवर एक चेकलिस्ट तुम्हाला निवडू देते, ज्यात संभाषण इतिहास सोडून सर्व काही निवडलेले असते; इतरत्र संभाषण इतिहास सोडून सर्व काही कॉपी होते.
- `~/.claude` रिकामा असला किंवा नसला तरी `~/.claude.json` मध्ये सर्व्हर असतील तर `--mcp` चालते.
- `-g` किंवा `--group` गट निवडतो.
- `-y` किंवा `--yes` न विचारता कॉपी करतो आणि एखादी फाइल वेगळी असल्यास गटाची प्रत ठेवतो.

MCP साइन-इन्स कॉपी होत नाहीत, कारण कॉपी केलेले साइन-इन ते असलेल्या इतर खात्यांना लॉग आउट करते; ज्या खात्यांनी `/mcp` ने साइन इन करायचे आहे त्यांची नावे योजना सांगते. [आधीचा सेटअप कॉपी करणे](/docs/sharing#copy-an-existing-setup) पाहा.

## MCP सर्व्हर

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` फक्त त्याच खात्याकडे असलेले सर्व्हर त्याच्या गटात हलवते, फक्त एकाच फोल्डरसाठी असलेले सर्व्हरसुद्धा (साधे `claude mcp add`), जे त्यानंतर प्रत्येक फोल्डरमध्ये लोड होतात. ज्या सर्व्हरचे नाव गट आधीच वेगळ्या व्याख्येसाठी वापरत असेल तो वगळला जातो; गटाची प्रत बदलण्यासाठी `--force` जोडा, ccam आधी तिचा बॅकअप घेते.

## वापर आणि खाते बदलणे

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

- `next --auto`, `--fallback` आणि `ccam pick` एकाच पद्धतीने निवडतात: आधी `auto_max_usage` (`config.toml`, डीफॉल्ट 100) खालची खाती, `auto_strategy` (`headroom`, `expiring` किंवा `balance`) नुसार, `--model` ची साप्ताहिक window मोजून. ते 100% वरचे खाते (Claude Code ने भरलेली म्हणून नोंदवलेली window तिच्या reset पर्यंत 100% मोजली जाते) किंवा API key लॉगिन कधीच निवडत नाहीत, आणि ते ज्या खात्यापासून सुरू करतात त्याच्या organization मध्येच राहतात (आणखी पुढे पाहण्यासाठी `--cross-org` किंवा `cross_org = true`).
- `--wait` त्याचे मूल्य `=` नंतर घेते (`--wait=2h`); space ने लिहिल्यास ते एका सूचनेसह नाकारले जाते. काहीच निवडता येत नसेल तर ते code 75 ने बाहेर पडतात, किंवा `--wait` सह पहिले खाते reset होण्याची वाट पाहतात.
- `ccam next --continue` इतिहास शेअर न करणाऱ्या खात्यांमध्येही संभाषण पुढे नेते: तेव्हा Claude Code ते नवीन खात्यात कॉपी करते. `--brief` सह, नवीन खाते त्याऐवजी नवीन संभाषण सुरू करते, ccam ने model शिवाय मागच्या संभाषणाचा लिहिलेल्या सारांशापासून ([`--brief`](/features/next#brief)).
- `ccam <name>` च्या उलट, `ccam next` स्वतःचे पर्याय वाचते, म्हणून Claude Code चे पर्याय `--` नंतर येतात; त्याशिवाय ccam `unknown flag` देऊन थांबते.
- `ccam agent` या फोल्डरमध्ये, ज्या खात्यातून ते बोलावले जाते त्याहून वेगळ्या खात्याने, बोलावणाऱ्या सत्राच्या permission mode आणि effort सह, `claude -p` ने काम चालवते आणि Claude चे शेवटचे उत्तर छापते. त्या खात्याची मर्यादा आल्यावर संभाषण दुसऱ्या खात्याने पुढे चालू राहते ([इतर खात्यांवरील subagents](/features/agent)).
- Claude Code मध्ये plugin चे `/ccam:switch [account]` संभाषण पुढे सोपवते, आणि `ccam hook` हे त्याचे hooks चालवतात.
- `config.toml` मध्ये `auto_switch = true` असेल (Linux आणि macOS; डीफॉल्टनुसार बंद), तर मर्यादेचा hook ccam ने सुरू केलेले interactive सत्रही संपवतो आणि जागा असलेल्या दुसऱ्या खात्यावर ते पुढे चालू ठेवतो, पण फक्त त्याच प्रकारच्या खात्यावर: त्याच organization ची दुसरी सीट, किंवा दुसरा वैयक्तिक प्लॅन ([`auto_switch`](/features/next#auto-switch)).

[पुढच्या खात्यावर जाणे](/features/next), [scripts साठी खाते निवडणे](/features/pick), [इतर खात्यांवरील subagents](/features/agent) आणि [प्रत्येक खात्याची पार्श्वभूमीतील सत्रे](/features/agents) पहा.

## देखभाल

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

`ccam update` ccam बदलण्यापूर्वी रिलीजचा SHA-256 checksum तपासते, आणि `--version v0.2.0` विशिष्ट रिलीज इन्स्टॉल करते. `.deb` पॅकेजमधून इन्स्टॉल केलेली ccam ची प्रत त्याऐवजी apt ने अपडेट करा. तपशीलासाठी [अपडेट](/docs/install#update) पाहा.

`ccam uninstall` काय काढणार आहे त्याची यादी दाखवते आणि आधी विचारते: प्रत्येक खाते `ccam remove` प्रमाणे, शॉर्टकट कमांड्स, ccam चे फोल्डर्स आणि `ccam` कमांड. इम्पोर्ट केलेला फोल्डर आहे तिथेच राहतो. [अनइन्स्टॉल](/docs/install#uninstall) पाहा.

## फाइल्स

| पाथ | मजकूर |
| --- | --- |
| `~/.config/ccam/config.toml` | खाती, गट आणि सेटिंग्ज (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). सेटिंग्ज सुरुवातीला, पहिल्या `[[group]]`, `[[profile]]` किंवा `[[folder]]` च्या वर येतात: table च्या खाली असलेली सेटिंग त्या table ची ठरते, म्हणून ccam प्रत्येक launch ला आणि `ccam doctor` मध्ये तिच्याबद्दल चेतावणी देते, आणि config सेव्ह करणाऱ्या कमांड्स ती हलवेपर्यंत नकार देतात. ccam फाइल सेव्ह करते तेव्हा comments राहत नाहीत |
| `~/.local/share/ccam/profiles/<name>/` | प्रत्येक खात्याचा `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | शेअर केलेल्या सेटिंग्ज, MCP सर्व्हर, एक्स्टेंशन्स आणि इतिहास |
| `~/.local/state/ccam/` | वापराचा cache, मर्यादेच्या नोंदी, locks, नोंदी आणि इतर स्थिती (टेबलखालील यादी पाहा) |
| `~/.local/bin/<shortcut>` | `ccwork` सारख्या शॉर्टकट कमांड्स, ccam च्या लिंक्स |

state फोल्डरमध्ये हे असते:

- वापराचा cache आणि मर्यादेच्या नोंदी (`usage/<name>.json`, `usage/<name>.limit.json`);
- शेवटच्या वापराच्या वेळा;
- session markers आणि locks (`locks/`);
- ccam ने तयार केलेले शॉर्टकट (`shortcuts.json`);
- पूर्ण न झालेला गटाचा बदल (`group-change.json`);
- `/ccam:switch` व मर्यादेचा hook मागे ठेवतात त्या नोंदी (`handoff/<account>/<session>.json`);
- `ccam pick` किंवा `auto_switch` ने प्रत्येक खाते शेवटचे कधी निवडले (`picks/<name>`; एक निवड 2 मिनिटे मोजली जाते);
- `--brief` चे सारांश, 7 दिवस ठेवलेले (`brief/<session>.md`).

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` किंवा `XDG_STATE_HOME` सेट केल्याने या फोल्डर्सची जागा बदलते. Windows वर `~` म्हणजे `%USERPROFILE%`, आणि शॉर्टकट `<shortcut>.exe` असतात.
