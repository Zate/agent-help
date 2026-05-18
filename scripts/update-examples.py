#!/usr/bin/env python3
"""Regenerate example raw output blocks from the live mem reference CLI."""

from __future__ import annotations

import os
import subprocess
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
GO_ENV = {
    **os.environ,
    "GOTOOLCHAIN": "local",
    "GOCACHE": str(ROOT / ".cache" / "go-build"),
    "GOMODCACHE": str(ROOT / ".cache" / "go-mod"),
}

EXAMPLES = {
    "examples/01-agent-help-index-simple.md": ["--agent-help"],
    "examples/02-agent-help-command-detail.md": ["node", "add", "--agent-help"],
    "examples/03-agent-help-error-hint.md": ["node", "add", "postgres 15 required", "--agent-out"],
    "examples/04-agent-out-list-results.md": ["node", "list", "--agent-out"],
    "examples/05-agent-out-single-object.md": ["project", "status", "--agent-out"],
    "examples/06-agent-out-paginated-warning.md": ["node", "list", "--limit", "3", "--agent-out"],
}


def run_mem(args: list[str]) -> str:
    result = subprocess.run(
        ["go", "run", ".", *args],
        cwd=ROOT / "impl",
        env=GO_ENV,
        check=True,
        text=True,
        stdout=subprocess.PIPE,
    )
    return result.stdout.rstrip("\n")


def replace_first_raw_text_block(markdown: str, replacement: str) -> str:
    lines = markdown.splitlines()
    saw_raw = False
    start = None
    end = None

    for idx, line in enumerate(lines):
        if line.startswith("Raw"):
            saw_raw = True
            continue
        if saw_raw and line == "```text":
            start = idx + 1
            continue
        if start is not None and line == "```":
            end = idx
            break

    if start is None or end is None:
        raise ValueError("could not find first Raw text block")

    new_lines = lines[:start] + replacement.splitlines() + lines[end:]
    return "\n".join(new_lines) + "\n"


def main() -> int:
    for rel_path, args in EXAMPLES.items():
        path = ROOT / rel_path
        current = path.read_text(encoding="utf-8")
        output = run_mem(args)
        updated = replace_first_raw_text_block(current, output)
        path.write_text(updated, encoding="utf-8")
        print(f"updated {rel_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

