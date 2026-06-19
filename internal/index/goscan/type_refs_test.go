package goscan

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

// TestTypeRefs verifies that types used only as parameter/field types (with no
// associated call edge) are recorded as "ref" relations. In the sample fixture
// CreateUser takes Repository, Publisher, and User as parameters.
func TestTypeRefs(t *testing.T) {
	t.Parallel()
	idx := New(filepath.Join("..", "..", "..", "test", "fixtures", "sample"))
	snap, err := idx.Build(context.Background())
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	name := map[int64]string{}
	for _, sym := range snap.Symbols {
		name[sym.ID] = sym.Name
	}

	refs := map[string]bool{}
	for _, rel := range snap.TypeRefs {
		if rel.Kind != "ref" {
			t.Errorf("type ref with unexpected kind %q", rel.Kind)
		}
		if rel.FromID == rel.ToID {
			t.Errorf("self type ref on %q", name[rel.FromID])
		}
		refs[name[rel.FromID]+"->"+name[rel.ToID]] = true
	}

	for _, want := range []string{
		"CreateUser->Repository",
		"CreateUser->Publisher",
		"CreateUser->User",
	} {
		if !refs[want] {
			t.Errorf("missing type ref %q; got %v", want, refs)
		}
	}

	// Builtins (error, string) are not indexed symbols, so they must not appear.
	for edge := range refs {
		if strings.HasSuffix(edge, "->error") || strings.HasSuffix(edge, "->string") {
			t.Errorf("builtin leaked into type refs: %q", edge)
		}
	}
}
