# Runtime Owner First Batch Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Start Role A implementation by delivering the Go baseline, Event contract, Provider contract, Agent state machine, and SQLite/FTS5 spike decision.

**Architecture:** Keep module IDs as documented (`internal/event-system`, `internal/model-provider`, `internal/runtime-core`) while using Go-safe package names (`eventsystem`, `modelprovider`, `runtimecore`). Implement only contract-level code owned by Runtime Owner; integration with tool-runtime, permission-engine, and CLI remains out of scope for this batch.

**Tech Stack:** Go 1.26.1 baseline with `go 1.22` compatibility target, standard library for core code, SQLite spike via a pure-Go driver if dependency resolution succeeds.

---

### Task 1: Go Baseline

**Files:**
- Create: `go.mod`
- Create: `internal/runtime-core/baseline_test.go`
- Modify: `AGENTS.md`
- Modify: `docs/planning/OPEN_QUESTIONS.md`
- Modify: `docs/tasks/runtime-core/TASKS.md`
- Modify: `docs/checklists/runtime-core/CHECKLIST.md`

- [x] Write a failing baseline test proving the module is not initialized.
- [x] Run `go test ./internal/runtime-core` and verify it fails because no module/package exists.
- [x] Create `go.mod` with module path `github.com/Oxidaner/ForgeCode` and `go 1.22`.
- [x] Add a minimal generic/slog/errors.Join baseline test.
- [x] Run `go test ./internal/runtime-core` and verify it passes.
- [x] Update Q1 and FC-RT-000 evidence.

### Task 2: Event Contract

**Files:**
- Create: `internal/event-system/types.go`
- Create: `internal/event-system/event_test.go`
- Modify: `docs/tasks/event-system/TASKS.md`
- Modify: `docs/checklists/event-system/CHECKLIST.md`

- [x] Write failing contract tests for Event envelope fields, EventType coverage, and EventClass mapping.
- [x] Run `go test ./internal/event-system` and verify tests fail because package symbols do not exist.
- [x] Implement Event, EventType, EventClass, SubscriptionFilter, Bus, Subscriber, and Unsubscribe.
- [x] Run `go test ./internal/event-system` and verify tests pass.
- [x] Update FC-EVT-001 evidence and checklist.

### Task 3: Provider Contract

**Files:**
- Create: `internal/model-provider/types.go`
- Create: `internal/model-provider/provider_test.go`
- Modify: `docs/tasks/model-provider/TASKS.md`
- Modify: `docs/checklists/model-provider/CHECKLIST.md`

- [x] Write failing tests for neutral request/response JSON behavior, ToolChoice/StopReason enums, and ProviderError retry metadata.
- [x] Run `go test ./internal/model-provider` and verify tests fail because package symbols do not exist.
- [x] Implement Provider, StreamReader, neutral message/tool/usage/capability/error types, and enums.
- [x] Run `go test ./internal/model-provider` and verify tests pass.
- [x] Update FC-PROV-001 evidence and checklist.

### Task 4: Agent State Machine

**Files:**
- Create: `internal/runtime-core/state.go`
- Create: `internal/runtime-core/state_test.go`
- Modify: `docs/tasks/runtime-core/TASKS.md`
- Modify: `docs/checklists/runtime-core/CHECKLIST.md`

- [x] Write failing tests for all legal transitions from runtime-core spec §8 and a set of illegal transitions.
- [x] Run `go test ./internal/runtime-core` and verify state tests fail because symbols do not exist.
- [x] Implement AgentStateName, Transition, TransitionEvent, and StateMachine with data-driven transitions.
- [x] Run `go test ./internal/runtime-core` and verify tests pass.
- [x] Update FC-RT-001 evidence and checklist.

### Task 5: SQLite/FTS5 Spike

**Files:**
- Create or modify: `internal/session-store/sqlite_spike_test.go`
- Modify: `go.mod`
- Modify: `go.sum`
- Modify: `docs/planning/OPEN_QUESTIONS.md`
- Modify: `docs/tasks/session-store/TASKS.md`
- Modify: `docs/checklists/session-store/CHECKLIST.md`

- [x] Attempt a pure-Go SQLite FTS5 spike with WAL and a simple FTS query.
- [x] If dependency download is blocked by sandbox/network, request approval and retry.
- [x] Run `go test ./internal/session-store` and verify the spike result.
- [x] Update Q2 and FC-SESS-000 evidence with the observed driver decision.

### Task 6: Full Verification

**Files:**
- All files touched above.

- [x] Run `gofmt` on all Go files.
- [x] Run `go test ./...`.
- [x] Run `go test -race ./...`.
- [x] Run `git status --short` and review the final diff.

## Self-Review

- Spec coverage: maps to FC-RT-000, FC-EVT-001, FC-PROV-001, FC-RT-001, and FC-SESS-000.
- Scope exclusions: no tool execution, permission decisions, CLI, provider adapters, runtime coordinator, or production session store.
- No placeholders: each task has concrete files, commands, and expected failure/pass checkpoints.
