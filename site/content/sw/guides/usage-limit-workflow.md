---
title: Cha kufanya Claude Code ikifikia kikomo cha matumizi
description: Subiri kikomo kianze upya au endeleza mazungumzo yale yale kwenye akaunti nyingine ya Claude Code mara moja. Jinsi vikomo vya saa 5 na vya wiki hufanya kazi.
updated: 2026-11-01
order: 3
translation: machine
sourceHash: "760a36724feb"
---

Akaunti inapofikia kikomo chake cha matumizi (mara nyingi huitwa rate limit), una machaguo mawili: kusubiri hadi kikomo kianze upya, au kuendeleza mazungumzo yale yale sasa hivi kwenye akaunti nyingine unayoimiliki. Kwa ccam, chaguo la pili huchukua amri moja, `ccam next --auto --continue`.

## Claude Code ina vikomo gani?

Kwenye mipango ya Pro na Max, Claude Code ina [kikomo cha kipindi cha saa 5 na kikomo cha wiki](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Vyote viwili hushirikiwa na Claude kwenye wavuti na kwenye programu, na matumizi huhesabiwa kwenye vyote viwili kwa wakati mmoja. Kimoja kinapoisha, Claude Code [huonyesha ujumbe](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) kama "You've hit your session limit · resets 3:45pm". Kubadili model hakusaidii, kwa sababu vikomo vinahusu model zote.

## Chaguo la 1: subiri kikomo kianze upya

Ikiwa una akaunti moja, subiri. Matoleo ya hivi karibuni ya Claude Code [husubiri ndani ya kipindi kilicho wazi](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) na huendelea na kazi yenyewe mara kikomo kinapoanza upya. `/usage` ndani ya Claude Code huonyesha kiasi ulichotumia cha vikomo vya mpango wako.

## Chaguo la 2: endelea kwenye akaunti nyingine sasa

Ikiwa una akaunti ya pili, unaweza kuendelea na mazungumzo yale yale mara moja. Sanidi hili mara moja tu:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Sasa mstari wa hali (status line) huonyesha akaunti na matumizi yake ya saa 5 wakati unafanya kazi, kwa hiyo kikomo hakikushtukizi.

Kikomo kinapofikiwa, toka kwenye Claude Code (`/exit`), kisha:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Fungua mazungumzo katika akaunti yenye nafasi kubwa zaidi:

```sh
ccam next --auto --continue
```

Claude Code huanza kama `work` ikiwa na mazungumzo ya hivi karibuni zaidi ya folda ya sasa, pamoja na kila kitu kilichosemwa hadi sasa. Kushiriki historia (`ccam group set default --add history` hapo juu) humaanisha akaunti zote mbili hufanya kazi kwenye mazungumzo yale yale; bila hilo, Claude Code huyanakili mazungumzo kwenye `work`, na ccam hukuambia ni tokeni ngapi jibu la kwanza litalazimika kusoma tena.

Ukiwa na [plugin ya ccam](/docs/install#claude-code-plugin) ndani ya Claude Code, kipindi kinachosimama kwenye kikomo chake huonyesha arifa na hurekodi kipindi cha muda kilichojaa, kwa hiyo `ccam next --auto` haichagui akaunti hiyo tena kabla ya kuanza upya. `/ccam:switch` (au `/ccam:switch work`) hurekodi mazungumzo kwa ajili ya `ccam next --continue`, ili iendeleze mazungumzo hayo hasa.

**Viti (seats) vya Team na Enterprise.** Uchaguzi wa akaunti wa kiotomatiki (`ccam next`, `--fallback`, `ccam pick`, kitufe `n` cha dashibodi na, ukishakiwasha, [`auto_switch`](/features/next#auto-switch)) pia hubadili kati ya viti vya shirika lile lile la Claude Team au Enterprise, jambo ambalo linaweza kukiuka masharti ya Anthropic. Soma [Maswali yanayoulizwa mara kwa mara](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) na uwasiliane na shirika lako na Anthropic kabla ya kulitegemea.

## Je, ccam hubadili akaunti kiotomatiki?

Ni pale tu unapoiwasha.

- Kwa chaguo-msingi, ccam haiendeshwi wakati unafanya kazi: huanzisha Claude Code na kujibadilisha nayo, kwa hiyo wewe ndiye unaamua wakati wa kubadili. `ccam next --auto` huchagua akaunti yenye matumizi ya chini zaidi ya saa 5, huruka akaunti zilizofikia 100% au ambazo hazijaingiwa, na hubaki ndani ya shirika moja la Team au Enterprise, au miongoni mwa mipango ya binafsi.
- Kwa `relay = true` katika `config.toml`, ccam badala yake huendelea kuendeshwa kando ya kipindi, kwa hiyo `/ccam:switch` huendeleza mazungumzo yenyewe, ikihifadhi chaguo za kipindi, hali ya ruhusa, kiwango cha juhudi (effort) na agent; kabla haijayahamisha mazungumzo kwenye akaunti ya shirika jingine, ccam hukuuliza.
- Kwa `auto_switch = true` (Linux na macOS; ya majaribio), ccam pia hufanya ubadilishaji yenyewe: Claude Code inaposimamisha kipindi kwenye kikomo chake na akaunti nyingine ya aina hiyo hiyo ina nafasi, ccam humaliza kipindi na huendeleza mazungumzo huko ([masharti na gharama](/features/next#auto-switch)).

Baadhi ya zana nyingine hubadilisha taarifa za kuingia za kipindi kinachoendelea badala yake; [ccam na zana nyingine za kubadili akaunti](/compare/account-switchers) inaeleza tofauti.

## Mambo ya kukumbuka

- Endesha `ccam next --continue` katika folda ile ile uliyofanyia kazi; mazungumzo huhifadhiwa kwa kila folda ya mradi.
- Epuka kuendesha akaunti mbili kwenye mazungumzo yale yale kwa wakati mmoja. ccam hukuonya ukijaribu.
- Bila historia inayoshirikiwa, mazungumzo hunakiliwa kwenye akaunti nyingine na nakala hizo mbili huendelea kila moja kwa njia yake; checkpoints ambazo `/rewind` hurejesha msimbo kutoka kwazo hubaki kwenye akaunti ya kwanza.
- Ujumbe wa kwanza kwenye akaunti mpya unaweza kutumia sehemu kubwa kuliko kawaida ya kikomo chake, kwa sababu akaunti hiyo bado haijahifadhi mazungumzo kwenye cache. `ccam next --continue --brief` badala yake huianzisha kwa muhtasari mfupi.
- Iwapo unaruhusiwa kubadili kati ya akaunti kwa namna hii inategemea masharti yanayohusu kila akaunti.

Tazama [Vikomo vya matumizi](/features/usage) na [Kubadili kwenda akaunti inayofuata](/features/next) kwa kila chaguo.
