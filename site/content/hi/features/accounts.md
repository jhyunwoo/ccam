---
title: अकाउंट प्रबंधित करना
description: ccam add, ccam login, ccam list, ccam rename, ccam remove और ccam import से Claude Code अकाउंट जोड़ें, लॉगिन करें, सूची देखें, नाम बदलें, हटाएँ व import करें।
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

हर ccam अकाउंट एक अलग Claude Code config फ़ोल्डर है, जिसका अपना लॉगिन होता है। इन कमांड से अकाउंट जोड़ें, उनमें साइन इन करें, उनकी सूची देखें, उनका नाम बदलें और उन्हें हटाएँ। ये आपके पास पहले से मौजूद फ़ोल्डर भी रजिस्टर करती हैं।

## उपयोग

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account that was never logged in or whose login expired
ccam logout side                           # log out: the login and its MCP sign-ins
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, logs out, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

जहाँ भी कोई कमांड `<name|number>` लेती है, वहाँ आप अकाउंट का नाम या उसका स्लॉट नंबर (`ccam list` में `#` कॉलम) दे सकते हैं।

`ccam use` केवल यह बदलता है कि डिफ़ॉल्ट अकाउंट कौन-सा है (`config.toml` में `default`)।

टर्मिनल पर बिना विकल्पों के `ccam add` एक विज़र्ड खोलता है। यह इन चरणों से गुज़रता है:

1. अकाउंट का नाम, अगर आपने नहीं दिया;
2. अकाउंट कैसे साइन इन करता है: अभी Claude लॉगिन, API key या बाद में Claude लॉगिन;
3. किन ग्रुप में शामिल होना है, जिनमें `default` चुना हुआ होता है (अगर अभी कोई ग्रुप नहीं है, तो यह दिखाता है कि `default` ग्रुप क्या साझा करेगा);
4. `Create a shortcut command for work? [y/N]`;
5. पुष्टि के लिए एक सारांश।

`Esc` एक कदम पीछे ले जाता है। `--group`, `--no-group`, `--shortcut` या `--no-login` में से कोई भी विकल्प, या टर्मिनल के बिना चलाना (उदाहरण के लिए, pipe से), विज़र्ड को छोड़ देता है। तब नाम देना ज़रूरी है, और कोई शॉर्टकट तब तक नहीं बनता जब तक `--shortcut` उसका नाम न दे।

`--group` के बिना नया अकाउंट `default` ग्रुप में शामिल होता है। पहला अकाउंट इसे बनाता है, जो MCP सर्वर, सेटिंग्स, `CLAUDE.md` और एक्सटेंशन साझा करता है; [ग्रुप के साथ साझा करना](/features/share) देखें। Windows पर फ़ाइलें केवल Developer Mode में लिंक हो सकती हैं। इसके बिना, पहला अकाउंट केवल `mcp` और `extensions` के साथ `default` बनाता है और यह बताता है (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), और ccam सेटिंग्स, memory या इतिहास साझा करने वाले ग्रुप के साथ `--group` अस्वीकार करता है।

`ccam list` स्लॉट, नाम, शॉर्टकट, ईमेल, ग्रुप और आख़िरी उपयोग दिखाता है, और डिफ़ॉल्ट अकाउंट के आगे `*` लगाता है। ईमेल वह है जिससे फ़ोल्डर अभी लॉगिन है (उसकी `.claude.json` से); जब वह ccam के दर्ज किए ईमेल से अलग हो, तो दर्ज किया गया ईमेल `(was …)` के रूप में बाद में आता है:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

```console
$ ccam list
   #  Name      Shortcut  Email                                     Groups   Last used
*  1  work      ccwork    other@company.com (was work@example.com)  default  2m ago
```

