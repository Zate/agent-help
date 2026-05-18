package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search memory nodes",
	Long:  "Search memory nodes by text or similarity.\n\nLLM agent? Use --agent-help for token-optimized usage.",
	Run: func(cmd *cobra.Command, args []string) {
		if AgentHelp {
			searchGroupAgentHelp()
			return
		}
		cmd.Help()
	},
}

func searchGroupAgentHelp() {
	ah2("mem", "search")
	use("mem search <subcommand>")
	cmdEntry("search query <text> [--type TYPE] [--limit int] [--cursor id]", "search nodes by text")
	cmdEntry("search similar <id> [--limit int] [--cursor id]", "find nodes similar to a given node")
	morePtr("mem", "search <subcommand>")
}

// --- search query ---

var (
	searchQueryType  string
	searchQueryLimit int
	searchQueryCursor string
)

var searchQueryCmd = &cobra.Command{
	Use:   "query <text>",
	Short: "Search nodes by text",
	Long:  "Full-text search across all memory nodes.\n\nLLM agent? Use --agent-help for token-optimized usage.",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if AgentHelp {
			searchQueryAgentHelp()
			return
		}
		if len(args) == 0 {
			if AgentOut {
				errLine("missing_arg", "arg=text")
				hint("provide a search string as the first argument")
				useLine("mem search query <text> [--type TYPE] [--limit int]")
			} else {
				fmt.Println("Error: search text is required")
				fmt.Println("  usage: mem search query <text> [--type TYPE] [--limit int]")
			}
			return
		}
		runSearchQuery(args[0])
	},
}

func searchQueryAgentHelp() {
	ah2("mem", "search query")
	use("mem search query <text> [--type TYPE] [--limit int] [--cursor id]")
	argEntry("text", "str", "req", "search string")
	flagEntry("type", fmt.Sprintf("enum(%s)", strings.Join(ValidNodeTypes, "|")), "opt", "filter by node type")
	flagEntryDefault("limit", "int", "opt", "10", "max results")
	flagEntry("cursor", "id", "opt", "resume after the given result ID")
	ex(`mem search query "postgres" --agent-out`)
	ex(`mem search query "decision" --type decision --agent-out`)
}

func runSearchQuery(query string) {
	// Validate type enum if provided.
	if searchQueryType != "" {
		validType := false
		for _, t := range ValidNodeTypes {
			if t == searchQueryType {
				validType = true
				break
			}
		}
		if !validType {
			if AgentOut {
				errLine("invalid_enum", "flag=--type", fmt.Sprintf("got=%s", searchQueryType))
				hint(fmt.Sprintf("--type enum(%s)", strings.Join(ValidNodeTypes, "|")))
				useLine("mem search query <text> [--type TYPE] [--limit int] [--cursor id]")
			} else {
				fmt.Printf("Error: invalid --type %q\n", searchQueryType)
				fmt.Printf("  valid types: %s\n", strings.Join(ValidNodeTypes, ", "))
			}
			return
		}
	}

	// Mock search: filter nodes where text contains the query string (case-insensitive).
	q := strings.ToLower(query)
	results := []Node{}
	for _, n := range MockNodes {
		if strings.Contains(strings.ToLower(n.Text), q) ||
			strings.Contains(strings.ToLower(n.Tags), q) {
			if searchQueryType == "" || n.Type == searchQueryType {
				results = append(results, n)
			}
		}
	}

	total := len(results)
	if searchQueryCursor != "" {
		start := -1
		for i := range results {
			if results[i].ID == searchQueryCursor {
				start = i
				break
			}
		}
		if start == -1 {
			if AgentOut {
				errLine("not_found", fmt.Sprintf("cursor=%s", searchQueryCursor))
				hint("use a cursor value from a previous search query response")
				useLine("mem search query <text> [--type TYPE] [--limit int] [--cursor id]")
			} else {
				fmt.Printf("Error: cursor %q not found\n", searchQueryCursor)
			}
			return
		}
		results = results[start:]
	}

	limit := searchQueryLimit
	if limit <= 0 {
		limit = 10
	}
	more := 0
	shown := len(results)
	cursor := ""
	if shown > limit {
		more = 1
		cursor = fmt.Sprintf("c_%s", results[limit].ID)
		results = results[:limit]
		shown = limit
	}

	if AgentOut {
		meta := []string{
			fmt.Sprintf("query=%q", query),
			fmt.Sprintf("count=%d", total),
			fmt.Sprintf("more=%d", more),
		}
		if cursor != "" {
			meta = append(meta, fmt.Sprintf("cursor=%s", cursor))
		}
		okLine("nodes", meta...)
		if shown == 0 {
			warnLine("no_results", fmt.Sprintf("query=%q", query))
			nextLine("broaden", "mem node list --agent-out")
			return
		}
		emitTOON(toViewList(results))
		if more == 1 {
			nextLine("", fmt.Sprintf("mem search query %q --limit %d --cursor %s --agent-out", query, limit, cursor))
		}
		if shown > 0 {
			nextLine("inspect", fmt.Sprintf("mem search similar %s --agent-out", results[0].ID))
		}
	} else {
		if shown == 0 {
			fmt.Printf("No results for %q\n", query)
			return
		}
		fmt.Printf("Results for %q (%d):\n", query, shown)
		for _, n := range results {
			fmt.Printf("  %s  %-12s  %s\n", n.ID, n.Type, n.Text)
		}
	}
}

