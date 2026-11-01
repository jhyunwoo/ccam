---
title: "Claude Code চালু করা"
description: "ccam <name|number>, ccam run বা ড্যাশবোর্ড দিয়ে যেকোনো অ্যাকাউন্টে Claude Code চালু করুন; আর্গুমেন্ট অপরিবর্তিত claude-এ যায়, --dry-run কমান্ডটি দেখায়।"
updated: 2026-11-01
order: 1
translation: machine
sourceHash: 11df08608a1e
---

`ccam work` `work` অ্যাকাউন্ট হিসেবে Claude Code চালু করে, আর নামের পরের সবকিছু অপরিবর্তিত অবস্থায় `claude`-এ যায়। যেখানে আপনি `claude` টাইপ করতেন সেখানেই এটি ব্যবহার করুন: টার্মিনালে, স্ক্রিপ্টে, বা পাইপ করা ইনপুটসহ।

## ব্যবহার

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

অ্যাকাউন্টের নাম বা `ccam list` থেকে তার নম্বর দিন। নামের পরের আর্গুমেন্টগুলো ccam পড়ে না, তাই `ccam work --help` Claude Code-এর সাহায্য দেখায়। নামের পরে `--` দেওয়া ঐচ্ছিক। স্ট্যান্ডার্ড ইনপুট, এক্সিট কোড ও সিগন্যাল `claude`-এর মতোই আচরণ করে।

প্রথম আর্গুমেন্ট `-` দিয়ে শুরু হলে এবং সেটি ccam-এর নিজস্ব কোনো ফ্ল্যাগ না হলে, ccam আপনার হয়ে অ্যাকাউন্ট বেছে নেয় এবং সব আর্গুমেন্ট Claude Code-এ পাঠায়। এটি [`ccam folder`](/features/folder) দিয়ে বর্তমান ফোল্ডারের জন্য নির্ধারিত অ্যাকাউন্ট ব্যবহার করে, না থাকলে ডিফল্ট অ্যাকাউন্ট: `ccam use` দিয়ে নির্ধারিত অ্যাকাউন্ট, অথবা আপনি কোনোটি নির্ধারণ না করলে সবচেয়ে ছোট নম্বরের অ্যাকাউন্ট।

### ড্যাশবোর্ড

স্ট্যান্ডার্ড ইনপুট ও আউটপুট টার্মিনাল হলে আর্গুমেন্ট ছাড়া `ccam` [ড্যাশবোর্ড](/features/dashboard) খোলে, অন্যথায় `ccam list` প্রিন্ট করে। একটি অ্যাকাউন্টে যান (অ্যারো কি দিয়ে, বা তার নম্বর টাইপ করে) এবং সেটি চালু করতে `Enter` চাপুন; `n` সবচেয়ে বেশি জায়গা থাকা অ্যাকাউন্টটি চালু করে, আর `o` `--continue`, `--resume` বা অন্য আর্গুমেন্ট দেওয়ার সুযোগ দেয়। Claude Code চালু হওয়ার আগে ড্যাশবোর্ড বন্ধ হয়, তাই টার্মিনাল আগে আগের অবস্থায় ফেরে।

