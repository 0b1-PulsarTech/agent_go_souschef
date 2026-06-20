package repocontext

import "context"

type LanguageIndexer interface {
	Build(ctx context.Context) (Snapshot, error)
	Source(ctx context.Context, query string) (string, error)
}

type SymbolStore interface {
	Reset(ctx context.Context) error
	Write(ctx context.Context, snap Snapshot) error
	Lookup(ctx context.Context, query string) ([]Symbol, error)
	Calls(ctx context.Context, id int64) ([]string, error)
	Callers(ctx context.Context, id int64) ([]string, error)
	Implementations(ctx context.Context, id int64) ([]string, error)
	References(ctx context.Context, id int64) ([]string, error)
	TypeUses(ctx context.Context, id int64) ([]string, error)
	Methods(ctx context.Context, id int64) ([]string, error)
	Shadows(ctx context.Context, scope string) ([]Shadow, error)
}

type ChangeReporter interface {
	Changed(ctx context.Context, scope string) (string, error)
}
