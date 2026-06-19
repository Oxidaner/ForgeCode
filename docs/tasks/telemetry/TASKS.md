# telemetry Tasks

模块：`telemetry`。Task 前缀 `FC-TEL`。相关需求 FR-TELEMETRY-001..004。

## FC-TEL-001 — Logger 与脱敏 Redactor
| Type | Architecture | Priority | P0 | Milestone | M1 | Status | Done | Size | M |
| Dependencies | - | Related Requirements | FR-TELEMETRY-001, NFR-SEC-002 | Spec | §6 |

**Description**：结构化 Logger + Redactor，脱敏密钥/Token/环境变量，结构化字段与自由文本均生效。
**Files**：`internal/telemetry/log.go`, `redact.go`。
**Security Considerations**：NFR-SEC-002。
**Tests Required**：Security Test（含密钥输入不泄露）。
**Acceptance Criteria**：
- [x] 脱敏字段不出现明文
- [x] 自由文本中的密钥模式被掩码
**Definition of Done**：安全测试通过。
**Evidence**：实现 `internal/telemetry` 的 `Logger`、`Redactor`、`MemoryMetrics`、`MemoryUsageMeter`；`go build ./...`、`go test ./...`、`go vet ./...`、`go test -race ./...` 通过。

## FC-TEL-002 — Metrics 与可选 Trace
| Type | Implementation | Priority | P1 | Milestone | M1 | Status | Backlog | Size | M |
| Dependencies | FC-TEL-001 | Related Requirements | FR-TELEMETRY-002 |

**Description**：Counter/Histogram（轮次/工具/Token/Cost/错误率），可选 Trace。
**Acceptance Criteria**：
- [ ] 关键指标可记录查询
- [ ] Trace 可开关

## FC-TEL-003 — AuditSink
| Type | Implementation | Priority | P0 | Milestone | M2 | Status | Backlog | Size | M |
| Dependencies | FC-TEL-001, FC-EVT-003 | Related Requirements | FR-TELEMETRY-003 |

**Description**：实现 eventsystem.Subscriber 消费审计类事件并落地（依赖反转，不 import Bus 实现）。
**Security Considerations**：append-only 防篡改。
**Tests Required**：Integration（审计事件落地）。
**Acceptance Criteria**：
- [ ] 审计事件完整落地
- [ ] telemetry 不 import 业务模块

## FC-TEL-004 — UsageMeter 记账与聚合
| Type | Implementation | Priority | P1 | Milestone | M3 | Status | Backlog | Size | M |
| Dependencies | FC-TEL-001 | Related Requirements | FR-TELEMETRY-004, NFR-COST-001 |

**Description**：UsageRecord 记录与按 Session/Agent/Team 聚合（支撑预算 RISK-013）。
**Acceptance Criteria**：
- [ ] 可按作用域聚合 Token/Cost
- [ ] 写失败不阻断主流程

## FC-TEL-005 — telemetry 测试套件
| Type | Test | Priority | P0 | Milestone | M2 | Status | Backlog | Size | S |
| Dependencies | FC-TEL-003 | Related Requirements | FR-TELEMETRY-001..004 |

**Description**：脱敏 Security、聚合 Unit、AuditSink Integration、Race 汇总。
**Acceptance Criteria**：
- [ ] 脱敏 Security 测试通过
- [ ] `go test -race` 通过
