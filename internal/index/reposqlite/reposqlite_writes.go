package reposqlite

import (
	"context"
	"fmt"
	"slices"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/reposqlite/db"
)

// Reset truncates every table so a fresh snapshot can be written.
func (s *Store) Reset(ctx context.Context) error {
	for _, op := range []struct {
		name string
		fn   func(context.Context) error
	}{
		{"relations", s.q.DeleteAllRelations},
		{"methods", s.q.DeleteAllMethods},
		{"symbols", s.q.DeleteAllSymbols},
		{"files", s.q.DeleteAllFiles},
	} {
		if err := op.fn(ctx); err != nil {
			return fmt.Errorf("reset %s: %w", op.name, err)
		}
	}
	return nil
}

// Write persists a snapshot. Relations from Calls, Implementations, and
// TypeRefs are flattened into a single insert pass.
func (s *Store) Write(ctx context.Context, snap repomodel.Snapshot) error {
	if err := s.insertSymbols(ctx, snap.Symbols); err != nil {
		return err
	}
	if err := s.insertRelations(
		ctx,
		slices.Concat(snap.Calls, snap.Implementations, snap.TypeRefs),
	); err != nil {
		return err
	}
	if err := s.insertFiles(ctx, snap.Files); err != nil {
		return err
	}
	return s.insertMethods(ctx, snap.Methods)
}

func (s *Store) insertSymbols(ctx context.Context, in []repomodel.Symbol) error {
	for _, sym := range in {
		if err := s.q.InsertSymbol(ctx, db.InsertSymbolParams{
			ID: sym.ID, Name: sym.Name, Kind: sym.Kind,
			Package: sym.Package, File: sym.File, Signature: sym.Signature,
		}); err != nil {
			return fmt.Errorf("insert symbol %q: %w", sym.Name, err)
		}
	}
	return nil
}

func (s *Store) insertRelations(ctx context.Context, in []repomodel.Relation) error {
	for _, rel := range in {
		if err := s.q.InsertRelation(ctx, db.InsertRelationParams{
			FromID: rel.FromID, ToID: rel.ToID, EdgeKind: rel.Kind,
		}); err != nil {
			return fmt.Errorf("insert relation %d→%d: %w", rel.FromID, rel.ToID, err)
		}
	}
	return nil
}

func (s *Store) insertFiles(ctx context.Context, in []repomodel.FileSummary) error {
	for _, f := range in {
		if err := s.q.InsertFile(ctx, db.InsertFileParams{
			Path: f.Path, Lang: f.Lang, Hash: f.Hash, Summary: f.Summary,
		}); err != nil {
			return fmt.Errorf("insert file %q: %w", f.Path, err)
		}
	}
	return nil
}

func (s *Store) insertMethods(ctx context.Context, in []repomodel.Method) error {
	for _, m := range in {
		if err := s.q.InsertMethod(ctx, db.InsertMethodParams{
			ParentID: m.ParentID, Name: m.Name, Signature: m.Signature, Kind: m.MemberKind,
		}); err != nil {
			return fmt.Errorf("insert method %q: %w", m.Name, err)
		}
	}
	return nil
}
