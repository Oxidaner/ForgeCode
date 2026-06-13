你现在是 ForgeCode 项目的首席架构师、技术负责人和交付规划负责人。

你的任务不是立即编写业务代码，而是基于当前仓库状态，为 ForgeCode 完成一次完整的架构分析、模块划分和开发规划，并在仓库中生成可以直接指导后续实现的：

1. 项目总览与总体架构 Spec
2. 每个模块独立的详细 Spec
3. 每个模块独立的 Task 文档
4. 每个模块独立的 Checklist
5. 跨模块依赖关系与实施顺序
6. 架构决策记录 ADR
7. 需求—模块—任务—验收标准追踪矩阵
8. 12 周实施路线图
9. 风险清单和待验证问题

本阶段只允许创建或修改 Markdown、YAML、JSON 等设计和规划文档。

除非当前仓库已经存在必要的最小骨架，否则不要实现生产代码、不要安装依赖、不要生成大量空 Go 文件、不要执行不可逆操作。

------

# 一、项目背景

项目名称：

```text
ForgeCode
```

项目定位：

```text
A Model-Agnostic Coding Agent Runtime Written in Go
```

ForgeCode 是一个使用 Go 从零实现的、模型无关的 Coding Agent Runtime。

它不是简单封装某个 Agent SDK，也不是一个只会调用工具的聊天机器人。项目需要自主实现 Agent Runtime 的核心控制逻辑，并支持安全、可恢复、可扩展、可观测和可并行的代码任务执行。

项目主要参考方向：

```text
Claude Code：
Agent Runtime
编程工具
权限控制
上下文管理
Hook
SubAgent
Git Worktree
Agent Teams

Hermes：
跨会话记忆
Skill 管理
经验沉淀
长期 Agent 能力

ForgeCode 自有场景：
PR 变更审查
SQL 变更审查
Kubernetes 配置变更审查
```

最终项目需要体现：

```text
不是调用 Agent SDK
而是自主实现一个：
可扩展
可审计
可恢复
可控制
支持安全并行执行
支持多 Agent 协作
的 Agent Runtime
```

------

# 二、必须涵盖的产品能力

以下能力全部必须进入总体设计，但不要求每项能力都成为独立模块。你需要根据职责边界、依赖方向和内聚性自行划分模块。

## 2.1 六大核心编程工具

必须支持：

```text
ReadFile
WriteFile
EditFile
Bash
Glob
Grep
```

至少覆盖：

- 文件读取
- 文件创建与覆盖
- 精确局部编辑
- Shell 命令执行
- 文件模式匹配
- 关键字和正则搜索
- 输出截断
- 超时控制
- 错误分类
- 权限检查
- 工具调用审计
- 工具执行前后 Hook
- 文件修改前 Checkpoint
- 修改结果 Diff
- 失败时保持原文件不变

所有内置工具必须实现统一 Tool 接口，并通过 Tool Registry 注册。

MCP 工具和内置工具进入 Agent Runtime 后，应尽量使用统一的 Tool Descriptor、调用流程、权限检查和审计机制。

------

## 2.2 自主 Agent Loop

需要基于工具调用循环实现自主任务执行。

Agent 能够：

```text
理解任务
制定或维护计划
调用模型
解析 Tool Call
检查权限
执行工具
读取 Observation
更新上下文
判断是否继续
压缩上下文
完成、暂停、取消或失败
```

Agent Loop 不应只是一个无法恢复的无限循环，而应设计成显式状态机或事件驱动运行时。

至少考虑以下状态：

```text
Created
Initializing
Thinking
ToolRequested
AwaitingApproval
ToolExecuting
Observing
Compacting
Paused
Completed
Failed
Cancelled
```

必须支持：

- 最大轮次
- 最大工具调用次数
- Token 预算
- 成本预算
- Deadline
- Context Cancellation
- 用户暂停
- 用户取消
- 会话恢复
- 重复工具调用检测
- 相同错误循环检测
- Provider 错误重试
- 非法 Tool Call 处理
- 中间状态持久化
- Append-only 事件记录

不要只使用“基于 ReAct”概括实现，必须给出清晰状态、事件、输入、输出、错误和恢复语义。

------

## 2.3 模型 Provider 抽象

ForgeCode 必须是 Model-Agnostic。

至少设计：

```text
OpenAI Provider
Anthropic Provider
OpenAI-Compatible Provider
Mock Provider
```

Provider 层至少处理：

- 普通响应
- Streaming
- Tool Calling
- 多 Tool Call
- Token Usage
- Stop Reason
- Provider Error
- Rate Limit
- Retry
- Timeout
- Model Capability
- Context Window
- Structured Output
- Provider-specific message conversion

不能让 Agent Runtime 直接依赖某一家模型提供商的请求和响应结构。

------

## 2.4 MCP 协议接入

ForgeCode 需要实现 MCP Client 能力，可以通过配置挂载外部 MCP Server。

至少规划：

```text
Server Lifecycle
Capability Negotiation
tools/list
tools/call
resources/list
resources/read
prompts/list
prompts/get
stdio transport
Streamable HTTP transport
```

需要考虑：

