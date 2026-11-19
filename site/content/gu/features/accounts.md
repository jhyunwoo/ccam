---
title: એકાઉન્ટનું સંચાલન
description: ccam add, ccam login, ccam list, ccam rename, ccam remove અને ccam import વડે Claude Code એકાઉન્ટ ઉમેરો, લૉગિન કરો, જુઓ, નામ બદલો, દૂર કરો કે ઇમ્પોર્ટ કરો.
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

દરેક ccam એકાઉન્ટ પોતાના લૉગિન સાથેનું અલગ Claude Code config ફોલ્ડર છે. એકાઉન્ટ ઉમેરવા, તેમાં સાઇન ઇન કરવા, તેમની યાદી જોવા, નામ બદલવા અને તેમને દૂર કરવા માટે આ કમાન્ડ વાપરો. તે તમારી પાસે પહેલેથી હોય તેવાં ફોલ્ડર પણ રજિસ્ટર કરે છે.

## ઉપયોગ

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

જ્યાં પણ કોઈ કમાન્ડ `<name|number>` લે છે, ત્યાં તમે એકાઉન્ટનું નામ કે તેનો સ્લૉટ નંબર (`ccam list` માં `#` કૉલમ) આપી શકો છો.

`ccam use` ફક્ત એટલું બદલે છે કે કયું એકાઉન્ટ ડિફૉલ્ટ છે (`config.toml` માં `default`).

ટર્મિનલ પર કોઈ ઓપ્શન વગરનું `ccam add` એક વિઝાર્ડ ખોલે છે. તે આ પગલાંમાંથી પસાર થાય છે:

1. એકાઉન્ટનું નામ, જો તમે આપ્યું ન હોય તો;
2. એકાઉન્ટ કેવી રીતે સાઇન ઇન કરે છે: હમણાં Claude લૉગિન, API key કે પછીથી Claude લૉગિન;
3. કયાં ગ્રુપમાં જોડાવું, `default` ચેક કરેલું હોય તે રીતે (હજી કોઈ ગ્રુપ ન હોય, તો `default` ગ્રુપ શું શેર કરશે તે બતાવે છે);
4. `Create a shortcut command for work? [y/N]`;
5. ખાતરી કરવા માટે સારાંશ.

`Esc` એક પગલું પાછળ જાય છે. `--group`, `--no-group`, `--shortcut` કે `--no-login` માંથી કોઈ પણ ઓપ્શન, અથવા ટર્મિનલ વગર ચલાવવું (ઉદાહરણ તરીકે, pipe માંથી), વિઝાર્ડ છોડી દે છે. ત્યારે નામ આપવું જરૂરી છે, અને `--shortcut` કોઈ નામ ન આપે ત્યાં સુધી કોઈ શૉર્ટકટ બનતો નથી.

`--group` વગર, નવું એકાઉન્ટ `default` ગ્રુપમાં જોડાય છે. પહેલું એકાઉન્ટ તેને બનાવે છે, જે MCP સર્વર, સેટિંગ્સ, `CLAUDE.md` અને એક્સટેન્શન્સ શેર કરે છે; જુઓ [ગ્રુપ વડે શેરિંગ](/features/share). Windows પર ફાઇલો ફક્ત Developer Mode માં જ લિંક થઈ શકે છે. તેના વગર, પહેલું એકાઉન્ટ ફક્ત `mcp` અને `extensions` સાથે `default` બનાવે છે અને તેમ જણાવે છે (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), અને સેટિંગ્સ, મેમરી કે ઇતિહાસ શેર કરતા ગ્રુપ સાથે ccam `--group` નો ઇનકાર કરે છે.

`ccam list` સ્લૉટ, નામ, શૉર્ટકટ, ઇમેઇલ, ગ્રુપ અને છેલ્લો ઉપયોગ બતાવે છે, અને ડિફૉલ્ટ એકાઉન્ટની આગળ `*` મૂકે છે. ઇમેઇલ એ છે જેનાથી ફોલ્ડર અત્યારે લૉગિન છે (તેની `.claude.json` માંથી); તે ccam એ નોંધેલા ઇમેઇલથી જુદો હોય ત્યારે, નોંધેલો ઇમેઇલ પાછળ `(was …)` તરીકે આવે છે:

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

