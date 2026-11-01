---
title: Pagsisimula ng Claude Code
description: Patakbuhin ang Claude Code bilang anumang account sa ccam <name|number>, ccam run o dashboard; buo sa claude ang argument; ipinapakita ng --dry-run ang command.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: 11df08608a1e
---

Sinisimulan ng `ccam work` ang Claude Code bilang `work` account, at ipinapasa sa `claude` nang walang pagbabago ang lahat ng kasunod ng pangalan. Gamitin ito saanman mo ita-type ang `claude`: sa terminal, sa script, o kapag may input na ipina-pipe.

## Paggamit

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Ibigay ang pangalan ng account o ang numero nito mula sa `ccam list`. Hindi binabasa ng ccam ang mga argument na kasunod ng pangalan, kaya ipinapakita ng `ccam work --help` ang help ng Claude Code. Opsyonal ang `--` pagkatapos ng pangalan. Kumikilos ang standard input, exit code at mga signal gaya ng sa `claude`.

Kapag nagsisimula sa `-` ang unang argument at hindi ito isa sa sariling flag ng ccam, pinipili ng ccam ang account para sa iyo at ipinapasa ang lahat ng argument sa Claude Code. Ginagamit nito ang account na itinakda para sa kasalukuyang folder gamit ang [`ccam folder`](/features/folder), o kung wala, ang default na account: ang itinakda gamit ang `ccam use`, o ang account na may pinakamababang numero kung wala kang itinakda.

### Ang dashboard

Binubuksan ng `ccam` na walang argument ang [dashboard](/features/dashboard) kapag terminal ang standard input at output, at ipi-print ang `ccam list` kung hindi. Pumunta sa isang account (gamit ang mga arrow, o sa pag-type ng numero nito) at pindutin ang `Enter` para simulan ito; sinisimulan ng `n` ang account na may pinakamaraming natitira at nag-aalok ang `o` ng `--continue`, `--resume` o iba pang argument. Nagsasara ang dashboard bago magsimula ang Claude Code, kaya naibabalik muna ang terminal.