- MCP Server 启动与关闭
- 连接健康状态
- 初始化失败
- 超时
- Server 重连
- Tool 名称冲突
- Namespace
- Schema 转换
- MCP Tool 权限等级
- MCP Tool 审计
- MCP Server 信任级别
- MCP Tool 输出大小限制
- 外部 Prompt 和 Resource 的安全边界

典型接入场景：

```text
GitHub
Slack
PostgreSQL
MySQL
监控平台
浏览器服务
12306 等第三方服务
```

设计时不要假设所有 MCP Server 都可信。

------

## 2.5 Skill 技能包系统

Skill 是可安装、可发现、可加载、可测试和可版本化的技能包。

Skill 不只是 Prompt 文件，应允许包含：

```text
SKILL.md
manifest.yaml
Prompt
工具声明
权限声明
参考资料
模板
脚本
示例
测试用例
其他资源文件
```

至少支持：

- 用户级 Skill
- 项目级 Skill
- 内置 Skill
- 安装
- 卸载
- 升级
- 版本锁定
- Skill Discovery
- 显式调用
- Agent 自动选择
- 延迟加载
- Tool 依赖检查
- MCP 依赖检查
- Skill 权限声明
- Skill 冲突处理
- Skill 测试
- Skill 版本兼容
- Skill 来源追踪

需要设计候选能力：

```text
成功任务轨迹
    ↓
提取可复用步骤
    ↓
生成候选 Skill
    ↓
静态检查
    ↓
回放评测
    ↓
人工审批
    ↓
安装到 Skill Registry
```

自动生成的 Skill 不允许未经评测和审批直接生效。

------

## 2.6 Slash Command 框架

至少支持：

```text
/help
/init
/model
/compact
/context
/cost
/memory
/skills
/mcp
/hooks
/permissions
/agents
/team
/worktree
/resume
/checkpoint
/rewind
/clear
/exit
```

命令分为：

```text
固定逻辑命令
Prompt 命令
Skill 命令
```

固定逻辑命令应直接由本地程序执行，不经过模型。

Prompt 和 Skill 命令可以展开为任务输入，再进入 Agent Runtime。

需要支持：

- 内置命令
- 用户级自定义命令
- 项目级自定义命令
- 命令参数
- 参数校验
- 自动补全所需元数据
- 命令冲突
- Alias
- Help
- 权限要求
- Command Hook

------

## 2.7 Hook 生命周期系统

Hook 应基于统一事件模型实现，不要在每个模块中散落硬编码回调。

至少支持以下生命周期事件：

```text
SessionStart
SessionEnd
UserPromptSubmit

PreModelCall
PostModelCall
ModelCallFailed

PreToolUse
PostToolUse
ToolFailure

ApprovalRequested
ApprovalResolved

PreCompact
PostCompact

MemoryRead
MemoryWrite

SubAgentStart
SubAgentStop

WorktreeCreate
WorktreeRemove

TeamCreated
TeamClosed

TaskCreated
TaskAssigned
TaskCompleted
TaskFailed
```

Hook 类型至少考虑：

```text
Internal Go Hook
Shell Hook
HTTP Hook
```

Hook 可以返回：

```text
Allow
Deny
Ask
Modify
Continue
```

需要考虑：

- Hook 执行顺序
- 优先级
- Matcher
- Timeout
- Hook 失败策略
- 是否允许修改输入
- 是否允许添加上下文
- 多个 Hook 决策冲突
- Hook 递归调用
- 安全边界
- 审计日志

------

## 2.8 五层纵深权限防御

必须设计成独立、可测试的 Permission Engine，而不是零散的字符串判断。

五层至少包括：

### 第一层：Schema 与输入校验

检查：

- 参数是否合法
- 必填字段
- 长度限制
- 路径格式
- 空字节
- 非法编码
- 命令字段结构
- JSON Schema
- Tool Call 注入风险

### 第二层：资源边界

检查：

- Workspace Root
- 可读目录
- 可写目录
- 敏感目录
- 路径穿越
- 符号链接逃逸
- 隐藏文件
- 密钥文件
- 环境变量
- 用户目录
- 系统目录

### 第三层：操作风险策略

需要支持风险等级：

```text
Low
Medium
High
Critical
```

需要支持决策：

```text
Allow
AskOnce
AskAlways
Deny
```

需要设计 Bash 命令分析，而不是只匹配完整命令字符串。

至少分析：

- 可执行程序
- 参数
- 管道
- 重定向
- 子 Shell
- 命令替换
- 网络访问
- 文件删除
- Git 危险操作
- Docker
- Kubernetes
- 数据库写入
- Force Push
- 下载后执行
- 提权命令

### 第四层：运行时沙箱

规划：

- Docker Sandbox
- 工作目录挂载
- 只读挂载
- 网络控制
- CPU 限制
- 内存限制
- PID 限制
- 执行时间限制
- 环境变量过滤
- 临时目录
- 进程回收

不要求第一版自行实现完整容器运行时。

### 第五层：人工审批与审计

记录：

- 谁发起操作
- 哪个 Session
- 哪个 Agent
- 哪个 Tool
- 原始参数
- 风险等级
- 规则命中原因
- 用户批准或拒绝
- 执行结果
- 文件变更
- 时间
- 审计事件 ID

必须明确权限优先级和冲突决策规则。

------

## 2.9 上下文压缩与 Token 管理

