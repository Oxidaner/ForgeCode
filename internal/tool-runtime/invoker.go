package toolruntime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

type PermissionEffect string

const (
	PermissionAllow     PermissionEffect = "Allow"
	PermissionAskOnce   PermissionEffect = "AskOnce"
	PermissionAskAlways PermissionEffect = "AskAlways"
	PermissionDeny      PermissionEffect = "Deny"
)

type PermissionReason struct {
	Layer  string
	RuleID string
	Reason string
}

type PermissionRequest struct {
	Descriptor ToolDescriptor
	Input      json.RawMessage
	Inv        InvocationContext
}

type PermissionDecision struct {
	Effect  PermissionEffect
	Risk    RiskLevel
	Reasons []PermissionReason
}

type PermissionChecker interface {
	DecideTool(ctx context.Context, req PermissionRequest) (PermissionDecision, error)
}

type HookStage string

const (
	HookStagePreToolUse  HookStage = "PreToolUse"
	HookStagePostToolUse HookStage = "PostToolUse"
	HookStageToolFailure HookStage = "ToolFailure"
)

type ToolEvent struct {
	Stage      HookStage
	Call       ToolCall
	Descriptor ToolDescriptor
	Decision   PermissionDecision
	Result     ToolResult
	Err        error
}

type HookDecision struct {
	Effect PermissionEffect
	Reason string
}

type Hook interface {
	HandleToolEvent(ctx context.Context, event ToolEvent) (HookDecision, error)
}

type AuditStage string

const (
	AuditStageSucceeded AuditStage = "Succeeded"
	AuditStageRejected  AuditStage = "Rejected"
	AuditStageFailed    AuditStage = "Failed"
)

type ToolAuditRecord struct {
	CallID     string
	ToolName   string
	Source     ToolSource
	Stage      AuditStage
	Category   ErrorCategory
	Input      json.RawMessage
	Decision   PermissionDecision
	Result     ToolResult
	Error      string
	OccurredAt time.Time
}

type AuditSink interface {
	RecordToolAudit(ctx context.Context, record ToolAuditRecord) error
}

type InvokerConfig struct {
	Registry       Registry
	Permission     PermissionChecker
	Hooks          []Hook
	Audit          AuditSink
	DefaultTimeout time.Duration
	MaxOutputBytes int
	Now            func() time.Time
}

type PipelineInvoker struct {
	config InvokerConfig
}

func NewInvoker(config InvokerConfig) *PipelineInvoker {
	if config.DefaultTimeout == 0 {
		config.DefaultTimeout = 60 * time.Second
	}
	if config.MaxOutputBytes == 0 {
		config.MaxOutputBytes = 100 << 10
	}
	if config.Now == nil {
		config.Now = func() time.Time { return time.Now().UTC() }
	}
	return &PipelineInvoker{config: config}
}

