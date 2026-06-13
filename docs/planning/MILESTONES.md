# MILESTONES

按 **依赖关系** 而非纯时间排序。每个里程碑给出目标、入口条件、退出条件（Definition of Done）。周数为参考，依赖未满足时顺延。

| 里程碑 | 周 | 版本 | 主题 | 退出条件（DoD） |
| --- | --- | --- | --- | --- |
| M1 | W1–2 | MVP | 单 Agent 可运行 | Provider 抽象+Streaming+Agent Loop+Tool Registry+ReadFile/Glob/Grep/Bash+Session+CLI 跑通只读任务，事件落库可恢复 |
| M2 | W3–4 | MVP | 编辑与安全 | WriteFile/EditFile/Diff/Checkpoint + Permission Engine(L1–L3+L5) + Approval + Hook Event Bus，危险 Bash 走审批 |
| M3 | W5 | MVP | 上下文与预算 | Token 估算 + Tool Result 截断 + Observation 压缩 + Auto Compaction + Cost Budget + Loop Detection，压缩可回滚 |
| M4 | W5 | MVP | 恢复闭环 | Session 暂停/恢复、用户取消、崩溃恢复 Recovery Test 通过（与 M3 并行收尾） |
| M5 | — | MVP | MVP 演示 | 6 条核心流程(§6.7 1–5,12)可复现 Demo；MVP 验收清单全过 |
| M6 | W6 | V0.2 | 扩展系统 | Skill 包加载 + Slash Command 完整 + Hook 完善 + 三个标杆命令注册 |
| M7 | W7 | V0.2 | MCP | MCP Lifecycle + stdio + Streamable HTTP + tools/resources/prompts + MCP 权限与审计 |
| M8 | W8 | V0.2 | 记忆与恢复 | User/Project Memory + FTS5 + 候选记忆审批 + Session Resume 增强 |
| M9 | W6–8 | V0.2 | 沙箱与评测 | Docker Sandbox(可选) + Replay + 三场景 Eval Case（与 M6–M8 并行） |
| M10 | W9 | V0.3 | SubAgent | Agent Definition + 独立上下文 + 委派 + 并发 + 结构化结果 + 预算/递归限制 |
| M11 | W10 | V0.3 | Worktree | Worktree 生命周期 + 分支 + Diff + Commit + 冲突 + 清理 + 孤儿回收 |
| M12 | W11 | V1.0 | Agent Teams | Team Lead + Task DAG + Mailbox + Artifact Store + Team Budget + 结果集成 |
| M13 | W12 | V1.0 | 展示与收尾 | Replay/Eval/Security Tests/Metrics/Demo/README/架构图/技术文章 |

## 关键入口依赖
- M2 依赖 M1（工具管线、事件）。
- M3 依赖 M1（Provider、上下文基础）。
- M6/M7 依赖 M2（统一工具/权限管线）。
- M9 Sandbox 依赖 M2（Permission L4 挂钩）。
- M10 依赖 M1–M3（Runtime 可复用启动子 Agent）。
- M11 依赖 M2（工具/事件）；M10 与 M11 可部分并行。
- M12 依赖 M10+M11（SubAgent + Worktree）。
