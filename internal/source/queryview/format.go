package queryview

import (
	"strings"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
)

func renderText(hits []repomodel.TextHit) string {
	if len(hits) == 0 {
		return "No matches."
	}
	var b strings.Builder
	b.WriteString("Found:\n")
	for _, hit := range hits {
		line(&b, "  ", hit.Snippet)
	}
	b.WriteString("Files:\n")
	for _, hit := range hits {
		line(&b, "  ", hit.Path)
	}
	return strings.TrimSpace(b.String())
}

func section(b *strings.Builder, title string, items []string) {
	if len(items) == 0 {
		return
	}
	b.WriteString("\n" + title + ":\n")
	for _, item := range items {
		line(b, "  ", item)
	}
}

func line(b *strings.Builder, prefix, value string) { b.WriteString(prefix + value + "\n") }
