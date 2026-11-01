---
title: Claude Code 시작하기
description: ccam <name|number>, ccam run, 대시보드로 원하는 계정의 Claude Code를 시작합니다. 인자는 그대로 claude에 전달되며, --dry-run은 실행할 명령을 보여 줍니다.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: "11df08608a1e"
---

`ccam work`는 `work` 계정으로 Claude Code를 시작하며, 이름 뒤의 모든 것은 그대로 `claude`에 전달됩니다. 터미널, 스크립트, 파이프로 입력을 넘길 때 등 `claude`를 입력하던 곳이면 어디서나 쓸 수 있습니다.

## 사용법

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

계정 이름이나 `ccam list`의 번호를 지정하세요. ccam은 이름 뒤의 인자를 해석하지 않으므로, `ccam work --help`는 Claude Code의 도움말을 보여 줍니다. 이름 뒤의 `--`는 넣어도 되고 생략해도 됩니다. 표준 입력, 종료 코드, 시그널은 `claude`를 직접 실행할 때와 똑같이 동작합니다.

첫 번째 인자가 `-`로 시작하고 ccam 자체 옵션이 아니면, ccam이 계정을 대신 고르고 모든 인자를 Claude Code에 넘깁니다. 이때는 [`ccam folder`](/features/folder)로 현재 폴더에 지정한 계정을, 없으면 기본 계정을 씁니다. 기본 계정은 `ccam use`로 지정한 계정이며, 지정하지 않았다면 번호가 가장 낮은 계정입니다.

### 대시보드

표준 입력과 출력이 터미널이면 인자 없이 `ccam`을 실행할 때 [대시보드](/features/dashboard)가 열리고, 터미널이 아니면 `ccam list`를 출력합니다. 계정으로 이동한 뒤(화살표 키로, 또는 번호를 입력해서) `Enter`를 누르면 그 계정이 시작됩니다. `n`은 여유가 가장 많은 계정을 시작하고, `o`는 `--continue`, `--resume`, 그 밖의 인자를 고를 수 있게 합니다. 대시보드는 Claude Code가 시작되기 전에 닫히므로, 터미널이 먼저 원래대로 돌아옵니다.

