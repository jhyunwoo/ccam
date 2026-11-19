---
title: Kusimamia akaunti
description: Ongeza, ingia, orodhesha, badilisha jina, ondoa na ingiza akaunti za Claude Code kwa ccam add, ccam login, ccam list, ccam rename, ccam remove na ccam import.
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

Kila akaunti ya ccam ni folda tofauti ya usanidi wa Claude Code yenye taarifa zake za kuingia. Tumia amri hizi kuongeza akaunti, kuziingiza, kuziorodhesha, kubadilisha majina yake na kuziondoa. Pia husajili folda ambazo tayari unazo.

## Matumizi

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

Popote amri inapopokea `<name|number>`, unaweza kutoa jina la akaunti au namba yake (slot) (safu ya `#` katika `ccam list`).

`ccam use` hubadilisha tu akaunti ipi ni chaguo-msingi (`default` katika `config.toml`).

`ccam add` kwenye terminal bila chaguo hufungua mchawi. Hupitia hatua hizi:

1. jina la akaunti, isipokuwa uwe umelitoa;
2. jinsi akaunti inavyoingia: kuingia kwa Claude sasa, API key au kuingia kwa Claude baadaye;
3. vikundi vya kujiunga, huku `default` ikiwa imeteuliwa (ikiwa bado hakuna vikundi, huonyesha kile ambacho kikundi cha `default` kitashiriki);
4. `Create a shortcut command for work? [y/N]`;
5. muhtasari wa kuthibitisha.

`Esc` hurudi hatua moja nyuma. Chaguo lolote kati ya `--group`, `--no-group`, `--shortcut` au `--no-login`, au kuendesha bila terminal (kwa mfano, kutoka kwenye pipe), huruka mchawi. Hapo jina linahitajika, na hakuna njia fupi inayoundwa isipokuwa `--shortcut` itaje moja.

Bila `--group`, akaunti mpya hujiunga na kikundi cha `default`. Akaunti ya kwanza hukiunda, kikishiriki seva za MCP, mipangilio, `CLAUDE.md` na viendelezi; tazama [kushiriki kwa vikundi](/features/share). Kwenye Windows, faili zinaweza kuunganishwa kwa viungo tu katika Developer Mode. Bila hiyo, akaunti ya kwanza huunda `default` ikiwa na `mcp` na `extensions` pekee na husema hivyo (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), na ccam hukataa `--group` yenye kikundi kinachoshiriki mipangilio, kumbukumbu (memory) au historia.

`ccam list` huonyesha namba, jina, njia fupi, barua pepe, vikundi na matumizi ya mwisho, na `*` mbele ya akaunti chaguo-msingi. Barua pepe ni ile ambayo folda imeingia nayo sasa (kutoka `.claude.json` yake); ikitofautiana na ile ambayo ccam ilirekodi, iliyorekodiwa hufuata kama `(was …)`:

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

Kwa `--json`, kila akaunti ni kitu chenye `slot`, `shortcut`, `name`, `email` (ya kuingia kwa sasa, kama kwenye jedwali), `default`, `dir`, `last_used` (huachwa ikiwa akaunti haijawahi kuanzishwa), `groups` na `share`. `share` huhusisha kila kipengele (`mcp`, `settings`, `memory`, `extensions`, `history`) na kikundi kinachotoka, au `"local"`.

## Chaguo

