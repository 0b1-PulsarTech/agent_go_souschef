package repocontext

import (
	"context"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/integrations/mcpkit"
)

// MCP input/output types. These travel as JSON over the MCP transport, so the
// field tags drive the published jsonschema the LLM sees.
type (
	QueryIn struct {
		Query  string `json:"query"            jsonschema:"Symbol name or free-text query"`
		Expand bool   `json:"expand,omitempty" jsonschema:"Return transitive callers/callees"`
	}
	SyncIn    struct{}
	SourceIn  struct {
		Query string `json:"query" jsonschema:"Symbol name to locate"`
	}
	ChangedIn struct {
		Scope string `json:"scope,omitempty" jsonschema:"Optional path filter"`
	}
	Result struct {
		Text string `json:"text" jsonschema:"Compact human/LLM-readable result"`
	}
)

// RegisterMCP wires every Service operation as a distinct MCP tool. The four
// tools share one stdio server but each carries its own name, description,
// and input schema — the LLM's tool catalog shows four entries, picking by
// purpose rather than from a single discriminator field.
//
// The closures are intentionally tiny: each one just adapts the typed input
// to the Service method and packages the string result. Splitting them across
// four files would add four package imports for no readability win.
func RegisterMCP(s *mcpkit.Server, svc *Service) {
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
