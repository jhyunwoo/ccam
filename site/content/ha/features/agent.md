---
title: Subagents a kan wasu asusai
description: ccam agent yana gudanar da ɓangaren aikin Claude Code a wani asusunka kuma yana dawo da amsar Claude, don haka dogon aiki yana kashe amfanin asusai da yawa.
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Subagents na Claude Code na kansa suna gudana a matsayin asusun zaman da ya fara su, don haka suna kashe amfanin wannan asusun. `ccam agent "<task>"` yana gudanar da aiki da `claude -p` a cikin foldar iri ɗaya a matsayin wani daga cikin asusanka kuma yana buga amsar Claude ta ƙarshe. Idan aka kira shi daga zaman Claude Code, yana aiki kamar subagent da ke kashe amfanin wani asusu dabam, don haka aiki ɗaya mai tsawo zai iya ci gaba bayan lokacin da asusun zaman na kansa zai ƙare.

**Kujeru (seats) na Team da Enterprise.** `ccam agent` yana zaɓar asusun da kansa. A cikin ƙungiyar Claude Team ko Enterprise ɗaya, yana gudanar da aikin a wata kujera dabam, wanda na iya saɓa wa sharuɗɗan Anthropic. Ka karanta [Tambayoyin da ake yawan yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) kuma ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi.

## Amfani

Claude Code yana gudanar da shi ta kayan aikinsa na Bash. Idan an shigar da [plugin ɗin ccam](/docs/install#claude-code-plugin), ka nema da kalmomi masu sauƙi, misali "yi amfani da sauran asusaina don sassan wannan aikin", ko ka rubuta haka a cikin CLAUDE.md ɗinka. Skill ɗin plugin ɗin yana gaya wa Claude Code ya fara kowane agent a bayan fage, ya rubuta aikin da zai tsaya da kansa, kuma ya ba agents da ke gyara fayiloli nasu ɓangare ko nasu git worktree.

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

Amsar Claude ta ƙarshe tana zuwa standard output; layukan da suka fara da `ccam:` suna zuwa standard error. `--json` yana buga `{account, session_id, result, is_error, accounts}` maimakon haka, inda `accounts` ke lissafa kowane asusun da aikin ya gudana a matsayinsa.

Exit codes: na claude na kansa (`0` idan agent ɗin ya gama), `75` (`EX_TEMPFAIL`) idan babu wani asusu banda mai kiran da ke da sarari, `1` don wasu kurakurai kamar agent da ke ƙoƙarin fara wani.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--model <model>` | Model ɗin agent ɗin; tagar mako ta wannan model ɗin ma tana ƙirguwa a cikin zaɓin | model na tsoho na asusun |
| `--resume <id>` | Ci gaba da tattaunawar wani agent na baya da aikin da ya biyo baya (session ID ɗin yana kan layinsa na farko da na ƙarshe) | a kashe |
| `--json` | Buga asusun, zaman, amsar da asusan da ya gudana a matsayinsu | a kashe |
| `--max-usage <1-100>` | Fifita asusan da ke ƙasa da wannan amfani | `auto_max_usage`, in ba haka ba 100 |
| `--strategy headroom\|expiring\|balance` | [Ƙa'idar zaɓi](/features/next#auto) | `auto_strategy`, in ba haka ba `headroom` |
| `--cross-org` | Gudana a matsayin asusan wasu ƙungiyoyi ma | `cross_org`, in ba haka ba a kashe |
| `--wait[=<duration>]` | Idan babu asusun da ke da sarari, jira na farko ya sake saitawa, har zuwa iyakar wannan tsawon | a kashe; `5h` idan an ba shi shi kaɗai |
| `--dry-run` | Buga umarnin don asusun da zai zaɓa, ba tare da gudanar da shi ba | a kashe |

## Yadda yake aiki

1. **Ana barin mai kiran a waje.** Asusun da aka kira `ccam agent` daga gare shi (na zaman da ccam ya fara, in ba haka ba wanda wannan foldar ke farawa) ba ya taɓa gudanar da agent ɗin, haka ma wata foldar da ta shiga asusun Claude iri ɗaya, wadda ke raba iyakokinsa. Agent ɗin yana zama a cikin [ƙungiyar](/features/next#auto) mai kiran sai dai idan `--cross-org` ya faɗi wani abu dabam.
2. **Yana zaɓa kamar `ccam pick`.** Ana neman amfani kamar yadda [`ccam pick`](/features/pick) ke yi, kuma agents da aka fara a cikin minti 2 da juna suna samun asusai daban-daban muddin isassu suna da sarari.
3. **Yana gudana da izinin zaman, ba fiye da haka ba.** Babu wanda ke wurin don ba da izinin kayan aiki a `claude -p`, don haka agent ɗin yana samun yanayin izini (permission mode) da matakin ƙoƙari (effort) na zaman da ya kira shi. Idan an gudanar da shi daga terminal, yanayin tsoho na Claude Code da ƙa'idodin izini na saitunanka ne ke aiki. `ccam agent` ba shi da zaɓin da ke ba da fiye da haka.
4. **A iyakar amfani, tattaunawar tana ci gaba.** Idan asusun agent ɗin ya kai iyakarsa, ccam yana rubuta iyakar kamar yadda [hook ɗin plugin](/features/usage#the-limit-record) zai yi kuma yana ci gaba da tattaunawar iri ɗaya a matsayin asusun da ya fi sarari: Claude Code yana kwafar ta zuwa can (`--resume <file> --fork-session`), tare da saƙon farko "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." Ana gwada kowane asusu sau ɗaya.
5. **Ana rubuta kowace tattaunawar agent** a `~/.local/state/ccam/agent/<session id>.json` (kwanaki 30). `--resume <id>` yana ci gaba da ita a foldarta, a matsayin asusun da ya gudanar da ita idan har yanzu yana da sarari, in ba haka ba a matsayin wani. [`ccam next --continue`](/features/next) yana barin tattaunawar agents a waje, don haka yana ci gaba da naka.

## Tsaro da yanayi na musamman

- Agent ba zai iya fara agents nasa ba: ccam yana yi masa alama da `CCAM_AGENT=1`, kuma `ccam agent` yana ƙin gudana a can.
- Dakatar da `ccam agent` (Ctrl+C, sigina, ko Claude Code da ya katse umarnin) yana dakatar da agent ɗin ma.
- Agent ɗin ba ya ganin tattaunawar da ta kira shi, sai aikinsa kawai.
- Agents da ke gyara fayiloli iri ɗaya a lokaci ɗaya suna goge aikin juna. Ka ba kowannensu nasa ɓangare na aikin, ko nasa git worktree.
- Idan zaman mai kiran ya koma wani asusu ([`auto_switch`](/features/next#auto-switch)), agents da ya fara suna tsayawa tare da shi. Tattaunawarsu tana nan: `--resume <id>` yana ci gaba da su.
- `agent` sunan umarni ne, don haka ba za a iya sa wa wani sabon asusu suna `agent` ba; wanda ke akwai tun da farko har yanzu yana farawa da `ccam run agent`.

## Masu alaƙa

- [Zaɓar asusu don rubutun umarni](/features/pick)
- [Sauyawa zuwa asusu na gaba](/features/next)
- [Iyakokin amfani](/features/usage)
- [Zaman bayan fage na kowane asusu](/features/agents)