| Chaguo | Kinachofanya | Chaguo-msingi |
| --- | --- | --- |
| `ccam add --no-login` | Unda akaunti bila kuingia; ingia baadaye kwa `ccam login <name>` | imezimwa |
| `ccam add --group <group>` | Jiunge na kikundi hiki badala ya `default`. Kinaweza kurudiwa; `default` huundwa ikiwa hakipo | `default` |
| `ccam add --no-group` | Usijiunge na kikundi chochote: akaunti haishiriki chochote | imezimwa |
| `ccam add --shortcut <command>` | Pia unda amri hii ya njia fupi katika `~/.local/bin` | hakuna |
| `ccam add --api --key-from <source>` | Akaunti inayotumia API key badala ya kuingia; tazama [Akaunti za API key](/features/api) | imezimwa |
| `ccam list --json` | Chapisha akaunti kama JSON | imezimwa |
| `ccam remove -y`, `--yes` | Ondoa bila kuuliza | imezimwa |
| `ccam logout -y`, `--yes` | Toka kwenye akaunti bila kuuliza | imezimwa |
| `ccam import --name <name>` | Jina la akaunti kwa folda unayotoa; hukataliwa bila folda (`--name needs a folder: ccam import <folder> --name ops`) | hutokana na jina la folda |
| `ccam import -y`, `--yes` | Ingiza kila folda iliyopatikana bila kuuliza | imezimwa |
| `--dry-run` | Hukubaliwa na `ccam list` na `ccam completion`, ambazo hazibadilishi kitu. Hukataliwa na `add`, `login`, `logout`, `use`, `rename`, `remove` na `import`, kwa ujumbe `` --dry-run is not supported by `ccam add`; nothing was changed `` (ukitaja amri husika) | imezimwa |

`ccam login` (mbali na `--expired`), `ccam use` na `ccam rename` hazina chaguo zao wenyewe. `ccam completion` hupokea `bash`, `zsh` au `powershell`; kitu kingine chochote hushindwa kwa `only bash, zsh and powershell are supported`.

Majina ya akaunti yana herufi 1 hadi 32 za `a-z`, `0-9`, `-` na `_`, yakianza na herufi au tarakimu. Jina haliwezi kuwa tarakimu tu, kwa sababu `ccam 2` inamaanisha namba 2. Jina pia haliwezi kuwa amri ya ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), ili `ccam <name>` isiwe na utata kamwe. Orodha hiyo ni `config.Reserved` katika `internal/config`.

## Jinsi inavyofanya kazi