// --- search similar ---

var searchSimilarLimit int
var searchSimilarCursor string

var searchSimilarCmd = &cobra.Command{
	Use:   "similar <id>",
	Short: "Find nodes similar to a given node",
	Long:  "Find semantically similar nodes based on a reference node ID.\n\nLLM agent? Use --agent-help for token-optimized usage.",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if AgentHelp {
			searchSimilarAgentHelp()
			return
		}
		if len(args) == 0 {
			if AgentOut {
				errLine("missing_arg", "arg=id")
				hint("provide a node ID as the first argument")
				useLine("mem search similar <id> [--limit int]")
			} else {
				fmt.Println("Error: node ID is required")
				fmt.Println("  usage: mem search similar <id>")
			}
			return
		}
		runSearchSimilar(args[0])
	},
}

func searchSimilarAgentHelp() {
	ah2("mem", "search similar")
	use("mem search similar <id> [--limit int] [--cursor id]")
	argEntry("id", "id", "req", "reference node ID")
	flagEntryDefault("limit", "int", "opt", "5", "max similar results")
	flagEntry("cursor", "id", "opt", "resume after the given result ID")
	ex("mem search similar n_102 --agent-out")
	ex("mem search similar n_105 --limit 3 --agent-out")
}

func runSearchSimilar(id string) {
	// Find the reference node.
	var ref *Node
	for i := range MockNodes {
		if MockNodes[i].ID == id {
			ref = &MockNodes[i]
			break
		}
	}
	if ref == nil {
		if AgentOut {
			errLine("not_found", fmt.Sprintf("id=%s", id))
			hint("check the node ID with: mem node list --agent-out")
			nextLine("list", "mem node list --agent-out")
		} else {
			fmt.Printf("Error: node %q not found\n", id)
			fmt.Println("  use: mem node list to see available nodes")
		}
		return
	}

	// Mock similar: return same-type nodes (excluding the reference node itself).
	limit := searchSimilarLimit
	if limit <= 0 {
		limit = 5
	}
	results := []Node{}
	for _, n := range MockNodes {
		if n.ID != id && n.Type == ref.Type {
			results = append(results, n)
		}
	}

	total := len(results)
	if searchSimilarCursor != "" {
		start := -1
		for i := range results {
			if results[i].ID == searchSimilarCursor {
				start = i
				break
			}
		}
		if start == -1 {
			if AgentOut {
				errLine("not_found", fmt.Sprintf("cursor=%s", searchSimilarCursor))
				hint("use a cursor value from a previous search similar response")
				useLine("mem search similar <id> [--limit int] [--cursor id]")
			} else {
				fmt.Printf("Error: cursor %q not found\n", searchSimilarCursor)
			}
			return
		}
		results = results[start:]
	}

	more := 0
	shown := len(results)
	cursor := ""
	if shown > limit {
		more = 1
		cursor = fmt.Sprintf("c_%s", results[limit].ID)
		results = results[:limit]
		shown = limit
	}

	if AgentOut {
		meta := []string{fmt.Sprintf("ref=%s", id), fmt.Sprintf("count=%d", total), fmt.Sprintf("shown=%d", shown), fmt.Sprintf("more=%d", more)}
		if cursor != "" {
			meta = append(meta, fmt.Sprintf("cursor=%s", cursor))
		}
		okLine("nodes", meta...)
		if shown == 0 {
			warnLine("no_similar", fmt.Sprintf("type=%s", ref.Type))
			nextLine("broaden", "mem node list --agent-out")
			return
		}
		emitTOON(toViewList(results))
		if more == 1 {
			nextLine("", fmt.Sprintf("mem search similar %s --limit %d --cursor %s --agent-out", id, limit, cursor))
		}
		nextLine("inspect", fmt.Sprintf("mem search query %q --type %s --agent-out", ref.Type, ref.Type))
	} else {
		if shown == 0 {
			fmt.Printf("No similar nodes found for %s (type: %s)\n", id, ref.Type)
			return
		}
		fmt.Printf("Similar to %s (%s):\n", id, ref.Type)
		for _, n := range results {
			fmt.Printf("  %s  %-12s  %s\n", n.ID, n.Type, n.Text)
		}
	}
}

func init() {
	searchQueryCmd.Flags().StringVar(&searchQueryType, "type", "", "filter by node type")
	searchQueryCmd.Flags().IntVar(&searchQueryLimit, "limit", 10, "max results")
	searchQueryCmd.Flags().StringVar(&searchQueryCursor, "cursor", "", "resume after result ID")

	searchSimilarCmd.Flags().IntVar(&searchSimilarLimit, "limit", 5, "max similar results")
	searchSimilarCmd.Flags().StringVar(&searchSimilarCursor, "cursor", "", "resume after result ID")

	searchCmd.AddCommand(searchQueryCmd)
	searchCmd.AddCommand(searchSimilarCmd)
	rootCmd.AddCommand(searchCmd)
}
