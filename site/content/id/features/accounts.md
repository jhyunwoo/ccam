---
title: Mengelola akun
description: Tambahkan, login, tampilkan, ganti nama, hapus, dan impor akun Claude Code dengan ccam add, ccam login, ccam list, ccam rename, ccam remove, dan ccam import.
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

Setiap akun ccam adalah folder konfigurasi Claude Code tersendiri dengan login-nya sendiri. Gunakan perintah-perintah ini untuk menambahkan akun, memasukkan login-nya, menampilkan, mengganti nama, dan menghapusnya. Perintah-perintah ini juga mendaftarkan folder yang sudah Anda miliki.

## Penggunaan

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

Di mana pun sebuah perintah menerima `<name|number>`, Anda dapat memberikan nama akun atau nomor slotnya (kolom `#` pada `ccam list`).

`ccam use` hanya mengubah akun mana yang menjadi default (`default` di `config.toml`).

`ccam add` di terminal tanpa opsi akan membuka wizard. Wizard melalui langkah-langkah berikut:

1. nama akun, kecuali Anda sudah memberikannya;
2. cara akun itu login: login Claude sekarang, kunci API, atau login Claude nanti;
3. grup yang akan diikuti, dengan `default` sudah dicentang (jika belum ada grup, wizard menampilkan apa yang akan dibagikan grup `default`);
4. `Create a shortcut command for work? [y/N]`;
5. ringkasan untuk dikonfirmasi.

`Esc` kembali satu langkah. Salah satu dari `--group`, `--no-group`, `--shortcut`, atau `--no-login`, atau menjalankan tanpa terminal (misalnya, dari pipe), akan melewati wizard. Dalam hal itu nama wajib diberikan, dan tidak ada pintasan yang dibuat kecuali `--shortcut` menyebutkannya.

Tanpa `--group`, akun baru bergabung ke grup `default`. Akun pertama membuat grup itu, yang membagikan server MCP, pengaturan, `CLAUDE.md`, dan ekstensi; lihat [berbagi dengan grup](/features/share). Di Windows, file hanya bisa ditautkan dalam Mode Pengembang (Developer Mode). Tanpa mode itu, akun pertama membuat `default` hanya dengan `mcp` dan `extensions` serta memberitahukannya (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), dan ccam menolak `--group` dengan grup yang membagikan pengaturan, memori, atau riwayat.

`ccam list` menampilkan slot, nama, pintasan, email, grup, dan waktu terakhir digunakan, dengan `*` di depan akun default. Emailnya adalah email yang saat ini dipakai folder itu untuk login (dari `.claude.json` miliknya); jika berbeda dari yang dicatat ccam, email yang dicatat menyusul sebagai `(was …)`:

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

Dengan `--json`, setiap akun adalah objek dengan `slot`, `shortcut`, `name`, `email` (milik login saat ini, seperti di tabel), `default`, `dir`, `last_used` (dihilangkan jika akun belum pernah dijalankan), `groups`, dan `share`. `share` memetakan setiap item (`mcp`, `settings`, `memory`, `extensions`, `history`) ke grup asalnya, atau ke `"local"`.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `ccam add --no-login` | Membuat akun tanpa login; login belakangan dengan `ccam login <name>` | nonaktif |
| `ccam add --group <group>` | Bergabung ke grup ini alih-alih `default`. Bisa diulang; `default` dibuat jika belum ada | `default` |
| `ccam add --no-group` | Tidak bergabung ke grup mana pun: akun tidak membagikan apa pun | nonaktif |
| `ccam add --shortcut <command>` | Juga membuat perintah pintasan ini di `~/.local/bin` | tidak ada |
| `ccam add --api --key-from <source>` | Akun yang memakai kunci API alih-alih login; lihat [akun kunci API](/features/api) | nonaktif |
| `ccam list --json` | Menampilkan akun sebagai JSON | nonaktif |
| `ccam remove -y`, `--yes` | Menghapus tanpa bertanya | nonaktif |
| `ccam logout -y`, `--yes` | Logout tanpa bertanya | nonaktif |
| `ccam import --name <name>` | Nama akun untuk folder yang Anda berikan; ditolak tanpa folder (`--name needs a folder: ccam import <folder> --name ops`) | diambil dari nama folder |
| `ccam import -y`, `--yes` | Mengimpor setiap folder yang ditemukan tanpa bertanya | nonaktif |
| `--dry-run` | Didukung oleh `ccam list` dan `ccam completion`, yang tidak mengubah apa pun. Ditolak oleh `add`, `login`, `logout`, `use`, `rename`, `remove`, dan `import`, dengan `` --dry-run is not supported by `ccam add`; nothing was changed `` (menyebutkan perintahnya) | nonaktif |

