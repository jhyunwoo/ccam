---
title: Mga madalas itanong
metaTitle: "ccam FAQ: maraming Claude Code account"
description: Mga sagot tungkol sa ccam, mga Claude Code account, usage limit, sharing, privacy at mga suportadong system.
updated: 2026-11-19
order: 6
translation: machine
sourceHash: "090ecc3a34bd"
---

## Opisyal bang tool ng Anthropic ang ccam?

Hindi. Ang ccam ay hindi opisyal na open-source na proyekto at hindi kaanib o inendorso ng Anthropic. Ginagamit nito ang setting na `CLAUDE_CONFIG_DIR` na [idinodokumento ng Claude Code](https://code.claude.com/docs/en/env-vars) para sa pagpapatakbo ng maraming account nang magkatabi.

## Pinapayagan ba akong gumamit ng maraming Claude Code account?

Nakadepende iyan sa terms ng bawat account at, para sa mga work account, sa patakaran ng iyong organisasyon: ang [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) ng Anthropic para sa Pro at Max, ang [Commercial Terms](https://www.anthropic.com/legal/commercial-terms) nito para sa Team at Enterprise, at ang [Usage Policy](https://www.anthropic.com/legal/aup) nito para sa lahat. Suriin muna ang mga ito bago ka gumamit ng maraming account. Hindi tinataasan, pinagsasama o itinatago ng ccam ang mga limit ng isang account, at humihinto pa rin ang bawat account sa sarili nitong limit. Inililipat nito ang trabaho mo sa iba mo pang account kapag hiniling mo, o kusa sa isang limit kapag binuksan mo ang [`auto_switch`](/features/next#auto-switch) (Linux at macOS). Alamin kung pinapayagan iyon para sa mga account mo bago mo ito gamitin; ikaw ang magpapasya.

**Mga seat ng Team at Enterprise.** Lumilipat din ang awtomatikong pagpili ng account (`ccam next`, `--fallback`, `ccam pick`, [`ccam agent`](/features/agent), ang `n` ng dashboard at ang [`auto_switch`](/features/next#auto-switch) kapag binuksan mo ito) sa pagitan ng mga seat ng iisang Claude Team o Enterprise na organisasyon. Sinasabi ng help center ng Anthropic na "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (nakatalaga ang bawat seat ng Team plan sa iisang partikular na tao, at siya lang ang puwedeng gumamit nito) at "Users can't route their work through another user's account to use that user's remaining usage" (hindi puwedeng idaan ng mga user ang trabaho nila sa account ng ibang user para magamit ang natitirang usage ng user na iyon) ([mga Team plan](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); ganoon din sa [mga Enterprise plan](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Kaya maaaring lumabag sa terms ng Anthropic ang paggamit ng awtomatikong pagpili sa pagitan ng mga seat. Ikaw ang magpapasya kung gagamitin mo ito roon: magtanong muna sa iyong organisasyon at sa Anthropic bago ka umasa rito. Saklaw ng talata sa itaas ang mga personal na plan (Pro, Max).

## Nalulusutan ba ng ccam ang mga usage limit?

Hindi lampas sa sariling limit ng anumang account: may sariling mga limit ang bawat account at humihinto ito sa mga iyon, at hindi tinataasan, pinagsasama o itinatago ng ccam ang mga ito. Ipinapakita ng ccam ang mga ito at hinahayaan kang magpatuloy sa iba mo pang account kapag naubos na ang isa; kapag naka-on ang `auto_switch`, kusa nitong ginagawa iyon (tingnan ang susunod na tanong). Sa pagitan ng mga seat ng iisang Team o Enterprise na organisasyon, ginagamit ng pagpapatuloy sa ibang seat ang natitirang usage ng seat na iyon, na ipinagbabawal ng help center ng Anthropic para sa seat ng ibang user (tingnan ang Mga seat ng Team at Enterprise sa itaas).

## Kaya bang kusang magpalit ng account ang ccam kapag naabot ng isang session ang limit nito?

Kapag binuksan mo lang ito. Kapag may `auto_switch = true` sa itaas ng `~/.config/ccam/config.toml` (Linux at macOS; naka-off bilang default at eksperimental), kapag pinahinto ng Claude Code ang isang interactive session (na sinimulan ng ccam) sa 5-oras, lingguhan, Opus, Sonnet o Fable na limit, tinatapos ng ccam ang session at ipinagpapatuloy ang usapan bilang iba mo pang account na may natitira: ibang seat ng parehong Team o Enterprise na organisasyon, o ibang personal na plan (Pro, Max). Hindi ito kailanman lumilipat bago ang limit, hindi kailanman binabantayan ang usage mo sa background, at hindi kailanman inilalabas ang trabaho mula sa isang Team o Enterprise na organisasyon o ipinapasok dito.

Sa mga seat, inililipat nito ang trabaho mo mula sa isang seat papunta sa iba, na maaaring lumabag sa terms ng Anthropic (tingnan ang Mga seat ng Team at Enterprise sa itaas). Para sa mga personal na plan, hindi ito malinaw na tinutugunan ng terms ng Anthropic: ipinagbabawal ng Consumer Terms ang pag-access sa Services "through automated or non-human means" (sa pamamagitan ng awtomatiko o hindi-tao na paraan) maliban kung tahasang pinapayagan ito ng Anthropic, at nagpapadala ang paglipat ng isang mensahe sa ngalan mo. Ikaw ang magpapasya kung bubuksan ito.

Inililista ng [`auto_switch`](/features/next#auto-switch) ang bawat kondisyon at ang kapalit ng paglipat.

## Paano gumamit ng dalawang Claude Code account sa iisang computer?

Bigyan ang bawat account ng sariling configuration folder gamit ang `CLAUDE_CONFIG_DIR`. Ginagawa ito ng ccam para sa iyo: patakbuhin ang `ccam add` para sa bawat account, saka simulan ang mga ito gamit ang `ccam personal` at `ccam work`, o gamit ang mga shortcut na ikaw mismo ang nagpangalan, gaya ng `ccwork`. Ipinapakita rin ng [gabay](/guides/multiple-claude-code-accounts) kung paano ito gawin nang mano-mano.

## Paano ko itutuloy ang isang usapan sa ibang account?

Patakbuhin ang `ccam next --continue` sa folder: dinadala nito ang huling usapan ng folder sa susunod na account (`--auto` para sa may pinakamaraming natitira). Sa loob ng Claude Code, ganoon din ang ginagawa ng `/ccam:switch` mula sa [ccam plugin](/docs/install#claude-code-plugin) para sa session na kinaroroonan mo. Kapag nagshe-share ng kasaysayan ng usapan ang dalawang account (`ccam group set default --add history`, o sa isang group na para lang sa kanila: `ccam group add pair --items history --members work,side`), binubuksan ng bagong account ang mismong usapang iyon; kung hindi, kinokopya ito ng Claude Code sa kasaysayan ng bagong account. Idagdag ang `--brief` para simulan ang bagong account sa maikling buod sa halip na sa buong usapan.

## Nakikita o ipinapadala ba ng ccam ang login ko?

Hindi kailanman ishe-share o kokopyahin ng ccam ang iyong Claude login sa pagitan ng mga account, at hindi nito kailanman kinokopya ang mga sign-in ng mga MCP server. Binabasa lang nito ang access token ng isang account para itanong sa usage endpoint ng Anthropic ang usage ng account na iyon: para sa `ccam usage`, sa dashboard, sa `ccam doctor` at sa awtomatikong pagpili (`ccam next --auto`, `ccam pick`, `--fallback`, at `auto_switch` kapag naabot ng isang session ang limit nito). Nire-renew lang nito ang expired na token gamit ang `ccam usage --refresh`; kapag `usage_source = "claude"`, wala itong ipinapadalang token kahit saan at binabasa lang nito ang naka-store na login (ang file, o ang entry sa Keychain ng macOS) para malaman kung gumagana pa ito. Binabasa ng `ccam ext` at `ccam doctor` ang mga pangalan ng mga MCP sign-in ng isang account, hindi kailanman ang mga token nila. Walang telemetry ang ccam. Bukod sa mga lookup na iyon, ang tanging request na kusa nitong ginagawa ay ang pagsusuri ng dashboard sa GitHub para sa mas bagong ccam release, na walang dalang anuman tungkol sa iyo o sa iyong mga account.

## Paano mag-sign in sa isang remote server?

Patakbuhin ang `ccam add` o `ccam login <name>` sa SSH session gaya ng dati. Naghihintay ang Claude Code ng sign-in sa server, kaya mapupunta ang browser sa computer mo sa isang `http://localhost:<port>/callback?code=…` na page na hindi naglo-load. Kopyahin ang address na iyon mula sa browser at i-paste sa terminal: ipapasa ito ng ccam sa Claude Code sa server, na siyang tatapos sa sign-in. May one-time code ang address, kaya huwag itong i-share.

## Saan naka-store ang mga account ko?

Sa `~/.local/share/ccam/profiles/<name>/`, isang folder bawat account. Sa macOS, iniimbak mismo ng Claude Code ang login sa Keychain, sa ilalim ng pangalang hango sa folder na iyon, kaya sa macOS, pinapanatili ng `ccam rename` ang lumang pangalan ng folder.

## Ano ang mangyayari sa dati kong ~/.claude?

Wala, maliban kung patakbuhin mo ang `ccam import` o `ccam migrate`. Nirerehistro ng import ang folder kung nasaan ito, kaya ginagamit pa rin ng karaniwang `claude` at ng `ccam <name>` ang parehong account at settings. Kinokopya ng `ccam migrate` ang settings, `CLAUDE.md`, mga extension at MCP server nito sa shared na group na `default`, at binabasa lang ang folder.

## Aling mga system ang suportado ng ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine at iba pa), macOS at Windows 10 at 11, direkta o sa loob ng WSL. Sa Windows, kailangan ng Developer Mode para mag-share ng mga file gaya ng `settings.json`; gumagana ang mga shared na folder kahit wala ito, at kapag wala ito, ginagawa ng unang `ccam add` ang group na `default` na may mga MCP server at extension lang ([sharing sa Windows](/features/share#on-windows)).

## Bakit sa iisang account lang lumalabas ang isang MCP server?

Ang mga server na idinagdag sa loob ng session gamit ang `claude mcp add` ay naka-store para sa account na iyon lang, at kapag walang `-s user`, para sa kasalukuyang folder lang. Sasabihin ito sa iyo ng ccam sa susunod na simulan mo ang account, at ipinapakita sila ng `ccam mcp list` at `ccam ext`; ilipat ang server sa shared group gamit ang `ccam mcp import <account> <server>`.

## Pinababagal ba ng ccam ang Claude Code?

Nagdadagdag ng mga 4 millisecond ang pagsisimula ng Claude Code sa pamamagitan ng ccam, na sinukat gamit ang `scripts/perf-check.sh` sa isang laptop. Nagbabasa ang ccam ng isang maliit na file, sinusuri ang mga shared na link at saka pinapalitan ang sarili nito ng `claude`, kaya bilang default, hindi ito tumatakbo habang nagtatrabaho ka. Hindi kayang palitan ng Windows ang isang process, kaya doon sinisimulan ng ccam ang Claude Code at hinihintay ito, nang walang ibang ginagawa. Kapag naka-on ang `relay` o `auto_switch`, ganoon din patuloy na tumatakbo ang ccam kasabay ng bawat interactive session, para maipagpatuloy nito ang usapan sa ibang account.

## Paano i-update ang ccam?

Patakbuhin ang `ccam update`. Dina-download nito ang pinakabagong release, sinusuri ang checksum nito at pinapalitan ang `ccam` file; sinasabi lang ng `ccam update --check` kung may mas bagong release. Sinasabi rin sa iyo ng dashboard kapag may mas bagong release, at ini-install ito ng `U` pagkatapos magtanong; hindi kailanman kusang nag-a-update ang ccam. Kung in-install mo ang `.deb` package, i-install na lang ang bagong `.deb` gamit ang apt. Tingnan ang [pag-update](/docs/install#update).

## Paano i-uninstall ang ccam?

Patakbuhin ang `ccam uninstall`. Inililista nito ang mga aalisin at nagtatanong muna, saka inaalis ang bawat account (ina-log out muna), ang mga shortcut command, `~/.config/ccam`, `~/.local/share/ccam`, `~/.local/state/ccam` at ang command na `ccam`. Nananatili ang in-import na `~/.claude`, kasama ang kopya ng shine-share nito sa mga group nito. Tingnan ang [install](/docs/install#uninstall).
