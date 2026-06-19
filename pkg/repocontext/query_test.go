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

// TestQueryTypeRef verifies that the type-reference edges flow through to query
// output. CreateUser takes Repository, Publisher, and User as parameters, so
// querying it must report those under "Uses types".
func TestQueryTypeRef(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	if _, err := svc.Sync(context.Background()); err != nil {
		t.Fatalf("sync: %v", err)
	}
	text, err := svc.Query(context.Background(), "CreateUser", false)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if !strings.Contains(text, "Uses types:") {
		t.Fatalf("expected a Uses types section: %s", text)
	}
	for _, want := range []string{"Repository", "Publisher", "User"} {
		if !strings.Contains(text, want) {
			t.Errorf("expected CreateUser to use type %q: %s", want, text)
		}
	}
}
