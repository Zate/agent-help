# Changelog

All notable project changes should be recorded here.

This project is pre-1.0. The AHF wire shape is intended to remain coherent within a draft series, but normative changes can still happen before v1.0.

## Unreleased

- Added repository verification with `make test`, including Go golden tests, example smoke tests, documented-output validation, markdown-link checks, stale-term linting, validator tests, and registry drift checks.
- Added `spec/ahf-v0.1.json` as the machine-readable AHF registry and shape manifest.
- Added `CONFORMANCE.md`, `VERSIONING.md`, `SECURITY.md`, `NOTICE`, `LICENSE-DOCS`, and `CODE_OF_CONDUCT.md`.
- Added parsing guidance, prior-art notes, a release checklist, and a Cobra implementation guide under `docs/`.
- Added a non-normative grammar sketch, Click guide, Commander guide, and validator fixtures.
- Added argparse and Clap framework guides, adopter checklist, and fixture documentation.
- Added implementer FAQ, publication checklist, framework snippets, and static-site link checking.
- Added argparse and Clap snippets, static-site deployment notes, site metadata, and validator scope notes.
- Added spec table-of-contents checking, `make release-check`, and a visible README draft-status callout.
- Softened RFC requirement language and simplified the public documentation surface.
- Added CI and pull request/issue templates for public contribution flow.
- Updated the Go/Cobra reference implementation with cursor-aware `--agent-out`, safer quoted `next` commands, and golden tests.
- Regenerated examples from live CLI output.

## v0.1 draft

- Initial public draft of Agent Help Format (AHF).
- Defined `--agent-help` AH1/AH2 records for agent-facing CLI discovery.
- Defined `--agent-out` as an AHF protocol envelope with TOON result bodies.
- Defined record prefix and scalar type registries.
- Included an agent implementation brief, quick reference, examples, and build-time `.agents/skills/ahf/SKILL.md`.
