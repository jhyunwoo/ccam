---
title: 使用量の上限
description: ccam usageはすべてのアカウントの5時間と週間の上限を一度に表示し、ccam statuslineはClaude Codeにアカウント名と使用量を表示します。
updated: 2026-11-01
order: 7
translation: machine
sourceHash: 54b296f7da35
---

`ccam usage`は、すべてのアカウントについて、各上限をどれだけ使ったかを表示します。対象は5時間の上限、週間の上限、モデルごとの週間上限（Fableなど）です。上限ごとに、バーとリセットまでの時間が表示されます。長い作業を始める前や、Claude Codeが上限に近いことを知らせたときに使い、どのアカウントに余裕が残っているかを確認します。`ccam statusline`はClaude Code用のステータスラインのコマンドです。Claude Codeの中にアカウント名と使用量を表示し、その数値を`ccam usage`のために保存します。使用量のもとになるトークンと推定コストについては、`ccam ccusage`がすべてのアカウントの会話に対して[ccusage](https://ccusage.com)を実行します（[後述](#tokens-and-costs-ccam-ccusage)）。

## 使い方

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

各アカウントには、番号、名前、数値がどれだけ前のものかを示す行があります。`(cached)`はエンドポイントではなく保存済みのファイルから取得した数値であることを意味します。その下には上限ごとに1行ずつ、20マスのバー、割合、枠がリセットされるまでの時間が表示されます。`5-hour`、`Weekly`（全モデル）、そして`Fable`のようなモデルごとの週間の枠の順です。アカウントにない枠は`-`と表示されます。割合は80%で黄色、100%で赤になります。

Claude Codeが時間枠が満杯だという理由でリクエストを拒否した場合、その時間枠は、保存済みの数値にかかわらず、リセットされるまで100%と表示されます（[上限の記録](#the-limit-record)を参照）。ターミナルでは、すべてのアカウントがまず`fetching…`として一度に表示され、応答があり次第その行が描画されます。

最後の行は、アカウントが2つ以上ある場合だけ表示され、[`ccam next --auto`](/features/next)が選ぶアカウントを示します。同じ規則に従います。

- `auto_max_usage`と`auto_strategy`
- 実行中のセッション
- `cross_org`を設定しない限り、最後に使ったアカウントと同じ組織のアカウントだけ
- 100%の枠はすべてスキップ。そのため、週間の枠が満杯のアカウントが示されることはありません

デフォルトの`headroom`戦略では`→ Most headroom: <name>`と表示され、`expiring`や`balance`では、`→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`のように理由が示されます。

**TeamとEnterpriseのシート。** 自動的なアカウント選択（`ccam next`、`--fallback`、`ccam pick`、ダッシュボードの`n`、有効にした場合の[`auto_switch`](/features/next#auto-switch)）は、同じClaude TeamまたはEnterpriseの組織に属するシート間でも切り替えます。これはAnthropicの規約に違反するおそれがあります。頼る前に[FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)を読み、所属する組織とAnthropicに確認してください。

ログインしていないアカウントには`not logged in → ccam login <name>`、APIキーのアカウントには`no limits (API)`と表示されます。アクセストークンの期限が切れている場合は、保存済みの数値があればその上に``access token expired; `ccam <name>` renews it when it starts``と表示されます。次の起動時にClaude Codeが更新します。リフレッシュトークンの期限が切れたか、サーバーに拒否された後にClaude Codeが消去したために、ログインの期限が完全に切れている場合は`` login expired → `ccam login <name>` ``と表示されます。

`ccam statusline install`の後は、Claude Codeに`[work] · Opus · 5h 24% · 7d 41%`のような行が表示されます。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--json` | 結果をJSONで出力します。`name`、`slot`、`status`、`limits`（アカウントで使用量クレジットがオンの場合は`credits`付き）、`source`、`at`、`error`を含み、上限の記録が時間枠を100%にしている場合は`limited`と`limited_windows`も含みます。アカウントがない場合は`[]`を出力します | オフ |
| `--watch[=N]` | Ctrl+Cを押すまで`N`秒ごとに描き直します。`--watch`だけの場合や、0以下の値の場合は60です。値は`=`の後に書きます。`--watch 10`はヒントとともに拒否されます | オフ |
| `-v`、`--verbose` | 各アカウントの下に`source:`の行を追加し、数値の取得元（使用量エンドポイント、Claude Code、ステータスライン、保存済みのファイル、またはClaude Codeが報告した使用量の上限）を示します。上限の記録が数値より優先される場合は、`Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset`のような行が続きます | オフ |
| `--no-cache` | 60秒のキャッシュを無視し、エンドポイントに問い合わせ直します | オフ |
| `--source live\|claude\|cache` | `live`はアカウントのトークンでエンドポイントに問い合わせます。`claude`はClaude Codeに問い合わせさせます（後述）。`cache`は保存済みの数値だけを表示します | `usage_source`、なければ`live` |
| `--refresh` | 期限切れのアクセストークンを更新します。対象は実行中でないアカウントだけです。更新には最大30秒かかることがあります | オフ |
| `--dry-run` | 受け付けます。`ccam usage`は自身のキャッシュ以外を変更しません。`--refresh`と一緒に使うと`--dry-run cannot be combined with --refresh, which rewrites login tokens`で拒否されます | オフ |

`--source`にそれ以外の値を指定すると、`--source must be live, claude or cache`で失敗します。`ccam usage`は引数を取りません。

### Claude Code自身による使用量の取得

`config.toml`の先頭（最初の`[[group]]`、`[[profile]]`、`[[folder]]`より上）で`usage_source = "claude"`を設定すると（または`--source claude`を指定すると）、ccamはログイントークンを一切読みません。アカウントごとに、アカウントのフォルダーで`claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence`を起動し、Claude Codeの`get_usage`リクエストを標準入力に送ります（Claude Code 2.1.270以降）。Claude Codeは自身のログインで使用量を問い合わせ、5時間、週間、モデルごとの枠を返します。

`--safe-mode`によってアカウントのフック、MCPサーバー、プラグインは起動せず、セッションも保存されません。アカウントごとにClaude Codeの起動（約1秒）が必要なので、この取得元は`live`より遅くなります。アカウントへの問い合わせは同時に行われ、それぞれ20秒でタイムアウトし、60秒のキャッシュも通常どおり適用されます。

アクセストークンの期限が切れているか、5分以内に切れるログインには50秒を与えます。Claude Codeは起動時にそれを更新し、途中で打ち切られた更新ではログインが失われることがあるためです。ccamが問い合わせをあきらめた場合は、SIGTERMでClaude Codeに終了を求めるので、Claude Codeは始めた更新を終えられます（WindowsにはSIGTERMがないため、そこではプロセスを強制終了します）。

ログインの期限が完全に切れている場合は、Claude Codeを起動せずに`` login expired → `ccam login <name>` ``と表示します。Claude Codeはこのリクエストを実験的なものとしているため、応答しなくなった場合、ccamは保存済みの数値を表示し、`ccam doctor`がそれを報告します。Claude Codeが自身の問い合わせに失敗したと応答した場合、ccamは上限のないプランとしてではなく、`the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`とともに保存済みの数値を表示します。

`--dry-run`では、ccamはClaude Codeを起動せず（Claude Codeがアカウントのフォルダーに書き込むため）、代わりに保存済みの数値を表示します。`usage_source`には`live`か`claude`を指定します。`cache`は`--source`でだけ選べます。

`ccam statusline`は、Claude Codeのステータスライン用のJSONを標準入力から読みます。オプションは`--wrap <command>`の1つだけで、既存のステータスラインのコマンドを指定すると、その出力の1行目の前にアカウント名が付きます。`ccam statusline install`はオプションを取りません。どちらも`--dry-run`を`` --dry-run is not supported by `ccam statusline`; nothing was changed ``で拒否します。

## トークンとコスト：`ccam ccusage`

`ccam usage`は、アカウントが各上限をどれだけ使ったかを表示します。これに対して[ccusage](https://ccusage.com)は、Claude Codeの会話ログにあるトークンと推定コストを集計します。ccusageは単独では、Claude Codeのログを`~/.claude`と`~/.config/claude`（または`CLAUDE_CONFIG_DIR`で指定したフォルダー）でしか探さないため、見えるccamのアカウントは、`~/.claude`にその場所のままインポートしたものだけです。`ccam ccusage`は、会話があるすべてのアカウントのフォルダーを`CLAUDE_CONFIG_DIR`に設定してccusageを実行するので、1つのレポートですべてのアカウントを扱えます。

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccamはまず自身のオプションを読みます。`-a`（`--account`）はアカウント名または番号を受け取り、カンマ区切りでも、オプションを繰り返しても指定できます。`--dry-run`はコマンドを実行せずに表示します。ccamのオプションでない最初の引数から後は、`bun x`と`npx`での`--config`を除き（後述）、すべて入力したとおりにccusageに渡ります。`ccam ccusage -- --help`ならccusage自身のヘルプを表示します。ccamは、`ccusage`がインストールされていればそれを、なければ`bun x ccusage`を、それもなければ`npx ccusage@latest`を実行します。後の2つはccusageをnpmから取得します。この2つは今いるフォルダーではなくccam自身のフォルダー（`~/.local/state/ccam/ccusage`）で起動します。そのため、今いるフォルダーのプロジェクトが、その`node_modules/.bin`、`package.json`の`bin`、`.npmrc`のレジストリによって、実行されるものを選ぶことはできません。`ccam --dry-run ccusage`では、そのフォルダーへの`cd`が表示されます。どれも使えない場合は、Node.js、Bun、ccusageのいずれかをインストールするよう伝えます。LinuxとmacOSではccamが自身をccusageに置き換え、Windowsではccamがccusageの終了を待ってその終了コードで終了します。WSLでは、Windowsのドライブ上にある`ccusage`、`bun`、`npx`は使わずに飛ばします。それらはアカウントのフォルダーのないWindowsのプログラムとして動いてしまうためです。

- レポートの対象は、フォルダーに会話がある（`projects`フォルダーがある）アカウントです。複数のアカウントがグループを通じて共有している会話は、1回だけ数えます。
- `-a`を使うと、グループを通じて履歴を共有しているアカウントからはグループ全体の会話が入り、ccamはそれが誰の会話かを伝えます: `ccam: work shares its conversation history with personal, so the report counts their conversations too`
- `blocks`は5時間の枠を表示します。この枠はアカウントごとに別々なので、上の例のようにアカウントを1つだけ指定します。
- ccusageは、ツールを1つ指定しない限り、見つけたほかのコーディングツール（Codexなど）も報告します。`ccam ccusage claude daily`ならClaude Codeだけです。
- ccamから実行したとき、ccusageは設定（`ccusage.json`）を`./.ccusage`とアカウントのフォルダーからだけ読みます。別のファイルは`--config`で指定します: `ccam ccusage daily --config ~/.config/claude/ccusage.json`。`bun x`と`npx`はccam自身のフォルダーで起動するため、`./.ccusage/ccusage.json`があり、ファイルを指定していないときは、ccamがそれを`--config`で渡し、相対パスの`--config`はフルパスにして渡します。このときccusageは`--pi-path`などほかの相対パスをccamのフォルダーから読むので、それらはフルパスで指定してください。
- ccusageはパスにカンマを含むフォルダーを2つのフォルダーとして読んでしまうため、ccamはそのフォルダーを拒否し、アカウント名を示します。そのアカウントは`-a`で外してください。
- `ccusage`はccamのコマンドなので、ccusageという名前のアカウントは`ccam run ccusage`で起動します。`ccam doctor`もそのことを伝えます。

## しくみ

`ccam usage`の数値の取得元は2つあります。各アカウント自身のログインで呼び出す、Anthropicの公開ドキュメントに記載されていない使用量エンドポイント`GET https://api.anthropic.com/api/oauth/usage`と、Claude Codeがステータスラインのコマンドに渡す値です。これらはすべて`internal/usage`にあります。

1. ccamは`config.toml`を読み込み、アカウントごとに1つの`usage.Account`を作ります（`cli.usageCmd`）。内容は、フォルダー`~/.local/share/ccam/profiles/<name>/`、セッションのマーカー（アカウントのフォルダーにあるClaude Code自身の記録と合わせて、実行中かどうかがわかります）、その場所のままインポートしたかどうかです。
2. ccamはすべてのアカウントを一度に問い合わせ、アカウントごとに1つのgoroutineを使います（`Fetcher.FetchAll`）。すべてが1つの`http.Client`を共有するので、エンドポイントへの接続は再利用されます。接続には最大2秒、各アカウントには合計3秒かかってもかまいません。リダイレクトするとトークンが別の場所に再送されてしまうため、クライアントはリダイレクトをたどらず、応答は最大1 MBまでしか読みません。
3. 各アカウントの処理が終わるたびに、その行が描き直されます。失敗したアカウントや時間内に応答しなかったアカウントは保存済みの数値にフォールバックするので、1つの遅いアカウントがほかのアカウントを待たせることはありません。

図では、各アカウントがそれぞれのgoroutineで問い合わせられ、タイムアウトしたアカウントは代わりに保存済みの数値を表示します。

<!-- motion:usage-fetch -->

`personal`と`work`は時間内に応答し、最新の数値を表示します。`side`は3秒より長くかかるので、ccamは代わりに保存済みの数値を表示します（`fallback()`）。続いてccamは、アカウントごとのブロックを描画し（`usage.Table()`）、最後の行に示すアカウントを選びます（`usage.Best()`）。

1つのアカウントについて、ccamは次の順序で確認します（`Fetcher.Fetch`）。

<!-- motion:usage-cache -->

1. `--source cache`の場合は、直接`fallback()`に進みます。
2. `--no-cache`が指定されていなければ、`~/.local/state/ccam/usage/<name>.json`にある、計測してから60秒（`CacheTTL`）未満の最新の結果を、現在の時点で判断して使います。その後リセットされた時間枠は0%と表示し、時計より先の日付のファイルは新しいものとみなしません。`ccam usage`を続けて2回実行しても、リクエストはアカウントごとに1回で済みます。
3. ccamは`.credentials.json`（またはmacOSのキーチェーン）を変更せずに読み、`claudeAiOauth`キーだけをデコードします（`readCredentials()`）。サブスクリプションのログインがない場合は`not logged in`、`.claude.json`にAPIキーがある場合は`no limits (API)`になります。リフレッシュトークンの期限が切れたログインや、Claude Codeが消去したログインは`` login expired → `ccam login <name>` ``となり、問い合わせは行いません。
4. 期限切れのアクセストークンは使いません。`--refresh`がない場合、そのアカウントは``access token expired; `ccam <name>` renews it when it starts``とともに保存済みの数値を表示します。
5. それ以外の場合、ccamはエンドポイントを呼び出します。401や403はログインの期限切れ（`login expired: authentication failed (401)`）、それ以外のステータスや読めない本文はエラーとして扱い、どちらの場合も保存済みの数値にフォールバックします。Ctrl+Cで終わった問い合わせには`cancelled`と表示します。429は、エンドポイントがこのアカウントについて受け取ったリクエストが多すぎることを意味します。ccamは`<name>.pause`に次に問い合わせる時刻（`Retry-After`、なければ5分、最長1時間）を書き込み、保存済みの数値を`rate limited by the usage endpoint; asking again after 15:10`とともに表示し、その時刻までは`--no-cache`を付けてもそのアカウントのエンドポイントを呼び出しません。すぐに再試行しても、レート制限が長引くだけだからです。
6. 正常な応答は、計測した時刻とともに`<name>.json`に書き込まれます。

保存済みの数値には、ccamは最も新しい保存済みの問い合わせ結果（`<name>.json`、またはClaude Codeによる問い合わせの`<name>.claude.json`）を使い、`<name>.statusline.json`がその後に書き込まれていれば、ステータスラインの5時間と週間の枠をその上に重ねます（`fallback()`）。ステータスラインにはないモデルごとの枠は、問い合わせ結果に保存されたままです。保存済みの数値を表示する前に、各時間枠のリセット時刻を確認します。数値の保存後にリセットされた時間枠は、古い割合がもう当てはまらないため0%と表示します。上の例で`personal`がリセット時刻なしで`0%`と表示されているのはこのためです。

エンドポイントの応答形式は2つあります。新しい形式には`limits`配列があり、`session`（5時間）、`weekly_all`、`weekly_scoped`（モデル別）の項目を含みます。古い形式には`five_hour`、`seven_day`、`seven_day_opus`、`seven_day_sonnet`のオブジェクトがあります。ccamは両方を読みます（`ParseLive`）。`limits`配列が優先され、それに含まれないものを古いオブジェクトで補います。上限がまったく含まれない本文は形式の変更として報告され、`ccam doctor`がそれを知らせます。

### 上限の記録

Claude Codeが使用量の上限でターンを止めると、[ccamプラグイン](/docs/install#claude-code-plugin)のフックが、拒否されたリクエストを会話から読み取ります。どの時間枠が満杯になったか（5時間、週間、Opus、Sonnet）と、それがいつリセットされるかです。フックはそれを`~/.local/state/ccam/usage/<name>.limit.json`に書き込みます。ccamは、あらゆる取得元から得たほかの情報の上にこの記録を適用します（`Fetch`）。記録にある各時間枠は、保存済みの数値がどれほど古くても低くても、そのリセットまで100%と数えるので、`ccam next --auto`、`ccam pick`、`--fallback`、[`auto_switch`](/features/next#auto-switch)がそれまでにそのアカウントを再び選ぶことはありません。

ほかに何もわからない場合は、記録だけがそのアカウントの表示になります（`-v`では`source: the usage limit Claude Code reported`）。無人のセッション（`-p`の実行とバックグラウンドセッション）は記録を書き込まず、組織が設定した停止（たとえば使用量が尽きた場合や、グループの上限が$0の場合）も、管理者がいつでも解除できるので記録しません。別のClaudeアカウントへの`ccam login`で記録は破棄され、ファイルを削除しても消えます。

### トークンの更新

リフレッシュトークンはローテーションします。更新すると新しいリフレッシュトークンが発行され、古いものは無効になります。同じアカウントでClaude Codeが実行中のときにccamがログインを更新すると、そのセッションはサインアウトされてしまいます。そのため、`ccam usage`はトークンを読むだけで、`--refresh`を付けない限り書き込みません。このオプションを付けると、ccamは次の処理を行います（`Fetcher.refresh`）。

1. `ccam import`でその場所のまま登録したアカウントは拒否します。ccamの知らないところで通常の`claude`がそのログインを使っている可能性があるためです。
2. macOSのキーチェーンに保存されたログインは拒否します。
3. ロックファイル`~/.local/state/ccam/locks/profile-<name>.refresh.lock`を取得し、2つのccamのプロセスが同じトークンを両方とも更新することがないようにします。ロックが保持されている間、そのアカウントには`token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)`と表示されます。`ccam rename`、`remove`、`logout`もこのロックを取り、更新の実行中は拒否します。
4. ロックを保持した状態で、アカウントのセッションが、どう起動されたかにかかわらず（ccam経由、通常の`claude`、IDE、またはバックグラウンドセッションとその背後にあり、自分でログインを更新するデーモン）実行中なら拒否し、別のプロセスがすでにトークンを更新していればリクエストをスキップします。
5. Claude Codeがアカウントのフォルダーで自身の更新用のロック（`.oauth_refresh.lock`、`<folder>.lock`）を保持している間は、何も送りません。Claude Codeがちょうどそのログインを更新しているためです。
6. リフレッシュトークンをトークンエンドポイントにPOSTし、Claude Codeが自身の更新に与えるのと同じ最大30秒、応答を待ちます。一度送ったリクエストは、Ctrl+Cでも中断しません。サーバーがすでにリフレッシュトークンを入れ替えていて、その応答だけが新しいログインを持っている可能性があるためです。Ctrl+Cやハングアップでほかの問い合わせは終わり、応答が保存された時点でccamは終了します（Ctrl+Cの後は130）。
7. Claude Codeの書き込みロック（`.storage-write.lock`）を取って保存し、その時点の`.credentials.json`のトークンの項目だけを変更します。その間にClaude Codeが書き込んだもの（MCPのサインイン）は残り、その間にClaude Codeが更新したログインはそのまま使われ、ログアウトが取り消されることもありません。応答に新しいリフレッシュトークンの有効期間が含まれている場合は、Claude Codeと同じく`refreshTokenExpiresAt`として保存します。

更新が行われなかった場合は、アカウントの行か保存済みの数値の上にその理由が示されます。たとえば`token not renewed: a session is running`です。リクエストを送ったものの、その応答を使えなかったり保存できなかったりした場合は`login may be lost: …`と表示され、その後Claude Codeにログインを求められたら`ccam login <name>`を実行するよう案内されます。

### ステータスラインの収集

Claude Codeはステータスラインのコマンドを実行し、JSONドキュメントをパイプで渡します。ccamはそこからモデル名と`rate_limits`を読み取ります。

<!-- motion:usage-statusline -->

1. ccamは、ccamがClaude Codeを起動するときに設定する`CCAM_PROFILE`からアカウントを特定します。ない場合は、`CLAUDE_CONFIG_DIR`をアカウントのフォルダーと照合します（`runStatusline`）。
2. ccamは`rate_limits.five_hour`と`rate_limits.seven_day`を読み取ります（`ParseStatusLine`）。APIキーのアカウントや、セッションの最初の応答より前では、JSONにレート制限が含まれません。その場合は何も保存しません。
3. ccamは`~/.local/state/ccam/usage/<name>.statusline.json`に書き込みます（`SaveStatusLine`）。エンドポイントが応答しないとき、`ccam usage`はこれを使います。
4. ccamは`[name] · model · 5h N% · 7d N%`を出力します。`--wrap`を付けた場合は、同じ入力でコマンドを`sh -c`で実行し、Claude Codeがステータスラインを表示するのと同じように、その出力を表示します。空でないすべての行を表示し、最初の行の前に`[name]`を付けます。失敗したコマンドでは、アカウント名だけが表示されます。Windowsでは、Claude Codeと同じ方法でコマンドを実行します。Git Bashがインストールされていればそれを使い、なければPowerShellを使います。
5. 5時間または週間の使用量が90%以上で、ほかのアカウントがある場合、行の最後に`→ ccam next --auto --continue`が付きます。たとえば`[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`のようになります。

ラップしたコマンド自体には制限時間がありません。Claude Codeはステータスラインのコマンドに10分を与え、新しい更新が始まると実行を止めるので、遅いコマンドもccamなしの場合と同じように表示されます。Claude Codeがccamを終了させると、コマンドも終了します（LinuxとmacOSでは、それが起動したすべてのプロセスも終了します）。

`ccam statusline install`は、各アカウントの`settings.json`に`statusLine`を書き込み、`padding`、`refreshInterval`、`hideVimModeIndicator`などのオブジェクトのほかのキーは残します。アカウントにすでにステータスラインのコマンドがある場合、ccamはそれを`--wrap '<old command>'`として残します。Windowsでは、どちらのシェルでも引用符が不要な`--wrap-base64 <base64>`として残します。ステータスラインがすでにccamを実行している場合（別のパスからのものや、古いバージョンが書いた引用符の付け方のもの）、ccamはラップせずに置き換え、それがラップしているコマンドは残します。

Windowsでは、ccamへのパスをスラッシュとホームフォルダーを表す`~`で書き込みます（`~/.local/bin/ccam.exe statusline`）。これはGit BashとPowerShellのどちらでも解釈できます。引用符が必要なパスは、Claude Codeが使うシェルに合わせて完全な形で書き込みます。Git Bashでは`'…'`、PowerShellでは`& '…'`です。グループで共有されている`settings.json`は、リンクを通じて1回だけ書き込まれるので、グループのすべてのアカウントに反映されます。もう一度実行すると、すでにccamを使っているファイルについては`(already set)`と表示します。

## 安全性と特殊なケース

- `--refresh`を付けない限り、`ccam usage`がトークンファイルを変更することはありません。アクセストークンの期限が切れたときは、`ccam <name>`でそのアカウントを一度起動すれば、Claude Codeが更新します。ログインの期限が完全に切れた場合は`ccam login <name>`が必要です。
- `--refresh`は、実行中のアカウント（ccamを使わずに起動したClaude Codeも含みます）、インポートしたアカウント、キーチェーンのログインには一切触れません。更新を拒否した場合は、保存済みの数値とその理由（`token not renewed: …`）を表示します。
- 使用量エンドポイントは公開ドキュメントに記載されておらず、変更される可能性があります。変更された場合でも保存済みの数値は表示され、`ccam doctor`が形式の変更を報告します。
- 保存済みの数値も応答もない場合、そのアカウントの行には`lookup failed:`と理由が表示されます。たとえば`--source cache`では`no saved usage`となります。
- キャッシュファイルに保存されるのは割合とリセット時刻だけで、トークンが保存されることはありません。

## 関連項目

- [使用量の上限とアカウントの切り替え](/docs/usage-limits)
- [次のアカウントへの切り替え](/features/next)
- [Claude Codeが使用量の上限に達したときの対処法](/guides/usage-limit-workflow)
- [コマンドリファレンス](/docs/commands)
