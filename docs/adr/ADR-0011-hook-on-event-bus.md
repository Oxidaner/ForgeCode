# ADR-0011：Hook 基于统一 Event Bus

| 字段 | 值 |
| --- | --- |
| Status | Accepted |
| Date | 2026-06-13 |
| Deciders | 架构组 |
| Related Requirements | FR-HOOK-001, FR-HOOK-002, FR-HOOK-003, FR-EVENT-001 |
| Related Modules | extension-system, event-system |

## Context
若在每个模块散落硬编码回调，Hook 将不可统一管理、不可审计且易遗漏。

## Decision
Hook 基于 **统一 Event Bus 与统一事件模型** 实现。HookDispatcher 订阅生命周期事件（§2.7），按 Matcher/优先级/顺序执行 Internal Go / Shell / HTTP Hook，返回 `Allow/Deny/Ask/Modify/Continue`。对可阻断决策点（如 PreToolUse）Hook 参与决策但 **不可静默提权**，受 Timeout、失败策略、递归防护与审计约束。

## Alternatives Considered
- **各模块硬编码回调**：散乱、不可审计——拒绝（反模式）。

## Consequences
- 正面：统一、可审计、可配置、可扩展；Hook 与事件一致。
- 负面：决策冲突需明确合并规则；外部 Hook 安全风险（RISK-012）。

## Security Impact
Hook 不可绕过事件总线，不可提权；外部 Hook 受沙箱/超时约束并审计。

## Operational Impact
Hook 配置集中；失败策略默认安全敏感事件 fail-closed。

## Revisit Conditions
若同步 Hook 影响主循环延迟，评估异步/旁路 Hook 分类。
