package goscan

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan/symbols"
)

func formatDecl(fset *token.FileSet, decl ast.Decl, query string) string {
	parts := strings.SplitN(query, ".", 2)
	forMethod := len(parts) == 2
	switch typed := decl.(type) {
	case *ast.FuncDecl:
		if matchesFunc(typed, parts, forMethod) {
			return nodeText(fset, typed)
		}
	case *ast.GenDecl:
		if !forMethod && matchesType(typed, query) {
			return nodeText(fset, typed)
		}
	}
	return ""
}

func matchesFunc(decl *ast.FuncDecl, parts []string, forMethod bool) bool {
	if !forMethod {
		return decl.Name.Name == parts[0]
	}
	return decl.Recv != nil && decl.Name.Name == parts[1] && symbols.RecvName(decl) == parts[0]
}

func matchesType(decl *ast.GenDecl, query string) bool {
	for _, spec := range decl.Specs {
		ts, ok := spec.(*ast.TypeSpec)
		if ok && ts.Name.Name == query {
			return true
		}
	}
	return false
}
