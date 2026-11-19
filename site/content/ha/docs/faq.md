---
title: Tambayoyin da ake yawan yi
metaTitle: "Tambayoyi game da ccam: asusun Claude Code da dama"
description: Amsoshi game da ccam, asusun Claude Code, iyakokin amfani, rabawa, sirri da tsarukan da ake tallafawa.
updated: 2026-11-19
order: 6
translation: machine
sourceHash: "090ecc3a34bd"
---

## Shin ccam kayan aikin Anthropic ne na hukuma?

A'a. ccam aikin open-source ne da ba na hukuma ba, ba shi da alaƙa da Anthropic, kuma Anthropic ba ta amince da shi ba. Yana amfani da saitin `CLAUDE_CONFIG_DIR` wanda [Claude Code ya rubuta bayaninsa](https://code.claude.com/docs/en/env-vars) don gudanar da asusai da dama gefe da gefe.

## Shin an yarda in yi amfani da asusun Claude Code da dama?

Wannan ya dogara da sharuɗɗan kowane asusu kuma, ga asusan aiki, da manufofin ƙungiyarka: [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) na Anthropic don Pro da Max, [Commercial Terms](https://www.anthropic.com/legal/commercial-terms) ɗinta don Team da Enterprise, da [Usage Policy](https://www.anthropic.com/legal/aup) ɗinta ga kowa. Ka duba su kafin ka yi amfani da asusai da dama. ccam ba ya ɗaga, haɗawa ko ɓoye iyakokin asusu, kuma har yanzu kowane asusu yana tsayawa a nasa iyakar. Yana motsa aikinka zuwa wani daga cikin asusanka idan ka nema, ko da kansa a iyaka idan ka kunna [`auto_switch`](/features/next#auto-switch) (Linux da macOS). Ka duba ko an yarda da hakan ga asusanka kafin ka yi amfani da shi; kai ne za ka yanke shawara.

**Kujeru (seats) na Team da Enterprise.** Zaɓin asusu na atomatik (`ccam next`, `--fallback`, `ccam pick`, [`ccam agent`](/features/agent), maɓallin `n` na dashboard, da [`auto_switch`](/features/next#auto-switch) idan ka kunna shi) yana kuma sauyawa tsakanin kujerun ƙungiyar Claude Team ko Enterprise iri ɗaya. Cibiyar taimako ta Anthropic ta ce "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (ana ba da kowace kujera a tsarin Team da sunan takamaiman mutum ɗaya, kuma wannan mutumin kaɗai ne aka yarda ya yi amfani da ita) da "Users can't route their work through another user's account to use that user's remaining usage" (masu amfani ba za su iya tura aikinsu ta asusun wani mai amfani ba domin su ci gajiyar ragowar amfanin wannan mai amfanin) ([Tsare-tsaren Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); haka ma ga [tsare-tsaren Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Don haka yin amfani da zaɓin asusu na atomatik tsakanin kujeru na iya saɓa wa sharuɗɗan Anthropic. Ko za ka yi amfani da shi a can ya rage naka: ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi. Don tsare-tsaren kai (Pro, Max), duba sakin layin da ke sama.

## Shin ccam yana kauce wa iyakokin amfani?

Ba ya wuce iyakar kowane asusu: kowane asusu yana riƙe iyakokinsa kuma yana tsayawa a kansu, kuma ccam ba ya ɗaga su, haɗa su ko ɓoye su. ccam yana nuna su kuma yana ba ka damar ci gaba a wani daga cikin asusanka idan iyakar ɗaya ta ƙare; idan `auto_switch` yana kunne, yana yin haka da kansa (duba tambaya ta gaba). Tsakanin kujerun ƙungiyar Team ko Enterprise ɗaya, ci gaba a wata kujera yana amfani da ragowar amfanin wannan kujerar, abin da cibiyar taimako ta Anthropic ta hana ga kujerar wani mai amfani (duba Kujeru na Team da Enterprise a sama).

## Shin ccam zai iya sauya asusu da kansa idan zama ya kai iyakarsa?

Sai idan ka kunna shi. Da `auto_switch = true` a saman `~/.config/ccam/config.toml` (Linux da macOS; a kashe ta tsohuwa kuma na gwaji ne), idan Claude Code ya dakatar da zaman mu'amala (wanda ccam ya fara) a iyakar awa 5, ta mako, ta Opus, ta Sonnet ko ta Fable, ccam yana kawo ƙarshen zaman kuma yana ci gaba da tattaunawar a matsayin wani daga cikin asusanka mai sarari: wata kujerar ƙungiyar Team ko Enterprise iri ɗaya, ko wani tsarin kai (Pro, Max). Ba ya taɓa sauyawa kafin iyakar, ba ya taɓa sa ido kan amfaninka a bayan fage, kuma ba ya taɓa fitar da aiki daga ƙungiyar Team ko Enterprise ko shigar da shi cikinta.

A kan kujeru yana motsa aikinka daga wata kujera zuwa wata, abin da na iya saɓa wa sharuɗɗan Anthropic (duba Kujeru na Team da Enterprise a sama). Ga tsare-tsaren kai, sharuɗɗan Anthropic ba su yi magana a fili kan batun ba: Consumer Terms sun hana samun damar Services "through automated or non-human means" (ta hanyoyin atomatik ko waɗanda ba na ɗan Adam ba) sai dai idan Anthropic ta ba da izini a fili, kuma motsi yana aika saƙo ɗaya a madadinka. Kunna shi ya rage naka.

[`auto_switch`](/features/next#auto-switch) yana lissafa kowane sharaɗi da abin da motsi ke ci.

## Ta yaya zan yi amfani da asusun Claude Code biyu a kwamfuta ɗaya?

Ba kowane asusu foldar saitunsa ta kansa da `CLAUDE_CONFIG_DIR`. ccam yana yi maka wannan: gudanar da `ccam add` ga kowane asusu, sannan ka fara su da `ccam personal` da `ccam work`, ko da gajerun umarni da kai da kanka ka sa musu suna, kamar `ccwork`. [Jagorar](/guides/multiple-claude-code-accounts) kuma tana nuna yadda ake yi da hannu.

## Ta yaya zan ci gaba da tattaunawa a wani asusu?

Gudanar da `ccam next --continue` a cikin foldar: yana kai tattaunawar ƙarshe ta foldar zuwa asusu na gaba (`--auto` don wanda ya fi sarari). A cikin Claude Code, `/ccam:switch` daga [plugin ɗin ccam](/docs/install#claude-code-plugin) yana yin haka ga zaman da kake ciki. Idan asusan biyu suna raba tarihin tattaunawa (`ccam group set default --add history`, ko rukuni nasu na kansu: `ccam group add pair --items history --members work,side`), sabon asusun yana buɗe ainihin wannan tattaunawar; in ba haka ba Claude Code yana kwafar ta cikin tarihin sabon asusun. Ƙara `--brief` don fara sabon asusun a kan taƙaitaccen bayani maimakon dukan tattaunawar.

## Shin ccam yana ganin bayanan shigata ko yana aika su?

ccam ba ya taɓa raba ko kwafar bayanan shigar Claude naka tsakanin asusu, kuma ba ya taɓa kwafar shigar sabobin MCP. Yana karanta access token na asusu ne kawai don tambayar endpoint na amfani na Anthropic game da amfanin wannan asusun: don `ccam usage`, dashboard, `ccam doctor` da zaɓin atomatik (`ccam next --auto`, `ccam pick`, `--fallback`, da `auto_switch` idan zama ya kai iyakarsa). Yana sabunta wanda ya ƙare ne kawai da `ccam usage --refresh`; da `usage_source = "claude"` ba ya aika token zuwa ko'ina, kuma yana karanta shigar da aka ajiye (fayil ɗin, ko shigarwar Keychain ta macOS) ne kawai don sanin ko har yanzu tana aiki. `ccam ext` da `ccam doctor` suna karanta sunayen shigar MCP na asusu, ba token ɗinsu ba. ccam ba shi da telemetry. Ban da waɗannan tambayoyin, buƙata ɗaya tak da yake yi da kansa ita ce binciken da dashboard ke yi a GitHub don sabon sakin ccam, wadda ba ta ɗauke da komai game da kai ko asusunka.

## Ta yaya zan shiga a sabar nesa?

Gudanar da `ccam add` ko `ccam login <name>` a cikin zaman SSH kamar yadda aka saba. Claude Code yana jiran shigarwar a kan sabar, don haka burauza a kwamfutarka yana ƙarewa a shafin `http://localhost:<port>/callback?code=…` wanda ba ya buɗewa. Kwafi wannan adireshin daga burauza ka liƙa shi a terminal: ccam yana miƙa shi ga Claude Code a kan sabar, wanda ke kammala shigarwar. Adireshin yana ɗauke da lambar amfani sau ɗaya, don haka kada ka raba shi.

## A ina ake ajiye asusuna?

A `~/.local/share/ccam/profiles/<name>/`, folda ɗaya ga kowane asusu. A macOS, Claude Code yana ajiye bayanan shigar kansu a Keychain, ƙarƙashin sunan da aka samo daga wannan foldar, don haka a can `ccam rename` yana riƙe tsohon sunan foldar.

## Me zai faru da ~/.claude da nake da shi?

Babu komai, sai idan ka gudanar da `ccam import` ko `ccam migrate`. Import yana yin rajistar foldar a inda take, don haka `claude` na yau da kullum da `ccam <name>` suna ci gaba da amfani da asusu da saituna iri ɗaya. `ccam migrate` yana kwafar saitunanta, `CLAUDE.md`, extensions da sabobin MCP zuwa rukunin `default` da aka raba, kuma yana karanta foldar ne kawai.

## Waɗanne tsaruka ne ccam yake tallafawa?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine da sauransu), macOS da Windows 10 da 11, kai tsaye ko a cikin WSL. A Windows, raba fayiloli kamar `settings.json` yana buƙatar Developer Mode; foldojin da aka raba suna aiki ba tare da shi ba, kuma ba tare da shi ba `ccam add` na farko yana ƙirƙirar rukunin `default` da sabobin MCP da extensions kawai ([rabawa a Windows](/features/share#on-windows)).

## Me yasa sabar MCP take bayyana a asusu ɗaya kawai?

Sabobin da aka ƙara a cikin zama (session) da `claude mcp add` ana ajiye su ga wannan asusun kawai, kuma ba tare da `-s user` ba ga foldar yanzu kawai. ccam yana gaya maka a gaba da ka fara asusun, kuma `ccam mcp list` da `ccam ext` suna nuna su; motsa sabar zuwa rukunin da aka raba da `ccam mcp import <account> <server>`.

## Shin ccam yana rage saurin Claude Code?

Fara Claude Code ta ccam yana ƙara kusan millisecond 4, wanda aka auna da `scripts/perf-check.sh` a kwamfutar tafi-da-gidanka. ccam yana karanta ƙaramin fayil ɗaya, yana duba mahaɗan da aka raba, sannan yana maye gurbin kansa da `claude`, don haka ta tsohuwa ba ya aiki yayin da kake aiki. Windows ba zai iya maye gurbin process ba, don haka a can ccam yana fara Claude Code kuma yana jiran sa, ba tare da yin wani abu ba. Idan an kunna `relay` ko `auto_switch`, ccam yana ci gaba da aiki ta wannan hanya tare da kowane zaman mu'amala, don ya iya ci gaba da tattaunawar a wani asusu.

## Ta yaya zan sabunta ccam?

Gudanar da `ccam update`. Yana sauke sakin baya-bayan nan, yana duba checksum nasa, sannan yana maye gurbin fayil ɗin `ccam`; `ccam update --check` yana gaya maka kawai ko akwai sabon saki. Dashboard ma yana gaya maka idan akwai sabon saki, kuma `U` yana girka shi bayan ya tambaya; ccam ba ya taɓa sabunta kansa da kansa. Idan ka girka kunshin `.deb`, ka girka sabon `.deb` da apt a maimakon haka. Duba [sabuntawa](/docs/install#update).

## Ta yaya zan cire ccam?

Gudanar da `ccam uninstall`. Yana lissafa abin da za a cire kuma yana tambaya tukuna, sannan yana cire kowane asusu (bayan fitar da shi daga shiga), gajerun umarnin, `~/.config/ccam`, `~/.local/share/ccam`, `~/.local/state/ccam` da umarnin `ccam`. `~/.claude` da aka shigo da shi yana nan, tare da kwafin abin da ya raba a rukunoninsa. Duba [shigarwa](/docs/install#uninstall).
