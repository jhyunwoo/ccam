---
title: दूसरे अकाउंट्स पर subagents
description: ccam agent, Claude Code के task का एक हिस्सा आपके दूसरे अकाउंट के रूप में चलाकर Claude का जवाब लौटाता है, ताकि एक लंबा task कई अकाउंट्स का उपयोग करे।
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Claude Code के अपने subagents उस session के अकाउंट के रूप में चलते हैं जो उन्हें शुरू करता है, इसलिए वे उसी अकाउंट का उपयोग ख़र्च करते हैं। `ccam agent "<task>"` उसी फ़ोल्डर में `claude -p` से एक task आपके किसी दूसरे अकाउंट के रूप में चलाता है और Claude का आख़िरी जवाब प्रिंट करता है। Claude Code के किसी session से बुलाए जाने पर यह ऐसे subagent की तरह काम करता है जो किसी दूसरे अकाउंट का उपयोग ख़र्च करता है, इसलिए एक लंबा task उस बिंदु के बाद भी चलता रह सकता है जहाँ session का अपना अकाउंट ख़त्म हो चुका होता।

**Team और Enterprise सीटें।** `ccam agent` अकाउंट ख़ुद चुनता है। एक ही Claude Team या Enterprise संगठन के भीतर यह task किसी दूसरी सीट पर चलाता है, जो Anthropic की शर्तों का उल्लंघन हो सकता है। इस पर भरोसा करने से पहले [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) पढ़ें और अपने संगठन और Anthropic से पता कर लें।

## उपयोग

