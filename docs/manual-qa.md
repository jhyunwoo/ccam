# 수동 점검표

실제 계정 2개로 Phase마다 확인합니다. 자동 테스트가 대신할 수 없는 부분(브라우저 로그인, 실제 Claude Code 화면)만 남겼습니다. 확인 날짜와 Claude Code 버전을 함께 적습니다.

## 설치와 추가

- [ ] 새 사용자(빈 HOME)에서 `install.sh` 한 줄로 설치, `ccam version` 출력
- [ ] `ccam add` 마법사: 이름 → 공유 항목 → `Create a shortcut command for <name>? [y/N]` → (y) `Shortcut name [cc<name>]:` (예: `ccwork`) → 브라우저 로그인 → `✓ <name> (<email>) added`와 `  shortcut: <cmd>` (N이면 `  no shortcut …`)
- [ ] 마법사에서 PATH의 다른 프로그램 이름(예: `ls`)이나 `list`를 단축 명령으로 주면 이유와 함께 다시 물음
- [ ] `ccam shortcut <name> <cmd>`로 만들기·이름 바꾸기, `--remove`로 지우기 → `~/.local/bin`의 링크가 따라 바뀌고 새 이름으로 실행됨
- [ ] 두 번째 계정도 같은 방식으로 추가, `ccam list`에 서로 다른 이메일
- [ ] `ccam import`로 기존 `~/.claude` 등록 후 `ccam main`과 그냥 `claude`가 같은 계정

## 업데이트

- [ ] 이전 릴리스를 `install.sh`로 설치(`CCAM_VERSION=v<이전>`)한 뒤 `ccam update --check` → ``ccam <새 버전> is available (you have <이전>). Run `ccam update`.``
- [ ] `ccam --dry-run update`가 받을 주소와 바꿀 경로만 보여주고 `ccam version`은 그대로
- [ ] `ccam update` → `✓ Updated ccam <이전> → <새 버전> (~/.local/bin/ccam)`, `ccam version`과 단축 명령 모두 새 버전, 다시 실행하면 `is up to date.`
- [ ] `.deb`로 설치한 `/usr/bin/ccam`에서 `ccam update` → 바꾸지 않고 `sudo apt install ./ccam_<버전>_linux_<아키텍처>.deb` 안내
- [ ] 네트워크를 끊고 `ccam update` → `could not reach github.com …` 오류, 바이너리 그대로

## 전환

- [ ] 단축 명령을 만든 두 계정이 각자의 단축 명령으로 서로 다른 이메일로 열림 (`/status`)
- [ ] `ccam work --resume`, `ccam 2 -p "hi"` 인자 전달
- [ ] `ccam` 선택 화면: 화살표·Enter, 숫자 키, 80열 터미널에서 줄바꿈 없음
- [ ] 터미널 제목이 `claude · <이름>`

## relay(`relay = true`)

- [ ] Linux·macOS: `ccam p1 --permission-mode plan --add-dir /tmp`로 연 세션에서 Shift+Tab으로 acceptEdits로 바꾸고 `/ccam:switch p2` → Claude Code가 끝나고 `ccam: carrying the conversation on (ccam next --continue)` 뒤 p2로 같은 대화가 열림. 권한 모드는 acceptEdits, `/tmp`가 추가 폴더로 보이고, 끝나는 동안 친 글자는 입력란에 없음(`ccam: what you typed during the switch was discarded.`)
- [ ] 한 계정으로 두 터미널(서로 다른 폴더)에서 대화를 열고 둘 다 `/ccam:switch` → 각 터미널이 자기 대화를 자기 폴더에서 이어 감
- [ ] 대화 내역을 공유하는 다른 계정이 같은 폴더에서 열려 있어도 hop은 `Run anyway?`를 묻지 않음
- [ ] `/ccam:switch` 직후 Claude Code가 끝나는 동안 Ctrl+C → `ccam: switch cancelled. …`로 끝나고 `ccam next --continue`가 그 대화를 이어 감. 그 사이 창을 닫으면 남는 ccam·claude 프로세스 없음(`ps`)
- [ ] 업데이트 전에 시작한 relay 세션(옛 ccam)에서 `/ccam:switch` → 세션이 닫히지 않고 /exit 안내
- [ ] Windows 11 호스트: `/ccam:switch p2` → `type /exit; ccam then carries this conversation on as p2.` → `/exit` → p2로 이어 감
- [ ] Windows 11 호스트: 위 hop 뒤 p2의 Claude Code에서 Ctrl+C(입력 지우기) → 셸 프롬프트가 돌아오지 않고 Claude Code가 그대로 받음. `/exit` 뒤 셸이 claude의 종료 코드를 받음(`$LASTEXITCODE`)
- [ ] Windows 11 호스트: 바로가기(`ccwork`, 심볼릭 링크와 `mklink /H` 하드 링크 둘 다)로 연 세션에서 `/ccam:switch` → `/exit` → 다른 계정으로 같은 대화가 이어 감
- [ ] Windows: npm으로 설치한 Claude Code(`claude.cmd`)로 `ccam p1 -p "a & echo X> x.txt"` → claude가 그 문장을 그대로 받고 `x.txt`가 생기지 않음. 작업 폴더의 `claude.bat`은 실행되지 않음

