# Adopter Checklist

Use this checklist when adding agent-help to an existing CLI.

## Discovery

- [ ] Human `--help` points agents to `--agent-help`.
- [ ] `tool --agent-help` emits AH1.
- [ ] `tool subcmd --agent-help` emits AH2.
- [ ] `--agent-help` is trailing in docs and examples.
- [ ] Leading `tool --agent-help subcmd` is not documented as the canonical form.

## AH1

- [ ] AH1 starts with `ah1 <tool> :: <purpose>`.
- [ ] AH1 lists primary invocable commands with `cmd`.
- [ ] AH1 omits aliases, decorative prose, version banners, and full flag detail.
- [ ] AH1 includes `more? <tool> <cmd> --agent-help`.

## AH2

- [ ] AH2 starts with `ah2 <tool> <command-path>`.
- [ ] AH2 includes one `use` record.
- [ ] Required positional args are listed with `arg`.
- [ ] Required flags and high-value optional flags are listed with `flag`.
- [ ] Scalar types use the AHF registry.
- [ ] Every `ex` record is tested or omitted.

## Runtime Output

- [ ] Structured result commands support `--agent-out`.
- [ ] Success output starts with `ok`.
- [ ] Error output starts with `err`.
- [ ] Structured success bodies use TOON.
- [ ] Warnings use `warn` before the TOON body.
- [ ] Follow-up actions use exact `next` commands.
- [ ] Pagination or truncation uses `more=1` or `truncated=1` plus `next`.

## Safety

- [ ] Agent-facing output redacts secrets.
- [ ] Agent-facing output avoids unnecessary personal data.
- [ ] Error messages give direct correction hints without leaking internals.
- [ ] `--agent-out` does not emit color, markdown, human tables, or prose paragraphs.

## Verification

- [ ] Golden tests cover AH1 and important AH2 outputs.
- [ ] Tests cover representative `--agent-out` success output.
- [ ] Tests cover retryable `err` output.
- [ ] Pagination or truncation continuations are tested.
- [ ] Examples are generated from, or checked against, live command output where possible.
