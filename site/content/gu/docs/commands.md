---
title: કમાન્ડ સંદર્ભ
metaTitle: ccam કમાન્ડ સંદર્ભ
description: ccam ના બધા કમાન્ડ અને ઓપ્શન, ccam add અને ccam <name> થી લઈને ccam usage, ccam next અને ccam doctor સુધી.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: 28e1ad1dcbad
---

કોઈ કમાન્ડ ખરેખર ચલાવ્યા વગર તે શું કરશે તે જોવા માટે તેની આગળ `--dry-run` મૂકો, ઉદાહરણ તરીકે `ccam --dry-run work`. તે Claude Code શરૂ કરતી વખતે (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) અને `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` તથા `doctor --fix` સાથે કામ કરે છે. ફક્ત વાંચતા કમાન્ડ, જેમ કે `list`, `usage`, `share` અને `ext`, તેને સ્વીકારે છે અને આમેય કંઈ બદલતા નથી. બાકીના કમાન્ડ, `group rename` સહિત, ચાલવાને બદલે તેનો ઇનકાર કરે છે.

<!-- commands -->

## Claude Code શરૂ કરવું

`ccam <name|number> [arguments]` તે એકાઉન્ટ તરીકે Claude Code શરૂ કરે છે. ccam નામ પછીના આર્ગ્યુમેન્ટ્સ વાંચતું નથી; તે તેમને એ જ ક્રમમાં `claude` ને આપે છે, અને standard input, exit code અને સિગ્નલ બરાબર `claude` જેવું જ વર્તન કરે છે. નામ પછી `--` મૂકવું વૈકલ્પિક છે.

આર્ગ્યુમેન્ટ્સ સાથે પણ એકાઉન્ટના નામ વગરનું `ccam`, જેમ કે `ccam -p "summarise"`, હાલના ફોલ્ડર માટે સેટ કરેલું એકાઉન્ટ, નહીં તો ડિફૉલ્ટ એકાઉન્ટ વાપરે છે. તેમને `ccam folder <name> [folder]` અને `ccam use <name>` વડે સેટ કરો.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

