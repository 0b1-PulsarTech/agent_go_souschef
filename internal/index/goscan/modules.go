package goscan

import (
	"io/fs"
	"os"
	"path/filepath"
)

// loadPlan is one packages.Load invocation: a working directory and the
// patterns to load from it.
type loadPlan struct {
	dir      string
	patterns []string
}

// loadPlans decides how to cover every Go module under root.
//
//   - go.work workspace: a single workspace-aware load from the root covering
//     every module, so inter-module dependencies resolve and cross-module call
//     edges link. A bare "./..." fails here because the root holds no go.mod.
//   - otherwise: one load per module — the root module plus any nested module
//     with its own go.mod (a database/entschema, say) — since each resolves its
//     dependencies on its own.
func loadPlans(root string) ([]loadPlan, error) {
	dirs, err := moduleDirs(root)
	if err != nil {
		return nil, err
	}
	if exists(filepath.Join(root, "go.work")) {
		return []loadPlan{{dir: root, patterns: workspacePatterns(root, dirs)}}, nil
	}
	if len(dirs) == 0 {
		return []loadPlan{{dir: root, patterns: []string{"./..."}}}, nil
	}
	plans := make([]loadPlan, 0, len(dirs))
	for _, dir := range dirs {
		plans = append(plans, loadPlan{dir: dir, patterns: []string{"./..."}})
	}
	return plans, nil
}

func workspacePatterns(root string, dirs []string) []string {
	patterns := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		rel, err := filepath.Rel(root, dir)
		if err != nil {
			continue
		}
		patterns = append(patterns, modulePattern(rel))
	}
	if len(patterns) == 0 {
		return []string{"./..."}
	}
	return patterns
}

func moduleDirs(root string) ([]string, error) {
	var dirs []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if path != root && skipDir(d.Name()) {
				return fs.SkipDir
			}
			return nil
		}
		if d.Name() == "go.mod" {
			dirs = append(dirs, filepath.Dir(path))
		}
		return nil
	})
	return dirs, err
}

func modulePattern(rel string) string {
	if rel == "." {
		return "./..."
	}
	return "./" + filepath.ToSlash(rel) + "/..."
}

func skipDir(name string) bool {
	switch name {
	case ".git", "vendor", "node_modules", ".repo-context":
		return true
	default:
		return false
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