`--json` के साथ हर अकाउंट एक ऑब्जेक्ट होता है, जिसमें `slot`, `shortcut`, `name`, `email` (मौजूदा लॉगिन का, जैसा table में), `default`, `dir`, `last_used` (जो अकाउंट कभी शुरू नहीं हुआ, उसमें यह नहीं होता), `groups` और `share` होते हैं। `share` हर आइटम (`mcp`, `settings`, `memory`, `extensions`, `history`) को उस ग्रुप से जोड़ता है जिससे वह आता है, या `"local"` से।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `ccam add --no-login` | लॉगिन किए बिना अकाउंट बनाएँ; बाद में `ccam login <name>` से लॉगिन करें | बंद |
| `ccam add --group <group>` | `default` के बजाय इस ग्रुप में शामिल हों। कई बार दिया जा सकता है; `default` न हो तो बनाया जाता है | `default` |
| `ccam add --no-group` | किसी ग्रुप में शामिल न हों: अकाउंट कुछ साझा नहीं करता | बंद |
| `ccam add --shortcut <command>` | `~/.local/bin` में यह शॉर्टकट कमांड भी बनाएँ | कोई नहीं |
| `ccam add --api --key-from <source>` | ऐसा अकाउंट जो लॉगिन के बजाय API key इस्तेमाल करता है; [API key अकाउंट](/features/api) देखें | बंद |
| `ccam list --json` | अकाउंट्स को JSON के रूप में प्रिंट करें | बंद |
| `ccam remove -y`, `--yes` | बिना पूछे हटाएँ | बंद |
| `ccam logout -y`, `--yes` | बिना पूछे लॉग आउट करें | बंद |
| `ccam import --name <name>` | आपके दिए फ़ोल्डर के लिए अकाउंट का नाम; फ़ोल्डर के बिना अस्वीकार होता है (`--name needs a folder: ccam import <folder> --name ops`) | फ़ोल्डर के नाम से बना |
| `ccam import -y`, `--yes` | मिले हर फ़ोल्डर को बिना पूछे import करें | बंद |
| `--dry-run` | `ccam list` और `ccam completion` इसका समर्थन करते हैं, क्योंकि वे कुछ नहीं बदलते। `add`, `login`, `logout`, `use`, `rename`, `remove` और `import` इसे `` --dry-run is not supported by `ccam add`; nothing was changed `` (कमांड के नाम के साथ) कहकर अस्वीकार करते हैं | बंद |

`ccam login` (`--expired` के अलावा), `ccam use` और `ccam rename` के अपने कोई विकल्प नहीं हैं। `ccam completion` `bash`, `zsh` या `powershell` लेता है; कुछ और देने पर यह `only bash, zsh and powershell are supported` के साथ विफल होता है।

अकाउंट के नाम 1 से 32 अक्षरों के होते हैं, जिनमें `a-z`, `0-9`, `-` और `_` आते हैं, और वे किसी अक्षर या अंक से शुरू होते हैं। नाम केवल अंकों का नहीं हो सकता, क्योंकि `ccam 2` का मतलब स्लॉट 2 है। नाम कोई ccam कमांड भी नहीं हो सकता (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), ताकि `ccam <name>` कभी अस्पष्ट न हो। यह सूची `internal/config` में `config.Reserved` है।

## यह कैसे काम करता है

अकाउंट जोड़ने पर, इसी क्रम में, उसका फ़ोल्डर बनता है, वह अपने ग्रुप में शामिल होता है, config सेव होता है, साइन-इन होता है और शॉर्टकट बनता है। जब आप `ccam add work --shortcut ccwork` चलाते हैं, तो ccam ये चरण लेता है (`internal/cli/profiles.go` में `addProfile`):

<!-- motion:accounts-add -->

