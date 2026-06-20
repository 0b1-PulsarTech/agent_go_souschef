# agent_go-souschef

Semantic symbol index for Go repositories, served to LLMs over the
Model Context Protocol. Lets Claude Code (or any MCP host) answer "what calls
this?", "where is X defined?", or "what changed?" in 30–200 tokens instead of
re-reading whole files.

## What you get

Five MCP tools on one stdio server:

| Tool | Purpose |
|---|---|
| `souschef_sync` | Build / refresh the index for the current workspace. |
| `souschef_query` | Look up a symbol; returns direct callers + callees. `expand=true` for transitive. |
| `souschef_source` | Return the file + source snippet for a named symbol. |
| `souschef_changed` | List modified files (optionally filtered by scope). |
| `souschef_shadows` | List identifiers that shadow a builtin, imported package, package symbol, or outer variable (optionally filtered by scope). |

The index is a throwaway cache kept under the OS temp dir
(`$TMPDIR/agent_go_souschef/<workspace-hash>/index.db`), so it never touches
the project you are indexing — nothing to add to `.gitignore`. The `mcp`
server builds it automatically on startup, so a manual `sync` is optional.

---

## Install (use in any Go project)

### Option A — native binary (recommended)

```sh
# 1. Install once
go install github.com/0b1-PulsarTech/agent_go_souschef/cmd/agent_go-souschef@latest

# Make sure GOBIN is on PATH:
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.bashrc  # or ~/.zshrc

# 2. Verify
agent_go-souschef --help

# 3. (Optional) pre-build the index — the mcp server also does this on startup
cd /path/to/your/go/project
agent_go-souschef sync          # builds the index under $TMPDIR
agent_go-souschef shadows       # one-shot: report shadowed builtins/imports/variables
```

Wire it into Claude Code — create `.claude/mcp.json` at your project root:

```json
{
  "mcpServers": {
    "souschef": {
      "command": "agent_go-souschef",
      "args": ["mcp"],
      "cwd": "${workspaceFolder}"
    }
  }
}
```

Restart Claude Code. The five tools show up in the catalog.

Full guide: [`docs/setup/claude-native.md`](docs/setup/claude-native.md).

### Option B — Docker (no Go toolchain on host)

```sh
# 1. Build the image (or pull, once we publish one)
docker build -t agent_go_souschef:latest -f build/Dockerfile .

# 2. (Optional) pre-build the index — the mcp server also does this on startup
cd /path/to/your/go/project
docker run --rm -v "$PWD:/workspace" -w /workspace agent_go_souschef:latest sync
```

Wire it into Claude Code — `.claude/mcp.json`:

```json
{
  "mcpServers": {
    "souschef": {
      "command": "docker",
      "args": [
        "run", "--rm", "-i",
        "-v", "${workspaceFolder}:/workspace",
        "-w", "/workspace",
        "agent_go_souschef:latest", "mcp"
      ]
    }
  }
}
```

Full guide (including SELinux/`:Z`, distroless notes):
[`docs/setup/claude-docker.md`](docs/setup/claude-docker.md).

---

## Verify it works

After installing and running `sync`, from the project root:

```sh
# Run the server (Ctrl-C to stop) and send any MCP client at it
agent_go-souschef mcp
```

Or, more directly inside Claude Code, ask:

> "Use souschef_query to find all callers of `HandleRequest`."

You should see the compact symbol/callers result — not a flood of file reads.

---

## Build from source

```sh
git clone https://github.com/0b1-PulsarTech/agent_go_souschef.git
cd agent_go_souschef
go build ./cmd/agent_go-souschef
go test ./...
```

Dev tasks (via `tools/Taskfile.yml`):

```sh
task -d tools run:checks         # golangci-lint
task -d tools run:formatters     # gofumpt + goimports + golines
task -d tools gen:code           # sqlc generate
task -d tools gen:modernize      # gopls modernize -fix
```

---

## How it works

The indexer walks the workspace with `golang.org/x/tools/go/packages`,
extracts every symbol + call-graph edge, and persists it to SQLite via
sqlc-typed queries. It is workspace-aware: a `go.work` monorepo is indexed
across all of its modules, not just the root. `agent_go-souschef mcp` builds
the index on startup and serves it as five MCP tools — the LLM calls them
instead of reading raw files.

Architecture: see [`AGENTS.md`](AGENTS.md) for layout and
[`docs/patterns/`](docs/patterns/) for the bootstrap / DI / MCP wiring patterns.

---

## License

See LICENSE.
