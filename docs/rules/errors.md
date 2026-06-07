# Errors

**Errors are values.** Functions return `error` as their last result; `panic` is reserved for
unrecoverable boot‑time failures inside `cmd/<app>/main.go`.

## The wrap rule — `wrapcheck`

Whenever you forward an error across a package boundary, wrap it with context:

```go
sf, err := uc.repo.GetSecretFriendByID(id)
if err != nil {
    return entities.SecretFriend{}, fmt.Errorf("get secret friend %s: %w", id, err)
}
```

- Use `%w`, never `%v` or `%s`, so `errors.Is` / `errors.As` keep working.
- Lead with the action you tried, not "error while ...". The wrapper supplies the trailing
  `: <cause>` when printed.
- The `wrapcheck` linter rejects naked returns of errors that originated outside the current
  module/package.

## Sentinel errors

Define sentinel errors as exported package vars when callers need to branch on them:

```go
var ErrNotFound = errors.New("not found")

if errors.Is(err, secretfriend.ErrNotFound) { ... }
```

- Use `errors.Is` for equality, `errors.As` for typed unwrapping.
- Don't compare on `err.Error()` strings — ever.

## Domain errors — `apperr`

For errors that need to surface to HTTP/gRPC with a stable code, use the `apperr` pattern
(`libs/tereckernel/dberrs` for DB errors and a per‑app `internal/domain/apperr` for business
errors):

```go
type Error struct {
    code, publicMessage string
    statusCode          int
    internalErr         error
    metadata            map[string]any
}

func (e *Error) Error() string  { return e.publicMessage }
func (e *Error) Unwrap() error  { return e.internalErr }
func (e *Error) Code() string   { return e.code }
func (e *Error) StatusCode() int { return e.statusCode }
```

Use the helper constructors (`apperr.NotFound`, `apperr.Forbidden`, `apperr.Conflict`,
`apperr.InternalError`, ...) so HTTP middleware can map them uniformly.

```go
return entities.SecretFriend{}, apperr.NotFound(
    "secret_friend_not_found",
    "secret friend not found",
    err,
)
```

The middleware reads `Code()` / `StatusCode()` / `DetailMsg()` / `Metadata()` and renders the
response. Never write the HTTP status from inside the use case.

## Forbidden

- Returning `nil` together with a non‑nil result and **no** error to signal "soft failure". Pick
  one: either `(zero, error)` or use a typed result.
- `panic` in libraries.
- `recover()` to mask bugs. The single legitimate use is at the top of a supervised long‑running
  goroutine (see [`concurrency.md`](concurrency.md)).
- Logging an error then returning the same error — log **or** return, not both. The top of the
  call chain logs.

## Boundaries — `try/except`‑equivalent

`recover` is allowed at exactly two boundaries:

1. The HTTP/gRPC server's panic middleware (already provided by `libs/webport` and `libs/fuegoport`).
2. A goroutine supervisor that restarts a worker loop.

Anywhere else, an unwound panic is a bug.

## Logging an error

Once, at the top of the call chain, with the structured pair `slog.String("error", err.Error())`
or `slog.Any("error", err)`. See [`logging.md`](logging.md).
