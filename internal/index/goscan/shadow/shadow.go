// Package shadow detects identifiers that hide one visible in an enclosing
// scope — a predeclared builtin, an imported package, a package-level symbol,
// or a variable from an outer function or block.
//
// It is expressed as a golang.org/x/tools/go/analysis Analyzer so detection
// runs under the maintained analysis driver (FileSet wiring, parallel
// per-package execution, diagnostics) instead of a hand-rolled walk, and so the
// same pass is reusable by any analysis host. Unlike x/tools' own
// passes/shadow — which is deliberately conservative (it ignores builtins and
// import names and suppresses shadows whose types differ) — this pass reports
// every kind, because the index exists to surface hidden names before a human
// picks one.
package shadow

import (
	"go/ast"
	"go/token"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/analysis"
)

// Origin classifies what a shadowing identifier hides. The values double as the
// human-facing label rendered in query output.
const (
	OriginBuiltin = "builtin"
	OriginImport  = "import"
	OriginPackage = "package"
	OriginOuter   = "outer"
)

// Finding is one shadowing site. Positions are already resolved (absolute), so
// callers need no FileSet to render them; the index layer maps the filenames to
// repo-relative paths.
type Finding struct {
	Pos        token.Position // the shadowing identifier
	Name       string
	Origin     string         // one of the Origin* constants
	ImportPath string         // set when Origin == OriginImport
	Hidden     token.Position // the hidden declaration; zero for builtin/import
}

// Analyzer reports every shadowing declaration. It carries no facts and
// requires no other pass, so the driver runs it once per package with no
// dependency fan-out. Run returns []Finding, surfaced through Action.Result.
var Analyzer = &analysis.Analyzer{
	Name:       "souschefshadow",
	Doc:        "report identifiers that shadow a builtin, import, package symbol, or outer variable",
	Run:        run,
	ResultType: reflect.TypeFor[[]Finding](),
}

func run(pass *analysis.Pass) (any, error) {
	var findings []Finding

	for ident, obj := range pass.TypesInfo.Defs {
		f, ok := finding(pass, ident, obj)
		if !ok {
			continue
		}

		pass.Report(analysis.Diagnostic{Pos: ident.Pos(), Message: message(f)})
		findings = append(findings, f)
	}

	return findings, nil
}

// finding resolves whether ident's object hides an identifier visible in an
// enclosing scope. It walks the type-checker's lexical scopes: a declaration's
// own scope (obj.Parent) sits inside the scope that must be searched for an
// outer binding (obj.Parent().Parent()).
func finding(pass *analysis.Pass, ident *ast.Ident, obj types.Object) (Finding, bool) {
	if obj == nil || obj.Name() == "_" {
		return Finding{}, false
	}

	inner := obj.Parent()
	if inner == nil { // struct fields and methods are not lexically scoped.
		return Finding{}, false
	}

	enclosing := inner.Parent()
	if enclosing == nil {
		return Finding{}, false
	}

	scope, hidden := enclosing.LookupParent(obj.Name(), obj.Pos())
	if hidden == nil || hidden == obj {
		return Finding{}, false
	}

	return classify(pass, ident, obj, scope, hidden), true
}
