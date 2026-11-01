---
title: "Claude Code شروع کرنا"
description: "ccam <name|number>، ccam run یا ڈیش بورڈ سے Claude Code کسی بھی اکاؤنٹ سے شروع کریں؛ آرگیومنٹس جوں کے توں claude کو جاتے ہیں اور --dry-run کمانڈ دکھاتا ہے۔"
updated: 2026-11-01
order: 1
translation: machine
sourceHash: 11df08608a1e
---

`ccam work`، Claude Code کو `work` اکاؤنٹ کے طور پر شروع کرتا ہے، اور نام کے بعد سب کچھ بغیر تبدیلی کے `claude` کو جاتا ہے۔ اسے ہر اس جگہ استعمال کریں جہاں آپ `claude` ٹائپ کرتے: ٹرمینل میں، اسکرپٹ میں، یا pipe سے آنے والے اِن پٹ کے ساتھ۔

## استعمال

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

اکاؤنٹ کا نام یا `ccam list` سے اس کا نمبر دیں۔ ccam نام کے بعد والے آرگیومنٹس نہیں پڑھتا، اس لیے `ccam work --help`، Claude Code کی مدد دکھاتا ہے۔ نام کے بعد `--` اختیاری ہے۔ اسٹینڈرڈ اِن پٹ، ایگزٹ کوڈ اور سگنلز ویسے ہی کام کرتے ہیں جیسے `claude` کے ساتھ۔

جب پہلا آرگیومنٹ `-` سے شروع ہو اور ccam کا اپنا کوئی فلیگ نہ ہو، تو ccam آپ کے لیے اکاؤنٹ چنتا ہے اور تمام آرگیومنٹس Claude Code کو دیتا ہے۔ یہ وہ اکاؤنٹ استعمال کرتا ہے جو [`ccam folder`](/features/folder) سے موجودہ فولڈر کے لیے مقرر ہو، ورنہ ڈیفالٹ اکاؤنٹ: وہ جو `ccam use` سے مقرر کیا گیا ہو، یا اگر آپ نے کوئی مقرر نہ کیا ہو تو سب سے کم نمبر والا اکاؤنٹ۔

### ڈیش بورڈ

بغیر آرگیومنٹس کے `ccam` [ڈیش بورڈ](/features/dashboard) کھولتا ہے جب اسٹینڈرڈ اِن پٹ اور آؤٹ پٹ ٹرمینل ہوں، ورنہ `ccam list` پرنٹ کرتا ہے۔ کسی اکاؤنٹ پر جائیں (ایرو کیز سے، یا اس کا نمبر ٹائپ کر کے) اور اسے شروع کرنے کے لیے `Enter` دبائیں؛ `n` سب سے زیادہ گنجائش والا اکاؤنٹ شروع کرتا ہے اور `o`، `--continue`، `--resume` یا دوسرے آرگیومنٹس پیش کرتا ہے۔ ڈیش بورڈ Claude Code شروع ہونے سے پہلے بند ہو جاتا ہے، اس لیے ٹرمینل پہلے بحال ہو جاتا ہے۔

