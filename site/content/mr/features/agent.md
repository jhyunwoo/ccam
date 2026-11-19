---
title: इतर खात्यांवरील subagents
description: ccam agent Claude Code च्या कामाचा एक भाग तुमच्या दुसऱ्या खात्याने चालवते आणि Claude चे उत्तर परत देते, त्यामुळे एक लांब काम अनेक खात्यांचा वापर करू शकते.
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Claude Code चे स्वतःचे subagents त्यांना सुरू करणाऱ्या सत्राच्या खात्याने चालतात, त्यामुळे ते त्याच खात्याचा वापर खर्च करतात. `ccam agent "<task>"` त्याच फोल्डरमध्ये तुमच्या दुसऱ्या खात्याने `claude -p` सह काम चालवते आणि Claude चे शेवटचे उत्तर छापते. Claude Code सत्रातून बोलावले असता ते दुसऱ्या खात्याचा वापर खर्च करणाऱ्या subagent सारखे काम करते, त्यामुळे सत्राचे स्वतःचे खाते संपले असते त्यानंतरही एक लांब काम पुढे चालू राहू शकते.

**Team आणि Enterprise सीट्स.** `ccam agent` खाते स्वतःच निवडते. एकाच Claude Team किंवा Enterprise organization मध्ये ते काम दुसऱ्या सीटवर चालवते, आणि यामुळे Anthropic च्या अटींचे उल्लंघन होऊ शकते. त्यावर अवलंबून राहण्यापूर्वी [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) वाचा आणि तुमच्या संस्थेकडे आणि Anthropic कडे चौकशी करा.

## वापर

