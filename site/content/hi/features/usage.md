---
title: उपयोग सीमाएँ
description: ccam usage हर अकाउंट की 5 घंटे और साप्ताहिक सीमाएँ एक साथ दिखाता है, और ccam statusline अकाउंट का नाम और उपयोग Claude Code में दिखाता है।
updated: 2026-11-01
order: 7
translation: machine
sourceHash: "54b296f7da35"
---

`ccam usage` हर अकाउंट के लिए दिखाता है कि हर सीमा का कितना हिस्सा इस्तेमाल हुआ है: 5 घंटे की सीमा, साप्ताहिक सीमा और हर मॉडल की अपनी साप्ताहिक सीमा (जैसे Fable)। हर सीमा को एक पट्टी और उसके रीसेट होने तक का समय मिलता है। कोई लंबा काम शुरू करने से पहले, या जब Claude Code बताए कि कोई सीमा पास है, तब इसे इस्तेमाल करें, ताकि पता चले कि किस अकाउंट में गुंजाइश बची है। `ccam statusline`, Claude Code के लिए स्टेटस लाइन कमांड है: यह Claude Code के अंदर अकाउंट का नाम और उपयोग दिखाता है और `ccam usage` के लिए संख्याएँ सेव करता है। उपयोग के पीछे के token और अनुमानित लागत के लिए, `ccam ccusage` हर अकाउंट की बातचीत पर [ccusage](https://ccusage.com) चलाता है ([नीचे](#tokens-and-costs-ccam-ccusage))।

## उपयोग

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

हर अकाउंट को एक पंक्ति मिलती है जिसमें उसका नंबर, नाम और यह होता है कि उसकी संख्याएँ कितनी पुरानी हैं; `(cached)` का मतलब है कि वे endpoint से नहीं, सेव की गई फ़ाइल से आई हैं। उसके नीचे हर सीमा की एक पंक्ति 20 खानों की पट्टी, प्रतिशत और विंडो के रीसेट होने तक का समय दिखाती है: `5-hour`, `Weekly` (सभी मॉडल), फिर हर मॉडल की अपनी साप्ताहिक विंडो, जैसे `Fable`। जो विंडो अकाउंट के पास नहीं है, वह `-` दिखाती है। प्रतिशत 80% पर पीला और 100% पर लाल हो जाता है।

अगर Claude Code ने किसी विंडो के भरे होने के कारण कोई request अस्वीकार की, तो वह विंडो अपने reset तक 100% दिखाती है, सेव की गई संख्याएँ चाहे जो कहें ([limit record](#the-limit-record) देखें)। टर्मिनल पर हर अकाउंट तुरंत `fetching…` के रूप में दिखता है, और उसके जवाब देते ही उसकी पंक्तियाँ बनती हैं।

आख़िरी पंक्ति, जो केवल तभी दिखती है जब आपके पास एक से ज़्यादा अकाउंट हों, उस अकाउंट का नाम बताती है जिसे [`ccam next --auto`](/features/next) चुनता। यह उन्हीं नियमों का पालन करती है:

- `auto_max_usage` और `auto_strategy`;
- चल रहे session;
- `cross_org` सेट न हो तो केवल पिछली बार इस्तेमाल हुए अकाउंट के संगठन के अकाउंट;
- 100% पर पहुँची हर विंडो छोड़ी जाती है, इसलिए जिस अकाउंट का हफ़्ता भर चुका है उसका नाम कभी नहीं आता।

डिफ़ॉल्ट `headroom` strategy के साथ इसमें `→ Most headroom: <name>` लिखा होता है; `expiring` या `balance` के साथ यह कारण बताती है, जैसे `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`।

**Team और Enterprise सीटें।** अपने-आप अकाउंट चुनना (`ccam next`, `--fallback`, `ccam pick`, डैशबोर्ड का `n` और, चालू करने के बाद, [`auto_switch`](/features/next#auto-switch)) एक ही Claude Team या Enterprise संगठन की सीटों के बीच भी स्विच करता है, जो Anthropic की शर्तों का उल्लंघन हो सकता है। इस पर भरोसा करने से पहले [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) पढ़ें और अपने संगठन और Anthropic से पता कर लें।

जिस अकाउंट में लॉगिन नहीं है, वह `not logged in → ccam login <name>` दिखाता है, और API key वाला अकाउंट `no limits (API)` दिखाता है। एक्सपायर हुआ access token ``access token expired; `ccam <name>` renews it when it starts`` दिखाता है, अगर सेव की गई संख्याएँ हों, तो उनके ऊपर: Claude Code अगली शुरुआत पर उसे renew करता है। जो लॉगिन पूरी तरह एक्सपायर हो गया हो, क्योंकि उसका refresh token एक्सपायर हो गया या सर्वर के मना करने के बाद Claude Code ने उसे हटा दिया, वह `` login expired → `ccam login <name>` `` दिखाता है।

`ccam statusline install` के बाद, Claude Code `[work] · Opus · 5h 24% · 7d 41%` जैसी पंक्ति दिखाता है।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--json` | रिपोर्ट को JSON के रूप में प्रिंट करें: `name`, `slot`, `status`, `limits` (जब अकाउंट में usage credits चालू हों तो `credits` के साथ), `source`, `at`, `error`, और `limited` तथा `limited_windows` जब कोई limit record विंडो को 100% पर सेट करता है। कोई अकाउंट न होने पर यह `[]` प्रिंट करता है | बंद |
| `--watch[=N]` | Ctrl+C दबाने तक हर `N` सेकंड में फिर से बनाएँ। अकेले `--watch`, या 0 या उससे कम मान का मतलब 60 है। मान `=` के बाद आता है: `--watch 10` एक संकेत के साथ अस्वीकार होता है | बंद |
| `-v`, `--verbose` | हर अकाउंट के नीचे एक `source:` पंक्ति जोड़ें: उसकी संख्याएँ कहाँ से आईं (usage endpoint, Claude Code, स्टेटस लाइन, सेव की गई फ़ाइल, या Claude Code की बताई उपयोग सीमा)। जब कोई limit record संख्याओं की जगह ले, तो उसके बाद `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` जैसी पंक्ति आती है | बंद |
| `--no-cache` | 60 सेकंड की cache अनदेखी करें और endpoint से फिर पूछें | बंद |
| `--source live\|claude\|cache` | `live` अकाउंट के token से endpoint से पूछता है; `claude` Claude Code से पता करवाता है (नीचे देखें); `cache` केवल सेव की गई संख्याएँ दिखाता है | `usage_source`, वरना `live` |
| `--refresh` | एक्सपायर हुए access tokens renew करें, केवल उन अकाउंट्स के लिए जो चल नहीं रहे। एक renewal में 30 सेकंड तक लग सकते हैं | बंद |
| `--dry-run` | स्वीकार किया जाता है: `ccam usage` अपनी cache के अलावा कुछ नहीं बदलता। `--refresh` के साथ अस्वीकार किया जाता है: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | बंद |

`--source` का कोई और मान `--source must be live, claude or cache` के साथ विफल होता है। `ccam usage` कोई आर्ग्युमेंट नहीं लेता।

### ख़ुद Claude Code से उपयोग

`config.toml` में सबसे ऊपर (पहले `[[group]]`, `[[profile]]` या `[[folder]]` से ऊपर) `usage_source = "claude"` होने पर, या `--source claude` के साथ, ccam लॉगिन token बिल्कुल नहीं पढ़ता। हर अकाउंट के लिए यह अकाउंट के फ़ोल्डर में `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` शुरू करता है और standard input पर Claude Code का `get_usage` request भेजता है (Claude Code 2.1.270 या उसके बाद)। Claude Code अपने ही लॉगिन से उपयोग पता करता है और 5 घंटे, साप्ताहिक और मॉडल विंडो के साथ जवाब देता है।

`--safe-mode` अकाउंट के hooks, MCP सर्वर और plugins को शुरू होने से रोकता है, और कोई session सेव नहीं होता। हर अकाउंट के लिए Claude Code को एक बार शुरू होना पड़ता है (लगभग एक सेकंड), इसलिए यह source `live` से धीमा है। अकाउंट्स से एक साथ पूछा जाता है, हर एक 20 सेकंड के timeout के साथ, और 60 सेकंड की cache हमेशा की तरह लागू होती है।

जिस लॉगिन का access token एक्सपायर हो गया है, या 5 मिनट के भीतर एक्सपायर होने वाला है, उसे 50 सेकंड मिलते हैं, क्योंकि Claude Code शुरू होते समय उसे renew करता है और बीच में कटा renewal लॉगिन खो सकता है। अगर ccam किसी खोज को छोड़ देता है, तो वह Claude Code से SIGTERM के साथ बंद होने को कहता है, ताकि Claude Code शुरू किया हुआ renewal पूरा कर ले (Windows में SIGTERM नहीं है, इसलिए वहाँ process को बंद कर दिया जाता है)।

जो लॉगिन पूरी तरह एक्सपायर हो गया हो, वह Claude Code शुरू किए बिना `` login expired → `ccam login <name>` `` दिखाता है। Claude Code इस request को experimental बताता है, इसलिए अगर यह जवाब देना बंद कर दे, तो ccam सेव की गई संख्याएँ दिखाता है और `ccam doctor` इसकी सूचना देता है। जब Claude Code जवाब देता है कि उसकी अपनी खोज विफल हुई, तो ccam सेव की गई संख्याएँ `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)` के साथ दिखाता है, और उसे बिना सीमाओं वाला प्लान नहीं मानता।

`--dry-run` के साथ ccam Claude Code शुरू नहीं करता, क्योंकि Claude Code अकाउंट के फ़ोल्डर में लिखता है; इसके बजाय यह सेव की गई संख्याएँ दिखाता है। `usage_source` में `live` या `claude` दिया जा सकता है; `cache` केवल `--source` का एक विकल्प है।

`ccam statusline`, standard input से Claude Code का स्टेटस लाइन JSON पढ़ता है। इसका एकमात्र विकल्प `--wrap <command>` है: एक मौजूदा स्टेटस लाइन कमांड, जिसके आउटपुट की पहली पंक्ति के आगे अकाउंट का नाम लगता है। `ccam statusline install` कोई विकल्प नहीं लेता। दोनों `--dry-run` को `` --dry-run is not supported by `ccam statusline`; nothing was changed `` के साथ अस्वीकार करते हैं।

## Token और लागत: `ccam ccusage`

`ccam usage` दिखाता है कि किसी अकाउंट ने हर सीमा का कितना हिस्सा इस्तेमाल किया है। इसके बजाय [ccusage](https://ccusage.com) Claude Code की बातचीत के logs में token और अनुमानित लागत गिनता है। अकेले चलने पर यह Claude Code के logs केवल `~/.claude` और `~/.config/claude` में (या `CLAUDE_CONFIG_DIR` जिन फ़ोल्डरों के नाम देता है, उनमें) ढूँढता है, इसलिए ccam का केवल वही अकाउंट इसे दिख सकता है जो `~/.claude` में अपनी जगह पर import किया गया हो। `ccam ccusage` इसे `CLAUDE_CONFIG_DIR` को बातचीत वाले हर अकाउंट के फ़ोल्डरों पर सेट करके चलाता है, ताकि एक ही रिपोर्ट उन सबको कवर करे:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam पहले अपने विकल्प पढ़ता है: `-a` (`--account`) अकाउंट के नाम या नंबर लेता है, कॉमा से अलग किए हुए या विकल्प दोहराकर, और `--dry-run` कमांड चलाने के बजाय उसे प्रिंट करता है। जो पहला आर्ग्युमेंट ccam का विकल्प नहीं है, वहाँ से आगे सब कुछ वैसे ही ccusage को जाता है जैसे आपने टाइप किया, `bun x` और `npx` के लिए `--config` को छोड़कर (नीचे देखें); `ccam ccusage -- --help` ccusage की अपनी help दिखाता है। ccam इंस्टॉल किया हुआ `ccusage` चलाता है, वह न हो तो `bun x ccusage`, वह भी न हो तो `npx ccusage@latest`; ये आख़िरी दोनों इसे npm से डाउनलोड करते हैं। ये दोनों ccam के अपने फ़ोल्डर (`~/.local/state/ccam/ccusage`) में शुरू होते हैं, उस फ़ोल्डर में नहीं जिसमें आप हैं, इसलिए वहाँ का कोई प्रोजेक्ट अपने `node_modules/.bin`, अपनी `package.json` के किसी `bin` या अपनी `.npmrc` की किसी registry से यह तय नहीं कर सकता कि ये क्या चलाते हैं; `ccam --dry-run ccusage` उस फ़ोल्डर में जाने वाला `cd` दिखाता है। जब इनमें से कोई भी मौजूद न हो, तो ccam Node.js, Bun या ccusage इंस्टॉल करने को कहता है। Linux और macOS पर ccusage, ccam की जगह ले लेता है; Windows पर ccam उसका इंतज़ार करता है और ccusage के exit code के साथ बाहर निकलता है। WSL में, ccam किसी Windows ड्राइव पर मौजूद `ccusage`, `bun` या `npx` को छोड़ देता है, क्योंकि वह अकाउंट फ़ोल्डरों के बिना Windows प्रोग्राम के रूप में चलेगा।

- रिपोर्ट उन अकाउंट्स को कवर करती है जिनके फ़ोल्डर में बातचीत है (एक `projects` फ़ोल्डर)। जो बातचीत अकाउंट्स किसी ग्रुप के ज़रिए साझा करते हैं, वह एक ही बार गिनी जाती है।
- `-a` के साथ, जो अकाउंट किसी ग्रुप के ज़रिए अपना इतिहास साझा करता है, वह पूरे ग्रुप की बातचीत साथ लाता है, और ccam बताता है कि वह किन अकाउंट्स की है: `ccam: work shares its conversation history with personal, so the report counts their conversations too`।
- `blocks` 5 घंटे की विंडो दिखाता है, और हर अकाउंट की अपनी विंडो होती है, इसलिए इसके लिए एक ही अकाउंट का नाम दें, जैसा ऊपर है।
- ccusage उसे मिलने वाले दूसरे coding tools की भी रिपोर्ट देता है, जैसे Codex, जब तक आप किसी एक का नाम न दें: `ccam ccusage claude daily` केवल Claude Code की रिपोर्ट देता है।
- जब ccam इसे चलाता है, तब ccusage अपनी सेटिंग्स (`ccusage.json`) केवल `./.ccusage` और अकाउंट फ़ोल्डरों से पढ़ता है। `--config` किसी दूसरी फ़ाइल का नाम देता है: `ccam ccusage daily --config ~/.config/claude/ccusage.json`। चूँकि `bun x` और `npx` ccam के अपने फ़ोल्डर में शुरू होते हैं, इसलिए जब `./.ccusage/ccusage.json` मौजूद हो और आप कोई फ़ाइल न बताएँ, तो ccam उन्हें वह फ़ाइल `--config` के साथ देता है, और relative `--config` को पूरे पाथ के रूप में देता है। तब ccusage दूसरे relative पाथ, जैसे `--pi-path`, ccam के फ़ोल्डर से पढ़ता है, इसलिए उन्हें पूरा लिखें।
- जिस फ़ोल्डर के path में कॉमा हो, उसे ccusage दो फ़ोल्डर के रूप में पढ़ेगा, इसलिए ccam उसे अस्वीकार करता है और अकाउंट का नाम बताता है; `-a` से उस अकाउंट को बाहर रखें।
- `ccusage` एक ccam कमांड है, इसलिए ccusage नाम का अकाउंट `ccam run ccusage` से शुरू होता है; `ccam doctor` यह बताता है।

## यह कैसे काम करता है

`ccam usage` की संख्याएँ दो जगहों से आती हैं: Anthropic का undocumented usage endpoint, `GET https://api.anthropic.com/api/oauth/usage`, जिसे हर अकाउंट के अपने लॉगिन से बुलाया जाता है, और वे मान जो Claude Code अपनी स्टेटस लाइन कमांड को देता है। यह सब `internal/usage` में है।

1. ccam `config.toml` लोड करता है और हर अकाउंट के लिए एक `usage.Account` बनाता है (`cli.usageCmd`): उसका फ़ोल्डर `~/.local/share/ccam/profiles/<name>/`, उसके session markers (अकाउंट फ़ोल्डर में Claude Code के अपने records के साथ, ये बताते हैं कि वह चल रहा है या नहीं) और यह कि वह अपनी जगह पर import किया गया था या नहीं।
2. ccam हर अकाउंट को एक साथ खोजता है, हर अकाउंट के लिए एक goroutine (`Fetcher.FetchAll`)। ये सभी एक ही `http.Client` साझा करते हैं, इसलिए endpoint से कनेक्शन दोबारा इस्तेमाल होते हैं। कनेक्ट होने में 2 s लग सकते हैं, और हर अकाउंट को कुल 3 s मिलते हैं। Client redirects का पालन नहीं करता, क्योंकि redirect token को कहीं और भेज देता, और जवाब का अधिकतम 1 MB पढ़ता है।
3. हर अकाउंट के पूरा होते ही उसकी पंक्ति फिर से बनती है। जो अकाउंट विफल होता है या समय पर जवाब नहीं देता, वह सेव की गई संख्याओं पर लौट आता है, इसलिए एक धीमा अकाउंट कभी दूसरों को नहीं रोकता।

आरेख में, हर अकाउंट की खोज उसके अपने goroutine में होती है, और जिसका समय ख़त्म हो जाता है, वह इसके बजाय अपनी सेव की गई संख्याएँ दिखाता है:

<!-- motion:usage-fetch -->

`personal` और `work` समय पर जवाब देते हैं और live संख्याएँ दिखाते हैं। `side` को 3 s से ज़्यादा लगते हैं, इसलिए ccam इसके बजाय उसकी सेव की गई संख्याएँ दिखाता है (`fallback()`)। फिर ccam हर अकाउंट का एक block बनाता है (`usage.Table()`) और आख़िरी पंक्ति में बताया जाने वाला अकाउंट चुनता है (`usage.Best()`)।

एक अकाउंट के लिए ccam ये जाँचें इस क्रम में करता है (`Fetcher.Fetch`):

<!-- motion:usage-cache -->

1. `--source cache` के साथ यह सीधे `fallback()` पर जाता है।
2. जब तक `--no-cache` सेट न हो, `~/.local/state/ccam/usage/<name>.json` में 60 सेकंड (`CacheTTL`) से कम पहले मापा गया live नतीजा इस्तेमाल होता है, अभी के हिसाब से आँककर: जो विंडो तब से रीसेट हो चुकी है वह 0% दिखाती है, और घड़ी से आगे की तारीख़ वाली फ़ाइल ताज़ा नहीं मानी जाती। लगातार दो बार `ccam usage` चलाने पर हर अकाउंट के लिए एक ही request लगती है।
3. ccam `.credentials.json` (या macOS कीचेन) को बदले बिना पढ़ता है और केवल `claudeAiOauth` key decode करता है (`readCredentials()`)। Subscription लॉगिन न होने का मतलब `not logged in` है, या `no limits (API)` जब `.claude.json` में API key हो। जिस लॉगिन का refresh token एक्सपायर हो गया हो, या जिसे Claude Code ने हटा दिया हो, वह `` login expired → `ccam login <name>` `` होता है और उसके लिए खोज नहीं की जाती।
4. एक्सपायर हुआ access token इस्तेमाल नहीं होता। `--refresh` के बिना, अकाउंट अपनी सेव की गई संख्याएँ ``access token expired; `ccam <name>` renews it when it starts`` के साथ दिखाता है।
5. वरना ccam endpoint को बुलाता है। 401 या 403 को एक्सपायर हुआ लॉगिन माना जाता है (`login expired: authentication failed (401)`), कोई दूसरा status या न पढ़ा जा सकने वाला body त्रुटि माना जाता है, और दोनों सेव की गई संख्याओं पर लौटते हैं। Ctrl+C से ख़त्म हुई खोज `cancelled` कहती है। 429 का मतलब है कि endpoint को इस अकाउंट के लिए बहुत ज़्यादा requests मिली हैं: ccam दोबारा पूछने के समय के साथ `<name>.pause` लिखता है (`Retry-After`, या 5 मिनट, ज़्यादा से ज़्यादा एक घंटा), `rate limited by the usage endpoint; asking again after 15:10` के साथ सेव की गई संख्याएँ दिखाता है, और तब तक उस अकाउंट के लिए endpoint को नहीं बुलाता, `--no-cache` के साथ भी। तुरंत दोबारा कोशिश करने से rate limit बस और लंबी हो जाती।
6. अच्छा जवाब मापे जाने के समय के साथ `<name>.json` में लिखा जाता है।

सेव की गई संख्याओं के लिए, ccam सबसे नई सेव की गई खोज लेता है (`<name>.json`, या Claude Code की खोजों से `<name>.claude.json`) और, जब `<name>.statusline.json` उसके बाद लिखी गई हो, तो स्टेटस लाइन की 5 घंटे और साप्ताहिक विंडो उसके ऊपर रखता है (`fallback()`); मॉडल की विंडो, जो स्टेटस लाइन में नहीं होतीं, वैसी ही रहती हैं जैसी खोज ने उन्हें सेव किया था। सेव की गई संख्याएँ दिखाने से पहले यह हर विंडो का रीसेट समय जाँचता है: जो विंडो संख्याएँ सेव होने के बाद रीसेट हो चुकी है, वह 0% दिखाती है, क्योंकि पुराना प्रतिशत अब सही नहीं है। इसीलिए ऊपर `personal` बिना रीसेट समय के `0%` दिखाता है।

Endpoint के दो response formats हैं। नए में `session` (5 घंटे), `weekly_all` और `weekly_scoped` (मॉडल के हिसाब से) entries वाला `limits` array होता है। पुराने में `five_hour`, `seven_day`, `seven_day_opus` और `seven_day_sonnet` ऑब्जेक्ट होते हैं। ccam दोनों पढ़ता है (`ParseLive`): `limits` array को प्राथमिकता मिलती है, और पुराने ऑब्जेक्ट वह भरते हैं जो वह कवर नहीं करता। बिना किसी सीमा वाला body बदले हुए format के रूप में बताया जाता है, जिसके बारे में `ccam doctor` आपको बताता है।

### Limit record

जब Claude Code किसी turn को उपयोग सीमा पर रोकता है, तो [ccam plugin](/docs/install#claude-code-plugin) का hook बातचीत से अस्वीकार हुई request पढ़ता है: कौन-सी विंडो भरी थी (5 घंटे, साप्ताहिक, Opus या Sonnet) और वह कब रीसेट होती है। यह इसे `~/.local/state/ccam/usage/<name>.limit.json` में लिखता है। ccam यह record हर source से मिली बाक़ी हर जानकारी के ऊपर लागू करता है (`Fetch`): इसमें दर्ज हर विंडो अपने reset तक 100% गिनी जाती है, सेव की गई संख्याएँ चाहे कितनी भी पुरानी या कम हों, इसलिए `ccam next --auto`, `ccam pick`, `--fallback` और [`auto_switch`](/features/next#auto-switch) उससे पहले उस अकाउंट को फिर से नहीं चुनते।

जब और कुछ पता न हो, तो अकाउंट केवल यही record दिखाता है (`-v` के साथ `source: the usage limit Claude Code reported`)। बिना निगरानी वाले sessions (`-p` runs और background sessions) कोई record नहीं लिखते, और न ही संगठन की लगाई रोक (उदाहरण के लिए, जब उपयोग ख़त्म हो जाए या किसी ग्रुप की सीमा $0 हो), जिसे उसके admins कभी भी हटा सकते हैं। किसी दूसरे Claude अकाउंट में `ccam login` करने पर record हट जाता है; फ़ाइल मिटाने से भी यह साफ़ हो जाता है।

### Tokens renew करना

Refresh tokens बदलते रहते हैं: एक को renew करने पर नया refresh token मिलता है और पुराना अमान्य हो जाता है। अगर ccam उसी अकाउंट से Claude Code चलते समय लॉगिन renew करता, तो वह session साइन आउट हो जाता। इसलिए `ccam usage` tokens पढ़ता है और कभी नहीं लिखता, जब तक आप `--refresh` न दें। तब ccam (`Fetcher.refresh`):

1. `ccam import` से अपनी जगह पर रजिस्टर किए गए अकाउंट अस्वीकार करता है, क्योंकि सादा `claude` ccam की जानकारी के बिना वह लॉगिन इस्तेमाल कर सकता है;
2. macOS कीचेन में रखे लॉगिन अस्वीकार करता है;
3. एक lock फ़ाइल लेता है, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, ताकि दो ccam processes कभी दोनों एक ही token renew न करें; जब तक यह lock लिया हुआ है, अकाउंट `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)` दिखाता है। `ccam rename`, `remove` और `logout` भी इसे लेते हैं, और renewal चलते रहने तक मना करते हैं;
4. lock के तहत, अकाउंट का कोई session चल रहा हो तो मना करता है, चाहे वह किसी भी तरह शुरू हुआ हो (ccam से, सादे `claude` से, किसी IDE से, या कोई background session और उसके पीछे का daemon, जो लॉगिन ख़ुद renew करता है), और अगर किसी दूसरे process ने token पहले ही renew कर दिया हो तो request छोड़ देता है;
5. जब तक Claude Code अकाउंट फ़ोल्डर में अपने renewal locks (`.oauth_refresh.lock`, `<folder>.lock`) रखे हुए है, कुछ नहीं भेजता: उस समय Claude Code ख़ुद वह लॉगिन renew कर रहा होता है;
6. refresh token को token endpoint पर post करता है और जवाब का 30 सेकंड तक इंतज़ार करता है, उतना ही समय जितना Claude Code अपने renewals को देता है। एक बार भेजी गई request कभी बीच में नहीं रोकी जाती, Ctrl+C से भी नहीं: हो सकता है सर्वर refresh token पहले ही बदल चुका हो, और केवल उसके जवाब में नया लॉगिन हो। Ctrl+C या hangup बाक़ी खोजों को ख़त्म कर देता है, और जवाब सेव होते ही ccam बाहर निकल जाता है (Ctrl+C के बाद 130 के साथ);
7. Claude Code के write lock (`.storage-write.lock`) के तहत सेव करता है और `.credentials.json` में, उस समय वह जैसी है, केवल token वाले fields बदलता है: इस बीच Claude Code ने जो लिखा (कोई MCP साइन-इन) वह बना रहता है, इस बीच Claude Code का renew किया लॉगिन रखा जाता है, और कोई logout कभी पलटा नहीं जाता। जब जवाब नए refresh token की अवधि बताता है, तो उसे `refreshTokenExpiresAt` के रूप में सेव किया जाता है, जैसे Claude Code करता है।

जो refresh पूरा नहीं होता, वह कारण बताता है, अकाउंट की पंक्ति पर या उसकी सेव की गई संख्याओं के ऊपर: उदाहरण के लिए `token not renewed: a session is running`, या `login may be lost: …` जब request भेजी गई लेकिन उसका जवाब इस्तेमाल या सेव नहीं हो सका, इस सलाह के साथ कि अगर इसके बाद Claude Code लॉगिन करने को कहे तो `ccam login <name>` चलाएँ।

### स्टेटस लाइन collector

Claude Code स्टेटस लाइन कमांड चलाता है और उसे pipe से एक JSON document देता है। ccam उससे मॉडल का नाम और `rate_limits` पढ़ता है:

<!-- motion:usage-statusline -->

1. ccam अकाउंट को `CCAM_PROFILE` से खोजता है, जिसे ccam Claude Code शुरू करते समय सेट करता है, या वरना `CLAUDE_CONFIG_DIR` का अकाउंट फ़ोल्डरों से मिलान करके (`runStatusline`)।
2. ccam `rate_limits.five_hour` और `rate_limits.seven_day` पढ़ता है (`ParseStatusLine`)। API key वाले अकाउंट्स के लिए, या session के पहले response से पहले, JSON में rate limits नहीं होतीं; तब कुछ सेव नहीं होता।
3. ccam `~/.local/state/ccam/usage/<name>.statusline.json` लिखता है (`SaveStatusLine`)। जब endpoint जवाब नहीं देता, तो `ccam usage` इसे इस्तेमाल करता है।
4. ccam `[name] · model · 5h N% · 7d N%` प्रिंट करता है। `--wrap` के साथ, यह आपकी कमांड को उसी input के साथ `sh -c` से चलाता है और उसका output उसी तरह प्रिंट करता है जैसे Claude Code स्टेटस लाइन दिखाता है: हर वह पंक्ति जो ख़ाली नहीं है, पहली पंक्ति के आगे `[name]` के साथ। विफल होने वाली कमांड केवल अकाउंट का नाम दिखाती है। Windows पर यह कमांड को उसी तरह चलाता है जैसे Claude Code: Git Bash इंस्टॉल हो तो उससे, वरना PowerShell से।
5. जब 5 घंटे या साप्ताहिक उपयोग 90% या उससे ज़्यादा हो और कोई दूसरा अकाउंट मौजूद हो, तो पंक्ति `→ ccam next --auto --continue` पर ख़त्म होती है, उदाहरण के लिए `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`।

Wrap की गई कमांड की अपनी कोई समय सीमा नहीं है: Claude Code स्टेटस लाइन कमांड को 10 मिनट देता है और नया अपडेट शुरू होने पर चल रहा run रोक देता है, इसलिए धीमी कमांड वैसे ही दिखती है जैसे ccam के बिना दिखती। जब Claude Code ccam को बंद करता है, तो कमांड भी बंद हो जाती है (Linux और macOS पर उसके शुरू किए हर process के साथ)।

`ccam statusline install` हर अकाउंट की `settings.json` में `statusLine` लिखता है और ऑब्जेक्ट की बाक़ी keys, जैसे `padding`, `refreshInterval` और `hideVimModeIndicator`, बनाए रखता है। जब किसी अकाउंट के पास पहले से स्टेटस लाइन कमांड हो, तो ccam उसे `--wrap '<old command>'` के रूप में रखता है; Windows पर `--wrap-base64 <base64>` के रूप में, जिसे किसी भी shell में quoting की ज़रूरत नहीं होती। अगर स्टेटस लाइन पहले से ccam चलाती है (किसी दूसरे पाथ से, या उस तरह quote करके जैसे किसी पुराने वर्ज़न ने लिखा था), तो ccam उसे wrap करने के बजाय बदल देता है, और जिस कमांड को वह wrap करती है उसे बनाए रखता है।

Windows पर ccam का पाथ forward slashes और home फ़ोल्डर के लिए `~` के साथ लिखा जाता है (`~/.local/bin/ccam.exe statusline`), जिसे Git Bash और PowerShell दोनों समझते हैं। जिस पाथ को quotes की ज़रूरत हो, वह Claude Code जिस shell को इस्तेमाल करता है उसके लिए पूरा लिखा जाता है: Git Bash के लिए `'…'`, PowerShell के लिए `& '…'`। ग्रुप के ज़रिए साझा की गई `settings.json` उसके लिंक से एक बार लिखी जाती है, इसलिए ग्रुप के हर अकाउंट को यह मिल जाता है। इसे फिर से चलाने पर, जो फ़ाइलें पहले से ccam इस्तेमाल करती हैं उनके लिए `(already set)` प्रिंट होता है।

## सुरक्षा और विशेष स्थितियाँ

- जब तक आप `--refresh` न दें, `ccam usage` कभी token फ़ाइल नहीं बदलता। जब कोई access token एक्सपायर हो जाए, तो अकाउंट को एक बार `ccam <name>` से चलाएँ, Claude Code उसे renew कर देगा। जो लॉगिन पूरी तरह एक्सपायर हो गया हो, उसके लिए `ccam login <name>` चाहिए।
- `--refresh` कभी चल रहे अकाउंट्स (तब भी, जब Claude Code ccam के बिना शुरू हुआ हो), import किए गए अकाउंट्स या कीचेन लॉगिन को नहीं छूता। अस्वीकार किया गया refresh सेव की गई संख्याएँ और कारण (`token not renewed: …`) दिखाता है।
- Usage endpoint सार्वजनिक रूप से documented नहीं है और बदल सकता है। ऐसा होने पर भी आपको सेव की गई संख्याएँ दिखती हैं, और `ccam doctor` format में बदलाव बताता है।
- सेव की गई संख्याएँ और जवाब दोनों न होने पर, अकाउंट की पंक्ति `lookup failed:` और कारण दिखाती है, उदाहरण के लिए `--source cache` के साथ `no saved usage`।
- Cache फ़ाइलों में केवल प्रतिशत और रीसेट समय होते हैं, कभी tokens नहीं।

## संबंधित

- [उपयोग सीमाएँ और अकाउंट स्विच करना](/docs/usage-limits)
- [अगले अकाउंट पर स्विच करना](/features/next)
- [Claude Code की usage limit पूरी होने पर क्या करें](/guides/usage-limit-workflow)
- [कमांड संदर्भ](/docs/commands)
