# ADR-0005：Permission Engine 独立于 Tool Executor

| 字段 | 值 |
| --- | --- |
| Status | Accepted |
| Date | 2026-06-13 |
| Deciders | 架构组 |
| Related Requirements | FR-PERM-001, FR-PERM-008, NFR-SEC-001 |
| Related Modules | permission-engine, tool-runtime, sandbox |

## Context
若权限判断散落在各工具或由执行器内联，将难以测试、易绕过、无法统一审计。master-plan §2.8 要求独立可测试的 Permission Engine。

## Decision
Permission Engine 为 **纯决策组件**：输入 `(ToolDescriptor, 参数, 上下文)`，输出 `Decision`（含命中原因），**不执行任何操作**。执行由 tool-runtime/sandbox 负责。五层逻辑均在引擎内、可独立单测。决策与执行分离，决策结果落审计事件。

## Alternatives Considered
- **执行器内联权限检查**：耦合、难测、易绕过——拒绝。
- **每个工具自带权限逻辑**：重复且不一致——拒绝。

## Consequences
- 正面：可独立测试、统一决策、消除"Permission Engine 直接执行命令"反模式。
- 负面：需在管线中强制调用，缺一即漏洞——由 tool-runtime 单点强制。

## Security Impact
集中决策点；Deny 优先与最严格生效规则统一实施；审批后执行前崩溃可重判。

## Operational Impact
权限策略可配置、可审计、可回放。

## Revisit Conditions
若性能成为瓶颈，评估决策缓存（保持纯函数语义）。
