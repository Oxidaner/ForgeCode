# event-system Tasks

模块：`event-system`。Task 前缀 `FC-EVT`。相关需求 FR-EVENT-001..003。ADR-0011。

## FC-EVT-001 — Event Envelope 与 EventType 定义
| Type | Architecture | Priority | P0 | Milestone | M1 | Status | Done | Size | M |
| Dependencies | - | Related Requirements | FR-EVENT-001, FR-EVENT-003 | Spec | §6 |

**Description**：实现 EVENT_MODEL.md 的 Event 结构、EventType 权威枚举、EventClass 分类。
**Implementation Notes**：Sequence 字段存在但不在此分配。SchemaVersion 支持演进。
**Files**：`internal/event-system/types.go`。
**Tests Required**：Contract Test（字段与 EVENT_MODEL 一致）。
**Acceptance Criteria**：
- [x] 字段与文档完全一致
- [x] 全部 EventType 与分类覆盖
**Definition of Done**：契约测试通过。
**Evidence**：
- 2026-06-22：`go test ./internal/event-system` 通过。
- 2026-06-22：`internal/event-system/event_test.go` 锁定 Event Envelope 字段、EventType 顺序与 EventClass bitmask 分类。

## FC-EVT-002 — 进程内 Event Bus
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Backlog | Size | L |
| Dependencies | FC-EVT-001 | Related Requirements | FR-EVENT-002 |

**Description**：Publish/Subscribe、按订阅者有序投递、订阅者间错误隔离、过滤。
**Tests Required**：Race Test、错误隔离 Failure Injection。
**Acceptance Criteria**：
- [ ] 单订阅者有序
- [ ] 订阅者失败被隔离
- [ ] `go test -race` 通过
**Definition of Done**：并发与隔离测试通过。

## FC-EVT-003 — 订阅过滤与分类路由
| Type | Implementation | Priority | P0 | Milestone | M1 | Status | Backlog | Size | S |
| Dependencies | FC-EVT-002 | Related Requirements | FR-EVENT-003 |

**Description**：按 EventType/EventClass 过滤；为 session-store/telemetry/extension-system 提供订阅入口。
**Acceptance Criteria**：
- [ ] Durable/Audit/Hook 过滤正确
- [ ] 多订阅者各取所需

## FC-EVT-004 — 订阅者超时与背压
| Type | Implementation | Priority | P1 | Milestone | M3 | Status | Backlog | Size | S |
| Dependencies | FC-EVT-002 | Related Requirements | NFR-LIMIT-001 |

**Description**：订阅者处理超时、Ephemeral 背压/丢弃策略与指标。
**Acceptance Criteria**：
- [ ] 超时订阅者被隔离并计数
- [ ] Ephemeral 背压不阻塞 Durable

## FC-EVT-005 — Event Bus 测试套件
| Type | Test | Priority | P0 | Milestone | M1 | Status | Backlog | Size | M |
| Dependencies | FC-EVT-003 | Related Requirements | FR-EVENT-001..003 |

**Description**：Unit + Race + Contract + Failure 汇总。
**Acceptance Criteria**：
- [ ] 契约测试锁定 Envelope
- [ ] race 与隔离测试通过
