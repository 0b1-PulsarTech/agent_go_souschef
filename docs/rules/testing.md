# Testing

## Layout

- **Tests live in the same package as the code under test.** Use `package <name>`, not `package <name>_test`. If a test needs an unexported helper it can use it directly without exporting.
- Tests are colocated: `foo_test.go` next to `foo.go` in the same directory.
- Fixture source files live in `fixtures/sample/` (own `go.mod`, module `sample`). Per-package fixtures go in `internal/<pkg>/fixtures/`.

## Unit tests — no I/O

Unit tests must not touch external systems. No DB, no network, no disk, no real time sleeps.
If the code under test needs one of those, depend on it through an interface and pass a fake.

- ❌ Real DB in a unit test — use `t.TempDir()` + `reposqlite.Open` for SQLite; it creates a real but ephemeral file.
- ❌ HTTP client hitting the network (use `httptest`)
- ❌ Reading from `/etc`, `~`, or any path you didn't create in the test
- ❌ Goroutines that outlive the test (`t.Cleanup` to stop them)

## Fixture convention

`goscan` tests (and anything that needs real Go source to parse) load from `fixtures/sample/`:

```go
idx := goscan.New(filepath.Join("..", "..", "fixtures", "sample"))
snap, err := idx.Build(context.Background())
```

The fixture module (`module sample`) keeps fixture source out of the main module build without
needing `//go:build ignore` tags. Never add build tags to files under `fixtures/`.

## Table-driven

Default to table-driven tests:

```go
func TestParseKind(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name  string
        input string
        want  string
    }{
        {"func", "func", "func"},
        {"unknown", "xyz", ""},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            if got := parseKind(tt.input); got != tt.want {
                t.Fatalf("got %q, want %q", got, tt.want)
            }
        })
    }
}
```

- `t.Parallel()` at both the outer and inner level.
- Go ≥ 1.22 scopes loop variables per-iteration; the `tt := tt` capture is not needed.

## Helpers

- Helpers that fail the test call `t.Helper()` first and accept `t testing.TB`.
- Set up via constructor injection. No package-level `var db = ...` in tests; use `t.Cleanup`.

## Mocks

Prefer real implementations (SQLite in-memory, `httptest.Server`) over mocks.
When a mock is unavoidable, write it by hand if the interface has ≤ 2 methods; use `mockgen`
otherwise. Output as `_mock_test.go` colocated with the test, same package.

## Coverage

Run `go test -cover ./...` locally; `task main:test` does this in CI.
No enforced coverage floor — the bar is "every business branch is exercised."

## Forbidden

- `time.Sleep` in tests. Use channels, callbacks, or `context.WithTimeout`.
- Sharing state between tests via package globals.
- `os.Exit` from a test.
- `//go:build ignore` on fixture files — use a separate `go.mod` instead.
