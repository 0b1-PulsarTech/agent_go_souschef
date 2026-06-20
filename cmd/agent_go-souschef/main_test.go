package main

import (
	"context"
	"os"
	"testing"
)

func TestRun_NoArgs(t *testing.T) {
	t.Parallel()

	code := run(context.Background(), nil, os.Stdout, os.Stderr)
	if code != 1 {
		t.Fatalf("expected exit 1 on empty args, got %d", code)
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	t.Parallel()

	code := run(context.Background(), []string{"definitely-not-real"}, os.Stdout, os.Stderr)
	if code != 1 {
		t.Fatalf("expected exit 1 on unknown command, got %d", code)
	}
}
