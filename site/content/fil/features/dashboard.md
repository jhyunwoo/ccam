---
title: Ang dashboard
description: Sa terminal, nagbubukas ang ccam ng full-screen dashboard ng mga account, group, folder rule, extension at pagsusuri; may plano muna ang bawat pagbabago.
updated: 2026-10-30
order: 15
translation: machine
sourceHash: "d65f1b5bb0c8"
---

Binubuksan ng `ccam` na walang argument ang isang full-screen dashboard kapag terminal ang standard input at output; kung hindi, ipi-print nito ang `ccam list`. Ipinapakita ng dashboard ang bawat account kasama ang usage at mga group nito, at ginagawa nito ang ginagawa ng mga command: magsimula, magdagdag, mag-log in, magbago ng mga group, magtakda ng mga folder rule, maglista ng mga background session, magsuri at mag-ayos. Nagsasara ito kapag nagsimula ang isang account o kapag kumabit ka sa isang background session; kapag nagdagdag ng account o nag-log in, ibinibigay muna nito sandali ang terminal, saka bumabalik sa dashboard.

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

Sinasabi ng linya sa ilalim ng mga pangalan ng tab, sa isang pangungusap, kung para saan ang kasalukuyang tab. Sa mas malapad na terminal, ipinapakita rin ng tab na Accounts ang email ng bawat account. Sa bawat tab, ipinapakita ng dulo ng linya ng mga tab kung naka-on ang [`auto_switch`](/features/next#auto-switch) (Linux at macOS; hindi ito pinapansin ng Windows). Kapag naka-on ito, inililista ng tab na Health kung ano pa ang kailangan nito: Claude Code 2.1.214 o mas bago, at ang ccam plugin na naka-on sa bawat account na nagsa-sign in gamit ang Claude subscription.

Tuwing nagbubukas ito, tinatanong ng dashboard sa GitHub sa background kung may mas bagong ccam release; hindi kailanman naghihintay rito ang unang frame, at walang ipinapakita kapag offline ka. Kapag mayroon, inirerekomenda ito ng isang linya sa ilalim ng pamagat:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

Nagtatanong muna ang `U`, saka nag-a-update gaya ng ginagawa ng [`ccam update`](/features/update): sinusuri ang release laban sa mga SHA-256 checksum nito at pinapalitan nito ang ccam sa isang hakbang. Lumang bersyon pa rin ang dashboard na kinaroroonan mo, kaya lumabas at simulan muli ang `ccam`. Kapag hindi magagawa rito ang update, halimbawa para sa kopya ng ccam na in-install mula sa `.deb` package, ipinapakita ng dashboard ang parehong payo gaya ng `ccam update` at pinananatili ang alok.

Iginuguhit ang unang frame mula sa `config.toml` lang, sa mas mababa pa nang husto sa 100 ms; napupunan ang usage habang sumasagot ang bawat account. May sariling bar ang bawat limit: ang 5-oras na window, ang lingguhang window para sa lahat ng model, at ang sariling lingguhang window ng Fable; ang ibig sabihin ng `-` ay walang ganoong window ang account. Nagiging dilaw ang porsiyento sa 80% at pula sa 100%. Kapag pumalya ang isang lookup, malabong ipinapakita ang mga naka-cache na value. 100% ang ipinapakita ng window na iniulat ng Claude Code na puno na hanggang sa reset nito. Nagpapakita ang expired na access token ng malabong `⚠ token expired` (nire-renew ito ng Claude Code sa susunod na pagsisimula), habang ang ibig sabihin ng `⚠ login expired` ay kailangan ng account ang `ccam login`.

## Mga tab

| Tab | Ano ang ipinapakita | Mga key |
| --- | --- | --- |
| Accounts | numero, pangalan, email, mga bar ng 5-oras, lingguhan at Fable, huling paggamit, mga group; `✔` ang default, `▸` ang account ng folder na ito, `•` tumatakbo (pati ang Claude Code na sinimulan nang walang ccam: karaniwang `claude`, isang IDE, mga background session at ang daemon nila) | sa ibaba |
| Agents | ang inililista ng [`ccam agents`](/features/agents): ang mga background session at bukas na session ng bawat account, kasama ang state, folder at oras ng pagsisimula nila (sa ibaba) | `enter` kumabit (gaya ng `ccam attach`), `x` ihinto, `r` i-reload |
| Groups | bawat group, kung ano ang shine-share nito at ang mga account nito, at ang mga account na walang group | `enter` baguhin, `n` bago, `e` palitan ang pangalan, `x` alisin |
| Folders | ang mga patakaran ng [`ccam folder`](/features/folder) at ang nalalapat dito, kahit narating mo ang folder sa pamamagitan ng symbolic link (sine-save ang mga bagong patakaran na naresolba ang mga link) | `n` gumamit ng account dito, `enter` baguhin, `x` alisin |
| Extensions | ang ipinapakita ng [`ccam ext`](/features/ext), bawat extension nang isang beses kasama ang mga group na pinanggagalingan nito, `off:` ang mga account kung saan hindi naka-on ang isang plugin, ang sariling mga server ng account na para lang sa folder na ito, `This folder's project`, at ang mga problemang hindi pag-aari ng anumang account | `r` i-reload |
| Health | ang iniuulat ng [`ccam doctor`](/features/doctor), bilang isang kahon bawat account (sa ibaba) | `f` ayusin ang napiling kahon, `F` ayusin lahat (`doctor --fix`), `↑`/`↓` lumipat sa pagitan ng mga kahon |

Nagpapalit ng tab ang `←`/`→` o `Tab`, inililista ng `?` ang mga key ng kasalukuyang tab, nire-refresh ng `r` ang usage nang nilalampasan ang maikling cache at nilo-load ulit ang mga tab na Agents, Extensions at Health, at lumalabas ang `q` o `Esc`.

## Mga key sa tab na Accounts

| Key | Ano ang ginagawa |
| --- | --- |
| `Enter` | Simulan ang Claude Code bilang ang napiling account |
| `1`–`9` | Pumunta sa account na may numerong iyon; i-type nang mabilis ang `1` `2` para sa 12. `Enter` lang ang nagsisimula ng account |
| `↑`/`↓`, `k`/`j` | Lumipat |
| `n` | Simulan ang account na pipiliin ng [`ccam next --auto`](/features/next) |
| `o` | Simulan nang may mga option: `--continue`, `--resume`, ibang argument, o dry run |
| `u` | Mga detalye ng usage: ang 5-oras, lingguhan at model na mga window, kasama ang mga reset nila at kung saan galing ang mga numero |
| `a` | Magdagdag ng account (ang wizard sa ibaba) |
| `l` | Mag-log in muli |
| `s` | Piliin ang mga group ng account |
| `d` | Gawin itong default na account |
| `f` | Gamitin ito sa kasalukuyang folder |
| `e` | Palitan ang pangalan nito |
| `c` | Itakda o alisin ang shortcut command nito |
| `x` | Alisin ito (nagtatanong muna, na nakapili ang No) |

**Mga seat ng Team at Enterprise.** Lumilipat din ang awtomatikong pagpili ng account (`ccam next`, `--fallback`, `ccam pick`, ang `n` ng dashboard at, kapag binuksan mo ito, ang [`auto_switch`](/features/next#auto-switch)) sa pagitan ng mga seat ng iisang Claude Team o Enterprise na organisasyon, na maaaring lumabag sa terms ng Anthropic. Basahin ang [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) at magtanong muna sa iyong organisasyon at sa Anthropic bago ka umasa rito.

## Ang tab na Health

Pinapatakbo ng tab na Health ang mga pagsusuri ng [`ccam doctor`](/features/doctor) kapag binuksan mo ito at iginuguhit ang parehong view: isang unang linya na may marka bawat account, saka isang kahon bawat account. Ang bawat hilera sa isang kahon ay isang bahaging sinuri, `Folder`, `Login` at `Share links` muna, na may markang `✔` OK, `ℹ` impormasyon, `⚠` babala o `✘` error, kaya makikita mo agad kung aling account ang may problema at saan. Napupunta sa kahong `General` ang mga natuklasang hindi kabilang sa anumang account.

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

Lumilipat ang `↑`/`↓` mula kahon papunta sa kahon, at may border na kulay accent ang napiling kahon. Kinukumpuni lang ng `f` ang napiling kahon, sa paraang gagawin ng `ccam doctor --fix`: sa kahon ng isang account, ang account na iyon lang ang kinukumpuni nito, at sa kahong `General`, ang mga natuklasang hindi kabilang sa anumang account (mga path ng plugin at mga shortcut); sinusuri lang ulit ang bawat iba pang account. Nagtatanong muna ito (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), at ang kahong walang kailangang kumpunihin ay nagsasabi lang ng `nothing to repair in <name>`. Kinukumpuni ng `F` ang lahat nang sabay-sabay, gaya ng `ccam doctor --fix`, at sinusuri ulit ng `r`.

## Ang tab na Agents

Ipinapakita ng tab na Agents ang ipi-print ng [`ccam agents`](/features/agents): ang mga background session at bukas na session ng bawat account, kasama ang account na kinabibilangan ng bawat isa. Sabay-sabay nitong tinatanong ang Claude Code ng bawat account (`claude agents --json`, bawat isa na may 15 segundong timeout) kapag una mong binuksan ang tab, sa `r`, pagkatapos ng pagpapahinto, at pagkatapos ng pagbabago sa config, at sinasabi nito ang `Asking every account's Claude Code for its sessions…` hanggang makasagot ang lahat.

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

Ang `State` ay ang state ng isang background session at kung ano ang ginagawa ng process nito, o `open` para sa interactive session, gaya sa `ccam agents`. Kung hindi matanong ang isang account, ipinapakita ng linyang `✘` sa ilalim ng listahan ang sinabi ng Claude Code, gaya ng `lab` sa itaas, na naka-off ang agent view. Hindi isinasama ang mga tapos nang background session (inililista ang mga ito ng `ccam agents --all`), at kapag sumagot ang bawat account at walang may session, sinasabi ng tab ang ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` Sa makitid na terminal, tinatanggal ng tab ang `Started`, saka ang `Name`, saka ang `Folder`, at pinuputol mula sa kaliwa ang mahahabang folder path.

Isinasara ng `Enter` sa isang background session ang dashboard at binubuksan ang session sa terminal na ito bilang ang account nito, gaya ng ginagawa ng `ccam attach <id>`; hindi na bumabalik ang dashboard (sa Claude Code, pumupunta ang `←` sa agent view nito at ang `Ctrl+Z` sa shell). Walang ID ang bukas na session, kaya sinasabi lang ng `Enter` at `x` na kailangan ng background session. Nagtatanong muna ang `x`, na nakapili ang `No, keep it`: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` Kapag oo, pinapatakbo ng ccam ang `claude stop <id>` bilang account na iyon, sinasabi ang ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` o kung bakit hindi nito nagawa, at binabasa ulit ang listahan.

## Pagbabago ng mga group

Sa tab na Groups, binubuksan ng `Enter` ang editor: ang mga item sa kaliwa, ang mga account sa kanan, `Tab` para magpalit, `Space` para mag-toggle. Kulay-abo, kasama ang dahilan, ang item na nakukuha na ng isa sa mga account mula sa ibang group, gayundin ang account na makakakuha ng isang item nang dalawang beses; iisang group lang ang puwedeng pinanggagalingan ng isang item. Ipinapakita ng `Enter` ang plano:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Kung naiiba ang isang file sa kopya ng group, itinatanong muna ng ccam kung ano ang gagawin (panatilihin ang sa group, gamitin ang akin o pagsamahin). Pagkatapos ay tumatakbo ang pagbabago sa background gamit ang parehong code ng [`ccam group`](/features/share): sine-save muna ang config, inililipat nang sunod-sunod ang data ng bawat account, at napupunta ang mga orihinal sa `.ccam-backup/`; pinapanatili ng account na umaalis sa pag-share ng MCP ang mga server ng group sa `.claude.json` nito.

Tinatanggihan, kasama ang pangalan ng account, ang pagbabagong maglilipat ng data ng tumatakbong account, kahit nagsimula ang account habang bukas ang plano, at tinatanggihan na may `the groups changed while this plan waited; nothing was moved, make the change again` ang planong naluma dahil binago ng ibang ccam process ang mga group. Habang nagpapalit ng mga group ang ibang ccam process, sinasabi ng bagong pagbabago, pagpapalit ng pangalan o pag-alis ang `another ccam process is changing the groups right now`; pagkatapos ng pagbabagong huminto sa kalagitnaan, naghihintay ang mga ito hanggang matapos ito ng `F` sa tab na Health (o ng `ccam doctor --fix`).

## Pagdagdag ng account

Binubuksan ng `a` ang wizard na ginagamit din ng `ccam add` sa terminal:

1. Ang pangalan ng account.
2. Kung paano ito magsa-sign in: Claude subscription ngayon (sa browser), [API key](/features/api), o Claude subscription mamaya.
3. Para sa API key: kung saan binabasa ng ccam ang key sa bawat pagsisimula (`env:NAME` o `command:COMMAND`), ang endpoint, kung paano ipinapadala ang key at ang model.
4. Ang mga group na sasalihan, na naka-check ang `default` (kapag wala pang group, kung ano ang ishe-share ng group na `default`).
5. Isang shortcut command, kung gusto mo.
6. Isang buod na kukumpirmahin.

Bumabalik nang isang hakbang ang `Esc`, at sa unang hakbang, isinasara nito ang wizard.

## Pag-alis ng account

Nagtatanong muna ang `x` bago may mangyari, na nakapili ang No, at sinasabi kung ano ang ginagawa ng pag-alis: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` Walang login na mababawi para sa API key account, at kapag hawak ng ibang account ang parehong login, hindi ito ina-log out. Pinapanatili ng in-import na account ang folder nito: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## Mga kulay

Sinusunod ng mga kulay ang mga papel ng kulay sa terminal UI ng Claude Code, kaya bagay na bagay ang ccam sa tabi nito; nakakakuha ng kulay ang mga pangalan ng group batay sa pangalan nila. Tinatanong ng ccam ang terminal tungkol sa background nito at pumipili ng maliwanag o madilim na theme. Puwede mo itong itakda mismo sa itaas ng `config.toml`, bago ang unang `[[group]]`, `[[profile]]` o `[[folder]]`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Para sa red-green color blindness ang mga theme na `-daltonized`, at ginagamit ng mga theme na `-ansi` ang sariling 16 na kulay ng terminal. Pinapatay ng `NO_COLOR` ang mga kulay.

## Kaligtasan at mga espesyal na kaso

- Tinatawag ng dashboard ang parehong mga function ng mga command, kaya nalalapat ang bawat patakaran: mga plano bago maglipat ng data, pagtanggi para sa mga tumatakbong account, mga backup, at mga login token na hindi kailanman umaalis sa isang account.
- Walang dry run ang dashboard: ipi-print ng `ccam --dry-run` na walang pangalan kung paano magsisimula ang account ng folder na ito, at nag-aalok ang `o` sa isang account ng dry run ng iisang pagsisimulang iyon.
- Naibabalik ang terminal bago magsimula ang Claude Code o bumukas ang login.
- Sa screen na masyadong maliit para sa isang dialog, unang tinatanggal ang mga hilera ng tab; pinuputol ang mga linya gamit ang `…` sa halip na i-wrap. Sa maikling screen, ang mga blangkong linya at ang paglalarawan ng tab ang susunod na tinatanggal, kaya nananatiling nakikita ang linya ng mga pahiwatig ng key. Kapag mas mababa sa 40 column o 8 hilera, `Too small: ccam needs 40x8 (now 36x7)` lang ang ipinapakita ng dashboard.
- Kapag pinindot ang `Ctrl+C` habang inilalapat ang isang pagbabago, pag-alis, update o pagkukumpuni, hinihintay muna nitong matapos iyon: sinasabi ng status line ang `ccam quits when it is done.`
- Pagkatapos ma-update ng `U` ang ccam, hindi na bumabalik ang alok kapag bumalik ang dashboard mula sa pagdagdag ng account o sa isang login.

## Kaugnay

- [Pag-share gamit ang mga group](/features/share)
- [Pagsisimula ng Claude Code](/features/launch)
- [Pamamahala ng mga account](/features/accounts)
- [Mga usage limit](/features/usage)
