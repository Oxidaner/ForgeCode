# evaluation Tasks

模块：`evaluation`。Task 前缀 `FC-EVAL`。相关需求 FR-EVAL-001..004。ADR-0010。RISK-018/019。

## FC-EVAL-001 — Replayer
| Type | Architecture | Priority | P1 | Milestone | M9 | Status | Ready | Size | M |
| Dependencies | FC-SESS-003 | Related Requirements | FR-EVAL-001 | Spec | §6 |

**Description**：从 session-store 事件重放 Session，不重执行外部副作用。
**Files**：`internal/evaluation/`。
**Tests Required**：Replay Integration。
**Acceptance Criteria**：
- [ ] 重放不触发副作用
- [ ] 轨迹可比对
**Definition of Done**：对真实 Session 重放通过。

## FC-EVAL-002 — Eval Case 框架与 Scorer
| Type | Implementation | Priority | P1 | Milestone | M9 | Status | Backlog | Size | L |
| Dependencies | FC-EVAL-001 | Related Requirements | FR-EVAL-002 |

**Description**：EvalCase/EvalSuite/EvalRunner/Scorer；规则断言 + 可选模型评审（Q15）。
**Acceptance Criteria**：
- [ ] 规则断言可评分
- [ ] 模型评审可选开关

## FC-EVAL-003 — 三标杆场景 Eval
| Type | Evaluation | Priority | P1 | Milestone | M9 | Status | Backlog | Size | L |
| Dependencies | FC-EVAL-002, FC-CMD-100 | Related Requirements | FR-EVAL-003 |

**Description**：review-pr/sql/k8s 场景检查项断言与得分（见 SPEC §10）。
**Security Considerations**：用例不含真实密钥。
**Acceptance Criteria**：
- [ ] 三场景检查项均有断言
- [ ] 产出得分报告（RISK-018/019）
**Definition of Done**：三场景 Eval 可复现。

## FC-EVAL-004 — 候选 Skill 回放评测
| Type | Implementation | Priority | P1 | Milestone | M10 | Status | Backlog | Size | M |
| Dependencies | FC-EVAL-002 | Related Requirements | FR-EVAL-004, FR-SKILL-004 |

**Description**：候选 Skill 在历史轨迹回放评测，产出通过/拒绝建议（ADR-0010 安全门）。
**Acceptance Criteria**：
- [ ] 回放产出建议
- [ ] 未通过不建议安装

## FC-EVAL-005 — EvalReport 导出与 CI 集成
| Type | Implementation | Priority | P2 | Milestone | M13 | Status | Backlog | Size | S |
| Dependencies | FC-EVAL-003 | Related Requirements | FR-EVAL-002 |

**Description**：EvalReport 导出，CI 可运行场景 Eval。
**Acceptance Criteria**：
- [ ] 报告可导出
- [ ] CI 可执行

## FC-EVAL-006 — evaluation 测试套件
| Type | Test | Priority | P1 | Milestone | M9 | Status | Backlog | Size | S |
| Dependencies | FC-EVAL-002 | Related Requirements | FR-EVAL-001..004 |

**Description**：Scorer Unit、Replay/场景 Integration、Golden 期望快照。
**Acceptance Criteria**：
- [ ] Golden 快照稳定
- [ ] `go test` 通过
