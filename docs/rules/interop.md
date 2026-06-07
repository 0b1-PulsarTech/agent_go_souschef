# Interop — cross-package calls

Any call that crosses an `internal/` package boundary must go through an **interface declared
in the consuming package**. The implementing struct lives in its own package and satisfies the
interface structurally.

## Why

- Prevents circular imports between packages.
- Keeps each package independently testable (pass a fake, not the concrete type).
- Makes the dependency graph explicit.
- Allows the implementation to change without touching the consumer.

## The rule

```
Consumer package declares the interface.
Provider package implements it.
main.go wires them together.
```

Example: `pkg/repocontext` consumes the SQLite store but must not import `internal/reposqlite`
for type purposes. It declares `SymbolStore` and the concrete `*reposqlite.Store` satisfies it:

```go
// pkg/repocontext/contracts.go — consumer declares the port
type SymbolStore interface {
    Reset(ctx context.Context) error
    Write(ctx context.Context, snap goscan.Snapshot) error
    Lookup(ctx context.Context, query string) ([]repomodel.Symbol, error)
    Calls(ctx context.Context, symbolID int64) ([]repomodel.Relation, error)
    Callers(ctx context.Context, symbolID int64) ([]repomodel.Relation, error)
}
```

```go
// internal/reposqlite/open.go — provider implements it, no knowledge of the interface
type Store struct{ q *db.Queries }

func (s *Store) Reset(ctx context.Context) error { ... }
func (s *Store) Write(ctx context.Context, snap goscan.Snapshot) error { ... }
// etc.
```

The wiring in `repocontext.New` simply passes `*reposqlite.Store` where `SymbolStore` is
expected — Go's structural typing does the rest.

## Shared types

Pure data types used by multiple packages belong in `internal/repomodel/`. They are imported
directly — the interface rule does not apply to data bags.

## Forbidden

- ❌ A package importing another `internal/` package's concrete type to call methods on it
  directly when an interface would suffice.
- ❌ Interfaces declared in the provider package and imported by the consumer (inverts the
  dependency direction; breaks independent testing).
- ❌ Circular imports between `internal/` packages — restructure or extract to `repomodel/`.

## See also

- [`types.md`](types.md) — consumer-defined interfaces.
- [`dependency-injection.md`](dependency-injection.md) — wiring in `main.go`.
