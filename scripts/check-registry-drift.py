#!/usr/bin/env python3
"""Check that AHF registries stay synchronized across public docs."""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SPEC = ROOT / "spec" / "ahf-v0.1.json"
RFC = ROOT / "AHF-RFC.md"
REFERENCE = ROOT / "references" / "REFERENCE.md"
LLMS_FULL = ROOT / "llms-full.txt"


def normalize(value: str) -> str:
    return value.strip().strip("`").replace(r"\|", "|")


def split_markdown_row(line: str) -> list[str]:
    cells: list[str] = []
    current: list[str] = []
    escaped = False
    for char in line.strip().strip("|"):
        if escaped:
            current.append("\\" + char)
            escaped = False
            continue
        if char == "\\":
            escaped = True
            continue
        if char == "|":
            cells.append("".join(current).strip())
            current = []
            continue
        current.append(char)
    if escaped:
        current.append("\\")
    cells.append("".join(current).strip())
    return cells


def extract_table_values(path: Path, heading: str, column_index: int) -> list[str]:
    text = path.read_text(encoding="utf-8")
    start = text.index(heading)
    section = text[start:]
    next_heading = re.search(r"\n##\s+", section[len(heading) :])
    if next_heading:
        section = section[: len(heading) + next_heading.start()]

    values: list[str] = []
    for line in section.splitlines():
        if not line.startswith("|"):
            continue
        cells = split_markdown_row(line)
        if not cells or cells[0] in {"---", "Prefix", "Type"}:
            continue
        if len(cells) <= column_index:
            continue
        values.append(normalize(cells[column_index]))
    return values


def extract_bullets(path: Path, heading: str) -> list[str]:
    text = path.read_text(encoding="utf-8")
    lines = text[text.index(heading) :].splitlines()
    values: list[str] = []
    started = False
    for line in lines[1:]:
        match = re.match(r"-\s+([^:]+):", line)
        if match:
            started = True
            values.append(normalize(match.group(1)))
        elif started:
            break
    return values


def compare(name: str, expected: list[str], actual: list[str], source: str) -> list[str]:
    if expected == actual:
        return []
    return [
        f"{name} drift in {source}",
        f"expected: {expected}",
        f"actual:   {actual}",
    ]


def main() -> int:
    spec = json.loads(SPEC.read_text(encoding="utf-8"))
    expected_prefixes = [row["prefix"] for row in spec["record_prefixes"]]
    expected_types = [row["type"] for row in spec["scalar_types"]]

    errors: list[str] = []
    errors.extend(
        compare(
            "record prefix registry",
            expected_prefixes,
            extract_table_values(RFC, "## 9. Record Prefix Registry", 0),
            "AHF-RFC.md",
        )
    )
    errors.extend(
        compare(
            "scalar type registry",
            expected_types,
            extract_table_values(RFC, "## 10. Scalar Type Registry", 0),
            "AHF-RFC.md",
        )
    )
    errors.extend(
        compare(
            "record prefix registry",
            expected_prefixes,
            extract_table_values(REFERENCE, "## AHF record prefixes", 0),
            "references/REFERENCE.md",
        )
    )
    errors.extend(
        compare(
            "scalar type registry",
            expected_types,
            extract_table_values(REFERENCE, "## Scalar types", 0),
            "references/REFERENCE.md",
        )
    )
    errors.extend(
        compare(
            "record prefix registry",
            expected_prefixes,
            extract_bullets(LLMS_FULL, "AHF record prefixes:"),
            "llms-full.txt",
        )
    )
    errors.extend(
        compare(
            "scalar type registry",
            expected_types,
            extract_bullets(LLMS_FULL, "Scalar types"),
            "llms-full.txt",
        )
    )

    if errors:
        print("\n".join(errors), file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
