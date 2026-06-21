package goscan

import (
	"fmt"
	"go/token"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan/shadow"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/repomodel"
)

// addShadows runs the shadow analysis pass over every loaded package and folds
// its findings into the snapshot. Findings are deduped by source position
// because one file may be type-checked under several package variants (its test
// build, say), which would otherwise report the same shadow twice.
func (b *snapshotBuilder) addShadows() error {
	findings, err := shadow.Run(b.pkgs)
	if err != nil {
		return fmt.Errorf("run shadow pass: %w", err)
	}

	seen := map[token.Position]struct{}{}

	for _, f := range findings {
		if _, dup := seen[f.Pos]; dup {
			continue
		}

		seen[f.Pos] = struct{}{}

		b.snapshot.Shadows = append(b.snapshot.Shadows, repomodel.Shadow{
			File:   relPath(b.root, f.Pos.Filename),
			Line:   f.Pos.Line,
			Column: f.Pos.Column,
			Name:   f.Name,
			Origin: f.Origin,
			Detail: shadowDetail(b.root, f),
		})
	}

	return nil
}

// shadowDetail renders the "what it hides" half of a finding. The position of
// the hidden declaration is made repo-relative to match the rest of the index.
func shadowDetail(root string, f shadow.Finding) string {
	switch f.Origin {
	case shadow.OriginBuiltin:
		return "predeclared"
	case shadow.OriginImport:
		return fmt.Sprintf("import %q", f.ImportPath)
	default:
		return fmt.Sprintf("declared at %s:%d", relPath(root, f.Hidden.Filename), f.Hidden.Line)
	}
}
