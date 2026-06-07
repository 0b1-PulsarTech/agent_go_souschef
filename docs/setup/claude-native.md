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

## 2. Bootstrap the index for your project

From the root of the Go project you want to expose:

```sh
cd /path/to/your/go/project
agent_go-souschef sync
```

This scans the workspace and writes the index to `./.repo-context/index.db`.
Add `.repo-context/` to `.gitignore`.

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
> index file at `./.repo-context/index.db` is found automatically. Subsequent
> sessions just call `souschef_sync` over MCP to refresh it.

## 4. (Optional) Install the PreToolUse hook

If you want Claude to consult the index before doing `Read`/`Grep` on indexed
Go files, install the hook too:

```sh
agent_go-souschef hook install --claude
```

This writes a `PreToolUse` entry to `~/.claude/settings.json`. Idempotent —
re-running won't duplicate it.

## Troubleshooting

| Symptom | Fix |
|---|---|
| Tools don't appear after restart | Check `.claude/mcp.json` is valid JSON; run `agent_go-souschef mcp` manually to see startup errors. |
| `souschef_query` returns "no symbols" | Run `souschef_sync` first (or `agent_go-souschef sync` in a terminal). |
| Permission denied on `.repo-context/index.db` | Delete the directory and resync; SQLite needs write access. |
| Wrong project indexed | Ensure `cwd` in `.claude/mcp.json` is `${workspaceFolder}` and Claude opened the right folder. |

## See also

- [`claude-docker.md`](claude-docker.md) — same setup, containerised.
- [`../patterns/mcp-server.md`](../patterns/mcp-server.md) — what the MCP tools do.
