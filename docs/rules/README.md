# `rules/`

Mandatory rules for `agent_go_souschef`. Each file describes one decision area; together they
form the engineering contract for any new or modified Go code.

The rules are designed to be enforceable. Where possible, each rule cites the linter from
`.golangci.yml` that automates it. Where no linter exists, the rule is enforced via code review.

## Lint baseline

```sh
golangci-lint run ./...
# or, via Taskfile:
task tools:lint
```

## Index

| Rule | Topic |
|---|---|
| [`effective-go.md`](effective-go.md) | The Effective Go subset we enforce |
| [`naming.md`](naming.md) | Go naming for packages, types, files |
| [`imports.md`](imports.md) | Three-group import layout + aliasing |
| [`errors.md`](errors.md) | Sentinel errors, `%w` wrapping |
| [`logging.md`](logging.md) | `log/slog` structured logging |
| [`commits.md`](commits.md) | Conventional commits |
| [`startup.md`](startup.md) | Zero side effects in `init()` |
| [`concurrency.md`](concurrency.md) | Goroutine ownership, `context` |
| [`types.md`](types.md) | Named structs and small interfaces |
| [`code-placement.md`](code-placement.md) | `cmd/` vs `internal/` vs `pkg/` vs `fixtures/` |
| [`testing.md`](testing.md) | Colocated tests, fixtures convention |
| [`dependency-injection.md`](dependency-injection.md) | Constructor injection, no DI container |
| [`interop.md`](interop.md) | Cross-package calls via consumer-defined interfaces |
