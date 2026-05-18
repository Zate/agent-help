# agent-help — Frequently Asked Questions

## Why does agent-help exist? Don't we have enough formats?

The short answer: existing formats solve adjacent problems, not this one.

CLIs already have two good output modes:
- **`--help`** — prose for humans
- **`--json`** — structured data for software

Neither works well for AI agents. Human output wastes tokens on decoration. JSON wastes tokens on punctuation and says nothing about how to discover commands, what to call next, or how to recover from errors.

agent-help defines a *third mode*: `--agent-help` for token-efficient discovery (using AHF line records) and `--agent-out` for dense runtime results (using TOON for data, AHF for the protocol envelope). The CLI exposes this surface directly; wrappers and servers can still build on top.

---

## How does agent-help relate to TOON?

They're complementary layers — designed to work together.

[TOON (Token-Oriented Object Notation)](https://github.com/toon-format/spec) is the recommended encoding for `--agent-out` result bodies. AHF and TOON each do their job:

| Layer | Format | Covers |
|---|---|---|
| Discovery | AHF | `--agent-help`, AH1/AH2, command signatures, args, flags, examples |
| Protocol | AHF | `ok`/`err`/`warn` status, `next`/`hint` follow-up records |
| Data | TOON | Result bodies, lists, objects, nested structures |

The golden path for `--agent-out`:

```text
ok nodes count=3 more=0                      ← AHF status line
nodes[#3]{id,type,tags,text}:                ← TOON result body
  n_102,fact,project:billing,postgres 15 required
  n_088,decision,project:billing,migrate read model
  n_061,pattern,db|ops,use connection pool max 20
next inspect "mem search similar n_102 --agent-out"    ← AHF follow-up
```

TOON format: `[#N]` length marker, two-space indented rows, values quoted only when containing commas or spaces.

TOON handles the data. AHF handles the protocol.

---

## Why not just output JSON with `--json`?

You should keep `--json` for software consumers — agent-help doesn't replace it. But JSON is suboptimal for agents:

1. **Token overhead.** Braces, brackets, commas, and quotes repeat constantly. A 20-field list in JSON uses 3–5× more tokens than the equivalent TOON.
2. **No discovery.** JSON encodes data. It says nothing about what commands exist, what args they take, or what to call next.
3. **No follow-up.** A JSON result doesn't tell an agent what command to run for the next page, to inspect a failure, or to recover from an error. AHF `next` records embed this directly.

---

## Why not just tell agents to read `--help`?

Human `--help` is designed for humans: prose, aligned columns, color escapes, examples chosen for clarity, aliases, version strings. Agents spend tokens on all of that decoration and still have to infer argument types, required flags, and valid enum values from prose.

`--agent-help` emits the same information as compact AHF records — derived from the same command metadata, without the presentation layer.

---

## Why not use man pages or shell completion scripts?

They encode some overlapping information, but:

- They're separate files, not embedded in the tool's stdout.
- They require platform-specific parsing.
- They don't encode follow-up commands, pagination cursors, or error hints.
- Completion scripts describe valid tokens, not valid typed invocations.

agent-help is in the tool's output stream itself — no side-channel files, no platform-specific parsing.

---

## How does agent-help relate to MCP or plugin protocols?

[MCP](https://modelcontextprotocol.io) and similar protocols wrap CLI capabilities in an RPC server. Powerful, but heavyweight:

- Separate server process to run and maintain.
- Server SDK to implement.
- Agent runtime must support MCP specifically.
- Deployment and configuration surface.

agent-help is complementary. It makes the CLI itself agent-readable, so an agent can use it directly. MCP wrappers, plugins, and skills can still be built on top of an agent-help CLI, but they are not required for basic discovery and invocation.

---

## What does "conforming" mean?

A typical conforming CLI:

1. Adds a breadcrumb or equivalent pointer to normal `--help`.
2. Implements trailing `--agent-help` that returns AH1 (index) or AH2 (command detail) in AHF format.
3. Optionally implements `--agent-out` on result-returning commands, emitting an AHF `ok`/`err` envelope followed by a TOON result body.

That's it. You can adopt discovery first and add `--agent-out` later. See [§15 of the spec](AHF-RFC.md#15-accuracy-and-conformance) for conformance notes.

---

## Does agent-help require a specific language or framework?

No. AHF is a text convention, not a library. TOON is simple enough to emit without a library for most result sets. Both work with any language and any CLI framework. The [spec](AHF-RFC.md#20-implementation-guidance) includes notes for Cobra (Go), Click (Python), Clap (Rust), Commander (Node), and argparse.

---

## What's the relationship to agentskills.io?

[agentskills.io](https://agentskills.io) defines a standard format for packaging reusable agent capabilities as skills (a `SKILL.md` file with YAML frontmatter and instructions). This repo ships a `SKILL.md` that follows that format — it's a **build-time development tool** that helps you implement agent-help in your own CLI using an agent assistant.

agent-help itself has no dependency on agentskills.io. Once your CLI conforms, agents can use it directly with no skill file, MCP server, or any other runtime component.

---

## Is v0.1 stable enough to implement against?

Yes, for trial implementations. The core wire shape, discovery convention, and `--agent-out` envelope are intended to be stable. Registry entries, naming details, and conformance levels may still change before v1.0. The [open questions](AHF-RFC.md#21-open-questions) are unlikely to affect a basic implementation.

If you implement agent-help, please open an implementation report issue — it helps track adoption and find spec gaps before 1.0.

---

## How do I contribute or raise a question?

See [CONTRIBUTING.md](CONTRIBUTING.md). The short version: open an issue using one of the templates (open question, errata, feature proposal, implementation report). For small fixes, a PR is welcome directly.
