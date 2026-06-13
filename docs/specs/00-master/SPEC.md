# ForgeCode 总体 Spec（Master SPEC）

## 6.1 文档元数据

| 字段 | 值 |
| --- | --- |
| Status | Draft（架构规划阶段） |
| Version | 0.1.0 |
| Owners | ForgeCode 核心架构组（当前单人，占位） |
| Last Updated | 2026-06-13 |
| Target Release | MVP（M1–M5）→ V0.2 → V0.3 → V1.0 |
| Related ADRs | ADR-0001 … ADR-0012 |

本文件是 ForgeCode 的总体规格说明，定义项目目标、非目标、功能/非功能需求、系统上下文、核心流程、数据模型与版本路线。模块级细节见 `docs/specs/<module-id>/SPEC.md`，模块划分理由见 `docs/architecture/MODULE_MAP.md`。

---

## 6.2 项目目标

### 要解决的问题
现有 Coding Agent 大多是对单一 Provider SDK 或第三方 Agent 框架的封装，运行时控制逻辑（循环、权限、恢复、并行）不透明、不可审计、不可移植。ForgeCode 从第一性原理自主实现一个 **Coding Agent Runtime 控制平面**，使任务执行 **安全、可恢复、可扩展、可观测、可并行**。

### 目标用户
- 需要在真实代码仓库上执行变更/审查任务的工程师；
- 希望理解并改造 Agent Runtime 内部机制的平台与基础设施团队；
- 将 ForgeCode 作为开源/面试展示项目的作者本人。

### 核心价值
- **Model-Agnostic**：Runtime 不依赖任何单一 Provider 的请求/响应结构。
- **Safe by default**：五层纵深权限防御 + 审批 + 审计 + Checkpoint。
- **Recoverable**：显式 Agent 状态机 + Append-only 事件 + Session 恢复。
- **Extensible**：统一 Tool 接口、MCP、Hook、Slash Command、Skill。
- **Parallel-capable**：SubAgent、Git Worktree 隔离、Agent Team。

### 面试与开源展示价值
项目展示的是 **自主实现的运行时与控制平面**，而非 SDK 编排：事件驱动状态机、统一工具/权限/审计管线、可恢复执行、并行隔离、安全模型。这些是区别于"调用 Agent SDK"的核心证据。

### 第三方库 vs 自主实现的边界

| 允许使用第三方库（基础设施） | 必须自主实现（核心逻辑） |
| --- | --- |
| HTTP/SSE 客户端、JSON、JSON Schema 校验 | Agent Loop 与状态机 |
| SQLite 驱动（如 `modernc.org/sqlite`）、FTS5 | 工具调用管线、权限引擎决策 |
| CLI 框架（如 `cobra`）、TUI 基础库 | 事件模型、Event Store、恢复语义 |
| 结构化日志、OpenTelemetry SDK | 上下文分层与压缩策略 |
| `go-git` 或 `git` CLI 包装（Worktree 底层操作） | Provider 抽象与消息转换、Tool Descriptor 统一 |
| Docker SDK（Sandbox 底层） | SubAgent / Team 调度、Task DAG、Artifact 协议 |

> 原则：**基础设施可借用，控制平面必须自研**。Provider 的私有 SDK 只能出现在 `model-provider` 适配器内部，不得泄漏到 `runtime-core`。

---

## 6.3 非目标（第一版）

- 不做完整 IDE 集成。
- 不做复杂 TUI（第一版为行式/分块输出 CLI）。
- 不支持任意深度多 Agent 嵌套（递归深度受限，见 FR-SUBAGENT）。
- 不做去中心化 Agent 协商或多 Agent 自由辩论（Team 为中心化调度）。
- 不默认引入向量数据库（记忆先用 SQLite FTS5）。
- 不实现自己的容器运行时（Sandbox 基于 Docker）。
- 不保证兼容所有 Provider 的私有能力（只统一公共能力子集）。
- 不允许 Agent 自动执行未经审批的高风险操作。
- 自动生成的 Skill 不允许未经评测与人工审批直接生效。

---

## 6.4 功能需求（Functional Requirements）

