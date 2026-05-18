# agent-help Repository — Agent Guidelines

These instructions apply to any agent working in this repository.

## What this repo is

This is the source for the **agent-help** open convention — a standard for building CLIs that AI agents can discover and use directly, while still allowing MCP servers, skills, and plugins to wrap the same CLI.

The convention has two layers:
- **AHF** (Agent Help Format) — the spec for `--agent-help` output and the `--agent-out` protocol envelope. Defined in `AHF-RFC.md`.
- **TOON** — the recommended encoding for `--agent-out` result bodies. External spec at https://github.com/toon-format/spec.

---

## Canonical source hierarchy

`AHF-RFC.md` is the single normative source of truth. Everything else derives from it:

```
AHF-RFC.md                          ← canonical spec (normative)
  ├── llms-full.txt                  ← self-contained implementation brief (derives from RFC)
  ├── llms.txt                       ← short orientation (derives from RFC, points to llms-full.txt)
  ├── .agents/skills/ahf/SKILL.md   ← build-time skill (derives from RFC + llms-full.txt)
  │     └── references/REFERENCE.md ← quick-lookup card (derives from RFC registries)
  ├── agent-help.ahf                 ← AH1 dogfood (derives from RFC AH1 format)
  ├── agent-help-full.ahf            ← AH2 dogfood (derives from RFC AH2 format)
  ├── examples/                      ← annotated wire samples (must match RFC + impl)
  └── site/                          ← human-readable HTML (derives from RFC + other .md files)
        ├── spec.html                ← rendered AHF-RFC.md
        ├── faq.html                 ← rendered FAQ.md
        ├── examples/index.html      ← rendered examples/
        ├── prompts.html             ← rendered prompts.md
        └── docs/implementation-guides.html ← rendered docs/IMPLEMENTATION_GUIDES.md
```

When you change `AHF-RFC.md`, propagate changes **down** through this hierarchy. Never propagate upward — if a derived file disagrees with the RFC, the RFC wins.

---

## Critical sync rules

| When you change | Also update |
|---|---|
| Record prefix registry (§9) | `llms-full.txt` prefixes list · `.agents/skills/ahf/references/REFERENCE.md` prefix table · `site/spec.html#s9` |
| Scalar type registry (§10) | `llms-full.txt` types list · `.agents/skills/ahf/references/REFERENCE.md` scalar types table · `site/spec.html#s10` |
| AH1 format or rules (§11.1) | `llms-full.txt` AH1 section · `.agents/skills/ahf/SKILL.md` AH1 section · `agent-help.ahf` · `examples/01-*` · `site/examples/index.html#ex01` |
| AH2 format or rules (§11.2) | `llms-full.txt` AH2 section · `.agents/skills/ahf/SKILL.md` AH2 section · `agent-help-full.ahf` · `examples/02-*` · `site/examples/index.html#ex02` |
| Error format (§13) | `llms-full.txt` error section · `.agents/skills/ahf/SKILL.md` AE1 section · `examples/03-*` · `site/examples/index.html#ex03` |
| `--agent-out` envelope (§12) | `llms-full.txt` agent-out section · `.agents/skills/ahf/SKILL.md` · `examples/04-* 05-* 06-*` · `site/examples/index.html` |
| TOON delegation rules (§12) | `llms-full.txt` TOON section · `.agents/skills/ahf/references/REFERENCE.md` |
| Pagination/truncation (§14) | `llms-full.txt` follow-up section · `examples/06-*` · `site/examples/index.html#ex06` |
| Conformance requirements (§15) | `llms-full.txt` implementation requirements section · `docs/CONFORMANCE.md` · `site/spec.html#s15` |
| Any normative text | `FAQ.md` if it answers a common question · `site/faq.html` to match |
| Core concept or shape | `llms.txt` if the orientation summary has changed |

---

## File purposes

### Normative spec

| File | Purpose | Audience |
|---|---|---|
| `AHF-RFC.md` | Full normative AHF spec — canonical source of truth | Spec readers, implementers |
| `spec/ahf-v0.1.json` | Machine-readable spec manifest and registry | Tooling, validators |

### Agent-facing files (for agents consuming or implementing agent-help)

| File | Purpose | Audience | Size |
|---|---|---|---|
| `llms.txt` | Orientation only — what is this, where to go next | Agent discovering the project | ~55 lines |
| `llms-full.txt` | Self-contained implementation brief — everything needed to implement agent-help in any CLI, no other file required | Agent implementing agent-help | ~311 lines |
| `.agents/skills/ahf/SKILL.md` | agentskills.io-compatible build-time skill — loaded by skill-aware platforms (Copilot, Cursor, etc.) to assist implementing agent-help | Agents on skill-aware platforms | ~287 lines |
| `.agents/skills/ahf/references/REFERENCE.md` | Quick-lookup card — record prefix table, scalar types, key rules, format shapes. Loaded on demand mid-task | Agents needing a registry lookup | ~87 lines |
| `agent-help.ahf` | AH1 dogfood — this project's own `--agent-help` root index in AHF wire format | Agent runtimes that fetch .ahf files | 7 lines |
| `agent-help-full.ahf` | AH2 dogfood — full AH2 detail records for each "command" in AHF wire format | Agent runtimes that fetch .ahf files | 23 lines |
| `prompts.md` | Copy-ready prompts for humans to paste into an agent | Humans getting started with an agent |  |

### Why we have both `llms-full.txt` and `SKILL.md`

They overlap in content but serve different use cases:

