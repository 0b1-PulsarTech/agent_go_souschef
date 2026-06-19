// Package mcpsvc registers a repocontext.Service as MCP tools. It lives apart
// from repocontext so the domain package never imports the MCP SDK.
package mcpsvc

import (
	"context"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/integrations/mcpkit"
	"github.com/0b1-PulsarTech/agent_go_souschef/pkg/repocontext"
)

// RegisterMCP wires every Service operation as its own MCP tool on s, so the
// LLM's catalog shows four entries picked by purpose.
func RegisterMCP(s *mcpkit.Server, svc repocontext.Service) {
	mcpkit.Tool(s, "souschef_sync",
		"Build or refresh the symbol index for the current workspace.",
		func(ctx context.Context, _ SyncIn) (Result, error) {
			text, err := svc.Sync(ctx)
			return Result{Text: text}, err
		})

	mcpkit.Tool(s, "souschef_query",
		"Look up a symbol and return its direct callers/callees. Set expand=true for transitive deps.",
		func(ctx context.Context, in QueryIn) (Result, error) {
			text, err := svc.Query(ctx, in.Query, in.Expand)
			return Result{Text: text}, err
		})

	mcpkit.Tool(s, "souschef_source",
		"Return the file path and source snippet for a named symbol.",
		func(ctx context.Context, in SourceIn) (Result, error) {
			text, err := svc.Source(ctx, in.Query)
			return Result{Text: text}, err
		})

	mcpkit.Tool(s, "souschef_changed",
		"List files modified in the workspace, optionally filtered by path scope.",
		func(ctx context.Context, in ChangedIn) (Result, error) {
			text, err := svc.Changed(ctx, in.Scope)
			return Result{Text: text}, err
		})
}
