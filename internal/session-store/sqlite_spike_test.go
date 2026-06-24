package sessionstore

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func TestSQLiteDriverSupportsWALAndFTS5(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "forgecode-spike.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	var journalMode string
	if err := db.QueryRowContext(ctx, "PRAGMA journal_mode=WAL").Scan(&journalMode); err != nil {
		t.Fatalf("enable WAL: %v", err)
	}
	if got, want := strings.ToLower(journalMode), "wal"; got != want {
		t.Fatalf("journal_mode = %q, want %q", got, want)
	}

	if _, err := db.ExecContext(ctx, `
CREATE TABLE events (
    event_id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    sequence INTEGER NOT NULL,
    payload TEXT NOT NULL
);
CREATE UNIQUE INDEX events_session_sequence ON events(session_id, sequence);
CREATE VIRTUAL TABLE memory_fts USING fts5(content);
`); err != nil {
		t.Fatalf("create tables and fts5 table: %v", err)
	}

	if _, err := db.ExecContext(ctx, `INSERT INTO memory_fts(content) VALUES (?)`, "runtime recovery uses append only events"); err != nil {
		t.Fatalf("insert fts5 row: %v", err)
	}

	var content string
	if err := db.QueryRowContext(ctx, `SELECT content FROM memory_fts WHERE memory_fts MATCH ?`, "recovery").Scan(&content); err != nil {
		t.Fatalf("query fts5 row: %v", err)
	}
	if got, want := content, "runtime recovery uses append only events"; got != want {
		t.Fatalf("fts5 content = %q, want %q", got, want)
	}
}
