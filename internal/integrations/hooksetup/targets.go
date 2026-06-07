package hooksetup

import (
	"os"
	"path/filepath"
)

// Target identifies an AI coding assistant that supports hook configuration.
type Target string

const (
	TargetClaude Target = "claude"
	TargetCodex  Target = "codex"
	TargetCursor Target = "cursor"
	TargetGemini Target = "gemini"
)

// ConfigPath returns the absolute path to the hook config file for this target.
func (t Target) ConfigPath() string {
	home, _ := os.UserHomeDir()
	switch t {
	case TargetClaude:
		return filepath.Join(home, ".claude", "settings.json")
	case TargetCodex:
		return filepath.Join(home, ".codex", "config.json")
	case TargetCursor:
		return filepath.Join(home, ".cursor", "mcp.json")
	case TargetGemini:
		return filepath.Join(home, ".gemini", "settings.json")
	default:
		return ""
	}
}

// allTargets lists every supported target for validation.
var allTargets = []Target{TargetClaude, TargetCodex, TargetCursor, TargetGemini}

func parseTarget(s string) (Target, bool) {
	for _, t := range allTargets {
		if string(t) == s {
			return t, true
		}
	}
	return "", false
}
