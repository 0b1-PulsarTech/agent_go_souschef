# Pattern — Hook install

`agent_go-souschef hook` installs a `PreToolUse` hook into AI coding assistants so that file
read requests against indexed files are answered from the semantic index rather than reading
raw source.

## Subcommand surface

```sh
agent_go-souschef hook install --claude          # inject hook into ~/.claude/settings.json
agent_go-souschef hook install --codex           # inject hook into Codex config
agent_go-souschef hook install --claude --codex  # both at once
agent_go-souschef hook run claude                # called by Claude Code's hook machinery (stdin → stdout)
```

## `internal/hooksetup/` layout

```
internal/hooksetup/
├── targets.go   # Target type + ConfigPath()
├── install.go   # Run() dispatcher + install logic
└── claude.go    # RunClaude() PreToolUse handler
```

## `targets.go`

```go
package hooksetup

type Target string

const (
    TargetClaude Target = "claude"
    TargetCodex  Target = "codex"
    TargetCursor Target = "cursor"
    TargetGemini Target = "gemini"
)

func (t Target) ConfigPath() string {
    home, _ := os.UserHomeDir()
    switch t {
    case TargetClaude:
        return filepath.Join(home, ".claude", "settings.json")
    case TargetCodex:
        return filepath.Join(home, ".codex", "config.json")
    // ...
    }
    return ""
}
```

## `install.go` — injecting the hook

Claude Code's `settings.json` accepts a `hooks` object with `PreToolUse` entries:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "agent_go-souschef hook run claude"
          }
        ]
      }
    ]
  }
}
```

`installClaude` reads the file (or starts from `{}`), merges the entry, and writes back
atomically via a temp file + rename. It is idempotent — running install twice does not
duplicate the entry.

## `claude.go` — the PreToolUse handler

Claude Code calls the hook with a JSON payload on stdin and reads a decision from stdout:

```go
type claudePreToolUse struct {
    ToolName  string          `json:"tool_name"`
    ToolInput json.RawMessage `json:"tool_input"`
}

type hookDecision struct {
    Decision string `json:"decision"`
    Reason   string `json:"reason,omitempty"`
}
```

`RunClaude` reads the event, inspects the tool name and input, and returns either:
- `{"decision": "continue"}` — pass through unmodified.
- `{"decision": "block", "reason": "<compact answer from index>"}` — answer from the index,
  preventing Claude from reading the raw file.

Unrecognised tool calls and malformed input both return `continue` so Claude is never blocked
unexpectedly.

## Adding a new target

1. Add a constant in `targets.go` and implement `ConfigPath()`.
2. Add a flag in `install.go`'s `runInstall` and an `installX(path)` function.
3. Add a `case` in `runHook` if the target has a distinct event format.
4. Document the config path in this file.

## See also

- [`cli-skeleton.md`](cli-skeleton.md) — how `hook` is dispatched from `main.go`.
- [`mcp-server.md`](mcp-server.md) — the MCP-based alternative integration.
