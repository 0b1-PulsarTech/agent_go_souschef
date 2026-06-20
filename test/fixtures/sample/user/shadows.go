package user

import "strings"

// shadowExamples intentionally hides an identifier of every kind so the shadow
// pass has deterministic findings to report. It is never called.
func shadowExamples(in User) string {
	name := strings.TrimSpace(in.Name) // genuine use of the strings import

	strings := name       // shadows the imported package "strings"
	CreateUser := strings // shadows the package-level func CreateUser
	string := CreateUser  // shadows the predeclared type string

	result := string
	if result != "" {
		result := result + "!" // shadows the outer variable result
		return result
	}
	return result
}