`ccam login` (selain `--expired`), `ccam use`, dan `ccam rename` tidak memiliki opsi sendiri. `ccam completion` menerima `bash`, `zsh`, atau `powershell`; selain itu gagal dengan `only bash, zsh and powershell are supported`.

Nama akun terdiri dari 1 sampai 32 karakter `a-z`, `0-9`, `-`, dan `_`, diawali huruf atau angka. Nama tidak boleh seluruhnya angka, karena `ccam 2` berarti slot 2. Nama juga tidak boleh berupa perintah ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`), sehingga `ccam <name>` tidak pernah ambigu. Daftarnya adalah `config.Reserved` di `internal/config`.

## Cara kerjanya

Menambahkan akun akan membuat foldernya, bergabung ke grup-grupnya, menyimpan konfigurasi, melakukan login, dan membuat pintasan, dengan urutan itu. Saat Anda menjalankan `ccam add work --shortcut ccwork`, ccam mengambil langkah-langkah ini (`addProfile` di `internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. ccam memeriksa nama (`config.ValidateName`) dan pintasan (`shim.Check`) sebelum membuat apa pun.
2. ccam membuat `~/.local/share/ccam/profiles/work/` dengan mode `0700` dan memberi akun nomor slot kosong terkecil, sehingga nomor akun yang dihapus dipakai ulang (`profile.Add`, `Config.NextSlot`). Akun pertama menjadi akun default.
3. Akun bergabung ke grup-grupnya (`default` kecuali `--group` atau `--no-group` menentukan lain; akun pertama membuat `default`), dan ccam menautkan folder baru yang masih kosong ke grup-grup itu, misalnya `settings.json` ke `groups/default/settings.json` (`applyShare`, dengan `share.Build` lalu `Apply`). Jika langkah ini gagal, folder baru dihapus, karena folder yang tertinggal akan menggagalkan setiap percobaan ulang dengan `folder already exists`.
4. Akun disimpan ke `~/.config/ccam/config.toml` melalui `fsutil.AtomicWrite`: file sementara, `fsync`, lalu `rename`, sehingga crash tidak pernah meninggalkan konfigurasi yang setengah jadi.
5. ccam menjalankan `claude auth login` dengan `CLAUDE_CONFIG_DIR` diatur ke folder akun, dan `--email` dengan email terakhir akun itu agar halaman login sudah terisi (`runLogin`). Jika login gagal, akun tetap ada, dan ccam menampilkan galatnya serta meminta Anda menjalankan `ccam login work` nanti.
6. ccam membaca `oauthAccount` dari `profiles/work/.claude.json` (`profile.ReadIdentity`): email, ID akun dan organisasi, nama organisasi, dan paketnya (`claude_max`, `claude_team` …). File itu bisa berukuran beberapa megabyte, sehingga ccam memindai satu kunci itu saja alih-alih mem-parsing seluruh file (`TopLevelValue`). ccam menyimpan detail ini di `[profile.login]` dalam `config.toml`.
7. ccam menautkan `~/.local/bin/ccwork` ke biner ccam, yang mencari `argv[0]` di `config.toml` saat dijalankan, dan mencatat tautan itu di `~/.local/state/ccam/shortcuts.json` (`shim.Install`). Jika tautan tidak bisa dibuat, akun tetap ditambahkan, ccam memberi tahu alasannya (`Could not create shortcuts: …`), dan ringkasan menampilkan `shortcut: ccwork (not created; see above)`.

Login bisa memakan waktu beberapa menit, dan proses ccam lain mungkin menyimpan konfigurasi pada saat itu. Karena itu, perintah yang menunggu memuat ulang `config.toml`, menerapkan satu perubahan itu (di sini, email), lalu menyimpan (`env.update`). `ccam remove` memuat ulang setelah pertanyaannya karena alasan yang sama.

`ccam login` menjalankan langkah 5 dan 6 untuk akun yang sudah ada. Claude Code menganggap email itu hanya sebagai petunjuk untuk halaman login, sehingga Anda bisa saja masuk ke akun Claude yang berbeda dari sebelumnya; ccam membandingkan login baru dengan `[profile.login]` dan memberitahukannya (`ccam: this login is other@company.com, but work was dev@company.com before.`). ccam mempertahankan login baru dan memberi tahu cara membatalkannya: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

