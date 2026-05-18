# agent-help Example 04: --agent-out list results

Source command:

```bash
mem node list --agent-out
```

Raw output (AHF envelope + TOON body):

```text
ok nodes count=8 shown=8 more=0
nodes[#8]{id,type,tags,text}:
  n_101,decision,"project:mem|arch",Use Cobra for the reference CLI implementation
  n_102,fact,"project:mem|db",Mock data lives in data.go; no real storage in the reference impl
  n_103,pattern,db|ops,Use connection pool max 20 in all environments
  n_104,fact,"project:mem|api",AHF ok/err envelope always precedes the TOON body
  n_105,decision,"project:mem|spec",TOON is the encoding for --agent-out result bodies
  n_106,task,"project:mem|todo",Write reference Go implementation for agent-help
  n_107,pattern,agent|cli,Always include next records for paginated results
  n_108,fact,"project:mem|spec",more? record is a pointer not a shell command
next inspect "mem search query <text> --agent-out"
```

Notes:

- `ok` is the AHF protocol status line. `count=` and `more=` are envelope metadata — always AHF.
- The TOON block (`nodes[#8]{...}:` followed by indented rows) is the result body.
- `[#8]` is the TOON length marker — the `#` prefix is toon-go's output for dynamic-length arrays.
- Rows are two-space indented per TOON spec. Values are quoted only when they contain commas or spaces.
- `next` is an AHF follow-up record — exact command the agent should call to act on this result.
- TOON handles the structured data; AHF handles the protocol wrapper and follow-up.