上下文必须分层管理。

至少考虑：

```text
System Context
Project Instructions
Active Skill
User Task
Current Plan
Working Memory
Recent Messages
Tool Results
Retrieved Memory
Compacted History
```

必须支持：

- 模型上下文窗口识别
- Token 估算
- Reserved Output
- Token Budget
- Cost Budget
- 工具输出截断
- Bash 输出头尾保留
- Grep 结果去重
- ReadFile 分页
- JSON 结果提取
- Observation 压缩
- 会话自动 Compaction
- 手动 `/compact`
- 压缩前 Checkpoint
- 压缩后可恢复
- 关键事实保护
- 文件路径和行号保护
- 未完成任务保护
- 用户原始目标保护

压缩结果至少保留：

```text
用户目标
当前计划
已完成步骤
关键事实
关键文件
代码修改
重要决策
权限决定
失败尝试
未解决问题
下一步
```

需要说明压缩错误如何检测和回滚。

------

## 2.10 跨会话记忆

至少规划四类记忆：

```text
User Memory
Project Memory
Episodic Memory
Procedural Memory
```

需要支持：

- 用户级记忆
- 项目级记忆
- 会话摘要
- 历史任务
- 项目架构
- 构建命令
- 代码规范
- 常见失败
- 用户偏好
- Skill 关联
- 关键字检索
- 可选语义检索
- 来源追踪
- 创建时间
- 最后验证时间
- 置信度
- 过期时间
- 手动编辑
- 手动删除
- 候选记忆审批
- 记忆污染控制
- 项目隔离
- 隐私和敏感信息控制

第一版优先考虑：

```text
文件系统
SQLite
SQLite FTS5
```

不要默认引入向量数据库。只有明确说明关键词检索不足且存在真实需求时，才将 Embedding 检索列为后续阶段。

------

## 2.11 SubAgent

SubAgent 必须具有独立：

```text
Agent ID
System Prompt
上下文
Token Budget
成本预算
Tool Allowlist
Tool Denylist
Skill
权限策略
执行状态
取消信号
结果结构
```

典型 Agent：

```text
ExploreAgent
TestAgent
ReviewAgent
ImplementAgent
SecurityAgent
DocumentationAgent
```

需要支持：

- 父 Agent 委派任务
- 结构化任务输入
- Expected Output
- 独立上下文
- 并发执行
- 并发数限制
- 超时
- 取消
- 失败重试
- 结果摘要
- Artifact 返回
- 工具输出隔离
- 不默认继承全部父上下文
- 父子 Agent 预算统计
- 递归深度限制

SubAgent 返回给父 Agent 的默认内容应是结构化摘要，而不是完整运行日志。

------

## 2.12 Git Worktree 并行隔离

会修改代码的并行 Agent 应支持独立 Git Worktree。

至少考虑：

```text
创建临时分支
创建 Worktree
在 Worktree 中执行
运行测试
生成 Diff
可选 Commit
提交主 Agent 审核
Merge
Cherry-pick
Discard
清理 Worktree
```

必须处理：

- 主仓库未提交修改
- 分支名称冲突
- Worktree 路径冲突
- Agent 没有产生修改
- 多 Agent 修改相同文件
- 合并冲突
- 未跟踪文件
- 测试失败
- 清理失败
- 进程仍占用目录
- 任务取消后的回收
- 崩溃后的孤儿 Worktree 清理
- Worktree 审计事件

------

## 2.13 Agent Teams

Agent Team 不是自由聊天群，而应建立在：

```text
Team Lead
Task DAG
Member Registry
Mailbox
Artifact Store
Shared Project State
Team Budget
```

之上。

需要支持：

- 团队创建和关闭
- 长期角色定义
- Lead 分解任务
- Task DAG
- 任务依赖
- 任务认领或分配
- Agent 状态
- 定向消息
- 广播消息
- 共享 Artifact
- 并发限制
- Team Token Budget
- Team Cost Budget
- 失败重试
- 超时
- 人工介入
- 结果集成
- 冲突处理
- 团队执行审计

任务状态至少包括：

```text
Pending
Blocked
Ready
Assigned
Running
Reviewing
Completed
Failed
Cancelled
```

需要明确区分：

```text
SubAgent：
一次性任务委派，执行完成后返回父 Agent

Agent Team：
具有长期角色、共享任务图、成员通信和结果集成
```

第一版 Agent Team 只需要支持受控的中心化调度，不要设计复杂的去中心化协商或多 Agent 自由辩论。

------

## 2.14 标杆 Skill 和评测场景

至少规划以下三个 Skill：

```text
/review-pr
/review-sql
/review-k8s
```

### Review PR

可能使用：

```text
ExploreAgent
ReviewAgent
SecurityAgent
TestAgent
```

输出：

- 变更范围
- 调用链影响
- 风险等级
- 证据
- 文件和行号
- 测试结果
- 建议修改
- 阻断项

### Review SQL

检查：

- DDL 风险
- 全表锁
- 大表修改
- 缺失索引
- 无条件 UPDATE
- 无条件 DELETE
- 数据兼容
- 回滚方案
- 发布顺序

### Review Kubernetes

检查：

