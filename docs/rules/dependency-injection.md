# Dependency injection

All wiring goes through `github.com/wrapped-owls/goremy-di/remy`. Every
collaborator a type depends on **must** be passed via its constructor —
nothing is read from package-level globals, nothing is fetched mid-method.
The injector exists at boot time, registers everything, and disappears into
the hands of constructors.

This is the foundation of every other rule (no globals, no `init()` side
effects, testable code).

## Vocabulary

| Term | What it is |
|---|---|
| `remy.Injector` | The container. Built once in `cmd/agent_go-souschef/main.go`. |
| `remy.Get[T](inj)` | Resolve a `T` from the injector. |
| `remy.RegisterInstance(inj, v)` | Register an already-built value. |
| `remy.RegisterSingleton(inj, binder)` | Lazy singleton; the binder runs once on first resolve. |
| `remy.RegisterConstructor[Err](inj, bind, ctor)` | No-arg constructor (with optional error). |
| `remy.RegisterConstructorArgsN(inj, bind, ctor)` | Constructor that takes N injected args. |
| `remy.Singleton[T]` (default here) | Cached instance, reused across `Get` calls. |
| `remy.Factory[T]` | Marker meaning "every `Get[T]` produces a fresh instance". |

## The boot sequence

`cmd/agent_go-souschef/main.go` does only this:

```go
cfg := bootstrap.Config{Root: ".", Version: version}
inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
bootstrap.DoInjections(inj, cfg)

switch sub {
case "mcp":  os.Exit(bootstrap.RunMCP(ctx, inj, cfg))
case "sync": os.Exit(bootstrap.RunSync(ctx, inj, stdout))
case "hook": os.Exit(hooksetup.Run(args))   // stateless, no injector
}
```

`internal/bootstrap/bootstrap.go` owns `DoInjections` — every collaborator
is registered there in one place. See the
[`bootstrap-and-di` pattern](../patterns/bootstrap-and-di.md) for the full code.

## `DuckTypeElements: true`

`repocontext.New(LanguageIndexer, SymbolStore, ChangeReporter)` asks for
interfaces, but we register concrete pointers (`*goscan.Indexer`,
`*reposqlite.Store`, `*gitprobe.Probe`). DuckTyping lets remy match the
concrete to the requested interface structurally — no extra factory closures
needed.

## Constructors

Constructors **list every dependency** as a parameter:

```go
func New(indexer LanguageIndexer, store SymbolStore, changes ChangeReporter) *Service {
    return &Service{indexer: indexer, store: store, changes: changes}
}
```

Registration:

```go
remy.RegisterConstructorArgs3(inj, remy.Singleton[*repocontext.Service], repocontext.New)
```

The arity (`Args3`) matches the parameter count. Use
`RegisterConstructorArgsNErr` when the constructor returns an error.

## Singleton vs Factory — when

| Lifetime | Choose `remy.Singleton[T]` when... | Choose `remy.Factory[T]` when... |
|---|---|---|
| Process | The instance is stateless or holds long-lived resources (DB pool, MCP server). Almost everything in this project. | — |
| Per-call | The instance needs to be reset per resolve — none of our current types qualify. | Use only if you genuinely need a fresh instance every time. |

In this codebase, everything is a singleton. Promote to `Factory` only when
you can name the per-call state that would leak otherwise.

## `remy.Module`

When a layer has more than ~5 registrations, group them into a `remy.Module`
and call it from `DoInjections`. Below that threshold, keep the registrations
inline — the indirection costs more than it saves.

## Testing

Build a fresh `remy.NewInjector` per test, register fakes, resolve the type
under test:

```go
inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
remy.RegisterInstance(inj, &fakeStore{})
remy.RegisterConstructorArgs1(inj, remy.Singleton[*Service], NewService)
svc, _ := remy.Get[*Service](inj)
```

`internal/bootstrap/bootstrap_test.go` shows the production wiring being
resolved end-to-end against a temp-dir workspace — the canonical "does every
type resolve?" smoke.

## What is forbidden

- `var globalStore SymbolStore` — no package-level singletons.
- Storing `remy.Injector` on a struct so methods can `Get` later. The
  injector is a boot-time tool. After boot, your struct holds typed fields.
- Calling `remy.Get` outside `cmd/`, `bootstrap/`, or constructors.
- Constructors that read env vars, open files, or hit the network. They
  receive their deps and return.
