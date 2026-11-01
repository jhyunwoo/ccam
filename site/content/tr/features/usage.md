---
title: Kullanım limitleri
description: ccam usage tüm hesapların 5 saatlik ve haftalık limitlerini aynı anda gösterir, ccam statusline ise hesap adını ve kullanımı Claude Code içine yerleştirir.
updated: 2026-11-01
order: 7
translation: machine
sourceHash: 54b296f7da35
---

`ccam usage`, her hesap için her limitin ne kadarının kullanıldığını gösterir: 5 saatlik limit, haftalık limit ve her modelin kendi haftalık limiti (Fable gibi). Her limit bir çubuk ve sıfırlanmasına kalan süreyi alır. Uzun bir işe başlamadan önce veya Claude Code bir limitin yaklaştığını söylediğinde, hangi hesapta pay kaldığını görmek için kullanın. `ccam statusline`, Claude Code için durum satırı komutudur: hesap adını ve kullanımı Claude Code içinde gösterir ve sayıları `ccam usage` için kaydeder. Kullanımın ardındaki tokenlar ve tahmini maliyetler için `ccam ccusage`, [ccusage](https://ccusage.com) aracını her hesabın konuşmaları üzerinde çalıştırır ([aşağıda](#tokens-and-costs-ccam-ccusage)).

## Kullanım

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts that are not running
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

Her hesap; numarası, adı ve sayılarının ne kadar eski olduğunu gösteren bir satır alır; `(cached)` sayıların uç noktadan değil, kayıtlı bir dosyadan geldiği anlamına gelir. Altında her limit için bir satır, 20 hücrelik bir çubuk, yüzdeyi ve pencerenin sıfırlanmasına kalan süreyi gösterir: `5-hour`, `Weekly` (tüm modeller), ardından her modelin kendi haftalık penceresi, örneğin `Fable`. Hesabın sahip olmadığı bir pencere `-` gösterir. Yüzde %80'de sarıya, %100'de kırmızıya döner.

Claude Code bir pencere dolu olduğu için bir isteği reddettiyse o pencere, kayıtlı sayılar ne derse desin, sıfırlanana kadar %100 gösterir (bkz. [limit kaydı](#the-limit-record)). Terminalde her hesap hemen `fetching…` olarak görünür ve satırları hesap yanıt verdikçe çizilir.

Yalnızca birden fazla hesabınız olduğunda görünen son satır, [`ccam next --auto`](/features/next) komutunun seçeceği hesabı adlandırır. Aynı kuralları izler:

- `auto_max_usage` ve `auto_strategy`;
- çalışan oturumlar;
- `cross_org` ayarlanmadıkça yalnızca son kullanılan hesabın kuruluşundaki hesaplar;
- %100'deki her pencere atlanır; böylece haftası dolmuş bir hesap hiçbir zaman adlandırılmaz.

Varsayılan `headroom` stratejisiyle `→ Most headroom: <name>` yazar; `expiring` veya `balance` ile nedeni verir, örneğin `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Team ve Enterprise koltukları.** Otomatik hesap seçimi (`ccam next`, `--fallback`, `ccam pick`, kontrol panelindeki `n` tuşu ve açtığınızda [`auto_switch`](/features/next#auto-switch)), aynı Claude Team veya Enterprise kuruluşunun koltukları arasında da geçiş yapar ve bu, Anthropic'in koşullarını ihlal edebilir. Buna güvenmeden önce [SSS'yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) okuyun ve kuruluşunuza ve Anthropic'e danışın.

Giriş yapılmamış bir hesap `not logged in → ccam login <name>`, API anahtarlı bir hesap ise `no limits (API)` gösterir. Süresi dolmuş bir erişim tokenı, varsa kayıtlı sayıların üstünde ``access token expired; `ccam <name>` renews it when it starts`` gösterir: Claude Code onu bir sonraki başlatmada yeniler. Yenileme tokenının süresi dolduğu ya da sunucu reddettikten sonra Claude Code onu sildiği için tümüyle süresi dolmuş bir giriş `` login expired → `ccam login <name>` `` gösterir.

`ccam statusline install` sonrasında Claude Code, `[work] · Opus · 5h 24% · 7d 41%` gibi bir satır gösterir.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--json` | Raporları JSON olarak yazdırır: `name`, `slot`, `status`, `limits` (hesapta kullanım kredileri açıksa `credits` ile), `source`, `at`, `error` ve bir limit kaydı pencereleri %100'e ayarladığında `limited` ile `limited_windows`. Hiç hesap yoksa `[]` yazdırır | kapalı |
| `--watch[=N]` | Ctrl+C'ye basana kadar her `N` saniyede bir yeniden çizer. Tek başına `--watch` veya 0 ya da daha küçük bir değer 60 anlamına gelir. Değer `=` işaretinden sonra yazılır: `--watch 10` bir ipucuyla reddedilir | kapalı |
| `-v`, `--verbose` | Her hesabın altına bir `source:` satırı ekler: sayılarının nereden geldiği (kullanım uç noktası, Claude Code, durum satırı, kayıtlı bir dosya ya da Claude Code'un bildirdiği kullanım limiti). Bir limit kaydı sayıları geçersiz kıldığında ardından `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` gibi bir satır gelir | kapalı |
| `--no-cache` | 60 saniyelik önbelleği yok sayar ve uç noktaya yeniden sorar | kapalı |
| `--source live\|claude\|cache` | `live` uç noktaya hesabın token'ıyla sorar; `claude` kullanımı Claude Code'a sorgulatır (aşağıda); `cache` yalnızca kayıtlı sayıları gösterir | `usage_source`, aksi halde `live` |
| `--refresh` | Süresi dolmuş erişim tokenlarını yalnızca çalışmayan hesaplar için yeniler. Bir yenileme 30 saniyeye kadar sürebilir | kapalı |
| `--dry-run` | Kabul edilir: `ccam usage` kendi önbelleği dışında hiçbir şeyi değiştirmez. `--refresh` ile birlikte reddedilir: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | kapalı |

Başka herhangi bir `--source` değeri `--source must be live, claude or cache` ile başarısız olur. `ccam usage` argüman almaz.

### Claude Code'un kendisinden kullanım

`config.toml` dosyasının en üstünde (ilk `[[group]]`, `[[profile]]` veya `[[folder]]` satırının üstünde) `usage_source = "claude"` ile ya da `--source claude` ile ccam giriş token'ını hiç okumaz. Her hesap için hesabın klasöründe `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` başlatır ve Claude Code'un `get_usage` isteğini standart girdiden gönderir (Claude Code 2.1.270 veya sonrası). Claude Code kullanımı kendi girişiyle sorgular ve 5 saatlik, haftalık ve model pencereleriyle yanıt verir.

`--safe-mode`, hesabın hook'larının, MCP sunucularının ve eklentilerinin başlamasını engeller ve hiçbir oturum kaydedilmez. Her hesap bir Claude Code başlatması (yaklaşık bir saniye) gerektirir; bu yüzden bu kaynak `live` kaynağından yavaştır. Hesaplara aynı anda, her birine 20 saniyelik bir zaman aşımıyla sorulur ve 60 saniyelik önbellek her zamanki gibi geçerlidir.

Erişim tokenının süresi dolmuş ya da 5 dakika içinde dolacak bir girişe 50 saniye verilir, çünkü Claude Code onu başlarken yeniler ve yarıda kesilen bir yenileme girişi kaybettirebilir. ccam bir sorgudan vazgeçerse Claude Code'dan SIGTERM ile çıkmasını ister; böylece Claude Code başladığı bir yenilemeyi tamamlar (Windows'ta SIGTERM yoktur; bu yüzden orada süreç sonlandırılır).

Tümüyle süresi dolmuş bir giriş, Claude Code başlatılmadan `` login expired → `ccam login <name>` `` gösterir. Claude Code bu isteği deneysel olarak işaretler; bu yüzden yanıt vermeyi bırakırsa ccam kayıtlı sayıları gösterir ve `ccam doctor` bunu bildirir. Claude Code kendi sorgusunun başarısız olduğunu yanıtladığında ccam, limitsiz bir plan değil, `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)` ile birlikte kayıtlı sayıları gösterir.

`--dry-run` ile ccam Claude Code'u başlatmaz, çünkü Claude Code hesabın klasörüne yazar; bunun yerine kayıtlı sayıları gösterir. `usage_source`, `live` veya `claude` alır; `cache` yalnızca bir `--source` seçeneğidir.

`ccam statusline`, Claude Code'un durum satırı JSON'unu standart girdiden okur. Tek seçeneği `--wrap <command>` seçeneğidir: ilk çıktı satırının başına hesap adı eklenen mevcut bir durum satırı komutu. `ccam statusline install` seçenek almaz. İkisi de `--dry-run` seçeneğini `` --dry-run is not supported by `ccam statusline`; nothing was changed `` ile reddeder.

## Tokenlar ve maliyetler: `ccam ccusage`

`ccam usage`, bir hesabın her limitin ne kadarını kullandığını gösterir. [ccusage](https://ccusage.com) ise Claude Code'un konuşma günlüklerindeki tokenları ve tahmini maliyetleri sayar. Tek başına Claude Code'un günlüklerini yalnızca `~/.claude` ve `~/.config/claude` klasörlerinde (ya da `CLAUDE_CONFIG_DIR` değişkeninin belirttiği klasörlerde) arar; bu yüzden görebileceği tek ccam hesabı, `~/.claude` konumunda yerinde içe aktarılmış olandır. `ccam ccusage` onu, `CLAUDE_CONFIG_DIR` değişkenini konuşması olan her hesabın klasörüne ayarlayarak çalıştırır; böylece tek bir rapor hepsini kapsar:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam önce kendi seçeneklerini okur: `-a` (`--account`) hesap adlarını veya numaralarını alır; bunları virgülle ayırarak ya da seçeneği tekrarlayarak verebilirsiniz. `--dry-run` ise komutu çalıştırmak yerine yazdırır. Bunlar dışındaki ilk argümandan itibaren her şey, `bun x` ve `npx` için `--config` dışında (aşağıya bakın), ccusage'a yazdığınız gibi gider; `ccam ccusage -- --help` ccusage'ın kendi yardımını gösterir. ccam kurulu bir `ccusage` varsa onu, yoksa `bun x ccusage` komutunu, o da yoksa `npx ccusage@latest` komutunu çalıştırır; son ikisi ccusage'ı npm'den indirir. Bu ikisi bulunduğunuz klasörde değil, ccam'in kendi klasöründe (`~/.local/state/ccam/ccusage`) başlar; böylece bulunduğunuz klasördeki bir proje, kendi `node_modules/.bin` klasörüyle, `package.json` dosyasındaki bir `bin` ile ya da `.npmrc` dosyasındaki bir kayıt deposuyla (registry) bunların ne çalıştıracağını seçemez; `ccam --dry-run ccusage`, o klasöre geçen `cd` komutunu gösterir. Hiçbiri yoksa Node.js, Bun veya ccusage kurmanızı söyler. Linux ve macOS'ta ccusage, ccam'in yerini alır; Windows'ta ccam onu bekler ve onun çıkış koduyla çıkar. WSL'de ccam, bir Windows sürücüsündeki `ccusage`, `bun` veya `npx` programını atlar, çünkü o program hesap klasörleri olmadan bir Windows programı olarak çalışırdı.

- Rapor, klasörlerinde konuşma (bir `projects` klasörü) bulunan hesapları kapsar. Hesapların bir grup üzerinden paylaştığı bir konuşma bir kez sayılır.
- `-a` ile, geçmişini bir grup üzerinden paylaşan bir hesap tüm grubun konuşmalarını beraberinde getirir ve ccam bunların kimlere ait olduğunu söyler: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- `blocks` 5 saatlik pencereleri gösterir ve her hesabın kendi pencereleri vardır; bu yüzden bu rapor için yukarıdaki gibi tek bir hesap belirtin.
- ccusage, siz birini belirtmedikçe bulduğu diğer kodlama araçlarını da (Codex gibi) raporlar: `ccam ccusage claude daily` yalnızca Claude Code'u kapsar.
- ccam onu çalıştırırken ccusage, ayarlarını (`ccusage.json`) yalnızca `./.ccusage` konumundan ve hesap klasörlerinden okur. `--config` başka bir dosya belirtir: `ccam ccusage daily --config ~/.config/claude/ccusage.json`. `bun x` ve `npx` ccam'in kendi klasöründe başladığı için, o dosya varsa ve siz bir dosya belirtmediyseniz ccam onlara `./.ccusage/ccusage.json` dosyasını `--config` ile verir ve göreli bir `--config` yolunu tam yol olarak verir. Bu durumda ccusage, `--pi-path` gibi diğer göreli yolları ccam'in klasöründen okur; bu yüzden onları tam yol olarak verin.
- ccusage, yolunda virgül olan bir klasörü iki klasör olarak okurdu; bu yüzden ccam böyle bir klasörü reddeder ve hesabın adını verir. O hesabı `-a` ile dışarıda bırakın.
- `ccusage` bir ccam komutu olduğundan ccusage adlı bir hesap `ccam run ccusage` ile başlatılır; `ccam doctor` da bunu söyler.

## Nasıl çalışır

`ccam usage` komutunun gösterdiği sayılar iki yerden gelir: her hesabın kendi girişiyle çağrılan, Anthropic'in belgelenmemiş kullanım uç noktası `GET https://api.anthropic.com/api/oauth/usage` ve Claude Code'un durum satırı komutuna verdiği değerler. Bunların hepsi `internal/usage` içindedir.

1. ccam `config.toml` dosyasını yükler ve her hesap için bir `usage.Account` oluşturur (`cli.usageCmd`): klasörü `~/.local/share/ccam/profiles/<name>/`, oturum işaretleri (hesap klasöründeki Claude Code'un kendi kayıtlarıyla birlikte, çalışıp çalışmadığını söylerler) ve yerinde içe aktarılıp aktarılmadığı.
2. ccam her hesabı aynı anda, hesap başına bir goroutine ile sorgular (`Fetcher.FetchAll`). Hepsi tek bir `http.Client` paylaşır; böylece uç noktaya yapılan bağlantılar yeniden kullanılır. Bağlanma 2 sn sürebilir ve her hesaba toplam 3 sn verilir. İstemci yönlendirmeleri izlemez, çünkü bir yönlendirme tokenı başka bir yere yeniden gönderir; yanıtın da en fazla 1 MB'ını okur.
3. Her hesap bittikçe satırı yeniden çizilir. Başarısız olan veya zamanında yanıt vermeyen bir hesap kayıtlı sayılara geri döner; böylece yavaş bir hesap diğerlerini hiçbir zaman bekletmez.

Diyagramda her hesap kendi goroutine'inde sorgulanır ve zaman aşımına uğrayan hesap bunun yerine kayıtlı sayılarını gösterir:

<!-- motion:usage-fetch -->

`personal` ve `work` zamanında yanıt verir ve canlı sayılar gösterir. `side` 3 sn'den uzun sürer; bu yüzden ccam bunun yerine kayıtlı sayılarını gösterir (`fallback()`). ccam ardından her hesap için bir blok çizer (`usage.Table()`) ve son satırda adı geçen hesabı seçer (`usage.Best()`).

Tek bir hesap için ccam şu kontrolleri sırayla yapar (`Fetcher.Fetch`):

<!-- motion:usage-cache -->

1. `--source cache` ile doğrudan `fallback()` fonksiyonuna gider.
2. `--no-cache` ayarlı değilse, `~/.local/state/ccam/usage/<name>.json` içinde 60 saniyeden (`CacheTTL`) daha kısa süre önce ölçülmüş canlı bir sonuç, şu ana göre değerlendirilerek kullanılır: o zamandan beri sıfırlanmış bir pencere %0 gösterir ve saatin ilerisine tarihli bir dosya taze sayılmaz. `ccam usage` komutunu art arda iki kez çalıştırmak hesap başına bir isteğe mal olur.
3. ccam `.credentials.json` dosyasını (veya macOS Anahtar Zinciri'ni) değiştirmeden okur ve yalnızca `claudeAiOauth` anahtarının kodunu çözer (`readCredentials()`). Abonelik girişi yoksa `not logged in`, `.claude.json` bir API anahtarı tutuyorsa `no limits (API)` gösterilir. Yenileme tokenının süresi dolmuş ya da Claude Code'un sildiği bir giriş `` login expired → `ccam login <name>` `` olur ve sorgulanmaz.
4. Süresi dolmuş bir erişim tokenı kullanılmaz. `--refresh` olmadan hesap, kayıtlı sayılarını ``access token expired; `ccam <name>` renews it when it starts`` ile gösterir.
5. Aksi halde ccam uç noktayı çağırır. 401 veya 403 süresi dolmuş bir giriş sayılır (`login expired: authentication failed (401)`), diğer her durum kodu veya okunamayan bir gövde hata sayılır ve her ikisi de kayıtlı sayılara geri döner. Ctrl+C ile sonlandırılan bir sorgu `cancelled` der. 429, uç noktanın bu hesap için çok fazla istek aldığı anlamına gelir: ccam `<name>.pause` dosyasına yeniden sorma zamanını (`Retry-After` ya da 5 dakika, en fazla bir saat) yazar, kayıtlı sayıları `rate limited by the usage endpoint; asking again after 15:10` ile gösterir ve o zamana kadar, `--no-cache` ile bile, o hesap için uç noktayı çağırmaz. Hemen yeniden denemek yalnızca hız sınırını uzatır.
6. İyi bir yanıt, ölçüldüğü zamanla birlikte `<name>.json` dosyasına yazılır.

Kayıtlı sayılar için ccam en yeni kayıtlı sorguyu (`<name>.json` ya da Claude Code sorgularından gelen `<name>.claude.json`) alır ve `<name>.statusline.json` ondan sonra yazılmışsa durum satırının 5 saatlik ve haftalık pencerelerini onun üzerine yerleştirir (`fallback()`); durum satırında bulunmayan model pencereleri sorgunun kaydettiği gibi kalır. Kayıtlı sayıları göstermeden önce her pencerenin sıfırlanma zamanını kontrol eder: sayılar kaydedildikten sonra sıfırlanmış bir pencere %0 gösterir, çünkü eski yüzde artık geçerli değildir. Yukarıdaki `personal` hesabının sıfırlanma zamanı olmadan `0%` göstermesinin nedeni budur.

Uç noktanın iki yanıt biçimi vardır. Yenisinde `session` (5 saatlik), `weekly_all` ve `weekly_scoped` (modele özgü) girdileri içeren bir `limits` dizisi bulunur. Eskisinde `five_hour`, `seven_day`, `seven_day_opus` ve `seven_day_sonnet` nesneleri bulunur. ccam ikisini de okur (`ParseLive`): `limits` dizisi önceliklidir, eski nesneler ise onun kapsamadığını tamamlar. Hiç limit içermeyen bir gövde, değişmiş bir biçim olarak bildirilir ve `ccam doctor` bunu size söyler.

### Limit kaydı

Claude Code bir turu kullanım limitinde durdurduğunda [ccam eklentisinin](/docs/install#claude-code-plugin) hook'u reddedilen isteği konuşmadan okur: hangi pencerenin dolduğunu (5 saatlik, haftalık, Opus veya Sonnet) ve ne zaman sıfırlanacağını. Bunu `~/.local/state/ccam/usage/<name>.limit.json` dosyasına yazar. ccam kaydı, her kaynaktan bildiği diğer her şeyin üzerine uygular (`Fetch`): içindeki her pencere, kayıtlı sayılar ne kadar eski veya düşük olursa olsun, sıfırlanana kadar %100 sayılır; böylece `ccam next --auto`, `ccam pick`, `--fallback` ve [`auto_switch`](/features/next#auto-switch) hesabı o zamandan önce yeniden seçmez.

Başka hiçbir şey bilinmiyorsa hesabın gösterdiği yalnızca kayıttır (`-v` ile `source: the usage limit Claude Code reported`). Gözetimsiz oturumlar (`-p` çalıştırmaları ve arka plan oturumları) kayıt yazmaz; yöneticilerinin her an kaldırabileceği, bir kuruluşun koyduğu durdurma (örneğin kullanımı bittiğinde ya da bir grubun limiti $0 olduğunda) da yazmaz. Başka bir Claude hesabına yapılan `ccam login` kaydı siler; dosyayı silmek de onu temizler.

### Tokenları yenileme

Yenileme tokenları döner: birini yenilemek size yeni bir yenileme tokenı verir ve eskisini geçersiz kılar. Claude Code aynı hesapla çalışırken ccam girişi yenileseydi o oturumdan çıkış yapılırdı. Bu yüzden `ccam usage` tokenları okur ve siz `--refresh` vermedikçe hiçbir zaman yazmaz. O durumda ccam (`Fetcher.refresh`):

1. `ccam import` ile yerinde kaydedilmiş hesapları reddeder, çünkü düz `claude` o girişi ccam'in haberi olmadan kullanıyor olabilir;
2. macOS Anahtar Zinciri'nde saklanan girişleri reddeder;
3. iki ccam sürecinin aynı tokenı hiçbir zaman birlikte yenilememesi için bir kilit dosyası alır: `~/.local/state/ccam/locks/profile-<name>.refresh.lock`; kilit tutulurken hesap `token not renewed: another ccam process is using this login right now (a refresh, rename, remove or logout)` gösterir. `ccam rename`, `remove` ve `logout` da bu kilidi alır ve bir yenileme sürerken reddeder;
4. kilit altındayken hesabın bir oturumu, nasıl başlatılmış olursa olsun (ccam üzerinden, düz `claude`, bir IDE ya da bir arka plan oturumu ve onun arkasındaki, girişi kendisi yenileyen daemon), çalışıyorsa reddeder ve token başka bir süreç tarafından zaten yenilenmişse isteği atlar;
5. Claude Code hesap klasöründe kendi yenileme kilitlerini (`.oauth_refresh.lock`, `<folder>.lock`) tutarken hiçbir şey göndermez: Claude Code o anda o girişi yeniliyordur;
6. yenileme tokenını token uç noktasına gönderir ve yanıtı, Claude Code'un kendi yenilemelerine verdiği süre olan 30 saniyeye kadar bekler. Gönderildikten sonra bir istek, Ctrl+C ile bile asla kesilmez: sunucu yenileme tokenını zaten değiştirmiş olabilir ve yeni girişi yalnızca onun yanıtı taşır. Ctrl+C veya bir kapanma diğer sorguları bitirir ve ccam yanıt kaydedilince çıkar (Ctrl+C sonrasında 130 ile);
7. Claude Code'un yazma kilidi (`.storage-write.lock`) altında kaydeder ve `.credentials.json` dosyasının o andaki halinde yalnızca token alanlarını değiştirir: Claude Code'un bu arada yazdıkları (bir MCP girişi) kalır, Claude Code'un bu arada yenilediği bir giriş korunur ve bir oturum kapatma asla geri alınmaz. Yanıt yeni yenileme tokenının ömrünü verdiğinde bu, Claude Code'un yaptığı gibi `refreshTokenExpiresAt` olarak saklanır.

Gerçekleşmeyen bir yenileme nedenini hesap satırında veya kayıtlı sayılarının üstünde söyler: örneğin `token not renewed: a session is running` ya da bir istek gönderilmiş ama yanıtı kullanılamamış veya kaydedilememişse `login may be lost: …`; ardından Claude Code sizden giriş yapmanızı isterse `ccam login <name>` çalıştırma önerisiyle birlikte.

### Durum satırı toplayıcısı

Claude Code durum satırı komutunu çalıştırır ve ona boruyla bir JSON belgesi gönderir. ccam bundan model adını ve `rate_limits` değerini okur:

<!-- motion:usage-statusline -->

1. ccam hesabı, ccam'in Claude Code'u başlatırken ayarladığı `CCAM_PROFILE` değişkeninden, aksi halde `CLAUDE_CONFIG_DIR` değerini hesap klasörleriyle eşleştirerek bulur (`runStatusline`).
2. ccam `rate_limits.five_hour` ve `rate_limits.seven_day` değerlerini okur (`ParseStatusLine`). API anahtarlı hesaplarda veya oturumun ilk yanıtından önce JSON'da hız limitleri yoktur; bu durumda hiçbir şey kaydedilmez.
3. ccam `~/.local/state/ccam/usage/<name>.statusline.json` dosyasını yazar (`SaveStatusLine`). `ccam usage`, uç nokta yanıt vermediğinde bunu kullanır.
4. ccam `[name] · model · 5h N% · 7d N%` yazdırır. `--wrap` ile komutunuzu aynı girdiyle `sh -c` üzerinden çalıştırır ve çıktısını Claude Code'un bir durum satırını gösterdiği şekilde yazdırır: boş olmayan her satır, ilkinin önünde `[name]` ile. Başarısız olan bir komut yalnızca hesap adını gösterir. Windows'ta komutu Claude Code'un yaptığı gibi çalıştırır: kuruluysa Git Bash ile, aksi halde PowerShell ile.
5. 5 saatlik veya haftalık kullanım %90 ya da üzerindeyse ve başka bir hesap varsa satır `→ ccam next --auto --continue` ile biter, örneğin `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Sarmalanan komutun kendine ait bir süre sınırı yoktur: Claude Code bir durum satırı komutuna 10 dakika verir ve daha yeni bir güncelleme başladığında bir çalıştırmayı durdurur; böylece yavaş bir komut ccam olmadan nasıl görünecekse öyle görünür. Claude Code ccam'i sonlandırdığında komut da sona erer (Linux ve macOS'ta başlattığı her süreçle birlikte).

`ccam statusline install`, her hesabın `settings.json` dosyasına `statusLine` yazar ve nesnenin `padding`, `refreshInterval` ve `hideVimModeIndicator` gibi diğer anahtarlarını korur. Bir hesabın zaten bir durum satırı komutu varsa ccam onu `--wrap '<old command>'` olarak, Windows'ta ise iki kabukta da tırnak gerektirmeyen `--wrap-base64 <base64>` olarak korur. Durum satırı zaten ccam çalıştırıyorsa (başka bir yoldan ya da daha eski bir sürümün yazdığı şekilde tırnaklanmış olarak), ccam onu sarmalamak yerine değiştirir ve sarmaladığı komutu korur.

Windows'ta ccam'in yolu eğik çizgilerle ve ev klasörü için `~` ile yazılır (`~/.local/bin/ccam.exe statusline`); bunu hem Git Bash hem PowerShell anlar. Tırnak gerektiren bir yol, Claude Code'un kullandığı kabuk için tam olarak yazılır: Git Bash için `'…'`, PowerShell için `& '…'`. Bir grup üzerinden paylaşılan `settings.json` bağlantısı aracılığıyla bir kez yazılır; böylece gruptaki her hesap onu alır. Yeniden çalıştırmak, zaten ccam kullanan dosyalar için `(already set)` yazdırır.

## Güvenlik ve uç durumlar

- `ccam usage`, siz `--refresh` vermedikçe hiçbir zaman bir token dosyasını değiştirmez. Bir erişim tokenının süresi dolduğunda hesabı bir kez `ccam <name>` ile çalıştırın, Claude Code onu yeniler. Tümüyle süresi dolmuş bir giriş `ccam login <name>` gerektirir.
- `--refresh` çalışan hesaplara (Claude Code ccam olmadan başlatıldığında bile), içe aktarılmış hesaplara veya Anahtar Zinciri girişlerine hiçbir zaman dokunmaz. Reddedilen bir yenileme, kayıtlı sayıları ve nedeni (`token not renewed: …`) gösterir.
- Kullanım uç noktası herkese açık olarak belgelenmemiştir ve değişebilir. Değiştiğinde kayıtlı sayıları yine görürsünüz ve `ccam doctor` biçim değişikliğini bildirir.
- Kayıtlı sayı ve yanıt olmadığında hesap satırı `lookup failed:` ve nedeni gösterir, örneğin `--source cache` ile `no saved usage`.
- Önbellek dosyaları yalnızca yüzdeleri ve sıfırlanma zamanlarını tutar, hiçbir zaman token tutmaz.

## İlgili

- [Kullanım limitleri ve hesap değiştirme](/docs/usage-limits)
- [Sonraki hesaba geçme](/features/next)
- [Claude Code kullanım limitine ulaştığında ne yapmalı](/guides/usage-limit-workflow)
- [Komut başvurusu](/docs/commands)