- 镜像变更
- Resource Requests 和 Limits
- 探针
- Service
- Ingress
- RBAC
- Secret
- ConfigMap
- 滚动更新
- 高可用
- Dry-run
- Schema 校验
- 回滚能力

这些场景必须成为后续 Eval 的主要来源。

------

# 三、工作方式

## 3.1 第一原则

本阶段先设计，后编码。

你必须：

1. 阅读当前仓库。
2. 阅读已有 `README`、`AGENTS.md`、设计文档、代码、配置和测试。
3. 判断项目处于空仓库、骨架阶段还是已有实现阶段。
4. 识别已有约束和可复用内容。
5. 自主划分模块。
6. 输出文档。
7. 检查所有文档之间的一致性。
8. 最后汇报创建和修改了哪些文件。

不要一开始机械地创建十三个模块。

模块划分应优先根据：

- 单一职责
- 高内聚
- 低耦合
- 清晰依赖方向
- 可独立测试
- 可替换性
- 安全边界
- 并发边界
- 数据所有权
- 生命周期
- 后续扩展性

同一能力可以横跨多个模块，但必须明确哪个模块拥有核心抽象，哪个模块只是调用方。

------

## 3.2 可使用 Subagent 的分析方式

在环境支持 Subagent 时，启动多个只读分析 Agent 并行完成以下工作：

### Architecture Analyst

分析：

- 当前仓库结构
- 模块候选
- 依赖方向
- 核心抽象
- 状态机
- 数据所有权

### Security Analyst

分析：

- Permission Engine
- Bash 风险
- 路径安全
- MCP 信任边界
- Sandbox
- 审批与审计
- 记忆隐私

### Runtime Analyst

分析：

- Agent Loop
- Context
- Session
- Event Store
- Recovery
- Provider
- Tool Runtime
- Hook

### Orchestration Analyst

分析：

- SubAgent
- Worktree
- Task DAG
- Agent Teams
- 并发
- 取消
- Artifact

### Quality Analyst

分析：

- 测试策略
- Eval
- Traceability
- Checklist
- Definition of Done
- 故障注入
- 性能指标

所有 Subagent 只读分析，不修改相同文件。

由主 Agent 汇总结果并负责最终文档，避免产生互相矛盾的设计。

如果当前环境不支持 Subagent，则按上述角色顺序自行完成分析。

------

# 四、期望的模块划分

以下只是候选边界，不是强制最终答案。

你需要根据当前仓库实际情况调整、合并或拆分，并在架构文档中说明理由。

候选模块：

```text
runtime-core
model-provider
tool-runtime
builtin-tools
permission-engine
sandbox
event-system
session-store
context-manager
memory-system
extension-system
mcp-client
agent-orchestration
git-worktree
telemetry
evaluation
cli
```

其中可能的内部能力：

```text
runtime-core
├── Agent Loop
├── Agent State Machine
├── Runtime Coordinator
├── Cancellation
└── Recovery

extension-system
├── Skill
├── Slash Command
└── Hook

agent-orchestration
├── Task
├── SubAgent
├── Agent Team
├── Mailbox
└── Artifact Store
```

你必须检查是否存在以下反模式：

- 模块过多且每个模块只有一个文件
- `common`、`utils`、`manager` 成为垃圾桶
- Runtime 直接依赖具体 Provider
- Tool 直接弹出 UI 审批
- Permission Engine 直接执行命令
- MCP 工具绕过权限系统
- Hook 绕过事件总线
- Session 同时承担 Memory
- Memory 自动写入所有模型输出
- SubAgent 复制父 Agent 完整上下文
- Agent Team 与 SubAgent 没有清晰区别
- Worktree 逻辑混入 Git Tool
- CLI 成为业务逻辑核心
- 多个模块分别维护自己的事件格式
- 领域层依赖 CLI、TUI 或具体数据库实现
- 核心接口提前抽象得过度复杂

------

# 五、需要生成的目录和文档

在没有更合理的现有文档规范时，使用以下结构：

```text
docs/
├── README.md
├── architecture/
│   ├── SYSTEM_OVERVIEW.md
│   ├── MODULE_MAP.md
│   ├── DEPENDENCY_GRAPH.md
│   ├── EVENT_MODEL.md
│   ├── DATA_OWNERSHIP.md
│   ├── SECURITY_MODEL.md
│   └── FAILURE_AND_RECOVERY.md
├── specs/
│   ├── 00-master/
│   │   └── SPEC.md
│   └── <module-id>/
│       └── SPEC.md
├── tasks/
│   ├── MASTER_TASKS.md
│   └── <module-id>/
│       └── TASKS.md
├── checklists/
│   ├── MASTER_CHECKLIST.md
│   └── <module-id>/
│       └── CHECKLIST.md
├── adr/
│   ├── README.md
│   ├── ADR-0001-*.md
│   └── ...
├── planning/
│   ├── ROADMAP.md
│   ├── MILESTONES.md
│   ├── TRACEABILITY_MATRIX.md
│   ├── RISK_REGISTER.md
│   ├── OPEN_QUESTIONS.md
│   └── GLOSSARY.md
└── templates/
    ├── SPEC_TEMPLATE.md
    ├── TASK_TEMPLATE.md
    ├── CHECKLIST_TEMPLATE.md
    └── ADR_TEMPLATE.md
```

