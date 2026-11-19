---
title: "অ্যাকাউন্ট পরিচালনা"
description: "ccam add, ccam login, ccam list, ccam rename, ccam remove ও ccam import দিয়ে Claude Code অ্যাকাউন্ট যোগ, লগইন, তালিকা, নাম বদল, অপসারণ ও ইমপোর্ট করুন।"
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

প্রতিটি ccam অ্যাকাউন্ট নিজস্ব লগইনসহ আলাদা একটি Claude Code কনফিগ ফোল্ডার। অ্যাকাউন্ট যোগ করতে, সেগুলোতে সাইন ইন করাতে, তালিকা দেখতে, নাম বদলাতে ও সরাতে এই কমান্ডগুলো ব্যবহার করুন। এগুলো আপনার আগে থেকে থাকা ফোল্ডারও নিবন্ধন করে।

## ব্যবহার

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

যেখানে কোনো কমান্ড `<name|number>` নেয়, সেখানে অ্যাকাউন্টের নাম বা তার স্লট নম্বর (`ccam list`-এর `#` কলাম) দিতে পারেন।

`ccam use` শুধু কোন অ্যাকাউন্টটি ডিফল্ট তা বদলায় (`config.toml`-এ `default`)।

টার্মিনালে কোনো অপশন ছাড়া `ccam add` একটি উইজার্ড খোলে। এটি এই ধাপগুলো পার হয়:

1. অ্যাকাউন্টের নাম, যদি না আপনি আগেই দিয়ে থাকেন;
2. অ্যাকাউন্টটি কীভাবে সাইন ইন করবে: এখনই একটি Claude লগইন, একটি API key বা পরে একটি Claude লগইন;
3. কোন কোন গ্রুপে যোগ দেবে, যেখানে `default` চেক করা থাকে (এখনও কোনো গ্রুপ না থাকলে এটি দেখায় `default` গ্রুপ কী শেয়ার করবে);
4. `Create a shortcut command for work? [y/N]`;
5. নিশ্চিত করার জন্য একটি সারাংশ।

`Esc` এক ধাপ পিছিয়ে যায়। `--group`, `--no-group`, `--shortcut` বা `--no-login`-এর যেকোনোটি দিলে, অথবা টার্মিনাল ছাড়া চালালে (যেমন পাইপ থেকে), উইজার্ড বাদ যায়। তখন নাম দেওয়া বাধ্যতামূলক, এবং `--shortcut` কোনো নাম না দিলে শর্টকাট তৈরি হয় না।

`--group` ছাড়া নতুন অ্যাকাউন্ট `default` গ্রুপে যোগ দেয়। প্রথম অ্যাকাউন্টটি এটি তৈরি করে, যা MCP সার্ভার, সেটিংস, `CLAUDE.md` ও এক্সটেনশন শেয়ার করে; দেখুন [গ্রুপ দিয়ে শেয়ার করা](/features/share)। Windows-এ ফাইল লিংক করা যায় শুধু Developer Mode-এ। এটি ছাড়া প্রথম অ্যাকাউন্টটি শুধু `mcp` ও `extensions` নিয়ে `default` তৈরি করে এবং তা জানায় (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), আর সেটিংস, মেমরি বা ইতিহাস শেয়ার করা কোনো গ্রুপসহ `--group` দিলে ccam তা প্রত্যাখ্যান করে।

`ccam list` স্লট, নাম, শর্টকাট, ইমেইল, গ্রুপ ও সর্বশেষ ব্যবহার দেখায়, ডিফল্ট অ্যাকাউন্টের আগে `*` থাকে। ইমেইলটি হলো ফোল্ডারটি এখন যে ইমেইলে লগইন করা (এর `.claude.json` থেকে); ccam যেটি রেকর্ড করেছিল তা থেকে আলাদা হলে রেকর্ড করাটি পরে `(was …)` হিসেবে আসে:

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

