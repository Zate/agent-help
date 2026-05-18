# Deployment

The public site is intended to publish at:

```text
https://zate.github.io/agent-help/
```

Use GitHub Pages with **Source: GitHub Actions**.

## GitHub Settings

1. Open the repository on GitHub.
2. Go to `Settings -> Pages`.
3. Set `Build and deployment -> Source` to `GitHub Actions`.
4. Push to `main` or run the `Pages` workflow manually.

The workflow in `.github/workflows/pages.yml` builds `_site/` with:

```bash
./scripts/build-site.sh _site
```

The artifact includes the static landing page plus the root docs that the site links to, including `AHF-RFC.md`, `llms.txt`, `llms-full.txt`, `SKILL.md`, examples, references, and AHF docs.

## Local Check

```bash
make build-site
make release-check
```

`_site/` is generated output and is ignored by Git.
