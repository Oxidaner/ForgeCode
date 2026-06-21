package toolruntime

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestInvokerRunsPipelineInOrder(t *testing.T) {
	var order []string
	registry := NewRegistry()
	if err := registry.Register(recordingTool{
		name:  "Echo",
		order: &order,
		result: ToolResult{
			Output: "hello",
		},
	}); err != nil {
		t.Fatal(err)
	}

	audit := &recordingAudit{order: &order}
	invoker := NewInvoker(InvokerConfig{
		Registry:   registry,
		Permission: &recordingPermission{effect: PermissionAllow, order: &order},
		Hooks:      []Hook{&recordingHook{order: &order}},
		Audit:      audit,
	})

	result, err := invoker.Invoke(context.Background(), ToolCall{
		ID:    "call-1",
		Name:  "Echo",
		Input: json.RawMessage(`{"message":"hello"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.CallID != "call-1" || result.Output != "hello" || result.IsError {
		t.Fatalf("result = %#v", result)
	}
	want := []string{"permission", "pre-hook", "execute", "post-hook", "audit"}
	assertOrder(t, order, want)
	if len(audit.records) != 1 || audit.records[0].Stage != AuditStageSucceeded {
		t.Fatalf("audit records = %#v", audit.records)
	}
}

func TestInvokerRejectsInvalidInputBeforePermissionAndAudits(t *testing.T) {
	var order []string
	registry := NewRegistry()
	if err := registry.Register(recordingTool{name: "Echo", order: &order}); err != nil {
		t.Fatal(err)
	}
	permission := &recordingPermission{effect: PermissionAllow, order: &order}
	audit := &recordingAudit{order: &order}
	invoker := NewInvoker(InvokerConfig{Registry: registry, Permission: permission, Audit: audit})

	result, err := invoker.Invoke(context.Background(), ToolCall{
		ID:    "call-1",
		Name:  "Echo",
		Input: json.RawMessage(`{}`),
	})
	if !IsCategory(err, ValidationError) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if !result.IsError || result.Category != ValidationError {
		t.Fatalf("result = %#v", result)
	}
	if permission.calls != 0 {
		t.Fatalf("permission was called %d times", permission.calls)
	}
	assertOrder(t, order, []string{"audit"})
}

func TestInvokerDeniesBeforeExecuteAndAudits(t *testing.T) {
	var order []string
	registry := NewRegistry()
	if err := registry.Register(recordingTool{name: "Echo", order: &order}); err != nil {
		t.Fatal(err)
	}
	audit := &recordingAudit{order: &order}
	invoker := NewInvoker(InvokerConfig{
		Registry:   registry,
		Permission: &recordingPermission{effect: PermissionDeny, order: &order},
		Hooks:      []Hook{&recordingHook{order: &order}},
		Audit:      audit,
	})

	result, err := invoker.Invoke(context.Background(), ToolCall{
		ID:    "call-1",
		Name:  "Echo",
		Input: json.RawMessage(`{"message":"hello"}`),
	})
	if !IsCategory(err, PermissionDenied) {
		t.Fatalf("expected PermissionDenied, got %v", err)
	}
	if !result.IsError || result.Category != PermissionDenied {
		t.Fatalf("result = %#v", result)
	}
	assertOrder(t, order, []string{"permission", "audit"})
	if len(audit.records) != 1 || audit.records[0].Stage != AuditStageRejected {
		t.Fatalf("audit records = %#v", audit.records)
	}
}

func TestInvokerReturnsApprovalRequiredBeforeExecute(t *testing.T) {
	var order []string
	registry := NewRegistry()
	if err := registry.Register(recordingTool{name: "Echo", order: &order}); err != nil {
		t.Fatal(err)
	}
	invoker := NewInvoker(InvokerConfig{
		Registry:   registry,
		Permission: &recordingPermission{effect: PermissionAskAlways, order: &order},
		Audit:      &recordingAudit{order: &order},
	})

	result, err := invoker.Invoke(context.Background(), ToolCall{
		ID:    "call-1",
		Name:  "Echo",
		Input: json.RawMessage(`{"message":"hello"}`),
	})
	if !IsCategory(err, ApprovalRequired) {
		t.Fatalf("expected ApprovalRequired, got %v", err)
	}
	if !result.IsError || result.Category != ApprovalRequired {
		t.Fatalf("result = %#v", result)
	}
	assertOrder(t, order, []string{"permission", "audit"})
}

func TestInvokerTruncatesOversizedOutput(t *testing.T) {
	registry := NewRegistry()
	if err := registry.Register(recordingTool{
		name:   "Echo",
		result: ToolResult{Output: "0123456789abcdef"},
	}); err != nil {
		t.Fatal(err)
	}
	invoker := NewInvoker(InvokerConfig{
		Registry:       registry,
		Permission:     &recordingPermission{effect: PermissionAllow},
		Audit:          &recordingAudit{},
		MaxOutputBytes: 8,
	})

	result, err := invoker.Invoke(context.Background(), ToolCall{
		ID:    "call-1",
		Name:  "Echo",
		Input: json.RawMessage(`{"message":"hello"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Truncated || result.Output != "01234567" {
		t.Fatalf("result = %#v, want truncated output", result)
	}
}

func TestInvokerClassifiesTimeout(t *testing.T) {
	registry := NewRegistry()
	if err := registry.Register(blockingTool{name: "Slow"}); err != nil {
		t.Fatal(err)
	}
	invoker := NewInvoker(InvokerConfig{
		Registry:       registry,
		Permission:     &recordingPermission{effect: PermissionAllow},
		Audit:          &recordingAudit{},
		DefaultTimeout: 10 * time.Millisecond,
	})

	result, err := invoker.Invoke(context.Background(), ToolCall{
		ID:    "call-1",
		Name:  "Slow",
		Input: json.RawMessage(`{"message":"hello"}`),
	})
	if !IsCategory(err, TimeoutError) {
		t.Fatalf("expected TimeoutError, got %v", err)
	}
	if !result.IsError || result.Category != TimeoutError {
		t.Fatalf("result = %#v", result)
	}
}

func TestInvokerClassifiesCancellation(t *testing.T) {
	registry := NewRegistry()
	if err := registry.Register(blockingTool{name: "Slow"}); err != nil {
		t.Fatal(err)
	}
	invoker := NewInvoker(InvokerConfig{
		Registry:       registry,
		Permission:     &recordingPermission{effect: PermissionAllow},
		Audit:          &recordingAudit{},
		DefaultTimeout: time.Second,
	})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	result, err := invoker.Invoke(ctx, ToolCall{
		ID:    "call-1",
		Name:  "Slow",
		Input: json.RawMessage(`{"message":"hello"}`),
	})
	if !IsCategory(err, CancelledError) {
		t.Fatalf("expected CancelledError, got %v", err)
	}
	if !result.IsError || result.Category != CancelledError {
		t.Fatalf("result = %#v", result)
	}
}

type recordingPermission struct {
	effect PermissionEffect
	order  *[]string
	calls  int
}

func (p *recordingPermission) DecideTool(ctx context.Context, req PermissionRequest) (PermissionDecision, error) {
	if err := ctx.Err(); err != nil {
		return PermissionDecision{}, err
	}
	p.calls++
	appendOrder(p.order, "permission")
	return PermissionDecision{Effect: p.effect, Risk: req.Descriptor.Risk}, nil
}

type recordingHook struct {
	order *[]string
	err   error
}

func (h *recordingHook) HandleToolEvent(ctx context.Context, event ToolEvent) (HookDecision, error) {
	if err := ctx.Err(); err != nil {
		return HookDecision{}, err
	}
	if h.err != nil {
		return HookDecision{}, h.err
	}
	switch event.Stage {
	case HookStagePreToolUse:
		appendOrder(h.order, "pre-hook")
	case HookStagePostToolUse:
		appendOrder(h.order, "post-hook")
	}
	return HookDecision{Effect: PermissionAllow}, nil
}

type recordingAudit struct {
	order   *[]string
	records []ToolAuditRecord
	err     error
}

func (a *recordingAudit) RecordToolAudit(ctx context.Context, record ToolAuditRecord) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if a.err != nil {
		return a.err
	}
	appendOrder(a.order, "audit")
	a.records = append(a.records, record)
	return nil
}

type recordingTool struct {
	name   string
	order  *[]string
	result ToolResult
	err    error
}

func (t recordingTool) Descriptor() ToolDescriptor {
	return ToolDescriptor{
		Name:        t.name,
		Source:      ToolSourceBuiltin,
		InputSchema: json.RawMessage(`{"type":"object","required":["message"],"properties":{"message":{"type":"string","minLength":1}},"additionalProperties":false}`),
		Risk:        RiskLow,
		Permission:  PermissionHint{Actions: []PermissionAction{PermissionRead}},
	}
}

func (t recordingTool) Execute(ctx context.Context, input json.RawMessage) (ToolResult, error) {
	if err := ctx.Err(); err != nil {
		return ToolResult{}, err
	}
	appendOrder(t.order, "execute")
	if t.err != nil {
		return t.result, t.err
	}
	return t.result, nil
}

type blockingTool struct {
	name string
}

func (t blockingTool) Descriptor() ToolDescriptor {
	return recordingTool{name: t.name}.Descriptor()
}

func (t blockingTool) Execute(ctx context.Context, input json.RawMessage) (ToolResult, error) {
	<-ctx.Done()
	return ToolResult{}, ctx.Err()
}

func appendOrder(order *[]string, value string) {
	if order != nil {
		*order = append(*order, value)
	}
}

func assertOrder(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("order = %#v, want %#v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("order = %#v, want %#v", got, want)
		}
	}
}
