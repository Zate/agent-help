# Prior Art and Related Work

agent-help overlaps with several existing CLI and tool-description surfaces. It is intended to complement them, not replace them.

## Human Help

`--help`, man pages, and README examples are optimized for people. They often include prose, layout, aliases, colors, and explanatory examples. Agents can read them, but they spend tokens on presentation and still have to infer required arguments, enum values, and follow-up commands.

agent-help keeps human help intact and adds a compact agent-facing discovery surface.

## JSON Output

`--json` is the right interface for deterministic software integrations. It is strict, widely supported, and easy to parse.

For agents, JSON has two gaps:

- It can be token-heavy for tabular output.
- It describes data, not how to discover commands or what to call next.

agent-help keeps `--json` for software and uses `--agent-out` for agent-facing runtime output.

## Shell Completion

Shell completion scripts expose valid tokens for interactive shells. They are useful but platform-specific and not designed as a command-discovery transcript.

agent-help focuses on typed invocation guidance, required values, examples, error recovery, and follow-up commands.

## Man Pages

Man pages are durable human documentation. They are not usually embedded in the binary output stream, and they do not define runtime result envelopes, pagination, or error-retry records.

agent-help is emitted directly by the CLI at the point an agent invokes it.

## MCP and Plugin Protocols

MCP servers, plugin protocols, and tool registries expose capabilities through adapters or RPC layers. They are powerful when an agent runtime supports them and when the additional deployment surface is acceptable.

agent-help makes the CLI itself agent-readable. A conforming CLI can still be wrapped by MCP, plugins, or skills; those wrappers can use `--agent-help` and `--agent-out` as a stable underlying surface.

## OpenAPI and JSON Schema

OpenAPI and JSON Schema describe HTTP APIs and structured data schemas. They are much broader and more formal than agent-help.

agent-help targets command-line invocation and stdout/stderr interaction: command discovery, arguments, flags, runtime status, TOON result bodies, and follow-up shell commands.

## CLI Framework Metadata

Frameworks such as Cobra, Click, Clap, Commander, and argparse already store command metadata. agent-help encourages implementers to derive AH1/AH2 output from that metadata rather than maintaining a separate command list.

