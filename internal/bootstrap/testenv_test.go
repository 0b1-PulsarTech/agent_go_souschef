package bootstrap

import (
	"os"
	"path/filepath"
	"testing"
)

// sampleWorkspace returns a fresh copy of the canonical fixtures/sample tree
// so each test mutates its own working dir without interfering with siblings.
func sampleWorkspace(t *testing.T) string {
	t.Helper()

	src := filepath.Join("..", "..", "test", "fixtures", "sample")
	dst := filepath.Join(t.TempDir(), "sample")
	copyTree(t, src, dst)

	return dst
}

func copyTree(t *testing.T, src, dst string) {
	t.Helper()

	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}

	if err := os.MkdirAll(dst, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	for _, entry := range entries {
		from, to := filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			copyTree(t, from, to)

			continue
		}

		data, err := os.ReadFile(from)
		if err != nil {
			t.Fatalf("read file: %v", err)
		}

		if err := os.WriteFile(to, data, 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}
}