`--json` দিলে প্রতিটি অ্যাকাউন্ট একটি অবজেক্ট, যাতে থাকে `slot`, `shortcut`, `name`, `email` (বর্তমান লগইনের, টেবিলের মতোই), `default`, `dir`, `last_used` (অ্যাকাউন্টটি কখনও চালু না হলে বাদ থাকে), `groups` ও `share`। `share` প্রতিটি আইটেমকে (`mcp`, `settings`, `memory`, `extensions`, `history`) সেই গ্রুপের সঙ্গে মেলায় যেখান থেকে আইটেমটি আসে, অথবা `"local"`-এর সঙ্গে।

## অপশন

| অপশন | কী করে | ডিফল্ট |
| --- | --- | --- |
| `ccam add --no-login` | লগইন না করে অ্যাকাউন্ট তৈরি করে; পরে `ccam login <name>` দিয়ে লগইন করুন | বন্ধ |
| `ccam add --group <group>` | `default`-এর বদলে এই গ্রুপে যোগ দেয়। একাধিকবার দেওয়া যায়; `default` না থাকলে তৈরি হয় | `default` |
| `ccam add --no-group` | কোনো গ্রুপে যোগ দেয় না: অ্যাকাউন্টটি কিছুই শেয়ার করে না | বন্ধ |
| `ccam add --shortcut <command>` | `~/.local/bin`-এ এই শর্টকাট কমান্ডটিও তৈরি করে | নেই |
| `ccam add --api --key-from <source>` | লগইনের বদলে API key ব্যবহার করা একটি অ্যাকাউন্ট; দেখুন [API key অ্যাকাউন্ট](/features/api) | বন্ধ |
| `ccam list --json` | অ্যাকাউন্টগুলো JSON হিসেবে প্রিন্ট করে | বন্ধ |
| `ccam remove -y`, `--yes` | জিজ্ঞেস না করেই সরায় | বন্ধ |
| `ccam logout -y`, `--yes` | জিজ্ঞেস না করেই লগআউট করে | বন্ধ |
| `ccam import --name <name>` | আপনার দেওয়া ফোল্ডারের জন্য অ্যাকাউন্টের নাম; ফোল্ডার ছাড়া প্রত্যাখ্যাত হয় (`--name needs a folder: ccam import <folder> --name ops`) | ফোল্ডারের নাম থেকে নেওয়া |
| `ccam import -y`, `--yes` | পাওয়া প্রতিটি ফোল্ডার জিজ্ঞেস না করেই ইমপোর্ট করে | বন্ধ |
| `--dry-run` | `ccam list` ও `ccam completion` এটি সমর্থন করে, কারণ এগুলো কিছু বদলায় না। `add`, `login`, `logout`, `use`, `rename`, `remove` ও `import` এটি প্রত্যাখ্যান করে, `` --dry-run is not supported by `ccam add`; nothing was changed `` বার্তাসহ (কমান্ডের নাম উল্লেখ করে) | বন্ধ |

`ccam login` (`--expired` ছাড়া), `ccam use` ও `ccam rename`-এর নিজস্ব কোনো অপশন নেই। `ccam completion` নেয় `bash`, `zsh` বা `powershell`; অন্য কিছু দিলে `only bash, zsh and powershell are supported` দিয়ে ব্যর্থ হয়।

অ্যাকাউন্টের নাম ১ থেকে ৩২ অক্ষরের, `a-z`, `0-9`, `-` ও `_` দিয়ে গঠিত, এবং কোনো অক্ষর বা অঙ্ক দিয়ে শুরু হয়। নাম শুধু অঙ্ক দিয়ে হতে পারে না, কারণ `ccam 2` মানে স্লট 2। নামটি কোনো ccam কমান্ডও হতে পারে না (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), তাই `ccam <name>` কখনও দ্ব্যর্থক হয় না। তালিকাটি `internal/config`-এর `config.Reserved`।

## এটি কীভাবে কাজ করে

অ্যাকাউন্ট যোগ করলে এই ক্রমে এর ফোল্ডার তৈরি হয়, এটি গ্রুপে যোগ দেয়, কনফিগ সংরক্ষিত হয়, সাইন-ইন হয় এবং শর্টকাট তৈরি হয়। আপনি `ccam add work --shortcut ccwork` চালালে ccam এই ধাপগুলো নেয় (`internal/cli/profiles.go`-এর `addProfile`):

