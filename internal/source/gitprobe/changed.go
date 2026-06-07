package gitprobe

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// Changed returns the worktree's modified-file list, optionally filtered to
// paths containing scope (case-insensitive substring). The signature matches
// repocontext.ChangeReporter so the Service can call it directly.
func (p *Probe) Changed(_ context.Context, scope string) (string, error) {
	repo, err := p.repoOrErr()
	if err != nil {
		return "No git diff available.", nil
	}
	wt, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("worktree: %w", err)
	}
	st, err := wt.Status()
	if err != nil {
		return "", fmt.Errorf("status: %w", err)
	}

	needle := strings.ToLower(scope)
	files := make([]string, 0, len(st))
	for path := range st {
		if needle == "" || strings.Contains(strings.ToLower(path), needle) {
			files = append(files, path)
		}
	}
	sort.Strings(files)
	return render(scope, files), nil
}

// render formats a sorted file list into the compact MCP-friendly response.
func render(scope string, files []string) string {
	if len(files) == 0 {
		return "No matching changes."
	}
	prefix := "Modified:\n"
	if scope != "" {
		prefix = scope + "\n\nModified:\n"
	}
	return prefix + "  " + strings.Join(files, "\n  ")
}