**Team ও Enterprise সিট।** স্বয়ংক্রিয় অ্যাকাউন্ট বাছাই (`ccam next`, `--fallback`, `ccam pick`, ড্যাশবোর্ডের `n` এবং আপনি চালু করলে [`auto_switch`](/features/next#auto-switch)) একই Claude Team বা Enterprise প্রতিষ্ঠানের সিটগুলোর মধ্যেও বদল করে, যা Anthropic-এর শর্তাবলি লঙ্ঘন করতে পারে। এর ওপর নির্ভর করার আগে [প্রশ্নোত্তর](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) পড়ুন এবং আপনার প্রতিষ্ঠান ও Anthropic-এর কাছ থেকে জেনে নিন।

শেলে `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` বা `ANTHROPIC_AUTH_TOKEN` সেট থাকলে Claude Code অ্যাকাউন্টের নিজস্ব লগইনের বদলে সেই key বা টোকেন ব্যবহার করতে পারে, তাই চালু করার আগে ccam স্ট্যান্ডার্ড এররে আপনাকে সতর্ক করে (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`)। এটি শুধু ভেরিয়েবলগুলো পড়ে; একটি [API key অ্যাকাউন্ট](/features/api) এর বদলে সেগুলো মুছে দেয়।

ডিফল্টভাবে ccam নিজের জায়গায় Claude Code বসিয়ে দেয়, তাই আর কিছুই চলতে থাকে না। `config.toml`-এর একেবারে ওপরে, প্রথম `[[group]]`, `[[profile]]` বা `[[folder]]`-এর আগে, `relay = true` (বা Linux ও macOS-এ [`auto_switch = true`](/features/next#auto-switch)) থাকলে ccam এর বদলে ইন্টারঅ্যাকটিভ সেশনগুলো চাইল্ড প্রসেস হিসেবে চালু করে এবং তাদের পাশাপাশি চলতে থাকে। এতে `/ccam:switch`, বা ব্যবহারের সীমায় `auto_switch`, নিজে থেকেই কথোপকথনটি অন্য অ্যাকাউন্টে চালিয়ে যেতে পারে, সেশনের অপশন, permission mode, effort level ও এজেন্ট-সহ; `/ccam:switch` কথোপকথনটি অন্য প্রতিষ্ঠানের কোনো অ্যাকাউন্টে সরানোর আগে ccam আপনাকে জিজ্ঞেস করে ([বিস্তারিত](/features/next#relay))।

## অপশন

| অপশন | কী করে | ডিফল্ট |
| --- | --- | --- |
| `--dry-run` | ccam যে এনভায়রনমেন্ট ও কমান্ড চালাত তা প্রিন্ট করে, তারপর বেরিয়ে যায় | বন্ধ |
| `--fallback` | `-p`/`--print` রানের জন্য: অ্যাকাউন্টটি ব্যবহারের সীমায় পৌঁছালে একই কমান্ড সবচেয়ে বেশি জায়গা থাকা লগইন করা অ্যাকাউন্ট হিসেবে আবার চালায় ([বিস্তারিত](/features/next#run-again-at-the-limit)) | বন্ধ |
| `-v`, `--version` | `ccam <version>` প্রিন্ট করে (শুধু প্রথম আর্গুমেন্ট হিসেবে) | |
| `-h`, `--help` | ccam-এর সাহায্য, `ccam --help` বা `ccam run --help` হিসেবে | |

আর্গুমেন্টের নিয়ম:

- ccam নিজের ফ্ল্যাগ শুধু অ্যাকাউন্টের নামের আগে পড়ে। নামের পরে প্রতিটি আর্গুমেন্ট Claude Code-এর, `--dry-run` ও `--help`-সহ।
- `--dry-run` কাজ করে `ccam --dry-run work`, `ccam run --dry-run work` ও `ccam --dry-run` (ফোল্ডারের অ্যাকাউন্ট বা ডিফল্টের জন্য) হিসেবে। `--dry-run` ও `--fallback` যেকোনো ক্রমে আসতে পারে।
- `ccam run <name|number> [claude args]` হলো `ccam <name>`-এর দীর্ঘ রূপ। নাম ছাড়া `ccam run` ফোল্ডারের অ্যাকাউন্ট বা ডিফল্ট অ্যাকাউন্ট ব্যবহার করে।
- অ্যাকাউন্টের নাম `add`, `list`, `run` বা `usage`-এর মতো ccam কমান্ডের নাম হতে পারে না, তাই `ccam <word>` কখনও দ্ব্যর্থক হয় না।

ড্রাই রান এমন একটি লাইন প্রিন্ট করে যা আপনি শেলে, বা Windows-এ PowerShell-এ, পেস্ট করতে পারেন। লাইনটি শুরু হয় চালু করাটি এনভায়রনমেন্টে কী বদলায় তা দিয়ে: সবসময় `CLAUDE_CONFIG_DIR` ও `CCAM_PROFILE`, যে ভেরিয়েবল এটি বাদ দেয় তার জন্য `env -u <name>` (যেমন উত্তরাধিকারে পাওয়া `CCAM_RELAY`), এবং একটি [API key অ্যাকাউন্টের](/features/api) জন্য তার এন্ডপয়েন্ট, তার মডেল ও key কোথা থেকে আসে (`ANTHROPIC_API_KEY="$WORK_KEY"`), key নিজে কখনও নয়:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

PowerShell লাইনটি পরে ভেরিয়েবলগুলো আগের অবস্থায় ফিরিয়ে দেয়, শেলে `VAR=x cmd` যেমন করে, এবং `Start-Process`-এর মাধ্যমে Claude Code-কে তার কমান্ড লাইন দেয়, তাই Windows PowerShell 5.1 ও PowerShell 7 এটি একইভাবে চালায়। চালু করাটি কোনো ভেরিয়েবলকে খালি মানে নির্ধারণ করলে, তার নিচের একটি `#` লাইন জানায় যে 7.5-এর আগের PowerShell এমন ভেরিয়েবল এর বদলে মুছে দেয়। npm দিয়ে ইনস্টল করা `claude.cmd`-এর ক্ষেত্রে লাইনটি সেই `cmd.exe` কমান্ড চালায় যা চালু করাটি ব্যবহার করত (নিচে দেখুন), এবং চালু করাটি সেখানে যে আর্গুমেন্ট প্রত্যাখ্যান করত, তাতে ড্রাই রানও ব্যর্থ হয়।

## এটি কীভাবে কাজ করে

প্রতিটি চালু করা কমান্ড লাইন থেকে `exec` পর্যন্ত একই ধাপগুলো নেয়:

<!-- motion:launch -->

প্রতিটি বাক্স তার ধাপ চলার সময় আলোকিত হয়, আর নিচের লাইনগুলো হলো সেই এনভায়রনমেন্ট ও আর্গুমেন্ট যা দিয়ে Claude Code চালু হয়।

1. ccam শুরু হয় `app.Main`-এ, যা `cmd/ccam/main.go` ডাকে। ccam অন্য কোনো নামে চালু হলে সেটি একটি [শর্টকাট](/features/shortcuts) এবং সরাসরি লঞ্চারে যায়।
2. cobra যুক্ত হওয়ার আগেই ccam `os.Args` দেখে (`app.parse`)। শুধু একা `ccam` (ড্যাশবোর্ড), কোনো ccam কমান্ডের নাম (`config.Reserved`), `-h`/`--help` ও `-v`/`--version` `cli.Execute`-এ যায়। বাকি সবকিছু চালু করা।
3. ccam `~/.config/ccam/config.toml` পড়ে, এটিই একমাত্র ফাইল যা এটি পার্স করে (`launcher.Launch`)।
4. ccam নাম বা নম্বর দিয়ে অ্যাকাউন্ট খুঁজে বের করে, সামনের একটি `--` বাদ দেয়, এবং `claude` এক্সিকিউটেবল খুঁজে বের করে (`launcher.Prepare`, `ResolveClaude`)। নাম না থাকলে এটি ফোল্ডারের অ্যাকাউন্ট নেয়, ওয়ার্কিং ডিরেক্টরিকে `config.toml`-এর `[[folder]]` এন্ট্রিগুলোর সঙ্গে স্ট্রিং হিসেবে তুলনা করে, macOS ও Windows-এ বড় ও ছোট হাতের অক্ষরের পার্থক্য উপেক্ষা করে (`HereProfile`)। PATH খোঁজার ফলাফল ফাইলের পরিবর্তনের সময় ও `PATH`-এর একটি হ্যাশসহ `~/.local/state/ccam/claude-path`-এ ক্যাশ করা থাকে, তাই ccam প্রতিবার চালু করার সময় `PATH` খোঁজে না। ফাইল বদলালে বা না থাকলে, অথবা `PATH` বদলালে, আবার খোঁজা হয়। `config.toml`-এ থাকা `claude_bin` পাথ (`~/` প্রসারিত করা হয়) একটি এক্সিকিউটেবল ফাইল হতে হবে; না হলে কিছু লেখার আগেই চালু করা `claude_bin X is not an executable file. Check claude_bin in config.toml` দিয়ে থেমে যায়। যে `claude` আসলে ccam নিজেই, ccam সেটি প্রত্যাখ্যান করে, inode দিয়ে তুলনা করে, কারণ তাতে লুপ তৈরি হতো।
5. ccam অ্যাকাউন্ট ফোল্ডার `~/.local/share/ccam/profiles/<name>/`-কে তার শেয়ার সেটিংসের সঙ্গে তুলনা করে, শুধু `lstat` ও `readlink` ডেকে, প্রতি আইটেমে একটি বা দুটি সিস্টেম কল (`checkLinks`, `share.Check`)। অনুপস্থিত লিংক সঙ্গে সঙ্গে তৈরি হয়। ডেটা হারাতে পারে এমন যেকোনো কিছু চালু করা থামিয়ে দেয়।
6. ccam শেয়ার করা MCP সার্ভার যোগ করে (`applyMCP`; নিচে দেখুন)।
7. ccam সেশনটি নথিভুক্ত করে (`session.Register`)। এটি তার PID অনুযায়ী নামের একটি মার্কার, যাতে প্রসেস শুরুর সময় থাকে, `~/.local/state/ccam/locks/profile-<name>/<pid>`-এ লেখে, `XDG_RUNTIME_DIR` যা-ই বলুক, তাই প্রতিটি শেল একই মার্কারগুলো দেখে। যে অ্যাকাউন্ট কথোপকথনের ইতিহাস শেয়ার করে, তার জন্য ccam প্রথমে `history-<group>-<hash>`-এ চলমান সেশন খোঁজে, যেখানে হ্যাশটি সিম্বলিক লিংক মিলিয়ে নেওয়া বর্তমান ফোল্ডারের, এবং সেখানেও নিবন্ধন করে (`historyCheck`)।

   `exec` PID একই রাখে, তাই মার্কারটি তখন Claude Code প্রসেসকেই নির্দেশ করে। ccam ধরে রাখা লকের বদলে একটি মার্কার ফাইল ব্যবহার করে, কারণ Claude Code যে প্রসেসগুলো চালু করে সেগুলো লক উত্তরাধিকারে পেত এবং সেশন শেষ হওয়ার পরেও ধরে রাখত। মার্কার লেখা না গেলে চালু করা থামে না: ccam বলে `ccam: could not mark the session as running (…); starting Claude Code anyway`।

   কোনো অ্যাকাউন্ট চলছে কি না, তা নিয়ে ccam-এর অন্য যাচাইগুলো (নাম বদল, সরানো, গ্রুপ পরিবর্তন, `doctor --fix`, MCP সিঙ্ক) অ্যাকাউন্ট ফোল্ডারে Claude Code-এর নিজস্ব রেকর্ডও পড়ে (`sessions/<pid>.json`, `daemon.lock`), তাই এগুলো ccam ছাড়া চালু হওয়া Claude Code-ও দেখে: সাধারণ `claude`, কোনো IDE, ব্যাকগ্রাউন্ড সেশন ও তাদের daemon।
8. ccam টার্মিনালের শিরোনাম `claude · <name>` করে (যখন স্ট্যান্ডার্ড এরর টার্মিনাল), `fsync` ছাড়া চালুর সময় `~/.local/state/ccam/last-used/<name>`-এ লেখে, এবং `syscall.Exec` ডাকে (`Plan.Exec`)।

Claude Code চলে অ্যাকাউন্ট ফোল্ডারে নির্ধারিত `CLAUDE_CONFIG_DIR` ও অ্যাকাউন্টের নামে নির্ধারিত `CCAM_PROFILE` দিয়ে। এই ভেরিয়েবলগুলোর উত্তরাধিকারে পাওয়া যেকোনো মান ccam মুছে দেয়, উত্তরাধিকারে পাওয়া `CCAM_RELAY`-ও। macOS-এ `ccam import` দিয়ে নিবন্ধিত `~/.claude` ফোল্ডারের জন্য এটি `CLAUDE_SECURESTORAGE_CONFIG_DIR`-ও একটি খালি মানে নির্ধারণ করে, যাতে অ্যাকাউন্টটি ও সাধারণ `claude` একই Keychain লগইন ব্যবহার করে; Windows-এ Claude Code লগইন Credential Manager-এ রাখলে এটি একই কাজ করে।

### কেন exec

ccam নিজের প্রসেসকে `claude` দিয়ে প্রতিস্থাপন করে (`syscall.Exec`)। ফলে:

- ccam মেমরিতে থাকে না;
- কীস্ট্রোক, সিগন্যাল ও এক্সিট কোড Claude Code-এর নিজের;
- আর্গুমেন্ট একটি অ্যারে থাকে, কখনও শেল স্ট্রিং নয়।

Windows-এ কোনো `exec` নেই। সেখানে ccam একই এনভায়রনমেন্ট, আর্গুমেন্ট ও কনসোলসহ `claude.exe`-কে একটি চাইল্ড প্রসেস হিসেবে চালু করে, তার জন্য অপেক্ষা করে এবং তার এক্সিট কোড দিয়ে বেরিয়ে যায় (`Plan.Exec`)। npm-এর ইনস্টল করা `claude.cmd` চলে `cmd.exe`-এ, এর পূর্ণ পাথ ও প্রতিটি আর্গুমেন্ট উদ্ধৃতিচিহ্নের মধ্যে রেখে, যাতে `cmd.exe` সেগুলোর কোনোটিকেই নিজের সিনট্যাক্স হিসেবে না পড়ে (`&`, `%PATH%`); লাইন ব্রেকসহ কোনো আর্গুমেন্ট এভাবে পাঠানো যায় না, তাই তা প্রত্যাখ্যান করা হয়, নেটিভ Claude Code ইনস্টল করার পরামর্শসহ।

Ctrl+C কনসোল থেকে সরাসরি Claude Code-এ পৌঁছায়, তাই ccam সেটি উপেক্ষা করে। সেশন মার্কারে ccam-এর PID থাকে, যা ঠিক Claude Code-এর সমান সময় বেঁচে থাকে, এবং Claude Code বেরিয়ে গেলে মার্কারটি সরানো হয়। অপেক্ষার সময় ccam আর কিছুই করে না।

সেশনের পরের কাজ, যেমন `claude mcp add` দিয়ে যোগ করা সার্ভার লক্ষ করা, পরের বার চালু করার সময় হয়: ccam `.claude.json`-এর পরিবর্তনের সময় `~/.local/state/ccam/mcp/<name>/seen`-এর সঙ্গে তুলনা করে এবং ফাইলটি বদলালেই শুধু পার্স করে।

### কমান্ড লাইন কোন পথে যায়

তিনটি কমান্ড লাইন, `internal/app`-এর মধ্য দিয়ে তিনটি পথ:

<!-- motion:launch-route -->

শর্টকাট মেলানো হয় `argv[0]` দিয়ে, `ccam <name>` প্রথম আর্গুমেন্ট দিয়ে, এবং দুটোই শেষ হয় `syscall.Exec`-এ। শুধু `ccam usage`-এর মতো ccam কমান্ড cobra-তে পৌঁছায়।

### শেয়ার করা MCP সার্ভার

ডিফল্ট `flag` কৌশলে ccam গ্রুপের MCP সার্ভারগুলো একটি আর্গুমেন্ট হিসেবে Claude Code-কে দেয়, আপনার আর্গুমেন্টের আগে বসিয়ে:

<!-- motion:launch-mcp -->

ccam শুধু গ্রুপের `mcp.json` থাকলে আর্গুমেন্টটি যোগ করে, এবং সেটি `claude`-এর ঠিক পরে বসায় (`mcp.FlagArgs`, `Plan.InsertArgs`)। আপনার নিজের `--mcp-config` বা `--strict-mcp-config` তখনও এর ওপরে প্রযোজ্য থাকে। এটি একটি আর্গুমেন্ট, `--mcp-config=<path>`, কারণ Claude Code-এর `--mcp-config` একাধিক মান নেয়: দুটি আর্গুমেন্ট হিসেবে লিখলে এটি পরের প্রম্পটটিকেও একটি কনফিগ ফাইল হিসেবে পড়ত।

ফ্ল্যাগটি কোথায় যেতে পারে না, তা ccam ঠিক করে (`claudeVerb`)। এটি আর্গুমেন্টগুলো থেকে Claude Code-এর সাবকমান্ড পড়ে, Claude Code যেভাবে পড়ে সেভাবেই, অপশনগুলোকে তাদের মানসহ বাদ দিয়ে। `remote-control` বা `rc` কোনো ফ্ল্যাগ পায় না (নিচে দেখুন), প্রথমে দেওয়া ব্যাকগ্রাউন্ড সেশনের কোনো সাবকমান্ডও নয় (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`)। `ccam work mcp list`, `get` ও `remove`-ও কোনো ফ্ল্যাগ পায় না, কারণ Claude Code সেখানে এটি প্রত্যাখ্যান করে; বদলে ccam স্ট্যান্ডার্ড এররে জানায় যে গ্রুপের সার্ভারগুলোও অ্যাকাউন্টের সেশনগুলোতে লোড হয়:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

বাকি প্রতিটি সাবকমান্ড ও প্রতিটি সেশন ফ্ল্যাগটি পায়।

`config.toml`-এ `mcp_strategy = "sync"` থাকলে ccam এর বদলে সার্ভারগুলো অ্যাকাউন্টের `.claude.json`-এ লেখে, এবং অ্যাকাউন্টটি আগে থেকেই চললে সেটি বাদ দেয়। শেষ যাচাইয়ের পর থেকে গ্রুপের `mcp.json` বা `.claude.json` কোনোটিই না বদলালেও চালু করার সময় কাজটি বাদ যায়। কোনো গ্রুপ সার্ভারের নামে অ্যাকাউন্টের নিজস্ব সার্ভার থাকলে সিঙ্ক অ্যাকাউন্টের সার্ভারটিই রাখে এবং একবার তা জানায়। সিঙ্ক `.claude.json` থেকে হারিয়ে যাওয়া কপিও আবার বসিয়ে দেয় (`put back group default's MCP servers …`), এবং শুধু নিজের লেখা সেই কপিগুলো বদলায় বা সরায়, যেগুলো এরপর কেউ বদলায়নি।

### Remote Control

`ccam work rc` (বা `ccam work remote-control`) `work` হিসেবে Claude Code-এর Remote Control চালু করে। `claude remote-control` তার সাবকমান্ডের আগে `--mcp-config` প্রত্যাখ্যান করে এবং পরে কোনোটিই নেয় না, এবং এর সেশনগুলো MCP সার্ভার শুধু অ্যাকাউন্টের `.claude.json` থেকে পড়ে। তাই এখানে ফ্ল্যাগ কৌশল ভিন্নভাবে কাজ করে:

<!-- motion:launch-rc -->

1. ccam Claude Code-এর সাবকমান্ড হিসেবে `remote-control` বা `rc` খুঁজে পায়, এমনকি অপশন ও তাদের মানের পরে এলেও (`ccam work --model opus rc`), এবং কোনো ফ্ল্যাগ যোগ করে না (`claudeVerb`, `mcp.FlagArgs`)।
2. ccam গ্রুপের `mcp.json`-এর পরিবর্তনের সময় অ্যাকাউন্টের `sync-stamp`-এর সঙ্গে তুলনা করে, এবং যাচাই করে যে প্রতিটি গ্রুপ সার্ভার এখনও `.claude.json`-এ আছে (`mcp.Synced`)। `.claude.json`-এ যদি আগে থেকেই গ্রুপের বর্তমান সার্ভারগুলো থাকে, বা গ্রুপে কোনো সার্ভার না থাকে, তাহলে Remote Control কোনো প্রশ্ন ছাড়াই চালু হয়।
3. না হলে ccam সেই সার্ভারগুলোর নাম জানায় যেগুলো সেশনগুলো পাবে না, এবং জিজ্ঞেস করে:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. `y` দিলে ccam সেগুলো `sync` কৌশলের মতো করেই কপি করে (`mcp.Sync`): Claude Code-এর লকের অধীনে, ব্যাকআপ নেওয়ার পরে, অ্যাকাউন্টের নিজস্ব সার্ভারগুলো রেখে। `N` দিলে, বা টার্মিনাল ছাড়া, Remote Control সেগুলো ছাড়াই চালু হয়, এবং ccam `mcp_strategy = "sync"` দেখিয়ে দেয়, যা প্রতিবার চালু করার আগে সেগুলো কপি করে।

সার্ভারগুলো কপি না হওয়া পর্যন্ত ccam প্রতিবার Remote Control চালু করার সময় জিজ্ঞেস করে। কপির পরে এটি আবার জিজ্ঞেস করে শুধু তখনই, যখন গ্রুপের `mcp.json` বদলায় বা কপি করা কোনো সার্ভার `.claude.json` থেকে হারিয়ে যায়। পরের স্বাভাবিক চালু করায় `.claude.json`-এর কপি `--mcp-config`-এর ওপর জেতে, তাই `ccam mcp add`, `remove` ও `import` অন্য অ্যাকাউন্টগুলোর কাছে থাকা কপিগুলো আপডেট করে ([শেয়ার করা MCP সার্ভার](/features/mcp))।

### ১০ ms বাজেট

চালু করার পথের বাজেট 95তম পার্সেন্টাইলে ১০ ms; বাস্তবে Claude Code চালু হওয়ার আগে ccam প্রায় ৪ ms যোগ করে। চালু করার পথ কখনও cobra কমান্ড ট্রি তৈরি করে না, Bubble Tea চালু করে না বা নেটওয়ার্ক সংযোগ খোলে না, এবং প্যাকেজের `init()` ফাংশনগুলো মোট ৫ ms-এর বেশি নিলে `TestPackageInitIsCheap` ব্যর্থ হয়। `scripts/perf-check.sh` `hyperfine` দিয়ে `ccam --dry-run work`, `ccam work` ও একটি শর্টকাট মাপে এবং ১০ ms-এর বেশি হলে ব্যর্থ হয়, আর `internal/launcher`-এর `BenchmarkLaunchPath` ধাপগুলো একটি লুপে চালায়। প্রতি ধাপের সময় দেখতে `CCAM_TRACE=1` নির্ধারণ করুন:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## নিরাপত্তা ও বিশেষ ক্ষেত্র

- অজানা নাম বা নম্বর `no account "<name>"` দিয়ে থেমে যায় এবং `ccam list` দেখিয়ে দেয়। কোনো অ্যাকাউন্টই না থাকলে ccam আপনাকে `ccam add <name>` দিয়ে একটি যোগ করতে বলে।
- কোনো শেয়ার করা পাথ সাধারণ ফাইলে পরিণত হলে বা ভুল জায়গা নির্দেশ করলে ccam `stopped because of share link problems:` দিয়ে থেমে যায়, সেগুলোর তালিকা দেয়, এবং আপনাকে `ccam doctor --fix` চালাতে বলে। চালু করার পথে এটি ডেটা সরায় না।
- এই অ্যাকাউন্টের সঙ্গে কথোপকথনের ইতিহাস শেয়ার করে (একই গ্রুপ) এমন অন্য একটি Claude Code সেশন একই ফোল্ডারে চলতে থাকলে ccam সতর্ক করে যে ইতিহাস মিশে যেতে পারে এবং জিজ্ঞেস করে `Run anyway? [y/N]`। ফোল্ডারটি সিম্বলিক লিংক মিলিয়ে তুলনা করা হয়। টার্মিনাল না থাকলে (`/dev/null` থেকে আসা স্ট্যান্ডার্ড ইনপুটও টার্মিনাল না থাকা হিসেবে গণ্য), বা ব্যাকগ্রাউন্ডের কোনো job থেকে চালালে, এটি সতর্কতা দেখিয়ে চালিয়ে যায়। নতুন `-p` রান নিজস্ব একটি কথোপকথন শুরু করে এবং কোনো সতর্কতা পায় না; `--version`, `--help` এবং যে সাবকমান্ডগুলো কোনো কথোপকথন খোলে না, যেমন `mcp`, `auth` বা `plugin`, সেগুলো একেবারেই যাচাই করা হয় না।
- `--dry-run` MCP সিঙ্ক, নতুন সার্ভারের বার্তা, `mcp list`-এর নোট, ইতিহাসের যাচাই, সেশন মার্কার, শিরোনাম ও সর্বশেষ ব্যবহারের সময় বাদ দেয়, এবং exec করে না। এটি তবুও শেয়ার লিংক যাচাই করে, এবং স্বাভাবিক চালু করার মতোই অনুপস্থিত লিংক তৈরি করে।
- `ccam work rc` চালালে, `work` যদি আগে থেকেই চলে, ccam `.claude.json`-এ হাত দেয় না, তা জানায়, এবং গ্রুপের সার্ভার ছাড়াই Remote Control চালু করে; সেগুলো কপি করতে সেশন শেষ হলে আবার চালান। `--dry-run` থাকলে এটি বার্তাটি প্রিন্ট করে, কিন্তু জিজ্ঞেসও করে না, কপিও করে না।
- `mcp_strategy = "sync"` থাকলে ব্যর্থ সিঙ্ক `ccam: MCP sync failed (launching anyway): …` প্রিন্ট করে এবং Claude Code আগে থেকে থাকা সার্ভারগুলো নিয়েই চালু হয়।

## সম্পর্কিত

- [শর্টকাট কমান্ড](/features/shortcuts)
- [ccam দিয়ে Remote Control ব্যবহার](/guides/remote-control)
- [অ্যাকাউন্ট পরিচালনা](/features/accounts)
- [কমান্ড রেফারেন্স](/docs/commands)
- [MCP সার্ভার, সেটিংস ও ইতিহাস শেয়ার করুন](/docs/sharing)
- [CLAUDE_CONFIG_DIR-এর ব্যাখ্যা](/guides/claude-config-dir)
