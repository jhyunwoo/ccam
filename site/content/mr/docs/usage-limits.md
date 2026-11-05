---
title: वापर मर्यादा आणि खाते बदलणे
description: ccam usage ने प्रत्येक Claude Code खात्याचा 5 तासांचा आणि साप्ताहिक वापर पाहा, आणि ccam next --auto ने सर्वात जास्त जागा असलेले खाते उघडा.
updated: 2026-11-05
order: 5
translation: machine
sourceHash: "06bd349073e3"
---

`ccam usage` प्रत्येक खात्यासाठी प्रत्येक मर्यादेपैकी किती वापरले ते दाखवते: 5 तासांची मर्यादा, साप्ताहिक मर्यादा आणि प्रत्येक मॉडेलची स्वतःची साप्ताहिक मर्यादा (उदा. Fable). `ccam next --auto` सर्वात जास्त जागा असलेले खाते उघडते, आणि `--continue` चालू संभाषण तिकडे नेते.

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

प्रत्येक खात्याला एक ओळ मिळते, ज्यात त्याचे आकडे किती जुने आहेत ते असते, आणि मग प्रत्येक मर्यादेसाठी एक bar, ती रीसेट होईपर्यंतच्या वेळेसह: 5 तासांची window, सर्व मॉडेल्ससाठीची साप्ताहिक window, आणि प्रत्येक मॉडेलची स्वतःची साप्ताहिक window, उदा. Fable. bar 80% वर पिवळा आणि 100% वर लाल होतो. शेवटची ओळ `ccam next --auto` जे खाते उघडेल त्याचे नाव देते; `auto_strategy = "expiring"` किंवा `"balance"` असताना त्याऐवजी ती `→ next --auto picks <name>: <reason>` अशी असते.

