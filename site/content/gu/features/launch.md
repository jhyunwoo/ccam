---
title: Claude Code શરૂ કરવું
description: ccam <name|number>, ccam run કે ડેશબોર્ડ વડે કોઈ પણ એકાઉન્ટ તરીકે Claude Code શરૂ કરો; આર્ગ્યુમેન્ટ્સ બદલાયા વગર claude ને જાય છે, --dry-run કમાન્ડ બતાવે છે.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: 11df08608a1e
---

`ccam work` Claude Code ને `work` એકાઉન્ટ તરીકે શરૂ કરે છે, અને નામ પછીનું બધું કોઈ ફેરફાર વગર `claude` ને જાય છે. જ્યાં તમે `claude` લખતા હો ત્યાં તેને વાપરો: ટર્મિનલમાં, સ્ક્રિપ્ટમાં, કે pipe થી input આપીને.

## ઉપયોગ

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

એકાઉન્ટનું નામ કે `ccam list` માંનો તેનો નંબર આપો. ccam નામ પછીના આર્ગ્યુમેન્ટ્સ વાંચતું નથી, તેથી `ccam work --help` Claude Code ની મદદ બતાવે છે. નામ પછી `--` મૂકવું વૈકલ્પિક છે. standard input, exit code અને સિગ્નલ `claude` જેવું જ વર્તન કરે છે.

જ્યારે પહેલો આર્ગ્યુમેન્ટ `-` થી શરૂ થાય અને ccam ના પોતાના flag માંનો ન હોય, ત્યારે ccam તમારા માટે એકાઉન્ટ પસંદ કરે છે અને બધા આર્ગ્યુમેન્ટ્સ Claude Code ને આપે છે. તે [`ccam folder`](/features/folder) વડે હાલના ફોલ્ડર માટે સેટ કરેલું એકાઉન્ટ વાપરે છે, નહીં તો ડિફૉલ્ટ એકાઉન્ટ: `ccam use` વડે સેટ કરેલું, અથવા તમે સેટ ન કર્યું હોય તો સૌથી નાના નંબરવાળું એકાઉન્ટ.

### ડેશબોર્ડ

standard input અને output ટર્મિનલ હોય ત્યારે આર્ગ્યુમેન્ટ વગરનું `ccam` [ડેશબોર્ડ](/features/dashboard) ખોલે છે, નહીં તો `ccam list` છાપે છે. કોઈ એકાઉન્ટ પર જાઓ (એરો કી વડે, કે તેનો નંબર લખીને) અને તેને શરૂ કરવા `Enter` દબાવો; `n` સૌથી વધુ જગ્યાવાળું એકાઉન્ટ શરૂ કરે છે અને `o` `--continue`, `--resume` કે બીજા આર્ગ્યુમેન્ટ્સ આપવાની તક આપે છે. Claude Code શરૂ થાય તે પહેલાં ડેશબોર્ડ બંધ થાય છે, તેથી ટર્મિનલ પહેલાં પાછું સામાન્ય થઈ જાય છે.