## 한도에서 자동 전환(`auto_switch = true`)

출시 관문으로 정한 확인입니다(계획 6장 Phase 4, D1). 소유자는 2026-10-25에 이 확인이 모두 열린 채로 v0.12.0을 내기로 정했습니다(auto_switch는 기본 끔, experimental). 아래 항목은 출시 뒤의 확인으로 남습니다. 실제 계정으로 하되 지금 쓰는 세션과 실제 프로필 폴더는 쓰지 않고 `HOME`·`CCAM_HOME`을 scratch 폴더로 둡니다. 시작 전후로 실제 파일의 sha256을 비교하고, 검증하려고 계정을 새로 만들지 않습니다.

출시 관문:

- [ ] Q1·Q3·Q4가 모두 "예"(아래 1·3·4. 결과는 날짜와 Claude Code 버전과 함께 docs/claude-layout.md에). 하나라도 "아니오"면 auto_switch를 내지 않고 계획 8장의 대체 동작을 검토
- [ ] 계정 사이 hop QA(아래 3 후반·6·9·11·13): 같은 pool의 실제 계정 둘 이상(개인 요금제 둘, 또는 한 Team·Enterprise 조직의 좌석 둘)으로, 검증하는 사람이 함. 좌석으로 할지는 소유자가 Anthropic 좌석 규칙(계획 9.2)을 보고 정함
- [ ] 실제 파일 sha256 전후 같음: `~/.config/ccam/config.toml`, `~/.local/state/ccam/handoff/` 목록, `~/.local/state/ccam/usage/*.limit.json`, 각 프로필의 `.credentials.json`·`settings.json`
- [ ] 2026-11-12 이후 재확인: https://www.anthropic.com/legal/aup(11-12 시행본에 문서가 인용한 "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products"가 그대로인지), https://www.anthropic.com/legal/consumer-terms(§3에서 문서가 인용한 "through automated or non-human means, whether through a bot, script, or otherwise"가 그대로인지), 좌석 도움말 셋(12004354, 9266767, 13393991), Pro·Max 도움말(8325606, 11049741)에 새 문장이 생겼는지. 결과와 날짜는 docs/PLAN.md에, 문구가 바뀌었으면 features/next.md의 약관 문단, docs/faq.md의 자동 전환 답과 30개 언어를 고침
- [ ] (선택, 계획 D2) Anthropic 지원 창구의 서면 답: "May one person's own Pro/Max accounts be switched automatically at a usage limit by a local tool that runs the unmodified Claude Code with each account's own login?"
- [x] 소유자의 출시 결정: 2026-10-25, 위 확인 전에 v0.12.0으로 내고 사이트 문서를 그 릴리스와 함께 올림. "experimental"을 뺄지는 위가 끝난 뒤 정함

점검 항목:

