package repocontext

import (
	"context"
	"strings"
	"testing"
)

func TestQuery(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	if _, err := svc.Sync(context.Background()); err != nil {
		t.Fatalf("sync: %v", err)
	}
	text, err := svc.Query(context.Background(), "CreateUser", false)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if !strings.Contains(text, "CreateUser") {
		t.Fatalf("unexpected text: %s", text)
	}
}
