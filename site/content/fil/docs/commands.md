---
title: Sanggunian ng mga command
metaTitle: Sanggunian ng mga command ng ccam
description: Bawat command at option ng ccam, mula sa ccam add at ccam <name> hanggang ccam usage, ccam next at ccam doctor.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: 28e1ad1dcbad
---

Ilagay ang `--dry-run` bago ang isang command para makita kung ano ang gagawin nito nang hindi ito talagang ginagawa, halimbawa `ccam --dry-run work`. Gumagana ito kapag sinisimulan ang Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) at sa `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` at `doctor --fix`. Tinatanggap ito ng mga command na nagbabasa lang, gaya ng `list`, `usage`, `share` at `ext`, pero wala naman silang binabago. Tinatanggihan ito ng iba sa halip na tumakbo, kasama ang `group rename`.

<!-- commands -->

## Pagsisimula ng Claude Code

Sinisimulan ng `ccam <name|number> [arguments]` ang Claude Code bilang account na iyon. Hindi binabasa ng ccam ang mga argument na kasunod ng pangalan; ipinapasa nito ang mga iyon sa `claude` sa parehong pagkakasunod-sunod, at kumikilos ang standard input, exit code at mga signal nang eksaktong gaya ng sa `claude` mismo. Opsyonal ang `--` pagkatapos ng pangalan.

