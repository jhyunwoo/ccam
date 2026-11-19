---
title: Bayanin umarnai
metaTitle: Bayanin umarnan ccam
description: Kowane umarni da zaɓi na ccam, daga ccam add da ccam <name> zuwa ccam usage, ccam next da ccam doctor.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: 28e1ad1dcbad
---

Saka `--dry-run` kafin umarni don ganin abin da zai yi ba tare da ya yi shi ba, misali `ccam --dry-run work`. Yana aiki lokacin fara Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) da kuma tare da `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` da `doctor --fix`. Umarnan da ke karantawa kawai, kamar `list`, `usage`, `share` da `ext`, suna karɓar sa kuma ba sa canza komai ko ta yaya. Sauran suna ƙin shi maimakon su gudana, har da `group rename`.

<!-- commands -->

## Fara Claude Code

`ccam <name|number> [arguments]` yana fara Claude Code a matsayin wannan asusun. ccam ba ya karanta arguments da ke bayan sunan; yana miƙa su ga `claude` a tsari ɗaya, kuma standard input, exit code da signals suna aiki daidai kamar da `claude` kansa. `--` bayan sunan ba dole ba ne.

`ccam` tare da arguments amma ba tare da sunan asusu ba, kamar `ccam -p "summarise"`, yana amfani da asusun da aka saita don foldar yanzu, in ba haka ba asusun tsoho. Saita su da `ccam folder <name> [folder]` da `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` kafin sunan asusun na gudun `-p`/`--print` ne. Idan asusun ya kai iyakar amfaninsa, ccam yana sake gudanar da umarnin iri ɗaya a matsayin asusun da aka shiga da shi mafi sarari, kuma yana buga sakamakon gudun da ya kammala kaɗai. Duba [Asusai bisa folda](/features/folder) da [Sake gudanarwa a iyaka](/features/next#run-again-at-the-limit).

**Kujeru (seats) na Team da Enterprise.** Zaɓin asusu na atomatik (`ccam next`, `--fallback`, `ccam pick`, maɓallin `n` na dashboard kuma, da zarar ka kunna shi, [`auto_switch`](/features/next#auto-switch)) yana kuma sauyawa tsakanin kujerun ƙungiyar Claude Team ko Enterprise iri ɗaya, wanda na iya saɓa wa sharuɗɗan Anthropic. Ka karanta [Tambayoyin da ake yawan yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) kuma ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi.

Gajeren umarni kamar `ccwork` yana fara asusunsa ta hanya ɗaya: `ccwork --resume` yana yin abin da `ccam work --resume` ke yi. Duba [Gajerun umarni](/docs/commands#shortcuts) a ƙasa.

## Ƙarawa da cire asusai

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` ba tare da folda ba yana neman `~/.claude` da `~/.claude-*` kuma yana tambaya game da kowanne; ana ba da rahoton foldar da ta gaza a matsayin `Skipped <path>: <reason>` kuma ana yin rajistar sauran. Idan sunan wata folda ya riga ya zama na wani asusu, ka shigo da wannan foldar ita kaɗai da `ccam import <folder> --name <name>`. Ana yin rajistar foldojin da aka shigo da su a inda suke; ccam ba ya taɓa motsawa ko kwafar bayanan shiga.

`ccam login --expired` yana kuma shiga asusan da Claude Code ya share shigarsu. `ccam remove` yana fitar da asusun da ccam ke sarrafawa (wanda ccam ya ƙirƙiri foldarsa, ba wanda aka yi rajistarsa da `ccam import` ba) daga shiga kafin foldarsa ta tafi shara, sai dai idan wani asusu yana riƙe da shiga iri ɗaya. Cire asusun da aka shigo da shi yana fitar da shi daga rukunoninsa kuma yana barin kwafin kowane abin da aka raba a foldarsa. A macOS, `ccam rename` yana riƙe sunan foldar, domin Keychain yana ajiye shigar ƙarƙashin hanyar foldar.

`ccam add --api` yana ƙara asusun da ke amfani da API key maimakon shigar Claude; ccam yana karanta key ɗin a kowane ƙaddamarwa kuma ba ya taɓa ajiye shi. Duba [Asusan API key](/features/api).

`ccam rename`, `ccam remove` da umarnan `ccam group` suna ƙin canza asusu yayin da Claude Code ke gudana a matsayin wannan asusun, ko ta yaya aka fara shi (har da zaman bayan fage). Ka rufe zaman (session) tukuna. Suna kuma ƙi idan canjin rukuni bai kammala ba: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` yana buga kowane asusu tare da ramin (slot) sa, gajeren umarninsa, imel da foldarsa. Array ɗinsa na `groups` yana lissafa rukunonin asusun, kuma abin `share` ɗinsa yana danganta kowane abu (`mcp`, `settings`, `memory`, `extensions`, `history`) da rukunin da yake fitowa daga gare shi, ko `local`.

## Gajerun umarni

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Gajeren umarni umarni ne da kai da kanka ka sa masa suna, wanda ke fara Claude Code a matsayin asusu ɗaya. Kamar `ccam work`, yana miƙa kowane argument ga Claude Code ba tare da canji ba, don haka `ccwork --resume` da `ccwork -p "hi"` suna aiki kamar yadda kake tsammani. Mahaɗi ne zuwa ccam a cikin `~/.local/bin`.

Gajerun umarni zaɓi ne. Mataimakin `ccam add` yana tambaya ko a ƙirƙiri ɗaya, `ccam add <name> --shortcut <command>` yana ƙirƙirar ɗaya ba tare da tambaya ba, kuma `ccam list` yana nuna gajeren umarnin kowane asusu, ko `-` idan ba shi da shi.

Sunaye suna amfani da haruffa a-z, 0-9, `.`, `-` da `_`. Suna ba zai iya zama umarnin ccam ba, ko gajeren umarnin wani asusu, ko sunan wani shiri a `PATH` naka, ko kuma fayil a `~/.local/bin` da ba ccam ne ya ƙirƙira shi ba, kamar kwafinka ko mahaɗinka na ccam. Gajeren umarnin da aka gyara da hannu a `config.toml` wanda ya karya waɗannan ƙa'idodi, ko shi ne `claude` ko `ccam`, ko an ba asusu biyu shi, yana sa `config.toml` ya kasa lodawa, tare da ambaton asusun.

`ccam rename` yana riƙe sunan gajeren umarnin, kuma `ccam remove` yana goge gajeren umarnin tare da asusun. `ccam shims install` yana sake ƙirƙirar mahaɗan dukkan asusai, misali bayan ka matsar da ccam; fayil da ke tare hanyar gajeren umarni ɗaya ba ya hana sauran.

## Rabawa

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Rukuni yana raba abubuwa tsakanin asusansa; asusu zai iya kasancewa a rukunoni da dama, amma kowane abu yana fitowa daga rukuni ɗaya kawai. Duba [rabawa da rukunoni](/features/share).

## Kwafar saitin da ke akwai

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` yana kwafar `settings.json`, `CLAUDE.md`, extensions da sabobin MCP na mai amfani daga foldar Claude Code zuwa rukuni, don kowane asusu da ke raba waɗannan abubuwa ya same su. Ana karanta foldar ne kawai.

- `--settings`, `--memory`, `--extensions`, `--mcp` da `--history` suna kwafar abubuwan da ka ambata kawai. Ba tare da su ba, a terminal jerin zaɓi (checklist) yana ba ka damar zaɓa, inda aka yi wa komai alama sai tarihin tattaunawa; a wasu wurare kuma ana kwafar komai sai tarihin tattaunawa.
- `--mcp` yana aiki ko da `~/.claude` babu komai a ciki ko babu shi, amma `~/.claude.json` yana da sabobi.
- `-g` ko `--group` yana zaɓar rukunin.
- `-y` ko `--yes` yana kwafa ba tare da tambaya ba kuma yana riƙe kwafin rukunin idan fayil ya bambanta.

Ba a kwafar shigar sabobin MCP, domin shigar da aka kwafa tana fitar da sauran masu riƙe da ita; tsarin yana ambaton asusan da za su shiga da `/mcp`. Duba [kwafar saitin da ke akwai](/docs/sharing#copy-an-existing-setup).

## Sabobin MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` yana motsa sabobin da wannan asusun kaɗai ke da su zuwa rukuninsa, har da waɗanda yake da su don folda ɗaya kawai (`claude mcp add` na yau da kullum), waɗanda daga nan suke lodawa a kowace folda. Ana tsallake sabar da rukunin ya riga ya yi amfani da sunanta don wata ma'ana daban; ƙara `--force` don maye gurbin kwafin rukunin, wanda ccam ke fara ajiye masa madadi.

## Amfani da sauyawa

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam sends no token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --continue --brief    # a new conversation from a summary of the last one
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agent "…"                  # from Claude Code: a subagent that runs as another account
ccam agent --resume <id> "…"    # follow that agent up
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
ccam ccusage claude daily       # Claude Code's tokens and costs, every account (ccusage)
ccam ccusage -a work blocks --active   # one account's current 5-hour block
```

- `next --auto`, `--fallback` da `ccam pick` suna zaɓa ta hanya ɗaya: na farko waɗanda ke ƙasa da `auto_max_usage` (`config.toml`, tsoho 100), bisa `auto_strategy` (`headroom`, `expiring` ko `balance`), suna ƙirga tagar mako ta `--model`. Ba sa taɓa zaɓar asusun da ya kai 100% (tagar da Claude Code ya ba da rahoton ta cika ana ƙirga ta a matsayin 100% har zuwa sake saitinta) ko shigar API key, kuma suna zama a cikin ƙungiya (organization) ta asusun da suka fara daga gare shi (`--cross-org` ko `cross_org = true` don duba gaba da haka).
- `--wait` yana ɗaukar ƙimarsa bayan `=` (`--wait=2h`); idan an rubuta shi da sarari ana ƙin shi tare da shawara. Idan ba za a iya zaɓar kowanne ba, suna fita da code 75, ko kuma, da `--wait`, suna jiran asusun farko ya sake saitawa.
- `ccam next --continue` yana kai tattaunawa gaba har ma tsakanin asusan da ba sa raba tarihi: a lokacin Claude Code yana kwafar ta cikin sabon asusun. Da `--brief`, sabon asusun yana fara sabuwar tattaunawa maimakon haka, daga taƙaitaccen bayanin da ccam ya rubuta na ta ƙarshe ba tare da model ba ([`--brief`](/features/next#brief)).
- Ba kamar `ccam <name>` ba, `ccam next` yana karanta zaɓuɓɓukansa da kansa, don haka zaɓuɓɓukan Claude Code suna zuwa bayan `--`; ba tare da shi ba ccam yana tsayawa da `unknown flag`.
- `ccam agent` yana gudanar da aiki da `claude -p` a cikin wannan foldar a matsayin wani asusu banda wanda aka kira shi daga gare shi, a cikin yanayin izini (permission mode) da matakin ƙoƙari (effort) na zaman da ya kira shi, kuma yana buga amsar Claude ta ƙarshe. A iyakar wannan asusun, tattaunawar tana ci gaba a matsayin wani asusu ([Subagents a kan wasu asusai](/features/agent)).
- A cikin Claude Code, `/ccam:switch [account]` na plugin ɗin yana miƙa tattaunawa gaba, kuma `ccam hook` shi ne abin da hooks ɗinsa ke gudanarwa.
- Da `auto_switch = true` a `config.toml` (Linux da macOS; a kashe ta tsohuwa), hook na iyaka yana kuma kawo ƙarshen zaman mu'amala da ccam ya fara kuma yana ci gaba da shi a wani asusu mai sarari, amma iri ɗaya kaɗai: wata kujerar ƙungiya iri ɗaya, ko wani tsarin kai ([`auto_switch`](/features/next#auto-switch)).

Duba [Sauyawa zuwa asusu na gaba](/features/next), [Zaɓar asusu don rubutun umarni](/features/pick), [Subagents a kan wasu asusai](/features/agent) da [Zaman bayan fage na kowane asusu](/features/agents).

## Kulawa

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam uninstall       # remove every account, ccam's folders and ccam itself
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` yana duba checksum na SHA-256 na sakin kafin ya maye gurbin ccam, kuma `--version v0.2.0` yana girka takamaiman saki. Kwafin ccam da aka girka daga kunshin `.deb` ana sabunta shi da apt a maimakon haka. Duba [sabuntawa](/docs/install#update).

`ccam uninstall` yana lissafa abin da yake cirewa kuma yana tambaya tukuna: kowane asusu kamar yadda `ccam remove` ke yi, gajerun umarnin, foldojin ccam da umarnin `ccam`. Foldar da aka shigo da ita tana nan a inda take. Duba [cirewa](/docs/install#uninstall).

## Fayiloli

| Hanya | Abin da ke ciki |
| --- | --- |
| `~/.config/ccam/config.toml` | Asusai, rukunoni da saituna (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). Saituna suna zuwa a sama, kafin `[[group]]`, `[[profile]]` ko `[[folder]]` na farko: saitin da ke ƙarƙashin tebur (table) na wannan tebur ne, don haka ccam yana gargaɗi game da shi a kowane ƙaddamarwa da a `ccam doctor`, kuma umarnan da ke ajiye config suna ƙi har sai an matsar da shi. Ba a riƙe sharhi (comments) idan ccam ya ajiye fayil ɗin |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` na kowane asusu |
| `~/.local/share/ccam/groups/<group>/` | Saituna, sabobin MCP, extensions da tarihi da aka raba |
| `~/.local/state/ccam/` | Cache na amfani, rikodin iyaka, kulle-kulle, bayanai da sauran yanayi (duba jerin da ke ƙarƙashin teburin) |
| `~/.local/bin/<shortcut>` | Gajerun umarni kamar `ccwork`, mahaɗai zuwa ccam |

Foldar yanayi (state) tana ɗauke da:

- cache na amfani da rikodin iyaka (`usage/<name>.json`, `usage/<name>.limit.json`);
- lokutan amfani na ƙarshe;
- alamomin zama da kulle-kulle (`locks/`);
- gajerun umarnin da ccam ya ƙirƙira (`shortcuts.json`);
- canjin rukuni da bai kammala ba (`group-change.json`);
- bayanan da `/ccam:switch` da hook na iyaka ke barwa (`handoff/<account>/<session>.json`);
- lokacin da `ccam pick` ko `auto_switch` suka zaɓi kowane asusu na ƙarshe (`picks/<name>`; zaɓi yana ƙirguwa na minti 2);
- taƙaitattun bayanan `--brief`, da ake ajiyewa na kwana 7 (`brief/<session>.md`).

Saita `XDG_CONFIG_HOME`, `XDG_DATA_HOME` ko `XDG_STATE_HOME` yana canza wurin waɗannan foldojin. A Windows, `~` shi ne `%USERPROFILE%`, kuma gajerun umarni su ne `<shortcut>.exe`.
