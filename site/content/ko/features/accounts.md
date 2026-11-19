---
title: 계정 관리
description: ccam add, ccam login, ccam list, ccam rename, ccam remove, ccam import로 Claude Code 계정을 추가, 로그인, 확인, 이름 변경, 삭제하고 가져옵니다.
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

ccam 계정은 각각 자체 로그인 정보를 가진 별도의 Claude Code 설정 폴더입니다. 다음 명령으로 계정을 추가하고, 로그인하고, 목록을 보고, 이름을 바꾸고, 삭제합니다. 이미 있는 폴더를 등록할 수도 있습니다.

## 사용법

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

명령이 `<name|number>`를 받는 곳에는 계정 이름이나 슬롯 번호(`ccam list`의 `#` 열)를 줄 수 있습니다.

`ccam use`는 기본 계정(`config.toml`의 `default`)만 바꿉니다.

터미널에서 옵션 없이 `ccam add`를 실행하면 설정 마법사가 열립니다. 마법사는 다음 단계를 거칩니다.

1. 계정 이름(직접 주지 않은 경우)
2. 로그인 방식: 지금 Claude 로그인, API 키, 나중에 Claude 로그인
3. 참여할 그룹. `default`가 선택된 상태로 보여 줍니다(아직 그룹이 없으면 `default` 그룹이 공유할 항목을 보여 줍니다)
4. `Create a shortcut command for work? [y/N]`
5. 확인할 요약

`Esc`를 누르면 한 단계 뒤로 갑니다. `--group`, `--no-group`, `--shortcut`, `--no-login` 중 하나라도 주거나, 터미널 없이(예: 파이프로) 실행하면 마법사를 건너뜁니다. 이때는 이름이 필요하며, `--shortcut`으로 이름을 지정하지 않으면 단축 명령을 만들지 않습니다.

