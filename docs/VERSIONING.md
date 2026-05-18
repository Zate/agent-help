# Versioning and Releases

agent-help has two related version streams:

- **Specification version**: the AHF draft version, currently `v0.1`.
- **Repository releases**: GitHub releases/tags for snapshots of this repository, examples, and reference implementation.

## Specification Stability

The core wire shape is intended to be stable enough for trial implementations:

- `--agent-help` emits AHF line records.
- Root help uses AH1 (`ah1`, `cmd`, `more?`).
- Command help uses AH2 (`ah2`, `use`, `arg`, `flag`, `ex`).
- `--agent-out` starts with an AHF `ok` or `err` envelope.
- Structured `--agent-out` bodies use TOON.
- Follow-up actions use AHF `next` records.

The following may still change before `v1.0`:

- Additional record prefixes.
- Additional scalar types.
- Conformance levels.
- Optional aliases for `--agent-out`.
- Formal grammar or media-type registration.

## Release Policy

Before a repository release:

1. `make test` must pass.
2. `AHF-RFC.md`, `llms.txt`, `llms-full.txt`, `.agents/skills/ahf/SKILL.md`, `references/REFERENCE.md`, and `examples/` must be checked for drift.
3. `spec/ahf-v0.1.json` and `docs/CONFORMANCE.md` must reflect registry and conformance changes.
4. The Go reference implementation must build and pass golden-output tests.
5. Documented examples must match current CLI output.
6. AHF examples must pass the lightweight validator.
7. Markdown links must resolve locally.
8. Run `make update-examples` if the reference CLI output changed.
9. The release notes must state the AHF draft version covered by the release.

Use semantic-style repository tags while the spec is pre-1.0:

- `v0.1.x`: clarifications, examples, test coverage, non-normative fixes.
- `v0.2.0`: backward-compatible spec additions.
- `v1.0.0`: first stable AHF release.
