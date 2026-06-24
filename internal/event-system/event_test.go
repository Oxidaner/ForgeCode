package eventsystem

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestEventEnvelopeFieldsMatchArchitectureDocument(t *testing.T) {
	eventType := reflect.TypeOf(Event{})
	fields := []struct {
		name string
		typ  reflect.Type
	}{
		{"EventID", reflect.TypeOf("")},
		{"EventType", reflect.TypeOf(EventType(""))},
		{"Timestamp", reflect.TypeOf(time.Time{})},
		{"SchemaVersion", reflect.TypeOf(0)},
		{"SessionID", reflect.TypeOf("")},
		{"AgentID", reflect.TypeOf("")},
		{"TaskID", reflect.TypeOf("")},
		{"TeamID", reflect.TypeOf("")},
		{"CorrelationID", reflect.TypeOf("")},
		{"CausationID", reflect.TypeOf("")},
		{"Sequence", reflect.TypeOf(int64(0))},
		{"Payload", reflect.TypeOf(json.RawMessage{})},
	}

	if got, want := eventType.NumField(), len(fields); got != want {
		t.Fatalf("Event field count = %d, want %d", got, want)
	}
	for _, field := range fields {
		got, ok := eventType.FieldByName(field.name)
		if !ok {
			t.Fatalf("Event missing field %s", field.name)
		}
		if got.Type != field.typ {
			t.Fatalf("Event.%s type = %s, want %s", field.name, got.Type, field.typ)
		}
	}
}

func TestAllDocumentedEventTypesAreDefinedInOrder(t *testing.T) {
	want := []EventType{
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

	got := AllEventTypes()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("AllEventTypes() = %#v, want %#v", got, want)
	}
	for _, eventType := range got {
		if string(eventType) == "" {
			t.Fatal("EventType string must not be empty")
		}
	}
}

func TestEventClassMappingSupportsMultipleClasses(t *testing.T) {
	tests := []struct {
		eventType EventType
		want      EventClass
	}{
		{EventSessionStart, EventClassDurable | EventClassHook},
		{EventAgentStateChanged, EventClassDurable | EventClassRecovery},
		{EventApprovalResolved, EventClassDurable | EventClassRecovery | EventClassAudit | EventClassHook},
		{EventToolObserved, EventClassDurable | EventClassRecovery},
		{EventAuditRecorded, EventClassDurable | EventClassAudit},
	}

	for _, tt := range tests {
		got := EventClasses(tt.eventType)
		if got != tt.want {
			t.Fatalf("EventClasses(%s) = %v, want %v", tt.eventType, got, tt.want)
		}
	}
}

func TestEventClassContains(t *testing.T) {
	classes := EventClassDurable | EventClassRecovery | EventClassAudit
	if !classes.Contains(EventClassDurable) {
		t.Fatal("expected combined class to contain Durable")
	}
	if !classes.Contains(EventClassRecovery | EventClassAudit) {
		t.Fatal("expected combined class to contain Recovery and Audit")
	}
	if classes.Contains(EventClassHook) {
		t.Fatal("did not expect combined class to contain Hook")
	}
}
