# MODULE_MAP

本文件是 ForgeCode 的 **权威模块清单**。所有 Spec/Task/Checklist 必须使用此处的 Module ID 与命名。模块划分依据：单一职责、高内聚、低耦合、清晰依赖方向、可独立测试、安全边界、并发边界、数据所有权、生命周期、扩展性。

## 划分说明与反模式规避

- 候选清单有 17 项，本项目采用 **17 个模块**，但通过明确"核心抽象拥有者 vs 调用方"避免碎片化。
- `extension-system` 合并 Skill / Slash Command / Hook 三个内部组件（共享扩展加载、来源追踪、权限声明机制），避免三个单文件模块。
- `agent-orchestration` 合并 Task / SubAgent / Agent Team / Mailbox / Artifact Store（共享调度、预算、取消语义），但 Spec 中明确区分 SubAgent（一次性委派）与 Team（长期角色+DAG）。
- `event-system`（契约 + 进程内 Bus）与 `session-store`（持久化）拆分：前者拥有事件**格式与分发**，后者拥有事件**durable 存储与恢复**。避免"多个模块各自维护事件格式"。
- `sandbox` 独立于 `permission-engine`：权限引擎只做**决策**，不执行；沙箱做**受控执行**，避免"Permission Engine 直接执行命令"。
- `git-worktree` 独立于 `builtin-tools`：避免"Worktree 逻辑混入 Git Tool"。
- `telemetry` 作为底层被广泛依赖，但不反向依赖业务模块，避免成为 `utils` 垃圾桶。
- 无 `common`/`utils`/`manager` 垃圾桶模块。

## 模块清单

| Module ID | 名称 | 职责 | 拥有的数据 | 公开接口（核心） | 依赖 | 调用方 | 风险等级 | MVP |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `runtime-core` | Runtime Core | Agent Loop、状态机、Runtime Coordinator、取消、恢复编排 | AgentInstance 运行态 | `Runtime`, `AgentLoop`, `Coordinator` | model-provider, tool-runtime, event-system, context-manager, session-store, permission-engine, telemetry | cli, agent-orchestration | High | Yes |
| `model-provider` | Model Provider | Provider 抽象与适配器、消息转换、错误归一化 | — | `Provider`, `ChatRequest/Response`, `StreamChunk` | telemetry | runtime-core, context-manager | High | Yes |
| `tool-runtime` | Tool Runtime | 统一 Tool 接口、Registry、调用管线、截断/超时/审计编排 | ToolCall/ToolResult 契约 | `Tool`, `ToolDescriptor`, `Registry`, `Invoker` | permission-engine, event-system, telemetry | runtime-core, builtin-tools, mcp-client | High | Yes |
| `builtin-tools` | Built-in Tools | ReadFile/WriteFile/EditFile/Bash/Glob/Grep | — | 各 Tool 实现 | tool-runtime, permission-engine, session-store(checkpoint) | runtime-core(经 registry) | High | Yes |
| `permission-engine` | Permission Engine | 五层纵深权限决策（不执行）、Bash 结构化分析、决策优先级 | Approval(契约) | `Decider`, `Decision`, `RiskLevel`, `BashAnalyzer` | event-system, telemetry | tool-runtime, mcp-client, extension-system, sandbox | Critical | Yes |
| `sandbox` | Sandbox | Docker 受控执行、资源限制、降级 | — | `Sandbox`, `ExecSpec`, `ExecResult` | permission-engine, telemetry | tool-runtime(Bash) | High | No (V0.2) |
| `event-system` | Event System | 统一 Event Envelope、进程内 Event Bus、事件分类 | Event(格式契约) | `Event`, `Bus`, `Subscriber`, `EventType` | telemetry | 全体 | High | Yes |
| `session-store` | Session Store | Append-only Event Store、Session/Checkpoint 持久化与恢复 | Session, Event(存储), Message, Checkpoint | `Store`, `EventLog`, `Checkpointer` | event-system | runtime-core, context-manager, cli, evaluation | High | Yes |
| `context-manager` | Context Manager | 分层上下文、Token/Cost 预算、截断、Compaction | 上下文组装态 | `ContextBuilder`, `Compactor`, `TokenEstimator`, `Budget` | model-provider, session-store, telemetry | runtime-core | High | Yes |
| `memory-system` | Memory System | 四类记忆、FTS5 检索、候选审批、污染控制 | Memory | `MemoryStore`, `Recall`, `CandidateReview` | session-store, event-system | runtime-core, context-manager | Medium | No (V0.2) |
| `extension-system` | Extension System | Skill / Slash Command / Hook（统一扩展加载与来源追踪） | SkillMetadata, Hook/Command 注册 | `CommandRegistry`, `HookDispatcher`, `SkillManager` | tool-runtime, permission-engine, event-system, mcp-client | cli, runtime-core | High | Partial（Command+Hook MVP；Skill V0.2） |
| `mcp-client` | MCP Client | MCP 协议客户端、Server 生命周期、Tool/Resource/Prompt 接入 | MCP Server 连接态 | `MCPClient`, `ServerHandle`, `Transport` | tool-runtime, permission-engine, event-system | extension-system, runtime-core | High | No (V0.2) |
| `agent-orchestration` | Agent Orchestration | Task、SubAgent、Agent Team、Mailbox、Artifact Store、中心化调度 | AgentDefinition, Task, Team, Artifact | `Delegator`, `TeamLead`, `TaskGraph`, `Mailbox`, `ArtifactStore` | runtime-core, git-worktree, event-system, telemetry | cli, runtime-core | High | No（SubAgent V0.3 / Team V1.0） |
| `git-worktree` | Git Worktree | Worktree 生命周期、分支、Diff、合并、清理、孤儿回收 | Worktree 登记表 | `WorktreeManager`, `Worktree` | event-system, telemetry | agent-orchestration | Medium | No (V0.3) |
| `telemetry` | Telemetry | 结构化日志、脱敏、指标、Trace、Audit Sink、Usage/Cost | UsageRecord, 日志/指标 | `Logger`, `Metrics`, `Tracer`, `AuditSink`, `UsageMeter` | （无业务依赖） | 全体 | Medium | Yes |
| `evaluation` | Evaluation | Replay、Eval Case 框架、三场景 Eval、Skill 候选回放 | Eval 结果 | `Replayer`, `EvalRunner`, `Scorer` | runtime-core, session-store | cli, ci | Medium | No (V0.2+) |
| `cli` | CLI | 交互入口、流式输出、审批交互、命令分派（无核心业务逻辑） | — | `App`, `REPL`, `Renderer` | runtime-core, extension-system, session-store, telemetry | User | Medium | Yes |

## MVP 模块集（M1–M5）
`runtime-core`、`model-provider`、`tool-runtime`、`builtin-tools`、`permission-engine`、`event-system`、`session-store`、`context-manager`、`telemetry`、`cli`，外加 `extension-system` 的 Command+Hook 部分。
