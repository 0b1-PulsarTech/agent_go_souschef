package goscan

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
	"golang.org/x/tools/go/packages"
)

// Build loads every module under the workspace root and folds their symbols
// and call edges into one snapshot. A module that fails to load is skipped
// with a warning rather than failing the whole index.
func (idx Indexer) Build(ctx context.Context) (repomodel.Snapshot, error) {
	plans, err := loadPlans(idx.root)
	if err != nil {
		return repomodel.Snapshot{}, fmt.Errorf("discover modules: %w", err)
	}

	builder := newBuilder(idx.root)

	for _, plan := range plans {
		cfg := &packages.Config{Context: ctx, Dir: plan.dir, Mode: pkgMode()}

		pkgs, loadErr := packages.Load(cfg, plan.patterns...)
		if loadErr != nil {
			slog.Warn("skip module", "dir", plan.dir, "err", loadErr)

			continue
		}

		if err = builder.addPackages(pkgs); err != nil {
			return repomodel.Snapshot{}, err
		}
	}

	builder.addImplementations()
	builder.addTypeRefs()

	if err = builder.addShadows(); err != nil {
		return repomodel.Snapshot{}, err
	}

	return builder.snapshot, nil
}

// pkgMode is the load mode shared by every extractor. NeedTypesSizes is
// required because the analysis driver behind the shadow pass rejects packages
// whose TypesSizes is unset.
func pkgMode() packages.LoadMode {
	return packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
		packages.NeedImports | packages.NeedSyntax | packages.NeedTypes |
		packages.NeedTypesSizes | packages.NeedTypesInfo
}
