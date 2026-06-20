package reposqlite

import (
	"context"
	"testing"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
)

func TestLookup(t *testing.T) {
	t.Parallel()

	store, err := Open(t.TempDir() + "/db.sqlite")
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	snap := repomodel.Snapshot{Symbols: []repomodel.Symbol{{ID: 1, Name: "CreateUser"}}}
	if err = store.Write(context.Background(), snap); err != nil {
		t.Fatalf("write: %v", err)
	}

	rows, err := store.Lookup(context.Background(), "CreateUser")
	if err != nil || len(rows) != 1 {
		t.Fatalf("lookup: %v len=%d", err, len(rows))
	}
}
