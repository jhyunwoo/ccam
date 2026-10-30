---
title: Dasbor
description: ccam di terminal membuka dasbor layar penuh berisi akun, grup, aturan folder, ekstensi, dan pemeriksaan; setiap perubahan menampilkan rencananya dulu.
updated: 2026-10-30
order: 15
translation: machine
sourceHash: "d65f1b5bb0c8"
---

`ccam` tanpa argumen membuka dasbor layar penuh jika input dan output standar adalah terminal; jika bukan, ccam mencetak `ccam list`. Dasbor menampilkan setiap akun beserta penggunaan dan grupnya, dan melakukan apa yang dilakukan perintah-perintahnya: menjalankan, menambahkan, login, mengubah grup, menetapkan aturan folder, menampilkan sesi latar belakang, memeriksa, dan memperbaiki. Menjalankan akun, atau menyambung (attach) ke sesi latar belakang, akan menutupnya; menambahkan akun atau login menyerahkan terminal sebentar, lalu kembali ke dasbor.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · r to refresh · a to add · ? for shortcuts
```

Baris di bawah nama-nama tab menjelaskan dalam satu kalimat kegunaan tab yang sedang dibuka. Di terminal yang lebih lebar, tab Accounts juga menampilkan email setiap akun. Di setiap tab, ujung baris nama-nama tab menunjukkan apakah [`auto_switch`](/features/next#auto-switch) aktif (Linux dan macOS; di Windows fitur ini diabaikan). Saat fitur ini aktif, tab Health mencantumkan apa yang masih dibutuhkannya: Claude Code 2.1.214 atau yang lebih baru, dan plugin ccam yang diaktifkan di setiap akun yang login dengan langganan Claude.

Setiap kali dibuka, dasbor bertanya ke GitHub di latar belakang apakah ada rilis ccam yang lebih baru; bingkai pertama tidak pernah menunggunya, dan saat Anda sedang offline, tidak ada yang ditampilkan. Jika ada, sebuah baris di bawah judul merekomendasikannya:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` bertanya dulu, lalu memperbarui seperti [`ccam update`](/features/update): rilis diperiksa terhadap checksum SHA-256-nya dan menggantikan ccam dalam satu langkah. Dasbor yang sedang Anda buka masih versi lama, jadi keluar lalu jalankan `ccam` lagi. Jika pembaruan tidak bisa dilakukan di sini, misalnya untuk salinan ccam yang diinstal dari paket `.deb`, dasbor menampilkan saran yang sama seperti `ccam update` dan tetap menyimpan tawarannya.

Bingkai pertama digambar hanya dari `config.toml`, jauh di bawah 100 ms; penggunaan terisi begitu setiap akun menjawab. Setiap batas punya bilahnya sendiri: jendela 5 jam, jendela mingguan untuk semua model, dan jendela mingguan khusus Fable; `-` berarti akun tidak memiliki jendela tersebut. Persentase berubah kuning pada 80% dan merah pada 100%. Jika pencarian gagal, nilai dari cache ditampilkan redup. Jendela yang dilaporkan penuh oleh Claude Code menampilkan 100% sampai di-reset. Token akses yang kedaluwarsa menampilkan `⚠ token expired` yang redup (Claude Code memperbaruinya pada saat dijalankan berikutnya), sedangkan `⚠ login expired` berarti akun itu memerlukan `ccam login`.

## Tab

