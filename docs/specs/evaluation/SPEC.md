# evaluation Spec

## 1. Module Info

| 字段 | 值 |
| --- | --- |
| Module ID | `evaluation` |
| Module Name | Evaluation |
| Status | Draft |
| Owner | 架构组（占位） |
| Dependencies | runtime-core, session-store |
| Dependents | cli, ci, extension-system（候选 Skill 回放） |
| Related Requirements | FR-EVAL-001..004 |
| Related ADRs | ADR-0010 |
| MVP | No（V0.2+） |

## 2. Purpose
evaluation 提供 Replay 与 Eval 能力：基于持久化事件重放 Session，运行固定 Eval Case 评分，并实现三个标杆场景（review-pr/sql/k8s）的评测，以及候选 Skill 的回放评测（支撑 ADR-0010）。它是证明系统能力与防止回归的主要手段。

## 3. Scope
- Replay：从持久化事件重放 Session（不重新执行外部副作用）。
- Eval Case 框架：固定输入 + 期望结果 + 评分（规则断言 + 可选模型评审）。
- 三个标杆场景 Eval：review-pr、review-sql、review-k8s。
- 候选 Skill 回放评测（FR-SKILL-004）。

## 4. Non-goals
- 不实现 Agent Loop（runtime-core）。
- 不实现 Skill 安装（extension-system）。
- 不做生产监控（telemetry）。
- 不以纯主观判断作为唯一评分（OPEN_QUESTIONS Q15）。

## 5. Responsibilities
- 拥有 Eval Case 定义与结果。
- 提供 Replayer：读 session-store 事件重放。
- 提供 EvalRunner + Scorer：执行用例并评分。
- 定义三个标杆场景的检查项与断言。
- 为 extension-system 候选 Skill 提供回放评测入口。

## 6. Public Interfaces

```go
type Replayer interface {
    Replay(ctx context.Context, sessionID string) (ReplayResult, error) // 不重执行副作用
}

type EvalRunner interface {
    Run(ctx context.Context, suite EvalSuite) (EvalReport, error)
}

type EvalCase struct {
    ID        string
    Input     TaskInput
    Expected  Expectation
    Assertions []Assertion // 规则断言
    Judge     *JudgeSpec   // 可选模型评审
}

type Scorer interface {
    Score(case EvalCase, actual Outcome) (Score, []AssertionResult)
}
```

## 7. Domain Model
- `EvalCase`、`EvalSuite`、`EvalReport`、`Assertion`、`AssertionResult`、`Score`、`JudgeSpec`、`ReplayResult`。
- 本模块拥有 Eval 定义与结果。

## 8. State Machine
Eval 运行（非持久）：`Pending → Running → Scored → Reported`。Replay：`Loading → Replaying → Done/Failed`。

## 9. Core Flows
- **Replay**：读 session-store 事件流 → 按序重放重建轨迹（使用已记录 ToolResult，不重执行）→ 输出可比对轨迹。
- **Eval**：加载 EvalSuite → 对每个 Case 用 Mock/真实 Provider 运行 → Scorer 规则断言 + 可选模型评审 → EvalReport。
- **标杆场景**：构造 PR/SQL/K8s 输入 → 运行对应 Skill/Agent → 断言检查项命中。
- **候选 Skill 回放**：对候选 Skill 在历史轨迹上回放评测，产出通过/拒绝建议（人工审批前置）。

## 10. 标杆场景检查项

### review-pr
变更范围、调用链影响、风险等级、证据、文件与行号、测试结果、建议修改、阻断项。

### review-sql
DDL 风险、全表锁、大表修改、缺失索引、无条件 UPDATE、无条件 DELETE、数据兼容、回滚方案、发布顺序。

### review-k8s
镜像变更、Resource Requests/Limits、探针、Service、Ingress、RBAC、Secret、ConfigMap、滚动更新、高可用、Dry-run、Schema 校验、回滚能力。

## 11. Configuration

| Key | 默认值 | 作用域 | 敏感 | 说明 |
| --- | --- | --- | --- | --- |
| `eval.provider` | mock | 全局 | 否 | Eval 用 Provider |
| `eval.judge_enabled` | false | 全局 | 否 | 是否启用模型评审 |
| `eval.suite_paths` | `eval/` | 全局 | 否 | 用例路径 |

## 12. Persistence
Eval 结果可落 SQLite/文件。Replay 只读 session-store 事件。

## 13. Concurrency
- 多 Case 可并发运行（受 Provider 限流）。
- Replay 只读、无副作用。
- 取消经 context 传播。

## 14. Error Model
`RecoveryError`（Replay 重放失败）、`ProviderError`（Eval 运行）、`ValidationError`（用例定义非法）、`TimeoutError`。

## 15. Security
- Replay 不重执行外部副作用（与恢复同语义）。
- 候选 Skill 回放是 ADR-0010 安全门的一环（防恶意 Skill）。
- Eval 用例不含真实密钥（占位）。

## 16. Observability
- 指标：Eval 通过率、断言命中、回放成功率、场景得分。
- 报告：EvalReport（可导出）。

## 17. Testing Strategy
- Unit：Scorer 断言、用例加载。
- Integration：Replay 真实 Session、标杆场景端到端。
- Golden：标杆场景期望输出快照。
- 自测：evaluation 自身用例。

## 18. Acceptance Criteria
- [ ] Replay 重放不触发外部副作用。
- [ ] Eval Case 框架支持规则断言 + 可选模型评审。
- [ ] 三个标杆场景检查项均有断言并产出得分。
- [ ] 候选 Skill 回放评测可产出通过/拒绝建议（FR-SKILL-004）。
- [ ] EvalReport 可导出供 CI 与展示。

## 19. Risks
RISK-018（Eval 不足）、RISK-019（Demo 过简）。

## 20. Open Questions
- 评分方式：规则 vs 模型评审权重（Q15）。
- 标杆场景的真实案例来源与难度选取。
- 模型评审的稳定性与成本。
