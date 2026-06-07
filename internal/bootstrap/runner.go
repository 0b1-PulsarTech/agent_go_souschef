package bootstrap

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/integrations/mcpkit"
	"github.com/0b1-PulsarTech/agent_go_souschef/pkg/repocontext"
)

// RunMCP resolves the Service from the injector, registers the four MCP tools
// on a new mcpkit.Server, and blocks until the client disconnects.
func RunMCP(ctx context.Context, inj remy.Injector, cfg Config) error {
	svc, err := remy.Get[*repocontext.Service](inj)
	if err != nil {
		return fmt.Errorf("resolve service: %w", err)
	}
	server := mcpkit.New("agent_go_souschef", cfg.Version)
	repocontext.RegisterMCP(server, svc)
	if err := server.Run(ctx); err != nil {
		slog.Error("mcp server", "err", err)
		return err
	}
	return nil
}

// RunSync builds the index once and prints a one-line summary. It exists so
// users can bootstrap the index from a terminal before connecting an MCP
// client (Claude Code, Codex, …).
func RunSync(ctx context.Context, inj remy.Injector, out io.Writer) error {
	svc, err := remy.Get[*repocontext.Service](inj)
	if err != nil {
		return fmt.Errorf("resolve service: %w", err)
	}
	summary, err := svc.Sync(ctx)
	if err != nil {
		return fmt.Errorf("sync: %w", err)
	}
	if _, err := fmt.Fprintln(out, summary); err != nil {
		return err
	}
	return nil
}
