---
title: Claude Code 사용량 한도에 도달했을 때 할 일
description: 한도가 초기화될 때까지 기다리거나 다른 Claude Code 계정에서 같은 대화를 바로 이어 가는 방법입니다. 5시간 한도와 주간 한도도 설명합니다.
updated: 2026-11-01
order: 3
translation: machine
sourceHash: "760a36724feb"
---

계정이 사용량 한도(흔히 rate limit이라고 부릅니다)에 도달하면 선택지는 두 가지입니다. 한도가 초기화될 때까지 기다리거나, 가지고 있는 다른 계정에서 같은 대화를 지금 이어 가는 것입니다. ccam에서는 두 번째 방법이 명령 하나 `ccam next --auto --continue`로 끝납니다.

## Claude Code에는 어떤 한도가 있나요?

Pro와 Max 요금제의 Claude Code에는 [5시간 세션 한도와 주간 한도](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan)가 있습니다. 두 한도 모두 웹과 앱의 Claude와 함께 쓰며, 사용량은 두 한도에 동시에 반영됩니다. 둘 중 하나가 소진되면 Claude Code는 "You've hit your session limit · resets 3:45pm" 같은 [메시지를 표시합니다](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit). 한도는 모든 모델에 적용되므로 모델을 바꿔도 소용없습니다.

## 방법 1: 초기화 기다리기

계정이 하나뿐이라면 기다리세요. 최신 버전의 Claude Code는 [열린 세션에서 기다렸다가](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) 한도가 초기화되면 작업을 스스로 이어 갑니다. Claude Code 안에서 `/usage`를 입력하면 요금제 한도를 얼마나 썼는지 볼 수 있습니다.

## 방법 2: 지금 다른 계정에서 이어 가기

두 번째 계정이 있다면 같은 대화를 바로 이어 갈 수 있습니다. 한 번만 설정해 두세요.

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

이제 작업하는 동안 상태 표시줄에 계정과 5시간 사용량이 표시되므로, 갑자기 한도에 걸려 당황할 일이 없습니다.

한도에 도달하면 Claude Code를 종료한 뒤(`/exit`) 다음을 실행합니다.

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

여유가 가장 많은 계정에서 대화를 엽니다.

```sh
ccam next --auto --continue
```

Claude Code가 `work` 계정으로 시작되고, 현재 폴더의 가장 최근 대화를 지금까지 나눈 내용 그대로 불러옵니다. 대화 내역을 공유하면(위의 `ccam group set default --add history`) 두 계정이 같은 대화로 작업합니다. 공유하지 않으면 Claude Code가 대화를 `work`로 복사하고, ccam은 첫 응답이 토큰을 얼마나 다시 읽어야 하는지 알려 줍니다.

Claude Code에 [ccam 플러그인](/docs/install#claude-code-plugin)을 설치해 두면, 한도에 걸려 멈춘 세션에서 알림이 뜨고 가득 찬 구간을 기록하므로 `ccam next --auto`가 초기화 전에 그 계정을 다시 고르지 않습니다. `/ccam:switch`(또는 `/ccam:switch work`)는 `ccam next --continue`를 위해 대화를 기록해 두므로, 정확히 그 대화를 이어 갑니다.

**Team·Enterprise 시트.** 자동 계정 선택(`ccam next`, `--fallback`, `ccam pick`, 대시보드의 `n`, 그리고 켜 두었을 때의 [`auto_switch`](/features/next#auto-switch))은 같은 Claude Team이나 Enterprise 조직의 시트 사이에서도 계정을 바꾸며, 이는 Anthropic의 약관을 위반할 수 있습니다. 의존하기 전에 [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)를 읽고 소속 조직과 Anthropic에 확인하세요.

## ccam이 자동으로 계정을 전환하나요?

직접 켰을 때만 그렇습니다.

- 기본적으로는 작업하는 동안 ccam이 실행되고 있지 않습니다. ccam은 Claude Code를 시작하면서 자기 자신을 Claude Code로 대체하므로 언제 전환할지는 사용자가 정합니다. `ccam next --auto`는 5시간 사용량이 가장 낮은 계정을 고르고, 100%에 도달했거나 로그인하지 않은 계정은 건너뛰며, 한 Team이나 Enterprise 조직 안에서만, 또는 개인 요금제끼리만 고릅니다.
- `config.toml`에 `relay = true`를 설정하면 ccam이 대신 세션과 함께 계속 실행되므로, `/ccam:switch`가 세션의 옵션, 권한 모드, 추론 수준, 에이전트를 유지한 채 스스로 대화를 이어 갑니다. 대화를 다른 조직의 계정으로 옮기기 전에는 ccam이 먼저 묻습니다.
- `auto_switch = true`(Linux와 macOS 전용, 실험적 기능)를 설정하면 ccam이 전환도 직접 합니다. Claude Code가 세션을 한도에서 멈추고 같은 종류의 다른 계정에 여유가 있으면, ccam이 세션을 끝내고 그 계정에서 대화를 이어 갑니다([조건과 비용](/features/next#auto-switch)).

다른 도구 중에는 실행 중인 세션의 로그인을 바꿔 끼우는 것도 있습니다. 그 차이는 [ccam과 다른 계정 전환 도구](/compare/account-switchers)에서 설명합니다.

## 알아 둘 점

- `ccam next --continue`는 작업하던 폴더에서 실행하세요. 대화는 프로젝트 폴더별로 저장됩니다.
- 같은 대화를 두 계정에서 동시에 실행하지 마세요. 그렇게 하려고 하면 ccam이 경고합니다.
- 대화 내역을 공유하지 않으면 대화가 다른 계정으로 복사되고, 그 뒤로 두 사본은 따로 이어집니다. `/rewind`가 코드를 되돌릴 때 쓰는 체크포인트는 처음 계정에 남습니다.
- 새 계정에서 보내는 첫 메시지는 평소보다 한도를 더 많이 쓸 수 있습니다. 그 계정에는 아직 대화가 캐시되지 않았기 때문입니다. `ccam next --continue --brief`를 쓰면 대신 짧은 요약으로 시작합니다.
- 이렇게 계정을 전환해도 되는지는 각 계정에 적용되는 약관에 따라 다릅니다.

모든 옵션은 [사용량 한도](/features/usage)와 [다음 계정으로 전환](/features/next)을 참고하세요.
