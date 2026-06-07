package goscan

import "testing"

func TestDataDir(t *testing.T) {
	t.Parallel()
	if got := dataDir("/tmp/root"); got == "" {
		t.Fatal("expected dir")
	}
}
