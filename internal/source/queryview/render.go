package queryview

import (
	"strings"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
)

func Render(hit repomodel.QueryHit, expanded bool) string {
	if hit.Symbol.Name == "" {
		return renderText(hit.TextHits)
	}
	var b strings.Builder
	line(&b, "Symbol: ", hit.Symbol.Name)
	line(&b, "Kind: ", title(hit.Symbol.Kind))
	section(&b, "Methods", hit.Methods)
	section(&b, "Calls", hit.Calls)
	section(&b, "Implementations", hit.Implementations)
	section(&b, "Referenced by", hit.Callers)
	if expanded && hit.Symbol.Signature != "" {
		line(&b, "Signature: ", hit.Symbol.Signature)
	}
	line(&b, "Defined: ", hit.Symbol.File)
	return strings.TrimSpace(b.String())
}

func title(text string) string {
	if text == "" {
		return text
	}
	return strings.ToUpper(text[:1]) + text[1:]
}
