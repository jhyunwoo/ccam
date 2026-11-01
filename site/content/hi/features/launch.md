---
title: Claude Code शुरू करना
description: ccam <name|number>, ccam run या डैशबोर्ड से किसी भी अकाउंट से Claude Code शुरू करें; आर्ग्युमेंट बिना बदलाव claude को जाते हैं, --dry-run कमांड दिखाता है।
updated: 2026-11-01
order: 1
translation: machine
sourceHash: "11df08608a1e"
---

`ccam work` Claude Code को `work` अकाउंट से शुरू करता है, और नाम के बाद का सब कुछ बिना बदलाव के `claude` को जाता है। इसे वहाँ इस्तेमाल करें जहाँ आप `claude` टाइप करते: टर्मिनल में, किसी स्क्रिप्ट में, या pipe से आए input के साथ।

## उपयोग

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

अकाउंट का नाम या `ccam list` से उसका नंबर दें। ccam नाम के बाद वाले आर्ग्युमेंट नहीं पढ़ता, इसलिए `ccam work --help` Claude Code की help दिखाता है। नाम के बाद `--` लगाना वैकल्पिक है। Standard input, exit code और signals वैसे ही काम करते हैं जैसे `claude` के साथ।

जब पहला आर्ग्युमेंट `-` से शुरू होता है और ccam के अपने flags में से नहीं है, तो ccam आपके लिए अकाउंट चुनता है और सभी आर्ग्युमेंट Claude Code को देता है। यह वह अकाउंट इस्तेमाल करता है जिसे [`ccam folder`](/features/folder) से मौजूदा फ़ोल्डर के लिए सेट किया गया, वरना डिफ़ॉल्ट अकाउंट: वह जिसे `ccam use` से सेट किया गया, या अगर आपने कोई सेट नहीं किया तो सबसे छोटे नंबर वाला अकाउंट।

### डैशबोर्ड

बिना आर्ग्युमेंट के `ccam`, standard input और output के टर्मिनल होने पर [डैशबोर्ड](/features/dashboard) खोलता है, वरना `ccam list` प्रिंट करता है। किसी अकाउंट पर जाएँ (arrow keys से, या उसका नंबर टाइप करके) और उसे शुरू करने के लिए `Enter` दबाएँ; `n` सबसे ज़्यादा गुंजाइश वाला अकाउंट शुरू करता है और `o` `--continue`, `--resume` या दूसरे आर्ग्युमेंट देने का विकल्प देता है। डैशबोर्ड Claude Code शुरू होने से पहले बंद हो जाता है, इसलिए टर्मिनल पहले अपनी सामान्य स्थिति में आ जाता है।

