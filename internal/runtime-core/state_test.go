package runtimecore

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	eventsystem "github.com/Oxidaner/ForgeCode/internal/event-system"
)

func TestAllAgentStatesMatchGlossaryOrder(t *testing.T) {
	want := []AgentStateName{
		AgentStateCreated,
		AgentStateInitializing,
		AgentStateThinking,
		AgentStateToolRequested,
		AgentStateAwaitingApproval,
		AgentStateToolExecuting,
		AgentStateObserving,
		AgentStateCompacting,
		AgentStatePaused,
		AgentStateCompleted,
		AgentStateFailed,
		AgentStateCancelled,
	}

	if got := AllAgentStates(); !reflect.DeepEqual(got, want) {
		t.Fatalf("AllAgentStates() = %#v, want %#v", got, want)
	}
}

func TestStateMachineAllowsOnlySpecTransitions(t *testing.T) {
	sm := NewStateMachine()
	legal := []struct {
		from AgentStateName
		to   AgentStateName
	}{
		{AgentStateCreated, AgentStateInitializing},
		{AgentStateInitializing, AgentStateThinking},
		{AgentStateThinking, AgentStateToolRequested},
		{AgentStateThinking, AgentStateCompleted},
		{AgentStateThinking, AgentStateCompacting},
		{AgentStateToolRequested, AgentStateAwaitingApproval},
		{AgentStateToolRequested, AgentStateToolExecuting},
		{AgentStateAwaitingApproval, AgentStateToolExecuting},
		{AgentStateAwaitingApproval, AgentStateObserving},
		{AgentStateToolExecuting, AgentStateObserving},
		{AgentStateObserving, AgentStateThinking},
		{AgentStateCompacting, AgentStateThinking},
		{AgentStateThinking, AgentStatePaused},
		{AgentStatePaused, AgentStateThinking},
		{AgentStateThinking, AgentStateCancelled},
		{AgentStateToolExecuting, AgentStateCancelled},
		{AgentStateThinking, AgentStateFailed},
	}

	for _, edge := range legal {
		if !sm.CanTransition(edge.from, edge.to) {
			t.Fatalf("expected transition %s -> %s to be legal", edge.from, edge.to)
		}
	}

	illegal := []struct {
		from AgentStateName
		to   AgentStateName
	}{
		{AgentStateCreated, AgentStateThinking},
		{AgentStateCompleted, AgentStateThinking},
		{AgentStateFailed, AgentStateThinking},
		{AgentStateCancelled, AgentStateThinking},
		{AgentStateObserving, AgentStateCompleted},
	}
	for _, edge := range illegal {
		if sm.CanTransition(edge.from, edge.to) {
			t.Fatalf("expected transition %s -> %s to be illegal", edge.from, edge.to)
		}
	}
}

func TestTransitionProducesAgentStateChangedEvent(t *testing.T) {
	sm := NewStateMachine()
	result, err := sm.Transition(AgentStateCreated, AgentStateInitializing, TransitionContext{
		SessionID:     "sess-1",
		AgentID:       "agent-1",
		CorrelationID: "corr-1",
		Reason:        "session start",
	})
	if err != nil {
		t.Fatalf("Transition returned error: %v", err)
	}

	if got, want := result.State, AgentStateInitializing; got != want {
		t.Fatalf("state = %s, want %s", got, want)
	}
	if got, want := result.Event.EventType, eventsystem.EventAgentStateChanged; got != want {
		t.Fatalf("event type = %s, want %s", got, want)
	}
	if got, want := result.Event.SessionID, "sess-1"; got != want {
		t.Fatalf("event session id = %q, want %q", got, want)
	}

	var payload AgentStateChangedPayload
	if err := json.Unmarshal(result.Event.Payload, &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if got, want := payload.From, AgentStateCreated; got != want {
		t.Fatalf("payload from = %s, want %s", got, want)
	}
	if got, want := payload.To, AgentStateInitializing; got != want {
		t.Fatalf("payload to = %s, want %s", got, want)
	}
	if got, want := payload.Reason, "session start"; got != want {
		t.Fatalf("payload reason = %q, want %q", got, want)
	}
}

func TestInvalidTransitionReturnsTypedError(t *testing.T) {
	sm := NewStateMachine()
	_, err := sm.Transition(AgentStateCompleted, AgentStateThinking, TransitionContext{})
	if err == nil {
		t.Fatal("expected invalid transition to return error")
	}
	var invalid InvalidTransitionError
	if !errors.As(err, &invalid) {
		t.Fatalf("error = %T, want InvalidTransitionError", err)
	}
	if got, want := invalid.From, AgentStateCompleted; got != want {
		t.Fatalf("invalid from = %s, want %s", got, want)
	}
	if got, want := invalid.To, AgentStateThinking; got != want {
		t.Fatalf("invalid to = %s, want %s", got, want)
	}
}
