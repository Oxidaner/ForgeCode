# TWO_PERSON_WORK_SPLIT

本文档定义 ForgeCode 在两名开发者协作时的默认分工、接口所有权、阶段计划与冲突规避规则。它是 `ROADMAP.md`、`MILESTONES.md` 与 `MASTER_TASKS.md` 的执行补充，不替代权威模块、任务或需求定义。

## 文档元数据

| 字段 | 内容 |
| --- | --- |
| Status | Draft |
| Last Updated | 2026-06-14 |
| Scope | 两人协作分工 |
| Related Docs | `docs/architecture/MODULE_MAP.md`, `docs/tasks/MASTER_TASKS.md`, `docs/planning/ROADMAP.md`, `docs/planning/MILESTONES.md` |

## 分工原则

- 按控制平面与执行安全平面拆分，避免两个人同时修改同一组核心接口。
- 模块 Owner 负责接口设计、实现推进和最终一致性；Reviewer 负责审查风险、兼容性和安全边界。
- 核心契约先冻结再并行开发，尤其是 Event、Tool、Provider、Permission、Runtime State。
- 安全相关改动必须双人 Review，包括 Bash、路径权限、审批、MCP Tool、Hook、Sandbox、Worktree。
- CLI 只作为入口和交互层，不承载核心业务逻辑。

## 人员分工

| 人 | 角色 | 主责 | 模块 |
| --- | --- | --- | --- |
| A | Runtime Owner | Agent 主循环、状态机、恢复、上下文、Provider、持久化、长期调度 | `runtime-core`, `model-provider`, `event-system`, `session-store`, `context-manager`, `memory-system`, `agent-orchestration` |
| B | Execution Owner | 工具、权限、安全、CLI、扩展、外部集成、执行隔离、评测 | `tool-runtime`, `builtin-tools`, `permission-engine`, `telemetry`, `cli`, `extension-system`, `mcp-client`, `sandbox`, `git-worktree`, `evaluation` |

## 接口所有权

| 契约 | Owner | Reviewer | 冻结条件 |
| --- | --- | --- | --- |
| Event Envelope / Event Type | A | B | `FC-EVT-001` 完成并通过事件命名评审 |
| Agent State / Runtime Coordinator | A | B | `FC-RT-001`、`FC-RT-002` 完成 |
| Session / Event Store / Recovery | A | B | `FC-SESS-001`、`FC-SESS-003` 完成 |
| Provider Interface | A | B | `FC-PROV-001` 完成，Mock Provider 可驱动 Runtime |
| Context / Token / Compaction | A | B | `FC-CTX-001` 至 `FC-CTX-004` 完成 |
| Tool Descriptor / Registry / Pipeline | B | A | `FC-TOOL-001`、`FC-TOOL-002` 完成 |
| Permission Decision / Approval | B | A | `FC-PERM-001`、`FC-PERM-007` 完成 |
| Built-in Tool Behavior | B | A | `FC-BT-001` 至 `FC-BT-005` 完成 |
| CLI Command / Rendering / Approval UI | B | A | `FC-CLI-001` 至 `FC-CLI-004` 完成 |
| Hook / Skill / Slash Command Contract | B | A | `FC-HOOK-001`、`FC-CMD-001`、`FC-SKILL-001` 完成 |
| MCP Tool Adapter Contract | B | A | `FC-MCP-001` 至 `FC-MCP-004` 完成 |
| SubAgent / Team Task Graph | A | B | `FC-SUB-001`、`FC-TEAM-001` 完成 |
| Worktree Lifecycle | B | A | `FC-WT-000` 至 `FC-WT-002` 完成 |

## MVP 阶段分工

| 阶段 | A：Runtime Owner | B：Execution Owner | 汇合点 |
| --- | --- | --- | --- |
| Week 1 | `event-system`、`model-provider` 接口、Mock Provider、Agent 状态机 | `telemetry`、`tool-runtime` 接口、`permission-engine` 最小接口 | 冻结 Event / Tool / Provider / Permission 四个核心契约 |
| Week 2 | `session-store`、Runtime Coordinator、Agent Loop | ReadFile / Glob / Grep / Bash、基础 CLI | 单 Agent 只读任务跑通 |
| Week 3-4 | Checkpoint / Recovery 编排、Runtime 暂停取消 | WriteFile / EditFile / Diff、权限 L1-L3、Approval、Hook 基础 | 编辑任务与危险操作审批闭环 |
| Week 5 | `context-manager`、Token 预算、Compaction、Loop Detection | CLI 展示成本/上下文状态、工具输出截断集成 | MVP 恢复与上下文闭环 |

