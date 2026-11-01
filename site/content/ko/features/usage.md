---
title: 사용량 한도
description: ccam usage는 모든 계정의 5시간 한도와 주간 한도를 한 번에 보여 주고, ccam statusline은 Claude Code 안에 계정 이름과 사용량을 표시합니다.
updated: 2026-11-01
order: 7
translation: machine
sourceHash: "54b296f7da35"
---

`ccam usage`는 모든 계정에 대해 각 한도를 얼마나 썼는지 보여 줍니다. 5시간 한도, 주간 한도, 그리고 모델별 주간 한도(Fable 등)입니다. 한도마다 막대 하나와 초기화까지 남은 시간이 표시됩니다. 긴 작업을 시작하기 전이나 Claude Code가 한도가 가까워졌다고 알릴 때 실행해, 어느 계정에 여유가 남았는지 확인하세요. `ccam statusline`은 Claude Code의 상태 표시줄 명령입니다. Claude Code 안에 계정 이름과 사용량을 표시하고, `ccam usage`가 쓸 수 있도록 수치를 저장합니다. 사용량 뒤에 있는 토큰과 추정 비용은 `ccam ccusage`가 모든 계정의 대화에 [ccusage](https://ccusage.com)를 실행해 보여 줍니다([아래](#tokens-and-costs-ccam-ccusage)).

## 사용법

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts that are not running
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

각 계정은 번호, 이름, 수치를 가져온 시점을 적은 줄로 시작합니다. `(cached)`는 엔드포인트가 아니라 저장된 파일에서 가져온 수치라는 뜻입니다. 그 아래에는 한도마다 한 줄씩 20칸 막대, 비율, 구간이 초기화되기까지 남은 시간이 표시됩니다. 순서는 `5-hour`, `Weekly`(모든 모델), 그다음 `Fable`처럼 모델마다 따로 있는 주간 구간입니다. 계정에 없는 구간은 `-`로 표시됩니다. 비율은 80%에서 노란색, 100%에서 빨간색으로 바뀝니다.

구간이 가득 차서 Claude Code가 요청을 거부했다면, 그 구간은 저장된 수치와 상관없이 초기화될 때까지 100%로 표시됩니다([한도 기록](#the-limit-record) 참고). 터미널에서는 모든 계정이 먼저 `fetching…`으로 한꺼번에 나타나고, 계정이 응답하는 대로 그 줄들을 그립니다.

마지막 줄은 계정이 둘 이상일 때만 나타나며, [`ccam next --auto`](/features/next)가 고를 계정을 알려 줍니다. 같은 규칙을 따릅니다.

- `auto_max_usage`와 `auto_strategy`
- 실행 중인 세션
- `cross_org`를 설정하지 않으면 마지막으로 사용한 계정과 같은 조직의 계정만
- 100%인 구간이 있는 계정은 건너뛰므로, 주간 한도가 찬 계정은 절대 나오지 않습니다

기본값인 `headroom` 전략에서는 `→ Most headroom: <name>`으로, `expiring`이나 `balance`에서는 `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`처럼 이유와 함께 표시됩니다.

**Team·Enterprise 시트.** 자동 계정 선택(`ccam next`, `--fallback`, `ccam pick`, 대시보드의 `n`, 그리고 켜 두었을 때의 [`auto_switch`](/features/next#auto-switch))은 같은 Claude Team이나 Enterprise 조직의 시트 사이에서도 계정을 바꾸며, 이는 Anthropic의 약관을 위반할 수 있습니다. 의존하기 전에 [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)를 읽고 소속 조직과 Anthropic에 확인하세요.

로그인되지 않은 계정은 `not logged in → ccam login <name>`, API 키 계정은 `no limits (API)`로 표시됩니다. 액세스 토큰이 만료되면 저장된 수치가 있으면 그 위에 ``access token expired; `ccam <name>` renews it when it starts``가 표시됩니다. 다음에 시작할 때 Claude Code가 갱신합니다. 리프레시 토큰이 만료되었거나, 서버가 거부한 뒤 Claude Code가 지워서 로그인이 완전히 만료되었으면 `` login expired → `ccam login <name>` ``이 표시됩니다.

`ccam statusline install`을 실행하면 Claude Code에 `[work] · Opus · 5h 24% · 7d 41%` 같은 줄이 표시됩니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--json` | 결과를 JSON으로 출력합니다: `name`, `slot`, `status`, `limits`(계정에 사용 크레딧이 켜져 있으면 `credits` 포함), `source`, `at`, `error`, 그리고 한도 기록이 구간을 100%로 만들 때 `limited`와 `limited_windows`. 계정이 없으면 `[]`를 출력합니다 | 꺼짐 |
| `--watch[=N]` | Ctrl+C를 누를 때까지 `N`초마다 다시 그립니다. `--watch`만 주거나 0 이하의 값을 주면 60초입니다. 값은 `=` 뒤에 씁니다. `--watch 10`은 힌트와 함께 거부합니다 | 꺼짐 |
| `-v`, `--verbose` | 각 계정 아래에 `source:` 줄을 추가해 수치를 어디서 가져왔는지(사용량 엔드포인트, Claude Code, 상태 표시줄, 저장된 파일, 또는 Claude Code가 알린 사용량 한도) 보여 줍니다. 한도 기록이 수치보다 우선하면 `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` 같은 줄이 이어집니다 | 꺼짐 |
| `--no-cache` | 60초 캐시를 무시하고 엔드포인트에 다시 요청합니다 | 꺼짐 |
| `--source live\|claude\|cache` | `live`는 계정의 토큰으로 엔드포인트에 요청하고, `claude`는 Claude Code가 조회하게 하며(아래 참고), `cache`는 저장된 수치만 보여 줍니다 | `usage_source`, 없으면 `live` |
| `--refresh` | 만료된 액세스 토큰을 갱신합니다. 실행 중이 아닌 계정만 대상입니다. 갱신에는 최대 30초가 걸릴 수 있습니다 | 꺼짐 |
| `--dry-run` | 받아들입니다. `ccam usage`는 자체 캐시 말고는 아무것도 바꾸지 않습니다. `--refresh`와 함께 쓰면 `--dry-run cannot be combined with --refresh, which rewrites login tokens`로 거부합니다 | 꺼짐 |

`--source`에 다른 값을 주면 `--source must be live, claude or cache`로 실패합니다. `ccam usage`는 인자를 받지 않습니다.

### Claude Code 자체에서 사용량 받기

`config.toml`의 맨 위(첫 `[[group]]`, `[[profile]]`, `[[folder]]`보다 위)에 `usage_source = "claude"`를 설정하거나 `--source claude`를 주면 ccam은 로그인 토큰을 전혀 읽지 않습니다. 계정마다 계정 폴더에서 `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence`를 시작하고, 표준 입력으로 Claude Code의 `get_usage` 요청을 보냅니다(Claude Code 2.1.270 이상). Claude Code가 자체 로그인으로 사용량을 조회해 5시간, 주간, 모델별 구간을 답합니다.

`--safe-mode` 덕분에 계정의 훅, MCP 서버, 플러그인은 시작되지 않으며, 세션도 저장되지 않습니다. 계정마다 Claude Code를 한 번 시작해야 하므로(약 1초) `live`보다 느립니다. 모든 계정에 동시에 요청하며 계정마다 타임아웃은 20초이고, 60초 캐시는 평소처럼 적용됩니다.

액세스 토큰이 만료되었거나 5분 안에 만료되는 로그인에는 50초를 줍니다. Claude Code가 시작하면서 토큰을 갱신하는데, 갱신이 중간에 끊기면 로그인을 잃을 수 있기 때문입니다. ccam이 조회를 포기하면 SIGTERM으로 Claude Code에 종료를 요청하므로, Claude Code는 시작한 갱신을 마칩니다(Windows에는 SIGTERM이 없으므로 프로세스를 강제로 끝냅니다).

로그인이 완전히 만료되었으면 Claude Code를 시작하지 않고 `` login expired → `ccam login <name>` ``을 표시합니다. Claude Code는 이 요청을 실험적 기능으로 표시하므로, 응답이 멈추면 ccam은 저장된 수치를 보여 주고 `ccam doctor`가 이를 보고합니다. Claude Code가 자체 조회에 실패했다고 답하면, ccam은 한도가 없는 요금제로 표시하지 않고 `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`와 함께 저장된 수치를 보여 줍니다.

`--dry-run`을 주면 ccam은 Claude Code를 시작하지 않고(Claude Code가 계정 폴더에 쓰기 때문입니다) 대신 저장된 수치를 보여 줍니다. `usage_source`에는 `live`나 `claude`를 줄 수 있으며, `cache`는 `--source`에서만 고를 수 있습니다.

`ccam statusline`은 표준 입력으로 Claude Code의 상태 표시줄 JSON을 읽습니다. 옵션은 `--wrap <command>` 하나뿐이며, 기존 상태 표시줄 명령을 지정하면 그 출력의 첫 줄 앞에 계정 이름을 붙입니다. `ccam statusline install`은 옵션을 받지 않습니다. 두 명령 모두 `--dry-run`을 주면 `` --dry-run is not supported by `ccam statusline`; nothing was changed ``를 출력하고 거부합니다.

## 토큰과 비용: `ccam ccusage`

`ccam usage`는 계정이 각 한도를 얼마나 썼는지 보여 줍니다. 반면 [ccusage](https://ccusage.com)는 Claude Code의 대화 로그에 있는 토큰과 추정 비용을 셉니다. ccusage는 그냥 실행하면 Claude Code의 로그를 `~/.claude`와 `~/.config/claude`(또는 `CLAUDE_CONFIG_DIR`에 지정한 폴더)에서만 찾으므로, ccam의 계정 중에서 볼 수 있는 것은 `~/.claude`를 제자리에 가져온 계정 하나뿐입니다. `ccam ccusage`는 대화가 있는 모든 계정의 폴더를 `CLAUDE_CONFIG_DIR`에 지정해 ccusage를 실행하므로, 보고서 하나로 모든 계정을 다룹니다.

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam은 먼저 자신의 옵션을 읽습니다. `-a`(`--account`)는 계정 이름이나 번호를 받으며, 쉼표로 구분하거나 옵션을 여러 번 지정할 수 있습니다. `--dry-run`은 명령을 실행하지 않고 출력합니다. ccam 옵션이 아닌 첫 인자부터는 `bun x`와 `npx`의 `--config`를 빼고(아래 참고) 모두 입력한 그대로 ccusage에 넘깁니다. `ccam ccusage -- --help`는 ccusage 자체의 도움말을 보여 줍니다. ccam은 `ccusage`가 설치되어 있으면 그것을, 없으면 `bun x ccusage`를, 그것도 안 되면 `npx ccusage@latest`를 실행합니다. 뒤의 둘은 ccusage를 npm에서 내려받습니다. 이 둘은 현재 폴더가 아니라 ccam 자신의 폴더(`~/.local/state/ccam/ccusage`)에서 시작하므로, 현재 폴더의 프로젝트가 자신의 `node_modules/.bin`, `package.json`의 `bin`, `.npmrc`의 레지스트리로 이 둘이 실행할 것을 정할 수 없습니다. `ccam --dry-run ccusage`는 그 폴더로 들어가는 `cd`를 보여 줍니다. 셋 다 쓸 수 없으면 Node.js, Bun, ccusage 중 하나를 설치하라고 알려 줍니다. Linux와 macOS에서는 ccusage가 ccam을 대체하고, Windows에서는 ccam이 ccusage가 끝나기를 기다렸다가 그 종료 코드로 종료합니다. WSL에서는 Windows 드라이브에 있는 `ccusage`, `bun`, `npx`를 건너뜁니다. 이들은 계정 폴더 없이 Windows 프로그램으로 실행되기 때문입니다.

- 보고서는 폴더에 대화가 있는(`projects` 폴더가 있는) 계정을 다룹니다. 계정들이 그룹을 통해 공유하는 대화는 한 번만 셉니다.
- `-a`를 쓰면, 그룹을 통해 대화 내역을 공유하는 계정은 그룹 전체의 대화를 함께 가져오며, ccam은 누구의 대화인지 알려 줍니다: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- `blocks`는 5시간 구간을 보여 주는데, 이 구간은 계정마다 따로 있으므로 위 예처럼 계정을 하나만 지정하세요.
- ccusage는 도구를 하나 지정하지 않으면 찾아낸 다른 코딩 도구(Codex 등)도 보고합니다. `ccam ccusage claude daily`는 Claude Code만 보여 줍니다.
- ccam이 실행할 때 ccusage는 설정(`ccusage.json`)을 `./.ccusage`와 계정 폴더에서만 읽습니다. 다른 파일은 `--config`로 지정합니다: `ccam ccusage daily --config ~/.config/claude/ccusage.json`. `bun x`와 `npx`는 ccam 자신의 폴더에서 시작하므로, `./.ccusage/ccusage.json`이 있고 파일을 따로 지정하지 않았으면 ccam이 그 파일을 `--config`로 넘기고, 상대 경로로 준 `--config`는 전체 경로로 바꿔 넘깁니다. 이때 ccusage는 `--pi-path` 같은 다른 상대 경로를 ccam의 폴더를 기준으로 읽으므로, 이런 경로는 전체 경로로 지정하세요.
- ccusage는 경로에 쉼표가 있는 폴더를 폴더 두 개로 읽어 버리므로, ccam은 그런 폴더를 거부하고 해당 계정 이름을 알려 줍니다. 그 계정은 `-a`로 빼세요.
- `ccusage`는 ccam 명령이므로 ccusage라는 이름의 계정은 `ccam run ccusage`로 시작합니다. `ccam doctor`도 이를 알려 줍니다.

## 동작 방식

`ccam usage`가 보여 주는 수치의 출처는 두 곳입니다. 하나는 각 계정 자체의 로그인 정보로 호출하는 Anthropic의 비공개 사용량 엔드포인트 `GET https://api.anthropic.com/api/oauth/usage`이고, 다른 하나는 Claude Code가 상태 표시줄 명령에 넘기는 값입니다. 관련 코드는 모두 `internal/usage`에 있습니다.

1. ccam이 `config.toml`을 읽고 계정마다 `usage.Account`를 하나씩 만듭니다(`cli.usageCmd`). 여기에는 계정 폴더 `~/.local/share/ccam/profiles/<name>/`, 세션 표시 파일(계정 폴더에 있는 Claude Code 자체의 기록과 함께 실행 중인지 알려 줍니다), 제자리에 가져온 계정인지 여부가 들어갑니다.
2. ccam이 계정마다 고루틴을 하나씩 시작해 모든 계정을 한꺼번에 조회합니다(`Fetcher.FetchAll`). 모든 고루틴이 `http.Client` 하나를 공유하므로 엔드포인트 연결이 재사용됩니다. 연결에는 최대 2초, 계정마다 전체 3초가 주어집니다. 리디렉션을 따라가면 토큰이 다른 곳으로 다시 전송되므로 클라이언트는 리디렉션을 따르지 않으며, 응답은 최대 1 MB까지만 읽습니다.
3. 계정마다 처리가 끝나는 대로 그 행을 다시 그립니다. 실패하거나 제때 응답하지 않은 계정은 저장된 수치로 대신하므로, 느린 계정 하나 때문에 다른 계정이 늦어지지 않습니다.

그림에서 각 계정은 자체 고루틴에서 조회되며, 시간 초과된 계정은 대신 저장된 수치를 보여 줍니다.

<!-- motion:usage-fetch -->

`personal`과 `work`는 제때 응답해 실시간 수치를 보여 줍니다. `side`는 3초보다 오래 걸리므로, ccam이 대신 저장된 수치를 보여 줍니다(`fallback()`). 그다음 ccam이 계정마다 블록을 그리고(`usage.Table()`), 마지막 줄에 표시할 계정을 고릅니다(`usage.Best()`).

계정 하나에 대해 ccam은 다음 순서로 확인합니다(`Fetcher.Fetch`).

<!-- motion:usage-cache -->

1. `--source cache`이면 바로 `fallback()`으로 갑니다.
2. `--no-cache`가 없으면, `~/.local/state/ccam/usage/<name>.json`에 있는 실시간 결과가 60초(`CacheTTL`)가 지나기 전에 측정된 것이면 지금 시점을 기준으로 판단해 씁니다. 그 사이에 초기화된 구간은 0%로 표시하고, 시계보다 앞선 시각이 적힌 파일은 최근 것으로 치지 않습니다. `ccam usage`를 연달아 두 번 실행해도 요청은 계정당 한 번입니다.
3. ccam이 `.credentials.json`(또는 macOS 키체인)을 바꾸지 않고 읽으며, `claudeAiOauth` 키만 해석합니다(`readCredentials()`). 구독 로그인이 없으면 `not logged in`이고, `.claude.json`에 API 키가 있으면 `no limits (API)`입니다. 리프레시 토큰이 만료되었거나 Claude Code가 지운 로그인은 `` login expired → `ccam login <name>` ``이며 조회하지 않습니다.
4. 만료된 액세스 토큰은 쓰지 않습니다. `--refresh`가 없으면 그 계정은 ``access token expired; `ccam <name>` renews it when it starts``와 함께 저장된 수치를 보여 줍니다.
5. 그렇지 않으면 ccam이 엔드포인트를 호출합니다. 401이나 403은 로그인 만료로(`login expired: authentication failed (401)`), 그 밖의 상태 코드나 읽을 수 없는 본문은 오류로 처리하며, 둘 다 저장된 수치로 대신합니다. Ctrl+C로 끝낸 조회는 `cancelled`로 표시합니다. 429는 엔드포인트가 이 계정에 대한 요청을 너무 많이 받았다는 뜻입니다. ccam은 `<name>.pause`에 다시 물어볼 시각(`Retry-After`, 없으면 5분, 최대 1시간)을 기록하고, 저장된 수치를 `rate limited by the usage endpoint; asking again after 15:10`과 함께 보여 주며, 그때까지는 `--no-cache`를 붙여도 그 계정의 엔드포인트를 호출하지 않습니다. 곧바로 다시 시도하면 요청 제한이 길어질 뿐입니다.
6. 정상 응답은 측정 시각과 함께 `<name>.json`에 기록합니다.

저장된 수치로는 ccam이 가장 최근에 저장된 조회(`<name>.json`, 또는 Claude Code 조회에서 나온 `<name>.claude.json`)를 쓰고, `<name>.statusline.json`이 그보다 나중에 기록되었으면 상태 표시줄의 5시간·주간 구간을 그 위에 덮어씁니다(`fallback()`). 상태 표시줄에는 모델별 구간이 없으므로, 모델별 구간은 조회가 저장한 그대로 남습니다. 저장된 수치를 보여 주기 전에 각 구간의 초기화 시각을 확인해, 수치를 저장한 뒤 초기화된 구간은 0%로 표시합니다. 이전 사용률은 더 이상 맞지 않기 때문입니다. 위 예시에서 `personal`이 초기화 시각 없이 `0%`로 표시되는 이유가 이것입니다.

엔드포인트의 응답 형식은 두 가지입니다. 새 형식에는 `session`(5시간), `weekly_all`, `weekly_scoped`(모델별) 항목이 담긴 `limits` 배열이 있습니다. 이전 형식에는 `five_hour`, `seven_day`, `seven_day_opus`, `seven_day_sonnet` 객체가 있습니다. ccam은 둘 다 읽습니다(`ParseLive`). `limits` 배열을 우선하고, 배열이 다루지 않는 부분은 이전 형식의 객체로 채웁니다. 한도가 하나도 없는 본문은 형식이 바뀐 것으로 보고하며, `ccam doctor`가 이를 알려 줍니다.

### 한도 기록

Claude Code가 사용량 한도에서 턴을 멈추면 [ccam 플러그인](/docs/install#claude-code-plugin)의 훅이 대화에서 거부된 요청을 읽어, 어느 구간(5시간, 주간, Opus, Sonnet)이 가득 찼고 언제 초기화되는지 알아냅니다. 이를 `~/.local/state/ccam/usage/<name>.limit.json`에 기록합니다. ccam은 모든 출처에서 알아낸 내용보다 이 기록을 우선 적용합니다(`Fetch`). 기록에 있는 구간은 저장된 수치가 아무리 오래되었거나 낮아도 초기화될 때까지 100%로 치므로, `ccam next --auto`, `ccam pick`, `--fallback`, [`auto_switch`](/features/next#auto-switch)는 그 전에 그 계정을 다시 고르지 않습니다.

다른 정보가 없으면 계정은 기록만으로 표시됩니다(`-v`를 주면 `source: the usage limit Claude Code reported`). 아무도 지켜보지 않는 세션(`-p` 실행과 백그라운드 세션)은 기록을 남기지 않으며, 관리자가 언제든 풀 수 있는 조직의 차단(예: 사용량이 소진되었거나 그룹 한도가 $0인 경우)도 기록하지 않습니다. 다른 Claude 계정으로 `ccam login`하면 기록이 지워지며, 파일을 지워도 기록이 사라집니다.

### 토큰 갱신

리프레시 토큰은 갱신할 때마다 바뀝니다. 갱신하면 새 리프레시 토큰을 받고 이전 토큰은 무효가 됩니다. 같은 계정으로 Claude Code가 실행 중일 때 ccam이 로그인을 갱신하면 그 세션은 로그아웃됩니다. 그래서 `ccam usage`는 `--refresh`를 주지 않는 한 토큰을 읽기만 하고 절대 쓰지 않습니다. 이 옵션을 주면 ccam은 다음과 같이 동작합니다(`Fetcher.refresh`).

1. `ccam import`로 제자리에 등록한 계정은 거부합니다. ccam이 모르는 사이에 일반 `claude`가 그 로그인을 쓸 수 있기 때문입니다.
2. macOS 키체인에 저장된 로그인은 거부합니다.
3. 잠금 파일 `~/.local/state/ccam/locks/profile-<name>.refresh.lock`을 잡아, 두 ccam 프로세스가 같은 토큰을 둘 다 갱신하는 일이 없게 합니다. 잠금이 잡혀 있는 동안 그 계정에는 `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)`가 표시됩니다. `ccam rename`, `remove`, `logout`도 이 잠금을 잡으며, 갱신이 진행 중이면 거부합니다.
4. 잠금을 잡은 상태에서 그 계정의 세션이 어떻게 시작되었든(ccam, 일반 `claude`, IDE, 또는 백그라운드 세션과 스스로 로그인을 갱신하는 그 뒤의 데몬) 실행 중이면 거부하고, 다른 프로세스가 이미 토큰을 갱신했으면 요청을 건너뜁니다.
5. Claude Code가 계정 폴더에서 자체 갱신 잠금(`.oauth_refresh.lock`, `<folder>.lock`)을 잡고 있으면 아무것도 보내지 않습니다. Claude Code가 바로 그때 그 로그인을 갱신하고 있기 때문입니다.
6. 리프레시 토큰을 토큰 엔드포인트에 보내고, Claude Code가 자체 갱신에 주는 시간인 최대 30초까지 응답을 기다립니다. 한 번 보낸 요청은 Ctrl+C로도 절대 끊지 않습니다. 서버가 이미 리프레시 토큰을 교체했을 수 있고, 새 로그인은 그 응답에만 담겨 있기 때문입니다. Ctrl+C나 행업은 다른 조회를 끝내며, ccam은 응답을 저장한 뒤 종료합니다(Ctrl+C 뒤에는 130).
7. Claude Code의 쓰기 잠금(`.storage-write.lock`)을 잡고 저장하며, 그 시점의 `.credentials.json`에서 토큰 필드만 바꿉니다. 그동안 Claude Code가 쓴 내용(MCP 로그인)은 남고, 그동안 Claude Code가 갱신한 로그인은 유지되며, 로그아웃은 절대 되돌리지 않습니다. 응답에 새 리프레시 토큰의 수명이 있으면 Claude Code처럼 `refreshTokenExpiresAt`으로 저장합니다.

갱신이 이루어지지 않으면 그 이유를 계정 줄이나 저장된 수치 위에 표시합니다. 예를 들어 `token not renewed: a session is running`, 또는 요청은 보냈지만 응답을 쓰거나 저장하지 못했을 때는 `login may be lost: …`이며, 이때 Claude Code가 로그인하라고 하면 `ccam login <name>`을 실행하라고 안내합니다.

### 상태 표시줄 수집기

Claude Code는 상태 표시줄 명령을 실행하면서 JSON 문서를 파이프로 넘깁니다. ccam은 여기서 모델 이름과 `rate_limits`를 읽습니다.

<!-- motion:usage-statusline -->

1. ccam은 Claude Code를 시작할 때 설정하는 `CCAM_PROFILE`로 계정을 찾고, 없으면 `CLAUDE_CONFIG_DIR`를 계정 폴더와 대조해 찾습니다(`runStatusline`).
2. ccam이 `rate_limits.five_hour`와 `rate_limits.seven_day`를 읽습니다(`ParseStatusLine`). API 키 계정이거나 세션의 첫 응답 전이면 JSON에 사용량 한도가 없으며, 이때는 아무것도 저장하지 않습니다.
3. ccam이 `~/.local/state/ccam/usage/<name>.statusline.json`을 기록합니다(`SaveStatusLine`). 엔드포인트가 응답하지 않으면 `ccam usage`가 이 파일을 씁니다.
4. ccam은 `[name] · model · 5h N% · 7d N%`를 출력합니다. `--wrap`을 주면 같은 입력으로 `sh -c`를 통해 지정한 명령을 실행하고, Claude Code가 상태 표시줄을 보여 주는 방식대로 출력합니다. 비어 있지 않은 모든 줄을 출력하며, 첫 줄 앞에 `[name]`을 붙입니다. 명령이 실패하면 계정 이름만 표시합니다. Windows에서는 Claude Code와 같은 방식으로 명령을 실행합니다. Git Bash가 설치되어 있으면 Git Bash로, 없으면 PowerShell로 실행합니다.
5. 5시간이나 주간 사용량이 90% 이상이고 다른 계정이 있으면, 줄 끝에 `→ ccam next --auto --continue`가 붙습니다. 예: `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

감싼 명령에는 자체 제한 시간이 없습니다. Claude Code가 상태 표시줄 명령에 10분을 주고 새 업데이트가 시작되면 실행을 멈추므로, 느린 명령은 ccam 없이 실행할 때와 똑같이 보입니다. Claude Code가 ccam을 끝내면 그 명령도 끝납니다(Linux와 macOS에서는 그 명령이 시작한 모든 프로세스까지).

`ccam statusline install`은 각 계정의 `settings.json`에 `statusLine`을 기록하며, `padding`, `refreshInterval`, `hideVimModeIndicator` 같은 그 객체의 다른 키는 그대로 둡니다. 계정에 이미 상태 표시줄 명령이 있으면 ccam은 그 명령을 `--wrap '<old command>'`로, Windows에서는 어느 셸에서도 따옴표가 필요 없는 `--wrap-base64 <base64>`로 보존합니다. 상태 표시줄이 이미 ccam을 실행하고 있으면(다른 경로의 ccam이거나 이전 버전이 따옴표로 감싸 기록한 경우) ccam은 감싸지 않고 바꾸며, 그것이 감싸던 명령은 유지합니다.

Windows에서는 ccam 경로를 슬래시와 홈 폴더를 뜻하는 `~`로 기록하며(`~/.local/bin/ccam.exe statusline`), Git Bash와 PowerShell 모두 이 경로를 이해합니다. 따옴표가 필요한 경로는 Claude Code가 쓰는 셸에 맞춰 전체 경로로 기록합니다. Git Bash에는 `'…'`, PowerShell에는 `& '…'`입니다. 그룹을 통해 공유되는 `settings.json`은 링크를 거쳐 한 번만 기록되므로, 그룹의 모든 계정에 적용됩니다. 다시 실행하면 이미 ccam을 쓰는 파일에는 `(already set)`를 출력합니다.

## 안전장치와 예외 상황

- `ccam usage`는 `--refresh`를 주지 않는 한 토큰 파일을 절대 바꾸지 않습니다. 액세스 토큰이 만료되었으면 `ccam <name>`으로 그 계정을 한 번 실행하세요. Claude Code가 갱신합니다. 로그인이 완전히 만료되었으면 `ccam login <name>`이 필요합니다.
- `--refresh`는 실행 중인 계정(ccam 없이 시작한 Claude Code도 포함), 가져온 계정, 키체인 로그인을 절대 건드리지 않습니다. 갱신이 거부되면 저장된 수치와 그 이유(`token not renewed: …`)를 보여 줍니다.
- 사용량 엔드포인트는 공개 문서가 없고, 바뀔 수도 있습니다. 바뀌더라도 저장된 수치는 계속 볼 수 있으며, `ccam doctor`가 형식 변경을 보고합니다.
- 저장된 수치도 없고 응답도 없으면 그 계정 줄에 `lookup failed:`와 이유를 표시합니다(예: `--source cache`일 때 `no saved usage`).
- 캐시 파일에는 사용률과 초기화 시각만 들어 있으며, 토큰은 절대 들어 있지 않습니다.

## 관련 문서

- [사용량 한도와 계정 전환](/docs/usage-limits)
- [다음 계정으로 전환](/features/next)
- [Claude Code 사용량 한도에 도달했을 때 할 일](/guides/usage-limit-workflow)
- [명령어 레퍼런스](/docs/commands)
