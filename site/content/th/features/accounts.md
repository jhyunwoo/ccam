---
title: การจัดการบัญชี
description: "เพิ่ม เข้าสู่ระบบ แสดงรายการ เปลี่ยนชื่อ ลบ และนำเข้าบัญชี Claude Code ด้วย ccam add, ccam login, ccam list, ccam rename, ccam remove และ ccam import"
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

บัญชี ccam แต่ละบัญชีคือโฟลเดอร์การตั้งค่า Claude Code แยกกันที่มีการเข้าสู่ระบบของตัวเอง ใช้คำสั่งเหล่านี้เพื่อเพิ่มบัญชี ลงชื่อเข้าใช้ แสดงรายการ เปลี่ยนชื่อ และลบบัญชี คำสั่งเหล่านี้ยังลงทะเบียนโฟลเดอร์ที่คุณมีอยู่แล้วได้ด้วย

## การใช้งาน

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

เมื่อคำสั่งรับ `<name|number>` คุณระบุชื่อบัญชีหรือหมายเลขสล็อตก็ได้ (คอลัมน์ `#` ใน `ccam list`)

`ccam use` เปลี่ยนเพียงบัญชีที่เป็นค่าเริ่มต้นเท่านั้น (`default` ใน `config.toml`)

`ccam add` บนเทอร์มินัลโดยไม่ใส่ตัวเลือกจะเปิดตัวช่วยตั้งค่า ซึ่งมีขั้นตอนดังนี้:

1. ชื่อบัญชี หากคุณยังไม่ได้ระบุ
2. วิธีที่บัญชีลงชื่อเข้าใช้: เข้าสู่ระบบ Claude ตอนนี้ ใช้ API key หรือเข้าสู่ระบบ Claude ภายหลัง
3. กลุ่มที่จะเข้าร่วม โดยเลือก `default` ไว้ (หากยังไม่มีกลุ่มใดเลย จะแสดงรายการที่กลุ่ม `default` จะใช้ร่วมกัน)
4. `Create a shortcut command for work? [y/N]`
5. สรุปให้ยืนยัน

`Esc` ใช้ย้อนกลับไปหนึ่งขั้น การใส่ `--group`, `--no-group`, `--shortcut` หรือ `--no-login` ตัวใดตัวหนึ่ง หรือการรันโดยไม่มีเทอร์มินัล (เช่น จากไปป์) จะข้ามตัวช่วยตั้งค่า ในกรณีนั้นต้องระบุชื่อ และจะไม่มีการสร้างคำสั่งลัดเว้นแต่ `--shortcut` จะระบุชื่อไว้

