# RISK_REGISTER

风险等级：Probability/Impact 取 Low/Medium/High。Status：Open/Mitigating/Accepted/Closed。

| Risk ID | Description | Category | Prob | Impact | Trigger | Mitigation | Contingency | Owner Module | Related Tasks | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| RISK-001 | 项目范围过大，无法在 12 周内闭环 | Scope | High | High | 里程碑持续延期 | 严格 MVP 边界；高级能力降级到 V0.2+；按依赖而非时间排期 | 砍 V1.0 Team，保 MVP+V0.2 | runtime-core | FC-RT-* | Mitigating |
| RISK-002 | 模块过度设计、提前抽象 | Design | Medium | Medium | 接口频繁返工 | 反模式清单评审；接口随需求演进；禁过度泛型 | 重构合并模块 | 全体 | — | Open |
| RISK-003 | Agent Loop 不稳定（死循环/卡死） | Runtime | Medium | High | 重复工具调用、同错误循环 | 最大轮次/工具数、循环检测、Deadline、状态机显式终止 | 强制中断+落事件 | runtime-core | FC-RT-004 | Mitigating |
| RISK-004 | Provider 间能力差异导致行为不一致 | Provider | High | Medium | 新增 Provider 行为异常 | 统一接口+能力元数据+Contract Test；只统一公共子集 | 标记不支持能力 | model-provider | FC-PROV-* | Mitigating |
| RISK-005 | Tool Schema 在内置与 MCP 间不一致 | Tool | Medium | Medium | MCP schema 转换失真 | 统一 Descriptor + 转换层 + Contract Test | 拒绝不兼容工具 | tool-runtime | FC-TOOL-004 | Open |
| RISK-006 | Bash 安全分析误判（漏判/过判） | Security | High | High | 危险命令放行或安全命令被拒 | 结构化分析 + 测试语料库 + 默认保守 + 审批兜底 | 收紧默认 Deny | permission-engine | FC-PERM-005 | Mitigating |
| RISK-007 | 路径逃逸（traversal/symlink） | Security | Medium | High | 越界读写 | 规范化+真实路径解析+Workspace 边界+安全测试 | 全局只读降级 | permission-engine | FC-PERM-003 | Mitigating |
| RISK-008 | 审批绕过 | Security | Low | High | 高危操作未经审批执行 | 决策与执行分离；统一管线；恢复后重判 | 紧急停用工具 | permission-engine | FC-PERM-007 | Open |
| RISK-009 | 上下文压缩丢失关键信息 | Context | High | High | 压缩后目标/计划/文件丢失 | 关键事实保护清单 + Golden Test + 压缩前 Checkpoint | 回滚 Checkpoint | context-manager | FC-CTX-004/005 | Mitigating |
| RISK-010 | Memory 污染 | Memory | Medium | High | 错误/恶意记忆被复用 | 候选审批、置信度、过期、项目隔离、不自动写入 | 清除受污染记忆 | memory-system | FC-MEM-* | Open |
| RISK-011 | MCP Server 不可信 | Security | High | High | 恶意工具/资源/Prompt | 信任级别、Namespace、输出限制、独立权限+审计 | 隔离/禁用 Server | mcp-client | FC-MCP-004 | Open |
| RISK-012 | Hook 执行任意代码 | Security | Medium | High | 恶意 Shell/HTTP Hook | 来源追踪、Timeout、不可提权、失败策略、审计 | 禁用外部 Hook | extension-system | FC-HOOK-003 | Open |
| RISK-013 | SubAgent Token 成本爆炸 | Cost | Medium | High | 大量并发子 Agent | 父子预算切分、并发上限、递归深度限制、Cost Budget | 终止超预算子树 | agent-orchestration | FC-SUB-004 | Open |
| RISK-014 | Worktree 孤儿资源 | Reliability | Medium | Medium | 崩溃/取消后残留 | 登记表 + 启动时孤儿扫描清理 + 审计 | 手动清理命令 | git-worktree | FC-WT-003 | Open |
| RISK-015 | Agent Team 调度复杂度过高 | Orchestration | High | Medium | DAG 死锁/饥饿 | 中心化调度、状态机、依赖校验、超时、人工介入 | 降级为串行 SubAgent | agent-orchestration | FC-TEAM-* | Open |
| RISK-016 | 事件存储无限增长 | Reliability | Medium | Medium | 长 Session 体积膨胀 | 归档/截断策略，保留恢复子集，输出截断 | 归档冷数据 | session-store | FC-SESS-004 | Open |
| RISK-017 | SQLite 并发写限制 | Persistence | Medium | Medium | 多 Agent 并发写冲突 | 单写者/WAL、写串行化、按 Session 分片 | 降低并发 | session-store | FC-SESS-001 | Open |
| RISK-018 | Eval 不足以证明能力 | Quality | Medium | Medium | Demo 说服力不足 | 三个真实标杆场景 + 规则断言 + Replay | 增加场景 | evaluation | FC-EVAL-003 | Open |
| RISK-019 | Demo 场景过于简单 | Quality | Medium | Medium | 展示价值低 | 选取有难度的 PR/SQL/K8s 真实案例 | 升级案例 | evaluation | FC-EVAL-003 | Open |
| RISK-020 | 多模块并行开发接口漂移 | Process | High | Medium | 并行 Agent 改动接口不一致 | 接口契约先行（本规划）+ Contract Test + 单一 Owner 写核心文件 | 接口冻结评审 | 全体 | FC-RT-001 等架构任务 | Mitigating |