| Tab | Yang ditampilkan | Tombol |
| --- | --- | --- |
| Accounts | nomor, nama, email, bilah 5 jam, mingguan, dan Fable, pemakaian terakhir, grup; `✔` akun default, `▸` akun folder ini, `•` sedang berjalan (juga Claude Code yang dijalankan tanpa ccam: `claude` biasa, IDE, sesi latar belakang, dan daemon-nya) | lihat di bawah |
| Agents | apa yang dicantumkan [`ccam agents`](/features/agents): sesi latar belakang dan sesi terbuka dari setiap akun, beserta status, folder, dan waktu mulainya (lihat di bawah) | `enter` menyambung (seperti `ccam attach`), `x` menghentikan, `r` memuat ulang |
| Groups | setiap grup, apa yang dibagikannya dan akun-akunnya, serta akun yang tidak ada di grup mana pun | `enter` sunting, `n` baru, `e` ganti nama, `x` hapus |
| Folders | aturan [`ccam folder`](/features/folder) dan aturan yang berlaku di sini, bahkan saat Anda mencapai folder itu melalui tautan simbolik (aturan baru disimpan dengan tautan yang sudah di-resolve) | `n` pakai sebuah akun di sini, `enter` ubah, `x` hapus |
| Extensions | apa yang ditampilkan [`ccam ext`](/features/ext), setiap ekstensi satu kali beserta grup asalnya, `off:` akun-akun tempat sebuah plugin tidak diaktifkan, server milik akun itu sendiri yang hanya berlaku untuk folder ini, `This folder's project`, dan masalah yang bukan milik akun mana pun | `r` muat ulang |
| Health | apa yang dilaporkan [`ccam doctor`](/features/doctor), sebagai satu kotak per akun (lihat di bawah) | `f` perbaiki kotak yang dipilih, `F` perbaiki semua (`doctor --fix`), `↑`/`↓` berpindah antarkotak |

`←`/`→` atau `Tab` berpindah tab, `?` menampilkan tombol-tombol tab yang sedang dibuka, `r` menyegarkan penggunaan tanpa melalui cache singkat dan memuat ulang tab Agents, Extensions, dan Health, dan `q` atau `Esc` keluar.

## Tombol di tab Accounts

| Tombol | Fungsinya |
| --- | --- |
| `Enter` | Menjalankan Claude Code sebagai akun yang dipilih |
| `1`–`9` | Pindah ke akun dengan nomor itu; ketik `1` `2` dengan cepat untuk 12. Hanya `Enter` yang menjalankan |
| `↑`/`↓`, `k`/`j` | Berpindah |
| `n` | Menjalankan akun yang akan dipilih [`ccam next --auto`](/features/next) |
| `o` | Menjalankan dengan opsi: `--continue`, `--resume`, argumen lain, atau uji coba (dry run) |
| `u` | Rincian penggunaan: jendela 5 jam, mingguan, dan per model beserta waktu reset-nya dan asal angka-angkanya |
| `a` | Menambahkan akun (wizard di bawah) |
| `l` | Login ulang |
| `s` | Memilih grup akun tersebut |
| `d` | Menjadikannya akun default |
| `f` | Memakainya di folder saat ini |
| `e` | Mengganti namanya |
| `c` | Menetapkan atau menghapus perintah pintasannya |
| `x` | Menghapusnya (bertanya dulu, dengan No terpilih) |

