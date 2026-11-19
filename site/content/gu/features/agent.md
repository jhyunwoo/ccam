---
title: બીજાં એકાઉન્ટ પર subagents
description: ccam agent, Claude Code ના task નો એક ભાગ તમારા બીજા એકાઉન્ટ તરીકે ચલાવી Claude નો જવાબ આપે છે, તેથી એક લાંબું task ઘણાં એકાઉન્ટનો વપરાશ વાપરે છે.
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Claude Code ના પોતાના subagents તેમને શરૂ કરનાર સેશનના એકાઉન્ટ તરીકે ચાલે છે, તેથી તેઓ એ એકાઉન્ટનો વપરાશ ખર્ચે છે. `ccam agent "<task>"` એ જ ફોલ્ડરમાં `claude -p` વડે એક task તમારાં બીજાં એકાઉન્ટમાંથી એક તરીકે ચલાવે છે અને Claude નો છેલ્લો જવાબ છાપે છે. Claude Code ના સેશનમાંથી બોલાવતાં, તે બીજા એકાઉન્ટનો વપરાશ ખર્ચતા subagent ની જેમ કામ કરે છે, તેથી સેશનનું પોતાનું એકાઉન્ટ ખૂટી ગયું હોત ત્યાર પછી પણ એક લાંબું task આગળ ચાલી શકે છે.

**Team અને Enterprise સીટ.** `ccam agent` એકાઉન્ટ જાતે પસંદ કરે છે. એક જ Claude Team કે Enterprise સંસ્થાની અંદર તે task બીજી સીટ પર ચલાવે છે, જેનાથી Anthropic ની શરતોનું ઉલ્લંઘન થઈ શકે છે. તેના પર આધાર રાખતાં પહેલાં [વારંવાર પૂછાતા પ્રશ્નો](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) વાંચો અને તમારી સંસ્થા અને Anthropic પાસેથી ખાતરી કરી લો.

## ઉપયોગ

