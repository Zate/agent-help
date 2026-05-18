#!/usr/bin/env python3
"""Minimal local Markdown link checker for repository docs."""

from __future__ import annotations

import re
import sys
from pathlib import Path


LINK = re.compile(r"(?<!!)\[[^\]]+\]\(([^)]+)\)")
ROOT = Path(__file__).resolve().parents[1]


def is_external(target: str) -> bool:
    return (
        "://" in target
        or target.startswith("mailto:")
        or target.startswith("#")
        or target.startswith("../../issues")
    )


def check_file(path: Path) -> list[str]:
    errors: list[str] = []
    text = path.read_text(encoding="utf-8")

    for lineno, line in enumerate(text.splitlines(), start=1):
        for match in LINK.finditer(line):
            raw = match.group(1).strip()
            if not raw or is_external(raw):
                continue
            target, _, anchor = raw.partition("#")
            if not target:
                continue
            resolved = (path.parent / target).resolve()
            try:
                resolved.relative_to(ROOT)
            except ValueError:
                errors.append(f"{path}:{lineno}: link escapes repository: {raw}")
                continue
            if not resolved.exists():
                errors.append(f"{path}:{lineno}: missing link target: {raw}")
            elif anchor and resolved.suffix.lower() == ".md":
                # Keep anchor validation intentionally simple; full GitHub slug parity is overkill here.
                continue

    return errors


def main() -> int:
    ignored = {".git", ".cache", "_site"}
    files = sorted(p for p in ROOT.rglob("*.md") if not ignored.intersection(p.parts))
    errors: list[str] = []
    for path in files:
        errors.extend(check_file(path))

    if errors:
        print("\n".join(str(e) for e in errors), file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
