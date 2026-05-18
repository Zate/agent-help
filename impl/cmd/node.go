package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Manage memory nodes",
	Long:  "Manage memory nodes — add and list facts, decisions, patterns, and tasks.\n\nLLM agent? Use --agent-help for token-optimized usage.",
	Run: func(cmd *cobra.Command, args []string) {
		if AgentHelp {
			nodeGroupAgentHelp()
			return
		}
		cmd.Help()
	},
}

func nodeGroupAgentHelp() {
	ah2("mem", "node")
	use("mem node <subcommand>")
	cmdEntry("node add <text> --type TYPE [--tag K:V...]", "store a new memory node")
	cmdEntry("node list [--type TYPE] [--limit int] [--cursor id]", "list memory nodes")
	morePtr("mem", "node <subcommand>")
}

// --- node add ---

var (
	nodeAddType string
	nodeAddTags []string
)

var nodeAddCmd = &cobra.Command{
	Use:   "add <text>",
	Short: "Store a new memory node",
	Long:  "Store a new memory node with a type and optional tags.\n\nLLM agent? Use --agent-help for token-optimized usage.",
	// Use Args validator that allows --agent-help with no positional args.
	Args: func(cmd *cobra.Command, args []string) error {
		if AgentHelp {
			return nil
		}
		return cobra.ExactArgs(1)(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if AgentHelp {
			nodeAddAgentHelp()
			return
		}
		runNodeAdd(args[0])
	},
}

func nodeAddAgentHelp() {
	ah2("mem", "node add")
	use("mem node add <text> --type TYPE [--tag K:V...]")
	argEntry("text", "str", "req", "node text content")
	flagEntry("type", fmt.Sprintf("enum(%s)", strings.Join(ValidNodeTypes, "|")), "req", "node type")
	flagEntry("tag", "kv", "repeat", "metadata key:value pairs")
	ex(`mem node add "postgres 15 required" --type fact --tag project:mem`)
	ex(`mem node add "use TOON for --agent-out bodies" --type decision --tag project:mem --tag spec:ahf`)
}

func runNodeAdd(text string) {
	if nodeAddType == "" {
		if AgentOut {
			errLine("missing_flag", "flag=--type")
			hint(fmt.Sprintf("--type enum(%s)", strings.Join(ValidNodeTypes, "|")))
			useLine("mem node add <text> --type TYPE [--tag K:V...]")
		} else {
			fmt.Println("Error: --type is required")
			fmt.Printf("  valid types: %s\n", strings.Join(ValidNodeTypes, ", "))
		}
		return
	}

	validType := false
	for _, t := range ValidNodeTypes {
		if t == nodeAddType {
			validType = true
			break
		}
	}
	if !validType {
		if AgentOut {
			errLine("invalid_enum", "flag=--type", fmt.Sprintf("got=%s", nodeAddType))
			hint(fmt.Sprintf("--type enum(%s)", strings.Join(ValidNodeTypes, "|")))
			useLine("mem node add <text> --type TYPE")
		} else {
			fmt.Printf("Error: invalid --type %q\n", nodeAddType)
			fmt.Printf("  valid types: %s\n", strings.Join(ValidNodeTypes, ", "))
		}
		return
	}

	// Generate a fake ID and emit result.
	id := fmt.Sprintf("n_%d", 100+len(MockNodes)+1)
	tags := "_"
	if len(nodeAddTags) > 0 {
		tags = strings.Join(nodeAddTags, "|")
	}
	created := time.Now().Format("2006-01-02")

	if AgentOut {
		okLine("node")
		kvLine("id", id)
		kvLine("type", nodeAddType)
		kvLine("tags", tags)
		kvLine("created", created)
		kvLine("text", text)
		nextLine("list", "mem node list --agent-out")
	} else {
		fmt.Printf("Created node %s\n", id)
		fmt.Printf("  type:    %s\n", nodeAddType)
		fmt.Printf("  tags:    %s\n", tags)
		fmt.Printf("  created: %s\n", created)
		fmt.Printf("  text:    %s\n", text)
	}
}

// --- node list ---

var (
	nodeListType  string
	nodeListLimit int
	nodeListCursor string
)

var nodeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List memory nodes",
	Long:  "List memory nodes, optionally filtered by type.\n\nLLM agent? Use --agent-help for token-optimized usage.",
	Run: func(cmd *cobra.Command, args []string) {
		if AgentHelp {
			nodeListAgentHelp()
			return
		}
		runNodeList()
	},
}

