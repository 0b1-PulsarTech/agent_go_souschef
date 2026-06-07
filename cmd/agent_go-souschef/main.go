// Command agent_go-souschef is the CLI entry point. It only dispatches to
// the three first-class subcommands — mcp, sync, hook — and delegates every
// substantive operation to internal/bootstrap (for index-backed runs) or
// internal/integrations/hooksetup (for the stateless hook handlers).
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/bootstrap"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/integrations/hooksetup"
)

// version is reported to MCP clients in the implementation handshake. Bumped
// alongside go.mod when the public surface changes.
const version = "0.1.0"

func main() {
	os.Exit(run(context.Background(), os.Args[1:], os.Stdout, os.Stderr))
}

func run(ctx context.Context, args []string, stdout, stderr *os.File) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 1
	}

	// hook is stateless and bypasses bootstrap (no index needed).
	if args[0] == "hook" {
		return hooksetup.Run(args[1:])
	}

	cfg := bootstrap.Config{Root: ".", Version: version}
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
  agent_go-souschef mcp                start the MCP stdio server (primary surface)
  agent_go-souschef sync               build/refresh the index once (bootstrap convenience)
  agent_go-souschef hook install --claude [--codex --cursor --gemini]
                                       wire up PreToolUse hooks
  agent_go-souschef hook run <target>  invoked by the host's hook machinery`)
}
