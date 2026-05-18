# Contributing to agent-help

agent-help is an open convention with an open draft specification (AHF). Contributions are welcome — whether that's clarifying the spec, adding examples, building reference implementations, or raising open questions.

## Ways to contribute

### Open an issue

The best way to contribute is to [open an issue](https://github.com/Zate/agent-help/issues/new/choose). We have templates for:

- **Open question** — raise a design question about the spec
- **Errata** — report an error or inconsistency in the spec
- **Feature proposal** — propose a new record type, flag, or convention
- **Implementation report** — share that you've implemented agent-help in a CLI

### Submit a pull request

For small changes (typos, clarifications, example improvements), a PR is welcome directly.

For substantive spec changes, please open an issue first so the change can be discussed before you invest time writing it up.

#### PR checklist

- [ ] `make test` passes locally
- [ ] Changes to `AHF-RFC.md` are reflected in `llms-full.txt` if they affect the implementation brief
- [ ] New examples follow the format in `examples/` (numbered, with source command and raw output)
- [ ] If adding a new AHF record prefix, update §9 (Record Prefix Registry) in `AHF-RFC.md`
- [ ] If adding a new scalar type, update §10 (Scalar Type Registry) in `AHF-RFC.md`
- [ ] If changing registries, update `spec/ahf-v0.1.json` and run `make check-drift`
- [ ] If changing public behavior, update `CHANGELOG.md`
- [ ] SKILL.md stays under 500 lines (move detail to `references/` if needed)

## Local Verification

```bash
make test
make demo
make verify-examples
make verify-doc-examples
make validate-ahf
make verify-fixtures
make lint-html-links
make check-spec-toc
make update-examples
```

## Open questions

The following design questions are explicitly open in the current draft. Each has a corresponding GitHub issue for discussion:

1. **Flag naming** — Should `--agent-out` remain as specified, or should a shorter alias be added?
2. **Formal grammar** — Should a formal ABNF grammar be added in v0.2?
3. **MIME type** — Should MIME type registration be pursued for an agent-help media type?
4. **Extension namespacing** — Should record prefix extension namespacing be defined?
5. **Conformance levels** — Should conformance levels distinguish `--agent-help`-only from full `--agent-help` + `--agent-out`?
6. **Test suite** — Should examples be packaged as a model comprehension test suite?

See [§21 of the spec](AHF-RFC.md#21-open-questions) for context on each question.

## Versioning

See [VERSIONING.md](VERSIONING.md) for the full draft stability and release policy. In short, AHF follows a simple versioning scheme:

- **Patch** (v0.1.x) — errata, clarifications, no normative changes
- **Minor** (v0.x) — new optional features, backward-compatible
- **Major** (vX.0) — breaking changes to the record format or required conventions

The current draft is **v0.1**. The core wire shape is stable enough for trial implementations, but registry entries, naming details, and conformance levels may change before v1.0.

## Scope

agent-help is specifically about **agent-native CLI interaction** — discovery, invocation help (AHF), and runtime output (TOON + AHF envelope). It is not:

- A general-purpose serialization format (that's TOON's job)
- A replacement for JSON between services
- An agent communication protocol or RPC layer

If your proposal is outside this scope, it may be better addressed in a separate companion spec. Open an issue to discuss.

## Code of conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## License

By contributing, you agree that your contributions will be licensed under the same terms as the project:
- Code and examples: [Apache 2.0](LICENSE)
- Documentation and specification: [CC-BY-4.0](LICENSE-DOCS)