**Kursi Team dan Enterprise.** Pemilihan akun otomatis (`ccam next`, `--fallback`, `ccam pick`, tombol `n` di dasbor, dan, setelah Anda mengaktifkannya, [`auto_switch`](/features/next#auto-switch)) juga berganti di antara kursi dalam organisasi Claude Team atau Enterprise yang sama, yang dapat melanggar ketentuan Anthropic. Baca [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) dan tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya.

## Tab Health

Tab Health menjalankan pemeriksaan [`ccam doctor`](/features/doctor) saat Anda membukanya dan menggambar tampilan yang sama: baris pertama dengan tanda untuk setiap akun, lalu satu kotak per akun. Setiap baris dalam kotak adalah satu area yang diperiksa, diawali `Folder`, `Login`, dan `Share links`, dengan tanda `✔` OK, `ℹ` informasi, `⚠` peringatan, atau `✘` galat, sehingga Anda langsung melihat akun mana yang bermasalah dan di mana. Temuan yang bukan milik akun mana pun masuk ke kotak `General`.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` berpindah dari kotak ke kotak, dan kotak yang dipilih memiliki bingkai berwarna aksen. `f` hanya memperbaiki kotak yang dipilih, dengan cara yang sama seperti `ccam doctor --fix`: pada kotak sebuah akun, ia hanya memperbaiki akun itu, dan pada kotak `General`, temuan yang bukan milik akun mana pun (path plugin dan pintasan); akun lainnya hanya diperiksa lagi. Ia bertanya dulu (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), dan kotak yang tidak punya apa pun untuk diperbaiki hanya menampilkan `nothing to repair in <name>`. `F` memperbaiki semuanya sekaligus, seperti `ccam doctor --fix`, dan `r` memeriksa lagi.

## Tab Agents

Tab Agents menampilkan apa yang dicetak [`ccam agents`](/features/agents): sesi latar belakang dan sesi terbuka dari setiap akun, beserta akun pemilik masing-masing. Tab ini bertanya ke Claude Code setiap akun sekaligus (`claude agents --json`, masing-masing dengan batas waktu 15 detik) saat Anda pertama kali membuka tab, saat `r`, setelah penghentian, dan setelah konfigurasi berubah, dan menampilkan `Asking every account's Claude Code for its sessions…` sampai semua akun menjawab.

```console
  Accounts   Agents   Groups   Folders   Extensions   Health    auto_switch off
 Background agents and open sessions of every account; enter attaches to one.

    Account   ID        Name   State          Folder      Started
 ❯  work      -         -      open (busy)    ~/code/web  now
    personal  49ec0028  probe  blocked, idle  ~/code/app  12m

 ✘ lab: could not list its sessions: exit status 1 'claude agents --json' is
   not available in this environment.

 enter to attach · x to stop · r to reload · ? for shortcuts
```

`State` adalah status sesi latar belakang dan apa yang sedang dilakukan prosesnya, atau `open` untuk sesi interaktif, seperti di `ccam agents`. Jika sebuah akun tidak bisa ditanyai, baris `✘` di bawah daftar menampilkan apa yang dikatakan Claude Code, seperti `lab` di atas, yang tampilan agennya (agent view) dinonaktifkan. Sesi latar belakang yang sudah selesai tidak ditampilkan (`ccam agents --all` mencantumkannya), dan jika setiap akun sudah menjawab dan tidak ada yang memiliki sesi, tab ini menampilkan ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` Di terminal yang sempit, tab ini menghilangkan `Started`, lalu `Name`, lalu `Folder`, dan memotong path folder yang panjang dari kiri.

`Enter` pada sesi latar belakang menutup dasbor dan membuka sesi itu di terminal ini sebagai akunnya, seperti yang dilakukan `ccam attach <id>`; dasbor tidak kembali (di Claude Code, `←` membuka tampilan agennya dan `Ctrl+Z` kembali ke shell). Sesi terbuka tidak memiliki ID, jadi `Enter` dan `x` hanya memberi tahu bahwa diperlukan sesi latar belakang. `x` bertanya lebih dulu, dengan `No, keep it` terpilih: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` Jika ya, ccam menjalankan `claude stop <id>` sebagai akun itu, menampilkan ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` atau alasan mengapa tidak bisa, lalu membaca daftar itu lagi.

## Mengubah grup

Di tab Groups, `Enter` membuka penyunting: item di kiri, akun di kanan, `Tab` untuk berpindah sisi, `Space` untuk mencentang atau menghapus centang. Item yang sudah didapat salah satu akun dari grup lain diredupkan beserta alasannya, begitu pula akun yang akan mendapat sebuah item dua kali; satu item hanya boleh berasal dari satu grup. `Enter` menampilkan rencananya:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Jika sebuah file berbeda dari salinan milik grup, ccam lebih dulu menanyakan apa yang harus dilakukan (pertahankan milik grup, pakai milik saya, atau gabungkan). Lalu perubahan berjalan di latar belakang dengan kode yang sama seperti [`ccam group`](/features/share): konfigurasi disimpan lebih dulu, data setiap akun dipindahkan secara berurutan, dan file aslinya masuk ke `.ccam-backup/`; akun yang keluar dari berbagi MCP mempertahankan server grup di `.claude.json` miliknya.

Perubahan yang akan memindahkan data akun yang sedang berjalan ditolak dengan menyebut nama akunnya, bahkan jika akun itu dimulai selagi rencana terbuka, dan rencana yang sudah usang karena proses ccam lain mengubah grup ditolak dengan `the groups changed while this plan waited; nothing was moved, make the change again`. Selama proses ccam lain mengubah grup, perubahan baru, penggantian nama, atau penghapusan menampilkan `another ccam process is changing the groups right now`; setelah perubahan yang terhenti di tengah jalan, semuanya menunggu sampai `F` di tab Health (atau `ccam doctor --fix`) menyelesaikannya.

## Menambahkan akun

`a` membuka wizard yang juga dipakai `ccam add` di terminal:

1. Nama akun.
2. Cara login: langganan Claude sekarang (di browser), [kunci API](/features/api), atau langganan Claude nanti.
3. Untuk kunci API: dari mana ccam membaca kunci itu setiap kali dijalankan (`env:NAME` atau `command:COMMAND`), endpoint, cara kunci dikirim, dan modelnya.
4. Grup yang akan diikuti, dengan `default` sudah dicentang (jika belum ada grup, apa yang akan dibagikan grup `default`).
5. Perintah pintasan, jika Anda menginginkannya.
6. Ringkasan untuk dikonfirmasi.

`Esc` kembali satu langkah, dan pada langkah pertama menutup wizard.

## Menghapus akun

`x` bertanya sebelum apa pun terjadi, dengan No terpilih, dan menjelaskan apa yang dilakukan penghapusan: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` Untuk akun kunci API tidak ada login yang perlu dicabut, dan jika akun lain memegang login yang sama, akun itu tidak di-logout. Akun yang diimpor mempertahankan foldernya: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## Warna

Warna mengikuti peran warna di UI terminal Claude Code, sehingga ccam terlihat serasi di sampingnya; nama grup mendapat warna berdasarkan namanya. ccam menanyakan warna latar belakang ke terminal lalu memilih tema terang atau gelap. Tetapkan sendiri di bagian atas `config.toml`, di atas `[[group]]`, `[[profile]]`, atau `[[folder]]` pertama:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Tema `-daltonized` ditujukan untuk buta warna merah-hijau, dan tema `-ansi` memakai 16 warna milik terminal itu sendiri. `NO_COLOR` mematikan warna.

## Keamanan dan kasus khusus

- Dasbor memanggil fungsi yang sama dengan perintah-perintahnya, jadi setiap aturan tetap berlaku: rencana sebelum data dipindahkan, penolakan untuk akun yang sedang berjalan, cadangan, dan token login yang tidak pernah meninggalkan akunnya.
- Dasbor tidak punya uji coba (dry run): `ccam --dry-run` tanpa nama mencetak bagaimana akun folder ini akan dijalankan, dan `o` pada sebuah akun menawarkan uji coba untuk satu kali menjalankan akun itu saja.
- Terminal dipulihkan sebelum Claude Code dimulai atau login dibuka.
- Layar yang terlalu kecil untuk sebuah dialog lebih dulu membuang baris-baris tab; baris dipotong dengan `…` alih-alih dibungkus ke baris berikutnya. Di layar yang pendek, baris kosong dan deskripsi tab dibuang berikutnya, sehingga baris petunjuk tombol tetap terlihat. Di bawah 40 kolom atau 8 baris, dasbor hanya menampilkan `Too small: ccam needs 40x8 (now 36x7)`.
- `Ctrl+C` saat sebuah perubahan, penghapusan, pembaruan, atau perbaikan sedang diterapkan akan menunggunya selesai: baris status menampilkan `ccam quits when it is done.`
- Setelah `U` memperbarui ccam, tawarannya tidak muncul lagi saat dasbor kembali dari menambahkan akun atau login.

## Terkait

- [Berbagi dengan grup](/features/share)
- [Menjalankan Claude Code](/features/launch)
- [Mengelola akun](/features/accounts)
- [Batas penggunaan](/features/usage)
