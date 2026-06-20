package symbols

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
)

// Kind returns the symbol kind string for a type declaration.
func Kind(spec *ast.TypeSpec) string {
	switch spec.Type.(type) {
	case *ast.InterfaceType:
		return "interface"
	case *ast.StructType:
		return "struct"
	default:
		return "type"
	}
}

// MethodRecord builds a Method entry for an interface method.
func MethodRecord(parentID int64, name string) repomodel.Method {
	return repomodel.Method{ParentID: parentID, Name: name, Signature: name, MemberKind: "method"}
}

// ShortPkg is a types.Qualifier that drops the package path, producing
// compact type signatures without fully-qualified prefixes.
func ShortPkg(*types.Package) string { return "" }

// RecvName extracts the receiver type name from a method declaration,
// unwrapping pointer receivers (*T) and generic ones (T[P], T[P, Q]). It
// returns "" if the receiver shape is unrecognised rather than panicking.
func RecvName(decl *ast.FuncDecl) string {
	if decl.Recv == nil || len(decl.Recv.List) == 0 {
		return ""
	}

	return baseTypeName(decl.Recv.List[0].Type)
}

func baseTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr: // *T
		return baseTypeName(t.X)
	case *ast.IndexExpr: // T[P]
		return baseTypeName(t.X)
	case *ast.IndexListExpr: // T[P, Q]
		return baseTypeName(t.X)
	case *ast.Ident:
		return t.Name
	default:
		return ""
	}
}

// FileSummary accumulates exported symbol names for one source file.
type FileSummary struct {
	Path    string
	Pkg     string
	Exports []string
}

// Text returns a compact human-readable summary of the file.
func (s FileSummary) Text() string {
	return "Package " + s.Pkg + ". Exports: " + strings.Join(s.Exports, ", ")
}
