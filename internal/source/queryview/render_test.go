package queryview

import (
	"strings"
	"testing"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
)

func TestRender(t *testing.T) {
	t.Parallel()
	text := Render(
		repomodel.QueryHit{
			Symbol: repomodel.Symbol{Name: "CreateUser", Kind: "function", File: "user.go"},
		},
		false,
	)
	if !strings.Contains(text, "CreateUser") {
		t.Fatalf("got %q", text)
	}
}
