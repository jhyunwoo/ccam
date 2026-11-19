---
title: 別のアカウントで動くサブエージェント
description: ccam agentは、Claude Codeのタスクの一部をあなたの別のアカウントで実行してClaudeの返答を返すので、1つの長いタスクで複数のアカウントの使用量を使えます。
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Claude Code自身のサブエージェントは、それを起動したセッションのアカウントで動くので、そのアカウントの使用量を消費します。`ccam agent "<task>"`は、同じフォルダーで、あなたの別のアカウントとして`claude -p`でタスクを実行し、Claudeの最後の返答を表示します。Claude Codeのセッションから呼び出すと、別のアカウントの使用量を消費するサブエージェントのように働くので、セッション自身のアカウントなら使い切っていたところを超えて、1つの長いタスクを続けられます。

**TeamとEnterpriseのシート。** `ccam agent`はアカウントを自分で選びます。1つのClaude TeamまたはEnterpriseの組織の中では、タスクを別のシートで実行します。これはAnthropicの規約に違反するおそれがあります。頼る前に[FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)を読み、所属する組織とAnthropicに確認してください。

## 使い方

Claude Codeは、これをBashツールで実行します。[ccamプラグイン](/docs/install#claude-code-plugin)をインストールしていれば、「このタスクの各部分に私の別のアカウントを使って」のように普通の言葉で頼むか、それをCLAUDE.mdに書いておきます。プラグインのスキルはClaude Codeに、各エージェントをバックグラウンドで起動すること、それだけで完結するタスクを書くこと、ファイルを編集するエージェントにはそれぞれ専用の担当部分かgit worktreeを与えることを指示します。

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

Claudeの最後の返答は標準出力に、`ccam:`で始まる行は標準エラー出力に出ます。`--json`では、代わりに`{account, session_id, result, is_error, accounts}`を表示します。`accounts`は、そのタスクを実行したすべてのアカウントの一覧です。

終了コード: claude自身の終了コード（エージェントが終了したときは`0`）、呼び出し元以外に余裕のあるアカウントがないときは`75`（`EX_TEMPFAIL`）、エージェントが別のエージェントを起動しようとしたときなど、そのほかのエラーでは`1`です。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--model <model>` | エージェントのモデル。選択では、そのモデルの週間の枠も数えます | アカウントのデフォルトのモデル |
| `--resume <id>` | 以前のエージェントの会話を、続きのタスクで続けます（セッションIDは、その最初と最後の行にあります） | オフ |
| `--json` | アカウント、セッション、返答、実行に使ったアカウントを表示します | オフ |
| `--max-usage <1-100>` | 使用量がこの値未満のアカウントを優先します | `auto_max_usage`、なければ100 |
| `--strategy headroom\|expiring\|balance` | [選択の規則](/features/next#auto) | `auto_strategy`、なければ`headroom` |
| `--cross-org` | ほかの組織のアカウントでも実行します | `cross_org`、なければオフ |
| `--wait[=<duration>]` | 余裕のあるアカウントがないとき、最初のアカウントがリセットされるまで最長でこの時間だけ待ちます | オフ。値なしで指定すると`5h` |
| `--dry-run` | 選ぶはずのアカウントで実行するコマンドを、実行せずに表示します | オフ |

## しくみ

1. **呼び出し元は除外されます。** `ccam agent`を呼び出したアカウント（ccamが起動したセッションのアカウント、なければこのフォルダーで起動されるアカウント）がエージェントを実行することはありません。同じClaudeアカウントにログインしている別のフォルダーも、上限を共有しているので同じく除外されます。`--cross-org`を指定しない限り、エージェントは呼び出し元の[組織](/features/next#auto)の中にとどまります。
2. **`ccam pick`と同じように選びます。** 使用量は[`ccam pick`](/features/pick)と同じ方法で問い合わせます。2分以内に相次いで起動したエージェントには、余裕のあるアカウントが十分にある間は、別々のアカウントを割り当てます。
3. **セッションの権限で動き、それを超えることはありません。** `claude -p`ではツールを許可する人がいないので、エージェントは呼び出したセッションの権限モードとeffortを受け継ぎます。ターミナルから実行した場合は、Claude Codeのデフォルトのモードと、設定の権限ルールが適用されます。`ccam agent`に、それ以上の権限を与えるオプションはありません。
4. **使用量の上限に達すると、会話が引き継がれます。** エージェントのアカウントが上限に達すると、ccamは[プラグインのフック](/features/usage#the-limit-record)と同じように上限を記録し、最も余裕があるアカウントとして同じ会話を続けます。Claude Codeが会話をそこへコピーし（`--resume <file> --fork-session`）、最初のメッセージは「Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.」です。各アカウントを試すのは1回だけです。
5. **エージェントの会話はすべて記録されます。** 記録先は`~/.local/state/ccam/agent/<session id>.json`です（30日間）。`--resume <id>`は、その会話をそのフォルダーで続けます。実行したアカウントにまだ余裕があればそのアカウントで、なければ別のアカウントで続けます。[`ccam next --continue`](/features/next)はエージェントの会話を対象外にするので、引き続きあなた自身の会話を続けます。

## 安全性と特殊なケース

- エージェントが自分のエージェントを起動することはできません。ccamはエージェントに`CCAM_AGENT=1`という印を付け、そこでは`ccam agent`が実行を拒否します。
- `ccam agent`を止めると（Ctrl+C、シグナル、またはClaude Codeによるコマンドの中断）、エージェントも止まります。
- エージェントには、自分を呼び出した会話は見えず、自分のタスクだけが見えます。
- 同じファイルを同時に編集するエージェントは、互いの作業を上書きしてしまいます。それぞれに、タスクの専用の担当部分か、専用のgit worktreeを与えてください。
- 呼び出し元のセッションが別のアカウントに移ると（[`auto_switch`](/features/next#auto-switch)）、そのセッションが起動したエージェントも一緒に止まります。その会話は残るので、`--resume <id>`で続けられます。
- `agent`はコマンド名なので、新しいアカウントに`agent`という名前は付けられません。既存のアカウントは、引き続き`ccam run agent`で起動できます。

## 関連項目

- [スクリプト用のアカウント選び](/features/pick)
- [次のアカウントへの切り替え](/features/next)
- [使用量の上限](/features/usage)
- [すべてのアカウントのバックグラウンドセッション](/features/agents)
