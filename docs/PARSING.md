# AHF Parsing Notes

AHF is intentionally line-oriented. Consumers do not need a full parser to get value from it, but these rules are enough for robust handling.

## Line Model

- AHF output is UTF-8 text.
- Each record is one line.
- Empty lines should be ignored by consumers.
- The first whitespace-delimited token is the record prefix.
- Unknown record prefixes should be tolerated unless strict validation is requested.

## Record Families

Help records:

- `ah1`, `cmd`, `more?`
- `ah2`, `use`, `arg`, `flag`, `ex`

Runtime protocol records:

- `ok`, `err`, `warn`, `next`, `hint`

Runtime data:

- Structured data after `ok` should be TOON.
- Bare scalar acknowledgements may use simple AHF key-value lines after `ok`.

## `::` Separator

The `::` separator is used for human-readable purpose text in help records:

```text
cmd node list [--limit int] :: list memory nodes
flag --limit:int opt default=10 :: max results to return
```

Consumers can split on the first ` :: ` when they need the signature and purpose separately.

## Metadata Fields

Header records may include `key=value` metadata:

```text
ok nodes count=8 shown=3 more=1 cursor=c_n_104
err invalid_enum flag=--type got=badval
```

Values should be treated as strings unless the key has known semantics. For example, `more=1` and `truncated=1` indicate that a continuation `next` record is required.

## Quoting

Values with spaces or reserved delimiters should be double-quoted. Consumers should handle backslash-escaped quotes and backslashes inside quoted values.

```text
next inspect "mem search query \"decision\" --type decision --agent-out"
```

## Lists and Nulls

- `_` means null, unknown, missing, or not applicable.
- `|` separates short lists inside one field value.

```text
tags project:mem|spec
owner _
```

## TOON Handoff

For `--agent-out`, parse AHF protocol records first:

```text
ok nodes count=3 more=0
nodes[#3]{id,type,text}:
  n_101,fact,example text
next inspect "mem search similar n_101 --agent-out"
```

The `ok`, `warn`, and `next` lines are AHF. The `nodes[#3]{...}:` block is TOON.

## Minimal Consumer Strategy

1. Read lines.
2. Inspect the first token of each line.
3. If the first line is `err`, use `hint`, `use`, and `next` records to recover.
4. If the first line is `ok`, parse metadata such as `count`, `more`, and `cursor`.
5. Treat TOON blocks as the result body.
6. Use `next` records for pagination or recommended follow-up actions.

## Validator Scope

The repository includes `scripts/ahf-validate.py` as a lightweight example and fixture validator. It is intentionally not a full formal parser.

It checks common producer mistakes:

- unknown AHF prefixes
- missing AH1 `more?`
- missing AH2 `use`
- malformed `arg` and `flag` records
- unknown scalar types
- leading `--agent-help` examples
- `warn` records after TOON bodies
- paginated or truncated `ok` output without `next`

It does not fully validate:

- shell quoting semantics
- every possible TOON body shape
- all command invocation syntax
- semantic correctness of examples
- whether documented flags actually exist in a real CLI

Use it as a guardrail for examples and fixtures, not as the only conformance test for a production CLI.
