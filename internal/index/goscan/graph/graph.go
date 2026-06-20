package graph

import (
	"go/ast"
	"go/types"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
	"golang.org/x/tools/go/packages"
)

// CalledObject resolves the types.Object referenced by a call expression.
// Returns nil when the callee cannot be determined statically.
func CalledObject(pkg *packages.Package, call *ast.CallExpr) types.Object {
	if pkg == nil {
		return nil
	}

	if ident, ok := call.Fun.(*ast.Ident); ok {
		return pkg.TypesInfo.Uses[ident]
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	if picked := pkg.TypesInfo.Selections[sel]; picked != nil {
		return picked.Obj()
	}

	return pkg.TypesInfo.Uses[sel.Sel]
}

// ImplementRelation returns a Relation representing an implements edge.
func ImplementRelation(fromID, toID int64) repomodel.Relation {
	return repomodel.Relation{FromID: fromID, ToID: toID, Kind: "implement"}
}

// TypeRefRelation returns a Relation representing a type-reference edge: the
// symbol fromID uses the type toID as a field, parameter, result, local, or
// embedded type.
func TypeRefRelation(fromID, toID int64) repomodel.Relation {
	return repomodel.Relation{FromID: fromID, ToID: toID, Kind: "ref"}
}

// ReferencedType returns the named type an identifier refers to, or nil when
// the identifier is not a use of a *types.TypeName (e.g. a value, builtin, or
// package qualifier).
func ReferencedType(pkg *packages.Package, ident *ast.Ident) *types.TypeName {
	if pkg == nil {
		return nil
	}

	name, _ := pkg.TypesInfo.Uses[ident].(*types.TypeName)

	return name
}

// FullName returns the fully-qualified name of a types.Object as "pkg.Name".
// Returns an empty string for nil or universe-scoped objects.
func FullName(obj types.Object) string {
	if obj == nil || obj.Pkg() == nil {
		return ""
	}

	return obj.Pkg().Path() + "." + obj.Name()
}
