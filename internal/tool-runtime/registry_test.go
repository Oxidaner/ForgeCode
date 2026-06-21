package toolruntime

import (
	"context"
	"encoding/json"
	"testing"
)

type fakeTool struct {
	descriptor ToolDescriptor
}

func (f fakeTool) Descriptor() ToolDescriptor {
	return f.descriptor
}

func (f fakeTool) Execute(ctx context.Context, input json.RawMessage) (ToolResult, error) {
	return ToolResult{Output: "ok"}, nil
}

func TestRegistryRegistersListsAndGetsTools(t *testing.T) {
	reg := NewRegistry()
	tool := fakeTool{descriptor: ToolDescriptor{
		Name:        "ReadFile",
		Source:      ToolSourceBuiltin,
		InputSchema: json.RawMessage(`{"type":"object"}`),
		Risk:        RiskLow,
		Permission:  PermissionHint{Actions: []PermissionAction{PermissionRead}},
	}}

	if err := reg.Register(tool); err != nil {
		t.Fatal(err)
	}
	if _, ok := reg.Get("ReadFile"); !ok {
		t.Fatal("expected registered tool to be discoverable")
	}
	list := reg.List()
	if len(list) != 1 || list[0].Name != "ReadFile" || list[0].Source != ToolSourceBuiltin {
		t.Fatalf("unexpected list: %#v", list)
	}
}

func TestRegistryRejectsConflictingToolNames(t *testing.T) {
	reg := NewRegistry()
	tool := fakeTool{descriptor: ToolDescriptor{
		Name:        "ReadFile",
		Source:      ToolSourceBuiltin,
		InputSchema: json.RawMessage(`{"type":"object"}`),
		Risk:        RiskLow,
	}}

	if err := reg.Register(tool); err != nil {
		t.Fatal(err)
	}
	err := reg.Register(tool)
	if !IsCategory(err, ConflictError) {
		t.Fatalf("expected ConflictError, got %v", err)
	}
}

func TestDescriptorRequiresValidSchema(t *testing.T) {
	reg := NewRegistry()
	err := reg.Register(fakeTool{descriptor: ToolDescriptor{
		Name:        "bad",
		Source:      ToolSourceBuiltin,
		InputSchema: json.RawMessage(`not-json`),
		Risk:        RiskLow,
	}})
	if !IsCategory(err, ValidationError) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}