หากไม่ใส่ `--group` บัญชีใหม่จะเข้าร่วมกลุ่ม `default` บัญชีแรกจะสร้างกลุ่มนี้ขึ้น โดยใช้เซิร์ฟเวอร์ MCP การตั้งค่า `CLAUDE.md` และส่วนขยายร่วมกัน ดู[การใช้ร่วมกันด้วยกลุ่ม](/features/share) บน Windows จะลิงก์ไฟล์ได้เฉพาะใน Developer Mode หากไม่ได้เปิด บัญชีแรกจะสร้าง `default` ที่มีเพียง `mcp` และ `extensions` และแจ้งให้ทราบ (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`) และ ccam จะปฏิเสธ `--group` ที่ระบุกลุ่มซึ่งใช้การตั้งค่า หน่วยความจำ หรือประวัติร่วมกัน

`ccam list` แสดงสล็อต ชื่อ คำสั่งลัด อีเมล กลุ่ม และเวลาที่ใช้ล่าสุด โดยมี `*` อยู่หน้าบัญชีเริ่มต้น อีเมลคืออีเมลที่โฟลเดอร์เข้าสู่ระบบอยู่ในตอนนี้ (จาก `.claude.json` ของโฟลเดอร์) เมื่ออีเมลนั้นต่างจากที่ ccam บันทึกไว้ อีเมลที่บันทึกไว้จะตามมาในรูป `(was …)`:

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

เมื่อใช้ `--json` แต่ละบัญชีจะเป็นออบเจ็กต์ที่มี `slot`, `shortcut`, `name`, `email` (ของการเข้าสู่ระบบปัจจุบัน เช่นเดียวกับในตาราง), `default`, `dir`, `last_used` (ไม่มีหากบัญชีไม่เคยถูกเริ่ม), `groups` และ `share` โดย `share` จะจับคู่แต่ละรายการ (`mcp`, `settings`, `memory`, `extensions`, `history`) กับกลุ่มที่รายการนั้นมาจาก หรือ `"local"`

## ตัวเลือก

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `ccam add --no-login` | สร้างบัญชีโดยไม่เข้าสู่ระบบ แล้วค่อยเข้าสู่ระบบภายหลังด้วย `ccam login <name>` | ปิด |
| `ccam add --group <group>` | เข้าร่วมกลุ่มนี้แทน `default` ใส่ซ้ำได้หลายครั้ง และ `default` จะถูกสร้างขึ้นเมื่อยังไม่มี | `default` |
| `ccam add --no-group` | ไม่เข้าร่วมกลุ่มใด: บัญชีจะไม่ใช้สิ่งใดร่วมกัน | ปิด |
| `ccam add --shortcut <command>` | สร้างคำสั่งลัดนี้ใน `~/.local/bin` ด้วย | ไม่มี |
| `ccam add --api --key-from <source>` | บัญชีที่ใช้ API key แทนการเข้าสู่ระบบ ดู[บัญชี API key](/features/api) | ปิด |
| `ccam list --json` | พิมพ์รายการบัญชีเป็น JSON | ปิด |
| `ccam remove -y`, `--yes` | ลบโดยไม่ถาม | ปิด |
| `ccam logout -y`, `--yes` | ออกจากระบบโดยไม่ถาม | ปิด |
| `ccam import --name <name>` | ชื่อบัญชีสำหรับโฟลเดอร์ที่คุณระบุ ถูกปฏิเสธหากไม่ระบุโฟลเดอร์ (`--name needs a folder: ccam import <folder> --name ops`) | ได้มาจากชื่อโฟลเดอร์ |
| `ccam import -y`, `--yes` | นำเข้าทุกโฟลเดอร์ที่พบโดยไม่ถาม | ปิด |
| `--dry-run` | `ccam list` และ `ccam completion` ซึ่งไม่เปลี่ยนแปลงสิ่งใดจะรองรับตัวเลือกนี้ ส่วน `add`, `login`, `logout`, `use`, `rename`, `remove` และ `import` จะปฏิเสธด้วยข้อความ `` --dry-run is not supported by `ccam add`; nothing was changed `` (โดยระบุชื่อคำสั่งนั้น) | ปิด |

`ccam login` (นอกจาก `--expired`), `ccam use` และ `ccam rename` ไม่มีตัวเลือกของตัวเอง `ccam completion` รับ `bash`, `zsh` หรือ `powershell` ค่าอื่นจะล้มเหลวด้วย `only bash, zsh and powershell are supported`

ชื่อบัญชียาว 1 ถึง 32 อักขระ ประกอบด้วย `a-z`, `0-9`, `-` และ `_` และต้องขึ้นต้นด้วยตัวอักษรหรือตัวเลข ชื่อต้องไม่เป็นตัวเลขล้วน เพราะ `ccam 2` หมายถึงสล็อต 2 และชื่อต้องไม่ซ้ำกับคำสั่งของ ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`) เพื่อให้ `ccam <name>` ไม่กำกวม รายการนี้คือ `config.Reserved` ใน `internal/config`

## วิธีการทำงาน

