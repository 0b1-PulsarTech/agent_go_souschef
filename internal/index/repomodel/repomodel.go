package repomodel

type Symbol struct {
	ID        int64
	Name      string
	Kind      string
	Package   string
	File      string
	Signature string
}

type Relation struct {
	FromID int64
	ToID   int64
	Kind   string
}

type Snapshot struct {
	Symbols         []Symbol
	Files           []FileSummary
	Calls           []Relation
	Implementations []Relation
	TypeRefs        []Relation
	Methods         []Method
}

type FileSummary struct {
	Path    string
	Lang    string
	Hash    string
	Summary string
}

type Method struct {
	ParentID   int64
	Name       string
	Signature  string
	MemberKind string
}

type TextHit struct {
	Path    string
	Snippet string
}

type QueryHit struct {
	Symbol          Symbol
	Calls           []string
	Callers         []string
	Implementations []string
	UsedBy          []string
	UsesTypes       []string
	Methods         []string
	TextHits        []TextHit
	Source          string
}