<!-- motion:accounts-add -->

1. কিছু তৈরি করার আগে ccam নাম (`config.ValidateName`) ও শর্টকাট (`shim.Check`) যাচাই করে।
2. ccam মোড `0700` দিয়ে `~/.local/share/ccam/profiles/work/` তৈরি করে এবং অ্যাকাউন্টটিকে সবচেয়ে ছোট খালি স্লট নম্বর দেয়, তাই সরানো অ্যাকাউন্টের নম্বর আবার ব্যবহার হয় (`profile.Add`, `Config.NextSlot`)। প্রথম অ্যাকাউন্টটি ডিফল্ট হয়।
3. অ্যাকাউন্টটি তার গ্রুপগুলোতে যোগ দেয় (`--group` বা `--no-group` অন্য কিছু না বললে `default`; প্রথম অ্যাকাউন্টটি `default` তৈরি করে), এবং ccam নতুন, খালি ফোল্ডারটিকে সেগুলোর সঙ্গে লিংক করে, যেমন `settings.json`-কে `groups/default/settings.json`-এর সঙ্গে (`applyShare`, `share.Build` ও তারপর `Apply` দিয়ে)। এটি ব্যর্থ হলে নতুন ফোল্ডারটি সরিয়ে দেওয়া হয়, কারণ রয়ে যাওয়া ফোল্ডার প্রতিটি পুনরায় চেষ্টাকে `folder already exists` দিয়ে আটকে দিত।
4. অ্যাকাউন্টটি `fsutil.AtomicWrite`-এর মাধ্যমে `~/.config/ccam/config.toml`-এ সংরক্ষিত হয়: একটি অস্থায়ী ফাইল, `fsync`, তারপর `rename`, তাই ক্র্যাশ হলেও কখনও অর্ধেক লেখা কনফিগ থেকে যায় না।
5. ccam অ্যাকাউন্ট ফোল্ডারে `CLAUDE_CONFIG_DIR` নির্ধারণ করে `claude auth login` চালু করে, অ্যাকাউন্টের শেষ ইমেইলসহ `--email` দিয়ে, যাতে লগইন পেজটি আগে থেকে পূরণ করা থাকে (`runLogin`)। লগইন ব্যর্থ হলে অ্যাকাউন্টটি থেকে যায়, এবং ccam ত্রুটিটি প্রিন্ট করে আপনাকে পরে `ccam login work` চালাতে বলে।
6. ccam `profiles/work/.claude.json` থেকে `oauthAccount` পড়ে (`profile.ReadIdentity`): ইমেইল, অ্যাকাউন্ট ও প্রতিষ্ঠানের ID, প্রতিষ্ঠানের নাম এবং প্ল্যান (`claude_max`, `claude_team` …)। ওই ফাইল কয়েক মেগাবাইট হতে পারে, তাই ccam পুরো ফাইল পার্স করার বদলে শুধু ওই একটি কী খোঁজে (`TopLevelValue`)। ccam এই তথ্যগুলো `config.toml`-এর `[profile.login]`-এ রাখে।
7. ccam `~/.local/bin/ccwork`-কে ccam বাইনারির সঙ্গে লিংক করে, যা চালু হওয়ার সময় `config.toml`-এ `argv[0]` খোঁজে, এবং লিংকটি `~/.local/state/ccam/shortcuts.json`-এ রেকর্ড করে (`shim.Install`)। লিংক তৈরি করা না গেলেও অ্যাকাউন্টটি যোগ হয়, ccam কারণ জানায় (`Could not create shortcuts: …`) এবং সারাংশে দেখায় `shortcut: ccwork (not created; see above)`।

একটি লগইনে কয়েক মিনিট লাগতে পারে, আর এর মধ্যে অন্য কোনো ccam প্রসেস কনফিগ সংরক্ষণ করতে পারে। তাই যে কমান্ডগুলো অপেক্ষা করে, সেগুলো `config.toml` আবার লোড করে, শুধু ওই একটি পরিবর্তন (এখানে ইমেইল) প্রয়োগ করে এবং সংরক্ষণ করে (`env.update`)। একই কারণে `ccam remove` প্রশ্নের পরে আবার লোড করে।

