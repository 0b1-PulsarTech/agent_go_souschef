package reposqlite

import (
	"context"
	"testing"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
)

func TestWrite(t *testing.T) {
	t.Parallel()
	store, err := Open(t.TempDir() + "/db.sqlite")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := store.Write(context.Background(), repomodel.Snapshot{}); err != nil {
		t.Fatalf("write: %v", err)
	}
}