如果当前仓库已经有成熟文档结构，应优先遵循现有规范，但必须保证以上信息全部被覆盖。

------

# 六、总 Spec 要求

创建：

```text
docs/specs/00-master/SPEC.md
```

总 Spec 必须包含：

## 6.1 文档元数据

```text
Status
Version
Owners
Last Updated
Target Release
Related ADRs
```

## 6.2 项目目标

明确：

- 要解决的问题
- 目标用户
- 核心价值
- 面试和开源展示价值
- 为什么不直接依赖现有 Agent SDK
- 哪些基础设施可以使用第三方库
- 哪些核心逻辑必须自主实现

## 6.3 非目标

至少说明：

- 第一版不做完整 IDE
- 第一版不做复杂 TUI
- 第一版不支持任意深度多 Agent 嵌套
- 第一版不做去中心化 Agent 协商
- 第一版不默认引入向量数据库
- 第一版不实现自己的容器运行时
- 第一版不保证兼容所有 Provider 私有能力
- 第一版不允许 Agent 自动执行未经审批的高风险操作

## 6.4 功能需求

所有需求使用稳定 ID，例如：

```text
FR-RUNTIME-001
FR-TOOL-001
FR-PERM-001
FR-MCP-001
FR-SKILL-001
FR-CONTEXT-001
FR-MEMORY-001
FR-SUBAGENT-001
FR-WORKTREE-001
FR-TEAM-001
```

每条功能需求必须包含：

- 描述
- 优先级
- 来源
- 所属模块
- 验收方式

## 6.5 非功能需求

至少涵盖：

```text
Security
Reliability
Recoverability
Performance
Observability
Testability
Maintainability
Portability
Compatibility
Resource Limits
Cost Control
```

非功能需求也必须使用 ID，例如：

```text
NFR-SEC-001
NFR-REL-001
NFR-PERF-001
```

尽量给出可验证指标，不要只写“高性能”“高可用”。

## 6.6 系统上下文

说明：

- 用户
- CLI
- Agent Runtime
- Model Provider
- Local Tools
- MCP Server
- Sandbox
- Git Repository
- SQLite
- File System
- External Services

## 6.7 核心流程

至少覆盖：

1. 普通只读任务
2. 代码修改任务
3. 危险 Bash 审批
4. 上下文自动压缩
5. Session 暂停与恢复
6. MCP Tool 调用
7. Skill 加载
8. SubAgent 委派
9. Worktree 修改
10. Agent Team 执行
11. 失败恢复
12. 用户取消

## 6.8 数据与持久化

定义：

- Session
- Event
- Message
- Tool Call
- Tool Result
- Approval
- Checkpoint
- Memory
- Skill Metadata
- Agent Definition
- Agent Instance
- Task
- Team
- Worktree
- Artifact
- Usage Record

## 6.9 版本路线

明确：

```text
MVP
V0.2
V0.3
V1.0
```

不要把所有高级能力全部定义为 MVP 阻断项。

------

# 七、每个模块 Spec 的统一格式

每个模块必须拥有独立：

```text
docs/specs/<module-id>/SPEC.md
```

每份模块 Spec 至少包含：

## 7.1 模块信息

```text
Module ID
Module Name
Status
Owner
Dependencies
Dependents
Related Requirements
Related ADRs
```

## 7.2 Purpose

说明模块为什么存在。

## 7.3 Scope

明确模块负责什么。

## 7.4 Non-goals

明确模块不负责什么。

## 7.5 Responsibilities

列出模块职责。

## 7.6 Public Interfaces

使用 Go 风格伪代码定义关键接口，但暂不实现。

例如：

```go
type Tool interface {
    Descriptor() ToolDescriptor
    Execute(ctx context.Context, input json.RawMessage) (ToolResult, error)
}
```

必须避免为未来假想需求过度抽象。

## 7.7 Domain Model

定义主要实体、值对象和枚举。

## 7.8 State Machine

存在生命周期的模块必须给出状态机。

## 7.9 Core Flows

描述正常流程和异常流程。

## 7.10 Configuration

列出配置项、默认值、作用域和敏感性。

## 7.11 Persistence

说明数据是否持久化、由谁拥有、如何迁移。

## 7.12 Concurrency

说明：

- 是否线程安全
- 锁或 Channel 的边界
- 并发限制
- 取消传播
- Race 风险
- 顺序保证
- 幂等要求

## 7.13 Error Model

不要只返回字符串错误。

需要定义错误类别，例如：

```text
ValidationError
PermissionDenied
ApprovalRequired
TimeoutError
CancelledError
ProviderError
ToolExecutionError
SandboxError
PersistenceError
ConflictError
RecoveryError
```

## 7.14 Security

说明：

- 信任边界
- 攻击面
- 权限检查点
- 敏感数据
- 日志脱敏
- Prompt Injection 风险
- Tool Output Injection 风险

## 7.15 Observability

说明：

- Log
- Metric
- Trace
- Audit Event
- Usage
- Cost

## 7.16 Testing Strategy

至少覆盖：

- Unit Test
- Integration Test
- Contract Test
- Race Test
- Failure Injection
- Security Test
- Golden Test
- Eval

## 7.17 Acceptance Criteria

必须是可检查、可测试的条件。