MVP 关键任务链：

```text
FC-EVT-001
  -> FC-SESS-001
  -> FC-PROV-001
  -> FC-TOOL-001
  -> FC-TOOL-002
  -> FC-PERM-001
  -> FC-RT-001
  -> FC-RT-002
  -> FC-RT-003
  -> FC-CLI-001
```

编辑与审批关键任务链：

```text
FC-PERM-001
  -> FC-PERM-005
  -> FC-PERM-007
  -> FC-BT-003
  -> FC-SESS-002
  -> FC-TOOL-004
```

## V0.2 到 V1.0 阶段分工

| 阶段 | A：Runtime Owner | B：Execution Owner | 目标 |
| --- | --- | --- | --- |
| Week 6 | Session Resume 增强、Memory 设计落地 | Slash Command、Skill 基础、Hook 完善 | 扩展系统可用 |
| Week 7-8 | Memory FTS、候选记忆审批 | MCP Client、Sandbox 初版、Evaluation Replay | V0.2 能力闭环 |
| Week 9-10 | SubAgent、独立上下文、调度 | Git Worktree、Diff / Commit / Cleanup | 并行 Agent 写任务闭环 |
| Week 11-12 | Agent Team、Task DAG、Team Budget | Eval、Demo、README、展示材料 | V1.0 展示收尾 |

## 共享文件与冲突规避

| 文件或目录 | 风险 | 规则 |
| --- | --- | --- |
| `docs/architecture/MODULE_MAP.md` | 模块命名漂移 | 只在双人确认后修改 |
| `docs/architecture/EVENT_MODEL.md` | 事件契约影响全局 | A 提案，B Review |
| `docs/architecture/DATA_OWNERSHIP.md` | 数据所有权冲突 | Owner 修改前先说明影响范围 |
| `docs/planning/GLOSSARY.md` | 状态/错误枚举不一致 | 新增枚举必须同步 Spec 与 Traceability |
| `docs/tasks/MASTER_TASKS.md` | Task DAG 冲突 | 修改关键路径需双人确认 |
| `runtime-core` 相关文件 | 状态机和恢复语义漂移 | A Owner，B 不直接改核心状态机 |
| `tool-runtime` 相关文件 | Tool 管线被绕过 | B Owner，A 不直接改核心管线 |
| `permission-engine` 相关文件 | 审批和权限绕过 | B Owner，A Review 安全完整性 |
| `event-system` 相关文件 | 多套事件格式 | A Owner，禁止模块私有事件格式 |

## 日常协作节奏

- 每周开始时从 `docs/tasks/MASTER_TASKS.md` 选择当周任务，确认 Owner、Reviewer 与阻塞关系。
- 每天同步一次阻塞点，只讨论接口变更、依赖阻塞、验收失败和安全风险。
- 接口变更需要先更新对应 Spec，再改实现。
- 完成任务时同步更新模块 `TASKS.md`、`CHECKLIST.md` 的 Evidence，并检查 `TRACEABILITY_MATRIX.md` 是否需要同步。
- 每个 Milestone 结束时执行一次一致性检查：Requirement、Spec、Task、Checklist、ADR、Event、Glossary 是否一致。

## 第一批任务建议

| 人 | 首批任务 | 输出 |
| --- | --- | --- |
| A | `FC-EVT-001`, `FC-PROV-001`, `FC-RT-001`, `FC-SESS-000` | Event 契约、Provider 契约、Agent 状态机、SQLite 决策 |
| B | `FC-TEL-001`, `FC-TOOL-001`, `FC-PERM-001`, `FC-BT-001` | 日志脱敏、Tool 契约、权限决策接口、内置工具边界 |

第一批任务完成后，双方共同 Review 四个冻结点：

```text
Event Envelope
Provider Interface
Tool Descriptor / Pipeline
Permission Decision
```

冻结后再进入并行实现，避免在 M1 后半段反复返工。
