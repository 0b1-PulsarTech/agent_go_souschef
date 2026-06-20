package repocontext

import "testing"

func TestContractsCompile(t *testing.T) {
	t.Parallel()

	var _ LanguageIndexer

	var _ SymbolStore

	var _ ChangeReporter
}
