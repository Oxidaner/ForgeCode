# mcp-client Tasks

模块：`mcp-client`。Task 前缀 `FC-MCP`。相关需求 FR-MCP-001..005。ADR-0004。RISK-011。

## FC-MCP-001 — MCP Client 与 Server 生命周期
| Type | Architecture | Priority | P1 | Milestone | M7 | Status | Ready | Size | L |
| Dependencies | FC-TOOL-006 | Related Requirements | FR-MCP-001 | Spec | §6 |

**Description**：MCPClient/ServerHandle、生命周期（MCP Server State）、能力协商、健康检查、重连。
**Files**：`internal/mcp-client/`。
**Tests Required**：生命周期 Unit、重连 Failure Injection。
**Acceptance Criteria**：
- [ ] 状态机按 GLOSSARY 转移
- [ ] 断开不阻塞其他 Server
**Definition of Done**：对 Mock MCP Server 跑通生命周期。

## FC-MCP-002 — stdio 与 Streamable HTTP transport
| Type | Implementation | Priority | P1 | Milestone | M7 | Status | Backlog | Size | L |
| Dependencies | FC-MCP-001 | Related Requirements | FR-MCP-002 |

**Description**：两种 transport 实现，凭证不入日志。
**Tests Required**：Contract Test（对 Mock Server）。
**Acceptance Criteria**：
- [ ] 两种 transport 通过同套 Contract Test
- [ ] HTTP 鉴权凭证脱敏

## FC-MCP-003 — tools/resources/prompts 接入
| Type | Implementation | Priority | P1 | Milestone | M7 | Status | Backlog | Size | M |
| Dependencies | FC-MCP-002 | Related Requirements | FR-MCP-003 |

**Description**：tools/list+call、resources/list+read、prompts/list+get。
**Acceptance Criteria**：
- [ ] 三类能力可用
- [ ] 外部 Prompt/Resource 标注来源边界

## FC-MCP-004 — 统一接入、Namespace、信任与权限
| Type | Security | Priority | P0 | Milestone | M7 | Status | Backlog | Size | L |
| Dependencies | FC-MCP-003, FC-TOOL-006, FC-PERM-001 | Related Requirements | FR-MCP-004 |

**Description**：Schema 转统一 Descriptor、Namespace、名称冲突、输出大小限制、信任级别、每工具权限等级、审计。
**Security Considerations**：MCP 工具不绕过权限（RISK-011）。
**Tests Required**：Security Test（越权/冒充/超大输出）。
**Acceptance Criteria**：
- [ ] MCP 工具经 tool-runtime 管线
- [ ] 不冒充内置工具
- [ ] 不可信 Server 工具更严格权限

## FC-MCP-005 — 外部 Prompt/Resource 安全边界
| Type | Security | Priority | P0 | Milestone | M7 | Status | Backlog | Size | M |
| Dependencies | FC-MCP-003 | Related Requirements | FR-MCP-005 |

**Description**：外部 Prompt/Resource 作为不可信输入隔离，注入上下文时标注，防 Prompt Injection。
**Acceptance Criteria**：
- [ ] 外部内容不作为可信指令
- [ ] 来源边界清晰

## FC-MCP-006 — mcp-client 测试套件
| Type | Test | Priority | P1 | Milestone | M7 | Status | Backlog | Size | M |
| Dependencies | FC-MCP-004, FC-MCP-005 | Related Requirements | FR-MCP-001..005 |

**Description**：Contract/Integration/Security/Failure 汇总（对 Mock MCP Server）。
**Acceptance Criteria**：
- [ ] Security 测试通过（RISK-011）
- [ ] `go test -race` 通过
