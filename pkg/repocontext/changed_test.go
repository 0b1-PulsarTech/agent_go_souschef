package repocontext

import (
	"context"
	"testing"
)

func TestChanged(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	text, err := svc.Changed(context.Background(), "user")
	if err != nil {
		t.Fatalf("changed: %v", err)
	}
	if text == "" {
		t.Fatal("expected response")
	}
}