**Team·Enterprise 시트.** 자동 계정 선택(`ccam next`, `--fallback`, `ccam pick`, 대시보드의 `n`, 그리고 켜 두었을 때의 [`auto_switch`](/features/next#auto-switch))은 같은 Claude Team이나 Enterprise 조직의 시트 사이에서도 계정을 바꾸며, 이는 Anthropic의 약관을 위반할 수 있습니다. 의존하기 전에 [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)를 읽고 소속 조직과 Anthropic에 확인하세요.

셸에 `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN` 중 하나가 설정되어 있으면 Claude Code가 계정 자체의 로그인 대신 그 키나 토큰을 쓸 수 있으므로, ccam은 시작하기 전에 표준 오류로 경고합니다(`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). ccam은 변수를 읽기만 하며, [API 키 계정](/features/api)은 대신 이 변수들을 지웁니다.

기본적으로 ccam은 자신을 Claude Code로 대체하므로, 계속 실행되는 것은 아무것도 없습니다. `config.toml`의 맨 위, 첫 `[[group]]`, `[[profile]]`, `[[folder]]`보다 위에 `relay = true`(또는 Linux와 macOS에서 [`auto_switch = true`](/features/next#auto-switch))를 설정하면, ccam은 대신 대화형 세션을 자식 프로세스로 시작하고 세션과 함께 계속 실행됩니다. 그래서 `/ccam:switch`나 사용량 한도에서의 `auto_switch`가 세션의 옵션, 권한 모드, 추론 수준, 에이전트를 유지한 채 다른 계정에서 스스로 대화를 이어 갈 수 있으며, `/ccam:switch`가 대화를 다른 조직의 계정으로 옮기기 전에는 ccam이 먼저 묻습니다([자세히](/features/next#relay)).

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--dry-run` | ccam이 실행할 환경 변수와 명령을 출력하고 종료합니다 | 꺼짐 |
| `--fallback` | `-p`/`--print` 실행용: 계정이 사용량 한도에 도달하면 로그인된 계정 중 여유가 가장 많은 계정으로 같은 명령을 다시 실행합니다([자세히](/features/next#run-again-at-the-limit)) | 꺼짐 |
| `-v`, `--version` | `ccam <version>`을 출력합니다(첫 번째 인자일 때만) | |
| `-h`, `--help` | ccam 도움말. `ccam --help` 또는 `ccam run --help`로 씁니다 | |

인자 규칙은 다음과 같습니다.

- ccam은 계정 이름 앞의 자체 옵션만 읽습니다. 이름 뒤의 모든 인자는 `--dry-run`과 `--help`를 포함해 Claude Code의 것입니다.
- `--dry-run`은 `ccam --dry-run work`, `ccam run --dry-run work`, `ccam --dry-run`(폴더의 계정이나 기본 계정)으로 쓸 수 있습니다. `--dry-run`과 `--fallback`은 순서를 바꿔 써도 됩니다.
- `ccam run <name|number> [claude args]`는 `ccam <name>`의 긴 형식입니다. 이름 없이 `ccam run`을 실행하면 폴더에 지정된 계정이나 기본 계정을 씁니다.
- 계정 이름에는 `add`, `list`, `run`, `usage` 같은 ccam 명령 이름을 쓸 수 없으므로, `ccam <word>`가 모호해지는 일은 없습니다.

시험 실행은 셸에, Windows에서는 PowerShell에 그대로 붙여 넣을 수 있는 한 줄을 출력합니다. 이 줄은 실행이 환경에서 바꾸는 것으로 시작합니다. `CLAUDE_CONFIG_DIR`와 `CCAM_PROFILE`은 항상, 지우는 변수(물려받은 `CCAM_RELAY` 등)는 `env -u <name>`으로, [API 키 계정](/features/api)이라면 엔드포인트, 모델, 그리고 키 자체가 아닌 키를 가져오는 곳(`ANTHROPIC_API_KEY="$WORK_KEY"`)을 보여 줍니다.

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

PowerShell 줄은 셸의 `VAR=x cmd`처럼 실행 뒤에 변수를 원래대로 돌려놓고, `Start-Process`로 Claude Code에 명령줄을 넘기므로 Windows PowerShell 5.1과 PowerShell 7에서 똑같이 실행됩니다. 실행이 변수를 빈 값으로 설정하면, 그 아래 `#` 줄이 PowerShell 7.5 이전에서는 그런 변수가 대신 지워진다고 알려 줍니다. npm이 설치한 `claude.cmd`라면, 이 줄은 실행이 쓸 `cmd.exe` 명령을 실행하며(아래 참고), 실행이 거기서 거부할 인자는 시험 실행에서도 실패합니다.

## 동작 방식

실행할 때마다 명령줄부터 `exec`까지 같은 단계를 거칩니다.

<!-- motion:launch -->

각 상자는 해당 단계가 실행될 때 강조되며, 아래 줄은 Claude Code가 시작될 때 받는 환경 변수와 인자입니다.

1. ccam은 `cmd/ccam/main.go`가 호출하는 `app.Main`에서 시작합니다. ccam이 다른 이름으로 시작되었다면 [단축 명령](/features/shortcuts)이므로 바로 실행기로 넘어갑니다.
2. ccam은 cobra가 관여하기 전에 `os.Args`를 살펴봅니다(`app.parse`). `ccam`만 입력한 경우(대시보드), ccam 명령 이름(`config.Reserved`), `-h`/`--help`, `-v`/`--version`만 `cli.Execute`로 갑니다. 나머지는 모두 실행입니다.
3. ccam이 `~/.config/ccam/config.toml`을 읽습니다. 해석하는 파일은 이것 하나뿐입니다(`launcher.Launch`).
4. ccam이 이름이나 번호로 계정을 찾고, 맨 앞의 `--` 하나를 제거한 뒤, `claude` 실행 파일을 찾습니다(`launcher.Prepare`, `ResolveClaude`). 이름이 없으면 작업 폴더를 `config.toml`의 `[[folder]]` 항목과 문자열로 비교해 폴더의 계정을 고르며, macOS와 Windows에서는 대소문자를 구분하지 않습니다(`HereProfile`). PATH 검색 결과는 파일의 수정 시각, `PATH`의 해시와 함께 `~/.local/state/ccam/claude-path`에 캐시되므로, ccam은 실행할 때마다 `PATH`를 검색하지 않습니다. 파일이 바뀌었거나 없어졌거나 `PATH`가 바뀌었으면 다시 검색합니다. `config.toml`의 `claude_bin` 경로(`~/`는 풀어 씁니다)는 실행 가능한 파일이어야 하며, 그렇지 않으면 무엇이든 쓰기 전에 `claude_bin X is not an executable file. Check claude_bin in config.toml`로 실행을 멈춥니다. `claude`가 ccam 자신이면 무한 반복이 되므로, inode로 비교해 거부합니다.
5. ccam이 계정 폴더 `~/.local/share/ccam/profiles/<name>/`를 공유 설정과 비교합니다. `lstat`과 `readlink`만 호출하므로 항목당 시스템 호출은 한두 번입니다(`checkLinks`, `share.Check`). 없는 링크는 그 자리에서 만듭니다. 데이터를 잃을 수 있는 문제가 있으면 실행을 멈춥니다.
6. ccam이 공유 MCP 서버를 추가합니다(`applyMCP`, 아래 참고).
7. ccam이 세션을 기록합니다(`session.Register`). PID를 이름으로 하고 프로세스 시작 시각을 담은 표시 파일을 `XDG_RUNTIME_DIR` 값과 상관없이 `~/.local/state/ccam/locks/profile-<name>/<pid>`에 기록하므로, 모든 셸이 같은 표시 파일을 봅니다. 대화 내역을 공유하는 계정이라면 ccam이 먼저 `history-<group>-<hash>`에서 실행 중인 세션을 찾고(해시는 심볼릭 링크를 푼 현재 폴더의 해시), 그곳에도 등록합니다(`historyCheck`).

   `exec`는 PID를 유지하므로, 이후 표시 파일은 Claude Code 프로세스 자체를 가리킵니다. ccam이 잠금을 잡아 두지 않고 표시 파일을 쓰는 이유는, Claude Code가 시작하는 프로세스가 잠금을 물려받아 세션이 끝난 뒤에도 쥐고 있게 되기 때문입니다. 표시 파일을 쓰지 못해도 실행은 멈추지 않으며, ccam은 `ccam: could not mark the session as running (…); starting Claude Code anyway`라고 알립니다.

   계정이 실행 중인지 확인하는 ccam의 다른 검사(이름 바꾸기, 삭제, 그룹 변경, `doctor --fix`, MCP 동기화)는 계정 폴더에 있는 Claude Code 자체의 기록(`sessions/<pid>.json`, `daemon.lock`)도 읽으므로, ccam 없이 시작한 Claude Code도 알아봅니다. 일반 `claude`, IDE, 백그라운드 세션과 그 데몬입니다.
8. ccam이 터미널 제목을 `claude · <name>`으로 설정하고(표준 오류가 터미널일 때), `fsync` 없이 실행 시각을 `~/.local/state/ccam/last-used/<name>`에 기록한 뒤 `syscall.Exec`를 호출합니다(`Plan.Exec`).

Claude Code는 `CLAUDE_CONFIG_DIR`가 계정 폴더로, `CCAM_PROFILE`이 계정 이름으로 설정된 상태에서 실행됩니다. ccam은 물려받은 이 변수들의 값과 물려받은 `CCAM_RELAY`를 지웁니다. macOS에서 `ccam import`로 등록한 `~/.claude` 폴더라면 `CLAUDE_SECURESTORAGE_CONFIG_DIR`도 빈 값으로 설정해, 그 계정과 일반 `claude`가 같은 키체인 로그인을 쓰게 합니다. Windows에서도 Claude Code가 로그인을 자격 증명 관리자에 보관하면 같은 일을 합니다.

### exec를 쓰는 이유

ccam은 자기 프로세스를 `claude`로 바꿉니다(`syscall.Exec`). 그 결과는 다음과 같습니다.

- ccam이 메모리에 남지 않습니다.
- 키 입력, 시그널, 종료 코드가 모두 Claude Code 자체의 것입니다.
- 인자는 셸 문자열이 아닌 배열로 유지됩니다.

Windows에는 `exec`가 없습니다. 그래서 Windows에서는 ccam이 같은 환경 변수, 인자, 콘솔로 `claude.exe`를 자식 프로세스로 시작하고, 끝날 때까지 기다린 뒤 그 종료 코드로 종료합니다(`Plan.Exec`). npm이 설치한 `claude.cmd`는 `cmd.exe`에서 전체 경로로 실행하며, `cmd.exe`가 어떤 인자도 자기 문법(`&`, `%PATH%`)으로 읽지 않도록 모든 인자를 따옴표로 감쌉니다. 줄바꿈이 들어간 인자는 이렇게 넘길 수 없으므로 거부하고, 네이티브 Claude Code를 설치하라고 안내합니다.

Ctrl+C는 콘솔에서 Claude Code로 바로 전달되므로 ccam은 이를 무시합니다. 세션 표시 파일에는 ccam의 PID가 들어가는데, ccam은 Claude Code와 정확히 같은 시간 동안 살아 있으며, 표시 파일은 Claude Code가 끝나면 지워집니다. 기다리는 동안 ccam은 다른 일을 하지 않습니다.

`claude mcp add`로 추가한 서버를 알아차리는 것처럼 세션 이후의 작업은 다음 실행 때 처리합니다. ccam은 `.claude.json`의 수정 시각을 `~/.local/state/ccam/mcp/<name>/seen`과 비교하고, 바뀌었을 때만 파일을 해석합니다.

### 명령줄이 처리되는 경로

세 가지 명령줄이 `internal/app`에서 세 가지 경로로 나뉩니다.

<!-- motion:launch-route -->

단축 명령은 `argv[0]`로, `ccam <name>`은 첫 번째 인자로 구분하며, 둘 다 `syscall.Exec`로 끝납니다. `ccam usage` 같은 ccam 명령만 cobra까지 갑니다.

### 공유 MCP 서버

기본값인 `flag` 방식에서는 ccam이 그룹의 MCP 서버를 인자로 Claude Code에 넘기며, 이 인자는 직접 넘긴 인자보다 앞에 놓입니다.

<!-- motion:launch-mcp -->

ccam은 그룹의 `mcp.json`이 있을 때만 인자를 추가하며, 그 인자를 `claude` 바로 뒤에 넣습니다(`mcp.FlagArgs`, `Plan.InsertArgs`). 그래서 직접 지정한 `--mcp-config`나 `--strict-mcp-config`도 그 위에 그대로 적용됩니다. Claude Code의 `--mcp-config`는 값을 여러 개 받으므로 `--mcp-config=<path>` 한 인자로 넘깁니다. 두 인자로 쓰면 뒤따르는 프롬프트까지 설정 파일로 읽게 됩니다.

ccam은 이 옵션을 넣을 수 없는 곳을 정합니다(`claudeVerb`). Claude Code와 같은 방식으로, 옵션과 그 값을 건너뛰며 인자에서 Claude Code의 하위 명령을 읽습니다. `remote-control`이나 `rc`에는 옵션을 넣지 않고(아래 참고), 맨 앞에 온 백그라운드 세션 동사(`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`)에도 넣지 않습니다. `ccam work mcp list`, `get`, `remove`도 Claude Code가 거기서 이 옵션을 거부하므로 넣지 않으며, 대신 ccam이 그룹의 서버도 그 계정의 세션에서 불러온다고 표준 오류로 알립니다.

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

다른 모든 하위 명령과 모든 세션에는 이 옵션이 그대로 들어갑니다.

`config.toml`에 `mcp_strategy = "sync"`를 설정하면 ccam은 대신 계정의 `.claude.json`에 서버를 기록하며, 그 계정이 이미 실행 중이면 이 단계를 건너뜁니다. 지난번 확인 이후 그룹의 `mcp.json`도 `.claude.json`도 바뀌지 않았으면 실행할 때 이 작업을 건너뜁니다. 계정에 그룹 서버와 이름이 같은 자체 서버가 있으면 동기화는 계정의 서버를 유지하고 한 번 알려 줍니다. 동기화는 또 `.claude.json`에서 사라진 복사본을 다시 넣으며(`put back group default's MCP servers …`), 자신이 쓴 복사본 중 그 뒤로 아무도 바꾸지 않은 것만 바꾸거나 지웁니다.

### Remote Control

`ccam work rc`(또는 `ccam work remote-control`)는 `work` 계정으로 Claude Code의 Remote Control을 시작합니다. `claude remote-control`은 동사 앞의 `--mcp-config`를 거부하고 동사 뒤에서도 받지 않으며, 그 세션은 계정의 `.claude.json`에서만 MCP 서버를 읽습니다. 그래서 여기서는 flag 방식이 다르게 동작합니다.

<!-- motion:launch-rc -->

1. ccam이 옵션과 그 값 뒤에 있더라도(`ccam work --model opus rc`) `remote-control`이나 `rc`를 Claude Code의 동사로 찾아내고, 옵션을 추가하지 않습니다(`claudeVerb`, `mcp.FlagArgs`).
2. ccam이 그룹 `mcp.json`의 수정 시각을 계정의 `sync-stamp`와 비교하고, 그룹의 모든 서버가 아직 `.claude.json`에 있는지 확인합니다(`mcp.Synced`). `.claude.json`에 이미 그룹의 현재 서버가 들어 있거나 그룹에 서버가 없으면, 묻지 않고 Remote Control을 시작합니다.
3. 그렇지 않으면 ccam은 세션에서 빠지게 될 서버 이름을 알려 주고 묻습니다.

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. `y`를 고르면 ccam이 `sync` 방식과 같은 방법으로 서버를 복사합니다(`mcp.Sync`). Claude Code의 잠금을 잡고, 백업을 만든 뒤, 계정 자체 서버는 그대로 둡니다. `N`을 고르거나 터미널이 아니면 이 서버 없이 Remote Control을 시작하고, 실행할 때마다 서버를 복사하는 `mcp_strategy = "sync"`를 안내합니다.

ccam은 서버가 복사될 때까지 Remote Control을 실행할 때마다 묻습니다. 복사한 뒤에는 그룹의 `mcp.json`이 바뀌거나 복사한 서버가 `.claude.json`에서 사라졌을 때만 다시 묻습니다. 이후 일반 실행에서는 `.claude.json`의 복사본이 `--mcp-config`보다 우선하므로, `ccam mcp add`, `remove`, `import`가 다른 계정이 가진 복사본도 갱신합니다([MCP 서버 공유](/features/mcp)).

### 10 ms 예산

실행 경로의 예산은 95번째 백분위수에서 10 ms이며, 실제로 ccam이 Claude Code 시작 전에 더하는 시간은 약 4 ms입니다. 실행 경로에서는 cobra 명령 트리를 만들지 않고, Bubble Tea를 시작하지 않으며, 네트워크 연결도 열지 않습니다. 패키지의 `init()` 함수가 모두 합쳐 5 ms보다 오래 걸리면 `TestPackageInitIsCheap`가 실패합니다. `scripts/perf-check.sh`는 `hyperfine`으로 `ccam --dry-run work`, `ccam work`, 단축 명령을 측정해 10 ms를 넘으면 실패하고, `internal/launcher`의 `BenchmarkLaunchPath`는 이 단계들을 반복 실행합니다. 단계별 시간을 보려면 `CCAM_TRACE=1`을 설정하세요.

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## 안전장치와 예외 상황

- 모르는 이름이나 번호를 주면 `no account "<name>"` 오류를 내고 멈추며, `ccam list`를 안내합니다. 계정이 하나도 없으면 `ccam add <name>`으로 추가하라고 안내합니다.
- 공유 경로가 일반 파일로 바뀌었거나 잘못된 곳을 가리키면 ccam은 `stopped because of share link problems:`와 함께 문제를 나열하고 멈추며, `ccam doctor --fix`를 실행하라고 안내합니다. 실행 경로에서는 데이터를 옮기지 않습니다.
- 이 계정과 대화 내역을 공유하는(같은 그룹의) 다른 Claude Code 세션이 같은 폴더에서 실행 중이면, ccam은 대화 내역이 섞일 수 있다고 경고하고 `Run anyway? [y/N]`를 묻습니다. 폴더는 심볼릭 링크를 푼 경로로 비교합니다. 터미널이 아니거나(`/dev/null`에서 오는 표준 입력도 터미널이 아닌 것으로 칩니다) 백그라운드 작업에서 실행하면 경고만 보여 주고 계속합니다. 새 `-p` 실행은 자기 대화를 따로 시작하므로 경고하지 않고, `--version`, `--help`, 그리고 `mcp`, `auth`, `plugin`처럼 대화를 열지 않는 하위 명령은 아예 확인하지 않습니다.
- `--dry-run`은 MCP 동기화, 새 서버 알림, `mcp list` 안내, 대화 내역 확인, 세션 표시 파일, 터미널 제목, 마지막 사용 시각을 건너뛰고 exec도 하지 않습니다. 다만 공유 링크는 여전히 확인하며, 없는 링크는 일반 실행처럼 만듭니다.
- `work`가 이미 실행 중일 때 `ccam work rc`를 실행하면 `.claude.json`을 건드리지 않고 그 사실을 알린 뒤, 그룹의 서버 없이 Remote Control을 시작합니다. 서버를 복사하려면 세션이 끝난 뒤 다시 실행하세요. `--dry-run`에서는 메시지만 출력하고 묻거나 복사하지 않습니다.
- `mcp_strategy = "sync"`에서 동기화가 실패하면 `ccam: MCP sync failed (launching anyway): …`를 출력하고, Claude Code는 기존 서버로 시작합니다.

## 관련 문서

- [단축 명령](/features/shortcuts)
- [ccam으로 Remote Control 사용하기](/guides/remote-control)
- [계정 관리](/features/accounts)
- [명령어 레퍼런스](/docs/commands)
- [MCP 서버, 설정, 대화 내역 공유](/docs/sharing)
- [CLAUDE_CONFIG_DIR 이해하기](/guides/claude-config-dir)
