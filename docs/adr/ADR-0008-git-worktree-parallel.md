# ADR-0008：并行写任务使用 Git Worktree

| 字段 | 值 |
| --- | --- |
| Status | Accepted |
| Date | 2026-06-13 |
| Deciders | 架构组 |
| Related Requirements | FR-WORKTREE-001, FR-WORKTREE-002, FR-SUBAGENT-002 |
| Related Modules | git-worktree, agent-orchestration |

## Context
多个写代码的 Agent 并行时若共用工作区，会互相覆盖、产生竞态与不可控冲突。

## Decision
会修改代码的并行 Agent 使用 **独立 Git Worktree + 临时分支** 隔离执行，完成后生成 Diff 提交主 Agent 审核，再 Merge/Cherry-pick/Discard。Worktree 生命周期、冲突与孤儿回收由独立 `git-worktree` 模块管理，不混入 Git 工具。

## Alternatives Considered
- **目录复制隔离**：丢失 Git 上下文、合并困难——拒绝。
- **共享工作区 + 锁**：串行化、易死锁、丧失并行收益——拒绝。

## Consequences
- 正面：真正并行写、隔离、可审查合并、保留 Git 语义。
- 负面：孤儿 Worktree（RISK-014）、合并冲突处理、主仓未提交修改边界。

## Security Impact
隔离降低并行 Agent 互相破坏风险；Worktree 操作落审计事件。

## Operational Impact
需启动时孤儿扫描与清理命令；登记表追踪 Worktree。

## Revisit Conditions
若 Git Worktree 在目标平台支持不佳，评估替代隔离方案。
