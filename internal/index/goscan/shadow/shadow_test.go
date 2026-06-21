package shadow_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan/shadow"
)

// TestAnalyzer exercises every origin the pass classifies (builtin, import,
// package symbol, outer variable) and the else-if init pattern whose nested
// scope is easy to treat as a sibling. The `// want` assertions live next to
// each shadowing site in testdata/src/shadowy.
func TestAnalyzer(t *testing.T) {
	t.Parallel()

	analysistest.Run(t, analysistest.TestData(), shadow.Analyzer, "shadowy")
}
