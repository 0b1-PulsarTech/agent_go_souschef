# Pattern — MCP server

`agent_go-souschef mcp` starts a
[Model Context Protocol](https://github.com/modelcontextprotocol/go-sdk)
stdio server. Tools are registered through a thin internal wrapper,
`mcpkit`, that mirrors `fuego.Post` / `fiber.Get` style registration.

## Two-layer split

```
internal/integrations/mcpkit/   ← the wrapper (Server, generic Tool[In,Out], Run)
pkg/repocontext/mcpsvc/          ← the handlers (RegisterMCP + IO schema, five tools)
```

The wrapper is internal so the rest of the codebase never touches the raw SDK.
The handlers live in `mcpsvc`, a subpackage of `repocontext`, so the MCP SDK
stays out of the domain package — anything importing `repocontext` as a library
gets the `Service` without the transport.

## `mcpkit/mcpkit.go` + `mcpkit/tool.go`

```go
type Server struct{ impl *mcp.Server }

func New(name, version string) *Server {
    return &Server{impl: mcp.NewServer(&mcp.Implementation{Name: name, Version: version}, nil)}
}

func (s *Server) Run(ctx context.Context) error {
    return s.impl.Run(ctx, &mcp.StdioTransport{})
}

// Generic Tool — one call per handler, no SDK shim by hand.
func Tool[In, Out any](s *Server, name, description string,
    handle func(context.Context, In) (Out, error),
) {
    mcp.AddTool(s.impl, &mcp.Tool{Name: name, Description: description},
        func(ctx context.Context, _ *mcp.CallToolRequest, in In) (*mcp.CallToolResult, Out, error) {
            out, err := handle(ctx, in)
            return nil, out, err
        })
}
```

## `pkg/repocontext/mcpsvc/tools.go`

```go
func RegisterMCP(s *mcpkit.Server, svc repocontext.Service, gate *SyncGate) {
    mcpkit.Tool(s, "souschef_sync", "Build or refresh the symbol index.",
        func(ctx context.Context, _ SyncIn) (Result, error) {
            text, err := gate.Force(ctx); return Result{Text: text}, err   // explicit: always sync
        })
    mcpkit.Tool(s, "souschef_query", "Look up a symbol and return callers/callees.",
        func(ctx context.Context, in QueryIn) (Result, error) {
            gate.Ensure(ctx)                                               // throttled auto-sync
            text, err := svc.Query(ctx, in.Query, in.Expand); return Result{Text: text}, err
        })
    mcpkit.Tool(s, "souschef_source", "Return the source location for a symbol.",
        func(ctx context.Context, in SourceIn) (Result, error) {
            text, err := svc.Source(ctx, in.Query); return Result{Text: text}, err
        })
    mcpkit.Tool(s, "souschef_changed", "List files modified in the workspace.",
        func(ctx context.Context, in ChangedIn) (Result, error) {
            text, err := svc.Changed(ctx, in.Scope); return Result{Text: text}, err
        })
    mcpkit.Tool(s, "souschef_shadows", "Report shadowed builtins/imports/variables.",
        func(ctx context.Context, in ShadowsIn) (Result, error) {
            text, err := svc.Shadows(ctx, in.Scope); return Result{Text: text}, err
        })
}
```

Four distinct tools on the same MCP server. The LLM sees four catalog entries
with their own descriptions and JSON schemas — better tool selection than
hiding behind a `command: ...` discriminator on a single tool.

## The five tools

| Tool | Input | What it returns |
|---|---|---|
| `souschef_sync` | `{}` | "Synced." (or an error) — forces a refresh now, bypassing the throttle. The server also syncs on startup and auto-refreshes before each read tool (see below). |
| `souschef_query` | `{query, expand?}` | Compact symbol + callers/callees. `expand: true` for transitive. |
| `souschef_source` | `{query}` | File path + source snippet for the named symbol. |
| `souschef_changed` | `{scope?}` | Modified-files list, optionally filtered by path. |
| `souschef_shadows` | `{scope?}` | Shadowed identifiers (builtin / import / package symbol / outer variable), grouped by file. |

## Auto-sync gate

`SyncGate` (`mcpsvc/autosync.go`) keeps the index fresh without re-indexing on
every call. Read tools call `gate.Ensure(ctx)` first; it syncs only when at
least `DefaultSyncInterval` (15s) has elapsed since the last sync, so a burst of
calls collapses to one refresh. `souschef_sync` calls `gate.Force(ctx)`, which
always syncs and resets the window. `RunMCP` builds the gate, primes it with one
`Force` at startup, and hands it to `RegisterMCP`.

The gate claims its window under a mutex but runs the sync outside the lock, so
concurrent calls don't pile up. A `time.Now` clock is injected through the
constructor so the throttle is unit-tested without sleeping. A failed auto-sync
is logged, not returned — a stale index still answers, matching the non-fatal
startup-sync policy.

## Adding a new tool

1. Add the method to `repocontext.Service` (`pkg/repocontext/<op>.go`).
2. Add Input/Output structs to `pkg/repocontext/mcpsvc/schema.go`.
3. Add a single `mcpkit.Tool(s, "souschef_<name>", desc, handler)` call in
   `pkg/repocontext/mcpsvc/tools.go` — call `gate.Ensure(ctx)` first if the
   tool reads the index.
4. Add a `t.Run` smoke to `mcpsvc/tools_test.go`.

No router boilerplate, no handler-per-file overhead.

## Why a wrapper at all?

- Type safety: `Tool[In, Out]` rejects malformed signatures at compile time.
- SDK isolation: an SDK upgrade or transport swap is a one-file change.
- Test discoverability: every tool is registered through one function, so
  `mcptools_test.go` can sanity-check registration without booting stdio.

## See also

- [`bootstrap-and-di.md`](bootstrap-and-di.md) — how `RegisterMCP` is wired by `RunMCP`.
- [`../setup/claude-native.md`](../setup/claude-native.md) — Claude Code launch config.
- [`../setup/claude-docker.md`](../setup/claude-docker.md) — Docker container launch config.
