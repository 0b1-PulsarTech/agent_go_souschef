# Pattern — CLI skeleton

The canonical shape for `agent_go-souschef`'s entry point. Two subcommands,
no business logic in `main.go`.

## `cmd/agent_go-souschef/main.go`

```go
const version = "0.1.0"

func main() { os.Exit(run(context.Background(), os.Args[1:], os.Stdout, os.Stderr)) }

func run(ctx context.Context, args []string, stdout, stderr *os.File) int {
    if len(args) == 0 {
        printUsage(stderr)
        return 1
    }

    root, err := filepath.Abs(".")
    if err != nil {
        fmt.Fprintln(stderr, err)
        return 1
    }
    cfg := bootstrap.Config{Root: root, Version: version}
    inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
    bootstrap.DoInjections(inj, cfg)

    switch args[0] {
    case "mcp":
        if err := bootstrap.RunMCP(ctx, inj, cfg); err != nil {
            fmt.Fprintln(stderr, err)
            return 1
        }
        return 0
    case "sync":
        if err := bootstrap.RunSync(ctx, inj, stdout); err != nil {
            fmt.Fprintln(stderr, err)
            return 1
        }
        return 0
    default:
        fmt.Fprintf(stderr, "unknown command %q\n", args[0])
        printUsage(stderr)
        return 1
    }
}
```

## The two subcommands

| Command | Stateful | Purpose |
|---|---|---|
| `mcp` | ✅ (index) | Start the stdio MCP server — the primary surface. Builds the index on startup. |
| `sync` | ✅ (index) | Build/refresh the index once. Optional — `mcp` already syncs on startup. |

The data commands (`query`, `source`, `changed`) are deliberately **not** in
the CLI. They are exposed only over MCP — that keeps one surface authoritative
and the binary small.

## Key constraints

- `main.go` is split-free dispatch only. Anything more substantial lives in
  `internal/bootstrap/`.
- `run` takes its writers as parameters so `main_test.go` can drive it
  without touching real stdout/stderr.
- The root is absolutised once so the index path and package loading are
  independent of how the process was launched.

## See also

- [`bootstrap-and-di.md`](bootstrap-and-di.md) — what's inside `DoInjections`.
- [`mcp-server.md`](mcp-server.md) — what the MCP server actually exposes.
- [`../setup/claude-native.md`](../setup/claude-native.md) — wiring the binary into Claude Code.
