---
title: Kullanım limitleri ve hesap değiştirme
description: Tüm Claude Code hesaplarının 5 saatlik ve haftalık kullanımını ccam usage ile görün, en çok pay kalan hesabı ccam next --auto ile açın.
updated: 2026-11-05
order: 5
translation: machine
sourceHash: "06bd349073e3"
---

`ccam usage`, her hesap için her limitin ne kadarının kullanıldığını gösterir: 5 saatlik limit, haftalık limit ve her modelin kendi haftalık limiti (Fable gibi). `ccam next --auto` en çok pay kalan hesabı açar, `--continue` ise mevcut konuşmayı o hesaba taşır.

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

Her hesap, sayılarının ne kadar eski olduğunu gösteren bir satır alır; ardından her limit için sıfırlanmaya kalan süreyle birlikte bir çubuk gelir: 5 saatlik pencere, tüm modeller için haftalık pencere ve her modelin kendi haftalık penceresi, örneğin Fable. Bir çubuk %80'de sarıya, %100'de kırmızıya döner. Son satır, `ccam next --auto` komutunun açacağı hesabı adlandırır; `auto_strategy = "expiring"` veya `"balance"` ile bunun yerine `→ next --auto picks <name>: <reason>` yazar.

**Team ve Enterprise koltukları.** Otomatik hesap seçimi (`ccam next`, `--fallback`, `ccam pick`, kontrol panelindeki `n` tuşu ve açtığınızda [`auto_switch`](/features/next#auto-switch)), aynı Claude Team veya Enterprise kuruluşunun koltukları arasında da geçiş yapar ve bu, Anthropic'in koşullarını ihlal edebilir. Buna güvenmeden önce [SSS'yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) okuyun ve kuruluşunuza ve Anthropic'e danışın.

## Sayılar nereden geliyor

ccam, her hesabın kendi giriş bilgilerini kullanarak tüm hesaplar için aynı anda Anthropic'in kullanım uç noktasına sorgu gönderir. Hesap başına en fazla 3 saniye bekler ve her hesabı yanıt verir vermez gösterir.

Bu uç nokta herkese açık olarak belgelenmemiştir ve değişebilir. Bir istek başarısız olduğunda ccam bilinen son sayıları gösterir ve bunları önbellekten alınmış olarak işaretler. Bu sayılar kaydedildikten sonra sıfırlanmış bir pencere %0 gösterir. Yanıt biçimi değiştiyse `ccam doctor` bunu size bildirir.

Claude Code bir oturumu kullanım limitinde durdurduğunda [ccam eklentisi](/docs/install#claude-code-plugin) hangi pencerenin dolduğunu ve ne zaman sıfırlanacağını kaydeder. O zamana kadar o pencere, diğer sayılar ne derse desin %100 gösterir; böylece `ccam next --auto` ve `--fallback` hesabı çok erken yeniden seçmez.

60 saniyeden daha kısa süre önce alınan sayılar yeniden kullanılır. Seçenekler:

| Seçenek | Etki |
| --- | --- |
| `--watch[=seconds]` | Sürekli yeniler, varsayılan olarak 60 saniyede bir |
| `--json` | Makinece okunabilir çıktı |
| `-v` | Her hesabın sayılarının nereden geldiğini de gösterir |
| `--no-cache` | Her zaman yeniden sorar |
| `--source cache` | Hiç sormaz; yalnızca kayıtlı sayıları gösterir |

## Süresi dolmuş girişler

Bir hesabın erişim tokenının süresi dolduğunda ccam onu yenilemez; çünkü aynı giriş bilgilerini bir Claude Code oturumu kullanırken token yenilenirse o oturumdan çıkış yapılabilir. Hesap ``access token expired; `ccam <name>` renews it when it starts`` gösterir: onu bir kez `ccam <name>` ile başlatın, Claude Code tokenı yeniler. Tümüyle süresi dolmuş ya da Claude Code tarafından silinmiş bir giriş `` login expired → `ccam login <name>` `` gösterir ve yeni bir giriş gerektirir.

`ccam usage --refresh` süresi dolmuş tokenları kendisi yeniler, ancak yalnızca, nasıl başlatılmış olursa olsun, çalışan bir Claude Code oturumu olmayan hesaplar için. macOS Anahtar Zinciri'nde veya içe aktarılmış bir `~/.claude` klasöründe saklanan girişleri hiçbir zaman yenilemez ve bir yenileme 30 saniyeye kadar sürebilir.

## Durum satırında kullanım

```sh
ccam statusline install
```

Bundan sonra Claude Code, durum satırında hesap adını ve kullanımını gösterir, örneğin `[work] · Opus · 5h 24% · 7d 41%`. %90'dan itibaren satır ayrıca `ccam next --auto --continue` komutunu önerir. Zaten bir durum satırı komutunuz varsa ccam onu korur ve hesap adını başına ekler. Durum satırı sayılarını da kaydeder; uç nokta yanıt vermediğinde `ccam usage` bu sayıları kullanır.

## Geçiş

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the most room
ccam next --continue  # continue this folder's conversation on the next account
```

`ccam next` giriş yapılmamış hesapları atlar, `--auto` ise %100'deki hesapları (sıfırlanana kadar Claude Code'un dolu bildirdiği bir pencereyi de) ve başka bir kuruluşun hesaplarını atlar. İki hesap konuşma geçmişini paylaşıyorsa `--continue` aynı konuşmayı açar; paylaşmıyorsa Claude Code konuşmayı yeni hesaba kopyalar. Bkz. [paylaşım](/docs/sharing) ve [sonraki hesaba geçme](/features/next#continue).

`config.toml` dosyasının en üstünde `auto_switch = true` olduğunda (Linux ve macOS; varsayılan olarak kapalı), Claude Code etkileşimli bir oturumu limitinde durdurduğunda ve aynı türden başka bir hesabın (aynı kuruluştaki bir koltuk ya da başka bir kişisel plan) payı kaldığında ccam bu geçişi kendisi yapar. Ne zaman devreye girdiği ve asla ne yapmadığı için bkz. [`auto_switch`](/features/next#auto-switch).

`claude -p` çalıştıran betikler ve zamanlanmış işler için hesap adından önce `--fallback` koyun. O hesap limitine ulaştığında ccam aynı komutu en çok payı kalan hesapla yeniden çalıştırır ve yalnızca tamamlanan çalıştırmanın çıktısı yazdırılır:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

ccam'in nasıl karar verdiği için [Limitte yeniden çalıştırma](/features/next#run-again-at-the-limit) bölümüne bakın.

[Kullanım limiti iş akışı kılavuzu](/guides/usage-limit-workflow) eksiksiz bir örneği adım adım anlatır.