એકાઉન્ટના નામ પહેલાં `--fallback` એ `-p`/`--print` runs માટે છે. એકાઉન્ટ તેની વપરાશ મર્યાદાએ પહોંચે ત્યારે ccam એ જ command સૌથી વધુ જગ્યાવાળા લૉગિન થયેલા એકાઉન્ટ તરીકે ફરી ચલાવે છે અને ફક્ત પૂરા થયેલા run નું output છાપે છે. જુઓ [ફોલ્ડર પ્રમાણે એકાઉન્ટ](/features/folder) અને [મર્યાદાએ ફરી ચલાવવું](/features/next#run-again-at-the-limit).

**Team અને Enterprise સીટ.** આપમેળે થતી એકાઉન્ટ પસંદગી (`ccam next`, `--fallback`, `ccam pick`, ડેશબોર્ડની `n` કી અને, તમે તેને ચાલુ કરો પછી, [`auto_switch`](/features/next#auto-switch)) એક જ Claude Team કે Enterprise સંસ્થાની સીટો વચ્ચે પણ સ્વિચ કરે છે, જેનાથી Anthropic ની શરતોનું ઉલ્લંઘન થઈ શકે છે. તેના પર આધાર રાખતાં પહેલાં [વારંવાર પૂછાતા પ્રશ્નો](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) વાંચો અને તમારી સંસ્થા અને Anthropic પાસેથી ખાતરી કરી લો.

`ccwork` જેવો શૉર્ટકટ તેના એકાઉન્ટને આ જ રીતે શરૂ કરે છે: `ccwork --resume` એ જ કરે છે જે `ccam work --resume` કરે છે. નીચે [શૉર્ટકટ](/docs/commands#shortcuts) જુઓ.

## એકાઉન્ટ ઉમેરવાં અને દૂર કરવાં

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

ફોલ્ડર વગરનું `ccam import` `~/.claude` અને `~/.claude-*` શોધે છે અને દરેક વિશે પૂછે છે; નિષ્ફળ જતું ફોલ્ડર `Skipped <path>: <reason>` તરીકે જણાવાય છે અને બાકીનાં રજિસ્ટર થાય છે. જો કોઈ ફોલ્ડરનું નામ પહેલેથી વપરાયેલું હોય, તો તે ફોલ્ડરને `ccam import <folder> --name <name>` વડે અલગથી ઇમ્પોર્ટ કરો. ઇમ્પોર્ટ કરેલાં ફોલ્ડર જ્યાં છે ત્યાં જ રજિસ્ટર થાય છે; ccam ક્યારેય લૉગિન ખસેડતું કે કૉપિ કરતું નથી.

`ccam login --expired` એવાં એકાઉન્ટમાં પણ લૉગિન કરે છે જેમનું લૉગિન Claude Code એ સાફ કરી નાખ્યું હોય. `ccam remove` મેનેજ કરેલા એકાઉન્ટને (જેનું ફોલ્ડર ccam એ બનાવ્યું હોય, `ccam import` વડે રજિસ્ટર કરેલું નહીં) તેનું ફોલ્ડર trash માં જાય તે પહેલાં લૉગઆઉટ કરે છે, સિવાય કે બીજું એકાઉન્ટ એ જ લૉગિન ધરાવતું હોય. ઇમ્પોર્ટ કરેલું એકાઉન્ટ કાઢી નાખવાથી તે તેનાં ગ્રુપમાંથી બહાર નીકળે છે અને દરેક શેર કરેલી આઇટમની કૉપિ તેના ફોલ્ડરમાં રહે છે. macOS પર `ccam rename` ફોલ્ડરનું નામ એ જ રાખે છે, કારણ કે Keychain લૉગિનને ફોલ્ડરના પાથ હેઠળ રાખે છે.

`ccam add --api` એવું એકાઉન્ટ ઉમેરે છે જે Claude લૉગિનને બદલે API key વાપરે છે; ccam દરેક લૉન્ચ વખતે key વાંચે છે અને તેને ક્યારેય સેવ કરતું નથી. [API key એકાઉન્ટ](/features/api) જુઓ.

Claude Code કોઈ એકાઉન્ટ તરીકે ચાલી રહ્યું હોય ત્યારે, ભલે તે ગમે તે રીતે શરૂ થયું હોય (બૅકગ્રાઉન્ડ સેશન સહિત), `ccam rename`, `ccam remove` અને `ccam group` કમાન્ડ તે એકાઉન્ટમાં ફેરફાર કરવાનો ઇનકાર કરે છે. પહેલાં સેશન બંધ કરો. ગ્રુપનો કોઈ ફેરફાર અધૂરો હોય ત્યારે પણ તેઓ ઇનકાર કરે છે: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` દરેક એકાઉન્ટને તેના સ્લૉટ, શૉર્ટકટ, ઇમેઇલ અને ફોલ્ડર સાથે છાપે છે. તેનો `groups` array એકાઉન્ટનાં ગ્રુપની યાદી આપે છે, અને તેનો `share` ઑબ્જેક્ટ દરેક આઇટમ (`mcp`, `settings`, `memory`, `extensions`, `history`) માટે તે જે ગ્રુપમાંથી આવે છે તે, અથવા `local` આપે છે.

## શૉર્ટકટ

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

શૉર્ટકટ એ તમે જાતે નામ આપેલો કમાન્ડ છે, જે Claude Code ને એક એકાઉન્ટ તરીકે શરૂ કરે છે. `ccam work` ની જેમ તે દરેક આર્ગ્યુમેન્ટ કોઈ ફેરફાર વગર Claude Code ને આપે છે, તેથી `ccwork --resume` અને `ccwork -p "hi"` અપેક્ષા મુજબ કામ કરે છે. તે `~/.local/bin` માં ccam ની લિંક છે.

શૉર્ટકટ વૈકલ્પિક છે. `ccam add` વિઝાર્ડ શૉર્ટકટ બનાવવો કે નહીં તે પૂછે છે, `ccam add <name> --shortcut <command>` પૂછ્યા વગર એક બનાવે છે, અને `ccam list` દરેક એકાઉન્ટનો શૉર્ટકટ બતાવે છે, અથવા શૉર્ટકટ ન હોય ત્યારે `-`.

નામમાં a-z, 0-9, `.`, `-` અને `_` અક્ષરો વાપરી શકાય છે. નામ કોઈ ccam કમાન્ડ, બીજા એકાઉન્ટનો શૉર્ટકટ, કે તમારા `PATH` પરના બીજા પ્રોગ્રામનું નામ ન હોઈ શકે, અને ccam એ ન બનાવી હોય તેવી `~/.local/bin` માંની ફાઇલ પણ નહીં, જેમ કે ccam ની તમારી પોતાની કૉપિ કે લિંક. `config.toml` માં હાથથી બદલેલો શૉર્ટકટ જો આ નિયમો તોડે, `claude` કે `ccam` હોય, અથવા બે એકાઉન્ટને આપેલો હોય, તો `config.toml` લોડ થઈ શકતી નથી, અને ભૂલ તે એકાઉન્ટનું નામ જણાવે છે.

`ccam rename` શૉર્ટકટનું નામ એ જ રાખે છે, અને `ccam remove` એકાઉન્ટની સાથે તેનો શૉર્ટકટ પણ કાઢી નાખે છે. `ccam shims install` બધાં એકાઉન્ટની લિંક ફરીથી બનાવે છે, ઉદાહરણ તરીકે ccam ને બીજે ખસેડ્યા પછી; કોઈ એક શૉર્ટકટના રસ્તામાં આવતી ફાઇલ બાકીના શૉર્ટકટને અટકાવતી નથી.

## શેરિંગ

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

ગ્રુપ તેનાં એકાઉન્ટ વચ્ચે આઇટમ શેર કરે છે; એક એકાઉન્ટ ઘણાં ગ્રુપમાં હોઈ શકે છે, પણ દરેક આઇટમ ફક્ત એક જ ગ્રુપમાંથી આવે છે. [ગ્રુપ વડે શેરિંગ](/features/share) જુઓ.

## હાલનું સેટઅપ કૉપિ કરવું

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` કોઈ Claude Code ફોલ્ડરમાંથી `settings.json`, `CLAUDE.md`, એક્સટેન્શન્સ અને user MCP સર્વર એક ગ્રુપમાં કૉપિ કરે છે, જેથી તે આઇટમ શેર કરતા દરેક એકાઉન્ટને તે મળે. ફોલ્ડર ફક્ત વંચાય છે.

- `--settings`, `--memory`, `--extensions`, `--mcp` અને `--history` ફક્ત તમે નામ આપેલી આઇટમ કૉપિ કરે છે. તેમના વગર, ટર્મિનલ પર એક ચેકલિસ્ટ તમને પસંદ કરવા દે છે, જેમાં વાતચીતના ઇતિહાસ સિવાય બધું ચેક કરેલું હોય છે; બીજે બધે વાતચીતના ઇતિહાસ સિવાય બધું કૉપિ થાય છે.
- `~/.claude` ખાલી હોય કે હોય જ નહીં, પણ `~/.claude.json` માં સર્વર હોય, ત્યારે પણ `--mcp` કામ કરે છે.
- `-g` કે `--group` ગ્રુપ પસંદ કરે છે.
- `-y` કે `--yes` પૂછ્યા વગર કૉપિ કરે છે અને કોઈ ફાઇલ જુદી હોય ત્યારે ગ્રુપની કૉપિ રાખે છે.

MCP સાઇન-ઇન કૉપિ થતાં નથી, કારણ કે કૉપિ કરેલું સાઇન-ઇન તેના બીજા ધારકોને લૉગઆઉટ કરી દે છે; યોજના એ એકાઉન્ટનાં નામ જણાવે છે જેમણે `/mcp` વડે સાઇન ઇન કરવાનું છે. [હાલનું સેટઅપ કૉપિ કરવું](/docs/sharing#copy-an-existing-setup) જુઓ.

## MCP સર્વર

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` ફક્ત તે એકાઉન્ટ પાસે જ હોય તેવા સર્વર તેના ગ્રુપમાં ખસેડે છે, તેની પાસે ફક્ત એક ફોલ્ડર માટે હોય તેવા સર્વર સહિત (સાદું `claude mcp add`), જે પછી દરેક ફોલ્ડરમાં લોડ થાય છે. જે સર્વરનું નામ ગ્રુપ પહેલેથી બીજી વ્યાખ્યા માટે વાપરે છે તેને છોડી દેવામાં આવે છે; ગ્રુપની કૉપિ બદલવા માટે `--force` ઉમેરો, ccam પહેલાં તેનો બૅકઅપ લે છે.

## વપરાશ અને સ્વિચિંગ

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

- `next --auto`, `--fallback` અને `ccam pick` એક જ રીતે પસંદ કરે છે: પહેલાં `auto_max_usage` (`config.toml`, ડિફૉલ્ટ 100) થી નીચેનાં, `auto_strategy` (`headroom`, `expiring` કે `balance`) પ્રમાણે, `--model` ની સાપ્તાહિક વિન્ડો ગણીને. તેઓ 100% પરનું એકાઉન્ટ (જે વિન્ડો ભરાઈ ગઈ હોવાનું Claude Code એ જણાવ્યું હોય તે તેના રીસેટ સુધી 100% ગણાય છે) કે API key લૉગિન ક્યારેય પસંદ કરતા નથી, અને જે એકાઉન્ટથી તેઓ શરૂ કરે તેની સંસ્થાની અંદર રહે છે (વધુ આગળ જોવા માટે `--cross-org` કે `cross_org = true`).
- `--wait` તેનું મૂલ્ય `=` પછી લે છે (`--wait=2h`); space સાથે લખો તો સૂચના સાથે તેનો ઇનકાર થાય છે. કોઈ પસંદ ન થઈ શકે ત્યારે તેઓ code 75 સાથે બહાર નીકળે છે, અથવા, `--wait` સાથે, પહેલું એકાઉન્ટ રીસેટ થાય તેની રાહ જુએ છે.
- `ccam next --continue` ઇતિહાસ શેર ન કરતાં એકાઉન્ટ વચ્ચે પણ વાતચીત આગળ લઈ જાય છે: ત્યારે Claude Code તેની નકલ નવા એકાઉન્ટમાં કરે છે. `--brief` સાથે, નવું એકાઉન્ટ તેના બદલે નવી વાતચીત શરૂ કરે છે, ccam એ કોઈ મૉડલ વગર છેલ્લી વાતચીતનો લખેલો સારાંશ લઈને ([`--brief`](/features/next#brief)).
- `ccam <name>` થી અલગ, `ccam next` પોતાના ઓપ્શન જાતે વાંચે છે, તેથી Claude Code ના ઓપ્શન `--` પછી આવે છે; તેના વગર ccam `unknown flag` સાથે અટકી જાય છે.
- `ccam agent` આ ફોલ્ડરમાં `claude -p` વડે એક task, જે એકાઉન્ટમાંથી તેને બોલાવાયું હોય તેના સિવાયના બીજા એકાઉન્ટ તરીકે, બોલાવનાર સેશનના permission mode અને effort માં ચલાવે છે, અને Claude નો છેલ્લો જવાબ છાપે છે. એ એકાઉન્ટની મર્યાદાએ વાતચીત બીજા એકાઉન્ટ તરીકે આગળ ચાલે છે ([બીજાં એકાઉન્ટ પર subagents](/features/agent)).
- Claude Code ની અંદર, plugin નું `/ccam:switch [account]` વાતચીત આગળ સોંપે છે, અને તેના hooks `ccam hook` ચલાવે છે.
- `config.toml` માં `auto_switch = true` સાથે (Linux અને macOS; મૂળભૂત રીતે બંધ), મર્યાદા hook ccam એ શરૂ કરેલું interactive સેશન બંધ પણ કરે છે અને તેને જગ્યા ધરાવતા બીજા એકાઉન્ટ પર ચાલુ રાખે છે, પણ ફક્ત એ જ પ્રકારના એકાઉન્ટ પર: એ જ સંસ્થાની બીજી સીટ, કે બીજો અંગત પ્લાન ([`auto_switch`](/features/next#auto-switch)).

જુઓ [આગલા એકાઉન્ટ પર સ્વિચ કરવું](/features/next), [scripts માટે એકાઉન્ટ પસંદ કરવું](/features/pick), [બીજાં એકાઉન્ટ પર subagents](/features/agent) અને [દરેક એકાઉન્ટનાં બૅકગ્રાઉન્ડ સેશન](/features/agents).

## જાળવણી

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

ccam ને બદલતાં પહેલાં `ccam update` રિલીઝનો SHA-256 checksum તપાસે છે, અને `--version v0.2.0` કોઈ ચોક્કસ રિલીઝ ઇન્સ્ટૉલ કરે છે. `.deb` પેકેજથી ઇન્સ્ટૉલ કરેલી ccam ની કૉપિ તેના બદલે apt વડે અપડેટ થાય છે. [અપડેટ](/docs/install#update) જુઓ.

`ccam uninstall` શું દૂર કરશે તેની યાદી બતાવે છે અને પહેલાં પૂછે છે: `ccam remove` ની જેમ દરેક એકાઉન્ટ, શૉર્ટકટ કમાન્ડ, ccam નાં ફોલ્ડર અને `ccam` કમાન્ડ. ઇમ્પોર્ટ કરેલું ફોલ્ડર જ્યાં છે ત્યાં જ રહે છે. [અનઇન્સ્ટૉલ](/docs/install#uninstall) જુઓ.

## ફાઇલો

| પાથ | સામગ્રી |
| --- | --- |
| `~/.config/ccam/config.toml` | એકાઉન્ટ, ગ્રુપ અને સેટિંગ્સ (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). સેટિંગ્સ સૌથી ઉપર, પહેલા `[[group]]`, `[[profile]]` કે `[[folder]]` ની ઉપર આવે છે: કોઈ table ની નીચેનું સેટિંગ તે table નું ગણાય છે, તેથી ccam દરેક લૉન્ચ વખતે અને `ccam doctor` માં તેના વિશે ચેતવણી આપે છે, અને config સેવ કરતા કમાન્ડ તેને ખસેડવામાં ન આવે ત્યાં સુધી ઇનકાર કરે છે. ccam ફાઇલ સેવ કરે ત્યારે કૉમેન્ટ સચવાતી નથી |
| `~/.local/share/ccam/profiles/<name>/` | દરેક એકાઉન્ટનું `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | શેર કરેલાં સેટિંગ્સ, MCP સર્વર, એક્સટેન્શન્સ અને ઇતિહાસ |
| `~/.local/state/ccam/` | વપરાશ કૅશ, મર્યાદા રેકોર્ડ, lock, નોંધો અને બીજી સ્થિતિ (કોષ્ટક નીચેની યાદી જુઓ) |
| `~/.local/bin/<shortcut>` | `ccwork` જેવા શૉર્ટકટ કમાન્ડ, ccam ની લિંક |

સ્ટેટ ફોલ્ડરમાં આ હોય છે:

- વપરાશ કૅશ અને મર્યાદા રેકોર્ડ (`usage/<name>.json`, `usage/<name>.limit.json`);
- છેલ્લે વાપર્યાના સમય;
- સેશન માર્કર અને lock (`locks/`);
- ccam એ બનાવેલા શૉર્ટકટ (`shortcuts.json`);
- પૂરો ન થયેલો ગ્રુપનો ફેરફાર (`group-change.json`);
- `/ccam:switch` તથા મર્યાદા hook છોડી જાય તે નોંધો (`handoff/<account>/<session>.json`);
- `ccam pick` કે `auto_switch` એ દરેક એકાઉન્ટ છેલ્લે ક્યારે પસંદ કર્યું (`picks/<name>`; એક પસંદગી 2 મિનિટ સુધી ગણાય છે);
- `--brief` ના સારાંશ, 7 દિવસ રખાય છે (`brief/<session>.md`).

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` કે `XDG_STATE_HOME` સેટ કરવાથી આ ફોલ્ડર ક્યાં છે તે બદલાય છે. Windows પર `~` એટલે `%USERPROFILE%`, અને શૉર્ટકટ `<shortcut>.exe` હોય છે.