优先级：P0=MVP 阻断 / P1=重要但可后置 / P2=增强。来源：master-plan §2 对应小节。验收方式列指向 Spec/Test/Eval。

### Runtime（runtime-core）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-RUNTIME-001 | 以显式状态机驱动 Agent Loop（非不可恢复无限循环），状态见 §状态枚举 | P0 | 2.2 | runtime-core | Unit + 状态转移 Golden Test |
| FR-RUNTIME-002 | 支持最大轮次、最大工具调用次数、Token/Cost 预算、Deadline 终止 | P0 | 2.2 | runtime-core | Unit（各上限触发）|
| FR-RUNTIME-003 | 支持 Context Cancellation、用户暂停/取消，并落事件 | P0 | 2.2 | runtime-core | Integration（取消传播）|
| FR-RUNTIME-004 | 解析 Tool Call、非法 Tool Call 处理、重复工具调用与同错误循环检测 | P0 | 2.2 | runtime-core | Unit + Failure Injection |
| FR-RUNTIME-005 | Provider 错误重试与退避，区分可重试/不可重试 | P1 | 2.2 | runtime-core | Failure Injection |
| FR-RUNTIME-006 | 中间状态持久化与 Session 恢复（崩溃后从事件重建） | P0 | 2.2 | runtime-core, session-store | Recovery Test |

### Model Provider（model-provider）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-PROVIDER-001 | 统一 Provider 接口，支持普通响应与 Streaming | P0 | 2.3 | model-provider | Contract Test |
| FR-PROVIDER-002 | Tool Calling 与多 Tool Call、Stop Reason、Token Usage 解析 | P0 | 2.3 | model-provider | Contract Test |
| FR-PROVIDER-003 | 提供 OpenAI、Anthropic、OpenAI-Compatible、Mock 适配器 | P0(Mock,OpenAI) / P1(其余) | 2.3 | model-provider | Contract + Mock |
| FR-PROVIDER-004 | Rate Limit / Retry / Timeout / Provider Error 归一化 | P0 | 2.3 | model-provider | Failure Injection |
| FR-PROVIDER-005 | Model Capability 与 Context Window 元数据、Structured Output | P1 | 2.3 | model-provider | Unit |
| FR-PROVIDER-006 | Provider-specific 消息转换，不向 Runtime 泄漏私有结构 | P0 | 2.3 | model-provider | Contract Test |

### Tool Runtime（tool-runtime）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-TOOL-001 | 统一 `Tool` 接口与 `ToolDescriptor`，Tool Registry 注册/发现 | P0 | 2.1 | tool-runtime | Unit + Contract |
| FR-TOOL-002 | 统一调用管线：Validation→Permission→PreHook→Execute→PostHook→Audit | P0 | 2.1 | tool-runtime | Integration |
| FR-TOOL-003 | 输出截断、超时控制、错误分类、调用审计 | P0 | 2.1 | tool-runtime | Unit |
| FR-TOOL-004 | 内置工具与 MCP 工具共用同一 Descriptor 与调用流程 | P0 | 2.1/2.4 | tool-runtime, mcp-client | Contract Test |

### Built-in Tools（builtin-tools）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-TOOL-101 | ReadFile：分页读取、二进制识别、超大文件保护 | P0 | 2.1 | builtin-tools | Unit + Golden |
| FR-TOOL-102 | WriteFile：创建/覆盖，写前 Checkpoint，失败保持原文件不变 | P0 | 2.1 | builtin-tools | Unit + Failure |
| FR-TOOL-103 | EditFile：精确局部替换，产生 Diff，唯一匹配校验 | P0 | 2.1 | builtin-tools | Unit + Golden |
| FR-TOOL-104 | Bash：命令执行、超时、输出头尾保留、退出码与错误分类 | P0 | 2.1 | builtin-tools | Unit + Failure |
| FR-TOOL-105 | Glob：文件模式匹配 | P0 | 2.1 | builtin-tools | Unit |
| FR-TOOL-106 | Grep：关键字与正则搜索，结果去重 | P0 | 2.1 | builtin-tools | Unit |