## 7.18 Risks

列出模块实现风险。

## 7.19 Open Questions

只保留真正需要后续验证的问题。

------

# 八、Task 文档要求

为每个模块创建：

```text
docs/tasks/<module-id>/TASKS.md
```

同时创建：

```text
docs/tasks/MASTER_TASKS.md
```

## 8.1 Task ID

使用稳定 ID：

```text
FC-RT-001
FC-TOOL-001
FC-PERM-001
FC-CTX-001
FC-MCP-001
FC-SKILL-001
FC-MEM-001
FC-SUB-001
FC-WT-001
FC-TEAM-001
```

## 8.2 每个 Task 必须包含

```text
ID
Title
Module
Priority
Milestone
Status
Size
Dependencies
Related Requirements
Related Spec Sections
Description
Implementation Notes
Files or Packages Likely Affected
Tests Required
Security Considerations
Acceptance Criteria
Definition of Done
Evidence
```

其中：

```text
Priority：P0 / P1 / P2
Size：XS / S / M / L / XL
Status：Backlog / Ready / In Progress / Blocked / Done
```

不要给出虚假的小时数。

## 8.3 Task 粒度

一个 Task 应满足：

- 能够独立理解
- 具有明确输入输出
- 具有可测试验收条件
- 通常只聚焦一个主要结果
- 不只是“实现某模块”
- 不把整个 Agent Runtime 塞进一个任务
- 不细碎到“创建一个文件”
- XL Task 必须继续拆分

## 8.4 Task 类型

可使用：

```text
Spike
Architecture
Implementation
Test
Security
Documentation
Evaluation
Migration
Refactor
```

Spike 必须包含：

- 需要回答的问题
- 最小实验
- 输出决策
- 结束条件

## 8.5 Master Task

`MASTER_TASKS.md` 必须：

- 按 Milestone 分组
- 标注关键路径
- 标注可并行任务
- 标注阻塞关系
- 给出 Task DAG
- 识别必须先完成的架构任务
- 识别适合不同 Codex Agent 并行完成的任务
- 避免让多个 Agent 同时修改相同核心文件

------

# 九、Checklist 文档要求

为每个模块创建：

```text
docs/checklists/<module-id>/CHECKLIST.md
```

同时创建：

```text
docs/checklists/MASTER_CHECKLIST.md
```

Checklist 不得只是重复 Task 标题。

每项使用：

```text
- [ ] 检查项
```

必要时附带：

```text
Evidence:
Related Task:
Related Requirement:
Blocking:
```

## 9.1 模块 Checklist 必须包括

### Design Ready

- 职责边界已明确
- 非目标已明确
- 接口已定义
- 状态机已定义
- 错误模型已定义
- 安全边界已定义
- 并发语义已定义
- 持久化所有权已定义
- 依赖方向无环
- ADR 已完成

### Implementation Ready

- Task 已拆分
- P0 依赖已满足
- 配置已定义
- 测试策略已定义
- Mock/Fake 边界已定义
- Migration 策略已定义
- 回滚策略已定义

### Implementation Complete

- 核心路径完成
- 异常路径完成
- Context Cancellation 生效
- 错误可识别
- 事件已记录
- 配置有默认值
- 敏感数据未写入普通日志

### Test Complete

- Unit Test
- Integration Test
- Race Test
- Failure Test
- Security Test
- Contract Test
- Eval Case
- 覆盖关键恢复路径

### Documentation Complete

- Spec 更新
- Task 状态更新
- ADR 更新
- README 更新
- 配置示例更新
- 操作手册更新
- 已知限制更新

### Release Ready

- 所有 P0 验收通过
- 没有未处理 Critical 风险
- 关键指标可观察
- 升级与回滚经过验证
- Demo 可复现
- Evidence 已记录

------

# 十、架构文档要求

## 10.1 SYSTEM_OVERVIEW.md

必须给出：

- 项目目标
- 系统边界
- 核心组件
- 关键流程
- 运行模式
- 部署模式
- 技术选型
- MVP 边界

## 10.2 MODULE_MAP.md

每个模块记录：

```text
名称
职责
拥有的数据
公开接口
依赖
调用方
风险等级
MVP 是否需要
```

## 10.3 DEPENDENCY_GRAPH.md

使用 Mermaid 输出模块依赖图。

依赖方向建议遵循：

```text
CLI
    ↓
Application / Runtime Coordinator
    ↓
Core Domain Interfaces
    ↓
Infrastructure Implementations
```

必须检查并说明是否存在循环依赖。

## 10.4 EVENT_MODEL.md

定义统一 Event Envelope：

```text
Event ID
Event Type
Timestamp
Session ID
Agent ID
Task ID
Team ID
Correlation ID
Causation ID
Sequence
Payload
Schema Version
```

说明：

- 哪些事件必须持久化
- 哪些事件用于恢复
- 哪些事件用于审计
- 哪些事件用于 Hook
- 哪些事件可以丢弃
- 顺序和幂等语义

## 10.5 DATA_OWNERSHIP.md

明确每类数据由哪个模块拥有。

任何核心实体不能被多个模块同时直接写入。

## 10.6 SECURITY_MODEL.md

至少覆盖：