Kuongeza akaunti huunda folda yake, hujiunga na vikundi vyake, huhifadhi usanidi, huingia na huunda njia fupi, kwa mpangilio huo. Unapoendesha `ccam add work --shortcut ccwork`, ccam huchukua hatua hizi (`addProfile` katika `internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. ccam hukagua jina (`config.ValidateName`) na njia fupi (`shim.Check`) kabla ya kuunda kitu chochote.
2. ccam huunda `~/.local/share/ccam/profiles/work/` kwa mode `0700` na huipa akaunti namba ndogo kabisa iliyo wazi, kwa hiyo namba ya akaunti iliyoondolewa hutumiwa tena (`profile.Add`, `Config.NextSlot`). Akaunti ya kwanza huwa chaguo-msingi.
3. Akaunti hujiunga na vikundi vyake (`default` isipokuwa `--group` au `--no-group` iseme vinginevyo; akaunti ya kwanza huunda `default`), na ccam huunganisha folda mpya iliyo tupu na vikundi hivyo, kwa mfano `settings.json` kwenda `groups/default/settings.json` (`applyShare`, kwa `share.Build` na kisha `Apply`). Hili likishindwa, folda mpya huondolewa, kwa sababu folda iliyobaki ingezuia kila jaribio jipya kwa `folder already exists`.
4. Akaunti huhifadhiwa katika `~/.config/ccam/config.toml` kupitia `fsutil.AtomicWrite`: faili ya muda, `fsync`, kisha `rename`, ili hitilafu ya ghafla isiache usanidi ulioandikwa nusu.
5. ccam huanzisha `claude auth login` huku `CLAUDE_CONFIG_DIR` ikiwa imewekwa kwenye folda ya akaunti, pamoja na `--email` yenye barua pepe ya mwisho ya akaunti ili ukurasa wa kuingia ujazwe tayari (`runLogin`). Kuingia kukishindwa, akaunti hubaki, na ccam huchapisha hitilafu na kukuambia uendeshe `ccam login work` baadaye.
6. ccam husoma `oauthAccount` kutoka `profiles/work/.claude.json` (`profile.ReadIdentity`): barua pepe, IDs za akaunti na shirika, jina la shirika na mpango (`claude_max`, `claude_team` …). Faili hiyo inaweza kuwa megabaiti kadhaa, kwa hiyo ccam hutafuta ufunguo huo mmoja badala ya kuchanganua faili nzima (`TopLevelValue`). ccam huhifadhi taarifa hizi katika `[profile.login]` ndani ya `config.toml`.
7. ccam huunganisha `~/.local/bin/ccwork` na binary ya ccam, ambayo hutafuta `argv[0]` katika `config.toml` inapoanza, na hurekodi kiungo hicho katika `~/.local/state/ccam/shortcuts.json` (`shim.Install`). Kiungo kisipoweza kuundwa, akaunti bado huongezwa, ccam husema sababu (`Could not create shortcuts: …`) na muhtasari huonyesha `shortcut: ccwork (not created; see above)`.

Kuingia kunaweza kuchukua dakika kadhaa, na mchakato mwingine wa ccam unaweza kuhifadhi usanidi wakati huo. Kwa hiyo amri zinazosubiri hupakia upya `config.toml`, hutumia badiliko hilo moja (hapa, barua pepe) na kuhifadhi (`env.update`). `ccam remove` hupakia upya baada ya swali lake kwa sababu hiyo hiyo.

`ccam login` hutekeleza hatua ya 5 na 6 kwa akaunti iliyopo. Claude Code huichukua barua pepe kama kidokezo tu kwa ukurasa wa kuingia, kwa hiyo unaweza kuishia umeingia kwenye akaunti ya Claude tofauti na ile ya awali; ccam hulinganisha kuingia kupya na `[profile.login]` na husema hivyo (`ccam: this login is other@company.com, but work was dev@company.com before.`). Huhifadhi kuingia kupya na hukuambia jinsi ya kulirudisha: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

Kuingia kwenye akaunti nyingine ya Claude pia hufuta [rekodi ya kikomo](/features/usage#the-limit-record) ya akaunti, kwa kuwa vikomo hivyo vilikuwa vya akaunti ile nyingine. Akaunti nyingine ya ccam ikiwa imeingia kwenye akaunti ile ile ya Claude, ccam huonya kwamba hizo mbili hushiriki vikomo vyao vya matumizi.

`ccam login --expired` huingiza, moja baada ya nyingine, kila akaunti ambayo haijawahi kuingia au ambayo kuingia kwake kumeisha muda: refresh token yake imeisha muda, au Claude Code ilifuta taarifa zake za kuingia baada ya seva kuzikataa. Kabla kivinjari cha kwanza hakijafunguka, hukagua kwamba `config.toml` inaweza kuhifadhiwa.

`ccam logout` huendesha `claude auth logout` ndani ya akaunti na hufuta barua pepe iliyohifadhiwa; rekodi ya kuingia hubaki, kwa hiyo kuingia kunakofuata hukaguliwa dhidi yake. Claude Code hubatilisha kuingia huko na hufuta kile inachohifadhi kwa ajili ya akaunti, ikiwemo kuingia kwenye seva za MCP. Kwa hiyo ccam hukuonya na huuliza kwanza (`--yes` huruka swali) wakati:

- kuingia huko kunashikilia kuingia kwa MCP;
- akaunti inaendeshwa (vipindi vyake navyo hutolewa nje ya akaunti);
- folda ya akaunti nyingine inashikilia kuingia kule kule (folda iliyonakiliwa), ambako nako hubatilishwa na kutoka huko.

Kwa mfano:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

Ili kuingia kama akaunti nyingine ya Claude, `ccam login <name>` peke yake inatosha, na huhifadhi kuingia kwa MCP.

`ccam list` huongeza safu ya `Organization` akaunti zinapokuwa za mashirika zaidi ya moja (mipango ya binafsi huonyesha `personal`), na `ccam list --json` ina `plan` na `org` kwa kila akaunti.

### Kuingia kwenye mashine ya mbali

Kupitia SSH, sehemu ya kivinjari ya `claude auth login` inaweza kushindwa. Claude Code husubiri kuingia kwenye seva ya callback kwenye mashine ya mbali (`http://localhost:<port>/callback`), lakini kivinjari kinachokamilisha kuingia kinaendeshwa kwenye kompyuta yako, ambako anwani hiyo haielekei popote na kivinjari huonyesha hitilafu kama `ERR_CONNECTION_REFUSED`. Anwani hiyo bado ina msimbo wa mara moja, kwa hiyo ccam inaweza kuuwasilisha:

<!-- motion:accounts-relay -->

ccam haikabidhi terminal yako kwa Claude Code (`runLogin`). Huunganisha standard input ya Claude Code na pipe, na hunakili ndani yake kila mstari unaouandika (`relayLogin`), kwa hiyo msimbo unaobandikwa kwenye swali la `Paste code here` la Claude Code bado hufika. Mstari unapokuwa anwani ya callback, ccam huikubali tu ikiwa host ni `localhost`, `127.0.0.1` au `[::1]`, port imetolewa, njia ni `/callback` na `code` pamoja na `state` zote zipo; `http://` inaweza kuachwa (`callbackURL`). Kisha ccam huomba anwani hiyo kwenye seva kwa muda wa mwisho wa sekunde 10 na bila kufuata redirects, kama kivinjari kingefanya (`deliverCallback`). Ombi hufika kwenye seva ya callback ya `claude auth login` inayosubiri (mstari wa vistari kwenye mchoro). Claude Code hupokea msimbo na huhifadhi taarifa za kuingia, kwenye Linux katika `.credentials.json` ndani ya folda ya akaunti. Mstari uliobandikwa hauwezi kamwe kuifanya ccam iwasiliane na host nyingine.

ccam huchapisha hiki kabla kuingia hakujaanza:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

Ikifaulu huchapisha `→ Passed the sign-in to Claude Code on this machine.`

Wakati kuingia kunasubiri ingizo kwenye terminal, ccam huzima mwangwi (echo) na hali ya mstari ya terminal na huonyesha `*` kwa kila herufi unayoandika, kwa hiyo msimbo au anwani ya callback iliyobandikwa haionekani kwenye skrini. ccam hufanya uhariri wa mstari yenyewe (`readLine`): `Enter` hutuma mstari, `Backspace` hufuta herufi moja, na `Ctrl+D` kwenye mstari mtupu humaliza ingizo. Terminal hurejeshwa kuingia kunapoisha, na pia kwa `Ctrl+C`, `SIGTERM` au `SIGHUP`, ambazo bado humaliza ccam. Bila terminal kwenye standard input na standard error, mistari hunakiliwa kama inavyokuja.

### Kubadilisha jina, kuondoa na kuingiza

`ccam rename` na `ccam remove` hufanya kazi kwenye folda ya akaunti na kwenye kile ambacho ccam huhifadhi mahali pengine kwa jina la akaunti. ccam huorodhesha njia hizo katika `nameState`:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- cache za matumizi `~/.local/state/ccam/usage/<name>.json` na `<name>.statusline.json`;
- rekodi ya kikomo `<name>.limit.json`;
- maelezo ya vipindi katika `~/.local/state/ccam/handoff/<name>/`.

`ccam remove work` huzishughulikia kwa mpangilio huu:

<!-- motion:accounts-remove -->

1. `ccam remove` kwanza hukataa ikiwa akaunti inaendeshwa, kisha huuliza, na hukagua tena baada ya kujibu (hili hugundua kipindi kilichoanzishwa wakati swali lilikuwa wazi).
2. Hukagua kwamba `config.toml` inaweza kuhifadhiwa kabla ya kubadilisha chochote.
3. Ikiwa akaunti inashiriki viendelezi, ccam huandika upya `installed_plugins.json` na `known_marketplaces.json` za kikundi (`RewritePluginPaths`), kwa sababu Claude Code huhifadhi njia kamili za plugins ambazo zinaweza kupita kwenye folda ya akaunti hii.
4. Kisha ccam hutoka kwenye akaunti inayosimamiwa (ile ambayo folda yake iliundwa na ccam, tofauti na ile iliyosajiliwa kwa `ccam import`): `claude auth logout` huendeshwa katika folda yake, kwa hiyo hakuna kuingia kunakofanya kazi kunakosubiri kwenye trash (au, kwenye macOS, katika kipengee cha Keychain kilichopewa jina la folda, ambacho akaunti ya baadaye yenye jina hilo ingekirithi). Hatua hii huachwa kwa akaunti ya API key na wakati folda ya akaunti nyingine inashikilia kuingia kule kule, ambako nako kungebatilishwa na kutoka huko. Kutoka kunaposhindwa huripotiwa, na uondoaji huendelea.
5. ccam huhamisha folda kwenda `~/.local/share/ccam/trash/<name>-<time>/` (`profile.Trash`), akaunti huondoka kwenye `config.toml`, njia za `nameState` hufutwa ili akaunti ya baadaye yenye jina lile lile isizirithi, na kiungo cha njia fupi huondolewa: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

Ikiwa hatimaye `config.toml` haiwezi kuhifadhiwa, folda hurudishwa kutoka kwenye trash. Ili kurejesha akaunti iliyoondolewa, hamisha folda kutoka kwenye trash na uingie tena.

`ccam rename` hubadilisha jina la folda na huandika upya njia za plugins katika `plugins/` yake kutoka folda ya zamani kwenda mpya. Huhamisha njia za `nameState` kwenda jina jipya na huiacha akaunti kuwa chaguo-msingi ikiwa ilikuwa hivyo. Pia husasisha maelezo ya vipindi, ili `ccam next --continue` bado iyapate. Amri ya njia fupi huhifadhi jina lake. Uandikaji upya wa plugins ukishindwa, au `config.toml` isipoweza kuhifadhiwa, folda hurudishiwa jina lake la zamani.

Kwenye macOS, na kwenye Windows Claude Code inapohifadhi kuingia katika Credential Manager, folda yenyewe hubaki na jina lake la zamani, kwa sababu kuingia na kuingia kwa MCP huhifadhiwa chini ya njia ya folda, na folda iliyohamishwa ingepoteza kuingia kwake. Kisha usanidi hurekodi folda kama `folder` kando ya jina jipya, na ccam husema `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

`ccam import` husajili folda pale ilipo, kama `dir` katika `config.toml`. Kamwe haiihamishi wala kuinakili: refresh token iliyonakiliwa ingesababisha nakala hizo zitoane nje ya akaunti.

Bila hoja ya folda, hutafuta folda za `~/.claude` na `~/.claude-*` ambazo hazijasajiliwa na huuliza kuhusu kila moja. Majina hutokana na folda (`~/.claude` ni `main`, `~/.claude-work` ni `work`, `~/.claude-2` ni `acct2`). Ikiwa ccam haiwezi kupata jina, au jina limechukuliwa, huorodhesha folda pamoja na amri ya `--name` ya kuendesha. Folda iliyopatikana inayoshindwa huripotiwa kama `Skipped <path>: <reason>` na nyingine huingizwa.

ccam huuliza maswali yake yote kwanza, kisha hupakia upya `config.toml` na husajili folda, kwa hiyo mabadiliko ambayo amri nyingine za ccam zilihifadhi katikati hubaki. Mistari ya `✓` huja baada ya kuhifadhi.

Bila `CLAUDE_CONFIG_DIR`, Claude Code huhifadhi `.claude.json` katika folda yako ya nyumbani kando ya `~/.claude`, kwa hiyo unapoingiza `~/.claude` na haina `.claude.json` yake yenyewe, ccam huunda kiungo `~/.claude/.claude.json` kwenda `~/.claude.json`. Kisha `claude` ya kawaida na `ccam main` huona akaunti ile ile.

Ikiwa `~/.claude` ina `.claude.json` yake yenyewe kando ya `~/.claude.json`, ccam huziacha zote mbili na husema hivyo. `ccam main` ingetumia ile ya ndani na `claude` ya kawaida ile ya nje, hivyo taarifa zao za kuingia na seva zao za MCP zingebaki tofauti. Ili ubaki na moja, badilisha faili ya ndani kwa kiungo ([doctor](/features/doctor) pia huiripoti).

Folda inaposhikilia kuingia kule kule kama akaunti iliyosajiliwa (folda iliyonakiliwa), ccam husema hivyo: ile itakayofanya upya kuingia huko kwanza huitoa nyingine nje ya akaunti, kwa hiyo endesha `ccam login <name>` kabla ya kuitumia.

Akaunti zilizoingizwa haziko katika kikundi chochote hadi uziongeze kwenye kimoja (`ccam group join <group> <name>`), na hakuna njia fupi inayoundwa.

`ccam completion bash|zsh` huchapisha script ya kukamilisha amri ambayo pia hukamilisha majina ya akaunti.

## Usalama na hali maalum

- `ccam rename` na `ccam remove` hukataa wakati akaunti inaendeshwa (`refuseRunning`), kwa ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``, kwa sababu kipindi kinaendelea kuandika ndani ya folda. Ukaguzi huu hushughulikia Claude Code bila kujali ilivyoanzishwa: kupitia ccam, `claude` ya kawaida, IDE, vipindi vya chinichini na daemon yao.
- `ccam rename`, `ccam remove` na `ccam logout` hukataa wakati mchakato mwingine wa ccam unafanya upya au unabadilisha kuingia huko: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- Wakati badiliko la vikundi linaendelea katika mchakato mwingine wa ccam, `ccam rename` na `ccam remove` hukataa kwa `another ccam process is changing the groups right now; try again when it is done`; baada ya badiliko kusimama nusu njia, kwa ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.
- `ccam remove` kwenye akaunti iliyoingizwa huhifadhi folda na kuingia kwake, kwa ajili ya `claude` ya kawaida. Ikiwa akaunti iko katika vikundi, kwanza huviacha kama `ccam group leave` inavyofanya: kila kipengele kinachoshirikiwa huwa nakala ya kile cha kikundi ndani ya folda, na seva za MCP za kikundi huingia katika `.claude.json` yake. Kisha huondolewa kwenye usajili: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. Nakala rudufu za `.claude.json` yake ambazo ccam iliunda huhamia `<folder>/.ccam-backup/mcp-<time>`; zisipoweza kuhamishwa, hubaki katika `~/.local/state/ccam/mcp/<name>/backups` na ujumbe husema hivyo. `~/.claude` yako haihamishwi kwenda kwenye trash.
- Akaunti mbili haziwezi kuwa na folda ile ile: `config.toml` inayotaja folda moja mara mbili hushindwa kupakiwa.
- `ccam add` hukataa jina ambalo folda yake tayari ipo katika `profiles/`, na huelekeza kwa `ccam import`. Kwenye macOS pia husimama Keychain inapokuwa bado inashikilia kuingia au kuingia kwa MCP kwa njia ya folda mpya, kulikoachwa na toleo la zamani la ccam, na hutaja amri ya `security delete-generic-password` inayoviondoa. Ikiwa `config.toml` haiwezi kuhifadhiwa, folda mpya huondolewa tena.
- `ccam import` hukataa folda ambayo tayari imesajiliwa, njia ambayo si folda, na jina lililokwisha kuchukuliwa. ccam pia huonya folda inapokuwa nakala ya `~/.claude` (kuingia kule kule), hata ikiwa `~/.claude` yenyewe haijasajiliwa.
- `ccam logout` huuliza kwanza akaunti inapoendeshwa; `ccam login` haikagui.
- Ikiwa relay haiwezi kufikia Claude Code, kwa mfano baada ya kuingia kumalizika, huchapisha `ccam: could not pass the sign-in to Claude Code:` pamoja na sababu. Hitilafu ya HTTP kwa kawaida inamaanisha msimbo umeisha muda; anza kuingia tena.

## Zinazohusiana

- [Kuanzisha Claude Code](/features/launch)
- [Amri za njia fupi](/features/shortcuts)
- [Shiriki seva za MCP, mipangilio na historia](/docs/sharing)
- [Tumia akaunti kadhaa za Claude Code kwenye kompyuta moja](/guides/multiple-claude-code-accounts)
