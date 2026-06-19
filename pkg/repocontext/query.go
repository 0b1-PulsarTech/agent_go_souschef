package repocontext

import (
	"context"
	"fmt"
	"strings"
)

func (svc Service) Query(ctx context.Context, query string, expanded bool) (string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return "", fmt.Errorf("query is required")
	}
	hit, err := svc.lookup(ctx, query)
	if err != nil {
		return "", err
	}
	return renderCompact(hit, expanded), nil
}

func (svc Service) lookup(ctx context.Context, query string) (QueryHit, error) {
	symbols, err := svc.store.Lookup(ctx, query)
	if err != nil {
		return QueryHit{}, err
	}
	if len(symbols) == 0 {
		return QueryHit{}, nil
	}
	return svc.lookupSymbol(ctx, symbols[0])
}
