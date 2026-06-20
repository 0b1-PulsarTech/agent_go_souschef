package goscan

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

func (idx Indexer) Source(_ context.Context, query string) (string, error) {
	fset := token.NewFileSet()
	err := filepath.WalkDir(idx.root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() || filepath.Ext(path) != ".go" {
			return walkErr
		}

		file, parseErr := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if parseErr != nil {
			return fmt.Errorf("parse %s: %w", path, parseErr)
		}

		if src := matchDecl(fset, file, query); src != "" {
			return stopSource(src)
		}

		return nil
	})

	var stop *sourceFound

	if err == nil {
		return "", fmt.Errorf("symbol source not found")
	}

	if ok := errorAsSource(err, &stop); ok {
		return stop.text, nil
	}

	return "", fmt.Errorf("scan sources: %w", err)
}

func matchDecl(fset *token.FileSet, file *ast.File, query string) string {
	for _, decl := range file.Decls {
		if out := formatDecl(fset, decl, query); out != "" {
			return out
		}
	}

	return ""
}

func nodeText(fset *token.FileSet, node ast.Node) string {
	var buf bytes.Buffer
	_ = format.Node(&buf, fset, node)

	return strings.TrimSpace(buf.String())
}

type sourceFound struct{ text string }

func (s *sourceFound) Error() string { return s.text }

func stopSource(text string) error { return &sourceFound{text: text} }

func errorAsSource(err error, target **sourceFound) bool {
	found, ok := err.(*sourceFound)
	if ok {
		*target = found
	}

	return ok
}
