---
title: Kontrol paneli
description: Terminalde ccam; hesapların, grupların, klasör kurallarının, uzantıların ve denetimlerin tam ekran kontrol panelini açar; her değişiklik önce planını gösterir.
updated: 2026-10-30
order: 15
translation: machine
sourceHash: d65f1b5bb0c8
---

Argümansız `ccam`, standart girdi ve çıktı bir terminal olduğunda tam ekran bir kontrol paneli açar; aksi halde `ccam list` çıktısını yazdırır. Kontrol paneli her hesabı kullanımı ve gruplarıyla birlikte gösterir ve komutların yaptığını yapar: başlatma, ekleme, giriş yapma, grupları değiştirme, klasör kuralları belirleme, arka plan oturumlarını listeleme, denetleme ve onarma. Başlatma ya da bir arka plan oturumuna bağlanma paneli kapatır; bir hesap eklemek veya giriş yapmak terminali kısa bir süre devreder, ardından kontrol paneline geri döner.

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

Sekme adlarının altındaki satır, geçerli sekmenin ne işe yaradığını tek bir cümleyle söyler. Daha geniş bir terminalde Accounts sekmesi her hesabın e-postasını da gösterir. Her sekmede, sekme satırının sonu [`auto_switch`](/features/next#auto-switch) ayarının açık olup olmadığını gösterir (Linux ve macOS; Windows'ta yok sayılır). Açıkken Health sekmesi, bu ayarın hâlâ neye ihtiyaç duyduğunu listeler: Claude Code 2.1.214 veya sonrası ve Claude aboneliğiyle giriş yapan her hesapta açık bir ccam eklentisi.

Kontrol paneli her açıldığında arka planda GitHub'a daha yeni bir ccam sürümü olup olmadığını sorar; ilk kare bunu asla beklemez ve çevrimdışıyken hiçbir şey gösterilmez. Yeni bir sürüm varsa başlığın altındaki bir satır onu önerir:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` önce sorar, sonra [`ccam update`](/features/update) gibi günceller: sürüm SHA-256 sağlama toplamlarıyla doğrulanır ve tek adımda ccam'in yerini alır. İçinde bulunduğunuz kontrol paneli hâlâ eski sürümdür; bu yüzden çıkın ve `ccam` komutunu yeniden başlatın. Güncelleme burada yapılamadığında, örneğin `.deb` paketinden kurulmuş bir ccam kopyası için, kontrol paneli `ccam update` ile aynı öneriyi gösterir ve teklifi korur.

İlk kare yalnızca `config.toml` dosyasından, 100 ms'nin çok altında çizilir; kullanım her hesap yanıt verdikçe dolar. Her limitin kendi çubuğu vardır: 5 saatlik pencere, tüm modeller için haftalık pencere ve Fable'ın kendi haftalık penceresi; `-` hesabın böyle bir penceresi olmadığı anlamına gelir. Bir yüzde %80'de sarıya, %100'de kırmızıya döner. Bir sorgu başarısız olursa önbellekteki değerler soluk gösterilir. Claude Code'un dolu bildirdiği bir pencere sıfırlanana kadar %100 gösterir. Süresi dolmuş bir erişim token'ı soluk bir `⚠ token expired` gösterir (Claude Code onu bir sonraki başlatmada yeniler); `⚠ login expired` ise hesabın `ccam login` gerektirdiği anlamına gelir.

## Sekmeler

| Sekme | Gösterdiği | Tuşlar |
| --- | --- | --- |
| Accounts | numara, ad, e-posta, 5 saatlik, haftalık ve Fable çubukları, son kullanım, gruplar; `✔` varsayılan, `▸` bu klasörün hesabı, `•` çalışıyor (ccam olmadan başlatılan Claude Code da: düz `claude`, bir IDE, arka plan oturumları ve onların daemon'ı) | aşağıda |
| Agents | [`ccam agents`](/features/agents) komutunun listeledikleri: her hesabın arka plan oturumları ve açık oturumları, durumları, klasörleri ve başlama zamanlarıyla (aşağıda) | `enter` bağlan (`ccam attach` gibi), `x` durdur, `r` yeniden yükle |
| Groups | her grup, neyi paylaştığı ve hesapları, ayrıca hiçbir grupta olmayan hesaplar | `enter` düzenle, `n` yeni, `e` yeniden adlandır, `x` kaldır |
| Folders | [`ccam folder`](/features/folder) kuralları ve burada geçerli olan kural; klasöre sembolik bir bağlantı üzerinden ulaştığınızda da (yeni kurallar bağlantılar çözülmüş olarak saklanır) | `n` burada bir hesap kullan, `enter` değiştir, `x` kaldır |
| Extensions | [`ccam ext`](/features/ext) komutunun gösterdikleri, her uzantı bir kez ve geldiği gruplarla; `off:` bir eklentinin açık olmadığı hesaplar, hesabın yalnızca bu klasör için geçerli olan kendi sunucuları, `This folder's project` ve hiçbir hesaba ait olmayan sorunlar | `r` yeniden yükle |
| Health | [`ccam doctor`](/features/doctor) komutunun bildirdikleri, hesap başına bir kutu olarak (aşağıda) | `f` seçili kutuyu onarır, `F` hepsini onarır (`doctor --fix`), `↑`/`↓` kutular arasında geçer |

`←`/`→` veya `Tab` sekmeler arasında geçer, `?` geçerli sekmenin tuşlarını listeler, `r` kısa önbelleği atlayarak kullanımı yeniler ve Agents, Extensions ve Health sekmelerini yeniden yükler, `q` veya `Esc` ise çıkar.

## Accounts sekmesindeki tuşlar

| Tuş | Ne yapar |
| --- | --- |
| `Enter` | Claude Code'u seçili hesapla başlatır |
| `1`–`9` | O numaralı hesaba gider; 12 için `1` `2` tuşlarına hızlıca basın. Yalnızca `Enter` başlatır |
| `↑`/`↓`, `k`/`j` | Hareket eder |
| `n` | [`ccam next --auto`](/features/next) komutunun seçeceği hesabı başlatır |
| `o` | Seçeneklerle başlatır: `--continue`, `--resume`, başka argümanlar ya da bir kuru çalıştırma |
| `u` | Kullanım ayrıntıları: sıfırlanma zamanlarıyla 5 saatlik, haftalık ve model pencereleri ve sayıların nereden geldiği |
| `a` | Bir hesap ekler (aşağıdaki sihirbaz) |
| `l` | Yeniden giriş yapar |
| `s` | Hesabın gruplarını seçer |
| `d` | Hesabı varsayılan hesap yapar |
| `f` | Hesabı geçerli klasörde kullanır |
| `e` | Hesabı yeniden adlandırır |
| `c` | Hesabın kısayol komutunu ayarlar veya kaldırır |
| `x` | Hesabı kaldırır (önce, Hayır seçili olarak sorar) |

**Team ve Enterprise koltukları.** Otomatik hesap seçimi (`ccam next`, `--fallback`, `ccam pick`, kontrol panelindeki `n` tuşu ve açtığınızda [`auto_switch`](/features/next#auto-switch)), aynı Claude Team veya Enterprise kuruluşunun koltukları arasında da geçiş yapar ve bu, Anthropic'in koşullarını ihlal edebilir. Buna güvenmeden önce [SSS'yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) okuyun ve kuruluşunuza ve Anthropic'e danışın.

## Health sekmesi

Health sekmesi, açtığınızda [`ccam doctor`](/features/doctor) kontrollerini çalıştırır ve aynı görünümü çizer: her hesap için bir işaret içeren ilk satır, ardından hesap başına bir kutu. Bir kutudaki her satır kontrol edilen bir alandır; önce `Folder`, `Login` ve `Share links` gelir ve `✔` sorun yok, `ℹ` bilgi, `⚠` uyarı veya `✘` hata ile işaretlenir; böylece hangi hesabın nerede sorunu olduğunu hemen görürsünüz. Hiçbir hesaba ait olmayan bulgular bir `General` kutusuna girer.

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

`↑`/`↓` kutudan kutuya geçer ve seçili kutunun kenarlığı vurgu rengindedir. `f` yalnızca seçili kutuyu, `ccam doctor --fix` komutunun yapacağı şekilde onarır: bir hesabın kutusunda yalnızca o hesabı, `General` kutusunda ise hiçbir hesaba ait olmayan bulguları (eklenti yolları ve kısayollar) onarır; diğer tüm hesaplar yalnızca yeniden kontrol edilir. Önce sorar (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``) ve onarılacak bir şey olmayan bir kutu için yalnızca `nothing to repair in <name>` der. `F` her şeyi bir kerede, `ccam doctor --fix` gibi onarır; `r` yeniden kontrol eder.

## Agents sekmesi

Agents sekmesi, [`ccam agents`](/features/agents) komutunun yazdırdıklarını gösterir: her hesabın arka plan oturumlarını ve açık oturumlarını, her birinin ait olduğu hesapla birlikte. Sekmeyi ilk açtığınızda, `r` tuşunda, bir durdurmadan sonra ve yapılandırmada bir değişiklikten sonra her hesabın Claude Code'una aynı anda sorar (`claude agents --json`, her biri 15 saniyelik bir zaman aşımıyla) ve hepsi yanıt verene kadar `Asking every account's Claude Code for its sessions…` yazar.

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

`State`, `ccam agents` komutunda olduğu gibi bir arka plan oturumunun durumu ve sürecinin ne yaptığıdır; etkileşimli bir oturum için ise `open` yazar. Bir hesap sorgulanamazsa listenin altındaki bir `✘` satırı Claude Code'un söylediğini gösterir; ajan görünümü kapalı olan yukarıdaki `lab` gibi. Biten arka plan oturumları dışarıda bırakılır (`ccam agents --all` onları listeler) ve her hesap yanıt verip hiçbirinin oturumu yoksa sekme ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` yazar. Dar bir terminalde sekme önce `Started`, sonra `Name`, sonra `Folder` sütununu kaldırır ve uzun klasör yollarını soldan kısaltır.

Bir arka plan oturumunda `Enter`, kontrol panelini kapatır ve oturumu bu terminalde, `ccam attach <id>` komutunun yaptığı gibi kendi hesabıyla açar; kontrol paneli geri dönmez (Claude Code içinde `←` ajan görünümüne, `Ctrl+Z` kabuğa götürür). Açık bir oturumun kimliği yoktur; bu yüzden `Enter` ve `x` yalnızca bir arka plan oturumu gerektiğini söyler. `x` önce `No, keep it` seçili olarak sorar: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` Evet derseniz ccam o hesapla `claude stop <id>` komutunu çalıştırır, ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` ya da neden durduramadığını yazar ve listeyi yeniden okur.

## Grupları değiştirme

Groups sekmesinde `Enter` düzenleyiciyi açar: solda öğeler, sağda hesaplar; `Tab` aralarında geçer, `Space` açıp kapatır. Hesaplardan birinin zaten başka bir gruptan aldığı bir öğe, nedeniyle birlikte gri gösterilir; bir öğeyi iki kez alacak bir hesap da öyle. Bir öğe yalnızca tek bir gruptan gelebilir. `Enter` planı gösterir:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Bir dosya grubun kopyasından farklıysa ccam önce ne yapılacağını sorar (grubunkini koru, benimkini kullan ya da birleştir). Ardından değişiklik, [`ccam group`](/features/share) ile aynı kodla arka planda çalışır: önce yapılandırma kaydedilir, her hesabın verisi sırayla taşınır ve orijinaller `.ccam-backup/` klasörüne gider; MCP paylaşımından ayrılan bir hesap grubun sunucularını kendi `.claude.json` dosyasında tutar.

Çalışan bir hesabın verisini taşıyacak bir değişiklik, hesap plan açıkken başlatılmış olsa bile, hesabın adıyla birlikte reddedilir; başka bir ccam süreci grupları değiştirdiği için güncelliğini yitirmiş bir plan ise `the groups changed while this plan waited; nothing was moved, make the change again` ile reddedilir. Başka bir ccam süreci grupları değiştirirken yeni bir değişiklik, yeniden adlandırma veya kaldırma `another ccam process is changing the groups right now` der; yarıda kalmış bir değişiklikten sonra ise Health sekmesindeki `F` (veya `ccam doctor --fix`) onu tamamlayana kadar beklerler.

## Hesap ekleme

`a`, `ccam add` komutunun da terminalde kullandığı sihirbazı açar:

1. Hesap adı.
2. Nasıl giriş yapacağı: şimdi bir Claude aboneliği (tarayıcıda), bir [API anahtarı](/features/api) ya da daha sonra bir Claude aboneliği.
3. API anahtarı için: ccam'in anahtarı her başlatmada nereden okuyacağı (`env:NAME` veya `command:COMMAND`), uç nokta, anahtarın nasıl gönderileceği ve model.
4. `default` işaretli olarak katılınacak gruplar (henüz grup yoksa `default` grubunun neyi paylaşacağı).
5. İsterseniz bir kısayol komutu.
6. Onaylamanız için bir özet.

`Esc` bir adım geri gider; ilk adımda ise sihirbazı kapatır.

## Hesap kaldırma

`x` herhangi bir şey olmadan önce, Hayır seçili olarak sorar ve kaldırmanın ne yapacağını söyler: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` API anahtarlı bir hesap için iptal edilecek bir giriş yoktur; aynı girişi başka bir hesap tutuyorsa da oturum kapatılmaz. İçe aktarılmış bir hesap klasörünü korur: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## Renkler

Renkler Claude Code'un terminal arayüzündeki rolleri izler; böylece ccam onun yanında yabancı durmaz. Grup adları, adlarından türetilen bir renk alır. ccam terminale arka planını sorar ve açık ya da koyu bir tema seçer. Temayı `config.toml` dosyasının en üstünde, ilk `[[group]]`, `[[profile]]` veya `[[folder]]` satırının üstünde kendiniz ayarlayın:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

`-daltonized` temaları kırmızı-yeşil renk körlüğü içindir; `-ansi` temaları ise terminalin kendi 16 rengini kullanır. `NO_COLOR` renkleri kapatır.

## Güvenlik ve uç durumlar

- Kontrol paneli komutlarla aynı işlevleri çağırır; bu yüzden her kural geçerlidir: veri taşınmadan önce planlar, çalışan hesaplar için reddetme, yedekler ve bir hesaptan hiçbir zaman çıkmayan giriş token'ları.
- Kontrol panelinin kuru çalıştırması yoktur: ad verilmeden `ccam --dry-run`, bu klasörün hesabının nasıl başlayacağını yazdırır; bir hesapta `o` ise yalnızca o tek başlatma için bir kuru çalıştırma sunar.
- Terminal, Claude Code başlamadan veya giriş açılmadan önce eski haline getirilir.
- Bir iletişim kutusu için fazla küçük bir ekran önce sekmenin satırlarını gizler; satırlar kaydırılmak yerine `…` ile kesilir. Kısa bir ekranda ardından boş satırlar ve sekme açıklaması gizlenir, böylece tuş ipucu satırı görünür kalır. 40 sütunun veya 8 satırın altında kontrol paneli yalnızca `Too small: ccam needs 40x8 (now 36x7)` gösterir.
- Bir değişiklik, kaldırma, güncelleme veya onarım uygulanırken `Ctrl+C` onu bekler: durum satırı `ccam quits when it is done.` der.
- `U` ccam'i güncelledikten sonra, kontrol paneli bir hesap eklemekten veya girişten döndüğünde teklif yeniden gelmez.

## İlgili

- [Gruplarla paylaşım](/features/share)
- [Claude Code'u başlatma](/features/launch)
- [Hesapları yönetme](/features/accounts)
- [Kullanım limitleri](/features/usage)
