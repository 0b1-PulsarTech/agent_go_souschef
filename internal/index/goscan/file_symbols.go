package goscan

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/types"
	"os"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan/graph"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan/symbols"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
	"golang.org/x/tools/go/packages"
)

func (b *snapshotBuilder) addFile(pkg *packages.Package, file *ast.File) error {
	path := rel(b.root, pkg.Fset, file.Pos())
	summary := symbols.FileSummary{Path: path, Pkg: pkg.PkgPath}

	for _, decl := range file.Decls {
		b.addDecl(pkg, path, decl, &summary)
	}

	content, err := os.ReadFile(pkg.Fset.Position(file.Pos()).Filename)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	sum := sha256.Sum256(content)
	b.snapshot.Files = append(b.snapshot.Files, repomodel.FileSummary{
		Path: path, Lang: "go", Hash: hex.EncodeToString(sum[:]), Summary: summary.Text(),
	})

	return nil
}

func (b *snapshotBuilder) register(obj types.Object, symbol repomodel.Symbol) int64 {
	b.nextID++
	symbol.ID = b.nextID
	b.snapshot.Symbols = append(b.snapshot.Symbols, symbol)
	b.ids[obj] = symbol.ID
	b.names[graph.FullName(obj)] = symbol.ID

	return symbol.ID
}
