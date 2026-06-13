# context-manager Tasks

模块：`context-manager`。Task 前缀 `FC-CTX`。相关需求 FR-CONTEXT-001..006。RISK-009。

## FC-CTX-001 — 分层上下文模型与组装
| Type | Architecture | Priority | P0 | Milestone | M3 | Status | Ready | Size | L |
| Dependencies | FC-PROV-007 | Related Requirements | FR-CONTEXT-001 | Spec | §6 |

**Description**：定义 Layer 枚举与 Builder，按层有序组装请求消息。
**Files**：`internal/context-manager/builder.go`。
**Tests Required**：分层组装 Unit。
**Acceptance Criteria**：
- [ ] 各层顺序与窗口约束生效
- [ ] 缺层降级合理
**Definition of Done**：与 runtime-core 集成组装请求。

## FC-CTX-002 — Token 估算与预算
| Type | Implementation | Priority | P0 | Milestone | M3 | Status | Backlog | Size | M |
| Dependencies | FC-CTX-001 | Related Requirements | FR-CONTEXT-002, NFR-COST-001 |

**Description**：TokenEstimator（启发式）、Reserved Output、Token/Cost Budget 计算与超限信号。
**Implementation Notes**：OPEN_QUESTIONS Q4。
**Acceptance Criteria**：
- [ ] 估算误差在可接受范围
- [ ] 超预算向 runtime-core 暴露信号

## FC-CTX-003 — 工具输出截断与压缩
| Type | Implementation | Priority | P0 | Milestone | M3 | Status | Backlog | Size | L |
| Dependencies | FC-CTX-001 | Related Requirements | FR-CONTEXT-003, NFR-LIMIT-001 |

**Description**：Bash 头尾保留、Grep 去重、ReadFile 分页、JSON 提取、Observation 压缩。
**Tests Required**：Golden Test 各规则。
**Acceptance Criteria**：
- [ ] 各截断规则按配置生效
- [ ] 截断标注保留

## FC-CTX-004 — 自动与手动 Compaction
| Type | Implementation | Priority | P0 | Milestone | M3 | Status | Backlog | Size | L |
| Dependencies | FC-CTX-002, FC-CTX-003, FC-SESS-003 | Related Requirements | FR-CONTEXT-004 |

**Description**：超阈值自动压缩 + `/compact`；压缩前 Checkpoint，生成 CompactedHistory。
**Acceptance Criteria**：
- [ ] 超阈值触发
- [ ] 压缩前必有 Checkpoint
**Definition of Done**：与 runtime-core Compacting 联动。

## FC-CTX-005 — 关键事实保护与压缩校验
| Type | Implementation | Priority | P0 | Milestone | M3 | Status | Backlog | Size | M |
| Dependencies | FC-CTX-004 | Related Requirements | FR-CONTEXT-005 |

**Description**：KeyFacts 提取与 Verify；缺失关键事实判定失败。
**Security Considerations**：KeyFacts 脱敏。
**Tests Required**：Golden（压缩前后 KeyFacts 一致）。
**Acceptance Criteria**：
- [ ] 目标/计划/文件行号/未完成任务/权限决定保留
- [ ] 缺失触发失败路径

## FC-CTX-006 — 压缩失败检测与回滚
| Type | Implementation | Priority | P1 | Milestone | M3 | Status | Backlog | Size | M |
| Dependencies | FC-CTX-005 | Related Requirements | FR-CONTEXT-006 |

**Description**：Verify 失败回滚 PreCompact Checkpoint，跳过压缩并告警（RISK-009）。
**Tests Required**：Failure Injection。
**Acceptance Criteria**：
- [ ] 校验失败回滚成功
- [ ] 回滚后可继续运行

## FC-CTX-007 — context-manager 测试套件
| Type | Test | Priority | P0 | Milestone | M3 | Status | Backlog | Size | M |
| Dependencies | FC-CTX-006 | Related Requirements | FR-CONTEXT-001..006 |

**Description**：Unit/Golden/Failure/Integration 汇总。
**Acceptance Criteria**：
- [ ] Golden 覆盖截断与 KeyFacts
- [ ] `go test -race` 通过