Claude Code इसे अपने Bash tool से चलाता है। [ccam plugin](/docs/install#claude-code-plugin) इंस्टॉल होने पर इसे सीधे शब्दों में माँगें, जैसे "इस task के हिस्सों के लिए मेरे दूसरे अकाउंट्स इस्तेमाल करो", या यह बात अपनी CLAUDE.md में लिख दें। plugin का skill Claude Code से कहता है कि हर agent को background में शुरू करे, ऐसा task लिखे जो अपने-आप में पूरा हो, और फ़ाइलें बदलने वाले agents को उनका अपना हिस्सा या अपना git worktree दे।

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

Claude का आख़िरी जवाब standard output पर जाता है; `ccam:` से शुरू होने वाली पंक्तियाँ standard error पर जाती हैं। `--json` इसके बजाय `{account, session_id, result, is_error, accounts}` प्रिंट करता है, जहाँ `accounts` उन सभी अकाउंट्स की सूची है जिनके रूप में task चला।

Exit codes: claude के अपने (`0` जब agent पूरा हो गया), `75` (`EX_TEMPFAIL`) जब बुलाने वाले के अलावा किसी अकाउंट में गुंजाइश न हो, और दूसरी त्रुटियों के लिए `1`, जैसे कोई agent किसी दूसरे agent को शुरू करने की कोशिश करे।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--model <model>` | agent का मॉडल; चुनाव में उस मॉडल की साप्ताहिक विंडो भी गिनी जाती है | अकाउंट का डिफ़ॉल्ट मॉडल |
| `--resume <id>` | किसी पिछले agent की बातचीत को एक अगले task के साथ जारी रखें (session ID उसकी पहली और आख़िरी पंक्ति पर होता है) | बंद |
| `--json` | अकाउंट, session, जवाब और वे अकाउंट्स प्रिंट करें जिनके रूप में यह चला | बंद |
| `--max-usage <1-100>` | इस उपयोग से नीचे वाले अकाउंट्स को प्राथमिकता दें | `auto_max_usage`, वरना 100 |
| `--strategy headroom\|expiring\|balance` | [चुनाव का नियम](/features/next#auto) | `auto_strategy`, वरना `headroom` |
| `--cross-org` | दूसरे संगठनों के अकाउंट्स के रूप में भी चलाएँ | `cross_org`, वरना बंद |
| `--wait[=<duration>]` | जब किसी अकाउंट में गुंजाइश न हो, तो पहले अकाउंट के reset होने का इंतज़ार करें, ज़्यादा से ज़्यादा इतनी देर | बंद; अकेले देने पर `5h` |
| `--dry-run` | जो अकाउंट चुना जाता, उसके लिए कमांड प्रिंट करें, उसे चलाए बिना | बंद |

## यह कैसे काम करता है

1. **बुलाने वाला अकाउंट छोड़ दिया जाता है।** जिस अकाउंट से `ccam agent` बुलाया जाता है (ccam के शुरू किए session का अकाउंट, वरना वह अकाउंट जिसे यह फ़ोल्डर शुरू करता है), वह agent कभी नहीं चलाता, और न ही उसी Claude अकाउंट में लॉगिन कोई दूसरा फ़ोल्डर, जो उसकी सीमाएँ साझा करता है। जब तक `--cross-org` कुछ और न कहे, agent बुलाने वाले के [संगठन](/features/next#auto) में ही रहता है।
2. **यह `ccam pick` की तरह चुनता है।** उपयोग उसी तरह पूछा जाता है जैसे [`ccam pick`](/features/pick) पूछता है, और एक-दूसरे के 2 मिनट के भीतर शुरू हुए agents को अलग-अलग अकाउंट मिलते हैं, जब तक काफ़ी अकाउंट्स में गुंजाइश हो।
3. **यह session की permissions के साथ चलता है, उससे ज़्यादा कभी नहीं।** `claude -p` में किसी tool की अनुमति देने के लिए कोई मौजूद नहीं होता, इसलिए agent को उसे बुलाने वाले session का permission mode और effort मिलता है। terminal से चलाने पर Claude Code का डिफ़ॉल्ट mode और आपकी settings के permission नियम लागू होते हैं। `ccam agent` में ऐसा कोई विकल्प नहीं है जो इससे ज़्यादा अनुमति दे।
4. **उपयोग सीमा पर बातचीत आगे बढ़ती है।** जब agent का अकाउंट अपनी सीमा पर पहुँचता है, तो ccam उस सीमा को वैसे ही दर्ज करता है जैसे [plugin hook](/features/usage#the-limit-record) करता, और उसी बातचीत को सबसे ज़्यादा गुंजाइश वाले अकाउंट के रूप में जारी रखता है: Claude Code उसे वहाँ कॉपी करता है (`--resume <file> --fork-session`), इस पहले संदेश के साथ: "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." हर अकाउंट एक बार आज़माया जाता है।
5. **हर agent बातचीत दर्ज होती है**, `~/.local/state/ccam/agent/<session id>.json` में (30 दिन)। `--resume <id>` उसे उसके फ़ोल्डर में जारी रखता है: उस अकाउंट के रूप में जिसने उसे चलाया था, अगर उसमें अब भी गुंजाइश हो, वरना किसी दूसरे के रूप में। [`ccam next --continue`](/features/next) agent बातचीतों को छोड़ देता है, इसलिए वह आपकी बातचीत को ही जारी रखता रहता है।

## सुरक्षा और विशेष स्थितियाँ

- कोई agent अपने agents शुरू नहीं कर सकता: ccam उसे `CCAM_AGENT=1` से चिह्नित करता है, और वहाँ `ccam agent` चलने से मना कर देता है।
- `ccam agent` को रोकने पर (Ctrl+C, कोई signal, या Claude Code का कमांड को बीच में रोकना) agent भी रुक जाता है।
- agent उस बातचीत को नहीं देखता जिसने उसे बुलाया, केवल अपना task देखता है।
- एक ही समय पर एक ही फ़ाइलें बदलने वाले agents एक-दूसरे का काम मिटा देते हैं। हर एक को task का अपना हिस्सा, या अपना git worktree दें।
- जब बुलाने वाला session किसी दूसरे अकाउंट पर चला जाता है ([`auto_switch`](/features/next#auto-switch)), तो उसके शुरू किए agents उसके साथ रुक जाते हैं। उनकी बातचीतें बनी रहती हैं: `--resume <id>` उन्हें जारी रखता है।
- `agent` एक कमांड का नाम है, इसलिए किसी नए अकाउंट का नाम `agent` नहीं रखा जा सकता; पहले से मौजूद ऐसा अकाउंट अब भी `ccam run agent` से शुरू होता है।

## संबंधित

- [scripts के लिए अकाउंट चुनना](/features/pick)
- [अगले अकाउंट पर स्विच करना](/features/next)
- [उपयोग सीमाएँ](/features/usage)
- [हर अकाउंट के background sessions](/features/agents)
