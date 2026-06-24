package modelprovider

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"
)

var ErrMockScriptExhausted = errors.New("mock provider script exhausted")

type MockStep struct {
	Response ChatResponse
	Stream   []StreamChunk
	Err      error
	Delay    time.Duration
}

type MockProvider struct {
	mu           sync.Mutex
	name         string
	steps        []MockStep
	next         int
	requests     []ChatRequest
	capabilities map[string]ModelCapability
}

func NewMockProvider(steps ...MockStep) *MockProvider {
	return &MockProvider{
		name:         "mock",
		steps:        append([]MockStep(nil), steps...),
		capabilities: make(map[string]ModelCapability),
	}
}

func (p *MockProvider) Name() string {
	return p.name
}

func (p *MockProvider) Chat(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	step, err := p.nextStep(ctx, req)
	if err != nil {
		return ChatResponse{}, err
	}
	if step.Err != nil {
		return ChatResponse{}, step.Err
	}
	return step.Response, nil
}

func (p *MockProvider) ChatStream(ctx context.Context, req ChatRequest) (StreamReader, error) {
	step, err := p.nextStep(ctx, req)
	if err != nil {
		return nil, err
	}
	if step.Err != nil {
		return nil, step.Err
	}
	return &mockStreamReader{
		chunks:   append([]StreamChunk(nil), step.Stream...),
		response: step.Response,
	}, nil
}

func (p *MockProvider) Capability(model string) (ModelCapability, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	capability, ok := p.capabilities[model]
	if !ok {
		return ModelCapability{}, nil
	}
	return capability, nil
}

func (p *MockProvider) SetCapability(model string, capability ModelCapability) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.capabilities[model] = capability
}

func (p *MockProvider) Requests() []ChatRequest {
	p.mu.Lock()
	defer p.mu.Unlock()

	out := make([]ChatRequest, len(p.requests))
	copy(out, p.requests)
	return out
}

func (p *MockProvider) nextStep(ctx context.Context, req ChatRequest) (MockStep, error) {
	step, err := p.reserveStep(req)
	if err != nil {
		return MockStep{}, err
	}
	if step.Delay <= 0 {
		return step, nil
	}

	timer := time.NewTimer(step.Delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return step, nil
	case <-ctx.Done():
		return MockStep{}, providerErrorFromContext(ctx.Err())
	}
}

func (p *MockProvider) reserveStep(req ChatRequest) (MockStep, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.requests = append(p.requests, req)
	if p.next >= len(p.steps) {
		return MockStep{}, ErrMockScriptExhausted
	}
	step := p.steps[p.next]
	p.next++
	return step, nil
}

func providerErrorFromContext(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ProviderError{Category: ErrorCategoryTimeout, Cause: err}
	}
	return ProviderError{Category: ErrorCategoryCancelled, Cause: err}
}

type mockStreamReader struct {
	mu       sync.Mutex
	chunks   []StreamChunk
	index    int
	response ChatResponse
	closed   bool
}

func (r *mockStreamReader) Recv() (StreamChunk, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return StreamChunk{}, io.EOF
	}
	if r.index >= len(r.chunks) {
		return StreamChunk{}, io.EOF
	}
	chunk := r.chunks[r.index]
	r.index++
	return chunk, nil
}

func (r *mockStreamReader) Response() (ChatResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.response, nil
}

func (r *mockStreamReader) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.closed = true
	return nil
}
