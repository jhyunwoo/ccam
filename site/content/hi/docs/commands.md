---
title: कमांड संदर्भ
metaTitle: "ccam कमांड संदर्भ"
description: ccam की हर कमांड और विकल्प, ccam add और ccam <name> से लेकर ccam usage, ccam next और ccam doctor तक।
updated: 2026-11-19
order: 3
translation: machine
sourceHash: "28e1ad1dcbad"
---

किसी कमांड को असल में चलाए बिना यह देखने के लिए कि वह क्या करेगी, उसके आगे `--dry-run` लगाएँ, उदाहरण के लिए `ccam --dry-run work`। यह Claude Code शुरू करते समय (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) और `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` तथा `doctor --fix` के साथ काम करता है। केवल पढ़ने वाली कमांड, जैसे `list`, `usage`, `share` और `ext`, इसे स्वीकार करती हैं और वैसे भी कुछ नहीं बदलतीं। बाक़ी कमांड इसे स्वीकार नहीं करतीं और चलने से मना कर देती हैं, जिनमें `group rename` भी शामिल है।

<!-- commands -->

## Claude Code शुरू करना

`ccam <name|number> [arguments]` Claude Code को उस अकाउंट से शुरू करता है। ccam नाम के बाद वाले आर्ग्युमेंट नहीं पढ़ता; वह उन्हें उसी क्रम में `claude` को दे देता है, और standard input, exit code और signals ठीक वैसे ही काम करते हैं जैसे ख़ुद `claude` के साथ। नाम के बाद `--` लगाना वैकल्पिक है।

आर्ग्युमेंट के साथ लेकिन बिना अकाउंट नाम के `ccam`, जैसे `ccam -p "summarise"`, मौजूदा फ़ोल्डर के लिए सेट किया गया अकाउंट इस्तेमाल करता है, वरना डिफ़ॉल्ट अकाउंट। इन्हें `ccam folder <name> [folder]` और `ccam use <name>` से सेट करें।

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