func (i *PipelineInvoker) Invoke(ctx context.Context, call ToolCall) (ToolResult, error) {
	if err := ctx.Err(); err != nil {
		return ToolResult{}, err
	}
	result := ToolResult{CallID: call.ID}
	if i.config.Registry == nil {
		err := NewError(ValidationError, "tool registry is required")
		return i.finish(ctx, call, ToolDescriptor{}, PermissionDecision{}, result, err, AuditStageRejected)
	}

	tool, ok := i.config.Registry.Get(call.Name)
	if !ok {
		err := NewError(ValidationError, "tool not registered: "+call.Name)
		return i.finish(ctx, call, ToolDescriptor{Name: call.Name}, PermissionDecision{}, result, err, AuditStageRejected)
	}
	descriptor := tool.Descriptor()
	if err := validateInputAgainstSchema(descriptor.InputSchema, call.Input); err != nil {
		return i.finish(ctx, call, descriptor, PermissionDecision{}, result, err, AuditStageRejected)
	}
	if i.config.Permission == nil {
		err := NewError(PermissionDenied, "permission checker is required")
		return i.finish(ctx, call, descriptor, PermissionDecision{}, result, err, AuditStageRejected)
	}

	decision, err := i.config.Permission.DecideTool(ctx, PermissionRequest{
		Descriptor: descriptor,
		Input:      append(json.RawMessage{}, call.Input...),
		Inv:        call.Ctx,
	})
	if err != nil {
		return i.finish(ctx, call, descriptor, decision, result, err, AuditStageFailed)
	}
	switch decision.Effect {
	case PermissionDeny:
		err := NewError(PermissionDenied, "tool invocation denied by permission decision")
		return i.finish(ctx, call, descriptor, decision, result, err, AuditStageRejected)
	case PermissionAskOnce, PermissionAskAlways:
		err := NewError(ApprovalRequired, "tool invocation requires approval")
		return i.finish(ctx, call, descriptor, decision, result, err, AuditStageRejected)
	case "", PermissionAllow:
	default:
		err := NewError(PermissionDenied, "unknown permission decision effect: "+string(decision.Effect))
		return i.finish(ctx, call, descriptor, decision, result, err, AuditStageRejected)
	}

	if err := i.runHooks(ctx, ToolEvent{Stage: HookStagePreToolUse, Call: call, Descriptor: descriptor, Decision: decision}); err != nil {
		return i.finish(ctx, call, descriptor, decision, result, err, AuditStageRejected)
	}

	execCtx, cancel := context.WithTimeout(ctx, i.config.DefaultTimeout)
	defer cancel()
	result, err = tool.Execute(execCtx, append(json.RawMessage{}, call.Input...))
	result.CallID = call.ID
	if err != nil {
		err = categorizeExecutionError(ctx, execCtx, err)
		result = markErrorResult(result, err)
		_ = i.runHooks(ctx, ToolEvent{Stage: HookStageToolFailure, Call: call, Descriptor: descriptor, Decision: decision, Result: result, Err: err})
		return i.finish(ctx, call, descriptor, decision, result, err, AuditStageFailed)
	}
	result = i.truncateResult(result)
	if err := i.runHooks(ctx, ToolEvent{Stage: HookStagePostToolUse, Call: call, Descriptor: descriptor, Decision: decision, Result: result}); err != nil {
		return i.finish(ctx, call, descriptor, decision, result, err, AuditStageFailed)
	}
	return i.finish(ctx, call, descriptor, decision, result, nil, AuditStageSucceeded)
}

func (i *PipelineInvoker) runHooks(ctx context.Context, event ToolEvent) error {
	for _, hook := range i.config.Hooks {
		if hook == nil {
			continue
		}
		decision, err := hook.HandleToolEvent(ctx, event)
		if err != nil {
			return WrapError(ToolExecutionError, "tool hook failed", err)
		}
		if decision.Effect == PermissionDeny {
			if decision.Reason == "" {
				decision.Reason = "tool hook denied invocation"
			}
			return NewError(PermissionDenied, decision.Reason)
		}
	}
	return nil
}

func (i *PipelineInvoker) finish(ctx context.Context, call ToolCall, descriptor ToolDescriptor, decision PermissionDecision, result ToolResult, err error, stage AuditStage) (ToolResult, error) {
	if err != nil {
		result = markErrorResult(result, err)
	}
	if auditErr := i.recordAudit(ctx, call, descriptor, decision, result, err, stage); auditErr != nil && err == nil {
		return markErrorResult(result, auditErr), auditErr
	}
	return result, err
}

func (i *PipelineInvoker) recordAudit(ctx context.Context, call ToolCall, descriptor ToolDescriptor, decision PermissionDecision, result ToolResult, err error, stage AuditStage) error {
	if i.config.Audit == nil {
		return nil
	}
	record := ToolAuditRecord{
		CallID:     call.ID,
		ToolName:   call.Name,
		Source:     descriptor.Source,
		Stage:      stage,
		Input:      append(json.RawMessage{}, call.Input...),
		Decision:   decision,
		Result:     result,
		OccurredAt: i.config.Now(),
	}
	if err != nil {
		record.Error = err.Error()
		if toolErr, ok := err.(*Error); ok {
			record.Category = toolErr.Category
		} else {
			record.Category = ToolExecutionError
		}
	}
	if record.ToolName == "" {
		record.ToolName = descriptor.Name
	}
	return i.config.Audit.RecordToolAudit(ctx, record)
}

