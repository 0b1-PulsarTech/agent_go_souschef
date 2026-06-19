package repocontext

import "context"

func (svc Service) lookupSymbol(ctx context.Context, symbol Symbol) (QueryHit, error) {
	calls, err := svc.store.Calls(ctx, symbol.ID)
	if err != nil {
		return QueryHit{}, err
	}
	callers, err := svc.store.Callers(ctx, symbol.ID)
	if err != nil {
		return QueryHit{}, err
	}
	impls, err := svc.store.Implementations(ctx, symbol.ID)
	if err != nil {
		return QueryHit{}, err
	}
	usedBy, err := svc.store.References(ctx, symbol.ID)
	if err != nil {
		return QueryHit{}, err
	}
	usesTypes, err := svc.store.TypeUses(ctx, symbol.ID)
	if err != nil {
		return QueryHit{}, err
	}
	methods, err := svc.store.Methods(ctx, symbol.ID)
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
