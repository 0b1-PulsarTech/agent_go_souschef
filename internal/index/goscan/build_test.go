package goscan

import (
	"context"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	t.Parallel()
	idx := New(filepath.Join("..", "..", "..", "test", "fixtures", "sample"))
	snap, err := idx.Build(context.Background())
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(snap.Symbols) == 0 {
		t.Fatal("expected symbols")
	}
}
