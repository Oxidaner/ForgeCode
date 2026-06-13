# agent-orchestration Tasks

模块：`agent-orchestration`。Task 前缀：`FC-SUB`（SubAgent）、`FC-TEAM`（Team）。相关需求 FR-SUBAGENT-001..004, FR-TEAM-001..003。ADR-0008/0009。RISK-013/015。

## FC-SUB-001 — AgentDefinition 与 SubAgentSpec
| Type | Architecture | Priority | P1 | Milestone | M10 | Status | Ready | Size | M |
| Dependencies | FC-RT-002 | Related Requirements | FR-SUBAGENT-001 | Spec | §6 |

**Description**：AgentDefinition、SubAgentSpec（独立身份/系统提示/预算/工具白名单/权限/Skill/Depth）。
**Files**：`internal/agent-orchestration/`。
**Tests Required**：Spec 构造 Unit。
**Acceptance Criteria**：
- [ ] 独立身份与隔离字段完整
- [ ] 白名单/权限只能等于或收窄父
**Definition of Done**：接口评审通过。

## FC-SUB-002 — Delegator 与 runtime-core 启动
| Type | Implementation | Priority | P1 | Milestone | M10 | Status | Backlog | Size | L |
| Dependencies | FC-SUB-001 | Related Requirements | FR-SUBAGENT-002 |

**Description**：Delegate/DelegateMany；经 runtime-core 启动隔离子 Agent；并发 limit、超时、取消、重试。
**Security Considerations**：子工具调用经统一管线（不绕过）。
**Tests Required**：并发 Race、超时/取消 Failure。
**Acceptance Criteria**：
- [ ] 并发受 limit
- [ ] 子工具经权限引擎
**Definition of Done**：委派端到端跑通。

## FC-SUB-003 — 结构化结果与上下文隔离
| Type | Implementation | Priority | P0 | Milestone | M10 | Status | Backlog | Size | M |
| Dependencies | FC-SUB-002 | Related Requirements | FR-SUBAGENT-003 |

**Description**：默认返回结构化摘要（非完整日志）；工具输出隔离；不默认继承全部父上下文。
**Acceptance Criteria**：
- [ ] 返回摘要而非完整日志
- [ ] 子上下文最小授权

## FC-SUB-004 — 预算统计与递归限制
| Type | Implementation | Priority | P0 | Milestone | M10 | Status | Backlog | Size | M |
| Dependencies | FC-SUB-002, FC-TEL-004 | Related Requirements | FR-SUBAGENT-004 |

**Description**：父子预算统计；递归深度限制；超预算终止子树。
**Security Considerations**：RISK-013。
**Acceptance Criteria**：
- [ ] 递归深度 ≤ 配置
- [ ] 超预算终止子树

## FC-TEAM-001 — Team、Task DAG 与 Member Registry
| Type | Architecture | Priority | P2 | Milestone | M12 | Status | Backlog | Size | L |
| Dependencies | FC-SUB-002 | Related Requirements | FR-TEAM-001 |

**Description**：Team/TeamLead、Task DAG（Task State）、Member Registry、Shared State。
**Acceptance Criteria**：
- [ ] DAG 依赖求解正确
- [ ] Team/Task 状态与 GLOSSARY 一致

## FC-TEAM-002 — 中心化调度、Mailbox、结果集成
| Type | Implementation | Priority | P2 | Milestone | M12 | Status | Backlog | Size | L |
| Dependencies | FC-TEAM-001 | Related Requirements | FR-TEAM-002 |

**Description**：Lead 分配 Ready 任务、定向/广播消息、Artifact 共享、结果集成、冲突处理（配合 git-worktree）。
**Implementation Notes**：中心化调度（ADR-0009），不去中心化。
**Acceptance Criteria**：
- [ ] Ready 任务被分配执行
- [ ] 结果集成与冲突处理

## FC-TEAM-003 — Team 预算、超时、重试与审计
| Type | Implementation | Priority | P2 | Milestone | M12 | Status | Backlog | Size | M |
| Dependencies | FC-TEAM-002, FC-TEL-004 | Related Requirements | FR-TEAM-003 |

**Description**：Team Token/Cost 预算、并发限制、失败重试、超时、人工介入、执行审计。
**Security Considerations**：RISK-015。
**Acceptance Criteria**：
- [ ] 预算与并发限制生效
- [ ] 成员失败可重试/人工介入

## FC-ORCH-900 — orchestration 测试套件
| Type | Test | Priority | P1 | Milestone | M10 | Status | Backlog | Size | M |
| Dependencies | FC-SUB-004 | Related Requirements | FR-SUBAGENT-001..004 |

**Description**：DAG/预算 Unit、委派 Integration、并发 Race、子绕过权限 Security。
**Acceptance Criteria**：
- [ ] 子绕过权限尝试被拒
- [ ] `go test -race` 通过
