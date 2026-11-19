---
title: よくある質問
metaTitle: "ccamのよくある質問：Claude Codeの複数アカウント"
description: ccam、Claude Codeのアカウント、使用量の上限、共有、プライバシー、対応システムについての回答です。
updated: 2026-11-19
order: 6
translation: machine
sourceHash: "090ecc3a34bd"
---

## ccamはAnthropicの公式ツールですか？

いいえ。ccamは非公式のオープンソースプロジェクトであり、Anthropicとは提携しておらず、Anthropicの承認も受けていません。ccamは、複数のアカウントを並行して使うために[Claude Codeがドキュメントで説明している](https://code.claude.com/docs/en/env-vars)`CLAUDE_CONFIG_DIR`設定を利用しています。

## 複数のClaude Codeアカウントを使ってもよいのですか？

それは各アカウントの規約によります。仕事用のアカウントの場合は、組織のポリシーにもよります。ProとMaxにはAnthropicの[消費者向け利用規約](https://www.anthropic.com/legal/consumer-terms)、TeamとEnterpriseには[商用利用規約](https://www.anthropic.com/legal/commercial-terms)、すべてに[利用ポリシー](https://www.anthropic.com/legal/aup)が適用されます。複数のアカウントを使う前に確認してください。ccamはアカウントの上限を引き上げたり、まとめたり、隠したりはせず、どのアカウントも自身の上限で止まります。ccamは、あなたが求めたとき、または[`auto_switch`](/features/next#auto-switch)（LinuxとmacOS）を有効にしていれば上限に達した時点で自動的に、あなたの別のアカウントに作業を移します。使う前に、それがあなたのアカウントで許されるかどうかを確認してください。判断するのはあなたです。

**TeamとEnterpriseのシート。** 自動的なアカウント選択（`ccam next`、`--fallback`、`ccam pick`、[`ccam agent`](/features/agent)、ダッシュボードの`n`、有効にした場合の[`auto_switch`](/features/next#auto-switch)）は、同じClaude TeamまたはEnterpriseの組織に属するシート間でも切り替えます。Anthropicのヘルプセンターには、「Each seat on a Team plan is assigned to one named individual, and only that person may use it」（Teamプランの各シートは1人の指定ユーザーに割り当てられ、そのユーザーのみが使用できます）、「Users can't route their work through another user's account to use that user's remaining usage」（ユーザーは別のユーザーのアカウントを通じて作業をルーティングして、そのユーザーの残りの使用量を使用することはできません）と書かれています（[Teamプラン](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)。[Enterpriseプラン](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)も同様です）。そのため、シート間で自動的なアカウント選択を使うと、Anthropicの規約に違反するおそれがあります。そこで使うかどうかはご自身の判断です。頼る前に、所属する組織とAnthropicに確認してください。個人プラン（Pro、Max）については、上の段落を参照してください。

## ccamで使用量の上限を回避できますか？

どのアカウントの上限も超えることはできません。各アカウントは自身の上限を保ち、そこで止まります。ccamが上限を引き上げたり、まとめたり、隠したりすることはありません。ccamは上限を表示し、あるアカウントの上限に達したときに、あなたの別のアカウントで作業を続けられるようにします。`auto_switch`を有効にしていれば、それを自動で行います（次の質問を参照）。1つのTeamまたはEnterpriseの組織のシート間では、別のシートで続けるとそのシートの残りの使用量を使うことになり、Anthropicのヘルプセンターは別のユーザーのシートについてそれを認めていません（上の「TeamとEnterpriseのシート」を参照）。

## セッションが上限に達したとき、ccamは自動でアカウントを切り替えられますか？

有効にした場合だけです。`~/.config/ccam/config.toml`の先頭で`auto_switch = true`を設定すると（LinuxとmacOS。デフォルトではオフで、実験的な機能です）、対話型のセッション（ccamが起動したもの）をClaude Codeが5時間、週間、Opus、Sonnet、Fableのいずれかの上限で止めたとき、ccamはセッションを終了させ、余裕のあるあなたの別のアカウントで会話を続けます。移す先は、同じTeamまたはEnterpriseの組織の別のシートか、別の個人プラン（Pro、Max）です。上限に達する前に切り替えることはなく、バックグラウンドで使用量を監視することもなく、TeamまたはEnterpriseの組織の外へ、あるいは組織の中へ作業を移すことも決してありません。

シートの場合は作業をあるシートから別のシートに移すことになり、Anthropicの規約に違反するおそれがあります（上の「TeamとEnterpriseのシート」を参照）。個人プランについては、Anthropicの規約はこれを明確には扱っていません。消費者向け利用規約は、Anthropicが明示的に許可しない限り「through automated or non-human means」（自動化された手段または人間以外の手段によって）サービスにアクセスすることを禁じており、移動の際にはあなたの代わりにメッセージが1つ送られます。有効にするかどうかはあなたの判断です。

すべての条件と、移動にかかるコストは[`auto_switch`](/features/next#auto-switch)にまとめています。

## 1台のコンピューターで2つのClaude Codeアカウントを使うにはどうすればよいですか？

`CLAUDE_CONFIG_DIR`を使って、各アカウントに専用の設定フォルダーを割り当てます。ccamはこれを自動で行います。アカウントごとに`ccam add`を実行し、`ccam personal`と`ccam work`、または`ccwork`のように自分で名前を付けたショートカットで起動してください。手動で設定する方法は[ガイド](/guides/multiple-claude-code-accounts)でも説明しています。

## 別のアカウントで会話を続けるにはどうすればよいですか？

そのフォルダーで`ccam next --continue`を実行します。フォルダーの最後の会話が次のアカウントに引き継がれます（最も余裕があるアカウントにするには`--auto`）。Claude Codeの中では、[ccamプラグイン](/docs/install#claude-code-plugin)の`/ccam:switch`が、いま使っているセッションに対して同じことをします。2つのアカウントが会話履歴を共有していれば（`ccam group set default --add history`、またはそのアカウント専用のグループ: `ccam group add pair --items history --members work,side`）、新しいアカウントはまさにその会話を開きます。共有していなければ、Claude Codeが会話を新しいアカウントの履歴にコピーします。`--brief`を付けると、会話全体ではなく短い要約から新しいアカウントを始めます。

## ccamはログイン情報を見たり送信したりしますか？

ccamがClaudeのログインをアカウント間で共有したりコピーしたりすることはなく、MCPサーバーのサインイン情報をコピーすることもありません。アカウントのアクセストークンを読み取るのは、Anthropicの使用量エンドポイントにそのアカウントの使用量を問い合わせるためだけです。対象は`ccam usage`、ダッシュボード、`ccam doctor`、そして自動的な選択（`ccam next --auto`、`ccam pick`、`--fallback`、セッションが上限に達したときの`auto_switch`）です。期限切れのトークンを更新するのは`ccam usage --refresh`を実行したときだけです。`usage_source = "claude"`では、トークンをどこにも送らず、保存されたログイン（ファイル、またはmacOSのキーチェーンの項目）は、それがまだ有効かどうかを確かめるためだけに読み取ります。`ccam ext`と`ccam doctor`が読み取るのはアカウントのMCPサインインの名前だけで、そのトークンは読み取りません。ccamにはテレメトリーがありません。これらの問い合わせ以外にccamが自分から送るリクエストは、ダッシュボードが新しいccamリリースがないかGitHubで確認するものだけで、そこにはあなたやあなたのアカウントに関する情報は一切含まれません。

## リモートサーバーでサインインするにはどうすればよいですか？

SSHセッションで、いつもどおり`ccam add`または`ccam login <name>`を実行します。Claude Codeはサーバー上でサインインを待つため、手元のコンピューターのブラウザーは最後に`http://localhost:<port>/callback?code=…`のページに移動しますが、このページは読み込まれません。そのアドレスをブラウザーからコピーしてターミナルに貼り付けてください。ccamがそれをサーバー上のClaude Codeに渡し、Claude Codeがサインインを完了します。アドレスには一度だけ使えるコードが含まれているので、他人と共有しないでください。

## アカウントはどこに保存されますか？

`~/.local/share/ccam/profiles/<name>/`に、アカウントごとに1つのフォルダーとして保存されます。macOSでは、Claude Codeはログイン情報そのものをキーチェーンに、そのフォルダーから導出した名前で保存します。そのため、`ccam rename`はそこでフォルダーの古い名前をそのまま使います。

## 既存の~/.claudeはどうなりますか？

`ccam import`か`ccam migrate`を実行しない限り、何も変わりません。インポートはフォルダーをその場所のまま登録するので、通常の`claude`も`ccam <name>`も、引き続き同じアカウントと設定を使います。`ccam migrate`は、フォルダーを読み取るだけで、その設定、`CLAUDE.md`、拡張機能、MCPサーバーを共有の`default`グループにコピーします。

## ccamはどのシステムに対応していますか？

Linux（Ubuntu、Debian、Fedora、Arch、Alpineなど）、macOS、Windows 10と11（ネイティブ、またはWSLの中）です。Windowsで`settings.json`のようなファイルを共有するには開発者モードが必要です。共有フォルダーは開発者モードなしで動作し、開発者モードがなければ、最初の`ccam add`はMCPサーバーと拡張機能だけを含む`default`グループを作成します（[Windowsでの共有](/features/share#on-windows)）。

## MCPサーバーが1つのアカウントにしか表示されないのはなぜですか？

セッション内で`claude mcp add`を使って追加したサーバーは、そのアカウントにだけ保存され、`-s user`を付けなければ現在のフォルダーにだけ保存されます。ccamは次にそのアカウントを起動したときにそのことを知らせ、`ccam mcp list`と`ccam ext`にもそれらが表示されます。`ccam mcp import <account> <server>`でサーバーを共有グループに移動してください。

## ccamを使うとClaude Codeは遅くなりますか？

ccam経由でClaude Codeを起動すると、起動時間が約4ミリ秒増えます。これはノートパソコンで`scripts/perf-check.sh`を使って測定した値です。ccamは小さなファイルを1つ読み、共有リンクを確認してから自身を`claude`に置き換えるため、デフォルトでは作業中にccamが動作していることはありません。Windowsではプロセスを置き換えられないため、ccamはClaude Codeを起動して終了を待つだけで、ほかには何もしません。`relay`か`auto_switch`を有効にすると、ccamは同じように対話型のセッションごとにそれと並んで動き続け、別のアカウントで会話を続けられるようにします。

## ccamをアップデートするにはどうすればよいですか？

`ccam update`を実行します。最新のリリースをダウンロードしてチェックサムを検証し、`ccam`ファイルを置き換えます。`ccam update --check`は、新しいリリースがあるかどうかを知らせるだけです。ダッシュボードも新しいリリースがあれば知らせ、`U`を押すと確認を取ってからインストールします。ccamが自動で更新されることはありません。`.deb`パッケージでインストールした場合は、代わりに新しい`.deb`をaptでインストールしてください。詳しくは[アップデート](/docs/install#update)を参照してください。

## ccamをアンインストールするにはどうすればよいですか？

`ccam uninstall`を実行します。削除するものを一覧で表示して確認を求めたうえで、すべてのアカウント（先にログアウトします）、ショートカットコマンド、`~/.config/ccam`、`~/.local/share/ccam`、`~/.local/state/ccam`、`ccam`コマンドを削除します。インポートした`~/.claude`は残り、グループで共有していたもののコピーがそこに戻ります。詳しくは[インストール](/docs/install#uninstall)を参照してください。
