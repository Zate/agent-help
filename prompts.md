# Prompts

These prompts are for **implementing** agent-help in a CLI, not for using an already agent-help-enabled CLI.

An agent using an enabled CLI should discover everything from the CLI itself:

```text
tool --help
tool --agent-help
tool subcmd --agent-help
tool subcmd args --agent-out
```

## Explain the idea to a human

```text
Explain agent-help in plain language. Treat it as MCP-lite for CLIs: two flags, --agent-help for compact command discovery and --agent-out for dense structured results. Make clear that the spec and llms-full.txt are implementation aids, not runtime instructions an agent needs before using an enabled CLI.
```

## Implement agent-help in this CLI

```text
Read https://raw.githubusercontent.com/Zate/agent-help/main/llms-full.txt and implement agent-help in this CLI. Preserve normal --help and --json behavior. Add --agent-help for compact command discovery and --agent-out for structured runtime results where appropriate. Generate outputs from existing command metadata where possible.
```

## Use the optional build-time skill

```text
Use this implementation skill while modifying the CLI: https://raw.githubusercontent.com/Zate/agent-help/main/.agents/skills/ahf/SKILL.md
```

## Sanity-check an implementation

```text
Check this CLI's agent-help implementation. Verify that tool --help points to --agent-help, tool --agent-help returns a compact command index, every listed command supports tool subcmd --agent-help, examples run as written, and structured result commands support --agent-out where appropriate.
```

## Raw implementation files

These files are for humans and coding agents adding agent-help support to a CLI.

| File | Purpose | Raw URL |
|---|---|---|
| `llms.txt` | Short orientation | https://raw.githubusercontent.com/Zate/agent-help/main/llms.txt |
| `llms-full.txt` | Full implementation brief | https://raw.githubusercontent.com/Zate/agent-help/main/llms-full.txt |
| `.agents/skills/ahf/SKILL.md` | Optional build-time skill | https://raw.githubusercontent.com/Zate/agent-help/main/.agents/skills/ahf/SKILL.md |
| `agent-help.ahf` | Dogfood AH1 example for these docs | https://raw.githubusercontent.com/Zate/agent-help/main/agent-help.ahf |
| `agent-help-full.ahf` | Dogfood AH2 examples for these docs | https://raw.githubusercontent.com/Zate/agent-help/main/agent-help-full.ahf |
