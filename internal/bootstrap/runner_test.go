package bootstrap

import (
	"bytes"
	"context"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"
)

func TestRunSync_WritesSummary(t *testing.T) {
	t.Parallel()
	cfg := Config{Root: sampleWorkspace(t), Version: "test"}
	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
	DoInjections(inj, cfg)

	var buf bytes.Buffer
	if err := RunSync(context.Background(), inj, &buf); err != nil {
		t.Fatalf("run sync: %v", err)
	}

	if buf.Len() == 0 {
		t.Fatal("expected non-empty summary")
	}
}
