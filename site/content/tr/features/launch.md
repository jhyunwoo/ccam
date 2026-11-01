---
title: Claude Code'u başlatma
description: Claude Code'u herhangi bir hesapla ccam <name|number>, ccam run veya kontrol paneliyle başlatın; argümanlar claude'a aynen gider, --dry-run komutu gösterir.
updated: 2026-11-01
order: 1
translation: machine
sourceHash: 11df08608a1e
---

`ccam work`, Claude Code'u `work` hesabıyla başlatır ve addan sonraki her şey `claude` komutuna değiştirilmeden iletilir. `claude` yazacağınız her yerde kullanın: terminalde, bir betikte veya girdiyi boruyla aktararak.

## Kullanım

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Bir hesabın adını veya `ccam list` çıktısındaki numarasını verin. ccam addan sonraki argümanları okumaz; bu yüzden `ccam work --help` Claude Code'un yardımını gösterir. Addan sonra `--` yazmak isteğe bağlıdır. Standart girdi, çıkış kodu ve sinyaller `claude` komutundaki gibi davranır.

İlk argüman `-` ile başlıyorsa ve ccam'in kendi seçeneklerinden biri değilse ccam hesabı sizin için seçer ve tüm argümanları Claude Code'a iletir. Geçerli klasör için [`ccam folder`](/features/folder) ile ayarlanan hesabı, yoksa varsayılan hesabı kullanır: `ccam use` ile ayarlanan hesap ya da ayarlamadıysanız en küçük numaralı hesap.

### Kontrol paneli

Argümansız `ccam`, standart girdi ve çıktı bir terminal olduğunda [kontrol panelini](/features/dashboard) açar, aksi halde `ccam list` çıktısını yazdırır. Bir hesaba gidin (ok tuşlarıyla ya da numarasını yazarak) ve başlatmak için `Enter` tuşuna basın; `n` en çok payı kalan hesabı başlatır, `o` ise `--continue`, `--resume` veya başka argümanlar sunar. Kontrol paneli Claude Code başlamadan önce kapanır; böylece terminal önce eski haline gelir.

