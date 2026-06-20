package goscan

import (
	"go/ast"
	"testing"
)

func TestMatchesType(t *testing.T) {
	t.Parallel()

	decl := &ast.GenDecl{Specs: []ast.Spec{&ast.TypeSpec{Name: &ast.Ident{Name: "User"}}}}
	if !matchesType(decl, "User") {
		t.Fatal("expected match")
	}
}
