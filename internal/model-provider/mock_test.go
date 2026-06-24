package modelprovider

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"testing"
	"time"
)

func TestMockProviderReplaysScriptedResponsesAndRecordsRequests(t *testing.T) {
	provider := NewMockProvider(
		MockStep{Response: ChatResponse{
			Message:    Message{Role: MessageRoleAssistant, Content: "first"},
			StopReason: StopReasonEndTurn,
			Model:      "mock-large",
		}},
		MockStep{Response: ChatResponse{
			Message:    Message{Role: MessageRoleAssistant, Content: "second"},
			StopReason: StopReasonEndTurn,
			Model:      "mock-large",
		}},
	)

	first, err := provider.Chat(context.Background(), ChatRequest{
		Model:    "mock-large",
		Messages: []Message{{Role: MessageRoleUser, Content: "one"}},
	})
	if err != nil {
		t.Fatalf("first Chat returned error: %v", err)
	}
	second, err := provider.Chat(context.Background(), ChatRequest{
		Model:    "mock-large",
		Messages: []Message{{Role: MessageRoleUser, Content: "two"}},
	})
	if err != nil {
		t.Fatalf("second Chat returned error: %v", err)
	}

	if got, want := first.Message.Content, "first"; got != want {
		t.Fatalf("first content = %q, want %q", got, want)
	}
	if got, want := second.Message.Content, "second"; got != want {
		t.Fatalf("second content = %q, want %q", got, want)
	}

	requests := provider.Requests()
	if got, want := len(requests), 2; got != want {
		t.Fatalf("recorded request count = %d, want %d", got, want)
	}
	if got, want := requests[0].Messages[0].Content, "one"; got != want {
		t.Fatalf("first recorded prompt = %q, want %q", got, want)
	}
	if got, want := requests[1].Messages[0].Content, "two"; got != want {
		t.Fatalf("second recorded prompt = %q, want %q", got, want)
	}
}

func TestMockProviderSupportsToolCallsAndScriptedErrors(t *testing.T) {
	scriptedErr := ProviderError{
		Category:  ErrorCategoryProvider,
		Kind:      ProviderErrServerError,
		Retryable: true,
		Cause:     errors.New("upstream unavailable"),
	}
	provider := NewMockProvider(
		MockStep{Response: ChatResponse{
			Message: Message{Role: MessageRoleAssistant, Content: "need tools"},
			ToolCalls: []ToolCall{
				{ID: "call-1", Name: "ReadFile", Arguments: json.RawMessage(`{"path":"README.md"}`)},
				{ID: "call-2", Name: "Grep", Arguments: json.RawMessage(`{"pattern":"ForgeCode"}`)},
			},
			StopReason: StopReasonToolUse,
		}},
		MockStep{Err: scriptedErr},
	)

	resp, err := provider.Chat(context.Background(), ChatRequest{Model: "mock-large"})
	if err != nil {
		t.Fatalf("tool call response returned error: %v", err)
	}
	if got, want := len(resp.ToolCalls), 2; got != want {
		t.Fatalf("tool call count = %d, want %d", got, want)
	}
	if got, want := resp.ToolCalls[1].ID, "call-2"; got != want {
		t.Fatalf("second tool call id = %q, want %q", got, want)
	}

	_, err = provider.Chat(context.Background(), ChatRequest{Model: "mock-large"})
	if !errors.Is(err, scriptedErr.Cause) {
		t.Fatalf("second Chat error = %v, want cause %v", err, scriptedErr.Cause)
	}
	var providerErr ProviderError
	if !errors.As(err, &providerErr) {
		t.Fatalf("second Chat error type = %T, want ProviderError", err)
	}
	if !providerErr.Retryable {
		t.Fatal("scripted ProviderError should preserve Retryable=true")
	}
}

func TestMockProviderDelayHonorsContextDeadline(t *testing.T) {
	provider := NewMockProvider(MockStep{
		Delay: 50 * time.Millisecond,
		Response: ChatResponse{
			Message:    Message{Role: MessageRoleAssistant, Content: "late"},
			StopReason: StopReasonEndTurn,
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	_, err := provider.Chat(ctx, ChatRequest{Model: "mock-large"})
	if err == nil {
		t.Fatal("expected delayed Chat to return context error")
	}
	var providerErr ProviderError
	if !errors.As(err, &providerErr) {
		t.Fatalf("error type = %T, want ProviderError", err)
	}
	if got, want := providerErr.Category, ErrorCategoryTimeout; got != want {
		t.Fatalf("error category = %s, want %s", got, want)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want context deadline cause", err)
	}
}

func TestMockProviderCapabilityCanBeConfigured(t *testing.T) {
	provider := NewMockProvider()
	provider.SetCapability("mock-large", ModelCapability{
		SupportsTools:         true,
		SupportsStreaming:     true,
		SupportsStructuredOut: true,
		ContextWindow: ContextWindow{
			MaxInputTokens:  128000,
			MaxOutputTokens: 4096,
		},
	})

	capability, err := provider.Capability("mock-large")
	if err != nil {
		t.Fatalf("Capability returned error: %v", err)
	}
	if !capability.SupportsTools || !capability.SupportsStreaming || !capability.SupportsStructuredOut {
		t.Fatalf("capability flags not preserved: %#v", capability)
	}
	if got, want := capability.ContextWindow.MaxInputTokens, 128000; got != want {
		t.Fatalf("MaxInputTokens = %d, want %d", got, want)
	}
}

func TestMockProviderStreamsScriptedChunks(t *testing.T) {
	final := ChatResponse{
		Message:    Message{Role: MessageRoleAssistant, Content: "hello"},
		StopReason: StopReasonEndTurn,
		Usage:      TokenUsage{InputTokens: 3, OutputTokens: 2, TotalTokens: 5},
	}
	provider := NewMockProvider(MockStep{
		Stream: []StreamChunk{
			{Kind: ChunkKindTextDelta, TextDelta: "hel"},
			{Kind: ChunkKindTextDelta, TextDelta: "lo"},
			{Kind: ChunkKindStopReason, StopReason: StopReasonEndTurn},
		},
		Response: final,
	})

	reader, err := provider.ChatStream(context.Background(), ChatRequest{Model: "mock-large"})
	if err != nil {
		t.Fatalf("ChatStream returned error: %v", err)
	}
	defer reader.Close()

	var chunks []StreamChunk
	for {
		chunk, err := reader.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("Recv returned error: %v", err)
		}
		chunks = append(chunks, chunk)
	}

	wantChunks := []StreamChunk{
		{Kind: ChunkKindTextDelta, TextDelta: "hel"},
		{Kind: ChunkKindTextDelta, TextDelta: "lo"},
		{Kind: ChunkKindStopReason, StopReason: StopReasonEndTurn},
	}
	if !reflect.DeepEqual(chunks, wantChunks) {
		t.Fatalf("chunks = %#v, want %#v", chunks, wantChunks)
	}

	resp, err := reader.Response()
	if err != nil {
		t.Fatalf("Response returned error: %v", err)
	}
	if !reflect.DeepEqual(resp, final) {
		t.Fatalf("response = %#v, want %#v", resp, final)
	}
}

func TestMockProviderReportsScriptExhaustion(t *testing.T) {
	provider := NewMockProvider()
	_, err := provider.Chat(context.Background(), ChatRequest{Model: "mock-large"})
	if !errors.Is(err, ErrMockScriptExhausted) {
		t.Fatalf("error = %v, want ErrMockScriptExhausted", err)
	}
}
