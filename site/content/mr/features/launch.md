---
title: Claude Code सुरू करणे
description: ccam <name|number>, ccam run किंवा डॅशबोर्डने कोणत्याही खात्याने Claude Code सुरू करा; आर्ग्युमेंट्स न बदलता claude ला जातात, आणि --dry-run कमांड दाखवते.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: 11df08608a1e
---

`ccam work` `work` खात्याने Claude Code सुरू करते, आणि नावानंतरचे सर्व काही न बदलता `claude` ला जाते. जिथे तुम्ही `claude` टाइप केले असते तिथे ते वापरा: टर्मिनलमध्ये, script मध्ये, किंवा pipe ने input देऊन.

## वापर

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

खात्याचे नाव किंवा `ccam list` मधील त्याचा क्रमांक द्या. नावानंतरचे आर्ग्युमेंट्स ccam वाचत नाही, त्यामुळे `ccam work --help` Claude Code ची मदत दाखवते. नावानंतर `--` लावणे ऐच्छिक आहे. standard input, exit code आणि signals `claude` प्रमाणेच वागतात.

पहिला आर्ग्युमेंट `-` ने सुरू होत असेल आणि तो ccam च्या स्वतःच्या flags पैकी नसेल, तर ccam तुमच्यासाठी खाते निवडते आणि सर्व आर्ग्युमेंट्स Claude Code ला देते. ते [`ccam folder`](/features/folder) ने सध्याच्या फोल्डरसाठी सेट केलेले खाते वापरते, नाहीतर डीफॉल्ट खाते: `ccam use` ने सेट केलेले खाते, किंवा तुम्ही ते सेट केले नसेल तर सर्वात लहान क्रमांकाचे खाते.

### डॅशबोर्ड

standard input आणि output टर्मिनल असताना आर्ग्युमेंट्सशिवाय `ccam` [डॅशबोर्ड](/features/dashboard) उघडते, आणि इतर वेळी `ccam list` छापते. एखाद्या खात्यावर जा (बाण की ने, किंवा त्याचा क्रमांक टाइप करून) आणि ते सुरू करण्यासाठी `Enter` दाबा; `n` सर्वात जास्त जागा असलेले खाते सुरू करते आणि `o` `--continue`, `--resume` किंवा इतर आर्ग्युमेंट्स देऊ करते. Claude Code सुरू होण्यापूर्वी डॅशबोर्ड बंद होतो, त्यामुळे टर्मिनल आधी पूर्ववत होतो.

