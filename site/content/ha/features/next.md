---
title: Sauyawa zuwa asusu na gaba
description: ccam next yana buɗe Claude Code a matsayin wani asusu; --auto yana zaɓar mafi sarari, --continue yana kai tattaunawar gaba, --fallback yana maimaita ayyukan -p.
updated: 2026-11-03
order: 8
translation: machine
sourceHash: "89e8dea9c55d"
---

`ccam next` yana buɗe Claude Code a matsayin wani asusu ba tare da ka ambaci sunansa ba. Yi amfani da shi idan asusu ya kai iyakar amfaninsa. Akwai hanyoyi huɗu na sauyawa:

- `ccam next` kawai yana ɗaukar asusun da ke bayan wanda ka yi amfani da shi na ƙarshe, kuma `--auto` yana ɗaukar wanda ya fi sarari. `--continue` yana ci gaba da tattaunawar ƙarshe ta wannan foldar a sabon asusun, ko asusun biyu suna raba tarihi ko ba sa raba.
- A cikin Claude Code, [plugin ɗin ccam](#inside-claude-code) yana sanar da kai idan iyaka ta dakatar da zama, kuma `/ccam:switch` yana miƙa tattaunawar gaba.
- Don rubutun umarni (scripts) masu gudanar da `claude -p`, `ccam --fallback` yana yin sauyawar da kansa (duba [Sake gudanarwa a iyaka](#run-again-at-the-limit)).
- Idan an kunna [`auto_switch`](#auto-switch) (Linux da macOS), ccam yana kuma motsa zaman mu'amala da Claude Code ya dakatar a iyakarsa.

**Kujeru (seats) na Team da Enterprise.** Zaɓin asusu na atomatik (`ccam next`, `--fallback`, `ccam pick`, maɓallin `n` na dashboard kuma, da zarar ka kunna shi, [`auto_switch`](/features/next#auto-switch)) yana kuma sauyawa tsakanin kujerun ƙungiyar Claude Team ko Enterprise iri ɗaya, wanda na iya saɓa wa sharuɗɗan Anthropic. Ka karanta [Tambayoyin da ake yawan yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) kuma ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi.


## Amfani

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # continue this folder's last conversation on the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next -c --brief            # a new conversation that starts from a summary of the last one
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the most weekly quota left per day
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work CCAM_PROFILE=work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam yana buga sauyawar a matsayin `→ <from> → <to>` a standard error, tare da dalilin zaɓin `--auto`, sannan yana fara Claude Code a matsayin sabon asusun kamar yadda `ccam <name>` ke yi. `<from>` shi ne asusun da ka yi amfani da shi na ƙarshe, ko da `--continue` asusun da tattaunawar take nasa. Arguments bayan `--` suna zuwa `claude` ba tare da canji ba.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--auto` | Zaɓi asusun da aka shiga da shi bisa amfani (duba ƙasa), ana tsallake asusan da suka kai 100% | a kashe |
| `--max-usage <1-100>` | Da `--auto`: fifita asusan da ke ƙasa da wannan amfani | `auto_max_usage`, in ba haka ba 100 |
| `--strategy headroom\|expiring\|balance` | Da `--auto`: mafi ƙarancin amfanin awa 5, sake saitin awa 5 mafi kusa ko mafi yawan kason mako da ya rage ga kowace rana | `auto_strategy`, in ba haka ba `headroom` |
| `--cross-org` | Da `--auto`: zaɓi asusan wata ƙungiya (organization) dabam ma (duba ƙasa) | `cross_org`, in ba haka ba a kashe |
| `--wait[=<duration>]` | Da `--auto`: idan babu asusun da ke da sarari, jira na farko ya sake saitawa, har zuwa iyakar wannan tsawon. Ƙimar tana zuwa bayan `=`: ana ƙin `--wait 2h` da shawarar `--wait takes its value after '=': --wait=2h` | a kashe; `5h` idan an ba shi shi kaɗai |
| `-c`, `--continue` | Ci gaba da tattaunawar ƙarshe ta wannan foldar a sabon asusun (duba ƙasa) | a kashe |
| `--brief` | Da `--continue`: fara sabuwar tattaunawa a matsayin sabon asusun daga taƙaitaccen bayanin ccam na ta ƙarshe, maimakon ci gaba da dukanta ([duba ƙasa](#brief)) | a kashe |
| `-- <claude args>` | Arguments da ake miƙa wa `claude` bayan na ccam kansa | babu |
| `--dry-run` | Ana tallafa masa: yana buga `→ <from> → <to>` da layin umarnin da zai gudanar, ba tare da fara Claude Code ko rubuta farawa ba | a kashe |

`--auto` da `--continue` suna aiki tare. Da `--continue`, ana saka arguments na `--resume` kafin naka.

Zaɓuɓɓukan Claude Code na kansa, kamar `--dangerously-skip-permissions` ko `--model`, suna zuwa bayan `--`. Ba kamar `ccam <name>` ba, wanda ke miƙa duk abin da ke bayan sunan gaba, `ccam next` yana karanta arguments ɗinsa da kansa, don haka ba tare da `--` ba yana tsayawa da `ccam: unknown flag: --dangerously-skip-permissions`.

## Yadda yake aiki

Kowane farawa ta ccam yana rubuta lokacin zuwa `~/.local/state/ccam/last-used/<name>`. `ccam next` yana karanta waɗannan fayiloli, yana duba asusan da aka shiga da su, kuma yana zaɓa.

1. ccam yana gina `next.Candidate` ɗaya ga kowane asusu daga `config.toml` (`cli.candidates`): asusun, lokacin farawarsa na ƙarshe (`launcher.LastUsed`) da ko yana da shiga mai aiki (`creds.Has`: shigar Claude a `.credentials.json`, ko a macOS a Keychain, wadda ba ta ƙare gaba ɗaya ba ko Claude Code bai share ta ba). Foldar da ke amfani da API key maimakon haka (shigar Console, ko asusun da aka ƙara da `ccam add --api`) ba ta da iyakokin tsari kuma ba a taɓa zaɓar ta ba.
2. Ba tare da `--auto` ba, ccam yana jera asusan bisa lambar rami kuma yana nemo wanda aka fara kwanan nan (`next.Rotate`). Wannan shi ne asusun da kake sauyawa daga gare shi. Idan ba a fara kowane asusu ba tukuna, shi ne mafi ƙanƙantar rami.
3. Daga nan yana tafiya gaba a tsarin ramuka, yana komawa daga rami na ƙarshe zuwa na farko, kuma yana mayar da asusun farko da aka shiga da shi.
4. ccam yana shirya kuma yana fara asusun da aka zaɓa daidai kamar yadda `ccam <name>` zai yi (`launcher.LaunchWith`): yana duba mahaɗan asusun, yana aiwatar da rabawar MCP, yana rubuta zaman kuma yana maye gurbin ccam da `claude` ta `syscall.Exec`.

Tafiyar tana farawa daga asusun da aka yi amfani da shi na ƙarshe kuma tana tsayawa a asusun farko da aka shiga da shi bayan sa:

<!-- motion:next-rotate -->

An fara `work` kwanan nan, don haka tafiyar tana farawa bayan sa. `side` ba shi da shiga kuma ana tsallake shi, kuma ana zaɓar `lab`. Da a ce `lab` ma an fita daga gare shi, tafiyar za ta zagaya zuwa `personal` (`% len` a zanen). Juyawa ba ta taɓa mayar da asusun da ka fara daga gare shi ba: idan babu wani asusu da aka shiga da shi, tana tsayawa da kuskure.

### --auto

`--auto` yana yanke shawara bisa amfani maimakon tsarin ramuka. ccam yana riƙe asusan da aka shiga da su, yana tambayar amfaninsu duka a lokaci ɗaya kamar yadda `ccam usage` ke yi (`Fetcher.FetchAll`), kuma yana zaɓa bisa ƙa'idodin da ke ƙasa (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

Ana auna asusu da tagoginsa na awa 5 da na mako, kuma, idan `--model` yana cikin arguments na Claude Code (bayan `--`), da tagar mako ta wannan model ɗin ma (`opus`, `sonnet`, `haiku` ko `fable` a cikin sunan). ccam yana raba asusan zuwa mataki uku (`Policy.Choose()`):

1. **Sarari**: kowace taga tana ƙasa da `auto_max_usage` (`config.toml`, tsoho 100). Dabarar ce ke zaɓa a cikinsu: `headroom` (tsoho) yana ɗaukar mafi ƙarancin amfanin awa 5, `expiring` yana ɗaukar tagar awa 5 da za ta sake saitawa da wuri, don a yi amfani da sauran quota ɗinta kafin ta ɓace; `balance` yana ɗaukar asusun da ke da mafi yawan quota na mako da ya rage ga kowace rana har makonsa ya sake saitawa, yana ƙirga kowane zaman ccam da ke gudana a kansa, don kada a bar quota na mako ba a yi amfani da shi ba kuma zaman da ke gudana a lokaci ɗaya su rarrabu.
2. **Sama**: ƙasa da 100% amma a ko sama da `auto_max_usage`. Ana amfani da su ne kawai idan babu asusun da ke da sarari; wanda ke da mafi ƙarancin amfanin awa 5 ne yake fara zuwa.
3. **Cike**: taga a 100%, ko amfanin awa 5 da ba a sani ba. Ba a taɓa zaɓar su ba.

Ana ƙirga taga a matsayin 100% har zuwa sake saitinta idan Claude Code ya ƙi buƙata a can: hook na [plugin ɗin ccam](#inside-claude-code) yana rubuta hakan (`~/.local/state/ccam/usage/<name>.limit.json`), ko me lambobin da aka ajiye suka ce. Don haka asusun da ya kai iyakarsa yanzu ba a sake zaɓar sa kafin sake saitinsa, ko da lambobinsa sun tsufa.

Saituna na dukkan asusai suna zuwa a saman `config.toml`, kafin `[[group]]`, `[[profile]]` ko `[[folder]]` na farko. Saitin da ke ƙarƙashin ɗaya daga cikinsu na wannan tebur (table) ne: ccam yana gargaɗi game da shi a kowane farawa, kuma umarnan da ke ajiye config suna ƙi har sai an matsar da shi.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Ƙungiyoyi (organizations).** Zaɓuɓɓukan atomatik suna zama a cikin asusai iri ɗaya da wanda kake barin: kujeru a ƙungiyar Team ko Enterprise iri ɗaya, ko tsare-tsaren kai (Pro, Max), don kada aikin asusun kamfani ya ci gaba a asusun kai ko akasin haka. ccam yana sanin ƙungiyar kowane asusu daga shigarsa (`ccam login` yana rubuta ta; ana karanta tsofaffin shige-shige daga `.claude.json`); asusun da ccam ba zai iya tantance ƙungiyarsa ba za a iya haɗa shi da kowane asusu. Foldoji biyu da suka shiga asusun Claude iri ɗaya suna raba iyakokinsa, don haka ɗaya ba ya taɓa karɓar aiki daga ɗayan. `--cross-org`, ko `cross_org = true` a `config.toml`, yana ɗage ƙa'idar ƙungiya, amma bayan [`/ccam:switch`](#inside-claude-code) motsi zuwa wata ƙungiya dabam yana jiran ka ce eh.

Idan ba za a iya zaɓar komai ba, ccam yana fita da code 75 (`EX_TEMPFAIL`), don rubutun umarni ya bambanta "a sake gwadawa daga baya" da gazawa. Da `--wait`, maimakon haka yana barci har tagar da ke toshe asusun farko da ya cika ta sake saitawa (daƙiƙa 30 bayan lokacin sake saitin, lokacin da endpoint ɗin ke nuna sabuwar tagar), yana buga `ccam: no account has room; waiting until 15:10`, kuma yana sake zaɓa; idan babu asusun da ya faɗi lokacin da zai sake saitawa, yana sake duba bayan minti 10. Idan babu abin da ya samu sarari kafin jiran ya ƙare, yana fita da 75.

Neman yana bin ƙa'idodin `ccam usage`: ana sake amfani da lambobin da ba su kai daƙiƙa 60 ba, neman kowane asusu yana ƙarewa bayan daƙiƙa 3, asusun da bai amsa ba, ko wanda access token ɗinsa ya ƙare, ana auna shi da lambobinsa da aka ajiye, kuma rikodin iyaka yana rinjayar su. Tagar da aka ajiye wadda lokacin sake saitinta ya wuce ana ƙirga ta a matsayin 0%, don haka ba a tsallake asusun da tagarsa ta sake saitawa tun daga lokacin ba. `--auto` ba ya taɓa sabunta token. Ba kamar juyawa ta yau da kullum ba, `--auto` zai iya zaɓar asusun da ka yi amfani da shi na ƙarshe idan har yanzu shi ne ya fi sarari.

### --continue

1. ccam yana neman tattaunawar ƙarshe ta wannan foldar a kowane asusu: a cikin fayilolin `<session>.jsonl` a `projects/<folder>/` (sunan foldar shi ne foldar aiki, bayan warware mahaɗai, inda aka mai da kowane hali da ba harafi ko lamba ba zuwa `-`, kamar yadda Claude Code ke sanya suna), wanda ke da prompt ko amsa mafi sabo. Ana tsallake zaman da ba su da ko ɗaya, kamar zaman da ya gudanar da `/mcp`, `/usage` ko `/status` kawai, ko fayil mara komai. Idan wani zama ya bar bayani (iyakar amfani ko `/ccam:switch`, duba ƙasa) kuma babu tattaunawar da ke da prompt ko amsa fiye da minti 2 bayan sa, ana amfani da ainihin wannan zaman maimakon haka; layukan da Claude Code ke rubutawa a `/exit` ba sa ƙirguwa.
2. Asusun da take nasa shi ne wanda kake sauyawa daga gare shi, kuma ana zaɓar sabon asusun kamar yadda yake a sama, a cikin ƙungiya ɗaya. Idan asusai da dama suna raba wannan tarihin, ta wanda `.claude.json` nasa ke cewa shi ne ya gudanar da tattaunawar na ƙarshe ce (`lastSessionId`), in ba haka ba ta wanda aka yi amfani da shi na ƙarshe.
3. Idan sabon asusun yana raba wannan tarihin (foldar rukuni ɗaya), ccam yana fara `claude --resume <session id>`, ainihin wannan tattaunawar (`--continue` zai buɗe duk wadda ta fi sabo a foldar).
4. In ba haka ba yana fara `claude --resume <the conversation's file> --fork-session`: Claude Code yana kwafar dukkan tattaunawar, ƙarƙashin sabon session ID, cikin tarihin sabon asusun kansa, kuma ta asali tana zama yadda take. Daga nan `--continue` na sabon asusun yana samun ta.
5. Sai dai idan ka ba da `--agent` bayan `--`, ccam yana kuma miƙa agent ɗin da tattaunawar ke gudana a matsayinsa: saitin `agent` na `.claude/settings.local.json` ko `.claude/settings.json` na foldar ko na `settings.json` na asusun, in ba haka ba agent ɗin da layin `agent-setting` na ƙarshe na tattaunawar ya ambata. Claude Code ba ya dawo da agent ɗin lokacin da yake kwafar tattaunawa (mataki na 4), kuma da zai ci gaba a matsayin agent ɗinsa na tsoho, da kowane kayan aiki. Idan sabon asusun ba ya raba agents, plugins da saitunan tsohon ba (foldojin `agents` da `plugins` iri ɗaya da `settings.json` iri ɗaya, ta rukuni), ccam yana ƙara layi bayan wanda ke magana kan motsin: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

ccam yana faɗin girman tattaunawar: sabon asusun yana sake karanta dukanta da amsarsa ta farko. Ba a taɓa raba prompt caches tsakanin ƙungiyoyi ba, kuma kowane tsarin kai ƙungiya ce ta kansa, don haka idan asusan biyu na ƙungiyoyi daban-daban ne ccam yana ƙara `prompt caches are not shared between organizations, so the first reply reads it all again`. Idan tsohon asusun ya shiga sabobin MCP da sabon bai shiga ba, ccam yana ambaton su: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. Raba `history` a rukuni yana guje wa kwafin, yana riƙe tattaunawa ɗaya, kuma yana raba checkpoints da `/rewind` ke maido da lamba daga gare su (`file-history/`).

### --brief

`--continue` yana ci gaba da dukan tattaunawar, kuma sabon asusun yana sake karanta dukanta da amsarsa ta farko. `ccam next --continue --brief` yana fara sabuwar tattaunawa maimakon haka, daga taƙaitaccen bayanin da ccam ya rubuta na tsohuwar ba tare da tambayar model ba:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

Taƙaitaccen bayanin yana ɗauke da waɗannan, a wannan tsari kuma sai inda akwai abin faɗi:

1. taken tattaunawar;
2. hanyar transcript ɗinta, don Claude ya bincika idan yana buƙatar wani dalla-dalla;
3. buƙatarka ta farko, idan sun fi huɗu;
4. taƙaitaccen bayanin da Claude Code ya rubuta lokacin da ya matse (compact) tattaunawar na ƙarshe;
5. buƙatunka huɗu na ƙarshe, tare da umarnan `!` da slash commands kamar yadda ka rubuta su;
6. amsar Claude ta ƙarshe;
7. jerin ayyuka (todo list);
8. fayiloli 10 da Claude ya canza na ƙarshe;
9. `git status --short --branch` na foldar, har zuwa haruffa 800.

Jerin ayyukan yana riƙe abubuwa 10 a mafi yawa: abubuwan da aka kammala ne ake fara cirewa, sannan matakan da suka fi nisa a gaba, don haka matakin yanzu da na gaba suna zama. Ana yanke dogayen sassa a tsakiya, kuma dukan taƙaitaccen bayanin yana zama ƙasa da haruffa 16,000.

ccam yana ajiye shi a matsayin `~/.local/state/ccam/brief/<session>.md`, wanda kai kaɗai za ka iya karantawa, kuma `--brief` na gaba yana cire taƙaitattun bayanan da suka wuce kwana 7. Claude Code yana farawa da prompt ɗaya, kuma yana haɗa fayil ɗin da yake ambata: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` Prompt da ka bayar bayan `--` yana maye gurbin jimla ta ƙarshe (`Then: <your prompt>`), kuma sauran zaɓuɓɓukanka bayan `--` suna bin sa.

Claude Code yana kawo ƙarshen ambato (mention) a `"` kuma yana yanke shi a `#`, don haka idan hanyar taƙaitaccen bayanin tana da ɗayansu, ccam yana ambaton sa da hanyarsa dangane da foldar aiki; idan wannan hanyar ma tana da ɗaya, prompt ɗin yana ambaton fayil ɗin da suna, kuma Claude yana karanta shi da kayan aikinsa na Read, wanda ke fara tambaya saboda fayil ɗin yana wajen foldojin aiki (a gudun `-p` a yanayin Manual, Claude Code yana ƙin karantawar kuma Claude yana ci gaba ba tare da taƙaitaccen bayanin ba).

`ccam --dry-run next --continue --brief` yana buga taƙaitaccen bayanin kuma ba ya ajiye komai. Tsohuwar tattaunawar tana zama yadda take, kuma sabuwar ta san abin da taƙaitaccen bayanin ya faɗa kawai, don haka yi amfani da `--brief` idan aikin ya fi muhimmanci fiye da kowane dalla-dalla na yadda ka kai can. Ba tare da `--continue` ba, ccam yana tsayawa da `--brief works with --continue`; idan ba zai iya karanta tattaunawar ba ko bai sami komai a cikinta da zai taƙaita ba, yana tsayawa kuma yana cewa `ccam next --continue` ba tare da `--brief` ba yana ci gaba da ita gaba ɗaya. [Relay](#relay) da [`auto_switch`](#auto-switch) koyaushe suna ci gaba da dukan tattaunawar. Zaɓin `--brief` na Claude Code da kansa yana zuwa bayan `--`, kamar sauran zaɓuɓɓukansa.

## A cikin Claude Code

[Plugin ɗin ccam](/docs/install#claude-code-plugin) yana ƙara hooks biyu (Claude Code 2.1.78 ko sabo):

- Idan juyi ya tsaya a iyakar amfani, yana rubuta bayanin tattaunawar kuma yana nuna sanarwar desktop a terminals da ke goyon bayanta (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  Yana sanin kowane saƙo da Claude Code ke ƙirgawa a matsayin iyakar amfani: iyakokin awa 5, na mako, na yau da kullum da na model (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` ga kason model), credits na amfani da iyakokin kashe kuɗi, da ƙungiyar da ta dakatar da asusun.

  Idan Claude Code ya rubuta buƙatar da aka ƙi a cikin tattaunawar, hook ɗin yana kuma rubuta cikakkiyar tagar har zuwa sake saitinta ([rikodin iyaka](/features/usage#the-limit-record)), sai dai dakatarwar da ƙungiya ta saita, wadda admins ɗinta za su iya ɗagewa a kowane lokaci. Takaita gudu na ɗan lokaci daga sabar (throttle) ko iyakar saurin (rate limit) na API key ba iyakar amfani ba ne, kuma ana barin su, haka ma zaman da babu mai kula da su (gudun `-p` da zaman bayan fage). Subagent ko agent na workflow da ya tsaya a iyakar amfani ana ƙirga shi a matsayin dakatarwar zamansa, ko da babban tattaunawar tana jiran sa ne kawai: Claude Code yana rubuta wannan buƙatar da aka ƙi a cikin fayil na agent ɗin na kansa a foldar `subagents` na zaman, kuma hook ɗin yana karanta ta a can. Idan agents da yawa suka tsaya a lokaci guda, zaman yana samun bayani ɗaya kawai da sanarwa ɗaya kawai.

  Claude Code 2.1.234 da sabo kuma yana jira a cikin zaman kuma yana ci gaba da kansa bayan sake saitin awa 5; sanarwar tana ba ka zaɓin sauyawa maimakon haka, kuma idan [`auto_switch`](#auto-switch) yana kunne (Linux da macOS), ccam yana sauyawa a madadinka idan wani asusu yana da sarari. Idan babu wani asusu da zai iya karɓa, saboda babu wani asusu da ke cikin ƙungiyar Team ko Enterprise iri ɗaya, ko babu wanda ke wani tsarin kai (idan `cross_org` yana kunne, kowane asusun Claude dabam zai iya), sanarwar tana ƙarewa da `You have no other personal plans, so ccam cannot move this conversation.` ko `You have no other accounts in Acme, so ccam cannot move this conversation.` maimakon haka.
- `/ccam:switch [account]` yana rubuta bayanin tattaunawar kuma yana kawo ƙarshen juyin kafin a yi kowace buƙata: rubuta `/exit`, sannan ka gudanar da `ccam next --continue` don ci gaba da ita a asusun da ka ambata (ko a na gaba). A zaman da ba shi da tattaunawa tukuna, yana cewa babu abin da za a ci gaba da shi. Abin da kai ka rubuta ne kawai ke gudanar da shi: daga plugin 0.3.2, Claude ba zai iya gudanar da `/ccam:switch` ba, kuma aikin da Claude ya tsara lokacinsa (CronCreate, ScheduleWakeup ko `.claude/scheduled_tasks.json`) yana miƙa shi ga Claude a matsayin rubutu na yau da kullum.

  Idan asusun da ka ambata yana wajen ƙungiyar tattaunawar (kujerar wata ƙungiyar Team ko Enterprise dabam, tsarin kai ga tattaunawar da ke kan kujera ko akasin haka, ko asusun API key a kowane gefe), `ccam next --continue` yana fara tambaya, ko me `cross_org` ya ce: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. Yana yin irin wannan tambayar kafin asusun da shi kansa ya zaɓa idan ba ka ambaci ko ɗaya ba, wanda `cross_org = true` ke ba shi damar ɗauka daga wata ƙungiya dabam: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. `y` ko `yes` ne kawai ke motsa ta; a Linux da macOS, ana fara watsar da maɓallan da aka danna kafin tambayar. Duk wata amsa dabam, ko rashin terminal da za a yi tambayar a kai, yana barin tattaunawar da bayanin a inda suke: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (idan ccam ne ya zaɓi asusun, `ccam next --continue` yana sake zaɓa). Kamar yadda yake a zaɓuɓɓukan atomatik, asusun da ccam ba zai iya tantance ƙungiyarsa ba ya dace da kowane.

Bayanin yana zama a `~/.local/state/ccam/handoff/<account>/<session>.json`, fayil ɗaya ga kowane zama, don haka zama biyu na asusu ɗaya ba sa taɓa rubuta kan na juna. Ana riƙe shi na awa 12 kuma ana amfani da shi sau ɗaya, lokacin da zama na gaba ke shirin farawa; farawar da ta gaza tana barin sa don ƙoƙari na gaba. Na foldar da zaman ya fara a cikinta ne, ko bayan `cd` a cikin shell na Claude Code. `ccam rename` yana sabunta sunayen da ke cikinsa; idan an cire asusun da ya ambata tun daga lokacin, ccam yana faɗin haka kuma yana zaɓa kamar yadda ya saba.

### Relay

Da `relay = true` a saman `config.toml`, ccam yana ci gaba da aiki tare da zaman mu'amala maimakon maye gurbin kansa da Claude Code, kuma bayan `/ccam:switch` yana ci gaba da tattaunawar wannan zaman da kansa (ba ya aiki da bayanin da wani zaman ya bari): hook ɗin yana rufe Claude Code, kuma ccam yana gudanar da `ccam next --continue` da bayanin wannan zaman. Idan ba a ambaci asusu ba, yana ɗaukar wanda ya fi sarari, yana ƙirga model na tattaunawar, kuma idan `cross_org = true` ya bar wannan ya zama asusun wata ƙungiya dabam, yana yin tambayar da ke sama kafin ya motsa tattaunawar zuwa can. A Windows hook ɗin ba zai iya rufe Claude Code ba: rubuta `/exit`, kuma ccam yana ci gaba da ita. Haka kuma hook ɗin ba ya rufe shi don asusun da ke wajen ƙungiyar tattaunawar: yana cewa `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`, kuma ccam yana yin tambayar da ke sama bayan ka rubuta `/exit`. Hakan yana cin process ɗaya na ccam mai jira ga kowane zama, shi ya sa yake a kashe ta tsohuwa.

Zaman da aka miƙa yana riƙe abin da ci gaba da tattaunawa (resume) ba ya maidowa: zaɓuɓɓuka kamar `--add-dir`, `--mcp-config`, jerin kayan aiki (tools), `--settings`, `--plugin-dir` da zaɓuɓɓukan system prompt, yanayin izini (permission mode) na tattaunawar, matakin ƙoƙari (effort level) na zaman da agent ɗin da tattaunawar ke gudana a matsayinsa (duba [`--continue`](#continue)). Claude Code da kansa yana maido da `--model` da `--name` daga tattaunawar. Duk wani abu dabam (`--worktree`, `--debug`, `--remote-control`, zaɓuɓɓukan da ccam bai sani ba) ana ambaton sa a `ccam: not carried over: …`. Sauyawar ba ta taɓa yin tambayar tarihin da aka raba ba, kuma abin da ka rubuta yayin da take gudana ana watsar da shi; tambayarta ɗaya tak ita ce wadda ke zuwa kafin asusun wata ƙungiya dabam.

Relay yana ci gaba da tattaunawar ne kawai idan Claude Code ya ƙare yadda `/ccam:switch` ya buƙace shi: ba idan ka rubuta `/exit` da farko ba, ko ka danna Ctrl+C yayin da yake ƙarewa, ko ccam ya sami hangup ko wani signal, kuma ba idan terminal ɗin ya tafi ba. Sannan bayanin yana zama, kuma `ccam next --continue` yana ci gaba da ita daga baya. Zaman da aka fara ƙarƙashin relay na tsohuwar sigar ccam, ko wanda babu mai kula da shi (gudun `-p` da zaman bayan fage), hook ɗin ba ya rufe shi; yana samun saƙon `/exit` maimakon haka. `claude remote-control` yana gudana ba tare da relay ba.

### `auto_switch`

Da `auto_switch = true` a saman `config.toml` (Linux da macOS), ccam yana motsa zama zuwa wani asusu da kansa idan Claude Code ya dakatar da shi a iyakar amfani: hook na plugin ɗin yana zaɓar wani asusu mai sarari kuma yana rufe Claude Code, kuma ccam, wanda ke ci gaba da aiki tare da zaman, yana ci gaba da tattaunawar a can, kamar yadda [relay](#relay) ke yi bayan `/ccam:switch`. A kashe yake ta tsohuwa kuma na gwaji ne.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam yana motsa zama ne kawai idan duk waɗannan sun cika:

1. ccam ne ya fara shi a terminal kuma yana ci gaba da aiki tare da shi. `auto_switch` yana kunna [relay](#relay), don haka wannan yana cika ga zaman mu'amala da ka fara bayan ka saita shi. Zaman da aka fara kafin haka, ko ƙarƙashin tsohon ccam, yana samun sanarwa da ke ƙarewa da ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` Ba a taɓa motsa zaman da babu mai kula da su (gudun `-p` da zaman bayan fage), `claude remote-control` da Claude Code da aka fara ba ta ccam ba.
2. Asusun yana da [plugin ɗin ccam](/docs/install#claude-code-plugin) a kunne kuma yana gudanar da Claude Code 2.1.214 ko sabo, wanda ke miƙa process ID ɗinsa ga hook ɗin (`CLAUDE_PID`).
3. Claude Code ya rubuta buƙatar da aka ƙi a cikin tattaunawar, tare da lokacin sake saitinta, don iyakar awa 5, ta mako, ta Opus ko ta Sonnet na wannan asusun. Iyakar kashe kuɗi, ƙarewar credits na amfani, dakatarwar da ƙungiya ta saita, kason Fable (`You've reached your Fable limit.`), takaita gudu daga sabar (throttle) da iyakar saurin (rate limit) na API key ba sa taɓa motsa zama. Iyakar subagent ko ta agent na workflow tana ƙirguwa, ko da babban tattaunawar tana jiran sa ne kawai; idan iyakar Opus ko Sonnet ta mako ce, zaman yana motsawa ne kawai idan tattaunawar kanta tana gudana a kan wannan model ɗin.
4. ccam ya san tsarin (plan) asusun, daga `.claude.json` ɗinsa ko daga abin da `ccam login` ya rubuta. Idan bai sani ba, sanarwar tana ƙarewa da `ccam does not know personal's plan, so auto_switch leaves this session.`
5. Zaman ba a haɗa shi da Remote Control ba. A matsayin wani asusu ba zai sake haɗuwa da manhajar da kake sarrafa shi daga gare ta ba, don haka yana zama, kuma sanarwar tana faɗin haka.
6. Iyakar tana sake saitawa fiye da minti 15 daga yanzu. Don jira mafi guntu ccam yana barin zaman ga Claude Code: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. Ba ka rubuta prompt ba tun bayan iyakar.
8. Wani asusu yana da sarari. ccam yana zaɓa kamar yadda [`--auto`](#auto) ke yi (`auto_max_usage`, `auto_strategy`, model na tattaunawar), a cikin asusai iri ɗaya da wanda kake barin: wata kujerar ƙungiyar Team ko Enterprise iri ɗaya, ko wani tsarin kai (Pro, Max). `cross_org` da `--cross-org` ba sa aiki a nan. Dole ne asusun:
   - ya kasance an shiga da shi da shigar da Claude Code zai iya sabuntawa;
   - kada ya kasance an shiga asusun Claude iri ɗaya da shi;
   - kada ya yi amfani da API key;
   - ya kasance yana da tsarin da ccam ya sani;
   - kada credits na amfani su kasance a kunne, a iya abin da binciken amfani na ccam da kansa ke nunawa;
   - ya riga ya amince da wannan foldar a Claude Code;
   - idan tattaunawar tana gudana a matsayin agent (saitin `agent` ko tattaunawar tana ambaton ɗaya), ya raba agents na asusun da kake barin: dole ne foldojinsa na `agents` da `plugins` da `settings.json` nasa kowanne ya zama ainihin iri ɗaya (ta rukuni), ko kuma babu su a duka biyun.

   Idan waɗannan uku na ƙarshe suka cire kowane asusu mai sarari, sanarwar tana ambaton su: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` ko `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` Tagar da ke da [rikodin iyaka](/features/usage#the-limit-record) ana ƙirga ta a matsayin cikakkiya, kuma asusun da wani zama ya motsa zuwa gare shi a cikin minti 2 na ƙarshe yana zuwa bayan sauran. Dole ne a yi zaɓin a cikin kusan daƙiƙa 6; idan binciken amfani ya ɗauki lokaci fiye da haka, zaman yana zama. Da `usage_source = "claude"`, `auto_switch` yana yanke shawara daga lambobin da aka ajiye kawai, kuma ccam ba zai iya sani daga binciken Claude Code ko layin matsayi ko credits na amfani suna kunne ba, don haka yana ɗauka cewa a kashe suke.

Sannan hook ɗin yana rufe Claude Code, kuma ccam, wanda ke aiki tare da shi, yana nuna sanarwar desktop a terminals da ke goyon bayanta, kamar `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` Yana buga layin farko na misalin output da ke sama kuma yana gudanar da `ccam next --continue` don wannan zaman a matsayin asusun da aka zaɓa, wanda ke buga sauran layukansa biyu kuma yana fara Claude Code tare da zaɓuɓɓukan zaman, yanayin izini (permission mode), matakin ƙoƙari (effort level) da agent, kamar yadda [relay](#relay) ke yi. Saƙon farko na sabon zaman wani rubutu ne da ba ya canzawa, abu ɗaya tak da ccam ke aikawa a madadinka (motsi bayan `/ccam:switch` ba ya aika komai):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

Bayan iyakar subagent, saƙon shi ne wannan, domin subagents, workflows da umarnin bayan fage sun tsaya tare da tsohon zaman:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

Idan wannan asusun ma ya kai iyakarsa, haka ne ke faruwa a can. Rikodin iyaka suna hana sarƙar komawa zuwa wani asusu kafin sake saitinsa, don haka idan kowane asusu ya cika zaman yana zama a na ƙarshe.

Hook ɗin yana sanin cewa tattaunawa tana gudana a matsayin agent daga saitin `agent`, ko daga layin da ke ambaton agent ɗin wanda Claude Code ke rubutawa cikin tattaunawar lokaci-lokaci da kuma lokacin da zaman ya ƙare. Idan har yanzu `--agent` na zaman ne kawai ke ambaton sa, kuma asusun da aka zaɓa ba ya raba agents ɗin, `ccam next --continue` ba ya fara komai kuma yana riƙe bayanin, kuma sanarwar desktop ta biyu tana faɗin haka: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

Idan babu wani asusu mai sarari, babu abin da ke motsawa: Claude Code yana kula da iyakar kamar yadda aka saba, yana jira a cikin zaman ko yana ba da zaɓuɓɓukansa, kuma sanarwar tana faɗin wane asusu ne zai fara samun sarari, misali `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` Ga kujeru tana cewa `None of your other accounts in Acme has room now`. Don dakatar da motsin da ya fara, danna Ctrl+C yayin da Claude Code ke rufewa: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

Motsi yana da farashi. Claude Code yana sake farawa, wanda ke ɗaukar 'yan daƙiƙoƙi, kuma shells na bayan fage, dev servers da agents da Claude ya fara suna tsayawa tare da tsohon zaman. Sabon asusun yana sake karanta dukan tattaunawar da amsarsa ta farko, kamar da [`--continue`](#continue), kuma ccam yana faɗin girmanta; motsi ba ya taɓa amfani da [`--brief`](#brief).

`auto_switch` ba ya taɓa sauyawa kafin Claude Code ya ba da rahoton iyaka, ba ya taɓa sa ido kan amfaninka a bayan fage, kuma ba ya taɓa kwafa, sabuntawa ko miƙa shiga. A Windows ana watsi da shi, kuma `ccam doctor` yana faɗin haka; a can, da `relay = true`, `/ccam:switch` da `/exit` har yanzu suna motsa tattaunawa zuwa wani asusu.

**Sharuɗɗan Anthropic.** A kan kujerun Team da Enterprise, `auto_switch` yana motsa aikinka zuwa wata kujerar ƙungiya iri ɗaya. Cibiyar taimako ta Anthropic ta ce "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (ana ba da kowace kujera a tsarin Team da sunan takamaiman mutum ɗaya, kuma wannan mutumin kaɗai ne aka yarda ya yi amfani da ita) da "Users can't route their work through another user's account to use that user's remaining usage" (masu amfani ba za su iya tura aikinsu ta asusun wani mai amfani ba domin su ci gajiyar ragowar amfanin wannan mai amfanin) ([Tsare-tsaren Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); haka ma ga [tsare-tsaren Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)).

Ga tsare-tsaren kai (Pro, Max), sharuɗɗan Anthropic ba su faɗi ko asusanka za su iya karɓar aiki daga junansu ta atomatik a iyaka ba. [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) sun hana samun damar Services "through automated or non-human means, whether through a bot, script, or otherwise" (ta hanyoyin atomatik ko waɗanda ba na ɗan Adam ba, ko ta bot, script, ko wata hanya) sai dai idan Anthropic ta ba da izini a fili, kuma motsi yana aika ɗaya daga cikin saƙonnin da ba sa canzawa na sama a madadinka. [Usage Policy](https://www.anthropic.com/legal/aup) da za ta fara aiki a ranar 12 ga Nuwamba, 2026 tana lissafa "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (kauce wa ko tsoma baki da gangan a cikin iyawa, ƙuntatawa, ko shingayen kariya da aka kafa a cikin kayayyakinmu) a cikin abubuwan da ta hana.

ccam ba ya iƙirarin cewa `auto_switch` ya bi waɗannan sharuɗɗa. Kunna shi ya rage naka, haka ma haɗarin: Anthropic za ta iya dakatar da asusu, ko damar ƙungiyar Team ko Enterprise, idan ta yi imani an karya sharuɗɗan.

Yayin da `auto_switch` ke kunne, `ccam doctor` yana gargaɗi game da asusun da ba shi da plugin ɗin ccam da kuma sigogin Claude Code da suka gaza 2.1.214, kuma ga kujeru yana maimaita ƙa'idar kujeru ta Anthropic ([bincike](/features/doctor)).

## Sake gudanarwa a iyaka

`ccam next` na zaman mu'amala ne da ya tsaya. Gudu wanda ba na mu'amala ba, `claude -p` a rubutun umarni ko a aikin da aka tsara lokacinsa, ba zai iya jiranka ba, don haka `--fallback` kafin sunan asusun yana bar ccam ya sauya a madadinsa:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` yana aiki ne kawai tare da `-p` ko `--print`; in ba haka ba ccam yana tsayawa kuma yana nuna `ccam next --auto --continue`.
2. Ana karanta standard input da aka tura ta bututu (pipe) sau ɗaya, har zuwa 10 MiB, kuma ana ba kowane yunƙuri shi. Kamar Claude Code, ccam yana jira daƙiƙa 3 a mafi yawa don ya fara (sai dai da `--input-format stream-json`); idan bututu ya kasance a buɗe amma bai aika komai ba, gudun ba ya samun input kuma ccam yana rubuta haka a standard error.
3. ccam yana fara `claude` a matsayin child process maimakon `exec`, yana riƙe standard output ɗinsa, kuma yana miƙa standard error yadda yake zuwa.
4. Idan `claude` ya fita da 0, ko ya gaza saboda wani dalili, ana miƙa output ɗin da aka riƙe da exit code ba tare da canji ba. Haka ma gudun da signal ya ƙare (exit code sama da 128): an nemi ya tsaya, kuma ba a sake gudanar da shi.
5. Gazawa tana ƙirguwa a matsayin iyakar amfani idan saƙon Claude Code ya faɗi haka, kowane daga cikin saƙonnin iyakar amfani da [hook ɗin](#inside-claude-code) ya sani, har a cikin output na `--output-format json` ko `stream-json`, ko idan amfanin asusun ya nuna tagarsa ta awa 5, ta mako ko ta `--model` a 100% (har da rikodin iyaka).
6. Sannan ana watsar da output ɗin wannan yunƙurin kuma umarnin iri ɗaya yana gudana a matsayin wani asusu da aka shiga da shi: wanda ke da mafi ƙarancin amfanin awa 5, sannan asusan da ba za a iya neman amfaninsu ba, bisa tsarin lambobi. Ana tsallake asusan da suka kai 100%, ba a taɓa gwada asusan da ke da shigar API key ba, kuma ana gwada kowane asusu sau ɗaya kawai.
7. Idan akwai `-c`, `--continue`, `-r` ko `--resume` a cikin arguments, asusan da ke rukunin tarihi ɗaya ne kawai ake gwadawa, don a sami tattaunawar.
8. Asusan ƙungiya ɗaya ne kawai ake gwadawa, kamar da `--auto` (`cross_org = true` yana ɗage wannan).

Ana zaɓar asusu na gaba bisa ƙa'idodi iri ɗaya da na `--auto` (`auto_max_usage`, `auto_strategy`, tagar `--model`). Idan babu sauran asusu, ccam yana buga output na yunƙurin ƙarshe, yana cewa babu asusun da aka shiga da shi da ke da sauran sarari, kuma yana fita da code 75. Da `--wait` kafin sunan asusun (`ccam --fallback --wait=2h work -p …`), maimakon haka yana jiran asusun farko ya sake saitawa, sannan zai iya sake gwada kowane asusu; koyaushe yana fara yin barci aƙalla sau ɗaya, don kada a sake gudanar da asusun da ya kai iyakarsa yanzu nan take. Tsawon yana iyakance jimillar lokacin jiran gudun, ba kowane jira ɗaya ba. Idan an rubuta shi da sarari (`--wait 30m`), ana ƙin sa da shawara iri ɗaya da ta `ccam next`. `--` na farko kafin arguments na Claude Code yana aiki kamar a kowane farawa (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` yana buga umarnin yunƙuri na farko.

Ba a sauya zaman mu'amala da kansu sai idan ka kunna [`auto_switch`](#auto-switch). Da `ccam statusline install`, layin matsayi yana ƙara `→ ccam next --auto --continue` da zarar amfanin awa 5 ko na mako ya kai 90% ([Iyakokin amfani](/features/usage)).

## Tsaro da yanayi na musamman

- Ba a taɓa zaɓar asusan da ba su da shiga. Idan babu wani asusu da aka shiga da shi, `ccam next` kawai yana tsayawa da `no logged-in account other than <name>`.
- Da `--continue` kuma babu tattaunawa a foldar: `no account has a conversation in <folder> to continue`.
- Idan ƙa'idar ƙungiya ba ta bar kowane asusu ba: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Da `--auto` kuma babu asusun da aka shiga da shi: `no logged-in account to pick`. Idan kowane asusu ya cika ko ba a san amfaninsa ba: `` no account with known usage and headroom (see `ccam usage`) ``. Dukansu suna ƙarewa da exit code 75.
- `--max-usage`, `--strategy`, `--wait` da `--cross-org` suna aiki ne kawai da `--auto`, kuma ana ƙin ƙima a wajen 1-100 ko dabarar da ba a sani ba.
- Fork yana kwafar tattaunawar, ba sauran yanayin asusun ba: checkpoints na `file-history/` na asusun asali ba sa zuwa tare sai idan an raba tarihi.
- `--auto` yana yin buƙatun network ko da `--dry-run`. Yana karantawa ne kawai, kuma ba ya taɓa rubuta fayilolin token.
- Asusun da aka zaɓa yana farawa ta hanya ɗaya da `ccam <name>`, don haka dubawa iri ɗaya suna aiki, har da gargaɗin idan wani zaman Claude Code da ke raba tarihi ɗaya yana aiki tuni a wannan foldar.
- Ana rubuta `last-used` ne kawai idan Claude Code ya fara da gaske. `--dry-run` yana barin sa, don haka dry run ba ya matsar da juyawar gaba.
- `--fallback` yana riƙe standard output har sai yunƙuri ya ƙare, don haka output mai gudana kamar `--output-format stream-json` yana zuwa gaba ɗaya a ƙarshe. Abokin aiki (client) da ke tura input na `stream-json` kuma yana jiran amsoshi kafin ya aika ƙari ba zai iya amfani da `--fallback` ba. Fayilolin da yunƙurin farko ya canza kafin ya kai iyaka suna nan a canje; yunƙuri na gaba yana sake fara prompt ɗin.

## Masu alaƙa

- [Iyakokin amfani](/features/usage)
- [Iyakokin amfani da sauya asusai](/docs/usage-limits)
- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Claude Code ya kai iyakar amfani? Ga abin da za ka yi](/guides/usage-limit-workflow)
