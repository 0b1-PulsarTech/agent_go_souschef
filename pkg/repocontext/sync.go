package repocontext

import (
	"context"
	"fmt"
)

func (svc Service) Sync(ctx context.Context) (string, error) {
	snap, err := svc.indexer.Build(ctx)
	if err != nil {
		return "", fmt.Errorf("build index: %w", err)
	}

	if err = svc.store.Reset(ctx); err != nil {
		return "", fmt.Errorf("reset store: %w", err)
	}

	if err = svc.store.Write(ctx, snap); err != nil {
		return "", fmt.Errorf("write snapshot: %w", err)
	}

	return "Synced.", nil
}
