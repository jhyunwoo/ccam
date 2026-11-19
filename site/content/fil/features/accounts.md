---
title: Pamamahala ng mga account
description: Magdagdag, mag-log in, maglista, magpalit ng pangalan, mag-alis at mag-import ng Claude Code account gamit ang ccam add, login, list, rename, remove at import.
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

Ang bawat ccam account ay hiwalay na Claude Code config folder na may sariling login. Gamitin ang mga command na ito para magdagdag ng mga account, i-sign in ang mga ito, ilista, palitan ang pangalan at alisin ang mga ito. Nagrerehistro rin ang mga ito ng mga folder na mayroon ka na.

## Paggamit

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account that was never logged in or whose login expired
ccam logout side                           # log out: the login and its MCP sign-ins
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, logs out, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Saanman tumatanggap ang isang command ng `<name|number>`, puwede mong ibigay ang pangalan ng account o ang slot number nito (ang column na `#` sa `ccam list`).

Binabago lang ng `ccam use` kung aling account ang default (`default` sa `config.toml`).

Kapag pinatakbo sa terminal nang walang option, nagbubukas ang `ccam add` ng wizard. Dumadaan ito sa mga hakbang na ito:

1. ang pangalan ng account, maliban kung nagbigay ka na;
2. kung paano magsa-sign in ang account: Claude login ngayon, API key o Claude login mamaya;
3. kung aling mga group ang sasalihan, na naka-check ang `default` (kapag wala pang group, ipinapakita nito kung ano ang ishe-share ng group na `default`);
4. `Create a shortcut command for work? [y/N]`;
5. isang buod na kukumpirmahin mo.

Bumabalik nang isang hakbang ang `Esc`. Nilalaktawan ang wizard kapag may `--group`, `--no-group`, `--shortcut` o `--no-login`, o kapag pinatakbo nang walang terminal (halimbawa, mula sa isang pipe). Sa ganitong kaso, kailangan ang pangalan, at walang shortcut na gagawin maliban kung may pinangalanan ang `--shortcut`.

Kapag walang `--group`, sumasali ang bagong account sa group na `default`. Ang unang account ang gumagawa nito, na nagshe-share ng mga MCP server, settings, `CLAUDE.md` at mga extension; tingnan ang [pag-share gamit ang mga group](/features/share). Sa Windows, nali-link lang ang mga file kapag naka-Developer Mode. Kapag wala ito, ginagawa ng unang account ang `default` na may `mcp` at `extensions` lang at sinasabi ito (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), at tinatanggihan ng ccam ang `--group` na may group na nagshe-share ng settings, memory o kasaysayan.

Ipinapakita ng `ccam list` ang slot, pangalan, shortcut, email, mga group at huling paggamit, na may `*` bago ang default na account. Ang email ay ang kasalukuyang naka-log in sa folder (mula sa `.claude.json` nito); kapag iba ito sa itinala ng ccam, susunod ang itinala bilang `(was …)`:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

```console
$ ccam list
   #  Name      Shortcut  Email                                     Groups   Last used
*  1  work      ccwork    other@company.com (was work@example.com)  default  2m ago
```

