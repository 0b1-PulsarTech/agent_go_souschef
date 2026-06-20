package goscan

import (
	"testing"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan/symbols"
)

func TestFileSummaryText(t *testing.T) {
	t.Parallel()

	txt := symbols.FileSummary{Pkg: "sample/user", Exports: []string{"CreateUser"}}.Text()
	if txt == "" {
		t.Fatal("expected summary")
	}
}
