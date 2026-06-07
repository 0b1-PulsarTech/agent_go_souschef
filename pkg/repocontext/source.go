package repocontext

import (
	"context"
	"fmt"
	"strings"
)

func (svc *Service) Source(ctx context.Context, query string) (string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return "", fmt.Errorf("query is required")
	}
	return svc.indexer.Source(ctx, query)
}
