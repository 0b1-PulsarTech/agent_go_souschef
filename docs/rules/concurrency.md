# Concurrency

## Ownership

Every goroutine has **one owner** that knows when it stops. Either:

- The owner blocks on the goroutine via `sync.WaitGroup`, `errgroup.Group`, or a result channel,
  **or**
- The goroutine is supervised by a long‑lived loop (e.g. a worker pool) that owns its lifecycle.

A "go and forget" call is a leak. The `noctx` linter catches the most common case (HTTP requests
without a context) but not all.

## `context.Context`

- The first parameter of every function that does I/O or that may block is `ctx context.Context`.
- Never store a `Context` on a struct field. Pass it through.
- Never pass `context.Background()` from inside a request handler — propagate the inbound context.
- Use `context.WithCancel` / `context.WithTimeout` and **always** `defer cancel()`.

## Fan‑out with `errgroup`

```go
import "golang.org/x/sync/errgroup"

g, gctx := errgroup.WithContext(ctx)
for _, target := range targets {
    target := target
    g.Go(func() error {
        return dispatch(gctx, target)
    })
}
if err := g.Wait(); err != nil {
    return fmt.Errorf("dispatch fan-out: %w", err)
}
```

`errgroup` cancels the shared context on the first error and waits for the rest to settle.

## Channels

- Prefer `errgroup` / `WaitGroup` over hand‑rolled channel coordination.
- When you do use a channel, document who **closes** it. The sender closes; the receiver never
  closes a channel it doesn't own.
- Buffered channels are an optimisation, not a synchronisation tool. Default to unbuffered.

## Long‑running workers

For consumer loops (Kafka, webhook poller, gRPC stream):

```go
func (w *Worker) Run(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case msg := <-w.in:
            if err := w.handle(ctx, msg); err != nil {
                slog.Error("handle failed", slog.String("error", err.Error()))
                // decide: continue, retry, or return
            }
        }
    }
}
```

The owning goroutine `Run`s; the parent supervises restarts.

## Mutexes

- Embed `sync.Mutex` as an unexported field, never as the first exported field.
- Use `sync.RWMutex` only when reads vastly outnumber writes — the write path becomes more
  expensive.
- Lock the smallest scope that preserves the invariant. Do not call back into user‑provided
  callbacks while holding a lock.

## Forbidden

- `time.Sleep` to wait for an event in production code. Use channels, `sync.Cond`, or context
  deadlines.
- Spawning goroutines inside a request handler without a `WaitGroup` / `errgroup`.
- Reading from a channel that may be closed without using the two‑value receive form
  (`v, ok := <-ch`).
- Using `select{}` to "block forever" instead of returning from `main`.

## Race detector

The CI test step runs with `-race`. Local runs:

```sh
go test -race ./...
```

Any data race is a bug; never silence the detector.
