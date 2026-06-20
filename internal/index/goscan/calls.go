package goscan

import (
	"go/ast"
	"go/types"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan/graph"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
	"golang.org/x/tools/go/packages"
)

func (b *snapshotBuilder) addCall(
	pkg *packages.Package,
	node ast.Node,
	from types.Object,
	fromID int64,
) bool {
	call, ok := node.(*ast.CallExpr)
	if !ok || node == nil {
		return true
	}

	target := graph.CalledObject(pkg, call)
	if target == nil {
		return true
	}

	toID := b.ids[target]
	if toID == 0 {
		toID = b.names[graph.FullName(target)]
	}

	if toID != 0 && target != from {
		b.snapshot.Calls = append(
			b.snapshot.Calls,
			repomodel.Relation{FromID: fromID, ToID: toID, Kind: "call"},
		)
	}

	return true
}
