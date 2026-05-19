# Framework Implementation Notes

These notes cover common CLI frameworks without repeating the full AHF shape. Use the [spec](https://zate.github.io/agent-help/spec.html) for the format, and use this file for framework-specific wiring.

## Minimum viable implementation

Do this first:

1. Keep human `--help` human-oriented.
2. Add this breadcrumb to normal help:
   `LLM agent? Use --agent-help for token-optimized usage.`
3. Add trailing `--agent-help`:
   - `tool --agent-help` emits the command index.
   - `tool subcmd --agent-help` emits command detail.
4. Generate that help from command metadata where practical.
5. Add `--agent-out` only to commands that return structured results.
6. Test the output with golden or snapshot tests.

That is enough. Avoid maturity labels, separate conformance levels, and hand-maintained docs that can drift.

## Tiny readiness checklist

Before calling an implementation done:

- [ ] `tool --help` points agents to `--agent-help`.
- [ ] `tool --agent-help` lists primary commands with `cmd` records.
- [ ] Every listed command supports `tool subcmd --agent-help`.
- [ ] Required args and flags in `--agent-help` match real validation.
- [ ] Examples either run successfully or are omitted.
- [ ] Structured result commands support `--agent-out`, or you intentionally skipped them.
- [ ] `--agent-out` starts with `ok` or `err` and includes exact `next` commands for pagination or useful follow-up.
- [ ] Agent-facing output avoids secrets, unnecessary personal data, colors, markdown tables, and prose paragraphs.

## Consuming output

Consumers do not need a formal parser to get value:

1. Read lines.
2. Inspect the first token: `ah1`, `cmd`, `ah2`, `use`, `arg`, `flag`, `ok`, `err`, `hint`, `next`, etc.
3. For `err`, use `hint`, `use`, and `next` records to recover.
4. For `ok`, read metadata like `count`, `more`, and `cursor`.
5. Treat the result body after `ok` as TOON when present.
6. Use `next` records as exact follow-up commands.

## Go / Cobra

Use a hidden persistent flag on the root command:

```go
var agentHelp bool

root.PersistentFlags().BoolVar(&agentHelp, "agent-help", false, "show token-optimized help for LLM agents")
_ = root.PersistentFlags().MarkHidden("agent-help")
```

Intercept before command execution and print AH1 for the root command or AH2 for the selected subcommand. Use a `PersistentPreRunE` hook:

```go
root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
    if agentHelp {
        if cmd == root {
            printAH1(root)
        } else {
            printAH2(cmd)
        }
        os.Exit(0)
    }
    return nil
}
```

Cobra-specific notes:

- `cmd.CommandPath()` gives the command path.
- `cmd.Use` gives the positional shape.
- `cmd.Short` gives the short purpose.
- Annotate scalar types explicitly; `pflag.Value.Type()` is not enough for good AHF.
- Commands with required positional args need an `Args` bypass when `--agent-help` is set.

```go
Args: func(cmd *cobra.Command, args []string) error {
    if AgentHelp {
        return nil
    }
    return cobra.ExactArgs(1)(cmd, args)
},
```

## Python / Click

Use a hidden eager option so Click can intercept before action execution:

```python
def agent_help_option(fn):
    return click.option(
        "--agent-help",
        is_flag=True,
        is_eager=True,
        expose_value=False,
        hidden=True,
        callback=emit_agent_help,
    )(fn)
```

Click-specific notes:

- Inspect the Click context's current command and parent command.
- Root groups emit AH1; concrete commands emit AH2.
- `command.params` exposes arguments and options.
- `click.Choice` maps naturally to `enum(a|b)`.
- Use custom attrs or a sidecar table for concise descriptions and examples.

## Python / argparse

Add a hidden flag to the root parser and to subparsers where needed:

```python
parser.add_argument("--agent-help", action="store_true", help=argparse.SUPPRESS)
```

Check the flag after `parse_known_args()`, before dispatch:

```python
args, remaining = parser.parse_known_args()
if args.agent_help:
    emit_ahf(parser, remaining)
    sys.exit(0)
```

argparse-specific notes:

- Subparser handling varies; `parse_known_args()` before dispatch can help identify the target subcommand.
- `argparse` internals are not ideal as the only metadata source — keep a small metadata table alongside parser construction for signatures, scalar types, and examples.
- Add `--agent-help` to each subparser as well as the root parser.

## Rust / Clap

Add a hidden global arg:

```rust
Arg::new("agent-help")
    .long("agent-help")
    .global(true)
    .hide(true)
    .action(ArgAction::SetTrue)
```

Check the flag in your dispatch logic before running the subcommand:

```rust
if matches.get_flag("agent-help") {
    emit_ahf(&cmd, &matches);
    std::process::exit(0);
}
```

Clap-specific notes:

- `Command::get_name()` gives command segments for building the command path.
- `Command::get_about()` gives short purpose text for AH1/AH2 headers.
- `Command::get_subcommands()` can walk the tree to build AH1.
- Possible values map naturally to `enum(a|b)` in AHF.
- Value parsers do not always map cleanly to AHF scalar types — explicit metadata is clearer.

## Node / Commander

Add a global option and detect it before normal actions:

```js
program.option("--agent-help", "show token-optimized help for LLM agents");

program.hook("preAction", (thisCommand) => {
  if (thisCommand.opts().agentHelp) {
    emitAHF(thisCommand);
    process.exit(0);
  }
});
```

Commander-specific notes:

- `program.addHelpText("after", "LLM agent? Use --agent-help for token-optimized usage.")` adds the breadcrumb.
- `program.commands` can be walked recursively to build AH1.
- `command.name()` gives the command segment; `command.description()` gives the purpose.
- Add explicit metadata for scalar types, defaults, and examples — Commander's option definitions don't carry AHF-ready type info.

## Reference Implementation

The Go/Cobra `mem` demo in [`impl/`](../impl/) is a complete, working implementation targeting `agent-help-full` conformance. It demonstrates:

- AH1 and AH2 help surfaces derived from Cobra command metadata
- AHF `ok`/`err`/`warn` protocol records
- TOON list and object bodies for `--agent-out`
- Pagination with continuation `next` records and cursor tokens
- Golden-output tests for key surfaces

Crib from it freely — it's Apache-2.0 licensed and designed to be copied.