- Threat Model
- Trust Boundaries
- Prompt Injection
- Tool Output Injection
- Path Traversal
- Symlink Escape
- Command Injection
- Secret Leakage
- MCP Supply Chain
- Malicious Skill
- Malicious Hook
- Sandbox Escape
- Approval Bypass
- Memory Poisoning
- Audit Tampering

## 10.7 FAILURE_AND_RECOVERY.md

至少分析：

- 模型请求中断
- Streaming 中断
- Tool 执行中断
- 进程崩溃
- SQLite 写入失败
- Worktree 创建失败
- MCP 连接断开
- Hook 超时
- SubAgent 失联
- Team Member 失败
- 用户取消
- Context Compaction 失败
- 审批后执行前崩溃
- 重启后的恢复流程

------

# 十一、ADR 要求

至少评估是否需要生成以下 ADR：

```text
ADR-0001：采用事件驱动 Agent Runtime
ADR-0002：使用 Append-only Event Store
ADR-0003：Provider 与 Runtime 解耦
ADR-0004：内置 Tool 和 MCP Tool 使用统一描述
ADR-0005：Permission Engine 独立于 Tool Executor
ADR-0006：SQLite 作为初始本地状态存储
ADR-0007：记忆先使用 FTS5 而非向量数据库
ADR-0008：并行写任务使用 Git Worktree
ADR-0009：Agent Team 使用中心化 Task DAG
ADR-0010：Skill 自动生成需要评测和人工审批
ADR-0011：Hook 基于统一 Event Bus
ADR-0012：第一版 Sandbox 使用 Docker
```

每份 ADR 使用：

```text
Title
Status
Context
Decision
Alternatives Considered
Consequences
Security Impact
Operational Impact
Revisit Conditions
```

不要为了凑数量创建没有实际决策价值的 ADR。

------

# 十二、追踪矩阵

创建：

```text
docs/planning/TRACEABILITY_MATRIX.md
```

至少包含：

| Requirement | Module | Spec Section | Task | Test/Eval | Checklist | Status |
| ----------- | ------ | ------------ | ---- | --------- | --------- | ------ |
|             |        |              |      |           |           |        |

所有 P0 Requirement 必须至少对应：

- 一个负责模块
- 一个 Spec 章节
- 一个 Task
- 一个测试或 Eval
- 一个 Checklist 验收项

不得出现：

- 没有 Task 的 P0 需求
- 没有验收标准的 Task
- 没有负责模块的需求
- Checklist 无法追溯
- Spec 中存在未进入规划的重要能力

------

# 十三、路线图

创建：

```text
docs/planning/ROADMAP.md
docs/planning/MILESTONES.md
```

按 12 周规划，但使用依赖关系而不是仅按时间排序。

参考阶段：

## Week 1—2：单 Agent 可运行

目标：

- Provider 抽象
- Streaming
- Agent Loop
- Tool Registry
- ReadFile
- Glob
- Grep
- Bash
- Session
- CLI

## Week 3—4：编辑与安全

目标：

- WriteFile
- EditFile
- Diff
- Checkpoint
- Permission Engine
- Approval
- Hook Event Bus
- Docker Sandbox

## Week 5：上下文与预算

目标：

- Token 估算
- Tool Result 截断
- Observation 压缩
- Auto Compaction
- Cost Budget
- Loop Detection

## Week 6：扩展系统

目标：

- Skill
- Slash Command
- Hook 完善
- Review Skills

## Week 7：MCP

目标：

- MCP Lifecycle
- stdio
- Streamable HTTP
- Tools
- Resources
- Prompts
- MCP Permission

## Week 8：记忆与恢复

目标：

- User Memory
- Project Memory
- Session Resume
- FTS
- Candidate Memory

## Week 9：SubAgent

目标：

- Agent Definition
- Independent Context
- Delegation
- Concurrency
- Structured Result

## Week 10：Worktree

目标：

- Worktree Lifecycle
- Branch
- Diff
- Commit
- Conflict
- Cleanup

## Week 11：Agent Teams

目标：

- Team Lead
- Task DAG
- Mailbox
- Artifact Store
- Team Budget
- Result Integration

## Week 12：Eval 和项目展示

目标：

- Replay
- Eval Cases
- Security Tests
- Metrics
- Demo
- README
- Architecture Diagram
- Technical Article

你可以根据实际依赖调整，但必须说明调整原因。

------

# 十四、风险清单

创建：

```text
docs/planning/RISK_REGISTER.md
```

每个风险至少包含：

```text
Risk ID
Description
Category
Probability
Impact
Trigger
Mitigation
Contingency
Owner Module
Related Tasks
Status
```

至少分析：

- 项目范围过大
- 模块过度设计
- Agent Loop 不稳定
- Provider 差异
- Tool Schema 不一致
- Bash 安全误判
- 路径逃逸
- 审批绕过
- 上下文压缩丢失信息
- Memory 污染
- MCP Server 不可信
- Hook 执行任意代码
- SubAgent Token 成本爆炸
- Worktree 孤儿资源
- Agent Team 调度复杂
- 事件存储无限增长
- SQLite 并发限制
- Eval 不足以证明能力
- Demo 场景过于简单
- 多模块并行开发产生接口漂移

------

# 十五、AGENTS.md

