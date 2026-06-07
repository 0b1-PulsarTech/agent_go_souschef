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

// RecvName extracts the receiver type name from a method declaration.
func RecvName(decl *ast.FuncDecl) string {
	if star, ok := decl.Recv.List[0].Type.(*ast.StarExpr); ok {
		return star.X.(*ast.Ident).Name
	}
	return decl.Recv.List[0].Type.(*ast.Ident).Name
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
