// Package b is a workspace fixture module that exercises generic receivers,
// which previously crashed the indexer's receiver-name extraction.
package b

// Box holds a single value of any type.
type Box[T any] struct{ V T }

// Get returns the boxed value (value receiver on a generic type).
func (b Box[T]) Get() T { return b.V }

// Set replaces the boxed value (pointer receiver on a generic type).
func (b *Box[T]) Set(v T) { b.V = v }
