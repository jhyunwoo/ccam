---
title: खाती व्यवस्थापित करणे
description: ccam add, ccam login, ccam list, ccam rename, ccam remove आणि ccam import ने Claude Code खाती जोडा, लॉग इन करा, यादी पाहा, नाव बदला, काढा आणि इम्पोर्ट करा.
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

प्रत्येक ccam खाते हा स्वतःचे लॉगिन असलेला वेगळा Claude Code config फोल्डर आहे. खाती जोडण्यासाठी, त्यांत साइन इन करण्यासाठी, त्यांची यादी पाहण्यासाठी, नाव बदलण्यासाठी आणि ती काढण्यासाठी या कमांड्स वापरा. त्या तुमच्याकडे आधीच असलेले फोल्डरही नोंदवतात.

## वापर

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

जिथे एखादी कमांड `<name|number>` घेते, तिथे तुम्ही खात्याचे नाव किंवा त्याचा slot क्रमांक (`ccam list` मधील `#` कॉलम) देऊ शकता.

`ccam use` फक्त कोणते खाते डीफॉल्ट आहे ते बदलते (`config.toml` मधील `default`).

कोणत्याही पर्यायाशिवाय टर्मिनलवर `ccam add` चालवल्यास विझार्ड उघडतो. तो हे टप्पे घेतो:

1. खात्याचे नाव, तुम्ही ते दिले नसेल तर;
2. खाते कसे साइन इन करेल: आत्ता Claude लॉगिन, API key किंवा नंतर Claude लॉगिन;
3. कोणत्या गटांमध्ये सामील व्हायचे, `default` निवडलेला ठेवून (अजून एकही गट नसेल तर, `default` गट काय शेअर करेल ते तो दाखवतो);
4. `Create a shortcut command for work? [y/N]`;
5. पुष्टी करण्यासाठी सारांश.

`Esc` एक टप्पा मागे नेतो. `--group`, `--no-group`, `--shortcut` किंवा `--no-login` यांपैकी कोणताही पर्याय, किंवा टर्मिनलशिवाय चालवणे (उदाहरणार्थ, pipe मधून), विझार्ड वगळतो. तेव्हा नाव देणे आवश्यक असते, आणि `--shortcut` ने नाव दिल्याशिवाय शॉर्टकट तयार होत नाही.

`--group` शिवाय नवीन खाते `default` गटात सामील होते. पहिले खाते तो गट तयार करते, जो MCP सर्व्हर, सेटिंग्ज, `CLAUDE.md` आणि extensions शेअर करतो; [गटांसह शेअरिंग](/features/share) पाहा. Windows वर फाइल्स फक्त Developer Mode मध्ये लिंक करता येतात. तो नसताना, पहिले खाते फक्त `mcp` आणि `extensions` सह `default` तयार करते आणि तसे सांगते (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), आणि settings, memory किंवा history शेअर करणाऱ्या गटासह `--group` ccam नाकारते.

`ccam list` slot, नाव, शॉर्टकट, ईमेल, गट आणि शेवटचा वापर दाखवते, आणि डीफॉल्ट खात्याच्या आधी `*` असतो. ईमेल म्हणजे फोल्डर सध्या ज्या ईमेलने लॉग इन आहे तो (त्याच्या `.claude.json` मधून); तो ccam ने नोंदवलेल्या ईमेलपेक्षा वेगळा असेल, तर नोंदवलेला ईमेल त्यानंतर `(was …)` म्हणून येतो:

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

