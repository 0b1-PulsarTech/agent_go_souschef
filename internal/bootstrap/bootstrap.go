// Package bootstrap wires the application together. It is the only place that
// knows about the concrete implementations behind each port — every other
// package depends on interfaces declared in pkg/repocontext.
//
// Layout mirrors amigonimo's internal/bootstrap: bootstrap.go owns the DI
// composition, runner.go owns the per-subcommand runners (RunMCP, RunSync).
package bootstrap

import (
	"path/filepath"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/goscan"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/index/reposqlite"
	"github.com/0b1-PulsarTech/agent_go_souschef/internal/source/gitprobe"
	"github.com/0b1-PulsarTech/agent_go_souschef/pkg/repocontext"
)

// Config is the bootstrap-time configuration the user provides on the CLI.
type Config struct {
	// Root is the workspace directory the indexer scans. Defaults to the cwd
	// in main.go.
	Root string
	// Version is reported to MCP clients in the implementation handshake.
	Version string
}

// DoInjections registers every collaborator on the injector. Layers stay
// inline — each is two lines — until any of them grows past ~5 registrations,
// at which point it earns its own remy.Module.
func DoInjections(inj remy.Injector, cfg Config) {
	remy.RegisterInstance(inj, cfg)

	// Storage: the SQLite store lives at <root>/.repo-context/index.db so
	// successive runs reuse the same index file.
	remy.RegisterConstructorErr(inj, remy.Singleton[*reposqlite.Store],
		func() (*reposqlite.Store, error) {
			return reposqlite.Open(filepath.Join(cfg.Root, ".repo-context", "index.db"))
		})

	// Source-level collaborators (cheap to build; singletons).
	remy.RegisterConstructor(inj, remy.Singleton[*goscan.Indexer],
		func() *goscan.Indexer { return goscan.New(cfg.Root) })
	remy.RegisterConstructor(inj, remy.Singleton[*gitprobe.Probe],
		func() *gitprobe.Probe { return gitprobe.New(cfg.Root) })

	// Service is the only thing that needs three deps wired together.
	// DuckTypeElements: true lets remy match the concrete types to the
	// LanguageIndexer / SymbolStore / ChangeReporter interface params.
	remy.RegisterConstructorArgs3(inj, remy.Singleton[*repocontext.Service],
		repocontext.New)
}
