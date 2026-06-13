# ADR-0002：使用 Append-only Event Store

| 字段 | 值 |
| --- | --- |
| Status | Accepted |
| Date | 2026-06-13 |
| Deciders | 架构组 |
| Related Requirements | FR-SESSION-001, FR-SESSION-003, NFR-REL-001, NFR-COMPAT-001 |
| Related Modules | session-store, event-system, runtime-core |

## Context
需要可靠重建会话状态、提供审计且支持 Replay。可变状态表难以审计且易在崩溃时产生不一致。

## Decision
采用 **Append-only Event Store**（SQLite 表）。状态通过重放事件得到；当前态可缓存但以事件为真相源。事件带单调 `Sequence` 与全局 `EventID`，写入幂等去重，带 `SchemaVersion` 支持演进。

## Alternatives Considered
- **可变状态表 + 审计日志双写**：双写易不一致——拒绝。
- **纯快照**：丢失中间过程，无法精细恢复/审计——拒绝（快照作为 Checkpoint 优化补充）。

## Consequences
- 正面：可恢复、可审计、可 Replay、天然支持 Hook 订阅。
- 负面：事件无限增长风险（RISK-016），需归档/截断策略与 Checkpoint 加速恢复。

## Security Impact
审计不可篡改（append-only + 序号 + 去重），支撑 Audit Tampering 防护。

## Operational Impact
需归档策略；恢复时间随事件量增长，用 Checkpoint 缓解。

## Revisit Conditions
单 Session 事件量或恢复时间超出可接受阈值时引入快照压缩/分段。
