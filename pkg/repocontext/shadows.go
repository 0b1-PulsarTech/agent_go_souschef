package repocontext

import (
	"context"
	"fmt"
	"strings"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/source/queryview"
)

// Shadows reports declarations that hide an outer identifier — a builtin, an
// imported package, a package-level symbol, or an enclosing variable. scope is
// an optional case-insensitive path filter, matching souschef_changed.
func (svc Service) Shadows(ctx context.Context, scope string) (string, error) {
	scope = strings.TrimSpace(scope)

	rows, err := svc.store.Shadows(ctx, scope)
	if err != nil {
		return "", fmt.Errorf("shadows: %w", err)
	}

	return queryview.RenderShadows(rows, scope), nil
}
