package goscan

import "path/filepath"

type Indexer struct {
	root string
}

func New(root string) *Indexer {
	return &Indexer{root: root}
}

func dataDir(root string) string {
	return filepath.Join(root, ".repo-context")
}

func best(names []string, query string) string {
	return BestMatch(names, query)
}
