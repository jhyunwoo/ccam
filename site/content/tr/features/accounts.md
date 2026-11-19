---
title: Hesapları yönetme
description: Claude Code hesaplarını ccam add, login, list, rename, remove ve import ile ekleyin, giriş yapın, listeleyin, yeniden adlandırın, kaldırın ve içe aktarın.
updated: 2026-11-19
order: 2
translation: machine
sourceHash: "3bb475df8b0e"
---

Her ccam hesabı, kendi girişine sahip ayrı bir Claude Code yapılandırma klasörüdür. Hesap eklemek, hesaplarda oturum açmak, hesapları listelemek, yeniden adlandırmak ve kaldırmak için bu komutları kullanın. Bu komutlar zaten sahip olduğunuz klasörleri de kaydeder.

## Kullanım

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

Bir komut `<name|number>` aldığı her yerde hesap adını veya yuva numarasını (`ccam list` çıktısındaki `#` sütunu) verebilirsiniz.

`ccam use` yalnızca hangi hesabın varsayılan olduğunu değiştirir (`config.toml` içindeki `default`).

Terminalde seçeneksiz çalıştırılan `ccam add` bir sihirbaz açar. Sihirbaz şu adımlardan geçer:

1. siz vermediyseniz hesap adı;
2. hesabın nasıl giriş yapacağı: şimdi bir Claude girişi, bir API anahtarı ya da daha sonra bir Claude girişi;
3. `default` işaretli olarak hangi gruplara katılınacağı (henüz grup yoksa `default` grubunun neyi paylaşacağını gösterir);
4. `Create a shortcut command for work? [y/N]`;
5. onaylamanız için bir özet.

`Esc` bir adım geri gider. `--group`, `--no-group`, `--shortcut` veya `--no-login` seçeneklerinden herhangi biri ya da terminal olmadan çalıştırmak (örneğin bir borudan) sihirbazı atlar. Bu durumda bir ad zorunludur ve `--shortcut` bir ad vermedikçe kısayol oluşturulmaz.

`--group` verilmezse yeni bir hesap `default` grubuna katılır. İlk hesap bu grubu MCP sunucularını, ayarları, `CLAUDE.md` dosyasını ve uzantıları paylaşacak şekilde oluşturur; bkz. [gruplarla paylaşım](/features/share). Windows'ta dosyalar yalnızca Geliştirici Modu'nda bağlanabilir. Bu mod olmadan ilk hesap `default` grubunu yalnızca `mcp` ve `extensions` ile oluşturur ve bunu söyler (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`); ccam ayarları, belleği veya geçmişi paylaşan bir grupla verilen `--group` seçeneğini reddeder.

`ccam list` yuvayı, adı, kısayolu, e-postayı, grupları ve son kullanımı gösterir; varsayılan hesabın önünde `*` bulunur. E-posta, klasörün şu anda giriş yaptığı e-postadır (klasörün `.claude.json` dosyasından); ccam'in kaydettiğinden farklıysa kaydedilen e-posta ardından `(was …)` olarak gelir:

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

`--json` ile her hesap `slot`, `shortcut`, `name`, `email` (tablodaki gibi geçerli girişin e-postası), `default`, `dir`, `last_used` (hesap hiç başlatılmadıysa yer almaz), `groups` ve `share` alanlarına sahip bir nesnedir. `share` her öğeyi (`mcp`, `settings`, `memory`, `extensions`, `history`) geldiği grupla veya `"local"` ile eşler.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `ccam add --no-login` | Hesabı giriş yapmadan oluşturur; daha sonra `ccam login <name>` ile giriş yapın | kapalı |
| `ccam add --group <group>` | `default` yerine bu gruba katılır. Tekrarlanabilir; `default` yoksa oluşturulur | `default` |
| `ccam add --no-group` | Hiçbir gruba katılmaz: hesap hiçbir şey paylaşmaz | kapalı |
| `ccam add --shortcut <command>` | Bu kısayol komutunu `~/.local/bin` içinde de oluşturur | yok |
| `ccam add --api --key-from <source>` | Giriş yerine API anahtarı kullanan bir hesap; bkz. [API anahtarlı hesaplar](/features/api) | kapalı |
| `ccam list --json` | Hesapları JSON olarak yazdırır | kapalı |
| `ccam remove -y`, `--yes` | Sormadan kaldırır | kapalı |
| `ccam logout -y`, `--yes` | Sormadan oturumu kapatır | kapalı |
| `ccam import --name <name>` | Verdiğiniz klasör için hesap adı; klasör verilmeden reddedilir (`--name needs a folder: ccam import <folder> --name ops`) | klasör adından türetilir |
| `ccam import -y`, `--yes` | Bulunan her klasörü sormadan içe aktarır | kapalı |
| `--dry-run` | Hiçbir şeyi değiştirmeyen `ccam list` ve `ccam completion` bunu destekler. `add`, `login`, `logout`, `use`, `rename`, `remove` ve `import` ise `` --dry-run is not supported by `ccam add`; nothing was changed `` (komutun adıyla) mesajıyla reddeder | kapalı |

`ccam login` (`--expired` dışında), `ccam use` ve `ccam rename` komutlarının kendi seçenekleri yoktur. `ccam completion` `bash`, `zsh` veya `powershell` alır; başka bir değer `only bash, zsh and powershell are supported` hatasıyla başarısız olur.

Hesap adları `a-z`, `0-9`, `-` ve `_` karakterlerinden oluşan, harf veya rakamla başlayan 1 ile 32 karakter arası adlardır. Bir ad yalnızca rakamlardan oluşamaz, çünkü `ccam 2` 2 numaralı yuva anlamına gelir. Bir ad bir ccam komutu da olamaz (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `uninstall`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agent`, `agents`, `attach`, `ccusage`); böylece `ccam <name>` hiçbir zaman belirsiz olmaz. Liste, `internal/config` içindeki `config.Reserved` değişkenidir.

