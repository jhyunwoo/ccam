---
title: アカウントの管理
description: ccam add、ccam login、ccam list、ccam rename、ccam remove、ccam importで、Claude Codeのアカウントを追加、ログイン、一覧、名前変更、削除、インポートします。
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

ccamの各アカウントは、独自のログインを持つ別々のClaude Codeの設定フォルダーです。ここで紹介するコマンドで、アカウントの追加、サインイン、一覧表示、名前変更、削除を行います。すでにあるフォルダーの登録もできます。

## 使い方

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account that was never logged in or whose login expired
ccam logout side                           # log out: the login and its MCP sign-ins
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, logs out, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

コマンドが`<name|number>`を受け取る場所では、アカウント名かスロット番号（`ccam list`の`#`列）を指定できます。

`ccam use`は、どのアカウントをデフォルトにするか（`config.toml`の`default`）を変えるだけです。

ターミナルでオプションなしに`ccam add`を実行すると、ウィザードが開きます。ウィザードは次の手順で進みます。

1. アカウント名（指定していない場合）
2. アカウントのサインイン方法: 今すぐClaudeにログイン、APIキー、あとでClaudeにログイン
3. 加わるグループ。`default`にチェックが入った状態で表示されます（グループがまだない場合は、`default`グループが共有する項目を表示します）
4. `Create a shortcut command for work? [y/N]`
5. 確認用の概要

`Esc`で1つ前の手順に戻ります。`--group`、`--no-group`、`--shortcut`、`--no-login`のいずれかを付けた場合や、ターミナルなしで（たとえばパイプから）実行した場合は、ウィザードは開きません。その場合は名前が必須で、`--shortcut`で名前を指定しない限りショートカットは作成されません。

`--group`を付けない場合、新しいアカウントは`default`グループに加わります。最初のアカウントがこのグループを作成し、MCPサーバー、設定、`CLAUDE.md`、拡張機能を共有します。[グループによる共有](/features/share)を参照してください。Windowsでは、ファイルをリンクできるのは開発者モードのときだけです。開発者モードがない場合、最初のアカウントは`mcp`と`extensions`だけを含む`default`を作成してその旨を知らせ（`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`）、ccamは設定、メモリー、履歴を共有するグループへの`--group`を拒否します。

`ccam list`は、スロット、名前、ショートカット、メールアドレス、グループ、最終使用日時を表示し、デフォルトのアカウントの前に`*`を付けます。メールアドレスは、そのフォルダーが現在ログインしているもの（その`.claude.json`から）です。それがccamの記録と異なる場合は、記録されたものが`(was …)`として後に続きます。

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

```console
$ ccam list
   #  Name      Shortcut  Email                                     Groups   Last used
*  1  work      ccwork    other@company.com (was work@example.com)  default  2m ago
```

`--json`を付けると、各アカウントは`slot`、`shortcut`、`name`、`email`（表と同じく現在のログインのもの）、`default`、`dir`、`last_used`（一度も起動していないアカウントでは省略されます）、`groups`、`share`を持つオブジェクトになります。`share`は、各項目（`mcp`、`settings`、`memory`、`extensions`、`history`）を、その項目の出どころのグループか`"local"`に対応付けます。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `ccam add --no-login` | ログインせずにアカウントを作成します。あとで`ccam login <name>`でログインします | オフ |
| `ccam add --group <group>` | `default`の代わりにこのグループに加わります。複数回指定できます。`default`は、存在しない場合に作成されます | `default` |
| `ccam add --no-group` | どのグループにも加わりません。アカウントは何も共有しません | オフ |
| `ccam add --shortcut <command>` | このショートカットコマンドも`~/.local/bin`に作成します | なし |
| `ccam add --api --key-from <source>` | ログインの代わりにAPIキーを使うアカウント。[APIキーのアカウント](/features/api)を参照 | オフ |
| `ccam list --json` | アカウントをJSONで出力します | オフ |
| `ccam remove -y`、`--yes` | 確認なしで削除します | オフ |
| `ccam logout -y`、`--yes` | 確認なしでログアウトします | オフ |
| `ccam import --name <name>` | 指定したフォルダーのアカウント名。フォルダーを指定しないと拒否されます（`--name needs a folder: ccam import <folder> --name ops`） | フォルダー名から決まります |
| `ccam import -y`、`--yes` | 見つかったすべてのフォルダーを確認なしでインポートします | オフ |
| `--dry-run` | 何も変更しない`ccam list`と`ccam completion`では使えます。`add`、`login`、`logout`、`use`、`rename`、`remove`、`import`では拒否され、`` --dry-run is not supported by `ccam add`; nothing was changed ``（コマンド名入り）と表示されます | オフ |

