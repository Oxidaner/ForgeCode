package builtin

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	toolruntime "forgecode/internal/tool-runtime"
)

func TestRegisterBuiltinsRegistersSixUniqueTools(t *testing.T) {
	reg := toolruntime.NewRegistry()
	if err := RegisterBuiltins(reg, Deps{WorkspaceRoot: t.TempDir()}); err != nil {
		t.Fatal(err)
	}

	descriptors := reg.List()
	if len(descriptors) != 6 {
		t.Fatalf("expected six built-in tools, got %d", len(descriptors))
	}

	seen := map[string]bool{}
	for _, descriptor := range descriptors {
		if seen[descriptor.Name] {
			t.Fatalf("duplicate descriptor: %s", descriptor.Name)
		}
		seen[descriptor.Name] = true
		if descriptor.Source != toolruntime.ToolSourceBuiltin {
			t.Fatalf("%s source = %s, want builtin", descriptor.Name, descriptor.Source)
		}
		if !json.Valid(descriptor.InputSchema) {
			t.Fatalf("%s schema is invalid JSON", descriptor.Name)
		}
	}
	for _, name := range []string{ToolReadFile, ToolWriteFile, ToolEditFile, ToolBash, ToolGlob, ToolGrep} {
		if !seen[name] {
			t.Fatalf("missing built-in tool %s", name)
		}
	}
}

func TestBashDescriptorDefaultsToHighRisk(t *testing.T) {
	for _, descriptor := range Descriptors() {
		if descriptor.Name == ToolBash {
			if descriptor.Risk != toolruntime.RiskHigh {
				t.Fatalf("Bash risk = %s, want High", descriptor.Risk)
			}
			return
		}
	}
	t.Fatal("Bash descriptor not found")
}

func TestRegisterBuiltinsDetectsConflicts(t *testing.T) {
	reg := toolruntime.NewRegistry()
	if err := RegisterBuiltins(reg, Deps{}); err != nil {
		t.Fatal(err)
	}
	err := RegisterBuiltins(reg, Deps{})
	if !toolruntime.IsCategory(err, toolruntime.ConflictError) {
		t.Fatalf("expected ConflictError, got %v", err)
	}
}

func TestRegisteredWriteAndEditToolsExecuteRealImplementations(t *testing.T) {
	root := t.TempDir()
	reg := toolruntime.NewRegistry()
	if err := RegisterBuiltins(reg, Deps{WorkspaceRoot: root, Checkpointer: &fakeCheckpointer{id: "cp"}}); err != nil {
		t.Fatal(err)
	}

	writeTool, ok := reg.Get(ToolWriteFile)
	if !ok {
		t.Fatal("WriteFile not registered")
	}
	if result, err := writeTool.Execute(context.Background(), mustJSON(t, WriteFileInput{Path: "notes.txt", Content: "old\n"})); err != nil {
		t.Fatal(err)
	} else if strings.Contains(result.Output, "not implemented") {
		t.Fatalf("WriteFile still uses placeholder: %#v", result)
	}

	editTool, ok := reg.Get(ToolEditFile)
	if !ok {
		t.Fatal("EditFile not registered")
	}
	result, err := editTool.Execute(context.Background(), mustJSON(t, EditFileInput{Path: "notes.txt", OldString: "old", NewString: "new"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.Output, "-old") || !strings.Contains(result.Output, "+new") {
		t.Fatalf("unexpected EditFile output: %q", result.Output)
	}
}
