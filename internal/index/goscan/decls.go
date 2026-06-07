package goscan

import (
	"go/ast"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan/symbols"
	"golang.org/x/tools/go/packages"
)

func (b *snapshotBuilder) addDecl(
	pkg *packages.Package,
	path string,
	decl ast.Decl,
	sum *symbols.FileSummary,
) {
	switch typed := decl.(type) {
	case *ast.FuncDecl:
		b.addFunc(pkg, path, typed, sum)
	case *ast.GenDecl:
		for _, spec := range typed.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if ok {
				b.addType(pkg, path, ts, sum)
			}
		}
	}
}
