# agent-help — FAQ

## What is agent-help?

agent-help is MCP-lite for command-line tools.

Instead of building an MCP server, plugin, or skill just so an AI agent can use your CLI, you add two flags:

- `--agent-help` — compact command discovery for agents
- `--agent-out` — compact structured results for agents

The CLI remains the source of truth. Humans still use `--help`. Scripts still use `--json`. Agents get a tiny native surface that is easy to discover and cheap to keep in context.

## What problem does it solve?

Agents are surprisingly bad at using ordinary CLIs from human help text. They parse prose, guess required flags, miss enum values, and waste context on tables and decoration.

agent-help gives them the part they actually need:

```text
ah1 tool :: example data CLI
cmd search <query> [--limit int] :: search records
cmd show <id> :: display one record
more? tool <cmd> --agent-help
```

Then, for a specific command:

```text
ah2 tool search
use tool search <query> [--limit int]
arg query:str req :: search text
flag --limit:int opt default=20 :: max rows
ex tool search postgres --limit 5
```

No prose. No guessing. No separate docs at runtime.

## Is this a replacement for MCP?

No. MCP is great when you need a real tool server, resources, long-lived sessions, auth mediation, or a protocol boundary.

agent-help is for the much smaller case: “I already have a CLI; I want agents to use it reliably.” It avoids the daemon, SDK, packaging, and deployment work of a full MCP integration.

You can still build an MCP server on top of an agent-help-enabled CLI later.

## Is this a replacement for skills?

No. Skills are useful for teaching an agent a workflow or domain.

But an agent should not need a skill just to learn how to invoke `tool search <query> --limit 5`. That knowledge belongs in the CLI itself.

This repository includes a `SKILL.md`, but it is only a build-time helper for coding agents that are implementing agent-help in a CLI. It is not required for agents that use an already-enabled CLI.

## Why not just use `--help`?

Human help is optimized for people: prose, sections, alignment, color, aliases, and examples. Agents can read it, but they waste tokens and still have to infer structure.

`--agent-help` is the same command truth in a smaller shape:

```text
cmd deploy <env> [--dry-run bool] :: deploy service
```

That one line tells an agent the command, required positional argument, useful optional flag, and purpose.

## Why not just use JSON?

Keep `--json` for software. agent-help does not replace it.

JSON is verbose for LLM context and does not solve command discovery. A JSON result can describe data, but it usually does not tell the agent:

- what command to run next
- whether output was truncated
- how to recover from an invalid invocation
- which flags or enum values are valid

`--agent-out` adds a small protocol envelope around dense data:

```text
ok issues count=2 more=1
issues[#2]{id,status,title}:
  BUG-1,open,login fails
  BUG-2,triage,slow search
next page "tool issues --cursor abc --agent-out"
```

## What is AHF?

AHF means Agent Help Format. It is the small line-record format used by `--agent-help` and by the `ok` / `err` / `warn` / `next` / `hint` envelope around `--agent-out`.

You mostly do not need to think about the name. The important part is: one useful record per line, minimal punctuation, exact commands when follow-up is needed.

## What is TOON?

TOON is the recommended dense data encoding for `--agent-out` result bodies.

AHF says whether the command succeeded and what to do next. TOON carries rows or objects.

```text
ok users count=2 more=0
users[#2]{id,name,role}:
  u1,Ada,admin
  u2,Lin,viewer
```

## What does an agent need before using an agent-help CLI?

Nothing special.

The expected path is:

1. Agent runs `tool --help`.
2. It sees `LLM agent? Use --agent-help for token-optimized usage.`
3. It runs `tool --agent-help`.
4. It follows the records from there.

If you have to give the agent this website, the spec, `llms-full.txt`, or a skill before it can use the CLI, the implementation has missed the point.

## Then why does this repo have a spec and `llms-full.txt`?

Because implementers need guidance.

Those files teach a human or coding agent how to add the two flags to a CLI and keep the output consistent. They are not meant to be runtime documentation for agents using your CLI.

## How much do I need to implement?

Start small:

1. Add the breadcrumb to `--help`.
2. Add `tool --agent-help` for a command index.
3. Add `tool subcmd --agent-help` for command detail.
4. Add `--agent-out` only where structured results matter.

You can stop after discovery and still get value.

## Does it require a library?

No. It is plain text. Generate it from your existing command metadata if you can.

Framework notes are available for Cobra, Click, argparse, Clap, and Commander in [`docs/IMPLEMENTATION_GUIDES.md`](docs/IMPLEMENTATION_GUIDES.md).

## Is v0.1 stable enough to try?

Yes. Treat it as a small draft convention for experiments and early adopters. The core shape is intentionally tiny; names and registry details may still change before v1.0.
