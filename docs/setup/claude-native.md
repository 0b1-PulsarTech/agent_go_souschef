# Setup — Claude Code (native binary)

The fastest way to use `agent_go-souschef` from Claude Code: build the binary
once, drop a `.claude/mcp.json` into your project.

## 1. Install the binary

```sh
go install github.com/0b1-PulsarTech/agent_go_souschef/cmd/agent_go-souschef@latest
```

Make sure `$(go env GOBIN)` (or `$(go env GOPATH)/bin`) is on your `PATH`:

```sh
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.bashrc  # or ~/.zshrc
```

Verify:

```sh
agent_go-souschef --help
```

## 2. (Optional) Bootstrap the index for your project

The `mcp` server builds the index automatically on startup, so this step is
optional. Run it only if you want the index ready before the first MCP call:

```sh
cd /path/to/your/go/project
agent_go-souschef sync
```

This scans the workspace — including every module of a `go.work` monorepo —
and writes the index to a throwaway location under the OS temp dir
(`$TMPDIR/agent_go_souschef/<workspace-hash>/index.db`). Nothing is written
into the project, so there is nothing to `.gitignore`.

## 3. Wire it into Claude Code

Create `.claude/mcp.json` at your project root:

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

Restart Claude Code — the four tools (`souschef_sync`, `souschef_query`,
`souschef_source`, `souschef_changed`) appear in the tool catalog.

> **Note**: Claude spawns the process **in the workspace folder**, so the
> index location is derived from that path and reused across sessions. The
> server runs an initial sync on startup; `souschef_sync` refreshes it on
> demand mid-session.

## Troubleshooting

| Symptom | Fix |
|---|---|
| Tools don't appear after restart | Check `.claude/mcp.json` is valid JSON; run `agent_go-souschef mcp` manually to see startup errors. |
| `souschef_query` returns "no symbols" | The startup sync may have failed — call `souschef_sync` again and check the server's stderr log. |
| Index seems stale | Call `souschef_sync` to rebuild; the temp index is keyed by workspace path and reused across runs. |
| Wrong project indexed | Ensure `cwd` in `.claude/mcp.json` is `${workspaceFolder}` and Claude opened the right folder. |

## See also

- [`claude-docker.md`](claude-docker.md) — same setup, containerised.
- [`../patterns/mcp-server.md`](../patterns/mcp-server.md) — what the MCP tools do.
