package goscan

import (
	"go/ast"
	"go/types"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan/symbols"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
	"golang.org/x/tools/go/packages"
)

func (b *snapshotBuilder) addType(
	pkg *packages.Package,
	path string,
	spec *ast.TypeSpec,
	sum *symbols.FileSummary,
) {
	obj := pkg.TypesInfo.Defs[spec.Name]
	if obj == nil {
		return
	}
	id := b.register(obj, repomodel.Symbol{
		Name: spec.Name.Name, Kind: symbols.Kind(spec), Package: pkg.PkgPath, File: path,
		Signature: types.ObjectString(obj, symbols.ShortPkg),
	})
	sum.Exports = append(sum.Exports, spec.Name.Name)
	iface, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		return
	}
	for _, field := range iface.Methods.List {
		for _, name := range field.Names {
			b.snapshot.Methods = append(b.snapshot.Methods, symbols.MethodRecord(id, name.Name))
		}
	}
}
