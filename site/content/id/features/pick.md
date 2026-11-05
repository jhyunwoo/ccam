---
title: Memilih akun untuk skrip
description: ccam pick mencetak akun dengan ruang terbanyak untuk skrip, memberi akun berbeda ke worker yang dimulai bersamaan, dan keluar 75 jika tak ada ruang.
updated: 2026-11-05
order: 12
translation: machine
sourceHash: "0ec25752ace9"
---

`ccam pick` mencetak nama akun yang sudah login dengan ruang terbanyak, sehingga skrip bisa menjalankan Claude Code dengannya. Worker yang dimulai pada waktu yang sama masing-masing mendapat akun berbeda selama masih cukup akun yang punya ruang, sehingga sekumpulan tugas `claude -p` tersebar ke akun-akun Anda alih-alih menghabiskan satu akun.

**Kursi Team dan Enterprise.** Pemilihan akun otomatis (`ccam next`, `--fallback`, `ccam pick`, tombol `n` di dasbor, dan, setelah Anda mengaktifkannya, [`auto_switch`](/features/next#auto-switch)) juga berganti di antara kursi dalam organisasi Claude Team atau Enterprise yang sama, yang dapat melanggar ketentuan Anthropic. Baca [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) dan tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya.

## Penggunaan

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # most weekly quota left per day; fewer sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

Kode keluar: `0` dengan nama di output standar, `75` (`EX_TEMPFAIL`) jika tidak ada akun yang punya ruang atau tidak ada yang sudah login, `1` untuk galat lain seperti flag yang salah.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--json` | Mencetak akun, penggunaan 5 jam dan mingguannya, asal angka-angka itu, dan alasan akun itu dipilih | nonaktif |
| `--model <model>` | Juga menghitung jendela mingguan model itu, seperti `opus` | tidak ada |
| `--max-usage <1-100>` | Utamakan akun di bawah penggunaan ini | `auto_max_usage`, jika tidak, 100 |
| `--strategy headroom\|expiring\|balance` | Penggunaan mingguan paling sedikit di antara akun yang penggunaan 5 jamnya paling banyak 3 poin di atas penggunaan 5 jam terendah, reset 5 jam yang paling cepat, atau sisa kuota mingguan terbanyak per hari | `auto_strategy`, jika tidak, `headroom` |
| `--cross-org` | Juga memilih akun dari organisasi lain | `cross_org`, jika tidak, nonaktif |
| `--wait[=<duration>]` | Jika tidak ada akun yang punya ruang, tunggu akun pertama di-reset, paling lama selama ini. Nilainya ditaruh setelah `=`: `--wait 30m` ditolak dengan petunjuk `--wait takes its value after '=': --wait=30m` | nonaktif; `5h` jika diberikan tanpa nilai |
| `--dry-run` | Memilih tanpa mencatat pilihan itu | nonaktif |

## Cara kerjanya

1. ccam membuat daftar akun yang sudah login (akun kunci API dan login Console tidak punya batas dan tidak diikutkan) di organisasi akun yang akan dijalankan folder ini (lihat [organisasi](/features/next#auto)); folder yang login ke akun Claude yang sama dengan akun itu juga tidak diikutkan.
2. ccam mengambil kunci (`~/.local/state/ccam/picks.lock`), sehingga pemilihan yang berjalan bersamaan dilakukan satu per satu.
3. ccam menanyakan penggunaan semua akun sekaligus, dengan cara yang sama seperti `ccam usage`, memakai cache 60 detik. Jendela yang dilaporkan penuh oleh Claude Code dihitung 100% sampai di-reset ([catatan batas](/features/usage#the-limit-record)), sehingga akun yang baru saja mencapai batasnya dikecualikan meskipun angka tersimpannya masih menunjukkan ada ruang.
4. Akun yang dipilih dalam 2 menit terakhir, oleh `ccam pick` atau oleh [`auto_switch`](/features/next#auto-switch) yang memindahkan sesi ke sana (`~/.local/state/ccam/picks/<name>`; pilihan dengan stempel waktu di masa depan tidak dihitung) dikeluarkan lebih dulu, lalu [aturan pemilihan](/features/next#auto) `ccam next --auto` memilih di antara sisanya: akun di bawah `auto_max_usage` lebih dulu, menurut strategi, dan tidak pernah akun yang sudah 100%. Jika hal ini melewati akun yang seharusnya dipilih oleh aturan tersebut, `reason` dari `--json` menyebutkannya: `30% of its 5 hours used; work was picked a moment ago`.
5. Jika setiap akun yang punya ruang baru saja dipilih, aturan itu memilih lagi di antara semuanya, sehingga dengan tiga akun, worker keempat berbagi salah satunya.
6. Pilihan dicatat dan kunci dilepas, lalu nama akun dicetak.

Jendela 2 menit itu penting karena penggunaan tugas baru butuh waktu untuk muncul di angka: tanpa jendela itu, lima worker yang dimulai bersamaan akan melihat akun yang sama sebagai akun dengan ruang terbanyak.

## Keamanan dan kasus khusus

- `ccam pick` hanya membaca penggunaan dan menulis waktu pemilihan; ia tidak pernah memperbarui atau membaca token lebih dari yang dilakukan `ccam usage`.
- Pemilihan tidak memesan akun. Dua skrip tetap bisa dimulai di akun yang sama setelah 2 menit; gunakan `--fallback` untuk pulih saat akun itu mencapai batasnya.
- Menunggu dengan `--wait` tidak memegang kunci, sehingga pemilihan lain tetap berjalan sementara itu.
- `ccam pick` mengabaikan [`ccam folder`](/features/folder): ia memilih hanya berdasarkan penggunaan.

## Terkait

- [Beralih ke akun berikutnya](/features/next)
- [Jalankan ulang saat batas tercapai](/features/next#run-again-at-the-limit)
- [Batas penggunaan](/features/usage)
- [Akun kunci API](/features/api)
