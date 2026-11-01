---
title: Iyakokin amfani
description: ccam usage yana nuna iyakokin awa 5 da na mako na kowane asusu a lokaci ɗaya, kuma ccam statusline yana saka sunan asusun da amfaninsa a Claude Code.
updated: 2026-11-01
order: 7
translation: machine
sourceHash: 54b296f7da35
---

`ccam usage` yana nuna, ga kowane asusu, nawa aka yi amfani da shi daga kowace iyaka: iyakar awa 5, iyakar mako da iyakar mako ta kowane model (kamar Fable). Kowace iyaka tana samun sanda da lokacin da ya rage kafin ta sake saitawa. Yi amfani da shi kafin ka fara dogon aiki, ko idan Claude Code ya gaya maka iyaka ta kusa, don ganin asusun da ke da sauran sarari. `ccam statusline` shi ne umarnin layin matsayi (status line) na Claude Code: yana nuna sunan asusun da amfaninsa a cikin Claude Code kuma yana ajiye lambobin don `ccam usage`. Don token da ƙiyasin kuɗin da ke bayan amfanin, `ccam ccusage` yana gudanar da [ccusage](https://ccusage.com) a kan tattaunawar kowane asusu ([a ƙasa](#tokens-and-costs-ccam-ccusage)).

## Amfani

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts that are not running
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

Kowane asusu yana samun layi mai lambarsa, sunansa da tsawon lokacin da lambobinsa suka yi; `(cached)` yana nufin sun fito daga fayil da aka ajiye, ba daga endpoint ba. A ƙarƙashinsa, layi ɗaya ga kowace iyaka yana nuna sanda mai ƙofofi 20, kashi da lokacin da ya rage kafin taga ta sake saitawa: `5-hour`, `Weekly` (dukkan model), sannan taga ta mako ta kowane model, kamar `Fable`. Tagar da asusun ba shi da ita tana nuna `-`. Kashi yana zama rawaya a 80% kuma ja a 100%.

Idan Claude Code ya ƙi buƙata saboda wata taga ta cika, wannan tagar tana nuna 100% har sai ta sake saitawa, ko me lambobin da aka ajiye suka ce (duba [rikodin iyaka](#the-limit-record)). A terminal kowane asusu yana bayyana nan take a matsayin `fetching…`, kuma ana zana layukansa yayin da ya amsa.

Layin ƙarshe, wanda ke bayyana ne kawai idan kana da asusu fiye da ɗaya, yana ambaton asusun da [`ccam next --auto`](/features/next) zai zaɓa. Yana bin ƙa'idodi iri ɗaya:

- `auto_max_usage` da `auto_strategy`;
- zaman da ke gudana;
- asusan ƙungiyar asusun da aka yi amfani da shi na ƙarshe kawai, sai dai idan an saita `cross_org`;
- ana tsallake duk tagar da ta kai 100%, don haka asusun da makonsa ya cika ba a taɓa ambaton sa ba.

Da tsohuwar dabarar `headroom` yana cewa `→ Most headroom: <name>`; da `expiring` ko `balance` yana ba da dalilin, kamar `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Kujeru (seats) na Team da Enterprise.** Zaɓin asusu na atomatik (`ccam next`, `--fallback`, `ccam pick`, maɓallin `n` na dashboard kuma, da zarar ka kunna shi, [`auto_switch`](/features/next#auto-switch)) yana kuma sauyawa tsakanin kujerun ƙungiyar Claude Team ko Enterprise iri ɗaya, wanda na iya saɓa wa sharuɗɗan Anthropic. Ka karanta [Tambayoyin da ake yawan yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) kuma ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi.


Asusun da ba a shiga da shi ba yana nuna `not logged in → ccam login <name>`, kuma asusun API key yana nuna `no limits (API)`. Access token da ya ƙare yana nuna ``access token expired; `ccam <name>` renews it when it starts``, a saman lambobin da aka ajiye, idan akwai: Claude Code yana sabunta shi a farawa ta gaba. Shigar da ta ƙare gaba ɗaya, domin refresh token ɗinta ya ƙare ko Claude Code ya share ta bayan sabar ta ƙi ta, tana nuna `` login expired → `ccam login <name>` ``.

Bayan `ccam statusline install`, Claude Code yana nuna layi kamar `[work] · Opus · 5h 24% · 7d 41%`.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--json` | Buga rahotannin a matsayin JSON: `name`, `slot`, `status`, `limits` (tare da `credits` idan asusun yana da credits na amfani a kunne), `source`, `at`, `error`, da `limited` da `limited_windows` idan rikodin iyaka ya saita tagogi zuwa 100%. Idan babu asusai yana buga `[]` | a kashe |
| `--watch[=N]` | Sake zanawa kowane daƙiƙa `N` har sai ka danna Ctrl+C. `--watch` shi kaɗai, ko ƙima ta 0 ko ƙasa da haka, yana nufin 60. Ƙimar tana zuwa bayan `=`: ana ƙin `--watch 10` tare da shawara | a kashe |
| `-v`, `--verbose` | Ƙara layin `source:` ƙarƙashin kowane asusu: inda lambobinsa suka fito (endpoint na amfani, Claude Code, layin matsayi, fayil da aka ajiye, ko iyakar amfani da Claude Code ya ba da rahoto). Idan rikodin iyaka yana rinjayar lambobin, layi kamar `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` yana biyo baya | a kashe |
| `--no-cache` | Yi watsi da cache na daƙiƙa 60 kuma sake tambayar endpoint | a kashe |
| `--source live\|claude\|cache` | `live` yana tambayar endpoint da token na asusun; `claude` yana sa Claude Code ya nemo shi (a ƙasa); `cache` yana nuna lambobin da aka ajiye kawai | `usage_source`, in ba haka ba `live` |
| `--refresh` | Sabunta access tokens da suka ƙare, ga asusan da ba sa aiki kawai. Sabuntawa na iya ɗaukar har daƙiƙa 30 | a kashe |
| `--dry-run` | Ana karɓar sa: `ccam usage` ba ya canza komai sai cache nasa. Ana ƙin sa tare da `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | a kashe |

Duk wata ƙimar `--source` dabam tana gazawa da `--source must be live, claude or cache`. `ccam usage` ba ya karɓar arguments.

### Amfani daga Claude Code kansa

Da `usage_source = "claude"` a saman `config.toml` (kafin `[[group]]`, `[[profile]]` ko `[[folder]]` na farko), ko `--source claude`, ccam ba ya karanta token na shiga ko kaɗan. Ga kowane asusu yana fara `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` a foldar asusun kuma yana aika buƙatar `get_usage` ta Claude Code a standard input (Claude Code 2.1.270 ko sabo). Claude Code yana nemo amfanin da nasa shigar kuma yana amsawa da tagogin awa 5, na mako da na model.

`--safe-mode` yana hana hooks, sabobin MCP da plugins na asusun farawa, kuma ba a ajiye zama. Kowane asusu yana buƙatar farawar Claude Code ɗaya (kusan daƙiƙa ɗaya), don haka wannan tushe ya fi `live` jinkiri. Ana tambayar asusan a lokaci ɗaya, kowanne da iyakar lokaci ta daƙiƙa 20, kuma cache na daƙiƙa 60 yana aiki kamar yadda aka saba.

Shigar da access token ɗinta ya ƙare, ko zai ƙare cikin minti 5, tana samun daƙiƙa 50, tunda Claude Code yana sabunta shi yayin da yake farawa kuma sabuntawar da aka yanke na iya rasa shigar. Idan ccam ya hakura da wani neman, yana roƙon Claude Code ya fita da SIGTERM, don Claude Code ya kammala sabuntawar da ya fara (Windows ba shi da SIGTERM, don haka a can ana dakatar da process ɗin).

Shigar da ta ƙare gaba ɗaya tana nuna `` login expired → `ccam login <name>` `` ba tare da fara Claude Code ba. Claude Code yana yi wa buƙatar alama a matsayin ta gwaji, don haka idan ta daina amsawa, ccam yana nuna lambobin da aka ajiye kuma `ccam doctor` yana bayar da rahoto. Idan Claude Code ya amsa cewa nasa neman ya gaza, ccam yana nuna lambobin da aka ajiye tare da `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`, ba tsari mara iyaka ba.

Da `--dry-run`, ccam ba ya fara Claude Code, domin Claude Code yana rubutu a foldar asusun; yana nuna lambobin da aka ajiye maimakon haka. `usage_source` yana karɓar `live` ko `claude`; `cache` zaɓi ne na `--source` kawai.

`ccam statusline` yana karanta JSON na layin matsayi na Claude Code daga standard input. Zaɓinsa ɗaya shi ne `--wrap <command>`: umarnin layin matsayi da ke akwai wanda layin fitarwa na farko ke samun sunan asusun a gabansa. `ccam statusline install` ba ya karɓar zaɓuɓɓuka. Duka biyun suna ƙin `--dry-run` da `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Token da kuɗi: `ccam ccusage`

`ccam usage` yana nuna nawa asusu ya yi amfani da shi daga kowace iyaka. [ccusage](https://ccusage.com) kuwa yana ƙirga token da ƙiyasin kuɗi a cikin bayanan tattaunawar da Claude Code ke ajiyewa. Shi kaɗai, yana neman bayanan da Claude Code ke ajiyewa a `~/.claude` da `~/.config/claude` kawai (ko a foldojin da `CLAUDE_CONFIG_DIR` ke ambato), don haka asusun ccam ɗaya tilo da zai iya gani shi ne wanda aka shigo da shi a wurinsa a `~/.claude`. `ccam ccusage` yana gudanar da shi da `CLAUDE_CONFIG_DIR` da aka saita zuwa foldojin kowane asusu da ke da tattaunawa, don rahoto ɗaya ya ƙunshi su duka:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam yana fara karanta zaɓuɓɓukansa: `-a` (`--account`) yana karɓar sunaye ko lambobin asusu, a raba su da waƙafi ko a maimaita zaɓin, kuma `--dry-run` yana buga umarnin maimakon gudanar da shi. Daga argument na farko da ba zaɓin ccam ba, komai yana zuwa ga ccusage yadda ka rubuta shi, sai dai `--config` ga `bun x` da `npx` (a ƙasa); `ccam ccusage -- --help` yana nuna taimakon ccusage kansa. ccam yana gudanar da `ccusage` idan an girka shi, in ba haka ba `bun x ccusage`, in kuma hakan bai yiwu ba `npx ccusage@latest`; waɗannan biyun na ƙarshe suna sauke shi daga npm. Suna farawa a foldar ccam kansa (`~/.local/state/ccam/ccusage`), ba a foldar da kake ciki ba, don haka wani project da ke can ba zai iya zaɓar abin da suke gudanarwa ta `node_modules/.bin` nasa, `bin` a cikin `package.json` nasa ko registry a cikin `.npmrc` nasa ba; `ccam --dry-run ccusage` yana nuna `cd` zuwa cikin wannan foldar. Idan babu ko ɗaya daga cikinsu, yana gaya maka ka girka Node.js, Bun ko ccusage. A Linux da macOS ccusage yana maye gurbin ccam; a Windows ccam yana jiran sa kuma yana fita da exit code na ccusage. A cikin WSL, ccam yana tsallake `ccusage`, `bun` ko `npx` da ke kan faifan Windows, domin zai gudana a matsayin shirin Windows ba tare da foldojin asusan ba.

- Rahoton yana ƙunshe da asusan da foldojinsu ke da tattaunawa (foldar `projects`). Tattaunawar da asusai ke raba ta wani rukuni ana ƙirga ta sau ɗaya kawai.
- Da `-a`, asusun da ke raba tarihinsa ta wani rukuni yana kawo tattaunawar rukunin gaba ɗaya, kuma ccam yana faɗin na waɗanne asusai ne: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- `blocks` yana nuna tagogin awa 5, kuma kowane asusu yana da nasa, don haka ka ambaci asusu ɗaya kawai a gare shi, kamar yadda yake a sama.
- ccusage yana kuma ba da rahoton sauran kayan aikin coding da ya samu, kamar Codex, sai dai idan ka ambaci ɗaya: `ccam ccusage claude daily` na Claude Code ne kawai.
- Yayin da ccam ke gudanar da shi, ccusage yana karanta saitunansa (`ccusage.json`) daga `./.ccusage` da foldojin asusan kawai. `--config` yana ambaton wani fayil dabam: `ccam ccusage daily --config ~/.config/claude/ccusage.json`. Tun da `bun x` da `npx` suna farawa a foldar ccam kansa, ccam yana ba su `./.ccusage/ccusage.json` da `--config` idan wannan fayil yana nan kuma ba ka ambaci ko ɗaya ba, kuma yana ba da `--config` mai dangantaka (relative) a matsayin cikakkiyar hanya. Sannan ccusage yana karanta sauran hanyoyi masu dangantaka, kamar `--pi-path`, daga foldar ccam, don haka ka ba da su cikakku.
- ccusage zai karanta foldar da hanyarta ke da waƙafi a matsayin foldoji biyu, don haka ccam yana ƙin ta kuma yana ambaton asusun; yi amfani da `-a` don barin wannan asusun a waje.
- `ccusage` umarnin ccam ne, don haka asusun da sunansa ccusage yana farawa da `ccam run ccusage`; `ccam doctor` yana faɗin haka.

## Yadda yake aiki

Lambobin `ccam usage` suna fitowa daga wurare biyu: endpoint na amfani na Anthropic da ba a rubuta bayaninsa ba, `GET https://api.anthropic.com/api/oauth/usage`, wanda ake kira da shigar kowane asusu, da ƙimomin da Claude Code ke miƙa wa umarnin layin matsayinsa. Duk wannan yana cikin `internal/usage`.

1. ccam yana loda `config.toml` kuma yana gina `usage.Account` ɗaya ga kowane asusu (`cli.usageCmd`): foldarsa `~/.local/share/ccam/profiles/<name>/`, alamomin zamansa (tare da bayanan Claude Code kansa a foldar asusun, suna faɗin ko yana aiki) da ko an shigo da shi a wurinsa.
2. ccam yana neman kowane asusu a lokaci ɗaya, goroutine ɗaya ga kowane asusu (`Fetcher.FetchAll`). Dukkansu suna raba `http.Client` ɗaya, don haka ana sake amfani da haɗi zuwa endpoint ɗin. Haɗawa na iya ɗaukar 2 s, kuma kowane asusu yana samun 3 s gaba ɗaya. Client ɗin ba ya bin turawa (redirects), tunda turawa za ta sake aika token zuwa wani wuri, kuma yana karanta 1 MB na amsar a mafi yawa.
3. Yayin da kowane asusu ya kammala, ana sake zana layinsa. Asusun da ya gaza ko bai amsa a kan lokaci ba yana komawa ga lambobin da aka ajiye, don haka asusu ɗaya mai jinkiri ba ya taɓa riƙe sauran.

A zanen, ana neman kowane asusu a goroutine nasa, kuma wanda lokacinsa ya ƙare yana nuna lambobinsa da aka ajiye maimakon haka:

<!-- motion:usage-fetch -->

`personal` da `work` suna amsawa a kan lokaci kuma suna nuna lambobi kai tsaye. `side` yana ɗaukar fiye da 3 s, don haka ccam yana nuna lambobinsa da aka ajiye maimakon haka (`fallback()`). Sannan ccam yana zana toshe ɗaya ga kowane asusu (`usage.Table()`) kuma yana zaɓar asusun da aka ambata a layin ƙarshe (`usage.Best()`).

Ga asusu ɗaya, ccam yana yin waɗannan dubawa a tsari (`Fetcher.Fetch`):

<!-- motion:usage-cache -->

1. Da `--source cache`, yana zuwa `fallback()` kai tsaye.
2. Sai dai idan an saita `--no-cache`, sakamako kai tsaye a `~/.local/state/ccam/usage/<name>.json` da aka auna ƙasa da daƙiƙa 60 da suka wuce (`CacheTTL`) ana amfani da shi, ana auna shi bisa yanzu: tagar da ta sake saitawa tun daga lokacin tana nuna 0%, kuma fayil mai kwanan wata da ya riga agogo ba a ƙirga shi a matsayin sabo. Gudanar da `ccam usage` sau biyu a jere yana kashe buƙata ɗaya ga kowane asusu.
3. ccam yana karanta `.credentials.json` (ko Keychain na macOS) ba tare da canza shi ba kuma maɓallin `claudeAiOauth` kaɗai yake karantawa (`readCredentials()`). Rashin shigar biyan kuɗi (subscription) yana nufin `not logged in`, ko `no limits (API)` idan `.claude.json` yana riƙe API key. Shigar da refresh token ɗinta ya ƙare, ko wadda Claude Code ya share, ita ce `` login expired → `ccam login <name>` `` kuma ba a neman ta.
4. Ba a amfani da access token da ya ƙare. Ba tare da `--refresh` ba, asusun yana nuna lambobinsa da aka ajiye tare da ``access token expired; `ccam <name>` renews it when it starts``.
5. In ba haka ba ccam yana kiran endpoint ɗin. 401 ko 403 ana ɗaukar su a matsayin shigar da ta ƙare (`login expired: authentication failed (401)`), duk wani matsayi dabam ko jikin amsa da ba zai iya karantawa ba ana ɗaukar sa a matsayin kuskure, kuma duka biyun suna komawa ga lambobin da aka ajiye. Neman da Ctrl+C ya ƙare yana cewa `cancelled`. 429 yana nufin endpoint ɗin ya karɓi buƙatu da yawa ga wannan asusun: ccam yana rubuta `<name>.pause` da lokacin sake tambaya (`Retry-After`, ko minti 5, ba fiye da awa ɗaya ba), yana nuna lambobin da aka ajiye tare da `rate limited by the usage endpoint; asking again after 15:10`, kuma ba ya kiran endpoint ɗin ga wannan asusun har zuwa lokacin, ko da `--no-cache`. Sake gwadawa nan take zai tsawaita iyakar saurin ne kawai.
6. Ana rubuta amsa mai kyau a `<name>.json` tare da lokacin da aka auna ta.

Don lambobin da aka ajiye, ccam yana ɗaukar neman da aka ajiye mafi sabo (`<name>.json`, ko `<name>.claude.json` daga neman Claude Code) kuma, idan an rubuta `<name>.statusline.json` bayan sa, yana ɗora tagogin awa 5 da na mako na layin matsayi a kansa (`fallback()`); tagogin model, waɗanda layin matsayi ba shi da su, suna zama yadda neman ya ajiye su. Kafin ya nuna lambobin da aka ajiye, yana duba lokacin sake saitin kowace taga: tagar da ta sake saitawa tun lokacin da aka ajiye lambobin tana nuna 0%, domin tsohon kashi ba ya aiki kuma. Shi ya sa `personal` a sama yana nuna `0%` ba tare da lokacin sake saiti ba.

Endpoint ɗin yana da tsarin amsa biyu. Sabon yana da array na `limits` mai shigarwar `session` (awa 5), `weekly_all` da `weekly_scoped` (na kowane model). Tsohon yana da abubuwan `five_hour`, `seven_day`, `seven_day_opus` da `seven_day_sonnet`. ccam yana karanta duka biyun (`ParseLive`): array na `limits` ne ke nasara, kuma tsofaffin abubuwan suna cike abin da bai ƙunsa ba. Jikin amsa da ba shi da iyaka ko ɗaya ana bayar da rahotonsa a matsayin tsarin da ya canza, wanda `ccam doctor` ke gaya maka game da shi.

### Rikodin iyaka

Idan Claude Code ya dakatar da juyi a iyakar amfani, hook na [plugin ɗin ccam](/docs/install#claude-code-plugin) yana karanta buƙatar da aka ƙi daga tattaunawar: wace taga ce ta cika (awa 5, mako, Opus ko Sonnet) da lokacin da za ta sake saitawa. Yana rubuta hakan a `~/.local/state/ccam/usage/<name>.limit.json`. ccam yana ɗora rikodin a kan duk abin da ya sani, daga kowane tushe (`Fetch`): kowace taga a cikinsa ana ƙirga ta a matsayin 100% har zuwa sake saitinta, ko yaya tsufa ko ƙarancin lambobin da aka ajiye suke, don haka `ccam next --auto`, `ccam pick`, `--fallback` da [`auto_switch`](/features/next#auto-switch) ba sa sake zaɓar asusun kafin lokacin.

Idan ba a san wani abu dabam ba, rikodin kaɗai ne abin da asusun ke nunawa (`source: the usage limit Claude Code reported` da `-v`). Zaman da babu mai kula da su (gudun `-p` da zaman bayan fage) ba sa rubuta rikodi, haka ma dakatarwar da ƙungiya ta saita (misali, idan amfaninta ya ƙare ko iyakar wani rukuni $0 ce), wadda admins ɗinta za su iya ɗagewa a kowane lokaci. `ccam login` cikin wani asusun Claude dabam yana cire rikodin; goge fayil ɗin ma yana share shi.

### Sabunta token

Refresh tokens suna juyawa: sabunta ɗaya yana ba ka sabon refresh token kuma yana soke tsohon. Idan ccam ya sabunta shigar yayin da Claude Code ke gudana a matsayin asusu ɗaya, za a fitar da wannan zaman. Don haka `ccam usage` yana karanta token kuma ba ya taɓa rubuta su, sai dai idan ka ba da `--refresh`. Sannan ccam (`Fetcher.refresh`):

1. yana ƙin asusan da aka yi wa rajista a wurinsu da `ccam import`, tunda `claude` na yau da kullum zai iya amfani da wannan shigar ba tare da ccam ya sani ba;
2. yana ƙin shigar da aka ajiye a Keychain na macOS;
3. yana ɗaukar fayil ɗin kulle, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, don haka processes biyu na ccam ba sa taɓa sabunta token ɗaya dukansu; yayin da ake riƙe da shi, asusun yana nuna `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)`. `ccam rename`, `remove` da `logout` su ma suna ɗaukar sa, kuma suna ƙi yayin da sabuntawa ke kan aiki;
4. ƙarƙashin kullen, yana ƙi idan zaman asusun yana gudana, ko ta yaya aka fara shi (ta ccam, `claude` na yau da kullum, IDE, ko zaman bayan fage da daemon da ke bayansa, wanda ke sabunta shigar da kansa), kuma yana tsallake buƙatar idan wani process ya riga ya sabunta token ɗin;
5. ba ya aika komai yayin da Claude Code ke riƙe da kullen sabuntawarsa a foldar asusun (`.oauth_refresh.lock`, `<folder>.lock`): Claude Code yana sabunta wannan shigar a daidai lokacin;
6. yana aika refresh token zuwa endpoint na token kuma yana jiran amsa har zuwa daƙiƙa 30, lokacin da Claude Code ke ba sabuntawarsa. Da zarar an aika, ba a taɓa katse buƙata, ko da Ctrl+C: sabar ta iya riga ta juya refresh token ɗin, kuma amsarta kaɗai ce ke ɗauke da sabuwar shigar. Ctrl+C ko hangup yana kawo ƙarshen sauran neman, kuma ccam yana fita (130 bayan Ctrl+C) da zarar an ajiye amsar;
7. yana ajiyewa ƙarƙashin kullen rubutu na Claude Code (`.storage-write.lock`) kuma yana canza filayen token na `.credentials.json` kawai kamar yadda yake a lokacin: abin da Claude Code ya rubuta a wannan lokacin (shigar MCP) yana nan, shigar da Claude Code ya sabunta a wannan lokacin ana riƙe ta, kuma ba a taɓa soke fita (logout). Idan amsar ta ba da tsawon rayuwar sabon refresh token, ana ajiye shi a matsayin `refreshTokenExpiresAt`, kamar yadda Claude Code ke yi.

Sabuntawar da ba ta yi nasara ba tana faɗin dalili, a layin asusun ko a saman lambobinsa da aka ajiye: `token not renewed: a session is running`, misali, ko `login may be lost: …` idan an aika buƙata amma ba a iya amfani da amsarta ko ajiye ta ba, tare da shawarar gudanar da `ccam login <name>` idan Claude Code ya nemi ka shiga daga baya.

### Mai tattara layin matsayi

Claude Code yana gudanar da umarnin layin matsayi kuma yana tura masa takardar JSON ta bututu (pipe). ccam yana karanta sunan model da `rate_limits` daga cikinta:

<!-- motion:usage-statusline -->

1. ccam yana nemo asusun daga `CCAM_PROFILE`, wanda ccam ke saitawa lokacin da yake fara Claude Code, ko kuma ta kwatanta `CLAUDE_CONFIG_DIR` da foldojin asusai (`runStatusline`).
2. ccam yana karanta `rate_limits.five_hour` da `rate_limits.seven_day` (`ParseStatusLine`). JSON ɗin ba shi da iyakokin amfani ga asusan API key, ko kafin amsar farko ta zaman; to ba a ajiye komai.
3. ccam yana rubuta `~/.local/state/ccam/usage/<name>.statusline.json` (`SaveStatusLine`). `ccam usage` yana amfani da shi idan endpoint ɗin bai amsa ba.
4. ccam yana buga `[name] · model · 5h N% · 7d N%`. Da `--wrap`, yana gudanar da umarninka ta `sh -c` da shigarwa iri ɗaya kuma yana buga fitarwarsa yadda Claude Code ke nuna layin matsayi: kowane layi da ba fanko ba ne, tare da `[name]` a gaban na farko. Umarnin da ya gaza yana nuna sunan asusun kaɗai. A Windows yana gudanar da umarnin yadda Claude Code ke yi: da Git Bash idan an girka shi, in ba haka ba da PowerShell.
5. Idan amfanin awa 5 ko na mako ya kai 90% ko fiye kuma akwai wani asusu, layin yana ƙarewa da `→ ccam next --auto --continue`, misali `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Umarnin da aka naɗe ba shi da iyakar lokaci na kansa: Claude Code yana ba umarnin layin matsayi minti 10 kuma yana dakatar da gudu idan sabon sabuntawa ya fara, don haka umarni mai jinkiri yana bayyana kamar yadda zai yi ba tare da ccam ba. Idan Claude Code ya ƙare ccam, umarnin ma yana ƙarewa (a Linux da macOS tare da kowane process da ya fara).

`ccam statusline install` yana rubuta `statusLine` a cikin `settings.json` na kowane asusu kuma yana riƙe sauran maɓallan abin, kamar `padding`, `refreshInterval` da `hideVimModeIndicator`. Idan asusu ya riga yana da umarnin layin matsayi, ccam yana riƙe shi a matsayin `--wrap '<old command>'`; a Windows a matsayin `--wrap-base64 <base64>`, wanda ba ya buƙatar alamomin zance a kowane shell. Idan layin matsayi ya riga yana gudanar da ccam (daga wata hanya, ko a cikin alamomin zance yadda tsohuwar siga ta rubuta shi), ccam yana maye gurbinsa maimakon ya naɗe shi, kuma yana riƙe umarnin da yake naɗewa.

A Windows ana rubuta hanyar ccam da slash na gaba da `~` don foldar gida (`~/.local/bin/ccam.exe statusline`), wanda Git Bash da PowerShell duka suke fahimta. Hanyar da ke buƙatar alamomin zance ana rubuta ta cikakkiya ga shell da Claude Code ke amfani da shi: `'…'` ga Git Bash, `& '…'` ga PowerShell. `settings.json` da aka raba ta rukuni ana rubuta shi sau ɗaya, ta mahaɗinsa, don haka kowane asusu a rukunin yana samun sa. Sake gudanar da shi yana buga `(already set)` ga fayilolin da suka riga suna amfani da ccam.

## Tsaro da yanayi na musamman

- `ccam usage` ba ya taɓa canza fayil ɗin token sai dai idan ka ba da `--refresh`. Idan access token ya ƙare, fara asusun sau ɗaya da `ccam <name>` kuma Claude Code zai sabunta shi. Shigar da ta ƙare gaba ɗaya tana buƙatar `ccam login <name>`.
- `--refresh` ba ya taɓa asusan da ke aiki (ko da an fara Claude Code ba ta ccam ba), asusan da aka shigo da su ko shigar Keychain. Sabuntawar da aka ƙi tana nuna lambobin da aka ajiye da dalilin (`token not renewed: …`).
- Ba a rubuta bayanin endpoint na amfani a bainar jama'a ba, kuma yana iya canzawa. Idan ya canza, har yanzu kana ganin lambobin da aka ajiye, kuma `ccam doctor` yana bayar da rahoton canjin tsarin.
- Idan babu lambobin da aka ajiye kuma babu amsa, layin asusun yana nuna `lookup failed:` da dalilin, misali `no saved usage` da `--source cache`.
- Fayilolin cache suna riƙe kashi da lokutan sake saiti kawai, ba taɓa token ba.

## Masu alaƙa

- [Iyakokin amfani da sauya asusai](/docs/usage-limits)
- [Sauyawa zuwa asusu na gaba](/features/next)
- [Claude Code ya kai iyakar amfani? Ga abin da za ka yi](/guides/usage-limit-workflow)
- [Bayanin umarnai](/docs/commands)
