package goscan

import (
	"context"
	"fmt"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
	"golang.org/x/tools/go/packages"
)

func (idx *Indexer) Build(ctx context.Context) (repomodel.Snapshot, error) {
	cfg := &packages.Config{Context: ctx, Dir: idx.root, Mode: pkgMode()}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return repomodel.Snapshot{}, fmt.Errorf("load packages: %w", err)
	}
	return collect(pkgs, idx.root)
}

func pkgMode() packages.LoadMode {
	return packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
		packages.NeedImports | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo
}
