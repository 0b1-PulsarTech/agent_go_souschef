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

func TestRenderTypeRefSections(t *testing.T) {
	t.Parallel()

	text := Render(
		repomodel.QueryHit{
			Symbol:    repomodel.Symbol{Name: "Repository", Kind: "interface", File: "service.go"},
			UsedBy:    []string{"CreateUser"},
			UsesTypes: []string{"User"},
		},
		false,
	)
	if !strings.Contains(text, "Used as type by:") || !strings.Contains(text, "CreateUser") {
		t.Errorf("missing incoming type-ref section: %q", text)
	}

	if !strings.Contains(text, "Uses types:") || !strings.Contains(text, "User") {
		t.Errorf("missing outgoing type-ref section: %q", text)
	}

	empty := Render(repomodel.QueryHit{Symbol: repomodel.Symbol{Name: "X", Kind: "type"}}, false)
	if strings.Contains(empty, "Used as type by") || strings.Contains(empty, "Uses types") {
		t.Errorf("empty sections should be omitted: %q", empty)
	}
}