`--json` સાથે, દરેક એકાઉન્ટ `slot`, `shortcut`, `name`, `email` (કોષ્ટકની જેમ, હાલના લૉગિનનો), `default`, `dir`, `last_used` (એકાઉન્ટ ક્યારેય શરૂ ન થયું હોય ત્યારે છોડી દેવાય છે), `groups` અને `share` ધરાવતો ઑબ્જેક્ટ છે. `share` દરેક આઇટમ (`mcp`, `settings`, `memory`, `extensions`, `history`) ને તે જે ગ્રુપમાંથી આવે છે તેની સાથે, અથવા `"local"` સાથે જોડે છે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `ccam add --no-login` | લૉગિન કર્યા વગર એકાઉન્ટ બનાવે છે; પછીથી `ccam login <name>` વડે લૉગિન કરો | off |
| `ccam add --group <group>` | `default` ને બદલે આ ગ્રુપમાં જોડાય છે. ઘણી વાર આપી શકાય; `default` અસ્તિત્વમાં ન હોય ત્યારે બનાવાય છે | `default` |
| `ccam add --no-group` | કોઈ ગ્રુપમાં જોડાતું નથી: એકાઉન્ટ કંઈ શેર કરતું નથી | off |
| `ccam add --shortcut <command>` | `~/.local/bin` માં આ શૉર્ટકટ કમાન્ડ પણ બનાવે છે | કોઈ નહીં |
| `ccam add --api --key-from <source>` | લૉગિનને બદલે API key વાપરતું એકાઉન્ટ; જુઓ [API key એકાઉન્ટ](/features/api) | off |
| `ccam list --json` | એકાઉન્ટને JSON તરીકે છાપે છે | off |
| `ccam remove -y`, `--yes` | પૂછ્યા વગર દૂર કરે છે | off |
| `ccam logout -y`, `--yes` | પૂછ્યા વગર લૉગઆઉટ કરે છે | off |
| `ccam import --name <name>` | તમે આપેલા ફોલ્ડર માટે એકાઉન્ટનું નામ; ફોલ્ડર વગર ઇનકાર થાય છે (`--name needs a folder: ccam import <folder> --name ops`) | ફોલ્ડરના નામ પરથી |
| `ccam import -y`, `--yes` | મળેલું દરેક ફોલ્ડર પૂછ્યા વગર ઇમ્પોર્ટ કરે છે | off |
| `--dry-run` | `ccam list` અને `ccam completion` તેને સમર્થન આપે છે, કારણ કે તેઓ કંઈ બદલતા નથી. `add`, `login`, `logout`, `use`, `rename`, `remove` અને `import` તેનો `` --dry-run is not supported by `ccam add`; nothing was changed `` (કમાન્ડના નામ સાથે) સંદેશ સાથે ઇનકાર કરે છે | off |

`ccam login` (`--expired` સિવાય), `ccam use` અને `ccam rename` ના પોતાના કોઈ ઓપ્શન નથી. `ccam completion` `bash`, `zsh` કે `powershell` લે છે; બીજું કંઈ પણ `only bash, zsh and powershell are supported` સાથે નિષ્ફળ જાય છે.

એકાઉન્ટનાં નામ `a-z`, `0-9`, `-` અને `_` ના 1 થી 32 અક્ષરનાં હોય છે, અને અક્ષર કે અંકથી શરૂ થાય છે. નામ ફક્ત અંકોનું ન હોઈ શકે, કારણ કે `ccam 2` નો અર્થ સ્લૉટ 2 છે. નામ કોઈ ccam કમાન્ડ પણ ન હોઈ શકે (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), તેથી `ccam <name>` ક્યારેય અસ્પષ્ટ નથી. આ યાદી `internal/config` માં `config.Reserved` છે.

## તે કેવી રીતે કામ કરે છે

એકાઉન્ટ ઉમેરવાથી આ ક્રમમાં તેનું ફોલ્ડર બને છે, તે તેનાં ગ્રુપમાં જોડાય છે, config સેવ થાય છે, સાઇન ઇન થાય છે અને શૉર્ટકટ બને છે. જ્યારે તમે `ccam add work --shortcut ccwork` ચલાવો, ત્યારે ccam આ પગલાં લે છે (`internal/cli/profiles.go` માં `addProfile`):

