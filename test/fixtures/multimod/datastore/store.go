// Package datastore is a nested module: it carries its own go.mod, so the root
// module's "./..." never reaches it and it must be loaded separately.
package datastore

// Store is a symbol that lives only in the nested module.
type Store struct{ DSN string }

// Open reports the configured DSN.
func (s Store) Open() string { return s.DSN }
