// Command agent_go-souschef is the CLI entry point. It dispatches to the two
// subcommands — mcp and sync — and delegates the work to internal/bootstrap.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/bootstrap"
)

// version is reported to MCP clients in the implementation handshake.
const version = "0.1.0"

func main() {
	os.Exit(run(context.Background(), os.Args[1:], os.Stdout, os.Stderr))
}

func run(ctx context.Context, args []string, stdout, stderr *os.File) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 1
	}

	root, err := filepath.Abs(".")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	cfg := bootstrap.Config{Root: root, Version: version}
	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
	bootstrap.DoInjections(inj, cfg)

	switch args[0] {
	case "mcp":
		if err := bootstrap.RunMCP(ctx, inj, cfg); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		return 0
	case "sync":
		if err := bootstrap.RunSync(ctx, inj, stdout); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		printUsage(stderr)
		return 1
	}
}

func printUsage(w *os.File) {
	fmt.Fprintln(w, `agent_go-souschef — semantic repository index for LLMs

Usage:
  agent_go-souschef mcp     start the MCP stdio server (builds the index on startup)
  agent_go-souschef sync    build/refresh the index once (optional; mcp does this too)`)
}
