---
title: Referensi perintah
metaTitle: Referensi perintah ccam
description: Semua perintah dan opsi ccam, dari ccam add dan ccam <name> hingga ccam usage, ccam next, dan ccam doctor.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: "28e1ad1dcbad"
---

Letakkan `--dry-run` sebelum perintah untuk melihat apa yang akan dilakukannya tanpa benar-benar menjalankannya, misalnya `ccam --dry-run work`. Opsi ini berfungsi saat menjalankan Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) dan dengan `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall`, serta `doctor --fix`. Perintah yang hanya membaca, seperti `list`, `usage`, `share`, dan `ext`, menerimanya dan memang tidak mengubah apa pun. Perintah lainnya menolaknya alih-alih berjalan, termasuk `group rename`.

<!-- commands -->

## Menjalankan Claude Code

`ccam <name|number> [arguments]` menjalankan Claude Code sebagai akun tersebut. ccam tidak membaca argumen setelah nama akun; ccam meneruskannya ke `claude` dengan urutan yang sama, dan input standar, kode keluar, serta sinyal berperilaku persis seperti pada `claude` itu sendiri. `--` setelah nama akun bersifat opsional.

`ccam` dengan argumen tetapi tanpa nama akun, seperti `ccam -p "summarise"`, menggunakan akun yang ditetapkan untuk folder saat ini, atau jika tidak ada, akun default. Tetapkan keduanya dengan `ccam folder <name> [folder]` dan `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` sebelum nama akun ditujukan untuk proses `-p`/`--print`. Saat akun mencapai batas penggunaannya, ccam menjalankan lagi perintah yang sama sebagai akun yang sudah login dengan sisa kuota terbanyak dan hanya mencetak output dari proses yang selesai. Lihat [Akun per folder](/features/folder) dan [Jalankan ulang saat batas tercapai](/features/next#run-again-at-the-limit).

**Kursi Team dan Enterprise.** Pemilihan akun otomatis (`ccam next`, `--fallback`, `ccam pick`, tombol `n` di dasbor, dan, setelah Anda mengaktifkannya, [`auto_switch`](/features/next#auto-switch)) juga berganti di antara kursi dalam organisasi Claude Team atau Enterprise yang sama, yang dapat melanggar ketentuan Anthropic. Baca [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) dan tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya.

Pintasan seperti `ccwork` menjalankan akunnya dengan cara yang sama: `ccwork --resume` melakukan hal yang sama dengan `ccam work --resume`. Lihat [Pintasan](/docs/commands#shortcuts) di bawah.

## Menambah dan menghapus akun

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` tanpa folder akan mencari `~/.claude` dan `~/.claude-*`, lalu menanyakan masing-masing; folder yang gagal dilaporkan sebagai `Skipped <path>: <reason>` dan yang lain tetap didaftarkan. Jika nama sebuah folder sudah dipakai, impor folder itu sendiri dengan `ccam import <folder> --name <name>`. Folder yang diimpor didaftarkan di tempatnya; ccam tidak pernah memindahkan atau menyalin login.

`ccam login --expired` juga me-login akun yang login-nya telah dihapus oleh Claude Code. `ccam remove` me-logout akun terkelola (akun yang foldernya dibuat oleh ccam, bukan yang didaftarkan dengan `ccam import`) sebelum foldernya masuk ke sampah, kecuali jika akun lain memegang login yang sama. Menghapus akun yang diimpor akan mengeluarkannya dari grup-grupnya dan meninggalkan salinan setiap item bersama di foldernya. Di macOS, `ccam rename` mempertahankan nama folder, karena Rantai Kunci menyimpan login di bawah path folder itu.

`ccam add --api` menambahkan akun yang memakai kunci API alih-alih login Claude; ccam membaca kunci itu setiap kali akun dijalankan dan tidak pernah menyimpannya. Lihat [akun kunci API](/features/api).

`ccam rename`, `ccam remove`, dan perintah `ccam group` menolak mengubah akun saat Claude Code sedang berjalan sebagai akun tersebut, bagaimanapun ia dijalankan (termasuk sesi latar belakang). Keluar dari sesi terlebih dahulu. Perintah-perintah itu juga menolak selama perubahan grup belum selesai: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` menampilkan setiap akun beserta slot, pintasan, email, dan foldernya. Array `groups` di dalamnya mencantumkan grup-grup akun tersebut, dan objek `share` memetakan setiap item (`mcp`, `settings`, `memory`, `extensions`, `history`) ke grup asal item itu, atau `local`.

## Pintasan

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Pintasan adalah perintah yang namanya Anda pilih sendiri dan menjalankan Claude Code sebagai satu akun. Seperti `ccam work`, pintasan meneruskan semua argumen ke Claude Code tanpa diubah, sehingga `ccwork --resume` dan `ccwork -p "hi"` bekerja seperti yang Anda harapkan. Pintasan adalah tautan ke ccam di `~/.local/bin`.

Pintasan bersifat opsional. Wizard `ccam add` menanyakan apakah perlu membuatnya, `ccam add <name> --shortcut <command>` membuatnya tanpa bertanya, dan `ccam list` menampilkan pintasan setiap akun, atau `-` jika akun tidak punya pintasan.

Nama memakai karakter a-z, 0-9, `.`, `-`, dan `_`. Nama tidak boleh berupa perintah ccam, pintasan akun lain, atau nama program lain di `PATH` Anda, juga bukan file di `~/.local/bin` yang tidak dibuat oleh ccam, seperti salinan atau tautan ccam milik Anda sendiri. Pintasan yang diubah secara manual di `config.toml` dan melanggar aturan ini, berupa `claude` atau `ccam`, atau diberikan ke dua akun membuat `config.toml` gagal dimuat, dengan menyebutkan akunnya.

`ccam rename` mempertahankan nama pintasan, dan `ccam remove` menghapus pintasan bersama akunnya. `ccam shims install` membuat ulang tautan semua akun, misalnya setelah Anda memindahkan ccam; file yang menghalangi satu pintasan tidak menghentikan pintasan lainnya.

## Pembagian

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Sebuah grup membagikan item di antara akun-akunnya; satu akun bisa masuk ke beberapa grup, tetapi setiap item hanya berasal dari satu grup. Lihat [berbagi dengan grup](/features/share).

## Menyalin penyiapan yang sudah ada

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` menyalin `settings.json`, `CLAUDE.md`, ekstensi, dan server MCP pengguna dari folder Claude Code ke sebuah grup, sehingga setiap akun yang membagikan item tersebut mendapatkannya. Folder itu hanya dibaca.

- `--settings`, `--memory`, `--extensions`, `--mcp`, dan `--history` hanya menyalin item yang Anda sebutkan. Tanpa opsi itu, di terminal sebuah daftar centang memungkinkan Anda memilih, dengan semua item kecuali riwayat percakapan sudah dicentang; di luar terminal, semuanya disalin kecuali riwayat percakapan.
- `--mcp` juga berfungsi saat `~/.claude` kosong atau tidak ada tetapi `~/.claude.json` memiliki server.
- `-g` atau `--group` memilih grup.
- `-y` atau `--yes` menyalin tanpa bertanya dan mempertahankan salinan grup jika ada file yang berbeda.

Info masuk MCP tidak disalin, karena info masuk yang disalin membuat pemegang lainnya ter-logout; rencana menyebutkan akun-akun yang perlu masuk dengan `/mcp`. Lihat [menyalin penyiapan yang sudah ada](/docs/sharing#copy-an-existing-setup).

## Server MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` memindahkan server yang hanya dimiliki akun itu ke grupnya, termasuk yang dimilikinya hanya untuk satu folder (`claude mcp add` biasa), yang kemudian dimuat di setiap folder. Server yang namanya sudah dipakai grup untuk definisi yang berbeda akan dilewati; tambahkan `--force` untuk mengganti salinan milik grup, yang dicadangkan ccam terlebih dahulu.

## Penggunaan dan perpindahan akun

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam sends no token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --continue --brief    # a new conversation from a summary of the last one
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agent "…"                  # from Claude Code: a subagent that runs as another account
ccam agent --resume <id> "…"    # follow that agent up
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
ccam ccusage claude daily       # Claude Code's tokens and costs, every account (ccusage)
ccam ccusage -a work blocks --active   # one account's current 5-hour block
```

- `next --auto`, `--fallback`, dan `ccam pick` memilih dengan cara yang sama: akun di bawah `auto_max_usage` (`config.toml`, bawaan 100) lebih dulu, menurut `auto_strategy` (`headroom`, `expiring`, atau `balance`), dengan menghitung jendela mingguan `--model`. Mereka tidak pernah memilih akun yang sudah 100% (jendela yang dilaporkan penuh oleh Claude Code dihitung 100% sampai di-reset) atau login kunci API, dan mereka tetap di dalam organisasi akun tempat mereka mulai (`--cross-org` atau `cross_org = true` untuk mencari lebih jauh).
- `--wait` menerima nilainya setelah `=` (`--wait=2h`); jika ditulis dengan spasi, opsi itu ditolak dengan petunjuk. Jika tidak ada yang bisa dipilih, semuanya keluar dengan kode 75, atau, dengan `--wait`, menunggu akun pertama di-reset.
- `ccam next --continue` membawa percakapan bahkan di antara akun yang tidak berbagi riwayat: Claude Code lalu menyalinnya ke akun baru. Dengan `--brief`, akun baru justru memulai percakapan baru, dari ringkasan percakapan terakhir yang ditulis ccam tanpa model ([`--brief`](/features/next#brief)).
- Tidak seperti `ccam <name>`, `ccam next` membaca opsinya sendiri, jadi opsi Claude Code ditaruh setelah `--`; tanpanya ccam berhenti dengan `unknown flag`.
- `ccam agent` menjalankan sebuah tugas dengan `claude -p` di folder ini sebagai akun lain selain akun tempat perintah itu dipanggil, dengan mode izin dan tingkat effort sesi pemanggil, lalu mencetak balasan terakhir Claude. Saat akun itu mencapai batasnya, percakapan berlanjut sebagai akun lain ([Subagen di akun lain](/features/agent)).
- Di dalam Claude Code, `/ccam:switch [account]` dari plugin meneruskan percakapan, dan `ccam hook` adalah yang dijalankan oleh hook-nya.
- Dengan `auto_switch = true` di `config.toml` (Linux dan macOS; nonaktif secara default), hook batas penggunaan juga mengakhiri sesi interaktif yang dimulai ccam dan melanjutkannya di akun lain yang punya ruang, tetapi hanya akun yang sejenis: kursi lain dari organisasi yang sama, atau paket pribadi lain ([`auto_switch`](/features/next#auto-switch)).

Lihat [Beralih ke akun berikutnya](/features/next), [Memilih akun untuk skrip](/features/pick), [Subagen di akun lain](/features/agent), dan [Sesi latar belakang semua akun](/features/agents).

## Pemeliharaan

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam uninstall       # remove every account, ccam's folders and ccam itself
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` memeriksa checksum SHA-256 rilis sebelum mengganti ccam, dan `--version v0.2.0` menginstal rilis tertentu. Salinan ccam yang diinstal dari paket `.deb` diperbarui dengan apt. Lihat [pembaruan](/docs/install#update).

`ccam uninstall` menampilkan apa saja yang dihapusnya dan bertanya terlebih dahulu: setiap akun seperti yang dilakukan `ccam remove`, perintah pintasan, folder ccam, dan perintah `ccam`. Folder yang diimpor tetap berada di tempatnya. Lihat [menghapus instalasi](/docs/install#uninstall).

## File

| Path | Isi |
| --- | --- |
| `~/.config/ccam/config.toml` | Akun, grup, dan pengaturan (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). Pengaturan ditaruh di bagian atas, di atas `[[group]]`, `[[profile]]`, atau `[[folder]]` pertama: pengaturan di bawah sebuah tabel menjadi milik tabel itu, sehingga ccam memperingatkannya di setiap peluncuran dan di `ccam doctor`, dan perintah yang menyimpan konfigurasi menolak sampai pengaturan itu dipindahkan. Komentar tidak dipertahankan saat ccam menyimpan file |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` milik setiap akun |
| `~/.local/share/ccam/groups/<group>/` | Pengaturan, server MCP, ekstensi, dan riwayat bersama |
| `~/.local/state/ccam/` | Cache penggunaan, catatan batas, kunci, catatan, dan status lainnya (lihat daftar di bawah tabel) |
| `~/.local/bin/<shortcut>` | Perintah pintasan seperti `ccwork`, berupa tautan ke ccam |

Folder status berisi:

- cache penggunaan dan catatan batas (`usage/<name>.json`, `usage/<name>.limit.json`);
- waktu terakhir digunakan;
- penanda sesi dan kunci (`locks/`);
- pintasan yang dibuat ccam (`shortcuts.json`);
- perubahan grup yang belum selesai (`group-change.json`);
- catatan yang ditinggalkan `/ccam:switch` serta hook batas penggunaan (`handoff/<account>/<session>.json`);
- kapan `ccam pick` atau `auto_switch` terakhir memilih setiap akun (`picks/<name>`; sebuah pilihan berlaku selama 2 menit);
- ringkasan dari `--brief`, disimpan 7 hari (`brief/<session>.md`).

Menetapkan `XDG_CONFIG_HOME`, `XDG_DATA_HOME`, atau `XDG_STATE_HOME` mengubah lokasi folder-folder ini. Di Windows, `~` adalah `%USERPROFILE%` dan pintasan berupa `<shortcut>.exe`.
