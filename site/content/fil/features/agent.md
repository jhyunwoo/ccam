---
title: Mga subagent sa ibang account
description: Pinapatakbo ng ccam agent ang bahagi ng Claude Code task sa ibang account mo at ibinabalik ang sagot ni Claude, para magamit ng mahabang task ang ilang account.
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Tumatakbo ang sariling mga subagent ng Claude Code bilang ang account ng session na nagsimula sa kanila, kaya ang usage ng account na iyon ang nagagamit nila. Pinapatakbo ng `ccam agent "<task>"` ang isang task gamit ang `claude -p` sa parehong folder bilang isa pa sa iyong mga account at ipi-print ang huling sagot ni Claude. Kapag tinawag mula sa isang session ng Claude Code, gumagana ito na parang subagent na gumagamit ng usage ng ibang account, kaya maipagpapatuloy ang isang mahabang task kahit naubos na sana ang sariling account ng session.

**Mga seat ng Team at Enterprise.** Kusang pumipili ng account ang `ccam agent`. Sa loob ng iisang Claude Team o Enterprise na organisasyon, pinapatakbo nito ang task sa ibang seat, na maaaring lumabag sa terms ng Anthropic. Basahin ang [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) at magtanong muna sa iyong organisasyon at sa Anthropic bago ka umasa rito.

## Paggamit

Pinapatakbo ito ng Claude Code sa pamamagitan ng Bash tool nito. Kapag naka-install ang [ccam plugin](/docs/install#claude-code-plugin), hilingin lang ito sa simpleng salita, halimbawa "gamitin ang iba kong account para sa mga bahagi ng task na ito", o ilagay iyon sa iyong CLAUDE.md. Sinasabi ng skill ng plugin sa Claude Code na simulan ang bawat agent sa background, na sumulat ng task na buo na sa sarili nito, at na bigyan ang mga agent na nag-e-edit ng mga file ng sarili nilang bahagi o sariling git worktree.

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

Napupunta sa standard output ang huling sagot ni Claude; napupunta sa standard error ang mga linyang nagsisimula sa `ccam:`. Sa halip, ipi-print ng `--json` ang `{account, session_id, result, is_error, accounts}`, kung saan inililista ng `accounts` ang bawat account na pinagpatakbuhan ng task.

Mga exit code: ang sa claude mismo (`0` kapag natapos ang agent), `75` (`EX_TEMPFAIL`) kapag walang account maliban sa tumawag ang may natitira, `1` para sa ibang error gaya ng agent na sumusubok magsimula ng isa pa.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--model <model>` | Ang model ng agent; kasama rin sa pagpili ang lingguhang window ng model na iyon | ang default na model ng account |
| `--resume <id>` | Ituloy ang usapan ng naunang agent gamit ang isang follow-up na task (ang session id na nasa una at huling linya nito) | off |
| `--json` | I-print ang account, session, sagot at ang mga account na pinagpatakbuhan nito | off |
| `--max-usage <1-100>` | Unahin ang mga account na mas mababa sa usage na ito | `auto_max_usage`, kung wala ay 100 |
| `--strategy headroom\|expiring\|balance` | Ang [patakaran ng pagpili](/features/next#auto) | `auto_strategy`, kung wala ay `headroom` |
| `--cross-org` | Patakbuhin din bilang mga account ng ibang organisasyon | `cross_org`, kung wala ay off |
| `--wait[=<duration>]` | Kapag walang account na may natitira, hintayin ang unang mag-reset, nang hindi hihigit sa tagal na ito | off; `5h` kapag ibinigay nang mag-isa |
| `--dry-run` | I-print ang command para sa account na pipiliin nito, nang hindi ito pinapatakbo | off |

## Paano ito gumagana

1. **Hindi isinasama ang tumawag.** Hindi kailanman pinapatakbo ang agent ng account kung saan tinawag ang `ccam agent` (ang sa session na sinimulan ng ccam, kung wala ay ang sinisimulan ng folder na ito), gayundin ng ibang folder na naka-log in sa parehong Claude account, na kahati nito sa mga limit. Nananatili ang agent sa [organisasyon](/features/next#auto) ng tumawag maliban kung iba ang sabihin ng `--cross-org`.
2. **Pumipili ito gaya ng `ccam pick`.** Tinitingnan ang usage sa paraang ginagawa ng [`ccam pick`](/features/pick), at magkakaibang account ang nakukuha ng mga agent na sinimulan nang hindi lalampas sa 2 minuto ang pagitan basta may sapat na account na may natitira.
3. **Tumatakbo ito gamit ang mga permission ng session, hindi kailanman higit pa.** Walang taong mag-a-allow ng tool sa `claude -p`, kaya nakukuha ng agent ang permission mode at effort level ng session na tumawag dito. Kapag pinatakbo mula sa terminal, ang default mode ng Claude Code at ang mga permission rule ng iyong settings ang umiiral. Walang option ang `ccam agent` na nagbibigay ng higit pa.
4. **Sa usage limit, nagpapatuloy ang usapan.** Kapag naabot ng account ng agent ang limit nito, itinatala ng ccam ang limit gaya ng gagawin ng [plugin hook](/features/usage#the-limit-record) at itinutuloy ang parehong usapan bilang ang account na may pinakamaraming natitira: kinokopya ito roon ng Claude Code (`--resume <file> --fork-session`), na may unang mensaheng "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." Isang beses lang sinusubukan ang bawat account.
5. **Itinatala ang bawat usapan ng agent** sa `~/.local/state/ccam/agent/<session id>.json` (30 araw). Itinutuloy ito ng `--resume <id>` sa folder nito, bilang ang account na nagpatakbo nito kung may natitira pa iyon, kung wala ay bilang ibang account. Hindi isinasama ng [`ccam next --continue`](/features/next) ang mga usapan ng agent, kaya ang mga usapan mo pa rin ang itinutuloy nito.

## Kaligtasan at mga espesyal na kaso

- Hindi makapagsisimula ang isang agent ng sarili nitong mga agent: minamarkahan ito ng ccam ng `CCAM_AGENT=1`, at tumatangging tumakbo roon ang `ccam agent`.
- Ang paghinto sa `ccam agent` (Ctrl+C, isang signal, o ang Claude Code na nag-i-interrupt sa command) ay nagpapahinto rin sa agent.
- Hindi nakikita ng agent ang usapang tumawag dito, ang task lang nito.
- Nag-o-overwrite sa trabaho ng isa't isa ang mga agent na sabay na nag-e-edit ng parehong mga file. Bigyan ang bawat isa ng sarili nitong bahagi ng task, o sarili nitong git worktree.
- Kapag lumipat sa ibang account ang tumawag na session ([`auto_switch`](/features/next#auto-switch)), humihinto kasama nito ang mga agent na sinimulan nito. Nananatili ang kanilang mga usapan: itinutuloy ang mga ito ng `--resume <id>`.
- Pangalan ng command ang `agent`, kaya walang bagong account na puwedeng tawaging `agent`; nagsisimula pa rin ang dati nang account sa `ccam run agent`.

## Kaugnay

- [Pagpili ng account para sa mga script](/features/pick)
- [Paglipat sa susunod na account](/features/next)
- [Mga usage limit](/features/usage)
- [Mga background session ng bawat account](/features/agents)
