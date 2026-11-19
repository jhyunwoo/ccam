---
title: Komut başvurusu
metaTitle: ccam komut başvurusu
description: ccam add ve ccam <name> komutlarından ccam usage, ccam next ve ccam doctor komutlarına kadar tüm ccam komutları ve seçenekleri.
updated: 2026-11-19
order: 3
translation: machine
sourceHash: 28e1ad1dcbad
---

Bir komutun ne yapacağını, yapmadan görmek için komuttan önce `--dry-run` yazın, örneğin `ccam --dry-run work`. Bu, Claude Code başlatılırken (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`, `ccam ccusage`) ve `group`, `migrate`, `folder`, `shortcut`, `pick`, `update`, `uninstall` ile `doctor --fix` komutlarında çalışır. `list`, `usage`, `share` ve `ext` gibi yalnızca okuyan komutlar bunu kabul eder ve zaten hiçbir şeyi değiştirmez. Geri kalanlar, `group rename` de dahil, çalışmak yerine bunu reddeder.

<!-- commands -->

## Claude Code'u başlatma

`ccam <name|number> [arguments]`, Claude Code'u o hesapla başlatır. ccam, addan sonraki argümanları okumaz; bunları aynı sırayla `claude` komutuna iletir. Standart girdi, çıkış kodu ve sinyaller de tam olarak `claude` komutunun kendisinde olduğu gibi davranır. Addan sonra `--` yazmak isteğe bağlıdır.

Hesap adı olmadan argümanlarla çalıştırılan `ccam`, örneğin `ccam -p "summarise"`, geçerli klasör için ayarlanan hesabı, yoksa varsayılan hesabı kullanır. Bunları `ccam folder <name> [folder]` ve `ccam use <name>` ile ayarlayın.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Hesap adından önce gelen `--fallback`, `-p`/`--print` çalıştırmaları içindir. Hesap kullanım limitine ulaştığında ccam aynı komutu en çok payı kalan girişli hesapla yeniden çalıştırır ve yalnızca tamamlanan çalıştırmanın çıktısını yazdırır. [Klasöre göre hesaplar](/features/folder) ve [Limitte yeniden çalıştırma](/features/next#run-again-at-the-limit) sayfalarına bakın.

**Team ve Enterprise koltukları.** Otomatik hesap seçimi (`ccam next`, `--fallback`, `ccam pick`, kontrol panelindeki `n` tuşu ve açtığınızda [`auto_switch`](/features/next#auto-switch)), aynı Claude Team veya Enterprise kuruluşunun koltukları arasında da geçiş yapar ve bu, Anthropic'in koşullarını ihlal edebilir. Buna güvenmeden önce [SSS'yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) okuyun ve kuruluşunuza ve Anthropic'e danışın.

`ccwork` gibi bir kısayol, kendi hesabını aynı şekilde başlatır: `ccwork --resume`, `ccam work --resume` ile aynı işi yapar. Aşağıdaki [Kısayollar](/docs/commands#shortcuts) bölümüne bakın.

## Hesap ekleme ve kaldırma

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

Klasör belirtilmeden çalıştırılan `ccam import`, `~/.claude` ve `~/.claude-*` klasörlerini arar ve her biri için sorar; başarısız olan bir klasör `Skipped <path>: <reason>` olarak bildirilir ve diğerleri kaydedilir. Bir klasörün adı zaten alınmışsa o klasörü `ccam import <folder> --name <name>` ile tek başına içe aktarın. İçe aktarılan klasörler bulundukları yerde kaydedilir; ccam hiçbir zaman bir girişi taşımaz veya kopyalamaz.

`ccam login --expired`, girişi Claude Code tarafından silinmiş hesaplarda da oturum açar. `ccam remove`, yönetilen bir hesabın (`ccam import` ile kaydedilen değil, klasörünü ccam'in oluşturduğu bir hesabın) klasörü çöp klasörüne gitmeden önce o hesabın oturumunu kapatır; aynı girişi başka bir hesap tutuyorsa kapatmaz. İçe aktarılmış bir hesabı kaldırmak onu gruplarından çıkarır ve paylaşılan her öğenin bir kopyasını klasöründe bırakır. macOS'ta `ccam rename` klasörün adını korur, çünkü Anahtar Zinciri girişi klasörün yoluyla saklar.

`ccam add --api`, Claude girişi yerine API anahtarı kullanan bir hesap ekler; ccam anahtarı her başlatmada okur ve hiçbir zaman saklamaz. Bkz. [API anahtarlı hesaplar](/features/api).

`ccam rename`, `ccam remove` ve `ccam group` komutları, Claude Code bir hesapla çalışırken, nasıl başlatılmış olursa olsun (arka plan oturumları da dahil), o hesabı değiştirmeyi reddeder. Önce oturumdan çıkın. Yarım kalmış bir grup değişikliği varken de reddederler: ``the last group change stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` her hesabı yuvası, kısayolu, e-postası ve klasörüyle birlikte yazdırır. Çıktıdaki `groups` dizisi hesabın gruplarını listeler; `share` nesnesi ise her öğeyi (`mcp`, `settings`, `memory`, `extensions`, `history`) geldiği grupla veya `local` ile eşler.

## Kısayollar

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Kısayol, adını kendiniz verdiğiniz ve Claude Code'u tek bir hesapla başlatan bir komuttur. `ccam work` gibi her argümanı Claude Code'a değiştirmeden iletir; bu yüzden `ccwork --resume` ve `ccwork -p "hi"` beklediğiniz gibi çalışır. `~/.local/bin` içinde ccam'e giden bir bağlantıdır.

Kısayollar isteğe bağlıdır. `ccam add` sihirbazı bir kısayol oluşturulup oluşturulmayacağını sorar, `ccam add <name> --shortcut <command>` sormadan oluşturur ve `ccam list` her hesabın kısayolunu, yoksa `-` gösterir.

Adlarda a-z, 0-9, `.`, `-` ve `_` karakterleri kullanılır. Bir ad bir ccam komutu, başka bir hesabın kısayolu veya `PATH` içindeki başka bir programın adı olamaz; `~/.local/bin` içinde ccam'in oluşturmadığı bir dosyanın, örneğin ccam'in kendi kopyanızın veya bağlantınızın adı da olamaz. `config.toml` içinde elle düzenlenmiş ve bu kurallara uymayan, `claude` ya da `ccam` olan veya iki hesaba verilmiş bir kısayol, `config.toml` dosyasının yüklenememesine yol açar ve hata hesabın adını verir.

`ccam rename` kısayolun adını korur, `ccam remove` ise kısayolu hesapla birlikte siler. `ccam shims install`, örneğin ccam'i taşıdıktan sonra, tüm hesapların bağlantılarını yeniden oluşturur; bir kısayolun yolundaki bir dosya diğerlerini durdurmaz.

## Paylaşım

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default starts sharing history
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Bir grup, öğeleri kendi hesapları arasında paylaştırır; bir hesap birkaç grupta olabilir, ancak her öğe yalnızca tek bir gruptan gelir. [Gruplarla paylaşım](/features/share) sayfasına bakın.

## Mevcut bir kurulumu kopyalama

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate`, `settings.json`, `CLAUDE.md`, uzantıları ve kullanıcı MCP sunucularını bir Claude Code klasöründen bir gruba kopyalar; böylece bu öğeleri paylaşan her hesap onlara sahip olur. Klasör yalnızca okunur.

- `--settings`, `--memory`, `--extensions`, `--mcp` ve `--history` yalnızca belirttiğiniz öğeleri kopyalar. Bunlar olmadan, terminalde bir kontrol listesinden seçim yaparsınız ve konuşma geçmişi dışındaki her şey işaretlidir; diğer durumlarda konuşma geçmişi dışındaki her şey kopyalanır.
- `--mcp`, `~/.claude` boş ya da hiç yokken ama `~/.claude.json` içinde sunucular varken de çalışır.
- `-g` veya `--group` grubu seçer.
- `-y` veya `--yes` sormadan kopyalar ve bir dosya farklıysa grubun kopyasını korur.

MCP girişleri kopyalanmaz, çünkü kopyalanan bir giriş onu tutan diğerlerinin oturumunu kapatır; plan, `/mcp` ile oturum açacak hesapları adlarıyla belirtir. [Mevcut bir kurulumu kopyalama](/docs/sharing#copy-an-existing-setup) bölümüne bakın.

## MCP sunucuları

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>`, yalnızca o hesapta bulunan sunucuları hesabın grubuna taşır; buna hesabın yalnızca tek bir klasör için tuttuğu sunucular da (düz `claude mcp add`) dahildir ve bunlar sonra her klasörde yüklenir. Adı grupta zaten farklı bir tanım için kullanılan sunucular atlanır; grubun kopyasını değiştirmek için `--force` ekleyin, ccam önce o kopyayı yedekler.

## Kullanım ve geçiş

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

- `next --auto`, `--fallback` ve `ccam pick` aynı şekilde seçer: önce `auto_max_usage` (`config.toml`, varsayılan 100) altındaki hesaplar, `auto_strategy` (`headroom`, `expiring` veya `balance`) stratejisine göre, `--model` haftalık penceresi de sayılarak. %100'deki bir hesabı (Claude Code'un dolu bildirdiği bir pencere sıfırlanana kadar %100 sayılır) veya API anahtarıyla giriş yapan bir hesabı asla seçmezler ve başladıkları hesabın kuruluşu içinde kalırlar (daha ötesine bakmak için `--cross-org` veya `cross_org = true`).
- `--wait` değerini `=` işaretinden sonra alır (`--wait=2h`); boşlukla yazıldığında bir ipucuyla reddedilir. Hiçbiri seçilemediğinde 75 koduyla çıkarlar ya da `--wait` ile ilk hesabın sıfırlanmasını beklerler.
- `ccam next --continue` bir konuşmayı geçmişi paylaşmayan hesaplar arasında bile taşır: Claude Code o zaman konuşmayı yeni hesaba kopyalar. `--brief` ile yeni hesap bunun yerine, ccam'in bir model kullanmadan sonuncusundan yazdığı bir özetle yeni bir konuşma başlatır ([`--brief`](/features/next#brief)).
- `ccam <name>` komutunun aksine `ccam next` kendi seçeneklerini okur, bu yüzden Claude Code'un seçenekleri `--` sonrasına yazılır; o olmadan ccam `unknown flag` hatasıyla durur.
- `ccam agent`, bir görevi bu klasörde `claude -p` ile, çağrıldığı hesaptan başka bir hesapla, çağıran oturumun izin modu ve efor düzeyiyle çalıştırır ve Claude'un son yanıtını yazdırır. O hesap limitine ulaştığında konuşma başka bir hesapla sürer ([Diğer hesaplarda alt ajanlar](/features/agent)).
- Claude Code'un içinde eklentinin `/ccam:switch [account]` komutu bir konuşmayı devreder; `ccam hook` ise eklentinin hook'larının çalıştırdığı komuttur.
- `config.toml` içinde `auto_switch = true` olduğunda (Linux ve macOS; varsayılan olarak kapalı), limit hook'u ayrıca ccam'in başlattığı etkileşimli bir oturumu sonlandırır ve ona payı kalan başka bir hesapta devam eder, ancak yalnızca aynı türden bir hesapta: aynı kuruluşun başka bir koltuğu ya da başka bir kişisel plan ([`auto_switch`](/features/next#auto-switch)).

Bkz. [Sonraki hesaba geçme](/features/next), [Betikler için hesap seçme](/features/pick), [Diğer hesaplarda alt ajanlar](/features/agent) ve [Her hesabın arka plan oturumları](/features/agents).

## Bakım

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

`ccam update`, ccam'ı değiştirmeden önce sürümün SHA-256 sağlama toplamını kontrol eder; `--version v0.2.0` ise belirli bir sürümü kurar. `.deb` paketinden kurulan bir ccam kopyası bunun yerine apt ile güncellenir. [Güncelleme](/docs/install#update) bölümüne bakın.

`ccam uninstall`, neleri kaldıracağını listeler ve önce sorar: `ccam remove` komutunun yaptığı gibi her hesap, kısayol komutları, ccam'in klasörleri ve `ccam` komutu. İçe aktarılmış bir klasör olduğu yerde kalır. [Kaldırma](/docs/install#uninstall) bölümüne bakın.

## Dosyalar

| Yol | İçerik |
| --- | --- |
| `~/.config/ccam/config.toml` | Hesaplar, gruplar ve ayarlar (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). Ayarlar en üste, ilk `[[group]]`, `[[profile]]` veya `[[folder]]` satırının üstüne yazılır: bir tablonun altındaki ayar o tabloya aittir; bu yüzden ccam her başlatmada ve `ccam doctor` içinde bu konuda uyarır, yapılandırmayı kaydeden komutlar da ayar taşınana kadar çalışmayı reddeder. ccam dosyayı kaydederken yorumlar korunmaz |
| `~/.local/share/ccam/profiles/<name>/` | Her hesabın `CLAUDE_CONFIG_DIR` klasörü |
| `~/.local/share/ccam/groups/<group>/` | Paylaşılan ayarlar, MCP sunucuları, uzantılar ve geçmiş |
| `~/.local/state/ccam/` | Kullanım önbelleği, limit kayıtları, kilitler, notlar ve diğer durum bilgileri (tablonun altındaki listeye bakın) |
| `~/.local/bin/<shortcut>` | `ccwork` gibi kısayol komutları, ccam'e giden bağlantılar |

Durum klasöründe şunlar bulunur:

- kullanım önbelleği ve limit kayıtları (`usage/<name>.json`, `usage/<name>.limit.json`);
- son kullanım zamanları;
- oturum işaretleri ve kilitler (`locks/`);
- ccam'in oluşturduğu kısayollar (`shortcuts.json`);
- tamamlanmamış bir grup değişikliği (`group-change.json`);
- `/ccam:switch` ile limit hook'unun bıraktığı notlar (`handoff/<account>/<session>.json`);
- `ccam pick` veya `auto_switch` özelliğinin her hesabı en son ne zaman seçtiği (`picks/<name>`; bir seçim 2 dakika boyunca sayılır);
- 7 gün saklanan `--brief` özetleri (`brief/<session>.md`).

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` veya `XDG_STATE_HOME` ayarlandığında bu klasörlerin yeri değişir. Windows'ta `~`, `%USERPROFILE%` anlamına gelir ve kısayollar `<shortcut>.exe` olur.
