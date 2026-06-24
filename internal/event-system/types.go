package eventsystem

import (
	"context"
	"encoding/json"
	"time"
)

type Event struct {
	EventID       string
	EventType     EventType
	Timestamp     time.Time
	SchemaVersion int
	SessionID     string
	AgentID       string
	TaskID        string
	TeamID        string
	CorrelationID string
	CausationID   string
	Sequence      int64
	Payload       json.RawMessage
}

type EventType string

const (
	EventSessionStart      EventType = "SessionStart"
	EventSessionEnd        EventType = "SessionEnd"
	EventUserPromptSubmit  EventType = "UserPromptSubmit"
	EventPreModelCall      EventType = "PreModelCall"
	EventPostModelCall     EventType = "PostModelCall"
	EventModelCallFailed   EventType = "ModelCallFailed"
	EventPreToolUse        EventType = "PreToolUse"
	EventPostToolUse       EventType = "PostToolUse"
	EventToolFailure       EventType = "ToolFailure"
	EventApprovalRequested EventType = "ApprovalRequested"
	EventApprovalResolved  EventType = "ApprovalResolved"
	EventPreCompact        EventType = "PreCompact"
	EventPostCompact       EventType = "PostCompact"
	EventMemoryRead        EventType = "MemoryRead"
	EventMemoryWrite       EventType = "MemoryWrite"
	EventSubAgentStart     EventType = "SubAgentStart"
	EventSubAgentStop      EventType = "SubAgentStop"
	EventWorktreeCreate    EventType = "WorktreeCreate"
	EventWorktreeRemove    EventType = "WorktreeRemove"
	EventTeamCreated       EventType = "TeamCreated"
	EventTeamClosed        EventType = "TeamClosed"
	EventTaskCreated       EventType = "TaskCreated"
	EventTaskAssigned      EventType = "TaskAssigned"
	EventTaskCompleted     EventType = "TaskCompleted"
	EventTaskFailed        EventType = "TaskFailed"

	EventAgentStateChanged  EventType = "AgentStateChanged"
	EventToolRequested      EventType = "ToolRequested"
	EventToolObserved       EventType = "ToolObserved"
	EventCheckpointCreated  EventType = "CheckpointCreated"
	EventCheckpointRestored EventType = "CheckpointRestored"
	EventBudgetExceeded     EventType = "BudgetExceeded"
	EventLoopDetected       EventType = "LoopDetected"
	EventProviderRetry      EventType = "ProviderRetry"
	EventSessionPaused      EventType = "SessionPaused"
	EventSessionResumed     EventType = "SessionResumed"
	EventAuditRecorded      EventType = "AuditRecorded"
)

var allEventTypes = []EventType{
	EventSessionStart,
	EventSessionEnd,
	EventUserPromptSubmit,
	EventPreModelCall,
	EventPostModelCall,
	EventModelCallFailed,
	EventPreToolUse,
	EventPostToolUse,
	EventToolFailure,
	EventApprovalRequested,
	EventApprovalResolved,
	EventPreCompact,
	EventPostCompact,
	EventMemoryRead,
	EventMemoryWrite,
	EventSubAgentStart,
	EventSubAgentStop,
	EventWorktreeCreate,
	EventWorktreeRemove,
	EventTeamCreated,
	EventTeamClosed,
	EventTaskCreated,
	EventTaskAssigned,
	EventTaskCompleted,
	EventTaskFailed,
	EventAgentStateChanged,
	EventToolRequested,
	EventToolObserved,
	EventCheckpointCreated,
	EventCheckpointRestored,
	EventBudgetExceeded,
	EventLoopDetected,
	EventProviderRetry,
	EventSessionPaused,
	EventSessionResumed,
	EventAuditRecorded,
}

func AllEventTypes() []EventType {
	out := make([]EventType, len(allEventTypes))
	copy(out, allEventTypes)
	return out
}

type EventClass uint8

const (
	EventClassDurable EventClass = 1 << iota
	EventClassRecovery
	EventClassAudit
	EventClassHook
	EventClassEphemeral
)

func (c EventClass) Contains(required EventClass) bool {
	return c&required == required
}

func EventClasses(eventType EventType) EventClass {
	switch eventType {
	case EventSessionStart, EventSessionEnd, EventUserPromptSubmit, EventPreModelCall, EventPostModelCall,
		EventModelCallFailed, EventPreToolUse, EventPostToolUse, EventToolFailure, EventApprovalRequested,
		EventPreCompact, EventPostCompact, EventMemoryRead, EventMemoryWrite, EventSubAgentStart,
		EventSubAgentStop, EventWorktreeCreate, EventWorktreeRemove, EventTeamCreated, EventTeamClosed,
		EventTaskCreated, EventTaskAssigned, EventTaskCompleted, EventTaskFailed:
		return lifecycleEventClasses(eventType)
	case EventApprovalResolved:
		return EventClassDurable | EventClassRecovery | EventClassAudit | EventClassHook
	case EventAgentStateChanged, EventToolRequested, EventToolObserved, EventCheckpointCreated,
		EventCheckpointRestored, EventSessionPaused, EventSessionResumed:
		return EventClassDurable | EventClassRecovery
	case EventBudgetExceeded, EventLoopDetected, EventProviderRetry:
		return EventClassDurable
	case EventAuditRecorded:
		return EventClassDurable | EventClassAudit
	default:
		return 0
	}
}

func lifecycleEventClasses(eventType EventType) EventClass {
	classes := EventClassHook
	switch eventType {
	case EventSessionStart, EventSessionEnd, EventUserPromptSubmit, EventPostModelCall, EventPostToolUse:
		classes |= EventClassDurable
	case EventApprovalRequested, EventPreToolUse, EventWorktreeCreate, EventWorktreeRemove:
		classes |= EventClassAudit
	}
	return classes
}

type Bus interface {
	Publish(ctx context.Context, e Event) error
	Subscribe(filter SubscriptionFilter, s Subscriber) (Unsubscribe, error)
}

type Subscriber interface {
	OnEvent(ctx context.Context, e Event) error
}

type Unsubscribe func() error

type SubscriptionFilter struct {
	Types []EventType
	Class EventClass
}
