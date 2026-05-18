# Release Checklist

Use this checklist before tagging or publicly announcing a repository release.

## Preflight

1. Confirm the intended AHF draft version.
2. Review open issues for release blockers.
3. Confirm `docs/VERSIONING.md` still describes the release policy.

## Sync Checks

1. If `AHF-RFC.md` changed, verify corresponding updates in:
   - `llms.txt`
   - `llms-full.txt`
   - `.agents/skills/ahf/SKILL.md`
   - `.agents/skills/ahf/references/REFERENCE.md`
   - `examples/`
   - `spec/ahf-v0.1.json`
   - `docs/CONFORMANCE.md`
2. If record prefixes changed, update every registry table and manifest.
3. If scalar types changed, update every scalar type table and manifest.
4. If `--agent-out` shape changed, update examples and golden tests.

## Verification

Run:

```bash
make test
make demo
make verify-examples
make verify-doc-examples
make validate-ahf
make verify-fixtures
make check-drift
make check-spec-toc
make lint-html-links
make release-check
```

Expected:

- Go golden tests pass.
- Example smoke checks pass.
- Documented example raw blocks match the live `mem` CLI.
- AHF examples pass the lightweight validator.
- Conformance fixtures pass or fail as expected.
- Registry tables match the machine-readable manifest.
- `AHF-RFC.md` table of contents matches numbered headings.
- Markdown links resolve locally.
- Static site links resolve locally.
- Stale-term lint passes.

## Release Notes

Include:

- AHF draft version.
- Compatibility notes.
- New record prefixes or scalar types.
- Reference implementation changes.
- Known open questions.
- `CHANGELOG.md` is updated.

## Tagging

Use semantic-style repository tags while pre-1.0:

```bash
git tag v0.1.x
git push origin v0.1.x
```

Use `v0.2.0` for backward-compatible spec additions and `v1.0.0` for the first stable release.
