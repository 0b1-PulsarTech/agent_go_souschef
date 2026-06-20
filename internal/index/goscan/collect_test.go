package goscan

import (
	"go/token"
	"testing"
)

func TestRel(t *testing.T) {
	t.Parallel()

	fset := token.NewFileSet()
	file := fset.AddFile("/tmp/root/a.go", -1, 10)

	if got := rel("/tmp/root", fset, file.Pos(1)); got != "a.go" {
		t.Fatalf("got %q", got)
	}
}
