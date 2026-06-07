# Pattern — Service method layout

`pkg/repocontext.Service` is the single public surface for all business operations. Each
operation lives in its own file. Adding a new operation means adding a new file, not touching
existing ones.

## File layout

```
pkg/repocontext/
├── repocontext.go      # Service struct + New() + SymbolStore interface
├── contracts.go        # all port interfaces
├── types.go            # QueryResult, ChangeSet, and other return types
├── sync.go             # Sync(ctx) — build and persist a fresh snapshot
├── query.go            # Query(ctx, q, expand) — semantic symbol search
├── query_lookup.go     # lookup helpers used by Query
├── source.go           # Source(ctx, q) — return source snippet for a symbol
└── changed.go          # Changed(ctx, scope) — list files changed since last sync
```

## The struct

```go
// pkg/repocontext/repocontext.go
package repocontext

type Service struct {
    root  string
    store SymbolStore
    scan  *goscan.Indexer
    git   *gitprobe.Probe
}

func New(root string) (*Service, error) { ... }
```

## One operation per file

```go
// pkg/repocontext/sync.go
package repocontext

func (s *Service) Sync(ctx context.Context) (string, error) {
    snap, err := s.scan.Build(ctx)
    if err != nil {
        return "", fmt.Errorf("scan: %w", err)
    }
    if err := s.store.Reset(ctx); err != nil {
        return "", fmt.Errorf("reset store: %w", err)
    }
    if err := s.store.Write(ctx, snap); err != nil {
        return "", fmt.Errorf("write snapshot: %w", err)
    }
    return fmt.Sprintf("indexed %d symbols", len(snap.Symbols)), nil
}
```

Each method:
- Takes a `context.Context` as its first argument.
- Returns a compact string result (passed to the LLM or printed) plus an `error`.
- Wraps errors with `fmt.Errorf("%w")` so callers can inspect them.

## Port interfaces

Interfaces consumed by `Service` are declared in `contracts.go`, not in the implementing
package:

```go
// pkg/repocontext/contracts.go
type SymbolStore interface {
    Reset(ctx context.Context) error
    Write(ctx context.Context, snap goscan.Snapshot) error
    Lookup(ctx context.Context, query string) ([]repomodel.Symbol, error)
    Calls(ctx context.Context, symbolID int64) ([]repomodel.Relation, error)
    Callers(ctx context.Context, symbolID int64) ([]repomodel.Relation, error)
    Implementations(ctx context.Context, symbolID int64) ([]repomodel.Relation, error)
    Methods(ctx context.Context, symbolID int64) ([]repomodel.Method, error)
}
```

The concrete `*reposqlite.Store` satisfies this interface structurally.

## Tests

```go
// pkg/repocontext/sync_test.go
func TestSync(t *testing.T) {
    t.Parallel()
    src := filepath.Join("..", "..", "fixtures", "sample")
    svc, err := repocontext.New(src)
    if err != nil {
        t.Fatal(err)
    }
    result, err := svc.Sync(context.Background())
    if err != nil {
        t.Fatalf("sync: %v", err)
    }
    if result == "" {
        t.Fatal("expected result")
    }
}
```

Use real implementations (SQLite + goscan) rather than mocks wherever practical.

## See also

- [`bootstrap-and-di.md`](bootstrap-and-di.md) — wiring `Service` in `main.go`.
- [`../rules/interop.md`](../rules/interop.md) — consumer-defined interfaces.
- [`../rules/errors.md`](../rules/errors.md) — error wrapping.
