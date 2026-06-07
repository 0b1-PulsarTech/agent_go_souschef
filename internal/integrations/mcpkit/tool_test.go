package mcpkit

import (
	"context"
	"testing"
)

// TestTool_RegistersWithoutPanic exercises the generic Tool wrapper end-to-end:
// if the SDK shim is mis-shaped the registration call panics or fails to compile.
// We don't drive the MCP transport here — that's covered by the integration
// smoke under pkg/repocontext/mcptools_test.go.
func TestTool_RegistersWithoutPanic(t *testing.T) {
	t.Parallel()
	type In struct{ Q string }
	type Out struct{ Answer string }

	s := New("test", "0.0.0")
	Tool(s, "echo", "echoes the query back",
		func(_ context.Context, in In) (Out, error) {
			return Out{Answer: in.Q}, nil
		})
}
