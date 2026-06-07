// Package gitprobe answers "what changed?" against a Git working tree without
// shelling out — go-git reads the index/refs directly.
package gitprobe

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
)

// Probe wraps a *git.Repository for the chosen worktree.
type Probe struct {
	root string
	repo *git.Repository
}

// New opens the Git repository rooted at root. If the directory is not a Git
// repository, callers receive a Probe in degraded mode: Changed reports an
// empty diff rather than failing, so a fresh checkout still works under MCP.
func New(root string) *Probe {
	repo, err := git.PlainOpenWithOptions(root, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		if !errors.Is(err, git.ErrRepositoryNotExists) {
			// Surface unexpected errors via Probe state — Changed will report them.
			return &Probe{root: root}
		}
		return &Probe{root: root}
	}
	return &Probe{root: root, repo: repo}
}

// repoOrErr is the small read accessor used by sibling files; centralising it
// here means new operations don't each re-check the degraded-mode flag.
func (p *Probe) repoOrErr() (*git.Repository, error) {
	if p.repo == nil {
		return nil, fmt.Errorf("%s is not a git repository", p.root)
	}
	return p.repo, nil
}