**Team ve Enterprise koltukları.** Otomatik hesap seçimi (`ccam next`, `--fallback`, `ccam pick`, kontrol panelindeki `n` tuşu ve açtığınızda [`auto_switch`](/features/next#auto-switch)), aynı Claude Team veya Enterprise kuruluşunun koltukları arasında da geçiş yapar ve bu, Anthropic'in koşullarını ihlal edebilir. Buna güvenmeden önce [SSS'yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) okuyun ve kuruluşunuza ve Anthropic'e danışın.

Kabukta `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` veya `ANTHROPIC_AUTH_TOKEN` ayarlıysa Claude Code, hesabın kendi girişi yerine o anahtarı ya da token'ı kullanabilir; bu yüzden ccam başlamadan önce standart hatada sizi uyarır (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Değişkenleri yalnızca okur; bir [API anahtarlı hesap](/features/api) ise onları temizler.

Varsayılan olarak ccam kendini Claude Code ile değiştirir, böylece başka hiçbir şey çalışmaya devam etmez. `config.toml` dosyasının en üstünde, ilk `[[group]]`, `[[profile]]` veya `[[folder]]` satırının üstünde `relay = true` (ya da Linux ve macOS'ta [`auto_switch = true`](/features/next#auto-switch)) olduğunda ccam bunun yerine etkileşimli oturumları bir alt süreç olarak başlatır ve yanlarında çalışmaya devam eder. Bu sayede `/ccam:switch` ya da kullanım limitinde `auto_switch`, konuşmaya oturumun seçenekleri, izin modu, efor düzeyi ve ajanıyla başka bir hesapta kendiliğinden devam edebilir; `/ccam:switch` konuşmayı başka bir kuruluşun hesabına taşımadan önce ccam size sorar ([ayrıntılar](/features/next#relay)).

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--dry-run` | ccam'in çalıştıracağı ortamı ve komutu yazdırır, ardından çıkar | kapalı |
| `--fallback` | `-p`/`--print` çalıştırmaları için: hesap kullanım limitine ulaştığında aynı komutu en çok payı kalan girişli hesapla yeniden çalıştırır ([ayrıntılar](/features/next#run-again-at-the-limit)) | kapalı |
| `-v`, `--version` | `ccam <version>` yazdırır (yalnızca ilk argüman olarak) | |
| `-h`, `--help` | `ccam --help` veya `ccam run --help` olarak ccam'in yardımı | |

Argüman kuralları:

- ccam kendi seçeneklerini yalnızca hesap adından önce okur. Addan sonra `--dry-run` ve `--help` dahil her argüman Claude Code'a aittir.
- `--dry-run`, `ccam --dry-run work`, `ccam run --dry-run work` ve `ccam --dry-run` (klasörün hesabı ya da varsayılan hesap için) olarak çalışır. `--dry-run` ve `--fallback` herhangi bir sırayla yazılabilir.
- `ccam run <name|number> [claude args]`, `ccam <name>` komutunun uzun biçimidir. Adsız `ccam run` klasörün hesabını ya da varsayılan hesabı kullanır.
- Hesap adları `add`, `list`, `run` veya `usage` gibi ccam komut adları olamaz; böylece `ccam <word>` hiçbir zaman belirsiz olmaz.

Kuru çalıştırma, bir kabuğa ya da Windows'ta PowerShell'e yapıştırabileceğiniz tek bir satır yazdırır. Satır, başlatmanın ortamda değiştirdikleriyle başlar: her zaman `CLAUDE_CONFIG_DIR` ve `CCAM_PROFILE`, kaldırdığı bir değişken için (devralınmış bir `CCAM_RELAY` gibi) `env -u <name>` ve bir [API anahtarlı hesap](/features/api) için uç noktası, modeli ve anahtarın nereden geldiği (`ANTHROPIC_API_KEY="$WORK_KEY"`); anahtarın kendisi asla yazılmaz:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

PowerShell satırı, bir kabuktaki `VAR=x cmd` gibi değişkenleri sonradan geri koyar ve komut satırını Claude Code'a `Start-Process` üzerinden verir; böylece Windows PowerShell 5.1 ve PowerShell 7 onu aynı şekilde çalıştırır. Bir başlatma bir değişkeni boş bir değere ayarladığında, altındaki bir `#` satırı PowerShell 7.5'ten önceki sürümlerin böyle bir değişkeni bunun yerine kaldırdığını söyler. npm ile kurulmuş bir `claude.cmd` için satır, başlatmanın kullanacağı `cmd.exe` komutunu çalıştırır (aşağıya bakın) ve başlatmanın orada reddedeceği bir argüman kuru çalıştırmayı da başarısız kılar.

## Nasıl çalışır

Her başlatma, komut satırından `exec` çağrısına kadar aynı adımları izler:

<!-- motion:launch -->

Her kutu, adımı çalışırken yanar; alttaki satırlar ise Claude Code'un başladığı ortam ve argümanlardır.

1. ccam, `cmd/ccam/main.go` dosyasının çağırdığı `app.Main` içinde başlar. ccam başka bir adla başlatıldıysa bu bir [kısayoldur](/features/shortcuts) ve doğrudan başlatıcıya gider.
2. ccam, cobra devreye girmeden önce `os.Args` değerine bakar (`app.parse`). Yalnızca tek başına `ccam` (kontrol paneli), bir ccam komut adı (`config.Reserved`), `-h`/`--help` ve `-v`/`--version`, `cli.Execute` fonksiyonuna gider. Geri kalan her şey bir başlatmadır.
3. ccam, ayrıştırdığı tek dosya olan `~/.config/ccam/config.toml` dosyasını okur (`launcher.Launch`).
4. ccam hesabı ada veya numaraya göre bulur, baştaki bir `--` değerini atar ve `claude` yürütülebilir dosyasını bulur (`launcher.Prepare`, `ResolveClaude`). Ad yoksa çalışma dizinini `config.toml` içindeki `[[folder]]` girdileriyle dize olarak, macOS ve Windows'ta büyük/küçük harfe bakmadan karşılaştırarak klasörün hesabını alır (`HereProfile`). PATH araması, dosyanın değiştirilme zamanı ve `PATH` değerinin bir hash'iyle birlikte `~/.local/state/ccam/claude-path` içinde önbelleğe alınır; böylece ccam her başlatmada `PATH` içinde arama yapmaz. Değişmiş veya eksik bir dosya ya da değişmiş bir `PATH` yeni bir aramayı tetikler. `config.toml` içindeki bir `claude_bin` yolu (`~/` genişletilir) çalıştırılabilir bir dosya olmalıdır; değilse başlatma, herhangi bir şey yazılmadan önce `claude_bin X is not an executable file. Check claude_bin in config.toml` ile durur. ccam, inode ile karşılaştırarak ccam'in kendisi olan bir `claude` dosyasını reddeder, çünkü bu bir döngüye yol açar.
5. ccam, hesap klasörü `~/.local/share/ccam/profiles/<name>/` ile paylaşım ayarlarını yalnızca `lstat` ve `readlink` çağırarak, öğe başına bir veya iki sistem çağrısıyla karşılaştırır (`checkLinks`, `share.Check`). Eksik bir bağlantı anında oluşturulur. Veri kaybettirebilecek her şey başlatmayı durdurur.
6. ccam paylaşılan MCP sunucularını ekler (`applyMCP`; aşağıya bakın).
7. ccam oturumu kaydeder (`session.Register`). PID'sinin adını taşıyan ve süreç başlangıç zamanını içeren bir işaret dosyasını, `XDG_RUNTIME_DIR` ne derse desin, `~/.local/state/ccam/locks/profile-<name>/<pid>` konumuna yazar; böylece her kabuk aynı işaretleri görür. Konuşma geçmişini paylaşan bir hesap için ccam önce `history-<group>-<hash>` içinde çalışan bir oturum arar (hash, sembolik bağlantıları çözülmüş geçerli klasörün hash'idir) ve oraya da kaydolur (`historyCheck`).

   `exec` PID'yi korur; böylece işaret dosyası artık Claude Code sürecinin kendisini adlandırır. ccam tutulan bir kilit yerine bir işaret dosyası kullanır; çünkü Claude Code'un başlattığı süreçler bir kilidi devralır ve oturum bittikten sonra da tutar. Yazılamayan bir işaret başlatmayı durdurmaz: ccam `ccam: could not mark the session as running (…); starting Claude Code anyway` der.

   ccam'in bir hesabın çalışıp çalışmadığına dair diğer denetimleri (yeniden adlandırma, kaldırma, grup değişiklikleri, `doctor --fix`, MCP eşitlemesi) hesap klasöründeki Claude Code'un kendi kayıtlarını da (`sessions/<pid>.json`, `daemon.lock`) okur; böylece ccam olmadan başlatılan Claude Code'u da görürler: düz `claude`, bir IDE, arka plan oturumları ve onların daemon'ı.
8. ccam terminal başlığını `claude · <name>` olarak ayarlar (standart hata bir terminal olduğunda), başlatma zamanını `fsync` olmadan `~/.local/state/ccam/last-used/<name>` dosyasına yazar ve `syscall.Exec` çağırır (`Plan.Exec`).

Claude Code, `CLAUDE_CONFIG_DIR` hesap klasörüne ve `CCAM_PROFILE` hesap adına ayarlanmış olarak çalışır. ccam, bu değişkenlerin devralınmış değerlerini ve devralınmış bir `CCAM_RELAY` değişkenini temizler. macOS'ta `ccam import` ile kaydedilmiş bir `~/.claude` klasörü için ayrıca `CLAUDE_SECURESTORAGE_CONFIG_DIR` değişkenini boş bir değere ayarlar; böylece hesap ve düz `claude` aynı Anahtar Zinciri girişini kullanır; Windows'ta Claude Code girişi Kimlik Bilgileri Yöneticisi'nde tuttuğunda da aynısını yapar.

### Neden exec

ccam kendi sürecini `claude` ile değiştirir (`syscall.Exec`). Bunun sonucunda:

- ccam bellekte kalmaz;
- tuş vuruşları, sinyaller ve çıkış kodu doğrudan Claude Code'a aittir;
- argümanlar hiçbir zaman bir kabuk dizesi olmadan dizi olarak kalır.

Windows'ta `exec` yoktur. Orada ccam, `claude.exe` dosyasını aynı ortam, argümanlar ve konsolla bir alt süreç olarak başlatır, onu bekler ve onun çıkış koduyla çıkar (`Plan.Exec`). npm'in kurduğu bir `claude.cmd`, tam yoluyla ve `cmd.exe` hiçbirini kendi sözdizimi (`&`, `%PATH%`) olarak okumasın diye her argümanı tırnak içine alınarak `cmd.exe` içinde çalışır; satır sonu içeren bir argüman bu şekilde iletilemez ve yerel Claude Code'u kurma önerisiyle reddedilir.

Ctrl+C, Claude Code'a doğrudan konsoldan ulaşır; bu yüzden ccam onu yok sayar. Oturum işaretçisi, tam olarak Claude Code kadar yaşayan ccam'in PID'sini taşır ve Claude Code çıktığında kaldırılır. Beklerken ccam başka hiçbir şey yapmaz.

Bir oturumdan sonraki işler, örneğin `claude mcp add` ile eklenen sunucuları fark etmek, bir sonraki başlatmada yapılır: ccam, `.claude.json` dosyasının değiştirilme zamanını `~/.local/state/ccam/mcp/<name>/seen` ile karşılaştırır ve dosyayı yalnızca değiştiyse ayrıştırır.

### Bir komut satırı nasıl yönlendirilir

Üç komut satırı, `internal/app` içinden üç yol:

<!-- motion:launch-route -->

Bir kısayol `argv[0]` ile, `ccam <name>` ise ilk argümanla eşleşir ve ikisi de `syscall.Exec` ile biter. Yalnızca `ccam usage` gibi bir ccam komutu cobra'ya ulaşır.

### Paylaşılan MCP sunucuları

Varsayılan `flag` stratejisiyle ccam, grubun MCP sunucularını Claude Code'a sizinkilerden önce yerleştirilen bir argüman olarak iletir:

<!-- motion:launch-mcp -->

ccam argümanı yalnızca grubun `mcp.json` dosyası varsa ekler ve onu `claude` komutunun hemen arkasına koyar (`mcp.FlagArgs`, `Plan.InsertArgs`). Kendi `--mcp-config` veya `--strict-mcp-config` seçeneğiniz bunun üzerine yine uygulanır. Bu tek bir argümandır, `--mcp-config=<path>`, çünkü Claude Code'un `--mcp-config` seçeneği birden fazla değer alır: iki argüman olarak yazılsaydı ardından gelen bir istemi de yapılandırma dosyası olarak okurdu.

Bayrağın nereye gidemeyeceğine ccam karar verir (`claudeVerb`). Claude Code'un alt komutunu, Claude Code'un yaptığı gibi, seçenekleri değerleriyle birlikte atlayarak argümanlardan okur. `remote-control` veya `rc` bayrak almaz (aşağıda); ilk sırada verilen bir arka plan oturumu alt komutu da (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`) almaz. Claude Code orada bayrağı reddettiği için `ccam work mcp list`, `get` ve `remove` de almaz; bunun yerine ccam standart hataya grubun sunucularının hesabın oturumlarında da yüklendiğini yazar:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Diğer her alt komut ve her oturum bayrağı korur.

`config.toml` içinde `mcp_strategy = "sync"` ile ccam bunun yerine sunucuları hesabın `.claude.json` dosyasına yazar ve hesap zaten çalışıyorsa bunu atlar. Son denetimden beri ne grubun `mcp.json` dosyası ne de `.claude.json` değiştiyse bir başlatma bu işi de atlar. Hesabın bir grup sunucusuyla aynı adı taşıyan kendi sunucusu varsa eşitleme hesabın sunucusunu korur ve bunu bir kez söyler. Eşitleme ayrıca `.claude.json` dosyasından kaybolmuş bir kopyayı geri koyar (`put back group default's MCP servers …`) ve yalnızca kendi yazdığı ve o zamandan beri kimsenin değiştirmediği kopyaları değiştirir veya kaldırır.

### Remote Control

`ccam work rc` (veya `ccam work remote-control`), Claude Code'un Remote Control özelliğini `work` hesabıyla başlatır. `claude remote-control`, `--mcp-config` seçeneğini alt komutundan önce reddeder, sonrasında da hiç almaz; oturumları ise MCP sunucularını yalnızca hesabın `.claude.json` dosyasından okur. Bu yüzden bayrak stratejisi burada farklı çalışır:

<!-- motion:launch-rc -->

1. ccam, seçeneklerden ve değerlerinden sonra bile (`ccam work --model opus rc`) Claude Code'un alt komutu olarak `remote-control` veya `rc` bulur ve bayrak eklemez (`claudeVerb`, `mcp.FlagArgs`).
2. ccam, grubun `mcp.json` dosyasının değiştirilme zamanını hesabın `sync-stamp` dosyasıyla karşılaştırır ve her grup sunucusunun hâlâ `.claude.json` içinde olduğunu denetler (`mcp.Synced`). `.claude.json` grubun güncel sunucularını zaten içeriyorsa veya grubun hiç sunucusu yoksa Remote Control soru sormadan başlar.
3. Aksi halde ccam, oturumlarda eksik kalacak sunucuları adlarıyla belirtir ve sorar:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

4. `y` yanıtında ccam bunları `sync` stratejisinin yaptığı gibi kopyalar (`mcp.Sync`): Claude Code'un kilidi altında, bir yedek aldıktan sonra ve hesabın kendi sunucularını koruyarak. `N` yanıtında veya terminal yoksa Remote Control onlar olmadan başlar ve ccam, bunları her başlatmadan önce kopyalayan `mcp_strategy = "sync"` ayarını önerir.

ccam, sunucular kopyalanana kadar her Remote Control başlatmasında sorar. Kopyalamadan sonra yalnızca grubun `mcp.json` dosyası değiştiğinde veya kopyalanmış bir sunucu `.claude.json` dosyasından kaybolduğunda yeniden sorar. `.claude.json` içindeki bir kopya, sonraki normal başlatmalarda `--mcp-config` seçeneğine üstün gelir; bu yüzden `ccam mcp add`, `remove` ve `import` diğer hesapların tuttuğu kopyaları günceller ([Paylaşılan MCP sunucuları](/features/mcp)).

### 10 ms bütçesi

Başlatma yolunun bütçesi 95. yüzdelikte 10 ms'dir; pratikte ccam, Claude Code başlamadan önce yaklaşık 4 ms ekler. Başlatma yolu hiçbir zaman cobra komut ağacını oluşturmaz, Bubble Tea başlatmaz veya bir ağ bağlantısı açmaz; paket `init()` fonksiyonları toplamda 5 ms'den uzun sürerse `TestPackageInitIsCheap` başarısız olur. `scripts/perf-check.sh`, `ccam --dry-run work`, `ccam work` ve bir kısayolu `hyperfine` ile ölçer ve 10 ms'nin üzerinde başarısız olur; `internal/launcher` içindeki `BenchmarkLaunchPath` ise adımları bir döngüde çalıştırır. Adım başına süreyi görmek için `CCAM_TRACE=1` ayarlayın:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Güvenlik ve uç durumlar

- Bilinmeyen bir ad veya numara `no account "<name>"` ile durur ve `ccam list` komutunu gösterir. Hiç hesap yoksa ccam, `ccam add <name>` ile bir hesap eklemenizi ister.
- Paylaşılan bir yol sıradan bir dosyaya dönüştüyse veya yanlış yere yöneliyorsa ccam `stopped because of share link problems:` ile durur, sorunları listeler ve `ccam doctor --fix` çalıştırmanızı ister. Başlatma yolunda veri taşımaz.
- Bu hesapla konuşma geçmişini paylaşan (aynı gruptaki) başka bir Claude Code oturumu aynı klasörde çalışıyorsa ccam, geçmişin karışabileceği konusunda uyarır ve `Run anyway? [y/N]` diye sorar. Klasör, sembolik bağlantıları çözülmüş olarak eşleştirilir. Terminal yoksa (`/dev/null` üzerinden gelen standart girdi de yok sayılır) ya da arka plandaki bir işten çalıştırıldığında uyarıyı gösterir ve devam eder. Yeni bir `-p` çalıştırması kendi konuşmasını başlatır ve uyarı almaz; `--version`, `--help` ve `mcp`, `auth` veya `plugin` gibi konuşma açmayan alt komutlar hiç denetlenmez.
- `--dry-run`; MCP eşitlemesini, yeni sunucu bildirimini, `mcp list` notunu, geçmiş kontrolünü, oturum işaretini, başlığı ve son kullanım zamanını atlar ve exec yapmaz. Paylaşım bağlantılarını yine kontrol eder ve eksik bir bağlantıyı normal bir başlatmanın yapacağı gibi oluşturur.
- `ccam work rc`, `work` zaten çalışırken `.claude.json` dosyasına dokunmaz, bunu belirtir ve Remote Control'ü grubun sunucuları olmadan başlatır; bunları kopyalamak için oturum bittikten sonra yeniden çalıştırın. `--dry-run` ile mesajı yazdırır ama ne sorar ne de kopyalar.
- `mcp_strategy = "sync"` ile başarısız bir eşitleme `ccam: MCP sync failed (launching anyway): …` yazdırır ve Claude Code zaten sahip olduğu sunucularla başlar.

## İlgili

- [Kısayol komutları](/features/shortcuts)
- [ccam ile Remote Control kullanma](/guides/remote-control)
- [Hesapları yönetme](/features/accounts)
- [Komut başvurusu](/docs/commands)
- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [CLAUDE_CONFIG_DIR nedir?](/guides/claude-config-dir)