`ccam login`（`--expired`以外）、`ccam use`、`ccam rename`には独自のオプションはありません。`ccam completion`は`bash`、`zsh`、`powershell`のいずれかを受け取り、それ以外を指定すると`only bash, zsh and powershell are supported`で失敗します。

アカウント名には、`a-z`、`0-9`、`-`、`_`を使って1～32文字で指定し、先頭は英字か数字にします。`ccam 2`はスロット2を意味するため、数字だけの名前は使えません。また、`ccam <name>`が曖昧にならないように、ccamのコマンド名（`add`、`list`、`ls`、`run`、`use`、`next`、`usage`、`u`、`share`、`mcp`、`login`、`logout`、`rename`、`remove`、`rm`、`import`、`migrate`、`doctor`、`shims`、`shortcut`、`completion`、`help`、`version`、`statusline`、`update`、`uninstall`、`folder`、`pick`、`ext`、`group`、`groups`、`hook`、`agent`、`agents`、`attach`、`ccusage`）も使えません。この一覧は`internal/config`の`config.Reserved`です。

## しくみ

アカウントを追加すると、フォルダーの作成、グループへの参加、設定の保存、サインイン、ショートカットの作成がこの順序で行われます。`ccam add work --shortcut ccwork`を実行すると、ccamは次の手順を実行します（`internal/cli/profiles.go`の`addProfile`）。

<!-- motion:accounts-add -->

1. ccamは、何かを作成する前に名前（`config.ValidateName`）とショートカット（`shim.Check`）を確認します。
2. ccamは`~/.local/share/ccam/profiles/work/`をモード`0700`で作成し、空いている最小のスロット番号をアカウントに割り当てます。そのため、削除したアカウントの番号は再利用されます（`profile.Add`、`Config.NextSlot`）。最初のアカウントがデフォルトになります。
3. アカウントがグループに加わり（`--group`か`--no-group`で指定しない限り`default`。最初のアカウントが`default`を作成します）、ccamは作成したばかりの空のフォルダーをそれらのグループにリンクします。たとえば`settings.json`を`groups/default/settings.json`にリンクします（`applyShare`。`share.Build`の後に`Apply`）。これが失敗した場合は、新しいフォルダーを削除します。フォルダーが残ると、再試行のたびに`folder already exists`で失敗するためです。
4. アカウントは`fsutil.AtomicWrite`を通じて`~/.config/ccam/config.toml`に保存されます。一時ファイルに書き込み、`fsync`してから`rename`するので、クラッシュしても設定が中途半端な状態で残ることはありません。
5. ccamは、`CLAUDE_CONFIG_DIR`をアカウントのフォルダーに設定し、ログインページに入力済みになるようにアカウントの前回のメールアドレスを`--email`に指定して、`claude auth login`を起動します（`runLogin`）。ログインに失敗してもアカウントは残り、ccamはエラーを表示して、あとで`ccam login work`を実行するよう案内します。
6. ccamは`profiles/work/.claude.json`から`oauthAccount`を読み取ります（`profile.ReadIdentity`）。メールアドレス、アカウントと組織のID、組織の名前、プラン（`claude_max`、`claude_team` …）です。このファイルは数メガバイトになることがあるため、ccamはファイル全体を解析せずに、その1つのキーだけを探します（`TopLevelValue`）。ccamはこれらを`config.toml`の`[profile.login]`に保存します。
7. ccamは`~/.local/bin/ccwork`をccamのバイナリにリンクし（ccamは起動時に`config.toml`で`argv[0]`を調べます）、そのリンクを`~/.local/state/ccam/shortcuts.json`に記録します（`shim.Install`）。リンクを作れなかった場合もアカウントは追加され、ccamがその理由を示し（`Could not create shortcuts: …`）、概要には`shortcut: ccwork (not created; see above)`と表示されます。

ログインには数分かかることがあり、その間に別のccamプロセスが設定を保存するかもしれません。そのため、待ち時間のあるコマンドは`config.toml`を読み込み直し、1つの変更（ここではメールアドレス）だけを適用して保存します（`env.update`）。`ccam remove`も同じ理由で、確認の後に読み込み直します。

