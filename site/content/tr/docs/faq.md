---
title: Sık sorulan sorular
metaTitle: "ccam SSS: birden fazla Claude Code hesabı"
description: ccam, Claude Code hesapları, kullanım limitleri, paylaşım, gizlilik ve desteklenen sistemler hakkında yanıtlar.
updated: 2026-11-19
order: 6
translation: machine
sourceHash: "090ecc3a34bd"
---

## ccam resmi bir Anthropic aracı mı?

Hayır. ccam resmi olmayan, açık kaynaklı bir projedir; Anthropic ile bağlantılı değildir ve Anthropic tarafından onaylanmamıştır. Birden fazla hesabı yan yana çalıştırmak için [Claude Code'un belgelediği](https://code.claude.com/docs/en/env-vars) `CLAUDE_CONFIG_DIR` ayarını kullanır.

## Birden fazla Claude Code hesabı kullanmama izin var mı?

Bu, her hesabın koşullarına ve iş hesapları için kurumunuzun politikasına bağlıdır: Pro ve Max için Anthropic'in [Tüketici Koşulları](https://www.anthropic.com/legal/consumer-terms), Team ve Enterprise için [Ticari Koşulları](https://www.anthropic.com/legal/commercial-terms) ve hepsi için [Kullanım Politikası](https://www.anthropic.com/legal/aup). Birden fazla hesap kullanmadan önce bunları kontrol edin. ccam bir hesabın limitlerini yükseltmez, birleştirmez veya gizlemez; her hesap yine kendi limitinde durur. İşinizi, siz istediğinizde ya da [`auto_switch`](/features/next#auto-switch) özelliğini (Linux ve macOS) açtığınızda bir limitte kendiliğinden, başka bir hesabınıza taşır. Kullanmadan önce bunun hesaplarınız için izinli olup olmadığını kontrol edin; karar sizindir.

**Team ve Enterprise koltukları.** Otomatik hesap seçimi (`ccam next`, `--fallback`, `ccam pick`, [`ccam agent`](/features/agent), kontrol panelindeki `n` tuşu ve açtığınızda [`auto_switch`](/features/next#auto-switch)), aynı Claude Team veya Enterprise kuruluşunun koltukları arasında da geçiş yapar. Anthropic'in yardım merkezinde şöyle belirtilir: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team planındaki her koltuk, adı belirli tek bir kişiye atanır ve onu yalnızca o kişi kullanabilir) ve "Users can't route their work through another user's account to use that user's remaining usage" (kullanıcılar, başka bir kullanıcının kalan kullanım hakkından yararlanmak için işlerini o kullanıcının hesabı üzerinden yürütemez) ([Team planları](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise planları](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) için de durum aynıdır). Bu nedenle otomatik hesap seçimini koltuklar arasında kullanmak Anthropic'in koşullarını ihlal edebilir. Orada kullanıp kullanmamak sizin kararınızdır: buna güvenmeden önce kuruluşunuza ve Anthropic'e danışın. Kişisel planlar (Pro, Max) yukarıdaki paragrafta ele alınır.

## ccam kullanım limitlerini atlatmanızı sağlar mı?

Hiçbir hesabın kendi limitinin ötesine geçirmez: her hesap kendi limitlerini korur ve onlarda durur; ccam bu limitleri yükseltmez, birleştirmez veya gizlemez. ccam bu limitleri gösterir ve bir hesabın limiti dolduğunda başka bir hesabınızda devam etmenizi sağlar; `auto_switch` açıksa bunu kendiliğinden yapar (bkz. sonraki soru). Aynı Team veya Enterprise kuruluşunun koltukları arasında başka bir koltukta devam etmek o koltuğun kalan kullanımını harcar; Anthropic'in yardım merkezi bunu başka bir kullanıcının koltuğu için yasaklar (bkz. yukarıdaki Team ve Enterprise koltukları).

## ccam, bir oturum limitine ulaştığında hesapları kendiliğinden değiştirebilir mi?

Yalnızca siz açarsanız. `~/.config/ccam/config.toml` dosyasının en üstünde `auto_switch = true` olduğunda (Linux ve macOS; varsayılan olarak kapalı ve deneysel), Claude Code etkileşimli bir oturumu (ccam'in başlattığı bir oturumu) 5 saatlik, haftalık, Opus, Sonnet veya Fable limitinde durdurduğunda ccam oturumu sonlandırır ve konuşmayı payı kalan başka bir hesabınızla sürdürür: aynı Team veya Enterprise kuruluşunun başka bir koltuğu ya da başka bir kişisel plan (Pro, Max). Limitten önce asla geçiş yapmaz, kullanımınızı arka planda asla izlemez ve işi asla bir Team veya Enterprise kuruluşunun dışına ya da içine taşımaz.

Koltuklarda işinizi bir koltuktan diğerine taşır; bu Anthropic'in koşullarını ihlal edebilir (bkz. yukarıdaki Team ve Enterprise koltukları). Kişisel planlar için Anthropic'in koşulları bu konuyu açıkça ele almaz: Tüketici Koşulları, Anthropic açıkça izin vermedikçe Hizmetlere "through automated or non-human means" (otomatik veya insan dışı yollarla) erişmeyi yasaklar ve bir taşıma sizin adınıza bir mesaj gönderir. Açıp açmamak sizin kararınızdır.

[`auto_switch`](/features/next#auto-switch) her koşulu ve bir taşımanın neye mal olduğunu listeler.

## Tek bilgisayarda iki Claude Code hesabını nasıl kullanırım?

Her hesaba `CLAUDE_CONFIG_DIR` ile kendi yapılandırma klasörünü verin. ccam bunu sizin için yapar: her hesap için `ccam add` komutunu çalıştırın, ardından hesapları `ccam personal` ve `ccam work` ile ya da adını kendiniz seçtiğiniz `ccwork` gibi kısayollarla başlatın. [Kılavuz](/guides/multiple-claude-code-accounts), bunun elle nasıl yapılacağını da gösterir.

## Bir konuşmaya başka bir hesapta nasıl devam ederim?

Klasörde `ccam next --continue` komutunu çalıştırın: klasörün son konuşmasını bir sonraki hesaba taşır (en çok payı kalan hesap için `--auto`). Claude Code'un içinde [ccam eklentisinin](/docs/install#claude-code-plugin) `/ccam:switch` komutu, içinde bulunduğunuz oturum için aynısını yapar. İki hesap konuşma geçmişini paylaşıyorsa (`ccam group set default --add history` ya da onlara özel bir grupla: `ccam group add pair --items history --members work,side`), yeni hesap tam olarak o konuşmayı açar; paylaşmıyorsa Claude Code konuşmayı yeni hesabın geçmişine kopyalar. Yeni hesabı konuşmanın tamamı yerine kısa bir özetle başlatmak için `--brief` ekleyin.

## ccam giriş bilgilerimi görüyor veya gönderiyor mu?

ccam, Claude girişinizi hesaplar arasında asla paylaşmaz veya kopyalamaz; MCP sunucularının girişlerini de asla kopyalamaz. Bir hesabın erişim tokenını yalnızca o hesabın kullanımını Anthropic'in kullanım uç noktasına sormak için okur: `ccam usage`, kontrol paneli, `ccam doctor` ve otomatik seçim için (`ccam next --auto`, `ccam pick`, `--fallback` ve bir oturum limitine ulaştığında `auto_switch`). Süresi dolmuş bir tokenı yalnızca `ccam usage --refresh` ile yeniler; `usage_source = "claude"` ile hiçbir yere token göndermez ve kayıtlı girişi (dosyayı ya da macOS Anahtar Zinciri kaydını) yalnızca hâlâ çalışıp çalışmadığını anlamak için okur. `ccam ext` ve `ccam doctor`, bir hesabın MCP girişlerinin yalnızca adlarını okur, tokenlarını asla okumaz. ccam'da telemetri yoktur. Bu sorguların dışında kendiliğinden yaptığı tek istek, kontrol panelinin GitHub'da daha yeni bir ccam sürümü olup olmadığını denetlemesidir ve bu istek sizin ya da hesaplarınız hakkında hiçbir şey taşımaz.

## Uzak bir sunucuda nasıl giriş yaparım?

SSH oturumunda her zamanki gibi `ccam add` veya `ccam login <name>` çalıştırın. Claude Code girişi sunucuda beklediği için bilgisayarınızdaki tarayıcı, yüklenmeyen bir `http://localhost:<port>/callback?code=…` sayfasında kalır. Bu adresi tarayıcıdan kopyalayıp terminale yapıştırın: ccam onu sunucudaki Claude Code'a iletir ve Claude Code girişi tamamlar. Adres tek kullanımlık bir kod içerir, bu yüzden kimseyle paylaşmayın.

## Hesaplarım nerede saklanıyor?

Her hesap için bir klasör olmak üzere `~/.local/share/ccam/profiles/<name>/` içinde. macOS üzerinde Claude Code, girişin kendisini Anahtar Zinciri'nde, o klasörden türetilen bir ad altında saklar; bu yüzden `ccam rename` orada klasörün eski adını korur.

## Mevcut ~/.claude klasörüme ne olur?

`ccam import` veya `ccam migrate` çalıştırmadığınız sürece hiçbir şey. İçe aktarma, klasörü bulunduğu yerde kaydeder; böylece düz `claude` ve `ccam <name>` aynı hesabı ve ayarları kullanmaya devam eder. `ccam migrate` ise klasörün ayarlarını, `CLAUDE.md` dosyasını, uzantılarını ve MCP sunucularını paylaşılan `default` grubuna kopyalar ve klasörü yalnızca okur.

## ccam hangi sistemleri destekliyor?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine ve diğerleri), macOS ve Windows 10 ile 11, yerel olarak veya WSL içinde. Windows'ta `settings.json` gibi dosyaları paylaşmak Geliştirici Modu gerektirir; paylaşılan klasörler onsuz da çalışır ve Geliştirici Modu olmadan ilk `ccam add`, yalnızca MCP sunucuları ve uzantılarla `default` grubunu oluşturur ([Windows'ta paylaşım](/features/share#on-windows)).

## Bir MCP sunucusu neden yalnızca bir hesapta görünüyor?

Bir oturum içinde `claude mcp add` ile eklenen sunucular yalnızca o hesap için, `-s user` olmadan da yalnızca geçerli klasör için saklanır. ccam, hesabı bir sonraki başlatışınızda bunu size bildirir; `ccam mcp list` ve `ccam ext` de bunları gösterir; sunucuyu `ccam mcp import <account> <server>` ile paylaşılan gruba taşıyın.

## ccam, Claude Code'u yavaşlatır mı?

Claude Code'u ccam üzerinden başlatmak yaklaşık 4 milisaniye ekler; bu süre bir dizüstü bilgisayarda `scripts/perf-check.sh` ile ölçülmüştür. ccam küçük bir dosya okur, paylaşılan bağlantıları kontrol eder ve ardından kendini `claude` ile değiştirir; dolayısıyla varsayılan olarak siz çalışırken çalışır durumda değildir. Windows bir süreci değiştiremez; bu yüzden orada ccam, Claude Code'u başlatır ve başka hiçbir şey yapmadan onu bekler. `relay` veya `auto_switch` açık olduğunda ccam, konuşmaya başka bir hesapta devam edebilmek için her etkileşimli oturumun yanında aynı şekilde çalışmaya devam eder.

## ccam'ı nasıl güncellerim?

`ccam update` komutunu çalıştırın. Bu komut en son sürümü indirir, sağlama toplamını kontrol eder ve `ccam` dosyasını değiştirir; `ccam update --check` ise yalnızca daha yeni bir sürüm olup olmadığını söyler. Kontrol paneli de daha yeni bir sürüm olduğunda bunu söyler ve `U` onu sorduktan sonra kurar; ccam hiçbir zaman kendiliğinden güncellenmez. `.deb` paketini kurduysanız bunun yerine yeni `.deb` dosyasını apt ile kurun. [Güncelleme](/docs/install#update) bölümüne bakın.

## ccam'ı nasıl kaldırırım?

`ccam uninstall` çalıştırın. Nelerin silineceğini listeler ve önce sorar, ardından her hesabı (önce oturumunu kapatarak), kısayol komutlarını, `~/.config/ccam`, `~/.local/share/ccam`, `~/.local/state/ccam` klasörlerini ve `ccam` komutunu kaldırır. İçe aktarılmış bir `~/.claude` yerinde kalır ve gruplarında paylaştığı şeylerin bir kopyası ona geri konur. [Kurulum](/docs/install#uninstall) sayfasına bakın.
