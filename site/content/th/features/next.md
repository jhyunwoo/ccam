---
title: การสลับไปยังบัญชีถัดไป
description: "ccam next เปิด Claude Code ด้วยบัญชีอื่น --auto เลือกบัญชีที่เหลือโควตามากที่สุด --continue พาการสนทนาไปต่อ และ --fallback รันงาน -p ซ้ำ"
updated: 2026-11-01
order: 8
translation: machine
sourceHash: "79300c521748"
---

`ccam next` เปิด Claude Code ด้วยบัญชีอื่นโดยที่คุณไม่ต้องระบุชื่อ ใช้เมื่อบัญชีถึงขีดจำกัดการใช้งาน การสลับทำได้สี่วิธี:

- `ccam next` แบบธรรมดาจะเลือกบัญชีถัดจากบัญชีที่คุณใช้ล่าสุด และ `--auto` จะเลือกบัญชีที่เหลือโควตามากที่สุด `--continue` จะสนทนาล่าสุดของโฟลเดอร์นี้ต่อในบัญชีใหม่ ไม่ว่าทั้งสองบัญชีจะใช้ประวัติร่วมกันหรือไม่
- ภายใน Claude Code [ปลั๊กอิน ccam](#inside-claude-code) จะบอกคุณเมื่อขีดจำกัดทำให้เซสชันหยุด และ `/ccam:switch` จะส่งต่อการสนทนาไป
- สำหรับสคริปต์ที่รัน `claude -p` นั้น `ccam --fallback` จะสลับบัญชีให้เอง (ดู [รันอีกครั้งเมื่อถึงขีดจำกัด](#run-again-at-the-limit))
- เมื่อเปิด [`auto_switch`](#auto-switch) ไว้ (Linux และ macOS) ccam จะย้ายเซสชันแบบโต้ตอบที่ Claude Code หยุดเพราะถึงขีดจำกัดให้ด้วย

**ที่นั่ง (seat) ของ Team และ Enterprise** การเลือกบัญชีอัตโนมัติ (`ccam next`, `--fallback`, `ccam pick`, ปุ่ม `n` ในแดชบอร์ด และ [`auto_switch`](/features/next#auto-switch) เมื่อคุณเปิดใช้) จะสลับระหว่างที่นั่งต่าง ๆ ในองค์กร Claude Team หรือ Enterprise เดียวกันด้วย ซึ่งอาจละเมิดข้อกำหนดของ Anthropic โปรดอ่าน[คำถามที่พบบ่อย](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) และสอบถามองค์กรของคุณและ Anthropic ก่อนที่จะพึ่งพาการเลือกอัตโนมัติ

## การใช้งาน

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

ccam พิมพ์การสลับเป็น `→ <from> → <to>` ทาง standard error พร้อมเหตุผลของการเลือกด้วย `--auto` จากนั้นเริ่ม Claude Code ด้วยบัญชีใหม่แบบเดียวกับ `ccam <name>` `<from>` คือบัญชีที่คุณใช้ล่าสุด หรือเมื่อใช้ `--continue` คือบัญชีที่เป็นเจ้าของการสนทนานั้น อาร์กิวเมนต์หลัง `--` จะส่งต่อให้ `claude` โดยไม่เปลี่ยนแปลง

## ตัวเลือก

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `--auto` | เลือกบัญชีที่เข้าสู่ระบบแล้วตามการใช้งาน (ดูด้านล่าง) โดยข้ามบัญชีที่ใช้ไปถึง 100% | ปิด |
| `--max-usage <1-100>` | ใช้กับ `--auto`: เลือกบัญชีที่ใช้งานต่ำกว่าค่านี้ก่อน | `auto_max_usage` หากไม่มีคือ 100 |
| `--strategy headroom\|expiring\|balance` | ใช้กับ `--auto`: การใช้งาน 5 ชั่วโมงต่ำที่สุด ช่วง 5 ชั่วโมงที่จะรีเซ็ตเร็วที่สุด หรือโควตารายสัปดาห์ที่เหลือต่อวันมากที่สุด | `auto_strategy` หากไม่มีคือ `headroom` |
| `--cross-org` | ใช้กับ `--auto`: เลือกบัญชีขององค์กรอื่นด้วย (ดูด้านล่าง) | `cross_org` หากไม่มีคือปิด |
| `--wait[=<duration>]` | ใช้กับ `--auto`: เมื่อไม่มีบัญชีใดเหลือโควตา ให้รอบัญชีแรกที่รีเซ็ต นานไม่เกินค่านี้ ค่าต้องอยู่หลัง `=`: `--wait 2h` จะถูกปฏิเสธพร้อมคำแนะนำ `--wait takes its value after '=': --wait=2h` | ปิด; `5h` เมื่อระบุโดยไม่มีค่า |
| `-c`, `--continue` | สนทนาล่าสุดของโฟลเดอร์นี้ต่อในบัญชีใหม่ (ดูด้านล่าง) | ปิด |
| `--brief` | ใช้กับ `--continue`: เริ่มการสนทนาใหม่ในบัญชีใหม่จากสรุปการสนทนาครั้งล่าสุดที่ ccam เขียนขึ้น แทนการเปิดการสนทนาทั้งหมดต่อ ([ดูด้านล่าง](#brief)) | ปิด |
| `-- <claude args>` | อาร์กิวเมนต์ที่ส่งให้ `claude` ต่อจากอาร์กิวเมนต์ของ ccam เอง | ไม่มี |
| `--dry-run` | รองรับ: พิมพ์ `→ <from> → <to>` และบรรทัดคำสั่งที่จะรัน โดยไม่เริ่ม Claude Code และไม่บันทึกการเริ่ม | ปิด |

`--auto` และ `--continue` ใช้ร่วมกันได้ เมื่อใช้ `--continue` อาร์กิวเมนต์ `--resume` จะถูกวางไว้ก่อนอาร์กิวเมนต์ของคุณ

ตัวเลือกของ Claude Code เอง เช่น `--dangerously-skip-permissions` หรือ `--model` ต้องวางไว้หลัง `--` ต่างจาก `ccam <name>` ที่ส่งต่อทุกอย่างหลังชื่อไปให้ `ccam next` อ่านอาร์กิวเมนต์เอง ดังนั้นหากไม่มี `--` จะหยุดพร้อมข้อความ `ccam: unknown flag: --dangerously-skip-permissions`

## วิธีการทำงาน

การเริ่มทุกครั้งผ่าน ccam จะเขียนเวลาลงใน `~/.local/state/ccam/last-used/<name>` `ccam next` อ่านไฟล์เหล่านี้ ตรวจว่าบัญชีใดเข้าสู่ระบบแล้ว แล้วจึงเลือก

1. ccam สร้าง `next.Candidate` หนึ่งรายการต่อบัญชีจาก `config.toml` (`cli.candidates`): ได้แก่ บัญชี เวลาที่เริ่มครั้งล่าสุด (`launcher.LastUsed`) และบัญชีนั้นมีการเข้าสู่ระบบที่ใช้งานได้หรือไม่ (`creds.Has`: การเข้าสู่ระบบ Claude ใน `.credentials.json` หรือบน macOS ในพวงกุญแจ (Keychain) ที่ยังไม่หมดอายุทั้งหมดและไม่ถูก Claude Code ล้างไป) โฟลเดอร์ที่ใช้ API key แทน (การเข้าสู่ระบบ Console หรือบัญชีที่เพิ่มด้วย `ccam add --api`) ไม่มีขีดจำกัดของแพ็กเกจและจะไม่ถูกเลือกเลย
2. หากไม่ใส่ `--auto` ccam จะเรียงบัญชีตามหมายเลขสล็อตและหาบัญชีที่เริ่มล่าสุด (`next.Rotate`) ซึ่งคือบัญชีที่คุณกำลังสลับออก หากยังไม่มีบัญชีใดถูกเริ่มเลย จะใช้สล็อตต่ำสุด
3. จากตรงนั้นจะไล่ไปข้างหน้าตามลำดับสล็อต วนจากสล็อตสุดท้ายกลับไปสล็อตแรก และคืนบัญชีแรกที่เข้าสู่ระบบแล้ว
4. ccam เตรียมและเริ่มบัญชีที่เลือกเหมือนกับที่ `ccam <name>` ทำทุกประการ (`launcher.LaunchWith`): ตรวจลิงก์ของบัญชี ใช้การตั้งค่า MCP ที่ใช้ร่วมกัน ลงทะเบียนเซสชัน และแทนที่ ccam ด้วย `claude` ผ่าน `syscall.Exec`

การไล่เริ่มที่บัญชีที่ใช้ล่าสุด และหยุดที่บัญชีแรกที่เข้าสู่ระบบแล้วซึ่งอยู่ถัดไป:

<!-- motion:next-rotate -->

`work` ถูกเริ่มล่าสุด การไล่จึงเริ่มต่อจากบัญชีนี้ `side` ไม่มีการเข้าสู่ระบบจึงถูกข้าม และเลือก `lab` หาก `lab` ออกจากระบบไปแล้วด้วย การไล่จะวนกลับไปที่ `personal` (`% len` ในแผนภาพ) การหมุนเวียนไม่มีทางคืนบัญชีที่คุณเริ่มต้นมา: หากไม่มีบัญชีอื่นที่เข้าสู่ระบบแล้ว จะหยุดพร้อมข้อผิดพลาด

### --auto

`--auto` ตัดสินจากการใช้งานแทนลำดับสล็อต ccam เก็บบัญชีที่เข้าสู่ระบบแล้วไว้ สอบถามการใช้งานของทุกบัญชีพร้อมกันแบบเดียวกับที่ `ccam usage` ทำ (`Fetcher.FetchAll`) และเลือกตามกฎด้านล่าง (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

บัญชีจะถูกตัดสินจากช่วง 5 ชั่วโมงและรายสัปดาห์ และเมื่อมี `--model` อยู่ในอาร์กิวเมนต์สำหรับ Claude Code (หลัง `--`) ก็จะดูช่วงรายสัปดาห์ของโมเดลนั้นด้วย (`opus`, `sonnet`, `haiku` หรือ `fable` ในชื่อ) ccam แบ่งบัญชีออกเป็นสามระดับ (`Policy.Choose()`):

1. **มีโควตาเหลือ**: ทุกช่วงต่ำกว่า `auto_max_usage` (`config.toml` ค่าเริ่มต้น 100) กลยุทธ์จะเลือกจากกลุ่มนี้: `headroom` (ค่าเริ่มต้น) เลือกบัญชีที่ใช้งาน 5 ชั่วโมงต่ำที่สุด `expiring` เลือกช่วง 5 ชั่วโมงที่จะรีเซ็ตเร็วที่สุด เพื่อให้ใช้โควตาที่เหลือก่อนที่จะหายไป ส่วน `balance` เลือกบัญชีที่เหลือโควตารายสัปดาห์ต่อวันมากที่สุดจนกว่าสัปดาห์ของบัญชีนั้นจะรีเซ็ต โดยนับเซสชัน ccam ที่กำลังทำงานแต่ละเซสชันหักออกจากบัญชีนั้นด้วย โควตารายสัปดาห์จึงไม่ถูกปล่อยทิ้ง และเซสชันที่ทำงานขนานกันจะกระจายออกไป
2. **เกิน**: ต่ำกว่า 100% แต่เท่ากับหรือเกิน `auto_max_usage` ใช้เฉพาะเมื่อไม่มีบัญชีใดมีโควตาเหลือ โดยบัญชีที่ใช้งาน 5 ชั่วโมงต่ำที่สุดจะมาก่อน
3. **เต็ม**: มีช่วงใดช่วงหนึ่งอยู่ที่ 100% หรือไม่ทราบการใช้งาน 5 ชั่วโมง ไม่ถูกเลือกเด็ดขาด

ช่วงหนึ่งจะนับเป็น 100% จนกว่าจะรีเซ็ต เมื่อ Claude Code ปฏิเสธคำขอในช่วงนั้น: hook ของ[ปลั๊กอิน ccam](#inside-claude-code)จะบันทึกไว้ (`~/.local/state/ccam/usage/<name>.limit.json`) ไม่ว่าตัวเลขที่บันทึกไว้จะเป็นเท่าใด ดังนั้นบัญชีที่เพิ่งถึงขีดจำกัดจะไม่ถูกเลือกอีกก่อนรีเซ็ต แม้ตัวเลขของบัญชีจะเก่าแล้วก็ตาม

การตั้งค่าสำหรับทุกบัญชีต้องอยู่ที่ด้านบนสุดของ `config.toml` เหนือ `[[group]]`, `[[profile]]` หรือ `[[folder]]` ตัวแรก การตั้งค่าที่อยู่ใต้ตารางใดตารางหนึ่งจะเป็นของตารางนั้น: ccam จะเตือนเรื่องนี้ทุกครั้งที่เริ่ม และคำสั่งที่บันทึกการตั้งค่าจะปฏิเสธจนกว่าจะย้ายมันขึ้นไป

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**องค์กร** การเลือกอัตโนมัติจะอยู่ในบรรดาบัญชีประเภทเดียวกับบัญชีที่คุณกำลังออกมา: ที่นั่งในองค์กร Team หรือ Enterprise เดียวกัน หรือแพ็กเกจส่วนตัว (Pro, Max) งานจากบัญชีบริษัทจึงไม่ไปทำต่อในบัญชีส่วนตัว และในทางกลับกันก็เช่นกัน ccam รู้องค์กรของแต่ละบัญชีจากการเข้าสู่ระบบ (`ccam login` บันทึกไว้ ส่วนการเข้าสู่ระบบรุ่นเก่าจะอ่านจาก `.claude.json`) บัญชีที่ ccam บอกองค์กรไม่ได้จะจับคู่กับบัญชีใดก็ได้ โฟลเดอร์สองโฟลเดอร์ที่เข้าสู่ระบบด้วยบัญชี Claude เดียวกันใช้ขีดจำกัดร่วมกัน จึงไม่มีโฟลเดอร์ใดรับช่วงต่อจากอีกโฟลเดอร์หนึ่ง `--cross-org` หรือ `cross_org = true` ใน `config.toml` จะยกเลิกกฎเรื่ององค์กรนี้ แต่หลัง [`/ccam:switch`](#inside-claude-code) การย้ายไปยังองค์กรอื่นยังคงต้องรอให้คุณตอบ yes

เมื่อเลือกบัญชีใดไม่ได้ ccam จะออกด้วยรหัส 75 (`EX_TEMPFAIL`) เพื่อให้สคริปต์แยก "ลองใหม่ภายหลัง" ออกจากความล้มเหลวได้ เมื่อใช้ `--wait` ccam จะพักรอจนกว่าช่วงที่ปิดกั้นของบัญชีที่เต็มบัญชีแรกจะรีเซ็ต (30 วินาทีหลังเวลารีเซ็ต ซึ่งเป็นตอนที่ endpoint แสดงช่วงใหม่) พิมพ์ `ccam: no account has room; waiting until 15:10` แล้วเลือกใหม่อีกครั้ง หากไม่มีบัญชีใดบอกเวลารีเซ็ต จะกลับมาดูอีกครั้งหลัง 10 นาที หากไม่มีบัญชีใดว่างก่อนหมดเวลารอ จะออกด้วยรหัส 75

การสอบถามทำตามกฎของ `ccam usage`: ตัวเลขที่มีอายุไม่ถึง 60 วินาทีจะถูกนำกลับมาใช้ การสอบถามของแต่ละบัญชีจะหมดเวลาหลัง 3 วินาที บัญชีที่ไม่ตอบกลับหรือโทเค็นการเข้าถึงหมดอายุแล้วจะถูกตัดสินจากตัวเลขที่บันทึกไว้ และบันทึกขีดจำกัดจะมีผลเหนือตัวเลขเหล่านั้น ช่วงเวลาที่บันทึกไว้ซึ่งเลยเวลารีเซ็ตไปแล้วจะนับเป็น 0% บัญชีที่ช่วงเวลารีเซ็ตไปแล้วจึงไม่ถูกข้าม `--auto` ไม่เคยต่ออายุโทเค็น ต่างจากการหมุนเวียนแบบธรรมดา `--auto` อาจเลือกบัญชีที่คุณใช้ล่าสุดได้หากบัญชีนั้นยังเหลือโควตามากที่สุด

### --continue

1. ccam ค้นหาการสนทนาล่าสุดของโฟลเดอร์นี้ในทุกบัญชี: จากไฟล์ `<session>.jsonl` ใน `projects/<folder>/` (ชื่อโฟลเดอร์คือโฟลเดอร์ที่ทำงานอยู่หลังแปลงลิงก์เป็นพาธจริงแล้ว โดยเปลี่ยนทุกอักขระที่ไม่ใช่ตัวอักษรหรือตัวเลขเป็น `-` ตามที่ Claude Code ตั้งชื่อ) จะเลือกไฟล์ที่มีพรอมต์หรือคำตอบล่าสุด เซสชันที่ไม่มีทั้งสองอย่าง เช่น เซสชันที่รันเพียง `/mcp`, `/usage` หรือ `/status` หรือไฟล์ว่าง จะถูกข้าม หากเซสชันทิ้งบันทึกไว้ (จากขีดจำกัดการใช้งานหรือ `/ccam:switch` ดูด้านล่าง) และไม่มีการสนทนาใดที่มีพรอมต์หรือคำตอบหลังจากบันทึกนั้นเกิน 2 นาที จะใช้เซสชันนั้นแทน แถวที่ Claude Code เขียนตอน `/exit` ไม่นับ
2. บัญชีที่เป็นเจ้าของการสนทนาคือบัญชีที่คุณสลับออก และบัญชีใหม่จะถูกเลือกตามข้างต้นภายในองค์กรเดียวกัน เมื่อหลายบัญชีใช้ประวัตินั้นร่วมกัน การสนทนาจะเป็นของบัญชีที่ `.claude.json` ระบุว่ารันการสนทนานั้นล่าสุด (`lastSessionId`) มิฉะนั้นจะเป็นของบัญชีที่ใช้ล่าสุด
3. หากบัญชีใหม่ใช้ประวัตินั้นร่วมกัน (โฟลเดอร์กลุ่มเดียวกัน) ccam จะเริ่ม `claude --resume <session id>` ซึ่งก็คือการสนทนานั้นเอง (`--continue` จะเปิดการสนทนาที่ใหม่ที่สุดในโฟลเดอร์ ไม่ว่าจะเป็นการสนทนาใด)
4. หากไม่ได้ใช้ร่วมกัน ccam จะเริ่ม `claude --resume <the conversation's file> --fork-session`: Claude Code จะคัดลอกการสนทนาทั้งหมดภายใต้ session ID ใหม่ไปยังประวัติของบัญชีใหม่เอง และต้นฉบับจะยังคงเหมือนเดิม จากนั้น `--continue` ของบัญชีใหม่จะพบการสนทนานั้น
5. หากคุณไม่ได้ใส่ `--agent` หลัง `--` ccam จะส่งต่อเอเจนต์ที่การสนทนาใช้อยู่ด้วย ได้แก่การตั้งค่า `agent` ใน `.claude/settings.local.json` หรือ `.claude/settings.json` ของโฟลเดอร์ หรือใน `settings.json` ของบัญชี หากไม่มีก็จะใช้เอเจนต์ที่แถว `agent-setting` แถวสุดท้ายของการสนทนาระบุ Claude Code ไม่ได้นำเอเจนต์กลับมาเมื่อคัดลอกการสนทนา (ขั้นตอนที่ 4) และจะทำงานต่อในฐานะเอเจนต์เริ่มต้นของตัวเองพร้อมเครื่องมือทุกตัว เมื่อบัญชีใหม่ไม่ได้ใช้เอเจนต์ ปลั๊กอิน และการตั้งค่าร่วมกับบัญชีเดิม (โฟลเดอร์ `agents` และ `plugins` และ `settings.json` เดียวกันผ่านกลุ่ม) ccam จะเพิ่มอีกบรรทัดต่อจากบรรทัดที่บอกเรื่องการย้าย: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`

ccam จะบอกว่าการสนทนามีขนาดเท่าไร: บัญชีใหม่จะอ่านการสนทนาทั้งหมดอีกครั้งในคำตอบแรก prompt cache ไม่เคยใช้ร่วมกันระหว่างองค์กร และแพ็กเกจส่วนตัวแต่ละแพ็กเกจก็เป็นองค์กรของตัวเอง ดังนั้นเมื่อสองบัญชีอยู่ต่างองค์กรกัน ccam จะเพิ่มข้อความ `prompt caches are not shared between organizations, so the first reply reads it all again` เมื่อบัญชีเดิมลงชื่อเข้าใช้เซิร์ฟเวอร์ MCP ที่บัญชีใหม่ยังไม่ได้ลงชื่อเข้าใช้ ccam จะระบุชื่อเซิร์ฟเวอร์เหล่านั้น: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)` การใช้ `history` ร่วมกันในกลุ่มช่วยเลี่ยงการคัดลอก ทำให้เป็นการสนทนาเดียว และยังใช้ checkpoint ที่ `/rewind` ใช้กู้คืนโค้ดร่วมกันด้วย (`file-history/`)

### --brief

`--continue` จะเปิดการสนทนาทั้งหมดต่อ และบัญชีใหม่จะอ่านการสนทนาทั้งหมดอีกครั้งในคำตอบแรก ส่วน `ccam next --continue --brief` จะเริ่มการสนทนาใหม่แทน จากสรุปที่ ccam เขียนขึ้นจากการสนทนาเดิมโดยไม่ถามโมเดล:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

สรุปนี้มีข้อมูลตามลำดับต่อไปนี้ และเฉพาะส่วนที่มีเนื้อหา:

1. ชื่อของการสนทนา
2. พาธของ transcript เพื่อให้ Claude ค้นหาเมื่อต้องการรายละเอียด
3. คำขอแรกของคุณ เมื่อมีคำขอมากกว่าสี่ข้อ
4. สรุปที่ Claude Code เขียนไว้เมื่อ compact การสนทนาครั้งล่าสุด
5. คำขอสี่ข้อล่าสุดของคุณ โดยคงคำสั่ง `!` และ slash command ไว้ตามที่คุณพิมพ์
6. คำตอบล่าสุดของ Claude
7. รายการสิ่งที่ต้องทำ (todo)
8. ไฟล์ 10 ไฟล์ที่ Claude แก้ไขล่าสุด
9. `git status --short --branch` ของโฟลเดอร์ ไม่เกิน 800 ตัวอักษร

รายการสิ่งที่ต้องทำจะเก็บไว้ไม่เกิน 10 รายการ: รายการที่เสร็จแล้วจะถูกตัดออกก่อน แล้วจึงเป็นขั้นตอนที่อยู่ไกลที่สุด ขั้นตอนปัจจุบันและขั้นตอนถัดไปจึงยังอยู่ ส่วนที่ยาวจะถูกตัดตรงกลาง และสรุปทั้งหมดจะไม่เกิน 16,000 ตัวอักษร

ccam บันทึกสรุปเป็น `~/.local/state/ccam/brief/<session>.md` ซึ่งอ่านได้เฉพาะคุณ และ `--brief` ครั้งถัดไปจะลบสรุปที่เก่ากว่า 7 วัน Claude Code จะเริ่มด้วย prompt เดียว และแนบไฟล์ที่ prompt กล่าวถึง: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` prompt ที่คุณให้หลัง `--` จะแทนที่ประโยคสุดท้าย (`Then: <your prompt>`) และตัวเลือกอื่นของคุณหลัง `--` จะตามมาหลังจากนั้น

Claude Code จะจบการกล่าวถึงไฟล์ที่ `"` และตัดที่ `#` ดังนั้นเมื่อพาธของสรุปมีอักขระใดอักขระหนึ่งนี้ ccam จะกล่าวถึงไฟล์ด้วยพาธที่สัมพันธ์กับโฟลเดอร์ทำงาน และเมื่อพาธนั้นก็มีอักขระเหล่านี้ด้วย prompt จะระบุชื่อไฟล์แทน แล้ว Claude จะอ่านไฟล์ด้วยเครื่องมือ Read ซึ่งจะถามก่อนเพราะไฟล์อยู่นอกโฟลเดอร์ทำงาน (ในการรัน `-p` ที่อยู่ในโหมด Manual นั้น Claude Code จะปฏิเสธการอ่าน และ Claude จะทำต่อโดยไม่มีสรุป)

`ccam --dry-run next --continue --brief` จะพิมพ์สรุปออกมาโดยไม่บันทึกอะไร การสนทนาเดิมยังคงอยู่เหมือนเดิม และการสนทนาใหม่รู้เพียงสิ่งที่อยู่ในสรุปเท่านั้น จึงควรใช้ `--brief` เมื่องานสำคัญกว่ารายละเอียดทุกอย่างของสิ่งที่ผ่านมา หากไม่มี `--continue` ccam จะหยุดพร้อม `--brief works with --continue` และเมื่ออ่านการสนทนาไม่ได้หรือไม่พบสิ่งใดให้สรุป ccam จะหยุดและบอกว่า `ccam next --continue` แบบไม่มี `--brief` จะเปิดการสนทนาทั้งหมดต่อ [การส่งต่อ](#relay) และ [`auto_switch`](#auto-switch) จะเปิดการสนทนาทั้งหมดต่อเสมอ ตัวเลือก `--brief` ของ Claude Code เองต้องวางไว้หลัง `--` เหมือนตัวเลือกอื่น ๆ ของ Claude Code

## ภายใน Claude Code

[ปลั๊กอิน ccam](/docs/install#claude-code-plugin) เพิ่ม hook สองตัว (Claude Code 2.1.78 ขึ้นไป):

- เมื่อเทิร์นหนึ่งหยุดเพราะถึงขีดจำกัดการใช้งาน hook จะจดการสนทนาไว้และแสดงการแจ้งเตือนบนเดสก์ท็อปในเทอร์มินัลที่รองรับ (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.``

  hook รู้จักทุกข้อความที่ Claude Code นับว่าเป็นขีดจำกัดการใช้งาน: ขีดจำกัด 5 ชั่วโมง รายสัปดาห์ แบบทั่วไป และของโมเดล (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` สำหรับโควตาของโมเดล) เครดิตการใช้งานและขีดจำกัดการใช้จ่าย และองค์กรที่หยุดบัญชีไว้

  เมื่อ Claude Code เขียนคำขอที่ถูกปฏิเสธลงในการสนทนาแล้ว hook จะบันทึกช่วงที่เต็มไว้จนกว่าจะรีเซ็ตด้วย ([บันทึกขีดจำกัด](/features/usage#the-limit-record)) ยกเว้นการหยุดที่องค์กรตั้งไว้ ซึ่งผู้ดูแลยกเลิกได้ทุกเมื่อ การจำกัดความเร็วสั้น ๆ ของเซิร์ฟเวอร์หรือ rate limit ของ API key ไม่ใช่ขีดจำกัดการใช้งานและจะไม่ถูกแตะต้อง เช่นเดียวกับการหยุดของ subagent และเซสชันที่ไม่มีคนดูแล (การรัน `-p` และเซสชันเบื้องหลัง)

  Claude Code 2.1.234 ขึ้นไปยังรออยู่ในเซสชันและทำต่อเองหลังช่วง 5 ชั่วโมงรีเซ็ตด้วย การแจ้งเตือนจึงให้คุณเลือกสลับบัญชีแทนได้ และเมื่อเปิด [`auto_switch`](#auto-switch) ไว้ (Linux และ macOS) ccam จะสลับให้คุณเมื่อมีอีกบัญชีหนึ่งที่ยังมีโควตาเหลือ เมื่อไม่มีบัญชีอื่นใดรับช่วงต่อได้ เพราะไม่มีบัญชีอื่นที่อยู่ในองค์กร Team หรือ Enterprise เดียวกัน หรือไม่มีบัญชีใดเป็นแพ็กเกจส่วนตัวอีกแพ็กเกจหนึ่ง (เมื่อเปิด `cross_org` บัญชี Claude อื่นใดก็รับช่วงต่อได้) การแจ้งเตือนจะลงท้ายด้วย `You have no other personal plans, so ccam cannot move this conversation.` หรือ `You have no other accounts in Acme, so ccam cannot move this conversation.` แทน
- `/ccam:switch [account]` จะจดการสนทนาไว้และจบเทิร์นก่อนที่จะมีคำขอใด ๆ: พิมพ์ `/exit` จากนั้นรัน `ccam next --continue` เพื่อสนทนาต่อในบัญชีที่คุณระบุ (หรือในบัญชีถัดไป) ในเซสชันที่ยังไม่มีการสนทนา คำสั่งนี้จะแจ้งว่าไม่มีอะไรให้สนทนาต่อ มีเพียงสิ่งที่คุณพิมพ์เองเท่านั้นที่รันคำสั่งนี้ได้: ตั้งแต่ปลั๊กอิน 0.3.2 Claude รัน `/ccam:switch` ไม่ได้ และงานที่ Claude ตั้งเวลาไว้ (CronCreate, ScheduleWakeup หรือ `.claude/scheduled_tasks.json`) จะส่งคำสั่งนี้ให้ Claude เป็นข้อความธรรมดา

  เมื่อบัญชีที่คุณระบุอยู่นอกองค์กรของการสนทนา (ที่นั่งขององค์กร Team หรือ Enterprise อื่น แพ็กเกจส่วนตัวแทนที่นั่งหรือในทางกลับกัน หรือบัญชี API key ฝั่งใดฝั่งหนึ่ง) `ccam next --continue` จะถามก่อน ไม่ว่า `cross_org` จะตั้งไว้อย่างไร: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]` และจะถามแบบเดียวกันก่อนย้ายไปยังบัญชีที่ ccam เลือกเองเมื่อคุณไม่ได้ระบุบัญชี ซึ่ง `cross_org = true` ทำให้เลือกจากองค์กรอื่นได้: `ccam: the account picked for /ccam:switch is home, a personal plan, …` จะย้ายก็ต่อเมื่อตอบ `y` หรือ `yes` เท่านั้น บน Linux และ macOS ปุ่มที่พิมพ์ไว้ก่อนคำถามจะถูกทิ้งไปก่อน คำตอบอื่นใด หรือการไม่มีเทอร์มินัลให้ถาม จะทำให้การสนทนาและบันทึกอยู่ที่เดิม: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (หาก ccam เลือกบัญชีเอง `ccam next --continue` จะเลือกใหม่อีกครั้ง) เช่นเดียวกับการเลือกอัตโนมัติ บัญชีที่ ccam บอกองค์กรไม่ได้จะจับคู่กับบัญชีใดก็ได้

บันทึกนี้อยู่ใน `~/.local/state/ccam/handoff/<account>/<session>.json` หนึ่งไฟล์ต่อเซสชัน สองเซสชันของบัญชีเดียวกันจึงไม่เขียนทับบันทึกของกันและกัน บันทึกจะถูกเก็บไว้ 12 ชั่วโมงและถูกใช้เพียงครั้งเดียวเมื่อเซสชันถัดไปกำลังจะเริ่ม หากการเริ่มล้มเหลว บันทึกจะยังอยู่ให้ลองครั้งถัดไป บันทึกเป็นของโฟลเดอร์ที่เซสชันเริ่มต้น แม้จะ `cd` ไปที่อื่นในเชลล์ของ Claude Code แล้วก็ตาม `ccam rename` จะอัปเดตชื่อในบันทึก หากบัญชีที่บันทึกระบุไว้ถูกลบไปแล้ว ccam จะแจ้งเรื่องนี้และเลือกบัญชีตามปกติ

### การส่งต่อ (relay)

เมื่อตั้ง `relay = true` ที่ด้านบนสุดของ `config.toml` ccam จะทำงานอยู่ควบคู่กับเซสชันแบบโต้ตอบแทนการแทนที่ตัวเองด้วย Claude Code และหลัง `/ccam:switch` ccam จะสนทนาของเซสชันนั้นต่อเอง (บันทึกที่เซสชันอื่นทิ้งไว้จะไม่ถูกนำมาใช้): hook จะจบ Claude Code แล้ว ccam จะรัน `ccam next --continue` พร้อมบันทึกของเซสชันนั้น หากไม่ได้ระบุบัญชี จะเลือกบัญชีที่เหลือโควตามากที่สุด โดยนับโมเดลของการสนทนาด้วย และเมื่อ `cross_org = true` ทำให้บัญชีนั้นเป็นบัญชีขององค์กรอื่น ccam จะถามคำถามข้างต้นก่อนย้ายการสนทนาไปยังบัญชีนั้น บน Windows hook จบ Claude Code ไม่ได้: ให้พิมพ์ `/exit` แล้ว ccam จะสนทนาต่อให้ และ hook ก็ไม่จบ Claude Code สำหรับบัญชีที่อยู่นอกองค์กรของการสนทนาเช่นกัน: hook จะแสดง `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).` และ ccam จะถามคำถามข้างต้นเมื่อคุณพิมพ์ `/exit` แล้ว วิธีนี้ต้องมีโปรเซส ccam ที่รออยู่หนึ่งโปรเซสต่อเซสชัน จึงปิดไว้เป็นค่าเริ่มต้น

เซสชันที่ถูกส่งต่อจะคงสิ่งที่การเปิดการสนทนาต่อไม่ได้คืนให้: ตัวเลือกอย่าง `--add-dir`, `--mcp-config` รายการเครื่องมือ `--settings`, `--plugin-dir` และตัวเลือกของ system prompt รวมถึงโหมดสิทธิ์ของการสนทนา ระดับ effort ของเซสชัน และเอเจนต์ที่การสนทนาใช้อยู่ (ดู [`--continue`](#continue)) Claude Code เองจะคืน `--model` และ `--name` จากการสนทนา ส่วนอย่างอื่น (`--worktree`, `--debug`, `--remote-control` และตัวเลือกที่ ccam ไม่รู้จัก) จะถูกระบุใน `ccam: not carried over: …` การสลับจะไม่ถามคำถามเรื่องประวัติที่ใช้ร่วมกัน และสิ่งที่คุณพิมพ์ระหว่างที่การสลับทำงานจะถูกทิ้ง คำถามเดียวที่การสลับถามคือคำถามก่อนย้ายไปยังบัญชีขององค์กรอื่น

การส่งต่อจะสนทนาต่อเฉพาะเมื่อ Claude Code จบลงตามที่ `/ccam:switch` ขอไว้: ไม่ทำเมื่อคุณพิมพ์ `/exit` ก่อน กด Ctrl+C ระหว่างที่กำลังจบ หรือ ccam ได้รับสัญญาณ hangup หรือสัญญาณอื่น และไม่ทำเมื่อเทอร์มินัลหายไปแล้ว บันทึกจะยังอยู่ และ `ccam next --continue` จะสนทนาต่อภายหลังได้ เซสชันที่เริ่มภายใต้การส่งต่อของ ccam เวอร์ชันเก่า หรือเซสชันที่ไม่มีคนดูแล (การรัน `-p` และเซสชันเบื้องหลัง) จะไม่ถูก hook จบ แต่จะได้รับข้อความให้ `/exit` แทน `claude remote-control` ทำงานโดยไม่มีการส่งต่อ

### `auto_switch`

เมื่อตั้ง `auto_switch = true` ที่ด้านบนสุดของ `config.toml` (Linux และ macOS) ccam จะย้ายเซสชันไปยังอีกบัญชีหนึ่งเองเมื่อ Claude Code หยุดเซสชันเพราะถึงขีดจำกัดการใช้งาน: hook ของปลั๊กอินจะเลือกอีกบัญชีหนึ่งที่ยังมีโควตาเหลือแล้วจบ Claude Code จากนั้น ccam ซึ่งทำงานอยู่ควบคู่กับเซสชันจะสนทนาต่อในบัญชีนั้น ในแบบเดียวกับที่[การส่งต่อ](#relay)ทำหลัง `/ccam:switch` ความสามารถนี้ปิดไว้โดยค่าเริ่มต้นและยังเป็นแบบทดลอง

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam จะย้ายเซสชันก็ต่อเมื่อเงื่อนไขต่อไปนี้เป็นจริงทั้งหมด:

1. ccam เป็นผู้เริ่มเซสชันในเทอร์มินัลและทำงานอยู่ควบคู่กับเซสชัน `auto_switch` จะเปิด[การส่งต่อ](#relay)ด้วย ข้อนี้จึงเป็นจริงสำหรับเซสชันแบบโต้ตอบที่คุณเริ่มหลังตั้งค่านี้ เซสชันที่เริ่มก่อนหน้านั้น หรือเริ่มภายใต้ ccam รุ่นเก่า จะได้รับการแจ้งเตือนที่ลงท้ายด้วย ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` เซสชันที่ไม่มีคนดูแล (การรัน `-p` และเซสชันเบื้องหลัง) `claude remote-control` และ Claude Code ที่เริ่มโดยไม่ผ่าน ccam จะไม่ถูกย้ายเด็ดขาด
2. บัญชีเปิดใช้[ปลั๊กอิน ccam](/docs/install#claude-code-plugin) และใช้ Claude Code 2.1.214 ขึ้นไป ซึ่งส่ง ID ของโปรเซสตัวเองให้ hook (`CLAUDE_PID`)
3. Claude Code เขียนคำขอที่ถูกปฏิเสธลงในการสนทนา พร้อมเวลาที่จะรีเซ็ต สำหรับขีดจำกัด 5 ชั่วโมง รายสัปดาห์ Opus หรือ Sonnet ของบัญชีนั้น ขีดจำกัดการใช้จ่าย เครดิตการใช้งานที่หมด การหยุดที่องค์กรตั้งไว้ โควตาของ Fable (`You've reached your Fable limit.`) การจำกัดความเร็วของเซิร์ฟเวอร์ rate limit ของ API key และการหยุดของ subagent จะไม่ทำให้เซสชันถูกย้ายเด็ดขาด
4. ccam รู้แพ็กเกจของบัญชี จาก `.claude.json` ของบัญชีหรือจากสิ่งที่ `ccam login` บันทึกไว้ หากไม่รู้ การแจ้งเตือนจะลงท้ายด้วย `ccam does not know personal's plan, so auto_switch leaves this session.`
5. เซสชันไม่ได้เชื่อมต่อกับ Remote Control เพราะเมื่อเป็นอีกบัญชีหนึ่ง เซสชันจะไม่เชื่อมต่อกลับไปยังแอปที่คุณใช้ควบคุม เซสชันจึงอยู่ที่เดิม และการแจ้งเตือนจะบอกไว้
6. ขีดจำกัดจะรีเซ็ตในอีกมากกว่า 15 นาทีนับจากตอนนี้ หากต้องรอน้อยกว่านั้น ccam จะปล่อยเซสชันไว้ให้ Claude Code: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. คุณยังไม่ได้พิมพ์ prompt ใดตั้งแต่ถึงขีดจำกัด
8. มีอีกบัญชีหนึ่งที่ยังมีโควตาเหลือ ccam เลือกแบบเดียวกับ [`--auto`](#auto) (`auto_max_usage`, `auto_strategy`, โมเดลของการสนทนา) จากบัญชีประเภทเดียวกับบัญชีที่คุณออกมา: อีกที่นั่งหนึ่งขององค์กร Team หรือ Enterprise เดียวกัน หรืออีกแพ็กเกจส่วนตัวหนึ่ง (Pro, Max) `cross_org` และ `--cross-org` ไม่มีผลในกรณีนี้ บัญชีนั้นต้อง:
   - เข้าสู่ระบบด้วยการเข้าสู่ระบบที่ Claude Code ต่ออายุได้
   - ไม่ได้เข้าสู่ระบบด้วยบัญชี Claude เดียวกัน
   - ไม่ใช้ API key
   - มีแพ็กเกจที่ ccam รู้จัก
   - ไม่ได้เปิดเครดิตการใช้งานไว้ เท่าที่การสอบถามการใช้งานของ ccam เองแสดงให้เห็น
   - เคยเชื่อถือ (trust) โฟลเดอร์นี้ใน Claude Code แล้ว
   - หากการสนทนาทำงานในฐานะเอเจนต์ (การตั้งค่า `agent` หรือการสนทนาระบุเอเจนต์ไว้) ต้องใช้เอเจนต์ร่วมกับบัญชีที่คุณออกมา: โฟลเดอร์ `agents` และ `plugins` และ `settings.json` ของบัญชีนั้นต้องเป็นตัวเดียวกันทุกตัว (ผ่านกลุ่ม) หรือไม่มีทั้งสองฝั่ง

   เมื่อสามเงื่อนไขสุดท้ายนี้ตัดบัญชีที่ยังมีโควตาเหลือออกทั้งหมด การแจ้งเตือนจะระบุชื่อบัญชีเหล่านั้น: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` หรือ `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` ช่วงที่มี[บันทึกขีดจำกัด](/features/usage#the-limit-record)จะนับว่าเต็ม และบัญชีที่เซสชันอื่นเพิ่งย้ายไปภายใน 2 นาทีที่ผ่านมาจะอยู่หลังบัญชีอื่น ๆ การเลือกต้องเสร็จภายในประมาณ 6 วินาที หากการสอบถามการใช้งานใช้เวลานานกว่านั้น เซสชันจะอยู่ที่เดิม เมื่อตั้ง `usage_source = "claude"` `auto_switch` จะตัดสินใจจากตัวเลขที่บันทึกไว้เท่านั้น และ ccam บอกไม่ได้จากการสอบถามของ Claude Code หรือจากแถบสถานะว่าเครดิตการใช้งานเปิดอยู่หรือไม่ จึงถือว่าปิดอยู่

จากนั้น hook จะจบ Claude Code และ ccam ซึ่งทำงานอยู่ควบคู่กันจะแสดงการแจ้งเตือนบนเดสก์ท็อปในเทอร์มินัลที่รองรับ เช่น `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` ccam จะพิมพ์บรรทัดแรกของตัวอย่างเอาต์พุตด้านบน แล้วรัน `ccam next --continue` สำหรับเซสชันนั้นในชื่อบัญชีที่เลือก ซึ่งจะพิมพ์อีกสองบรรทัดของตัวอย่างและเริ่ม Claude Code พร้อมตัวเลือกของเซสชัน โหมดสิทธิ์ ระดับ effort และเอเจนต์ แบบเดียวกับที่[การส่งต่อ](#relay)ทำ ข้อความแรกของเซสชันใหม่คือข้อความตายตัวนี้ ซึ่งเป็นสิ่งเดียวที่ ccam ส่งแทนคุณ (การย้ายหลัง `/ccam:switch` ไม่ส่งอะไรเลย):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

หากบัญชีนั้นถึงขีดจำกัดด้วย ก็จะเกิดแบบเดียวกันที่บัญชีนั้น บันทึกขีดจำกัดจะป้องกันไม่ให้การย้ายต่อกันเป็นทอด ๆ ย้อนกลับไปยังบัญชีก่อนที่บัญชีนั้นจะรีเซ็ต ดังนั้นเมื่อทุกบัญชีเต็ม เซสชันจะอยู่ที่บัญชีสุดท้าย

hook รู้ว่าการสนทนาทำงานในฐานะเอเจนต์จากการตั้งค่า `agent` หรือจากแถวที่ระบุเอเจนต์ซึ่ง Claude Code เขียนลงในการสนทนาเป็นระยะ ๆ และเมื่อเซสชันจบ หากจนถึงตอนนั้นมีเพียง `--agent` ของเซสชันที่ระบุเอเจนต์ และบัญชีที่เลือกไม่ได้ใช้เอเจนต์ร่วมกัน `ccam next --continue` จะไม่เริ่มอะไรเลยและเก็บบันทึกไว้ และการแจ้งเตือนบนเดสก์ท็อปครั้งที่สองจะแจ้งเรื่องนี้: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``

เมื่อไม่มีบัญชีอื่นที่ยังมีโควตาเหลือ จะไม่มีอะไรถูกย้าย: Claude Code จะจัดการกับขีดจำกัดตามปกติ โดยรออยู่ในเซสชันหรือเสนอตัวเลือกของตัวเอง และการแจ้งเตือนจะบอกว่าบัญชีใดจะว่างก่อน เช่น `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` สำหรับที่นั่งจะแสดง `None of your other accounts in Acme has room now` หากต้องการหยุดการย้ายที่เริ่มไปแล้ว ให้กด Ctrl+C ระหว่างที่ Claude Code กำลังปิด: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

การย้ายมีค่าใช้จ่าย Claude Code จะเริ่มใหม่ ซึ่งใช้เวลาไม่กี่วินาที และเชลล์เบื้องหลัง dev server และเอเจนต์ที่ Claude เริ่มไว้จะหยุดไปพร้อมกับเซสชันเดิม บัญชีใหม่จะอ่านการสนทนาทั้งหมดอีกครั้งในคำตอบแรก เหมือนกับ [`--continue`](#continue) และ ccam จะบอกว่าการสนทนามีขนาดเท่าไร การย้ายไม่ใช้ [`--brief`](#brief) เด็ดขาด

`auto_switch` จะไม่สลับก่อนที่ Claude Code จะรายงานขีดจำกัด ไม่เฝ้าดูการใช้งานของคุณอยู่เบื้องหลัง และไม่คัดลอก ต่ออายุ หรือส่งต่อการเข้าสู่ระบบใด ๆ เด็ดขาด บน Windows ค่านี้จะถูกละเว้น และ `ccam doctor` จะแจ้งไว้ บน Windows เมื่อตั้ง `relay = true` แล้ว `/ccam:switch` และ `/exit` ยังย้ายการสนทนาไปยังอีกบัญชีหนึ่งได้

**ข้อกำหนดของ Anthropic** สำหรับที่นั่งของ Team และ Enterprise นั้น `auto_switch` จะย้ายงานของคุณไปยังอีกที่นั่งหนึ่งขององค์กรเดียวกัน ศูนย์ช่วยเหลือของ Anthropic ระบุว่า "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (แต่ละที่นั่งในแพ็กเกจ Team จะกำหนดให้กับบุคคลที่ระบุชื่อเพียงคนเดียว และมีเพียงบุคคลนั้นที่ใช้ได้) และ "Users can't route their work through another user's account to use that user's remaining usage" (ผู้ใช้ไม่สามารถส่งงานของตนผ่านบัญชีของผู้ใช้คนอื่นเพื่อใช้ปริมาณการใช้งานที่เหลืออยู่ของผู้ใช้คนนั้นได้) ([แพ็กเกจ Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129) ส่วน [แพ็กเกจ Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) ก็เช่นเดียวกัน)

สำหรับแพ็กเกจส่วนตัว (Pro, Max) ข้อกำหนดของ Anthropic ไม่ได้ระบุว่าบัญชีของคุณเองรับช่วงต่อจากกันโดยอัตโนมัติเมื่อถึงขีดจำกัดได้หรือไม่ [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) ห้ามการเข้าถึงบริการ "through automated or non-human means, whether through a bot, script, or otherwise" (ด้วยวิธีอัตโนมัติหรือไม่ใช่มนุษย์ ไม่ว่าจะผ่านบอต สคริปต์ หรือวิธีอื่น) เว้นแต่ Anthropic จะอนุญาตไว้อย่างชัดแจ้ง และการย้ายแต่ละครั้งจะส่งข้อความตายตัวข้างต้นแทนคุณ ส่วน [Usage Policy](https://www.anthropic.com/legal/aup) ที่จะมีผลในวันที่ 12 พฤศจิกายน 2026 ระบุ "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (จงใจหลบเลี่ยงหรือแทรกแซงความสามารถ ข้อจำกัด หรือมาตรการป้องกันที่กำหนดไว้ในผลิตภัณฑ์ของเรา) ไว้ในสิ่งที่ห้ามทำ

ccam ไม่ได้อ้างว่า `auto_switch` เป็นไปตามข้อกำหนดเหล่านี้ การเปิดใช้เป็นการตัดสินใจของคุณ และความเสี่ยงก็เป็นของคุณเช่นกัน: Anthropic อาจระงับบัญชี หรือการเข้าถึงขององค์กร Team หรือ Enterprise เมื่อเห็นว่ามีการละเมิดข้อกำหนด

ขณะที่ `auto_switch` เปิดอยู่ `ccam doctor` จะเตือนเกี่ยวกับบัญชีที่ไม่มีปลั๊กอิน ccam และ Claude Code เวอร์ชันที่เก่ากว่า 2.1.214 และสำหรับที่นั่งจะย้ำข้อกำหนดเรื่องที่นั่งของ Anthropic ([การตรวจ](/features/doctor))

## รันอีกครั้งเมื่อถึงขีดจำกัด

`ccam next` ใช้สำหรับเซสชันแบบโต้ตอบที่หยุดไปแล้ว การรันแบบไม่โต้ตอบ เช่น `claude -p` ในสคริปต์หรืองานตามกำหนดเวลา รอคุณไม่ได้ การใส่ `--fallback` ไว้หน้าชื่อบัญชีจึงให้ ccam สลับบัญชีแทนให้:

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

1. `--fallback` ใช้ได้เฉพาะร่วมกับ `-p` หรือ `--print` เท่านั้น มิฉะนั้น ccam จะหยุดและแนะนำ `ccam next --auto --continue` แทน
2. อินพุตมาตรฐานที่ส่งผ่านไปป์จะถูกอ่านครั้งเดียว ไม่เกิน 10 MiB และส่งให้ทุกครั้งที่ลอง เช่นเดียวกับ Claude Code ccam จะรอให้อินพุตเริ่มเข้ามาไม่เกิน 3 วินาที (ยกเว้นเมื่อใช้ `--input-format stream-json`) หากไปป์เปิดค้างไว้แต่ไม่ส่งอะไรมา การรันจะไม่มีอินพุต และ ccam จะแจ้งเรื่องนี้ทาง standard error
3. ccam เริ่ม `claude` เป็นโปรเซสลูกแทนการใช้ `exec` กักเอาต์พุตมาตรฐานไว้ก่อน และส่ง standard error ผ่านออกไปทันทีที่ได้รับ
4. เมื่อ `claude` ออกด้วยรหัส 0 หรือล้มเหลวด้วยสาเหตุอื่น เอาต์พุตที่กักไว้และรหัสออกจะถูกส่งต่อโดยไม่เปลี่ยนแปลง การรันที่จบเพราะสัญญาณ (รหัสออกเกิน 128) ก็เช่นกัน: การรันนั้นถูกขอให้หยุด และจะไม่ถูกรันอีก
5. ความล้มเหลวจะนับว่าเป็นขีดจำกัดการใช้งานเมื่อข้อความของ Claude Code ระบุเช่นนั้น ซึ่งคือข้อความขีดจำกัดการใช้งานใดก็ได้ที่ [hook](#inside-claude-code) รู้จัก รวมถึงภายในเอาต์พุต `--output-format json` หรือ `stream-json` หรือเมื่อการใช้งานของบัญชีแสดงว่าช่วง 5 ชั่วโมง รายสัปดาห์ หรือของ `--model` อยู่ที่ 100% (รวมถึงบันทึกขีดจำกัด)
6. จากนั้นเอาต์พุตของการลองครั้งนั้นจะถูกทิ้ง และคำสั่งเดิมจะรันด้วยบัญชีอื่นที่เข้าสู่ระบบอยู่: บัญชีที่ใช้งาน 5 ชั่วโมงน้อยที่สุด ตามด้วยบัญชีที่ดูการใช้งานไม่ได้ เรียงตามหมายเลข บัญชีที่อยู่ที่ 100% จะถูกข้าม บัญชีที่เข้าสู่ระบบด้วย API key จะไม่ถูกลองเลย และแต่ละบัญชีจะถูกลองไม่เกินหนึ่งครั้ง
7. หากมี `-c`, `--continue`, `-r` หรือ `--resume` อยู่ในอาร์กิวเมนต์ จะลองเฉพาะบัญชีที่อยู่ในกลุ่มประวัติเดียวกัน เพื่อให้หาการสนทนาเจอ
8. จะลองเฉพาะบัญชีขององค์กรเดียวกัน เช่นเดียวกับ `--auto` (`cross_org = true` จะยกเลิกข้อนี้)

บัญชีถัดไปจะถูกเลือกด้วยกฎเดียวกับ `--auto` (`auto_max_usage`, `auto_strategy` และช่วงของ `--model`) เมื่อไม่มีบัญชีเหลือ ccam จะพิมพ์เอาต์พุตของการลองครั้งสุดท้าย แจ้งว่าไม่มีบัญชีที่เข้าสู่ระบบอยู่บัญชีใดเหลือโควตา และออกด้วยรหัส 75 เมื่อใส่ `--wait` ไว้ก่อนชื่อบัญชี (`ccam --fallback --wait=2h work -p …`) ccam จะรอให้บัญชีแรกรีเซ็ตแทน แล้วอาจลองทุกบัญชีอีกครั้ง โดยจะพักรออย่างน้อยหนึ่งครั้งก่อนเสมอ บัญชีที่เพิ่งถึงขีดจำกัดจึงไม่ถูกรันซ้ำทันที ระยะเวลานี้จำกัดเวลารอทั้งหมดของการรัน ไม่ใช่การรอแต่ละครั้ง หากเขียนโดยเว้นวรรค (`--wait 30m`) จะถูกปฏิเสธพร้อมคำแนะนำเดียวกับ `ccam next` ส่วน `--` ที่นำหน้าอาร์กิวเมนต์ของ Claude Code ใช้ได้เหมือนการเริ่มทุกแบบ (`ccam --fallback work -- -p …`) `ccam --dry-run --fallback work -p …` จะพิมพ์คำสั่งของการลองครั้งแรก

เซสชันแบบโต้ตอบจะไม่ถูกสลับเอง เว้นแต่คุณจะเปิด [`auto_switch`](#auto-switch) เมื่อใช้ `ccam statusline install` แถบสถานะจะเพิ่ม `→ ccam next --auto --continue` เมื่อการใช้งาน 5 ชั่วโมงหรือรายสัปดาห์ถึง 90% ([ขีดจำกัดการใช้งาน](/features/usage))

## ความปลอดภัยและกรณีพิเศษ

- บัญชีที่ไม่มีการเข้าสู่ระบบจะไม่ถูกเลือกเลย หากไม่มีบัญชีอื่นที่เข้าสู่ระบบแล้ว `ccam next` แบบธรรมดาจะหยุดด้วย `no logged-in account other than <name>`
- เมื่อใช้ `--continue` และไม่มีการสนทนาในโฟลเดอร์: `no account has a conversation in <folder> to continue`
- เมื่อกฎเรื่ององค์กรทำให้ไม่เหลือบัญชีให้เลือก: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``
- เมื่อใช้ `--auto` และไม่มีบัญชีที่เข้าสู่ระบบแล้ว: `no logged-in account to pick` เมื่อทุกบัญชีเต็มหรือไม่ทราบการใช้งาน: `` no account with known usage and headroom (see `ccam usage`) `` ทั้งสองกรณีจบด้วยรหัสออก 75
- `--max-usage`, `--strategy`, `--wait` และ `--cross-org` ใช้ได้เฉพาะกับ `--auto` และค่าที่อยู่นอกช่วง 1-100 หรือกลยุทธ์ที่ไม่รู้จักจะถูกปฏิเสธ
- fork คัดลอกการสนทนา ไม่ได้คัดลอกสถานะอื่นของบัญชี: checkpoint ใน `file-history/` ของบัญชีเดิมจะไม่ตามมาด้วย เว้นแต่จะใช้ประวัติร่วมกัน
- `--auto` ส่งคำขอผ่านเครือข่ายแม้จะใส่ `--dry-run` โดยจะอ่านเท่านั้น และไม่เคยเขียนไฟล์โทเค็น
- บัญชีที่เลือกจะเริ่มผ่านเส้นทางเดียวกับ `ccam <name>` จึงผ่านการตรวจแบบเดียวกัน รวมถึงคำเตือนเมื่อมีเซสชัน Claude Code อื่นที่ใช้ประวัติเดียวกันกำลังทำงานอยู่ในโฟลเดอร์นี้
- `last-used` จะถูกเขียนเฉพาะเมื่อ Claude Code เริ่มทำงานจริง `--dry-run` จะไม่แตะไฟล์นี้ การทดลองรันจึงไม่ทำให้การหมุนเวียนเลื่อนไปข้างหน้า
- `--fallback` กักเอาต์พุตมาตรฐานไว้จนกว่าการลองแต่ละครั้งจะจบ เอาต์พุตแบบสตรีมอย่าง `--output-format stream-json` จึงมาถึงทีเดียวตอนท้าย ไคลเอนต์ที่สตรีมอินพุต `stream-json` และรอคำตอบก่อนส่งเพิ่มจะใช้ `--fallback` ไม่ได้ ไฟล์ที่การลองครั้งแรกแก้ไขไปก่อนถึงขีดจำกัดจะยังคงถูกแก้ไขอยู่ การลองครั้งถัดไปจะเริ่มพรอมต์ใหม่ตั้งแต่ต้น

## หัวข้อที่เกี่ยวข้อง

- [ขีดจำกัดการใช้งาน](/features/usage)
- [ขีดจำกัดการใช้งานและการสลับบัญชี](/docs/usage-limits)
- [ใช้เซิร์ฟเวอร์ MCP การตั้งค่า และประวัติร่วมกัน](/docs/sharing)
- [ทำอย่างไรเมื่อ Claude Code ถึงขีดจำกัดการใช้งาน](/guides/usage-limit-workflow)
