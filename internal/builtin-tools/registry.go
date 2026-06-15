package builtin

import (
	"context"
	"encoding/json"

	toolruntime "forgecode/internal/tool-runtime"
)

func RegisterBuiltins(reg toolruntime.Registry, deps Deps) error {
	deps = deps.withDefaults()
	for _, descriptor := range Descriptors() {
		if err := reg.Register(newPlaceholderTool(descriptor, deps)); err != nil {
			return err
		}
	}
	return nil
}

func NewReadFileTool(d Deps) toolruntime.Tool {
	return newPlaceholderTool(Descriptors()[0], d.withDefaults())
}

func NewWriteFileTool(d Deps) toolruntime.Tool {
	return newPlaceholderTool(Descriptors()[1], d.withDefaults())
}

func NewEditFileTool(d Deps) toolruntime.Tool {
	return newPlaceholderTool(Descriptors()[2], d.withDefaults())
}

func NewBashTool(d Deps) toolruntime.Tool {
	return newPlaceholderTool(Descriptors()[3], d.withDefaults())
}

func NewGlobTool(d Deps) toolruntime.Tool {
	return newPlaceholderTool(Descriptors()[4], d.withDefaults())
}

func NewGrepTool(d Deps) toolruntime.Tool {
	return newPlaceholderTool(Descriptors()[5], d.withDefaults())
}

type placeholderTool struct {
	descriptor toolruntime.ToolDescriptor
	deps       Deps
}

func newPlaceholderTool(descriptor toolruntime.ToolDescriptor, deps Deps) toolruntime.Tool {
	return placeholderTool{descriptor: descriptor, deps: deps}
}

func (t placeholderTool) Descriptor() toolruntime.ToolDescriptor {
	return t.descriptor
}

func (t placeholderTool) Execute(ctx context.Context, input json.RawMessage) (toolruntime.ToolResult, error) {
	if err := ctx.Err(); err != nil {
		return toolruntime.ToolResult{}, err
	}
	return toolruntime.ToolResult{
		Output:   t.descriptor.Name + " execution is not implemented in FC-BT-001",
		IsError:  true,
		Category: toolruntime.ToolExecutionError,
		Meta: map[string]any{
			"workspace_root": t.deps.WorkspaceRoot,
			"task":           "FC-BT-001",
		},
	}, nil
}
