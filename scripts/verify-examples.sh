#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$repo_root/impl"

run() {
  GOTOOLCHAIN=local \
    GOCACHE="$repo_root/.cache/go-build" \
    GOMODCACHE="$repo_root/.cache/go-mod" \
    go run . "$@"
}

assert_contains() {
  local output="$1"
  local expected="$2"
  if [[ "$output" != *"$expected"* ]]; then
    printf 'missing expected text:\n%s\n\noutput:\n%s\n' "$expected" "$output" >&2
    exit 1
  fi
}

help_output="$(run --help)"
assert_contains "$help_output" "LLM agent? Use --agent-help for token-optimized usage."

index_output="$(run --agent-help)"
assert_contains "$index_output" "ah1 mem ::"
assert_contains "$index_output" "cmd node list [--type TYPE] [--limit int] [--cursor id]"
assert_contains "$index_output" "more? mem <cmd> --agent-help"

detail_output="$(run node list --agent-help)"
assert_contains "$detail_output" "ah2 mem node list"
assert_contains "$detail_output" "flag --cursor:id opt :: resume after the given node ID"

list_output="$(run node list --agent-out)"
assert_contains "$list_output" "ok nodes count=8 shown=8 more=0"
assert_contains "$list_output" "nodes[#8]{id,type,tags,text}:"

page_output="$(run node list --limit 3 --agent-out)"
assert_contains "$page_output" "ok nodes count=8 shown=3 more=1 cursor=c_n_104"
assert_contains "$page_output" 'next "mem node list --limit 3 --cursor c_n_104 --agent-out"'

error_output="$(run node add "test node" --type badval --agent-out)"
assert_contains "$error_output" "err invalid_enum flag=--type got=badval"
assert_contains "$error_output" "use mem node add <text> --type TYPE"

printf 'examples verified\n'
