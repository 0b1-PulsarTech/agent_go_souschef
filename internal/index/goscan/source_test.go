package goscan

import "testing"

func TestSourcePath(t *testing.T) {
	t.Parallel()
	if got := sourcePath("/tmp", "a.go"); got == "" {
		t.Fatal("expected path")
	}
}
