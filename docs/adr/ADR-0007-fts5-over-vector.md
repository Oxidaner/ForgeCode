# ADR-0007：记忆先使用 FTS5 而非向量数据库

| 字段 | 值 |
| --- | --- |
| Status | Accepted |
| Date | 2026-06-13 |
| Deciders | 架构组 |
| Related Requirements | FR-MEMORY-002, FR-MEMORY-001 |
| Related Modules | memory-system |

## Context
跨会话记忆需要检索。引入向量数据库会增加依赖、运维与 Embedding 成本，第一版不一定必要。

## Decision
记忆检索 **第一版仅使用 SQLite FTS5 关键词检索**。不默认引入向量数据库/Embedding。仅当有证据表明关键词检索不足且存在真实需求时，才将语义检索列入后续阶段。

## Alternatives Considered
- **向量数据库 + Embedding**：能力强但成本/依赖高，过早优化——推迟。
- **无检索、全量注入**：上下文爆炸——拒绝。

## Consequences
- 正面：零额外依赖、可解释、低成本。
- 负面：语义相近但用词不同的记忆可能漏检——通过良好标签/字段缓解。

## Security Impact
无外部 Embedding 服务，避免敏感记忆外泄。

## Operational Impact
FTS5 随 SQLite 内置，无额外组件。

## Revisit Conditions
当关键词检索召回率被 Eval 证明不足且影响任务效果时，评估本地 Embedding 方案。