### Permission Engine（permission-engine）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-PERM-001 | 独立 Permission Engine，输入 Descriptor+参数+上下文，输出决策 | P0 | 2.8 | permission-engine | Unit + Contract |
| FR-PERM-002 | 第一层 Schema/输入校验（必填、长度、路径格式、空字节、注入风险） | P0 | 2.8 | permission-engine | Security Test |
| FR-PERM-003 | 第二层资源边界（Workspace Root、读写目录、敏感/密钥文件、路径穿越、符号链接逃逸） | P0 | 2.8 | permission-engine | Security Test |
| FR-PERM-004 | 第三层操作风险策略：风险等级 Low/Medium/High/Critical，决策 Allow/AskOnce/AskAlways/Deny | P0 | 2.8 | permission-engine | Unit |
| FR-PERM-005 | Bash 命令结构化分析（可执行程序、参数、管道、重定向、子 Shell、命令替换、危险操作识别），非整串匹配 | P0 | 2.8 | permission-engine | Security Test |
| FR-PERM-006 | 第四层运行时沙箱挂钩（委托 sandbox 模块） | P1 | 2.8 | permission-engine, sandbox | Integration |
| FR-PERM-007 | 第五层人工审批与审计记录（含规则命中原因、原始参数、结果） | P0 | 2.8 | permission-engine, telemetry | Audit Test |
| FR-PERM-008 | 明确权限优先级与冲突决策规则（Deny 优先、最严格生效） | P0 | 2.8 | permission-engine | Unit |

### Sandbox（sandbox）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-SANDBOX-001 | Docker 沙箱执行 Bash：工作目录挂载、只读挂载、网络控制 | V0.2/P1 | 2.8 | sandbox | Integration |
| FR-SANDBOX-002 | 资源限制：CPU/内存/PID/执行时间，环境变量过滤，进程回收 | V0.2/P1 | 2.8 | sandbox | Integration |
| FR-SANDBOX-003 | 沙箱不可用时的降级策略（拒绝高风险或回退本地受限执行） | V0.2/P1 | 2.8 | sandbox, permission-engine | Failure Injection |

### Event System（event-system）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-EVENT-001 | 统一 Event Envelope（见 `EVENT_MODEL.md`），全局唯一事件格式 | P0 | 2.7/10.4 | event-system | Contract Test |
| FR-EVENT-002 | 进程内 Event Bus：发布/订阅、顺序保证、错误隔离 | P0 | 2.7 | event-system | Race Test |
| FR-EVENT-003 | 事件分类：持久化/恢复/审计/Hook/可丢弃 | P0 | 10.4 | event-system | Unit |

### Session Store（session-store）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-SESSION-001 | Append-only Event Store（SQLite），按 Session 顺序写入 | P0 | 2.2/10.4 | session-store | Integration |
| FR-SESSION-002 | Session 元数据、Checkpoint 创建与读取 | P0 | 2.1/2.9 | session-store | Unit |
| FR-SESSION-003 | 从事件流重建 Session 状态（Recovery） | P0 | 2.2 | session-store, runtime-core | Recovery Test |
| FR-SESSION-004 | 事件存储增长控制（归档/截断策略，不阻断恢复） | P1 | 14 | session-store | Integration |

### Context Manager（context-manager）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-CONTEXT-001 | 分层上下文模型（见 §6.7/模块 Spec），按层组装请求 | P0 | 2.9 | context-manager | Unit |
| FR-CONTEXT-002 | Token 估算、Reserved Output、Token/Cost Budget 计算 | P0 | 2.9 | context-manager | Unit |
| FR-CONTEXT-003 | 工具输出截断、Bash 头尾保留、Grep 去重、ReadFile 分页、Observation 压缩 | P0 | 2.9 | context-manager, builtin-tools | Golden Test |
| FR-CONTEXT-004 | 自动 Compaction + 手动 `/compact`，压缩前 Checkpoint，压缩后可恢复 | P0 | 2.9 | context-manager | Recovery Test |
| FR-CONTEXT-005 | 关键事实保护（用户目标、计划、文件/行号、未完成任务、权限决定） | P0 | 2.9 | context-manager | Golden Test |
| FR-CONTEXT-006 | 压缩错误检测与回滚（压缩结果不可信时回退 Checkpoint） | P1 | 2.9 | context-manager | Failure Injection |

