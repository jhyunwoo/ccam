# CLAUDE.md

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.

---

# ccam — Claude Code 다중 계정 관리 도구

전체 계획과 결정 기록은 docs/PLAN.md, Claude Code 폴더 조사 결과는 docs/claude-layout.md를 먼저 읽는다.

## 규칙
- 실행 파일과 모든 명령은 ccam 하나다
- ccam <이름|번호>는 run의 줄임이다. 프로필 이름은 하위 명령 이름(config.Reserved)과 겹칠 수 없다
- 프로필 이름 뒤의 인자는 파싱하지 않고 배열 그대로 claude에 넘긴다 (internal/app이 cobra 이전에 os.Args 처리)
- Go 1.26 (Charm v2 요구), 모듈 경로 github.com/jhyunwoo/ccam
- 코드는 internal/ 아래 패키지로 나누고 cmd/ccam/main.go는 얇게 유지
- 모든 경로는 CCAM_HOME(테스트용)과 XDG 변수를 따른다
- 사용자의 실제 ~/.claude, ~/.local/share/ccam은 절대 수정하지 않는다
- 파일 쓰기는 fsutil.AtomicWrite, 데이터 이전은 share.Build → Apply 순서
- 토큰 파일(.credentials.json)은 어떤 경우에도 공유 그룹에 링크하지 않는다 (share.Layout 허용 목록)
- 메모(state/handoff)는 세션마다 한 파일(`<계정>/<세션 id>.json`)이다. relay는 자기 `CCAM_RELAY`(`2:<pid>`)와 자식이 맞는 메모로만 hop하고, 훅은 이 형식의 relay 자식에게만 TERM을 보낸다
- 한도 표시(state/usage/<이름>.limit.json)는 Claude Code가 기록한 한도 행에서만 쓰고, 모든 사용량 출처 위에 reset까지 100%로 덮는다. 예외는 한도 5분(usage.liftDelay) 뒤에 잰 사용량 엔드포인트 값(live·claude 저장값, 상태 줄은 아님)이 그 창을 100% 아래로 보일 때뿐이다(리셋권 /limit-reset 등으로 일찍 풀림). Fable 허용량(seven_day_overage_included)은 Fable 모델 창으로 기록한다
- auto_switch(기본 끔, Windows는 무시)는 Claude Code가 한도 행을 기록한 뒤에만, 한 pool 안에서만(같은 Team·Enterprise 조직의 좌석끼리, 또는 개인 요금제끼리) 세션을 옮긴다. 다른 조직으로는 옮기지 않고(cross_org는 적용하지 않음), 같은 Claude 계정의 다른 폴더로도 옮기지 않는다
- ccam agent는 부른 계정(`CCAM_PROFILE`, 없으면 폴더의 계정)과 그 Claude 계정의 다른 폴더로는 돌지 않고, 부른 계정의 pool 안에서만(`--cross-org` 제외) 고른다. 작업자에게는 부른 세션의 권한 모드·effort만 넘기고 그보다 큰 권한을 주는 길(권한 플래그, claude 인자 전달)은 두지 않는다. 작업자는 `CCAM_AGENT=1`을 받고 그 안에서 `ccam agent`는 거절한다. 작업자 대화는 `state/agent/<세션 id>.json`에 기록하고 `next --continue`는 그 대화를 건너뛴다
- 프로그램 출력(도움말, 메시지, 오류, 질문, TUI, install.sh·install.ps1)은 모두 영어다. 한글은 절대 출력하지 않는다. 코드·주석·커밋도 영어. internal/e2e/english_test.go의 TestSourceHasNoHangul(cmd/·internal/의 테스트 외 Go 소스와 install.sh·install.ps1·scripts/*.sh에 한글 금지, 이스케이프 포함)과 TestOutputHasNoHangul(빌드한 바이너리의 모든 도움말과 주요 명령 출력 검사)이 강제한다

## 명령
- 빌드: go build ./cmd/ccam
- 테스트: go test ./...  (internal/e2e는 실제 바이너리를 빌드해 fake claude로 돌린다)
- Windows 테스트(WSL에서): scripts/windows-test.sh  (교차 컴파일한 테스트를 Windows 호스트에서 실행, CI는 windows-latest)
- 린트: golangci-lint run
- 벤치마크: go test -bench=. -benchmem ./...
- 실행 경로 측정: scripts/perf-check.sh (hyperfine), 회귀 비교: scripts/bench-compare.sh <base>

## 성능 규칙 (docs/PLAN.md의 성능 설계)
- ccam <이름> 실행 경로는 10 ms 이하. 이 경로에서는 cobra·bubbletea·lipgloss·net/http를 쓰지 않는다
- 패키지 init()에서 파일·터미널·네트워크에 접근하지 않는다. 의존성 추가 후 TestPackageInitIsCheap 확인 (huh는 이 이유로 뺐다)
- 실행 경로에서 읽는 파일은 config.toml 하나, 링크 확인은 lstat만
- claude 실행은 항상 syscall.Exec, 후처리는 다음 실행 때. 예외는 exec이 없는 Windows, --fallback(-p 재시도), 그리고 relay = true이거나 auto_switch = true(Linux·macOS, relay를 함께 켬)인 대화형 세션뿐이다. 셋 다 launcher.Plan.Run으로 자식 프로세스를 띄우고 종료 코드를 그대로 돌려주며, relay의 hop(`ccam next --continue --note <메모>`)은 Linux·macOS에서는 같은 PID의 exec 사슬, Windows에서는 기다리는 자식 프로세스다
- .claude.json은 통째로 파싱하지 않는다 (profile.TopLevelValue 스캐너)
- 네트워크 호출은 병렬 + 제한 시간 + 캐시
- 실행 경로나 사용량 조회를 바꾸면 benchstat 결과를 커밋 메시지에 남긴다
