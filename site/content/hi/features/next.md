---
title: अगले अकाउंट पर स्विच करना
description: ccam next किसी दूसरे अकाउंट से Claude Code खोलता है; --auto सबसे ज़्यादा गुंजाइश वाला चुनता है, --continue बातचीत आगे ले जाता है, --fallback -p jobs दोहराता है।
updated: 2026-11-01
order: 8
translation: machine
sourceHash: "79300c521748"
---

`ccam next` आपके नाम बताए बिना Claude Code को किसी दूसरे अकाउंट से खोलता है। इसे तब इस्तेमाल करें जब कोई अकाउंट अपनी उपयोग सीमा तक पहुँच जाए। स्विच करने के चार तरीक़े हैं:

- सादा `ccam next` पिछली बार इस्तेमाल किए गए अकाउंट के बाद वाला अकाउंट लेता है, और `--auto` वह लेता है जिसमें सबसे ज़्यादा गुंजाइश हो। `--continue` इस फ़ोल्डर की आख़िरी बातचीत को नए अकाउंट पर जारी रखता है, चाहे दोनों इतिहास साझा करते हों या नहीं।
- Claude Code के अंदर, [ccam plugin](#inside-claude-code) बताता है कि कब कोई सीमा session को रोकती है, और `/ccam:switch` बातचीत आगे सौंप देता है।
- `claude -p` चलाने वाली scripts के लिए, `ccam --fallback` यह स्विच अपने-आप करता है ([सीमा पर फिर से चलाना](#run-again-at-the-limit) देखें)।
- [`auto_switch`](#auto-switch) चालू होने पर (Linux और macOS), ccam उस interactive session को भी दूसरे अकाउंट पर ले जाता है जिसे Claude Code उसकी सीमा पर रोक देता है।

**Team और Enterprise सीटें।** अपने-आप अकाउंट चुनना (`ccam next`, `--fallback`, `ccam pick`, डैशबोर्ड का `n` और, चालू करने के बाद, [`auto_switch`](/features/next#auto-switch)) एक ही Claude Team या Enterprise संगठन की सीटों के बीच भी स्विच करता है, जो Anthropic की शर्तों का उल्लंघन हो सकता है। इस पर भरोसा करने से पहले [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) पढ़ें और अपने संगठन और Anthropic से पता कर लें।

## उपयोग

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

ccam स्विच को standard error पर `→ <from> → <to>` के रूप में प्रिंट करता है, `--auto` वाले चुनाव के कारण के साथ, फिर नए अकाउंट से Claude Code उसी तरह शुरू करता है जैसे `ccam <name>` करता है। `<from>` वह अकाउंट है जिसे आपने पिछली बार इस्तेमाल किया, या `--continue` के साथ वह अकाउंट जिसकी बातचीत है। `--` के बाद वाले आर्ग्युमेंट बिना बदलाव के `claude` को जाते हैं।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--auto` | उपयोग के हिसाब से लॉगिन किया हुआ अकाउंट चुनें (नीचे देखें), 100% वाले अकाउंट छोड़कर | बंद |
| `--max-usage <1-100>` | `--auto` के साथ: इस उपयोग से नीचे वाले अकाउंट्स को प्राथमिकता दें | `auto_max_usage`, वरना 100 |
| `--strategy headroom\|expiring\|balance` | `--auto` के साथ: सबसे कम 5 घंटे का उपयोग, सबसे जल्दी होने वाला 5 घंटे का reset या प्रति दिन सबसे ज़्यादा बचा साप्ताहिक कोटा | `auto_strategy`, वरना `headroom` |
| `--cross-org` | `--auto` के साथ: किसी दूसरे संगठन के अकाउंट भी चुनें (नीचे देखें) | `cross_org`, वरना बंद |
| `--wait[=<duration>]` | `--auto` के साथ: जब किसी अकाउंट में गुंजाइश न हो, तो पहले अकाउंट के reset होने का इंतज़ार करें, ज़्यादा से ज़्यादा इतनी देर। मान `=` के बाद आता है: `--wait 2h` इस संकेत के साथ अस्वीकार होता है: `--wait takes its value after '=': --wait=2h` | बंद; अकेले देने पर `5h` |
| `-c`, `--continue` | इस फ़ोल्डर की आख़िरी बातचीत को नए अकाउंट पर जारी रखें (नीचे देखें) | बंद |
| `--brief` | `--continue` के साथ: पूरी बातचीत resume करने के बजाय, नए अकाउंट से एक नई बातचीत शुरू करें, जो पिछली बातचीत के ccam के सारांश से शुरू होती है ([नीचे देखें](#brief)) | बंद |
| `-- <claude args>` | ccam के अपने आर्ग्युमेंट के बाद `claude` को दिए जाने वाले आर्ग्युमेंट | कोई नहीं |
| `--dry-run` | समर्थित: Claude Code शुरू किए बिना या शुरुआत दर्ज किए बिना `→ <from> → <to>` और चलाई जाने वाली कमांड लाइन प्रिंट करता है | बंद |

`--auto` और `--continue` एक साथ काम करते हैं। `--continue` के साथ, `--resume` वाले आर्ग्युमेंट आपके अपने आर्ग्युमेंट से पहले रखे जाते हैं।

Claude Code के अपने विकल्प, जैसे `--dangerously-skip-permissions` या `--model`, `--` के बाद आते हैं। `ccam <name>` के उलट, जो नाम के बाद की हर चीज़ आगे भेज देता है, `ccam next` अपने आर्ग्युमेंट ख़ुद पढ़ता है, इसलिए `--` के बिना यह `ccam: unknown flag: --dangerously-skip-permissions` के साथ रुक जाता है।

## यह कैसे काम करता है

ccam से होने वाली हर शुरुआत समय को `~/.local/state/ccam/last-used/<name>` में लिखती है। `ccam next` इन फ़ाइलों को पढ़ता है, जाँचता है कि कौन से अकाउंट लॉगिन किए हुए हैं, और चुनता है।

1. ccam `config.toml` से हर अकाउंट के लिए एक `next.Candidate` बनाता है (`cli.candidates`): अकाउंट, उसके आख़िरी बार शुरू होने का समय (`launcher.LastUsed`) और यह कि उसका लॉगिन काम करता है या नहीं (`creds.Has`: `.credentials.json` में, या macOS पर कीचेन में, ऐसा Claude लॉगिन जो पूरी तरह एक्सपायर नहीं हुआ और जिसे Claude Code ने हटाया नहीं)। जो फ़ोल्डर इसके बजाय API key इस्तेमाल करता है (Console लॉगिन, या `ccam add --api` से जोड़ा गया अकाउंट), उसकी कोई प्लान सीमा नहीं होती और उसे कभी नहीं चुना जाता।
2. `--auto` के बिना, ccam अकाउंट्स को स्लॉट नंबर से क्रम में रखता है और सबसे हाल में शुरू हुआ अकाउंट खोजता है (`next.Rotate`)। यही वह अकाउंट है जिससे आप स्विच कर रहे हैं। अगर अभी तक कोई अकाउंट शुरू नहीं हुआ, तो यह सबसे छोटा स्लॉट होता है।
3. वहाँ से यह स्लॉट क्रम में आगे बढ़ता है, आख़िरी स्लॉट से वापस पहले पर घूमते हुए, और पहला लॉगिन किया हुआ अकाउंट लौटाता है।
4. ccam चुने गए अकाउंट को ठीक वैसे तैयार करके शुरू करता है जैसे `ccam <name>` करता (`launcher.LaunchWith`): यह अकाउंट के लिंक जाँचता है, MCP साझा करना लागू करता है, session रजिस्टर करता है और `syscall.Exec` के ज़रिए ccam की जगह `claude` रखता है।

खोज पिछली बार इस्तेमाल किए गए अकाउंट से शुरू होती है और उसके बाद पहले लॉगिन किए हुए अकाउंट पर रुकती है:

<!-- motion:next-rotate -->

`work` सबसे हाल में शुरू हुआ था, इसलिए खोज उसके बाद से शुरू होती है। `side` का लॉगिन नहीं है और उसे छोड़ दिया जाता है, और `lab` चुना जाता है। अगर `lab` भी लॉग आउट होता, तो खोज घूमकर `personal` पर आती (आरेख में `% len`)। रोटेशन कभी वह अकाउंट नहीं लौटाता जिससे आपने शुरू किया: कोई दूसरा लॉगिन किया हुआ अकाउंट न होने पर यह त्रुटि के साथ रुकता है।

### --auto

`--auto` स्लॉट क्रम की जगह उपयोग से तय करता है। ccam लॉगिन किए हुए अकाउंट रखता है, उसी तरह जैसे `ccam usage` करता है, एक साथ उन सबका उपयोग पूछता है (`Fetcher.FetchAll`), और नीचे दिए नियमों से चुनता है (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

हर अकाउंट को उसकी 5 घंटे और साप्ताहिक विंडो से आँका जाता है, और जब Claude Code के आर्ग्युमेंट में `--model` हो (`--` के बाद), तो उस मॉडल की साप्ताहिक विंडो से भी (नाम में `opus`, `sonnet`, `haiku` या `fable`)। ccam अकाउंट्स को तीन स्तरों में बाँटता है (`Policy.Choose()`):

1. **गुंजाइश**: हर विंडो `auto_max_usage` (`config.toml`, डिफ़ॉल्ट 100) से नीचे। इनमें से strategy चुनती है: `headroom` (डिफ़ॉल्ट) सबसे कम 5 घंटे का उपयोग लेती है; `expiring` वह 5 घंटे की विंडो लेती है जो सबसे जल्दी reset होती है, ताकि उसका बचा हुआ quota ख़त्म होने से पहले इस्तेमाल हो जाए; `balance` वह अकाउंट लेती है जिसके हफ़्ते के reset होने तक हर दिन के हिसाब से सबसे ज़्यादा साप्ताहिक quota बचा हो, और उस अकाउंट पर चल रहे हर ccam session को उसके ख़िलाफ़ गिनती है, ताकि साप्ताहिक quota बिना इस्तेमाल के न रह जाए और समानांतर sessions बँट जाएँ।
2. **सीमा से ऊपर**: 100% से नीचे, लेकिन `auto_max_usage` पर या उससे ऊपर। केवल तब इस्तेमाल होते हैं जब किसी अकाउंट में गुंजाइश न हो; सबसे कम 5 घंटे के उपयोग वाला पहले आता है।
3. **भरे हुए**: कोई विंडो 100% पर, या अज्ञात 5 घंटे का उपयोग। कभी नहीं चुने जाते।

जब Claude Code ने किसी विंडो में कोई request अस्वीकार की हो, तो वह विंडो अपने reset तक 100% गिनी जाती है: [ccam plugin](#inside-claude-code) का hook इसे दर्ज करता है (`~/.local/state/ccam/usage/<name>.limit.json`), सेव की गई संख्याएँ चाहे जो कहें। इसलिए जो अकाउंट अभी-अभी अपनी सीमा पर पहुँचा है, वह अपने reset से पहले फिर से नहीं चुना जाता, भले ही उसकी संख्याएँ पुरानी हों।

सभी अकाउंट्स की सेटिंग्स `config.toml` में सबसे ऊपर, पहले `[[group]]`, `[[profile]]` या `[[folder]]` से ऊपर आती हैं। इनमें से किसी के नीचे लिखी सेटिंग उसी table की होती है: ccam हर launch पर उसके बारे में चेतावनी देता है, और config सेव करने वाली कमांड तब तक मना करती हैं जब तक उसे ऊपर न ले जाया जाए।

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**संगठन।** अपने-आप होने वाले चुनाव उसी तरह के अकाउंट्स के बीच रहते हैं जिस तरह का अकाउंट आप छोड़ रहे हैं: उसी Team या Enterprise संगठन की सीटें, या personal प्लान (Pro, Max), ताकि किसी कंपनी अकाउंट का काम किसी निजी अकाउंट पर, या इसका उल्टा, जारी न हो। ccam हर अकाउंट का संगठन उसके लॉगिन से जानता है (`ccam login` इसे दर्ज करता है; पुराने लॉगिन `.claude.json` से पढ़े जाते हैं); जिस अकाउंट का संगठन ccam नहीं बता पाता, उसे किसी भी अकाउंट के साथ मिलाया जा सकता है। एक ही Claude अकाउंट में लॉगिन दो फ़ोल्डर उसकी सीमाएँ साझा करते हैं, इसलिए उनमें से एक कभी दूसरे की जगह नहीं लेता। `--cross-org`, या `config.toml` में `cross_org = true`, संगठन वाला नियम हटा देता है, लेकिन [`/ccam:switch`](#inside-claude-code) के बाद किसी दूसरे संगठन में स्विच अब भी आपकी हाँ का इंतज़ार करता है।

जब कुछ भी नहीं चुना जा सकता, तो ccam code 75 (`EX_TEMPFAIL`) के साथ बाहर निकलता है, ताकि script "बाद में फिर कोशिश करें" को विफलता से अलग पहचान सके। `--wait` के साथ, यह इसके बजाय तब तक रुकता है जब तक पहले भरे हुए अकाउंट की रोकने वाली विंडो reset न हो जाए (reset समय के 30 सेकंड बाद, जब endpoint नई विंडो दिखाता है), `ccam: no account has room; waiting until 15:10` प्रिंट करता है, और फिर से चुनता है; जब कोई अकाउंट नहीं बताता कि वह कब reset होगा, तो यह 10 मिनट बाद फिर देखता है। अगर इंतज़ार ख़त्म होने से पहले कुछ ख़ाली नहीं होता, तो यह 75 के साथ बाहर निकलता है।

यह खोज `ccam usage` के नियमों का पालन करती है: 60 सेकंड से कम पुरानी संख्याएँ दोबारा इस्तेमाल होती हैं, हर अकाउंट की खोज 3 सेकंड के बाद time out हो जाती है, जो अकाउंट जवाब नहीं देता, या जिसका access token एक्सपायर हो गया है, उसे उसकी सेव की गई संख्याओं से आँका जाता है, और limit record उन संख्याओं की जगह लेता है। सेव की गई जिस विंडो का रीसेट समय बीत चुका है, उसे 0% माना जाता है, ताकि जिस अकाउंट की विंडो तब से रीसेट हो चुकी है, वह छूट न जाए। `--auto` कभी tokens renew नहीं करता। सादे रोटेशन के उलट, `--auto` पिछली बार इस्तेमाल किया गया अकाउंट भी चुन सकता है, अगर उसमें अब भी सबसे ज़्यादा गुंजाइश हो।

### --continue

1. ccam हर अकाउंट में इस फ़ोल्डर की आख़िरी बातचीत खोजता है: `projects/<folder>/` की `<session>.jsonl` फ़ाइलों में से वह जिसमें सबसे नया prompt या जवाब हो (फ़ोल्डर का नाम काम करने वाला फ़ोल्डर होता है, links resolve करके, जिसमें अक्षर या अंक के अलावा हर character `-` में बदल जाता है, जैसे Claude Code इसका नाम रखता है)। जिन sessions में ऐसा कुछ नहीं है, जैसे ऐसा session जिसने केवल `/mcp`, `/usage` या `/status` चलाया, या ख़ाली फ़ाइल, उन्हें छोड़ दिया जाता है। अगर किसी session ने कोई नोट छोड़ा है (उपयोग सीमा या `/ccam:switch`, नीचे देखें) और किसी बातचीत में उसके 2 मिनट से ज़्यादा बाद का कोई prompt या जवाब नहीं है, तो इसके बजाय ठीक वही session इस्तेमाल होता है; Claude Code `/exit` पर जो पंक्तियाँ लिखता है, वे नहीं गिनी जातीं।
2. जिस अकाउंट की वह बातचीत है, उसी से आप स्विच करते हैं, और नया अकाउंट ऊपर बताए तरीक़े से, उसी संगठन के भीतर, चुना जाता है। जब कई अकाउंट वह इतिहास साझा करते हैं, तो बातचीत उस अकाउंट की मानी जाती है जिसकी `.claude.json` बताती है कि उसने वह बातचीत आख़िरी बार चलाई (`lastSessionId`), नहीं तो उस अकाउंट की जो आख़िरी बार इस्तेमाल हुआ।
3. अगर नया अकाउंट वह इतिहास साझा करता है (वही ग्रुप फ़ोल्डर), तो ccam `claude --resume <session id>` शुरू करता है, यानी ठीक वही बातचीत (`--continue` फ़ोल्डर की सबसे नई बातचीत खोलता, चाहे वह कोई भी हो)।
4. नहीं तो यह `claude --resume <the conversation's file> --fork-session` शुरू करता है: Claude Code पूरी बातचीत को एक नए session ID के तहत नए अकाउंट के अपने इतिहास में कॉपी करता है, और मूल बातचीत जैसी थी वैसी रहती है। उसके बाद से नए अकाउंट का `--continue` उसे ढूँढ लेता है।
5. जब तक आप `--` के बाद `--agent` न दें, ccam वह agent भी आगे ले जाता है जिसके रूप में बातचीत चलती है: फ़ोल्डर की `.claude/settings.local.json` या `.claude/settings.json` की, या अकाउंट की `settings.json` की `agent` सेटिंग, नहीं तो वह agent जिसका नाम बातचीत की आख़िरी `agent-setting` पंक्ति में है। Claude Code बातचीत कॉपी करते समय (चरण 4) agent को वापस नहीं लाता, और अपने डिफ़ॉल्ट agent के रूप में, हर tool के साथ, आगे चलता। जब नया अकाउंट पुराने अकाउंट के agents, plugins और settings साझा नहीं करता (किसी ग्रुप के ज़रिए वही `agents` और `plugins` फ़ोल्डर और `settings.json`), तो ccam स्विच वाली पंक्ति के बाद एक पंक्ति जोड़ता है: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`।

ccam बताता है कि बातचीत कितनी बड़ी है: नया अकाउंट अपने पहले जवाब के साथ उसे पूरा फिर से पढ़ता है। Prompt caches कभी संगठनों के बीच साझा नहीं होते, और हर personal प्लान अपना अलग संगठन है, इसलिए जब दोनों अकाउंट अलग-अलग संगठनों के हों, तो ccam `prompt caches are not shared between organizations, so the first reply reads it all again` जोड़ता है। जब पुराना अकाउंट ऐसे MCP सर्वरों में साइन इन हो जिनमें नया अकाउंट नहीं है, तो ccam उनके नाम बताता है: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`। किसी ग्रुप में `history` साझा करने से कॉपी की ज़रूरत नहीं पड़ती, एक ही बातचीत रहती है, और वे checkpoints भी साझा होते हैं जिनसे `/rewind` कोड वापस लाता है (`file-history/`)।

### --brief

`--continue` पूरी बातचीत resume करता है, और नया अकाउंट अपने पहले जवाब के साथ उसे पूरा फिर से पढ़ता है। `ccam next --continue --brief` इसके बजाय एक नई बातचीत शुरू करता है, उस सारांश से जो ccam किसी मॉडल से पूछे बिना पुरानी बातचीत का लिखता है:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

सारांश में, इसी क्रम में और केवल वहाँ जहाँ कुछ कहने को हो, ये होते हैं:

1. बातचीत का शीर्षक;
2. उसके transcript का path, ताकि Claude को किसी ब्योरे की ज़रूरत पड़ने पर वह उसमें खोज सके;
3. आपकी पहली request, जब चार से ज़्यादा रही हों;
4. वह सारांश जो Claude Code ने बातचीत को पिछली बार compact करते समय लिखा;
5. आपकी आख़िरी चार requests, `!` commands और slash commands के साथ, जैसे आपने उन्हें टाइप किया;
6. Claude का आख़िरी जवाब;
7. todo सूची;
8. वे 10 फ़ाइलें जिन्हें Claude ने सबसे आख़िर में बदला;
9. फ़ोल्डर का `git status --short --branch`, 800 characters तक।

Todo सूची में ज़्यादा से ज़्यादा 10 items रहते हैं: पूरे हो चुके items पहले हटते हैं, फिर सबसे आगे के steps, ताकि मौजूदा step और उसके बाद वाले steps बने रहें। लंबे हिस्से बीच से काटे जाते हैं, और पूरा सारांश 16,000 characters से कम रहता है।

ccam इसे `~/.local/state/ccam/brief/<session>.md` के रूप में सेव करता है, जिसे केवल आप पढ़ सकते हैं, और अगला `--brief` 7 दिन से पुराने सारांश हटा देता है। Claude Code एक prompt के साथ शुरू होता है, और उसमें बताई गई फ़ाइल जोड़ लेता है: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` `--` के बाद दिया गया आपका prompt आख़िरी वाक्य की जगह लेता है (`Then: <your prompt>`), और `--` के बाद आपके दूसरे विकल्प उसके बाद आते हैं।

Claude Code किसी mention को `"` पर ख़त्म करता है और `#` पर काट देता है, इसलिए जब सारांश के path में इनमें से कोई हो, तो ccam उसे काम करने वाले फ़ोल्डर के सापेक्ष उसके path के ज़रिए mention करता है; जब उस path में भी कोई हो, तो prompt फ़ाइल का नाम बताता है, और Claude उसे अपने Read tool से पढ़ता है, जो पहले पूछता है क्योंकि फ़ाइल काम करने वाले फ़ोल्डरों से बाहर है (Manual mode में `-p` run के दौरान, Claude Code पढ़ने से मना कर देता है और Claude सारांश के बिना आगे बढ़ता है)।

`ccam --dry-run next --continue --brief` सारांश प्रिंट करता है और कुछ सेव नहीं करता। पुरानी बातचीत जैसी थी वैसी रहती है, और नई बातचीत केवल वही जानती है जो सारांश में है, इसलिए `--brief` तब इस्तेमाल करें जब काम, वहाँ तक पहुँचने के हर ब्योरे से ज़्यादा मायने रखता हो। `--continue` के बिना, ccam `--brief works with --continue` के साथ रुक जाता है; जब वह बातचीत पढ़ नहीं पाता या उसमें सारांश लायक़ कुछ नहीं मिलता, तो वह रुकता है और बताता है कि `--brief` के बिना `ccam next --continue` उसे पूरा resume करता है। [Relay](#relay) और [`auto_switch`](#auto-switch) हमेशा पूरी बातचीत resume करते हैं। Claude Code का अपना `--brief` विकल्प, उसके दूसरे विकल्पों की तरह, `--` के बाद आता है।

## Claude Code के अंदर

[ccam plugin](/docs/install#claude-code-plugin) दो hooks जोड़ता है (Claude Code 2.1.78 या उसके बाद):

- जब कोई turn उपयोग सीमा पर रुकता है, तो यह बातचीत को दर्ज करता है और उन टर्मिनलों में desktop notification दिखाता है जो इसे समर्थन करते हैं (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  यह हर वह संदेश पहचानता है जिसे Claude Code उपयोग सीमा मानता है: 5 घंटे की, साप्ताहिक, सामान्य और मॉडल की सीमाएँ (`You've hit your session limit`, `You've hit your limit`, किसी मॉडल के हिस्से के लिए `You've reached your …`), usage credits और खर्च की सीमाएँ, और ऐसा संगठन जिसने अकाउंट को रोक दिया हो।

  जब Claude Code ने अस्वीकार हुई request बातचीत में लिख दी हो, तो hook भरी हुई विंडो को उसके reset तक दर्ज भी करता है ([limit record](/features/usage#the-limit-record)), सिवाय संगठन की लगाई रोक के, जिसे उसके admins कभी भी हटा सकते हैं। सर्वर का थोड़ी देर का throttle या किसी API key की rate limit उपयोग सीमा नहीं है और उसे छोड़ दिया जाता है, और इसी तरह किसी subagent का रुकना और बिना निगरानी वाले sessions भी (`-p` runs और background sessions)।

  Claude Code 2.1.234 और उसके बाद के वर्ज़न session में इंतज़ार भी करते हैं और 5 घंटे के reset के बाद अपने-आप जारी रखते हैं; notification आपको इसके बजाय स्विच करने का विकल्प देता है, और [`auto_switch`](#auto-switch) चालू होने पर (Linux और macOS), जब किसी दूसरे अकाउंट में गुंजाइश हो तो ccam आपके लिए स्विच कर देता है। जब कोई दूसरा अकाउंट जगह नहीं ले सकता, क्योंकि कोई दूसरा अकाउंट उसी Team या Enterprise संगठन का नहीं है, या कोई दूसरा personal प्लान नहीं है (`cross_org` चालू होने पर कोई भी दूसरा Claude अकाउंट ले सकता है), तो notification इसके बजाय `You have no other personal plans, so ccam cannot move this conversation.` या `You have no other accounts in Acme, so ccam cannot move this conversation.` पर ख़त्म होता है।
- `/ccam:switch [account]` बातचीत को दर्ज करता है और कोई request भेजे जाने से पहले turn ख़त्म कर देता है: `/exit` टाइप करें, फिर उसे आपके बताए अकाउंट पर (या अगले अकाउंट पर) जारी रखने के लिए `ccam next --continue` चलाएँ। जिस session में अभी कोई बातचीत नहीं है, उसमें यह बताता है कि जारी रखने के लिए कुछ नहीं है। इसे केवल वही चलाता है जो आप टाइप करते हैं: plugin 0.3.2 से, Claude `/ccam:switch` नहीं चला सकता, और Claude की शेड्यूल की हुई job (CronCreate, ScheduleWakeup या `.claude/scheduled_tasks.json`) इसे Claude को सादे टेक्स्ट के रूप में देती है।

  जब आपका बताया अकाउंट बातचीत के संगठन से बाहर हो (किसी दूसरे Team या Enterprise संगठन की सीट, सीट की जगह personal प्लान या इसका उल्टा, या किसी भी तरफ़ API key अकाउंट), तो `cross_org` चाहे जो कहे, `ccam next --continue` पहले पूछता है: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`। जब आप कोई नाम नहीं देते, तब ख़ुद चुने गए अकाउंट से पहले भी यह यही पूछता है, जिसे `cross_org = true` होने पर यह किसी दूसरे संगठन से ले सकता है: `ccam: the account picked for /ccam:switch is home, a personal plan, …`। केवल `y` या `yes` बातचीत को ले जाता है; Linux और macOS पर, सवाल से पहले टाइप की गई keys पहले हटा दी जाती हैं। कोई दूसरा जवाब, या पूछने के लिए टर्मिनल न होना, बातचीत और नोट को वहीं छोड़ देता है जहाँ वे हैं: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (अपने-आप चुनाव के बाद, `ccam next --continue` फिर से चुनता है)। अपने-आप होने वाले चुनावों की तरह, जिस अकाउंट का संगठन ccam नहीं बता पाता, वह किसी के साथ भी मेल खाता है।

नोट `~/.local/state/ccam/handoff/<account>/<session>.json` में रहता है, हर session के लिए एक फ़ाइल, ताकि एक अकाउंट के दो sessions कभी एक-दूसरे का नोट न मिटाएँ। यह 12 घंटे तक रखा जाता है और एक बार इस्तेमाल होता है, जब अगला session शुरू होने वाला हो; जो शुरुआत विफल हो जाए, वह नोट को अगली कोशिश के लिए छोड़ देती है। नोट उसी फ़ोल्डर का होता है जिसमें session शुरू हुआ था, Claude Code के shell में `cd` करने के बाद भी। `ccam rename` उसमें दिए नाम अपडेट करता है; अगर नोट में जिस अकाउंट का नाम है, उसे तब से हटा दिया गया हो, तो ccam यह बताता है और हमेशा की तरह चुनता है।

### Relay

`config.toml` में सबसे ऊपर `relay = true` होने पर ccam ख़ुद को Claude Code से बदलने के बजाय interactive sessions के साथ-साथ चलता रहता है, और `/ccam:switch` के बाद उस session की बातचीत अपने-आप जारी रखता है (किसी दूसरे session के छोड़े नोट पर यह कुछ नहीं करता): hook Claude Code को बंद करता है, और ccam उस session के नोट के साथ `ccam next --continue` चलाता है। अगर किसी अकाउंट का नाम नहीं दिया गया, तो यह बातचीत के मॉडल को गिनते हुए सबसे ज़्यादा गुंजाइश वाला अकाउंट लेता है, और जब `cross_org = true` के कारण वह किसी दूसरे संगठन का अकाउंट हो सकता है, तो बातचीत को वहाँ ले जाने से पहले यह ऊपर वाला सवाल पूछता है। Windows पर hook Claude Code को बंद नहीं कर सकता: `/exit` टाइप करें, और ccam उसे जारी रखता है। बातचीत के संगठन से बाहर के अकाउंट के लिए भी hook उसे बंद नहीं करता: यह `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).` कहता है, और आपके `/exit` टाइप करने के बाद ccam ऊपर वाला सवाल पूछता है। इसके लिए हर session पर एक इंतज़ार करता ccam process लगता है, इसीलिए यह डिफ़ॉल्ट रूप से बंद है।

आगे ले जाया गया session वह सब बनाए रखता है जो बातचीत को resume करने से वापस नहीं आता: `--add-dir`, `--mcp-config`, tools की सूचियाँ, `--settings`, `--plugin-dir` और system prompt वाले विकल्प, बातचीत का permission mode, session का effort level और वह agent जिसके रूप में बातचीत चलती है ([`--continue`](#continue) देखें)। `--model` और `--name` को Claude Code ख़ुद बातचीत से वापस लाता है। बाक़ी सब कुछ (`--worktree`, `--debug`, `--remote-control`, ऐसे विकल्प जिन्हें ccam नहीं जानता) `ccam: not carried over: …` में बताया जाता है। स्विच कभी साझा इतिहास वाला सवाल नहीं पूछता, और उसके चलते समय आपने जो टाइप किया, वह छोड़ दिया जाता है; उसका इकलौता सवाल वह है जो किसी दूसरे संगठन के अकाउंट से पहले आता है।

Relay बातचीत को तभी जारी रखता है जब Claude Code उसी तरह बंद हुआ हो जैसे `/ccam:switch` ने उससे कहा था: तब नहीं जब आपने पहले ही `/exit` टाइप कर दिया, उसके बंद होते समय Ctrl+C दबाया, या ccam को hangup या कोई दूसरा signal मिला, और तब भी नहीं जब टर्मिनल ही न रहा हो। तब नोट बना रहता है, और `ccam next --continue` उसे बाद में जारी रखता है। ccam के किसी पुराने वर्ज़न के relay में शुरू हुआ session, या बिना निगरानी वाला session (`-p` runs और background sessions), hook से बंद नहीं होता; उसे इसके बजाय `/exit` वाला संदेश मिलता है। `claude remote-control` बिना relay के चलता है।

### `auto_switch`

`config.toml` में सबसे ऊपर `auto_switch = true` होने पर (Linux और macOS), जब Claude Code किसी session को उपयोग सीमा पर रोकता है, तो ccam उसे अपने-आप दूसरे अकाउंट पर ले जाता है: plugin का hook गुंजाइश वाला कोई दूसरा अकाउंट चुनता है और Claude Code को बंद करता है, और session के साथ-साथ चल रहा ccam बातचीत को वहाँ जारी रखता है, उसी तरह जैसे [relay](#relay) `/ccam:switch` के बाद करता है। यह डिफ़ॉल्ट रूप से बंद है और प्रायोगिक है।

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam किसी session को केवल तभी ले जाता है जब ये सभी बातें सही हों:

1. ccam ने उसे किसी टर्मिनल में शुरू किया और उसके साथ-साथ चलता रहता है। `auto_switch` [relay](#relay) को चालू करता है, इसलिए यह उन interactive sessions के लिए सही है जिन्हें आप इसे सेट करने के बाद शुरू करते हैं। उससे पहले, या किसी पुराने ccam में, शुरू हुए session को ऐसा notification मिलता है जो ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` पर ख़त्म होता है। बिना निगरानी वाले sessions (`-p` runs और background sessions), `claude remote-control` और ccam के बिना शुरू हुआ Claude Code कभी नहीं ले जाए जाते।
2. अकाउंट में [ccam plugin](/docs/install#claude-code-plugin) चालू है और वह Claude Code 2.1.214 या उसके बाद का वर्ज़न चलाता है, जो hook को अपना process ID देता है (`CLAUDE_PID`)।
3. Claude Code ने अस्वीकार हुई request को, उसके reset होने के समय के साथ, बातचीत में लिखा, और वह उस अकाउंट की 5 घंटे, साप्ताहिक, Opus या Sonnet सीमा के लिए थी। खर्च की सीमा, usage credits ख़त्म होना, संगठन की लगाई रोक, Fable का हिस्सा (`You've reached your Fable limit.`), सर्वर का throttle, किसी API key की rate limit और किसी subagent का रुकना कभी session को नहीं ले जाते।
4. ccam अकाउंट का प्लान जानता है, उसकी `.claude.json` से या उससे जो `ccam login` ने दर्ज किया। नहीं तो notification `ccam does not know personal's plan, so auto_switch leaves this session.` पर ख़त्म होता है।
5. Session, Remote Control से जुड़ा नहीं है। दूसरे अकाउंट के रूप में वह उस app से फिर नहीं जुड़ेगा जिससे आप उसे चलाते हैं, इसलिए वह वहीं रहता है, और notification यह बताता है।
6. सीमा अब से 15 मिनट से ज़्यादा बाद reset होती है। इससे कम इंतज़ार के लिए ccam session को Claude Code पर छोड़ देता है: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. सीमा के बाद से आपने कोई prompt टाइप नहीं किया है।
8. किसी दूसरे अकाउंट में गुंजाइश है। ccam वैसे ही चुनता है जैसे [`--auto`](#auto) (`auto_max_usage`, `auto_strategy`, बातचीत का मॉडल), उन अकाउंट्स में से जो उसी तरह के हैं जिस तरह का अकाउंट आप छोड़ रहे हैं: उसी Team या Enterprise संगठन की कोई दूसरी सीट, या कोई दूसरा personal प्लान (Pro, Max)। `cross_org` और `--cross-org` यहाँ लागू नहीं होते। अकाउंट को:
   - ऐसे लॉगिन से लॉग इन होना चाहिए जिसे Claude Code renew कर सके;
   - उसी Claude अकाउंट में लॉग इन नहीं होना चाहिए;
   - API key इस्तेमाल नहीं करनी चाहिए;
   - ऐसा प्लान रखना चाहिए जो ccam को पता हो;
   - usage credits चालू नहीं रखने चाहिए, जहाँ तक ccam की अपनी उपयोग खोज दिखाती है;
   - Claude Code में इस फ़ोल्डर पर भरोसा किया होना चाहिए;
   - अगर बातचीत किसी agent के रूप में चलती है (`agent` सेटिंग या बातचीत किसी agent का नाम देती है), तो जिस अकाउंट को आप छोड़ रहे हैं उसके agents साझा करने चाहिए: उसके `agents` और `plugins` फ़ोल्डर और उसकी `settings.json` में से हर एक बिल्कुल वही होना चाहिए (किसी ग्रुप के ज़रिए), या दोनों में मौजूद न हो।

   जब ये आख़िरी तीन शर्तें गुंजाइश वाले हर अकाउंट को बाहर कर देती हैं, तो notification उनके नाम बताता है: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` या `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` जिस विंडो का [limit record](/features/usage#the-limit-record) है, वह भरी हुई गिनी जाती है, और जिस अकाउंट पर पिछले 2 मिनट में कोई दूसरा session ले जाया गया, वह बाक़ियों के बाद आता है। चुनाव लगभग 6 सेकंड के भीतर होना चाहिए; जब उपयोग की खोज में ज़्यादा समय लगता है, तो session वहीं रहता है। `usage_source = "claude"` के साथ, `auto_switch` केवल सेव की गई संख्याओं से तय करता है, और ccam Claude Code की खोज या स्टेटस लाइन से यह नहीं जान सकता कि usage credits चालू हैं या नहीं, इसलिए वह मान लेता है कि वे बंद हैं।

फिर hook Claude Code को बंद करता है, और उसके साथ-साथ चल रहा ccam उन टर्मिनलों में desktop notification दिखाता है जो इसे समर्थन करते हैं, जैसे `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` यह ऊपर के उदाहरण output की पहली पंक्ति प्रिंट करता है और चुने गए अकाउंट के रूप में उस session के लिए `ccam next --continue` चलाता है, जो उसकी बाक़ी दो पंक्तियाँ प्रिंट करता है और session के विकल्पों, permission mode, effort level और agent के साथ Claude Code शुरू करता है, जैसे [relay](#relay) करता है। नए session का पहला संदेश यह तय टेक्स्ट होता है, जो अकेली चीज़ है जिसे ccam आपकी ओर से भेजता है (`/ccam:switch` के बाद होने वाला स्विच कुछ नहीं भेजता):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

अगर वह अकाउंट भी अपनी सीमा पर पहुँचता है, तो वहाँ भी यही होता है। Limit records इस सिलसिले को किसी अकाउंट पर उसके reset से पहले वापस जाने से रोकते हैं, इसलिए जब हर अकाउंट भरा हो, तो session आख़िरी अकाउंट पर रहता है।

`agent` सेटिंग से, या agent का नाम बताने वाली उस पंक्ति से जिसे Claude Code समय-समय पर और session ख़त्म होने पर बातचीत में लिखता है, hook जानता है कि बातचीत किसी agent के रूप में चलती है। अगर अब तक केवल session का `--agent` उसका नाम देता है, और चुना गया अकाउंट agents साझा नहीं करता, तो `ccam next --continue` कुछ भी शुरू नहीं करता और नोट रख लेता है, और एक दूसरा desktop notification यह बताता है: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``।

जब किसी दूसरे अकाउंट में गुंजाइश नहीं होती, तो कुछ भी स्विच नहीं होता: Claude Code सीमा को हमेशा की तरह संभालता है, session में इंतज़ार करते हुए या अपने विकल्प देते हुए, और notification बताता है कि कौन सा अकाउंट सबसे पहले ख़ाली होगा, उदाहरण के लिए `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` सीटों के लिए यह `None of your other accounts in Acme has room now` कहता है। शुरू हो चुके स्विच को रोकने के लिए, Claude Code के बंद होते समय Ctrl+C दबाएँ: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

स्विच की क़ीमत होती है। Claude Code फिर से शुरू होता है, जिसमें कुछ सेकंड लगते हैं, और Claude के शुरू किए background shells, dev servers और agents पुराने session के साथ रुक जाते हैं। नया अकाउंट अपने पहले जवाब के साथ पूरी बातचीत फिर से पढ़ता है, जैसे [`--continue`](#continue) के साथ, और ccam बताता है कि वह कितनी बड़ी है; स्विच कभी [`--brief`](#brief) इस्तेमाल नहीं करता।

`auto_switch` Claude Code के सीमा बताने से पहले कभी स्विच नहीं करता, background में कभी आपका उपयोग नहीं देखता, और कभी कोई लॉगिन कॉपी, renew या आगे नहीं देता। Windows पर इसे अनदेखा किया जाता है, और `ccam doctor` यह बताता है; वहाँ, `relay = true` के साथ, `/ccam:switch` और `/exit` अब भी बातचीत को दूसरे अकाउंट पर ले जाते हैं।

**Anthropic की शर्तें।** Team और Enterprise सीटों पर, `auto_switch` आपका काम उसी संगठन की किसी दूसरी सीट पर ले जाता है। Anthropic का सहायता केंद्र कहता है कि "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team प्लान की हर सीट एक नामित व्यक्ति को सौंपी जाती है, और केवल वही व्यक्ति उसका इस्तेमाल कर सकता है) और "Users can't route their work through another user's account to use that user's remaining usage" (उपयोगकर्ता किसी दूसरे उपयोगकर्ता के बचे हुए उपयोग का इस्तेमाल करने के लिए अपना काम उसके अकाउंट के ज़रिए नहीं चला सकते) ([Team प्लान](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise प्लान](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) के लिए भी इसी तरह)।

Personal प्लान (Pro, Max) के लिए, Anthropic की शर्तें यह नहीं बतातीं कि आपके अपने अकाउंट सीमा पर अपने-आप एक-दूसरे की जगह ले सकते हैं या नहीं। [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) Services को "through automated or non-human means, whether through a bot, script, or otherwise" (स्वचालित या ग़ैर-मानवीय तरीक़ों से, चाहे bot, script या किसी और तरह) इस्तेमाल करने से मना करती हैं, जब तक Anthropic साफ़ तौर पर इसकी अनुमति न दे, और स्विच आपकी ओर से ऊपर वाला तय संदेश भेजता है। 12 नवंबर 2026 से लागू होने वाली [Usage Policy](https://www.anthropic.com/legal/aup) "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (हमारे products में बनाई गई क्षमताओं, पाबंदियों या सुरक्षा उपायों को जानबूझकर दरकिनार करना या उनमें दख़ल देना) को उन बातों में गिनाती है जिन पर वह रोक लगाती है।

ccam यह दावा नहीं करता कि `auto_switch` इन शर्तों का पालन करता है। इसे चालू करना आपका फ़ैसला है, और जोखिम भी आपका है: जब Anthropic को लगे कि शर्तें तोड़ी गई हैं, तो वह किसी अकाउंट को, या किसी Team या Enterprise संगठन की पहुँच को, निलंबित कर सकता है।

जब `auto_switch` चालू हो, तो `ccam doctor` ccam plugin के बिना वाले अकाउंट और 2.1.214 से पुराने Claude Code वर्ज़न के बारे में चेतावनी देता है, और सीटों के लिए Anthropic का सीट वाला नियम दोहराता है ([जाँचें](/features/doctor))।

## सीमा पर फिर से चलाना

`ccam next` उस interactive session के लिए है जो रुक गया है। एक non-interactive run, यानी किसी script या scheduled job में `claude -p`, आपका इंतज़ार नहीं कर सकता, इसलिए अकाउंट नाम से पहले `--fallback` लगाने से ccam उसके लिए स्विच करता है:

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

1. `--fallback` केवल `-p` या `--print` के साथ काम करता है; वरना ccam रुक जाता है और `ccam next --auto --continue` की ओर इशारा करता है।
2. Pipe से आया standard input एक बार, 10 MiB तक, पढ़ा जाता है और हर प्रयास को दिया जाता है। Claude Code की तरह, ccam उसके शुरू होने का ज़्यादा से ज़्यादा 3 सेकंड इंतज़ार करता है (`--input-format stream-json` के साथ छोड़कर); अगर pipe खुला रहता है लेकिन कुछ नहीं भेजता, तो run को कोई input नहीं मिलता और ccam standard error पर यह बताता है।
3. ccam `exec` के बजाय `claude` को एक child process के रूप में शुरू करता है, उसका standard output रोककर रखता है, और standard error को आते ही आगे भेजता है।
4. जब `claude` 0 के साथ बाहर निकलता है, या किसी दूसरे कारण से विफल होता है, तो रोका गया output और exit code बिना बदलाव आगे दिए जाते हैं। किसी signal से ख़त्म हुए run (128 से ऊपर का exit code) के साथ भी यही होता है: उसे रुकने को कहा गया था, और वह दोबारा नहीं चलता।
5. विफलता को उपयोग सीमा तब माना जाता है जब Claude Code का संदेश यही कहे, यानी उपयोग सीमा का कोई भी ऐसा संदेश जिसे [hook](#inside-claude-code) पहचानता है, `--output-format json` या `stream-json` output के अंदर भी, या जब अकाउंट का उपयोग उसकी 5 घंटे, साप्ताहिक या `--model` विंडो को 100% पर दिखाए (limit record समेत)।
6. तब उस प्रयास का output छोड़ दिया जाता है और वही कमांड किसी दूसरे लॉग इन अकाउंट के रूप में चलती है: सबसे कम 5 घंटे के उपयोग वाला, फिर वे अकाउंट जिनका उपयोग पता नहीं किया जा सकता, नंबर के क्रम में। 100% वाले अकाउंट छोड़ दिए जाते हैं, API key लॉगिन वाले अकाउंट कभी नहीं आज़माए जाते, और हर अकाउंट को ज़्यादा से ज़्यादा एक बार आज़माया जाता है।
7. आर्ग्युमेंट में `-c`, `--continue`, `-r` या `--resume` होने पर, केवल उसी इतिहास ग्रुप के अकाउंट आज़माए जाते हैं, ताकि बातचीत मिल सके।
8. `--auto` की तरह, केवल उसी संगठन के अकाउंट आज़माए जाते हैं (`cross_org = true` यह पाबंदी हटाता है)।

अगला अकाउंट `--auto` वाले नियमों से ही चुना जाता है (`auto_max_usage`, `auto_strategy`, `--model` की विंडो)। जब कोई अकाउंट नहीं बचता, तो ccam आख़िरी प्रयास का output छापता है, बताता है कि किसी लॉग इन अकाउंट में गुंजाइश नहीं बची, और code 75 के साथ बाहर निकलता है। अकाउंट नाम से पहले `--wait` के साथ (`ccam --fallback --wait=2h work -p …`), यह इसके बजाय पहले अकाउंट के reset होने का इंतज़ार करता है, फिर हर अकाउंट को दोबारा आज़मा सकता है; यह हमेशा पहले कम से कम एक बार रुकता है, ताकि अभी-अभी सीमा पर पहुँचा अकाउंट तुरंत दोबारा न चले। यह अवधि run के कुल इंतज़ार के समय को सीमित करती है, हर इंतज़ार को अलग-अलग नहीं। Space के साथ लिखने पर (`--wait 30m`), यह उसी संकेत के साथ अस्वीकार होता है जो `ccam next` देता है। Claude Code के आर्ग्युमेंट से पहले शुरुआत का `--` किसी भी launch की तरह काम करता है (`ccam --fallback work -- -p …`)। `ccam --dry-run --fallback work -p …` पहले प्रयास की कमांड छापता है।

Interactive sessions तब तक अपने-आप स्विच नहीं होते जब तक आप [`auto_switch`](#auto-switch) चालू न करें। `ccam statusline install` के साथ, 5 घंटे या साप्ताहिक उपयोग 90% तक पहुँचते ही स्टेटस लाइन `→ ccam next --auto --continue` जोड़ देती है ([उपयोग सीमाएँ](/features/usage))।

## सुरक्षा और विशेष स्थितियाँ

- बिना लॉगिन वाले अकाउंट कभी नहीं चुने जाते। कोई दूसरा लॉगिन किया हुआ अकाउंट न होने पर, सादा `ccam next` `no logged-in account other than <name>` के साथ रुकता है।
- `--continue` के साथ और फ़ोल्डर में कोई बातचीत न हो: `no account has a conversation in <folder> to continue`।
- जब संगठन वाला नियम कोई अकाउंट नहीं छोड़ता: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``।
- `--auto` के साथ और कोई लॉगिन किया हुआ अकाउंट न हो: `no logged-in account to pick`। जब हर अकाउंट भरा हो या उसका उपयोग अज्ञात हो: `` no account with known usage and headroom (see `ccam usage`) ``। दोनों exit code 75 के साथ ख़त्म होते हैं।
- `--max-usage`, `--strategy`, `--wait` और `--cross-org` केवल `--auto` के साथ काम करते हैं, और 1-100 से बाहर का मान या कोई अज्ञात strategy अस्वीकार की जाती है।
- Fork बातचीत को कॉपी करता है, अकाउंट की दूसरी state को नहीं: मूल अकाउंट के `file-history/` के checkpoints तब तक साथ नहीं आते जब तक इतिहास साझा न हो।
- `--auto`, `--dry-run` के साथ भी नेटवर्क requests भेजता है। यह केवल पढ़ता है, और token फ़ाइलों में कभी नहीं लिखता।
- चुना गया अकाउंट उसी रास्ते से शुरू होता है जिससे `ccam <name>`, इसलिए वही जाँचें लागू होती हैं, उस चेतावनी सहित जब वही इतिहास साझा करने वाला कोई दूसरा Claude Code session इस फ़ोल्डर में पहले से चल रहा हो।
- `last-used` केवल तभी लिखा जाता है जब Claude Code वास्तव में शुरू होता है। `--dry-run` इसे नहीं छूता, इसलिए dry run रोटेशन को आगे नहीं बढ़ाता।
- `--fallback` किसी प्रयास के ख़त्म होने तक standard output रोककर रखता है, इसलिए `--output-format stream-json` जैसा streaming output आख़िर में एक साथ आता है। जो client `stream-json` input stream करता है और आगे भेजने से पहले जवाबों का इंतज़ार करता है, वह `--fallback` इस्तेमाल नहीं कर सकता। सीमा तक पहुँचने से पहले पहले प्रयास ने जो फ़ाइलें बदलीं, वे बदली ही रहती हैं; अगला प्रयास prompt को फिर से शुरू करता है।

## संबंधित

- [उपयोग सीमाएँ](/features/usage)
- [उपयोग सीमाएँ और अकाउंट स्विच करना](/docs/usage-limits)
- [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing)
- [Claude Code की usage limit पूरी होने पर क्या करें](/guides/usage-limit-workflow)
