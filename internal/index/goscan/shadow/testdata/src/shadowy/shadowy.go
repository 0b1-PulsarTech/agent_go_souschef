// Package shadowy feeds analysistest: every shadowing site is annotated with
// the diagnostic the pass must emit.
package shadowy

import "strings"

// Pkg and pkgVar are package-level symbols a local can hide.
func Pkg() {}

var pkgVar = 1

func builtinShadow() int {
	len := 3 // want `"len" shadows a builtin`

	return len
}

func importShadow() string {
	s := strings.TrimSpace(" x ") // genuine use of the import
	strings := s                  // want `"strings" shadows import "strings"`

	return strings
}

func packageShadow() int {
	Pkg := 2    // want `"Pkg" shadows a package-level declaration`
	pkgVar := 3 // want `"pkgVar" shadows a package-level declaration`

	return Pkg + pkgVar
}

func outerShadow() int {
	a := 1
	{
		a := 2 // want `"a" shadows an outer declaration`

		return a
	}

	return a
}

// elseIfShadow is the case the scope walk must not miss: the else-if init
// re-declares the names the if init introduced, so the inner pair shadows the
// outer one.
func elseIfShadow(m map[string]int) int {
	if a, ok := m["x"]; ok {
		return a
	} else if a, ok := m["y"]; ok { // want `"a" shadows an outer declaration` `"ok" shadows an outer declaration`
		return a + btoi(ok)
	}

	return 0
}

func btoi(b bool) int {
	if b {
		return 1
	}

	return 0
}
