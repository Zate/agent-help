# Framework Implementation Notes

These notes cover common CLI frameworks without repeating the full AHF shape. Use the spec for the format, and use this file for framework-specific wiring.

## Shared Pattern

- Keep human `--help` human-oriented.
- Add a short breadcrumb or equivalent pointer to `--agent-help`.
- Treat documented `--agent-help` placement as trailing: `tool subcmd --agent-help`.
- Generate AH1/AH2 from command metadata where practical.
- Add `--agent-out` only to commands that return structured runtime results.
- Protect the surface with golden or snapshot tests.

## Go / Cobra

Use a hidden persistent flag on the root command:

```go
var agentHelp bool

root.PersistentFlags().BoolVar(&agentHelp, "agent-help", false, "show token-optimized help for LLM agents")
_ = root.PersistentFlags().MarkHidden("agent-help")
```

Intercept before command execution and print AH1 for the root command or AH2 for the selected command.

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

argparse-specific notes:

- Subparser handling varies; `parse_known_args()` before dispatch can help.
- `argparse` internals are not ideal as the only metadata source.
- Keep a small metadata table next to parser construction for signatures, scalar types, and examples.

## Rust / Clap

Add a hidden global arg:

```rust
Arg::new("agent-help")
    .long("agent-help")
    .global(true)
    .hide(true)
    .action(ArgAction::SetTrue)
```

Clap-specific notes:

- `Command::get_name()` gives command segments.
- `Command::get_about()` gives short purpose text.
- `Command::get_subcommands()` can walk command paths.
- Possible values map naturally to `enum(a|b)`.
- Value parsers do not always map cleanly to AHF scalar types; explicit metadata is clearer.

## Node / Commander

Add a global option and detect it before normal actions:

```js
program.option("--agent-help", "show token-optimized help for LLM agents");
```

Commander-specific notes:

- `program.addHelpText("after", "...")` can add the breadcrumb.
- `program.commands` can be walked recursively.
- `command.name()` gives the command segment.
- `command.description()` gives the purpose.
- Add explicit metadata for scalar types, defaults, and examples.
