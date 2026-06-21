package shadow

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// classify labels a confirmed shadow by what it hides. Import names are checked
// before the package scope because an import lives in the file scope (a child of
// the package scope) yet reads more usefully as "import" than "package".
func classify(
	pass *analysis.Pass,
	ident *ast.Ident,
	obj types.Object,
	scope *types.Scope,
	hidden types.Object,
) Finding {
	f := Finding{
		Pos:  pass.Fset.Position(ident.Pos()),
		Name: obj.Name(),
	}

	switch {
	case scope == types.Universe:
		f.Origin = OriginBuiltin
	case isPkgName(hidden):
		f.Origin = OriginImport
		f.ImportPath = hidden.(*types.PkgName).Imported().Path()
	case scope == pass.Pkg.Scope():
		f.Origin = OriginPackage
		f.Hidden = pass.Fset.Position(hidden.Pos())
	default:
		f.Origin = OriginOuter
		f.Hidden = pass.Fset.Position(hidden.Pos())
	}

	return f
}

func isPkgName(obj types.Object) bool {
	_, ok := obj.(*types.PkgName)

	return ok
}

// message renders the diagnostic text. It is intentionally path-free so the
// same string is stable across machines (the structured Finding carries the
// hidden declaration's location for callers that need it).
func message(f Finding) string {
	switch f.Origin {
	case OriginBuiltin:
		return fmt.Sprintf("%q shadows a builtin", f.Name)
	case OriginImport:
		return fmt.Sprintf("%q shadows import %q", f.Name, f.ImportPath)
	case OriginPackage:
		return fmt.Sprintf("%q shadows a package-level declaration", f.Name)
	default:
		return fmt.Sprintf("%q shadows an outer declaration", f.Name)
	}
}
