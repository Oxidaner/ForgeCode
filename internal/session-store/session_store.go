package sessionstore

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrSessionConflict          = errors.New("session-store: session already exists")
	ErrSessionNotFound          = errors.New("session-store: session not found")
	ErrInvalidSessionTransition = errors.New("session-store: invalid session state transition")
)

type Store interface {
	CreateSession(ctx context.Context, meta SessionMeta) (*Session, error)
	GetSession(ctx context.Context, id string) (*Session, error)
	ListSessions(ctx context.Context, f SessionFilter) ([]Session, error)
	UpdateState(ctx context.Context, id string, s SessionState) error
}

type SessionState string

const (
	SessionStateActive    SessionState = "Active"
	SessionStatePaused    SessionState = "Paused"
	SessionStateCompleted SessionState = "Completed"
	SessionStateFailed    SessionState = "Failed"
	SessionStateCancelled SessionState = "Cancelled"
)

type SessionMeta struct {
	ID            string
	WorkspaceRoot string
	Model         string
	UserTask      string
	CreatedAt     time.Time
}

type Session struct {
	ID        string
	State     SessionState
	Meta      SessionMeta
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SessionFilter struct {
	States []SessionState
	Limit  int
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(ctx context.Context, db *sql.DB) (*SQLiteStore, error) {
	if db == nil {
		return nil, errors.New("session-store: nil sqlite database")
	}
	if err := initSessionSchema(ctx, db); err != nil {
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) CreateSession(ctx context.Context, meta SessionMeta) (*Session, error) {
	if meta.ID == "" {
		meta.ID = newSessionID()
	}
	if meta.CreatedAt.IsZero() {
		meta.CreatedAt = time.Now().UTC()
	}
	meta.CreatedAt = meta.CreatedAt.UTC()

	_, err := s.db.ExecContext(ctx, `
INSERT INTO sessions (
    session_id, state, workspace_root, model, user_task, created_at_utc, updated_at_utc
) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		meta.ID,
		string(SessionStateActive),
		meta.WorkspaceRoot,
		meta.Model,
		meta.UserTask,
		formatSessionTime(meta.CreatedAt),
		formatSessionTime(meta.CreatedAt),
	)
	if err != nil {
		if isSQLiteConstraint(err) {
			return nil, ErrSessionConflict
		}
		return nil, fmt.Errorf("session-store: create session: %w", err)
	}
	return s.GetSession(ctx, meta.ID)
}

func (s *SQLiteStore) GetSession(ctx context.Context, id string) (*Session, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT session_id, state, workspace_root, model, user_task, created_at_utc, updated_at_utc
FROM sessions
WHERE session_id = ?`, id)
	session, err := scanSession(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return session, nil
}

func (s *SQLiteStore) ListSessions(ctx context.Context, f SessionFilter) ([]Session, error) {
	query := `
SELECT session_id, state, workspace_root, model, user_task, created_at_utc, updated_at_utc
FROM sessions`
	args := make([]any, 0, len(f.States))
	if len(f.States) > 0 {
		placeholders := make([]string, 0, len(f.States))
		for _, state := range f.States {
			placeholders = append(placeholders, "?")
			args = append(args, string(state))
		}
		query += " WHERE state IN (" + strings.Join(placeholders, ", ") + ")"
	}
	query += " ORDER BY created_at_utc ASC, session_id ASC"
	if f.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, f.Limit)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("session-store: list sessions: %w", err)
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		session, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, *session)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("session-store: iterate sessions: %w", err)
	}
	return sessions, nil
}

func (s *SQLiteStore) UpdateState(ctx context.Context, id string, next SessionState) error {
	current, err := s.GetSession(ctx, id)
	if err != nil {
		return err
	}
	if !CanTransitionSession(current.State, next) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidSessionTransition, current.State, next)
	}

	result, err := s.db.ExecContext(ctx, `
UPDATE sessions
SET state = ?, updated_at_utc = ?
WHERE session_id = ?`, string(next), formatSessionTime(time.Now().UTC()), id)
	if err != nil {
		return fmt.Errorf("session-store: update session state: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("session-store: inspect session update: %w", err)
	}
	if affected == 0 {
		return ErrSessionNotFound
	}
	return nil
}

func CanTransitionSession(from, to SessionState) bool {
	if from == to {
		return true
	}
	switch from {
	case SessionStateActive:
		return to == SessionStatePaused ||
			to == SessionStateCompleted ||
			to == SessionStateFailed ||
			to == SessionStateCancelled
	case SessionStatePaused:
		return to == SessionStateActive ||
			to == SessionStateCancelled
	default:
		return false
	}
}

func initSessionSchema(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `PRAGMA journal_mode=WAL`); err != nil {
		return fmt.Errorf("session-store: enable WAL: %w", err)
	}
	_, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS sessions (
    session_id TEXT PRIMARY KEY,
    state TEXT NOT NULL,
    workspace_root TEXT NOT NULL,
    model TEXT NOT NULL,
    user_task TEXT NOT NULL,
    created_at_utc TEXT NOT NULL,
    updated_at_utc TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS sessions_state_created
    ON sessions(state, created_at_utc ASC, session_id ASC);
`)
	if err != nil {
		return fmt.Errorf("session-store: initialize session schema: %w", err)
	}
	return nil
}

type sessionScanner interface {
	Scan(dest ...any) error
}

func scanSession(row sessionScanner) (*Session, error) {
	var session Session
	var state string
	var createdAt string
	var updatedAt string
	err := row.Scan(
		&session.ID,
		&state,
		&session.Meta.WorkspaceRoot,
		&session.Meta.Model,
		&session.Meta.UserTask,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("session-store: scan session: %w", err)
	}
	parsedCreatedAt, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, fmt.Errorf("session-store: parse session created_at: %w", err)
	}
	parsedUpdatedAt, err := time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("session-store: parse session updated_at: %w", err)
	}
	session.State = SessionState(state)
	session.Meta.ID = session.ID
	session.Meta.CreatedAt = parsedCreatedAt
	session.CreatedAt = parsedCreatedAt
	session.UpdatedAt = parsedUpdatedAt
	return &session, nil
}

func isSQLiteConstraint(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "constraint")
}

func formatSessionTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

func newSessionID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err == nil {
		return hex.EncodeToString(bytes[:])
	}
	return fmt.Sprintf("session-%d", time.Now().UnixNano())
}
