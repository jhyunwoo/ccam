---
title: Dashboard
description: ccam a terminal yana buɗe dashboard mai cike allo na asusai, rukunoni, ƙa'idodin folda, extensions da dubawa; kowane canji yana fara nuna shirinsa.
updated: 2026-10-30
order: 15
translation: machine
sourceHash: d65f1b5bb0c8
---

`ccam` ba tare da arguments ba yana buɗe dashboard mai cike allo idan standard input da output terminal ne; in ba haka ba yana buga `ccam list`. Dashboard yana nuna kowane asusu tare da amfaninsa da rukunoninsa, kuma yana yin abin da umarnan ke yi: farawa, ƙarawa, shiga, canza rukunoni, saita ƙa'idodin folda, lissafa zaman bayan fage, dubawa da gyarawa. Farawa, ko haɗawa da zaman bayan fage, yana rufe shi; ƙara asusu ko shiga yana miƙa terminal na ɗan lokaci, sannan ya dawo kan dashboard.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · r to refresh · a to add · ? for shortcuts
```

Layin da ke ƙarƙashin sunayen shafuka yana faɗin a jimla ɗaya abin da shafin yanzu yake yi. A terminal mai faɗi, shafin Accounts yana kuma nuna imel ɗin kowane asusu. A kowane shafi, ƙarshen layin shafuka yana nuna ko [`auto_switch`](/features/next#auto-switch) yana a kunne (Linux da macOS; a Windows ana watsi da shi). Idan yana a kunne, shafin Health yana lissafa abin da har yanzu yake buƙata: Claude Code 2.1.214 ko sabo, da plugin ɗin ccam a kunne a kowane asusu da ke shiga da biyan kuɗin Claude.

Duk lokacin da ya buɗe, dashboard yana tambayar GitHub a bayan fage ko akwai sabon sakin ccam; firam na farko ba ya taɓa jiran sa, kuma idan ba ka da intanet, ba a nuna komai. Idan akwai, layi a ƙarƙashin take yana ba da shawarar sa:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` yana tambaya da farko, sannan yana sabuntawa yadda [`ccam update`](/features/update) ke yi: ana duba sakin da SHA-256 checksums nasa kuma yana maye gurbin ccam a mataki ɗaya. Dashboard ɗin da kake ciki har yanzu tsohon sigar ne, don haka ka fita ka sake fara `ccam`. Idan ba za a iya yin sabuntawar a nan ba, misali ga kwafin ccam da aka girka daga kunshin `.deb`, dashboard yana nuna shawara ɗaya da ta `ccam update` kuma yana riƙe tayin.

Firam na farko yana zana daga `config.toml` kaɗai, cikin ƙasa da 100 ms sosai; amfani yana cika yayin da kowane asusu ya amsa. Kowace iyaka tana da nata sandar: taga ta awa 5, taga ta mako ga dukkan model, da taga ta mako ta Fable kanta; `-` yana nufin asusun ba shi da irin wannan taga. Kashi yana zama rawaya a 80% kuma ja a 100%. Idan nema ya gaza, ana nuna ƙimomin cache a dusashe. Tagar da Claude Code ya ba da rahoton ta cika tana nuna 100% har zuwa sake saitinta. Access token da ya ƙare yana nuna `⚠ token expired` a dusashe (Claude Code yana sabunta shi a farawa ta gaba), yayin da `⚠ login expired` yana nufin asusun yana buƙatar `ccam login`.

## Shafuka

