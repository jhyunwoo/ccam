---
title: ડેશબોર્ડ
description: ટર્મિનલ પર ccam એકાઉન્ટ, ગ્રુપ, ફોલ્ડર નિયમો, એક્સટેન્શન્સ અને તપાસનું ફુલ-સ્ક્રીન ડેશબોર્ડ ખોલે છે; દરેક ફેરફાર પહેલાં પોતાની યોજના બતાવે છે.
updated: 2026-10-30
order: 15
translation: machine
sourceHash: d65f1b5bb0c8
---

standard input અને output ટર્મિનલ હોય ત્યારે આર્ગ્યુમેન્ટ વગરનું `ccam` ફુલ-સ્ક્રીન ડેશબોર્ડ ખોલે છે; નહીં તો તે `ccam list` છાપે છે. ડેશબોર્ડ દરેક એકાઉન્ટને તેના વપરાશ અને ગ્રુપ સાથે બતાવે છે, અને કમાન્ડ જે કરે છે તે જ કરે છે: શરૂ કરવું, ઉમેરવું, લૉગિન કરવું, ગ્રુપ બદલવાં, ફોલ્ડર નિયમો સેટ કરવા, બૅકગ્રાઉન્ડ સેશનની યાદી બતાવવી, તપાસવું અને સુધારવું. શરૂ કરવાથી, કે બૅકગ્રાઉન્ડ સેશન સાથે જોડાવાથી, તે બંધ થાય છે; એકાઉન્ટ ઉમેરવા કે લૉગિન કરવા માટે તે થોડી વાર માટે ટર્મિનલ સોંપી દે છે, પછી ડેશબોર્ડ પર પાછું આવે છે.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · r to refresh · a to add · ? for shortcuts
```

ટૅબનાં નામ નીચેની લાઇન એક વાક્યમાં કહે છે કે હાલનું ટૅબ શેના માટે છે. વધુ પહોળા ટર્મિનલ પર Accounts ટૅબ દરેક એકાઉન્ટનો ઇમેઇલ પણ બતાવે છે. દરેક ટૅબ પર, ટૅબનાં નામવાળી લાઇનનો છેડો જણાવે છે કે [`auto_switch`](/features/next#auto-switch) ચાલુ છે કે નહીં (Linux અને macOS; Windows પર તેને અવગણવામાં આવે છે). તે ચાલુ હોય ત્યારે, Health ટૅબ તેને હજી શું જોઈએ છે તેની યાદી આપે છે: Claude Code 2.1.214 કે પછીનું, અને Claude subscription વડે સાઇન ઇન કરતા દરેક એકાઉન્ટમાં ચાલુ કરેલું ccam plugin.

દરેક વખતે ખૂલે ત્યારે ડેશબોર્ડ બૅકગ્રાઉન્ડમાં GitHub ને પૂછે છે કે ccam ની નવી રિલીઝ છે કે નહીં; પહેલી frame ક્યારેય તેની રાહ જોતી નથી, અને offline હોય ત્યારે કંઈ બતાવાતું નથી. નવી રિલીઝ હોય ત્યારે, શીર્ષક નીચેની એક લાઇન તેની ભલામણ કરે છે:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` પહેલાં પૂછે છે, પછી [`ccam update`](/features/update) ની જેમ અપડેટ કરે છે: રિલીઝને તેનાં SHA-256 checksums સામે તપાસવામાં આવે છે અને તે એક જ પગલામાં ccam ને બદલે છે. તમે જે ડેશબોર્ડમાં છો તે હજી જૂનું વર્ઝન છે, તેથી બહાર નીકળો અને `ccam` ફરી શરૂ કરો. જ્યારે અપડેટ અહીં થઈ શકતું ન હોય, દાખલા તરીકે `.deb` પેકેજમાંથી ઇન્સ્ટૉલ કરેલી ccam ની કૉપિ માટે, ત્યારે ડેશબોર્ડ `ccam update` જેવી જ સલાહ બતાવે છે અને ઑફર જાળવી રાખે છે.

