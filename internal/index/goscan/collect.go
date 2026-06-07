package goscan

import (
	"fmt"
	"go/token"
	"go/types"
	"path/filepath"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
	"golang.org/x/tools/go/packages"
)

func collect(pkgs []*packages.Package, root string) (repomodel.Snapshot, error) {
	builder := snapshotBuilder{ids: map[types.Object]int64{}, names: map[string]int64{}, root: root}
	for _, pkg := range pkgs {
		if err := builder.addPackage(pkg); err != nil {
			return repomodel.Snapshot{}, err
		}
	}
	builder.addImplementations()
	return builder.snapshot, nil
}

type snapshotBuilder struct {
	nextID   int64
	root     string
	ids      map[types.Object]int64
	names    map[string]int64
	snapshot repomodel.Snapshot
}

func (b *snapshotBuilder) addPackage(pkg *packages.Package) error {
	for _, file := range pkg.Syntax {
		if err := b.addFile(pkg, file); err != nil {
			return fmt.Errorf("add file: %w", err)
		}
	}
	return nil
}

func rel(root string, fset *token.FileSet, pos token.Pos) string {
	path := fset.Position(pos).Filename
	result, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return result
}
