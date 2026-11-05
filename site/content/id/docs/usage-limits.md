---
title: Batas penggunaan dan perpindahan akun
description: Lihat penggunaan 5 jam dan mingguan setiap akun Claude Code dengan ccam usage, lalu buka akun dengan sisa kuota terbanyak lewat ccam next --auto.
updated: 2026-11-05
order: 5
translation: machine
sourceHash: "06bd349073e3"
---

`ccam usage` menampilkan, untuk setiap akun, seberapa banyak setiap batas yang sudah digunakan: batas 5 jam, batas mingguan, dan batas mingguan khusus setiap model (seperti Fable). `ccam next --auto` membuka akun dengan sisa kuota terbanyak, dan `--continue` membawa percakapan saat ini ke akun tersebut.

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

Setiap akun mendapat satu baris yang menunjukkan umur angka-angkanya, lalu satu bilah per batas beserta waktu hingga direset: jendela 5 jam, jendela mingguan untuk semua model, dan jendela mingguan khusus setiap model, seperti Fable. Bilah berubah kuning pada 80% dan merah pada 100%. Baris terakhir menyebut akun yang akan dibuka `ccam next --auto`; dengan `auto_strategy = "expiring"` atau `"balance"`, baris itu menjadi `→ next --auto picks <name>: <reason>`.

**Kursi Team dan Enterprise.** Pemilihan akun otomatis (`ccam next`, `--fallback`, `ccam pick`, tombol `n` di dasbor, dan, setelah Anda mengaktifkannya, [`auto_switch`](/features/next#auto-switch)) juga berganti di antara kursi dalam organisasi Claude Team atau Enterprise yang sama, yang dapat melanggar ketentuan Anthropic. Baca [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) dan tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya.

## Dari mana angkanya berasal

ccam menanyakan endpoint penggunaan Anthropic untuk semua akun secara bersamaan, menggunakan login masing-masing akun. ccam menunggu paling lama 3 detik per akun dan menampilkan setiap akun begitu akun itu menjawab.

Endpoint ini tidak didokumentasikan secara publik dan dapat berubah. Jika permintaan gagal, ccam menampilkan angka terakhir yang diketahui dan menandainya sebagai cache. Jendela yang sudah direset sejak angka-angka itu disimpan menampilkan 0%. `ccam doctor` memberi tahu Anda jika format responsnya berubah.

Saat Claude Code menghentikan sebuah sesi karena batas penggunaan, [plugin ccam](/docs/install#claude-code-plugin) mencatat jendela mana yang penuh dan kapan jendela itu direset. Sampai saat itu, jendela tersebut menampilkan 100% apa pun kata angka lainnya, sehingga `ccam next --auto` dan `--fallback` tidak memilih akun itu lagi terlalu cepat.

Angka yang diambil kurang dari 60 detik yang lalu akan dipakai ulang. Opsi:

| Opsi | Efek |
| --- | --- |
| `--watch[=seconds]` | Memperbarui terus-menerus, setiap 60 detik secara default |
| `--json` | Keluaran yang dapat dibaca mesin |
| `-v` | Juga menampilkan asal angka setiap akun |
| `--no-cache` | Selalu bertanya ulang |
| `--source cache` | Tidak pernah bertanya; hanya menampilkan angka yang tersimpan |

## Login yang kedaluwarsa

Saat token akses sebuah akun sudah kedaluwarsa, ccam tidak memperbaruinya, karena memperbaruinya saat sesi Claude Code sedang memakai login yang sama dapat membuat sesi itu keluar. Akun itu menampilkan ``access token expired; `ccam <name>` renews it when it starts``: jalankan sekali dengan `ccam <name>` dan Claude Code akan memperbaruinya. Login yang sudah kedaluwarsa sepenuhnya, atau yang dihapus oleh Claude Code, menampilkan `` login expired → `ccam login <name>` `` dan memerlukan login baru.

`ccam usage --refresh` memperbarui token yang kedaluwarsa sendiri, tetapi hanya untuk akun yang tidak memiliki sesi Claude Code yang sedang berjalan, bagaimanapun sesi itu dijalankan. Perintah ini tidak pernah memperbarui login yang tersimpan di Rantai Kunci (Keychain) macOS atau di `~/.claude` yang diimpor, dan satu pembaruan bisa memakan waktu hingga 30 detik.

## Penggunaan di status line

```sh
ccam statusline install
```

Claude Code kemudian menampilkan nama akun dan penggunaannya di status line, misalnya `[work] · Opus · 5h 24% · 7d 41%`. Mulai dari 90%, baris itu juga menyarankan `ccam next --auto --continue`. Jika Anda sudah memiliki perintah status line, ccam mempertahankannya dan menaruh nama akun di depannya. Status line juga menyimpan angkanya, yang dipakai `ccam usage` saat endpoint tidak menjawab.

## Berpindah akun

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the most room
ccam next --continue  # continue this folder's conversation on the next account
```

`ccam next` melewati akun yang belum login, dan `--auto` melewati akun yang sudah mencapai 100% (termasuk jendela yang dilaporkan penuh oleh Claude Code, sampai di-reset) serta akun dari organisasi lain. Jika kedua akun berbagi riwayat percakapan, `--continue` membuka percakapan yang sama; jika tidak, Claude Code menyalinnya ke akun baru. Lihat [pembagian](/docs/sharing) dan [beralih ke akun berikutnya](/features/next#continue).

Dengan `auto_switch = true` di bagian atas `config.toml` (Linux dan macOS; nonaktif secara default), ccam melakukan perpindahan ini sendiri saat Claude Code menghentikan sesi interaktif di batasnya dan akun lain yang sejenis (kursi di organisasi yang sama, atau paket pribadi lain) punya ruang. Lihat [`auto_switch`](/features/next#auto-switch) untuk kapan fitur ini bertindak dan apa yang tidak pernah dilakukannya.

Untuk skrip dan tugas terjadwal yang menjalankan `claude -p`, letakkan `--fallback` sebelum nama akun. Saat akun itu mencapai batasnya, ccam menjalankan lagi perintah yang sama sebagai akun dengan sisa kuota terbanyak, dan hanya output dari proses yang selesai yang dicetak:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Lihat [Jalankan ulang saat batas tercapai](/features/next#run-again-at-the-limit) untuk cara ccam memutuskannya.

[Panduan alur kerja batas penggunaan](/guides/usage-limit-workflow) menjelaskan contoh lengkapnya.