પહેલી frame ફક્ત `config.toml` પરથી દોરાય છે, 100 ms થી ઘણા ઓછા સમયમાં; દરેક એકાઉન્ટ જવાબ આપે તેમ વપરાશ ભરાતો જાય છે. દરેક મર્યાદાની પોતાની bar છે: 5 કલાકની વિન્ડો, બધાં મૉડલ માટેની સાપ્તાહિક વિન્ડો, અને Fable ની પોતાની સાપ્તાહિક વિન્ડો; `-` નો અર્થ છે કે એકાઉન્ટ પાસે એવી વિન્ડો નથી. ટકાવારી 80% પર પીળી અને 100% પર લાલ થાય છે. lookup નિષ્ફળ જાય તો cache માંનાં મૂલ્યો ઝાંખાં બતાવાય છે. જે વિન્ડો ભરાઈ ગઈ હોવાનું Claude Code એ જણાવ્યું હોય તે તેના રીસેટ સુધી 100% બતાવે છે. સમાપ્ત થયેલો access token ઝાંખું `⚠ token expired` બતાવે છે (Claude Code તેને આગલી શરૂઆતે રિન્યૂ કરે છે), જ્યારે `⚠ login expired` નો અર્થ છે કે એકાઉન્ટને `ccam login` જોઈએ.

## ટૅબ

| ટૅબ | શું બતાવે છે | કી |
| --- | --- | --- |
| Accounts | નંબર, નામ, ઇમેઇલ, 5 કલાકની, સાપ્તાહિક અને Fable bar, છેલ્લો ઉપયોગ, ગ્રુપ; `✔` ડિફૉલ્ટ, `▸` આ ફોલ્ડરનું એકાઉન્ટ, `•` ચાલુ (ccam વગર શરૂ થયેલું Claude Code પણ: સાદું `claude`, IDE, બૅકગ્રાઉન્ડ સેશન અને તેમનો daemon) | નીચે |
| Agents | [`ccam agents`](/features/agents) જે બતાવે છે તે: દરેક એકાઉન્ટનાં બૅકગ્રાઉન્ડ સેશન અને ખુલ્લાં સેશન, તેમની સ્થિતિ, ફોલ્ડર અને શરૂઆતના સમય સાથે (નીચે) | `enter` જોડાઓ (`ccam attach` ની જેમ), `x` બંધ કરો, `r` ફરી લોડ કરો |
| Groups | દરેક ગ્રુપ, તે શું શેર કરે છે અને તેનાં એકાઉન્ટ, અને કોઈ ગ્રુપમાં ન હોય તેવાં એકાઉન્ટ | `enter` ફેરફાર, `n` નવું, `e` નામ બદલો, `x` દૂર કરો |
| Folders | [`ccam folder`](/features/folder) ના નિયમો અને અહીં લાગુ પડતો નિયમ, તમે ફોલ્ડર સુધી symbolic link દ્વારા પહોંચ્યા હો ત્યારે પણ (નવા નિયમો લિંક ઉકેલીને સેવ થાય છે) | `n` અહીં એકાઉન્ટ વાપરો, `enter` બદલો, `x` દૂર કરો |
| Extensions | [`ccam ext`](/features/ext) જે બતાવે છે તે, દરેક એક્સટેન્શન એક વાર, તે જે ગ્રુપમાંથી આવે છે તેની સાથે, `off:` જે એકાઉન્ટમાં plugin ચાલુ નથી તે, એકાઉન્ટના પોતાના સર્વર જે ફક્ત આ ફોલ્ડરને લાગુ પડે છે, `This folder's project`, અને કોઈ એકાઉન્ટની ન હોય તેવી સમસ્યાઓ | `r` ફરી લોડ કરો |
| Health | [`ccam doctor`](/features/doctor) જે જણાવે છે તે, દરેક એકાઉન્ટ માટે એક બૉક્સ તરીકે (નીચે) | `f` પસંદ કરેલું બૉક્સ સુધારો, `F` બધું સુધારો (`doctor --fix`), `↑`/`↓` બૉક્સ વચ્ચે ખસો |

