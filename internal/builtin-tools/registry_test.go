package builtin

import (
	"encoding/json"
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
