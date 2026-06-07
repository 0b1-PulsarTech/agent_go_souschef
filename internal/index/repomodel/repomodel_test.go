package repomodel

import "testing"

func TestZeroSnapshot(t *testing.T) {
	t.Parallel()
	var snap Snapshot
	if len(snap.Symbols) != 0 {
		t.Fatal("expected zero snapshot")
	}
}
