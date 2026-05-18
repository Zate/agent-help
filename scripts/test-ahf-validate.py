#!/usr/bin/env python3
"""Regression tests for scripts/ahf-validate.py."""

from __future__ import annotations

import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
VALIDATOR = ROOT / "scripts" / "ahf-validate.py"


def run_validator(text: str) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [sys.executable, str(VALIDATOR)],
        input=text,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        check=False,
    )


def expect_valid(name: str, text: str) -> list[str]:
    result = run_validator(text)
    if result.returncode == 0:
        return []
    return [f"{name}: expected valid, got {result.returncode}: {result.stderr.strip()}"]


def expect_invalid(name: str, text: str, needle: str) -> list[str]:
    result = run_validator(text)
    if result.returncode != 0 and needle in result.stderr:
        return []
    return [
        f"{name}: expected invalid containing {needle!r}",
        f"returncode: {result.returncode}",
        f"stderr: {result.stderr.strip()}",
    ]


def main() -> int:
    errors: list[str] = []

    errors.extend(
        expect_valid(
            "valid AH2",
            """ah2 mem node list
use mem node list [--type TYPE] [--limit int]
flag --type:enum(decision|fact) opt :: filter by type
flag --limit:int opt default=20 :: max rows
ex mem node list --type fact
""",
        )
    )
    errors.extend(
        expect_valid(
            "valid paginated agent-out",
            """ok nodes count=8 shown=3 more=1 cursor=c_n_104
warn truncated shown=3 total=8
nodes[#3]{id,type,text}:
  n_101,fact,alpha
  n_102,fact,beta
  n_103,fact,gamma
next "mem node list --limit 3 --cursor c_n_104 --agent-out"
""",
        )
    )
    errors.extend(
        expect_invalid(
            "AH2 missing use",
            """ah2 mem node list
flag --limit:int opt :: max rows
""",
            "AH2 output must include a use record",
        )
    )
    errors.extend(
        expect_invalid(
            "unknown scalar",
            """ah2 mem node list
use mem node list [--limit int]
flag --limit:integer opt :: max rows
""",
            "unknown flag scalar type",
        )
    )
    errors.extend(
        expect_invalid(
            "leading agent-help example",
            """ah2 mem node list
use mem node list
ex mem --agent-help node list
""",
            "leading --agent-help placement",
        )
    )
    errors.extend(
        expect_invalid(
            "pagination needs next",
            """ok nodes count=8 shown=3 more=1 cursor=c_n_104
nodes[#3]{id,type,text}:
  n_101,fact,alpha
""",
            "must include a next continuation",
        )
    )

    if errors:
        print("\n".join(errors), file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
