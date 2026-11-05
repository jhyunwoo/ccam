---
title: Kapag naabot mo ang Claude Code usage limit
description: Maghintay ng reset o ituloy agad ang parehong usapan sa ibang Claude Code account. Paano gumagana ang 5-oras at lingguhang limit, hakbang-hakbang.
updated: 2026-11-05
order: 3
translation: machine
sourceHash: "53b1e3634ea9"
---

Kapag naabot ng isang account ang usage limit nito (madalas tawaging rate limit), may dalawa kang pagpipilian: maghintay hanggang mag-reset ang limit, o ituloy ngayon din ang parehong usapan sa isa pang account na pag-aari mo. Sa ccam, isang command lang ang pangalawa, `ccam next --auto --continue`.

## Anong mga limit mayroon ang Claude Code?

Sa mga Pro at Max plan, may [5-oras na session limit at lingguhang limit](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan) ang Claude Code. Parehong kahati ng mga ito ang Claude sa web at sa mga app, at sabay na ibinabawas ang usage sa dalawa. Kapag naubos ang isa, [nagpapakita ng mensahe](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) ang Claude Code gaya ng "You've hit your session limit · resets 3:45pm". Hindi makakatulong ang paglipat ng model, dahil saklaw ng mga limit ang lahat ng model.

## Opsyon 1: hintayin ang reset

Kung iisa lang ang account mo, maghintay. Ang mga bagong bersyon ng Claude Code ay [naghihintay sa bukas na session](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) at kusang itinutuloy ang gawain kapag nag-reset na ang limit. Ipinapakita ng `/usage` sa loob ng Claude Code kung gaano na karami ang nagamit mo sa mga limit ng iyong plan.

## Opsyon 2: ituloy ngayon sa ibang account

Kung may pangalawa kang account, puwede mong ituloy agad ang parehong usapan. I-set up ito nang isang beses:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Ipinapakita na ngayon ng status line ang account at ang 5-oras na usage nito habang nagtatrabaho ka, kaya hindi ka mabibigla sa limit.

Kapag naabot ang limit, lumabas sa Claude Code (`/exit`), pagkatapos:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Buksan ang usapan sa account na may pinakamaraming natitira:

```sh
ccam next --auto --continue
```

Magsisimula ang Claude Code bilang `work` kasama ang pinakahuling usapan ng kasalukuyang folder, kasama ang lahat ng nasabi na. Kapag shared ang kasaysayan (ang `ccam group set default --add history` sa itaas), iisang usapan ang ginagamit ng dalawang account; kung hindi, kinokopya ng Claude Code ang usapan sa `work`, at sinasabi ng ccam kung ilang token ang kailangang basahin ulit ng unang sagot.

Kapag nasa Claude Code ang [ccam plugin](/docs/install#claude-code-plugin), nagpapakita ng notification ang session na huminto sa limit nito at itinatala nito ang punong window, kaya hindi pinipili ulit ng `ccam next --auto` ang account na iyon bago ang reset nito. Itinatala ng `/ccam:switch` (o `/ccam:switch work`) ang usapan para sa `ccam next --continue`, kaya eksaktong ang usapang iyon ang ipinagpapatuloy nito.

**Mga seat ng Team at Enterprise.** Lumilipat din ang awtomatikong pagpili ng account (`ccam next`, `--fallback`, `ccam pick`, ang `n` ng dashboard at, kapag binuksan mo ito, ang [`auto_switch`](/features/next#auto-switch)) sa pagitan ng mga seat ng iisang Claude Team o Enterprise na organisasyon, na maaaring lumabag sa terms ng Anthropic. Basahin ang [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) at magtanong muna sa iyong organisasyon at sa Anthropic bago ka umasa rito.

## Kusa bang nagpapalit ng account ang ccam?

Kapag binuksan mo lang ito.

- Bilang default, hindi tumatakbo ang ccam habang nagtatrabaho ka: sinisimulan nito ang Claude Code at pinapalitan ang sarili nito ng Claude Code, kaya ikaw ang nagpapasya kung kailan lilipat. Pumipili ang `ccam next --auto` ayon sa usage, kinukuha ang may pinakakaunting lingguhang usage sa mga account na nasa loob ng 3 puntos ng pinakamababang 5-oras na usage; nilalaktawan nito ang mga account na nasa 100% o hindi naka-log in, at nananatili sa iisang Team o Enterprise na organisasyon, o sa mga personal na plan.
- Kapag may `relay = true` sa `config.toml`, patuloy namang tumatakbo ang ccam kasabay ng session, kaya kusang ipinagpapatuloy ng `/ccam:switch` ang usapan, at pinapanatili nito ang mga option, permission mode, effort level at agent ng session; bago nito ilipat ang usapan sa account ng ibang organisasyon, nagtatanong muna sa iyo ang ccam.
- Kapag may `auto_switch = true` (Linux at macOS; eksperimental), ang ccam mismo rin ang lumilipat: kapag pinahinto ng Claude Code ang session sa limit nito at may natitira ang ibang account na kauri nito, tinatapos ng ccam ang session at ipinagpapatuloy doon ang usapan ([mga kondisyon at kapalit](/features/next#auto-switch)).

Ang ilang ibang tool ay nagpapalit ng login ng tumatakbong session; ipinapaliwanag ng [ccam at iba pang account switcher](/compare/account-switchers) ang pagkakaiba.

## Mga dapat tandaan

- Patakbuhin ang `ccam next --continue` sa parehong folder kung saan ka nagtrabaho; naka-store ang mga usapan ayon sa project folder.
- Iwasang magpatakbo ng dalawang account sa parehong usapan nang sabay. Babalaan ka ng ccam kung susubukan mo.
- Kapag hindi shared ang kasaysayan, kinokopya ang usapan sa kabilang account at magkakahiwalay na ang landas ng dalawang kopya; nananatili sa unang account ang mga checkpoint na pinagbabalikan ng code ng `/rewind`.
- Puwedeng mas malaki ang mabawas sa limit ng bagong account sa unang mensahe kaysa karaniwan, dahil hindi pa naka-cache sa account na iyon ang usapan. Sa halip, sinisimulan ito ng `ccam next --continue --brief` sa maikling buod.
- Nakadepende sa terms na sumasaklaw sa bawat account kung puwede kang lumipat sa pagitan ng mga account nang ganito.

Tingnan ang [Mga usage limit](/features/usage) at [Paglipat sa susunod na account](/features/next) para sa bawat option.