| Shafi | Abin da yake nunawa | Maɓallai |
| --- | --- | --- |
| Accounts | lamba, suna, imel, sandunan awa 5, na mako da na Fable, amfani na ƙarshe, rukunoni; `✔` tsoho, `▸` asusun wannan foldar, `•` mai aiki (har da Claude Code da aka fara ba ta ccam ba: `claude` na yau da kullum, IDE, zaman bayan fage da daemon ɗinsu) | a ƙasa |
| Agents | abin da [`ccam agents`](/features/agents) ke lissafawa: zaman bayan fage da zaman da ke buɗe na kowane asusu, tare da yanayinsu, foldarsu da lokacin farawarsu (a ƙasa) | `enter` haɗawa (kamar `ccam attach`), `x` dakatarwa, `r` sake lodawa |
| Groups | kowane rukuni, abin da yake rabawa da asusansa, da asusan da ba sa cikin kowane rukuni | `enter` gyara, `n` sabo, `e` sake suna, `x` cire |
| Folders | ƙa'idodin [`ccam folder`](/features/folder) da wanda ke aiki a nan, har ma idan ka isa foldar ta mahaɗi na alama (symbolic link) (ana ajiye sababbin ƙa'idodi tare da warware mahaɗai) | `n` yi amfani da asusu a nan, `enter` canza, `x` cire |
| Extensions | abin da [`ccam ext`](/features/ext) ke nunawa, kowane extension sau ɗaya tare da rukunonin da yake fitowa daga gare su, `off:` asusan da ba a kunna plugin a cikinsu ba, sabobin asusun na kansa da ke aiki ga wannan foldar kawai, `This folder's project`, da matsalolin da ba na wani asusu ba | `r` sake lodawa |
| Health | abin da [`ccam doctor`](/features/doctor) ke bayar da rahoto, a matsayin akwati ɗaya ga kowane asusu (a ƙasa) | `f` gyara akwatin da aka zaɓa, `F` gyara duka (`doctor --fix`), `↑`/`↓` matsawa tsakanin akwatuna |

`←`/`→` ko `Tab` suna sauya shafuka, `?` yana lissafa maɓallan shafin yanzu, `r` yana sabunta amfani, yana tsallake gajeren cache, kuma yana sake lodawa shafukan Agents, Extensions da Health, kuma `q` ko `Esc` yana fita.

## Maɓallai a shafin Accounts

| Maɓalli | Abin da yake yi |
| --- | --- |
| `Enter` | Fara Claude Code a matsayin asusun da aka zaɓa |
| `1`–`9` | Je kan asusun da ke da wannan lambar; rubuta `1` `2` da sauri don 12. `Enter` ne kaɗai ke farawa |
| `↑`/`↓`, `k`/`j` | Matsawa |
| `n` | Fara asusun da [`ccam next --auto`](/features/next) zai zaɓa |
| `o` | Fara da zaɓuɓɓuka: `--continue`, `--resume`, wasu arguments, ko dry run |
| `u` | Cikakken bayanin amfani: tagogin awa 5, na mako da na model tare da lokutan sake saitinsu da inda lambobin suka fito |
| `a` | Ƙara asusu (mataimakin da ke ƙasa) |
| `l` | Sake shiga |
| `s` | Zaɓi rukunonin asusun |
| `d` | Mai da shi asusun tsoho |
| `f` | Yi amfani da shi a foldar yanzu |
| `e` | Sake masa suna |
| `c` | Saita ko cire gajeren umarninsa |
| `x` | Cire shi (yana fara tambaya, tare da No a zaɓe) |