func nodeListAgentHelp() {
	ah2("mem", "node list")
	use("mem node list [--type TYPE] [--limit int] [--cursor id]")
	flagEntry("type", fmt.Sprintf("enum(%s)", strings.Join(ValidNodeTypes, "|")), "opt", "filter by node type")
	flagEntryDefault("limit", "int", "opt", "10", "max results to return")
	flagEntry("cursor", "id", "opt", "resume after the given node ID")
	ex("mem node list --agent-out")
	ex("mem node list --type decision --agent-out")
	ex("mem node list --limit 3 --agent-out")
}

func runNodeList() {
	// Filter by type if requested.
	nodes := MockNodes
	if nodeListType != "" {
		// Validate enum.
		validType := false
		for _, t := range ValidNodeTypes {
			if t == nodeListType {
				validType = true
				break
			}
		}
		if !validType {
			if AgentOut {
				errLine("invalid_enum", "flag=--type", fmt.Sprintf("got=%s", nodeListType))
				hint(fmt.Sprintf("--type enum(%s)", strings.Join(ValidNodeTypes, "|")))
				useLine("mem node list [--type TYPE] [--limit int] [--cursor id]")
			} else {
				fmt.Printf("Error: invalid --type %q\n", nodeListType)
				fmt.Printf("  valid types: %s\n", strings.Join(ValidNodeTypes, ", "))
			}
			return
		}
		filtered := []Node{}
		for _, n := range nodes {
			if n.Type == nodeListType {
				filtered = append(filtered, n)
			}
		}
		nodes = filtered
	}

	total := len(nodes)
	if nodeListCursor != "" {
		start := -1
		for i := range nodes {
			if nodes[i].ID == nodeListCursor {
				start = i
				break
			}
		}
		if start == -1 {
			if AgentOut {
				errLine("not_found", fmt.Sprintf("cursor=%s", nodeListCursor))
				hint("use a cursor value from a previous node list response")
				useLine("mem node list [--type TYPE] [--limit int] [--cursor id]")
			} else {
				fmt.Printf("Error: cursor %q not found\n", nodeListCursor)
			}
			return
		}
		nodes = nodes[start:]
	}

	// Apply limit and pagination.
	limit := nodeListLimit
	if limit <= 0 {
		limit = 10
	}
	more := 0
	cursor := ""
	shown := len(nodes)
	if shown > limit {
		more = 1
		cursor = fmt.Sprintf("c_%s", nodes[limit].ID)
		nodes = nodes[:limit]
		shown = limit
	}

	if AgentOut {
		meta := []string{
			fmt.Sprintf("count=%d", total),
			fmt.Sprintf("shown=%d", shown),
			fmt.Sprintf("more=%d", more),
		}
		if cursor != "" {
			meta = append(meta, fmt.Sprintf("cursor=%s", cursor))
		}
		okLine("nodes", meta...)
		if more == 1 {
			warnLine("truncated", fmt.Sprintf("shown=%d", shown), fmt.Sprintf("total=%d", total))
		}
		emitTOON(toViewList(nodes))
		if more == 1 {
			nextLine("", fmt.Sprintf("mem node list --limit %d --cursor %s --agent-out", limit, cursor))
		}
		nextLine("inspect", "mem search query <text> --agent-out")
	} else {
		fmt.Printf("Nodes (%d of %d):\n", shown, total)
		for _, n := range nodes {
			fmt.Printf("  %s  %-12s  %-30s  %s\n", n.ID, n.Type, n.Tags, n.Text)
		}
		if more == 1 {
			fmt.Printf("\n  (showing %d of %d — use --limit to see more)\n", shown, total)
		}
	}
}

func init() {
	// node add flags
	nodeAddCmd.Flags().StringVar(&nodeAddType, "type", "", "node type")
	nodeAddCmd.Flags().StringArrayVar(&nodeAddTags, "tag", []string{}, "metadata key:value")

	// node list flags
	nodeListCmd.Flags().StringVar(&nodeListType, "type", "", "filter by node type")
	nodeListCmd.Flags().IntVar(&nodeListLimit, "limit", 10, "max results")
	nodeListCmd.Flags().StringVar(&nodeListCursor, "cursor", "", "resume after node ID")

	nodeCmd.AddCommand(nodeAddCmd)
	nodeCmd.AddCommand(nodeListCmd)
	rootCmd.AddCommand(nodeCmd)
}