अकाउंट नाम से पहले `--fallback`, `-p`/`--print` runs के लिए है। जब अकाउंट अपनी उपयोग सीमा तक पहुँचता है, तो ccam वही कमांड सबसे ज़्यादा गुंजाइश वाले लॉग इन अकाउंट के रूप में फिर से चलाता है और सिर्फ़ उस run का output छापता है जो पूरा हुआ। [फ़ोल्डर के हिसाब से अकाउंट](/features/folder) और [सीमा पर फिर से चलाना](/features/next#run-again-at-the-limit) देखें।

**Team और Enterprise सीटें।** अपने-आप अकाउंट चुनना (`ccam next`, `--fallback`, `ccam pick`, डैशबोर्ड का `n` और, चालू करने के बाद, [`auto_switch`](/features/next#auto-switch)) एक ही Claude Team या Enterprise संगठन की सीटों के बीच भी स्विच करता है, जो Anthropic की शर्तों का उल्लंघन हो सकता है। इस पर भरोसा करने से पहले [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) पढ़ें और अपने संगठन और Anthropic से पता कर लें।

`ccwork` जैसा शॉर्टकट अपने अकाउंट को इसी तरह शुरू करता है: `ccwork --resume` वही करता है जो `ccam work --resume` करता है। नीचे [शॉर्टकट](/docs/commands#shortcuts) देखें।

## अकाउंट जोड़ना और हटाना

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

फ़ोल्डर के बिना `ccam import`, `~/.claude` और `~/.claude-*` खोजता है और हर एक के बारे में पूछता है; जो फ़ोल्डर विफल होता है उसे `Skipped <path>: <reason>` के रूप में बताया जाता है और बाक़ी रजिस्टर हो जाते हैं, अगर किसी फ़ोल्डर का नाम पहले से लिया जा चुका है, तो उस फ़ोल्डर को अलग से `ccam import <folder> --name <name>` से import करें। Import किए गए फ़ोल्डर जहाँ हैं वहीं रजिस्टर होते हैं; ccam कभी लॉगिन को move या copy नहीं करता।

`ccam login --expired` उन अकाउंट्स में भी लॉग इन करता है जिनका लॉगिन Claude Code ने हटा दिया। `ccam remove` किसी managed अकाउंट (जिसका फ़ोल्डर ccam ने बनाया हो, न कि `ccam import` से रजिस्टर किया गया अकाउंट) को उसका फ़ोल्डर trash में जाने से पहले लॉग आउट करता है, जब तक कोई दूसरा अकाउंट वही लॉगिन न रखता हो। import किए गए अकाउंट को हटाने पर वह अपने ग्रुप से बाहर हो जाता है और उसके फ़ोल्डर में हर साझा आइटम की एक कॉपी रह जाती है। macOS पर `ccam rename` फ़ोल्डर का नाम बनाए रखता है, क्योंकि कीचेन लॉगिन को फ़ोल्डर के path के तहत रखता है।

`ccam add --api` ऐसा अकाउंट जोड़ता है जो Claude लॉगिन के बजाय API key इस्तेमाल करता है; ccam हर launch पर key पढ़ता है और उसे कभी सेव नहीं करता। [API key अकाउंट](/features/api) देखें।

जब Claude Code उस अकाउंट के रूप में चल रहा हो, तब `ccam rename`, `ccam remove` और `ccam group` कमांड उस अकाउंट को बदलने से मना कर देते हैं, चाहे वह किसी भी तरह शुरू हुआ हो (background sessions समेत)। पहले session बंद करें। जब ग्रुप का कोई बदलाव अधूरा हो, तब भी वे मना कर देते हैं: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``।

`ccam list --json` हर अकाउंट को उसके स्लॉट, शॉर्टकट, ईमेल और फ़ोल्डर के साथ दिखाता है। इसकी `groups` array अकाउंट के ग्रुप बताती है, और इसका `share` ऑब्जेक्ट हर आइटम (`mcp`, `settings`, `memory`, `extensions`, `history`) को उस ग्रुप से जोड़ता है जिससे वह आता है, या `local` से।

## शॉर्टकट

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

शॉर्टकट एक कमांड है जिसका नाम आप खुद रखते हैं और जो Claude Code को एक अकाउंट से शुरू करती है। `ccam work` की तरह यह हर आर्ग्युमेंट बिना बदलाव के Claude Code को देती है, इसलिए `ccwork --resume` और `ccwork -p "hi"` वैसे ही काम करते हैं जैसी आप उम्मीद करते हैं। यह `~/.local/bin` में ccam का एक लिंक है।

शॉर्टकट वैकल्पिक हैं। `ccam add` का विज़र्ड पूछता है कि शॉर्टकट बनाना है या नहीं, `ccam add <name> --shortcut <command>` बिना पूछे एक शॉर्टकट बनाता है, और `ccam list` हर अकाउंट का शॉर्टकट दिखाता है, या शॉर्टकट न होने पर `-`।

नामों में a-z, 0-9, `.`, `-` और `_` अक्षर इस्तेमाल होते हैं। नाम कोई ccam कमांड, किसी दूसरे अकाउंट का शॉर्टकट, या आपके `PATH` पर मौजूद किसी दूसरे प्रोग्राम का नाम नहीं हो सकता, और न ही `~/.local/bin` में कोई ऐसी फ़ाइल जिसे ccam ने नहीं बनाया, जैसे ccam की आपकी अपनी कॉपी या लिंक। `config.toml` में हाथ से बदला गया ऐसा शॉर्टकट जो इन नियमों को तोड़ता है, `claude` या `ccam` है, या दो अकाउंट्स को दिया गया है, `config.toml` को लोड होने से रोकता है, और अकाउंट का नाम बताता है।

`ccam rename` शॉर्टकट का नाम वही रखता है, और `ccam remove` अकाउंट के साथ उसका शॉर्टकट भी हटा देता है। `ccam shims install` सभी अकाउंट्स के लिंक फिर से बनाता है, जैसे ccam को दूसरी जगह ले जाने के बाद; किसी एक शॉर्टकट के रास्ते में आई फ़ाइल बाक़ी शॉर्टकट को नहीं रोकती।

## साझा करना

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

एक ग्रुप अपने अकाउंट्स के बीच आइटम साझा करता है; एक अकाउंट कई ग्रुप में हो सकता है, लेकिन हर आइटम केवल एक ही ग्रुप से आता है। [ग्रुप के साथ साझा करना](/features/share) देखें।

## मौजूदा सेटअप कॉपी करना

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` किसी Claude Code फ़ोल्डर से `settings.json`, `CLAUDE.md`, एक्सटेंशन और user MCP सर्वर एक ग्रुप में कॉपी करता है, ताकि उन आइटम को साझा करने वाले हर अकाउंट को वे मिल जाएँ। फ़ोल्डर को केवल पढ़ा जाता है।

- `--settings`, `--memory`, `--extensions`, `--mcp` और `--history` केवल उन्हीं आइटम को कॉपी करते हैं जिनका आप नाम देते हैं। इनके बिना, टर्मिनल पर एक चेकलिस्ट आपको चुनने देती है, जिसमें बातचीत के इतिहास को छोड़कर सब कुछ चुना हुआ होता है; बाकी जगहों पर बातचीत के इतिहास को छोड़कर सब कुछ कॉपी होता है।
- `--mcp` तब भी काम करता है जब `~/.claude` खाली हो या मौजूद न हो, लेकिन `~/.claude.json` में सर्वर हों।
- `-g` या `--group` ग्रुप चुनता है।
- `-y` या `--yes` बिना पूछे कॉपी करता है और कोई फ़ाइल अलग होने पर ग्रुप की कॉपी रखता है।

MCP साइन-इन कॉपी नहीं होते, क्योंकि कॉपी किया गया साइन-इन बाक़ी धारकों को लॉग आउट कर देता है; plan उन अकाउंट्स के नाम बताता है जो `/mcp` से साइन इन करते हैं। [मौजूदा सेटअप कॉपी करना](/docs/sharing#copy-an-existing-setup) देखें।

## MCP सर्वर

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` उन सर्वरों को उसके ग्रुप में move करता है जो केवल उसी अकाउंट के पास हैं, जिनमें वे भी शामिल हैं जो उसके पास केवल एक फ़ोल्डर के लिए हैं (सादा `claude mcp add`), जो फिर हर फ़ोल्डर में लोड होते हैं। जिस सर्वर का नाम ग्रुप पहले से किसी अलग definition के लिए इस्तेमाल करता है, उसे छोड़ दिया जाता है; ग्रुप की कॉपी बदलने के लिए `--force` जोड़ें, जिसका ccam पहले बैकअप लेता है।

## उपयोग और स्विच करना

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

- `next --auto`, `--fallback` और `ccam pick` एक ही तरह से चुनते हैं: पहले `auto_max_usage` (`config.toml`, डिफ़ॉल्ट 100) से नीचे वाले अकाउंट, `auto_strategy` (`headroom`, `expiring` या `balance`) के हिसाब से, `--model` की साप्ताहिक विंडो गिनते हुए। वे कभी 100% वाला अकाउंट (जिस window को Claude Code ने भरा हुआ बताया, वह अपने reset तक 100% गिनी जाती है) या API key लॉगिन नहीं चुनते, और उसी अकाउंट के संगठन के भीतर रहते हैं जिससे वे शुरू करते हैं (आगे देखने के लिए `--cross-org` या `cross_org = true`)।
- `--wait` अपना मान `=` के बाद लेता है (`--wait=2h`); space के साथ लिखने पर यह एक संकेत के साथ अस्वीकार होता है। जब कोई नहीं चुना जा सकता, तो वे code 75 के साथ बाहर निकलते हैं, या, `--wait` के साथ, पहले अकाउंट के reset होने का इंतज़ार करते हैं।
- `ccam next --continue` बातचीत को उन अकाउंट्स के बीच भी आगे ले जाता है जो इतिहास साझा नहीं करते: तब Claude Code उसे नए अकाउंट में कॉपी कर देता है। `--brief` के साथ, नया अकाउंट इसके बजाय एक नई बातचीत शुरू करता है, उस सारांश से जो ccam किसी मॉडल के बिना पिछली बातचीत का लिखता है ([`--brief`](/features/next#brief))।
- `ccam <name>` के उलट, `ccam next` अपने विकल्प ख़ुद पढ़ता है, इसलिए Claude Code के विकल्प `--` के बाद आते हैं; इसके बिना ccam `unknown flag` के साथ रुक जाता है।
- `ccam agent` इस फ़ोल्डर में `claude -p` से एक task उस अकाउंट के अलावा किसी दूसरे अकाउंट के रूप में चलाता है जिससे उसे बुलाया गया है, बुलाने वाले session के permission mode और effort में, और Claude का आख़िरी जवाब प्रिंट करता है। उस अकाउंट की सीमा पर बातचीत किसी दूसरे अकाउंट के रूप में आगे चलती है ([दूसरे अकाउंट्स पर subagents](/features/agent))।
- Claude Code के अंदर, plugin का `/ccam:switch [account]` बातचीत आगे सौंपता है, और `ccam hook` वह है जिसे उसके hooks चलाते हैं।
- `config.toml` में `auto_switch = true` होने पर (Linux और macOS; डिफ़ॉल्ट रूप से बंद), limit hook ccam से शुरू हुए interactive session को भी बंद करता है और उसे गुंजाइश वाले किसी दूसरे अकाउंट पर जारी रखता है, लेकिन केवल उसी तरह के अकाउंट पर: उसी संगठन की कोई दूसरी सीट, या कोई दूसरा personal प्लान ([`auto_switch`](/features/next#auto-switch))।

[अगले अकाउंट पर स्विच करना](/features/next), [scripts के लिए अकाउंट चुनना](/features/pick), [दूसरे अकाउंट्स पर subagents](/features/agent) और [हर अकाउंट के background sessions](/features/agents) देखें।

## रखरखाव

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

ccam को बदलने से पहले `ccam update` रिलीज़ का SHA-256 checksum जाँचता है, और `--version v0.2.0` कोई खास रिलीज़ इंस्टॉल करता है। `.deb` पैकेज से इंस्टॉल की गई ccam की कॉपी इसके बजाय apt से अपडेट होता है। [अपडेट](/docs/install#update) देखें।

`ccam uninstall` दिखाता है कि वह क्या हटाएगा और पहले पूछता है: `ccam remove` की तरह हर अकाउंट, शॉर्टकट कमांड, ccam के फ़ोल्डर और `ccam` कमांड। import किया गया फ़ोल्डर अपनी जगह पर बना रहता है। [अनइंस्टॉल](/docs/install#uninstall) देखें।

## फ़ाइलें

| पाथ | सामग्री |
| --- | --- |
| `~/.config/ccam/config.toml` | अकाउंट्स, ग्रुप और सेटिंग्स (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …)। सेटिंग्स सबसे ऊपर, पहले `[[group]]`, `[[profile]]` या `[[folder]]` से ऊपर आती हैं: किसी table के नीचे वाली सेटिंग उसी table की होती है, इसलिए ccam हर launch पर और `ccam doctor` में उसके बारे में चेतावनी देता है, और config सेव करने वाली कमांड तब तक मना करती हैं जब तक उसे ऊपर न ले जाया जाए। ccam फ़ाइल सेव करते समय comments नहीं रखता |
| `~/.local/share/ccam/profiles/<name>/` | हर अकाउंट का `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | साझा सेटिंग्स, MCP सर्वर, एक्सटेंशन और इतिहास |
| `~/.local/state/ccam/` | उपयोग का cache, limit records, locks, नोट्स और दूसरी state (table के नीचे की सूची देखें) |
| `~/.local/bin/<shortcut>` | `ccwork` जैसी शॉर्टकट कमांड, ccam के लिंक |

state फ़ोल्डर में ये रहते हैं:

- उपयोग का cache और limit records (`usage/<name>.json`, `usage/<name>.limit.json`);
- आख़िरी इस्तेमाल का समय;
- session markers और locks (`locks/`);
- ccam के बनाए शॉर्टकट (`shortcuts.json`);
- ग्रुप का कोई बदलाव जो पूरा नहीं हुआ (`group-change.json`);
- `/ccam:switch` और limit hook के छोड़े गए नोट्स (`handoff/<account>/<session>.json`);
- `ccam pick` या `auto_switch` ने हर अकाउंट को आख़िरी बार कब चुना (`picks/<name>`; एक चुनाव 2 मिनट तक गिना जाता है);
- `--brief` के सारांश, 7 दिन तक रखे जाते हैं (`brief/<session>.md`)।

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` या `XDG_STATE_HOME` सेट करने से इन फ़ोल्डरों की जगह बदल जाती है। Windows पर `~` का मतलब `%USERPROFILE%` है और शॉर्टकट `<shortcut>.exe` होते हैं।