### Memory System（memory-system）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-MEMORY-001 | 四类记忆：User / Project / Episodic / Procedural | V0.2/P1 | 2.10 | memory-system | Unit |
| FR-MEMORY-002 | SQLite + FTS5 关键词检索，来源/创建时间/最后验证/置信度/过期元数据 | V0.2/P1 | 2.10 | memory-system | Integration |
| FR-MEMORY-003 | 手动编辑/删除、候选记忆审批、项目隔离、敏感信息控制 | V0.2/P1 | 2.10 | memory-system | Security Test |
| FR-MEMORY-004 | 记忆污染控制：不自动写入所有模型输出，仅经审批的候选生效 | V0.2/P0 | 2.10 | memory-system | Security Test |

### Extension System — Skill / Command / Hook（extension-system）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-CMD-001 | Slash Command 框架：固定逻辑命令本地执行（不经模型），Prompt/Skill 命令展开为任务 | P0 | 2.6 | extension-system | Unit |
| FR-CMD-002 | 内置/用户级/项目级命令、参数校验、Alias、Help、冲突处理、Command Hook、权限要求 | P0 | 2.6 | extension-system | Unit |
| FR-HOOK-001 | 基于统一 Event Bus 的 Hook 系统（不散落硬编码回调），支持 §2.7 全部生命周期事件 | P0 | 2.7 | extension-system, event-system | Integration |
| FR-HOOK-002 | Hook 类型：Internal Go / Shell / HTTP；返回 Allow/Deny/Ask/Modify/Continue | P0 | 2.7 | extension-system | Unit |
| FR-HOOK-003 | Hook 顺序/优先级/Matcher/Timeout/失败策略/决策冲突/递归防护/安全边界/审计 | P0 | 2.7 | extension-system | Security Test |
| FR-SKILL-001 | Skill 包结构（SKILL.md+manifest.yaml+资源），用户/项目/内置级 | V0.2/P1 | 2.5 | extension-system | Unit |
| FR-SKILL-002 | 安装/卸载/升级/版本锁定/Discovery/显式调用/Agent 自动选择/延迟加载 | V0.2/P1 | 2.5 | extension-system | Integration |
| FR-SKILL-003 | Skill 权限/Tool/MCP 依赖声明与检查，Skill 不可扩大自身权限 | V0.2/P0 | 2.5 | extension-system, permission-engine | Security Test |
| FR-SKILL-004 | 候选 Skill 自动生成流水线：轨迹→候选→静态检查→回放评测→人工审批→安装 | V0.3/P1 | 2.5 | extension-system, evaluation | Eval |
| FR-CMD-100 | 标杆命令 `/review-pr`、`/review-sql`、`/review-k8s` 注册与触发 | V0.2/P1 | 2.14 | extension-system, evaluation | Eval |

### MCP Client（mcp-client）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-MCP-001 | MCP Client：Server Lifecycle、Capability Negotiation、健康状态、重连 | V0.2/P1 | 2.4 | mcp-client | Integration |
| FR-MCP-002 | stdio 与 Streamable HTTP transport | V0.2/P1 | 2.4 | mcp-client | Contract Test |
| FR-MCP-003 | tools/list、tools/call、resources/list/read、prompts/list/get | V0.2/P1 | 2.4 | mcp-client | Contract Test |
| FR-MCP-004 | Namespace、名称冲突、Schema 转换、输出大小限制、信任级别、权限等级、审计 | V0.2/P0 | 2.4 | mcp-client, permission-engine | Security Test |
| FR-MCP-005 | 外部 Prompt/Resource 安全边界（不可信输入隔离） | V0.2/P0 | 2.4 | mcp-client, permission-engine | Security Test |

