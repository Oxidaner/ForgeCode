# sandbox Tasks

模块：`sandbox`。Task 前缀 `FC-SBX`。相关需求 FR-SANDBOX-001..003, FR-PERM-006。ADR-0012。

## FC-SBX-001 — Sandbox 接口与 Docker 执行
| Type | Architecture | Priority | P1 | Milestone | M9 | Status | Ready | Size | L |
| Dependencies | FC-PERM-006 | Related Requirements | FR-SANDBOX-001 | Spec | §6 |

**Description**：Sandbox/ExecSpec/ExecResult；Docker 执行：工作目录挂载、只读挂载、网络控制。
**Implementation Notes**：封装 Docker SDK/CLI，不自研运行时（ADR-0012）。
**Files**：`internal/sandbox/`。
**Tests Required**：Docker Integration（有 Docker 时）。
**Acceptance Criteria**：
- [ ] 命令在容器执行
- [ ] 挂载与网络策略生效
**Definition of Done**：基础容器执行跑通。

## FC-SBX-002 — 资源限制与环境过滤
| Type | Implementation | Priority | P1 | Milestone | M9 | Status | Backlog | Size | M |
| Dependencies | FC-SBX-001 | Related Requirements | FR-SANDBOX-002 |

**Description**：CPU/内存/PID/执行时间限制、环境变量白名单过滤、临时目录、进程回收。
**Security Considerations**：防密钥泄露到容器。
**Acceptance Criteria**：
- [ ] 资源限制生效
- [ ] env 仅白名单放行
- [ ] 无孤儿容器

## FC-SBX-003 — 降级策略
| Type | Implementation | Priority | P1 | Milestone | M9 | Status | Backlog | Size | M |
| Dependencies | FC-SBX-001, FC-PERM-006 | Related Requirements | FR-SANDBOX-003 |

**Description**：Docker 不可用时 Refuse（拒绝高风险）或 LocalRestricted，并审计告警。
**Tests Required**：Failure Injection（Docker 不可用）。
**Acceptance Criteria**：
- [ ] 不可用按策略降级
- [ ] 降级被审计

## FC-SBX-004 — 取消与超时回收
| Type | Implementation | Priority | P1 | Milestone | M9 | Status | Backlog | Size | S |
| Dependencies | FC-SBX-001 | Related Requirements | FR-SANDBOX-002 |

**Description**：超时/取消 Kill 容器并回收，结果标注 TimedOut。
**Acceptance Criteria**：
- [ ] 超时终止并回收
- [ ] 取消传播

## FC-SBX-005 — sandbox 测试套件
| Type | Test | Priority | P1 | Milestone | M9 | Status | Backlog | Size | M |
| Dependencies | FC-SBX-002, FC-SBX-003 | Related Requirements | FR-SANDBOX-001..003 |

**Description**：Docker Integration、降级 Failure、网络/挂载/env Security。
**Acceptance Criteria**：
- [ ] 网络隔离与挂载越界 Security 通过
- [ ] 降级测试通过
