# tool-runtime Tasks

模块：`tool-runtime`。Task 前缀 `FC-TOOL`。相关需求 FR-TOOL-001..004。ADR-0004/0005。RISK-005。

## FC-TOOL-001 — Tool 接口、Descriptor 与 Registry
| Type | Architecture | Priority | P0 | Milestone | M1 | Status | Done | Size | M |
| Dependencies | - | Related Requirements | FR-TOOL-001 | Spec | §6 |

**Description**：定义 `Tool`、`ToolDescriptor`、`Registry`，含 Namespace 与命名冲突处理。
**Files**：`internal/tool-runtime/tool.go`, `registry.go`。
**Tests Required**：注册/发现/冲突 Unit。
**Acceptance Criteria**：
- [x] 命名冲突返回 ConflictError
- [x] Descriptor 含来源与 schema
**Definition of Done**：接口评审 + 测试通过。
**Evidence**：实现 `internal/tool-runtime` 的 `Tool`、`ToolDescriptor`、`Registry`、`ToolCall`、`ToolResult` 与错误分类；注册/发现/冲突/非法 schema 单测通过。`go build ./...`、`go test ./...`、`go vet ./...` 通过；race 因缺 `gcc` 未执行。

## FC-TOOL-002 — 统一调用管线 Invoker
| Type | Architecture | Priority | P0 | Milestone | M2 | Status | Backlog | Size | L |
| Dependencies | FC-TOOL-001, FC-PERM-001, FC-EVT-002 | Related Requirements | FR-TOOL-002 |

**Description**：实现 Validation→Permission→PreHook→Execute→PostHook→Audit，顺序不可跳过。
**Security Considerations**：单点强制权限（NFR-SEC-001）。
**Tests Required**：管线顺序 Integration + 绕过尝试 Security Test。
**Acceptance Criteria**：
- [ ] 缺任一阶段测试失败
- [ ] Deny 阻止执行
**Definition of Done**：与 permission-engine + 一个内置工具跑通。

## FC-TOOL-003 — 截断、超时与错误分类
| Type | Implementation | Priority | P0 | Milestone | M2 | Status | Backlog | Size | M |
| Dependencies | FC-TOOL-002 | Related Requirements | FR-TOOL-003, NFR-LIMIT-001 |

**Description**：输出硬截断（标注 Truncated）、执行超时、错误归类到 GLOSSARY。
**Tests Required**：截断/超时/取消 Failure Injection。
**Acceptance Criteria**：
- [ ] 超 max_output_bytes 截断
- [ ] 超时归 TimeoutError 并落 ToolFailure
**Definition of Done**：失败注入测试通过。

## FC-TOOL-004 — 审批结果回填与 ApprovalRequired 协议
| Type | Implementation | Priority | P0 | Milestone | M2 | Status | Backlog | Size | M |
| Dependencies | FC-TOOL-002 | Related Requirements | FR-PERM-007, FR-RUNTIME-001 |

**Description**：Decide 返回 ApprovalRequired 时上抛 runtime-core，批准后重入管线执行段。
**Acceptance Criteria**：
- [ ] 未批准不进入 Execute
- [ ] 批准后从 PreHook 续行
**Definition of Done**：与 runtime-core AwaitingApproval 集成。

## FC-TOOL-005 — 审计编排与事件点位
| Type | Implementation | Priority | P0 | Milestone | M2 | Status | Backlog | Size | S |
| Dependencies | FC-TOOL-002, FC-TEL-003 | Related Requirements | FR-PERM-007, FR-TELEMETRY-003 |

**Description**：在管线点位产生 PreToolUse/PostToolUse/ToolFailure/AuditRecorded，输入脱敏后审计。
**Acceptance Criteria**：
- [ ] 每次调用产生完整审计事件
- [ ] 敏感输入脱敏（NFR-SEC-002）

## FC-TOOL-006 — MCP 工具统一接入契约
| Type | Implementation | Priority | P1 | Milestone | M7 | Status | Backlog | Size | M |
| Dependencies | FC-TOOL-001, FC-TOOL-002 | Related Requirements | FR-TOOL-004, FR-MCP-004 |

**Description**：定义 mcp-client 注册 MCP 工具为统一 Descriptor 的契约与输出大小限制挂钩。
**Security Considerations**：MCP 工具不得绕过管线（RISK-011）。
**Acceptance Criteria**：
- [ ] MCP 工具经同一 Invoker
- [ ] Namespace 防冒充内置工具

## FC-TOOL-007 — Contract Test 套件
| Type | Test | Priority | P0 | Milestone | M2 | Status | Backlog | Size | M |
| Dependencies | FC-TOOL-003 | Related Requirements | FR-TOOL-001..004 |

**Description**：对任意 Tool 实现统一运行的管线行为契约测试（内置与 MCP 共用）。
**Acceptance Criteria**：
- [ ] 内置与 MCP 工具通过同套用例（RISK-005）
- [ ] `go test -race` 通过
