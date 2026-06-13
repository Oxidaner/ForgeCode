# <Module Name> Spec

> 复制本模板创建 `docs/specs/<module-id>/SPEC.md`。所有标题为英文技术术语，正文使用简体中文。

## 1. Module Info

| 字段 | 值 |
| --- | --- |
| Module ID | `<module-id>` |
| Module Name | <Module Name> |
| Status | Draft / Reviewed / Stable |
| Owner | <负责人或占位> |
| Dependencies | `<module-id>`, ... |
| Dependents | `<module-id>`, ... |
| Related Requirements | FR-XXX-001, NFR-XXX-001, ... |
| Related ADRs | ADR-00XX |
| MVP | Yes / No / Partial |

## 2. Purpose
模块为什么存在，解决什么问题。

## 3. Scope
模块负责什么。

## 4. Non-goals
模块明确不负责什么（区别于全局非目标）。

## 5. Responsibilities
- 职责 1
- 职责 2

## 6. Public Interfaces
使用 Go 风格伪代码定义关键接口（本阶段不要求可编译，不过度抽象）。

```go
// 关键接口
```

## 7. Domain Model
主要实体、值对象、枚举。标注每个实体的拥有模块（应与 `DATA_OWNERSHIP.md` 一致）。

## 8. State Machine
存在生命周期的模块给出状态枚举与合法转移（Mermaid `stateDiagram-v2`）。

## 9. Core Flows
正常流程与异常流程（可用 Mermaid `sequenceDiagram`）。

## 10. Configuration

| Key | 默认值 | 作用域 | 敏感 | 说明 |
| --- | --- | --- | --- | --- |

## 11. Persistence
是否持久化、由谁拥有、存储介质、迁移策略。

## 12. Concurrency
线程安全性、锁/Channel 边界、并发上限、取消传播、Race 风险、顺序与幂等要求。

## 13. Error Model
列出错误类别（引用 `GLOSSARY.md` 的标准错误分类），说明触发条件与调用方处理方式。

## 14. Security
信任边界、攻击面、权限检查点、敏感数据、日志脱敏、Prompt/Tool Output Injection 风险。

## 15. Observability
Log / Metric / Trace / Audit Event / Usage / Cost。

## 16. Testing Strategy
Unit / Integration / Contract / Race / Failure Injection / Security / Golden / Eval。

## 17. Acceptance Criteria
可检查、可测试的条件（与 Checklist、Traceability 对应）。

## 18. Risks
模块实现风险（引用 `RISK_REGISTER.md` 的 Risk ID）。

## 19. Open Questions
真正需要后续验证的问题（同步到 `OPEN_QUESTIONS.md`）。