1. कुछ भी बनाने से पहले ccam नाम (`config.ValidateName`) और शॉर्टकट (`shim.Check`) जाँचता है।
2. ccam मोड `0700` के साथ `~/.local/share/ccam/profiles/work/` बनाता है और अकाउंट को सबसे छोटा ख़ाली स्लॉट नंबर देता है, इसलिए हटाए गए अकाउंट का नंबर दोबारा इस्तेमाल होता है (`profile.Add`, `Config.NextSlot`)। पहला अकाउंट डिफ़ॉल्ट बनता है।
3. अकाउंट अपने ग्रुप्स में शामिल होता है (`default`, जब तक `--group` या `--no-group` कुछ और न कहे; पहला अकाउंट `default` बनाता है), और ccam नए, ख़ाली फ़ोल्डर को उनसे लिंक करता है, उदाहरण के लिए `settings.json` को `groups/default/settings.json` से (`applyShare`, `share.Build` और फिर `Apply` के साथ)। अगर यह विफल हो, तो नया फ़ोल्डर हटा दिया जाता है, क्योंकि बचा हुआ फ़ोल्डर हर अगली कोशिश को `folder already exists` के साथ रोक देता।
4. अकाउंट `fsutil.AtomicWrite` के ज़रिए `~/.config/ccam/config.toml` में सेव होता है: एक अस्थायी फ़ाइल, `fsync`, फिर `rename`, ताकि crash होने पर भी आधा लिखा config कभी न बचे।
5. ccam अकाउंट फ़ोल्डर पर सेट `CLAUDE_CONFIG_DIR` के साथ, और अकाउंट के पिछले ईमेल वाले `--email` के साथ `claude auth login` शुरू करता है, ताकि लॉगिन पेज पहले से भरा हो (`runLogin`)। अगर लॉगिन विफल हो, तो अकाउंट बना रहता है, और ccam त्रुटि प्रिंट करके बताता है कि बाद में `ccam login work` चलाएँ।
6. ccam `profiles/work/.claude.json` से `oauthAccount` पढ़ता है (`profile.ReadIdentity`): ईमेल, अकाउंट और संगठन के ID, संगठन का नाम और प्लान (`claude_max`, `claude_team` …)। वह फ़ाइल कई मेगाबाइट की हो सकती है, इसलिए ccam पूरी फ़ाइल parse करने के बजाय केवल उस एक key को खोजता है (`TopLevelValue`)। ccam ये विवरण `config.toml` में `[profile.login]` में रखता है।
7. ccam `~/.local/bin/ccwork` को ccam बाइनरी से लिंक करता है, जो शुरू होते समय `config.toml` में `argv[0]` खोजती है, और लिंक को `~/.local/state/ccam/shortcuts.json` में दर्ज करता है (`shim.Install`)। अगर लिंक नहीं बन सकता, तो अकाउंट फिर भी जुड़ जाता है, ccam कारण बताता है (`Could not create shortcuts: …`) और सारांश `shortcut: ccwork (not created; see above)` दिखाता है।

लॉगिन में कई मिनट लग सकते हैं, और इस बीच कोई दूसरा ccam process config सेव कर सकता है। इसलिए इंतज़ार करने वाली कमांड `config.toml` दोबारा लोड करती हैं, केवल वह एक बदलाव (यहाँ, ईमेल) लागू करती हैं और सेव करती हैं (`env.update`)। इसी कारण `ccam remove` भी अपने सवाल के बाद दोबारा लोड करता है।

`ccam login` किसी मौजूदा अकाउंट के लिए चरण 5 और 6 चलाता है। Claude Code ईमेल को लॉगिन पेज के लिए केवल एक संकेत की तरह लेता है, इसलिए हो सकता है कि आप पहले से अलग किसी Claude अकाउंट में साइन इन हो जाएँ; ccam नए लॉगिन की तुलना `[profile.login]` से करता है और यह बताता है (`ccam: this login is other@company.com, but work was dev@company.com before.`)। यह नया लॉगिन रखता है और बताता है कि इसे कैसे पूर्ववत करें: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