**Team और Enterprise सीटें।** अपने-आप अकाउंट चुनना (`ccam next`, `--fallback`, `ccam pick`, डैशबोर्ड का `n` और, चालू करने के बाद, [`auto_switch`](/features/next#auto-switch)) एक ही Claude Team या Enterprise संगठन की सीटों के बीच भी स्विच करता है, जो Anthropic की शर्तों का उल्लंघन हो सकता है। इस पर भरोसा करने से पहले [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) पढ़ें और अपने संगठन और Anthropic से पता कर लें।

जब shell में `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` या `ANTHROPIC_AUTH_TOKEN` सेट हो, तो Claude Code अकाउंट के अपने लॉगिन के बजाय वह key या token इस्तेमाल कर सकता है, इसलिए ccam शुरू करने से पहले standard error पर आपको चेतावनी देता है (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`)। यह केवल variables पढ़ता है; [API key अकाउंट](/features/api) इसके बजाय उन्हें हटा देता है।

डिफ़ॉल्ट रूप से ccam ख़ुद को Claude Code से बदल देता है, इसलिए और कुछ चलता नहीं रहता। `config.toml` में सबसे ऊपर, पहले `[[group]]`, `[[profile]]` या `[[folder]]` से ऊपर, `relay = true` (या Linux और macOS पर [`auto_switch = true`](/features/next#auto-switch)) होने पर, ccam इसके बजाय interactive sessions को child process के रूप में शुरू करता है और उनके साथ-साथ चलता रहता है। इससे `/ccam:switch`, या उपयोग सीमा पर `auto_switch`, अपने-आप बातचीत को दूसरे अकाउंट पर जारी रख सकता है, session के विकल्पों, permission mode, effort level और agent के साथ; `/ccam:switch` बातचीत को किसी दूसरे संगठन के अकाउंट पर ले जाए, उससे पहले ccam आपसे पूछता है ([विवरण](/features/next#relay))।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--dry-run` | वह environment और कमांड प्रिंट करें जो ccam चलाता, फिर बाहर निकलें | बंद |
| `--fallback` | `-p`/`--print` runs के लिए: जब अकाउंट अपनी उपयोग सीमा तक पहुँचता है, तो वही कमांड सबसे ज़्यादा गुंजाइश वाले लॉग इन अकाउंट के रूप में फिर से चलाएँ ([विवरण](/features/next#run-again-at-the-limit)) | बंद |
| `-v`, `--version` | `ccam <version>` प्रिंट करें (केवल पहले आर्ग्युमेंट के रूप में) | |
| `-h`, `--help` | ccam की help, `ccam --help` या `ccam run --help` के रूप में | |

आर्ग्युमेंट के नियम:

- ccam अपने flags केवल अकाउंट के नाम से पहले पढ़ता है। नाम के बाद हर आर्ग्युमेंट Claude Code का है, `--dry-run` और `--help` भी।
- `--dry-run`, `ccam --dry-run work`, `ccam run --dry-run work` और `ccam --dry-run` (फ़ोल्डर के अकाउंट या डिफ़ॉल्ट के लिए) के रूप में काम करता है। `--dry-run` और `--fallback` किसी भी क्रम में आ सकते हैं।
- `ccam run <name|number> [claude args]`, `ccam <name>` का लंबा रूप है। बिना नाम के `ccam run` फ़ोल्डर का अकाउंट या डिफ़ॉल्ट अकाउंट इस्तेमाल करता है।
- अकाउंट के नाम `add`, `list`, `run` या `usage` जैसे ccam कमांड के नाम नहीं हो सकते, इसलिए `ccam <word>` कभी अस्पष्ट नहीं होता।

Dry run एक पंक्ति प्रिंट करता है जिसे आप shell में, या Windows पर PowerShell में, paste कर सकते हैं। यह उससे शुरू होती है जो launch environment में बदलता है: हमेशा `CLAUDE_CONFIG_DIR` और `CCAM_PROFILE`, हटाए जाने वाले variable के लिए `env -u <name>` (जैसे विरासत में मिला `CCAM_RELAY`), और, [API key अकाउंट](/features/api) के लिए, उसका endpoint, उसका मॉडल और यह कि key कहाँ से आती है (`ANTHROPIC_API_KEY="$WORK_KEY"`), कभी key ख़ुद नहीं:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

PowerShell वाली पंक्ति बाद में variables को वापस रख देती है, जैसे shell में `VAR=x cmd` करता है, और Claude Code को उसकी कमांड लाइन `Start-Process` से देती है, ताकि Windows PowerShell 5.1 और PowerShell 7 उसे एक ही तरह चलाएँ। जब कोई launch किसी variable को ख़ाली मान पर सेट करता है, तो उसके नीचे एक `#` पंक्ति बताती है कि 7.5 से पहले का PowerShell ऐसे variable को इसके बजाय हटा देता है। npm से इंस्टॉल किए गए `claude.cmd` के लिए, पंक्ति वही `cmd.exe` कमांड चलाती है जो launch इस्तेमाल करता (नीचे देखें), और जिस आर्ग्युमेंट को launch वहाँ अस्वीकार करता, उस पर dry run भी विफल होता है।

## यह कैसे काम करता है

हर बार शुरू करने में, कमांड लाइन से `exec` तक, वही चरण होते हैं:

<!-- motion:launch -->

हर बॉक्स अपना चरण चलते समय जलता है, और नीचे की पंक्तियाँ वह environment और आर्ग्युमेंट हैं जिनके साथ Claude Code शुरू होता है।

1. ccam `app.Main` में शुरू होता है, जिसे `cmd/ccam/main.go` बुलाता है। अगर ccam किसी दूसरे नाम से शुरू हुआ, तो वह एक [शॉर्टकट](/features/shortcuts) है और सीधे launcher पर जाता है।
2. ccam cobra के शामिल होने से पहले `os.Args` देखता है (`app.parse`)। केवल अकेला `ccam` (डैशबोर्ड), कोई ccam कमांड नाम (`config.Reserved`), `-h`/`--help` और `-v`/`--version` ही `cli.Execute` पर जाते हैं। बाकी सब कुछ शुरू करना है।
3. ccam `~/.config/ccam/config.toml` पढ़ता है, जो इकलौती फ़ाइल है जिसे यह parse करता है (`launcher.Launch`)।
4. ccam नाम या नंबर से अकाउंट खोजता है, शुरुआत का एक `--` हटाता है, और `claude` executable खोजता है (`launcher.Prepare`, `ResolveClaude`)। नाम न होने पर यह फ़ोल्डर का अकाउंट लेता है, working directory की तुलना `config.toml` की `[[folder]]` entries से strings के रूप में करके, macOS और Windows पर बड़े-छोटे अक्षरों की परवाह किए बिना (`HereProfile`)। PATH की खोज फ़ाइल के modification time और `PATH` के एक hash के साथ `~/.local/state/ccam/claude-path` में cache होती है, इसलिए ccam हर बार शुरू होने पर `PATH` नहीं खोजता। बदली हुई या गायब फ़ाइल, या बदला हुआ `PATH`, नई खोज करवाता है। `config.toml` में दिया `claude_bin` पाथ (`~/` फैलाया जाता है) executable फ़ाइल होना चाहिए; वरना कुछ भी लिखे जाने से पहले launch `claude_bin X is not an executable file. Check claude_bin in config.toml` के साथ रुक जाता है। ccam ऐसे `claude` को अस्वीकार करता है जो ख़ुद ccam हो, inode से तुलना करके, क्योंकि उससे loop बन जाता।
5. ccam अकाउंट फ़ोल्डर `~/.local/share/ccam/profiles/<name>/` की तुलना उसकी साझा सेटिंग्स से करता है, केवल `lstat` और `readlink` बुलाकर, यानी हर आइटम के लिए एक या दो system call (`checkLinks`, `share.Check`)। गायब लिंक तुरंत बना दिया जाता है। जिस भी चीज़ से डेटा खो सकता है, वह शुरू करना रोक देती है।
6. ccam साझा MCP सर्वर जोड़ता है (`applyMCP`; नीचे देखें)।
7. ccam session दर्ज करता है (`session.Register`)। यह अपने PID के नाम वाला एक marker, जिसमें process के शुरू होने का समय होता है, `~/.local/state/ccam/locks/profile-<name>/<pid>` में लिखता है, `XDG_RUNTIME_DIR` चाहे जो कहे, ताकि हर shell वही markers देखे। बातचीत का इतिहास साझा करने वाले अकाउंट के लिए, ccam पहले `history-<group>-<hash>` में कोई चल रहा session खोजता है, जहाँ hash symbolic links resolve किए गए मौजूदा फ़ोल्डर का है, और वहाँ भी रजिस्टर करता है (`historyCheck`)।

   `exec` PID बनाए रखता है, इसलिए marker फिर ख़ुद Claude Code process का नाम बताता है। ccam पकड़े हुए lock के बजाय marker फ़ाइल इस्तेमाल करता है, क्योंकि Claude Code जो processes शुरू करता है वे lock विरासत में पाकर session ख़त्म होने के बाद भी उसे रखे रहते। जो marker लिखा नहीं जा सकता, वह launch नहीं रोकता: ccam कहता है `ccam: could not mark the session as running (…); starting Claude Code anyway`।

   कोई अकाउंट चल रहा है या नहीं, इसकी ccam की दूसरी जाँचें (rename, remove, ग्रुप बदलाव, `doctor --fix`, MCP sync) अकाउंट फ़ोल्डर में Claude Code के अपने records (`sessions/<pid>.json`, `daemon.lock`) भी पढ़ती हैं, इसलिए वे ccam के बिना शुरू हुए Claude Code को भी देखती हैं: सादा `claude`, कोई IDE, background sessions और उनका daemon।
8. ccam टर्मिनल का शीर्षक `claude · <name>` सेट करता है (जब standard error टर्मिनल हो), शुरू करने का समय `fsync` के बिना `~/.local/state/ccam/last-used/<name>` में लिखता है, और `syscall.Exec` बुलाता है (`Plan.Exec`)।

Claude Code, अकाउंट फ़ोल्डर पर सेट `CLAUDE_CONFIG_DIR` और अकाउंट के नाम पर सेट `CCAM_PROFILE` के साथ चलता है। ccam इन variables के विरासत में मिले मान हटा देता है, और विरासत में मिला `CCAM_RELAY` भी। macOS पर `ccam import` से रजिस्टर किए गए `~/.claude` फ़ोल्डर के लिए, यह `CLAUDE_SECURESTORAGE_CONFIG_DIR` को ख़ाली मान पर भी सेट करता है, ताकि अकाउंट और सादा `claude` एक ही कीचेन लॉगिन इस्तेमाल करें; Windows पर यह तब यही करता है जब Claude Code लॉगिन Credential Manager में रखता है।

### exec क्यों

ccam अपने process की जगह `claude` को रख देता है (`syscall.Exec`)। इसका नतीजा:

- ccam memory में नहीं रहता;
- keystrokes, signals और exit code Claude Code के अपने होते हैं;
- आर्ग्युमेंट एक array बने रहते हैं, कभी shell string नहीं बनते।

Windows में `exec` नहीं है। वहाँ ccam `claude.exe` को उसी environment, आर्ग्युमेंट और console के साथ एक child process के रूप में शुरू करता है, उसका इंतज़ार करता है और उसी के exit code के साथ बाहर निकलता है (`Plan.Exec`)। npm से इंस्टॉल किया गया `claude.cmd` `cmd.exe` में अपने पूरे पाथ के साथ चलता है, और हर आर्ग्युमेंट quote किया जाता है ताकि `cmd.exe` उनमें से किसी को अपनी syntax (`&`, `%PATH%`) न समझे; line break वाला आर्ग्युमेंट इस तरह नहीं दिया जा सकता और अस्वीकार होता है, native Claude Code इंस्टॉल करने की सलाह के साथ।

Ctrl+C console से सीधे Claude Code तक पहुँचता है, इसलिए ccam उसे अनदेखा करता है। Session marker में ccam का PID होता है, जो ठीक उतनी देर चलता है जितनी देर Claude Code, और Claude Code के बाहर निकलने पर हटा दिया जाता है। इंतज़ार करते समय ccam और कुछ नहीं करता।

Session के बाद का काम, जैसे `claude mcp add` से जोड़े गए सर्वर पहचानना, अगली बार शुरू होने पर होता है: ccam `.claude.json` के modification time की तुलना `~/.local/state/ccam/mcp/<name>/seen` से करता है और फ़ाइल को केवल तभी parse करता है जब वह बदली हो।

### कमांड लाइन का रास्ता कैसे तय होता है

तीन कमांड लाइनें, `internal/app` से होकर तीन रास्ते:

<!-- motion:launch-route -->

शॉर्टकट का मिलान `argv[0]` पर होता है, `ccam <name>` का पहले आर्ग्युमेंट पर, और दोनों `syscall.Exec` पर ख़त्म होते हैं। केवल `ccam usage` जैसी ccam कमांड ही cobra तक पहुँचती है।

### साझा MCP सर्वर

डिफ़ॉल्ट `flag` strategy के साथ, ccam ग्रुप के MCP सर्वर Claude Code को एक आर्ग्युमेंट के रूप में देता है, आपके आर्ग्युमेंट से पहले:

<!-- motion:launch-mcp -->

ccam आर्ग्युमेंट केवल तभी जोड़ता है जब ग्रुप की `mcp.json` मौजूद हो, और उसे `claude` के ठीक बाद रखता है (`mcp.FlagArgs`, `Plan.InsertArgs`)। तब आपका अपना `--mcp-config` या `--strict-mcp-config` उसके ऊपर फिर भी लागू होता है। यह एक आर्ग्युमेंट है, `--mcp-config=<path>`, क्योंकि Claude Code का `--mcp-config` कई मान लेता है: दो आर्ग्युमेंट के रूप में लिखने पर वह बाद में आने वाले prompt को भी config फ़ाइल समझ लेता।

ccam तय करता है कि flag कहाँ नहीं जा सकता (`claudeVerb`)। यह आर्ग्युमेंट में से Claude Code का subcommand वैसे ही पढ़ता है जैसे Claude Code पढ़ता है, विकल्पों को उनके मानों के साथ छोड़ते हुए। `remote-control` या `rc` को कोई flag नहीं मिलता (नीचे देखें), न ही पहले दी गई background-session verb को (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`)। `ccam work mcp list`, `get` और `remove` को भी कोई नहीं मिलता, क्योंकि Claude Code वहाँ इसे अस्वीकार करता है; इसके बजाय ccam standard error पर बताता है कि ग्रुप के सर्वर अकाउंट के sessions में भी लोड होते हैं:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

बाक़ी हर subcommand और हर session flag बनाए रखता है।

`config.toml` में `mcp_strategy = "sync"` होने पर, ccam इसके बजाय सर्वर अकाउंट की `.claude.json` में लिखता है, और अगर अकाउंट पहले से चल रहा हो तो यह छोड़ देता है। जब पिछली जाँच के बाद से न ग्रुप की `mcp.json` बदली हो न `.claude.json`, तब भी launch यह काम छोड़ देता है। अगर अकाउंट का अपना कोई सर्वर किसी ग्रुप सर्वर के नाम वाला है, तो sync अकाउंट का सर्वर रखता है और एक बार यह बताता है। Sync `.claude.json` से ग़ायब हुई कॉपी भी वापस रखता है (`put back group default's MCP servers …`), और केवल वही कॉपियाँ बदलता या हटाता है जो उसने लिखीं और जिन्हें तब से किसी ने नहीं बदला।

### Remote Control

`ccam work rc` (या `ccam work remote-control`) Claude Code का Remote Control `work` के रूप में शुरू करता है। `claude remote-control` अपनी verb से पहले `--mcp-config` अस्वीकार करता है और उसके बाद कोई नहीं लेता, और उसके sessions MCP सर्वर केवल अकाउंट की `.claude.json` से पढ़ते हैं। इसलिए यहाँ flag strategy अलग तरह से काम करती है:

<!-- motion:launch-rc -->

1. ccam Claude Code की verb के रूप में `remote-control` या `rc` पहचानता है, विकल्पों और उनके मानों के बाद भी (`ccam work --model opus rc`), और कोई flag नहीं जोड़ता (`claudeVerb`, `mcp.FlagArgs`)।
2. ccam ग्रुप की `mcp.json` के modification time की तुलना अकाउंट के `sync-stamp` से करता है, और जाँचता है कि हर ग्रुप सर्वर अभी भी `.claude.json` में है (`mcp.Synced`)। अगर `.claude.json` में पहले से ग्रुप के मौजूदा सर्वर हों, या ग्रुप में कोई सर्वर न हो, तो Remote Control बिना कुछ पूछे शुरू होता है।
3. वरना ccam उन सर्वरों के नाम बताता है जो sessions को नहीं मिलेंगे, और पूछता है:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. `y` पर, ccam उन्हें वैसे ही कॉपी करता है जैसे `sync` strategy करती है (`mcp.Sync`): Claude Code के lock के तहत, बैकअप के बाद, अकाउंट के अपने सर्वर बनाए रखते हुए। `N` पर, या टर्मिनल के बिना, Remote Control उनके बिना शुरू होता है, और ccam `mcp_strategy = "sync"` की ओर इशारा करता है, जो हर बार शुरू करने से पहले उन्हें कॉपी करता है।

सर्वर कॉपी होने तक ccam हर बार Remote Control शुरू करने पर पूछता है। कॉपी के बाद, यह फिर से केवल तब पूछता है जब ग्रुप की `mcp.json` बदले या कोई कॉपी किया गया सर्वर `.claude.json` से ग़ायब हो जाए। बाद में सामान्य रूप से शुरू करने पर `.claude.json` की कॉपी `--mcp-config` पर जीतती है, इसलिए `ccam mcp add`, `remove` और `import` दूसरे अकाउंट्स के पास रखी कॉपियाँ अपडेट करते हैं ([साझा MCP सर्वर](/features/mcp))।

### 10 ms का बजट

शुरू करने के रास्ते का बजट 95th percentile पर 10 ms है; असल में Claude Code शुरू होने से पहले ccam लगभग 4 ms जोड़ता है। शुरू करने का रास्ता कभी cobra का कमांड ट्री नहीं बनाता, Bubble Tea शुरू नहीं करता और कोई नेटवर्क कनेक्शन नहीं खोलता, और अगर package के `init()` functions कुल मिलाकर 5 ms से ज़्यादा लें तो `TestPackageInitIsCheap` विफल होता है। `scripts/perf-check.sh`, `hyperfine` से `ccam --dry-run work`, `ccam work` और एक शॉर्टकट को मापता है और 10 ms से ऊपर विफल होता है, और `internal/launcher` में `BenchmarkLaunchPath` इन चरणों को loop में चलाता है। हर चरण का समय देखने के लिए `CCAM_TRACE=1` सेट करें:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## सुरक्षा और विशेष स्थितियाँ

- अज्ञात नाम या नंबर `no account "<name>"` के साथ रुकता है और `ccam list` की ओर इशारा करता है। अगर कोई अकाउंट ही न हो, तो ccam `ccam add <name>` से एक जोड़ने को कहता है।
- अगर कोई साझा पाथ सामान्य फ़ाइल बन गया हो या ग़लत जगह इशारा करता हो, तो ccam `stopped because of share link problems:` के साथ रुकता है, उनकी सूची दिखाता है, और `ccam doctor --fix` चलाने को कहता है। शुरू करने के रास्ते पर यह डेटा move नहीं करता।
- अगर इस अकाउंट के साथ बातचीत का इतिहास साझा करने वाला (उसी ग्रुप का) कोई दूसरा Claude Code session उसी फ़ोल्डर में चल रहा हो, तो ccam चेतावनी देता है कि इतिहास मिल सकता है और `Run anyway? [y/N]` पूछता है। फ़ोल्डर का मिलान symbolic links resolve करके होता है। टर्मिनल के बिना (`/dev/null` से standard input को टर्मिनल नहीं माना जाता), या background के किसी job से, यह चेतावनी दिखाकर आगे बढ़ता है। नया `-p` run अपनी अलग बातचीत शुरू करता है और उसे कोई चेतावनी नहीं मिलती; `--version`, `--help` और ऐसे subcommands जो कोई बातचीत नहीं खोलते, जैसे `mcp`, `auth` या `plugin`, बिल्कुल नहीं जाँचे जाते।
- `--dry-run`, MCP sync, नए सर्वर की सूचना, `mcp list` वाला नोट, इतिहास की जाँच, session marker, शीर्षक और आख़िरी उपयोग का समय छोड़ देता है, और exec नहीं करता। यह फिर भी साझा लिंक जाँचता है, और सामान्य शुरुआत की तरह गायब लिंक बनाता है।
- `ccam work rc`, जब `work` पहले से चल रहा हो, `.claude.json` को नहीं छूता, यह बताता है, और Remote Control को ग्रुप के सर्वरों के बिना शुरू करता है; उन्हें कॉपी करने के लिए session ख़त्म होने के बाद इसे फिर से चलाएँ। `--dry-run` के साथ यह संदेश प्रिंट करता है, लेकिन न पूछता है न कॉपी करता है।
- `mcp_strategy = "sync"` के साथ, विफल sync `ccam: MCP sync failed (launching anyway): …` प्रिंट करता है और Claude Code उन्हीं सर्वरों के साथ शुरू होता है जो उसके पास पहले से थे।

## संबंधित

- [शॉर्टकट कमांड](/features/shortcuts)
- [ccam के साथ Remote Control का उपयोग](/guides/remote-control)
- [अकाउंट प्रबंधित करना](/features/accounts)
- [कमांड संदर्भ](/docs/commands)
- [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing)
- [CLAUDE_CONFIG_DIR क्या है](/guides/claude-config-dir)