`ccam login` একটি বিদ্যমান অ্যাকাউন্টের জন্য ধাপ ৫ ও ৬ চালায়। Claude Code ইমেইলটিকে শুধু লগইন পেজের জন্য একটি ইঙ্গিত হিসেবে নেয়, তাই আপনি আগের চেয়ে ভিন্ন কোনো Claude অ্যাকাউন্টে সাইন ইন হয়ে যেতে পারেন; ccam নতুন লগইনকে `[profile.login]`-এর সঙ্গে তুলনা করে এবং সেটি জানায় (`ccam: this login is other@company.com, but work was dev@company.com before.`)। এটি নতুন লগইনটি রেখে দেয় এবং কীভাবে তা পূর্বাবস্থায় ফেরাবেন তা জানায়: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

অন্য Claude অ্যাকাউন্টে লগইন করলে অ্যাকাউন্টটির [সীমার রেকর্ড](/features/usage#the-limit-record)-ও বাদ যায়, কারণ ওই সীমাগুলো ছিল অন্য অ্যাকাউন্টের। অন্য কোনো ccam অ্যাকাউন্ট একই Claude অ্যাকাউন্টে লগইন করা থাকলে ccam সতর্ক করে যে দুটি তাদের ব্যবহারের সীমা ভাগ করে নেয়।

`ccam login --expired` একের পর এক এমন প্রতিটি অ্যাকাউন্টে লগইন করায় যা কখনও লগইন করা হয়নি বা যার লগইনের মেয়াদ শেষ: তার রিফ্রেশ টোকেনের মেয়াদ শেষ হয়েছে, অথবা সার্ভার প্রত্যাখ্যান করার পর Claude Code লগইনটি মুছে দিয়েছে। প্রথম ব্রাউজার খোলার আগে এটি যাচাই করে যে `config.toml` সংরক্ষণ করা যায়।

`ccam logout` অ্যাকাউন্টে `claude auth logout` চালায় এবং সংরক্ষিত ইমেইল মুছে দেয়; লগইনের রেকর্ড থেকে যায়, তাই পরের লগইন তার সঙ্গে মিলিয়ে দেখা হয়। Claude Code লগইনটি বাতিল করে এবং অ্যাকাউন্টের জন্য যা রাখে তা মুছে দেয়, MCP সার্ভারের সাইন-ইনসহ। তাই নিচের ক্ষেত্রগুলোতে ccam আগে সতর্ক করে এবং জিজ্ঞেস করে (`--yes` প্রশ্নটি বাদ দেয়):

- লগইনে MCP সাইন-ইন আছে;
- অ্যাকাউন্টটি চলছে (এর সেশনগুলোও সাইন আউট হয়);
- অন্য কোনো অ্যাকাউন্টের ফোল্ডারে একই লগইন আছে (কপি করা ফোল্ডার), যা লগআউট সেটিকেও বাতিল করে।

উদাহরণ:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

অন্য Claude অ্যাকাউন্ট হিসেবে লগইন করতে শুধু `ccam login <name>`-ই যথেষ্ট, এবং এটি MCP সাইন-ইনগুলো রেখে দেয়।

অ্যাকাউন্টগুলো একাধিক প্রতিষ্ঠানের হলে `ccam list` একটি `Organization` কলাম যোগ করে (ব্যক্তিগত প্ল্যান দেখায় `personal`), আর `ccam list --json`-এ প্রতিটি অ্যাকাউন্টের জন্য `plan` ও `org` থাকে।

### দূরের মেশিনে সাইন ইন

SSH-এর মাধ্যমে `claude auth login`-এর ব্রাউজার অংশটি ভেঙে যেতে পারে। Claude Code রিমোট মেশিনের একটি কলব্যাক সার্ভারে (`http://localhost:<port>/callback`) সাইন-ইনের জন্য অপেক্ষা করে, কিন্তু যে ব্রাউজার সাইন-ইন শেষ করে সেটি চলে আপনার কম্পিউটারে, যেখানে ওই ঠিকানা কোথাও পৌঁছায় না এবং ব্রাউজার `ERR_CONNECTION_REFUSED`-এর মতো ত্রুটি দেখায়। ঠিকানাটিতে তবুও এককালীন কোডটি থাকে, তাই ccam সেটি পৌঁছে দিতে পারে:

<!-- motion:accounts-relay -->

ccam আপনার টার্মিনাল Claude Code-কে দিয়ে দেয় না (`runLogin`)। এটি Claude Code-এর স্ট্যান্ডার্ড ইনপুট একটি পাইপের সঙ্গে যুক্ত করে এবং আপনার টাইপ করা প্রতিটি লাইন তাতে কপি করে (`relayLogin`), তাই Claude Code-এর `Paste code here` প্রম্পটে পেস্ট করা কোডও পৌঁছায়। কোনো লাইন কলব্যাক ঠিকানা হলে ccam সেটি গ্রহণ করে শুধু তখনই, যখন হোস্ট `localhost`, `127.0.0.1` বা `[::1]`, একটি পোর্ট দেওয়া আছে, পাথ `/callback`, এবং `code` ও `state` দুটোই আছে; `http://` বাদ দেওয়া যায় (`callbackURL`)। তারপর ccam সার্ভারে ১০ সেকেন্ডের টাইমআউট দিয়ে এবং রিডাইরেক্ট অনুসরণ না করে ওই ঠিকানায় অনুরোধ পাঠায়, ঠিক যেভাবে ব্রাউজার পাঠাত (`deliverCallback`)। অনুরোধটি অপেক্ষারত `claude auth login`-এর কলব্যাক সার্ভারে পৌঁছায় (ডায়াগ্রামের ড্যাশ দেওয়া রেখা)। Claude Code কোডটি পায় এবং লগইন সংরক্ষণ করে, Linux-এ অ্যাকাউন্ট ফোল্ডারের `.credentials.json`-এ। পেস্ট করা কোনো লাইন কখনও ccam-কে অন্য কোনো হোস্টে যোগাযোগ করাতে পারে না।

লগইন শুরুর আগে ccam এটি প্রিন্ট করে:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

সফল হলে এটি প্রিন্ট করে `→ Passed the sign-in to Claude Code on this machine.`

লগইন টার্মিনালে ইনপুটের জন্য অপেক্ষা করার সময় ccam টার্মিনালের echo ও line mode বন্ধ করে দেয় এবং আপনার টাইপ করা প্রতিটি অক্ষরের জন্য একটি `*` দেখায়, যাতে পেস্ট করা কোড বা callback ঠিকানা স্ক্রিনে দেখা না যায়। ccam নিজেই লাইন সম্পাদনা করে (`readLine`): `Enter` লাইনটি পাঠায়, `Backspace` একটি অক্ষর মোছে, এবং খালি লাইনে `Ctrl+D` ইনপুট শেষ করে। লগইন শেষ হলে টার্মিনাল আগের অবস্থায় ফেরে, এবং `Ctrl+C`, `SIGTERM` বা `SIGHUP`-এও ফেরে, যেগুলো তবুও ccam বন্ধ করে। স্ট্যান্ডার্ড ইনপুট ও স্ট্যান্ডার্ড এররে টার্মিনাল না থাকলে লাইনগুলো যেমন আসে তেমনই কপি হয়।

### নাম বদল, সরানো ও ইমপোর্ট

`ccam rename` ও `ccam remove` অ্যাকাউন্টের ফোল্ডারের ওপর এবং অ্যাকাউন্টের নামে ccam অন্য জায়গায় যা রাখে তার ওপর কাজ করে। ccam ওই পাথগুলোর তালিকা `nameState`-এ রাখে:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- ব্যবহারের ক্যাশ `~/.local/state/ccam/usage/<name>.json` ও `<name>.statusline.json`;
- সীমার রেকর্ড `<name>.limit.json`;
- `~/.local/state/ccam/handoff/<name>/`-এর সেশন নোটগুলো।

`ccam remove work` এগুলো এই ক্রমে সামলায়:

<!-- motion:accounts-remove -->

1. `ccam remove` প্রথমে অ্যাকাউন্টটি চলতে থাকলে প্রত্যাখ্যান করে, তারপর জিজ্ঞেস করে, এবং আপনি উত্তর দেওয়ার পর আবার যাচাই করে (এতে প্রশ্নটি খোলা থাকার সময় চালু হওয়া সেশন ধরা পড়ে)।
2. কিছু বদলানোর আগে এটি যাচাই করে যে `config.toml` সংরক্ষণ করা যায়।
3. অ্যাকাউন্টটি এক্সটেনশন শেয়ার করলে ccam গ্রুপের `installed_plugins.json` ও `known_marketplaces.json` আবার লেখে (`RewritePluginPaths`), কারণ Claude Code প্লাগইনের পূর্ণ পাথ রাখে যা এই অ্যাকাউন্টের ফোল্ডারের মধ্য দিয়ে যেতে পারে।
4. এরপর পরিচালিত অ্যাকাউন্টটি (যার ফোল্ডার ccam তৈরি করেছে, `ccam import` দিয়ে নিবন্ধিত অ্যাকাউন্ট নয়) লগআউট হয়: এর ফোল্ডারে `claude auth logout` চলে, যাতে কোনো কার্যকর লগইন ট্র্যাশে (বা macOS-এ ফোল্ডারের নামে রাখা Keychain আইটেমে, যা পরে ওই নামের অ্যাকাউন্ট উত্তরাধিকারে পেত) অপেক্ষা না করে। API key অ্যাকাউন্টের জন্য এবং যখন অন্য কোনো অ্যাকাউন্টের ফোল্ডারে একই লগইন থাকে, যা লগআউট সেটিকেও বাতিল করত, তখন এটি বাদ যায়। লগআউট ব্যর্থ হলে তা জানানো হয়, এবং সরানো চলতে থাকে।
5. ccam ফোল্ডারটি `~/.local/share/ccam/trash/<name>-<time>/`-এ সরায় (`profile.Trash`), অ্যাকাউন্টটি `config.toml` থেকে বাদ যায়, `nameState` পাথগুলো মুছে ফেলা হয় যাতে পরে একই নামের অ্যাকাউন্ট সেগুলো উত্তরাধিকারে না পায়, এবং শর্টকাট লিংক সরানো হয়: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`।

শেষ পর্যন্ত `config.toml` সংরক্ষণ করা না গেলে ফোল্ডারটি ট্র্যাশ থেকে ফিরে আসে। সরানো অ্যাকাউন্ট ফিরিয়ে আনতে ফোল্ডারটি ট্র্যাশ থেকে আবার সরিয়ে আনুন এবং আবার লগইন করুন।

`ccam rename` ফোল্ডারের নাম বদলায় এবং এর `plugins/`-এ থাকা প্লাগইন পাথগুলো পুরোনো ফোল্ডার থেকে নতুনটিতে আবার লেখে। এটি `nameState` পাথগুলো নতুন নামে সরায় এবং অ্যাকাউন্টটি ডিফল্ট থাকলে ডিফল্টই রাখে। এটি সেশন নোটগুলোও হালনাগাদ করে, যাতে `ccam next --continue` তবুও সেগুলো খুঁজে পায়। শর্টকাট কমান্ডের নাম একই থাকে। প্লাগইন আবার লেখা ব্যর্থ হলে, বা `config.toml` সংরক্ষণ করা না গেলে, ফোল্ডারের নাম আগের মতো ফিরিয়ে দেওয়া হয়।

macOS-এ, এবং Windows-এ যখন Claude Code Credential Manager-এ লগইন রাখে, ফোল্ডারটি নিজে তার পুরোনো নাম রাখে, কারণ লগইন ও MCP সাইন-ইনগুলো ফোল্ডারের পাথের নামে রাখা থাকে এবং সরানো ফোল্ডার লগআউট হয়ে যেত। তখন কনফিগ নতুন নামের পাশে ফোল্ডারটিকে `folder` হিসেবে রেকর্ড করে, এবং ccam জানায় `its folder stays …: the macOS Keychain keeps the login under the folder's path`।

`ccam import` একটি ফোল্ডার যেখানে আছে সেখানেই নিবন্ধন করে, `config.toml`-এ `dir` হিসেবে। এটি কখনও ফোল্ডার সরায় বা কপি করে না: কপি করা রিফ্রেশ টোকেন কপিগুলোকে একে অপরকে লগ আউট করিয়ে দিত।

ফোল্ডার আর্গুমেন্ট ছাড়া এটি অনিবন্ধিত `~/.claude` ও `~/.claude-*` ফোল্ডার খোঁজে এবং প্রতিটির বিষয়ে জিজ্ঞেস করে। নাম আসে ফোল্ডার থেকে (`~/.claude` হলো `main`, `~/.claude-work` হলো `work`, `~/.claude-2` হলো `acct2`)। ccam কোনো নাম বের করতে না পারলে, বা নামটি ইতিমধ্যে নেওয়া থাকলে, ফোল্ডারটি চালানোর মতো `--name` কমান্ডসহ তালিকায় দেখায়। পাওয়া কোনো ফোল্ডার ব্যর্থ হলে তা `Skipped <path>: <reason>` হিসেবে জানানো হয় এবং বাকিগুলো ইমপোর্ট হয়।

ccam আগে সব প্রশ্ন করে নেয়, তারপর `config.toml` আবার লোড করে ফোল্ডারগুলো নিবন্ধন করে, তাই এর মধ্যে অন্য ccam কমান্ড যে পরিবর্তন সংরক্ষণ করেছে তা থেকে যায়। `✓` লাইনগুলো সংরক্ষণের পরে আসে।

`CLAUDE_CONFIG_DIR` ছাড়া Claude Code `.claude.json` আপনার হোম ফোল্ডারে `~/.claude`-এর পাশে রাখে, তাই আপনি `~/.claude` ইমপোর্ট করলে এবং তার নিজের কোনো `.claude.json` না থাকলে ccam `~/.claude.json`-এর দিকে `~/.claude/.claude.json` লিংকটি তৈরি করে। তখন সাধারণ `claude` ও `ccam main` একই অ্যাকাউন্ট দেখে।

`~/.claude`-এ `~/.claude.json`-এর পাশাপাশি নিজস্ব একটি `.claude.json` থাকলে ccam দুটিই রেখে দেয় এবং তা জানায়। `ccam main` ভেতরেরটি আর সাধারণ `claude` বাইরেরটি ব্যবহার করত, ফলে তাদের লগইনের তথ্য ও MCP সার্ভার আলাদা থাকত। একটিই রাখতে ভেতরের ফাইলটির জায়গায় একটি লিংক দিন ([doctor](/features/doctor)-ও এটি জানায়)।

ফোল্ডারটিতে নিবন্ধিত কোনো অ্যাকাউন্টের একই লগইন থাকলে (কপি করা ফোল্ডার) ccam তা জানায়: যেটি আগে নবায়ন করবে সেটি অন্যটিকে লগআউট করে দেবে, তাই ব্যবহারের আগে `ccam login <name>` চালান।

ইমপোর্ট করা অ্যাকাউন্ট কোনো গ্রুপে থাকে না যতক্ষণ না আপনি সেটিকে একটিতে যোগ করেন (`ccam group join <group> <name>`), এবং কোনো শর্টকাট তৈরি হয় না।

`ccam completion bash|zsh` একটি কমপ্লিশন স্ক্রিপ্ট প্রিন্ট করে, যা অ্যাকাউন্টের নামও পূরণ করে।

## নিরাপত্তা ও বিশেষ ক্ষেত্র

- অ্যাকাউন্ট চলতে থাকলে `ccam rename` ও `ccam remove` প্রত্যাখ্যান করে (`refuseRunning`), ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)`` বার্তাসহ, কারণ সেশনটি ফোল্ডারে লিখতে থাকে। যাচাইটি Claude Code যেভাবেই চালু হোক তা ধরে: ccam-এর মাধ্যমে, সাধারণ `claude`, কোনো IDE, ব্যাকগ্রাউন্ড সেশন ও তাদের daemon।
- অন্য কোনো ccam প্রসেস সেই লগইন নবায়ন করছে বা বদলাচ্ছে এমন অবস্থায় `ccam rename`, `ccam remove` ও `ccam logout` প্রত্যাখ্যান করে: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``।
- অন্য কোনো ccam প্রসেসে গ্রুপের পরিবর্তন চলার সময় `ccam rename` ও `ccam remove` প্রত্যাখ্যান করে `another ccam process is changing the groups right now; try again when it is done` দিয়ে; কোনোটি মাঝপথে থেমে যাওয়ার পর, ``the last group change stopped halfway; finish it with `ccam doctor --fix` first`` দিয়ে।
- ইমপোর্ট করা অ্যাকাউন্টে `ccam remove` ফোল্ডার ও তার লগইন রেখে দেয়, সাধারণ `claude`-এর জন্য। অ্যাকাউন্টটি গ্রুপে থাকলে এটি আগে `ccam group leave`-এর মতো করে সেগুলো ছেড়ে যায়: প্রতিটি শেয়ার করা আইটেম ফোল্ডারে গ্রুপেরটির একটি কপি হয়ে যায়, আর গ্রুপের MCP সার্ভারগুলো এর `.claude.json`-এ যায়। তারপর নিবন্ধন বাতিল হয়: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`। ccam এর `.claude.json`-এর যে ব্যাকআপগুলো তৈরি করেছিল সেগুলো `<folder>/.ccam-backup/mcp-<time>`-এ সরে যায়; সরানো না গেলে সেগুলো `~/.local/state/ccam/mcp/<name>/backups`-এ থাকে এবং বার্তায় তা বলা হয়। আপনার `~/.claude` ট্র্যাশে সরানো হয় না।
- দুটি অ্যাকাউন্টের একই ফোল্ডার থাকতে পারে না: যে `config.toml` একটি ফোল্ডারের নাম দুবার দেয়, সেটি লোড হয় না।
- `profiles/`-এ যে নামের ফোল্ডার আগে থেকেই আছে, `ccam add` সেই নাম প্রত্যাখ্যান করে এবং `ccam import` দেখিয়ে দেয়। macOS-এ Keychain-এ যদি নতুন ফোল্ডারের পাথের জন্য ccam-এর কোনো পুরোনো সংস্করণের রেখে যাওয়া লগইন বা MCP সাইন-ইন তখনও থাকে, তাহলেও এটি থেমে যায় এবং সেগুলো সরানোর `security delete-generic-password` কমান্ডটির নাম জানায়। `config.toml` সংরক্ষণ করা না গেলে নতুন ফোল্ডারটি আবার সরানো হয়।
- `ccam import` আগে থেকে নিবন্ধিত ফোল্ডার, ফোল্ডার নয় এমন পাথ, এবং নেওয়া হয়ে গেছে এমন নাম প্রত্যাখ্যান করে। ফোল্ডারটি `~/.claude`-এর কপি হলে (একই লগইন) ccam সতর্কও করে, এমনকি `~/.claude` নিজে নিবন্ধিত না থাকলেও।
- অ্যাকাউন্টটি চলতে থাকলে `ccam logout` আগে জিজ্ঞেস করে; `ccam login` যাচাই করে না।
- রিলে Claude Code-এ পৌঁছাতে না পারলে, যেমন লগইন শেষ হয়ে যাওয়ার পরে, এটি কারণসহ `ccam: could not pass the sign-in to Claude Code:` প্রিন্ট করে। HTTP ত্রুটির মানে সাধারণত কোডের মেয়াদ শেষ; আবার লগইন শুরু করুন।

## সম্পর্কিত

- [Claude Code চালু করা](/features/launch)
- [শর্টকাট কমান্ড](/features/shortcuts)
- [MCP সার্ভার, সেটিংস ও ইতিহাস শেয়ার করুন](/docs/sharing)
- [এক কম্পিউটারে একাধিক Claude Code অ্যাকাউন্ট ব্যবহার করুন](/guides/multiple-claude-code-accounts)