Ang `ccam` na may mga argument pero walang pangalan ng account, gaya ng `ccam -p "summarise"`, ay gumagamit ng account na itinakda para sa kasalukuyang folder, o kung wala, ng default na account. Itakda ang mga ito gamit ang `ccam folder <name> [folder]` at `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Para sa mga `-p`/`--print` run ang `--fallback` bago ang pangalan ng account. Kapag naabot ng account ang usage limit nito, pinapatakbo ulit ng ccam ang parehong command bilang naka-log in na account na may pinakamaraming natitira at ipi-print lang ang output ng run na natapos. Tingnan ang [Mga account ayon sa folder](/features/folder) at [Patakbuhin ulit sa limit](/features/next#run-again-at-the-limit).

**Mga seat ng Team at Enterprise.** Lumilipat din ang awtomatikong pagpili ng account (`ccam next`, `--fallback`, `ccam pick`, ang `n` ng dashboard at, kapag binuksan mo ito, ang [`auto_switch`](/features/next#auto-switch)) sa pagitan ng mga seat ng iisang Claude Team o Enterprise na organisasyon, na maaaring lumabag sa terms ng Anthropic. Basahin ang [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) at magtanong muna sa iyong organisasyon at sa Anthropic bago ka umasa rito.

Sa parehong paraan sinisimulan ng isang shortcut gaya ng `ccwork` ang account nito: ginagawa ng `ccwork --resume` ang ginagawa ng `ccam work --resume`. Tingnan ang [Mga shortcut](/docs/commands#shortcuts) sa ibaba.

## Pagdagdag at pag-alis ng mga account

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

Kapag walang folder, hinahanap ng `ccam import` ang `~/.claude` at `~/.claude-*` at nagtatanong tungkol sa bawat isa; iniuulat bilang `Skipped <path>: <reason>` ang folder na pumalya at nirerehistro ang iba. Kung kinuha na ang pangalan ng isang folder, i-import nang hiwalay ang folder na iyon gamit ang `ccam import <folder> --name <name>`. Nirerehistro ang mga in-import na folder kung nasaan sila; hindi kailanman naglilipat o kumokopya ng login ang ccam.

Ina-log in din ng `ccam login --expired` ang mga account na ang login ay binura ng Claude Code. Ina-log out ng `ccam remove` ang account na pinamamahalaan ng ccam (ang account na ang folder ay ginawa ng ccam, hindi ang nirehistro gamit ang `ccam import`) bago mapunta sa trash ang folder nito, maliban kung hawak ng ibang account ang parehong login. Kapag inalis ang in-import na account, inaalis ito sa mga group nito at nag-iiwan ng kopya ng bawat shared na item sa folder nito. Sa macOS, pinapanatili ng `ccam rename` ang pangalan ng folder, dahil iniimbak ng Keychain ang login sa ilalim ng path ng folder.

Nagdadagdag ang `ccam add --api` ng account na gumagamit ng API key sa halip na Claude login; binabasa ng ccam ang key sa bawat pagsisimula at hindi ito kailanman iniimbak. Tingnan ang [Mga API key account](/features/api).

Tumatanggi ang `ccam rename`, `ccam remove` at ang mga command na `ccam group` na baguhin ang isang account habang tumatakbo ang Claude Code bilang account na iyon, paano man ito sinimulan (pati ang mga background session). Isara muna ang session. Tumatanggi rin ang mga ito habang hindi pa tapos ang pagpapalit ng group: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

Ipinapakita ng `ccam list --json` ang bawat account kasama ang slot, shortcut, email at folder nito. Inililista ng `groups` array nito ang mga group ng account, at itinutugma ng `share` object nito ang bawat item (`mcp`, `settings`, `memory`, `extensions`, `history`) sa group na pinanggagalingan nito, o sa `local`.

## Mga shortcut

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Ang shortcut ay isang command na ikaw mismo ang nagpangalan at nagsisimula ng Claude Code bilang isang account. Gaya ng `ccam work`, ipinapasa nito nang walang pagbabago sa Claude Code ang bawat argument, kaya gumagana ang `ccwork --resume` at `ccwork -p "hi"` gaya ng inaasahan mo. Isa itong link sa ccam sa `~/.local/bin`.

Opsyonal ang mga shortcut. Nagtatanong ang wizard ng `ccam add` kung gagawa ng isa, gumagawa ng isa ang `ccam add <name> --shortcut <command>` nang hindi nagtatanong, at ipinapakita ng `ccam list` ang shortcut ng bawat account, o `-` kapag wala.

Puwedeng gamitin sa mga pangalan ang mga character na a-z, 0-9, `.`, `-` at `_`. Hindi puwedeng maging ccam command ang pangalan, shortcut ng ibang account, o pangalan ng ibang program sa iyong `PATH`, ni file sa `~/.local/bin` na hindi ginawa ng ccam, gaya ng sarili mong kopya o link ng ccam. Ang shortcut na mano-manong binago sa `config.toml` na lumalabag sa mga patakarang ito, na `claude` o `ccam`, o na ibinigay sa dalawang account ay nagpapapalya sa pag-load ng `config.toml`, at pinapangalanan ang account.

Pinapanatili ng `ccam rename` ang pangalan ng shortcut, at binubura ng `ccam remove` ang shortcut kasama ng account. Ginagawa muli ng `ccam shims install` ang mga link ng lahat ng account, halimbawa pagkatapos mong ilipat ang ccam; hindi napapahinto ang iba ng file na nakaharang sa isang shortcut.

## Sharing

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Nagshe-share ang isang group ng mga item sa pagitan ng mga account nito; puwedeng kasapi ang isang account sa ilang group, pero iisang group lang ang pinanggagalingan ng bawat item. Tingnan ang [pag-share gamit ang mga group](/features/share).

## Pagkopya ng dati nang setup

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

Kinokopya ng `ccam migrate` ang `settings.json`, `CLAUDE.md`, mga extension at user MCP server mula sa isang Claude Code folder papunta sa isang group, kaya makukuha ang mga ito ng bawat account na nagshe-share ng mga item na iyon. Binabasa lang ang folder.

- Kinokopya lang ng `--settings`, `--memory`, `--extensions`, `--mcp` at `--history` ang mga item na pinangalanan mo. Kung wala ang mga ito, sa terminal ay may checklist para makapili ka, na naka-check ang lahat maliban sa kasaysayan ng usapan; sa ibang lugar, kinokopya ang lahat maliban sa kasaysayan ng usapan.
- Gumagana rin ang `--mcp` kapag walang laman o wala ang `~/.claude` pero may mga server ang `~/.claude.json`.
- Pinipili ng `-g` o `--group` ang group.
- Kumokopya ang `-y` o `--yes` nang hindi nagtatanong at pinapanatili ang kopya ng group kapag magkaiba ang isang file.

Hindi kinokopya ang mga MCP sign-in, dahil nagla-log out sa ibang may hawak ang kinopyang sign-in; pinapangalanan ng plano ang mga account na magsa-sign in gamit ang `/mcp`. Tingnan ang [pagkopya ng dati nang setup](/docs/sharing#copy-an-existing-setup).

## Mga MCP server

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

Inililipat ng `ccam mcp import <account>` sa group nito ang mga server na ang account na iyon lang ang mayroon, pati ang mga mayroon ito para sa iisang folder lang (karaniwang `claude mcp add`), na naglo-load na sa bawat folder pagkatapos. Nilalaktawan ang server na ang pangalan ay ginagamit na ng group para sa ibang definition; idagdag ang `--force` para palitan ang kopya ng group, na bina-back up muna ng ccam.

## Usage at paglipat

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

- Pareho ang paraan ng pagpili ng `next --auto`, `--fallback` at `ccam pick`: inuuna ang mga account na mas mababa sa `auto_max_usage` (`config.toml`, default 100), ayon sa `auto_strategy` (`headroom`, `expiring` o `balance`), binibilang ang lingguhang window ng `--model`. Hindi sila kailanman pumipili ng account na nasa 100% (itinuturing na 100% ang window na iniulat ng Claude Code na puno na hanggang sa reset nito) o ng API key login, at nananatili sila sa loob ng organisasyon ng account na pinagmumulan nila (`--cross-org` o `cross_org = true` para maghanap pa nang mas malayo).
- Inilalagay ang value ng `--wait` pagkatapos ng `=` (`--wait=2h`); kapag isinulat nang may espasyo, tinatanggihan ito kasama ang isang pahiwatig. Kapag walang mapipili, lumalabas sila na may code 75, o, sa `--wait`, naghihintay sa unang account na mag-reset.
- Dinadala ng `ccam next --continue` ang isang usapan kahit sa pagitan ng mga account na hindi nagshe-share ng kasaysayan: kinokopya ito ng Claude Code sa bagong account. Sa `--brief`, nagsisimula sa halip ang bagong account ng bagong usapan, mula sa buod na isinusulat ng ccam tungkol sa huling usapan nang walang model ([`--brief`](/features/next#brief)).
- Hindi tulad ng `ccam <name>`, binabasa ng `ccam next` ang sarili nitong mga option, kaya inilalagay ang mga option ng Claude Code pagkatapos ng `--`; kapag wala ito, humihinto ang ccam sa `unknown flag`.
- Pinapatakbo ng `ccam agent` ang isang task gamit ang `claude -p` sa folder na ito bilang ibang account kaysa sa pinagtawagan dito, sa permission mode at effort level ng tumawag na session, at ipi-print ang huling sagot ni Claude. Sa limit ng account na iyon, nagpapatuloy ang usapan bilang ibang account ([Mga subagent sa ibang account](/features/agent)).
- Sa loob ng Claude Code, ipinapasa ng `/ccam:switch [account]` ng plugin ang isang usapan, at ang `ccam hook` ang pinapatakbo ng mga hook nito.
- Kapag may `auto_switch = true` sa `config.toml` (Linux at macOS; naka-off bilang default), tinatapos din ng limit hook ang interactive session na sinimulan ng ccam at ipinagpapatuloy ito sa ibang account na may natitira, pero sa kauri lang nito: ibang seat ng parehong organisasyon, o ibang personal na plan ([`auto_switch`](/features/next#auto-switch)).

Tingnan ang [Paglipat sa susunod na account](/features/next), ang [Pagpili ng account para sa mga script](/features/pick), ang [Mga subagent sa ibang account](/features/agent) at ang [Mga background session ng bawat account](/features/agents).

## Maintenance

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

Sinusuri ng `ccam update` ang SHA-256 checksum ng release bago nito palitan ang ccam, at nag-i-install ang `--version v0.2.0` ng partikular na release. Ang kopya ng ccam na in-install mula sa `.deb` package ay ina-update na lang gamit ang apt. Tingnan ang [pag-update](/docs/install#update).

Inililista ng `ccam uninstall` ang mga aalisin nito at nagtatanong muna: ang bawat account gaya ng ginagawa ng `ccam remove`, ang mga shortcut command, ang mga folder ng ccam at ang command na `ccam`. Mananatili kung nasaan ito ang in-import na folder. Tingnan ang [pag-uninstall](/docs/install#uninstall).

## Mga file

| Path | Nilalaman |
| --- | --- |
| `~/.config/ccam/config.toml` | Mga account, group at setting (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). Inilalagay ang mga setting sa itaas, bago ang unang `[[group]]`, `[[profile]]` o `[[folder]]`: ang setting na nasa ibaba ng isang table ay pag-aari ng table na iyon, kaya nagbababala tungkol dito ang ccam sa bawat launch at sa `ccam doctor`, at tumatanggi ang mga command na nagse-save ng config hanggang mailipat ito. Hindi pinapanatili ang mga comment kapag sine-save ng ccam ang file |
| `~/.local/share/ccam/profiles/<name>/` | Ang `CLAUDE_CONFIG_DIR` ng bawat account |
| `~/.local/share/ccam/groups/<group>/` | Mga shared na setting, MCP server, extension at kasaysayan |
| `~/.local/state/ccam/` | Usage cache, mga tala ng limit, lock, tala at iba pang state (tingnan ang listahan sa ibaba ng table) |
| `~/.local/bin/<shortcut>` | Mga shortcut command gaya ng `ccwork`, mga link sa ccam |

Ang state folder ay naglalaman ng:

- ang usage cache at mga tala ng limit (`usage/<name>.json`, `usage/<name>.limit.json`);
- mga oras ng huling paggamit;
- mga session marker at lock (`locks/`);
- ang mga shortcut na ginawa ng ccam (`shortcuts.json`);
- ang pagpapalit ng group na hindi pa natatapos (`group-change.json`);
- ang mga talang iniiwan ng `/ccam:switch` at ng limit hook (`handoff/<account>/<session>.json`);
- kung kailan huling pinili ng `ccam pick` o ng `auto_switch` ang bawat account (`picks/<name>`; binibilang ang isang pagpili nang 2 minuto);
- ang mga buod ng `--brief`, na itinatago nang 7 araw (`brief/<session>.md`).

Binabago ng pagtatakda ng `XDG_CONFIG_HOME`, `XDG_DATA_HOME` o `XDG_STATE_HOME` kung nasaan ang mga folder na ito. Sa Windows, ang `~` ay `%USERPROFILE%` at ang mga shortcut ay `<shortcut>.exe`.
