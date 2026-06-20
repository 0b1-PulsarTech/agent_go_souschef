package gitprobe

import "testing"

func TestRender(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name, scope string
		files       []string
		want        string
	}{
		{"empty", "", nil, "No matching changes."},
		{"scoped empty", "auth", nil, "No matching changes."},
		{"single", "", []string{"a.go"}, "Modified:\n  a.go"},
		{"scoped header", "auth", []string{"auth/x.go"}, "auth\n\nModified:\n  auth/x.go"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := render(tc.scope, tc.files); got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}