`←`/`→` કે `Tab` ટૅબ બદલે છે, `?` હાલના ટૅબની કી બતાવે છે, `r` ટૂંકા cache ને બાયપાસ કરીને વપરાશ તાજો કરે છે અને Agents, Extensions અને Health ટૅબ ફરી લોડ કરે છે, અને `q` કે `Esc` બહાર નીકળે છે.

## Accounts ટૅબ પરની કી

| કી | શું કરે છે |
| --- | --- |
| `Enter` | પસંદ કરેલા એકાઉન્ટ તરીકે Claude Code શરૂ કરે છે |
| `1`–`9` | એ નંબરવાળા એકાઉન્ટ પર જાય છે; 12 માટે ઝડપથી `1` `2` લખો. ફક્ત `Enter` શરૂ કરે છે |
| `↑`/`↓`, `k`/`j` | ખસે છે |
| `n` | [`ccam next --auto`](/features/next) જે પસંદ કરે તે એકાઉન્ટ શરૂ કરે છે |
| `o` | ઓપ્શન સાથે શરૂ કરે છે: `--continue`, `--resume`, બીજા આર્ગ્યુમેન્ટ્સ, કે dry run |
| `u` | વપરાશની વિગતો: 5 કલાકની, સાપ્તાહિક અને મૉડલની વિન્ડો, તેમના રીસેટ અને આંકડા ક્યાંથી આવે છે |
| `a` | એકાઉન્ટ ઉમેરે છે (નીચેનો વિઝાર્ડ) |
| `l` | ફરીથી લૉગિન કરે છે |
| `s` | એકાઉન્ટનાં ગ્રુપ પસંદ કરે છે |
| `d` | તેને ડિફૉલ્ટ એકાઉન્ટ બનાવે છે |
| `f` | તેને હાલના ફોલ્ડરમાં વાપરે છે |
| `e` | તેનું નામ બદલે છે |
| `c` | તેનો શૉર્ટકટ કમાન્ડ સેટ કરે છે કે દૂર કરે છે |
| `x` | તેને દૂર કરે છે (પહેલાં પૂછે છે, No પસંદ કરેલું હોય છે) |

