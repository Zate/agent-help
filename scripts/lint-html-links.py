#!/usr/bin/env python3
"""Minimal local HTML link checker for the static site."""

from __future__ import annotations

import re
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SITE = ROOT / "site"
HREF = re.compile(r"""href=(["'])(.*?)\1""")


def is_external(target: str) -> bool:
    return (
        "://" in target
        or target.startswith("mailto:")
        or target.startswith("#")
    )


def candidates(path: Path, raw: str) -> list[Path]:
    target, _, _anchor = raw.partition("#")
    if not target:
        return []

    primary = (path.parent / target).resolve()
    result = [primary]

    # GitHub Pages serves the repository root at the site root. The source site
    # lives in site/, so root-doc links are valid even when not present in site/.
    site_relative = (SITE / target).resolve()
    root_relative = (ROOT / target).resolve()
    if site_relative == primary and root_relative != primary:
        result.append(root_relative)

    return result


def check_file(path: Path) -> list[str]:
    errors: list[str] = []
    text = path.read_text(encoding="utf-8")

    for lineno, line in enumerate(text.splitlines(), start=1):
        for match in HREF.finditer(line):
            raw = match.group(2).strip()
            if not raw or is_external(raw):
                continue
            resolved = candidates(path, raw)
            if not resolved:
                continue
            for candidate in resolved:
                try:
                    candidate.relative_to(ROOT)
                except ValueError:
                    errors.append(f"{path}:{lineno}: link escapes repository: {raw}")
                    break
            else:
                if not any(candidate.exists() for candidate in resolved):
                    errors.append(f"{path}:{lineno}: missing link target: {raw}")

    return errors


def main() -> int:
    errors: list[str] = []
    for path in sorted(SITE.rglob("*.html")):
        errors.extend(check_file(path))

    if errors:
        print("\n".join(errors), file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
