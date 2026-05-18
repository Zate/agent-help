# agent-help

**An open convention for building CLIs that AI agents can discover and use directly.**

[![License: Apache-2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Docs License: CC-BY-4.0](https://img.shields.io/badge/docs-CC--BY--4.0-lightgrey.svg)](LICENSE-DOCS)
[![Spec: Draft v0.1](https://img.shields.io/badge/Spec-Draft_v0.1-green.svg)](AHF-RFC.md)
[![CI](https://github.com/Zate/agent-help-dev/actions/workflows/ci.yml/badge.svg)](../../actions/workflows/ci.yml)
[![Docs](https://img.shields.io/badge/docs-site%2F-0ea5e9.svg)](site/)

---

> **Draft status:** agent-help/AHF is currently Draft v0.1. The wire shape is ready for trial implementations, but naming details and conformance levels may change before v1.0. Treat early integrations as experimental. See [`VERSIONING.md`](docs/docs/VERSIONING.md) and [`CHANGELOG.md`](CHANGELOG.md).

---

## The problem

CLIs have two established output modes: human-readable prose (`--help`, colorized tables) and machine-readable serialization (`--json`). Neither is right for AI agents. Human output wastes tokens on decoration. JSON wastes tokens on punctuation and structure an agent doesn't need. Agents need a third mode: dense, line-oriented, action-ready text exposed by the CLI itself.

## The solution

The agent-help convention defines two CLI surfaces:

| Convention | Flag | Format | Purpose |
|---|---|---|---|
| **Agent Help** | `--agent-help` | AHF | Token-efficient invocation help — what to call, with what args |
| **Agent Output** | `--agent-out` | TOON + AHF envelope | Dense structured runtime results |

```
human:    tool --help
software: tool --json
agent:    tool --agent-help
agent:    tool search query --agent-out
```

The goal: any agent can discover and use any conforming CLI directly. Wrappers, MCP servers, plugins, and skills can still build on top of that surface, but they are not required for basic use.

**AHF** (Agent Help Format) is the line-record spec for `--agent-help` output and the `ok`/`err` protocol envelope for `--agent-out`. **TOON** is the recommended encoding for `--agent-out` result bodies. They are separate, complementary layers.

---

## Quick start

For a human skim, start with the short orientation:

```bash
cat llms.txt
```

For an agent, point it directly at the raw files:

- `https://raw.githubusercontent.com/Zate/agent-help/main/llms.txt`
- `https://raw.githubusercontent.com/Zate/agent-help/main/llms-full.txt`
- `https://raw.githubusercontent.com/Zate/agent-help/main/.agents/skills/ahf/SKILL.md`

Useful prompts:

```text
Read https://raw.githubusercontent.com/Zate/agent-help/main/llms-full.txt and explain what agent-help is for.
```

```text
Read https://raw.githubusercontent.com/Zate/agent-help/main/llms-full.txt and implement agent-help in this CLI.
```

If your agent supports skills, point it at:

```text
https://raw.githubusercontent.com/Zate/agent-help/main/.agents/skills/ahf/SKILL.md
```

These files are intentionally optimized for agents. They may feel terse to humans; that is the point.

Try the reference implementation:

```bash
make demo
```

Verify the repository before proposing changes:

```bash
make test
```

---

## Shape at a glance

```text
tool --agent-help
  └─ AHF help records
     ah1 / cmd / more?
     ah2 / use / arg / flag / ex

tool subcmd args --agent-out
  ├─ AHF protocol envelope
  │  ok / err / warn
  ├─ TOON result body
  │  records[#N]{field,...}:
  └─ AHF follow-up records
     next / hint
```

`more?` is a pointer record, not a command. The question mark is intentional: it tells agents "more detail is available with this shape" without looking like a shell token to copy.

---

## Quick example

```text
$ tool search postgres --agent-out
ok nodes count=2 more=0
nodes[#2]{id,type,tags,text}:
  n_102,fact,"project:billing|db",postgres 15 required
  n_088,decision,project:billing,migrate read model
next inspect "tool show n_102 --agent-out"
```

```text
$ tool --agent-help
ah1 tool :: example data CLI
cmd search <query> [--type TYPE] [--limit int] :: search records
cmd show <id> :: display one record
cmd rm <id> :: delete a record
more? tool <cmd> --agent-help
```

---

## Repository contents

| File | Purpose |
|---|---|
| [`AHF-RFC.md`](AHF-RFC.md) | Full AHF draft specification (Internet-Draft style) |
| [`FAQ.md`](FAQ.md) | Why agent-help? How does AHF relate to TOON, JSON, MCP? |
| [`.agents/skills/ahf/SKILL.md`](.agents/skills/ahf/SKILL.md) | agentskills.io-compatible skill for implementing agent-help |
| [`llms.txt`](llms.txt) | Short agent orientation file (point agents here first) |
| [`llms-full.txt`](llms-full.txt) | Full implementation brief (paste to ask an agent to implement agent-help) |
| [`agent-help.ahf`](agent-help.ahf) | AHF-style docs entry point |
| [`agent-help-full.ahf`](agent-help-full.ahf) | AH2-style details for agent-help doc flows |
| [`prompts.md`](prompts.md) | Copy/paste prompts for agents |
| [`references/REFERENCE.md`](references/REFERENCE.md) | Quick-reference card |
| [`examples/`](examples/) | Standalone examples for testing model comprehension |
| [`CONFORMANCE.md`](docs/docs/CONFORMANCE.md) | Conformance levels for implementers |
| [`spec/ahf-v0.1.json`](spec/ahf-v0.1.json) | Machine-readable registry and shape manifest |
| [`docs/`](docs/) | Parsing notes, adopter checklist, framework notes, prior art, and release checklist |
| [`site/`](site/) | Static landing page source |
| [`impl/`](impl/) | Go/Cobra reference implementation with golden tests |
| [`tests/`](tests/) | Validator fixtures for accepted and rejected AHF snippets |
| [`VERSIONING.md`](docs/docs/VERSIONING.md) | Draft stability and release policy |
| [`CHANGELOG.md`](CHANGELOG.md) | Project change history |

---

## Implementing agent-help in your CLI

### 1. Add the breadcrumb to `--help`

```text
LLM agent? Use --agent-help for token-optimized usage.
```

Equivalent wording is fine; the point is that agents can discover the agent-facing surface from normal help.

### 2. Add `--agent-help` (trailing global flag)

```text
tool --agent-help          → AH1 index of all commands
tool subcmd --agent-help   → AH2 detail for that command
```

### 3. Add `--agent-out` to result-returning commands

```text
tool subcmd args --agent-out   → agent-help runtime output
```

### Framework guides

- **Go / Cobra** — hidden persistent flag, generate from `Command.Use` + annotations
- **Python / Click** — hidden eager global option, derive from `Command.params`
- **Rust / Clap** — hidden global arg, intercept after parse
- **Node / Commander** — global option, `preAction` hook
- **Python / argparse** — global flag, generate from parser/subparser metadata

See [`AHF-RFC.md §20`](AHF-RFC.md#20-implementation-guidance) for general details and the framework guides in [`docs/`](docs/).

### Verify this repository

```bash
make test
make demo
make verify-examples
make verify-doc-examples
make validate-ahf
make update-examples
make release-check
```

Use `make release-check` before tagging or publishing. It runs tests, validates examples and fixtures, checks registry drift, checks docs/site links, runs demo commands, and verifies the spec table of contents.

---

## Implementing via an agent

Point your agent at `llms-full.txt` and ask it to add agent-help support:

```
Read llms-full.txt and implement agent-help in this CLI.
```

Or use the `.agents/skills/ahf/SKILL.md` directly — it follows the [agentskills.io](https://agentskills.io) format, so any agent platform that supports skills can load it to help you build agent-help support into your CLI.

---

## Specification

The AHF spec lives in [`AHF-RFC.md`](AHF-RFC.md). It is written in Internet-Draft style with:

- RFC 2119-style guidance, with hard requirements used sparingly
- A complete AHF record prefix registry
- A scalar type registry
- Conformance requirements for producers
- `--agent-out` TOON delegation convention
- Security and privacy considerations
- Versioning policy

**Current status:** Draft v0.1 — the core wire shape is stable enough for trial implementations, while naming details, registry additions, and conformance levels may change before v1.0. See [§21 Open Questions](AHF-RFC.md#21-open-questions) and the [issue tracker](../../issues) for active discussions.

---

## The `.agents/skills/ahf/SKILL.md` in this repo

This repo ships a `.agents/skills/ahf/SKILL.md` that follows the [agentskills.io](https://agentskills.io) spec format. It is a **build-time tool** — a skill for agent-assisted development that helps you implement agent-help (`--agent-help` + `--agent-out`) in your own CLI. It is not required at runtime. Once your CLI conforms, agents can use it directly with no skill file needed.

---

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md). Open questions and proposals live in the [issue tracker](../../issues).

---

## License

Code and examples: [Apache 2.0](LICENSE)  
Documentation and specification: [CC-BY-4.0](LICENSE-DOCS)  
Notices: [NOTICE](NOTICE)
