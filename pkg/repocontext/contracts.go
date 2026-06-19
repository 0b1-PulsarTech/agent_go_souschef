package repocontext

import "context"

type LanguageIndexer interface {
	Build(context.Context) (Snapshot, error)
	Source(context.Context, string) (string, error)
}

type SymbolStore interface {
	Reset(context.Context) error
	Write(context.Context, Snapshot) error
	Lookup(context.Context, string) ([]Symbol, error)
	Calls(context.Context, int64) ([]string, error)
	Callers(context.Context, int64) ([]string, error)
	Implementations(context.Context, int64) ([]string, error)
	References(context.Context, int64) ([]string, error)
	TypeUses(context.Context, int64) ([]string, error)
	Methods(context.Context, int64) ([]string, error)
}

type ChangeReporter interface {
	Changed(context.Context, string) (string, error)
}
