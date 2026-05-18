.PHONY: test release-check build-site test-go demo update-examples verify-examples verify-doc-examples validate-ahf test-validator verify-fixtures check-drift check-spec-toc lint-docs lint-links lint-html-links

GO_ENV = GOTOOLCHAIN=local GOCACHE=$(CURDIR)/.cache/go-build GOMODCACHE=$(CURDIR)/.cache/go-mod

test: test-go verify-examples verify-doc-examples validate-ahf test-validator verify-fixtures check-drift lint-docs lint-links lint-html-links

release-check: test demo check-spec-toc

build-site:
	./scripts/build-site.sh _site

test-go:
	cd impl && $(GO_ENV) go test ./...

demo:
	cd impl && $(GO_ENV) go run . --help
	cd impl && $(GO_ENV) go run . --agent-help
	cd impl && $(GO_ENV) go run . node list --agent-out
	cd impl && $(GO_ENV) go run . node list --limit 3 --agent-out

update-examples:
	./scripts/update-examples.py

verify-examples:
	./scripts/verify-examples.sh

verify-doc-examples:
	./scripts/verify-doc-examples.sh

validate-ahf:
	./scripts/verify-doc-examples.sh --validate-only
	./scripts/ahf-validate.py agent-help.ahf
	./scripts/ahf-validate.py agent-help-full.ahf

test-validator:
	./scripts/test-ahf-validate.py

verify-fixtures:
	./scripts/verify-fixtures.sh

check-drift:
	./scripts/check-registry-drift.py

check-spec-toc:
	./scripts/check-spec-toc.py

lint-docs:
	@! grep -RInE 'AOF|AOF-RFC|Agent Output Format|\bctx\b|zberg|mmiller|aziz' \
		README.md AHF-RFC.md FAQ.md CONTRIBUTING.md SKILL.md llms.txt llms-full.txt \
		references examples impl site .github AGENTS.md CONFORMANCE.md VERSIONING.md spec scripts docs

lint-links:
	./scripts/lint-markdown-links.py

lint-html-links:
	./scripts/lint-html-links.py
