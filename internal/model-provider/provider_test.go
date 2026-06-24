package modelprovider

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestChatRequestUsesNeutralSerializableTypes(t *testing.T) {
	temperature := 0.2
	maxTokens := 512
	req := ChatRequest{
		Model: "mock-large",
		Messages: []Message{
			{Role: MessageRoleSystem, Content: "You are a coding agent."},
			{Role: MessageRoleUser, Content: "Read README.md"},
		},
		Tools: []ToolSchema{{
			Name:        "ReadFile",
			Description: "Read a workspace file",
			Parameters:  json.RawMessage(`{"type":"object","required":["path"]}`),
		}},
		ToolChoice:  ToolChoiceAuto,
		Temperature: &temperature,
		MaxTokens:   &maxTokens,
		Stop:        []string{"</done>"},
	}

	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal ChatRequest: %v", err)
	}

	var decoded ChatRequest
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal ChatRequest: %v", err)
	}

	if got, want := decoded.Messages[0].Role, MessageRoleSystem; got != want {
		t.Fatalf("decoded role = %s, want %s", got, want)
	}
	if got, want := decoded.Tools[0].Name, "ReadFile"; got != want {
		t.Fatalf("decoded tool name = %q, want %q", got, want)
	}
	if got, want := decoded.ToolChoice, ToolChoiceAuto; got != want {
		t.Fatalf("decoded tool choice = %s, want %s", got, want)
	}
}

func TestChatResponseRepresentsToolUseWithoutProviderSpecificFields(t *testing.T) {
	resp := ChatResponse{
		Message: Message{
			Role:    MessageRoleAssistant,
			Content: "I need to inspect a file.",
		},
		ToolCalls: []ToolCall{{
			ID:        "call-1",
			Name:      "ReadFile",
			Arguments: json.RawMessage(`{"path":"README.md"}`),
		}},
		StopReason: StopReasonToolUse,
		Usage: TokenUsage{
			InputTokens:  10,
			OutputTokens: 5,
			TotalTokens:  15,
		},
		Model: "mock-large",
	}

	if got, want := resp.StopReason, StopReasonToolUse; got != want {
		t.Fatalf("StopReason = %s, want %s", got, want)
	}
	if got, want := json.Valid(resp.ToolCalls[0].Arguments), true; got != want {
		t.Fatalf("tool call arguments valid = %v, want %v", got, want)
	}
	if got, want := resp.Usage.TotalTokens, 15; got != want {
		t.Fatalf("total tokens = %d, want %d", got, want)
	}
}

func TestProviderErrorCarriesRetryMetadataAndCause(t *testing.T) {
	cause := errors.New("upstream returned 429")
	err := ProviderError{
		Category:   ErrorCategoryProvider,
		Kind:       ProviderErrRateLimit,
		Retryable:  true,
		RetryAfter: 2 * time.Second,
		Cause:      cause,
	}

	if !err.Retryable {
		t.Fatal("rate limit error should be retryable")
	}
	if got, want := err.RetryAfter, 2*time.Second; got != want {
		t.Fatalf("RetryAfter = %s, want %s", got, want)
	}
	if !errors.Is(err, cause) {
		t.Fatal("ProviderError should unwrap the original cause")
	}
	if got, want := err.Error(), "ProviderError: RateLimit: upstream returned 429"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestProviderInterfaceAcceptsNeutralTypes(t *testing.T) {
	var provider Provider = staticProvider{name: "mock"}

	resp, err := provider.Chat(context.Background(), ChatRequest{Model: "mock-large"})
	if err != nil {
		t.Fatalf("Chat returned error: %v", err)
	}
	if got, want := resp.StopReason, StopReasonEndTurn; got != want {
		t.Fatalf("StopReason = %s, want %s", got, want)
	}
	if got, want := provider.Name(), "mock"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

type staticProvider struct {
	name string
}

func (p staticProvider) Chat(context.Context, ChatRequest) (ChatResponse, error) {
	return ChatResponse{StopReason: StopReasonEndTurn}, nil
}

func (p staticProvider) ChatStream(context.Context, ChatRequest) (StreamReader, error) {
	return nil, nil
}

func (p staticProvider) Capability(string) (ModelCapability, error) {
	return ModelCapability{SupportsStreaming: true}, nil
}

func (p staticProvider) Name() string {
	return p.name
}
