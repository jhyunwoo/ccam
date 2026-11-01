---
title: Claude Codeが使用量の上限に達したときの対処法
description: リセットを待つか、別のClaude Codeアカウントですぐに同じ会話を続けます。5時間の上限と週間の上限のしくみと対処の手順を説明します。
updated: 2026-11-01
order: 3
translation: machine
sourceHash: 760a36724feb
---

アカウントが使用量の上限（レート制限とも呼ばれます）に達したときの選択肢は2つです。上限がリセットされるまで待つか、自分が所有する別のアカウントで同じ会話をすぐに続けるかです。ccamを使う場合、後者は`ccam next --auto --continue`という1つのコマンドで済みます。

## Claude Codeにはどんな上限がありますか？

ProプランとMaxプランでは、Claude Codeに[5時間のセッション上限と週間の上限](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan)があります。どちらもWebやアプリのClaudeと共通で、使用量は両方の上限に同時に加算されます。どちらかを使い切ると、Claude Codeは「You've hit your session limit · resets 3:45pm」のような[メッセージを表示します](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit)。上限はすべてのモデルに適用されるため、モデルを切り替えても解決しません。

## 方法1：リセットを待つ

アカウントが1つなら、待つことになります。最近のバージョンのClaude Codeは、[開いているセッションのまま待機し](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset)、上限がリセットされると自動的にタスクを再開します。Claude Codeの中で`/usage`を実行すると、プランの上限をどれだけ使ったかを確認できます。

## 方法2：今すぐ別のアカウントで続ける

2つ目のアカウントがあれば、同じ会話をすぐに続けられます。一度だけ設定しておきます。

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

これで、作業中はステータスラインにアカウントと5時間の使用量が表示されるので、上限に達して驚くことはありません。

上限に達したら、Claude Codeを終了（`/exit`）してから、次のコマンドを実行します。

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

最も余裕があるアカウントで会話を開きます。

```sh
ccam next --auto --continue
```

Claude Codeが`work`として起動し、現在のフォルダーで最も新しい会話を、それまでのやり取りをすべて含めて開きます。履歴を共有していれば（上の`ccam group set default --add history`）、両方のアカウントが同じ会話で作業します。共有していなければ、Claude Codeが会話を`work`にコピーし、ccamは最初の応答で読み直す必要があるトークン数を表示します。

Claude Codeに[ccamプラグイン](/docs/install#claude-code-plugin)を入れておくと、上限で止まったセッションが通知を表示し、満杯の時間枠を記録するので、`ccam next --auto`がリセット前にそのアカウントを再び選ぶことはありません。`/ccam:switch`（または`/ccam:switch work`）は、`ccam next --continue`のために会話をメモに残すので、まさにその会話が続きます。

**TeamとEnterpriseのシート。** 自動的なアカウント選択（`ccam next`、`--fallback`、`ccam pick`、ダッシュボードの`n`、有効にした場合の[`auto_switch`](/features/next#auto-switch)）は、同じClaude TeamまたはEnterpriseの組織に属するシート間でも切り替えます。これはAnthropicの規約に違反するおそれがあります。頼る前に[FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)を読み、所属する組織とAnthropicに確認してください。

## ccamは自動でアカウントを切り替えますか？

有効にした場合だけです。

- デフォルトでは、作業中にccamは動いていません。ccamはClaude Codeを起動し、自身をClaude Codeに置き換えるので、切り替えるタイミングはあなたが決めます。`ccam next --auto`は、5時間の使用量が最も少ないアカウントを選び、100%に達したアカウントとログインしていないアカウントは飛ばし、1つのTeamまたはEnterpriseの組織の中、あるいは個人プランの中に留まります。
- `config.toml`で`relay = true`を設定すると、ccamは代わりにセッションと並んで動き続けるので、`/ccam:switch`だけで、セッションのオプション、権限モード、effortレベル、エージェントを保ったまま、会話が自動で続きます。会話を別の組織のアカウントへ移す前には、ccamがあなたに確認します。
- `auto_switch = true`を設定すると（LinuxとmacOS。実験的な機能）、ccamは切り替えそのものも自動で行います。Claude Codeがセッションを上限で止め、同じ種類の別のアカウントに余裕があるとき、ccamはセッションを終了させ、そのアカウントで会話を続けます（[条件とコスト](/features/next#auto-switch)）。

実行中のセッションのログインを入れ替えるツールもあります。違いは[ccamとほかのアカウント切り替えツール](/compare/account-switchers)で説明しています。

## 注意点

- `ccam next --continue`は、作業していたのと同じフォルダーで実行してください。会話はプロジェクトのフォルダーごとに保存されます。
- 2つのアカウントで同じ会話を同時に使うのは避けてください。そうしようとすると、ccamが警告します。
- 履歴を共有していない場合、会話はもう一方のアカウントにコピーされ、2つのコピーはそれぞれ別々に進みます。`/rewind`がコードを復元するためのチェックポイントは、最初のアカウントに残ります。
- 新しいアカウントではまだ会話がキャッシュされていないため、最初のメッセージで通常より多く上限を消費することがあります。`ccam next --continue --brief`なら、代わりに短い要約から始めます。
- このようにアカウントを切り替えてよいかどうかは、それぞれのアカウントに適用される規約によります。

すべてのオプションについては、[使用量の上限](/features/usage)と[次のアカウントへの切り替え](/features/next)を参照してください。
