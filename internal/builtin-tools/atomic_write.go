package builtin

import (
	"context"
	"os"
	"path/filepath"

	toolruntime "forgecode/internal/tool-runtime"
)

func writeFileAtomically(ctx context.Context, target string, data []byte) error {
	if err := ctx.Err(); err != nil {
		return toolruntime.WrapError(toolruntime.CancelledError, "write cancelled", err)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
		return toolruntime.WrapError(toolruntime.ToolExecutionError, "create parent directory", err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(target), "."+filepath.Base(target)+".tmp-*")
	if err != nil {
		return toolruntime.WrapError(toolruntime.ToolExecutionError, "create temporary file", err)
	}
	tmpName := tmp.Name()
	closed := false
	defer func() {
		if !closed {
			_ = tmp.Close()
		}
		_ = os.Remove(tmpName)
	}()

	if _, err := tmp.Write(data); err != nil {
		return toolruntime.WrapError(toolruntime.ToolExecutionError, "write temporary file", err)
	}
	if err := tmp.Sync(); err != nil {
		return toolruntime.WrapError(toolruntime.ToolExecutionError, "sync temporary file", err)
	}
	if err := tmp.Close(); err != nil {
		closed = true
		return toolruntime.WrapError(toolruntime.ToolExecutionError, "close temporary file", err)
	}
	closed = true
	if err := ctx.Err(); err != nil {
		return toolruntime.WrapError(toolruntime.CancelledError, "write cancelled", err)
	}
	if err := os.Rename(tmpName, target); err != nil {
		return toolruntime.WrapError(toolruntime.ToolExecutionError, "replace target file", err)
	}
	return nil
}

func checkpointBeforeWrite(ctx context.Context, cp Checkpointer, reason string) (string, error) {
	if cp == nil {
		return "", toolruntime.NewError(toolruntime.PersistenceError, "checkpoint is required before overwriting existing file")
	}
	id, err := cp.CreateCheckpoint(ctx, "", reason)
	if err != nil {
		return "", toolruntime.WrapError(toolruntime.PersistenceError, "create checkpoint", err)
	}
	return id, nil
}
