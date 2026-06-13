# git-worktree Tasks

模块：`git-worktree`。Task 前缀 `FC-WT`。相关需求 FR-WORKTREE-001..003。ADR-0008。RISK-014。

## FC-WT-000 — Spike：git CLI vs go-git Worktree 支持
| Type | Spike | Priority | P1 | Milestone | M11 | Status | Ready | Size | XS |
| Dependencies | - | Related Requirements | FR-WORKTREE-001 |

**需回答**：Worktree 操作用 git CLI 还是 go-git（OPEN_QUESTIONS Q3）。
**最小实验**：两种方式创建/列出/删除 Worktree。
**输出决策**：选定方式，更新 Q3。
**结束条件**：决策写入文档。

## FC-WT-001 — WorktreeManager 与生命周期
| Type | Architecture | Priority | P1 | Milestone | M11 | Status | Backlog | Size | L |
| Dependencies | FC-WT-000 | Related Requirements | FR-WORKTREE-001 | Spec | §6/§8 |

**Description**：WorktreeManager、临时分支、创建、登记表、Worktree State 状态机。
**Files**：`internal/git-worktree/`。
**Tests Required**：生命周期 Integration（真实 git）。
**Acceptance Criteria**：
- [ ] 创建临时分支与 Worktree
- [ ] 状态与 GLOSSARY 一致
**Definition of Done**：基础生命周期跑通。

## FC-WT-002 — Diff、Merge、Cherry-pick、Discard
| Type | Implementation | Priority | P1 | Milestone | M11 | Status | Backlog | Size | L |
| Dependencies | FC-WT-001 | Related Requirements | FR-WORKTREE-001 |

**Description**：隔离目录 Diff 生成；Merge/Cherry-pick/Discard；合并冲突回待处理。
**Acceptance Criteria**：
- [ ] Diff 正确
- [ ] 合并冲突不丢失变更

## FC-WT-003 — 边界条件与孤儿回收
| Type | Security | Priority | P0 | Milestone | M11 | Status | Backlog | Size | L |
| Dependencies | FC-WT-001 | Related Requirements | FR-WORKTREE-002, FR-WORKTREE-003 |

**Description**：主仓未提交、分支/路径冲突、无修改、同文件冲突、未跟踪文件、测试失败、清理失败、占用、取消回收、崩溃孤儿扫描清理。
**Security Considerations**：RISK-014；路径受 permission-engine 边界。
**Tests Required**：Failure Injection 覆盖各边界 + Recovery（孤儿）。
**Acceptance Criteria**：
- [ ] 各边界条件被正确处理
- [ ] 崩溃孤儿被扫描清理

## FC-WT-004 — Worktree 审计与指标
| Type | Implementation | Priority | P1 | Milestone | M11 | Status | Backlog | Size | S |
| Dependencies | FC-WT-001, FC-EVT-002 | Related Requirements | FR-WORKTREE-003 |

**Description**：WorktreeCreate/Remove 事件、活跃/孤儿/冲突指标。
**Acceptance Criteria**：
- [ ] 操作产生审计事件
- [ ] 指标可观察

## FC-WT-005 — git-worktree 测试套件
| Type | Test | Priority | P1 | Milestone | M11 | Status | Backlog | Size | M |
| Dependencies | FC-WT-002, FC-WT-003 | Related Requirements | FR-WORKTREE-001..003 |

**Description**：Integration（真实 git）、边界 Failure、孤儿 Recovery 汇总。
**Acceptance Criteria**：
- [ ] 边界与孤儿测试通过
- [ ] `go test -race` 通过
