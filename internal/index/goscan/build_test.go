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

// TestBuildWorkspace covers a go.work monorepo: the indexer must walk every
// module (a bare "./..." fails at a workspace root) and must not panic on
// generic receivers such as Box[T].
func TestBuildWorkspace(t *testing.T) {
	t.Parallel()

	idx := New(filepath.Join("..", "..", "..", "test", "fixtures", "workspace"))

	snap, err := idx.Build(context.Background())
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	want := map[string]bool{"Greeter.Hello": false, "Box.Get": false, "Box.Set": false}
	pkgs := map[string]bool{}

	for _, sym := range snap.Symbols {
		if _, ok := want[sym.Name]; ok {
			want[sym.Name] = true
		}

		pkgs[sym.Package] = true
	}

	for name, found := range want {
		if !found {
			t.Errorf("symbol %q from one of the workspace modules was not indexed", name)
		}
	}

	if !pkgs["wsmod/a"] || !pkgs["wsmod/b"] {
		t.Errorf("expected symbols from both modules, got packages %v", pkgs)
	}
}

// TestBuildMultiModule covers a repo with a nested module but no go.work:
// the root module's "./..." cannot reach the nested module, so each must be
// loaded on its own.
func TestBuildMultiModule(t *testing.T) {
	t.Parallel()

	idx := New(filepath.Join("..", "..", "..", "test", "fixtures", "multimod"))

	snap, err := idx.Build(context.Background())
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	pkgs := map[string]bool{}
	for _, sym := range snap.Symbols {
		pkgs[sym.Package] = true
	}

	if !pkgs["multimod"] || !pkgs["multimod/datastore"] {
		t.Errorf("expected root and nested-module symbols, got packages %v", pkgs)
	}
}
