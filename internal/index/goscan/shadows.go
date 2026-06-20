package goscan

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
	"golang.org/x/tools/go/packages"
)

// addShadows records every declaration that hides an identifier visible in an
// enclosing scope: a predeclared builtin/type (len, error, string, min…), an
// imported package name, a package-level symbol, or a variable/parameter from
// an outer function or block. Resolution walks the type-checker's lexical
// scopes, so it needs no external linter. Findings are deduped by source
// position because one file may be type-checked under several package variants
// (its test build, say), which would otherwise report the same shadow twice.
func (b *snapshotBuilder) addShadows() {
	seen := map[token.Position]struct{}{}

	for _, pkg := range b.pkgs {
		for ident, obj := range pkg.TypesInfo.Defs {
			b.addShadow(pkg, ident, obj, seen)
		}
	}
}

func (b *snapshotBuilder) addShadow(
	pkg *packages.Package,
	ident *ast.Ident,
	obj types.Object,
	seen map[token.Position]struct{},
) {
	if obj == nil || obj.Name() == "_" {
		return
	}

	inner := obj.Parent()
	if inner == nil { // struct fields and methods are not lexically scoped.
		return
	}

	enclosing := inner.Parent()
	if enclosing == nil {
		return
	}

	scope, hidden := enclosing.LookupParent(obj.Name(), obj.Pos())
	if hidden == nil || hidden == obj {
		return
	}

	pos := pkg.Fset.Position(ident.Pos())
	if _, dup := seen[pos]; dup {
		return
	}

	seen[pos] = struct{}{}

	b.snapshot.Shadows = append(b.snapshot.Shadows, repomodel.Shadow{
		File:   rel(b.root, pkg.Fset, ident.Pos()),
		Line:   pos.Line,
		Column: pos.Column,
		Name:   obj.Name(),
		Origin: shadowOrigin(pkg, scope, hidden),
		Detail: b.shadowDetail(pkg, scope, hidden),
	})
}

func shadowOrigin(pkg *packages.Package, scope *types.Scope, hidden types.Object) string {
	switch {
	case scope == types.Universe:
		return "builtin"
	case isPkgName(hidden):
		return "import"
	case scope == pkg.Types.Scope():
		return "package"
	default:
		return "outer"
	}
}

func (b *snapshotBuilder) shadowDetail(
	pkg *packages.Package,
	scope *types.Scope,
	hidden types.Object,
) string {
	if scope == types.Universe {
		return "predeclared"
	}

	if pkgName, ok := hidden.(*types.PkgName); ok {
		return fmt.Sprintf("import %q", pkgName.Imported().Path())
	}

	loc := rel(b.root, pkg.Fset, hidden.Pos())

	return fmt.Sprintf("declared at %s:%d", loc, pkg.Fset.Position(hidden.Pos()).Line)
}

func isPkgName(obj types.Object) bool {
	_, ok := obj.(*types.PkgName)

	return ok
}
