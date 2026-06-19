// Package a is a workspace fixture module. It exists only so the indexer can
// be exercised against a go.work monorepo with more than one module.
package a

// Greeter produces greetings for a name.
type Greeter struct{ Name string }

// Hello returns a greeting for the receiver's name.
func (g Greeter) Hello() string { return "hi " + g.Name }
