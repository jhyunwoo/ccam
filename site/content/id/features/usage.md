---
title: Batas penggunaan
description: ccam usage menampilkan batas 5 jam dan mingguan semua akun sekaligus, dan ccam statusline menampilkan nama akun dan penggunaannya di Claude Code.
updated: 2026-11-01
order: 7
translation: machine
sourceHash: "54b296f7da35"
---

`ccam usage` menampilkan, untuk setiap akun, seberapa banyak setiap batas yang sudah digunakan: batas 5 jam, batas mingguan, dan batas mingguan khusus setiap model (seperti Fable). Setiap batas mendapat satu bilah beserta waktu hingga direset. Gunakan sebelum memulai tugas panjang, atau saat Claude Code memberi tahu bahwa batas sudah dekat, untuk melihat akun mana yang masih punya ruang. `ccam statusline` adalah perintah status line untuk Claude Code: perintah ini menampilkan nama akun dan penggunaan di dalam Claude Code serta menyimpan angkanya untuk `ccam usage`. Untuk token dan perkiraan biaya di balik penggunaan itu, `ccam ccusage` menjalankan [ccusage](https://ccusage.com) atas percakapan semua akun ([lihat di bawah](#tokens-and-costs-ccam-ccusage)).

## Penggunaan

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts that are not running
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

Setiap akun mendapat satu baris berisi nomornya, namanya, dan umur angka-angkanya; `(cached)` berarti angkanya berasal dari file yang tersimpan, bukan dari endpoint. Di bawahnya, satu baris per batas menampilkan bilah 20 sel, persentase, dan waktu hingga jendela direset: `5-hour`, `Weekly` (semua model), lalu jendela mingguan khusus setiap model, seperti `Fable`. Jendela yang tidak dimiliki akun menampilkan `-`. Persentase berubah kuning pada 80% dan merah pada 100%.

Jika Claude Code menolak permintaan karena sebuah jendela sudah penuh, jendela itu menampilkan 100% sampai di-reset, apa pun kata angka yang tersimpan (lihat [catatan batas](#the-limit-record)). Di terminal, setiap akun langsung muncul sebagai `fetching…`, dan baris-barisnya digambar begitu akun itu menjawab.

Baris terakhir, yang hanya muncul jika Anda memiliki lebih dari satu akun, menyebut akun yang akan dipilih [`ccam next --auto`](/features/next). Baris itu mengikuti aturan yang sama:

- `auto_max_usage` dan `auto_strategy`;
- sesi yang sedang berjalan;
- hanya akun dari organisasi akun yang terakhir dipakai, kecuali `cross_org` diatur;
- setiap jendela yang sudah 100% dilewati, sehingga akun yang kuota mingguannya sudah penuh tidak pernah disebut.

Dengan strategi bawaan `headroom`, baris itu berbunyi `→ Most headroom: <name>`; dengan `expiring` atau `balance`, baris itu memberi alasannya, seperti `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Kursi Team dan Enterprise.** Pemilihan akun otomatis (`ccam next`, `--fallback`, `ccam pick`, tombol `n` di dasbor, dan, setelah Anda mengaktifkannya, [`auto_switch`](/features/next#auto-switch)) juga berganti di antara kursi dalam organisasi Claude Team atau Enterprise yang sama, yang dapat melanggar ketentuan Anthropic. Baca [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) dan tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya.

Akun yang belum login menampilkan `not logged in → ccam login <name>`, dan akun dengan kunci API menampilkan `no limits (API)`. Token akses yang kedaluwarsa menampilkan ``access token expired; `ccam <name>` renews it when it starts``, di atas angka yang tersimpan, jika ada: Claude Code memperbaruinya pada saat dijalankan berikutnya. Login yang sudah kedaluwarsa sepenuhnya, karena refresh token-nya kedaluwarsa atau Claude Code menghapusnya setelah server menolaknya, menampilkan `` login expired → `ccam login <name>` ``.

Setelah `ccam statusline install`, Claude Code menampilkan baris seperti `[work] · Opus · 5h 24% · 7d 41%`.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--json` | Menampilkan laporan sebagai JSON: `name`, `slot`, `status`, `limits` (dengan `credits` jika akun mengaktifkan kredit penggunaan), `source`, `at`, `error`, serta `limited` dan `limited_windows` jika catatan batas menetapkan jendela menjadi 100%. Tanpa akun, perintah ini mencetak `[]` | nonaktif |
| `--watch[=N]` | Menggambar ulang setiap `N` detik sampai Anda menekan Ctrl+C. `--watch` saja, atau nilai 0 atau kurang, berarti 60. Nilainya ditaruh setelah `=`: `--watch 10` ditolak dengan petunjuk | nonaktif |
| `-v`, `--verbose` | Menambahkan baris `source:` di bawah setiap akun: asal angka-angkanya (endpoint penggunaan, Claude Code, status line, file yang tersimpan, atau batas penggunaan yang dilaporkan Claude Code). Jika catatan batas menggantikan angka-angka itu, baris seperti `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` menyusul | nonaktif |
| `--no-cache` | Mengabaikan cache 60 detik dan bertanya ulang ke endpoint | nonaktif |
| `--source live\|claude\|cache` | `live` bertanya ke endpoint dengan token akun; `claude` meminta Claude Code mencarinya (lihat di bawah); `cache` hanya menampilkan angka yang tersimpan | `usage_source`, jika tidak, `live` |
| `--refresh` | Memperbarui token akses yang kedaluwarsa, hanya untuk akun yang tidak sedang berjalan. Satu pembaruan bisa memakan waktu hingga 30 detik | nonaktif |
| `--dry-run` | Diterima: `ccam usage` tidak mengubah apa pun selain cache-nya sendiri. Ditolak jika dipakai bersama `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | nonaktif |

Nilai `--source` lainnya gagal dengan `--source must be live, claude or cache`. `ccam usage` tidak menerima argumen.

### Penggunaan dari Claude Code sendiri

Dengan `usage_source = "claude"` di bagian atas `config.toml` (di atas `[[group]]`, `[[profile]]`, atau `[[folder]]` pertama), atau `--source claude`, ccam sama sekali tidak membaca token login. Untuk setiap akun, ccam menjalankan `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` di folder akun itu dan mengirim permintaan `get_usage` milik Claude Code lewat input standar (Claude Code 2.1.270 atau lebih baru). Claude Code mencari penggunaannya dengan login miliknya sendiri dan menjawab dengan jendela 5 jam, mingguan, dan per model.

`--safe-mode` mencegah hook, server MCP, dan plugin akun itu dimulai, dan tidak ada sesi yang disimpan. Setiap akun memerlukan satu kali start Claude Code (sekitar satu detik), sehingga sumber ini lebih lambat daripada `live`. Semua akun ditanya sekaligus, masing-masing dengan batas waktu 20 detik, dan cache 60 detik berlaku seperti biasa.

Login yang token aksesnya sudah kedaluwarsa, atau akan kedaluwarsa dalam 5 menit, mendapat 50 detik, karena Claude Code memperbaruinya saat dimulai dan pembaruan yang terputus bisa menghilangkan login. Jika ccam menyerah pada sebuah pencarian, ccam meminta Claude Code keluar dengan SIGTERM, sehingga Claude Code menyelesaikan pembaruan yang sudah dimulainya (Windows tidak punya SIGTERM, jadi di sana prosesnya dihentikan paksa).

Login yang sudah kedaluwarsa sepenuhnya menampilkan `` login expired → `ccam login <name>` `` tanpa menjalankan Claude Code. Claude Code menandai permintaan ini sebagai eksperimental, jadi jika permintaan itu berhenti dijawab, ccam menampilkan angka yang tersimpan dan `ccam doctor` melaporkannya. Jika Claude Code menjawab bahwa pencariannya sendiri gagal, ccam menampilkan angka yang tersimpan dengan `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`, bukan paket tanpa batas.

Dengan `--dry-run`, ccam tidak menjalankan Claude Code, karena Claude Code menulis ke folder akun; sebagai gantinya ccam menampilkan angka yang tersimpan. `usage_source` menerima `live` atau `claude`; `cache` hanya bisa dipilih lewat `--source`.

`ccam statusline` membaca JSON status line milik Claude Code dari input standar. Satu-satunya opsinya adalah `--wrap <command>`: perintah status line yang sudah ada, yang baris keluaran pertamanya diberi nama akun di depannya. `ccam statusline install` tidak menerima opsi. Keduanya menolak `--dry-run` dengan `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Token dan biaya: `ccam ccusage`

`ccam usage` menampilkan seberapa banyak setiap batas yang sudah digunakan sebuah akun. Sebaliknya, [ccusage](https://ccusage.com) menghitung token dan perkiraan biaya dalam log percakapan Claude Code. Jika dijalankan sendiri, ccusage hanya mencari log Claude Code di `~/.claude` dan `~/.config/claude` (atau di folder yang ditunjuk `CLAUDE_CONFIG_DIR`), sehingga satu-satunya akun ccam yang dapat dilihatnya adalah akun yang diimpor di tempatnya di `~/.claude`. `ccam ccusage` menjalankannya dengan `CLAUDE_CONFIG_DIR` diatur ke folder setiap akun yang memiliki percakapan, sehingga satu laporan mencakup semuanya:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam membaca opsinya sendiri terlebih dahulu: `-a` (`--account`) menerima nama atau nomor akun, dipisahkan koma atau dengan mengulang opsinya, dan `--dry-run` menampilkan perintahnya alih-alih menjalankannya. Mulai dari argumen pertama yang bukan opsi ccam, semuanya diteruskan ke ccusage persis seperti yang Anda ketik, kecuali `--config` untuk `bun x` dan `npx` (lihat di bawah); `ccam ccusage -- --help` menampilkan bantuan ccusage sendiri. ccam menjalankan `ccusage` jika sudah terinstal, jika tidak `bun x ccusage`, dan jika itu pun gagal `npx ccusage@latest`; dua yang terakhir mengunduhnya dari npm. Keduanya dimulai di folder milik ccam sendiri (`~/.local/state/ccam/ccusage`), bukan di folder yang sedang Anda pakai, sehingga proyek di folder itu tidak dapat menentukan apa yang dijalankan keduanya lewat `node_modules/.bin` miliknya, `bin` di `package.json` miliknya, atau registri di `.npmrc` miliknya; `ccam --dry-run ccusage` menampilkan `cd` ke folder ccam tersebut. Jika tidak ada satu pun yang tersedia, ccam meminta Anda menginstal Node.js, Bun, atau ccusage. Di Linux dan macOS, ccusage menggantikan ccam; di Windows, ccam menunggunya dan keluar dengan kode keluar ccusage. Di WSL, ccam melewati `ccusage`, `bun`, atau `npx` yang berada di drive Windows, karena program itu akan berjalan sebagai program Windows tanpa folder akun.

- Laporan mencakup akun-akun yang foldernya berisi percakapan (folder `projects`). Percakapan yang dibagikan antarakun melalui sebuah grup hanya dihitung sekali.
- Dengan `-a`, akun yang berbagi riwayatnya melalui grup membawa serta percakapan seluruh grup, dan ccam menyebutkan milik akun mana saja percakapan itu: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- `blocks` menampilkan jendela 5 jam, dan setiap akun punya jendelanya sendiri, jadi sebutkan satu akun saja untuk laporan itu, seperti contoh di atas.
- ccusage juga melaporkan alat coding lain yang ditemukannya, seperti Codex, kecuali Anda menyebut salah satunya: `ccam ccusage claude daily` hanya mencakup Claude Code.
- Saat dijalankan oleh ccam, ccusage hanya membaca pengaturannya (`ccusage.json`) dari `./.ccusage` dan folder akun. `--config` menunjuk file lain: `ccam ccusage daily --config ~/.config/claude/ccusage.json`. Karena `bun x` dan `npx` dimulai di folder milik ccam sendiri, ccam meneruskan `./.ccusage/ccusage.json` kepada keduanya dengan `--config` jika file itu ada dan Anda tidak menunjuk file apa pun, dan memberikan `--config` yang relatif sebagai path lengkap. Dalam hal ini, ccusage membaca path relatif lain, seperti `--pi-path`, dari folder ccam, jadi tuliskan path tersebut secara lengkap.
- ccusage akan membaca folder yang path-nya mengandung koma sebagai dua folder, jadi ccam menolaknya dan menyebutkan akunnya; pakai `-a` agar akun itu tidak disertakan.
- `ccusage` adalah perintah ccam, jadi akun bernama ccusage dijalankan dengan `ccam run ccusage`; `ccam doctor` memberitahukan hal ini.

## Cara kerjanya

Angka yang ditampilkan `ccam usage` berasal dari dua tempat: endpoint penggunaan Anthropic yang tidak didokumentasikan, `GET https://api.anthropic.com/api/oauth/usage`, yang dipanggil dengan login masing-masing akun, dan nilai yang diberikan Claude Code ke perintah status line-nya. Semuanya ada di `internal/usage`.

1. ccam memuat `config.toml` dan membuat satu `usage.Account` per akun (`cli.usageCmd`): foldernya `~/.local/share/ccam/profiles/<name>/`, penanda sesinya (bersama catatan milik Claude Code sendiri di folder akun, penanda itu menunjukkan apakah akun sedang berjalan), dan apakah akun itu diimpor di tempatnya.
2. ccam mencari semua akun sekaligus, satu goroutine per akun (`Fetcher.FetchAll`). Semuanya berbagi satu `http.Client`, sehingga koneksi ke endpoint dipakai ulang. Menyambung boleh memakan waktu 2 detik, dan setiap akun mendapat total 3 detik. Klien tidak mengikuti pengalihan, karena pengalihan akan mengirim ulang token ke tempat lain, dan hanya membaca paling banyak 1 MB jawaban.
3. Begitu setiap akun selesai, barisnya digambar ulang. Akun yang gagal atau tidak menjawab tepat waktu beralih ke angka yang tersimpan, sehingga satu akun yang lambat tidak pernah menahan akun lain.

Dalam diagram, setiap akun dicari di goroutine-nya sendiri, dan akun yang melewati batas waktu menampilkan angka tersimpannya sebagai gantinya:

<!-- motion:usage-fetch -->

`personal` dan `work` menjawab tepat waktu dan menampilkan angka langsung. `side` membutuhkan lebih dari 3 detik, sehingga ccam menampilkan angka tersimpannya sebagai gantinya (`fallback()`). ccam lalu menggambar satu blok per akun (`usage.Table()`) dan memilih akun yang disebut di baris terakhir (`usage.Best()`).

Untuk satu akun, ccam melakukan pemeriksaan ini secara berurutan (`Fetcher.Fetch`):

<!-- motion:usage-cache -->

1. Dengan `--source cache`, langsung menuju `fallback()`.
2. Kecuali `--no-cache` diatur, hasil langsung di `~/.local/state/ccam/usage/<name>.json` yang diukur kurang dari 60 detik yang lalu (`CacheTTL`) dipakai, dinilai menurut saat ini: jendela yang sudah direset sejak itu menampilkan 0%, dan file yang bertanggal di depan jam tidak dihitung sebagai baru. Menjalankan `ccam usage` dua kali berturut-turut hanya memakan satu permintaan per akun.
3. ccam membaca `.credentials.json` (atau Rantai Kunci (Keychain) macOS) tanpa mengubahnya dan hanya mendekode kunci `claudeAiOauth` (`readCredentials()`). Tidak adanya login langganan berarti `not logged in`, atau `no limits (API)` jika `.claude.json` berisi kunci API. Login yang refresh token-nya sudah kedaluwarsa, atau yang dihapus Claude Code, menjadi `` login expired → `ccam login <name>` `` dan tidak dicari.
4. Token akses yang kedaluwarsa tidak dipakai. Tanpa `--refresh`, akun menampilkan angka tersimpannya dengan ``access token expired; `ccam <name>` renews it when it starts``.
5. Jika tidak, ccam memanggil endpoint. 401 atau 403 dianggap sebagai login yang kedaluwarsa (`login expired: authentication failed (401)`), status lain atau isi yang tidak dapat dibaca dianggap galat, dan keduanya beralih ke angka yang tersimpan. Pencarian yang diakhiri dengan Ctrl+C menampilkan `cancelled`. 429 berarti endpoint sudah menerima terlalu banyak permintaan untuk akun ini: ccam menulis `<name>.pause` berisi waktu untuk bertanya lagi (`Retry-After`, atau 5 menit, paling lama satu jam), menampilkan angka yang tersimpan dengan `rate limited by the usage endpoint; asking again after 15:10`, dan tidak memanggil endpoint untuk akun itu sampai saat tersebut, bahkan dengan `--no-cache`. Mencoba lagi seketika hanya akan memperpanjang pembatasan itu.
6. Jawaban yang baik ditulis ke `<name>.json` bersama waktu pengukurannya.

Untuk angka tersimpan, ccam mengambil pencarian tersimpan yang terbaru (`<name>.json`, atau `<name>.claude.json` dari pencarian Claude Code) dan, jika `<name>.statusline.json` ditulis setelahnya, menimpakan jendela 5 jam dan mingguan dari status line di atasnya (`fallback()`); jendela per model, yang tidak dimiliki status line, tetap seperti yang disimpan pencarian itu. Sebelum menampilkan angka tersimpan, ccam memeriksa waktu reset setiap jendela: jendela yang sudah direset sejak angka disimpan menampilkan 0%, karena persentase lama sudah tidak berlaku. Itulah sebabnya `personal` di atas menampilkan `0%` tanpa waktu reset.

Endpoint memiliki dua format respons. Yang lebih baru memiliki array `limits` dengan entri `session` (5 jam), `weekly_all`, dan `weekly_scoped` (per model). Yang lebih lama memiliki objek `five_hour`, `seven_day`, `seven_day_opus`, dan `seven_day_sonnet`. ccam membaca keduanya (`ParseLive`): array `limits` yang menang, dan objek lama mengisi apa yang tidak dicakupnya. Isi tanpa batas sama sekali dilaporkan sebagai format yang berubah, yang akan diberitahukan oleh `ccam doctor`.

### Catatan batas

Saat Claude Code menghentikan sebuah giliran di batas penggunaan, hook dari [plugin ccam](/docs/install#claude-code-plugin) membaca permintaan yang ditolak dari percakapan: jendela mana yang penuh (5 jam, mingguan, Opus, atau Sonnet) dan kapan jendela itu direset. Hook menulisnya ke `~/.local/state/ccam/usage/<name>.limit.json`. ccam menerapkan catatan itu di atas apa pun yang diketahuinya, dari setiap sumber (`Fetch`): setiap jendela di dalamnya dihitung 100% sampai di-reset, seberapa pun lama atau rendahnya angka yang tersimpan, sehingga `ccam next --auto`, `ccam pick`, `--fallback`, dan [`auto_switch`](/features/next#auto-switch) tidak memilih akun itu lagi sebelum saat itu.

Jika tidak ada hal lain yang diketahui, catatan itu saja yang ditampilkan untuk akun tersebut (`source: the usage limit Claude Code reported` dengan `-v`). Sesi tanpa pengawasan (proses `-p` dan sesi latar belakang) tidak menulis catatan, begitu pula penghentian yang ditetapkan oleh organisasi (misalnya saat organisasi kehabisan penggunaan atau batas sebuah grup diatur ke $0), yang bisa dicabut adminnya kapan saja. `ccam login` ke akun Claude lain membuang catatan itu; menghapus file itu juga menghapusnya.

### Memperbarui token

Refresh token berotasi: memperbaruinya memberi Anda refresh token baru dan membatalkan yang lama. Jika ccam memperbarui login saat Claude Code sedang berjalan sebagai akun yang sama, sesi itu akan keluar. Karena itu `ccam usage` membaca token dan tidak pernah menulisnya, kecuali Anda memberikan `--refresh`. Dalam hal itu ccam (`Fetcher.refresh`):

1. menolak akun yang didaftarkan di tempatnya dengan `ccam import`, karena `claude` biasa mungkin memakai login itu tanpa sepengetahuan ccam;
2. menolak login yang tersimpan di Rantai Kunci (Keychain) macOS;
3. mengambil file lock, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, sehingga dua proses ccam tidak pernah sama-sama memperbarui token yang sama; selama lock dipegang, akun menampilkan `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)`. `ccam rename`, `remove`, dan `logout` juga mengambilnya, dan menolak selama pembaruan sedang berlangsung;
4. di bawah lock, menolak jika ada sesi akun yang sedang berjalan, bagaimanapun sesi itu dijalankan (melalui ccam, `claude` biasa, IDE, atau sesi latar belakang dan daemon di belakangnya, yang memperbarui login itu sendiri), dan melewati permintaan jika proses lain sudah memperbarui token;
5. tidak mengirim apa pun selama Claude Code memegang lock pembaruannya sendiri di folder akun (`.oauth_refresh.lock`, `<folder>.lock`): Claude Code sedang memperbarui login itu saat itu juga;
6. mengirim refresh token ke endpoint token dan menunggu jawabannya hingga 30 detik, waktu yang diberikan Claude Code untuk pembaruannya sendiri. Setelah terkirim, permintaan tidak pernah diputus, bahkan oleh Ctrl+C: server mungkin sudah merotasi refresh token, dan hanya jawabannya yang membawa login baru. Ctrl+C atau hangup mengakhiri pencarian lainnya, dan ccam keluar (130 setelah Ctrl+C) begitu jawabannya tersimpan;
7. menyimpan di bawah lock tulis milik Claude Code (`.storage-write.lock`) dan hanya mengubah kolom token di `.credentials.json` sebagaimana isinya saat itu: apa yang ditulis Claude Code selama itu (info masuk MCP) tetap ada, login yang diperbarui Claude Code selama itu dipertahankan, dan logout tidak pernah dibatalkan. Jika jawabannya menyebutkan masa berlaku refresh token baru, nilai itu disimpan sebagai `refreshTokenExpiresAt`, seperti yang dilakukan Claude Code.

Pembaruan yang tidak berhasil memberi tahu alasannya, di baris akun atau di atas angka tersimpannya: misalnya `token not renewed: a session is running`, atau `login may be lost: …` jika permintaan sudah terkirim tetapi jawabannya tidak bisa dipakai atau disimpan, dengan saran untuk menjalankan `ccam login <name>` jika Claude Code kemudian meminta Anda login.

### Pengumpul status line

Claude Code menjalankan perintah status line dan mengirimkan dokumen JSON kepadanya melalui pipe. ccam membaca nama model dan `rate_limits` darinya:

<!-- motion:usage-statusline -->

1. ccam menemukan akun dari `CCAM_PROFILE`, yang diatur ccam saat menjalankan Claude Code, atau jika tidak ada, dengan mencocokkan `CLAUDE_CONFIG_DIR` dengan folder akun (`runStatusline`).
2. ccam membaca `rate_limits.five_hour` dan `rate_limits.seven_day` (`ParseStatusLine`). JSON tidak berisi batas laju untuk akun kunci API, atau sebelum respons pertama sesi; dalam hal itu tidak ada yang disimpan.
3. ccam menulis `~/.local/state/ccam/usage/<name>.statusline.json` (`SaveStatusLine`). `ccam usage` memakainya jika endpoint tidak menjawab.
4. ccam menampilkan `[name] · model · 5h N% · 7d N%`. Dengan `--wrap`, ccam menjalankan perintah Anda melalui `sh -c` dengan input yang sama dan menampilkan keluarannya seperti cara Claude Code menampilkan status line: setiap baris yang tidak kosong, dengan `[name]` di depan baris pertama. Perintah yang gagal hanya menampilkan nama akun. Di Windows, ccam menjalankan perintah itu dengan cara yang sama seperti Claude Code: dengan Git Bash jika terpasang, atau dengan PowerShell jika tidak.
5. Saat penggunaan 5 jam atau mingguan mencapai 90% atau lebih dan ada akun lain, baris itu diakhiri dengan `→ ccam next --auto --continue`, misalnya `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Perintah yang dibungkus tidak memiliki batas waktu sendiri: Claude Code memberi perintah status line waktu 10 menit dan menghentikan proses saat pembaruan yang lebih baru dimulai, sehingga perintah yang lambat tampil seperti tanpa ccam. Saat Claude Code mengakhiri ccam, perintah itu juga berakhir (di Linux dan macOS beserta setiap proses yang dimulainya).

`ccam statusline install` menulis `statusLine` ke `settings.json` setiap akun dan mempertahankan kunci lain objek itu, seperti `padding`, `refreshInterval`, dan `hideVimModeIndicator`. Jika sebuah akun sudah memiliki perintah status line, ccam mempertahankannya sebagai `--wrap '<old command>'`; di Windows sebagai `--wrap-base64 <base64>`, yang tidak memerlukan tanda kutip di shell mana pun. Jika status line sudah menjalankan ccam (dari path lain, atau dengan tanda kutip seperti yang ditulis versi lama), ccam menggantinya alih-alih membungkusnya, dan mempertahankan perintah yang dibungkusnya.

Di Windows, path ke ccam ditulis dengan garis miring biasa dan `~` untuk folder home (`~/.local/bin/ccam.exe statusline`), yang dipahami oleh Git Bash maupun PowerShell. Path yang memerlukan tanda kutip ditulis lengkap untuk shell yang dipakai Claude Code: `'…'` untuk Git Bash, `& '…'` untuk PowerShell. `settings.json` yang dibagikan melalui grup ditulis sekali, melalui tautannya, sehingga setiap akun dalam grup mendapatkannya. Menjalankannya lagi menampilkan `(already set)` untuk file yang sudah memakai ccam.

## Keamanan dan kasus khusus

- `ccam usage` tidak pernah mengubah file token kecuali Anda memberikan `--refresh`. Jika token akses sudah kedaluwarsa, jalankan akun tersebut sekali dengan `ccam <name>` dan Claude Code akan memperbaruinya. Login yang sudah kedaluwarsa sepenuhnya memerlukan `ccam login <name>`.
- `--refresh` tidak pernah menyentuh akun yang sedang berjalan (bahkan jika Claude Code dijalankan tanpa ccam), akun yang diimpor, atau login Rantai Kunci (Keychain). Pembaruan yang ditolak menampilkan angka tersimpan beserta alasannya (`token not renewed: …`).
- Endpoint penggunaan tidak didokumentasikan secara publik dan dapat berubah. Jika berubah, Anda tetap melihat angka yang tersimpan, dan `ccam doctor` melaporkan perubahan format.
- Jika tidak ada angka tersimpan dan tidak ada jawaban, baris akun menampilkan `lookup failed:` beserta alasannya, misalnya `no saved usage` dengan `--source cache`.
- File cache hanya berisi persentase dan waktu reset, tidak pernah token.

## Terkait

- [Batas penggunaan dan perpindahan akun](/docs/usage-limits)
- [Beralih ke akun berikutnya](/features/next)
- [Saat Claude Code mencapai batas penggunaan](/guides/usage-limit-workflow)
- [Referensi perintah](/docs/commands)
