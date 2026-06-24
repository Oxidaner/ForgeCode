# session-store Tasks

模块：`session-store`。Task 前缀 `FC-SESS`。相关需求 FR-SESSION-001..004。ADR-0002/0006。RISK-016/017。

## FC-SESS-000 — Spike：SQLite 驱动与 FTS5
| Type | Spike | Priority | P0 | Milestone | M1 | Status | Done | Size | XS |
| Dependencies | - | Related Requirements | FR-SESSION-001 |

**需回答**：驱动选型（OPEN_QUESTIONS Q2），是否免 CGO，FTS5 可用性（供 memory-system）。
**最小实验**：建表 + WAL + 简单 FTS5 查询。
**输出决策**：确定驱动，更新 Q2。
**结束条件**：决策写入文档。
**Evidence**：
- 2026-06-22：选用 `modernc.org/sqlite v1.34.5`，保持 `go.mod` 的 `go 1.22` 基线。
- 2026-06-22：`go test ./internal/session-store` 通过，覆盖 SQLite 建表、`PRAGMA journal_mode=WAL` 与 FTS5 `MATCH` 查询。

## FC-SESS-001 — Append-only Event Store
| Type | Architecture | Priority | P0 | Milestone | M1 | Status | Done | Size | L |
| Dependencies | FC-SESS-000, FC-EVT-001 | Related Requirements | FR-SESSION-001 | Spec | §6 |

**Description**：EventLog.Append/Read，事务内分配单调 Sequence，EventID 去重，WAL。
**Implementation Notes**：单 Session 写串行化（RISK-017）。
**Tests Required**：Sequence 单调/去重 Unit、并发写 Race。
**Acceptance Criteria**：
- [x] Sequence 严格单调连续
- [x] 重复 EventID 不重复写
**Definition of Done**：Unit + Race 通过。
**Evidence**：
- 2026-06-22：`go test ./internal/session-store` 通过。
- 2026-06-22：`internal/session-store/event_log_test.go` 覆盖单 Session 连续 Sequence、重复 EventID 幂等返回既有 Sequence、`Read(sessionID, from)` 有序读取、跨 Session 独立 Sequence、Envelope 字段往返与并发 append 连续性。

## FC-SESS-002 — Session 元数据与状态
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Done | Size | M |
| Dependencies | FC-SESS-001 | Related Requirements | FR-SESSION-002 |

**Description**：Store.CreateSession/Get/List/UpdateState，Session State 枚举（GLOSSARY）。
**Acceptance Criteria**：
- [x] 状态转移合法
- [x] 可列出未完成 Session（供 Resume）
**Evidence**：
- 2026-06-22：`go test ./internal/session-store` 通过。
- 2026-06-22：`internal/session-store/session_store_test.go` 覆盖 Create/Get/List/UpdateState、重复 SessionID 冲突、合法/非法状态转移、终态保护、未完成 Session 过滤与 NotFound。

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
