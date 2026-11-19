---
title: Sarrafa asusai
description: Ƙara, shiga, lissafa, sake suna, cire da shigo da asusun Claude Code da ccam add, ccam login, ccam list, ccam rename, ccam remove da ccam import.
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

Kowane asusun ccam foldar saitin Claude Code ce ta daban mai nata shiga. Yi amfani da waɗannan umarnai don ƙara asusai, shigar da su, lissafa su, sake musu suna da cire su. Suna kuma yin rajistar foldojin da kake da su tuni.

## Amfani

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

Duk inda umarni ke karɓar `<name|number>`, za ka iya ba da sunan asusun ko lambar ramin (slot) sa (ginshiƙin `#` a `ccam list`).

`ccam use` yana canza wane asusu ne tsoho kawai (`default` a `config.toml`).

`ccam add` a terminal ba tare da zaɓuɓɓuka ba yana buɗe mataimaki (wizard). Yana bin waɗannan matakai:

1. sunan asusun, sai dai idan ka bayar da shi;
2. yadda asusun zai shiga: shigar Claude yanzu, API key ko shigar Claude daga baya;
3. rukunonin da za a shiga, tare da alama a kan `default` (idan babu rukunoni tukuna, yana nuna abin da rukunin `default` zai raba);
4. `Create a shortcut command for work? [y/N]`;
5. taƙaitawa don ka tabbatar.

`Esc` yana komawa mataki ɗaya baya. Duk wani daga `--group`, `--no-group`, `--shortcut` ko `--no-login`, ko gudanarwa ba tare da terminal ba (misali, daga bututu (pipe)), yana tsallake mataimakin. Sannan dole a ba da suna, kuma ba a ƙirƙirar gajeren umarni sai idan `--shortcut` ya ambaci ɗaya.

Ba tare da `--group` ba, sabon asusu yana shiga rukunin `default`. Asusun farko ne yake ƙirƙirar sa, yana raba sabobin MCP, saituna, `CLAUDE.md` da extensions; duba [rabawa da rukunoni](/features/share). A Windows, ana iya haɗa fayiloli ne kawai a Developer Mode. Ba tare da shi ba, asusun farko yana ƙirƙirar `default` da `mcp` da `extensions` kawai kuma yana faɗin haka (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), kuma ccam yana ƙin `--group` tare da rukunin da ke raba saituna, memory ko tarihi.

`ccam list` yana nuna ramin, suna, gajeren umarni, imel, rukunoni da amfani na ƙarshe, tare da `*` a gaban asusun tsoho. Imel ɗin shi ne wanda foldar ke shiga da shi yanzu (daga `.claude.json` nata); idan ya bambanta da wanda ccam ya rubuta, wanda aka rubuta yana biyo baya a matsayin `(was …)`:

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

Da `--json`, kowane asusu abu ne mai `slot`, `shortcut`, `name`, `email` (na shigar yanzu, kamar a tebur), `default`, `dir`, `last_used` (ba ya nan idan ba a taɓa fara asusun ba), `groups` da `share`. `share` yana danganta kowane abu (`mcp`, `settings`, `memory`, `extensions`, `history`) da rukunin da yake fitowa daga gare shi, ko da `"local"`.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `ccam add --no-login` | Ƙirƙiri asusun ba tare da shiga ba; shiga daga baya da `ccam login <name>` | a kashe |
| `ccam add --group <group>` | Shiga wannan rukunin maimakon `default`. Ana iya maimaitawa; ana ƙirƙirar `default` idan babu shi | `default` |
| `ccam add --no-group` | Kada a shiga kowane rukuni: asusun ba ya raba komai | a kashe |
| `ccam add --shortcut <command>` | Kuma ƙirƙiri wannan gajeren umarni a `~/.local/bin` | babu |
| `ccam add --api --key-from <source>` | Asusun da ke amfani da API key maimakon shiga; duba [Asusan API key](/features/api) | a kashe |
| `ccam list --json` | Buga asusan a matsayin JSON | a kashe |
| `ccam remove -y`, `--yes` | Cire ba tare da tambaya ba | a kashe |
| `ccam logout -y`, `--yes` | Fita daga shiga ba tare da tambaya ba | a kashe |
| `ccam import --name <name>` | Sunan asusu na foldar da ka bayar; ana ƙin shi ba tare da folda ba (`--name needs a folder: ccam import <folder> --name ops`) | daga sunan foldar |
| `ccam import -y`, `--yes` | Shigo da kowace foldar da aka samo ba tare da tambaya ba | a kashe |
| `--dry-run` | `ccam list` da `ccam completion` suna tallafa masa, domin ba sa canza komai. `add`, `login`, `logout`, `use`, `rename`, `remove` da `import` suna ƙin shi, da `` --dry-run is not supported by `ccam add`; nothing was changed `` (tare da sunan umarnin) | a kashe |

