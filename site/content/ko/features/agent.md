---
title: 다른 계정의 하위 에이전트
description: ccam agent는 Claude Code 작업의 일부를 다른 계정으로 실행하고 Claude의 응답을 돌려주므로, 긴 작업 하나가 여러 계정의 사용량을 쓸 수 있습니다.
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Claude Code 자체의 하위 에이전트는 자신을 시작한 세션의 계정으로 실행되므로 그 계정의 사용량을 씁니다. `ccam agent "<task>"`는 같은 폴더에서 다른 계정으로 `claude -p`를 실행해 작업을 처리하고 Claude의 마지막 응답을 출력합니다. Claude Code 세션에서 호출하면 다른 계정의 사용량을 쓰는 하위 에이전트처럼 동작하므로, 세션 자체의 계정이 바닥났을 시점 이후에도 긴 작업 하나를 이어 갈 수 있습니다.

**Team·Enterprise 시트.** `ccam agent`는 계정을 스스로 고릅니다. 하나의 Claude Team이나 Enterprise 조직 안에서는 작업을 다른 시트에서 실행하며, 이는 Anthropic의 약관을 위반할 수 있습니다. 의존하기 전에 [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)를 읽고 소속 조직과 Anthropic에 확인하세요.

## 사용법

Claude Code는 Bash 도구로 이 명령을 실행합니다. [ccam 플러그인](/docs/install#claude-code-plugin)을 설치했다면 "이 작업의 일부는 내 다른 계정들로 처리해 줘"처럼 평범한 말로 요청하거나, 그 내용을 CLAUDE.md에 적어 두세요. 플러그인의 스킬은 Claude Code에게 각 에이전트를 백그라운드에서 시작하고, 그것만 읽어도 이해되는 작업을 쓰고, 파일을 고치는 에이전트에게는 각자의 몫이나 별도의 git worktree를 주라고 알려 줍니다.

```sh
ccam agent "find why TestParse fails and fix it; reply with what you changed"
ccam agent --model opus -- "-v breaks the parser: find out why"
ccam agent --json < task.md            # the task from standard input
ccam agent --resume 7f3c9a2e-… "now add a test for it"
ccam agent --wait=1h "…"               # nobody has room: wait for a reset
ccam --dry-run agent "…"               # the account and command it would run
```

```console
$ ccam agent "summarize what internal/usage does"
ccam: agent runs as work (12% of its 5 hours used); session 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f
internal/usage looks up the 5-hour and weekly usage of every account …
ccam: agent finished as work; follow up with: ccam agent --resume 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f "…"
```

Claude의 마지막 응답은 표준 출력으로, `ccam:`으로 시작하는 줄은 표준 오류로 나갑니다. `--json`은 대신 `{account, session_id, result, is_error, accounts}`를 출력하며, `accounts`에는 작업이 실행된 모든 계정이 나열됩니다.

종료 코드: claude 자체의 종료 코드(에이전트가 끝나면 `0`), 호출한 계정 말고는 여유 있는 계정이 없으면 `75`(`EX_TEMPFAIL`), 에이전트가 다른 에이전트를 시작하려는 경우 같은 그 밖의 오류는 `1`입니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--model <model>` | 에이전트의 모델입니다. 계정을 고를 때 그 모델의 주간 구간도 따집니다 | 계정의 기본 모델 |
| `--resume <id>` | 이전 에이전트의 대화를 후속 작업으로 이어 갑니다(첫 줄과 마지막 줄에 나오는 세션 ID) | 꺼짐 |
| `--json` | 계정, 세션, 응답, 작업이 실행된 계정들을 출력합니다 | 꺼짐 |
| `--max-usage <1-100>` | 사용량이 이 값 미만인 계정을 우선합니다 | `auto_max_usage`, 없으면 100 |
| `--strategy headroom\|expiring\|balance` | [선택 규칙](/features/next#auto) | `auto_strategy`, 없으면 `headroom` |
| `--cross-org` | 다른 조직의 계정으로도 실행합니다 | `cross_org`, 없으면 꺼짐 |
| `--wait[=<duration>]` | 여유 있는 계정이 없으면 첫 계정이 초기화될 때까지 최대 이 시간만큼 기다립니다 | 꺼짐, 값 없이 주면 `5h` |
| `--dry-run` | 고를 계정에 대한 명령을 실행하지 않고 출력합니다 | 꺼짐 |

## 동작 방식

1. **호출한 계정은 빠집니다.** `ccam agent`를 호출한 계정(ccam이 시작한 세션의 계정, 없으면 이 폴더가 시작하는 계정)은 에이전트를 절대 실행하지 않으며, 같은 Claude 계정으로 로그인해 한도를 함께 쓰는 다른 폴더도 마찬가지입니다. `--cross-org`를 주지 않으면 에이전트는 호출한 계정의 [조직](/features/next#auto) 안에 머뭅니다.
2. **`ccam pick`처럼 고릅니다.** 사용량은 [`ccam pick`](/features/pick)과 같은 방식으로 조회하며, 서로 2분 안에 시작한 에이전트들은 여유 있는 계정이 충분한 동안 서로 다른 계정을 받습니다.
3. **세션의 권한으로 실행하며, 그 이상은 절대 받지 않습니다.** `claude -p`에는 도구 사용을 허용해 줄 사람이 없으므로, 에이전트는 자신을 호출한 세션의 권한 모드와 추론 수준을 받습니다. 터미널에서 실행하면 Claude Code의 기본 모드와 설정의 권한 규칙이 적용됩니다. `ccam agent`에는 더 많은 권한을 주는 옵션이 없습니다.
4. **사용량 한도에 도달하면 대화가 다른 계정으로 이어집니다.** 에이전트의 계정이 한도에 도달하면, ccam은 [플러그인 훅](/features/usage#the-limit-record)이 하듯이 한도를 기록하고 여유가 가장 많은 계정으로 같은 대화를 이어 갑니다. Claude Code가 대화를 그 계정으로 복사하며(`--resume <file> --fork-session`), 첫 메시지는 "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete."입니다. 계정마다 한 번씩만 시도합니다.
5. **모든 에이전트 대화를 기록합니다.** 기록은 `~/.local/state/ccam/agent/<session id>.json`에 30일 동안 남습니다. `--resume <id>`는 그 대화를 원래 폴더에서, 대화를 실행한 계정에 아직 여유가 있으면 그 계정으로, 없으면 다른 계정으로 이어 갑니다. [`ccam next --continue`](/features/next)는 에이전트 대화를 빼므로 계속 사용자의 대화를 이어 갑니다.

## 안전장치와 예외 상황

- 에이전트는 자기 에이전트를 시작할 수 없습니다. ccam이 에이전트에 `CCAM_AGENT=1` 표시를 남기고, 그 안에서는 `ccam agent`가 실행을 거부합니다.
- `ccam agent`를 멈추면(Ctrl+C, 시그널, 또는 Claude Code가 명령을 중단할 때) 에이전트도 멈춥니다.
- 에이전트는 자신을 호출한 대화를 보지 못하고 작업 내용만 봅니다.
- 같은 파일을 동시에 고치는 에이전트들은 서로의 작업을 덮어씁니다. 각자에게 작업의 다른 몫이나 별도의 git worktree를 주세요.
- 호출한 세션이 다른 계정으로 옮겨 가면([`auto_switch`](/features/next#auto-switch)) 그 세션이 시작한 에이전트도 함께 멈춥니다. 대화는 남아 있으므로 `--resume <id>`로 이어 갈 수 있습니다.
- `agent`는 명령 이름이므로 새 계정의 이름을 `agent`로 지을 수 없습니다. 이미 있는 계정은 `ccam run agent`로 계속 시작할 수 있습니다.

## 관련 문서

- [스크립트용 계정 고르기](/features/pick)
- [다음 계정으로 전환](/features/next)
- [사용량 한도](/features/usage)
- [모든 계정의 백그라운드 세션](/features/agents)
