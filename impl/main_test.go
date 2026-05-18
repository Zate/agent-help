package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func runMem(t *testing.T, args ...string) string {
	t.Helper()

	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run . %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

func TestHumanHelpIncludesAgentBreadcrumb(t *testing.T) {
	out := runMem(t, "--help")
	if !strings.Contains(out, "LLM agent? Use --agent-help for token-optimized usage.") {
		t.Fatalf("--help output missing agent-help breadcrumb:\n%s", out)
	}
}

func TestAgentHelpIndexGolden(t *testing.T) {
	want := `ah1 mem :: project memory — store and query facts, decisions, and tasks
cmd node add <text> --type TYPE [--tag K:V...] :: store a new memory node
cmd node list [--type TYPE] [--limit int] [--cursor id] :: list memory nodes
cmd search query <text> [--type TYPE] [--limit int] [--cursor id] :: search nodes by text
cmd search similar <id> [--limit int] [--cursor id] :: find nodes similar to a given node
cmd project status :: show current project settings
cmd project set --key KEY --value VALUE :: set a project setting
more? mem <cmd> --agent-help
`

	if got := runMem(t, "--agent-help"); got != want {
		t.Fatalf("--agent-help mismatch\nwant:\n%sgot:\n%s", want, got)
	}
}

func TestAgentHelpDetailGolden(t *testing.T) {
	want := `ah2 mem node list
use mem node list [--type TYPE] [--limit int] [--cursor id]
flag --type:enum(decision|fact|pattern|task|observation) opt :: filter by node type
flag --limit:int opt default=10 :: max results to return
flag --cursor:id opt :: resume after the given node ID
ex mem node list --agent-out
ex mem node list --type decision --agent-out
ex mem node list --limit 3 --agent-out
`

	if got := runMem(t, "node", "list", "--agent-help"); got != want {
		t.Fatalf("node list --agent-help mismatch\nwant:\n%sgot:\n%s", want, got)
	}
}

func TestAgentOutListGolden(t *testing.T) {
	want := `ok nodes count=8 shown=8 more=0
nodes[#8]{id,type,tags,text}:
  n_101,decision,"project:mem|arch",Use Cobra for the reference CLI implementation
  n_102,fact,"project:mem|db",Mock data lives in data.go; no real storage in the reference impl
  n_103,pattern,db|ops,Use connection pool max 20 in all environments
  n_104,fact,"project:mem|api",AHF ok/err envelope always precedes the TOON body
  n_105,decision,"project:mem|spec",TOON is the encoding for --agent-out result bodies
  n_106,task,"project:mem|todo",Write reference Go implementation for agent-help
  n_107,pattern,agent|cli,Always include next records for paginated results
  n_108,fact,"project:mem|spec",more? record is a pointer not a shell command
next inspect "mem search query <text> --agent-out"
`

	if got := runMem(t, "node", "list", "--agent-out"); got != want {
		t.Fatalf("node list --agent-out mismatch\nwant:\n%sgot:\n%s", want, got)
	}
}

func TestAgentOutPaginationGolden(t *testing.T) {
	want := `ok nodes count=8 shown=3 more=1 cursor=c_n_104
warn truncated shown=3 total=8
nodes[#3]{id,type,tags,text}:
  n_101,decision,"project:mem|arch",Use Cobra for the reference CLI implementation
  n_102,fact,"project:mem|db",Mock data lives in data.go; no real storage in the reference impl
  n_103,pattern,db|ops,Use connection pool max 20 in all environments
next "mem node list --limit 3 --cursor c_n_104 --agent-out"
next inspect "mem search query <text> --agent-out"
`

	if got := runMem(t, "node", "list", "--limit", "3", "--agent-out"); got != want {
		t.Fatalf("paginated node list mismatch\nwant:\n%sgot:\n%s", want, got)
	}
}

func TestAgentOutErrorGolden(t *testing.T) {
	want := `err invalid_enum flag=--type got=badval
hint --type enum(decision|fact|pattern|task|observation)
use mem node add <text> --type TYPE
`

	if got := runMem(t, "node", "add", "test node", "--type", "badval", "--agent-out"); got != want {
		t.Fatalf("invalid enum error mismatch\nwant:\n%sgot:\n%s", want, got)
	}
}
