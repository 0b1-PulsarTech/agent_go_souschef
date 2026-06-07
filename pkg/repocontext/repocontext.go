package repocontext

import "github.com/0b1-PulsarTech/agent_go_souschef/internal/source/queryview"

// Service orchestrates symbol lookups against an indexer, a store, and a
// change-reporter. Collaborators arrive via the constructor — no globals,
// no service-locator. Each collaborator already knows its own root path.
type Service struct {
	indexer LanguageIndexer
	store   SymbolStore
	changes ChangeReporter
}

// New wires a Service from already-built collaborators. Concrete construction
// (open SQLite, build goscan, open git repo) is the bootstrap layer's job.
// The three-argument shape lets remy auto-resolve every dep with
// RegisterConstructorArgs3 — no closure-binders required.
func New(indexer LanguageIndexer, store SymbolStore, changes ChangeReporter) *Service {
	return &Service{indexer: indexer, store: store, changes: changes}
}

func renderCompact(hit QueryHit, expanded bool) string {
	return queryview.Render(hit, expanded)
}
