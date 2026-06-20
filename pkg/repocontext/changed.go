package repocontext

import (
	"context"
	"fmt"
)

func (svc Service) Changed(ctx context.Context, scope string) (string, error) {
	out, err := svc.changes.Changed(ctx, scope)
	if err != nil {
		return "", fmt.Errorf("changed: %w", err)
	}

	return out, nil
}
