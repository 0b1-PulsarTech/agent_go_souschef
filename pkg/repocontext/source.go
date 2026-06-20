package repocontext

import (
	"context"
	"fmt"
	"strings"
)

func (svc Service) Source(ctx context.Context, query string) (string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return "", fmt.Errorf("query is required")
	}

	src, err := svc.indexer.Source(ctx, query)
	if err != nil {
		return "", fmt.Errorf("source %q: %w", query, err)
	}

	return src, nil
}
