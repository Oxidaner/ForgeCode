# OPEN_QUESTIONS

记录架构规划阶段做出的假设与需后续验证的问题。每项含状态与影响模块。规划阶段对缺失信息采用合理假设并在此登记，不阻塞文档产出。

| ID | 问题 | 当前假设 | 影响模块 | 状态 | 决策方式 |
| --- | --- | --- | --- | --- | --- |
| Q1 | 目标 Go 版本 | 假设 Go 1.22+（泛型、`log/slog`、`errors.Join`） | 全体 | Open | 首个 Spike（FC-RT-000）确认 |
| Q2 | SQLite 驱动选型 | 候选 `modernc.org/sqlite`（纯 Go，免 CGO）；需确认 FTS5 支持 | session-store, memory-system | Open | Spike 验证 FTS5 |
| Q3 | Git 操作方式 | 候选 `git` CLI 包装（Worktree 兼容性好）vs `go-git`（Worktree 支持有限） | git-worktree | Open | V0.3 前 Spike |
| Q4 | Token 估算精度 | MVP 用启发式（字符/词近似 + 模型系数），非精确 tokenizer | context-manager | Open | 误差超 15% 时引入 tokenizer |
| Q5 | CLI 框架 | 候选 cobra；是否需要 TUI 待定 | cli | Open | MVP 评审 |
| Q6 | Sandbox 默认开关 | MVP 默认关闭，Bash 本地受限执行；V0.2 引入 Docker 可选开启 | sandbox, permission-engine | Assumed | ADR-0012 |
| Q7 | Provider 优先级 | MVP 实现 Mock + OpenAI；Anthropic、OpenAI-Compatible 次之 | model-provider | Assumed | 路线图 |
| Q8 | 审批交互通道 | MVP 经 CLI 同步交互；未来可经 Hook/HTTP | permission-engine, cli | Assumed | — |
| Q9 | 事件存储归档策略 | 假设按 Session 大小阈值归档冷事件，保留恢复所需子集 | session-store | Open | V0.2 性能评估 |
| Q10 | 记忆是否需语义检索 | 默认仅 FTS5；除非证明关键词检索不足 | memory-system | Assumed | ADR-0007 重访条件 |
| Q11 | Structured Output 实现方式 | 经 Tool Calling / JSON Schema 约束，而非各 Provider 私有 JSON mode | model-provider | Open | Contract Test |
| Q12 | Team 持久化粒度 | 假设 Task/Team 落 SQLite，成员运行态走事件 | agent-orchestration | Open | V1.0 设计 |
| Q13 | Hook 默认失败策略 | 安全敏感事件 fail-closed，通知类 fail-open | extension-system | Assumed | 安全评审 |
| Q14 | 多 Agent 预算分配 | 父预算切分给子，子超额不回溯父；递归深度默认 ≤ 3 | agent-orchestration | Assumed | — |
| Q15 | Eval 评分方式 | 标杆场景用规则断言 + 可选模型评审；避免纯主观 | evaluation | Open | V0.2 |