Sa `--json`, object ang bawat account na may `slot`, `shortcut`, `name`, `email` (ang sa kasalukuyang login, gaya ng sa table), `default`, `dir`, `last_used` (wala kapag hindi pa kailanman sinimulan ang account), `groups` at `share`. Itinutugma ng `share` ang bawat item (`mcp`, `settings`, `memory`, `extensions`, `history`) sa group na pinanggagalingan nito, o sa `"local"`.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `ccam add --no-login` | Gumawa ng account nang hindi nagla-log in; mag-log in mamaya gamit ang `ccam login <name>` | off |
| `ccam add --group <group>` | Sumali sa group na ito sa halip na sa `default`. Puwedeng ulitin; ginagawa ang `default` kapag wala pa ito | `default` |
| `ccam add --no-group` | Huwag sumali sa anumang group: walang shine-share ang account | off |
| `ccam add --shortcut <command>` | Gumawa rin ng shortcut command na ito sa `~/.local/bin` | wala |
| `ccam add --api --key-from <source>` | Account na gumagamit ng API key sa halip na login; tingnan ang [Mga API key account](/features/api) | off |
| `ccam list --json` | I-print ang mga account bilang JSON | off |
| `ccam remove -y`, `--yes` | Mag-alis nang hindi nagtatanong | off |
| `ccam logout -y`, `--yes` | Mag-log out nang hindi nagtatanong | off |
| `ccam import --name <name>` | Pangalan ng account para sa folder na ibinigay mo; tinatanggihan kapag walang folder (`--name needs a folder: ccam import <folder> --name ops`) | kinukuha mula sa pangalan ng folder |
| `ccam import -y`, `--yes` | I-import ang bawat nahanap na folder nang hindi nagtatanong | off |
| `--dry-run` | Sinusuportahan ng `ccam list` at `ccam completion`, na walang binabago. Tinatanggihan ng `add`, `login`, `logout`, `use`, `rename`, `remove` at `import`, na may `` --dry-run is not supported by `ccam add`; nothing was changed `` (kasama ang pangalan ng command) | off |

Walang sariling option ang `ccam login` (bukod sa `--expired`), `ccam use` at `ccam rename`. Tumatanggap ang `ccam completion` ng `bash`, `zsh` o `powershell`; pumapalya ang iba na may `only bash, zsh and powershell are supported`.

Ang mga pangalan ng account ay 1 hanggang 32 character ng `a-z`, `0-9`, `-` at `_`, na nagsisimula sa titik o numero. Hindi puwedeng puro numero ang pangalan, dahil ang ibig sabihin ng `ccam 2` ay slot 2. Hindi rin puwedeng ccam command ang pangalan (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), kaya hindi kailanman malabo ang `ccam <name>`. Ang listahan ay `config.Reserved` sa `internal/config`.

## Paano ito gumagana

