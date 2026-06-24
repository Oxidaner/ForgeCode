package sessionstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"sync"
	"testing"
	"time"

	eventsystem "github.com/Oxidaner/ForgeCode/internal/event-system"
	_ "modernc.org/sqlite"
)

func TestSQLiteEventLogAssignsContiguousSequencePerSession(t *testing.T) {
	ctx := context.Background()
	log := newTestEventLog(t, ctx)

	first := testEvent("event-1", "session-a", eventsystem.EventSessionStart)
	second := testEvent("event-2", "session-a", eventsystem.EventUserPromptSubmit)
	third := testEvent("event-3", "session-a", eventsystem.EventAgentStateChanged)

	seq1, err := log.Append(ctx, first)
	if err != nil {
		t.Fatalf("append first: %v", err)
	}
	seq2, err := log.Append(ctx, second)
	if err != nil {
		t.Fatalf("append second: %v", err)
	}
	seq3, err := log.Append(ctx, third)
	if err != nil {
		t.Fatalf("append third: %v", err)
	}

	if got, want := []int64{seq1, seq2, seq3}, []int64{1, 2, 3}; !equalInt64s(got, want) {
		t.Fatalf("sequences = %v, want %v", got, want)
	}

	events, err := log.Read(ctx, "session-a", 0)
	if err != nil {
		t.Fatalf("read events: %v", err)
	}
	if got, want := len(events), 3; got != want {
		t.Fatalf("event count = %d, want %d", got, want)
	}
	for i, event := range events {
		if got, want := event.Sequence, int64(i+1); got != want {
			t.Fatalf("event %d sequence = %d, want %d", i, got, want)
		}
	}
}

func TestSQLiteEventLogReturnsExistingSequenceForDuplicateEventID(t *testing.T) {
	ctx := context.Background()
	log := newTestEventLog(t, ctx)

	original := testEvent("event-duplicate", "session-a", eventsystem.EventSessionStart)
	duplicate := testEvent("event-duplicate", "session-a", eventsystem.EventSessionEnd)

	seq1, err := log.Append(ctx, original)
	if err != nil {
		t.Fatalf("append original: %v", err)
	}
	seq2, err := log.Append(ctx, duplicate)
	if err != nil {
		t.Fatalf("append duplicate: %v", err)
	}
	if got, want := seq2, seq1; got != want {
		t.Fatalf("duplicate sequence = %d, want existing %d", got, want)
	}

	events, err := log.Read(ctx, "session-a", 0)
	if err != nil {
		t.Fatalf("read events: %v", err)
	}
	if got, want := len(events), 1; got != want {
		t.Fatalf("event count after duplicate append = %d, want %d", got, want)
	}
	if got, want := events[0].EventType, original.EventType; got != want {
		t.Fatalf("stored event type = %s, want original %s", got, want)
	}
}

func TestSQLiteEventLogReadsFromSequenceInOrder(t *testing.T) {
	ctx := context.Background()
	log := newTestEventLog(t, ctx)

	for i := 1; i <= 5; i++ {
		if _, err := log.Append(ctx, testEvent(fmt.Sprintf("event-%d", i), "session-a", eventsystem.EventType(fmt.Sprintf("TestEvent%d", i)))); err != nil {
			t.Fatalf("append event %d: %v", i, err)
		}
	}

	events, err := log.Read(ctx, "session-a", 3)
	if err != nil {
		t.Fatalf("read events: %v", err)
	}
	gotSequences := make([]int64, 0, len(events))
	for _, event := range events {
		gotSequences = append(gotSequences, event.Sequence)
	}
	if want := []int64{3, 4, 5}; !equalInt64s(gotSequences, want) {
		t.Fatalf("read sequences = %v, want %v", gotSequences, want)
	}
}

