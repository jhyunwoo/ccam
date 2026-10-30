---
title: डॅशबोर्ड
description: टर्मिनलवर ccam खाती, गट, फोल्डर नियम, extensions आणि तपासण्यांचा पूर्ण स्क्रीन डॅशबोर्ड उघडते; प्रत्येक बदल आधी आपली योजना दाखवतो.
updated: 2026-10-30
order: 15
translation: machine
sourceHash: "d65f1b5bb0c8"
---

standard input आणि output टर्मिनल असताना आर्ग्युमेंट्सशिवाय `ccam` पूर्ण स्क्रीन डॅशबोर्ड उघडते; इतर वेळी ते `ccam list` छापते. डॅशबोर्ड प्रत्येक खाते त्याच्या वापर आणि गटांसह दाखवतो, आणि कमांड्स जे करतात ते करतो: सुरू करणे, जोडणे, लॉग इन करणे, गट बदलणे, फोल्डर नियम सेट करणे, background सत्रांची यादी करणे, तपासणे आणि दुरुस्त करणे. खाते सुरू केल्यावर, किंवा background सत्राशी attach केल्यावर, तो बंद होतो; खाते जोडताना किंवा लॉग इन करताना तो थोड्या वेळासाठी टर्मिनल सोपवतो, आणि नंतर डॅशबोर्डवर परत येतो.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · r to refresh · a to add · ? for shortcuts
```

tab नावांखालची ओळ सध्याचा tab कशासाठी आहे ते एका वाक्यात सांगते. जास्त रुंद टर्मिनलवर Accounts tab प्रत्येक खात्याचा ईमेलही दाखवतो. प्रत्येक tab वर, tab नावांच्या ओळीच्या शेवटी [`auto_switch`](/features/next#auto-switch) चालू आहे की नाही ते दिसते (Linux आणि macOS; Windows वर त्याकडे दुर्लक्ष केले जाते). ते चालू असताना, त्याला अजून काय हवे आहे त्याची यादी Health tab दाखवतो: Claude Code 2.1.214 किंवा नंतरचे, आणि Claude subscription ने साइन इन करणाऱ्या प्रत्येक खात्यात चालू असलेले ccam plugin.

प्रत्येक वेळी उघडताना डॅशबोर्ड पार्श्वभूमीत GitHub ला विचारतो की नवीन ccam रिलीज आहे का; पहिली frame त्यासाठी कधीही थांबत नाही, आणि offline असताना काहीही दाखवले जात नाही. नवीन रिलीज असल्यास, शीर्षकाखालची एक ओळ तो सुचवते:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` आधी विचारतो, मग [`ccam update`](/features/update) प्रमाणेच अपडेट करतो: रिलीज त्याच्या SHA-256 checksums शी तपासला जातो आणि एकाच पायरीत ccam ची जागा घेतो. तुम्ही ज्या डॅशबोर्डमध्ये आहात तो अजूनही जुनी आवृत्ती आहे, त्यामुळे बाहेर पडा आणि `ccam` पुन्हा सुरू करा. अपडेट इथे करता येत नसेल, उदाहरणार्थ `.deb` पॅकेजमधून इन्स्टॉल केलेल्या ccam च्या प्रतीसाठी, तर डॅशबोर्ड `ccam update` सारखाच सल्ला दाखवतो आणि ऑफर कायम ठेवतो.

पहिली frame फक्त `config.toml` वरून, 100 ms पेक्षा खूप कमी वेळात काढली जाते; प्रत्येक खाते उत्तर देईल तसा वापर भरला जातो. प्रत्येक मर्यादेचा स्वतःचा bar असतो: 5 तासांची window, सर्व मॉडेल्ससाठीची साप्ताहिक window, आणि Fable ची स्वतःची साप्ताहिक window; `-` म्हणजे खात्याकडे अशी window नाही. टक्केवारी 80% वर पिवळी आणि 100% वर लाल होते. शोध अयशस्वी झाल्यास, cache मधील मूल्ये फिकट दाखवली जातात. Claude Code ने भरलेली म्हणून नोंदवलेली window तिच्या reset पर्यंत 100% दाखवते. कालबाह्य access token फिकट `⚠ token expired` दाखवते (पुढच्या सुरुवातीला Claude Code त्याचे नूतनीकरण करते), तर `⚠ login expired` म्हणजे खात्याला `ccam login` लागते.

