# Logging

**`log/slog` only.** No `fmt.Println`, no `log.Printf`, no third‑party loggers. The `forbidigo`
rule blocks `fmt.Print*` outside `*_test.go`, and `sloglint` enforces the call shape below.

## Call shape

- Use the package‑level convenience funcs (`slog.Info`, `slog.Error`, ...) only at the very top
  of `main()` and inside boot helpers.
- Inside libraries and request handlers, take a `*slog.Logger` from constructor or context.
- Always use **typed attributes**, never positional `Sprintf`:

```go
slog.Info("secret friend created",
    slog.String("secret_friend_id", sf.ID.String()),
    slog.String("owner_id", sf.OwnerID.String()),
)

if err != nil {
    slog.Error("create secret friend failed",
        slog.String("error", err.Error()),
        slog.String("owner_id", uc.associatedUser.ID.String()),
    )
}
```

`sloglint` (with `attr-only`) rejects `slog.Info("msg", "key", value)` — always use
`slog.String`, `slog.Int`, `slog.Any`, `slog.Group`, etc.

## Levels

| Level | Meaning |
|---|---|
| `Debug` | Verbose tracing, off by default in production. |
| `Info` | Normal operation, business‑level events (request received, message sent). |
| `Warn` | Unexpected but recoverable (retryable downstream failure). |
| `Error` | Action failed; operator attention may be needed. |

Reserve `Error` for things you would page on. Use `Warn` for transient failures that the system
already retries.

## Context

- A logger may live on the `context.Context` (`fuegoport`/`webport` populate one with
  `request_id`, `user_id`, `route`).
- Inside handlers and use cases, derive the logger from context rather than from a global, so
  request‑scoped attributes appear automatically.

## PII and secrets

- Never log secrets, passwords, tokens, or raw message bodies that may contain PII (phone numbers
  belonging to end customers, document numbers, addresses, ...).
- IDs (`user_id`, `contact_id`, `message_id`, `tenant_id`) are fine; raw payloads are not.
- When in doubt, log a redacted summary (`slog.Int("body_size", len(payload))`).

## Format

- Production runs JSON output (`slog.NewJSONHandler`); local dev may use the text handler.
- Log message strings are short, lowercase, no trailing punctuation: `"message dispatched"`,
  not `"Message dispatched."`.
- Keys are `snake_case`, lowercase. Stable keys across services so dashboards line up.

## Don't log + return

When a function returns an error, **do not log it** unless this is the top of the call chain. The
caller will log once, with the wrapped context — see [`errors.md`](errors.md).
