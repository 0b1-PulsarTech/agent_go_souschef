// Package reposqlite is the SQLite-backed implementation of repocontext.SymbolStore.
// Reads/writes are split into _reads.go / _writes.go; this file owns the connection
// lifecycle and the per-row mapper helpers reused on both sides.
package reposqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/reposqlite/db"
)

//go:embed sql/schema.sql
var schemaDDL string

// dirPerm is the permission for the index's parent directory: owner rwx, group
// rx, no world access (gosec rejects anything looser than 0o750).
const dirPerm = 0o750

// Store wraps sqlc-generated queries. The embedded *sql.DB is kept so Close can
// hand it back to the runtime when the process exits.
type Store struct {
	conn *sql.DB
	q    *db.Queries
}

// Open creates (or opens) the SQLite database at path, applies the embedded
// schema, and returns a ready-to-use Store.
func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}

	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if _, err = conn.ExecContext(context.Background(), schemaDDL); err != nil {
		_ = conn.Close()

		return nil, fmt.Errorf("apply schema: %w", err)
	}

	return &Store{conn: conn, q: db.New(conn)}, nil
}

// Close releases the underlying database connection.
func (s *Store) Close() error {
	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("close db: %w", err)
	}

	return nil
}

// toSymbol is the canonical sqlc-row → domain mapper, used by every read path
// that returns symbols. Centralising it keeps reads/writes one-liners.
func toSymbol(r db.Symbol) repomodel.Symbol {
	return repomodel.Symbol{
		ID: r.ID, Name: r.Name, Kind: r.Kind,
		Package: r.Package, File: r.File, Signature: r.Signature,
	}
}

func toShadow(r db.ListShadowsRow) repomodel.Shadow {
	return repomodel.Shadow{
		File: r.File, Line: int(r.Line), Column: int(r.Col),
		Name: r.Name, Origin: r.Origin, Detail: r.Detail,
	}
}
