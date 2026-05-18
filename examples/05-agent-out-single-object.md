# agent-help Example 05: --agent-out single object

Source command:

```bash
mem project status --agent-out
```

Raw output (AHF envelope + TOON body):

```text
ok project
project[#6]{key,value}:
  name,agent-help
  version,0.1.0
  status,draft
  owner,team-cli
  spec,AHF-RFC.md
  repo,github.com/Zate/agent-help
next nodes "mem node list --agent-out"
```

Notes:

- `ok project` is the AHF status line with no envelope metadata — this is a single object, not a list.
- The TOON block encodes a vertical key/value table — one row per setting.
- This is a valid TOON pattern for structured settings or config objects.
- `next` gives the agent a natural follow-up action.

---

Contrast: bare AHF key-value fallback (§12.4)

When a result has no structured tabular body — a bare acknowledgement, a scalar status — AHF
key-value lines may be used directly after `ok`, without a TOON block:

```text
ok node
id n_109
type fact
tags spec:ahf
created 2026-05-12
text "toon-go library works"
next list "mem node list --agent-out"
```

Use this form only when there is no tabular or object data to encode. For any list or structured
object result, TOON is preferred.
