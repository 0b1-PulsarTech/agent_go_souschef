package goscan

import "go/types"

func (b *snapshotBuilder) namedTypes() (map[*types.Named]int64, map[*types.Named]int64) {
	structs := map[*types.Named]int64{}
	interfaces := map[*types.Named]int64{}
	for obj, id := range b.ids {
		typeName, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}
		named, ok := typeName.Type().(*types.Named)
		if !ok {
			continue
		}
		switch named.Underlying().(type) {
		case *types.Struct:
			structs[named] = id
		case *types.Interface:
			interfaces[named] = id
		}
	}
	return structs, interfaces
}