## Nasıl çalışır

Bir hesap eklemek sırasıyla klasörünü oluşturur, gruplarına katılır, yapılandırmayı kaydeder, giriş yapar ve kısayolu oluşturur. `ccam add work --shortcut ccwork` çalıştırdığınızda ccam şu adımları izler (`internal/cli/profiles.go` içindeki `addProfile`):

<!-- motion:accounts-add -->

1. ccam herhangi bir şey oluşturmadan önce adı (`config.ValidateName`) ve kısayolu (`shim.Check`) kontrol eder.
2. ccam `~/.local/share/ccam/profiles/work/` klasörünü `0700` moduyla oluşturur ve hesaba en küçük boş yuva numarasını verir; böylece kaldırılmış bir hesabın numarası yeniden kullanılır (`profile.Add`, `Config.NextSlot`). İlk hesap varsayılan olur.
3. Hesap gruplarına katılır (`--group` veya `--no-group` aksini söylemedikçe `default`; ilk hesap `default` grubunu oluşturur) ve ccam yeni, boş klasörü bu gruplara bağlar; örneğin `settings.json` dosyasını `groups/default/settings.json` dosyasına (`applyShare`, önce `share.Build` ve ardından `Apply` ile). Bu başarısız olursa yeni klasör kaldırılır, çünkü geride kalan bir klasör her yeniden denemeyi `folder already exists` ile engeller.
4. Hesap, `fsutil.AtomicWrite` aracılığıyla `~/.config/ccam/config.toml` dosyasına kaydedilir: önce geçici bir dosya, `fsync`, ardından `rename`; böylece bir çökme hiçbir zaman yarım yazılmış bir yapılandırma bırakmaz.
5. ccam, `CLAUDE_CONFIG_DIR` hesap klasörüne ayarlanmış olarak `claude auth login` komutunu, giriş sayfası dolu gelsin diye hesabın son e-postasıyla birlikte `--email` ile başlatır (`runLogin`). Giriş başarısız olursa hesap yerinde kalır; ccam hatayı yazdırır ve daha sonra `ccam login work` çalıştırmanızı söyler.
6. ccam `oauthAccount` değerini `profiles/work/.claude.json` dosyasından okur (`profile.ReadIdentity`): e-posta, hesap ve kuruluş kimlikleri, kuruluşun adı ve plan (`claude_max`, `claude_team` …). Bu dosya birkaç megabayt olabilir; bu yüzden ccam tüm dosyayı ayrıştırmak yerine yalnızca o anahtarı tarar (`TopLevelValue`). ccam bu ayrıntıları `config.toml` içinde `[profile.login]` altında saklar.
7. ccam `~/.local/bin/ccwork` bağlantısını ccam ikili dosyasına yönlendirir; ccam başladığında `argv[0]` değerini `config.toml` içinde arar. Bağlantı `~/.local/state/ccam/shortcuts.json` dosyasına kaydedilir (`shim.Install`). Bağlantı oluşturulamazsa hesap yine de eklenir, ccam nedenini söyler (`Could not create shortcuts: …`) ve özet `shortcut: ccwork (not created; see above)` gösterir.

