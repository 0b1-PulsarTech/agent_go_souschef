// Package multimod is the root module of a fixture repo that has a nested
// module (./datastore) with its own go.mod but no go.work — the layout the
// indexer must cover by loading each module on its own.
package multimod

// App is a symbol that lives in the root module.
type App struct{ Name string }

// Run reports the app name.
func (a App) Run() string { return a.Name }
