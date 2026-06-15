package toolruntime

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ToolSource string

const (
	ToolSourceBuiltin ToolSource = "builtin"
	ToolSourceMCP     ToolSource = "mcp"
	ToolSourceSkill   ToolSource = "skill"
)

type RiskLevel string

const (
	RiskLow      RiskLevel = "Low"
	RiskMedium   RiskLevel = "Medium"
	RiskHigh     RiskLevel = "High"
	RiskCritical RiskLevel = "Critical"
)

type PermissionAction string

const (
	PermissionRead    PermissionAction = "read"
	PermissionWrite   PermissionAction = "write"
	PermissionExecute PermissionAction = "execute"
	PermissionSearch  PermissionAction = "search"
)

type PermissionHint struct {
	Actions []PermissionAction `json:"actions"`
}

type ToolDescriptor struct {
	Name        string          `json:"name"`
	Source      ToolSource      `json:"source"`
	InputSchema json.RawMessage `json:"input_schema"`
	Risk        RiskLevel       `json:"risk"`
	Permission  PermissionHint  `json:"permission"`
}

func (d ToolDescriptor) Validate() error {
	if d.Name == "" {
		return NewError(ValidationError, "tool descriptor name is required")
	}
	if d.Source == "" {
		return NewError(ValidationError, fmt.Sprintf("tool %q source is required", d.Name))
	}
	if d.Risk == "" {
		return NewError(ValidationError, fmt.Sprintf("tool %q risk is required", d.Name))
	}
	if len(d.InputSchema) == 0 || !json.Valid(d.InputSchema) {
		return NewError(ValidationError, fmt.Sprintf("tool %q input schema must be valid JSON", d.Name))
	}
	return nil
}

type InvocationContext struct {
	SessionID     string
	AgentID       string
	TeamID        string
	WorkspaceRoot string
	Source        string
}

type ToolCall struct {
	ID    string
	Name  string
	Input json.RawMessage
	Ctx   InvocationContext
}

type ToolResult struct {
	CallID    string
	Output    string
	Truncated bool
	IsError   bool
	Category  ErrorCategory
	Meta      map[string]any
}

type ErrorCategory string

const (
	ValidationError    ErrorCategory = "ValidationError"
	PermissionDenied   ErrorCategory = "PermissionDenied"
	ApprovalRequired   ErrorCategory = "ApprovalRequired"
	TimeoutError       ErrorCategory = "TimeoutError"
	CancelledError     ErrorCategory = "CancelledError"
	ProviderError      ErrorCategory = "ProviderError"
	ToolExecutionError ErrorCategory = "ToolExecutionError"
	SandboxError       ErrorCategory = "SandboxError"
	PersistenceError   ErrorCategory = "PersistenceError"
	ConflictError      ErrorCategory = "ConflictError"
	RecoveryError      ErrorCategory = "RecoveryError"
)

type Error struct {
	Category ErrorCategory
	Message  string
	Err      error
}

func NewError(category ErrorCategory, message string) *Error {
	return &Error{Category: category, Message: message}
}

func WrapError(category ErrorCategory, message string, err error) *Error {
	return &Error{Category: category, Message: message, Err: err}
}

func (e *Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s: %s", e.Category, e.Message)
	}
	return fmt.Sprintf("%s: %s: %v", e.Category, e.Message, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func IsCategory(err error, category ErrorCategory) bool {
	var toolErr *Error
	return errors.As(err, &toolErr) && toolErr.Category == category
}
