---
title: वापर मर्यादा
description: ccam usage प्रत्येक खात्याच्या 5 तासांच्या आणि साप्ताहिक मर्यादा एकाच वेळी दाखवते, आणि ccam statusline खात्याचे नाव आणि वापर Claude Code मध्ये दाखवते.
updated: 2026-11-01
order: 7
translation: machine
sourceHash: "54b296f7da35"
---

`ccam usage` प्रत्येक खात्यासाठी प्रत्येक मर्यादेपैकी किती वापरले ते दाखवते: 5 तासांची मर्यादा, साप्ताहिक मर्यादा आणि प्रत्येक मॉडेलची स्वतःची साप्ताहिक मर्यादा (उदा. Fable). प्रत्येक मर्यादेला एक bar आणि ती रीसेट होईपर्यंतचा वेळ मिळतो. कोणत्या खात्यात जागा उरली आहे ते पाहण्यासाठी, एखादे मोठे काम सुरू करण्यापूर्वी किंवा Claude Code मर्यादा जवळ आल्याचे सांगेल तेव्हा ते वापरा. `ccam statusline` ही Claude Code साठी स्टेटस लाइन कमांड आहे: ती Claude Code मध्ये खात्याचे नाव आणि वापर दाखवते आणि `ccam usage` साठी आकडे साठवते. वापरामागचे token आणि अंदाजे खर्च पाहण्यासाठी, `ccam ccusage` प्रत्येक खात्याच्या संभाषणांवर [ccusage](https://ccusage.com) चालवते ([खाली](#tokens-and-costs-ccam-ccusage)).

## वापर

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts that are not running
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

प्रत्येक खात्याला एक ओळ मिळते, ज्यात त्याचा क्रमांक, त्याचे नाव आणि त्याचे आकडे किती जुने आहेत ते असते; `(cached)` म्हणजे ते endpoint कडून नाही, तर साठवलेल्या फाइलमधून आले आहेत. त्याखाली, प्रत्येक मर्यादेसाठी एक ओळ 20 cells चा bar, टक्केवारी आणि window रीसेट होईपर्यंतचा वेळ दाखवते: `5-hour`, `Weekly` (सर्व मॉडेल्स), मग प्रत्येक मॉडेलची स्वतःची साप्ताहिक window, उदा. `Fable`. खात्याकडे नसलेली window `-` दाखवते. टक्केवारी 80% वर पिवळी आणि 100% वर लाल होते.

window भरल्यामुळे Claude Code ने request नाकारली असेल, तर ती window, साठवलेले आकडे काहीही सांगोत, रीसेट होईपर्यंत 100% दाखवते ([मर्यादेची नोंद](#the-limit-record) पाहा). टर्मिनलवर प्रत्येक खाते लगेच `fetching…` म्हणून दिसते, आणि ते उत्तर देईल तशा त्याच्या ओळी काढल्या जातात.

तुमच्याकडे एकापेक्षा जास्त खाती असतील तेव्हाच दिसणारी शेवटची ओळ [`ccam next --auto`](/features/next) जे खाते निवडेल त्याचे नाव देते. ती त्याच नियमांचे पालन करते:

- `auto_max_usage` आणि `auto_strategy`;
- चालू सत्रे;
- `cross_org` सेट नसेल तर, फक्त शेवटी वापरलेल्या खात्याच्या organization ची खाती;
- 100% वरची प्रत्येक window वगळली जाते, त्यामुळे ज्या खात्याचा आठवडा भरला आहे त्याचे नाव कधीच येत नाही.

डीफॉल्ट `headroom` strategy सह ती `→ Most headroom: <name>` अशी असते; `expiring` किंवा `balance` सह ती कारण देते, जसे `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Team आणि Enterprise सीट्स.** आपोआप होणारी खाते निवड (`ccam next`, `--fallback`, `ccam pick`, डॅशबोर्डमधील `n` आणि, तुम्ही चालू केल्यावर, [`auto_switch`](/features/next#auto-switch)) एकाच Claude Team किंवा Enterprise organization मधील सीट्सदरम्यानही खाते बदलते, आणि यामुळे Anthropic च्या अटींचे उल्लंघन होऊ शकते. त्यावर अवलंबून राहण्यापूर्वी [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) वाचा आणि तुमच्या संस्थेकडे आणि Anthropic कडे चौकशी करा.

लॉग इन नसलेले खाते `not logged in → ccam login <name>` दाखवते, आणि API key खाते `no limits (API)` दाखवते. कालबाह्य access token ``access token expired; `ccam <name>` renews it when it starts`` दाखवते, साठवलेले आकडे असतील तर त्यांच्या वर: पुढच्या सुरुवातीला Claude Code त्याचे नूतनीकरण करते. refresh token कालबाह्य झाल्यामुळे, किंवा सर्व्हरने नकार दिल्यानंतर Claude Code ने काढून टाकल्यामुळे, संपूर्णपणे कालबाह्य झालेले लॉगिन `` login expired → `ccam login <name>` `` दाखवते.

`ccam statusline install` नंतर Claude Code `[work] · Opus · 5h 24% · 7d 41%` सारखी ओळ दाखवते.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--json` | reports JSON म्हणून छापा: `name`, `slot`, `status`, `limits` (खात्यावर usage credits चालू असल्यास `credits` सह), `source`, `at`, `error`, आणि मर्यादेची नोंद windows 100% करत असेल तेव्हा `limited` आणि `limited_windows`. एकही खाते नसल्यास ते `[]` छापते | बंद |
| `--watch[=N]` | तुम्ही Ctrl+C दाबेपर्यंत दर `N` सेकंदांनी पुन्हा काढा. फक्त `--watch`, किंवा 0 किंवा त्यापेक्षा कमी मूल्य, म्हणजे 60. मूल्य `=` नंतर येते: `--watch 10` एका सूचनेसह नाकारले जाते | बंद |
| `-v`, `--verbose` | प्रत्येक खात्याखाली `source:` ओळ जोडा: त्याचे आकडे कुठून आले (usage endpoint, Claude Code, स्टेटस लाइन, साठवलेली फाइल, किंवा Claude Code ने नोंदवलेली वापर मर्यादा). मर्यादेची नोंद आकड्यांची जागा घेत असेल, तर त्यानंतर `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` सारखी ओळ येते | बंद |
| `--no-cache` | 60 सेकंदांच्या cache कडे दुर्लक्ष करून endpoint ला पुन्हा विचारा | बंद |
| `--source live\|claude\|cache` | `live` खात्याच्या token ने endpoint ला विचारते; `claude` Claude Code कडून वापर पाहून घेते (खाली पाहा); `cache` फक्त साठवलेले आकडे दाखवते | `usage_source`, अन्यथा `live` |
| `--refresh` | कालबाह्य access tokens नूतनीकरण करा, फक्त चालू नसलेल्या खात्यांसाठी. एका नूतनीकरणाला 30 सेकंदांपर्यंत लागू शकतात | बंद |
| `--dry-run` | स्वीकारला जातो: `ccam usage` स्वतःच्या cache शिवाय काहीही बदलत नाही. `--refresh` सोबत नाकारला जातो: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | बंद |

`--source` चे इतर कोणतेही मूल्य `--source must be live, claude or cache` देऊन अयशस्वी होते. `ccam usage` कोणतेही आर्ग्युमेंट घेत नाही.

### Claude Code कडूनच वापर

`config.toml` च्या सुरुवातीला (पहिल्या `[[group]]`, `[[profile]]` किंवा `[[folder]]` च्या वर) `usage_source = "claude"` असेल, किंवा `--source claude` दिले असेल, तर ccam लॉगिन token अजिबात वाचत नाही. प्रत्येक खात्यासाठी ते खात्याच्या फोल्डरमध्ये `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` सुरू करते आणि standard input वर Claude Code ची `get_usage` विनंती पाठवते (Claude Code 2.1.270 किंवा नंतरचे). Claude Code स्वतःच्या लॉगिनने वापर पाहते आणि 5 तासांच्या, साप्ताहिक आणि model windows सह उत्तर देते.

`--safe-mode` खात्याचे hooks, MCP सर्व्हर आणि plugins सुरू होऊ देत नाही, आणि कोणतेही सत्र साठवले जात नाही. प्रत्येक खात्यासाठी Claude Code एकदा सुरू करावे लागते (सुमारे एक सेकंद), त्यामुळे हा source `live` पेक्षा हळू आहे. खात्यांना एकाच वेळी विचारले जाते, प्रत्येकासाठी 20 सेकंदांचा timeout ठेवून, आणि 60 सेकंदांचा cache नेहमीप्रमाणे लागू होतो.

ज्या लॉगिनचे access token कालबाह्य झाले आहे, किंवा 5 मिनिटांत कालबाह्य होणार आहे, त्याला 50 सेकंद मिळतात, कारण Claude Code सुरू होताना ते नूतनीकरण करते आणि अर्ध्यावर तोडलेल्या नूतनीकरणामुळे लॉगिन गमावले जाऊ शकते. ccam ने एखादा शोध सोडून दिला, तर ते Claude Code ला SIGTERM ने बाहेर पडायला सांगते, त्यामुळे Claude Code सुरू केलेले नूतनीकरण पूर्ण करते (Windows वर SIGTERM नसतो, त्यामुळे तिथे process संपवली जाते).

संपूर्णपणे कालबाह्य झालेले लॉगिन Claude Code सुरू न करता `` login expired → `ccam login <name>` `` दाखवते. Claude Code ही विनंती experimental म्हणून चिन्हांकित करते, त्यामुळे तिने उत्तर देणे थांबवले तर ccam साठवलेले आकडे दाखवते आणि `ccam doctor` ते नोंदवते. Claude Code ने आपला स्वतःचा शोध अयशस्वी झाल्याचे उत्तर दिले, तर ccam मर्यादा नसलेला plan न दाखवता, `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)` सह साठवलेले आकडे दाखवते.

`--dry-run` सह ccam Claude Code सुरू करत नाही, कारण Claude Code खात्याच्या फोल्डरमध्ये लिहिते; त्याऐवजी ते साठवलेले आकडे दाखवते. `usage_source` `live` किंवा `claude` घेतो; `cache` फक्त `--source` चा पर्याय आहे.

`ccam statusline` standard input वरून Claude Code चा स्टेटस लाइन JSON वाचते. त्याचा एकच पर्याय `--wrap <command>` आहे: आधीपासून असलेली स्टेटस लाइन कमांड, जिच्या आउटपुटच्या पहिल्या ओळीसमोर खात्याचे नाव लावले जाते. `ccam statusline install` कोणतेही पर्याय घेत नाही. दोन्ही `--dry-run` `` --dry-run is not supported by `ccam statusline`; nothing was changed `` देऊन नाकारतात.

## Tokens आणि खर्च: `ccam ccusage`

`ccam usage` एखाद्या खात्याने प्रत्येक मर्यादेपैकी किती वापरले ते दाखवते. [ccusage](https://ccusage.com) त्याऐवजी Claude Code च्या संभाषण logs मधील token आणि अंदाजे खर्च मोजते. स्वतःहून ते Claude Code चे logs फक्त `~/.claude` आणि `~/.config/claude` मध्ये (किंवा `CLAUDE_CONFIG_DIR` मध्ये दिलेल्या फोल्डरमध्ये) शोधते, त्यामुळे ccam च्या खात्यांपैकी त्याला दिसू शकणारे एकमेव खाते म्हणजे `~/.claude` येथे जागेवरच इम्पोर्ट केलेले खाते. `ccam ccusage` संभाषणे असलेल्या प्रत्येक खात्याच्या फोल्डरवर `CLAUDE_CONFIG_DIR` सेट करून ते चालवते, त्यामुळे एकाच अहवालात सर्व खाती येतात:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam आधी स्वतःचे पर्याय वाचते: `-a` (`--account`) खात्यांची नावे किंवा क्रमांक घेतो, स्वल्पविरामाने वेगळे करून किंवा पर्याय अनेकदा देऊन, आणि `--dry-run` कमांड चालवण्याऐवजी ती छापतो. ccam चा पर्याय नसलेल्या पहिल्या आर्ग्युमेंटपासून पुढे सर्व काही तुम्ही टाइप केले तसेच ccusage ला जाते, फक्त `bun x` आणि `npx` साठी `--config` सोडून (खाली पहा); `ccam ccusage -- --help` ccusage ची स्वतःची मदत दाखवते. `ccusage` इन्स्टॉल असेल तर ccam ते चालवते, नाहीतर `bun x ccusage`, आणि तेही नसेल तर `npx ccusage@latest`; हे शेवटचे दोन ccusage npm वरून आणतात. ते ccam च्या स्वतःच्या फोल्डरमध्ये (`~/.local/state/ccam/ccusage`) सुरू होतात, तुम्ही असलेल्या फोल्डरमध्ये नाही, त्यामुळे तिथला प्रकल्प त्याच्या `node_modules/.bin`, त्याच्या `package.json` मधील `bin` किंवा त्याच्या `.npmrc` मधील registry ने ते काय चालवतात हे ठरवू शकत नाही; `ccam --dry-run ccusage` त्या फोल्डरमध्ये जाणारा `cd` दाखवते. यांपैकी काहीच उपलब्ध नसेल तेव्हा ते तुम्हाला Node.js, Bun किंवा ccusage इन्स्टॉल करायला सांगते. Linux आणि macOS वर ccusage, ccam ची जागा घेते; Windows वर ccam त्याची वाट पाहते आणि त्याच्याच exit code ने संपते. WSL मध्ये, Windows drive वरील `ccusage`, `bun` किंवा `npx` ccam वगळते, कारण ते खात्यांच्या फोल्डरशिवाय Windows प्रोग्राम म्हणून चालेल.

- ज्या खात्यांच्या फोल्डरमध्ये संभाषणे आहेत (एक `projects` फोल्डर) ती खाती अहवालात येतात. खाती गटामार्फत शेअर करतात ते संभाषण एकदाच मोजले जाते.
- `-a` सह, गटामार्फत आपला इतिहास शेअर करणारे खाते संपूर्ण गटाची संभाषणे घेऊन येते, आणि ती कोणाची आहेत ते ccam सांगते: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- `blocks` 5 तासांच्या windows दाखवते, आणि प्रत्येक खात्याच्या windows वेगळ्या असतात, त्यामुळे वरीलप्रमाणे त्यासाठी एकच खाते द्या.
- तुम्ही एखादे tool नमूद केले नाही तर, ccusage त्याला सापडणाऱ्या इतर coding tools चाही (उदा. Codex) अहवाल देते: `ccam ccusage claude daily` म्हणजे फक्त Claude Code.
- ccam ते चालवत असताना, ccusage आपल्या सेटिंग्ज (`ccusage.json`) फक्त `./.ccusage` आणि खात्यांच्या फोल्डरमधून वाचते. `--config` दुसरी फाइल देतो: `ccam ccusage daily --config ~/.config/claude/ccusage.json`. `bun x` आणि `npx` ccam च्या स्वतःच्या फोल्डरमध्ये सुरू होत असल्याने, `./.ccusage/ccusage.json` असेल आणि तुम्ही कोणतीही फाइल दिली नसेल, तर ccam त्यांना ती `--config` सह देते, आणि relative `--config` पूर्ण पाथ म्हणून देते. मग ccusage इतर relative पाथ, जसे `--pi-path`, ccam च्या फोल्डरमधून वाचते, त्यामुळे ते पूर्ण द्या.
- ज्या फोल्डरच्या पाथमध्ये स्वल्पविराम आहे तो ccusage दोन फोल्डर म्हणून वाचेल, त्यामुळे ccam तो नाकारते आणि त्या खात्याचे नाव सांगते; ते खाते `-a` ने वगळा.
- `ccusage` ही ccam ची कमांड आहे, त्यामुळे ccusage नावाचे खाते `ccam run ccusage` ने सुरू होते; `ccam doctor` तसे सांगते.

## हे कसे काम करते

`ccam usage` चे आकडे दोन ठिकाणांहून येतात: प्रत्येक खात्याच्या स्वतःच्या लॉगिनने बोलावलेला Anthropic चा दस्तऐवजीकृत नसलेला usage endpoint, `GET https://api.anthropic.com/api/oauth/usage`, आणि Claude Code आपल्या स्टेटस लाइन कमांडला देत असलेली मूल्ये. हे सर्व `internal/usage` मध्ये आहे.

1. ccam `config.toml` लोड करते आणि प्रत्येक खात्यासाठी एक `usage.Account` तयार करते (`cli.usageCmd`): त्याचा फोल्डर `~/.local/share/ccam/profiles/<name>/`, त्याचे session markers (खात्याच्या फोल्डरमधील Claude Code च्या स्वतःच्या नोंदींसह, ते खाते चालू आहे का ते सांगतात) आणि ते जागेवरच इम्पोर्ट केले आहे का.
2. ccam प्रत्येक खात्याचा वापर एकाच वेळी शोधते, प्रत्येक खात्यासाठी एक goroutine (`Fetcher.FetchAll`). ते सर्व एकच `http.Client` वापरतात, त्यामुळे endpoint ची connections पुन्हा वापरली जातात. connect होण्यास 2 s लागू शकतात, आणि प्रत्येक खात्याला एकूण 3 s मिळतात. client redirects पाळत नाही, कारण redirect token दुसरीकडे पुन्हा पाठवेल, आणि उत्तराचा जास्तीत जास्त 1 MB वाचतो.
3. प्रत्येक खाते पूर्ण होईल तशी त्याची ओळ पुन्हा काढली जाते. अयशस्वी होणारे किंवा वेळेत उत्तर न देणारे खाते साठवलेल्या आकड्यांवर येते, त्यामुळे एक हळू खाते इतरांना कधीच अडवत नाही.

आकृतीमध्ये प्रत्येक खात्याचा शोध त्याच्या स्वतःच्या goroutine मध्ये होतो, आणि वेळ संपलेले खाते त्याऐवजी आपले साठवलेले आकडे दाखवते:

<!-- motion:usage-fetch -->

`personal` आणि `work` वेळेत उत्तर देतात आणि live आकडे दाखवतात. `side` ला 3 s पेक्षा जास्त वेळ लागतो, म्हणून ccam त्याऐवजी त्याचे साठवलेले आकडे दाखवते (`fallback()`). मग ccam प्रत्येक खात्यासाठी एक block काढते (`usage.Table()`) आणि शेवटच्या ओळीत नाव दिलेले खाते निवडते (`usage.Best()`).

एका खात्यासाठी ccam या तपासण्या क्रमाने करते (`Fetcher.Fetch`):

<!-- motion:usage-cache -->

1. `--source cache` सह ते थेट `fallback()` कडे जाते.
2. `--no-cache` सेट नसेल तर, `~/.local/state/ccam/usage/<name>.json` मधील 60 सेकंदांपेक्षा (`CacheTTL`) कमी वेळापूर्वी मोजलेला live परिणाम वापरला जातो, आत्ताच्या वेळेनुसार तपासून: त्यानंतर रीसेट झालेली window 0% दाखवते, आणि घड्याळापेक्षा पुढची तारीख असलेली फाइल ताजी मानली जात नाही. `ccam usage` लागोपाठ दोनदा चालवल्यास प्रत्येक खात्यासाठी एकच request लागते.
3. ccam `.credentials.json` (किंवा macOS Keychain) न बदलता वाचते आणि फक्त `claudeAiOauth` key decode करते (`readCredentials()`). subscription लॉगिन नसल्यास `not logged in`, किंवा `.claude.json` मध्ये API key असल्यास `no limits (API)`. refresh token कालबाह्य झालेले, किंवा Claude Code ने काढून टाकलेले लॉगिन `` login expired → `ccam login <name>` `` असते आणि त्याचा शोध घेतला जात नाही.
4. कालबाह्य access token वापरले जात नाही. `--refresh` शिवाय, खाते ``access token expired; `ccam <name>` renews it when it starts`` सह आपले साठवलेले आकडे दाखवते.
5. नाहीतर ccam endpoint बोलावते. 401 किंवा 403 कालबाह्य लॉगिन मानले जाते (`login expired: authentication failed (401)`), इतर कोणताही status किंवा वाचता न येणारा body त्रुटी मानली जाते, आणि दोन्ही वेळा साठवलेले आकडे वापरले जातात. Ctrl+C ने संपलेला शोध `cancelled` सांगतो. 429 म्हणजे endpoint ला या खात्यासाठी खूप जास्त requests मिळाल्या आहेत: ccam पुन्हा विचारण्याच्या वेळेसह `<name>.pause` लिहिते (`Retry-After`, किंवा 5 मिनिटे, जास्तीत जास्त एक तास), `rate limited by the usage endpoint; asking again after 15:10` सह साठवलेले आकडे दाखवते, आणि तोपर्यंत त्या खात्यासाठी endpoint बोलावत नाही, `--no-cache` असले तरी. लगेच पुन्हा प्रयत्न केल्याने rate limit फक्त लांबले असते.
6. चांगले उत्तर मोजल्याच्या वेळेसह `<name>.json` मध्ये लिहिले जाते.

साठवलेल्या आकड्यांसाठी ccam सर्वात नवीन साठवलेला शोध घेते (`<name>.json`, किंवा Claude Code च्या शोधांमधील `<name>.claude.json`) आणि, `<name>.statusline.json` त्यानंतर लिहिली गेली असेल तर, त्यावर स्टेटस लाइनच्या 5 तासांच्या आणि साप्ताहिक windows ठेवते (`fallback()`); स्टेटस लाइनकडे नसलेल्या model windows शोधाने साठवल्या तशाच राहतात. साठवलेले आकडे दाखवण्यापूर्वी ते प्रत्येक window ची reset time तपासते: आकडे साठवल्यानंतर रीसेट झालेली window 0% दाखवते, कारण जुनी टक्केवारी आता लागू होत नाही. म्हणूनच वरती `personal` reset time शिवाय `0%` दाखवते.

endpoint चे दोन response formats आहेत. नवीन format मध्ये `session` (5 तास), `weekly_all` आणि `weekly_scoped` (मॉडेलनुसार) नोंदी असलेला `limits` array असतो. जुन्यामध्ये `five_hour`, `seven_day`, `seven_day_opus` आणि `seven_day_sonnet` objects असतात. ccam दोन्ही वाचते (`ParseLive`): `limits` array जिंकतो, आणि तो जे व्यापत नाही ते जुने objects भरतात. कोणत्याही मर्यादा नसलेला body बदललेला format म्हणून नोंदवला जातो, ज्याबद्दल `ccam doctor` तुम्हाला सांगते.

### मर्यादेची नोंद

Claude Code एखादा turn वापर मर्यादेवर थांबवते, तेव्हा [ccam plugin](/docs/install#claude-code-plugin) चा hook संभाषणातून नाकारलेली request वाचतो: कोणती window भरली होती (5 तास, साप्ताहिक, Opus किंवा Sonnet) आणि ती कधी रीसेट होते. hook ही माहिती `~/.local/state/ccam/usage/<name>.limit.json` मध्ये लिहितो. ccam ही नोंद प्रत्येक source मधून त्याला माहीत असलेल्या इतर सर्व गोष्टींवर लागू करते (`Fetch`): तिच्यातील प्रत्येक window तिच्या reset पर्यंत 100% मोजली जाते, साठवलेले आकडे कितीही जुने किंवा कमी असले तरी, त्यामुळे `ccam next --auto`, `ccam pick`, `--fallback` आणि [`auto_switch`](/features/next#auto-switch) तोपर्यंत ते खाते पुन्हा निवडत नाहीत.

दुसरे काहीच माहीत नसेल, तर खाते फक्त हीच नोंद दाखवते (`-v` सह `source: the usage limit Claude Code reported`). कोणीही समोर नसताना चालणारी सत्रे (`-p` runs आणि background सत्रे) कोणतीही नोंद लिहीत नाहीत, तसेच organization ने लावलेल्या थांब्यासाठीही (उदाहरणार्थ, usage संपल्यावर किंवा group ची मर्यादा $0 असताना) नोंद लिहिली जात नाही, कारण तिचे admins ते कधीही उठवू शकतात. दुसऱ्या Claude खात्यात `ccam login` केल्याने नोंद काढून टाकली जाते; फाइल हटवल्यानेही ती पुसली जाते.

### Tokens नूतनीकरण करणे

Refresh tokens बदलत राहतात: एक नूतनीकरण केल्यावर तुम्हाला नवीन refresh token मिळते आणि जुने अवैध होते. Claude Code त्याच खात्याने चालू असताना ccam ने लॉगिन नूतनीकरण केले, तर ते session साइन आउट होईल. म्हणून `ccam usage` tokens वाचते आणि कधीच लिहीत नाही, तुम्ही `--refresh` दिल्याशिवाय. तेव्हा ccam (`Fetcher.refresh`):

1. `ccam import` ने जागेवरच नोंदवलेली खाती नाकारते, कारण साधे `claude` ccam ला न कळता ते लॉगिन वापरत असू शकते;
2. macOS Keychain मध्ये साठवलेली लॉगिन्स नाकारते;
3. एक lock फाइल घेते, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, त्यामुळे दोन ccam processes कधीच दोघेही तोच token नूतनीकरण करत नाहीत; तो धरलेला असताना खाते `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)` दाखवते. `ccam rename`, `remove` आणि `logout` सुद्धा तो घेतात, आणि नूतनीकरण चालू असताना नकार देतात;
4. lock खाली, खात्याचे session चालू असेल तर नकार देते, ते कसेही सुरू केलेले असो (ccam मार्फत, साधे `claude`, IDE, किंवा background सत्र आणि त्यामागचा daemon, जो स्वतः लॉगिन नूतनीकरण करतो), आणि दुसऱ्या process ने token आधीच नूतनीकरण केले असेल तर request वगळते;
5. Claude Code खात्याच्या फोल्डरमध्ये स्वतःचे नूतनीकरणाचे locks (`.oauth_refresh.lock`, `<folder>.lock`) धरून असताना काहीही पाठवत नाही: तेव्हा Claude Code नेमके तेच लॉगिन नूतनीकरण करत असते;
6. refresh token token endpoint ला post करते आणि उत्तराची 30 सेकंदांपर्यंत वाट पाहते, जितका वेळ Claude Code स्वतःच्या नूतनीकरणांना देते. एकदा पाठवलेली request कधीही तोडली जात नाही, Ctrl+C ने सुद्धा नाही: सर्व्हरने refresh token आधीच बदललेला असू शकतो, आणि फक्त त्याच्या उत्तरातच नवीन लॉगिन असते. Ctrl+C किंवा hangup इतर शोध संपवतात, आणि उत्तर सेव्ह झाल्यावर ccam बाहेर पडते (Ctrl+C नंतर 130);
7. Claude Code च्या write lock (`.storage-write.lock`) खाली सेव्ह करते आणि त्या वेळी असलेल्या `.credentials.json` मधील फक्त token fields बदलते: दरम्यान Claude Code ने लिहिलेले (MCP साइन-इन) तसेच राहते, दरम्यान Claude Code ने नूतनीकरण केलेले लॉगिन ठेवले जाते, आणि लॉगआउट कधीही उलटवले जात नाही. उत्तरात नवीन refresh token चे आयुष्य दिलेले असेल, तर ते Claude Code प्रमाणेच `refreshTokenExpiresAt` म्हणून साठवले जाते.

पूर्ण न झालेले refresh, खात्याच्या ओळीवर किंवा त्याच्या साठवलेल्या आकड्यांच्या वर, त्याचे कारण सांगते: उदाहरणार्थ `token not renewed: a session is running`, किंवा request पाठवली गेली पण तिचे उत्तर वापरता किंवा सेव्ह करता आले नाही तेव्हा `login may be lost: …`, आणि त्यानंतर Claude Code ने लॉग इन करायला सांगितले तर `ccam login <name>` चालवण्याच्या सल्ल्यासह.

### स्टेटस लाइन collector

Claude Code स्टेटस लाइन कमांड चालवते आणि तिला pipe ने एक JSON document देते. ccam त्यातून मॉडेलचे नाव आणि `rate_limits` वाचते:

<!-- motion:usage-statusline -->

1. ccam `CCAM_PROFILE` वरून खाते शोधते, जो ccam Claude Code सुरू करताना सेट करते, नाहीतर `CLAUDE_CONFIG_DIR` खात्यांच्या फोल्डर्सशी जुळवून (`runStatusline`).
2. ccam `rate_limits.five_hour` आणि `rate_limits.seven_day` वाचते (`ParseStatusLine`). API key खात्यांसाठी, किंवा session च्या पहिल्या response आधी, JSON मध्ये rate limits नसतात; तेव्हा काहीही साठवले जात नाही.
3. ccam `~/.local/state/ccam/usage/<name>.statusline.json` लिहिते (`SaveStatusLine`). endpoint उत्तर देत नाही तेव्हा `ccam usage` ती वापरते.
4. ccam `[name] · model · 5h N% · 7d N%` छापते. `--wrap` सह, ते तुमची कमांड त्याच input सह `sh -c` मार्फत चालवते आणि Claude Code स्टेटस लाइन दाखवते तसा तिचा output छापते: रिकामी नसलेली प्रत्येक ओळ, पहिल्या ओळीसमोर `[name]` लावून. अयशस्वी होणारी कमांड फक्त खात्याचे नाव दाखवते. Windows वर ते कमांड Claude Code प्रमाणेच चालवते: Git Bash इन्स्टॉल असल्यास त्याने, नाहीतर PowerShell ने.
5. 5 तासांचा किंवा साप्ताहिक वापर 90% किंवा त्याहून जास्त असेल आणि दुसरे खाते असेल, तर ओळीच्या शेवटी `→ ccam next --auto --continue` येते, उदाहरणार्थ `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

wrap केलेल्या कमांडला स्वतःची वेळेची मर्यादा नसते: Claude Code स्टेटस लाइन कमांडला 10 मिनिटे देते आणि नवीन update सुरू झाल्यावर चालू run थांबवते, त्यामुळे हळू कमांड ccam शिवाय जशी दिसली असती तशीच दिसते. Claude Code ccam संपवते तेव्हा कमांडही संपते (Linux आणि macOS वर तिने सुरू केलेल्या प्रत्येक process सह).

`ccam statusline install` प्रत्येक खात्याच्या `settings.json` मध्ये `statusLine` लिहिते आणि त्या object च्या इतर keys, जसे `padding`, `refreshInterval` आणि `hideVimModeIndicator`, ठेवते. खात्याकडे आधीच स्टेटस लाइन कमांड असेल, तर ccam ती `--wrap '<old command>'` म्हणून ठेवते; Windows वर `--wrap-base64 <base64>` म्हणून, ज्याला कोणत्याही shell मध्ये quoting लागत नाही. स्टेटस लाइन आधीच ccam चालवत असेल (दुसऱ्या पाथवरून, किंवा जुन्या आवृत्तीने लिहिल्याप्रमाणे quote केलेली), तर ccam ती wrap करण्याऐवजी बदलते, आणि ती जी कमांड wrap करते ती ठेवते.

Windows वर ccam चा पाथ forward slashes आणि home फोल्डरसाठी `~` वापरून लिहिला जातो (`~/.local/bin/ccam.exe statusline`), जो Git Bash आणि PowerShell दोन्ही समजतात. quotes लागणारा पाथ Claude Code वापरत असलेल्या shell साठी पूर्ण लिहिला जातो: Git Bash साठी `'…'`, PowerShell साठी `& '…'`. गटामार्फत शेअर केलेली `settings.json` तिच्या लिंकमार्फत एकदाच लिहिली जाते, त्यामुळे गटातील प्रत्येक खात्याला ती मिळते. ते पुन्हा चालवल्यास आधीच ccam वापरणाऱ्या फाइल्ससाठी `(already set)` छापले जाते.

## सुरक्षितता आणि अपवादात्मक स्थिती

- तुम्ही `--refresh` दिल्याशिवाय `ccam usage` कधीच token फाइल बदलत नाही. access token कालबाह्य झाल्यावर, `ccam <name>` ने खाते एकदा चालवा, म्हणजे Claude Code ते नूतनीकरण करते. संपूर्णपणे कालबाह्य झालेल्या लॉगिनला `ccam login <name>` लागते.
- `--refresh` चालू खाती (Claude Code ccam शिवाय सुरू केले असले तरीही), इम्पोर्ट केलेली खाती किंवा Keychain लॉगिन्सना कधीच हात लावत नाही. नाकारलेले refresh साठवलेले आकडे आणि कारण दाखवते (`token not renewed: …`).
- usage endpoint सार्वजनिकरित्या दस्तऐवजीकृत नाही आणि बदलू शकतो. तो बदलला तरी तुम्हाला साठवलेले आकडे दिसतात, आणि `ccam doctor` format मधील बदल नोंदवते.
- साठवलेले आकडे आणि उत्तर दोन्ही नसल्यास, खात्याची ओळ `lookup failed:` आणि कारण दाखवते, उदाहरणार्थ `--source cache` सह `no saved usage`.
- cache फाइल्समध्ये फक्त टक्केवारी आणि reset times असतात, कधीच tokens नसतात.

## संबंधित

- [वापर मर्यादा आणि खाते बदलणे](/docs/usage-limits)
- [पुढच्या खात्यावर जाणे](/features/next)
- [Claude Code वापर मर्यादेला पोहोचल्यावर काय करायचे](/guides/usage-limit-workflow)
- [कमांड संदर्भ](/docs/commands)