func (i *PipelineInvoker) truncateResult(result ToolResult) ToolResult {
	if i.config.MaxOutputBytes <= 0 || len(result.Output) <= i.config.MaxOutputBytes {
		return result
	}
	original := len(result.Output)
	result.Output = string([]byte(result.Output)[:i.config.MaxOutputBytes])
	result.Truncated = true
	if result.Meta == nil {
		result.Meta = map[string]any{}
	}
	result.Meta["original_output_bytes"] = original
	result.Meta["max_output_bytes"] = i.config.MaxOutputBytes
	return result
}

func markErrorResult(result ToolResult, err error) ToolResult {
	result.IsError = true
	if toolErr, ok := err.(*Error); ok {
		result.Category = toolErr.Category
	} else {
		result.Category = ToolExecutionError
	}
	return result
}

func categorizeExecutionError(parent, execCtx context.Context, err error) error {
	if IsCategory(err, ValidationError) ||
		IsCategory(err, PermissionDenied) ||
		IsCategory(err, ApprovalRequired) ||
		IsCategory(err, TimeoutError) ||
		IsCategory(err, CancelledError) ||
		IsCategory(err, ToolExecutionError) ||
		IsCategory(err, PersistenceError) {
		return err
	}
	if execCtx.Err() == context.DeadlineExceeded || err == context.DeadlineExceeded {
		return WrapError(TimeoutError, "tool execution timed out", err)
	}
	if parent.Err() != nil || execCtx.Err() == context.Canceled || err == context.Canceled {
		return WrapError(CancelledError, "tool execution cancelled", err)
	}
	return WrapError(ToolExecutionError, "tool execution failed", err)
}

type invocationSchema struct {
	Type                 string                              `json:"type"`
	Required             []string                            `json:"required"`
	Properties           map[string]invocationSchemaProperty `json:"properties"`
	AdditionalProperties *bool                               `json:"additionalProperties"`
}

type invocationSchemaProperty struct {
	Type      string `json:"type"`
	MinLength *int   `json:"minLength"`
	Minimum   *int   `json:"minimum"`
}

func validateInputAgainstSchema(schema json.RawMessage, raw json.RawMessage) error {
	if len(raw) == 0 {
		return NewError(ValidationError, "tool input is required")
	}
	if !utf8.Valid(raw) || strings.Contains(string(raw), "\x00") {
		return NewError(ValidationError, "tool input must be valid UTF-8 without null bytes")
	}
	if !json.Valid(raw) {
		return NewError(ValidationError, "tool input must be valid JSON")
	}
	var input map[string]any
	if err := json.Unmarshal(raw, &input); err != nil {
		return WrapError(ValidationError, "decode tool input", err)
	}
	var spec invocationSchema
	if err := json.Unmarshal(schema, &spec); err != nil {
		return WrapError(ValidationError, "decode tool input schema", err)
	}
	if spec.Type != "" && spec.Type != "object" {
		return NewError(ValidationError, "only object input schemas are supported")
	}
	for _, field := range spec.Required {
		if _, ok := input[field]; !ok {
			return NewError(ValidationError, "missing required field: "+field)
		}
	}
	if spec.AdditionalProperties != nil && !*spec.AdditionalProperties {
		for key := range input {
			if _, ok := spec.Properties[key]; !ok {
				return NewError(ValidationError, "unexpected field: "+key)
			}
		}
	}
	for key, value := range input {
		prop, ok := spec.Properties[key]
		if !ok {
			continue
		}
		if err := validateInputProperty(key, value, prop); err != nil {
			return err
		}
	}
	return nil
}

func validateInputProperty(key string, value any, prop invocationSchemaProperty) error {
	switch prop.Type {
	case "string":
		str, ok := value.(string)
		if !ok {
			return NewError(ValidationError, key+" must be a string")
		}
		if prop.MinLength != nil && len(str) < *prop.MinLength {
			return NewError(ValidationError, key+" is too short")
		}
	case "integer":
		num, ok := value.(float64)
		if !ok || num != float64(int(num)) {
			return NewError(ValidationError, key+" must be an integer")
		}
		if prop.Minimum != nil && int(num) < *prop.Minimum {
			return NewError(ValidationError, fmt.Sprintf("%s must be >= %d", key, *prop.Minimum))
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return NewError(ValidationError, key+" must be a boolean")
		}
	}
	return nil
}
