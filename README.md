# agent-help

**MCP-lite for CLIs: two small flags that let AI agents discover, call, and parse your tool without a server, plugin, or skill file.**

[![License: Apache-2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Docs License: CC-BY-4.0](https://img.shields.io/badge/docs-CC--BY--4.0-lightgrey.svg)](LICENSE-DOCS)
[![Draft: v0.1](https://img.shields.io/badge/Draft-v0.1-green.svg)](SPEC.md)

## The whole idea

```text
agent-help is a lightweight, token-efficient alternative to MCP servers and SKILL.md files for teaching AI agents how to talk to CLIs.

Instead of spinning up infrastructure or maintaining separate docs, you add two flags:

  --agent-help  machine-readable command discovery
  --agent-out   dense structured runtime results

Then agents can discover, invoke, and parse your tool without hallucinating.

It's MCP-lite: no daemons, no JSON bloat, just dense text that fits in context and replaces both the learning curve of skills and the engineering overhead of full MCP integration.
```

That is the product. Everything else in this repository exists only to help people and coding agents **implement those two flags correctly**.

## What agents need at runtime

An agent using an agent-help-enabled CLI should not need this repository, the spec, `llms-full.txt`, a skill, or any prior instructions.

It should be able to do this:

```text
$ tool --help
...
LLM agent? Use --agent-help for token-optimized usage.

$ tool --agent-help
ah1 tool :: example data CLI
cmd search <query> [--limit int] :: search records
cmd show <id> :: display one record
more? tool <cmd> --agent-help

$ tool search --agent-help
ah2 tool search
use tool search <query> [--limit int]
arg query:str req :: search text
flag --limit:int opt default=20 :: max rows
ex tool search postgres --limit 5

$ tool search postgres --agent-out
ok records count=2 more=0
records[#2]{id,type,text}:
  r_102,fact,postgres 15 required
  r_088,decision,migrate read model
next inspect "tool show r_102 --agent-out"
```

If an agent has to read a tutorial before it can use the CLI, agent-help has failed.

## What implementers add

1. Keep normal human help:

   ```text
   tool --help
   ```

2. Add a short breadcrumb to that help:

   ```text
   LLM agent? Use --agent-help for token-optimized usage.
   ```

3. Add `--agent-help` as a trailing global flag:

   ```text
   tool --agent-help          # command index
   tool subcmd --agent-help   # command detail
   ```

4. Add `--agent-out` to commands that return structured results:

   ```text
   tool subcmd args --agent-out
   ```

5. Generate the output from your CLI's existing command metadata where possible so it cannot drift from reality.

## Tiny shape reference

`--agent-help` is dense invocation help:

```text
ah1 tool :: what this CLI does
cmd subcmd <arg> [--flag type] :: what this command does
more? tool <cmd> --agent-help

ah2 tool subcmd
use tool subcmd <arg> [--flag type]
arg arg:str req :: what the arg is
flag --flag:int opt default=10 :: what the flag does
ex tool subcmd thing --flag 5
```

`--agent-out` is dense runtime output:

```text
ok kind count=2 more=0
kind[#2]{id,status,summary}:
  a1,open,first thing
  a2,done,second thing
next inspect "tool show a1 --agent-out"
```

Errors should help the agent recover:

```text
err missing_flag flag=--type
hint --type enum(fact|decision|task)
use tool add <text> --type TYPE
```

## Repository map

Human-first docs:

- [`site/`](site/) — static website source
- [`FAQ.md`](FAQ.md) — short explanation and tradeoffs
- [`docs/IMPLEMENTATION_GUIDES.md`](docs/IMPLEMENTATION_GUIDES.md) — framework wiring notes
- [`examples/`](examples/) — small annotated examples

Implementation-only references:

- [`SPEC.md`](SPEC.md) — technical wire-format spec for implementers and validators
- [`llms-full.txt`](llms-full.txt) — paste this into a coding agent when asking it to add agent-help to a CLI
- [`.agents/skills/ahf/SKILL.md`](.agents/skills/ahf/SKILL.md) — optional build-time skill for agent-assisted implementation
- [`llms.txt`](llms.txt) — short orientation that points agents to implementation resources
- [`agent-help.ahf`](agent-help.ahf) / [`agent-help-full.ahf`](agent-help-full.ahf) — dogfood examples for these docs
- [`spec/ahf-v0.1.json`](spec/ahf-v0.1.json) — machine-readable registry for tooling

Code and tests:

- [`impl/`](impl/) — Go/Cobra reference implementation
- [`tests/`](tests/) — valid and invalid AHF fixtures
- [`scripts/`](scripts/) — validation and site checks

## Try it

```bash
make demo
make test
```

## Use a coding agent to add agent-help to a CLI

```text
Read https://raw.githubusercontent.com/Zate/agent-help/main/llms-full.txt and implement agent-help in this CLI. Preserve normal --help and --json behavior. Add --agent-help for discovery and --agent-out for structured runtime results where appropriate.
```

## Status

Draft v0.1. The core idea is intentionally small and ready for trial implementations. Names and registry details may change before v1.0.

## License

Code and examples: [Apache 2.0](LICENSE)  
Documentation and specification: [CC-BY-4.0](LICENSE-DOCS)  
Notices: [NOTICE](NOTICE)