func TestSQLiteEventLogRoundTripsEventEnvelopeFields(t *testing.T) {
	ctx := context.Background()
	log := newTestEventLog(t, ctx)

	event := testEvent("event-roundtrip", "session-a", eventsystem.EventPostModelCall)
	event.AgentID = "agent-roundtrip"
	event.TaskID = "task-roundtrip"
	event.TeamID = "team-roundtrip"
	event.CorrelationID = "corr-roundtrip"
	event.CausationID = "cause-roundtrip"
	event.Timestamp = time.Date(2026, 6, 22, 12, 34, 56, 789, time.UTC)
	event.Payload = json.RawMessage(`{"message":"ok","tokens":3}`)

	if _, err := log.Append(ctx, event); err != nil {
		t.Fatalf("append event: %v", err)
	}

	events, err := log.Read(ctx, "session-a", 1)
	if err != nil {
		t.Fatalf("read events: %v", err)
	}
	if got, want := len(events), 1; got != want {
		t.Fatalf("event count = %d, want %d", got, want)
	}
	got := events[0]
	if got.EventID != event.EventID ||
		got.EventType != event.EventType ||
		!got.Timestamp.Equal(event.Timestamp) ||
		got.SchemaVersion != event.SchemaVersion ||
		got.SessionID != event.SessionID ||
		got.AgentID != event.AgentID ||
		got.TaskID != event.TaskID ||
		got.TeamID != event.TeamID ||
		got.CorrelationID != event.CorrelationID ||
		got.CausationID != event.CausationID ||
		string(got.Payload) != string(event.Payload) {
		t.Fatalf("round-tripped event = %#v, want fields from %#v", got, event)
	}
}

func TestSQLiteEventLogSequencesAreIndependentPerSession(t *testing.T) {
	ctx := context.Background()
	log := newTestEventLog(t, ctx)

	seqA1, err := log.Append(ctx, testEvent("a-1", "session-a", eventsystem.EventSessionStart))
	if err != nil {
		t.Fatalf("append session-a first: %v", err)
	}
	seqB1, err := log.Append(ctx, testEvent("b-1", "session-b", eventsystem.EventSessionStart))
	if err != nil {
		t.Fatalf("append session-b first: %v", err)
	}
	seqA2, err := log.Append(ctx, testEvent("a-2", "session-a", eventsystem.EventSessionEnd))
	if err != nil {
		t.Fatalf("append session-a second: %v", err)
	}

	if got, want := []int64{seqA1, seqB1, seqA2}, []int64{1, 1, 2}; !equalInt64s(got, want) {
		t.Fatalf("sequences = %v, want %v", got, want)
	}
}

func TestSQLiteEventLogConcurrentAppendsStayContiguous(t *testing.T) {
	ctx := context.Background()
	log := newTestEventLog(t, ctx)

	const count = 32
	var wg sync.WaitGroup
	errs := make(chan error, count)
	sequences := make(chan int64, count)
	for i := 0; i < count; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			seq, err := log.Append(ctx, testEvent(fmt.Sprintf("event-%02d", i), "session-a", eventsystem.EventUserPromptSubmit))
			if err != nil {
				errs <- err
				return
			}
			sequences <- seq
		}()
	}
	wg.Wait()
	close(errs)
	close(sequences)

	for err := range errs {
		t.Fatalf("append error: %v", err)
	}

	got := make([]int64, 0, count)
	for seq := range sequences {
		got = append(got, seq)
	}
	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
	want := make([]int64, 0, count)
	for i := 1; i <= count; i++ {
		want = append(want, int64(i))
	}
	if !equalInt64s(got, want) {
		t.Fatalf("concurrent sequences = %v, want %v", got, want)
	}
}

func newTestEventLog(t *testing.T, ctx context.Context) *SQLiteEventLog {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "events.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close sqlite: %v", err)
		}
	})

	log, err := NewSQLiteEventLog(ctx, db)
	if err != nil {
		t.Fatalf("new event log: %v", err)
	}
	return log
}

func testEvent(eventID string, sessionID string, eventType eventsystem.EventType) eventsystem.Event {
	payload, err := json.Marshal(map[string]string{"event_id": eventID})
	if err != nil {
		panic(err)
	}
	return eventsystem.Event{
		EventID:       eventID,
		EventType:     eventType,
		Timestamp:     time.Date(2026, 6, 22, 10, 0, 0, 0, time.UTC),
		SchemaVersion: 1,
		SessionID:     sessionID,
		AgentID:       "agent-a",
		TaskID:        "task-a",
		TeamID:        "team-a",
		CorrelationID: "corr-a",
		CausationID:   "cause-a",
		Sequence:      999,
		Payload:       payload,
	}
}

func equalInt64s(got []int64, want []int64) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
