package repocontext

import (
	"context"
	"testing"
)

func TestSync(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	text, err := svc.Sync(context.Background())
	if err != nil {
		t.Fatalf("sync: %v", err)
	}
	if text != "Synced." {
		t.Fatalf("got %q", text)
	}
}
