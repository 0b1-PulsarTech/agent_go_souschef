package reposqlite

import (
	"context"
	"fmt"
	"strings"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/reposqlite/db"
)

// Lookup returns symbols whose name matches the query (exact, contains, or
// "Type.Method" suffix).
func (s *Store) Lookup(ctx context.Context, query string) ([]repomodel.Symbol, error) {
	rows, err := s.q.ListSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("list symbols: %w", err)
	}
	find := strings.ToLower(query)
	matches := repomodel.Filter(rows, func(r db.Symbol) bool {
		name := strings.ToLower(r.Name)
		return name == find || strings.Contains(name, find) || strings.HasSuffix(name, "."+find)
	})
	return repomodel.Map(matches, toSymbol), nil
}

func (s *Store) Calls(ctx context.Context, id int64) ([]string, error) {
	out, err := s.q.GetCallsFrom(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("calls: %w", err)
	}
	return out, nil
}

func (s *Store) Callers(ctx context.Context, id int64) ([]string, error) {
	out, err := s.q.GetCallersOf(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("callers: %w", err)
	}
	return out, nil
}

func (s *Store) Implementations(ctx context.Context, id int64) ([]string, error) {
	out, err := s.q.GetImplementationsOf(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("implementations: %w", err)
	}
	return out, nil
}

func (s *Store) References(ctx context.Context, id int64) ([]string, error) {
	out, err := s.q.GetReferencesOf(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("references: %w", err)
	}
	return out, nil
}

func (s *Store) TypeUses(ctx context.Context, id int64) ([]string, error) {
	out, err := s.q.GetTypeUsesFrom(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("type uses: %w", err)
	}
	return out, nil
}

func (s *Store) Methods(ctx context.Context, id int64) ([]string, error) {
	out, err := s.q.GetMethodsOf(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("methods: %w", err)
	}
	return out, nil
}
