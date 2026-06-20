package queryview

import (
	"strings"
	"testing"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
)

func TestRenderShadowsEmpty(t *testing.T) {
	t.Parallel()

	if got := RenderShadows(nil, ""); got != "No variable shadowing detected." {
		t.Errorf("empty render = %q", got)
	}

	if got := RenderShadows(nil, "pkg/x"); !strings.Contains(got, "pkg/x") {
		t.Errorf("scoped empty render should mention scope, got %q", got)
	}
}

func TestRenderShadows(t *testing.T) {
	t.Parallel()

	out := RenderShadows([]repomodel.Shadow{
		{File: "a.go", Line: 3, Column: 2, Name: "len", Origin: "builtin", Detail: "predeclared"},
		{
			File:   "a.go",
			Line:   5,
			Column: 1,
			Name:   "ctx",
			Origin: "import",
			Detail: `import "context"`,
		},
		{
			File:   "b.go",
			Line:   9,
			Column: 4,
			Name:   "x",
			Origin: "outer",
			Detail: "declared at b.go:7",
		},
	}, "")

	for _, want := range []string{
		"Shadowing — 3 findings",
		"a.go",
		"L3:2  len → shadows a builtin/predeclared identifier",
		`L5:1  ctx → shadows import "context"`,
		"b.go",
		"L9:4  x → shadows an outer variable (declared at b.go:7)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("render missing %q in:\n%s", want, out)
		}
	}
	// Each file header appears once even with multiple findings.
	if n := strings.Count(out, "a.go\n"); n != 1 {
		t.Errorf("a.go header repeated %d times:\n%s", n, out)
	}
}
