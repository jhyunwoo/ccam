---
title: उपयोग सीमाएँ और अकाउंट स्विच करना
description: ccam usage से हर Claude Code अकाउंट का 5 घंटे और साप्ताहिक उपयोग देखें, और ccam next --auto से सबसे ज़्यादा गुंजाइश वाला अकाउंट खोलें।
updated: 2026-11-05
order: 5
translation: machine
sourceHash: "06bd349073e3"
---

`ccam usage` हर अकाउंट के लिए दिखाता है कि हर सीमा का कितना हिस्सा इस्तेमाल हुआ है: 5 घंटे की सीमा, साप्ताहिक सीमा और हर मॉडल की अपनी साप्ताहिक सीमा (जैसे Fable)। `ccam next --auto` सबसे ज़्यादा गुंजाइश वाला अकाउंट खोलता है, और `--continue` मौजूदा बातचीत को साथ ले जाता है।

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

हर अकाउंट को एक पंक्ति मिलती है जो बताती है कि उसकी संख्याएँ कितनी पुरानी हैं, फिर हर सीमा के लिए एक पट्टी, रीसेट होने तक के समय के साथ: 5 घंटे की विंडो, सभी मॉडल की साप्ताहिक विंडो, और हर मॉडल की अपनी साप्ताहिक विंडो, जैसे Fable। पट्टी 80% पर पीली और 100% पर लाल हो जाती है। आख़िरी पंक्ति उस अकाउंट का नाम बताती है जिसे `ccam next --auto` खोलता; `auto_strategy = "expiring"` या `"balance"` के साथ उसकी जगह `→ next --auto picks <name>: <reason>` लिखा होता है।

