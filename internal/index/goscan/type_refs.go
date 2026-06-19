package goscan

import (
	"go/ast"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan/graph"
	"golang.org/x/tools/go/packages"
)

// addTypeRefs records type-reference edges in a final pass, after every symbol
// has been registered. Resolving here (rather than inline during the call walk)
// is required: type references usually point at types declared in a sibling
// file or another package, which are not yet registered when their referencing
// declaration is first visited.
func (b *snapshotBuilder) addTypeRefs() {
	seen := map[[2]int64]struct{}{}
	for _, pkg := range b.pkgs {
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				b.refsInDecl(pkg, decl, seen)
			}
		}
	}
}

func (b *snapshotBuilder) refsInDecl(
	pkg *packages.Package,
	decl ast.Decl,
	seen map[[2]int64]struct{},
) {
	switch typed := decl.(type) {
	case *ast.FuncDecl:
		b.refsFromNode(pkg, typed, b.declID(pkg, typed.Name), seen)
	case *ast.GenDecl:
		for _, spec := range typed.Specs {
			if ts, ok := spec.(*ast.TypeSpec); ok {
				b.refsFromNode(pkg, ts, b.declID(pkg, ts.Name), seen)
			}
		}
	}
}

func (b *snapshotBuilder) declID(pkg *packages.Package, name *ast.Ident) int64 {
	return b.ids[pkg.TypesInfo.Defs[name]]
}

// refsFromNode walks node and records a ref edge for every identifier that uses
// an indexed named type.
func (b *snapshotBuilder) refsFromNode(
	pkg *packages.Package,
	node ast.Node,
	fromID int64,
	seen map[[2]int64]struct{},
) {
	if fromID == 0 {
		return
	}
	ast.Inspect(node, func(n ast.Node) bool {
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		target := graph.ReferencedType(pkg, ident)
		if target == nil {
			return true
		}
		toID := b.ids[target]
		if toID == 0 {
			toID = b.names[graph.FullName(target)]
		}
		b.addTypeRef(fromID, toID, seen)
		return true
	})
}

func (b *snapshotBuilder) addTypeRef(fromID, toID int64, seen map[[2]int64]struct{}) {
	if toID == 0 || toID == fromID {
		return
	}
	key := [2]int64{fromID, toID}
	if _, dup := seen[key]; dup {
		return
	}
	seen[key] = struct{}{}
	b.snapshot.TypeRefs = append(b.snapshot.TypeRefs, graph.TypeRefRelation(fromID, toID))
}
