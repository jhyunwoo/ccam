---
title: Subagents kwenye akaunti nyingine
description: ccam agent huendesha sehemu ya kazi ya Claude Code kama akaunti yako nyingine na hurudisha jibu la Claude, ili kazi ndefu itumie matumizi ya akaunti kadhaa.
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Subagents wa Claude Code wenyewe huendeshwa kama akaunti ya kipindi kinachowaanzisha, kwa hiyo hutumia matumizi ya akaunti hiyo. `ccam agent "<task>"` huendesha kazi kwa `claude -p` katika folda ile ile kama akaunti yako nyingine na huchapisha jibu la mwisho la Claude. Ikiitwa kutoka kipindi cha Claude Code, hufanya kazi kama subagent anayetumia matumizi ya akaunti nyingine, kwa hiyo kazi moja ndefu inaweza kuendelea hata pale ambapo akaunti ya kipindi chenyewe ingekuwa imeishiwa.

**Viti (seats) vya Team na Enterprise.** `ccam agent` huchagua akaunti yenyewe. Ndani ya shirika moja la Claude Team au Enterprise huendesha kazi kwenye kiti kingine, jambo ambalo linaweza kukiuka masharti ya Anthropic. Soma [Maswali yanayoulizwa mara kwa mara](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) na uwasiliane na shirika lako na Anthropic kabla ya kulitegemea.

## Matumizi