**Team और Enterprise सीटें।** अपने-आप अकाउंट चुनना (`ccam next`, `--fallback`, `ccam pick`, डैशबोर्ड का `n` और, चालू करने के बाद, [`auto_switch`](/features/next#auto-switch)) एक ही Claude Team या Enterprise संगठन की सीटों के बीच भी स्विच करता है, जो Anthropic की शर्तों का उल्लंघन हो सकता है। इस पर भरोसा करने से पहले [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) पढ़ें और अपने संगठन और Anthropic से पता कर लें।

## संख्याएँ कहाँ से आती हैं

ccam हर अकाउंट के अपने लॉगिन से, सभी अकाउंट्स के लिए एक साथ Anthropic के usage endpoint से पूछता है। यह हर अकाउंट के लिए अधिकतम 3 सेकंड इंतज़ार करता है और हर अकाउंट को जवाब मिलते ही दिखा देता है।

यह endpoint सार्वजनिक रूप से documented नहीं है और बदल सकता है। जब कोई request विफल होती है, तो ccam आख़िरी ज्ञात संख्याएँ दिखाता है और उन्हें cached के रूप में चिह्नित करता है। जो विंडो उन संख्याओं के सहेजे जाने के बाद रीसेट हो चुकी है, वह 0% दिखाती है। अगर response का format बदल गया हो, तो `ccam doctor` आपको बताता है।

जब Claude Code किसी session को उपयोग सीमा पर रोकता है, तो [ccam plugin](/docs/install#claude-code-plugin) दर्ज करता है कि कौन-सी विंडो भरी थी और वह कब रीसेट होती है। तब तक वह विंडो 100% दिखाती है, बाक़ी संख्याएँ चाहे जो कहें, ताकि `ccam next --auto` और `--fallback` उस अकाउंट को बहुत जल्दी फिर से न चुनें।

60 सेकंड से कम पहले लाई गई संख्याएँ दोबारा इस्तेमाल होती हैं। विकल्प:

| विकल्प | असर |
| --- | --- |
| `--watch[=seconds]` | लगातार रीफ़्रेश करें, डिफ़ॉल्ट रूप से हर 60 सेकंड |
| `--json` | मशीन से पढ़ने लायक आउटपुट |
| `-v` | यह भी दिखाएँ कि हर अकाउंट की संख्याएँ कहाँ से आईं |
| `--no-cache` | हमेशा फिर से पूछें |
| `--source cache` | कभी न पूछें; केवल सेव की गई संख्याएँ दिखाएँ |

## एक्सपायर हुए लॉगिन

जब किसी अकाउंट का access token एक्सपायर हो जाता है, तो ccam उसे renew नहीं करता, क्योंकि अगर कोई Claude Code session वही लॉगिन इस्तेमाल कर रहा हो, तो renew करने से वह session साइन आउट हो सकता है। अकाउंट पर ``access token expired; `ccam <name>` renews it when it starts`` दिखता है: उसे एक बार `ccam <name>` से शुरू करें, Claude Code उसे renew कर देगा। जो लॉगिन पूरी तरह एक्सपायर हो गया हो, या जिसे Claude Code ने हटा दिया हो, उस पर `` login expired → `ccam login <name>` `` दिखता है और उसे नया लॉगिन चाहिए।

`ccam usage --refresh` एक्सपायर हुए tokens ख़ुद renew करता है, लेकिन केवल उन अकाउंट्स के लिए जिनका कोई Claude Code session नहीं चल रहा, चाहे वह session किसी भी तरह शुरू हुआ हो। यह macOS कीचेन में या import किए गए `~/.claude` में रखे लॉगिन कभी renew नहीं करता, और एक renewal में 30 सेकंड तक लग सकते हैं।

## स्टेटस लाइन में उपयोग

```sh
ccam statusline install
```

इसके बाद Claude Code स्टेटस लाइन में अकाउंट का नाम और उसका उपयोग दिखाता है, उदाहरण के लिए `[work] · Opus · 5h 24% · 7d 41%`। 90% से आगे, लाइन `ccam next --auto --continue` भी सुझाती है। अगर आपके पास पहले से कोई status line कमांड है, तो ccam उसे रखता है और उसके आगे अकाउंट का नाम जोड़ता है। स्टेटस लाइन अपनी संख्याएँ सेव भी करती है, जिन्हें `ccam usage` तब इस्तेमाल करता है जब endpoint जवाब नहीं देता।

## स्विच करना

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the most room
ccam next --continue  # continue this folder's conversation on the next account
```

`ccam next` उन अकाउंट्स को छोड़ देता है जिनमें लॉगिन नहीं है, और `--auto` उन अकाउंट्स को छोड़ देता है जो 100% पर हैं (वह विंडो भी जिसे Claude Code ने भरा हुआ बताया, उसके reset तक) और जो किसी दूसरे संगठन के हैं। जब दोनों अकाउंट्स बातचीत का इतिहास साझा करते हैं, तो `--continue` वही बातचीत खोलता है; नहीं तो Claude Code उसे नए अकाउंट में कॉपी कर देता है। [साझा करना](/docs/sharing) और [अगले अकाउंट पर स्विच करना](/features/next#continue) देखें।

`config.toml` में सबसे ऊपर `auto_switch = true` होने पर (Linux और macOS; डिफ़ॉल्ट रूप से बंद), जब Claude Code किसी interactive session को उसकी सीमा पर रोकता है और उसी तरह के किसी दूसरे अकाउंट (उसी संगठन की कोई सीट, या कोई दूसरा personal प्लान) में गुंजाइश होती है, तो ccam यह स्विच ख़ुद करता है। यह कब काम करता है और क्या कभी नहीं करता, इसके लिए [`auto_switch`](/features/next#auto-switch) देखें।

`claude -p` चलाने वाली scripts और scheduled jobs के लिए, अकाउंट नाम से पहले `--fallback` लगाएँ। जब वह अकाउंट अपनी सीमा तक पहुँचता है, तो ccam वही कमांड सबसे ज़्यादा गुंजाइश वाले अकाउंट के रूप में फिर से चलाता है, और सिर्फ़ उस run का output छपता है जो पूरा हुआ:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

ccam यह कैसे तय करता है, इसके लिए [सीमा पर फिर से चलाना](/features/next#run-again-at-the-limit) देखें।

[उपयोग सीमा वर्कफ़्लो गाइड](/guides/usage-limit-workflow) एक पूरा उदाहरण क़दम-दर-क़दम दिखाती है।
