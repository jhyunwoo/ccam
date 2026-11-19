---
title: பிற கணக்குகளில் subagents
description: ccam agent ஒரு Claude Code பணியின் பகுதியை மற்றொரு கணக்காக இயக்கி Claude-இன் பதிலைத் தருகிறது, அதனால் நீண்ட பணி பல கணக்குகளின் பயன்பாட்டைச் செலவிடுகிறது.
updated: 2026-11-19
order: 17
translation: machine
sourceHash: "c08d3cce6103"
---

Claude Code-இன் சொந்த subagents அவற்றைத் தொடங்கும் அமர்வின் கணக்காகவே இயங்குகின்றன, அதனால் அவை அந்தக் கணக்கின் பயன்பாட்டையே செலவழிக்கின்றன. `ccam agent "<task>"` அதே கோப்புறையில் உங்கள் மற்றொரு கணக்காக `claude -p` மூலம் ஒரு பணியை இயக்கி, Claude-இன் கடைசிப் பதிலை அச்சிடுகிறது. ஒரு Claude Code அமர்விலிருந்து அழைக்கப்படும்போது, அது மற்றொரு கணக்கின் பயன்பாட்டைச் செலவழிக்கும் subagent போலச் செயல்படுகிறது, அதனால் அமர்வின் சொந்தக் கணக்கு தீர்ந்துபோயிருக்கும் நிலையைத் தாண்டியும் ஒரு நீண்ட பணி தொடர முடியும்.

**Team மற்றும் Enterprise இருக்கைகள்.** `ccam agent` தானே கணக்கைத் தேர்ந்தெடுக்கிறது. ஒரே Claude Team அல்லது Enterprise நிறுவனத்துக்குள், அது பணியை மற்றொரு இருக்கையில் இயக்குகிறது, இது Anthropic-இன் விதிமுறைகளை மீறக்கூடும். அதை நம்பியிருக்கும் முன், [அடிக்கடி கேட்கப்படும் கேள்விகளைப்](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) படித்து, உங்கள் நிறுவனத்திடமும் Anthropic-இடமும் கேட்டு உறுதிசெய்துகொள்ளுங்கள்.

## பயன்பாடு

