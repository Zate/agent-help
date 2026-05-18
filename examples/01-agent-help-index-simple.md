# AHF Example 01: agent-help index

Source command:

```bash
mem --agent-help
```

Raw AHF output:

```text
ah1 mem :: project memory — store and query facts, decisions, and tasks
cmd node add <text> --type TYPE [--tag K:V...] :: store a new memory node
cmd node list [--type TYPE] [--limit int] [--cursor id] :: list memory nodes
cmd search query <text> [--type TYPE] [--limit int] [--cursor id] :: search nodes by text
cmd search similar <id> [--limit int] [--cursor id] :: find nodes similar to a given node
cmd project status :: show current project settings
cmd project set --key KEY --value VALUE :: set a project setting
more? mem <cmd> --agent-help
```
