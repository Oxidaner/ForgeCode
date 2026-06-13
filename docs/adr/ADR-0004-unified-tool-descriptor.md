# ADR-0004：内置 Tool 与 MCP Tool 使用统一描述

| 字段 | 值 |
| --- | --- |
| Status | Accepted |
| Date | 2026-06-13 |
| Deciders | 架构组 |
| Related Requirements | FR-TOOL-001, FR-TOOL-004, FR-MCP-004 |
| Related Modules | tool-runtime, mcp-client, builtin-tools |

## Context
内置工具与 MCP 工具若各走独立调用路径，将导致权限/审计/截断逻辑重复且易出现绕过权限的缺口。

## Decision
定义统一 `ToolDescriptor`（名称、JSON Schema、风险标注、权限要求、来源/Namespace）与统一 `Tool` 接口。MCP 工具经转换层映射为同一 Descriptor，进入 **同一调用管线**（Validation→Permission→Hook→Execute→Audit）。

## Alternatives Considered
- **MCP 工具独立通道**：实现快，但产生权限绕过风险（反模式）——拒绝。

## Consequences
- 正面：单一权限/审计/截断路径；新增工具来源成本低。
- 负面：MCP schema 转换需处理差异与冲突（RISK-005），Namespace 防冲突。

## Security Impact
消除"MCP 工具绕过权限系统"的反模式；统一审计。

## Operational Impact
工具来源对 Runtime 透明，便于扩展。

## Revisit Conditions
若某类工具语义无法用统一 Descriptor 表达，评估扩展字段。
