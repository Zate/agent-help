# agent-help Repository — Agent Guidelines

These instructions apply to any agent working in this repository.

## What this repo is

This is the source for the **agent-help** open convention — a standard for building CLIs that AI agents can discover and use directly, while still allowing MCP servers, skills, and plugins to wrap the same CLI.

The convention has two layers:
- **AHF** (Agent Help Format) — the spec for `--agent-help` output and the `--agent-out` protocol envelope. Defined in `AHF-RFC.md`.
- **TOON** — the recommended encoding for `--agent-out` result bodies. External spec at https://github.com/toon-format/spec.

## Critical sync rule

**`llms.txt` and `llms-full.txt` MUST stay in sync with `AHF-RFC.md` at all times.**

These files are not supplementary docs — they are the primary interface through which agents learn to implement agent-help. If they drift from the spec, implementations will be wrong.

Specifically:

| When you change | Also update |
|---|---|
| Record prefix registry (§9) | `llms-full.txt` AHF record prefixes list, `references/REFERENCE.md` prefix table |
| Scalar type registry (§10) | `llms-full.txt` types list, `references/REFERENCE.md` scalar types table |
| AH1 format or rules (§11.1) | `llms-full.txt` AH1 section, `.agents/skills/ahf/SKILL.md` AH1 section, examples/01-* |
| AH2 format or rules (§11.2) | `llms-full.txt` AH2 section, `.agents/skills/ahf/SKILL.md` AH2 section, examples/02-* |
| Error format (§13) | `llms-full.txt` error hint section, `.agents/skills/ahf/SKILL.md` AE1 section, examples/03-* |
| `--agent-out` envelope (§12) | `llms-full.txt` --agent-out section, `.agents/skills/ahf/SKILL.md`, examples/04-* 05-* 06-* |
| TOON delegation rules (§12) | `llms-full.txt` TOON section, `references/REFERENCE.md` |
| Pagination/truncation (§14) | `llms-full.txt` follow-up commands section, examples/06-* |
| Conformance requirements (§15) | `llms-full.txt` implementation requirements section |
| Any normative text | `FAQ.md` if it answers a common question |
| Any of the above | `llms.txt` if the core concept or shape has changed |

## File purposes

| File | Purpose | Audience |
|---|---|---|
| `AHF-RFC.md` | Full normative AHF spec | Spec readers, implementers |
| `llms.txt` | Short orientation (~50 lines) | Agents discovering what agent-help is |
| `llms-full.txt` | Full implementation brief | Agents implementing agent-help in a CLI |
| `.agents/skills/ahf/SKILL.md` | agentskills.io-compatible build skill | Agents using skill-aware platforms |
| `FAQ.md` | Conversational "why agent-help?" | Skeptical humans, first-time visitors |
| `references/REFERENCE.md` | Quick-reference card | Agents needing a registry lookup |
| `examples/` | Annotated output samples | Agents and implementers |
| `docs/CONFORMANCE.md` | Conformance levels | Implementers |
| `spec/ahf-v0.1.json` | Machine-readable spec manifest | Tooling, validators |
| `docs/PARSING.md` | Practical parsing notes | Implementers |
| `docs/PRIOR_ART.md` | Related work and positioning | Humans, reviewers |
| `docs/RELEASE_CHECKLIST.md` | Release checklist | Maintainers |
| `SECURITY.md` | Security reporting and redaction guidance | Humans, implementers |
| `site/` | Static landing page source | Web visitors |
| `Makefile` | Local verification entrypoint | Contributors |
| `docs/VERSIONING.md` | Draft stability and release policy | Implementers, contributors |
| `LICENSE-DOCS` | Documentation/spec license | Humans, legal review |
| `NOTICE` | Project notices and license split | Humans, legal review |

## `llms.txt` vs `llms-full.txt`

- **`llms.txt`** — orientation only. Tells an agent what agent-help is and where to find more. ~50 lines. Does not contain enough to implement agent-help.
- **`llms-full.txt`** — self-contained implementation brief. Contains everything needed to implement `--agent-help` and `--agent-out` in any CLI without reading any other file. Must be kept complete and accurate.

When in doubt: if a change affects how a CLI author would implement agent-help, it belongs in `llms-full.txt`. If it only affects spec readers, it can stay in `AHF-RFC.md` alone.

## Key spec decisions (do not reverse without a spec change)

- `--agent-help` is **always trailing**: `tool subcmd --agent-help` ✓ — `tool --agent-help subcmd` ✗
- `more?` is the AH1 pointer record (not a shell command — the `?` is intentional)
- `more=0|1` on `ok` headers is a key=value boolean (unambiguous, leave as-is)
- `next` is the follow-up command record for results and errors
- `--agent-out` result bodies use **TOON** — not AHF's own row format
- `ok`/`err`/`warn`/`next`/`hint` are **AHF** protocol records — always present regardless of TOON
- `_` means null/unknown/not applicable
- `|` separates short lists inside a single field value
- Agents get `--agent-help` and `--agent-out`. Humans get `--help`. Software gets `--json`. These are separate surfaces.

## Canonical URLs

All public-facing files should use `https://zate.github.io/agent-help/` as the canonical base URL.
Repository: `https://github.com/Zate/agent-help`.
Do not hardcode other URLs without updating both.

## Before finishing any change

- [ ] `AHF-RFC.md` section numbers in cross-references still correct?
- [ ] `llms-full.txt` reflects the change?
- [ ] `llms.txt` still accurately describes the shape of agent-help?
- [ ] `references/REFERENCE.md` tables still match the registries?
- [ ] Examples in `examples/` still valid per the current spec?
- [ ] `FAQ.md` still accurate?
- [ ] `.agents/skills/ahf/SKILL.md` still consistent with the spec?
