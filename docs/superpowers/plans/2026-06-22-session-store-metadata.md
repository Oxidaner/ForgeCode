# Session Store Metadata Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver `FC-SESS-002` with SQLite-backed Session metadata, Session State enum, legal state transitions, and list filters for resumable sessions.

**Architecture:** Add a focused `SQLiteStore` in `internal/session-store` for session metadata. It initializes a `sessions` table, creates sessions in `Active`, reads by ID, lists by state filters, and updates state only through legal transitions defined by the session-store spec.

**Tech Stack:** Go standard library `database/sql`, `modernc.org/sqlite v1.34.5`, existing `sessionstore` package.

---

### Task 1: Session Store Behavior

**Files:**
- Create: `internal/session-store/session_store_test.go`
- Create: `internal/session-store/session_store.go`
- Modify: `docs/tasks/session-store/TASKS.md`
- Modify: `docs/checklists/session-store/CHECKLIST.md`
- Modify: `docs/planning/TRACEABILITY_MATRIX.md`

- [x] Write failing tests for `CreateSession`, `GetSession`, `ListSessions`, legal/illegal `UpdateState`, terminal-state protection, duplicate session ID conflicts, and listing unfinished sessions.
- [x] Run `go test ./internal/session-store` and verify tests fail because `SQLiteStore` symbols do not exist.
- [x] Implement `Store`, `Session`, `SessionMeta`, `SessionFilter`, `SessionState`, `SQLiteStore`, schema initialization, state validation, and typed errors.
- [x] Run `go test ./internal/session-store` and verify tests pass.
- [x] Update `FC-SESS-002` task status, checklist evidence, traceability status, and plan checkboxes.

### Task 2: Verification

**Files:**
- All files touched above.

- [x] Run `gofmt` on new Go files.
- [x] Run `go test ./...`.
- [x] Run `go test -race ./...`.
- [x] Review `git status --short` and `git diff --stat`.
