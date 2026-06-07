# Startup — no side effects on import

Importing a package must **only define** symbols. It must not open connections, start goroutines,
read files, register global handlers, or schedule timers. All of that happens inside `main()` and
the `bootstrap` helpers it calls.

## Why

- Imports stay fast and predictable.
- Unit tests can import any package without paying for I/O.
- `go vet`, `gopls`, `golangci-lint`, `staticcheck` all import every package they analyse — they
  must not trigger network calls.
- Boot order becomes explicit and traceable.

## The rule

```go
// ❌ Bad — connection at package load
var DB = openDatabase()

// ❌ Bad — background goroutine at import
func init() {
    go pollWebhooks()
}

// ❌ Bad — env read at import
var apiKey = os.Getenv("API_KEY")
```

```go
// ✅ Good — define only
type Service struct {
    db *sql.DB
}

func New(db *sql.DB) *Service {
    return &Service{db: db}
}
```

Wiring happens in `cmd/<app>/main.go`:

```go
func main() {
    conf, err := bootstrap.Config()
    if err != nil { panic(err) }

    db := bootstrap.OpenDatabase(conf.Database)
    defer db.Close()

    if os.Getenv("RUN_MIGRATIONS") == "true" {
        if err := bootstrap.RunMigrations(conf.Database, migrations.VersionedMigrationsFS()); err != nil {
            slog.Error("migration failed", slog.String("error", err.Error()))
            os.Exit(1)
        }
    }

    inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
    remy.RegisterInstance(inj, db)
    bootstrap.DoInjections(inj, conf)

    server, err := bootstrap.NewWebServer(conf, jwtPublicKey, inj)
    if err != nil { panic(err) }
    defer server.Close()

    if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
        panic(err)
    }
}
```

## `init()` is almost always wrong

`init()` is acceptable in exactly two cases:

1. **Driver registration** in a `_` import inside `cmd/<app>/main.go` (e.g. `_
   "github.com/go-sql-driver/mysql"`). Document each blank import with a one‑line comment.
2. **Compile‑time invariant checks** that have no I/O:

   ```go
   var _ Repository = (*MySQLRepo)(nil) // interface assertion
   ```

Anything else — connections, env reads, goroutines, registrations — goes in `main` or in a
function the wiring layer calls explicitly.

## Test collection imports too

`go test ./...` imports every package under the cwd. A package with import‑time I/O makes the
test collector pay that cost (and may break CI on networks that can't reach the resource). Yet
another reason to keep imports pure.

## Where structured logging is configured

In `cmd/<app>/main.go`, **before** any other call:

```go
slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: conf.LogLevel(),
})))
```

Libraries never call `slog.SetDefault`.

## See also

- [`config.md`](config.md) — how configuration is loaded once, then handed off to the injector.
- [`dependency-injection.md`](dependency-injection.md) — how `bootstrap.DoInjections(inj, conf)`
  wires every collaborator.
