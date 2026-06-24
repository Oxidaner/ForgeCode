package sessionstore

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestSQLiteStoreCreatesAndGetsSession(t *testing.T) {
	ctx := context.Background()
	store := newTestSessionStore(t, ctx)
	createdAt := time.Date(2026, 6, 22, 13, 0, 0, 0, time.UTC)

	session, err := store.CreateSession(ctx, SessionMeta{
		ID:            "session-1",
		WorkspaceRoot: "/workspace/forge",
		Model:         "mock-large",
		UserTask:      "read repository",
		CreatedAt:     createdAt,
	})
	if err != nil {
		t.Fatalf("CreateSession returned error: %v", err)
	}
	if got, want := session.ID, "session-1"; got != want {
		t.Fatalf("session ID = %q, want %q", got, want)
	}
	if got, want := session.State, SessionStateActive; got != want {
		t.Fatalf("session state = %s, want %s", got, want)
	}

	got, err := store.GetSession(ctx, "session-1")
	if err != nil {
		t.Fatalf("GetSession returned error: %v", err)
	}
	if got.ID != session.ID ||
		got.Meta.WorkspaceRoot != "/workspace/forge" ||
		got.Meta.Model != "mock-large" ||
		got.Meta.UserTask != "read repository" ||
		!got.CreatedAt.Equal(createdAt) {
		t.Fatalf("got session = %#v, want created metadata %#v", got, session)
	}
}

func TestSQLiteStoreRejectsDuplicateSessionID(t *testing.T) {
	ctx := context.Background()
	store := newTestSessionStore(t, ctx)

	if _, err := store.CreateSession(ctx, SessionMeta{ID: "session-1"}); err != nil {
		t.Fatalf("CreateSession first returned error: %v", err)
	}
	_, err := store.CreateSession(ctx, SessionMeta{ID: "session-1"})
	if !errors.Is(err, ErrSessionConflict) {
		t.Fatalf("duplicate error = %v, want ErrSessionConflict", err)
	}
}

func TestSQLiteStoreUpdatesOnlyLegalStateTransitions(t *testing.T) {
	ctx := context.Background()
	store := newTestSessionStore(t, ctx)

	if _, err := store.CreateSession(ctx, SessionMeta{ID: "session-1"}); err != nil {
		t.Fatalf("CreateSession returned error: %v", err)
	}

	legal := []SessionState{
		SessionStatePaused,
		SessionStateActive,
		SessionStateCompleted,
	}
	for _, state := range legal {
		if err := store.UpdateState(ctx, "session-1", state); err != nil {
			t.Fatalf("UpdateState(%s) returned error: %v", state, err)
		}
	}

	session, err := store.GetSession(ctx, "session-1")
	if err != nil {
		t.Fatalf("GetSession returned error: %v", err)
	}
	if got, want := session.State, SessionStateCompleted; got != want {
		t.Fatalf("final session state = %s, want %s", got, want)
	}

	err = store.UpdateState(ctx, "session-1", SessionStateActive)
	if !errors.Is(err, ErrInvalidSessionTransition) {
		t.Fatalf("terminal transition error = %v, want ErrInvalidSessionTransition", err)
	}
}

func TestSQLiteStoreRejectsInvalidSessionTransition(t *testing.T) {
	ctx := context.Background()
	store := newTestSessionStore(t, ctx)

	if _, err := store.CreateSession(ctx, SessionMeta{ID: "session-1"}); err != nil {
		t.Fatalf("CreateSession returned error: %v", err)
	}

	err := store.UpdateState(ctx, "session-1", SessionStateCompleted)
	if err != nil {
		t.Fatalf("complete active session: %v", err)
	}
	err = store.UpdateState(ctx, "session-1", SessionStatePaused)
	if !errors.Is(err, ErrInvalidSessionTransition) {
		t.Fatalf("completed->paused error = %v, want ErrInvalidSessionTransition", err)
	}
}

func TestSQLiteStoreListsSessionsByState(t *testing.T) {
	ctx := context.Background()
	store := newTestSessionStore(t, ctx)

	for _, id := range []string{"active-1", "paused-1", "completed-1", "failed-1"} {
		if _, err := store.CreateSession(ctx, SessionMeta{ID: id}); err != nil {
			t.Fatalf("CreateSession(%s): %v", id, err)
		}
	}
	if err := store.UpdateState(ctx, "paused-1", SessionStatePaused); err != nil {
		t.Fatalf("pause session: %v", err)
	}
	if err := store.UpdateState(ctx, "completed-1", SessionStateCompleted); err != nil {
		t.Fatalf("complete session: %v", err)
	}
	if err := store.UpdateState(ctx, "failed-1", SessionStateFailed); err != nil {
		t.Fatalf("fail session: %v", err)
	}

	sessions, err := store.ListSessions(ctx, SessionFilter{
		States: []SessionState{SessionStateActive, SessionStatePaused},
	})
	if err != nil {
		t.Fatalf("ListSessions returned error: %v", err)
	}
	if got, want := sessionIDs(sessions), []string{"active-1", "paused-1"}; !equalStrings(got, want) {
		t.Fatalf("listed session IDs = %v, want %v", got, want)
	}
}

func TestSQLiteStoreReturnsNotFoundForMissingSession(t *testing.T) {
	ctx := context.Background()
	store := newTestSessionStore(t, ctx)

	_, err := store.GetSession(ctx, "missing")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("missing get error = %v, want ErrSessionNotFound", err)
	}
	err = store.UpdateState(ctx, "missing", SessionStateCancelled)
	if !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("missing update error = %v, want ErrSessionNotFound", err)
	}
}

func newTestSessionStore(t *testing.T, ctx context.Context) *SQLiteStore {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "sessions.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close sqlite: %v", err)
		}
	})

	store, err := NewSQLiteStore(ctx, db)
	if err != nil {
		t.Fatalf("new sqlite store: %v", err)
	}
	return store
}

func sessionIDs(sessions []Session) []string {
	out := make([]string, 0, len(sessions))
	for _, session := range sessions {
		out = append(out, session.ID)
	}
	return out
}

func equalStrings(got []string, want []string) bool {
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