การเพิ่มบัญชีจะสร้างโฟลเดอร์ เข้าร่วมกลุ่ม บันทึกการตั้งค่า ลงชื่อเข้าใช้ และสร้างคำสั่งลัด ตามลำดับนี้ เมื่อคุณรัน `ccam add work --shortcut ccwork` ccam จะทำตามขั้นตอนต่อไปนี้ (`addProfile` ใน `internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. ccam ตรวจสอบชื่อ (`config.ValidateName`) และคำสั่งลัด (`shim.Check`) ก่อนที่จะสร้างสิ่งใด
2. ccam สร้าง `~/.local/share/ccam/profiles/work/` ด้วยโหมด `0700` และให้หมายเลขสล็อตที่ว่างต่ำสุดแก่บัญชี ดังนั้นหมายเลขของบัญชีที่ลบไปแล้วจะถูกนำกลับมาใช้ (`profile.Add`, `Config.NextSlot`) บัญชีแรกจะเป็นบัญชีเริ่มต้น
3. บัญชีเข้าร่วมกลุ่มของตัวเอง (`default` เว้นแต่ `--group` หรือ `--no-group` จะระบุเป็นอย่างอื่น บัญชีแรกจะสร้าง `default`) และ ccam ลิงก์โฟลเดอร์ใหม่ที่ยังว่างเข้ากับกลุ่มเหล่านั้น เช่น ลิงก์ `settings.json` ไปยัง `groups/default/settings.json` (`applyShare` ด้วย `share.Build` แล้วตามด้วย `Apply`) หากขั้นตอนนี้ล้มเหลว โฟลเดอร์ใหม่จะถูกลบ เพราะโฟลเดอร์ที่ค้างอยู่จะทำให้การลองใหม่ทุกครั้งล้มเหลวด้วย `folder already exists`
4. บัญชีจะถูกบันทึกลงใน `~/.config/ccam/config.toml` ผ่าน `fsutil.AtomicWrite`: เขียนไฟล์ชั่วคราว เรียก `fsync` แล้วจึง `rename` ดังนั้นการแครชจะไม่ทิ้งการตั้งค่าที่เขียนไว้ครึ่ง ๆ กลาง ๆ
5. ccam เริ่ม `claude auth login` โดยตั้ง `CLAUDE_CONFIG_DIR` เป็นโฟลเดอร์ของบัญชี และใส่ `--email` เป็นอีเมลล่าสุดของบัญชี เพื่อให้หน้าเข้าสู่ระบบกรอกไว้ให้ (`runLogin`) หากเข้าสู่ระบบไม่สำเร็จ บัญชีจะยังคงอยู่ และ ccam จะพิมพ์ข้อผิดพลาดพร้อมบอกให้คุณรัน `ccam login work` ภายหลัง
6. ccam อ่าน `oauthAccount` จาก `profiles/work/.claude.json` (`profile.ReadIdentity`) ได้แก่ อีเมล ID ของบัญชีและองค์กร ชื่อองค์กร และแพ็กเกจ (`claude_max`, `claude_team` …) ไฟล์นั้นอาจมีขนาดหลายเมกะไบต์ ccam จึงสแกนหาเพียงคีย์เดียวแทนการพาร์สทั้งไฟล์ (`TopLevelValue`) ccam เก็บรายละเอียดเหล่านี้ไว้ใน `[profile.login]` ใน `config.toml`
7. ccam ลิงก์ `~/.local/bin/ccwork` ไปยังไบนารีของ ccam ซึ่งจะค้นหา `argv[0]` ใน `config.toml` เมื่อเริ่มทำงาน และบันทึกลิงก์ไว้ใน `~/.local/state/ccam/shortcuts.json` (`shim.Install`) หากสร้างลิงก์ไม่ได้ บัญชีจะยังถูกเพิ่ม ccam จะบอกเหตุผล (`Could not create shortcuts: …`) และสรุปผลจะแสดง `shortcut: ccwork (not created; see above)`

การเข้าสู่ระบบอาจใช้เวลาหลายนาที และโปรเซส ccam อื่นอาจบันทึกการตั้งค่าในระหว่างนั้น คำสั่งที่ต้องรอจึงโหลด `config.toml` ใหม่ ใช้การเปลี่ยนแปลงเพียงรายการเดียว (ในที่นี้คืออีเมล) แล้วบันทึก (`env.update`) `ccam remove` ก็โหลดใหม่หลังจากถามด้วยเหตุผลเดียวกัน

`ccam login` ทำขั้นตอนที่ 5 และ 6 สำหรับบัญชีที่มีอยู่แล้ว Claude Code ใช้อีเมลเป็นเพียงคำใบ้สำหรับหน้าเข้าสู่ระบบ คุณจึงอาจลงชื่อเข้าใช้บัญชี Claude อื่นที่ไม่ใช่บัญชีเดิมได้ ccam จะเปรียบเทียบการเข้าสู่ระบบใหม่กับ `[profile.login]` และแจ้งให้ทราบ (`ccam: this login is other@company.com, but work was dev@company.com before.`) ccam จะเก็บการเข้าสู่ระบบใหม่ไว้และบอกวิธีย้อนกลับ: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

การเข้าสู่ระบบไปยังบัญชี Claude อื่นยังทำให้[บันทึกขีดจำกัด](/features/usage#the-limit-record)ของบัญชีถูกทิ้งไปด้วย เพราะขีดจำกัดเหล่านั้นเป็นของอีกบัญชีหนึ่ง เมื่อบัญชี ccam อื่นเข้าสู่ระบบด้วยบัญชี Claude เดียวกัน ccam จะเตือนว่าทั้งสองใช้ขีดจำกัดการใช้งานร่วมกัน

`ccam login --expired` จะเข้าสู่ระบบทีละบัญชีให้กับทุกบัญชีที่ไม่เคยเข้าสู่ระบบ หรือที่การเข้าสู่ระบบหมดอายุแล้ว: refresh token หมดอายุ หรือ Claude Code ล้างการเข้าสู่ระบบหลังจากเซิร์ฟเวอร์ปฏิเสธ ก่อนเปิดเบราว์เซอร์ครั้งแรก คำสั่งนี้จะตรวจสอบว่าบันทึก `config.toml` ได้

`ccam logout` รัน `claude auth logout` ในบัญชีนั้นและล้างอีเมลที่บันทึกไว้ แต่บันทึกการเข้าสู่ระบบยังคงอยู่ การเข้าสู่ระบบครั้งถัดไปจึงถูกตรวจเทียบกับบันทึกนั้น Claude Code จะเพิกถอนการเข้าสู่ระบบและลบสิ่งที่เก็บไว้สำหรับบัญชีนั้น รวมถึงการลงชื่อเข้าใช้เซิร์ฟเวอร์ MCP ด้วย ccam จึงเตือนและถามก่อน (`--yes` ข้ามคำถาม) เมื่อ:

- การเข้าสู่ระบบมีการลงชื่อเข้าใช้ MCP อยู่
- บัญชีกำลังทำงาน (เซสชันของบัญชีจะถูกออกจากระบบด้วย)
- โฟลเดอร์ของบัญชีอื่นถือการเข้าสู่ระบบเดียวกัน (โฟลเดอร์ที่ถูกคัดลอก) ซึ่งการออกจากระบบจะเพิกถอนด้วยเช่นกัน

ตัวอย่างเช่น:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

หากต้องการเข้าสู่ระบบเป็นบัญชี Claude อื่น ใช้เพียง `ccam login <name>` ก็พอ และจะคงการลงชื่อเข้าใช้ MCP ไว้

`ccam list` จะเพิ่มคอลัมน์ `Organization` เมื่อบัญชีอยู่ในองค์กรมากกว่าหนึ่งแห่ง (แพ็กเกจส่วนตัวแสดงเป็น `personal`) และ `ccam list --json` มี `plan` และ `org` สำหรับทุกบัญชี

### การลงชื่อเข้าใช้บนเครื่องระยะไกล

เมื่อใช้งานผ่าน SSH ส่วนที่ทำงานในเบราว์เซอร์ของ `claude auth login` อาจใช้ไม่ได้ Claude Code รอการลงชื่อเข้าใช้บนเซิร์ฟเวอร์ callback บนเครื่องระยะไกล (`http://localhost:<port>/callback`) แต่เบราว์เซอร์ที่ทำการลงชื่อเข้าใช้จนเสร็จทำงานอยู่บนคอมพิวเตอร์ของคุณ ซึ่งที่อยู่นั้นไม่ได้ชี้ไปที่ใด เบราว์เซอร์จึงแสดงข้อผิดพลาด เช่น `ERR_CONNECTION_REFUSED` อย่างไรก็ตาม ที่อยู่นั้นยังมีรหัสแบบใช้ครั้งเดียวอยู่ ccam จึงส่งต่อรหัสให้ได้:

<!-- motion:accounts-relay -->

ccam ไม่ได้ส่งเทอร์มินัลของคุณให้ Claude Code โดยตรง (`runLogin`) แต่เชื่อมอินพุตมาตรฐานของ Claude Code เข้ากับไปป์ และคัดลอกแต่ละบรรทัดที่คุณพิมพ์ลงไป (`relayLogin`) ดังนั้นรหัสที่วางในพรอมต์ `Paste code here` ของ Claude Code ก็ยังส่งถึงตามปกติ เมื่อบรรทัดหนึ่งเป็นที่อยู่ callback ccam จะยอมรับก็ต่อเมื่อโฮสต์เป็น `localhost`, `127.0.0.1` หรือ `[::1]` มีการระบุพอร์ต พาธเป็น `/callback` และมีทั้ง `code` และ `state` โดยจะละ `http://` ไว้ก็ได้ (`callbackURL`) จากนั้น ccam จะส่งคำขอไปยังที่อยู่นั้นบนเซิร์ฟเวอร์ โดยหมดเวลาที่ 10 วินาทีและไม่ตาม redirect เหมือนที่เบราว์เซอร์จะทำ (`deliverCallback`) คำขอจะไปถึงเซิร์ฟเวอร์ callback ของ `claude auth login` ที่กำลังรออยู่ (เส้นประในแผนภาพ) Claude Code ได้รับรหัสและบันทึกการเข้าสู่ระบบ ซึ่งบน Linux จะอยู่ใน `.credentials.json` ในโฟลเดอร์ของบัญชี บรรทัดที่วางไม่มีทางทำให้ ccam ติดต่อโฮสต์อื่นได้

ccam พิมพ์ข้อความนี้ก่อนเริ่มเข้าสู่ระบบ:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

เมื่อสำเร็จจะพิมพ์ `→ Passed the sign-in to Claude Code on this machine.`

ขณะที่การเข้าสู่ระบบรออินพุตบนเทอร์มินัล ccam จะปิด echo และโหมดบรรทัดของเทอร์มินัล และแสดง `*` แทนแต่ละตัวอักษรที่คุณพิมพ์ ดังนั้นรหัสหรือที่อยู่ callback ที่วางจะไม่ปรากฏบนหน้าจอ ccam จัดการแก้ไขบรรทัดเอง (`readLine`): `Enter` ส่งบรรทัด `Backspace` ลบหนึ่งตัวอักษร และ `Ctrl+D` บนบรรทัดว่างจะจบอินพุต เทอร์มินัลจะถูกคืนค่าเมื่อการเข้าสู่ระบบจบ และเมื่อได้รับ `Ctrl+C`, `SIGTERM` หรือ `SIGHUP` ด้วย ซึ่งยังคงจบ ccam ตามเดิม หากอินพุตมาตรฐานและเอาต์พุตข้อผิดพลาดมาตรฐานไม่ได้เชื่อมกับเทอร์มินัล บรรทัดจะถูกคัดลอกต่อไปตามที่เข้ามา

### การเปลี่ยนชื่อ การลบ และการนำเข้า

`ccam rename` และ `ccam remove` ทำงานกับโฟลเดอร์ของบัญชี และกับสิ่งที่ ccam เก็บไว้ภายใต้ชื่อบัญชีในที่อื่น ccam ระบุพาธเหล่านั้นไว้ใน `nameState`:

- `~/.local/state/ccam/last-used/<name>`
- `~/.local/state/ccam/mcp/<name>`
- แคชการใช้งาน `~/.local/state/ccam/usage/<name>.json` และ `<name>.statusline.json`
- บันทึกขีดจำกัด `<name>.limit.json`
- บันทึกของเซสชันใน `~/.local/state/ccam/handoff/<name>/`

`ccam remove work` จัดการตามลำดับนี้:

<!-- motion:accounts-remove -->

1. `ccam remove` จะปฏิเสธก่อนหากบัญชีกำลังทำงาน แล้วจึงถาม และตรวจสอบอีกครั้งหลังคุณตอบ (เพื่อตรวจพบเซสชันที่เริ่มขึ้นระหว่างที่คำถามเปิดอยู่)
2. คำสั่งนี้ตรวจสอบว่าบันทึก `config.toml` ได้ก่อนจะเปลี่ยนแปลงสิ่งใด
3. หากบัญชีใช้ส่วนขยายร่วมกัน ccam จะเขียน `installed_plugins.json` และ `known_marketplaces.json` ของกลุ่มใหม่ (`RewritePluginPaths`) เพราะ Claude Code เก็บพาธปลั๊กอินแบบสัมบูรณ์ซึ่งอาจผ่านโฟลเดอร์ของบัญชีนี้
4. จากนั้นบัญชีที่ ccam จัดการ (บัญชีที่ ccam สร้างโฟลเดอร์ให้ ไม่ใช่บัญชีที่ลงทะเบียนด้วย `ccam import`) จะถูกออกจากระบบ: `claude auth logout` จะรันในโฟลเดอร์ของบัญชี เพื่อไม่ให้มีการเข้าสู่ระบบที่ยังใช้งานได้ค้างอยู่ในถังขยะ (หรือบน macOS ในรายการพวงกุญแจ (Keychain) ที่ตั้งชื่อตามโฟลเดอร์ ซึ่งบัญชีที่สร้างภายหลังด้วยชื่อนั้นจะได้รับสืบทอดไป) ขั้นตอนนี้จะถูกข้ามสำหรับบัญชี API key และเมื่อโฟลเดอร์ของบัญชีอื่นถือการเข้าสู่ระบบเดียวกัน ซึ่งการออกจากระบบจะเพิกถอนไปด้วย การออกจากระบบที่ล้มเหลวจะถูกรายงาน และการลบจะดำเนินต่อไป
5. ccam ย้ายโฟลเดอร์ไปยัง `~/.local/share/ccam/trash/<name>-<time>/` (`profile.Trash`) บัญชีจะถูกนำออกจาก `config.toml` พาธใน `nameState` จะถูกลบเพื่อไม่ให้บัญชีที่สร้างภายหลังด้วยชื่อเดียวกันได้รับสืบทอดไป และลิงก์คำสั่งลัดจะถูกลบ: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`

หากสุดท้ายบันทึก `config.toml` ไม่ได้ โฟลเดอร์จะกลับมาจากถังขยะ หากต้องการกู้คืนบัญชีที่ลบไป ให้ย้ายโฟลเดอร์กลับมาจากถังขยะแล้วเข้าสู่ระบบอีกครั้ง

`ccam rename` เปลี่ยนชื่อโฟลเดอร์ และเขียนพาธปลั๊กอินใน `plugins/` ของบัญชีใหม่จากโฟลเดอร์เดิมเป็นโฟลเดอร์ใหม่ คำสั่งนี้ย้ายพาธใน `nameState` ไปใช้ชื่อใหม่ และคงสถานะบัญชีเริ่มต้นไว้หากบัญชีนั้นเป็นบัญชีเริ่มต้น นอกจากนี้ยังปรับบันทึกของเซสชันให้ `ccam next --continue` ยังหาเจอ คำสั่งลัดยังคงชื่อเดิม หากการเขียนพาธปลั๊กอินใหม่ล้มเหลว หรือบันทึก `config.toml` ไม่ได้ โฟลเดอร์จะถูกเปลี่ยนชื่อกลับ

บน macOS และบน Windows เมื่อ Claude Code เก็บการเข้าสู่ระบบไว้ใน Credential Manager ตัวโฟลเดอร์จะคงชื่อเดิมไว้ เพราะการเข้าสู่ระบบและการลงชื่อเข้าใช้ MCP ถูกเก็บไว้ภายใต้พาธของโฟลเดอร์ และโฟลเดอร์ที่ถูกย้ายจะหลุดจากระบบ การตั้งค่าจึงบันทึกโฟลเดอร์ไว้เป็น `folder` ข้างชื่อใหม่ และ ccam จะแจ้งว่า `its folder stays …: the macOS Keychain keeps the login under the folder's path`

`ccam import` ลงทะเบียนโฟลเดอร์ไว้ในตำแหน่งเดิม โดยเก็บเป็น `dir` ใน `config.toml` ccam ไม่ย้ายหรือคัดลอกโฟลเดอร์เลย เพราะ refresh token ที่ถูกคัดลอกจะทำให้สำเนาแต่ละชุดทำให้กันและกันออกจากระบบ

หากไม่ระบุโฟลเดอร์ ccam จะค้นหาโฟลเดอร์ `~/.claude` และ `~/.claude-*` ที่ยังไม่ได้ลงทะเบียนแล้วถามทีละโฟลเดอร์ ชื่อบัญชีได้มาจากชื่อโฟลเดอร์ (`~/.claude` เป็น `main`, `~/.claude-work` เป็น `work`, `~/.claude-2` เป็น `acct2`) หาก ccam หาชื่อไม่ได้ หรือชื่อถูกใช้ไปแล้ว จะแสดงโฟลเดอร์นั้นพร้อมคำสั่ง `--name` ที่ต้องรัน โฟลเดอร์ที่พบแต่ไม่สำเร็จจะถูกรายงานเป็น `Skipped <path>: <reason>` และโฟลเดอร์อื่นจะถูกนำเข้าตามปกติ

ccam จะถามคำถามทั้งหมดก่อน แล้วจึงโหลด `config.toml` ใหม่และลงทะเบียนโฟลเดอร์ การเปลี่ยนแปลงที่คำสั่ง ccam อื่นบันทึกไว้ระหว่างนั้นจึงยังคงอยู่ บรรทัด `✓` จะแสดงหลังการบันทึก

หากไม่มี `CLAUDE_CONFIG_DIR` Claude Code จะเก็บ `.claude.json` ไว้ในโฟลเดอร์ home ข้าง ๆ `~/.claude` ดังนั้นเมื่อคุณนำเข้า `~/.claude` และโฟลเดอร์นั้นไม่มี `.claude.json` ของตัวเอง ccam จะสร้างลิงก์ `~/.claude/.claude.json` ไปยัง `~/.claude.json` จากนั้น `claude` แบบปกติและ `ccam main` จะเห็นบัญชีเดียวกัน

หาก `~/.claude` มี `.claude.json` ของตัวเองอยู่ข้าง `~/.claude.json` ccam จะปล่อยทั้งสองไฟล์ไว้และแจ้งให้ทราบ `ccam main` จะใช้ไฟล์ด้านใน ส่วน `claude` แบบปกติจะใช้ไฟล์ด้านนอก ทำให้รายละเอียดการเข้าสู่ระบบและเซิร์ฟเวอร์ MCP แยกจากกัน หากต้องการให้เหลือไฟล์เดียว ให้แทนที่ไฟล์ด้านในด้วยลิงก์ ([doctor](/features/doctor) ก็รายงานเรื่องนี้เช่นกัน)

เมื่อโฟลเดอร์ถือการเข้าสู่ระบบเดียวกับบัญชีที่ลงทะเบียนไว้ (โฟลเดอร์ที่ถูกคัดลอก) ccam จะแจ้งให้ทราบ: บัญชีใดต่ออายุก่อนจะทำให้อีกบัญชีหลุดจากระบบ ดังนั้นให้รัน `ccam login <name>` ก่อนใช้งาน

บัญชีที่นำเข้าจะไม่อยู่ในกลุ่มใดจนกว่าคุณจะเพิ่มเข้าไปในกลุ่มหนึ่ง (`ccam group join <group> <name>`) และจะไม่มีการสร้างคำสั่งลัด

`ccam completion bash|zsh` พิมพ์สคริปต์ completion ที่เติมชื่อบัญชีให้ด้วย

## ความปลอดภัยและกรณีพิเศษ

- `ccam rename` และ `ccam remove` จะปฏิเสธขณะที่บัญชีกำลังทำงาน (`refuseRunning`) ด้วยข้อความ ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)`` เพราะเซสชันยังคงเขียนลงในโฟลเดอร์ การตรวจสอบนี้ครอบคลุม Claude Code ไม่ว่าจะเริ่มด้วยวิธีใด: ผ่าน ccam, `claude` แบบปกติ, IDE, เซสชันเบื้องหลังและ daemon ของเซสชันเหล่านั้น
- `ccam rename`, `ccam remove` และ `ccam logout` จะปฏิเสธขณะที่โปรเซส ccam อื่นกำลังต่ออายุหรือเปลี่ยนการเข้าสู่ระบบนั้น: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``
- ขณะที่โปรเซส ccam อื่นกำลังเปลี่ยนกลุ่มอยู่ `ccam rename` และ `ccam remove` จะปฏิเสธด้วย `another ccam process is changing the groups right now; try again when it is done` และหลังจากการเปลี่ยนกลุ่มหยุดกลางคัน จะปฏิเสธด้วย ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``
- `ccam remove` กับบัญชีที่นำเข้ามาจะคงโฟลเดอร์และการเข้าสู่ระบบไว้สำหรับ `claude` แบบปกติ หากบัญชีอยู่ในกลุ่ม จะออกจากกลุ่มก่อนแบบเดียวกับ `ccam group leave`: ทุกรายการที่ใช้ร่วมกันจะกลายเป็นสำเนาของกลุ่มในโฟลเดอร์ และเซิร์ฟเวอร์ MCP ของกลุ่มจะไปอยู่ใน `.claude.json` ของโฟลเดอร์ จากนั้นจึงยกเลิกการลงทะเบียน: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)` ข้อมูลสำรองของ `.claude.json` ที่ ccam สร้างไว้จะย้ายไปที่ `<folder>/.ccam-backup/mcp-<time>` หากย้ายไม่ได้ จะยังอยู่ใน `~/.local/state/ccam/mcp/<name>/backups` และข้อความจะแจ้งไว้ `~/.claude` ของคุณจะไม่ถูกย้ายไปถังขยะ
- สองบัญชีจะใช้โฟลเดอร์เดียวกันไม่ได้: `config.toml` ที่ระบุโฟลเดอร์เดียวซ้ำสองครั้งจะโหลดไม่สำเร็จ
- `ccam add` จะปฏิเสธชื่อที่มีโฟลเดอร์อยู่แล้วใน `profiles/` และแนะนำให้ใช้ `ccam import` บน macOS คำสั่งนี้ยังหยุดเมื่อพวงกุญแจ (Keychain) ยังเก็บการเข้าสู่ระบบหรือการลงชื่อเข้าใช้ MCP สำหรับพาธของโฟลเดอร์ใหม่ไว้ ซึ่ง ccam เวอร์ชันเก่าทิ้งไว้ และจะระบุคำสั่ง `security delete-generic-password` ที่ใช้ลบรายการเหล่านั้น หากบันทึก `config.toml` ไม่ได้ โฟลเดอร์ใหม่จะถูกลบออกอีกครั้ง
- `ccam import` จะปฏิเสธโฟลเดอร์ที่ลงทะเบียนไว้แล้ว พาธที่ไม่ใช่โฟลเดอร์ และชื่อที่ถูกใช้ไปแล้ว ccam ยังเตือนเมื่อโฟลเดอร์เป็นสำเนาของ `~/.claude` (การเข้าสู่ระบบเดียวกัน) แม้ว่า `~/.claude` เองจะยังไม่ได้ลงทะเบียน
- `ccam logout` จะถามก่อนเมื่อบัญชีกำลังทำงานอยู่ ส่วน `ccam login` ไม่ตรวจสอบ
- หากตัวส่งต่อติดต่อ Claude Code ไม่ได้ เช่น หลังจากการเข้าสู่ระบบสิ้นสุดไปแล้ว จะพิมพ์ `ccam: could not pass the sign-in to Claude Code:` พร้อมเหตุผล ข้อผิดพลาด HTTP มักหมายความว่ารหัสหมดอายุแล้ว ให้เริ่มเข้าสู่ระบบใหม่

## หัวข้อที่เกี่ยวข้อง

- [การเริ่ม Claude Code](/features/launch)
- [คำสั่งลัด](/features/shortcuts)
- [ใช้เซิร์ฟเวอร์ MCP การตั้งค่า และประวัติร่วมกัน](/docs/sharing)
- [ใช้ Claude Code หลายบัญชีบนคอมพิวเตอร์เครื่องเดียว](/guides/multiple-claude-code-accounts)
