package runtimecore

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	eventsystem "github.com/Oxidaner/ForgeCode/internal/event-system"
)

type AgentStateName string

const (
	AgentStateCreated          AgentStateName = "Created"
	AgentStateInitializing     AgentStateName = "Initializing"
	AgentStateThinking         AgentStateName = "Thinking"
	AgentStateToolRequested    AgentStateName = "ToolRequested"
	AgentStateAwaitingApproval AgentStateName = "AwaitingApproval"
	AgentStateToolExecuting    AgentStateName = "ToolExecuting"
	AgentStateObserving        AgentStateName = "Observing"
	AgentStateCompacting       AgentStateName = "Compacting"
	AgentStatePaused           AgentStateName = "Paused"
	AgentStateCompleted        AgentStateName = "Completed"
	AgentStateFailed           AgentStateName = "Failed"
	AgentStateCancelled        AgentStateName = "Cancelled"
)

var allAgentStates = []AgentStateName{
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

func AllAgentStates() []AgentStateName {
	out := make([]AgentStateName, len(allAgentStates))
	copy(out, allAgentStates)
	return out
}

type StateMachine struct {
	transitions map[AgentStateName]map[AgentStateName]struct{}
}

func NewStateMachine() StateMachine {
	return StateMachine{transitions: buildTransitions([]stateTransition{
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
	})}
}

func (m StateMachine) CanTransition(from, to AgentStateName) bool {
	targets, ok := m.transitions[from]
	if !ok {
		return false
	}
	_, ok = targets[to]
	return ok
}

func (m StateMachine) Transition(from, to AgentStateName, ctx TransitionContext) (TransitionResult, error) {
	if !m.CanTransition(from, to) {
		return TransitionResult{}, InvalidTransitionError{From: from, To: to}
	}

	payload := AgentStateChangedPayload{
		From:   from,
		To:     to,
		Reason: ctx.Reason,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return TransitionResult{}, fmt.Errorf("marshal state transition payload: %w", err)
	}

	return TransitionResult{
		State: to,
		Event: eventsystem.Event{
			EventID:       newRuntimeEventID(),
			EventType:     eventsystem.EventAgentStateChanged,
			Timestamp:     time.Now().UTC(),
			SchemaVersion: 1,
			SessionID:     ctx.SessionID,
			AgentID:       ctx.AgentID,
			TaskID:        ctx.TaskID,
			TeamID:        ctx.TeamID,
			CorrelationID: ctx.CorrelationID,
			CausationID:   ctx.CausationID,
			Payload:       raw,
		},
	}, nil
}

type TransitionContext struct {
	SessionID     string
	AgentID       string
	TaskID        string
	TeamID        string
	CorrelationID string
	CausationID   string
	Reason        string
}

type TransitionResult struct {
	State AgentStateName
	Event eventsystem.Event
}

type AgentStateChangedPayload struct {
	From   AgentStateName
	To     AgentStateName
	Reason string
}

type InvalidTransitionError struct {
	From AgentStateName
	To   AgentStateName
}

func (e InvalidTransitionError) Error() string {
	return fmt.Sprintf("invalid agent state transition: %s -> %s", e.From, e.To)
}

type stateTransition struct {
	from AgentStateName
	to   AgentStateName
}

func buildTransitions(edges []stateTransition) map[AgentStateName]map[AgentStateName]struct{} {
	transitions := make(map[AgentStateName]map[AgentStateName]struct{}, len(allAgentStates))
	for _, edge := range edges {
		if transitions[edge.from] == nil {
			transitions[edge.from] = make(map[AgentStateName]struct{})
		}
		transitions[edge.from][edge.to] = struct{}{}
	}
	return transitions
}

func newRuntimeEventID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err == nil {
		return hex.EncodeToString(bytes[:])
	}
	return fmt.Sprintf("runtime-%d", time.Now().UnixNano())
}
