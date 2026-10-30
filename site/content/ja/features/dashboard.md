---
title: ダッシュボード
description: ターミナルでccamを実行すると、アカウント、グループ、フォルダーのルール、拡張機能、点検結果の全画面ダッシュボードが開きます。どの変更も先に計画を示します。
updated: 2026-10-30
order: 15
translation: machine
sourceHash: d65f1b5bb0c8
---

引数なしの`ccam`は、標準入力と標準出力がターミナルであれば全画面のダッシュボードを開き、そうでなければ`ccam list`を出力します。ダッシュボードはすべてのアカウントを使用量とグループとともに表示し、起動、追加、ログイン、グループの変更、フォルダーのルールの設定、バックグラウンドセッションの一覧、点検と修復といった、コマンドと同じことを行えます。起動するか、バックグラウンドセッションにアタッチすると、ダッシュボードは閉じます。アカウントの追加やログインでは、少しの間ターミナルを渡し、その後ダッシュボードに戻ります。

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

タブ名の下の行は、現在のタブが何のためのものかを1文で説明します。幅の広いターミナルでは、Accountsタブに各アカウントのメールアドレスも表示されます。タブ名の行の末尾には、どのタブでも、[`auto_switch`](/features/next#auto-switch)が有効かどうかが表示されます（LinuxとmacOS。Windowsでは無視されます）。有効なときは、それにまだ必要なもの（Claude Code 2.1.214以降と、Claudeのサブスクリプションでサインインするすべてのアカウントでオンになっているccamプラグイン）をHealthタブが一覧表示します。

ダッシュボードは開くたびに、新しいccamのリリースがあるかどうかをバックグラウンドでGitHubに問い合わせます。最初の画面がこれを待つことはなく、オフラインのときは何も表示しません。新しいリリースがあれば、タイトルの下の行でそれを勧めます。

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U`はまず確認を取ってから、[`ccam update`](/features/update)と同じ方法で更新します。リリースをSHA-256チェックサムと照合し、1回の操作でccamを置き換えます。開いているダッシュボードは古いバージョンのままなので、終了してから`ccam`をもう一度起動してください。`.deb`パッケージでインストールしたccamのコピーなど、ここで更新できない場合は、`ccam update`と同じ案内を表示し、提案はそのまま残します。

最初の画面は`config.toml`だけから描画され、100 msを大きく下回る時間で表示されます。使用量は、各アカウントの応答があり次第表示されます。上限ごとに専用のバーがあります。5時間の枠、全モデル共通の週間の枠、Fable専用の週間の枠です。`-`はそのアカウントにその枠がないことを示します。割合は80%で黄色、100%で赤になります。取得に失敗した場合は、キャッシュの値が薄く表示されます。Claude Codeが上限に達したと報告した時間枠は、そのリセットまで100%と表示されます。アクセストークンの期限切れは薄い`⚠ token expired`で示され（次の起動時にClaude Codeが更新します）、`⚠ login expired`はそのアカウントに`ccam login`が必要なことを示します。

## タブ

| タブ | 表示内容 | キー |
| --- | --- | --- |
| Accounts | 番号、名前、メールアドレス、5時間・週間・Fableのバー、最終使用、グループ。`✔`はデフォルト、`▸`はこのフォルダーのアカウント、`•`は実行中（ccamを使わずに起動したClaude Code、つまり通常の`claude`、IDE、バックグラウンドセッションとそのデーモンも含みます） | 下記を参照 |
| Agents | [`ccam agents`](/features/agents)の一覧と同じ内容。すべてのアカウントのバックグラウンドセッションと開いているセッションを、状態、フォルダー、開始時刻とともに表示（下記） | `enter`でアタッチ（`ccam attach`と同じ）、`x`で停止、`r`で再読み込み |
| Groups | 各グループと、その共有項目とアカウント、どのグループにも入っていないアカウント | `enter`で編集、`n`で新規作成、`e`で名前変更、`x`で削除 |
| Folders | [`ccam folder`](/features/folder)のルールと、ここで適用されるルール。シンボリックリンクを通ってそのフォルダーに来た場合も同じです（新しいルールはリンクを解決して保存されます） | `n`でここでアカウントを使う、`enter`で変更、`x`で削除 |
| Extensions | [`ccam ext`](/features/ext)の表示内容。各拡張機能を、出どころのグループとともに一度だけ表示。`off:`はプラグインがオンになっていないアカウント、このフォルダーだけに適用されるアカウント自身のサーバー、`This folder's project`、どのアカウントにも属さない問題も表示 | `r`で再読み込み |
| Health | [`ccam doctor`](/features/doctor)の報告内容をアカウントごとのボックスで表示（下記） | `f`で選択中のボックスを修復、`F`ですべて修復（`doctor --fix`）、`↑`/`↓`でボックス間を移動 |

`←`/`→`か`Tab`でタブを切り替え、`?`で現在のタブのキーを一覧表示し、`r`で短期キャッシュを使わずに使用量を更新してAgents・Extensions・Healthタブを読み込み直し、`q`か`Esc`で終了します。

## Accountsタブのキー

| キー | 動作 |
| --- | --- |
| `Enter` | 選択したアカウントとしてClaude Codeを起動します |
| `1`–`9` | その番号のアカウントに移動します。12に移動するには`1` `2`とすばやく入力します。起動するのは`Enter`だけです |
| `↑`/`↓`、`k`/`j` | 移動します |
| `n` | [`ccam next --auto`](/features/next)が選ぶアカウントを起動します |
| `o` | オプションを付けて起動します: `--continue`、`--resume`、そのほかの引数、ドライラン |
| `u` | 使用量の詳細: 5時間、週間、モデルごとの期間と、それぞれのリセット時刻、数値の出どころ |
| `a` | アカウントを追加します（後述のウィザード） |
| `l` | もう一度ログインします |
| `s` | アカウントのグループを選びます |
| `d` | デフォルトのアカウントにします |
| `f` | 現在のフォルダーで使うようにします |
| `e` | 名前を変更します |
| `c` | ショートカットコマンドを設定または削除します |
| `x` | 削除します（先に確認し、最初はNoが選ばれています） |

**TeamとEnterpriseのシート。** 自動的なアカウント選択（`ccam next`、`--fallback`、`ccam pick`、ダッシュボードの`n`、有効にした場合の[`auto_switch`](/features/next#auto-switch)）は、同じClaude TeamまたはEnterpriseの組織に属するシート間でも切り替えます。これはAnthropicの規約に違反するおそれがあります。頼る前に[FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)を読み、所属する組織とAnthropicに確認してください。

## Healthタブ

Healthタブは、開いたときに[`ccam doctor`](/features/doctor)の点検を実行し、同じ表示を描画します。最初の行にアカウントごとの記号が並び、その下にアカウントごとのボックスがあります。ボックスの各行は点検した領域1つで、`Folder`、`Login`、`Share links`が先頭に来ます。記号は`✔`が正常、`ℹ`が情報、`⚠`が警告、`✘`がエラーで、どのアカウントのどこに問題があるかがひと目でわかります。どのアカウントにも属さない結果は`General`ボックスに入ります。

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

`↑`/`↓`でボックス間を移動でき、選択中のボックスはアクセントカラーの枠線で表示されます。`f`は選択中のボックスだけを`ccam doctor --fix`と同じやり方で修復します。アカウントのボックスではそのアカウントだけを、`General`ボックスではどのアカウントにも属さない結果（プラグインのパスとショートカットコマンド）を修復し、ほかのアカウントはもう一度点検するだけです。修復の前に確認し（``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``）、修復するものがないボックスでは`nothing to repair in <name>`と表示するだけです。`F`は`ccam doctor --fix`と同じようにすべてを一度に修復し、`r`はもう一度点検します。

## Agentsタブ

Agentsタブは、[`ccam agents`](/features/agents)が出力する内容、つまりすべてのアカウントのバックグラウンドセッションと開いているセッションを、それぞれが属するアカウントとともに表示します。初めてタブを開いたとき、`r`を押したとき、停止の後、設定の変更の後に、すべてのアカウントのClaude Codeに同時に問い合わせ（`claude agents --json`、それぞれ15秒のタイムアウト）、すべてが応答するまで`Asking every account's Claude Code for its sessions…`と表示します。

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

`State`は、`ccam agents`と同じく、バックグラウンドセッションの状態とそのプロセスが行っていること、または対話型のセッションなら`open`です。アカウントに問い合わせられない場合は、一覧の下の`✘`の行にClaude Codeの返した内容が示されます。上の`lab`がその例で、エージェントビューがオフになっています。終了したバックグラウンドセッションは表示されず（`ccam agents --all`で一覧できます）、すべてのアカウントが応答してどれにもセッションがない場合、タブには``No sessions. Start one in the background with `ccam <account> --bg "…"`.``と表示されます。幅の狭いターミナルでは、`Started`、`Name`、`Folder`の順に列を省き、長いフォルダーのパスは左側から切り詰めます。

バックグラウンドセッションで`Enter`を押すと、ダッシュボードを閉じ、`ccam attach <id>`と同じように、そのセッションをそのアカウントとしてこのターミナルで開きます。ダッシュボードには戻りません（Claude Codeでは、`←`でエージェントビューに、`Ctrl+Z`でシェルに移ります）。開いているセッションにはIDがないので、`Enter`と`x`はバックグラウンドセッションが必要だと伝えるだけです。`x`は、`No, keep it`が選ばれた状態で先に確認します: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.``。「はい」を選ぶと、ccamはそのアカウントとして`claude stop <id>`を実行し、``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again``と、または停止できなかった理由を示して、一覧を読み込み直します。

## グループの変更

Groupsタブで`Enter`を押すとエディターが開きます。左に項目、右にアカウントが並び、`Tab`で切り替え、`Space`でオンとオフを切り替えます。いずれかのアカウントがすでに別のグループから受け取っている項目は、理由とともにグレー表示になります。項目を二重に受け取ることになるアカウントも同様です。1つの項目を受け取れるのは1つのグループからだけです。`Enter`で計画が表示されます。

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

ファイルがグループのコピーと異なる場合、ccamは先にどうするかを尋ねます（グループのものを残す、自分のものを使う、マージする）。その後、変更は[`ccam group`](/features/share)と同じコードでバックグラウンドで実行されます。まず設定を保存し、各アカウントのデータを順に移動し、元のファイルは`.ccam-backup/`に移します。MCPの共有から抜けるアカウントは、グループのサーバーを自分の`.claude.json`に持ち続けます。

実行中のアカウントのデータを移動することになる変更は、計画を開いている間にアカウントが起動した場合も含め、そのアカウントの名前を示して拒否されます。ほかのccamプロセスがグループを変更したために古くなった計画は、`the groups changed while this plan waited; nothing was moved, make the change again`で拒否されます。ほかのccamプロセスがグループを変更している間は、新しい変更、名前の変更、削除は`another ccam process is changing the groups right now`と表示します。途中で止まった変更の後は、Healthタブの`F`（または`ccam doctor --fix`）がそれを完了させるまで待ちます。

## アカウントの追加

`a`を押すと、ターミナルで`ccam add`を実行したときと同じウィザードが開きます。

1. アカウント名。
2. サインインの方法: 今すぐClaudeのサブスクリプションで（ブラウザーで）、[APIキー](/features/api)、またはあとでClaudeのサブスクリプションで。
3. APIキーの場合: 起動のたびにccamがキーを読む場所（`env:NAME`か`command:COMMAND`）、エンドポイント、キーの送り方、モデル。
4. 加わるグループ（`default`にはチェックが入っています）。グループがまだない場合は、`default`グループが共有する予定の項目。
5. 必要であれば、ショートカットコマンド。
6. 確認用の概要。

`Esc`で1つ前の手順に戻り、最初の手順ではウィザードを閉じます。

## アカウントの削除

`x`は、何かを行う前に、Noを選択した状態で確認を求め、削除で何が行われるかを示します: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` APIキーのアカウントには取り消すログインがなく、ほかのアカウントが同じログインを持つ場合はログアウトされません。インポートしたアカウントはフォルダーを残します: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## 色

色はClaude CodeのターミナルUIの役割に合わせているので、ccamはClaude Codeと並べても違和感がありません。グループ名には、名前から決まる色が付きます。ccamはターミナルに背景色を問い合わせ、ライトかダークのテーマを選びます。`config.toml`の先頭、最初の`[[group]]`、`[[profile]]`、`[[folder]]`より上で自分で設定することもできます。

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

`-daltonized`のテーマは赤と緑を見分けにくい色覚向けで、`-ansi`のテーマはターミナル自身の16色を使います。`NO_COLOR`を設定すると色を使いません。

## 安全性と特殊なケース

- ダッシュボードはコマンドと同じ関数を呼び出すので、すべての規則がそのまま適用されます。データを移動する前の計画、実行中のアカウントに対する拒否、バックアップ、そしてログイントークンがアカウントの外に出ないことです。
- ダッシュボードにはドライランがありません。名前を付けない`ccam --dry-run`は、このフォルダーのアカウントがどう起動するかを表示します。アカウントで`o`を押すと、その1回の起動だけをドライランにできます。
- Claude Codeが起動する前やログインが開く前に、ターミナルは元の状態に戻ります。
- ダイアログに対して画面が小さすぎる場合は、まずタブの行を減らします。行は折り返さずに`…`で切り詰めます。高さが足りない画面では、次に空行とタブの説明を省くので、キーのヒントの行は表示されたままです。40列または8行より小さいと、ダッシュボードは`Too small: ccam needs 40x8 (now 36x7)`だけを表示します。
- 変更、削除、アップデート、修復の適用中に`Ctrl+C`を押すと、その完了を待ちます。ステータス行には`ccam quits when it is done.`と表示されます。
- `U`でccamをアップデートした後は、アカウントの追加やログインからダッシュボードに戻っても、アップデートの案内は再び表示されません。

## 関連項目

- [グループによる共有](/features/share)
- [Claude Codeの起動](/features/launch)
- [アカウントの管理](/features/accounts)
- [使用量の上限](/features/usage)
