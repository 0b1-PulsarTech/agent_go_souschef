# Setup — Claude Code (Docker container)

Run `agent_go-souschef` inside a container so the host doesn't need a Go
toolchain. Source files are exposed via a volume mount; the index file is
written into the same mount so it persists across container restarts.

## 1. Build (or pull) the image

A minimal `build/Dockerfile`:

```dockerfile
# build/Dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/agent_go-souschef ./cmd/agent_go-souschef

FROM gcr.io/distroless/static:nonroot
WORKDIR /workspace
COPY --from=builder /out/agent_go-souschef /usr/local/bin/agent_go-souschef
ENTRYPOINT ["agent_go-souschef"]
```

Build:

```sh
docker build -t agent_go_souschef:latest -f build/Dockerfile .
```

Or, when the image lands on a registry:

```sh
docker pull ghcr.io/0b1-pulsartech/agent_go_souschef:latest
```

## 2. Wire it into Claude Code

`.claude/mcp.json` at your project root:

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

Key flags:

| Flag | Why |
|---|---|
| `--rm` | Remove the container on exit (every MCP session is fresh). |
| `-i` | Keep stdin open — MCP uses stdio for transport. |
| `-v ${workspaceFolder}:/workspace` | Expose the project source so `go/packages` can load it. |
| `-w /workspace` | Make `/workspace` the cwd so `.repo-context/index.db` lands inside the volume. |

## 3. Bootstrap the index

From the host:

```sh
docker run --rm -v "$PWD:/workspace" -w /workspace agent_go_souschef:latest sync
```

This populates `./.repo-context/index.db` inside the project, accessible to
every subsequent containerised `mcp` invocation.

## Caveats

- The container must be able to read the volume as the user the daemon runs
  as. On SELinux/Podman, use `:Z` (`-v ${workspaceFolder}:/workspace:Z`).
- Distroless images don't ship `git`. `gitprobe` uses the
  [go-git](https://github.com/go-git/go-git) library, so the missing binary
  doesn't matter — but the `.git/` directory **must** be in the mounted
  volume for `souschef_changed` to work.
- The hook subcommand needs write access to `~/.claude/settings.json`, which
  isn't in the container's filesystem. Install hooks from the host with the
  native binary instead.

## See also

- [`claude-native.md`](claude-native.md) — non-containerised setup.
- [`../patterns/mcp-server.md`](../patterns/mcp-server.md) — what the MCP tools do.
