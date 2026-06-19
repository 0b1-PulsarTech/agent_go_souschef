package repocontext

import "github.com/0b1-PulsarTech/agent_go_souschef/internal/source/queryview"

// Service orchestrates symbol lookups against an indexer, a store, and a
// change-reporter. It holds only interface values, so it is passed by value.
type Service struct {
	indexer LanguageIndexer
	store   SymbolStore
	changes ChangeReporter
}

// New wires a Service from already-built collaborators. The three-argument
// shape lets remy auto-resolve every dep with RegisterConstructorArgs3.
func New(indexer LanguageIndexer, store SymbolStore, changes ChangeReporter) Service {
	return Service{indexer: indexer, store: store, changes: changes}
}

func renderCompact(hit QueryHit, expanded bool) string {
	return queryview.Render(hit, expanded)
}
