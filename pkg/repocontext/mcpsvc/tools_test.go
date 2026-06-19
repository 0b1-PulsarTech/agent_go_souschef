package mcpsvc

import (
	"testing"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/integrations/mcpkit"
	"github.com/0b1-PulsarTech/agent_go_souschef/pkg/repocontext"
)

// TestRegisterMCP_NoPanic catches any tool shape mismatch with the generic
// mcpkit.Tool wrapper at registration time. A zero Service is enough: the
// closures are only registered here, never invoked.
func TestRegisterMCP_NoPanic(t *testing.T) {
	t.Parallel()
	server := mcpkit.New("test", "0.0.0")
	RegisterMCP(server, repocontext.Service{})
}