**Team અને Enterprise સીટ.** આપમેળે થતી એકાઉન્ટ પસંદગી (`ccam next`, `--fallback`, `ccam pick`, ડેશબોર્ડની `n` કી અને, તમે તેને ચાલુ કરો પછી, [`auto_switch`](/features/next#auto-switch)) એક જ Claude Team કે Enterprise સંસ્થાની સીટો વચ્ચે પણ સ્વિચ કરે છે, જેનાથી Anthropic ની શરતોનું ઉલ્લંઘન થઈ શકે છે. તેના પર આધાર રાખતાં પહેલાં [વારંવાર પૂછાતા પ્રશ્નો](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) વાંચો અને તમારી સંસ્થા અને Anthropic પાસેથી ખાતરી કરી લો.

જ્યારે શેલમાં `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` કે `ANTHROPIC_AUTH_TOKEN` સેટ હોય, ત્યારે Claude Code એકાઉન્ટના પોતાના લૉગિનને બદલે એ key કે token વાપરી શકે છે, તેથી ccam શરૂ કરતાં પહેલાં standard error પર તમને ચેતવે છે (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). તે ફક્ત variables વાંચે છે; [API key એકાઉન્ટ](/features/api) તેના બદલે તેમને સાફ કરે છે.

ડિફૉલ્ટ રીતે, ccam પોતાની જગ્યા Claude Code ને આપી દે છે, તેથી બીજું કંઈ ચાલતું રહેતું નથી. `config.toml` માં સૌથી ઉપર, પહેલા `[[group]]`, `[[profile]]` કે `[[folder]]` ની ઉપર, `relay = true` (અથવા Linux અને macOS પર [`auto_switch = true`](/features/next#auto-switch)) સાથે, ccam તેના બદલે interactive સેશનને child process તરીકે શરૂ કરે છે અને તેમની સાથે ચાલતું રહે છે. આથી `/ccam:switch`, અથવા વપરાશ મર્યાદાએ `auto_switch`, વાતચીતને જાતે બીજા એકાઉન્ટ પર ચાલુ રાખી શકે છે, સેશનના ઓપ્શન, permission mode, effort level અને agent સાથે; `/ccam:switch` વાતચીતને બીજી સંસ્થાના એકાઉન્ટ પર ખસેડે તે પહેલાં, ccam તમને પૂછે છે ([વિગતો](/features/next#relay)).

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--dry-run` | ccam જે environment અને કમાન્ડ ચલાવશે તે છાપીને બહાર નીકળે છે | off |
| `--fallback` | `-p`/`--print` runs માટે: એકાઉન્ટ તેની વપરાશ મર્યાદાએ પહોંચે ત્યારે એ જ command સૌથી વધુ જગ્યાવાળા લૉગિન થયેલા એકાઉન્ટ તરીકે ફરી ચલાવે છે ([વિગતો](/features/next#run-again-at-the-limit)) | off |
| `-v`, `--version` | `ccam <version>` છાપે છે (ફક્ત પહેલા આર્ગ્યુમેન્ટ તરીકે) | |
| `-h`, `--help` | ccam ની મદદ, `ccam --help` કે `ccam run --help` તરીકે | |

આર્ગ્યુમેન્ટના નિયમો:

- ccam પોતાના flag ફક્ત એકાઉન્ટના નામ પહેલાં જ વાંચે છે. નામ પછી, `--dry-run` અને `--help` સહિત દરેક આર્ગ્યુમેન્ટ Claude Code નો છે.
- `--dry-run` `ccam --dry-run work`, `ccam run --dry-run work` અને `ccam --dry-run` (ફોલ્ડરના એકાઉન્ટ કે ડિફૉલ્ટ માટે) તરીકે કામ કરે છે. `--dry-run` અને `--fallback` કોઈ પણ ક્રમમાં આવી શકે.
- `ccam run <name|number> [claude args]` એ `ccam <name>` નું લાંબું સ્વરૂપ છે. નામ વગરનું `ccam run` ફોલ્ડરનું એકાઉન્ટ કે ડિફૉલ્ટ એકાઉન્ટ વાપરે છે.
- એકાઉન્ટનાં નામ `add`, `list`, `run` કે `usage` જેવાં ccam કમાન્ડનાં નામ ન હોઈ શકે, તેથી `ccam <word>` ક્યારેય અસ્પષ્ટ નથી.

dry run એક લાઇન છાપે છે, જે તમે shell માં, કે Windows પર PowerShell માં, પેસ્ટ કરી શકો. તે લૉન્ચ environment માં જે બદલે છે તેનાથી શરૂ થાય છે: હંમેશાં `CLAUDE_CONFIG_DIR` અને `CCAM_PROFILE`, તે જે variable કાઢી નાખે છે તેના માટે `env -u <name>` (જેમ કે વારસામાં મળેલો `CCAM_RELAY`), અને, [API key એકાઉન્ટ](/features/api) માટે, તેનું endpoint, તેનું મૉડલ અને key ક્યાંથી આવે છે તે (`ANTHROPIC_API_KEY="$WORK_KEY"`), key પોતે ક્યારેય નહીં:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

PowerShell લાઇન પછીથી variables પાછા મૂકે છે, જેમ shell માં `VAR=x cmd` કરે છે, અને Claude Code ને તેની કમાન્ડ લાઇન `Start-Process` દ્વારા આપે છે, તેથી Windows PowerShell 5.1 અને PowerShell 7 તેને એક જ રીતે ચલાવે છે. જ્યારે લૉન્ચ કોઈ variable ને ખાલી મૂલ્ય પર સેટ કરે, ત્યારે તેની નીચેની `#` લાઇન જણાવે છે કે 7.5 પહેલાંનું PowerShell આવા variable ને તેના બદલે દૂર કરે છે. npm એ ઇન્સ્ટૉલ કરેલા `claude.cmd` માટે, લાઇન એ `cmd.exe` કમાન્ડ ચલાવે છે જે લૉન્ચ વાપરે (નીચે જુઓ), અને લૉન્ચ ત્યાં જે આર્ગ્યુમેન્ટનો ઇનકાર કરે તે dry run ને પણ નિષ્ફળ કરે છે.

## તે કેવી રીતે કામ કરે છે

દરેક શરૂઆત કમાન્ડ લાઇનથી `exec` સુધી એ જ પગલાં લે છે:

<!-- motion:launch -->

દરેક બૉક્સ તેનું પગલું ચાલે ત્યારે પ્રકાશિત થાય છે, અને નીચેની લાઇનો એ environment અને આર્ગ્યુમેન્ટ્સ છે જેની સાથે Claude Code શરૂ થાય છે.

1. ccam `app.Main` માં શરૂ થાય છે, જેને `cmd/ccam/main.go` બોલાવે છે. જો ccam બીજા નામે શરૂ થયું હોય, તો તે [શૉર્ટકટ](/features/shortcuts) છે અને સીધું launcher પર જાય છે.
2. cobra સામેલ થાય તે પહેલાં ccam `os.Args` જુએ છે (`app.parse`). ફક્ત એકલું `ccam` (ડેશબોર્ડ), ccam કમાન્ડનું નામ (`config.Reserved`), `-h`/`--help` અને `-v`/`--version` `cli.Execute` પર જાય છે. બાકીનું બધું શરૂઆત છે.
3. ccam `~/.config/ccam/config.toml` વાંચે છે, જે એકમાત્ર ફાઇલ તે parse કરે છે (`launcher.Launch`).
4. ccam નામ કે નંબરથી એકાઉન્ટ શોધે છે, શરૂઆતનું એક `--` કાઢી નાખે છે, અને `claude` executable શોધે છે (`launcher.Prepare`, `ResolveClaude`). નામ ન હોય તો, તે ફોલ્ડરનું એકાઉન્ટ લે છે, working directory ને `config.toml` ની `[[folder]]` entries સાથે strings તરીકે સરખાવીને, macOS અને Windows પર નાના-મોટા અક્ષરોનો ભેદ રાખ્યા વગર (`HereProfile`). PATH lookup ફાઇલના modification time અને `PATH` ના hash સાથે `~/.local/state/ccam/claude-path` માં કૅશ થાય છે, તેથી ccam દરેક શરૂઆતે `PATH` શોધતું નથી. બદલાયેલી કે ગુમ થયેલી ફાઇલ, કે બદલાયેલો `PATH`, નવું lookup કરાવે છે. `config.toml` માંનો `claude_bin` પાથ (`~/` વિસ્તારાય છે) ચલાવી શકાય તેવી ફાઇલ હોવો જોઈએ; નહીં તો કંઈ પણ લખાય તે પહેલાં શરૂઆત `claude_bin X is not an executable file. Check claude_bin in config.toml` સાથે અટકે છે. જે `claude` પોતે ccam હોય તેનો ccam inode વડે સરખામણી કરીને ઇનકાર કરે છે, કારણ કે તે loop બનાવે.
5. ccam એકાઉન્ટ ફોલ્ડર `~/.local/share/ccam/profiles/<name>/` ની તેનાં શેરિંગ સેટિંગ્સ સાથે સરખામણી કરે છે, ફક્ત `lstat` અને `readlink` ને બોલાવીને, આઇટમ દીઠ એક કે બે system call (`checkLinks`, `share.Check`). ખૂટતી લિંક ત્યાં જ બને છે. ડેટા ગુમાવી શકે તેવું કંઈ પણ શરૂઆત અટકાવે છે.
6. ccam શેર કરેલા MCP સર્વર ઉમેરે છે (`applyMCP`; નીચે જુઓ).
7. ccam સેશન નોંધે છે (`session.Register`). તે પોતાના PID ના નામવાળું, process ના શરૂ થવાના સમય સાથેનું marker `~/.local/state/ccam/locks/profile-<name>/<pid>` માં લખે છે, `XDG_RUNTIME_DIR` ગમે તે કહે, જેથી દરેક shell એ જ markers જુએ. વાતચીતનો ઇતિહાસ શેર કરતા એકાઉન્ટ માટે, ccam પહેલાં `history-<group>-<hash>` માં ચાલતું સેશન શોધે છે, જ્યાં hash symbolic links ઉકેલેલા હાલના ફોલ્ડરનો છે, અને ત્યાં પણ નોંધણી કરે છે (`historyCheck`).

   `exec` PID રાખે છે, તેથી પછી marker Claude Code process ને જ દર્શાવે છે. ccam પકડી રાખેલા lock ને બદલે marker ફાઇલ વાપરે છે, કારણ કે Claude Code જે process શરૂ કરે તે lock વારસામાં મેળવે અને સેશન પૂરું થયા પછી પણ રાખે. જે marker લખી ન શકાય તે શરૂઆત અટકાવતું નથી: ccam `ccam: could not mark the session as running (…); starting Claude Code anyway` કહે છે.

   એકાઉન્ટ ચાલે છે કે નહીં તેની ccam ની બીજી તપાસો (નામ બદલવું, દૂર કરવું, ગ્રુપ ફેરફાર, `doctor --fix`, MCP sync) એકાઉન્ટ ફોલ્ડરમાંના Claude Code ના પોતાના રેકોર્ડ (`sessions/<pid>.json`, `daemon.lock`) પણ વાંચે છે, તેથી તેઓ ccam વગર શરૂ થયેલું Claude Code પણ જુએ છે: સાદું `claude`, IDE, બૅકગ્રાઉન્ડ સેશન અને તેમનો daemon.
8. ccam ટર્મિનલનું શીર્ષક `claude · <name>` કરે છે (standard error ટર્મિનલ હોય ત્યારે), શરૂઆતનો સમય `fsync` વગર `~/.local/state/ccam/last-used/<name>` માં લખે છે, અને `syscall.Exec` ને બોલાવે છે (`Plan.Exec`).

Claude Code `CLAUDE_CONFIG_DIR` ને એકાઉન્ટ ફોલ્ડર પર અને `CCAM_PROFILE` ને એકાઉન્ટના નામ પર સેટ કરીને ચાલે છે. ccam આ variables નાં વારસામાં મળેલાં મૂલ્યો, અને વારસામાં મળેલો `CCAM_RELAY`, સાફ કરે છે. macOS પર `ccam import` વડે રજિસ્ટર કરેલા `~/.claude` ફોલ્ડર માટે, તે `CLAUDE_SECURESTORAGE_CONFIG_DIR` ને ખાલી મૂલ્ય પર પણ સેટ કરે છે, જેથી એકાઉન્ટ અને સાદું `claude` એક જ Keychain લૉગિન વાપરે; Claude Code લૉગિન Credential Manager માં રાખે ત્યારે Windows પર પણ તે એવું જ કરે છે.

### exec શા માટે

ccam પોતાની process ને `claude` થી બદલી નાખે છે (`syscall.Exec`). પરિણામે:

- ccam memory માં રહેતું નથી;
- keystroke, સિગ્નલ અને exit code Claude Code ના પોતાના હોય છે;
- આર્ગ્યુમેન્ટ્સ array જ રહે છે, ક્યારેય shell string બનતા નથી.

Windows માં `exec` નથી. ત્યાં ccam એ જ environment, આર્ગ્યુમેન્ટ્સ અને console સાથે `claude.exe` ને child process તરીકે શરૂ કરે છે, તેની રાહ જુએ છે અને તેના exit code સાથે બહાર નીકળે છે (`Plan.Exec`). npm એ ઇન્સ્ટૉલ કરેલું `claude.cmd` `cmd.exe` માં તેના પૂરા પાથ સાથે ચાલે છે, અને દરેક આર્ગ્યુમેન્ટ quote કરેલો હોય છે જેથી `cmd.exe` તેમાંના કોઈને પોતાની syntax (`&`, `%PATH%`) તરીકે ન વાંચે; line break વાળો આર્ગ્યુમેન્ટ આ રીતે આપી શકાતો નથી અને native Claude Code ઇન્સ્ટૉલ કરવાની સલાહ સાથે તેનો ઇનકાર થાય છે.

Ctrl+C console માંથી સીધું Claude Code સુધી પહોંચે છે, તેથી ccam તેને અવગણે છે. સેશન marker માં ccam નો PID હોય છે, જે બરાબર Claude Code જેટલો જ સમય જીવે છે, અને Claude Code બહાર નીકળે ત્યારે marker દૂર થાય છે. રાહ જોતી વખતે ccam બીજું કંઈ કરતું નથી.

સેશન પછીનું કામ, જેમ કે `claude mcp add` વડે ઉમેરેલા સર્વર ધ્યાનમાં લેવા, આગલી શરૂઆતે થાય છે: ccam `.claude.json` ના modification time ની `~/.local/state/ccam/mcp/<name>/seen` સાથે સરખામણી કરે છે અને ફાઇલ બદલાઈ હોય ત્યારે જ તેને parse કરે છે.

### કમાન્ડ લાઇન કયા રસ્તે જાય છે

ત્રણ કમાન્ડ લાઇન, `internal/app` માંથી ત્રણ રસ્તા:

<!-- motion:launch-route -->

શૉર્ટકટ `argv[0]` પરથી ઓળખાય છે, `ccam <name>` પહેલા આર્ગ્યુમેન્ટ પરથી, અને બંને `syscall.Exec` માં પૂરા થાય છે. ફક્ત `ccam usage` જેવો ccam કમાન્ડ જ cobra સુધી પહોંચે છે.

### શેર કરેલા MCP સર્વર

ડિફૉલ્ટ `flag` વ્યૂહરચના સાથે, ccam ગ્રુપના MCP સર્વર Claude Code ને આર્ગ્યુમેન્ટ તરીકે, તમારા આર્ગ્યુમેન્ટ્સની પહેલાં, આપે છે:

<!-- motion:launch-mcp -->

ગ્રુપની `mcp.json` હોય તો જ ccam આર્ગ્યુમેન્ટ ઉમેરે છે, અને તેને `claude` પછી તરત મૂકે છે (`mcp.FlagArgs`, `Plan.InsertArgs`). તમારું પોતાનું `--mcp-config` કે `--strict-mcp-config` પછી પણ તેના પર લાગુ પડે છે. તે એક જ આર્ગ્યુમેન્ટ, `--mcp-config=<path>`, છે, કારણ કે Claude Code નું `--mcp-config` ઘણાં મૂલ્યો લે છે: બે આર્ગ્યુમેન્ટ તરીકે લખેલું હોય, તો તે પછીના prompt ને પણ config ફાઇલ તરીકે વાંચે.

flag ક્યાં ન જઈ શકે તે ccam નક્કી કરે છે (`claudeVerb`). તે Claude Code ની જેમ જ આર્ગ્યુમેન્ટ્સમાંથી Claude Code નો subcommand વાંચે છે, ઓપ્શનને તેમનાં મૂલ્યો સાથે છોડીને. `remote-control` કે `rc` ને flag મળતો નથી (નીચે), અને પહેલાં આપેલા બૅકગ્રાઉન્ડ-સેશન verb (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`) ને પણ નહીં. `ccam work mcp list`, `get` અને `remove` ને પણ મળતો નથી, કારણ કે Claude Code ત્યાં તેનો ઇનકાર કરે છે; તેના બદલે ccam standard error પર જણાવે છે કે ગ્રુપના સર્વર એકાઉન્ટનાં સેશનમાં પણ લોડ થાય છે:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

બાકીનો દરેક subcommand અને દરેક સેશન flag રાખે છે.

`config.toml` માં `mcp_strategy = "sync"` સાથે, ccam તેના બદલે સર્વર એકાઉન્ટની `.claude.json` માં લખે છે, અને એકાઉન્ટ પહેલેથી ચાલતું હોય તો એ છોડી દે છે. છેલ્લી તપાસ પછી ગ્રુપની `mcp.json` કે `.claude.json` બેમાંથી કોઈ ન બદલાયું હોય ત્યારે પણ લૉન્ચ આ કામ છોડી દે છે. જો એકાઉન્ટ પાસે ગ્રુપના સર્વર જેવા જ નામનો પોતાનો સર્વર હોય, તો sync એકાઉન્ટનો સર્વર રાખે છે અને એક વાર તેમ જણાવે છે. sync `.claude.json` માંથી ગાયબ થયેલી કૉપિ પણ પાછી મૂકે છે (`put back group default's MCP servers …`), અને ફક્ત પોતે લખેલી, ત્યાર પછી કોઈએ ન બદલેલી કૉપિ જ બદલે છે કે દૂર કરે છે.

### Remote Control

`ccam work rc` (અથવા `ccam work remote-control`) Claude Code નું Remote Control `work` તરીકે શરૂ કરે છે. `claude remote-control` તેના verb પહેલાં `--mcp-config` નકારે છે અને પછી એક પણ લેતું નથી, અને તેનાં સેશન MCP સર્વર ફક્ત એકાઉન્ટની `.claude.json` માંથી વાંચે છે. તેથી અહીં flag વ્યૂહરચના જુદી રીતે કામ કરે છે:

<!-- motion:launch-rc -->

1. ccam `remote-control` કે `rc` ને Claude Code ના verb તરીકે શોધે છે, ઓપ્શન અને તેમનાં મૂલ્યો પછી પણ (`ccam work --model opus rc`), અને કોઈ flag ઉમેરતું નથી (`claudeVerb`, `mcp.FlagArgs`).
2. ccam ગ્રુપની `mcp.json` ના modification time ની એકાઉન્ટના `sync-stamp` સાથે સરખામણી કરે છે, અને તપાસે છે કે ગ્રુપનો દરેક સર્વર હજી `.claude.json` માં છે (`mcp.Synced`). જો `.claude.json` માં પહેલેથી ગ્રુપના હાલના સર્વર હોય, અથવા ગ્રુપમાં કોઈ સર્વર ન હોય, તો Remote Control કોઈ પ્રશ્ન વગર શરૂ થાય છે.
3. નહીંતર સેશનને જે સર્વર નહીં મળે તેમનાં નામ ccam જણાવે છે અને પૂછે છે:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. `y` પર, ccam તેમને `sync` વ્યૂહરચનાની જેમ જ કૉપિ કરે છે (`mcp.Sync`): Claude Code ના lock હેઠળ, બૅકઅપ પછી, એકાઉન્ટના પોતાના સર્વર જાળવીને. `N` પર, અથવા ટર્મિનલ વગર, Remote Control તેમના વગર શરૂ થાય છે, અને ccam `mcp_strategy = "sync"` તરફ ધ્યાન દોરે છે, જે દરેક શરૂઆત પહેલાં તેમને કૉપિ કરે છે.

સર્વર કૉપિ ન થાય ત્યાં સુધી ccam દરેક Remote Control શરૂઆતે પૂછે છે. કૉપિ પછી, તે ફરી ફક્ત ત્યારે જ પૂછે છે જ્યારે ગ્રુપની `mcp.json` બદલાય કે કૉપિ કરેલો સર્વર `.claude.json` માંથી ગાયબ થાય. પછીની સામાન્ય શરૂઆતોમાં `.claude.json` ની કૉપિ `--mcp-config` સામે જીતે છે, તેથી `ccam mcp add`, `remove` અને `import` બીજાં એકાઉન્ટ પાસે રહેલી કૉપિ અપડેટ કરે છે ([શેર કરેલા MCP સર્વર](/features/mcp)).

### 10 ms ની મર્યાદા

શરૂઆતના પાથનું બજેટ 95મા percentile પર 10 ms છે; વ્યવહારમાં Claude Code શરૂ થાય તે પહેલાં ccam લગભગ 4 ms ઉમેરે છે. શરૂઆતનો પાથ ક્યારેય cobra command tree બનાવતો નથી, Bubble Tea શરૂ કરતો નથી કે નેટવર્ક કનેક્શન ખોલતો નથી, અને package ના `init()` functions કુલ 5 ms થી વધુ લે તો `TestPackageInitIsCheap` નિષ્ફળ જાય છે. `scripts/perf-check.sh` `hyperfine` વડે `ccam --dry-run work`, `ccam work` અને એક શૉર્ટકટ માપે છે અને 10 ms થી ઉપર નિષ્ફળ જાય છે, અને `internal/launcher` માંનું `BenchmarkLaunchPath` પગલાં loop માં ચલાવે છે. દરેક પગલાનો સમય જોવા માટે `CCAM_TRACE=1` સેટ કરો:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## સલામતી અને ખાસ કિસ્સા

- અજાણ્યું નામ કે નંબર `no account "<name>"` સાથે અટકે છે અને `ccam list` તરફ નિર્દેશ કરે છે. કોઈ એકાઉન્ટ જ ન હોય, ત્યારે ccam તમને `ccam add <name>` વડે એક ઉમેરવા કહે છે.
- જો કોઈ શેર કરેલો પાથ સામાન્ય ફાઇલ બની ગયો હોય કે ખોટી જગ્યાએ નિર્દેશ કરતો હોય, તો ccam `stopped because of share link problems:` સાથે અટકે છે, તેમની યાદી આપે છે, અને `ccam doctor --fix` ચલાવવા કહે છે. શરૂઆતના પાથ પર તે ડેટા ખસેડતું નથી.
- જો આ એકાઉન્ટ સાથે વાતચીતનો ઇતિહાસ શેર કરતું (એ જ ગ્રુપનું) બીજું Claude Code સેશન એ જ ફોલ્ડરમાં ચાલી રહ્યું હોય, તો ccam ચેતવણી આપે છે કે ઇતિહાસ ભળી શકે છે અને `Run anyway? [y/N]` પૂછે છે. ફોલ્ડર symbolic links ઉકેલીને મેળવાય છે. ટર્મિનલ વગર (`/dev/null` માંથી આવતું standard input ટર્મિનલ ગણાતું નથી), અથવા બૅકગ્રાઉન્ડમાંના job માંથી, તે ચેતવણી બતાવીને આગળ વધે છે. નવો `-p` run પોતાની વાતચીત શરૂ કરે છે અને તેને કોઈ ચેતવણી મળતી નથી; `--version`, `--help` અને કોઈ વાતચીત ન ખોલતા subcommands, જેમ કે `mcp`, `auth` કે `plugin`, બિલકુલ તપાસાતા નથી.
- `--dry-run` MCP sync, નવા સર્વરની સૂચના, `mcp list` ની નોંધ, ઇતિહાસની તપાસ, સેશન marker, શીર્ષક અને છેલ્લા ઉપયોગનો સમય છોડી દે છે, અને exec કરતું નથી. તે તો પણ શેર લિંક તપાસે છે, અને સામાન્ય શરૂઆતની જેમ ખૂટતી લિંક બનાવે છે.
- `ccam work rc`, જ્યારે `work` પહેલેથી ચાલતું હોય, ત્યારે `.claude.json` ને અડતું નથી, એ જણાવે છે, અને ગ્રુપના સર્વર વગર Remote Control શરૂ કરે છે; તેમને કૉપિ કરવા માટે સેશન પૂરું થયા પછી તેને ફરી ચલાવો. `--dry-run` સાથે તે સંદેશ છાપે છે પણ ન તો પૂછે છે, ન કૉપિ કરે છે.
- `mcp_strategy = "sync"` સાથે, નિષ્ફળ sync `ccam: MCP sync failed (launching anyway): …` છાપે છે અને Claude Code પાસે પહેલેથી હતા તે સર્વર સાથે શરૂ થાય છે.

## સંબંધિત

- [શૉર્ટકટ કમાન્ડ](/features/shortcuts)
- [ccam સાથે Remote Control વાપરવું](/guides/remote-control)
- [એકાઉન્ટનું સંચાલન](/features/accounts)
- [કમાન્ડ સંદર્ભ](/docs/commands)
- [MCP સર્વર, સેટિંગ્સ અને ઇતિહાસ શેર કરો](/docs/sharing)
- [CLAUDE_CONFIG_DIR સમજૂતી](/guides/claude-config-dir)
