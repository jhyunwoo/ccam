---
title: Kubadili kwenda akaunti inayofuata
description: ccam next hufungua Claude Code kama akaunti nyingine; --auto huchagua nafasi kubwa zaidi, --continue huhamisha mazungumzo, --fallback hurudia kazi za -p.
updated: 2026-11-05
order: 8
translation: machine
sourceHash: "a320c160caef"
---

`ccam next` hufungua Claude Code kama akaunti nyingine bila wewe kuitaja. Itumie akaunti inapofikia kikomo chake cha matumizi. Kuna njia nne za kubadili:

- `ccam next` ya kawaida huchukua akaunti inayofuata ile uliyotumia mwisho, na `--auto` huchukua ile yenye nafasi kubwa zaidi. `--continue` huendeleza mazungumzo ya mwisho ya folda hii kwenye akaunti mpya, iwe hizo mbili zinashiriki historia au la.
- Ndani ya Claude Code, [plugin ya ccam](#inside-claude-code) hukujulisha kikomo kinaposimamisha kipindi, na `/ccam:switch` hukabidhi mazungumzo.
- Kwa script zinazoendesha `claude -p`, `ccam --fallback` hufanya ubadilishaji yenyewe (tazama [Endesha tena kwenye kikomo](#run-again-at-the-limit)).
- [`auto_switch`](#auto-switch) ikiwa imewashwa (Linux na macOS), ccam pia huhamisha kipindi shirikishi ambacho Claude Code inakisimamisha kwenye kikomo chake.

**Viti (seats) vya Team na Enterprise.** Uchaguzi wa akaunti wa kiotomatiki (`ccam next`, `--fallback`, `ccam pick`, kitufe `n` cha dashibodi na, ukishakiwasha, [`auto_switch`](/features/next#auto-switch)) pia hubadili kati ya viti vya shirika lile lile la Claude Team au Enterprise, jambo ambalo linaweza kukiuka masharti ya Anthropic. Soma [Maswali yanayoulizwa mara kwa mara](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) na uwasiliane na shirika lako na Anthropic kabla ya kulitegemea.

## Matumizi

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

ccam huchapisha ubadilishaji kama `→ <from> → <to>` kwenye standard error, pamoja na sababu ya chaguo la `--auto`, kisha huanzisha Claude Code kama akaunti mpya kwa njia ile ile ambayo `ccam <name>` hufanya. `<from>` ni akaunti uliyotumia mwisho, au kwa `--continue` akaunti ambayo mazungumzo ni yake. Hoja zinazokuja baada ya `--` huenda kwa `claude` bila kubadilishwa.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--auto` | Huchagua akaunti iliyoingia kwa matumizi (tazama hapa chini), ikiruka akaunti zilizofikia 100% | imezimwa |
| `--max-usage <1-100>` | Pamoja na `--auto`: hupendelea akaunti zilizo chini ya matumizi haya | `auto_max_usage`, vinginevyo 100 |
| `--strategy headroom\|expiring\|balance` | Pamoja na `--auto`: matumizi ya chini zaidi ya wiki ndani ya pointi 3 za asilimia kutoka matumizi ya chini zaidi ya saa 5, kipindi cha saa 5 kinachoanza upya mapema zaidi au kiasi kikubwa zaidi cha wiki kilichobaki kwa kila siku | `auto_strategy`, vinginevyo `headroom` |
| `--cross-org` | Pamoja na `--auto`: chagua pia akaunti za shirika jingine (tazama hapa chini) | `cross_org`, vinginevyo imezimwa |
| `--wait[=<duration>]` | Pamoja na `--auto`: hakuna akaunti yenye nafasi, husubiri ya kwanza ianze upya, si zaidi ya muda huu. Thamani huwekwa baada ya `=`: `--wait 2h` hukataliwa pamoja na kidokezo `--wait takes its value after '=': --wait=2h` | imezimwa; `5h` ikitolewa peke yake |
| `-c`, `--continue` | Huendeleza mazungumzo ya mwisho ya folda hii kwenye akaunti mpya (tazama hapa chini) | imezimwa |
| `--brief` | Pamoja na `--continue`: huanzisha mazungumzo mapya kama akaunti mpya kutoka kwa muhtasari wa ccam wa mazungumzo ya mwisho, badala ya kuyaendeleza yote ([tazama hapa chini](#brief)) | imezimwa |
| `-- <claude args>` | Hoja zinazopitishwa kwa `claude` baada ya zile za ccam yenyewe | hakuna |
| `--dry-run` | Inakubaliwa: huchapisha `→ <from> → <to>` na mstari wa amri ambao ingeendesha, bila kuanzisha Claude Code wala kurekodi uanzishaji | imezimwa |

`--auto` na `--continue` hufanya kazi pamoja. Kwa `--continue`, hoja za `--resume` huwekwa kabla ya hoja zako mwenyewe.

Chaguo za Claude Code yenyewe, kama vile `--dangerously-skip-permissions` au `--model`, huwekwa baada ya `--`. Tofauti na `ccam <name>`, ambayo hupitisha kila kitu kilicho baada ya jina, `ccam next` husoma hoja zake yenyewe, kwa hiyo bila `--` husimama na `ccam: unknown flag: --dangerously-skip-permissions`.

## Jinsi inavyofanya kazi

Kila uanzishaji kupitia ccam huandika muda katika `~/.local/state/ccam/last-used/<name>`. `ccam next` husoma faili hizi, hukagua ni akaunti zipi zimeingia, na huchagua.

1. ccam huunda `next.Candidate` moja kwa kila akaunti kutoka `config.toml` (`cli.candidates`): akaunti, muda wake wa mwisho wa kuanzishwa (`launcher.LastUsed`) na kama ina kuingia kunakofanya kazi (`creds.Has`: kuingia kwa Claude katika `.credentials.json`, au kwenye macOS katika Keychain, ambako hakujaisha muda kabisa wala kufutwa na Claude Code). Folda inayotumia API key badala yake (kuingia kwa Console, au akaunti iliyoongezwa kwa `ccam add --api`) haina vikomo vya mpango na haichaguliwi kamwe.
2. Bila `--auto`, ccam hupanga akaunti kwa namba ya slot na hutafuta ile iliyoanzishwa hivi karibuni zaidi (`next.Rotate`). Hiyo ndiyo akaunti unayoiacha. Ikiwa hakuna akaunti iliyowahi kuanzishwa, ni slot ya chini zaidi.
3. Kutoka hapo huendelea mbele kwa mpangilio wa slot, ikirudi kutoka slot ya mwisho hadi ya kwanza, na hurudisha akaunti ya kwanza iliyoingia.
4. ccam huandaa na kuanzisha akaunti iliyochaguliwa sawasawa kama `ccam <name>` ingefanya (`launcher.LaunchWith`): hukagua viungo vya akaunti, hutumia kushiriki kwa MCP, husajili kipindi na hubadilisha ccam kwa `claude` kupitia `syscall.Exec`.

Mzunguko huanza kwenye akaunti iliyotumiwa mwisho na husimama kwenye akaunti ya kwanza iliyoingia baada yake:

<!-- motion:next-rotate -->

`work` ndiyo iliyoanzishwa hivi karibuni zaidi, kwa hiyo mzunguko huanza baada yake. `side` haijaingia na inarukwa, na `lab` huchaguliwa. Kama `lab` nayo ingekuwa haijaingia, mzunguko ungerudi hadi `personal` (`% len` kwenye mchoro). Mzunguko kamwe haurudishi akaunti uliyoanzia: ikiwa hakuna akaunti nyingine iliyoingia, husimama na kosa.

### --auto

`--auto` huamua kwa matumizi badala ya mpangilio wa slot. ccam hubakiza akaunti zilizoingia, huomba matumizi yao yote kwa wakati mmoja kwa njia ile ile ambayo `ccam usage` hufanya (`Fetcher.FetchAll`), na huchagua kwa kanuni zilizo hapa chini (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

Akaunti hupimwa kwa vipindi vyake vya saa 5 na vya wiki, na, `--model` ikiwa miongoni mwa hoja za Claude Code (baada ya `--`), pia kwa kipindi cha wiki cha modeli hiyo (`opus`, `sonnet`, `haiku` au `fable` katika jina). ccam hupanga akaunti katika ngazi tatu (`Policy.Choose()`):

1. **Nafasi**: kila kipindi kiko chini ya `auto_max_usage` (`config.toml`, chaguo-msingi 100). Mkakati huchagua kati yao: `headroom` (chaguo-msingi) huchukua, kati ya akaunti zilizo ndani ya pointi 3 za asilimia kutoka matumizi ya chini zaidi ya saa 5, ile iliyotumia sehemu ndogo zaidi ya wiki yake, kwa hiyo akaunti yenye 12% hushinda ile yenye 10% ikiwa imetumia sehemu ndogo zaidi ya wiki yake, lakini ile yenye 14% haishindi; `expiring` huchukua kipindi cha saa 5 kitakachoanza upya mapema zaidi, ili kiasi chake kilichobaki kitumike kabla hakijapotea; `balance` huchukua akaunti yenye kiasi kikubwa zaidi cha wiki kilichobaki kwa kila siku hadi wiki yake ianze upya, ikihesabu kila kipindi cha ccam kinachoendeshwa dhidi yake, ili kiasi cha wiki kisibaki bila kutumika na vipindi sambamba vienee.
2. **Juu**: chini ya 100% lakini kwenye au juu ya `auto_max_usage`. Hutumiwa tu wakati hakuna akaunti yenye nafasi; ile ambayo `headroom` ingeichukua kati yao huja kwanza.
3. **Imejaa**: kipindi kiko 100%, au matumizi ya saa 5 hayajulikani. Haichaguliwi kamwe.

Claude Code ilipokataa ombi katika kipindi cha muda, kipindi hicho huhesabiwa kama 100% hadi kianze upya: hook ya [plugin ya ccam](#inside-claude-code) hukirekodi (`~/.local/state/ccam/usage/<name>.limit.json`), bila kujali namba zilizohifadhiwa zinasema nini. Kwa hiyo akaunti iliyofikia kikomo chake sasa hivi haichaguliwi tena kabla ya kuanza upya, hata namba zake zikiwa za zamani.

Mipangilio ya akaunti zote huwekwa juu kabisa ya `config.toml`, kabla ya `[[group]]`, `[[profile]]` au `[[folder]]` ya kwanza. Mpangilio ulio chini ya mojawapo ni wa jedwali hilo: ccam huonya kuuhusu kila inapoanzishwa, na amri zinazohifadhi usanidi hukataa hadi uhamishwe.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Mashirika.** Chaguo za kiotomatiki hubaki kati ya akaunti za aina ile ile kama unayoiacha: viti katika shirika lile lile la Team au Enterprise, au mipango ya binafsi (Pro, Max), ili kazi kutoka akaunti ya kampuni isiendelee kwenye akaunti ya binafsi au kinyume chake. ccam hujua shirika la kila akaunti kutokana na kuingia kwake (`ccam login` hulirekodi; kuingia kwa zamani husomwa kutoka `.claude.json`); akaunti ambayo ccam haiwezi kutambua shirika lake inaweza kulinganishwa na akaunti yoyote. Folda mbili zilizoingia kwenye akaunti ile ile ya Claude hushiriki vikomo vyake, kwa hiyo moja kamwe haichukui nafasi ya nyingine. `--cross-org`, au `cross_org = true` katika `config.toml`, huondoa kanuni ya shirika, lakini baada ya [`/ccam:switch`](#inside-claude-code) uhamisho kwenda shirika jingine bado husubiri useme ndiyo.

Isipowezekana kuchagua chochote, ccam hutoka na msimbo 75 (`EX_TEMPFAIL`), ili script iweze kutofautisha "jaribu baadaye" na kushindwa. Kwa `--wait`, badala yake hulala hadi kipindi kinachozuia akaunti ya kwanza iliyojaa kianze upya (sekunde 30 baada ya muda wa kuanza upya, endpoint inapoonyesha kipindi kipya), huchapisha `ccam: no account has room; waiting until 15:10`, na huchagua tena; ikiwa hakuna akaunti inayosema itaanza upya lini, huangalia tena baada ya dakika 10. Ikiwa hakuna kinachopata nafasi kabla ya kusubiri kuisha, hutoka na 75.

Utafutaji hufuata kanuni za `ccam usage`: namba zenye umri wa chini ya sekunde 60 hutumiwa tena, utafutaji wa kila akaunti hukoma baada ya sekunde 3, akaunti isiyojibu, au ambayo tokeni yake ya ufikiaji imeisha muda, hupimwa kwa namba zake zilizohifadhiwa, na rekodi ya kikomo huzipita namba hizo. Kipindi kilichohifadhiwa ambacho muda wake wa kuanza upya umepita huhesabiwa kama 0%, ili akaunti ambayo kipindi chake kimeanza upya tangu hapo isirukwe. `--auto` kamwe haifanyi upya tokeni. Tofauti na mzunguko wa kawaida, `--auto` inaweza kuchagua akaunti uliyotumia mwisho ikiwa bado ina nafasi kubwa zaidi.

### --continue

1. ccam hutafuta mazungumzo ya mwisho ya folda hii katika kila akaunti: kati ya faili za `<session>.jsonl` katika `projects/<folder>/` (jina la folda ni folda ya kazi, viungo vikiwa vimefuatwa, huku kila herufi isiyo herufi ya alfabeti wala tarakimu ikigeuzwa kuwa `-`, kama Claude Code inavyoipa jina), ile yenye prompt au jibu jipya zaidi. Vipindi visivyo navyo, kama kipindi kilichoendesha tu `/mcp`, `/usage` au `/status`, au faili tupu, hurukwa. Ikiwa kipindi kiliacha maelezo (kikomo cha matumizi au `/ccam:switch`, tazama hapa chini) na hakuna mazungumzo yenye prompt au jibu zaidi ya dakika 2 baada yake, kipindi hicho hasa ndicho hutumiwa badala yake; safu ambazo Claude Code huandika wakati wa `/exit` hazihesabiwi.
2. Akaunti ambayo mazungumzo ni yake ndiyo unayobadili kutoka kwayo, na akaunti mpya huchaguliwa kama ilivyo hapo juu, ndani ya shirika lile lile. Akaunti kadhaa zinaposhiriki historia hiyo, mazungumzo ni ya ile ambayo `.claude.json` yake inasema iliendesha mazungumzo hayo mwisho (`lastSessionId`), vinginevyo ni ya ile iliyotumiwa mwisho.
3. Akaunti mpya ikishiriki historia hiyo (folda ile ile ya kikundi), ccam huanzisha `claude --resume <session id>`, yale mazungumzo hasa (`--continue` ingefungua yoyote yaliyo mapya zaidi katika folda).
4. Vinginevyo huanzisha `claude --resume <the conversation's file> --fork-session`: Claude Code hunakili mazungumzo yote, chini ya session ID mpya, kwenye historia ya akaunti mpya yenyewe, na ya asili hubaki kama ilivyokuwa. Kuanzia hapo `--continue` ya akaunti mpya huyapata.
5. Usipotoa `--agent` baada ya `--`, ccam pia hupitisha agent ambaye mazungumzo yanaendeshwa kama yeye: mpangilio wa `agent` wa `.claude/settings.local.json` au `.claude/settings.json` ya folda au wa `settings.json` ya akaunti, vinginevyo agent anayetajwa na safu ya mwisho ya `agent-setting` ya mazungumzo. Claude Code hairudishi agent inaponakili mazungumzo (hatua ya 4), na ingeendelea kama agent wake chaguo-msingi, ikiwa na kila zana. Akaunti mpya isiposhiriki agents, plugins na mipangilio ya ile ya zamani (folda zile zile za `agents` na `plugins` na `settings.json` ile ile, kupitia kikundi), ccam huongeza mstari baada ya ule unaohusu uhamisho: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

ccam husema ukubwa wa mazungumzo: akaunti mpya huyasoma yote tena kwa jibu lake la kwanza. Prompt caches kamwe hazishirikiwi kati ya mashirika, na kila mpango wa binafsi ni shirika lake lenyewe, kwa hiyo akaunti hizo mbili zinapokuwa za mashirika tofauti ccam huongeza `prompt caches are not shared between organizations, so the first reply reads it all again`. Akaunti ya zamani inapokuwa imeingia kwenye seva za MCP ambazo mpya haijaingia, ccam huzitaja: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. Kushiriki `history` katika kikundi huepuka kunakili, huhifadhi mazungumzo mamoja, na pia hushiriki checkpoints ambazo `/rewind` hurejesha msimbo kutoka kwazo (`file-history/`).

### --brief

`--continue` huendeleza mazungumzo yote, na akaunti mpya huyasoma yote tena kwa jibu lake la kwanza. `ccam next --continue --brief` badala yake huanzisha mazungumzo mapya, kutoka kwa muhtasari ambao ccam huandika wa yale ya zamani bila kuuliza model:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

Muhtasari huwa na, kwa mpangilio huu na pale tu panapokuwa na la kusema:

1. kichwa cha mazungumzo;
2. njia ya transcript yake, ili Claude aitafute anapohitaji undani fulani;
3. ombi lako la kwanza, ikiwa yalikuwa zaidi ya manne;
4. muhtasari ambao Claude Code iliandika ilipobana (compact) mazungumzo mara ya mwisho;
5. maombi yako manne ya mwisho, pamoja na amri za `!` na amri za slash kama ulivyoziandika;
6. jibu la mwisho la Claude;
7. orodha ya kazi (todo);
8. faili 10 ambazo Claude alibadilisha mwisho;
9. `git status --short --branch` ya folda, hadi herufi 800.

Orodha ya kazi huhifadhi vipengele 10 kwa zaidi: vilivyokamilika huondolewa kwanza, kisha hatua zilizo mbali zaidi mbele, ili hatua ya sasa na zinazofuata zibaki. Sehemu ndefu hukatwa katikati, na muhtasari wote hubaki chini ya herufi 16,000.

ccam huuhifadhi kama `~/.local/state/ccam/brief/<session>.md`, unaosomeka na wewe pekee, na `--brief` inayofuata huondoa muhtasari wenye umri wa zaidi ya siku 7. Claude Code huanza na prompt moja, na huambatisha faili inayoitaja: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` Prompt unayotoa baada ya `--` huchukua nafasi ya sentensi ya mwisho (`Then: <your prompt>`), na chaguo zako nyingine baada ya `--` huifuata.

Claude Code humaliza utajo (mention) kwenye `"` na huukata kwenye `#`, kwa hiyo njia ya muhtasari inapokuwa na mojawapo, ccam huutaja kwa njia yake inayohusiana na folda ya kazi; njia hiyo nayo ikiwa na mojawapo, prompt hutaja faili, na Claude huisoma kwa zana yake ya Read, ambayo huuliza kwanza kwa sababu faili iko nje ya folda za kazi (katika uendeshaji wa `-p` kwenye hali ya Manual, Claude Code hukataa usomaji huo na Claude huendelea bila muhtasari).

`ccam --dry-run next --continue --brief` huchapisha muhtasari na haihifadhi chochote. Mazungumzo ya zamani hubaki kama yalivyokuwa, na mapya hujua tu kile ambacho muhtasari unasema, kwa hiyo tumia `--brief` kazi inapokuwa muhimu kuliko kila undani wa jinsi ulivyofika hapo. Bila `--continue`, ccam husimama na `--brief works with --continue`; isipoweza kusoma mazungumzo au isipopata chochote ndani yake cha kufupisha, husimama na kusema kwamba `ccam next --continue` bila `--brief` huyaendeleza yote. [Relay](#relay) na [`auto_switch`](#auto-switch) kila mara huendeleza mazungumzo yote. Chaguo la `--brief` la Claude Code yenyewe huwekwa baada ya `--`, kama chaguo zake nyingine.

## Ndani ya Claude Code

[Plugin ya ccam](/docs/install#claude-code-plugin) huongeza hooks mbili (Claude Code 2.1.78 au baadaye):

- Zamu inaposimama kwenye kikomo cha matumizi, hurekodi mazungumzo na huonyesha arifa ya desktop katika terminal zinazoiunga mkono (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  Hujua kila ujumbe ambao Claude Code huuhesabu kama kikomo cha matumizi: vikomo vya saa 5, vya wiki, vya kawaida na vya modeli (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` kwa kiasi cha modeli), credits za matumizi na vikomo vya matumizi ya fedha, na shirika lililosimamisha akaunti.

  Claude Code inapokuwa imeandika ombi lililokataliwa ndani ya mazungumzo, hook pia hurekodi kipindi kilichojaa hadi kianze upya ([rekodi ya kikomo](/features/usage#the-limit-record)), isipokuwa kwa usimamishaji uliowekwa na shirika, ambao wasimamizi wake wanaweza kuuondoa wakati wowote. Kizuizi kifupi cha seva au kikomo cha kasi cha API key si kikomo cha matumizi na huachwa, vivyo hivyo vipindi visivyo na mtu anayevifuatilia (uendeshaji wa `-p` na vipindi vya chinichini). Subagent au agent wa workflow anayesimama kwenye kikomo cha matumizi huhesabiwa kama usimamishaji wa kipindi chake, hata wakati mazungumzo makuu yanamsubiri tu: Claude Code huandika ombi hilo lililokataliwa ndani ya faili ya agent huyo mwenyewe katika folda ya `subagents` ya kipindi, na hook hulisoma hapo. Agents kadhaa wanaposimama kwa wakati mmoja, kipindi huandikiwa maelezo mara moja tu na hupata arifa moja tu.

  Claude Code 2.1.234 na baadaye pia husubiri ndani ya kipindi na huendelea yenyewe baada ya kipindi cha saa 5 kuanza upya; arifa hukupa chaguo la kubadili badala yake, na [`auto_switch`](#auto-switch) ikiwa imewashwa (Linux na macOS), ccam hukubadilishia akaunti nyingine inapokuwa na nafasi. Pale ambapo hakuna akaunti nyingine inayoweza kuchukua nafasi, kwa sababu hakuna akaunti nyingine iliyo ya shirika lile lile la Team au Enterprise, au hakuna iliyo mpango mwingine wa binafsi (`cross_org` ikiwa imewashwa, akaunti yoyote nyingine ya Claude inaweza), arifa badala yake huishia na `You have no other personal plans, so ccam cannot move this conversation.` au `You have no other accounts in Acme, so ccam cannot move this conversation.`
- `/ccam:switch [account]` hurekodi mazungumzo na humaliza zamu kabla ombi lolote halijatumwa: andika `/exit`, kisha endesha `ccam next --continue` ili kuyaendeleza kwenye akaunti uliyoitaja (au kwenye inayofuata). Katika kipindi ambacho bado hakina mazungumzo, husema hakuna cha kuendeleza. Ni kile unachoandika wewe tu kinachoiendesha: kuanzia plugin 0.3.2, Claude hawezi kuendesha `/ccam:switch`, na kazi ambayo Claude huipangia ratiba (CronCreate, ScheduleWakeup au `.claude/scheduled_tasks.json`) huipitisha kwa Claude kama maandishi ya kawaida.

  Akaunti unayoitaja inapokuwa nje ya shirika la mazungumzo (kiti cha shirika jingine la Team au Enterprise, mpango wa binafsi kwa mazungumzo yaliyo kwenye kiti au kinyume chake, au akaunti ya API key upande wowote), `ccam next --continue` huuliza kwanza, bila kujali `cross_org` inasema nini: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. Huuliza vivyo hivyo kabla ya akaunti inayoichagua yenyewe usipotaja yoyote, ambayo `cross_org = true` huiruhusu kuichukua kutoka shirika jingine: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. Ni `y` au `yes` pekee inayoyahamisha; kwenye Linux na macOS, vitufe vilivyobonyezwa kabla ya swali hutupwa kwanza. Jibu jingine lolote, au kukosekana kwa terminal ya kuulizia, huacha mazungumzo na maelezo pale yalipo: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (ikiwa akaunti ilichaguliwa na ccam, `ccam next --continue` huchagua tena). Kama ilivyo kwa chaguo za kiotomatiki, akaunti ambayo ccam haiwezi kutambua shirika lake hufaa yoyote.

Maelezo hayo huhifadhiwa katika `~/.local/state/ccam/handoff/<account>/<session>.json`, faili moja kwa kila kipindi, kwa hiyo vipindi viwili vya akaunti moja kamwe havifuti maelezo ya kingine. Huhifadhiwa kwa saa 12 na hutumiwa mara moja, kipindi kijacho kinapokaribia kuanza; uanzishaji unaoshindwa huyaacha kwa jaribio linalofuata. Ni ya folda ambamo kipindi kilianzia, hata baada ya `cd` katika shell ya Claude Code. `ccam rename` husasisha majina yaliyomo; ikiwa akaunti inayotajwa ndani yake imeondolewa tangu hapo, ccam husema hivyo na huchagua kama kawaida.

### Relay

Kwa `relay = true` juu kabisa ya `config.toml`, ccam huendelea kuendeshwa kando ya vipindi shirikishi badala ya kujibadilisha na Claude Code, na baada ya `/ccam:switch` huendeleza mazungumzo ya kipindi hicho yenyewe (maelezo yaliyoachwa na kipindi kingine hayafanyiwi kazi): hook humaliza Claude Code, na ccam huendesha `ccam next --continue` pamoja na maelezo ya kipindi hicho. Bila akaunti iliyotajwa, huchukua ile yenye nafasi kubwa zaidi, ikihesabu modeli ya mazungumzo, na `cross_org = true` inaporuhusu hiyo kuwa akaunti ya shirika jingine, huuliza swali lililo hapo juu kabla ya kuyahamishia mazungumzo huko. Kwenye Windows hook haiwezi kumaliza Claude Code: andika `/exit`, na ccam huyaendeleza. Wala hook haimalizi Claude Code kwa akaunti iliyo nje ya shirika la mazungumzo: husema `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`, na ccam huuliza swali lililo hapo juu mara tu unapokuwa umeandika `/exit`. Hilo hugharimu mchakato mmoja wa ccam unaosubiri kwa kila kipindi, ndiyo sababu umezimwa kwa chaguo-msingi.

Kipindi kinachoendelezwa huhifadhi kile ambacho kuendelea na mazungumzo (resume) hakurejeshi: chaguo kama `--add-dir`, `--mcp-config`, orodha za zana, `--settings`, `--plugin-dir` na chaguo za system prompt, hali ya ruhusa ya mazungumzo, kiwango cha juhudi (effort) cha kipindi na agent ambaye mazungumzo yanaendeshwa kama yeye (tazama [`--continue`](#continue)). Claude Code yenyewe hurejesha `--model` na `--name` kutoka kwenye mazungumzo. Kitu kingine chochote (`--worktree`, `--debug`, `--remote-control`, chaguo ambazo ccam haizijui) hutajwa katika `ccam: not carried over: …`. Ubadilishaji kamwe hauulizi swali la historia inayoshirikiwa, na ulichoandika wakati ubadilishaji ukiendelea hutupwa; swali lake pekee ni lile linaloulizwa kabla ya akaunti ya shirika jingine.

Relay huendeleza mazungumzo tu wakati Claude Code ilimalizika kwa njia ambayo `/ccam:switch` iliiomba: si ulipoandika `/exit` kwanza, ulipobonyeza Ctrl+C ilipokuwa ikimalizika, au ccam ilipopata hangup au ishara nyingine, na si wakati terminal imetoweka. Hapo maelezo hubaki, na `ccam next --continue` huyaendeleza baadaye. Kipindi kilichoanzishwa chini ya relay ya toleo la zamani la ccam, au kisicho na mtu anayekifuatilia (uendeshaji wa `-p` na vipindi vya chinichini), hakimalizwi na hook; badala yake hupata ujumbe wa `/exit`. `claude remote-control` huendeshwa bila relay.

### `auto_switch`

Kwa `auto_switch = true` juu kabisa ya `config.toml` (Linux na macOS), ccam huhamisha kipindi kwenda akaunti nyingine yenyewe Claude Code inapokisimamisha kwenye kikomo cha matumizi: hook ya plugin huchagua akaunti nyingine yenye nafasi na humaliza Claude Code, na ccam, inayoendeshwa kando ya kipindi, huendeleza mazungumzo huko, kama [relay](#relay) inavyofanya baada ya `/ccam:switch`. Imezimwa kwa chaguo-msingi na ni ya majaribio.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam huhamisha kipindi pale tu yote haya yanapotimia:

1. ccam ilikianzisha katika terminal na huendelea kuendeshwa kando yake. `auto_switch` huwasha [relay](#relay), kwa hiyo hili hutimia kwa vipindi shirikishi unavyoanzisha baada ya kuiweka. Kipindi kilichoanzishwa kabla, au chini ya ccam ya zamani, hupata arifa inayoishia na ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` Vipindi visivyo na mtu anayevifuatilia (uendeshaji wa `-p` na vipindi vya chinichini), `claude remote-control` na Claude Code iliyoanzishwa bila ccam kamwe havihamishwi.
2. Akaunti ina [plugin ya ccam](/docs/install#claude-code-plugin) ikiwa imewashwa na huendesha Claude Code 2.1.214 au baadaye, ambayo hupitisha ID ya mchakato wake kwa hook (`CLAUDE_PID`).
3. Claude Code iliandika ombi lililokataliwa ndani ya mazungumzo, pamoja na muda linapoanza upya, kwa kikomo cha saa 5, cha wiki, cha Opus au cha Sonnet cha akaunti hiyo. Kikomo cha matumizi ya fedha, credits za matumizi kuisha, usimamishaji uliowekwa na shirika, kiasi cha Fable (`You've reached your Fable limit.`), kizuizi cha seva na kikomo cha kasi cha API key kamwe havihamishi kipindi. Kikomo cha subagent au cha agent wa workflow huhesabiwa, hata wakati mazungumzo makuu yanamsubiri tu; kinapokuwa kikomo cha wiki cha Opus au cha Sonnet, kipindi huhamishwa tu ikiwa mazungumzo yenyewe yanaendeshwa kwenye modeli hiyo.
4. ccam inajua mpango wa akaunti, kutoka `.claude.json` yake au kutoka kile ambacho `ccam login` ilirekodi. Isipojua, arifa huishia na `ccam does not know personal's plan, so auto_switch leaves this session.`
5. Kipindi hakijaunganishwa na Remote Control. Kama akaunti nyingine kisingeunganishwa tena na programu unayokiendesha kutoka kwayo, kwa hiyo hubaki, na arifa husema hivyo.
6. Kikomo huanza upya zaidi ya dakika 15 kuanzia sasa. Kwa kusubiri kufupi zaidi ccam huiachia Claude Code kipindi: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. Hujaandika prompt tangu kikomo.
8. Akaunti nyingine ina nafasi. ccam huchagua kama [`--auto`](#auto) inavyofanya (`auto_max_usage`, `auto_strategy`, model ya mazungumzo), kati ya akaunti za aina ile ile unayoiacha: kiti kingine cha shirika lile lile la Team au Enterprise, au mpango mwingine wa binafsi (Pro, Max). `cross_org` na `--cross-org` hazitumiki hapa. Akaunti lazima:
   - iwe imeingia kwa kuingia ambako Claude Code inaweza kufanya upya;
   - isiwe imeingia kwenye akaunti ile ile ya Claude;
   - isitumie API key;
   - iwe na mpango ambao ccam inaujua;
   - isiwe na credits za matumizi zilizowashwa, kwa kadiri utafutaji wa matumizi wa ccam yenyewe unavyoonyesha;
   - iwe imeiamini (trust) folda hii katika Claude Code;
   - ikiwa mazungumzo yanaendeshwa kama agent (mpangilio wa `agent` au mazungumzo yanamtaja mmoja), ishiriki agents wa akaunti unayoiacha: folda zake za `agents` na `plugins` na `settings.json` yake lazima kila moja iwe ile ile hasa (kupitia kikundi), au isiwepo katika zote mbili.

   Masharti hayo matatu ya mwisho yanapoacha nje kila akaunti yenye nafasi, arifa huzitaja: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` au `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` Kipindi chenye [rekodi ya kikomo](/features/usage#the-limit-record) huhesabiwa kama kimejaa, na akaunti ambayo kipindi kingine kilihamia ndani ya dakika 2 zilizopita huja baada ya nyingine. Uchaguzi lazima ufanyike ndani ya takriban sekunde 6; utafutaji wa matumizi unapochukua muda mrefu zaidi, kipindi hubaki. Kwa `usage_source = "claude"`, `auto_switch` huamua kutoka kwa namba zilizohifadhiwa pekee, na ccam haiwezi kujua kutoka kwenye utafutaji wa Claude Code wala mstari wa hali kama credits za matumizi zimewashwa, kwa hiyo huchukulia kuwa zimezimwa.

Kisha hook humaliza Claude Code, na ccam, inayoendeshwa kando yake, huonyesha arifa ya desktop katika terminal zinazoiunga mkono, kama `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` Huchapisha mstari wa kwanza wa matokeo ya mfano hapo juu na huendesha `ccam next --continue` kwa kipindi hicho kama akaunti iliyochaguliwa, ambayo huchapisha mistari yake mingine miwili na huanzisha Claude Code pamoja na chaguo za kipindi, hali ya ruhusa, kiwango cha juhudi (effort) na agent, kama [relay](#relay) inavyofanya. Ujumbe wa kwanza wa kipindi kipya ni maandishi yasiyobadilika, kitu pekee ambacho ccam hutuma kwa niaba yako (uhamisho baada ya `/ccam:switch` hautumi chochote):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

Baada ya kikomo cha subagent ni maandishi haya, kwa sababu subagents, workflows na amri za chinichini zilisimama pamoja na kipindi cha zamani:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

Akaunti hiyo nayo ikifikia kikomo chake, jambo hilo hilo hutokea huko. Rekodi za vikomo huzuia mnyororo usirudi kwenye akaunti kabla haijaanza upya, kwa hiyo kila akaunti ikiwa imejaa kipindi hubaki kwenye ya mwisho.

Hook hujua kwamba mazungumzo yanaendeshwa kama agent kutokana na mpangilio wa `agent`, au kutokana na safu inayomtaja agent ambayo Claude Code huiandika ndani ya mazungumzo mara kwa mara na kipindi kinapoisha. Ikiwa hadi sasa ni `--agent` ya kipindi pekee inayomtaja, na akaunti iliyochaguliwa haishiriki agents, `ccam next --continue` haianzishi chochote na huhifadhi maelezo, na arifa ya pili ya desktop husema hivyo: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

Pale ambapo hakuna akaunti nyingine yenye nafasi, hakuna kinachohamishwa: Claude Code hushughulikia kikomo kama kawaida, ikisubiri ndani ya kipindi au ikitoa chaguo zake, na arifa husema ni akaunti ipi itakayopata nafasi kwanza, kwa mfano `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` Kwa viti husema `None of your other accounts in Acme has room now`. Ili kusimamisha uhamisho ulioanza, bonyeza Ctrl+C wakati Claude Code inafungwa: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

Uhamisho una gharama. Claude Code huanza upya, jambo linalochukua sekunde chache, na shell za chinichini, seva za maendeleo (dev servers) na agents ambazo Claude alianzisha husimama pamoja na kipindi cha zamani. Akaunti mpya husoma mazungumzo yote tena kwa jibu lake la kwanza, kama ilivyo kwa [`--continue`](#continue), na ccam husema ukubwa wake; uhamisho kamwe hautumii [`--brief`](#brief).

`auto_switch` kamwe haibadili kabla Claude Code haijaripoti kikomo, kamwe haifuatilii matumizi yako chinichini, na kamwe hainakili, haifanyi upya wala haipitishi taarifa za kuingia. Kwenye Windows hupuuzwa, na `ccam doctor` husema hivyo; huko, kwa `relay = true`, `/ccam:switch` na `/exit` bado huhamisha mazungumzo kwenda akaunti nyingine.

**Masharti ya Anthropic.** Kwenye viti vya Team na Enterprise, `auto_switch` huhamisha kazi yako kwenda kiti kingine cha shirika lile lile. Kituo cha usaidizi cha Anthropic kinasema "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (kila kiti katika mpango wa Team hutolewa kwa mtu mmoja aliyetajwa kwa jina, na ni mtu huyo pekee anayeruhusiwa kukitumia) na "Users can't route their work through another user's account to use that user's remaining usage" (watumiaji hawaruhusiwi kuelekeza kazi yao kupitia akaunti ya mtumiaji mwingine ili kutumia kiasi cha matumizi kilichobaki cha mtumiaji huyo) ([mipango ya Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); vivyo hivyo kwa [mipango ya Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)).

Kwa mipango ya binafsi (Pro, Max), masharti ya Anthropic hayasemi kama akaunti zako mwenyewe zinaweza kuchukua nafasi ya nyingine kiotomatiki kwenye kikomo. [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) hukataza kufikia Huduma "through automated or non-human means, whether through a bot, script, or otherwise" (kwa njia za kiotomatiki au zisizo za binadamu, iwe kupitia bot, script au vinginevyo) isipokuwa Anthropic iruhusu waziwazi, na uhamisho hutuma mmoja wa jumbe zisizobadilika zilizo hapo juu kwa niaba yako. [Usage Policy](https://www.anthropic.com/legal/aup) inayoanza kutumika tarehe 12 Novemba 2026 huorodhesha "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (kukwepa au kuingilia kwa makusudi uwezo, vizuizi au kinga zilizowekwa ndani ya bidhaa zetu) miongoni mwa mambo inayoyakataza.

ccam haidai kwamba `auto_switch` inatii masharti haya. Kuiwasha ni uamuzi wako, na hatari pia ni yako: Anthropic inaweza kusimamisha akaunti, au ufikiaji wa shirika la Team au Enterprise, inapoamini masharti yamevunjwa.

Wakati `auto_switch` imewashwa, `ccam doctor` huonya kuhusu akaunti isiyo na plugin ya ccam na kuhusu matoleo ya Claude Code ya zamani kuliko 2.1.214, na kwa viti hurudia kanuni ya viti ya Anthropic ([ukaguzi](/features/doctor)).

## Endesha tena kwenye kikomo

`ccam next` ni kwa kipindi shirikishi kilichosimama. Uendeshaji usio shirikishi, `claude -p` katika script au kazi iliyopangwa, hauwezi kukusubiri, kwa hiyo `--fallback` kabla ya jina la akaunti huiruhusu ccam ibadili kwa ajili yake:

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

1. `--fallback` hufanya kazi tu pamoja na `-p` au `--print`; vinginevyo ccam husimama na huelekeza kwa `ccam next --auto --continue`.
2. Standard input inayokuja kupitia pipe husomwa mara moja, hadi MiB 10, na hupewa kila jaribio. Kama Claude Code, ccam husubiri si zaidi ya sekunde 3 ianze kufika (isipokuwa kwa `--input-format stream-json`); pipe ikibaki wazi bila kutuma chochote, uendeshaji haupati ingizo na ccam huandika hilo kwenye standard error.
3. ccam huanzisha `claude` kama mchakato mtoto badala ya `exec`, huzuia standard output yake, na hupitisha standard error kama inavyokuja.
4. `claude` inapotoka kwa 0, au inaposhindwa kwa sababu nyingine, output iliyozuiwa na exit code hupitishwa bila kubadilishwa. Vivyo hivyo uendeshaji uliokomeshwa na ishara (exit code iliyo juu ya 128): uliombwa kusimama, na hauendeshwi tena.
5. Kushindwa huhesabiwa kama kikomo cha matumizi ujumbe wa Claude Code unaposema hivyo, yaani ujumbe wowote wa kikomo cha matumizi ambao [hook](#inside-claude-code) inaujua, hata ndani ya output ya `--output-format json` au `stream-json`, au matumizi ya akaunti yanapoonyesha kipindi chake cha saa 5, cha wiki au cha `--model` kiko 100% (pamoja na rekodi ya kikomo).
6. Kisha output ya jaribio hilo hutupwa na amri ile ile huendeshwa kama akaunti nyingine iliyoingia: ile ambayo `--auto` ingeichagua, kisha akaunti ambazo matumizi yao hayawezi kupatikana, kwa mpangilio wa namba. Akaunti zilizofikia 100% hurukwa, akaunti zenye kuingia kwa API key hazijaribiwi kamwe, na kila akaunti hujaribiwa si zaidi ya mara moja.
7. Kukiwa na `-c`, `--continue`, `-r` au `--resume` miongoni mwa hoja, akaunti zilizo katika kikundi kile kile cha historia pekee ndizo hujaribiwa, ili mazungumzo yaweze kupatikana.
8. Akaunti za shirika lile lile pekee ndizo hujaribiwa, kama ilivyo kwa `--auto` (`cross_org = true` huondoa hili).

Akaunti inayofuata huchaguliwa kwa kanuni zile zile kama `--auto` (`auto_max_usage`, `auto_strategy`, kipindi cha `--model`). Hakuna akaunti inapobaki, ccam huchapisha output ya jaribio la mwisho, husema kwamba hakuna akaunti iliyoingia yenye nafasi iliyobaki, na hutoka kwa code 75. Kwa `--wait` kabla ya jina la akaunti (`ccam --fallback --wait=2h work -p …`), badala yake husubiri akaunti ya kwanza ianze upya, kisha inaweza kujaribu kila akaunti tena; kila mara hulala angalau mara moja kwanza, ili akaunti iliyofikia kikomo chake sasa hivi isiendeshwe tena mara moja. Muda huo huweka kikomo cha jumla ya muda wote wa kusubiri wa uendeshaji, si wa kila kusubiri. Ukiandikwa kwa nafasi (`--wait 30m`), hukataliwa pamoja na kidokezo kile kile kama cha `ccam next`. `--` ya mwanzo kabla ya hoja za Claude Code hufanya kazi kama katika uanzishaji wowote (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` huchapisha amri ya jaribio la kwanza.

Vipindi shirikishi havibadilishwi vyenyewe isipokuwa uwashe [`auto_switch`](#auto-switch). Kwa `ccam statusline install`, mstari wa hali huongeza `→ ccam next --auto --continue` matumizi ya saa 5 au ya wiki yanapofikia 90% ([Vikomo vya matumizi](/features/usage)).

## Usalama na hali maalum

- Akaunti ambazo hazijaingia hazichaguliwi kamwe. Ikiwa hakuna akaunti nyingine iliyoingia, `ccam next` ya kawaida husimama na `no logged-in account other than <name>`.
- Kwa `--continue` na hakuna mazungumzo katika folda: `no account has a conversation in <folder> to continue`.
- Kanuni ya shirika isipoacha akaunti yoyote: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Kwa `--auto` na hakuna akaunti iliyoingia: `no logged-in account to pick`. Kila akaunti ikiwa imejaa au matumizi yake hayajulikani: `` no account with known usage and headroom (see `ccam usage`) ``. Zote mbili huishia na exit code 75.
- `--max-usage`, `--strategy`, `--wait` na `--cross-org` hufanya kazi tu pamoja na `--auto`, na thamani nje ya 1-100 au mkakati usiojulikana hukataliwa.
- Fork hunakili mazungumzo, si hali nyingine ya akaunti: checkpoints za `file-history/` ya akaunti ya asili haziji isipokuwa historia inashirikiwa.
- `--auto` hufanya maombi ya mtandao hata kwa `--dry-run`. Husoma tu, na kamwe haiandiki faili za tokeni.
- Akaunti iliyochaguliwa huanza kupitia njia ile ile kama `ccam <name>`, kwa hiyo ukaguzi ule ule hutumika, pamoja na onyo wakati kipindi kingine cha Claude Code kinachoshiriki historia ile ile tayari kinaendeshwa katika folda hii.
- `last-used` huandikwa tu Claude Code inapoanza kweli. `--dry-run` huiacha bila kuigusa, kwa hiyo jaribio halisogezi mzunguko mbele.
- `--fallback` huzuia standard output hadi jaribio liishe, kwa hiyo output ya mtiririko kama `--output-format stream-json` hufika kwa kipande kimoja mwishoni. Mteja anayetuma ingizo la `stream-json` kwa mtiririko na kusubiri majibu kabla ya kutuma zaidi hawezi kutumia `--fallback`. Faili ambazo jaribio la kwanza lilibadilisha kabla ya kufikia kikomo hubaki zimebadilishwa; jaribio linalofuata huanza prompt upya.

## Zinazohusiana

- [Vikomo vya matumizi](/features/usage)
- [Vikomo vya matumizi na kubadili akaunti](/docs/usage-limits)
- [Shiriki seva za MCP, mipangilio na historia](/docs/sharing)
- [Cha kufanya Claude Code ikifikia kikomo cha matumizi](/guides/usage-limit-workflow)
