package goscan

import "testing"

func TestBestMatch(t *testing.T) {
	t.Parallel()
	got := BestMatch([]string{"Store.UpsertSymbol", "CreateUser"}, "UpsertSymbol")
	if got != "Store.UpsertSymbol" {
		t.Fatalf("got %q", got)
	}
}