Ang pagdaragdag ng account ay gumagawa ng folder nito, sumasali sa mga group nito, nagse-save ng config, nagsa-sign in at gumagawa ng shortcut, sa ganitong pagkakasunod-sunod. Kapag pinatakbo mo ang `ccam add work --shortcut ccwork`, ginagawa ng ccam ang mga hakbang na ito (`addProfile` sa `internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. Sinusuri ng ccam ang pangalan (`config.ValidateName`) at ang shortcut (`shim.Check`) bago ito may gawin.
2. Ginagawa ng ccam ang `~/.local/share/ccam/profiles/work/` na may mode na `0700` at ibinibigay sa account ang pinakamababang bakanteng slot number, kaya muling ginagamit ang numero ng inalis na account (`profile.Add`, `Config.NextSlot`). Ang unang account ang nagiging default.
3. Sumasali ang account sa mga group nito (`default` maliban kung iba ang sinasabi ng `--group` o `--no-group`; ang unang account ang gumagawa ng `default`), at nili-link ng ccam ang bago at walang lamang folder sa mga iyon, halimbawa ang `settings.json` sa `groups/default/settings.json` (`applyShare`, gamit ang `share.Build` at saka ang `Apply`). Kapag pumalya ito, inaalis ang bagong folder, dahil haharangin ng natirang folder ang bawat pagsubok ulit na may `folder already exists`.
4. Sine-save ang account sa `~/.config/ccam/config.toml` sa pamamagitan ng `fsutil.AtomicWrite`: isang pansamantalang file, `fsync`, saka `rename`, kaya hindi kailanman mag-iiwan ng kalahating naisulat na config ang isang crash.
5. Sinisimulan ng ccam ang `claude auth login` na nakatakda ang `CLAUDE_CONFIG_DIR` sa folder ng account, at ang `--email` na may huling email ng account para napunan na ang login page (`runLogin`). Kapag pumalya ang login, nananatili ang account, at ipi-print ng ccam ang error at sasabihing patakbuhin mo ang `ccam login work` mamaya.
6. Binabasa ng ccam ang `oauthAccount` mula sa `profiles/work/.claude.json` (`profile.ReadIdentity`): ang email, ang mga ID ng account at organisasyon, ang pangalan ng organisasyon at ang plan (`claude_max`, `claude_team` …). Puwedeng ilang megabyte ang file na iyon, kaya hinahanap ng ccam ang iisang key sa halip na i-parse ang buong file (`TopLevelValue`). Itinatabi ng ccam ang mga detalyeng ito sa `[profile.login]` sa `config.toml`.
7. Nili-link ng ccam ang `~/.local/bin/ccwork` sa ccam binary, na hinahanap ang `argv[0]` sa `config.toml` kapag nagsisimula ito, at itinatala ang link sa `~/.local/state/ccam/shortcuts.json` (`shim.Install`). Kapag hindi magawa ang link, idinadagdag pa rin ang account, sinasabi ng ccam kung bakit (`Could not create shortcuts: …`) at ipinapakita ng buod ang `shortcut: ccwork (not created; see above)`.

Puwedeng umabot nang ilang minuto ang login, at baka may ibang ccam process na mag-save ng config habang naghihintay. Kaya binabasa ulit ng mga command na naghihintay ang `config.toml`, inilalapat ang iisang pagbabago (dito, ang email) at nagse-save (`env.update`). Binabasa ulit ng `ccam remove` ang config pagkatapos ng tanong nito sa parehong dahilan.

Pinapatakbo ng `ccam login` ang hakbang 5 at 6 para sa isang umiiral na account. Pahiwatig lang para sa login page ang email para sa Claude Code, kaya puwede kang mapunta sa pag-sign in sa ibang Claude account kaysa dati; ikinukumpara ng ccam ang bagong login sa `[profile.login]` at sinasabi iyon (`ccam: this login is other@company.com, but work was dev@company.com before.`). Pinapanatili nito ang bagong login at sinasabi kung paano ito babawiin: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

Inaalis din ng login sa ibang Claude account ang [tala ng limit](/features/usage#the-limit-record) ng account, dahil sa ibang account ang mga limit na iyon. Kapag naka-log in ang isa pang ccam account sa parehong Claude account, nagbababala ang ccam na iisa ang usage limit ng dalawa.

Isa-isang ina-log in ng `ccam login --expired` ang bawat account na hindi pa kailanman naka-log in o nag-expire na ang login: nag-expire ang refresh token nito, o binura ng Claude Code ang login matapos itong tanggihan ng server. Bago bumukas ang unang browser, sinusuri nito kung masa-save ang `config.toml`.

Pinapatakbo ng `ccam logout` ang `claude auth logout` sa account at binubura ang naka-save na email; nananatili ang tala ng login, kaya sinusuri laban dito ang susunod na login. Binabawi ng Claude Code ang login at binubura ang iniimbak nito para sa account, pati ang mga sign-in sa MCP server. Kaya binabalaan ka muna ng ccam at nagtatanong ito (nilalaktawan ng `--yes` ang tanong) kapag:

- may hawak na mga MCP sign-in ang login;
- tumatakbo ang account (masa-sign out din ang mga session nito);
- hawak ng folder ng ibang account ang parehong login (isang kinopyang folder), na binabawi rin ng logout.

Halimbawa:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

Para mag-log in bilang ibang Claude account, sapat na ang `ccam login <name>` lang, at pinapanatili nito ang mga MCP sign-in.

Nagdadagdag ang `ccam list` ng column na `Organization` kapag kabilang ang mga account sa higit sa isang organisasyon (`personal` ang ipinapakita ng mga personal na plan), at may `plan` at `org` ang `ccam list --json` para sa bawat account.

### Pag-sign in sa isang remote na makina

Sa SSH, puwedeng masira ang bahagi ng `claude auth login` na nasa browser. Naghihintay ang Claude Code sa sign-in sa isang callback server sa remote na makina (`http://localhost:<port>/callback`), pero ang browser na tumatapos sa sign-in ay tumatakbo sa iyong computer, kung saan walang pinupuntahan ang address na iyon at nagpapakita ang browser ng error gaya ng `ERR_CONNECTION_REFUSED`. Dala pa rin ng address ang one-time code, kaya kaya itong ihatid ng ccam:

<!-- motion:accounts-relay -->

Hindi ibinibigay ng ccam ang iyong terminal sa Claude Code (`runLogin`). Ikinokonekta nito ang standard input ng Claude Code sa isang pipe at kinokopya doon ang bawat linyang tina-type mo (`relayLogin`), kaya dumarating pa rin ang code na idinikit sa prompt na `Paste code here` ng Claude Code. Kapag callback address ang isang linya, tinatanggap lang ito ng ccam kung `localhost`, `127.0.0.1` o `[::1]` ang host, may port, `/callback` ang path at parehong naroon ang `code` at `state`; puwedeng wala ang `http://` (`callbackURL`). Pagkatapos ay hinihiling ng ccam ang address na iyon sa server na may 10 segundong timeout at nang hindi sumusunod sa mga redirect, gaya ng gagawin sana ng browser (`deliverCallback`). Umaabot ang request sa callback server ng naghihintay na `claude auth login` (ang putol-putol na linya sa diagram). Natatanggap ng Claude Code ang code at sine-save ang login, sa Linux sa `.credentials.json` sa folder ng account. Hindi kailanman mapapakontak ng isang idinikit na linya ang ccam sa ibang host.

Ipi-print ito ng ccam bago magsimula ang login:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

Kapag nagtagumpay, ipi-print nito ang `→ Passed the sign-in to Claude Code on this machine.`

Habang naghihintay ng input ang login sa isang terminal, pinapatay ng ccam ang echo at line mode ng terminal at nagpapakita ng `*` para sa bawat character na tina-type mo, kaya hindi lumalabas sa screen ang idinikit na code o callback address. Ang ccam mismo ang gumagawa ng line editing (`readLine`): ipinapadala ng `Enter` ang linya, binubura ng `Backspace` ang isang character, at tinatapos ng `Ctrl+D` sa walang lamang linya ang input. Ibinabalik sa dati ang terminal kapag natapos ang login, at gayundin sa `Ctrl+C`, `SIGTERM` o `SIGHUP`, na tumatapos pa rin sa ccam. Kapag walang terminal sa standard input at standard error, kinokopya ang mga linya habang dumarating.

### Pagpapalit ng pangalan, pag-alis at pag-import

Kumikilos ang `ccam rename` at `ccam remove` sa folder ng account at sa mga iniimbak ng ccam sa ibang lugar sa ilalim ng pangalan ng account. Inililista ng ccam ang mga path na iyon sa `nameState`:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- ang mga usage cache na `~/.local/state/ccam/usage/<name>.json` at `<name>.statusline.json`;
- ang tala ng limit na `<name>.limit.json`;
- ang mga tala ng session sa `~/.local/state/ccam/handoff/<name>/`.

Hinaharap ng `ccam remove work` ang mga ito sa ganitong pagkakasunod-sunod:

<!-- motion:accounts-remove -->

1. Tumatanggi muna ang `ccam remove` kung tumatakbo ang account, saka nagtatanong, at sumusuri ulit pagkatapos mong sumagot (nahuhuli nito ang session na sinimulan habang bukas ang tanong).
2. Sinusuri nito kung masa-save ang `config.toml` bago ito magbago ng anuman.
3. Kung nagshe-share ng mga extension ang account, isinusulat ulit ng ccam ang `installed_plugins.json` at `known_marketplaces.json` ng group (`RewritePluginPaths`), dahil nag-iimbak ang Claude Code ng mga absolute na plugin path na puwedeng dumaan sa folder ng account na ito.
4. Susunod na ina-log out ang account na pinamamahalaan ng ccam (ang account na ang folder ay ginawa ng ccam, hindi ang nirehistro gamit ang `ccam import`): pinapatakbo ang `claude auth logout` sa folder nito, kaya walang gumaganang login na naghihintay sa trash (o, sa macOS, sa isang Keychain item na ipinangalan sa folder, na mamanahin ng susunod na account na may ganoong pangalan). Nilalaktawan ito para sa API key account at kapag hawak ng folder ng ibang account ang parehong login, na babawiin din ng logout. Iniuulat ang logout na pumalya, at nagpapatuloy ang pag-alis.
5. Inililipat ng ccam ang folder sa `~/.local/share/ccam/trash/<name>-<time>/` (`profile.Trash`), inaalis ang account sa `config.toml`, binubura ang mga path ng `nameState` para hindi ito mamana ng susunod na account na may parehong pangalan, at inaalis ang link ng shortcut: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

Kapag hindi pa rin masa-save ang `config.toml`, ibinabalik ang folder mula sa trash. Para maibalik ang inalis na account, ilipat pabalik ang folder mula sa trash at mag-log in ulit.

Pinapalitan ng `ccam rename` ang pangalan ng folder at isinusulat ulit ang mga plugin path sa `plugins/` nito mula sa lumang folder papunta sa bago. Inililipat nito ang mga path ng `nameState` sa bagong pangalan at pinapanatiling default ang account kung default ito dati. Ina-update din nito ang mga tala ng session, kaya mahahanap pa rin ang mga ito ng `ccam next --continue`. Pinapanatili ng shortcut command ang pangalan nito. Kapag pumalya ang pagsulat ulit ng plugin, o hindi masa-save ang `config.toml`, ibinabalik ang lumang pangalan ng folder.

Sa macOS, at sa Windows kapag iniimbak ng Claude Code ang mga login sa Credential Manager, pinapanatili ng folder mismo ang lumang pangalan nito, dahil naka-store ang login at ang mga MCP sign-in sa ilalim ng path ng folder at mala-log out ang inilipat na folder. Itinatala na ng config ang folder bilang `folder` katabi ng bagong pangalan, at sinasabi ng ccam ang `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

Nirerehistro ng `ccam import` ang isang folder kung nasaan ito, bilang `dir` sa `config.toml`. Hindi nito kailanman inililipat o kinokopya ito: kapag kinopya ang refresh token, magla-log out sa isa't isa ang mga kopya.

Kapag walang folder argument, hinahanap nito ang mga hindi pa nakarehistrong `~/.claude` at `~/.claude-*` folder at nagtatanong tungkol sa bawat isa. Kinukuha ang mga pangalan mula sa folder (`main` ang `~/.claude`, `work` ang `~/.claude-work`, `acct2` ang `~/.claude-2`). Kung hindi makakuha ang ccam ng pangalan, o kinuha na ang pangalan, inililista nito ang folder kasama ang `--name` command na patatakbuhin. Iniuulat bilang `Skipped <path>: <reason>` ang nahanap na folder na pumalya at ini-import ang iba.

Itinatanong muna ng ccam ang lahat ng tanong nito, saka binabasa ulit ang `config.toml` at nirerehistro ang mga folder, kaya napapanatili ang mga pagbabagong ini-save ng ibang ccam command sa pagitan. Lumalabas ang mga linyang `✓` pagkatapos ng pag-save.

Kapag walang `CLAUDE_CONFIG_DIR`, iniimbak ng Claude Code ang `.claude.json` sa iyong home folder katabi ng `~/.claude`, kaya kapag in-import mo ang `~/.claude` at wala itong sariling `.claude.json`, gagawa ang ccam ng link na `~/.claude/.claude.json` papunta sa `~/.claude.json`. Sa ganitong paraan iisang account ang nakikita ng karaniwang `claude` at ng `ccam main`.

Kung may sariling `.claude.json` ang `~/.claude` katabi ng `~/.claude.json`, hinahayaan ng ccam ang dalawa at sinasabi ito. Gagamitin ng `ccam main` ang nasa loob at ng karaniwang `claude` ang nasa labas, na naghihiwalay sa mga detalye ng login at mga MCP server nila. Para iisa na lang, palitan ng link ang file na nasa loob (iniuulat din ito ng [doctor](/features/doctor)).

Kapag hawak ng folder ang parehong login ng isang nakarehistrong account (isang kinopyang folder), sinasabi ito ng ccam: ang unang mag-renew nito ang magla-log out sa isa pa, kaya patakbuhin ang `ccam login <name>` bago ito gamitin.

Hindi kasapi sa anumang group ang mga in-import na account hanggang idagdag mo sila sa isa (`ccam group join <group> <name>`), at walang shortcut na ginagawa.

Nagpi-print ang `ccam completion bash|zsh` ng completion script na nagko-complete rin ng mga pangalan ng account.

## Kaligtasan at mga espesyal na kaso

- Tumatanggi ang `ccam rename` at `ccam remove` habang tumatakbo ang account (`refuseRunning`), na may ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``, dahil patuloy na nagsusulat ang session sa folder. Saklaw ng pagsusuri ang Claude Code paano man ito sinimulan: sa pamamagitan ng ccam, karaniwang `claude`, isang IDE, mga background session at ang daemon nila.
- Tumatanggi ang `ccam rename`, `ccam remove` at `ccam logout` habang nire-renew o binabago ng ibang ccam process ang login na iyon: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- Habang may isinasagawang pagpapalit ng group sa ibang ccam process, tumatanggi ang `ccam rename` at `ccam remove` na may `another ccam process is changing the groups right now; try again when it is done`; pagkatapos ng pagpapalit na huminto sa kalagitnaan, na may ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.
- Sa isang in-import na account, pinapanatili ng `ccam remove` ang folder at ang login nito, para sa karaniwang `claude`. Kung nasa mga group ang account, aalis muna ito sa mga iyon gaya ng ginagawa ng `ccam group leave`: nagiging kopya ng sa group ang bawat shared na item sa folder, at napupunta ang mga MCP server ng group sa `.claude.json` nito. Saka inaalis ang rehistro nito: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. Inililipat sa `<folder>/.ccam-backup/mcp-<time>` ang mga backup ng `.claude.json` nito na ginawa ng ccam; kung hindi mailipat ang mga ito, nananatili ang mga ito sa `~/.local/state/ccam/mcp/<name>/backups` at sinasabi ito ng mensahe. Hindi inililipat sa trash ang iyong `~/.claude`.
- Hindi puwedeng iisa ang folder ng dalawang account: pumapalya sa pag-load ang `config.toml` na dalawang beses na nagpapangalan sa iisang folder.
- Tinatanggihan ng `ccam add` ang pangalan na ang folder ay mayroon na sa `profiles/`, at itinuturo ang `ccam import`. Sa macOS, humihinto rin ito kapag may hawak pa ang Keychain na login o mga MCP sign-in para sa path ng bagong folder, na iniwan ng mas lumang bersyon ng ccam, at pinapangalanan ang `security delete-generic-password` command na nag-aalis sa mga ito. Kapag hindi masa-save ang `config.toml`, inaalis ulit ang bagong folder.
- Tinatanggihan ng `ccam import` ang folder na nakarehistro na, ang path na hindi folder, at ang pangalan na gamit na. Nagbababala rin ang ccam kapag kopya ng `~/.claude` ang folder (parehong login), kahit hindi nakarehistro ang `~/.claude` mismo.
- Nagtatanong muna ang `ccam logout` kapag tumatakbo ang account; hindi ito sinusuri ng `ccam login`.
- Kapag hindi maabot ng relay ang Claude Code, halimbawa pagkatapos matapos ang login, ipi-print nito ang `ccam: could not pass the sign-in to Claude Code:` kasama ang dahilan. Kadalasan, ang HTTP error ay nangangahulugang nag-expire na ang code; simulan ulit ang login.

## Kaugnay

- [Pagsisimula ng Claude Code](/features/launch)
- [Mga shortcut command](/features/shortcuts)
- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Maraming Claude Code account sa iisang computer](/guides/multiple-claude-code-accounts)