- [ ] 1 (Q1) 한 계정의 실제 5시간 한도, subagent 셋이 도는 중: StopFailure 호출 수와 `agent_id`, `effort`, `last_assistant_message` 끝의 ` (Asia/Seoul)`, 훅 시작 뒤 0/200/500/1000 ms에 대화 기록의 한도 행 유무. Claude Code를 업데이트할 때마다 다시
- [ ] 2 (Q2) 화면 녹화: StopFailure와 Claude Code 자체 대기·메뉴의 순서. 훅이 판단하는 동안 Esc로 대기를 취소하면 대화 기록 행이나 Notification `quota_auto_resume_disabled`가 남는지
- [ ] 3 (Q3) 같은 계정의 세션에 TERM: 종료 코드 143, `stty -a` 전후 같음, "Resume this session with" 뒤의 ccam 줄, SessionEnd reason other, `ps`에 앞 세션 프로세스 없음, TERM 뒤 훅의 terminalSequence가 보이는지(참고). hop 검증: 한도부터 새 프롬프트까지 걸린 시간(목표 10초 이하)
- [ ] 4 (Q4) 같은 계정으로 `--resume <id> "Reply with OK only."`와 fork 형식(`--resume <파일> --fork-session "…"`)이 Enter 없이 바로 시작하는지, 이어 가기 프롬프트 뒤 끝낸 단계를 되풀이하지 않는지. "아니오"면 D3에 따라 프롬프트 상수(`handoff.ContinuePrompt`)와 문서(features/next.md의 `text` 블록, SKILL.md)에서 뺌
- [ ] 5 (Q5) acceptEdits, plan, auto, bypassPermissions 세션의 resume 뒤 권한 모드(공유·fork 둘 다), effort 유지, 한도 행의 `perTurnEffort`가 훅 입력의 `effort`와 같은지, 게이트 없는 계정의 `--permission-mode auto`, `--effort max`
- [ ] 6 (Q6) 그 폴더를 연 적 없는 대상(신뢰 대화상자)과 상위 폴더 신뢰, 공유 settings 없는 bypass 수락, 주간 한도(메뉴가 열린 상태)에서 hop, 긴 대기 뒤 "Resume this conversation?" 대화상자, Fable 허용량에서는 옮기지 않음
- [ ] 7 (Q7) `run_in_background`로 띄운 `sleep 600`과 `python3 -m http.server`가 TERM 뒤 남는지(문서는 멈춘다고 씀)
- [ ] 8 Claude Code가 끝나는 동안 Ctrl+C → `ccam: switch cancelled. …`, 그때 창 닫기 → 남는 ccam·claude 프로세스 없음, hop 뒤 Ctrl+Z와 fg, 한 계정의 두 터미널이 같은 폴더에서 동시에 한도
- [ ] 9 (Q9) 대상의 첫 assistant 행의 `cache_read_input_tokens`·`cache_creation_input_tokens`와 ccam이 알린 크기: 개인 요금제 사이 hop, 한 조직의 좌석 사이 hop(좌석끼리 캐시가 남는지), 비교로 같은 계정 `--resume`. 모두 파일을 고친 뒤
- [ ] 10 macOS(Terminal.app, iTerm2, 계정 셋 이상): 훅 시간, Keychain 질문, Ctrl+Z
- [ ] 11 사슬 qa1 → qa2 → qa3: 마지막 계정은 세션을 그대로 두고, 알림이 가장 먼저 풀리는 계정과 시각을 알리며, reset 전에 앞 계정으로 돌아가지 않음
- [ ] 12 (Q10) 같은 organizationUuid·결제 수단의 개인 계정이 첫 계정의 한도 뒤 100%인지(가능할 때). (Q11) Remote Control에 연결된 세션은 옮기지 않고 알림이 그렇게 말함(실제 세션의 `CLAUDE_CODE_BRIDGE_SESSION_ID`로도, 대화 기록의 bridge-session 행으로도). 손으로 다른 계정에서 resume한 뒤 휴대폰 앱 표시와 `bridge-session.ownerAccountUuid`. (Q8, Windows 단계 전) Windows 11 호스트에서 가짜 자식으로 `CTRL_C_EVENT` 방송
- [ ] 13 D7 (b): 한 조직의 좌석 사이에는 hop, 다른 조직의 좌석이나 좌석과 개인 요금제 사이에는 hop 없음(가짜 claude로는 e2e가 지킴). auto_switch가 켜져 있고 좌석일 때 `ccam doctor`가 Anthropic 좌석 규칙 문장과 링크를 보임
- [ ] 14 Windows 11 호스트: `auto_switch = true`여도 세션을 끝내지 않고 `ccam doctor`가 무시된다고 알림. `relay = true`면 `/ccam:switch` 뒤 `/exit`로 이어 감
- [ ] 15 `auto_switch = true`에서 `ccam doctor`: ccam 플러그인이 없거나 꺼진 계정의 경고 줄(그 계정 상자), 좌석이면 Health 탭 80열에서 좌석 인용 줄이 잘리지 않음, 2.1.214보다 낮은 Claude Code가 있으면 그 경고. 다른 계정이 없는 계정(조직에 혼자인 좌석, 좌석들 사이의 개인 요금제 하나)의 실제 한도에서 알림이 `You have no other … to carry it on as.`로 끝남
- [ ] 16 2.1.310에서 한도 뒤 `!` 명령을 치면(기본 `respondToBashCommands`로 모델이 답함) ccam이 세션을 옮기지 않음. Remote Control에 연결된 Claude Code의 셸에서 시작한 ccam(물려받은 `CLAUDE_CODE_BRIDGE_SESSION_ID`)의 세션은 연결되지 않았으면 옮김

