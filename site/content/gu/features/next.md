---
title: આગલા એકાઉન્ટ પર સ્વિચ કરવું
description: ccam next બીજા એકાઉન્ટ તરીકે Claude Code ખોલે છે; --auto સૌથી વધુ જગ્યાવાળું લે છે, --continue વાતચીત આગળ લઈ જાય છે, --fallback -p jobs ફરી ચલાવે છે.
updated: 2026-11-01
order: 8
translation: machine
sourceHash: 79300c521748
---

`ccam next` તમે નામ આપ્યા વગર Claude Code ને બીજા એકાઉન્ટ તરીકે ખોલે છે. જ્યારે કોઈ એકાઉન્ટ તેની વપરાશ મર્યાદાએ પહોંચે ત્યારે તેને વાપરો. સ્વિચ કરવાની ચાર રીતો છે:

- સાદું `ccam next` છેલ્લે વાપરેલા એકાઉન્ટ પછીનું એકાઉન્ટ લે છે, અને `--auto` સૌથી વધુ જગ્યાવાળું એકાઉન્ટ લે છે. `--continue` આ ફોલ્ડરની છેલ્લી વાતચીત નવા એકાઉન્ટ પર ચાલુ રાખે છે, બંને ઇતિહાસ શેર કરતાં હોય કે નહીં.
- Claude Code ની અંદર, [ccam plugin](#inside-claude-code) કોઈ મર્યાદા સેશનને અટકાવે ત્યારે તમને જણાવે છે, અને `/ccam:switch` વાતચીત આગળ સોંપે છે.
- `claude -p` ચલાવતી scripts માટે, `ccam --fallback` સ્વિચ જાતે કરે છે (જુઓ [મર્યાદાએ ફરી ચલાવવું](#run-again-at-the-limit)).
- [`auto_switch`](#auto-switch) ચાલુ હોય ત્યારે (Linux અને macOS), Claude Code મર્યાદાએ અટકાવે તે interactive સેશનને પણ ccam ખસેડે છે.

**Team અને Enterprise સીટ.** આપમેળે થતી એકાઉન્ટ પસંદગી (`ccam next`, `--fallback`, `ccam pick`, ડેશબોર્ડની `n` કી અને, તમે તેને ચાલુ કરો પછી, [`auto_switch`](/features/next#auto-switch)) એક જ Claude Team કે Enterprise સંસ્થાની સીટો વચ્ચે પણ સ્વિચ કરે છે, જેનાથી Anthropic ની શરતોનું ઉલ્લંઘન થઈ શકે છે. તેના પર આધાર રાખતાં પહેલાં [વારંવાર પૂછાતા પ્રશ્નો](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) વાંચો અને તમારી સંસ્થા અને Anthropic પાસેથી ખાતરી કરી લો.

## ઉપયોગ

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # continue this folder's last conversation on the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next -c --brief            # a new conversation that starts from a summary of the last one
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the most weekly quota left per day
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work CCAM_PROFILE=work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam સ્વિચને standard error પર `→ <from> → <to>` તરીકે છાપે છે, `--auto` ની પસંદગી હોય તો તેના કારણ સાથે, પછી `ccam <name>` ની જેમ જ Claude Code ને નવા એકાઉન્ટ તરીકે શરૂ કરે છે. `<from>` એ તમે છેલ્લે વાપરેલું એકાઉન્ટ છે, અથવા `--continue` સાથે, જે એકાઉન્ટની વાતચીત છે તે. `--` પછીના આર્ગ્યુમેન્ટ્સ કોઈ ફેરફાર વગર `claude` ને જાય છે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--auto` | વપરાશ પ્રમાણે લૉગિન થયેલું એકાઉન્ટ પસંદ કરે છે (નીચે જુઓ), 100% પરનાં એકાઉન્ટ છોડીને | off |
| `--max-usage <1-100>` | `--auto` સાથે: આ વપરાશથી નીચેનાં એકાઉન્ટને પ્રાથમિકતા આપે છે | `auto_max_usage`, નહીં તો 100 |
| `--strategy headroom\|expiring\|balance` | `--auto` સાથે: સૌથી ઓછો 5 કલાકનો વપરાશ, સૌથી વહેલો 5 કલાકનો રીસેટ કે દરરોજ માટે સૌથી વધુ બાકી સાપ્તાહિક quota | `auto_strategy`, નહીં તો `headroom` |
| `--cross-org` | `--auto` સાથે: બીજી સંસ્થાનાં એકાઉન્ટ પણ પસંદ કરે છે (નીચે જુઓ) | `cross_org`, નહીં તો off |
| `--wait[=<duration>]` | `--auto` સાથે: કોઈ એકાઉન્ટમાં જગ્યા ન હોય ત્યારે, પહેલું એકાઉન્ટ રીસેટ થાય તેની રાહ જુએ છે, વધુમાં વધુ આટલો સમય. મૂલ્ય `=` પછી આવે છે: `--wait 2h` નો `--wait takes its value after '=': --wait=2h` સૂચના સાથે ઇનકાર થાય છે | off; એકલું આપો તો `5h` |
| `-c`, `--continue` | આ ફોલ્ડરની છેલ્લી વાતચીત નવા એકાઉન્ટ પર ચાલુ રાખે છે (નીચે જુઓ) | off |
| `--brief` | `--continue` સાથે: આખી જૂની વાતચીત ફરી શરૂ કરવાને બદલે, ccam ના છેલ્લી વાતચીતના સારાંશથી નવા એકાઉન્ટ તરીકે નવી વાતચીત શરૂ કરે છે ([નીચે જુઓ](#brief)) | off |
| `-- <claude args>` | ccam ના પોતાના આર્ગ્યુમેન્ટ્સ પછી `claude` ને આપવામાં આવતા આર્ગ્યુમેન્ટ્સ | કોઈ નહીં |
| `--dry-run` | સમર્થિત: Claude Code શરૂ કર્યા વગર કે શરૂઆત નોંધ્યા વગર, `→ <from> → <to>` અને જે કમાન્ડ લાઇન ચલાવશે તે છાપે છે | off |

`--auto` અને `--continue` સાથે કામ કરે છે. `--continue` સાથે, `--resume` આર્ગ્યુમેન્ટ્સ તમારા પોતાના આર્ગ્યુમેન્ટ્સ પહેલાં મુકાય છે.

Claude Code ના પોતાના ઓપ્શન, જેમ કે `--dangerously-skip-permissions` કે `--model`, `--` પછી આવે છે. `ccam <name>` નામ પછીનું બધું આગળ મોકલે છે, પણ `ccam next` પોતાના આર્ગ્યુમેન્ટ્સ જાતે વાંચે છે, તેથી `--` વગર તે `ccam: unknown flag: --dangerously-skip-permissions` સાથે અટકી જાય છે.

## તે કેવી રીતે કામ કરે છે

ccam દ્વારા થતી દરેક શરૂઆત સમયને `~/.local/state/ccam/last-used/<name>` માં લખે છે. `ccam next` આ ફાઇલો વાંચે છે, કયાં એકાઉન્ટ લૉગિન થયેલાં છે તે તપાસે છે, અને પસંદ કરે છે.

1. ccam `config.toml` પરથી દરેક એકાઉન્ટ માટે એક `next.Candidate` બનાવે છે (`cli.candidates`): એકાઉન્ટ, તેની છેલ્લી શરૂઆતનો સમય (`launcher.LastUsed`) અને તેની પાસે કામ કરતું લૉગિન છે કે નહીં (`creds.Has`: `.credentials.json` માં, કે macOS પર Keychain માં, એવું Claude લૉગિન જેની મુદત પૂરેપૂરી પૂરી ન થઈ હોય અને જેને Claude Code એ સાફ ન કર્યું હોય). તેના બદલે API key વાપરતું ફોલ્ડર (Console લૉગિન, કે `ccam add --api` વડે ઉમેરેલું એકાઉન્ટ) પ્લાનની કોઈ મર્યાદા ધરાવતું નથી અને ક્યારેય પસંદ થતું નથી.
2. `--auto` વગર, ccam એકાઉન્ટને સ્લૉટ નંબર પ્રમાણે ગોઠવે છે અને સૌથી તાજેતરમાં શરૂ થયેલું એકાઉન્ટ શોધે છે (`next.Rotate`). એ જ એકાઉન્ટ છે જેમાંથી તમે સ્વિચ કરી રહ્યા છો. જો હજી કોઈ એકાઉન્ટ શરૂ ન થયું હોય, તો તે સૌથી નાનો સ્લૉટ છે.
3. ત્યાંથી તે સ્લૉટના ક્રમમાં આગળ વધે છે, છેલ્લા સ્લૉટ પછી ફરી પહેલા પર આવે છે, અને લૉગિન થયેલું પહેલું એકાઉન્ટ પાછું આપે છે.
4. ccam પસંદ કરેલું એકાઉન્ટ બરાબર `ccam <name>` ની જેમ તૈયાર કરે છે અને શરૂ કરે છે (`launcher.LaunchWith`): તે એકાઉન્ટની લિંક તપાસે છે, MCP શેરિંગ લાગુ કરે છે, સેશન નોંધે છે અને `syscall.Exec` દ્વારા ccam ને `claude` થી બદલે છે.

આ શોધ છેલ્લે વાપરેલા એકાઉન્ટથી શરૂ થાય છે અને તેના પછીના લૉગિન થયેલા પહેલા એકાઉન્ટ પર અટકે છે:

<!-- motion:next-rotate -->

`work` સૌથી તાજેતરમાં શરૂ થયું હતું, તેથી શોધ તેના પછીથી શરૂ થાય છે. `side` પાસે લૉગિન નથી અને તે છોડી દેવાય છે, અને `lab` પસંદ થાય છે. જો `lab` પણ લૉગઆઉટ હોત, તો શોધ ફરીને `personal` પર આવત (ડાયાગ્રામમાં `% len`). rotation ક્યારેય તમે જ્યાંથી શરૂ કર્યું તે એકાઉન્ટ પાછું આપતું નથી: લૉગિન થયેલું બીજું કોઈ એકાઉન્ટ ન હોય, ત્યારે તે ભૂલ સાથે અટકે છે.

### --auto

`--auto` સ્લૉટના ક્રમને બદલે વપરાશ પરથી નક્કી કરે છે. ccam લૉગિન થયેલાં એકાઉન્ટ રાખે છે, `ccam usage` ની જેમ જ તેમનો બધો વપરાશ એકસાથે પૂછે છે (`Fetcher.FetchAll`), અને નીચેના નિયમોથી પસંદ કરે છે (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

એકાઉન્ટને તેની 5 કલાકની અને સાપ્તાહિક વિન્ડો પરથી મૂલવવામાં આવે છે, અને જ્યારે Claude Code માટેના આર્ગ્યુમેન્ટ્સમાં `--model` હોય (`--` પછી) ત્યારે એ મૉડલની સાપ્તાહિક વિન્ડો પરથી પણ (નામમાં `opus`, `sonnet`, `haiku` કે `fable`). ccam એકાઉન્ટને ત્રણ સ્તરમાં વહેંચે છે (`Policy.Choose()`):

1. **જગ્યા**: દરેક વિન્ડો `auto_max_usage` (`config.toml`, ડિફૉલ્ટ 100) થી નીચે. તેમાંથી strategy પસંદ કરે છે: `headroom` (ડિફૉલ્ટ) સૌથી ઓછો 5 કલાકનો વપરાશ લે છે; `expiring` સૌથી વહેલી રીસેટ થતી 5 કલાકની વિન્ડો લે છે, જેથી તેનો બાકીનો quota ગુમાવાય તે પહેલાં વપરાઈ જાય; `balance` એવું એકાઉન્ટ લે છે જેની પાસે તેનું અઠવાડિયું રીસેટ થાય ત્યાં સુધી દરરોજ માટે સૌથી વધુ સાપ્તાહિક quota બાકી હોય, ચાલુ દરેક ccam સેશનને તેની સામે ગણીને, જેથી સાપ્તાહિક quota વણવપરાયેલો ન રહે અને સમાંતર સેશન વહેંચાઈ જાય.
2. **ઉપર**: 100% થી નીચે પણ `auto_max_usage` પર કે તેથી વધુ. કોઈ એકાઉન્ટમાં જગ્યા ન હોય ત્યારે જ વપરાય છે; સૌથી ઓછા 5 કલાકના વપરાશવાળું પહેલાં આવે છે.
3. **ભરેલું**: 100% પરની કોઈ વિન્ડો, અથવા અજાણ્યો 5 કલાકનો વપરાશ. ક્યારેય પસંદ થતું નથી.

જ્યારે Claude Code એ કોઈ વિન્ડોમાં request નકારી હોય, ત્યારે તે વિન્ડો તેના રીસેટ સુધી 100% ગણાય છે: [ccam plugin](#inside-claude-code) નો hook તેની નોંધ કરે છે (`~/.local/state/ccam/usage/<name>.limit.json`), સેવ કરેલા આંકડા ભલે ગમે તે કહે. તેથી હમણાં જ મર્યાદાએ પહોંચેલું એકાઉન્ટ તેના રીસેટ પહેલાં ફરી પસંદ થતું નથી, તેના આંકડા જૂના હોય ત્યારે પણ.

બધાં એકાઉન્ટ માટેનાં સેટિંગ્સ `config.toml` માં સૌથી ઉપર, પહેલા `[[group]]`, `[[profile]]` કે `[[folder]]` ની ઉપર આવે છે. તેમાંના કોઈની નીચેનું સેટિંગ તે table નું ગણાય છે: ccam દરેક શરૂઆતે તેના વિશે ચેતવણી આપે છે, અને config સેવ કરતા કમાન્ડ તેને ખસેડવામાં ન આવે ત્યાં સુધી ઇનકાર કરે છે.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**સંસ્થાઓ.** આપમેળે થતી પસંદગી તમે છોડો છો તે એકાઉન્ટ જેવા જ પ્રકારનાં એકાઉન્ટમાં રહે છે: એ જ Team કે Enterprise સંસ્થાની સીટો, અથવા અંગત પ્લાન (Pro, Max), જેથી કંપનીના એકાઉન્ટનું કામ ખાનગી એકાઉન્ટ પર ચાલુ ન થાય, કે ઊલટું. ccam દરેક એકાઉન્ટની સંસ્થા તેના લૉગિન પરથી જાણે છે (`ccam login` તેની નોંધ કરે છે; જૂનાં લૉગિન `.claude.json` માંથી વંચાય છે); જે એકાઉન્ટની સંસ્થા ccam કહી ન શકે તે કોઈ પણ એકાઉન્ટ સાથે મેળ ખાઈ શકે છે. એ જ Claude એકાઉન્ટમાં લૉગિન થયેલાં બે ફોલ્ડર તેની મર્યાદા શેર કરે છે, તેથી એક ક્યારેય બીજાનું કામ સંભાળતું નથી. `--cross-org`, અથવા `config.toml` માં `cross_org = true`, સંસ્થાનો નિયમ હટાવે છે, પણ [`/ccam:switch`](#inside-claude-code) પછી બીજી સંસ્થામાં ખસેડ હજી પણ તમારી હા ની રાહ જુએ છે.

કંઈ પસંદ ન થઈ શકે ત્યારે ccam code 75 (`EX_TEMPFAIL`) સાથે બહાર નીકળે છે, જેથી script "પછી પ્રયાસ કરો" ને નિષ્ફળતાથી અલગ ઓળખી શકે. `--wait` સાથે, તેના બદલે તે પહેલા ભરેલા એકાઉન્ટની અટકાવતી વિન્ડો રીસેટ થાય ત્યાં સુધી ઊંઘે છે (રીસેટ સમયના 30 સેકન્ડ પછી, જ્યારે endpoint નવી વિન્ડો બતાવે), `ccam: no account has room; waiting until 15:10` છાપે છે, અને ફરી પસંદ કરે છે; જ્યારે કોઈ એકાઉન્ટ ક્યારે રીસેટ થશે તે ન કહે, ત્યારે તે 10 મિનિટ પછી ફરી જુએ છે. રાહ પૂરી થાય તે પહેલાં કંઈ ખાલી ન થાય, તો તે 75 સાથે બહાર નીકળે છે.

lookup `ccam usage` ના નિયમો અનુસરે છે: 60 સેકન્ડથી ઓછા જૂના આંકડા ફરી વપરાય છે, દરેક એકાઉન્ટનું lookup 3 સેકન્ડ પછી timeout થાય છે, જે એકાઉન્ટ જવાબ ન આપે કે જેના access token ની મુદત પૂરી થઈ ગઈ હોય, તેને તેના સેવ કરેલા આંકડા પરથી મૂલવવામાં આવે છે, અને મર્યાદા રેકોર્ડ તેમના પર અગ્રતા લે છે. જે સેવ કરેલી વિન્ડોનો રીસેટ સમય વીતી ગયો હોય તે 0% ગણાય છે, તેથી ત્યાર પછી જેની વિન્ડો રીસેટ થઈ હોય તે એકાઉન્ટ છોડી દેવાતું નથી. `--auto` ક્યારેય ટોકન રિન્યૂ કરતું નથી. સાદા rotation થી વિપરીત, `--auto` છેલ્લે વાપરેલું એકાઉન્ટ પણ પસંદ કરી શકે છે, જો તેની પાસે હજી સૌથી વધુ જગ્યા હોય.

### --continue

1. ccam દરેક એકાઉન્ટમાં આ ફોલ્ડરની છેલ્લી વાતચીત શોધે છે: `projects/<folder>/` માંની `<session>.jsonl` ફાઇલોમાંથી, જેમાં સૌથી નવો prompt કે જવાબ હોય તે (ફોલ્ડરનું નામ એ લિંક ઉકેલેલું કામનું ફોલ્ડર છે, જેમાં અક્ષર કે અંક સિવાયનો દરેક અક્ષર `-` માં ફેરવાય છે, જેમ Claude Code તેને નામ આપે છે). જેમાં એવું કંઈ ન હોય તેવાં સેશન, જેમ કે ફક્ત `/mcp`, `/usage` કે `/status` ચલાવેલું સેશન, કે ખાલી ફાઇલ, છોડી દેવાય છે. જો કોઈ સેશને નોંધ છોડી હોય (વપરાશ મર્યાદા કે `/ccam:switch`, નીચે જુઓ) અને કોઈ વાતચીતમાં તેના 2 મિનિટથી વધુ પછીનો prompt કે જવાબ ન હોય, તો તેના બદલે બરાબર એ જ સેશન વપરાય છે; `/exit` વખતે Claude Code જે હરોળો લખે છે તે ગણાતી નથી.
2. તે જે એકાઉન્ટની છે તે જ એ એકાઉન્ટ છે જેમાંથી તમે સ્વિચ કરો છો, અને નવું એકાઉન્ટ ઉપર પ્રમાણે, એ જ સંસ્થાની અંદર પસંદ થાય છે. જ્યારે ઘણાં એકાઉન્ટ એ ઇતિહાસ શેર કરતાં હોય, ત્યારે તે એ એકાઉન્ટની ગણાય છે જેની `.claude.json` કહે છે કે તેણે વાતચીત છેલ્લે ચલાવી હતી (`lastSessionId`), નહીં તો છેલ્લે વાપરેલા એકાઉન્ટની.
3. જો નવું એકાઉન્ટ એ ઇતિહાસ શેર કરતું હોય (એ જ ગ્રુપ ફોલ્ડર), તો ccam `claude --resume <session id>` શરૂ કરે છે, એટલે કે બરાબર એ જ વાતચીત (`--continue` વાપરવાથી ફોલ્ડરમાંની સૌથી નવી વાતચીત ખૂલે).
4. નહીં તો તે `claude --resume <the conversation's file> --fork-session` શરૂ કરે છે: Claude Code આખી વાતચીતની નકલ નવા session ID હેઠળ નવા એકાઉન્ટના પોતાના ઇતિહાસમાં કરે છે, અને મૂળ જેમ હતી તેમ રહે છે. ત્યાર પછી નવા એકાઉન્ટનું `--continue` તેને શોધી લે છે.
5. જો તમે `--` પછી `--agent` ન આપો, તો વાતચીત જે agent તરીકે ચાલે છે તે agent પણ ccam આગળ આપે છે: ફોલ્ડરની `.claude/settings.local.json` કે `.claude/settings.json` નું, અથવા એકાઉન્ટની `settings.json` નું `agent` સેટિંગ, નહીં તો વાતચીતની છેલ્લી `agent-setting` હરોળ જે agent નું નામ આપે તે. વાતચીતની નકલ કરતી વખતે (પગલું 4) Claude Code agent ને પાછો લાવતું નથી, અને તે પોતાના ડિફૉલ્ટ agent તરીકે, દરેક tool સાથે, આગળ ચાલત. જ્યારે નવું એકાઉન્ટ જૂના એકાઉન્ટના agents, plugins અને સેટિંગ્સ શેર ન કરતું હોય (ગ્રુપ દ્વારા એ જ `agents` અને `plugins` ફોલ્ડર અને `settings.json`), ત્યારે ccam ખસેડ વિશેની લાઇન પછી એક લાઇન ઉમેરે છે: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

ccam જણાવે છે કે વાતચીત કેટલી મોટી છે: નવું એકાઉન્ટ તેના પહેલા જવાબ સાથે તે આખી ફરી વાંચે છે. prompt caches ક્યારેય સંસ્થાઓ વચ્ચે શેર થતા નથી, અને દરેક અંગત પ્લાન પોતે એક અલગ સંસ્થા છે, તેથી જ્યારે બંને એકાઉન્ટ જુદી સંસ્થાનાં હોય ત્યારે ccam `prompt caches are not shared between organizations, so the first reply reads it all again` ઉમેરે છે. જ્યારે જૂનું એકાઉન્ટ એવા MCP સર્વરમાં સાઇન ઇન હોય જેમાં નવું નથી, ત્યારે ccam તેમનાં નામ આપે છે: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. ગ્રુપમાં `history` શેર કરવાથી નકલ ટળે છે, એક જ વાતચીત રહે છે, અને `/rewind` જેમાંથી કોડ પાછો લાવે છે તે checkpoints (`file-history/`) પણ શેર થાય છે.

### --brief

`--continue` આખી વાતચીત ફરી શરૂ કરે છે, અને નવું એકાઉન્ટ તેના પહેલા જવાબ સાથે તે આખી ફરી વાંચે છે. `ccam next --continue --brief` તેના બદલે નવી વાતચીત શરૂ કરે છે, એક સારાંશથી જે ccam કોઈ મૉડલને પૂછ્યા વગર જૂની વાતચીતમાંથી લખે છે:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

સારાંશમાં, આ ક્રમમાં અને ફક્ત જ્યાં કહેવા જેવું કંઈ હોય ત્યાં, આ હોય છે:

1. વાતચીતનું શીર્ષક;
2. તેના transcript નો પાથ, જેથી Claude ને કોઈ વિગત જોઈએ ત્યારે તેમાં શોધી શકે;
3. તમારી પહેલી વિનંતી, જ્યારે ચારથી વધુ હોય;
4. Claude Code એ છેલ્લે વાતચીત compact કરી ત્યારે લખેલો સારાંશ;
5. તમારી છેલ્લી ચાર વિનંતીઓ, `!` commands અને slash commands તમે લખ્યા હતા તે રીતે;
6. Claude નો છેલ્લો જવાબ;
7. todo યાદી;
8. Claude એ છેલ્લે બદલેલી 10 ફાઇલો;
9. ફોલ્ડરનું `git status --short --branch`, 800 અક્ષરો સુધી.

todo યાદી વધુમાં વધુ 10 આઇટમ રાખે છે: પહેલાં પૂરી થયેલી આઇટમ કાઢી નખાય છે, પછી સૌથી આગળનાં પગલાં, જેથી હાલનું પગલું અને પછીનાં પગલાં રહે. લાંબા ભાગ વચ્ચેથી કપાય છે, અને આખો સારાંશ 16,000 અક્ષરોથી ઓછો રહે છે.

ccam તેને `~/.local/state/ccam/brief/<session>.md` તરીકે સેવ કરે છે, જે ફક્ત તમે વાંચી શકો, અને આગલું `--brief` 7 દિવસથી જૂના સારાંશ કાઢી નાખે છે. Claude Code એક prompt સાથે શરૂ થાય છે, અને તેમાં ઉલ્લેખેલી ફાઇલ જોડે છે: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` `--` પછી તમે આપેલો prompt છેલ્લા વાક્યની જગ્યા લે છે (`Then: <your prompt>`), અને `--` પછીના તમારા બીજા ઓપ્શન તેની પાછળ આવે છે.

Claude Code ઉલ્લેખને `"` પર પૂરો કરે છે અને `#` પર કાપે છે, તેથી જ્યારે સારાંશના પાથમાં આમાંથી કોઈ હોય, ત્યારે ccam કામના ફોલ્ડરની સાપેક્ષે તેના પાથ વડે તેનો ઉલ્લેખ કરે છે; જ્યારે એ પાથમાં પણ એક હોય, ત્યારે prompt ફાઇલનું નામ આપે છે, અને Claude તેને પોતાના Read tool વડે વાંચે છે, જે પહેલાં પૂછે છે કારણ કે ફાઇલ કામનાં ફોલ્ડરની બહાર છે (Manual mode માં `-p` run માં, Claude Code વાંચવાની મનાઈ કરે છે અને Claude સારાંશ વગર આગળ વધે છે).

`ccam --dry-run next --continue --brief` સારાંશ છાપે છે અને કંઈ સેવ કરતું નથી. જૂની વાતચીત જેમ હતી તેમ રહે છે, અને નવી વાતચીત ફક્ત એટલું જ જાણે છે જેટલું સારાંશ કહે છે, તેથી જ્યારે તમે ત્યાં કેવી રીતે પહોંચ્યા તેની દરેક વિગત કરતાં કામ વધુ મહત્ત્વનું હોય ત્યારે `--brief` વાપરો. `--continue` વગર, ccam `--brief works with --continue` સાથે અટકે છે; જ્યારે તે વાતચીત વાંચી ન શકે કે તેમાં સારાંશ આપવા જેવું કંઈ ન મળે, ત્યારે તે અટકે છે અને કહે છે કે `--brief` વગરનું `ccam next --continue` તેને આખી ફરી શરૂ કરે છે. [relay](#relay) અને [`auto_switch`](#auto-switch) હંમેશાં આખી વાતચીત ફરી શરૂ કરે છે. Claude Code નો પોતાનો `--brief` ઓપ્શન, તેના બીજા ઓપ્શનની જેમ, `--` પછી આવે છે.

## Claude Code ની અંદર

[ccam plugin](/docs/install#claude-code-plugin) બે hooks ઉમેરે છે (Claude Code 2.1.78 કે પછીનું):

- જ્યારે કોઈ turn વપરાશ મર્યાદાએ અટકે, ત્યારે તે વાતચીતની નોંધ કરે છે અને તેને ટેકો આપતાં ટર્મિનલ (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty) માં desktop સૂચના બતાવે છે: ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  Claude Code જે દરેક સંદેશને વપરાશ મર્યાદા ગણે છે તે તે ઓળખે છે: 5 કલાકની, સાપ્તાહિક, સાદી અને મૉડલની મર્યાદાઓ (`You've hit your session limit`, `You've hit your limit`, મૉડલના ભાગ માટે `You've reached your …`), usage credits અને spend limits, અને એકાઉન્ટને અટકાવનારી સંસ્થા.

  જ્યારે Claude Code એ નકારેલી request વાતચીતમાં લખી હોય, ત્યારે hook ભરાયેલી વિન્ડોની પણ તેના રીસેટ સુધી નોંધ કરે છે ([મર્યાદા રેકોર્ડ](/features/usage#the-limit-record)), સંસ્થાએ લગાવેલા અટકાવ સિવાય, કારણ કે તેના admins તે ગમે ત્યારે હટાવી શકે છે. સર્વરનું ટૂંકું throttling કે API key ની rate limit વપરાશ મર્યાદા નથી, અને તેને એમ જ છોડી દેવાય છે, અને subagent નું અટકવું તથા કોઈ સામે બેઠું ન હોય તેવાં સેશન (`-p` runs અને બૅકગ્રાઉન્ડ સેશન) પણ.

  Claude Code 2.1.234 અને પછીનાં વર્ઝન સેશનમાં રાહ પણ જુએ છે અને 5 કલાકના રીસેટ પછી જાતે ચાલુ રાખે છે; સૂચના તમને તેના બદલે સ્વિચ કરવાનો વિકલ્પ આપે છે, અને [`auto_switch`](#auto-switch) ચાલુ હોય ત્યારે (Linux અને macOS), બીજા એકાઉન્ટમાં જગ્યા હોય તો ccam તમારા માટે સ્વિચ કરે છે. જ્યારે બીજું કોઈ એકાઉન્ટ કામ સંભાળી ન શકે, કારણ કે બીજું કોઈ એકાઉન્ટ એ જ Team કે Enterprise સંસ્થાનું નથી, કે બીજું કોઈ અંગત પ્લાન નથી (`cross_org` ચાલુ હોય ત્યારે, બીજું કોઈ પણ Claude એકાઉન્ટ સંભાળી શકે), ત્યારે સૂચના તેના બદલે `You have no other personal plans, so ccam cannot move this conversation.` કે `You have no other accounts in Acme, so ccam cannot move this conversation.` થી પૂરી થાય છે.
- `/ccam:switch [account]` વાતચીતની નોંધ કરે છે અને કોઈ request થાય તે પહેલાં turn પૂરો કરે છે: `/exit` લખો, પછી તેને તમે નામ આપેલા એકાઉન્ટ પર (કે આગલા એકાઉન્ટ પર) ચાલુ રાખવા `ccam next --continue` ચલાવો. જે સેશનમાં હજી કોઈ વાતચીત ન હોય, તેમાં તે જણાવે છે કે ચાલુ રાખવા જેવું કંઈ નથી. ફક્ત તમે જે લખો તે જ તેને ચલાવે છે: plugin 0.3.2 થી Claude `/ccam:switch` ચલાવી શકતું નથી, અને Claude એ શેડ્યૂલ કરેલું કામ (CronCreate, ScheduleWakeup કે `.claude/scheduled_tasks.json`) તેને Claude ને સાદા લખાણ તરીકે આપે છે.

  જ્યારે તમે નામ આપેલું એકાઉન્ટ વાતચીતની સંસ્થાની બહાર હોય (બીજી Team કે Enterprise સંસ્થાની સીટ, સીટને બદલે અંગત પ્લાન કે ઊલટું, કે કોઈ પણ બાજુએ API key એકાઉન્ટ), ત્યારે `cross_org` ગમે તે કહે, `ccam next --continue` પહેલાં પૂછે છે: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. તમે કોઈ નામ ન આપો ત્યારે તે જાતે પસંદ કરેલા એકાઉન્ટ પહેલાં પણ એ જ પૂછે છે, જે `cross_org = true` હોય તો તે બીજી સંસ્થામાંથી લઈ શકે છે: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. ફક્ત `y` કે `yes` તેને ખસેડે છે; Linux અને macOS પર, પ્રશ્ન પહેલાં દબાવેલી keys પહેલાં કાઢી નખાય છે. બીજો કોઈ પણ જવાબ, કે પૂછવા માટે ટર્મિનલ ન હોય તો, વાતચીત અને નોંધ જ્યાં છે ત્યાં જ રહે છે: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (જાતે પસંદ કર્યા પછી, `ccam next --continue` ફરી પસંદ કરે છે). આપમેળે થતી પસંદગીની જેમ, જે એકાઉન્ટની સંસ્થા ccam કહી ન શકે તે કોઈ પણ સાથે મેળ ખાય છે.

નોંધ `~/.local/state/ccam/handoff/<account>/<session>.json` માં રહે છે, દરેક સેશન માટે એક ફાઇલ, તેથી એક એકાઉન્ટનાં બે સેશન ક્યારેય એકબીજાની નોંધ પર લખતાં નથી. તે 12 કલાક રખાય છે અને આગલું સેશન શરૂ થવાનું હોય ત્યારે એક જ વાર વપરાય છે; નિષ્ફળ જતી શરૂઆત તેને આગલા પ્રયાસ માટે રહેવા દે છે. તે જે ફોલ્ડરમાં સેશન શરૂ થયું હતું તે ફોલ્ડરની જ રહે છે, Claude Code ના shell માં `cd` કર્યા પછી પણ. `ccam rename` તેમાંનાં નામ અપડેટ કરે છે; જો તેમાં લખેલું એકાઉન્ટ ત્યાર પછી કાઢી નખાયું હોય, તો ccam તે જણાવે છે અને હંમેશની જેમ પસંદ કરે છે.

### Relay

`config.toml` માં સૌથી ઉપર `relay = true` સાથે, ccam પોતાની જગ્યા Claude Code ને આપવાને બદલે interactive સેશનની સાથે ચાલતું રહે છે, અને `/ccam:switch` પછી એ સેશનની વાતચીત જાતે ચાલુ રાખે છે (બીજા સેશને છોડેલી નોંધ પર તે કંઈ કરતું નથી): hook Claude Code ને બંધ કરે છે, અને ccam એ સેશનની નોંધ સાથે `ccam next --continue` ચલાવે છે. એકાઉન્ટનું નામ ન આપ્યું હોય, ત્યારે તે વાતચીતના મૉડલને ગણીને સૌથી વધુ જગ્યાવાળું એકાઉન્ટ લે છે, અને જ્યારે `cross_org = true` તેને બીજી સંસ્થાનું એકાઉન્ટ બનવા દે, ત્યારે વાતચીતને ત્યાં ખસેડતાં પહેલાં તે ઉપરનો પ્રશ્ન પૂછે છે. Windows પર hook Claude Code ને બંધ કરી શકતો નથી: `/exit` લખો, અને ccam તેને ચાલુ રાખે છે. વાતચીતની સંસ્થાની બહારના એકાઉન્ટ માટે પણ hook તેને બંધ કરતો નથી: તે `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).` કહે છે, અને તમે `/exit` લખો તે પછી ccam ઉપરનો પ્રશ્ન પૂછે છે. તેનો ખર્ચ દરેક સેશન દીઠ એક રાહ જોતી ccam process છે, તેથી તે ડિફૉલ્ટ રીતે બંધ છે.

relay થી આગળ ચાલતું સેશન એ રાખે છે જે વાતચીત ફરી શરૂ કરવાથી પાછું આવતું નથી: `--add-dir`, `--mcp-config`, tool યાદીઓ, `--settings`, `--plugin-dir` અને system prompt ઓપ્શન જેવા ઓપ્શન, વાતચીતનો permission mode, સેશનનું effort level અને વાતચીત જે agent તરીકે ચાલે છે તે agent ([`--continue`](#continue) જુઓ). Claude Code પોતે વાતચીતમાંથી `--model` અને `--name` પાછા લાવે છે. બીજું કંઈ પણ (`--worktree`, `--debug`, `--remote-control`, ccam ન જાણતું હોય તેવા ઓપ્શન) `ccam: not carried over: …` માં જણાવાય છે. સ્વિચ ક્યારેય શેર કરેલા ઇતિહાસનો પ્રશ્ન પૂછતું નથી, અને તે ચાલતું હતું ત્યારે તમે જે લખ્યું તે છોડી દેવાય છે; તેનો એકમાત્ર પ્રશ્ન બીજી સંસ્થાના એકાઉન્ટ પહેલાંનો પ્રશ્ન છે.

relay ત્યારે જ વાતચીત ચાલુ રાખે છે જ્યારે Claude Code `/ccam:switch` એ કહ્યું હતું તે રીતે બંધ થયું હોય: તમે પહેલાં `/exit` લખ્યું હોય, તે બંધ થતું હતું ત્યારે Ctrl+C દબાવ્યું હોય, કે ccam ને hangup કે બીજો signal મળ્યો હોય ત્યારે નહીં, અને ટર્મિનલ ન રહ્યું હોય ત્યારે પણ નહીં. ત્યારે નોંધ રહે છે, અને `ccam next --continue` પછીથી તેને ચાલુ રાખે છે. ccam ના જૂના વર્ઝનના relay હેઠળ શરૂ થયેલું સેશન, કે કોઈ સામે બેઠું ન હોય તેવું સેશન (`-p` runs અને બૅકગ્રાઉન્ડ સેશન), hook વડે બંધ થતું નથી; તેને તેના બદલે `/exit` નો સંદેશ મળે છે. `claude remote-control` relay વગર ચાલે છે.

### `auto_switch`

`config.toml` માં સૌથી ઉપર `auto_switch = true` સાથે (Linux અને macOS), જ્યારે Claude Code કોઈ સેશનને વપરાશ મર્યાદાએ અટકાવે ત્યારે ccam તેને જાતે બીજા એકાઉન્ટ પર ખસેડે છે: plugin નો hook જગ્યા ધરાવતું બીજું એકાઉન્ટ પસંદ કરે છે અને Claude Code બંધ કરે છે, અને સેશનની સાથે ચાલતું ccam વાતચીત ત્યાં ચાલુ રાખે છે, જે રીતે `/ccam:switch` પછી [relay](#relay) કરે છે. તે મૂળભૂત રીતે બંધ છે અને પ્રાયોગિક છે.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam સેશન ત્યારે જ ખસેડે છે જ્યારે આ બધી શરતો સાચી હોય:

1. ccam એ તેને ટર્મિનલમાં શરૂ કર્યું હોય અને તેની સાથે ચાલતું રહેતું હોય. `auto_switch` [relay](#relay) ચાલુ કરે છે, તેથી તે સેટ કર્યા પછી તમે શરૂ કરો તે interactive સેશન માટે આ સાચું છે. પહેલાં શરૂ થયેલા, કે જૂના ccam હેઠળ શરૂ થયેલા, સેશનને એવી સૂચના મળે છે જે ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` થી પૂરી થાય છે. કોઈ સામે બેઠું ન હોય તેવાં સેશન (`-p` runs અને બૅકગ્રાઉન્ડ સેશન), `claude remote-control` અને ccam વગર શરૂ થયેલું Claude Code ક્યારેય ખસેડાતાં નથી.
2. એકાઉન્ટમાં [ccam plugin](/docs/install#claude-code-plugin) ચાલુ હોય અને તે Claude Code 2.1.214 કે પછીનું ચલાવતું હોય, જે hook ને તેની process ID આપે છે (`CLAUDE_PID`).
3. Claude Code એ એ એકાઉન્ટની 5 કલાકની, સાપ્તાહિક, Opus કે Sonnet મર્યાદા માટે, રીસેટના સમય સાથે, નકારેલી request વાતચીતમાં લખી હોય. spend limit, usage credits ખૂટી જવી, સંસ્થાએ લગાવેલો અટકાવ, Fable નો ભાગ (`You've reached your Fable limit.`), સર્વરનું throttling, API key ની rate limit અને subagent નું અટકવું ક્યારેય સેશન ખસેડતાં નથી.
4. ccam એકાઉન્ટનો પ્લાન જાણતું હોય, તેની `.claude.json` પરથી કે `ccam login` એ નોંધ્યું હોય તે પરથી. નહીં તો સૂચના `ccam does not know personal's plan, so auto_switch leaves this session.` થી પૂરી થાય છે.
5. સેશન Remote Control સાથે જોડાયેલું ન હોય. બીજા એકાઉન્ટ તરીકે તે જે app થી તમે તેને ચલાવો છો તેની સાથે ફરી જોડાય નહીં, તેથી તે રહે છે, અને સૂચના તે જણાવે છે.
6. મર્યાદા હમણાંથી 15 મિનિટથી વધુ સમય પછી રીસેટ થતી હોય. ઓછી રાહ માટે ccam સેશન Claude Code પર છોડે છે: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. મર્યાદા પછી તમે કોઈ prompt લખ્યો ન હોય.
8. બીજા એકાઉન્ટમાં જગ્યા હોય. ccam [`--auto`](#auto) ની જેમ પસંદ કરે છે (`auto_max_usage`, `auto_strategy`, વાતચીતનું મૉડલ), તમે જે છોડો છો તેના જેવા જ પ્રકારનાં એકાઉન્ટમાંથી: એ જ Team કે Enterprise સંસ્થાની બીજી સીટ, કે બીજો અંગત પ્લાન (Pro, Max). `cross_org` અને `--cross-org` અહીં લાગુ પડતાં નથી. એકાઉન્ટ:
   - Claude Code રિન્યૂ કરી શકે તેવા લૉગિનથી લૉગિન થયેલું હોવું જોઈએ;
   - એ જ Claude એકાઉન્ટમાં લૉગિન થયેલું ન હોવું જોઈએ;
   - API key વાપરતું ન હોવું જોઈએ;
   - ccam જાણતું હોય તેવો પ્લાન ધરાવતું હોવું જોઈએ;
   - ccam ના પોતાના વપરાશ lookup માં દેખાય ત્યાં સુધી, તેમાં usage credits ચાલુ ન હોવી જોઈએ;
   - Claude Code માં આ ફોલ્ડર પર ભરોસો કરી ચૂક્યું હોવું જોઈએ;
   - જો વાતચીત કોઈ agent તરીકે ચાલતી હોય (`agent` સેટિંગ કે વાતચીત કોઈ agent નું નામ આપે), તો તમે છોડો છો તે એકાઉન્ટના agents શેર કરતું હોવું જોઈએ: તેનાં `agents` અને `plugins` ફોલ્ડર અને તેની `settings.json` દરેક બરાબર એ જ હોવાં જોઈએ (ગ્રુપ દ્વારા), અથવા બંનેમાં ન હોવાં જોઈએ.

   જ્યારે આ છેલ્લી ત્રણ શરતો જગ્યા ધરાવતું દરેક એકાઉન્ટ બાકાત કરે, ત્યારે સૂચના તેમનાં નામ આપે છે: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` કે `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` [મર્યાદા રેકોર્ડ](/features/usage#the-limit-record) વાળી વિન્ડો ભરાયેલી ગણાય છે, અને છેલ્લી 2 મિનિટમાં બીજું સેશન જે એકાઉન્ટમાં ગયું હોય તે બીજાં પછી આવે છે. પસંદગી લગભગ 6 સેકન્ડમાં થવી જોઈએ; વપરાશ lookup તેનાથી વધુ સમય લે, ત્યારે સેશન રહે છે. `usage_source = "claude"` સાથે, `auto_switch` ફક્ત સેવ કરેલા આંકડા પરથી નક્કી કરે છે, અને usage credits ચાલુ છે કે નહીં તે ccam Claude Code ના lookup કે સ્ટેટસ લાઇન પરથી કહી શકતું નથી, તેથી તે ધારે છે કે તે બંધ છે.

પછી hook Claude Code બંધ કરે છે, અને તેની સાથે ચાલતું ccam, તેને ટેકો આપતાં ટર્મિનલમાં desktop સૂચના બતાવે છે, જેમ કે `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` તે ઉપરના ઉદાહરણ output ની પહેલી લાઇન છાપે છે અને એ સેશન માટે પસંદ કરેલા એકાઉન્ટ તરીકે `ccam next --continue` ચલાવે છે, જે તેની બાકીની બે લાઇન છાપે છે અને [relay](#relay) ની જેમ સેશનના ઓપ્શન, permission mode, effort level અને agent સાથે Claude Code શરૂ કરે છે. નવા સેશનનો પહેલો સંદેશ આ નક્કી કરેલું લખાણ છે, ccam તમારા વતી મોકલે તે એકમાત્ર વસ્તુ (`/ccam:switch` પછીનો ખસેડ કંઈ મોકલતો નથી):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

જો એ એકાઉન્ટ પણ તેની મર્યાદાએ પહોંચે, તો ત્યાં પણ એવું જ થાય છે. મર્યાદા રેકોર્ડ આ સાંકળને કોઈ એકાઉન્ટ પર તેના રીસેટ પહેલાં પાછી જતી રોકે છે, તેથી જ્યારે દરેક એકાઉન્ટ ભરાઈ જાય ત્યારે સેશન છેલ્લા એકાઉન્ટ પર રહે છે.

વાતચીત કોઈ agent તરીકે ચાલે છે તે hook ને `agent` સેટિંગ પરથી, અથવા Claude Code સમયાંતરે અને સેશન પૂરું થાય ત્યારે વાતચીતમાં લખે છે તે agent નું નામ આપતી હરોળ પરથી ખબર પડે છે. જો અત્યાર સુધી ફક્ત સેશનનું `--agent` તેનું નામ આપતું હોય, અને પસંદ કરેલું એકાઉન્ટ agents શેર ન કરતું હોય, તો `ccam next --continue` કંઈ શરૂ કરતું નથી અને નોંધ રાખી મૂકે છે, અને બીજી desktop સૂચના એ જણાવે છે: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

જ્યારે બીજા કોઈ એકાઉન્ટમાં જગ્યા ન હોય, ત્યારે કંઈ ખસતું નથી: Claude Code મર્યાદાને હંમેશની જેમ સંભાળે છે, સેશનમાં રાહ જોઈને કે તેના વિકલ્પો આપીને, અને સૂચના કહે છે કે કયું એકાઉન્ટ પહેલું મુક્ત થાય છે, ઉદાહરણ તરીકે `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` સીટો માટે તે `None of your other accounts in Acme has room now` કહે છે. શરૂ થયેલો ખસેડ રોકવા, Claude Code બંધ થતું હોય ત્યારે Ctrl+C દબાવો: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

ખસેડની કિંમત હોય છે. Claude Code ફરી શરૂ થાય છે, જેમાં થોડી સેકન્ડ લાગે છે, અને Claude એ શરૂ કરેલા બૅકગ્રાઉન્ડ shells, dev servers અને agents જૂના સેશન સાથે બંધ થાય છે. નવું એકાઉન્ટ તેના પહેલા જવાબ સાથે આખી વાતચીત ફરી વાંચે છે, [`--continue`](#continue) ની જેમ, અને ccam જણાવે છે કે તે કેટલી મોટી છે; ખસેડ ક્યારેય [`--brief`](#brief) વાપરતો નથી.

`auto_switch` Claude Code મર્યાદા જણાવે તે પહેલાં ક્યારેય સ્વિચ કરતું નથી, બૅકગ્રાઉન્ડમાં ક્યારેય તમારા વપરાશ પર નજર રાખતું નથી, અને ક્યારેય લૉગિન કૉપિ, રિન્યૂ કે બીજે આપતું નથી. Windows પર તેને અવગણવામાં આવે છે, અને `ccam doctor` તે જણાવે છે; ત્યાં, `relay = true` સાથે, `/ccam:switch` અને `/exit` હજી પણ વાતચીતને બીજા એકાઉન્ટ પર ખસેડે છે.

**Anthropic ની શરતો.** Team અને Enterprise સીટ પર, `auto_switch` તમારું કામ એ જ સંસ્થાની બીજી સીટ પર ખસેડે છે. Anthropic નું સહાય કેન્દ્ર કહે છે કે "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team પ્લાનની દરેક સીટ નામથી નક્કી કરેલી એક વ્યક્તિને ફાળવવામાં આવે છે, અને ફક્ત એ વ્યક્તિ જ તેનો ઉપયોગ કરી શકે છે) અને "Users can't route their work through another user's account to use that user's remaining usage" (વપરાશકર્તાઓ બીજા વપરાશકર્તાના એકાઉન્ટ મારફતે પોતાનું કામ ચલાવીને એ વપરાશકર્તાનો બાકી રહેલો વપરાશ વાપરી શકતા નથી) ([Team પ્લાન](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise પ્લાન](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) માટે પણ એવું જ).

અંગત પ્લાન (Pro, Max) માટે, Anthropic ની શરતો એ કહેતી નથી કે તમારાં પોતાનાં એકાઉન્ટ મર્યાદાએ આપમેળે એકબીજાનું કામ સંભાળી શકે કે નહીં. Anthropic સ્પષ્ટ છૂટ ન આપે ત્યાં સુધી [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) સેવાઓને "through automated or non-human means, whether through a bot, script, or otherwise" (સ્વચાલિત કે માનવ સિવાયના માધ્યમથી, bot, script કે બીજી કોઈ રીતે) વાપરવાની મનાઈ કરે છે, અને ખસેડ ઉપરનો નક્કી કરેલો સંદેશ તમારા વતી મોકલે છે. નવેમ્બર 12, 2026 થી અમલમાં આવતી [Usage Policy](https://www.anthropic.com/legal/aup) જે બાબતોની મનાઈ કરે છે તેમાં "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (અમારાં ઉત્પાદનોમાં સ્થાપિત ક્ષમતાઓ, પ્રતિબંધો કે સુરક્ષા-મર્યાદાઓને ઇરાદાપૂર્વક ટાળવી કે તેમાં દખલ કરવી) ગણાવે છે.

ccam એવો દાવો કરતું નથી કે `auto_switch` આ શરતોનું પાલન કરે છે. તે ચાલુ કરવું એ તમારો નિર્ણય છે, અને જોખમ પણ તમારું છે: Anthropic ને લાગે કે શરતોનો ભંગ થયો છે ત્યારે તે એકાઉન્ટ, કે Team કે Enterprise સંસ્થાની ઍક્સેસ, સ્થગિત કરી શકે છે.

`auto_switch` ચાલુ હોય ત્યારે, `ccam doctor` ccam plugin વગરના એકાઉન્ટ વિશે અને 2.1.214 કરતાં જૂનાં Claude Code વર્ઝન વિશે ચેતવે છે, અને સીટો માટે Anthropic નો સીટ અંગેનો નિયમ ફરી જણાવે છે ([તપાસો](/features/doctor)).

## મર્યાદાએ ફરી ચલાવવું

`ccam next` અટકી ગયેલા interactive સેશન માટે છે. non-interactive run, એટલે કે script કે શેડ્યૂલ કરેલા job માંનું `claude -p`, તમારી રાહ જોઈ શકતું નથી, તેથી એકાઉન્ટના નામ પહેલાં `--fallback` મૂકવાથી ccam તેના માટે સ્વિચ કરે છે:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` ફક્ત `-p` કે `--print` સાથે જ કામ કરે છે; નહીં તો ccam અટકી જાય છે અને `ccam next --auto --continue` તરફ ધ્યાન દોરે છે.
2. pipe કરેલું standard input એક વાર, 10 MiB સુધી, વંચાય છે અને દરેક પ્રયાસને અપાય છે. Claude Code ની જેમ, ccam તે શરૂ થાય તેની વધુમાં વધુ 3 સેકન્ડ રાહ જુએ છે (`--input-format stream-json` સિવાય); જો pipe ખુલ્લી રહે પણ કંઈ ન મોકલે, તો run ને કોઈ input મળતું નથી અને ccam standard error પર તેની નોંધ કરે છે.
3. ccam `exec` ને બદલે `claude` ને child process તરીકે શરૂ કરે છે, તેનું standard output રોકી રાખે છે, અને standard error આવે તેમ આગળ આપે છે.
4. `claude` 0 સાથે બહાર નીકળે, કે બીજા કોઈ કારણે નિષ્ફળ જાય, ત્યારે રોકેલું output અને exit code બદલ્યા વગર આગળ અપાય છે. signal થી પૂરા થયેલા run (128 થી ઉપરનો exit code) નું પણ એવું જ છે: તેને બંધ થવા કહેવાયું હતું, અને તે ફરી ચલાવાતો નથી.
5. નિષ્ફળતાને વપરાશ મર્યાદા ત્યારે ગણાય છે જ્યારે Claude Code નો સંદેશ એમ કહે, એટલે કે [hook](#inside-claude-code) જાણે છે તે વપરાશ-મર્યાદાના સંદેશમાંથી કોઈ પણ, `--output-format json` કે `stream-json` output ની અંદર પણ, અથવા જ્યારે એકાઉન્ટનો વપરાશ તેની 5 કલાકની, સાપ્તાહિક કે `--model` ની વિન્ડો 100% પર બતાવે (મર્યાદા રેકોર્ડ સહિત).
6. પછી એ પ્રયાસનું output છોડી દેવાય છે અને એ જ command બીજા લૉગિન થયેલા એકાઉન્ટ તરીકે ચાલે છે: સૌથી ઓછા 5 કલાકના વપરાશવાળું, પછી જેમનો વપરાશ જાણી ન શકાય તે એકાઉન્ટ, નંબરના ક્રમમાં. 100% પરનાં એકાઉન્ટ છોડી દેવાય છે, API key લૉગિનવાળાં એકાઉન્ટ ક્યારેય અજમાવાતાં નથી, અને દરેક એકાઉન્ટ વધુમાં વધુ એક વાર અજમાવાય છે.
7. આર્ગ્યુમેન્ટ્સમાં `-c`, `--continue`, `-r` કે `--resume` હોય ત્યારે ફક્ત એ જ history ગ્રુપનાં એકાઉન્ટ અજમાવાય છે, જેથી વાતચીત મળી શકે.
8. `--auto` ની જેમ, ફક્ત એ જ સંસ્થાનાં એકાઉન્ટ અજમાવાય છે (`cross_org = true` આ હટાવે છે).

આગલું એકાઉન્ટ `--auto` ના જ નિયમોથી પસંદ થાય છે (`auto_max_usage`, `auto_strategy`, `--model` ની વિન્ડો). કોઈ એકાઉન્ટ ન બચે ત્યારે ccam છેલ્લા પ્રયાસનું output છાપે છે, કહે છે કે કોઈ લૉગિન થયેલા એકાઉન્ટમાં જગ્યા બાકી નથી, અને code 75 સાથે બહાર નીકળે છે. એકાઉન્ટના નામ પહેલાં `--wait` સાથે (`ccam --fallback --wait=2h work -p …`), તે તેના બદલે પહેલું એકાઉન્ટ રીસેટ થાય તેની રાહ જુએ છે, પછી બધાં એકાઉન્ટ ફરી અજમાવી શકે છે; તે હંમેશાં પહેલાં ઓછામાં ઓછું એક વાર ઊંઘે છે, જેથી હમણાં જ મર્યાદાએ પહોંચેલું એકાઉન્ટ તરત ફરી ન ચાલે. સમયગાળો run ના રાહ જોવાના કુલ સમયને મર્યાદિત કરે છે, દરેક રાહને અલગથી નહીં. space સાથે લખેલું હોય (`--wait 30m`), તો `ccam next` જેવી જ સૂચના સાથે તેનો ઇનકાર થાય છે. Claude Code ના આર્ગ્યુમેન્ટ્સ પહેલાં શરૂઆતનું `--` કોઈ પણ શરૂઆતની જેમ કામ કરે છે (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` પહેલા પ્રયાસનો command છાપે છે.

તમે [`auto_switch`](#auto-switch) ચાલુ ન કરો ત્યાં સુધી interactive સેશન જાતે સ્વિચ થતાં નથી. `ccam statusline install` સાથે, 5 કલાકનો કે સાપ્તાહિક વપરાશ 90% પર પહોંચે એટલે સ્ટેટસ લાઇન `→ ccam next --auto --continue` ઉમેરે છે ([વપરાશ મર્યાદા](/features/usage)).

## સલામતી અને ખાસ કિસ્સા

- લૉગિન વગરનાં એકાઉન્ટ ક્યારેય પસંદ થતાં નથી. લૉગિન થયેલું બીજું કોઈ એકાઉન્ટ ન હોય, ત્યારે સાદું `ccam next` `no logged-in account other than <name>` સાથે અટકે છે.
- `--continue` સાથે અને ફોલ્ડરમાં કોઈ વાતચીત ન હોય ત્યારે: `no account has a conversation in <folder> to continue`.
- જ્યારે સંસ્થાનો નિયમ કોઈ એકાઉન્ટ બાકી ન રાખે: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- `--auto` સાથે અને લૉગિન થયેલું કોઈ એકાઉન્ટ ન હોય ત્યારે: `no logged-in account to pick`. જ્યારે દરેક એકાઉન્ટ ભરાઈ ગયું હોય કે તેનો વપરાશ અજાણ્યો હોય: `` no account with known usage and headroom (see `ccam usage`) ``. બંને exit code 75 સાથે પૂરા થાય છે.
- `--max-usage`, `--strategy`, `--wait` અને `--cross-org` ફક્ત `--auto` સાથે કામ કરે છે, અને 1-100 ની બહારનું મૂલ્ય કે અજાણી strategy નકારવામાં આવે છે.
- fork વાતચીતની નકલ કરે છે, એકાઉન્ટની બાકીની સ્થિતિની નહીં: ઇતિહાસ શેર ન થતો હોય તો મૂળ એકાઉન્ટના `file-history/` ના checkpoints સાથે આવતા નથી.
- `--auto` `--dry-run` સાથે પણ નેટવર્ક request કરે છે. તે ફક્ત વાંચે છે, અને ક્યારેય ટોકન ફાઇલો લખતું નથી.
- પસંદ કરેલું એકાઉન્ટ `ccam <name>` ના જ પાથથી શરૂ થાય છે, તેથી એ જ તપાસો લાગુ પડે છે, જેમાં એ જ ઇતિહાસ શેર કરતું બીજું Claude Code સેશન આ ફોલ્ડરમાં પહેલેથી ચાલતું હોય ત્યારની ચેતવણી પણ સામેલ છે.
- `last-used` ત્યારે જ લખાય છે જ્યારે Claude Code ખરેખર શરૂ થાય. `--dry-run` તેને અડતું નથી, તેથી dry run rotation ને આગળ ખસેડતું નથી.
- `--fallback` પ્રયાસ પૂરો થાય ત્યાં સુધી standard output રોકી રાખે છે, તેથી `--output-format stream-json` જેવું streaming output અંતે એકસાથે આવે છે. જે client `stream-json` input stream કરે છે અને વધુ મોકલતાં પહેલાં જવાબોની રાહ જુએ છે, તે `--fallback` વાપરી શકતું નથી. મર્યાદાએ પહોંચતાં પહેલાં પહેલા પ્રયાસે બદલેલી ફાઇલો બદલાયેલી જ રહે છે; આગલો પ્રયાસ prompt ફરીથી શરૂ કરે છે.

## સંબંધિત

- [વપરાશ મર્યાદા](/features/usage)
- [વપરાશ મર્યાદા અને એકાઉન્ટ સ્વિચ કરવું](/docs/usage-limits)
- [MCP સર્વર, સેટિંગ્સ અને ઇતિહાસ શેર કરો](/docs/sharing)
- [Claude Code વપરાશ મર્યાદાએ પહોંચે ત્યારે શું કરવું](/guides/usage-limit-workflow)
