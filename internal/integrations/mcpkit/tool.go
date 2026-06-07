package mcpkit

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Handler is the user-facing handler signature: take a typed input, return a
// typed output (or error). The wrapper handles the SDK's CallToolRequest /
// CallToolResult plumbing.
type Handler[In, Out any] func(ctx context.Context, in In) (Out, error)

// Tool registers a typed MCP tool on the server. Mirrors fuego.Post:
//
//	mcpkit.Tool(srv, "souschef_query", "Look up a symbol …",
//	    func(ctx context.Context, in QueryIn) (Result, error) { ... })
//
// One generic call replaces the four-argument SDK shim every handler would
// otherwise write inline.
func Tool[In, Out any](s *Server, name, description string, handle Handler[In, Out]) {
	mcp.AddTool(s.impl, &mcp.Tool{Name: name, Description: description},
		func(ctx context.Context, _ *mcp.CallToolRequest, in In) (*mcp.CallToolResult, Out, error) {
			out, err := handle(ctx, in)
			return nil, out, err
		})
}
