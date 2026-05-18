# Test Fixtures

This directory contains lightweight conformance fixtures for AHF and `--agent-out` examples.

## Layout

- `fixtures/valid/` — snippets that must pass `scripts/ahf-validate.py`
- `fixtures/invalid/` — snippets that must fail `scripts/ahf-validate.py`

Run:

```bash
make verify-fixtures
```

The fixture suite is intentionally small. It protects the project from accidental validator regressions and documents representative protocol mistakes.

## Adding Valid Fixtures

Add valid fixtures when a protocol shape should be accepted by the lightweight validator.

Good candidates:

- AH1 command index output
- AH2 command detail output
- `ok` output with TOON body
- `err` output with `hint` and `use`
- pagination with `next`

## Adding Invalid Fixtures

Add invalid fixtures when a common implementation mistake should stay rejected.

Good candidates:

- AH1 without `more?`
- AH2 without `use`
- unknown scalar types
- leading `--agent-help` examples
- paginated `ok` output without `next`
- `warn` records after the TOON body

Invalid fixtures only assert that validation fails. If the exact error text matters, add a case to `scripts/test-ahf-validate.py`.
