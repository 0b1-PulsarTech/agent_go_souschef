package repocontext

import (
	"context"
	"fmt"
)

func (svc Service) lookupSymbol(ctx context.Context, symbol Symbol) (QueryHit, error) {
	// load runs one relation query and tags failures with the relation name, so
	// the six near-identical calls below stay one-liners and wrap consistently.
	load := func(name string, fn func(context.Context, int64) ([]string, error)) ([]string, error) {
		out, err := fn(ctx, symbol.ID)
		if err != nil {
			return nil, fmt.Errorf("load %s for %q: %w", name, symbol.Name, err)
		}

		return out, nil
	}

	calls, err := load("calls", svc.store.Calls)
	if err != nil {
		return QueryHit{}, err
	}

	callers, err := load("callers", svc.store.Callers)
	if err != nil {
		return QueryHit{}, err
	}

	impls, err := load("implementations", svc.store.Implementations)
	if err != nil {
		return QueryHit{}, err
	}

	usedBy, err := load("references", svc.store.References)
	if err != nil {
		return QueryHit{}, err
	}

	usesTypes, err := load("type uses", svc.store.TypeUses)
	if err != nil {
		return QueryHit{}, err
	}

	methods, err := load("methods", svc.store.Methods)
	if err != nil {
		return QueryHit{}, err
	}

	return QueryHit{
		Symbol:          symbol,
		Calls:           calls,
		Callers:         callers,
		Implementations: impls,
		UsedBy:          usedBy,
		UsesTypes:       usesTypes,
		Methods:         methods,
	}, nil
}