`--json` सह प्रत्येक खाते एक object असते, ज्यात `slot`, `shortcut`, `name`, `email` (टेबलप्रमाणेच सध्याच्या लॉगिनचा), `default`, `dir`, `last_used` (खाते कधीच सुरू केले नसेल तर वगळले जाते), `groups` आणि `share` असतात. `share` प्रत्येक आयटम (`mcp`, `settings`, `memory`, `extensions`, `history`) तो ज्या गटातून येतो त्या गटाशी, किंवा `"local"` शी जोडतो.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `ccam add --no-login` | लॉग इन न करता खाते तयार करा; नंतर `ccam login <name>` ने लॉग इन करा | बंद |
| `ccam add --group <group>` | `default` ऐवजी या गटात सामील व्हा. अनेकदा देता येतो; `default` अस्तित्वात नसेल तर तो तयार होतो | `default` |
| `ccam add --no-group` | कोणत्याही गटात सामील होऊ नका: खाते काहीही शेअर करत नाही | बंद |
| `ccam add --shortcut <command>` | `~/.local/bin` मध्ये ही शॉर्टकट कमांडही तयार करा | काहीही नाही |
| `ccam add --api --key-from <source>` | लॉगिनऐवजी API key वापरणारे खाते; [API key खाती](/features/api) पहा | बंद |
| `ccam list --json` | खाती JSON म्हणून छापा | बंद |
| `ccam remove -y`, `--yes` | न विचारता काढा | बंद |
| `ccam logout -y`, `--yes` | न विचारता लॉग आउट करा | बंद |
| `ccam import --name <name>` | तुम्ही दिलेल्या फोल्डरसाठी खात्याचे नाव; फोल्डरशिवाय नाकारले जाते (`--name needs a folder: ccam import <folder> --name ops`) | फोल्डरच्या नावावरून |
| `ccam import -y`, `--yes` | सापडलेला प्रत्येक फोल्डर न विचारता इम्पोर्ट करा | बंद |
| `--dry-run` | काहीही न बदलणाऱ्या `ccam list` आणि `ccam completion` ला समर्थित. `add`, `login`, `logout`, `use`, `rename`, `remove` आणि `import` तो `` --dry-run is not supported by `ccam add`; nothing was changed `` (कमांडच्या नावासह) देऊन नाकारतात | बंद |

`ccam login` (`--expired` सोडून), `ccam use` आणि `ccam rename` यांना स्वतःचे पर्याय नाहीत. `ccam completion` `bash`, `zsh` किंवा `powershell` घेते; इतर काहीही दिल्यास ते `only bash, zsh and powershell are supported` देऊन अयशस्वी होते.

खात्यांची नावे `a-z`, `0-9`, `-` आणि `_` यांपासून बनलेली 1 ते 32 अक्षरांची असतात, आणि अक्षराने किंवा अंकाने सुरू होतात. नाव फक्त अंकांचे असू शकत नाही, कारण `ccam 2` म्हणजे slot 2. नाव ccam ची कमांडही असू शकत नाही (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), त्यामुळे `ccam <name>` कधीच संदिग्ध होत नाही. ही यादी `internal/config` मधील `config.Reserved` आहे.

## हे कसे काम करते

खाते जोडताना त्याचा फोल्डर तयार होतो, ते आपल्या गटांमध्ये सामील होते, config सेव्ह होते, साइन इन होते आणि शॉर्टकट तयार होतो, याच क्रमाने. तुम्ही `ccam add work --shortcut ccwork` चालवता तेव्हा ccam हे टप्पे घेते (`internal/cli/profiles.go` मधील `addProfile`):

<!-- motion:accounts-add -->

