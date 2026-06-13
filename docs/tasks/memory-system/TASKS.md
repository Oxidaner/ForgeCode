# memory-system Tasks

模块：`memory-system`。Task 前缀 `FC-MEM`。相关需求 FR-MEMORY-001..004。ADR-0007。RISK-010。

## FC-MEM-001 — Memory 模型与存储
| Type | Architecture | Priority | P1 | Milestone | M8 | Status | Ready | Size | M |
| Dependencies | FC-SESS-000 | Related Requirements | FR-MEMORY-001 | Spec | §6 |

**Description**：四类 Memory、MemoryStore 接口、SQLite 表（与 session-store 分离）、项目隔离。
**Files**：`internal/memory-system/`。
**Tests Required**：CRUD/隔离 Unit。
**Acceptance Criteria**：
- [ ] 四类记忆可存储
- [ ] ProjectID 隔离生效
**Definition of Done**：与 session-store 同库不同表跑通。

## FC-MEM-002 — FTS5 关键词检索
| Type | Implementation | Priority | P1 | Milestone | M8 | Status | Backlog | Size | M |
| Dependencies | FC-MEM-001 | Related Requirements | FR-MEMORY-002 |

**Description**：FTS5 检索 + 置信度/时效过滤，供 context-manager RetrievedMemory。
**Implementation Notes**：不引入向量库（ADR-0007）。
**Acceptance Criteria**：
- [ ] 关键词检索可用
- [ ] 按置信度/过期过滤

## FC-MEM-003 — 元数据与生命周期
| Type | Implementation | Priority | P1 | Milestone | M8 | Status | Backlog | Size | S |
| Dependencies | FC-MEM-001 | Related Requirements | FR-MEMORY-002 |

**Description**：来源/创建/最后验证/置信度/过期；手动编辑删除；过期失效。
**Acceptance Criteria**：
- [ ] 元数据完整
- [ ] 过期自动失效

## FC-MEM-004 — 候选审批与污染控制
| Type | Security | Priority | P0 | Milestone | M8 | Status | Backlog | Size | M |
| Dependencies | FC-MEM-001, FC-EVT-002 | Related Requirements | FR-MEMORY-003, FR-MEMORY-004 |

**Description**：CandidateReview 流；不自动采信模型输出；仅 Approved 入库 Active。
**Security Considerations**：RISK-010。
**Tests Required**：Security Test（未审批不生效）。
**Acceptance Criteria**：
- [ ] 候选未审批不 Active
- [ ] 重复候选去重

## FC-MEM-005 — 敏感信息控制
| Type | Security | Priority | P0 | Milestone | M8 | Status | Backlog | Size | S |
| Dependencies | FC-MEM-001 | Related Requirements | FR-MEMORY-003 |

**Description**：密钥/PII 脱敏或拒绝入库。
**Acceptance Criteria**：
- [ ] 敏感内容被脱敏/拒绝
- [ ] 跨项目访问被拒

## FC-MEM-006 — memory-system 测试套件
| Type | Test | Priority | P1 | Milestone | M8 | Status | Backlog | Size | M |
| Dependencies | FC-MEM-004, FC-MEM-005 | Related Requirements | FR-MEMORY-001..004 |

**Description**：检索 Unit、污染/隔离 Security、与 context-manager Integration、Migration。
**Acceptance Criteria**：
- [ ] Security 测试通过（RISK-010）
- [ ] `go test -race` 通过