### SubAgent / Agent Team（agent-orchestration）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-SUBAGENT-001 | SubAgent 独立身份：Agent ID、System Prompt、上下文、Token/Cost 预算、Tool Allow/Denylist、Skill、权限策略、取消信号、结果结构 | V0.3/P1 | 2.11 | agent-orchestration | Unit |
| FR-SUBAGENT-002 | 父 Agent 结构化任务委派、Expected Output、并发执行与并发上限、超时、取消、失败重试 | V0.3/P1 | 2.11 | agent-orchestration | Integration |
| FR-SUBAGENT-003 | SubAgent 默认返回结构化摘要而非完整日志；工具输出隔离；不默认继承全部父上下文 | V0.3/P0 | 2.11 | agent-orchestration | Contract Test |
| FR-SUBAGENT-004 | 父子预算统计与递归深度限制 | V0.3/P0 | 2.11 | agent-orchestration | Unit |
| FR-TEAM-001 | Agent Team：Team Lead、Task DAG、Member Registry、Mailbox、Artifact Store、Shared State、Team Budget | V1.0/P2 | 2.13 | agent-orchestration | Integration |
| FR-TEAM-002 | 中心化调度：Lead 分解任务、依赖、认领/分配、定向/广播消息、结果集成、冲突处理 | V1.0/P2 | 2.13 | agent-orchestration | Integration |
| FR-TEAM-003 | Team Token/Cost Budget、并发限制、失败重试、超时、人工介入、执行审计 | V1.0/P2 | 2.13 | agent-orchestration, telemetry | Integration |

### Git Worktree（git-worktree）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-WORKTREE-001 | Worktree 生命周期：临时分支、创建、执行、测试、Diff、可选 Commit、Merge/Cherry-pick/Discard、清理 | V0.3/P1 | 2.12 | git-worktree | Integration |
| FR-WORKTREE-002 | 处理边界条件：主仓未提交修改、分支/路径冲突、无修改、相同文件冲突、合并冲突、未跟踪文件、测试失败、清理失败、占用、取消回收 | V0.3/P0 | 2.12 | git-worktree | Failure Injection |
| FR-WORKTREE-003 | 崩溃后孤儿 Worktree 清理与 Worktree 审计事件 | V0.3/P1 | 2.12 | git-worktree | Recovery Test |

### Telemetry（telemetry）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-TELEMETRY-001 | 结构化日志、敏感数据脱敏 | P0 | 10.6 | telemetry | Security Test |
| FR-TELEMETRY-002 | 指标（轮次、工具调用、Token、Cost、错误率）与可选 Trace | P1 | 2.x | telemetry | Unit |
| FR-TELEMETRY-003 | Audit Event Sink（消费 event-system 的审计事件） | P0 | 2.8 | telemetry, event-system | Audit Test |
| FR-TELEMETRY-004 | Usage/Cost 记录（按 Session/Agent/Team 聚合） | P1 | 2.9 | telemetry | Unit |

### Evaluation（evaluation）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-EVAL-001 | Replay：基于持久化事件重放 Session | V0.2/P1 | 12 | evaluation, session-store | Integration |
| FR-EVAL-002 | Eval Case 框架：固定输入+期望结果+评分 | V0.2/P1 | 2.14 | evaluation | Unit |
| FR-EVAL-003 | 三个标杆场景 Eval：review-pr / review-sql / review-k8s | V0.2/P1 | 2.14 | evaluation, extension-system | Eval |
| FR-EVAL-004 | Skill 候选回放评测（支撑 FR-SKILL-004） | V0.3/P1 | 2.5 | evaluation | Eval |

### CLI（cli）
| ID | 描述 | 优先级 | 来源 | 模块 | 验收方式 |
| --- | --- | --- | --- | --- | --- |
| FR-CLI-001 | CLI 入口：启动 Session、提交任务、流式输出、审批交互 | P0 | 2.6 | cli | Integration |
| FR-CLI-002 | 内置固定逻辑命令：/help /model /context /cost /clear /exit /resume /checkpoint | P0 | 2.6 | cli, extension-system | Unit |
| FR-CLI-003 | 命令分派到 extension-system，不在 CLI 层实现核心业务逻辑 | P0 | 4 | cli | 代码审查 + 依赖检查 |
| FR-CLI-004 | `/init` 项目初始化、`/rewind` 回退 Checkpoint | P1 | 2.6 | cli, session-store | Integration |

---

