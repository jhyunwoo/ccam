---
title: 次のアカウントへの切り替え
description: ccam nextは別のアカウントでClaude Codeを開きます。--autoは余裕が最も多いものを選び、--continueは会話を持ち越し、--fallbackは-pのジョブを再実行します。
updated: 2026-11-05
order: 8
translation: machine
sourceHash: "a320c160caef"
---

`ccam next`は、アカウント名を指定しなくても、別のアカウントとしてClaude Codeを開きます。アカウントが使用量の上限に達したときに使います。切り替える方法は4つあります。

- オプションなしの`ccam next`は最後に使ったアカウントの次のアカウントを、`--auto`は最も余裕があるアカウントを選びます。`--continue`は、2つのアカウントが履歴を共有しているかどうかにかかわらず、このフォルダーの最後の会話を新しいアカウントで続けます。
- Claude Codeの中では、[ccamプラグイン](#inside-claude-code)が、上限でセッションが止まったことを知らせ、`/ccam:switch`が会話を引き渡します。
- `claude -p`を実行するスクリプトでは、`ccam --fallback`がこの切り替えを自動で行います（[上限に達したら再実行](#run-again-at-the-limit)を参照）。
- [`auto_switch`](#auto-switch)を有効にすると（LinuxとmacOS）、Claude Codeが上限で止めた対話型セッションもccamが移します。

**TeamとEnterpriseのシート。** 自動的なアカウント選択（`ccam next`、`--fallback`、`ccam pick`、ダッシュボードの`n`、有効にした場合の[`auto_switch`](/features/next#auto-switch)）は、同じClaude TeamまたはEnterpriseの組織に属するシート間でも切り替えます。これはAnthropicの規約に違反するおそれがあります。頼る前に[FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)を読み、所属する組織とAnthropicに確認してください。

## 使い方

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

ccamは切り替えを`→ <from> → <to>`として標準エラー出力に表示し（`--auto`で選んだ場合はその理由も表示します）、`ccam <name>`と同じ方法で新しいアカウントとしてClaude Codeを起動します。`<from>`は最後に使ったアカウントで、`--continue`の場合は会話が属するアカウントです。`--`の後の引数は、そのまま`claude`に渡されます。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--auto` | ログイン済みのアカウントを使用量で選びます（下記参照）。100%のアカウントはスキップします | オフ |
| `--max-usage <1-100>` | `--auto`と併用。使用量がこの値未満のアカウントを優先します | `auto_max_usage`、なければ100 |
| `--strategy headroom\|expiring\|balance` | `--auto`と併用。5時間の使用量が最小値から3ポイント以内のアカウントの中で週間の使用量が最も少ないもの、5時間の枠のリセットが最も早いもの、または1日あたりの残り週間クォータが最も多いもの | `auto_strategy`、なければ`headroom` |
| `--cross-org` | `--auto`と併用。別の組織のアカウントも選びます（下記参照） | `cross_org`、なければオフ |
| `--wait[=<duration>]` | `--auto`と併用。余裕のあるアカウントがないとき、最初のアカウントがリセットされるまで最長でこの時間だけ待ちます。値は`=`の後に書きます。`--wait 2h`は、`--wait takes its value after '=': --wait=2h`というヒントとともに拒否されます | オフ。値なしで指定すると`5h` |
| `-c`、`--continue` | このフォルダーの最後の会話を、新しいアカウントで続けます（下記参照） | オフ |
| `--brief` | `--continue`と併用。前の会話全体を再開する代わりに、ccamが書いたその要約から、新しいアカウントで新しい会話を始めます（[下記参照](#brief)） | オフ |
| `-- <claude args>` | ccam自身の引数の後に`claude`に渡す引数 | なし |
| `--dry-run` | 対応しています。`→ <from> → <to>`と実行するコマンドラインを表示し、Claude Codeの起動も起動の記録も行いません | オフ |

`--auto`と`--continue`は一緒に使えます。`--continue`を付けると、`--resume`の引数は自分で指定した引数より前に置かれます。

`--dangerously-skip-permissions`や`--model`のようなClaude Code自身のオプションは、`--`の後に置きます。名前の後ろをすべてそのまま渡す`ccam <name>`と違い、`ccam next`は引数を自分で読むので、`--`がないと`ccam: unknown flag: --dangerously-skip-permissions`で止まります。

## しくみ

ccamを通じて起動するたびに、その時刻が`~/.local/state/ccam/last-used/<name>`に書き込まれます。`ccam next`はこれらのファイルを読み、どのアカウントがログイン済みかを確認して、アカウントを選びます。

1. ccamは、`config.toml`からアカウントごとに1つの`next.Candidate`を作ります（`cli.candidates`）。内容は、アカウント、最後に起動した時刻（`launcher.LastUsed`）、使えるログインがあるかどうか（`creds.Has`: `.credentials.json`、またはmacOSではキーチェーンにあるClaudeのログインで、期限が完全には切れておらず、Claude Codeに消去されてもいないもの）です。代わりにAPIキーを使うフォルダー（Consoleのログインや、`ccam add --api`で追加したアカウント）にはプランの上限がないので、選ばれることはありません。
2. `--auto`がない場合、ccamはアカウントをスロット番号順に並べ、最も最近起動したアカウントを見つけます（`next.Rotate`）。それが切り替え元のアカウントです。まだどのアカウントも起動していない場合は、最も小さいスロットになります。
3. そこからスロット順に進み、最後のスロットの次は最初のスロットに戻って、最初に見つかったログイン済みのアカウントを返します。
4. ccamは、`ccam <name>`とまったく同じように、選んだアカウントを準備して起動します（`launcher.LaunchWith`）。アカウントのリンクを確認し、MCPの共有を適用し、セッションを登録して、`syscall.Exec`でccamを`claude`に置き換えます。

探索は最後に使ったアカウントから始まり、その後で最初に見つかったログイン済みのアカウントで止まります。

<!-- motion:next-rotate -->

`work`が最も最近起動されたので、探索はその次から始まります。`side`はログインしていないのでスキップされ、`lab`が選ばれます。`lab`もログアウトしていた場合は、探索は`personal`に戻ります（図の`% len`）。ローテーションが開始元のアカウントを返すことはありません。ほかにログイン済みのアカウントがない場合は、エラーで停止します。

### --auto

`--auto`は、スロット順ではなく使用量で決めます。ccamはログイン済みのアカウントを残し、`ccam usage`と同じ方法（`Fetcher.FetchAll`）で、それらすべての使用量を一度に問い合わせ、次の規則で選びます（`cli.pickAuto`、`usage.Policy.Choose()`）。

<!-- motion:next-auto -->

アカウントは5時間と週間の枠で判断され、Claude Codeへの引数（`--`の後）に`--model`がある場合は、そのモデルの週間の枠でも判断されます（名前に`opus`、`sonnet`、`haiku`、`fable`を含むもの）。ccamはアカウントを3つの段階に分けます（`Policy.Choose()`）。

1. **余裕あり**: すべての枠が`auto_max_usage`（`config.toml`、デフォルトは100）未満。この中から戦略で選びます。`headroom`（デフォルト）は、5時間の使用量が最小値から3ポイント以内のアカウントの中で、週間の使用量が最も少ないものを選びます。そのため、週間の使用量が少なければ12%のアカウントが10%のアカウントより優先されますが、14%のアカウントは優先されません。`expiring`は5時間の枠のリセットが最も早いものを選び、残りのクォータが失われる前に使い切ります。`balance`は、週がリセットされるまでの1日あたりの残り週間クォータが最も多いアカウントを選びます。その際、実行中のccamのセッションをそれぞれ差し引いて数えるので、週間のクォータが使われないまま残ることがなく、並行するセッションも分散されます。
2. **超過**: 100%未満だが`auto_max_usage`以上。余裕のあるアカウントがないときだけ使い、その中で`headroom`が選ぶものが先になります。
3. **満杯**: 100%の枠がある、または5時間の使用量が不明。選ばれることはありません。

Claude Codeがある時間枠でリクエストを拒否した場合、その枠はリセットまで100%と数えます。[ccamプラグイン](#inside-claude-code)のフックがそれを記録し（`~/.local/state/ccam/usage/<name>.limit.json`）、保存済みの数値にかかわらず適用されます。そのため、上限に達したばかりのアカウントは、数値が古くても、リセット前に再び選ばれることはありません。

すべてのアカウントに共通の設定は、`config.toml`の先頭、最初の`[[group]]`、`[[profile]]`、`[[folder]]`より上に書きます。それらの下に書いた設定はそのテーブルに属するので、ccamは起動のたびにそれについて警告し、設定を保存するコマンドは、それが移されるまで拒否します。

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**組織。** 自動的な選択は、切り替え元のアカウントと同じ種類のアカウントの中で行われます。同じTeamやEnterpriseの組織のシートどうし、または個人プラン（Pro、Max）どうしです。そのため、会社のアカウントでの作業が個人のアカウントで続くことも、その逆もありません。ccamは各アカウントの組織をログインから知ります（`ccam login`が記録し、以前のログインは`.claude.json`から読みます）。組織を判別できないアカウントは、どのアカウントとも組み合わせられます。同じClaudeアカウントにログインした2つのフォルダーは上限を共有するので、一方が他方の代わりになることはありません。`--cross-org`か`config.toml`の`cross_org = true`で、組織の規則を外せます。ただし、[`/ccam:switch`](#inside-claude-code)の後の別の組織への移動は、それでもあなたがyesと答えるまで行われません。

選べるアカウントがない場合、ccamは終了コード75（`EX_TEMPFAIL`）で終了するので、スクリプトは「あとで再試行」と失敗を区別できます。`--wait`を付けると、代わりに満杯のアカウントのうち、上限になっている枠が最も早くリセットされるまで（エンドポイントに新しい枠が表示される、リセット時刻の30秒後まで）待ち、`ccam: no account has room; waiting until 15:10`と表示してから、もう一度選びます。どのアカウントもリセット時刻を示さない場合は、10分後にもう一度確認します。待機時間が終わるまでに余裕ができなければ、75で終了します。

問い合わせは`ccam usage`の規則に従います。取得してから60秒未満の数値は再利用し、各アカウントの問い合わせは3秒でタイムアウトし、応答がないアカウントやアクセストークンの期限が切れたアカウントは保存済みの数値で判断し、上限の記録がそれより優先されます。保存済みの時間枠のリセット時刻が過ぎている場合は0%として扱うので、その後リセットされたアカウントが見送られることはありません。`--auto`がトークンを更新することはありません。通常のローテーションとは異なり、`--auto`は、最後に使ったアカウントに最も余裕が残っていれば、そのアカウントを選ぶこともあります。

### --continue

1. ccamは、このフォルダーの最後の会話をすべてのアカウントで探します。対象は、`projects/<folder>/`の中の`<session>.jsonl`ファイルのうち、最も新しいプロンプトか応答を持つものです（フォルダー名は、リンクを解決した作業フォルダーの英字と数字以外のすべての文字を`-`に置き換えたもので、Claude Codeの命名と同じです）。`/mcp`、`/usage`、`/status`を実行しただけのセッションや空のファイルのように、プロンプトも応答もないセッションは飛ばします。セッションがメモ（使用量の上限か`/ccam:switch`によるもの、後述）を残していて、メモから2分を超えた後のプロンプトや応答を持つ会話がない場合は、代わりにまさにそのセッションが使われます。`/exit`のときにClaude Codeが書く行は数えません。
2. 会話が属するアカウントが切り替え元になり、新しいアカウントは上記のとおり、同じ組織の中から選ばれます。複数のアカウントがその履歴を共有している場合、会話は、`.claude.json`でその会話を最後に実行したとされているアカウント（`lastSessionId`）に属し、それがなければ最後に使ったアカウントに属します。
3. 新しいアカウントがその履歴を共有している場合（同じグループのフォルダー）、ccamは`claude --resume <session id>`を起動して、まさにその会話を開きます（`--continue`だと、フォルダーの中で最も新しい会話が開きます）。
4. そうでない場合は`claude --resume <the conversation's file> --fork-session`を起動します。Claude Codeは会話全体を新しいセッションIDで新しいアカウント自身の履歴にコピーし、元の会話はそのまま残ります。それ以降は、新しいアカウントの`--continue`でその会話が見つかります。
5. `--`の後に`--agent`を指定しない限り、ccamは会話が使っているエージェントも引き継ぎます。フォルダーの`.claude/settings.local.json`か`.claude/settings.json`、またはアカウントの`settings.json`の`agent`設定、それがなければ会話の最後の`agent-setting`行が示すエージェントです。Claude Codeは会話をコピーするとき（手順4）エージェントを復元しないので、そのままではデフォルトのエージェントとして、すべてのツールを使える状態で会話を続けてしまいます。新しいアカウントが古いアカウントのエージェント、プラグイン、設定を共有していない場合（グループを通じて同じ`agents`フォルダーと`plugins`フォルダー、`settings.json`を使っていない場合）、ccamは移動についての行の後に次の行を加えます: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`

ccamは会話の大きさを表示します。新しいアカウントは最初の応答で会話全体を読み直すためです。プロンプトキャッシュが組織間で共有されることはなく、個人プランはそれぞれが独自の組織なので、2つのアカウントが別の組織に属する場合、ccamは`prompt caches are not shared between organizations, so the first reply reads it all again`と付け加えます。元のアカウントがサインインしているMCPサーバーに新しいアカウントがサインインしていない場合は、その名前を示します: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`。グループで`history`を共有すればコピーが不要になり、会話は1つのままで、`/rewind`がコードを復元するためのチェックポイント（`file-history/`）も共有されます。

### --brief

`--continue`は会話全体を再開するので、新しいアカウントは最初の応答で会話全体を読み直します。`ccam next --continue --brief`は代わりに、ccamがモデルに頼らずに書いた元の会話の要約から、新しい会話を始めます。

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

要約には、次のものがこの順に、書くことがある場合だけ入ります。

1. 会話のタイトル
2. Claudeが細部を必要としたときに検索するための会話記録のパス
3. 依頼が4つより多かった場合は、あなたの最初の依頼
4. Claude Codeが最後に会話を圧縮したときに書いた要約
5. あなたの最後の4つの依頼（`!`コマンドとスラッシュコマンドは入力したとおり）
6. Claudeの最後の応答
7. ToDoリスト
8. Claudeが最後に変更した10個のファイル
9. フォルダーの`git status --short --branch`（最大800文字）

ToDoリストは最大10項目までで、完了した項目から先に、次に最も先のステップから省くので、現在のステップとその次のステップは残ります。長い部分は中央が切り詰められ、要約全体は16,000文字未満に収まります。

ccamは要約を`~/.local/state/ccam/brief/<session>.md`に、あなただけが読める形で保存し、次の`--brief`が7日より古い要約を削除します。Claude Codeは1つのプロンプトで始まり、そこで言及したファイルを添付します: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` `--`の後にプロンプトを指定すると、それが最後の文の代わりになり（`Then: <your prompt>`）、`--`の後のほかのオプションはその後に続きます。

Claude Codeは言及を`"`で終え、`#`で切ってしまうので、要約のパスにどちらかが含まれる場合、ccamは作業フォルダーからの相対パスで言及します。そのパスにも含まれる場合は、プロンプトでファイル名を示し、ClaudeがReadツールでそれを読みます。ファイルが作業フォルダーの外にあるので、このときは先に確認が入ります（Manualモードの`-p`の実行では、Claude Codeが読み取りを拒否し、Claudeは要約なしで続けます）。

`ccam --dry-run next --continue --brief`は要約を表示し、何も保存しません。元の会話はそのまま残り、新しい会話は要約に書かれたことしか知らないので、そこに至るまでの細部すべてよりもタスクのほうが大事なときに`--brief`を使ってください。`--continue`なしでは、ccamは`--brief works with --continue`で止まります。会話を読めない場合や、要約するものが会話の中に見つからない場合は、止まって、`--brief`なしの`ccam next --continue`なら会話全体を再開できると伝えます。[リレー](#relay)と[`auto_switch`](#auto-switch)は、常に会話全体を再開します。Claude Code自身の`--brief`オプションは、ほかのオプションと同じく`--`の後に置きます。

## Claude Codeの中で

[ccamプラグイン](/docs/install#claude-code-plugin)は、2つのフックを追加します（Claude Code 2.1.78以降）。

- ターンが使用量の上限で止まると、会話をメモに残し、デスクトップ通知に対応したターミナル（iTerm2、WezTerm、Windows Terminal、kitty、Ghostty）で通知を表示します: ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  フックは、Claude Codeが使用量の上限とみなすすべてのメッセージを知っています。5時間、週間、通常、モデルの上限（`You've hit your session limit`、`You've hit your limit`、モデルの割り当てについての`You've reached your …`）、使用量クレジットと支出の上限、そしてアカウントを止めた組織です。

  Claude Codeが拒否されたリクエストを会話に書き込んでいる場合、フックはその満杯の時間枠もリセットまで記録します（[上限の記録](/features/usage#the-limit-record)）。ただし、組織が設定した停止は管理者がいつでも解除できるので記録しません。サーバーの短時間のスロットリングやAPIキーのレート制限は使用量の上限ではないので何もしません。無人のセッション（`-p`の実行とバックグラウンドセッション）も同じです。サブエージェントやワークフローのエージェントが使用量の上限で止まった場合は、メインの会話がそれを待っているだけでも、そのセッションの停止として扱います。Claude Codeはその拒否されたリクエストを、セッションの`subagents`フォルダーにあるそのエージェント自身のファイルに書き込み、フックはそこから読み取ります。複数のエージェントが同時に止まっても、セッションに残るメモと通知は1つずつです。

  Claude Code 2.1.234以降は、5時間の枠がリセットされるとセッション内で待って自動で続行します。通知によって、代わりに切り替えることも選べます。[`auto_switch`](#auto-switch)を有効にしていれば（LinuxとmacOS）、別のアカウントに余裕があるときにccamが代わりに切り替えます。引き継げるアカウントがほかにない場合、つまり同じTeamまたはEnterpriseの組織に属するアカウントがほかにないか、別の個人プランがない場合（`cross_org`を有効にしていれば、ほかのどのClaudeアカウントでも引き継げます）、通知は代わりに`You have no other personal plans, so ccam cannot move this conversation.`または`You have no other accounts in Acme, so ccam cannot move this conversation.`で終わります。
- `/ccam:switch [account]`は会話をメモに残し、リクエストを送る前にターンを終えます。`/exit`と入力した後、`ccam next --continue`を実行すると、指定したアカウント（または次のアカウント）で会話を続けられます。まだ会話のないセッションでは、続けるものがないことを伝えます。これを実行できるのは、あなた自身が入力したときだけです。プラグイン0.3.2以降、Claudeは`/ccam:switch`を実行できず、Claudeがスケジュールしたジョブ（CronCreate、ScheduleWakeup、`.claude/scheduled_tasks.json`）はこれをただのテキストとしてClaudeに渡します。

  指定したアカウントが会話の組織の外にある場合（別のTeamやEnterpriseの組織のシート、シートから個人プランへの移動やその逆、またはどちらか一方がAPIキーのアカウントの場合）、`cross_org`の設定にかかわらず、`ccam next --continue`はまず確認します: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`。アカウントを指定しなかったときに、ccamが自分で選んだアカウント（`cross_org = true`なら別の組織からも選べます）へ移す前にも、同じように確認します: `ccam: the account picked for /ccam:switch is home, a personal plan, …`。移すのは`y`か`yes`と答えたときだけです。LinuxとmacOSでは、質問より前に入力されたキーを先に捨てます。それ以外の答えや、確認に使えるターミナルがない場合は、会話とメモはそのまま残ります: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work``（ccamが自分で選んだ場合は、`ccam next --continue`がもう一度選びます）。自動的な選択と同じく、組織を判別できないアカウントは、どのアカウントとも組み合わせられます。

メモはセッションごとに1つのファイルとして`~/.local/state/ccam/handoff/<account>/<session>.json`に置かれるので、1つのアカウントの2つのセッションが互いのメモを上書きすることはありません。メモは12時間保存され、次のセッションが始まる直前に一度だけ使われます。起動に失敗した場合は、次の試みのために残ります。メモはセッションを開始したフォルダーのものとして扱われ、Claude Codeのシェルで`cd`した後も変わりません。`ccam rename`はメモの中の名前を更新します。メモにあるアカウントがその後削除されていた場合、ccamはそのことを伝え、通常どおりに選びます。

### リレー

`config.toml`の先頭で`relay = true`を設定すると、ccamは自身をClaude Codeに置き換えずに対話型のセッションと並んで動き続け、`/ccam:switch`の後にそのセッションの会話を自動で続けます（別のセッションが残したメモには反応しません）。フックがClaude Codeを終了させ、ccamがそのセッションのメモを使って`ccam next --continue`を実行します。アカウントを指定しない場合は、会話のモデルも数えて、最も余裕があるアカウントを選びます。`cross_org = true`のためにそれが別の組織のアカウントになるときは、会話をそこへ移す前に上の質問をします。WindowsではフックがClaude Codeを終了させられないので、`/exit`と入力すると、ccamが会話を続けます。会話の組織の外にあるアカウントの場合も、フックはClaude Codeを終了させず、`ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`と表示します。あなたが`/exit`と入力すると、ccamが上の質問をします。その分、セッションごとに待機するccamのプロセスが1つ必要になるので、デフォルトではオフです。

引き継いだセッションは、会話の再開では復元されないものを引き継ぎます。`--add-dir`、`--mcp-config`、ツールの一覧、`--settings`、`--plugin-dir`、システムプロンプトのオプションなどのオプション、会話の権限モード、セッションのeffortレベル、会話が使っているエージェント（[`--continue`](#continue)を参照）です。`--model`と`--name`は、Claude Code自身が会話から復元します。それ以外のもの（`--worktree`、`--debug`、`--remote-control`、ccamが知らないオプション）は、`ccam: not carried over: …`で名前が示されます。切り替えで履歴の共有についての質問をすることはなく、切り替え中に入力したものは捨てられます。切り替えでする質問は、別の組織のアカウントに移す前の質問だけです。

リレーが会話を続けるのは、`/ccam:switch`が求めたとおりにClaude Codeが終了した場合だけです。先に`/exit`と入力した場合、終了中にCtrl+Cを押した場合、ccamがハングアップやほかのシグナルを受け取った場合、ターミナルがなくなった場合は続けません。そのときメモは残るので、あとで`ccam next --continue`が会話を続けます。古いバージョンのccamのリレーのもとで始まったセッションや、無人のセッション（`-p`の実行とバックグラウンドセッション）は、フックによって終了されず、代わりに`/exit`のメッセージが表示されます。`claude remote-control`はリレーなしで実行されます。

### `auto_switch`

`config.toml`の先頭で`auto_switch = true`を設定すると（LinuxとmacOS）、Claude Codeがセッションを使用量の上限で止めたときに、ccamがそのセッションを自動で別のアカウントに移します。プラグインのフックが余裕のある別のアカウントを選んでClaude Codeを終了させ、セッションと並んで動いているccamが、[リレー](#relay)が`/ccam:switch`の後に行うのと同じように、そのアカウントで会話を続けます。デフォルトではオフで、実験的な機能です。

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccamがセッションを移すのは、次のすべてが成り立つときだけです。

1. ccamがターミナルでそのセッションを起動し、それと並んで動き続けていること。`auto_switch`は[リレー](#relay)を有効にするので、設定した後に起動した対話型のセッションはこれを満たします。設定する前に、または古いccamのもとで起動したセッションには、``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.``で終わる通知が表示されます。無人のセッション（`-p`の実行とバックグラウンドセッション）、`claude remote-control`、ccamを使わずに起動したClaude Codeが移されることはありません。
2. そのアカウントで[ccamプラグイン](/docs/install#claude-code-plugin)がオンになっており、フックにプロセスIDを渡す（`CLAUDE_PID`）Claude Code 2.1.214以降を実行していること。
3. Claude Codeが、拒否されたリクエストを、リセットの時刻とともに会話に書き込んでいること。対象はそのアカウントの5時間、週間、Opus、Sonnetの上限です。支出の上限、使用量クレジットの使い切り、組織が設定した停止、Fableの割り当て（`You've reached your Fable limit.`）、サーバーのスロットリング、APIキーのレート制限で、セッションが移されることはありません。サブエージェントやワークフローのエージェントの上限は、メインの会話がそれを待っているだけでも対象になります。それが週間のOpusまたはSonnetの上限の場合、セッションが移されるのは、会話そのものがそのモデルで動いているときだけです。
4. ccamがそのアカウントのプランを、`.claude.json`から、または`ccam login`が記録した内容から知っていること。知らない場合、通知は`ccam does not know personal's plan, so auto_switch leaves this session.`で終わります。
5. セッションがRemote Controlに接続されていないこと。別のアカウントになると、操作に使っているアプリに再接続しないので、セッションはそのまま残り、通知がそのことを伝えます。
6. 上限のリセットが今から15分より先であること。待ち時間がそれより短い場合、ccamはセッションをClaude Codeに任せます: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. 上限に達してから、あなたがプロンプトを入力していないこと。
8. 別のアカウントに余裕があること。ccamは[`--auto`](#auto)と同じように（`auto_max_usage`、`auto_strategy`、会話のモデル）、元のアカウントと同じ種類のアカウント、つまり同じTeamまたはEnterpriseの組織の別のシートか、別の個人プラン（Pro、Max）の中から選びます。ここでは`cross_org`と`--cross-org`は適用されません。移す先のアカウントは、次の条件を満たす必要があります。
   - Claude Codeが更新できるログインでログインしていること
   - 同じClaudeアカウントにログインしていないこと
   - APIキーを使っていないこと
   - ccamがプランを知っていること
   - ccam自身の使用量の問い合わせでわかる範囲で、使用量クレジットがオンになっていないこと
   - Claude Codeでこのフォルダーを信頼済みであること
   - 会話がエージェントとして動いている場合（`agent`設定か会話がエージェントを示している場合）、元のアカウントのエージェントを共有していること（`agents`フォルダー、`plugins`フォルダー、`settings.json`のそれぞれが、グループを通じてまったく同じものか、両方にないこと）

   最後の3つの条件によって余裕のあるアカウントがすべて外れる場合、通知がそのアカウントを示します: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`、`side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.`または`This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` [上限の記録](/features/usage#the-limit-record)がある時間枠は満杯と数え、直近2分以内に別のセッションの移動先になったアカウントは、ほかのアカウントの後に回されます。選択は約6秒以内に終える必要があり、使用量の問い合わせがそれより長くかかった場合、セッションはそのまま残ります。`usage_source = "claude"`では、`auto_switch`は保存済みの数値だけで判断し、ccamはClaude Codeの問い合わせやステータスラインから使用量クレジットがオンかどうかを判断できないので、オフとみなします。

その後、フックがClaude Codeを終了させ、それと並んで動いているccamが、デスクトップ通知に対応したターミナルで`ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).`のような通知を表示します。ccamは上の出力例の1行目を表示し、選んだアカウントとしてそのセッションの`ccam next --continue`を実行します。これが残りの2行を表示し、[リレー](#relay)と同じように、セッションのオプション、権限モード、effortレベル、エージェントとともにClaude Codeを起動します。新しいセッションの最初のメッセージは決まった文で、これがccamがあなたの代わりに送る唯一のものです（`/ccam:switch`の後の移動では何も送りません）。

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

サブエージェントが上限に達した後は、サブエージェント、ワークフロー、バックグラウンドのコマンドが古いセッションとともに止まっているので、次の文になります。

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

そのアカウントも上限に達した場合は、そこで同じことが起こります。上限の記録があるので、この連鎖がリセット前のアカウントに戻ることはなく、すべてのアカウントが満杯になると、セッションは最後のアカウントに留まります。

フックは、会話がエージェントとして動いていることを、`agent`設定から、またはClaude Codeがときどき、およびセッションの終了時に会話へ書き込む、エージェントを示す行から知ります。これまでにセッションの`--agent`だけがエージェントを示していて、選んだアカウントがエージェントを共有していない場合、`ccam next --continue`は何も起動せずにメモを残し、2つ目のデスクトップ通知がそのことを伝えます: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``

ほかのアカウントに余裕がない場合は何も移しません。Claude Codeは通常どおり上限を扱い、セッション内で待つか選択肢を示します。通知は、どのアカウントが最初に空くかを伝えます。たとえば`ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.`です。シートの場合は`None of your other accounts in Acme has room now`となります。始まった移動を止めるには、Claude Codeが閉じている間にCtrl+Cを押します: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

移動にはコストがあります。Claude Codeが再起動するので数秒かかり、Claudeが起動したバックグラウンドのシェル、開発サーバー、エージェントは古いセッションとともに止まります。新しいアカウントは、[`--continue`](#continue)と同じく最初の応答で会話全体を読み直し、ccamはその大きさを表示します。移動で[`--brief`](#brief)が使われることはありません。

`auto_switch`は、Claude Codeが上限を報告する前に切り替えることはなく、バックグラウンドで使用量を監視することもなく、ログインをコピー、更新、受け渡しすることも決してありません。Windowsでは無視され、`ccam doctor`がそのことを伝えます。Windowsでも、`relay = true`なら`/ccam:switch`と`/exit`で会話を別のアカウントに移せます。

**Anthropicの規約。** TeamとEnterpriseのシートでは、`auto_switch`はあなたの作業を同じ組織の別のシートに移します。Anthropicのヘルプセンターには、「Each seat on a Team plan is assigned to one named individual, and only that person may use it」（Teamプランの各シートは1人の指定ユーザーに割り当てられ、そのユーザーのみが使用できます）、「Users can't route their work through another user's account to use that user's remaining usage」（ユーザーは別のユーザーのアカウントを通じて作業をルーティングして、そのユーザーの残りの使用量を使用することはできません）と書かれています（[Teamプラン](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)。[Enterpriseプラン](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)も同様です）。

個人プラン（Pro、Max）については、上限に達したときに自分のアカウント同士が自動で引き継ぎ合ってよいかどうかを、Anthropicの規約は述べていません。[消費者向け利用規約](https://www.anthropic.com/legal/consumer-terms)は、Anthropicが明示的に許可しない限り、「through automated or non-human means, whether through a bot, script, or otherwise」（ボット、スクリプト、その他の方法による、自動化された手段または人間以外の手段によって）サービスにアクセスすることを禁じており、移動の際には上の決まったメッセージのいずれかがあなたの代わりに送られます。2026年11月12日に発効する[利用ポリシー](https://www.anthropic.com/legal/aup)は、禁止事項の1つに「Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products」（当社の製品内に設けられた機能、制限、ガードレールを意図的に回避または妨害すること）を挙げています。

ccamは、`auto_switch`がこれらの規約に適合しているとは主張しません。有効にするかどうかはあなたの判断であり、そのリスクもあなたが負います。Anthropicは、規約に違反していると判断した場合、アカウントや、TeamまたはEnterpriseの組織のアクセスを停止することがあります。

`auto_switch`が有効な間、`ccam doctor`は、ccamプラグインのないアカウントと2.1.214より古いClaude Codeについて警告し、シートについてはAnthropicのシートに関する規則を繰り返し示します（[点検](/features/doctor)）。

## 上限に達したら再実行

`ccam next`は、止まった対話型セッションのためのものです。非対話型の実行、つまりスクリプトや定期実行ジョブの中の`claude -p`はあなたを待てないので、アカウント名の前に`--fallback`を付けると、ccamがその実行のために切り替えます。

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

1. `--fallback`は`-p`か`--print`と一緒の場合にだけ動作します。そうでない場合、ccamは停止して`ccam next --auto --continue`を案内します。
2. パイプからの標準入力は、最大10 MiBまで一度だけ読み込まれ、すべての試行に渡されます。Claude Codeと同じく、ccamは入力が始まるのを最大3秒待ちます（`--input-format stream-json`の場合を除きます）。パイプが開いたまま何も送らない場合、実行は入力なしになり、ccamは標準エラー出力でそのことを知らせます。
3. ccamは`exec`ではなく子プロセスとして`claude`を起動し、標準出力を保留して、標準エラー出力は届いたとおりに流します。
4. `claude`が0で終了した場合や、別の理由で失敗した場合は、保留した出力と終了コードをそのまま返します。シグナルで終了した実行（128より大きい終了コード）も同じです。停止を求められたものなので、再実行はしません。
5. Claude Codeのメッセージが使用量の上限を示している場合（[フック](#inside-claude-code)が知っている使用量の上限のメッセージのいずれかで、`--output-format json`や`stream-json`の出力の中にあるものも含みます）や、アカウントの使用量で5時間、週間、または`--model`の枠が100%になっている場合（上限の記録も含みます）、その失敗は使用量の上限と見なされます。
6. その場合、その試行の出力は破棄され、同じコマンドが別のログイン済みのアカウントで実行されます。順序は、`--auto`が選ぶアカウント、次に使用量を問い合わせられないアカウントを番号順です。100%のアカウントはスキップされ、APIキーでログインしたアカウントは試されず、各アカウントは最大1回だけ試されます。
7. 引数に`-c`、`--continue`、`-r`、`--resume`がある場合は、会話を見つけられるように、同じ履歴グループのアカウントだけを試します。
8. `--auto`と同じく、同じ組織のアカウントだけを試します（`cross_org = true`でこの制限を外せます）。

次のアカウントは`--auto`と同じ規則（`auto_max_usage`、`auto_strategy`、`--model`の枠）で選ばれます。試せるアカウントが残っていない場合、ccamは最後の試行の出力を表示し、余裕のあるログイン済みのアカウントがないことを伝えて、終了コード75で終了します。アカウント名の前に`--wait`を付けると（`ccam --fallback --wait=2h work -p …`）、代わりに最初のアカウントがリセットされるまで待ち、そのあとすべてのアカウントをもう一度試すことがあります。必ず先に一度は待つので、上限に達したばかりのアカウントがすぐに再実行されることはありません。指定した時間は、待機1回ごとではなく、その実行のすべての待機を合わせた上限です。スペースで区切って書くと（`--wait 30m`）、`ccam next`と同じヒントとともに拒否されます。Claude Codeの引数の前の`--`は、ほかの起動と同じように使えます（`ccam --fallback work -- -p …`）。`ccam --dry-run --fallback work -p …`は、最初の試行のコマンドを表示します。

[`auto_switch`](#auto-switch)を有効にしない限り、対話型セッションが自動で切り替わることはありません。`ccam statusline install`を実行しておくと、5時間または週間の使用量が90%に達したときに、ステータスラインに`→ ccam next --auto --continue`が追加されます（[使用量の上限](/features/usage)）。

## 安全性と特殊なケース

- ログインしていないアカウントが選ばれることはありません。ほかにログイン済みのアカウントがない場合、オプションなしの`ccam next`は`no logged-in account other than <name>`で停止します。
- `--continue`を付けて、フォルダーに会話がない場合: `no account has a conversation in <folder> to continue`
- 組織の規則で候補のアカウントがなくなった場合: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``
- `--auto`を付けて、ログイン済みのアカウントがない場合: `no logged-in account to pick`。すべてのアカウントが上限に達しているか使用量が不明な場合: `` no account with known usage and headroom (see `ccam usage`) ``。どちらも終了コード75で終了します。
- `--max-usage`、`--strategy`、`--wait`、`--cross-org`は`--auto`と一緒にしか使えず、1-100以外の値や不明な戦略は拒否されます。
- フォークでコピーされるのは会話であり、アカウントのほかの状態はコピーされません。履歴を共有していない限り、元のアカウントの`file-history/`にあるチェックポイントは引き継がれません。
- `--auto`は`--dry-run`を付けてもネットワークリクエストを行います。読み取るだけで、トークンファイルに書き込むことはありません。
- 選ばれたアカウントは`ccam <name>`と同じ経路で起動するので、同じ確認が行われます。同じ履歴を共有する別のClaude Codeセッションがすでにこのフォルダーで実行中の場合の警告も含まれます。
- `last-used`は、Claude Codeが実際に起動したときにだけ書き込まれます。`--dry-run`はこれを変更しないので、ドライランでローテーションが進むことはありません。
- `--fallback`は試行が終わるまで標準出力を保留するので、`--output-format stream-json`のようなストリーミング出力は最後にまとめて届きます。`stream-json`の入力をストリーミングし、応答を待ってから続きを送るクライアントは、`--fallback`を使えません。最初の試行が上限に達する前に変更したファイルは変更されたままで、次の試行はプロンプトを最初からやり直します。

## 関連項目

- [使用量の上限](/features/usage)
- [使用量の上限とアカウントの切り替え](/docs/usage-limits)
- [MCPサーバー、設定、履歴の共有](/docs/sharing)
- [Claude Codeが使用量の上限に達したときの対処法](/guides/usage-limit-workflow)
