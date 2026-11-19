---
title: Diğer hesaplarda alt ajanlar
description: ccam agent, Claude Code görevinin bir kısmını başka bir hesabınızla çalıştırıp Claude'un yanıtını döndürür; uzun bir görev birkaç hesabın kullanımını harcar.
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Claude Code'un kendi alt ajanları, onları başlatan oturumun hesabıyla çalışır; bu yüzden o hesabın kullanımını harcar. `ccam agent "<task>"`, bir görevi aynı klasörde `claude -p` ile diğer hesaplarınızdan biriyle çalıştırır ve Claude'un son yanıtını yazdırır. Bir Claude Code oturumundan çağrıldığında başka bir hesabın kullanımını harcayan bir alt ajan gibi çalışır; böylece uzun bir görev, oturumun kendi hesabının payı tükenecek noktadan sonra da sürebilir.

**Team ve Enterprise koltukları.** `ccam agent` hesabı kendisi seçer. Tek bir Claude Team veya Enterprise kuruluşu içinde görevi başka bir koltukta çalıştırır ve bu, Anthropic'in koşullarını ihlal edebilir. Buna güvenmeden önce [SSS'yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) okuyun ve kuruluşunuza ve Anthropic'e danışın.

## Kullanım

Claude Code onu Bash aracıyla çalıştırır. [ccam eklentisi](/docs/install#claude-code-plugin) kuruluyken bunu sade sözcüklerle isteyin, örneğin "bu görevin parçaları için diğer hesaplarımı kullan", ya da bunu CLAUDE.md dosyanıza yazın. Eklentinin becerisi Claude Code'a her ajanı arka planda başlatmasını, tek başına anlaşılır bir görev yazmasını ve dosya düzenleyen ajanlara kendi parçalarını ya da kendi git worktree'lerini vermesini söyler.

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

Claude'un son yanıtı standart çıktıya gider; `ccam:` ile başlayan satırlar standart hataya gider. `--json` bunun yerine `{account, session_id, result, is_error, accounts}` yazdırır; burada `accounts`, görevin çalıştığı her hesabı listeler.

Çıkış kodları: claude'un kendi kodu (ajan bittiğinde `0`), çağıran dışında hiçbir hesapta pay kalmadığında `75` (`EX_TEMPFAIL`), bir ajanın başka bir ajan başlatmaya çalışması gibi diğer hatalar için `1`.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--model <model>` | Ajanın modeli; o modelin haftalık penceresi de seçimde sayılır | hesabın varsayılan modeli |
| `--resume <id>` | Önceki bir ajanın konuşmasını bir devam göreviyle sürdürür (ilk ve son satırındaki oturum kimliği) | kapalı |
| `--json` | Hesabı, oturumu, yanıtı ve görevin çalıştığı hesapları yazdırır | kapalı |
| `--max-usage <1-100>` | Bu kullanımın altındaki hesapları tercih eder | `auto_max_usage`, aksi halde 100 |
| `--strategy headroom\|expiring\|balance` | [Seçim kuralı](/features/next#auto) | `auto_strategy`, aksi halde `headroom` |
| `--cross-org` | Başka kuruluşların hesaplarıyla da çalıştırır | `cross_org`, aksi halde kapalı |
| `--wait[=<duration>]` | Hiçbir hesapta pay yoksa ilk hesabın sıfırlanmasını en fazla bu süre kadar bekler | kapalı; tek başına verilirse `5h` |
| `--dry-run` | Seçeceği hesap için komutu, çalıştırmadan yazdırır | kapalı |

## Nasıl çalışır

1. **Çağıran dışarıda bırakılır.** `ccam agent` komutunun çağrıldığı hesap (ccam'in başlattığı oturumun hesabı, yoksa bu klasörün başlattığı hesap) ajanı hiçbir zaman çalıştırmaz; o hesapla limitlerini paylaşan, aynı Claude hesabına giriş yapmış başka bir klasör de çalıştırmaz. `--cross-org` aksini söylemedikçe ajan, çağıranın [kuruluşunda](/features/next#auto) kalır.
2. **`ccam pick` gibi seçer.** Kullanım, [`ccam pick`](/features/pick) komutunun yaptığı gibi sorgulanır ve birbirinden en fazla 2 dakika arayla başlatılan ajanlar, payı kalan yeterli hesap olduğu sürece farklı hesaplar alır.
3. **Oturumun izinleriyle çalışır, asla daha fazlasıyla değil.** `claude -p` içinde bir araca izin verecek kimse yoktur; bu yüzden ajan, onu çağıran oturumun izin modunu ve efor düzeyini alır. Bir terminalden çalıştırıldığında Claude Code'un varsayılan modu ve ayarlarınızdaki izin kuralları geçerli olur. `ccam agent` komutunun daha fazlasını veren bir seçeneği yoktur.
4. **Bir kullanım limitinde konuşma devam eder.** Ajanın hesabı limitine ulaştığında ccam limiti [eklentinin hook'unun](/features/usage#the-limit-record) yapacağı gibi kaydeder ve aynı konuşmayı en çok payı kalan hesapla sürdürür: Claude Code konuşmayı oraya kopyalar (`--resume <file> --fork-session`) ve ilk mesaj "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." olur. Her hesap bir kez denenir.
5. **Her ajan konuşması kaydedilir:** `~/.local/state/ccam/agent/<session id>.json` (30 gün). `--resume <id>` konuşmayı kendi klasöründe, onu çalıştıran hesapta hâlâ pay varsa o hesapla, yoksa başka bir hesapla sürdürür. [`ccam next --continue`](/features/next) ajan konuşmalarını dışarıda bırakır; böylece sizin konuşmanızı sürdürmeye devam eder.

## Güvenlik ve uç durumlar

- Bir ajan kendi ajanlarını başlatamaz: ccam onu `CCAM_AGENT=1` ile işaretler ve `ccam agent` orada çalışmayı reddeder.
- `ccam agent` durdurulursa (Ctrl+C, bir sinyal ya da Claude Code'un komutu kesmesi) ajan da durur.
- Ajan, onu çağıran konuşmayı görmez; yalnızca kendi görevini görür.
- Aynı dosyaları aynı anda düzenleyen ajanlar birbirlerinin işinin üzerine yazar. Her birine görevin kendi parçasını ya da kendi git worktree'sini verin.
- Çağıran oturum başka bir hesaba taşındığında ([`auto_switch`](/features/next#auto-switch)) başlattığı ajanlar da onunla birlikte durur. Konuşmaları kalır: `--resume <id>` onları sürdürür.
- `agent` bir komut adıdır, bu yüzden hiçbir yeni hesaba `agent` adı verilemez; bu adı taşıyan mevcut bir hesap yine `ccam run agent` ile başlar.

## İlgili

- [Betikler için hesap seçme](/features/pick)
- [Sonraki hesaba geçme](/features/next)
- [Kullanım limitleri](/features/usage)
- [Her hesabın arka plan oturumları](/features/agents)
