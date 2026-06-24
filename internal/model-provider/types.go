package modelprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Provider interface {
	Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)
	ChatStream(ctx context.Context, req ChatRequest) (StreamReader, error)
	Capability(model string) (ModelCapability, error)
	Name() string
}

type StreamReader interface {
	Recv() (StreamChunk, error)
	Response() (ChatResponse, error)
	Close() error
}

type ChatRequest struct {
	Model         string
	Messages      []Message
	Tools         []ToolSchema
	ToolChoice    ToolChoice
	Temperature   *float64
	MaxTokens     *int
	StructuredOut *StructuredSpec
	Stop          []string
}

type ChatResponse struct {
	Message    Message
	ToolCalls  []ToolCall
	StopReason StopReason
	Usage      TokenUsage
	Model      string
}

type Message struct {
	Role       MessageRole
	Content    string
	ToolCallID string
	ToolCalls  []ToolCall
}

type MessageRole string

const (
	MessageRoleSystem    MessageRole = "system"
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleTool      MessageRole = "tool"
)

type ToolSchema struct {
	Name        string
	Description string
	Parameters  json.RawMessage
}

type ToolChoice string

const (
	ToolChoiceAuto     ToolChoice = "Auto"
	ToolChoiceNone     ToolChoice = "None"
	ToolChoiceRequired ToolChoice = "Required"
	ToolChoiceNamed    ToolChoice = "Named"
)

type StructuredSpec struct {
	Name   string
	Schema json.RawMessage
	Strict bool
}

type StreamChunk struct {
	Kind          ChunkKind
	TextDelta     string
	ToolCallDelta *ToolCallDelta
	StopReason    StopReason
	Usage         *TokenUsage
}

type ChunkKind string

const (
	ChunkKindTextDelta     ChunkKind = "TextDelta"
	ChunkKindToolCallDelta ChunkKind = "ToolCallDelta"
	ChunkKindStopReason    ChunkKind = "StopReason"
	ChunkKindUsageUpdate   ChunkKind = "UsageUpdate"
)

type ToolCall struct {
	ID        string
	Name      string
	Arguments json.RawMessage
}

type ToolCallDelta struct {
	ID             string
	NameDelta      string
	ArgumentsDelta string
	Index          int
}

type TokenUsage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
}

type StopReason string

const (
	StopReasonEndTurn       StopReason = "EndTurn"
	StopReasonToolUse       StopReason = "ToolUse"
	StopReasonMaxTokens     StopReason = "MaxTokens"
	StopReasonStopSequence  StopReason = "StopSequence"
	StopReasonContentFilter StopReason = "ContentFilter"
)

type ModelCapability struct {
	SupportsTools         bool
	SupportsStreaming     bool
	SupportsStructuredOut bool
	SupportsVision        bool
	ContextWindow         ContextWindow
}

type ContextWindow struct {
	MaxInputTokens  int
	MaxOutputTokens int
}

type ErrorCategory string

const (
	ErrorCategoryValidation ErrorCategory = "ValidationError"
	ErrorCategoryProvider   ErrorCategory = "ProviderError"
	ErrorCategoryTimeout    ErrorCategory = "TimeoutError"
	ErrorCategoryCancelled  ErrorCategory = "CancelledError"
)

type ProviderErrKind string

const (
	ProviderErrRateLimit   ProviderErrKind = "RateLimit"
	ProviderErrAuth        ProviderErrKind = "Auth"
	ProviderErrBadRequest  ProviderErrKind = "BadRequest"
	ProviderErrServerError ProviderErrKind = "ServerError"
	ProviderErrNetwork     ProviderErrKind = "Network"
	ProviderErrOverloaded  ProviderErrKind = "Overloaded"
)

type ProviderError struct {
	Category   ErrorCategory
	Kind       ProviderErrKind
	Retryable  bool
	RetryAfter time.Duration
	Cause      error
}

func (e ProviderError) Error() string {
	if e.Cause == nil {
		return fmt.Sprintf("%s: %s", e.Category, e.Kind)
	}
	return fmt.Sprintf("%s: %s: %s", e.Category, e.Kind, e.Cause)
}

func (e ProviderError) Unwrap() error {
	return e.Cause
}
