package goscan

import (
	"fmt"
	"go/token"
	"go/types"
	"path/filepath"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
	"golang.org/x/tools/go/packages"
)

type snapshotBuilder struct {
	nextID   int64
	root     string
	ids      map[types.Object]int64
	names    map[string]int64
	pkgs     []*packages.Package
	snapshot repomodel.Snapshot
}

func newBuilder(root string) *snapshotBuilder {
	return &snapshotBuilder{root: root, ids: map[types.Object]int64{}, names: map[string]int64{}}
}

// addPackages folds one module's loaded packages into the snapshot. It is
// called once per module so several loads accumulate into a single index.
func (b *snapshotBuilder) addPackages(pkgs []*packages.Package) error {
	b.pkgs = append(b.pkgs, pkgs...)
	for _, pkg := range pkgs {
		if err := b.addPackage(pkg); err != nil {
			return err
		}
	}
	return nil
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