<!-- motion:accounts-add -->

1. કંઈ પણ બનાવતાં પહેલાં ccam નામ (`config.ValidateName`) અને શૉર્ટકટ (`shim.Check`) તપાસે છે.
2. ccam `~/.local/share/ccam/profiles/work/` ને mode `0700` સાથે બનાવે છે અને એકાઉન્ટને સૌથી નાનો ખાલી સ્લૉટ નંબર આપે છે, તેથી દૂર કરેલા એકાઉન્ટનો નંબર ફરી વપરાય છે (`profile.Add`, `Config.NextSlot`). પહેલું એકાઉન્ટ ડિફૉલ્ટ બને છે.
3. એકાઉન્ટ તેનાં ગ્રુપમાં જોડાય છે (`--group` કે `--no-group` બીજું કંઈ ન કહે તો `default`; પહેલું એકાઉન્ટ `default` બનાવે છે), અને ccam નવા, ખાલી ફોલ્ડરને તેમની સાથે લિંક કરે છે, ઉદાહરણ તરીકે `settings.json` ને `groups/default/settings.json` સાથે (`applyShare`, `share.Build` અને પછી `Apply` વડે). જો આ નિષ્ફળ જાય તો નવું ફોલ્ડર દૂર થાય છે, કારણ કે બાકી રહેલું ફોલ્ડર દરેક ફરી પ્રયાસને `folder already exists` સાથે અટકાવી દે.
4. એકાઉન્ટ `fsutil.AtomicWrite` દ્વારા `~/.config/ccam/config.toml` માં સેવ થાય છે: એક કામચલાઉ ફાઇલ, `fsync`, પછી `rename`, તેથી crash થાય તો પણ config ક્યારેય અડધું લખાયેલું રહેતું નથી.
5. ccam `CLAUDE_CONFIG_DIR` ને એકાઉન્ટ ફોલ્ડર પર સેટ કરીને `claude auth login` શરૂ કરે છે, અને લૉગિન પેજ ભરાયેલું મળે તે માટે એકાઉન્ટના છેલ્લા ઇમેઇલ સાથે `--email` આપે છે (`runLogin`). જો લૉગિન નિષ્ફળ જાય, તો એકાઉન્ટ રહે છે, અને ccam ભૂલ છાપે છે અને પછીથી `ccam login work` ચલાવવા કહે છે.
6. ccam `profiles/work/.claude.json` માંથી `oauthAccount` વાંચે છે (`profile.ReadIdentity`): ઇમેઇલ, એકાઉન્ટ અને સંસ્થાનાં IDs, સંસ્થાનું નામ અને પ્લાન (`claude_max`, `claude_team` …). એ ફાઇલ કેટલાક મેગાબાઇટની હોઈ શકે છે, તેથી ccam આખી ફાઇલ parse કરવાને બદલે એ એક key શોધે છે (`TopLevelValue`). ccam આ વિગતો `config.toml` માં `[profile.login]` માં રાખે છે.
7. ccam `~/.local/bin/ccwork` ને ccam binary સાથે લિંક કરે છે, જે શરૂ થાય ત્યારે `config.toml` માં `argv[0]` શોધે છે, અને લિંકને `~/.local/state/ccam/shortcuts.json` માં નોંધે છે (`shim.Install`). જો લિંક ન બની શકે, તો પણ એકાઉન્ટ ઉમેરાય છે, ccam કારણ જણાવે છે (`Could not create shortcuts: …`) અને સારાંશ `shortcut: ccwork (not created; see above)` બતાવે છે.

લૉગિનમાં મિનિટો લાગી શકે છે, અને એ દરમિયાન બીજી ccam process config સેવ કરી શકે છે. તેથી રાહ જોતા કમાન્ડ `config.toml` ફરી લોડ કરે છે, એક ફેરફાર (અહીં, ઇમેઇલ) લાગુ કરે છે અને સેવ કરે છે (`env.update`). આ જ કારણે `ccam remove` પોતાના પ્રશ્ન પછી ફરી લોડ કરે છે.