**Team आणि Enterprise सीट्स.** आपोआप होणारी खाते निवड (`ccam next`, `--fallback`, `ccam pick`, डॅशबोर्डमधील `n` आणि, तुम्ही चालू केल्यावर, [`auto_switch`](/features/next#auto-switch)) एकाच Claude Team किंवा Enterprise organization मधील सीट्सदरम्यानही खाते बदलते, आणि यामुळे Anthropic च्या अटींचे उल्लंघन होऊ शकते. त्यावर अवलंबून राहण्यापूर्वी [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) वाचा आणि तुमच्या संस्थेकडे आणि Anthropic कडे चौकशी करा.

## आकडे कुठून येतात

ccam प्रत्येक खात्याचे स्वतःचे लॉगिन वापरून, सर्व खात्यांसाठी एकाच वेळी Anthropic च्या usage endpoint ला विचारते. ते प्रत्येक खात्यासाठी जास्तीत जास्त 3 सेकंद थांबते आणि प्रत्येक खात्याचे उत्तर येताच ते दाखवते.

हा endpoint सार्वजनिकरित्या दस्तऐवजीकृत नाही आणि बदलू शकतो. एखादी request अयशस्वी झाली तर ccam शेवटचे माहीत असलेले आकडे दाखवते आणि त्यांना cached म्हणून चिन्हांकित करते. ते आकडे साठवल्यानंतर जी window रीसेट झाली आहे ती 0% दाखवते. response format बदलला असेल तर `ccam doctor` तुम्हाला सांगते.

Claude Code एखादे सत्र वापर मर्यादेवर थांबवते, तेव्हा [ccam plugin](/docs/install#claude-code-plugin) कोणती window भरली होती आणि ती कधी रीसेट होते याची नोंद ठेवते. तोपर्यंत इतर आकडे काहीही सांगोत, ती window 100% दाखवते, त्यामुळे `ccam next --auto` आणि `--fallback` ते खाते फार लवकर पुन्हा निवडत नाहीत.

60 सेकंदांपेक्षा कमी वेळापूर्वी आणलेले आकडे पुन्हा वापरले जातात. पर्याय:

| पर्याय | परिणाम |
| --- | --- |
| `--watch[=seconds]` | सतत रिफ्रेश करा, डीफॉल्टनुसार दर 60 सेकंदांनी |
| `--json` | मशीनला वाचता येईल असा आउटपुट |
| `-v` | प्रत्येक खात्याचे आकडे कुठून आले तेही दाखवा |
| `--no-cache` | नेहमी पुन्हा विचारा |
| `--source cache` | कधीच विचारू नका; फक्त साठवलेले आकडे दाखवा |

## कालबाह्य झालेले लॉगिन

खात्याचे access token कालबाह्य झाले असेल तर ccam ते नूतनीकरण करत नाही, कारण Claude Code चे session तेच लॉगिन वापरत असताना ते नूतनीकरण केल्यास ते session साइन आउट होऊ शकते. खाते ``access token expired; `ccam <name>` renews it when it starts`` दाखवते: `ccam <name>` ने ते एकदा सुरू करा, म्हणजे Claude Code ते नूतनीकरण करते. संपूर्ण लॉगिनच कालबाह्य झाले असेल, किंवा Claude Code ने ते काढून टाकले असेल, तर `` login expired → `ccam login <name>` `` दिसते आणि नवीन लॉगिन लागते.

`ccam usage --refresh` कालबाह्य tokens स्वतः नूतनीकरण करते, पण फक्त ज्या खात्यांचे कोणतेही Claude Code session चालू नाही त्यांच्यासाठी, ते session कसेही सुरू केलेले असो. macOS Keychain मध्ये किंवा इम्पोर्ट केलेल्या `~/.claude` मध्ये साठवलेले लॉगिन ते कधीच नूतनीकरण करत नाही, आणि एका नूतनीकरणाला 30 सेकंदांपर्यंत लागू शकतात.

## स्टेटस लाइनमध्ये वापर

```sh
ccam statusline install
```

त्यानंतर Claude Code स्टेटस लाइनमध्ये खात्याचे नाव आणि त्याचा वापर दाखवते, उदाहरणार्थ `[work] · Opus · 5h 24% · 7d 41%`. 90% पासून पुढे ही ओळ `ccam next --auto --continue` सुद्धा सुचवते. तुमच्याकडे आधीच स्टेटस लाइन कमांड असेल, तर ccam ती ठेवते आणि तिच्या आधी खात्याचे नाव लावते. स्टेटस लाइन आपले आकडेही साठवते, जे endpoint उत्तर देत नाही तेव्हा `ccam usage` वापरते.

## खाते बदलणे

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the most room
ccam next --continue  # continue this folder's conversation on the next account
```

`ccam next` लॉगिन नसलेली खाती वगळते, आणि `--auto` 100% वर असलेली खाती (Claude Code ने भरलेली म्हणून नोंदवलेली window सुद्धा, तिच्या reset पर्यंत) आणि दुसऱ्या organization ची खाती वगळते. दोन्ही खाती संभाषण इतिहास शेअर करत असतील तर `--continue` तेच संभाषण उघडते; नाहीतर Claude Code ते नवीन खात्यात कॉपी करते. [शेअरिंग](/docs/sharing) आणि [पुढच्या खात्यावर जाणे](/features/next#continue) पाहा.

`config.toml` च्या सुरुवातीला `auto_switch = true` असेल (Linux आणि macOS; डीफॉल्टनुसार बंद), तर Claude Code interactive सत्र त्याच्या मर्यादेवर थांबवते आणि त्याच प्रकारच्या दुसऱ्या खात्यात (त्याच organization मधील सीट, किंवा दुसरा वैयक्तिक प्लॅन) जागा असते, तेव्हा ccam हा बदल स्वतः करते. ते केव्हा काम करते आणि काय कधीच करत नाही, ते [`auto_switch`](/features/next#auto-switch) मध्ये पाहा.

`claude -p` चालवणाऱ्या scripts आणि scheduled jobs साठी, खात्याच्या नावाआधी `--fallback` लावा. ते खाते मर्यादेला पोहोचले की ccam तीच कमांड सर्वात जास्त जागा असलेल्या खात्याने पुन्हा चालवते, आणि फक्त पूर्ण झालेल्या run चा output छापला जातो:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

ccam कसे ठरवते ते [मर्यादेवर पुन्हा चालवणे](/features/next#run-again-at-the-limit) मध्ये पाहा.

[वापर मर्यादा workflow मार्गदर्शक](/guides/usage-limit-workflow) एक संपूर्ण उदाहरण टप्प्याटप्प्याने दाखवते.
