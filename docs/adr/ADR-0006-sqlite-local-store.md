# ADR-0006：SQLite 作为初始本地状态存储

| 字段 | 值 |
| --- | --- |
| Status | Accepted |
| Date | 2026-06-13 |
| Deciders | 架构组 |
| Related Requirements | FR-SESSION-001, FR-MEMORY-002, NFR-PORT-001 |
| Related Modules | session-store, memory-system |

## Context
需要本地、零运维、支持事务与全文检索的存储，承载事件、Session、Checkpoint、记忆与 Usage。

## Decision
第一版统一使用 **SQLite**（含 FTS5 用于记忆检索），文件系统存放 Skill 包与大 Artifact。启用 WAL 以改善并发读。驱动选型见 OPEN_QUESTIONS Q2。

## Alternatives Considered
- **外部 DB（Postgres）**：增加运维与部署复杂度，违背本地优先——拒绝。
- **纯文件 + JSON 行**：缺事务与查询能力——拒绝（仅用于大对象）。

## Consequences
- 正面：零运维、可移植、事务、FTS5。
- 负面：单写者并发限制（RISK-017），需写串行化/分片；大体积需归档（RISK-016）。

## Security Impact
本地存储；敏感数据需脱敏与访问控制，记忆需项目隔离。

## Operational Impact
单文件便于备份；需迁移脚本支持 Schema 版本演进。

## Revisit Conditions
出现高并发多 Agent 写瓶颈或跨机共享需求时，评估外部存储后端。