**Team आणि Enterprise सीट्स.** आपोआप होणारी खाते निवड (`ccam next`, `--fallback`, `ccam pick`, डॅशबोर्डमधील `n` आणि, तुम्ही चालू केल्यावर, [`auto_switch`](/features/next#auto-switch)) एकाच Claude Team किंवा Enterprise organization मधील सीट्सदरम्यानही खाते बदलते, आणि यामुळे Anthropic च्या अटींचे उल्लंघन होऊ शकते. त्यावर अवलंबून राहण्यापूर्वी [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) वाचा आणि तुमच्या संस्थेकडे आणि Anthropic कडे चौकशी करा.

shell मध्ये `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` किंवा `ANTHROPIC_AUTH_TOKEN` सेट असेल, तर Claude Code खात्याच्या स्वतःच्या लॉगिनऐवजी ती key किंवा token वापरू शकते, म्हणून ccam सुरू होण्यापूर्वी standard error वर तुम्हाला इशारा देते (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). ते फक्त variables वाचते; [API key खाते](/features/api) त्याऐवजी ते काढून टाकते.

डीफॉल्टनुसार ccam स्वतःची जागा Claude Code ला देते, त्यामुळे दुसरे काहीही चालू राहत नाही. `config.toml` च्या सुरुवातीला, पहिल्या `[[group]]`, `[[profile]]` किंवा `[[folder]]` च्या वर, `relay = true` (किंवा Linux आणि macOS वर [`auto_switch = true`](/features/next#auto-switch)) असेल, तर ccam त्याऐवजी interactive सत्रे child process म्हणून सुरू करते आणि त्यांच्यासोबत चालू राहते. त्यामुळे `/ccam:switch`, किंवा वापर मर्यादेवर `auto_switch`, सत्राचे पर्याय, permission mode, effort level आणि agent यांसह संभाषण आपोआप दुसऱ्या खात्यावर पुढे चालू ठेवू शकते; `/ccam:switch` ते दुसऱ्या organization च्या खात्यावर हलवण्यापूर्वी ccam तुम्हाला विचारते ([तपशील](/features/next#relay)).

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--dry-run` | ccam चालवेल ते environment आणि कमांड छापा, मग बाहेर पडा | बंद |
| `--fallback` | `-p`/`--print` runs साठी: खाते आपल्या वापर मर्यादेला पोहोचले की तीच कमांड सर्वात जास्त जागा असलेल्या लॉगिन केलेल्या खात्याने पुन्हा चालवा ([तपशील](/features/next#run-again-at-the-limit)) | बंद |
| `-v`, `--version` | `ccam <version>` छापा (फक्त पहिला आर्ग्युमेंट म्हणून) | |
| `-h`, `--help` | ccam ची मदत, `ccam --help` किंवा `ccam run --help` म्हणून | |

आर्ग्युमेंट्सचे नियम:

- ccam आपले स्वतःचे flags फक्त खात्याच्या नावाआधी वाचते. नावानंतरचा प्रत्येक आर्ग्युमेंट Claude Code चा असतो, `--dry-run` आणि `--help` सुद्धा.
- `--dry-run` `ccam --dry-run work`, `ccam run --dry-run work` आणि `ccam --dry-run` (फोल्डरच्या खात्यासाठी किंवा डीफॉल्टसाठी) असे चालते. `--dry-run` आणि `--fallback` कोणत्याही क्रमाने येऊ शकतात.
- `ccam run <name|number> [claude args]` हे `ccam <name>` चे लांब रूप आहे. नावाशिवाय `ccam run` फोल्डरचे खाते किंवा डीफॉल्ट खाते वापरते.
- खात्यांची नावे `add`, `list`, `run` किंवा `usage` सारखी ccam कमांड्सची नावे असू शकत नाहीत, त्यामुळे `ccam <word>` कधीच संदिग्ध होत नाही.

dry run एक ओळ छापतो, जी तुम्ही shell मध्ये, किंवा Windows वर PowerShell मध्ये, paste करू शकता. ती launch environment मध्ये काय बदलते त्याने सुरू होते: `CLAUDE_CONFIG_DIR` आणि `CCAM_PROFILE` नेहमी, launch काढून टाकत असलेल्या variable साठी `env -u <name>` (जसे वारशाने मिळालेले `CCAM_RELAY`), आणि [API key खात्यासाठी](/features/api) त्याचा endpoint, त्याचे model आणि key कुठून येते ते (`ANTHROPIC_API_KEY="$WORK_KEY"`), key स्वतः कधीच नाही:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

PowerShell ओळ नंतर variables पूर्वीसारखे करते, जसे shell मध्ये `VAR=x cmd` करते, आणि Claude Code ला त्याची कमांड लाइन `Start-Process` मार्फत देते, त्यामुळे Windows PowerShell 5.1 आणि PowerShell 7 ती एकाच प्रकारे चालवतात. launch एखादा variable रिकाम्या मूल्यावर सेट करत असेल, तर तिच्या खालची `#` ओळ सांगते की 7.5 पूर्वीचे PowerShell असा variable त्याऐवजी काढून टाकते. npm ने इन्स्टॉल केलेल्या `claude.cmd` साठी, ही ओळ launch वापरेल तीच `cmd.exe` कमांड चालवते (खाली पाहा), आणि launch तिथे नाकारेल असा आर्ग्युमेंट dry run सुद्धा अयशस्वी करतो.

## हे कसे काम करते

प्रत्येक वेळी सुरू करताना, कमांड लाइनपासून `exec` पर्यंत तेच टप्पे होतात:

<!-- motion:launch -->

प्रत्येक टप्पा चालताना त्याचा box उजळतो, आणि खालच्या ओळी म्हणजे Claude Code ज्या environment आणि आर्ग्युमेंट्ससह सुरू होते ते.

1. ccam `app.Main` मध्ये सुरू होते, जे `cmd/ccam/main.go` बोलावते. ccam दुसऱ्या नावाने सुरू झाले असेल, तर तो [शॉर्टकट](/features/shortcuts) आहे आणि थेट launcher कडे जातो.
2. cobra सहभागी होण्यापूर्वी ccam `os.Args` पाहते (`app.parse`). फक्त एकटे `ccam` (डॅशबोर्ड), ccam कमांडचे नाव (`config.Reserved`), `-h`/`--help` आणि `-v`/`--version` `cli.Execute` कडे जातात. बाकी सर्व काही म्हणजे Claude Code सुरू करणे.
3. ccam `~/.config/ccam/config.toml` वाचते, ही एकमेव फाइल जी ते parse करते (`launcher.Launch`).
4. ccam नाव किंवा क्रमांकावरून खाते शोधते, सुरुवातीचा एक `--` काढते, आणि `claude` executable शोधते (`launcher.Prepare`, `ResolveClaude`). नाव नसल्यास ते फोल्डरचे खाते घेते, working directory ची `config.toml` मधील `[[folder]]` नोंदींशी strings म्हणून तुलना करून, macOS आणि Windows वर लहान-मोठ्या अक्षरांचा विचार न करता (`HereProfile`). PATH मधील शोध फाइलच्या modification time आणि `PATH` च्या hash सह `~/.local/state/ccam/claude-path` मध्ये cache होतो, त्यामुळे ccam प्रत्येक वेळी `PATH` शोधत नाही. फाइल बदलली किंवा नसली, किंवा `PATH` बदलला, तर नवीन शोध होतो. `config.toml` मधील `claude_bin` पाथ (`~/` विस्तारित केला जातो) executable फाइल असला पाहिजे; नाहीतर काहीही लिहिले जाण्यापूर्वी launch `claude_bin X is not an executable file. Check claude_bin in config.toml` ने थांबतो. जो `claude` स्वतः ccam आहे, तो ccam inode ने तुलना करून नाकारते, कारण त्याने loop होईल.
5. ccam खात्याचा फोल्डर `~/.local/share/ccam/profiles/<name>/` त्याच्या शेअरिंग सेटिंग्जशी तुलना करते, फक्त `lstat` आणि `readlink` बोलावून, प्रत्येक आयटमसाठी एक किंवा दोन system calls (`checkLinks`, `share.Check`). नसलेली लिंक तिथल्या तिथे तयार होते. डेटा गमावू शकेल अशी कोणतीही गोष्ट सुरू करणे थांबवते.
6. ccam शेअर केलेले MCP सर्व्हर जोडते (`applyMCP`; खाली पाहा).
7. ccam session नोंदवते (`session.Register`). ते आपल्या PID च्या नावाचा, process सुरू झाल्याची वेळ असलेला marker `~/.local/state/ccam/locks/profile-<name>/<pid>` मध्ये लिहिते, `XDG_RUNTIME_DIR` काहीही सांगो, त्यामुळे प्रत्येक shell ला तेच markers दिसतात. संभाषण इतिहास शेअर करणाऱ्या खात्यासाठी, ccam आधी `history-<group>-<hash>` मध्ये चालू session शोधते, जिथे hash symbolic links सोडवलेल्या सध्याच्या फोल्डरचा असतो, आणि तिथेही नोंदवते (`historyCheck`).

   `exec` PID तोच ठेवते, त्यामुळे marker नंतर Claude Code process चेच नाव सांगतो. ccam धरलेल्या lock ऐवजी marker फाइल वापरते, कारण Claude Code सुरू करत असलेल्या processes ना lock वारसाने मिळेल आणि session संपल्यानंतरही ते राहील. marker लिहिता आला नाही तरी launch थांबत नाही: ccam `ccam: could not mark the session as running (…); starting Claude Code anyway` सांगते.

   खाते चालू आहे का याच्या ccam च्या इतर तपासण्या (नाव बदलणे, काढणे, गटांचे बदल, `doctor --fix`, MCP sync) खात्याच्या फोल्डरमधील Claude Code च्या स्वतःच्या नोंदीही वाचतात (`sessions/<pid>.json`, `daemon.lock`), त्यामुळे त्यांना ccam शिवाय सुरू केलेले Claude Code सुद्धा दिसते: साधे `claude`, IDE, background सत्रे आणि त्यांचा daemon.
8. ccam टर्मिनलचे title `claude · <name>` करते (standard error टर्मिनल असताना), सुरू केल्याची वेळ `fsync` शिवाय `~/.local/state/ccam/last-used/<name>` मध्ये लिहिते, आणि `syscall.Exec` बोलावते (`Plan.Exec`).

Claude Code खात्याच्या फोल्डरवर सेट केलेल्या `CLAUDE_CONFIG_DIR` सह आणि खात्याच्या नावावर सेट केलेल्या `CCAM_PROFILE` सह चालते. या variables ची वारसाने मिळालेली कोणतीही मूल्ये ccam काढून टाकते, आणि वारसाने मिळालेले `CCAM_RELAY` सुद्धा. macOS वर `ccam import` ने नोंदवलेल्या `~/.claude` फोल्डरसाठी ते `CLAUDE_SECURESTORAGE_CONFIG_DIR` रिकाम्या मूल्यावर सेट करते, जेणेकरून ते खाते आणि साधे `claude` एकच Keychain लॉगिन वापरतील; Claude Code लॉगिन Credential Manager मध्ये ठेवत असेल तेव्हा Windows वरही ते असेच करते.

### exec का

ccam स्वतःच्या process च्या जागी `claude` ठेवते (`syscall.Exec`). त्यामुळे:

- ccam memory मध्ये राहत नाही;
- keystrokes, signals आणि exit code Claude Code चेच असतात;
- आर्ग्युमेंट्स array राहतात, कधीच shell string बनत नाहीत.

Windows मध्ये `exec` नाही. तिथे ccam `claude.exe` ला त्याच environment, आर्ग्युमेंट्स आणि console सह child process म्हणून सुरू करते, त्याची वाट पाहते आणि त्याच्याच exit code ने बाहेर पडते (`Plan.Exec`). npm ने इन्स्टॉल केलेले `claude.cmd` `cmd.exe` मध्ये चालते, त्याचा पूर्ण पाथ आणि प्रत्येक आर्ग्युमेंट quote करून, जेणेकरून `cmd.exe` त्यांपैकी काहीही स्वतःची syntax म्हणून वाचत नाही (`&`, `%PATH%`); line break असलेला आर्ग्युमेंट अशा प्रकारे देता येत नाही आणि तो नाकारला जातो, native Claude Code इन्स्टॉल करण्याच्या सल्ल्यासह.

Ctrl+C console मधून थेट Claude Code पर्यंत पोहोचतो, त्यामुळे ccam त्याकडे दुर्लक्ष करते. session marker मध्ये ccam चा PID असतो, जो Claude Code इतकाच काळ जिवंत राहतो, आणि Claude Code बाहेर पडल्यावर marker काढला जातो. वाट पाहत असताना ccam दुसरे काहीही करत नाही.

session नंतरचे काम, जसे `claude mcp add` ने जोडलेले सर्व्हर लक्षात घेणे, पुढच्या वेळी सुरू करताना होते: ccam `.claude.json` च्या modification time ची `~/.local/state/ccam/mcp/<name>/seen` शी तुलना करते आणि फाइल बदलली असेल तरच ती parse करते.

### कमांड लाइन कुठे जाते

तीन कमांड लाइन्स, `internal/app` मधून तीन मार्ग:

<!-- motion:launch-route -->

शॉर्टकट `argv[0]` वरून ओळखला जातो, `ccam <name>` पहिल्या आर्ग्युमेंटवरून, आणि दोन्ही `syscall.Exec` मध्ये संपतात. फक्त `ccam usage` सारखी ccam कमांड cobra पर्यंत पोहोचते.

### शेअर केलेले MCP सर्व्हर

डीफॉल्ट `flag` strategy सह, ccam गटाचे MCP सर्व्हर Claude Code ला एक आर्ग्युमेंट म्हणून, तुमच्या आर्ग्युमेंट्सच्या आधी देते:

<!-- motion:launch-mcp -->

गटाची `mcp.json` असेल तरच ccam आर्ग्युमेंट जोडते, आणि तो `claude` च्या लगेच नंतर ठेवते (`mcp.FlagArgs`, `Plan.InsertArgs`). तुमचे स्वतःचे `--mcp-config` किंवा `--strict-mcp-config` त्यावर लागू होत राहतात. तो `--mcp-config=<path>` हा एकच आर्ग्युमेंट आहे, कारण Claude Code चा `--mcp-config` अनेक मूल्ये घेतो: दोन आर्ग्युमेंट्स म्हणून लिहिल्यास तो पुढचा prompt सुद्धा config फाइल म्हणून वाचेल.

flag कुठे जाऊ शकत नाही ते ccam ठरवते (`claudeVerb`). ते Claude Code प्रमाणेच आर्ग्युमेंट्समधून Claude Code ची subcommand वाचते, पर्याय त्यांच्या मूल्यांसह वगळून. `remote-control` किंवा `rc` ला flag मिळत नाही (खाली), तसेच पहिला आर्ग्युमेंट म्हणून दिलेल्या background-session verb ला सुद्धा नाही (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`). `ccam work mcp list`, `get` आणि `remove` यांनाही flag मिळत नाही, कारण Claude Code तिथे तो नाकारते; त्याऐवजी ccam standard error वर सांगते की गटाचे सर्व्हर खात्याच्या सत्रांमध्येही लोड होतात:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

इतर प्रत्येक subcommand आणि प्रत्येक सत्र flag ठेवतात.

`config.toml` मध्ये `mcp_strategy = "sync"` असल्यास, ccam त्याऐवजी सर्व्हर खात्याच्या `.claude.json` मध्ये लिहिते, आणि खाते आधीच चालू असेल तर ते वगळते. मागच्या तपासणीनंतर गटाची `mcp.json` किंवा `.claude.json` दोन्हीपैकी काहीही बदलले नसेल, तरीही launch हे काम वगळतो. खात्याकडे गटातील सर्व्हरच्याच नावाचा स्वतःचा सर्व्हर असेल, तर sync खात्याचा सर्व्हर ठेवते आणि एकदा तसे सांगते. sync `.claude.json` मधून नाहीशी झालेली प्रतही परत ठेवते (`put back group default's MCP servers …`), आणि फक्त स्वतः लिहिलेल्या, आणि त्यानंतर कोणीही न बदललेल्या, प्रती बदलते किंवा काढते.

### Remote Control

`ccam work rc` (किंवा `ccam work remote-control`) Claude Code चे Remote Control `work` म्हणून सुरू करते. `claude remote-control` त्याच्या verb च्या आधी `--mcp-config` नाकारतो आणि नंतर एकही घेत नाही, आणि त्याचे sessions MCP सर्व्हर फक्त खात्याच्या `.claude.json` मधून वाचतात. म्हणून इथे flag strategy वेगळ्या प्रकारे काम करते:

<!-- motion:launch-rc -->

1. ccam `remote-control` किंवा `rc` हा Claude Code चा verb ओळखते, पर्याय आणि त्यांच्या मूल्यांनंतर आला तरीही (`ccam work --model opus rc`), आणि कोणताही flag जोडत नाही (`claudeVerb`, `mcp.FlagArgs`).
2. ccam गटाच्या `mcp.json` ची modification time खात्याच्या `sync-stamp` शी तुलना करते, आणि गटाचा प्रत्येक सर्व्हर अजून `.claude.json` मध्ये आहे का ते तपासते (`mcp.Synced`). `.claude.json` मध्ये आधीच गटाचे सध्याचे सर्व्हर असतील, किंवा गटात एकही सर्व्हर नसेल, तर Remote Control प्रश्न न विचारता सुरू होते.
3. नाहीतर sessions ना जे सर्व्हर मिळणार नाहीत त्यांची नावे ccam सांगते आणि विचारते:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. `y` वर, ccam ते `sync` strategy प्रमाणेच कॉपी करते (`mcp.Sync`): Claude Code च्या lock खाली, बॅकअप घेतल्यानंतर, खात्याचे स्वतःचे सर्व्हर तसेच ठेवून. `N` वर, किंवा टर्मिनल नसल्यास, Remote Control त्यांच्याशिवाय सुरू होते, आणि ccam `mcp_strategy = "sync"` सुचवते, जे प्रत्येक वेळी सुरू करण्यापूर्वी ते कॉपी करते.

सर्व्हर कॉपी होईपर्यंत ccam प्रत्येक वेळी Remote Control सुरू करताना विचारते. कॉपीनंतर ते फक्त गटाची `mcp.json` बदलल्यावर किंवा कॉपी केलेला सर्व्हर `.claude.json` मधून नाहीसा झाल्यावरच पुन्हा विचारते. नंतर सामान्यपणे सुरू करताना `.claude.json` मधील कॉपी `--mcp-config` वर जिंकते, म्हणून `ccam mcp add`, `remove` आणि `import` इतर खात्यांकडील कॉपी अपडेट करतात ([शेअर केलेले MCP सर्व्हर](/features/mcp)).

### 10 ms ची मर्यादा

सुरू करण्याच्या मार्गासाठी 95 व्या percentile ला 10 ms ची मर्यादा आहे; प्रत्यक्षात Claude Code सुरू होण्यापूर्वी ccam सुमारे 4 ms जोडते. सुरू करण्याचा मार्ग कधीच cobra कमांड tree तयार करत नाही, Bubble Tea सुरू करत नाही किंवा नेटवर्क कनेक्शन उघडत नाही, आणि package च्या `init()` functions ना एकूण 5 ms पेक्षा जास्त लागल्यास `TestPackageInitIsCheap` अयशस्वी होते. `scripts/perf-check.sh` `hyperfine` ने `ccam --dry-run work`, `ccam work` आणि एक शॉर्टकट मोजते आणि 10 ms पेक्षा जास्त झाल्यास अयशस्वी होते, आणि `internal/launcher` मधील `BenchmarkLaunchPath` हे टप्पे loop मध्ये चालवते. प्रत्येक टप्प्याचा वेळ पाहण्यासाठी `CCAM_TRACE=1` सेट करा:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## सुरक्षितता आणि अपवादात्मक स्थिती

- अज्ञात नाव किंवा क्रमांक `no account "<name>"` ने थांबतो आणि `ccam list` कडे निर्देश करतो. एकही खाते नसल्यास, ccam `ccam add <name>` ने एक जोडायला सांगते.
- एखादा शेअर केलेला पाथ साधी फाइल बनला असेल किंवा चुकीच्या जागी निर्देश करत असेल, तर ccam `stopped because of share link problems:` ने थांबते, त्यांची यादी देते, आणि `ccam doctor --fix` चालवायला सांगते. सुरू करण्याच्या मार्गावर ते डेटा हलवत नाही.
- या खात्याशी संभाषण इतिहास शेअर करणारे (त्याच गटातील) दुसरे Claude Code session त्याच फोल्डरमध्ये चालू असेल, तर ccam इतिहास मिसळू शकतो असा इशारा देते आणि `Run anyway? [y/N]` विचारते. फोल्डरची जुळणी symbolic links सोडवून होते. टर्मिनल नसताना (`/dev/null` मधून येणारे standard input टर्मिनल नसल्यासारखेच मानले जाते), किंवा background मधील job मधून, ते इशारा दाखवते आणि पुढे जाते. नवीन `-p` run स्वतःचे संभाषण सुरू करते आणि त्याला इशारा मिळत नाही; `--version`, `--help` आणि कोणतेही संभाषण न उघडणाऱ्या subcommands, जसे `mcp`, `auth` किंवा `plugin`, मुळीच तपासल्या जात नाहीत.
- `--dry-run` MCP sync, नवीन सर्व्हरची सूचना, `mcp list` ची टीप, इतिहासाची तपासणी, session marker, title आणि शेवटच्या वापराची वेळ वगळते, आणि exec करत नाही. तरीही ते शेअर लिंक्स तपासते, आणि सामान्यपणे सुरू करताना होईल तशी नसलेली लिंक तयार करते.
- `ccam work rc` हे `work` आधीच चालू असताना `.claude.json` ला हात लावत नाही, तसे सांगते, आणि गटाच्या सर्व्हरशिवाय Remote Control सुरू करते; ते कॉपी करण्यासाठी session संपल्यानंतर ते पुन्हा चालवा. `--dry-run` सह ते संदेश छापते पण विचारत नाही आणि कॉपीही करत नाही.
- `mcp_strategy = "sync"` सह sync अयशस्वी झाल्यास `ccam: MCP sync failed (launching anyway): …` छापले जाते आणि Claude Code त्याच्याकडे आधीपासून असलेल्या सर्व्हरसह सुरू होते.

## संबंधित

- [शॉर्टकट कमांड्स](/features/shortcuts)
- [ccam सोबत Remote Control वापरणे](/guides/remote-control)
- [खाती व्यवस्थापित करणे](/features/accounts)
- [कमांड संदर्भ](/docs/commands)
- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [CLAUDE_CONFIG_DIR समजून घ्या](/guides/claude-config-dir)
