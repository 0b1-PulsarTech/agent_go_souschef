package reposqlite

import "testing"

func TestOpen(t *testing.T) {
	t.Parallel()
	store, err := Open(t.TempDir() + "/index.db")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if store.q == nil {
		t.Fatal("expected queries")
	}
}
