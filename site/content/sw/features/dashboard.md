---
title: Dashibodi
description: ccam kwenye terminal hufungua dashibodi ya skrini nzima ya akaunti, vikundi, kanuni za folda, viendelezi na ukaguzi; kila badiliko huonyesha mpango wake kwanza.
updated: 2026-10-30
order: 15
translation: machine
sourceHash: "d65f1b5bb0c8"
---

`ccam` bila hoja hufungua dashibodi ya skrini nzima wakati ingizo na matokeo ya kawaida ni terminal; vinginevyo huchapisha `ccam list`. Dashibodi huonyesha kila akaunti pamoja na matumizi yake na vikundi vyake, na hufanya kile ambacho amri hufanya: kuanzisha, kuongeza, kuingia, kubadilisha vikundi, kuweka kanuni za folda, kuorodhesha vipindi vya chinichini, kukagua na kurekebisha. Kuanzisha, au kuunganishwa na kipindi cha chinichini, huifunga; kuongeza akaunti au kuingia hukabidhi terminal kwa muda mfupi, kisha hurudi kwenye dashibodi.

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

Mstari ulio chini ya majina ya vichupo husema kwa sentensi moja kichupo cha sasa ni cha nini. Kwenye terminal pana zaidi, kichupo cha Accounts pia huonyesha barua pepe ya kila akaunti. Kwenye kila kichupo, mwisho wa mstari wa vichupo huonyesha kama [`auto_switch`](/features/next#auto-switch) imewashwa (Linux na macOS; kwenye Windows hupuuzwa). Ikiwa imewashwa, kichupo cha Health huorodhesha bado inahitaji nini: Claude Code 2.1.214 au baadaye, na plugin ya ccam ikiwa imewashwa katika kila akaunti inayoingia kwa usajili wa Claude.

Kila inapofunguliwa, dashibodi huuliza GitHub chinichini kama kuna toleo jipya zaidi la ccam; fremu ya kwanza kamwe haisubiri jibu hilo, na ukiwa nje ya mtandao, hakuna kinachoonyeshwa. Likiwepo, mstari ulio chini ya kichwa hulipendekeza:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` huuliza kwanza, kisha husasisha kama [`ccam update`](/features/update) inavyofanya: toleo hukaguliwa dhidi ya checksum zake za SHA-256 na huchukua nafasi ya ccam kwa hatua moja. Dashibodi uliyomo bado ni toleo la zamani, kwa hiyo toka kisha uanzishe `ccam` tena. Sasisho lisipoweza kufanywa hapa, kwa mfano kwa nakala ya ccam iliyosakinishwa kutoka kifurushi cha `.deb`, dashibodi huonyesha ushauri ule ule kama `ccam update` na huendelea kulitoa toleo hilo.

Fremu ya kwanza huchorwa kutoka `config.toml` pekee, kwa muda mfupi sana chini ya 100 ms; matumizi hujazwa kadiri kila akaunti inavyojibu. Kila kikomo kina upau wake: dirisha la saa 5, dirisha la wiki kwa model zote, na dirisha la wiki la Fable lenyewe; `-` inamaanisha akaunti haina dirisha hilo. Asilimia hugeuka manjano kwenye 80% na nyekundu kwenye 100%. Utafutaji ukishindwa, thamani kutoka kwenye akiba (cache) huonyeshwa kwa rangi hafifu. Dirisha ambalo Claude Code iliripoti kuwa limejaa huonyesha 100% hadi liwekwe upya. Tokeni ya ufikiaji iliyoisha muda huonyesha `⚠ token expired` hafifu (Claude Code huifanya upya katika uanzishaji unaofuata), ilhali `⚠ login expired` inamaanisha akaunti inahitaji `ccam login`.

## Vichupo

| Kichupo | Kinachoonyesha | Vitufe |
| --- | --- | --- |
| Accounts | namba, jina, barua pepe, pau za saa 5, za wiki na za Fable, matumizi ya mwisho, vikundi; `✔` akaunti chaguo-msingi, `▸` akaunti ya folda hii, `•` inaendeshwa (pia Claude Code iliyoanzishwa bila ccam: `claude` ya kawaida, IDE, vipindi vya chinichini na daemon yao) | hapa chini |
| Agents | kile ambacho [`ccam agents`](/features/agents) huorodhesha: vipindi vya chinichini na vipindi vilivyo wazi vya kila akaunti, pamoja na hali yake, folda na muda wa kuanza (hapa chini) | `enter` unganisha (kama `ccam attach`), `x` simamisha, `r` pakia upya |
| Groups | kila kikundi, kile kinachoshiriki na akaunti zake, pamoja na akaunti zisizo katika kikundi chochote | `enter` hariri, `n` kipya, `e` badilisha jina, `x` ondoa |
| Folders | kanuni za [`ccam folder`](/features/folder) na ile inayotumika hapa, hata ulipofika kwenye folda kupitia symbolic link (kanuni mpya huhifadhiwa kwa njia halisi, viungo vikiwa vimefuatwa) | `n` tumia akaunti hapa, `enter` badilisha, `x` ondoa |
| Extensions | kile ambacho [`ccam ext`](/features/ext) huonyesha, kila kiendelezi mara moja pamoja na vikundi kinakotoka, `off:` akaunti ambazo plugin haijawashwa, seva za akaunti yenyewe zinazotumika kwa folda hii pekee, `This folder's project`, na matatizo yasiyo ya akaunti yoyote | `r` pakia upya |
| Health | kile ambacho [`ccam doctor`](/features/doctor) huripoti, kama kisanduku kimoja kwa kila akaunti (hapa chini) | `f` rekebisha kisanduku kilichochaguliwa, `F` rekebisha vyote (`doctor --fix`), `↑`/`↓` sogea kati ya visanduku |

`←`/`→` au `Tab` hubadilisha vichupo, `?` huorodhesha vitufe vya kichupo cha sasa, `r` huonyesha upya matumizi, ikipita akiba fupi, na hupakia upya vichupo vya Agents, Extensions na Health, na `q` au `Esc` hufunga.

## Vitufe kwenye kichupo cha Accounts

| Kitufe | Inachofanya |
| --- | --- |
| `Enter` | Anzisha Claude Code kama akaunti iliyochaguliwa |
| `1`–`9` | Nenda kwenye akaunti yenye namba hiyo; andika `1` `2` haraka kwa 12. `Enter` pekee ndicho huanzisha |
| `↑`/`↓`, `k`/`j` | Sogea |
| `n` | Anzisha akaunti ambayo [`ccam next --auto`](/features/next) ingechagua |
| `o` | Anzisha kwa chaguo: `--continue`, `--resume`, hoja nyingine, au jaribio (dry run) |
| `u` | Maelezo ya matumizi: madirisha ya saa 5, ya wiki na ya model pamoja na nyakati zake za kuwekwa upya na mahali namba zinapotoka |
| `a` | Ongeza akaunti (mchawi ulio hapa chini) |
| `l` | Ingia tena |
| `s` | Chagua vikundi vya akaunti |
| `d` | Ifanye akaunti chaguo-msingi |
| `f` | Itumie katika folda ya sasa |
| `e` | Badilisha jina lake |
| `c` | Weka au ondoa amri yake ya njia fupi |
| `x` | Iondoe (huuliza kwanza, No ikiwa imechaguliwa) |

**Viti (seats) vya Team na Enterprise.** Uchaguzi wa akaunti wa kiotomatiki (`ccam next`, `--fallback`, `ccam pick`, kitufe `n` cha dashibodi na, ukishakiwasha, [`auto_switch`](/features/next#auto-switch)) pia hubadili kati ya viti vya shirika lile lile la Claude Team au Enterprise, jambo ambalo linaweza kukiuka masharti ya Anthropic. Soma [Maswali yanayoulizwa mara kwa mara](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) na uwasiliane na shirika lako na Anthropic kabla ya kulitegemea.

## Kichupo cha Health

Kichupo cha Health huendesha ukaguzi wa [`ccam doctor`](/features/doctor) unapokifungua na huchora mwonekano ule ule: mstari wa kwanza wenye alama kwa kila akaunti, kisha kisanduku kimoja kwa kila akaunti. Kila safu katika kisanduku ni eneo moja lililokaguliwa, `Folder`, `Login` na `Share links` kwanza, likiwa na alama `✔` sawa, `ℹ` taarifa, `⚠` onyo au `✘` hitilafu, kwa hiyo unaona mara moja ni akaunti ipi ina tatizo na wapi. Matokeo yasiyo ya akaunti yoyote huenda kwenye kisanduku cha `General`.

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

`↑`/`↓` husogea kutoka kisanduku kimoja hadi kingine, na kisanduku kilichochaguliwa kina mpaka wa rangi ya msisitizo. `f` hurekebisha kisanduku kilichochaguliwa pekee, kwa namna ile ile ambayo `ccam doctor --fix` hurekebisha: kwenye kisanduku cha akaunti hurekebisha akaunti hiyo pekee, na kwenye kisanduku cha `General` matokeo yasiyo ya akaunti yoyote (njia za plugins na njia fupi); kila akaunti nyingine hukaguliwa tena tu. Huuliza kwanza (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), na kisanduku kisicho na chochote cha kurekebisha husema tu `nothing to repair in <name>`. `F` hurekebisha kila kitu kwa mara moja, kama `ccam doctor --fix`, na `r` hukagua tena.

## Kichupo cha Agents

Kichupo cha Agents huonyesha kile ambacho [`ccam agents`](/features/agents) huchapisha: vipindi vya chinichini na vipindi vilivyo wazi vya kila akaunti, pamoja na akaunti ambayo kila kimoja ni chake. Huuliza Claude Code ya kila akaunti kwa wakati mmoja (`claude agents --json`, kila moja ikiwa na muda wa mwisho wa sekunde 15) unapofungua kichupo kwa mara ya kwanza, kwa `r`, baada ya kusimamisha, na baada ya badiliko la config, na husema `Asking every account's Claude Code for its sessions…` hadi zote zijibu.

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

`State` ni hali ya kipindi cha chinichini na kile ambacho mchakato wake unafanya, au `open` kwa kipindi shirikishi, kama ilivyo katika `ccam agents`. Ikiwa akaunti haiwezi kuulizwa, mstari wa `✘` chini ya orodha huonyesha kile ambacho Claude Code ilisema, kama `lab` hapo juu, ambayo mwonekano wake wa agents (agent view) umezimwa. Vipindi vya chinichini vilivyokwisha huachwa nje (`ccam agents --all` huviorodhesha), na kila akaunti ikiwa imejibu na hakuna yenye kipindi, kichupo husema ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` Kwenye terminal nyembamba kichupo huondoa `Started`, kisha `Name`, kisha `Folder`, na hufupisha njia ndefu za folda kuanzia kushoto.

`Enter` kwenye kipindi cha chinichini hufunga dashibodi na hufungua kipindi hicho katika terminal hii kama akaunti yake, kama `ccam attach <id>` inavyofanya; dashibodi hairudi (ndani ya Claude Code, `←` huenda kwenye mwonekano wake wa agents na `Ctrl+Z` kwenye shell). Kipindi kilicho wazi hakina ID, kwa hiyo `Enter` na `x` husema tu kwamba kipindi cha chinichini kinahitajika. `x` huuliza kwanza, `No, keep it` ikiwa imechaguliwa: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` Ukijibu ndiyo, ccam huendesha `claude stop <id>` kama akaunti hiyo, husema ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` au kwa nini haikuweza, na husoma orodha tena.

## Kubadilisha vikundi

Kwenye kichupo cha Groups, `Enter` hufungua kihariri: vipengele upande wa kushoto, akaunti upande wa kulia, `Tab` kubadili, `Space` kuwasha au kuzima. Kipengele ambacho mojawapo ya akaunti tayari inakipata kutoka kikundi kingine hufifishwa pamoja na sababu, vivyo hivyo akaunti ambayo ingepata kipengele mara mbili; kipengele kinaweza kutoka kwenye kikundi kimoja tu. `Enter` huonyesha mpango:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Faili ikitofautiana na nakala ya kikundi, ccam kwanza huuliza la kufanya (kuhifadhi ya kikundi, kutumia yangu au kuunganisha). Kisha badiliko huendeshwa chinichini kwa msimbo ule ule wa [`ccam group`](/features/share): usanidi huhifadhiwa kwanza, data ya kila akaunti huhamishwa kwa mpangilio, na faili asili huenda `.ccam-backup/`; akaunti inayoacha kushiriki MCP huhifadhi seva za kikundi katika `.claude.json` yake.

Badiliko ambalo lingehamisha data ya akaunti inayoendeshwa hukataliwa pamoja na jina la akaunti hiyo, hata pale akaunti ilipoanzishwa wakati mpango ulikuwa wazi, na mpango uliopitwa na wakati kwa sababu mchakato mwingine wa ccam ulibadilisha vikundi hukataliwa kwa `the groups changed while this plan waited; nothing was moved, make the change again`. Wakati mchakato mwingine wa ccam unabadilisha vikundi, badiliko jipya, kubadilisha jina au kuondoa husema `another ccam process is changing the groups right now`; baada ya badiliko lililosimama nusu njia, husubiri hadi `F` kwenye kichupo cha Health (au `ccam doctor --fix`) ilikamilishe.

## Kuongeza akaunti

`a` hufungua mchawi ambao `ccam add` pia hutumia kwenye terminal:

1. Jina la akaunti.
2. Jinsi inavyoingia: usajili wa Claude sasa (kwenye kivinjari), [API key](/features/api), au usajili wa Claude baadaye.
3. Kwa API key: mahali ccam husoma ufunguo katika kila uanzishaji (`env:NAME` au `command:COMMAND`), endpoint, jinsi ufunguo unavyotumwa na model.
4. Vikundi vya kujiunga, `default` ikiwa imeteuliwa (ikiwa bado hakuna vikundi, kile ambacho kikundi cha `default` kitashiriki).
5. Amri ya njia fupi, ukiitaka.
6. Muhtasari wa kuthibitisha.

`Esc` hurudi hatua moja nyuma, na kwenye hatua ya kwanza hufunga mchawi.

## Kuondoa akaunti

`x` huuliza kabla ya chochote kutokea, No ikiwa imechaguliwa, na husema uondoaji hufanya nini: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` Kwa akaunti ya API key hakuna kuingia kunakohitaji kubatilishwa, na akaunti nyingine inaposhikilia kuingia kule kule, ccam haitoki kwenye akaunti hiyo. Akaunti iliyoingizwa huhifadhi folda yake: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## Rangi

Rangi hufuata majukumu ya kiolesura cha terminal cha Claude Code, ili ccam ionekane ikiendana nayo; majina ya vikundi hupata rangi kutokana na jina lake. ccam huuliza terminal kuhusu mandharinyuma yake na huchagua mandhari angavu au meusi. Weka moja mwenyewe juu kabisa ya `config.toml`, kabla ya `[[group]]`, `[[profile]]` au `[[folder]]` ya kwanza:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Mandhari ya `-daltonized` ni kwa ajili ya upofu wa rangi nyekundu na kijani, na mandhari ya `-ansi` hutumia rangi 16 za terminal yenyewe. `NO_COLOR` huzima rangi.

## Usalama na hali maalum

- Dashibodi huita kazi zile zile ambazo amri huita, kwa hiyo kila kanuni hubaki: mipango kabla data haijahamishwa, kukataa kwa akaunti zinazoendeshwa, nakala rudufu, na tokeni za kuingia ambazo kamwe haziondoki kwenye akaunti.
- Dashibodi haina jaribio (dry run): `ccam --dry-run` bila jina huchapisha jinsi akaunti ya folda hii ingeanza, na `o` kwenye akaunti hutoa jaribio la uanzishaji huo mmoja tu.
- Terminal hurejeshwa kabla Claude Code haijaanza au kuingia hakujafunguka.
- Skrini iliyo ndogo mno kwa kidirisha huondoa kwanza safu za kichupo; mistari hukatwa kwa `…` badala ya kufungwa kwenye mstari unaofuata. Kwenye skrini fupi, mistari mitupu na maelezo ya kichupo huondolewa baadaye, ili mstari wa vidokezo vya vitufe ubaki ukionekana. Chini ya safu wima 40 au mistari 8, dashibodi huonyesha tu `Too small: ccam needs 40x8 (now 36x7)`.
- `Ctrl+C` wakati badiliko, uondoaji, sasisho au marekebisho yanatekelezwa husubiri yakamilike: mstari wa hali husema `ccam quits when it is done.`
- Baada ya `U` kusasisha ccam, pendekezo hilo halirudi dashibodi inaporudi kutoka kuongeza akaunti au kuingia.

## Zinazohusiana

- [Kushiriki kwa vikundi](/features/share)
- [Kuanzisha Claude Code](/features/launch)
- [Kusimamia akaunti](/features/accounts)
- [Vikomo vya matumizi](/features/usage)