Claude Code ते त्याच्या Bash tool मधून चालवते. [ccam plugin](/docs/install#claude-code-plugin) install केलेले असेल, तर साध्या शब्दांत सांगा, उदाहरणार्थ "या कामाच्या भागांसाठी माझी इतर खाती वापर", किंवा ते तुमच्या CLAUDE.md मध्ये लिहा. Plugin चे skill Claude Code ला प्रत्येक agent पार्श्वभूमीत सुरू करायला, इतर संदर्भाशिवाय समजेल असे काम लिहायला, आणि फाइल्स बदलणाऱ्या agents ना त्यांचा स्वतःचा भाग किंवा स्वतःचा git worktree द्यायला सांगते.

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

Claude चे शेवटचे उत्तर standard output वर जाते; `ccam:` ने सुरू होणाऱ्या ओळी standard error वर जातात. `--json` त्याऐवजी `{account, session_id, result, is_error, accounts}` छापते, जिथे `accounts` मध्ये काम ज्या ज्या खात्यांनी चालले ती सर्व खाती असतात.

Exit codes: claude चा स्वतःचा (agent पूर्ण झाल्यावर `0`), बोलावणाऱ्या खात्याशिवाय इतर कोणत्याही खात्यात जागा नसेल तर `75` (`EX_TEMPFAIL`), आणि agent ने दुसरा agent सुरू करण्याचा प्रयत्न करणे यासारख्या इतर त्रुटींसाठी `1`.

## पर्याय

| पर्याय | काय करते | डीफॉल्ट |
| --- | --- | --- |
| `--model <model>` | agent चे model; निवडीत त्या model ची साप्ताहिक window सुद्धा मोजली जाते | खात्याचे डीफॉल्ट model |
| `--resume <id>` | आधीच्या agent चे संभाषण पुढच्या कामासह चालू ठेवा (त्याच्या पहिल्या आणि शेवटच्या ओळीवरचा session id) | बंद |
| `--json` | खाते, सत्र, उत्तर आणि ते ज्या खात्यांनी चालले ती खाती छापा | बंद |
| `--max-usage <1-100>` | या वापराखालील खात्यांना प्राधान्य द्या | `auto_max_usage`, अन्यथा 100 |
| `--strategy headroom\|expiring\|balance` | [निवडीचा नियम](/features/next#auto) | `auto_strategy`, अन्यथा `headroom` |
| `--cross-org` | इतर organizations च्या खात्यांनीही चालवा | `cross_org`, अन्यथा बंद |
| `--wait[=<duration>]` | कोणत्याही खात्यात जागा नसेल तर पहिले खाते reset होण्याची वाट पाहा, जास्तीत जास्त इतका वेळ | बंद; एकटे दिल्यास `5h` |
| `--dry-run` | ते जे खाते निवडेल त्यासाठीची कमांड न चालवता छापा | बंद |

## हे कसे काम करते

1. **बोलावणारे खाते वगळले जाते.** ज्या खात्यातून `ccam agent` बोलावले जाते (ccam ने सुरू केलेल्या सत्राचे खाते, नाहीतर हा फोल्डर जे खाते सुरू करतो ते) ते कधीच agent चालवत नाही, आणि त्याच Claude खात्यात लॉग इन असलेला दुसरा फोल्डरही नाही, कारण तो त्याच मर्यादा शेअर करतो. `--cross-org` वेगळे सांगत नाही तोपर्यंत agent बोलावणाऱ्या खात्याच्या [organization](/features/next#auto) मध्येच राहतो.
2. **ते `ccam pick` सारखे निवडते.** वापर [`ccam pick`](/features/pick) प्रमाणेच विचारला जातो, आणि एकमेकांपासून 2 मिनिटांच्या आत सुरू झालेल्या agents ना, पुरेशा खात्यांत जागा असेपर्यंत, वेगवेगळी खाती मिळतात.
3. **ते सत्राच्या परवानग्यांनी चालते, त्याहून जास्त कधीच नाही.** `claude -p` मध्ये tool ला परवानगी द्यायला कोणीही नसते, त्यामुळे agent ला त्याला बोलावणाऱ्या सत्राचा permission mode आणि effort level मिळतो. terminal मधून चालवल्यास Claude Code चा डीफॉल्ट mode आणि तुमच्या settings मधील परवानगीचे नियम लागू होतात. `ccam agent` मध्ये जास्त परवानगी देणारा कोणताही पर्याय नाही.
4. **वापर मर्यादेवर संभाषण पुढे जाते.** agent चे खाते आपली मर्यादा गाठते तेव्हा ccam ती मर्यादा [plugin hook](/features/usage#the-limit-record) नोंदवेल तशीच नोंदवते आणि तेच संभाषण सर्वाधिक जागा असलेल्या खात्याने पुढे चालू ठेवते: Claude Code ते तिथे कॉपी करते (`--resume <file> --fork-session`), आणि पहिला संदेश असतो "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." प्रत्येक खाते एकदाच वापरून पाहिले जाते.
5. **agent चे प्रत्येक संभाषण नोंदवले जाते** `~/.local/state/ccam/agent/<session id>.json` मध्ये (30 दिवस). `--resume <id>` ते त्याच्या फोल्डरमध्ये पुढे चालू ठेवते, ज्या खात्याने ते चालवले त्या खात्यात अजून जागा असेल तर त्याने, नाहीतर दुसऱ्या खात्याने. [`ccam next --continue`](/features/next) agent संभाषणे वगळते, त्यामुळे ते तुमचीच संभाषणे पुढे चालू ठेवत राहते.

## सुरक्षितता आणि अपवादात्मक स्थिती

- agent स्वतःचे agents सुरू करू शकत नाही: ccam त्याला `CCAM_AGENT=1` ने चिन्हांकित करते, आणि तिथे `ccam agent` चालायला नकार देते.
- `ccam agent` थांबवल्यास (Ctrl+C, signal, किंवा Claude Code ने कमांड थांबवल्यास) agent सुद्धा थांबतो.
- agent ला त्याला बोलावणारे संभाषण दिसत नाही, फक्त त्याचे काम दिसते.
- एकाच वेळी त्याच फाइल्स बदलणारे agents एकमेकांचे काम overwrite करतात. प्रत्येकाला कामाचा स्वतःचा भाग द्या, किंवा स्वतःचा git worktree द्या.
- बोलावणारे सत्र दुसऱ्या खात्यावर जाते तेव्हा ([`auto_switch`](/features/next#auto-switch)), त्याने सुरू केलेले agents त्याच्यासोबत थांबतात. त्यांची संभाषणे राहतात: `--resume <id>` ती पुढे चालू ठेवते.
- `agent` हे कमांडचे नाव आहे, त्यामुळे कोणत्याही नवीन खात्याचे नाव `agent` ठेवता येत नाही; आधीपासून असलेले खाते अजूनही `ccam run agent` ने सुरू होते.

## संबंधित

- [scripts साठी खाते निवडणे](/features/pick)
- [पुढच्या खात्यावर जाणे](/features/next)
- [वापर मर्यादा](/features/usage)
- [प्रत्येक खात्याची पार्श्वभूमीतील सत्रे](/features/agents)
