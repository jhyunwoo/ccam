---
title: 스크립트용 계정 고르기
description: ccam pick은 스크립트에 여유가 가장 많은 계정을 출력하고, 함께 시작한 워커에는 서로 다른 계정을 주며, 여유 있는 계정이 없으면 75로 종료합니다.
updated: 2026-11-05
order: 12
translation: machine
sourceHash: "0ec25752ace9"
---

`ccam pick`은 로그인된 계정 중 여유가 가장 많은 계정의 이름을 출력하므로, 스크립트가 그 계정으로 Claude Code를 시작할 수 있습니다. 여유 있는 계정이 충분하면 동시에 시작한 워커는 저마다 다른 계정을 받으므로, `claude -p` 작업 묶음이 한 계정을 다 써 버리지 않고 여러 계정에 나뉩니다.

**Team·Enterprise 시트.** 자동 계정 선택(`ccam next`, `--fallback`, `ccam pick`, 대시보드의 `n`, 그리고 켜 두었을 때의 [`auto_switch`](/features/next#auto-switch))은 같은 Claude Team이나 Enterprise 조직의 시트 사이에서도 계정을 바꾸며, 이는 Anthropic의 약관을 위반할 수 있습니다. 의존하기 전에 [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)를 읽고 소속 조직과 Anthropic에 확인하세요.

## 사용법

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # most weekly quota left per day; fewer sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

종료 코드: 표준 출력에 이름을 내보내면 `0`, 여유 있는 계정이 없거나 로그인된 계정이 없으면 `75`(`EX_TEMPFAIL`), 잘못된 플래그 같은 그 밖의 오류는 `1`입니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--json` | 계정, 그 계정의 5시간·주간 사용량, 수치의 출처, 선택한 이유를 출력합니다 | 꺼짐 |
| `--model <model>` | `opus`처럼 그 모델의 주간 구간도 따집니다 | 없음 |
| `--max-usage <1-100>` | 사용량이 이 값 미만인 계정을 우선합니다 | `auto_max_usage`, 없으면 100 |
| `--strategy headroom\|expiring\|balance` | 5시간 사용량이 최저치에서 3%p 이내인 계정 중 주간 사용량이 가장 낮은 계정, 5시간 구간이 가장 먼저 초기화되는 계정, 또는 하루당 남은 주간 할당량이 가장 많은 계정 | `auto_strategy`, 없으면 `headroom` |
| `--cross-org` | 다른 조직의 계정도 고릅니다 | `cross_org`, 없으면 꺼짐 |
| `--wait[=<duration>]` | 여유 있는 계정이 없으면 첫 계정이 초기화될 때까지 최대 이 시간만큼 기다립니다. 값은 `=` 뒤에 씁니다. `--wait 30m`은 `--wait takes its value after '=': --wait=30m`라는 힌트와 함께 거부합니다 | 꺼짐, 값 없이 주면 `5h` |
| `--dry-run` | 고르기만 하고 선택을 기록하지 않습니다 | 꺼짐 |

## 동작 방식

1. ccam이 이 폴더에서 시작될 계정과 같은 조직에 속한 로그인된 계정 목록을 만듭니다(API 키 계정과 Console 로그인은 한도가 없으므로 뺍니다. [조직](/features/next#auto) 참고). 그 계정과 같은 Claude 계정으로 로그인한 폴더도 뺍니다.
2. 잠금(`~/.local/state/ccam/picks.lock`)을 잡으므로, 동시에 실행된 선택은 하나씩 차례로 처리됩니다.
3. `ccam usage`와 같은 방식으로, 60초 캐시를 쓰며 모든 계정의 사용량을 한꺼번에 조회합니다. Claude Code가 가득 찼다고 알린 구간은 초기화될 때까지 100%로 치므로([한도 기록](/features/usage#the-limit-record)), 방금 한도에 도달한 계정은 저장된 수치에 여유가 있어도 빠집니다.
4. 최근 2분 안에 `ccam pick`이 고르거나 [`auto_switch`](/features/next#auto-switch)가 세션을 옮겨 간 계정(`~/.local/state/ccam/picks/<name>`. 미래 시각이 적힌 선택은 치지 않습니다)을 먼저 빼고, 나머지 중에서 `ccam next --auto`의 [선택 규칙](/features/next#auto)으로 고릅니다. `auto_max_usage` 미만인 계정을 먼저, 전략에 따라 고르며, 100%인 계정은 절대 고르지 않습니다. 이 때문에 규칙이 원래 고를 계정을 건너뛰면 `--json`의 `reason`이 그렇다고 알려 줍니다: `30% of its 5 hours used; work was picked a moment ago`.
5. 여유 있는 계정이 모두 최근에 선택되었다면 규칙이 그 전체 중에서 다시 고르므로, 계정이 셋일 때 네 번째 워커는 한 계정을 함께 씁니다.
6. 선택을 기록하고 잠금을 푼 다음 이름을 출력합니다.

2분이라는 구간이 중요한 이유는 새 작업의 사용량이 수치에 나타나기까지 시간이 걸리기 때문입니다. 이 구간이 없으면 함께 시작한 워커 다섯 개가 모두 같은 계정을 여유가 가장 많은 계정으로 보게 됩니다.

## 안전장치와 예외 상황

- `ccam pick`은 사용량을 읽고 선택 시각을 기록할 뿐입니다. `ccam usage`가 하는 것 이상으로 토큰을 갱신하거나 읽지 않습니다.
- 선택했다고 계정이 예약되지는 않습니다. 2분이 지나면 두 스크립트가 같은 계정에서 시작할 수도 있으니, 그 계정이 한도에 도달했을 때 복구하려면 `--fallback`을 쓰세요.
- `--wait`로 기다리는 동안에는 잠금을 잡지 않으므로, 그사이 다른 선택도 계속 진행됩니다.
- `ccam pick`은 [`ccam folder`](/features/folder)를 무시하고 사용량만으로 고릅니다.

## 관련 문서

- [다음 계정으로 전환](/features/next)
- [한도에서 다시 실행](/features/next#run-again-at-the-limit)
- [사용량 한도](/features/usage)
- [API 키 계정](/features/api)