## 요약으로 이어 가기(`ccam next --continue --brief`)

- [ ] 실제 계정으로 `ccam next --continue --brief` → 첫 메시지에 요약 파일이 첨부되어 보이고(`Read @"…"`), Claude가 어디까지 왔는지 한두 문장으로 말한 뒤 기다림. `-- "<요청>"`을 주면 그 요청을 따름
- [ ] 새 대화를 `claude --resume`으로 다시 열어도 요약을 앎(첨부 행이 대화 기록에 남음)
- [ ] 공백이 든 state 경로(`@"…"`로 첨부), Windows에서 npm으로 설치한 `claude.cmd`
- [ ] 경로의 `#`: 작업 폴더와 공유하는 폴더에 `#`가 있으면(Windows 사용자 `C#dev`, 또는 `XDG_STATE_HOME=/tmp/c#dev/state`와 `/tmp/c#dev/proj`의 대화) 첫 메시지에 요약이 첨부되어 보임(`@"../state/ccam/brief/…"`). 공유하지 않는 `#`(`XDG_STATE_HOME=/tmp/c#dev`, 대화는 다른 곳)면 메시지가 파일 이름만 대고 Claude Code가 읽기 전에 물음(Manual 모드는 매번, auto 모드는 계정마다 한 번: "Allow this read outside the working directories?"). Manual 모드에서 `-- -p "<요청>"`이면 읽기가 거부되고 Claude는 요약 없이 진행
- [ ] `--restricted`나 `permissions.blockReadsOutsideWorkingDirectories`에서는 Claude가 파일을 읽을 수 없다고 말함(첫 프롬프트의 언급 첨부는 코드로 읽었을 뿐 실행해 보지 않음)

## 대시보드 Agents 탭

- [ ] 실제 계정으로 `ccam work --bg "…"` → Agents 탭에 보임. Enter로 붙고, Ctrl+Z로 셸로 돌아감(대시보드로 돌아오지 않음)
- [ ] `x` → No면 그대로, Yes면 실제 `claude stop`으로 멈추고 행이 사라짐(탭은 `--all` 없이 보여 줌). `ccam attach <id>`로 다시 열림
- [ ] agent view가 꺼진 계정은 `✘ <계정>: could not list its sessions: …` 줄로 보임. job id 없는 `bg` 기록이 얼마나 자주 보이는지

## 다른 계정의 작업자(`ccam agent`)

- [ ] 플러그인 0.4.0을 설치한 계정의 세션에서 "이 작업의 일부는 내 다른 계정들로 처리해 줘" → Claude Code가 `ccam agent "…"`를 Bash의 `run_in_background`로 띄우고, 끝난 뒤 출력의 답을 읽어 이어 감
- [ ] 부른 세션이 bypassPermissions·acceptEdits·기본 모드일 때 작업자가 각각 같은 모드로 돎(기본 모드면 편집이 거절되고 답에 그렇게 나옴). 긴 턴 중간(대화 파일 끝 1 MB에 모드 행이 없음)에서도 같음
- [ ] 작업자 계정이 실제로 한도에 걸릴 때: `-p --output-format json`의 결과 줄이 한도 문장이고, 작업자 대화에 한도 행이 남아 `state/usage/<이름>.limit.json`이 기록되며, `ccam: <이름> reached its usage limit; the agent goes on as <다음>` 뒤 다음 계정에서 같은 대화가 이어짐
- [ ] `ccam agent --resume <id> "…"`가 같은 계정에서(여유가 있으면) 앞 대화를 기억한 채 이어 감
- [ ] 백그라운드 작업자가 도는 동안 Bash 작업을 멈추면(Claude Code에서 그 작업 중지) 작업자 claude도 끝남(`ccam agents`에서 사라짐)
- [ ] 작업자가 돈 폴더에서 `ccam next --continue`가 사람의 대화를 이어 감