Login ke akun Claude lain juga membuang [catatan batas](/features/usage#the-limit-record) akun itu, karena batas tersebut milik akun yang lain. Jika akun ccam lain sudah login ke akun Claude yang sama, ccam memperingatkan bahwa keduanya berbagi batas penggunaan.

`ccam login --expired` melakukan login, satu per satu, untuk setiap akun yang belum pernah login atau yang login-nya sudah kedaluwarsa: refresh token-nya kedaluwarsa, atau Claude Code menghapus login itu setelah server menolaknya. Sebelum browser pertama terbuka, perintah ini memeriksa bahwa `config.toml` bisa disimpan.

`ccam logout` menjalankan `claude auth logout` di akun tersebut dan menghapus email yang tersimpan; catatan login tetap ada, sehingga login berikutnya diperiksa terhadapnya. Claude Code mencabut login itu dan menghapus apa yang disimpannya untuk akun tersebut, termasuk info masuk server MCP. Karena itu ccam memperingatkan Anda dan bertanya lebih dulu (`--yes` melewati pertanyaan itu) jika:

- login memegang info masuk MCP;
- akun sedang berjalan (sesinya juga ikut ter-logout);
- folder akun lain memegang login yang sama (folder yang disalin), yang juga ikut dicabut oleh logout.

Contohnya:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

Untuk login sebagai akun Claude lain, `ccam login <name>` saja sudah cukup, dan info masuk MCP tetap dipertahankan.

`ccam list` menambahkan kolom `Organization` jika akun-akun berasal dari lebih dari satu organisasi (paket pribadi ditampilkan sebagai `personal`), dan `ccam list --json` memiliki `plan` dan `org` untuk setiap akun.

### Login di mesin jarak jauh

Melalui SSH, bagian browser dari `claude auth login` bisa gagal. Claude Code menunggu login di server callback pada mesin jarak jauh (`http://localhost:<port>/callback`), tetapi browser yang menyelesaikan login berjalan di komputer Anda, tempat alamat itu tidak mengarah ke mana pun dan browser menampilkan galat seperti `ERR_CONNECTION_REFUSED`. Alamat itu tetap membawa kode sekali pakai, sehingga ccam dapat menyampaikannya:

<!-- motion:accounts-relay -->

ccam tidak menyerahkan terminal Anda ke Claude Code (`runLogin`). ccam menghubungkan input standar Claude Code ke sebuah pipe dan menyalin setiap baris yang Anda ketik ke dalamnya (`relayLogin`), sehingga kode yang ditempel di prompt `Paste code here` milik Claude Code tetap sampai. Jika sebuah baris berupa alamat callback, ccam hanya menerimanya jika host-nya `localhost`, `127.0.0.1`, atau `[::1]`, port disebutkan, path-nya `/callback`, serta `code` dan `state` keduanya ada; `http://` boleh dihilangkan (`callbackURL`). ccam lalu meminta alamat itu di server dengan batas waktu 10 detik dan tanpa mengikuti pengalihan, seperti yang akan dilakukan browser (`deliverCallback`). Permintaan itu sampai ke server callback milik `claude auth login` yang sedang menunggu (garis putus-putus di diagram). Claude Code menerima kode dan menyimpan login, di Linux dalam `.credentials.json` di folder akun. Baris yang ditempel tidak akan pernah membuat ccam menghubungi host lain.

ccam menampilkan ini sebelum login dimulai:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

Jika berhasil, ccam menampilkan `→ Passed the sign-in to Claude Code on this machine.`

Selama login menunggu input di terminal, ccam mematikan echo dan mode baris terminal serta menampilkan `*` untuk setiap karakter yang Anda ketik, sehingga kode atau alamat callback yang ditempel tidak muncul di layar. ccam menangani penyuntingan baris sendiri (`readLine`): `Enter` mengirim baris, `Backspace` menghapus satu karakter, dan `Ctrl+D` di baris kosong mengakhiri input. Terminal dipulihkan saat login berakhir, dan juga saat `Ctrl+C`, `SIGTERM`, atau `SIGHUP`, yang tetap mengakhiri ccam. Jika input standar dan error standar tidak terhubung ke terminal, baris disalin apa adanya saat masuk.

### Mengganti nama, menghapus, dan mengimpor

`ccam rename` dan `ccam remove` bekerja pada folder akun dan pada apa yang disimpan ccam atas nama akun itu di tempat lain. ccam mendaftar path tersebut di `nameState`:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- cache penggunaan `~/.local/state/ccam/usage/<name>.json` dan `<name>.statusline.json`;
- catatan batas `<name>.limit.json`;
- catatan sesi di `~/.local/state/ccam/handoff/<name>/`.

`ccam remove work` menanganinya dengan urutan ini:

<!-- motion:accounts-remove -->

1. `ccam remove` pertama-tama menolak jika akun sedang berjalan, lalu bertanya, dan memeriksa lagi setelah Anda menjawab (ini mendeteksi sesi yang dimulai selagi pertanyaan terbuka).
2. Perintah ini memeriksa bahwa `config.toml` bisa disimpan sebelum mengubah apa pun.
3. Jika akun membagikan ekstensi, ccam menulis ulang `installed_plugins.json` dan `known_marketplaces.json` milik grup (`RewritePluginPaths`), karena Claude Code menyimpan path plugin absolut yang mungkin melewati folder akun ini.
4. Akun terkelola (akun yang foldernya dibuat oleh ccam, bukan yang didaftarkan dengan `ccam import`) kemudian di-logout: `claude auth logout` dijalankan di foldernya, sehingga tidak ada login yang masih berfungsi menunggu di tempat sampah (atau, di macOS, di item Rantai Kunci yang dinamai sesuai folder itu, yang akan diwarisi oleh akun berikutnya dengan nama yang sama). Langkah itu dilewati untuk akun kunci API dan jika folder akun lain memegang login yang sama, yang juga akan ikut dicabut oleh logout. Logout yang gagal dilaporkan, dan penghapusan tetap berlanjut.
5. ccam memindahkan folder ke `~/.local/share/ccam/trash/<name>-<time>/` (`profile.Trash`), akun dikeluarkan dari `config.toml`, path `nameState` dihapus agar akun berikutnya dengan nama yang sama tidak mewarisinya, dan tautan pintasan dihapus: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

Jika `config.toml` ternyata tidak bisa disimpan, folder dikembalikan dari tempat sampah. Untuk memulihkan akun yang dihapus, pindahkan kembali folder itu dari tempat sampah dan login lagi.

`ccam rename` mengganti nama folder dan menulis ulang path plugin di `plugins/` miliknya dari folder lama ke folder baru. Perintah ini memindahkan path `nameState` ke nama baru dan mempertahankan akun sebagai default jika sebelumnya default. Perintah ini juga memperbarui catatan sesi, sehingga `ccam next --continue` tetap menemukannya. Perintah pintasan tetap memakai namanya. Jika penulisan ulang plugin gagal, atau `config.toml` tidak bisa disimpan, nama folder dikembalikan.

Di macOS, dan di Windows jika Claude Code menyimpan login di Credential Manager, folder itu sendiri tetap memakai nama lamanya, karena login dan info masuk MCP disimpan di bawah path folder itu dan folder yang dipindahkan akan ter-logout. Konfigurasi kemudian mencatat folder itu sebagai `folder` di samping nama baru, dan ccam menampilkan `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

`ccam import` mendaftarkan folder di tempatnya, sebagai `dir` di `config.toml`. Perintah ini tidak pernah memindahkan atau menyalinnya: refresh token yang tersalin akan membuat salinan-salinan itu saling mengeluarkan login satu sama lain.

Tanpa argumen folder, perintah ini mencari folder `~/.claude` dan `~/.claude-*` yang belum terdaftar dan menanyakan masing-masing. Nama diambil dari folder (`~/.claude` menjadi `main`, `~/.claude-work` menjadi `work`, `~/.claude-2` menjadi `acct2`). Jika ccam tidak dapat menurunkan nama, atau nama itu sudah dipakai, ccam menampilkan folder itu beserta perintah `--name` yang perlu dijalankan. Folder yang ditemukan tetapi gagal dilaporkan sebagai `Skipped <path>: <reason>` dan yang lain tetap diimpor.

ccam mengajukan semua pertanyaannya lebih dulu, lalu memuat ulang `config.toml` dan mendaftarkan folder-folder itu, sehingga perubahan yang disimpan perintah ccam lain pada saat itu tetap ada. Baris `✓` muncul setelah penyimpanan.

Tanpa `CLAUDE_CONFIG_DIR`, Claude Code menyimpan `.claude.json` di folder home Anda di samping `~/.claude`, sehingga saat Anda mengimpor `~/.claude` dan folder itu tidak memiliki `.claude.json` sendiri, ccam membuat tautan `~/.claude/.claude.json` ke `~/.claude.json`. `claude` biasa dan `ccam main` kemudian melihat akun yang sama.

Jika `~/.claude` memiliki `.claude.json` sendiri di samping `~/.claude.json`, ccam membiarkan keduanya dan memberitahukannya. `ccam main` akan memakai yang di dalam dan `claude` biasa yang di luar, sehingga detail login dan server MCP keduanya terpisah. Untuk mempertahankan satu saja, ganti file yang di dalam dengan tautan ([doctor](/features/doctor) juga melaporkannya).

Jika folder itu memegang login yang sama dengan akun yang sudah terdaftar (folder yang disalin), ccam memberitahukannya: mana pun yang memperbaruinya lebih dulu akan membuat yang lain ter-logout, jadi jalankan `ccam login <name>` sebelum memakainya.

Akun yang diimpor tidak berada di grup mana pun sampai Anda menambahkannya ke salah satu grup (`ccam group join <group> <name>`), dan tidak ada pintasan yang dibuat.

`ccam completion bash|zsh` menampilkan skrip pelengkapan yang juga melengkapi nama akun.

## Keamanan dan kasus khusus

- `ccam rename` dan `ccam remove` menolak saat akun sedang berjalan (`refuseRunning`), dengan ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``, karena sesi terus menulis ke folder tersebut. Pemeriksaan ini mencakup Claude Code bagaimanapun ia dijalankan: melalui ccam, `claude` biasa, IDE, sesi latar belakang, dan daemon-nya.
- `ccam rename`, `ccam remove`, dan `ccam logout` menolak selama proses ccam lain sedang memperbarui atau mengubah login itu: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- Selama perubahan grup sedang berlangsung di proses ccam lain, `ccam rename` dan `ccam remove` menolak dengan `another ccam process is changing the groups right now; try again when it is done`; setelah perubahan grup terhenti di tengah jalan, dengan ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.
- `ccam remove` pada akun yang diimpor mempertahankan folder dan login-nya, untuk `claude` biasa. Jika akun itu berada di grup, akun lebih dulu keluar dari grup-grup itu seperti yang dilakukan `ccam group leave`: setiap item bersama menjadi salinan milik grup di dalam folder, dan server MCP grup masuk ke `.claude.json` miliknya. Lalu pendaftarannya dibatalkan: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. Cadangan `.claude.json` miliknya yang dibuat ccam dipindahkan ke `<folder>/.ccam-backup/mcp-<time>`; jika tidak bisa dipindahkan, cadangan itu tetap di `~/.local/state/ccam/mcp/<name>/backups` dan pesannya memberitahukan hal itu. `~/.claude` Anda tidak dipindahkan ke tempat sampah.
- Dua akun tidak boleh memiliki folder yang sama: `config.toml` yang menyebut satu folder dua kali gagal dimuat.
- `ccam add` menolak nama yang foldernya sudah ada di `profiles/`, dan mengarahkan ke `ccam import`. Di macOS perintah ini juga berhenti jika Rantai Kunci masih menyimpan login atau info masuk MCP untuk path folder baru itu, yang ditinggalkan oleh versi ccam yang lebih lama, dan menyebutkan perintah `security delete-generic-password` yang menghapusnya. Jika `config.toml` tidak bisa disimpan, folder baru dihapus lagi.
- `ccam import` menolak folder yang sudah terdaftar, path yang bukan folder, dan nama yang sudah dipakai. ccam juga memperingatkan jika folder itu adalah salinan `~/.claude` (login yang sama), meskipun `~/.claude` sendiri tidak terdaftar.
- `ccam logout` bertanya lebih dulu jika akun sedang berjalan; `ccam login` tidak memeriksanya.
- Jika penerusan tidak dapat menjangkau Claude Code, misalnya setelah login berakhir, ccam menampilkan `ccam: could not pass the sign-in to Claude Code:` beserta alasannya. Galat HTTP biasanya berarti kode sudah kedaluwarsa; mulai login lagi.

## Terkait

- [Menjalankan Claude Code](/features/launch)
- [Perintah pintasan](/features/shortcuts)
- [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing)
- [Gunakan beberapa akun Claude Code di satu komputer](/guides/multiple-claude-code-accounts)
