package repomodel

// Map applies fn to each element of in and returns a same-length result slice.
// Reused by reposqlite read methods to translate sqlc rows into domain types
// without writing the per-method for-loop seven times.
func Map[A, B any](in []A, fn func(A) B) []B {
	out := make([]B, len(in))
	for i, v := range in {
		out[i] = fn(v)
	}

	return out
}

// Filter returns the elements for which keep returns true.
func Filter[A any](in []A, keep func(A) bool) []A {
	out := make([]A, 0, len(in))

	for _, v := range in {
		if keep(v) {
			out = append(out, v)
		}
	}

	return out
}
