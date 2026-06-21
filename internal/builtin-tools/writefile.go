package builtin

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	toolruntime "forgecode/internal/tool-runtime"
)

type writeFileTool struct {
	descriptor toolruntime.ToolDescriptor
	deps       Deps
}

type WriteFileInput struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func (t writeFileTool) Descriptor() toolruntime.ToolDescriptor {
	return t.descriptor
}

func (t writeFileTool) Execute(ctx context.Context, raw json.RawMessage) (toolruntime.ToolResult, error) {
	var input WriteFileInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.ValidationError, "invalid WriteFile input", err)
	}
	target, err := resolveWorkspacePath(t.deps.WorkspaceRoot, input.Path)
	if err != nil {
		return toolruntime.ToolResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.CancelledError, "WriteFile cancelled", err)
	}

	exists := false
	if info, err := os.Stat(target); err == nil {
		if info.IsDir() {
			return toolruntime.ToolResult{}, toolruntime.NewError(toolruntime.ValidationError, "WriteFile path is a directory")
		}
		exists = true
	} else if !os.IsNotExist(err) {
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.ToolExecutionError, "stat target file", err)
	}

	checkpointID := ""
	if exists {
		checkpointID, err = checkpointBeforeWrite(ctx, t.deps.Checkpointer, "pre-write WriteFile "+input.Path)
		if err != nil {
			return toolruntime.ToolResult{}, err
		}
	}
	if err := writeFileAtomically(ctx, target, []byte(input.Content)); err != nil {
		return toolruntime.ToolResult{}, err
	}

	meta := map[string]any{
		"path":    input.Path,
		"bytes":   len(input.Content),
		"created": !exists,
	}
	if checkpointID != "" {
		meta["checkpoint_id"] = checkpointID
	}
	return toolruntime.ToolResult{
		Output: "wrote " + strconv.Itoa(len(input.Content)) + " bytes to " + input.Path,
		Meta:   meta,
	}, nil
}
