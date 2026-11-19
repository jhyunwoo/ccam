---
title: "اکاؤنٹس کا انتظام"
description: "ccam add، ccam login، ccam list، ccam rename، ccam remove اور ccam import: Claude Code اکاؤنٹ شامل کریں، لاگ اِن کریں، دیکھیں، نام بدلیں، ہٹائیں، امپورٹ کریں۔"
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

ccam کا ہر اکاؤنٹ Claude Code کا ایک الگ کنفیگ فولڈر ہے جس کا اپنا لاگ اِن ہوتا ہے۔ ان کمانڈز سے اکاؤنٹس شامل کریں، ان میں سائن اِن کریں، ان کی فہرست دیکھیں، ان کے نام بدلیں اور انہیں ہٹائیں۔ یہ ان فولڈرز کو بھی رجسٹر کرتی ہیں جو آپ کے پاس پہلے سے ہیں۔

## استعمال

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

جہاں بھی کوئی کمانڈ `<name|number>` لیتی ہے، آپ اکاؤنٹ کا نام یا اس کا سلاٹ نمبر (`ccam list` کا `#` کالم) دے سکتے ہیں۔

`ccam use` صرف یہ بدلتا ہے کہ کون سا اکاؤنٹ ڈیفالٹ ہے (`config.toml` میں `default`)۔

ٹرمینل پر بغیر آپشنز کے `ccam add` ایک وزرڈ کھولتا ہے۔ یہ ان مراحل سے گزرتا ہے:

1. اکاؤنٹ کا نام، اگر آپ نے نہ دیا ہو؛
2. اکاؤنٹ کیسے سائن اِن کرے گا: ابھی Claude لاگ اِن، API key یا بعد میں Claude لاگ اِن؛
3. کن گروپس میں شامل ہونا ہے، جن میں `default` چیک ہوتا ہے (اگر ابھی کوئی گروپ نہ ہو تو یہ دکھاتا ہے کہ `default` گروپ کیا شیئر کرے گا)؛
4. `Create a shortcut command for work? [y/N]`؛
5. تصدیق کے لیے ایک خلاصہ۔

`Esc` ایک قدم پیچھے جاتا ہے۔ `--group`، `--no-group`، `--shortcut` یا `--no-login` میں سے کوئی بھی، یا ٹرمینل کے بغیر چلانا (مثلاً pipe سے)، وزرڈ کو چھوڑ دیتا ہے۔ تب نام دینا لازمی ہے، اور کوئی شارٹ کٹ نہیں بنتا جب تک `--shortcut` کسی کا نام نہ دے۔

