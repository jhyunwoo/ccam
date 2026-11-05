---
title: Mga usage limit at paglipat ng account
description: Tingnan ang 5-oras at lingguhang usage ng bawat Claude Code account gamit ang ccam usage, at buksan ang may pinakamaraming natitira gamit ang ccam next --auto.
updated: 2026-11-05
order: 5
translation: machine
sourceHash: "06bd349073e3"
---

Ipinapakita ng `ccam usage`, para sa bawat account, kung gaano kalaki na ang nagamit sa bawat limit: ang 5-oras na limit, ang lingguhang limit at ang sariling lingguhang limit ng bawat model (gaya ng Fable). Binubuksan ng `ccam next --auto` ang account na may pinakamaraming natitira, at dinadala ng `--continue` ang kasalukuyang usapan.

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

Bawat account ay may linyang nagsasabi kung gaano na katagal ang mga numero nito, at pagkatapos ay isang bar bawat limit kasama ang oras bago ito mag-reset: ang 5-oras na window, ang lingguhang window para sa lahat ng model, at ang sariling lingguhang window ng bawat model, gaya ng Fable. Nagiging dilaw ang bar sa 80% at pula sa 100%. Pinapangalanan ng huling linya ang account na bubuksan ng `ccam next --auto`; kapag `auto_strategy = "expiring"` o `"balance"`, `→ next --auto picks <name>: <reason>` na lang ang nakasulat dito.

**Mga seat ng Team at Enterprise.** Lumilipat din ang awtomatikong pagpili ng account (`ccam next`, `--fallback`, `ccam pick`, ang `n` ng dashboard at, kapag binuksan mo ito, ang [`auto_switch`](/features/next#auto-switch)) sa pagitan ng mga seat ng iisang Claude Team o Enterprise na organisasyon, na maaaring lumabag sa terms ng Anthropic. Basahin ang [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) at magtanong muna sa iyong organisasyon at sa Anthropic bago ka umasa rito.

## Saan galing ang mga numero

Sabay-sabay na nagtatanong ang ccam sa usage endpoint ng Anthropic para sa bawat account, gamit ang sariling login ng bawat account. Naghihintay ito nang hanggang 3 segundo bawat account at ipinapakita ang bawat account sa sandaling sumagot ito.

Hindi pampublikong dokumentado ang endpoint na ito at puwede itong magbago. Kapag pumalya ang isang request, ipinapakita ng ccam ang huling alam na mga numero at minamarkahan ang mga ito bilang cached. Ang window na nag-reset na mula nang i-save ang mga numerong iyon ay nagpapakita ng 0%. Sasabihin sa iyo ng `ccam doctor` kung nagbago ang format ng response.

Kapag pinahinto ng Claude Code ang isang session sa usage limit, itinatala ng [ccam plugin](/docs/install#claude-code-plugin) kung aling window ang puno at kung kailan ito magre-reset. Hanggang doon, 100% ang ipinapakita ng window na iyon anuman ang sabihin ng ibang mga numero, kaya hindi pinipili ulit nang masyadong maaga ng `ccam next --auto` at `--fallback` ang account.

Muling ginagamit ang mga numerong kinuha wala pang 60 segundo ang nakalipas. Mga option:

| Option | Epekto |
| --- | --- |
| `--watch[=seconds]` | Tuloy-tuloy na i-refresh, tuwing 60 segundo bilang default |
| `--json` | Output na nababasa ng makina |
| `-v` | Ipakita rin kung saan galing ang mga numero ng bawat account |
| `--no-cache` | Laging magtanong ulit |
| `--source cache` | Huwag kailanman magtanong; ipakita lang ang mga naka-save na numero |

## Mga expired na login

Kapag nag-expire na ang access token ng isang account, hindi ito nire-renew ng ccam, dahil kapag ni-renew ito habang ginagamit ng isang Claude Code session ang parehong login, puwedeng ma-sign out ang session na iyon. Ipinapakita ng account ang ``access token expired; `ccam <name>` renews it when it starts``: simulan ito nang isang beses gamit ang `ccam <name>` at ire-renew ito ng Claude Code. Ang login na buong nag-expire, o na binura ng Claude Code, ay nagpapakita ng `` login expired → `ccam login <name>` `` at nangangailangan ng bagong login.

Nire-renew mismo ng `ccam usage --refresh` ang mga expired na token, pero para lang sa mga account na walang tumatakbong Claude Code session, paano man sinimulan ang session na iyon. Hindi nito kailanman nire-renew ang mga login na naka-store sa macOS Keychain o sa isang in-import na `~/.claude`, at puwedeng tumagal nang hanggang 30 segundo ang isang pag-renew.

## Usage sa status line

```sh
ccam statusline install
```

Ipapakita na ng Claude Code ang pangalan ng account at ang usage nito sa status line, halimbawa `[work] · Opus · 5h 24% · 7d 41%`. Mula 90% pataas, iminumungkahi rin ng linya ang `ccam next --auto --continue`. Kung may status line command ka na, pinapanatili ito ng ccam at inilalagay sa unahan ang pangalan ng account. Sine-save din ng status line ang mga numero nito, na ginagamit ng `ccam usage` kapag hindi sumasagot ang endpoint.

## Paglipat

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the most room
ccam next --continue  # continue this folder's conversation on the next account
```

Nilalaktawan ng `ccam next` ang mga account na hindi naka-log in, at nilalaktawan ng `--auto` ang mga account na nasa 100% (pati ang window na iniulat ng Claude Code na puno na, hanggang sa reset nito) at ang mga account ng ibang organisasyon. Kapag nagshe-share ng kasaysayan ng usapan ang dalawang account, binubuksan ng `--continue` ang parehong usapan; kung hindi, kinokopya ito ng Claude Code sa bagong account. Tingnan ang [sharing](/docs/sharing) at ang [paglipat sa susunod na account](/features/next#continue).

Kapag may `auto_switch = true` sa itaas ng `config.toml` (Linux at macOS; naka-off bilang default), ang ccam mismo ang gumagawa ng paglipat na ito kapag pinahinto ng Claude Code ang isang interactive session sa limit nito at may natitira ang ibang account na kauri nito (seat sa parehong organisasyon, o ibang personal na plan). Tingnan ang [`auto_switch`](/features/next#auto-switch) para malaman kung kailan ito kumikilos at kung ano ang hindi nito kailanman ginagawa.

Para sa mga script at naka-iskedyul na job na nagpapatakbo ng `claude -p`, ilagay ang `--fallback` bago ang pangalan ng account. Kapag naabot ng account na iyon ang limit nito, pinapatakbo ulit ng ccam ang parehong command bilang account na may pinakamaraming natitira, at ang output lang ng run na natapos ang ipi-print:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Tingnan ang [Patakbuhin ulit sa limit](/features/next#run-again-at-the-limit) para sa kung paano nagpapasya ang ccam.

Ang [gabay sa workflow kapag naabot ang usage limit](/guides/usage-limit-workflow) ay nagpapakita ng isang buong halimbawa.
