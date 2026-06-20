package repocontext

import (
	"context"
	"strings"
	"testing"
)

// TestShadows runs the full path — build, persist, read, render — against the
// sample fixture's shadows.go and checks every kind of finding survives the
// round-trip through SQLite.
func TestShadows(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	if _, err := svc.Sync(context.Background()); err != nil {
		t.Fatalf("sync: %v", err)
	}

	text, err := svc.Shadows(context.Background(), "")
	if err != nil {
		t.Fatalf("shadows: %v", err)
	}

	for _, want := range []string{
		`strings → shadows import "strings"`,
		"CreateUser → shadows a package-level symbol",
		"string → shadows a builtin/predeclared identifier",
		"result → shadows an outer variable",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("shadow report missing %q:\n%s", want, text)
		}
	}
}

// TestShadowsScope narrows the report by path; a non-matching scope yields the
// empty-report message rather than findings.
func TestShadowsScope(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	if _, err := svc.Sync(context.Background()); err != nil {
		t.Fatalf("sync: %v", err)
	}

	in, err := svc.Shadows(context.Background(), "shadows.go")
	if err != nil {
		t.Fatalf("shadows: %v", err)
	}

	if !strings.Contains(in, "scope: shadows.go") || !strings.Contains(in, "shadows import") {
		t.Errorf("scoped report should keep matching findings:\n%s", in)
	}

	out, err := svc.Shadows(context.Background(), "no/such/path")
	if err != nil {
		t.Fatalf("shadows: %v", err)
	}

	if !strings.Contains(out, "No variable shadowing") {
		t.Errorf("non-matching scope should report nothing, got:\n%s", out)
	}
}