检查仓库根目录是否存在 `AGENTS.md`。

如果不存在，创建一个简洁版本。

如果存在，只在确有必要时更新，不覆盖已有有效规则。

`AGENTS.md` 应只存放长期有效的仓库规则，不要把全部 Spec 复制进去。

至少考虑写入：

```text
项目定位
主要语言和版本
目录约定
架构依赖规则
核心逻辑不得依赖 Agent SDK
测试命令
格式化命令
静态检查命令
安全规则
文档同步规则
Task 和 Requirement ID 规则
禁止绕过 Permission Engine
禁止 MCP Tool 绕过统一 Tool Runtime
禁止在 CLI 层实现核心业务逻辑
修改 Spec 时同步 Traceability Matrix
完成任务时更新 Checklist 和 Evidence
```

如果当前尚未确定构建命令，在 `AGENTS.md` 中明确标记为待建立，不要虚构命令。

------

# 十六、质量要求

所有文档使用简体中文。

以下内容使用英文：

- 代码标识符
- 接口名
- 类型名
- 文件名
- Module ID
- Requirement ID
- Task ID
- 状态枚举
- 固有技术术语

文档必须：

- 信息具体
- 可执行
- 可追踪
- 内部一致
- 不堆砌口号
- 不虚构当前仓库已有能力
- 区分“已存在”“计划实现”“候选设计”
- 明确 MVP 与后续版本
- 对不确定事项记录 Assumption 或 Open Question
- 不使用“后续完善”代替设计
- 不使用“高性能”“高可用”而不给验收方式
- 不把所有任务都设为 P0
- 不把所有模块都定义为 MVP 阻断项

Mermaid 图必须语法有效。

伪代码必须保持 Go 风格，但本阶段不要求可编译。

------

# 十七、一致性检查

生成文档后，必须执行一次完整的自检。

检查：

## 17.1 覆盖性

确认十三项核心能力全部被某个模块负责。

## 17.2 唯一所有权

确认每个核心实体只有一个主拥有模块。

## 17.3 依赖方向

确认不存在明显循环依赖。

## 17.4 接口一致性

确认同一接口、实体、事件在不同文档中的命名一致。

## 17.5 状态一致性

确认 Session、Agent、Task、Team、Worktree 状态定义一致。

## 17.6 Requirement 追踪

确认所有 P0 需求都有：

- 模块
- Spec
- Task
- Test/Eval
- Checklist

## 17.7 Task 可执行性

确认没有：

- “实现整个模块”这种超大任务
- 没有验收条件的任务
- 依赖不存在的任务
- 循环依赖任务
- 只有文档没有测试的核心任务

## 17.8 安全完整性

确认所有工具调用路径都经过：

```text
Validation
Permission
Hook
Execution
Audit
```

确认：

- MCP Tool 不绕过权限
- SubAgent 不绕过权限
- Team Member 不绕过权限
- Slash Command 不绕过权限
- Skill 不可扩大自身权限
- Hook 不可静默提升权限

## 17.9 恢复完整性

确认关键状态能够通过持久化事件或 Checkpoint 恢复。

## 17.10 MVP 可完成性

确认 MVP 能够在不依赖 Agent Teams 完整实现的情况下形成一个可演示闭环。

发现不一致时，直接修改对应文档，不要只在最终报告中指出。

------

# 十八、最终输出格式

完成文档后，在终端回复中给出：

## 18.1 仓库判断

说明：

- 当前仓库成熟度
- 已有可复用内容
- 缺失内容
- 采用了哪些假设

## 18.2 最终模块列表

使用表格：

| Module ID | Module | Responsibility | MVP  | Dependencies |
| --------- | ------ | -------------- | ---- | ------------ |
|           |        |                |      |              |

## 18.3 创建或修改的文件

列出所有文件及用途。

## 18.4 关键架构决策

列出最重要的 5—10 项决策。

## 18.5 关键路径

说明实现 MVP 的关键任务链。

## 18.6 可并行工作

说明哪些模块可以交给不同 Codex Agent 并行执行，并指出共享文件冲突风险。

## 18.7 最高风险

列出当前最重要的五个风险。

## 18.8 下一步

只给出下一阶段最合理的动作：

```text
从 MASTER_TASKS.md 中选择第一个 Ready 的 P0 Task
为该 Task 创建独立实现线程
实现前重新读取对应 Spec、ADR 和 Checklist
```

------

# 十九、执行约束

在本次任务中：

- 不实现生产代码。
- 不安装依赖。
- 不大规模重构已有代码。
- 不删除已有文件。
- 不执行危险命令。
- 不创建大量无内容的代码骨架。
- 不因为缺少少量信息而停止。
- 对缺失信息做合理假设，并记录到 `OPEN_QUESTIONS.md`。
- 不要求用户逐个确认模块。
- 自主完成模块划分和文档生成。
- 对现有仓库内容保持尊重，避免覆盖有效设计。
- 最终必须实际将文档写入仓库，而不是只在对话中描述。

现在开始：

1. 阅读仓库。
2. 输出内部执行计划。
3. 必要时启动只读 Subagent 分析。
4. 完成模块划分。
5. 创建全部 Spec、Task、Checklist 和配套架构文档。
6. 执行一致性检查。
7. 给出最终报告。