Claude Code huiendesha kupitia zana yake ya Bash. Ukiwa umesakinisha [plugin ya ccam](/docs/install#claude-code-plugin), iombe kwa maneno ya kawaida, kwa mfano "tumia akaunti zangu nyingine kwa sehemu za kazi hii", au weka hilo katika CLAUDE.md yako. Skill ya plugin huiambia Claude Code ianzishe kila agent chinichini, iandike kazi inayojitosheleza, na iwape agents wanaohariri faili sehemu yao wenyewe au git worktree yao wenyewe.

```sh
ccam agent "find why TestParse fails and fix it; reply with what you changed"
ccam agent --model opus -- "-v breaks the parser: find out why"
ccam agent --json < task.md            # the task from standard input
ccam agent --resume 7f3c9a2e-… "now add a test for it"
ccam agent --wait=1h "…"               # nobody has room: wait for a reset
ccam --dry-run agent "…"               # the account and command it would run
```

```console
$ ccam agent "summarize what internal/usage does"
ccam: agent runs as work (12% of its 5 hours used); session 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f
internal/usage looks up the 5-hour and weekly usage of every account …
ccam: agent finished as work; follow up with: ccam agent --resume 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f "…"
```

Jibu la mwisho la Claude huenda kwenye standard output; mistari inayoanza na `ccam:` huenda kwenye standard error. `--json` huchapisha `{account, session_id, result, is_error, accounts}` badala yake, ambapo `accounts` huorodhesha akaunti zote ambazo kazi iliendeshwa nazo.

Exit codes: ya claude yenyewe (`0` agent alipomaliza), `75` (`EX_TEMPFAIL`) wakati hakuna akaunti yenye nafasi isipokuwa ile inayoita, `1` kwa makosa mengine kama agent anayejaribu kuanzisha agent mwingine.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--model <model>` | Modeli ya agent; kipindi cha wiki cha modeli hiyo pia huhesabiwa katika uchaguzi | modeli chaguo-msingi ya akaunti |
| `--resume <id>` | Huendeleza mazungumzo ya agent wa awali kwa kazi ya ufuatiliaji (session ID iliyo kwenye mstari wake wa kwanza na wa mwisho) | imezimwa |
| `--json` | Huchapisha akaunti, kipindi, jibu na akaunti ambazo iliendeshwa nazo | imezimwa |
| `--max-usage <1-100>` | Hupendelea akaunti zilizo chini ya matumizi haya | `auto_max_usage`, vinginevyo 100 |
| `--strategy headroom\|expiring\|balance` | [Kanuni ya uchaguzi](/features/next#auto) | `auto_strategy`, vinginevyo `headroom` |
| `--cross-org` | Huendesha pia kama akaunti za mashirika mengine | `cross_org`, vinginevyo imezimwa |
| `--wait[=<duration>]` | Hakuna akaunti yenye nafasi, husubiri ya kwanza ianze upya, si zaidi ya muda huu | imezimwa; `5h` ikitolewa peke yake |
| `--dry-run` | Huchapisha amri kwa akaunti ambayo ingechagua, bila kuiendesha | imezimwa |

## Jinsi inavyofanya kazi

1. **Akaunti inayoita huachwa nje.** Akaunti ambayo `ccam agent` inaitwa kutoka kwayo (ile ya kipindi ambacho ccam ilianzisha, vinginevyo ile ambayo folda hii huanzisha) kamwe haiendeshi agent, wala folda nyingine iliyoingia kwenye akaunti ile ile ya Claude, ambayo hushiriki vikomo vyake. Agent hubaki katika [shirika](/features/next#auto) la akaunti inayoita isipokuwa `--cross-org` iseme vinginevyo.
2. **Huchagua kama `ccam pick`.** Matumizi huombwa kwa njia ile ile kama [`ccam pick`](/features/pick) inavyofanya, na agents walioanzishwa ndani ya dakika 2 kati yao hupata akaunti tofauti maadamu akaunti za kutosha zina nafasi.
3. **Huendeshwa kwa ruhusa za kipindi, kamwe si zaidi.** Hakuna mtu wa kuruhusu zana katika `claude -p`, kwa hiyo agent hupata hali ya ruhusa na kiwango cha juhudi (effort) cha kipindi kilichomwita. Ikiendeshwa kutoka terminal, hali chaguo-msingi ya Claude Code na kanuni za ruhusa za mipangilio yako ndizo hutumika. `ccam agent` haina chaguo linalotoa zaidi.
4. **Kwenye kikomo cha matumizi, mazungumzo huendelea.** Akaunti ya agent inapofikia kikomo chake, ccam hurekodi kikomo kama ambavyo [hook ya plugin](/features/usage#the-limit-record) ingefanya na huendeleza mazungumzo yale yale kama akaunti yenye nafasi kubwa zaidi: Claude Code huyanakili huko (`--resume <file> --fork-session`), pamoja na ujumbe wa kwanza "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." Kila akaunti hujaribiwa mara moja.
5. **Kila mazungumzo ya agent hurekodiwa** katika `~/.local/state/ccam/agent/<session id>.json` (siku 30). `--resume <id>` huyaendeleza katika folda yake, kama akaunti iliyoyaendesha ikiwa bado ina nafasi, vinginevyo kama nyingine. [`ccam next --continue`](/features/next) huacha nje mazungumzo ya agents, kwa hiyo huendelea kuendeleza yako.

## Usalama na hali maalum

- Agent hawezi kuanzisha agents wake mwenyewe: ccam humwekea alama `CCAM_AGENT=1`, na `ccam agent` hukataa kuendeshwa hapo.
- Kusimamisha `ccam agent` (Ctrl+C, ishara, au Claude Code kukatiza amri) humsimamisha agent pia.
- Agent haoni mazungumzo yaliyomwita, huona kazi yake tu.
- Agents wanaohariri faili zile zile kwa wakati mmoja huandika juu ya kazi ya wenzao. Mpe kila mmoja sehemu yake ya kazi, au git worktree yake mwenyewe.
- Kipindi kinachoita kinapohamia akaunti nyingine ([`auto_switch`](/features/next#auto-switch)), agents kilichowaanzisha husimama pamoja nacho. Mazungumzo yao hubaki: `--resume <id>` huyaendeleza.
- `agent` ni jina la amri, kwa hiyo hakuna akaunti mpya inayoweza kuitwa `agent`; iliyopo bado huanza kwa `ccam run agent`.

## Zinazohusiana

- [Kuchagua akaunti kwa ajili ya script](/features/pick)
- [Kubadili kwenda akaunti inayofuata](/features/next)
- [Vikomo vya matumizi](/features/usage)
- [Vipindi vya chinichini vya kila akaunti](/features/agents)
