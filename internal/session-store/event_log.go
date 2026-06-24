package sessionstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	eventsystem "github.com/Oxidaner/ForgeCode/internal/event-system"
)

type EventLog interface {
	Append(ctx context.Context, e eventsystem.Event) (seq int64, err error)
	Read(ctx context.Context, sessionID string, from int64) ([]eventsystem.Event, error)
}

type SQLiteEventLog struct {
	db *sql.DB
	mu sync.Mutex
}

func NewSQLiteEventLog(ctx context.Context, db *sql.DB) (*SQLiteEventLog, error) {
	if db == nil {
		return nil, errors.New("session-store: nil sqlite database")
	}
	if err := initEventLogSchema(ctx, db); err != nil {
		return nil, err
	}
	return &SQLiteEventLog{db: db}, nil
}

func (l *SQLiteEventLog) Append(ctx context.Context, e eventsystem.Event) (seq int64, err error) {
	if err := validateEventForAppend(e); err != nil {
		return 0, err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("session-store: begin append event transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	existing, found, err := lookupSequenceByEventID(ctx, tx, e.EventID)
	if err != nil {
		return 0, err
	}
	if found {
		if err = tx.Commit(); err != nil {
			return 0, fmt.Errorf("session-store: commit duplicate event lookup: %w", err)
		}
		return existing, nil
	}

	next, err := nextSequence(ctx, tx, e.SessionID)
	if err != nil {
		return 0, err
	}
	e.Sequence = next
	if err = insertEvent(ctx, tx, e); err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("session-store: commit append event: %w", err)
	}
	return next, nil
}

func (l *SQLiteEventLog) Read(ctx context.Context, sessionID string, from int64) ([]eventsystem.Event, error) {
	rows, err := l.db.QueryContext(ctx, `
SELECT event_id, event_type, timestamp_utc, schema_version, session_id, agent_id, task_id, team_id,
       correlation_id, causation_id, sequence, payload
FROM events
WHERE session_id = ? AND sequence >= ?
ORDER BY sequence ASC`, sessionID, from)
	if err != nil {
		return nil, fmt.Errorf("session-store: read events: %w", err)
	}
	defer rows.Close()

	var events []eventsystem.Event
	for rows.Next() {
		var event eventsystem.Event
		var eventType string
		var timestamp string
		var payload []byte
		if err := rows.Scan(
			&event.EventID,
			&eventType,
			&timestamp,
			&event.SchemaVersion,
			&event.SessionID,
			&event.AgentID,
			&event.TaskID,
			&event.TeamID,
			&event.CorrelationID,
			&event.CausationID,
			&event.Sequence,
			&payload,
		); err != nil {
			return nil, fmt.Errorf("session-store: scan event: %w", err)
		}
		parsedTime, err := time.Parse(time.RFC3339Nano, timestamp)
		if err != nil {
			return nil, fmt.Errorf("session-store: parse event timestamp: %w", err)
		}
		event.EventType = eventsystem.EventType(eventType)
		event.Timestamp = parsedTime
		event.Payload = json.RawMessage(append([]byte(nil), payload...))
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("session-store: iterate events: %w", err)
	}
	return events, nil
}

func initEventLogSchema(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `PRAGMA journal_mode=WAL`); err != nil {
		return fmt.Errorf("session-store: enable WAL: %w", err)
	}
	_, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS events (
    event_id TEXT PRIMARY KEY,
    event_type TEXT NOT NULL,
    timestamp_utc TEXT NOT NULL,
    schema_version INTEGER NOT NULL,
    session_id TEXT NOT NULL,
    agent_id TEXT NOT NULL,
    task_id TEXT NOT NULL,
    team_id TEXT NOT NULL,
    correlation_id TEXT NOT NULL,
    causation_id TEXT NOT NULL,
    sequence INTEGER NOT NULL,
    payload BLOB NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS events_session_sequence
    ON events(session_id, sequence);
CREATE INDEX IF NOT EXISTS events_session_read
    ON events(session_id, sequence ASC);
`)
	if err != nil {
		return fmt.Errorf("session-store: initialize event log schema: %w", err)
	}
	return nil
}

func validateEventForAppend(e eventsystem.Event) error {
	if e.EventID == "" {
		return errors.New("session-store: event_id is required")
	}
	if e.SessionID == "" {
		return errors.New("session-store: session_id is required")
	}
	if e.EventType == "" {
		return errors.New("session-store: event_type is required")
	}
	if e.SchemaVersion <= 0 {
		return errors.New("session-store: schema_version must be positive")
	}
	if !json.Valid(e.Payload) {
		return errors.New("session-store: payload must be valid JSON")
	}
	return nil
}

func lookupSequenceByEventID(ctx context.Context, tx *sql.Tx, eventID string) (int64, bool, error) {
	var seq int64
	err := tx.QueryRowContext(ctx, `SELECT sequence FROM events WHERE event_id = ?`, eventID).Scan(&seq)
	if err == nil {
		return seq, true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}
	return 0, false, fmt.Errorf("session-store: lookup event_id: %w", err)
}

func nextSequence(ctx context.Context, tx *sql.Tx, sessionID string) (int64, error) {
	var seq int64
	err := tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(sequence), 0) + 1 FROM events WHERE session_id = ?`, sessionID).Scan(&seq)
	if err != nil {
		return 0, fmt.Errorf("session-store: allocate sequence: %w", err)
	}
	return seq, nil
}

func insertEvent(ctx context.Context, tx *sql.Tx, e eventsystem.Event) error {
	_, err := tx.ExecContext(ctx, `
INSERT INTO events (
    event_id, event_type, timestamp_utc, schema_version, session_id, agent_id, task_id, team_id,
    correlation_id, causation_id, sequence, payload
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.EventID,
		string(e.EventType),
		e.Timestamp.UTC().Format(time.RFC3339Nano),
		e.SchemaVersion,
		e.SessionID,
		e.AgentID,
		e.TaskID,
		e.TeamID,
		e.CorrelationID,
		e.CausationID,
		e.Sequence,
		[]byte(e.Payload),
	)
	if err != nil {
		return fmt.Errorf("session-store: insert event: %w", err)
	}
	return nil
}
