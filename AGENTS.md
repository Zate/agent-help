# agent-help Repository — Agent Guidelines

This repo is for **agent-help**: MCP-lite for CLIs.

The idea is intentionally small:

- `--agent-help` gives AI agents compact command discovery.
- `--agent-out` gives AI agents dense structured runtime results.
- Humans keep `--help`.
- Scripts keep `--json`.
- Agents using an enabled CLI should not need this repo, the spec, `llms-full.txt`, or a skill at runtime.

## Minimum viable docs rule

Do not add standalone docs unless they have a clear reader and job-to-be-done.

Prefer these homes:

- Human overview -> `README.md`
- Why / tradeoffs -> `FAQ.md`
- Implementation wiring -> `docs/IMPLEMENTATION_GUIDES.md`
- Exact wire-format contract -> `SPEC.md`
- Copy/paste implementation prompts -> `prompts.md`
- Small output samples -> `examples/`
- Rendered site versions -> `site/`

If a proposed doc is only a checklist, conformance label, release ritual, prior-art essay, or maintainer note, first try to fold the useful part into one of the files above.

## Canonical source hierarchy

`SPEC.md` is the detailed implementer contract. Keep derived files in sync when the wire shape changes:

```text
SPEC.md
  ├── llms-full.txt                  implementation brief for coding agents
  ├── llms.txt                       short orientation
  ├── .agents/skills/ahf/SKILL.md    optional build-time skill
  │     └── references/REFERENCE.md  quick registry lookup
  ├── agent-help.ahf                 dogfood AH1 example
  ├── agent-help-full.ahf            dogfood AH2 examples
  ├── examples/                      raw wire samples
  └── site/                          human-readable HTML
```

When `SPEC.md` changes, propagate down. Do not propagate a derived file upward unless the user explicitly decides to change the contract.

## Sync rules

| When you change | Also update |
|---|---|
| Record prefixes | `llms-full.txt` · `.agents/skills/ahf/references/REFERENCE.md` · `spec/ahf-v0.1.json` · `site/spec.html` |
| Scalar types | `llms-full.txt` · `.agents/skills/ahf/references/REFERENCE.md` · `spec/ahf-v0.1.json` · `site/spec.html` |
| `--agent-help` AH1/AH2 shape | `llms-full.txt` · `.agents/skills/ahf/SKILL.md` · `agent-help*.ahf` · `examples/01-*` / `examples/02-*` · `site/examples/index.html` |
| Error format | `llms-full.txt` · `.agents/skills/ahf/SKILL.md` · `examples/03-*` · `site/examples/index.html` |
| `--agent-out` envelope | `llms-full.txt` · `.agents/skills/ahf/SKILL.md` · `examples/04-*` to `06-*` · `site/examples/index.html` |
| Implementation guidance | `docs/IMPLEMENTATION_GUIDES.md` · `site/docs/implementation-guides.html` |
| Human explanation | `README.md` / `FAQ.md` · matching `site/` page |
| Prompts | `prompts.md` · `site/prompts.html` |

## Key decisions

- `--agent-help` is documented in trailing form: `tool subcmd --agent-help`.
- `more?` is the AH1 pointer record.
- `more=0|1` on `ok` headers is key=value metadata.
- `next` is the follow-up command record for results and errors.
- `--agent-out` result bodies use TOON.
- `ok`, `err`, `warn`, `next`, and `hint` are AHF protocol records.
- `_` means null, unknown, or not applicable.
- `|` separates short lists inside a field value.

## Public URLs

Canonical site: `https://zate.github.io/agent-help/`  
Repository: `https://github.com/Zate/agent-help`

## Before finishing changes

- Run focused checks for the files you touched; prefer `make test` before final handoff.
- Run link lint after deleting or moving docs.
- Keep site HTML pages in sync with their source markdown.
- Leave unrelated untracked files alone.