`ccam login` (banda `--expired`), `ccam use` da `ccam rename` ba su da zaɓuɓɓukan kansu. `ccam completion` yana karɓar `bash`, `zsh` ko `powershell`; duk wani abu dabam yana gazawa da `only bash, zsh and powershell are supported`.

Sunayen asusai haruffa 1 zuwa 32 ne na `a-z`, `0-9`, `-` da `_`, suna farawa da harafi ko lamba. Suna ba zai iya zama lambobi kawai ba, domin `ccam 2` yana nufin ramin 2. Haka kuma suna ba zai iya zama umarnin ccam ba (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), don haka `ccam <name>` ba ya taɓa yin shubuha. Jerin shi ne `config.Reserved` a `internal/config`.

## Yadda yake aiki

Ƙara asusu yana ƙirƙirar foldarsa, yana shiga rukunoninsa, yana ajiye saitin, yana shiga kuma yana ƙirƙirar gajeren umarni, a wannan tsari. Idan ka gudanar da `ccam add work --shortcut ccwork`, ccam yana bin waɗannan matakai (`addProfile` a `internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. ccam yana duba sunan (`config.ValidateName`) da gajeren umarnin (`shim.Check`) kafin ya ƙirƙiri komai.
2. ccam yana ƙirƙirar `~/.local/share/ccam/profiles/work/` da mode `0700` kuma yana ba asusun mafi ƙanƙantar lambar rami da ke a kwance, don haka ana sake amfani da lambar asusun da aka cire (`profile.Add`, `Config.NextSlot`). Asusun farko yana zama tsoho.
3. Asusun yana shiga rukunoninsa (`default` sai dai idan `--group` ko `--no-group` ya faɗi wani abu dabam; asusun farko yana ƙirƙirar `default`), kuma ccam yana haɗa sabuwar foldar mara komai da su, misali `settings.json` zuwa `groups/default/settings.json` (`applyShare`, ta `share.Build` sannan `Apply`). Idan wannan ya gaza, ana cire sabuwar foldar, domin foldar da ta rage za ta hana kowace sake gwadawa da `folder already exists`.
4. Ana ajiye asusun a `~/.config/ccam/config.toml` ta `fsutil.AtomicWrite`: fayil na wucin gadi, `fsync`, sannan `rename`, don haka faɗuwar shiri ba ta taɓa barin saitin da aka rubuta rabi ba.
5. ccam yana fara `claude auth login` tare da `CLAUDE_CONFIG_DIR` da aka saita zuwa foldar asusun, da `--email` tare da imel na ƙarshe na asusun don a cike shafin shiga (`runLogin`). Idan shigar ta gaza, asusun yana nan, kuma ccam yana buga kuskuren kuma yana gaya maka ka gudanar da `ccam login work` daga baya.
6. ccam yana karanta `oauthAccount` daga `profiles/work/.claude.json` (`profile.ReadIdentity`): imel, ID na asusu da na ƙungiya (organization), sunan ƙungiyar da tsarin biyan kuɗi (`claude_max`, `claude_team` …). Wannan fayil zai iya kai megabytes da dama, don haka ccam yana neman maɓalli ɗaya kawai maimakon ya karanta dukkan fayil ɗin (`TopLevelValue`). ccam yana ajiye waɗannan bayanai a `[profile.login]` a cikin `config.toml`.
7. ccam yana haɗa `~/.local/bin/ccwork` da binary na ccam, wanda ke neman `argv[0]` a `config.toml` lokacin da ya fara, kuma yana rubuta mahaɗin a `~/.local/state/ccam/shortcuts.json` (`shim.Install`). Idan ba za a iya ƙirƙirar mahaɗin ba, ana ƙara asusun duk da haka, ccam yana faɗin dalili (`Could not create shortcuts: …`) kuma taƙaitawar tana nuna `shortcut: ccwork (not created; see above)`.

Shiga na iya ɗaukar mintuna, kuma wani process na ccam na iya ajiye saitin a wannan lokacin. Don haka umarnin da ke jira suna sake loda `config.toml`, suna amfani da canji ɗaya (a nan, imel ɗin) kuma suna ajiyewa (`env.update`). `ccam remove` yana sake lodawa bayan tambayarsa saboda dalili ɗaya.

`ccam login` yana gudanar da mataki na 5 da na 6 ga asusun da ke akwai. Claude Code yana ɗaukar imel ɗin ne kawai a matsayin alama ga shafin shiga, don haka za ka iya ƙarewa a shiga cikin wani asusun Claude dabam da na da; ccam yana kwatanta sabuwar shigar da `[profile.login]` kuma yana faɗin haka (`ccam: this login is other@company.com, but work was dev@company.com before.`). Yana riƙe sabuwar shigar kuma yana gaya maka yadda za ka soke ta: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

Shiga cikin wani asusun Claude dabam yana kuma cire [rikodin iyaka](/features/usage#the-limit-record) na asusun, tunda waɗannan iyakokin na ɗayan asusun ne. Idan wani asusun ccam ya shiga asusun Claude iri ɗaya, ccam yana gargaɗi cewa su biyun suna raba iyakokin amfaninsu.

`ccam login --expired` yana shiga, ɗaya bayan ɗaya, kowane asusu da ba a taɓa shiga da shi ba ko wanda shigarsa ta ƙare: refresh token ɗinsa ya ƙare, ko Claude Code ya share shigar bayan sabar ta ƙi ta. Kafin burauza ta farko ta buɗe, yana duba cewa za a iya ajiye `config.toml`.

`ccam logout` yana gudanar da `claude auth logout` a cikin asusun kuma yana share imel ɗin da aka ajiye; bayanin shigar yana nan, don haka ana duba shiga ta gaba da shi. Claude Code yana soke shigar kuma yana goge abin da yake ajiyewa don asusun, har da shigar sabobin MCP. Don haka ccam yana gargaɗe ka kuma yana fara tambaya (`--yes` yana tsallake tambayar) idan:

- shigar tana riƙe da shigar MCP;
- asusun yana aiki (ana fitar da zamansa ma);
- foldar wani asusu tana riƙe da shiga iri ɗaya (foldar da aka kwafa), wadda fitar ke soke ita ma.

Misali:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

Don shiga a matsayin wani asusun Claude dabam, `ccam login <name>` kaɗai ya isa, kuma yana riƙe shigar MCP.

`ccam list` yana ƙara ginshiƙin `Organization` idan asusan na ƙungiyoyi fiye da ɗaya ne (tsare-tsaren kai suna nuna `personal`), kuma `ccam list --json` yana da `plan` da `org` ga kowane asusu.

### Shiga a kwamfutar nesa

Ta SSH, ɓangaren burauza na `claude auth login` zai iya lalacewa. Claude Code yana jiran shigar a sabar callback a kan kwamfutar nesa (`http://localhost:<port>/callback`), amma burauzar da ke kammala shigar tana aiki a kwamfutarka, inda wannan adireshin ba ya kai ko'ina kuma burauzar tana nuna kuskure kamar `ERR_CONNECTION_REFUSED`. Adireshin har yanzu yana ɗauke da lambar sau ɗaya, don haka ccam zai iya isar da ita:

<!-- motion:accounts-relay -->

ccam ba ya miƙa terminal ɗinka ga Claude Code (`runLogin`). Yana haɗa standard input na Claude Code da bututu (pipe), kuma yana kwafar kowane layi da ka rubuta zuwa cikinsa (`relayLogin`), don haka lambar da aka liƙa a tambayar `Paste code here` ta Claude Code har yanzu tana isa. Idan layi adireshin callback ne, ccam yana karɓar sa ne kawai idan host ɗin `localhost`, `127.0.0.1` ko `[::1]` ne, an ba da port, hanyar `/callback` ce kuma duka `code` da `state` suna nan; ana iya barin `http://` (`callbackURL`). Sannan ccam yana neman wannan adireshin a kan sabar da iyakar lokaci ta daƙiƙa 10 kuma ba tare da bin redirects ba, kamar yadda burauzar za ta yi (`deliverCallback`). Buƙatar tana isa sabar callback na `claude auth login` da ke jira (layin da aka yanke a zanen). Claude Code yana karɓar lambar kuma yana ajiye shigar, a Linux a `.credentials.json` cikin foldar asusun. Layin da aka liƙa ba zai taɓa sa ccam ya tuntuɓi wani host ba.

ccam yana buga wannan kafin shigar ta fara:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

Idan ya yi nasara yana buga `→ Passed the sign-in to Claude Code on this machine.`

Yayin da shigar ke jiran shigarwa a terminal, ccam yana kashe echo da yanayin layi (line mode) na terminal kuma yana nuna `*` ga kowane harafi da ka rubuta, don haka lambar da aka liƙa ko adireshin callback ba ya bayyana a allo. ccam yana gyara layin da kansa (`readLine`): `Enter` yana aika layin, `Backspace` yana goge harafi ɗaya, kuma `Ctrl+D` a layi mara komai yana kawo ƙarshen shigarwa. Ana mayar da terminal yadda yake idan shigar ta ƙare, haka kuma a `Ctrl+C`, `SIGTERM` ko `SIGHUP`, waɗanda har yanzu suke ƙare ccam. Idan babu terminal a standard input da standard error, ana kwafar layuka yadda suke zuwa.

### Sake suna, cirewa da shigowa

`ccam rename` da `ccam remove` suna aiki a kan foldar asusun da kuma abin da ccam ke ajiyewa da sunan asusun a wasu wurare. ccam yana lissafa waɗannan hanyoyi a `nameState`:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- cache na amfani `~/.local/state/ccam/usage/<name>.json` da `<name>.statusline.json`;
- rikodin iyaka `<name>.limit.json`;
- bayanan zama a `~/.local/state/ccam/handoff/<name>/`.

`ccam remove work` yana sarrafa su a wannan tsari:

<!-- motion:accounts-remove -->

1. `ccam remove` da farko yana ƙi idan asusun yana aiki, sannan yana tambaya, kuma yana sake dubawa bayan ka amsa (wannan yana kama zaman da aka fara yayin da tambayar ke buɗe).
2. Yana duba cewa za a iya ajiye `config.toml` kafin ya canza komai.
3. Idan asusun yana raba extensions, ccam yana sake rubuta `installed_plugins.json` da `known_marketplaces.json` na rukunin (`RewritePluginPaths`), domin Claude Code yana ajiye cikakkun hanyoyin plugin waɗanda za su iya ratsa foldar wannan asusun.
4. Sannan ana fitar da asusun da ccam ke sarrafawa (wanda ccam ya ƙirƙiri foldarsa, ba wanda aka yi rajistarsa da `ccam import` ba) daga shiga: `claude auth logout` yana gudana a cikin foldarsa, don kada wata shiga mai aiki ta jira a shara (ko, a macOS, a cikin abin Keychain mai sunan foldar, wanda asusun da za a yi daga baya mai wannan suna zai gada). Ana tsallake hakan ga asusun API key da kuma idan foldar wani asusu tana riƙe da shiga iri ɗaya, wadda fitar za ta soke ita ma. Ana ba da rahoton fitar da ta gaza, kuma cirewar tana ci gaba.
5. ccam yana matsar da foldar zuwa `~/.local/share/ccam/trash/<name>-<time>/` (`profile.Trash`), asusun yana fita daga `config.toml`, ana goge hanyoyin `nameState` don kada asusun da za a yi daga baya mai suna iri ɗaya ya gaje su, kuma ana cire mahaɗin gajeren umarnin: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

Idan daga ƙarshe ba za a iya ajiye `config.toml` ba, foldar tana dawowa daga shara. Don maido da asusun da aka cire, mayar da foldar daga shara kuma ka sake shiga.

`ccam rename` yana sake sunan foldar kuma yana sake rubuta hanyoyin plugin a cikin `plugins/` ɗinta daga tsohuwar foldar zuwa sabuwa. Yana matsar da hanyoyin `nameState` zuwa sabon sunan kuma yana riƙe asusun a matsayin tsoho idan shi ne tsoho. Yana kuma sabunta bayanan zama, don `ccam next --continue` ya ci gaba da samun su. Gajeren umarnin yana riƙe sunansa. Idan sake rubuta plugin ɗin ya gaza, ko ba za a iya ajiye `config.toml` ba, ana mayar da sunan foldar.

A macOS, da kuma a Windows idan Claude Code yana ajiye bayanan shiga a Credential Manager, foldar kanta tana riƙe tsohon sunanta, domin ana ajiye shigar da shigar MCP ƙarƙashin hanyar foldar, kuma foldar da aka matsar za a fitar da ita daga shiga. Sannan config yana rubuta foldar a matsayin `folder` kusa da sabon sunan, kuma ccam yana cewa `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

`ccam import` yana yin rajistar folda a inda take, a matsayin `dir` a `config.toml`. Ba ya taɓa matsar ko kwafar ta: refresh token da aka kwafa zai sa kwafin su fitar da juna.

Ba tare da argument na folda ba, yana neman foldojin `~/.claude` da `~/.claude-*` da ba a yi wa rajista ba kuma yana tambaya game da kowace. Sunaye suna fitowa daga foldar (`~/.claude` shi ne `main`, `~/.claude-work` shi ne `work`, `~/.claude-2` shi ne `acct2`). Idan ccam ba zai iya samun suna ba, ko an riga an ɗauki sunan, yana lissafa foldar tare da umarnin `--name` da za a gudanar. Ana ba da rahoton foldar da aka samo wadda ta gaza a matsayin `Skipped <path>: <reason>` kuma ana shigo da sauran.

ccam yana fara yin dukkan tambayoyinsa, sannan yana sake loda `config.toml` kuma yana yin rajistar foldojin, don haka canje-canjen da wasu umarnan ccam suka ajiye a wannan lokacin suna nan. Layukan `✓` suna zuwa bayan ajiyewa.

Ba tare da `CLAUDE_CONFIG_DIR` ba, Claude Code yana ajiye `.claude.json` a foldar gidanka kusa da `~/.claude`, don haka idan ka shigo da `~/.claude` kuma ba shi da `.claude.json` nasa, ccam yana ƙirƙirar mahaɗin `~/.claude/.claude.json` zuwa `~/.claude.json`. Sannan `claude` na yau da kullum da `ccam main` suna ganin asusu ɗaya.

Idan `~/.claude` yana da `.claude.json` nasa kusa da `~/.claude.json`, ccam yana barin dukansu kuma yana faɗin haka. `ccam main` zai yi amfani da na ciki, `claude` na yau da kullum kuma na waje, kowanne da bayanan shigarsa da sabobin MCP nasa daban. Don riƙe ɗaya, maye gurbin fayil na ciki da mahaɗi ([doctor](/features/doctor) shi ma yana ba da rahoton sa).

Idan foldar tana riƙe da shiga iri ɗaya da asusun da aka yi wa rajista (foldar da aka kwafa), ccam yana faɗin haka: duk wanda ya fara sabunta ta yana fitar da ɗayan, don haka gudanar da `ccam login <name>` kafin ka yi amfani da ita.

Asusan da aka shigo da su ba sa cikin kowane rukuni har sai ka ƙara su a ɗaya (`ccam group join <group> <name>`), kuma ba a ƙirƙirar gajeren umarni.

`ccam completion bash|zsh` yana buga rubutun kammalawa (completion) wanda kuma ke kammala sunayen asusai.

## Tsaro da yanayi na musamman

- `ccam rename` da `ccam remove` suna ƙi yayin da asusun ke aiki (`refuseRunning`), da ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``, domin zaman yana ci gaba da rubutu a cikin foldar. Binciken yana rufe Claude Code ko ta yaya aka fara shi: ta ccam, `claude` na yau da kullum, IDE, zaman bayan fage da daemon ɗinsu.
- `ccam rename`, `ccam remove` da `ccam logout` suna ƙi yayin da wani process na ccam ke sabunta ko canza wannan shigar: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- Yayin da canjin rukuni ke gudana a wani process na ccam, `ccam rename` da `ccam remove` suna ƙi da `another ccam process is changing the groups right now; try again when it is done`; bayan wanda ya tsaya rabi, da ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.
- `ccam remove` a kan asusun da aka shigo da shi yana riƙe foldar da shigarta, don `claude` na yau da kullum. Idan asusun yana cikin rukunoni, da farko yana fita daga cikinsu kamar yadda `ccam group leave` ke yi: kowane abu da aka raba yana zama kwafin na rukunin a cikin foldar, kuma sabobin MCP na rukunin suna shiga `.claude.json` nasa. Sannan ana soke rajistarsa: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. Kwafin ajiyar `.claude.json` nasa da ccam ya ƙirƙira suna matsawa zuwa `<folder>/.ccam-backup/mcp-<time>`; idan ba za a iya matsar da su ba, suna zama a `~/.local/state/ccam/mcp/<name>/backups` kuma saƙon yana faɗin haka. Ba a matsar da `~/.claude` naka zuwa shara ba.
- Asusai biyu ba za su iya samun folda ɗaya ba: `config.toml` da ya ambaci folda ɗaya sau biyu yana kasa lodawa.
- `ccam add` yana ƙin sunan da foldarsa ta riga ta kasance a `profiles/`, kuma yana nuna `ccam import`. A macOS yana kuma tsayawa idan Keychain har yanzu yana riƙe da shiga ko shigar MCP don hanyar sabuwar foldar, wanda tsohuwar sigar ccam ta bari, kuma yana ambaton umarnin `security delete-generic-password` da ke cire su. Idan ba za a iya ajiye `config.toml` ba, ana sake cire sabuwar foldar.
- `ccam import` yana ƙin foldar da aka riga aka yi wa rajista, hanyar da ba folda ba ce, da sunan da aka riga aka ɗauka. ccam yana kuma gargaɗi idan foldar kwafin `~/.claude` ce (shiga iri ɗaya), ko da ba a yi wa `~/.claude` kansa rajista ba.
- `ccam logout` yana fara tambaya idan asusun yana aiki; `ccam login` ba ya dubawa.
- Idan relay ɗin ba zai iya kai wa Claude Code ba, misali bayan shigar ta ƙare, yana buga `ccam: could not pass the sign-in to Claude Code:` tare da dalilin. Kuskuren HTTP yawanci yana nufin lambar ta ƙare; sake fara shigar.

## Masu alaƙa

- [Fara Claude Code](/features/launch)
- [Gajerun umarni](/features/shortcuts)
- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Amfani da asusun Claude Code da dama a kwamfuta ɗaya](/guides/multiple-claude-code-accounts)