`ccam login` હાલના એકાઉન્ટ માટે પગલાં 5 અને 6 ચલાવે છે. Claude Code ઇમેઇલને ફક્ત લૉગિન પેજ માટેના સંકેત તરીકે લે છે, તેથી તમે પહેલાં કરતાં બીજા Claude એકાઉન્ટમાં સાઇન ઇન થઈ શકો છો; ccam નવા લૉગિનની `[profile.login]` સાથે સરખામણી કરે છે અને તેમ જણાવે છે (`ccam: this login is other@company.com, but work was dev@company.com before.`). તે નવું લૉગિન રાખે છે અને તેને કેવી રીતે પાછું ફેરવવું તે જણાવે છે: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

બીજા Claude એકાઉન્ટમાં થયેલું લૉગિન એકાઉન્ટનો [મર્યાદા રેકોર્ડ](/features/usage#the-limit-record) પણ કાઢી નાખે છે, કારણ કે એ મર્યાદાઓ બીજા એકાઉન્ટની હતી. જ્યારે બીજું ccam એકાઉન્ટ એ જ Claude એકાઉન્ટમાં લૉગિન હોય, ત્યારે ccam ચેતવે છે કે બંને તેમની વપરાશ મર્યાદા શેર કરે છે.

`ccam login --expired` એક પછી એક દરેક એવા એકાઉન્ટમાં લૉગિન કરે છે જે ક્યારેય લૉગિન ન થયું હોય અથવા જેના લૉગિનની મુદત પૂરી થઈ ગઈ હોય: તેના refresh token ની મુદત પૂરી થઈ હોય, કે સર્વરે ઇનકાર કર્યા પછી Claude Code એ લૉગિન સાફ કરી નાખ્યું હોય. પહેલું બ્રાઉઝર ખૂલે તે પહેલાં, તે તપાસે છે કે `config.toml` સેવ થઈ શકે છે.

`ccam logout` એકાઉન્ટમાં `claude auth logout` ચલાવે છે અને સેવ કરેલો ઇમેઇલ સાફ કરે છે; લૉગિનની નોંધ રહે છે, તેથી આગલા લૉગિનની તેની સાથે સરખામણી થાય છે. Claude Code લૉગિન રદ કરે છે અને એકાઉન્ટ માટે જે કંઈ સાચવે છે તે કાઢી નાખે છે, MCP સર્વરનાં સાઇન-ઇન સહિત. તેથી નીચેની સ્થિતિમાં ccam પહેલાં ચેતવે છે અને પૂછે છે (`--yes` પ્રશ્ન છોડી દે છે):

- લૉગિનમાં MCP સાઇન-ઇન હોય;
- એકાઉન્ટ ચાલી રહ્યું હોય (તેનાં સેશન પણ સાઇન આઉટ થાય છે);
- બીજા એકાઉન્ટનું ફોલ્ડર એ જ લૉગિન ધરાવતું હોય (કૉપિ કરેલું ફોલ્ડર), જેને લૉગઆઉટ પણ રદ કરે છે.

ઉદાહરણ તરીકે:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

બીજા Claude એકાઉન્ટ તરીકે લૉગિન કરવા માટે, ફક્ત `ccam login <name>` પૂરતું છે, અને તે MCP સાઇન-ઇન રાખે છે.

જ્યારે એકાઉન્ટ એકથી વધુ સંસ્થાનાં હોય ત્યારે `ccam list` એક `Organization` કૉલમ ઉમેરે છે (અંગત પ્લાન `personal` બતાવે છે), અને `ccam list --json` માં દરેક એકાઉન્ટ માટે `plan` અને `org` હોય છે.

### રિમોટ મશીન પર સાઇન ઇન કરવું

SSH પર, `claude auth login` નો બ્રાઉઝરવાળો ભાગ તૂટી શકે છે. Claude Code રિમોટ મશીન પરના callback સર્વર (`http://localhost:<port>/callback`) પર સાઇન-ઇનની રાહ જુએ છે, પણ સાઇન-ઇન પૂરું કરતું બ્રાઉઝર તમારા કમ્પ્યુટર પર ચાલે છે, જ્યાં એ સરનામું ક્યાંય જતું નથી અને બ્રાઉઝર `ERR_CONNECTION_REFUSED` જેવી ભૂલ બતાવે છે. સરનામામાં તો પણ એક વખતનો કોડ હોય છે, તેથી ccam તે પહોંચાડી શકે છે:

<!-- motion:accounts-relay -->

ccam તમારું ટર્મિનલ Claude Code ને સોંપતું નથી (`runLogin`). તે Claude Code ના standard input ને એક pipe સાથે જોડે છે અને તમે લખો છો તે દરેક લાઇન તેમાં કૉપિ કરે છે (`relayLogin`), તેથી Claude Code ના `Paste code here` પ્રશ્ન પર પેસ્ટ કરેલો કોડ પણ પહોંચે છે. જ્યારે કોઈ લાઇન callback સરનામું હોય, ત્યારે ccam તેને ફક્ત ત્યારે જ સ્વીકારે છે જ્યારે host `localhost`, `127.0.0.1` કે `[::1]` હોય, port આપેલો હોય, path `/callback` હોય અને `code` તથા `state` બંને હાજર હોય; `http://` છોડી શકાય છે (`callbackURL`). પછી ccam સર્વર પર એ સરનામાને 10 સેકન્ડના timeout સાથે અને redirect અનુસર્યા વગર request મોકલે છે, જેમ બ્રાઉઝરે કર્યું હોત (`deliverCallback`). request રાહ જોતા `claude auth login` ના callback સર્વર સુધી પહોંચે છે (ડાયાગ્રામમાં તૂટક લાઇન). Claude Code કોડ મેળવે છે અને લૉગિન સેવ કરે છે, Linux પર એકાઉન્ટ ફોલ્ડરની `.credentials.json` માં. પેસ્ટ કરેલી લાઇન ccam ને ક્યારેય બીજા host નો સંપર્ક કરાવી શકતી નથી.

લૉગિન શરૂ થાય તે પહેલાં ccam આ છાપે છે:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

સફળ થાય ત્યારે તે `→ Passed the sign-in to Claude Code on this machine.` છાપે છે.

લૉગિન ટર્મિનલ પર input ની રાહ જુએ ત્યારે, ccam ટર્મિનલનો echo અને line mode બંધ કરે છે અને તમે ટાઇપ કરો તે દરેક અક્ષર માટે એક `*` બતાવે છે, જેથી paste કરેલો code કે callback સરનામું સ્ક્રીન પર ન દેખાય. ccam line editing જાતે કરે છે (`readLine`): `Enter` લાઇન મોકલે છે, `Backspace` એક અક્ષર ભૂંસે છે, અને ખાલી લાઇન પર `Ctrl+D` input પૂરું કરે છે. લૉગિન પૂરું થાય ત્યારે ટર્મિનલ પહેલાં જેવું થઈ જાય છે, અને `Ctrl+C`, `SIGTERM` કે `SIGHUP` પર પણ, જે હજી પણ ccam ને બંધ કરે છે. standard input અને standard error પર ટર્મિનલ ન હોય ત્યારે, લાઇનો જેમ આવે તેમ કૉપિ થાય છે.

### નામ બદલવું, દૂર કરવું અને ઇમ્પોર્ટ કરવું

`ccam rename` અને `ccam remove` એકાઉન્ટના ફોલ્ડર પર અને ccam બીજે એકાઉન્ટના નામે જે રાખે છે તેના પર કામ કરે છે. ccam એ પાથની યાદી `nameState` માં રાખે છે:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- વપરાશ કૅશ `~/.local/state/ccam/usage/<name>.json` તથા `<name>.statusline.json`;
- મર્યાદા રેકોર્ડ `<name>.limit.json`;
- `~/.local/state/ccam/handoff/<name>/` માંની સેશન નોંધો.

`ccam remove work` તેમને આ ક્રમમાં સંભાળે છે:

<!-- motion:accounts-remove -->

1. `ccam remove` પહેલાં એકાઉન્ટ ચાલતું હોય તો ઇનકાર કરે છે, પછી પૂછે છે, અને તમે જવાબ આપો પછી ફરી તપાસે છે (આથી પ્રશ્ન ખુલ્લો હતો ત્યારે શરૂ થયેલું સેશન પકડાય છે).
2. કંઈ પણ બદલતાં પહેલાં તે તપાસે છે કે `config.toml` સેવ થઈ શકે છે.
3. જો એકાઉન્ટ એક્સટેન્શન્સ શેર કરતું હોય, તો ccam ગ્રુપની `installed_plugins.json` અને `known_marketplaces.json` ફરીથી લખે છે (`RewritePluginPaths`), કારણ કે Claude Code plugin ના absolute પાથ સાચવે છે જે આ એકાઉન્ટના ફોલ્ડરમાંથી પસાર થતા હોઈ શકે.
4. પછી મેનેજ કરેલું એકાઉન્ટ (જેનું ફોલ્ડર ccam એ બનાવ્યું હોય, `ccam import` વડે રજિસ્ટર કરેલું નહીં) લૉગઆઉટ થાય છે: તેના ફોલ્ડરમાં `claude auth logout` ચાલે છે, જેથી trash માં કોઈ કામ કરતું લૉગિન પડ્યું ન રહે (કે, macOS પર, ફોલ્ડરના નામવાળી Keychain આઇટમમાં, જે પછીથી એ જ નામનું એકાઉન્ટ વારસામાં મેળવે). API key એકાઉન્ટ માટે, અને જ્યારે બીજા એકાઉન્ટનું ફોલ્ડર એ જ લૉગિન ધરાવતું હોય ત્યારે આ પગલું છોડી દેવાય છે, કારણ કે લૉગઆઉટ તેને પણ રદ કરી દે. નિષ્ફળ જતું લૉગઆઉટ જણાવાય છે, અને દૂર કરવાનું ચાલુ રહે છે.
5. ccam ફોલ્ડરને `~/.local/share/ccam/trash/<name>-<time>/` માં ખસેડે છે (`profile.Trash`), એકાઉન્ટ `config.toml` માંથી નીકળી જાય છે, `nameState` પાથ કાઢી નાખવામાં આવે છે જેથી પછીથી એ જ નામનું એકાઉન્ટ તેમને વારસામાં ન મેળવે, અને શૉર્ટકટની લિંક દૂર થાય છે: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

જો છેવટે `config.toml` સેવ ન થઈ શકે, તો ફોલ્ડર trash માંથી પાછું આવે છે. દૂર કરેલું એકાઉન્ટ પાછું મેળવવા માટે, ફોલ્ડરને trash માંથી પાછું ખસેડો અને ફરી લૉગિન કરો.

`ccam rename` ફોલ્ડરનું નામ બદલે છે અને તેના `plugins/` માંના plugin પાથ જૂના ફોલ્ડરથી નવા ફોલ્ડર પર ફરીથી લખે છે. તે `nameState` પાથ નવા નામ પર ખસેડે છે, અને એકાઉન્ટ ડિફૉલ્ટ હોય તો તેને ડિફૉલ્ટ જ રાખે છે. તે સેશન નોંધો પણ અપડેટ કરે છે, જેથી `ccam next --continue` તેમને હજી પણ શોધી શકે. શૉર્ટકટ કમાન્ડ પોતાનું નામ રાખે છે. જો plugin પાથ ફરીથી લખવામાં નિષ્ફળ જાય, કે `config.toml` સેવ ન થઈ શકે, તો ફોલ્ડરનું નામ પાછું બદલાય છે.

macOS પર, અને Claude Code લૉગિન Credential Manager માં રાખે ત્યારે Windows પર, ફોલ્ડર પોતે પોતાનું જૂનું નામ રાખે છે, કારણ કે લૉગિન અને MCP સાઇન-ઇન ફોલ્ડરના પાથ હેઠળ સચવાય છે અને ખસેડેલું ફોલ્ડર લૉગઆઉટ થઈ જાય. પછી config ફોલ્ડરને નવા નામની બાજુમાં `folder` તરીકે નોંધે છે, અને ccam `its folder stays …: the macOS Keychain keeps the login under the folder's path` જણાવે છે.

`ccam import` ફોલ્ડરને જ્યાં છે ત્યાં જ, `config.toml` માં `dir` તરીકે રજિસ્ટર કરે છે. તે તેને ક્યારેય ખસેડતું કે કૉપિ કરતું નથી: કૉપિ થયેલો refresh token નકલોને એકબીજાને લૉગઆઉટ કરાવી દે.

ફોલ્ડર આર્ગ્યુમેન્ટ વગર, તે રજિસ્ટર ન થયેલાં `~/.claude` અને `~/.claude-*` ફોલ્ડર શોધે છે અને દરેક વિશે પૂછે છે. નામ ફોલ્ડર પરથી આવે છે (`~/.claude` એ `main`, `~/.claude-work` એ `work`, `~/.claude-2` એ `acct2`). જો ccam નામ કાઢી ન શકે, કે નામ પહેલેથી વપરાયેલું હોય, તો તે ફોલ્ડરને ચલાવવાના `--name` કમાન્ડ સાથે બતાવે છે. મળેલું જે ફોલ્ડર નિષ્ફળ જાય તે `Skipped <path>: <reason>` તરીકે જણાવાય છે અને બાકીનાં ઇમ્પોર્ટ થાય છે.

ccam પહેલાં બધા પ્રશ્નો પૂછે છે, પછી `config.toml` ફરી લોડ કરે છે અને ફોલ્ડર રજિસ્ટર કરે છે, તેથી એ દરમિયાન બીજા ccam કમાન્ડે સેવ કરેલા ફેરફારો સચવાઈ રહે છે. `✓` લાઇનો સેવ પછી આવે છે.

`CLAUDE_CONFIG_DIR` વગર, Claude Code `.claude.json` ને તમારા home ફોલ્ડરમાં `~/.claude` ની બાજુમાં રાખે છે, તેથી જ્યારે તમે `~/.claude` ઇમ્પોર્ટ કરો અને તેની પોતાની `.claude.json` ન હોય, ત્યારે ccam `~/.claude/.claude.json` થી `~/.claude.json` ની લિંક બનાવે છે. પછી સાદું `claude` અને `ccam main` એક જ એકાઉન્ટ જુએ છે.

જો `~/.claude` પાસે `~/.claude.json` ની બાજુમાં પોતાની `.claude.json` હોય, તો ccam બંનેને રહેવા દે છે અને તેમ જણાવે છે. `ccam main` અંદરની ફાઇલ અને સાદું `claude` બહારની ફાઇલ વાપરે, જેથી તેમની લૉગિન વિગતો અને MCP સર્વર અલગ રહે. એક જ રાખવા માટે, અંદરની ફાઇલને લિંકથી બદલો ([doctor](/features/doctor) પણ તે જણાવે છે).

જ્યારે ફોલ્ડર કોઈ રજિસ્ટર થયેલા એકાઉન્ટ જેવું જ લૉગિન ધરાવતું હોય (કૉપિ કરેલું ફોલ્ડર), ત્યારે ccam તેમ જણાવે છે: જે પહેલું તેને રિન્યૂ કરે તે બીજાને લૉગઆઉટ કરી દે છે, તેથી તેને વાપરતાં પહેલાં `ccam login <name>` ચલાવો.

તમે તેમને કોઈ ગ્રુપમાં ન ઉમેરો (`ccam group join <group> <name>`) ત્યાં સુધી ઇમ્પોર્ટ કરેલાં એકાઉન્ટ કોઈ ગ્રુપમાં હોતાં નથી, અને કોઈ શૉર્ટકટ બનતો નથી.

`ccam completion bash|zsh` એક completion સ્ક્રિપ્ટ છાપે છે, જે એકાઉન્ટનાં નામ પણ પૂરાં કરે છે.

## સલામતી અને ખાસ કિસ્સા

- એકાઉન્ટ ચાલતું હોય ત્યારે `ccam rename` અને `ccam remove` (`refuseRunning`) ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)`` સાથે ઇનકાર કરે છે, કારણ કે સેશન ફોલ્ડરમાં લખવાનું ચાલુ રાખે છે. આ તપાસ Claude Code ગમે તે રીતે શરૂ થયું હોય તેને આવરી લે છે: ccam દ્વારા, સાદું `claude`, IDE, બૅકગ્રાઉન્ડ સેશન અને તેમનો daemon.
- બીજી ccam process તે લૉગિન રિન્યૂ કરતી કે બદલતી હોય ત્યારે `ccam rename`, `ccam remove` અને `ccam logout` ઇનકાર કરે છે: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- બીજી ccam process માં ગ્રુપનો ફેરફાર ચાલુ હોય ત્યારે `ccam rename` અને `ccam remove` `another ccam process is changing the groups right now; try again when it is done` સાથે ઇનકાર કરે છે; કોઈ ફેરફાર અધવચ્ચે અટકી ગયો હોય પછી, ``the last group change stopped halfway; finish it with `ccam doctor --fix` first`` સાથે.
- ઇમ્પોર્ટ કરેલા એકાઉન્ટ પર `ccam remove` સાદા `claude` માટે ફોલ્ડર અને તેનું લૉગિન રાખે છે. જો એકાઉન્ટ ગ્રુપમાં હોય, તો તે પહેલાં `ccam group leave` ની જેમ તેમને છોડે છે: દરેક શેર કરેલી આઇટમ ફોલ્ડરમાં ગ્રુપની આઇટમની કૉપિ બને છે, અને ગ્રુપના MCP સર્વર તેની `.claude.json` માં જાય છે. પછી તેનું રજિસ્ટ્રેશન રદ થાય છે: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. ccam એ બનાવેલા તેની `.claude.json` ના બૅકઅપ `<folder>/.ccam-backup/mcp-<time>` માં ખસે છે; જો તે ખસેડી ન શકાય, તો તે `~/.local/state/ccam/mcp/<name>/backups` માં રહે છે અને સંદેશ તેમ જણાવે છે. તમારું `~/.claude` trash માં ખસેડાતું નથી.
- બે એકાઉન્ટનું ફોલ્ડર એક ન હોઈ શકે: એક ફોલ્ડરનું નામ બે વાર આપતી `config.toml` લોડ થતી નથી.
- જેનું ફોલ્ડર `profiles/` માં પહેલેથી હોય તેવા નામનો `ccam add` ઇનકાર કરે છે, અને `ccam import` તરફ નિર્દેશ કરે છે. macOS પર, ccam ના જૂના વર્ઝને છોડેલું, નવા ફોલ્ડરના પાથ માટેનું લૉગિન કે MCP સાઇન-ઇન Keychain માં હજી હોય ત્યારે પણ તે અટકે છે, અને તેમને દૂર કરતો `security delete-generic-password` કમાન્ડ જણાવે છે. જો `config.toml` સેવ ન થઈ શકે, તો નવું ફોલ્ડર ફરી દૂર થાય છે.
- `ccam import` પહેલેથી રજિસ્ટર થયેલા ફોલ્ડરનો, ફોલ્ડર ન હોય તેવા પાથનો અને વપરાયેલા નામનો ઇનકાર કરે છે. ફોલ્ડર `~/.claude` ની કૉપિ હોય (એ જ લૉગિન) ત્યારે પણ ccam ચેતવે છે, ભલે `~/.claude` પોતે રજિસ્ટર ન હોય.
- એકાઉન્ટ ચાલી રહ્યું હોય ત્યારે `ccam logout` પહેલાં પૂછે છે; `ccam login` તપાસતું નથી.
- જો relay Claude Code સુધી ન પહોંચી શકે, ઉદાહરણ તરીકે લૉગિન પૂરું થઈ ગયા પછી, તો તે કારણ સાથે `ccam: could not pass the sign-in to Claude Code:` છાપે છે. HTTP ભૂલનો અર્થ સામાન્ય રીતે એ છે કે કોડની મુદત પૂરી થઈ ગઈ; લૉગિન ફરીથી શરૂ કરો.

## સંબંધિત

- [Claude Code શરૂ કરવું](/features/launch)
- [શૉર્ટકટ કમાન્ડ](/features/shortcuts)
- [MCP સર્વર, સેટિંગ્સ અને ઇતિહાસ શેર કરો](/docs/sharing)
- [એક કમ્પ્યુટર પર ઘણાં Claude Code એકાઉન્ટ વાપરો](/guides/multiple-claude-code-accounts)
