# AHF Example 02: agent-help command detail

Source command:

```bash
mem node add --agent-help
```

Raw AHF output:

```text
ah2 mem node add
use mem node add <text> --type TYPE [--tag K:V...]
arg text:str req :: node text content
flag --type:enum(decision|fact|pattern|task|observation) req :: node type
flag --tag:kv repeat :: metadata key:value pairs
ex mem node add "postgres 15 required" --type fact --tag project:mem
ex mem node add "use TOON for --agent-out bodies" --type decision --tag project:mem --tag spec:ahf
```
