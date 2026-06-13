# ROADMAP

12 周路线图，按 **依赖关系** 组织。每阶段标注目标、负责模块、关键任务与可并行项。里程碑映射见 `MILESTONES.md`，任务详情见 `tasks/MASTER_TASKS.md`。

## Week 1–2：单 Agent 可运行（M1）
- 目标：Provider 抽象、Streaming、Agent Loop、Tool Registry、ReadFile、Glob、Grep、Bash、Session、CLI。
- 模块：model-provider, runtime-core, tool-runtime, builtin-tools(只读+Bash), event-system, session-store, telemetry, cli。
- 关键路径：FC-EVT-001 → FC-SESS-001 → FC-PROV-001 → FC-TOOL-001 → FC-RT-001 → FC-CLI-001。
- 可并行：builtin-tools 只读工具、telemetry 日志、Provider 适配器。

## Week 3–4：编辑与安全（M2）
- 目标：WriteFile、EditFile、Diff、Checkpoint、Permission Engine、Approval、Hook Event Bus、Docker Sandbox（挂钩，实现延后到 M9）。
- 模块：builtin-tools(写), permission-engine, extension-system(Hook+Command), session-store(Checkpoint)。
- 关键路径：FC-PERM-001 → FC-PERM-005 → FC-TOOL-002(管线集成) → FC-PERM-007(审批) → FC-CTX/Checkpoint。
- 可并行：Hook Event Bus、Slash Command 框架、EditFile/WriteFile。

## Week 5：上下文与预算（M3/M4）
- 目标：Token 估算、Tool Result 截断、Observation 压缩、Auto Compaction、Cost Budget、Loop Detection、暂停/恢复/取消。
- 模块：context-manager, runtime-core(恢复)。
- 关键路径：FC-CTX-001 → FC-CTX-002 → FC-CTX-003 → FC-CTX-004 → FC-RT-006(恢复)。

## Week 6：扩展系统（M6）
- 目标：Skill 包、Slash Command 完善、Hook 完善、三个 Review Skills 注册骨架。
- 模块：extension-system, evaluation(场景定义)。
- 可并行：三个 Review Skill 内容（review-pr/sql/k8s）。

## Week 7：MCP（M7）
- 目标：MCP Lifecycle、stdio、Streamable HTTP、Tools、Resources、Prompts、MCP Permission。
- 模块：mcp-client, permission-engine(MCP 权限等级)。
- 关键路径：FC-MCP-001 → FC-MCP-002 → FC-MCP-003 → FC-MCP-004。

## Week 8：记忆与恢复（M8）
- 目标：User/Project Memory、Session Resume、FTS、Candidate Memory。
- 模块：memory-system。
- 可并行：FTS 检索、候选审批流。

## Week 9：SubAgent（M10）
- 目标：Agent Definition、独立上下文、委派、并发、结构化结果。
- 模块：agent-orchestration(SubAgent)。
- 关键路径：FC-SUB-001 → FC-SUB-002 → FC-SUB-003 → FC-SUB-004。

## Week 10：Worktree（M11）
- 目标：Worktree 生命周期、分支、Diff、Commit、冲突、清理。
- 模块：git-worktree。
- 与 Week 9 部分并行（不同模块、不同文件）。

## Week 11：Agent Teams（M12）
- 目标：Team Lead、Task DAG、Mailbox、Artifact Store、Team Budget、结果集成。
- 模块：agent-orchestration(Team)。
- 依赖 SubAgent + Worktree。

## Week 12：Eval 和项目展示（M13）
- 目标：Replay、Eval Cases、Security Tests、Metrics、Demo、README、架构图、技术文章。
- 模块：evaluation, telemetry, 文档。

## 调整说明
- Sandbox 完整实现从 W3–4 推迟到 W6–8（M9），仅在 M2 预留 L4 挂钩接口——因 MVP 演示不强依赖容器隔离（本地受限执行足以闭环），优先保证单 Agent 安全闭环。
- Memory、MCP、Skill 均为 V0.2，不进入 MVP 关键路径，以控制 RISK-001。