**Team અને Enterprise સીટ.** આપમેળે થતી એકાઉન્ટ પસંદગી (`ccam next`, `--fallback`, `ccam pick`, ડેશબોર્ડની `n` કી અને, તમે તેને ચાલુ કરો પછી, [`auto_switch`](/features/next#auto-switch)) એક જ Claude Team કે Enterprise સંસ્થાની સીટો વચ્ચે પણ સ્વિચ કરે છે, જેનાથી Anthropic ની શરતોનું ઉલ્લંઘન થઈ શકે છે. તેના પર આધાર રાખતાં પહેલાં [વારંવાર પૂછાતા પ્રશ્નો](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) વાંચો અને તમારી સંસ્થા અને Anthropic પાસેથી ખાતરી કરી લો.

## Health ટૅબ

Health ટૅબ ખોલતાં જ તે [`ccam doctor`](/features/doctor) ની તપાસો ચલાવે છે અને એ જ દૃશ્ય દોરે છે: દરેક એકાઉન્ટ માટે એક ચિહ્નવાળી પહેલી લાઇન, પછી દરેક એકાઉન્ટ માટે એક બૉક્સ. બૉક્સની દરેક હરોળ તપાસાયેલો એક વિસ્તાર છે, પહેલાં `Folder`, `Login` અને `Share links`, અને તેના પર `✔` બરાબર, `ℹ` માહિતી, `⚠` ચેતવણી કે `✘` ભૂલનું ચિહ્ન હોય છે, તેથી તમે તરત જોઈ શકો કે કયા એકાઉન્ટમાં સમસ્યા છે અને ક્યાં. કોઈ એકાઉન્ટનાં ન હોય તેવાં તારણ `General` બૉક્સમાં જાય છે.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` એક બૉક્સથી બીજા બૉક્સ પર ખસે છે, અને પસંદ કરેલા બૉક્સની કિનારી accent રંગની હોય છે. `f` ફક્ત પસંદ કરેલું બૉક્સ સુધારે છે, જે રીતે `ccam doctor --fix` સુધારે છે: કોઈ એકાઉન્ટના બૉક્સ પર તે ફક્ત એ એકાઉન્ટ સુધારે છે, અને `General` બૉક્સ પર કોઈ એકાઉન્ટનાં ન હોય તેવાં તારણ (plugin પાથ અને શૉર્ટકટ); બાકીનું દરેક એકાઉન્ટ ફક્ત ફરીથી તપાસાય છે. તે પહેલાં પૂછે છે (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), અને જેમાં સુધારવા જેવું કંઈ ન હોય તે બૉક્સ ફક્ત `nothing to repair in <name>` કહે છે. `F` બધું એકસાથે સુધારે છે, `ccam doctor --fix` ની જેમ, અને `r` ફરીથી તપાસે છે.

## Agents ટૅબ

Agents ટૅબ [`ccam agents`](/features/agents) જે છાપે છે તે બતાવે છે: દરેક એકાઉન્ટનાં બૅકગ્રાઉન્ડ સેશન અને ખુલ્લાં સેશન, દરેક જે એકાઉન્ટનું છે તેની સાથે. તમે ટૅબ પહેલી વાર ખોલો ત્યારે, `r` પર, સેશન બંધ કર્યા પછી, અને config માં ફેરફાર પછી, તે દરેક એકાઉન્ટના Claude Code ને એકસાથે પૂછે છે (`claude agents --json`, દરેક 15 સેકન્ડના timeout સાથે), અને બધાં જવાબ આપે ત્યાં સુધી `Asking every account's Claude Code for its sessions…` કહે છે.

```console
  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Background agents and open sessions of every account; enter attaches to one.

    Account   ID        Name   State          Folder      Started
 ❯  work      -         -      open (busy)    ~/code/web  now
    personal  49ec0028  probe  blocked, idle  ~/code/app  12m

 ✘ lab: could not list its sessions: exit status 1 'claude agents --json' is
   not available in this environment.

 enter to attach · x to stop · r to reload · ? for shortcuts
```

`State` બૅકગ્રાઉન્ડ સેશનની સ્થિતિ અને તેની process શું કરી રહી છે તે બતાવે છે, અથવા interactive સેશન માટે `open`, `ccam agents` ની જેમ. જો કોઈ એકાઉન્ટને પૂછી ન શકાય, તો યાદીની નીચે એક `✘` લાઇન Claude Code એ જે કહ્યું તે બતાવે છે, ઉપરના `lab` ની જેમ, જેનું agent view બંધ કરેલું છે. પૂરાં થયેલાં બૅકગ્રાઉન્ડ સેશન બાકાત રખાય છે (`ccam agents --all` તેમને બતાવે છે), અને જ્યારે દરેક એકાઉન્ટે જવાબ આપ્યો હોય અને કોઈ પાસે સેશન ન હોય, ત્યારે ટૅબ ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` કહે છે. સાંકડા ટર્મિનલ પર ટૅબ પહેલાં `Started`, પછી `Name`, પછી `Folder` કાઢી નાખે છે, અને લાંબા ફોલ્ડર પાથને ડાબી બાજુથી કાપે છે.

બૅકગ્રાઉન્ડ સેશન પર `Enter` ડેશબોર્ડ બંધ કરે છે અને સેશનને આ ટર્મિનલમાં તેના એકાઉન્ટ તરીકે ખોલે છે, જે રીતે `ccam attach <id>` કરે છે; ડેશબોર્ડ પાછું આવતું નથી (Claude Code માં, `←` તેના agent view પર અને `Ctrl+Z` shell પર લઈ જાય છે). ખુલ્લા સેશનનો કોઈ ID હોતો નથી, તેથી `Enter` અને `x` ફક્ત એટલું કહે છે કે બૅકગ્રાઉન્ડ સેશન જોઈએ. `x` પહેલાં પૂછે છે, `No, keep it` પસંદ કરેલું હોય છે: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` હા કહો તો, ccam એ એકાઉન્ટ તરીકે `claude stop <id>` ચલાવે છે, ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` કે તે કેમ ન થઈ શક્યું તે કહે છે, અને યાદી ફરી વાંચે છે.

## ગ્રુપ બદલવાં

Groups ટૅબ પર, `Enter` editor ખોલે છે: ડાબી બાજુ આઇટમ, જમણી બાજુ એકાઉન્ટ, બદલવા માટે `Tab`, ટૉગલ કરવા `Space`. જે આઇટમ કોઈ એકાઉન્ટને પહેલેથી બીજા ગ્રુપમાંથી મળતી હોય તે કારણ સાથે ઝાંખી બતાવાય છે, અને જે એકાઉન્ટને કોઈ આઇટમ બે વાર મળે તે પણ; એક આઇટમ ફક્ત એક જ ગ્રુપમાંથી આવી શકે છે. `Enter` યોજના બતાવે છે:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

જો કોઈ ફાઇલ ગ્રુપની કૉપિથી જુદી હોય, તો ccam પહેલાં પૂછે છે કે શું કરવું (ગ્રુપની રાખવી, મારી વાપરવી કે મર્જ કરવી). પછી ફેરફાર [`ccam group`](/features/share) ના જ કોડ સાથે બૅકગ્રાઉન્ડમાં ચાલે છે: પહેલાં config સેવ થાય છે, દરેક એકાઉન્ટનો ડેટા ક્રમમાં ખસે છે, અને મૂળ ફાઇલો `.ccam-backup/` માં જાય છે; MCP શેરિંગ છોડતું એકાઉન્ટ ગ્રુપના સર્વર તેની `.claude.json` માં રાખે છે.

જે ફેરફાર ચાલુ એકાઉન્ટનો ડેટા ખસેડે, તે એકાઉન્ટના નામ સાથે નકારાય છે, યોજના ખુલ્લી હતી ત્યારે એકાઉન્ટ શરૂ થયું હોય ત્યારે પણ, અને બીજી ccam process એ ગ્રુપ બદલ્યાં હોવાથી જૂની પડી ગયેલી યોજના `the groups changed while this plan waited; nothing was moved, make the change again` સાથે નકારાય છે. બીજી ccam process ગ્રુપ બદલતી હોય ત્યારે નવો ફેરફાર, નામ બદલવું કે દૂર કરવું `another ccam process is changing the groups right now` કહે છે; અધવચ્ચે અટકેલા ફેરફાર પછી, Health ટૅબ પર `F` (કે `ccam doctor --fix`) તેને પૂરો ન કરે ત્યાં સુધી તેઓ રાહ જુએ છે.

## એકાઉન્ટ ઉમેરવું

`a` એ વિઝાર્ડ ખોલે છે જે `ccam add` પણ ટર્મિનલ પર વાપરે છે:

1. એકાઉન્ટનું નામ.
2. તે કેવી રીતે સાઇન ઇન કરે છે: હમણાં Claude subscription (browser માં), [API key](/features/api), કે પછીથી Claude subscription.
3. API key માટે: ccam દરેક શરૂઆતે key ક્યાંથી વાંચે છે (`env:NAME` કે `command:COMMAND`), endpoint, key કેવી રીતે મોકલાય છે અને મૉડલ.
4. જોડાવા માટેનાં ગ્રુપ, `default` ચેક કરેલું (હજી કોઈ ગ્રુપ ન હોય, તો `default` ગ્રુપ શું શેર કરશે તે).
5. શૉર્ટકટ કમાન્ડ, જો તમારે જોઈતો હોય.
6. ખાતરી કરવા માટે સારાંશ.

`Esc` એક પગલું પાછળ જાય છે, અને પહેલા પગલે વિઝાર્ડ બંધ કરે છે.

## એકાઉન્ટ દૂર કરવું

`x` કંઈ પણ થાય તે પહેલાં પૂછે છે, No પસંદ કરેલું હોય છે, અને દૂર કરવાથી શું થશે તે જણાવે છે: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` API key એકાઉન્ટ માટે રદ કરવા જેવું કોઈ લૉગિન નથી, અને બીજું એકાઉન્ટ એ જ લૉગિન ધરાવતું હોય ત્યારે તે લૉગઆઉટ થતું નથી. ઇમ્પોર્ટ કરેલું એકાઉન્ટ પોતાનું ફોલ્ડર રાખે છે: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## રંગો

રંગો Claude Code ના ટર્મિનલ UI ની ભૂમિકાઓ અનુસરે છે, તેથી ccam તેની બાજુમાં સ્વાભાવિક લાગે છે; ગ્રુપનાં નામને તેમના નામ પરથી રંગ મળે છે. ccam ટર્મિનલને તેનું background પૂછે છે અને light કે dark theme પસંદ કરે છે. `config.toml` માં સૌથી ઉપર, પહેલા `[[group]]`, `[[profile]]` કે `[[folder]]` ની ઉપર, તમે જાતે પણ એક સેટ કરી શકો છો:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

`-daltonized` themes લાલ-લીલા રંગઅંધત્વ માટે છે અને `-ansi` themes ટર્મિનલના પોતાના 16 રંગો વાપરે છે. `NO_COLOR` રંગો બંધ કરે છે.

## સલામતી અને ખાસ કિસ્સા

- ડેશબોર્ડ કમાન્ડ જે functions બોલાવે છે તે જ બોલાવે છે, તેથી દરેક નિયમ લાગુ પડે છે: ડેટા ખસે તે પહેલાં યોજના, ચાલુ એકાઉન્ટ માટે ઇનકાર, બૅકઅપ, અને ક્યારેય એકાઉન્ટ ન છોડતા લૉગિન ટોકન.
- ડેશબોર્ડમાં dry run નથી: નામ વગરનું `ccam --dry-run` આ ફોલ્ડરનું એકાઉન્ટ કેવી રીતે શરૂ થશે તે છાપે છે, અને કોઈ એકાઉન્ટ પર `o` ફક્ત એ એક શરૂઆતનો dry run આપે છે.
- Claude Code શરૂ થાય કે લૉગિન ખૂલે તે પહેલાં ટર્મિનલ પાછું સામાન્ય થાય છે.
- dialog માટે નાની પડતી સ્ક્રીન પહેલાં ટૅબની હરોળો છોડે છે; લાઇનો wrap થવાને બદલે `…` સાથે કપાય છે. ઓછી ઊંચાઈની સ્ક્રીન પર તે પછી ખાલી લાઇનો અને ટૅબનું વર્ણન છોડાય છે, જેથી કીની સૂચનાવાળી લાઇન દેખાતી રહે. 40 કૉલમ કે 8 હરોળથી નાની સ્ક્રીન પર ડેશબોર્ડ ફક્ત `Too small: ccam needs 40x8 (now 36x7)` બતાવે છે.
- ફેરફાર, દૂર કરવું, અપડેટ કે સમારકામ લાગુ થતું હોય ત્યારે `Ctrl+C` તેની રાહ જુએ છે: સ્ટેટસ લાઇન `ccam quits when it is done.` કહે છે.
- `U` એ ccam અપડેટ કરી દીધા પછી, એકાઉન્ટ ઉમેરીને કે લૉગિન કરીને ડેશબોર્ડ પાછું આવે ત્યારે ઑફર ફરી દેખાતી નથી.

## સંબંધિત

- [ગ્રુપ વડે શેરિંગ](/features/share)
- [Claude Code શરૂ કરવું](/features/launch)
- [એકાઉન્ટનું સંચાલન](/features/accounts)
- [વપરાશ મર્યાદા](/features/usage)
