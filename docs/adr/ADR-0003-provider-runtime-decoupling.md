# ADR-0003：Provider 与 Runtime 解耦

| 字段 | 值 |
| --- | --- |
| Status | Accepted |
| Date | 2026-06-13 |
| Deciders | 架构组 |
| Related Requirements | FR-PROVIDER-001, FR-PROVIDER-006, NFR-MAINT-001 |
| Related Modules | model-provider, runtime-core, context-manager |

## Context
ForgeCode 必须 Model-Agnostic。若 Runtime 直接依赖某家 Provider 的请求/响应结构，将无法替换模型且违背项目核心价值。

## Decision
定义中立的 `Provider` 接口与中立消息/工具调用/usage 结构。各家适配器在 model-provider 内部完成 **provider-specific 消息转换**，私有 SDK/结构不得泄漏到 runtime-core。runtime-core 仅依赖接口（依赖反转，构造注入）。

## Alternatives Considered
- **直接用某 SDK 的类型贯穿全栈**：开发快但锁定供应商——拒绝。
- **OpenAI 格式作为通用格式**：仍是供应商耦合，且 Anthropic 等差异大——拒绝。

## Consequences
- 正面：可插拔 Provider、可 Mock 测试、Contract Test 保障一致性。
- 负面：需维护转换层与能力差异矩阵（RISK-004）。

## Security Impact
Provider 响应在转换层统一校验，便于注入防护集中处理。

## Operational Impact
新增 Provider 仅需实现适配器与通过 Contract Test。

## Revisit Conditions
若需深度使用某 Provider 私有能力且无法抽象，评估能力扩展点机制。
