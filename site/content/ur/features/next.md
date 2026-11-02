---
title: "اگلے اکاؤنٹ پر جانا"
description: "ccam next کسی دوسرے اکاؤنٹ سے Claude Code کھولتا ہے؛ --auto زیادہ گنجائش والا چنتا ہے، --continue گفتگو ساتھ لے جاتا ہے، --fallback -p جابز دوبارہ چلاتا ہے۔"
updated: 2026-11-03
order: 8
translation: machine
sourceHash: 89e8dea9c55d
---

`ccam next`، Claude Code کو کسی دوسرے اکاؤنٹ کے طور پر کھولتا ہے، بغیر اس کے کہ آپ اس کا نام لیں۔ اسے تب استعمال کریں جب کوئی اکاؤنٹ اپنی استعمال کی حد تک پہنچ جائے۔ اکاؤنٹ بدلنے کے چار طریقے ہیں:

- سادہ `ccam next` آخری بار استعمال کیے گئے اکاؤنٹ کے بعد والا اکاؤنٹ لیتا ہے، اور `--auto` وہ اکاؤنٹ لیتا ہے جس میں سب سے زیادہ گنجائش ہو۔ `--continue` اس فولڈر کی آخری گفتگو نئے اکاؤنٹ پر جاری رکھتا ہے، چاہے دونوں ہسٹری شیئر کرتے ہوں یا نہیں۔
- Claude Code کے اندر، [ccam پلگ اِن](#inside-claude-code) آپ کو بتاتا ہے جب کوئی حد سیشن کو روک دے، اور `/ccam:switch` گفتگو آگے سونپ دیتا ہے۔
- `claude -p` چلانے والے اسکرپٹس کے لیے `ccam --fallback` یہ تبدیلی خود کرتا ہے (دیکھیں [حد پر دوبارہ چلانا](#run-again-at-the-limit))۔
- [`auto_switch`](#auto-switch) آن ہو (Linux اور macOS) تو ccam اس انٹرایکٹو سیشن کو بھی منتقل کرتا ہے جسے Claude Code اس کی حد پر روک دے۔

**Team اور Enterprise سیٹیں۔** خودکار اکاؤنٹ انتخاب (`ccam next`، `--fallback`، `ccam pick`، ڈیش بورڈ کی `n` کلید اور، جب آپ اسے آن کریں تو، [`auto_switch`](/features/next#auto-switch)) ایک ہی Claude Team یا Enterprise تنظیم کی سیٹوں کے درمیان بھی اکاؤنٹ بدلتا ہے، جس سے Anthropic کی شرائط کی خلاف ورزی ہو سکتی ہے۔ اس پر انحصار کرنے سے پہلے [عام سوالات](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) پڑھیں اور اپنی تنظیم اور Anthropic سے تصدیق کر لیں۔

## استعمال

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

ccam تبدیلی کو اسٹینڈرڈ ایرر پر `→ <from> → <to>` کی صورت میں پرنٹ کرتا ہے، `--auto` کے انتخاب کی وجہ کے ساتھ، پھر Claude Code کو نئے اکاؤنٹ کے طور پر اسی طرح شروع کرتا ہے جیسے `ccam <name>` کرتا ہے۔ `<from>` وہ اکاؤنٹ ہے جو آپ نے آخری بار استعمال کیا، یا `--continue` کے ساتھ وہ اکاؤنٹ جس کی یہ گفتگو ہے۔ `--` کے بعد والے آرگیومنٹس بغیر تبدیلی کے `claude` کو جاتے ہیں۔

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `--auto` | لاگ اِن شدہ اکاؤنٹ استعمال کی بنیاد پر چنتا ہے (نیچے دیکھیں)، اور 100% والے اکاؤنٹس چھوڑ دیتا ہے | بند |
| `--max-usage <1-100>` | `--auto` کے ساتھ: اس استعمال سے کم والے اکاؤنٹس کو ترجیح دیتا ہے | `auto_max_usage`، ورنہ 100 |
| `--strategy headroom\|expiring\|balance` | `--auto` کے ساتھ: سب سے کم 5 گھنٹے کا استعمال، 5 گھنٹے کا سب سے جلد ری سیٹ یا فی دن سب سے زیادہ بچا ہوا ہفتہ وار کوٹا | `auto_strategy`، ورنہ `headroom` |
| `--cross-org` | `--auto` کے ساتھ: کسی دوسری تنظیم کے اکاؤنٹس بھی چنتا ہے (نیچے دیکھیں) | `cross_org`، ورنہ بند |
| `--wait[=<duration>]` | `--auto` کے ساتھ: جب کسی اکاؤنٹ میں گنجائش نہ ہو تو پہلے اکاؤنٹ کے ری سیٹ ہونے کا انتظار کرتا ہے، زیادہ سے زیادہ اتنی دیر۔ قدر `=` کے بعد آتی ہے: `--wait 2h` اس اشارے کے ساتھ رد کر دیا جاتا ہے `--wait takes its value after '=': --wait=2h` | بند؛ اکیلا دیا جائے تو `5h` |
| `-c`، `--continue` | اس فولڈر کی آخری گفتگو نئے اکاؤنٹ کے طور پر جاری رکھتا ہے (نیچے دیکھیں) | بند |
| `--brief` | `--continue` کے ساتھ: پوری گفتگو دوبارہ شروع کرنے کے بجائے، نئے اکاؤنٹ کے طور پر پچھلی گفتگو کے ccam کے لکھے خلاصے سے نئی گفتگو شروع کرتا ہے ([نیچے دیکھیں](#brief)) | بند |
| `-- <claude args>` | ccam کے اپنے آرگیومنٹس کے بعد `claude` کو دیے جانے والے آرگیومنٹس | کوئی نہیں |
| `--dry-run` | قبول کیا جاتا ہے: Claude Code شروع کیے یا لانچ ریکارڈ کیے بغیر `→ <from> → <to>` اور وہ کمانڈ لائن پرنٹ کرتا ہے جو چلائی جاتی | بند |

`--auto` اور `--continue` ایک ساتھ کام کرتے ہیں۔ `--continue` کے ساتھ، `--resume` والے آرگیومنٹس آپ کے اپنے آرگیومنٹس سے پہلے رکھے جاتے ہیں۔

Claude Code کے اپنے آپشنز، جیسے `--dangerously-skip-permissions` یا `--model`، `--` کے بعد آتے ہیں۔ `ccam <name>` کے برعکس، جو نام کے بعد سب کچھ آگے بھیج دیتا ہے، `ccam next` اپنے آرگیومنٹس خود پڑھتا ہے، اس لیے `--` کے بغیر یہ `ccam: unknown flag: --dangerously-skip-permissions` کے ساتھ رک جاتا ہے۔

## یہ کیسے کام کرتا ہے

ccam کے ذریعے ہر لانچ وقت کو `~/.local/state/ccam/last-used/<name>` میں لکھتا ہے۔ `ccam next` یہ فائلیں پڑھتا ہے، دیکھتا ہے کہ کون سے اکاؤنٹس لاگ اِن ہیں، اور انتخاب کرتا ہے۔

1. ccam، `config.toml` سے ہر اکاؤنٹ کے لیے ایک `next.Candidate` بناتا ہے (`cli.candidates`): اکاؤنٹ، اس کے آخری لانچ کا وقت (`launcher.LastUsed`) اور یہ کہ اس کا کارآمد لاگ اِن ہے یا نہیں (`creds.Has`: `.credentials.json` میں، یا macOS پر Keychain میں، Claude لاگ اِن جس کی پوری میعاد ختم نہ ہوئی ہو اور جسے Claude Code نے مٹایا نہ ہو)۔ جو فولڈر اس کے بجائے API key استعمال کرتا ہے (Console لاگ اِن، یا `ccam add --api` سے شامل کیا گیا اکاؤنٹ) اس کی کوئی پلان حدیں نہیں ہوتیں اور وہ کبھی نہیں چنا جاتا۔
2. `--auto` کے بغیر، ccam اکاؤنٹس کو سلاٹ نمبر کے لحاظ سے ترتیب دیتا ہے اور سب سے حال میں لانچ کیا گیا اکاؤنٹ تلاش کرتا ہے (`next.Rotate`)۔ یہی وہ اکاؤنٹ ہے جس سے آپ ہٹ رہے ہیں۔ اگر ابھی تک کوئی اکاؤنٹ لانچ نہیں ہوا تو یہ سب سے کم سلاٹ والا ہوتا ہے۔
3. وہاں سے یہ سلاٹ کی ترتیب میں آگے بڑھتا ہے، آخری سلاٹ سے واپس پہلے پر گھومتا ہے، اور پہلا لاگ اِن شدہ اکاؤنٹ لوٹاتا ہے۔
4. ccam چنے گئے اکاؤنٹ کو بالکل اسی طرح تیار اور شروع کرتا ہے جیسے `ccam <name>` کرتا (`launcher.LaunchWith`): یہ اکاؤنٹ کے لنکس جانچتا ہے، MCP شیئرنگ لاگو کرتا ہے، سیشن رجسٹر کرتا ہے اور `syscall.Exec` کے ذریعے ccam کی جگہ `claude` چلا دیتا ہے۔

یہ گردش آخری بار استعمال کیے گئے اکاؤنٹ سے شروع ہوتی ہے اور اس کے بعد پہلے لاگ اِن شدہ اکاؤنٹ پر رک جاتی ہے:

<!-- motion:next-rotate -->

`work` سب سے حال میں لانچ ہوا تھا، اس لیے گردش اس کے بعد سے شروع ہوتی ہے۔ `side` کا کوئی لاگ اِن نہیں، اس لیے اسے چھوڑ دیا جاتا ہے، اور `lab` چنا جاتا ہے۔ اگر `lab` بھی لاگ آؤٹ ہوتا تو گردش گھوم کر `personal` پر آ جاتی (ڈایاگرام میں `% len`)۔ گردش کبھی وہ اکاؤنٹ نہیں لوٹاتی جہاں سے آپ نے شروع کیا: اگر کوئی دوسرا لاگ اِن شدہ اکاؤنٹ نہ ہو تو یہ غلطی کے ساتھ رک جاتی ہے۔

### --auto

`--auto` سلاٹ کی ترتیب کے بجائے استعمال کی بنیاد پر فیصلہ کرتا ہے۔ ccam لاگ اِن شدہ اکاؤنٹس رکھتا ہے، ان سب کا استعمال اسی طرح ایک ساتھ پوچھتا ہے جیسے `ccam usage` پوچھتا ہے (`Fetcher.FetchAll`)، اور نیچے دیے گئے اصولوں سے چنتا ہے (`cli.pickAuto`، `usage.Policy.Choose()`):

<!-- motion:next-auto -->

ہر اکاؤنٹ کو اس کی 5 گھنٹے اور ہفتہ وار ونڈوز سے پرکھا جاتا ہے، اور جب Claude Code کے آرگیومنٹس میں `--model` ہو (`--` کے بعد) تو اس ماڈل کی ہفتہ وار ونڈو سے بھی (نام میں `opus`، `sonnet`، `haiku` یا `fable`)۔ ccam اکاؤنٹس کو تین درجوں میں بانٹتا ہے (`Policy.Choose()`):

1. **گنجائش**: ہر ونڈو `auto_max_usage` (`config.toml`، ڈیفالٹ 100) سے کم۔ ان میں سے حکمتِ عملی چنتی ہے: `headroom` (ڈیفالٹ) سب سے کم 5 گھنٹے کا استعمال لیتا ہے؛ `expiring` وہ 5 گھنٹے کی ونڈو لیتا ہے جو سب سے جلد ری سیٹ ہو، تاکہ اس کا بچا ہوا کوٹا ضائع ہونے سے پہلے استعمال ہو جائے؛ `balance` وہ اکاؤنٹ لیتا ہے جس کے ہفتے کے ری سیٹ تک فی دن سب سے زیادہ ہفتہ وار کوٹا بچا ہو، اور ہر چلتے ccam سیشن کو اس کے حساب میں شمار کرتا ہے، تاکہ ہفتہ وار کوٹا بغیر استعمال کے نہ رہ جائے اور متوازی سیشنز بٹ جائیں۔
2. **حد سے اوپر**: 100% سے کم لیکن `auto_max_usage` کے برابر یا اس سے زیادہ۔ صرف تب استعمال ہوتا ہے جب کسی اکاؤنٹ میں گنجائش نہ ہو؛ سب سے کم 5 گھنٹے کے استعمال والا پہلے آتا ہے۔
3. **بھرا ہوا**: کوئی ونڈو 100% پر، یا نامعلوم 5 گھنٹے کا استعمال۔ کبھی نہیں چنا جاتا۔

جب Claude Code نے کسی ونڈو میں کوئی درخواست رد کی ہو تو وہ ونڈو اپنے ری سیٹ تک 100% شمار ہوتی ہے: [ccam پلگ اِن](#inside-claude-code) کا ہُک اسے ریکارڈ کرتا ہے (`~/.local/state/ccam/usage/<name>.limit.json`)، چاہے محفوظ اعداد کچھ بھی کہیں۔ اس لیے جو اکاؤنٹ ابھی اپنی حد کو پہنچا ہو وہ اپنے ری سیٹ سے پہلے دوبارہ نہیں چنا جاتا، چاہے اس کے اعداد پرانے ہوں۔

تمام اکاؤنٹس کی سیٹنگز `config.toml` میں سب سے اوپر، پہلے `[[group]]`، `[[profile]]` یا `[[folder]]` سے اوپر آتی ہیں۔ ان میں سے کسی کے نیچے لکھی گئی سیٹنگ اسی ٹیبل کی ہوتی ہے: ccam ہر لانچ پر اس کے بارے میں خبردار کرتا ہے، اور کنفیگ محفوظ کرنے والی کمانڈز اس کے منتقل ہونے تک انکار کرتی ہیں۔

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**تنظیمیں۔** خودکار انتخاب ان اکاؤنٹس میں رہتے ہیں جو اسی قسم کے ہوں جس قسم کا اکاؤنٹ آپ چھوڑ رہے ہیں: اسی Team یا Enterprise تنظیم کی سیٹیں، یا ذاتی پلانز (Pro، Max)، تاکہ کمپنی کے اکاؤنٹ کا کام کسی نجی اکاؤنٹ پر جاری نہ ہو اور نہ اس کے الٹ۔ ccam ہر اکاؤنٹ کی تنظیم اس کے لاگ اِن سے جانتا ہے (`ccam login` اسے ریکارڈ کرتا ہے؛ پرانے لاگ اِنز کے لیے یہ `.claude.json` سے پڑھی جاتی ہے)؛ جس اکاؤنٹ کی تنظیم ccam نہ بتا سکے اسے کسی بھی اکاؤنٹ کے ساتھ ملایا جا سکتا ہے۔ ایک ہی Claude اکاؤنٹ میں لاگ اِن دو فولڈرز اس کی حدیں شیئر کرتے ہیں، اس لیے ایک کبھی دوسرے کی جگہ نہیں لیتا۔ `--cross-org`، یا `config.toml` میں `cross_org = true`، تنظیم کا اصول ہٹا دیتا ہے، لیکن [`/ccam:switch`](#inside-claude-code) کے بعد کسی دوسری تنظیم میں منتقلی پھر بھی آپ کی ہاں کا انتظار کرتی ہے۔

جب کچھ نہ چنا جا سکے تو ccam کوڈ 75 (`EX_TEMPFAIL`) کے ساتھ باہر نکلتا ہے، تاکہ اسکرپٹ "بعد میں کوشش کریں" کو ناکامی سے الگ پہچان سکے۔ `--wait` کے ساتھ یہ اس کے بجائے تب تک رکتا ہے جب تک پہلے بھرے ہوئے اکاؤنٹ کی روکنے والی ونڈو ری سیٹ نہ ہو جائے (ری سیٹ کے وقت سے 30 سیکنڈ بعد، جب اینڈ پوائنٹ نئی ونڈو دکھاتا ہے)، `ccam: no account has room; waiting until 15:10` پرنٹ کرتا ہے، اور دوبارہ چنتا ہے؛ جب کوئی اکاؤنٹ یہ نہ بتائے کہ وہ کب ری سیٹ ہوگا، تو یہ 10 منٹ بعد دوبارہ دیکھتا ہے۔ اگر انتظار ختم ہونے سے پہلے کچھ خالی نہ ہو تو یہ 75 کے ساتھ باہر نکلتا ہے۔

یہ تلاش `ccam usage` کے اصولوں پر چلتی ہے: 60 سیکنڈ سے کم پرانے اعداد دوبارہ استعمال ہوتے ہیں، ہر اکاؤنٹ کی تلاش 3 سیکنڈ بعد ٹائم آؤٹ ہو جاتی ہے، جو اکاؤنٹ جواب نہ دے، یا جس کے ایکسیس ٹوکن کی میعاد ختم ہو چکی ہو، اس کا فیصلہ اس کے محفوظ اعداد سے ہوتا ہے، اور حد کا ریکارڈ ان اعداد کی جگہ لے لیتا ہے۔ جس محفوظ ونڈو کا ری سیٹ کا وقت گزر چکا ہو وہ 0% شمار ہوتی ہے، اس لیے جس اکاؤنٹ کی ونڈو اس دوران ری سیٹ ہو چکی ہو اسے چھوڑا نہیں جاتا۔ `--auto` کبھی ٹوکنز کی تجدید نہیں کرتا۔ سادہ گردش کے برعکس، `--auto` آخری بار استعمال کیا گیا اکاؤنٹ بھی چن سکتا ہے اگر اس میں اب بھی سب سے زیادہ گنجائش ہو۔

### --continue

1. ccam ہر اکاؤنٹ میں اس فولڈر کی آخری گفتگو تلاش کرتا ہے: `projects/<folder>/` کی `<session>.jsonl` فائلوں میں سے (فولڈر کا نام ورکنگ فولڈر ہوتا ہے، لنکس حل کر کے، جس میں حرف یا ہندسے کے سوا ہر کیریکٹر `-` میں بدل دیا جاتا ہے، جیسے Claude Code اسے نام دیتا ہے) وہ جس میں سب سے نیا پرامپٹ یا جواب ہو۔ جن سیشنز میں ایسا کچھ نہ ہو، جیسے وہ سیشن جس نے صرف `/mcp`، `/usage` یا `/status` چلایا ہو، یا خالی فائل، انہیں چھوڑ دیا جاتا ہے۔ اگر کسی سیشن نے کوئی نوٹ چھوڑا ہو (استعمال کی حد یا `/ccam:switch`، نیچے دیکھیں) اور کسی گفتگو میں اس کے 2 منٹ سے زیادہ بعد کا کوئی پرامپٹ یا جواب نہ ہو، تو اس کے بجائے ٹھیک وہی سیشن استعمال ہوتا ہے؛ وہ قطاریں جو Claude Code `/exit` پر لکھتا ہے شمار نہیں ہوتیں۔
2. جس اکاؤنٹ کی یہ گفتگو ہے وہی وہ اکاؤنٹ ہے جس سے آپ ہٹ رہے ہیں، اور نیا اکاؤنٹ اوپر کی طرح، اسی تنظیم کے اندر، چنا جاتا ہے۔ جب کئی اکاؤنٹس وہ ہسٹری شیئر کرتے ہوں تو یہ اس اکاؤنٹ کی مانی جاتی ہے جس کی `.claude.json` کہتی ہو کہ اس نے یہ گفتگو آخری بار چلائی (`lastSessionId`)، ورنہ اس کی جو آخری بار استعمال ہوا۔
3. اگر نیا اکاؤنٹ وہ ہسٹری شیئر کرتا ہے (ایک ہی گروپ فولڈر)، تو ccam `claude --resume <session id>` شروع کرتا ہے، یعنی ٹھیک وہی گفتگو (`--continue` فولڈر میں جو بھی گفتگو سب سے نئی ہو اسے کھول دیتا)۔
4. ورنہ یہ `claude --resume <the conversation's file> --fork-session` شروع کرتا ہے: Claude Code پوری گفتگو ایک نئی سیشن ID کے تحت نئے اکاؤنٹ کی اپنی ہسٹری میں کاپی کر دیتا ہے، اور اصل جوں کی توں رہتی ہے۔ اس کے بعد نئے اکاؤنٹ کا `--continue` اسے ڈھونڈ لیتا ہے۔
5. جب تک آپ `--` کے بعد `--agent` نہ دیں، ccam وہ ایجنٹ بھی آگے پہنچاتا ہے جس کے طور پر گفتگو چل رہی ہے: فولڈر کی `.claude/settings.local.json` یا `.claude/settings.json` میں، یا اکاؤنٹ کی `settings.json` میں `agent` سیٹنگ، ورنہ وہ ایجنٹ جس کا نام گفتگو کی آخری `agent-setting` قطار میں ہو۔ Claude Code گفتگو کاپی کرتے وقت (مرحلہ 4) ایجنٹ واپس نہیں لاتا، اور اپنے ڈیفالٹ ایجنٹ کے طور پر، ہر ٹول کے ساتھ، چلتا رہتا۔ جب نیا اکاؤنٹ پرانے اکاؤنٹ کے ایجنٹس، پلگ اِنز اور سیٹنگز شیئر نہ کرتا ہو (کسی گروپ کے ذریعے وہی `agents` اور `plugins` فولڈرز اور وہی `settings.json`)، تو ccam منتقلی والی لائن کے بعد ایک لائن شامل کرتا ہے: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`۔

ccam بتاتا ہے کہ گفتگو کتنی بڑی ہے: نیا اکاؤنٹ اپنے پہلے جواب کے ساتھ یہ پوری دوبارہ پڑھتا ہے۔ پرامپٹ کیش کبھی تنظیموں کے درمیان شیئر نہیں ہوتے، اور ہر ذاتی پلان اپنی الگ تنظیم ہے، اس لیے جب دونوں اکاؤنٹس مختلف تنظیموں کے ہوں تو ccam یہ شامل کرتا ہے: `prompt caches are not shared between organizations, so the first reply reads it all again`۔ جب پرانا اکاؤنٹ ایسے MCP سرورز میں سائن اِن ہو جن میں نیا نہیں، تو ccam ان کے نام بتاتا ہے: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`۔ کسی گروپ میں `history` شیئر کرنے سے کاپی کی ضرورت نہیں پڑتی، ایک ہی گفتگو رہتی ہے، اور وہ checkpoints بھی شیئر ہوتے ہیں جن سے `/rewind` کوڈ بحال کرتا ہے (`file-history/`)۔

### --brief

`--continue` پوری گفتگو دوبارہ شروع کرتا ہے، اور نیا اکاؤنٹ اپنے پہلے جواب کے ساتھ یہ پوری دوبارہ پڑھتا ہے۔ `ccam next --continue --brief` اس کے بجائے ایک نئی گفتگو شروع کرتا ہے، ایک ایسے خلاصے سے جو ccam کسی ماڈل سے پوچھے بغیر پرانی گفتگو کا لکھتا ہے:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

خلاصے میں، اسی ترتیب سے اور صرف وہیں جہاں کچھ کہنے کو ہو، یہ ہوتا ہے:

1. گفتگو کا عنوان؛
2. اس کے ٹرانسکرپٹ کا path، تاکہ کسی تفصیل کی ضرورت ہو تو Claude اس میں تلاش کر سکے؛
3. آپ کی پہلی درخواست، جب چار سے زیادہ ہوں؛
4. وہ خلاصہ جو Claude Code نے آخری بار گفتگو compact کرتے وقت لکھا؛
5. آپ کی آخری چار درخواستیں، `!` کمانڈز اور سلیش کمانڈز کے ساتھ ویسے ہی جیسے آپ نے ٹائپ کیں؛
6. Claude کا آخری جواب؛
7. todo فہرست؛
8. وہ 10 فائلیں جو Claude نے آخر میں بدلیں؛
9. فولڈر کا `git status --short --branch`، زیادہ سے زیادہ 800 حروف تک۔

todo فہرست زیادہ سے زیادہ 10 آئٹمز رکھتی ہے: مکمل آئٹمز پہلے نکالے جاتے ہیں، پھر سب سے آگے والے مراحل، تاکہ موجودہ مرحلہ اور اگلے مراحل باقی رہیں۔ لمبے حصے بیچ سے کاٹے جاتے ہیں، اور پورا خلاصہ 16,000 حروف سے کم رہتا ہے۔

ccam اسے `~/.local/state/ccam/brief/<session>.md` کے طور پر محفوظ کرتا ہے، جسے صرف آپ پڑھ سکتے ہیں، اور اگلا `--brief` 7 دن سے پرانے خلاصے ہٹا دیتا ہے۔ Claude Code ایک پرامپٹ کے ساتھ شروع ہوتا ہے، اور اس میں ذکر کی گئی فائل منسلک کر دیتا ہے: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` جو پرامپٹ آپ `--` کے بعد دیں وہ آخری جملے کی جگہ لیتا ہے (`Then: <your prompt>`)، اور `--` کے بعد آپ کے باقی آپشنز اس کے بعد آتے ہیں۔

Claude Code کسی ذکر (mention) کو `"` پر ختم کر دیتا ہے اور `#` پر کاٹ دیتا ہے، اس لیے جب خلاصے کے path میں ان میں سے کوئی ہو تو ccam ورکنگ فولڈر کے لحاظ سے اس کے path کے ذریعے ذکر کرتا ہے؛ جب اس path میں بھی کوئی ہو، تو پرامپٹ فائل کا نام لیتا ہے، اور Claude اسے اپنے Read ٹول سے پڑھتا ہے، جو پہلے پوچھتا ہے کیونکہ فائل ورکنگ فولڈرز سے باہر ہے (Manual موڈ میں `-p` رن کے دوران Claude Code یہ پڑھنا رد کر دیتا ہے اور Claude خلاصے کے بغیر آگے بڑھتا ہے)۔

`ccam --dry-run next --continue --brief` خلاصہ پرنٹ کرتا ہے اور کچھ محفوظ نہیں کرتا۔ پرانی گفتگو جیسی تھی ویسی رہتی ہے، اور نئی گفتگو صرف وہی جانتی ہے جو خلاصے میں ہے، اس لیے `--brief` تب استعمال کریں جب کام اس سے زیادہ اہم ہو کہ آپ وہاں تک کیسے پہنچے اس کی ہر تفصیل۔ `--continue` کے بغیر ccam `--brief works with --continue` کے ساتھ رک جاتا ہے؛ جب یہ گفتگو پڑھ نہ سکے یا اس میں خلاصہ کرنے کو کچھ نہ ملے، تو یہ رک جاتا ہے اور بتاتا ہے کہ `--brief` کے بغیر `ccam next --continue` اسے پوری دوبارہ شروع کرتا ہے۔ [ریلے](#relay) اور [`auto_switch`](#auto-switch) ہمیشہ پوری گفتگو دوبارہ شروع کرتے ہیں۔ Claude Code کا اپنا `--brief` آپشن، اس کے دوسرے آپشنز کی طرح، `--` کے بعد آتا ہے۔

## Claude Code کے اندر

[ccam پلگ اِن](/docs/install#claude-code-plugin) دو hooks شامل کرتا ہے (Claude Code 2.1.78 یا بعد کا):

- جب کوئی باری استعمال کی حد پر رک جائے، تو یہ گفتگو نوٹ کر لیتا ہے اور ان ٹرمینلز میں جو اس کی سہولت دیتے ہیں (iTerm2، WezTerm، Windows Terminal، kitty، Ghostty) ڈیسک ٹاپ اطلاع دکھاتا ہے: ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  یہ ہر وہ پیغام پہچانتا ہے جسے Claude Code استعمال کی حد مانتا ہے: 5 گھنٹے کی، ہفتہ وار، سادہ اور ماڈل کی حدیں (`You've hit your session limit`، `You've hit your limit`، کسی ماڈل کے کوٹے کے لیے `You've reached your …`)، usage credits اور خرچ کی حدیں، اور کسی تنظیم کا اکاؤنٹ کو روک دینا۔

  جب Claude Code نے رد کی گئی درخواست گفتگو میں لکھ دی ہو، تو ہُک بھری ہوئی ونڈو کو اس کے ری سیٹ تک ریکارڈ بھی کرتا ہے ([حد کا ریکارڈ](/features/usage#the-limit-record))، سوائے تنظیم کی طرف سے لگائی گئی روک کے، جسے اس کے ایڈمن کسی بھی وقت ہٹا سکتے ہیں۔ سرور کی مختصر رکاوٹ یا کسی API key کی rate limit استعمال کی حد نہیں ہے اور اسے چھیڑا نہیں جاتا، اور نہ ہی بغیر نگرانی والے سیشنز کو (`-p` رنز اور بیک گراؤنڈ سیشنز)۔ کوئی subagent یا ورک فلو ایجنٹ استعمال کی حد پر رک جائے تو یہ اس کے سیشن کا رکنا شمار ہوتا ہے، چاہے مرکزی گفتگو صرف اس کا انتظار کر رہی ہو: Claude Code وہ رد کی گئی درخواست سیشن کے `subagents` فولڈر میں ایجنٹ کی اپنی فائل میں لکھتا ہے، اور ہُک اسے وہیں سے پڑھتا ہے۔ جب کئی ایجنٹ ایک ساتھ رکیں تو سیشن کو ایک ہی نوٹ اور ایک ہی اطلاع ملتی ہے۔

  Claude Code 2.1.234 اور بعد کے ورژن سیشن میں انتظار بھی کرتے ہیں اور 5 گھنٹے کے ری سیٹ کے بعد خود جاری رکھتے ہیں؛ اطلاع آپ کو اس کے بجائے اکاؤنٹ بدلنے کا اختیار دیتی ہے، اور [`auto_switch`](#auto-switch) آن ہو (Linux اور macOS) تو جب کسی دوسرے اکاؤنٹ میں گنجائش ہو، ccam آپ کے لیے اکاؤنٹ بدل دیتا ہے۔ جب کوئی دوسرا اکاؤنٹ جگہ نہ لے سکے، کیونکہ کوئی دوسرا اکاؤنٹ اسی Team یا Enterprise تنظیم کا نہیں، یا کوئی دوسرا ذاتی پلان نہیں (`cross_org` آن ہو تو کوئی بھی دوسرا Claude اکاؤنٹ لے سکتا ہے)، تو اطلاع اس کے بجائے `You have no other personal plans, so ccam cannot move this conversation.` یا `You have no other accounts in Acme, so ccam cannot move this conversation.` پر ختم ہوتی ہے۔
- `/ccam:switch [account]` گفتگو نوٹ کرتا ہے اور کوئی درخواست بھیجے جانے سے پہلے باری ختم کر دیتا ہے: `/exit` ٹائپ کریں، پھر اسے اس اکاؤنٹ پر جس کا آپ نے نام لیا (یا اگلے اکاؤنٹ پر) جاری رکھنے کے لیے `ccam next --continue` چلائیں۔ جس سیشن میں ابھی کوئی گفتگو نہ ہو، اس میں یہ بتاتا ہے کہ جاری رکھنے کو کچھ نہیں۔ اسے صرف وہی چلاتا ہے جو آپ ٹائپ کریں: پلگ اِن 0.3.2 سے، Claude `/ccam:switch` نہیں چلا سکتا، اور Claude کا شیڈول کیا ہوا کوئی کام (CronCreate، ScheduleWakeup یا `.claude/scheduled_tasks.json`) اسے Claude کو سادہ متن کے طور پر دیتا ہے۔

  جب آپ کا بتایا ہوا اکاؤنٹ گفتگو کی تنظیم سے باہر ہو (کسی دوسری Team یا Enterprise تنظیم کی سیٹ، سیٹ پر چلنے والی گفتگو کے لیے ذاتی پلان یا اس کے الٹ، یا کسی بھی طرف API key والا اکاؤنٹ)، تو `ccam next --continue` پہلے پوچھتا ہے، `cross_org` چاہے کچھ بھی کہے: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`۔ جب آپ کوئی نام نہ دیں تو یہ اسی طرح اس اکاؤنٹ سے پہلے بھی پوچھتا ہے جو یہ خود چنتا ہے، جسے `cross_org = true` کسی دوسری تنظیم سے لینے دیتا ہے: `ccam: the account picked for /ccam:switch is home, a personal plan, …`۔ صرف `y` یا `yes` اسے منتقل کرتا ہے؛ Linux اور macOS پر سوال سے پہلے دبائی گئی کیز پہلے ہٹا دی جاتی ہیں۔ کوئی اور جواب، یا پوچھنے کے لیے ٹرمینل کا نہ ہونا، گفتگو اور نوٹ کو وہیں رہنے دیتا ہے جہاں وہ ہیں: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (اگر اکاؤنٹ ccam نے خود چنا تھا تو `ccam next --continue` دوبارہ چنتا ہے)۔ خودکار انتخاب کی طرح، جس اکاؤنٹ کی تنظیم ccam نہ بتا سکے وہ کسی کے ساتھ بھی میل کھاتا ہے۔

نوٹ `~/.local/state/ccam/handoff/<account>/<session>.json` میں رہتا ہے، ہر سیشن کے لیے ایک فائل، اس لیے ایک اکاؤنٹ کے دو سیشنز کبھی ایک دوسرے کا نوٹ اوور رائٹ نہیں کرتے۔ یہ 12 گھنٹے رکھا جاتا ہے اور ایک بار استعمال ہوتا ہے، جب اگلا سیشن شروع ہونے والا ہو؛ ناکام لانچ اسے اگلی کوشش کے لیے چھوڑ دیتا ہے۔ یہ اسی فولڈر کا ہوتا ہے جس میں سیشن شروع ہوا تھا، چاہے Claude Code کے شیل میں `cd` کر لیا گیا ہو۔ `ccam rename` اس میں نام اپ ڈیٹ کرتا ہے؛ اگر اس میں درج اکاؤنٹ بعد میں ہٹا دیا گیا ہو، تو ccam یہ بتا دیتا ہے اور معمول کے مطابق انتخاب کرتا ہے۔

### ریلے

`config.toml` میں سب سے اوپر `relay = true` کے ساتھ، ccam خود کو Claude Code سے بدلنے کے بجائے انٹرایکٹو سیشنز کے ساتھ ساتھ چلتا رہتا ہے، اور `/ccam:switch` کے بعد اس سیشن کی گفتگو خود جاری رکھتا ہے (کسی دوسرے سیشن کے چھوڑے ہوئے نوٹ پر عمل نہیں کیا جاتا): ہُک Claude Code کو ختم کرتا ہے، اور ccam اس سیشن کے نوٹ کے ساتھ `ccam next --continue` چلاتا ہے۔ اکاؤنٹ کا نام نہ دیا گیا ہو تو یہ سب سے زیادہ گنجائش والا اکاؤنٹ لیتا ہے، گفتگو کے ماڈل کو شمار کرتے ہوئے، اور جب `cross_org = true` اسے کسی دوسری تنظیم کا اکاؤنٹ ہونے دے تو گفتگو وہاں منتقل کرنے سے پہلے یہ اوپر والا سوال پوچھتا ہے۔ Windows پر ہُک Claude Code کو ختم نہیں کر سکتا: `/exit` ٹائپ کریں، اور ccam گفتگو جاری رکھتا ہے۔ گفتگو کی تنظیم سے باہر کے اکاؤنٹ کے لیے بھی ہُک اسے ختم نہیں کرتا: یہ کہتا ہے `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`، اور آپ کے `/exit` ٹائپ کرنے کے بعد ccam اوپر والا سوال پوچھتا ہے۔ اس کی قیمت ہر سیشن کے لیے ایک انتظار کرتا ccam پروسیس ہے، اسی لیے یہ ڈیفالٹ طور پر بند ہے۔

ریلے کیا گیا سیشن وہ چیزیں برقرار رکھتا ہے جو گفتگو دوبارہ شروع کرنے سے بحال نہیں ہوتیں: `--add-dir`، `--mcp-config`، ٹولز کی فہرستیں، `--settings`، `--plugin-dir` اور سسٹم پرامپٹ کے آپشنز جیسے آپشنز، گفتگو کا اجازتوں کا موڈ، سیشن کا effort لیول اور وہ ایجنٹ جس کے طور پر گفتگو چل رہی ہے (دیکھیں [`--continue`](#continue))۔ Claude Code خود گفتگو سے `--model` اور `--name` بحال کرتا ہے۔ باقی سب کچھ (`--worktree`، `--debug`، `--remote-control`، وہ آپشنز جنہیں ccam نہیں جانتا) `ccam: not carried over: …` میں بتایا جاتا ہے۔ اکاؤنٹ بدلنے کا یہ عمل کبھی شیئر شدہ ہسٹری والا سوال نہیں پوچھتا، اور اس کے چلنے کے دوران آپ نے جو ٹائپ کیا وہ چھوڑ دیا جاتا ہے؛ اس کا واحد سوال وہی ہے جو کسی دوسری تنظیم کے اکاؤنٹ سے پہلے پوچھا جاتا ہے۔

ریلے گفتگو صرف تب جاری رکھتا ہے جب Claude Code اسی طرح ختم ہوا ہو جیسے `/ccam:switch` نے کہا تھا: تب نہیں جب آپ نے پہلے `/exit` ٹائپ کیا، ختم ہوتے وقت Ctrl+C دبایا، یا ccam کو hangup یا کوئی اور سگنل ملا، اور تب بھی نہیں جب ٹرمینل ہی نہ رہا ہو۔ تب نوٹ باقی رہتا ہے، اور `ccam next --continue` اسے بعد میں جاری رکھتا ہے۔ ccam کے کسی پرانے ورژن کے ریلے کے تحت شروع ہونے والا سیشن، یا بغیر نگرانی والا سیشن (`-p` رنز اور بیک گراؤنڈ سیشنز)، ہُک ختم نہیں کرتا؛ اسے اس کے بجائے `/exit` والا پیغام ملتا ہے۔ `claude remote-control` ریلے کے بغیر چلتا ہے۔

### `auto_switch`

`config.toml` میں سب سے اوپر `auto_switch = true` کے ساتھ (Linux اور macOS)، جب Claude Code کسی سیشن کو استعمال کی حد پر روک دے تو ccam اسے خود کسی دوسرے اکاؤنٹ پر منتقل کرتا ہے: پلگ اِن کا ہُک گنجائش والا کوئی دوسرا اکاؤنٹ چنتا ہے اور Claude Code کو ختم کر دیتا ہے، اور ccam، جو سیشن کے ساتھ ساتھ چل رہا ہوتا ہے، گفتگو وہاں جاری رکھتا ہے، اسی طرح جیسے [ریلے](#relay) `/ccam:switch` کے بعد کرتا ہے۔ یہ ڈیفالٹ طور پر بند ہے اور تجرباتی ہے۔

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam کوئی سیشن صرف تب منتقل کرتا ہے جب یہ سب شرائط پوری ہوں:

1. ccam نے اسے ٹرمینل میں شروع کیا ہو اور اس کے ساتھ ساتھ چل رہا ہو۔ `auto_switch` [ریلے](#relay) کو آن کر دیتا ہے، اس لیے یہ شرط ان انٹرایکٹو سیشنز کے لیے پوری ہوتی ہے جو آپ اسے سیٹ کرنے کے بعد شروع کریں۔ جو سیشن پہلے، یا کسی پرانے ccam کے تحت شروع ہوا ہو، اسے ایک اطلاع ملتی ہے جو اس پر ختم ہوتی ہے: ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` بغیر نگرانی والے سیشنز (`-p` رنز اور بیک گراؤنڈ سیشنز)، `claude remote-control` اور ccam کے بغیر شروع کیا گیا Claude Code کبھی منتقل نہیں ہوتے۔
2. اکاؤنٹ میں [ccam پلگ اِن](/docs/install#claude-code-plugin) آن ہو اور وہ Claude Code 2.1.214 یا بعد کا ورژن چلاتا ہو، جو ہُک کو اپنی پروسیس ID دیتا ہے (`CLAUDE_PID`)۔
3. Claude Code نے رد کی گئی درخواست، اس کے ری سیٹ کے وقت کے ساتھ، اس اکاؤنٹ کی 5 گھنٹے، ہفتہ وار، Opus یا Sonnet کی حد کے لیے گفتگو میں لکھی ہو۔ خرچ کی حد، usage credits کا ختم ہونا، کسی تنظیم کی طرف سے لگائی گئی روک، Fable کا کوٹہ (`You've reached your Fable limit.`)، سرور کی رکاوٹ اور کسی API key کی rate limit کبھی سیشن منتقل نہیں کرتے۔ کسی subagent یا ورک فلو ایجنٹ کی حد شمار ہوتی ہے، چاہے مرکزی گفتگو صرف اس کا انتظار کر رہی ہو؛ جب وہ Opus یا Sonnet کی ہفتہ وار حد ہو تو سیشن صرف تب منتقل ہوتا ہے جب گفتگو خود اسی ماڈل پر چل رہی ہو۔
4. ccam اکاؤنٹ کا پلان جانتا ہو، اس کی `.claude.json` سے یا اس سے جو `ccam login` نے درج کیا۔ ورنہ اطلاع اس پر ختم ہوتی ہے: `ccam does not know personal's plan, so auto_switch leaves this session.`
5. سیشن Remote Control سے جڑا ہوا نہ ہو۔ کسی دوسرے اکاؤنٹ کے طور پر یہ اس ایپ سے دوبارہ نہیں جڑے گا جس سے آپ اسے چلاتے ہیں، اس لیے یہ وہیں رہتا ہے، اور اطلاع یہ بتاتی ہے۔
6. حد اب سے 15 منٹ سے زیادہ بعد ری سیٹ ہو۔ اس سے کم انتظار کے لیے ccam سیشن Claude Code پر چھوڑ دیتا ہے: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. حد کے بعد آپ نے کوئی پرامپٹ ٹائپ نہ کیا ہو۔
8. کسی دوسرے اکاؤنٹ میں گنجائش ہو۔ ccam ویسے ہی چنتا ہے جیسے [`--auto`](#auto) (`auto_max_usage`، `auto_strategy`، گفتگو کا ماڈل)، ان اکاؤنٹس میں سے جو اسی قسم کے ہوں جس قسم کا اکاؤنٹ آپ چھوڑ رہے ہیں: اسی Team یا Enterprise تنظیم کی کوئی دوسری سیٹ، یا کوئی دوسرا ذاتی پلان (Pro، Max)۔ `cross_org` اور `--cross-org` یہاں لاگو نہیں ہوتے۔ اکاؤنٹ کو:
   - ایسے لاگ اِن کے ساتھ لاگ اِن ہونا چاہیے جس کی تجدید Claude Code کر سکے؛
   - اسی Claude اکاؤنٹ میں لاگ اِن نہیں ہونا چاہیے؛
   - API key استعمال نہیں کرنی چاہیے؛
   - ایسا پلان رکھنا چاہیے جو ccam کو معلوم ہو؛
   - usage credits آن نہیں رکھنے چاہئیں، جہاں تک ccam کی اپنی usage معلومات دکھائیں؛
   - Claude Code میں اس فولڈر پر بھروسہ (trust) کیا ہونا چاہیے؛
   - اگر گفتگو کسی ایجنٹ کے طور پر چل رہی ہو (`agent` سیٹنگ یا گفتگو کسی ایجنٹ کا نام لے)، تو جس اکاؤنٹ کو آپ چھوڑ رہے ہیں اس کے ایجنٹس شیئر کرنے چاہئیں: اس کے `agents` اور `plugins` فولڈرز اور اس کی `settings.json` میں سے ہر ایک بالکل وہی ہونا چاہیے (کسی گروپ کے ذریعے)، یا دونوں میں موجود نہ ہو۔

   جب آخری تین شرائط گنجائش والے ہر اکاؤنٹ کو باہر کر دیں تو اطلاع ان کے نام بتاتی ہے: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`، `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` یا `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` جس ونڈو کا [حد کا ریکارڈ](/features/usage#the-limit-record) ہو وہ بھری ہوئی شمار ہوتی ہے، اور جس اکاؤنٹ میں پچھلے 2 منٹ میں کوئی دوسرا سیشن منتقل ہوا ہو وہ باقیوں کے بعد آتا ہے۔ انتخاب تقریباً 6 سیکنڈ کے اندر ہونا ضروری ہے؛ جب usage معلوم کرنے میں اس سے زیادہ وقت لگے تو سیشن وہیں رہتا ہے۔ `usage_source = "claude"` کے ساتھ `auto_switch` صرف محفوظ اعداد سے فیصلہ کرتا ہے، اور ccam، Claude Code کی usage معلومات یا اسٹیٹس لائن سے یہ نہیں جان سکتا کہ usage credits آن ہیں یا نہیں، اس لیے یہ فرض کرتا ہے کہ وہ بند ہیں۔

پھر ہُک Claude Code کو ختم کر دیتا ہے، اور ccam، جو اس کے ساتھ ساتھ چل رہا ہوتا ہے، ان ٹرمینلز میں جو اس کی سہولت دیتے ہیں ایک ڈیسک ٹاپ اطلاع دکھاتا ہے، جیسے `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` یہ اوپر والے نمونہ آؤٹ پٹ کی پہلی لائن پرنٹ کرتا ہے اور چنے گئے اکاؤنٹ کے طور پر اس سیشن کے لیے `ccam next --continue` چلاتا ہے، جو اپنی باقی دو لائنیں پرنٹ کرتا ہے اور Claude Code کو سیشن کے آپشنز، اجازتوں کے موڈ، effort لیول اور ایجنٹ کے ساتھ شروع کرتا ہے، جیسے [ریلے](#relay) کرتا ہے۔ نئے سیشن کا پہلا پیغام ایک مقررہ متن ہوتا ہے، واحد چیز جو ccam آپ کی طرف سے بھیجتا ہے (`/ccam:switch` کے بعد ہونے والی منتقلی کچھ نہیں بھیجتی):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

کسی subagent کی حد کے بعد یہ والا متن ہوتا ہے، کیونکہ subagents، ورک فلوز اور بیک گراؤنڈ کمانڈز پرانے سیشن کے ساتھ ہی رک گئے تھے:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

اگر وہ اکاؤنٹ بھی اپنی حد تک پہنچ جائے تو وہاں بھی یہی ہوتا ہے۔ حد کے ریکارڈ اس سلسلے کو کسی اکاؤنٹ کے ری سیٹ سے پہلے اس پر واپس جانے سے روکتے ہیں، اس لیے جب ہر اکاؤنٹ بھرا ہو تو سیشن آخری اکاؤنٹ پر رہتا ہے۔

ہُک کو `agent` سیٹنگ سے، یا اس قطار سے جو ایجنٹ کا نام بتاتی ہے اور جسے Claude Code وقتاً فوقتاً اور سیشن ختم ہونے پر گفتگو میں لکھتا ہے، پتا چلتا ہے کہ گفتگو کسی ایجنٹ کے طور پر چل رہی ہے۔ اگر اب تک صرف سیشن کا `--agent` اس کا نام بتاتا ہو، اور چنا گیا اکاؤنٹ ایجنٹس شیئر نہ کرتا ہو، تو `ccam next --continue` کچھ شروع نہیں کرتا اور نوٹ رکھ لیتا ہے، اور ایک دوسری ڈیسک ٹاپ اطلاع یہی بتاتی ہے: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``۔

جب کسی دوسرے اکاؤنٹ میں گنجائش نہ ہو تو کچھ منتقل نہیں ہوتا: Claude Code حد کو معمول کے مطابق سنبھالتا ہے، سیشن میں انتظار کرتا ہے یا اپنے اختیارات پیش کرتا ہے، اور اطلاع بتاتی ہے کہ کون سا اکاؤنٹ سب سے پہلے خالی ہوگا، مثلاً `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` سیٹوں کے لیے یہ `None of your other accounts in Acme has room now` کہتی ہے۔ شروع ہو چکی منتقلی روکنے کے لیے، Claude Code کے بند ہوتے وقت Ctrl+C دبائیں: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

منتقلی کی کچھ لاگت ہے۔ Claude Code دوبارہ شروع ہوتا ہے، جس میں چند سیکنڈ لگتے ہیں، اور Claude کے شروع کیے ہوئے بیک گراؤنڈ شیلز، dev سرورز اور ایجنٹس پرانے سیشن کے ساتھ رک جاتے ہیں۔ نیا اکاؤنٹ اپنے پہلے جواب کے ساتھ پوری گفتگو دوبارہ پڑھتا ہے، جیسے [`--continue`](#continue) کے ساتھ، اور ccam بتاتا ہے کہ یہ کتنی بڑی ہے؛ منتقلی کبھی [`--brief`](#brief) استعمال نہیں کرتی۔

`auto_switch` کبھی Claude Code کے حد بتانے سے پہلے اکاؤنٹ نہیں بدلتا، پس منظر میں کبھی آپ کے استعمال پر نظر نہیں رکھتا، اور کبھی کوئی لاگ اِن کاپی، تجدید یا آگے منتقل نہیں کرتا۔ Windows پر اسے نظر انداز کیا جاتا ہے، اور `ccam doctor` یہ بتاتا ہے؛ وہاں، `relay = true` کے ساتھ، `/ccam:switch` اور `/exit` اب بھی گفتگو کو کسی دوسرے اکاؤنٹ پر منتقل کرتے ہیں۔

**Anthropic کی شرائط۔** Team اور Enterprise سیٹوں پر `auto_switch` آپ کا کام اسی تنظیم کی کسی دوسری سیٹ پر منتقل کرتا ہے۔ Anthropic کا ہیلپ سینٹر کہتا ہے کہ "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (کسی Team پلان کی ہر سیٹ ایک نامزد فرد کو تفویض کی جاتی ہے، اور صرف اسی شخص کو اسے استعمال کرنے کی اجازت ہے) اور "Users can't route their work through another user's account to use that user's remaining usage" (صارفین کسی دوسرے صارف کا باقی ماندہ استعمال کام میں لانے کے لیے اپنا کام اس کے اکاؤنٹ کے ذریعے نہیں چلا سکتے) ([Team پلانز](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)؛ اسی طرح [Enterprise پلانز](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) کے لیے بھی)۔

ذاتی پلانز (Pro، Max) کے لیے Anthropic کی شرائط یہ نہیں بتاتیں کہ آیا آپ کے اپنے اکاؤنٹس کسی حد پر خود بخود ایک دوسرے کی جگہ لے سکتے ہیں۔ [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) سروسز تک "through automated or non-human means, whether through a bot, script, or otherwise" (خودکار یا غیر انسانی ذرائع سے، چاہے بوٹ، اسکرپٹ یا کسی اور طریقے سے) رسائی سے منع کرتی ہیں جب تک Anthropic واضح طور پر اس کی اجازت نہ دے، اور ایک منتقلی اوپر والے مقررہ پیغامات میں سے ایک آپ کی طرف سے بھیجتی ہے۔ [Usage Policy](https://www.anthropic.com/legal/aup)، جو 12 نومبر 2026 سے نافذ ہوگی، ممنوعہ کاموں میں "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (ہماری پروڈکٹس میں قائم صلاحیتوں، پابندیوں یا حفاظتی حدود کو جان بوجھ کر بائی پاس کرنا یا ان میں مداخلت کرنا) شامل کرتی ہے۔

ccam یہ دعویٰ نہیں کرتا کہ `auto_switch` ان شرائط کے مطابق ہے۔ اسے آن کرنا آپ کا فیصلہ ہے، اور خطرہ بھی آپ کا ہے: جب Anthropic سمجھے کہ شرائط توڑی گئی ہیں تو وہ کسی اکاؤنٹ کو، یا کسی Team یا Enterprise تنظیم کی رسائی کو، معطل کر سکتا ہے۔

جب تک `auto_switch` آن ہو، `ccam doctor` ccam پلگ اِن کے بغیر اکاؤنٹ اور 2.1.214 سے پرانے Claude Code ورژنز کے بارے میں انتباہ دیتا ہے، اور سیٹوں کے لیے Anthropic کا سیٹ والا اصول دہراتا ہے ([جانچیں](/features/doctor))۔

## حد پر دوبارہ چلانا

`ccam next` اس انٹرایکٹو سیشن کے لیے ہے جو رک گیا ہو۔ غیر انٹرایکٹو رن، یعنی کسی اسکرپٹ یا شیڈول شدہ جاب میں `claude -p`، آپ کا انتظار نہیں کر سکتا، اس لیے اکاؤنٹ کے نام سے پہلے `--fallback` لگانے سے ccam اس کے لیے اکاؤنٹ بدل دیتا ہے:

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

1. `--fallback` صرف `-p` یا `--print` کے ساتھ کام کرتا ہے؛ ورنہ ccam رک جاتا ہے اور `ccam next --auto --continue` کی طرف اشارہ کرتا ہے۔
2. pipe سے آنے والا اسٹینڈرڈ اِن پٹ ایک بار، 10 MiB تک، پڑھا جاتا ہے اور ہر کوشش کو دیا جاتا ہے۔ Claude Code کی طرح، ccam اس کے شروع ہونے کا زیادہ سے زیادہ 3 سیکنڈ انتظار کرتا ہے (سوائے `--input-format stream-json` کے)؛ اگر pipe کھلا رہے لیکن کچھ نہ بھیجے تو رن کو کوئی اِن پٹ نہیں ملتا اور ccam اسٹینڈرڈ ایرر پر یہ بتا دیتا ہے۔
3. ccam، `exec` کے بجائے `claude` کو چائلڈ پروسیس کے طور پر شروع کرتا ہے، اس کا اسٹینڈرڈ آؤٹ پٹ روک کر رکھتا ہے، اور اسٹینڈرڈ ایرر کو جیسے جیسے آئے آگے دیتا ہے۔
4. جب `claude` 0 کے ساتھ باہر نکلے، یا کسی اور وجہ سے ناکام ہو، تو روکا ہوا آؤٹ پٹ اور ایگزٹ کوڈ بغیر تبدیلی کے آگے دیے جاتے ہیں۔ کسی سگنل سے ختم ہونے والے رن کے ساتھ بھی یہی ہوتا ہے (128 سے اوپر کا ایگزٹ کوڈ): اسے رکنے کو کہا گیا تھا، اور یہ دوبارہ نہیں چلایا جاتا۔
5. ناکامی کو استعمال کی حد مانا جاتا ہے جب Claude Code کا پیغام یہ کہے، یعنی استعمال کی حد کا کوئی بھی پیغام جو [ہُک](#inside-claude-code) پہچانتا ہے، `--output-format json` یا `stream-json` آؤٹ پٹ کے اندر بھی، یا جب اکاؤنٹ کا استعمال اس کی 5 گھنٹے، ہفتہ وار یا `--model` کی ونڈو 100% پر دکھائے (حد کے ریکارڈ سمیت)۔
6. پھر اس کوشش کا آؤٹ پٹ چھوڑ دیا جاتا ہے اور وہی کمانڈ کسی دوسرے لاگ اِن اکاؤنٹ سے چلتی ہے: سب سے کم 5 گھنٹے کے استعمال والا، پھر وہ اکاؤنٹس جن کا استعمال معلوم نہ ہو سکے، نمبر کی ترتیب سے۔ 100% والے اکاؤنٹس چھوڑ دیے جاتے ہیں، API key لاگ اِن والے اکاؤنٹس کبھی نہیں آزمائے جاتے، اور ہر اکاؤنٹ زیادہ سے زیادہ ایک بار آزمایا جاتا ہے۔
7. اگر آرگیومنٹس میں `-c`، `--continue`، `-r` یا `--resume` ہو تو صرف اسی ہسٹری گروپ کے اکاؤنٹس آزمائے جاتے ہیں، تاکہ گفتگو مل سکے۔
8. صرف اسی تنظیم کے اکاؤنٹس آزمائے جاتے ہیں، جیسے `--auto` کے ساتھ (`cross_org = true` یہ پابندی ہٹا دیتا ہے)۔

اگلا اکاؤنٹ انہی اصولوں سے چنا جاتا ہے جن سے `--auto` (`auto_max_usage`، `auto_strategy`، `--model` کی ونڈو)۔ جب کوئی اکاؤنٹ باقی نہ رہے تو ccam آخری کوشش کا آؤٹ پٹ پرنٹ کرتا ہے، بتاتا ہے کہ کسی لاگ اِن اکاؤنٹ میں گنجائش باقی نہیں، اور کوڈ 75 کے ساتھ باہر نکلتا ہے۔ اکاؤنٹ کے نام سے پہلے `--wait` کے ساتھ (`ccam --fallback --wait=2h work -p …`)، یہ اس کے بجائے پہلے اکاؤنٹ کے ری سیٹ ہونے کا انتظار کرتا ہے، پھر ہر اکاؤنٹ دوبارہ آزما سکتا ہے؛ یہ ہمیشہ پہلے کم از کم ایک بار رکتا ہے، تاکہ جو اکاؤنٹ ابھی اپنی حد کو پہنچا ہو وہ فوراً دوبارہ نہ چلے۔ یہ مدت رن کے انتظار کے کل وقت کو محدود کرتی ہے، ہر انتظار کو الگ نہیں۔ اسپیس کے ساتھ لکھا جائے (`--wait 30m`) تو یہ اسی اشارے کے ساتھ رد کر دیا جاتا ہے جو `ccam next` دیتا ہے۔ Claude Code کے آرگیومنٹس سے پہلے شروع میں `--` کسی بھی لانچ کی طرح کام کرتا ہے (`ccam --fallback work -- -p …`)۔ `ccam --dry-run --fallback work -p …` پہلی کوشش کی کمانڈ پرنٹ کرتا ہے۔

انٹرایکٹو سیشنز خود بخود نہیں بدلتے، جب تک آپ [`auto_switch`](#auto-switch) آن نہ کریں۔ `ccam statusline install` کے ساتھ، جب 5 گھنٹے یا ہفتہ وار استعمال 90% تک پہنچے تو اسٹیٹس لائن `→ ccam next --auto --continue` شامل کر دیتی ہے ([استعمال کی حدیں](/features/usage))۔

## حفاظت اور خاص صورتیں

- لاگ اِن کے بغیر اکاؤنٹس کبھی نہیں چنے جاتے۔ اگر کوئی دوسرا لاگ اِن شدہ اکاؤنٹ نہ ہو تو سادہ `ccam next`، `no logged-in account other than <name>` کے ساتھ رک جاتا ہے۔
- `--continue` کے ساتھ اور جب فولڈر میں کوئی گفتگو نہ ہو: `no account has a conversation in <folder> to continue`۔
- جب تنظیم کا اصول کوئی اکاؤنٹ نہ چھوڑے: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``۔
- `--auto` کے ساتھ اور جب کوئی لاگ اِن شدہ اکاؤنٹ نہ ہو: `no logged-in account to pick`۔ جب ہر اکاؤنٹ بھرا ہوا ہو یا اس کا استعمال نامعلوم ہو: `` no account with known usage and headroom (see `ccam usage`) ``۔ دونوں ایگزٹ کوڈ 75 پر ختم ہوتے ہیں۔
- `--max-usage`، `--strategy`، `--wait` اور `--cross-org` صرف `--auto` کے ساتھ کام کرتے ہیں، اور 1-100 سے باہر کی قدر یا کوئی نامعلوم حکمتِ عملی رد کر دی جاتی ہے۔
- فورک گفتگو کاپی کرتا ہے، اکاؤنٹ کی باقی حالت نہیں: اصل اکاؤنٹ کے `file-history/` کے checkpoints ساتھ نہیں آتے جب تک ہسٹری شیئر نہ ہو۔
- `--auto`، `--dry-run` کے ساتھ بھی نیٹ ورک درخواستیں بھیجتا ہے۔ یہ صرف پڑھتا ہے، اور کبھی ٹوکن فائلیں نہیں لکھتا۔
- چنا گیا اکاؤنٹ اسی راستے سے شروع ہوتا ہے جس سے `ccam <name>`، اس لیے وہی جانچ لاگو ہوتی ہے، بشمول اس انتباہ کے جب اسی ہسٹری کو شیئر کرنے والا کوئی دوسرا Claude Code سیشن پہلے سے اس فولڈر میں چل رہا ہو۔
- `last-used` صرف تب لکھا جاتا ہے جب Claude Code واقعی شروع ہو۔ `--dry-run` اسے نہیں چھیڑتا، اس لیے ڈرائی رن گردش کو آگے نہیں بڑھاتا۔
- `--fallback` کوشش ختم ہونے تک اسٹینڈرڈ آؤٹ پٹ روک کر رکھتا ہے، اس لیے `--output-format stream-json` جیسا اسٹریمنگ آؤٹ پٹ آخر میں ایک ساتھ آتا ہے۔ جو کلائنٹ `stream-json` اِن پٹ اسٹریم کرتا ہو اور مزید بھیجنے سے پہلے جوابات کا انتظار کرتا ہو وہ `--fallback` استعمال نہیں کر سکتا۔ جو فائلیں پہلی کوشش نے حد تک پہنچنے سے پہلے بدلیں وہ بدلی ہی رہتی ہیں؛ اگلی کوشش پرامپٹ کو نئے سرے سے شروع کرتی ہے۔

## متعلقہ

- [استعمال کی حدیں](/features/usage)
- [استعمال کی حدیں اور اکاؤنٹ تبدیل کرنا](/docs/usage-limits)
- [MCP سرورز، سیٹنگز اور ہسٹری شیئر کریں](/docs/sharing)
- [جب Claude Code استعمال کی حد تک پہنچ جائے تو کیا کریں](/guides/usage-limit-workflow)
