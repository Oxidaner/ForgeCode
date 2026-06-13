# agent-orchestration Checklist

模块：`agent-orchestration`。相关需求 FR-SUBAGENT-001..004, FR-TEAM-001..003。RISK-013/015。

## Design Ready
- [ ] AgentDefinition/SubAgentSpec/SubAgentResult 已定义（FC-SUB-001）
- [ ] Delegator/TeamLead/Mailbox/ArtifactStore 接口已定义
- [ ] Task State / Team State 与 GLOSSARY 一致
- [ ] SubAgent（一次性）与 Team（长期+DAG）清晰区分
- [ ] 中心化调度（ADR-0009）、Worktree 隔离（ADR-0008）已对齐
- [ ] AgentInstance 运行态属 runtime-core，本模块只读调度

## Implementation Ready
- [ ] 任务已拆分（SubAgent 系列 + Team 系列）
- [ ] 并发上限/递归深度/预算切分策略已定义
- [ ] git-worktree 隔离接口已约定
- [ ] Team 持久化粒度已决策（OPEN_QUESTIONS Q12）

## Implementation Complete
- [ ] SubAgent 独立身份/上下文/预算/白名单，不继承全部父上下文（FR-SUBAGENT-001/003）
- [ ] 默认返回结构化摘要而非完整日志（FR-SUBAGENT-003）
- [ ] 并发受 limit，递归深度受限（FR-SUBAGENT-004）
- [ ] 子/成员工具调用经统一管线 + 权限（无绕过，NFR-SEC-001）
- [ ] 父子/团队预算统计正确，超限终止子树（RISK-013）
- [ ] Team Lead 按 DAG 依赖调度并集成结果（FR-TEAM-001/002）

## Test Complete
- [ ] DAG 调度/预算 Unit Test
- [ ] 委派/Team Integration
- [ ] 并发委派 Race Test
- [ ] 子超时/失败/取消、成员失败阻塞 Failure Injection
- [ ] 子绕过权限尝试被拒 Security Test
- [ ] `go test -race` 通过

## Documentation Complete
- [ ] SPEC 与实现一致
- [ ] TASKS 状态更新
- [ ] ADR-0008/0009 与实现一致
- [ ] SubAgent/Team 配置示例

## Release Ready
- [ ] SubAgent（V0.3）/Team（V1.0）验收通过
- [ ] 无未处理 Critical 风险（RISK-013/015 缓解）
- [ ] 并发/预算/任务吞吐指标可观察
- [ ] 委派与 Team Demo 可复现
