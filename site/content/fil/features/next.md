---
title: Paglipat sa susunod na account
description: Sa ccam next, bubukas ang Claude Code bilang ibang account; --auto para sa pinakamaluwag, --continue para dalhin ang usapan, --fallback para ulitin ang -p job.
updated: 2026-11-05
order: 8
translation: machine
sourceHash: "de1fd573f1ac"
---

Binubuksan ng `ccam next` ang Claude Code bilang ibang account nang hindi mo ito pinapangalanan. Gamitin ito kapag naabot ng isang account ang usage limit nito. May apat na paraan para lumipat:

- Kinukuha ng simpleng `ccam next` ang account na kasunod ng huli mong ginamit, at kinukuha ng `--auto` ang may pinakamaraming natitira. Ipinagpapatuloy ng `--continue` ang huling usapan ng folder na ito sa bagong account, nagshe-share man o hindi ng kasaysayan ang dalawa.
- Sa loob ng Claude Code, sinasabi sa iyo ng [ccam plugin](#inside-claude-code) kapag pinahinto ng limit ang isang session, at ipinapasa ng `/ccam:switch` ang usapan.
- Para sa mga script na nagpapatakbo ng `claude -p`, kusang ginagawa ng `ccam --fallback` ang paglipat (tingnan ang [Patakbuhin ulit sa limit](#run-again-at-the-limit)).
- Kapag naka-on ang [`auto_switch`](#auto-switch) (Linux at macOS), inililipat din ng ccam ang interactive session na pinahinto ng Claude Code sa limit nito.

**Mga seat ng Team at Enterprise.** Lumilipat din ang awtomatikong pagpili ng account (`ccam next`, `--fallback`, `ccam pick`, ang `n` ng dashboard at, kapag binuksan mo ito, ang [`auto_switch`](/features/next#auto-switch)) sa pagitan ng mga seat ng iisang Claude Team o Enterprise na organisasyon, na maaaring lumabag sa terms ng Anthropic. Basahin ang [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) at magtanong muna sa iyong organisasyon at sa Anthropic bago ka umasa rito.

## Paggamit

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

Ipi-print ng ccam ang paglipat bilang `→ <from> → <to>` sa standard error, kasama ang dahilan ng pagpili ng `--auto`, saka sisimulan ang Claude Code bilang ang bagong account sa parehong paraan ng `ccam <name>`. Ang `<from>` ay ang account na huli mong ginamit, o sa `--continue`, ang account na kinabibilangan ng usapan. Ipinapasa sa `claude` nang walang pagbabago ang mga argument na kasunod ng `--`.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--auto` | Piliin ang naka-log in na account ayon sa usage (tingnan sa ibaba), at laktawan ang mga account na nasa 100% | off |
| `--max-usage <1-100>` | Kasama ang `--auto`: unahin ang mga account na mas mababa sa usage na ito | `auto_max_usage`, kung wala ay 100 |
| `--strategy headroom\|expiring\|balance` | Kasama ang `--auto`: ang pinakamababang 5-oras na usage, ang pinakamalapit na 5-oras na reset o ang pinakamaraming natitirang lingguhang quota bawat araw | `auto_strategy`, kung wala ay `headroom` |
| `--cross-org` | Kasama ang `--auto`: pumili rin ng mga account ng ibang organisasyon (tingnan sa ibaba) | `cross_org`, kung wala ay off |
| `--wait[=<duration>]` | Kasama ang `--auto`: kapag walang account na may natitira, hintayin ang unang mag-reset, nang hindi hihigit sa tagal na ito. Inilalagay ang value pagkatapos ng `=`: tinatanggihan ang `--wait 2h` kasama ang pahiwatig na `--wait takes its value after '=': --wait=2h` | off; `5h` kapag ibinigay nang mag-isa |
| `-c`, `--continue` | Ipagpatuloy ang huling usapan ng folder na ito sa bagong account (tingnan sa ibaba) | off |
| `--brief` | Kasama ang `--continue`: magsimula ng bagong usapan bilang bagong account mula sa buod ng ccam ng huling usapan, sa halip na i-resume ang buong usapan ([tingnan sa ibaba](#brief)) | off |
| `-- <claude args>` | Mga argument na ipinapasa sa `claude` pagkatapos ng sariling argument ng ccam | wala |
| `--dry-run` | Sinusuportahan: ipi-print ang `→ <from> → <to>` at ang command line na patatakbuhin nito, nang hindi sinisimulan ang Claude Code o itinatala ang pagsisimula | off |

Gumagana nang magkasama ang `--auto` at `--continue`. Sa `--continue`, inilalagay ang mga argument ng `--resume` bago ang sarili mong mga argument.

Ang sariling mga option ng Claude Code, gaya ng `--dangerously-skip-permissions` o `--model`, ay inilalagay pagkatapos ng `--`. Hindi tulad ng `ccam <name>`, na ipinapasa ang lahat ng kasunod ng pangalan, binabasa mismo ng `ccam next` ang mga argument nito, kaya kapag walang `--` ay humihinto ito sa `ccam: unknown flag: --dangerously-skip-permissions`.

## Paano ito gumagana

Isinusulat ng bawat pagsisimula sa pamamagitan ng ccam ang oras sa `~/.local/state/ccam/last-used/<name>`. Binabasa ng `ccam next` ang mga file na ito, tinitingnan kung aling mga account ang naka-log in, at pumipili.

1. Bumubuo ang ccam ng isang `next.Candidate` bawat account mula sa `config.toml` (`cli.candidates`): ang account, ang huling oras ng pagsisimula nito (`launcher.LastUsed`) at kung may gumaganang login ito (`creds.Has`: isang Claude login sa `.credentials.json`, o sa macOS sa Keychain, na hindi pa buong nag-expire o binura ng Claude Code). Walang plan limit ang folder na gumagamit ng API key sa halip (isang Console login, o account na idinagdag gamit ang `ccam add --api`) at hindi ito kailanman pinipili.
2. Kapag walang `--auto`, inaayos ng ccam ang mga account ayon sa slot number at hinahanap ang pinakahuling sinimulan (`next.Rotate`). Iyon ang account na nililisan mo. Kung wala pang account na nasimulan, ito ang pinakamababang slot.
3. Mula roon, sumusulong ito ayon sa pagkakasunod-sunod ng slot, bumabalik sa una mula sa huling slot, at ibinabalik ang unang account na naka-log in.
4. Inihahanda at sinisimulan ng ccam ang napiling account nang eksaktong gaya ng gagawin ng `ccam <name>` (`launcher.LaunchWith`): sinusuri nito ang mga link ng account, inilalapat ang MCP sharing, nirerehistro ang session at pinapalitan ang ccam ng `claude` sa pamamagitan ng `syscall.Exec`.

Nagsisimula ang pagdaan sa account na huling ginamit at humihinto sa unang naka-log in na account kasunod nito:

<!-- motion:next-rotate -->

Ang `work` ang pinakahuling sinimulan, kaya nagsisimula ang pagdaan pagkatapos nito. Walang login ang `side` at nilalaktawan ito, at pinipili ang `lab`. Kung naka-log out din sana ang `lab`, babalik ang pagdaan sa `personal` (`% len` sa diagram). Hindi kailanman ibinabalik ng rotation ang account na pinagmulan mo: kapag walang ibang naka-log in na account, humihinto ito na may error.

### --auto

Nagpapasya ang `--auto` ayon sa usage sa halip na sa pagkakasunod-sunod ng slot. Pinapanatili ng ccam ang mga naka-log in na account, sabay-sabay na hinihingi ang usage nilang lahat sa parehong paraan ng `ccam usage` (`Fetcher.FetchAll`), at pumipili ayon sa mga patakaran sa ibaba (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

Hinuhusgahan ang isang account ayon sa 5-oras at lingguhang window nito, at, kapag kasama ang `--model` sa mga argument para sa Claude Code (pagkatapos ng `--`), ayon din sa lingguhang window ng model na iyon (`opus`, `sonnet`, `haiku` o `fable` sa pangalan). Hinahati ng ccam ang mga account sa tatlong antas (`Policy.Choose()`):

1. **May natitira**: bawat window ay mas mababa sa `auto_max_usage` (`config.toml`, default 100). Ang strategy ang pumipili sa kanila: kinukuha ng `headroom` (default) ang pinakamababang 5-oras na usage, at sa mga account na pareho ang 5-oras na usage, ang may pinakakaunting nagamit sa linggo nito; kinukuha ng `expiring` ang 5-oras na window na pinakamalapit mag-reset, para magamit ang natitirang quota nito bago ito mawala; kinukuha ng `balance` ang account na may pinakamaraming natitirang lingguhang quota bawat araw hanggang mag-reset ang linggo nito, na ibinibilang laban dito ang bawat tumatakbong ccam session, para hindi masayang ang lingguhang quota at kumalat ang mga parallel session.
2. **Lampas**: mas mababa sa 100% pero nasa o lampas sa `auto_max_usage`. Ginagamit lang kapag walang account na may natitira; nauuna ang may pinakamababang 5-oras na usage, at kapag tabla, ang may pinakakaunting nagamit sa linggo nito.
3. **Puno**: may window na nasa 100%, o hindi alam ang 5-oras na usage. Hindi kailanman pinipili.

Itinuturing na 100% ang isang window hanggang sa reset nito kapag tinanggihan doon ng Claude Code ang isang request: itinatala ito ng hook ng [ccam plugin](#inside-claude-code) (`~/.local/state/ccam/usage/<name>.limit.json`), anuman ang sabihin ng mga naka-save na numero. Kaya hindi pinipili ulit bago ang reset nito ang account na kaaabot lang sa limit nito, kahit luma na ang mga numero nito.

Inilalagay ang mga setting para sa lahat ng account sa itaas ng `config.toml`, bago ang unang `[[group]]`, `[[profile]]` o `[[folder]]`. Ang setting na nasa ibaba ng isa sa mga iyon ay pag-aari ng table na iyon: nagbababala tungkol dito ang ccam sa bawat pagsisimula, at tumatanggi ang mga command na nagse-save ng config hanggang mailipat ito.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Mga organisasyon.** Nananatili ang mga awtomatikong pagpili sa mga account na kauri ng nililisan mo: mga seat sa parehong Team o Enterprise na organisasyon, o mga personal na plan (Pro, Max), kaya hindi naitutuloy sa pribadong account ang trabaho mula sa account ng kompanya, o ang kabaligtaran. Alam ng ccam ang organisasyon ng bawat account mula sa login nito (itinatala ito ng `ccam login`; binabasa mula sa `.claude.json` ang mas lumang mga login); puwedeng itugma sa anumang account ang account na hindi matukoy ng ccam ang organisasyon. Iisa ang limit ng dalawang folder na naka-log in sa parehong Claude account, kaya hindi kailanman pumapalit ang isa sa isa pa. Inaalis ng `--cross-org`, o ng `cross_org = true` sa `config.toml`, ang patakaran ng organisasyon, pero pagkatapos ng [`/ccam:switch`](#inside-claude-code), hinihintay pa rin ng paglipat sa ibang organisasyon ang pag-oo mo.

Kapag walang mapipili, lumalabas ang ccam na may code 75 (`EX_TEMPFAIL`), para matukoy ng script ang "subukan mamaya" mula sa isang pagkabigo. Sa `--wait`, natutulog ito sa halip hanggang mag-reset ang humaharang na window ng unang punong account (30 segundo pagkatapos ng oras ng reset, kapag ipinapakita na ng endpoint ang bagong window), nagpi-print ng `ccam: no account has room; waiting until 15:10`, at pumipili ulit; kapag walang account na nagsasabi kung kailan ito magre-reset, titingin ulit ito pagkalipas ng 10 minuto. Kung walang lumuwag bago matapos ang paghihintay, lumalabas ito na may 75.

Sinusunod ng paghahanap ang mga patakaran ng `ccam usage`: muling ginagamit ang mga numerong wala pang 60 segundo ang tanda, nagta-time out ang paghahanap ng bawat account pagkalipas ng 3 segundo, ang account na hindi sumasagot, o nag-expire ang access token, ay hinuhusgahan ayon sa mga naka-save na numero nito, at nangingibabaw sa mga iyon ang tala ng limit. Itinuturing na 0% ang naka-save na window na lumipas na ang oras ng reset, kaya hindi nilalaktawan ang account na nag-reset na ang window mula noon. Hindi kailanman nire-renew ng `--auto` ang mga token. Hindi tulad ng simpleng rotation, puwedeng piliin ng `--auto` ang account na huli mong ginamit kung ito pa rin ang may pinakamaraming natitira.

### --continue

1. Hinahanap ng ccam ang huling usapan ng folder na ito sa bawat account: sa mga `<session>.jsonl` file sa `projects/<folder>/` (ang pangalan ng folder ay ang working folder, na naresolba ang mga link, na ginawang `-` ang bawat character na hindi titik o numero, gaya ng pagpapangalan ng Claude Code), ang may pinakabagong prompt o sagot. Nilalaktawan ang mga session na wala nito, gaya ng session na nagpatakbo lang ng `/mcp`, `/usage` o `/status`, o ang walang lamang file. Kung nag-iwan ng tala ang isang session (isang usage limit o `/ccam:switch`, tingnan sa ibaba) at walang usapang may prompt o sagot na lampas sa 2 minuto pagkatapos nito, ang eksaktong session na iyon ang ginagamit; hindi binibilang ang mga hilerang isinusulat ng Claude Code sa `/exit`.
2. Ang account na kinabibilangan nito ang account na nililisan mo, at pinipili ang bagong account gaya ng nasa itaas, sa loob ng parehong organisasyon. Kapag ilang account ang nagshe-share ng kasaysayang iyon, kabilang ito sa account na ang `.claude.json` ay nagsasabing ito ang huling nagpatakbo ng usapan (`lastSessionId`), kung hindi ay sa huling ginamit.
3. Kung shine-share ng bagong account ang kasaysayang iyon (ang parehong group folder), sinisimulan ng ccam ang `claude --resume <session id>`, ang mismong usapang iyon (bubuksan ng `--continue` kung alinman ang pinakabago sa folder).
4. Kung hindi, sinisimulan nito ang `claude --resume <the conversation's file> --fork-session`: kinokopya ng Claude Code ang buong usapan, sa ilalim ng bagong session ID, sa sariling kasaysayan ng bagong account, at nananatili ang orihinal gaya ng dati. Mula noon, nahahanap na ito ng `--continue` ng bagong account.
5. Maliban kung magbibigay ka ng `--agent` pagkatapos ng `--`, ipinapasa rin ng ccam ang agent na ginagamit ng usapan: ang setting na `agent` ng `.claude/settings.local.json` o `.claude/settings.json` ng folder o ng `settings.json` ng account, o kung wala, ang agent na pinapangalanan ng huling `agent-setting` na hilera ng usapan. Hindi ibinabalik ng Claude Code ang agent kapag kinokopya nito ang isang usapan (hakbang 4), at magpapatuloy sana ito bilang default na agent nito, na may lahat ng tool. Kapag hindi shared ng bagong account ang mga agent, plugin at setting ng luma (ang parehong mga folder na `agents` at `plugins` at `settings.json`, sa pamamagitan ng isang group), nagdadagdag ang ccam ng linya pagkatapos ng linya tungkol sa paglipat: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

Sinasabi ng ccam kung gaano kalaki ang usapan: binabasa ulit ng bagong account ang lahat ng ito sa unang sagot nito. Hindi kailanman shine-share ang mga prompt cache sa pagitan ng mga organisasyon, at sariling organisasyon ang bawat personal na plan, kaya kapag magkaiba ang kinabibilangan ng dalawang account, idinadagdag ng ccam ang `prompt caches are not shared between organizations, so the first reply reads it all again`. Kapag naka-sign in ang lumang account sa mga MCP server na hindi pa naka-sign in ang bago, pinapangalanan ng ccam ang mga ito: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. Kapag shine-share ang `history` sa isang group, naiiwasan ang pagkopya, nananatiling iisa ang usapan, at shine-share din ang mga checkpoint na pinagbabalikan ng code ng `/rewind` (`file-history/`).

### --brief

Nire-resume ng `--continue` ang buong usapan, at binabasa ulit ng bagong account ang lahat ng ito sa unang sagot nito. Sa halip, nagsisimula ang `ccam next --continue --brief` ng bagong usapan mula sa buod na isinusulat ng ccam tungkol sa luma, nang hindi nagtatanong sa isang model:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

Nilalaman ng buod, sa ganitong pagkakasunod-sunod at kapag may masasabi lang:

1. ang pamagat ng usapan;
2. ang path ng transcript nito, para hanapan ni Claude kapag kailangan nito ng detalye;
3. ang una mong request, kapag higit sa apat ang mga ito;
4. ang buod na isinulat ng Claude Code noong huli nitong kinompact ang usapan;
5. ang huli mong apat na request, kasama ang mga `!` command at slash command gaya ng pagkaka-type mo;
6. ang huling sagot ni Claude;
7. ang todo list;
8. ang 10 file na huling binago ni Claude;
9. ang `git status --short --branch` ng folder, hanggang 800 character.

Hanggang 10 item lang ang itinatago ng todo list: nauunang tinatanggal ang mga tapos na item, saka ang mga hakbang na pinakamalayo pa, kaya naiiwan ang kasalukuyang hakbang at ang mga kasunod. Pinuputol sa gitna ang mahahabang bahagi, at nananatiling mas mababa sa 16,000 character ang buong buod.

Sine-save ito ng ccam bilang `~/.local/state/ccam/brief/<session>.md`, na ikaw lang ang makakabasa, at tinatanggal ng susunod na `--brief` ang mga buod na mas matanda sa 7 araw. Nagsisimula ang Claude Code sa iisang prompt, at ina-attach ang file na binabanggit nito: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` Pumapalit sa huling pangungusap ang prompt na ibinigay mo pagkatapos ng `--` (`Then: <your prompt>`), at sumusunod dito ang iba mo pang option pagkatapos ng `--`.

Tinatapos ng Claude Code ang isang mention sa `"` at pinuputol ito sa `#`, kaya kapag may alinman dito ang path ng buod, binabanggit ito ng ccam sa path nito na relative sa working folder; kapag mayroon din nito ang path na iyon, pinapangalanan ng prompt ang file, at binabasa ito ni Claude gamit ang Read tool nito, na nagtatanong muna dahil nasa labas ng mga working folder ang file (sa isang `-p` run sa Manual mode, tinatanggihan ng Claude Code ang pagbasa at nagpapatuloy si Claude nang walang buod).

Ipi-print ng `ccam --dry-run next --continue --brief` ang buod at walang sine-save. Nananatili ang lumang usapan gaya ng dati, at ang sinasabi lang ng buod ang alam ng bago, kaya gamitin ang `--brief` kapag mas mahalaga ang gawain kaysa sa bawat detalye kung paano ka nakarating doon. Kapag walang `--continue`, humihinto ang ccam na may `--brief works with --continue`; kapag hindi nito mabasa ang usapan o wala itong makitang maibubuod dito, humihinto ito at sinasabing nire-resume nang buo ng `ccam next --continue` na walang `--brief` ang usapan. Laging nire-resume ng [relay](#relay) at ng [`auto_switch`](#auto-switch) ang buong usapan. Inilalagay pagkatapos ng `--` ang sariling option na `--brief` ng Claude Code, gaya ng iba pa nitong option.

## Sa loob ng Claude Code

Nagdadagdag ang [ccam plugin](/docs/install#claude-code-plugin) ng dalawang hook (Claude Code 2.1.78 o mas bago):

- Kapag huminto ang isang turn sa usage limit, itinatala nito ang usapan at nagpapakita ng desktop notification sa mga terminal na sumusuporta nito (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  Kilala nito ang bawat mensaheng itinuturing ng Claude Code na usage limit: ang 5-oras, lingguhan, karaniwan at model na limit (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` para sa allowance ng isang model), ang mga usage credit at spend limit, at ang organisasyong nagpahinto sa account.

  Kapag naisulat na ng Claude Code sa usapan ang tinanggihang request, itinatala rin ng hook ang punong window hanggang sa reset nito ([ang tala ng limit](/features/usage#the-limit-record)), maliban sa pagpapahintong itinakda ng isang organisasyon, na puwedeng alisin ng mga admin nito anumang oras. Hindi usage limit ang maikling throttle ng server o ang rate limit ng isang API key, at hinahayaan lang ang mga iyon, gayundin ang mga session na walang nagbabantay (mga `-p` run at background session). Itinuturing na paghinto ng session nito ang isang subagent o workflow agent na huminto sa usage limit, kahit naghihintay lang dito ang pangunahing usapan: isinusulat ng Claude Code ang tinanggihang request na iyon sa sariling file ng agent sa `subagents` folder ng session, at doon ito binabasa ng hook. Kapag sabay-sabay na huminto ang ilang agent, iisang tala at iisang notification lang ang natatanggap ng session.

  Naghihintay rin sa session ang Claude Code 2.1.234 at mas bago, at kusang nagpapatuloy pagkatapos ng 5-oras na reset; binibigyan ka ng notification ng pagpipiliang lumipat na lang, at kapag naka-on ang [`auto_switch`](#auto-switch) (Linux at macOS), ang ccam ang lilipat para sa iyo kapag may natitira ang ibang account. Kapag walang ibang account na puwedeng pumalit, dahil walang ibang account na kabilang sa parehong Team o Enterprise na organisasyon, o walang ibang personal na plan (kapag naka-on ang `cross_org`, puwede ang anumang ibang Claude account), nagtatapos ang notification sa `You have no other personal plans, so ccam cannot move this conversation.` o `You have no other accounts in Acme, so ccam cannot move this conversation.` sa halip.
- Itinatala ng `/ccam:switch [account]` ang usapan at tinatapos ang turn bago gumawa ng anumang request: i-type ang `/exit`, saka patakbuhin ang `ccam next --continue` para ipagpatuloy ito sa account na pinangalanan mo (o sa kasunod na account). Sa session na wala pang usapan, sinasabi nitong walang maipagpapatuloy. Tanging ang tina-type mo ang nagpapatakbo nito: mula sa plugin 0.3.2, hindi kayang patakbuhin ni Claude ang `/ccam:switch`, at ipinapasa ito ng job na ini-schedule ni Claude (CronCreate, ScheduleWakeup o `.claude/scheduled_tasks.json`) kay Claude bilang plain text.

  Kapag nasa labas ng organisasyon ng usapan ang account na pinangalanan mo (seat ng ibang Team o Enterprise na organisasyon, personal na plan para sa usapang nasa seat o ang kabaligtaran, o API key account sa alinmang panig), nagtatanong muna ang `ccam next --continue`, anuman ang sabihin ng `cross_org`: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. Ganoon din ang itinatanong nito bago ang account na pinili nito mismo kapag wala kang pinangalanan, na pinapayagan ng `cross_org = true` na kunin mula sa ibang organisasyon: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. `y` o `yes` lang ang naglilipat nito; sa Linux at macOS, itinatapon muna ang mga key na na-type bago ang tanong. Anumang ibang sagot, o kapag walang terminal na mapagtatanungan, ay iniiwan ang usapan at ang tala kung nasaan sila: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (kung ang ccam ang pumili ng account, pipili ulit ang `ccam next --continue`). Gaya sa mga awtomatikong pagpili, bagay sa alinman ang account na hindi matukoy ng ccam ang organisasyon.

Nasa `~/.local/state/ccam/handoff/<account>/<session>.json` ang tala, isang file bawat session, kaya hindi kailanman napapatungan ng dalawang session ng iisang account ang tala ng isa't isa. Itinatago ito nang 12 oras at isang beses lang ginagamit, kapag magsisimula na ang susunod na session; kapag pumalya ang pagsisimula, naiiwan ito para sa susunod na subok. Nakatali ito sa folder kung saan nagsimula ang session, kahit pagkatapos ng `cd` sa shell ng Claude Code. Ina-update ng `ccam rename` ang mga pangalan dito; kung natanggal na ang account na binabanggit nito, sasabihin iyon ng ccam at pipili ito gaya ng dati.

### Relay

Kapag may `relay = true` sa itaas ng `config.toml`, patuloy na tumatakbo ang ccam kasabay ng mga interactive session sa halip na palitan ang sarili nito ng Claude Code, at pagkatapos ng `/ccam:switch`, kusa nitong ipinagpapatuloy ang usapan ng session na iyon (hindi nito pinapansin ang talang iniwan ng ibang session): tinatapos ng hook ang Claude Code, at pinapatakbo ng ccam ang `ccam next --continue` gamit ang tala ng session na iyon. Kapag walang pinangalanang account, kinukuha nito ang may pinakamaraming natitira, na isinasaalang-alang ang model ng usapan, at kapag pinapayagan ng `cross_org = true` na maging account ng ibang organisasyon iyon, itinatanong nito ang tanong sa itaas bago nito ilipat doon ang usapan. Sa Windows, hindi kayang tapusin ng hook ang Claude Code: i-type ang `/exit`, at ipagpapatuloy ito ng ccam. Hindi rin ito tinatapos ng hook para sa account na nasa labas ng organisasyon ng usapan: sinasabi nito ang `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`, at itinatanong ng ccam ang tanong sa itaas kapag na-type mo na ang `/exit`. Isang naghihintay na ccam process bawat session ang kapalit nito, kaya naka-off ito bilang default.

Pinapanatili ng ipinasang session ang hindi naibabalik ng pag-resume ng usapan: mga option gaya ng `--add-dir`, `--mcp-config`, ang mga listahan ng tool, `--settings`, `--plugin-dir` at ang mga option ng system prompt, ang permission mode ng usapan, ang effort level ng session at ang agent na ginagamit ng usapan (tingnan ang [`--continue`](#continue)). Ang Claude Code mismo ang nagbabalik ng `--model` at `--name` mula sa usapan. Pinapangalanan sa `ccam: not carried over: …` ang anumang iba pa (`--worktree`, `--debug`, `--remote-control`, mga option na hindi kilala ng ccam). Hindi kailanman itinatanong ng paglipat ang tanong tungkol sa shared na kasaysayan, at itinatapon ang anumang na-type mo habang tumatakbo ito; ang tanging tanong nito ay ang tanong bago ang account ng ibang organisasyon.

Ipinagpapatuloy lang ng relay ang usapan kapag natapos ang Claude Code sa paraang hiniling ng `/ccam:switch`: hindi kapag nauna kang nag-type ng `/exit`, pumindot ng Ctrl+C habang natatapos ito, o nakatanggap ang ccam ng hangup o ibang signal, at hindi kapag wala na ang terminal. Sa ganoong kaso, naiiwan ang tala, at ipinagpapatuloy ito ng `ccam next --continue` sa ibang pagkakataon. Hindi tinatapos ng hook ang session na sinimulan sa ilalim ng relay ng mas lumang bersyon ng ccam, o ang session na walang nagbabantay (mga `-p` run at background session); ang mensaheng nagsasabing i-type ang `/exit` ang nakukuha nito sa halip. Tumatakbo nang walang relay ang `claude remote-control`.

### `auto_switch`

Kapag may `auto_switch = true` sa itaas ng `config.toml` (Linux at macOS), kusang inililipat ng ccam ang isang session sa ibang account kapag pinahinto ito ng Claude Code sa usage limit: pumipili ang hook ng plugin ng ibang account na may natitira at tinatapos ang Claude Code, at ang ccam, na tumatakbo kasabay ng session, ang nagpapatuloy ng usapan doon, gaya ng ginagawa ng [relay](#relay) pagkatapos ng `/ccam:switch`. Naka-off ito bilang default at eksperimental.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

Inililipat lang ng ccam ang isang session kapag totoo ang lahat ng ito:

1. Sinimulan ito ng ccam sa isang terminal at patuloy na tumatakbo ang ccam kasabay nito. Binubuksan ng `auto_switch` ang [relay](#relay), kaya totoo ito para sa mga interactive session na sinimulan mo pagkatapos itong itakda. Ang session na nagsimula bago iyon, o sa ilalim ng mas lumang ccam, ay nakakakuha ng notification na nagtatapos sa ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` Hindi kailanman inililipat ang mga session na walang nagbabantay (mga `-p` run at background session), ang `claude remote-control` at ang Claude Code na sinimulan nang walang ccam.
2. Naka-on sa account ang [ccam plugin](/docs/install#claude-code-plugin) at Claude Code 2.1.214 o mas bago ang pinapatakbo nito, na nagpapasa ng process ID nito sa hook (`CLAUDE_PID`).
3. Naisulat ng Claude Code sa usapan ang tinanggihang request, kasama ang oras ng reset nito, para sa 5-oras, lingguhan, Opus o Sonnet na limit ng account na iyon. Hindi kailanman naglilipat ng session ang spend limit, ang pagkaubos ng mga usage credit, ang pagpapahintong itinakda ng isang organisasyon, ang Fable allowance (`You've reached your Fable limit.`), ang throttle ng server at ang rate limit ng isang API key. Ibinibilang ang limit ng isang subagent o workflow agent, kahit naghihintay lang dito ang pangunahing usapan; kapag lingguhang Opus o Sonnet limit iyon, ililipat lang ang session kung ang usapan mismo ang tumatakbo sa model na iyon.
4. Alam ng ccam ang plan ng account, mula sa `.claude.json` nito o sa itinala ng `ccam login`. Kung hindi, nagtatapos ang notification sa `ccam does not know personal's plan, so auto_switch leaves this session.`
5. Hindi nakakonekta sa Remote Control ang session. Bilang ibang account, hindi na ito muling kokonekta sa app na ginagamit mo para patakbuhin ito, kaya nananatili ito, at sinasabi iyon ng notification.
6. Magre-reset ang limit nang lampas 15 minuto mula ngayon. Para sa mas maikling paghihintay, ipinapaubaya ng ccam ang session sa Claude Code: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. Hindi ka pa nagta-type ng prompt mula nang maabot ang limit.
8. May natitira ang ibang account. Pumipili ang ccam gaya ng [`--auto`](#auto) (`auto_max_usage`, `auto_strategy`, ang model ng usapan), sa mga account na kauri ng nililisan mo: ibang seat ng parehong Team o Enterprise na organisasyon, o ibang personal na plan (Pro, Max). Hindi nalalapat dito ang `cross_org` at `--cross-org`. Ang account ay dapat:
   - naka-log in gamit ang login na kayang i-renew ng Claude Code;
   - hindi naka-log in sa parehong Claude account;
   - hindi gumagamit ng API key;
   - may plan na alam ng ccam;
   - walang naka-on na usage credit, ayon sa ipinapakita ng sariling paghahanap ng usage ng ccam;
   - pinagkatiwalaan na ang folder na ito sa Claude Code;
   - kung tumatakbo ang usapan bilang isang agent (pinapangalanan ito ng setting na `agent` o ng usapan), naka-share dito ang mga agent ng account na nililisan mo: dapat iisa talaga ang bawat isa sa mga folder nitong `agents` at `plugins` at ang `settings.json` nito (sa pamamagitan ng isang group), o wala ang mga ito sa pareho.

   Kapag inaalis ng huling tatlong iyon ang bawat account na may natitira, pinapangalanan sila ng notification: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` o `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` Itinuturing na puno ang window na may [tala ng limit](/features/usage#the-limit-record), at nahuhuli sa iba ang account na nilipatan ng ibang session sa nakaraang 2 minuto. Kailangang mapili ito sa loob ng mga 6 na segundo; kapag mas matagal ang mga paghahanap ng usage, nananatili ang session. Sa `usage_source = "claude"`, nagpapasya ang `auto_switch` mula sa mga naka-save na numero lang, at hindi matukoy ng ccam mula sa paghahanap ng Claude Code o sa status line kung naka-on ang mga usage credit, kaya ipinapalagay nitong naka-off ang mga ito.

Saka tinatapos ng hook ang Claude Code, at ang ccam, na tumatakbo kasabay nito, ay nagpapakita ng desktop notification sa mga terminal na sumusuporta nito, gaya ng `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` Ipi-print nito ang unang linya ng halimbawang output sa itaas at pinapatakbo ang `ccam next --continue` para sa session na iyon bilang ang napiling account, na nagpi-print ng dalawa pa nitong linya at sinisimulan ang Claude Code kasama ang mga option, permission mode, effort level at agent ng session, gaya ng ginagawa ng [relay](#relay). Ang unang mensahe ng bagong session ay isang nakapirming teksto, ang tanging ipinapadala ng ccam sa ngalan mo (walang ipinapadala ang paglipat pagkatapos ng `/ccam:switch`):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

Pagkatapos ng limit ng isang subagent, ito ang mensahe, dahil huminto kasabay ng lumang session ang mga subagent, workflow at background command:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

Kung maabot din ng account na iyon ang limit nito, ganoon din ang mangyayari roon. Pinipigilan ng mga tala ng limit ang kadena na bumalik sa isang account bago ang reset nito, kaya kapag puno ang bawat account, nananatili ang session sa huli.

Nalalaman ng hook na tumatakbo bilang isang agent ang usapan mula sa setting na `agent`, o mula sa hilerang nagpapangalan sa agent na isinusulat ng Claude Code sa usapan paminsan-minsan at kapag natapos ang session. Kung sa ngayon ay ang `--agent` lang ng session ang nagpapangalan dito, at hindi shared ng napiling account ang mga agent, walang sinisimulan ang `ccam next --continue` at itinatago nito ang tala, at sinasabi iyon ng ikalawang desktop notification: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

Kapag walang ibang account na may natitira, walang lumilipat: hinahawakan ng Claude Code ang limit gaya ng dati, naghihintay sa session o nag-aalok ng mga option nito, at sinasabi ng notification kung aling account ang unang luluwag, halimbawa `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` Para sa mga seat, sinasabi nito ang `None of your other accounts in Acme has room now`. Para pigilan ang paglipat na nagsimula na, pindutin ang Ctrl+C habang nagsasara ang Claude Code: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

May kapalit ang paglipat. Nagre-restart ang Claude Code, na tumatagal nang ilang segundo, at humihinto kasama ng lumang session ang mga background shell, dev server at agent na sinimulan ni Claude. Binabasa ulit ng bagong account ang buong usapan sa unang sagot nito, gaya sa [`--continue`](#continue), at sinasabi ng ccam kung gaano ito kalaki; hindi kailanman gumagamit ng [`--brief`](#brief) ang paglipat.

Hindi kailanman lumilipat ang `auto_switch` bago mag-ulat ng limit ang Claude Code, hindi kailanman binabantayan ang usage mo sa background, at hindi kailanman kumokopya, nagre-renew o nagpapasa ng login. Sa Windows, hindi ito pinapansin, at sinasabi iyon ng `ccam doctor`; doon, kapag may `relay = true`, inililipat pa rin ng `/ccam:switch` at `/exit` ang usapan sa ibang account.

**Ang terms ng Anthropic.** Sa mga seat ng Team at Enterprise, inililipat ng `auto_switch` ang trabaho mo sa ibang seat ng parehong organisasyon. Sinasabi ng help center ng Anthropic na "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (nakatalaga ang bawat seat ng Team plan sa iisang partikular na tao, at siya lang ang puwedeng gumamit nito) at "Users can't route their work through another user's account to use that user's remaining usage" (hindi puwedeng idaan ng mga user ang trabaho nila sa account ng ibang user para magamit ang natitirang usage ng user na iyon) ([mga Team plan](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); ganoon din sa [mga Enterprise plan](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)).

Para sa mga personal na plan (Pro, Max), hindi sinasabi ng terms ng Anthropic kung puwedeng awtomatikong pumalit sa isa't isa ang sarili mong mga account sa limit. Ipinagbabawal ng [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) ang pag-access sa Services "through automated or non-human means, whether through a bot, script, or otherwise" (sa pamamagitan ng awtomatiko o hindi-tao na paraan, gamit man ang bot, script o iba pa) maliban kung tahasang pinapayagan ito ng Anthropic, at ipinapadala ng paglipat ang isa sa mga nakapirming mensahe sa itaas sa ngalan mo. Inililista ng [Usage Policy](https://www.anthropic.com/legal/aup) na magkakabisa sa Nobyembre 12, 2026 ang "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (sadyang lampasan o panghimasukan ang mga kakayahan, paghihigpit o guardrail na itinakda sa loob ng mga produkto nito) sa mga bagay na ipinagbabawal nito.

Hindi sinasabi ng ccam na sumusunod ang `auto_switch` sa mga terms na ito. Ikaw ang magpapasya kung bubuksan ito, at ikaw rin ang mananagot sa panganib: puwedeng suspindihin ng Anthropic ang isang account, o ang access ng isang Team o Enterprise na organisasyon, kapag naniniwala itong nilabag ang terms.

Habang naka-on ang `auto_switch`, nagbababala ang `ccam doctor` tungkol sa account na walang ccam plugin at tungkol sa mga bersyon ng Claude Code na mas luma sa 2.1.214, at para sa mga seat, inuulit nito ang patakaran ng Anthropic sa seat ([mga pagsusuri](/features/doctor)).

## Patakbuhin ulit sa limit

Para sa interactive session na huminto ang `ccam next`. Hindi ka kayang hintayin ng non-interactive na run, ang `claude -p` sa isang script o naka-iskedyul na job, kaya hinahayaan ng `--fallback` bago ang pangalan ng account na ang ccam ang lumipat para dito:

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

1. Gumagana lang ang `--fallback` kasama ang `-p` o `--print`; kung hindi, humihinto ang ccam at itinuturo ang `ccam next --auto --continue`.
2. Isang beses binabasa ang piped na standard input, hanggang 10 MiB, at ibinibigay sa bawat pagtatangka. Gaya ng Claude Code, naghihintay ang ccam nang hindi hihigit sa 3 segundo para magsimula ito (maliban sa `--input-format stream-json`); kung nananatiling bukas ang pipe pero walang ipinapadala, walang input na makukuha ang run at itinatala ito ng ccam sa standard error.
3. Sinisimulan ng ccam ang `claude` bilang child process sa halip na `exec`, pinipigil ang standard output nito, at ipinapasa ang standard error habang dumarating ito.
4. Kapag lumabas ang `claude` nang may 0, o pumalya dahil sa ibang dahilan, ipinapasa nang walang pagbabago ang pinigil na output at ang exit code. Ganoon din ang run na tinapos ng isang signal (exit code na lampas sa 128): hiniling itong huminto, at hindi na ito pinapatakbo ulit.
5. Itinuturing na usage limit ang pagkabigo kapag sinasabi iyon ng mensahe ng Claude Code, alinman sa mga mensahe ng usage limit na kilala ng [hook](#inside-claude-code), pati sa loob ng output na `--output-format json` o `stream-json`, o kapag ipinapakita ng usage ng account na nasa 100% ang 5-oras, lingguhan o `--model` na window nito (kasama ang tala ng limit).
6. Saka itinatapon ang output ng pagtatangkang iyon at pinapatakbo ang parehong command bilang ibang naka-log in na account: ang may pinakamababang 5-oras na usage, saka ang mga account na hindi matingnan ang usage, ayon sa pagkakasunod-sunod ng numero. Nilalaktawan ang mga account na nasa 100%, hindi kailanman sinusubukan ang mga account na may API key login, at sinusubukan ang bawat account nang hindi hihigit sa isang beses.
7. Kapag kasama sa mga argument ang `-c`, `--continue`, `-r` o `--resume`, ang mga account lang sa parehong history group ang sinusubukan, para mahanap ang usapan.
8. Ang mga account lang ng parehong organisasyon ang sinusubukan, gaya sa `--auto` (inaalis ito ng `cross_org = true`).

Pinipili ang susunod na account ayon sa parehong mga patakaran ng `--auto` (`auto_max_usage`, `auto_strategy`, ang window ng `--model`). Kapag wala nang natitirang account, ipi-print ng ccam ang output ng huling pagtatangka, sasabihing walang naka-log in na account na may natitira pa, at lalabas na may code 75. Sa `--wait` bago ang pangalan ng account (`ccam --fallback --wait=2h work -p …`), hinihintay nito sa halip na mag-reset ang unang account, saka puwedeng subukan ulit ang bawat account; lagi itong natutulog nang kahit isang beses muna, kaya hindi agad pinapatakbo ulit ang account na kaaabot lang sa limit nito. Nililimitahan ng tagal ang kabuuang oras ng paghihintay ng run, hindi ang bawat paghihintay. Kapag isinulat nang may espasyo (`--wait 30m`), tinatanggihan ito kasama ang parehong pahiwatig ng `ccam next`. Gumagana ang `--` sa unahan ng mga argument ng Claude Code gaya ng sa anumang pagsisimula (`ccam --fallback work -- -p …`). Ipi-print ng `ccam --dry-run --fallback work -p …` ang command ng unang pagtatangka.

Hindi kusang inililipat ang mga interactive session maliban kung bubuksan mo ang [`auto_switch`](#auto-switch). Kapag may `ccam statusline install`, idinadagdag ng status line ang `→ ccam next --auto --continue` kapag umabot sa 90% ang 5-oras o lingguhang usage ([Mga usage limit](/features/usage)).

## Kaligtasan at mga espesyal na kaso

- Hindi kailanman pinipili ang mga account na walang login. Kapag walang ibang naka-log in na account, humihinto ang simpleng `ccam next` na may `no logged-in account other than <name>`.
- Sa `--continue` at walang usapan sa folder: `no account has a conversation in <folder> to continue`.
- Kapag walang natirang account dahil sa patakaran ng organisasyon: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Sa `--auto` at walang naka-log in na account: `no logged-in account to pick`. Kapag puno ang bawat account o hindi alam ang usage nito: `` no account with known usage and headroom (see `ccam usage`) ``. Parehong nagtatapos ang dalawa sa exit code 75.
- Gumagana lang ang `--max-usage`, `--strategy`, `--wait` at `--cross-org` kasama ang `--auto`, at tinatanggihan ang value na labas sa 1-100 o ang hindi kilalang strategy.
- Kinokopya ng fork ang usapan, hindi ang iba pang state ng account: hindi kasamang dinadala ang mga checkpoint ng `file-history/` ng orihinal na account maliban kung shine-share ang kasaysayan.
- Gumagawa ng mga network request ang `--auto` kahit may `--dry-run`. Nagbabasa lang ito, at hindi kailanman sumusulat sa mga token file.
- Nagsisimula ang napiling account sa parehong daan ng `ccam <name>`, kaya nalalapat ang parehong mga pagsusuri, kasama ang babala kapag may tumatakbo nang ibang Claude Code session na nagshe-share ng parehong kasaysayan sa folder na ito.
- Isinusulat lang ang `last-used` kapag talagang nagsimula ang Claude Code. Hindi ito ginagalaw ng `--dry-run`, kaya hindi isinusulong ng dry run ang rotation.
- Pinipigil ng `--fallback` ang standard output hanggang matapos ang isang pagtatangka, kaya dumarating nang buo sa dulo ang streaming output gaya ng `--output-format stream-json`. Hindi magagamit ang `--fallback` ng client na nagsi-stream ng `stream-json` na input at naghihintay ng mga sagot bago magpadala ng higit pa. Nananatiling binago ang mga file na binago ng unang pagtatangka bago nito naabot ang limit; sinisimulan ulit ng susunod na pagtatangka ang prompt.

## Kaugnay

- [Mga usage limit](/features/usage)
- [Mga usage limit at paglipat ng account](/docs/usage-limits)
- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Kapag naabot mo ang Claude Code usage limit](/guides/usage-limit-workflow)
