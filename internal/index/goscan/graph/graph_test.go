package graph

import (
	"go/ast"
	"go/types"
	"testing"
)

func TestFullNameNil(t *testing.T) {
	t.Parallel()
	if got := FullName(nil); got != "" {
		t.Fatalf("got %q", got)
	}
}

func TestFullNameNoPkg(t *testing.T) {
	t.Parallel()
	// types.Universe objects have no package
	obj := types.Universe.Lookup("error")
	if got := FullName(obj); got != "" {
		t.Fatalf("expected empty for universe object, got %q", got)
	}
}

func TestCalledObjectNilPkg(t *testing.T) {
	t.Parallel()
	if got := CalledObject(nil, &ast.CallExpr{}); got != nil {
		t.Fatal("expected nil for nil pkg")
	}
}

func TestImplementRelation(t *testing.T) {
	t.Parallel()
	r := ImplementRelation(1, 2)
	if r.FromID != 1 || r.ToID != 2 || r.Kind != "implement" {
		t.Fatalf("unexpected %+v", r)
	}
}
