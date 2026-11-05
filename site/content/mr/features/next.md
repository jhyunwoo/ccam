---
title: पुढच्या खात्यावर जाणे
description: ccam next दुसऱ्या खात्याने Claude Code उघडते; --auto सर्वाधिक जागेचे खाते निवडते, --continue संभाषण पुढे नेते, --fallback -p jobs पुन्हा चालवते.
updated: 2026-11-05
order: 8
translation: machine
sourceHash: "bc9837b7a6ae"
---

`ccam next` तुम्ही नाव न देता दुसऱ्या खात्याने Claude Code उघडते. एखादे खाते वापर मर्यादेला पोहोचल्यावर ते वापरा. खाते बदलण्याचे चार मार्ग आहेत:

- साधे `ccam next` तुम्ही शेवटी वापरलेल्या खात्यानंतरचे खाते घेते, आणि `--auto` सर्वात जास्त जागा असलेले खाते घेते. `--continue` या फोल्डरमधील शेवटचे संभाषण नवीन खात्यावर पुढे चालू ठेवते, दोन्ही खाती इतिहास शेअर करत असोत वा नसोत.
- Claude Code मध्ये, मर्यादेमुळे सत्र थांबले की [ccam plugin](#inside-claude-code) तुम्हाला सांगते, आणि `/ccam:switch` संभाषण पुढे सोपवते.
- `claude -p` चालवणाऱ्या scripts साठी `ccam --fallback` हा बदल आपोआप करते ([मर्यादेवर पुन्हा चालवणे](#run-again-at-the-limit) पाहा).
- [`auto_switch`](#auto-switch) चालू असेल (Linux आणि macOS), तर Claude Code ने मर्यादेवर थांबवलेले interactive सत्रही ccam हलवते.

**Team आणि Enterprise सीट्स.** आपोआप होणारी खाते निवड (`ccam next`, `--fallback`, `ccam pick`, डॅशबोर्डमधील `n` आणि, तुम्ही चालू केल्यावर, [`auto_switch`](/features/next#auto-switch)) एकाच Claude Team किंवा Enterprise organization मधील सीट्सदरम्यानही खाते बदलते, आणि यामुळे Anthropic च्या अटींचे उल्लंघन होऊ शकते. त्यावर अवलंबून राहण्यापूर्वी [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) वाचा आणि तुमच्या संस्थेकडे आणि Anthropic कडे चौकशी करा.

## वापर

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

ccam हा बदल standard error वर `→ <from> → <to>` असा छापते, `--auto` ने निवड केली असल्यास तिच्या कारणासह, मग `ccam <name>` प्रमाणेच नवीन खात्याने Claude Code सुरू करते. `<from>` म्हणजे तुम्ही शेवटी वापरलेले खाते, किंवा `--continue` सह ते संभाषण ज्या खात्याचे आहे ते खाते. `--` नंतरचे आर्ग्युमेंट्स न बदलता `claude` ला जातात.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--auto` | वापरावरून लॉग इन केलेले खाते निवडा (खाली पहा), 100% वरची खाती वगळून | बंद |
| `--max-usage <1-100>` | `--auto` सह: या वापराखालील खात्यांना प्राधान्य द्या | `auto_max_usage`, अन्यथा 100 |
| `--strategy headroom\|expiring\|balance` | `--auto` सह: सर्वात कमी 5 तासांचा वापर, सर्वात लवकर होणारा 5 तासांचा reset, किंवा दररोज सर्वात जास्त उरणारा साप्ताहिक quota | `auto_strategy`, अन्यथा `headroom` |
| `--cross-org` | `--auto` सह: दुसऱ्या organization ची खातीसुद्धा निवडा (खाली पहा) | `cross_org`, अन्यथा बंद |
| `--wait[=<duration>]` | `--auto` सह: कोणत्याही खात्यात जागा नसेल तर पहिले खाते reset होण्याची वाट पाहा, जास्तीत जास्त इतका वेळ. मूल्य `=` नंतर येते: `--wait 2h` हे `--wait takes its value after '=': --wait=2h` या सूचनेसह नाकारले जाते | बंद; एकटे दिल्यास `5h` |
| `-c`, `--continue` | या फोल्डरमधील शेवटचे संभाषण नवीन खात्यावर पुढे चालू ठेवा (खाली पहा) | बंद |
| `--brief` | `--continue` सह: संपूर्ण संभाषण resume करण्याऐवजी, ccam ने मागच्या संभाषणाचा लिहिलेल्या सारांशापासून नवीन खात्याने नवीन संभाषण सुरू करा ([खाली पहा](#brief)) | बंद |
| `-- <claude args>` | ccam च्या स्वतःच्या आर्ग्युमेंट्सनंतर `claude` ला दिले जाणारे आर्ग्युमेंट्स | काहीही नाही |
| `--dry-run` | समर्थित: Claude Code सुरू न करता किंवा सुरू केल्याची नोंद न करता `→ <from> → <to>` आणि चालवली जाणारी कमांड लाइन छापतो | बंद |

`--auto` आणि `--continue` एकत्र चालतात. `--continue` सह, `--resume` चे आर्ग्युमेंट्स तुमच्या स्वतःच्या आर्ग्युमेंट्सच्या आधी ठेवले जातात.

Claude Code चे स्वतःचे पर्याय, जसे `--dangerously-skip-permissions` किंवा `--model`, `--` नंतर येतात. `ccam <name>` नावानंतरचे सर्व काही पुढे देते, पण `ccam next` त्याचे आर्ग्युमेंट्स स्वतः वाचते, म्हणून `--` शिवाय ते `ccam: unknown flag: --dangerously-skip-permissions` देऊन थांबते.

## हे कसे काम करते

ccam मार्फत प्रत्येक वेळी सुरू करताना वेळ `~/.local/state/ccam/last-used/<name>` मध्ये लिहिली जाते. `ccam next` या फाइल्स वाचते, कोणती खाती लॉग इन आहेत ते तपासते, आणि निवड करते.

1. ccam `config.toml` वरून प्रत्येक खात्यासाठी एक `next.Candidate` तयार करते (`cli.candidates`): खाते, ते शेवटी सुरू केल्याची वेळ (`launcher.LastUsed`) आणि त्याला चालणारे लॉगिन आहे का (`creds.Has`: `.credentials.json` मध्ये, किंवा macOS वर Keychain मध्ये, असलेले Claude लॉगिन, जे संपूर्णपणे कालबाह्य झालेले नाही किंवा Claude Code ने काढून टाकलेले नाही). त्याऐवजी API key वापरणाऱ्या फोल्डरला (Console लॉगिन, किंवा `ccam add --api` ने जोडलेले खाते) plan मर्यादा नसतात आणि ते कधीच निवडले जात नाही.
2. `--auto` शिवाय, ccam खाती slot क्रमांकानुसार लावते आणि सर्वात अलीकडे सुरू केलेले खाते शोधते (`next.Rotate`). ज्या खात्यावरून तुम्ही दुसरीकडे जात आहात ते हेच. अजून एकही खाते सुरू केले नसेल, तर ते सर्वात लहान slot असते.
3. तिथून ते slot क्रमाने पुढे जाते, शेवटच्या slot नंतर पुन्हा पहिल्यावर येते, आणि लॉग इन असलेले पहिले खाते परत देते.
4. ccam निवडलेले खाते अगदी `ccam <name>` प्रमाणे तयार करून सुरू करते (`launcher.LaunchWith`): ते खात्याच्या लिंक्स तपासते, MCP शेअरिंग लागू करते, session नोंदवते आणि `syscall.Exec` मार्फत ccam च्या जागी `claude` ठेवते.

शोध शेवटी वापरलेल्या खात्यापासून सुरू होतो आणि त्यानंतरच्या पहिल्या लॉग इन असलेल्या खात्यावर थांबतो:

<!-- motion:next-rotate -->

`work` सर्वात अलीकडे सुरू केले होते, म्हणून शोध त्याच्यानंतर सुरू होतो. `side` ला लॉगिन नाही आणि ते वगळले जाते, आणि `lab` निवडले जाते. `lab` सुद्धा लॉग आउट असते, तर शोध फिरून `personal` वर आला असता (आकृतीतील `% len`). rotation कधीच तुम्ही जिथून सुरू केले ते खाते परत देत नाही: दुसरे लॉग इन असलेले खाते नसल्यास ते त्रुटीसह थांबते.

### --auto

`--auto` slot क्रमाऐवजी वापरावरून ठरवते. ccam लॉग इन असलेली खाती ठेवते, `ccam usage` प्रमाणेच त्या सर्वांचा वापर एकाच वेळी विचारते (`Fetcher.FetchAll`), आणि खालील नियमांनी निवडते (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

खाते त्याच्या 5 तासांच्या आणि साप्ताहिक windows वरून ठरवले जाते, आणि Claude Code च्या arguments मध्ये (`--` नंतर) `--model` असेल तर त्या model च्या साप्ताहिक window वरूनही (नावात `opus`, `sonnet`, `haiku` किंवा `fable`). ccam खात्यांची तीन स्तरांत विभागणी करते (`Policy.Choose()`):

1. **जागा आहे**: प्रत्येक window `auto_max_usage` (`config.toml`, डीफॉल्ट 100) खाली. त्यांच्यातून strategy निवडते: `headroom` (डीफॉल्ट) सर्वात कमी 5 तासांच्या वापरापासून 3 टक्केवारी अंकांच्या आत असलेल्या खात्यांपैकी ज्याच्या आठवड्याचा सर्वात कमी भाग वापरला गेला आहे ते खाते घेते, म्हणजे आठवड्याचा कमी भाग वापरला गेला असेल तर 12% वरचे खाते 10% वरच्या खात्यावर मात करते, पण 14% वरचे खाते करत नाही; `expiring` सर्वात लवकर reset होणारी 5 तासांची window घेते, म्हणजे तिचा उरलेला quota वाया जाण्याआधी वापरला जातो; `balance` आठवडा reset होईपर्यंत दररोज सर्वात जास्त साप्ताहिक quota उरलेले खाते घेते, आणि चालू असलेले प्रत्येक ccam सत्र त्याच्या विरुद्ध मोजते, त्यामुळे साप्ताहिक quota न वापरता राहत नाही आणि समांतर सत्रे विखुरली जातात.
2. **मर्यादेवर**: 100% खाली, पण `auto_max_usage` वर किंवा त्याहून जास्त. कोणत्याही खात्यात जागा नसेल तेव्हाच वापरली जातात; त्यांच्यापैकी `headroom` जे खाते घेईल ते आधी येते.
3. **भरलेली**: एखादी window 100% वर, किंवा अज्ञात 5 तासांचा वापर. कधीच निवडली जात नाहीत.

Claude Code ने एखाद्या window मध्ये request नाकारली असेल, तर ती window तिच्या reset पर्यंत 100% मोजली जाते: [ccam plugin](#inside-claude-code) चा hook ती नोंदवतो (`~/.local/state/ccam/usage/<name>.limit.json`), साठवलेले आकडे काहीही सांगोत. त्यामुळे नुकतीच मर्यादा गाठलेले खाते त्याच्या reset पूर्वी पुन्हा निवडले जात नाही, त्याचे आकडे जुने असले तरीही.

सर्व खात्यांसाठीच्या सेटिंग्ज `config.toml` च्या सुरुवातीला, पहिल्या `[[group]]`, `[[profile]]` किंवा `[[folder]]` च्या वर येतात. त्यांपैकी एखाद्याच्या खाली असलेली सेटिंग त्या table ची ठरते: ccam प्रत्येक launch ला तिच्याबद्दल इशारा देते, आणि config सेव्ह करणाऱ्या कमांड्स ती हलवेपर्यंत नकार देतात.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organizations नुसार मर्यादा.** आपोआप होणाऱ्या निवडी तुम्ही सोडत असलेल्या खात्याच्याच प्रकारच्या खात्यांमध्ये राहतात: त्याच Team किंवा Enterprise organization मधील सीट्स, किंवा personal plans (Pro, Max), त्यामुळे कंपनीच्या खात्यावरील काम खासगी खात्यावर पुढे चालू होत नाही, किंवा उलटही नाही. ccam ला प्रत्येक खात्याची organization त्याच्या लॉगिनवरून कळते (`ccam login` ती नोंदवते; जुनी लॉगिन `.claude.json` मधून वाचली जातात); ज्या खात्याची organization ccam सांगू शकत नाही ते कोणत्याही खात्याशी जुळवता येते. एकाच Claude खात्यात लॉग इन असलेले दोन फोल्डर त्याच्या मर्यादा शेअर करतात, त्यामुळे एक दुसऱ्याची जागा कधीच घेत नाही. `--cross-org`, किंवा `config.toml` मधील `cross_org = true`, organization चा नियम उठवतो, पण [`/ccam:switch`](#inside-claude-code) नंतर दुसऱ्या organization मध्ये होणारा बदल तरीही तुमच्या होकाराची वाट पाहतो.

काहीच निवडता येत नसेल तर ccam code 75 (`EX_TEMPFAIL`) ने बाहेर पडते, जेणेकरून script "नंतर पुन्हा प्रयत्न करा" आणि अपयश यांतील फरक ओळखू शकेल. `--wait` सह, ते त्याऐवजी पहिल्या भरलेल्या खात्याची अडवणारी window reset होईपर्यंत थांबते (reset time नंतर 30 सेकंदांनी, जेव्हा endpoint नवीन window दाखवतो), `ccam: no account has room; waiting until 15:10` छापते, आणि पुन्हा निवडते; कोणतेही खाते कधी reset होणार ते सांगत नसेल तर ते 10 मिनिटांनी पुन्हा पाहते. वाट पाहणे संपण्याआधी काहीच मोकळे झाले नाही तर ते 75 ने बाहेर पडते.

हा शोध `ccam usage` चे नियम पाळतो: 60 सेकंदांपेक्षा कमी जुने आकडे पुन्हा वापरले जातात, प्रत्येक खात्याचा शोध 3 सेकंदांनंतर timeout होतो, उत्तर न देणारे किंवा access token कालबाह्य झालेले खाते त्याच्या साठवलेल्या आकड्यांवरून ठरवले जाते, आणि मर्यादेची नोंद असेल तर ती त्या आकड्यांची जागा घेते. reset time उलटून गेलेली साठवलेली window 0% मोजली जाते, त्यामुळे ज्या खात्याची window दरम्यान रीसेट झाली आहे ते वगळले जात नाही. `--auto` कधीच tokens नूतनीकरण करत नाही. साध्या rotation च्या विपरीत, `--auto` तुम्ही शेवटी वापरलेले खाते निवडू शकते, जर त्यात अजूनही सर्वात जास्त जागा असेल.

### --continue

1. ccam प्रत्येक खात्यात या फोल्डरचे शेवटचे संभाषण शोधते: `projects/<folder>/` मधील `<session>.jsonl` फाइल्सपैकी सर्वात नवीन prompt किंवा उत्तर असलेली फाइल (फोल्डरचे नाव म्हणजे links सोडवलेला कामाचा फोल्डर, ज्यातील अक्षर किंवा अंक नसलेले प्रत्येक चिन्ह `-` मध्ये बदललेले असते, जसे Claude Code नाव देते). prompt किंवा उत्तर नसलेली सत्रे, जसे फक्त `/mcp`, `/usage` किंवा `/status` चालवलेले सत्र, किंवा रिकामी फाइल, वगळली जातात. एखाद्या सत्राने नोंद ठेवली असेल (वापर मर्यादा किंवा `/ccam:switch`, खाली पहा) आणि त्या नोंदीनंतर 2 मिनिटांहून अधिक काळाने कोणत्याही संभाषणात prompt किंवा उत्तर आलेले नसेल, तर त्याऐवजी नेमके तेच सत्र वापरले जाते; `/exit` वेळी Claude Code लिहित असलेल्या ओळी मोजल्या जात नाहीत.
2. ते संभाषण ज्या खात्याचे आहे त्याच खात्यावरून तुम्ही बदलत असता, आणि नवीन खाते वरीलप्रमाणे, त्याच organization मधून निवडले जाते. अनेक खाती तो इतिहास शेअर करत असतील, तर ते संभाषण अशा खात्याचे मानले जाते ज्याची `.claude.json` सांगते की त्याने ते संभाषण शेवटी चालवले (`lastSessionId`), नाहीतर शेवटी वापरलेल्या खात्याचे.
3. नवीन खाते तो इतिहास शेअर करत असेल (तोच गट फोल्डर), तर ccam `claude --resume <session id>` सुरू करते, म्हणजे नेमके तेच संभाषण (`--continue` फोल्डरमधील जे सर्वात नवीन असेल ते उघडेल).
4. नाहीतर ते `claude --resume <the conversation's file> --fork-session` सुरू करते: Claude Code संपूर्ण संभाषण नवीन session ID खाली नवीन खात्याच्या स्वतःच्या इतिहासात कॉपी करते, आणि मूळ संभाषण होते तसेच राहते. त्यानंतर नवीन खात्याचा `--continue` ते शोधतो.
5. तुम्ही `--` नंतर `--agent` दिला नसेल, तर संभाषण ज्या agent म्हणून चालते तो agent देखील ccam पुढे देते: फोल्डरच्या `.claude/settings.local.json` किंवा `.claude/settings.json` मधील, किंवा खात्याच्या `settings.json` मधील `agent` सेटिंग, नाहीतर संभाषणाची शेवटची `agent-setting` ओळ ज्या agent चे नाव देते तो. संभाषण कॉपी करताना (पायरी 4) Claude Code agent परत आणत नाही, आणि ते आपल्या डीफॉल्ट agent म्हणून, सर्व tools सह, पुढे चालेल. नवीन खाते जुन्या खात्याचे agents, plugins आणि सेटिंग्ज शेअर करत नसेल (गटाद्वारे तेच `agents` आणि `plugins` फोल्डर आणि `settings.json`), तर ccam बदलाविषयीच्या ओळीनंतर एक ओळ जोडते: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

संभाषण किती मोठे आहे ते ccam सांगते: नवीन खाते आपल्या पहिल्या उत्तराबरोबर ते सगळे पुन्हा वाचते. Prompt caches कधीही organizations मध्ये शेअर होत नाहीत, आणि प्रत्येक personal plan ही स्वतःची वेगळी organization असते, त्यामुळे दोन्ही खाती वेगवेगळ्या organizations ची असतील तेव्हा ccam `prompt caches are not shared between organizations, so the first reply reads it all again` जोडते. जुने खाते ज्या MCP सर्व्हरमध्ये साइन इन आहे त्यांत नवीन खाते साइन इन नसेल, तर ccam त्यांची नावे सांगते: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. गटात `history` शेअर केल्याने ही प्रत टळते, एकच संभाषण राहते, आणि `/rewind` ज्या checkpoints मधून code परत आणते ते (`file-history/`) सुद्धा शेअर होतात.

### --brief

`--continue` संपूर्ण संभाषण resume करते, आणि नवीन खाते आपल्या पहिल्या उत्तराबरोबर ते सगळे पुन्हा वाचते. `ccam next --continue --brief` त्याऐवजी नवीन संभाषण सुरू करते, ccam ने कोणत्याही model ला न विचारता जुन्या संभाषणाचा लिहिलेल्या सारांशापासून:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

सारांशात, याच क्रमाने आणि फक्त सांगण्यासारखे काही असेल तिथे, हे असते:

1. संभाषणाचे शीर्षक;
2. त्याच्या transcript चा path, एखादा तपशील हवा असेल तेव्हा Claude ने शोधण्यासाठी;
3. तुमची पहिली विनंती, चारपेक्षा जास्त विनंत्या असतील तेव्हा;
4. Claude Code ने संभाषण शेवटचे compact केले तेव्हा लिहिलेला सारांश;
5. तुमच्या शेवटच्या चार विनंत्या, `!` commands आणि slash commands तुम्ही टाइप केल्या तशाच;
6. Claude चे शेवटचे उत्तर;
7. todo list;
8. Claude ने शेवटी बदललेल्या 10 फाइल्स;
9. फोल्डरचा `git status --short --branch`, 800 characters पर्यंत.

todo list मध्ये जास्तीत जास्त 10 items राहतात: पूर्ण झालेले items आधी काढले जातात, मग सर्वात पुढच्या पायऱ्या, त्यामुळे सध्याची पायरी आणि त्यानंतरच्या पायऱ्या राहतात. लांब भाग मध्यभागी कापले जातात, आणि संपूर्ण सारांश 16,000 characters च्या आत राहतो.

ccam तो `~/.local/state/ccam/brief/<session>.md` म्हणून साठवते, जो फक्त तुम्हीच वाचू शकता, आणि पुढचे `--brief` 7 दिवसांपेक्षा जुने सारांश काढून टाकते. Claude Code एका prompt ने सुरू होते, आणि त्यात उल्लेख केलेली फाइल जोडते: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` तुम्ही `--` नंतर दिलेला prompt शेवटच्या वाक्याची जागा घेतो (`Then: <your prompt>`), आणि `--` नंतरचे तुमचे इतर पर्याय त्याच्या मागे येतात.

Claude Code उल्लेख `"` वर संपवते आणि `#` वर कापते, त्यामुळे सारांशाच्या path मध्ये यापैकी एखादे असेल, तर ccam कार्यरत फोल्डरच्या सापेक्ष path ने त्याचा उल्लेख करते; त्या path मध्येही असेल, तर prompt फाइलचे नाव सांगतो, आणि Claude ती आपल्या Read tool ने वाचते, जे आधी विचारते कारण फाइल कार्यरत फोल्डर्सच्या बाहेर आहे (Manual mode मधील `-p` run मध्ये Claude Code वाचणे नाकारते आणि Claude सारांशाशिवाय पुढे जाते).

`ccam --dry-run next --continue --brief` सारांश छापते आणि काहीही साठवत नाही. जुने संभाषण होते तसेच राहते, आणि नवीन संभाषणाला फक्त सारांशात असलेलेच माहीत असते, त्यामुळे तुम्ही तिथे कसे पोहोचलात याच्या प्रत्येक तपशीलापेक्षा काम जास्त महत्त्वाचे असेल तेव्हा `--brief` वापरा. `--continue` शिवाय ccam `--brief works with --continue` देऊन थांबते; संभाषण वाचता आले नाही किंवा त्यात सारांश करण्यासारखे काहीच सापडले नाही, तर ते थांबते आणि सांगते की `--brief` शिवायचे `ccam next --continue` ते संपूर्ण resume करते. [relay](#relay) आणि [`auto_switch`](#auto-switch) नेहमी संपूर्ण संभाषण resume करतात. Claude Code चा स्वतःचा `--brief` पर्याय, त्याच्या इतर पर्यायांप्रमाणे, `--` नंतर येतो.

## Claude Code मध्ये

[ccam plugin](/docs/install#claude-code-plugin) दोन hooks जोडते (Claude Code 2.1.78 किंवा नंतरचे):

- एखादा turn वापर मर्यादेवर थांबला की ते संभाषणाची नोंद ठेवते आणि सूचना देऊ शकणाऱ्या terminals मध्ये (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty) desktop सूचना दाखवते: ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  Claude Code वापर मर्यादा म्हणून मोजत असलेला प्रत्येक संदेश त्याला ओळखता येतो: 5 तासांची, साप्ताहिक, साधी आणि model ची मर्यादा (`You've hit your session limit`, `You've hit your limit`, model च्या मुभेसाठी `You've reached your …`), usage credits आणि खर्चाच्या मर्यादा, आणि खाते थांबवलेली organization.

  Claude Code ने नाकारलेली request संभाषणात लिहिली असेल, तर hook भरलेली window तिच्या reset पर्यंत नोंदवतो ([मर्यादेची नोंद](/features/usage#the-limit-record)), organization ने लावलेला थांबा सोडून, कारण तिचे admins तो कधीही उठवू शकतात. सर्व्हरचे थोड्या वेळाचे throttling किंवा API key ची rate limit ही वापर मर्यादा नाही, आणि तिला हात लावला जात नाही, तसेच कोणीही समोर नसताना चालणारी सत्रे (`-p` runs आणि background सत्रे) यांनाही नाही. वापर मर्यादेवर थांबलेल्या subagent किंवा workflow agent चे थांबणे हे त्याच्या सत्राचेच थांबणे मानले जाते, मुख्य संभाषण फक्त त्याची वाट पाहत असले तरीही: Claude Code ती नाकारलेली request सत्राच्या `subagents` फोल्डरमधील त्या agent च्या स्वतःच्या फाइलमध्ये लिहिते, आणि hook ती तिथूनच वाचतो. अनेक agents एकाच वेळी थांबले, तर सत्राला एकच नोंद आणि एकच सूचना मिळते.

  Claude Code 2.1.234 आणि नंतरचे सत्रात थांबून राहते आणि 5 तासांच्या reset नंतर आपोआप पुढे चालू ठेवते; सूचना तुम्हाला त्याऐवजी खाते बदलण्याचा पर्याय देते, आणि [`auto_switch`](#auto-switch) चालू असेल (Linux आणि macOS), तर दुसऱ्या खात्यात जागा असताना ccam तुमच्यासाठी खाते बदलते. दुसरे कोणतेही खाते जबाबदारी घेऊ शकत नसेल, कारण दुसरे कोणतेही खाते त्याच Team किंवा Enterprise organization चे नाही, किंवा दुसरा कोणताही वैयक्तिक प्लॅन नाही (`cross_org` चालू असल्यास कोणतेही दुसरे Claude खाते घेऊ शकते), तर सूचनेच्या शेवटी त्याऐवजी `You have no other personal plans, so ccam cannot move this conversation.` किंवा `You have no other accounts in Acme, so ccam cannot move this conversation.` येते.
- `/ccam:switch [account]` संभाषणाची नोंद ठेवते आणि कोणतीही request होण्यापूर्वी turn संपवते: `/exit` टाइप करा, मग ते तुम्ही नाव दिलेल्या खात्यावर (किंवा पुढच्या खात्यावर) पुढे चालू ठेवण्यासाठी `ccam next --continue` चालवा. अजून एकही संभाषण नसलेल्या सत्रात ते सांगते की पुढे चालू ठेवण्यासारखे काहीच नाही. ते फक्त तुम्ही टाइप केल्यावरच चालते: plugin 0.3.2 पासून Claude `/ccam:switch` चालवू शकत नाही, आणि Claude ने शेड्यूल केलेले काम (CronCreate, ScheduleWakeup किंवा `.claude/scheduled_tasks.json`) ते Claude ला साध्या मजकुराच्या रूपात देते.

  तुम्ही नाव दिलेले खाते संभाषणाच्या organization बाहेरचे असेल (दुसऱ्या Team किंवा Enterprise organization ची सीट, सीटऐवजी वैयक्तिक प्लॅन किंवा उलट, किंवा कोणत्याही बाजूला API key खाते), तर `cross_org` काहीही सांगो, `ccam next --continue` आधी विचारते: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. तुम्ही नाव दिले नाही तर ते स्वतः निवडलेल्या खात्यापूर्वीही हेच विचारते, जे `cross_org = true` असल्यास ते दुसऱ्या organization मधून घेऊ शकते: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. फक्त `y` किंवा `yes` ते हलवते; Linux आणि macOS वर, प्रश्नापूर्वी टाइप केलेल्या keys आधी टाकून दिल्या जातात. इतर कोणतेही उत्तर, किंवा विचारण्यासाठी terminal नसणे, संभाषण आणि नोंद जिथे आहेत तिथेच ठेवते: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (स्वतः निवडल्यानंतर, `ccam next --continue` पुन्हा निवडते). आपोआप होणाऱ्या निवडींप्रमाणे, ज्या खात्याची organization ccam सांगू शकत नाही ते कोणत्याही खात्याशी जुळते.

ही नोंद `~/.local/state/ccam/handoff/<account>/<session>.json` मध्ये असते, प्रत्येक सत्रासाठी एक फाइल, त्यामुळे एकाच खात्याची दोन सत्रे एकमेकांची नोंद कधीही overwrite करत नाहीत. ती 12 तास ठेवली जाते आणि पुढचे सत्र सुरू होणार असताना एकदाच वापरली जाते; सुरू करणे अयशस्वी झाले तर ती पुढच्या प्रयत्नासाठी तशीच राहते. ती सत्र ज्या फोल्डरमध्ये सुरू झाले त्या फोल्डरची असते, Claude Code च्या shell मध्ये `cd` केल्यानंतरही. `ccam rename` तिच्यातील नावे अद्ययावत करते; तिच्यात नाव असलेले खाते नंतर काढून टाकले असेल, तर ccam तसे सांगते आणि नेहमीप्रमाणे निवडते.

### Relay

`config.toml` च्या सुरुवातीला `relay = true` असेल, तर ccam स्वतःची जागा Claude Code ला देण्याऐवजी interactive सत्रांसोबत चालू राहते, आणि `/ccam:switch` नंतर त्या सत्राचे संभाषण आपोआप पुढे चालू ठेवते (दुसऱ्या सत्राने ठेवलेल्या नोंदीवर कृती केली जात नाही): hook Claude Code संपवतो, आणि ccam त्या सत्राच्या नोंदीसह `ccam next --continue` चालवते. खात्याचे नाव दिले नसेल, तर ते संभाषणाचे model मोजून सर्वात जास्त जागा असलेले खाते घेते, आणि `cross_org = true` मुळे ते दुसऱ्या organization चे खाते असू शकत असेल, तर संभाषण तिथे हलवण्यापूर्वी ते वरचा प्रश्न विचारते. Windows वर hook Claude Code संपवू शकत नाही: `/exit` टाइप करा, आणि ccam ते पुढे चालू ठेवते. संभाषणाच्या organization बाहेरच्या खात्यासाठीही hook ते संपवत नाही: तो `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).` असे सांगतो, आणि तुम्ही `/exit` टाइप केल्यावर ccam वरचा प्रश्न विचारते. त्यासाठी प्रत्येक सत्रामागे वाट पाहणारी एक ccam process लागते, म्हणूनच हे डीफॉल्टनुसार बंद असते.

relay ने पुढे नेलेले सत्र, संभाषण resume केल्याने परत न येणाऱ्या गोष्टी ठेवते: `--add-dir`, `--mcp-config`, tool lists, `--settings`, `--plugin-dir` आणि system prompt चे पर्याय यांसारखे पर्याय, संभाषणाचा permission mode, सत्राचा effort level आणि संभाषण ज्या agent म्हणून चालते तो agent ([`--continue`](#continue) पहा). `--model` आणि `--name` Claude Code स्वतः संभाषणातून परत आणते. इतर कोणतेही (`--worktree`, `--debug`, `--remote-control`, ccam ला माहीत नसलेले पर्याय) `ccam: not carried over: …` मध्ये नमूद केले जातात. हा बदल शेअर केलेल्या इतिहासाचा प्रश्न कधीच विचारत नाही, आणि तो चालू असताना तुम्ही टाइप केलेले टाकून दिले जाते; त्याचा एकमेव प्रश्न म्हणजे दुसऱ्या organization च्या खात्यापूर्वी विचारला जाणारा प्रश्न.

`/ccam:switch` ने सांगितल्याप्रमाणे Claude Code संपले तरच relay संभाषण पुढे चालू ठेवतो: तुम्ही आधीच `/exit` टाइप केले असेल, ते संपत असताना Ctrl+C दाबले असेल, किंवा ccam ला hangup किंवा दुसरा signal मिळाला असेल, तेव्हा नाही, आणि terminal नाहीसे झाले असेल तेव्हाही नाही. तेव्हा नोंद तशीच राहते, आणि `ccam next --continue` ते नंतर पुढे चालू ठेवते. ccam च्या जुन्या आवृत्तीच्या relay खाली सुरू झालेले सत्र, किंवा कोणीही समोर नसताना चालणारे सत्र (`-p` runs आणि background सत्रे), hook संपवत नाही; त्याला त्याऐवजी `/exit` चा संदेश मिळतो. `claude remote-control` relay शिवाय चालते.

### `auto_switch`

`config.toml` च्या सुरुवातीला `auto_switch = true` असेल (Linux आणि macOS), तर Claude Code सत्र वापर मर्यादेवर थांबवते तेव्हा ccam ते स्वतः दुसऱ्या खात्यावर हलवते: plugin चा hook जागा असलेले दुसरे खाते निवडतो आणि Claude Code संपवतो, आणि सत्रासोबत चालू असलेले ccam संभाषण तिथे पुढे चालू ठेवते, `/ccam:switch` नंतर [relay](#relay) करते त्याच प्रकारे. हे डीफॉल्टनुसार बंद आहे आणि experimental आहे.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam सत्र फक्त तेव्हाच हलवते जेव्हा या सर्व गोष्टी खऱ्या असतात:

1. ccam ने ते terminal मध्ये सुरू केले आहे आणि त्याच्यासोबत चालू आहे. `auto_switch` [relay](#relay) चालू करते, त्यामुळे ते सेट केल्यानंतर तुम्ही सुरू केलेल्या interactive सत्रांसाठी हे खरे असते. आधी, किंवा जुन्या ccam खाली सुरू झालेल्या सत्राला अशी सूचना मिळते जिच्या शेवटी ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` असते. कोणीही समोर नसताना चालणारी सत्रे (`-p` runs आणि background सत्रे), `claude remote-control` आणि ccam शिवाय सुरू केलेले Claude Code कधीच हलवले जात नाहीत.
2. खात्यात [ccam plugin](/docs/install#claude-code-plugin) चालू आहे आणि ते Claude Code 2.1.214 किंवा नंतरचे चालवते, जे hook ला आपला process ID देते (`CLAUDE_PID`).
3. Claude Code ने नाकारलेली request, ती reset होण्याच्या वेळेसह, त्या खात्याच्या 5 तासांच्या, साप्ताहिक, Opus किंवा Sonnet मर्यादेसाठी संभाषणात लिहिली आहे. खर्चाची मर्यादा, usage credits संपणे, organization ने लावलेला थांबा, Fable ची मुभा (`You've reached your Fable limit.`), सर्व्हरचे throttling आणि API key ची rate limit यांमुळे सत्र कधीच हलत नाही. subagent किंवा workflow agent ची मर्यादा मात्र गणली जाते, मुख्य संभाषण फक्त त्याची वाट पाहत असले तरीही; ती Opus किंवा Sonnet ची साप्ताहिक मर्यादा असेल, तर संभाषण स्वतः त्याच model वर चालत असेल तरच सत्र हलते.
4. ccam ला खात्याचा plan माहीत आहे, त्याच्या `.claude.json` मधून किंवा `ccam login` ने नोंदवलेल्या माहितीतून. नसेल, तर सूचनेच्या शेवटी `ccam does not know personal's plan, so auto_switch leaves this session.` असते.
5. सत्र Remote Control ला जोडलेले नाही. दुसऱ्या खात्याने ते तुम्ही ज्या app मधून चालवता त्याला पुन्हा जोडले जाणार नाही, त्यामुळे ते तसेच राहते, आणि सूचना तसे सांगते.
6. मर्यादा आतापासून 15 मिनिटांपेक्षा जास्त वेळाने reset होते. यापेक्षा कमी वाट पाहायची असेल तर ccam सत्र Claude Code वर सोडते: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. मर्यादेनंतर तुम्ही कोणताही prompt टाइप केलेला नाही.
8. दुसऱ्या खात्यात जागा आहे. ccam [`--auto`](#auto) प्रमाणे निवडते (`auto_max_usage`, `auto_strategy`, संभाषणाचे model), तुम्ही सोडत असलेल्या खात्याच्याच प्रकारच्या खात्यांमधून: त्याच Team किंवा Enterprise organization ची दुसरी सीट, किंवा दुसरा वैयक्तिक प्लॅन (Pro, Max). `cross_org` आणि `--cross-org` इथे लागू होत नाहीत. खात्याने:
   - Claude Code नूतनीकरण करू शकेल अशा लॉगिनने लॉग इन असले पाहिजे;
   - त्याच Claude खात्यात लॉग इन नसले पाहिजे;
   - API key वापरत नसले पाहिजे;
   - ccam ला माहीत असलेला plan असला पाहिजे;
   - ccam च्या स्वतःच्या वापर-शोधात दिसते तितपत, usage credits चालू नसले पाहिजेत;
   - Claude Code मध्ये हा फोल्डर trust केलेला असला पाहिजे;
   - संभाषण एखाद्या agent म्हणून चालत असेल (`agent` सेटिंग किंवा संभाषण एखाद्या agent चे नाव देते), तर तुम्ही सोडत असलेल्या खात्याचे agents शेअर केलेले असले पाहिजेत: त्याचे `agents` आणि `plugins` फोल्डर आणि त्याची `settings.json` प्रत्येकी अगदी तीच असली पाहिजे (गटाद्वारे), किंवा दोन्हीकडे नसली पाहिजे.

   या शेवटच्या तीन अटींमुळे जागा असलेले प्रत्येक खाते वगळले गेले, तर सूचना त्यांची नावे देते: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` किंवा `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` [मर्यादेची नोंद](/features/usage#the-limit-record) असलेली window भरलेली मानली जाते, आणि मागील 2 मिनिटांत दुसरे सत्र ज्या खात्यावर गेले ते खाते इतरांनंतर येते. निवड सुमारे 6 सेकंदांत व्हायला हवी; वापर-शोधांना जास्त वेळ लागला, तर सत्र तसेच राहते. `usage_source = "claude"` सह `auto_switch` फक्त साठवलेल्या आकड्यांवरून ठरवते, आणि usage credits चालू आहेत की नाही हे ccam ला Claude Code च्या शोधातून किंवा स्टेटस लाइनमधून कळू शकत नाही, त्यामुळे ते त्यांना बंद आहेत असे गृहीत धरते.

मग hook Claude Code संपवतो, आणि त्याच्यासोबत चालू असलेले ccam सूचना देऊ शकणाऱ्या terminals मध्ये desktop सूचना दाखवते, जसे `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` ते वरच्या उदाहरण output ची पहिली ओळ छापते आणि निवडलेल्या खात्याने त्या सत्रासाठी `ccam next --continue` चालवते, जे त्याच्या इतर दोन ओळी छापते आणि सत्राचे पर्याय, permission mode, effort level आणि agent यांसह Claude Code सुरू करते, [relay](#relay) करते तसे. नवीन सत्राचा पहिला संदेश एक ठरलेला मजकूर असतो, ccam तुमच्या वतीने पाठवत असलेली ही एकमेव गोष्ट (`/ccam:switch` नंतरचा बदल काहीही पाठवत नाही):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

subagent च्या मर्यादेनंतर तो हा असतो, कारण subagents, workflows आणि background कमांड्स जुन्या सत्रासोबतच थांबले:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

ते खातेही मर्यादेला पोहोचले, तर तिथेही असेच होते. मर्यादेच्या नोंदी ही साखळी एखाद्या खात्याच्या reset पूर्वी त्याच्याकडे परत जाऊ देत नाहीत, त्यामुळे प्रत्येक खाते भरलेले असेल तेव्हा सत्र शेवटच्या खात्यावर राहते.

संभाषण एखाद्या agent म्हणून चालते हे hook ला `agent` सेटिंगवरून कळते, किंवा Claude Code वेळोवेळी आणि सत्र संपताना संभाषणात लिहित असलेल्या, agent चे नाव देणाऱ्या ओळीवरून. आतापर्यंत फक्त सत्राचा `--agent` त्याचे नाव देत असेल, आणि निवडलेले खाते agents शेअर करत नसेल, तर `ccam next --continue` काहीही सुरू करत नाही आणि नोंद तशीच ठेवते, आणि दुसरी desktop सूचना तसे सांगते: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

दुसऱ्या कोणत्याही खात्यात जागा नसेल, तर काहीही हलत नाही: Claude Code नेहमीप्रमाणे मर्यादा हाताळते, सत्रात थांबून किंवा त्याचे पर्याय देऊन, आणि कोणते खाते आधी मोकळे होते ते सूचना सांगते, उदाहरणार्थ `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` सीट्ससाठी ती `None of your other accounts in Acme has room now` असे सांगते. सुरू झालेला बदल थांबवण्यासाठी, Claude Code बंद होत असताना Ctrl+C दाबा: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

बदलाचा खर्च असतो. Claude Code पुन्हा सुरू होते, ज्याला काही सेकंद लागतात, आणि Claude ने सुरू केलेले background shells, dev servers आणि agents जुन्या सत्राबरोबर थांबतात. नवीन खाते आपल्या पहिल्या उत्तराबरोबर संपूर्ण संभाषण पुन्हा वाचते, [`--continue`](#continue) प्रमाणे, आणि ते किती मोठे आहे ते ccam सांगते; बदल कधीही [`--brief`](#brief) वापरत नाही.

`auto_switch` Claude Code ने मर्यादा नोंदवण्यापूर्वी कधीच खाते बदलत नाही, background मध्ये तुमच्या वापरावर कधीच लक्ष ठेवत नाही, आणि लॉगिन कधीच कॉपी, नूतनीकरण किंवा पुढे देत नाही. Windows वर त्याकडे दुर्लक्ष केले जाते, आणि `ccam doctor` तसे सांगते; तिथे `relay = true` सह `/ccam:switch` आणि `/exit` तरीही संभाषण दुसऱ्या खात्यावर हलवतात.

**Anthropic च्या अटी.** Team आणि Enterprise सीट्सवर `auto_switch` तुमचे काम त्याच organization च्या दुसऱ्या सीटवर हलवते. Anthropic च्या मदत केंद्रात म्हटले आहे की "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team प्लॅनमधील प्रत्येक सीट नावानिशी एका व्यक्तीला नेमून दिलेली असते, आणि फक्त तीच व्यक्ती ती वापरू शकते) आणि "Users can't route their work through another user's account to use that user's remaining usage" (दुसऱ्या वापरकर्त्याच्या उरलेल्या वापराचा लाभ घेण्यासाठी वापरकर्ते आपले काम त्याच्या खात्यामार्फत चालवू शकत नाहीत) ([Team प्लॅन](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise प्लॅन](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) साठीही असेच).

वैयक्तिक प्लॅनसाठी (Pro, Max), तुमची स्वतःची खाती मर्यादेवर एकमेकांची जागा आपोआप घेऊ शकतात की नाही, हे Anthropic च्या अटी सांगत नाहीत. [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) Anthropic ने स्पष्टपणे परवानगी दिल्याशिवाय Services "through automated or non-human means, whether through a bot, script, or otherwise" (स्वयंचलित किंवा मानवेतर मार्गाने, bot, script किंवा इतर कशानेही) वापरण्यास मनाई करतात, आणि एक बदल वरच्या ठरलेल्या संदेशांपैकी एक तुमच्या वतीने पाठवतो. November 12, 2026 पासून लागू होणारी [Usage Policy](https://www.anthropic.com/legal/aup) ती ज्या गोष्टींना मनाई करते त्यांत "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (आमच्या उत्पादनांमध्ये स्थापित क्षमता, निर्बंध किंवा सुरक्षा-मर्यादा जाणूनबुजून टाळणे किंवा त्यांत ढवळाढवळ करणे) याचा समावेश करते.

`auto_switch` या अटींचे पालन करते असा दावा ccam करत नाही. ते चालू करायचे की नाही हा तुमचा निर्णय आहे, आणि धोकाही तुमचाच: अटी मोडल्या गेल्या असे Anthropic ला वाटले, तर ते एखादे खाते, किंवा Team किंवा Enterprise organization चा प्रवेश, निलंबित करू शकते.

`auto_switch` चालू असताना `ccam doctor` ccam plugin नसलेल्या खात्याबद्दल आणि 2.1.214 पेक्षा जुन्या Claude Code आवृत्त्यांबद्दल इशारा देते, आणि सीट्ससाठी Anthropic चा सीटविषयीचा नियम पुन्हा सांगते ([तपासण्या](/features/doctor)).

## मर्यादेवर पुन्हा चालवणे

`ccam next` हे थांबलेल्या interactive session साठी आहे. non-interactive run, म्हणजे script किंवा scheduled job मधील `claude -p`, तुमची वाट पाहू शकत नाही, त्यामुळे खात्याच्या नावाआधी `--fallback` लावल्यास ccam त्याच्यासाठी खाते बदलते:

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

1. `--fallback` फक्त `-p` किंवा `--print` सोबतच चालते; नाहीतर ccam थांबते आणि `ccam next --auto --continue` कडे निर्देश करते.
2. pipe ने आलेला standard input एकदाच, 10 MiB पर्यंत, वाचला जातो आणि प्रत्येक प्रयत्नाला दिला जातो. Claude Code प्रमाणेच, ccam तो सुरू होण्याची जास्तीत जास्त 3 सेकंद वाट पाहते (`--input-format stream-json` सोडून); pipe उघडा राहिला पण त्याने काहीच पाठवले नाही, तर run ला कोणताही input मिळत नाही आणि ccam याची नोंद standard error वर करते.
3. ccam `exec` ऐवजी `claude` ला child process म्हणून सुरू करते, त्याचा standard output राखून ठेवते, आणि standard error येईल तसा पुढे देते.
4. `claude` 0 ने बाहेर पडल्यास, किंवा दुसऱ्या कारणाने अयशस्वी झाल्यास, राखून ठेवलेला output आणि exit code न बदलता पुढे दिले जातात. signal ने संपलेल्या run चेही (128 पेक्षा मोठा exit code) तसेच होते: त्याला थांबायला सांगितले गेले होते, आणि तो पुन्हा चालवला जात नाही.
5. Claude Code चा संदेश तसे सांगत असेल, म्हणजे [hook](#inside-claude-code) ला माहीत असलेल्या वापर-मर्यादेच्या संदेशांपैकी कोणताही, `--output-format json` किंवा `stream-json` output मध्येही, किंवा खात्याच्या वापरात त्याची 5 तासांची, साप्ताहिक किंवा `--model` window 100% वर दिसत असेल (मर्यादेच्या नोंदीसह), तर ते अपयश वापर मर्यादा म्हणून मोजले जाते.
6. मग त्या प्रयत्नाचा output टाकून दिला जातो आणि तीच कमांड दुसऱ्या लॉग इन असलेल्या खात्याने चालते: सर्वात कमी 5 तासांचा वापर असलेले खाते, मग ज्यांचा वापर शोधता येत नाही ती खाती, क्रमांकाच्या क्रमाने. 100% वरची खाती वगळली जातात, API key लॉगिन असलेली खाती कधीच वापरून पाहिली जात नाहीत, आणि प्रत्येक खाते जास्तीत जास्त एकदाच वापरून पाहिले जाते.
7. आर्ग्युमेंट्समध्ये `-c`, `--continue`, `-r` किंवा `--resume` असल्यास, फक्त त्याच इतिहास गटातील खाती वापरून पाहिली जातात, जेणेकरून संभाषण सापडेल.
8. `--auto` प्रमाणेच फक्त त्याच organization ची खाती वापरून पाहिली जातात (`cross_org = true` हे उठवते).

पुढचे खाते `--auto` च्याच नियमांनी निवडले जाते (`auto_max_usage`, `auto_strategy`, `--model` ची window). एकही खाते उरले नाही तर ccam शेवटच्या प्रयत्नाचा output छापते, कोणत्याही लॉग इन असलेल्या खात्यात जागा उरलेली नाही असे सांगते, आणि code 75 ने बाहेर पडते. खात्याच्या नावाआधी `--wait` दिल्यास (`ccam --fallback --wait=2h work -p …`), ते त्याऐवजी पहिले खाते reset होण्याची वाट पाहते, मग प्रत्येक खाते पुन्हा वापरून पाहू शकते; ते नेहमी आधी किमान एकदा थांबते, जेणेकरून नुकतीच मर्यादा गाठलेले खाते लगेच पुन्हा चालवले जात नाही. हा कालावधी run मधील एकूण वाट पाहण्याच्या वेळेला मर्यादा घालतो, प्रत्येक वाट पाहण्याला वेगळी नाही. space ने लिहिल्यास (`--wait 30m`), ते `ccam next` सारख्याच सूचनेसह नाकारले जाते. Claude Code च्या आर्ग्युमेंट्सआधीचा `--` कोणत्याही launch प्रमाणेच चालतो (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` पहिल्या प्रयत्नाची कमांड छापते.

तुम्ही [`auto_switch`](#auto-switch) चालू केल्याशिवाय interactive sessions आपोआप बदलले जात नाहीत. `ccam statusline install` सह, 5 तासांचा किंवा साप्ताहिक वापर 90% ला पोहोचला की स्टेटस लाइन `→ ccam next --auto --continue` जोडते ([वापर मर्यादा](/features/usage)).

## सुरक्षितता आणि अपवादात्मक स्थिती

- लॉगिन नसलेली खाती कधीच निवडली जात नाहीत. दुसरे लॉग इन असलेले खाते नसल्यास, साधे `ccam next` `no logged-in account other than <name>` ने थांबते.
- `--continue` सह आणि फोल्डरमध्ये एकही संभाषण नसल्यास: `no account has a conversation in <folder> to continue`.
- organization च्या नियमामुळे एकही खाते उरले नाही तर: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- `--auto` सह आणि एकही लॉग इन असलेले खाते नसल्यास: `no logged-in account to pick`. प्रत्येक खाते भरलेले असेल किंवा त्याचा वापर अज्ञात असेल तेव्हा: `` no account with known usage and headroom (see `ccam usage`) ``. दोन्ही exit code 75 ने संपतात.
- `--max-usage`, `--strategy`, `--wait` आणि `--cross-org` फक्त `--auto` सह काम करतात, आणि 1-100 बाहेरचे मूल्य किंवा अज्ञात strategy नाकारली जाते.
- fork संभाषण कॉपी करतो, खात्याची इतर स्थिती नाही: इतिहास शेअर केलेला नसेल तर मूळ खात्याच्या `file-history/` मधील checkpoints सोबत येत नाहीत.
- `--dry-run` सह सुद्धा `--auto` नेटवर्क requests करते. ते फक्त वाचते, आणि token फाइल्समध्ये कधीच लिहीत नाही.
- निवडलेले खाते `ccam <name>` च्याच मार्गाने सुरू होते, त्यामुळे त्याच तपासण्या लागू होतात, यात तोच इतिहास शेअर करणारे दुसरे Claude Code session या फोल्डरमध्ये आधीच चालू असल्याचा इशाराही येतो.
- `last-used` फक्त Claude Code प्रत्यक्ष सुरू झाल्यावर लिहिले जाते. `--dry-run` त्याला हात लावत नाही, त्यामुळे dry run rotation पुढे सरकवत नाही.
- `--fallback` प्रयत्न संपेपर्यंत standard output राखून ठेवते, त्यामुळे `--output-format stream-json` सारखा streaming output शेवटी एकाच वेळी येतो. `stream-json` input stream करणारा आणि आणखी पाठवण्यापूर्वी उत्तरांची वाट पाहणारा client `--fallback` वापरू शकत नाही. पहिल्या प्रयत्नाने मर्यादा गाठण्यापूर्वी बदललेल्या फाइल्स बदललेल्याच राहतात; पुढचा प्रयत्न prompt पुन्हा सुरुवातीपासून सुरू करतो.

## संबंधित

- [वापर मर्यादा](/features/usage)
- [वापर मर्यादा आणि खाते बदलणे](/docs/usage-limits)
- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [Claude Code वापर मर्यादेला पोहोचल्यावर काय करायचे](/guides/usage-limit-workflow)
