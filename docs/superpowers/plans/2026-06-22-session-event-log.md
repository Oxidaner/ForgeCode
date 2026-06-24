# Session Event Log Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver `FC-SESS-001` Append-only Event Store with SQLite-backed sequence allocation, EventID idempotency, and ordered reads.

**Architecture:** Add a focused `SQLiteEventLog` in `internal/session-store` that owns only durable event storage. It reuses the `eventsystem.Event` envelope, stores each envelope field in SQLite, allocates per-session sequence numbers inside a serialized transaction, and returns existing sequences for duplicate EventID appends.

**Tech Stack:** Go standard library `database/sql`, `modernc.org/sqlite v1.34.5`, existing `event-system` package.

---

### Task 1: EventLog Behavior

**Files:**
- Create: `internal/session-store/event_log_test.go`
- Create: `internal/session-store/event_log.go`
- Modify: `docs/tasks/session-store/TASKS.md`
- Modify: `docs/checklists/session-store/CHECKLIST.md`

- [x] Write failing tests for monotonic per-session Sequence, duplicate EventID idempotency, ordered `Read(sessionID, from)`, cross-session independent sequences, and concurrent append race safety.
- [x] Run `go test ./internal/session-store` and verify tests fail because `SQLiteEventLog` symbols do not exist.
- [x] Implement `EventLog`, `SQLiteEventLog`, schema initialization, append transaction, duplicate lookup, and ordered read mapping.
- [x] Run `go test ./internal/session-store` and verify tests pass.
- [x] Update `FC-SESS-001` task status, acceptance criteria, checklist evidence, and plan checkboxes.

### Task 2: Verification

**Files:**
- All files touched above.

- [x] Run `gofmt` on new Go files.
- [x] Run `go test ./...`.
- [x] Run `go test -race ./...`.
- [x] Review `git status --short` and `git diff --stat`.