**Kujeru (seats) na Team da Enterprise.** Zaɓin asusu na atomatik (`ccam next`, `--fallback`, `ccam pick`, maɓallin `n` na dashboard kuma, da zarar ka kunna shi, [`auto_switch`](/features/next#auto-switch)) yana kuma sauyawa tsakanin kujerun ƙungiyar Claude Team ko Enterprise iri ɗaya, wanda na iya saɓa wa sharuɗɗan Anthropic. Ka karanta [Tambayoyin da ake yawan yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) kuma ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi.

## Shafin Health

Shafin Health yana gudanar da dubawar [`ccam doctor`](/features/doctor) idan ka buɗe shi kuma yana zana kallo iri ɗaya: layin farko mai alama ga kowane asusu, sannan akwati ɗaya ga kowane asusu. Kowane layi a akwati fanni ɗaya ne da aka duba, `Folder`, `Login` da `Share links` da farko, tare da alamar `✔` lafiya, `ℹ` bayani, `⚠` gargaɗi ko `✘` kuskure, don haka kana gani nan take wane asusu ke da matsala da kuma inda. Sakamakon da ba na wani asusu ba yana shiga akwatin `General`.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` suna matsawa daga akwati zuwa akwati, kuma akwatin da aka zaɓa yana da iyaka mai launin accent. `f` yana gyara akwatin da aka zaɓa kawai, kamar yadda `ccam doctor --fix` ke yi: a akwatin asusu yana gyara asusun nan kaɗai, a akwatin `General` kuma sakamakon da ba na wani asusu ba (hanyoyin plugins da gajerun umarni); kowane sauran asusu ana sake duba shi ne kawai. Yana fara tambaya (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), kuma akwatin da babu abin gyarawa a cikinsa yana cewa `nothing to repair in <name>` kawai. `F` yana gyara komai a lokaci ɗaya, kamar `ccam doctor --fix`, kuma `r` yana sake dubawa.

## Shafin Agents

Shafin Agents yana nuna abin da [`ccam agents`](/features/agents) ke bugawa: zaman bayan fage da zaman da ke buɗe na kowane asusu, tare da asusun da kowanne ke ciki. Yana tambayar Claude Code na kowane asusu a lokaci ɗaya (`claude agents --json`, kowanne da iyakar lokaci ta daƙiƙa 15) idan ka fara buɗe shafin, a `r`, bayan dakatarwa, da bayan canji a config, kuma yana cewa `Asking every account's Claude Code for its sessions…` har sai duka sun amsa.

```console
  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Background agents and open sessions of every account; enter attaches to one.

    Account   ID        Name   State          Folder      Started
 ❯  work      -         -      open (busy)    ~/code/web  now
    personal  49ec0028  probe  blocked, idle  ~/code/app  12m

 ✘ lab: could not list its sessions: exit status 1 'claude agents --json' is
   not available in this environment.

 enter to attach · x to stop · r to reload · ? for shortcuts
```

`State` shi ne yanayin zaman bayan fage da abin da process ɗinsa ke yi, ko `open` ga zaman mu'amala, kamar a `ccam agents`. Idan ba za a iya tambayar wani asusu ba, layin `✘` a ƙarƙashin jerin yana nuna abin da Claude Code ya faɗa, kamar `lab` a sama, wanda aka kashe agent view ɗinsa. Ba a haɗa zaman bayan fage da suka ƙare ba (`ccam agents --all` yana lissafa su), kuma idan kowane asusu ya amsa kuma babu wanda ke da zama, shafin yana cewa ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` A ƙunƙuntaccen terminal shafin yana cire `Started`, sannan `Name`, sannan `Folder`, kuma yana yanke dogayen hanyoyin folda daga hagu.

`Enter` a kan zaman bayan fage yana rufe dashboard kuma yana buɗe zaman a wannan terminal a matsayin asusunsa, kamar yadda `ccam attach <id>` ke yi; dashboard ba ya dawowa (a cikin Claude Code, `←` yana zuwa agent view ɗinsa, `Ctrl+Z` kuma zuwa shell). Zaman da ke buɗe ba shi da ID, don haka `Enter` da `x` suna faɗin cewa ana buƙatar zaman bayan fage ne kawai. `x` yana fara tambaya, tare da `No, keep it` a zaɓe: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` Idan ka ce eh, ccam yana gudanar da `claude stop <id>` a matsayin wannan asusun, yana cewa ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` ko dalilin da ya sa bai iya ba, kuma yana sake karanta jerin.

## Canza rukunoni

A shafin Groups, `Enter` yana buɗe mai gyarawa: abubuwan a hagu, asusan a dama, `Tab` don sauyawa, `Space` don kunnawa ko kashewa. Abin da ɗaya daga cikin asusan ke samu tuni daga wani rukuni ana dusashe shi tare da dalilin, haka ma asusun da zai sami abu sau biyu; abu zai iya fitowa daga rukuni ɗaya kawai. `Enter` yana nuna shirin:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Idan fayil ya bambanta da kwafin rukunin, ccam yana fara tambayar abin da za a yi (riƙe na rukunin, yi amfani da nawa ko haɗa su). Sannan canjin yana gudana a bayan fage da lamba iri ɗaya da ta [`ccam group`](/features/share): ana fara ajiye saitin, bayanan kowane asusu suna motsawa a tsari, kuma fayilolin asali suna zuwa `.ccam-backup/`; asusun da ya bar rabon MCP yana riƙe sabobin rukunin a `.claude.json` nasa.

Ana ƙin canjin da zai matsar da bayanan asusun da ke aiki, tare da sunan asusun, har ma idan asusun ya fara yayin da shirin ke buɗe, kuma ana ƙin shirin da ya tsufa saboda wani process na ccam ya canza rukunoni da `the groups changed while this plan waited; nothing was moved, make the change again`. Yayin da wani process na ccam ke canza rukunoni, sabon canji, sake suna ko cirewa yana cewa `another ccam process is changing the groups right now`; bayan canjin da ya tsaya rabi, suna jira har sai `F` a shafin Health (ko `ccam doctor --fix`) ya kammala shi.

## Ƙara asusu

`a` yana buɗe mataimakin da `ccam add` ma ke amfani da shi a terminal:

1. Sunan asusun.
2. Yadda yake shiga: biyan kuɗin Claude yanzu (a burauza), [API key](/features/api), ko biyan kuɗin Claude daga baya.
3. Ga API key: inda ccam ke karanta key ɗin a kowane farawa (`env:NAME` ko `command:COMMAND`), endpoint, yadda ake aika key ɗin da model.
4. Rukunonin da za a shiga, tare da alama a kan `default` (idan babu rukunoni tukuna, abin da rukunin `default` zai raba).
5. Gajeren umarni, idan kana son ɗaya.
6. Taƙaitawa don tabbatarwa.

`Esc` yana komawa mataki ɗaya baya, kuma a mataki na farko yana rufe mataimakin.

## Cire asusu

`x` yana tambaya kafin wani abu ya faru, tare da No a zaɓe, kuma yana faɗin abin da cirewar ke yi: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` Ga asusun API key babu shigar da za a soke, kuma idan wani asusu yana riƙe da shiga iri ɗaya ba a fitar da shi. Asusun da aka shigo da shi yana riƙe foldarsa: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## Launuka

Launukan suna bin matsayin launuka na UI na terminal na Claude Code, don haka ccam yana dacewa a gefensa; sunayen rukunoni suna samun launi daga sunansu. ccam yana tambayar terminal bangonsa kuma yana zaɓar jigo mai haske ko mai duhu. Saita ɗaya da kanka a saman `config.toml`, kafin `[[group]]`, `[[profile]]` ko `[[folder]]` na farko:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Jigogin `-daltonized` na makantar launin ja da kore ne, kuma jigogin `-ansi` suna amfani da launuka 16 na terminal kansa. `NO_COLOR` yana kashe launuka.

## Tsaro da yanayi na musamman

- Dashboard yana kiran ayyuka (functions) iri ɗaya da umarnan, don haka kowace ƙa'ida tana aiki: shiri kafin bayanai su motsa, ƙin canza asusan da ke aiki, kwafin ajiya, da tokens na shiga da ba sa taɓa barin asusu.
- Dashboard ba shi da dry run: `ccam --dry-run` ba tare da suna ba yana buga yadda asusun wannan foldar zai fara, kuma `o` a kan asusu yana ba da dry run na wannan farawa ɗaya kawai.
- Ana maido da terminal kafin Claude Code ya fara ko shiga ta buɗe.
- Allon da ya yi ƙanƙanta ga akwatin tattaunawa (dialog) yana fara cire layukan shafin; ana yanke layuka da `…` maimakon naɗa su. A gajeren allo, ana cire layukan da babu komai da bayanin shafin na gaba, don layin alamun maɓallai ya ci gaba da bayyana. Ƙasa da ginshiƙai 40 ko layuka 8, dashboard yana nuna `Too small: ccam needs 40x8 (now 36x7)` kawai.
- `Ctrl+C` yayin da ake aiwatar da canji, cirewa, sabuntawa ko gyara yana jiran sa: layin matsayi yana cewa `ccam quits when it is done.`
- Bayan `U` ya sabunta ccam, tayin ba ya dawowa idan dashboard ya dawo daga ƙara asusu ko shiga.

## Masu alaƙa

- [Rabawa da rukunoni](/features/share)
- [Fara Claude Code](/features/launch)
- [Sarrafa asusai](/features/accounts)
- [Iyakokin amfani](/features/usage)
