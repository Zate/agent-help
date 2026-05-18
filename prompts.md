# Agent Prompts

Use these prompts with an AI coding agent.

## Explain agent-help

```text
Read https://raw.githubusercontent.com/Zate/agent-help/main/llms-full.txt and explain what agent-help is for. Keep the explanation grounded in how an agent would use a CLI.
```

## Implement agent-help in a CLI

```text
Read https://raw.githubusercontent.com/Zate/agent-help/main/llms-full.txt and implement agent-help in this CLI. Preserve normal --help and --json behavior. Add --agent-help for invocation help and --agent-out for structured runtime results where appropriate.
```

## Use the skill file

```text
Use this skill: https://raw.githubusercontent.com/Zate/agent-help/main/.agents/skills/ahf/SKILL.md
```

## Inspect the AHF docs

```text
Read https://raw.githubusercontent.com/Zate/agent-help/main/agent-help.ahf and https://raw.githubusercontent.com/Zate/agent-help/main/agent-help-full.ahf. Use those AHF records to decide which agent-help docs to read next.
```

---

## Raw agent files

These files are designed to be fetched directly by agents. Point your agent at the raw URL — no HTML parsing required.

| File | Purpose | Raw URL |
|---|---|---|
| `llms.txt` | Short orientation (~55 lines) | https://raw.githubusercontent.com/Zate/agent-help/main/llms.txt |
| `llms-full.txt` | Full implementation brief (~311 lines) | https://raw.githubusercontent.com/Zate/agent-help/main/llms-full.txt |
| `.agents/skills/ahf/SKILL.md` | agentskills.io-format skill | https://raw.githubusercontent.com/Zate/agent-help/main/.agents/skills/ahf/SKILL.md |
| `agent-help.ahf` | AH1 index dogfood | https://raw.githubusercontent.com/Zate/agent-help/main/agent-help.ahf |
| `agent-help-full.ahf` | AH2 detail dogfood | https://raw.githubusercontent.com/Zate/agent-help/main/agent-help-full.ahf |
