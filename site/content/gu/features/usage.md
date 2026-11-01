---
title: વપરાશ મર્યાદા
description: ccam usage દરેક એકાઉન્ટની 5 કલાકની અને સાપ્તાહિક મર્યાદા એકસાથે બતાવે છે, અને ccam statusline એકાઉન્ટનું નામ અને વપરાશ Claude Code માં મૂકે છે.
updated: 2026-11-01
order: 7
translation: machine
sourceHash: 54b296f7da35
---

`ccam usage` દરેક એકાઉન્ટ માટે બતાવે છે કે દરેક મર્યાદામાંથી કેટલું વપરાયું છે: 5 કલાકની મર્યાદા, સાપ્તાહિક મર્યાદા અને દરેક મૉડલની પોતાની સાપ્તાહિક મર્યાદા (જેમ કે Fable). દરેક મર્યાદાને એક બાર અને તે રીસેટ થવામાં બાકી સમય મળે છે. કોઈ લાંબું કામ શરૂ કરતાં પહેલાં, અથવા Claude Code મર્યાદા નજીક હોવાનું કહે ત્યારે, કયા એકાઉન્ટ પાસે જગ્યા બાકી છે તે જોવા તેને વાપરો. `ccam statusline` એ Claude Code માટેનો સ્ટેટસ લાઇન કમાન્ડ છે: તે Claude Code ની અંદર એકાઉન્ટનું નામ અને વપરાશ બતાવે છે અને `ccam usage` માટે આંકડા સેવ કરે છે. વપરાશ પાછળના token અને અંદાજિત ખર્ચ માટે, `ccam ccusage` દરેક એકાઉન્ટની વાતચીત પર [ccusage](https://ccusage.com) ચલાવે છે ([નીચે](#tokens-and-costs-ccam-ccusage)).

## ઉપયોગ

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts that are not running
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

દરેક એકાઉન્ટને એક લાઇન મળે છે જેમાં તેનો નંબર, તેનું નામ અને તેના આંકડા કેટલા જૂના છે તે હોય છે; `(cached)` નો અર્થ છે કે તે endpoint માંથી નહીં, સેવ કરેલી ફાઇલમાંથી આવે છે. તેની નીચે, દરેક મર્યાદા માટે એક લાઇન 20 ખાનાંનો બાર, ટકાવારી અને વિન્ડો રીસેટ થવામાં બાકી સમય બતાવે છે: `5-hour`, `Weekly` (બધાં મૉડલ), પછી દરેક મૉડલની પોતાની સાપ્તાહિક વિન્ડો, જેમ કે `Fable`. જે વિન્ડો એકાઉન્ટ પાસે ન હોય તે `-` બતાવે છે. ટકાવારી 80% પર પીળી અને 100% પર લાલ થાય છે.

જો Claude Code એ કોઈ વિન્ડો ભરાઈ ગઈ હોવાથી request નકારી હોય, તો એ વિન્ડો, સેવ કરેલા આંકડા ભલે ગમે તે કહે, તેના રીસેટ સુધી 100% બતાવે છે ([મર્યાદા રેકોર્ડ](#the-limit-record) જુઓ). ટર્મિનલ પર દરેક એકાઉન્ટ તરત `fetching…` તરીકે દેખાય છે, અને તે જવાબ આપે તેમ તેની લાઇનો દોરાય છે.

છેલ્લી લાઇન, જે તમારી પાસે એક કરતાં વધુ એકાઉન્ટ હોય ત્યારે જ દેખાય છે, એ એકાઉન્ટનું નામ આપે છે જે [`ccam next --auto`](/features/next) પસંદ કરે. તે એ જ નિયમો અનુસરે છે:

- `auto_max_usage` અને `auto_strategy`;
- ચાલી રહેલાં સેશન;
- `cross_org` સેટ ન હોય તો, ફક્ત છેલ્લે વપરાયેલા એકાઉન્ટની સંસ્થાનાં એકાઉન્ટ;
- 100% પરની દરેક વિન્ડો છોડી દેવાય છે, તેથી જેનું અઠવાડિયું ભરાઈ ગયું હોય તે એકાઉન્ટનું નામ ક્યારેય આવતું નથી.

ડિફૉલ્ટ `headroom` strategy સાથે તે `→ Most headroom: <name>` બતાવે છે; `expiring` કે `balance` સાથે તે કારણ આપે છે, જેમ કે `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Team અને Enterprise સીટ.** આપમેળે થતી એકાઉન્ટ પસંદગી (`ccam next`, `--fallback`, `ccam pick`, ડેશબોર્ડની `n` કી અને, તમે તેને ચાલુ કરો પછી, [`auto_switch`](/features/next#auto-switch)) એક જ Claude Team કે Enterprise સંસ્થાની સીટો વચ્ચે પણ સ્વિચ કરે છે, જેનાથી Anthropic ની શરતોનું ઉલ્લંઘન થઈ શકે છે. તેના પર આધાર રાખતાં પહેલાં [વારંવાર પૂછાતા પ્રશ્નો](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) વાંચો અને તમારી સંસ્થા અને Anthropic પાસેથી ખાતરી કરી લો.

લૉગિન ન થયેલું એકાઉન્ટ `not logged in → ccam login <name>` બતાવે છે, અને API key વાળું એકાઉન્ટ `no limits (API)` બતાવે છે. મુદત પૂરી થયેલો access token ``access token expired; `ccam <name>` renews it when it starts`` બતાવે છે, સેવ કરેલા આંકડા હોય તો તેમની ઉપર: Claude Code આગલી શરૂઆતે તેને રિન્યૂ કરે છે. જે લૉગિનની મુદત પૂરેપૂરી પૂરી થઈ ગઈ હોય, કારણ કે તેના refresh token ની મુદત પૂરી થઈ હોય કે સર્વરે ઇનકાર કર્યા પછી Claude Code એ તેને સાફ કરી નાખ્યું હોય, તે `` login expired → `ccam login <name>` `` બતાવે છે.

`ccam statusline install` પછી, Claude Code `[work] · Opus · 5h 24% · 7d 41%` જેવી લાઇન બતાવે છે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--json` | રિપોર્ટ JSON તરીકે છાપે છે: `name`, `slot`, `status`, `limits` (એકાઉન્ટમાં usage credits ચાલુ હોય ત્યારે `credits` સાથે), `source`, `at`, `error`, અને જ્યારે મર્યાદા રેકોર્ડ વિન્ડોને 100% કરે ત્યારે `limited` અને `limited_windows`. કોઈ એકાઉન્ટ ન હોય ત્યારે તે `[]` છાપે છે | off |
| `--watch[=N]` | તમે Ctrl+C દબાવો ત્યાં સુધી દર `N` સેકન્ડે ફરીથી દોરે છે. એકલું `--watch`, અથવા 0 કે તેથી ઓછું મૂલ્ય, એટલે 60. મૂલ્ય `=` પછી આવે છે: `--watch 10` નો સૂચના સાથે ઇનકાર થાય છે | off |
| `-v`, `--verbose` | દરેક એકાઉન્ટ નીચે `source:` લાઇન ઉમેરે છે: તેના આંકડા ક્યાંથી આવ્યા (usage endpoint, Claude Code, સ્ટેટસ લાઇન, સેવ કરેલી ફાઇલ, કે Claude Code એ જણાવેલી વપરાશ મર્યાદા). જ્યારે મર્યાદા રેકોર્ડ આંકડા પર અગ્રતા લે, ત્યારે પછી `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` જેવી લાઇન આવે છે | off |
| `--no-cache` | 60 સેકન્ડનો કૅશ અવગણીને endpoint ને ફરીથી પૂછે છે | off |
| `--source live\|claude\|cache` | `live` એકાઉન્ટના token વડે endpoint ને પૂછે છે; `claude` Claude Code પાસે તે જોવડાવે છે (નીચે જુઓ); `cache` ફક્ત સેવ કરેલા આંકડા બતાવે છે | `usage_source`, નહીં તો `live` |
| `--refresh` | મુદત પૂરી થયેલા access token રિન્યૂ કરે છે, ફક્ત ચાલતાં ન હોય તેવાં એકાઉન્ટ માટે. રિન્યૂ કરવામાં 30 સેકન્ડ સુધી લાગી શકે છે | off |
| `--dry-run` | સ્વીકારાય છે: `ccam usage` પોતાના કૅશ સિવાય કંઈ બદલતું નથી. `--refresh` સાથે ઇનકાર થાય છે: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | off |

`--source` નું બીજું કોઈ પણ મૂલ્ય `--source must be live, claude or cache` સાથે નિષ્ફળ જાય છે. `ccam usage` કોઈ આર્ગ્યુમેન્ટ લેતું નથી.

### Claude Code પાસેથી જ વપરાશ

`config.toml` માં સૌથી ઉપર (પહેલા `[[group]]`, `[[profile]]` કે `[[folder]]` ની ઉપર) `usage_source = "claude"` સાથે, અથવા `--source claude` સાથે, ccam લૉગિન token બિલકુલ વાંચતું નથી. દરેક એકાઉન્ટ માટે તે એકાઉન્ટના ફોલ્ડરમાં `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` શરૂ કરે છે અને standard input પર Claude Code ની `get_usage` request મોકલે છે (Claude Code 2.1.270 કે પછીનું). Claude Code પોતાના લૉગિન વડે વપરાશ જુએ છે અને 5 કલાકની, સાપ્તાહિક અને મૉડલ વિન્ડો સાથે જવાબ આપે છે.

`--safe-mode` એકાઉન્ટના hooks, MCP સર્વર અને plugins ને શરૂ થતાં રોકે છે, અને કોઈ સેશન સેવ થતું નથી. દરેક એકાઉન્ટને એક Claude Code શરૂઆત (લગભગ એક સેકન્ડ) જોઈએ, તેથી આ source `live` કરતાં ધીમો છે. એકાઉન્ટને એકસાથે પૂછવામાં આવે છે, દરેક 20 સેકન્ડના timeout સાથે, અને 60 સેકન્ડનો cache હંમેશની જેમ લાગુ પડે છે.

જે લૉગિનના access token ની મુદત પૂરી થઈ ગઈ હોય, કે 5 મિનિટની અંદર પૂરી થવાની હોય, તેને 50 સેકન્ડ મળે છે, કારણ કે Claude Code શરૂ થતાં જ તેને રિન્યૂ કરે છે અને અધવચ્ચે કપાયેલું રિન્યૂઅલ લૉગિન ગુમાવી શકે છે. જો ccam કોઈ lookup જતો કરે, તો તે Claude Code ને SIGTERM વડે બંધ થવા કહે છે, જેથી Claude Code શરૂ કરેલું રિન્યૂઅલ પૂરું કરે (Windows માં SIGTERM નથી, તેથી ત્યાં process ને બંધ કરી દેવાય છે).

જે લૉગિનની મુદત પૂરેપૂરી પૂરી થઈ ગઈ હોય તે Claude Code શરૂ કર્યા વગર `` login expired → `ccam login <name>` `` બતાવે છે. Claude Code આ request ને પ્રાયોગિક ગણાવે છે, તેથી જો તે જવાબ આપવાનું બંધ કરે, તો ccam સેવ કરેલા આંકડા બતાવે છે અને `ccam doctor` તેની જાણ કરે છે. જ્યારે Claude Code જવાબ આપે કે તેનો પોતાનો lookup નિષ્ફળ ગયો, ત્યારે ccam મર્યાદા વગરનો પ્લાન નહીં, પણ `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)` સાથે સેવ કરેલા આંકડા બતાવે છે.

`--dry-run` સાથે, ccam Claude Code શરૂ કરતું નથી, કારણ કે Claude Code એકાઉન્ટના ફોલ્ડરમાં લખે છે; તેના બદલે તે સેવ કરેલા આંકડા બતાવે છે. `usage_source` માં `live` અથવા `claude` આપી શકાય છે; `cache` ફક્ત `--source` નો વિકલ્પ છે.

`ccam statusline` standard input માંથી Claude Code નું સ્ટેટસ લાઇન JSON વાંચે છે. તેનો એકમાત્ર ઓપ્શન `--wrap <command>` છે: હાલનો સ્ટેટસ લાઇન કમાન્ડ, જેની આઉટપુટની પહેલી લાઇનની આગળ એકાઉન્ટનું નામ મુકાય છે. `ccam statusline install` કોઈ ઓપ્શન લેતું નથી. બંને `` --dry-run is not supported by `ccam statusline`; nothing was changed `` સાથે `--dry-run` નો ઇનકાર કરે છે.

## Token અને ખર્ચ: `ccam ccusage`

`ccam usage` બતાવે છે કે કોઈ એકાઉન્ટે દરેક મર્યાદામાંથી કેટલું વાપર્યું છે. તેના બદલે [ccusage](https://ccusage.com) Claude Code ની વાતચીતની log ફાઇલોમાંના token અને અંદાજિત ખર્ચ ગણે છે. એકલું ચાલે ત્યારે તે Claude Code ની log ફાઇલો ફક્ત `~/.claude` અને `~/.config/claude` માં (અથવા `CLAUDE_CONFIG_DIR` જે ફોલ્ડરનાં નામ આપે તેમાં) શોધે છે, તેથી ccam નું જે એકમાત્ર એકાઉન્ટ તેને દેખાઈ શકે તે `~/.claude` માં જ્યાં છે ત્યાં જ ઇમ્પોર્ટ થયેલું એકાઉન્ટ છે. `ccam ccusage` તેને `CLAUDE_CONFIG_DIR` માં વાતચીત ધરાવતા દરેક એકાઉન્ટનાં ફોલ્ડર સેટ કરીને ચલાવે છે, જેથી એક જ રિપોર્ટ એ બધાંને આવરી લે:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam પહેલાં પોતાના ઓપ્શન વાંચે છે: `-a` (`--account`) એકાઉન્ટનાં નામ કે નંબર લે છે, અલ્પવિરામથી અલગ કરેલાં કે ઓપ્શન ફરી આપીને, અને `--dry-run` કમાન્ડ ચલાવવાને બદલે તેને છાપે છે. ccam ના ઓપ્શન સિવાયના પહેલા આર્ગ્યુમેન્ટથી શરૂ કરીને બધું જ તમે ટાઇપ કર્યું હોય તેમ જ ccusage ને જાય છે, ફક્ત `bun x` અને `npx` માટે `--config` સિવાય (નીચે જુઓ); `ccam ccusage -- --help` ccusage ની પોતાની મદદ બતાવે છે. ccam ઇન્સ્ટૉલ કરેલું `ccusage` ચલાવે છે, તે ન હોય તો `bun x ccusage`, અને તે પણ ન હોય તો `npx ccusage@latest`; આ છેલ્લાં બે તેને npm પરથી ડાઉનલોડ કરે છે. આ બંને ccam ના પોતાના ફોલ્ડર (`~/.local/state/ccam/ccusage`) માં શરૂ થાય છે, તમે જે ફોલ્ડરમાં છો તેમાં નહીં, તેથી ત્યાંનો કોઈ પ્રોજેક્ટ તેના `node_modules/.bin`, તેના `package.json` માંના કોઈ `bin` કે તેના `.npmrc` માંની કોઈ registry વડે તે શું ચલાવે તે નક્કી કરી શકતો નથી; `ccam --dry-run ccusage` તે ફોલ્ડરમાં જતો `cd` બતાવે છે. આમાંનું એક પણ ન હોય ત્યારે, તે Node.js, Bun કે ccusage ઇન્સ્ટૉલ કરવાનું કહે છે. Linux અને macOS પર ccusage, ccam ની જગ્યા લે છે; Windows પર ccam તેની રાહ જુએ છે અને ccusage ના exit code સાથે બહાર નીકળે છે. WSL ની અંદર, ccam Windows ડ્રાઇવ પર આવેલા `ccusage`, `bun` કે `npx` ને છોડી દે છે, કારણ કે તે એકાઉન્ટ ફોલ્ડર વગર Windows પ્રોગ્રામ તરીકે ચાલશે.

- રિપોર્ટ એ એકાઉન્ટને આવરી લે છે જેમના ફોલ્ડરમાં વાતચીત હોય (`projects` ફોલ્ડર). ગ્રુપ દ્વારા એકાઉન્ટ જે વાતચીત શેર કરે છે તે એક જ વાર ગણાય છે.
- `-a` સાથે, ગ્રુપ દ્વારા પોતાનો ઇતિહાસ શેર કરતું એકાઉન્ટ આખા ગ્રુપની વાતચીત સાથે લાવે છે, અને ccam જણાવે છે કે તે કોની છે: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- `blocks` 5 કલાકની વિન્ડો બતાવે છે, અને દરેક એકાઉન્ટની પોતાની વિન્ડો હોય છે, તેથી તેના માટે ઉપર બતાવ્યા પ્રમાણે એક જ એકાઉન્ટનું નામ આપો.
- તમે કોઈ એકનું નામ ન આપો ત્યાં સુધી, ccusage તેને મળતાં બીજાં coding ટૂલ, જેમ કે Codex, નો પણ રિપોર્ટ આપે છે: `ccam ccusage claude daily` ફક્ત Claude Code નો રિપોર્ટ આપે છે.
- ccam તેને ચલાવે ત્યારે, ccusage તેનાં સેટિંગ્સ (`ccusage.json`) ફક્ત `./.ccusage` અને એકાઉન્ટ ફોલ્ડરમાંથી વાંચે છે. `--config` બીજી ફાઇલનું નામ આપે છે: `ccam ccusage daily --config ~/.config/claude/ccusage.json`. `bun x` અને `npx` ccam ના પોતાના ફોલ્ડરમાં શરૂ થતાં હોવાથી, `./.ccusage/ccusage.json` હોય અને તમે કોઈ ફાઇલનું નામ ન આપો ત્યારે ccam તેમને તે `--config` સાથે આપે છે, અને relative `--config` ને પૂરા પાથ તરીકે આપે છે. પછી ccusage બીજા relative પાથ, જેમ કે `--pi-path`, ccam ના ફોલ્ડરમાંથી વાંચે છે, તેથી તે પૂરા આપો.
- જે ફોલ્ડરના પાથમાં અલ્પવિરામ હોય તેને ccusage બે ફોલ્ડર તરીકે વાંચે, તેથી ccam તેનો ઇનકાર કરે છે અને એકાઉન્ટનું નામ જણાવે છે; `-a` વડે એ એકાઉન્ટને બહાર રાખો.
- `ccusage` ccam નો કમાન્ડ છે, તેથી ccusage નામનું એકાઉન્ટ `ccam run ccusage` વડે શરૂ થાય છે; `ccam doctor` આ જણાવે છે.

## તે કેવી રીતે કામ કરે છે

`ccam usage` ના આંકડા બે જગ્યાએથી આવે છે: Anthropic નું દસ્તાવેજીકૃત ન થયેલું usage endpoint, `GET https://api.anthropic.com/api/oauth/usage`, જેને દરેક એકાઉન્ટના પોતાના લૉગિન સાથે બોલાવવામાં આવે છે, અને Claude Code પોતાના સ્ટેટસ લાઇન કમાન્ડને જે મૂલ્યો આપે છે તે. આ બધું `internal/usage` માં છે.

1. ccam `config.toml` લોડ કરે છે અને દરેક એકાઉન્ટ માટે એક `usage.Account` બનાવે છે (`cli.usageCmd`): તેનું ફોલ્ડર `~/.local/share/ccam/profiles/<name>/`, તેના સેશન markers (એકાઉન્ટ ફોલ્ડરમાંના Claude Code ના પોતાના રેકોર્ડ સાથે, તે જણાવે છે કે તે ચાલી રહ્યું છે કે નહીં) અને તે જ્યાં છે ત્યાં જ ઇમ્પોર્ટ થયું હતું કે નહીં.
2. ccam બધાં એકાઉન્ટ એકસાથે જુએ છે, દરેક એકાઉન્ટ માટે એક goroutine (`Fetcher.FetchAll`). તે બધાં એક જ `http.Client` વાપરે છે, તેથી endpoint સાથેનાં કનેક્શન ફરી વપરાય છે. કનેક્ટ થવામાં 2 s લાગી શકે છે, અને દરેક એકાઉન્ટને કુલ 3 s મળે છે. client redirect અનુસરતો નથી, કારણ કે redirect ટોકનને બીજે ક્યાંક ફરીથી મોકલી દે, અને જવાબનો વધુમાં વધુ 1 MB વાંચે છે.
3. દરેક એકાઉન્ટ પૂરું થાય તેમ તેની હરોળ ફરીથી દોરાય છે. જે એકાઉન્ટ નિષ્ફળ જાય કે સમયસર જવાબ ન આપે તે સેવ કરેલા આંકડા પર આધાર રાખે છે, તેથી એક ધીમું એકાઉન્ટ ક્યારેય બીજાંને રોકતું નથી.

ડાયાગ્રામમાં, દરેક એકાઉન્ટ પોતાની goroutine માં જોવાય છે, અને જેનો સમય પૂરો થાય તે તેના બદલે પોતાના સેવ કરેલા આંકડા બતાવે છે:

<!-- motion:usage-fetch -->

`personal` અને `work` સમયસર જવાબ આપે છે અને live આંકડા બતાવે છે. `side` ને 3 s કરતાં વધુ લાગે છે, તેથી ccam તેના બદલે તેના સેવ કરેલા આંકડા બતાવે છે (`fallback()`). પછી ccam દરેક એકાઉન્ટ માટે એક બ્લૉક દોરે છે (`usage.Table()`) અને છેલ્લી લાઇનમાં જેનું નામ આવે છે તે એકાઉન્ટ પસંદ કરે છે (`usage.Best()`).

એક એકાઉન્ટ માટે, ccam આ તપાસો ક્રમમાં કરે છે (`Fetcher.Fetch`):

<!-- motion:usage-cache -->

1. `--source cache` સાથે, તે સીધું `fallback()` પર જાય છે.
2. `--no-cache` સેટ ન હોય ત્યાં સુધી, `~/.local/state/ccam/usage/<name>.json` માંનું 60 સેકન્ડ (`CacheTTL`) કરતાં ઓછા સમય પહેલાં માપેલું live પરિણામ વપરાય છે, હાલના સમય પ્રમાણે મૂલવીને: ત્યાર પછી રીસેટ થયેલી વિન્ડો 0% બતાવે છે, અને ઘડિયાળથી આગળની તારીખવાળી ફાઇલ તાજી ગણાતી નથી. `ccam usage` સળંગ બે વાર ચલાવવાથી એકાઉન્ટ દીઠ એક જ request થાય છે.
3. ccam `.credentials.json` (કે macOS Keychain) ને બદલ્યા વગર વાંચે છે અને ફક્ત `claudeAiOauth` key decode કરે છે (`readCredentials()`). subscription લૉગિન ન હોય તો `not logged in`, અથવા `.claude.json` માં API key હોય ત્યારે `no limits (API)`. જે લૉગિનના refresh token ની મુદત પૂરી થઈ હોય, કે જેને Claude Code એ સાફ કર્યું હોય, તે `` login expired → `ccam login <name>` `` છે અને તેનો lookup થતો નથી.
4. મુદત પૂરી થયેલો access token વપરાતો નથી. `--refresh` વગર, એકાઉન્ટ ``access token expired; `ccam <name>` renews it when it starts`` સાથે પોતાના સેવ કરેલા આંકડા બતાવે છે.
5. નહીં તો ccam endpoint ને બોલાવે છે. 401 કે 403 મુદત પૂરી થયેલું લૉગિન ગણાય છે (`login expired: authentication failed (401)`), બીજું કોઈ પણ status કે ન વાંચી શકાય તેવો body ભૂલ ગણાય છે, અને બંને સેવ કરેલા આંકડા પર આધાર રાખે છે. Ctrl+C થી અટકાવેલો lookup `cancelled` કહે છે. 429 નો અર્થ છે કે endpoint ને આ એકાઉન્ટ માટે બહુ વધારે requests મળી છે: ccam ફરી પૂછવાના સમય (`Retry-After`, અથવા 5 મિનિટ, વધુમાં વધુ એક કલાક) સાથે `<name>.pause` લખે છે, `rate limited by the usage endpoint; asking again after 15:10` સાથે સેવ કરેલા આંકડા બતાવે છે, અને ત્યાં સુધી, `--no-cache` સાથે પણ, તે એકાઉન્ટ માટે endpoint ને બોલાવતું નથી. તરત ફરી પ્રયાસ કરવાથી rate limit ફક્ત લંબાય.
6. સારો જવાબ માપ્યાના સમય સાથે `<name>.json` માં લખાય છે.

સેવ કરેલા આંકડા માટે, ccam સૌથી નવો સેવ કરેલો lookup (`<name>.json`, કે Claude Code ના lookups માંથી `<name>.claude.json`) લે છે અને, જ્યારે `<name>.statusline.json` તેના પછી લખાઈ હોય, ત્યારે તેના પર સ્ટેટસ લાઇનની 5 કલાકની અને સાપ્તાહિક વિન્ડો મૂકે છે (`fallback()`); મૉડલની વિન્ડો, જે સ્ટેટસ લાઇન પાસે નથી, lookup એ સેવ કરી હતી તેમ જ રહે છે. સેવ કરેલા આંકડા બતાવતાં પહેલાં, તે દરેક વિન્ડોનો રીસેટ સમય તપાસે છે: આંકડા સેવ થયા પછી રીસેટ થયેલી વિન્ડો 0% બતાવે છે, કારણ કે જૂની ટકાવારી હવે લાગુ પડતી નથી. તેથી જ ઉપરનું `personal` રીસેટ સમય વગર `0%` બતાવે છે.

endpoint ના બે જવાબ ફૉર્મેટ છે. નવામાં `session` (5 કલાક), `weekly_all` અને `weekly_scoped` (મૉડલ પ્રમાણે) એન્ટ્રીઓ સાથેનો `limits` array છે. જૂનામાં `five_hour`, `seven_day`, `seven_day_opus` અને `seven_day_sonnet` objects છે. ccam બંને વાંચે છે (`ParseLive`): `limits` array જીતે છે, અને તે જે આવરી ન લે તે જૂના objects પૂરું કરે છે. કોઈ મર્યાદા વગરનો body બદલાયેલા ફૉર્મેટ તરીકે નોંધાય છે, જે વિશે `ccam doctor` તમને જણાવે છે.

### મર્યાદા રેકોર્ડ

જ્યારે Claude Code કોઈ turn ને વપરાશ મર્યાદાએ અટકાવે, ત્યારે [ccam plugin](/docs/install#claude-code-plugin) નો hook વાતચીતમાંથી નકારાયેલી request વાંચે છે: કઈ વિન્ડો ભરાઈ ગઈ હતી (5 કલાકની, સાપ્તાહિક, Opus કે Sonnet) અને તે ક્યારે રીસેટ થાય છે. તે આ `~/.local/state/ccam/usage/<name>.limit.json` માં લખે છે. ccam દરેક source માંથી તે જે કંઈ જાણે છે તેના પર આ રેકોર્ડ લાગુ કરે છે (`Fetch`): તેમાંની દરેક વિન્ડો તેના રીસેટ સુધી 100% ગણાય છે, સેવ કરેલા આંકડા ગમે તેટલા જૂના કે ઓછા હોય, તેથી `ccam next --auto`, `ccam pick`, `--fallback` અને [`auto_switch`](/features/next#auto-switch) ત્યાં સુધી એ એકાઉન્ટ ફરી પસંદ કરતાં નથી.

બીજું કંઈ જાણીતું ન હોય, ત્યારે એકાઉન્ટ ફક્ત રેકોર્ડ જ બતાવે છે (`-v` સાથે `source: the usage limit Claude Code reported`). કોઈ સામે બેઠું ન હોય તેવાં સેશન (`-p` runs અને બૅકગ્રાઉન્ડ સેશન) કોઈ રેકોર્ડ લખતાં નથી, અને સંસ્થાએ લગાવેલો અટકાવ (ઉદાહરણ તરીકે, જ્યારે તેનો વપરાશ ખૂટી જાય કે કોઈ ગ્રુપની મર્યાદા $0 હોય) પણ નહીં, કારણ કે તેના admins તે ગમે ત્યારે હટાવી શકે છે. બીજા Claude એકાઉન્ટમાં `ccam login` રેકોર્ડ કાઢી નાખે છે; ફાઇલ કાઢી નાખવાથી પણ તે સાફ થાય છે.

### ટોકન રિન્યૂ કરવા

refresh token ફરતા રહે છે: એક રિન્યૂ કરવાથી તમને નવો refresh token મળે છે અને જૂનો અમાન્ય થાય છે. જો Claude Code એ જ એકાઉન્ટ તરીકે ચાલતું હોય ત્યારે ccam લૉગિન રિન્યૂ કરે, તો એ સેશન સાઇન આઉટ થઈ જાય. તેથી `ccam usage` ટોકન વાંચે છે અને ક્યારેય લખતું નથી, સિવાય કે તમે `--refresh` આપો. ત્યારે ccam (`Fetcher.refresh`):

1. `ccam import` વડે જ્યાં છે ત્યાં જ રજિસ્ટર થયેલાં એકાઉન્ટનો ઇનકાર કરે છે, કારણ કે સાદું `claude` ccam ની જાણ વગર એ લૉગિન વાપરતું હોઈ શકે;
2. macOS Keychain માં રાખેલાં લૉગિનનો ઇનકાર કરે છે;
3. lock ફાઇલ, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, લે છે, જેથી બે ccam process ક્યારેય એક જ ટોકન બંને રિન્યૂ ન કરે; તે પકડાયેલી હોય ત્યારે, એકાઉન્ટ `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)` બતાવે છે. `ccam rename`, `remove` અને `logout` પણ તે લે છે, અને રિન્યૂઅલ ચાલુ હોય ત્યારે ઇનકાર કરે છે;
4. lock હેઠળ, એકાઉન્ટનું કોઈ સેશન ચાલતું હોય તો ઇનકાર કરે છે, ભલે તે ગમે તે રીતે શરૂ થયું હોય (ccam દ્વારા, સાદું `claude`, IDE, કે બૅકગ્રાઉન્ડ સેશન અને તેની પાછળનો daemon, જે લૉગિન જાતે રિન્યૂ કરે છે), અને બીજા process એ ટોકન પહેલેથી રિન્યૂ કર્યો હોય તો request છોડી દે છે;
5. Claude Code એકાઉન્ટ ફોલ્ડરમાં પોતાનાં રિન્યૂઅલ lock (`.oauth_refresh.lock`, `<folder>.lock`) પકડી રાખે ત્યારે કંઈ મોકલતું નથી: ત્યારે Claude Code એ જ સમયે તે લૉગિન રિન્યૂ કરી રહ્યું હોય છે;
6. refresh token ને token endpoint પર post કરે છે અને જવાબ માટે 30 સેકન્ડ સુધી રાહ જુએ છે, જેટલો સમય Claude Code પોતાનાં રિન્યૂઅલને આપે છે. એક વાર મોકલ્યા પછી, request ક્યારેય અટકાવાતી નથી, Ctrl+C થી પણ નહીં: સર્વરે refresh token પહેલેથી બદલી નાખ્યો હોઈ શકે, અને ફક્ત તેના જવાબમાં જ નવું લૉગિન હોય છે. Ctrl+C કે hangup બીજા lookups પૂરા કરે છે, અને જવાબ સેવ થાય પછી ccam બહાર નીકળે છે (Ctrl+C પછી 130);
7. Claude Code ના write lock (`.storage-write.lock`) હેઠળ સેવ કરે છે અને તે સમયની `.credentials.json` નાં ફક્ત token fields જ બદલે છે: એ દરમિયાન Claude Code એ લખેલું (MCP સાઇન-ઇન) રહે છે, એ દરમિયાન Claude Code એ રિન્યૂ કરેલું લૉગિન રખાય છે, અને લૉગઆઉટ ક્યારેય પાછું ફેરવાતું નથી. જ્યારે જવાબ નવા refresh token નું આયુષ્ય આપે, ત્યારે તે Claude Code ની જેમ જ `refreshTokenExpiresAt` તરીકે સંગ્રહાય છે.

જે refresh પૂરું ન થાય તે એકાઉન્ટ લાઇન પર કે તેના સેવ કરેલા આંકડાની ઉપર કારણ જણાવે છે: ઉદાહરણ તરીકે `token not renewed: a session is running`, અથવા request મોકલાઈ હોય પણ તેનો જવાબ વાપરી કે સેવ કરી ન શકાયો હોય ત્યારે `login may be lost: …`, અને પછી Claude Code તમને લૉગિન કરવા કહે તો `ccam login <name>` ચલાવવાની સલાહ સાથે.

### સ્ટેટસ લાઇન collector

Claude Code સ્ટેટસ લાઇન કમાન્ડ ચલાવે છે અને તેને pipe દ્વારા JSON document આપે છે. ccam તેમાંથી મૉડલનું નામ અને `rate_limits` વાંચે છે:

<!-- motion:usage-statusline -->

1. ccam `CCAM_PROFILE` પરથી એકાઉન્ટ શોધે છે, જે ccam Claude Code શરૂ કરતી વખતે સેટ કરે છે, અથવા નહીં તો `CLAUDE_CONFIG_DIR` ને એકાઉન્ટ ફોલ્ડર સાથે મેળવીને (`runStatusline`).
2. ccam `rate_limits.five_hour` અને `rate_limits.seven_day` વાંચે છે (`ParseStatusLine`). API key વાળાં એકાઉન્ટ માટે, અથવા સેશનના પહેલા જવાબ પહેલાં, JSON માં rate limit હોતી નથી; ત્યારે કંઈ સેવ થતું નથી.
3. ccam `~/.local/state/ccam/usage/<name>.statusline.json` લખે છે (`SaveStatusLine`). endpoint જવાબ ન આપે ત્યારે `ccam usage` તેને વાપરે છે.
4. ccam `[name] · model · 5h N% · 7d N%` છાપે છે. `--wrap` સાથે, તે એ જ input સાથે તમારો કમાન્ડ `sh -c` દ્વારા ચલાવે છે અને Claude Code સ્ટેટસ લાઇન બતાવે છે તે રીતે તેનું આઉટપુટ છાપે છે: ખાલી ન હોય તેવી દરેક લાઇન, પહેલીની આગળ `[name]` સાથે. નિષ્ફળ જતો કમાન્ડ ફક્ત એકાઉન્ટનું નામ બતાવે છે. Windows પર તે કમાન્ડ Claude Code ની જેમ ચલાવે છે: Git Bash ઇન્સ્ટૉલ હોય તો તેનાથી, નહીં તો PowerShell થી.
5. જ્યારે 5 કલાકનો કે સાપ્તાહિક વપરાશ 90% કે વધુ હોય અને બીજું એકાઉન્ટ હોય, ત્યારે લાઇનના અંતે `→ ccam next --auto --continue` આવે છે, ઉદાહરણ તરીકે `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

wrap કરેલા કમાન્ડની પોતાની કોઈ સમય મર્યાદા નથી: Claude Code સ્ટેટસ લાઇન કમાન્ડને 10 મિનિટ આપે છે અને નવું update શરૂ થાય ત્યારે ચાલતો run અટકાવે છે, તેથી ધીમો કમાન્ડ ccam વગર જેમ દેખાત તેમ જ દેખાય છે. જ્યારે Claude Code ccam ને બંધ કરે, ત્યારે કમાન્ડ પણ બંધ થાય છે (Linux અને macOS પર તેણે શરૂ કરેલી દરેક process સાથે).

`ccam statusline install` દરેક એકાઉન્ટની `settings.json` માં `statusLine` લખે છે અને object ની બીજી keys, જેમ કે `padding`, `refreshInterval` અને `hideVimModeIndicator`, રાખે છે. જ્યારે કોઈ એકાઉન્ટ પાસે પહેલેથી સ્ટેટસ લાઇન કમાન્ડ હોય, ત્યારે ccam તેને `--wrap '<old command>'` તરીકે રાખે છે; Windows પર `--wrap-base64 <base64>` તરીકે, જેને બંનેમાંથી કોઈ shell માં quoting જોઈતું નથી. જો સ્ટેટસ લાઇન પહેલેથી ccam ચલાવતી હોય (બીજા પાથથી, કે જૂના વર્ઝને લખ્યું હતું તે રીતે quote કરેલું), તો ccam તેને wrap કરવાને બદલે બદલી નાખે છે, અને તે જે કમાન્ડ wrap કરે છે તે રાખે છે.

Windows પર ccam નો પાથ forward slashes અને હોમ ફોલ્ડર માટે `~` સાથે લખાય છે (`~/.local/bin/ccam.exe statusline`), જે Git Bash અને PowerShell બંને સમજે છે. જે પાથને quotes જોઈએ તે Claude Code જે shell વાપરે છે તેના માટે આખો લખાય છે: Git Bash માટે `'…'`, PowerShell માટે `& '…'`. ગ્રુપ દ્વારા શેર થયેલી `settings.json` તેની લિંક દ્વારા એક જ વાર લખાય છે, તેથી ગ્રુપના દરેક એકાઉન્ટને તે મળે છે. તેને ફરી ચલાવવાથી પહેલેથી ccam વાપરતી ફાઇલો માટે `(already set)` છપાય છે.

## સલામતી અને ખાસ કિસ્સા

- તમે `--refresh` ન આપો ત્યાં સુધી `ccam usage` ક્યારેય ટોકન ફાઇલ બદલતું નથી. જ્યારે access token ની મુદત પૂરી થઈ જાય, ત્યારે એકાઉન્ટને `ccam <name>` વડે એક વાર ચલાવો, અને Claude Code તેને રિન્યૂ કરે છે. જે લૉગિનની મુદત પૂરેપૂરી પૂરી થઈ ગઈ હોય તેને `ccam login <name>` જોઈએ.
- `--refresh` ક્યારેય ચાલતાં એકાઉન્ટ (Claude Code ccam વગર શરૂ થયું હોય ત્યારે પણ), ઇમ્પોર્ટ કરેલાં એકાઉન્ટ કે Keychain લૉગિનને અડતું નથી. ઇનકાર થયેલું refresh સેવ કરેલા આંકડા અને કારણ (`token not renewed: …`) બતાવે છે.
- usage endpoint જાહેર રીતે દસ્તાવેજીકૃત નથી અને બદલાઈ શકે છે. જ્યારે તે બદલાય, ત્યારે પણ તમે સેવ કરેલા આંકડા જુઓ છો, અને `ccam doctor` ફૉર્મેટના ફેરફાર વિશે જણાવે છે.
- સેવ કરેલા આંકડા અને જવાબ બંને ન હોય ત્યારે, એકાઉન્ટ લાઇન `lookup failed:` અને કારણ બતાવે છે, ઉદાહરણ તરીકે `--source cache` સાથે `no saved usage`.
- કૅશ ફાઇલોમાં ફક્ત ટકાવારી અને રીસેટ સમય હોય છે, ક્યારેય ટોકન નહીં.

## સંબંધિત

- [વપરાશ મર્યાદા અને એકાઉન્ટ સ્વિચ કરવું](/docs/usage-limits)
- [આગલા એકાઉન્ટ પર સ્વિચ કરવું](/features/next)
- [Claude Code વપરાશ મર્યાદાએ પહોંચે ત્યારે શું કરવું](/guides/usage-limit-workflow)
- [કમાન્ડ સંદર્ભ](/docs/commands)
