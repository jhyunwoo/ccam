---
title: 다음 계정으로 전환
description: ccam next는 다른 계정으로 Claude Code를 엽니다. --auto는 여유가 가장 많은 계정을 고르고, --continue는 대화를 넘겨받으며, --fallback은 -p 작업을 재실행합니다.
updated: 2026-11-05
order: 8
translation: machine
sourceHash: "bc9837b7a6ae"
---

`ccam next`는 계정 이름을 지정하지 않아도 다른 계정으로 Claude Code를 엽니다. 계정이 사용량 한도에 도달했을 때 쓰면 됩니다. 전환하는 방법은 네 가지입니다.

- 그냥 `ccam next`를 실행하면 마지막으로 사용한 계정의 다음 계정을, `--auto`를 붙이면 여유가 가장 많은 계정을 고릅니다. `--continue`를 붙이면 두 계정이 대화 내역을 공유하든 안 하든 이 폴더의 마지막 대화를 새 계정에서 이어 갑니다.
- Claude Code 안에서는 [ccam 플러그인](#inside-claude-code)이 한도 때문에 세션이 멈추면 알려 주고, `/ccam:switch`로 대화를 넘길 수 있습니다.
- `claude -p`를 실행하는 스크립트에서는 `ccam --fallback`이 이 전환을 알아서 합니다([한도에서 다시 실행](#run-again-at-the-limit) 참고).
- [`auto_switch`](#auto-switch)를 켜 두면(Linux와 macOS) Claude Code가 한도에서 멈춘 대화형 세션도 ccam이 옮깁니다.

**Team·Enterprise 시트.** 자동 계정 선택(`ccam next`, `--fallback`, `ccam pick`, 대시보드의 `n`, 그리고 켜 두었을 때의 [`auto_switch`](/features/next#auto-switch))은 같은 Claude Team이나 Enterprise 조직의 시트 사이에서도 계정을 바꾸며, 이는 Anthropic의 약관을 위반할 수 있습니다. 의존하기 전에 [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)를 읽고 소속 조직과 Anthropic에 확인하세요.

## 사용법

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

ccam은 표준 오류에 전환 내용을 `→ <from> → <to>` 형식으로 출력하고(`--auto`로 골랐다면 그 이유도 함께), `ccam <name>`과 같은 방식으로 새 계정의 Claude Code를 시작합니다. `<from>`은 마지막으로 사용한 계정이며, `--continue`를 쓰면 대화가 속한 계정입니다. `--` 뒤의 인자는 그대로 `claude`에 전달됩니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--auto` | 로그인된 계정을 사용량으로 고릅니다(아래 참고). 사용량이 100%인 계정은 건너뜁니다 | 꺼짐 |
| `--max-usage <1-100>` | `--auto`와 함께: 사용량이 이 값 미만인 계정을 우선합니다 | `auto_max_usage`, 없으면 100 |
| `--strategy headroom\|expiring\|balance` | `--auto`와 함께: 5시간 사용량이 가장 낮은 계정, 5시간 구간이 가장 먼저 초기화되는 계정, 또는 하루당 남은 주간 할당량이 가장 많은 계정 | `auto_strategy`, 없으면 `headroom` |
| `--cross-org` | `--auto`와 함께: 다른 조직의 계정도 고릅니다(아래 참고) | `cross_org`, 없으면 꺼짐 |
| `--wait[=<duration>]` | `--auto`와 함께: 여유 있는 계정이 없으면 첫 계정이 초기화될 때까지 최대 이 시간만큼 기다립니다. 값은 `=` 뒤에 씁니다. `--wait 2h`는 `--wait takes its value after '=': --wait=2h`라는 힌트와 함께 거부합니다 | 꺼짐, 값 없이 주면 `5h` |
| `-c`, `--continue` | 이 폴더의 마지막 대화를 새 계정에서 이어 갑니다(아래 참고) | 꺼짐 |
| `--brief` | `--continue`와 함께: 대화 전체를 재개하는 대신, ccam이 이전 대화를 정리한 요약에서 새 계정으로 새 대화를 시작합니다([아래 참고](#brief)) | 꺼짐 |
| `-- <claude args>` | ccam 자체 인자 뒤에 `claude`로 넘길 인자 | 없음 |
| `--dry-run` | 지원합니다. Claude Code를 시작하거나 실행 기록을 남기지 않고, `→ <from> → <to>`와 실행할 명령줄을 출력합니다 | 꺼짐 |

`--auto`와 `--continue`는 함께 쓸 수 있습니다. `--continue`를 쓰면 `--resume` 인자가 직접 넘긴 인자보다 앞에 놓입니다.

`--dangerously-skip-permissions`나 `--model` 같은 Claude Code 자체 옵션은 `--` 뒤에 둡니다. 이름 뒤의 모든 것을 그대로 넘기는 `ccam <name>`과 달리 `ccam next`는 인자를 직접 읽으므로, `--`가 없으면 `ccam: unknown flag: --dangerously-skip-permissions` 오류를 내고 멈춥니다.

## 동작 방식

ccam으로 실행할 때마다 그 시각이 `~/.local/state/ccam/last-used/<name>`에 기록됩니다. `ccam next`는 이 파일들을 읽고, 어떤 계정이 로그인되어 있는지 확인한 뒤 계정을 고릅니다.

1. ccam이 `config.toml`의 계정마다 `next.Candidate`를 하나씩 만듭니다(`cli.candidates`). 여기에는 계정, 마지막 실행 시각(`launcher.LastUsed`), 동작하는 로그인이 있는지(`creds.Has`: `.credentials.json`, macOS에서는 키체인에 있는 Claude 로그인으로, 통째로 만료되지 않았고 Claude Code가 지우지 않은 것)가 들어갑니다. 대신 API 키를 쓰는 폴더(Console 로그인이나 `ccam add --api`로 추가한 계정)에는 요금제 한도가 없으므로 절대 고르지 않습니다.
2. `--auto`가 없으면 ccam이 계정을 슬롯 번호순으로 정렬하고 가장 최근에 실행된 계정을 찾습니다(`next.Rotate`). 그 계정이 지금 전환해 떠나는 계정입니다. 아직 실행된 계정이 없으면 가장 낮은 슬롯이 기준이 됩니다.
3. 그 지점부터 슬롯 순서대로 앞으로 진행하며, 마지막 슬롯 다음에는 첫 슬롯으로 돌아갑니다. 처음 만나는 로그인된 계정을 반환합니다.
4. ccam이 고른 계정을 `ccam <name>`과 똑같이 준비하고 시작합니다(`launcher.LaunchWith`). 계정의 링크를 확인하고, MCP 공유를 적용하고, 세션을 등록한 뒤 `syscall.Exec`로 ccam을 `claude`로 바꿉니다.

탐색은 마지막으로 사용한 계정에서 시작해, 그 뒤에 처음 나오는 로그인된 계정에서 멈춥니다.

<!-- motion:next-rotate -->

`work`가 가장 최근에 실행되었으므로 탐색은 그 다음부터 시작합니다. `side`는 로그인되어 있지 않아 건너뛰고, `lab`이 선택됩니다. `lab`도 로그아웃되어 있었다면 탐색은 처음으로 돌아가 `personal`로 갑니다(그림의 `% len`). 순환 방식은 출발한 계정을 절대 반환하지 않습니다. 로그인된 다른 계정이 없으면 오류를 내고 멈춥니다.

### --auto

`--auto`는 슬롯 순서가 아니라 사용량으로 결정합니다. ccam은 로그인된 계정만 남기고, `ccam usage`와 같은 방식으로 모든 계정의 사용량을 한꺼번에 요청한 뒤(`Fetcher.FetchAll`), 아래 규칙으로 고릅니다(`cli.pickAuto`, `usage.Policy.Choose()`).

<!-- motion:next-auto -->

계정은 5시간 구간과 주간 구간으로 판단하며, Claude Code에 넘기는 인자(`--` 뒤)에 `--model`이 있으면 그 모델의 주간 구간도 따집니다(이름에 `opus`, `sonnet`, `haiku`, `fable`이 들어간 경우). ccam은 계정을 세 단계로 나눕니다(`Policy.Choose()`).

1. **여유**: 모든 구간이 `auto_max_usage`(`config.toml`, 기본값 100) 미만인 계정. 이 중에서는 전략에 따라 고릅니다. `headroom`(기본값)은 5시간 사용량이 최저치에서 3%p 이내인 계정 중에서 주간 사용량이 가장 낮은 계정을 고릅니다. 그래서 주간 사용량이 더 낮으면 12%인 계정이 10%인 계정보다 앞서지만, 14%인 계정은 그렇지 않습니다. `expiring`은 5시간 구간이 가장 먼저 초기화되는 계정을 골라, 남은 할당량이 사라지기 전에 쓰게 합니다. `balance`는 주간 구간이 초기화될 때까지 하루당 남은 주간 할당량이 가장 많은 계정을 고르며, 실행 중인 ccam 세션도 그 계정의 몫으로 계산합니다. 그래서 주간 할당량이 쓰이지 못하고 남지 않으며, 병렬 세션이 여러 계정에 고루 나뉩니다.
2. **초과**: 100% 미만이지만 `auto_max_usage` 이상인 계정. 여유 있는 계정이 하나도 없을 때만 쓰며, 이 중 `headroom`이 고를 계정을 먼저 고릅니다.
3. **가득 참**: 100%인 구간이 있거나 5시간 사용량을 알 수 없는 계정. 절대 고르지 않습니다.

Claude Code가 어떤 구간에서 요청을 거부했으면, 저장된 수치와 상관없이 그 구간은 초기화될 때까지 100%로 칩니다. [ccam 플러그인](#inside-claude-code)의 훅이 이를 기록합니다(`~/.local/state/ccam/usage/<name>.limit.json`). 그래서 방금 한도에 도달한 계정은 수치가 오래되었더라도 초기화 전에 다시 고르지 않습니다.

모든 계정에 적용되는 설정은 `config.toml`의 맨 위, 첫 `[[group]]`, `[[profile]]`, `[[folder]]`보다 위에 둡니다. 이들 아래에 둔 설정은 그 테이블에 속하므로, ccam은 실행할 때마다 그 설정에 대해 경고하고, 설정을 저장하는 명령은 옮길 때까지 거부합니다.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**조직.** 자동으로 고를 때는 떠나는 계정과 같은 종류의 계정 안에서만 고릅니다. 같은 Team이나 Enterprise 조직의 시트끼리, 또는 개인 요금제(Pro, Max) 계정끼리 고르므로, 회사 계정의 작업이 개인 계정에서 이어지거나 그 반대가 되는 일이 없습니다. ccam은 각 계정의 조직을 로그인에서 알아냅니다(`ccam login`이 기록하며, 이전 로그인은 `.claude.json`에서 읽습니다). 조직을 알 수 없는 계정은 어느 계정과도 짝지을 수 있습니다. 같은 Claude 계정으로 로그인한 두 폴더는 한도를 함께 쓰므로, 한쪽이 다른 쪽을 이어받는 일은 없습니다. `--cross-org`나 `config.toml`의 `cross_org = true`를 쓰면 조직 규칙이 풀리지만, [`/ccam:switch`](#inside-claude-code) 뒤에 다른 조직으로 옮기는 것은 그래도 사용자가 yes라고 답할 때까지 기다립니다.

고를 수 있는 계정이 없으면 ccam은 종료 코드 75(`EX_TEMPFAIL`)로 끝나므로, 스크립트는 "나중에 다시 시도"와 실패를 구별할 수 있습니다. `--wait`를 주면 대신, 가득 찬 계정 중 막고 있는 구간이 가장 먼저 초기화될 때까지(엔드포인트에 새 구간이 나타나는 초기화 시각 30초 뒤) 기다리며 `ccam: no account has room; waiting until 15:10`을 출력하고, 다시 고릅니다. 초기화 시각을 알려 주는 계정이 없으면 10분 뒤에 다시 확인합니다. 기다리는 시간이 끝날 때까지 여유가 생기지 않으면 75로 종료합니다.

조회는 `ccam usage`의 규칙을 따릅니다. 가져온 지 60초가 안 된 수치는 다시 쓰고, 계정마다 조회 타임아웃은 3초이며, 응답하지 않거나 액세스 토큰이 만료된 계정은 저장된 수치로 판단하고, 한도 기록이 그 위에 덮입니다. 저장된 구간의 초기화 시각이 지났으면 0%로 계산하므로, 그 사이에 구간이 초기화된 계정이 밀려나지 않습니다. `--auto`는 토큰을 절대 갱신하지 않습니다. 일반 순환과 달리 `--auto`는 마지막으로 사용한 계정의 여유가 여전히 가장 많으면 그 계정을 고를 수도 있습니다.

### --continue

1. ccam은 모든 계정에서 이 폴더의 마지막 대화를 찾습니다. `projects/<folder>/`의 `<session>.jsonl` 파일 중 가장 최근의 프롬프트나 응답이 있는 파일입니다(폴더 이름은 Claude Code가 짓는 방식대로, 링크를 푼 작업 폴더 경로에서 문자와 숫자가 아닌 모든 글자를 `-`로 바꾼 것입니다). `/mcp`, `/usage`, `/status`만 실행한 세션이나 빈 파일처럼 프롬프트나 응답이 없는 세션은 건너뜁니다. 세션이 메모를 남겼고(사용량 한도나 `/ccam:switch`, 아래 참고) 그보다 2분 넘게 뒤에 프롬프트나 응답이 있는 대화가 없으면 그 세션을 그대로 씁니다. Claude Code가 `/exit`에서 쓰는 줄은 치지 않습니다.
2. 대화가 속한 계정이 전환해 떠나는 계정이 되고, 새 계정은 위와 같은 방식으로 같은 조직 안에서 고릅니다. 여러 계정이 그 대화 내역을 공유하면, `.claude.json`에 그 대화를 마지막으로 실행했다고 적힌(`lastSessionId`) 계정의 것으로, 없으면 마지막으로 사용한 계정의 것으로 봅니다.
3. 새 계정이 그 대화 내역을 공유하면(같은 그룹 폴더) ccam은 `claude --resume <session id>`로 바로 그 대화를 시작합니다(`--continue`로는 그 폴더에서 가장 최근의 대화가 열립니다).
4. 그렇지 않으면 `claude --resume <the conversation's file> --fork-session`으로 시작합니다. Claude Code가 대화 전체를 새 세션 ID로 새 계정 자체의 대화 내역에 복사하며, 원본은 그대로 남습니다. 그 뒤로는 새 계정의 `--continue`가 이 대화를 찾습니다.
5. `--` 뒤에 `--agent`를 주지 않으면, ccam은 대화가 쓰는 에이전트도 넘겨줍니다. 폴더의 `.claude/settings.local.json`이나 `.claude/settings.json`, 또는 계정의 `settings.json`에 있는 `agent` 설정이며, 없으면 대화의 마지막 `agent-setting` 줄이 가리키는 에이전트입니다. Claude Code는 대화를 복사할 때(4단계) 에이전트를 되살리지 않으므로, 그대로 두면 기본 에이전트로 모든 도구를 쓰며 대화를 이어 가게 됩니다. 새 계정이 이전 계정의 에이전트, 플러그인, 설정을 공유하지 않으면(그룹을 통해 같은 `agents` 폴더와 `plugins` 폴더, `settings.json`을 쓰지 않으면), ccam은 이동에 관한 줄 뒤에 다음 줄을 덧붙입니다: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

ccam은 대화의 크기를 알려 줍니다. 새 계정은 첫 응답에서 대화 전체를 다시 읽기 때문입니다. 프롬프트 캐시는 조직 사이에 절대 공유되지 않고 개인 요금제는 각각 별도의 조직이므로, 두 계정이 서로 다른 조직에 속하면 ccam이 `prompt caches are not shared between organizations, so the first reply reads it all again`을 덧붙입니다. 이전 계정은 로그인했지만 새 계정은 로그인하지 않은 MCP 서버가 있으면 ccam이 이를 알려 줍니다: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. 그룹에서 `history`를 공유하면 복사할 필요가 없어 대화가 하나로 유지되고, `/rewind`가 코드를 되돌릴 때 쓰는 체크포인트(`file-history/`)도 함께 공유됩니다.

### --brief

`--continue`는 대화 전체를 재개하므로, 새 계정은 첫 응답에서 그 대화를 전부 다시 읽습니다. `ccam next --continue --brief`는 대신 새 대화를 시작하며, ccam이 모델에 묻지 않고 이전 대화를 정리해 쓴 요약에서 출발합니다:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

요약에는 다음 내용이 이 순서대로, 적을 것이 있을 때만 들어갑니다.

1. 대화 제목
2. Claude가 세부 내용이 필요할 때 검색할 수 있는 대화 기록 파일의 경로
3. 요청이 넷보다 많았을 때의 첫 요청
4. Claude Code가 마지막으로 대화를 압축할 때 쓴 요약
5. 입력한 그대로의 `!` 명령과 슬래시 명령을 포함한 마지막 요청 네 개
6. Claude의 마지막 응답
7. 할 일 목록
8. Claude가 마지막으로 바꾼 파일 10개
9. 최대 800자까지의 폴더의 `git status --short --branch`

할 일 목록은 최대 10개 항목만 남깁니다. 끝난 항목부터 빼고, 그다음 가장 멀리 있는 단계를 빼므로 현재 단계와 바로 다음 단계들은 남습니다. 긴 부분은 가운데를 잘라 내며, 요약 전체는 16,000자를 넘지 않습니다.

ccam은 요약을 본인만 읽을 수 있는 `~/.local/state/ccam/brief/<session>.md`로 저장하며, 다음 `--brief`가 7일보다 오래된 요약을 지웁니다. Claude Code는 프롬프트 하나로 시작하고, 그 프롬프트가 가리키는 파일을 첨부합니다: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` `--` 뒤에 프롬프트를 주면 마지막 문장을 대신하고(`Then: <your prompt>`), `--` 뒤의 다른 옵션은 그 뒤에 붙습니다.

Claude Code는 멘션을 `"`에서 끝내고 `#`에서 자르므로, 요약 경로에 둘 중 하나가 있으면 ccam은 작업 폴더 기준 경로로 멘션합니다. 그 경로에도 있으면 프롬프트에 파일 이름을 적고, Claude가 Read 도구로 읽습니다. 이 파일은 작업 폴더 밖에 있으므로 Read 도구가 먼저 묻습니다(Manual 모드의 `-p` 실행에서는 Claude Code가 읽기를 거부하고, Claude는 요약 없이 진행합니다).

`ccam --dry-run next --continue --brief`는 요약을 출력하고 아무것도 저장하지 않습니다. 이전 대화는 그대로 남고, 새 대화는 요약에 적힌 것만 알므로, 지금까지 어떻게 왔는지의 모든 세부보다 작업 자체가 더 중요할 때 `--brief`를 쓰세요. `--continue` 없이 쓰면 ccam은 `--brief works with --continue`를 내고 멈춥니다. 대화를 읽을 수 없거나 요약할 내용을 찾지 못하면, 멈추고 `--brief` 없이 `ccam next --continue`를 쓰면 대화 전체를 재개한다고 알려 줍니다. [릴레이](#relay)와 [`auto_switch`](#auto-switch)는 항상 대화 전체를 재개합니다. Claude Code 자체의 `--brief` 옵션은 다른 옵션처럼 `--` 뒤에 둡니다.

## Claude Code 안에서

[ccam 플러그인](/docs/install#claude-code-plugin)은 훅 두 개를 추가합니다(Claude Code 2.1.78 이상).

- 턴이 사용량 한도에 걸려 멈추면 대화를 기록해 두고, 데스크톱 알림을 지원하는 터미널(iTerm2, WezTerm, Windows Terminal, kitty, Ghostty)에서 알림을 띄웁니다: ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  훅은 Claude Code가 사용량 한도로 치는 모든 메시지를 압니다. 5시간, 주간, 일반, 모델 한도(`You've hit your session limit`, `You've hit your limit`, 모델 할당량의 `You've reached your …`), 사용 크레딧과 지출 한도, 그리고 조직이 계정을 막은 경우입니다.

  Claude Code가 거부된 요청을 대화에 기록했으면, 훅은 가득 찬 구간도 초기화될 때까지 기록합니다([한도 기록](/features/usage#the-limit-record)). 단, 관리자가 언제든 풀 수 있는 조직의 차단은 기록하지 않습니다. 서버의 짧은 요청 제한이나 API 키의 속도 제한은 사용량 한도가 아니므로 건드리지 않으며, 아무도 지켜보지 않는 세션(`-p` 실행과 백그라운드 세션)도 마찬가지입니다. 하위 에이전트나 워크플로 에이전트가 사용량 한도에서 멈추면, 메인 대화가 그 에이전트를 기다리기만 하는 중이어도 그 세션이 멈춘 것으로 봅니다. Claude Code는 그 거부된 요청을 세션의 `subagents` 폴더에 있는 에이전트 자신의 파일에 기록하고, 훅은 거기서 읽습니다. 여러 에이전트가 한꺼번에 멈춰도 세션에는 메모와 알림이 하나씩만 생깁니다.

  Claude Code 2.1.234 이상은 세션 안에서 기다렸다가 5시간 구간이 초기화되면 스스로 이어 가기도 합니다. 알림은 그 대신 전환할 수 있는 선택지를 주며, [`auto_switch`](#auto-switch)를 켜 두면(Linux와 macOS) 다른 계정에 여유가 있을 때 ccam이 대신 전환합니다. 이어받을 수 있는 다른 계정이 없으면, 곧 같은 Team이나 Enterprise 조직에 속한 다른 계정이 없거나 다른 개인 요금제가 없으면(`cross_org`를 켜면 다른 어떤 Claude 계정이든 이어받을 수 있습니다), 알림은 대신 `You have no other personal plans, so ccam cannot move this conversation.` 또는 `You have no other accounts in Acme, so ccam cannot move this conversation.`로 끝납니다.
- `/ccam:switch [account]`는 대화를 기록해 두고, 요청을 보내기 전에 턴을 끝냅니다. `/exit`를 입력한 뒤 `ccam next --continue`를 실행하면, 지정한 계정(또는 다음 계정)에서 대화를 이어 갑니다. 아직 대화가 없는 세션에서는 이어 갈 것이 없다고 알려 줍니다. 이 명령은 사용자가 직접 입력할 때만 실행됩니다. 플러그인 0.3.2부터 Claude는 `/ccam:switch`를 실행할 수 없으며, Claude가 예약한 작업(CronCreate, ScheduleWakeup, `.claude/scheduled_tasks.json`)은 이 명령을 일반 텍스트로 Claude에게 넘깁니다.

  지정한 계정이 대화의 조직 밖에 있으면(다른 Team이나 Enterprise 조직의 시트, 시트에서 개인 요금제로 옮기거나 그 반대인 경우, 또는 어느 한쪽이 API 키 계정인 경우) `cross_org` 설정과 관계없이 `ccam next --continue`가 먼저 묻습니다: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. 계정을 지정하지 않았을 때 ccam이 직접 고른 계정(`cross_org = true`이면 다른 조직에서도 고를 수 있습니다)으로 옮기기 전에도 똑같이 묻습니다: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. `y`나 `yes`라고 답해야만 옮깁니다. Linux와 macOS에서는 질문 전에 입력한 키를 먼저 버립니다. 다른 답을 하거나 물어볼 터미널이 없으면 대화와 메모는 그대로 남습니다: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work``(ccam이 직접 고른 경우에는 `ccam next --continue`가 다시 고릅니다). 자동으로 고를 때와 마찬가지로, 조직을 알 수 없는 계정은 어느 계정과도 짝지을 수 있습니다.

메모는 세션마다 한 파일씩 `~/.local/state/ccam/handoff/<account>/<session>.json`에 남으므로, 한 계정의 두 세션이 서로의 메모를 덮어쓰지 않습니다. 메모는 12시간 동안 남아 있으며, 다음 세션이 시작되려 할 때 한 번만 쓰입니다. 실행에 실패하면 메모는 다음 시도를 위해 그대로 남습니다. 메모는 세션이 시작된 폴더에 속하며, Claude Code의 셸에서 `cd`로 옮겨 가도 마찬가지입니다. `ccam rename`은 메모 안의 이름을 고쳐 씁니다. 메모에 적힌 계정이 그 사이에 삭제되었으면 ccam은 그 사실을 알리고 평소처럼 고릅니다.

### 릴레이

`config.toml`의 맨 위에 `relay = true`를 설정하면 ccam은 자신을 Claude Code로 대체하지 않고 대화형 세션과 함께 계속 실행되며, `/ccam:switch` 뒤에 그 세션의 대화를 스스로 이어 갑니다(다른 세션이 남긴 메모는 처리하지 않습니다). 훅이 Claude Code를 끝내면 ccam이 그 세션의 메모로 `ccam next --continue`를 실행합니다. 계정을 지정하지 않으면 대화의 모델까지 따져 여유가 가장 많은 계정을 고르며, `cross_org = true` 때문에 그 계정이 다른 조직의 계정이 되면 대화를 그리로 옮기기 전에 위의 질문을 합니다. Windows에서는 훅이 Claude Code를 끝낼 수 없으니 `/exit`를 입력하세요. 그러면 ccam이 이어 갑니다. 대화의 조직 밖에 있는 계정일 때도 훅은 Claude Code를 끝내지 않고 `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`라고 알리며, `/exit`를 입력하면 ccam이 위의 질문을 합니다. 세션마다 기다리는 ccam 프로세스가 하나씩 생기므로 기본값은 꺼짐입니다.

이어 받은 세션은 대화를 재개해도 복원되지 않는 것을 유지합니다. `--add-dir`, `--mcp-config`, 도구 목록, `--settings`, `--plugin-dir`, 시스템 프롬프트 옵션 같은 옵션과 대화의 권한 모드, 세션의 추론 수준, 대화가 쓰는 에이전트([`--continue`](#continue) 참고)입니다. `--model`과 `--name`은 Claude Code가 대화에서 직접 복원합니다. 그 밖의 것(`--worktree`, `--debug`, `--remote-control`, ccam이 모르는 옵션)은 `ccam: not carried over: …`로 알려 줍니다. 전환할 때는 대화 내역 공유에 관한 질문을 절대 하지 않으며, 전환하는 동안 입력한 내용은 버려집니다. 전환이 하는 질문은 다른 조직의 계정으로 옮기기 전의 질문 하나뿐입니다.

릴레이는 Claude Code가 `/ccam:switch`가 요청한 방식으로 끝났을 때만 이어 갑니다. 먼저 `/exit`를 입력했거나, 끝나는 동안 Ctrl+C를 눌렀거나, ccam이 행업이나 다른 시그널을 받았거나, 터미널이 사라졌을 때는 이어 가지 않습니다. 이때 메모는 남아 있으므로, 나중에 `ccam next --continue`로 대화를 이어 갈 수 있습니다. 이전 버전의 ccam 릴레이로 시작한 세션이나 아무도 지켜보지 않는 세션(`-p` 실행과 백그라운드 세션)은 훅이 끝내지 않고, 대신 `/exit` 메시지를 보여 줍니다. `claude remote-control`은 릴레이 없이 실행됩니다.

### `auto_switch`

`config.toml` 맨 위에 `auto_switch = true`를 두면(Linux와 macOS) Claude Code가 사용량 한도에서 세션을 멈출 때 ccam이 스스로 세션을 다른 계정으로 옮깁니다. 플러그인의 훅이 여유가 있는 다른 계정을 골라 Claude Code를 끝내고, 세션과 함께 실행 중이던 ccam이 [릴레이](#relay)가 `/ccam:switch` 뒤에 하는 것처럼 그 계정에서 대화를 이어 갑니다. 기본값은 꺼짐이며, 실험적 기능입니다.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam은 다음 조건을 모두 만족할 때만 세션을 옮깁니다:

1. ccam이 터미널에서 시작했고 ccam이 함께 계속 실행 중인 세션입니다. `auto_switch`는 [릴레이](#relay)를 켜므로, 설정한 뒤에 시작한 대화형 세션은 이 조건을 만족합니다. 그 전에 시작했거나 이전 ccam으로 시작한 세션에는 ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.``로 끝나는 알림이 갑니다. 아무도 지켜보지 않는 세션(`-p` 실행과 백그라운드 세션), `claude remote-control`, ccam 없이 시작한 Claude Code는 절대 옮기지 않습니다.
2. 그 계정에 [ccam 플러그인](/docs/install#claude-code-plugin)이 켜져 있고, 훅에 자기 프로세스 ID(`CLAUDE_PID`)를 넘겨주는 Claude Code 2.1.214 이상을 씁니다.
3. Claude Code가 그 계정의 5시간, 주간, Opus 또는 Sonnet 한도에 대해 거부된 요청을 초기화 시각과 함께 대화에 기록했습니다. 지출 한도, 사용 크레딧 소진, 조직의 차단, Fable 할당량(`You've reached your Fable limit.`), 서버의 요청 제한, API 키의 속도 제한은 절대 세션을 옮기지 않습니다. 하위 에이전트나 워크플로 에이전트의 한도는 메인 대화가 그 에이전트를 기다리기만 하는 중이어도 여기에 포함되며, 그 한도가 주간 Opus 또는 Sonnet 한도이면 대화 자체가 그 모델로 실행될 때만 세션을 옮깁니다.
4. ccam이 계정의 요금제를 압니다. 계정의 `.claude.json`이나 `ccam login`이 기록한 내용에서 알아냅니다. 모르면 알림이 `ccam does not know personal's plan, so auto_switch leaves this session.`으로 끝납니다.
5. 세션이 Remote Control에 연결되어 있지 않습니다. 다른 계정으로 옮기면 세션을 조작하던 앱에 다시 연결되지 않으므로 세션은 그대로 남고, 알림이 그렇다고 알려 줍니다.
6. 한도가 지금부터 15분보다 뒤에 초기화됩니다. 그보다 짧게 기다리면 되는 경우 ccam은 세션을 Claude Code에 맡깁니다: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. 한도 이후 프롬프트를 입력하지 않았습니다.
8. 다른 계정에 여유가 있습니다. ccam은 [`--auto`](#auto)처럼(`auto_max_usage`, `auto_strategy`, 대화의 모델) 떠나는 계정과 같은 종류의 계정 가운데서 고릅니다. 같은 Team이나 Enterprise 조직의 다른 시트, 또는 다른 개인 요금제(Pro, Max)입니다. 여기에는 `cross_org`와 `--cross-org`가 적용되지 않습니다. 고를 계정은 다음 조건을 갖춰야 합니다.
   - Claude Code가 갱신할 수 있는 로그인으로 로그인되어 있음
   - 같은 Claude 계정으로 로그인되어 있지 않음
   - API 키를 쓰지 않음
   - ccam이 아는 요금제임
   - ccam 자체의 사용량 조회로 보이는 한에서 사용 크레딧이 꺼져 있음
   - Claude Code에서 이 폴더를 신뢰한 적이 있음
   - 대화가 에이전트로 실행되면(`agent` 설정이나 대화가 에이전트를 가리키면) 떠나는 계정의 에이전트를 공유함: `agents` 폴더와 `plugins` 폴더, `settings.json`이 각각 (그룹을 통해) 똑같은 것이거나 양쪽에 모두 없어야 함

   마지막 세 조건 때문에 여유 있는 계정이 모두 빠지면 알림이 그 계정을 밝힙니다: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` 또는 `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` [한도 기록](/features/usage#the-limit-record)이 있는 구간은 가득 찬 것으로 치며, 최근 2분 안에 다른 세션이 옮겨 간 계정은 다른 계정들보다 뒤로 밀립니다. 선택은 약 6초 안에 끝나야 하며, 사용량 조회가 그보다 오래 걸리면 세션은 그대로 남습니다. `usage_source = "claude"`이면 `auto_switch`는 저장된 수치만으로 판단하며, ccam은 Claude Code의 조회나 상태 표시줄로는 사용 크레딧이 켜져 있는지 알 수 없으므로 꺼져 있다고 가정합니다.

그러면 훅이 Claude Code를 끝내고, 함께 실행 중이던 ccam이 데스크톱 알림을 지원하는 터미널에서 `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` 같은 알림을 띄웁니다. 그리고 위 출력 예의 첫 줄을 출력하고, 고른 계정으로 그 세션에 대해 `ccam next --continue`를 실행합니다. 이 명령이 출력 예의 나머지 두 줄을 출력하고, [릴레이](#relay)처럼 세션의 옵션, 권한 모드, 추론 수준, 에이전트를 유지한 채 Claude Code를 시작합니다. 새 세션의 첫 메시지는 고정 문구이며, ccam이 사용자 대신 보내는 것은 이것뿐입니다(`/ccam:switch` 뒤의 이동은 아무것도 보내지 않습니다):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

하위 에이전트가 한도에 걸렸을 때는 하위 에이전트, 워크플로, 백그라운드 명령이 이전 세션과 함께 멈췄으므로 다음 문구를 보냅니다:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

그 계정도 한도에 도달하면 거기서도 같은 일이 일어납니다. 한도 기록 덕분에 이 사슬은 초기화 전의 계정으로 되돌아가지 않으므로, 모든 계정이 가득 차면 세션은 마지막 계정에 남습니다.

훅은 대화가 에이전트로 실행된다는 것을 `agent` 설정에서, 또는 Claude Code가 때때로, 그리고 세션이 끝날 때 대화에 기록하는, 에이전트를 가리키는 줄에서 알아냅니다. 지금까지 세션의 `--agent`만 에이전트를 지정했고 고른 계정이 에이전트를 공유하지 않으면, `ccam next --continue`는 아무것도 시작하지 않고 메모를 남겨 두며, 두 번째 데스크톱 알림이 그렇다고 알려 줍니다: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

여유 있는 다른 계정이 없으면 아무것도 옮기지 않습니다. Claude Code가 평소처럼 한도를 처리하며(세션 안에서 기다리거나 선택지를 보여 줌), 알림은 어느 계정이 가장 먼저 여유가 생기는지 알려 줍니다. 예: `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` 시트라면 `None of your other accounts in Acme has room now`라고 알립니다. 이미 시작된 이동을 멈추려면 Claude Code가 닫히는 동안 Ctrl+C를 누르세요: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

이동에는 비용이 있습니다. Claude Code가 다시 시작하는 데 몇 초가 걸리고, Claude가 시작한 백그라운드 셸, 개발 서버, 에이전트는 이전 세션과 함께 멈춥니다. 새 계정은 [`--continue`](#continue)처럼 첫 응답에서 대화 전체를 다시 읽으며, ccam이 그 크기를 알려 줍니다. 이동에는 [`--brief`](#brief)를 절대 쓰지 않습니다.

`auto_switch`는 Claude Code가 한도를 알리기 전에 전환하는 일이 절대 없고, 백그라운드에서 사용량을 지켜보지 않으며, 로그인을 복사하거나 갱신하거나 넘기지 않습니다. Windows에서는 무시되며 `ccam doctor`가 그렇다고 알려 줍니다. Windows에서도 `relay = true`이면 `/ccam:switch`와 `/exit`로 대화를 다른 계정으로 옮길 수 있습니다.

**Anthropic의 약관.** Team·Enterprise 시트에서 `auto_switch`는 사용자의 작업을 같은 조직의 다른 시트로 옮깁니다. Anthropic 고객센터는 "Each seat on a Team plan is assigned to one named individual, and only that person may use it"(팀 플랜의 각 시트는 한 명의 지정된 개인에게 할당되며, 그 사람만 사용할 수 있습니다), "Users can't route their work through another user's account to use that user's remaining usage"(사용자는 다른 사용자의 계정을 통해 작업을 라우팅하여 해당 사용자의 남은 사용량을 사용할 수 없습니다)라고 안내합니다([Team 요금제](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129), [Enterprise 요금제](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)도 마찬가지).

개인 요금제(Pro, Max)에 대해서는, 사용자 자신의 계정끼리 한도에서 자동으로 이어받아도 되는지 Anthropic의 약관이 정하고 있지 않습니다. [Consumer Terms](https://www.anthropic.com/legal/consumer-terms)는 Anthropic이 명시적으로 허용하지 않는 한 "through automated or non-human means, whether through a bot, script, or otherwise"(봇, 스크립트 등 자동화된 수단이나 사람이 아닌 수단으로) 서비스에 접근하는 것을 금지하며, 이동할 때 ccam은 위의 고정 메시지 중 하나를 사용자 대신 보냅니다. 2026년 11월 12일부터 시행되는 [Usage Policy](https://www.anthropic.com/legal/aup)는 금지 사항 가운데 하나로 "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products"(제품에 마련된 기능, 제한, 안전장치를 의도적으로 우회하거나 방해하는 것)를 듭니다.

ccam은 `auto_switch`가 이 약관들을 지킨다고 주장하지 않습니다. 켤지는 사용자가 판단할 일이며, 그 위험도 사용자가 집니다. Anthropic은 약관이 위반되었다고 판단하면 계정을, 또는 Team이나 Enterprise 조직의 접근을 정지할 수 있습니다.

`auto_switch`가 켜져 있는 동안 `ccam doctor`는 ccam 플러그인이 없는 계정과 2.1.214보다 오래된 Claude Code 버전에 대해 경고하고, 시트에 대해서는 Anthropic의 시트 규칙을 다시 알려 줍니다([점검 항목](/features/doctor)).

## 한도에서 다시 실행

`ccam next`는 멈춘 대화형 세션을 위한 명령입니다. 스크립트나 예약 작업에서 실행하는 `claude -p` 같은 비대화형 실행은 사용자를 기다릴 수 없으므로, 계정 이름 앞에 `--fallback`을 붙이면 ccam이 대신 전환합니다.

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

1. `--fallback`은 `-p`나 `--print`와 함께 쓸 때만 동작합니다. 그렇지 않으면 ccam은 멈추고 `ccam next --auto --continue`를 안내합니다.
2. 파이프로 들어온 표준 입력은 최대 10 MiB까지 한 번 읽어 두고 모든 시도에 똑같이 넘깁니다. Claude Code와 마찬가지로 입력이 시작되기를 최대 3초까지 기다리며(`--input-format stream-json`일 때는 제외), 파이프가 열려 있지만 아무것도 보내지 않으면 입력 없이 실행하고 ccam이 표준 오류에 이를 알립니다.
3. ccam은 `exec` 대신 `claude`를 자식 프로세스로 시작하고, 표준 출력은 붙잡아 두며, 표준 오류는 나오는 대로 그대로 전달합니다.
4. `claude`가 0으로 종료하거나 다른 이유로 실패하면, 붙잡아 둔 출력과 종료 코드를 바꾸지 않고 그대로 전달합니다. 시그널로 끝난 실행(128보다 큰 종료 코드)도 마찬가지입니다. 멈추라는 요청을 받은 것이므로 다시 실행하지 않습니다.
5. Claude Code의 메시지가 사용량 한도를 알릴 때(`--output-format json`이나 `stream-json` 출력 안에 있는 경우를 포함해 [훅](#inside-claude-code)이 아는 사용량 한도 메시지 중 하나), 또는 계정의 사용량에서 5시간, 주간, `--model` 구간이 100%로 나올 때(한도 기록 포함) 그 실패를 사용량 한도로 봅니다.
6. 그러면 그 시도의 출력은 버리고, 로그인된 다른 계정으로 같은 명령을 실행합니다. 5시간 사용량이 가장 낮은 계정부터, 그다음 사용량을 조회할 수 없는 계정을 번호 순서대로 시도합니다. 100%인 계정은 건너뛰고, API 키로 로그인한 계정은 절대 시도하지 않으며, 계정마다 최대 한 번만 시도합니다.
7. 인자에 `-c`, `--continue`, `-r`, `--resume`이 있으면 대화를 찾을 수 있도록 같은 대화 내역 그룹의 계정만 시도합니다.
8. `--auto`와 마찬가지로 같은 조직의 계정만 시도합니다(`cross_org = true`로 이 제한을 풀 수 있습니다).

다음 계정은 `--auto`와 같은 규칙(`auto_max_usage`, `auto_strategy`, `--model` 구간)으로 고릅니다. 남은 계정이 없으면 ccam은 마지막 시도의 출력을 표시하고, 여유가 남은 로그인 계정이 없다고 알린 뒤, 종료 코드 75로 종료합니다. 계정 이름 앞에 `--wait`를 주면(`ccam --fallback --wait=2h work -p …`) 대신 첫 번째 계정이 초기화될 때까지 기다린 뒤 모든 계정을 다시 시도할 수 있습니다. 항상 먼저 한 번은 기다리므로, 방금 한도에 도달한 계정을 곧바로 다시 실행하지는 않습니다. 이 시간은 기다림 하나하나가 아니라 실행 전체의 총 대기 시간을 제한합니다. 띄어 쓰면(`--wait 30m`) `ccam next`와 같은 힌트와 함께 거부합니다. 다른 실행과 마찬가지로 Claude Code 인자 앞에 `--`를 붙여도 됩니다(`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …`는 첫 번째 시도의 명령을 출력합니다.

[`auto_switch`](#auto-switch)를 켜지 않으면 대화형 세션은 자동으로 전환하지 않습니다. `ccam statusline install`을 해 두면, 5시간이나 주간 사용량이 90%에 도달했을 때 상태 표시줄 끝에 `→ ccam next --auto --continue`가 붙습니다([사용량 한도](/features/usage)).

## 안전장치와 예외 상황

- 로그인되지 않은 계정은 절대 선택되지 않습니다. 로그인된 다른 계정이 없으면 일반 `ccam next`는 `no logged-in account other than <name>` 오류를 내고 멈춥니다.
- `--continue`를 썼는데 이 폴더에 대화가 없으면 `no account has a conversation in <folder> to continue`가 출력됩니다.
- 조직 규칙 때문에 고를 계정이 남지 않으면 `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``가 출력됩니다.
- `--auto`를 썼는데 로그인된 계정이 없으면 `no logged-in account to pick`가 출력됩니다. 모든 계정이 가득 찼거나 사용량을 알 수 없으면 `` no account with known usage and headroom (see `ccam usage`) ``가 출력됩니다. 둘 다 종료 코드 75로 끝납니다.
- `--max-usage`, `--strategy`, `--wait`, `--cross-org`는 `--auto`와 함께만 쓸 수 있으며, 1-100 밖의 값이나 알 수 없는 전략은 거부됩니다.
- 포크는 대화만 복사할 뿐 계정의 다른 상태는 복사하지 않습니다. 대화 내역을 공유하지 않으면 원래 계정의 `file-history/`에 있는 체크포인트는 따라오지 않습니다.
- `--auto`는 `--dry-run`을 붙여도 네트워크 요청을 보냅니다. 읽기만 하며 토큰 파일에는 절대 쓰지 않습니다.
- 선택된 계정은 `ccam <name>`과 같은 경로로 시작하므로 같은 검사가 적용됩니다. 같은 대화 내역을 공유하는 다른 Claude Code 세션이 이 폴더에서 이미 실행 중일 때 나오는 경고도 마찬가지입니다.
- `last-used`는 Claude Code가 실제로 시작될 때만 기록됩니다. `--dry-run`은 이 파일을 건드리지 않으므로, 시험 실행으로 순환 순서가 넘어가지 않습니다.
- `--fallback`은 시도가 끝날 때까지 표준 출력을 붙잡아 두므로, `--output-format stream-json` 같은 스트리밍 출력도 마지막에 한꺼번에 나옵니다. `stream-json` 입력을 스트리밍하면서 응답을 기다린 뒤에 더 보내는 클라이언트는 `--fallback`을 쓸 수 없습니다. 첫 번째 시도가 한도에 도달하기 전에 바꾼 파일은 그대로 남고, 다음 시도는 프롬프트를 처음부터 다시 시작합니다.

## 관련 문서

- [사용량 한도](/features/usage)
- [사용량 한도와 계정 전환](/docs/usage-limits)
- [MCP 서버, 설정, 대화 내역 공유](/docs/sharing)
- [Claude Code 사용량 한도에 도달했을 때 할 일](/guides/usage-limit-workflow)