1. काहीही तयार करण्यापूर्वी ccam नाव (`config.ValidateName`) आणि शॉर्टकट (`shim.Check`) तपासते.
2. ccam `0700` mode सह `~/.local/share/ccam/profiles/work/` तयार करते आणि खात्याला सर्वात लहान मोकळा slot क्रमांक देते, त्यामुळे काढलेल्या खात्याचा क्रमांक पुन्हा वापरला जातो (`profile.Add`, `Config.NextSlot`). पहिले खाते डीफॉल्ट होते.
3. खाते त्याच्या गटांमध्ये सामील होते (`--group` किंवा `--no-group` अन्यथा सांगत नसेल तर `default`; पहिले खाते `default` तयार करते), आणि ccam नवीन, रिकामा फोल्डर त्यांच्याशी लिंक करते, उदाहरणार्थ `settings.json` ला `groups/default/settings.json` शी (`applyShare`, `share.Build` आणि मग `Apply` सह). हे अयशस्वी झाले तर नवीन फोल्डर काढला जातो, कारण उरलेला फोल्डर प्रत्येक पुढच्या प्रयत्नाला `folder already exists` ने अडवेल.
4. खाते `fsutil.AtomicWrite` मार्फत `~/.config/ccam/config.toml` मध्ये सेव्ह होते: आधी तात्पुरती फाइल, मग `fsync`, मग `rename`, त्यामुळे crash झाला तरी अर्धवट लिहिलेली config कधीच उरत नाही.
5. ccam खात्याच्या फोल्डरवर सेट केलेल्या `CLAUDE_CONFIG_DIR` सह `claude auth login` सुरू करते, आणि लॉगिन पान आधीच भरलेले असावे म्हणून खात्याच्या शेवटच्या ईमेलसह `--email` देते (`runLogin`). लॉगिन अयशस्वी झाले तरी खाते राहते, आणि ccam त्रुटी छापून नंतर `ccam login work` चालवायला सांगते.
6. ccam `profiles/work/.claude.json` मधून `oauthAccount` वाचते (`profile.ReadIdentity`): ईमेल, खाते आणि organization चे IDs, organization चे नाव आणि plan (`claude_max`, `claude_team` …). ती फाइल अनेक मेगाबाइट्सची असू शकते, म्हणून ccam संपूर्ण फाइल parse करण्याऐवजी फक्त ती एक key शोधते (`TopLevelValue`). ccam हा तपशील `config.toml` मधील `[profile.login]` मध्ये ठेवते.
7. ccam `~/.local/bin/ccwork` ला ccam binary शी लिंक करते, जी सुरू होताना `config.toml` मध्ये `argv[0]` शोधते, आणि ती लिंक `~/.local/state/ccam/shortcuts.json` मध्ये नोंदवते (`shim.Install`). लिंक तयार करता आली नाही तरी खाते जोडले जाते, ccam त्याचे कारण सांगते (`Could not create shortcuts: …`) आणि सारांश `shortcut: ccwork (not created; see above)` दाखवतो.

लॉगिनला काही मिनिटे लागू शकतात, आणि त्या दरम्यान दुसरी ccam process config सेव्ह करू शकते. म्हणून थांबणाऱ्या कमांड्स `config.toml` पुन्हा लोड करतात, तो एक बदल (इथे ईमेल) लागू करतात आणि सेव्ह करतात (`env.update`). याच कारणासाठी `ccam remove` प्रश्न विचारल्यानंतर पुन्हा लोड करते.

`ccam login` आधीपासून असलेल्या खात्यासाठी टप्पे 5 आणि 6 चालवते. Claude Code ईमेल फक्त लॉगिन पानासाठी सूचना म्हणून घेते, त्यामुळे तुम्ही आधीपेक्षा वेगळ्याच Claude खात्यात साइन इन होऊ शकता; ccam नवीन लॉगिनची `[profile.login]` शी तुलना करते आणि तसे सांगते (`ccam: this login is other@company.com, but work was dev@company.com before.`). ते नवीन लॉगिन ठेवते आणि ते कसे उलटवायचे ते सांगते: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

