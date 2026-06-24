# Mock Provider Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver `FC-PROV-002` so runtime tests can drive deterministic model responses without real provider dependencies.

**Architecture:** Add a script-driven `MockProvider` inside `internal/model-provider`, reusing the neutral `Provider` contract from `FC-PROV-001`. The mock consumes `MockStep` values in order, records requests, supports response/error/delay/stream chunks, and exposes deterministic state through copy-returning accessors.

**Tech Stack:** Go standard library only; existing `modelprovider` package.

---

### Task 1: Mock Provider Unit Behavior

**Files:**
- Create: `internal/model-provider/mock_test.go`
- Create: `internal/model-provider/mock.go`
- Modify: `docs/tasks/model-provider/TASKS.md`
- Modify: `docs/checklists/model-provider/CHECKLIST.md`

- [x] Write failing tests for scripted response replay, multi Tool Call replay, scripted errors, context timeout through delay, request recording, capability metadata, and stream chunks.
- [x] Run `go test ./internal/model-provider` and verify tests fail because `MockProvider` symbols do not exist.
- [x] Implement `MockProvider`, `MockStep`, script exhaustion error, request accessors, capability configuration, and `mockStreamReader`.
- [x] Run `go test ./internal/model-provider` and verify tests pass.
- [x] Update `FC-PROV-002` task status, acceptance criteria, checklist evidence, and plan checkboxes.

### Task 2: Verification

**Files:**
- All files touched above.

- [x] Run `gofmt` on new Go files.
- [x] Run `go test ./...`.
- [x] Run `go test -race ./...`.
- [x] Review `git status --short` and `git diff --stat`.
