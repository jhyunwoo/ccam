---
title: Beralih ke akun berikutnya
description: ccam next membuka Claude Code sebagai akun lain; --auto memilih ruang terbanyak, --continue membawa percakapan, dan --fallback mengulang tugas -p.
updated: 2026-11-05
order: 8
translation: machine
sourceHash: "bc9837b7a6ae"
---

`ccam next` membuka Claude Code sebagai akun lain tanpa Anda perlu menyebutkan namanya. Gunakan saat sebuah akun mencapai batas penggunaannya. Ada empat cara untuk berpindah:

- `ccam next` biasa mengambil akun setelah akun yang terakhir Anda pakai, dan `--auto` mengambil akun dengan ruang terbanyak. `--continue` melanjutkan percakapan terakhir di folder ini di akun baru, baik kedua akun berbagi riwayat maupun tidak.
- Di dalam Claude Code, [plugin ccam](#inside-claude-code) memberi tahu Anda saat sebuah batas menghentikan sesi, dan `/ccam:switch` meneruskan percakapan.
- Untuk skrip yang menjalankan `claude -p`, `ccam --fallback` melakukan pergantian itu dengan sendirinya (lihat [Jalankan ulang saat batas tercapai](#run-again-at-the-limit)).
- Dengan [`auto_switch`](#auto-switch) diaktifkan (Linux dan macOS), ccam juga memindahkan sesi interaktif yang dihentikan Claude Code di batasnya.

**Kursi Team dan Enterprise.** Pemilihan akun otomatis (`ccam next`, `--fallback`, `ccam pick`, tombol `n` di dasbor, dan, setelah Anda mengaktifkannya, [`auto_switch`](/features/next#auto-switch)) juga berganti di antara kursi dalam organisasi Claude Team atau Enterprise yang sama, yang dapat melanggar ketentuan Anthropic. Baca [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) dan tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya.

## Penggunaan

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

ccam menampilkan perpindahan sebagai `→ <from> → <to>` di error standar, beserta alasan pilihan `--auto`, lalu menjalankan Claude Code sebagai akun baru dengan cara yang sama seperti `ccam <name>`. `<from>` adalah akun yang terakhir Anda pakai, atau dengan `--continue`, akun pemilik percakapan itu. Argumen setelah `--` diteruskan ke `claude` tanpa diubah.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--auto` | Memilih akun yang sudah login berdasarkan penggunaan (lihat di bawah), melewati akun yang sudah 100% | nonaktif |
| `--max-usage <1-100>` | Dengan `--auto`: utamakan akun di bawah penggunaan ini | `auto_max_usage`, jika tidak, 100 |
| `--strategy headroom\|expiring\|balance` | Dengan `--auto`: penggunaan 5 jam terendah, reset 5 jam yang paling cepat, atau sisa kuota mingguan terbanyak per hari | `auto_strategy`, jika tidak, `headroom` |
| `--cross-org` | Dengan `--auto`: juga memilih akun dari organisasi lain (lihat di bawah) | `cross_org`, jika tidak, nonaktif |
| `--wait[=<duration>]` | Dengan `--auto`: jika tidak ada akun yang punya ruang, tunggu akun pertama di-reset, paling lama selama ini. Nilainya ditaruh setelah `=`: `--wait 2h` ditolak dengan petunjuk `--wait takes its value after '=': --wait=2h` | nonaktif; `5h` jika diberikan tanpa nilai |
| `-c`, `--continue` | Melanjutkan percakapan terakhir di folder ini di akun baru (lihat di bawah) | nonaktif |
| `--brief` | Dengan `--continue`: memulai percakapan baru sebagai akun baru dari ringkasan yang dibuat ccam atas percakapan terakhir, alih-alih melanjutkan seluruhnya ([lihat di bawah](#brief)) | nonaktif |
| `-- <claude args>` | Argumen yang diteruskan ke `claude` setelah argumen milik ccam | tidak ada |
| `--dry-run` | Didukung: menampilkan `→ <from> → <to>` dan baris perintah yang akan dijalankan, tanpa menjalankan Claude Code atau mencatat peluncuran | nonaktif |

`--auto` dan `--continue` dapat dipakai bersama. Dengan `--continue`, argumen `--resume` ditempatkan sebelum argumen Anda sendiri.

Opsi milik Claude Code sendiri, seperti `--dangerously-skip-permissions` atau `--model`, ditaruh setelah `--`. Tidak seperti `ccam <name>`, yang meneruskan semua yang ada setelah nama, `ccam next` membaca argumennya sendiri, jadi tanpa `--` ccam berhenti dengan `ccam: unknown flag: --dangerously-skip-permissions`.

## Cara kerjanya

Setiap peluncuran melalui ccam menulis waktunya ke `~/.local/state/ccam/last-used/<name>`. `ccam next` membaca file-file ini, memeriksa akun mana yang sudah login, lalu memilih.

1. ccam membuat satu `next.Candidate` per akun dari `config.toml` (`cli.candidates`): akun, waktu peluncuran terakhirnya (`launcher.LastUsed`), dan apakah akun itu memiliki login yang berfungsi (`creds.Has`: login Claude di `.credentials.json`, atau di macOS di Rantai Kunci (Keychain), yang belum kedaluwarsa sepenuhnya atau dihapus oleh Claude Code). Folder yang memakai kunci API sebagai gantinya (login Console, atau akun yang ditambahkan dengan `ccam add --api`) tidak memiliki batas paket dan tidak pernah dipilih.
2. Tanpa `--auto`, ccam mengurutkan akun berdasarkan nomor slot dan mencari akun yang paling baru diluncurkan (`next.Rotate`). Itulah akun yang Anda tinggalkan. Jika belum ada akun yang pernah diluncurkan, itu adalah slot terkecil.
3. Dari sana fungsi ini bergerak maju sesuai urutan slot, kembali dari slot terakhir ke slot pertama, dan mengembalikan akun pertama yang sudah login.
4. ccam menyiapkan dan menjalankan akun yang dipilih persis seperti `ccam <name>` (`launcher.LaunchWith`): memeriksa tautan akun, menerapkan pembagian MCP, mendaftarkan sesi, dan mengganti ccam dengan `claude` melalui `syscall.Exec`.

Penelusuran dimulai dari akun yang terakhir dipakai dan berhenti di akun pertama yang sudah login setelahnya:

<!-- motion:next-rotate -->

`work` paling baru diluncurkan, sehingga penelusuran dimulai setelahnya. `side` tidak memiliki login dan dilewati, lalu `lab` dipilih. Seandainya `lab` juga belum login, penelusuran akan kembali ke `personal` (`% len` di diagram). Rotasi tidak pernah mengembalikan akun asal: jika tidak ada akun lain yang sudah login, rotasi berhenti dengan galat.

### --auto

`--auto` memutuskan berdasarkan penggunaan alih-alih urutan slot. ccam mempertahankan akun yang sudah login, menanyakan penggunaan semuanya sekaligus dengan cara yang sama seperti `ccam usage` (`Fetcher.FetchAll`), lalu memilih berdasarkan aturan di bawah ini (`cli.pickAuto`, `usage.Policy.Choose()`):

<!-- motion:next-auto -->

Akun dinilai berdasarkan jendela 5 jam dan mingguannya, dan, jika `--model` ada di antara argumen untuk Claude Code (setelah `--`), juga berdasarkan jendela mingguan model itu (`opus`, `sonnet`, `haiku`, atau `fable` di namanya). ccam membagi akun menjadi tiga tingkatan (`Policy.Choose()`):

1. **Ada ruang**: setiap jendela di bawah `auto_max_usage` (`config.toml`, bawaan 100). Strategi memilih di antaranya: `headroom` (bawaan) mengambil, di antara akun yang penggunaan 5 jamnya paling banyak 3 poin di atas penggunaan 5 jam terendah, akun dengan penggunaan mingguan paling sedikit, sehingga akun di 12% mengalahkan akun di 10% jika penggunaan mingguannya lebih sedikit, tetapi akun di 14% tidak; `expiring` mengambil jendela 5 jam yang paling cepat di-reset, sehingga sisa kuotanya terpakai sebelum hangus; `balance` mengambil akun dengan sisa kuota mingguan terbanyak per hari hingga minggunya di-reset, dengan menghitung setiap sesi ccam yang sedang berjalan sebagai beban akun itu, sehingga kuota mingguan tidak dibiarkan tak terpakai dan sesi paralel tersebar.
2. **Lewat batas**: di bawah 100% tetapi sama dengan atau di atas `auto_max_usage`. Hanya dipakai jika tidak ada akun yang punya ruang; akun yang akan diambil `headroom` di antara mereka didahulukan.
3. **Penuh**: sebuah jendela di 100%, atau penggunaan 5 jam yang tidak diketahui. Tidak pernah dipilih.

Sebuah jendela dihitung 100% sampai di-reset jika Claude Code menolak permintaan di sana: hook [plugin ccam](#inside-claude-code) mencatatnya (`~/.local/state/ccam/usage/<name>.limit.json`), apa pun kata angka yang tersimpan. Jadi akun yang baru saja mencapai batasnya tidak dipilih lagi sebelum di-reset, bahkan jika angkanya sudah lama.

Pengaturan untuk semua akun ditaruh di bagian atas `config.toml`, di atas `[[group]]`, `[[profile]]`, atau `[[folder]]` pertama. Pengaturan di bawah salah satunya menjadi milik tabel itu: ccam memperingatkannya di setiap peluncuran, dan perintah yang menyimpan konfigurasi menolak sampai pengaturan itu dipindahkan.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organisasi.** Pilihan otomatis tetap berada di antara akun yang sejenis dengan akun yang Anda tinggalkan: kursi di organisasi Team atau Enterprise yang sama, atau paket pribadi (Pro, Max), sehingga pekerjaan dari akun perusahaan tidak berlanjut di akun pribadi, atau sebaliknya. ccam mengetahui organisasi setiap akun dari login-nya (`ccam login` mencatatnya; login lama dibaca dari `.claude.json`); akun yang organisasinya tidak bisa dipastikan ccam dapat dicocokkan dengan akun mana pun. Dua folder yang login ke akun Claude yang sama berbagi batasnya, sehingga yang satu tidak pernah menggantikan yang lain. `--cross-org`, atau `cross_org = true` di `config.toml`, mencabut aturan organisasi ini, tetapi setelah [`/ccam:switch`](#inside-claude-code), perpindahan ke organisasi lain tetap menunggu jawaban yes dari Anda.

Jika tidak ada yang bisa dipilih, ccam keluar dengan kode 75 (`EX_TEMPFAIL`), sehingga skrip bisa membedakan "coba lagi nanti" dari kegagalan. Dengan `--wait`, ccam justru tidur sampai jendela penghalang milik akun penuh pertama di-reset (30 detik setelah waktu reset, saat endpoint menampilkan jendela baru), mencetak `ccam: no account has room; waiting until 15:10`, lalu memilih lagi; jika tidak ada akun yang menyebutkan kapan di-reset, ccam memeriksa lagi setelah 10 menit. Jika tidak ada yang kosong sebelum waktu tunggu berakhir, ccam keluar dengan 75.

Pencarian mengikuti aturan `ccam usage`: angka yang berumur kurang dari 60 detik dipakai ulang, pencarian setiap akun habis waktunya setelah 3 detik, akun yang tidak menjawab, atau yang token aksesnya sudah kedaluwarsa, dinilai berdasarkan angka yang tersimpan, dan catatan batas menggantikan angka-angka itu. Jendela tersimpan yang waktu reset-nya sudah lewat dihitung 0%, sehingga akun yang jendelanya sudah direset sejak itu tidak dilewatkan. `--auto` tidak pernah memperbarui token. Tidak seperti rotasi biasa, `--auto` dapat memilih akun yang terakhir Anda pakai jika akun itu masih memiliki ruang terbanyak.

### --continue

1. ccam mencari percakapan terakhir di folder ini di setiap akun: di antara file `<session>.jsonl` di `projects/<folder>/` (nama foldernya adalah folder kerja, dengan tautan di-resolve, dan setiap karakter selain huruf atau angka diubah menjadi `-`, seperti cara Claude Code menamainya), yang memiliki prompt atau balasan terbaru. Sesi tanpa prompt atau balasan, seperti sesi yang hanya menjalankan `/mcp`, `/usage`, atau `/status`, atau file kosong, dilewati. Jika sebuah sesi meninggalkan catatan (batas penggunaan atau `/ccam:switch`, lihat di bawah) dan tidak ada percakapan yang memiliki prompt atau balasan lebih dari 2 menit setelahnya, sesi itulah yang dipakai; baris yang ditulis Claude Code saat `/exit` tidak dihitung.
2. Akun pemilik percakapan itu adalah akun yang Anda tinggalkan, dan akun baru dipilih seperti di atas, di dalam organisasi yang sama. Jika beberapa akun berbagi riwayat itu, percakapan itu milik akun yang `.claude.json`-nya menyatakan bahwa akun itulah yang terakhir menjalankan percakapan tersebut (`lastSessionId`), atau jika tidak, milik akun yang terakhir dipakai.
3. Jika akun baru berbagi riwayat itu (folder grup yang sama), ccam menjalankan `claude --resume <session id>`, tepat percakapan itu (`--continue` akan membuka percakapan mana pun yang terbaru di folder tersebut).
4. Jika tidak, ccam menjalankan `claude --resume <the conversation's file> --fork-session`: Claude Code menyalin seluruh percakapan, dengan ID sesi baru, ke riwayat milik akun baru, dan aslinya tetap seperti semula. Sejak itu `--continue` di akun baru akan menemukannya.
5. Kecuali Anda memberikan `--agent` setelah `--`, ccam juga meneruskan agen yang dipakai percakapan itu: pengaturan `agent` dari `.claude/settings.local.json` atau `.claude/settings.json` milik folder itu atau dari `settings.json` milik akun, atau jika tidak ada, agen yang disebutkan oleh baris `agent-setting` terakhir dalam percakapan. Claude Code tidak mengembalikan agen saat menyalin percakapan (langkah 4), dan akan melanjutkan sebagai agen default-nya, dengan semua alat. Jika akun baru tidak berbagi agen, plugin, dan pengaturan akun lama (folder `agents` dan `plugins` serta `settings.json` yang sama, melalui grup), ccam menambahkan satu baris setelah baris tentang perpindahan: `work does not share personal's agents, plugins and settings, so it may not have agent "reviewer" or may define it differently: the conversation may go on without that agent's tool limits`.

ccam menyebutkan seberapa besar percakapannya: akun baru membaca seluruhnya lagi pada balasan pertamanya. Cache prompt tidak pernah dibagikan antarorganisasi, dan setiap paket pribadi adalah organisasi tersendiri, jadi jika kedua akun berasal dari organisasi yang berbeda, ccam menambahkan `prompt caches are not shared between organizations, so the first reply reads it all again`. Jika akun lama sudah masuk ke server MCP yang belum dimasuki akun baru, ccam menyebutkannya: `work is not signed in to MCP servers that personal is signed in to: linear (sign in with /mcp there; sign-ins are not shared)`. Berbagi `history` dalam sebuah grup menghindari penyalinan, mempertahankan satu percakapan, dan juga membagikan checkpoint tempat `/rewind` memulihkan kode (`file-history/`).

### --brief

`--continue` melanjutkan seluruh percakapan, dan akun baru membaca semuanya lagi pada balasan pertamanya. `ccam next --continue --brief` justru memulai percakapan baru, dari ringkasan percakapan lama yang ditulis ccam tanpa bertanya kepada model:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

Ringkasan itu memuat, dengan urutan ini dan hanya jika ada yang perlu disampaikan:

1. judul percakapan;
2. path transkripnya, agar Claude bisa mencarinya saat membutuhkan suatu detail;
3. permintaan pertama Anda, jika ada lebih dari empat;
4. ringkasan yang ditulis Claude Code saat terakhir kali memadatkan (compact) percakapan;
5. empat permintaan terakhir Anda, dengan perintah `!` dan perintah garis miring seperti yang Anda ketik;
6. balasan terakhir Claude;
7. daftar tugas (todo);
8. 10 file yang terakhir diubah Claude;
9. `git status --short --branch` folder itu, hingga 800 karakter.

Daftar tugas menyimpan paling banyak 10 item: item yang sudah selesai dibuang lebih dulu, lalu langkah yang paling jauh di depan, sehingga langkah saat ini dan langkah berikutnya tetap ada. Bagian yang panjang dipotong di tengah, dan seluruh ringkasan tetap di bawah 16.000 karakter.

ccam menyimpannya sebagai `~/.local/state/ccam/brief/<session>.md`, yang hanya bisa dibaca oleh Anda, dan `--brief` berikutnya menghapus ringkasan yang lebih lama dari 7 hari. Claude Code dimulai dengan satu prompt, dan melampirkan file yang disebut di dalamnya: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` Prompt yang Anda berikan setelah `--` menggantikan kalimat terakhir (`Then: <your prompt>`), dan opsi Anda yang lain setelah `--` mengikutinya. Claude Code mengakhiri sebuah mention di `"` dan memotongnya di `#`, jadi jika path ringkasan mengandung salah satunya, ccam menyebutnya dengan path relatif dari folder kerja; jika path itu juga mengandung salah satunya, prompt menyebutkan nama filenya, dan Claude membacanya dengan alat Read-nya, yang meminta izin lebih dulu karena file itu berada di luar folder kerja (dalam proses `-p` dengan mode Manual, Claude Code menolak pembacaan itu dan Claude melanjutkan tanpa ringkasan).

`ccam --dry-run next --continue --brief` mencetak ringkasan dan tidak menyimpan apa pun. Percakapan lama tetap seperti semula, dan percakapan baru hanya mengetahui apa yang dikatakan ringkasan, jadi gunakan `--brief` jika tugasnya lebih penting daripada setiap detail tentang bagaimana Anda sampai di sana. Tanpa `--continue`, ccam berhenti dengan `--brief works with --continue`; jika ccam tidak bisa membaca percakapan atau tidak menemukan apa pun di dalamnya untuk diringkas, ccam berhenti dan memberi tahu bahwa `ccam next --continue` tanpa `--brief` melanjutkannya secara utuh. [Relay](#relay) dan [`auto_switch`](#auto-switch) selalu melanjutkan seluruh percakapan. Opsi `--brief` milik Claude Code sendiri ditaruh setelah `--`, seperti opsinya yang lain.

## Di dalam Claude Code

[Plugin ccam](/docs/install#claude-code-plugin) menambahkan dua hook (Claude Code 2.1.78 atau yang lebih baru):

- Saat sebuah giliran berhenti di batas penggunaan, hook mencatat percakapan itu dan menampilkan notifikasi desktop di terminal yang mendukungnya (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` continues the conversation on another account.`` Hook mengenali setiap pesan yang dianggap Claude Code sebagai batas penggunaan: batas 5 jam, mingguan, biasa, dan model (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` untuk jatah sebuah model), kredit penggunaan dan batas pengeluaran, serta organisasi yang menghentikan akun. Jika Claude Code sudah menulis permintaan yang ditolak itu ke percakapan, hook juga mencatat jendela yang penuh sampai di-reset ([catatan batas](/features/usage#the-limit-record)), kecuali untuk penghentian yang ditetapkan oleh organisasi, yang bisa dicabut adminnya kapan saja. Pembatasan singkat dari server atau batas laju (rate limit) kunci API bukanlah batas penggunaan dan dibiarkan, begitu pula sesi tanpa pengawasan (proses `-p` dan sesi latar belakang). Subagen atau agen workflow yang berhenti di batas penggunaan dihitung sebagai penghentian sesinya, bahkan saat percakapan utama hanya menunggunya: Claude Code menulis permintaan yang ditolak itu ke file milik agen itu sendiri di folder `subagents` sesi tersebut, dan hook membacanya di sana. Jika beberapa agen berhenti sekaligus, sesi itu hanya mendapat satu catatan dan satu notifikasi. Claude Code 2.1.234 dan yang lebih baru juga menunggu di dalam sesi dan melanjutkan dengan sendirinya setelah reset 5 jam; notifikasi itu memberi Anda pilihan untuk berpindah akun saja, dan dengan [`auto_switch`](#auto-switch) aktif (Linux dan macOS), ccam berpindah untuk Anda jika akun lain punya ruang. Jika tidak ada akun lain yang boleh mengambil alih, karena tidak ada akun lain yang termasuk dalam organisasi Team atau Enterprise yang sama, atau tidak ada paket pribadi lain (dengan `cross_org` aktif, akun Claude lain mana pun boleh), notifikasi itu justru diakhiri dengan `You have no other personal plans, so ccam cannot move this conversation.` atau `You have no other accounts in Acme, so ccam cannot move this conversation.`
- `/ccam:switch [account]` mencatat percakapan dan mengakhiri giliran sebelum permintaan apa pun dibuat: ketik `/exit`, lalu jalankan `ccam next --continue` untuk melanjutkannya di akun yang Anda sebutkan (atau di akun berikutnya). Di sesi yang belum memiliki percakapan, perintah ini memberi tahu bahwa tidak ada yang bisa dilanjutkan. Hanya yang Anda ketik yang menjalankannya: sejak plugin 0.3.2, Claude tidak dapat menjalankan `/ccam:switch`, dan tugas yang dijadwalkan Claude (CronCreate, ScheduleWakeup, atau `.claude/scheduled_tasks.json`) meneruskannya ke Claude sebagai teks biasa.

  Jika akun yang Anda sebutkan berada di luar organisasi percakapan itu (kursi organisasi Team atau Enterprise lain, paket pribadi untuk percakapan di sebuah kursi atau sebaliknya, atau akun kunci API di salah satu sisi), `ccam next --continue` bertanya dulu, apa pun pengaturan `cross_org`: `ccam: /ccam:switch named home, a personal plan, but this conversation is on work, a Claude Team seat in Acme. Move it to home? [y/N]`. ccam juga bertanya hal yang sama sebelum pindah ke akun yang dipilihnya sendiri saat Anda tidak menyebutkan akun, yang dengan `cross_org = true` boleh diambil dari organisasi lain: `ccam: the account picked for /ccam:switch is home, a personal plan, …`. Hanya `y` atau `yes` yang memindahkannya; di Linux dan macOS, tombol yang diketik sebelum pertanyaan dibuang terlebih dahulu. Jawaban lain apa pun, atau tidak adanya terminal untuk bertanya, membiarkan percakapan dan catatannya tetap di tempatnya: ``the conversation stays on work, a Claude Team seat in Acme: it moves to home, a personal plan, only when you answer yes on a terminal. `ccam next --continue` asks again, and `ccam work --resume <session id>` continues it as work`` (setelah ccam memilih sendiri, `ccam next --continue` memilih lagi). Seperti pada pilihan otomatis, akun yang organisasinya tidak bisa dipastikan ccam dapat dicocokkan dengan akun mana pun.

Catatan itu disimpan di `~/.local/state/ccam/handoff/<account>/<session>.json`, satu file per sesi, sehingga dua sesi dari satu akun tidak pernah saling menimpa. Catatan itu disimpan selama 12 jam dan dipakai sekali, saat sesi berikutnya akan dimulai; jika peluncurannya gagal, catatan itu tetap ada untuk percobaan berikutnya. Catatan itu milik folder tempat sesi dimulai, bahkan setelah `cd` di shell Claude Code. `ccam rename` memperbarui nama di dalamnya; jika akun yang disebut di dalamnya sejak itu dihapus, ccam memberitahukannya dan memilih seperti biasa.

### Relay

Dengan `relay = true` di bagian atas `config.toml`, ccam tetap berjalan di samping sesi interaktif alih-alih menggantikan dirinya dengan Claude Code, dan setelah `/ccam:switch` ccam melanjutkan percakapan sesi itu dengan sendirinya (catatan yang ditinggalkan sesi lain tidak ditindaklanjuti): hook mengakhiri Claude Code, dan ccam menjalankan `ccam next --continue` dengan catatan sesi itu. Tanpa nama akun, ccam mengambil akun dengan sisa kuota terbanyak, dengan menghitung model percakapan itu, dan jika berkat `cross_org = true` akun itu berasal dari organisasi lain, ccam mengajukan pertanyaan di atas sebelum memindahkan percakapan ke sana. Di Windows, hook tidak bisa mengakhiri Claude Code: ketik `/exit`, dan ccam melanjutkannya. Hook juga tidak mengakhirinya untuk akun di luar organisasi percakapan itu: hook menampilkan `ccam: type /exit; ccam then asks whether to continue this conversation on home, a personal plan (work is a Claude Team seat in Acme).`, dan ccam mengajukan pertanyaan di atas setelah Anda mengetik `/exit`. Itu memerlukan satu proses ccam yang menunggu per sesi, karena itulah fitur ini nonaktif secara default.

Sesi yang diteruskan mempertahankan apa yang tidak dipulihkan saat melanjutkan percakapan: opsi seperti `--add-dir`, `--mcp-config`, daftar alat, `--settings`, `--plugin-dir`, dan opsi system prompt, mode izin percakapan, tingkat effort sesi, dan agen yang dipakai percakapan itu (lihat [`--continue`](#continue)). Claude Code sendiri memulihkan `--model` dan `--name` dari percakapan. Selain itu (`--worktree`, `--debug`, `--remote-control`, opsi yang tidak dikenal ccam) disebutkan dalam `ccam: not carried over: …`. Perpindahan itu tidak pernah menanyakan pertanyaan riwayat bersama, dan apa yang Anda ketik selama perpindahan berjalan dibuang; satu-satunya pertanyaannya adalah pertanyaan sebelum pindah ke akun organisasi lain.

Relay hanya melanjutkan jika Claude Code berakhir dengan cara yang diminta `/ccam:switch`: tidak jika Anda mengetik `/exit` lebih dulu, menekan Ctrl+C saat Claude Code sedang berakhir, atau ccam menerima hangup atau sinyal lain, dan tidak jika terminalnya sudah tidak ada. Catatan itu lalu tetap ada, dan `ccam next --continue` melanjutkan percakapan itu nanti. Sesi yang dimulai di bawah relay versi ccam yang lebih lama, atau sesi tanpa pengawasan (proses `-p` dan sesi latar belakang), tidak diakhiri oleh hook; sesi itu mendapat pesan `/exit` sebagai gantinya. `claude remote-control` berjalan tanpa relay.

### `auto_switch`

Dengan `auto_switch = true` di bagian atas `config.toml` (Linux dan macOS), ccam memindahkan sesi ke akun lain dengan sendirinya saat Claude Code menghentikannya di batas penggunaan: hook plugin memilih akun lain yang punya ruang dan mengakhiri Claude Code, lalu ccam, yang tetap berjalan di samping sesi, melanjutkan percakapan di sana, seperti yang dilakukan [relay](#relay) setelah `/ccam:switch`. Fitur ini nonaktif secara default dan masih eksperimental.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam memindahkan sesi hanya jika semua hal berikut terpenuhi:

1. ccam memulainya di terminal dan tetap berjalan di sampingnya. `auto_switch` mengaktifkan [relay](#relay), sehingga ini terpenuhi untuk sesi interaktif yang Anda mulai setelah mengaturnya. Sesi yang dimulai sebelumnya, atau di bawah ccam versi lama, mendapat notifikasi yang diakhiri dengan ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` Sesi tanpa pengawasan (proses `-p` dan sesi latar belakang), `claude remote-control`, dan Claude Code yang dimulai tanpa ccam tidak pernah dipindahkan.
2. Akun itu mengaktifkan [plugin ccam](/docs/install#claude-code-plugin) dan menjalankan Claude Code 2.1.214 atau yang lebih baru, yang meneruskan ID prosesnya ke hook (`CLAUDE_PID`).
3. Claude Code menulis permintaan yang ditolak ke percakapan, beserta waktu reset-nya, untuk batas 5 jam, mingguan, Opus, atau Sonnet akun itu. Batas pengeluaran, kredit penggunaan yang habis, penghentian yang ditetapkan oleh organisasi, jatah Fable (`You've reached your Fable limit.`), pembatasan dari server, dan batas laju kunci API tidak pernah memindahkan sesi. Batas subagen atau agen workflow ikut dihitung, bahkan saat percakapan utama hanya menunggunya; jika itu batas mingguan Opus atau Sonnet, sesi hanya dipindahkan jika percakapan itu sendiri berjalan pada model tersebut.
4. ccam mengetahui paket akun itu, dari `.claude.json`-nya atau dari apa yang dicatat `ccam login`. Jika tidak, notifikasi diakhiri dengan `ccam does not know personal's plan, so auto_switch leaves this session.`
5. Sesi tidak terhubung ke Remote Control. Sebagai akun lain, sesi itu tidak akan tersambung kembali ke aplikasi yang Anda pakai untuk mengendalikannya, jadi sesi itu tetap di tempat, dan notifikasi memberitahukannya.
6. Batasnya di-reset lebih dari 15 menit dari sekarang. Untuk waktu tunggu yang lebih singkat, ccam menyerahkan sesi itu kepada Claude Code: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. Anda belum mengetik prompt sejak batas tercapai.
8. Akun lain punya ruang. ccam memilih seperti [`--auto`](#auto) (`auto_max_usage`, `auto_strategy`, model percakapan), di antara akun yang sejenis dengan akun yang Anda tinggalkan: kursi lain dari organisasi Team atau Enterprise yang sama, atau paket pribadi lain (Pro, Max). `cross_org` dan `--cross-org` tidak berlaku di sini. Akun itu harus:
   - login dengan login yang bisa diperbarui Claude Code;
   - tidak login ke akun Claude yang sama;
   - tidak memakai kunci API;
   - memiliki paket yang diketahui ccam;
   - tidak mengaktifkan kredit penggunaan, sejauh yang ditunjukkan oleh pencarian penggunaan milik ccam sendiri;
   - sudah memercayai folder ini di Claude Code;
   - jika percakapan berjalan sebagai agen (pengaturan `agent` atau percakapannya menyebutkan agen), berbagi agen dengan akun yang Anda tinggalkan: folder `agents` dan `plugins` serta `settings.json`-nya masing-masing harus benar-benar sama (melalui grup), atau tidak ada di keduanya.

   Jika tiga syarat terakhir itu menyingkirkan setiap akun yang punya ruang, notifikasi menyebutkannya: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.`, `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` atau `This conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so auto_switch did not move it there.` Jendela dengan [catatan batas](/features/usage#the-limit-record) dihitung penuh, dan akun yang menjadi tujuan perpindahan sesi lain dalam 2 menit terakhir ditempatkan setelah akun lainnya. Pilihan harus dibuat dalam sekitar 6 detik; jika pencarian penggunaan memakan waktu lebih lama, sesi itu tetap di tempat. Dengan `usage_source = "claude"`, `auto_switch` memutuskan hanya dari angka yang tersimpan, dan ccam tidak dapat mengetahui dari pencarian Claude Code atau dari baris status apakah kredit penggunaan aktif, sehingga ccam menganggapnya nonaktif.

Kemudian hook mengakhiri Claude Code, dan ccam, yang tetap berjalan di sampingnya, menampilkan notifikasi desktop di terminal yang mendukungnya, misalnya `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` ccam mencetak baris pertama dari contoh output di atas dan menjalankan `ccam next --continue` untuk sesi itu sebagai akun yang dipilih, yang mencetak dua baris lainnya dan memulai Claude Code dengan opsi, mode izin, tingkat effort, dan agen sesi itu, seperti yang dilakukan [relay](#relay). Pesan pertama sesi baru adalah sebuah teks tetap, satu-satunya hal yang dikirim ccam atas nama Anda (perpindahan setelah `/ccam:switch` tidak mengirim apa pun):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

Setelah batas subagen, pesannya adalah yang ini, karena subagen, workflow, dan perintah latar belakang ikut berhenti bersama sesi lama:

```text
A subagent hit your usage limit and ccam moved this conversation to another account. The subagents, workflows and background commands that were running then have stopped, and their results may be missing: run again what the task still needs, then continue it; do not repeat work that is already complete.
```

Jika akun itu juga mencapai batasnya, hal yang sama terjadi di sana. Catatan batas mencegah rantai itu kembali ke akun sebelum akun itu di-reset, sehingga jika setiap akun penuh, sesi tetap di akun terakhir.

Hook mengetahui bahwa percakapan berjalan sebagai agen dari pengaturan `agent`, atau dari baris yang menyebutkan agen yang ditulis Claude Code ke dalam percakapan sesekali dan saat sesi berakhir. Jika sejauh ini hanya `--agent` sesi itu yang menyebutkannya, dan akun yang dipilih tidak berbagi agen, `ccam next --continue` tidak memulai apa pun dan menyimpan catatannya, dan notifikasi desktop kedua memberitahukannya: ``ccam: this conversation runs as agent "reviewer", and side does not share personal's agents, plugins and settings, so ccam did not move it there. Its note stays: `ccam next --continue` continues it later, and `ccam personal --resume <session id>` continues it as personal after the reset``.

Jika tidak ada akun lain yang punya ruang, tidak ada yang dipindahkan: Claude Code menangani batas itu seperti biasa, menunggu di dalam sesi atau menawarkan pilihannya, dan notifikasi menyebutkan akun mana yang paling dulu kosong, misalnya `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees up first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` Untuk kursi, notifikasinya berbunyi `None of your other accounts in Acme has room now`. Untuk menghentikan perpindahan yang sudah dimulai, tekan Ctrl+C saat Claude Code sedang ditutup: ``ccam: switch cancelled. Run `ccam next --continue` later to continue the conversation.``

Perpindahan ada biayanya. Claude Code dimulai ulang, yang memakan beberapa detik, dan shell latar belakang, server pengembangan, serta agen yang dimulai Claude ikut berhenti bersama sesi lama. Akun baru membaca seluruh percakapan lagi pada balasan pertamanya, seperti dengan [`--continue`](#continue), dan ccam menyebutkan seberapa besar percakapannya; perpindahan tidak pernah memakai [`--brief`](#brief).

`auto_switch` tidak pernah berpindah sebelum Claude Code melaporkan batas, tidak pernah memantau penggunaan Anda di latar belakang, dan tidak pernah menyalin, memperbarui, atau meneruskan login. Di Windows fitur ini diabaikan, dan `ccam doctor` memberitahukannya; di sana, dengan `relay = true`, `/ccam:switch` dan `/exit` tetap memindahkan percakapan ke akun lain.

**Ketentuan Anthropic.** Pada kursi Team dan Enterprise, `auto_switch` memindahkan pekerjaan Anda ke kursi lain dari organisasi yang sama. Pusat bantuan Anthropic menyatakan "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (setiap kursi pada paket Team diberikan kepada satu orang tertentu, dan hanya orang itu yang boleh menggunakannya) dan "Users can't route their work through another user's account to use that user's remaining usage" (pengguna tidak boleh menyalurkan pekerjaannya melalui akun pengguna lain untuk memakai sisa kuota pengguna tersebut) ([paket Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); begitu pula untuk [paket Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)).

Untuk paket pribadi (Pro, Max), ketentuan Anthropic tidak menyebutkan apakah akun-akun Anda sendiri boleh saling mengambil alih secara otomatis saat mencapai batas. [Ketentuan Konsumen](https://www.anthropic.com/legal/consumer-terms) melarang mengakses Layanan "through automated or non-human means, whether through a bot, script, or otherwise" (melalui cara otomatis atau non-manusia, baik melalui bot, skrip, maupun cara lain) kecuali Anthropic secara tegas mengizinkannya, dan sebuah perpindahan mengirim salah satu pesan tetap di atas atas nama Anda. [Kebijakan Penggunaan](https://www.anthropic.com/legal/aup) yang mulai berlaku pada 12 November 2026 mencantumkan "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" (dengan sengaja melewati atau mengganggu kemampuan, pembatasan, atau pengaman yang ditetapkan dalam produk kami) di antara hal-hal yang dilarangnya.

ccam tidak mengklaim bahwa `auto_switch` mematuhi ketentuan-ketentuan ini. Mengaktifkannya adalah keputusan Anda, begitu pula risikonya: Anthropic dapat menangguhkan akun, atau akses organisasi Team atau Enterprise, jika Anthropic menilai ketentuannya dilanggar.

Selama `auto_switch` aktif, `ccam doctor` memperingatkan tentang akun tanpa plugin ccam dan tentang versi Claude Code yang lebih lama dari 2.1.214, dan untuk kursi mengulangi aturan kursi Anthropic ([pemeriksaan](/features/doctor)).

## Jalankan ulang saat batas tercapai

`ccam next` ditujukan untuk sesi interaktif yang sudah berhenti. Proses non-interaktif, yaitu `claude -p` di skrip atau tugas terjadwal, tidak bisa menunggu Anda, jadi `--fallback` sebelum nama akun membuat ccam berganti akun untuknya:

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

1. `--fallback` hanya berfungsi bersama `-p` atau `--print`; jika tidak, ccam berhenti dan mengarahkan ke `ccam next --auto --continue`.
2. Input standar dari pipe dibaca sekali, hingga 10 MiB, dan diberikan ke setiap percobaan. Seperti Claude Code, ccam menunggu paling lama 3 detik sampai input mulai masuk (kecuali dengan `--input-format stream-json`); jika pipe tetap terbuka tetapi tidak mengirim apa pun, proses berjalan tanpa input dan ccam mencatatnya di error standar.
3. ccam menjalankan `claude` sebagai proses anak alih-alih `exec`, menahan output standarnya, dan meneruskan error standar begitu muncul.
4. Saat `claude` keluar dengan 0, atau gagal karena alasan lain, output yang ditahan dan kode keluarnya diteruskan tanpa diubah. Begitu pula proses yang diakhiri sinyal (kode keluar di atas 128): proses itu diminta berhenti, dan tidak dijalankan lagi.
5. Kegagalan dianggap sebagai batas penggunaan jika pesan Claude Code menyatakannya, yaitu salah satu pesan batas penggunaan yang dikenali [hook](#inside-claude-code), juga di dalam output `--output-format json` atau `stream-json`, atau jika penggunaan akun menunjukkan jendela 5 jam, mingguan, atau `--model`-nya di 100% (termasuk catatan batas).
6. Kemudian output percobaan itu dibuang dan perintah yang sama dijalankan sebagai akun lain yang sudah login: akun dengan penggunaan 5 jam terendah, lalu akun yang penggunaannya tidak bisa dicek, sesuai urutan nomor. Akun di 100% dilewati, akun dengan login kunci API tidak pernah dicoba, dan setiap akun dicoba paling banyak sekali.
7. Jika ada `-c`, `--continue`, `-r`, atau `--resume` di antara argumen, hanya akun dalam grup riwayat yang sama yang dicoba, sehingga percakapannya bisa ditemukan.
8. Hanya akun dari organisasi yang sama yang dicoba, seperti pada `--auto` (`cross_org = true` mencabut aturan ini).

Akun berikutnya dipilih dengan aturan yang sama seperti `--auto` (`auto_max_usage`, `auto_strategy`, jendela `--model`). Jika tidak ada akun yang tersisa, ccam mencetak output percobaan terakhir, menyatakan bahwa tidak ada akun yang sudah login yang masih punya ruang, dan keluar dengan kode 75. Dengan `--wait` sebelum nama akun (`ccam --fallback --wait=2h work -p …`), ccam justru menunggu akun pertama di-reset, lalu dapat mencoba setiap akun lagi; ccam selalu tidur setidaknya sekali terlebih dahulu, sehingga akun yang baru saja mencapai batasnya tidak langsung dijalankan lagi. Durasinya membatasi total waktu tunggu proses itu, bukan setiap kali menunggu. Jika ditulis dengan spasi (`--wait 30m`), opsi itu ditolak dengan petunjuk yang sama seperti `ccam next`. `--` di depan argumen Claude Code berfungsi seperti pada peluncuran mana pun (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` mencetak perintah percobaan pertama.

Sesi interaktif tidak berganti akun dengan sendirinya kecuali Anda mengaktifkan [`auto_switch`](#auto-switch). Dengan `ccam statusline install`, baris status menambahkan `→ ccam next --auto --continue` begitu penggunaan 5 jam atau mingguan mencapai 90% ([Batas penggunaan](/features/usage)).

## Keamanan dan kasus khusus

- Akun tanpa login tidak pernah dipilih. Jika tidak ada akun lain yang sudah login, `ccam next` biasa berhenti dengan `no logged-in account other than <name>`.
- Dengan `--continue` dan tidak ada percakapan di folder ini: `no account has a conversation in <folder> to continue`.
- Jika aturan organisasi tidak menyisakan akun apa pun: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Dengan `--auto` dan tidak ada akun yang sudah login: `no logged-in account to pick`. Jika setiap akun penuh atau penggunaannya tidak diketahui: `` no account with known usage and headroom (see `ccam usage`) ``. Keduanya berakhir dengan kode keluar 75.
- `--max-usage`, `--strategy`, `--wait`, dan `--cross-org` hanya berfungsi dengan `--auto`, dan nilai di luar 1-100 atau strategi yang tidak dikenal ditolak.
- Fork menyalin percakapan, bukan status akun lainnya: checkpoint di `file-history/` milik akun asal tidak ikut terbawa kecuali riwayat dibagikan.
- `--auto` membuat permintaan jaringan bahkan dengan `--dry-run`. Opsi ini hanya membaca, dan tidak pernah menulis file token.
- Akun yang dipilih dijalankan melalui jalur yang sama dengan `ccam <name>`, sehingga pemeriksaan yang sama berlaku, termasuk peringatan jika sesi Claude Code lain yang berbagi riwayat yang sama sudah berjalan di folder ini.
- `last-used` hanya ditulis saat Claude Code benar-benar dimulai. `--dry-run` tidak menyentuhnya, sehingga dry run tidak memajukan rotasi.
- `--fallback` menahan output standar sampai sebuah percobaan selesai, sehingga output streaming seperti `--output-format stream-json` tiba sekaligus di akhir. Klien yang mengalirkan input `stream-json` dan menunggu balasan sebelum mengirim lebih banyak tidak dapat memakai `--fallback`. File yang diubah percobaan pertama sebelum mencapai batas tetap berubah; percobaan berikutnya memulai prompt dari awal.

## Terkait

- [Batas penggunaan](/features/usage)
- [Batas penggunaan dan perpindahan akun](/docs/usage-limits)
- [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing)
- [Saat Claude Code mencapai batas penggunaan](/guides/usage-limit-workflow)