## 6.5 非功能需求（Non-Functional Requirements）

可验证指标优先，避免空洞描述。

| ID | 类别 | 描述与可验证指标 | 优先级 |
| --- | --- | --- | --- |
| NFR-SEC-001 | Security | 所有工具调用必经 Validation→Permission→Hook→Execute→Audit；安全测试覆盖路径穿越、符号链接逃逸、命令注入、密钥泄露、审批绕过，CI 中全部通过 | P0 |
| NFR-SEC-002 | Security | 普通日志中不出现密钥/Token/完整环境变量；脱敏由 telemetry 强制，安全测试验证 | P0 |
| NFR-REL-001 | Reliability | 进程崩溃后可从持久化事件恢复 Session，恢复后状态与崩溃前一致（Recovery Test 通过） | P0 |
| NFR-REL-002 | Reliability | Provider 瞬时错误自动重试（指数退避，默认最多 3 次），不可重试错误快速失败 | P1 |
| NFR-RECOV-001 | Recoverability | 压缩与高风险执行前必有 Checkpoint，可 `/rewind` 回退 | P0 |
| NFR-PERF-001 | Performance | 单工具调用管线（不含工具自身执行）开销 < 5ms（基准测试） | P1 |
| NFR-PERF-002 | Performance | Token 估算与上下文组装在 100 条消息下 < 20ms | P2 |
| NFR-OBS-001 | Observability | 每个 Session 可导出完整事件时间线；关键指标（轮次/工具/Token/Cost/错误率）可查询 | P1 |
| NFR-TEST-001 | Testability | 核心模块支持 Mock/Fake（Provider、Tool、Clock、FS 边界）；`go test -race` 无数据竞争 | P0 |
| NFR-MAINT-001 | Maintainability | 依赖方向无环（CLI→Runtime→Domain Interfaces→Infra）；`go vet`/lint 通过 | P0 |
| NFR-PORT-001 | Portability | 在 macOS 与 Linux 上构建与运行；Sandbox 缺失时核心可降级运行 | P1 |
| NFR-COMPAT-001 | Compatibility | Event Envelope 与持久化 Schema 带版本号，旧事件可被新版本读取（向后兼容） | P1 |
| NFR-LIMIT-001 | Resource Limits | 工具输出、MCP 输出、上下文均有显式上限并截断，不会无界增长 | P0 |
| NFR-COST-001 | Cost Control | 支持 Token/Cost Budget，超限触发 Compaction 或安全终止并落事件 | P0 |

---

## 6.6 系统上下文

参与方与边界（详见 `architecture/SYSTEM_OVERVIEW.md`）：

- **User**：通过 CLI 提交任务、审批高风险操作。
- **CLI**：交互入口，不含核心业务逻辑。
- **Agent Runtime（runtime-core）**：事件驱动状态机，编排模型调用与工具执行。
- **Model Provider**：外部模型 API（不可信网络边界）。
- **Local Tools（builtin-tools）**：在本机或 Sandbox 中执行。
- **MCP Server**：外部进程/服务（**默认不可信**）。
- **Sandbox**：Docker 受控执行环境。
- **Git Repository**：工作区与 Worktree 操作对象。
- **SQLite**：Session/Event/Checkpoint/Memory 本地存储。
- **File System**：受 Workspace Root 与权限边界约束。
- **External Services**：经 MCP 接入（GitHub/Slack/DB/监控等）。

信任边界：User↔CLI（半信任）、Runtime↔Provider（不可信响应，可能含注入）、Runtime↔MCP（不可信）、Runtime↔Tool 执行（受权限/沙箱约束）。详见 `SECURITY_MODEL.md`。

---

## 6.7 核心流程

以下流程在各模块 Spec 与 `FAILURE_AND_RECOVERY.md` 中细化。