किसी दूसरे Claude अकाउंट में लॉगिन करने पर अकाउंट का [limit record](/features/usage#the-limit-record) भी हट जाता है, क्योंकि वे सीमाएँ दूसरे अकाउंट की थीं। जब कोई दूसरा ccam अकाउंट उसी Claude अकाउंट में लॉगिन हो, तो ccam चेतावनी देता है कि दोनों अपनी उपयोग सीमाएँ साझा करते हैं।

`ccam login --expired` एक के बाद एक हर उस अकाउंट में लॉगिन करता है जिसमें कभी लॉगिन नहीं हुआ या जिसका लॉगिन समाप्त हो गया है: उसका refresh token समाप्त हो गया, या सर्वर के मना करने के बाद Claude Code ने लॉगिन हटा दिया। पहला ब्राउज़र खुलने से पहले यह जाँचता है कि `config.toml` सेव हो सकता है।

`ccam logout` अकाउंट में `claude auth logout` चलाता है और सेव किया गया ईमेल हटा देता है; लॉगिन का रिकॉर्ड बना रहता है, इसलिए अगला लॉगिन उससे जाँचा जाता है। Claude Code लॉगिन को revoke करता है और अकाउंट के लिए जो कुछ रखता है उसे मिटा देता है, MCP सर्वर साइन-इन समेत। इसलिए इन स्थितियों में ccam पहले चेतावनी देता है और पूछता है (`--yes` सवाल छोड़ देता है):

- लॉगिन में MCP साइन-इन हों;
- अकाउंट चल रहा हो (उसके sessions भी साइन आउट हो जाते हैं);
- किसी दूसरे अकाउंट का फ़ोल्डर वही लॉगिन रखता हो (कॉपी किया गया फ़ोल्डर), जिसे logout भी revoke कर देता है।

उदाहरण के लिए:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

किसी दूसरे Claude अकाउंट के रूप में लॉगिन करने के लिए अकेला `ccam login <name>` काफ़ी है, और यह MCP साइन-इन बनाए रखता है।

जब अकाउंट एक से ज़्यादा संगठनों के हों, तो `ccam list` एक `Organization` कॉलम जोड़ता है (personal प्लान `personal` दिखाते हैं), और `ccam list --json` में हर अकाउंट के लिए `plan` और `org` होते हैं।

### रिमोट मशीन पर साइन इन करना

SSH पर `claude auth login` का ब्राउज़र वाला हिस्सा टूट सकता है। Claude Code रिमोट मशीन पर एक callback सर्वर (`http://localhost:<port>/callback`) पर साइन-इन का इंतज़ार करता है, लेकिन साइन-इन पूरा करने वाला ब्राउज़र आपके कंप्यूटर पर चलता है, जहाँ वह पता कहीं नहीं ले जाता और ब्राउज़र `ERR_CONNECTION_REFUSED` जैसी त्रुटि दिखाता है। पते में एक बार इस्तेमाल होने वाला कोड फिर भी मौजूद रहता है, इसलिए ccam उसे पहुँचा सकता है:

<!-- motion:accounts-relay -->

ccam आपका टर्मिनल Claude Code को नहीं सौंपता (`runLogin`)। यह Claude Code के standard input को एक pipe से जोड़ता है और आपकी टाइप की गई हर पंक्ति उसमें कॉपी करता है (`relayLogin`), इसलिए Claude Code के `Paste code here` prompt पर paste किया गया कोड फिर भी पहुँचता है। जब कोई पंक्ति callback पता हो, तो ccam उसे केवल तभी स्वीकार करता है जब host `localhost`, `127.0.0.1` या `[::1]` हो, port दिया गया हो, path `/callback` हो और `code` तथा `state` दोनों मौजूद हों; `http://` छोड़ा जा सकता है (`callbackURL`)। फिर ccam सर्वर पर 10 सेकंड के timeout के साथ और redirects का पालन किए बिना उस पते पर request भेजता है, ठीक वैसे जैसे ब्राउज़र भेजता (`deliverCallback`)। Request इंतज़ार कर रहे `claude auth login` के callback सर्वर तक पहुँचती है (आरेख में डैश वाली रेखा)। Claude Code कोड पाता है और लॉगिन सेव करता है, Linux पर अकाउंट फ़ोल्डर की `.credentials.json` में। Paste की गई कोई पंक्ति ccam से कभी किसी दूसरे host से संपर्क नहीं करवा सकती।

लॉगिन शुरू होने से पहले ccam यह प्रिंट करता है:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

सफल होने पर यह `→ Passed the sign-in to Claude Code on this machine.` प्रिंट करता है।

जब लॉगिन टर्मिनल पर input का इंतज़ार करता है, तब ccam टर्मिनल का echo और line mode बंद कर देता है और आपके टाइप किए हर अक्षर के लिए एक `*` दिखाता है, ताकि paste किया गया code या callback पता स्क्रीन पर न दिखे। ccam line editing ख़ुद करता है (`readLine`): `Enter` पंक्ति भेजता है, `Backspace` एक अक्षर मिटाता है, और ख़ाली पंक्ति पर `Ctrl+D` input ख़त्म करता है। लॉगिन ख़त्म होने पर टर्मिनल पहले जैसी स्थिति में आ जाता है, और `Ctrl+C`, `SIGTERM` या `SIGHUP` पर भी, जो फिर भी ccam को बंद करते हैं। Standard input और standard error पर टर्मिनल न होने पर, पंक्तियाँ जैसे आती हैं वैसे ही कॉपी होती हैं।

### नाम बदलना, हटाना और import करना

`ccam rename` और `ccam remove` अकाउंट के फ़ोल्डर पर, और उस पर काम करते हैं जो ccam अकाउंट के नाम से दूसरी जगहों पर रखता है। ccam उन पाथ की सूची `nameState` में रखता है:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- उपयोग की cache फ़ाइलें `~/.local/state/ccam/usage/<name>.json` और `<name>.statusline.json`;
- limit record `<name>.limit.json`;
- `~/.local/state/ccam/handoff/<name>/` में session के नोट।

`ccam remove work` इन्हें इस क्रम में संभालता है:

<!-- motion:accounts-remove -->

1. `ccam remove` पहले अकाउंट चल रहा हो तो मना कर देता है, फिर पूछता है, और आपके जवाब के बाद फिर से जाँचता है (इससे सवाल खुला रहने के दौरान शुरू हुआ session पकड़ में आ जाता है)।
2. कुछ भी बदलने से पहले यह जाँचता है कि `config.toml` सेव हो सकता है।
3. अगर अकाउंट एक्सटेंशन साझा करता है, तो ccam ग्रुप की `installed_plugins.json` और `known_marketplaces.json` को फिर से लिखता है (`RewritePluginPaths`), क्योंकि Claude Code plugins के absolute पाथ सेव करता है जो इस अकाउंट के फ़ोल्डर से होकर जा सकते हैं।
4. इसके बाद managed अकाउंट (जिसका फ़ोल्डर ccam ने बनाया हो, न कि `ccam import` से रजिस्टर किया गया अकाउंट) लॉग आउट किया जाता है: उसके फ़ोल्डर में `claude auth logout` चलता है, ताकि कोई काम करता लॉगिन trash में (या macOS पर, फ़ोल्डर के नाम वाले कीचेन आइटम में, जिसे बाद में उसी नाम का अकाउंट विरासत में पाता) इंतज़ार न करे। API key अकाउंट के लिए, और जब किसी दूसरे अकाउंट का फ़ोल्डर वही लॉगिन रखता हो, जिसे logout भी revoke कर देता, यह छोड़ दिया जाता है। विफल logout बताया जाता है, और हटाना जारी रहता है।
5. ccam फ़ोल्डर को `~/.local/share/ccam/trash/<name>-<time>/` में move करता है (`profile.Trash`), अकाउंट `config.toml` से हट जाता है, `nameState` पाथ मिटा दिए जाते हैं ताकि बाद में उसी नाम का अकाउंट उन्हें विरासत में न पाए, और शॉर्टकट लिंक हटा दिया जाता है: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`।

अगर आख़िर में `config.toml` सेव नहीं हो पाता, तो फ़ोल्डर trash से वापस आ जाता है। हटाए गए अकाउंट को वापस पाने के लिए फ़ोल्डर को trash से वापस move करें और फिर से लॉगिन करें।

`ccam rename` फ़ोल्डर का नाम बदलता है और उसके `plugins/` में plugin पाथ को पुराने फ़ोल्डर से नए फ़ोल्डर पर फिर से लिखता है। यह `nameState` पाथ को नए नाम पर move करता है और अगर अकाउंट डिफ़ॉल्ट था तो उसे डिफ़ॉल्ट बनाए रखता है। यह session के नोट भी अपडेट करता है, ताकि `ccam next --continue` उन्हें फिर भी ढूँढ सके। शॉर्टकट कमांड का नाम वही रहता है। अगर plugin पाथ फिर से लिखना विफल हो, या `config.toml` सेव न हो सके, तो फ़ोल्डर का नाम वापस बदल दिया जाता है।

macOS पर, और Windows पर जब Claude Code लॉगिन Credential Manager में रखता है, फ़ोल्डर ख़ुद अपना पुराना नाम रखता है, क्योंकि लॉगिन और MCP साइन-इन फ़ोल्डर के path के तहत रखे जाते हैं और move किया गया फ़ोल्डर लॉग आउट हो जाता। तब config फ़ोल्डर को नए नाम के बगल में `folder` के रूप में दर्ज करता है, और ccam कहता है `its folder stays …: the macOS Keychain keeps the login under the folder's path`।

`ccam import` किसी फ़ोल्डर को जहाँ वह है वहीं, `config.toml` में `dir` के रूप में रजिस्टर करता है। यह उसे कभी move या copy नहीं करता: कॉपी किया गया refresh token कॉपियों को एक-दूसरे से लॉग आउट करवा देता।

फ़ोल्डर आर्ग्युमेंट के बिना यह बिना रजिस्टर वाले `~/.claude` और `~/.claude-*` फ़ोल्डर खोजता है और हर एक के बारे में पूछता है। नाम फ़ोल्डर से आते हैं (`~/.claude` का `main`, `~/.claude-work` का `work`, `~/.claude-2` का `acct2`)। अगर ccam नाम नहीं बना पाता, या नाम पहले से लिया जा चुका है, तो वह फ़ोल्डर को चलाने वाली `--name` कमांड के साथ दिखाता है। मिला हुआ जो फ़ोल्डर विफल होता है, उसे `Skipped <path>: <reason>` के रूप में बताया जाता है और बाक़ी import हो जाते हैं।

ccam पहले अपने सारे सवाल पूछता है, फिर `config.toml` दोबारा लोड करता है और फ़ोल्डर रजिस्टर करता है, इसलिए इस बीच दूसरी ccam कमांड ने जो बदलाव सेव किए वे बने रहते हैं। `✓` वाली पंक्तियाँ सेव के बाद आती हैं।

`CLAUDE_CONFIG_DIR` के बिना Claude Code `.claude.json` को आपके home फ़ोल्डर में `~/.claude` के बगल में रखता है, इसलिए जब आप `~/.claude` import करते हैं और उसकी अपनी कोई `.claude.json` नहीं होती, तो ccam `~/.claude/.claude.json` से `~/.claude.json` का लिंक बनाता है। तब सादा `claude` और `ccam main` एक ही अकाउंट देखते हैं।

अगर `~/.claude` में `~/.claude.json` के बगल में अपनी `.claude.json` है, तो ccam दोनों को छोड़ देता है और यह बताता है। `ccam main` अंदर वाली और सादा `claude` बाहर वाली इस्तेमाल करता, जिससे उनके लॉगिन विवरण और MCP सर्वर अलग रहते। एक रखने के लिए अंदर वाली फ़ाइल को लिंक से बदलें ([doctor](/features/doctor) भी यह बताता है)।

जब फ़ोल्डर में किसी रजिस्टर अकाउंट वाला ही लॉगिन हो (कॉपी किया गया फ़ोल्डर), तो ccam यह बताता है: जो पहले उसे renew करता है वह दूसरे को लॉग आउट कर देता है, इसलिए इस्तेमाल करने से पहले `ccam login <name>` चलाएँ।

Import किए गए अकाउंट तब तक किसी ग्रुप में नहीं होते जब तक आप उन्हें किसी ग्रुप में न जोड़ें (`ccam group join <group> <name>`), और कोई शॉर्टकट नहीं बनता।

`ccam completion bash|zsh` एक completion स्क्रिप्ट प्रिंट करता है जो अकाउंट के नाम भी पूरे करती है।

## सुरक्षा और विशेष स्थितियाँ

- जब अकाउंट चल रहा हो, तब `ccam rename` और `ccam remove` मना कर देते हैं (`refuseRunning`), ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)`` के साथ, क्योंकि session फ़ोल्डर में लिखता रहता है। यह जाँच Claude Code को पकड़ती है चाहे वह किसी भी तरह शुरू हुआ हो: ccam से, सादे `claude` से, किसी IDE से, background sessions और उनके daemon से।
- जब कोई दूसरा ccam process उस लॉगिन को renew कर रहा हो या बदल रहा हो, तब `ccam rename`, `ccam remove` और `ccam logout` मना कर देते हैं: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``।
- जब किसी दूसरे ccam process में ग्रुप का बदलाव चल रहा हो, तो `ccam rename` और `ccam remove` `another ccam process is changing the groups right now; try again when it is done` के साथ मना कर देते हैं; और जब कोई बदलाव बीच में रुक गया हो, तो ``the last group change stopped halfway; finish it with `ccam doctor --fix` first`` के साथ।
- Import किए गए अकाउंट पर `ccam remove` फ़ोल्डर और उसका लॉगिन सादे `claude` के लिए बनाए रखता है। अगर अकाउंट ग्रुप में है, तो वह पहले उन्हें वैसे ही छोड़ता है जैसे `ccam group leave` करता है: हर साझा आइटम फ़ोल्डर में ग्रुप वाले की कॉपी बन जाता है, और ग्रुप के MCP सर्वर उसकी `.claude.json` में चले जाते हैं। फिर उसका रजिस्ट्रेशन हटता है: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`। ccam ने उसकी `.claude.json` के जो बैकअप बनाए थे, वे `<folder>/.ccam-backup/mcp-<time>` में चले जाते हैं; अगर वे move नहीं हो सकते, तो वे `~/.local/state/ccam/mcp/<name>/backups` में रहते हैं और संदेश यह बताता है। आपका `~/.claude` trash में नहीं जाता।
- दो अकाउंट्स का एक ही फ़ोल्डर नहीं हो सकता: एक फ़ोल्डर का नाम दो बार देने वाला `config.toml` लोड नहीं होता।
- `ccam add` ऐसा नाम अस्वीकार करता है जिसका फ़ोल्डर `profiles/` में पहले से मौजूद है, और `ccam import` की ओर इशारा करता है। macOS पर यह तब भी रुक जाता है जब कीचेन में नए फ़ोल्डर के path के लिए ccam के किसी पुराने वर्ज़न का छोड़ा लॉगिन या MCP साइन-इन अभी भी हो, और उन्हें हटाने वाली `security delete-generic-password` कमांड बताता है। अगर `config.toml` सेव नहीं हो सकता, तो नया फ़ोल्डर फिर से हटा दिया जाता है।
- `ccam import` पहले से रजिस्टर फ़ोल्डर, ऐसा पाथ जो फ़ोल्डर नहीं है, और पहले से लिया गया नाम अस्वीकार करता है। जब फ़ोल्डर `~/.claude` की कॉपी हो (वही लॉगिन), तो ccam चेतावनी भी देता है, चाहे `~/.claude` ख़ुद रजिस्टर न हो।
- अकाउंट चल रहा हो तो `ccam logout` पहले पूछता है; `ccam login` नहीं जाँचता।
- अगर relay Claude Code तक नहीं पहुँच पाता, उदाहरण के लिए लॉगिन ख़त्म होने के बाद, तो यह कारण के साथ `ccam: could not pass the sign-in to Claude Code:` प्रिंट करता है। HTTP त्रुटि का आमतौर पर मतलब है कि कोड एक्सपायर हो गया; लॉगिन फिर से शुरू करें।

## संबंधित

- [Claude Code शुरू करना](/features/launch)
- [शॉर्टकट कमांड](/features/shortcuts)
- [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing)
- [एक कंप्यूटर पर कई Claude Code अकाउंट इस्तेमाल करें](/guides/multiple-claude-code-accounts)
