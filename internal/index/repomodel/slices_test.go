package repomodel

import "testing"

func TestMap(t *testing.T) {
	t.Parallel()

	got := Map([]int{1, 2, 3}, func(v int) int { return v * 2 })
	if len(got) != 3 || got[0] != 2 || got[2] != 6 {
		t.Fatalf("got %v", got)
	}
}

func TestMapEmpty(t *testing.T) {
	t.Parallel()

	got := Map([]string(nil), func(s string) int { return len(s) })
	if len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	got := Filter([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 })
	if len(got) != 2 || got[0] != 2 || got[1] != 4 {
		t.Fatalf("got %v", got)
	}
}
