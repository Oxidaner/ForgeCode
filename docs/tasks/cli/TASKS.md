# cli Tasks

模块：`cli`。Task 前缀 `FC-CLI`。相关需求 FR-CLI-001..004。

## FC-CLI-001 — REPL 入口与任务提交
| Type | Architecture | Priority | P0 | Milestone | M1 | Status | Ready | Size | M |
| Dependencies | FC-RT-002 | Related Requirements | FR-CLI-001 | Spec | §6 |

**Description**：App/REPL 骨架，启动 Session、提交任务、订阅流式事件渲染。
**Implementation Notes**：CLI 框架见 OPEN_QUESTIONS Q5。
**Files**：`cmd/forge/`, `internal/cli/`。
**Tests Required**：输入解析 Unit、提交 Integration。
**Acceptance Criteria**：
- [ ] 可提交任务并流式渲染
- [ ] 任务 vs 命令正确区分
**Definition of Done**：与 runtime-core 跑通只读任务。

## FC-CLI-002 — 流式渲染与工具活动展示
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Backlog | Size | M |
| Dependencies | FC-CLI-001 | Related Requirements | FR-CLI-001 |

**Description**：Renderer.Stream/ToolActivity，渲染模型输出与工具调用/结果摘要。
**Acceptance Criteria**：
- [ ] 流式增量渲染
- [ ] 工具调用可见且不泄露敏感参数

## FC-CLI-003 — 审批交互
| Type | Implementation | Priority | P0 | Milestone | M2 | Status | Backlog | Size | M |
| Dependencies | FC-CLI-001, FC-PERM-007 | Related Requirements | FR-CLI-001, FR-PERM-007 |

**Description**：ApprovalPrompt 展示风险等级/命中原因，收集 Allow/Deny 回传。
**Security Considerations**：脱敏原始敏感参数。
**Acceptance Criteria**：
- [ ] 审批信息充分且脱敏
- [ ] 决策正确回传 runtime-core

## FC-CLI-004 — 内置固定逻辑命令
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Backlog | Size | M |
| Dependencies | FC-CLI-001, FC-CMD-001 | Related Requirements | FR-CLI-002, FR-CLI-003 |

**Description**：/help /model /context /cost /clear /exit /resume /checkpoint 本地执行或分派，不经模型。
**Security Considerations**：命令字符串不注入模型。
**Acceptance Criteria**：
- [ ] 固定命令本地执行
- [ ] 业务逻辑在 extension-system/runtime-core 而非 cli（FR-CLI-003）

## FC-CLI-005 — 取消与中断处理
| Type | Implementation | Priority | P0 | Milestone | M4 | Status | Backlog | Size | S |
| Dependencies | FC-CLI-001, FC-RT-006 | Related Requirements | FR-RUNTIME-003 |

**Description**：Ctrl-C → context 取消 → runtime-core.Cancel，安全终止渲染。
**Acceptance Criteria**：
- [ ] Ctrl-C 传播取消
- [ ] 终止后回到 Idle

## FC-CLI-006 — Resume / Checkpoint / Rewind 入口
| Type | Implementation | Priority | P1 | Milestone | M4 | Status | Backlog | Size | M |
| Dependencies | FC-CLI-001, FC-SESS-003 | Related Requirements | FR-CLI-004 |

**Description**：列出 Session/Checkpoint，/resume /checkpoint /rewind /init。
**Acceptance Criteria**：
- [ ] 可恢复历史 Session
- [ ] /rewind 回退 Checkpoint

## FC-CLI-007 — cli 测试与依赖约束
| Type | Test | Priority | P0 | Milestone | M2 | Status | Backlog | Size | S |
| Dependencies | FC-CLI-004 | Related Requirements | FR-CLI-003 |

**Description**：输入解析 Unit、提交 Integration、渲染 Golden、依赖检查（无核心业务逻辑）。
**Acceptance Criteria**：
- [ ] cli 不 import 实现核心逻辑（仅接口）
- [ ] `go test` 通过
