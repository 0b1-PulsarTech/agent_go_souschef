// Package mcpsvc registers a repocontext.Service as MCP tools. It lives apart
// from repocontext so the domain package never imports the MCP SDK.
package mcpsvc

import (
	"context"
	"fmt"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/integrations/mcpkit"
	"github.com/0b1-PulsarTech/agent_go_souschef/pkg/repocontext"
)

// RegisterMCP wires every Service operation as its own MCP tool on s, so the
// LLM's catalog shows five entries picked by purpose. gate keeps the index
// fresh: read tools refresh through it (throttled), and souschef_sync forces a
// refresh regardless of the throttle.
func RegisterMCP(s *mcpkit.Server, svc repocontext.Service, gate *SyncGate) {
	mcpkit.Tool(s, "souschef_sync",
		"Build or refresh the symbol index for the current workspace.",
		func(ctx context.Context, _ SyncIn) (Result, error) {
			text, err := gate.Force(ctx)

			return result("souschef_sync", text, err)
		})

	mcpkit.Tool(
		s,
		"souschef_query",
		"Look up a symbol and return its direct callers/callees. Set expand=true for transitive deps.",
		func(ctx context.Context, in QueryIn) (Result, error) {
			gate.Ensure(ctx)
			text, err := svc.Query(ctx, in.Query, in.Expand)

			return result("souschef_query", text, err)
		},
	)

	mcpkit.Tool(s, "souschef_source",
		"Return the file path and source snippet for a named symbol.",
		func(ctx context.Context, in SourceIn) (Result, error) {
			gate.Ensure(ctx)
			text, err := svc.Source(ctx, in.Query)

			return result("souschef_source", text, err)
		})

	mcpkit.Tool(s, "souschef_changed",
		"List files modified in the workspace, optionally filtered by path scope.",
		func(ctx context.Context, in ChangedIn) (Result, error) {
			gate.Ensure(ctx)
			text, err := svc.Changed(ctx, in.Scope)

			return result("souschef_changed", text, err)
		})

	mcpkit.Tool(s, "souschef_shadows",
		"Report identifiers that shadow a builtin, an imported package, a "+
			"package-level symbol, or an outer variable. Use before naming locals "+
			"to avoid hiding stdlib names or enclosing variables.",
		func(ctx context.Context, in ShadowsIn) (Result, error) {
			gate.Ensure(ctx)
			text, err := svc.Shadows(ctx, in.Scope)

			return result("souschef_shadows", text, err)
		})
}

// result packages a Service call's output for MCP: it tags any error with the
// tool name so failures are attributable, and never wraps a nil error.
func result(tool, text string, err error) (Result, error) {
	if err != nil {
		return Result{}, fmt.Errorf("%s: %w", tool, err)
	}

	return Result{Text: text}, nil
}
