package repocontext

import "context"

func (svc *Service) Sync(ctx context.Context) (string, error) {
	snap, err := svc.indexer.Build(ctx)
	if err != nil {
		return "", err
	}
	if err = svc.store.Reset(ctx); err != nil {
		return "", err
	}
	if err = svc.store.Write(ctx, snap); err != nil {
		return "", err
	}
	return "Synced.", nil
}
