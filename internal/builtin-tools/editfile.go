package builtin

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	toolruntime "forgecode/internal/tool-runtime"
)

type editFileTool struct {
	descriptor toolruntime.ToolDescriptor
	deps       Deps
}

type EditFileInput struct {
	Path      string `json:"path"`
	OldString string `json:"old_string"`
	NewString string `json:"new_string"`
}

func (t editFileTool) Descriptor() toolruntime.ToolDescriptor {
	return t.descriptor
}

func (t editFileTool) Execute(ctx context.Context, raw json.RawMessage) (toolruntime.ToolResult, error) {
	var input EditFileInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.ValidationError, "invalid EditFile input", err)
	}
	if input.OldString == "" {
		return toolruntime.ToolResult{}, toolruntime.NewError(toolruntime.ValidationError, "old_string is required")
	}
	target, err := resolveWorkspacePath(t.deps.WorkspaceRoot, input.Path)
	if err != nil {
		return toolruntime.ToolResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.CancelledError, "EditFile cancelled", err)
	}
	info, err := os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return toolruntime.ToolResult{}, toolruntime.NewError(toolruntime.ValidationError, "EditFile target does not exist")
		}
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.ToolExecutionError, "stat target file", err)
	}
	if info.IsDir() {
		return toolruntime.ToolResult{}, toolruntime.NewError(toolruntime.ValidationError, "EditFile path is a directory")
	}

	originalBytes, err := os.ReadFile(target)
	if err != nil {
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.ToolExecutionError, "read target file", err)
	}
	if err := ctx.Err(); err != nil {
		return toolruntime.ToolResult{}, toolruntime.WrapError(toolruntime.CancelledError, "EditFile cancelled", err)
	}
	original := string(originalBytes)
	matches := strings.Count(original, input.OldString)
	if matches != 1 {
		return toolruntime.ToolResult{}, toolruntime.NewError(toolruntime.ValidationError, "old_string must match exactly once; got "+strconv.Itoa(matches))
	}

	checkpointID, err := checkpointBeforeWrite(ctx, t.deps.Checkpointer, "pre-edit EditFile "+input.Path)
	if err != nil {
		return toolruntime.ToolResult{}, err
	}
	updated := strings.Replace(original, input.OldString, input.NewString, 1)
	if err := writeFileAtomically(ctx, target, []byte(updated)); err != nil {
		return toolruntime.ToolResult{}, err
	}

	diff := unifiedDiff(input.Path, original, updated)
	return toolruntime.ToolResult{
		Output: diff,
		Meta: map[string]any{
			"path":          input.Path,
			"checkpoint_id": checkpointID,
			"replacements":  1,
			"diff":          diff,
		},
	}, nil
}

func unifiedDiff(path, before, after string) string {
	var builder strings.Builder
	builder.WriteString("--- " + path + "\n")
	builder.WriteString("+++ " + path + "\n")
	beforeLines := splitDiffLines(before)
	afterLines := splitDiffLines(after)
	max := len(beforeLines)
	if len(afterLines) > max {
		max = len(afterLines)
	}
	for i := 0; i < max; i++ {
		var oldLine, newLine string
		if i < len(beforeLines) {
			oldLine = beforeLines[i]
		}
		if i < len(afterLines) {
			newLine = afterLines[i]
		}
		switch {
		case i >= len(beforeLines):
			builder.WriteString("+" + newLine + "\n")
		case i >= len(afterLines):
			builder.WriteString("-" + oldLine + "\n")
		case oldLine == newLine:
			builder.WriteString(" " + oldLine + "\n")
		default:
			builder.WriteString("-" + oldLine + "\n")
			builder.WriteString("+" + newLine + "\n")
		}
	}
	return builder.String()
}

func splitDiffLines(text string) []string {
	text = strings.TrimSuffix(text, "\n")
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}
