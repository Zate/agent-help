# agent-help Example 06: --agent-out paginated with warning

Source command:

```bash
mem node list --limit 3 --agent-out
```

Raw output (AHF envelope + TOON body):

```text
ok nodes count=8 shown=3 more=1 cursor=c_n_104
warn truncated shown=3 total=8
nodes[#3]{id,type,tags,text}:
  n_101,decision,"project:mem|arch",Use Cobra for the reference CLI implementation
  n_102,fact,"project:mem|db",Mock data lives in data.go; no real storage in the reference impl
  n_103,pattern,db|ops,Use connection pool max 20 in all environments
next "mem node list --limit 3 --cursor c_n_104 --agent-out"
next inspect "mem search query <text> --agent-out"
```

Notes:

- `ok` with `more=1` signals paginated/truncated output; include a `next` continuation.
- `count=8` is total available; `shown=3` is rows in this response; `cursor=` is the page token.
- `warn truncated` is an AHF record — non-fatal, appears after `ok`, before the TOON body.
- The TOON body carries only the rows for this page. `[#3]` reflects the actual row count.
- Two `next` records: one for the next page (unlabelled), one for an inspection action (labelled).
- The `cursor=` value from the `ok` envelope is passed through verbatim in the `next` command.
- An agent that sees `more=1` with no `next` record SHOULD treat this as a producer error.
