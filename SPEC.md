# agent-help SPEC v0.1

Status: Draft  
Audience: CLI authors, tool authors, validators, and coding agents implementing agent-help  
Canonical URL: https://zate.github.io/agent-help/SPEC.md  
Repository: https://github.com/Zate/agent-help

This is the technical contract for agent-help. It exists so implementers can add two small CLI surfaces consistently:

- `--agent-help` for compact command discovery.
- `--agent-out` for dense structured runtime results.

Agents using an already-enabled CLI should not need this file. They should discover the CLI from the CLI itself.

## Table of Contents

1. Core Idea
2. Terminology
3. CLI Flags
4. AHF Line Format
5. Record Prefix Registry
6. Scalar Type Registry
7. `--agent-help` Output
8. `--agent-out` Output
9. Errors
10. Pagination and Follow-Up
11. Accuracy Requirements
12. Security and Privacy
13. Versioning
14. Examples
15. Implementation Notes
16. Open Questions

## 1. Core Idea

agent-help adds an agent-facing surface beside existing human and software surfaces:

```text
human:    tool --help
software: tool --json
agent:    tool --agent-help
agent:    tool subcmd args --agent-out
```

`--agent-help` and the `--agent-out` protocol envelope use AHF: a small line-record format. Structured result bodies in `--agent-out` SHOULD use [TOON](https://github.com/toon-format/spec).

AHF is not a general serialization format and does not replace JSON. It is the discovery, invocation, status, and follow-up layer for agents talking to CLIs.

## 2. Terminology

- **agent-help**: the overall two-flag convention.
- **AHF**: Agent Help Format; line records used by `--agent-help` and the `--agent-out` envelope.
- **AH1**: root command index returned by `tool --agent-help`.
- **AH2**: command detail returned by `tool subcmd --agent-help`.
- **producer**: a CLI that emits agent-help output.
- **consumer**: an agent or runtime reading agent-help output.
- **TOON**: Token-Oriented Object Notation, recommended for `--agent-out` result bodies.

Requirement words such as MUST, SHOULD, and MAY describe interoperability expectations. Use common sense: the goal is small, accurate, easy-to-consume output.

## 3. CLI Flags

### 3.1 Human help breadcrumb

Human `--help` SHOULD include this breadcrumb or an equivalent pointer:

```text
LLM agent? Use --agent-help for token-optimized usage.
```

### 3.2 `--agent-help`

`--agent-help` is documented as a trailing global flag:

```text
tool --agent-help
tool subcmd --agent-help
tool group subcmd --agent-help
```

Do not document leading placement as the primary form:

```text
tool --agent-help subcmd
```

Implementations MAY accept leading placement for compatibility, but examples and docs SHOULD teach the trailing form.

### 3.3 `--agent-out`

`--agent-out` requests agent-optimized runtime output:

```text
tool subcmd args --agent-out
```

When `--agent-out` is present, the command SHOULD emit an AHF status envelope followed by a TOON result body when structured data exists.

## 4. AHF Line Format

AHF is UTF-8 plain text. ASCII is preferred when practical.

Rules:

- One record per line.
- First whitespace-delimited token is the record prefix.
- Use one space between tokens when emitting.
- Empty lines SHOULD be ignored by consumers.
- Unknown prefixes SHOULD be tolerated unless strict validation was requested.
- Use `_` for null, unknown, missing, or not applicable.
- Use `|` for short lists inside one field value.
- Avoid color, markdown tables, decorative headings, comments, and prose paragraphs.

Metadata fields use `key=value`:

```text
ok nodes count=8 shown=3 more=1 cursor=c_n_104
err invalid_enum flag=--type got=note
```

Purpose text in help records uses ` :: `:

```text
cmd node list [--limit int] :: list memory nodes
flag --limit:int opt default=10 :: max rows
```

Values containing whitespace, quotes, commas, or reserved delimiters SHOULD be double-quoted. Within quoted strings, escape quotes and backslashes with backslash.

```text
next inspect "tool search \"billing db\" --agent-out"
```

## 5. Record Prefix Registry

| Prefix | Meaning | Context |
|---|---|---|
| `ah1` | agent-help index | root `--agent-help` |
| `ah2` | agent-help detail | subcommand `--agent-help` |
| `ok` | success status | `--agent-out` envelope |
| `err` | error status | any agent-facing error |
| `warn` | non-fatal warning | `--agent-out` envelope |
| `cmd` | command entry | AH1 |
| `use` | canonical invocation | AH2 or error |
| `arg` | positional argument definition | AH2 |
| `flag` | flag definition | AH2 |
| `ex` | valid example invocation | AH2 |
| `hint` | direct correction | after `err` |
| `next` | follow-up or continuation command | after `ok` or `err` |
| `more?` | pointer to AH2 detail; not a shell command | AH1 |

The `?` in `more?` is intentional. It signals “want more information?” and helps prevent agents from treating the line as a shell command.

## 6. Scalar Type Registry

Use these scalar types in `arg` and `flag` records:

| Type | Meaning |
|---|---|
| `str` | free text string |
| `int` | integer |
| `num` | float |
| `bool` | true or false |
| `path` | filesystem path |
| `url` | URL |
| `id` | opaque identifier |
| `ts` | Unix timestamp |
| `date` | ISO 8601 date |
| `dur` | duration string |
| `kv` | key:value pair |
| `enum(a\|b)` | one of the listed values |

## 7. `--agent-help` Output

### 7.1 AH1: command index

Returned by `tool --agent-help`.

```text
ah1 <tool> :: <purpose>
cmd <command-signature> :: <purpose>
cmd <group> <command-signature> :: <purpose>
more? <tool> <cmd> --agent-help
```

Rules:

- Start with one `ah1` header.
- Include one `cmd` line per primary invocable command.
- Flatten nested commands into command paths.
- Put required args inline in `<angle brackets>`.
- Put only highest-value optional flags in `[square brackets]`.
- Omit aliases, author, version banners, prose, examples, and full flag details.
- Include `more?` to show how to request AH2 detail.
- Target fewer than 300 tokens.

### 7.2 AH2: command detail

Returned by `tool subcmd --agent-help`.

```text
ah2 <tool> <command-path>
use <canonical invocation>
arg <name>:<type> <req|opt> :: <purpose>
flag --<name>:<type> <req|opt|repeat> [default=<value>] :: <purpose>
ex <valid example invocation>
```

Rules:

- Start with one `ah2` header.
- Include one `use` record.
- Include all required args and required flags.
- Include optional flags that materially change behavior.
- Mark repeatable flags with `repeat`.
- State defaults only when useful or non-obvious.
- Every `ex` SHOULD work as written; omit fragile examples.
- Target fewer than 150 tokens.

## 8. `--agent-out` Output

`--agent-out` output has two layers:

1. AHF envelope: `ok`, `err`, `warn`, `next`, `hint` records.
2. Result body: SHOULD be TOON for structured lists/objects.

Successful output starts with `ok`:

```text
ok <kind> [key=value...]
```

Common metadata:

- `count=<n>` total known count
- `shown=<n>` number shown in this output
- `more=0|1` whether another page exists
- `cursor=<id>` pagination cursor
- `truncated=1` output was truncated

List result example:

```text
ok issues count=2 more=0
issues[#2]{key,status,summary}:
  PROJ-1,open,Add agent-help
  PROJ-2,done,Document examples
next inspect "tool issue PROJ-1 --agent-out"
```

Single-object results SHOULD still include stable identifiers and useful follow-up commands.

If there is no structured body, simple key-value AHF lines MAY follow `ok`:

```text
ok deploy
env prod
status healthy
next logs "tool logs --env prod --agent-out"
```

Use this fallback only when TOON would add no value.

Warnings are non-fatal and SHOULD appear immediately after `ok`:

```text
ok nodes count=8 shown=3 more=1
warn truncated shown=3 total=8
```

## 9. Errors

Agent-facing errors start with `err`:

```text
err <code> [key=value...]
```

Errors SHOULD include enough information for an agent to retry without reading human `--help`.

Use `hint` and `use` for direct corrections:

```text
err invalid_enum flag=--type got=note
hint --type enum(decision|fact|pattern|task|observation)
use mem node add <text> --type TYPE [--tag K:V...]
```

Rules:

- Do not respond only with “run --help”.
- For missing values, name the exact missing arg or flag.
- For enum errors, list valid values.
- Include `use` when it prevents another discovery call.
- Include `next` when there is an exact recovery or inspection command.

## 10. Pagination and Follow-Up

Use `next` for exact follow-up commands:

```text
next "tool list --cursor c_123 --agent-out"
next inspect "tool show n_102 --agent-out"
```

If output is paginated, truncated, or partial:

- include `more=1` or `truncated=1` on `ok`; and
- include an exact continuation command in `next`.

## 11. Accuracy Requirements

Producer output SHOULD match real CLI behavior.

Minimum checks:

- `tool --help` includes the breadcrumb.
- `tool --agent-help` emits AH1.
- Every AH1 `cmd` has AH2 output.
- AH2 required args and flags match validation.
- AH2 examples run successfully or are omitted.
- Structured result commands emit `ok` or `err` first under `--agent-out`.
- Paginated or truncated output includes `next`.
- Agent-facing output avoids colors, markdown tables, and prose paragraphs.

Generate AHF from existing command metadata where practical. Do not maintain a stale parallel command list if the framework can provide the data.

## 12. Security and Privacy

Agent-facing output may be logged by agent runtimes. Producers SHOULD:

- redact secrets, tokens, credentials, and keys by default;
- avoid unnecessary personal data;
- avoid leaking internals in retryable errors;
- treat `--agent-out` as potentially less private than terminal output.

## 13. Versioning

Current draft: v0.1.

Consumers SHOULD tolerate unknown record prefixes and metadata keys. Producers SHOULD keep field names and order stable when possible.

Optional future versions may add record prefixes, scalar types, or aliases without breaking existing consumers.

## 14. Examples

### 14.1 AH1 index

```text
ah1 mem :: project memory — store and query facts, decisions, and tasks
cmd node add <text> --type TYPE [--tag K:V...] :: store a new memory node
cmd node list [--type TYPE] [--limit int] [--cursor id] :: list memory nodes
cmd search query <text> [--type TYPE] [--limit int] [--cursor id] :: search nodes by text
cmd project status :: show current project settings
more? mem <cmd> --agent-help
```

### 14.2 AH2 detail

```text
ah2 mem node add
use mem node add <text> --type TYPE [--tag K:V...]
arg text:str req :: node text content
flag --type:enum(decision|fact|pattern|task|observation) req :: node type
flag --tag:kv repeat :: metadata key:value
ex mem node add "postgres 15 required" --type fact --tag project:mem
```

### 14.3 Error hint

```text
err invalid_enum flag=--type got=note
hint --type enum(decision|fact|pattern|observation)
use mem node add <text> --type TYPE
```

### 14.4 Runtime list

```text
ok nodes count=8 shown=3 more=0
nodes[#3]{id,type,tags,text}:
  n_102,fact,"project:billing|db",postgres 15 required for billing service
  n_088,decision,project:billing,migrate billing read model to postgres
  n_061,pattern,db|ops,use connection pool max 20 in all environments
next inspect "mem search similar n_102 --agent-out"
```

### 14.5 Runtime object

```text
ok project
project[#5]{key,value}:
  name,agent-help
  version,0.1.0
  status,draft
  owner,team-cli
  spec,SPEC.md
next nodes "mem node list --agent-out"
```

### 14.6 Paginated list with warning

```text
ok nodes count=8 shown=3 more=1 cursor=c_n_104
warn truncated shown=3 total=8
nodes[#3]{id,type,tags,text}:
  n_101,decision,"project:mem|arch",Use Cobra for the reference CLI implementation
  n_102,fact,"project:mem|db",Mock data lives in data.go
  n_103,pattern,db|ops,Use connection pool max 20 in all environments
next "mem node list --limit 3 --cursor c_n_104 --agent-out"
```

## 15. Implementation Notes

Framework pattern:

- Add hidden global/persistent flags for `--agent-help` and `--agent-out`.
- Intercept `--agent-help` before normal command execution.
- Emit AH1 for root and AH2 for the selected subcommand.
- Bypass required-arg validation when serving `--agent-help`.
- Derive command paths, args, flags, defaults, and enum values from framework metadata where practical.
- Add explicit metadata for scalar types and examples if the framework does not store them.

Common framework hooks:

- Cobra: persistent hidden flags; derive from `Command.Use`, flags, args, annotations.
- Click: hidden eager option; derive from `Command.params` and context path.
- argparse: hidden flag on root and subparsers; keep a small metadata table if needed.
- Clap: hidden global arg; inspect `Command` and matches before dispatch.
- Commander: global option and `preAction` hook; derive from command metadata.

## 16. Open Questions

- Should `--agent-out` have a shorter alias?
- Should AHF define an optional extension namespace?
- Should there be an official model comprehension test suite for examples?
- Should a media type such as `text/ahf` be registered later?