दुसऱ्या Claude खात्यात केलेले लॉगिन खात्याची [मर्यादेची नोंद](/features/usage#the-limit-record) सुद्धा काढून टाकते, कारण त्या मर्यादा दुसऱ्या खात्याच्या होत्या. दुसरे ccam खाते त्याच Claude खात्यात लॉग इन असेल, तर ccam इशारा देते की ती दोन्ही आपल्या वापर मर्यादा शेअर करतात.

`ccam login --expired` कधीच लॉग इन न झालेल्या किंवा ज्यांचे लॉगिन कालबाह्य झाले आहे अशा प्रत्येक खात्यात एकामागून एक लॉग इन करते: त्याचा refresh token कालबाह्य झाला, किंवा सर्व्हरने नकार दिल्यानंतर Claude Code ने लॉगिन काढून टाकले. पहिला browser उघडण्यापूर्वी ते `config.toml` सेव्ह करता येते का ते तपासते.

`ccam logout` खात्यात `claude auth logout` चालवते आणि सेव्ह केलेला ईमेल काढून टाकते; लॉगिनची नोंद राहते, त्यामुळे पुढचे लॉगिन तिच्याशी तपासले जाते. Claude Code लॉगिन रद्द करते आणि त्या खात्यासाठी साठवलेले सर्व काही हटवते, MCP सर्व्हरच्या साइन-इन्ससह. म्हणून पुढील स्थितींमध्ये ccam आधी इशारा देते आणि विचारते (`--yes` हा प्रश्न वगळतो):

- लॉगिनमध्ये MCP साइन-इन्स आहेत;
- खाते चालू आहे (त्याची सत्रेही साइन आउट होतात);
- दुसऱ्या खात्याच्या फोल्डरमध्ये तेच लॉगिन आहे (कॉपी केलेला फोल्डर), जे लॉगआउट सुद्धा रद्द करते.

उदाहरणार्थ:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

दुसऱ्या Claude खात्याने लॉग इन करण्यासाठी फक्त `ccam login <name>` पुरेसे आहे, आणि ते MCP साइन-इन्स ठेवते.

खाती एकापेक्षा जास्त organizations ची असतील तेव्हा `ccam list` एक `Organization` कॉलम जोडते (personal plans `personal` दाखवतात), आणि `ccam list --json` मध्ये प्रत्येक खात्यासाठी `plan` आणि `org` असतात.

### रिमोट मशीनवर साइन इन करणे

SSH वरून `claude auth login` चा browser वाला भाग बिघडू शकतो. Claude Code रिमोट मशीनवरील callback सर्व्हरवर (`http://localhost:<port>/callback`) साइन-इनची वाट पाहते, पण साइन-इन पूर्ण करणारा browser तुमच्या कॉम्प्युटरवर चालतो, जिथे तो पत्ता कुठेच जात नाही आणि browser `ERR_CONNECTION_REFUSED` सारखी त्रुटी दाखवतो. तरीही त्या पत्त्यात एकदाच वापरता येणारा code असतो, त्यामुळे ccam तो पोहोचवू शकते:

<!-- motion:accounts-relay -->

ccam तुमचा टर्मिनल Claude Code च्या हाती देत नाही (`runLogin`). ते Claude Code चे standard input एका pipe ला जोडते आणि तुम्ही टाइप केलेली प्रत्येक ओळ त्यात कॉपी करते (`relayLogin`), त्यामुळे Claude Code च्या `Paste code here` prompt वर paste केलेला code तरीही पोहोचतो. एखादी ओळ callback पत्ता असेल, तर ccam ती फक्त तेव्हाच स्वीकारते जेव्हा host `localhost`, `127.0.0.1` किंवा `[::1]` असेल, port दिलेला असेल, path `/callback` असेल आणि `code` व `state` दोन्ही असतील; `http://` वगळता येते (`callbackURL`). मग ccam तो पत्ता सर्व्हरवर 10 सेकंदांच्या timeout सह आणि redirects न पाळता request करते, जसे browser ने केले असते (`deliverCallback`). ही request वाट पाहणाऱ्या `claude auth login` च्या callback सर्व्हरपर्यंत पोहोचते (आकृतीतील तुटक रेषा). Claude Code ला code मिळतो आणि ते लॉगिन सेव्ह करते, Linux वर खात्याच्या फोल्डरमधील `.credentials.json` मध्ये. paste केलेली ओळ ccam ला कधीही दुसऱ्या host शी संपर्क करायला लावू शकत नाही.

लॉगिन सुरू होण्यापूर्वी ccam हे छापते:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

यशस्वी झाल्यावर ते `→ Passed the sign-in to Claude Code on this machine.` छापते.

लॉगिन टर्मिनलवर input ची वाट पाहत असताना, ccam टर्मिनलचा echo आणि line mode बंद करते आणि तुम्ही टाइप केलेल्या प्रत्येक अक्षरासाठी एक `*` दाखवते, जेणेकरून paste केलेला code किंवा callback पत्ता स्क्रीनवर दिसत नाही. ccam line editing स्वतः करते (`readLine`): `Enter` ओळ पाठवते, `Backspace` एक अक्षर पुसते, आणि रिकाम्या ओळीवर `Ctrl+D` input संपवते. लॉगिन संपल्यावर टर्मिनल पूर्वस्थितीत येतो, आणि `Ctrl+C`, `SIGTERM` किंवा `SIGHUP` वरही येतो, जे तरीही ccam बंद करतात. standard input आणि standard error वर टर्मिनल नसल्यास, ओळी जशा येतात तशा कॉपी होतात.

### नाव बदलणे, काढणे आणि इम्पोर्ट करणे

`ccam rename` आणि `ccam remove` खात्याच्या फोल्डरवर आणि ccam खात्याच्या नावाखाली इतरत्र ठेवत असलेल्या गोष्टींवर काम करतात. ccam ते पाथ `nameState` मध्ये नोंदवते:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- वापराचे caches `~/.local/state/ccam/usage/<name>.json` आणि `<name>.statusline.json`;
- मर्यादेची नोंद `<name>.limit.json`;
- `~/.local/state/ccam/handoff/<name>/` मधील सत्रांच्या नोंदी.

`ccam remove work` ते या क्रमाने हाताळते:

<!-- motion:accounts-remove -->

1. खाते चालू असेल तर `ccam remove` आधी नकार देते, मग विचारते, आणि तुम्ही उत्तर दिल्यानंतर पुन्हा तपासते (यामुळे प्रश्न उघडा असताना सुरू झालेले session पकडले जाते).
2. काहीही बदलण्यापूर्वी ते `config.toml` सेव्ह करता येते का ते तपासते.
3. खाते एक्स्टेंशन्स शेअर करत असेल, तर ccam गटाच्या `installed_plugins.json` आणि `known_marketplaces.json` पुन्हा लिहिते (`RewritePluginPaths`), कारण Claude Code plugins चे absolute पाथ साठवते, जे या खात्याच्या फोल्डरमधून जाऊ शकतात.
4. यानंतर व्यवस्थापित खाते (ज्याचा फोल्डर ccam ने तयार केला, `ccam import` ने नोंदवलेले नव्हे) लॉग आउट केले जाते: त्याच्या फोल्डरमध्ये `claude auth logout` चालते, त्यामुळे trash मध्ये कोणतेही चालू लॉगिन पडून राहत नाही (किंवा macOS वर, फोल्डरच्या नावाच्या Keychain item मध्ये, जो नंतरच्या त्याच नावाच्या खात्याला वारसाने मिळाला असता). API key खात्यासाठी, आणि दुसऱ्या खात्याच्या फोल्डरमध्ये तेच लॉगिन असेल तेव्हा, हे वगळले जाते, कारण लॉगआउट ते लॉगिनही रद्द करेल. अयशस्वी झालेले लॉगआउट कळवले जाते, आणि काढणे पुढे चालू राहते.
5. ccam फोल्डर `~/.local/share/ccam/trash/<name>-<time>/` मध्ये हलवते (`profile.Trash`), खाते `config.toml` मधून निघते, `nameState` पाथ हटवले जातात जेणेकरून नंतरच्या त्याच नावाच्या खात्याला ते वारसाने मिळू नयेत, आणि शॉर्टकट लिंक काढली जाते: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

तरीही `config.toml` सेव्ह करता आले नाही, तर फोल्डर trash मधून परत येतो. काढलेले खाते परत मिळवण्यासाठी फोल्डर trash मधून परत हलवा आणि पुन्हा लॉग इन करा.

`ccam rename` फोल्डरचे नाव बदलते आणि त्याच्या `plugins/` मधील plugin पाथ जुन्या फोल्डरवरून नवीन फोल्डरवर पुन्हा लिहिते. ते `nameState` पाथ नवीन नावावर हलवते आणि खाते डीफॉल्ट असेल तर ते डीफॉल्टच ठेवते. ते सत्रांच्या नोंदीही अद्ययावत करते, त्यामुळे `ccam next --continue` ला त्या तरीही सापडतात. शॉर्टकट कमांडचे नाव तेच राहते. plugin पाथ पुन्हा लिहिणे अयशस्वी झाले, किंवा `config.toml` सेव्ह करता आले नाही, तर फोल्डरचे नाव परत पूर्वीसारखे केले जाते.

macOS वर, आणि Claude Code लॉगिन Credential Manager मध्ये साठवते तेव्हा Windows वर, फोल्डर स्वतः आपले जुने नाव ठेवतो, कारण लॉगिन आणि MCP साइन-इन्स फोल्डरच्या पाथखाली साठवलेले असतात आणि हलवलेला फोल्डर लॉग आउट होईल. मग config नवीन नावाशेजारी फोल्डर `folder` म्हणून नोंदवतो, आणि ccam `its folder stays …: the macOS Keychain keeps the login under the folder's path` असे सांगते.

`ccam import` फोल्डर आहे त्याच जागी, `config.toml` मध्ये `dir` म्हणून नोंदवते. ते त्याला कधीही हलवत किंवा कॉपी करत नाही: कॉपी केलेल्या refresh token मुळे प्रती एकमेकींना लॉग आउट करतील.

फोल्डर न दिल्यास ते न नोंदवलेले `~/.claude` आणि `~/.claude-*` फोल्डर शोधते आणि प्रत्येकाबद्दल विचारते. नावे फोल्डरवरून ठरतात (`~/.claude` म्हणजे `main`, `~/.claude-work` म्हणजे `work`, `~/.claude-2` म्हणजे `acct2`). ccam ला नाव ठरवता आले नाही, किंवा नाव आधीच घेतलेले असेल, तर ते फोल्डर चालवायच्या `--name` कमांडसह दाखवते. सापडलेला एखादा फोल्डर अयशस्वी झाला तर `Skipped <path>: <reason>` असे सांगितले जाते आणि बाकीचे इम्पोर्ट केले जातात.

ccam आधी सर्व प्रश्न विचारते, मग `config.toml` पुन्हा लोड करून फोल्डर नोंदवते, त्यामुळे दरम्यान इतर ccam कमांड्सनी सेव्ह केलेले बदल तसेच राहतात. `✓` ओळी सेव्ह झाल्यानंतर येतात.

`CLAUDE_CONFIG_DIR` नसताना Claude Code `.claude.json` तुमच्या home फोल्डरमध्ये `~/.claude` शेजारी ठेवते, म्हणून तुम्ही `~/.claude` इम्पोर्ट करता आणि त्यात स्वतःची `.claude.json` नसते, तेव्हा ccam `~/.claude/.claude.json` ते `~/.claude.json` अशी लिंक तयार करते. मग साधे `claude` आणि `ccam main` एकच खाते पाहतात.

`~/.claude` मध्ये स्वतःची `.claude.json` असेल आणि शेजारी `~/.claude.json` सुद्धा असेल, तर ccam दोन्ही तशाच ठेवते आणि तसे सांगते. `ccam main` आतली फाइल वापरेल आणि साधे `claude` बाहेरची, त्यामुळे त्यांचे लॉगिन तपशील आणि MCP सर्व्हर वेगळे राहतील. एकच ठेवण्यासाठी आतल्या फाइलच्या जागी लिंक ठेवा ([doctor](/features/doctor) सुद्धा हे कळवतो).

फोल्डरमध्ये नोंदवलेल्या एखाद्या खात्याचेच लॉगिन असेल (कॉपी केलेला फोल्डर), तर ccam तसे सांगते: जे आधी त्याचे नूतनीकरण करेल ते दुसऱ्याला लॉग आउट करते, म्हणून वापरण्यापूर्वी `ccam login <name>` चालवा.

इम्पोर्ट केलेली खाती तुम्ही त्यांना एखाद्या गटात जोडेपर्यंत (`ccam group join <group> <name>`) कोणत्याही गटात नसतात, आणि शॉर्टकट तयार होत नाही.

`ccam completion bash|zsh` एक completion script छापते, जी खात्यांची नावेही पूर्ण करते.

## सुरक्षितता आणि अपवादात्मक स्थिती

- खाते चालू असताना `ccam rename` आणि `ccam remove` (`refuseRunning`) ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)`` देऊन नकार देतात, कारण session फोल्डरमध्ये लिहीत राहते. ही तपासणी Claude Code कसेही सुरू केलेले असले तरी लागू होते: ccam मार्फत, साधे `claude`, IDE, background सत्रे आणि त्यांचा daemon.
- दुसरी ccam process ते लॉगिन नूतनीकरण करत असताना किंवा बदलत असताना `ccam rename`, `ccam remove` आणि `ccam logout` नकार देतात: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- दुसऱ्या ccam process मध्ये गटाचा बदल चालू असताना `ccam rename` आणि `ccam remove` `another ccam process is changing the groups right now; try again when it is done` देऊन नकार देतात; असा बदल अर्ध्यावर थांबला असेल तर ``the last group change stopped halfway; finish it with `ccam doctor --fix` first`` देऊन.
- इम्पोर्ट केलेल्या खात्यावर `ccam remove` फोल्डर आणि त्याचे लॉगिन साध्या `claude` साठी तसेच ठेवते. खाते गटांमध्ये असेल, तर ते आधी `ccam group leave` प्रमाणे ते गट सोडते: प्रत्येक शेअर केलेला आयटम फोल्डरमध्ये गटाच्या आयटमची प्रत बनतो, आणि गटाचे MCP सर्व्हर त्याच्या `.claude.json` मध्ये जातात. मग त्याची नोंदणी रद्द होते: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. ccam ने घेतलेले त्याच्या `.claude.json` चे बॅकअप `<folder>/.ccam-backup/mcp-<time>` मध्ये हलवले जातात; ते हलवता आले नाहीत, तर ते `~/.local/state/ccam/mcp/<name>/backups` मध्ये राहतात आणि संदेश तसे सांगतो. तुमचा `~/.claude` trash मध्ये हलवला जात नाही.
- दोन खात्यांचा फोल्डर एकच असू शकत नाही: एकच फोल्डर दोनदा नमूद करणारी `config.toml` लोड होत नाही.
- ज्या नावाचा फोल्डर `profiles/` मध्ये आधीच आहे ते नाव `ccam add` नाकारते, आणि `ccam import` कडे निर्देश करते. macOS वर, ccam च्या जुन्या आवृत्तीने मागे ठेवलेले नवीन फोल्डरच्या पाथसाठीचे लॉगिन किंवा MCP साइन-इन्स Keychain मध्ये अजून असतील, तर तेही थांबते, आणि ते काढणारी `security delete-generic-password` कमांड सांगते. `config.toml` सेव्ह करता आले नाही, तर नवीन फोल्डर पुन्हा काढला जातो.
- आधीच नोंदवलेला फोल्डर, फोल्डर नसलेला पाथ आणि आधीच घेतलेले नाव `ccam import` नाकारते. फोल्डर `~/.claude` ची प्रत (तेच लॉगिन) असेल, तर `~/.claude` स्वतः नोंदवलेला नसला तरी ccam इशारा देते.
- खाते चालू असेल तर `ccam logout` आधी विचारते; `ccam login` तपासत नाही.
- relay ला Claude Code पर्यंत पोहोचता आले नाही, उदाहरणार्थ लॉगिन संपल्यानंतर, तर ते कारणासह `ccam: could not pass the sign-in to Claude Code:` छापते. HTTP त्रुटीचा अर्थ सहसा code कालबाह्य झाला असा असतो; लॉगिन पुन्हा सुरू करा.

## संबंधित

- [Claude Code सुरू करणे](/features/launch)
- [शॉर्टकट कमांड्स](/features/shortcuts)
- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [एका कॉम्प्युटरवर अनेक Claude Code खाती वापरा](/guides/multiple-claude-code-accounts)
