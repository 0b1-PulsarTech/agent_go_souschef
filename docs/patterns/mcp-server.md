# Pattern — MCP server

`agent_go-souschef mcp` starts a
[Model Context Protocol](https://github.com/modelcontextprotocol/go-sdk)
stdio server. Tools are registered through a thin internal wrapper,
`mcpkit`, that mirrors `fuego.Post` / `fiber.Get` style registration.

## Two-layer split

```
internal/integrations/mcpkit/   ← the wrapper (Server, generic Tool[In,Out], Run)
pkg/repocontext/mcptools.go     ← the handlers (one RegisterMCP call, four tools)
```

The wrapper is internal so the rest of the codebase never touches the raw SDK.
The handlers live in `pkg/` so they're part of the public API surface —
mirroring amigonimo's `pkg/web/handlers/*ctrl/router.go` consuming
`pkg/web.RouterContract`.

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

## `pkg/repocontext/mcptools.go`

```go
func RegisterMCP(s *mcpkit.Server, svc *Service) {
    mcpkit.Tool(s, "souschef_sync", "Build or refresh the symbol index.",
        func(ctx context.Context, _ SyncIn) (Result, error) {
            text, err := svc.Sync(ctx); return Result{Text: text}, err
        })
    mcpkit.Tool(s, "souschef_query", "Look up a symbol and return callers/callees.",
        func(ctx context.Context, in QueryIn) (Result, error) {
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
}
```

Four distinct tools on the same MCP server. The LLM sees four catalog entries
with their own descriptions and JSON schemas — better tool selection than
hiding behind a `command: ...` discriminator on a single tool.

## The four tools

| Tool | Input | What it returns |
|---|---|---|
| `souschef_sync` | `{}` | "Synced." (or an error) — call once per session before querying. |
| `souschef_query` | `{query, expand?}` | Compact symbol + callers/callees. `expand: true` for transitive. |
| `souschef_source` | `{query}` | File path + source snippet for the named symbol. |
| `souschef_changed` | `{scope?}` | Modified-files list, optionally filtered by path. |

## Adding a new tool

1. Add the method to `*repocontext.Service` (`pkg/repocontext/<op>.go`).
2. Add Input/Output structs near the top of `pkg/repocontext/mcptools.go`.
3. Add a single `mcpkit.Tool(s, "souschef_<name>", desc, handler)` call.
4. Add a `t.Run` smoke to `mcptools_test.go`.

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
