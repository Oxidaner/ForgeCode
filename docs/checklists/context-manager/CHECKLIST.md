# context-manager Checklist

模块：`context-manager`。相关需求 FR-CONTEXT-001..006。RISK-009。

## Design Ready
- [ ] Layer 分层模型已定义（FC-CTX-001）
- [ ] Builder/Compactor/TokenEstimator 接口已定义
- [ ] Budget（Token/Cost/Reserved Output）模型已定义
- [ ] KeyFacts 保护清单已定义
- [ ] 压缩回滚（依赖 session-store Checkpoint）语义已定义
- [ ] 错误模型映射 GLOSSARY

## Implementation Ready
- [ ] 任务已拆分（分层/估算/截断/压缩/保护/回滚）
- [ ] Token 估算方式已决策（OPEN_QUESTIONS Q4）
- [ ] 截断默认参数（头尾行数/分页/阈值）已定义
- [ ] Golden 语料已准备

## Implementation Complete
- [ ] 组装遵守窗口 − Reserved Output（FR-CONTEXT-001/002）
- [ ] Bash 头尾/Grep 去重/ReadFile 分页/Observation 压缩生效（FR-CONTEXT-003）
- [ ] 自动 + `/compact` 触发，压缩前 Checkpoint（FR-CONTEXT-004）
- [ ] KeyFacts 压缩后完整保留（FR-CONTEXT-005）
- [ ] Verify 失败回滚 Checkpoint（FR-CONTEXT-006, RISK-009）
- [ ] KeyFacts 脱敏，不含密钥明文

## Test Complete
- [ ] 分层组装/预算 Unit Test
- [ ] 截断/KeyFacts Golden Test
- [ ] 压缩失败回滚 Failure Injection
- [ ] 与 runtime-core/session-store Integration
- [ ] `go test -race` 通过

## Documentation Complete
- [ ] SPEC 与实现一致
- [ ] TASKS 状态更新
- [ ] 配置示例（阈值/分页）更新
- [ ] 已知限制（估算精度）记录

## Release Ready
- [ ] P0 验收通过
- [ ] 无未处理 Critical 风险（RISK-009 缓解）
- [ ] 压缩前后 token 与回滚次数可观察
- [ ] 自动压缩 Demo 可复现
