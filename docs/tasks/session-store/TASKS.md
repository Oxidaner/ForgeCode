# session-store Tasks

模块：`session-store`。Task 前缀 `FC-SESS`。相关需求 FR-SESSION-001..004。ADR-0002/0006。RISK-016/017。

## FC-SESS-000 — Spike：SQLite 驱动与 FTS5
| Type | Spike | Priority | P0 | Milestone | M1 | Status | Ready | Size | XS |
| Dependencies | - | Related Requirements | FR-SESSION-001 |

**需回答**：驱动选型（OPEN_QUESTIONS Q2），是否免 CGO，FTS5 可用性（供 memory-system）。
**最小实验**：建表 + WAL + 简单 FTS5 查询。
**输出决策**：确定驱动，更新 Q2。
**结束条件**：决策写入文档。

## FC-SESS-001 — Append-only Event Store
| Type | Architecture | Priority | P0 | Milestone | M1 | Status | Backlog | Size | L |
| Dependencies | FC-SESS-000, FC-EVT-001 | Related Requirements | FR-SESSION-001 | Spec | §6 |

**Description**：EventLog.Append/Read，事务内分配单调 Sequence，EventID 去重，WAL。
**Implementation Notes**：单 Session 写串行化（RISK-017）。
**Tests Required**：Sequence 单调/去重 Unit、并发写 Race。
**Acceptance Criteria**：
- [ ] Sequence 严格单调连续
- [ ] 重复 EventID 不重复写
**Definition of Done**：Unit + Race 通过。

## FC-SESS-002 — Session 元数据与状态
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Backlog | Size | M |
| Dependencies | FC-SESS-001 | Related Requirements | FR-SESSION-002 |

**Description**：Store.CreateSession/Get/List/UpdateState，Session State 枚举（GLOSSARY）。
**Acceptance Criteria**：
- [ ] 状态转移合法
- [ ] 可列出未完成 Session（供 Resume）

## FC-SESS-003 — Checkpoint 与恢复读流
| Type | Implementation | Priority | P0 | Milestone | M4 | Status | Backlog | Size | L |
| Dependencies | FC-SESS-001 | Related Requirements | FR-SESSION-002, FR-SESSION-003 |

**Description**：Checkpointer 读写；按 Session/from Sequence 读流供 runtime-core 重放。
**Security Considerations**：重放不触发外部副作用。
**Tests Required**：Recovery Test。
**Acceptance Criteria**：
- [ ] Checkpoint 定位到 Sequence
- [ ] 杀进程后读流完整（NFR-REL-001）
**Definition of Done**：与 runtime-core FC-RT-006 集成。

## FC-SESS-004 — 事件归档与增长控制
| Type | Implementation | Priority | P1 | Milestone | M8 | Status | Backlog | Size | M |
| Dependencies | FC-SESS-001 | Related Requirements | FR-SESSION-004 |

**Description**：超阈值归档冷事件，保留恢复所需子集（RISK-016）。
**Acceptance Criteria**：
- [ ] 归档后仍可恢复
- [ ] DB 大小受控

## FC-SESS-005 — Schema 版本与迁移
| Type | Migration | Priority | P1 | Milestone | M1 | Status | Backlog | Size | S |
| Dependencies | FC-SESS-001 | Related Requirements | NFR-COMPAT-001 |

**Description**：Schema 版本表与启动迁移；旧事件向后兼容读取。
**Acceptance Criteria**：
- [ ] 旧 Schema 自动迁移
- [ ] 迁移失败安全中止并告警

## FC-SESS-006 — 持久化失败处理
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Backlog | Size | S |
| Dependencies | FC-SESS-001 | Related Requirements | NFR-REL-001 |

**Description**：写失败重试，持续失败暂停 Session 告警，不静默丢事件。
**Tests Required**：Failure Injection（写失败/锁）。
**Acceptance Criteria**：
- [ ] 写失败重试
- [ ] 不丢事件

## FC-SESS-007 — session-store 测试套件
| Type | Test | Priority | P0 | Milestone | M4 | Status | Backlog | Size | M |
| Dependencies | FC-SESS-003 | Related Requirements | FR-SESSION-001..004 |

**Description**：Unit/Integration/Recovery/Race/Migration 汇总。
**Acceptance Criteria**：
- [ ] `go test -race` 通过
- [ ] Recovery 与 Migration 测试通过
