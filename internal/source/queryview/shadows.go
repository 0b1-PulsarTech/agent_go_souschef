package queryview

import (
	"fmt"
	"strings"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
)

// RenderShadows formats shadowing findings grouped by file, compact enough for
// an MCP response. The input is expected pre-sorted by file/line. scope, when
// set, is echoed so the caller sees which filter produced the list.
func RenderShadows(shadows []repomodel.Shadow, scope string) string {
	if len(shadows) == 0 {
		if scope != "" {
			return "No variable shadowing in " + scope + "."
		}

		return "No variable shadowing detected."
	}

	var b strings.Builder

	fmt.Fprintf(&b, "Shadowing — %d finding%s", len(shadows), plural(len(shadows)))

	if scope != "" {
		fmt.Fprintf(&b, " (scope: %s)", scope)
	}

	b.WriteByte('\n')

	current := ""
	for _, sh := range shadows {
		if sh.File != current {
			current = sh.File
			b.WriteString(current)
			b.WriteByte('\n')
		}

		fmt.Fprintf(&b, "  L%d:%d  %s → %s\n", sh.Line, sh.Column, sh.Name, shadowPhrase(sh))
	}

	return strings.TrimSpace(b.String())
}

func shadowPhrase(sh repomodel.Shadow) string {
	switch sh.Origin {
	case "builtin":
		return "shadows a builtin/predeclared identifier"
	case "import":
		return "shadows " + sh.Detail
	case "package":
		return "shadows a package-level symbol (" + sh.Detail + ")"
	default:
		return "shadows an outer variable (" + sh.Detail + ")"
	}
}

func plural(n int) string {
	if n == 1 {
		return ""
	}

	return "s"
}
