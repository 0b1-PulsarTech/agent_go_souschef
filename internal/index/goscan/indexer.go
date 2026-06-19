package goscan

// Indexer scans a Go workspace rooted at root. It holds only the root path, so
// it is passed by value.
type Indexer struct {
	root string
}

func New(root string) Indexer {
	return Indexer{root: root}
}
