# Architecture Decision Records

记录 ForgeCode 的关键架构决策。每份 ADR 使用 `docs/templates/ADR_TEMPLATE.md` 格式。仅记录有实际决策价值的条目。

| ADR | 标题 | Status | 相关模块 |
| --- | --- | --- | --- |
| [ADR-0001](ADR-0001-event-driven-runtime.md) | 采用事件驱动 Agent Runtime | Accepted | runtime-core, event-system |
| [ADR-0002](ADR-0002-append-only-event-store.md) | 使用 Append-only Event Store | Accepted | session-store, event-system |
| [ADR-0003](ADR-0003-provider-runtime-decoupling.md) | Provider 与 Runtime 解耦 | Accepted | model-provider, runtime-core |
| [ADR-0004](ADR-0004-unified-tool-descriptor.md) | 内置 Tool 与 MCP Tool 统一描述 | Accepted | tool-runtime, mcp-client |
| [ADR-0005](ADR-0005-permission-engine-independent.md) | Permission Engine 独立于 Tool Executor | Accepted | permission-engine, tool-runtime |
| [ADR-0006](ADR-0006-sqlite-local-store.md) | SQLite 作为初始本地状态存储 | Accepted | session-store |
| [ADR-0007](ADR-0007-fts5-over-vector.md) | 记忆先用 FTS5 而非向量数据库 | Accepted | memory-system |
| [ADR-0008](ADR-0008-git-worktree-parallel.md) | 并行写任务使用 Git Worktree | Accepted | git-worktree, agent-orchestration |
| [ADR-0009](ADR-0009-team-centralized-dag.md) | Agent Team 使用中心化 Task DAG | Accepted | agent-orchestration |
| [ADR-0010](ADR-0010-skill-autogen-review.md) | Skill 自动生成需评测与人工审批 | Accepted | extension-system, evaluation |
| [ADR-0011](ADR-0011-hook-on-event-bus.md) | Hook 基于统一 Event Bus | Accepted | extension-system, event-system |
| [ADR-0012](ADR-0012-docker-sandbox.md) | 第一版 Sandbox 使用 Docker | Accepted | sandbox |
