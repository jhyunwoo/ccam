---
title: 대시보드
description: 터미널에서 ccam을 실행하면 계정, 그룹, 폴더 규칙, 확장, 점검을 보여 주는 전체 화면 대시보드가 열리며, 모든 변경은 계획을 먼저 보여 줍니다.
updated: 2026-10-30
order: 15
translation: machine
sourceHash: "d65f1b5bb0c8"
---

표준 입력과 출력이 터미널이면 인자 없는 `ccam`은 전체 화면 대시보드를 엽니다. 터미널이 아니면 `ccam list`를 출력합니다. 대시보드는 모든 계정을 사용량, 그룹과 함께 보여 주며, 명령으로 하는 일을 그대로 할 수 있습니다. 실행, 추가, 로그인, 그룹 변경, 폴더 규칙 지정, 백그라운드 세션 목록 보기, 점검과 복구가 모두 됩니다. 실행하거나 백그라운드 세션에 연결하면 대시보드가 닫히고, 계정을 추가하거나 로그인할 때는 잠시 터미널을 넘겨준 뒤 대시보드로 돌아옵니다.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · r to refresh · a to add · ? for shortcuts
```

탭 이름 아래 줄은 현재 탭이 무엇을 하는지 한 문장으로 알려 줍니다. 터미널이 넓으면 Accounts 탭에 각 계정의 이메일도 표시됩니다. 어느 탭에서든 탭 이름 줄 끝에 [`auto_switch`](/features/next#auto-switch)가 켜져 있는지 표시됩니다(Linux와 macOS 전용, Windows에서는 무시됨). 켜져 있으면 Health 탭이 이 기능에 아직 필요한 것을 나열합니다. Claude Code 2.1.214 이상, 그리고 Claude 구독으로 로그인하는 모든 계정에서 켜져 있는 ccam 플러그인입니다.

대시보드는 열릴 때마다 더 새로운 ccam 릴리스가 있는지 백그라운드에서 GitHub에 묻습니다. 첫 화면은 이 응답을 기다리지 않으며, 오프라인이면 아무것도 표시하지 않습니다. 새 릴리스가 있으면 제목 아래 줄에서 업데이트를 권합니다.

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U`는 먼저 확인을 받은 뒤 [`ccam update`](/features/update)와 같은 방식으로 업데이트합니다. 릴리스를 SHA-256 체크섬과 대조한 뒤 한 번에 ccam을 교체합니다. 지금 보고 있는 대시보드는 여전히 이전 버전이므로, 종료한 뒤 `ccam`을 다시 실행하세요. `.deb` 패키지로 설치한 ccam 사본처럼 여기서 업데이트할 수 없을 때는 `ccam update`와 같은 안내를 보여 주고 제안은 그대로 남겨 둡니다.

첫 화면은 `config.toml`만 읽어 100 ms보다 훨씬 짧은 시간에 그리며, 사용량은 계정마다 응답이 오는 대로 채워집니다. 한도마다 막대가 따로 있습니다. 5시간 구간, 모든 모델에 대한 주간 구간, Fable 전용 주간 구간이며, `-`는 그 계정에 해당 구간이 없다는 뜻입니다. 비율은 80%에서 노란색, 100%에서 빨간색으로 바뀝니다. 조회에 실패하면 캐시된 값을 흐리게 표시합니다. Claude Code가 가득 찼다고 알린 구간은 초기화될 때까지 100%로 표시됩니다. 액세스 토큰이 만료되면 흐린 `⚠ token expired`가 표시되고(다음에 시작할 때 Claude Code가 갱신합니다), `⚠ login expired`는 그 계정에 `ccam login`이 필요하다는 뜻입니다.

## 탭

