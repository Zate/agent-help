# Agent Help Format (AHF) v0.1

Status: Internet-Draft style proposal  
Category: Informational / Implementation Convention  
Intended audience: CLI authors, tool authors, agent runtime authors, LLM agents  
Canonical URL: https://zate.github.io/agent-help/AHF-RFC.md  
Repository: https://github.com/Zate/agent-help  
Date: 2026-05-11

## Abstract

Agent Help Format (AHF) is the help and discovery sub-specification of the **agent-help** convention. It defines how command-line tools expose their interface to LLM agents via `--agent-help`, and how they signal runtime results and errors via `--agent-out`.

AHF defines two agent-facing surfaces:

1. `--agent-help` — a trailing global flag for command discovery and invocation construction. Output uses the AHF line-record format defined in this document.
2. `--agent-out` — a runtime output mode signalling that the tool should emit agent-optimized output. Runtime result bodies SHOULD use [TOON (Token-Oriented Object Notation)](https://github.com/toon-format/spec). Protocol-level status lines (`ok`, `err`) and follow-up commands (`next`, `hint`) are defined here.

AHF is not a general-purpose serialization format, a replacement for JSON, or a competitor to TOON. It defines the discovery and invocation layer — the conversation between an agent and a CLI — not the encoding of data within that conversation.

The goal: any agent should be able to discover and invoke any agent-help conforming CLI directly. The CLI itself carries the discovery and invocation surface agents need, while MCP servers, skills, and plugins can still wrap that same CLI.

## Status of This Memo

This document is an early draft specification for discussion and implementation experimentation. It uses requirement keywords from RFC 2119 and RFC 8174 sparingly: hard requirements describe the wire shape needed for interoperability, while most implementation guidance is intentionally advisory. Registry entries, naming details, and conformance levels may change before a 1.0 version is published.

## Table of Contents

1. Introduction
2. Terminology
3. Requirement Keywords
4. Design Goals
5. Audience Model
6. CLI Discovery Convention
7. AHF Data Model
8. Lexical Syntax
9. Record Prefix Registry
10. Scalar Type Registry
11. `--agent-help` Output
12. `--agent-out` Convention
13. Error Handling
14. Pagination, Truncation, and Follow-Up Commands
15. Accuracy and Conformance
16. Security and Privacy Considerations
17. Accessibility and Internationalization
18. Versioning
19. Examples (19.1–19.7)
20. Implementation Guidance
21. Open Questions
22. Relationship to Other Formats

## 1. Introduction

Command-line tools commonly expose two output modes:

- human-readable output, optimized for people;
- machine-readable output, commonly JSON, optimized for software parsers.

LLM agents are neither traditional humans nor deterministic software parsers. Agents need concise, stable, low-ambiguity text that provides enough structure to construct correct commands, inspect results, recover from errors, and choose next actions. Human-oriented output often contains prose, alignment, color, headings, aliases, and examples that consume tokens without improving agent actionability. JSON is excellent for software integrations but contains punctuation and structural redundancy that is often unnecessary for language models.

The agent-help convention defines a small, composable set of CLI surfaces for agent-native interaction:

```text
human:    tool --help
software: tool --json
agent:    tool --agent-help
agent:    tool subcmd args --agent-out
```

AHF is the specification for the `--agent-help` output format and the `--agent-out` protocol envelope. For `--agent-out` result bodies, AHF delegates to TOON.

## 2. Terminology

**agent-help** — the overall convention and project name.  
**AHF** — Agent Help Format; the specific line-record specification for `--agent-help` output and the `--agent-out` protocol envelope.  
**AH1** — Agent Help Index; the top-level help record set returned by `tool --agent-help`.  
**AH2** — Agent Help Detail; the command-level help record set returned by `tool subcmd --agent-help`.  
**producer** — a CLI tool that implements the agent-help convention.  
**consumer** — an LLM agent or agent runtime that reads AHF output.  
**TOON** — Token-Oriented Object Notation; the recommended encoding for `--agent-out` result bodies. See https://github.com/toon-format/spec.

## 3. Requirement Keywords

The key words MUST, MUST NOT, REQUIRED, SHALL, SHALL NOT, SHOULD, SHOULD NOT, RECOMMENDED, MAY, and OPTIONAL in this document are to be interpreted as described in RFC 2119 and RFC 8174.

## 4. Design Goals

AHF is designed to achieve:

1. low token count;
2. unambiguous machine parseability without a formal grammar;
3. human readability as a secondary property;
4. derivability from existing CLI framework metadata;
5. graceful degradation when truncated mid-stream.

AHF explicitly does not attempt to achieve:

1. arbitrary nested data serialization (delegated to TOON);
2. binary encoding efficiency;
3. schema validation.

## 5. Audience Model

| Audience | Interface | Expected Format |
|---|---|---|
| Humans | default output / `--help` | prose, tables, color optional |
| Software | `--json` | strict JSON |
| Agents | `--agent-help`, `--agent-out` | AHF (help/protocol) + TOON (data) |

## 6. CLI Discovery Convention

### 6.1 Human Help Breadcrumb

Implementations SHOULD add the following breadcrumb, or an equivalent pointer, to normal human `--help` output:

```text
LLM agent? Use --agent-help for token-optimized usage.
```

### 6.2 `--agent-help` Flag Placement

`--agent-help` is designed as a trailing global flag. Implementations SHOULD accept the following positions:

```text
tool --agent-help
tool subcmd --agent-help
tool group subcmd --agent-help
```

The following forms are discouraged and SHOULD NOT be documented as the primary form:

```text
tool --agent-help subcmd
tool --agent-help group subcmd
```

Implementations MAY accept leading placement for compatibility, but SHOULD avoid teaching agents that form.

### 6.3 `--agent-out` Flag Placement

`--agent-out` requests agent-optimized runtime output. Implementations SHOULD support it as a global flag and SHOULD document the trailing form:

```text
tool subcmd args --agent-out
```

When `--agent-out` is present, result bodies SHOULD be emitted as TOON. The protocol envelope (`ok`/`err` status lines, `next`/`hint`/`warn` records) uses AHF regardless.

Tools MAY also support `--format toon` or equivalent aliases, but `--agent-out` is the RECOMMENDED flag name.

## 7. AHF Data Model

AHF `--agent-help` output is a sequence of records. Each record is one line. The first token on each line identifies the record type. The remainder of the line contains record-specific fields.

AHF supports two structural patterns:

1. metadata records using `key=value` fields;
2. key-value object lines for single-record output (used in AH2).

AHF does not define its own tabular row format for runtime data. Use TOON for `--agent-out` result bodies.

## 8. Lexical Syntax

### 8.1 Character Encoding

AHF output is UTF-8. ASCII-only output is RECOMMENDED when practical.

### 8.2 Lines

Each AHF record is emitted on a single line terminated by LF (`\n`) or CRLF (`\r\n`). Consumers SHOULD tolerate a final line without a trailing newline.

### 8.3 Tokens

Tokens are separated by one or more spaces. Producers SHOULD emit a single space between tokens.

### 8.4 Key-Value Metadata

Metadata fields use this form:

```text
key=value
```

Keys SHOULD use lowercase ASCII letters, digits, `_`, or `-`. Values SHOULD be unquoted when they do not contain whitespace or reserved delimiters.

### 8.5 Quoting

Values containing whitespace, quotes, or reserved delimiters SHOULD be quoted using double quotes.

```text
summary "Implement AHF support"
```

Within quoted strings, producers SHOULD escape double quotes and backslashes using backslash escapes.

```text
summary "agent said \"retry\""
```

Consumers SHOULD be liberal in accepting quoted values and SHOULD NOT fail the entire output if one quoted field cannot be decoded.

### 8.6 Null and Unknown Values

The underscore token (`_`) represents null, unknown, missing, or not applicable. Producers SHOULD prefer `_` over empty strings in positional fields.

### 8.7 Lists

Short lists inside one field SHOULD use the pipe character (`|`).

```text
tags project:billing|db|critical
```

Long lists SHOULD use a follow-up `next` command pointing to a TOON result.

### 8.8 Comments and Decorative Text

AHF output SHOULD avoid decorative headings, color escape codes, markdown formatting, prose paragraphs, and trailing comments. Producers aiming for strict interoperability should make every line a parseable record.

## 9. Record Prefix Registry

The following record prefixes are defined for AHF:

| Prefix | Meaning | Context |
|---|---|---|
| `ah1` | agent-help index | `--agent-help` on root |
| `ah2` | agent-help detail | `--agent-help` on subcommand |
| `ok` | success status | `--agent-out` protocol envelope |
| `err` | error | any |
| `warn` | non-fatal warning | `--agent-out` protocol envelope |
| `cmd` | command entry | AH1 |
| `use` | canonical invocation | AH2 |
| `arg` | positional argument definition | AH2 |
| `flag` | flag definition | AH2 |
| `ex` | valid example invocation | AH2 |
| `hint` | direct correction | after `err` |
| `next` | follow-up or continuation command | after `ok` or `err` |
| `more?` | pointer to AH2 detail; not a shell command | AH1 |

The `?` in `more?` is intentional. It makes the record visibly advisory rather than executable: it points to the shape of a follow-up `--agent-help` request, but agents should not copy the `more?` token into a shell command.

## 10. Scalar Type Registry

The following scalar types are used in `arg` and `flag` records:

| Type | Meaning |
|---|---|
| `str` | free text string |
| `int` | integer |
| `num` | float |
| `bool` | true / false |
| `path` | filesystem path |
| `url` | URL |
| `id` | opaque identifier |
| `ts` | Unix timestamp |
| `date` | ISO 8601 date |
| `dur` | duration string (e.g. `30s`, `5m`, `2h`) |
| `kv` | key:value pair |
| `enum(a\|b)` | one of the listed values |

## 11. `--agent-help` Output

### 11.1 AH1: Agent Help Index

AH1 is returned by `tool --agent-help`. It provides a compact index of all invocable commands.

Format:

```text
ah1 <tool> :: <purpose>
cmd <command-signature> :: <purpose>
cmd <group> <command-signature> :: <purpose>
more? <tool> <cmd> --agent-help
```

AH1 output starts with an `ah1` header. It SHOULD include at least one `cmd` record unless the tool has no subcommands, and SHOULD include a `more?` record describing how to request AH2 detail. The `more?` record is a pointer, not a shell command. The `?` is intentional and signals "want more information?".

Rules:

- One `cmd` line per invocable command.
- Flatten nested commands into command paths.
- Include required args inline in `<angle brackets>`.
- Include only highest-value optional flags in `[square brackets]`.
- No aliases, author, version, prose, examples, or flag details.
- Sort by likely agent usage when known.
- Target fewer than 300 tokens.

### 11.2 AH2: Agent Help Detail

AH2 is returned by `tool subcmd --agent-help`. It describes one invocable command.

Format:

```text
ah2 <tool> <command-path>
use <canonical invocation>
arg <name>:<type> <req|opt> :: <purpose>
flag --<name>:<type> <req|opt|repeat> [default=<value>] :: <purpose>
ex <valid example invocation>
```

AH2 output starts with an `ah2` header and normally includes a `use` record. It SHOULD include all required arguments and required flags. It SHOULD include optional flags that materially change behavior, and MAY omit obscure optional flags if including them would materially reduce agent comprehension.

Every `ex` record SHOULD be valid as written in the documented environment; omit examples that are likely to drift or fail.

AH2 output SHOULD target fewer than 150 tokens.

### 11.3 Accuracy of Help Output

`--agent-help` is the agent-facing source of truth. Implementations SHOULD derive AH1 and AH2 from the same command metadata used for dispatch when possible. Producers SHOULD NOT knowingly emit command signatures, flags, defaults, enum values, or examples that contradict actual behavior.

## 12. `--agent-out` Convention

### 12.1 Overview

`--agent-out` signals that the tool should emit agent-optimized output. The output consists of two layers:

1. **Protocol envelope** — AHF status and control records (`ok`, `err`, `warn`, `next`, `hint`). These are always AHF.
2. **Result body** — the actual data. This SHOULD be TOON. See https://github.com/toon-format/spec.

### 12.2 Success Envelope

Successful runtime output begins with an `ok` record:

```text
ok <kind> [key=value...]
```

`kind` SHOULD be a plural noun for list results and a singular noun for single-object results. Common metadata keys: `count=`, `more=0|1`, `shown=`, `cursor=`, `truncated=1`.

Example envelope + TOON body:

```text
ok issues count=3 more=0
issues[#3]{key,status,assignee,summary}:
  PROJ-123,open,team-cli,Add AHF support
  PROJ-124,done,team-docs,Document CLI flags
  PROJ-125,open,team-runtime,Write TOON encoder
next inspect "forge issue PROJ-123 --agent-out"
```

TOON format: `[#N]` length marker, two-space indented rows, values quoted only when containing commas or spaces. See https://github.com/toon-format/toon-go for the Go library.

### 12.3 Warning Records

Non-fatal warnings SHOULD use `warn` records immediately after the `ok` line:

```text
warn <code> [key=value...]
```

Warnings do not change success semantics. A successful command with warnings still begins with `ok`.

### 12.4 When TOON Is Insufficient

When TOON cannot represent the result (e.g. a single scalar status, a bare acknowledgement, or a result with no structured data), producers MAY emit key-value lines directly after `ok`:

```text
ok deploy
env prod
status healthy
version 3.2.1
next logs "tool logs --env prod --agent-out"
```

This form SHOULD be used only when there is no tabular or structured data to encode. For any list or object result, TOON is RECOMMENDED.

## 13. Error Handling

### 13.1 Error Envelope

Invalid invocations and runtime errors begin with an `err` record:

```text
err <code> [key=value...]
```

### 13.2 TOON Error Bodies

Where TOON provides a suitable error encoding, producers SHOULD emit a TOON body after `err`. Where TOON is insufficient (e.g. a simple missing-flag correction), producers SHOULD use AHF `hint` and `use` records:

```text
err <code> [key=value...]
hint <direct correction>
use <canonical invocation if useful>
```

Error outputs SHOULD include enough information for an agent to retry without calling human `--help`. Tools SHOULD NOT respond to agent-facing errors only with `run --help`.

For missing required values, errors SHOULD identify the exact missing argument or flag. For enum violations, errors SHOULD list valid values.

Example (simple AHF-only error):

```text
err invalid_enum flag=--type got=note
hint --type enum(decision|fact|pattern|task|observation)
use mem node add <text> --type TYPE [--tag K:V...]
```

## 14. Pagination, Truncation, and Follow-Up Commands

`next` records MAY include a label:

```text
next "tool list --cursor c_123 --agent-out"
next inspect "tool show n_102 --agent-out"
```

If output is paginated, truncated, or partial, the `ok` record SHOULD include metadata such as `more=1`, `truncated=1`, `shown=<n>`, `count=<n>`, or `cursor=<value>`. When `more=1` or `truncated=1`, output SHOULD include an exact continuation command in a `next` record.

## 15. Accuracy and Conformance

### 15.1 Conforming AHF Producer

A conforming AHF producer:

- emits UTF-8 line-oriented AHF records;
- uses registered record prefixes according to this document;
- emits AH1 for `tool --agent-help`;
- emits AH2 for `tool subcmd --agent-help` when subcommands exist;
- emits `ok` or `err` as the first record for `--agent-out`;
- uses TOON for `--agent-out` result bodies where practical;
- avoids human formatting in agent-facing output;
- preserves stable field names where practical.

### 15.2 Conforming Agent-Help CLI

A conforming agent-help CLI:

- supports trailing `--agent-help`;
- includes the human help breadcrumb;
- keeps AH1 and AH2 synchronized with actual command behavior;
- tests every AH1 command path has AH2 output;
- tests every AH2 example or omits examples.

### 15.3 Conforming Agent-Out CLI

A conforming agent-out CLI:

- supports `--agent-out` for structured runtime results;
- emits TOON for result bodies;
- emits AHF `ok`/`err` protocol envelope before the TOON body;
- includes identifiers required for follow-up operations;
- includes exact continuation commands for pagination and truncation;
- redacts sensitive values by default.

## 16. Security and Privacy Considerations

AHF output is plain text and may be consumed by LLM agents with broad access. Producers SHOULD:

- redact secrets, tokens, credentials, and keys from all agent-facing output by default;
- avoid including personal data beyond what is necessary for the agent to act;
- treat `--agent-out` output as potentially logged by agent runtimes.

Producers SHOULD NOT rely on AHF output being private.

## 17. Accessibility and Internationalization

AHF output is ASCII-preferred plain text. It does not use color, emoji, or locale-specific formatting. Non-ASCII field values are UTF-8. Record type tokens are ASCII.

## 18. Versioning

AHF uses semantic versioning.

- **Patch** (0.1.x) — errata and clarifications; no normative changes.
- **Minor** (0.x) — new optional record types or conventions; backward-compatible.
- **Major** (x.0) — breaking changes to record format or required conventions.

Producers MAY advertise their supported AHF version in a `version` key on the `ah1` or `ok` record. Consumers SHOULD tolerate unknown record prefixes gracefully.

## 19. Examples

### 19.1 AH1 Index

```text
ah1 mem :: project memory — store and query facts, decisions, and tasks
cmd node add <text> --type TYPE [--tag K:V...] :: store a new memory node
cmd node list [--type TYPE] [--limit int] [--cursor id] :: list memory nodes
cmd search query <text> [--type TYPE] [--limit int] [--cursor id] :: search nodes by text
cmd search similar <id> [--limit int] [--cursor id] :: find nodes similar to a given node
cmd project status :: show current project settings
cmd project set --key KEY --value VALUE :: set a project setting
more? mem <cmd> --agent-help
```

### 19.2 AH2 Detail

```text
ah2 mem node add
use mem node add <text> --type TYPE [--tag K:V...]
arg text:str req :: node text content
flag --type:enum(decision|fact|pattern|task|observation) req :: node type
flag --tag:kv repeat :: metadata key:value
ex mem node add "postgres 15 required" --type fact --tag project:mem
```

### 19.3 Error Hint

```text
err invalid_enum flag=--type got=note
hint --type enum(decision|fact|pattern|observation)
use mem node add <text> --type TYPE
```

### 19.4 Runtime List (TOON body)

```text
ok nodes count=8 shown=3 more=0
nodes[#3]{id,type,tags,text}:
  n_102,fact,"project:billing|db",postgres 15 required for billing service
  n_088,decision,project:billing,migrate billing read model to postgres
  n_061,pattern,db|ops,use connection pool max 20 in all environments
next inspect "mem search similar n_102 --agent-out"
```

TOON format notes: `[#N]` is the length marker; rows are two-space indented; values are quoted only when they contain commas or spaces.

### 19.5 Runtime Object (TOON key-value table)

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

### 19.6 Scalar Acknowledgement (AHF key-value fallback, §12.4)

When a result has no structured body, AHF key-value lines MAY follow `ok` directly:

```text
ok node
id n_109
type fact
tags spec:ahf
created 2026-05-12
text "toon-go library works"
next list "mem node list --agent-out"
```

### 19.7 Paginated List with Warning

```text
ok nodes count=8 shown=3 more=1 cursor=c_n_104
warn truncated shown=3 total=8
nodes[#3]{id,type,tags,text}:
  n_101,decision,"project:mem|arch",Use Cobra for the reference CLI implementation
  n_102,fact,"project:mem|db",Mock data lives in data.go
  n_103,pattern,db|ops,Use connection pool max 20 in all environments
next "mem node list --limit 3 --cursor c_n_104 --agent-out"
next inspect "mem search query <text> --agent-out"
```

## 20. Implementation Guidance

Implementations SHOULD generate AHF help from existing command metadata rather than maintaining a separate handwritten command list. Common approaches:

- **Cobra (Go)** — hidden persistent flags, derive from `Command.Use`, flags, args, annotations.
- **Click (Python)** — hidden eager options, derive from `Command.params` and context command path.
- **Clap (Rust)** — hidden global args, derive from `Command::get_subcommands()` and arguments.
- **Commander (Node)** — global options and `preAction` hooks, derive from command metadata.
- **Argparse (Python)** — global flags and parser/subparser/action metadata.

For `--agent-out`, use your language's TOON library or generate TOON directly — the format is simple enough to emit without a library for flat result sets.

Implementations SHOULD add tests that verify:

- human `--help` includes the breadcrumb;
- AH1 includes every primary command;
- every AH1 command has AH2 output;
- AH2 required flags and args match validation;
- AH2 examples run successfully;
- `--agent-out` begins with `ok` or `err`;
- TOON body follows the `ok` envelope for list and object results;
- paginated outputs include `next`.

## 21. Open Questions

This draft intentionally leaves several issues open for experimentation:

1. Should `--agent-out` become `--format toon`, `--toon`, or remain as specified?
2. Should a formal ABNF grammar be added for AHF in v0.2?
3. Should MIME type registration be pursued, for example `text/ahf`?
4. Should record prefix extension namespacing be defined?
5. Should conformance levels distinguish `--agent-help`-only from full `--agent-help` plus `--agent-out`?
6. Should examples be packaged as a model comprehension test suite?

## 22. Relationship to Other Formats

### 22.1 AHF and TOON

[TOON (Token-Oriented Object Notation)](https://github.com/toon-format/spec) is the RECOMMENDED encoding for `--agent-out` result bodies. AHF and TOON are complementary layers:

| Layer | Format | Covers |
|---|---|---|
| Discovery | AHF | `--agent-help`, AH1/AH2, command signatures, args, flags, examples |
| Protocol | AHF | `ok`/`err`/`warn` status, `next`/`hint` follow-up |
| Data | TOON | Result bodies, lists, objects, nested structures |

Use TOON when you have structured result data. Use AHF key-value lines only when the result is a bare scalar acknowledgement with no structured body.

### 22.2 Why not JSON for `--agent-out`?

JSON is suboptimal for agents because:

1. **Punctuation overhead.** Braces, brackets, commas, and quotes consume tokens without adding agent-actionable meaning.
2. **No discovery protocol.** JSON encodes data; it says nothing about what commands exist, what arguments they take, or what to call next.
3. **No follow-up commands.** JSON results do not carry `next` commands, pagination cursors, or error hints.

Use `--json` when software is the consumer. Use `--agent-out` when an agent is.

### 22.3 Why not plain text / `--help`?

Human `--help` is optimized for human reading: prose descriptions, aligned columns, color codes, aliases, version strings. Agents spend tokens on decoration and must infer argument types, required flags, and valid values from prose. `--agent-help` replaces this with a compact, machine-structured surface derived from the same command metadata.

### 22.4 Why not man pages or completion scripts?

Man pages and shell completion scripts encode some overlapping information but are not designed for runtime agent consumption:

- Separate files, not embedded in the binary's output stream.
- Platform-specific parsing required.
- No follow-up commands, pagination, or error hints.
- Completion scripts describe valid tokens, not valid typed invocations.

### 22.5 How does this relate to MCP or plugin protocols?

[Model Context Protocol (MCP)](https://modelcontextprotocol.io) and similar protocols wrap capabilities in an RPC server layer. This is powerful but heavyweight:

- Separate server process required.
- Server SDK implementation required.
- Agent runtime must support MCP specifically.
- Deployment and maintenance surface added.

agent-help is complementary: it makes the CLI itself agent-readable. MCP wrappers, plugins, and skills can still be built on top of an agent-help CLI, but they are not required for basic agent-native discovery and invocation.

## Appendix A. Minimal Agent Implementation Prompt

When asking an agent to implement agent-help in a CLI, provide `llms-full.txt` or the following minimal instruction:

```text
Implement the agent-help convention. Add trailing global --agent-help for AH1/AH2 invocation help using AHF line records. Add --agent-out for structured runtime results: emit an ok/err AHF status line, then a TOON body for result data. Humans keep --help. Software keeps --json. Agents get --agent-help and --agent-out. Follow AHF-RFC.md.
```
