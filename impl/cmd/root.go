// Package cmd implements the mem CLI — the agent-help reference implementation.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// AgentHelp and AgentOut are the global flags for agent-native surfaces.
var AgentHelp bool
var AgentOut bool

// rootCmd is the top-level command.
var rootCmd = &cobra.Command{
	Use:   "mem",
	Short: "Project memory CLI — store and query facts, decisions, and tasks",
	Long: `mem — project memory CLI

Store and retrieve facts, decisions, patterns, and tasks about your project.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// --agent-help is handled per-command before execution.
		// This hook ensures the breadcrumb always appears in --help.
	},
	Run: func(cmd *cobra.Command, args []string) {
		if AgentHelp {
			rootAgentHelp()
			return
		}
		cmd.Help()
	},
}

// rootAgentHelp emits AH1 for the root command.
func rootAgentHelp() {
	ah1("mem", "project memory — store and query facts, decisions, and tasks")
	cmdEntry("node add <text> --type TYPE [--tag K:V...]", "store a new memory node")
	cmdEntry("node list [--type TYPE] [--limit int] [--cursor id]", "list memory nodes")
	cmdEntry("search query <text> [--type TYPE] [--limit int] [--cursor id]", "search nodes by text")
	cmdEntry("search similar <id> [--limit int] [--cursor id]", "find nodes similar to a given node")
	cmdEntry("project status", "show current project settings")
	cmdEntry("project set --key KEY --value VALUE", "set a project setting")
	morePtr("mem", "<cmd>")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global hidden flags for agent-native surfaces.
	rootCmd.PersistentFlags().BoolVar(&AgentHelp, "agent-help", false, "")
	rootCmd.PersistentFlags().BoolVar(&AgentOut, "agent-out", false, "")
	_ = rootCmd.PersistentFlags().MarkHidden("agent-help")
	_ = rootCmd.PersistentFlags().MarkHidden("agent-out")

	// Append the discovery breadcrumb to the long description so it appears in --help.
	rootCmd.Long = strings.TrimSpace(rootCmd.Long) + "\n\nLLM agent? Use --agent-help for token-optimized usage."
}
