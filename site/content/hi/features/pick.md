---
title: scripts के लिए अकाउंट चुनना
description: ccam pick script के लिए सबसे ज़्यादा गुंजाइश वाला अकाउंट प्रिंट करता है, साथ शुरू हुए workers को अलग अकाउंट देता है, और गुंजाइश न होने पर 75 देता है।
updated: 2026-11-05
order: 12
translation: machine
sourceHash: "0ec25752ace9"
---

`ccam pick` सबसे ज़्यादा गुंजाइश वाले लॉगिन किए हुए अकाउंट का नाम प्रिंट करता है, ताकि कोई script उससे Claude Code शुरू कर सके। जब तक काफ़ी अकाउंट्स में गुंजाइश हो, एक ही समय पर शुरू हुए हर worker को अलग अकाउंट मिलता है, इसलिए `claude -p` jobs का एक batch किसी एक अकाउंट को ख़ाली करने के बजाय आपके अकाउंट्स में बँट जाता है।

**Team और Enterprise सीटें।** अपने-आप अकाउंट चुनना (`ccam next`, `--fallback`, `ccam pick`, डैशबोर्ड का `n` और, चालू करने के बाद, [`auto_switch`](/features/next#auto-switch)) एक ही Claude Team या Enterprise संगठन की सीटों के बीच भी स्विच करता है, जो Anthropic की शर्तों का उल्लंघन हो सकता है। इस पर भरोसा करने से पहले [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) पढ़ें और अपने संगठन और Anthropic से पता कर लें।

## उपयोग

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # most weekly quota left per day; fewer sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

Exit codes: standard output पर नाम के साथ `0`, जब किसी अकाउंट में गुंजाइश न हो या कोई लॉगिन न हो तो `75` (`EX_TEMPFAIL`), और ग़लत flag जैसी दूसरी त्रुटियों के लिए `1`।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--json` | अकाउंट, उसका 5 घंटे और साप्ताहिक उपयोग, संख्याएँ कहाँ से आईं और उसे क्यों चुना गया, यह प्रिंट करें | बंद |
| `--model <model>` | उस मॉडल की साप्ताहिक विंडो भी गिनें, जैसे `opus` | कोई नहीं |
| `--max-usage <1-100>` | इस उपयोग से नीचे वाले अकाउंट्स को प्राथमिकता दें | `auto_max_usage`, वरना 100 |
| `--strategy headroom\|expiring\|balance` | सबसे कम 5 घंटे के उपयोग से 3 प्रतिशत अंक के भीतर वाले अकाउंट्स में से सबसे कम साप्ताहिक उपयोग, सबसे जल्दी होने वाला 5 घंटे का reset या प्रति दिन सबसे ज़्यादा बचा साप्ताहिक कोटा | `auto_strategy`, वरना `headroom` |
| `--cross-org` | दूसरे संगठनों के अकाउंट भी चुनें | `cross_org`, वरना बंद |
| `--wait[=<duration>]` | जब किसी अकाउंट में गुंजाइश न हो, तो पहले अकाउंट के reset होने का इंतज़ार करें, ज़्यादा से ज़्यादा इतनी देर। मान `=` के बाद आता है: `--wait 30m` इस संकेत के साथ अस्वीकार होता है: `--wait takes its value after '=': --wait=30m` | बंद; अकेले देने पर `5h` |
| `--dry-run` | चुनाव दर्ज किए बिना चुनें | बंद |

## यह कैसे काम करता है

1. ccam उस अकाउंट के संगठन के लॉगिन किए हुए अकाउंट्स की सूची बनाता है जिसे यह फ़ोल्डर शुरू करता (देखें [संगठन](/features/next#auto)); API key अकाउंट्स और Console लॉगिन की कोई सीमा नहीं होती, इसलिए वे छोड़ दिए जाते हैं, और जो फ़ोल्डर उसी Claude अकाउंट में लॉगिन है जिसमें वह अकाउंट, वह भी छोड़ दिया जाता है।
2. यह एक lock (`~/.local/state/ccam/picks.lock`) लेता है, ताकि एक साथ चल रहे picks एक के बाद एक हों।
3. यह हर अकाउंट का उपयोग एक साथ पूछता है, उसी तरह जैसे `ccam usage` करता है, 60 सेकंड के cache के साथ। जिस window को Claude Code ने भरा हुआ बताया, वह अपने reset तक 100% गिनी जाती है ([limit record](/features/usage#the-limit-record)), इसलिए जो अकाउंट अभी अपनी सीमा पर पहुँचा है वह बाहर रखा जाता है, भले ही उसके सहेजे गए आँकड़े अभी भी गुंजाइश दिखाएँ।
4. पिछले 2 मिनट में चुने गए अकाउंट, `ccam pick` द्वारा या [`auto_switch`](/features/next#auto-switch) द्वारा वहाँ कोई session ले जाने से (`~/.local/state/ccam/picks/<name>`; भविष्य के timestamp वाला pick नहीं गिना जाता), पहले छोड़ दिए जाते हैं, और बाक़ी में से `ccam next --auto` के [चुनाव के नियम](/features/next#auto) चुनते हैं: पहले `auto_max_usage` से नीचे वाले, strategy के हिसाब से, और 100% वाला अकाउंट कभी नहीं। अगर इससे वह अकाउंट छूट जाता है जिसे नियम वरना चुनते, तो `--json` का `reason` यह बताता है: `30% of its 5 hours used; work was picked a moment ago`।
5. अगर गुंजाइश वाला हर अकाउंट हाल ही में चुना गया हो, तो नियम फिर से उन सब में से चुनते हैं, इसलिए तीन अकाउंट्स होने पर चौथा worker उनमें से किसी एक को साझा करता है।
6. चुनाव दर्ज होता है और lock छोड़ा जाता है, फिर नाम प्रिंट होता है।

2 मिनट की विंडो इसलिए ज़रूरी है क्योंकि नए job का उपयोग संख्याओं में दिखने में कुछ समय लगता है: इसके बिना, एक साथ शुरू हुए पाँच workers सभी एक ही अकाउंट को सबसे ज़्यादा गुंजाइश वाला मानते।

## सुरक्षा और विशेष स्थितियाँ

- `ccam pick` केवल उपयोग पढ़ता है और चुनाव का समय लिखता है; यह `ccam usage` जितना करता है उससे ज़्यादा tokens कभी renew नहीं करता और न पढ़ता है।
- चुनाव अकाउंट को आरक्षित नहीं करता। 2 मिनट बाद दो scripts फिर भी एक ही अकाउंट पर शुरू हो सकती हैं; जब वह अकाउंट अपनी सीमा तक पहुँचे तो सँभलने के लिए `--fallback` इस्तेमाल करें।
- `--wait` से इंतज़ार करते समय कोई lock नहीं रखा जाता, इसलिए इस बीच दूसरे picks चलते रहते हैं।
- `ccam pick` [`ccam folder`](/features/folder) को नज़रअंदाज़ करता है: यह केवल उपयोग के हिसाब से चुनता है।

## संबंधित

- [अगले अकाउंट पर स्विच करना](/features/next)
- [सीमा पर फिर से चलाना](/features/next#run-again-at-the-limit)
- [उपयोग सीमाएँ](/features/usage)
- [API key अकाउंट](/features/api)
