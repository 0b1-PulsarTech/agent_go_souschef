package mcpsvc

// JSON IO types for the MCP tools. The field tags drive the published
// jsonschema the LLM sees.
type (
	QueryIn struct {
		Query  string `json:"query"            jsonschema:"Symbol name or free-text query"`
		Expand bool   `json:"expand,omitempty" jsonschema:"Return transitive callers/callees"`
	}
	SyncIn   struct{}
	SourceIn struct {
		Query string `json:"query" jsonschema:"Symbol name to locate"`
	}
	ChangedIn struct {
		Scope string `json:"scope,omitempty" jsonschema:"Optional path filter"`
	}
	ShadowsIn struct {
		Scope string `json:"scope,omitempty" jsonschema:"Optional path filter to narrow findings"`
	}
	Result struct {
		Text string `json:"text" jsonschema:"Compact human/LLM-readable result"`
	}
)
