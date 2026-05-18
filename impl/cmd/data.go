// Package cmd provides mock data for the mem reference implementation.
// In a real CLI these would come from storage.
package cmd

// Node represents a memory node — a fact, decision, pattern, or task.
// TOON struct tags drive emitTOON serialization for --agent-out results.
type Node struct {
	ID      string `toon:"id"`
	Type    string `toon:"type"`
	Tags    string `toon:"tags"`
	Text    string `toon:"text"`
	Created string `toon:"created"`
}

// NodeList is the TOON envelope for full list results (includes created).
type NodeList struct {
	Nodes []Node `toon:"nodes"`
}

// NodeView is a compact node for list/search results (omits created).
type NodeView struct {
	ID   string `toon:"id"`
	Type string `toon:"type"`
	Tags string `toon:"tags"`
	Text string `toon:"text"`
}

// NodeViewList is the TOON envelope for compact list/search output.
type NodeViewList struct {
	Nodes []NodeView `toon:"nodes"`
}

// toViewList converts a slice of Nodes to NodeViewList for compact output.
func toViewList(nodes []Node) NodeViewList {
	views := make([]NodeView, len(nodes))
	for i, n := range nodes {
		views[i] = NodeView{n.ID, n.Type, n.Tags, n.Text}
	}
	return NodeViewList{Nodes: views}
}

// ProjectRow is a single project setting key/value pair for TOON output.
type ProjectRow struct {
	Key   string `toon:"key"`
	Value string `toon:"value"`
}

// ProjectList is the TOON envelope for project status output.
type ProjectList struct {
	Project []ProjectRow `toon:"project"`
}

// MockNodes is the full set of nodes returned by node list / search.
var MockNodes = []Node{
	{"n_101", "decision", "project:mem|arch", "Use Cobra for the reference CLI implementation", "2026-05-11"},
	{"n_102", "fact", "project:mem|db", "Mock data lives in data.go; no real storage in the reference impl", "2026-05-11"},
	{"n_103", "pattern", "db|ops", "Use connection pool max 20 in all environments", "2026-05-10"},
	{"n_104", "fact", "project:mem|api", "AHF ok/err envelope always precedes the TOON body", "2026-05-11"},
	{"n_105", "decision", "project:mem|spec", "TOON is the encoding for --agent-out result bodies", "2026-05-11"},
	{"n_106", "task", "project:mem|todo", "Write reference Go implementation for agent-help", "2026-05-11"},
	{"n_107", "pattern", "agent|cli", "Always include next records for paginated results", "2026-05-10"},
	{"n_108", "fact", "project:mem|spec", "more? record is a pointer not a shell command", "2026-05-11"},
}

// ValidNodeTypes is the enum for the --type flag.
var ValidNodeTypes = []string{"decision", "fact", "pattern", "task", "observation"}

var MockProject = []ProjectRow{
	{"name", "agent-help"},
	{"version", "0.1.0"},
	{"status", "draft"},
	{"owner", "team-cli"},
	{"spec", "AHF-RFC.md"},
	{"repo", "github.com/Zate/agent-help"},
}
