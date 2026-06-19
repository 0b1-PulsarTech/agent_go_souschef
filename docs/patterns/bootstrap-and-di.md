# Pattern — Bootstrap and DI

How `agent_go-souschef` wires itself from `main()` down to a running MCP server,
using `github.com/wrapped-owls/goremy-di/remy` for constructor injection.

## Goal

`cmd/agent_go-souschef/main.go` stays thin. It only:

1. Parses `os.Args` for the subcommand.
2. Builds a `remy.Injector` and calls `bootstrap.DoInjections`.
3. Hands the injector to `bootstrap.RunMCP` or `bootstrap.RunSync`.

Every concrete construction (open SQLite, build the goscan indexer, open the
git repo, build the `Service`) lives in `internal/bootstrap/bootstrap.go`.

## File layout

```
internal/bootstrap/
├── bootstrap.go        # Config + DoInjections — the wiring
├── bootstrap_test.go   # asserts every registered type resolves
├── runner.go           # RunMCP + RunSync — the per-subcommand drivers
├── runner_test.go
└── testenv_test.go     # sample-workspace helpers for tests
```

## `bootstrap.go` — the wiring

```go
type Config struct {
    Root    string // workspace directory; defaults to cwd in main.go
    Version string // reported to MCP clients
}

func DoInjections(inj remy.Injector, cfg Config) {
    remy.RegisterInstance(inj, cfg)

    remy.RegisterConstructorErr(inj, remy.Singleton[*reposqlite.Store],
        func() (*reposqlite.Store, error) {
            // Throwaway cache under $TMPDIR keyed by workspace path — never
            // inside the project being indexed (see indexPath in bootstrap.go).
            return reposqlite.Open(indexPath(cfg.Root))
        })

    remy.RegisterConstructor(inj, remy.Singleton[*goscan.Indexer],
        func() *goscan.Indexer { return goscan.New(cfg.Root) })
    remy.RegisterConstructor(inj, remy.Singleton[*gitprobe.Probe],
        func() *gitprobe.Probe { return gitprobe.New(cfg.Root) })

    remy.RegisterConstructorArgs3(inj, remy.Singleton[*repocontext.Service],
        repocontext.New)
}
```

Each layer's registrations stay inline. Promote a `remy.Module` only when a
layer grows past ~5 registrations — until then the flat file is clearer.

### Why `DuckTypeElements: true`?

`repocontext.New(indexer LanguageIndexer, store SymbolStore, changes ChangeReporter)`
asks for interface types, but we register concrete types (`*goscan.Indexer`
etc.). DuckTyping lets remy match the registered concrete to the requested
interface structurally — no extra factory closures.

## `runner.go` — the dispatchers

```go
func RunMCP(ctx context.Context, inj remy.Injector, cfg Config) error {
    svc, err := remy.Get[*repocontext.Service](inj)
    if err != nil {
        return fmt.Errorf("resolve service: %w", err)
    }
    server := mcpkit.New("agent_go_souschef", cfg.Version)
    repocontext.RegisterMCP(server, svc)
    return server.Run(ctx)
}

func RunSync(ctx context.Context, inj remy.Injector, out io.Writer) error {
    svc, err := remy.Get[*repocontext.Service](inj)
    if err != nil {
        return fmt.Errorf("resolve service: %w", err)
    }
    summary, err := svc.Sync(ctx)
    if err != nil {
        return fmt.Errorf("sync: %w", err)
    }
    _, err = fmt.Fprintln(out, summary)
    return err
}
```

`main.go` then becomes a two-screen file: parse args, dispatch.

## Checklist when adding a new subcommand

1. Add the method to `*repocontext.Service` (one file in `pkg/repocontext/`).
2. Add `mcpkit.Tool(s, "souschef_<name>", desc, handler)` to
   `pkg/repocontext/mcptools.go`.
3. If the subcommand needs CLI exposure beyond MCP, add a `Run<Name>` to
   `internal/bootstrap/runner.go` and a `case` in `cmd/agent_go-souschef/main.go`.
4. If you introduce a new collaborator, register it in `DoInjections`.

## See also

- [`cli-skeleton.md`](cli-skeleton.md) — the full `main.go` shape.
- [`mcp-server.md`](mcp-server.md) — the mcpkit wrapper and tool registration.
- [`../rules/dependency-injection.md`](../rules/dependency-injection.md) — the underlying rule.