## 공유

- [ ] settings 공유: 한 계정에서 `/config` 변경 → 다른 계정에 반영, 링크 유지
- [ ] `ccam mcp add …` 후 두 계정 모두 `/mcp`에 보임
- [ ] 세션 안 `claude mcp add` → 다음 실행 때 안내 → `ccam mcp import`
- [ ] 대화 내역 공유 켜기 → 1번에서 대화 → 종료 → `ccam next --continue`가 2번에서 이어감
- [ ] 같은 폴더에서 두 계정 동시 실행 시 경고
- [ ] 기존 `~/.claude`(설정·CLAUDE.md·스킬·`~/.claude.json`의 MCP)가 있는 상태에서 첫 `ccam add` → `ccam migrate` 안내 한 줄. `ccam --dry-run migrate`로 계획 확인 → `ccam migrate` → 두 계정의 `/config`·`/mcp`·스킬 목록에 기존 설정이 보이고, `~/.claude`와 그냥 `claude`는 그대로

## 사용량

- [ ] `ccam usage`가 claude.ai 설정 화면의 사용률과 일치
- [ ] `ccam statusline install` 후 상태줄에 `[account] · model · 5h N%`
- [ ] 네트워크를 끊고 `ccam usage` → 캐시 값과 "(cached)" 표시

## 점검

- [ ] 공유 항목 자리에 실제 파일을 만들고 실행 → 멈춤 → `ccam doctor --fix` → 백업 후 정상
- [ ] `ccam remove`가 폴더를 휴지통으로 옮기고 단축 명령 정리
- [ ] `ccam remove <계정>` → `logged out first`, 휴지통 폴더에 `.credentials.json` 없음, 그 계정의 refresh token이 폐기됨(다른 계정은 그대로)

## 로그인 보존 (테스트 계정으로만)

- [ ] `ccam usage --refresh`(만료된 테스트 계정): 갱신 뒤 `.credentials.json`의 mcpOAuth가 그대로, 바로 `ccam <계정>`이 로그인 없이 열림. 갱신 중 Ctrl+C도 같음
- [ ] `ccam logout <계정>`: MCP 로그인이 있으면 서버 이름과 함께 묻고, n이면 그대로
- [ ] macOS: `ccam add t1` → 로그인 → `ccam rename t1 t2` → `its folder stays …` 출력, `ccam t2`가 로그인된 채 열리고 MCP 로그인도 남음
- [ ] macOS: `ccam remove t2` → 키체인에서 `Claude Code-credentials-<해시>` 항목이 사라짐(`security find-generic-password -s …`가 44로 끝남), `ccam add t1 --no-login`이 로그아웃 상태
- [ ] macOS: 예전 ccam으로 지운 계정의 항목이 남은 상태에서 `ccam add <그 이름>` → `security delete-generic-password …` 안내로 멈춤
- [ ] macOS: `USER`를 지운 셸(`env -u USER`)에서 로그인한 계정을 `ccam usage`가 같은 셸에서 읽음
- [ ] Windows(Developer Mode 꺼짐): `ccam shortcut p1 ccp1` → install.ps1을 다시 실행 → `OK 1 shortcuts now run the new ccam`, `ccam doctor`의 Shortcuts가 OK. 세션이 열린 채로 install.ps1을 두 번 더 실행해도 성공(`ccam.exe.old-<숫자>`)
- [ ] Windows: `CLAUDE_CODE_FORCE_WINDOWS_CREDMAN=1`로 로그인한 계정이 `ccam doctor`에서 Info(Credential Manager), `ccam next`의 후보에 남고 `ccam login --expired`에 들지 않음

## 기록

| 날짜 | Claude Code | 결과 | 메모 |
| --- | --- | --- | --- |
| 2026-10-18 | 2.1.302 | 부분 | 실제 Claude Code로 MCP 공유(flag·sync), `next --continue`, 상태줄 수집, 사용량 조회, tmux에서 선택 화면·마법사 확인. 브라우저 로그인과 두 번째 실제 계정은 미확인 |