`ccam login`は、既存のアカウントについて手順5と6を実行します。Claude Codeはメールアドレスをログインページのヒントとしてしか使わないので、以前とは別のClaudeアカウントにサインインしてしまうことがあります。ccamは新しいログインを`[profile.login]`と比較し、その旨を表示します（`ccam: this login is other@company.com, but work was dev@company.com before.`）。新しいログインはそのまま残し、元に戻す方法を案内します: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

別のClaudeアカウントへのログインでは、アカウントの[上限の記録](/features/usage#the-limit-record)も破棄されます。その上限は別のアカウントのものだったためです。別のccamのアカウントが同じClaudeアカウントにログインしている場合、ccamは2つのアカウントが使用量の上限を共有することを警告します。

`ccam login --expired`は、一度もログインしていないアカウントと、ログインの期限が切れたアカウントに、1つずつ順にログインします。期限切れとは、リフレッシュトークンの期限が切れた場合か、サーバーに拒否された後にClaude Codeがログインを消去した場合です。最初のブラウザーを開く前に、`config.toml`を保存できるかを確認します。

`ccam logout`はアカウント内で`claude auth logout`を実行し、保存済みのメールアドレスを消去します。ログインの記録は残るので、次のログインはそれと照合されます。Claude Codeはログインを取り消し、そのアカウントについて保存しているもの（MCPサーバーのサインイン情報を含む）を削除します。そのため、次の場合、ccamは警告して先に確認を求めます（`--yes`で確認を省略します）。

- ログインがMCPのサインイン情報を持つ場合
- アカウントが実行中の場合（そのセッションもサインアウトされます）
- ほかのアカウントのフォルダーが同じログインを持つ場合（コピーしたフォルダー）。ログアウトでそのログインも取り消されます

例:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

別のClaudeアカウントとしてログインするには、`ccam login <name>`だけで十分で、MCPのサインイン情報も残ります。

`ccam list`は、アカウントが複数の組織に属している場合に`Organization`列を追加し（個人プランは`personal`と表示されます）、`ccam list --json`はすべてのアカウントに`plan`と`org`を含めます。

### リモートマシンでのサインイン

SSH経由では、`claude auth login`のブラウザー側の処理がうまくいかないことがあります。Claude Codeはリモートマシン上のコールバックサーバー（`http://localhost:<port>/callback`）でサインインを待ちますが、サインインを完了するブラウザーは手元のコンピューターで動いているため、そのアドレスはどこにもつながらず、ブラウザーには`ERR_CONNECTION_REFUSED`のようなエラーが表示されます。それでもアドレスにはワンタイムコードが含まれているので、ccamがそれを届けられます。

<!-- motion:accounts-relay -->

ccamはターミナルをClaude Codeに渡しません（`runLogin`）。Claude Codeの標準入力をパイプにつなぎ、入力された各行をそこにコピーします（`relayLogin`）。そのため、Claude Codeの`Paste code here`プロンプトに貼り付けたコードもそのまま届きます。行がコールバックアドレスの場合、ccamは、ホストが`localhost`、`127.0.0.1`、`[::1]`のいずれかで、ポートが指定され、パスが`/callback`で、`code`と`state`の両方がある場合にだけ受け付けます。`http://`は省略できます（`callbackURL`）。続いてccamは、ブラウザーが本来行うのと同じように、10秒のタイムアウトを付け、リダイレクトをたどらずに、サーバー上でそのアドレスにリクエストを送ります（`deliverCallback`）。リクエストは、待機中の`claude auth login`のコールバックサーバー（図の破線）に届きます。Claude Codeはコードを受け取ってログインを保存します。Linuxではアカウントのフォルダー内の`.credentials.json`に保存されます。貼り付けた行によって、ccamがほかのホストに接続することはありません。

ccamはログインを始める前に次のメッセージを表示します。

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

成功すると`→ Passed the sign-in to Claude Code on this machine.`と表示します。

ログインがターミナルで入力を待っている間、ccamはターミナルのエコーと行モードをオフにし、入力した文字ごとに`*`を表示します。そのため、貼り付けたコードやコールバックアドレスは画面に表示されません。行の編集はccamが自分で行います（`readLine`）。`Enter`で行を送信し、`Backspace`で1文字消し、空の行で`Ctrl+D`を押すと入力が終わります。ターミナルの設定は、ログインが終わったときに元に戻ります。`Ctrl+C`、`SIGTERM`、`SIGHUP`を受けたときも元に戻り、これらは従来どおりccamを終了します。標準入力と標準エラー出力がターミナルでない場合は、行を届いたとおりにコピーします。

### 名前の変更、削除、インポート

`ccam rename`と`ccam remove`は、アカウントのフォルダーと、ccamがアカウント名で別の場所に保持しているものを対象にします。ccamはそれらのパスを`nameState`で列挙します。

- `~/.local/state/ccam/last-used/<name>`
- `~/.local/state/ccam/mcp/<name>`
- 使用量のキャッシュ`~/.local/state/ccam/usage/<name>.json`と`<name>.statusline.json`
- 上限の記録`<name>.limit.json`
- `~/.local/state/ccam/handoff/<name>/`にあるセッションのメモ

`ccam remove work`は次の順序でこれらを処理します。

<!-- motion:accounts-remove -->

1. `ccam remove`は、まずアカウントが実行中なら拒否し、そうでなければ確認を求め、回答の後にもう一度確認します（質問の間に起動したセッションも見逃しません）。
2. 何かを変更する前に、`config.toml`を保存できるかを確認します。
3. アカウントが拡張機能を共有している場合、ccamはグループの`installed_plugins.json`と`known_marketplaces.json`を書き換えます（`RewritePluginPaths`）。Claude Codeはプラグインの絶対パスを保存しており、それがこのアカウントのフォルダーを経由している可能性があるためです。
4. 管理しているアカウント（`ccam import`で登録したものではなく、ccamがフォルダーを作成したアカウント）は、次にログアウトされます。そのフォルダーで`claude auth logout`を実行するので、使えるログインがごみ箱に残ることはありません（macOSでは、フォルダーにちなんだ名前のキーチェーン項目にも残りません。残っていると、あとで同じ名前で作ったアカウントがそれを引き継いでしまいます）。APIキーのアカウントと、ほかのアカウントのフォルダーが同じログインを持つ場合はこれを省略します。後者では、ログアウトがそのログインも取り消してしまうためです。ログアウトに失敗した場合は報告し、削除は続けます。
5. ccamはフォルダーを`~/.local/share/ccam/trash/<name>-<time>/`に移動し（`profile.Trash`）、アカウントは`config.toml`から外れます。あとで同じ名前のアカウントを作ったときに引き継がれないように`nameState`のパスは削除され、ショートカットのリンクも削除されます: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`。

結局`config.toml`を保存できなかった場合は、フォルダーがごみ箱から戻されます。削除したアカウントを元に戻すには、ごみ箱からフォルダーを戻し、もう一度ログインします。

`ccam rename`はフォルダーの名前を変更し、その`plugins/`内のプラグインのパスを古いフォルダーから新しいフォルダーに書き換えます。`nameState`のパスを新しい名前に移動し、アカウントがデフォルトだった場合はデフォルトのままにします。セッションのメモも更新されるので、`ccam next --continue`はそれを引き続き見つけられます。ショートカットコマンドの名前は変わりません。プラグインのパスの書き換えに失敗した場合や、`config.toml`を保存できない場合は、フォルダーの名前を元に戻します。

macOSと、Claude Codeがログインを資格情報マネージャーに保存している場合のWindowsでは、フォルダー自体は古い名前のままです。ログインとMCPのサインイン情報はフォルダーのパスで保存されており、フォルダーを移動するとログアウトされてしまうためです。その場合、設定には新しい名前と並んでフォルダーが`folder`として記録され、ccamは`its folder stays …: the macOS Keychain keeps the login under the folder's path`と表示します。

`ccam import`は、フォルダーをその場所のまま`config.toml`の`dir`として登録します。移動もコピーもしません。リフレッシュトークンをコピーすると、コピー同士が互いをログアウトさせてしまうためです。

フォルダーを指定しない場合は、登録されていない`~/.claude`と`~/.claude-*`のフォルダーを探し、それぞれについて確認します。名前はフォルダーから決まります（`~/.claude`は`main`、`~/.claude-work`は`work`、`~/.claude-2`は`acct2`）。ccamが名前を決められない場合や、名前がすでに使われている場合は、実行すべき`--name`コマンドとともにそのフォルダーを一覧表示します。見つかったフォルダーのうち失敗したものは`Skipped <path>: <reason>`と報告され、ほかのフォルダーはインポートされます。

ccamは先にすべての質問をしてから、`config.toml`を読み込み直してフォルダーを登録するので、その間にほかのccamのコマンドが保存した内容は残ります。`✓`の行は保存の後に表示されます。

`CLAUDE_CONFIG_DIR`がない場合、Claude Codeは`.claude.json`をホームフォルダーの`~/.claude`の隣に置きます。そのため、`~/.claude`をインポートしたときに独自の`.claude.json`がなければ、ccamは`~/.claude.json`を指すリンク`~/.claude/.claude.json`を作成します。これで通常の`claude`と`ccam main`が同じアカウントを参照します。

`~/.claude`に、`~/.claude.json`とは別に独自の`.claude.json`がある場合、ccamは両方を残し、その旨を知らせます。`ccam main`は内側のファイルを、通常の`claude`は外側のファイルを使うので、ログインの詳細とMCPサーバーが別々になります。1つにまとめるには、内側のファイルをリンクに置き換えてください（[doctor](/features/doctor)もこれを報告します）。

フォルダーが登録済みのアカウントと同じログインを持つ場合（コピーしたフォルダー）も、ccamはその旨を知らせます。先に更新した方がもう一方をログアウトさせるので、使う前に`ccam login <name>`を実行してください。

インポートしたアカウントは、グループに追加する（`ccam group join <group> <name>`）まではどのグループにも入らず、ショートカットも作成されません。

`ccam completion bash|zsh`は、アカウント名も補完する補完スクリプトを出力します。

## 安全性と特殊なケース

- `ccam rename`と`ccam remove`は、アカウントの実行中は``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``と表示して拒否します（`refuseRunning`）。セッションがフォルダーに書き込み続けるためです。この確認は、Claude Codeがどう起動されたかにかかわらず対象になります。ccam経由、通常の`claude`、IDE、バックグラウンドセッションとそのデーモンです。
- `ccam rename`、`ccam remove`、`ccam logout`は、ほかのccamプロセスがそのログインを更新または変更している間は拒否します: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``。
- ほかのccamプロセスでグループの変更が進んでいる間、`ccam rename`と`ccam remove`は`another ccam process is changing the groups right now; try again when it is done`で拒否します。変更が途中で止まった後は``the last group change stopped halfway; finish it with `ccam doctor --fix` first``で拒否します。

- インポートしたアカウントに`ccam remove`を実行しても、通常の`claude`のために、フォルダーとそのログインは残ります。アカウントがグループに入っている場合は、まず`ccam group leave`と同じようにグループから抜けます。共有していた各項目は、フォルダー内のグループのもののコピーになり、グループのMCPサーバーはその`.claude.json`に入ります。その後、登録が解除されます: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`。ccamが作ったその`.claude.json`のバックアップは`<folder>/.ccam-backup/mcp-<time>`に移ります。移せない場合は`~/.local/state/ccam/mcp/<name>/backups`に残り、メッセージにその旨が示されます。`~/.claude`がごみ箱に移動されることはありません。
- 2つのアカウントが同じフォルダーを持つことはできません。`config.toml`で1つのフォルダーを2回指定すると、読み込みに失敗します。
- `ccam add`は、`profiles/`にすでにフォルダーがある名前を拒否し、`ccam import`を案内します。macOSでは、古いバージョンのccamが残した、新しいフォルダーのパスに対するログインやMCPのサインイン情報がキーチェーンにまだある場合も停止し、それらを削除する`security delete-generic-password`のコマンドを示します。`config.toml`を保存できない場合は、新しいフォルダーをまた削除します。
- `ccam import`は、すでに登録されているフォルダー、フォルダーではないパス、使用済みの名前を拒否します。また、フォルダーが`~/.claude`のコピー（同じログイン）である場合は、`~/.claude`自体が登録されていなくても警告します。
- `ccam logout`は、アカウントが実行中の場合は先に確認を求めます。`ccam login`は確認しません。
- ログインが終了した後などで、中継がClaude Codeに届かない場合は、`ccam: could not pass the sign-in to Claude Code:`と理由が表示されます。HTTPエラーは通常、コードの期限切れを意味します。ログインをやり直してください。

## 関連項目

- [Claude Codeの起動](/features/launch)
- [ショートカットコマンド](/features/shortcuts)
- [MCPサーバー、設定、履歴の共有](/docs/sharing)
- [1台のPCで複数のClaude Codeアカウントを使う](/guides/multiple-claude-code-accounts)