`--group` کے بغیر نیا اکاؤنٹ `default` گروپ میں شامل ہوتا ہے۔ پہلا اکاؤنٹ اسے بناتا ہے، جو MCP سرورز، سیٹنگز، `CLAUDE.md` اور ایکسٹینشنز شیئر کرتا ہے؛ دیکھیں [گروپس کے ساتھ شیئرنگ](/features/share)۔ Windows پر فائلیں صرف Developer Mode میں لنک ہو سکتی ہیں۔ اس کے بغیر پہلا اکاؤنٹ `default` کو صرف `mcp` اور `extensions` کے ساتھ بناتا ہے اور یہ بتا دیتا ہے (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`)، اور ccam کسی ایسے گروپ کے ساتھ `--group` مسترد کر دیتا ہے جو سیٹنگز، میموری یا ہسٹری شیئر کرتا ہو۔

`ccam list` سلاٹ، نام، شارٹ کٹ، ای میل، گروپس اور آخری استعمال دکھاتا ہے، اور ڈیفالٹ اکاؤنٹ سے پہلے `*` ہوتا ہے۔ ای میل وہ ہے جس سے فولڈر اس وقت لاگ اِن ہے (اس کی `.claude.json` سے)؛ جب یہ اس ای میل سے مختلف ہو جو ccam نے ریکارڈ کی تھی، تو ریکارڈ شدہ ای میل اس کے بعد `(was …)` کے طور پر آتی ہے:

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

`--json` کے ساتھ ہر اکاؤنٹ ایک آبجیکٹ ہے جس میں `slot`، `shortcut`، `name`، `email` (موجودہ لاگ اِن کی، جیسے ٹیبل میں)، `default`، `dir`، `last_used` (جب اکاؤنٹ کبھی شروع نہ ہوا ہو تو شامل نہیں)، `groups` اور `share` ہوتے ہیں۔ `share` ہر آئٹم (`mcp`، `settings`، `memory`، `extensions`، `history`) کو اس گروپ سے جوڑتا ہے جس سے وہ آتا ہے، یا `"local"` سے۔

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `ccam add --no-login` | لاگ اِن کیے بغیر اکاؤنٹ بناتا ہے؛ بعد میں `ccam login <name>` سے لاگ اِن کریں | بند |
| `ccam add --group <group>` | `default` کے بجائے اس گروپ میں شامل ہوتا ہے۔ کئی بار دیا جا سکتا ہے؛ `default` موجود نہ ہو تو بنا دیا جاتا ہے | `default` |
| `ccam add --no-group` | کسی گروپ میں شامل نہیں ہوتا: اکاؤنٹ کچھ شیئر نہیں کرتا | بند |
| `ccam add --shortcut <command>` | `~/.local/bin` میں یہ شارٹ کٹ کمانڈ بھی بناتا ہے | کوئی نہیں |
| `ccam add --api --key-from <source>` | ایسا اکاؤنٹ جو لاگ اِن کے بجائے API key استعمال کرتا ہے؛ دیکھیں [API key اکاؤنٹس](/features/api) | بند |
| `ccam list --json` | اکاؤنٹس کو JSON میں پرنٹ کرتا ہے | بند |
| `ccam remove -y`، `--yes` | پوچھے بغیر ہٹاتا ہے | بند |
| `ccam logout -y`، `--yes` | پوچھے بغیر لاگ آؤٹ کرتا ہے | بند |
| `ccam import --name <name>` | آپ کے دیے گئے فولڈر کے لیے اکاؤنٹ کا نام؛ فولڈر کے بغیر مسترد (`--name needs a folder: ccam import <folder> --name ops`) | فولڈر کے نام سے اخذ |
| `ccam import -y`، `--yes` | ملنے والا ہر فولڈر پوچھے بغیر امپورٹ کرتا ہے | بند |
| `--dry-run` | `ccam list` اور `ccam completion` اسے قبول کرتے ہیں، کیونکہ وہ کچھ تبدیل نہیں کرتے۔ `add`، `login`، `logout`، `use`، `rename`، `remove` اور `import` اسے `` --dry-run is not supported by `ccam add`; nothing was changed `` کے ساتھ مسترد کرتے ہیں (کمانڈ کے نام کے ساتھ) | بند |

`ccam login` (`--expired` کے علاوہ)، `ccam use` اور `ccam rename` کے اپنے کوئی آپشنز نہیں۔ `ccam completion`، `bash`، `zsh` یا `powershell` لیتا ہے؛ اس کے علاوہ کچھ بھی `only bash, zsh and powershell are supported` کے ساتھ ناکام ہوتا ہے۔

اکاؤنٹ کے نام `a-z`، `0-9`، `-` اور `_` کے 1 سے 32 حروف ہوتے ہیں، جو کسی حرف یا ہندسے سے شروع ہوں۔ نام صرف ہندسوں پر مشتمل نہیں ہو سکتا، کیونکہ `ccam 2` کا مطلب سلاٹ 2 ہے۔ نام کوئی ccam کمانڈ بھی نہیں ہو سکتا (`add`، `list`، `ls`، `run`، `use`، `next`، `usage`، `u`، `share`، `mcp`، `login`، `logout`، `rename`، `remove`، `rm`، `import`، `migrate`، `doctor`، `shims`، `shortcut`، `completion`، `help`، `version`، `statusline`، `update`، `uninstall`، `folder`، `pick`، `ext`، `group`، `groups`، `hook`، `agent`، `agents`، `attach`، `ccusage`)، اس لیے `ccam <name>` کبھی مبہم نہیں ہوتا۔ یہ فہرست `internal/config` میں `config.Reserved` ہے۔

## یہ کیسے کام کرتا ہے

اکاؤنٹ شامل کرنا اس ترتیب سے اس کا فولڈر بناتا ہے، اس کے گروپس میں شامل ہوتا ہے، کنفیگ محفوظ کرتا ہے، سائن اِن کرتا ہے اور شارٹ کٹ بناتا ہے۔ جب آپ `ccam add work --shortcut ccwork` چلاتے ہیں تو ccam یہ مراحل طے کرتا ہے (`internal/cli/profiles.go` میں `addProfile`):

<!-- motion:accounts-add -->

1. کچھ بھی بنانے سے پہلے ccam نام (`config.ValidateName`) اور شارٹ کٹ (`shim.Check`) جانچتا ہے۔
2. ccam، `~/.local/share/ccam/profiles/work/` کو موڈ `0700` کے ساتھ بناتا ہے اور اکاؤنٹ کو سب سے کم خالی سلاٹ نمبر دیتا ہے، اس لیے ہٹائے گئے اکاؤنٹ کا نمبر دوبارہ استعمال ہوتا ہے (`profile.Add`، `Config.NextSlot`)۔ پہلا اکاؤنٹ ڈیفالٹ بن جاتا ہے۔
3. اکاؤنٹ اپنے گروپس میں شامل ہوتا ہے (`default`، جب تک `--group` یا `--no-group` کچھ اور نہ کہے؛ پہلا اکاؤنٹ `default` بناتا ہے)، اور ccam نئے، خالی فولڈر کو ان سے لنک کرتا ہے، مثلاً `settings.json` کو `groups/default/settings.json` سے (`applyShare`، `share.Build` اور پھر `Apply` کے ساتھ)۔ اگر یہ ناکام ہو تو نیا فولڈر ہٹا دیا جاتا ہے، کیونکہ بچا ہوا فولڈر ہر دوبارہ کوشش کو `folder already exists` کے ساتھ روک دے گا۔
4. اکاؤنٹ `fsutil.AtomicWrite` کے ذریعے `~/.config/ccam/config.toml` میں محفوظ ہوتا ہے: ایک عارضی فائل، `fsync`، پھر `rename`، اس لیے کریش کبھی آدھا لکھا ہوا کنفیگ نہیں چھوڑتا۔
5. ccam، `CLAUDE_CONFIG_DIR` کو اکاؤنٹ فولڈر پر مقرر کر کے `claude auth login` شروع کرتا ہے، اور اکاؤنٹ کی آخری ای میل کے ساتھ `--email` دیتا ہے تاکہ لاگ اِن صفحہ پہلے سے بھرا ہو (`runLogin`)۔ اگر لاگ اِن ناکام ہو تو اکاؤنٹ برقرار رہتا ہے، اور ccam غلطی پرنٹ کر کے بتاتا ہے کہ بعد میں `ccam login work` چلائیں۔
6. ccam، `profiles/work/.claude.json` سے `oauthAccount` پڑھتا ہے (`profile.ReadIdentity`): ای میل، اکاؤنٹ اور تنظیم کی IDs، تنظیم کا نام اور پلان (`claude_max`، `claude_team` …)۔ یہ فائل کئی میگا بائٹ کی ہو سکتی ہے، اس لیے ccam پوری فائل پارس کرنے کے بجائے صرف اسی ایک key کو تلاش کرتا ہے (`TopLevelValue`)۔ ccam یہ تفصیلات `config.toml` میں `[profile.login]` میں محفوظ کرتا ہے۔
7. ccam، `~/.local/bin/ccwork` کو ccam بائنری سے لنک کرتا ہے، جو شروع ہوتے وقت `argv[0]` کو `config.toml` میں تلاش کرتی ہے، اور اس لنک کو `~/.local/state/ccam/shortcuts.json` میں ریکارڈ کرتا ہے (`shim.Install`)۔ اگر لنک نہ بن سکے تو اکاؤنٹ پھر بھی شامل ہو جاتا ہے، ccam وجہ بتاتا ہے (`Could not create shortcuts: …`) اور خلاصہ `shortcut: ccwork (not created; see above)` دکھاتا ہے۔

لاگ اِن میں کئی منٹ لگ سکتے ہیں، اور اس دوران کوئی دوسرا ccam پروسیس کنفیگ محفوظ کر سکتا ہے۔ اس لیے انتظار کرنے والی کمانڈز `config.toml` دوبارہ لوڈ کرتی ہیں، صرف ایک تبدیلی (یہاں، ای میل) لاگو کرتی ہیں اور محفوظ کرتی ہیں (`env.update`)۔ اسی وجہ سے `ccam remove` اپنے سوال کے بعد دوبارہ لوڈ کرتا ہے۔

`ccam login` کسی موجودہ اکاؤنٹ کے لیے مراحل 5 اور 6 چلاتا ہے۔ Claude Code ای میل کو صرف لاگ اِن صفحے کے لیے ایک اشارے کے طور پر لیتا ہے، اس لیے ہو سکتا ہے آپ پہلے سے مختلف کسی Claude اکاؤنٹ میں سائن اِن ہو جائیں؛ ccam نئے لاگ اِن کا `[profile.login]` سے موازنہ کرتا ہے اور یہ بتا دیتا ہے (`ccam: this login is other@company.com, but work was dev@company.com before.`)۔ یہ نیا لاگ اِن رکھتا ہے اور بتاتا ہے کہ اسے کیسے واپس کریں: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

کسی دوسرے Claude اکاؤنٹ میں لاگ اِن اکاؤنٹ کا [حد کا ریکارڈ](/features/usage#the-limit-record) بھی ہٹا دیتا ہے، کیونکہ وہ حدیں دوسرے اکاؤنٹ کی تھیں۔ جب کوئی دوسرا ccam اکاؤنٹ اسی Claude اکاؤنٹ میں لاگ اِن ہو تو ccam خبردار کرتا ہے کہ دونوں اپنی استعمال کی حدیں شیئر کرتے ہیں۔

`ccam login --expired` ایک کے بعد ایک ہر اس اکاؤنٹ میں لاگ اِن کرتا ہے جس میں کبھی لاگ اِن نہیں ہوا یا جس کا لاگ اِن ختم ہو چکا ہو: اس کا ریفریش ٹوکن ختم ہو گیا ہو، یا سرور کے انکار کے بعد Claude Code نے لاگ اِن مٹا دیا ہو۔ پہلا براؤزر کھلنے سے پہلے یہ جانچتا ہے کہ `config.toml` محفوظ ہو سکتی ہے۔

`ccam logout` اکاؤنٹ میں `claude auth logout` چلاتا ہے اور محفوظ ای میل صاف کر دیتا ہے؛ لاگ اِن کا ریکارڈ باقی رہتا ہے، اس لیے اگلا لاگ اِن اس سے جانچا جاتا ہے۔ Claude Code لاگ اِن منسوخ کر دیتا ہے اور اکاؤنٹ کے لیے جو کچھ محفوظ رکھتا ہے وہ حذف کر دیتا ہے، بشمول MCP سرورز کے سائن اِن۔ اس لیے ccam ان صورتوں میں آپ کو خبردار کرتا ہے اور پہلے پوچھتا ہے (`--yes` سوال چھوڑ دیتا ہے):

- لاگ اِن میں MCP سائن اِن ہوں؛
- اکاؤنٹ چل رہا ہو (اس کے سیشنز بھی سائن آؤٹ ہو جاتے ہیں)؛
- کسی دوسرے اکاؤنٹ کے فولڈر میں وہی لاگ اِن ہو (کاپی کیا گیا فولڈر)، جسے لاگ آؤٹ اسی طرح منسوخ کر دیتا ہے۔

مثال کے طور پر:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

کسی دوسرے Claude اکاؤنٹ کے طور پر لاگ اِن کرنے کے لیے صرف `ccam login <name>` کافی ہے، اور یہ MCP سائن اِن برقرار رکھتا ہے۔

جب اکاؤنٹس ایک سے زیادہ تنظیموں کے ہوں تو `ccam list` ایک `Organization` کالم شامل کرتا ہے (ذاتی پلان `personal` دکھاتے ہیں)، اور `ccam list --json` میں ہر اکاؤنٹ کے لیے `plan` اور `org` ہوتے ہیں۔

### ریموٹ مشین پر سائن اِن

SSH پر `claude auth login` کا براؤزر والا حصہ ٹوٹ سکتا ہے۔ Claude Code ریموٹ مشین پر ایک callback سرور (`http://localhost:<port>/callback`) پر سائن اِن کا انتظار کرتا ہے، لیکن جو براؤزر سائن اِن مکمل کرتا ہے وہ آپ کے کمپیوٹر پر چلتا ہے، جہاں یہ پتا کہیں نہیں جاتا اور براؤزر `ERR_CONNECTION_REFUSED` جیسی غلطی دکھاتا ہے۔ اس پتے میں پھر بھی ایک بار استعمال ہونے والا کوڈ ہوتا ہے، اس لیے ccam اسے پہنچا سکتا ہے:

<!-- motion:accounts-relay -->

ccam آپ کا ٹرمینل Claude Code کے حوالے نہیں کرتا (`runLogin`)۔ یہ Claude Code کے اسٹینڈرڈ اِن پٹ کو ایک pipe سے جوڑتا ہے اور آپ کی ٹائپ کی ہوئی ہر لائن اس میں کاپی کرتا ہے (`relayLogin`)، اس لیے Claude Code کے `Paste code here` پرامپٹ پر چسپاں کیا گیا کوڈ پھر بھی پہنچ جاتا ہے۔ جب کوئی لائن callback پتا ہو تو ccam اسے صرف تب قبول کرتا ہے جب ہوسٹ `localhost`، `127.0.0.1` یا `[::1]` ہو، پورٹ دیا گیا ہو، راستہ `/callback` ہو اور `code` اور `state` دونوں موجود ہوں؛ `http://` چھوڑا جا سکتا ہے (`callbackURL`)۔ پھر ccam سرور پر اس پتے کی درخواست 10 سیکنڈ کے ٹائم آؤٹ کے ساتھ اور ری ڈائریکٹس پر عمل کیے بغیر کرتا ہے، جیسے براؤزر کرتا (`deliverCallback`)۔ درخواست انتظار کرنے والے `claude auth login` کے callback سرور تک پہنچتی ہے (ڈایاگرام میں ڈیش والی لائن)۔ Claude Code کوڈ وصول کر کے لاگ اِن محفوظ کرتا ہے، Linux پر اکاؤنٹ فولڈر کی `.credentials.json` میں۔ چسپاں کی گئی کوئی لائن کبھی ccam کو کسی دوسرے ہوسٹ سے رابطہ نہیں کروا سکتی۔

لاگ اِن شروع ہونے سے پہلے ccam یہ پرنٹ کرتا ہے:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

کامیابی پر یہ `→ Passed the sign-in to Claude Code on this machine.` پرنٹ کرتا ہے۔

جب لاگ اِن ٹرمینل پر اِن پٹ کا انتظار کر رہا ہو تو ccam ٹرمینل کی echo اور لائن موڈ بند کر دیتا ہے اور آپ کے ٹائپ کیے ہوئے ہر حرف کے لیے ایک `*` دکھاتا ہے، اس لیے چسپاں کیا گیا کوڈ یا callback پتا اسکرین پر نظر نہیں آتا۔ ccam لائن کی ایڈیٹنگ خود کرتا ہے (`readLine`): `Enter` لائن بھیجتا ہے، `Backspace` ایک حرف مٹاتا ہے، اور خالی لائن پر `Ctrl+D` اِن پٹ ختم کرتا ہے۔ لاگ اِن ختم ہونے پر ٹرمینل بحال ہو جاتا ہے، اور `Ctrl+C`، `SIGTERM` یا `SIGHUP` پر بھی، جو اب بھی ccam کو ختم کرتے ہیں۔ اگر اسٹینڈرڈ اِن پٹ اور اسٹینڈرڈ ایرر پر ٹرمینل نہ ہو تو لائنیں جیسے آتی ہیں ویسے ہی کاپی ہوتی ہیں۔

### نام بدلنا، ہٹانا اور امپورٹ کرنا

`ccam rename` اور `ccam remove` اکاؤنٹ کے فولڈر پر اور ان چیزوں پر عمل کرتے ہیں جو ccam اکاؤنٹ کے نام سے دوسری جگہوں پر رکھتا ہے۔ ccam ان راستوں کی فہرست `nameState` میں رکھتا ہے:

- `~/.local/state/ccam/last-used/<name>`؛
- `~/.local/state/ccam/mcp/<name>`؛
- استعمال کے کیش `~/.local/state/ccam/usage/<name>.json` اور `<name>.statusline.json`؛
- حد کا ریکارڈ `<name>.limit.json`؛
- `~/.local/state/ccam/handoff/<name>/` میں سیشن کے نوٹس۔

`ccam remove work` انہیں اس ترتیب سے سنبھالتا ہے:

<!-- motion:accounts-remove -->

1. `ccam remove` پہلے انکار کرتا ہے اگر اکاؤنٹ چل رہا ہو، پھر پوچھتا ہے، اور آپ کے جواب کے بعد دوبارہ جانچتا ہے (اس سے وہ سیشن بھی پکڑا جاتا ہے جو سوال کھلا رہنے کے دوران شروع ہوا ہو)۔
2. کچھ بھی تبدیل کرنے سے پہلے یہ جانچتا ہے کہ `config.toml` محفوظ ہو سکتی ہے۔
3. اگر اکاؤنٹ ایکسٹینشنز شیئر کرتا ہو تو ccam گروپ کی `installed_plugins.json` اور `known_marketplaces.json` دوبارہ لکھتا ہے (`RewritePluginPaths`)، کیونکہ Claude Code پلگ اِنز کے مکمل راستے محفوظ کرتا ہے جو اس اکاؤنٹ کے فولڈر سے گزر سکتے ہیں۔
4. اس کے بعد منیجڈ اکاؤنٹ (جس کا فولڈر ccam نے بنایا ہو، نہ کہ `ccam import` سے رجسٹر کیا گیا) لاگ آؤٹ کیا جاتا ہے: اس کے فولڈر میں `claude auth logout` چلتا ہے، تاکہ کوئی کارآمد لاگ اِن trash میں انتظار نہ کرے (یا macOS پر فولڈر کے نام والے Keychain آئٹم میں، جو بعد میں اسی نام کا اکاؤنٹ وراثت میں پاتا)۔ API key اکاؤنٹ کے لیے یہ مرحلہ چھوڑ دیا جاتا ہے، اور تب بھی جب کسی دوسرے اکاؤنٹ کے فولڈر میں وہی لاگ اِن ہو، جسے لاگ آؤٹ بھی منسوخ کر دیتا۔ ناکام لاگ آؤٹ کی اطلاع دی جاتی ہے، اور ہٹانے کا عمل جاری رہتا ہے۔
5. ccam فولڈر کو `~/.local/share/ccam/trash/<name>-<time>/` میں منتقل کرتا ہے (`profile.Trash`)، اکاؤنٹ `config.toml` سے نکل جاتا ہے، `nameState` کے راستے حذف ہو جاتے ہیں تاکہ بعد میں اسی نام کا اکاؤنٹ انہیں وراثت میں نہ پائے، اور شارٹ کٹ کا لنک ہٹا دیا جاتا ہے: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`۔

اگر آخر میں `config.toml` محفوظ نہ ہو سکے تو فولڈر trash سے واپس آ جاتا ہے۔ ہٹائے گئے اکاؤنٹ کو واپس لانے کے لیے فولڈر کو trash سے واپس منتقل کریں اور دوبارہ لاگ اِن کریں۔

`ccam rename` فولڈر کا نام بدلتا ہے اور اس کے `plugins/` میں پلگ اِن کے راستے پرانے فولڈر سے نئے فولڈر پر دوبارہ لکھتا ہے۔ یہ `nameState` کے راستوں کو نئے نام پر منتقل کرتا ہے اور اگر اکاؤنٹ ڈیفالٹ تھا تو اسے ڈیفالٹ رکھتا ہے۔ یہ سیشن کے نوٹس بھی اپ ڈیٹ کرتا ہے، تاکہ `ccam next --continue` انہیں پھر بھی ڈھونڈ سکے۔ شارٹ کٹ کمانڈ اپنا نام برقرار رکھتی ہے۔ اگر پلگ اِن کے راستے دوبارہ لکھنا ناکام ہو، یا `config.toml` محفوظ نہ ہو سکے، تو فولڈر کا نام واپس بدل دیا جاتا ہے۔

macOS پر، اور Windows پر جب Claude Code لاگ اِنز Credential Manager میں رکھتا ہو، فولڈر خود اپنا پرانا نام برقرار رکھتا ہے، کیونکہ لاگ اِن اور MCP سائن اِن فولڈر کے راستے کے تحت محفوظ ہوتے ہیں اور منتقل کیا گیا فولڈر لاگ آؤٹ ہو جاتا۔ تب کنفیگ فولڈر کو نئے نام کے ساتھ `folder` کے طور پر ریکارڈ کرتا ہے، اور ccam کہتا ہے `its folder stays …: the macOS Keychain keeps the login under the folder's path`۔

`ccam import` کسی فولڈر کو اس کی جگہ پر ہی، `config.toml` میں `dir` کے طور پر رجسٹر کرتا ہے۔ یہ اسے کبھی منتقل یا کاپی نہیں کرتا: کاپی کیا گیا ریفریش ٹوکن کاپیوں کو ایک دوسرے کو لاگ آؤٹ کروا دے گا۔

فولڈر آرگیومنٹ کے بغیر یہ غیر رجسٹر شدہ `~/.claude` اور `~/.claude-*` فولڈرز تلاش کرتا ہے اور ہر ایک کے بارے میں پوچھتا ہے۔ نام فولڈر سے آتے ہیں (`~/.claude` سے `main`، `~/.claude-work` سے `work`، `~/.claude-2` سے `acct2`)۔ اگر ccam نام اخذ نہ کر سکے، یا نام پہلے سے لیا جا چکا ہو، تو وہ فولڈر کو چلانے کے لیے `--name` کمانڈ کے ساتھ دکھاتا ہے۔ ملنے والا جو فولڈر ناکام ہو اس کی اطلاع `Skipped <path>: <reason>` کے طور پر دی جاتی ہے اور باقی امپورٹ ہو جاتے ہیں۔

ccam پہلے اپنے سارے سوال پوچھتا ہے، پھر `config.toml` دوبارہ لوڈ کر کے فولڈرز رجسٹر کرتا ہے، اس لیے اس دوران دوسری ccam کمانڈز نے جو تبدیلیاں محفوظ کیں وہ برقرار رہتی ہیں۔ `✓` والی لائنیں محفوظ کرنے کے بعد آتی ہیں۔

`CLAUDE_CONFIG_DIR` کے بغیر Claude Code، `.claude.json` کو آپ کے ہوم فولڈر میں `~/.claude` کے ساتھ رکھتا ہے، اس لیے جب آپ `~/.claude` امپورٹ کریں اور اس میں اپنی `.claude.json` نہ ہو تو ccam، `~/.claude/.claude.json` سے `~/.claude.json` کا لنک بناتا ہے۔ تب سادہ `claude` اور `ccam main` ایک ہی اکاؤنٹ دیکھتے ہیں۔

اگر `~/.claude` میں `~/.claude.json` کے علاوہ اپنی ایک `.claude.json` بھی ہو تو ccam دونوں کو رہنے دیتا ہے اور یہ بتا دیتا ہے۔ `ccam main` اندر والی فائل استعمال کرتا اور سادہ `claude` باہر والی، جس سے ان کی لاگ اِن تفصیلات اور MCP سرورز الگ الگ رہتے۔ ایک ہی فائل رکھنے کے لیے اندر والی فائل کو لنک سے بدل دیں ([doctor](/features/doctor) بھی اس کی اطلاع دیتا ہے)۔

جب فولڈر میں وہی لاگ اِن ہو جو کسی رجسٹر شدہ اکاؤنٹ کا ہے (کاپی کیا گیا فولڈر)، تو ccam یہ بتا دیتا ہے: جو بھی پہلے اس کی تجدید کرے وہ دوسرے کو لاگ آؤٹ کر دیتا ہے، اس لیے اسے استعمال کرنے سے پہلے `ccam login <name>` چلائیں۔

امپورٹ کیے گئے اکاؤنٹس کسی گروپ میں نہیں ہوتے جب تک آپ انہیں کسی گروپ میں شامل نہ کریں (`ccam group join <group> <name>`)، اور کوئی شارٹ کٹ نہیں بنتا۔

`ccam completion bash|zsh` ایک completion اسکرپٹ پرنٹ کرتا ہے جو اکاؤنٹس کے نام بھی مکمل کرتا ہے۔

## حفاظت اور خاص صورتیں

- `ccam rename` اور `ccam remove` جب اکاؤنٹ چل رہا ہو (`refuseRunning`) تو ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)`` کے ساتھ انکار کرتے ہیں، کیونکہ سیشن فولڈر میں لکھتا رہتا ہے۔ یہ جانچ Claude Code کو پکڑتی ہے چاہے وہ کسی بھی طرح شروع ہوا ہو: ccam کے ذریعے، سادہ `claude`، کوئی IDE، بیک گراؤنڈ سیشنز اور ان کا daemon۔
- جب کوئی دوسرا ccam پروسیس اس لاگ اِن کی تجدید یا تبدیلی کر رہا ہو تو `ccam rename`، `ccam remove` اور `ccam logout` انکار کرتے ہیں: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``۔
- جب کسی دوسرے ccam پروسیس میں گروپ کی تبدیلی جاری ہو تو `ccam rename` اور `ccam remove` `another ccam process is changing the groups right now; try again when it is done` کے ساتھ انکار کرتے ہیں؛ اور جب کوئی تبدیلی ادھوری رک گئی ہو تو ``the last group change stopped halfway; finish it with `ccam doctor --fix` first`` کے ساتھ۔
- امپورٹ کیے گئے اکاؤنٹ پر `ccam remove` فولڈر اور اس کا لاگ اِن سادہ `claude` کے لیے برقرار رکھتا ہے۔ اگر اکاؤنٹ گروپس میں ہو تو یہ پہلے انہیں اسی طرح چھوڑتا ہے جیسے `ccam group leave` کرتا ہے: ہر شیئر کیا گیا آئٹم فولڈر میں گروپ والے آئٹم کی کاپی بن جاتا ہے، اور گروپ کے MCP سرورز اس کی `.claude.json` میں چلے جاتے ہیں۔ پھر اس کی رجسٹریشن ختم کی جاتی ہے: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`۔ اس کی `.claude.json` کے جو بیک اپ ccam نے بنائے تھے وہ `<folder>/.ccam-backup/mcp-<time>` میں منتقل ہو جاتے ہیں؛ اگر وہ منتقل نہ ہو سکیں تو `~/.local/state/ccam/mcp/<name>/backups` میں رہتے ہیں اور پیغام یہ بتا دیتا ہے۔ آپ کا `~/.claude`، trash میں منتقل نہیں ہوتا۔
- دو اکاؤنٹس کا ایک ہی فولڈر نہیں ہو سکتا: جو `config.toml` ایک فولڈر کا نام دو بار لے وہ لوڈ نہیں ہوتی۔
- `ccam add` ایسا نام مسترد کرتا ہے جس کا فولڈر `profiles/` میں پہلے سے موجود ہو، اور `ccam import` کی طرف رہنمائی کرتا ہے۔ macOS پر یہ تب بھی رک جاتا ہے جب Keychain میں نئے فولڈر کے راستے کے لیے ccam کے کسی پرانے ورژن کا چھوڑا ہوا لاگ اِن یا MCP سائن اِن اب بھی موجود ہوں، اور وہ `security delete-generic-password` کمانڈ بتاتا ہے جو انہیں ہٹاتی ہے۔ اگر `config.toml` محفوظ نہ ہو سکے تو نیا فولڈر دوبارہ ہٹا دیا جاتا ہے۔
- `ccam import` پہلے سے رجسٹر شدہ فولڈر، ایسا راستہ جو فولڈر نہ ہو، اور پہلے سے لیا گیا نام مسترد کرتا ہے۔ جب فولڈر `~/.claude` کی کاپی ہو (وہی لاگ اِن) تو ccam خبردار بھی کرتا ہے، چاہے `~/.claude` خود رجسٹر نہ ہو۔
- جب اکاؤنٹ چل رہا ہو تو `ccam logout` پہلے پوچھتا ہے؛ `ccam login` یہ نہیں جانچتا۔
- اگر ریلے Claude Code تک نہ پہنچ سکے، مثلاً لاگ اِن ختم ہونے کے بعد، تو یہ وجہ کے ساتھ `ccam: could not pass the sign-in to Claude Code:` پرنٹ کرتا ہے۔ HTTP غلطی کا عام طور پر مطلب ہے کہ کوڈ کی میعاد ختم ہو گئی؛ لاگ اِن دوبارہ شروع کریں۔

## متعلقہ

- [Claude Code شروع کرنا](/features/launch)
- [شارٹ کٹ کمانڈز](/features/shortcuts)
- [MCP سرورز، سیٹنگز اور ہسٹری شیئر کریں](/docs/sharing)
- [ایک کمپیوٹر پر Claude Code کے کئی اکاؤنٹس چلائیں](/guides/multiple-claude-code-accounts)
