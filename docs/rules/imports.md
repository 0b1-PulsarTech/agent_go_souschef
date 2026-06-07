# Imports

Three import groups, separated by a single blank line, in this order:

1. **Standard library**
2. **External modules** (third‑party + sibling repos like `terectek_commsproto` and
   `wrapped-owls/...`)
3. **This module** (`github.com/0b1-PulsarTech/terectek_comms/...`)

`gofumpt` plus `golangci-lint`'s `goimports`/`gci` formatter (run via `task tools:fmt`) maintain
this layout. Do not reorder by hand.

```go
import (
    "context"
    "errors"
    "log/slog"

    "github.com/wrapped-owls/goremy-di/remy"
    commsv1 "github.com/0b1-PulsarTech/terectek_commsproto/gen/commstekproto-go/terekchat/v1"

    "github.com/0b1-PulsarTech/terectek_comms/libs/tereckernel/confloader"
)
```

## Aliasing

- Don't alias unless there is a name collision or the package name reads poorly at the call site.
- When you do alias generated proto packages, follow `<service>v<n>` (`commsv1`,
  `commsv1grpc`) consistently across the repo. The `importas` linter pins these aliases.
- Never alias the standard library (`time`, `context`, `errors`).

## Forbidden

- Dot imports (`import . "x"`) outside `*_test.go` files. They obscure the call site and break
  IDE navigation.
- Blank imports (`import _ "x"`) outside the `cmd/<app>/main.go` of an app, where they are used
  to register database drivers. Document each `_` import with a one‑line comment.
- Importing from `apps/...` inside `libs/...`. Direction is one‑way: apps depend on libs, never
  the reverse. See [`code-placement.md`](code-placement.md).
- Importing across apps (`apps/foo` → `apps/bar`). Share via `libs/`.
- Importing `internal/` from outside the parent module. The Go compiler enforces this; don't try
  to work around it with replace directives.

## Cyclic imports

A cycle means the package boundary is wrong. Resolutions, in order of preference:

1. Move the shared symbol down into a package both sides depend on.
2. Invert the direction by introducing a small interface in the consumer package — see
   [`types.md`](types.md).
3. Merge the two packages if the split was artificial.

Never break a cycle by importing inside a function.
