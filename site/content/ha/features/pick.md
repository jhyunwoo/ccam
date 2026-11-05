---
title: Zaɓar asusu don rubutun umarni
description: ccam pick yana buga asusun da ya fi sarari don rubutun umarni, yana ba workers da aka fara tare asusai daban-daban, kuma yana fita da 75 idan babu sarari.
updated: 2026-11-05
order: 12
translation: machine
sourceHash: "0ec25752ace9"
---

`ccam pick` yana buga sunan asusun da aka shiga da shi wanda ya fi sarari, don rubutun umarni (script) ya iya fara Claude Code da shi. Workers da aka fara a lokaci ɗaya kowannensu yana samun asusu daban muddin isassun asusai suna da sarari, don haka rukunin ayyukan `claude -p` yana rarrabuwa a kan asusanka maimakon ya shanye ɗaya.

**Kujeru (seats) na Team da Enterprise.** Zaɓin asusu na atomatik (`ccam next`, `--fallback`, `ccam pick`, maɓallin `n` na dashboard kuma, da zarar ka kunna shi, [`auto_switch`](/features/next#auto-switch)) yana kuma sauyawa tsakanin kujerun ƙungiyar Claude Team ko Enterprise iri ɗaya, wanda na iya saɓa wa sharuɗɗan Anthropic. Ka karanta [Tambayoyin da ake yawan yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) kuma ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi.


## Amfani

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

Exit codes: `0` tare da suna a standard output, `75` (`EX_TEMPFAIL`) idan babu asusun da ke da sarari ko babu wanda aka shiga da shi, `1` don wasu kurakurai kamar flag mara kyau.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--json` | Buga asusun, amfaninsa na awa 5 da na mako, inda lambobin suka fito da dalilin da ya sa aka zaɓe shi | a kashe |
| `--model <model>` | Ƙirga tagar mako ta wannan model ɗin ma, kamar `opus` | babu |
| `--max-usage <1-100>` | Fifita asusan da ke ƙasa da wannan amfani | `auto_max_usage`, in ba haka ba 100 |
| `--strategy headroom\|expiring\|balance` | Mafi ƙarancin amfanin mako a cikin asusan da amfanin awa 5 ɗinsu bai zarce mafi ƙarancin amfanin awa 5 da maki 3 ba, sake saitin awa 5 mafi kusa ko mafi yawan kason mako da ya rage ga kowace rana | `auto_strategy`, in ba haka ba `headroom` |
| `--cross-org` | Zaɓi asusan wasu ƙungiyoyi (organizations) ma | `cross_org`, in ba haka ba a kashe |
| `--wait[=<duration>]` | Idan babu asusun da ke da sarari, jira na farko ya sake saitawa, har zuwa iyakar wannan tsawon. Ƙimar tana zuwa bayan `=`: ana ƙin `--wait 30m` da shawarar `--wait takes its value after '=': --wait=30m` | a kashe; `5h` idan an ba shi shi kaɗai |
| `--dry-run` | Zaɓa ba tare da rubuta zaɓin ba | a kashe |

## Yadda yake aiki

1. ccam yana lissafa asusan da aka shiga da su (asusan API key da shigar Console ba su da iyaka kuma ana barin su a waje) a cikin ƙungiyar asusun da wannan foldar za ta fara (duba [ƙungiyoyi](/features/next#auto)); ana barin foldar da ta shiga asusun Claude iri ɗaya da wannan a waje ita ma.
2. Yana ɗaukar kulle (`~/.local/state/ccam/picks.lock`), don zaɓuɓɓukan da ke gudana a lokaci ɗaya su faru ɗaya bayan ɗaya.
3. Yana neman amfanin kowane asusu a lokaci ɗaya, kamar yadda `ccam usage` ke yi, tare da cache na daƙiƙa 60. Tagar da Claude Code ya ba da rahoton ta cika ana ƙirga ta a matsayin 100% har zuwa sake saitinta ([rikodin iyaka](/features/usage#the-limit-record)), don haka ana cire asusun da ya kai iyakarsa yanzu ko da lambobinsa da aka ajiye har yanzu suna nuna sarari.
4. Ana fara barin asusan da aka zaɓa a cikin minti 2 da suka gabata, ta `ccam pick` ko ta [`auto_switch`](/features/next#auto-switch) lokacin da ya motsa zama zuwa can (`~/.local/state/ccam/picks/<name>`; zaɓin da ke da lokaci na gaba (future timestamp) ba ya ƙirguwa) a waje, kuma [ƙa'idodin zaɓi](/features/next#auto) na `ccam next --auto` suna zaɓa a cikin sauran: waɗanda ke ƙasa da `auto_max_usage` da farko, bisa dabarar, kuma ba a taɓa zaɓar asusun da ya kai 100% ba. Idan wannan ya tsallake asusun da ƙa'idodin za su zaɓa in ba haka ba, `reason` na `--json` yana faɗin haka: `30% of its 5 hours used; work was picked a moment ago`.
5. Idan an zaɓi kowane asusu mai sarari kwanan nan, ƙa'idodin suna sake zaɓa a cikin dukansu, don haka da asusai uku, worker na huɗu yana raba ɗaya daga cikinsu.
6. Ana rubuta zaɓin kuma a saki kullen, sannan a buga sunan.

Tagar minti 2 tana da muhimmanci domin amfanin sabon aiki yana ɗaukar ɗan lokaci kafin ya bayyana a lambobin: in babu ita, workers biyar da aka fara tare duk za su ga asusu ɗaya a matsayin wanda ya fi sarari.

## Tsaro da yanayi na musamman

- `ccam pick` yana karanta amfani ne kawai kuma yana rubuta lokacin zaɓi; ba ya taɓa sabunta ko karanta token fiye da abin da `ccam usage` ke yi.
- Zaɓi ba ya keɓe asusun. Rubutun umarni biyu za su iya fara a asusu ɗaya bayan minti 2; yi amfani da `--fallback` don farfaɗowa idan wannan asusun ya kai iyakarsa.
- Jira da `--wait` ba ya riƙe kulle, don haka sauran zaɓuɓɓuka suna ci gaba a halin yanzu.
- `ccam pick` yana watsi da [`ccam folder`](/features/folder): yana zaɓa bisa amfani kawai.

## Masu alaƙa

- [Sauyawa zuwa asusu na gaba](/features/next)
- [Sake gudanarwa a iyaka](/features/next#run-again-at-the-limit)
- [Iyakokin amfani](/features/usage)
- [Asusan API key](/features/api)
