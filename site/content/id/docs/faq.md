---
title: Pertanyaan yang sering diajukan
metaTitle: "FAQ ccam: beberapa akun Claude Code"
description: Jawaban seputar ccam, akun Claude Code, batas penggunaan, pembagian, privasi, dan sistem yang didukung.
updated: 2026-11-19
order: 6
translation: machine
sourceHash: "090ecc3a34bd"
---

## Apakah ccam alat resmi dari Anthropic?

Tidak. ccam adalah proyek open source tidak resmi dan tidak berafiliasi dengan atau didukung oleh Anthropic. ccam menggunakan pengaturan `CLAUDE_CONFIG_DIR` yang [didokumentasikan Claude Code](https://code.claude.com/docs/en/env-vars) untuk menjalankan beberapa akun secara berdampingan.

## Apakah saya boleh menggunakan beberapa akun Claude Code?

Itu bergantung pada ketentuan setiap akun dan, untuk akun kerja, kebijakan organisasi Anda: [Ketentuan Konsumen](https://www.anthropic.com/legal/consumer-terms) Anthropic untuk Pro dan Max, [Ketentuan Komersial](https://www.anthropic.com/legal/commercial-terms)-nya untuk Team dan Enterprise, dan [Kebijakan Penggunaan](https://www.anthropic.com/legal/aup)-nya untuk semuanya. Periksa hal tersebut sebelum Anda menggunakan beberapa akun. ccam tidak menaikkan, menggabungkan, atau menyembunyikan batas sebuah akun, dan setiap akun tetap berhenti di batasnya sendiri. ccam memindahkan pekerjaan Anda ke akun Anda yang lain saat Anda memintanya, atau dengan sendirinya saat mencapai batas jika Anda mengaktifkan [`auto_switch`](/features/next#auto-switch) (Linux dan macOS). Periksa apakah hal itu diizinkan untuk akun-akun Anda sebelum menggunakannya; keputusannya ada di tangan Anda.

**Kursi Team dan Enterprise.** Pemilihan akun otomatis (`ccam next`, `--fallback`, `ccam pick`, [`ccam agent`](/features/agent), tombol `n` di dasbor, dan [`auto_switch`](/features/next#auto-switch) bila Anda mengaktifkannya) juga berganti di antara kursi dalam organisasi Claude Team atau Enterprise yang sama. Pusat bantuan Anthropic menyatakan "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (setiap kursi pada paket Team diberikan kepada satu orang tertentu, dan hanya orang itu yang boleh menggunakannya) dan "Users can't route their work through another user's account to use that user's remaining usage" (pengguna tidak boleh menyalurkan pekerjaannya melalui akun pengguna lain untuk memakai sisa kuota pengguna tersebut) ([paket Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); begitu pula untuk [paket Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Karena itu, memakai pemilihan otomatis di antara kursi dapat melanggar ketentuan Anthropic. Apakah Anda memakainya di sana adalah keputusan Anda: tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya. Paket pribadi (Pro, Max) dibahas dalam paragraf di atas.

## Apakah ccam mengakali batas penggunaan?

Tidak melewati batas milik akun mana pun: setiap akun tetap memiliki batasnya dan berhenti di sana, dan ccam tidak menaikkan, menggabungkan, atau menyembunyikannya. ccam menampilkan batas tersebut dan memungkinkan Anda melanjutkan di akun Anda yang lain saat satu akun habis kuotanya; dengan `auto_switch` aktif, ccam melakukannya dengan sendirinya (lihat pertanyaan berikutnya). Di antara kursi dalam satu organisasi Team atau Enterprise, melanjutkan di kursi lain memakai sisa kuota kursi itu, dan pusat bantuan Anthropic melarang hal itu untuk kursi milik pengguna lain (lihat Kursi Team dan Enterprise di atas).

## Bisakah ccam berpindah akun dengan sendirinya saat sebuah sesi mencapai batasnya?

Hanya jika Anda mengaktifkannya. Dengan `auto_switch = true` di bagian atas `~/.config/ccam/config.toml` (Linux dan macOS; nonaktif secara default dan masih eksperimental), saat Claude Code menghentikan sesi interaktif (yang dimulai oleh ccam) di batas 5 jam, mingguan, Opus, Sonnet, atau Fable, ccam mengakhiri sesi itu dan melanjutkan percakapan sebagai akun Anda yang lain yang punya ruang: kursi lain dari organisasi Team atau Enterprise yang sama, atau paket pribadi lain (Pro, Max). ccam tidak pernah berpindah sebelum batas tercapai, tidak pernah memantau penggunaan Anda di latar belakang, dan tidak pernah memindahkan pekerjaan keluar dari organisasi Team atau Enterprise ataupun ke dalamnya.

Pada kursi, ccam memindahkan pekerjaan Anda dari satu kursi ke kursi lain, yang dapat melanggar ketentuan Anthropic (lihat Kursi Team dan Enterprise di atas). Untuk paket pribadi, ketentuan Anthropic tidak membahasnya dengan jelas: Ketentuan Konsumen melarang mengakses Layanan "through automated or non-human means" (melalui cara otomatis atau non-manusia) kecuali Anthropic secara tegas mengizinkannya, dan sebuah perpindahan mengirim satu pesan atas nama Anda. Mengaktifkannya adalah keputusan Anda.

[`auto_switch`](/features/next#auto-switch) mencantumkan setiap syarat dan biaya sebuah perpindahan.

## Bagaimana cara menggunakan dua akun Claude Code di satu komputer?

Beri setiap akun folder konfigurasinya sendiri dengan `CLAUDE_CONFIG_DIR`. ccam melakukannya untuk Anda: jalankan `ccam add` untuk setiap akun, lalu jalankan akun-akun itu dengan `ccam personal` dan `ccam work`, atau dengan pintasan yang namanya Anda pilih sendiri, seperti `ccwork`. [Panduan ini](/guides/multiple-claude-code-accounts) juga menunjukkan cara melakukannya secara manual.

## Bagaimana cara melanjutkan percakapan di akun lain?

Jalankan `ccam next --continue` di folder itu: perintah ini membawa percakapan terakhir folder tersebut ke akun berikutnya (`--auto` untuk akun dengan sisa kuota terbanyak). Di dalam Claude Code, `/ccam:switch` dari [plugin ccam](/docs/install#claude-code-plugin) melakukan hal yang sama untuk sesi yang sedang Anda jalani. Jika kedua akun berbagi riwayat percakapan (`ccam group set default --add history`, atau dengan grup tersendiri: `ccam group add pair --items history --members work,side`), akun baru membuka percakapan yang sama persis; jika tidak, Claude Code menyalinnya ke riwayat akun baru. Tambahkan `--brief` untuk memulai akun baru dengan ringkasan singkat alih-alih seluruh percakapan.

## Apakah ccam melihat atau mengirim login saya?

ccam tidak pernah membagikan atau menyalin login Claude Anda antarakun, dan tidak pernah menyalin info masuk server MCP. ccam membaca token akses sebuah akun hanya untuk menanyakan penggunaan akun itu ke endpoint penggunaan Anthropic: untuk `ccam usage`, dasbor, `ccam doctor`, dan pilihan otomatis (`ccam next --auto`, `ccam pick`, `--fallback`, dan `auto_switch` saat sebuah sesi mencapai batasnya). ccam memperbarui token yang kedaluwarsa hanya dengan `ccam usage --refresh`; dengan `usage_source = "claude"`, ccam tidak mengirim token ke mana pun dan membaca login yang tersimpan (file-nya, atau entri Rantai Kunci (Keychain) di macOS) hanya untuk mengetahui apakah login itu masih berfungsi. `ccam ext` dan `ccam doctor` membaca nama-nama info masuk MCP sebuah akun, tidak pernah tokennya. ccam tidak memiliki telemetri. Selain pengecekan itu, satu-satunya permintaan yang dibuatnya sendiri adalah pengecekan dasbor ke GitHub untuk rilis ccam yang lebih baru, yang tidak membawa apa pun tentang Anda atau akun-akun Anda.

## Bagaimana cara masuk di server jarak jauh?

Jalankan `ccam add` atau `ccam login <name>` di sesi SSH seperti biasa. Claude Code menunggu proses masuk di server, sehingga browser di komputer Anda berakhir di halaman `http://localhost:<port>/callback?code=…` yang tidak bisa dimuat. Salin alamat itu dari browser lalu tempelkan ke terminal: ccam meneruskannya ke Claude Code di server, yang kemudian menyelesaikan proses masuk. Alamat itu berisi kode sekali pakai, jadi jangan bagikan kepada siapa pun.

## Di mana akun saya disimpan?

Di `~/.local/share/ccam/profiles/<name>/`, satu folder per akun. Di macOS, Claude Code menyimpan login itu sendiri di Rantai Kunci (Keychain), dengan nama yang diturunkan dari folder tersebut, sehingga `ccam rename` mempertahankan nama lama folder itu di sana.

## Apa yang terjadi pada ~/.claude saya yang sudah ada?

Tidak ada, kecuali Anda menjalankan `ccam import` atau `ccam migrate`. Impor mendaftarkan folder itu di tempatnya, sehingga `claude` biasa dan `ccam <name>` tetap menggunakan akun dan pengaturan yang sama. `ccam migrate` menyalin pengaturan, `CLAUDE.md`, ekstensi, dan server MCP dari folder itu ke grup bersama `default` dan hanya membaca folder tersebut.

## Sistem apa saja yang didukung ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine, dan lainnya), macOS, serta Windows 10 dan 11, secara native atau di dalam WSL. Di Windows, membagikan file seperti `settings.json` memerlukan Mode Pengembang (Developer Mode); folder bersama berfungsi tanpanya, dan tanpanya `ccam add` pertama membuat grup `default` yang hanya berisi server MCP dan ekstensi ([berbagi di Windows](/features/share#on-windows)).

## Mengapa server MCP hanya muncul di satu akun?

Server yang ditambahkan di dalam sesi dengan `claude mcp add` hanya disimpan untuk akun tersebut, dan tanpa `-s user` hanya untuk folder saat ini. ccam memberi tahu Anda saat berikutnya Anda menjalankan akun itu, dan `ccam mcp list` serta `ccam ext` menampilkannya; pindahkan server ke grup bersama dengan `ccam mcp import <account> <server>`.

## Apakah ccam memperlambat Claude Code?

Menjalankan Claude Code melalui ccam menambah sekitar 4 milidetik, diukur dengan `scripts/perf-check.sh` di sebuah laptop. ccam membaca satu file kecil, memeriksa tautan bersama, lalu menggantikan dirinya dengan `claude`, sehingga secara default ccam tidak berjalan saat Anda bekerja. Windows tidak bisa mengganti sebuah proses, jadi di sana ccam menjalankan Claude Code dan menunggunya, tanpa melakukan hal lain. Dengan `relay` atau `auto_switch` diaktifkan, ccam tetap berjalan dengan cara yang sama di samping setiap sesi interaktif, sehingga ccam bisa melanjutkan percakapan di akun lain.

## Bagaimana cara memperbarui ccam?

Jalankan `ccam update`. Perintah ini mengunduh rilis terbaru, memeriksa checksum-nya, lalu mengganti file `ccam`; `ccam update --check` hanya memberi tahu apakah ada rilis yang lebih baru. Dasbor juga memberi tahu jika ada rilis yang lebih baru, dan `U` menginstalnya setelah bertanya; ccam tidak pernah memperbarui dirinya sendiri. Jika Anda menginstal paket `.deb`, instal `.deb` yang baru dengan apt. Lihat [pembaruan](/docs/install#update).

## Bagaimana cara menghapus instalasi ccam?

Jalankan `ccam uninstall`. Perintah ini menampilkan apa saja yang akan dihapus dan bertanya terlebih dahulu, lalu menghapus setiap akun (di-logout terlebih dahulu), perintah pintasan, `~/.config/ccam`, `~/.local/share/ccam`, `~/.local/state/ccam`, dan perintah `ccam`. `~/.claude` yang diimpor tetap ada, dengan salinan dari apa yang dibagikannya di grup-grupnya. Lihat [instalasi](/docs/install#uninstall).
