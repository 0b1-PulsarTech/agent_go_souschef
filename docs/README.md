# `docs/`

Engineering documentation for `agent_go_souschef`.

> **Rules tell you what. Patterns show you how.**

## How to navigate

1. Start at [`AGENTS.md`](../AGENTS.md) at the repository root.
2. End-user installation: [`../README.md`](../README.md) and [`setup/`](setup/).
3. For mandatory rules (what to do, what not to do), read [`rules/`](rules/).
4. For implementation recipes (how to do it, with code), read [`patterns/`](patterns/).

## Contents

### [`rules/`](rules/) — Rules (imperative, enforceable)

Mandatory standards for new and modified code.

- [`effective-go.md`](rules/effective-go.md) — the Effective Go subset we enforce
- [`naming.md`](rules/naming.md) — Go identifier and package naming
- [`imports.md`](rules/imports.md) — three-group import layout and aliasing
- [`errors.md`](rules/errors.md) — sentinel errors, `%w` wrapping
- [`logging.md`](rules/logging.md) — `log/slog` only, structured key/value
- [`commits.md`](rules/commits.md) — emoji + conventional commit style
- [`startup.md`](rules/startup.md) — zero side effects in `init()`
- [`concurrency.md`](rules/concurrency.md) — goroutine ownership and `context`
- [`types.md`](rules/types.md) — named structs; consumer-side interfaces
- [`code-placement.md`](rules/code-placement.md) — `cmd/` vs `internal/` vs `pkg/`
- [`testing.md`](rules/testing.md) — colocated tests, fixture convention
- [`dependency-injection.md`](rules/dependency-injection.md) — wiring via `remy`
- [`interop.md`](rules/interop.md) — cross-package calls via ports

### [`patterns/`](patterns/) — Patterns (cookbook)

Implementation recipes with copyable templates.

- [`cli-skeleton.md`](patterns/cli-skeleton.md) — `cmd/agent_go-souschef/main.go` shape
- [`bootstrap-and-di.md`](patterns/bootstrap-and-di.md) — `DoInjections` + `RunMCP`/`RunSync`
- [`mcp-server.md`](patterns/mcp-server.md) — `mcpkit.Tool[In,Out]` + `RegisterMCP`
- [`hook-install.md`](patterns/hook-install.md) — `PreToolUse` install + handler
- [`usecase-layout.md`](patterns/usecase-layout.md) — `*Service` method layout
- [`app-skeleton.md`](patterns/app-skeleton.md) — adding a new subcommand end-to-end

### [`setup/`](setup/) — End-user installation

- [`claude-native.md`](setup/claude-native.md) — `go install` + `.claude/mcp.json`
- [`claude-docker.md`](setup/claude-docker.md) — Docker image + volume mount