- **`llms-full.txt`** — universal. Works with any agent, any platform. Paste it, link to it, or `curl` it. No platform support needed. Format is plain prose + code blocks.
- **`SKILL.md`** — platform-aware. Follows the agentskills.io YAML frontmatter format. Loaded automatically by skill-aware platforms (`npx skills add Zate/agent-help`). Includes frontmatter metadata that platforms use to index and activate the skill.

If you are implementing agent-help in a CLI and are not on a skill-aware platform, use `llms-full.txt`. If you are on a skill-aware platform, the skill will be loaded automatically.

### Why we have both `agent-help.ahf` and `agent-help-full.ahf`

They are an intentional AH1/AH2 pair, not duplicates:

- **`agent-help.ahf`** — AH1 index (root-level `tool --agent-help` output). Lists all "commands" with one-line summaries.
- **`agent-help-full.ahf`** — AH2 detail (per-command `tool subcmd --agent-help` output). Full arg/flag/example records for each command.

Both must stay in sync with the AH1/AH2 format rules in `AHF-RFC.md §11`.

### Human-facing docs

| File | Purpose | Audience |
|---|---|---|
| `FAQ.md` | Conversational "why agent-help?" | Skeptical humans, first-time visitors |
| `README.md` | Project overview and quick-start | GitHub visitors |
| `docs/CONFORMANCE.md` | Conformance level definitions (L1/L2/L3) | Implementers |
| `docs/VERSIONING.md` | Draft stability and release policy | Implementers, contributors |
| `docs/IMPLEMENTATION_GUIDES.md` | Framework guides (Cobra, Click, Clap, Commander, argparse) | Implementers |
| `docs/PARSING.md` | Practical parsing notes | Implementers |
| `docs/PRIOR_ART.md` | Related work and positioning | Humans, reviewers |
| `docs/ADOPTER_CHECKLIST.md` | Readiness checklist for adopters | Implementers |
| `docs/RELEASE_CHECKLIST.md` | Release checklist | Maintainers |
| `SECURITY.md` | Security reporting and redaction guidance | Humans, implementers |
| `CODE_OF_CONDUCT.md` | Community standards | Contributors |
| `CONTRIBUTING.md` | How to contribute | Contributors |
| `CHANGELOG.md` | Release history | All |

### Site (human-readable HTML, served at zate.github.io/agent-help/)

| File | Derives from | Notes |
|---|---|---|
| `site/index.html` | Overall project | Main landing page |
| `site/spec.html` | `AHF-RFC.md` | Must stay in sync with RFC content and section numbers |
| `site/faq.html` | `FAQ.md` | Must stay in sync with FAQ.md |
| `site/examples/index.html` | `examples/` | Must stay in sync with example files |
| `site/prompts.html` | `prompts.md` | Must stay in sync with prompts.md |
| `site/docs/implementation-guides.html` | `docs/IMPLEMENTATION_GUIDES.md` | Must stay in sync |
| `site/styles.css` | — | Shared stylesheet for all site pages |
| `site/copy.js` | — | Shared clipboard copy button script |

**Site sync rule:** When content changes in any source `.md` file that has a corresponding `site/` page, update both. The site pages are intentional HTML renderings of the markdown — not auto-generated, so drift is possible. Check both when editing.

### Infrastructure

| File | Purpose |
|---|---|
| `examples/` | Raw annotated AHF wire format examples (01–06) |
| `impl/` | Go/Cobra reference implementation |
| `tests/` | AHF fixture tests |
| `scripts/` | Validation scripts |
| `Makefile` | Local verification entrypoint |
| `spec/ahf-v0.1.json` | Machine-readable registry |
| `LICENSE` | Apache-2.0 (code) |
| `LICENSE-DOCS` | CC-BY-4.0 (spec/docs) |
| `NOTICE` | Project notices and license split |

---

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

---

## Canonical URLs

All public-facing files should use `https://zate.github.io/agent-help/` as the canonical base URL.
Repository: `https://github.com/Zate/agent-help`.
Do not hardcode other URLs without updating both.

---

## Before finishing any change

### Spec changes
- [ ] `AHF-RFC.md` section numbers in cross-references still correct?
- [ ] `llms-full.txt` reflects the change?
- [ ] `llms.txt` still accurately describes the shape of agent-help?
- [ ] `.agents/skills/ahf/SKILL.md` still consistent with the spec?
- [ ] `.agents/skills/ahf/references/REFERENCE.md` tables still match the registries?
- [ ] `agent-help.ahf` still valid AH1 per current AH1 rules?
- [ ] `agent-help-full.ahf` still valid AH2 per current AH2 rules?
- [ ] Examples in `examples/` still valid per the current spec?
- [ ] `spec/ahf-v0.1.json` still reflects the current registry?

### Content changes
- [ ] `FAQ.md` updated if normative text changed?
- [ ] `site/faq.html` in sync with `FAQ.md`?
- [ ] `site/spec.html` in sync with `AHF-RFC.md`?
- [ ] `site/examples/index.html` in sync with `examples/`?
- [ ] `site/prompts.html` in sync with `prompts.md`?
- [ ] `site/docs/implementation-guides.html` in sync with `docs/IMPLEMENTATION_GUIDES.md`?

### File moves or renames
- [ ] All internal cross-references updated (grep for old filename)?
- [ ] All site HTML links updated?
- [ ] `README.md` file table updated?
- [ ] `AGENTS.md` file purpose table updated?
- [ ] GitHub URLs in site pages updated (they use `blob/main/<path>`)?
