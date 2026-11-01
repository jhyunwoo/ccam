---
title: Mga usage limit
description: Ipinapakita ng ccam usage ang 5-oras at lingguhang limit ng bawat account nang sabay; inilalagay ng ccam statusline ang account at usage nito sa Claude Code.
updated: 2026-11-01
order: 7
translation: machine
sourceHash: "54b296f7da35"
---

Ipinapakita ng `ccam usage`, para sa bawat account, kung gaano kalaki na ang nagamit sa bawat limit: ang 5-oras na limit, ang lingguhang limit at ang sariling lingguhang limit ng bawat model (gaya ng Fable). May bar ang bawat limit at ang oras bago ito mag-reset. Gamitin ito bago ka magsimula ng mahabang gawain, o kapag sinabi ng Claude Code na malapit na ang isang limit, para makita kung aling account ang may natitira pa. Ang `ccam statusline` ay ang status line command para sa Claude Code: ipinapakita nito ang pangalan ng account at usage sa loob ng Claude Code at sine-save ang mga numero para sa `ccam usage`. Para sa mga token at tinatayang gastos sa likod ng usage, pinapatakbo ng `ccam ccusage` ang [ccusage](https://ccusage.com) sa mga usapan ng bawat account ([sa ibaba](#tokens-and-costs-ccam-ccusage)).

## Paggamit

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

Bawat account ay may linyang may numero nito, pangalan nito at kung gaano na katagal ang mga numero nito; ang ibig sabihin ng `(cached)` ay galing ang mga ito sa naka-save na file, hindi sa endpoint. Sa ilalim nito, ipinapakita ng isang linya bawat limit ang 20-cell na bar, ang porsiyento at ang oras bago mag-reset ang window: `5-hour`, `Weekly` (lahat ng model), at pagkatapos ang sariling lingguhang window ng bawat model, gaya ng `Fable`. Nagpapakita ng `-` ang window na wala sa account. Nagiging dilaw ang porsiyento sa 80% at pula sa 100%.

Kung tinanggihan ng Claude Code ang isang request dahil puno ang isang window, nagpapakita ang window na iyon ng 100% hanggang mag-reset ito, anuman ang sabihin ng mga naka-save na numero (tingnan ang [tala ng limit](#the-limit-record)). Sa terminal, lumalabas agad ang bawat account bilang `fetching…`, at iginuguhit ang mga linya nito habang sumasagot ito.

Pinapangalanan ng huling linya, na lumalabas lang kapag mahigit sa isa ang account mo, ang account na pipiliin ng [`ccam next --auto`](/features/next). Sinusunod nito ang parehong mga patakaran:

- `auto_max_usage` at `auto_strategy`;
- ang mga tumatakbong session;
- mga account lang ng organisasyon ng huling ginamit na account, maliban kung nakatakda ang `cross_org`;
- nilalaktawan ang bawat window na nasa 100%, kaya hindi kailanman pinapangalanan ang account na puno na ang linggo.

Sa default na `headroom` strategy, `→ Most headroom: <name>` ang nakasulat dito; sa `expiring` o `balance`, ibinibigay nito ang dahilan, gaya ng `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Mga seat ng Team at Enterprise.** Lumilipat din ang awtomatikong pagpili ng account (`ccam next`, `--fallback`, `ccam pick`, ang `n` ng dashboard at, kapag binuksan mo ito, ang [`auto_switch`](/features/next#auto-switch)) sa pagitan ng mga seat ng iisang Claude Team o Enterprise na organisasyon, na maaaring lumabag sa terms ng Anthropic. Basahin ang [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) at magtanong muna sa iyong organisasyon at sa Anthropic bago ka umasa rito.

Nagpapakita ng `not logged in → ccam login <name>` ang account na hindi naka-log in, at ng `no limits (API)` ang account na API key ang gamit. Nagpapakita ang expired na access token ng ``access token expired; `ccam <name>` renews it when it starts``, sa itaas ng mga naka-save na numero, kung mayroon: nire-renew ito ng Claude Code sa susunod na pagsisimula. Ang login na buong nag-expire, dahil nag-expire ang refresh token nito o binura ito ng Claude Code matapos itong tanggihan ng server, ay nagpapakita ng `` login expired → `ccam login <name>` ``.

Pagkatapos ng `ccam statusline install`, magpapakita ang Claude Code ng linya gaya ng `[work] · Opus · 5h 24% · 7d 41%`.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--json` | I-print ang mga ulat bilang JSON: `name`, `slot`, `status`, `limits` (kasama ang `credits` kapag naka-on ang mga usage credit ng account), `source`, `at`, `error`, at `limited` at `limited_windows` kapag itinatakda ng tala ng limit sa 100% ang mga window. Kapag walang account, ipi-print nito ang `[]` | off |
| `--watch[=N]` | Iguhit ulit tuwing `N` segundo hanggang pindutin mo ang Ctrl+C. Ang `--watch` lang, o value na 0 o mas mababa, ay nangangahulugang 60. Inilalagay ang value pagkatapos ng `=`: tinatanggihan ang `--watch 10` kasama ang isang pahiwatig | off |
| `-v`, `--verbose` | Magdagdag ng linyang `source:` sa ilalim ng bawat account: kung saan galing ang mga numero nito (ang usage endpoint, ang Claude Code, ang status line, isang naka-save na file, o ang usage limit na iniulat ng Claude Code). Kapag nangingibabaw sa mga numero ang tala ng limit, susunod ang linyang gaya ng `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` | off |
| `--no-cache` | Huwag pansinin ang 60 segundong cache at magtanong ulit sa endpoint | off |
| `--source live\|claude\|cache` | Nagtatanong sa endpoint ang `live` gamit ang token ng account; ipinapahanap ng `claude` ito sa Claude Code (tingnan sa ibaba); ipinapakita lang ng `cache` ang mga naka-save na numero | `usage_source`, kung wala ay `live` |
| `--refresh` | I-renew ang mga expired na access token, para lang sa mga account na hindi tumatakbo. Puwedeng tumagal nang hanggang 30 segundo ang isang pag-renew | off |
| `--dry-run` | Tinatanggap: walang binabago ang `ccam usage` maliban sa sarili nitong cache. Tinatanggihan kapag kasama ang `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | off |

Pumapalya ang anumang ibang value ng `--source` na may `--source must be live, claude or cache`. Walang tinatanggap na argument ang `ccam usage`.

### Usage mula mismo sa Claude Code

Kapag may `usage_source = "claude"` sa itaas ng `config.toml` (bago ang unang `[[group]]`, `[[profile]]` o `[[folder]]`), o `--source claude`, hindi binabasa ng ccam ang login token kahit kailan. Para sa bawat account, sinisimulan nito ang `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` sa folder ng account at ipinapadala ang `get_usage` request ng Claude Code sa standard input (Claude Code 2.1.270 o mas bago). Hinahanap ng Claude Code ang usage gamit ang sarili nitong login at sumasagot ito ng 5-oras, lingguhan at model na mga window.

Pinipigilan ng `--safe-mode` na magsimula ang mga hook, MCP server at plugin ng account, at walang session na sine-save. Kailangan ng bawat account ng isang pagsisimula ng Claude Code (mga isang segundo), kaya mas mabagal ang source na ito kaysa sa `live`. Sabay-sabay na tinatanong ang mga account, bawat isa na may 20 segundong timeout, at gumagana ang 60-segundong cache gaya ng dati.

Binibigyan ng 50 segundo ang login na nag-expire na ang access token, o mag-e-expire sa loob ng 5 minuto, dahil nire-renew ito ng Claude Code habang nagsisimula at puwedeng mawala ang login kapag naputol ang pag-renew. Kapag sinukuan ng ccam ang isang lookup, hinihiling nito sa Claude Code na lumabas gamit ang SIGTERM, para matapos ng Claude Code ang pag-renew na nasimulan na nito (walang SIGTERM ang Windows, kaya doon, tinatapos ang process).

Para sa login na buong nag-expire, ipinapakita ang `` login expired → `ccam login <name>` `` at hindi na sinisimulan ang Claude Code. Itinuturing ng Claude Code na experimental ang request, kaya kapag huminto itong sumagot, ipinapakita ng ccam ang mga naka-save na numero at iniuulat ito ng `ccam doctor`. Kapag sumagot ang Claude Code na pumalya ang sarili nitong lookup, ipinapakita ng ccam ang mga naka-save na numero na may `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`, hindi isang plan na walang limit.

Sa `--dry-run`, hindi sinisimulan ng ccam ang Claude Code, dahil nagsusulat ang Claude Code sa folder ng account; ipinapakita na lang nito ang mga naka-save na numero. Tumatanggap ang `usage_source` ng `live` o `claude`; pagpipilian lang ng `--source` ang `cache`.

Binabasa ng `ccam statusline` ang status line JSON ng Claude Code mula sa standard input. Iisa ang option nito, ang `--wrap <command>`: isang umiiral na status line command na ang unang linya ng output ay lalagyan ng pangalan ng account sa unahan. Walang option ang `ccam statusline install`. Parehong tinatanggihan ng dalawa ang `--dry-run` na may `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Mga token at gastos: `ccam ccusage`

Ipinapakita ng `ccam usage` kung gaano kalaki na ang nagamit ng isang account sa bawat limit. Ang [ccusage](https://ccusage.com) naman ay nagbibilang ng mga token at tinatayang gastos sa mga log ng usapan ng Claude Code. Kapag mag-isa, sa `~/.claude` at `~/.config/claude` lang nito hinahanap ang mga log ng Claude Code (o sa mga folder na pinapangalanan ng `CLAUDE_CONFIG_DIR`), kaya ang tanging ccam account na nakikita nito ay ang in-import sa kinalalagyan nito sa `~/.claude`. Pinapatakbo ito ng `ccam ccusage` nang nakatakda ang `CLAUDE_CONFIG_DIR` sa mga folder ng bawat account na may mga usapan, kaya iisang ulat ang sumasaklaw sa kanilang lahat:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

Binabasa muna ng ccam ang sarili nitong mga option: tumatanggap ang `-a` (`--account`) ng mga pangalan o numero ng account, pinaghihiwalay ng kuwit o inuulit, at ipi-print ng `--dry-run` ang command sa halip na patakbuhin ito. Mula sa unang argument na hindi option ng ccam, ipinapasa ang lahat sa ccusage gaya ng pagkaka-type mo, maliban sa `--config` para sa `bun x` at `npx` (sa ibaba); ipinapakita ng `ccam ccusage -- --help` ang sariling help ng ccusage. Pinapatakbo ng ccam ang `ccusage` kung naka-install ito, kung hindi ay ang `bun x ccusage`, at kung wala pa rin ay ang `npx ccusage@latest`; kinukuha ito ng huling dalawa mula sa npm. Nagsisimula ang mga ito sa sariling folder ng ccam (`~/.local/state/ccam/ccusage`), hindi sa folder kung nasaan ka, kaya hindi mapipili ng isang project doon kung ano ang pinapatakbo ng mga ito gamit ang `node_modules/.bin` nito, isang `bin` sa `package.json` nito o isang registry sa `.npmrc` nito; ipinapakita ng `ccam --dry-run ccusage` ang `cd` papunta sa folder na iyon. Kapag walang available sa mga ito, sinasabihan ka ng ccam na i-install ang Node.js, Bun o ccusage. Sa Linux at macOS, pumapalit ang ccusage sa ccam; sa Windows, hinihintay ito ng ccam at lumalabas ang ccam gamit ang exit code ng ccusage. Sa WSL, nilalaktawan ng ccam ang `ccusage`, `bun` o `npx` na nasa isang Windows drive, dahil tatakbo ito bilang Windows program nang wala ang mga folder ng account.

- Sakop ng ulat ang mga account na may mga usapan ang folder (isang `projects` folder). Isang beses lang binibilang ang usapang shine-share ng mga account sa pamamagitan ng isang group.
- Sa `-a`, dinadala ng account na nagshe-share ng kasaysayan nito sa pamamagitan ng isang group ang mga usapan ng buong group, at sinasabi ng ccam kung kanino ang mga ito: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- Ipinapakita ng `blocks` ang mga 5-oras na window, at may sariling window ang bawat account, kaya pangalanan ang iisang account para rito, gaya ng nasa itaas.
- Iniuulat din ng ccusage ang iba pang coding tool na nakikita nito, gaya ng Codex, maliban kung pangalanan mo ang isa: Claude Code lang ang sakop ng `ccam ccusage claude daily`.
- Habang pinapatakbo ito ng ccam, binabasa lang ng ccusage ang mga setting nito (`ccusage.json`) mula sa `./.ccusage` at sa mga folder ng account. Pinapangalanan ng `--config` ang ibang file: `ccam ccusage daily --config ~/.config/claude/ccusage.json`. Dahil nagsisimula ang `bun x` at `npx` sa sariling folder ng ccam, ipinapasa sa kanila ng ccam ang `./.ccusage/ccusage.json` gamit ang `--config` kapag umiiral ang file na iyon at wala kang pinangalanan, at ibinibigay nito ang relative na `--config` bilang buong path. Pagkatapos, binabasa ng ccusage ang iba pang relative na path, gaya ng `--pi-path`, mula sa folder ng ccam, kaya ibigay nang buo ang mga iyon.
- Babasahin ng ccusage bilang dalawang folder ang folder na may kuwit sa path nito, kaya tinatanggihan ito ng ccam at pinapangalanan ang account; gamitin ang `-a` para hindi isama ang account na iyon.
- Command ng ccam ang `ccusage`, kaya sinisimulan ang account na pinangalanang ccusage gamit ang `ccam run ccusage`; sinasabi iyon ng `ccam doctor`.

## Paano ito gumagana

Galing sa dalawang lugar ang mga numero ng `ccam usage`: ang hindi dokumentadong usage endpoint ng Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, na tinatawag gamit ang sariling login ng bawat account, at ang mga value na ibinibigay ng Claude Code sa status line command nito. Nasa `internal/usage` ang lahat ng ito.

1. Nilo-load ng ccam ang `config.toml` at bumubuo ng isang `usage.Account` bawat account (`cli.usageCmd`): ang folder nito na `~/.local/share/ccam/profiles/<name>/`, ang mga session marker nito (kasama ang sariling mga tala ng Claude Code sa folder ng account, sinasabi ng mga ito kung tumatakbo ito) at kung in-import ito sa kinalalagyan nito.
2. Sabay-sabay na hinahanap ng ccam ang usage ng bawat account, isang goroutine bawat account (`Fetcher.FetchAll`). Iisang `http.Client` ang gamit nilang lahat, kaya muling ginagamit ang mga koneksyon sa endpoint. Puwedeng umabot nang 2 s ang pagkonekta, at may kabuuang 3 s ang bawat account. Hindi sumusunod sa mga redirect ang client, dahil ipapadala ulit ng redirect ang token sa ibang lugar, at hanggang 1 MB lang ng sagot ang binabasa nito.
3. Habang natatapos ang bawat account, iginuguhit ulit ang hilera nito. Ang account na pumapalya o hindi sumasagot sa oras ay bumabalik sa mga naka-save na numero, kaya hindi kailanman napapabagal ng isang mabagal na account ang iba.

Sa diagram, hinahanap ang bawat account sa sarili nitong goroutine, at ang nag-time out ay nagpapakita ng mga naka-save na numero nito sa halip:

<!-- motion:usage-fetch -->

Sumasagot sa oras ang `personal` at `work` at nagpapakita ng mga live na numero. Lumalampas sa 3 s ang `side`, kaya ipinapakita na lang ng ccam ang mga naka-save na numero nito (`fallback()`). Pagkatapos ay iginuguhit ng ccam ang isang block bawat account (`usage.Table()`) at pinipili ang account na pinapangalanan sa huling linya (`usage.Best()`).

Para sa isang account, ginagawa ng ccam ang mga pagsusuring ito nang sunod-sunod (`Fetcher.Fetch`):

<!-- motion:usage-cache -->

1. Sa `--source cache`, direkta itong pumupunta sa `fallback()`.
2. Maliban kung nakatakda ang `--no-cache`, ginagamit ang live na resulta sa `~/.local/state/ccam/usage/<name>.json` na sinukat wala pang 60 segundo ang nakalipas (`CacheTTL`), na hinuhusgahan ayon sa kasalukuyan: nagpapakita ng 0% ang window na nag-reset na mula noon, at hindi itinuturing na bago ang file na may petsang lampas sa orasan. Isang request bawat account lang ang gastos kapag dalawang beses na sunod na pinatakbo ang `ccam usage`.
3. Binabasa ng ccam ang `.credentials.json` (o ang macOS Keychain) nang hindi ito binabago, at ang key na `claudeAiOauth` lang ang dine-decode nito (`readCredentials()`). Kapag walang subscription login, `not logged in` ang resulta, o `no limits (API)` kapag may API key ang `.claude.json`. Ang login na nag-expire na ang refresh token, o na binura ng Claude Code, ay `` login expired → `ccam login <name>` `` at hindi na hinahanapan ng usage.
4. Hindi ginagamit ang expired na access token. Kapag walang `--refresh`, ipinapakita ng account ang mga naka-save na numero nito na may ``access token expired; `ccam <name>` renews it when it starts``.
5. Kung hindi, tinatawag ng ccam ang endpoint. Itinuturing na expired na login ang 401 o 403 (`login expired: authentication failed (401)`), itinuturing na error ang anumang ibang status o body na hindi nito mabasa, at parehong bumabalik ang mga ito sa mga naka-save na numero. Sinasabi ng lookup na tinapos ng Ctrl+C ang `cancelled`. Ang 429 ay nangangahulugang nakatanggap ang endpoint ng sobrang daming request para sa account na ito: isinusulat ng ccam ang `<name>.pause` kasama ang oras kung kailan magtatanong ulit (`Retry-After`, o 5 minuto, hindi hihigit sa isang oras), ipinapakita ang mga naka-save na numero na may `rate limited by the usage endpoint; asking again after 15:10`, at hindi tinatawag ang endpoint para sa account na iyon hanggang sa oras na iyon, kahit may `--no-cache`. Lalo lang mapapatagal ng agarang pagsubok ulit ang rate limit.
6. Isinusulat ang maayos na sagot sa `<name>.json` kasama ang oras ng pagsukat.

Para sa mga naka-save na numero, kinukuha ng ccam ang pinakabagong naka-save na lookup (`<name>.json`, o `<name>.claude.json` mula sa mga lookup ng Claude Code) at, kapag isinulat ang `<name>.statusline.json` pagkatapos nito, ipinapatong dito ang 5-oras at lingguhang window ng status line (`fallback()`); nananatili gaya ng pagkaka-save ng lookup ang mga model window, na wala sa status line. Bago nito ipakita ang mga naka-save na numero, sinusuri nito ang oras ng reset ng bawat window: nagpapakita ng 0% ang window na nag-reset na mula nang i-save ang mga numero, dahil hindi na totoo ang lumang porsiyento. Kaya ang `personal` sa itaas ay nagpapakita ng `0%` na walang oras ng reset.

May dalawang format ng response ang endpoint. May array na `limits` ang mas bago na may mga entry na `session` (5-oras), `weekly_all` at `weekly_scoped` (bawat model). May mga object na `five_hour`, `seven_day`, `seven_day_opus` at `seven_day_sonnet` ang mas luma. Binabasa ng ccam ang dalawa (`ParseLive`): nananaig ang array na `limits`, at pinupunan ng mga lumang object ang hindi nito sakop. Iniuulat bilang nagbagong format ang body na walang kahit anong limit, at sasabihin iyon sa iyo ng `ccam doctor`.

### Ang tala ng limit

Kapag pinahinto ng Claude Code ang isang turn sa usage limit, binabasa ng hook ng [ccam plugin](/docs/install#claude-code-plugin) ang tinanggihang request mula sa usapan: kung aling window ang puno (5-oras, lingguhan, Opus o Sonnet) at kung kailan ito magre-reset. Isinusulat nito iyon sa `~/.local/state/ccam/usage/<name>.limit.json`. Ipinapatong ng ccam ang tala sa ibabaw ng anumang iba pang alam nito, mula sa bawat source (`Fetch`): itinuturing na 100% ang bawat window na nasa tala hanggang sa reset nito, gaano man kaluma o kababa ang mga naka-save na numero, kaya hindi pinipili ulit ng `ccam next --auto`, `ccam pick`, `--fallback` at [`auto_switch`](/features/next#auto-switch) ang account bago iyon.

Kapag wala nang ibang alam, ang tala lang ang ipinapakita ng account (`source: the usage limit Claude Code reported` kasama ang `-v`). Walang talang isinusulat ang mga session na walang nagbabantay (mga `-p` run at background session), gayundin ang pagpapahintong itinakda ng isang organisasyon (halimbawa, kapag naubos na ang usage nito o $0 ang limit ng isang group), na puwedeng alisin ng mga admin nito anumang oras. Inaalis ng `ccam login` sa ibang Claude account ang tala; nabubura rin ito kapag binura ang file.

### Pag-renew ng mga token

Nagro-rotate ang mga refresh token: kapag ni-renew ang isa, makakakuha ka ng bagong refresh token at magiging invalid ang luma. Kung ire-renew ng ccam ang login habang tumatakbo ang Claude Code bilang ang parehong account, masa-sign out ang session na iyon. Kaya binabasa lang ng `ccam usage` ang mga token at hindi kailanman sinusulatan, maliban kung ipapasa mo ang `--refresh`. Pagkatapos, ang ccam (`Fetcher.refresh`) ay:

1. tumatanggi sa mga account na nirehistro sa kinalalagyan nila gamit ang `ccam import`, dahil baka ginagamit ng karaniwang `claude` ang login na iyon nang hindi alam ng ccam;
2. tumatanggi sa mga login na naka-store sa macOS Keychain;
3. kumukuha ng lock file, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, para hindi kailanman sabay na ma-renew ng dalawang ccam process ang parehong token; habang hawak ito, ipinapakita ng account ang `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)`. Kinukuha rin ito ng `ccam rename`, `remove` at `logout`, at tumatanggi sila habang may isinasagawang pag-renew;
4. habang hawak ang lock, tumatanggi kung may tumatakbong session ng account, paano man ito sinimulan (sa pamamagitan ng ccam, karaniwang `claude`, isang IDE, o isang background session at ang daemon sa likod nito, na mismong nagre-renew ng login), at nilalaktawan ang request kung ni-renew na ng ibang process ang token;
5. walang ipinapadala habang hawak ng Claude Code ang sarili nitong mga renewal lock sa folder ng account (`.oauth_refresh.lock`, `<folder>.lock`): nire-renew ng Claude Code ang login na iyon sa mismong sandaling iyon;
6. ipinapadala ang refresh token sa token endpoint at naghihintay nang hanggang 30 segundo para sa sagot, ang oras na ibinibigay ng Claude Code sa sarili nitong mga pag-renew. Kapag naipadala na, hindi kailanman pinuputol ang request, kahit ng Ctrl+C: baka na-rotate na ng server ang refresh token, at ang response na lang nito ang may dala ng bagong login. Tinatapos ng Ctrl+C o ng hangup ang ibang mga lookup, at lalabas ang ccam (130 pagkatapos ng Ctrl+C) kapag na-save na ang sagot;
7. nagse-save sa ilalim ng write lock ng Claude Code (`.storage-write.lock`) at binabago lang ang mga token field ng `.credentials.json` gaya ng kalagayan nito sa sandaling iyon: nananatili ang isinulat ng Claude Code sa pagitan (isang MCP sign-in), pinapanatili ang login na ni-renew ng Claude Code sa pagitan, at hindi kailanman binabawi ang isang logout. Kapag ibinigay ng sagot ang tagal ng bisa ng bagong refresh token, iniimbak ito bilang `refreshTokenExpiresAt`, gaya ng ginagawa ng Claude Code.

Sinasabi ng refresh na hindi natuloy kung bakit, sa linya ng account o sa itaas ng mga naka-save na numero nito: `token not renewed: a session is running`, halimbawa, o `login may be lost: …` kapag naipadala ang request pero hindi nagamit o na-save ang sagot nito, kasama ang payong patakbuhin ang `ccam login <name>` kung hihilingin ng Claude Code na mag-log in ka.

### Ang status line collector

Pinapatakbo ng Claude Code ang status line command at ipina-pipe dito ang isang JSON document. Binabasa ng ccam mula rito ang pangalan ng model at ang `rate_limits`:

<!-- motion:usage-statusline -->

1. Hinahanap ng ccam ang account mula sa `CCAM_PROFILE`, na itinatakda ng ccam kapag sinisimulan nito ang Claude Code, o kung hindi ay sa pagtutugma ng `CLAUDE_CONFIG_DIR` sa mga folder ng account (`runStatusline`).
2. Binabasa ng ccam ang `rate_limits.five_hour` at `rate_limits.seven_day` (`ParseStatusLine`). Walang rate limit ang JSON para sa mga account na API key ang gamit, o bago ang unang response ng session; sa ganitong kaso, walang sine-save.
3. Isinusulat ng ccam ang `~/.local/state/ccam/usage/<name>.statusline.json` (`SaveStatusLine`). Ginagamit ito ng `ccam usage` kapag hindi sumasagot ang endpoint.
4. Ipi-print ng ccam ang `[name] · model · 5h N% · 7d N%`. Sa `--wrap`, pinapatakbo nito ang iyong command sa pamamagitan ng `sh -c` na may parehong input at ipi-print ang output nito sa paraang ipinapakita ng Claude Code ang status line: bawat linyang hindi blangko, na may `[name]` sa unahan ng una. Pangalan lang ng account ang ipinapakita ng command na pumalya. Sa Windows, pinapatakbo nito ang command sa paraan ng Claude Code: gamit ang Git Bash kapag naka-install ito, o kung hindi, gamit ang PowerShell.
5. Kapag nasa 90% o higit pa ang 5-oras o lingguhang usage at may iba pang account, nagtatapos ang linya sa `→ ccam next --auto --continue`, halimbawa `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Walang sariling time limit ang naka-wrap na command: binibigyan ng Claude Code ng 10 minuto ang status line command at pinapahinto ang isang run kapag nagsimula ang mas bagong update, kaya lumalabas ang mabagal na command gaya ng paglabas nito kung walang ccam. Kapag tinapos ng Claude Code ang ccam, natatapos din ang command (sa Linux at macOS, kasama ang bawat process na sinimulan nito).

Isinusulat ng `ccam statusline install` ang `statusLine` sa `settings.json` ng bawat account at pinapanatili ang ibang key ng object, gaya ng `padding`, `refreshInterval` at `hideVimModeIndicator`. Kapag may status line command na ang isang account, pinapanatili ito ng ccam bilang `--wrap '<old command>'`; sa Windows bilang `--wrap-base64 <base64>`, na hindi nangangailangan ng quoting sa alinmang shell. Kung nagpapatakbo na ng ccam ang status line (mula sa ibang path, o naka-quote sa paraang isinulat ito ng mas lumang bersyon), pinapalitan ito ng ccam sa halip na i-wrap, at pinapanatili ang command na naka-wrap dito.

Sa Windows, isinusulat ang path papunta sa ccam gamit ang forward slash at `~` para sa home folder (`~/.local/bin/ccam.exe statusline`), na naiintindihan ng Git Bash at ng PowerShell. Isinusulat nang buo ang path na nangangailangan ng quote para sa shell na ginagamit ng Claude Code: `'…'` para sa Git Bash, `& '…'` para sa PowerShell. Isang beses lang isinusulat ang `settings.json` na shine-share sa pamamagitan ng group, sa pamamagitan ng link nito, kaya nakukuha ito ng bawat account sa group. Kapag pinatakbo ulit, ipi-print ang `(already set)` para sa mga file na gumagamit na ng ccam.

## Kaligtasan at mga espesyal na kaso

- Hindi kailanman binabago ng `ccam usage` ang token file maliban kung ipapasa mo ang `--refresh`. Kapag nag-expire ang isang access token, patakbuhin nang isang beses ang account gamit ang `ccam <name>` at ire-renew ito ng Claude Code. Kailangan ng `ccam login <name>` ang login na buong nag-expire.
- Hindi kailanman ginagalaw ng `--refresh` ang mga tumatakbong account (kahit sinimulan ang Claude Code nang walang ccam), in-import na account o Keychain login. Ipinapakita ng tinanggihang refresh ang mga naka-save na numero at ang dahilan (`token not renewed: …`).
- Hindi pampublikong dokumentado ang usage endpoint at puwede itong magbago. Kapag nangyari iyon, makikita mo pa rin ang mga naka-save na numero, at iuulat ng `ccam doctor` ang pagbabago sa format.
- Kapag walang naka-save na numero at walang sagot, ipinapakita ng linya ng account ang `lookup failed:` at ang dahilan, halimbawa `no saved usage` kapag may `--source cache`.
- Mga porsiyento at oras ng reset lang ang laman ng mga cache file, hindi kailanman mga token.

## Kaugnay

- [Mga usage limit at paglipat ng account](/docs/usage-limits)
- [Paglipat sa susunod na account](/features/next)
- [Kapag naabot mo ang Claude Code usage limit](/guides/usage-limit-workflow)
- [Sanggunian ng mga command](/docs/commands)
