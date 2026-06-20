package symbols

import (
	"go/ast"
	"testing"
)

func TestKindInterface(t *testing.T) {
	t.Parallel()

	spec := &ast.TypeSpec{Type: &ast.InterfaceType{}}
	if got := Kind(spec); got != "interface" {
		t.Fatalf("got %q", got)
	}
}

func TestKindStruct(t *testing.T) {
	t.Parallel()

	spec := &ast.TypeSpec{Type: &ast.StructType{}}
	if got := Kind(spec); got != "struct" {
		t.Fatalf("got %q", got)
	}
}

func TestKindOther(t *testing.T) {
	t.Parallel()

	spec := &ast.TypeSpec{Type: &ast.Ident{}}
	if got := Kind(spec); got != "type" {
		t.Fatalf("got %q", got)
	}
}

func TestRecvName(t *testing.T) {
	t.Parallel()

	decl := &ast.FuncDecl{Recv: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "Svc"}}}}}
	if got := RecvName(decl); got != "Svc" {
		t.Fatalf("got %q", got)
	}
}

func TestRecvNamePointer(t *testing.T) {
	t.Parallel()

	decl := &ast.FuncDecl{Recv: &ast.FieldList{List: []*ast.Field{
		{Type: &ast.StarExpr{X: &ast.Ident{Name: "Repo"}}},
	}}}
	if got := RecvName(decl); got != "Repo" {
		t.Fatalf("got %q", got)
	}
}

func TestFileSummaryText(t *testing.T) {
	t.Parallel()

	s := FileSummary{Pkg: "user", Exports: []string{"CreateUser", "DeleteUser"}}

	text := s.Text()
	if text == "" {
		t.Fatal("empty summary")
	}
}

func TestMethodRecord(t *testing.T) {
	t.Parallel()

	m := MethodRecord(42, "Save")
	if m.ParentID != 42 || m.Name != "Save" {
		t.Fatalf("unexpected %+v", m)
	}
}
