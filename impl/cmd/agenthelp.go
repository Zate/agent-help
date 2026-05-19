// Package cmd implements the agent-help convention for the mem CLI.
// This file provides the core AHF output helpers used by all commands.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/toon-format/toon-go"
)

// AHF record helpers — all --agent-help and --agent-out output goes through these.

func ah1(tool, purpose string)     { fmt.Printf("ah1 %s :: %s\n", tool, purpose) }
func ah2(tool, path string)        { fmt.Printf("ah2 %s %s\n", tool, path) }
func cmdEntry(sig, purpose string) { fmt.Printf("cmd %s :: %s\n", sig, purpose) }
func use(invocation string)        { fmt.Printf("use %s\n", invocation) }
func morePtr(tool, path string)    { fmt.Printf("more? %s %s --agent-help\n", tool, path) }
func argEntry(name, typ, req, purpose string) {
	fmt.Printf("arg %s:%s %s :: %s\n", name, typ, req, purpose)
}
func flagEntry(name, typ, req, purpose string) {
	fmt.Printf("flag --%s:%s %s :: %s\n", name, typ, req, purpose)
}
func flagEntryDefault(name, typ, req, def, purpose string) {
	fmt.Printf("flag --%s:%s %s default=%s :: %s\n", name, typ, req, def, purpose)
}
func ex(invocation string) { fmt.Printf("ex %s\n", invocation) }

// --agent-out AHF protocol envelope records.

func okLine(kind string, meta ...string) {
	if len(meta) == 0 {
		fmt.Printf("ok %s\n", kind)
	} else {
		fmt.Printf("ok %s %s\n", kind, strings.Join(meta, " "))
	}
}
func errLine(code string, meta ...string) {
	if len(meta) == 0 {
		fmt.Printf("err %s\n", code)
	} else {
		fmt.Printf("err %s %s\n", code, strings.Join(meta, " "))
	}
}
func warnLine(code string, meta ...string) {
	if len(meta) == 0 {
		fmt.Printf("warn %s\n", code)
	} else {
		fmt.Printf("warn %s %s\n", code, strings.Join(meta, " "))
	}
}
func hint(msg string)           { fmt.Printf("hint %s\n", msg) }
func useLine(invocation string) { fmt.Printf("use %s\n", invocation) }
func nextLine(label, cmd string) {
	cmd = strings.ReplaceAll(cmd, `\`, `\\`)
	cmd = strings.ReplaceAll(cmd, `"`, `\"`)
	if label == "" {
		fmt.Printf("next \"%s\"\n", cmd)
	} else {
		fmt.Printf("next %s \"%s\"\n", label, cmd)
	}
}

// kvLine emits a simple key-value line for bare scalar results
// where a TOON block would be overkill (§8 of SPEC.md).
func kvLine(key, value string) {
	if strings.ContainsAny(value, " \t") {
		fmt.Printf("%s %q\n", key, value)
	} else {
		fmt.Printf("%s %s\n", key, value)
	}
}

// -- TOON output via github.com/toon-format/toon-go --
//
// emitTOON marshals v using the toon-go library and writes it to stdout.
// v should be a struct (or slice of structs) with `toon:"field"` tags.
// The ok/err envelope line must be emitted BEFORE calling emitTOON.
//
// Example:
//
//	type NodeRow struct {
//	    ID   string `toon:"id"`
//	    Type string `toon:"type"`
//	    Text string `toon:"text"`
//	}
//	type NodeList struct {
//	    Nodes []NodeRow `toon:"nodes"`
//	}
//	okLine("nodes", "count=3", "more=0")
//	emitTOON(NodeList{Nodes: rows})

func emitTOON(v any) {
	out, err := toon.MarshalString(v, toon.WithLengthMarkers(true))
	if err != nil {
		fmt.Fprintf(os.Stderr, "toon marshal error: %v\n", err)
		return
	}
	fmt.Print(out)
	// Ensure next/hint AHF records always start on a new line after the TOON block.
	if len(out) > 0 && out[len(out)-1] != '\n' {
		fmt.Println()
	}
}
