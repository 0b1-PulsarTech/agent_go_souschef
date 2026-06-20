package repocontext

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/reposqlite"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/source/gitprobe"
)

// sampleRoot returns a fresh copy of fixtures/sample under t.TempDir.
func sampleRoot(t *testing.T) string {
	t.Helper()

	src := filepath.Join("..", "..", "test", "fixtures", "sample")
	dst := filepath.Join(t.TempDir(), "sample")
	copyTree(t, src, dst)

	return dst
}

// newTestService wires a Service against the sample fixture with real
// collaborators (goscan + SQLite + gitprobe). Mirrors the production wiring
// in internal/bootstrap but without remy.
func newTestService(t *testing.T) Service {
	t.Helper()
	root := sampleRoot(t)

	store, err := reposqlite.Open(filepath.Join(root, ".repo-context", "index.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}

	t.Cleanup(func() { _ = store.Close() })

	return New(goscan.New(root), store, gitprobe.New(root))
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
		from := filepath.Join(src, entry.Name())
		to := filepath.Join(dst, entry.Name())

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
