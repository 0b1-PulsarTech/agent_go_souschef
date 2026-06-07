package goscan

import (
	"go/ast"
	"go/types"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan/symbols"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
	"golang.org/x/tools/go/packages"
)

func (b *snapshotBuilder) addFunc(
	pkg *packages.Package,
	path string,
	decl *ast.FuncDecl,
	sum *symbols.FileSummary,
) {
	obj := pkg.TypesInfo.Defs[decl.Name]
	if obj == nil {
		return
	}
	name := decl.Name.Name
	kind := "function"
	if decl.Recv != nil {
		name = symbols.RecvName(decl) + "." + name
		kind = "method"
	}
	id := b.register(obj, repomodel.Symbol{
		Name: name, Kind: kind, Package: pkg.PkgPath, File: path,
		Signature: types.ObjectString(obj, symbols.ShortPkg),
	})
	sum.Exports = append(sum.Exports, name)
	ast.Inspect(decl, func(node ast.Node) bool { return b.addCall(pkg, node, obj, id) })
}