Bir giriş dakikalar sürebilir ve bu arada başka bir ccam süreci yapılandırmayı kaydedebilir. Bu yüzden bekleyen komutlar `config.toml` dosyasını yeniden yükler, tek değişikliği (burada e-postayı) uygular ve kaydeder (`env.update`). `ccam remove` da aynı nedenle sorusundan sonra yeniden yükler.

`ccam login`, mevcut bir hesap için 5. ve 6. adımları çalıştırır. Claude Code e-postayı yalnızca giriş sayfası için bir ipucu olarak alır; bu yüzden öncekinden farklı bir Claude hesabına giriş yapmış olabilirsiniz. ccam yeni girişi `[profile.login]` ile karşılaştırır ve bunu söyler (`ccam: this login is other@company.com, but work was dev@company.com before.`). Yeni girişi korur ve bunu nasıl geri alacağınızı söyler: ``To undo it, run `ccam login work` again and sign in as dev@company.com.``

Başka bir Claude hesabına yapılan giriş, hesabın [limit kaydını](/features/usage#the-limit-record) da siler, çünkü o limitler diğer hesabındı. Başka bir ccam hesabı aynı Claude hesabına giriş yapmışsa ccam, ikisinin kullanım limitlerini paylaştığı konusunda uyarır.

`ccam login --expired`, hiç giriş yapmamış ya da girişinin süresi dolmuş her hesapta sırayla giriş yapar: yenileme token'ının süresi dolmuştur ya da sunucu reddettikten sonra Claude Code girişi silmiştir. İlk tarayıcı açılmadan önce `config.toml` dosyasının kaydedilebildiğini denetler.

`ccam logout`, hesapta `claude auth logout` komutunu çalıştırır ve kayıtlı e-postayı temizler; giriş kaydı kalır, böylece bir sonraki giriş onunla karşılaştırılır. Claude Code girişi iptal eder ve hesap için sakladıklarını, MCP sunucusu girişleri de dahil, siler. Bu yüzden ccam şu durumlarda sizi uyarır ve önce sorar (`--yes` soruyu atlar):

- giriş MCP girişleri tutuyorsa;
- hesap çalışıyorsa (oturumlarından da çıkış yapılır);
- başka bir hesabın klasörü aynı girişi tutuyorsa (kopyalanmış bir klasör); oturumu kapatmak onu da iptal eder.

Örneğin:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

Başka bir Claude hesabı olarak giriş yapmak için tek başına `ccam login <name>` yeterlidir ve MCP girişlerini korur.

Hesaplar birden fazla kuruluşa aitse `ccam list` bir `Organization` sütunu ekler (kişisel planlar `personal` gösterir) ve `ccam list --json` her hesap için `plan` ve `org` alanlarını içerir.

### Uzak bir makinede oturum açma

SSH üzerinden `claude auth login` komutunun tarayıcı tarafı bozulabilir. Claude Code, oturum açmayı uzak makinedeki bir geri çağırma sunucusunda (`http://localhost:<port>/callback`) bekler; ancak oturum açmayı tamamlayan tarayıcı sizin bilgisayarınızda çalışır. Orada bu adres hiçbir yere gitmez ve tarayıcı `ERR_CONNECTION_REFUSED` gibi bir hata gösterir. Adres yine de tek kullanımlık kodu taşıdığı için ccam onu iletebilir:

<!-- motion:accounts-relay -->

ccam terminalinizi Claude Code'a vermez (`runLogin`). Claude Code'un standart girdisini bir boruya bağlar ve yazdığınız her satırı bu boruya kopyalar (`relayLogin`); böylece Claude Code'un `Paste code here` istemine yapıştırılan bir kod yine ulaşır. Bir satır geri çağırma adresi olduğunda ccam onu yalnızca ana makine `localhost`, `127.0.0.1` veya `[::1]` ise, bir port verilmişse, yol `/callback` ise ve hem `code` hem `state` mevcutsa kabul eder; `http://` yazılmayabilir (`callbackURL`). Ardından ccam, tarayıcının yapacağı gibi, bu adresi sunucuda 10 saniyelik bir zaman aşımıyla ve yönlendirmeleri izlemeden ister (`deliverCallback`). İstek, bekleyen `claude auth login` komutunun geri çağırma sunucusuna ulaşır (diyagramdaki kesikli çizgi). Claude Code kodu alır ve girişi kaydeder; Linux'ta hesap klasöründeki `.credentials.json` dosyasına. Yapıştırılan bir satır ccam'in başka bir ana makineyle bağlantı kurmasına asla yol açamaz.

ccam, giriş başlamadan önce şunu yazdırır:

```text
On a remote machine? If your browser ends up on an http://localhost:…/callback page that does not load, paste that address here.
```

Başarılı olduğunda `→ Passed the sign-in to Claude Code on this machine.` yazdırır.

Oturum açma bir terminalde girdi beklerken ccam terminalin yankısını ve satır modunu kapatır ve yazdığınız her karakter için bir `*` gösterir; böylece yapıştırılan bir kod veya geri çağırma adresi ekranda görünmez. ccam satır düzenlemeyi kendisi yapar (`readLine`): `Enter` satırı gönderir, `Backspace` bir karakter siler ve boş bir satırda `Ctrl+D` girdiyi bitirir. Terminal, oturum açma bittiğinde ve ayrıca ccam'i yine sonlandıran `Ctrl+C`, `SIGTERM` veya `SIGHUP` durumunda geri yüklenir. Standart girdi ve standart hata bir terminale bağlı değilse satırlar geldikleri gibi kopyalanır.

### Yeniden adlandırma, kaldırma ve içe aktarma

`ccam rename` ve `ccam remove`, hesabın klasörü ve ccam'in hesabın adı altında başka yerlerde tuttukları üzerinde işlem yapar. ccam bu yolları `nameState` içinde listeler:

- `~/.local/state/ccam/last-used/<name>`;
- `~/.local/state/ccam/mcp/<name>`;
- kullanım önbellekleri `~/.local/state/ccam/usage/<name>.json` ve `<name>.statusline.json`;
- limit kaydı `<name>.limit.json`;
- `~/.local/state/ccam/handoff/<name>/` içindeki oturum notları.

`ccam remove work` bunları şu sırayla işler:

<!-- motion:accounts-remove -->

1. `ccam remove`, hesap çalışıyorsa önce reddeder, ardından sorar ve siz yanıtladıktan sonra yeniden denetler (bu, soru açıkken başlatılan bir oturumu yakalar).
2. Herhangi bir şeyi değiştirmeden önce `config.toml` dosyasının kaydedilebildiğini denetler.
3. Hesap uzantıları paylaşıyorsa ccam grubun `installed_plugins.json` ve `known_marketplaces.json` dosyalarını yeniden yazar (`RewritePluginPaths`), çünkü Claude Code, bu hesabın klasöründen geçebilecek mutlak eklenti yollarını saklar.
4. Ardından yönetilen bir hesabın (`ccam import` ile kaydedilen değil, klasörünü ccam'in oluşturduğu bir hesabın) oturumu kapatılır: `claude auth logout` onun klasöründe çalışır; böylece çöp klasöründe (ya da macOS'ta, o adla sonra eklenen bir hesabın devralacağı, klasörün adını taşıyan bir Anahtar Zinciri öğesinde) çalışan bir giriş beklemez. API anahtarlı bir hesapta ve başka bir hesabın klasörü aynı girişi tuttuğunda bu adım atlanır, çünkü oturumu kapatmak onu da iptal ederdi. Başarısız olan bir oturum kapatma bildirilir ve kaldırma devam eder.
5. ccam klasörü `~/.local/share/ccam/trash/<name>-<time>/` konumuna taşır (`profile.Trash`), hesap `config.toml` dosyasından çıkarılır, aynı adlı sonraki bir hesap bunları devralmasın diye `nameState` yolları silinir ve kısayol bağlantısı kaldırılır: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`.

`config.toml` sonunda yine de kaydedilemezse klasör çöp kutusundan geri gelir. Kaldırılmış bir hesabı kurtarmak için klasörü çöp kutusundan geri taşıyın ve yeniden giriş yapın.

`ccam rename` klasörü yeniden adlandırır ve `plugins/` içindeki eklenti yollarını eski klasörden yenisine göre yeniden yazar. `nameState` yollarını yeni ada taşır ve hesap varsayılansa varsayılan olarak kalmasını sağlar. Oturum notlarını da günceller, böylece `ccam next --continue` onları yine bulur. Kısayol komutu adını korur. Eklenti yollarının yeniden yazılması başarısız olursa ya da `config.toml` kaydedilemezse klasör eski adına geri döndürülür.

macOS'ta ve Claude Code girişleri Kimlik Bilgileri Yöneticisi'nde sakladığında Windows'ta klasörün kendisi eski adını korur, çünkü giriş ve MCP girişleri klasörün yoluyla saklanır ve taşınan bir klasörün oturumu kapanırdı. Yapılandırma bu durumda klasörü yeni adın yanında `folder` olarak kaydeder ve ccam `its folder stays …: the macOS Keychain keeps the login under the folder's path` der.

`ccam import` bir klasörü bulunduğu yerde, `config.toml` içinde `dir` olarak kaydeder. Onu hiçbir zaman taşımaz veya kopyalamaz: kopyalanmış bir yenileme tokenı, kopyaların birbirinin oturumunu kapatmasına neden olur.

Klasör argümanı olmadan kaydedilmemiş `~/.claude` ve `~/.claude-*` klasörlerini arar ve her biri için sorar. Adlar klasörden gelir (`~/.claude` için `main`, `~/.claude-work` için `work`, `~/.claude-2` için `acct2`). ccam bir ad türetemezse ya da ad alınmışsa klasörü, çalıştırılacak `--name` komutuyla birlikte listeler. Bulunan bir klasör başarısız olursa `Skipped <path>: <reason>` olarak bildirilir ve diğerleri içe aktarılır.

ccam önce tüm sorularını sorar, ardından `config.toml` dosyasını yeniden yükler ve klasörleri kaydeder; böylece diğer ccam komutlarının bu arada kaydettiği değişiklikler korunur. `✓` satırları kayıttan sonra gelir.

`CLAUDE_CONFIG_DIR` olmadan Claude Code, `.claude.json` dosyasını ev klasörünüzde `~/.claude` klasörünün yanında tutar; bu yüzden `~/.claude` klasörünü içe aktardığınızda ve kendine ait bir `.claude.json` dosyası yoksa ccam, `~/.claude/.claude.json` bağlantısını `~/.claude.json` dosyasına yönlendirerek oluşturur. Böylece düz `claude` ve `ccam main` aynı hesabı görür.

`~/.claude` klasörünün `~/.claude.json` dosyasının yanında kendine ait bir `.claude.json` dosyası varsa ccam ikisini de bırakır ve bunu söyler. `ccam main` içtekini, düz `claude` ise dıştakini kullanır ve giriş ayrıntılarıyla MCP sunucularını ayrı tutar. Tek bir dosya tutmak için içteki dosyayı bir bağlantıyla değiştirin ([doctor](/features/doctor) da bunu bildirir).

Klasör kayıtlı bir hesapla aynı girişi tutuyorsa (kopyalanmış bir klasör) ccam bunu söyler: girişi ilk yenileyen diğerinin oturumunu kapatır, bu yüzden onu kullanmadan önce `ccam login <name>` çalıştırın.

İçe aktarılan hesaplar siz onları bir gruba ekleyene kadar (`ccam group join <group> <name>`) hiçbir grupta değildir ve kısayol oluşturulmaz.

`ccam completion bash|zsh`, hesap adlarını da tamamlayan bir tamamlama betiği yazdırır.

## Güvenlik ve uç durumlar

- `ccam rename` ve `ccam remove`, hesap çalışırken ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)`` mesajıyla reddeder (`refuseRunning`), çünkü oturum klasöre yazmaya devam eder. Denetim, nasıl başlatılmış olursa olsun Claude Code'u kapsar: ccam üzerinden, düz `claude`, bir IDE, arka plan oturumları ve onların daemon'ı.
- `ccam rename`, `ccam remove` ve `ccam logout`, başka bir ccam süreci o girişi yenilerken veya değiştirirken reddeder: ``another ccam process is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- Başka bir ccam sürecinde bir grup değişikliği sürerken `ccam rename` ve `ccam remove` `another ccam process is changing the groups right now; try again when it is done` mesajıyla, yarıda kalmış bir değişiklikten sonra ise ``the last group change stopped halfway; finish it with `ccam doctor --fix` first`` mesajıyla reddeder.
- İçe aktarılmış bir hesapta `ccam remove`, düz `claude` için klasörü ve girişini korur. Hesap gruplardaysa önce `ccam group leave` gibi onlardan ayrılır: paylaşılan her öğe klasörde grubunkinin bir kopyası olur ve grubun MCP sunucuları onun `.claude.json` dosyasına girer. Ardından kaydı silinir: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. ccam'in onun `.claude.json` dosyası için oluşturduğu yedekler `<folder>/.ccam-backup/mcp-<time>` içine taşınır; taşınamazlarsa `~/.local/state/ccam/mcp/<name>/backups` içinde kalırlar ve mesaj bunu söyler. `~/.claude` klasörünüz çöp kutusuna taşınmaz.
- İki hesap aynı klasöre sahip olamaz: bir klasörü iki kez belirten `config.toml` yüklenemez.
- `ccam add`, klasörü `profiles/` içinde zaten bulunan bir adı reddeder ve `ccam import` komutunu önerir. macOS'ta, Anahtar Zinciri yeni klasörün yolu için ccam'in eski bir sürümünden kalma bir giriş veya MCP girişleri hâlâ tutuyorsa da durur ve bunları kaldıran `security delete-generic-password` komutunu belirtir. `config.toml` kaydedilemezse yeni klasör yeniden kaldırılır.
- `ccam import`, zaten kayıtlı bir klasörü, klasör olmayan bir yolu ve kullanılmakta olan bir adı reddeder. ccam, klasör `~/.claude` klasörünün bir kopyasıysa (aynı giriş), `~/.claude` kendisi kayıtlı olmasa bile uyarır.
- `ccam logout` hesap çalışıyorsa önce sorar; `ccam login` denetlemez.
- Aktarma Claude Code'a ulaşamazsa, örneğin giriş sona erdikten sonra, nedeniyle birlikte `ccam: could not pass the sign-in to Claude Code:` yazdırır. Bir HTTP hatası genellikle kodun süresinin dolduğu anlamına gelir; girişi yeniden başlatın.

## İlgili

- [Claude Code'u başlatma](/features/launch)
- [Kısayol komutları](/features/shortcuts)
- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [Tek bilgisayarda birden fazla Claude Code hesabı kullanma](/guides/multiple-claude-code-accounts)
