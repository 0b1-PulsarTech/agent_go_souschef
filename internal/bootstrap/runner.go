package bootstrap

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/integrations/mcpkit"
	"github.com/0b1-PulsarTech/agent_go_souschef/pkg/repocontext"
	"github.com/0b1-PulsarTech/agent_go_souschef/pkg/repocontext/mcpsvc"
)

// RunMCP resolves the Service from the injector, registers the five MCP tools
// on a new mcpkit.Server, and blocks until the client disconnects. One SyncGate
// owns every sync: the startup refresh below, the throttled refresh before each
// read tool, and the explicit souschef_sync.
func RunMCP(ctx context.Context, inj remy.Injector, cfg Config) error {
	svc, err := remy.Get[repocontext.Service](inj)
	if err != nil {
		return fmt.Errorf("resolve service: %w", err)
	}

	gate := mcpsvc.NewSyncGate(svc, time.Now, mcpsvc.DefaultSyncInterval)
	// Build the index up-front so the first query has data to answer from.
	// A failure here is non-fatal: the server still starts and the client can
	// re-run souschef_sync once the workspace is fixed. We only log it.
	if summary, syncErr := gate.Force(ctx); syncErr != nil {
		slog.Warn("initial sync failed", "err", syncErr)
	} else {
		slog.Info("initial sync", "result", summary)
	}

	server := mcpkit.New("agent_go_souschef", cfg.Version)
	mcpsvc.RegisterMCP(server, svc, gate)

	if err := server.Run(ctx); err != nil {
		slog.Error("mcp server", "err", err)

		return fmt.Errorf("run mcp server: %w", err)
	}

	return nil
}

// RunSync builds the index once and prints a one-line summary. It exists so
// users can bootstrap the index from a terminal before connecting an MCP
// client (Claude Code, Codex, …).
func RunSync(ctx context.Context, inj remy.Injector, out io.Writer) error {
	svc, err := remy.Get[repocontext.Service](inj)
	if err != nil {
		return fmt.Errorf("resolve service: %w", err)
	}

	summary, err := svc.Sync(ctx)
	if err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	if _, err := fmt.Fprintln(out, summary); err != nil {
		return fmt.Errorf("write summary: %w", err)
	}

	return nil
}

// RunShadows builds the index (so the report reflects the current tree) and
// prints every shadowing finding, optionally narrowed to scope.
func RunShadows(ctx context.Context, inj remy.Injector, out io.Writer, scope string) error {
	svc, err := remy.Get[repocontext.Service](inj)
	if err != nil {
		return fmt.Errorf("resolve service: %w", err)
	}

	if _, err := svc.Sync(ctx); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	report, err := svc.Shadows(ctx, scope)
	if err != nil {
		return fmt.Errorf("shadows: %w", err)
	}

	if _, err := fmt.Fprintln(out, report); err != nil {
		return fmt.Errorf("write report: %w", err)
	}

	return nil
}
