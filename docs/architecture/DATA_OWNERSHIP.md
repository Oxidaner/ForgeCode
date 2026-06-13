# DATA_OWNERSHIP

每类核心实体只有 **一个主拥有模块**（唯一写入者）。其他模块通过接口读取或请求变更，**不直接写入**他人拥有的实体。这避免"Session 同时承担 Memory""多个模块写同一实体"等反模式。

| 实体 | 主拥有模块（唯一写入） | 读取方 | 持久化 | 说明 |
| --- | --- | --- | --- | --- |
| Session | `session-store` | runtime-core, cli, context-manager, evaluation | SQLite | 会话根；状态字段只由 session-store 写 |
| Event（存储） | `session-store` | evaluation, telemetry(经 Bus) | SQLite append-only | 格式契约属 event-system，存储属 session-store |
| Event（格式契约/EventType） | `event-system` | 全体（只读引用） | — | 类型与 Envelope 定义 |
| Message | `session-store` | context-manager, runtime-core | SQLite | 模型消息 |
| ToolCall（契约） | `tool-runtime` | runtime-core, session-store | — | 调用请求结构 |
| ToolResult | `tool-runtime`→落库 `session-store` | runtime-core, context-manager | SQLite | 结果（截断后） |
| Approval | `permission-engine`（决策）→落库 `session-store` | runtime-core, cli, telemetry | SQLite | 审批记录 |
| Checkpoint | `session-store` | runtime-core, context-manager, cli | SQLite | 回退点 |
| Memory | `memory-system` | runtime-core, context-manager | SQLite+FTS5 | 四类记忆，候选经审批入库 |
| SkillMetadata | `extension-system` | runtime-core, cli | SQLite+FS | Skill manifest 索引 |
| Command/Hook 注册 | `extension-system` | cli, runtime-core | 内存+配置 | 命令与 Hook 注册表 |
| AgentDefinition | `agent-orchestration` | runtime-core | SQLite/配置 | SubAgent 与角色定义 |
| AgentInstance（运行态） | `runtime-core` | agent-orchestration(只读状态) | 事件 | 运行中的 Agent 状态 |
| Task | `agent-orchestration` | runtime-core, cli | SQLite | Team Task 节点 |
| Team | `agent-orchestration` | cli | SQLite | 团队 |
| Worktree（登记） | `git-worktree` | agent-orchestration | SQLite/FS | Worktree 登记表 |
| Artifact | `agent-orchestration`（ArtifactStore） | runtime-core, cli | FS+SQLite 引用 | 产物 |
| MCP Server 连接态 | `mcp-client` | tool-runtime | 内存 | 连接/健康状态 |
| UsageRecord | `telemetry` | cli, agent-orchestration(预算) | SQLite | Token/Cost 记账 |
| 日志/指标/Trace | `telemetry` | 外部 | 文件/导出 | 脱敏后 |

## 关键约束

- **AgentInstance 运行态** 由 `runtime-core` 拥有；`agent-orchestration` 只读其状态用于调度，不直接修改。子 Agent 的运行态同样属其各自 Runtime。
- **Approval**：决策由 `permission-engine` 产生，但持久化交给 `session-store`（统一事件/记录存储）。这是"决策与存储分离"。
- **Event**：`event-system` 拥有**格式**，`session-store` 拥有**存储**。两者不冲突——一个定义 schema，一个负责落盘与读取。
- **ToolResult / Approval** 的落库统一走 `session-store`，避免 tool-runtime / permission-engine 各自接 DB。
- 跨模块变更通过 **发布事件** 或 **调用拥有者接口** 完成，禁止直接写他人表。
