#!/usr/bin/env python3
"""Lightweight AHF/agent-out validator for examples and CLI output."""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path


PREFIXES = {
    "ah1",
    "ah2",
    "ok",
    "err",
    "warn",
    "cmd",
    "use",
    "arg",
    "flag",
    "ex",
    "hint",
    "next",
    "more?",
}

TOON_HEADER = re.compile(r"^[A-Za-z_][A-Za-z0-9_]*\[#\d+\]\{[A-Za-z0-9_,]+\}:$")
ARG_RE = re.compile(r"^arg\s+[A-Za-z0-9_-]+:([^ ]+)\s+(req|opt)\s+::\s+.+$")
FLAG_RE = re.compile(r"^flag\s+--[A-Za-z0-9_-]+:([^ ]+)\s+(req|opt|repeat)(\s+default=[^ ]+)?\s+::\s+.+$")
SCALAR_TYPES = {"str", "int", "num", "bool", "path", "url", "id", "ts", "date", "dur", "kv"}


def validate_text(name: str, text: str) -> list[str]:
    errors: list[str] = []
    lines = text.splitlines()
    records = [(idx + 1, line) for idx, line in enumerate(lines) if line.strip()]

    if not records:
        return [f"{name}: empty input"]

    first = records[0][1].split(maxsplit=1)[0]
    if first not in {"ah1", "ah2", "ok", "err"}:
        errors.append(f"{name}:1: first AHF record should be ah1, ah2, ok, or err, got {first!r}")

    has_more = False
    has_next = False
    has_toon = False
    has_use = False
    saw_toon = False

    for lineno, line in records:
        token = line.split(maxsplit=1)[0]

        if token in PREFIXES:
            pass
        elif TOON_HEADER.match(line) or line.startswith("  "):
            has_toon = True
            saw_toon = True
            continue
        elif first == "ok" and re.match(r"^[a-z][a-z0-9_-]*\s+", line):
            # AHF key-value fallback for scalar acknowledgements.
            continue
        else:
            errors.append(f"{name}:{lineno}: unknown AHF/TOON line prefix {token!r}")
            continue

        if token == "warn":
            if first != "ok":
                errors.append(f"{name}:{lineno}: warn records should only appear after ok")
            if saw_toon:
                errors.append(f"{name}:{lineno}: warn records should appear before the TOON body")
        if token == "use":
            has_use = True
        if token == "more?":
            has_more = True
            if "--agent-help" not in line:
                errors.append(f"{name}:{lineno}: more? should point to a --agent-help command shape")
        if token == "next":
            has_next = True
            if has_more and first == "ok" and "--agent-out" not in line:
                errors.append(f"{name}:{lineno}: continuation next after more=1 should include --agent-out")
        if token == "ex" and re.search(r"\s--agent-help\s+\S", line):
            errors.append(f"{name}:{lineno}: ex should not document leading --agent-help placement")
        if token == "arg":
            match = ARG_RE.match(line)
            if not match:
                errors.append(f"{name}:{lineno}: malformed arg record")
            elif not valid_type(match.group(1)):
                errors.append(f"{name}:{lineno}: unknown arg scalar type {match.group(1)!r}")
        if token == "flag":
            match = FLAG_RE.match(line)
            if not match:
                errors.append(f"{name}:{lineno}: malformed flag record")
            elif not valid_type(match.group(1)):
                errors.append(f"{name}:{lineno}: unknown flag scalar type {match.group(1)!r}")

        if token == "ok":
            if re.search(r"\b(more=1|truncated=1)\b", line):
                has_more = True
        if token == "err":
            if lineno == 1:
                continue

    if first == "ah1" and not has_more:
        errors.append(f"{name}: AH1 output should include a more? pointer")

    if first == "ah2" and not has_use:
        errors.append(f"{name}: AH2 output must include a use record")

    if has_more and first == "ok" and not has_next:
        errors.append(f"{name}: paginated/truncated ok output must include a next continuation")

    if first == "ok" and not (has_toon or any(line.split(maxsplit=1)[0] not in PREFIXES for _, line in records[1:])):
        errors.append(f"{name}: ok output should include a TOON body or scalar fallback lines")

    return errors


def valid_type(value: str) -> bool:
    if value in SCALAR_TYPES:
        return True
    return bool(re.fullmatch(r"enum\([A-Za-z0-9_.:-]+(\|[A-Za-z0-9_.:-]+)+\)", value))


def read_inputs(paths: list[str]) -> list[tuple[str, str]]:
    if not paths or paths == ["-"]:
        return [("<stdin>", sys.stdin.read())]
    return [(path, Path(path).read_text(encoding="utf-8")) for path in paths]


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate lightweight AHF/agent-out output.")
    parser.add_argument("paths", nargs="*", help="files to validate, or stdin when omitted")
    args = parser.parse_args()

    errors: list[str] = []
    for name, text in read_inputs(args.paths):
        errors.extend(validate_text(name, text))

    if errors:
        print("\n".join(errors), file=sys.stderr)
        return 1

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
