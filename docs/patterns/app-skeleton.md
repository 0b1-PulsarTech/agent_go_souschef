# Pattern — CLI skeleton

How to extend `agent_go-souschef` with a new subcommand.

## File layout (existing)

```
agent_go_souschef/
├── cmd/agent_go-souschef/
│   └── main.go
├── internal/
│   ├── goscan/
│   │   ├── symbols/
│   │   └── graph/
│   ├── reposqlite/
│   │   ├── sql/
│   │   └── db/             ← sqlc-generated
│   ├── textgrep/
│   ├── gitprobe/
│   ├── queryview/
│   ├── repomodel/
│   └── mcpserver/
├── pkg/repocontext/
├── fixtures/sample/
├── go.mod
├── go.sum
└── sqlc.yaml
```

## Adding a new subcommand

### 1. Stateless subcommand (no index needed)

If the subcommand needs no repository index (e.g. `version`, `help`):

```go
// cmd/agent_go-souschef/main.go
case "version":
    fmt.Println("agent_go-souschef 0.1.0")
    os.Exit(0)
```

No new files required.

### 2. Index-backed subcommand

If the subcommand reads or writes the repo index:

**a)** Add a method to `*repocontext.Service` in `pkg/repocontext/<operation>.go`:

```go
// pkg/repocontext/stats.go
package repocontext

func (s *Service) Stats(ctx context.Context) (string, error) {
    symbols, err := s.store.ListSymbols(ctx)
    if err != nil {
        return "", fmt.Errorf("list: %w", err)
    }
    return fmt.Sprintf("%d symbols indexed", len(symbols)), nil
}
```

**b)** If `Stats` needs a new store query, add it to `internal/reposqlite/sql/queries.sql`
and run `sqlc generate` to regenerate `internal/reposqlite/db/`.

**c)** Wire the subcommand in `main.go`:

```go
case "stats":
    result, err := svc.Stats(context.Background())
    exitOnErr(err)
    fmt.Println(result)
```

### 3. New internal package

Only create a new package in `internal/` when the logic is substantial enough to warrant it:

```
internal/myfeature/
├── myfeature.go
└── myfeature_test.go
```

Import it from `pkg/repocontext/` or, if it's a subcommand handler, from `main.go`.

## Verification

```sh
go build ./cmd/agent_go-souschef/...
go test ./...
task tools:lint
```

## Checklist

- [ ] New method in `pkg/repocontext/` if the subcommand reads the index.
- [ ] New SQL query in `internal/reposqlite/sql/queries.sql` + `sqlc generate` if needed.
- [ ] `case` added in `cmd/agent_go-souschef/main.go`.
- [ ] Test for the new method in `pkg/repocontext/`.
- [ ] Build and tests clean.

## See also

- [`bootstrap-and-di.md`](bootstrap-and-di.md) — wiring and DI pattern.
- [`mcp-server.md`](mcp-server.md) — exposing operations over MCP.
