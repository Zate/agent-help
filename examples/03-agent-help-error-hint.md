# AHF Example 03: agent-help error hint

Source command:

```bash
mem node add "postgres 15 required" --agent-out
```

Raw AHF output:

```text
err missing_flag flag=--type
hint --type enum(decision|fact|pattern|task|observation)
use mem node add <text> --type TYPE [--tag K:V...]
```
