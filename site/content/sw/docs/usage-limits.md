---
title: Vikomo vya matumizi na kubadili akaunti
description: Ona matumizi ya saa 5 na ya wiki ya kila akaunti ya Claude Code kwa ccam usage, na ufungue akaunti yenye nafasi kubwa zaidi kwa ccam next --auto.
updated: 2026-11-05
order: 5
translation: machine
sourceHash: "06bd349073e3"
---

`ccam usage` huonyesha, kwa kila akaunti, kiasi kilichotumika cha kila kikomo: kikomo cha saa 5, kikomo cha wiki na kikomo cha wiki cha kila model (kama Fable). `ccam next --auto` hufungua akaunti yenye nafasi kubwa zaidi, na `--continue` huhamisha mazungumzo ya sasa.

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

Kila akaunti hupata mstari unaoonyesha umri wa namba zake, kisha upau mmoja kwa kila kikomo pamoja na muda hadi kianze upya: kipindi cha saa 5, kipindi cha wiki kwa model zote, na kipindi cha wiki cha kila model yenyewe, kama Fable. Upau hugeuka manjano kwenye 80% na nyekundu kwenye 100%. Mstari wa mwisho hutaja akaunti ambayo `ccam next --auto` ingefungua; kwa `auto_strategy = "expiring"` au `"balance"` badala yake husomeka `→ next --auto picks <name>: <reason>`.

**Viti (seats) vya Team na Enterprise.** Uchaguzi wa akaunti wa kiotomatiki (`ccam next`, `--fallback`, `ccam pick`, kitufe `n` cha dashibodi na, ukishakiwasha, [`auto_switch`](/features/next#auto-switch)) pia hubadili kati ya viti vya shirika lile lile la Claude Team au Enterprise, jambo ambalo linaweza kukiuka masharti ya Anthropic. Soma [Maswali yanayoulizwa mara kwa mara](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) na uwasiliane na shirika lako na Anthropic kabla ya kulitegemea.

## Namba zinatoka wapi

ccam huuliza endpoint ya matumizi ya Anthropic kwa ajili ya kila akaunti kwa wakati mmoja, ikitumia taarifa za kuingia za kila akaunti yenyewe. Husubiri sekunde 3 kwa akaunti kwa kiwango cha juu na huonyesha kila akaunti mara tu inapojibu.

Endpoint hii haijaandikiwa nyaraka hadharani na inaweza kubadilika. Ombi linaposhindwa, ccam huonyesha namba za mwisho zinazojulikana na kuziwekea alama kuwa ni za cache. Kipindi cha muda kilichoanza upya tangu namba hizo zihifadhiwe huonyesha 0%. `ccam doctor` hukujulisha ikiwa muundo wa jibu umebadilika.

Claude Code inaposimamisha kipindi kwenye kikomo cha matumizi, [plugin ya ccam](/docs/install#claude-code-plugin) hurekodi ni kipindi kipi cha muda kilichojaa na lini kinaanza upya. Hadi wakati huo kipindi hicho huonyesha 100% bila kujali namba nyingine zinasema nini, kwa hiyo `ccam next --auto` na `--fallback` hazichagui akaunti hiyo tena mapema mno.

Namba zilizopatikana chini ya sekunde 60 zilizopita hutumiwa tena. Chaguo:

| Chaguo | Athari |
| --- | --- |
| `--watch[=seconds]` | Onyesha upya mfululizo, kila sekunde 60 kwa chaguo-msingi |
| `--json` | Matokeo yanayosomeka na mashine |
| `-v` | Onyesha pia namba za kila akaunti zilikotoka |
| `--no-cache` | Uliza upya kila mara |
| `--source cache` | Usiulize kamwe; onyesha namba zilizohifadhiwa tu |

## Taarifa za kuingia zilizoisha muda

Tokeni ya ufikiaji ya akaunti inapoisha muda, ccam haiifanyi upya, kwa sababu kuifanya upya wakati kipindi cha Claude Code kinatumia taarifa zile zile za kuingia kunaweza kusababisha kipindi hicho kitolewe nje ya akaunti (sign out). Akaunti huonyesha ``access token expired; `ccam <name>` renews it when it starts``: ianzishe mara moja kwa `ccam <name>` na Claude Code itaifanya upya. Taarifa za kuingia zilizoisha muda kabisa, au ambazo Claude Code ilizifuta, huonyesha `` login expired → `ccam login <name>` `` na huhitaji kuingia upya.

`ccam usage --refresh` hufanya upya tokeni zilizoisha muda yenyewe, lakini kwa akaunti zisizo na kipindi cha Claude Code kinachoendelea tu, bila kujali kipindi hicho kilianzishwa vipi. Kamwe haifanyi upya taarifa za kuingia zilizohifadhiwa katika Keychain ya macOS au katika `~/.claude` iliyoingizwa, na kufanya upya kunaweza kuchukua hadi sekunde 30.

## Matumizi kwenye mstari wa hali

```sh
ccam statusline install
```

Kisha Claude Code huonyesha jina la akaunti na matumizi yake kwenye mstari wa hali (status line), kwa mfano `[work] · Opus · 5h 24% · 7d 41%`. Kuanzia 90%, mstari pia hupendekeza `ccam next --auto --continue`. Ikiwa tayari una amri ya mstari wa hali, ccam huihifadhi na huweka jina la akaunti mbele yake. Mstari wa hali pia huhifadhi namba zake, ambazo `ccam usage` hutumia wakati endpoint haijibu.

## Kubadili

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the most room
ccam next --continue  # continue this folder's conversation on the next account
```

`ccam next` huruka akaunti ambazo hazijaingia, na `--auto` huruka akaunti zilizofikia 100% (ikiwemo kipindi ambacho Claude Code iliripoti kuwa kimejaa, hadi kianze upya) na akaunti za shirika jingine. Akaunti hizo mbili zikishiriki historia ya mazungumzo, `--continue` hufungua mazungumzo yale yale; vinginevyo Claude Code huyanakili kwenye akaunti mpya. Tazama [kushiriki](/docs/sharing) na [kubadili kwenda akaunti inayofuata](/features/next#continue).

Kwa `auto_switch = true` juu kabisa ya `config.toml` (Linux na macOS; imezimwa kwa chaguo-msingi), ccam hufanya ubadilishaji huu yenyewe Claude Code inaposimamisha kipindi shirikishi kwenye kikomo chake na akaunti nyingine ya aina hiyo hiyo (kiti katika shirika lile lile, au mpango mwingine wa binafsi) ina nafasi. Tazama [`auto_switch`](/features/next#auto-switch) ili kujua lini hutenda na kile ambacho kamwe haifanyi.

Kwa script na kazi zilizopangwa zinazoendesha `claude -p`, weka `--fallback` kabla ya jina la akaunti. Akaunti hiyo inapofikia kikomo chake, ccam huendesha amri ile ile tena kama akaunti yenye nafasi kubwa zaidi, na matokeo ya uendeshaji uliokamilika pekee ndiyo huchapishwa:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Tazama [Endesha tena kwenye kikomo](/features/next#run-again-at-the-limit) kuona jinsi ccam inavyoamua.

[Mwongozo wa hatua za kufuata kikomo cha matumizi kinapofikiwa](/guides/usage-limit-workflow) unapitia mfano kamili.
