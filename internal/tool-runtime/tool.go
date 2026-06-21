package toolruntime

import (
	"context"
	"encoding/json"
)

type Tool interface {
	Descriptor() ToolDescriptor
	Execute(ctx context.Context, input json.RawMessage) (ToolResult, error)
}

type Registry interface {
	Register(t Tool) error
	Get(name string) (Tool, bool)
	List() []ToolDescriptor
}

type Invoker interface {
	Invoke(ctx context.Context, call ToolCall) (ToolResult, error)
}
