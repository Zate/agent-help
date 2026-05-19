# mem — agent-help reference implementation

A minimal Go/Cobra CLI that demonstrates the [agent-help convention](../SPEC.md) in practice.

`mem` is a fake "project memory" tool — it stores facts, decisions, patterns, and tasks. The domain is simple enough to understand immediately, rich enough to show every agent-help scenario.

## Build

```bash
cd impl
go build -o mem .
./mem --help
```

The demo depends on `github.com/toon-format/toon-go`. At the time of this draft, that module has no tagged release visible to `go list -m -versions`, so `go.mod` uses a Go pseudo-version. Replace it with a tagged version when one is published.

## Commands

| Command | Subcommands | Demonstrates |
|---|---|---|
| `mem node` | `add`, `list` | AH2 group, creating records, TOON list output, pagination |
| `mem search` | `query`, `similar` | Required args, enum flags, no-results warning, not-found error |
| `mem project` | `status`, `set` | Single-object TOON, key-value output, missing flag errors |

## Try it: agent-help surfaces

```bash
# Discovery breadcrumb in normal --help
./mem --help

# AH1: index of all commands
./mem --agent-help

# AH2: group-level help
./mem node --agent-help
./mem search --agent-help
./mem project --agent-help

# AH2: command detail
./mem node add --agent-help
./mem node list --agent-help
./mem search query --agent-help
./mem search similar --agent-help
./mem project status --agent-help
./mem project set --agent-help
```

## Try it: --agent-out runtime output

```bash
# List all nodes (TOON body)
./mem node list --agent-out

# Paginated list (more=1 + next record)
./mem node list --limit 3 --agent-out

# Add a node
./mem node add "AHF ok/err envelope precedes TOON body" --type fact --tag spec:ahf --agent-out

# Search — with results
./mem search query "project" --agent-out

# Search — no results (warn + next)
./mem search query "zzznomatch" --agent-out

# Similarity search
./mem search similar n_101 --agent-out

# Project settings (single-object TOON)
./mem project status --agent-out

# Update a setting
./mem project set --key status --value active --agent-out
```

## Try it: error handling

```bash
# Missing required flag → err/hint/use
./mem node add "test node" --agent-out

# Invalid enum value → err/hint/use
./mem node add "test node" --type badval --agent-out

# Missing required positional arg → err/hint/use
./mem search query --agent-out

# Node not found → err/hint/next
./mem search similar n_999 --agent-out

# Missing both required flags → err/hint/use
./mem project set --agent-out

# Missing one required flag → err/hint/use
./mem project set --key status --agent-out
```

## What to observe

**AH1 index** — flat list of all commands with signatures. An agent can orient itself in one call.

**AH2 detail** — exact `use` line, typed `arg`/`flag` records, working `ex` examples. An agent can construct a correct invocation without guessing.

**TOON bodies** — compact tabular output. No JSON punctuation overhead. Field names declared once in the header, data rows follow.

**`ok`/`err` envelope** — always the first line. An agent knows immediately whether the command succeeded before parsing the body.

**`more?` pointer** — tells the agent how to get AH2 detail without being a shell command.

**`next` records** — exact follow-up commands embedded in the output. An agent never has to infer what to do next.

**`warn` records** — non-fatal signals (truncation, no results) that don't change success semantics.

## Pattern: --agent-help bypass for Args

Commands with required positional args use a custom `Args` validator that bypasses validation when `--agent-help` is set. This is the recommended Cobra pattern:

```go
Args: func(cmd *cobra.Command, args []string) error {
    if AgentHelp {
        return nil
    }
    return cobra.ExactArgs(1)(cmd, args)
},
```

This ensures `tool subcmd --agent-help` always works, even when the command normally requires arguments.