## Tabs

| Tab | काय दाखवते | की |
| --- | --- | --- |
| Accounts | क्रमांक, नाव, ईमेल, 5 तासांचा, साप्ताहिक आणि Fable bars, शेवटचा वापर, गट; `✔` डीफॉल्ट, `▸` या फोल्डरचे खाते, `•` चालू (ccam शिवाय सुरू केलेले Claude Code सुद्धा: साधे `claude`, IDE, background सत्रे आणि त्यांचा daemon) | खाली |
| Agents | [`ccam agents`](/features/agents) जे यादीत दाखवते ते: प्रत्येक खात्याची background सत्रे आणि उघडी सत्रे, त्यांची स्थिती, फोल्डर आणि सुरू होण्याची वेळ यांसह (खाली) | `enter` attach करा (`ccam attach` प्रमाणे), `x` थांबवा, `r` पुन्हा लोड करा |
| Groups | प्रत्येक गट, तो काय शेअर करतो आणि त्याची खाती, आणि कोणत्याही गटात नसलेली खाती | `enter` संपादित करा, `n` नवीन, `e` नाव बदला, `x` काढा |
| Folders | [`ccam folder`](/features/folder) चे नियम आणि इथे लागू होणारा नियम, तुम्ही symbolic link मधून फोल्डरमध्ये पोहोचला असाल तरीही (नवीन नियम links resolve करून साठवले जातात) | `n` इथे एखादे खाते वापरा, `enter` बदला, `x` काढा |
| Extensions | [`ccam ext`](/features/ext) जे दाखवते ते, प्रत्येक extension एकदाच, ते ज्या गटांमधून येते त्यांच्यासह, `off:` म्हणजे ज्या खात्यांमध्ये plugin चालू नाही ती खाती, फक्त या फोल्डरला लागू होणारे खात्याचे स्वतःचे सर्व्हर, `This folder's project`, आणि कोणत्याही खात्याशी संबंधित नसलेल्या समस्या | `r` पुन्हा लोड करा |
| Health | [`ccam doctor`](/features/doctor) जे नोंदवते ते, प्रत्येक खात्यासाठी एक box म्हणून (खाली) | `f` निवडलेला box दुरुस्त करा, `F` सर्व दुरुस्त करा (`doctor --fix`), `↑`/`↓` boxes मध्ये फिरा |

`←`/`→` किंवा `Tab` tabs बदलतात, `?` सध्याच्या tab च्या की दाखवतो, `r` छोटा cache टाळून वापर ताजा करतो आणि Agents, Extensions आणि Health tabs पुन्हा लोड करतो, आणि `q` किंवा `Esc` बाहेर पडतो.

## Accounts tab वरील की

| की | काय करते |
| --- | --- |
| `Enter` | निवडलेल्या खात्याने Claude Code सुरू करा |
| `1`–`9` | त्या क्रमांकाच्या खात्यावर जा; 12 साठी `1` `2` पटकन टाइप करा. फक्त `Enter` सुरू करतो |
| `↑`/`↓`, `k`/`j` | हलवा |
| `n` | [`ccam next --auto`](/features/next) जे खाते निवडेल ते सुरू करा |
| `o` | पर्यायांसह सुरू करा: `--continue`, `--resume`, इतर आर्ग्युमेंट्स, किंवा dry run |
| `u` | वापराचा तपशील: 5 तासांच्या, साप्ताहिक आणि model windows, त्यांचे resets आणि आकडे कुठून येतात |
| `a` | खाते जोडा (खालील विझार्ड) |
| `l` | पुन्हा लॉग इन करा |
| `s` | खात्याचे गट निवडा |
| `d` | ते डीफॉल्ट खाते करा |
| `f` | ते सध्याच्या फोल्डरमध्ये वापरा |
| `e` | त्याचे नाव बदला |
| `c` | त्याची शॉर्टकट कमांड सेट करा किंवा काढा |
| `x` | ते काढा (आधी विचारते, No निवडलेले असते) |

