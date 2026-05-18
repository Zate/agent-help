#!/usr/bin/env python3
"""Check that AHF-RFC.md numbered TOC entries match numbered headings."""

from __future__ import annotations

import re
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SPEC = ROOT / "AHF-RFC.md"
TOC_ENTRY = re.compile(r"^(\d+)\.\s+(.+)$")
HEADING = re.compile(r"^##\s+(\d+)\.\s+(.+)$")


def normalize(title: str) -> str:
    title = re.sub(r"\s*\([^)]*\)\s*$", "", title.strip())
    return title.replace("`", "")


def main() -> int:
    lines = SPEC.read_text(encoding="utf-8").splitlines()

    try:
        start = lines.index("## Table of Contents") + 1
    except ValueError:
        print("AHF-RFC.md: missing Table of Contents heading", file=sys.stderr)
        return 1

    toc: list[tuple[str, str]] = []
    for line in lines[start:]:
        if line.startswith("## "):
            break
        match = TOC_ENTRY.match(line)
        if match:
            toc.append((match.group(1), normalize(match.group(2))))

    headings: list[tuple[str, str]] = []
    for line in lines:
        match = HEADING.match(line)
        if match:
            headings.append((match.group(1), normalize(match.group(2))))

    errors: list[str] = []
    if toc != headings:
        errors.append("AHF-RFC.md table of contents drift")
        errors.append(f"toc:      {toc}")
        errors.append(f"headings: {headings}")

    if errors:
        print("\n".join(errors), file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
