package mcpsvc

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// fakeSyncer counts Sync calls and can be told to fail. The mutex keeps the
// counter race-free under the -race detector even though the gate's tests are
// single-threaded.
type fakeSyncer struct {
	mu    sync.Mutex
	calls int
	err   error
}

func (f *fakeSyncer) Sync(context.Context) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.calls++

	return "synced", f.err
}

func (f *fakeSyncer) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.calls
}

// fakeClock is a hand-advanced clock so the throttle is tested without sleeping.
type fakeClock struct{ t time.Time }

func (c *fakeClock) now() time.Time          { return c.t }
func (c *fakeClock) advance(d time.Duration) { c.t = c.t.Add(d) }

func TestSyncGateEnsureThrottles(t *testing.T) {
	t.Parallel()

	clock := &fakeClock{t: time.Unix(1000, 0)}
	syncer := &fakeSyncer{}
	gate := NewSyncGate(syncer, clock.now, 15*time.Second)

	gate.Ensure(context.Background())

	if got := syncer.count(); got != 1 {
		t.Fatalf("first Ensure synced %d times, want 1", got)
	}

	clock.advance(14 * time.Second)
	gate.Ensure(context.Background())

	if got := syncer.count(); got != 1 {
		t.Fatalf("Ensure inside the window synced %d times, want 1", got)
	}

	clock.advance(time.Second) // 15s since the last sync
	gate.Ensure(context.Background())

	if got := syncer.count(); got != 2 {
		t.Fatalf("Ensure after the window synced %d times, want 2", got)
	}
}

func TestSyncGateForceResetsWindow(t *testing.T) {
	t.Parallel()

	clock := &fakeClock{t: time.Unix(1000, 0)}
	syncer := &fakeSyncer{}
	gate := NewSyncGate(syncer, clock.now, 15*time.Second)

	if _, err := gate.Force(context.Background()); err != nil {
		t.Fatalf("Force: %v", err)
	}

	if got := syncer.count(); got != 1 {
		t.Fatalf("Force synced %d times, want 1", got)
	}

	// A read inside the window Force just opened must not sync again.
	clock.advance(5 * time.Second)
	gate.Ensure(context.Background())

	if got := syncer.count(); got != 1 {
		t.Fatalf("Ensure after Force inside the window synced %d times, want 1", got)
	}
}

func TestSyncGateEnsureSwallowsError(t *testing.T) {
	t.Parallel()

	clock := &fakeClock{t: time.Unix(1000, 0)}
	syncer := &fakeSyncer{err: errors.New("boom")}
	gate := NewSyncGate(syncer, clock.now, 15*time.Second)

	gate.Ensure(context.Background()) // a failed sync is logged, never propagated

	if got := syncer.count(); got != 1 {
		t.Fatalf("Ensure attempted %d syncs, want 1", got)
	}
}
