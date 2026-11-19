---
title: 자주 묻는 질문
metaTitle: "ccam FAQ: Claude Code 여러 계정 쓰기"
description: ccam, Claude Code 계정, 사용량 한도, 공유, 개인정보, 지원 시스템에 관한 질문과 답변입니다.
updated: 2026-11-19
order: 6
translation: machine
sourceHash: "090ecc3a34bd"
---

## ccam은 Anthropic 공식 도구인가요?

아닙니다. ccam은 비공식 오픈소스 프로젝트로, Anthropic과 제휴 관계가 없으며 Anthropic의 보증을 받지 않았습니다. ccam은 [Claude Code 문서](https://code.claude.com/docs/en/env-vars)에서 여러 계정을 함께 실행하는 방법으로 안내하는 `CLAUDE_CONFIG_DIR` 설정을 사용합니다.

## Claude Code 계정을 여러 개 써도 되나요?

각 계정의 약관과, 회사 계정이라면 소속 조직의 정책에 따라 다릅니다. Pro와 Max는 Anthropic의 [Consumer Terms](https://www.anthropic.com/legal/consumer-terms), Team과 Enterprise는 [Commercial Terms](https://www.anthropic.com/legal/commercial-terms), 그리고 모든 요금제에 [Usage Policy](https://www.anthropic.com/legal/aup)가 적용됩니다. 여러 계정을 쓰기 전에 확인하세요. ccam은 계정의 한도를 올리거나 합치거나 숨기지 않으며, 모든 계정은 여전히 자기 한도에서 멈춥니다. ccam은 사용자가 요청할 때, 또는 [`auto_switch`](/features/next#auto-switch)(Linux와 macOS)를 켜 두었다면 한도에서 스스로, 작업을 사용자의 다른 계정으로 옮깁니다. 쓰기 전에 이것이 사용자의 계정에 허용되는지 확인하세요. 결정은 사용자의 몫입니다.

**Team·Enterprise 시트.** 자동 계정 선택(`ccam next`, `--fallback`, `ccam pick`, [`ccam agent`](/features/agent), 대시보드의 `n`, 그리고 켜 두었을 때의 [`auto_switch`](/features/next#auto-switch))은 같은 Claude Team이나 Enterprise 조직의 시트 사이에서도 계정을 바꿉니다. Anthropic 고객센터는 "Each seat on a Team plan is assigned to one named individual, and only that person may use it"(팀 플랜의 각 시트는 한 명의 지정된 개인에게 할당되며, 그 사람만 사용할 수 있습니다), "Users can't route their work through another user's account to use that user's remaining usage"(사용자는 다른 사용자의 계정을 통해 작업을 라우팅하여 해당 사용자의 남은 사용량을 사용할 수 없습니다)라고 안내합니다([Team 요금제](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129), [Enterprise 요금제](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)도 마찬가지). 따라서 시트 사이에서 자동 계정 선택을 쓰면 Anthropic의 약관을 위반할 수 있습니다. 시트 사이에서 이 기능을 쓸지는 사용자가 판단할 일입니다. 의존하기 전에 소속 조직과 Anthropic에 확인하세요. 개인 요금제(Pro, Max)에 대해서는 위 문단을 참고하세요.

## ccam으로 사용량 한도를 우회할 수 있나요?

어떤 계정도 자기 한도를 넘어서 쓸 수는 없습니다. 각 계정은 자기 한도를 그대로 지니고 그 한도에서 멈추며, ccam은 한도를 올리거나 합치거나 숨기지 않습니다. ccam은 한도를 보여 주고, 한 계정의 한도가 소진되면 사용자의 다른 계정에서 작업을 이어 가게 해 줍니다. `auto_switch`를 켜 두면 이 일을 스스로 합니다(다음 질문 참고). 한 Team이나 Enterprise 조직의 시트 사이에서는, 다른 시트에서 이어 가면 그 시트의 남은 사용량을 쓰게 되는데, Anthropic 고객센터는 다른 사용자의 시트에 대해 이를 허용하지 않습니다(위의 Team·Enterprise 시트 안내 참고).

## ccam이 세션이 한도에 걸렸을 때 스스로 계정을 바꿀 수 있나요?

직접 켰을 때만 그렇습니다. `~/.config/ccam/config.toml` 맨 위에 `auto_switch = true`를 두면(Linux와 macOS 전용, 기본값은 꺼짐, 실험적 기능), Claude Code가 대화형 세션(ccam이 시작한 세션)을 5시간, 주간, Opus, Sonnet 또는 Fable 한도에서 멈출 때 ccam이 그 세션을 끝내고 여유가 있는 사용자의 다른 계정에서 대화를 이어 갑니다. 옮겨 가는 곳은 같은 Team이나 Enterprise 조직의 다른 시트, 또는 다른 개인 요금제(Pro, Max)입니다. 한도 전에 전환하는 일은 절대 없고, 백그라운드에서 사용량을 지켜보지도 않으며, Team이나 Enterprise 조직 밖으로 또는 그 안으로 작업을 옮기지도 않습니다.

시트에서는 한 시트에서 다른 시트로 작업을 옮기므로 Anthropic의 약관을 위반할 수 있습니다(위의 Team·Enterprise 시트 안내 참고). 개인 요금제에 대해서는 Anthropic의 약관이 이를 분명히 정하지 않습니다. Consumer Terms는 Anthropic이 명시적으로 허용하지 않는 한 "through automated or non-human means"(자동화된 수단이나 사람이 아닌 수단으로) 서비스에 접근하는 것을 금지하며, 한 번 옮길 때마다 ccam이 메시지 하나를 사용자 대신 보냅니다. 켤지는 사용자가 판단할 일입니다.

모든 조건과 한 번 옮길 때 드는 비용은 [`auto_switch`](/features/next#auto-switch)에 정리되어 있습니다.

## 컴퓨터 한 대에서 Claude Code 계정 두 개를 쓰려면 어떻게 하나요?

`CLAUDE_CONFIG_DIR`로 계정마다 설정 폴더를 따로 지정하면 됩니다. 이 작업은 ccam이 대신 해 줍니다. 계정마다 `ccam add`를 실행한 뒤 `ccam personal`, `ccam work`로 시작하거나, `ccwork`처럼 직접 이름을 붙인 단축 명령으로 시작하세요. [가이드](/guides/multiple-claude-code-accounts)에서는 직접 설정하는 방법도 설명합니다.

## 다른 계정에서 대화를 이어 가려면 어떻게 하나요?

폴더에서 `ccam next --continue`를 실행하세요. 그 폴더의 마지막 대화를 다음 계정으로 넘깁니다(여유가 가장 많은 계정으로 넘기려면 `--auto`). Claude Code 안에서는 [ccam 플러그인](/docs/install#claude-code-plugin)의 `/ccam:switch`가 지금 있는 세션에 대해 같은 일을 합니다. 두 계정이 대화 내역을 공유하면(`ccam group set default --add history`, 또는 두 계정만의 그룹: `ccam group add pair --items history --members work,side`) 새 계정이 바로 그 대화를 엽니다. 공유하지 않으면 Claude Code가 대화를 새 계정의 내역으로 복사합니다. `--brief`를 붙이면 새 계정이 대화 전체 대신 짧은 요약으로 시작합니다.

## ccam이 제 로그인 정보를 보거나 외부로 보내나요?

ccam은 계정 사이에 Claude 로그인 정보를 절대 공유하거나 복사하지 않으며, MCP 서버의 로그인도 복사하지 않습니다. 계정의 액세스 토큰은 Anthropic의 사용량 엔드포인트에 해당 계정의 사용량을 조회하는 용도로만 읽습니다. `ccam usage`, 대시보드, `ccam doctor`, 자동 선택(`ccam next --auto`, `ccam pick`, `--fallback`, 그리고 세션이 한도에 걸렸을 때의 `auto_switch`)이 그 경우입니다. 만료된 토큰은 `ccam usage --refresh`를 실행할 때만 갱신합니다. `usage_source = "claude"`이면 토큰을 어디에도 보내지 않으며, 저장된 로그인(파일 또는 macOS 키체인 항목)은 아직 유효한지 확인하는 용도로만 읽습니다. `ccam ext`와 `ccam doctor`는 계정의 MCP 로그인 이름만 읽고 토큰은 읽지 않습니다. ccam에는 원격 측정(telemetry) 기능이 없습니다. 이 조회 외에 ccam이 스스로 보내는 요청은 대시보드가 GitHub에서 새 ccam 릴리스가 있는지 확인하는 것 하나뿐이며, 이 요청에는 사용자나 계정에 관한 정보가 전혀 담기지 않습니다.

## 원격 서버에서는 어떻게 로그인하나요?

SSH 세션에서 평소처럼 `ccam add`나 `ccam login <name>`을 실행하세요. Claude Code는 서버에서 로그인을 기다리므로, 내 컴퓨터의 브라우저는 결국 로드되지 않는 `http://localhost:<port>/callback?code=…` 페이지에서 멈춥니다. 브라우저에서 그 주소를 복사해 터미널에 붙여 넣으세요. ccam이 주소를 서버의 Claude Code에 넘기고, Claude Code가 로그인을 마칩니다. 주소에는 일회용 코드가 들어 있으니 다른 사람과 공유하지 마세요.

## 계정은 어디에 저장되나요?

`~/.local/share/ccam/profiles/<name>/`에 계정마다 폴더가 하나씩 만들어집니다. macOS에서는 Claude Code가 로그인 정보 자체를 키체인에 보관하며, 항목 이름은 그 폴더를 바탕으로 정해집니다. 그래서 `ccam rename`은 폴더의 원래 이름을 그대로 둡니다.

## 기존 ~/.claude는 어떻게 되나요?

`ccam import`나 `ccam migrate`를 실행하지 않는 한 아무것도 바뀌지 않습니다. 가져오기는 폴더를 원래 위치에서 그대로 등록하므로, `claude`를 그냥 실행하든 `ccam <name>`으로 실행하든 같은 계정과 설정을 계속 사용합니다. `ccam migrate`는 폴더를 읽기만 하면서 설정, `CLAUDE.md`, 확장 기능, MCP 서버를 공유 `default` 그룹으로 복사합니다.

## ccam은 어떤 시스템을 지원하나요?

Linux(Ubuntu, Debian, Fedora, Arch, Alpine 등), macOS, Windows 10과 11을 지원하며, Windows에서는 기본 환경과 WSL 안 모두에서 쓸 수 있습니다. Windows에서 `settings.json` 같은 파일을 공유하려면 개발자 모드가 필요하고, 폴더 공유는 개발자 모드 없이도 됩니다. 개발자 모드가 없으면 처음 `ccam add`를 실행할 때 MCP 서버와 확장 기능만 공유하는 `default` 그룹을 만듭니다([Windows에서 공유](/features/share#on-windows)).

## 왜 MCP 서버가 한 계정에서만 보이나요?

세션 안에서 `claude mcp add`로 추가한 서버는 해당 계정에만 저장되고, `-s user` 없이 추가했다면 현재 폴더에만 저장됩니다. 다음에 그 계정을 시작하면 ccam이 알려 주고 `ccam mcp list`와 `ccam ext`에도 표시되니, `ccam mcp import <account> <server>`로 서버를 공유 그룹으로 옮기세요.

## ccam을 쓰면 Claude Code가 느려지나요?

ccam을 거쳐 Claude Code를 시작하면 약 4밀리초가 더 걸립니다. 노트북에서 `scripts/perf-check.sh`로 측정한 값입니다. ccam은 작은 파일 하나를 읽고 공유 링크를 확인한 다음 자기 프로세스를 `claude`로 바꾸므로, 기본적으로는 작업하는 동안 실행되고 있지 않습니다. Windows에서는 프로세스를 바꿀 수 없으므로, ccam이 Claude Code를 시작하고 끝날 때까지 기다릴 뿐 다른 일은 하지 않습니다. `relay`나 `auto_switch`를 켜 두면 ccam은 다른 계정에서 대화를 이어 갈 수 있도록 같은 방식으로 각 대화형 세션과 함께 계속 실행됩니다.

## ccam은 어떻게 업데이트하나요?

`ccam update`를 실행하세요. 최신 릴리스를 내려받아 체크섬을 확인한 뒤 `ccam` 파일을 교체합니다. `ccam update --check`는 새 릴리스가 있는지 알려 주기만 합니다. 대시보드도 새 릴리스가 있으면 알려 주며, `U`를 누르면 확인을 받은 뒤 설치합니다. ccam이 스스로 업데이트하는 일은 없습니다. `.deb` 패키지로 설치했다면 새 `.deb`를 apt로 설치하세요. 자세한 내용은 [업데이트](/docs/install#update)를 참고하세요.

## ccam은 어떻게 제거하나요?

`ccam uninstall`을 실행하세요. 지울 항목을 보여 주고 먼저 확인을 받은 뒤, 모든 계정(먼저 로그아웃), 단축 명령, `~/.config/ccam`, `~/.local/share/ccam`, `~/.local/state/ccam`, `ccam` 명령을 지웁니다. 가져온 `~/.claude`는 남고, 그룹에서 공유하던 항목의 사본을 돌려받습니다. 자세한 내용은 [설치](/docs/install#uninstall)를 참고하세요.
