# Conformance Levels

agent-help conformance is split into levels so CLI authors can adopt the convention incrementally.

## Level 1: `agent-help-core`

A CLI conforms to `agent-help-core` when it implements agent-facing discovery.

Baseline expectations:

- Human `--help` includes this breadcrumb or an equivalent pointer:
  `LLM agent? Use --agent-help for token-optimized usage.`
- `tool --agent-help` emits AH1 output.
- `tool subcmd --agent-help` emits AH2 output for every primary subcommand.
- `--agent-help` is documented and supported as a trailing global flag.
- Agent-facing help uses registered AHF records only.
- AH1 command signatures and AH2 args/flags match actual command behavior.
- AH2 examples either run successfully or are omitted.

## Level 2: `agent-help-runtime`

A CLI conforms to `agent-help-runtime` when it implements `agent-help-core` plus agent-facing runtime output.

Baseline expectations:

- Structured result commands support `--agent-out`.
- `--agent-out` output begins with an AHF `ok` or `err` record.
- Structured success bodies use TOON.
- Runtime errors include enough information for an agent to retry without human `--help`.
- Non-fatal conditions use AHF `warn` records after `ok`.
- Follow-up actions use AHF `next` records.

## Level 3: `agent-help-full`

A CLI conforms to `agent-help-full` when it implements `agent-help-runtime` plus pagination, truncation, and verification discipline.

Baseline expectations:

- Paginated or truncated results include `more=1` or `truncated=1` metadata.
- Paginated or truncated results include an exact continuation command in a `next` record.
- Agent-facing outputs redact secrets and unnecessary personal data.
- The project tests AH1/AH2 output, AH2 examples, `--agent-out` success output, and representative `err` output.
- The implementation keeps `--agent-help` generated from command metadata where practical.

## Reference Implementation

The Go/Cobra `mem` demo in `impl/` is intended to demonstrate `agent-help-full` for a small CLI:

- AH1 and AH2 help surfaces.
- AHF `ok`/`err`/`warn` protocol records.
- TOON list and object bodies.
- Pagination with continuation `next` records.
- Golden-output tests for key surfaces.
