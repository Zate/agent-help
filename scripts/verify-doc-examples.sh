#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
validate_only=0
if [[ "${1:-}" == "--validate-only" ]]; then
  validate_only=1
fi

extract_first_text_block_after_raw() {
  local file="$1"
  awk '
    /^Raw/ { seen_raw=1; next }
    seen_raw && /^```text$/ { in_block=1; next }
    in_block && /^```$/ { exit }
    in_block { print }
  ' "$file"
}

run_mem() {
  local command="$1"
  local args="${command#mem }"
  (
    cd "$repo_root/impl"
    GOTOOLCHAIN=local \
      GOCACHE="$repo_root/.cache/go-build" \
      GOMODCACHE="$repo_root/.cache/go-mod" \
      bash -lc "go run . $args"
  )
}

check_example() {
  local file="$1"
  local command="$2"
  local expected actual

  expected="$(extract_first_text_block_after_raw "$repo_root/$file")"
  actual="$(run_mem "$command")"

  if [[ "$actual" != "$expected" ]]; then
    printf 'example mismatch: %s\ncommand: %s\n\nexpected:\n%s\n\nactual:\n%s\n' \
      "$file" "$command" "$expected" "$actual" >&2
    exit 1
  fi

  printf '%s\n' "$expected" | "$repo_root/scripts/ahf-validate.py" -
}

validate_example() {
  local file="$1"
  extract_first_text_block_after_raw "$repo_root/$file" | "$repo_root/scripts/ahf-validate.py" -
}

if [[ "$validate_only" == "1" ]]; then
  validate_example "examples/01-agent-help-index-simple.md"
  validate_example "examples/02-agent-help-command-detail.md"
  validate_example "examples/03-agent-help-error-hint.md"
  validate_example "examples/04-agent-out-list-results.md"
  validate_example "examples/05-agent-out-single-object.md"
  validate_example "examples/06-agent-out-paginated-warning.md"
  printf 'documented examples validate as AHF/agent-out\n'
  exit 0
fi

check_example "examples/01-agent-help-index-simple.md" "mem --agent-help"
check_example "examples/02-agent-help-command-detail.md" "mem node add --agent-help"
check_example "examples/03-agent-help-error-hint.md" 'mem node add "postgres 15 required" --agent-out'
check_example "examples/04-agent-out-list-results.md" "mem node list --agent-out"
check_example "examples/05-agent-out-single-object.md" "mem project status --agent-out"
check_example "examples/06-agent-out-paginated-warning.md" "mem node list --limit 3 --agent-out"

printf 'documented examples match CLI output\n'