Claude Code அதைத் தன் Bash tool மூலம் இயக்குகிறது. [ccam plugin](/docs/install#claude-code-plugin) நிறுவப்பட்டிருந்தால், எளிய சொற்களில் கேளுங்கள், எடுத்துக்காட்டாக "இந்தப் பணியின் பகுதிகளுக்கு என் மற்ற கணக்குகளைப் பயன்படுத்து", அல்லது அதை உங்கள் CLAUDE.md-இல் எழுதுங்கள். ஒவ்வொரு agent-ஐயும் பின்னணியில் தொடங்கவும், தனியாகப் படித்தாலே புரியும் பணியை எழுதவும், கோப்புகளைத் திருத்தும் agents-க்கு அவற்றுக்கென ஒரு பகுதியையோ தனி git worktree-ஐயோ தரவும் plugin-இன் skill Claude Code-இடம் சொல்கிறது.

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

Claude-இன் கடைசிப் பதில் standard output-க்குச் செல்கிறது; `ccam:` என்று தொடங்கும் வரிகள் standard error-க்குச் செல்கின்றன. `--json` அதற்குப் பதிலாக `{account, session_id, result, is_error, accounts}`-ஐ அச்சிடுகிறது; இதில் `accounts` பணி இயங்கிய ஒவ்வொரு கணக்கையும் பட்டியலிடுகிறது.

வெளியேறும் codes: claude-இன் சொந்த code (agent முடிந்தால் `0`), அழைத்த கணக்கைத் தவிர வேறு எந்தக் கணக்கிலும் இடம் இல்லாதபோது `75` (`EX_TEMPFAIL`), ஒரு agent மற்றொன்றைத் தொடங்க முயல்வது போன்ற பிற பிழைகளுக்கு `1`.

## Options

| Option | என்ன செய்யும் | இயல்புநிலை |
| --- | --- | --- |
| `--model <model>` | agent-இன் model; தேர்வில் அந்த model-இன் வாராந்திரக் கால அளவும் கணக்கில் கொள்ளப்படும் | கணக்கின் இயல்புநிலை model |
| `--resume <id>` | முந்தைய agent-இன் உரையாடலை ஒரு தொடர் பணியுடன் தொடரும் (அதன் முதல் மற்றும் கடைசி வரியில் உள்ள session id) | off |
| `--json` | கணக்கு, அமர்வு, பதில், அது இயங்கிய கணக்குகள் ஆகியவற்றை அச்சிடும் | off |
| `--max-usage <1-100>` | இந்தப் பயன்பாட்டுக்குக் கீழ் உள்ள கணக்குகளுக்கு முன்னுரிமை தரும் | `auto_max_usage`, இல்லையெனில் 100 |
| `--strategy headroom\|expiring\|balance` | [தேர்வு விதி](/features/next#auto) | `auto_strategy`, இல்லையெனில் `headroom` |
| `--cross-org` | பிற நிறுவனங்களின் கணக்குகளாகவும் இயங்கும் | `cross_org`, இல்லையெனில் off |
| `--wait[=<duration>]` | எந்தக் கணக்கிலும் இடம் இல்லாதபோது, முதல் கணக்கு மீட்டமைக்கப்படும் வரை, அதிகபட்சம் இவ்வளவு நேரம் காத்திருக்கும் | off; தனியாகக் கொடுத்தால் `5h` |
| `--dry-run` | அது தேர்ந்தெடுக்கும் கணக்குக்கான கட்டளையை இயக்காமல் அச்சிடும் | off |

## இது எப்படி வேலை செய்கிறது

1. **அழைக்கும் கணக்கு விலக்கப்படுகிறது.** `ccam agent` அழைக்கப்படும் கணக்கு (ccam தொடங்கிய அமர்வின் கணக்கு, இல்லையெனில் இந்தக் கோப்புறை தொடங்கும் கணக்கு) ஒருபோதும் agent-ஐ இயக்காது; அதே Claude கணக்கில் உள்நுழைந்து அதன் வரம்புகளைப் பகிரும் மற்றொரு கோப்புறையும் இயக்காது. `--cross-org` வேறுவிதமாகச் சொல்லாவிட்டால், agent அழைக்கும் கணக்கின் [நிறுவனத்துக்குள்](/features/next#auto) இருக்கும்.
2. **அது `ccam pick` போலத் தேர்ந்தெடுக்கிறது.** [`ccam pick`](/features/pick) செய்வது போலவே பயன்பாடு பார்க்கப்படுகிறது, ஒன்றுக்கொன்று 2 நிமிடங்களுக்குள் தொடங்கும் agents, போதுமான கணக்குகளில் இடம் இருக்கும் வரை, வெவ்வேறு கணக்குகளைப் பெறுகின்றன.
3. **அது அமர்வின் அனுமதிகளுடன் இயங்குகிறது, ஒருபோதும் அதற்கு மேல் அல்ல.** `claude -p`-இல் ஒரு tool-ஐ அனுமதிக்க யாரும் இல்லை, அதனால் agent தன்னை அழைத்த அமர்வின் permission mode-ஐயும் effort level-ஐயும் பெறுகிறது. ஒரு டெர்மினலிலிருந்து இயக்கினால், Claude Code-இன் இயல்புநிலை mode-உம் உங்கள் settings-இன் அனுமதி விதிகளும் பொருந்தும். `ccam agent`-இல் அதிக அனுமதி தரும் option எதுவும் இல்லை.
4. **பயன்பாட்டு வரம்பில், உரையாடல் தொடர்கிறது.** agent-இன் கணக்கு தன் வரம்பை எட்டும்போது, [plugin hook](/features/usage#the-limit-record) செய்வது போலவே ccam அந்த வரம்பைப் பதிவுசெய்து, அதிக இடம் உள்ள கணக்காக அதே உரையாடலைத் தொடர்கிறது: Claude Code அதை அங்கே நகலெடுக்கிறது (`--resume <file> --fork-session`), அதன் முதல் செய்தி "Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete." ஒவ்வொரு கணக்கும் ஒருமுறை மட்டுமே முயலப்படுகிறது.
5. **ஒவ்வொரு agent உரையாடலும் பதிவுசெய்யப்படுகிறது**, `~/.local/state/ccam/agent/<session id>.json`-இல் (30 நாட்கள்). `--resume <id>` அதை அதன் கோப்புறையில், அதை இயக்கிய கணக்கில் இன்னும் இடம் இருந்தால் அந்தக் கணக்காக, இல்லையெனில் மற்றொன்றாகத் தொடர்கிறது. [`ccam next --continue`](/features/next) agent உரையாடல்களை விலக்குகிறது, அதனால் அது உங்கள் உரையாடல்களையே தொடர்ந்து தொடர்கிறது.

## பாதுகாப்பும் விளிம்பு நிலைகளும்

- ஒரு agent தன் சொந்த agents-ஐத் தொடங்க முடியாது: ccam அதை `CCAM_AGENT=1` மூலம் குறிக்கிறது, அங்கே `ccam agent` இயங்க மறுக்கிறது.
- `ccam agent`-ஐ நிறுத்துவது (Ctrl+C, ஒரு signal, அல்லது Claude Code கட்டளையை இடைமறிப்பது) agent-ஐயும் நிறுத்துகிறது.
- agent தன்னை அழைத்த உரையாடலைப் பார்ப்பதில்லை, தன் பணியை மட்டுமே பார்க்கிறது.
- ஒரே கோப்புகளை ஒரே நேரத்தில் திருத்தும் agents ஒன்றின் வேலையை மற்றொன்று மேலெழுதிவிடும். ஒவ்வொன்றுக்கும் பணியின் தனிப் பகுதியையோ தனி git worktree-ஐயோ கொடுங்கள்.
- அழைக்கும் அமர்வு மற்றொரு கணக்குக்கு நகரும்போது ([`auto_switch`](/features/next#auto-switch)), அது தொடங்கிய agents அதனுடன் சேர்ந்து நின்றுவிடும். அவற்றின் உரையாடல்கள் இருக்கும்: `--resume <id>` அவற்றைத் தொடர்கிறது.
- `agent` ஒரு கட்டளையின் பெயர், அதனால் எந்தப் புதிய கணக்குக்கும் `agent` என்று பெயரிட முடியாது; ஏற்கெனவே உள்ள ஒன்று இன்னும் `ccam run agent` மூலம் தொடங்கும்.

## தொடர்புடையவை

- [scripts-க்கு ஒரு கணக்கைத் தேர்ந்தெடுத்தல்](/features/pick)
- [அடுத்த கணக்குக்கு மாறுதல்](/features/next)
- [பயன்பாட்டு வரம்புகள்](/features/usage)
- [எல்லாக் கணக்குகளின் பின்னணி அமர்வுகள்](/features/agents)