**Team اور Enterprise سیٹیں۔** خودکار اکاؤنٹ انتخاب (`ccam next`، `--fallback`، `ccam pick`، ڈیش بورڈ کی `n` کلید اور، جب آپ اسے آن کریں تو، [`auto_switch`](/features/next#auto-switch)) ایک ہی Claude Team یا Enterprise تنظیم کی سیٹوں کے درمیان بھی اکاؤنٹ بدلتا ہے، جس سے Anthropic کی شرائط کی خلاف ورزی ہو سکتی ہے۔ اس پر انحصار کرنے سے پہلے [عام سوالات](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) پڑھیں اور اپنی تنظیم اور Anthropic سے تصدیق کر لیں۔

جب شیل میں `CLAUDE_CODE_OAUTH_TOKEN`، `ANTHROPIC_API_KEY` یا `ANTHROPIC_AUTH_TOKEN` سیٹ ہو تو Claude Code اکاؤنٹ کے اپنے لاگ اِن کے بجائے وہ key یا ٹوکن استعمال کر سکتا ہے، اس لیے ccam شروع کرنے سے پہلے اسٹینڈرڈ ایرر پر آپ کو خبردار کرتا ہے (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`)۔ یہ صرف ویری ایبلز پڑھتا ہے؛ [API key اکاؤنٹ](/features/api) اس کے بجائے انہیں صاف کر دیتا ہے۔

ڈیفالٹ طور پر ccam خود کو Claude Code سے بدل دیتا ہے، اس لیے اور کچھ نہیں چلتا رہتا۔ `config.toml` میں سب سے اوپر، پہلے `[[group]]`، `[[profile]]` یا `[[folder]]` سے اوپر، `relay = true` (یا Linux اور macOS پر [`auto_switch = true`](/features/next#auto-switch)) کے ساتھ ccam اس کے بجائے انٹرایکٹو سیشنز کو چائلڈ پروسیس کے طور پر شروع کرتا ہے اور ان کے ساتھ ساتھ چلتا رہتا ہے۔ اس سے `/ccam:switch`، یا استعمال کی حد پر `auto_switch`، گفتگو کو خود کسی دوسرے اکاؤنٹ پر جاری رکھ سکتا ہے، سیشن کے آپشنز، اجازتوں کے موڈ، effort لیول اور ایجنٹ کے ساتھ؛ اس سے پہلے کہ `/ccam:switch` اسے کسی دوسری تنظیم کے اکاؤنٹ پر منتقل کرے، ccam آپ سے پوچھتا ہے ([تفصیلات](/features/next#relay))۔

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `--dry-run` | وہ انوائرنمنٹ اور کمانڈ پرنٹ کرتا ہے جو ccam چلاتا، پھر باہر نکل جاتا ہے | بند |
| `--fallback` | `-p`/`--print` رنز کے لیے: جب اکاؤنٹ اپنے استعمال کی حد کو پہنچ جائے تو وہی کمانڈ سب سے زیادہ گنجائش والے لاگ اِن اکاؤنٹ سے دوبارہ چلاتا ہے ([تفصیل](/features/next#run-again-at-the-limit)) | بند |
| `-v`، `--version` | `ccam <version>` پرنٹ کرتا ہے (صرف پہلے آرگیومنٹ کے طور پر) | |
| `-h`، `--help` | ccam کی مدد، `ccam --help` یا `ccam run --help` کے طور پر | |

آرگیومنٹس کے اصول:

- ccam اپنے فلیگز صرف اکاؤنٹ کے نام سے پہلے پڑھتا ہے۔ نام کے بعد ہر آرگیومنٹ Claude Code کا ہے، بشمول `--dry-run` اور `--help`۔
- `--dry-run`، `ccam --dry-run work`، `ccam run --dry-run work` اور `ccam --dry-run` (فولڈر کے اکاؤنٹ یا ڈیفالٹ کے لیے) کے طور پر کام کرتا ہے۔ `--dry-run` اور `--fallback` کسی بھی ترتیب میں آ سکتے ہیں۔
- `ccam run <name|number> [claude args]`، `ccam <name>` کی لمبی شکل ہے۔ نام کے بغیر `ccam run` فولڈر کا اکاؤنٹ یا ڈیفالٹ اکاؤنٹ استعمال کرتا ہے۔
- اکاؤنٹس کے نام `add`، `list`، `run` یا `usage` جیسے ccam کمانڈز کے نام نہیں ہو سکتے، اس لیے `ccam <word>` کبھی مبہم نہیں ہوتا۔

ڈرائی رن ایک لائن پرنٹ کرتا ہے جسے آپ شیل میں، یا Windows پر PowerShell میں، چسپاں کر سکتے ہیں۔ یہ اس سے شروع ہوتی ہے جو لانچ انوائرنمنٹ میں بدلتا ہے: ہمیشہ `CLAUDE_CONFIG_DIR` اور `CCAM_PROFILE`، ہٹائے جانے والے ویری ایبل کے لیے `env -u <name>` (جیسے وراثت میں ملا `CCAM_RELAY`)، اور [API key اکاؤنٹ](/features/api) کے لیے اس کا اینڈ پوائنٹ، اس کا ماڈل اور یہ کہ key کہاں سے آتی ہے (`ANTHROPIC_API_KEY="$WORK_KEY"`)، خود key کبھی نہیں:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

PowerShell والی لائن بعد میں ویری ایبلز کو پہلے جیسا کر دیتی ہے، جیسے شیل میں `VAR=x cmd` کرتا ہے، اور Claude Code کو اس کی کمانڈ لائن `Start-Process` کے ذریعے دیتی ہے، اس لیے Windows PowerShell 5.1 اور PowerShell 7 اسے ایک ہی طرح چلاتے ہیں۔ جب لانچ کسی ویری ایبل کو خالی قدر پر مقرر کرے تو اس کے نیچے ایک `#` لائن بتاتی ہے کہ 7.5 سے پہلے کا PowerShell ایسے ویری ایبل کو اس کے بجائے ہٹا دیتا ہے۔ npm سے انسٹال کیے گئے `claude.cmd` کے لیے، لائن وہی `cmd.exe` کمانڈ چلاتی ہے جو لانچ استعمال کرتا (نیچے دیکھیں)، اور جو آرگیومنٹ لانچ وہاں مسترد کرتا وہ ڈرائی رن کو بھی ناکام کر دیتا ہے۔

## یہ کیسے کام کرتا ہے

ہر لانچ کمانڈ لائن سے `exec` تک وہی مراحل اٹھاتا ہے:

<!-- motion:launch -->

ہر خانہ اپنا مرحلہ چلنے پر روشن ہوتا ہے، اور نیچے کی لائنیں وہ انوائرنمنٹ اور آرگیومنٹس ہیں جن کے ساتھ Claude Code شروع ہوتا ہے۔

1. ccam، `app.Main` میں شروع ہوتا ہے، جسے `cmd/ccam/main.go` کال کرتا ہے۔ اگر ccam کسی اور نام سے شروع ہوا ہو تو یہ ایک [شارٹ کٹ](/features/shortcuts) ہے اور سیدھا لانچر پر جاتا ہے۔
2. ccam، cobra کے شامل ہونے سے پہلے `os.Args` دیکھتا ہے (`app.parse`)۔ صرف اکیلا `ccam` (ڈیش بورڈ)، ccam کمانڈ کا نام (`config.Reserved`)، `-h`/`--help` اور `-v`/`--version`، `cli.Execute` پر جاتے ہیں۔ باقی سب کچھ لانچ ہے۔
3. ccam، `~/.config/ccam/config.toml` پڑھتا ہے، جو واحد فائل ہے جسے یہ پارس کرتا ہے (`launcher.Launch`)۔
4. ccam اکاؤنٹ کو نام یا نمبر سے تلاش کرتا ہے، شروع کا ایک `--` ہٹاتا ہے، اور `claude` قابلِ عمل فائل تلاش کرتا ہے (`launcher.Prepare`، `ResolveClaude`)۔ نام نہ ہو تو یہ فولڈر کا اکاؤنٹ لیتا ہے، ورکنگ ڈائریکٹری کو `config.toml` کے `[[folder]]` اندراجات سے اسٹرنگ کے طور پر ملا کر، macOS اور Windows پر بڑے چھوٹے حروف کا لحاظ کیے بغیر (`HereProfile`)۔ PATH کی تلاش فائل کے ترمیم کے وقت اور `PATH` کے ایک hash کے ساتھ `~/.local/state/ccam/claude-path` میں کیش ہوتی ہے، اس لیے ccam ہر لانچ پر `PATH` تلاش نہیں کرتا۔ بدلی ہوئی یا غائب فائل، یا بدلا ہوا `PATH`، نئی تلاش شروع کرتا ہے۔ `config.toml` میں `claude_bin` کا راستہ (`~/` پھیلا دیا جاتا ہے) ایک قابلِ عمل فائل ہونا چاہیے؛ ورنہ کچھ بھی لکھے جانے سے پہلے لانچ `claude_bin X is not an executable file. Check claude_bin in config.toml` کے ساتھ رک جاتا ہے۔ ccam ایسا `claude` مسترد کرتا ہے جو خود ccam ہو، inode سے موازنہ کر کے، کیونکہ اس سے لوپ بن جائے گا۔
5. ccam اکاؤنٹ فولڈر `~/.local/share/ccam/profiles/<name>/` کا اس کی شیئرنگ سیٹنگز سے موازنہ کرتا ہے، صرف `lstat` اور `readlink` کال کرتے ہوئے، ہر آئٹم کے لیے ایک یا دو سسٹم کالز (`checkLinks`، `share.Check`)۔ غائب لنک فوراً بنا دیا جاتا ہے۔ جس چیز سے ڈیٹا کھو سکتا ہو وہ لانچ روک دیتی ہے۔
6. ccam شیئر کیے گئے MCP سرورز شامل کرتا ہے (`applyMCP`؛ نیچے دیکھیں)۔
7. ccam سیشن ریکارڈ کرتا ہے (`session.Register`)۔ یہ اپنے PID کے نام سے ایک نشان، جس میں پروسیس کے شروع ہونے کا وقت ہوتا ہے، `~/.local/state/ccam/locks/profile-<name>/<pid>` میں لکھتا ہے، `XDG_RUNTIME_DIR` کچھ بھی کہے، تاکہ ہر شیل وہی نشان دیکھے۔ گفتگو کی ہسٹری شیئر کرنے والے اکاؤنٹ کے لیے، ccam پہلے `history-<group>-<hash>` میں کوئی چلتا ہوا سیشن تلاش کرتا ہے، جہاں hash موجودہ فولڈر کا ہوتا ہے، symbolic links حل کر کے، اور وہاں بھی رجسٹر کرتا ہے (`historyCheck`)۔

   `exec`، PID برقرار رکھتا ہے، اس لیے نشان پھر خود Claude Code پروسیس کا نام بتاتا ہے۔ ccam پکڑے ہوئے لاک کے بجائے نشان والی فائل استعمال کرتا ہے، کیونکہ Claude Code جو پروسیس شروع کرتا ہے وہ لاک وراثت میں پا کر سیشن ختم ہونے کے بعد بھی اسے رکھے رہیں گے۔ جو نشان لکھا نہ جا سکے وہ لانچ نہیں روکتا: ccam کہتا ہے `ccam: could not mark the session as running (…); starting Claude Code anyway`۔

   یہ جاننے کے لیے کہ کوئی اکاؤنٹ چل رہا ہے یا نہیں، ccam کی دوسری جانچیں (نام بدلنا، ہٹانا، گروپس کی تبدیلیاں، `doctor --fix`، MCP sync) اکاؤنٹ فولڈر میں Claude Code کے اپنے ریکارڈ بھی پڑھتی ہیں (`sessions/<pid>.json`، `daemon.lock`)، اس لیے وہ ccam کے بغیر شروع ہونے والا Claude Code بھی دیکھتی ہیں: سادہ `claude`، کوئی IDE، بیک گراؤنڈ سیشنز اور ان کا daemon۔
8. ccam ٹرمینل کا عنوان `claude · <name>` مقرر کرتا ہے (جب اسٹینڈرڈ ایرر ٹرمینل ہو)، لانچ کا وقت `fsync` کے بغیر `~/.local/state/ccam/last-used/<name>` میں لکھتا ہے، اور `syscall.Exec` کال کرتا ہے (`Plan.Exec`)۔

Claude Code اس طرح چلتا ہے کہ `CLAUDE_CONFIG_DIR` اکاؤنٹ فولڈر پر اور `CCAM_PROFILE` اکاؤنٹ کے نام پر مقرر ہوتا ہے۔ ccam ان ویری ایبلز کی وراثت میں ملی قدریں صاف کر دیتا ہے، اور وراثت میں ملا `CCAM_RELAY` بھی۔ macOS پر `ccam import` سے رجسٹر کیے گئے `~/.claude` فولڈر کے لیے یہ `CLAUDE_SECURESTORAGE_CONFIG_DIR` کو بھی خالی قدر پر مقرر کرتا ہے، تاکہ اکاؤنٹ اور سادہ `claude` ایک ہی Keychain لاگ اِن استعمال کریں؛ Windows پر یہ یہی کرتا ہے جب Claude Code لاگ اِن Credential Manager میں رکھتا ہو۔

### exec کیوں

ccam اپنے پروسیس کو `claude` سے بدل دیتا ہے (`syscall.Exec`)۔ اس کے نتیجے میں:

- ccam میموری میں نہیں رہتا؛
- کی اسٹروکس، سگنلز اور ایگزٹ کوڈ Claude Code کے اپنے ہوتے ہیں؛
- آرگیومنٹس ایک array رہتے ہیں، کبھی شیل اسٹرنگ نہیں بنتے۔

Windows میں `exec` نہیں ہوتا۔ وہاں ccam، `claude.exe` کو اسی انوائرنمنٹ، آرگیومنٹس اور کنسول کے ساتھ چائلڈ پروسیس کے طور پر شروع کرتا ہے، اس کا انتظار کرتا ہے اور اسی کے ایگزٹ کوڈ کے ساتھ باہر نکلتا ہے (`Plan.Exec`)۔ npm کا انسٹال کیا ہوا `claude.cmd`، `cmd.exe` میں اپنے پورے راستے کے ساتھ چلتا ہے اور ہر آرگیومنٹ quotes میں ہوتا ہے تاکہ `cmd.exe` ان میں سے کسی کو اپنا syntax نہ سمجھے (`&`، `%PATH%`)؛ لائن بریک والا آرگیومنٹ اس طرح نہیں دیا جا سکتا اور مسترد کر دیا جاتا ہے، اس مشورے کے ساتھ کہ native Claude Code انسٹال کریں۔

Ctrl+C کنسول سے سیدھا Claude Code تک پہنچتا ہے، اس لیے ccam اسے نظر انداز کرتا ہے۔ سیشن مارکر میں ccam کا PID ہوتا ہے، جو ٹھیک اتنی دیر زندہ رہتا ہے جتنی دیر Claude Code، اور Claude Code کے باہر نکلنے پر ہٹا دیا جاتا ہے۔ انتظار کے دوران ccam اور کچھ نہیں کرتا۔

سیشن کے بعد کا کام، جیسے `claude mcp add` سے شامل کیے گئے سرورز کو نوٹ کرنا، اگلے لانچ پر ہوتا ہے: ccam، `.claude.json` کے ترمیم کے وقت کو `~/.local/state/ccam/mcp/<name>/seen` سے موازنہ کرتا ہے اور فائل صرف تب پارس کرتا ہے جب وہ بدلی ہو۔

### کمانڈ لائن کس راستے پر جاتی ہے

تین کمانڈ لائنیں، `internal/app` سے گزرنے کے تین راستے:

<!-- motion:launch-route -->

شارٹ کٹ `argv[0]` پر، اور `ccam <name>` پہلے آرگیومنٹ پر پہچانا جاتا ہے، اور دونوں `syscall.Exec` پر ختم ہوتے ہیں۔ صرف `ccam usage` جیسی ccam کمانڈ cobra تک پہنچتی ہے۔

### شیئر کیے گئے MCP سرورز

ڈیفالٹ `flag` حکمتِ عملی کے ساتھ، ccam گروپ کے MCP سرورز Claude Code کو ایک آرگیومنٹ کے طور پر دیتا ہے، جو آپ کے آرگیومنٹس سے پہلے رکھا جاتا ہے:

<!-- motion:launch-mcp -->

ccam یہ آرگیومنٹ صرف تب شامل کرتا ہے جب گروپ کی `mcp.json` موجود ہو، اور اسے `claude` کے فوراً بعد رکھتا ہے (`mcp.FlagArgs`، `Plan.InsertArgs`)۔ آپ کا اپنا `--mcp-config` یا `--strict-mcp-config` پھر بھی اس کے اوپر لاگو ہوتا ہے۔ یہ ایک ہی آرگیومنٹ ہے، `--mcp-config=<path>`، کیونکہ Claude Code کا `--mcp-config` کئی قدریں لیتا ہے: دو آرگیومنٹس کے طور پر لکھا جائے تو یہ بعد میں آنے والے پرامپٹ کو بھی کنفیگ فائل سمجھ کر پڑھے گا۔

ccam طے کرتا ہے کہ فلیگ کہاں نہیں جا سکتا (`claudeVerb`)۔ یہ آرگیومنٹس سے Claude Code کی سب کمانڈ اسی طرح پڑھتا ہے جیسے Claude Code پڑھتا ہے، آپشنز کو ان کی قدروں سمیت چھوڑتے ہوئے۔ `remote-control` یا `rc` کو کوئی فلیگ نہیں ملتا (نیچے)، اور نہ ہی پہلے دیے گئے کسی بیک گراؤنڈ سیشن والے فعل کو (`attach`، `logs`، `stop`، `kill`، `respawn`، `rm`)۔ `ccam work mcp list`، `get` اور `remove` کو بھی نہیں ملتا، کیونکہ Claude Code وہاں اسے مسترد کرتا ہے؛ اس کے بجائے ccam اسٹینڈرڈ ایرر پر بتاتا ہے کہ گروپ کے سرورز اکاؤنٹ کے سیشنز میں بھی لوڈ ہوتے ہیں:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

باقی ہر سب کمانڈ اور ہر سیشن کو فلیگ ملتا رہتا ہے۔

`config.toml` میں `mcp_strategy = "sync"` کے ساتھ، ccam اس کے بجائے سرورز اکاؤنٹ کی `.claude.json` میں لکھتا ہے، اور اگر اکاؤنٹ پہلے سے چل رہا ہو تو یہ چھوڑ دیتا ہے۔ جب پچھلی جانچ کے بعد نہ گروپ کی `mcp.json` بدلی ہو نہ `.claude.json`، تو لانچ یہ کام بھی چھوڑ دیتا ہے۔ اگر اکاؤنٹ کا اپنا کوئی سرور گروپ کے کسی سرور کے نام سے ہو تو sync اکاؤنٹ کا سرور رکھتا ہے اور ایک بار یہ بتا دیتا ہے۔ sync، `.claude.json` سے غائب ہو جانے والی کاپی بھی واپس رکھ دیتا ہے (`put back group default's MCP servers …`)، اور صرف وہی کاپیاں بدلتا یا ہٹاتا ہے جو اس نے لکھیں اور جنہیں بعد میں کسی نے نہیں بدلا۔

### Remote Control

`ccam work rc` (یا `ccam work remote-control`)، Claude Code کا Remote Control اکاؤنٹ `work` کے طور پر شروع کرتا ہے۔ `claude remote-control` اپنے فعل (verb) سے پہلے `--mcp-config` کو مسترد کرتا ہے اور اس کے بعد کوئی نہیں لیتا، اور اس کے سیشنز MCP سرورز صرف اکاؤنٹ کی `.claude.json` سے پڑھتے ہیں۔ اس لیے یہاں فلیگ والی حکمتِ عملی مختلف طریقے سے کام کرتی ہے:

<!-- motion:launch-rc -->

1. ccam، `remote-control` یا `rc` کو Claude Code کے فعل کے طور پر پہچانتا ہے، آپشنز اور ان کی قدروں کے بعد بھی (`ccam work --model opus rc`)، اور کوئی فلیگ شامل نہیں کرتا (`claudeVerb`، `mcp.FlagArgs`)۔
2. ccam گروپ کی `mcp.json` کے ترمیم کے وقت کا موازنہ اکاؤنٹ کے `sync-stamp` سے کرتا ہے، اور جانچتا ہے کہ گروپ کا ہر سرور اب بھی `.claude.json` میں ہے (`mcp.Synced`)۔ اگر `.claude.json` میں پہلے سے گروپ کے موجودہ سرورز ہوں، یا گروپ کے پاس کوئی سرور نہ ہو، تو Remote Control بغیر سوال کے شروع ہو جاتا ہے۔
3. ورنہ ccam ان سرورز کے نام بتاتا ہے جو سیشنز کو نہیں ملیں گے اور پوچھتا ہے:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. `y` پر، ccam انہیں اسی طرح کاپی کرتا ہے جیسے `sync` حکمتِ عملی کرتی ہے (`mcp.Sync`): Claude Code کے لاک کے تحت، بیک اپ کے بعد، اور اکاؤنٹ کے اپنے سرورز برقرار رکھتے ہوئے۔ `N` پر، یا ٹرمینل کے بغیر، Remote Control ان کے بغیر شروع ہوتا ہے، اور ccam، `mcp_strategy = "sync"` کی طرف اشارہ کرتا ہے، جو ہر لانچ سے پہلے انہیں کاپی کرتی ہے۔

ccam ہر Remote Control لانچ پر پوچھتا ہے جب تک سرورز کاپی نہ ہو جائیں۔ کاپی کے بعد یہ دوبارہ صرف تب پوچھتا ہے جب گروپ کی `mcp.json` بدلے یا کاپی کیا گیا کوئی سرور `.claude.json` سے غائب ہو جائے۔ بعد کے عام لانچز میں `.claude.json` کی کاپی کو `--mcp-config` پر ترجیح ملتی ہے، اس لیے `ccam mcp add`، `remove` اور `import` وہ کاپیاں اپ ڈیٹ کرتے ہیں جو دوسرے اکاؤنٹس کے پاس ہیں ([شیئر کیے گئے MCP سرورز](/features/mcp))۔

### 10 ms کا بجٹ

لانچ کے راستے کا بجٹ 95ویں پرسنٹائل پر 10 ms ہے؛ عملاً ccam، Claude Code شروع ہونے سے پہلے تقریباً 4 ms شامل کرتا ہے۔ لانچ کا راستہ کبھی cobra کمانڈ ٹری نہیں بناتا، Bubble Tea شروع نہیں کرتا اور نیٹ ورک کنکشن نہیں کھولتا، اور اگر پیکیجز کے `init()` فنکشنز کل ملا کر 5 ms سے زیادہ لیں تو `TestPackageInitIsCheap` ناکام ہوتا ہے۔ `scripts/perf-check.sh`، `hyperfine` سے `ccam --dry-run work`، `ccam work` اور ایک شارٹ کٹ کو ناپتا ہے اور 10 ms سے اوپر ناکام ہوتا ہے، اور `internal/launcher` میں `BenchmarkLaunchPath` مراحل کو لوپ میں چلاتا ہے۔ ہر مرحلے کا وقت دیکھنے کے لیے `CCAM_TRACE=1` مقرر کریں:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## حفاظت اور خاص صورتیں

- نامعلوم نام یا نمبر `no account "<name>"` کے ساتھ رک جاتا ہے اور `ccam list` کی طرف رہنمائی کرتا ہے۔ اگر کوئی اکاؤنٹ ہی نہ ہو تو ccam کہتا ہے کہ `ccam add <name>` سے ایک شامل کریں۔
- اگر کوئی شیئر کیا گیا راستہ عام فائل بن گیا ہو یا غلط جگہ اشارہ کرتا ہو، تو ccam، `stopped because of share link problems:` کے ساتھ رک جاتا ہے، ان کی فہرست دیتا ہے، اور `ccam doctor --fix` چلانے کو کہتا ہے۔ یہ لانچ کے راستے پر ڈیٹا منتقل نہیں کرتا۔
- اگر کوئی دوسرا Claude Code سیشن جو اس اکاؤنٹ کے ساتھ گفتگو کی ہسٹری شیئر کرتا ہو (ایک ہی گروپ) اسی فولڈر میں چل رہا ہو، تو ccam انتباہ دیتا ہے کہ ہسٹری آپس میں مل سکتی ہے اور `Run anyway? [y/N]` پوچھتا ہے۔ فولڈر کا موازنہ symbolic links حل کر کے کیا جاتا ہے۔ ٹرمینل کے بغیر (`/dev/null` سے آنے والا اسٹینڈرڈ اِن پٹ بھی ٹرمینل نہیں مانا جاتا)، یا پس منظر کی کسی job سے، یہ انتباہ دکھا کر جاری رکھتا ہے۔ نیا `-p` رن اپنی الگ گفتگو شروع کرتا ہے اور اسے کوئی انتباہ نہیں ملتا؛ `--version`، `--help` اور وہ سب کمانڈز جو کوئی گفتگو نہیں کھولتیں، جیسے `mcp`، `auth` یا `plugin`، بالکل نہیں جانچی جاتیں۔
- `--dry-run`، MCP sync، نئے سرور کی اطلاع، `mcp list` والا نوٹ، ہسٹری کی جانچ، سیشن کا نشان، عنوان اور آخری استعمال کا وقت چھوڑ دیتا ہے، اور exec نہیں کرتا۔ یہ پھر بھی شیئر لنکس جانچتا ہے، اور عام لانچ کی طرح غائب لنک بنا دیتا ہے۔
- `ccam work rc`، جب `work` پہلے سے چل رہا ہو، `.claude.json` کو نہیں چھیڑتا، یہ بتا دیتا ہے، اور Remote Control کو گروپ کے سرورز کے بغیر شروع کرتا ہے؛ انہیں کاپی کرنے کے لیے سیشن ختم ہونے کے بعد اسے دوبارہ چلائیں۔ `--dry-run` کے ساتھ یہ پیغام پرنٹ کرتا ہے لیکن نہ پوچھتا ہے نہ کاپی کرتا ہے۔
- `mcp_strategy = "sync"` کے ساتھ، ناکام sync، `ccam: MCP sync failed (launching anyway): …` پرنٹ کرتا ہے اور Claude Code ان سرورز کے ساتھ شروع ہوتا ہے جو اس کے پاس پہلے سے تھے۔

## متعلقہ

- [شارٹ کٹ کمانڈز](/features/shortcuts)
- [ccam کے ساتھ Remote Control استعمال کرنا](/guides/remote-control)
- [اکاؤنٹس کا انتظام](/features/accounts)
- [کمانڈز کا حوالہ](/docs/commands)
- [MCP سرورز، سیٹنگز اور ہسٹری شیئر کریں](/docs/sharing)
- [CLAUDE_CONFIG_DIR کی وضاحت](/guides/claude-config-dir)
