# `patterns/`

Implementation recipes with working, self-contained Go code.

## Index

| Pattern | What it solves |
|---|---|
| [`cli-skeleton.md`](cli-skeleton.md) | Thin `main.go` dispatching to three subcommands (mcp/sync/hook). |
| [`bootstrap-and-di.md`](bootstrap-and-di.md) | `DoInjections` wiring + `RunMCP`/`RunSync` runners via remy. |
| [`mcp-server.md`](mcp-server.md) | `mcpkit.Tool[In,Out]` wrapper + `RegisterMCP` handler registration. |
| [`hook-install.md`](hook-install.md) | `PreToolUse` install + handler for Claude Code / Codex / Cursor / Gemini. |
| [`usecase-layout.md`](usecase-layout.md) | `*Service` struct + one-method-per-file convention. |
| [`app-skeleton.md`](app-skeleton.md) | Adding a new subcommand or MCP tool end-to-end. |
