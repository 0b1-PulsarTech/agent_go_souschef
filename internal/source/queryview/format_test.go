package queryview

import (
	"strings"
	"testing"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
)

func TestRenderText(t *testing.T) {
	t.Parallel()

	text := renderText([]repomodel.TextHit{{Path: "a.go", Snippet: "CreateUser"}})
	if !strings.Contains(text, "a.go") {
		t.Fatalf("got %q", text)
	}
}
