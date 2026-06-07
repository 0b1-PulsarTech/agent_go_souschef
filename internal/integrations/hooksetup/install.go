package hooksetup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Run is the entry point for the `hook` subcommand.
// args examples: ["install", "--claude"], ["install", "--claude", "--codex"], ["run", "claude"]
func Run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(
			os.Stderr,
			"usage: agent_go-souschef hook install --<target> | hook run <target>",
		)
		return 1
	}
	switch args[0] {
	case "install":
		return runInstall(args[1:])
	case "run":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: agent_go-souschef hook run <target>")
			return 1
		}
		return runHook(args[1])
	default:
		fmt.Fprintf(os.Stderr, "unknown hook subcommand: %s\n", args[0])
		return 1
	}
}

func runInstall(flags []string) int {
	targets := make([]Target, 0, len(flags))
	for _, f := range flags {
		name := f
		if len(f) > 2 && f[:2] == "--" {
			name = f[2:]
		}
		t, ok := parseTarget(name)
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown target: %s\n", f)
			return 1
		}
		targets = append(targets, t)
	}
	if len(targets) == 0 {
		fmt.Fprintln(
			os.Stderr,
			"specify at least one target: --claude, --codex, --cursor, --gemini",
		)
		return 1
	}
	for _, t := range targets {
		if err := install(t); err != nil {
			fmt.Fprintf(os.Stderr, "install %s: %v\n", t, err)
			return 1
		}
		fmt.Printf("Installed hook for %s\n", t)
	}
	return 0
}

func install(t Target) error {
	path := t.ConfigPath()
	if path == "" {
		return fmt.Errorf("unsupported target: %s", t)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	switch t {
	case TargetClaude:
		return installClaude(path)
	default:
		return fmt.Errorf("install not yet implemented for %s", t)
	}
}

// claudeSettings is a minimal representation of ~/.claude/settings.json
type claudeSettings struct {
	Hooks map[string][]claudeHookGroup `json:"hooks,omitempty"`
	Extra map[string]json.RawMessage   `json:"-"`
}

type claudeHookGroup struct {
	Matcher string       `json:"matcher"`
	Hooks   []claudeHook `json:"hooks"`
}

type claudeHook struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

func installClaude(path string) error {
	raw := map[string]json.RawMessage{}
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &raw)
	}

	hookEntry := claudeHookGroup{
		Matcher: "Bash",
		Hooks:   []claudeHook{{Type: "command", Command: "agent_go-souschef hook run claude"}},
	}

	var existing []claudeHookGroup
	if hooksRaw, ok := raw["hooks"]; ok {
		var hooksMap map[string][]claudeHookGroup
		if err := json.Unmarshal(hooksRaw, &hooksMap); err == nil {
			existing = hooksMap["PreToolUse"]
		}
	}

	// Deduplicate: don't add if already present
	for _, g := range existing {
		for _, h := range g.Hooks {
			if h.Command == hookEntry.Hooks[0].Command {
				fmt.Printf("Hook already installed for claude\n")
				return nil
			}
		}
	}
	existing = append(existing, hookEntry)

	hooksMap := map[string][]claudeHookGroup{"PreToolUse": existing}
	hooksJSON, err := json.Marshal(hooksMap)
	if err != nil {
		return fmt.Errorf("marshal hooks: %w", err)
	}
	raw["hooks"] = hooksJSON

	out, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}
	return os.WriteFile(path, out, 0o644)
}