Claude Code તેને પોતાના Bash tool વડે ચલાવે છે. [ccam plugin](/docs/install#claude-code-plugin) ઇન્સ્ટૉલ હોય ત્યારે, તેને સાદા શબ્દોમાં માગો, ઉદાહરણ તરીકે "આ task ના ભાગો માટે મારાં બીજાં એકાઉન્ટ વાપરો", અથવા એ વાત તમારી CLAUDE.md માં લખો. plugin નું skill Claude Code ને કહે છે કે દરેક agent ને બૅકગ્રાઉન્ડમાં શરૂ કરે, પોતાની રીતે પૂરું સમજાય તેવું task લખે, અને ફાઇલો બદલતા agents ને તેમનો પોતાનો ભાગ કે પોતાનું git worktree આપે.

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

Claude નો છેલ્લો જવાબ standard output પર જાય છે; `ccam:` થી શરૂ થતી લીટીઓ standard error પર જાય છે. `--json` તેના બદલે `{account, session_id, result, is_error, accounts}` છાપે છે, જ્યાં `accounts` એ બધાં એકાઉન્ટની યાદી છે જેના તરીકે task ચાલ્યું.

Exit code: claude ના પોતાના (agent પૂરો થાય ત્યારે `0`), બોલાવનાર સિવાય કોઈ એકાઉન્ટમાં જગ્યા ન હોય ત્યારે `75` (`EX_TEMPFAIL`), અને બીજી ભૂલો માટે `1`, જેમ કે કોઈ agent બીજો agent શરૂ કરવાનો પ્રયાસ કરે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--model <model>` | agent નું મૉડલ; પસંદગીમાં એ મૉડલની સાપ્તાહિક વિન્ડો પણ ગણાય છે | એકાઉન્ટનું ડિફૉલ્ટ મૉડલ |
| `--resume <id>` | અગાઉના agent ની વાતચીત એક આગળના task સાથે ચાલુ રાખે છે (session ID તેની પહેલી અને છેલ્લી લીટી પર હોય છે) | off |
| `--json` | એકાઉન્ટ, સેશન, જવાબ અને જે એકાઉન્ટ તરીકે તે ચાલ્યું તે છાપે છે | off |
| `--max-usage <1-100>` | આ વપરાશથી નીચેનાં એકાઉન્ટને પ્રાથમિકતા આપે છે | `auto_max_usage`, નહીં તો 100 |
| `--strategy headroom\|expiring\|balance` | [પસંદગીનો નિયમ](/features/next#auto) | `auto_strategy`, નહીં તો `headroom` |
| `--cross-org` | બીજી સંસ્થાનાં એકાઉન્ટ તરીકે પણ ચલાવે છે | `cross_org`, નહીં તો off |
| `--wait[=<duration>]` | કોઈ એકાઉન્ટમાં જગ્યા ન હોય ત્યારે, પહેલું એકાઉન્ટ રીસેટ થાય તેની રાહ જુએ છે, વધુમાં વધુ આટલો સમય | off; એકલું આપો તો `5h` |
| `--dry-run` | જે એકાઉન્ટ પસંદ કરત તેના માટેનો કમાન્ડ છાપે છે, તેને ચલાવ્યા વગર | off |

## તે કેવી રીતે કામ કરે છે

1. **બોલાવનાર બહાર રહે છે.** જે એકાઉન્ટમાંથી `ccam agent` બોલાવાય છે (ccam એ શરૂ કરેલા સેશનનું એકાઉન્ટ, નહીં તો આ ફોલ્ડર જે શરૂ કરે તે એકાઉન્ટ) તે ક્યારેય agent ચલાવતું નથી, અને એ જ Claude એકાઉન્ટમાં લૉગિન થયેલું બીજું ફોલ્ડર પણ નહીં, કારણ કે તે તેની મર્યાદાઓ વહેંચે છે. `--cross-org` બીજું કંઈ ન કહે ત્યાં સુધી agent બોલાવનારની [સંસ્થા](/features/next#auto) માં જ રહે છે.
2. **તે `ccam pick` ની જેમ પસંદ કરે છે.** વપરાશ [`ccam pick`](/features/pick) ની જેમ જ જોવાય છે, અને એકબીજાથી 2 મિનિટની અંદર શરૂ થયેલા agents ને, જ્યાં સુધી પૂરતાં એકાઉન્ટમાં જગ્યા હોય ત્યાં સુધી, જુદાં એકાઉન્ટ મળે છે.
3. **તે સેશનની પરવાનગીઓ સાથે ચાલે છે, ક્યારેય વધુ નહીં.** `claude -p` માં કોઈ tool ને મંજૂરી આપવા માટે કોઈ હાજર હોતું નથી, તેથી agent ને તેને બોલાવનાર સેશનનો permission mode અને effort મળે છે. terminal માંથી ચલાવતાં, Claude Code નો ડિફૉલ્ટ mode અને તમારાં settings ના permission નિયમો લાગુ પડે છે. `ccam agent` પાસે વધુ પરવાનગી આપતો કોઈ ઓપ્શન નથી.
4. **વપરાશ મર્યાદાએ, વાતચીત આગળ વધે છે.** agent નું એકાઉન્ટ તેની મર્યાદાએ પહોંચે ત્યારે, ccam [plugin hook](/features/usage#the-limit-record) કરત તેમ મર્યાદાની નોંધ કરે છે અને એ જ વાતચીત સૌથી વધુ જગ્યાવાળા એકાઉન્ટ તરીકે ચાલુ રાખે છે: Claude Code તેની નકલ ત્યાં કરે છે (`--resume <file> --fork-session`), આ પહેલા સંદેશ સાથે: "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." દરેક એકાઉન્ટ એક વાર અજમાવાય છે.
5. **દરેક agent વાતચીત નોંધાય છે**, `~/.local/state/ccam/agent/<session id>.json` માં (30 દિવસ). `--resume <id>` તેને તેના ફોલ્ડરમાં ચાલુ રાખે છે: જે એકાઉન્ટે તે ચલાવી હતી તેમાં હજી જગ્યા હોય તો તે એકાઉન્ટ તરીકે, નહીં તો બીજા તરીકે. [`ccam next --continue`](/features/next) agent વાતચીતોને બહાર રાખે છે, તેથી તે તમારી વાતચીત જ ચાલુ રાખતું રહે છે.

## સલામતી અને ખાસ કિસ્સા

- કોઈ agent પોતાના agents શરૂ કરી શકતો નથી: ccam તેને `CCAM_AGENT=1` વડે ચિહ્નિત કરે છે, અને ત્યાં `ccam agent` ચાલવાનો ઇનકાર કરે છે.
- `ccam agent` ને રોકવાથી (Ctrl+C, કોઈ signal, કે Claude Code કમાન્ડને અટકાવે) agent પણ અટકે છે.
- agent તેને બોલાવનાર વાતચીત જોતો નથી, ફક્ત પોતાનું task જુએ છે.
- એકસાથે એ જ ફાઇલો બદલતા agents એકબીજાનું કામ ભૂંસી નાખે છે. દરેકને task નો પોતાનો ભાગ, કે પોતાનું git worktree આપો.
- બોલાવનાર સેશન બીજા એકાઉન્ટ પર જાય ત્યારે ([`auto_switch`](/features/next#auto-switch)), તેણે શરૂ કરેલા agents તેની સાથે અટકે છે. તેમની વાતચીતો રહે છે: `--resume <id>` તેમને ચાલુ રાખે છે.
- `agent` એક કમાન્ડનું નામ છે, તેથી કોઈ નવા એકાઉન્ટનું નામ `agent` રાખી શકાતું નથી; પહેલેથી હોય તેવું એકાઉન્ટ હજી પણ `ccam run agent` થી શરૂ થાય છે.

## સંબંધિત

- [scripts માટે એકાઉન્ટ પસંદ કરવું](/features/pick)
- [આગલા એકાઉન્ટ પર સ્વિચ કરવું](/features/next)
- [વપરાશ મર્યાદા](/features/usage)
- [દરેક એકાઉન્ટનાં બૅકગ્રાઉન્ડ સેશન](/features/agents)
