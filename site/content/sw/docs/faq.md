---
title: Maswali yanayoulizwa mara kwa mara
metaTitle: "Maswali kuhusu ccam: akaunti kadhaa za Claude Code"
description: Majibu kuhusu ccam, akaunti za Claude Code, vikomo vya matumizi, kushiriki, faragha na mifumo inayotumika.
updated: 2026-11-19
order: 6
translation: machine
sourceHash: "090ecc3a34bd"
---

## Je, ccam ni zana rasmi ya Anthropic?

Hapana. ccam ni mradi usio rasmi wa open source na hauhusiani na Anthropic wala haujaidhinishwa nayo. Inatumia mpangilio wa `CLAUDE_CONFIG_DIR` ambao [Claude Code inauelezea katika nyaraka zake](https://code.claude.com/docs/en/env-vars) kwa ajili ya kuendesha akaunti kadhaa sambamba.

## Je, ninaruhusiwa kutumia akaunti kadhaa za Claude Code?

Hilo linategemea masharti ya kila akaunti na, kwa akaunti za kazi, sera ya shirika lako: [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) za Anthropic kwa Pro na Max, [Commercial Terms](https://www.anthropic.com/legal/commercial-terms) zake kwa Team na Enterprise, na [Usage Policy](https://www.anthropic.com/legal/aup) yake kwa zote. Yakague kabla ya kutumia akaunti kadhaa. ccam haiongezi, haiunganishi wala haifichi vikomo vya akaunti, na kila akaunti bado husimama kwenye vikomo vyake yenyewe. Huhamisha kazi yako kwenda akaunti nyingine kati ya akaunti zako unapoomba, au yenyewe kwenye kikomo unapowasha [`auto_switch`](/features/next#auto-switch) (Linux na macOS). Hakikisha kama hilo linaruhusiwa kwa akaunti zako kabla ya kulitumia; uamuzi ni wako.

**Viti (seats) vya Team na Enterprise.** Uchaguzi wa akaunti wa kiotomatiki (`ccam next`, `--fallback`, `ccam pick`, [`ccam agent`](/features/agent), kitufe `n` cha dashibodi na [`auto_switch`](/features/next#auto-switch) unapokiwasha) pia hubadili kati ya viti vya shirika lile lile la Claude Team au Enterprise. Kituo cha usaidizi cha Anthropic kinasema "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (kila kiti katika mpango wa Team hutolewa kwa mtu mmoja aliyetajwa kwa jina, na ni mtu huyo pekee anayeruhusiwa kukitumia) na "Users can't route their work through another user's account to use that user's remaining usage" (watumiaji hawaruhusiwi kuelekeza kazi yao kupitia akaunti ya mtumiaji mwingine ili kutumia kiasi cha matumizi kilichobaki cha mtumiaji huyo) ([mipango ya Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); vivyo hivyo kwa [mipango ya Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Kwa hiyo, kutumia uchaguzi wa kiotomatiki kati ya viti kunaweza kukiuka masharti ya Anthropic. Kama utazitumia huko ni uamuzi wako: wasiliana na shirika lako na Anthropic kabla ya kuzitegemea. Mipango ya binafsi (Pro, Max) inaelezwa katika aya iliyo hapo juu.

## Je, ccam hukwepa vikomo vya matumizi?

Si zaidi ya kikomo cha akaunti yoyote: kila akaunti huhifadhi vikomo vyake na husimama kwenye vikomo hivyo, na ccam haiviongezi, haiviunganishi wala haivifichi. ccam huvionyesha na hukuwezesha kuendelea kwenye akaunti nyingine kati ya akaunti zako pale moja inapoishiwa; `auto_switch` ikiwa imewashwa, hufanya hivyo yenyewe (tazama swali linalofuata). Kati ya viti vya shirika moja la Team au Enterprise, kuendelea kwenye kiti kingine hutumia matumizi yaliyobaki ya kiti hicho, jambo ambalo kituo cha usaidizi cha Anthropic hukataza kwa kiti cha mtumiaji mwingine (tazama Viti vya Team na Enterprise hapo juu).

## Je, ccam inaweza kubadili akaunti yenyewe kipindi kinapofikia kikomo chake?

Ni pale tu unapoiwasha. Kwa `auto_switch = true` juu kabisa ya `~/.config/ccam/config.toml` (Linux na macOS; imezimwa kwa chaguo-msingi na ni ya majaribio), Claude Code inaposimamisha kipindi shirikishi (kilichoanzishwa na ccam) kwenye kikomo cha saa 5, cha wiki, cha Opus, cha Sonnet au cha Fable, ccam humaliza kipindi na huendeleza mazungumzo kama akaunti nyingine kati ya akaunti zako yenye nafasi: kiti kingine cha shirika lile lile la Team au Enterprise, au mpango mwingine wa binafsi (Pro, Max). Kamwe haibadili kabla ya kikomo, kamwe haifuatilii matumizi yako chinichini, na kamwe haihamishi kazi kutoka kwenye shirika la Team au Enterprise wala kuingia ndani ya moja.

Kwenye viti huhamisha kazi yako kutoka kiti kimoja kwenda kingine, jambo ambalo linaweza kukiuka masharti ya Anthropic (tazama Viti vya Team na Enterprise hapo juu). Kwa mipango ya binafsi, masharti ya Anthropic hayalizungumzii jambo hili waziwazi: Consumer Terms hukataza kufikia Huduma "through automated or non-human means" (kwa njia za kiotomatiki au zisizo za binadamu) isipokuwa Anthropic iruhusu waziwazi, na uhamisho hutuma ujumbe mmoja kwa niaba yako. Kuiwasha ni uamuzi wako.

[`auto_switch`](/features/next#auto-switch) huorodhesha kila sharti na gharama ya uhamisho.

## Nitatumiaje akaunti mbili za Claude Code kwenye kompyuta moja?

Ipe kila akaunti folda yake ya usanidi kwa `CLAUDE_CONFIG_DIR`. ccam hukufanyia hivyo: endesha `ccam add` kwa kila akaunti, kisha uzianzishe kwa `ccam personal` na `ccam work`, au kwa njia fupi unazozipa majina mwenyewe, kama `ccwork`. [Mwongozo](/guides/multiple-claude-code-accounts) pia unaonyesha jinsi ya kufanya hivyo kwa mkono.

## Nitaendeleaje na mazungumzo kwenye akaunti nyingine?

Endesha `ccam next --continue` ndani ya folda: huhamisha mazungumzo ya mwisho ya folda hiyo kwenda akaunti inayofuata (`--auto` kwa ile yenye nafasi kubwa zaidi). Ndani ya Claude Code, `/ccam:switch` kutoka [plugin ya ccam](/docs/install#claude-code-plugin) hufanya vivyo hivyo kwa kipindi ulichomo. Akaunti hizo mbili zinaposhiriki historia ya mazungumzo (`ccam group set default --add history`, au kikundi chao wenyewe: `ccam group add pair --items history --members work,side`), akaunti mpya hufungua mazungumzo hayo hayo; vinginevyo Claude Code huyanakili kwenye historia ya akaunti mpya. Ongeza `--brief` ili kuanzisha akaunti mpya kwa muhtasari mfupi badala ya mazungumzo yote.

## Je, ccam huona au kutuma taarifa zangu za kuingia?

ccam kamwe haishiriki wala kunakili taarifa zako za kuingia za Claude kati ya akaunti, na kamwe hainakili kuingia kwa seva za MCP. Husoma tokeni ya ufikiaji ya akaunti ili tu kuuliza endpoint ya matumizi ya Anthropic kuhusu matumizi ya akaunti hiyo: kwa `ccam usage`, dashibodi, `ccam doctor` na chaguo za kiotomatiki (`ccam next --auto`, `ccam pick`, `--fallback`, na `auto_switch` kipindi kinapofikia kikomo chake). Hufanya upya tokeni iliyokwisha muda kwa `ccam usage --refresh` pekee; kwa `usage_source = "claude"` haitumi tokeni popote na husoma taarifa za kuingia zilizohifadhiwa (faili, au kipengele cha Keychain ya macOS) ili tu kujua kama bado zinafanya kazi. `ccam ext` na `ccam doctor` husoma majina ya kuingia kwa MCP ya akaunti, kamwe si tokeni zake. ccam haina telemetry. Mbali na maswali hayo, ombi pekee inalotuma yenyewe ni ukaguzi wa dashibodi kwenye GitHub wa toleo jipya zaidi la ccam, ambalo halibebi chochote kukuhusu wewe au akaunti zako.

## Nitaingiaje kwenye seva ya mbali?

Endesha `ccam add` au `ccam login <name>` kwenye kipindi cha SSH kama kawaida. Claude Code husubiri kuingia kwenye seva, kwa hiyo kivinjari kwenye kompyuta yako huishia kwenye ukurasa wa `http://localhost:<port>/callback?code=…` ambao haufunguki. Nakili anwani hiyo kutoka kwenye kivinjari na uibandike kwenye terminal: ccam huipitisha kwa Claude Code kwenye seva, ambayo hukamilisha kuingia. Anwani hiyo ina msimbo wa matumizi ya mara moja, kwa hiyo usiishiriki.

## Akaunti zangu zimehifadhiwa wapi?

Katika `~/.local/share/ccam/profiles/<name>/`, folda moja kwa kila akaunti. Kwenye macOS, Claude Code huhifadhi taarifa za kuingia zenyewe katika Keychain, chini ya jina linalotokana na folda hiyo, kwa hiyo huko `ccam rename` huhifadhi jina la zamani la folda.

## Nini kinatokea kwa ~/.claude yangu iliyopo?

Hakuna kitu, isipokuwa uendeshe `ccam import` au `ccam migrate`. Import husajili folda pale ilipo, kwa hiyo `claude` ya kawaida na `ccam <name>` huendelea kutumia akaunti na mipangilio ile ile. `ccam migrate` hunakili mipangilio yake, `CLAUDE.md`, viendelezi na seva za MCP kwenye kikundi cha `default` kinachoshirikiwa, na husoma folda tu.

## ccam inatumika kwenye mifumo gani?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine na mingineyo), macOS na Windows 10 na 11, moja kwa moja au ndani ya WSL. Kwenye Windows, kushiriki faili kama `settings.json` kunahitaji Developer Mode; folda zinazoshirikiwa hufanya kazi bila hiyo, na bila hiyo `ccam add` ya kwanza huunda kikundi cha `default` chenye seva za MCP na viendelezi pekee ([kushiriki kwenye Windows](/features/share#on-windows)).

## Kwa nini seva ya MCP inaonekana katika akaunti moja tu?

Seva zilizoongezwa ndani ya kipindi kwa `claude mcp add` huhifadhiwa kwa akaunti hiyo tu, na bila `-s user` kwa folda ya sasa tu. ccam hukujulisha wakati ujao utakapoanzisha akaunti hiyo, na `ccam mcp list` na `ccam ext` huzionyesha; hamishia seva kwenye kikundi kinachoshiriki kwa `ccam mcp import <account> <server>`.

## Je, ccam hupunguza kasi ya Claude Code?

Kuanzisha Claude Code kupitia ccam huongeza takriban milisekunde 4, zilizopimwa kwa `scripts/perf-check.sh` kwenye kompyuta mpakato. ccam husoma faili moja ndogo, hukagua viungo vinavyoshirikiwa kisha hujibadilisha yenyewe kuwa `claude`, kwa hiyo kwa chaguo-msingi haiendeshwi wakati unafanya kazi. Windows haiwezi kubadilisha mchakato, kwa hiyo huko ccam huanzisha Claude Code na kuisubiri, bila kufanya kitu kingine. `relay` au `auto_switch` ikiwa imewashwa, ccam huendelea kuendeshwa kwa njia hiyo hiyo kando ya kila kipindi shirikishi, ili iweze kuendeleza mazungumzo kwenye akaunti nyingine.

## Nitasasishaje ccam?

Endesha `ccam update`. Hupakua toleo jipya zaidi, hukagua checksum yake na hubadilisha faili ya `ccam`, huku `ccam update --check` ikikuambia tu kama kuna toleo jipya zaidi. Dashibodi pia hukuambia kunapokuwa na toleo jipya zaidi, na `U` hulisakinisha baada ya kuuliza; ccam kamwe haijisasishi yenyewe. Ikiwa ulisakinisha kifurushi cha `.deb`, badala yake sakinisha `.deb` mpya kwa apt. Tazama [kusasisha ccam](/docs/install#update).

## Nitaondoaje ccam?

Endesha `ccam uninstall`. Huorodhesha kitakachoondolewa na huuliza kwanza, kisha huondoa kila akaunti (baada ya kutoka kwanza kwenye akaunti hiyo), amri za njia fupi, `~/.config/ccam`, `~/.local/share/ccam`, `~/.local/state/ccam` na amri ya `ccam`. `~/.claude` iliyoingizwa hubaki, ikiwa na nakala ya kile ilichoshiriki katika vikundi vyake. Tazama [usakinishaji](/docs/install#uninstall).