| 탭 | 보여 주는 내용 | 키 |
| --- | --- | --- |
| Accounts | 번호, 이름, 이메일, 5시간·주간·Fable 막대, 마지막 사용, 그룹. `✔`는 기본 계정, `▸`는 이 폴더의 계정, `•`는 실행 중(ccam 없이 시작한 Claude Code도 포함: 일반 `claude`, IDE, 백그라운드 세션과 그 데몬) | 아래 참고 |
| Agents | [`ccam agents`](/features/agents)가 보여 주는 내용: 모든 계정의 백그라운드 세션과 열려 있는 세션, 그리고 각각의 상태, 폴더, 시작 시각(아래 참고) | `enter` 연결(`ccam attach`처럼), `x` 멈추기, `r` 다시 불러오기 |
| Groups | 각 그룹, 그룹이 공유하는 항목과 계정, 어느 그룹에도 속하지 않은 계정 | `enter` 편집, `n` 새로 만들기, `e` 이름 바꾸기, `x` 삭제 |
| Folders | [`ccam folder`](/features/folder) 규칙과 여기에 적용되는 규칙. 심볼릭 링크를 거쳐 폴더에 들어온 경우도 포함(새 규칙은 링크를 푼 경로로 저장) | `n` 여기서 쓸 계정 지정, `enter` 변경, `x` 삭제 |
| Extensions | [`ccam ext`](/features/ext)가 보여 주는 내용. 확장마다 한 번씩, 출처 그룹과 함께 표시. 플러그인이 켜져 있지 않은 계정을 나타내는 `off:`, 이 폴더에만 적용되는 계정 자체 서버, `This folder's project`, 어느 계정에도 속하지 않는 문제도 표시 | `r` 다시 불러오기 |
| Health | [`ccam doctor`](/features/doctor)가 보고하는 내용을 계정마다 상자로 표시(아래 참고) | `f` 선택한 상자 복구, `F` 전체 복구(`doctor --fix`), `↑`/`↓` 상자 사이 이동 |

`←`/`→`나 `Tab`으로 탭을 바꾸고, `?`는 현재 탭의 키를 보여 주며, `r`은 짧은 캐시를 건너뛰고 사용량을 새로 고치며 Agents, Extensions, Health 탭도 다시 불러오고, `q`나 `Esc`는 종료합니다.

## Accounts 탭의 키

| 키 | 동작 |
| --- | --- |
| `Enter` | 선택한 계정으로 Claude Code 시작 |
| `1`–`9` | 그 번호의 계정으로 이동. 12번은 `1` `2`를 빠르게 입력합니다. 실행은 `Enter`로만 합니다 |
| `↑`/`↓`, `k`/`j` | 이동 |
| `n` | [`ccam next --auto`](/features/next)가 고를 계정으로 시작 |
| `o` | 옵션을 주고 시작: `--continue`, `--resume`, 그 밖의 인자, 또는 시험 실행 |
| `u` | 사용량 자세히 보기: 5시간, 주간, 모델별 구간과 초기화 시각, 수치의 출처 |
| `a` | 계정 추가(아래 마법사) |
| `l` | 다시 로그인 |
| `s` | 계정의 그룹 고르기 |
| `d` | 기본 계정으로 지정 |
| `f` | 현재 폴더에서 쓸 계정으로 지정 |
| `e` | 이름 바꾸기 |
| `c` | 단축 명령 지정 또는 삭제 |
| `x` | 삭제(No가 선택된 상태로 먼저 묻습니다) |

