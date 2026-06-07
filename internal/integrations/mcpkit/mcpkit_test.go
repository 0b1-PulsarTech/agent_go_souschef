package mcpkit

import "testing"

func TestNew(t *testing.T) {
	t.Parallel()
	s := New("agent_go_souschef_test", "0.0.1")
	if s == nil || s.impl == nil {
		t.Fatal("expected initialized server")
	}
}
