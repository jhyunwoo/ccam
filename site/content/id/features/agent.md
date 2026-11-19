---
title: Subagen di akun lain
description: ccam agent menjalankan sebagian tugas Claude Code sebagai akun lain Anda dan memberikan balasan Claude, sehingga tugas panjang memakai penggunaan beberapa akun.
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Subagen milik Claude Code sendiri berjalan sebagai akun dari sesi yang memulainya, sehingga mereka memakai penggunaan akun itu. `ccam agent "<task>"` menjalankan sebuah tugas dengan `claude -p` di folder yang sama sebagai salah satu akun Anda yang lain dan mencetak balasan terakhir Claude. Jika dipanggil dari sesi Claude Code, perintah ini bekerja seperti subagen yang memakai penggunaan akun lain, sehingga satu tugas panjang bisa terus berjalan setelah akun sesi itu sendiri seharusnya sudah kehabisan.

**Kursi Team dan Enterprise.** `ccam agent` memilih akun dengan sendirinya. Di dalam satu organisasi Claude Team atau Enterprise, perintah ini menjalankan tugas di kursi lain, yang dapat melanggar ketentuan Anthropic. Baca [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) dan tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya.

## Penggunaan

Claude Code menjalankannya melalui alat Bash miliknya. Dengan [plugin ccam](/docs/install#claude-code-plugin) terpasang, mintalah dengan kata-kata biasa, misalnya "pakai akun-akun saya yang lain untuk bagian-bagian tugas ini", atau tuliskan itu di CLAUDE.md Anda. Skill plugin itu meminta Claude Code memulai setiap agen di latar belakang, menulis tugas yang bisa berdiri sendiri, dan memberi agen yang mengedit file bagiannya sendiri atau git worktree-nya sendiri.

```sh
ccam agent "find why TestParse fails and fix it; reply with what you changed"
ccam agent --model opus -- "-v breaks the parser: find out why"
ccam agent --json < task.md            # the task from standard input
ccam agent --resume 7f3c9a2e-… "now add a test for it"
ccam agent --wait=1h "…"               # nobody has room: wait for a reset
ccam --dry-run agent "…"               # the account and command it would run
```

```console
$ ccam agent "summarize what internal/usage does"
ccam: agent runs as work (12% of its 5 hours used); session 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f
internal/usage looks up the 5-hour and weekly usage of every account …
ccam: agent finished as work; follow up with: ccam agent --resume 7f3c9a2e-5b1d-4c2a-9e8f-1a2b3c4d5e6f "…"
```

Balasan terakhir Claude dikirim ke output standar; baris yang diawali `ccam:` dikirim ke error standar. Sebagai gantinya, `--json` mencetak `{account, session_id, result, is_error, accounts}`, dengan `accounts` berisi setiap akun yang menjalankan tugas itu.

Kode keluar: kode milik claude sendiri (`0` jika agen selesai), `75` (`EX_TEMPFAIL`) jika tidak ada akun selain akun pemanggil yang punya ruang, `1` untuk galat lain seperti agen yang mencoba memulai agen lain.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--model <model>` | Model agen; jendela mingguan model itu juga ikut dihitung dalam pemilihan | model default akun |
| `--resume <id>` | Melanjutkan percakapan agen sebelumnya dengan tugas lanjutan (ID sesi ada di baris pertama dan terakhirnya) | nonaktif |
| `--json` | Mencetak akun, sesi, balasan, dan akun-akun yang menjalankannya | nonaktif |
| `--max-usage <1-100>` | Utamakan akun di bawah penggunaan ini | `auto_max_usage`, jika tidak, 100 |
| `--strategy headroom\|expiring\|balance` | [Aturan pemilihan](/features/next#auto) | `auto_strategy`, jika tidak, `headroom` |
| `--cross-org` | Juga menjalankan sebagai akun dari organisasi lain | `cross_org`, jika tidak, nonaktif |
| `--wait[=<duration>]` | Jika tidak ada akun yang punya ruang, tunggu akun pertama di-reset, paling lama selama ini | nonaktif; `5h` jika diberikan tanpa nilai |
| `--dry-run` | Mencetak perintah untuk akun yang akan dipilih, tanpa menjalankannya | nonaktif |

## Cara kerjanya

1. **Pemanggil tidak diikutkan.** Akun tempat `ccam agent` dipanggil (akun dari sesi yang dimulai ccam, atau jika tidak ada, akun yang dijalankan folder ini) tidak pernah menjalankan agen, begitu pula folder lain yang login ke akun Claude yang sama, karena folder itu berbagi batas dengannya. Agen tetap berada di [organisasi](/features/next#auto) pemanggil, kecuali `--cross-org` menyatakan lain.
2. **Pemilihannya seperti `ccam pick`.** Penggunaan ditanyakan dengan cara yang sama seperti [`ccam pick`](/features/pick), dan agen yang dimulai dalam selang 2 menit satu sama lain mendapat akun berbeda selama masih cukup akun yang punya ruang.
3. **Agen berjalan dengan izin sesi itu, tidak pernah lebih.** Di `claude -p` tidak ada orang yang bisa mengizinkan sebuah alat, jadi agen mendapat mode izin dan tingkat effort dari sesi yang memanggilnya. Jika dijalankan dari terminal, yang berlaku adalah mode default Claude Code dan aturan izin di pengaturan Anda. `ccam agent` tidak punya opsi yang memberi izin lebih.
4. **Saat batas penggunaan tercapai, percakapan berlanjut.** Saat akun agen mencapai batasnya, ccam mencatat batas itu seperti yang akan dilakukan [hook plugin](/features/usage#the-limit-record) dan melanjutkan percakapan yang sama sebagai akun dengan ruang terbanyak: Claude Code menyalinnya ke sana (`--resume <file> --fork-session`), dengan pesan pertama "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." Setiap akun dicoba satu kali.
5. **Setiap percakapan agen dicatat** di `~/.local/state/ccam/agent/<session id>.json` (30 hari). `--resume <id>` melanjutkannya di foldernya, sebagai akun yang menjalankannya jika akun itu masih punya ruang, jika tidak, sebagai akun lain. [`ccam next --continue`](/features/next) tidak mengikutkan percakapan agen, sehingga perintah itu tetap melanjutkan percakapan Anda.

## Keamanan dan kasus khusus

- Agen tidak bisa memulai agennya sendiri: ccam menandainya dengan `CCAM_AGENT=1`, dan `ccam agent` menolak berjalan di sana.
- Menghentikan `ccam agent` (Ctrl+C, sinyal, atau Claude Code yang menyela perintah itu) juga menghentikan agennya.
- Agen tidak melihat percakapan yang memanggilnya, hanya tugasnya.
- Agen yang mengedit file yang sama secara bersamaan saling menimpa hasil kerja. Beri masing-masing bagian tugasnya sendiri, atau git worktree-nya sendiri.
- Saat sesi pemanggil pindah ke akun lain ([`auto_switch`](/features/next#auto-switch)), agen yang dimulainya ikut berhenti. Percakapan mereka tetap ada: `--resume <id>` melanjutkannya.
- `agent` adalah nama perintah, jadi tidak ada akun baru yang bisa diberi nama `agent`; akun yang sudah ada dengan nama itu tetap bisa dijalankan dengan `ccam run agent`.

## Terkait

- [Memilih akun untuk skrip](/features/pick)
- [Beralih ke akun berikutnya](/features/next)
- [Batas penggunaan](/features/usage)
- [Sesi latar belakang semua akun](/features/agents)
