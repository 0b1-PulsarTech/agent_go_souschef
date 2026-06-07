package hooksetup

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// claudePreToolUse is the JSON payload Claude Code sends to a PreToolUse hook.
type claudePreToolUse struct {
	ToolName  string          `json:"tool_name"`
	ToolInput json.RawMessage `json:"tool_input"`
}

// hookDecision is what the hook writes to stdout to control Claude Code.
type hookDecision struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason,omitempty"`
}

// runHook dispatches to the correct hook handler based on target name.
func runHook(target string) int {
	switch target {
	case string(TargetClaude):
		if err := RunClaude(os.Stdin, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	default:
		fmt.Fprintf(os.Stderr, "hook run: unsupported target %q\n", target)
		return 1
	}
}

// RunClaude handles a single Claude Code PreToolUse hook invocation.
// It reads a JSON event from r and writes a decision to w.
// Unrecognised tool calls receive {"decision":"continue"} (pass-through).
func RunClaude(r io.Reader, w io.Writer) error {
	var event claudePreToolUse
	if err := json.NewDecoder(r).Decode(&event); err != nil {
		// Malformed input — pass through so Claude Code continues unblocked
		return writeDecision(w, "continue", "")
	}
	return writeDecision(w, "continue", "")
}

func writeDecision(w io.Writer, decision, reason string) error {
	return json.NewEncoder(w).Encode(hookDecision{Decision: decision, Reason: reason})
}
