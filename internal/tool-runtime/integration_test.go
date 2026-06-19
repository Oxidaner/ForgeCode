package toolruntime_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	builtin "forgecode/internal/builtin-tools"
	permission "forgecode/internal/permission-engine"
	toolruntime "forgecode/internal/tool-runtime"
)

func TestInvokerRunsBuiltinToolThroughPermissionEngine(t *testing.T) {
	root := t.TempDir()
	writeIntegrationFile(t, root, "notes.txt", "hello\n")
	registry := toolruntime.NewRegistry()
	if err := builtin.RegisterBuiltins(registry, builtin.Deps{WorkspaceRoot: root}); err != nil {
		t.Fatal(err)
	}
	audit := &integrationAudit{}
	invoker := toolruntime.NewInvoker(toolruntime.InvokerConfig{
		Registry:   registry,
		Permission: permission.NewPolicyDecider(permission.PolicyConfig{WorkspaceRoot: root}),
		Audit:      audit,
	})

	result, err := invoker.Invoke(context.Background(), toolruntime.ToolCall{
		ID:    "call-read",
		Name:  builtin.ToolReadFile,
		Input: []byte(`{"path":"notes.txt"}`),
		Ctx:   toolruntime.InvocationContext{WorkspaceRoot: root},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.Output, "1: hello") {
		t.Fatalf("unexpected output: %q", result.Output)
	}
	if len(audit.records) != 1 || audit.records[0].Stage != toolruntime.AuditStageSucceeded {
		t.Fatalf("audit records = %#v", audit.records)
	}
}

func TestInvokerBlocksBuiltinToolWhenPermissionDenies(t *testing.T) {
	root := t.TempDir()
	registry := toolruntime.NewRegistry()
	if err := builtin.RegisterBuiltins(registry, builtin.Deps{WorkspaceRoot: root}); err != nil {
		t.Fatal(err)
	}
	invoker := toolruntime.NewInvoker(toolruntime.InvokerConfig{
		Registry:   registry,
		Permission: permission.NewPolicyDecider(permission.PolicyConfig{WorkspaceRoot: root}),
		Audit:      &integrationAudit{},
	})

	result, err := invoker.Invoke(context.Background(), toolruntime.ToolCall{
		ID:    "call-escape",
		Name:  builtin.ToolReadFile,
		Input: []byte(`{"path":"../outside.txt"}`),
		Ctx:   toolruntime.InvocationContext{WorkspaceRoot: root},
	})
	if !toolruntime.IsCategory(err, toolruntime.PermissionDenied) {
		t.Fatalf("expected PermissionDenied, got %v", err)
	}
	if !result.IsError || result.Category != toolruntime.PermissionDenied {
		t.Fatalf("result = %#v", result)
	}
}

func TestInvokerRequiresApprovalForHighRiskBuiltinBash(t *testing.T) {
	root := t.TempDir()
	registry := toolruntime.NewRegistry()
	if err := builtin.RegisterBuiltins(registry, builtin.Deps{WorkspaceRoot: root}); err != nil {
		t.Fatal(err)
	}
	invoker := toolruntime.NewInvoker(toolruntime.InvokerConfig{
		Registry:   registry,
		Permission: permission.NewPolicyDecider(permission.PolicyConfig{WorkspaceRoot: root}),
		Audit:      &integrationAudit{},
	})

	result, err := invoker.Invoke(context.Background(), toolruntime.ToolCall{
		ID:    "call-bash",
		Name:  builtin.ToolBash,
		Input: []byte(`{"command":"echo hello"}`),
		Ctx:   toolruntime.InvocationContext{WorkspaceRoot: root},
	})
	if !toolruntime.IsCategory(err, toolruntime.ApprovalRequired) {
		t.Fatalf("expected ApprovalRequired, got %v", err)
	}
	if !result.IsError || result.Category != toolruntime.ApprovalRequired {
		t.Fatalf("result = %#v", result)
	}
}

type integrationAudit struct {
	records []toolruntime.ToolAuditRecord
}

func (a *integrationAudit) RecordToolAudit(ctx context.Context, record toolruntime.ToolAuditRecord) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	a.records = append(a.records, record)
	return nil
}

func writeIntegrationFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
