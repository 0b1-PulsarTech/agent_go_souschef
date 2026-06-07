package goscan

import (
	"go/types"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan/graph"
)

func (b *snapshotBuilder) addImplementations() {
	structs, interfaces := b.namedTypes()
	for iface, ifaceID := range interfaces {
		contract := iface.Underlying().(*types.Interface).Complete()
		for impl, implID := range structs {
			if types.Implements(impl, contract) ||
				types.Implements(types.NewPointer(impl), contract) {
				b.snapshot.Implementations = append(
					b.snapshot.Implementations,
					graph.ImplementRelation(ifaceID, implID),
				)
			}
		}
	}
}
