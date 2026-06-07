// Package mcpkit wraps github.com/modelcontextprotocol/go-sdk with a
// fuego/fiber-style registration API: callers say `mcpkit.Tool(server, name,
// desc, handler)` and never write the (*CallToolRequest, In) → (*CallToolResult,
// Out, error) shim by hand. The wrapper is the only place that touches the
// raw SDK, so swapping transports or SDK versions is a one-file change.
package mcpkit

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server is the public stdio-MCP server. Tools are registered with Tool[In,Out]
// (see tool.go) and the process runs until the MCP client disconnects.
type Server struct {
	impl *mcp.Server
}

// New constructs an empty server. Tools are registered via Tool[In,Out]; the
// server starts answering once Run is called.
func New(name, version string) *Server {
	return &Server{
		impl: mcp.NewServer(&mcp.Implementation{Name: name, Version: version}, nil),
	}
}

// Run blocks until the MCP client disconnects, using stdio as the transport.
// Cancellation through ctx unwinds the SDK's read loop cleanly.
func (s *Server) Run(ctx context.Context) error {
	return s.impl.Run(ctx, &mcp.StdioTransport{})
}
