package shadow

import (
	"fmt"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/checker"
	"golang.org/x/tools/go/packages"
)

// Run executes the shadow Analyzer over pkgs using the analysis driver and
// returns every finding. The driver builds each analysis.Pass (FileSet, type
// info) and runs the pass per package in parallel; we read each root action's
// typed Result. Packages must be loaded with type info and TypesSizes
// (packages.NeedTypes|NeedTypesSizes|NeedTypesInfo|NeedSyntax), or the driver
// rejects them.
//
// Findings are returned as produced; one source position may appear more than
// once when a file is type-checked under several package variants (its test
// build, say), so callers that fold findings into a single index dedupe by
// position.
func Run(pkgs []*packages.Package) ([]Finding, error) {
	graph, err := checker.Analyze([]*analysis.Analyzer{Analyzer}, pkgs, nil)
	if err != nil {
		return nil, fmt.Errorf("analyze shadows: %w", err)
	}

	var findings []Finding

	for _, root := range graph.Roots {
		if root.Err != nil {
			return nil, fmt.Errorf("shadow pass on %s: %w", root.Package.PkgPath, root.Err)
		}

		if result, ok := root.Result.([]Finding); ok {
			findings = append(findings, result...)
		}
	}

	return findings, nil
}
