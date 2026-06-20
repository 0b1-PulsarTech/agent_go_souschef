package repocontext

import (
	"context"
	"strings"
	"testing"
)

func TestSource(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)

	text, err := svc.Source(context.Background(), "CreateUser")
	if err != nil {
		t.Fatalf("source: %v", err)
	}

	if !strings.Contains(text, "func CreateUser") {
		t.Fatalf("unexpected source: %s", text)
	}
}
