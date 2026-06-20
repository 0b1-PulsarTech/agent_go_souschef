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
		_, _ = fmt.Fprintln(stderr, err)

		return 1
	}

	cfg := bootstrap.Config{Root: root, Version: version}
	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
	bootstrap.DoInjections(inj, cfg)

	switch args[0] {
	case "mcp":
		if err := bootstrap.RunMCP(ctx, inj, cfg); err != nil {
			_, _ = fmt.Fprintln(stderr, err)

			return 1
		}

		return 0
	case "sync":
		if err := bootstrap.RunSync(ctx, inj, stdout); err != nil {
			_, _ = fmt.Fprintln(stderr, err)

			return 1
		}

		return 0
	case "shadows":
		if err := bootstrap.RunShadows(ctx, inj, stdout, scopeArg(args)); err != nil {
			_, _ = fmt.Fprintln(stderr, err)

			return 1
		}

		return 0
	default:
		_, _ = fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		printUsage(stderr)

		return 1
	}
}

// scopeArg returns the optional path filter passed after a subcommand, e.g.
// `shadows internal/index`. It is empty when no scope is given.
func scopeArg(args []string) string {
	const scopeIdx = 1 // args[0] is the subcommand; args[1], if present, is the scope.
	if len(args) <= scopeIdx {
		return ""
	}

	return args[scopeIdx]
}

func printUsage(w *os.File) {
	_, _ = fmt.Fprintln(w, `agent_go-souschef — semantic repository index for LLMs

Usage:
  agent_go-souschef mcp              start the MCP stdio server (builds the index on startup)
  agent_go-souschef sync             build/refresh the index once (optional; mcp does this too)
  agent_go-souschef shadows [scope]  report shadowed builtins/imports/variables (optional path filter)`)
}
