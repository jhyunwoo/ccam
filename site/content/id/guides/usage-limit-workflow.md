---
title: Saat Claude Code mencapai batas penggunaan
description: Tunggu reset atau langsung lanjutkan percakapan yang sama di akun Claude Code lain. Cara kerja batas 5 jam dan mingguan, langkah demi langkah.
updated: 2026-11-01
order: 3
translation: machine
sourceHash: "760a36724feb"
---

Saat sebuah akun mencapai batas penggunaannya (sering disebut rate limit), Anda punya dua pilihan: menunggu sampai batasnya direset, atau melanjutkan percakapan yang sama sekarang juga di akun lain milik Anda. Dengan ccam, pilihan kedua cukup satu perintah, `ccam next --auto --continue`.

## Batas apa saja yang dimiliki Claude Code?

Pada paket Pro dan Max, Claude Code memiliki [batas sesi 5 jam dan batas mingguan](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Keduanya dipakai bersama dengan Claude di web dan di aplikasi, dan penggunaan dihitung terhadap keduanya sekaligus. Saat salah satunya habis, Claude Code [menampilkan pesan](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) seperti "You've hit your session limit · resets 3:45pm". Berganti model tidak membantu, karena batasnya berlaku untuk semua model.

## Opsi 1: tunggu reset

Jika Anda hanya punya satu akun, tunggu. Versi terbaru Claude Code [menunggu di sesi yang terbuka](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) dan melanjutkan tugas dengan sendirinya setelah batasnya direset. `/usage` di dalam Claude Code menampilkan seberapa banyak batas paket Anda yang sudah terpakai.

## Opsi 2: lanjutkan di akun lain sekarang

Jika Anda punya akun kedua, Anda bisa melanjutkan percakapan yang sama saat itu juga. Siapkan sekali:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Status line kini menampilkan akun dan penggunaan 5 jamnya selama Anda bekerja, sehingga batas tidak datang tiba-tiba.

Saat batas tercapai, keluar dari Claude Code (`/exit`), lalu:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Buka percakapan di akun dengan sisa kuota terbanyak:

```sh
ccam next --auto --continue
```

Claude Code dimulai sebagai `work` dengan percakapan terbaru dari folder saat ini, termasuk semua yang sudah dibicarakan. Berbagi riwayat (`ccam group set default --add history` di atas) berarti kedua akun mengerjakan percakapan yang sama; tanpa itu, Claude Code menyalin percakapan ke `work`, dan ccam memberi tahu berapa token yang harus dibaca ulang oleh balasan pertama.

Dengan [plugin ccam](/docs/install#claude-code-plugin) di Claude Code, sesi yang berhenti karena batasnya menampilkan notifikasi dan mencatat jendela yang penuh, sehingga `ccam next --auto` tidak memilih akun itu lagi sebelum di-reset. `/ccam:switch` (atau `/ccam:switch work`) mencatat percakapan itu untuk `ccam next --continue`, sehingga perintah itu melanjutkan tepat percakapan tersebut.

**Kursi Team dan Enterprise.** Pemilihan akun otomatis (`ccam next`, `--fallback`, `ccam pick`, tombol `n` di dasbor, dan, setelah Anda mengaktifkannya, [`auto_switch`](/features/next#auto-switch)) juga berganti di antara kursi dalam organisasi Claude Team atau Enterprise yang sama, yang dapat melanggar ketentuan Anthropic. Baca [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) dan tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya.

## Apakah ccam berganti akun secara otomatis?

Hanya jika Anda mengaktifkannya.

- Secara default, ccam tidak berjalan selama Anda bekerja: ccam menjalankan Claude Code lalu menggantikan dirinya dengan Claude Code, jadi Anda yang menentukan kapan berganti. `ccam next --auto` memilih akun dengan penggunaan 5 jam terendah, melewati akun yang sudah 100% atau belum login, dan tetap berada dalam satu organisasi Team atau Enterprise, atau di antara paket pribadi.
- Dengan `relay = true` di `config.toml`, ccam justru tetap berjalan di samping sesi, sehingga `/ccam:switch` melanjutkan percakapan dengan sendirinya, dengan mempertahankan opsi, mode izin, tingkat effort, dan agen sesi tersebut; sebelum memindahkan percakapan ke akun organisasi lain, ccam bertanya kepada Anda.
- Dengan `auto_switch = true` (Linux dan macOS; eksperimental), ccam juga melakukan perpindahan itu sendiri: saat Claude Code menghentikan sesi di batasnya dan akun lain yang sejenis punya ruang, ccam mengakhiri sesi itu dan melanjutkan percakapan di sana ([syarat dan biaya](/features/next#auto-switch)).

Beberapa alat lain justru menukar login sesi yang sedang berjalan; [ccam dan pengalih akun lainnya](/compare/account-switchers) menjelaskan perbedaannya.

## Hal yang perlu diperhatikan

- Jalankan `ccam next --continue` di folder yang sama tempat Anda bekerja; percakapan disimpan per folder proyek.
- Hindari menjalankan dua akun pada percakapan yang sama secara bersamaan. ccam memperingatkan Anda jika Anda mencobanya.
- Tanpa riwayat bersama, percakapan disalin ke akun lain dan kedua salinan itu berjalan sendiri-sendiri; checkpoint yang dipakai `/rewind` untuk memulihkan kode tetap berada di akun pertama.
- Pesan pertama di akun baru bisa memakai lebih banyak batasnya daripada biasanya, karena akun itu belum menyimpan percakapan tersebut di cache. `ccam next --continue --brief` justru memulainya dengan ringkasan singkat.
- Boleh tidaknya Anda berganti akun seperti ini bergantung pada ketentuan yang berlaku untuk masing-masing akun.

Lihat [Batas penggunaan](/features/usage) dan [Beralih ke akun berikutnya](/features/next) untuk setiap opsi.