**Mga seat ng Team at Enterprise.** Lumilipat din ang awtomatikong pagpili ng account (`ccam next`, `--fallback`, `ccam pick`, ang `n` ng dashboard at, kapag binuksan mo ito, ang [`auto_switch`](/features/next#auto-switch)) sa pagitan ng mga seat ng iisang Claude Team o Enterprise na organisasyon, na maaaring lumabag sa terms ng Anthropic. Basahin ang [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) at magtanong muna sa iyong organisasyon at sa Anthropic bago ka umasa rito.

Kapag nakatakda sa shell ang `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` o `ANTHROPIC_AUTH_TOKEN`, maaaring gamitin ng Claude Code ang key o token na iyon sa halip na ang sariling login ng account, kaya binabalaan ka ng ccam sa standard error bago ito magsimula (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Binabasa lang nito ang mga variable; binubura naman ang mga ito ng isang [API key account](/features/api).

Bilang default, pinapalitan ng ccam ang sarili nito ng Claude Code, kaya walang ibang patuloy na tumatakbo. Kapag may `relay = true` (o [`auto_switch = true`](/features/next#auto-switch) sa Linux at macOS) sa itaas ng `config.toml`, bago ang unang `[[group]]`, `[[profile]]` o `[[folder]]`, sinisimulan na lang ng ccam ang mga interactive session bilang child process at patuloy itong tumatakbo kasabay nila. Dahil dito, kusang maipagpapatuloy ng `/ccam:switch`, o ng `auto_switch` sa usage limit, ang usapan sa ibang account, kasama ang mga option, permission mode, effort level at agent ng session; bago ito ilipat ng `/ccam:switch` sa account ng ibang organisasyon, nagtatanong muna sa iyo ang ccam ([mga detalye](/features/next#relay)).

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--dry-run` | I-print ang environment at command na patatakbuhin ng ccam, saka lumabas | off |
| `--fallback` | Para sa mga `-p`/`--print` run: kapag naabot ng account ang usage limit nito, patakbuhin ulit ang parehong command bilang naka-log in na account na may pinakamaraming natitira ([mga detalye](/features/next#run-again-at-the-limit)) | off |
| `-v`, `--version` | I-print ang `ccam <version>` (bilang unang argument lang) | |
| `-h`, `--help` | Ang help ng ccam, bilang `ccam --help` o `ccam run --help` | |

Mga patakaran sa argument:

- Binabasa lang ng ccam ang sarili nitong mga flag bago ang pangalan ng account. Pagkatapos ng pangalan, sa Claude Code na ang bawat argument, kasama ang `--dry-run` at `--help`.
- Gumagana ang `--dry-run` bilang `ccam --dry-run work`, `ccam run --dry-run work` at `ccam --dry-run` (para sa account ng folder o sa default). Puwedeng magkasunod ang `--dry-run` at `--fallback` sa anumang pagkakasunod-sunod.
- Ang `ccam run <name|number> [claude args]` ay ang mahabang anyo ng `ccam <name>`. Ginagamit ng `ccam run` na walang pangalan ang account ng folder o ang default.
- Hindi puwedeng pangalan ng ccam command ang mga pangalan ng account, gaya ng `add`, `list`, `run` o `usage`, kaya hindi kailanman malabo ang `ccam <word>`.

Nagpi-print ang dry run ng isang linya na puwede mong idikit sa shell, o sa PowerShell sa Windows. Nagsisimula ito sa mga binabago ng pagsisimula sa environment: laging `CLAUDE_CONFIG_DIR` at `CCAM_PROFILE`, `env -u <name>` para sa variable na inaalis nito (gaya ng namanang `CCAM_RELAY`), at, para sa isang [API key account](/features/api), ang endpoint nito, ang model nito at kung saan galing ang key (`ANTHROPIC_API_KEY="$WORK_KEY"`), hindi kailanman ang key mismo:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

Ibinabalik ng PowerShell na linya ang mga variable pagkatapos, gaya ng ginagawa ng `VAR=x cmd` sa shell, at ibinibigay sa Claude Code ang command line nito sa pamamagitan ng `Start-Process`, kaya pareho itong pinapatakbo ng Windows PowerShell 5.1 at PowerShell 7. Kapag nagtakda ang pagsisimula ng variable sa walang lamang value, may linyang `#` sa ilalim nito na nagsasabing inaalis na lang ng PowerShell na mas luma sa 7.5 ang ganoong variable. Para sa `claude.cmd` na in-install ng npm, pinapatakbo ng linya ang `cmd.exe` command na gagamitin ng pagsisimula (tingnan sa ibaba), at pumapalya rin ang dry run sa argument na tatanggihan doon ng pagsisimula.

## Paano ito gumagana

Pareho ang mga hakbang ng bawat pagsisimula, mula sa command line hanggang `exec`:

<!-- motion:launch -->

Umiilaw ang bawat kahon habang tumatakbo ang hakbang nito, at ang mga linya sa ilalim ay ang environment at mga argument na kasama ng pagsisimula ng Claude Code.

1. Nagsisimula ang ccam sa `app.Main`, na tinatawag ng `cmd/ccam/main.go`. Kung sinimulan ang ccam sa ibang pangalan, isa itong [shortcut](/features/shortcuts) at direktang napupunta sa launcher.
2. Tinitingnan ng ccam ang `os.Args` bago pa makisangkot ang cobra (`app.parse`). Ang `ccam` lang na mag-isa (ang dashboard), ang pangalan ng ccam command (`config.Reserved`), `-h`/`--help` at `-v`/`--version` ang napupunta sa `cli.Execute`. Pagsisimula ang lahat ng iba pa.
3. Binabasa ng ccam ang `~/.config/ccam/config.toml`, ang tanging file na pina-parse nito (`launcher.Launch`).
4. Hinahanap ng ccam ang account ayon sa pangalan o numero, inaalis ang isang `--` sa unahan, at hinahanap ang `claude` executable (`launcher.Prepare`, `ResolveClaude`). Kapag walang pangalan, kinukuha nito ang account ng folder, na ikinukumpara ang working directory sa mga `[[folder]]` entry ng `config.toml` bilang mga string, nang hindi pinapansin ang malaki at maliit na titik sa macOS at Windows (`HereProfile`). Naka-cache ang paghahanap sa PATH sa `~/.local/state/ccam/claude-path` kasama ang modification time ng file at isang hash ng `PATH`, kaya hindi hinahanap ng ccam ang `PATH` sa bawat pagsisimula. Naghahanap ulit kapag nagbago o nawala ang file, o kapag nagbago ang `PATH`. Dapat executable file ang path na `claude_bin` sa `config.toml` (pinapalawak ang `~/`); kung hindi, humihinto ang pagsisimula na may `claude_bin X is not an executable file. Check claude_bin in config.toml` bago may maisulat. Tinatanggihan ng ccam ang `claude` na ccam mismo, na ikinukumpara ayon sa inode, dahil magdudulot iyon ng loop.
5. Ikinukumpara ng ccam ang folder ng account na `~/.local/share/ccam/profiles/<name>/` sa sharing settings nito, na tumatawag lang ng `lstat` at `readlink`, isa o dalawang system call bawat item (`checkLinks`, `share.Check`). Ginagawa agad ang nawawalang link. Pinapahinto ang pagsisimula ng anumang puwedeng magpawala ng data.
6. Idinadagdag ng ccam ang mga shared na MCP server (`applyMCP`; tingnan sa ibaba).
7. Itinatala ng ccam ang session (`session.Register`). Nagsusulat ito ng marker na ipinangalan sa PID nito, na naglalaman ng oras ng pagsisimula ng proseso, sa `~/.local/state/ccam/locks/profile-<name>/<pid>`, anuman ang sabihin ng `XDG_RUNTIME_DIR`, kaya pareho ang mga marker na nakikita ng bawat shell. Para sa account na nagshe-share ng kasaysayan ng usapan, hinahanap muna ng ccam ang tumatakbong session sa `history-<group>-<hash>`, kung saan ang hash ay sa kasalukuyang folder na naresolba ang mga symbolic link, at nagrerehistro rin doon (`historyCheck`).

   Pinapanatili ng `exec` ang PID, kaya ang marker ay tumutukoy na sa mismong proseso ng Claude Code. Gumagamit ang ccam ng marker file sa halip na hawak na lock, dahil mamanahin ng mga prosesong sinisimulan ng Claude Code ang lock at hahawakan ito kahit tapos na ang session. Hindi napapahinto ang pagsisimula ng marker na hindi maisulat: sinasabi ng ccam ang `ccam: could not mark the session as running (…); starting Claude Code anyway`.

   Binabasa rin ng iba pang pagsusuri ng ccam kung tumatakbo ang isang account (rename, remove, pagpapalit ng group, `doctor --fix`, ang MCP sync) ang sariling mga tala ng Claude Code sa folder ng account (`sessions/<pid>.json`, `daemon.lock`), kaya nakikita rin nila ang Claude Code na sinimulan nang walang ccam: karaniwang `claude`, isang IDE, mga background session at ang daemon nila.
8. Itinatakda ng ccam ang pamagat ng terminal sa `claude · <name>` (kapag terminal ang standard error), isinusulat ang oras ng pagsisimula sa `~/.local/state/ccam/last-used/<name>` nang walang `fsync`, at tinatawag ang `syscall.Exec` (`Plan.Exec`).

Tumatakbo ang Claude Code na nakatakda ang `CLAUDE_CONFIG_DIR` sa folder ng account at ang `CCAM_PROFILE` sa pangalan ng account. Binubura ng ccam ang anumang namanang value ng mga variable na ito, at ang namanang `CCAM_RELAY`. Para sa `~/.claude` folder na nirehistro gamit ang `ccam import` sa macOS, itinatakda rin nito sa walang lamang value ang `CLAUDE_SECURESTORAGE_CONFIG_DIR`, para iisang Keychain login ang gamitin ng account at ng karaniwang `claude`; sa Windows, ganoon din ang ginagawa nito kapag iniimbak ng Claude Code ang login sa Credential Manager.

### Bakit exec

Pinapalitan ng ccam ang sarili nitong proseso ng `claude` (`syscall.Exec`). Bilang resulta:

- hindi nananatili sa memory ang ccam;
- sa Claude Code mismo ang mga keystroke, signal at exit code;
- nananatiling array ang mga argument, hindi kailanman shell string.

Walang `exec` ang Windows. Doon, sinisimulan ng ccam ang `claude.exe` bilang child process na may parehong environment, mga argument at console, hinihintay ito at lumalabas gamit ang exit code nito (`Plan.Exec`). Ang `claude.cmd` na in-install ng npm ay pinapatakbo sa `cmd.exe` kasama ang buong path nito at naka-quote ang bawat argument para walang mabasa ang `cmd.exe` bilang sarili nitong syntax (`&`, `%PATH%`); hindi maipapasa sa ganoong paraan ang argument na may line break at tinatanggihan ito, kasama ang payong i-install ang native na Claude Code.

Direktang umaabot sa Claude Code mula sa console ang Ctrl+C, kaya binabalewala ito ng ccam. Dala ng session marker ang PID ng ccam, na nabubuhay nang eksaktong kasinghaba ng Claude Code, at inaalis ito kapag lumabas ang Claude Code. Habang naghihintay, walang ibang ginagawa ang ccam.

Ang mga gawain pagkatapos ng session, gaya ng pagpansin sa mga server na idinagdag gamit ang `claude mcp add`, ay nangyayari sa susunod na pagsisimula: ikinukumpara ng ccam ang modification time ng `.claude.json` sa `~/.local/state/ccam/mcp/<name>/seen` at pina-parse lang ang file kapag nagbago ito.

### Paano dinadala ang isang command line

Tatlong command line, tatlong ruta sa `internal/app`:

<!-- motion:launch-route -->

Itinutugma ang shortcut sa `argv[0]`, ang `ccam <name>` sa unang argument, at parehong nagtatapos sa `syscall.Exec`. Ang ccam command lang gaya ng `ccam usage` ang umaabot sa cobra.

### Mga shared na MCP server

Sa default na `flag` strategy, ipinapasa ng ccam sa Claude Code ang mga MCP server ng group bilang argument, na inilalagay bago ang sa iyo:

<!-- motion:launch-mcp -->

Idinadagdag lang ng ccam ang argument kung mayroong `mcp.json` ang group, at inilalagay ito kasunod mismo ng `claude` (`mcp.FlagArgs`, `Plan.InsertArgs`). Sa ganitong paraan, nalalapat pa rin sa ibabaw nito ang sarili mong `--mcp-config` o `--strict-mcp-config`. Isa itong argument, `--mcp-config=<path>`, dahil tumatanggap ng ilang value ang `--mcp-config` ng Claude Code: kapag isinulat bilang dalawang argument, babasahin din nito bilang config file ang kasunod na prompt.

Ang ccam ang nagpapasya kung saan hindi puwedeng ilagay ang flag (`claudeVerb`). Binabasa nito ang subcommand ng Claude Code mula sa mga argument sa paraang ginagawa ng Claude Code, nilalaktawan ang mga option kasama ang mga value nito. Walang flag na nakukuha ang `remote-control` o `rc` (sa ibaba), gayundin ang verb ng background session na unang ibinigay (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`). Wala ring nakukuha ang `ccam work mcp list`, `get` at `remove`, dahil tinatanggihan ito roon ng Claude Code; sa halip, sinasabi ng ccam sa standard error na naglo-load din ang mga server ng group sa mga session ng account:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Pinapanatili ng bawat ibang subcommand at bawat session ang flag.

Kapag may `mcp_strategy = "sync"` sa `config.toml`, isinusulat na lang ng ccam ang mga server sa `.claude.json` ng account, at nilalaktawan iyon kung tumatakbo na ang account. Nilalaktawan din ng pagsisimula ang gawaing iyon kapag hindi nagbago ang `mcp.json` ng group o ang `.claude.json` mula sa huling pagsusuri. Kung may sariling server ang account na kapareho ng pangalan ng isang server ng group, pinapanatili ng sync ang server ng account at sinasabi ito nang isang beses. Ibinabalik din ng sync ang kopyang nawala sa `.claude.json` (`put back group default's MCP servers …`), at binabago o inaalis lang nito ang mga kopyang isinulat nito na walang nagbago mula noon.

### Remote Control

Sinisimulan ng `ccam work rc` (o `ccam work remote-control`) ang Remote Control ng Claude Code bilang `work`. Tinatanggihan ng `claude remote-control` ang `--mcp-config` bago ang verb nito at wala itong tinatanggap pagkatapos nito, at binabasa lang ng mga session nito ang mga MCP server mula sa `.claude.json` ng account. Kaya iba ang takbo ng flag strategy dito:

<!-- motion:launch-rc -->

1. Nahahanap ng ccam ang `remote-control` o `rc` bilang verb ng Claude Code, kahit pagkatapos ng mga option at mga value nito (`ccam work --model opus rc`), at walang idinadagdag na flag (`claudeVerb`, `mcp.FlagArgs`).
2. Ikinukumpara ng ccam ang modification time ng `mcp.json` ng group sa `sync-stamp` ng account, at sinusuri kung nasa `.claude.json` pa ang bawat server ng group (`mcp.Synced`). Kung nasa `.claude.json` na ang kasalukuyang mga server ng group, o walang server ang group, nagsisimula ang Remote Control nang walang tanong.
3. Kung hindi, pinapangalanan ng ccam ang mga server na mawawala sa mga session at nagtatanong:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. Sa `y`, kinokopya ng ccam ang mga ito gaya ng ginagawa ng `sync` strategy (`mcp.Sync`): sa ilalim ng lock ng Claude Code, pagkatapos ng backup, at pinapanatili ang sariling server ng account. Sa `N`, o kapag walang terminal, nagsisimula ang Remote Control nang wala ang mga ito, at itinuturo ng ccam ang `mcp_strategy = "sync"`, na kumokopya sa mga ito bago ang bawat pagsisimula.

Nagtatanong ang ccam sa bawat pagsisimula ng Remote Control hanggang makopya ang mga server. Pagkatapos ng kopya, magtatanong lang ulit ito kapag nagbago ang `mcp.json` ng group o nawala sa `.claude.json` ang isang kinopyang server. Nananaig ang kopya sa `.claude.json` sa `--mcp-config` sa mga susunod na normal na pagsisimula, kaya ina-update ng `ccam mcp add`, `remove` at `import` ang mga kopyang hawak ng ibang account ([Mga shared na MCP server](/features/mcp)).

### Ang 10 ms na budget

10 ms sa 95th percentile ang budget para sa launch path; sa praktika, mga 4 ms ang idinadagdag ng ccam bago magsimula ang Claude Code. Hindi kailanman binubuo ng launch path ang command tree ng cobra, sinisimulan ang Bubble Tea o nagbubukas ng network connection, at pumapalya ang `TestPackageInitIsCheap` kapag lumampas sa kabuuang 5 ms ang mga function na `init()` ng package. Sinusukat ng `scripts/perf-check.sh` ang `ccam --dry-run work`, `ccam work` at isang shortcut gamit ang `hyperfine` at pumapalya kapag lampas sa 10 ms, at pinapatakbo nang paulit-ulit ng `BenchmarkLaunchPath` sa `internal/launcher` ang mga hakbang. Para makita ang oras ng bawat hakbang, itakda ang `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Kaligtasan at mga espesyal na kaso

- Humihinto ang hindi kilalang pangalan o numero na may `no account "<name>"` at itinuturo ang `ccam list`. Kapag wala pang kahit isang account, hihilingin ng ccam na magdagdag ka gamit ang `ccam add <name>`.
- Kung naging karaniwang file ang isang shared na path o tumuturo ito sa maling lugar, humihinto ang ccam na may `stopped because of share link problems:`, inililista ang mga ito, at hinihiling na patakbuhin mo ang `ccam doctor --fix`. Hindi ito naglilipat ng data sa launch path.
- Kung may ibang Claude Code session na nagshe-share ng kasaysayan ng usapan sa account na ito (parehong group) at tumatakbo sa parehong folder, nagbababala ang ccam na puwedeng maghalo ang kasaysayan at nagtatanong ng `Run anyway? [y/N]`. Itinutugma ang folder na naresolba ang mga symbolic link. Kapag walang terminal (itinuturing na walang terminal kapag galing sa `/dev/null` ang standard input), o mula sa isang job sa background, ipinapakita nito ang babala at nagpapatuloy. Nagsisimula ng sarili nitong usapan ang bagong `-p` run at walang nakukuhang babala; hindi talaga sinusuri ang `--version`, `--help` at ang mga subcommand na hindi nagbubukas ng usapan, gaya ng `mcp`, `auth` o `plugin`.
- Nilalaktawan ng `--dry-run` ang MCP sync, ang abiso tungkol sa bagong server, ang paalala sa `mcp list`, ang pagsusuri ng kasaysayan, ang session marker, ang pamagat at ang oras ng huling paggamit, at hindi ito nag-e-exec. Sinusuri pa rin nito ang mga share link, at ginagawa ang nawawalang link gaya ng sa karaniwang pagsisimula.
- Kapag pinatakbo ang `ccam work rc` habang tumatakbo na ang `work`, hindi nito ginagalaw ang `.claude.json`, sinasabi iyon, at sinisimulan ang Remote Control nang wala ang mga server ng group; patakbuhin ulit pagkatapos lumabas ng session para makopya ang mga ito. Sa `--dry-run`, ipi-print nito ang mensahe pero hindi ito magtatanong o kokopya.
- Sa `mcp_strategy = "sync"`, nagpi-print ang pumalyang sync ng `ccam: MCP sync failed (launching anyway): …` at nagsisimula ang Claude Code kasama ang mga server na mayroon na ito.

## Kaugnay

- [Mga shortcut command](/features/shortcuts)
- [Paggamit ng Remote Control sa ccam](/guides/remote-control)
- [Pamamahala ng mga account](/features/accounts)
- [Sanggunian ng mga command](/docs/commands)
- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Ipinaliwanag ang CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
