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
	Shadows         []Shadow
}

// Shadow is one declaration that hides an identifier visible in an enclosing
// scope. Origin classifies what was hidden: "builtin" (a predeclared name such
// as len/error/string), "import" (an imported package name), "package" (a
// package-level symbol), or "outer" (a variable or parameter from an enclosing
// function or block). Detail says where the hidden identifier comes from.
type Shadow struct {
	File   string
	Line   int
	Column int
	Name   string
	Origin string
	Detail string
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
