#!/usr/bin/env bash
set -euo pipefail

out="${1:-_site}"

mkdir -p "$out"
cp -R site/. "$out/"

cp \
  AHF-RFC.md \
  CHANGELOG.md \
  CODE_OF_CONDUCT.md \
  CONTRIBUTING.md \
  FAQ.md \
  LICENSE \
  LICENSE-DOCS \
  NOTICE \
  README.md \
  SECURITY.md \
  agent-help.ahf \
  agent-help-full.ahf \
  llms-full.txt \
  llms.txt \
  prompts.md \
  "$out/"

mkdir -p "$out/docs" "$out/examples" "$out/spec" "$out/tests"
cp -R docs/. "$out/docs/"
cp -R examples/. "$out/examples/"
cp -R spec/. "$out/spec/"
cp -R tests/. "$out/tests/"

# Copy skill files so agents can fetch them at their canonical path
mkdir -p "$out/.agents/skills/ahf/references"
cp .agents/skills/ahf/SKILL.md "$out/.agents/skills/ahf/SKILL.md"
cp .agents/skills/ahf/references/REFERENCE.md "$out/.agents/skills/ahf/references/REFERENCE.md"

touch "$out/.nojekyll"
