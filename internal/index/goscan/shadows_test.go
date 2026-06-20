package goscan

import (
	"context"
	"path/filepath"
	"testing"
)

// TestShadows verifies the scope-based detector classifies each kind of
// shadowing in the sample fixture's shadows.go: a builtin, an imported
// package, a package-level symbol, and an outer variable.
func TestShadows(t *testing.T) {
	t.Parallel()

	idx := New(filepath.Join("..", "..", "..", "test", "fixtures", "sample"))

	snap, err := idx.Build(context.Background())
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	// Keyed by "name/origin" so each finding is unambiguous.
	got := map[string]string{}
	for _, sh := range snap.Shadows {
		got[sh.Name+"/"+sh.Origin] = sh.Detail

		if sh.Line == 0 || sh.File == "" {
			t.Errorf("shadow %q missing position: %+v", sh.Name, sh)
		}
	}

	for _, want := range []string{
		"string/builtin",     // local `string` hides the predeclared type
		"strings/import",     // local `strings` hides the imported package
		"CreateUser/package", // local `CreateUser` hides the package-level func
		"result/outer",       // inner `result` hides the outer variable
	} {
		if _, ok := got[want]; !ok {
			t.Errorf("missing shadow %q; got %v", want, got)
		}
	}

	if detail := got["strings/import"]; detail != `import "strings"` {
		t.Errorf("import shadow detail = %q, want %q", detail, `import "strings"`)
	}

	if got["string/builtin"] != "predeclared" {
		t.Errorf("builtin shadow detail = %q, want %q", got["string/builtin"], "predeclared")
	}
}
