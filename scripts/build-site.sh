#!/usr/bin/env bash
set -euo pipefail

out="${1:-_site}"

mkdir -p "$out"
cp -R site/. "$out/"

cp \
  AHF-RFC.md \
  CHANGELOG.md \
  CODE_OF_CONDUCT.md \
  CONFORMANCE.md \
  CONTRIBUTING.md \
  FAQ.md \
  LICENSE \
  LICENSE-DOCS \
  NOTICE \
  README.md \
  SECURITY.md \
  SKILL.md \
  VERSIONING.md \
  agent-help.ahf \
  agent-help-full.ahf \
  llms-full.txt \
  llms.txt \
  prompts.md \
  "$out/"

mkdir -p "$out/docs" "$out/examples" "$out/references" "$out/spec" "$out/tests"
cp -R docs/. "$out/docs/"
cp -R examples/. "$out/examples/"
cp -R references/. "$out/references/"
cp -R spec/. "$out/spec/"
cp -R tests/. "$out/tests/"

touch "$out/.nojekyll"
