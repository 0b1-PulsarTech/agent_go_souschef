package repocontext

import "context"

func (svc Service) Changed(ctx context.Context, scope string) (string, error) {
	return svc.changes.Changed(ctx, scope)
}