**Team·Enterprise 시트.** 자동 계정 선택(`ccam next`, `--fallback`, `ccam pick`, 대시보드의 `n`, 그리고 켜 두었을 때의 [`auto_switch`](/features/next#auto-switch))은 같은 Claude Team이나 Enterprise 조직의 시트 사이에서도 계정을 바꾸며, 이는 Anthropic의 약관을 위반할 수 있습니다. 의존하기 전에 [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)를 읽고 소속 조직과 Anthropic에 확인하세요.

## Health 탭

Health 탭은 열 때 [`ccam doctor`](/features/doctor)의 점검을 실행하고 같은 화면을 그립니다. 첫 줄에 계정마다 표시가 하나씩 있고, 그 아래에 계정마다 상자가 하나씩 있습니다. 상자의 각 행은 점검한 영역 하나이며 `Folder`, `Login`, `Share links`가 먼저 나오고, `✔`는 정상, `ℹ`는 정보, `⚠`는 경고, `✘`는 오류를 뜻합니다. 그래서 어느 계정의 어디에 문제가 있는지 한눈에 알 수 있습니다. 어느 계정에도 속하지 않는 결과는 `General` 상자에 들어갑니다.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓`로 상자 사이를 이동하며, 선택한 상자는 강조색 테두리로 표시됩니다. `f`는 선택한 상자만 `ccam doctor --fix`와 같은 방식으로 복구합니다. 계정 상자에서는 그 계정만 복구하고, `General` 상자에서는 어느 계정에도 속하지 않는 결과(플러그인 경로와 단축 명령)를 복구하며, 나머지 계정은 다시 점검만 합니다. 먼저 물어보며(``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), 복구할 것이 없는 상자에서는 `nothing to repair in <name>`이라고만 표시합니다. `F`는 `ccam doctor --fix`처럼 모든 것을 한 번에 복구하고, `r`은 다시 점검합니다.

## Agents 탭

Agents 탭은 [`ccam agents`](/features/agents)가 출력하는 내용, 곧 모든 계정의 백그라운드 세션과 열려 있는 세션을 각각 속한 계정과 함께 보여 줍니다. 탭을 처음 열 때, `r`을 누를 때, 세션을 멈춘 뒤, 설정이 바뀐 뒤에 모든 계정의 Claude Code에 한꺼번에 묻고(`claude agents --json`, 계정마다 타임아웃 15초), 모두 답할 때까지 `Asking every account's Claude Code for its sessions…`를 표시합니다.

```console
  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Background agents and open sessions of every account; enter attaches to one.

    Account   ID        Name   State          Folder      Started
 ❯  work      -         -      open (busy)    ~/code/web  now
    personal  49ec0028  probe  blocked, idle  ~/code/app  12m

 ✘ lab: could not list its sessions: exit status 1 'claude agents --json' is
   not available in this environment.

 enter to attach · x to stop · r to reload · ? for shortcuts
```

`State`는 `ccam agents`에서처럼 백그라운드 세션의 상태와 그 프로세스가 하고 있는 일이며, 대화형 세션이면 `open`입니다. 계정에 물어볼 수 없으면 목록 아래 `✘` 줄에 Claude Code가 한 말이 나옵니다. 위의 `lab`이 그런 경우로, 에이전트 보기가 꺼져 있습니다. 끝난 백그라운드 세션은 빠지며(`ccam agents --all`이 보여 줍니다), 모든 계정이 답했는데 세션이 하나도 없으면 탭에 ``No sessions. Start one in the background with `ccam <account> --bg "…"`.``가 나옵니다. 터미널이 좁으면 `Started`, `Name`, `Folder` 순서로 열을 빼고, 긴 폴더 경로는 왼쪽을 잘라 냅니다.

백그라운드 세션에서 `Enter`를 누르면 대시보드가 닫히고, `ccam attach <id>`처럼 그 세션을 해당 계정으로 이 터미널에서 엽니다. 대시보드는 다시 돌아오지 않습니다(Claude Code에서 `←`는 에이전트 보기로, `Ctrl+Z`는 셸로 갑니다). 열려 있는 세션에는 ID가 없으므로 `Enter`와 `x`는 백그라운드 세션이어야 한다고만 알려 줍니다. `x`는 먼저 `No, keep it`이 선택된 상태로 물어봅니다: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` 예를 고르면 ccam은 그 계정으로 `claude stop <id>`를 실행하고, ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` 또는 멈추지 못한 이유를 알린 뒤 목록을 다시 읽습니다.

## 그룹 바꾸기

Groups 탭에서 `Enter`를 누르면 편집기가 열립니다. 왼쪽에는 항목, 오른쪽에는 계정이 있으며, `Tab`으로 오가고 `Space`로 켜고 끕니다. 계정 중 하나가 이미 다른 그룹에서 받는 항목은 이유와 함께 흐리게 표시되고, 같은 항목을 두 번 받게 되는 계정도 마찬가지입니다. 항목은 한 그룹에서만 받을 수 있습니다. `Enter`를 누르면 계획이 표시됩니다.

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

그룹의 사본과 내용이 다른 파일이 있으면 먼저 어떻게 할지 묻습니다(그룹의 것을 유지, 내 것을 사용, 병합). 그다음 [`ccam group`](/features/share)과 같은 코드로 백그라운드에서 변경을 실행합니다. 설정을 먼저 저장하고, 계정마다 데이터를 차례로 옮기며, 원본은 `.ccam-backup/`으로 갑니다. MCP 공유에서 빠지는 계정은 그룹의 서버를 자기 `.claude.json`에 남겨 둡니다.

실행 중인 계정의 데이터를 옮기게 되는 변경은 그 계정 이름과 함께 거부되며, 계획을 열어 둔 동안 계정이 시작된 경우도 마찬가지입니다. 다른 ccam 프로세스가 그룹을 바꿔서 낡은 계획은 `the groups changed while this plan waited; nothing was moved, make the change again`로 거부됩니다. 다른 ccam 프로세스가 그룹을 바꾸는 동안 새 변경, 이름 바꾸기, 삭제는 `another ccam process is changing the groups right now`를 표시하고, 변경이 중간에 멈춘 뒤에는 Health 탭의 `F`(또는 `ccam doctor --fix`)가 마무리할 때까지 기다립니다.

## 계정 추가

`a`를 누르면 터미널에서 `ccam add`도 쓰는 마법사가 열립니다.

1. 계정 이름.
2. 로그인 방식: 지금 Claude 구독으로 로그인(브라우저에서), [API 키](/features/api), 또는 나중에 Claude 구독으로 로그인.
3. API 키를 고른 경우: 실행할 때마다 ccam이 키를 읽을 곳(`env:NAME` 또는 `command:COMMAND`), 엔드포인트, 키를 보내는 방식, 모델.
4. 들어갈 그룹. `default`가 선택되어 있습니다(아직 그룹이 없으면 `default` 그룹이 공유할 항목).
5. 원하면 단축 명령.
6. 확인할 요약.

`Esc`는 한 단계 뒤로 가며, 첫 단계에서는 마법사를 닫습니다.

## 계정 삭제

`x`는 무엇이든 하기 전에 No를 선택한 상태로 묻고, 삭제가 무엇을 하는지 알려 줍니다: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` API 키 계정에는 폐기할 로그인이 없고, 다른 계정이 같은 로그인을 쓰고 있으면 로그아웃하지 않습니다. 가져온 계정은 폴더를 그대로 둡니다: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## 색

색은 Claude Code 터미널 UI의 역할을 따르므로, ccam이 Claude Code 옆에서도 어색하지 않습니다. 그룹 이름은 이름에 따라 색이 정해집니다. ccam은 터미널에 배경색을 물어 밝은 테마나 어두운 테마를 고릅니다. 직접 정하려면 `config.toml`의 맨 위, 첫 `[[group]]`, `[[profile]]`, `[[folder]]`보다 위에 지정하세요.

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

`-daltonized` 테마는 적록 색각 이상을 위한 것이고, `-ansi` 테마는 터미널 자체의 16색을 씁니다. `NO_COLOR`를 설정하면 색을 끕니다.

## 안전장치와 예외 상황

- 대시보드는 명령과 같은 함수를 호출하므로 모든 규칙이 그대로 적용됩니다. 데이터를 옮기기 전에 계획을 보여 주고, 실행 중인 계정은 거부하고, 백업을 남기며, 로그인 토큰은 계정 밖으로 나가지 않습니다.
- 대시보드에는 시험 실행이 없습니다. 이름 없이 `ccam --dry-run`을 실행하면 이 폴더의 계정이 어떻게 시작될지 출력하고, 계정에서 `o`를 누르면 그 한 번의 실행만 시험 실행할 수 있습니다.
- Claude Code가 시작되거나 로그인 화면이 열리기 전에 터미널을 원래 상태로 되돌립니다.
- 화면이 대화 상자를 띄우기에 너무 작으면 탭의 행부터 줄입니다. 줄은 바꿈하지 않고 `…`로 자릅니다. 화면 높이가 낮으면 빈 줄과 탭 설명을 그다음으로 줄여 키 안내 줄은 남깁니다. 40열이나 8행보다 작으면 대시보드는 `Too small: ccam needs 40x8 (now 36x7)`만 표시합니다.
- 변경, 삭제, 업데이트, 복구를 적용하는 중에 `Ctrl+C`를 누르면 끝날 때까지 기다립니다. 상태 줄에 `ccam quits when it is done.`가 표시됩니다.
- `U`로 ccam을 업데이트한 뒤에는 계정 추가나 로그인을 마치고 대시보드로 돌아와도 업데이트 제안이 다시 나오지 않습니다.

## 관련 문서

- [그룹으로 공유하기](/features/share)
- [Claude Code 시작하기](/features/launch)
- [계정 관리](/features/accounts)
- [사용량 한도](/features/usage)