1. **普通只读任务**：User→CLI→Runtime(Thinking)→Provider→ToolRequested(Read/Glob/Grep)→Permission(Allow)→Execute→Observing→（无更多调用）→Completed。
2. **代码修改任务**：…→ToolRequested(Write/Edit)→Permission(AskOnce/Allow)→写前 Checkpoint→Execute→Diff→Observing→Completed。
3. **危险 Bash 审批**：…→ToolRequested(Bash rm/force-push)→Permission(High/Critical→Ask)→AwaitingApproval→User 决策→（Allow→Execute / Deny→记录并返回 Observation）。
4. **上下文自动压缩**：Token 接近窗口→PreCompact Checkpoint→Compacting→压缩并校验→PostCompact→继续；校验失败→回滚 Checkpoint。
5. **Session 暂停与恢复**：Paused→持久化→进程退出；重启→读取事件→重建状态→恢复至暂停点。
6. **MCP Tool 调用**：与内置工具同管线，附加 Namespace 解析、信任级别、输出大小限制。
7. **Skill 加载**：Discovery→依赖检查（Tool/MCP/权限）→延迟加载→注入上下文层。
8. **SubAgent 委派**：父 Agent 构造结构化任务+预算+Allowlist→子 Runtime 独立执行→返回结构化摘要+Artifact。
9. **Worktree 修改**：创建临时分支+Worktree→在隔离目录执行→测试→Diff→提交主 Agent 审核→Merge/Discard→清理。
10. **Agent Team 执行**：Lead 构建 Task DAG→分配 Ready 任务→成员执行→Mailbox/Artifact 交换→Reviewing→结果集成。
11. **失败恢复**：见 `FAILURE_AND_RECOVERY.md`（模型/工具/进程/SQLite/MCP/Hook/SubAgent/Team 各类失败）。
12. **用户取消**：Context Cancellation→传播到当前工具/子 Agent→落 Cancelled 事件→可恢复或终止。

---

## 6.8 数据与持久化

核心实体（唯一所有权见 `DATA_OWNERSHIP.md`）：

| 实体 | 拥有模块 | 持久化 | 说明 |
| --- | --- | --- | --- |
| Session | session-store | 是 | 一次交互会话的根 |
| Event | event-system(契约)/session-store(存储) | 是(append-only) | 统一 Envelope |
| Message | session-store | 是 | 模型消息（user/assistant/tool） |
| ToolCall | tool-runtime(契约)/session-store(存储) | 是 | 工具调用请求 |
| ToolResult | tool-runtime/session-store | 是 | 工具结果（可截断） |
| Approval | permission-engine/session-store | 是 | 审批决策记录 |
| Checkpoint | session-store | 是 | 可回退快照点 |
| Memory | memory-system | 是 | 四类记忆 |
| SkillMetadata | extension-system | 是 | Skill manifest 索引 |
| AgentDefinition | agent-orchestration | 是 | SubAgent/角色定义 |
| AgentInstance | agent-orchestration | 运行态+事件 | 运行中的 Agent |
| Task | agent-orchestration | 是 | Team Task 节点 |
| Team | agent-orchestration | 是 | 团队 |
| Worktree | git-worktree | 是(登记) | Worktree 登记表 |
| Artifact | agent-orchestration | 是 | 产物存储引用 |
| UsageRecord | telemetry | 是 | Token/Cost 记账 |

存储介质：第一版统一 SQLite（含 FTS5 用于 Memory），文件系统用于 Skill 包与 Artifact 大对象。Schema 带版本号（NFR-COMPAT-001）。

---

## 6.9 版本路线

不把所有高级能力都列为 MVP 阻断项。

- **MVP（M1–M5）**：单 Agent 闭环。runtime-core、model-provider（Mock+OpenAI）、tool-runtime、builtin-tools、permission-engine（前三层+审批审计）、event-system、session-store、context-manager、telemetry、cli。可演示：只读任务、代码修改、危险 Bash 审批、自动压缩、Session 恢复、取消。
- **V0.2**：sandbox（Docker）、mcp-client、extension-system（Skill + 标杆命令）、memory-system、evaluation（Replay + 三场景 Eval）。
- **V0.3**：agent-orchestration（SubAgent）、git-worktree、Skill 自动生成流水线。
- **V1.0**：agent-orchestration（Agent Team 完整）、Team 预算/审计/结果集成、完整 Eval 与展示材料。

详见 `planning/ROADMAP.md` 与 `planning/MILESTONES.md`。
