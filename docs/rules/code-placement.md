# Code placement

Where each file goes.

> ## ⛔ STOP
>
> New code goes in `cmd/`, `internal/`, `pkg/`, `tools/`, or `test/`.
> Nothing else. `internal/` is not importable outside this module; `pkg/`
> is the stable public surface.

## Top level

| Directory | Role |
|---|---|
| `cmd/agent_go-souschef/` | Binary entry point. Thin `main.go` only — parse args, build injector, dispatch. |
| `internal/` | Implementation packages, grouped by domain (`index/`, `source/`, `integrations/`, `bootstrap/`). |
| `pkg/repocontext/` | Stable public API — `*Service`, `RegisterMCP`, the ports. |
| `test/fixtures/` | Test fixtures. Each subdirectory has its own `go.mod` so it stays out of the main module build. |
| `tools/` | Build/dev tooling — separate Go module with `Taskfile.yml`, `.golangci.yml`, tool directives for sqlc/mockgen/gopls/modernize. |
| `docs/` | This documentation tree. |

## `cmd/agent_go-souschef/`

```
cmd/agent_go-souschef/
├── main.go        # parse os.Args, build injector, dispatch (mcp|sync|hook)
└── main_test.go   # smoke test on the dispatch switch
```

`main.go` holds zero business logic. Three switch cases, every substantial
operation lives in `internal/bootstrap/`.

## `internal/` — grouped by domain

```
internal/
├── bootstrap/                              # remy wiring + per-subcommand runners
│   ├── bootstrap.go    runner.go
│   └── *_test.go
├── index/                                  # symbol index + persistence
│   ├── goscan/{symbols/, graph/, …}        # parse Go source, build snapshot
│   ├── repomodel/{repomodel.go, slices.go} # shared domain types + generic Map/Filter
│   └── reposqlite/{reposqlite.go,_reads.go,_writes.go,sql/,db/}
├── source/                                 # source-level helpers (no shell-out)
│   ├── gitprobe/{gitprobe.go, changed.go}  # go-git/v5 based
│   └── queryview/                          # result rendering
└── integrations/
    ├── mcpkit/{mcpkit.go, tool.go}         # generic MCP wrapper
    └── hooksetup/                          # PreToolUse hook install + handler
```

The three buckets capture the conceptual layers:

- **`index/`** — everything that builds, stores, or describes symbols.
- **`source/`** — everything that reads the workspace directly (git, files).
- **`integrations/`** — everything that talks to an external host (MCP, hooks).

`bootstrap/` sits outside them because it composes all three.

## `pkg/repocontext/`

```
pkg/repocontext/
├── repocontext.go    # Service struct + New constructor
├── contracts.go      # SymbolStore, LanguageIndexer, ChangeReporter
├── types.go          # shared domain aliases
├── sync.go           # one operation per file
├── query.go
├── query_lookup.go
├── source.go
├── changed.go
├── mcptools.go       # RegisterMCP(s, svc) — four mcpkit.Tool registrations
└── *_test.go
```

Everything exported here is part of the stable surface. Add deliberately.

## File-level rules

- **Hard cap: 150 lines per file.** Past that, split by concern.
- **Package-named entry file.** Every package has a `<pkgname>.go` holding
  the main type, its constructor, and any generic helpers. Concern-specific
  operations go in `<pkg>_reads.go`, `<pkg>_writes.go`, etc.
- **No single-file packages.** A package always has ≥2 files plus tests. If a
  candidate package would naturally hold only one file, fold it into a
  neighbour.
- **Every new file gets a `_test.go` sibling.** Same package, same directory.

## Forbidden moves

- Business logic in `cmd/` beyond argument parsing and dispatch.
- Packages in `internal/` importing from `pkg/` — direction is one-way outward.
- Raw SQL strings outside `internal/index/reposqlite/sql/*.sql` — every
  query goes through sqlc.
- Adding a new top-level directory without updating this doc and `AGENTS.md`.
- Shelling out to `git`, `grep`, `find`, or any external CLI when stdlib /
  x/tools / go-git already covers it.
