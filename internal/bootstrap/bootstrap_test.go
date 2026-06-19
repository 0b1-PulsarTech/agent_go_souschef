package bootstrap

import (
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/reposqlite"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/source/gitprobe"
	"github.com/0b1-PulsarTech/agent_go_souschef/pkg/repocontext"
)

// TestDoInjections_ResolvesEveryType is the canonical bootstrap smoke: if any
// registration is missing or shape-mismatched, one of these Get calls fails.
// We point Config at a temp dir so opening the SQLite store touches a real
// (but empty) file.
func TestDoInjections_ResolvesEveryType(t *testing.T) {
	t.Parallel()
	cfg := Config{Root: t.TempDir(), Version: "test"}
	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
	DoInjections(inj, cfg)

	if _, err := remy.Get[*reposqlite.Store](inj); err != nil {
		t.Fatalf("store: %v", err)
	}
	if _, err := remy.Get[goscan.Indexer](inj); err != nil {
		t.Fatalf("indexer: %v", err)
	}
	if _, err := remy.Get[*gitprobe.Probe](inj); err != nil {
		t.Fatalf("probe: %v", err)
	}
	if _, err := remy.Get[repocontext.Service](inj); err != nil {
		t.Fatalf("service: %v", err)
	}
}
