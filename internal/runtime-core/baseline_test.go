package runtimecore

import (
	"errors"
	"log/slog"
	"testing"
)

func TestGoBaselineSupportsRuntimeCorePrerequisites(t *testing.T) {
	values := mapSlice([]int{1, 2, 3}, func(v int) int {
		return v * 2
	})
	if got, want := values[2], 6; got != want {
		t.Fatalf("generic helper result = %d, want %d", got, want)
	}

	joined := errors.Join(errors.New("provider failed"), errors.New("store failed"))
	if joined == nil {
		t.Fatal("errors.Join returned nil for non-nil errors")
	}

	attrs := []slog.Attr{slog.String("module", "runtime-core")}
	if got, want := attrs[0].Key, "module"; got != want {
		t.Fatalf("slog attr key = %q, want %q", got, want)
	}
}

func mapSlice[T any, U any](items []T, fn func(T) U) []U {
	out := make([]U, 0, len(items))
	for _, item := range items {
		out = append(out, fn(item))
	}
	return out
}