**Team आणि Enterprise सीट्स.** आपोआप होणारी खाते निवड (`ccam next`, `--fallback`, `ccam pick`, डॅशबोर्डमधील `n` आणि, तुम्ही चालू केल्यावर, [`auto_switch`](/features/next#auto-switch)) एकाच Claude Team किंवा Enterprise organization मधील सीट्सदरम्यानही खाते बदलते, आणि यामुळे Anthropic च्या अटींचे उल्लंघन होऊ शकते. त्यावर अवलंबून राहण्यापूर्वी [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) वाचा आणि तुमच्या संस्थेकडे आणि Anthropic कडे चौकशी करा.

## Health tab

तुम्ही Health tab उघडता तेव्हा तो [`ccam doctor`](/features/doctor) च्या तपासण्या चालवतो आणि तेच दृश्य काढतो: प्रत्येक खात्यासाठी एक चिन्ह असलेली पहिली ओळ, आणि मग प्रत्येक खात्यासाठी एक box. box मधील प्रत्येक ओळ तपासलेले एक क्षेत्र आहे, आधी `Folder`, `Login` आणि `Share links`, ज्यावर `✔` ठीक, `ℹ` माहिती, `⚠` इशारा किंवा `✘` त्रुटी असे चिन्ह असते, त्यामुळे कोणत्या खात्यात समस्या आहे आणि कुठे ते तुम्हाला लगेच दिसते. कोणत्याही खात्याशी संबंधित नसलेले निष्कर्ष `General` box मध्ये जातात.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` एका box मधून दुसऱ्या box मध्ये नेतात, आणि निवडलेल्या box ची border accent रंगाची असते. `f` फक्त निवडलेला box दुरुस्त करतो, `ccam doctor --fix` ज्या प्रकारे करतो त्याच प्रकारे: एखाद्या खात्याच्या box वर तो फक्त तेच खाते दुरुस्त करतो, आणि `General` box वर कोणत्याही खात्याशी संबंधित नसलेले निष्कर्ष (plugin paths आणि शॉर्टकट्स); बाकी प्रत्येक खाते फक्त पुन्हा तपासले जाते. तो आधी विचारतो (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), आणि दुरुस्त करण्यासारखे काहीच नसलेला box फक्त `nothing to repair in <name>` असे सांगतो. `F` सर्व काही एकाच वेळी दुरुस्त करतो, `ccam doctor --fix` प्रमाणे, आणि `r` पुन्हा तपासतो.

## Agents tab

Agents tab [`ccam agents`](/features/agents) जे छापते ते दाखवतो: प्रत्येक खात्याची background सत्रे आणि उघडी सत्रे, प्रत्येक सत्र कोणत्या खात्याचे आहे त्यासह. तुम्ही tab पहिल्यांदा उघडता तेव्हा, `r` वर, एखादे सत्र थांबवल्यानंतर, आणि config बदलल्यानंतर तो प्रत्येक खात्याच्या Claude Code ला एकाच वेळी विचारतो (`claude agents --json`, प्रत्येकासाठी 15 सेकंदांचा timeout), आणि सर्वांनी उत्तर देईपर्यंत `Asking every account's Claude Code for its sessions…` असे सांगतो.

```console
  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Background agents and open sessions of every account; enter attaches to one.

    Account   ID        Name   State          Folder      Started
 ❯  work      -         -      open (busy)    ~/code/web  now
    personal  49ec0028  probe  blocked, idle  ~/code/app  12m

 ✘ lab: could not list its sessions: exit status 1 'claude agents --json' is
   not available in this environment.

 enter to attach · x to stop · r to reload · ? for shortcuts
```

`State` म्हणजे background सत्राची स्थिती आणि त्याची process काय करत आहे ते, किंवा interactive सत्रासाठी `open`, `ccam agents` प्रमाणे. एखाद्या खात्याला विचारता आले नाही, तर यादीखालील एक `✘` ओळ Claude Code ने काय सांगितले ते दाखवते, वरच्या `lab` प्रमाणे, ज्याचा agent view बंद केलेला आहे. पूर्ण झालेली background सत्रे वगळली जातात (`ccam agents --all` ती दाखवते), आणि प्रत्येक खात्याने उत्तर दिले आणि कोणाचेही सत्र नसेल, तर tab ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` असे सांगतो. अरुंद टर्मिनलवर tab आधी `Started`, मग `Name`, मग `Folder` काढून टाकतो, आणि लांब फोल्डर पाथ डावीकडून कापतो.

background सत्रावर `Enter` डॅशबोर्ड बंद करतो आणि ते सत्र या टर्मिनलमध्ये त्याच्या खात्याने उघडतो, `ccam attach <id>` प्रमाणे; डॅशबोर्ड परत येत नाही (Claude Code मध्ये `←` त्याच्या agent view वर जातो आणि `Ctrl+Z` shell वर). उघड्या सत्राला ID नसतो, त्यामुळे `Enter` आणि `x` फक्त background सत्र लागते एवढेच सांगतात. `x` आधी विचारतो, `No, keep it` निवडलेले ठेवून: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` हो म्हटल्यावर ccam त्या खात्याने `claude stop <id>` चालवते, ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` किंवा ते का होऊ शकले नाही ते सांगते, आणि यादी पुन्हा वाचते.

## गट बदलणे

Groups tab वर `Enter` editor उघडतो: डावीकडे आयटम, उजवीकडे खाती, बदलण्यासाठी `Tab`, टॉगल करण्यासाठी `Space`. एखाद्या खात्याला दुसऱ्या गटातून आधीच मिळणारा आयटम कारणासह राखाडी दाखवला जातो, आणि एखादा आयटम दोनदा मिळेल असे खातेही; एक आयटम फक्त एकाच गटातून येऊ शकतो. `Enter` योजना दाखवतो:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

एखादी फाइल गटाच्या प्रतीपेक्षा वेगळी असेल, तर ccam आधी काय करायचे ते विचारते (गटाची ठेवा, माझी वापरा, किंवा merge करा). मग बदल [`ccam group`](/features/share) च्याच code ने पार्श्वभूमीत चालतो: आधी config सेव्ह होते, प्रत्येक खात्याचा डेटा क्रमाने हलतो, आणि मूळ फाइल्स `.ccam-backup/` मध्ये जातात; MCP शेअरिंग सोडणारे खाते गटाचे सर्व्हर आपल्या `.claude.json` मध्ये ठेवते.

चालू खात्याचा डेटा हलवणारा बदल त्या खात्याच्या नावासह नाकारला जातो, योजना उघडी असताना ते खाते सुरू झाले तरीही, आणि दुसऱ्या ccam process ने गट बदलल्यामुळे जुनी झालेली योजना `the groups changed while this plan waited; nothing was moved, make the change again` देऊन नाकारली जाते. दुसरी ccam process गट बदलत असताना, नवीन बदल, नाव बदलणे किंवा काढणे `another ccam process is changing the groups right now` असे सांगते; एखादा बदल अर्ध्यावर थांबल्यानंतर, Health tab वरील `F` (किंवा `ccam doctor --fix`) तो पूर्ण करेपर्यंत ते थांबतात.

## खाते जोडणे

`a` तोच विझार्ड उघडतो जो `ccam add` सुद्धा टर्मिनलवर वापरते:

1. खात्याचे नाव.
2. ते कसे साइन इन करते: आत्ता Claude subscription (browser मध्ये), [API key](/features/api), किंवा नंतर Claude subscription.
3. API key साठी: प्रत्येक launch ला ccam key कुठून वाचते (`env:NAME` किंवा `command:COMMAND`), endpoint, key कशी पाठवली जाते आणि model.
4. सामील होण्यासाठी गट, ज्यात `default` निवडलेला असतो (अजून एकही गट नसेल तर, `default` गट काय शेअर करेल ते).
5. हवी असल्यास शॉर्टकट कमांड.
6. पुष्टी करण्यासाठी सारांश.

`Esc` एक टप्पा मागे जातो, आणि पहिल्या टप्प्यावर विझार्ड बंद करतो.

## खाते काढणे

काहीही होण्यापूर्वी `x` विचारतो, No निवडलेले ठेवून, आणि काढल्याने काय होते ते सांगतो: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` API key खात्यासाठी रद्द करण्यासारखे लॉगिन नसते, आणि दुसऱ्या खात्याकडे तेच लॉगिन असेल तेव्हा ते लॉग आउट केले जात नाही. इम्पोर्ट केलेले खाते आपला फोल्डर ठेवते: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## रंग

रंग Claude Code च्या terminal UI मधील भूमिकांनुसार असतात, त्यामुळे ccam त्याच्या शेजारी अगदी घरच्यासारखे दिसते; गटांच्या नावांना त्यांच्या नावावरून रंग मिळतो. ccam टर्मिनलला त्याची पार्श्वभूमी विचारते आणि light किंवा dark theme निवडते. `config.toml` च्या सुरुवातीला, पहिल्या `[[group]]`, `[[profile]]` किंवा `[[folder]]` च्या वर, स्वतः theme सेट करा:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

`-daltonized` themes लाल-हिरव्या रंगांधळेपणासाठी आहेत, आणि `-ansi` themes टर्मिनलचे स्वतःचे 16 रंग वापरतात. `NO_COLOR` रंग बंद करते.

## सुरक्षितता आणि अपवादात्मक स्थिती

- डॅशबोर्ड कमांड्सचीच functions वापरतो, त्यामुळे प्रत्येक नियम लागू राहतो: डेटा हलण्यापूर्वी योजना, चालू खात्यांसाठी नकार, बॅकअप, आणि खात्याबाहेर कधीच न जाणारे लॉगिन tokens.
- डॅशबोर्डमध्ये dry run नाही: नावाशिवाय `ccam --dry-run` या फोल्डरचे खाते कसे सुरू होईल ते छापते, आणि एखाद्या खात्यावर `o` दाबल्यास फक्त त्या एका सुरुवातीचा dry run करण्याचा पर्याय मिळतो.
- Claude Code सुरू होण्यापूर्वी किंवा लॉगिन उघडण्यापूर्वी टर्मिनल पूर्ववत होतो.
- dialog साठी स्क्रीन खूप लहान असेल तर आधी tab च्या ओळी वगळल्या जातात; ओळी wrap होण्याऐवजी `…` ने कापल्या जातात. कमी उंचीच्या स्क्रीनवर त्यानंतर रिकाम्या ओळी आणि tab चे वर्णन वगळले जाते, त्यामुळे keys सुचवणारी ओळ दिसत राहते. 40 columns किंवा 8 rows पेक्षा कमी असल्यास डॅशबोर्ड फक्त `Too small: ccam needs 40x8 (now 36x7)` दाखवतो.
- बदल, काढणे, अपडेट किंवा दुरुस्ती लागू होत असताना दाबलेले `Ctrl+C` ते पूर्ण होण्याची वाट पाहते: स्टेटस लाइन `ccam quits when it is done.` असे सांगते.
- `U` ने ccam अपडेट केल्यानंतर, खाते जोडून किंवा लॉग इन करून डॅशबोर्ड परत आल्यावर ऑफर पुन्हा दिसत नाही.

## संबंधित

- [गटांसह शेअरिंग](/features/share)
- [Claude Code सुरू करणे](/features/launch)
- [खाती व्यवस्थापित करणे](/features/accounts)
- [वापर मर्यादा](/features/usage)
