# ADR-0001：采用事件驱动 Agent Runtime

| 字段 | 值 |
| --- | --- |
| Status | Accepted |
| Date | 2026-06-13 |
| Deciders | 架构组 |
| Related Requirements | FR-RUNTIME-001, FR-RUNTIME-003, FR-RUNTIME-006, NFR-REL-001 |
| Related Modules | runtime-core, event-system, session-store |

## Context
Agent Loop 若实现为不可恢复的无限循环，将无法支持暂停、取消、恢复、审计与并行。master-plan §2.2 明确要求显式状态机或事件驱动运行时。

## Decision
将 runtime-core 实现为 **显式状态机 + 事件驱动**：状态转移由事件触发并产生新事件，所有关键转移落 Append-only 事件。Agent 状态枚举见 GLOSSARY。Loop 的每一步（Thinking/ToolRequested/Observing 等）都是可观察、可持久化、可恢复的状态。

## Alternatives Considered
- **纯 ReAct 循环**：简单，但不可恢复、难审计、难并行——拒绝。
- **外部工作流引擎（如 Temporal）**：可恢复，但引入重型依赖且违背"自主实现控制平面"目标——拒绝。

## Consequences
- 正面：可恢复、可审计、可暂停/取消、便于 SubAgent/Team 复用。
- 负面：状态机与事件设计前期成本高，需严格保证转移合法性。
- 后续：需 Golden Test 覆盖状态转移；与 ADR-0002 配合。

## Security Impact
状态与决策可审计；审批后执行前崩溃可被检测并安全重判。

## Operational Impact
事件时间线可导出，便于排障与 Replay（evaluation）。

## Revisit Conditions
若事件量导致性能瓶颈且无法通过归档解决，重新评估状态持久化粒度。
