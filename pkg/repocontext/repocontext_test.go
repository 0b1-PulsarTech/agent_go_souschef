package repocontext

import "testing"

func TestNew(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	if svc.indexer == nil || svc.store == nil || svc.changes == nil {
		t.Fatal("expected fully-wired Service")
	}
}
