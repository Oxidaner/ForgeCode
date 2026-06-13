# ADR-0009：Agent Team 使用中心化 Task DAG

| 字段 | 值 |
| --- | --- |
| Status | Accepted |
| Date | 2026-06-13 |
| Deciders | 架构组 |
| Related Requirements | FR-TEAM-001, FR-TEAM-002 |
| Related Modules | agent-orchestration |

## Context
多 Agent 协作若采用去中心化自由协商/辩论，复杂度与不可预测性极高（RISK-015），难以保证可控与可审计。

## Decision
第一版 Agent Team 采用 **中心化调度**：Team Lead 维护 **Task DAG**，按依赖将 Ready 任务分配给成员，通过 Mailbox 定向/广播通信、Artifact Store 共享产物，最终由 Lead 集成结果。不做去中心化协商或自由辩论。Task/Team 状态见 GLOSSARY。

## Alternatives Considered
- **去中心化协商/辩论**：不可控、难审计、成本高——拒绝（第一版）。
- **纯串行 SubAgent**：无并行协作能力——作为降级方案保留。

## Consequences
- 正面：可控、可审计、调度可测试、失败可定位。
- 负面：Lead 成为复杂度与单点；DAG 死锁/饥饿需防护。

## Security Impact
集中调度便于统一权限与审计；成员仍走统一工具管线，不绕过权限。

## Operational Impact
需 Team Budget、超时、人工介入与结果集成机制。

## Revisit Conditions
当中心化调度成为瓶颈且有明确协作模式需求时，评估更灵活的协作模型。