`--group`을 주지 않으면 새 계정은 `default` 그룹에 들어갑니다. 첫 번째 계정이 이 그룹을 만들며, 그룹은 MCP 서버, 설정, `CLAUDE.md`, 확장 기능을 공유합니다. [그룹으로 공유하기](/features/share)를 참고하세요. Windows에서는 개발자 모드에서만 파일을 링크할 수 있습니다. 개발자 모드가 없으면 첫 번째 계정이 `mcp`와 `extensions`만 공유하는 `default`를 만들고 그 사실을 알리며(`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), 설정, 메모리, 대화 내역을 공유하는 그룹으로 `--group`을 지정하면 ccam이 거부합니다.

`ccam list`는 슬롯, 이름, 단축 명령, 이메일, 그룹, 마지막 사용 시각을 보여 주며, 기본 계정 앞에는 `*`가 붙습니다. 이메일은 폴더가 지금 로그인된 계정의 것(폴더의 `.claude.json`에서 읽음)이며, ccam이 기록해 둔 이메일과 다르면 기록된 이메일이 `(was …)`로 뒤에 붙습니다.

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

`--json`을 붙이면 각 계정이 `slot`, `shortcut`, `name`, `email`(표와 마찬가지로 현재 로그인의 이메일), `default`, `dir`, `last_used`(한 번도 시작하지 않은 계정에서는 생략), `groups`, `share`를 가진 객체로 출력됩니다. `share`는 각 항목(`mcp`, `settings`, `memory`, `extensions`, `history`)을 그 항목을 받는 그룹이나 `"local"`에 대응시킵니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `ccam add --no-login` | 로그인하지 않고 계정을 만듭니다. 나중에 `ccam login <name>`으로 로그인합니다 | 꺼짐 |
| `ccam add --group <group>` | `default` 대신 이 그룹에 들어갑니다. 여러 번 지정할 수 있으며, `default`가 없으면 만듭니다 | `default` |
| `ccam add --no-group` | 어느 그룹에도 들어가지 않습니다. 계정은 아무것도 공유하지 않습니다 | 꺼짐 |
| `ccam add --shortcut <command>` | `~/.local/bin`에 이 단축 명령도 만듭니다 | 없음 |
| `ccam add --api --key-from <source>` | 로그인 대신 API 키를 쓰는 계정. [API 키 계정](/features/api) 참고 | 꺼짐 |
| `ccam list --json` | 계정을 JSON으로 출력합니다 | 꺼짐 |
| `ccam remove -y`, `--yes` | 묻지 않고 삭제합니다 | 꺼짐 |
| `ccam logout -y`, `--yes` | 묻지 않고 로그아웃합니다 | 꺼짐 |
| `ccam import --name <name>` | 지정한 폴더의 계정 이름. 폴더 없이 주면 거부합니다(`--name needs a folder: ccam import <folder> --name ops`) | 폴더 이름에서 만듦 |
| `ccam import -y`, `--yes` | 찾은 폴더를 모두 묻지 않고 가져옵니다 | 꺼짐 |
| `--dry-run` | 아무것도 바꾸지 않는 `ccam list`와 `ccam completion`에서 지원합니다. `add`, `login`, `logout`, `use`, `rename`, `remove`, `import`는 `` --dry-run is not supported by `ccam add`; nothing was changed ``(해당 명령 이름이 들어감)를 출력하고 거부합니다 | 꺼짐 |

`ccam login`(`--expired` 외), `ccam use`, `ccam rename`에는 자체 옵션이 없습니다. `ccam completion`은 `bash`, `zsh`, `powershell`을 받으며, 그 밖의 값은 `only bash, zsh and powershell are supported`로 실패합니다.

계정 이름은 `a-z`, `0-9`, `-`, `_`로 이루어진 1~32자이며, 문자나 숫자로 시작해야 합니다. `ccam 2`는 슬롯 2를 뜻하므로 숫자로만 된 이름은 쓸 수 없습니다. ccam 명령 이름(`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`)도 쓸 수 없으므로, `ccam <name>`이 모호해지는 일은 없습니다. 이 목록은 `internal/config`의 `config.Reserved`입니다.

## 동작 방식

계정을 추가하면 폴더를 만들고, 그룹에 들어가고, 설정을 저장하고, 로그인하고, 단축 명령을 만드는 일을 이 순서대로 합니다. `ccam add work --shortcut ccwork`를 실행하면 ccam은 다음 단계를 거칩니다(`internal/cli/profiles.go`의 `addProfile`).

<!-- motion:accounts-add -->

1. ccam은 무엇이든 만들기 전에 이름(`config.ValidateName`)과 단축 명령(`shim.Check`)을 검사합니다.
2. ccam이 `0700` 권한으로 `~/.local/share/ccam/profiles/work/`를 만들고, 비어 있는 가장 낮은 슬롯 번호를 계정에 줍니다. 그래서 삭제한 계정의 번호가 다시 쓰입니다(`profile.Add`, `Config.NextSlot`). 첫 번째 계정은 기본 계정이 됩니다.
3. 계정이 자신의 그룹에 들어가고(`--group`이나 `--no-group`이 따로 정하지 않으면 `default`, 첫 번째 계정은 `default`를 만듭니다), ccam이 새 빈 폴더를 그 그룹에 링크합니다. 예를 들어 `settings.json`을 `groups/default/settings.json`에 링크합니다(`applyShare`, `share.Build` 다음 `Apply`). 이 단계가 실패하면 새 폴더를 삭제합니다. 폴더가 남아 있으면 이후의 모든 재시도가 `folder already exists`로 막히기 때문입니다.
4. 계정을 `fsutil.AtomicWrite`로 `~/.config/ccam/config.toml`에 저장합니다. 임시 파일에 쓰고 `fsync`한 뒤 `rename`하므로, 도중에 비정상 종료되어도 설정이 절반만 남지 않습니다.
5. ccam이 `CLAUDE_CONFIG_DIR`를 계정 폴더로 설정하고 `claude auth login`을 시작합니다. 로그인 페이지가 미리 채워지도록 계정의 마지막 이메일을 `--email`로 함께 넘깁니다(`runLogin`). 로그인이 실패해도 계정은 남으며, ccam은 오류를 출력하고 나중에 `ccam login work`를 실행하라고 안내합니다.
6. ccam이 `profiles/work/.claude.json`에서 `oauthAccount`를 읽습니다(`profile.ReadIdentity`). 이메일, 계정 ID와 조직 ID, 조직 이름, 요금제(`claude_max`, `claude_team` …)가 들어 있습니다. 이 파일은 수 MB가 될 수 있으므로, ccam은 파일 전체를 해석하지 않고 그 키 하나만 찾습니다(`TopLevelValue`). ccam은 이 값들을 `config.toml`의 `[profile.login]`에 보관합니다.
7. ccam이 `~/.local/bin/ccwork`를 ccam 실행 파일에 링크합니다. ccam은 시작할 때 `config.toml`에서 `argv[0]`을 찾습니다. 링크는 `~/.local/state/ccam/shortcuts.json`에 기록됩니다(`shim.Install`). 링크를 만들 수 없어도 계정은 추가되며, ccam이 이유를 알리고(`Could not create shortcuts: …`) 요약에 `shortcut: ccwork (not created; see above)`를 표시합니다.

로그인에는 몇 분이 걸릴 수 있고, 그동안 다른 ccam 프로세스가 설정을 저장할 수도 있습니다. 그래서 기다리는 명령은 `config.toml`을 다시 읽고, 바꿀 내용 하나(여기서는 이메일)만 적용해 저장합니다(`env.update`). 같은 이유로 `ccam remove`도 확인 질문 뒤에 설정을 다시 읽습니다.

`ccam login`은 기존 계정에 대해 5단계와 6단계를 실행합니다. Claude Code는 이메일을 로그인 페이지의 힌트로만 쓰므로, 이전과 다른 Claude 계정으로 로그인될 수 있습니다. ccam은 새 로그인을 `[profile.login]`과 비교해 그 사실을 알려 줍니다(`ccam: this login is other@company.com, but work was dev@company.com before.`). 새 로그인은 그대로 두고, 되돌리는 방법을 안내합니다: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

다른 Claude 계정으로 로그인하면 그 계정의 [한도 기록](/features/usage#the-limit-record)도 지웁니다. 그 한도는 다른 계정의 것이었기 때문입니다. 다른 ccam 계정이 같은 Claude 계정에 로그인되어 있으면, 두 계정이 사용량 한도를 함께 쓴다고 경고합니다.

`ccam login --expired`는 한 번도 로그인하지 않았거나 로그인이 만료된 모든 계정에 차례로 로그인합니다. 리프레시 토큰이 만료된 경우와, 서버가 거부한 뒤 Claude Code가 로그인을 지운 경우입니다. 첫 브라우저를 열기 전에 `config.toml`을 저장할 수 있는지 확인합니다.

`ccam logout`은 계정 안에서 `claude auth logout`을 실행하고 저장된 이메일을 지웁니다. 로그인 기록은 남으므로, 다음 로그인도 이 기록과 비교합니다. Claude Code는 로그인을 폐기하고 그 계정에 저장한 것을 MCP 서버 로그인까지 모두 지웁니다. 그래서 다음과 같은 경우 ccam은 먼저 경고하고 확인을 받습니다(`--yes`는 질문을 건너뜁니다).

- 로그인에 MCP 로그인이 들어 있을 때
- 계정이 실행 중일 때(그 세션도 로그아웃됩니다)
- 다른 계정의 폴더에 같은 로그인이 있을 때(복사한 폴더이며, 로그아웃하면 그 로그인도 폐기됩니다)

예를 들면 다음과 같습니다.

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

다른 Claude 계정으로 로그인하려면 `ccam login <name>`만으로 충분하며, MCP 로그인도 유지됩니다.

계정들이 둘 이상의 조직에 속하면 `ccam list`에 `Organization` 열이 추가되며(개인 요금제는 `personal`로 표시), `ccam list --json`에는 계정마다 `plan`과 `org`가 들어 있습니다.

### 원격 머신에서 로그인하기

SSH로 접속했을 때는 `claude auth login`의 브라우저 쪽 과정이 실패할 수 있습니다. Claude Code는 원격 머신의 콜백 서버(`http://localhost:<port>/callback`)에서 로그인을 기다리지만, 로그인을 마치는 브라우저는 내 컴퓨터에서 실행되므로 그 주소로는 아무 데도 연결되지 않고 `ERR_CONNECTION_REFUSED` 같은 오류가 표시됩니다. 그래도 주소에는 일회용 코드가 들어 있으므로 ccam이 대신 전달할 수 있습니다.

<!-- motion:accounts-relay -->

ccam은 터미널을 Claude Code에 그대로 넘기지 않습니다(`runLogin`). Claude Code의 표준 입력을 파이프에 연결하고, 입력한 줄을 하나씩 그 파이프로 복사합니다(`relayLogin`). 그래서 Claude Code의 `Paste code here` 프롬프트에 붙여 넣은 코드도 그대로 전달됩니다. 입력한 줄이 콜백 주소이면, 호스트가 `localhost`, `127.0.0.1`, `[::1]` 중 하나이고, 포트가 있고, 경로가 `/callback`이며, `code`와 `state`가 모두 있을 때만 ccam이 받아들입니다. `http://`는 생략해도 됩니다(`callbackURL`). 그다음 ccam이 브라우저가 했을 것처럼 서버에서 그 주소로 요청을 보냅니다. 타임아웃은 10초이고 리디렉션은 따라가지 않습니다(`deliverCallback`). 요청은 기다리고 있던 `claude auth login`의 콜백 서버에 도착합니다(그림의 점선). Claude Code가 코드를 받아 로그인 정보를 저장하며, Linux에서는 계정 폴더의 `.credentials.json`에 저장됩니다. 붙여 넣은 줄 때문에 ccam이 다른 호스트에 접속하는 일은 절대 없습니다.

ccam은 로그인을 시작하기 전에 다음을 출력합니다.

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

성공하면 `→ Passed the sign-in to Claude Code on this machine.`를 출력합니다.

로그인이 터미널에서 입력을 기다리는 동안 ccam은 터미널의 에코와 줄 단위 입력 모드를 끄고, 입력한 글자마다 `*`를 하나씩 보여 줍니다. 그래서 붙여 넣은 코드나 콜백 주소가 화면에 나타나지 않습니다. 줄 편집은 ccam이 직접 처리합니다(`readLine`). `Enter`는 줄을 보내고, `Backspace`는 한 글자를 지우며, 빈 줄에서 `Ctrl+D`를 누르면 입력이 끝납니다. 터미널 설정은 로그인이 끝날 때 복원되며, `Ctrl+C`, `SIGTERM`, `SIGHUP`을 받았을 때도 복원됩니다. 이 신호들은 여전히 ccam을 종료합니다. 표준 입력과 표준 오류가 터미널이 아니면 줄을 들어오는 그대로 복사합니다.

### 이름 변경, 삭제, 가져오기

`ccam rename`과 `ccam remove`는 계정 폴더와, ccam이 다른 곳에 계정 이름으로 보관하는 파일을 함께 다룹니다. ccam은 그 경로들을 `nameState`에 나열합니다.

- `~/.local/state/ccam/last-used/<name>`
- `~/.local/state/ccam/mcp/<name>`
- 사용량 캐시 `~/.local/state/ccam/usage/<name>.json`과 `<name>.statusline.json`
- 한도 기록 `<name>.limit.json`
- `~/.local/state/ccam/handoff/<name>/`의 세션 메모

`ccam remove work`는 이들을 다음 순서로 처리합니다.

<!-- motion:accounts-remove -->

1. `ccam remove`는 먼저 계정이 실행 중이면 거부하고, 그다음 확인을 받으며, 답을 받은 뒤 한 번 더 확인합니다(질문이 떠 있는 동안 시작한 세션도 잡아냅니다).
2. 무엇이든 바꾸기 전에 `config.toml`을 저장할 수 있는지 확인합니다.
3. 계정이 확장 기능을 공유하면 ccam이 그룹의 `installed_plugins.json`과 `known_marketplaces.json`을 다시 씁니다(`RewritePluginPaths`). Claude Code는 플러그인의 절대 경로를 저장하는데, 그 경로가 이 계정의 폴더를 거칠 수 있기 때문입니다.
4. 다음으로 관리하는 계정(`ccam import`로 등록한 계정이 아니라 ccam이 폴더를 만든 계정)을 로그아웃합니다. 폴더에서 `claude auth logout`을 실행하므로, 동작하는 로그인이 휴지통에 남지 않습니다(macOS에서는 폴더 이름을 딴 키체인 항목에도 남지 않으며, 남았다면 나중에 같은 이름으로 만든 계정이 물려받았을 것입니다). API 키 계정이거나, 다른 계정의 폴더에 같은 로그인이 있어 로그아웃하면 그 로그인도 폐기될 때는 이 단계를 건너뜁니다. 로그아웃이 실패하면 알리고 삭제는 계속합니다.
5. ccam이 폴더를 `~/.local/share/ccam/trash/<name>-<time>/`으로 옮기고(`profile.Trash`), 계정을 `config.toml`에서 빼고, 나중에 같은 이름으로 만든 계정이 물려받지 않도록 `nameState` 경로를 삭제한 뒤, 단축 명령 링크를 제거합니다: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

그래도 `config.toml`을 저장하지 못하면 폴더를 휴지통에서 되돌립니다. 삭제한 계정을 되살리려면 휴지통에서 폴더를 다시 옮기고 다시 로그인하세요.

`ccam rename`은 폴더 이름을 바꾸고, 폴더의 `plugins/`에 있는 플러그인 경로를 이전 폴더에서 새 폴더로 다시 씁니다. `nameState` 경로를 새 이름으로 옮기고, 기본 계정이었다면 기본 계정으로 유지합니다. 세션 메모도 고쳐 쓰므로 `ccam next --continue`가 계속 찾을 수 있습니다. 단축 명령 이름은 그대로 유지됩니다. 플러그인 경로를 다시 쓰다 실패하거나 `config.toml`을 저장하지 못하면 폴더 이름을 원래대로 되돌립니다.

macOS에서, 그리고 Claude Code가 로그인을 자격 증명 관리자에 보관하는 Windows에서는 폴더 이름 자체는 그대로 둡니다. 로그인과 MCP 로그인이 폴더 경로로 보관되므로, 폴더를 옮기면 로그아웃되기 때문입니다. 이때 설정에는 새 이름 옆에 폴더가 `folder`로 기록되고, ccam은 `its folder stays …: the macOS Keychain keeps the login under the folder's path`라고 알립니다.

`ccam import`는 폴더를 원래 위치에 그대로 두고 `config.toml`의 `dir`로 등록합니다. 폴더를 옮기거나 복사하지 않습니다. 리프레시 토큰을 복사하면 사본끼리 서로를 로그아웃시키기 때문입니다.

폴더를 지정하지 않으면 등록되지 않은 `~/.claude`와 `~/.claude-*` 폴더를 찾아 폴더마다 하나씩 묻습니다. 이름은 폴더에서 정해집니다(`~/.claude`는 `main`, `~/.claude-work`는 `work`, `~/.claude-2`는 `acct2`). 이름을 정할 수 없거나 이미 쓰이는 이름이면, 실행할 `--name` 명령과 함께 그 폴더를 나열합니다. 찾은 폴더 중 실패한 폴더는 `Skipped <path>: <reason>`으로 알리고 나머지는 가져옵니다.

ccam은 먼저 질문을 모두 마친 뒤 `config.toml`을 다시 읽고 폴더를 등록하므로, 그동안 다른 ccam 명령이 저장한 변경 사항도 남습니다. `✓` 줄은 저장한 뒤에 나옵니다.

`CLAUDE_CONFIG_DIR`가 없으면 Claude Code는 `.claude.json`을 `~/.claude` 옆, 즉 홈 폴더에 둡니다. 그래서 `~/.claude`를 가져올 때 그 안에 자체 `.claude.json`이 없으면, ccam은 `~/.claude.json`을 가리키는 `~/.claude/.claude.json` 링크를 만듭니다. 그러면 일반 `claude`와 `ccam main`이 같은 계정을 봅니다.

`~/.claude`에 `~/.claude.json`과 별도로 자체 `.claude.json`이 있으면, ccam은 둘 다 그대로 두고 그 사실을 알립니다. `ccam main`은 안쪽 파일을, 일반 `claude`는 바깥쪽 파일을 써서 로그인 정보와 MCP 서버가 따로 놀게 됩니다. 하나만 쓰려면 안쪽 파일을 링크로 바꾸세요([doctor](/features/doctor)도 이를 알려 줍니다).

폴더에 이미 등록된 계정과 같은 로그인이 있으면(복사한 폴더) ccam이 그 사실을 알립니다. 먼저 갱신하는 쪽이 다른 쪽을 로그아웃시키므로, 쓰기 전에 `ccam login <name>`을 실행하세요.

가져온 계정은 직접 그룹에 넣기 전까지(`ccam group join <group> <name>`) 어느 그룹에도 속하지 않으며, 단축 명령도 만들지 않습니다.

`ccam completion bash|zsh`는 계정 이름까지 완성해 주는 자동 완성 스크립트를 출력합니다.

## 안전장치와 예외 상황

- 세션이 폴더에 계속 쓰므로, 계정이 실행 중이면 `ccam rename`과 `ccam remove`는 ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``를 출력하고 거부합니다(`refuseRunning`). 이 검사는 Claude Code를 어떻게 시작했든 적용됩니다. ccam, 일반 `claude`, IDE, 백그라운드 세션과 그 데몬 모두입니다.
- 다른 ccam 프로세스가 그 로그인을 갱신하거나 바꾸는 중이면 `ccam rename`, `ccam remove`, `ccam logout`은 ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``를 출력하고 거부합니다.
- 다른 ccam 프로세스에서 그룹 변경이 진행 중이면 `ccam rename`과 `ccam remove`는 `another ccam process is changing the groups right now; try again when it is done`로, 그룹 변경이 중간에 멈춘 뒤라면 ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``로 거부합니다.
- 가져온 계정에 `ccam remove`를 실행하면 일반 `claude`를 위해 폴더와 로그인을 그대로 둡니다. 계정이 그룹에 속해 있으면 먼저 `ccam group leave`처럼 그룹에서 빠집니다. 공유하던 항목은 모두 폴더 안의 그룹 사본이 되고, 그룹의 MCP 서버는 계정의 `.claude.json`에 들어갑니다. 그다음 등록을 해제합니다: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. ccam이 만든 `.claude.json` 백업은 `<folder>/.ccam-backup/mcp-<time>`으로 옮기며, 옮길 수 없으면 `~/.local/state/ccam/mcp/<name>/backups`에 남기고 메시지로 알립니다. `~/.claude`는 휴지통으로 옮겨지지 않습니다.
- 두 계정이 같은 폴더를 쓸 수는 없습니다. 한 폴더를 두 번 적은 `config.toml`은 불러오지 못합니다.
- `ccam add`는 `profiles/`에 폴더가 이미 있는 이름을 거부하고 `ccam import`를 안내합니다. macOS에서는 이전 버전의 ccam이 남긴, 새 폴더 경로의 로그인이나 MCP 로그인이 키체인에 아직 있을 때도 멈추고, 이를 지우는 `security delete-generic-password` 명령을 알려 줍니다. `config.toml`을 저장하지 못하면 새 폴더를 다시 지웁니다.
- `ccam import`는 이미 등록된 폴더, 폴더가 아닌 경로, 이미 쓰이는 이름을 거부합니다. 폴더가 `~/.claude`의 사본(같은 로그인)이면 `~/.claude` 자체가 등록되어 있지 않아도 경고합니다.
- `ccam logout`은 계정이 실행 중이면 먼저 묻고, `ccam login`은 확인하지 않습니다.
- 로그인이 끝난 뒤처럼 전달 과정이 Claude Code에 닿지 못하면 `ccam: could not pass the sign-in to Claude Code:`와 이유를 출력합니다. HTTP 오류는 대개 코드가 만료되었다는 뜻이니 로그인을 다시 시작하세요.

## 관련 문서

- [Claude Code 시작하기](/features/launch)
- [단축 명령](/features/shortcuts)
- [MCP 서버, 설정, 대화 내역 공유](/docs/sharing)
- [컴퓨터 한 대에서 Claude Code 여러 계정 사용하기](/guides/multiple-claude-code-accounts)
