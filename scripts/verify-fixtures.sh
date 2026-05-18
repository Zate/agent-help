#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

for file in "$repo_root"/tests/fixtures/valid/*.ahf; do
  "$repo_root/scripts/ahf-validate.py" "$file"
done

for file in "$repo_root"/tests/fixtures/invalid/*.ahf; do
  if "$repo_root/scripts/ahf-validate.py" "$file" >/tmp/agent-help-fixture.out 2>/tmp/agent-help-fixture.err; then
    printf 'invalid fixture passed validation: %s\n' "$file" >&2
    exit 1
  fi
done

printf 'fixtures validate expected pass/fail behavior\n'
