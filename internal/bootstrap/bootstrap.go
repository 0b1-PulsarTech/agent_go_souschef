// Package bootstrap wires the application together. It is the only place that
// knows about the concrete implementations behind each port — every other
// package depends on interfaces declared in pkg/repocontext.
//
// Layout mirrors amigonimo's internal/bootstrap: bootstrap.go owns the DI
// composition, runner.go owns the per-subcommand runners (RunMCP, RunSync).
package bootstrap

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
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

	// Storage: the index is a throwaway cache, so it lives under the OS temp
	// dir keyed by a hash of the workspace path — never inside the project.
	// Successive runs on the same workspace still reuse the same file.
	remy.RegisterConstructorErr(inj, remy.Singleton[*reposqlite.Store],
		func() (*reposqlite.Store, error) {
			return reposqlite.Open(indexPath(cfg.Root))
		})

	// Source-level collaborators (cheap to build; singletons).
	remy.RegisterConstructor(inj, remy.Singleton[goscan.Indexer],
		func() goscan.Indexer { return goscan.New(cfg.Root) })
	remy.RegisterConstructor(inj, remy.Singleton[*gitprobe.Probe],
		func() *gitprobe.Probe { return gitprobe.New(cfg.Root) })

	// DuckTypeElements: true lets remy match the concrete types to the
	// LanguageIndexer / SymbolStore / ChangeReporter interface params.
	remy.RegisterConstructorArgs3(inj, remy.Singleton[repocontext.Service],
		repocontext.New)
}

// indexPath returns the throwaway SQLite path for a workspace: a stable,
// collision-free location under the OS temp dir derived from the workspace
// path. Keeping the index out of the project tree means souschef never dirties
// the repo it is indexing.
func indexPath(root string) string {
	abs, err := filepath.Abs(root)
	if err != nil {
		abs = root
	}

	sum := sha256.Sum256([]byte(abs))

	return filepath.Join(os.TempDir(), "agent_go_souschef", hex.EncodeToString(sum[:8]), "index.db")
}
