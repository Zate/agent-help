package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage project settings",
	Long:  "View and update project-level settings.\n\nLLM agent? Use --agent-help for token-optimized usage.",
	Run: func(cmd *cobra.Command, args []string) {
		if AgentHelp {
			projectGroupAgentHelp()
			return
		}
		cmd.Help()
	},
}

func projectGroupAgentHelp() {
	ah2("mem", "project")
	use("mem project <subcommand>")
	cmdEntry("project status", "show current project settings")
	cmdEntry("project set --key KEY --value VALUE", "set a project setting")
	morePtr("mem", "project <subcommand>")
}

// --- project status ---

var projectStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current project settings",
	Long:  "Display all current project settings as a single object result.\n\nLLM agent? Use --agent-help for token-optimized usage.",
	Run: func(cmd *cobra.Command, args []string) {
		if AgentHelp {
			projectStatusAgentHelp()
			return
		}
		runProjectStatus()
	},
}

func projectStatusAgentHelp() {
	ah2("mem", "project status")
	use("mem project status")
	ex("mem project status --agent-out")
}

func runProjectStatus() {
	if AgentOut {
		okLine("project")
		emitTOON(ProjectList{Project: MockProject})
		nextLine("nodes", "mem node list --agent-out")
	} else {
		fmt.Println("Project settings:")
		for _, s := range MockProject {
			fmt.Printf("  %-10s  %s\n", s.Key, s.Value)
		}
	}
}

// --- project set ---

var (
	projectSetKey   string
	projectSetValue string
)

var projectSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a project setting",
	Long:  "Update a project setting by key.\n\nLLM agent? Use --agent-help for token-optimized usage.",
	Run: func(cmd *cobra.Command, args []string) {
		if AgentHelp {
			projectSetAgentHelp()
			return
		}
		runProjectSet()
	},
}

func projectSetAgentHelp() {
	ah2("mem", "project set")
	use("mem project set --key KEY --value VALUE")
	flagEntry("key", "str", "req", "setting key to update")
	flagEntry("value", "str", "req", "new value for the setting")
	ex(`mem project set --key status --value active --agent-out`)
	ex(`mem project set --key owner --value alice --agent-out`)
}

func runProjectSet() {
	// Validate required flags.
	if projectSetKey == "" && projectSetValue == "" {
		if AgentOut {
			errLine("missing_flags", "flags=--key,--value")
			hint("both --key and --value are required")
			useLine("mem project set --key KEY --value VALUE")
		} else {
			fmt.Println("Error: --key and --value are required")
			fmt.Println("  usage: mem project set --key KEY --value VALUE")
		}
		return
	}
	if projectSetKey == "" {
		if AgentOut {
			errLine("missing_flag", "flag=--key")
			hint("--key str req :: the setting key to update")
			useLine("mem project set --key KEY --value VALUE")
		} else {
			fmt.Println("Error: --key is required")
		}
		return
	}
	if projectSetValue == "" {
		if AgentOut {
			errLine("missing_flag", "flag=--value")
			hint("--value str req :: the new value to set")
			useLine("mem project set --key KEY --value VALUE")
		} else {
			fmt.Println("Error: --value is required")
		}
		return
	}

	// Find and update the setting (in-memory only — this is a mock).
	found := false
	for i := range MockProject {
		if MockProject[i].Key == projectSetKey {
			MockProject[i].Value = projectSetValue
			found = true
			break
		}
	}
	if !found {
		// Add new setting.
		MockProject = append(MockProject, ProjectRow{projectSetKey, projectSetValue})
	}

	if AgentOut {
		okLine("setting")
		kvLine("key", projectSetKey)
		kvLine("value", projectSetValue)
		kvLine("status", "updated")
		nextLine("verify", "mem project status --agent-out")
	} else {
		if found {
			fmt.Printf("Updated %s = %s\n", projectSetKey, projectSetValue)
		} else {
			fmt.Printf("Created %s = %s\n", projectSetKey, projectSetValue)
		}
		fmt.Println("  use: mem project status to verify")
	}
}

func init() {
	projectSetCmd.Flags().StringVar(&projectSetKey, "key", "", "setting key")
	projectSetCmd.Flags().StringVar(&projectSetValue, "value", "", "new value")

	projectCmd.AddCommand(projectStatusCmd)
	projectCmd.AddCommand(projectSetCmd)
	rootCmd.AddCommand(projectCmd)
}
