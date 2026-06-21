package mcpsvc

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// DefaultSyncInterval is the minimum gap between automatic index refreshes
// triggered by tool calls. Calls arriving within the window answer from the
// index already in place rather than re-syncing.
const DefaultSyncInterval = 15 * time.Second

// Syncer is the slice of repocontext.Service the gate drives. Declaring it on
// the consumer side keeps the gate unit-testable with a fake and keeps the
// throttling concern out of the domain package.
type Syncer interface {
	Sync(ctx context.Context) (string, error)
}

// SyncGate keeps the index fresh without re-syncing on every MCP call. Ensure
// runs a sync only when the interval has elapsed since the last one; Force
// always syncs. The clock is injected so the throttle is testable without
// sleeping.
type SyncGate struct {
	syncer   Syncer
	now      func() time.Time
	interval time.Duration

	mu     sync.Mutex
	last   time.Time
	primed bool
}

// NewSyncGate builds a gate over syncer. now supplies the current time (pass
// time.Now in production); interval is the minimum gap between syncs.
func NewSyncGate(syncer Syncer, now func() time.Time, interval time.Duration) *SyncGate {
	return &SyncGate{syncer: syncer, now: now, interval: interval}
}

// Ensure runs a throttled sync before a read tool answers. It claims the window
// under the lock and syncs outside it, so concurrent calls collapse to one sync
// instead of piling up. A sync failure is logged, not returned: a stale index
// still answers, matching the non-fatal startup-sync policy.
func (g *SyncGate) Ensure(ctx context.Context) {
	if !g.claim() {
		return
	}

	if _, err := g.syncer.Sync(ctx); err != nil {
		slog.Warn("auto-sync failed", "err", err)
	}
}

// Force syncs unconditionally and resets the window. It backs the startup sync
// and the explicit souschef_sync tool, where the caller wants a refresh now
// regardless of the throttle and expects the summary or error surfaced.
func (g *SyncGate) Force(ctx context.Context) (string, error) {
	summary, err := g.syncer.Sync(ctx)

	g.mu.Lock()
	g.last = g.now()
	g.primed = true
	g.mu.Unlock()

	if err != nil {
		return summary, fmt.Errorf("force sync: %w", err)
	}

	return summary, nil
}

// claim reports whether the caller owns the current sync window, recording the
// attempt so a concurrent caller backs off. The sync itself runs after the lock
// is released.
func (g *SyncGate) claim() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := g.now()
	if g.primed && now.Sub(g.last) < g.interval {
		return false
	}

	g.last = now
	g.primed = true

	return true
}
