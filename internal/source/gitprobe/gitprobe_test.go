package gitprobe

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
)

func TestNewNonRepo(t *testing.T) {
	t.Parallel()

	p := New(t.TempDir())
	if p == nil {
		t.Fatal("expected probe")
	}

	got, err := p.Changed(context.Background(), "")
	if err != nil {
		t.Fatalf("changed: %v", err)
	}

	if got != msgNoRepo {
		t.Fatalf("expected degraded message, got %q", got)
	}
}

func TestChangedDetectsDirty(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if _, err := git.PlainInit(root, false); err != nil {
		t.Fatalf("init: %v", err)
	}

	if err := os.WriteFile(
		filepath.Join(root, "dirty.go"),
		[]byte("package x"),
		0o644,
	); err != nil {
		t.Fatalf("write: %v", err)
	}

	p := New(root)

	got, err := p.Changed(context.Background(), "")
	if err != nil {
		t.Fatalf("changed: %v", err)
	}

	if got == msgNoChanges || got == msgNoRepo {
		t.Fatalf("expected dirty file listed, got %q", got)
	}
}
