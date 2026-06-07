package repocontext

import (
	"testing"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/integrations/mcpkit"
)

// TestRegisterMCP_NoPanic catches any tool shape mismatch with the generic
// mcpkit.Tool wrapper at registration time — if Sync/Query/Source/Changed
// don't fit the (ctx, In) (Out, error) shape the call panics or fails to
// compile, both surfaced here without booting the SDK transport.
func TestRegisterMCP_NoPanic(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	server := mcpkit.New("test", "0.0.0")
	RegisterMCP(server, svc)
}
